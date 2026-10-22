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

type mvpR9ResumeCall struct {
	target string
	input  map[string]any
}

type mvpR9ResumeApp struct {
	calls  []mvpR9ResumeCall
	result map[string]any
}

func (a *mvpR9ResumeApp) Call(_ context.Context, target string, input map[string]any) (map[string]any, error) {
	a.calls = append(a.calls, mvpR9ResumeCall{target: target, input: input})
	return a.result, nil
}

func (a *mvpR9ResumeApp) Note(context.Context, string, ...string) {}

func writeMVPR9Checkpoint(t *testing.T, repoPath, issueName string) string {
	t.Helper()
	artifactsDir := filepath.Join(repoPath, ".artifacts")
	checkpointPath := filepath.Join(artifactsDir, "execution", "checkpoint.json")
	if err := os.MkdirAll(filepath.Dir(checkpointPath), 0o755); err != nil {
		t.Fatal(err)
	}
	checkpoint := map[string]any{
		"artifacts_dir":         artifactsDir,
		"all_issues":            []any{map[string]any{"name": issueName, "depends_on": []any{}}},
		"levels":                []any{[]any{issueName}},
		"completed_issues":      []any{map[string]any{"issue_name": "already-done", "outcome": "completed"}},
		"original_plan_summary": "r9 adversarial replay fixture",
		"last_resume": map[string]any{
			"owner_node_id": "swe-planner",
			"resume_id":     "resume-001",
			"effect_id":     "effect-001",
			"status":        "interrupted",
		},
	}
	raw, err := json.Marshal(checkpoint)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(checkpointPath, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	return checkpointPath
}

func TestMVPR9ResumeNativePositivePreservesTypedPersistedResult(t *testing.T) {
	repoPath := t.TempDir()
	writeMVPR9Checkpoint(t, repoPath, "native-positive")
	result := map[string]any{
		"status": "succeeded",
		"result": map[string]any{
			"persisted": true,
			"effects":   []any{"native-positive"},
		},
	}
	app := &mvpR9ResumeApp{result: result}
	deps := &Deps{App: app, NodeID: "swe-planner"}

	got, err := ResumeBuildHandler(context.Background(), deps, map[string]any{"repo_path": repoPath})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, result) {
		t.Fatalf("resume result = %#v, want persisted typed result %#v", got, result)
	}
	if len(app.calls) != 1 || app.calls[0].target != "swe-planner.execute" {
		t.Fatalf("execute calls = %#v, want one swe-planner.execute call", app.calls)
	}
	if app.calls[0].input["resume"] != true {
		t.Fatalf("resume flag = %#v, want true", app.calls[0].input["resume"])
	}
}

func TestMVPR9ResumeCanceledRequestFailsClosedBeforeDispatch(t *testing.T) {
	repoPath := t.TempDir()
	writeMVPR9Checkpoint(t, repoPath, "canceled")
	app := &mvpR9ResumeApp{result: map[string]any{"status": "unexpected"}}
	deps := &Deps{App: app, NodeID: "swe-planner"}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := ResumeBuildHandler(ctx, deps, map[string]any{"repo_path": repoPath})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("resume error = %v, want context.Canceled", err)
	}
	if got != nil {
		t.Fatalf("resume result = %#v, want nil", got)
	}
	if len(app.calls) != 0 {
		t.Fatalf("canceled resume dispatched %d execute calls", len(app.calls))
	}
}

func TestMVPR9ResumeInvalidPayloadFailsClosedBeforeDispatch(t *testing.T) {
	app := &mvpR9ResumeApp{result: map[string]any{"status": "unexpected"}}
	deps := &Deps{App: app, NodeID: "swe-planner"}

	got, err := ResumeBuildHandler(context.Background(), deps, map[string]any{})
	if err == nil {
		t.Fatal("invalid resume payload unexpectedly succeeded")
	}
	if got != nil {
		t.Fatalf("invalid resume result = %#v, want nil", got)
	}
	if len(app.calls) != 0 {
		t.Fatalf("invalid resume dispatched %d execute calls", len(app.calls))
	}
}

