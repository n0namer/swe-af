package dag

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
)

func TestTerraHighNativeCheckpointPersistsAndSelectsRecoverableIssues(t *testing.T) {
	artifactsDir := t.TempDir()
	state := terraHighCheckpointState(artifactsDir, []map[string]any{
		{"id": "accepted-effect", "status": "completed", "effect_id": "write:src/recovery.go"},
		{"id": "retry-verification", "status": "unknown", "result": "interrupted"},
		{"id": "write-terminal-evidence", "status": "pending"},
	}, 1)

	saveCheckpoint(state, nil)
	loaded := loadCheckpoint(artifactsDir)
	if loaded == nil {
		t.Fatal("loadCheckpoint returned nil for a persisted checkpoint")
	}
	if loaded.GitInitialCommit != state.GitInitialCommit {
		t.Fatalf("GitInitialCommit = %q, want %q", loaded.GitInitialCommit, state.GitInitialCommit)
	}
	if got, want := recoverableIssueIDs(loaded), []string{"retry-verification", "write-terminal-evidence"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("recoverableIssueIDs() = %v, want %v", got, want)
	}
}

func TestTerraHighNativeCheckpointRejectsStaleOverwrite(t *testing.T) {
	artifactsDir := t.TempDir()
	persisted := terraHighCheckpointState(artifactsDir, []map[string]any{
		{"id": "accepted-effect", "status": "completed"},
		{"id": "retry-verification", "status": "unknown"},
	}, 2)
	saveCheckpoint(persisted, nil)

	stale := terraHighCheckpointState(artifactsDir, []map[string]any{
		{"id": "stale-replay", "status": "pending"},
	}, 1)
	saveCheckpoint(stale, nil)

	loaded := loadCheckpoint(artifactsDir)
	if loaded == nil {
		t.Fatal("loadCheckpoint returned nil after stale write")
	}
	if loaded.CurrentLevel != persisted.CurrentLevel {
		t.Fatalf("CurrentLevel = %d, want persisted %d", loaded.CurrentLevel, persisted.CurrentLevel)
	}
	if got, want := recoverableIssueIDs(loaded), []string{"retry-verification"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("stale checkpoint replaced persisted recovery state: got %v, want %v", got, want)
	}
}

func TestTerraHighNativeCheckpointFailsClosedForCorruptState(t *testing.T) {
	artifactsDir := t.TempDir()
	path := filepath.Join(artifactsDir, "execution", "checkpoint.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"all_issues":`), 0o644); err != nil {
		t.Fatal(err)
	}

	if loaded := loadCheckpoint(artifactsDir); loaded != nil {
		t.Fatalf("loadCheckpoint returned state for corrupt checkpoint: %#v", loaded)
	}
}

func TestTerraHighNativeCheckpointRetryExcludesAcceptedEffect(t *testing.T) {
	artifactsDir := t.TempDir()
	interrupted := terraHighCheckpointState(artifactsDir, []map[string]any{
		{"id": "accepted-effect", "status": "completed", "effect_id": "write:src/recovery.go"},
		{"id": "retry-verification", "status": "unknown", "result": "interrupted"},
	}, 1)
	saveCheckpoint(interrupted, nil)

	loaded := loadCheckpoint(artifactsDir)
	if loaded == nil {
		t.Fatal("loadCheckpoint returned nil for interrupted retry checkpoint")
	}
	if got, want := recoverableIssueIDs(loaded), []string{"retry-verification"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("retry selection = %v, want %v; accepted effect must not repeat", got, want)
	}

	loaded.AllIssues[1]["status"] = "completed"
	loaded.CurrentLevel = 2
	saveCheckpoint(loaded, nil)
	resumed := loadCheckpoint(artifactsDir)
	if resumed == nil {
		t.Fatal("loadCheckpoint returned nil after completing interrupted retry")
	}
	if got := recoverableIssueIDs(resumed); len(got) != 0 {
		t.Fatalf("completed retry remained recoverable: %v", got)
	}
}

func terraHighCheckpointState(artifactsDir string, issues []map[string]any, currentLevel int) *schemas.DAGState {
	state := initDAGState(map[string]any{
		"artifacts_dir": artifactsDir,
		"issues":        issues,
		"levels":        [][]string{{"accepted-effect"}, {"retry-verification"}, {"write-terminal-evidence"}},
	}, "/bounded/native-dag", nil, "terra-high-native-conformance")
	state.CurrentLevel = currentLevel
	state.GitInitialCommit = "064a7226a500365005b5b662044fc9014a328722"
	return state
}
