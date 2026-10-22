from __future__ import annotations

import pytest

from swe_af.runtime.providers import (
    runtime_to_harness_adapter,
    runtime_to_harness_provider,
)


def test_runtime_to_harness_adapter_supports_codex() -> None:
    assert runtime_to_harness_adapter("codex") == "codex"


@pytest.mark.parametrize("alias", ["claude", "claude-code", " claude_code "])
def test_runtime_provider_normalizes_claude_aliases(alias: str) -> None:
    assert runtime_to_harness_provider(alias) == "claude"


def test_runtime_provider_rejects_unknown_runtime() -> None:
    with pytest.raises(ValueError, match="Unsupported runtime provider"):
        runtime_to_harness_provider("unknown")


def test_execution_agents_source_uses_shared_runtime_adapter() -> None:
    import inspect
    import swe_af.reasoners.execution_agents as execution_agents

    source = inspect.getsource(execution_agents)
    assert "runtime_to_harness_adapter" in source
    assert 'provider = "claude-code" if ai_provider == "claude" else ai_provider' not in source


def test_pipeline_source_uses_shared_runtime_adapter() -> None:
    import inspect
    import swe_af.reasoners.pipeline as pipeline

    source = inspect.getsource(pipeline)
    assert "runtime_to_harness_adapter" in source
    assert 'provider = "claude-code" if ai_provider == "claude" else ai_provider' not in source


def test_fast_planner_source_uses_shared_runtime_adapter() -> None:
    import inspect
    import swe_af.fast.planner as planner

    source = inspect.getsource(planner)
    assert "runtime_to_harness_adapter" in source
    assert 'provider = "claude-code" if ai_provider == "claude" else ai_provider' not in source
