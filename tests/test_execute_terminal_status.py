"""Transport terminal status must reflect genuine DAG completion, not just return."""
from unittest.mock import AsyncMock, patch
import pytest
from swe_af.app import execute, ReasonerFailed
from swe_af.execution.schemas import DAGState, IssueResult, IssueOutcome


@pytest.mark.asyncio
async def test_execute_all_failed_is_terminal_failure_with_receipt():
    state = DAGState(
        repo_path="R:/fixture",
        all_issues=[{"name": "broken"}],
        levels=[["broken"]],
        failed_issues=[IssueResult(issue_name="broken", outcome=IssueOutcome.FAILED_UNRECOVERABLE)],
    )
    with patch("swe_af.execution.dag_executor.run_dag", new=AsyncMock(return_value=state)):
        with pytest.raises(ReasonerFailed) as exc:
            await execute(plan_result={}, repo_path="R:/fixture")
    assert exc.value.result["failed_issues"][0]["issue_name"] == "broken"


@pytest.mark.asyncio
async def test_execute_partial_failure_is_not_green():
    state = DAGState(
        repo_path="R:/fixture",
        all_issues=[{"name": "done"}, {"name": "blocked"}],
        completed_issues=[IssueResult(issue_name="done", outcome=IssueOutcome.COMPLETED)],
        failed_issues=[IssueResult(issue_name="blocked", outcome=IssueOutcome.FAILED_UNRECOVERABLE)],
    )
    with patch("swe_af.execution.dag_executor.run_dag", new=AsyncMock(return_value=state)):
        with pytest.raises(ReasonerFailed) as exc:
            await execute(plan_result={}, repo_path="R:/fixture")
    assert len(exc.value.result["completed_issues"]) == 1


@pytest.mark.asyncio
async def test_execute_all_completed_returns_normal_result():
    state = DAGState(
        repo_path="R:/fixture",
        all_issues=[{"name": "done"}],
        completed_issues=[IssueResult(issue_name="done", outcome=IssueOutcome.COMPLETED)],
    )
    with patch("swe_af.execution.dag_executor.run_dag", new=AsyncMock(return_value=state)):
        result = await execute(plan_result={}, repo_path="R:/fixture")
    assert result["failed_issues"] == []
    assert len(result["completed_issues"]) == 1
