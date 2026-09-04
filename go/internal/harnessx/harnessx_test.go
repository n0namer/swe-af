package harnessx

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/sdk/go/harness"

	"github.com/Agent-Field/SWE-AF/go/internal/fatal"
	"github.com/Agent-Field/SWE-AF/go/internal/hitl"
	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
)

// --- test fixtures ----------------------------------------------------------

// mockHarness is the HarnessCaller seam the Python tests get by patching
// router.harness. It records what Run passed and returns a scripted result.
type mockHarness struct {
	fn        func(ctx context.Context, prompt string, schema map[string]any, dest any, opts harness.Options) (*harness.Result, error)
	gotOpts   harness.Options
	gotSchema map[string]any
	gotPrompt string
}

func (m *mockHarness) Harness(ctx context.Context, prompt string, schema map[string]any, dest any, opts harness.Options) (*harness.Result, error) {
	m.gotOpts = opts
	m.gotSchema = schema
	m.gotPrompt = prompt
	return m.fn(ctx, prompt, schema, dest, opts)
}

// seededResult carries non-zero pydantic-parity defaults via UnmarshalJSON,
// exactly like the real schemas structs (schemas §2.2). Used to prove the
// Parsed==nil path returns seeded defaults, not the Go zero value.
type seededResult struct {
	Complete bool   `json:"complete"`
	Scope    string `json:"estimated_scope"`
}

func (s *seededResult) UnmarshalJSON(b []byte) error {
	*s = seededResult{Complete: true, Scope: "medium"}
	type alias seededResult
	return json.Unmarshal(b, (*alias)(s))
}

// childItem / parentSchema exercise nested $defs, array items, and enum output.
type childItem struct {
	Name string `json:"name"`
}

type parentSchema struct {
	Outcome  string      `json:"outcome" jsonschema:"enum=completed,enum=failed"`
	Children []childItem `json:"children"`
}

// --- schema generation ------------------------------------------------------

// Contract: schema for a nested struct emits $defs / items / enum for the SDK's
// consumption.
func TestSchemaForNestedStructEmitsDefsItemsEnum(t *testing.T) {
	m := schemaFor[parentSchema]()

	props, ok := m["properties"].(map[string]any)
	if !ok {
		t.Fatalf("expected top-level properties (ExpandedStruct), got: %v", m)
	}

	// enum on the outcome field.
	outcome, ok := props["outcome"].(map[string]any)
	if !ok {
		t.Fatalf("expected properties.outcome, got: %v", props)
	}
	enum, ok := outcome["enum"].([]any)
	if !ok {
		t.Fatalf("expected properties.outcome.enum, got: %v", outcome)
	}
	if !containsStr(enum, "completed") || !containsStr(enum, "failed") {
		t.Fatalf("enum missing expected values: %v", enum)
	}

	// array items on the children field.
	children, ok := props["children"].(map[string]any)
	if !ok {
		t.Fatalf("expected properties.children, got: %v", props)
	}
	if children["type"] != "array" {
		t.Fatalf("expected children.type=array, got: %v", children["type"])
	}
	if _, ok := children["items"]; !ok {
		t.Fatalf("expected children.items, got: %v", children)
	}

	// $defs for the nested struct type.
	defs, ok := m["$defs"].(map[string]any)
	if !ok || len(defs) == 0 {
		t.Fatalf("expected non-empty $defs for nested type, got: %v", m["$defs"])
	}
}

func TestSchemaForIsCached(t *testing.T) {
	a := schemaFor[parentSchema]()
	b := schemaFor[parentSchema]()
	// Same underlying cached map instance (pointer identity via a mutation probe).
	a["__probe__"] = 1
	if _, ok := b["__probe__"]; !ok {
		t.Fatalf("expected schemaFor to return the cached map instance on repeat calls")
	}
	delete(a, "__probe__")
}

// --- Run: fatal propagation -------------------------------------------------

// Contract: a fatal harness error propagates as *FatalHarnessError via errors.As.
func TestRunFatalErrorPropagates(t *testing.T) {
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
			return &harness.Result{IsError: true, ErrorMessage: "Credit balance is too low"}, nil
		},
	}

	out, res, err := Run[seededResult](context.Background(), mh, "prompt", harness.Options{})
	if err == nil {
		t.Fatal("expected an error for a fatal harness result")
	}
	var fe *fatal.FatalHarnessError
	if !errors.As(err, &fe) {
		t.Fatalf("expected *fatal.FatalHarnessError, got %T: %v", err, err)
	}
	if fe.OriginalMessage != "Credit balance is too low" {
		t.Fatalf("expected original message preserved, got %q", fe.OriginalMessage)
	}
	if out != nil {
		t.Fatalf("expected nil value on fatal error, got %v", out)
	}
	if res == nil {
		t.Fatal("expected the *harness.Result to be returned alongside the fatal error")
	}
}

