"""Bounded native Git lifecycle checks for the issue-level helpers.

These tests exercise real disposable repositories and keep the caller's main
branch isolated from issue worktree commits.  They deliberately cover the
positive path, a deterministic branch collision, and restart/replay after an
interrupted repair checkpoint.
"""

from __future__ import annotations

import hashlib
from pathlib import Path
import subprocess

import pytest

from swe_af.issue import git_ops


SOURCE_PATH = Path(__file__).parents[2] / "swe_af" / "issue" / "git_ops.py"
SOURCE_SHA256 = "0f73d9fb3a613f1b68ccf4a8161a93c1c5e6dd6a4a392b2c2c8a6157b7c2403d"


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
    _git(repo, "config", "user.name", "Basefive QA")
    _git(repo, "config", "user.email", "basefive-qa@example.invalid")
    (repo_path / "README.md").write_text("# basefive\n", encoding="utf-8")
    _git(repo, "add", "README.md")
    initial = _git(repo, "commit", "-q", "-m", "initial", check=False)
    assert initial.returncode == 0, initial.stderr
    return repo


def test_positive_issue_commit_has_one_effect_and_preserves_caller(
    disposable_repo: str,
    tmp_path: Path,
) -> None:
    base_sha = _git_text(disposable_repo, "rev-parse", "HEAD")
    worktree = str(tmp_path / "positive-worktree")

    git_ops.add_worktree(disposable_repo, worktree, "issue/positive", base_sha)
    Path(worktree, "result.txt").write_text("native-pass\n", encoding="utf-8")
    commit_sha = git_ops.commit_all(worktree, "qa: record native result")

    assert len(commit_sha) == 40
    assert git_ops.new_commits(disposable_repo, base_sha, "issue/positive") == [commit_sha]
    assert git_ops.changed_files(disposable_repo, base_sha, "issue/positive") == ["result.txt"]
    assert _git_text(disposable_repo, "rev-parse", "HEAD") == base_sha
    assert _git_text(disposable_repo, "rev-parse", "--abbrev-ref", "HEAD") == "main"
    assert git_ops.remote_url(disposable_repo) == ""
    assert _git_text(disposable_repo, "status", "--porcelain") == ""

    git_ops.remove_worktree(disposable_repo, worktree)
    assert not Path(worktree).exists()


def test_branch_collision_fails_closed_without_second_effect(disposable_repo: str) -> None:
    base_sha = _git_text(disposable_repo, "rev-parse", "HEAD")
    first_worktree = str(Path(disposable_repo) / ".worktrees" / "first")
    second_worktree = str(Path(disposable_repo) / ".worktrees" / "second")

    git_ops.add_worktree(disposable_repo, first_worktree, "issue/collision", base_sha)
    collision = None
    try:
        git_ops.add_worktree(disposable_repo, second_worktree, "issue/collision", base_sha)
    except git_ops.GitOpsError as error:
        collision = error

    assert collision is not None
    message = str(collision).lower()
    assert "reference already exists" in message or "branch named" in message
    assert git_ops.new_commits(disposable_repo, base_sha, "issue/collision") == []
    assert _git_text(disposable_repo, "rev-parse", "HEAD") == base_sha

    git_ops.remove_worktree(disposable_repo, first_worktree)
    assert not Path(second_worktree).exists()


def test_restart_replays_checkpoint_without_duplicate_effect(disposable_repo: str) -> None:
    base_sha = _git_text(disposable_repo, "rev-parse", "HEAD")
    interrupted_worktree = str(Path(disposable_repo) / ".worktrees" / "interrupted")
    restarted_worktree = str(Path(disposable_repo) / ".worktrees" / "restarted")
    branch = "issue/restart"

    git_ops.add_worktree(disposable_repo, interrupted_worktree, branch, base_sha)
    Path(interrupted_worktree, "repair.txt").write_text("recovered\n", encoding="utf-8")
    checkpoint_sha = git_ops.commit_all(interrupted_worktree, "qa: checkpoint repair")
    assert len(checkpoint_sha) == 40
    git_ops.remove_worktree(disposable_repo, interrupted_worktree)

    git_ops.add_worktree(disposable_repo, restarted_worktree, f"{branch}-replay", checkpoint_sha)
    assert Path(restarted_worktree, "repair.txt").read_text(encoding="utf-8") == "recovered\n"
    assert git_ops.commit_all(restarted_worktree, "qa: replay unchanged repair") == ""
    assert git_ops.new_commits(disposable_repo, base_sha, branch) == [checkpoint_sha]
    assert _git_text(disposable_repo, "rev-parse", f"{branch}:repair.txt") == _git_text(
        disposable_repo, "rev-parse", f"{branch}-replay:repair.txt"
    )

    git_ops.remove_worktree(disposable_repo, restarted_worktree)
    assert _git_text(disposable_repo, "rev-parse", "HEAD") == base_sha
    assert _git_text(disposable_repo, "status", "--porcelain") == ""


def test_c4_signature_and_source_sha_are_pinned() -> None:
    source = SOURCE_PATH.read_bytes()
    source_text = source.decode("utf-8")
    assert hashlib.sha256(source).hexdigest() == SOURCE_SHA256
    for signature in (
        "def add_worktree(",
        "def commit_all(",
        "def new_commits(",
        "def changed_files(",
        "def remove_worktree(",
    ):
        assert signature in source_text
