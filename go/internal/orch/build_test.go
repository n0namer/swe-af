package orch

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Agent-Field/agentfield/sdk/go/agent"

	"github.com/Agent-Field/SWE-AF/go/internal/workspace"
)

// buildHandler routes mock reasoner responses by target suffix. Overridable
// per-reasoner via the exec/verify hooks.
func buildHandler(execResp, verifyResp func(input map[string]any) map[string]any) func(context.Context, string, map[string]any) (map[string]any, error) {
	return func(_ context.Context, target string, input map[string]any) (map[string]any, error) {
		switch {
		case strings.HasSuffix(target, ".plan"):
			return map[string]any{
				"prd": map[string]any{}, "issues": []any{}, "artifacts_dir": input["artifacts_dir"],
			}, nil
		case strings.HasSuffix(target, ".run_git_init"):
			return map[string]any{"success": false, "error_message": "no remote"}, nil
		case strings.HasSuffix(target, ".execute"):
			return execResp(input), nil
		case strings.HasSuffix(target, ".run_verifier"):
			return verifyResp(input), nil
		case strings.HasSuffix(target, ".run_repo_finalize"):
			return map[string]any{"success": true, "summary": ""}, nil
		default:
			return map[string]any{}, nil
		}
	}
}

func TestHumanApprovalEnabledByConfiguredSurface(t *testing.T) {
	t.Setenv("HAX_API_KEY", "")
	t.Setenv("SWE_OPENCLAW_HITL", "0")
	if humanApprovalEnabled() {
		t.Fatal("human approval should be disabled with neither HAX nor OpenClaw")
	}
	t.Setenv("HAX_API_KEY", "hax-test")
	if !humanApprovalEnabled() {
		t.Fatal("HAX should enable human approval")
	}
	t.Setenv("HAX_API_KEY", "")
	t.Setenv("SWE_OPENCLAW_HITL", "1")
	if !humanApprovalEnabled() {
		t.Fatal("OpenClaw should enable human approval")
	}
}

func emptyExec(map[string]any) map[string]any {
	return map[string]any{
		"completed_issues": []any{}, "merged_branches": []any{}, "all_issues": []any{},
		"failed_issues": []any{}, "skipped_issues": []any{}, "accumulated_debt": []any{},
	}
}

// TestBuildEmptyGuardReportsFailed maps to test_empty_build_guard.py: a build
// that ships nothing and fails verification must return the SDK's result-carrying
// &agent.ReasonerFailed (Message = the failure summary, Result = the BuildResult
// map) so the async handler records status=failed while preserving the outcome.
func TestBuildEmptyGuardReportsFailed(t *testing.T) {
	defer withExecCtx("run-1", "exec-1")()

	app := &mockApp{handler: buildHandler(emptyExec, func(map[string]any) map[string]any {
		return map[string]any{"passed": false, "criteria_results": []any{}, "summary": "nope"}
	})}
	deps := &Deps{App: app, NodeID: "swe-planner"}

	out, err := Build(context.Background(), deps, map[string]any{
		"goal":      "do a thing",
		"repo_path": t.TempDir(),
		"config":    map[string]any{"git_init_max_retries": 1},
	})
	if err == nil {
		t.Fatal("empty build must return an error")
	}
	var rf *agent.ReasonerFailed
	if !errors.As(err, &rf) {
		t.Fatalf("error must be *agent.ReasonerFailed, got %T: %v", err, err)
	}
	if rf.Message != "Build failed: 0/0 issues completed, no branches merged" {
		t.Fatalf("ReasonerFailed.Message = %q", rf.Message)
	}
	res, ok := rf.Result.(map[string]any)
	if !ok {
		t.Fatalf("ReasonerFailed.Result not a map: %T", rf.Result)
	}
	if asBool(res["success"]) {
		t.Fatal("carried result success should be false")
	}
	if _, has := res["dag_state"]; !has {
		t.Fatal("carried result must carry dag_state")
	}
	// Build also returns the same BuildResult map as the value.
	if m, ok := out.(map[string]any); !ok || asBool(m["success"]) {
		t.Fatalf("build return should be a BuildResult map with success=false, got %T", out)
	}
}

// TestBuildPartialNotEmpty: a build that completed an issue and merged a branch
// is NOT empty even when verification fails — it returns normally (no error).
func TestBuildPartialNotEmpty(t *testing.T) {
	defer withExecCtx("run-2", "exec-2")()

	exec := func(map[string]any) map[string]any {
		return map[string]any{
			"completed_issues": []any{map[string]any{"name": "i1"}},
			"merged_branches":  []any{"issue/x"},
			"all_issues":       []any{map[string]any{"name": "i1"}},
			"failed_issues":    []any{}, "skipped_issues": []any{}, "accumulated_debt": []any{},
		}
	}
	app := &mockApp{handler: buildHandler(exec, func(map[string]any) map[string]any {
		return map[string]any{"passed": false, "criteria_results": []any{}, "summary": "partial"}
	})}
	deps := &Deps{App: app, NodeID: "swe-planner"}

	out, err := Build(context.Background(), deps, map[string]any{
		"goal":      "thing",
		"repo_path": t.TempDir(),
		"config":    map[string]any{"git_init_max_retries": 1},
	})
	if err != nil {
		t.Fatalf("partial build should not error: %v", err)
	}
	m := out.(map[string]any)
	if asBool(m["success"]) {
		t.Fatal("success should be false for a failed-verification partial build")
	}
}