func TestMVPR9ResumeForeignRepositoryCannotReplacePinnedCheckpoint(t *testing.T) {
	ownerRepo := t.TempDir()
	foreignRepo := t.TempDir()
	ownerCheckpoint := writeMVPR9Checkpoint(t, ownerRepo, "owner-only")
	foreignCheckpoint := writeMVPR9Checkpoint(t, foreignRepo, "foreign-only")
	foreignBefore, err := os.ReadFile(foreignCheckpoint)
	if err != nil {
		t.Fatal(err)
	}
	app := &mvpR9ResumeApp{result: map[string]any{"status": "succeeded"}}
	deps := &Deps{App: app, NodeID: "swe-planner"}

	if _, err := ResumeBuildHandler(context.Background(), deps, map[string]any{
		"repo_path": ownerRepo,
	}); err != nil {
		t.Fatal(err)
	}
	if len(app.calls) != 1 {
		t.Fatalf("execute calls = %d, want 1", len(app.calls))
	}
	plan, ok := app.calls[0].input["plan_result"].(map[string]any)
	if !ok {
		t.Fatalf("plan_result type = %T, want map[string]any", app.calls[0].input["plan_result"])
	}
	issues, ok := plan["issues"].([]any)
	if !ok || len(issues) != 1 {
		t.Fatalf("pinned issues = %#v, want one owner issue", plan["issues"])
	}
	issue, ok := issues[0].(map[string]any)
	if !ok || issue["name"] != "owner-only" {
		t.Fatalf("pinned issue = %#v, want owner-only", issues[0])
	}
	lastResume, ok := plan["last_resume"].(map[string]any)
	if !ok {
		t.Fatalf("last_resume type = %T, want map[string]any", plan["last_resume"])
	}
	wantLastResume := map[string]any{
		"owner_node_id": "swe-planner",
		"resume_id":     "resume-001",
		"effect_id":     "effect-001",
		"status":        "interrupted",
	}
	if !reflect.DeepEqual(lastResume, wantLastResume) {
		t.Fatalf("last_resume = %#v, want %#v", lastResume, wantLastResume)
	}
	foreignAfter, err := os.ReadFile(foreignCheckpoint)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(foreignAfter, foreignBefore) {
		t.Fatalf("foreign checkpoint changed during owner resume: before %q after %q", foreignBefore, foreignAfter)
	}
	ownerAfter, err := os.ReadFile(ownerCheckpoint)
	if err != nil {
		t.Fatal(err)
	}
	if len(ownerAfter) == 0 {
		t.Fatal("owner checkpoint was unexpectedly emptied")
	}
}

func TestMVPR9ResumeRestartReplaysWithoutCheckpointMutationOrExtraDispatch(t *testing.T) {
	repoPath := t.TempDir()
	checkpointPath := writeMVPR9Checkpoint(t, repoPath, "restart-safe")
	before, err := os.ReadFile(checkpointPath)
	if err != nil {
		t.Fatal(err)
	}
	app := &mvpR9ResumeApp{result: map[string]any{
		"status": "succeeded",
		"result": map[string]any{"persisted": true},
	}}
	deps := &Deps{App: app, NodeID: "swe-planner"}
	input := map[string]any{"repo_path": repoPath}

	first, err := ResumeBuildHandler(context.Background(), deps, input)
	if err != nil {
		t.Fatal(err)
	}
	second, err := ResumeBuildHandler(context.Background(), deps, input)
	if err != nil {
		t.Fatal(err)
	}
	if len(app.calls) != 2 {
		t.Fatalf("restart dispatch count = %d, want one dispatch per resume attempt", len(app.calls))
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("restart result changed: first %#v second %#v", first, second)
	}
	firstPlan := app.calls[0].input["plan_result"]
	secondPlan := app.calls[1].input["plan_result"]
	if !reflect.DeepEqual(firstPlan, secondPlan) {
		t.Fatalf("restart reconstructed different plans: first %#v second %#v", firstPlan, secondPlan)
	}
	after, err := os.ReadFile(checkpointPath)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("restart mutated durable checkpoint: before %q after %q", before, after)
	}
}
