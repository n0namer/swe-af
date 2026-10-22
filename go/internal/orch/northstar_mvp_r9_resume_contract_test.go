package orch

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type mvpR9Runner struct {
	effects map[string]int
	calls   int
}

func (runner *mvpR9Runner) Call(_ context.Context, target string, input map[string]any) (map[string]any, error) {
	if target != "swe-planner.execute" {
		return nil, errors.New("unexpected target")
	}
	runner.calls++
	for _, raw := range input["resume_dispatch"].([]map[string]any) {
		key := mapStr(raw, "key", "")
		if runner.effects[key] == 0 {
			runner.effects[key]++
		}
	}
	return map[string]any{
		"status": "succeeded",
		"result": map[string]any{"typed": true, "effects": len(runner.effects)},
	}, nil
}

func (*mvpR9Runner) Note(context.Context, string, ...string) {}

func writeMVP9Checkpoint(t *testing.T, owner string, epoch int) string {
	t.Helper()
	repo := t.TempDir()
	execution := filepath.Join(repo, ".artifacts", "execution")
	if err := os.MkdirAll(execution, 0o755); err != nil {
		t.Fatal(err)
	}
	checkpoint := map[string]any{
		"artifacts_dir": filepath.Join(repo, ".artifacts"),
		"owner_id":      owner,
		"lease_epoch":   epoch,
		"all_issues": []any{
			map[string]any{"name": "prepare", "depends_on": []any{}},
			map[string]any{"name": "implement", "depends_on": []any{"prepare"}},
		},
		"levels": []any{[]any{"prepare"}, []any{"implement"}},
		"completed_issues": []any{map[string]any{
			"issue_name": "prepare",
			"outcome":    "completed",
		}},
	}
	raw, err := json.Marshal(checkpoint)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(execution, "checkpoint.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	return repo
}

func TestMVPR9ResumePersistsTypedResultAndOnlyDependencyReadyEffect(t *testing.T) {
	repo := writeMVP9Checkpoint(t, "owner-a", 7)
	runner := &mvpR9Runner{effects: map[string]int{}}
	deps := &Deps{App: runner, NodeID: "swe-planner"}

	result, err := ResumeBuildHandler(context.Background(), deps, map[string]any{
		"repo_path":  repo,
		"owner_id":   "owner-a",
		"lease_epoch": 7,
		"run_id":     "run-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.(map[string]any)["status"] != "succeeded" || runner.calls != 1 || len(runner.effects) != 1 {
		t.Fatalf("result=%v calls=%d effects=%d", result, runner.calls, len(runner.effects))
	}

	raw, err := os.ReadFile(filepath.Join(repo, ".artifacts", "execution", "checkpoint.json"))
	if err != nil {
		t.Fatal(err)
	}
	var persisted map[string]any
	if err := json.Unmarshal(raw, &persisted); err != nil {
		t.Fatal(err)
	}
	if persisted["last_resume_result"].(map[string]any)["status"] != "succeeded" {
		t.Fatalf("typed result was not persisted: %v", persisted["last_resume_result"])
	}
	lastDispatch := persisted["last_resume_dispatch"].([]any)
	if len(lastDispatch) != 1 || lastDispatch[0].(map[string]any)["issue"] != "implement" {
		t.Fatalf("ready dispatch=%v, want implement only", lastDispatch)
	}
}

func TestMVPR9ResumeRejectsStaleOwnerWithoutEffect(t *testing.T) {
	repo := writeMVP9Checkpoint(t, "owner-a", 7)
	runner := &mvpR9Runner{effects: map[string]int{}}
	_, err := ResumeBuildHandler(context.Background(), &Deps{App: runner, NodeID: "swe-planner"}, map[string]any{
		"repo_path": repo, "owner_id": "owner-b", "lease_epoch": 7,
	})
	if !errors.Is(err, errResumeStaleOwner) {
		t.Fatalf("error=%v, want stale owner", err)
	}
	if runner.calls != 0 || len(runner.effects) != 0 {
		t.Fatalf("stale owner caused effect: calls=%d effects=%d", runner.calls, len(runner.effects))
	}
}

func TestMVPR9ResumeCancellationAndReplayAreFailClosedAndIdempotent(t *testing.T) {
	repo := writeMVP9Checkpoint(t, "owner-a", 7)
	runner := &mvpR9Runner{effects: map[string]int{}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := ResumeBuildHandler(ctx, &Deps{App: runner, NodeID: "swe-planner"}, map[string]any{
		"repo_path": repo, "owner_id": "owner-a", "lease_epoch": 7,
	}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel error=%v, want context.Canceled", err)
	}
	if runner.calls != 0 {
		t.Fatalf("cancelled resume dispatched %d calls", runner.calls)
	}

	input := map[string]any{"repo_path": repo, "owner_id": "owner-a", "lease_epoch": 7, "run_id": "run-1"}
	if _, err := ResumeBuildHandler(context.Background(), &Deps{App: runner, NodeID: "swe-planner"}, input); err != nil {
		t.Fatal(err)
	}
	if _, err := ResumeBuildHandler(context.Background(), &Deps{App: runner, NodeID: "swe-planner"}, input); err != nil {
		t.Fatal(err)
	}
	if len(runner.effects) != 1 {
		t.Fatalf("replay duplicated effect: %v", runner.effects)
	}
}
