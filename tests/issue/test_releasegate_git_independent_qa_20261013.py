"""Independent adversarial QA for the disposable Git release-gate evidence.

This file deliberately does not import ``tests.issue.conftest``.  It reproduces
that producer test's collection hazard separately, then challenges the public
``implement_issue`` entrypoint with real disposable Git state.
"""

from __future__ import annotations

import asyncio
import os
from pathlib import Path
import subprocess
import sys
from types import SimpleNamespace
from typing import Callable

import pytest

from swe_af.issue import build as issue_build


PRODUCER_TEST = Path(__file__).with_name("test_releasegate_git_20261013.py")
BUGGY_SOURCE = "def normalize(items):\n    return sorted(set(items))\n"
FIXED_SOURCE = "def normalize(items):\n    return sorted(items)\n"
ADD_ISSUE = {
    "name": "independent-add-normalizer",
    "title": "Add stable normalizer",
    "description": "Sort values while preserving duplicates.",
    "acceptance_criteria": ["normalize preserves duplicates while sorting"],
    "files_to_create": ["releasegate.py"],
}
REPAIR_ISSUE = {
    "name": "independent-repair-normalizer",
    "title": "Repair stable normalizer",
    "description": "Replace the duplicate-dropping implementation.",
    "acceptance_criteria": ["normalize preserves duplicates while sorting"],
    "files_to_modify": ["releasegate.py"],
}


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
def independent_git_repo(tmp_path: Path) -> str:
    repo_path = tmp_path / "repo"
    repo_path.mkdir()
    repo = str(repo_path)
    _git(repo, "init", "-q")
    _git(repo, "checkout", "-q", "-b", "main")
    _git(repo, "config", "user.email", "independent-qa@example.invalid")
    _git(repo, "config", "user.name", "Independent QA")
    (repo_path / "README.md").write_text("# independent qa\n", encoding="utf-8")
    _git(repo, "add", "README.md")
    _git(repo, "commit", "-q", "-m", "initial")
    return repo


class ScriptedAgent:
    node_id = "releasegate-independent-qa"

    def __init__(self, handler: Callable):
        self.handler = handler
        self.calls: list[str] = []

    async def call(self, target: str, **kwargs):
        self.calls.append(target.rsplit(".", 1)[-1])
        return await self.handler(target.rsplit(".", 1)[-1], **kwargs)

    def note(self, message: str, tags: list[str] | None = None) -> None:
        del message, tags


def _invoke(monkeypatch, agent: ScriptedAgent, issue: dict, repo: str, base_branch: str) -> dict:
    monkeypatch.setattr(issue_build, "issue_router", SimpleNamespace(app=agent))
    return asyncio.run(
        issue_build.implement_issue(
            issue=issue,
            repo_path=repo,
            base_branch=base_branch,
            artifacts_dir=".independent-releasegate-artifacts",
            config={"max_coding_iterations": 2},
        )
    )


def _write_source(worktree_path: str, source: str) -> None:
    Path(worktree_path, "releasegate.py").write_text(source, encoding="utf-8")


def _probe(repo_path: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [
            sys.executable,
            "-c",
            (
                "from releasegate import normalize; "
                "actual = normalize([3, 1, 3]); "
                "assert actual == [1, 3, 3], actual"
            ),
        ],
        cwd=repo_path,
        capture_output=True,
        text=True,
        check=False,
    )


def test_producer_collection_import_is_shadowable(tmp_path: Path) -> None:
    """A regular third-party ``tests`` package deterministically breaks collection."""
    producer_source = PRODUCER_TEST.read_text(encoding="utf-8")
    assert "from tests.issue.conftest import run_git" in producer_source
    assert not PRODUCER_TEST.parents[1].joinpath("__init__.py").exists()

    shadow_root = tmp_path / "shadow"
    shadow_tests = shadow_root / "tests"
    shadow_tests.mkdir(parents=True)
    (shadow_tests / "__init__.py").write_text("ORIGIN = 'shadow'\n", encoding="utf-8")
    environment = os.environ.copy()
    environment["PYTHONPATH"] = os.pathsep.join(
        item for item in [str(shadow_root), environment.get("PYTHONPATH", "")] if item
    )
    reproduction = subprocess.run(
        [sys.executable, "-c", "from tests.issue.conftest import run_git"],
        cwd=PRODUCER_TEST.parents[2],
        env=environment,
        capture_output=True,
        text=True,
        check=False,
    )

    assert reproduction.returncode != 0
    assert "No module named 'tests.issue'" in reproduction.stderr


