// Package harnessx is the single choke-point every SWE-AF role reasoner uses to
// call the AgentField harness. It replaces the Python monkeypatch of
// app.harness (swe_af/app.py:80-93) with an explicit generic wrapper:
//
//   - schemaFor[T] reflects a Go struct into the JSON-schema map the harness
//     consumes to build its OUTPUT REQUIREMENTS prompt suffix (design §2.3).
//   - Run[T] injects the build's run-scoped credentials into the subprocess env
//     (scoped creds win over the base, mirroring Python precedence), calls the
//     harness, classifies fatal API errors, and on a schema parse-failure hands
//     the caller a default-seeded value so it can apply its deterministic
//     fallback (design §4.1).
//
// This is the ONLY way roles should reach the harness — it guarantees uniform
// credential injection and fatal-error handling across all 22 role reasoners.
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

	"github.com/Agent-Field/agentfield/sdk/go/harness"
	invjsonschema "github.com/invopop/jsonschema"
	tekjsonschema "github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/Agent-Field/SWE-AF/go/internal/fatal"
	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
)

// schemaCache memoizes the reflected JSON-schema map per concrete type T so the
// (non-trivial) reflection + marshal round-trip runs once per type. Keyed by
// reflect.Type; the stored map[string]any is treated as immutable by callers
// (the harness only ever marshals/reads it, never mutates), so sharing the
// cached value across goroutines is safe.
var schemaCache sync.Map // reflect.Type -> map[string]any

// schemaFor reflects T into the JSON-schema map the Go SDK harness consumes.
//
// How the SDK consumes this map (verified against sdk/go/harness/schema.go):
// the map is NOT used for programmatic validation — validity is defined purely
// by json.Unmarshal into the dest struct succeeding. The harness uses the map
// only to (a) embed a pretty-printed schema in the BuildPromptSuffix /
// BuildFollowupPrompt OUTPUT REQUIREMENTS instruction (harness/schema.go:36-74,
// :322-354) and (b) list expected top-level keys in DiagnoseOutputFailure
// (schema.go:306-318, which reads map["properties"]). Keys are alphabetized by
// json.MarshalIndent, so field ordering and `required` completeness are
// cosmetic. This means invopop's output — with $defs, items, and enum — is more
// than sufficient, and far richer than the SDK's own shallow StructToJSONSchema
// (which drops nested props/items/enums; design §2.3 says do NOT use it).
//
// Reflector configuration:
//   - ExpandedStruct: inline the root type's own properties at the top level so
//     map["properties"] is populated for DiagnoseOutputFailure (rather than a
//     bare $ref to $defs).
//   - DoNotReference=false (default): emit a $defs map for nested struct types.
//   - Anonymous: suppress the auto-generated $id derived from the package path.
func schemaFor[T any]() map[string]any {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if cached, ok := schemaCache.Load(t); ok {
		return cached.(map[string]any)
	}

	r := &invjsonschema.Reflector{
		ExpandedStruct: true,  // root properties inline at top level
		DoNotReference: false, // emit $defs for nested types
		Anonymous:      true,  // no auto-generated $id from PkgPath
	}
	schema := r.ReflectFromType(t)

	b, err := json.Marshal(schema)
	if err != nil {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]any{}
	}

	// Make the reflected schema faithful to the source Pydantic model: the SDK
	// now validates harness output against this map, so invopop's all-fields
	// `required`, `additionalProperties:false`, and non-null Optional subschemas
	// would OVER-reject valid output. (Enums are added on the enum types via
	// JSONSchemaExtend, already present in `m`.)
	m = schemas.MakePydanticFaithful(m, t.Name())

	schemaCache.Store(t, m)
	return m
}

