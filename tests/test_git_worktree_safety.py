"""Safety regressions: no false worktree reuse or destructive cleanup."""
import os
import pathlib
import subprocess
import pytest
from swe_af.execution import git_fast_path


def git(repo, *args, check=True):
    p = subprocess.run(["git", "-C", str(repo), *args], capture_output=True, text=True, check=check)
    return p.stdout.strip()


@pytest.fixture
def repo(tmp_path):
    r = tmp_path / "repo"
    r.mkdir()
    git(r, "init", "-q")
    git(r, "config", "user.email", "test@local")
    git(r, "config", "user.name", "QA")
    (r / "base.txt").write_text("base\n")
    git(r, "add", ".")
    git(r, "commit", "-qm", "base")
    git(r, "branch", "integration")
    return r


def test_orphan_directory_must_not_be_reported_as_worktree(repo):
    branch = "issue/build-00-a"
    git(repo, "branch", branch, "integration")
    directory = repo / ".worktrees" / "issue-build-00-a"
    directory.mkdir(parents=True)
    sentinel = directory / "preserve-me.txt"
    sentinel.write_text("user evidence\n")
    with pytest.raises(git_fast_path.GitFastPathError):
        git_fast_path.setup_worktrees(str(repo), "integration", [{"name":"a"}], str(repo / ".worktrees"), "build")
    assert sentinel.read_text() == "user evidence\n"


def test_cleanup_must_preserve_dirty_worktree_and_unmerged_branch(repo):
    wt_dir = repo / ".worktrees"
    setup = git_fast_path.setup_worktrees(str(repo), "integration", [{"name":"x"}], str(wt_dir))
    w = pathlib.Path(setup["workspaces"][0]["worktree_path"])
    (w / "untracked.txt").write_text("valuable\n")
    result = git_fast_path.cleanup_worktrees(str(repo), str(wt_dir), ["issue/00-x"])
    assert result["success"] is False
    assert (w / "untracked.txt").read_text() == "valuable\n"
    assert git(repo, "branch", "--list", "issue/00-x") != ""
