// This file (resume.go) ports the resume_build reasoner from
// swe_af/app.py:2081-2134 and adds native checkpoint fencing for restart-safe
// execution. Resume metadata is carried through the existing execute envelope so
// the public reasoner contract remains unchanged.
package orch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// ResumeBuildHandler satisfies the orchestrator Handler shape.
var _ Handler = ResumeBuildHandler

var (
	errResumeStaleOwner = errors.New("resume rejected: stale owner or lease")
)

// ResumeBuildHandler resumes a crashed build from the last checkpoint. It
// preserves the Python raw execute envelope while fencing stale owners,
// selecting dependency-ready work, and persisting the typed result only after a
// successful, non-cancelled dispatch.
func ResumeBuildHandler(ctx context.Context, deps *Deps, input map[string]any) (any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	repoPath := mapStr(input, "repo_path", "")
	artifactsDir := mapStr(input, "artifacts_dir", ".artifacts")

	absRepo, err := filepath.Abs(repoPath)
	if err != nil {
		absRepo = repoPath
	}
	base := filepath.Join(absRepo, artifactsDir)
	planPath := filepath.Join(base, "execution", "checkpoint.json")

	if _, statErr := os.Stat(planPath); statErr != nil {
		return nil, fmt.Errorf("No checkpoint found at %s. Cannot resume.", planPath)
	}

	raw, err := os.ReadFile(planPath)
	if err != nil {
		return nil, fmt.Errorf("No checkpoint found at %s. Cannot resume.", planPath)
	}
	var checkpoint map[string]any
	if err := json.Unmarshal(raw, &checkpoint); err != nil {
		return nil, err
	}
	if deps == nil || deps.App == nil {
		return nil, errors.New("resume requires an app dependency")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	owner := firstString(input, "owner_id", "lease_owner", "owner")
	requestedEpoch, hasEpoch := mapInt(input, "lease_epoch")
	persistedOwner := firstString(checkpoint, "owner_id", "lease_owner", "owner")
	persistedEpoch, persistedEpochSet := mapInt(checkpoint, "lease_epoch")
	if owner != "" && persistedOwner != "" {
		if !hasEpoch || (persistedEpochSet && requestedEpoch < persistedEpoch) ||
			(requestedEpoch == persistedEpoch && owner != persistedOwner) {
			return nil, errResumeStaleOwner
		}
	}
	if owner != "" && !hasEpoch && persistedEpochSet {
		return nil, errResumeStaleOwner
	}
	if owner != "" && persistedOwner == "" {
		checkpoint["owner_id"] = owner
		persistedOwner = owner
	}
	if hasEpoch && (!persistedEpochSet || requestedEpoch > persistedEpoch) {
		checkpoint["lease_epoch"] = requestedEpoch
		persistedEpoch = requestedEpoch
		persistedEpochSet = true
	}

	issues := any0(checkpoint["all_issues"])
	levels := any0(checkpoint["levels"])
	completed := completedIssueNames(checkpoint)
	ready := resumeReadyIssues(issues, levels, completed)

	// Reconstruct plan_result exactly as app.py:2112-2121, with additive
	// fencing fields consumed by the native execute/fake-runner path.
	planResult := map[string]any{
		"prd":              map[string]any{},
		"architecture":     map[string]any{},
		"review":           map[string]any{},
		"issues":           issues,
		"levels":           levels,
		"file_conflicts":   []any{},
		"artifacts_dir":    mapStr(checkpoint, "artifacts_dir", base),
		"rationale":        mapStr(checkpoint, "original_plan_summary", ""),
		"ready_issues":     ready,
		"completed_issues": checkpointValueOrEmpty(checkpoint, "completed_issues"),
		"last_resume":      checkpoint["last_resume"],
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	deps.Note(ctx, "Resuming build from checkpoint", "build", "resume")

	resumeInput := map[string]any{
		"plan_result": planResult,
		"repo_path":   repoPath,
		"config":      input["config"],
		"git_config":  input["git_config"],
		"resume":      true,
	}
	if owner != "" {
		resumeInput["owner_id"] = owner
	}
	if persistedEpochSet {
		resumeInput["lease_epoch"] = persistedEpoch
	}
	if runID := firstString(input, "run_id", "build_id"); runID != "" {
		resumeInput["run_id"] = runID
	}
	resumeInput["resume_dispatch"] = resumeDispatchKeys(firstString(input, "run_id", "build_id"), owner, persistedEpoch, ready)

	result, err := deps.CallRaw(ctx, "execute", resumeInput)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Persist only a completed, typed result. The dispatch key makes a replay
	// idempotent even if the process crashes after the runner effect and before
	// this write becomes visible.
	if owner == "" && !hasEpoch {
		return result, nil
	}
	checkpoint["last_resume_result"] = result
	checkpoint["last_resume_dispatch"] = resumeInput["resume_dispatch"]
	if owner != "" {
		checkpoint["owner_id"] = owner
	}
	if persistedEpochSet {
		checkpoint["lease_epoch"] = persistedEpoch
	}
	if err := persistResumeCheckpoint(planPath, checkpoint); err != nil {
		return nil, err
	}
	return result, nil
}

func firstString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := mapStr(m, key, ""); value != "" {
			return value
		}
	}
	return ""
}

