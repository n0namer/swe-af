package node

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/sdk/go/agent"
	"github.com/Agent-Field/agentfield/sdk/go/harness"

	"github.com/Agent-Field/SWE-AF/go/internal/fast"
	"github.com/Agent-Field/SWE-AF/go/internal/furrow"
	"github.com/Agent-Field/SWE-AF/go/internal/workspace"
)

func TestFurrowRootResolutionPrecedence(t *testing.T) {
	t.Setenv("AGENTFIELD_HOME", "")
	t.Setenv("SWE_FURROW_DATA_DIR", "")
	t.Setenv("SWE_FURROW_REMOTES_ROOT", "")
	store, remotes := furrowRoots()
	if store != filepath.Join(workspace.Root(), ".furrow-store") || remotes != filepath.Join(workspace.Root(), ".furrow-remotes") {
		t.Fatalf("legacy roots = (%q, %q)", store, remotes)
	}

	home := t.TempDir()
	t.Setenv("AGENTFIELD_HOME", home)
	store, remotes = furrowRoots()
	if store != filepath.Join(home, "furrow", "store") || remotes != filepath.Join(home, "furrow", "remotes") {
		t.Fatalf("AGENTFIELD_HOME roots = (%q, %q)", store, remotes)
	}

	t.Setenv("SWE_FURROW_DATA_DIR", "/override/store")
	t.Setenv("SWE_FURROW_REMOTES_ROOT", "/override/remotes")
	store, remotes = furrowRoots()
	if store != "/override/store" || remotes != "/override/remotes" {
		t.Fatalf("override roots = (%q, %q)", store, remotes)
	}
}

// pythonRoleSurface is the independent parity checklist: the exact 25 role
// reasoner names the Python swe_af.reasoners.router registers (pipeline.py's 5
// planning roles + execution_agents.py's 20 execution roles). It is written from
// the Python inventory — NOT derived from the Go Handlers() maps — so the test
// catches a drift in either direction (a missing or an extra Go registration).
var pythonRoleSurface = []string{
	// pipeline.py — planning roles
	"run_product_manager",
	"run_environment_scout",
	"run_architect",
	"run_tech_lead",
	"run_sprint_planner",
	// execution_agents.py — coding roles
	"run_coder",
	"run_qa",
	"run_code_reviewer",
	"run_qa_synthesizer",
	// execution_agents.py — git/workspace roles
	"run_git_init",
	"run_workspace_setup",
	"run_workspace_cleanup",
	"run_merger",
	"run_integration_tester",
	"run_repo_finalize",
	"run_github_pr",
	// execution_agents.py — advisor/verify roles
	"run_retry_advisor",
	"run_issue_advisor",
	"run_replanner",
	"run_issue_writer",
	"run_verifier",
	"generate_fix_issues",
	// execution_agents.py — CI/resolve roles
	"run_ci_watcher",
	"run_ci_fixer",
	"run_pr_resolver",
}

// pythonOrchestrators is the 5 orchestrator reasoners defined on swe_af.app
// (app.py @app.reasoner()): build, plan, execute, resolve, resume_build.
// get_workspace_handle is deliberately NOT here: it is gated on furrow being
// switched on, and TestWorkspaceHandleReasonerIsGatedOnFurrow owns it.
var pythonOrchestrators = []string{"build", "plan", "execute", "resolve", "resume_build"}

// pythonFastReasoners is the 4 first-class fast reasoners: fast/app.py's build
// plus fast_plan_tasks / fast_execute_tasks / fast_verify.
var pythonFastReasoners = []string{"build", "fast_plan_tasks", "fast_execute_tasks", "fast_verify"}

// pythonIssueReasoners is the issue-level entry point (swe_af/issue/build.py),
// registered on BOTH nodes via the shared issue_router.
var pythonIssueReasoners = []string{"implement_issue"}