// --- Run: scoped credential injection ---------------------------------------

// Contract: scoped creds appear in opts.Env overriding base env.
func TestRunInjectsScopedCredsOverridingBase(t *testing.T) {
	const runID = "run-scoped-123"
	hitl.StoreScopedCredentials(runID, map[string]string{"RAILWAY_TOKEN": "fresh"})
	defer hitl.ClearScopedCredentials(runID)

	restore := runIDFromContext
	runIDFromContext = func(context.Context) string { return runID }
	defer func() { runIDFromContext = restore }()

	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, dest any, _ harness.Options) (*harness.Result, error) {
			return &harness.Result{Parsed: dest}, nil
		},
	}

	base := harness.Options{Env: map[string]string{"RAILWAY_TOKEN": "stale", "KEEP": "yes"}}
	if _, _, err := Run[seededResult](context.Background(), mh, "prompt", base); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := mh.gotOpts.Env["RAILWAY_TOKEN"]; got != "fresh" {
		t.Fatalf("expected scoped cred to override base, got RAILWAY_TOKEN=%q", got)
	}
	if got := mh.gotOpts.Env["KEEP"]; got != "yes" {
		t.Fatalf("expected non-overlapping base env preserved, got KEEP=%q", got)
	}
	// Base map must not be mutated in place (InjectCredentialsIntoEnv returns a copy).
	if base.Env["RAILWAY_TOKEN"] != "stale" {
		t.Fatalf("base env was mutated: %v", base.Env)
	}
}

// --- Run: Parsed==nil fallback path -----------------------------------------

// Contract: on Result.Parsed == nil (schema parse failure), Run returns the
// default-seeded value plus the Result, NOT an error.
func TestRunParsedNilReturnsSeededDefaults(t *testing.T) {
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
			// Non-fatal error result with no parsed output.
			return &harness.Result{IsError: true, ErrorMessage: "schema validation failed after retries", Parsed: nil}, nil
		},
	}

	out, res, err := Run[seededResult](context.Background(), mh, "prompt", harness.Options{})
	if err != nil {
		t.Fatalf("expected no error on Parsed==nil, got %v", err)
	}
	if out == nil {
		t.Fatal("expected a seeded default value, got nil")
	}
	if !out.Complete || out.Scope != "medium" {
		t.Fatalf("expected seeded defaults (Complete=true, Scope=medium), got %+v", *out)
	}
	if res == nil || !res.IsError {
		t.Fatal("expected the failing Result returned so the caller can inspect IsError")
	}
}

func TestRunRecoversSchemaValidJSONWithBraceInsideString(t *testing.T) {
	const weakModelText = `I finished the task. {"complete":false,"estimated_scope":"medium ${foo {"} trailing prose`
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
			return &harness.Result{
				Result:       weakModelText,
				Parsed:       nil,
				IsError:      true,
				ErrorMessage: "schema validation failed after retries",
				FailureType:  harness.FailureSchema,
			}, nil
		},
	}

	out, res, err := Run[seededResult](context.Background(), mh, "prompt", harness.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil || out.Complete || out.Scope != "medium ${foo {" {
		t.Fatalf("expected recovered structured result, got %+v", out)
	}
	if res == nil || res.Parsed == nil || res.IsError || res.FailureType != harness.FailureNone {
		t.Fatalf("expected recovery to normalize harness result, got %+v", res)
	}
}

func TestRunRecoveryRejectsSchemaInvalidJSON(t *testing.T) {
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
			return &harness.Result{
				Result:       `prefix {"complete":false,"estimated_scope":42} suffix`,
				Parsed:       nil,
				IsError:      true,
				ErrorMessage: "schema validation failed after retries",
				FailureType:  harness.FailureSchema,
			}, nil
		},
	}

	out, res, err := Run[seededResult](context.Background(), mh, "prompt", harness.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil || !out.Complete || out.Scope != "medium" {
		t.Fatalf("schema-invalid text must fall back to seeded defaults, got %+v", out)
	}
	if res == nil || res.Parsed != nil || !res.IsError || res.FailureType != harness.FailureSchema {
		t.Fatalf("schema-invalid text must remain a harness failure, got %+v", res)
	}
}

