package dag

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
)

func r24ClaimPlan(root string) map[string]any {
	return map[string]any{
		"artifacts_dir": root,
		"issues": []map[string]any{
			{"id": "r24-claim", "name": "r24-claim", "status": "pending"},
		},
		"levels": [][]string{{"r24-claim"}},
	}
}

func r24ClaimRun(root, revision string, resume bool, effect ExecuteFn) (*schemas.DAGState, error) {
	return RunDAG(context.Background(), r24ClaimPlan(root), root, nil, "", nil,
		WithResume(resume), WithBuildID("r24-build"),
		WithGitConfig(map[string]any{"initial_commit_sha": revision}), WithExecuteFn(effect))
}

func TestR24ClaimCheckpointCompletedReplay(t *testing.T) {
	root := t.TempDir()
	var effects int32
	effect := func(context.Context, map[string]any, *schemas.DAGState) (map[string]any, error) {
		atomic.AddInt32(&effects, 1)
		return map[string]any{"outcome": "completed"}, nil
	}
	state, err := r24ClaimRun(root, "source-a", false, effect)
	if err != nil || state == nil || len(state.CompletedIssues) != 1 {
		t.Fatalf("initial execution: state=%+v err=%v", state, err)
	}
	state, err = r24ClaimRun(root, "source-a", true, effect)
	if err != nil || state == nil || len(state.CompletedIssues) != 1 {
		t.Fatalf("completed replay: state=%+v err=%v", state, err)
	}
	if got := atomic.LoadInt32(&effects); got != 1 {
		t.Fatalf("completed replay repeated effect: count=%d", got)
	}
}

func TestR24ClaimCheckpointRejectsInvalidResumeSource(t *testing.T) {
	for _, invalid := range []string{"foreign-revision", "corrupt-checkpoint"} {
		t.Run(invalid, func(t *testing.T) {
			root := t.TempDir()
			initial := initDAGState(r24ClaimPlan(root), root,
				map[string]any{"initial_commit_sha": "source-a"}, "r24-build")
			saveCheckpoint(initial, nil)
			if loaded := loadCheckpoint(root); loaded == nil || loaded.GitInitialCommit != "source-a" {
				t.Fatal("fixture checkpoint was not persisted")
			}
			revision := "source-b"
			if invalid == "corrupt-checkpoint" {
				revision = "source-a"
				if err := os.WriteFile(checkpointPath(initial), []byte("{corrupt"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			var effects int32
			_, err := r24ClaimRun(root, revision, true,
				func(context.Context, map[string]any, *schemas.DAGState) (map[string]any, error) {
					atomic.AddInt32(&effects, 1)
					return map[string]any{"outcome": "completed"}, nil
				})
			if got := atomic.LoadInt32(&effects); got != 0 {
				t.Errorf("invalid resume executed an external effect: count=%d", got)
			}
			if err == nil {
				t.Error("invalid resume must return an error rather than silently restart or trust foreign state")
			}
		})
	}
}

func TestR24ClaimCheckpointCrashRestart(t *testing.T) {
	root := t.TempDir()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	child := exec.CommandContext(ctx, binary, "-test.run=^TestR24ClaimCheckpointCrashChild$", "-test.count=1")
	child.Env = append(os.Environ(), "SWE_AF_R24_CLAIM_CRASH_ROOT="+root)
	if err := child.Start(); err != nil {
		t.Fatalf("start native crash child: %v", err)
	}
	waited := false
	defer func() {
		if !waited {
			_ = child.Process.Kill()
			_ = child.Wait()
		}
	}()
	ready := filepath.Join(root, "effect-ready")
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		} else if !os.IsNotExist(err) {
			t.Fatal(err)
		}
		if time.Now().After(deadline) {
			t.Fatal("child did not reach the durable effect / uncommitted result crash boundary")
		}
		time.Sleep(10 * time.Millisecond)
	}
	before := loadCheckpoint(root)
	if before == nil || len(before.InFlightIssues) != 1 || len(before.CompletedIssues) != 0 {
		t.Fatalf("expected durable in-flight checkpoint before kill: %+v", before)
	}
	if err := child.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	waitErr := child.Wait()
	waited = true
	if waitErr == nil {
		t.Fatal("crash child unexpectedly exited successfully")
	}
	if ctx.Err() != nil {
		t.Fatalf("child timed out instead of reaching controlled kill: %v", ctx.Err())
	}
	_, resumeErr := r24ClaimRun(root, "source-a", true,
		func(context.Context, map[string]any, *schemas.DAGState) (map[string]any, error) {
			if err := r24ClaimAppendEffect(root); err != nil {
				return nil, err
			}
			return map[string]any{"outcome": "completed"}, nil
		})
	effects, err := os.ReadFile(filepath.Join(root, "effects.log"))
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(effects), "effect\n"); got != 1 {
		t.Fatalf("crash restart duplicated external effect: count=%d resumeErr=%v", got, resumeErr)
	}
	after := loadCheckpoint(root)
	if after == nil || after.CurrentLevel < before.CurrentLevel || len(after.CompletedIssues) < len(before.CompletedIssues) {
		t.Fatalf("checkpoint regressed after restart: before=%+v after=%+v", before, after)
	}
	if resumeErr == nil && len(after.CompletedIssues) != 1 {
		t.Fatal("successful recovery did not durably complete the issue")
	}
}

func TestR24ClaimCheckpointCrashChild(t *testing.T) {
	root := os.Getenv("SWE_AF_R24_CLAIM_CRASH_ROOT")
	if root == "" {
		return
	}
	_, err := r24ClaimRun(root, "source-a", false,
		func(ctx context.Context, _ map[string]any, _ *schemas.DAGState) (map[string]any, error) {
			if err := r24ClaimAppendEffect(root); err != nil {
				return nil, err
			}
			if err := os.WriteFile(filepath.Join(root, "effect-ready"), []byte("ready"), 0o644); err != nil {
				return nil, err
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Minute):
				return nil, fmt.Errorf("crash child was not killed at the effect boundary")
			}
		})
	t.Fatalf("crash child returned before kill: %v", err)
}

func r24ClaimAppendEffect(root string) error {
	file, err := os.OpenFile(filepath.Join(root, "effects.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprint(file, "effect\n"); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
