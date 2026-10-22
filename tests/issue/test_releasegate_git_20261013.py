"""Release-gate evidence for ``implement_issue`` on disposable real Git state.

The two logical issues exercise the public entrypoint with scripted role calls
while keeping Git, worktrees, branches, commits, diffs, checkpoint salvage, and
verification probes real.
"""

from __future__ import annotations

import asyncio
import os
import subprocess
import sys
from types import SimpleNamespace
from typing import Callable

from swe_af.issue import build as issue_build
from tests.issue.conftest import run_git


ADD_ISSUE = {
    "name": "add-stable-normalizer",
    "title": "Add stable normalizer",
    "description": "Add a normalizer that sorts values without dropping duplicates.",
    "acceptance_criteria": [
        "normalize([3, 1, 3]) returns [1, 3, 3]",
    ],
    "files_to_create": ["releasegate.py"],
}

REPAIR_ISSUE = {
    "name": "repair-stable-normalizer",
    "title": "Repair stable normalizer",
    "description": "Preserve duplicates while sorting values.",
    "acceptance_criteria": [
        "normalize([3, 1, 3]) returns [1, 3, 3]",
    ],
    "files_to_modify": ["releasegate.py"],
}

BUGGY_SOURCE = "def normalize(items):\n    return sorted(set(items))\n"
FIXED_SOURCE = "def normalize(items):\n    return sorted(items)\n"


class ScriptedAgent:
    """Minimal app surface used by the real public ``implement_issue`` wrapper."""

    node_id = "releasegate-git-test"

    def __init__(self, handler: Callable):
        self.handler = handler
        self.calls: list[str] = []
        self.notes: list[tuple[str, tuple[str, ...]]] = []

    async def call(self, target: str, **kwargs):
        self.calls.append(target)
        return await self.handler(target, **kwargs)

    def note(self, message: str, tags: list[str] | None = None) -> None:
        self.notes.append((message, tuple(tags or [])))


def _invoke(
    monkeypatch,
    *,
    agent: ScriptedAgent,
    issue: dict,
    repo_path: str,
    base_branch: str,
) -> dict:
    monkeypatch.setattr(issue_build, "issue_router", SimpleNamespace(app=agent))
    return asyncio.run(
        issue_build.implement_issue(
            issue=issue,
            repo_path=repo_path,
            base_branch=base_branch,
            artifacts_dir=".releasegate-artifacts",
            config={"max_coding_iterations": 2},
        )
    )


def _write_source(worktree_path: str, source: str) -> None:
    with open(os.path.join(worktree_path, "releasegate.py"), "w", encoding="utf-8") as stream:
        stream.write(source)


def _probe(worktree_path: str) -> subprocess.CompletedProcess[str]:
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
        cwd=worktree_path,
        capture_output=True,
        text=True,
        check=False,
    )


def _role(target: str) -> str:
    return target.split(".", 1)[1]


def _blob_sha(repo_path: str, branch: str) -> str:
    return run_git(repo_path, "rev-parse", f"{branch}:releasegate.py")


def _git_object_exists(repo_path: str, object_name: str) -> bool:
    probe = subprocess.run(
        ["git", "-C", repo_path, "cat-file", "-e", object_name],
        capture_output=True,
        text=True,
        check=False,
    )
    return probe.returncode == 0


