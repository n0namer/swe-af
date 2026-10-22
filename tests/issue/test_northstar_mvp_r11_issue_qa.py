from __future__ import annotations

import asyncio
import hashlib
import subprocess
from pathlib import Path
from types import SimpleNamespace

import pytest

from swe_af.execution.schemas import IssueOutcome
from swe_af.issue import build


_FROZEN_R10_REVISION = "064a7226a500365005b5b662044fc9014a328722"


def _qa_git(repo: Path, *args: str) -> str:
    result = subprocess.run(
        ["git", *args],
        cwd=repo,
        check=True,
        capture_output=True,
        text=True,
    )
    return result.stdout.strip()


def _new_qa_repo(tmp_path: Path) -> Path:
    repo = tmp_path / "qa-caller"
    repo.mkdir()
    _qa_git(repo, "init", "-b", "main")
    _qa_git(repo, "config", "user.email", "qa@example.invalid")
    _qa_git(repo, "config", "user.name", "Independent QA")
    (repo / "README.md").write_text("qa base\n", encoding="utf-8")
    _qa_git(repo, "add", "README.md")
    _qa_git(repo, "commit", "-m", "qa base")
    return repo


def _result(outcome: object, summary: str) -> SimpleNamespace:
    return SimpleNamespace(
        outcome=outcome,
        result_summary=summary,
        error_message="" if outcome is IssueOutcome.COMPLETED else summary,
        attempts=1,
        iteration_history=[],
        debt_items=[],
    )


def _run_issue(repo: Path, issue: dict, monkeypatch: pytest.MonkeyPatch, loop) -> dict:
    monkeypatch.setattr(build, "run_coding_loop", loop)
    return asyncio.run(
        build._implement_issue_impl(
            issue=issue,
            repo_path=str(repo),
            base_branch="main",
            artifacts_dir=".qa-artifacts",
            additional_context="independent C4 QA",
            config={
                "verify": False,
                "enable_github_pr": False,
                "keep_worktree": False,
            },
            call_fn=lambda *_args, **_kwargs: None,
            note_fn=None,
            node_id="r11-independent-qa",
        )
    )


def test_r11_independent_qa_positive_red_green_and_frozen_source(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    project_root = Path(build.__file__).parents[2]
    assert _qa_git(project_root, "rev-parse", "HEAD") == _FROZEN_R10_REVISION
    source_sha = hashlib.sha256(Path(build.__file__).read_bytes()).hexdigest()
    assert len(source_sha) == 64

    repo = _new_qa_repo(tmp_path)
    base_sha = _qa_git(repo, "rev-parse", "HEAD")
    calls = 0

    async def red_then_green_loop(*, dag_state, **_: object) -> SimpleNamespace:
        nonlocal calls
        calls += 1
        (Path(dag_state.repo_path) / "qa-repair.txt").write_text(
            f"qa-attempt-{calls}\n", encoding="utf-8"
        )
        if calls == 1:
            return _result(SimpleNamespace(value="rejected"), "independent QA red")
        return _result(IssueOutcome.COMPLETED, "independent QA green")

    issue = {
        "title": "Independent C4 QA repair",
        "description": "Exercise one isolated issue through red then green output.",
        "acceptance_criteria": ["A disposable branch contains the repair."],
        "name": "r11-independent-qa",
    }
    rejected = _run_issue(repo, issue, monkeypatch, red_then_green_loop)
    completed = _run_issue(repo, issue, monkeypatch, red_then_green_loop)

    assert rejected["success"] is False
    assert rejected["outcome"] == "rejected"
    assert completed["success"] is True
    assert completed["outcome"] == IssueOutcome.COMPLETED.value
    assert rejected["base_sha"] == completed["base_sha"] == base_sha
    assert len(rejected["commits"]) == len(completed["commits"]) == 1
    assert rejected["branch"] != completed["branch"]
    assert calls == 2
    assert _qa_git(repo, "status", "--porcelain") == ""
    assert not (repo / "qa-repair.txt").exists()
    assert not (repo / ".worktrees").exists() or not list((repo / ".worktrees").iterdir())


def test_r11_independent_qa_invalid_payload_and_cancelled_replay(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo = _new_qa_repo(tmp_path)
    base_sha = _qa_git(repo, "rev-parse", "HEAD")

    async def should_not_run(**_: object) -> SimpleNamespace:
        raise AssertionError("invalid payload must fail before the coding loop")

    with pytest.raises(Exception):
        _run_issue(repo, {"description": "title is required"}, monkeypatch, should_not_run)
    assert _qa_git(repo, "rev-parse", "HEAD") == base_sha
    assert _qa_git(repo, "status", "--porcelain") == ""

    async def cancelled_loop(*, dag_state, **_: object) -> SimpleNamespace:
        (Path(dag_state.repo_path) / "qa-interrupted.txt").write_text(
            "checkpoint this repair\n", encoding="utf-8"
        )
        raise asyncio.CancelledError()

    interrupted_issue = {
        "title": "Interrupted independent QA repair",
        "description": "Checkpoint work without mutating the caller checkout.",
        "name": "r11-independent-qa-interrupted",
    }
    with pytest.raises(asyncio.CancelledError):
        _run_issue(repo, interrupted_issue, monkeypatch, cancelled_loop)

    interrupted_branches = _qa_git(repo, "branch", "--list", "issue/*").splitlines()
    assert len(interrupted_branches) == 1
    interrupted_branch = interrupted_branches[0].strip()
    assert _qa_git(repo, "show", f"{interrupted_branch}:qa-interrupted.txt") == "checkpoint this repair"
    assert _qa_git(repo, "rev-parse", "HEAD") == base_sha
    assert _qa_git(repo, "status", "--porcelain") == ""
    assert not (repo / "qa-interrupted.txt").exists()

    async def replay_loop(*, dag_state, **_: object) -> SimpleNamespace:
        (Path(dag_state.repo_path) / "qa-replayed.txt").write_text(
            "fresh replay\n", encoding="utf-8"
        )
        return _result(IssueOutcome.COMPLETED, "replayed after cancellation")

    replay = _run_issue(
        repo,
        {
            "title": "Replay independent QA repair",
            "description": "Use a fresh branch after a cancelled repair.",
            "name": "r11-independent-qa-replay",
        },
        monkeypatch,
        replay_loop,
    )

    assert replay["success"] is True
    assert replay["base_sha"] == base_sha
    assert len(replay["commits"]) == 1
    assert replay["branch"] != interrupted_branch
    assert _qa_git(repo, "rev-parse", "HEAD") == base_sha
    assert _qa_git(repo, "status", "--porcelain") == ""
    assert not (repo / "qa-replayed.txt").exists()
