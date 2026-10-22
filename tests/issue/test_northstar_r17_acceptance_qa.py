"""Independent R17 QA; fixture dispatch is not a native task/result receipt."""

import asyncio
import hashlib
import importlib.util
import json
from pathlib import Path
import subprocess
import sys

import pytest


ROOT = Path(__file__).resolve().parents[2]
R16_PATH = ROOT / "tests/issue/test_northstar_r16_git_cycle.py"
R16_SPEC = importlib.util.spec_from_file_location("r17_r16_fixture", R16_PATH)
r16 = importlib.util.module_from_spec(R16_SPEC)
R16_SPEC.loader.exec_module(r16)
repository = r16.repository

PINNED_SHA256 = {
    "swe_af/issue/build.py": "440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd",
    "go/internal/orch/build.go": "05c39a2e3b41869b01a1482543c9875630a0777fbea24ea706c05727fe640cbc",
    "tests/issue/test_northstar_r16_git_cycle.py": "3a550d5d87bd68bb71d70a8d7a86cde97687726f06a6a651f3660347fbb1e648",
    "go/internal/orch/northstar_r16_resume_canary_test.go": "8356d7644b233c422eed1bd215d7bcfe8e0c315e37dc10132cd9d9f1a0c1182c",
    "docs/architecture/basefive-r14-contracts.json": "cc0eeb34cbc3031785787caf1e72ce8393551533376d46403b71a393c3330029",
}


def test_source_test_identity_and_c4_anchors():
    for relative, expected in PINNED_SHA256.items():
        assert hashlib.sha256((ROOT / relative).read_bytes()).hexdigest() == expected, relative
    contract = json.loads((ROOT / "docs/architecture/basefive-r14-contracts.json").read_text(encoding="utf-8"))
    sources = {
        relative: (ROOT / relative).read_text(encoding="utf-8")
        for relative in ("swe_af/issue/build.py", "go/internal/orch/build.go")
    }
    for relative in sources:
        assert contract["sources"][relative] == PINNED_SHA256[relative]
    checked = []
    missing = []

    def check_anchors(value):
        if isinstance(value, dict):
            source = value.get("source")
            if source in sources and "anchor" in value:
                if value["anchor"] not in sources[source]:
                    missing.append((source, value["anchor"]))
                checked.append((source, value["anchor"]))
            for child in value.values():
                check_anchors(child)
        elif isinstance(value, list):
            for child in value:
                check_anchors(child)

    check_anchors(contract)
    assert not missing, f"C4 anchors missing from pinned sources: {missing!r}"
    assert len(checked) >= 20
    assert {source for source, _ in checked} == set(sources)
    assert contract["dto"]["wireCompatibility"] == "UNVERIFIED"
    assert contract["coverage"]["graphSourceFreshness"] == "UNVERIFIED"
    assert contract["coverage"]["archsteerGoFileCount"] == 0
    assert contract["acceptance"] == "UNACCEPTED"


def install_verdict(monkeypatch, verdict):
    observations = r16.configure(monkeypatch)

    async def call(target, **kwargs):
        observations["calls"].append(target)
        assert target == "fixture.run_verifier"
        assert Path(kwargs["repo_path"], "value.txt").read_text(encoding="utf-8") == "repair\n"
        return verdict

    monkeypatch.setattr(r16.build.issue_router.app, "call", call)
    return observations


def run_local_verifier(value_path):
    return subprocess.run(
        [sys.executable, "-I", "-c",
         "from pathlib import Path; import sys; "
         "sys.exit(0 if Path(sys.argv[1]).read_bytes() == sys.argv[2].encode('utf-8') else 1)",
         str(value_path), "repair\n"],
        capture_output=True, text=True, timeout=30, check=False,
    )


@pytest.mark.parametrize("content, expected_status", [
    (b"repair\n", 0),
    (b"repair\\n", 1),
    (b"unrepaired\n", 1),
])
def test_local_verifier_fixture_controls(tmp_path, content, expected_status):
    value_path = tmp_path / "value.txt"
    value_path.write_bytes(content)
    process = run_local_verifier(value_path)
    assert process.returncode == expected_status, process.stderr
    assert process.stdout == ""
    assert process.stderr == ""


def test_positive_git_delivery_with_real_verifier_process(repository, tmp_path, monkeypatch):
    observations = r16.configure(monkeypatch)

    async def call(target, **kwargs):
        observations["calls"].append(target)
        assert target == "fixture.run_verifier"
        process = run_local_verifier(Path(kwargs["repo_path"], "value.txt"))
        assert process.returncode == 0, process.stderr
        return {"passed": True, "summary": "Independent local verifier process"}

    monkeypatch.setattr(r16.build.issue_router.app, "call", call)
    result = r16.run_issue(repository, tmp_path)
    assert result["success"] is True
    assert result["verification"]["passed"] is True
    assert len(observations["loops"]) == 1
    assert observations["calls"] == ["fixture.run_verifier"]
    r16.assert_delivery(result, repository)


@pytest.mark.parametrize("verdict", [None, [], {}, {"summary": "no verdict"}, {"passed": False}])
def test_missing_or_rejected_verdict_fails_closed(repository, tmp_path, monkeypatch, verdict):
    observations = install_verdict(monkeypatch, verdict)
    result = r16.run_issue(repository, tmp_path)
    assert result["success"] is False
    assert observations["calls"] == ["fixture.run_verifier"]
    r16.assert_delivery(result, repository)


@pytest.mark.parametrize("passed", ["false", "true", 1, [False], {"passed": False}])
def test_non_boolean_verdict_cannot_authorize_success(repository, tmp_path, monkeypatch, passed):
    observations = install_verdict(monkeypatch, {"passed": passed})
    result = r16.run_issue(repository, tmp_path)
    assert observations["calls"] == ["fixture.run_verifier"]
    r16.assert_delivery(result, repository)
    assert result["success"] is False, "Only a literal boolean True is a trusted verifier approval"


def test_interruption_then_two_restarts_preserve_one_effect(repository, tmp_path, monkeypatch):
    repo, base, _ = repository
    interrupted = r16.configure(monkeypatch, interrupt=True)
    with pytest.raises(asyncio.CancelledError, match="fixture interruption"):
        r16.run_issue(repository, tmp_path)
    assert interrupted["calls"] == []
    branch = interrupted["loops"][0]["branch"]
    head = r16.git(repo, "rev-parse", branch)
    blob = r16.git(repo, "rev-parse", f"{branch}:value.txt")
    r16.assert_isolated(repository)
    for _ in range(2):
        restarted = r16.configure(monkeypatch)
        result = r16.run_issue(repository, tmp_path, base_branch=branch)
        assert restarted["loops"][0]["before"] == "repair\n"
        assert result["base_sha"] == head
        assert result["success"] is False
        assert result["commits"] == []
        assert result["branch"] == ""
        assert result["verification"] is None
        assert restarted["calls"] == []
        assert r16.git(repo, "rev-parse", branch) == head
        assert r16.git(repo, "rev-parse", f"{branch}:value.txt") == blob
        assert r16.git(repo, "rev-list", "--count", f"{base}..{branch}") == "1"
        assert r16.git(repo, "branch", "--list", restarted["loops"][0]["branch"]) == ""
        r16.assert_isolated(repository)