// TestBuildVerifiedSuccess: verification passes → success true, no carrier.
func TestBuildVerifiedSuccess(t *testing.T) {
	defer withExecCtx("run-3", "exec-3")()

	exec := func(map[string]any) map[string]any {
		return map[string]any{
			"completed_issues": []any{map[string]any{"name": "i1"}},
			"merged_branches":  []any{"issue/x"},
			"all_issues":       []any{map[string]any{"name": "i1"}},
			"failed_issues":    []any{}, "skipped_issues": []any{}, "accumulated_debt": []any{},
		}
	}
	app := &mockApp{handler: buildHandler(exec, func(map[string]any) map[string]any {
		return map[string]any{"passed": true, "criteria_results": []any{}, "summary": "ok"}
	})}
	deps := &Deps{App: app, NodeID: "swe-planner"}
	out, err := Build(context.Background(), deps, map[string]any{
		"goal": "thing", "repo_path": t.TempDir(),
		"config": map[string]any{"git_init_max_retries": 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !asBool(out.(map[string]any)["success"]) {
		t.Fatal("success should be true")
	}
}

func TestBuildStopsImmediatelyOnAmbiguousExecuteEffect(t *testing.T) {
	defer withExecCtx("run-ambiguous", "exec-ambiguous")()

	var verifyCalls, finalizeCalls int
	app := &mockApp{handler: func(_ context.Context, target string, input map[string]any) (map[string]any, error) {
		switch {
		case strings.HasSuffix(target, ".plan"):
			return map[string]any{
				"prd": map[string]any{}, "issues": []any{}, "artifacts_dir": input["artifacts_dir"],
			}, nil
		case strings.HasSuffix(target, ".run_git_init"):
			return map[string]any{"success": false, "error_message": "no remote"}, nil
		case strings.HasSuffix(target, ".execute"):
			return nil, errors.New("AMBIGUOUS_EFFECT: coder timed out after possible mutation")
		case strings.HasSuffix(target, ".run_verifier"):
			verifyCalls++
			return map[string]any{"passed": true}, nil
		case strings.HasSuffix(target, ".run_repo_finalize"):
			finalizeCalls++
			return map[string]any{"success": true}, nil
		default:
			return map[string]any{}, nil
		}
	}}

	out, err := Build(context.Background(), &Deps{App: app, NodeID: "swe-planner"}, map[string]any{
		"goal": "ambiguous mutation", "repo_path": t.TempDir(),
		"config": map[string]any{"git_init_max_retries": 1, "enable_github_pr": false, "check_ci": false},
	})
	if err == nil || !strings.Contains(err.Error(), "AMBIGUOUS_EFFECT") {
		t.Fatalf("Build result=%v err=%v, want AMBIGUOUS_EFFECT", out, err)
	}
	if verifyCalls != 0 || finalizeCalls != 0 {
		t.Fatalf("downstream stages ran after ambiguous effect: verify=%d finalize=%d", verifyCalls, finalizeCalls)
	}
}

func TestBuildVerifierFailureGeneratesFixAndReverifies(t *testing.T) {
	defer withExecCtx("run-fix", "exec-fix")()

	var executeCalls, verifyCalls, generateFixCalls int
	app := &mockApp{handler: func(_ context.Context, target string, input map[string]any) (map[string]any, error) {
		switch {
		case strings.HasSuffix(target, ".plan"):
			return map[string]any{
				"prd":           map[string]any{"acceptance_criteria": []any{"AC-1"}},
				"issues":        []any{},
				"artifacts_dir": input["artifacts_dir"],
			}, nil
		case strings.HasSuffix(target, ".run_git_init"):
			return map[string]any{"success": false, "error_message": "no remote"}, nil
		case strings.HasSuffix(target, ".execute"):
			executeCalls++
			name := "initial"
			if executeCalls > 1 {
				name = "fix-ac-1"
			}
			return map[string]any{
				"completed_issues": []any{map[string]any{"name": name}},
				"merged_branches":  []any{"issue/" + name},
				"all_issues":       []any{map[string]any{"name": name}},
				"failed_issues":    []any{},
				"skipped_issues":   []any{},
				"accumulated_debt": []any{},
			}, nil
		case strings.HasSuffix(target, ".run_verifier"):
			verifyCalls++
			if verifyCalls == 1 {
				return map[string]any{
					"passed": false,
					"criteria_results": []any{map[string]any{
						"criterion": "AC-1", "passed": false, "evidence": "missing",
					}},
					"summary": "needs fix",
				}, nil
			}
			return map[string]any{"passed": true, "criteria_results": []any{}, "summary": "fixed"}, nil
		case strings.HasSuffix(target, ".generate_fix_issues"):
			generateFixCalls++
			failed := asMapList(input["failed_criteria"])
			if len(failed) != 1 || mapStr(failed[0], "criterion", "") != "AC-1" {
				t.Fatalf("generate_fix_issues failed_criteria = %v", input["failed_criteria"])
			}
			return map[string]any{
				"fix_issues": []any{map[string]any{
					"name": "fix-ac-1", "title": "Fix AC-1", "description": "repair acceptance",
				}},
				"debt_items": []any{},
			}, nil
		case strings.HasSuffix(target, ".run_repo_finalize"):
			return map[string]any{"success": true, "summary": ""}, nil
		default:
			return map[string]any{}, nil
		}
	}}

	out, err := Build(context.Background(), &Deps{App: app, NodeID: "swe-planner"}, map[string]any{
		"goal":      "repair verifier failure",
		"repo_path": t.TempDir(),
		"config": map[string]any{
			"git_init_max_retries":  1,
			"max_verify_fix_cycles": 1,
			"enable_github_pr":      false,
			"check_ci":              false,
		},
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if !asBool(out.(map[string]any)["success"]) {
		t.Fatalf("Build success=false after verifier repair: %v", out)
	}
	if executeCalls != 2 {
		t.Fatalf("execute calls = %d, want initial + fix", executeCalls)
	}
	if verifyCalls != 2 {
		t.Fatalf("verifier calls = %d, want fail + reverify", verifyCalls)
	}
	if generateFixCalls != 1 {
		t.Fatalf("generate_fix_issues calls = %d, want 1", generateFixCalls)
	}
}

// TestBuildRequiresRepoPathOrURL maps to the ValueError branch.
func TestBuildRequiresRepoPathOrURL(t *testing.T) {
	defer withExecCtx("r", "e")()
	deps := &Deps{App: &mockApp{handler: func(context.Context, string, map[string]any) (map[string]any, error) {
		return map[string]any{}, nil
	}}, NodeID: "swe-planner"}
	_, err := Build(context.Background(), deps, map[string]any{"goal": "x"})
	if err == nil || !strings.Contains(err.Error(), "Either repo_path or repo_url") {
		t.Fatalf("expected repo_path/url error, got %v", err)
	}
}

// TestBuildIsolationConcurrent: two concurrent builds must receive distinct
// build_ids (no shared mutable state). Maps to test_build_isolation.py intent.
func TestBuildIsolationConcurrent(t *testing.T) {
	defer withExecCtx("run", "exec")()

	var mu sync.Mutex
	buildIDs := map[string]bool{}
	exec := func(input map[string]any) map[string]any {
		if id, ok := input["build_id"].(string); ok {
			mu.Lock()
			buildIDs[id] = true
			mu.Unlock()
		}
		return map[string]any{
			"completed_issues": []any{map[string]any{"name": "i"}},
			"merged_branches":  []any{"b"}, "all_issues": []any{map[string]any{"name": "i"}},
			"failed_issues": []any{}, "skipped_issues": []any{}, "accumulated_debt": []any{},
		}
	}
	verify := func(map[string]any) map[string]any {
		return map[string]any{"passed": true, "criteria_results": []any{}, "summary": "ok"}
	}
	app := &mockApp{handler: buildHandler(exec, verify)}

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			deps := &Deps{App: app, NodeID: "swe-planner"}
			_, _ = Build(context.Background(), deps, map[string]any{
				"goal": "g", "repo_path": t.TempDir(),
				"config": map[string]any{"git_init_max_retries": 1},
			})
		}()
	}
	wg.Wait()
	if len(buildIDs) != 2 {
		t.Fatalf("expected 2 distinct build_ids across concurrent builds, got %d: %v", len(buildIDs), buildIDs)
	}
}

// TestBuildScopedPathIncludesBuildID maps to test_build_isolation.py's
// two-builds-same-repo assertion (derived paths differ by build_id).
func TestBuildScopedPathIncludesBuildID(t *testing.T) {
	repoURL := "https://github.com/example/my-repo.git"
	name := deriveRepoName(repoURL)
	a := filepath.Join(workspace.Root(), fmt.Sprintf("%s-%s", name, newBuildID()))
	b := filepath.Join(workspace.Root(), fmt.Sprintf("%s-%s", name, newBuildID()))
	if a == b {
		t.Fatal("two builds for the same repo must produce different workspace paths")
	}
	if !strings.Contains(a, name) || !strings.Contains(b, name) {
		t.Fatalf("derived paths must contain repo name: %s %s", a, b)
	}
}

func TestNewBuildIDFormat(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		id := newBuildID()
		if len(id) != 8 {
			t.Fatalf("build_id length = %d, want 8 (%q)", len(id), id)
		}
		for _, c := range id {
			if !strings.ContainsRune("0123456789abcdef", c) {
				t.Fatalf("build_id not hex: %q", id)
			}
		}
		seen[id] = true
	}
	if len(seen) < 190 {
		t.Fatalf("build_ids not sufficiently unique: %d/200", len(seen))
	}
}
