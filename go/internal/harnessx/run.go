package harnessx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Agent-Field/agentfield/sdk/go/agent"
	"github.com/Agent-Field/agentfield/sdk/go/harness"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/Agent-Field/SWE-AF/go/internal/fatal"
	"github.com/Agent-Field/SWE-AF/go/internal/hitl"
	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
)

// runIDFromContext extracts the build's run ID from the reasoner execution
// context. In production this reads agent.ExecutionContextFrom(ctx).RunID, which
// the SDK populates on the handler's ctx before dispatch. It is a package var
// (not a direct call) purely so tests can inject a run ID — the SDK's context
// key is unexported, so there is no public way to seed ExecutionContext into a
// ctx from an external package.
var runIDFromContext = func(ctx context.Context) string {
	return agent.ExecutionContextFrom(ctx).RunID
}

// HarnessCaller is the minimal method set Run needs from *agent.Agent. Declaring
// it as an interface (rather than depending on the concrete *agent.Agent) lets
// tests supply a mock harness without a live subprocess — the same seam the
// Python tests get by patching router.harness. *agent.Agent satisfies it via
// its Harness method (sdk/go/agent/harness.go:84).
type HarnessCaller interface {
	Harness(ctx context.Context, prompt string, schema map[string]any, dest any, opts harness.Options) (*harness.Result, error)
}

// Run is the single generic entry point every role reasoner uses to invoke the
// harness for structured output of type T.
//
// Sequence (design §4.1, §2.3, §4.3):
//  1. Reflect T into the JSON schema the harness consumes (cached per type).
//  2. Inject the build's run-scoped credentials into opts.Env, scoped creds
//     overriding the base env — mirroring the Python precedence where a freshly
//     minted scout token beats a stale value inherited from os.environ. The run
//     ID comes from agent.ExecutionContextFrom(ctx).RunID.
//  3. Call app.Harness with a fresh *T dest.
//  4. Classify fatal (non-retryable) API errors FIRST, before the Parsed==nil
//     fallback, so the real billing/auth message surfaces past every retry layer
//     as a *fatal.FatalHarnessError (callers must propagate it, not swallow it).
//  5. On Result.Parsed == nil (the harness could not parse valid JSON into T),
//     return a default-seeded T plus the Result — NOT an error — so the caller
//     inspects Result.IsError and applies its role-specific deterministic
//     fallback. The seed comes from unmarshaling "{}" into T, which triggers
//     T's UnmarshalJSON default-seeding (schemas §2.2) when present, and yields
//     the Go zero value otherwise.
//
// Returns (*T, *harness.Result, error). The Result is returned even alongside a
// non-nil error so callers can inspect diagnostics.
func Run[T any](ctx context.Context, app HarnessCaller, prompt string, opts harness.Options) (*T, *harness.Result, error) {
	schema := schemaFor[T]()

	runID := runIDFromContext(ctx)
	opts.Env = hitl.InjectCredentialsIntoEnv(opts.Env, runID)

	var dest T
	stopOutputCapture := startStructuredOutputCapture[T](ctx, opts.ProjectDir, schema)
	result, err := app.Harness(ctx, prompt, schema, &dest, opts)
	capturedOutput := stopOutputCapture()
	if err != nil {
		// A provider CLI can stall after already emitting a complete structured
		// result. The no-progress watchdog is a liveness signal, not proof that
		// the result is unusable. Recover only this narrow case and only after
		// exact schema validation; generic/provider errors remain fatal.
		if strings.Contains(err.Error(), "CLI command made no progress") {
			if result != nil {
				for _, text := range structuredResultCandidates(result) {
					var recovered T
					if recoverErr := recoverStructuredText(text, schema, &recovered); recoverErr != nil {
						continue
					}
					schemas.EmptyForNilSlices(&recovered)
					result.Parsed = &recovered
					result.IsError = false
					result.ErrorMessage = ""
					result.FailureType = harness.FailureNone
					return &recovered, result, nil
				}
			}
			if len(capturedOutput) > 0 {
				var recovered T
				if recoverErr := recoverStructuredText(string(capturedOutput), schema, &recovered); recoverErr == nil {
					schemas.EmptyForNilSlices(&recovered)
					if result == nil {
						result = &harness.Result{}
					}
					result.Parsed = &recovered
					result.IsError = false
					result.ErrorMessage = ""
					result.FailureType = harness.FailureNone
					return &recovered, result, nil
				}
			}
		}
		return nil, result, err
	}

	// Fatal-error classification comes before the Parsed==nil fallback so the
	// real non-retryable message is not masked by a generic fallback struct.
	if fErr := fatal.CheckFatalHarnessError(result); fErr != nil {
		return nil, result, fErr
	}

	// Recover one important weak-model failure mode before giving up: OpenCode may
	// preserve a useful schema-shaped JSON object in an earlier text event even
	// when the file protocol fails and later schema retries produce unrelated
	// text. The SDK returns only the LAST retry text in Result, but it preserves
	// all attempt events in Messages. Search only assistant text events (never tool
	// output) and recover only schema/no-output failures; provider/transport errors
	// remain fail-closed.
	if result != nil && result.Parsed == nil &&
		(result.FailureType == harness.FailureSchema || result.FailureType == harness.FailureNoOutput || result.FailureType == harness.FailureNone) {
		for _, text := range structuredResultCandidates(result) {
			var recovered T
			if recoverErr := recoverStructuredText(text, schema, &recovered); recoverErr != nil {
				continue
			}
			schemas.EmptyForNilSlices(&recovered)
			result.Parsed = &recovered
			result.IsError = false
			result.ErrorMessage = ""
			result.FailureType = harness.FailureNone
			return &recovered, result, nil
		}
	}

	// Schema parse failure: hand the caller a default-seeded value plus the
	// Result so it can apply its own deterministic fallback. Not an error.
	if result == nil || result.Parsed == nil {
		seeded := seedDefaults[T]()
		schemas.EmptyForNilSlices(&seeded)
		return &seeded, result, nil
	}

	// Normalise nil slices to empty slices on the parsed result so every role
	// result serialises its list fields as [] (pydantic parity) — matching the
	// checkpoint/DAGState normalisation. Maps and nil pointers are left null.
	schemas.EmptyForNilSlices(&dest)
	return &dest, result, nil
}

