"""Independent native Git branch/worktree contract oracle.

These tests use disposable local repositories and replace only the documented
coding/verifier seam. They intentionally do not invoke remotes or live tasks.
"""

from __future__ import annotations

import asyncio
import re
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
    _git(repo, "config", "user.name", "Native Git Oracle")
    _git(repo, "config", "user.email", "native-oracle@example.invalid")
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
            result_summary="deterministic native Git contract loop",
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
        node_id="native-git-contract-oracle",
    )
    return result, verifier_calls


def test_safe_branch_name_edit_verify_and_dirty_parent(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo, base_sha = _new_repo(tmp_path, "positive")
    dirty = repo / "parent-only.txt"
    dirty.write_text("caller checkout state\n", encoding="utf-8")
    status_before = _status(repo)

    result, verifier_calls = asyncio.run(
        _run_issue(repo, monkeypatch, write_name="solution.txt", write_value="final repair\n")
    )

    assert result["success"] is True
    assert len(result["commits"]) == 1
    assert verifier_calls == ["run_verifier"]
    assert re.fullmatch(r"issue/[0-9a-f]+-restart-oracle", result["branch"])
    assert ".." not in result["branch"]
    assert result["branch"].endswith("-restart-oracle")
    assert git_ops.new_commits(str(repo), base_sha, result["branch"]) == result["commits"]
    assert git_ops.changed_files(str(repo), base_sha, result["branch"]) == ["solution.txt"]
    assert _git(repo, "show", f"{result['commits'][0]}:solution.txt").stdout == "final repair\n"
    assert _git(repo, "remote").stdout == ""
    assert _status(repo) == status_before
    assert dirty.read_text(encoding="utf-8") == "caller checkout state\n"


def test_branch_collision_is_rejected_without_duplicate_worktree_effect(
    tmp_path: Path,
) -> None:
    repo, base_sha = _new_repo(tmp_path, "collision")
    worktree = tmp_path / "worktree"
    branch = "issue/collision-safe"
    git_ops.add_worktree(str(repo), str(worktree), branch, base_sha)
    try:
        with pytest.raises(git_ops.GitOpsError, match="worktree add failed"):
            git_ops.add_worktree(str(repo), str(tmp_path / "duplicate"), branch, base_sha)
        assert _git(repo, "branch", "--list", branch).stdout.strip() == f"+ {branch}"
        assert not (tmp_path / "duplicate").exists()
        assert len(git_ops.new_commits(str(repo), base_sha, branch)) == 0
    finally:
        git_ops.remove_worktree(str(repo), str(worktree))
        git_ops.delete_branch(str(repo), branch)


def test_path_traversal_input_fails_closed_and_never_escapes_branch_namespace(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo, _ = _new_repo(tmp_path, "traversal")
    result, verifier_calls = asyncio.run(
        _run_issue(
            repo,
            monkeypatch,
            write_name="safe.txt",
            write_value="safe artifact\n",
            outcome=IssueOutcome.FAILED_UNRECOVERABLE,
        )
    )

    assert result["success"] is False
    assert result["outcome"] == IssueOutcome.FAILED_UNRECOVERABLE.value
    assert verifier_calls == []
    assert re.fullmatch(r"issue/[0-9a-f]+-restart-oracle", result["branch"])
    assert ".." not in result["branch"]
    assert not result["branch"].startswith("/")
    assert result["commits"]


def test_interruption_salvages_once_and_restart_has_one_commit(
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
        _run_issue(repo, monkeypatch, write_name="solution.txt", write_value="final repair\n")
    )
    assert restarted["success"] is True
    assert len(restarted["commits"]) == 1
    assert verifier_calls == ["run_verifier"]
    assert len(git_ops.new_commits(str(repo), base_sha, restarted["branch"])) == 1
    assert git_ops.new_commits(str(repo), base_sha, interrupted["branch"]) == interrupted["commits"]
    assert _git(repo, "show", f"{restarted['commits'][0]}:solution.txt").stdout == "final repair\n"
