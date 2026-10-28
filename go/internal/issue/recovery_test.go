package issue

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResumeRejectsWorktreeIdentityMismatch(t *testing.T) {
	repo := initRepo(t)
	base := gitT(t, repo, "rev-parse", "HEAD")
	wt := filepath.Join(repo, ".worktrees", "deadbeef-add-retry-helper")
	if err := addWorktree(repo, wt, "issue/deadbeef-wrong-issue", base); err != nil {
		t.Fatal(err)
	}
	_, err := ImplementIssue(context.Background(), &Deps{Call: func(context.Context, string, map[string]any) (map[string]any, error) {
		t.Fatal("unexpected coder call")
		return nil, nil
	}}, map[string]any{
		"repo_path": repo, "resume_build_id": "deadbeef", "issue": map[string]any{"title": "Add retry helper", "description": "d"},
	})
	if err == nil || !strings.Contains(err.Error(), "branch mismatch") {
		t.Fatalf("error = %v", err)
	}
}

func TestDeclaredScopeViolationFailsClosed(t *testing.T) {
	repo := initRepo(t)
	call := func(ctx context.Context, target string, kwargs map[string]any) (map[string]any, error) {
		switch target[strings.LastIndex(target, ".")+1:] {
		case "run_coder":
			wt := kwargs["worktree_path"].(string)
			for _, name := range []string{"allowed.go", "extra.go"} {
				if err := os.WriteFile(filepath.Join(wt, name), []byte("package x\n"), 0644); err != nil {
					return nil, err
				}
			}
			gitT(t, wt, "add", "allowed.go", "extra.go")
			gitT(t, wt, "commit", "-q", "-m", "change")
			return map[string]any{"files_changed": []any{"allowed.go", "extra.go"}, "complete": true, "summary": "done"}, nil
		case "run_code_reviewer":
			return map[string]any{"approved": true}, nil
		default:
			return nil, fmt.Errorf("unexpected call %s", target)
		}
	}
	r := runImplement(t, repo, call, map[string]any{"issue": map[string]any{"title": "Scoped", "description": "d", "files_to_modify": []any{"allowed.go"}}, "config": map[string]any{"verify": true, "keep_worktree": true}})
	if r["success"] != false || r["outcome"] != "failed_unrecoverable" || !strings.Contains(r["error_message"].(string), "extra.go") {
		t.Fatalf("result = %#v", r)
	}
}

func TestCancellationWithNoObservedProgressDoesNotReportSuccess(t *testing.T) {
	repo := initRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	call := func(context.Context, string, map[string]any) (map[string]any, error) {
		cancel()
		return nil, context.Canceled
	}
	raw, err := ImplementIssue(ctx, &Deps{Call: call, NodeID: "test-node"}, map[string]any{"repo_path": repo, "issue": map[string]any{"title": "Canceled", "description": "d"}, "config": map[string]any{"verify": false}})
	if err == nil {
		t.Fatalf("expected cancellation without observed progress, got %#v", raw)
	}
}

func TestCancellationPreservesProgressAndResumeContinuesSameWorktree(t *testing.T) {
	repo := initRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	call := func(_ context.Context, target string, kwargs map[string]any) (map[string]any, error) {
		if !strings.HasSuffix(target, ".run_coder") {
			return nil, fmt.Errorf("unexpected call before interruption: %s", target)
		}
		wt := kwargs["worktree_path"].(string)
		if err := os.WriteFile(filepath.Join(wt, "partial.go"), []byte("package partial\n"), 0o644); err != nil {
			return nil, err
		}
		gitT(t, wt, "add", "partial.go")
		gitT(t, wt, "commit", "-q", "-m", "checkpoint partial implementation")
		cancel()
		return nil, context.Canceled
	}
	raw, err := ImplementIssue(ctx, &Deps{Call: call, NodeID: "test-node"}, map[string]any{
		"repo_path": repo, "issue": map[string]any{"title": "Recover issue", "description": "resume without losing progress"},
		"config": map[string]any{"verify": false, "keep_worktree": true},
	})
	if err != nil {
		t.Fatalf("interrupted ImplementIssue: %v", err)
	}
	interrupted := raw.(map[string]any)
	if interrupted["success"] != false || interrupted["outcome"] != "interrupted_with_progress" {
		t.Fatalf("interrupted result = %#v", interrupted)
	}
	buildID, _ := interrupted["build_id"].(string)
	wt, _ := interrupted["worktree_path"].(string)
	branch, _ := interrupted["branch"].(string)
	if len(buildID) != 8 || wt == "" || branch == "" {
		t.Fatalf("missing resumable identity: %#v", interrupted)
	}
	if _, err := os.Stat(filepath.Join(wt, "partial.go")); err != nil {
		t.Fatalf("partial worktree not preserved: %v", err)
	}

	resume := func(_ context.Context, target string, kwargs map[string]any) (map[string]any, error) {
		name := target[strings.LastIndex(target, ".")+1:]
		if name == "run_coder" {
			resumedWT := kwargs["worktree_path"].(string)
			if resumedWT != wt {
				return nil, fmt.Errorf("worktree = %q, want %q", resumedWT, wt)
			}
			if _, err := os.Stat(filepath.Join(resumedWT, "partial.go")); err != nil {
				return nil, fmt.Errorf("partial work lost: %w", err)
			}
			if err := os.WriteFile(filepath.Join(resumedWT, "final.go"), []byte("package final\n"), 0o644); err != nil {
				return nil, err
			}
			gitT(t, resumedWT, "add", "final.go")
			gitT(t, resumedWT, "commit", "-q", "-m", "finish implementation")
			return map[string]any{"files_changed": []any{"final.go"}, "complete": true, "summary": "finished"}, nil
		}
		if name == "run_code_reviewer" {
			return map[string]any{"approved": true, "blocking": false, "summary": "approved"}, nil
		}
		return nil, fmt.Errorf("unexpected resume call: %s", target)
	}
	completed := runImplement(t, repo, resume, map[string]any{
		"resume_build_id": buildID,
		"issue":           map[string]any{"title": "Recover issue", "description": "resume without losing progress"},
		"config":          map[string]any{"verify": false, "keep_worktree": true},
	})
	if completed["success"] != true || completed["build_id"] != buildID || completed["branch"] != branch {
		t.Fatalf("resumed result = %#v", completed)
	}
	if got := gitT(t, wt, "show", "HEAD:partial.go"); got != "package partial" {
		t.Fatalf("prior commit lost: %q", got)
	}
	if got := gitT(t, wt, "show", "HEAD:final.go"); got != "package final" {
		t.Fatalf("resumed change missing: %q", got)
	}
}