def test_real_git_rejection_interruption_restart_has_one_repair_effect(
    independent_git_repo: str,
    monkeypatch,
) -> None:
    """Challenge verifier authority, crash salvage, replay, and source identity."""
    initial_head = _git_text(independent_git_repo, "rev-parse", "HEAD")

    async def rejected_handler(role: str, **kwargs):
        if role == "run_coder":
            _write_source(kwargs["worktree_path"], BUGGY_SOURCE)
            return {"files_changed": ["releasegate.py"], "summary": "buggy", "complete": True}
        if role == "run_code_reviewer":
            return {"approved": True, "blocking": False, "summary": "misleading approval"}
        if role == "run_verifier":
            probe = _probe(kwargs["repo_path"])
            return {
                "passed": probe.returncode == 0,
                "summary": probe.stderr or "unexpected pass",
                "criteria_results": [{"criterion": "preserve duplicates", "passed": probe.returncode == 0}],
            }
        raise AssertionError(f"unexpected role: {role}")

    rejected_agent = ScriptedAgent(rejected_handler)
    rejected = _invoke(monkeypatch, rejected_agent, ADD_ISSUE, independent_git_repo, "main")

    assert rejected["success"] is False
    assert rejected["outcome"] == "completed"
    assert rejected["verification"]["passed"] is False
    assert rejected["base_sha"] == initial_head
    assert len(rejected["commits"]) == 1
    assert rejected_agent.calls == ["run_coder", "run_code_reviewer", "run_verifier"]
    assert _git_text(independent_git_repo, "show", f"{rejected['branch']}:releasegate.py") == BUGGY_SOURCE.rstrip()

    async def interrupted_handler(role: str, **kwargs):
        if role == "run_coder":
            _write_source(kwargs["worktree_path"], FIXED_SOURCE)
            raise RuntimeError("independent simulated crash after source write")
        raise AssertionError(f"crash attempt reached unexpected role: {role}")

    interrupted_agent = ScriptedAgent(interrupted_handler)
    interrupted = _invoke(
        monkeypatch,
        interrupted_agent,
        REPAIR_ISSUE,
        independent_git_repo,
        rejected["branch"],
    )

    assert interrupted["success"] is False
    assert interrupted["outcome"] == "failed_unrecoverable"
    assert interrupted["verification"] is None
    assert interrupted["base_sha"] == rejected["commits"][-1]
    assert len(interrupted["commits"]) == 1
    assert interrupted_agent.calls == ["run_coder"]

    async def repaired_handler(role: str, **kwargs):
        if role == "run_coder":
            _write_source(kwargs["worktree_path"], FIXED_SOURCE)
            return {"files_changed": ["releasegate.py"], "summary": "replayed repair", "complete": True}
        if role == "run_code_reviewer":
            return {"approved": True, "blocking": False, "summary": "repair reviewed"}
        if role == "run_verifier":
            probe = _probe(kwargs["repo_path"])
            return {
                "passed": probe.returncode == 0,
                "summary": probe.stderr or "oracle passed",
                "criteria_results": [{"criterion": "preserve duplicates", "passed": probe.returncode == 0}],
            }
        raise AssertionError(f"unexpected role: {role}")

    repaired_agent = ScriptedAgent(repaired_handler)
    repaired = _invoke(
        monkeypatch,
        repaired_agent,
        REPAIR_ISSUE,
        independent_git_repo,
        rejected["branch"],
    )

    assert repaired["success"] is True
    assert repaired["outcome"] == "completed"
    assert repaired["verification"]["passed"] is True
    assert repaired["base_sha"] == rejected["commits"][-1]
    assert len(repaired["commits"]) == 1
    assert repaired_agent.calls == ["run_coder", "run_code_reviewer", "run_verifier"]

    rejected_blob = _git_text(independent_git_repo, "rev-parse", f"{rejected['branch']}:releasegate.py")
    interrupted_blob = _git_text(independent_git_repo, "rev-parse", f"{interrupted['branch']}:releasegate.py")
    repaired_blob = _git_text(independent_git_repo, "rev-parse", f"{repaired['branch']}:releasegate.py")
    repaired_source = _git_text(independent_git_repo, "show", f"{repaired['branch']}:releasegate.py")

    assert all(len(value) == 40 for value in [rejected_blob, interrupted_blob, repaired_blob])
    assert rejected_blob != repaired_blob
    assert interrupted_blob == repaired_blob
    assert repaired_source.count("return sorted(items)") == 1
    assert "set(items)" not in repaired_source
    assert len({rejected["branch"], interrupted["branch"], repaired["branch"]}) == 3
    assert _git_text(independent_git_repo, "rev-parse", "HEAD") == initial_head
    assert _git_text(independent_git_repo, "rev-parse", "--abbrev-ref", "HEAD") == "main"
    assert _git_text(independent_git_repo, "status", "--porcelain") == ""
    assert _git_text(independent_git_repo, "remote") == ""
    assert _git(independent_git_repo, "cat-file", "-e", "main:releasegate.py", check=False).returncode != 0
    worktrees = Path(independent_git_repo, ".worktrees")
    assert not worktrees.exists() or list(worktrees.iterdir()) == []