def test_releasegate_real_git_reject_repair_interrupt_restart_without_duplicates(
    git_repo: str,
    monkeypatch,
) -> None:
    """Two issues produce bounded, independently inspectable Git evidence."""
    main_head = run_git(git_repo, "rev-parse", "HEAD")

    async def add_handler(target: str, **kwargs):
        role = _role(target)
        if role == "run_coder":
            _write_source(kwargs["worktree_path"], BUGGY_SOURCE)
            return {
                "files_changed": ["releasegate.py"],
                "summary": "added first implementation",
                "complete": True,
            }
        if role == "run_code_reviewer":
            return {
                "approved": True,
                "blocking": False,
                "summary": "source is reviewable; executable verifier owns the criterion",
            }
        if role == "run_verifier":
            probe = _probe(kwargs["repo_path"])
            return {
                "passed": probe.returncode == 0,
                "summary": probe.stderr or "unexpected pass",
                "criteria_results": [
                    {
                        "criterion": ADD_ISSUE["acceptance_criteria"][0],
                        "passed": probe.returncode == 0,
                    }
                ],
            }
        raise AssertionError(f"unexpected role: {role}")

    add_agent = ScriptedAgent(add_handler)
    rejected = _invoke(
        monkeypatch,
        agent=add_agent,
        issue=ADD_ISSUE,
        repo_path=git_repo,
        base_branch="main",
    )

    assert rejected["success"] is False
    assert rejected["outcome"] == "completed"
    assert rejected["verification"]["passed"] is False
    assert rejected["branch"].startswith("issue/")
    assert rejected["base_sha"] == main_head
    assert len(rejected["commits"]) == 1
    assert rejected["files_changed"] == ["releasegate.py"]
    assert run_git(git_repo, "show", f"{rejected['branch']}:releasegate.py") == BUGGY_SOURCE.rstrip()
    assert [target.rsplit(".", 1)[-1] for target in add_agent.calls] == [
        "run_coder",
        "run_code_reviewer",
        "run_verifier",
    ]

    async def interrupted_repair_handler(target: str, **kwargs):
        role = _role(target)
        if role == "run_coder":
            _write_source(kwargs["worktree_path"], FIXED_SOURCE)
            raise RuntimeError("simulated interruption after repair write")
        raise AssertionError(f"interrupted repair must not call {role}")

    interrupted_agent = ScriptedAgent(interrupted_repair_handler)
    interrupted = _invoke(
        monkeypatch,
        agent=interrupted_agent,
        issue=REPAIR_ISSUE,
        repo_path=git_repo,
        base_branch=rejected["branch"],
    )

    assert interrupted["success"] is False
    assert interrupted["outcome"] == "failed_unrecoverable"
    assert interrupted["verification"] is None
    assert interrupted["branch"].startswith("issue/")
    assert interrupted["base_sha"] == rejected["commits"][-1]
    assert len(interrupted["commits"]) == 1
    assert interrupted["files_changed"] == ["releasegate.py"]
    assert run_git(git_repo, "show", f"{interrupted['branch']}:releasegate.py") == FIXED_SOURCE.rstrip()
    assert [target.rsplit(".", 1)[-1] for target in interrupted_agent.calls] == ["run_coder"]

    async def restarted_repair_handler(target: str, **kwargs):
        role = _role(target)
        if role == "run_coder":
            _write_source(kwargs["worktree_path"], FIXED_SOURCE)
            return {
                "files_changed": ["releasegate.py"],
                "summary": "replayed the bounded repair after interruption",
                "complete": True,
            }
        if role == "run_code_reviewer":
            return {"approved": True, "blocking": False, "summary": "repair accepted"}
        if role == "run_verifier":
            probe = _probe(kwargs["repo_path"])
            return {
                "passed": probe.returncode == 0,
                "summary": probe.stderr or "executable criterion passed",
                "criteria_results": [
                    {
                        "criterion": REPAIR_ISSUE["acceptance_criteria"][0],
                        "passed": probe.returncode == 0,
                    }
                ],
            }
        raise AssertionError(f"unexpected role: {role}")

    restarted_agent = ScriptedAgent(restarted_repair_handler)
    repaired = _invoke(
        monkeypatch,
        agent=restarted_agent,
        issue=REPAIR_ISSUE,
        repo_path=git_repo,
        base_branch=rejected["branch"],
    )

    assert repaired["success"] is True
    assert repaired["outcome"] == "completed"
    assert repaired["verification"]["passed"] is True
    assert repaired["iterations"] == 1
    assert repaired["base_sha"] == rejected["commits"][-1]
    assert len(repaired["commits"]) == 1
    assert repaired["files_changed"] == ["releasegate.py"]
    assert [target.rsplit(".", 1)[-1] for target in restarted_agent.calls] == [
        "run_coder",
        "run_code_reviewer",
        "run_verifier",
    ]

    repaired_source = run_git(git_repo, "show", f"{repaired['branch']}:releasegate.py")
    assert repaired_source == FIXED_SOURCE.rstrip()
    assert repaired_source.count("return sorted(items)") == 1
    assert repaired_source.count("set(items)") == 0

    rejected_blob = _blob_sha(git_repo, rejected["branch"])
    interrupted_blob = _blob_sha(git_repo, interrupted["branch"])
    repaired_blob = _blob_sha(git_repo, repaired["branch"])
    assert len(rejected_blob) == len(interrupted_blob) == len(repaired_blob) == 40
    assert rejected_blob != repaired_blob
    assert interrupted_blob == repaired_blob

    branches = {rejected["branch"], interrupted["branch"], repaired["branch"]}
    assert len(branches) == 3
    assert all(len(run_git(git_repo, "rev-parse", branch)) == 40 for branch in branches)
    assert run_git(git_repo, "rev-parse", "HEAD") == main_head
    assert run_git(git_repo, "rev-parse", "--abbrev-ref", "HEAD") == "main"
    assert run_git(git_repo, "status", "--porcelain") == ""
    assert run_git(git_repo, "remote") == ""
    assert _git_object_exists(git_repo, "main:releasegate.py") is False

    worktrees_dir = os.path.join(git_repo, ".worktrees")
    assert not os.path.isdir(worktrees_dir) or os.listdir(worktrees_dir) == []
