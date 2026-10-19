"""Long Codex harness runs must keep their execution record active."""
import asyncio

import pytest

from swe_af.harness_activity import run_with_activity_heartbeat


def test_live_harness_emits_notes_and_stops_on_completion():
    async def scenario():
        notes = []
        async def work():
            await asyncio.sleep(0.065)
            return "coded"
        result = await run_with_activity_heartbeat(
            work(), lambda message, tags: notes.append((message, tags)),
            label="coder", interval=0.01,
        )
        seen = len(notes)
        await asyncio.sleep(0.04)
        assert result == "coded"
        assert seen >= 2
        assert len(notes) == seen
    asyncio.run(scenario())


def test_failed_harness_cleans_up_heartbeat():
    async def scenario():
        notes = []
        async def fail():
            await asyncio.sleep(0.03)
            raise RuntimeError("provider")
        with pytest.raises(RuntimeError, match="provider"):
            await run_with_activity_heartbeat(
                fail(), lambda message, tags: notes.append(message),
                label="pm", interval=0.01,
            )
        count = len(notes)
        await asyncio.sleep(0.035)
        assert len(notes) == count
    asyncio.run(scenario())


def test_note_failure_does_not_abort_coder():
    async def scenario():
        def note(message, tags):
            raise OSError("control plane briefly offline")
        return await run_with_activity_heartbeat(
            asyncio.sleep(0.03, result="saved"), note, label="coder", interval=0.01,
        )
    assert asyncio.run(scenario()) == "saved"
