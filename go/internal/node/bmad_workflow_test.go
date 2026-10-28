package node

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Agent-Field/agentfield/sdk/go/harness"

	"github.com/Agent-Field/SWE-AF/go/internal/furrow"
)

func TestBMADWorkflowReasonersAreOptIn(t *testing.T) {
	names := []string{bmadReviewAdversarialGeneral, bmadReviewEdgeCaseHunter}
	for _, tc := range []struct {
		enabled string
		want    bool
	}{{enabled: "", want: false}, {enabled: "1", want: true}} {
		t.Run("enabled="+tc.enabled, func(t *testing.T) {
			t.Setenv("SWE_PRO_ENGINE", "")
			t.Setenv("SWE_BMAD_ENABLED", tc.enabled)
			t.Setenv(furrow.EnvEnabled, "")
			t.Setenv(furrow.EnvPublicAddr, "")
			n, err := BuildAgent("swe-planner", "8005", "Autonomous SWE planning pipeline")
			if err != nil {
				t.Fatalf("BuildAgent: %v", err)
			}
			n.RegisterPlanner()
			for _, name := range names {
				if got := toSet(n.RegisteredNames())[name]; got != tc.want {
					t.Fatalf("BMAD workflow %s registered=%v, want %v", name, got, tc.want)
				}
				if tc.want {
					meta := n.RegisteredMeta()[name]
					if !toSet(meta.Tags)["bmad"] || !toSet(meta.Tags)[tagEntrypoint] {
						t.Fatalf("BMAD workflow %s tags=%v", name, meta.Tags)
					}
				}
			}
		})
	}
}

func TestBMADMethodPreservesStepOrderAndState(t *testing.T) {
	method := bmadMethod{ID: "test-method", Source: "deadbeef", Steps: []bmadStep{
		{ID: "one", Text: "one"}, {ID: "two", Text: "two"}, {ID: "three", Text: "three"},
	}}
	seen := []string{}
	exec := func(_ context.Context, _ bmadMethod, step bmadStep, state, artifacts map[string]any) (*bmadStepResult, error) {
		seen = append(seen, step.ID)
		switch step.ID {
		case "one":
			return &bmadStepResult{Status: "completed", Summary: "one", State: map[string]any{"from_one": true}, Artifacts: map[string]any{"artifact_one": "x"}}, nil
		case "two":
			if state["from_one"] != true || artifacts["artifact_one"] != "x" {
				t.Fatalf("step two did not receive prior state/artifacts: state=%v artifacts=%v", state, artifacts)
			}
			return &bmadStepResult{Status: "completed", Summary: "two", State: map[string]any{"from_two": true}}, nil
		default:
			if state["from_two"] != true {
				t.Fatalf("step three did not receive step two state: %v", state)
			}
			return &bmadStepResult{Status: "completed", Summary: "three", Output: "final"}, nil
		}
	}

	got, err := runBMADMethod(context.Background(), method, map[string]any{"initial": 1}, exec)
	if err != nil {
		t.Fatalf("runBMADMethod: %v", err)
	}
	if strings.Join(seen, ",") != "one,two,three" {
		t.Fatalf("step order = %v", seen)
	}
	if got.Status != "completed" || got.Output != "final" || got.SourceCommit != "deadbeef" {
		t.Fatalf("unexpected result: %#v", got)
	}
	if strings.Join(got.CompletedSteps, ",") != "one,two,three" {
		t.Fatalf("completed steps = %v", got.CompletedSteps)
	}
}

func TestBMADMethodStopsAtBlockedStep(t *testing.T) {
	method := bmadMethod{ID: "test-method", Source: "deadbeef", Steps: []bmadStep{
		{ID: "one", Text: "one"}, {ID: "two", Text: "two"}, {ID: "must-not-run", Text: "three"},
	}}
	seen := []string{}
	exec := func(_ context.Context, _ bmadMethod, step bmadStep, _, _ map[string]any) (*bmadStepResult, error) {
		seen = append(seen, step.ID)
		if step.ID == "two" {
			return &bmadStepResult{Status: "blocked", Summary: "needs evidence"}, nil
		}
		return &bmadStepResult{Status: "completed", Summary: "ok"}, nil
	}

	got, err := runBMADMethod(context.Background(), method, nil, exec)
	if err != nil {
		t.Fatalf("runBMADMethod: %v", err)
	}
	if got.Status != "blocked" || strings.Join(seen, ",") != "one,two" {
		t.Fatalf("blocked workflow continued: result=%#v seen=%v", got, seen)
	}
	if strings.Join(got.CompletedSteps, ",") != "one" {
		t.Fatalf("blocked step counted as completed: %v", got.CompletedSteps)
	}
}

