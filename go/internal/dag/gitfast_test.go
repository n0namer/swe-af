package dag

// Contract tests for the deterministic git fast-paths — the Go mirror of
// tests/test_git_fast_path.py. Real git repos; only the agent CallFn is
// scripted.

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func gitfastT(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// initFastRepo creates a repo with an "integration" branch checked out.
func initFastRepo(t *testing.T) string {
	t.Helper()
	repo := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	gitfastT(t, repo, "init", "-q")
	gitfastT(t, repo, "checkout", "-q", "-b", "main")
	gitfastT(t, repo, "config", "user.email", "t@example.com")
	gitfastT(t, repo, "config", "user.name", "T")
	if err := os.WriteFile(filepath.Join(repo, "base.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitfastT(t, repo, "add", "base.txt")
	gitfastT(t, repo, "commit", "-q", "-m", "base")
	gitfastT(t, repo, "checkout", "-q", "-b", "integration")
	return repo
}

func branchWithFile(t *testing.T, repo, branch, filename, content string) {
	t.Helper()
	gitfastT(t, repo, "checkout", "-q", "-b", branch, "integration")
	if err := os.WriteFile(filepath.Join(repo, filename), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	gitfastT(t, repo, "add", filename)
	gitfastT(t, repo, "commit", "-q", "-m", "work on "+branch)
	gitfastT(t, repo, "checkout", "-q", "integration")
}

func recordingCallFn(response map[string]any) (func(context.Context, string, map[string]any) (map[string]any, error), *[]string, *[]map[string]any) {
	var mu sync.Mutex
	targets := &[]string{}
	kwargsLog := &[]map[string]any{}
	fn := func(_ context.Context, target string, kwargs map[string]any) (map[string]any, error) {
		mu.Lock()
		defer mu.Unlock()
		*targets = append(*targets, target)
		*kwargsLog = append(*kwargsLog, kwargs)
		return response, nil
	}
	return fn, targets, kwargsLog
}

func TestFastSetupWorktreesNaming(t *testing.T) {
	repo := initFastRepo(t)
	wtDir := filepath.Join(repo, ".worktrees")
	setup, err := fastSetupWorktrees(repo, "integration", []map[string]any{
		{"name": "lexer", "sequence_number": 1},
	}, wtDir, "ab12cd34")
	if err != nil {
		t.Fatalf("fastSetupWorktrees: %v", err)
	}
	ws := setup["workspaces"].([]any)[0].(map[string]any)
	if ws["branch_name"] != "issue/ab12cd34-01-lexer" {
		t.Errorf("branch_name = %v", ws["branch_name"])
	}
	wantPath := filepath.Join(wtDir, "issue-ab12cd34-01-lexer")
	if ws["worktree_path"] != wantPath {
		t.Errorf("worktree_path = %v, want %v", ws["worktree_path"], wantPath)
	}
	// Resume: second call reuses branch+worktree.
	again, err := fastSetupWorktrees(repo, "integration", []map[string]any{
		{"name": "lexer", "sequence_number": 1},
	}, wtDir, "ab12cd34")
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if again["workspaces"].([]any)[0].(map[string]any)["branch_name"] != "issue/ab12cd34-01-lexer" {
		t.Errorf("resume workspaces mismatch")
	}
}

func TestFastMergeBranchesCleanAndConflict(t *testing.T) {
	repo := initFastRepo(t)
	branchWithFile(t, repo, "issue/01-a", "a.txt", "a\n")
	branchWithFile(t, repo, "issue/02-b", "same.txt", "b\n")
	branchWithFile(t, repo, "issue/03-c", "same.txt", "c\n")

	result, err := fastMergeBranches(repo, "integration",
		[]string{"issue/01-a", "issue/02-b", "issue/03-c"}, 1)
	if err != nil {
		t.Fatalf("fastMergeBranches: %v", err)
	}
	if got := asStringSlice(result["merged_branches"]); len(got) != 2 {
		t.Errorf("merged = %v", got)
	}
	if got := asStringSlice(result["failed_branches"]); len(got) != 1 || got[0] != "issue/03-c" {
		t.Errorf("failed = %v", got)
	}
	if asBool(result["needs_integration_test"]) != true {
		t.Error("needs_integration_test should be true for >1 merged branch")
	}
	// Integration branch clean: no in-progress merge.
	if got := gitfastT(t, repo, "status", "--porcelain"); got != "" {
		t.Errorf("dirty after conflict abort: %q", got)
	}
}

func TestFastMergeSingleBranchSkipsIntegrationTest(t *testing.T) {
	repo := initFastRepo(t)
	branchWithFile(t, repo, "issue/01-a", "a.txt", "a\n")
	result, err := fastMergeBranches(repo, "integration", []string{"issue/01-a"}, 0)
	if err != nil {
		t.Fatalf("fastMergeBranches: %v", err)
	}
	if asBool(result["needs_integration_test"]) {
		t.Error("single-branch merge must not request integration tests")
	}
}

func TestFastCleanupWorktrees(t *testing.T) {
	repo := initFastRepo(t)
	wtDir := filepath.Join(repo, ".worktrees")
	if _, err := fastSetupWorktrees(repo, "integration",
		[]map[string]any{{"name": "a", "sequence_number": 1}}, wtDir, "ab12cd34"); err != nil {
		t.Fatal(err)
	}
	worktree := filepath.Join(wtDir, "issue-ab12cd34-01-a")
	if err := os.WriteFile(filepath.Join(worktree, "change.txt"), []byte("merged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitfastT(t, worktree, "add", "change.txt")
	gitfastT(t, worktree, "commit", "-q", "-m", "issue work")
	gitfastT(t, repo, "merge", "--no-ff", "-q", "issue/ab12cd34-01-a", "-m", "integrate issue")
	result, err := fastCleanupWorktrees(repo, wtDir, []string{"issue/ab12cd34-01-a"}, "ab12cd34")
	if err != nil {
		t.Fatalf("fastCleanupWorktrees: %v", err)
	}
	if got := asStringSlice(result["cleaned"]); len(got) != 1 {
		t.Errorf("cleaned = %v", got)
	}
	if got := gitfastT(t, repo, "branch", "--list", "issue/ab12cd34-01-a"); got != "" {
		t.Errorf("branch survived cleanup: %q", got)
	}
	// Non-repo path errors so the caller can fall back to the agent.
	if _, err := fastCleanupWorktrees(t.TempDir(), wtDir, []string{"issue/ab12cd34-01-x"}, "ab12cd34"); err == nil {
		t.Error("expected error for non-repo path")
	}
}

func TestFastCleanupRefusesUnsafeWorktrees(t *testing.T) {
	tests := []struct {
		name   string
		branch string
		setup  func(*testing.T, string, string)
	}{
		{
			name:   "dirty worktree",
			branch: "issue/ab12cd34-01-a",
			setup: func(t *testing.T, repo, wtDir string) {
				t.Helper()
				if _, err := fastSetupWorktrees(repo, "integration", []map[string]any{{"name": "a", "sequence_number": 1}}, wtDir, "ab12cd34"); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(wtDir, "issue-ab12cd34-01-a", "untracked.txt"), []byte("keep"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:   "unmerged branch",
			branch: "issue/ab12cd34-01-a",
			setup: func(t *testing.T, repo, wtDir string) {
				t.Helper()
				if _, err := fastSetupWorktrees(repo, "integration", []map[string]any{{"name": "a", "sequence_number": 1}}, wtDir, "ab12cd34"); err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(wtDir, "issue-ab12cd34-01-a")
				if err := os.WriteFile(filepath.Join(path, "change.txt"), []byte("change\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				gitfastT(t, path, "add", "change.txt")
				gitfastT(t, path, "commit", "-q", "-m", "unmerged")
			},
		},
		{
			name:   "foreign build",
			branch: "issue/cdef1234-01-a",
			setup: func(t *testing.T, repo, wtDir string) {
				t.Helper()
				if _, err := fastSetupWorktrees(repo, "integration", []map[string]any{{"name": "a", "sequence_number": 1}}, wtDir, "cdef1234"); err != nil {
					t.Fatal(err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := initFastRepo(t)
			wtDir := filepath.Join(repo, ".worktrees")
			tt.setup(t, repo, wtDir)
			_, err := fastCleanupWorktrees(repo, wtDir, []string{tt.branch}, "ab12cd34")
			var unsafe *unsafeCleanupError
			if !errors.As(err, &unsafe) {
				t.Fatalf("expected unsafe cleanup refusal, got %v", err)
			}
			path := filepath.Join(wtDir, strings.ReplaceAll(tt.branch, "/", "-"))
			if _, statErr := os.Stat(path); statErr != nil {
				t.Errorf("worktree was removed despite refusal: %v", statErr)
			}
			if got := gitfastT(t, repo, "branch", "--list", tt.branch); got == "" {
				t.Error("branch was removed despite refusal")
			}
		})
	}
}

func TestCleanupSingleRepoUnsafeOwnershipDoesNotFallback(t *testing.T) {
	repo := initFastRepo(t)
	callFn, targets, _ := recordingCallFn(map[string]any{"success": true})
	err := cleanupSingleRepo(context.Background(), callFn, "node", repo, filepath.Join(repo, ".worktrees"),
		[]string{"issue/01-a"}, "", "", 0, "", "", true, nil)
	var unsafe *unsafeCleanupError
	if !errors.As(err, &unsafe) {
		t.Fatalf("expected unsafe cleanup refusal, got %v", err)
	}
	if len(*targets) != 0 {
		t.Fatalf("unsafe cleanup was delegated to agent: %v", *targets)
	}
}

func TestCleanupSingleRepoUsesTrustedBuildID(t *testing.T) {
	repo := initFastRepo(t)
	wtDir := filepath.Join(repo, ".worktrees")
	branch := "issue/cdef1234-01-a"
	if _, err := fastSetupWorktrees(repo, "integration",
		[]map[string]any{{"name": "a", "sequence_number": 1}}, wtDir, "cdef1234"); err != nil {
		t.Fatal(err)
	}
	worktree := filepath.Join(wtDir, "issue-cdef1234-01-a")
	if err := os.WriteFile(filepath.Join(worktree, "change.txt"), []byte("merged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitfastT(t, worktree, "add", "change.txt")
	gitfastT(t, worktree, "commit", "-q", "-m", "issue work")
	gitfastT(t, repo, "merge", "--no-ff", "-q", branch, "-m", "integrate issue")

	callFn, targets, _ := recordingCallFn(map[string]any{"success": true})
	err := cleanupSingleRepo(context.Background(), callFn, "node", repo, wtDir,
		[]string{branch}, "ab12cd34", "", 0, "", "", true, nil)
	var unsafe *unsafeCleanupError
	if !errors.As(err, &unsafe) {
		t.Fatalf("expected ownership refusal, got %v", err)
	}
	if len(*targets) != 0 {
		t.Fatalf("foreign-build cleanup was delegated to agent: %v", *targets)
	}
	if _, statErr := os.Stat(worktree); statErr != nil {
		t.Errorf("foreign worktree was removed: %v", statErr)
	}
	if got := gitfastT(t, repo, "branch", "--list", branch); got == "" {
		t.Error("foreign branch was removed")
	}
}

func TestDispatchMergeHandsOnlyConflictsToAgent(t *testing.T) {
	repo := initFastRepo(t)
	branchWithFile(t, repo, "issue/01-a", "same.txt", "a\n")
	branchWithFile(t, repo, "issue/02-b", "same.txt", "b\n")

	callFn, targets, kwargsLog := recordingCallFn(map[string]any{
		"success": true, "merged_branches": []any{"issue/02-b"},
		"failed_branches": []any{}, "summary": "agent resolved",
	})
	cfg := testCfg(t, nil)
	branches := []map[string]any{
		{"branch_name": "issue/01-a"}, {"branch_name": "issue/02-b"},
	}
	result, err := dispatchMerge(context.Background(), callFn, "node", cfg,
		repo, "integration", branches,
		map[string]any{"branches_to_merge": branches}, 0, nil)
	if err != nil {
		t.Fatalf("dispatchMerge: %v", err)
	}
	if len(*targets) != 1 || (*targets)[0] != "node.run_merger" {
		t.Fatalf("agent calls = %v", *targets)
	}
	sent := (*kwargsLog)[0]["branches_to_merge"].([]map[string]any)
	if len(sent) != 1 || sent[0]["branch_name"] != "issue/02-b" {
		t.Errorf("agent got %v, want only the conflicted branch", sent)
	}
	merged := asStringSlice(result["merged_branches"])
	if len(merged) != 2 {
		t.Errorf("combined merged = %v", merged)
	}
}

func TestDispatchWorkspaceSetupNoAgentOnRealRepo(t *testing.T) {
	repo := initFastRepo(t)
	callFn, targets, _ := recordingCallFn(map[string]any{"success": true})
	cfg := testCfg(t, nil)
	setup, err := dispatchWorkspaceSetup(context.Background(), callFn, "node", cfg,
		repo, "integration",
		[]map[string]any{{"name": "a", "sequence_number": 1}},
		filepath.Join(repo, ".worktrees"), "", 0, "", nil)
	if err != nil {
		t.Fatalf("dispatchWorkspaceSetup: %v", err)
	}
	if !asBool(setup["success"]) || len(*targets) != 0 {
		t.Errorf("expected fast-path success with no agent calls; calls=%v", *targets)
	}
	// Fake repo → agent fallback.
	callFn2, targets2, _ := recordingCallFn(map[string]any{"success": true, "workspaces": []any{}})
	if _, err := dispatchWorkspaceSetup(context.Background(), callFn2, "node", cfg,
		filepath.Join(t.TempDir(), "nope"), "integration",
		[]map[string]any{{"name": "a"}}, filepath.Join(t.TempDir(), "wt"), "", 0, "", nil); err != nil {
		t.Fatalf("fallback: %v", err)
	}
	if len(*targets2) != 1 || (*targets2)[0] != "node.run_workspace_setup" {
		t.Errorf("fallback calls = %v", *targets2)
	}
}