func TestVerifierGetsGitChangedFilesInsteadOfCoderClaim(t *testing.T) {
	repo := initRepo(t)
	call := func(_ context.Context, target string, kwargs map[string]any) (map[string]any, error) {
		switch target[strings.LastIndex(target, ".")+1:] {
		case "run_coder":
			wt := kwargs["worktree_path"].(string)
			if err := os.WriteFile(filepath.Join(wt, "actual.go"), []byte("package actual\n"), 0o644); err != nil {
				return nil, err
			}
			gitT(t, wt, "add", "actual.go")
			gitT(t, wt, "commit", "-q", "-m", "actual source change")
			return map[string]any{"files_changed": []any{"claimed.go"}, "complete": true, "summary": "done"}, nil
		case "run_code_reviewer":
			return map[string]any{"approved": true, "blocking": false}, nil
		case "run_verifier":
			issues := kwargs["completed_issues"].([]any)
			files := issues[0].(map[string]any)["files_changed"].([]any)
			if len(files) != 1 || files[0] != "actual.go" {
				return nil, fmt.Errorf("verifier got files_changed %#v", files)
			}
			return map[string]any{"passed": true, "summary": "verified actual change"}, nil
		default:
			return nil, fmt.Errorf("unexpected call: %s", target)
		}
	}
	r := runImplement(t, repo, call, map[string]any{"config": map[string]any{"verify": true, "keep_worktree": true}})
	if r["success"] != true {
		t.Fatalf("result = %#v", r)
	}
}

func TestReviewerSeesCleanCheckpointWithCoderBytes(t *testing.T) {
	repo := initRepo(t)
	const contents = "package reviewed\n\nconst Value = 42\n"
	call := func(_ context.Context, target string, kwargs map[string]any) (map[string]any, error) {
		switch target[strings.LastIndex(target, ".")+1:] {
		case "run_coder":
			wt := kwargs["worktree_path"].(string)
			if err := os.WriteFile(filepath.Join(wt, "reviewed.go"), []byte(contents), 0o644); err != nil {
				return nil, err
			}
			return map[string]any{"files_changed": []any{"reviewed.go"}, "complete": true, "summary": "implemented"}, nil
		case "run_code_reviewer":
			worktree := kwargs["worktree_path"].(string)
			if status := gitT(t, worktree, "status", "--porcelain"); status != "" {
				return nil, fmt.Errorf("reviewer saw dirty worktree: %s", status)
			}
			got, err := os.ReadFile(filepath.Join(worktree, "reviewed.go"))
			if err != nil {
				return nil, err
			}
			if string(got) != contents {
				return nil, fmt.Errorf("reviewer bytes = %q", got)
			}
			return map[string]any{"approved": true, "blocking": false, "summary": "approved committed bytes"}, nil
		default:
			return nil, fmt.Errorf("unexpected call: %s", target)
		}
	}
	r := runImplement(t, repo, call, map[string]any{"config": map[string]any{"verify": false, "keep_worktree": true}})
	if r["success"] != true {
		t.Fatalf("result = %#v", r)
	}
}
