"""Utility for unwrapping execution envelopes from app.call() results.

The AgentField SDK's ``Agent.call()`` has two execution paths:

- **Async path**: correctly unwraps the result and raises on failure.
- **Sync fallback path**: may return the full execution envelope when the
  inner ``result`` is ``None`` (e.g. on failed executions).

This module provides a single ``unwrap_call_result`` helper that normalises
both cases so pipeline code can always expect the raw reasoner output.
"""

from __future__ import annotations

from swe_af.execution.fatal_error import FatalHarnessError, is_fatal_error

# Keys present in the execution envelope returned by _build_execute_response.
_ENVELOPE_KEYS = frozenset({
    "execution_id", "run_id", "node_id", "type", "target",
    "status", "duration_ms", "timestamp", "result",
    "error_message", "cost",
    "capability_id", "definition_id", "source_revision",
    "parent_id", "root_id", "thread_id",
    "parent_execution_id", "root_workflow_id", "provenance", "correlation",
})

_BINDING_KEYS = ("capability_id", "definition_id", "source_revision",
                 "parent_id", "root_id", "thread_id")


def binding_from(value: dict | None) -> dict[str, str]:
    """Return the canonical optional G8 binding from a flat or nested map."""
    value = value if isinstance(value, dict) else {}
    out = {key: str(value.get(key) or "").strip() for key in _BINDING_KEYS}
    out["parent_id"] = out["parent_id"] or str(value.get("parent_execution_id") or "").strip()
    out["root_id"] = out["root_id"] or str(value.get("root_workflow_id") or "").strip()
    for section, keys in (("provenance", _BINDING_KEYS[:3]), ("correlation", _BINDING_KEYS[3:])):
        nested = value.get(section)
        if isinstance(nested, dict):
            for key in keys:
                out[key] = out[key] or str(nested.get(key) or "").strip()
    return out


def has_binding(value: dict | None) -> bool:
    value = value if isinstance(value, dict) else {}
    return any(key in value for key in _BINDING_KEYS + ("parent_execution_id", "root_workflow_id", "provenance", "correlation"))


def validate_binding(actual: dict | None, expected: dict | None) -> None:
    """Raise ValueError when supplied provenance/correlation values disagree."""
    if not has_binding(expected):
        return
    got, want = binding_from(actual), binding_from(expected)
    for key in _BINDING_KEYS:
        if want[key] and got[key] != want[key]:
            raise ValueError(
                f"provenance/correlation mismatch for {key}: got {got[key]!r}, want {want[key]!r}"
            )


def unwrap_call_result(result, label: str = "call", expected_binding: dict | None = None):
    """Extract the actual reasoner output from an ``app.call()`` response.

    Parameters
    ----------
    result:
        The value returned by ``app.call()`` (or the equivalent ``call_fn``).
    label:
        Human-readable label for error messages (e.g. the target name).

    Returns
    -------
    dict | Any
        The inner reasoner output, free of envelope metadata.

    Raises
    ------
    RuntimeError
        If the envelope indicates a terminal failure status.
    """
    if not isinstance(result, dict):
        return result

    # Callers that know the expected execution lineage can validate it without
    # changing the legacy two-argument API.
    if expected_binding is not None:
        validate_binding(result, expected_binding)

    # Fast path: already unwrapped (no envelope keys present)
    if not _ENVELOPE_KEYS.intersection(result):
        return result

    # Looks like the execution envelope — check for errors first
    status = str(result.get("status", "")).lower()
    if status in ("failed", "error", "cancelled", "timeout"):
        err = result.get("error_message") or result.get("error") or "unknown"
        if is_fatal_error(str(err)):
            raise FatalHarnessError(str(err))
        raise RuntimeError(f"{label} failed (status={status}): {err}")

    inner = result.get("result")
    if inner is not None:
        return inner

    # Envelope present but result is None and status isn't a known failure —
    # return as-is (caller should validate).
    return result
