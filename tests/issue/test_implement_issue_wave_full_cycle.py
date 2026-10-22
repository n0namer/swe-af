"""Deterministic full-cycle evidence for the issue-level build path."""

from __future__ import annotations

import asyncio
import os
import subprocess

from swe_af.issue.build import _implement_issue_impl
from tests.issue.conftest import run_git


ISSUE = {
    "title": "Add retry helper",
    "description": "Add a retry helper with bounded exponential backoff.",
    "acceptance_criteria": ["retry() retries a failing operation up to three times"],
    "files_to_create": ["retryable.py"],
}


def _run(repo_path: str, call_fn, *, config: dict | None = None) -> dict:
    return asyncio.run(
        _implement_issue_impl(
            issue=ISSUE,
            repo_path=repo_path,
            base_branch="main",
            artifacts_dir=".artifacts",
            additional_context="",
            config=config or {},
            call_fn=call_fn,
            note_fn=None,
            node_id="wave-test",
        )
    )


def _retry_source() -> str:
    return (
        "def retry(operation, attempts=3):\n"
        "    last_error = None\n"
        "    for _ in range(attempts):\n"
        "        try:\n"
        "            return operation()\n"
        "        except Exception as error:\n"
        "            last_error = error\n"
        "    raise last_error\n"
    )


def test_full_cycle_keeps_real_artifact_and_fail_closed_verifier(git_repo: str) -> None:
    calls: list[str] = []

    async def call_fn(target: str, **kwargs):
        calls.append(target)
        role = target.split(".", 1)[1]
        if role == "run_coder":
            path = os.path.join(kwargs["worktree_path"], "retryable.py")
            with open(path, "w", encoding="utf-8") as stream:
                stream.write(_retry_source())
            return {"files_changed": ["retryable.py"], "summary": "implemented", "complete": True}
        if role == "run_code_reviewer":
            return {"approved": True, "blocking": False, "summary": "accepted"}
        if role == "run_verifier":
            worktree = kwargs["repo_path"]
            probe = subprocess.run(
                [
                    "python",
                    "-c",
                    "from retryable import retry; state=[]; assert retry(lambda: state.append(1) or (42 if len(state) == 3 else (_ for _ in ()).throw(ValueError()))) == 42",
                ],
                cwd=worktree,
                capture_output=True,
                text=True,
            )
            return {
                "passed": probe.returncode == 0,
                "summary": probe.stderr or "oracle passed",
                "criteria_results": [{"criterion": ISSUE["acceptance_criteria"][0], "passed": probe.returncode == 0}],
            }
        raise AssertionError(f"unexpected call target: {target}")

    result = _run(repo_path=git_repo, call_fn=call_fn)

    assert result["success"] is True
    assert result["verification"]["passed"] is True
    assert result["commits"]
    assert result["files_changed"] == ["retryable.py"]
    assert run_git(git_repo, "show", f"{result['branch']}:retryable.py") == _retry_source().rstrip(chr(10))
    assert any(item.endswith(".run_verifier") for item in calls)

    failing_calls: list[str] = []

    async def failing_verifier(target: str, **kwargs):
        failing_calls.append(target)
        role = target.split(".", 1)[1]
        if role == "run_coder":
            path = os.path.join(kwargs["worktree_path"], "retryable.py")
            with open(path, "w", encoding="utf-8") as stream:
                stream.write("def retry(operation):\n    return None\n")
            return {"files_changed": ["retryable.py"], "summary": "bad implementation", "complete": True}
        if role == "run_code_reviewer":
            return {"approved": True, "blocking": False, "summary": "review accepted"}
        if role == "run_verifier":
            return {"passed": False, "summary": "oracle rejected behavior", "criteria_results": []}
        raise AssertionError(f"unexpected call target: {target}")

    failed = _run(repo_path=git_repo, call_fn=failing_verifier)
    assert failed["success"] is False
    assert failed["verification"]["passed"] is False
    assert failed["branch"]
    assert run_git(git_repo, "show", f"{failed['branch']}:retryable.py")
    assert any(item.endswith(".run_verifier") for item in failing_calls)


def test_interruption_checkpoints_uncommitted_work_for_recovery(git_repo: str) -> None:
    async def interrupted(target: str, **kwargs):
        if target.endswith("run_coder"):
            path = os.path.join(kwargs["worktree_path"], "retryable.py")
            with open(path, "w", encoding="utf-8") as stream:
                stream.write(_retry_source())
            raise RuntimeError("simulated interruption")
        raise AssertionError(f"unexpected call target: {target}")

    result = _run(repo_path=git_repo, call_fn=interrupted, config={"verify": False})

    assert result["success"] is False
    assert result["outcome"] == "failed_unrecoverable"
    assert result["commits"]
    assert result["branch"]
    assert "retryable.py" in result["files_changed"]
    assert run_git(git_repo, "show", f"{result['branch']}:retryable.py") == _retry_source().rstrip(chr(10))
