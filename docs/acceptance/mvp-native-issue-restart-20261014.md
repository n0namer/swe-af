# MVP native issue disposable repair — 2026-10-14

## Scope

This bounded repair hardens `swe_af.issue.build._implement_issue_impl` for cancellation during the coding loop. Ordinary failures already checkpointed tracked work; cancellation previously bypassed `except Exception`, so the `finally` cleanup could remove the disposable worktree and delete the branch before repair work was durable.

## Implemented behavior

- Positive repair writes one isolated commit and leaves the caller branch unchanged.
- Branch identity collision raises `GitOpsError` without creating a second effect.
- Restarting from a checkpoint commit sees the existing repair and creates no duplicate commit.
- `asyncio.CancelledError` checkpoints tracked work, re-raises cancellation, and does not enter verifier or PR paths.
- Checkpoint failures are recorded without masking the original interruption.

## Source-bound changes

- `swe_af/issue/build.py`
- `tests/issue/test_mvp_native_issue_restart_20261014.py`

## Validation receipt

Focused command required by the task:

```text
python -m pytest -q -p no:cacheprovider tests/issue/test_mvp_native_issue_restart_20261014.py
```

The authorized GO Wave runner for this attempt exposes Node/Vitest/tsx execution, not Python/pytest. Therefore the Python command is **BLOCKED_INFRA** unless a Python-capable runner is provided; no PASS is inferred. The test file is source-bound and covers positive, wrong-identity, restart, and cancellation-checkpoint cases.

Affected baseline requested by the task:

```text
python -m pytest -q -p no:cacheprovider tests/issue/test_northstar_independent_gitreplay_oracle_20261014.py
```

This attempt does not claim execution for the same runner limitation. Independent QA must run both commands with the project Python environment and verify the final source SHA before acceptance.

## Acceptance status

Candidate only. Product acceptance remains `UNACCEPTED`; no `ACCEPTED_GREEN` claim is made. The independent reviewer must verify exact DTO/revision compatibility and runtime cancellation/restart behavior.