func TestBMADAdversarialMethodPinnedAndOrdered(t *testing.T) {
	if adversarialGeneralMethod.Source != bmadSourceCommit || len(adversarialGeneralMethod.Steps) != 3 {
		t.Fatalf("unexpected pinned method: %#v", adversarialGeneralMethod)
	}
	want := "receive-content,adversarial-analysis,present-findings"
	ids := make([]string, 0, len(adversarialGeneralMethod.Steps))
	for _, step := range adversarialGeneralMethod.Steps {
		ids = append(ids, step.ID)
	}
	if strings.Join(ids, ",") != want {
		t.Fatalf("step order = %v, want %s", ids, want)
	}
}

func TestBMADEdgeCaseMethodPinnedAndOrdered(t *testing.T) {
	if edgeCaseHunterMethod.Source != bmadSourceCommit || len(edgeCaseHunterMethod.Steps) != 5 {
		t.Fatalf("unexpected pinned edge method: %#v", edgeCaseHunterMethod)
	}
	want := "receive-content,exhaustive-path-analysis,validate-completeness,deletion-check,present-findings"
	ids := make([]string, 0, len(edgeCaseHunterMethod.Steps))
	for _, step := range edgeCaseHunterMethod.Steps {
		ids = append(ids, step.ID)
	}
	if strings.Join(ids, ",") != want {
		t.Fatalf("edge step order = %v, want %s", ids, want)
	}
}

type bmadHarnessStub struct {
	fn func(schema map[string]any, dest any, opts harness.Options) (*harness.Result, error)
}

func (b *bmadHarnessStub) Harness(_ context.Context, _ string, schema map[string]any, dest any, opts harness.Options) (*harness.Result, error) {
	return b.fn(schema, dest, opts)
}

func TestBMADTextStepFailsClosedOnInvalidTextEnvelope(t *testing.T) {
	t.Setenv("SWE_DEFAULT_RUNTIME", "claude_code")
	t.Setenv("SWE_DEFAULT_MODEL", "sonnet")
	h := &bmadHarnessStub{fn: func(schema map[string]any, dest any, _ harness.Options) (*harness.Result, error) {
		if schema != nil || dest != nil {
			t.Fatalf("BMAD text path unexpectedly requested schema output")
		}
		return &harness.Result{Result: "not-json"}, nil
	}}
	_, err := runBMADTextStep(context.Background(), h, adversarialGeneralMethod, adversarialGeneralMethod.Steps[0], map[string]any{"content": "diff"}, nil)
	if err == nil || !strings.Contains(err.Error(), "decode BMAD step JSON") {
		t.Fatalf("error=%v, want strict JSON envelope failure", err)
	}
}

func TestBMADTextStepPropagatesProviderFailure(t *testing.T) {
	t.Setenv("SWE_DEFAULT_RUNTIME", "claude_code")
	t.Setenv("SWE_DEFAULT_MODEL", "sonnet")
	h := &bmadHarnessStub{fn: func(_ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
		return &harness.Result{IsError: true, ErrorMessage: "provider failed"}, nil
	}}
	_, err := runBMADTextStep(context.Background(), h, adversarialGeneralMethod, adversarialGeneralMethod.Steps[0], map[string]any{"content": "diff"}, nil)
	if err == nil || !strings.Contains(err.Error(), "provider failed") {
		t.Fatalf("error=%v", err)
	}
}

