"""The unavailable independent code reviewer must never be an approval."""
import asyncio

from swe_af.reasoners import execution_agents


def test_unavailable_code_reviewer_never_approves(monkeypatch):
    class UnavailableAgent:
        def note(self, *args, **kwargs):
            pass

        async def harness(self, *args, **kwargs):
            raise RuntimeError("reviewer provider offline")

    monkeypatch.setattr(execution_agents.router, "_agent", UnavailableAgent())

    result = asyncio.run(execution_agents.run_code_reviewer(
        worktree_path=".",
        coder_result={"files_changed": ["answer.go"], "complete": True},
        issue={"name":"issue-reviewer-offline", "title":"Review actual change"},
        ai_provider="codex",
    ))
    assert result["approved"] is False
    assert result["blocking"] is True
