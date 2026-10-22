"""One-shot native SWE-AF plan -> execute handoff, not a scheduler.

Default is read-only. --cutover moves a single verified finished plan from the
old Full build to swe-factory-30.execute, without re-running planning.
"""
import argparse
import json
import subprocess
import sys
import urllib.parse
import urllib.request
from pathlib import Path

BASE = "http://127.0.0.1:18080"
OPS = Path(__file__).resolve().parent
CFG = json.loads((OPS / "project-profiles.json").read_text(encoding="utf-8"))
ROOTS = {
    "pr-af": ("run_20261020_020554_t6av71xu",
              "D:/Temp/swe-full-pr-af-20261020"),
    "af-skill-intelligence": ("run_20261020_020554_4emb1fv5",
                              "D:/Temp/swe-full-skill-intelligence-20261020"),
    "codeaf": ("run_20261020_020554_7t8xp8zl",
               "D:/Temp/swe-full-codeaf-20261020"),
}


def api(path, payload=None):
    data = json.dumps(payload).encode() if payload is not None else None
    with urllib.request.urlopen(urllib.request.Request(
            BASE + path, data=data,
            headers={"Content-Type": "application/json"} if data else {},
            method="POST" if data else "GET"), timeout=20) as r:
        return json.load(r)


def git(repo, *args):
    return subprocess.check_output(["git", "-C", repo, *args], text=True).strip()


def main():
    arg = argparse.ArgumentParser()
    arg.add_argument("project", choices=sorted(ROOTS))
    arg.add_argument("--cutover", action="store_true")
    a = arg.parse_args()
    run_id, repo = ROOTS[a.project]
    if not Path(repo).is_dir():
        raise RuntimeError(f"missing repo: {repo}")
    q = urllib.parse.quote(run_id)
    executions = api("/api/ui/v1/agents/swe-factory/executions?pageSize=100&workflowId=" + q).get("executions", [])
    plan = [e for e in executions if e.get("reasoner_id") == "plan"]
    old_root = [e for e in executions if e.get("reasoner_id") == "build"]
    old_exec = [e for e in executions if e.get("reasoner_id") == "execute"]
    if len(plan) != 1 or len(old_root) != 1:
        raise RuntimeError("Cannot identify exactly one original plan and build")
    state = {
        "project": a.project, "source_run": run_id,
        "plan_status": plan[0]["status"],
        "old_build_status": old_root[0]["status"],
        "old_execute_count": len(old_exec),
        "new_target": "swe-factory-30.execute",
        "max_concurrent_issues": CFG["full_dag_issue_limits"][a.project],
        "cutover_requested": a.cutover,
    }
    if plan[0]["status"] != "succeeded":
        print(json.dumps({**state, "decision": "NOT_READY_KEEP_PLANNING"}))
        return
    if old_exec:
        print(json.dumps({**state, "decision": "REFUSE_EXISTING_EXECUTION"}))
        return
    result = api("/api/ui/v1/executions/" + plan[0]["execution_id"] + "/details").get("output_data")
    if not isinstance(result, dict) or not result.get("issues") or not result.get("levels"):
        raise RuntimeError("Incomplete plan_result, refusing cutoff")
    if len({x["name"] for x in result["issues"]}) != len(result["issues"]):
        raise RuntimeError("duplicate issue names")
    original_sha = git(repo, "rev-parse", "HEAD")
    branch = git(repo, "symbolic-ref", "--short", "HEAD")
    tracked_dirty = git(repo, "status", "--porcelain", "--untracked-files=no")
    if tracked_dirty:
        raise RuntimeError("tracked files have uncommitted changes; preserve and reconcile first")
    source_root = str(CFG["swe_af_code_root"])
    sys.path.insert(0, source_root)
    from swe_af.reasoners.pipeline import (
        _compute_levels, _serialize_file_conflicts, _validate_file_conflicts,
    )
    from swe_af.reasoners.schemas import PlanResult
    owned = json.loads(json.dumps(result))
    extra_edges = _serialize_file_conflicts(owned["issues"])
    owned["levels"] = _compute_levels(owned["issues"])
    owned["file_conflicts"] = _validate_file_conflicts(owned["issues"], owned["levels"])
    if owned["file_conflicts"]:
        raise RuntimeError("unresolved parallel file conflict")
    PlanResult.model_validate(owned)
    state.update({
        "issue_count": len(owned["issues"]),
        "ready_front_width": len(owned["levels"][0]),
        "levels_count": len(owned["levels"]),
        "serialized_conflict_edges": len(extra_edges),
        "branch": branch,
        "sha": original_sha,
    })
    if not a.cutover:
        print(json.dumps({**state, "decision": "READY_FOR_EXPLICIT_CUTOVER"}))
        return
    active = api("/api/v1/executions/active?limit=100").get("runs", [])
    if not any(run.get("run_id") == run_id for run in active):
        raise RuntimeError("Old build is not active; inspect terminal result before a new execute")
    for run in active:
        if run.get("run_id") != run_id and str(run.get("target", "")).endswith(".execute"):
            raise RuntimeError("Another DAG execute is active; refuse potential duplicate")
    cancelled = api("/api/v1/workflows/" + q + "/cancel-tree",
                    {"reason": "Switching verified native plan to 30-slot execute"})
    if cancelled.get("error_count"):
        raise RuntimeError(f"cancel-tree errors, do not dispatch: {cancelled['error_count']}")
    recheck = api("/api/ui/v1/agents/swe-factory/executions?pageSize=100&workflowId=" + q).get("executions", [])
    if any(e.get("reasoner_id") == "execute" for e in recheck):
        raise RuntimeError("Old run had already started execute; cutover stopped fail-closed")
    if any(e.get("status") in ("running", "queued", "waiting")
           for e in recheck):
        raise RuntimeError("Old run still active after cancel-tree")
    models = {
        "default": "gpt-6-luna", "coder": "gpt-6-luna",
        "qa": "gpt-6.1-sol", "code_reviewer": "gpt-6.1-sol",
        "qa_synthesizer": "gpt-6-luna", "issue_advisor": "gpt-6.1-sol",
        "replan": "gpt-6.1-sol", "retry_advisor": "gpt-6.1-sol",
        "verifier": "gpt-6.1-sol", "git": "gpt-6-luna",
        "merger": "gpt-6.1-sol", "integration_tester": "gpt-6.1-sol",
    }
    config = {
        "runtime": "codex", "models": models,
        "max_concurrent_issues": CFG["full_dag_issue_limits"][a.project],
        "max_coding_iterations": 6, "agent_max_turns": 40,
        "agent_timeout_seconds": 3600, "max_retries_per_issue": 2,
        "max_advisor_invocations": 2, "enable_issue_advisor": True,
        "enable_replanning": True, "max_replans": 3,
        "deterministic_git": True, "enable_integration_testing": True,
        "check_ci": False, "enable_learning": False,
    }
    git_cfg = {
        "integration_branch": branch, "original_branch": branch,
        "initial_commit_sha": original_sha, "mode": "existing",
        "remote_url": "", "remote_default_branch": "",
    }
    dispatch = api("/api/v1/execute/async/swe-factory-30.execute", {
        "input": {
            "plan_result": owned, "repo_path": repo,
            "config": config, "git_config": git_cfg,
            "build_id": f"handoff-{a.project}-20261020",
            "resume": False,
        }
    })
    print(json.dumps({**state, "decision": "DISPATCHED_NATIVE_EXECUTE",
                      "new_execution": dispatch}, ensure_ascii=False))


if __name__ == "__main__":
    main()