func TestBMADTextStepUsesReadOnlyTextPolicy(t *testing.T) {
	t.Setenv("SWE_DEFAULT_RUNTIME", "open_code")
	t.Setenv("SWE_DEFAULT_MODEL", "test-model")
	var got harness.Options
	h := &bmadHarnessStub{fn: func(schema map[string]any, dest any, opts harness.Options) (*harness.Result, error) {
		if schema != nil || dest != nil {
			t.Fatalf("schema=%v dest=%T, want nil/nil", schema, dest)
		}
		got = opts
		return &harness.Result{Result: `{"status":"completed","summary":"loaded"}`}, nil
	}}
	if _, err := runBMADTextStep(context.Background(), h, adversarialGeneralMethod, adversarialGeneralMethod.Steps[0], map[string]any{"content": "diff"}, nil); err != nil {
		t.Fatalf("runBMADTextStep: %v", err)
	}
	var policy map[string]any
	if err := json.Unmarshal([]byte(got.Env["OPENCODE_CONFIG_CONTENT"]), &policy); err != nil {
		t.Fatalf("invalid BMAD OpenCode policy: %v", err)
	}
	permission := policy["permission"].(map[string]any)
	if permission["*"] != "deny" {
		t.Fatalf("wildcard permission=%v, want deny", permission["*"])
	}
	for _, key := range []string{"read", "glob", "grep", "list"} {
		if permission[key] != "allow" {
			t.Fatalf("%s permission=%v, want allow", key, permission[key])
		}
	}
	if strings.Join(got.Tools, ",") != "Read" {
		t.Fatalf("tools=%v, want Read only", got.Tools)
	}
	if got.PermissionMode != "plan" {
		t.Fatalf("permission_mode=%q, want plan", got.PermissionMode)
	}
	if got.Timeout != 300 {
		t.Fatalf("timeout=%d, want 300", got.Timeout)
	}
	if got.Cwd == "" || !strings.Contains(filepath.Base(got.Cwd), "swe-bmad-review-") {
		t.Fatalf("cwd=%q", got.Cwd)
	}
	if _, err := os.Stat(got.Cwd); !os.IsNotExist(err) {
		t.Fatalf("BMAD temp cwd survived cleanup: stat err=%v", err)
	}
}

func TestBMADTextStepRejectsCodexWithoutReadOnlyToolAllowlist(t *testing.T) {
	t.Setenv("SWE_DEFAULT_RUNTIME", "codex")
	t.Setenv("SWE_DEFAULT_MODEL", "test-model")
	h := &bmadHarnessStub{fn: func(_ map[string]any, _ any, _ harness.Options) (*harness.Result, error) {
		t.Fatal("Codex harness must not be invoked when read-only tool allowlist cannot be enforced")
		return nil, nil
	}}
	_, err := runBMADTextStep(context.Background(), h, adversarialGeneralMethod, adversarialGeneralMethod.Steps[0], map[string]any{"content": "diff"}, nil)
	if err == nil || !strings.Contains(err.Error(), "does not support Codex") {
		t.Fatalf("error=%v, want fail-closed Codex rejection", err)
	}
}

func TestBMADEdgeEmptyInputReturnsCanonicalFinding(t *testing.T) {
	t.Setenv("SWE_PRO_ENGINE", "")
	t.Setenv("SWE_BMAD_ENABLED", "1")
	t.Setenv(furrow.EnvEnabled, "")
	t.Setenv(furrow.EnvPublicAddr, "")
	n, err := BuildAgent("swe-planner", "8005", "Autonomous SWE planning pipeline")
	if err != nil {
		t.Fatalf("BuildAgent: %v", err)
	}
	n.RegisterPlanner()
	out, err := n.App.Execute(context.Background(), bmadReviewEdgeCaseHunter, map[string]any{"content": ""})
	if err != nil {
		t.Fatalf("edge empty input: %v", err)
	}
	got, ok := out.(*bmadRunResult)
	if !ok {
		t.Fatalf("result type=%T", out)
	}
	if got.Status != "completed" || strings.Join(got.CompletedSteps, ",") != "receive-content" {
		t.Fatalf("unexpected edge empty result: %#v", got)
	}
	if err := validateBMADRunOutput(edgeCaseHunterMethod, got); err != nil {
		t.Fatalf("canonical edge empty output invalid: %v", err)
	}
}

