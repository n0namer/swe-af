from __future__ import annotations

import asyncio
import hashlib
import subprocess
from pathlib import Path
from types import SimpleNamespace

import pytest

from swe_af.execution.schemas import IssueOutcome
from swe_af.issue import build


def _git(repo: Path, *args: str) -> str:
    result = subprocess.run(
        ["git", *args],
        cwd=repo,
        check=True,
        capture_output=True,
        text=True,
    )
    return result.stdout.strip()


def _new_repo(tmp_path: Path) -> Path:
    repo = tmp_path / "caller"
    repo.mkdir()
    _git(repo, "init", "-b", "main")
    _git(repo, "config", "user.email", "test@example.invalid")
    _git(repo, "config", "user.name", "Northstar Test")
    (repo / "README.md").write_text("base\n", encoding="utf-8")
    _git(repo, "add", "README.md")
    _git(repo, "commit", "-m", "base")
    return repo


def _loop_result(outcome: object, summary: str) -> SimpleNamespace:
    return SimpleNamespace(
        outcome=outcome,
        result_summary=summary,
        error_message="" if outcome is IssueOutcome.COMPLETED else summary,
        attempts=1,
        iteration_history=[],
        debt_items=[],
    )


def test_disposable_tasks_red_green_and_source_bound_receipt(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo = _new_repo(tmp_path)
    base_sha = _git(repo, "rev-parse", "HEAD")
    build_source_sha = hashlib.sha256(
        Path(build.__file__).read_bytes()
    ).hexdigest()
    calls = 0

    async def fake_loop(*, dag_state, **_: object) -> SimpleNamespace:
        nonlocal calls
        calls += 1
        (Path(dag_state.repo_path) / "repair.txt").write_text(
            f"attempt-{calls}\n", encoding="utf-8"
        )
        if calls == 1:
            return _loop_result(SimpleNamespace(value="rejected"), "rejected output")
        return _loop_result(IssueOutcome.COMPLETED, "bounded repair complete")

    monkeypatch.setattr(build, "run_coding_loop", fake_loop)

    issue = {
        "title": "Northstar disposable task",
        "description": "Repair one isolated issue.",
        "acceptance_criteria": ["The repair is committed on a disposable branch."],
        "name": "northstar-disposable",
    }
    config = {
        "verify": False,
        "enable_github_pr": False,
        "keep_worktree": False,
    }

    first = asyncio.run(
        build._implement_issue_impl(
            issue=issue,
            repo_path=str(repo),
            base_branch="main",
            artifacts_dir=".artifacts",
            additional_context="",
            config=config,
            call_fn=lambda *_args, **_kwargs: None,
            note_fn=None,
            node_id="test",
        )
    )
    second = asyncio.run(
        build._implement_issue_impl(
            issue=issue,
            repo_path=str(repo),
            base_branch="main",
            artifacts_dir=".artifacts",
            additional_context="",
            config=config,
            call_fn=lambda *_args, **_kwargs: None,
            note_fn=None,
            node_id="test",
        )
    )

    assert first["success"] is False
    assert first["outcome"] == "rejected"
    assert first["base_sha"] == base_sha
    assert second["success"] is True
    assert second["base_sha"] == base_sha
    assert len(first["commits"]) == 1
    assert len(second["commits"]) == 1
    assert first["branch"] != second["branch"]
    assert build_source_sha and len(build_source_sha) == 64
    assert _git(repo, "status", "--porcelain") == ""
    assert not (repo / "repair.txt").exists()
    assert not (repo / ".worktrees").exists() or not list((repo / ".worktrees").iterdir())


def test_invalid_payload_and_cancelled_replay_fail_closed(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    repo = _new_repo(tmp_path)
    base_sha = _git(repo, "rev-parse", "HEAD")

    with pytest.raises(Exception):
        asyncio.run(
            build._implement_issue_impl(
                issue={"description": "missing title"},
                repo_path=str(repo),
                base_branch="main",
                artifacts_dir=".artifacts",
                additional_context="",
                config={"verify": False, "enable_github_pr": False},
                call_fn=lambda *_args, **_kwargs: None,
                note_fn=None,
                node_id="test",
            )
        )

    async def cancelled_loop(*, dag_state, **_: object) -> SimpleNamespace:
        (Path(dag_state.repo_path) / "interrupted.txt").write_text(
            "checkpointed\n", encoding="utf-8"
        )
        raise asyncio.CancelledError()

    monkeypatch.setattr(build, "run_coding_loop", cancelled_loop)
    with pytest.raises(asyncio.CancelledError):
        asyncio.run(
            build._implement_issue_impl(
                issue={
                    "title": "cancelled disposable task",
                    "description": "Preserve interruption work without touching caller.",
                    "name": "cancelled-disposable",
                },
                repo_path=str(repo),
                base_branch="main",
                artifacts_dir=".artifacts",
                additional_context="",
                config={"verify": False, "enable_github_pr": False},
                call_fn=lambda *_args, **_kwargs: None,
                note_fn=None,
                node_id="test",
            )
        )

    assert _git(repo, "rev-parse", "HEAD") == base_sha
    assert _git(repo, "status", "--porcelain") == ""
    assert not (repo / "interrupted.txt").exists()
    assert len(_git(repo, "branch", "--list", "issue/*").splitlines()) == 1

    calls = 0

    async def replay_loop(*, dag_state, **_: object) -> SimpleNamespace:
        nonlocal calls
        calls += 1
        (Path(dag_state.repo_path) / "replayed.txt").write_text(
            "replayed once\n", encoding="utf-8"
        )
        return _loop_result(IssueOutcome.COMPLETED, "replayed")

    monkeypatch.setattr(build, "run_coding_loop", replay_loop)
    replay = asyncio.run(
        build._implement_issue_impl(
            issue={
                "title": "replayed disposable task",
                "description": "Replay after cancellation.",
                "name": "replayed-disposable",
            },
            repo_path=str(repo),
            base_branch="main",
            artifacts_dir=".artifacts",
            additional_context="",
            config={"verify": False, "enable_github_pr": False},
            call_fn=lambda *_args, **_kwargs: None,
            note_fn=None,
            node_id="test",
        )
    )

    assert replay["success"] is True
    assert len(replay["commits"]) == 1
    assert calls == 1
    assert _git(repo, "rev-parse", "HEAD") == base_sha
    assert _git(repo, "status", "--porcelain") == ""