// executeStructured owns SWE's structured-output reliability policy around the
// upstream AgentField harness call. Run[T] intentionally stays a thin
// integration seam so recovery behavior can evolve without spreading changes
// across role code or increasing rebase pressure on the base call path.
func executeStructured[T any](ctx context.Context, app HarnessCaller, prompt string, schema map[string]any, opts harness.Options) (*T, *harness.Result, error) {
	// Weak OpenCode-backed models are materially more reliable when the SDK
	// builds the structured envelope incrementally. Keep that policy here, next
	// to recovery/validation, so role code and the base Run seam do not know how
	// structured output is repaired.
	if opts.Provider == "opencode" && opts.SchemaMode == "" {
		opts.SchemaMode = "incremental"
	}

	var dest T
	stopOutputCapture := startStructuredOutputCapture[T](ctx, opts.ProjectDir, schema)
	result, err := app.Harness(ctx, prompt, schema, &dest, opts)
	capturedOutput := stopOutputCapture()
	if err != nil {
		// A no-progress watchdog is a liveness signal, not proof that an already
		// completed exact-schema result is unusable. Generic/provider errors remain
		// fail-closed.
		if strings.Contains(err.Error(), "CLI command made no progress") {
			if recovered, ok := recoverStructuredResult[T](result, schema); ok {
				return recovered, result, nil
			}
			if len(capturedOutput) > 0 {
				var recovered T
				if recoverErr := recoverStructuredText(string(capturedOutput), schema, &recovered); recoverErr == nil {
					schemas.EmptyForNilSlices(&recovered)
					result = normalizeRecoveredResult(result, &recovered)
					return &recovered, result, nil
				}
			}
		}
		return nil, result, err
	}

	// Fatal API errors outrank schema fallback so billing/auth/provider failures
	// cannot be hidden behind a default-seeded result.
	if fErr := fatal.CheckFatalHarnessError(result); fErr != nil {
		return nil, result, fErr
	}

	if result != nil && result.Parsed == nil &&
		(result.FailureType == harness.FailureSchema || result.FailureType == harness.FailureNoOutput || result.FailureType == harness.FailureNone) {
		if recovered, ok := recoverStructuredResult[T](result, schema); ok {
			return recovered, result, nil
		}
	}

	if result == nil || result.Parsed == nil {
		seeded := seedDefaults[T]()
		schemas.EmptyForNilSlices(&seeded)
		return &seeded, result, nil
	}

	schemas.EmptyForNilSlices(&dest)
	return &dest, result, nil
}

func recoverStructuredResult[T any](result *harness.Result, schema map[string]any) (*T, bool) {
	if result == nil {
		return nil, false
	}
	for _, text := range structuredResultCandidates(result) {
		var recovered T
		if err := recoverStructuredText(text, schema, &recovered); err != nil {
			continue
		}
		schemas.EmptyForNilSlices(&recovered)
		normalizeRecoveredResult(result, &recovered)
		return &recovered, true
	}
	return nil, false
}

func normalizeRecoveredResult[T any](result *harness.Result, recovered *T) *harness.Result {
	if result == nil {
		result = &harness.Result{}
	}
	result.Parsed = recovered
	result.IsError = false
	result.ErrorMessage = ""
	result.FailureType = harness.FailureNone
	return result
}

// structuredResultCandidates returns only assistant text surfaces that can
// legitimately contain the final structured result. Tool outputs are excluded
// so repository/file JSON cannot be mistaken for the orchestration envelope.
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
// returning a CLI no-progress error, so this narrow monitor lets SWE salvage a
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
// unmarshals into dest. Malformed or schema-invalid text stays a failure.
func recoverStructuredText[T any](text string, schema map[string]any, dest *T) error {
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("marshal schema: %w", err)
	}
	compiler := tekjsonschema.NewCompiler()
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
		// Let its custom UnmarshalJSON normalize its one known weak-model shape,
		// then validate the fully materialized typed object. Other schemas stay
		// strict; required fields are never synthesized generically.
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
// ignoring braces inside quoted strings. Larger candidates win over nested
// examples so the orchestration envelope is preferred.
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

// seedDefaults returns a T seeded with its pydantic-parity defaults. Custom
// UnmarshalJSON implementations can materialize non-zero defaults; plain
// structs simply retain their Go zero value.
func seedDefaults[T any]() T {
	var v T
	_ = json.Unmarshal([]byte("{}"), &v)
	return v
}