// structuredResultCandidates returns only assistant text surfaces that can
// legitimately contain the final structured result. Result is the latest retry
// text; Messages retains earlier OpenCode attempts, which is critical when a
// useful first answer was followed by broken schema retries. Tool outputs are
// deliberately excluded so repository/file JSON cannot be mistaken for the
// orchestration result.
func structuredResultCandidates(result *harness.Result) []string {
	if result == nil {
		return nil
	}
	seen := make(map[string]struct{})
	out := make([]string, 0, 4)
	appendUnique := func(text string) {
		if text == "" {
			return
		}
		if _, ok := seen[text]; ok {
			return
		}
		seen[text] = struct{}{}
		out = append(out, text)
	}

	appendUnique(result.Result)
	for i := len(result.Messages) - 1; i >= 0; i-- {
		msg := result.Messages[i]
		kind, _ := msg["type"].(string)
		if kind != "text" && kind != "assistant" && kind != "result" {
			continue
		}
		if text, ok := msg["text"].(string); ok {
			appendUnique(text)
		}
		if part, ok := msg["part"].(map[string]any); ok {
			if text, ok := part["text"].(string); ok {
				appendUnique(text)
			}
		}
	}
	return out
}

// startStructuredOutputCapture watches only output directories created after
// this invocation starts and keeps the latest exact-schema-valid output bytes
// in memory. The AgentField SDK removes its temporary output directory before
// returning a CLI no-progress error, so this narrow monitor lets Run salvage a
// completed result without weakening validation. If more than one new output
// directory appears, capture becomes ambiguous and fails closed.
func startStructuredOutputCapture[T any](ctx context.Context, projectDir string, schema map[string]any) func() []byte {
	if projectDir == "" || schema == nil {
		return func() []byte { return nil }
	}

	pattern := filepath.Join(projectDir, ".agentfield-out-*")
	existingDirs, _ := filepath.Glob(pattern)
	existing := make(map[string]struct{}, len(existingDirs))
	for _, dir := range existingDirs {
		existing[dir] = struct{}{}
	}

	watchCtx, cancel := context.WithCancel(ctx)
	var mu sync.Mutex
	var latest []byte
	ambiguous := false
	done := make(chan struct{})

	go func() {
		defer close(done)
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()

		scan := func() {
			dirs, _ := filepath.Glob(pattern)
			newDirs := make([]string, 0, 1)
			for _, dir := range dirs {
				if _, ok := existing[dir]; !ok {
					newDirs = append(newDirs, dir)
				}
			}
			if len(newDirs) > 1 {
				mu.Lock()
				ambiguous = true
				latest = nil
				mu.Unlock()
				return
			}
			if len(newDirs) != 1 {
				return
			}

			mu.Lock()
			isAmbiguous := ambiguous
			mu.Unlock()
			if isAmbiguous {
				return
			}

			b, err := os.ReadFile(filepath.Join(newDirs[0], ".agentfield_output.json"))
			if err != nil || len(b) == 0 {
				return
			}
			var recovered T
			if err := recoverStructuredText(string(b), schema, &recovered); err != nil {
				return
			}
			mu.Lock()
			latest = append(latest[:0], b...)
			mu.Unlock()
		}

		for {
			scan()
			select {
			case <-watchCtx.Done():
				return
			case <-ticker.C:
			}
		}
	}()

	return func() []byte {
		cancel()
		<-done
		mu.Lock()
		defer mu.Unlock()
		if ambiguous || len(latest) == 0 {
			return nil
		}
		return append([]byte(nil), latest...)
	}
}

