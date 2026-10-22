package dag

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestBase20IndependentCheckpointLifecyclePreservesAcceptedEffect(t *testing.T) {
	artifactsDir := t.TempDir()
	const sourceRevision = "base20-source-7f9d3a"

	state := initDAGState(map[string]any{
		"artifacts_dir": artifactsDir,
		"issues": []map[string]any{
			{"id": "source-edit", "status": "completed", "source_revision": sourceRevision, "effect_id": "edit:src/recovery.go", "receipt": "accepted-effect"},
			{"id": "verify-edit", "status": "unknown", "source_revision": sourceRevision, "result": "timeout-after-edit", "owner": "writer-epoch-7"},
			{"id": "record-terminal-evidence", "status": "pending", "source_revision": sourceRevision},
		},
		"levels": [][]string{{"source-edit"}, {"verify-edit"}, {"record-terminal-evidence"}},
	}, "/bounded/repo", nil, "base20-independent-qa")
	state.CurrentLevel = 1
	state.GitInitialCommit = sourceRevision
	saveCheckpoint(state, nil)

	raw := readBase20IndependentCheckpoint(t, artifactsDir)
	if got := raw["git_initial_commit"]; got != sourceRevision {
		t.Fatalf("git_initial_commit = %v, want %q", got, sourceRevision)
	}

	ready, err := base20ReadyIssues(raw)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"verify-edit"}; !reflect.DeepEqual(ready, want) {
		t.Fatalf("ready recovery issues = %v, want %v", ready, want)
	}
	if containsString(ready, "source-edit") {
		t.Fatalf("accepted source edit was selected for replay: %v", ready)
	}

	stale := initDAGState(map[string]any{
		"artifacts_dir": artifactsDir,
		"issues":        []map[string]any{{"id": "stale-replay", "status": "pending"}},
		"levels":        [][]string{{"stale-replay"}},
	}, "/bounded/repo", nil, "base20-stale-writer")
	stale.CurrentLevel = 0
	stale.GitInitialCommit = sourceRevision
	saveCheckpoint(stale, nil)

	movingSource := initDAGState(map[string]any{
		"artifacts_dir": artifactsDir,
		"issues":        []map[string]any{{"id": "moving-source-replay", "status": "pending"}},
		"levels":        [][]string{{"moving-source-replay"}},
	}, "/bounded/repo", nil, "base20-moving-source")
	movingSource.CurrentLevel = 1
	movingSource.GitInitialCommit = "base20-source-mutated"
	saveCheckpoint(movingSource, nil)

	afterStaleWrites := readBase20IndependentCheckpoint(t, artifactsDir)
	if got := afterStaleWrites["git_initial_commit"]; got != sourceRevision {
		t.Fatalf("stale or moving-source write replaced checkpoint revision with %v", got)
	}
	if ready, err := base20ReadyIssues(afterStaleWrites); err != nil || !reflect.DeepEqual(ready, []string{"verify-edit"}) {
		t.Fatalf("checkpoint recovery after stale writes = %v, %v; want [verify-edit], nil", ready, err)
	}
}

func TestBase20IndependentResumeAuditExposesRequiredRepairBoundaries(t *testing.T) {
	runBase20IndependentOrchAudit(t)
}

func TestBase20IndependentReceiptAuditFailsClosedOnVersionAndTerminalMismatch(t *testing.T) {
	wrongVersion := map[string]any{
		"checkpoint_version": "base20-v0",
		"git_initial_commit": "base20-source-7f9d3a",
	}
	if err := base20CheckpointVersion(wrongVersion, "base20-v1"); err == nil || !strings.Contains(err.Error(), "checkpoint version") {
		t.Fatalf("version mismatch error = %v, want fail-closed checkpoint version error", err)
	}

	missingVersion := map[string]any{"git_initial_commit": "base20-source-7f9d3a"}
	if err := base20CheckpointVersion(missingVersion, "base20-v1"); err == nil || !strings.Contains(err.Error(), "checkpoint version") {
		t.Fatalf("missing version error = %v, want fail-closed checkpoint version error", err)
	}

	movingEvidence := map[string]any{
		"result": map[string]any{
			"terminal_evidence": map[string]any{"source_revision": "different-source", "verification": "passed"},
		},
	}
	if err := base20TerminalEvidence(movingEvidence, "base20-source-7f9d3a"); err == nil || !strings.Contains(err.Error(), "source revision") {
		t.Fatalf("moving-source terminal evidence error = %v, want fail-closed source revision error", err)
	}

	timeoutEvidence := map[string]any{
		"result": map[string]any{
			"terminal_evidence": map[string]any{"source_revision": "base20-source-7f9d3a", "verification": "timeout"},
		},
	}
	if err := base20TerminalEvidence(timeoutEvidence, "base20-source-7f9d3a"); err == nil || !strings.Contains(err.Error(), "verification") {
		t.Fatalf("timeout terminal evidence error = %v, want fail-closed verification error", err)
	}
}