func TestRegisterPlannerExactSurface(t *testing.T) {
	// Pin the pro engine and furrow off so an inherited SWE_PRO_ENGINE,
	// SWE_FURROW_ENABLED or FURROW_PUBLIC_ADDR (which auto-enables mirroring
	// when the enable flag is unconfigured) cannot widen the surface under
	// test (each gated surface has its own test).
	t.Setenv("SWE_PRO_ENGINE", "")
	t.Setenv("SWE_BMAD_ENABLED", "")
	t.Setenv(furrow.EnvEnabled, "")
	t.Setenv(furrow.EnvPublicAddr, "")
	n, err := BuildAgent("swe-planner", "8005", "Autonomous SWE planning pipeline")
	if err != nil {
		t.Fatalf("BuildAgent: %v", err)
	}
	n.RegisterPlanner()

	// swe-planner surface = 25 roles + 5 orchestrators + implement_issue
	// = 31 unique names.
	want := append(append([]string(nil), pythonRoleSurface...), pythonOrchestrators...)
	want = append(want, pythonIssueReasoners...)
	assertSurface(t, "swe-planner", n.RegisteredNames(), want)
}

// stubAttacher is an enabled furrow that never mirrors anything: enough to
// open the registration gate, nothing more.
type stubAttacher struct{ enabled bool }

func (s stubAttacher) Enabled() bool                                         { return s.enabled }
func (s stubAttacher) Attach(string, string, string) (*furrow.Handle, error) { return nil, nil }
func (s stubAttacher) Publish(string, string) error                          { return nil }
func (s stubAttacher) Handle(string) *furrow.Handle                          { return nil }
func (s stubAttacher) Detach(string) error                                   { return nil }
func (s stubAttacher) Sweep(time.Duration, int64) (int, error)               { return 0, nil }

// get_workspace_handle hands out the connection details for a live workspace
// mirror. Mirroring is opt-in, so on a node that never makes a mirror the
// reasoner must not be advertised at all — an entrypoint-tagged surface that
// can only ever answer {"available": false} is an invitation to route to it.
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

func TestWorkspaceHandleReasonerIsGatedOnFurrow(t *testing.T) {
	const name = "get_workspace_handle"
	for _, tc := range []struct {
		label    string
		attacher furrow.Attacher
		want     bool
	}{
		{label: "furrow absent", attacher: nil, want: false},
		{label: "furrow present but disabled", attacher: stubAttacher{}, want: false},
		{label: "furrow mirroring", attacher: stubAttacher{enabled: true}, want: true},
	} {
		t.Run(tc.label, func(t *testing.T) {
			t.Setenv("SWE_PRO_ENGINE", "")
			t.Setenv(furrow.EnvEnabled, "")
			t.Setenv(furrow.EnvPublicAddr, "")
			n, err := BuildAgent("swe-planner", "8005", "Autonomous SWE planning pipeline")
			if err != nil {
				t.Fatalf("BuildAgent: %v", err)
			}
			n.Furrow = tc.attacher
			n.RegisterPlanner()
			if got := toSet(n.RegisteredNames())[name]; got != tc.want {
				t.Fatalf("%s registered = %v, want %v", name, got, tc.want)
			}
		})
	}
}

// get_workspace_handle answers anyone who can reach the node and name a run —
// there is no per-caller authorization anywhere on that path. The handle's Key
// decrypts the workspace and its Token authenticates to furrowd read-write, so
// neither may be the default answer to an unauthenticated question.
func TestWorkspaceHandleRedactsSecretsUnlessOperatorOptsIn(t *testing.T) {
	handle := &furrow.Handle{
		Version: furrow.HandleVersion, Remote: "ssh://node.internal:8802", Namespace: "run-1",
		Key: "0123456789abcdef", Token: "transport-token", RepoPath: "/work/repo",
	}
	for _, tc := range []struct {
		env  string
		want bool // secrets present
	}{
		{env: "", want: false},
		{env: "0", want: false},
		{env: "no", want: false},
		{env: "1", want: true},
		{env: "true", want: true},
	} {
		t.Run("SWE_FURROW_EXPOSE_SECRETS="+tc.env, func(t *testing.T) {
			t.Setenv(furrow.EnvExposeSecrets, tc.env)
			result := workspaceHandleResult(handle)

			_, gotKey := result["key"]
			_, gotToken := result["token"]
			if gotKey != tc.want || gotToken != tc.want {
				t.Fatalf("key present = %v, token present = %v, want both %v", gotKey, gotToken, tc.want)
			}
			if result["secrets_redacted"] != !tc.want {
				t.Errorf("secrets_redacted = %v, want %v", result["secrets_redacted"], !tc.want)
			}
			// Redacting must not blind the caller: what a mirror IS stays.
			for _, key := range []string{"v", "remote", "namespace", "repo_path"} {
				if _, ok := result[key]; !ok {
					t.Errorf("result dropped %q, which carries no secret", key)
				}
			}
			// Nothing may smuggle the secrets back through another field.
			rendered, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			if leaked := strings.Contains(string(rendered), handle.Key) ||
				strings.Contains(string(rendered), handle.Token); leaked != tc.want {
				t.Fatalf("secret material in payload = %v, want %v: %s", leaked, tc.want, rendered)
			}
		})
	}
}

