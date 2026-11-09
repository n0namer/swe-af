package issue

// build.go ports swe_af/issue/build.py — the implement_issue orchestrator.

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Agent-Field/SWE-AF/go/internal/afx"
	"github.com/Agent-Field/SWE-AF/go/internal/coding"
	"github.com/Agent-Field/SWE-AF/go/internal/config"
	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
)

// CallFn dispatches to a reasoner by target with the same keyword args Python
// passes to app.call. Structurally identical to coding.CallFn / fast.CallFn, so
// the node's single agent.Call + envelope-unwrap closure satisfies all three.
type CallFn func(ctx context.Context, target string, kwargs map[string]any) (map[string]any, error)

// Noter matches the SDK's *agent.Agent.Note method; tests supply a recorder.
type Noter interface {
	Note(ctx context.Context, message string, tags ...string)
}

// Handler is the registration signature, matching the roles/fast packages.
type Handler func(ctx context.Context, deps *Deps, input map[string]any) (any, error)

// Handlers returns the name→handler map, keyed by the EXACT Python reasoner
// name. The wiring registers it on BOTH nodes (Python includes issue_router in
// swe_af/app.py and swe_af/fast/app.py).
func Handlers() map[string]Handler {
	return map[string]Handler{"implement_issue": ImplementIssue}
}

// Deps carries the seams implement_issue needs.
type Deps struct {
	Call   CallFn
	Note   Noter
	NodeID string
}

func (d *Deps) note(ctx context.Context, msg string, tags ...string) {
	if d != nil && d.Note != nil {
		d.Note.Note(ctx, msg, tags...)
	}
}

type implementInput struct {
	Issue             map[string]any `json:"issue"`
	RepoPath          string         `json:"repo_path"`
	BaseBranch        string         `json:"base_branch"`
	ResumeBuildID     string         `json:"resume_build_id"`
	ArtifactsDir      string         `json:"artifacts_dir"`
	AdditionalContext string         `json:"additional_context"`
	Config            map[string]any `json:"config"`
}

var completedOutcomes = map[string]bool{
	"completed":           true,
	"completed_with_debt": true,
}

func declaredScopeViolations(spec *Spec, files []string) []string {
	if spec == nil {
		return nil
	}
	allowed := map[string]struct{}{}
	add := func(paths []string) {
		for _, raw := range paths {
			path := filepath.ToSlash(filepath.Clean(strings.TrimSpace(raw)))
			if path == "" || path == "." {
				continue
			}
			allowed[path] = struct{}{}
		}
	}
	add(spec.FilesToCreate)
	add(spec.FilesToModify)
	if len(allowed) == 0 {
		return nil
	}
	var violations []string
	for _, raw := range files {
		path := filepath.ToSlash(filepath.Clean(strings.TrimSpace(raw)))
		if path == "" || path == "." {
			continue
		}
		if _, ok := allowed[path]; !ok {
			violations = append(violations, path)
		}
	}
	return violations
}

func newBuildID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%08x", time.Now().UnixNano()&0xffffffff)
	}
	return hex.EncodeToString(b)
}

