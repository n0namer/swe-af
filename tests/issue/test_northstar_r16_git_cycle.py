"""Real disposable Git coverage for the public issue-build orchestration boundary."""

import asyncio
import hashlib
import inspect
import json
import os
from pathlib import Path
import subprocess
from types import SimpleNamespace

import pytest

from swe_af.execution.schemas import IssueOutcome
from swe_af.issue import build


SOURCE_SHA256 = "440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd"


def git(repo, *arguments):
    environment = os.environ.copy()
    environment.update({
        "GIT_CONFIG_NOSYSTEM": "1",
        "GIT_CONFIG_GLOBAL": os.devnull,
        "GIT_TERMINAL_PROMPT": "0",
    })
    for variable in ("GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE"):
        environment.pop(variable, None)
    result = subprocess.run(
        ["git", "-C", str(repo), *arguments],
        env=environment, text=True, capture_output=True, check=True, timeout=30,
    )
    return result.stdout.strip()


@pytest.fixture
def repository(tmp_path):
    repo = tmp_path / "repository"
    repo.mkdir()
    git(repo, "init", "-b", "main")
    git(repo, "config", "user.name", "Disposable fixture")
    git(repo, "config", "user.email", "fixture@example.invalid")
    git(repo, "config", "commit.gpgsign", "false")
    (repo / "value.txt").write_text("base\n", encoding="utf-8")
    git(repo, "add", "value.txt")
    git(repo, "commit", "-m", "fixture base")
    base = git(repo, "rev-parse", "HEAD")
    (repo / "value.txt").write_text("caller dirty\n", encoding="utf-8")
    return repo, base, git(repo, "status", "--porcelain")


def assert_isolated(repository):
    repo, base, dirty = repository
    assert git(repo, "branch", "--show-current") == "main"
    assert git(repo, "rev-parse", "HEAD") == base
    assert git(repo, "status", "--porcelain") == dirty
    assert (repo / "value.txt").read_text(encoding="utf-8") == "caller dirty\n"
    assert git(repo, "remote") == ""
    worktrees = git(repo, "worktree", "list", "--porcelain")
    assert sum(line.startswith("worktree ") for line in worktrees.splitlines()) == 1


def configure(monkeypatch, *, outcome=IssueOutcome.COMPLETED, verifier="pass", interrupt=False):
    observations = {"loops": [], "calls": []}

    async def coding_loop(**kwargs):
        worktree = Path(kwargs["issue"]["worktree_path"])
        observations["loops"].append({
            "branch": kwargs["issue"]["branch_name"],
            "before": (worktree / "value.txt").read_text(encoding="utf-8"),
        })
        (worktree / "value.txt").write_text("repair\n", encoding="utf-8")
        if interrupt:
            raise asyncio.CancelledError("fixture interruption")
        return SimpleNamespace(
            outcome=outcome, result_summary="fixture coding result",
            error_message="", attempts=1, iteration_history=[], debt_items=[],
        )

    async def call(target, **kwargs):
        observations["calls"].append(target)
        assert target == "fixture.run_verifier"
        assert kwargs["prd"]["acceptance_criteria"] == ["value.txt contains repair"]
        assert Path(kwargs["repo_path"]).is_dir()
        if verifier == "unavailable":
            raise RuntimeError("fixture verifier unavailable")
        assert Path(kwargs["repo_path"], "value.txt").read_text(encoding="utf-8") == "repair\n"
        return {"passed": verifier == "pass", "summary": "fixture verification"}

    monkeypatch.setattr(build, "run_coding_loop", coding_loop)
    monkeypatch.setattr(build.issue_router, "app", SimpleNamespace(
        call=call, note=lambda *args, **kwargs: None, node_id="fixture",
    ))
    return observations


def run_issue(repository, tmp_path, *, base_branch="main"):
    repo, _, _ = repository
    return asyncio.run(build.implement_issue(
        issue={
            "name": "r16-git-cycle", "title": "Repair fixture value",
            "description": "Repair only value.txt in a disposable local repository.",
            "acceptance_criteria": ["value.txt contains repair"],
            "files_to_modify": ["value.txt"],
        },
        repo_path=str(repo), base_branch=base_branch,
        artifacts_dir=str(tmp_path / "artifacts"),
        config={"verify": True, "enable_github_pr": False, "keep_worktree": False},
    ))