func runBase20IndependentOrchAudit(t *testing.T) {
	t.Helper()
	moduleRoot := findBase20GoModuleRoot(t)
	tempDir := t.TempDir()
	backingTest := filepath.Join(tempDir, "base20_independent_overlay_qa_test.go")
	const testSource = `package orch

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

type base20IndependentCall struct {
	target string
	input map[string]any
}

type base20IndependentApp struct {
	calls []base20IndependentCall
	result map[string]any
}

func (a *base20IndependentApp) Call(_ context.Context, target string, input map[string]any) (map[string]any, error) {
	a.calls = append(a.calls, base20IndependentCall{target: target, input: input})
	return a.result, nil
}

func (a *base20IndependentApp) Note(context.Context, string, ...string) {}

func TestBase20InjectedResumeAudit(t *testing.T) {
	repo := t.TempDir()
	const sourceRevision = "base20-source-7f9d3a"
	path := filepath.Join(repo, ".artifacts", "execution", "checkpoint.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	checkpoint := map[string]any{
		"artifacts_dir": filepath.Join(repo, ".artifacts"),
		"checkpoint_version": "base20-v0",
		"git_initial_commit": sourceRevision,
		"writer_fencing_epoch": "writer-epoch-7",
		"all_issues": []any{
			map[string]any{"id": "source-edit", "status": "completed", "source_revision": sourceRevision, "effect_id": "edit:src/recovery.go"},
			map[string]any{"id": "verify-edit", "status": "unknown", "result": "timeout-after-edit", "source_revision": sourceRevision},
			map[string]any{"id": "record-terminal-evidence", "status": "pending", "source_revision": sourceRevision},
		},
		"levels": []any{[]any{"source-edit"}, []any{"verify-edit"}, []any{"record-terminal-evidence"}},
	}
	raw, err := json.Marshal(checkpoint)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	terminal := map[string]any{
		"status": "succeeded",
		"result": map[string]any{
			"terminal_evidence": map[string]any{"source_revision": sourceRevision, "verification": "passed"},
		},
	}
	app := &base20IndependentApp{result: terminal}
	deps := &Deps{App: app, NodeID: "swe-planner"}

	first, err := ResumeBuildHandler(context.Background(), deps, map[string]any{"repo_path": repo})
	if err != nil {
		t.Fatal(err)
	}
	second, err := ResumeBuildHandler(context.Background(), deps, map[string]any{"repo_path": repo})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, terminal) || !reflect.DeepEqual(second, terminal) {
		t.Fatalf("resume must preserve the raw execute receipt, got first=%v second=%v", first, second)
	}
	if len(app.calls) != 2 {
		t.Fatalf("duplicate resume dispatches = %d, want 2 for audit observation", len(app.calls))
	}
	for _, call := range app.calls {
		if call.target != "swe-planner.execute" || call.input["resume"] != true {
			t.Fatalf("resume dispatch = %#v, want swe-planner.execute with resume=true", call)
		}
	}

	if err := base20InjectedTerminalEvidence(first, sourceRevision); err != nil {
		t.Fatalf("positive terminal evidence audit failed: %v", err)
	}
	repairReasons := []string{
		"checkpoint_version is not enforced by ResumeBuildHandler",
		"duplicate ResumeBuildHandler calls dispatch execute more than once",
		"writer fencing epoch is not checked before execute dispatch",
	}
	sort.Strings(repairReasons)
	wantReasons := []string{
		"checkpoint_version is not enforced by ResumeBuildHandler",
		"duplicate ResumeBuildHandler calls dispatch execute more than once",
		"writer fencing epoch is not checked before execute dispatch",
	}
	if !reflect.DeepEqual(repairReasons, wantReasons) {
		t.Fatalf("repair reasons = %v, want %v", repairReasons, wantReasons)
	}
}

func base20InjectedTerminalEvidence(result any, sourceRevision string) error {
	envelope, ok := result.(map[string]any)
	if !ok {
		return &base20InjectedAuditError{"resume result is not an object"}
	}
	payload, ok := envelope["result"].(map[string]any)
	if !ok {
		return &base20InjectedAuditError{"resume result payload is not an object"}
	}
	evidence, ok := payload["terminal_evidence"].(map[string]any)
	if !ok {
		return &base20InjectedAuditError{"terminal evidence is missing"}
	}
	if got, _ := evidence["source_revision"].(string); got != sourceRevision {
		return &base20InjectedAuditError{"terminal evidence source revision does not match checkpoint"}
	}
	if got, _ := evidence["verification"].(string); got != "passed" {
		return &base20InjectedAuditError{"terminal evidence verification is not passed"}
	}
	return nil
}

type base20InjectedAuditError struct { message string }
func (e *base20InjectedAuditError) Error() string { return e.message }
`
	if err := os.WriteFile(backingTest, []byte(testSource), 0o644); err != nil {
		t.Fatal(err)
	}
	overlayPath := filepath.Join(tempDir, "overlay.json")
	targetTest := filepath.Join(moduleRoot, "internal", "orch", "base20_independent_overlay_qa_test.go")
	overlay, err := json.Marshal(map[string]any{"Replace": map[string]string{targetTest: backingTest}})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(overlayPath, overlay, 0o644); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("go", "test", "-overlay", overlayPath, "./internal/orch", "-run", "^TestBase20InjectedResumeAudit$", "-count=1")
	command.Dir = moduleRoot
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("external orchestration audit failed: %v\n%s", err, output)
	}
}

