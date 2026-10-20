"""Protect long native Full planning from AgentField's stale-execution reaper."""
import ast
import asyncio
import importlib
import inspect

from swe_af import harness_activity


def test_every_planning_harness_is_heartbeat_wrapped():
    module = importlib.import_module("swe_af.reasoners.pipeline")
    tree = ast.parse(inspect.getsource(module))
    for name in ["run_environment_scout", "run_architect",
                 "run_tech_lead", "run_sprint_planner"]:
        fn = next(
            node for node in ast.walk(tree)
            if isinstance(node, ast.AsyncFunctionDef) and node.name == name
        )
        heartbeat = [
            call for call in ast.walk(fn)
            if isinstance(call, ast.Call)
            and isinstance(call.func, ast.Name)
            and call.func.id == "run_with_activity_heartbeat"
        ]
        assert heartbeat, name
        assert any(
            call.args and isinstance(call.args[0], ast.Call)
            and isinstance(call.args[0].func, ast.Attribute)
            and call.args[0].func.attr == "harness"
            for call in heartbeat
        ), name


def test_parent_wait_heartbeat_renews_and_stops(monkeypatch):
    mod = importlib.import_module("swe_af.app")
    notes = []

    async def fake_call(target, **kwargs):
        assert target == "child.run_architect"
        await asyncio.sleep(0.06)
        return {"completed": True}

    def note(message, tags):
        notes.append((message, tuple(tags)))

    native = harness_activity.run_with_activity_heartbeat

    async def accelerated(awaitable, note_fn, *, label):
        return await native(awaitable, note_fn, label=label, interval=0.01)

    monkeypatch.setattr(mod, "_original_call", fake_call)
    monkeypatch.setattr(mod.app, "note", note)
    monkeypatch.setattr(harness_activity, "run_with_activity_heartbeat", accelerated)

    async def scenario():
        result = await mod._call_with_execution_heartbeat("child.run_architect")
        count = len(notes)
        await asyncio.sleep(0.04)
        return result, count

    result, count = asyncio.run(scenario())
    assert result == {"completed": True}
    assert count >= 3
    assert len(notes) == count
    assert all("parent_wait:child.run_architect" in message for message, _ in notes)


def test_parent_wait_failure_does_not_leave_heartbeat(monkeypatch):
    mod = importlib.import_module("swe_af.app")
    notes = []

    async def failed(target, **kwargs):
        await asyncio.sleep(0.035)
        raise ValueError("child failure")

    native = harness_activity.run_with_activity_heartbeat

    async def accelerated(awaitable, note_fn, *, label):
        return await native(awaitable, note_fn, label=label, interval=0.01)

    monkeypatch.setattr(mod, "_original_call", failed)
    monkeypatch.setattr(mod.app, "note", lambda message, tags: notes.append(message))
    monkeypatch.setattr(harness_activity, "run_with_activity_heartbeat", accelerated)

    async def scenario():
        try:
            await mod._call_with_execution_heartbeat("child.run_architect")
        except ValueError:
            pass
        count = len(notes)
        await asyncio.sleep(0.04)
        return count

    count = asyncio.run(scenario())
    assert count >= 1
    assert len(notes) == count