func TestBMADInputSchemaIsValidAndBounded(t *testing.T) {
	for _, allowEmpty := range []bool{false, true} {
		raw := bmadInputSchema(allowEmpty)
		if !json.Valid(raw) {
			t.Fatalf("allowEmpty=%v invalid schema: %s", allowEmpty, raw)
		}
		var doc map[string]any
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatal(err)
		}
		props := doc["properties"].(map[string]any)
		content := props["content"].(map[string]any)
		if int(content["maxLength"].(float64)) != bmadMaxContentChars {
			t.Fatalf("content bound=%v", content)
		}
		if allowEmpty {
			if _, ok := content["minLength"]; ok {
				t.Fatalf("edge content unexpectedly requires nonempty input: %v", content)
			}
		} else if int(content["minLength"].(float64)) != 1 {
			t.Fatalf("adversarial minLength=%v", content)
		}
		also := props["also_consider"].(map[string]any)
		if int(also["maxLength"].(float64)) != bmadMaxAlsoConsiderChars {
			t.Fatalf("also bound=%v", also)
		}
	}
}

func TestBMADMethodRejectsDuplicateStepIDs(t *testing.T) {
	method := bmadMethod{ID: "dup", Source: "deadbeef", Steps: []bmadStep{{ID: "same", Text: "one"}, {ID: "same", Text: "two"}}}
	_, err := runBMADMethod(context.Background(), method, nil, func(_ context.Context, _ bmadMethod, _ bmadStep, _, _ map[string]any) (*bmadStepResult, error) {
		return &bmadStepResult{Status: "completed", Summary: "ok"}, nil
	})
	if err == nil || !strings.Contains(err.Error(), "duplicate step") {
		t.Fatalf("error=%v, want duplicate step rejection", err)
	}
}

