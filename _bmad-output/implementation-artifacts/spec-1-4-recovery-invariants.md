---
title: 'SWE-AF native Full DAG checkpoint and safe crash recovery'
type: 'bugfix'
created: '2026-10-24'
status: 'in-progress'
review_loop_iteration: 1
context:
  - 'AGENTS.md'
  - 'PLAN.md'
  - '_bmad-output/planning-artifacts/epics-sweaf-mvp.md'
---

<frozen-after-approval reason="human-owned intent — do not modify unless human renegotiates">

## Intent

**Problem:** On VP49, the existing Full DAG could silently start new work when `resume=True` lacked a checkpoint, overwrite a valid checkpoint during an interrupted write, or clear in-flight ownership before persisting completed issue receipts. These states risk duplicated effects after a crash.

**Approach:** Retain the existing `swe-planner.execute`, `run_dag`, `resume_build`, Git/worktree logic and the A55-derived scheduled recovery watcher. Harden only the checkpoint write and resume admission, then prove an actual native process exits and restarts with no duplicated completed work.

## Boundaries & Constraints

**Always:** Keep original Git repositories and foreign dirty files unchanged; match resume repository/build/artifacts/plan identity; fail closed on absent, invalid, in-flight, unmerged, or unresolved checkpoints. Persist issue receipts before clearing the in-flight fence. Run independent executable regression evidence on exact local source.

**Ask First:** Destructive reset, cancellation of somebody else's active execution, overwriting private work, or external deploy/push/PR.

**Never:** Introduce another scheduler, invoker or orchestrator; silently replay ambiguous non-idempotent work; treat a passing unit test or transport `succeeded` as production Full acceptance.

## I/O & Edge-Case Matrix

| Scenario | Input / State | Expected Output / Behavior | Error Handling |
|---|---|---|---|
| Missing checkpoint | `resume=True`, no saved file | No worker invocation | Explicit RuntimeError |
| Interrupted checkpoint write | Existing valid checkpoint, publish fails | Earlier valid JSON remains intact | Remove temp, propagate failure |
| Mismatched run | Wrong repo/build/plan or in-flight issue | No new execution | Fail closed |
| Unmerged Git effects | Completed issue without merge receipt or unresolved branch | No unsafe replay | Manual reconciliation |
| Crashed native executor | Process dies after completed first level | Fresh process executes second only | Exact effects remain once |

</frozen-after-approval>

## Code Map

- `swe_af/execution/dag_executor.py` — existing native checkpoint writer, admission and level execution barrier.
- `swe_af/app.py` — existing `resume_build` entrypoint and build ID forwarding.
- `tests/test_native_dag_checkpoint_recovery_e2e.py` — targeted missing/malformed/identity/merge and level-resume tests.
- `tests/test_native_dag_subprocess_crash_recovery.py` — independent process `os._exit(77)` followed by fresh-process replay.
- VP49 `swe-af-local-ops/vp49-full-recovery-watch.py` — previously installed A55-derived Windows Scheduler watcher; not reimplemented.

## Tasks & Acceptance

**Execution:**
- [x] Write the checkpoint to a synced sibling temp file and atomically publish it.
- [x] Require exact original resume identity; block missing/in-flight/unmerged checkpoint states before any worker executes.
- [x] Persist completed issue receipts before clearing the in-flight fence at the DAG barrier.
- [x] Preserve exact build ID through existing `resume_build` entrypoint.
- [x] Confirm an actual independent Python process crash and replay with each issue effect exactly once.
- [ ] Run controlled **AgentField + scheduled watcher** positive recovery E2E on a disposable/isolated Full run with durable post-state receipts. Do not mutate a real active run.

**Acceptance Criteria:**
- Given a completed level in a valid checkpoint, when the native executor crashes and resumes in a new process, then that level's non-idempotent effects are not repeated.
- Given a foreign, incomplete, or interrupted checkpoint, when `resume=True` is requested, then no coder is dispatched.
- Given a transport-interrupted real Full run with a stopped old writer, when the scheduled watcher admits recovery, then the exact original execution identity and final Git/test receipts remain traceable. **Not yet proven**.

## Spec Change Log

- 2026-10-24: Edge-case review added fail-closed handling for pending/unmerged Git branches; no change to normal Full planning architecture.

## Verification

**Commands:**
- `python -m pytest -q tests/test_native_dag_checkpoint_recovery_e2e.py tests/test_native_dag_subprocess_crash_recovery.py tests/test_generated_fix_contract.py tests/test_llm_shape_normalization.py tests/test_build_isolation.py tests/test_empty_build_guard.py tests/test_planner_pipeline.py tests/test_dag_executor_multi_repo.py tests/test_execute_workspace_manifest_dag_pipeline.py` — **119 PASS** on VP49.
- `git diff --check` — PASS.
- Local candidate: `C:\\Users\\Рафик\\Documents\\DEV\\universal-solver-site\\swe-af-recovery-e2e-20261024`, branch `fix/sweaf-native-dag-checkpoint-20261024`, HEAD `044ef82800de580e999ffe0a426de1816b62692a` (clean).
- The running `swe-planner` has **not** been switched to this candidate; do not mark CrewSync 04 DONE before live recovery acceptance.
