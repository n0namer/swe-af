"""Partial C4/API probes; these do not establish native router acceptance."""

import importlib
import inspect
import json
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[2]
CONTRACT_PATH = PROJECT_ROOT / "docs/architecture/basefive-r14-contracts.json"
SOURCE_SHA = "440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd"
GIT_RECEIPT = "064a7226a500365005b5b662044fc9014a328722"


def load_contract():
    return json.loads(CONTRACT_PATH.read_text(encoding="utf-8"))


def test_r24_c4_issue_source_binding():
    contract = load_contract()
    assert contract["schema"] == "swe-af.basefive-r14-contracts.v1"
    assert contract["project"] == "swe-af"
    assert contract["git_receipt"] == GIT_RECEIPT
    assert contract["sources"]["swe_af/issue/build.py"] == SOURCE_SHA


def test_r24_public_issue_api_matches_documented_signature():
    build = importlib.import_module("swe_af.issue.build")
    assert inspect.iscoroutinefunction(build.implement_issue)
    parameters = inspect.signature(build.implement_issue).parameters
    assert list(parameters) == [
        "issue", "repo_path", "base_branch", "artifacts_dir",
        "additional_context", "config",
    ]
    assert parameters["issue"].default is inspect.Parameter.empty
    assert parameters["repo_path"].default is inspect.Parameter.empty
    assert parameters["base_branch"].default == ""
    assert parameters["artifacts_dir"].default == ".artifacts"
    assert parameters["additional_context"].default == ""
    assert parameters["config"].default is None
    assert all(
        parameter.kind is inspect.Parameter.POSITIONAL_OR_KEYWORD
        for parameter in parameters.values()
    )


def test_r24_contract_does_not_promote_graph_evidence_to_runtime_proof():
    contract = load_contract()
    assert contract["acceptance"] == "UNACCEPTED"
    assert contract["proof_level"] == "source-contract-only"
    assert contract["dto"]["runtimeExecution"] == "UNVERIFIED"
    assert contract["dto"]["wireCompatibility"] == "UNVERIFIED"
    assert contract["dto"]["wireSchemaVersion"] is None
    assert contract["coverage"]["archsteerGoFileCount"] == 0
    assert contract["coverage"]["graphSourceFreshness"] == "UNVERIFIED"
