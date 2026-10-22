"""Pydantic schemas for the swe-fast single-pass build node."""

from __future__ import annotations

from typing import Literal

from pydantic import BaseModel, ConfigDict, Field
from swe_af.execution.schemas import (
    _OPENROUTER_AUTO_DEFAULT_MODEL,
    _default_runtime,
    resolve_runtime_models,
)

# Compatibility aliases for callers/tests that inspect the historical fast
# defaults; the actual values now come from the shared resolver's table.
_CLAUDE_CODE_DEFAULT = "haiku"
_OPEN_CODE_DEFAULT = _OPENROUTER_AUTO_DEFAULT_MODEL

# All four roles resolved by fast_resolve_models()
_FAST_ROLES: tuple[str, ...] = ("pm_model", "coder_model", "verifier_model", "git_model")

# Mapping from role name (in models dict) to resolved key
_ROLE_KEY_MAP: dict[str, str] = {
    "pm": "pm_model",
    "coder": "coder_model",
    "verifier": "verifier_model",
    "git": "git_model",
}


# ---------------------------------------------------------------------------
# Task-level schemas
# ---------------------------------------------------------------------------


class FastTask(BaseModel):
    """A single task in the flat fast-build decomposition."""

    model_config = ConfigDict(extra="forbid")

    name: str                           # kebab-case slug
    title: str                          # human-readable title
    description: str                    # self-contained description for the coder
    acceptance_criteria: list[str]      # task-specific acceptance criteria
    files_to_create: list[str] = []
    files_to_modify: list[str] = []
    estimated_minutes: int = 5


class FastPlanResult(BaseModel):
    """Output of the fast planner reasoner."""

    tasks: list[FastTask]
    rationale: str = ""
    fallback_used: bool = False


class FastTaskResult(BaseModel):
    """Result of executing a single FastTask."""

    task_name: str
    outcome: str                # "completed" | "failed" | "timeout"
    files_changed: list[str] = []
    summary: str = ""
    error: str = ""


class FastExecutionResult(BaseModel):
    """Aggregate result of executing all tasks."""

    task_results: list[FastTaskResult]
    completed_count: int
    failed_count: int
    timed_out: bool = False


class FastVerificationResult(BaseModel):
    """Result of the single verification pass."""

    passed: bool
    summary: str = ""
    criteria_results: list[dict] = []
    suggested_fixes: list[str] = []


# ---------------------------------------------------------------------------
# Build-level schemas
# ---------------------------------------------------------------------------

class FastBuildConfig(BaseModel):
    """Configuration for a fast single-pass build run."""

    model_config = ConfigDict(extra="forbid")

    runtime: Literal["claude_code", "open_code", "codex"] = Field(default_factory=_default_runtime)
    models: dict[str, str] | None = None
    max_tasks: int = 10
    task_timeout_seconds: int = 300
    build_timeout_seconds: int = 600
    enable_github_pr: bool = True
    github_pr_base: str = ""
    permission_mode: str = ""
    repo_url: str = ""
    agent_max_turns: int = 50


class FastBuildResult(BaseModel):
    """Top-level result returned by the fast build reasoner."""

    plan_result: dict
    execution_result: dict
    verification: dict | None = None
    success: bool
    summary: str
    pr_url: str = ""


# ---------------------------------------------------------------------------
# Model resolution helper
# ---------------------------------------------------------------------------


def fast_resolve_models(config: FastBuildConfig) -> dict[str, str]:
    """Resolve the four role model strings for a fast build run.

    Resolution order (last wins):
      1. Runtime default (haiku or the shared open_code default, per runtime)
      2. Env cascade: ``SWE_DEFAULT_MODEL`` → ``AI_MODEL`` → ``HARNESS_MODEL``
         (``HARNESS_MODEL`` only on the ``open_code`` runtime)
      3. ``models["default"]`` — overrides all roles
      4. ``models["<role>"]`` — overrides a specific role (pm, coder, verifier, git)

    Args:
        config: A :class:`FastBuildConfig` instance.

    Returns:
        A dict with keys ``pm_model``, ``coder_model``, ``verifier_model``,
        ``git_model`` mapping to model name strings.

    Raises:
        ValueError: If ``config.models`` contains a key that is not ``"default"``
            and not one of the four known role names.
    """
    # Fast mode deliberately keeps its Claude role default (haiku), while all
    # runtime, gateway, env/tier, and explicit role precedence comes from the
    # main resolver.  The base override is below every shared override layer.
    base_models = None
    if config.runtime == "claude_code":
        base_models = {role: _CLAUDE_CODE_DEFAULT for role in _FAST_ROLES}
    if config.models:
        valid_keys = {"default"} | set(_ROLE_KEY_MAP)
        unknown = sorted(set(config.models) - valid_keys)
        if unknown:
            raise ValueError(
                f"Unknown role key {unknown[0]!r} in models dict. "
                f"Valid keys are: {sorted(valid_keys)}"
            )
    return resolve_runtime_models(
        runtime=config.runtime,
        models=config.models,
        field_names=list(_FAST_ROLES),
        base_models=base_models,
    )
