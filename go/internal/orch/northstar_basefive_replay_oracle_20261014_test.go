package orch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const baseFiveGitSHA = "064a7226a500365005b5b662044fc9014a328722"
const baseFiveResumeSourceSHA256 = "02c1a49916f19a1e489ad8a94d6f76aea2ff39869b3bec9351bd96ca33cb7f48"

func baseFiveResumeSource(t *testing.T) string {
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

func TestBaseFiveC4SignatureAndSourceIdentity(t *testing.T) {
	source := baseFiveResumeSource(t)
	digest := sha256.Sum256([]byte(source))
	if actual := hex.EncodeToString(digest[:]); actual != baseFiveResumeSourceSHA256 {
		t.Fatalf("resume.go source SHA changed: got %s want %s", actual, baseFiveResumeSourceSHA256)
	}
	for _, required := range []string{
		"func ResumeBuildHandler(ctx context.Context, deps *Deps, input map[string]any) (any, error)",
		"deps.CallRaw(ctx, \"execute\"",
		"\"resume\":      true",
		"return result, nil",
	} {
		if !strings.Contains(source, required) {
			t.Fatalf("C4 contract marker missing %q (base git %s)", required, baseFiveGitSHA)
		}
	}
}

func TestBaseFiveReplayOraclePositiveAndDependencyReadyOrdering(t *testing.T) {
	source := baseFiveResumeSource(t)
	checkpoint := map[string]any{
		"all_issues":           []any{map[string]any{"id": "issue-a", "status": "pending"}},
		"levels":               []any{map[string]any{"level": 1}},
		"artifacts_dir":        ".artifacts",
		"original_plan_summary": "base-five replay",
	}
	encoded, err := json.Marshal(checkpoint)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"all_issues", "levels", "artifacts_dir", "original_plan_summary"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("positive replay checkpoint lost field %q", key)
		}
	}
	checkpointRead := strings.Index(source, "json.Unmarshal(raw, &checkpoint)")
	dependencyCall := strings.Index(source, "deps.CallRaw(ctx, \"execute\"")
	if checkpointRead < 0 || dependencyCall < 0 || checkpointRead >= dependencyCall {
		t.Fatal("checkpoint reconstruction must precede dependency execution")
	}
	if !strings.Contains(source, "mapStr(checkpoint, \"artifacts_dir\", base)") {
		t.Fatal("replay must restore artifacts_dir before dependency execution")
	}
}

func TestBaseFiveReplayOracleAdversarialRestartFailures(t *testing.T) {
	t.Run("missing checkpoint fails closed", func(t *testing.T) {
		repo := t.TempDir()
		_, err := ResumeBuildHandler(context.Background(), nil, map[string]any{
			"repo_path":     repo,
			"artifacts_dir": ".artifacts",
		})
		if err == nil || !strings.Contains(err.Error(), "No checkpoint found") {
			t.Fatalf("missing checkpoint error = %v", err)
		}
	})

	t.Run("malformed checkpoint fails closed", func(t *testing.T) {
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
	})

	t.Run("canceled restart stops before replay", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := ResumeBuildHandler(ctx, nil, map[string]any{})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled restart error = %v", err)
		}
	})
}

func TestBaseFiveStaleOwnerFencingMarkers(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	checkpointSource, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "..", "dag", "checkpoint.go"))
	if err != nil {
		t.Fatalf("read checkpoint.go: %v", err)
	}
	for _, required := range []string{
		"staleLevel := previous.CurrentLevel > dagState.CurrentLevel",
		"revisionChanged := previous.GitInitialCommit != \"\"",
		"terminalRegression := checkpointRegressesTerminalIssue(previous, dagState)",
		"Checkpoint ignored: stale, revision-mismatched, or terminal-regressive write",
	} {
		if !strings.Contains(string(checkpointSource), required) {
			t.Fatalf("stale-owner fencing marker missing %q", required)
		}
	}
}
