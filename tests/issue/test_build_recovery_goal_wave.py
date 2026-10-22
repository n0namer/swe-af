"""Executable recovery contract for two bounded issue-build tasks.

The scripted agent boundary is offline, but all Git branches, worktrees, commits,
and verifier subprocesses are real disposable local-repository operations.
"""

from __future__ import annotations

import asyncio
import os
import subprocess
from collections.abc import Callable

import pytest

from swe_af.issue import git_ops
from swe_af.issue.build import _implement_issue_impl


def _git(repo: str, *args: str) -> str:
    return subprocess.run(
        ["git", "-C", repo, *args],
        check=True,
        capture_output=True,
        text=True,
    ).stdout.strip()


def _init_repo(tmp_path) -> str:
    repo = tmp_path / "repo"
    repo.mkdir()
    subprocess.run(["git", "init", "-q", "-b", "main", str(repo)], check=True)
    _git(str(repo), "config", "user.name", "Test User")
    _git(str(repo), "config", "user.email", "test@example.invalid")
    (repo / "README.md").write_text("base\n", encoding="utf-8")
    _git(str(repo), "add", "README.md")
    _git(str(repo), "commit", "-qm", "base")
    return str(repo)


def _run(repo: str, issue: dict, call_fn: Callable, *, config: dict | None = None) -> dict:
    return asyncio.run(
        _implement_issue_impl(
            issue=issue,
            repo_path=repo,
            base_branch="main",
            artifacts_dir=".artifacts",
            additional_context="",
            config=config or {},
            call_fn=call_fn,
            note_fn=None,
            node_id="goal-wave",
        )
    )


def _write_retry_module(worktree: str, *, working: bool) -> None:
    source = (
        "def retry(operation, attempts=3):\n"
        "    last_error = None\n"
        "    for _ in range(attempts):\n"
        "        try:\n"
        "            return operation()\n"
        "        except Exception as error:\n"
        "            last_error = error\n"
        "    raise last_error\n"
        if working
        else "def retry(operation, attempts=3):\n    return operation()\n"
    )
    with open(os.path.join(worktree, "retryable.py"), "w", encoding="utf-8") as stream:
        stream.write(source)


def test_real_issue_cycles_repair_review_and_preserve_interrupted_work(tmp_path) -> None:
    repo = _init_repo(tmp_path)
    with open(os.path.join(repo, "foreign-note.txt"), "w", encoding="utf-8") as stream:
        stream.write("must stay in caller workspace\n")

    repair_issue = {
        "title": "Add retry helper",
        "description": "Implement bounded retry behavior.",
        "acceptance_criteria": ["retry retries twice before succeeding"],
        "files_to_create": ["retryable.py"],
    }
    coder_calls = 0
    review_calls = 0

    async def repair_call(target: str, **kwargs):
        nonlocal coder_calls, review_calls
        role = target.rsplit(".", 1)[1]
        if role == "run_coder":
            coder_calls += 1
            _write_retry_module(kwargs["worktree_path"], working=coder_calls == 2)
            return {"files_changed": ["retryable.py"], "summary": "coder pass", "complete": True}
        if role == "run_code_reviewer":
            review_calls += 1
            return {
                "approved": review_calls == 2,
                "blocking": False,
                "summary": "first review rejects missing retry; repair accepted",
            }
        if role == "run_verifier":
            probe = subprocess.run(
                [
                    "python",
                    "-c",
                    "from retryable import retry; state=[]; "
                    "assert retry(lambda: state.append(1) or (42 if len(state) == 3 else (_ for _ in ()).throw(ValueError()))) == 42",
                ],
                cwd=kwargs["repo_path"],
                capture_output=True,
                text=True,
            )
            return {"passed": probe.returncode == 0, "summary": probe.stderr or "oracle passed"}
        raise AssertionError(f"unexpected target: {target}")

    repaired = _run(repo, repair_issue, repair_call, config={"max_coding_iterations": 2})

    assert repaired["success"] is True
    assert repaired["verification"]["passed"] is True
    assert coder_calls == review_calls == 2
    assert repaired["commits"] and repaired["branch"]
    assert _git(repo, "show", f"{repaired['branch']}:retryable.py").startswith("def retry")
    assert _git(repo, "status", "--porcelain") == "?? foreign-note.txt"

    interrupted_issue = {
        "title": "Preserve checkpoint after interruption",
        "description": "Write an artifact before a simulated worker interruption.",
        "acceptance_criteria": ["artifact survives interruption on the issue branch"],
        "files_to_create": ["recovery.txt"],
    }

    async def interruption_call(target: str, **kwargs):
        if target.endswith("run_coder"):
            with open(os.path.join(kwargs["worktree_path"], "recovery.txt"), "w", encoding="utf-8") as stream:
                stream.write("preserved\n")
            raise RuntimeError("simulated timeout after write")
        raise AssertionError(f"unexpected target: {target}")

    interrupted = _run(repo, interrupted_issue, interruption_call, config={"verify": False})

    assert interrupted["success"] is False
    assert interrupted["commits"] and interrupted["branch"]
    assert "recovery.txt" in interrupted["files_changed"]
    assert _git(repo, "show", f"{interrupted['branch']}:recovery.txt") == "preserved"
    assert _git(repo, "status", "--porcelain") == "?? foreign-note.txt"


