"""Real disposable Git integration; coding-loop boundary is deterministic.

These checks do not attest provider execution or versioned C4 DTO compatibility.
"""

import asyncio
from pathlib import Path
import subprocess
import sys
from types import SimpleNamespace

import pytest

from swe_af.issue import build


BROKEN = "def answer():\n    return 0\n"
FIXED = "def answer():\n    return 42\n"
PARTIAL = "def answer():\n    return 41\n"


def git(repo, *arguments):
    return subprocess.run(
        ["git", "-C", str(repo), *arguments],
        check=True, capture_output=True, text=True,
    ).stdout.strip()


@pytest.fixture
def repository(tmp_path):
    repo = tmp_path / "caller"
    repo.mkdir()
    git(repo, "init", "-b", "main")
    git(repo, "config", "user.name", "Scoped Test")
    git(repo, "config", "user.email", "scoped-test@example.invalid")
    git(repo, "config", "commit.gpgsign", "false")
    (repo / "solution.py").write_text(BROKEN, encoding="utf-8")
    git(repo, "add", "solution.py")
    git(repo, "commit", "-m", "seed broken implementation")
    return repo


def verifier_exit(worktree):
    return subprocess.run(
        [sys.executable, "-c", "from solution import answer; assert answer() == 42"],
        cwd=worktree, capture_output=True, text=True,
    ).returncode


def assert_isolated(repo, original_head):
    assert git(repo, "rev-parse", "HEAD") == original_head
    assert git(repo, "branch", "--show-current") == "main"
    assert (repo / "solution.py").read_text(encoding="utf-8") == BROKEN
    assert git(repo, "status", "--porcelain") == ""
    assert git(repo, "remote") == ""
    assert git(repo, "worktree", "list", "--porcelain").count("worktree ") == 1


def run_issue(monkeypatch, repo, *, source, outcome="completed", base="", interrupt=False):
    events = []

    async def coding_loop(**kwargs):
        worktree = Path(kwargs["issue"]["worktree_path"])
        events.append(("coding", (worktree / "solution.py").read_text(encoding="utf-8")))
        (worktree / "solution.py").write_text(source, encoding="utf-8")
        if interrupt:
            raise RuntimeError("deterministic interruption after tracked repair")
        loop_outcome = (
            build.IssueOutcome.COMPLETED
            if outcome == "completed"
            else SimpleNamespace(value=outcome)
        )
        return SimpleNamespace(
            outcome=loop_outcome, result_summary="deterministic tracked repair",
            error_message="" if outcome == "completed" else "review rejected",
            attempts=1, iteration_history=[], debt_items=[],
        )

    async def call_fn(target, **kwargs):
        assert target == "scoped.run_verifier", "unexpected external agent effect"
        events.append(("verifier", kwargs["repo_path"]))
        return {"passed": verifier_exit(kwargs["repo_path"]) == 0}

    monkeypatch.setattr(build, "run_coding_loop", coding_loop)
    result = asyncio.run(build._implement_issue_impl(
        issue={
            "title": "Return the required answer",
            "description": "Repair the tracked solution without touching the caller.",
            "acceptance_criteria": ["answer() returns 42"],
        },
        repo_path=str(repo), base_branch=base, artifacts_dir=".artifacts",
        additional_context="", config={"verify": True, "enable_github_pr": False},
        call_fn=call_fn, note_fn=None, node_id="scoped",
    ))
    return result, events


def test_positive_real_git_and_verifier(monkeypatch, repository):
    original_head = git(repository, "rev-parse", "HEAD")
    assert verifier_exit(repository) != 0
    result, events = run_issue(monkeypatch, repository, source=FIXED)
    assert result["success"] is True
    assert result["base_sha"] == original_head
    assert result["verification"]["passed"] is True
    assert len(result["commits"]) == 1
    assert result["files_changed"] == ["solution.py"]
    assert git(repository, "show", f"{result['branch']}:solution.py") == FIXED.strip()
    assert [event[0] for event in events] == ["coding", "verifier"]
    assert result["pr_url"] == ""
    assert_isolated(repository, original_head)


def test_rejected_loop_cannot_become_success(monkeypatch, repository):
    original_head = git(repository, "rev-parse", "HEAD")
    result, events = run_issue(monkeypatch, repository, source=FIXED, outcome="rejected")
    assert result["success"] is False
    assert result["verification"] is None
    assert len(result["commits"]) == 1
    assert [event[0] for event in events] == ["coding"]
    assert result["pr_url"] == ""
    assert_isolated(repository, original_head)


def test_completed_loop_with_failed_verifier_cannot_succeed(monkeypatch, repository):
    original_head = git(repository, "rev-parse", "HEAD")
    result, events = run_issue(monkeypatch, repository, source=PARTIAL)
    assert result["success"] is False
    assert result["verification"]["passed"] is False
    assert len(result["commits"]) == 1
    assert [event[0] for event in events] == ["coding", "verifier"]
    assert result["pr_url"] == ""
    assert_isolated(repository, original_head)


def test_interrupted_checkpoint_restart_has_no_duplicate_effect(monkeypatch, repository):
    original_head = git(repository, "rev-parse", "HEAD")
    interrupted, first_events = run_issue(
        monkeypatch, repository, source=PARTIAL, interrupt=True,
    )
    assert interrupted["success"] is False
    assert "deterministic interruption" in interrupted["error_message"]
    assert len(interrupted["commits"]) == 1
    assert [event[0] for event in first_events] == ["coding"]
    assert_isolated(repository, original_head)
    checkpoint = git(repository, "rev-parse", interrupted["branch"])
    partial_blob = git(repository, "rev-parse", f"{checkpoint}:solution.py")
    restarted, second_events = run_issue(
        monkeypatch, repository, source=FIXED, base=interrupted["branch"],
    )
    assert restarted["success"] is True
    assert restarted["base_sha"] == checkpoint
    assert second_events[0] == ("coding", PARTIAL)
    assert [event[0] for event in second_events] == ["coding", "verifier"]
    assert len(restarted["commits"]) == 1
    assert git(repository, "rev-list", "--count", f"{original_head}..{restarted['branch']}") == "2"
    assert git(repository, "rev-parse", interrupted["branch"]) == checkpoint
    assert git(repository, "rev-parse", f"{checkpoint}:solution.py") == partial_blob
    assert restarted["pr_url"] == interrupted["pr_url"] == ""
    assert_isolated(repository, original_head)
