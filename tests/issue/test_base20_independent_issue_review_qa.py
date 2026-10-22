"""Independent QA contracts for the Base20 issue-build and BMAD-review ladders.

The tests use real disposable Git repositories and the real issue entry-point,
while replacing only the model/AgentField boundary.  Their verdicts are local
review evidence: neither a process exit nor a reviewer response is treated as
product acceptance.
"""

from __future__ import annotations

import asyncio
import os
import subprocess
from dataclasses import dataclass
from pathlib import Path

import pytest

from swe_af.issue.build import _implement_issue_impl


@dataclass(frozen=True)
class ReviewReceipt:
    verdict: str
    base_sha: str
    current_sha: str
    branch: str
    branch_sha: str
    oracle_passed: bool
    reason: str


def _git(repo_path: str, *args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        ["git", "-C", repo_path, *args],
        check=check,
        capture_output=True,
        text=True,
    )


def _repo(tmp_path: Path) -> str:
    repo_path = tmp_path / "repository"
    repo_path.mkdir()
    _git(str(repo_path), "init", "-q")
    _git(str(repo_path), "checkout", "-q", "-b", "main")
    _git(str(repo_path), "config", "user.name", "QA Test User")
    _git(str(repo_path), "config", "user.email", "qa@example.invalid")
    (repo_path / "README.md").write_text("base\n", encoding="utf-8")
    _git(str(repo_path), "add", "README.md")
    _git(str(repo_path), "commit", "-q", "-m", "base")
    return str(repo_path)


def _run(repo_path: str, issue: dict, call_fn, *, config: dict | None = None) -> dict:
    return asyncio.run(
        _implement_issue_impl(
            issue=issue,
            repo_path=repo_path,
            base_branch="main",
            artifacts_dir=".artifacts",
            additional_context="",
            config=config or {"max_coding_iterations": 2},
            call_fn=call_fn,
            note_fn=None,
            node_id="base20-independent-qa",
        )
    )


def _receipt(repo_path: str, result: dict) -> ReviewReceipt:
    base_sha = result.get("base_sha", "")
    branch = result.get("branch", "")
    current_sha = _git(repo_path, "rev-parse", "main").stdout.strip()
    branch_sha = _git(repo_path, "rev-parse", branch).stdout.strip() if branch else ""
    verification = result.get("verification")
    oracle_passed = isinstance(verification, dict) and verification.get("passed") is True
    if not result.get("success"):
        reason = "implementation result is not successful"
    elif not oracle_passed:
        reason = "independent oracle did not pass"
    elif base_sha != current_sha:
        reason = "source branch advanced after the frozen base snapshot"
    elif not branch or not branch_sha:
        reason = "deliverable branch receipt is missing"
    else:
        reason = "frozen source and executable oracle are consistent"
    return ReviewReceipt(
        verdict="ACCEPT_SCOPED" if reason.startswith("frozen") else "RED",
        base_sha=base_sha,
        current_sha=current_sha,
        branch=branch,
        branch_sha=branch_sha,
        oracle_passed=oracle_passed,
        reason=reason,
    )


def _repairing_call(filename: str, expected: int, calls: list[str]):
    attempts = 0

    async def call_fn(target: str, **kwargs):
        nonlocal attempts
        role = target.rsplit(".", 1)[1]
        calls.append(role)
        if role == "run_coder":
            attempts += 1
            value = 0 if attempts == 1 else expected
            Path(kwargs["worktree_path"], filename).write_text(
                f"def answer():\n    return {value}\n", encoding="utf-8"
            )
            return {"files_changed": [filename], "summary": f"attempt {attempts}", "complete": True}
        if role == "run_code_reviewer":
            return {
                "approved": attempts > 1,
                "blocking": False,
                "summary": "deliberate RED" if attempts == 1 else "repair reviewed",
            }
        if role == "run_verifier":
            probe = subprocess.run(
                ["python", "-c", f"from {filename[:-3]} import answer; assert answer() == {expected}"],
                cwd=kwargs["repo_path"],
                capture_output=True,
                text=True,
            )
            return {
                "passed": probe.returncode == 0,
                "summary": probe.stderr or "oracle passed",
                "criteria_results": [{"criterion": "answer repaired", "passed": probe.returncode == 0}],
            }
        raise AssertionError(f"unexpected target: {target}")

    return call_fn


@pytest.mark.parametrize(("filename", "expected"), [("alpha_qa.py", 7), ("beta_qa.py", 11)])
def test_two_safe_tasks_freeze_evidence_after_same_task_repair(
    tmp_path: Path, filename: str, expected: int
) -> None:
    repo_path = _repo(tmp_path)
    calls: list[str] = []
    result = _run(
        repo_path,
        {
            "title": f"Implement {filename}",
            "description": "Write a deterministic answer function.",
            "acceptance_criteria": [f"{filename} returns its expected value"],
            "files_to_create": [filename],
        },
        _repairing_call(filename, expected, calls),
    )

    receipt = _receipt(repo_path, result)

    assert receipt.verdict == "ACCEPT_SCOPED"
    assert result["verification"]["passed"] is True
    assert calls == [
        "run_coder",
        "run_code_reviewer",
        "run_coder",
        "run_code_reviewer",
        "run_verifier",
    ]
    assert _git(repo_path, "show", f"{receipt.branch}:{filename}").stdout == (
        f"def answer():\n    return {expected}\n"
    )
    assert _git(repo_path, "status", "--porcelain").stdout == ""


