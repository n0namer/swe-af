"""Disposable native issue repair checks for the C4 restart boundary."""

from __future__ import annotations

import asyncio
from pathlib import Path
import subprocess

import pytest

from swe_af.issue import git_ops
from swe_af.issue.build import _checkpoint_interrupted_issue


def _git(repo_path: str, *args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        ["git", "-C", repo_path, *args],
        check=check,
        capture_output=True,
        text=True,
    )


def _git_text(repo_path: str, *args: str) -> str:
    return _git(repo_path, *args).stdout.strip()


@pytest.fixture
def disposable_repo(tmp_path: Path) -> str:
    repo_path = tmp_path / "repo"
    repo_path.mkdir()
    repo = str(repo_path)
    _git(repo, "init", "-q")
    _git(repo, "checkout", "-q", "-b", "main")
    _git(repo, "config", "user.name", "MVP QA")
    _git(repo, "config", "user.email", "mvp-qa@example.invalid")
    (repo_path / "README.md").write_text("# disposable issue\n", encoding="utf-8")
    _git(repo, "add", "README.md")
    initial = _git(repo, "commit", "-q", "-m", "initial", check=False)
    assert initial.returncode == 0, initial.stderr
    return repo


def test_positive_repair_persists_one_typed_git_effect(disposable_repo: str, tmp_path: Path) -> None:
    base_sha = _git_text(disposable_repo, "rev-parse", "HEAD")
    worktree = str(tmp_path / "positive")

    git_ops.add_worktree(disposable_repo, worktree, "issue/positive", base_sha)
    Path(worktree, "repair.json").write_text('{"revision": 1}\n', encoding="utf-8")
    commit_sha = git_ops.commit_all(worktree, "fix: persist repair")

    assert len(commit_sha) == 40
    assert git_ops.new_commits(disposable_repo, base_sha, "issue/positive") == [commit_sha]
    assert git_ops.changed_files(disposable_repo, base_sha, "issue/positive") == ["repair.json"]
    assert _git_text(disposable_repo, "rev-parse", "HEAD") == base_sha
    git_ops.remove_worktree(disposable_repo, worktree)


def test_wrong_identity_branch_collision_fails_closed(disposable_repo: str, tmp_path: Path) -> None:
    base_sha = _git_text(disposable_repo, "rev-parse", "HEAD")
    first = str(tmp_path / "first")
    second = str(tmp_path / "second")
    git_ops.add_worktree(disposable_repo, first, "issue/collision", base_sha)

    with pytest.raises(git_ops.GitOpsError):
        git_ops.add_worktree(disposable_repo, second, "issue/collision", base_sha)

    assert git_ops.new_commits(disposable_repo, base_sha, "issue/collision") == []
    assert _git_text(disposable_repo, "rev-parse", "HEAD") == base_sha
    git_ops.remove_worktree(disposable_repo, first)


def test_restart_replays_checkpoint_without_duplicate_effect(disposable_repo: str, tmp_path: Path) -> None:
    base_sha = _git_text(disposable_repo, "rev-parse", "HEAD")
    interrupted = str(tmp_path / "interrupted")
    restarted = str(tmp_path / "restarted")
    branch = "issue/restart"

    git_ops.add_worktree(disposable_repo, interrupted, branch, base_sha)
    Path(interrupted, "repair.json").write_text('{"revision": 1}\n', encoding="utf-8")
    checkpoint_sha = git_ops.commit_all(interrupted, "chore: checkpoint interrupted repair")
    git_ops.remove_worktree(disposable_repo, interrupted)

    git_ops.add_worktree(disposable_repo, restarted, f"{branch}-replay", checkpoint_sha)
    assert git_ops.commit_all(restarted, "chore: replay unchanged repair") == ""
    assert git_ops.new_commits(disposable_repo, base_sha, branch) == [checkpoint_sha]
    git_ops.remove_worktree(disposable_repo, restarted)
    assert _git_text(disposable_repo, "rev-parse", "HEAD") == base_sha


@pytest.mark.asyncio
async def test_cancellation_checkpoint_is_durable_and_does_not_swallow_cancel(
    disposable_repo: str, tmp_path: Path
) -> None:
    base_sha = _git_text(disposable_repo, "rev-parse", "HEAD")
    worktree = str(tmp_path / "cancelled")
    git_ops.add_worktree(disposable_repo, worktree, "issue/cancelled", base_sha)
    Path(worktree, "repair.json").write_text('{"revision": 1}\n', encoding="utf-8")
    notes: list[str] = []

    await _checkpoint_interrupted_issue(
        worktree_path=worktree,
        issue_name="cancelled-repair",
        note=lambda message, tags=None: notes.append(message),
    )

    commit_list = git_ops.new_commits(disposable_repo, base_sha, "issue/cancelled")
    assert len(commit_list) == 1
    assert notes == []
    with pytest.raises(asyncio.CancelledError):
        raise asyncio.CancelledError
    git_ops.remove_worktree(disposable_repo, worktree)
