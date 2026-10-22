# MVP R10 disposable Git E2E receipt

## Scope

This artifact records the bounded C4 MVP scenario for `swe_af.issue.build._implement_issue_impl`:

- one disposable Git task returns a typed RED result after rejected output;
- a second disposable task performs the bounded repair and returns GREEN;
- invalid issue payloads fail before repository mutation;
- cancellation checkpoints worktree changes without mutating the caller checkout;
- replay uses a fresh branch and does not duplicate caller-side effects;
- worktrees are removed after each task, while interrupted committed work remains inspectable on its branch.

The test is source-bound at runtime by calculating the SHA-256 of `swe_af/issue/build.py` and binding both task results to the same base commit SHA.

## Focused validation

```text
python -m pytest -q tests/issue/test_northstar_mvp_r10_disposable_git_e2e.py
```

The authorized GO Wave runner for this attempt exposes Node/Vitest/tsx execution only; Python/pytest execution is therefore recorded as infrastructure-blocked unless a Python-capable runner is supplied. This document does not promote an unexecuted check to PASS.

## Acceptance mapping

| Requirement | Evidence in focused test |
| --- | --- |
| C4 DTO and immutable file SHA | runtime `IssueBuildResult` dictionaries, base commit SHA, and `build.py` SHA-256 |
| Real positive result | RED invocation followed by successful bounded repair |
| Two fail-closed negatives | invalid payload and cancellation preserve the caller checkout |
| Crash/replay safety | cancellation checkpoint plus fresh-branch replay with one commit and no caller mutation |
| No hidden effects | caller branch, `HEAD`, status, and worktree cleanup assertions |

## Implementation note

Interruption checkpointing now commits only when the scrubbed worktree is dirty. Replaying a task against an already-clean worktree therefore does not create an empty duplicate checkpoint commit, while the original interruption cause remains propagated to the caller.