def test_verifier_false_result_is_red_even_when_reviewer_approves(tmp_path: Path) -> None:
    repo_path = _repo(tmp_path)

    async def false_oracle(target: str, **kwargs):
        role = target.rsplit(".", 1)[1]
        if role == "run_coder":
            Path(kwargs["worktree_path"], "wrong.py").write_text("value = 0\n", encoding="utf-8")
            return {"files_changed": ["wrong.py"], "summary": "incorrect artifact", "complete": True}
        if role == "run_code_reviewer":
            return {"approved": True, "blocking": False, "summary": "misleading green review"}
        if role == "run_verifier":
            return {"passed": False, "summary": "oracle rejects value", "criteria_results": []}
        raise AssertionError(f"unexpected target: {target}")

    result = _run(
        repo_path,
        {"title": "Reject false green", "description": "Verifier is authoritative.", "acceptance_criteria": ["oracle passes"]},
        false_oracle,
    )

    receipt = _receipt(repo_path, result)

    assert result["success"] is False
    assert receipt.verdict == "RED"
    assert receipt.oracle_passed is False
    assert receipt.reason == "implementation result is not successful"


def test_moving_main_branch_is_red_despite_a_completed_isolated_branch(tmp_path: Path) -> None:
    repo_path = _repo(tmp_path)

    async def move_main_after_coding(target: str, **kwargs):
        role = target.rsplit(".", 1)[1]
        if role == "run_coder":
            Path(kwargs["worktree_path"], "artifact.py").write_text("value = 5\n", encoding="utf-8")
            Path(repo_path, "README.md").write_text("advanced\n", encoding="utf-8")
            _git(repo_path, "add", "README.md")
            _git(repo_path, "commit", "-q", "-m", "foreign advance")
            return {"files_changed": ["artifact.py"], "summary": "artifact created", "complete": True}
        if role == "run_code_reviewer":
            return {"approved": True, "blocking": False, "summary": "reviewed"}
        if role == "run_verifier":
            return {"passed": True, "summary": "local oracle passed", "criteria_results": []}
        raise AssertionError(f"unexpected target: {target}")

    result = _run(
        repo_path,
        {"title": "Freeze source", "description": "Do not accept a stale base.", "acceptance_criteria": ["frozen source"]},
        move_main_after_coding,
    )

    receipt = _receipt(repo_path, result)

    assert result["success"] is True
    assert receipt.base_sha != receipt.current_sha
    assert receipt.verdict == "RED"
    assert receipt.reason == "source branch advanced after the frozen base snapshot"
    assert _git(repo_path, "show", f"{receipt.branch}:artifact.py").stdout == "value = 5\n"


def test_interrupted_uncommitted_work_is_checkpointed_as_red_evidence(tmp_path: Path) -> None:
    repo_path = _repo(tmp_path)

    async def interrupted(target: str, **kwargs):
        if target.rsplit(".", 1)[1] == "run_coder":
            Path(kwargs["worktree_path"], "interrupted.py").write_text(
                "value = 'preserved'\n", encoding="utf-8"
            )
            raise RuntimeError("simulated interruption")
        raise AssertionError(f"unexpected target: {target}")

    result = _run(
        repo_path,
        {"title": "Checkpoint interrupted work", "description": "Preserve partial evidence.", "acceptance_criteria": ["artifact survives interruption"]},
        interrupted,
        config={"verify": False},
    )

    receipt = _receipt(repo_path, result)

    assert result["success"] is False
    assert receipt.verdict == "RED"
    assert result["commits"]
    assert _git(repo_path, "show", f"{result['branch']}:interrupted.py").stdout == "value = 'preserved'\n"


def test_interruption_and_exhausted_budget_do_not_duplicate_effects(tmp_path: Path) -> None:
    repo_path = _repo(tmp_path)
    calls: list[str] = []

    async def interrupted_then_rejected(target: str, **kwargs):
        role = target.rsplit(".", 1)[1]
        calls.append(role)
        if role == "run_coder":
            Path(kwargs["worktree_path"], "partial.py").write_text("value = 'checkpoint'\n", encoding="utf-8")
            return {"files_changed": ["partial.py"], "summary": "one incomplete attempt", "complete": True}
        if role == "run_code_reviewer":
            return {"approved": False, "blocking": True, "summary": "budget exhausted; repair required"}
        raise AssertionError(f"unexpected target: {target}")

    result = _run(
        repo_path,
        {"title": "Preserve only one partial effect", "description": "Reject without retry.", "acceptance_criteria": ["no duplicate commits"]},
        interrupted_then_rejected,
        config={"max_coding_iterations": 1, "verify": False},
    )

    receipt = _receipt(repo_path, result)

    assert result["success"] is False
    assert receipt.verdict == "RED"
    assert calls == ["run_coder", "run_code_reviewer"]
    assert result["commits"]
    assert len(result["commits"]) == len(set(result["commits"]))
    assert _git(repo_path, "show", f"{result['branch']}:partial.py").stdout == "value = 'checkpoint'\n"
