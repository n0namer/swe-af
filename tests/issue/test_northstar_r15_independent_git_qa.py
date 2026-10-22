"""Independent Git oracle for the native issue entrypoint, not an agent-quality test."""

from __future__ import annotations

import ast
import asyncio
import hashlib
import importlib
import json
import os
from pathlib import Path
import subprocess
import sys
from types import SimpleNamespace

import pytest


ROOT = Path(__file__).resolve().parents[2]
SOURCE = ROOT / "swe_af/issue/build.py"
CONTRACT = ROOT / "docs/architecture/basefive-r14-contracts.json"
SOURCE_SHA = "440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd"
BROKEN = "answer = 41\n"
REPAIRED = "answer = 42\n"
ISSUE = {
    "name": "r15-independent-repair",
    "title": "Repair the answer",
    "description": "Change answer from 41 to 42 without modifying the caller.",
    "acceptance_criteria": ["The independent Python assertion exits zero."],
    "files_to_modify": ["answer.py"],
}


def digest(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def git(repo: Path, *arguments: str) -> str:
    environment = os.environ.copy()
    environment.update(
        GIT_CONFIG_NOSYSTEM="1",
        GIT_CONFIG_GLOBAL=os.devnull,
        GIT_TERMINAL_PROMPT="0",
    )
    result = subprocess.run(
        ["git", "-C", str(repo), *arguments],
        env=environment,
        capture_output=True,
        text=True,
        timeout=20,
        check=True,
    )
    return result.stdout.strip()


def oracle(worktree: Path) -> subprocess.CompletedProcess:
    environment = os.environ.copy()
    environment["PYTHONDONTWRITEBYTECODE"] = "1"
    return subprocess.run(
        [sys.executable, "-I", "-c", "exec(open('answer.py').read()); assert answer == 42"],
        cwd=worktree,
        env=environment,
        capture_output=True,
        text=True,
        timeout=20,
        check=False,
    )


@pytest.fixture
def native():
    return importlib.import_module("swe_af.issue.build")


@pytest.fixture
def repository(tmp_path):
    repo = tmp_path / "repository"
    repo.mkdir()
    git(repo, "init", "-b", "main")
    git(repo, "config", "user.name", "Disposable R15 QA")
    git(repo, "config", "user.email", "r15-qa@example.invalid")
    git(repo, "config", "commit.gpgsign", "false")
    git(repo, "config", "core.autocrlf", "false")
    (repo / "answer.py").write_text(BROKEN, encoding="utf-8")
    git(repo, "add", "answer.py")
    git(repo, "commit", "-m", "fixture: broken answer")
    base_sha = git(repo, "rev-parse", "HEAD")
    (repo / "answer.py").write_text("answer = 99\n", encoding="utf-8")
    return SimpleNamespace(repo=repo, base_sha=base_sha)


def assert_isolated(state, expected_worktrees=1):
    assert git(state.repo, "rev-parse", "HEAD") == state.base_sha
    assert git(state.repo, "branch", "--show-current") == "main"
    assert (state.repo / "answer.py").read_text(encoding="utf-8") == "answer = 99\n"
    assert git(state.repo, "status", "--porcelain") == "M answer.py"
    assert git(state.repo, "remote") == ""
    worktrees = git(state.repo, "worktree", "list", "--porcelain")
    assert sum(line.startswith("worktree ") for line in worktrees.splitlines()) == expected_worktrees


def install_boundaries(monkeypatch, native, state, loop, *, verifier_unavailable=False):
    calls = []

    async def call(target, **kwargs):
        calls.append(target)
        assert target == "r15-qa.run_verifier", "Unexpected external agent/send boundary"
        worktree = Path(kwargs["repo_path"]).resolve()
        assert worktree.is_relative_to(state.repo.resolve() / ".worktrees")
        if verifier_unavailable:
            raise RuntimeError("independent verifier unavailable")
        checked = oracle(worktree)
        return {"passed": checked.returncode == 0, "summary": checked.stderr}

    app = SimpleNamespace(call=call, node_id="r15-qa", note=lambda *args, **kwargs: None)
    monkeypatch.setattr(native, "issue_router", SimpleNamespace(app=app))
    monkeypatch.setattr(native, "run_coding_loop", loop)
    return calls


def loop_result(native, *, completed=True):
    outcome = native.IssueOutcome.COMPLETED if completed else SimpleNamespace(value="failed")
    return SimpleNamespace(
        outcome=outcome,
        result_summary="Independent deterministic coder fixture",
        error_message="" if completed else "Independent rejection",
        attempts=1,
        iteration_history=[],
        debt_items=[],
    )


async def implement(native, state):
    return await native.implement_issue(
        issue=dict(ISSUE),
        repo_path=str(state.repo),
        base_branch="main",
        config={"verify": True, "enable_github_pr": False, "keep_worktree": False},
    )


def assert_deliverable(state, result):
    assert result["base_sha"] == state.base_sha
    assert result["branch"].startswith("issue/")
    assert git(state.repo, "rev-list", "--count", f"{state.base_sha}..{result['branch']}") == "1"
    assert len(result["commits"]) == 1
    assert result["commits"][0] == git(state.repo, "rev-parse", result["branch"])
    assert result["files_changed"] == ["answer.py"]
    assert git(state.repo, "show", f"{result['branch']}:answer.py") == REPAIRED.strip()
    assert result["pr_url"] == ""
    assert_isolated(state)


def test_source_bound_c4_signature_and_fingerprint():
    assert digest(SOURCE) == SOURCE_SHA, "Source changed: require fresh independent QA"
    contract = json.loads(CONTRACT.read_text(encoding="utf-8"))
    assert contract["sources"]["swe_af/issue/build.py"] == digest(SOURCE)
    tree = ast.parse(SOURCE.read_text(encoding="utf-8"))
    entry = next(node for node in tree.body if isinstance(node, ast.AsyncFunctionDef) and node.name == "implement_issue")
    assert [argument.arg for argument in entry.args.args] == [
        "issue", "repo_path", "base_branch", "artifacts_dir", "additional_context", "config"
    ]
    assert [ast.unparse(argument.annotation) for argument in entry.args.args] == [
        "dict", "str", "str", "str", "str", "dict | None"
    ]
    assert [ast.literal_eval(default) for default in entry.args.defaults] == ["", ".artifacts", "", None]
    assert ast.unparse(entry.returns) == "dict"
    assert contract["dto"]["comparison"] == "DIFFERENT_ENTRYPOINTS_NOT_PARITY"
    assert contract["dto"]["wireSchemaVersion"] is None
    print(json.dumps({"source_sha256": digest(SOURCE), "test_sha256": digest(Path(__file__)), "contract_sha256": digest(CONTRACT)}, sort_keys=True))


def test_native_red_repair_green(monkeypatch, native, repository):
    observations = []

    async def loop(**kwargs):
        worktree = Path(kwargs["issue"]["worktree_path"])
        assert (worktree / "answer.py").read_text(encoding="utf-8") == BROKEN
        observations.append(oracle(worktree).returncode)
        (worktree / "answer.py").write_text(REPAIRED, encoding="utf-8")
        observations.append(oracle(worktree).returncode)
        return loop_result(native)

    calls = install_boundaries(monkeypatch, native, repository, loop)
    result = asyncio.run(implement(native, repository))
    assert observations[0] != 0 and observations[1] == 0
    assert result["success"] is True
    assert result["verification"]["passed"] is True
    assert calls == ["r15-qa.run_verifier"]
    assert_deliverable(repository, result)


@pytest.mark.parametrize("rejection", ["coding", "verification", "unavailable"])
def test_native_rejection_never_success(monkeypatch, native, repository, rejection):
    async def loop(**kwargs):
        worktree = Path(kwargs["issue"]["worktree_path"])
        replacement = "answer = 40\n" if rejection == "verification" else REPAIRED
        (worktree / "answer.py").write_text(replacement, encoding="utf-8")
        return loop_result(native, completed=rejection != "coding")

    calls = install_boundaries(
        monkeypatch, native, repository, loop, verifier_unavailable=rejection == "unavailable"
    )
    result = asyncio.run(implement(native, repository))
    assert result["success"] is False
    assert len(result["commits"]) == 1
    if rejection == "coding":
        assert calls == []
        assert result["verification"] is None
        assert result["outcome"] == "failed"
    else:
        assert calls == ["r15-qa.run_verifier"]
        assert result["verification"]["passed"] is False
    assert_isolated(repository)


def test_native_crash_salvage_and_fresh_replay(monkeypatch, native, repository):
    async def crash(**kwargs):
        worktree = Path(kwargs["issue"]["worktree_path"])
        (worktree / "answer.py").write_text(REPAIRED, encoding="utf-8")
        raise RuntimeError("r15 injected coder interruption")

    calls = install_boundaries(monkeypatch, native, repository, crash)
    interrupted = asyncio.run(implement(native, repository))
    assert interrupted["success"] is False
    assert interrupted["outcome"] == "error"
    assert "r15 injected coder interruption" in interrupted["error_message"]
    assert calls == []
    assert_deliverable(repository, interrupted)
    salvaged_blob = git(repository.repo, "rev-parse", f"{interrupted['branch']}:answer.py")

    async def replay(**kwargs):
        worktree = Path(kwargs["issue"]["worktree_path"])
        assert oracle(worktree).returncode != 0
        (worktree / "answer.py").write_text(REPAIRED, encoding="utf-8")
        return loop_result(native)

    calls = install_boundaries(monkeypatch, native, repository, replay)
    recovered = asyncio.run(implement(native, repository))
    assert recovered["success"] is True
    assert recovered["branch"] != interrupted["branch"]
    assert recovered["build_id"] != interrupted["build_id"]
    assert calls == ["r15-qa.run_verifier"]
    assert_deliverable(repository, recovered)
    assert git(repository.repo, "rev-parse", f"{recovered['branch']}:answer.py") == salvaged_blob
    assert git(repository.repo, "rev-parse", interrupted["branch"]) == interrupted["commits"][0]


def test_native_cancellation_salvages_before_propagation(monkeypatch, native, repository):
    captured = {}

    async def scenario():
        started = asyncio.Event()
        hold = asyncio.Event()

        async def loop(**kwargs):
            captured.update(kwargs["issue"])
            worktree = Path(captured["worktree_path"])
            (worktree / "answer.py").write_text(REPAIRED, encoding="utf-8")
            started.set()
            await hold.wait()
            raise AssertionError("Cancelled coder must not complete")

        calls = install_boundaries(monkeypatch, native, repository, loop)
        task = asyncio.create_task(implement(native, repository))
        try:
            await asyncio.wait_for(started.wait(), timeout=10)
            task.cancel()
            with pytest.raises(asyncio.CancelledError):
                await asyncio.wait_for(task, timeout=10)
        finally:
            if not task.done():
                task.cancel()
                await asyncio.gather(task, return_exceptions=True)
        assert calls == []

    asyncio.run(scenario())
    branch = captured["branch_name"]
    assert git(repository.repo, "rev-list", "--count", f"{repository.base_sha}..{branch}") == "1"
    assert git(repository.repo, "show", f"{branch}:answer.py") == REPAIRED.strip()
    assert not Path(captured["worktree_path"]).exists()
    assert_isolated(repository)


def test_native_preserves_unrelated_stale_worktree(monkeypatch, native, repository, tmp_path):
    stale_path = tmp_path / "abandoned-worktree"
    git(repository.repo, "worktree", "add", "-b", "r15-abandoned", str(stale_path), repository.base_sha)
    stale_head = git(stale_path, "rev-parse", "HEAD")
    (stale_path / "answer.py").write_text("answer = -1\n", encoding="utf-8")

    async def loop(**kwargs):
        worktree = Path(kwargs["issue"]["worktree_path"])
        assert worktree != stale_path
        (worktree / "answer.py").write_text(REPAIRED, encoding="utf-8")
        return loop_result(native)

    calls = install_boundaries(monkeypatch, native, repository, loop)
    result = asyncio.run(implement(native, repository))
    assert result["success"] is True
    assert calls == ["r15-qa.run_verifier"]
    assert git(stale_path, "rev-parse", "HEAD") == stale_head
    assert (stale_path / "answer.py").read_text(encoding="utf-8") == "answer = -1\n"
    assert git(stale_path, "status", "--porcelain") == "M answer.py"
    assert_isolated(repository, expected_worktrees=2)
