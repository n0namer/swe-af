package advisor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Agent-Field/SWE-AF/go/internal/config"
	"github.com/Agent-Field/SWE-AF/go/internal/dagutil"
	"github.com/Agent-Field/SWE-AF/go/internal/harnessx"
	"github.com/Agent-Field/SWE-AF/go/internal/prompts/advisor"
	"github.com/Agent-Field/SWE-AF/go/internal/prompts/coding"
	"github.com/Agent-Field/SWE-AF/go/internal/runtimex"
	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
)

// ---------------------------------------------------------------------------
// run_verifier
// ---------------------------------------------------------------------------

type verifierInput struct {
	PRD               map[string]any   `json:"prd"`
	RepoPath          string           `json:"repo_path"`
	ArtifactsDir      string           `json:"artifacts_dir"`
	CompletedIssues   []map[string]any `json:"completed_issues"`
	FailedIssues      []map[string]any `json:"failed_issues"`
	SkippedIssues     []string         `json:"skipped_issues"`
	Model             string           `json:"model"`
	PermissionMode    string           `json:"permission_mode"`
	AIProvider        string           `json:"ai_provider"`
	WorkspaceManifest map[string]any   `json:"workspace_manifest"`
}

// RunVerifier runs final acceptance verification against the PRD. Ports
// run_verifier. On agent failure it returns an inconclusive VerificationResult.
func RunVerifier(ctx context.Context, deps *Deps, input map[string]any) (any, error) {
	in := verifierInput{Model: "sonnet", AIProvider: "claude"}
	if err := decodeInput(input, &in); err != nil {
		return nil, err
	}

	deps.note(ctx, "Verifier starting", "verifier", "start")

	wsManifest, err := maybeWorkspaceManifest(in.WorkspaceManifest)
	if err != nil {
		return nil, err
	}

	taskPrompt := coding.VerifierTaskPrompt(coding.VerifierTaskPromptOpts{
		PRD:               in.PRD,
		ArtifactsDir:      in.ArtifactsDir,
		CompletedIssues:   in.CompletedIssues,
		FailedIssues:      in.FailedIssues,
		SkippedIssues:     in.SkippedIssues,
		WorkspaceManifest: wsManifest,
	})

	provider, err := runtimex.RuntimeToHarnessAdapter(in.AIProvider)
	if err != nil {
		return nil, err
	}

	verifierEnv := map[string]string{}
	if provider == "opencode" {
		// OpenCode does not currently enforce harness.Options.Tools. Verification
		// must remain noninteractive and repository-local: hidden subagents can
		// otherwise stall on external-directory permission prompts.
		verifierEnv["OPENCODE_CONFIG_CONTENT"] = `{"permission":{"task":"deny","external_directory":"deny"}}`
		for _, venvName := range []string{".venv", "venv"} {
			venvBin := filepath.Join(in.RepoPath, venvName, "bin")
			if info, statErr := os.Stat(venvBin); statErr == nil && info.IsDir() {
				verifierEnv["PATH"] = venvBin + string(os.PathListSeparator) + os.Getenv("PATH")
				break
			}
		}
	}
	opts := harnessx.RoleOptions{
		Provider:       provider,
		Model:          in.Model,
		MaxTurns:       config.DefaultAgentMaxTurns,
		Tools:          advisorTools,
		PermissionMode: in.PermissionMode,
		SystemPrompt:   coding.VerifierSystemPrompt,
		Cwd:            in.RepoPath,
		Env:            verifierEnv,
		SchemaMode:     "single",
	}.ToOptions()

	untrackedBefore, untrackedErr := verifierUntrackedPaths(in.RepoPath)
	if untrackedErr != nil {
		deps.note(ctx, fmt.Sprintf("Verifier scratch tracking unavailable: %v", untrackedErr), "verifier", "warning")
	}

	parsed, result, err := harnessx.Run[schemas.VerificationResult](ctx, deps.Harness, taskPrompt, opts)
	if untrackedErr == nil {
		archived, archiveRoot, archiveErr := archiveNewVerifierScratch(in.RepoPath, untrackedBefore)
		if archiveErr != nil {
			return nil, fmt.Errorf("verifier scratch archival failed: %w", archiveErr)
		}
		if len(archived) > 0 {
			deps.note(ctx, fmt.Sprintf("Verifier archived scratch files: %s -> %s", strings.Join(archived, ", "), archiveRoot), "verifier", "scratch_archived")
		}
	}
	if err != nil {
		if isFatal(err) {
			return nil, err
		}
		deps.note(ctx, fmt.Sprintf("Verifier agent failed: %v", err), "verifier", "error")
	} else if result != nil && result.Parsed != nil {
		deps.note(ctx, fmt.Sprintf(
			"Verifier complete: passed=%s, summary=%s", pyBool(parsed.Passed), parsed.Summary),
			"verifier", "complete")
		return parsed, nil
	}

	// Fallback: verification inconclusive.
	return schemas.VerificationResult{
		Passed:          false,
		CriteriaResults: []schemas.CriterionResult{},
		Summary:         "Verifier agent failed to produce a valid result.",
		SuggestedFixes:  []string{"Re-run verification manually."},
	}, nil
}