// ImplementIssue implements one fully-scoped issue on an isolated branch (the
// sub-harness entry point). Ports build.implement_issue/_implement_issue_impl:
// setup/validation errors return an error (they are the caller's to fix);
// failures after branch creation come back as a structured result.
func ImplementIssue(ctx context.Context, deps *Deps, input map[string]any) (any, error) {
	in, err := afx.Bind[implementInput](input)
	if err != nil {
		return nil, err
	}
	if in.RepoPath == "" {
		return nil, fmt.Errorf("implement_issue: repo_path is required")
	}
	if in.ArtifactsDir == "" {
		in.ArtifactsDir = ".artifacts"
	}

	cfg, err := config.LoadIssueBuildConfig(in.Config)
	if err != nil {
		return nil, err
	}
	spec, err := ParseSpec(in.Issue)
	if err != nil {
		return nil, err
	}
	planned := spec.ToPlannedIssue(in.AdditionalContext)
	plannedName, _ := planned["name"].(string)

	// --- Setup (validation errors raise: they are the caller's to fix) ------
	repoPath, err := filepath.Abs(in.RepoPath)
	if err != nil {
		return nil, err
	}
	if err := ensureIssueReadyRepo(repoPath); err != nil {
		return nil, err
	}
	baseRef, baseSHA, err := resolveBase(repoPath, in.BaseBranch)
	if err != nil {
		return nil, err
	}
	excludes := []string{".worktrees/"}
	if !filepath.IsAbs(in.ArtifactsDir) {
		top := topPathSegment(in.ArtifactsDir)
		if top != "" && top != "." && top != ".." {
			excludes = append(excludes, top+"/")
		}
	}
	if err := ensureLocalExcludes(repoPath, excludes); err != nil {
		return nil, err
	}
	if isDirty(repoPath) {
		deps.note(ctx,
			"Caller repo has uncommitted changes; the issue branch is created "+
				"from the committed base state and will not see them",
			"issue_build", "warning", "dirty_tree",
		)
	}

	buildID := newBuildID()
	resuming := false
	if in.ResumeBuildID != "" {
		if len(in.ResumeBuildID) != 8 {
			return nil, fmt.Errorf("implement_issue: invalid resume_build_id %q", in.ResumeBuildID)
		}
		if _, err := hex.DecodeString(in.ResumeBuildID); err != nil {
			return nil, fmt.Errorf("implement_issue: invalid resume_build_id %q", in.ResumeBuildID)
		}
		buildID = in.ResumeBuildID
		resuming = true
	}
	branch := cfg.BranchPrefix + buildID + "-" + plannedName
	worktreePath := filepath.Join(repoPath, ".worktrees", buildID+"-"+plannedName)
	var absArtifacts string
	if filepath.IsAbs(in.ArtifactsDir) {
		absArtifacts = filepath.Join(in.ArtifactsDir, "issue-builds", buildID)
	} else {
		absArtifacts = filepath.Join(repoPath, in.ArtifactsDir, "issue-builds", buildID)
	}
	if err := os.MkdirAll(absArtifacts, 0o755); err != nil {
		return nil, err
	}

	if resuming {
		if info, err := os.Stat(worktreePath); err != nil || !info.IsDir() {
			return nil, fmt.Errorf("implement_issue: resume worktree missing for build %s", buildID)
		}
		gotBranch, err := currentBranch(worktreePath)
		if err != nil || gotBranch != branch {
			return nil, fmt.Errorf("implement_issue: resume worktree branch mismatch: got %q, want %q", gotBranch, branch)
		}
		if _, _, code := runGit(repoPath, "merge-base", "--is-ancestor", baseSHA, branch); code != 0 {
			return nil, fmt.Errorf("implement_issue: resume branch %s is not based on %s", branch, baseSHA)
		}
		deps.note(ctx, fmt.Sprintf("Issue build %s: resuming existing worktree %s", buildID, worktreePath),
			"issue_build", "resume")
	} else if err := addWorktree(repoPath, worktreePath, branch, baseSHA); err != nil {
		return nil, err
	}
	deps.note(ctx,
		fmt.Sprintf("Issue build %s: %s on %s (base %s @ %.12s)",
			buildID, plannedName, branch, baseRef, baseSHA),
		"issue_build", "start",
	)

	planned["worktree_path"] = worktreePath
	planned["branch_name"] = branch

	execCfg, err := config.LoadExecutionConfig(cfg.ToExecutionRaw())
	if err != nil {
		if !resuming {
			removeWorktree(repoPath, worktreePath)
			deleteBranch(repoPath, branch)
		}
		return nil, err
	}
	dagState := &schemas.DAGState{
		RepoPath:     worktreePath,
		ArtifactsDir: absArtifacts,
		BuildID:      buildID,
		AllIssues:    []map[string]any{planned},
	}

	// --- Execute (failures after branch creation → structured result) ------
	outcomeValue := "error"
	loopSummary := ""
	errorMessage := ""
	iterations := 0
	var iterationHistory []map[string]any
	var debtItems []map[string]any
	var verification map[string]any
	prURL := ""
	var commits []string
	var filesChanged []string
	stat := ""

	noteFn := func(msg string, tags []string) { deps.note(ctx, msg, tags...) }
	callFn := coding.CallFn(deps.Call)

	loopResult, loopErr := coding.RunCodingLoop(
		ctx, planned, dagState, callFn, deps.NodeID, execCfg, noteFn, nil,
	)
	if loopErr != nil && ctx.Err() != nil {
		// A cancelled/expired transport can race with a completed mutation. Observe
		// authoritative Git post-state before deciding whether continuation is
		// still required. Never claim completion here: reviewer/QA/verifier may not
		// have run. When an effect is observable, preserve the worktree and return a
		// structured recoverable result so an explicit resume can continue the same
		// logical build without repeating the mutation.
		commits = newCommits(repoPath, baseSHA, branch)
		dirty := isDirty(worktreePath)
		if len(commits) > 0 || dirty {
			filesChanged = changedFiles(repoPath, baseSHA, branch)
			stat = diffStat(repoPath, baseSHA, branch)
			if commits == nil {
				commits = []string{}
			}
			if filesChanged == nil {
				filesChanged = []string{}
			}
			return map[string]any{
				"success":           false,
				"outcome":           "interrupted_with_progress",
				"summary":           "Execution was interrupted after observable issue progress; resume the same build_id to continue verification.",
				"build_id":          buildID,
				"branch":            branch,
				"base_branch":       baseRef,
				"base_sha":          baseSHA,
				"commits":           commits,
				"files_changed":     filesChanged,
				"diff_stat":         stat,
				"iterations":        0,
				"iteration_history": []map[string]any{},
				"debt_items":        []map[string]any{},
				"verification":      nil,
				"pr_url":            "",
				"error_message":     loopErr.Error(),
				"effect_observed":   true,
				"worktree_path":     worktreePath,
			}, nil
		}
		// No observable effect: keep the previous cancellation semantics. A resumed
		// build still preserves its existing worktree; a fresh untouched build can
		// be removed safely.
		if !resuming {
			removeWorktree(repoPath, worktreePath)
		}
		return nil, loopErr
	}
	if loopErr != nil {
		// Fatal harness errors become a structured failure (Python's generic
		// `except Exception` path).
		errorMessage = loopErr.Error()
		deps.note(ctx, fmt.Sprintf("Issue build failed: %v", loopErr), "issue_build", "error")
		commits = newCommits(repoPath, baseSHA, branch)
		filesChanged = changedFiles(repoPath, baseSHA, branch)
		stat = diffStat(repoPath, baseSHA, branch)
	} else {
		outcomeValue = string(loopResult.Outcome)
		loopSummary = loopResult.ResultSummary
		if loopSummary == "" {
			loopSummary = loopResult.ErrorMessage
		}
		errorMessage = loopResult.ErrorMessage
		iterations = loopResult.Attempts
		iterationHistory = loopResult.IterationHistory
		debtItems = loopResult.DebtItems

		if _, err := scrubTrackedJunk(worktreePath, plannedName); err != nil {
			deps.note(ctx, fmt.Sprintf("Junk scrub failed: %v", err),
				"issue_build", "error")
		}
		if _, err := commitAll(
			worktreePath,
			fmt.Sprintf("chore(%s): checkpoint uncommitted issue work", plannedName),
		); err != nil {
			deps.note(ctx, fmt.Sprintf("Checkpoint commit failed: %v", err),
				"issue_build", "error")
		}
		commits = newCommits(repoPath, baseSHA, branch)
		filesChanged = changedFiles(repoPath, baseSHA, branch)
		stat = diffStat(repoPath, baseSHA, branch)

		if violations := declaredScopeViolations(spec, filesChanged); len(violations) > 0 {
			outcomeValue = "failed_unrecoverable"
			errorMessage = fmt.Sprintf("out-of-scope files changed: %s", strings.Join(violations, ", "))
			loopSummary = errorMessage
			deps.note(ctx, errorMessage, "issue_build", "scope", "error")
		}

		codingOK := completedOutcomes[outcomeValue]
		if cfg.Verify && codingOK && len(commits) > 0 {
			verification = runVerification(ctx, deps, cfg, execCfg, spec, planned,
				worktreePath, absArtifacts, loopSummary)
		}
		if cfg.EnableGithubPR && codingOK && len(commits) > 0 {
			prURL = maybeCreatePR(ctx, deps, cfg, execCfg, spec, planned,
				repoPath, worktreePath, branch, baseRef, absArtifacts,
				loopSummary, debtItems)
		}
	}

	if !cfg.KeepWorktree && !resuming {
		removeWorktree(repoPath, worktreePath)
	}
	if len(commits) == 0 && !resuming {
		// A fresh branch with zero commits is pure noise for the caller. A resumed
		// branch may contain the only copy of partial work and must be preserved.
		deleteBranch(repoPath, branch)
	}

	codingOK := completedOutcomes[outcomeValue]
	verifyOK := verification == nil || isTruthyBool(verification["passed"])
	success := codingOK && len(commits) > 0 && verifyOK

	summary := fmt.Sprintf("%s: %s, %d commit(s), %d file(s) changed",
		plannedName, outcomeValue, len(commits), len(filesChanged))
	if len(commits) > 0 {
		summary += " on " + branch
	} else {
		summary += "; no commits produced"
	}
	if verification != nil {
		if isTruthyBool(verification["passed"]) {
			summary += "; verification passed"
		} else {
			summary += "; verification FAILED"
		}
	}
	if loopSummary != "" {
		summary += " — " + loopSummary
	}

	completeTag := "complete"
	if !success {
		completeTag = "failed"
	}
	deps.note(ctx, summary, "issue_build", completeTag)

	resultBranch := branch
	if len(commits) == 0 {
		resultBranch = ""
	}
	finalError := errorMessage
	if success {
		finalError = ""
	}
	if commits == nil {
		commits = []string{}
	}
	if filesChanged == nil {
		filesChanged = []string{}
	}
	if iterationHistory == nil {
		iterationHistory = []map[string]any{}
	}
	if debtItems == nil {
		debtItems = []map[string]any{}
	}

	return map[string]any{
		"success":           success,
		"outcome":           outcomeValue,
		"summary":           summary,
		"build_id":          buildID,
		"branch":            resultBranch,
		"base_branch":       baseRef,
		"base_sha":          baseSHA,
		"commits":           commits,
		"files_changed":     filesChanged,
		"diff_stat":         stat,
		"iterations":        iterations,
		"iteration_history": iterationHistory,
		"debt_items":        debtItems,
		"verification":      anyOrNil(verification),
		"pr_url":            prURL,
		"error_message":     finalError,
	}, nil
}

