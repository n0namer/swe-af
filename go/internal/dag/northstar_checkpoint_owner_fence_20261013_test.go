package dag

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
)

func TestNorthstarCheckpointOwnerFenceNativePositivePath(t *testing.T) {
	artifactsDir := t.TempDir()
	initial := checkpointFenceState(artifactsDir, "revision-a", 1, "running")

	saveCheckpoint(initial, nil)
	loaded := loadCheckpoint(artifactsDir)
	assertCheckpointFenceState(t, loaded, "revision-a", 1, "running")

	terminal := checkpointFenceState(artifactsDir, "revision-a", 2, "completed")
	saveCheckpoint(terminal, nil)
	loaded = loadCheckpoint(artifactsDir)
	assertCheckpointFenceState(t, loaded, "revision-a", 2, "completed")
}

func TestNorthstarCheckpointOwnerFenceRejectsConcurrentStaleRestarts(t *testing.T) {
	artifactsDir := t.TempDir()
	terminal := checkpointFenceState(artifactsDir, "revision-a", 4, "accepted")
	saveCheckpoint(terminal, nil)

	candidates := []*schemas.DAGState{
		checkpointFenceState(artifactsDir, "revision-a", 3, "running"),
		checkpointFenceState(artifactsDir, "revision-a", 5, "running"),
		checkpointFenceState(artifactsDir, "revision-b", 5, "accepted"),
	}

	var waitGroup sync.WaitGroup
	var notesMu sync.Mutex
	notes := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		candidate := candidate
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			saveCheckpoint(candidate, func(message string, _ []string) {
				notesMu.Lock()
				notes = append(notes, message)
				notesMu.Unlock()
			})
		}()
	}
	waitGroup.Wait()

	loaded := loadCheckpoint(artifactsDir)
	assertCheckpointFenceState(t, loaded, "revision-a", 4, "accepted")
	if len(notes) != len(candidates) {
		t.Fatalf("expected %d rejected-write notes, got %d: %v", len(candidates), len(notes), notes)
	}
	for _, note := range notes {
		if !strings.Contains(note, "Checkpoint ignored") {
			t.Fatalf("expected rejected-write note, got %q", note)
		}
	}
}

func TestNorthstarCheckpointOwnerFenceFailsClosedAndRecoversOnRetry(t *testing.T) {
	artifactsDir := t.TempDir()
	if loaded := loadCheckpoint(artifactsDir); loaded != nil {
		t.Fatalf("missing checkpoint must fail closed, got %#v", loaded)
	}

	path := filepath.Join(artifactsDir, "execution", "checkpoint.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create checkpoint directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("{corrupt"), 0o644); err != nil {
		t.Fatalf("write corrupt checkpoint: %v", err)
	}
	if loaded := loadCheckpoint(artifactsDir); loaded != nil {
		t.Fatalf("corrupt checkpoint must fail closed, got %#v", loaded)
	}

	retry := checkpointFenceState(artifactsDir, "revision-retry", 2, "failed")
	saveCheckpoint(retry, nil)
	loaded := loadCheckpoint(artifactsDir)
	assertCheckpointFenceState(t, loaded, "revision-retry", 2, "failed")
}

func checkpointFenceState(artifactsDir, revision string, level int, status string) *schemas.DAGState {
	return &schemas.DAGState{
		ArtifactsDir:     artifactsDir,
		CurrentLevel:     level,
		GitInitialCommit: revision,
		AllIssues: []map[string]any{
			{"id": "C4-owner-fence", "status": status},
		},
	}
}

func assertCheckpointFenceState(t *testing.T, state *schemas.DAGState, revision string, level int, status string) {
	t.Helper()
	if state == nil {
		t.Fatal("expected checkpoint state, got nil")
	}
	if state.GitInitialCommit != revision {
		t.Fatalf("expected revision %q, got %q", revision, state.GitInitialCommit)
	}
	if state.CurrentLevel != level {
		t.Fatalf("expected level %d, got %d", level, state.CurrentLevel)
	}
	if len(state.AllIssues) != 1 {
		t.Fatalf("expected one issue, got %d", len(state.AllIssues))
	}
	if got := mapGetStr(state.AllIssues[0], "id", ""); got != "C4-owner-fence" {
		t.Fatalf("expected issue ID %q, got %q", "C4-owner-fence", got)
	}
	if got := mapGetStr(state.AllIssues[0], "status", ""); got != status {
		t.Fatalf("expected terminal status %q, got %q", status, got)
	}
}