func TestRegisterFastExactSurface(t *testing.T) {
	n, err := BuildAgent("swe-fast", "8006", "fast desc")
	if err != nil {
		t.Fatalf("BuildAgent: %v", err)
	}
	n.RegisterFast()

	// swe-fast surface = 25 roles + 4 fast reasoners + implement_issue
	// = 30 unique names. It must NOT contain plan/execute/resolve/resume_build
	// (those live only on swe-planner) — assertSurface's extra-name check
	// enforces that.
	want := append(append([]string(nil), pythonRoleSurface...), pythonFastReasoners...)
	want = append(want, pythonIssueReasoners...)
	assertSurface(t, "swe-fast", n.RegisteredNames(), want)
}

// TestFastWrappersAreBackedByRoles verifies the seven delegating wrappers
// (fast/__init__.py) are present on the swe-fast surface — each is one of the
// role names, backed by the full-pipeline role handler (fast.Wrappers identity).
func TestFastWrappersAreBackedByRoles(t *testing.T) {
	n, err := BuildAgent("swe-fast", "8006", "fast desc")
	if err != nil {
		t.Fatalf("BuildAgent: %v", err)
	}
	n.RegisterFast()

	got := toSet(n.RegisteredNames())
	for _, w := range fast.WrapperNames() {
		if !got[w] {
			t.Errorf("fast wrapper %q not registered on swe-fast surface", w)
		}
	}
	// Every wrapper must also be one of the role names (identity delegation).
	roleSet := toSet(pythonRoleSurface)
	for _, w := range fast.WrapperNames() {
		if !roleSet[w] {
			t.Errorf("fast wrapper %q is not a role name — delegation is not identity", w)
		}
	}
}

// TestRegHandlerRoutesToPackageHandler validates the adapter closure: it must
// invoke the package handler with the captured Deps and the request input, and
// propagate the handler's return value. This exercises the single routing seam
// every reasoner registration flows through, without a control plane.
func TestRegHandlerRoutesToPackageHandler(t *testing.T) {
	app, err := agent.New(agent.Config{NodeID: "test", Version: "1.0.0", ListenAddress: ":0"})
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}
	n := &Node{App: app}

	type fakeDeps struct{ marker string }
	deps := &fakeDeps{marker: "wired"}

	var gotDeps *fakeDeps
	var gotInput map[string]any
	h := func(_ context.Context, d *fakeDeps, in map[string]any) (any, error) {
		gotDeps = d
		gotInput = in
		return map[string]any{"echo": in["x"]}, nil
	}

	regHandler(n, "probe", deps, h)

	if len(n.registered) != 1 || n.registered[0] != "probe" {
		t.Fatalf("regHandler did not record name: %v", n.registered)
	}

	out, err := app.Execute(context.Background(), "probe", map[string]any{"x": 42})
	if err != nil {
		t.Fatalf("Execute(probe): %v", err)
	}
	if gotDeps != deps {
		t.Errorf("handler received deps %v, want the captured %v", gotDeps, deps)
	}
	if gotInput["x"] != 42 {
		t.Errorf("handler received input %v, want x=42", gotInput)
	}
	m, ok := out.(map[string]any)
	if !ok || m["echo"] != 42 {
		t.Errorf("handler return not propagated: got %#v", out)
	}
}