func TestBMADFinalOutputValidationFailsClosed(t *testing.T) {
	for _, tc := range []struct {
		name    string
		method  bmadMethod
		output  string
		wantErr bool
	}{
		{name: "edge valid empty array", method: edgeCaseHunterMethod, output: `[]`, wantErr: false},
		{name: "edge prose", method: edgeCaseHunterMethod, output: `looks fine`, wantErr: true},
		{name: "edge missing field", method: edgeCaseHunterMethod, output: `[{"location":"x"}]`, wantErr: true},
		{name: "edge multiline guard", method: edgeCaseHunterMethod, output: "[{\"location\":\"x\",\"trigger_condition\":\"trigger\",\"guard_snippet\":\"a\\nb\",\"potential_consequence\":\"breakage\"}]", wantErr: true},
		{name: "edge long trigger", method: edgeCaseHunterMethod, output: `[{"location":"x","trigger_condition":"one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen","guard_snippet":"guard","potential_consequence":"breakage"}]`, wantErr: true},
		{name: "unknown method", method: bmadMethod{ID: "unknown"}, output: `anything`, wantErr: true},
		{name: "adversarial ten", method: adversarialGeneralMethod, output: "- f1\n- f2\n- f3\n- f4\n- f5\n- f6\n- f7\n- f8\n- f9\n- f10", wantErr: false},
		{name: "adversarial nine", method: adversarialGeneralMethod, output: "- f1\n- f2\n- f3\n- f4\n- f5\n- f6\n- f7\n- f8\n- f9", wantErr: true},
		{name: "adversarial prose", method: adversarialGeneralMethod, output: `no problems`, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := validateBMADRunOutput(tc.method, &bmadRunResult{Status: "completed", Output: tc.output})
			if (err != nil) != tc.wantErr {
				t.Fatalf("error=%v, wantErr=%v", err, tc.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// assertSurface fails if got (the registered names) does not equal want as a
// set, or if got contains duplicates. Reporting missing and extra names
// separately makes a parity drift immediately diagnosable.

func TestBMADMethodRejectsImmutableInputMutation(t *testing.T) {
	method := bmadMethod{ID: "immut", Source: "deadbeef", Steps: []bmadStep{{ID: "one", Text: "one"}}}
	_, err := runBMADMethod(context.Background(), method, map[string]any{"content": "original", "also_consider": "context"}, func(_ context.Context, _ bmadMethod, _ bmadStep, _, _ map[string]any) (*bmadStepResult, error) {
		return &bmadStepResult{Status: "completed", Summary: "mutate", State: map[string]any{"content": "rewritten"}}, nil
	})
	if err == nil || !strings.Contains(err.Error(), "attempted to mutate immutable input content") {
		t.Fatalf("error=%v, want immutable input rejection", err)
	}
}

func TestBMADMethodFinalOutputCannotReuseIntermediateOutput(t *testing.T) {
	method := bmadMethod{ID: "fresh-output", Source: "deadbeef", Steps: []bmadStep{{ID: "one", Text: "one"}, {ID: "two", Text: "two"}}}
	got, err := runBMADMethod(context.Background(), method, nil, func(_ context.Context, _ bmadMethod, step bmadStep, _, _ map[string]any) (*bmadStepResult, error) {
		if step.ID == "one" {
			return &bmadStepResult{Status: "completed", Summary: "one", Output: "stale"}, nil
		}
		return &bmadStepResult{Status: "completed", Summary: "two"}, nil
	})
	if err != nil {
		t.Fatalf("runBMADMethod: %v", err)
	}
	if got.Output != "" {
		t.Fatalf("final output=%q, want empty rather than stale intermediate output", got.Output)
	}
}

func TestBMADMethodStopsOnCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	method := bmadMethod{ID: "cancel", Source: "deadbeef", Steps: []bmadStep{{ID: "one", Text: "one"}, {ID: "two", Text: "two"}}}
	calls := 0
	_, err := runBMADMethod(ctx, method, nil, func(_ context.Context, _ bmadMethod, _ bmadStep, _, _ map[string]any) (*bmadStepResult, error) {
		calls++
		cancel()
		return &bmadStepResult{Status: "completed", Summary: "done"}, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v, want context.Canceled", err)
	}
	if calls != 1 {
		t.Fatalf("calls=%d, want exactly one step before cancellation", calls)
	}
}

func TestBMADReasonerRejectsWrongInputTypes(t *testing.T) {
	t.Setenv("SWE_PRO_ENGINE", "")
	t.Setenv("SWE_BMAD_ENABLED", "1")
	t.Setenv(furrow.EnvEnabled, "")
	t.Setenv(furrow.EnvPublicAddr, "")
	n, err := BuildAgent("swe-planner", "8005", "Autonomous SWE planning pipeline")
	if err != nil {
		t.Fatalf("BuildAgent: %v", err)
	}
	n.RegisterPlanner()
	for _, tc := range []map[string]any{{"content": 42}, {"content": "diff", "also_consider": 42}} {
		if _, err := n.App.Execute(context.Background(), bmadReviewEdgeCaseHunter, tc); err == nil {
			t.Fatalf("input=%v unexpectedly accepted", tc)
		}
	}
}

func TestBMADMethodCancellationAfterFinalStepStillFails(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	method := bmadMethod{ID: "cancel-final", Source: "deadbeef", Steps: []bmadStep{{ID: "one", Text: "one"}}}
	_, err := runBMADMethod(ctx, method, nil, func(_ context.Context, _ bmadMethod, _ bmadStep, _, _ map[string]any) (*bmadStepResult, error) {
		cancel()
		return &bmadStepResult{Status: "completed", Summary: "done"}, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v, want context.Canceled", err)
	}
}

func TestBMADMethodRejectsOversizedStepEnvelope(t *testing.T) {
	method := bmadMethod{ID: "oversized", Source: "deadbeef", Steps: []bmadStep{{ID: "one", Text: "one"}}}
	_, err := runBMADMethod(context.Background(), method, nil, func(_ context.Context, _ bmadMethod, _ bmadStep, _, _ map[string]any) (*bmadStepResult, error) {
		return &bmadStepResult{Status: "completed", Summary: "done", Output: strings.Repeat("x", bmadMaxStepEnvelopeBytes)}, nil
	})
	if err == nil || !strings.Contains(err.Error(), "envelope exceeds") {
		t.Fatalf("error=%v, want envelope size rejection", err)
	}
}