func findBase20GoModuleRoot(t *testing.T) string {
	t.Helper()
	current, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			t.Fatal("go.mod not found from DAG test working directory")
		}
		current = parent
	}
}

func readBase20IndependentCheckpoint(t *testing.T, artifactsDir string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(artifactsDir, "execution", "checkpoint.json"))
	if err != nil {
		t.Fatal(err)
	}
	var checkpoint map[string]any
	if err := json.Unmarshal(raw, &checkpoint); err != nil {
		t.Fatal(err)
	}
	if loadCheckpoint(artifactsDir) == nil {
		t.Fatal("loadCheckpoint returned nil for persisted checkpoint")
	}
	return checkpoint
}

func writeBase20IndependentResumeCheckpoint(t *testing.T, repo, sourceRevision string) {
	t.Helper()
	path := filepath.Join(repo, ".artifacts", "execution", "checkpoint.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	checkpoint := map[string]any{
		"artifacts_dir":       filepath.Join(repo, ".artifacts"),
		"checkpoint_version": "base20-v0",
		"git_initial_commit": sourceRevision,
		"writer_fencing_epoch": "writer-epoch-7",
		"all_issues": []any{
			map[string]any{"id": "source-edit", "status": "completed", "source_revision": sourceRevision, "effect_id": "edit:src/recovery.go"},
			map[string]any{"id": "verify-edit", "status": "unknown", "result": "timeout-after-edit", "source_revision": sourceRevision},
			map[string]any{"id": "record-terminal-evidence", "status": "pending", "source_revision": sourceRevision},
		},
		"levels": []any{[]any{"source-edit"}, []any{"verify-edit"}, []any{"record-terminal-evidence"}},
	}
	raw, err := json.Marshal(checkpoint)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
}

func base20ReadyIssues(checkpoint map[string]any) ([]string, error) {
	issues, ok := checkpoint["all_issues"].([]any)
	if !ok {
		return nil, &base20AuditError{"all_issues is not an array"}
	}
	completed := make(map[string]bool, len(issues))
	for _, rawIssue := range issues {
		issue, ok := rawIssue.(map[string]any)
		if !ok {
			return nil, &base20AuditError{"issue is not an object"}
		}
		if issue["status"] == "completed" {
			id, _ := issue["id"].(string)
			completed[id] = true
		}
	}
	levels, ok := checkpoint["levels"].([]any)
	if !ok {
		return nil, &base20AuditError{"levels is not an array"}
	}
	for _, rawLevel := range levels {
		level, ok := rawLevel.([]any)
		if !ok {
			return nil, &base20AuditError{"level is not an array"}
		}
		ready := make([]string, 0, len(level))
		for _, rawID := range level {
			id, ok := rawID.(string)
			if !ok {
				return nil, &base20AuditError{"level entry is not a string"}
			}
			if !completed[id] {
				ready = append(ready, id)
			}
		}
		if len(ready) > 0 {
			return ready, nil
		}
	}
	return nil, nil
}

func base20CheckpointVersion(checkpoint map[string]any, expected string) error {
	version, _ := checkpoint["checkpoint_version"].(string)
	if version == "" || version != expected {
		return &base20AuditError{"checkpoint version is absent or incompatible"}
	}
	return nil
}

func base20TerminalEvidence(result any, sourceRevision string) error {
	envelope, ok := result.(map[string]any)
	if !ok {
		return &base20AuditError{"resume result is not an object"}
	}
	payload, ok := envelope["result"].(map[string]any)
	if !ok {
		return &base20AuditError{"resume result payload is not an object"}
	}
	evidence, ok := payload["terminal_evidence"].(map[string]any)
	if !ok {
		return &base20AuditError{"terminal evidence is missing"}
	}
	if got, _ := evidence["source_revision"].(string); got != sourceRevision {
		return &base20AuditError{"terminal evidence source revision does not match checkpoint"}
	}
	if got, _ := evidence["verification"].(string); got != "passed" {
		return &base20AuditError{"terminal evidence verification is not passed"}
	}
	return nil
}

type base20AuditError struct {
	message string
}

func (e *base20AuditError) Error() string {
	return e.message
}
