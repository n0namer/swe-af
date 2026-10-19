"""Regression tests for the QA synthesizer dropping the AI's decision (#113).

``run_qa_synthesizer`` uses ``router.ai()`` for supported LiteLLM
models, which returns the validated schema directly. Codex CLI model ids do
not have a LiteLLM provider prefix, so Codex must instead use
``router.harness()``, which returns a ``HarnessResult.parsed`` instance.
Both paths must preserve a model's BLOCK/FIX/APPROVE decision rather than
silently discarding it in the heuristic fallback.
"""

from __future__ import annotations

from types import SimpleNamespace
from unittest.mock import AsyncMock, MagicMock

from swe_af.execution.schemas import QASynthesisAction, QASynthesisResult
from swe_af.reasoners import execution_agents

# Inputs the heuristic fallback resolves to APPROVE, so any other action in the
# result can only have come from the synthesizer itself.
FALLBACK_APPROVES = {
    "qa_result": {"passed": True},
    "review_result": {"approved": True, "blocking": False},
}


def _error_notes(router: MagicMock) -> list:
    return [
        c for c in router.note.call_args_list if "error" in c.kwargs.get("tags", [])
    ]


async def test_ai_decision_wins_over_heuristic_fallback(monkeypatch) -> None:
    """A BLOCK from the synthesizer survives inputs the fallback would approve."""
    decision = QASynthesisResult(
        action=QASynthesisAction.BLOCK,
        summary="Tests pass but the fix regresses the public API.",
        stuck=True,
    )
    router = MagicMock(ai=AsyncMock(return_value=decision))
    monkeypatch.setattr(execution_agents, "router", router)

    out = await execution_agents.run_qa_synthesizer(
        iteration_history=[], iteration_id="iter-7", **FALLBACK_APPROVES
    )

    assert router.ai.await_args.kwargs["schema"] is QASynthesisResult
    assert out["action"] == QASynthesisAction.BLOCK
    assert out["summary"] == decision.summary
    assert out["stuck"] is True
    assert out["iteration_id"] == "iter-7"
    assert _error_notes(router) == []


async def test_non_schema_response_falls_back_and_says_so(monkeypatch) -> None:
    """An unparseable response reaches the heuristic, but not silently.

    ``ai()`` hands back a ``ToolCallResponse`` when the model's content is not
    parseable JSON.  Issue #113 was diagnosed from the note this path emits, so
    the fallback must stay loud.
    """
    router = MagicMock(ai=AsyncMock(return_value=object()))
    monkeypatch.setattr(execution_agents, "router", router)

    out = await execution_agents.run_qa_synthesizer(
        iteration_history=[], iteration_id="iter-8", **FALLBACK_APPROVES
    )

    assert out["action"] == QASynthesisAction.APPROVE
    assert _error_notes(router), "operators lose their only signal for this path"


async def test_codex_synthesis_routes_to_structured_harness(monkeypatch, tmp_path) -> None:
    """Codex CLI model names must never be passed to the LiteLLM-only ai path."""
    decision = QASynthesisResult(
        action=QASynthesisAction.BLOCK,
        summary="Reviewer found an API regression.",
        stuck=True,
    )
    router = MagicMock()
    router.ai = AsyncMock(side_effect=AssertionError("unqualified codex model reached LiteLLM"))
    router.harness = AsyncMock(return_value=SimpleNamespace(parsed=decision))
    monkeypatch.setattr(execution_agents, "router", router)

    out = await execution_agents.run_qa_synthesizer(
        iteration_history=[],
        iteration_id="codex-1",
        worktree_path=str(tmp_path),
        model="gpt-6.1-sol",
        ai_provider="codex",
        **FALLBACK_APPROVES,
    )

    assert out["action"] == QASynthesisAction.BLOCK
    assert out["summary"] == decision.summary
    assert out["stuck"] is True
    assert out["iteration_id"] == "codex-1"
    assert router.ai.await_count == 0
    call = router.harness.await_args.kwargs
    assert call["provider"] == "codex"
    assert call["model"] == "gpt-6.1-sol"
    assert call["schema"] is QASynthesisResult
    assert call["cwd"] == str(tmp_path)
    assert _error_notes(router) == []
