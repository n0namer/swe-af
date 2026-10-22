package dag

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
)

var (
	_ func(*schemas.DAGState, noteFunc) = saveCheckpoint
	_ func(string) *schemas.DAGState = loadCheckpoint
	_ func(*schemas.DAGState, *schemas.DAGState) bool = checkpointRegressesTerminalIssue
	_ func(*schemas.DAGState) []string = recoverableIssueIDs
)

func TestNorthstarLongBundle(t *testing.T) {
	t.Run("positive checkpoint write and restart recovery", func(t *testing.T) {
		artifactsDir := t.TempDir()
		state := &schemas.DAGState{
			ArtifactsDir:     artifactsDir,
			CurrentLevel:     4,
			GitInitialCommit: "base-sha",
			AllIssues: []map[string]any{
				{"id": "issue-b", "status": "running"},
				{"id": "issue-a", "status": "pending"},
				{"id": "issue-d", "status": "accepted"},
				{"id": "issue-c", "status": "failed"},
			},
		}

		saveCheckpoint(state, nil)

		checkpointFile := filepath.Join(artifactsDir, "execution", "checkpoint.json")
		if _, err := os.Stat(checkpointFile); err != nil {
			t.Fatalf("checkpoint file missing: %v", err)
		}

		resumed := loadCheckpoint(artifactsDir)
		if resumed == nil {
			t.Fatal("loadCheckpoint returned nil after a successful save")
		}
		if resumed.CurrentLevel != 4 || resumed.GitInitialCommit != "base-sha" {
			t.Fatalf("resumed checkpoint = level %d, commit %q; want level 4, base-sha", resumed.CurrentLevel, resumed.GitInitialCommit)
		}
		if got, want := recoverableIssueIDs(resumed), []string{"issue-a", "issue-b"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("recoverableIssueIDs = %v, want %v", got, want)
		}
	})

	t.Run("stale level and revision writes are ignored", func(t *testing.T) {
		artifactsDir := t.TempDir()
		baseline := &schemas.DAGState{
			ArtifactsDir:     artifactsDir,
			CurrentLevel:     3,
			GitInitialCommit: "base-sha",
			AllIssues:        []map[string]any{{"id": "issue-a", "status": "running"}},
		}
		saveCheckpoint(baseline, nil)

		staleLevel := *baseline
		staleLevel.CurrentLevel = 2
		saveCheckpoint(&staleLevel, nil)
		assertCheckpointUnchanged(t, artifactsDir, baseline)

		staleRevision := *baseline
		staleRevision.GitInitialCommit = "different-base-sha"
		saveCheckpoint(&staleRevision, nil)
		assertCheckpointUnchanged(t, artifactsDir, baseline)
	})

	t.Run("terminal accepted effect cannot regress or duplicate", func(t *testing.T) {
		previous := &schemas.DAGState{
			AllIssues: []map[string]any{{"id": "issue-a", "status": "accepted"}},
		}
		acceptedAgain := &schemas.DAGState{
			AllIssues: []map[string]any{{"id": "issue-a", "status": "accepted"}},
		}
		if checkpointRegressesTerminalIssue(previous, acceptedAgain) {
			t.Fatal("repeating the same terminal accepted effect was treated as regression")
		}

		regressed := &schemas.DAGState{
			AllIssues: []map[string]any{{"id": "issue-a", "status": "running"}},
		}
		if !checkpointRegressesTerminalIssue(previous, regressed) {
			t.Fatal("terminal accepted effect regression was not fenced")
		}
	})
}

func assertCheckpointUnchanged(t *testing.T, artifactsDir string, want *schemas.DAGState) {
	t.Helper()
	path := filepath.Join(artifactsDir, "execution", "checkpoint.json")
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read checkpoint: %v", err)
	}
	var got schemas.DAGState
	if err := json.Unmarshal(contents, &got); err != nil {
		t.Fatalf("decode checkpoint: %v", err)
	}
	if got.CurrentLevel != want.CurrentLevel || got.GitInitialCommit != want.GitInitialCommit {
		t.Fatalf("checkpoint changed to level %d, commit %q; want level %d, commit %q", got.CurrentLevel, got.GitInitialCommit, want.CurrentLevel, want.GitInitialCommit)
	}
}