func verifierUntrackedPaths(repoPath string) (map[string]struct{}, error) {
	cmd := exec.Command("git", "-C", repoPath, "ls-files", "--others", "--exclude-standard", "-z")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	paths := map[string]struct{}{}
	for _, rel := range strings.Split(string(out), "\x00") {
		if rel != "" {
			paths[rel] = struct{}{}
		}
	}
	return paths, nil
}

func archiveNewVerifierScratch(repoPath string, before map[string]struct{}) ([]string, string, error) {
	after, err := verifierUntrackedPaths(repoPath)
	if err != nil {
		return nil, "", err
	}
	var created []string
	for rel := range after {
		if _, existed := before[rel]; !existed {
			created = append(created, rel)
		}
	}
	if len(created) == 0 {
		return nil, "", nil
	}
	sort.Strings(created)

	archiveBase := filepath.Join(filepath.Dir(repoPath), ".verifier-scratch", verifierScratchSegment(filepath.Base(repoPath)))
	if err := os.MkdirAll(archiveBase, 0o755); err != nil {
		return nil, "", err
	}
	archiveRoot, err := os.MkdirTemp(archiveBase, "run-")
	if err != nil {
		return nil, "", err
	}
	for _, rel := range created {
		clean := filepath.Clean(rel)
		if clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
			return nil, archiveRoot, fmt.Errorf("unsafe verifier scratch path %q", rel)
		}
		src := filepath.Join(repoPath, clean)
		dst := filepath.Join(archiveRoot, clean)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return nil, archiveRoot, err
		}
		if err := os.Rename(src, dst); err != nil {
			return nil, archiveRoot, fmt.Errorf("archive verifier scratch %q: %w", rel, err)
		}
	}
	return created, archiveRoot, nil
}