func TestBMADMethodPreservesStepOrderAndState(t *testing.T) {
	method := bmadMethod{ID: "test-method", Source: "deadbeef", Steps: []bmadStep{{ID: "one", Text: "one"}, {ID: "two", Text: "two"}, {ID: "three", Text: "three"}}}
	seen := []string{}
	exec := func(_ context.Context, _ bmadMethod, step bmadStep, state, artifacts map[string]any) (*bmadStepResult, error) {
		seen = append(seen, step.ID)
		switch step.ID {
		case "one":
			return &bmadStepResult{Status: "completed", Summary: "one", State: map[string]any{"from_one": true}, Artifacts: map[string]any{"artifact_one": "x"}}, nil
		case "two":
			if state["from_one"] != true || artifacts["artifact_one"] != "x" { t.Fatalf("step two did not receive prior state/artifacts: state=%v artifacts=%v", state, artifacts) }
			return &bmadStepResult{Status: "completed", Summary: "two", State: map[string]any{"from_two": true}}, nil
		default:
			if state["from_two"] != true { t.Fatalf("step three did not receive step two state: %v", state) }
			return &bmadStepResult{Status: "completed", Summary: "three", Output: "final"}, nil
		}
	}
	got, err := runBMADMethod(context.Background(), method, map[string]any{"initial": 1}, exec)
	if err != nil { t.Fatalf("runBMADMethod: %v", err) }
	if strings.Join(seen, ",") != "one,two,three" { t.Fatalf("step order = %v", seen) }
	if got.Status != "completed" || got.Output != "final" || got.SourceCommit != "deadbeef" { t.Fatalf("unexpected result: %#v", got) }
	if strings.Join(got.CompletedSteps, ",") != "one,two,three" { t.Fatalf("completed steps = %v", got.CompletedSteps) }
}

func TestBMADMethodStopsAtBlockedStep(t *testing.T) {
	method := bmadMethod{ID: "test-method", Source: "deadbeef", Steps: []bmadStep{{ID: "one", Text: "one"}, {ID: "two", Text: "two"}, {ID: "must-not-run", Text: "three"}}}
	seen := []string{}
	exec := func(_ context.Context, _ bmadMethod, step bmadStep, _, _ map[string]any) (*bmadStepResult, error) { seen = append(seen, step.ID); if step.ID == "two" { return &bmadStepResult{Status: "blocked", Summary: "needs evidence"}, nil }; return &bmadStepResult{Status: "completed", Summary: "ok"}, nil }
	got, err := runBMADMethod(context.Background(), method, nil, exec)
	if err != nil { t.Fatalf("runBMADMethod: %v", err) }
	if got.Status != "blocked" || strings.Join(seen, ",") != "one,two" { t.Fatalf("blocked workflow continued: result=%#v seen=%v", got, seen) }
	if strings.Join(got.CompletedSteps, ",") != "one" { t.Fatalf("blocked step counted as completed: %v", got.CompletedSteps) }
}

func TestBMADAdversarialMethodPinnedAndOrdered(t *testing.T) {
	if adversarialGeneralMethod.Source != bmadSourceCommit || len(adversarialGeneralMethod.Steps) != 3 { t.Fatalf("unexpected pinned method: %#v", adversarialGeneralMethod) }
	want := "receive-content,adversarial-analysis,present-findings"
	ids := make([]string, 0, len(adversarialGeneralMethod.Steps)); for _, step := range adversarialGeneralMethod.Steps { ids = append(ids, step.ID) }
	if strings.Join(ids, ",") != want { t.Fatalf("step order = %v, want %s", ids, want) }
}

func TestBMADEdgeCaseMethodPinnedAndOrdered(t *testing.T) {
	if edgeCaseHunterMethod.Source != bmadSourceCommit || len(edgeCaseHunterMethod.Steps) != 5 { t.Fatalf("unexpected pinned edge method: %#v", edgeCaseHunterMethod) }
	want := "receive-content,exhaustive-path-analysis,validate-completeness,deletion-check,present-findings"
	ids := make([]string, 0, len(edgeCaseHunterMethod.Steps)); for _, step := range edgeCaseHunterMethod.Steps { ids = append(ids, step.ID) }
	if strings.Join(ids, ",") != want { t.Fatalf("edge step order = %v, want %s", ids, want) }
}

type bmadHarnessStub struct { fn func(schema map[string]any, dest any, opts harness.Options) (*harness.Result, error) }
func (b *bmadHarnessStub) Harness(_ context.Context, _ string, schema map[string]any, dest any, opts harness.Options) (*harness.Result, error) { return b.fn(schema, dest, opts) }

func TestBMADTextStepFailsClosedOnInvalidTextEnvelope(t *testing.T) {
	t.Setenv("SWE_DEFAULT_RUNTIME", "claude_code"); t.Setenv("SWE_DEFAULT_MODEL", "sonnet")
	h := &bmadHarnessStub{fn: func(schema map[string]any, dest any, _ harness.Options) (*harness.Result, error) { if schema != nil || dest != nil { t.Fatalf("BMAD text path unexpectedly requested schema output") }; return &harness.Result{Result: "not-json"}, nil }}
	_, err := runBMADTextStep(context.Background(), h, adversarialGeneralMethod, adversarialGeneralMethod.Steps[0], map[string]any{"content": "diff"}, nil)
	if err == nil || !strings.Contains(err.Error(), "decode BMAD step JSON") { t.Fatalf("error=%v, want strict JSON envelope failure", err) }
}

