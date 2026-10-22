"""Fail-closed auto-resume admission for the existing SWE-AF Recovery Watch.

One durable dispatch attempt per failed root. No new scheduler. Never assume a
timeout implies the old writer stopped; unknown evidence blocks automatic resume.
"""
import json
import os
import subprocess
from pathlib import Path

def git(repo, *args):
    p = subprocess.run(["git", "-C", str(repo), *args], capture_output=True, text=True, timeout=10)
    if p.returncode:
        raise RuntimeError("git failed: " + " ".join(args))
    return p.stdout.strip()

def guarded_payload(project, root, build_id, repo, detail, checkpoint, plan, *, active_runs, child_executions, process_commands):
    """Return (payload, blocker). The caller owns the single-instance lock."""
    if active_runs:
        return None, "active_agentfield_execution"
    if detail.get("status") not in ("failed", "timeout"):
        return None, "not_failed_transport"
    error = (detail.get("error_message") or "").lower()
    if not any(word in error for word in ("timed out", "timeout", "cancelled by upstream client", "connection lost", "crashed")):
        return None, "logical_failure_not_recoverable"
    if not isinstance(checkpoint, dict) or checkpoint.get("build_id") != build_id:
        return None, "checkpoint_build_mismatch"
    if os.path.normcase(str(checkpoint.get("repo_path", ""))) != os.path.normcase(str(repo)):
        return None, "checkpoint_repo_mismatch"
    if checkpoint.get("completed_issues") and not checkpoint.get("merge_results"):
        return None, "completed_without_merge_receipts"
    if checkpoint.get("failed_issues"):
        return None, "logical_issue_failures"
    if not checkpoint.get("all_issues") or not checkpoint.get("levels"):
        return None, "empty_checkpoint"
    if not isinstance(plan, dict) or not plan.get("issues") or not plan.get("levels"):
        return None, "invalid_original_plan"
    if {x.get("name") for x in plan["issues"]} != {x.get("name") for x in checkpoint["all_issues"]}:
        return None, "plan_checkpoint_drift"
    inp = (detail.get("input_data") or {}).get("input") or {}
    if not inp or inp.get("build_id") != build_id or inp.get("repo_path", "").lower() != str(repo).lower():
        return None, "missing_original_execution_identity"
    if inp.get("plan_result", {}).get("issues") != plan.get("issues"):
        return None, "original_plan_not_equal"
    if any(x.get("status") in ("running", "queued", "pending", "waiting", "paused") for x in child_executions):
        return None, "live_or_ambiguous_child"
    if any(str(repo).lower() in c.lower() for c in process_commands):
        return None, "live_or_ambiguous_process"
    if checkpoint.get("in_flight_issues"):
        return None, "inflight_checkpoint_ambiguous"
    if git(repo, "status", "--porcelain"):
        return None, "dirty_integration_worktree"
    manifest = git(repo, "worktree", "list", "--porcelain")
    # Every attached issue worktree must be clean, registered, and on a real branch.
    for row in manifest.split("\n\n"):
        lines = row.splitlines()
        if not lines:
            continue
        path = lines[0].removeprefix("worktree ")
        if ".worktrees" not in path.lower():
            continue
        if not any(x.startswith("branch refs/heads/") for x in lines):
            return None, "detached_issue_worktree"
        if git(path, "status", "--porcelain"):
            return None, "dirty_issue_worktree"
    cfg = inp.get("git_config") or {}
    if not cfg.get("integration_branch") or git(repo, "branch", "--show-current") != cfg["integration_branch"]:
        return None, "integration_branch_mismatch"
    if not cfg.get("initial_commit_sha"):
        return None, "missing_base_sha"
    # Preserve full original plan and config; use native SWE-AF execute, not build().
    return {"input": {
        "plan_result": inp["plan_result"], "repo_path": repo,
        "config": inp.get("config"), "git_config": cfg,
        "resume": True, "build_id": build_id,
        "workspace_manifest": inp.get("workspace_manifest"),
    }}, None
