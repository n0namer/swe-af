---
title: 'SWE-AF native Full DAG checkpoint and safe crash recovery'
type: 'bugfix'
created: '2026-10-24'
status: 'in-progress'
review_loop_iteration: 2
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
- `ops/swe-af-local-ops/vp49-full-recovery-watch.py` — existing A55-derived watcher now tracked in SWE-AF; active VP49 Windows Task points to the synchronized copy. The original `safe_resume_admission.py`, wrapper, and regression tests are also preserved under `ops/swe-af-local-ops/`.

## Tasks & Acceptance

**Execution:**
- [x] Write the checkpoint to a synced sibling temp file and atomically publish it.
- [x] Require exact original resume identity; block missing/in-flight/unmerged checkpoint states before any worker executes.
- [x] Persist completed issue receipts before clearing the in-flight fence at the DAG barrier.
- [x] Preserve exact build ID through existing `resume_build` entrypoint.
- [x] Confirm an actual independent Python process crash and replay with each issue effect exactly once.
- [x] Exercise positive watcher admission and exactly one mocked AgentField POST against an isolated real Git/checkpoint fixture, with persisted one-shot marker.
- [x] Block dispatch when AgentField activity appears during inspection; retain ambiguous POST marker and require readback if the response lacks execution identity.
- [ ] Run controlled **AgentField + scheduled watcher** positive recovery E2E on a disposable/isolated Full run with durable post-state receipts. Do not mutate a real active run.

**Acceptance Criteria:**
- Given a completed level in a valid checkpoint, when the native executor crashes and resumes in a new process, then that level's non-idempotent effects are not repeated.
- Given a foreign, incomplete, or interrupted checkpoint, when `resume=True` is requested, then no coder is dispatched.
- Given a transport-interrupted real Full run with a stopped old writer, when the scheduled watcher admits recovery, then the exact original execution identity and final Git/test receipts remain traceable. **Not yet proven**.

## Spec Change Log

- 2026-10-24: Edge-case review added fail-closed handling for pending/unmerged Git branches; no change to normal Full planning architecture.
- 2026-10-24 watcher-positive-20261024: adversarial review found activity-to-POST race and ambiguous acceptance response. Existing A55 watcher was preserved in the SWE repo, guarded by a final activity read and durable ambiguous marker; no new scheduler or orchestration. Positive mocked API route tested against a real isolated Git repository with exact checkpoint.
- 2026-10-24 level-barrier-20261024: RED test showed checkpoint level=1 but in_flight=['second'] after merge; added atomic checkpoint after current_level advances, before next issue dispatch. Real subprocess crash at published level=1 then resume executes each issue exactly once. Isolated candidate `40a922582af4c6fece1231be734d4e7383d17262`; 120/120 targeted pytest PASS; does not change active Full's loaded source.

## Verification

