"""Source-contract preflight only; not rollback/restart runtime acceptance."""

import json
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[2]
CONTRACT_PATH = PROJECT_ROOT / "docs/architecture/basefive-r14-contracts.json"
EXPECTED_GIT = "064a7226a500365005b5b662044fc9014a328722"
EXPECTED_BUILD_SHA = "440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd"
EXPECTED_SIGNATURE = (
    'async def implement_issue( issue: dict, repo_path: str, '
    'base_branch: str = "", artifacts_dir: str = ".artifacts", '
    'additional_context: str = "", config: dict | None = None, ) -> dict:'
)


def test_r24_declared_source_contract_and_recovery_call_owners():
    contract = json.loads(CONTRACT_PATH.read_text(encoding="utf-8"))
    assert contract["schema"] == "swe-af.basefive-r14-contracts.v1"
    assert contract["project"] == "swe-af"
    assert contract["git_receipt"] == EXPECTED_GIT
    assert contract["sources"]["swe_af/issue/build.py"] == EXPECTED_BUILD_SHA
    assert contract["api"]["pythonSignature"] == EXPECTED_SIGNATURE
    assert {"success", "commits", "verification", "base_sha", "branch"} <= set(
        contract["api"]["pythonResultFields"]
    )
    assert contract["acceptance"] == "UNACCEPTED"
    assert contract["proof_level"] == "source-contract-only"
    assert contract["dto"]["wireCompatibility"] == "UNVERIFIED"
    assert contract["dto"]["wireSchemaVersion"] is None
    assert contract["coverage"]["archsteerGoFileCount"] == 0
    assert contract["coverage"]["graphSourceFreshness"] == "UNVERIFIED"

    component = next(
        diagram for diagram in contract["diagrams"]
        if diagram["level"] == "component"
    )
    nodes = {node["id"]: node for node in component["nodes"]}
    assert nodes["entry"]["anchor"] == "async def implement_issue("
    assert nodes["impl"]["anchor"] == "async def _implement_issue_impl("
    assert nodes["checkpoint"]["anchor"] == "async def _checkpoint_if_dirty("
    assert nodes["verify"]["anchor"] == "async def _run_verification("
    assert nodes["result"]["anchor"] == "return IssueBuildResult("
    for node_id in ("entry", "impl", "checkpoint", "verify", "result"):
        assert nodes[node_id]["source"] == "swe_af/issue/build.py"
    edges = {(edge["from"], edge["to"]): edge for edge in component["edges"]}
    assert edges[("entry", "impl")]["anchor"] == (
        "return await _implement_issue_impl("
    )
    assert edges[("impl", "checkpoint")]["anchor"] == (
        "await _checkpoint_interrupted_issue("
    )
    assert edges[("impl", "verify")]["anchor"] == (
        "verification = await _run_verification("
    )
    assert ("impl", "result") in edges