def assert_delivery(result, repository):
    repo, base, _ = repository
    assert result["base_sha"] == base
    assert result["base_branch"] == "main"
    assert result["branch"].startswith("issue/")
    assert len(result["commits"]) == 1
    assert git(repo, "rev-list", "--count", f"{base}..{result['branch']}") == "1"
    assert result["files_changed"] == ["value.txt"]
    assert git(repo, "show", f"{result['branch']}:value.txt") == "repair"
    assert result["pr_url"] == ""
    assert_isolated(repository)


def test_source_and_c4_contract():
    root = Path(__file__).resolve().parents[2]
    source = root / "swe_af" / "issue" / "build.py"
    assert hashlib.sha256(source.read_bytes()).hexdigest() == SOURCE_SHA256
    contract = json.loads((root / "docs" / "architecture" / "basefive-r14-contracts.json").read_text(encoding="utf-8"))
    assert contract["sources"]["swe_af/issue/build.py"] == SOURCE_SHA256
    assert contract["api"]["pythonSignature"].startswith("async def implement_issue(")
    signature = inspect.signature(build.implement_issue)
    assert list(signature.parameters) == [
        "issue", "repo_path", "base_branch", "artifacts_dir", "additional_context", "config",
    ]
    assert signature.parameters["base_branch"].default == ""
    assert signature.parameters["artifacts_dir"].default == ".artifacts"
    assert signature.parameters["config"].default is None
    assert contract["dto"]["wireCompatibility"] == "UNVERIFIED"
    assert contract["dto"]["wireSchemaVersion"] is None
    assert contract["acceptance"] == "UNACCEPTED"


def test_completed_issue_has_one_checkpoint_and_verification(repository, tmp_path, monkeypatch):
    observations = configure(monkeypatch)
    result = run_issue(repository, tmp_path)
    assert result["success"] is True
    assert result["outcome"] == "completed"
    assert result["verification"]["passed"] is True
    assert observations["calls"] == ["fixture.run_verifier"]
    assert observations["loops"][0]["before"] == "base\n"
    assert_delivery(result, repository)


def test_rejected_coding_does_not_become_success(repository, tmp_path, monkeypatch):
    rejected = next(
        outcome for outcome in IssueOutcome
        if outcome.value not in ("completed", "completed_with_debt")
    )
    observations = configure(monkeypatch, outcome=rejected)
    result = run_issue(repository, tmp_path)
    assert result["success"] is False
    assert result["outcome"] == rejected.value
    assert result["verification"] is None
    assert observations["calls"] == []
    assert_delivery(result, repository)


@pytest.mark.parametrize("verifier", ["reject", "unavailable"])
def test_failed_verification_cannot_downgrade_trust(repository, tmp_path, monkeypatch, verifier):
    observations = configure(monkeypatch, verifier=verifier)
    result = run_issue(repository, tmp_path)
    assert result["success"] is False
    assert result["verification"]["passed"] is False
    assert observations["calls"] == ["fixture.run_verifier"]
    assert_delivery(result, repository)


def test_cancellation_salvage_and_restart_do_not_duplicate_effects(repository, tmp_path, monkeypatch):
    repo, base, _ = repository
    interrupted = configure(monkeypatch, interrupt=True)
    with pytest.raises(asyncio.CancelledError, match="fixture interruption"):
        run_issue(repository, tmp_path)
    assert interrupted["calls"] == []
    branch = interrupted["loops"][0]["branch"]
    assert git(repo, "rev-list", "--count", f"{base}..{branch}") == "1"
    assert git(repo, "show", f"{branch}:value.txt") == "repair"
    assert_isolated(repository)
    salvaged_head = git(repo, "rev-parse", branch)
    salvaged_blob = git(repo, "rev-parse", f"{branch}:value.txt")

    restarted = configure(monkeypatch)
    result = run_issue(repository, tmp_path, base_branch=branch)
    assert restarted["loops"][0]["before"] == "repair\n"
    assert result["base_sha"] == salvaged_head
    assert result["base_branch"] == branch
    assert result["success"] is False
    assert result["commits"] == []
    assert result["branch"] == ""
    assert result["verification"] is None
    assert restarted["calls"] == []
    assert git(repo, "rev-parse", branch) == salvaged_head
    assert git(repo, "rev-parse", f"{branch}:value.txt") == salvaged_blob
    assert git(repo, "rev-list", "--count", f"{base}..{branch}") == "1"
    assert git(repo, "branch", "--list", restarted["loops"][0]["branch"]) == ""
    assert_isolated(repository)