**Commands:**
- `python -m pytest -q tests/test_native_dag_checkpoint_recovery_e2e.py tests/test_native_dag_subprocess_crash_recovery.py tests/test_generated_fix_contract.py tests/test_llm_shape_normalization.py tests/test_build_isolation.py tests/test_empty_build_guard.py tests/test_planner_pipeline.py tests/test_dag_executor_multi_repo.py tests/test_execute_workspace_manifest_dag_pipeline.py` — **119 PASS** on VP49.
- `git diff --check` — PASS.
- Local candidate: `C:\\Users\\Рафик\\Documents\\DEV\\universal-solver-site\\swe-af-recovery-e2e-20261024`, branch `fix/sweaf-native-dag-checkpoint-20261024`, HEAD `044ef82800de580e999ffe0a426de1816b62692a` (clean).
- Added `ops/swe-af-local-ops/test_safe_resume_admission.py` and `test_vp49_full_recovery_watch.py`: **19/19 unittest PASS**, including one positive eligible resume to mock API, new-active-job block and ambiguous POST protection; `119/119` native pytest still PASS.
- Final local Git candidate `4325e159f3e985a1533609cb3a2e60276fa7272e`, clean on VP49 in `swe-af-recovery-e2e-20261024`; predecessor `044ef82800de580e999ffe0a426de1816b62692a` retains native DAG checkpoint repairs.
- Existing scheduled `VP49-SWE-AF-Full-Recovery-Watch` active source updated with exact SHA256 match/backup and independent Windows Task Scheduler readback: result `0`, no stderr, no dispatch markers, old logically failed Full safely blocked.
- Runtime cutover for E2E: VP49 scheduled `swe-planner` PID 17808 was explicitly restarted on verified `4325e159f3e985a1533609cb3a2e60276fa7272e`, without other active executions. A new real Full started on an isolated repo as `exec_20261024_005304_qeoab5h2` / run `run_20261024_005304_8hgg2st4`. Product Manager, Architect, Tech Lead, Sprint Planner and issue writing ran; coder/reviewer entered the repair loop. The newer barrier fix `40a9225` remains an un-loaded isolated candidate; do not restart the active execution to pick it up.
- Pending positive **live** recovery E2E: after a first level has durable merge receipts and no ambiguous in-flight writer, safely interrupt/resume and prove no duplicated effects; if the window cannot be proven, do not force interruption. CrewSync 04 remains OPEN until that evidence is obtained.

## 2026-10-24 live Full canary verdict

**DoD: NOT GREEN.** Genuine `swe-planner.build` execution `exec_20261024_005304_qeoab5h2` / workflow `run_20261024_005304_8hgg2st4` ran against isolated `C:/Temp/sweaf-native-full-recovery-20261024`, baseline Git `3dcb8d3543c0b82b93d17d5a641612cea90ecea9`. Exact loaded SWE source was `4325e159f3e985a1533609cb3a2e60276fa7272e`; production/foreign working trees were not edited.

The real Full completed **transport `succeeded`, application `success=false`**. Architect/Tech Lead/PM/planner, issue writers, coder/reviewer repair and independent verifier all ran. Host independent unittest on the issue worktree passed **5/5**, but the Codex review sandbox could not launch Python or commit to the shared Git worktree index. Its native deterministic merge reported `success=true` although `merge_commit_sha == pre_merge_sha == d461812739c2b69eba9be9c182e40112ddf4da83`. The second dependent issue started without the claimed first issue code in integration HEAD. Verifier correctly refused acceptance (missing `src/outcome.py`, missing tests); finalizer returned `success=false`, left `.gitignore` dirty and removed `.artifacts/execution/checkpoint.json`. No crash was injected or resume attempted because the second issue was already in flight and no safe first-level merge boundary existed.

**Tested source repairs, isolated and NOT loaded into the completed canary:**
- `40a922582af4c6fece1231be734d4e7383d17262` — checkpoint at the completed/merged DAG level barrier. RED reproduced level 1 marked with second issue already in flight; **120/120 pytest PASS** including `os._exit(77)` and exact-once local replay.
- `ec4797a86c01497b5e0e81473fc85cd9246aee9e` — fail-closed on Git no-op merges and missing issue branch receipts. Actual Git RED reproduced the false-green merge; **101/101 pytest PASS** including real positive merges and revised unit fixtures.

**Next minimum DoD:** reuse existing host-owned deterministic Git commit path to commit the issue worktree's intended changes before merge, rather than asking Codex sandbox to write a shared parent Git index; prove actual integration HEAD includes required files and exact test SHA. Then run one controlled real Full interruption/resume with saved checkpoint and no duplicate effects, followed by independent runtime and Git readback. Avoid new orchestration, PR, deployment and CI until GREEN. Keep CrewSync `sweaf:mvp:04-recovery` OPEN and `05-full-build` BLOCKED.

## 2026-10-24 native Git handoff repair

**Framework defect:** In the first real Full canary, Codex could write `src/outcome.py` and tests to its issue worktree, but `workspace-write` sandbox denied Git `index.lock` in the parent shared Git directory. A false-positive "Already up to date" merge advanced the DAG without delivering code.