func TestBMADTextStepPropagatesProviderFailure(t *testing.T) {
	t.Setenv("SWE_DEFAULT_RUNTIME", "claude_code"); t.Setenv("SWE_DEFAULT_MODEL", "sonnet")
	h := &bmadHarnessStub{fn: func(_ map[string]any, _ any, _ harness.Options) (*harness.Result, error) { return &harness.Result{IsError: true, ErrorMessage: "provider failed"}, nil }}
	_, err := runBMADTextStep(context.Background(), h, adversarialGeneralMethod, adversarialGeneralMethod.Steps[0], map[string]any{"content": "diff"}, nil)
	if err == nil || !strings.Contains(err.Error(), "provider failed") { t.Fatalf("error=%v", err) }
}

func TestBMADTextStepUsesReadOnlyTextPolicy(t *testing.T) {
	t.Setenv("SWE_DEFAULT_RUNTIME", "open_code")
	t.Setenv("SWE_DEFAULT_MODEL", "test-model")
	var got harness.Options
	h := &bmadHarnessStub{fn: func(schema map[string]any, dest any, opts harness.Options) (*harness.Result, error) {
		if schema != nil || dest != nil { t.Fatalf("schema=%v dest=%T, want nil/nil", schema, dest) }
		got = opts
		return &harness.Result{Result: `{"status":"completed","summary":"loaded"}`}, nil
	}}
	if _, err := runBMADTextStep(context.Background(), h, adversarialGeneralMethod, adversarialGeneralMethod.Steps[0], map[string]any{"content": "diff"}, nil); err != nil { t.Fatalf("runBMADTextStep: %v", err) }
	var policy map[string]any
	if err := json.Unmarshal([]byte(got.Env["OPENCODE_CONFIG_CONTENT"]), &policy); err != nil { t.Fatalf("invalid BMAD OpenCode policy: %v", err) }
	permission := policy["permission"].(map[string]any)
	if permission["bash"] != "deny" { t.Fatalf("bash permission=%v", permission["bash"]) }
	edit := permission["edit"].(map[string]any)
	if edit["*"] != "deny" { t.Fatalf("edit policy=%v", edit) }
	if strings.Join(got.Tools, ",") != "Read" { t.Fatalf("tools=%v, want Read only", got.Tools) }
	if got.PermissionMode != "plan" { t.Fatalf("permission_mode=%q, want plan", got.PermissionMode) }
	if got.Timeout != 300 { t.Fatalf("timeout=%d, want 300", got.Timeout) }
	if got.Cwd == "" || !strings.Contains(filepath.Base(got.Cwd), "swe-bmad-review-") { t.Fatalf("cwd=%q", got.Cwd) }
	if _, err := os.Stat(got.Cwd); !os.IsNotExist(err) { t.Fatalf("BMAD temp cwd survived cleanup: stat err=%v", err) }
}

func TestDecodeBMADStepResultIsStrict(t *testing.T) {
	for _, tc := range []struct { name, text string; wantErr bool }{
		{name: "valid", text: `{"status":"completed","summary":"ok"}`},
		{name: "unknown field", text: `{"status":"completed","summary":"ok","extra":1}`, wantErr: true},
		{name: "trailing value", text: `{"status":"completed","summary":"ok"} {}`, wantErr: true},
		{name: "empty", text: ``, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) { _, err := decodeBMADStepResult(tc.text); if (err != nil) != tc.wantErr { t.Fatalf("err=%v wantErr=%v", err, tc.wantErr) } })
	}
}

