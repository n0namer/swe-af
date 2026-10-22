"""Independent native Git replay oracle for the C4 release-gate boundary.

The oracle intentionally avoids ``tests.issue.conftest`` so collection does not
inherit the producer test's import-shadowing hazard. Each case uses disposable
real Git state and checks both the subprocess exit and the resulting refs.
"""

from __future__ import annotations

import hashlib
from pathlib import Path
import subprocess

import pytest

from swe_af.issue import git_ops


SOURCE_PATH = Path(__file__).parents[2] / "swe_af" / "issue" / "git_ops.py"
SOURCE_SHA256 = "0f73d9fb3a613f1b68ccf4a8161a93c1c5e6dd6a4a392b2c2c8a6157b7c2403d"
PRODUCER_TEST = Path(__file__).with_name("test_releasegate_git_20261013.py")


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
    _git(repo, "config", "user.name", "Northstar QA")
    _git(repo, "config", "user.email", "northstar-qa@example.invalid")
    (repo_path / "README.md").write_text("# native oracle\n", encoding="utf-8")
    _git(repo, "add", "README.md")
    initial = _git(repo, "commit", "-q", "-m", "initial", check=False)
    assert initial.returncode == 0, initial.stderr
    return repo


def test_native_git_positive_commit_has_exact_source_and_one_effect(
    disposable_repo: str,
    tmp_path: Path,
) -> None:
    """The positive path creates one isolated commit and no caller mutation."""
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
    assert _git_text(disposable_repo, "remote") == ""
    assert _git_text(disposable_repo, "status", "--porcelain") == ""

    git_ops.remove_worktree(disposable_repo, worktree)
    assert not Path(worktree).exists()


def test_native_git_negative_branch_collision_fails_closed(disposable_repo: str) -> None:
    """A deterministic branch collision is rejected without a second mutation."""
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
    assert "reference already exists" in str(collision).lower() or "branch named" in str(collision).lower()
    assert git_ops.new_commits(disposable_repo, base_sha, "issue/collision") == []
    assert _git_text(disposable_repo, "rev-parse", "HEAD") == base_sha

    git_ops.remove_worktree(disposable_repo, first_worktree)
    assert not Path(second_worktree).exists()


def test_native_git_restart_replays_checkpoint_once_after_interruption(
    disposable_repo: str,
) -> None:
    """A restart from the interrupted branch preserves exactly one repair effect."""
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


def test_c4_source_signature_and_producer_receipt_comparison() -> None:
    """Pin the inspected C4 helper surface and expose producer evidence gaps."""
    source = SOURCE_PATH.read_bytes()
    assert hashlib.sha256(source).hexdigest() == SOURCE_SHA256
    for signature in (
        "def add_worktree(",
        "def commit_all(",
        "def new_commits(",
        "def changed_files(",
        "def remove_worktree(",
    ):
        assert signature in source.decode("utf-8")

    assert PRODUCER_TEST.exists()
    producer_text = PRODUCER_TEST.read_text(encoding="utf-8")
    assert "implement_issue" in producer_text
    assert "worktree" in producer_text
    assert "verification" in producer_text
    # No producer worker-result receipt is present in the permitted source scope;
    # keep this comparison fail-closed rather than inventing a task identity.
    assert "worker-result" not in producer_text