**Fix:** Reuse existing `swe_af.issue.git_ops.commit_all()` inside native `_merge_level_branches`, only after checking the exact worktree is a registered, branch-matching direct child of the build-owned `.worktrees` folder. Refresh `issue_by_name` after native workspace setup so the merge uses the canonical enriched path. Keep Codex `permission_mode=workspace-write`; no new Invoker, Git agent, scheduler or global sandbox bypass.

**Exact code:** Local VP49 branch `fix/full-native-issue-git-handoff-20261024` @ `bf27224d36044b7a809f795edab344dc137aa3d3`, clean. RED tested missing native commits/foreign worktree. GREEN: 144/144 pytest, including actual two-level `run_dag` with native worktree setup, first code commit/merge, second dependent tests, and independent unittest. Git diff check PASS.

**Loaded runtime proof:** Existing `VP49-AgentField-swe-planner` restarted only with zero active runs; listener PID `21600` loaded this exact source checkout. Isolated new repo `C:/Temp/sweaf-live-native-handoff-20261024` uses native Git init and original planner-generated two-level plan. Real AgentField execution `exec_20261024_024051_od67wi3u` / run `run_20261024_024051_9elzm1iv`: first coder and reviewer passed, native merge SHA `13dd93ac9e4c703ddd07de4ca1aa42ed982dcafa` contains `src/outcome.py` and `tests/test_outcome.py`, checkpoint level=1 completed=1 merge_results=1. Host independent unittest on *exact merged SHA* PASS 4/4. Second dependent issue running with host unittest 4/4 on its pending worktree; no operator code edits.

**Next gate:** accept the terminal outcome for second issue, verify final integration SHA and exact delivered tests. Then controlled native Full recovery watch E2E, with immutable source identity. Task `sweaf:mvp:04-recovery` stays OPEN until actual interrupted AgentField recovery; `05-full-build` remains blocked in CrewSync. No PR/deploy/CI.

## 2026-10-24 exact Git delivery acceptance

Native Git handoff repair is **GREEN as a component, not full recovery GREEN**.

- Verified local framework branch `fix/full-native-issue-git-handoff-20261024` SHA `bf27224d36044b7a809f795edab344dc137aa3d3`, clean VP49 source. `144/144` targeted pytest PASS.
- Existing scheduled `VP49-AgentField-swe-planner` switched only after zero active runs, PID `21600`, exact checkout readback. Codex remained `workspace-write`.
- Real AgentField native DAG `swe-planner.execute`: `exec_20261024_024051_od67wi3u` / `run_20261024_024051_9elzm1iv`. Two dependent issues completed, two native issue Git checkpoint commits (`5bfc31e`, `d1bd895`) and two actual merge commits (`13dd93a`, `8967bf1`); no errors/in-flight, checkpoint at terminal level 2.
- Final delivered repo: `C:/Temp/sweaf-live-native-handoff-20261024`, exact Git HEAD `8967bf11f53daad66a09e479605c426c282431eb`, clean. Independent host Python unittest on this exact SHA: **4/4 PASS**. No duplicate workers, no second plan/coder execution after completion.
- The component repair uses `git_ops.commit_all` from native SWE-AF, after verifying a registered direct-child build-owned Git worktree and matching branch. No global sandbox bypass, manual task-product fix, PR, deploy, or CI.
- A separate **terminal** checkpoint resume attempt was inconclusive: its one-shot local script did not produce a receipt or new visible AgentField execution and was terminated after readback. It was not retried. This is not a successful crash recovery test.
- **CrewSync `sweaf:mvp:04-recovery` remains OPEN**, and `05-full-build` remains dependent on 04. Next required DoD: reproducible interrupted real Full → A55-derived scheduled watcher admission → safe resume on the same exact loaded SHA, plus independent final Git/effect oracle. `execute` E2E is issue-level native evidence; it does not count as the 3/3 full `build` production streak.