def test_commit_all_reports_no_effect_without_creating_a_revision(tmp_path) -> None:
    repo = _init_repo(tmp_path)
    head_before = _git(repo, "rev-parse", "HEAD")

    assert git_ops.commit_all(repo, "chore: should not commit") == ""
    assert _git(repo, "rev-parse", "HEAD") == head_before
    assert _git(repo, "status", "--porcelain") == ""


def test_branch_collision_does_not_retry_or_sleep(tmp_path, monkeypatch) -> None:
    repo = _init_repo(tmp_path)
    base_sha = _git(repo, "rev-parse", "HEAD")
    branch = "issue/already-exists"
    _git(repo, "branch", branch)

    def unexpected_sleep(_: float) -> None:
        pytest.fail("deterministic branch collision must not be retried")

    monkeypatch.setattr(git_ops.time, "sleep", unexpected_sleep)

    with pytest.raises(git_ops.GitOpsError, match="worktree add failed"):
        git_ops.add_worktree(repo, os.path.join(repo, "collision-worktree"), branch, base_sha)

    assert not os.path.exists(os.path.join(repo, "collision-worktree"))
    assert _git(repo, "rev-parse", branch) == base_sha


def test_exhausted_nonblocking_review_budget_returns_debt_without_duplicate_effects(tmp_path) -> None:
    repo = _init_repo(tmp_path)
    issue = {
        "title": "Bounded review budget",
        "description": "Make one local artifact while reviews remain nonblocking.",
        "acceptance_criteria": ["artifact is committed exactly once"],
        "files_to_create": ["bounded.txt"],
    }
    coder_calls = 0

    async def nonblocking_review(target: str, **kwargs):
        nonlocal coder_calls
        role = target.rsplit(".", 1)[1]
        if role == "run_coder":
            coder_calls += 1
            with open(os.path.join(kwargs["worktree_path"], "bounded.txt"), "w", encoding="utf-8") as stream:
                stream.write(f"attempt {coder_calls}\n")
            return {"files_changed": ["bounded.txt"], "summary": "work written", "complete": True}
        if role == "run_code_reviewer":
            return {"approved": False, "blocking": False, "summary": "needs another pass"}
        raise AssertionError(f"unexpected target: {target}")

    result = _run(
        repo,
        issue,
        nonblocking_review,
        config={"verify": False, "max_coding_iterations": 2},
    )

    assert result["success"] is True
    assert result["outcome"] == "completed_with_debt"
    assert result["iterations"] == coder_calls == 2
    assert len(result["commits"]) == 1
    assert _git(repo, "show", f"{result['branch']}:bounded.txt") == "attempt 2"