func TestRunRecoveryDoesNotMaskProviderFailure(t *testing.T) {
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
			return &harness.Result{
				Result:       `{"complete":false,"estimated_scope":"medium"}`,
				Parsed:       nil,
				IsError:      true,
				ErrorMessage: "Stream error occurred",
				FailureType:  harness.FailureAPIError,
			}, nil
		},
	}

	out, res, err := Run[seededResult](context.Background(), mh, "prompt", harness.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil || !out.Complete || out.Scope != "medium" {
		t.Fatalf("provider failure must fall back to seeded defaults, got %+v", out)
	}
	if res == nil || res.Parsed != nil || !res.IsError || res.FailureType != harness.FailureAPIError {
		t.Fatalf("provider failure must remain visible, got %+v", res)
	}
}

func TestRunRecoversObservedWeakCoderJSONShape(t *testing.T) {
	// Reproduces the 2026-09-03 L2 failure exactly: the useful first OpenCode
	// attempt emitted near-schema JSON in a text event, then two schema retries
	// produced unrelated final text. Runner.Result therefore contains retry
	// garbage while Runner.Messages still preserves the earlier useful JSON.
	const observed = `{"files_changed":["calculator.py","cli.py","test_calculator.py"],"summary":"Added power and modulo operations.","complete":true,"tests_passed":true,"test_summary":"All tests passed.","codebase_learnings":["Test framework is unittest."],"agent_retro":"The implementation was straightforward."}`
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
			return &harness.Result{
				Result: "OpenCode is an open-source AI coding agent.",
				Messages: []map[string]any{
					{"type": "text", "part": map[string]any{"text": observed}},
					{"type": "text", "part": map[string]any{"text": "OpenCode is an open-source AI coding agent."}},
				},
				Parsed:       nil,
				IsError:      true,
				ErrorMessage: "schema validation failed after retries",
				FailureType:  harness.FailureSchema,
			}, nil
		},
	}

	out, res, err := Run[schemas.CoderResult](context.Background(), mh, "prompt", harness.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil || !out.Complete || out.Summary == "" || len(out.FilesChanged) != 3 {
		t.Fatalf("expected observed coder JSON to recover, got %+v", out)
	}
	if got, ok := out.AgentRetro["summary"].(string); !ok || got != "The implementation was straightforward." {
		t.Fatalf("expected string agent_retro to normalize under summary, got %#v", out.AgentRetro)
	}
	if out.IterationID != "" || out.RepoName != "" {
		t.Fatalf("omitted defaulted fields should materialize as empty defaults, got iteration=%q repo=%q", out.IterationID, out.RepoName)
	}
	if res == nil || res.Parsed == nil || res.IsError || res.FailureType != harness.FailureNone {
		t.Fatalf("expected recovered result to be marked successful, got %+v", res)
	}
}

func TestRunRecoversValidStructuredTextAfterNoProgressWatchdog(t *testing.T) {
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
			return &harness.Result{
				Messages: []map[string]any{{
					"type": "text",
					"part": map[string]any{"text": `{"complete":false,"estimated_scope":"medium"}`},
				}},
				Parsed:  nil,
				IsError: true,
			}, errors.New("CLI command made no progress for 300s: opencode run --format json")
		},
	}

	out, res, err := Run[seededResult](context.Background(), mh, "prompt", harness.Options{})
	if err != nil {
		t.Fatalf("watchdog should not discard schema-valid assistant output: %v", err)
	}
	if out == nil || out.Complete || out.Scope != "medium" {
		t.Fatalf("expected recovered watchdog result, got %+v", out)
	}
	if res == nil || res.Parsed == nil || res.IsError || res.FailureType != harness.FailureNone {
		t.Fatalf("expected recovered watchdog result normalized, got %+v", res)
	}
}

func TestRunRecoversValidOutputFileDeletedBeforeWatchdogReturns(t *testing.T) {
	projectDir := t.TempDir()
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, _ any, opts harness.Options) (*harness.Result, error) {
			outDir, err := os.MkdirTemp(opts.ProjectDir, ".agentfield-out-")
			if err != nil {
				return nil, err
			}
			outputPath := filepath.Join(outDir, ".agentfield_output.json")
			if err := os.WriteFile(outputPath, []byte(`{"complete":false,"estimated_scope":"large"}`), 0o600); err != nil {
				return nil, err
			}
			// Simulate the SDK lifecycle: the file is valid while the provider is
			// still running, then its temp directory is removed before Harness
			// returns the watchdog error to the SWE wrapper.
			time.Sleep(100 * time.Millisecond)
			_ = os.RemoveAll(outDir)
			return nil, errors.New("CLI command made no progress for 300s: opencode run --format json")
		},
	}

	out, res, err := Run[seededResult](context.Background(), mh, "prompt", harness.Options{ProjectDir: projectDir})
	if err != nil {
		t.Fatalf("watchdog should salvage exact-schema-valid output file: %v", err)
	}
	if out == nil || out.Complete || out.Scope != "large" {
		t.Fatalf("expected output-file recovery, got %+v", out)
	}
	if res == nil || res.Parsed == nil || res.IsError || res.FailureType != harness.FailureNone {
		t.Fatalf("expected recovered watchdog file result normalized, got %+v", res)
	}
}