// runVerification ports build._run_verification — one verifier pass whose
// unavailability reports passed=false, never success.
func runVerification(
	ctx context.Context,
	deps *Deps,
	cfg *config.IssueBuildConfig,
	execCfg *config.ExecutionConfig,
	spec *Spec,
	planned map[string]any,
	worktreePath, artifactsDir, loopSummary string,
) map[string]any {
	callCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.AgentTimeoutSeconds)*time.Second)
	defer cancel()

	prd := map[string]any{
		"validated_description": spec.Title + "\n\n" + spec.Description,
		"acceptance_criteria":   toAnySlice(spec.AcceptanceCriteria),
		"must_have":             toAnySlice(spec.AcceptanceCriteria),
		"nice_to_have":          []any{},
		"out_of_scope":          []any{},
	}
	result, err := deps.Call(callCtx, deps.NodeID+".run_verifier", map[string]any{
		"prd":           prd,
		"repo_path":     worktreePath,
		"artifacts_dir": artifactsDir,
		"completed_issues": []any{map[string]any{
			"issue_name":     planned["name"],
			"result_summary": loopSummary,
		}},
		"failed_issues":   []any{},
		"skipped_issues":  []any{},
		"model":           execCfg.VerifierModel(),
		"permission_mode": cfg.PermissionMode,
		"ai_provider":     execCfg.AIProvider(),
	})
	if err != nil {
		deps.note(ctx, fmt.Sprintf("Verifier unavailable: %v", err),
			"issue_build", "verify", "error")
		return map[string]any{
			"passed":  false,
			"summary": fmt.Sprintf("Verification unavailable: %v", err),
		}
	}
	return result
}

