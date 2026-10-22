"""Real native _execute_level concurrency limiter with fake non-LLM workers."""
import asyncio
import pytest
from swe_af.execution import dag_executor
from swe_af.execution.schemas import DAGState, ExecutionConfig, IssueResult, IssueOutcome


@pytest.mark.asyncio
@pytest.mark.parametrize("limit, expected_peak", [(30, 30), (10, 10), (3, 3)])
async def test_native_execution_level_fills_limit(limit, expected_peak, monkeypatch):
    running, peak, completed = 0, 0, []

    async def fake_execute(issue, dag_state, execute_fn, config, **kwargs):
        nonlocal running, peak
        running += 1
        peak = max(running, peak)
        await asyncio.sleep(0.04)
        running -= 1
        completed.append(issue["name"])
        return IssueResult(issue_name=issue["name"], outcome=IssueOutcome.COMPLETED)

    monkeypatch.setattr(dag_executor, "_execute_single_issue", fake_execute)
    tasks = [{"name": f"issue-{i:02}", "files_to_modify": [f"pkg/file{i:02}.go"]} for i in range(30)]
    state = DAGState(repo_path=".", artifacts_dir=".", all_issues=tasks)
    cfg = ExecutionConfig(runtime="codex", max_concurrent_issues=limit)
    result = await dag_executor._execute_level(
        tasks, None, state, cfg, level_index=0, call_fn=None)
    assert peak == expected_peak
    assert len(completed) == 30
    assert len(result.completed) == 30
    assert not result.failed
