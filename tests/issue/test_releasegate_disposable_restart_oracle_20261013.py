from __future__ import annotations

import importlib
import json
import subprocess
import sys
from pathlib import Path

import pytest


COMMAND_TIMEOUT_SECONDS = 10


def run_command(
    command: list[str],
    *,
    cwd: Path,
    timeout: float = COMMAND_TIMEOUT_SECONDS,
) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        command,
        cwd=cwd,
        check=True,
        capture_output=True,
        text=True,
        timeout=timeout,
    )


def git(repo: Path, *arguments: str) -> str:
    return run_command(["git", *arguments], cwd=repo).stdout.strip()


def write_artifact(repo: Path, value: str) -> None:
    (repo / "artifact.txt").write_text(f"{value}\n", encoding="utf-8")


def initialize_rejected_repository(root: Path) -> tuple[Path, str, str]:
    repo = root / "releasegate-repo"
    repo.mkdir()
    git(repo, "init")
    git(repo, "config", "user.name", "ReleaseGate Test")
    git(repo, "config", "user.email", "releasegate@example.invalid")
    git(repo, "checkout", "-b", "main")

    write_artifact(repo, "base")
    git(repo, "add", "artifact.txt")
    git(repo, "commit", "-m", "base")
    base_sha = git(repo, "rev-parse", "HEAD")

    git(repo, "checkout", "-b", "build-review-repair")
    write_artifact(repo, "rejected")
    git(repo, "add", "artifact.txt")
    git(repo, "commit", "-m", "rejected build")
    rejected_sha = git(repo, "rev-parse", "HEAD")

    assert artifact_passes_review(repo) is False
    return repo, base_sha, rejected_sha


def artifact_passes_review(repo: Path) -> bool:
    return (repo / "artifact.txt").read_text(encoding="utf-8").strip() == "repaired"


def capture_interrupted_repair(repo: Path, checkpoint: Path) -> str:
    write_artifact(repo, "repaired")
    repair_blob = git(repo, "hash-object", "artifact.txt")
    checkpoint.write_text(
        json.dumps(
            {
                "schema": "releasegate-repair-checkpoint.v1",
                "branch": git(repo, "branch", "--show-current"),
                "rejected_head": git(repo, "rev-parse", "HEAD"),
                "repair_blob": repair_blob,
            },
            sort_keys=True,
        ),
        encoding="utf-8",
    )
    return repair_blob


def recover_interrupted_repair(repo: Path, checkpoint: Path) -> str:
    payload = json.loads(checkpoint.read_text(encoding="utf-8"))
    if payload.get("schema") != "releasegate-repair-checkpoint.v1":
        raise ValueError("unsupported repair checkpoint schema")
    if payload.get("branch") != git(repo, "branch", "--show-current"):
        raise ValueError("repair checkpoint branch does not match")

    repair_blob = payload.get("repair_blob")
    if not isinstance(repair_blob, str) or not repair_blob:
        raise ValueError("repair checkpoint has no repair blob")

    committed_blob = git(repo, "rev-parse", "HEAD:artifact.txt")
    if committed_blob == repair_blob:
        return git(repo, "rev-parse", "HEAD")

    if payload.get("rejected_head") != git(repo, "rev-parse", "HEAD"):
        raise ValueError("repair checkpoint head does not match")

    working_blob = git(repo, "hash-object", "artifact.txt")
    if working_blob != repair_blob:
        raise ValueError("interrupted repair worktree does not match checkpoint")

    git(repo, "add", "artifact.txt")
    git(repo, "commit", "-m", "repair rejected build")
    return git(repo, "rev-parse", "HEAD")


def assert_single_repair_commit(repo: Path, rejected_sha: str, repaired_sha: str) -> None:
    assert repaired_sha != rejected_sha
    assert git(repo, "rev-list", "--count", f"{rejected_sha}..{repaired_sha}") == "1"
    assert git(repo, "log", "-1", "--format=%s", repaired_sha) == "repair rejected build"
    assert artifact_passes_review(repo) is True
    assert git(repo, "status", "--porcelain") == ""


def test_public_releasegate_entrypoint_imports_without_custom_fixtures() -> None:
    module = importlib.import_module("swe_af.issue.build")
    assert callable(getattr(module, "implement_issue"))


def test_restart_salvages_interrupted_repair_once(tmp_path: Path) -> None:
    repo, base_sha, rejected_sha = initialize_rejected_repository(tmp_path)
    checkpoint = tmp_path / "repair-checkpoint.json"
    repair_blob = capture_interrupted_repair(repo, checkpoint)

    assert git(repo, "rev-parse", "HEAD") == rejected_sha
    assert git(repo, "rev-parse", "main") == base_sha
    assert artifact_passes_review(repo) is True
    assert git(repo, "status", "--porcelain") == " M artifact.txt"

    repaired_sha = recover_interrupted_repair(repo, checkpoint)
    repeated_sha = recover_interrupted_repair(repo, checkpoint)

    assert repeated_sha == repaired_sha
    assert git(repo, "rev-parse", "HEAD:artifact.txt") == repair_blob
    assert_single_repair_commit(repo, rejected_sha, repaired_sha)
    assert git(repo, "rev-parse", "main") == base_sha
    assert git(repo, "remote") == ""


def test_restart_rejects_changed_worktree_without_committing(tmp_path: Path) -> None:
    repo, base_sha, rejected_sha = initialize_rejected_repository(tmp_path)
    checkpoint = tmp_path / "repair-checkpoint.json"
    capture_interrupted_repair(repo, checkpoint)
    write_artifact(repo, "unexpected mutation")

    with pytest.raises(
        ValueError,
        match="interrupted repair worktree does not match checkpoint",
    ):
        recover_interrupted_repair(repo, checkpoint)

    assert git(repo, "rev-parse", "HEAD") == rejected_sha
    assert git(repo, "rev-list", "--count", f"{rejected_sha}..HEAD") == "0"
    assert git(repo, "rev-parse", "main") == base_sha
    assert artifact_passes_review(repo) is False


def test_timeout_then_restart_recovers_without_duplicate_commit(tmp_path: Path) -> None:
    repo, _, rejected_sha = initialize_rejected_repository(tmp_path)
    checkpoint = tmp_path / "repair-checkpoint.json"
    capture_interrupted_repair(repo, checkpoint)

    with pytest.raises(subprocess.TimeoutExpired):
        run_command(
            [sys.executable, "-c", "import time; time.sleep(2)"],
            cwd=repo,
            timeout=0.05,
        )

    assert git(repo, "rev-parse", "HEAD") == rejected_sha
    repaired_sha = recover_interrupted_repair(repo, checkpoint)
    assert recover_interrupted_repair(repo, checkpoint) == repaired_sha
    assert_single_repair_commit(repo, rejected_sha, repaired_sha)
