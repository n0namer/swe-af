"""C4 end-to-end evidence for the deterministic issue build lifecycle.

The repository and Git operations are real. Only the LLM coding-loop boundary is
controlled so the test can deterministically cover success, rejection, verifier
failure, and interruption/restart without network or provider dependencies.
"""

from __future__ import annotations

import asyncio
import inspect
import subprocess
from pathlib import Path
from types import SimpleNamespace

import pytest

from swe_af.execution import schemas as execution_schemas
from swe_af.issue import build, git_ops


def _git(repo: Path, *args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    result = subprocess.run(
        ["git", "-C", str(repo), *args],
        capture_output=True,
        text=True,
        check=False,
    )
    if check and result.returncode != 0:
        raise AssertionError(result.stderr or result.stdout)
    return result


def _new_repo(root: Path) -> tuple[Path, str]:
    repo = root / "repo"
    repo.mkdir()
    _git(repo, "init")
    _git(repo, "config", "user.name", "C4 Acceptance")
    _git(repo, "config", "user.email", "c4@example.invalid")
    (repo / "README.md").write_text("base\n", encoding="utf-8")
    _git(repo, "add", "README.md")
    _git(repo, "commit", "-m", "base")
    return repo, _git(repo, "rev-parse", "HEAD").stdout.strip()


def _issue() -> dict:
    return {
        "name": "northstar-repair",
        "title": "Repair the issue",
        "description": "Apply the bounded repair and preserve the caller checkout.",
        "acceptance_criteria": ["The isolated branch contains the repair."],
    }


def _result(outcome: object, summary: str = "loop result") -> SimpleNamespace:
    return SimpleNamespace(
        outcome=outcome,
        result_summary=summary,
        error_message="",
        attempts=1,
        iteration_history=[{"attempt": 1, "summary": summary}],
        debt_items=[],
    )


def _run_impl(
    monkeypatch: pytest.MonkeyPatch,
    repo: Path,
    loop_factory,
    *,
    verify: bool = False,
    verifier_passed: bool = True,
) -> dict:
    async def fake_call(_target: str, **_kwargs):
        return {"passed": verifier_passed, "summary": "deterministic verifier"}

    monkeypatch.setattr(build, "run_coding_loop", loop_factory)
    return asyncio.run(
        build._implement_issue_impl(
            issue=_issue(),
            repo_path=str(repo),
            base_branch="",
            artifacts_dir=".artifacts",
            additional_context="",
            config={
                "verify": verify,
                "keep_worktree": False,
                "enable_github_pr": False,
            },
            call_fn=fake_call,
            note_fn=None,
            node_id="test",
        )
    )


def _cleanup_branch(repo: Path, result: dict) -> None:
    if result["branch"]:
        git_ops.delete_branch(str(repo), result["branch"])


def test_signature_dto_and_positive_cycle_preserve_immutable_base(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo, base_sha = _new_repo(tmp_path)
    dirty_file = repo / "caller-dirty.txt"
    dirty_file.write_text("preserve\n", encoding="utf-8")
    status_before = _git(repo, "status", "--porcelain").stdout

    assert list(inspect.signature(build.implement_issue).parameters) == [
        "issue",
        "repo_path",
        "base_branch",
        "artifacts_dir",
        "additional_context",
        "config",
    ]

    async def completed_loop(*, dag_state, **_kwargs):
        (Path(dag_state.repo_path) / "solution.txt").write_text(
            "repaired\n", encoding="utf-8"
        )
        return _result(execution_schemas.IssueOutcome.COMPLETED, "repair completed")

    result = _run_impl(monkeypatch, repo, completed_loop, verify=True)
    try:
        assert result["success"] is True
        assert result["base_sha"] == base_sha
        assert len(result["commits"]) == 1
        assert result["files_changed"] == ["solution.txt"]
        assert result["verification"] == {
            "passed": True,
            "summary": "deterministic verifier",
        }
        assert _git(repo, "show", f"{result['commits'][0]}:solution.txt").stdout == "repaired\n"
        assert _git(repo, "rev-parse", result["branch"]).stdout.strip() != base_sha
        assert _git(repo, "status", "--porcelain").stdout == status_before
        assert dirty_file.read_text(encoding="utf-8") == "preserve\n"
    finally:
        _cleanup_branch(repo, result)


def test_rejected_loop_returns_failure_without_duplicate_effects(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo, base_sha = _new_repo(tmp_path)
    status_before = _git(repo, "status", "--porcelain").stdout

    rejected_outcome = next(
        outcome
        for outcome in execution_schemas.IssueOutcome
        if outcome.value not in {"completed", "completed_with_debt"}
    )

    async def rejected_loop(*_args, **_kwargs):
        return _result(rejected_outcome, "acceptance rejected")

    result = _run_impl(monkeypatch, repo, rejected_loop)
    assert result["success"] is False
    assert result["outcome"] == rejected_outcome.value
    assert result["commits"] == []
    assert result["branch"] == ""
    assert _git(repo, "rev-parse", "HEAD").stdout.strip() == base_sha
    assert _git(repo, "status", "--porcelain").stdout == status_before


def test_verifier_rejection_keeps_receipt_but_never_reports_success(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo, base_sha = _new_repo(tmp_path)

    async def completed_loop(*, dag_state, **_kwargs):
        (Path(dag_state.repo_path) / "candidate.txt").write_text(
            "candidate\n", encoding="utf-8"
        )
        return _result(execution_schemas.IssueOutcome.COMPLETED)

    result = _run_impl(
        monkeypatch,
        repo,
        completed_loop,
        verify=True,
        verifier_passed=False,
    )
    try:
        assert result["success"] is False
        assert result["base_sha"] == base_sha
        assert len(result["commits"]) == 1
        assert result["verification"]["passed"] is False
        assert _git(repo, "show", f"{result['commits'][0]}:candidate.txt").stdout == "candidate\n"
    finally:
        _cleanup_branch(repo, result)


def test_crash_salvage_then_restart_has_distinct_single_commit_receipts(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo, base_sha = _new_repo(tmp_path)

    async def crashing_loop(*, dag_state, **_kwargs):
        (Path(dag_state.repo_path) / "partial.txt").write_text(
            "interrupted\n", encoding="utf-8"
        )
        raise RuntimeError("simulated coder interruption")

    crashed = _run_impl(monkeypatch, repo, crashing_loop)
    try:
        assert crashed["success"] is False
        assert "simulated coder interruption" in crashed["error_message"]
        assert len(crashed["commits"]) == 1
        assert crashed["files_changed"] == ["partial.txt"]
        assert _git(repo, "show", f"{crashed['commits'][0]}:partial.txt").stdout == "interrupted\n"
    finally:
        _cleanup_branch(repo, crashed)

    async def restarted_loop(*, dag_state, **_kwargs):
        (Path(dag_state.repo_path) / "solution.txt").write_text(
            "final\n", encoding="utf-8"
        )
        return _result(execution_schemas.IssueOutcome.COMPLETED, "restart completed")

    restarted = _run_impl(monkeypatch, repo, restarted_loop)
    try:
        assert restarted["success"] is True
        assert restarted["base_sha"] == base_sha
        assert len(restarted["commits"]) == 1
        assert restarted["branch"] != crashed["branch"]
        assert restarted["files_changed"] == ["solution.txt"]
        assert _git(repo, "show", f"{restarted['commits'][0]}:solution.txt").stdout == "final\n"
        assert _git(repo, "show", f"{restarted['commits'][0]}:partial.txt", check=False).returncode != 0
        assert _git(repo, "rev-parse", "HEAD").stdout.strip() == base_sha
    finally:
        _cleanup_branch(repo, restarted)
