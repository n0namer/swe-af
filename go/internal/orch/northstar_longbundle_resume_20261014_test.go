package orch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const expectedResumeSourceSHA256 = "02c1a49916f19a1e489ad8a94d6f76aea2ff39869b3bec9351bd96ca33cb7f48"

func northstarResumeSource(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	source, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "resume.go"))
	if err != nil {
		t.Fatalf("read resume.go: %v", err)
	}
	return string(source)
}

func TestNorthstarLongBundleResumeSignatureAndImmutableSource(t *testing.T) {
	source := northstarResumeSource(t)
	digest := sha256.Sum256([]byte(source))
	if actual := hex.EncodeToString(digest[:]); actual != expectedResumeSourceSHA256 {
		t.Fatalf("resume.go source SHA changed: got %s want %s", actual, expectedResumeSourceSHA256)
	}
	for _, required := range []string{
		"func ResumeBuildHandler(ctx context.Context, deps *Deps, input map[string]any) (any, error)",
		"deps.CallRaw(ctx, \"execute\"",
		"\"resume\":      true",
		"return result, nil",
	} {
		if !strings.Contains(source, required) {
			t.Fatalf("resume.go missing required contract marker %q", required)
		}
	}
}

func TestNorthstarLongBundleMissingCheckpointFailsClosed(t *testing.T) {
	repo := t.TempDir()
	_, err := ResumeBuildHandler(context.Background(), nil, map[string]any{
		"repo_path":      repo,
		"artifacts_dir":  ".artifacts",
		"config":         map[string]any{},
		"git_config":     map[string]any{},
	})
	if err == nil || !strings.Contains(err.Error(), "No checkpoint found") {
		t.Fatalf("missing checkpoint error = %v", err)
	}
}

func TestNorthstarLongBundleMalformedCheckpointFailsClosed(t *testing.T) {
	repo := t.TempDir()
	checkpoint := filepath.Join(repo, ".artifacts", "execution", "checkpoint.json")
	if err := os.MkdirAll(filepath.Dir(checkpoint), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(checkpoint, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ResumeBuildHandler(context.Background(), nil, map[string]any{
		"repo_path":     repo,
		"artifacts_dir": ".artifacts",
	})
	if err == nil {
		t.Fatal("malformed checkpoint unexpectedly resumed")
	}
}

func TestNorthstarLongBundleCrashRestartAndDependencyOrdering(t *testing.T) {
	source := northstarResumeSource(t)
	checkpointRead := strings.Index(source, "json.Unmarshal(raw, &checkpoint)")
	dependencyCall := strings.Index(source, "deps.CallRaw(ctx, \"execute\"")
	if checkpointRead < 0 || dependencyCall < 0 || checkpointRead >= dependencyCall {
		t.Fatalf("checkpoint reconstruction must precede execute dependency call")
	}
	if strings.Count(source, "if err := ctx.Err(); err != nil") < 2 {
		t.Fatal("resume path must check cancellation before and after checkpoint reconstruction")
	}
	for _, required := range []string{
		"checkpoint[\"all_issues\"]",
		"checkpoint[\"levels\"]",
		"mapStr(checkpoint, \"artifacts_dir\", base)",
		"mapStr(checkpoint, \"original_plan_summary\", \"\")",
	} {
		if !strings.Contains(source, required) {
			t.Fatalf("restart reconstruction missing %q", required)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := ResumeBuildHandler(ctx, nil, map[string]any{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled restart error = %v", err)
	}
}

func TestNorthstarLongBundleCheckpointRejectsDuplicateOrRegressiveEdits(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	checkpointSource, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "..", "dag", "checkpoint.go"))
	if err != nil {
		t.Fatalf("read checkpoint.go: %v", err)
	}
	source := string(checkpointSource)
	for _, required := range []string{
		"staleLevel := previous.CurrentLevel > dagState.CurrentLevel",
		"revisionChanged := previous.GitInitialCommit != \"\"",
		"terminalRegression := checkpointRegressesTerminalIssue(previous, dagState)",
		"Checkpoint ignored: stale, revision-mismatched, or terminal-regressive write",
	} {
		if !strings.Contains(source, required) {
			t.Fatalf("checkpoint safety guard missing %q", required)
		}
	}
}
