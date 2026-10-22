package orch

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNorthstarNativeResumeRejectsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := ResumeBuildHandler(ctx, nil, map[string]any{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ResumeBuildHandler() error = %v, want context.Canceled", err)
	}
	if result != nil {
		t.Fatalf("ResumeBuildHandler() result = %#v, want nil", result)
	}
}

func TestNorthstarNativeResumeRejectsMissingCheckpoint(t *testing.T) {
	repoPath := t.TempDir()

	result, err := ResumeBuildHandler(context.Background(), nil, map[string]any{
		"repo_path": repoPath,
	})
	if err == nil {
		t.Fatal("ResumeBuildHandler() error = nil, want missing-checkpoint error")
	}
	if result != nil {
		t.Fatalf("ResumeBuildHandler() result = %#v, want nil", result)
	}

	checkpointPath := filepath.Join(repoPath, ".artifacts", "execution", "checkpoint.json")
	if !strings.Contains(err.Error(), checkpointPath) || !strings.Contains(err.Error(), "Cannot resume") {
		t.Fatalf("ResumeBuildHandler() error = %q, want checkpoint path and Cannot resume", err)
	}
}

func TestNorthstarNativeResumeRejectsMalformedCheckpoint(t *testing.T) {
	repoPath := t.TempDir()
	checkpointPath := writeNorthstarCheckpoint(t, repoPath, []byte(`{"all_issues":`))

	result, err := ResumeBuildHandler(context.Background(), nil, map[string]any{
		"repo_path": repoPath,
	})
	if err == nil {
		t.Fatal("ResumeBuildHandler() error = nil, want JSON decode error")
	}
	if result != nil {
		t.Fatalf("ResumeBuildHandler() result = %#v, want nil", result)
	}
	if !strings.Contains(err.Error(), "unexpected end of JSON input") {
		t.Fatalf("ResumeBuildHandler() error = %q, want JSON decode failure", err)
	}
	if _, statErr := os.Stat(checkpointPath); statErr != nil {
		t.Fatalf("malformed checkpoint should remain for diagnosis: %v", statErr)
	}
}

func TestNorthstarNativeResumeRejectsStaleCheckpointRemovedBeforeRestart(t *testing.T) {
	repoPath := t.TempDir()
	checkpointPath := writeNorthstarCheckpoint(t, repoPath, []byte(`{"all_issues":[],"levels":[]}`))
	if err := os.Remove(checkpointPath); err != nil {
		t.Fatalf("remove stale checkpoint: %v", err)
	}

	result, err := ResumeBuildHandler(context.Background(), nil, map[string]any{
		"repo_path": repoPath,
	})
	if err == nil || !strings.Contains(err.Error(), "No checkpoint found") {
		t.Fatalf("ResumeBuildHandler() error = %v, want stale checkpoint rejection", err)
	}
	if result != nil {
		t.Fatalf("ResumeBuildHandler() result = %#v, want nil", result)
	}
}

func TestNorthstarNativeResumeMalformedRetryIsIdempotent(t *testing.T) {
	repoPath := t.TempDir()
	checkpointPath := writeNorthstarCheckpoint(t, repoPath, []byte(`{"all_issues":[{"id":"foreign-owner"}]`))
	before, err := os.ReadFile(checkpointPath)
	if err != nil {
		t.Fatalf("read checkpoint before retry: %v", err)
	}

	var firstError string
	for attempt := 0; attempt < 2; attempt++ {
		result, resumeErr := ResumeBuildHandler(context.Background(), nil, map[string]any{
			"repo_path": repoPath,
		})
		if resumeErr == nil {
			t.Fatalf("attempt %d error = nil, want malformed-checkpoint rejection", attempt+1)
		}
		if result != nil {
			t.Fatalf("attempt %d result = %#v, want nil", attempt+1, result)
		}
		if attempt == 0 {
			firstError = resumeErr.Error()
		} else if resumeErr.Error() != firstError {
			t.Fatalf("retry error = %q, want stable error %q", resumeErr, firstError)
		}
	}

	after, err := os.ReadFile(checkpointPath)
	if err != nil {
		t.Fatalf("read checkpoint after retry: %v", err)
	}
	if string(after) != string(before) {
		t.Fatalf("failed retries mutated durable checkpoint: before %q after %q", before, after)
	}
}

func TestNorthstarNativeResumeSourceContractPreservesDependencyOrderingAndReplayFlag(t *testing.T) {
	source, err := os.ReadFile("resume.go")
	if err != nil {
		t.Fatalf("read resume.go: %v", err)
	}

	requiredFragments := []string{
		`filepath.Join(base, "execution", "checkpoint.json")`,
		`"issues":         any0(checkpoint["all_issues"])`,
		`"levels":         any0(checkpoint["levels"])`,
		`deps.CallRaw(ctx, "execute"`,
		`"resume":      true`,
		`return result, nil`,
	}
	for _, fragment := range requiredFragments {
		if !strings.Contains(string(source), fragment) {
			t.Errorf("resume.go missing required C4 source contract fragment %q", fragment)
		}
	}
}

func writeNorthstarCheckpoint(t *testing.T, repoPath string, content []byte) string {
	t.Helper()

	checkpointPath := filepath.Join(repoPath, ".artifacts", "execution", "checkpoint.json")
	if err := os.MkdirAll(filepath.Dir(checkpointPath), 0o755); err != nil {
		t.Fatalf("create checkpoint directory: %v", err)
	}
	if err := os.WriteFile(checkpointPath, content, 0o600); err != nil {
		t.Fatalf("write checkpoint: %v", err)
	}
	return checkpointPath
}
