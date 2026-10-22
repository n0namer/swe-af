"""Independent, fail-closed oracle for the native ReleaseGate Git scenario.

This test intentionally delegates the end-to-end request, rejection, interruption,
persistence, restart, retry, isolation, and cleanup exercise to the source-bound
native scenario reviewed by the C4 reducer. It never substitutes mocks or an
in-memory fixture for that scenario.
"""

from __future__ import annotations

import hashlib
import importlib
import inspect
import re
import subprocess
from pathlib import Path

import pytest


PROJECT_ROOT = Path(__file__).resolve().parents[2]
NATIVE_TEST = PROJECT_ROOT / "tests" / "issue" / "test_releasegate_git_20261013.py"
EXPECTED_NATIVE_TEST_SHA256 = (
    "8e9c356a188e43f003286fc07d85b4d4fcc1e6b252f0faca6b02ef2956adf0eb"
)
PUBLIC_MODULE = "swe_af.issue.build"
PUBLIC_ENTRYPOINT = "implement_issue"


def _sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as source:
        for chunk in iter(lambda: source.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def _project_python() -> Path:
    candidates = (
        PROJECT_ROOT / ".venv" / "Scripts" / "python.exe",
        PROJECT_ROOT / ".venv" / "bin" / "python",
    )
    for candidate in candidates:
        if candidate.is_file():
            return candidate
    pytest.fail(
        "BLOCKED[PYTHON_VENV_UNAVAILABLE]: regular project .venv Python is required; "
        "do not substitute the worker runtime or report GREEN",
        pytrace=False,
    )


def _assert_source_bound_native_scenario() -> None:
    assert NATIVE_TEST.is_file(), (
        "BLOCKED[NATIVE_SCENARIO_MISSING]: expected source-backed native scenario at "
        f"{NATIVE_TEST.relative_to(PROJECT_ROOT)}"
    )
    actual_sha256 = _sha256(NATIVE_TEST)
    assert actual_sha256 == EXPECTED_NATIVE_TEST_SHA256, (
        "BLOCKED[NATIVE_SCENARIO_DRIFT]: independently review the changed native "
        f"scenario before updating this oracle; expected={EXPECTED_NATIVE_TEST_SHA256} "
        f"actual={actual_sha256}"
    )


def _assert_real_public_boundary() -> str:
    try:
        module = importlib.import_module(PUBLIC_MODULE)
        entrypoint = getattr(module, PUBLIC_ENTRYPOINT)
    except (ImportError, AttributeError) as error:
        pytest.fail(
            f"FAIL[C4_BOUNDARY_UNAVAILABLE]: {PUBLIC_MODULE}.{PUBLIC_ENTRYPOINT}: {error}",
            pytrace=False,
        )

    source_file = inspect.getsourcefile(entrypoint)
    assert source_file, "FAIL[UNINSPECTABLE_ENTRYPOINT]: product source file is required"
    source_path = Path(source_file).resolve()
    implementation_root = (PROJECT_ROOT / "swe_af").resolve()

    assert callable(entrypoint), f"{PUBLIC_MODULE}.{PUBLIC_ENTRYPOINT} is not callable"
    assert source_path.is_relative_to(implementation_root), (
        "FAIL[NON_PRODUCT_ENTRYPOINT]: C4 boundary must resolve to product source, "
        f"not a test/fake adapter: {source_path}"
    )
    return str(inspect.signature(entrypoint))


def _assert_real_pytest_pass(completed: subprocess.CompletedProcess[str]) -> None:
    combined = f"{completed.stdout}\n{completed.stderr}".strip()
    assert completed.returncode == 0, (
        "FAIL[NATIVE_ORACLE_FAILED]: source-backed request/rejection/restart scenario "
        f"returned {completed.returncode}\n{combined}"
    )
    assert re.search(r"\b[1-9]\d* passed\b", combined), (
        "FAIL[NO_PASS_RECEIPT]: pytest exited zero without a positive pass count\n"
        f"{combined}"
    )
    assert not re.search(r"\b(?:skipped|xfailed|xpassed)\b", combined, re.IGNORECASE), (
        "FAIL[NON_EXECUTED_PATH]: skip/xfail output cannot satisfy the native gate\n"
        f"{combined}"
    )


def test_c4_public_boundary_is_real_product_code() -> None:
    """Bind the C4 request boundary to the implementation, not an invented DTO."""

    signature = _assert_real_public_boundary()
    assert signature, "FAIL[EMPTY_SIGNATURE]: inspectable C4 callable signature required"


def test_native_request_receipt_persistence_retry_restart_oracle() -> None:
    """Execute the pinned real-Git scenario and reject blocked or synthetic GREEN."""

    _assert_source_bound_native_scenario()
    signature = _assert_real_public_boundary()
    project_python = _project_python()
    command = (
        str(project_python),
        "-m",
        "pytest",
        "-q",
        str(NATIVE_TEST.relative_to(PROJECT_ROOT)),
    )
    try:
        completed = subprocess.run(
            command,
            cwd=PROJECT_ROOT,
            text=True,
            capture_output=True,
            timeout=300,
            check=False,
        )
    except subprocess.TimeoutExpired as error:
        pytest.fail(
            "FAIL[NATIVE_ORACLE_TIMEOUT]: nested native scenario exceeded 300 seconds; "
            f"stdout={error.stdout!r} stderr={error.stderr!r}",
            pytrace=False,
        )

    try:
        _assert_real_pytest_pass(completed)
    except AssertionError as error:
        raise AssertionError(
            f"C4 boundary {PUBLIC_MODULE}.{PUBLIC_ENTRYPOINT}{signature}\n{error}"
        ) from error
