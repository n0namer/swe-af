package dag

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestWaveCheckpointPreservesCrashRecoveryFence(t *testing.T) {
	artifactsDir := t.TempDir()
	const sourceRevision = "9f8e7d6c5b4a3210"

	state := initDAGState(map[string]any{
		"artifacts_dir": artifactsDir,
		"issues": []map[string]any{
			{
				"id":              "edit-file",
				"status":          "completed",
				"source_revision": sourceRevision,
				"effect_id":       "write:src/widget.go@9f8e7d6",
			},
			{
				"id":              "run-tests",
				"status":          "unknown",
				"source_revision": sourceRevision,
				"owner":           "stale-worker",
				"result":          "ambiguous-after-crash",
			},
			{
				"id":              "accept-result",
				"status":          "pending",
				"source_revision": sourceRevision,
			},
		},
		"levels": [][]string{{"edit-file"}, {"run-tests"}, {"accept-result"}},
	}, "/bounded/repo", nil, "wave-build")
	state.CurrentLevel = 1
	state.GitInitialCommit = sourceRevision

	saveCheckpoint(state, nil)

	checkpoint := readGoalWaveCheckpoint(t, artifactsDir)
	assertCheckpointIdentity(t, checkpoint, sourceRevision)
	assertCheckpointIssues(t, checkpoint, []string{"edit-file", "run-tests", "accept-result"})

	loaded := loadCheckpoint(artifactsDir)
	if loaded == nil {
		t.Fatal("loadCheckpoint returned nil")
	}
	if got, want := recoverableIssueIDs(loaded), []string{"accept-result", "run-tests"}; !equalStrings(got, want) {
		t.Fatalf("recovery IDs = %v, want %v", got, want)
	}

	revisionMismatch := initDAGState(map[string]any{
		"artifacts_dir": artifactsDir,
		"issues":        []map[string]any{{"id": "revision-replay", "status": "pending"}},
		"levels":        [][]string{{"revision-replay"}},
	}, "/bounded/repo", nil, "revision-build")
	revisionMismatch.CurrentLevel = 1
	revisionMismatch.GitInitialCommit = "mutated-source-revision"
	saveCheckpoint(revisionMismatch, nil)

	stale := initDAGState(map[string]any{
		"artifacts_dir": artifactsDir,
		"issues":        []map[string]any{{"id": "stale-replay", "status": "pending"}},
		"levels":        [][]string{{"stale-replay"}},
	}, "/bounded/repo", nil, "stale-build")
	stale.CurrentLevel = 0
	stale.GitInitialCommit = sourceRevision
	saveCheckpoint(stale, nil)

	terminalRegression := initDAGState(map[string]any{
		"artifacts_dir": artifactsDir,
		"issues": []map[string]any{
			{"id": "edit-file", "status": "pending"},
			{"id": "run-tests", "status": "unknown"},
			{"id": "accept-result", "status": "pending"},
		},
		"levels": [][]string{{"edit-file"}, {"run-tests"}, {"accept-result"}},
	}, "/bounded/repo", nil, "stale-terminal-build")
	terminalRegression.CurrentLevel = 1
	terminalRegression.GitInitialCommit = sourceRevision
	saveCheckpoint(terminalRegression, nil)

	missingTerminal := initDAGState(map[string]any{
		"artifacts_dir": artifactsDir,
		"issues": []map[string]any{
			{"id": "run-tests", "status": "unknown"},
			{"id": "accept-result", "status": "pending"},
		},
		"levels": [][]string{{"edit-file"}, {"run-tests"}, {"accept-result"}},
	}, "/bounded/repo", nil, "missing-terminal-build")
	missingTerminal.CurrentLevel = 1
	missingTerminal.GitInitialCommit = sourceRevision
	saveCheckpoint(missingTerminal, nil)

	checkpoint = readGoalWaveCheckpoint(t, artifactsDir)
	assertCheckpointIdentity(t, checkpoint, sourceRevision)
	assertCheckpointIssues(t, checkpoint, []string{"edit-file", "run-tests", "accept-result"})
}

func readGoalWaveCheckpoint(t *testing.T, artifactsDir string) map[string]any {
	t.Helper()
	path := filepath.Join(artifactsDir, "execution", "checkpoint.json")
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read checkpoint: %v", err)
	}
	var checkpoint map[string]any
	if err := json.Unmarshal(contents, &checkpoint); err != nil {
		t.Fatalf("unmarshal checkpoint: %v", err)
	}
	if loaded := loadCheckpoint(artifactsDir); loaded == nil {
		t.Fatal("loadCheckpoint returned nil")
	}
	return checkpoint
}

func assertCheckpointIdentity(t *testing.T, checkpoint map[string]any, sourceRevision string) {
	t.Helper()
	if got := checkpoint["git_initial_commit"]; got != sourceRevision {
		t.Errorf("git_initial_commit = %v, want %q", got, sourceRevision)
	}
	if got := checkpoint["current_level"]; got != float64(1) {
		t.Errorf("current_level = %v, want 1", got)
	}
}

func assertCheckpointIssues(t *testing.T, checkpoint map[string]any, wantIDs []string) {
	t.Helper()
	issues, ok := checkpoint["all_issues"].([]any)
	if !ok {
		t.Fatalf("all_issues = %T, want JSON array", checkpoint["all_issues"])
	}
	gotIDs := make([]string, 0, len(issues))
	byID := make(map[string]map[string]any, len(issues))
	for _, rawIssue := range issues {
		issue, ok := rawIssue.(map[string]any)
		if !ok {
			t.Fatalf("issue = %T, want JSON object", rawIssue)
		}
		id, _ := issue["id"].(string)
		gotIDs = append(gotIDs, id)
		byID[id] = issue
	}
	sort.Strings(gotIDs)
	want := append([]string(nil), wantIDs...)
	sort.Strings(want)
	if !equalStrings(gotIDs, want) {
		t.Fatalf("checkpoint node IDs = %v, want %v", gotIDs, want)
	}
	if got := byID["edit-file"]["status"]; got != "completed" {
		t.Errorf("edit-file status = %v, want completed", got)
	}
	if got := byID["run-tests"]["owner"]; got != "stale-worker" {
		t.Errorf("run-tests owner = %v, want stale-worker", got)
	}
	if got := byID["run-tests"]["result"]; got != "ambiguous-after-crash" {
		t.Errorf("run-tests result = %v, want ambiguous-after-crash", got)
	}
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
