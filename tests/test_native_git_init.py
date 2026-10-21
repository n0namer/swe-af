from __future__ import annotations

import subprocess

import pytest

from swe_af.issue.git_ops import GitOpsError, native_git_init


def _git(path, *args):
    return subprocess.run(["git", "-C", str(path), *args], check=True, capture_output=True, text=True).stdout.strip()


def test_native_git_init_existing_repo_creates_isolated_integration_branch(tmp_path):
    _git(tmp_path, "init", "-q")
    (tmp_path / "README.md").write_text("test\n")
    _git(tmp_path, "add", "README.md")
    _git(tmp_path, "-c", "user.name=QA", "-c", "user.email=qa@localhost", "commit", "-qm", "initial")
    initial = _git(tmp_path, "rev-parse", "HEAD")
    original = _git(tmp_path, "branch", "--show-current")

    result = native_git_init(str(tmp_path), "run42")

    assert result["success"] is True
    assert result["mode"] == "existing"
    assert result["initial_commit_sha"] == initial
    assert result["original_branch"] == original
    assert result["integration_branch"] == "run42-integration"
    assert _git(tmp_path, "branch", "--show-current") == "run42-integration"
    assert (tmp_path / ".worktrees").is_dir()
    assert ".worktrees/" in (tmp_path / ".gitignore").read_text()


def test_native_git_init_refuses_dirty_source(tmp_path):
    _git(tmp_path, "init", "-q")
    (tmp_path / "README.md").write_text("initial")
    _git(tmp_path, "add", ".")
    _git(tmp_path, "-c", "user.name=QA", "-c", "user.email=qa@localhost", "commit", "-qm", "initial")
    (tmp_path / "README.md").write_text("changed")
    with pytest.raises(GitOpsError, match="dirty"):
        native_git_init(str(tmp_path), "run43")
    assert not _git(tmp_path, "branch", "--list", "run43-integration")


def test_native_git_init_fresh_repo(tmp_path):
    result = native_git_init(str(tmp_path), "run44")
    assert result["success"] is True
    assert result["mode"] == "fresh"
    assert result["initial_commit_sha"]
    assert result["integration_branch"] == "run44-integration"
