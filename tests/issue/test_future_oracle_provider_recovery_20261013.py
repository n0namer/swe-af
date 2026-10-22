from __future__ import annotations

import hashlib
import importlib
import inspect
import json
import subprocess
import sys
import urllib.error
import urllib.request
from pathlib import Path

import pytest

PROJECT_ROOT = Path(__file__).resolve().parents[2]
EXPECTED_GIT_SHA = "064a7226a500365005b5b662044fc9014a328722"
EXPECTED_FUTURE_ORACLE_SHA256 = (
    "03119b752aadff57cbc18a45bef58cf52afc7fcfcc3c44f14322425dd548b6a8"
)
FUTURE_ORACLE = PROJECT_ROOT / "tests/issue/test_releasegate_future_git_qa_20261013.py"
HEALTH_URL = "http://127.0.0.1:10110/health"


def _sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _git_head() -> str:
    completed = subprocess.run(
        ["git", "rev-parse", "HEAD"],
        cwd=PROJECT_ROOT,
        check=False,
        capture_output=True,
        text=True,
        timeout=15,
    )
    assert completed.returncode == 0, (
        f"git identity unavailable: stdout={completed.stdout!r} "
        f"stderr={completed.stderr!r}"
    )
    return completed.stdout.strip()


def _health_receipt() -> dict[str, object]:
    request = urllib.request.Request(HEALTH_URL, method="GET")
    try:
        with urllib.request.urlopen(request, timeout=10) as response:
            return {
                "url": HEALTH_URL,
                "status": response.status,
                "body_prefix": response.read(256).decode("utf-8", errors="replace"),
            }
    except urllib.error.HTTPError as error:
        return {
            "url": HEALTH_URL,
            "status": error.code,
            "body_prefix": error.read(256).decode("utf-8", errors="replace"),
        }
    except OSError as error:
        return {"url": HEALTH_URL, "status": None, "error": repr(error)}


def _entrypoint_receipt() -> dict[str, object]:
    function = importlib.import_module("swe_af.issue.build").implement_issue
    source_file = Path(inspect.getsourcefile(function) or "").resolve()
    product_root = (PROJECT_ROOT / "swe_af").resolve()
    assert source_file.is_relative_to(product_root), (
        f"implement_issue must resolve below {product_root}, got {source_file}"
    )
    signature = inspect.signature(function)
    return {
        "qualified_name": f"{function.__module__}.{function.__qualname__}",
        "signature": str(signature),
        "parameters": {
            name: {
                "kind": str(parameter.kind),
                "annotation": repr(parameter.annotation),
                "default": repr(parameter.default),
            }
            for name, parameter in signature.parameters.items()
        },
        "return_annotation": repr(signature.return_annotation),
        "source_file": str(source_file.relative_to(PROJECT_ROOT)),
        "source_sha256": _sha256(source_file),
    }


def _require_future_oracle(expected_sha256: str) -> str:
    assert FUTURE_ORACLE.is_file(), f"future oracle missing: {FUTURE_ORACLE}"
    actual_sha256 = _sha256(FUTURE_ORACLE)
    assert actual_sha256 == expected_sha256, (
        "future oracle source drift: "
        f"expected={expected_sha256} actual={actual_sha256}"
    )
    return actual_sha256


def _run_future_oracle(label: str) -> dict[str, object]:
    command = [sys.executable, "-m", "pytest", "-q", "-s", str(FUTURE_ORACLE)]
    completed = subprocess.run(
        command,
        cwd=PROJECT_ROOT,
        check=False,
        capture_output=True,
        text=True,
        timeout=180,
    )
    combined_output = f"{completed.stdout}\n{completed.stderr}".lower()
    receipt = {
        "label": label,
        "command": command,
        "exit_code": completed.returncode,
        "stdout": completed.stdout,
        "stderr": completed.stderr,
    }
    assert completed.returncode == 0, json.dumps(receipt, indent=2, sort_keys=True)
    assert "passed" in combined_output, json.dumps(receipt, indent=2, sort_keys=True)
    for forbidden in (" skipped", " xfailed", " xpassed"):
        assert forbidden not in combined_output, json.dumps(
            receipt, indent=2, sort_keys=True
        )
    return receipt


def test_provider_health_is_recorded_but_not_treated_as_native_pass() -> None:
    health = _health_receipt()
    assert health["status"] == 200, json.dumps(health, indent=2, sort_keys=True)

    with pytest.raises(AssertionError, match="future oracle source drift"):
        _require_future_oracle("0" * 64)

    receipt = {
        "health_only_is_insufficient": True,
        "health": health,
        "negative_gate": "deliberate source drift rejected",
    }
    print("PROVIDER_RECOVERY_HEALTH_RECEIPT=" + json.dumps(receipt, sort_keys=True))


def test_native_future_oracle_positive_and_restart_replay() -> None:
    git_sha = _git_head()
    assert git_sha == EXPECTED_GIT_SHA, (
        f"git source drift: expected={EXPECTED_GIT_SHA} actual={git_sha}"
    )
    oracle_sha256 = _require_future_oracle(EXPECTED_FUTURE_ORACLE_SHA256)
    entrypoint = _entrypoint_receipt()

    first_run = _run_future_oracle("positive")
    replay_run = _run_future_oracle("restart-replay")

    receipt = {
        "git_sha": git_sha,
        "future_oracle_sha256": oracle_sha256,
        "entrypoint": entrypoint,
        "positive": first_run,
        "restart_replay": replay_run,
    }
    print("PROVIDER_RECOVERY_NATIVE_RECEIPT=" + json.dumps(receipt, sort_keys=True))
