"""Independent disposable-Git oracle for the issue-level build lifecycle.

The coding loop is replaced only at its documented orchestration seam. All
repository, worktree, branch, commit, verifier, cleanup, and restart checks use
real local Git state; no remote or live AgentField task is invoked.
"""

from __future__ import annotations

import asyncio
import subprocess
from pathlib import Path
from types import SimpleNamespace

import pytest

from swe_af.execution.schemas import IssueOutcome
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


def _new_repo(root: Path, name: str) -> tuple[Path, str]:
    repo = root / name
    repo.mkdir()
    _git(repo, "init")
    _git(repo, "config", "user.name", "Luna High Oracle")
    _git(repo, "config", "user.email", "luna-oracle@example.invalid")
    (repo / "README.md").write_text("base\n", encoding="utf-8")
    _git(repo, "add", "README.md")
    _git(repo, "commit", "-m", "base")
    return repo, _git(repo, "rev-parse", "HEAD").stdout.strip()


def _status(repo: Path) -> str:
    return _git(repo, "status", "--porcelain").stdout


async def _run_issue(
    repo: Path,
    monkeypatch: pytest.MonkeyPatch,
    *,
    write_name: str,
    write_value: str,
    outcome: IssueOutcome = IssueOutcome.COMPLETED,
    raise_after_write: bool = False,
) -> tuple[dict, list[str]]:
    verifier_calls: list[str] = []

    async def fake_loop(*, dag_state, **_kwargs):
        target = Path(dag_state.repo_path) / write_name
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(write_value, encoding="utf-8")
        if raise_after_write:
            raise RuntimeError("simulated coder interruption")
        return SimpleNamespace(
            outcome=outcome,
            result_summary="deterministic oracle coding loop",
            error_message="",
            attempts=1,
            iteration_history=[],
            debt_items=[],
        )

    async def fake_call(*_args, **_kwargs):
        verifier_calls.append("run_verifier")
        return {"passed": True, "summary": "verified disposable artifact"}

    monkeypatch.setattr(build, "run_coding_loop", fake_loop)
    result = await build._implement_issue_impl(
        issue={
            "title": "Repair malicious ../../restart path",
            "description": "Create one persistent issue artifact.",
            "name": "../../restart oracle",
            "acceptance_criteria": ["one persistent branch artifact"],
        },
        repo_path=str(repo),
        base_branch="",
        artifacts_dir=".artifacts",
        additional_context="",
        config={"verify": True, "enable_github_pr": False},
        call_fn=fake_call,
        note_fn=None,
        node_id="luna-high-oracle",
    )
    return result, verifier_calls


def test_positive_edit_verify_persistent_branch_and_dirty_parent(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo, base_sha = _new_repo(tmp_path, "positive")
    dirty = repo / "parent-only.txt"
    dirty.write_text("must remain in caller checkout\n", encoding="utf-8")
    status_before = _status(repo)

    result, verifier_calls = asyncio.run(
        _run_issue(
            repo, monkeypatch, write_name="solution.txt", write_value="final repair\n"
        )
    )

    assert result["success"] is True
    assert len(result["commits"]) == 1
    assert verifier_calls == ["run_verifier"]
    assert result["branch"].startswith("issue/")
    assert ".." not in result["branch"]
    assert "/restart-oracle" in result["branch"]
    assert git_ops.new_commits(str(repo), base_sha, result["branch"]) == result["commits"]
    assert git_ops.changed_files(str(repo), base_sha, result["branch"]) == ["solution.txt"]
    assert _git(repo, "show", f"{result['commits'][0]}:solution.txt").stdout == "final repair\n"
    assert _git(repo, "remote").stdout == ""
    assert _status(repo) == status_before
    assert dirty.read_text(encoding="utf-8") == "must remain in caller checkout\n"


def test_interrupted_build_salvages_once_then_restart_has_no_duplicate_commit(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo, base_sha = _new_repo(tmp_path, "restart")

    interrupted, _ = asyncio.run(
        _run_issue(
            repo,
            monkeypatch,
            write_name="partial.txt",
            write_value="salvaged checkpoint\n",
            raise_after_write=True,
        )
    )

    assert interrupted["success"] is False
    assert len(interrupted["commits"]) == 1
    assert interrupted["verification"] is None
    assert _git(repo, "show", f"{interrupted['commits'][0]}:partial.txt").stdout == "salvaged checkpoint\n"

    restarted, verifier_calls = asyncio.run(
        _run_issue(
            repo, monkeypatch, write_name="solution.txt", write_value="final repair\n"
        )
    )

    assert restarted["success"] is True
    assert len(restarted["commits"]) == 1
    assert verifier_calls == ["run_verifier"]
    assert len(git_ops.new_commits(str(repo), base_sha, restarted["branch"])) == 1
    assert git_ops.new_commits(str(repo), base_sha, interrupted["branch"]) == interrupted["commits"]
    assert _git(repo, "show", f"{restarted['commits'][0]}:solution.txt").stdout == "final repair\n"


def test_non_completed_loop_is_fail_closed_and_does_not_verify(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo, _ = _new_repo(tmp_path, "rejected")

    result, verifier_calls = asyncio.run(
        _run_issue(
            repo,
            monkeypatch,
            write_name="rejected.txt",
            write_value="candidate\n",
            outcome=IssueOutcome.FAILED_UNRECOVERABLE,
        )
    )

    assert result["success"] is False
    assert result["outcome"] == IssueOutcome.FAILED_UNRECOVERABLE.value
    assert result["commits"]
    assert result["verification"] is None
    assert verifier_calls == []
    assert git_ops.new_commits(str(repo), result["base_sha"], result["branch"]) == result["commits"]


def test_commit_all_is_idempotent_after_accepted_artifact(tmp_path: Path) -> None:
    repo, base_sha = _new_repo(tmp_path, "idempotent")
    worktree = tmp_path / "worktree"
    branch = "issue/idempotent"
    git_ops.add_worktree(str(repo), str(worktree), branch, base_sha)
    try:
        (worktree / "solution.txt").write_text("one commit\n", encoding="utf-8")
        first = git_ops.commit_all(str(worktree), "fix: accepted artifact")
        second = git_ops.commit_all(str(worktree), "fix: duplicate attempt")
        assert first
        assert second == ""
        assert git_ops.new_commits(str(repo), base_sha, branch) == [first]
    finally:
        git_ops.remove_worktree(str(repo), str(worktree))
        git_ops.delete_branch(str(repo), branch)