func verifierScratchSegment(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	if b.Len() == 0 {
		return "unknown"
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// generate_fix_issues
// ---------------------------------------------------------------------------

type fixGeneratorInput struct {
	FailedCriteria    []map[string]any `json:"failed_criteria"`
	DAGState          map[string]any   `json:"dag_state"`
	PRD               map[string]any   `json:"prd"`
	ArtifactsDir      string           `json:"artifacts_dir"`
	Model             string           `json:"model"`
	PermissionMode    string           `json:"permission_mode"`
	AIProvider        string           `json:"ai_provider"`
	WorkspaceManifest map[string]any   `json:"workspace_manifest"`
}

// fixGeneratorOutput is the local output schema of generate_fix_issues (the
// Python reasoner declares it inline). The list defaults are seeded to empty
// (non-nil) so a parsed result serializes fix_issues/debt_items as [] rather
// than null, matching pydantic model_dump.
type fixGeneratorOutput struct {
	FixIssues []map[string]any `json:"fix_issues"`
	DebtItems []map[string]any `json:"debt_items"`
	Summary   string           `json:"summary"`
}

func defaultFixGeneratorOutput() fixGeneratorOutput {
	return fixGeneratorOutput{FixIssues: []map[string]any{}, DebtItems: []map[string]any{}}
}

// UnmarshalJSON seeds the empty-list defaults before decoding.
func (o *fixGeneratorOutput) UnmarshalJSON(b []byte) error {
	*o = defaultFixGeneratorOutput()
	type alias fixGeneratorOutput
	return json.Unmarshal(b, (*alias)(o))
}

// GenerateFixIssues generates targeted fix issues from failed verification
// criteria. Ports generate_fix_issues. On agent failure it records all criteria
// as debt. Uses the verifier's model per the caller (the reasoner input is
// simply `model`).
func GenerateFixIssues(ctx context.Context, deps *Deps, input map[string]any) (any, error) {
	in := fixGeneratorInput{Model: "sonnet", AIProvider: "claude"}
	if err := decodeInput(input, &in); err != nil {
		return nil, err
	}

	deps.note(ctx, fmt.Sprintf("Fix generator starting: %d failed criteria", len(in.FailedCriteria)), "fix_generator", "start")

	repoPath := getStr(in.DAGState, "repo_path", ".")
	taskPrompt := advisor.FixGeneratorTaskPrompt(advisor.FixGeneratorTaskOptions{
		FailedCriteria:  in.FailedCriteria,
		DAGStateSummary: in.DAGState,
		PRD:             in.PRD,
	})

	// If multi-repo, instruct the agent to set target_repo on each fix issue.
	wsManifest, err := maybeWorkspaceManifest(in.WorkspaceManifest)
	if err != nil {
		return nil, err
	}
	if wsManifest != nil && len(wsManifest.Repos) > 1 {
		taskPrompt += "\n\n## Multi-Repo Context\n" +
			"This workspace spans multiple repositories. For each fix issue you generate, " +
			"include a `target_repo` field specifying which repository the fix should be " +
			"applied to. Available repos:\n"
		for _, repo := range wsManifest.Repos {
			taskPrompt += fmt.Sprintf("- **%s** (role: %s): `%s`\n", repo.RepoName, repo.Role, repo.AbsolutePath)
		}
	}

	provider, err := runtimex.RuntimeToHarnessAdapter(in.AIProvider)
	if err != nil {
		return nil, err
	}

	opts := harnessx.RoleOptions{
		Provider:       provider,
		Model:          in.Model,
		MaxTurns:       config.DefaultAgentMaxTurns,
		Tools:          advisorTools,
		PermissionMode: in.PermissionMode,
		SystemPrompt:   advisor.FixGeneratorSystemPrompt,
		Cwd:            repoPath,
	}.ToOptions()

	parsed, result, err := harnessx.Run[fixGeneratorOutput](ctx, deps.Harness, taskPrompt, opts)
	if err != nil {
		if isFatal(err) {
			return nil, err
		}
		deps.note(ctx, fmt.Sprintf("Fix generator agent failed: %v", err), "fix_generator", "error")
	} else if result != nil && result.Parsed != nil {
		deps.note(ctx, fmt.Sprintf(
			"Fix generator complete: %d fix issues, %d debt items", len(parsed.FixIssues), len(parsed.DebtItems)),
			"fix_generator", "complete")
		// fix_issues are raw dicts headed for DAGState.all_issues — coerce
		// LLM scalar shapes (str acceptance_criteria etc.) at the boundary.
		for i, fi := range parsed.FixIssues {
			parsed.FixIssues[i] = dagutil.NormalizeIssueDict(fi)
		}
		return parsed, nil
	}

	// Fallback: record all criteria as debt.
	debtItems := make([]map[string]any, 0, len(in.FailedCriteria))
	for _, c := range in.FailedCriteria {
		debtItems = append(debtItems, map[string]any{
			"criterion": getStr(c, "criterion", ""),
			"reason":    "Fix generator failed to analyze",
			"severity":  "high",
		})
	}
	return map[string]any{
		"fix_issues": []any{},
		"debt_items": debtItems,
		"summary":    "Fix generator failed — all criteria recorded as debt",
	}, nil
}

// ---------------------------------------------------------------------------
// run_issue_writer
// ---------------------------------------------------------------------------

type issueWriterInput struct {
	Issue               map[string]any   `json:"issue"`
	PRDSummary          string           `json:"prd_summary"`
	ArchitectureSummary string           `json:"architecture_summary"`
	IssuesDir           string           `json:"issues_dir"`
	RepoPath            string           `json:"repo_path"`
	PRDPath             string           `json:"prd_path"`
	ArchitecturePath    string           `json:"architecture_path"`
	SiblingIssues       []map[string]any `json:"sibling_issues"`
	Model               string           `json:"model"`
	PermissionMode      string           `json:"permission_mode"`
	AIProvider          string           `json:"ai_provider"`
	WorkspaceManifest   map[string]any   `json:"workspace_manifest"`
}

// issueWriterOutput is the local output schema of run_issue_writer.
type issueWriterOutput struct {
	IssueName     string `json:"issue_name"`
	IssueFilePath string `json:"issue_file_path"`
	Success       bool   `json:"success"`
}

// RunIssueWriter writes a lean issue-*.md file for a new or updated issue. Ports
// run_issue_writer. On agent failure it returns success=false without blocking.
func issueWriterPlanningMaxTurnsDefault() int {
	if n, err := strconv.Atoi(os.Getenv("SWE_PLANNING_MAX_TURNS")); err == nil && n > 0 {
		return n
	}
	return 2
}

func RunIssueWriter(ctx context.Context, deps *Deps, input map[string]any) (any, error) {
	in := issueWriterInput{Model: "sonnet", AIProvider: "claude"}
	if err := decodeInput(input, &in); err != nil {
		return nil, err
	}

	issueName := getStr(in.Issue, "name", "unknown")
	deps.note(ctx, fmt.Sprintf("Issue writer starting for %s", issueName), "issue_writer", "start")

	wsManifest, err := maybeWorkspaceManifest(in.WorkspaceManifest)
	if err != nil {
		return nil, err
	}

	taskPrompt := coding.IssueWriterTaskPrompt(coding.IssueWriterTaskPromptOpts{
		Issue:               in.Issue,
		PRDSummary:          in.PRDSummary,
		ArchitectureSummary: in.ArchitectureSummary,
		IssuesDir:           in.IssuesDir,
		PRDPath:             in.PRDPath,
		ArchitecturePath:    in.ArchitecturePath,
		SiblingIssues:       in.SiblingIssues,
		WorkspaceManifest:   wsManifest,
	})

	provider, err := runtimex.RuntimeToHarnessAdapter(in.AIProvider)
	if err != nil {
		return nil, err
	}

	opts := harnessx.RoleOptions{
		Provider:       provider,
		Model:          in.Model,
		MaxTurns:       issueWriterPlanningMaxTurnsDefault(),
		Tools:          []string{"Read", "Write", "Glob", "Grep"},
		PermissionMode: in.PermissionMode,
		SystemPrompt:   coding.IssueWriterSystemPrompt,
		Cwd:            in.RepoPath,
	}.ToOptions()

	parsed, result, err := harnessx.Run[issueWriterOutput](ctx, deps.Harness, taskPrompt, opts)
	if err != nil {
		if isFatal(err) {
			return nil, err
		}
		deps.note(ctx, fmt.Sprintf("Issue writer failed for %s: %v", issueName, err), "issue_writer", "error")
	} else if result != nil && result.Parsed != nil {
		deps.note(ctx, fmt.Sprintf("Issue writer complete for %s: %s", issueName, parsed.IssueFilePath), "issue_writer", "complete")
		return parsed, nil
	}

	// Fallback: issue file wasn't written but we don't block on it.
	return map[string]any{
		"issue_name":      issueName,
		"issue_file_path": "",
		"success":         false,
	}, nil
}