// recoverStructuredText extracts candidate JSON objects with a string-aware
// scanner, validates them against the exact generated schema, and only then
// unmarshals into dest. It is intentionally strict: malformed or schema-invalid
// text stays a failure and is handled by the normal retry/error path.
func recoverStructuredText[T any](text string, schema map[string]any, dest *T) error {
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("marshal schema: %w", err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("mem://swe/schema.json", bytes.NewReader(schemaBytes)); err != nil {
		return fmt.Errorf("add schema: %w", err)
	}
	compiled, err := compiler.Compile("mem://swe/schema.json")
	if err != nil {
		return fmt.Errorf("compile schema: %w", err)
	}

	for _, candidate := range extractJSONObjectCandidates(text) {
		var data any
		if err := json.Unmarshal([]byte(candidate), &data); err != nil {
			continue
		}
		if err := compiled.Validate(data); err == nil {
			if err := json.Unmarshal([]byte(candidate), dest); err == nil {
				return nil
			}
		}

		// CoderResult mirrors a Pydantic model where every field has a default.
		// Weak models commonly omit those defaulted fields and sometimes emit the
		// advisory agent_retro as a string. Let its custom UnmarshalJSON normalize
		// that one known shape, then validate the fully-materialized typed object.
		// Do not do this generically: required fields on other schemas must never be
		// synthesized from Go zero values.
		if reflect.TypeOf((*T)(nil)).Elem().Name() == "CoderResult" {
			var normalized T
			if err := json.Unmarshal([]byte(candidate), &normalized); err != nil {
				continue
			}
			normalizedBytes, err := json.Marshal(normalized)
			if err != nil {
				continue
			}
			var normalizedData any
			if err := json.Unmarshal(normalizedBytes, &normalizedData); err != nil {
				continue
			}
			if err := compiled.Validate(normalizedData); err != nil {
				continue
			}
			*dest = normalized
			return nil
		}
	}
	return fmt.Errorf("no schema-valid JSON object found in final text")
}

// extractJSONObjectCandidates finds balanced top-level JSON objects while
// ignoring braces that appear inside quoted strings. Candidates are returned
// largest-first so an enclosing result object wins over nested examples.
func extractJSONObjectCandidates(text string) []string {
	candidates := make([]string, 0, 2)
	depth := 0
	start := -1
	inString := false
	escaped := false

	for i := 0; i < len(text); i++ {
		ch := text[i]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			if depth == 0 {
				continue
			}
			depth--
			if depth == 0 && start >= 0 {
				candidates = append(candidates, text[start:i+1])
				start = -1
			}
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return len(candidates[i]) > len(candidates[j])
	})
	return candidates
}