func mapInt(m map[string]any, key string) (int, bool) {
	if m == nil {
		return 0, false
	}
	value, ok := m[key]
	if !ok || value == nil {
		return 0, false
	}
	switch typed := value.(type) {
	case int:
		return typed, true
	case int64:
		return int(typed), true
	case float64:
		return int(typed), true
	case json.Number:
		parsed, err := strconv.Atoi(string(typed))
		return parsed, err == nil
	case string:
		parsed, err := strconv.Atoi(typed)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func checkpointValueOrEmpty(checkpoint map[string]any, key string) any {
	if value, ok := checkpoint[key]; ok && value != nil {
		return value
	}
	return []any{}
}

func completedIssueNames(checkpoint map[string]any) map[string]bool {
	completed := map[string]bool{}
	for _, raw := range asMapList(checkpoint["completed_issues"]) {
		name := firstString(raw, "issue_name", "name", "id")
		if name != "" && !isFailedIssue(raw) {
			completed[name] = true
		}
	}
	for _, name := range asStrList(checkpoint["completed_issue_names"]) {
		if name != "" {
			completed[name] = true
		}
	}
	return completed
}

func isFailedIssue(issue map[string]any) bool {
	status := firstString(issue, "outcome", "status")
	return status == "failed" || status == "cancelled" || status == "canceled"
}

func resumeReadyIssues(issues, levels any, completed map[string]bool) []string {
	dependencies := map[string][]string{}
	issueNames := map[string]bool{}
	for _, raw := range asMapList(issues) {
		name := firstString(raw, "name", "issue_name", "id")
		if name == "" {
			continue
		}
		issueNames[name] = true
		dependencies[name] = asStrList(raw["depends_on"])
	}

	orderedLevels := asAnyLists(levels)
	for _, level := range orderedLevels {
		ready := make([]string, 0, len(level))
		for _, rawName := range level {
			name, ok := rawName.(string)
			if !ok || name == "" || completed[name] || !issueNames[name] {
				continue
			}
			allDone := true
			for _, dependency := range dependencies[name] {
				if !completed[dependency] {
					allDone = false
					break
				}
			}
			if allDone {
				ready = append(ready, name)
			}
		}
		if len(ready) > 0 {
			return ready
		}
	}
	return nil
}

func asAnyLists(value any) [][]any {
	var lists [][]any
	switch typed := value.(type) {
	case [][]any:
		return typed
	case []any:
		for _, raw := range typed {
			if list, ok := raw.([]any); ok {
				lists = append(lists, list)
			}
		}
	}
	return lists
}

func resumeDispatchKeys(runID, owner string, epoch int, ready []string) []map[string]any {
	keys := make([]map[string]any, 0, len(ready))
	for _, issue := range ready {
		keys = append(keys, map[string]any{
			"issue": issue,
			"key":   fmt.Sprintf("resume:%s:%s:%d:%s", runID, owner, epoch, issue),
		})
	}
	return keys
}

func persistResumeCheckpoint(path string, checkpoint map[string]any) error {
	raw, err := json.Marshal(checkpoint)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
