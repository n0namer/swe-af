package orch

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type resumeGoalWaveCall struct {
	target string
	input  map[string]any
}

type resumeGoalWaveApp struct {
	calls  []resumeGoalWaveCall
	result map[string]any
}

func (a *resumeGoalWaveApp) Call(_ context.Context, target string, input map[string]any) (map[string]any, error) {
	a.calls = append(a.calls, resumeGoalWaveCall{target: target, input: input})
	return a.result, nil
}

func (a *resumeGoalWaveApp) Note(context.Context, string, ...string) {}

func writeResumeGoalWaveCheckpoint(t *testing.T, repo string) string {
	t.Helper()
	artifactsDir := filepath.Join(repo, ".artifacts")
	executionDir := filepath.Join(artifactsDir, "execution")
	if err := os.MkdirAll(executionDir, 0o755); err != nil {
		t.Fatal(err)
	}

	checkpoint := map[string]any{
		"artifacts_dir": artifactsDir,
		"all_issues": []any{
			map[string]any{"name": "prepare", "depends_on": []any{}},
			map[string]any{"name": "implement", "depends_on": []any{"prepare"}},
			map[string]any{"name": "verify", "depends_on": []any{"implement"}},
		},
		"levels": []any{
			[]any{"prepare"},
			[]any{"implement"},
			[]any{"verify"},
		},
		"current_level": 1,
		"completed_issues": []any{
			map[string]any{
				"issue_name":     "prepare",
				"outcome":        "completed",
				"result_summary": "accepted setup effect",
			},
		},
		"original_plan_summary": "prepare, implement, then verify one issue",
	}
	raw, err := json.Marshal(checkpoint)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(executionDir, "checkpoint.json")
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	return artifactsDir
}

func readyWaveOracle(issues, levels []any, completed map[string]bool) []string {
	dependencies := make(map[string][]string, len(issues))
	for _, rawIssue := range issues {
		issue, ok := rawIssue.(map[string]any)
		if !ok {
			continue
		}
		name, ok := issue["name"].(string)
		if !ok || name == "" {
			continue
		}
		dependenciesRaw, ok := issue["depends_on"].([]any)
		if !ok {
			continue
		}
		for _, rawDependency := range dependenciesRaw {
			if dependency, ok := rawDependency.(string); ok {
				dependencies[name] = append(dependencies[name], dependency)
			}
		}
	}

	for _, rawLevel := range levels {
		level, ok := rawLevel.([]any)
		if !ok {
			continue
		}
		ready := make([]string, 0, len(level))
		for _, rawName := range level {
			name, ok := rawName.(string)
			if !ok || completed[name] {
				continue
			}
			allDependenciesCompleted := true
			for _, dependency := range dependencies[name] {
				if !completed[dependency] {
					allDependenciesCompleted = false
					break
				}
			}
			if allDependenciesCompleted {
				ready = append(ready, name)
			}
		}
		if len(ready) > 0 {
			return ready
		}
	}
	return nil
}

func TestResumeGoalWaveCancelledRequestDoesNotDispatch(t *testing.T) {
	repo := t.TempDir()
	writeResumeGoalWaveCheckpoint(t, repo)
	app := &resumeGoalWaveApp{result: map[string]any{"status": "succeeded"}}
	deps := &Deps{App: app, NodeID: "swe-planner"}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	out, err := ResumeBuildHandler(ctx, deps, map[string]any{"repo_path": repo})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if out != nil {
		t.Fatalf("output = %v, want nil", out)
	}
	if len(app.calls) != 0 {
		t.Fatalf("cancelled resume dispatched %d execute call(s)", len(app.calls))
	}
}

func TestResumeGoalWaveReturnsAmbiguousWorkerResultWithoutAcceptance(t *testing.T) {
	repo := t.TempDir()
	writeResumeGoalWaveCheckpoint(t, repo)
	result := map[string]any{
		"status": "timed_out",
		"result": map[string]any{
			"worker_state": "lost",
			"dispatched":   []any{"implement"},
		},
	}
	app := &resumeGoalWaveApp{result: result}
	deps := &Deps{App: app, NodeID: "swe-planner"}

	out, err := ResumeBuildHandler(context.Background(), deps, map[string]any{"repo_path": repo})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out, result) {
		t.Fatalf("resume must preserve an ambiguous execute result, got %v", out)
	}
	if len(app.calls) != 1 {
		t.Fatalf("execute calls = %d, want 1", len(app.calls))
	}
	if status := out.(map[string]any)["status"]; status != "timed_out" {
		t.Fatalf("resume synthesized status %v, want timed_out", status)
	}
}

func TestResumeGoalWaveRestartsFromCheckpointWithoutSynthesizingAcceptance(t *testing.T) {
	repo := t.TempDir()
	artifactsDir := writeResumeGoalWaveCheckpoint(t, repo)
	result := map[string]any{
		"status": "succeeded",
		"result": map[string]any{
			"completed_issues": []any{map[string]any{"issue_name": "prepare", "outcome": "completed"}},
			"dispatched":       []any{"implement"},
		},
	}
	app := &resumeGoalWaveApp{result: result}
	deps := &Deps{App: app, NodeID: "swe-planner"}

	out, err := ResumeBuildHandler(context.Background(), deps, map[string]any{"repo_path": repo})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out, result) {
		t.Fatalf("resume must preserve the raw execute result, got %v", out)
	}
	if len(app.calls) != 1 {
		t.Fatalf("execute calls = %d, want 1", len(app.calls))
	}
	call := app.calls[0]
	if call.target != "swe-planner.execute" {
		t.Fatalf("target = %q, want swe-planner.execute", call.target)
	}
	if call.input["resume"] != true {
		t.Fatalf("resume = %v, want true", call.input["resume"])
	}

	plan, ok := call.input["plan_result"].(map[string]any)
	if !ok {
		t.Fatalf("plan_result type = %T, want map[string]any", call.input["plan_result"])
	}
	if got := mapStr(plan, "artifacts_dir", ""); got != artifactsDir {
		t.Fatalf("artifacts_dir = %q, want persisted checkpoint dir %q", got, artifactsDir)
	}
	levels, ok := plan["levels"].([]any)
	if !ok {
		t.Fatalf("levels type = %T, want []any", plan["levels"])
	}
	issues, ok := plan["issues"].([]any)
	if !ok {
		t.Fatalf("issues type = %T, want []any", plan["issues"])
	}
	if got := readyWaveOracle(issues, levels, map[string]bool{"prepare": true}); !reflect.DeepEqual(got, []string{"implement"}) {
		t.Fatalf("independent ready-node oracle = %v, want [implement]", got)
	}
	if got := readyWaveOracle(issues, []any{[]any{"implement", "verify"}}, map[string]bool{"prepare": true}); !reflect.DeepEqual(got, []string{"implement"}) {
		t.Fatalf("dependency-aware ready-node oracle = %v, want [implement]", got)
	}

	rawResult, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("result type = %T, want map[string]any", out)
	}
	if dispatched := rawResult["result"].(map[string]any)["dispatched"]; !reflect.DeepEqual(dispatched, []any{"implement"}) {
		t.Fatalf("typed execute result dispatched = %v, want [implement]", dispatched)
	}
}