func TestBMADEdgeEmptyInputReturnsCanonicalFinding(t *testing.T) {
	t.Setenv("SWE_PRO_ENGINE", ""); t.Setenv("SWE_BMAD_ENABLED", "1"); t.Setenv(furrow.EnvEnabled, ""); t.Setenv(furrow.EnvPublicAddr, "")
	n, err := BuildAgent("swe-planner", "8005", "Autonomous SWE planning pipeline"); if err != nil { t.Fatalf("BuildAgent: %v", err) }
	n.RegisterPlanner()
	out, err := n.App.Execute(context.Background(), bmadReviewEdgeCaseHunter, map[string]any{"content": ""}); if err != nil { t.Fatalf("edge empty input: %v", err) }
	got, ok := out.(*bmadRunResult); if !ok { t.Fatalf("result type=%T", out) }
	if got.Status != "completed" || strings.Join(got.CompletedSteps, ",") != "receive-content" { t.Fatalf("unexpected edge empty result: %#v", got) }
	if err := validateBMADRunOutput(edgeCaseHunterMethod, got); err != nil { t.Fatalf("canonical edge empty output invalid: %v", err) }
}

func TestBMADInputSchemaIsValidAndBounded(t *testing.T) {
	for _, allowEmpty := range []bool{false, true} {
		raw := bmadInputSchema(allowEmpty)
		if !json.Valid(raw) { t.Fatalf("allowEmpty=%v invalid schema: %s", allowEmpty, raw) }
		var doc map[string]any; if err := json.Unmarshal(raw, &doc); err != nil { t.Fatal(err) }
		props := doc["properties"].(map[string]any); content := props["content"].(map[string]any)
		if int(content["maxLength"].(float64)) != bmadMaxContentChars { t.Fatalf("content bound=%v", content) }
		if allowEmpty { if _, ok := content["minLength"]; ok { t.Fatalf("edge content unexpectedly requires nonempty input: %v", content) } } else if int(content["minLength"].(float64)) != 1 { t.Fatalf("adversarial minLength=%v", content) }
		also := props["also_consider"].(map[string]any); if int(also["maxLength"].(float64)) != bmadMaxAlsoConsiderChars { t.Fatalf("also bound=%v", also) }
	}
}

func TestBMADMethodRejectsDuplicateStepIDs(t *testing.T) {
	method := bmadMethod{ID: "dup", Source: "deadbeef", Steps: []bmadStep{{ID: "same", Text: "one"}, {ID: "same", Text: "two"}}}
	_, err := runBMADMethod(context.Background(), method, nil, func(_ context.Context, _ bmadMethod, _ bmadStep, _, _ map[string]any) (*bmadStepResult, error) { return &bmadStepResult{Status: "completed", Summary: "ok"}, nil })
	if err == nil || !strings.Contains(err.Error(), "duplicate step") { t.Fatalf("error=%v, want duplicate step rejection", err) }
}

func TestBMADFinalOutputValidationFailsClosed(t *testing.T) {
	for _, tc := range []struct { name string; method bmadMethod; output string; wantErr bool }{
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
		t.Run(tc.name, func(t *testing.T) { err := validateBMADRunOutput(tc.method, &bmadRunResult{Status: "completed", Output: tc.output}); if (err != nil) != tc.wantErr { t.Fatalf("error=%v, wantErr=%v", err, tc.wantErr) } })
	}
}

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

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// assertSurface fails if got (the registered names) does not equal want as a
// set, or if got contains duplicates. Reporting missing and extra names
// separately makes a parity drift immediately diagnosable.
func assertSurface(t *testing.T, node string, got, want []string) {
	t.Helper()

	// Duplicate guard: RegisterReasoner dedupes by name in its map, so a
	// duplicate in the recorded slice means two registrations collided on one
	// name (a silent surface bug the set comparison would otherwise hide).
	seen := map[string]int{}
	for _, name := range got {
		seen[name]++
	}
	for name, c := range seen {
		if c > 1 {
			t.Errorf("[%s] reasoner %q registered %d times (collision)", node, name, c)
		}
	}

	gotSet := toSet(got)
	wantSet := toSet(want)

	var missing, extra []string
	for name := range wantSet {
		if !gotSet[name] {
			missing = append(missing, name)
		}
	}
	for name := range gotSet {
		if !wantSet[name] {
			extra = append(extra, name)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)

	if len(missing) > 0 {
		t.Errorf("[%s] missing reasoners (in Python, not registered): %v", node, missing)
	}
	if len(extra) > 0 {
		t.Errorf("[%s] extra reasoners (registered, not in Python): %v", node, extra)
	}
	if len(gotSet) != len(wantSet) {
		t.Errorf("[%s] surface size = %d, want %d", node, len(gotSet), len(wantSet))
	}
}

func toSet(names []string) map[string]bool {
	s := make(map[string]bool, len(names))
	for _, n := range names {
		s[n] = true
	}
	return s
}
