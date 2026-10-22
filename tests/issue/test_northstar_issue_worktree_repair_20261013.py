"""C4 native Git acceptance tests for issue worktree creation and repair.

Live tests use disposable repositories and the real ``git`` executable. The
single monkeypatched test is explicitly limited to retry-budget bookkeeping.
"""

from __future__ import annotations

import subprocess
from pathlib import Path

import pytest

from swe_af.issue import git_ops
from swe_af.issue.git_ops import GitOpsError


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
    _git(repo, "config", "user.name", "C4 Acceptance")
    _git(repo, "config", "user.email", "c4@example.invalid")
    (repo / "README.md").write_text("base\n", encoding="utf-8")
    _git(repo, "add", "README.md")
    _git(repo, "commit", "-m", "base")
    return repo, _git(repo, "rev-parse", "HEAD").stdout.strip()


def _status(repo: Path) -> str:
    return _git(repo, "status", "--porcelain").stdout


def test_live_worktree_commit_preserves_dirty_parent_and_creates_one_commit(
    tmp_path: Path,
) -> None:
    repo, base_sha = _new_repo(tmp_path, "parent")
    dirty_file = repo / "parent-dirty.txt"
    dirty_file.write_text("must survive\n", encoding="utf-8")
    status_before = _status(repo)
    worktree = tmp_path / "worktrees" / "issue-17"
    branch = "issue/17-repair"

    base_ref, resolved_sha = git_ops.resolve_base(str(repo))
    git_ops.add_worktree(str(repo), str(worktree), branch, resolved_sha)
    try:
        (worktree / "solution.txt").write_text("repaired\n", encoding="utf-8")
        commit_sha = git_ops.commit_all(str(worktree), "fix(issue-17): repair")

        assert base_ref == git_ops.current_branch(str(repo))
        assert resolved_sha == base_sha
        assert commit_sha
        assert git_ops.new_commits(str(repo), base_sha, branch) == [commit_sha]
        assert git_ops.changed_files(str(repo), base_sha, branch) == ["solution.txt"]
        assert _git(worktree, "rev-parse", f"{commit_sha}^").stdout.strip() == base_sha
        assert _git(worktree, "show", f"{commit_sha}:solution.txt").stdout == "repaired\n"
        assert _git(repo, "remote").stdout == ""
        assert _status(repo) == status_before
        assert dirty_file.read_text(encoding="utf-8") == "must survive\n"
    finally:
        git_ops.remove_worktree(str(repo), str(worktree))

    assert not worktree.exists()
    assert _status(repo) == status_before


def test_live_duplicate_branch_fails_once_without_mutating_parent(
    tmp_path: Path,
) -> None:
    repo, base_sha = _new_repo(tmp_path, "duplicate")
    dirty_file = repo / "keep.txt"
    dirty_file.write_text("keep\n", encoding="utf-8")
    status_before = _status(repo)
    first = tmp_path / "worktrees" / "first"
    second = tmp_path / "worktrees" / "second"
    branch = "issue/duplicate"

    git_ops.add_worktree(str(repo), str(first), branch, base_sha)
    try:
        with pytest.raises(GitOpsError, match=r"failed after 1 attempts"):
            git_ops.add_worktree(str(repo), str(second), branch, base_sha, attempts=4)
        assert not second.exists()
        assert _status(repo) == status_before
        assert dirty_file.read_text(encoding="utf-8") == "keep\n"
    finally:
        git_ops.remove_worktree(str(repo), str(first))
        git_ops.delete_branch(str(repo), branch)


def test_live_foreign_base_is_rejected_without_branch_or_worktree(
    tmp_path: Path,
) -> None:
    repo, _ = _new_repo(tmp_path, "target")
    foreign_repo, foreign_sha = _new_repo(tmp_path, "foreign")
    (foreign_repo / "foreign.txt").write_text("foreign\n", encoding="utf-8")
    _git(foreign_repo, "add", "foreign.txt")
    _git(foreign_repo, "commit", "-m", "foreign")
    foreign_sha = _git(foreign_repo, "rev-parse", "HEAD").stdout.strip()
    dirty_file = repo / "preserve.txt"
    dirty_file.write_text("preserve\n", encoding="utf-8")
    status_before = _status(repo)
    worktree = tmp_path / "worktrees" / "foreign"
    branch = "issue/foreign-base"

    with pytest.raises(GitOpsError, match=r"failed after 1 attempts"):
        git_ops.add_worktree(str(repo), str(worktree), branch, foreign_sha, attempts=3)

    assert not worktree.exists()
    assert _git(repo, "show-ref", "--verify", f"refs/heads/{branch}", check=False).returncode != 0
    assert _status(repo) == status_before
    assert dirty_file.read_text(encoding="utf-8") == "preserve\n"


def test_live_stale_branch_requires_cleanup_then_restart_has_one_final_commit(
    tmp_path: Path,
) -> None:
    repo, base_sha = _new_repo(tmp_path, "restart")
    interrupted = tmp_path / "worktrees" / "interrupted"
    restarted = tmp_path / "worktrees" / "restarted"
    branch = "issue/restart"

    git_ops.add_worktree(str(repo), str(interrupted), branch, base_sha)
    (interrupted / "partial.txt").write_text("interrupted\n", encoding="utf-8")
    git_ops.remove_worktree(str(repo), str(interrupted))

    with pytest.raises(GitOpsError, match=r"failed after 1 attempts"):
        git_ops.add_worktree(str(repo), str(restarted), branch, base_sha, attempts=5)

    git_ops.delete_branch(str(repo), branch)
    git_ops.add_worktree(str(repo), str(restarted), branch, base_sha)
    try:
        assert not (restarted / "partial.txt").exists()
        (restarted / "solution.txt").write_text("final repair\n", encoding="utf-8")
        final_sha = git_ops.commit_all(str(restarted), "fix(restart): final repair")

        assert final_sha
        assert git_ops.new_commits(str(repo), base_sha, branch) == [final_sha]
        assert _git(restarted, "show", f"{final_sha}:solution.txt").stdout == "final repair\n"
    finally:
        git_ops.remove_worktree(str(repo), str(restarted))


def test_retry_budget_contract_uses_all_transient_attempts(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    """Unit seam only: retry counting; this is not labelled as a live Git path."""
    calls: list[tuple[str, ...]] = []
    sleeps: list[float] = []

    def fail_with_lock(_repo_path: str, *args: str, check: bool = True):
        del check
        calls.append(args)
        return subprocess.CompletedProcess(
            ["git", *args],
            128,
            stdout="",
            stderr="fatal: could not lock config file: File exists",
        )

    monkeypatch.setattr(git_ops, "_git", fail_with_lock)
    monkeypatch.setattr(git_ops.time, "sleep", sleeps.append)

    with pytest.raises(GitOpsError, match=r"failed after 3 attempts"):
        git_ops.add_worktree(
            str(tmp_path / "repo"),
            str(tmp_path / "worktrees" / "retry"),
            "issue/retry",
            "0" * 40,
            attempts=3,
        )

    assert len(calls) == 3
    assert sleeps == [0.5, 1.0]