// seedDefaults returns a T seeded with its pydantic-parity defaults. Unmarshaling
// an empty JSON object invokes T's UnmarshalJSON (which seeds non-zero defaults,
// schemas §2.2) when T implements it; for a plain struct it leaves the Go zero
// value. Any unmarshal error is ignored — the zero value is an acceptable floor.
func seedDefaults[T any]() T {
	var v T
	_ = json.Unmarshal([]byte("{}"), &v)
	return v
}

// RoleOptions is the role→harness parameter mapping (design §4.1). Each role
// reasoner fills it from its resolved config, then ToOptions produces the
// harness.Options passed to Run. Centralizing the mapping here keeps every role
// consistent — the field set mirrors the keyword arguments the Python role
// reasoners pass to router.harness (system_prompt, schema, model, provider,
// tools, cwd, max_turns, permission_mode).
type RoleOptions struct {
	// Provider is the harness ADAPTER string (not the provider), e.g. the output
	// of runtimex.RuntimeToHarnessAdapter — "claude-code", "opencode", "codex".
	// Python passes provider=runtime_to_harness_adapter(ai_provider).
	Provider string

	// Model is the resolved role model identifier.
	Model string

	// MaxTurns caps agent iterations (Python DEFAULT_AGENT_MAX_TURNS per role).
	MaxTurns int

	// Tools is the allowed-tool list (e.g. ["Read","Write","Glob","Grep","Bash"]).
	Tools []string

	// PermissionMode maps to Python's permission_mode; empty means the harness
	// default (Python passes `permission_mode or None`, and harness.Options
	// treats "" as "use default").
	PermissionMode string

	// SystemPrompt is the role's module-level system prompt.
	SystemPrompt string

	// Cwd is the working directory for the subprocess (repo path / worktree).
	Cwd string

	// Env is the base environment for the subprocess. Run overlays the build's
	// scoped credentials on top of this before invoking the harness.
	Env map[string]string
}

// ToOptions converts a RoleOptions into a harness.Options. Run injects scoped
// credentials into Env afterwards, so callers leave Env as the base env only.
func (r RoleOptions) ToOptions() harness.Options {
	binPath := ""
	schemaMode := ""
	if r.Provider == "opencode" {
		// SWE owns its harness binary contract. Cross-component fallbacks (for
		// example SEC_AF_OPENCODE_BIN) couple independently deployable agents and
		// belong in fleet orchestration, not SWE source.
		binPath = os.Getenv("SWE_OPENCODE_BIN")

		// Weak OpenCode-backed models are materially more reliable when the
		// structured envelope is built field-by-field. Incremental mode also keeps
		// the original task/schema context on validation retries, while the SDK's
		// single-shot retry prompt contains only the output-file diagnosis. That
		// distinction is critical when a model otherwise drifts into generic prose
		// after completing the coding work.
		schemaMode = "incremental"
	}
	return harness.Options{
		Provider:       r.Provider,
		BinPath:        binPath,
		Model:          r.Model,
		MaxTurns:       r.MaxTurns,
		Tools:          r.Tools,
		PermissionMode: r.PermissionMode,
		SystemPrompt:   r.SystemPrompt,
		Cwd:            r.Cwd,
		ProjectDir:     r.Cwd,
		Env:            r.Env,
		SchemaMode:     schemaMode,
	}
}