// maybeCreatePR ports build._maybe_create_pr — best-effort PR creation.
func maybeCreatePR(
	ctx context.Context,
	deps *Deps,
	cfg *config.IssueBuildConfig,
	execCfg *config.ExecutionConfig,
	spec *Spec,
	planned map[string]any,
	repoPath, worktreePath, branch, baseRef, artifactsDir, loopSummary string,
	debtItems []map[string]any,
) string {
	if remoteURL(repoPath) == "" {
		deps.note(ctx, "No origin remote; skipping PR creation",
			"issue_build", "github_pr", "skip")
		return ""
	}
	baseForPR := cfg.GithubPRBase
	if baseForPR == "" {
		baseForPR = defaultRemoteBranch(repoPath)
	}
	if baseForPR == "" {
		baseForPR = baseRef
		if baseForPR == "HEAD" {
			baseForPR = "main"
		}
	}
	buildSummary := loopSummary
	if buildSummary == "" {
		buildSummary = spec.Title
	}
	debtAny := make([]any, len(debtItems))
	for i, d := range debtItems {
		debtAny[i] = d
	}

	callCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.AgentTimeoutSeconds)*time.Second)
	defer cancel()
	result, err := deps.Call(callCtx, deps.NodeID+".run_github_pr", map[string]any{
		"repo_path":          worktreePath,
		"integration_branch": branch,
		"base_branch":        baseForPR,
		"goal":               spec.Title,
		"build_summary":      buildSummary,
		"completed_issues": []any{map[string]any{
			"issue_name":     planned["name"],
			"result_summary": loopSummary,
		}},
		"accumulated_debt": debtAny,
		"artifacts_dir":    artifactsDir,
		"model":            execCfg.GitModel(),
		"permission_mode":  cfg.PermissionMode,
		"ai_provider":      execCfg.AIProvider(),
	})
	if err != nil {
		deps.note(ctx, fmt.Sprintf("PR creation failed (non-fatal): %v", err),
			"issue_build", "github_pr", "error")
		return ""
	}
	if url, ok := result["pr_url"].(string); ok {
		return url
	}
	return ""
}

// topPathSegment returns the first path component of a relative path — the
// pattern added to .git/info/exclude for the artifacts dir.
func topPathSegment(p string) string {
	p = strings.Trim(p, "/")
	if i := strings.IndexByte(p, '/'); i >= 0 {
		return p[:i]
	}
	return p
}

func isTruthyBool(v any) bool {
	b, ok := v.(bool)
	return ok && b
}

func anyOrNil(m map[string]any) any {
	if m == nil {
		return nil
	}
	return m
}
