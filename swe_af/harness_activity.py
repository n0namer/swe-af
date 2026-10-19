"""Keep long-running native harness executions fresh in AgentField.

AgentField's existing execution-note handler updates execution.updated_at.
Node heartbeats alone do not renew the execution's stale-reaper clock.
"""
from __future__ import annotations

import asyncio
import logging
from contextlib import suppress
from typing import Awaitable, Callable, TypeVar

T = TypeVar("T")
logger = logging.getLogger(__name__)


async def run_with_activity_heartbeat(
    awaitable: Awaitable[T],
    note: Callable[..., None],
    *,
    label: str,
    interval: float = 60.0,
) -> T:
    """Run harness and renew its actual execution until terminal completion.

    A note is progress/liveness evidence only, never proof of a coding effect
    or an accepted terminal result. The heartbeat is always cancelled.
    """
    if interval <= 0:
        raise ValueError("heartbeat interval must be positive")

    async def pulse() -> None:
        while True:
            await asyncio.sleep(interval)
            try:
                note(f"{label}: harness running", tags=["harness", "heartbeat"])
            except Exception:
                logger.warning("AgentField execution heartbeat failed", exc_info=True)

    task = asyncio.create_task(pulse())
    try:
        return await awaitable
    finally:
        task.cancel()
        with suppress(asyncio.CancelledError):
            await task