func TestRunDoesNotRecoverGenericTransportErrorFromText(t *testing.T) {
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
			return &harness.Result{
				Messages: []map[string]any{{"type": "text", "part": map[string]any{"text": `{"complete":false,"estimated_scope":"medium"}`}}},
			}, errors.New("network blip")
		},
	}

	_, _, err := Run[seededResult](context.Background(), mh, "prompt", harness.Options{})
	if err == nil || !strings.Contains(err.Error(), "network blip") {
		t.Fatalf("generic transport errors must remain fail-closed, got %v", err)
	}
}

// --- Run: success path ------------------------------------------------------

func TestRunSuccessReturnsParsed(t *testing.T) {
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, dest any, _ harness.Options) (*harness.Result, error) {
			d := dest.(*seededResult)
			d.Complete = false
			d.Scope = "large"
			return &harness.Result{Parsed: dest}, nil
		},
	}

	out, _, err := Run[seededResult](context.Background(), mh, "prompt", harness.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil || out.Complete != false || out.Scope != "large" {
		t.Fatalf("expected parsed value returned, got %+v", out)
	}
}

func TestRunPropagatesTransportError(t *testing.T) {
	sentinel := errors.New("boom")
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
			return nil, sentinel
		},
	}
	_, _, err := Run[seededResult](context.Background(), mh, "prompt", harness.Options{})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected transport error propagated, got %v", err)
	}
}

// --- RoleOptions mapping -----------------------------------------------------

func TestRoleOptionsOpenCodeUsesOnlySWEOwnedBinary(t *testing.T) {
	t.Setenv("SWE_OPENCODE_BIN", "/opt/swe/opencode")
	t.Setenv("SEC_AF_OPENCODE_BIN", "/opt/sec/opencode")

	o := (RoleOptions{Provider: "opencode"}).ToOptions()
	if o.BinPath != "/opt/swe/opencode" {
		t.Fatalf("expected SWE-owned binary, got %q", o.BinPath)
	}
	if o.SchemaMode != "" {
		t.Fatalf("role options must not own structured-output recovery policy, got schema_mode=%q", o.SchemaMode)
	}

	t.Setenv("SWE_OPENCODE_BIN", "")
	o = (RoleOptions{Provider: "opencode"}).ToOptions()
	if o.BinPath != "" {
		t.Fatalf("SEC-AF binary must not leak into SWE options, got %q", o.BinPath)
	}
}

func TestRunCentralizesOpenCodeIncrementalSchemaMode(t *testing.T) {
	seenMode := ""
	mh := &mockHarness{
		fn: func(_ context.Context, _ string, _ map[string]any, dest any, opts harness.Options) (*harness.Result, error) {
			seenMode = opts.SchemaMode
			parsed := dest.(*seededResult)
			parsed.Complete = false
			parsed.Scope = "large"
			return &harness.Result{Parsed: parsed}, nil
		},
	}

	out, _, err := Run[seededResult](context.Background(), mh, "prompt", (RoleOptions{Provider: "opencode"}).ToOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if seenMode != "incremental" {
		t.Fatalf("central structured controller must enable incremental mode, got %q", seenMode)
	}
	if out == nil || out.Scope != "large" {
		t.Fatalf("unexpected parsed result: %+v", out)
	}
}

func TestRoleOptionsToOptions(t *testing.T) {
	ro := RoleOptions{
		Provider:       "claude-code",
		Model:          "sonnet",
		MaxTurns:       12,
		Tools:          []string{"Read", "Write", "Bash"},
		PermissionMode: "auto",
		SystemPrompt:   "sys",
		Cwd:            "/repo",
		Env:            map[string]string{"A": "1"},
	}
	o := ro.ToOptions()
	if o.Provider != "claude-code" || o.Model != "sonnet" || o.MaxTurns != 12 ||
		o.PermissionMode != "auto" || o.SystemPrompt != "sys" || o.Cwd != "/repo" {
		t.Fatalf("scalar option fields not mapped: %+v", o)
	}
	if len(o.Tools) != 3 || o.Tools[2] != "Bash" {
		t.Fatalf("tools not mapped: %v", o.Tools)
	}
	if o.Env["A"] != "1" {
		t.Fatalf("env not mapped: %v", o.Env)
	}
}

// --- helpers ----------------------------------------------------------------

func containsStr(xs []any, want string) bool {
	for _, x := range xs {
		if s, ok := x.(string); ok && s == want {
			return true
		}
	}
	return false
}
