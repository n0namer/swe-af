---
title: 'Autonomous Full Build with two dependent real issues'
type: 'feature'
created: '2026-10-24'
status: 'in-progress'
review_loop_iteration: 0
context:
  - 'AGENTS.md'
  - '_bmad-output/planning-artifacts/epics-sweaf-mvp.md'
  - '_bmad-output/implementation-artifacts/spec-1-4-recovery-invariants.md'
---

<frozen-after-approval reason="human-owned intent — do not modify unless human renegotiates">

## Intent

**Problem:** The SWE-AF North Star is an autonomous `swe-planner.build` that plans, codes, reviews, repairs and independently verifies a real change. A green `swe-planner.execute` with a preexisting plan is useful evidence, but is not the complete `build` acceptance. The earlier Full incorrectly counted no-op Git merges as delivered code.

**Approach:** Reuse the native SWE-AF build pipeline, existing AgentField node, Codex `workspace-write`, `git_ops.commit_all`, deterministic Git merge, checkpoint and verifier. Require real delivered integration HEAD and independent tests before declaring logical success.

## Boundaries & Constraints

**Always:** Work on a disposable local repository; verify source SHA, exact issue dependencies, original/new Git commits, clean integration branch, application `success=true`, independent executable tests, and no duplicate work. All Git commits happen via the existing native SWE-AF owner after Codex edits; no global sandbox bypass.

**Ask First:** Changes to an unrelated or dirty checkout; external PR/push/deploy; abort of another user's active execution.

**Never:** Add another orchestration layer, scheduler, custom LLM routing, extra agents, or a second parallel roadmap. Never conflate AgentField transport `succeeded` with `BuildResult.success=true`.

## I/O & Edge-Case Matrix

| Scenario | Input | Expected result | Rejection condition |
|---|---|---|---|
| Basic | Goal describing classifier and Python unittest | Planner creates two dependent issues and Codex delivers both on integration HEAD | No merge receipt, incorrect source |
| Repair | Reviewer detects missing tests/contract | Existing coder-reviewer loop fixes and rechecks | Unresolved blocking review or unverified tests |
| Sandbox Git | Codex can edit worktree but shared Git index is protected | Native SWE-AF Git owner commits and merges | Missing or foreign worktree, no-op merge |
| Acceptance | Complete build | Logical success, test oracle PASS, exact SHA/provenance | Only transport success or dirty workspace |

</frozen-after-approval>

## Code Map

- `swe_af/app.py` — existing Full Build planning, verification and logical result.
- `swe_af/execution/dag_executor.py` — native issue DAG, worktree creation, branch merge and checkpoint.
- `swe_af/execution/git_fast_path.py` — deterministic Git delivery, no-op guards and native handoff.
- `swe_af/issue/git_ops.py` — reusable host-owned `commit_all`.
- `ops/windows/SWE.ps1` — existing controlled Full submission, receipt ID and `-NoPR`.

## Tasks & Acceptance

**Execution:**
- [x] Reject completed issues whose changes were not committed to the Git issue branch; commit via existing native Git owner.
- [x] Demonstrate 2/2 dependent issues in real AgentField `execute`, each with independent Git commit/merge, canonical tests on final SHA.
- [ ] Run unmodified `swe-planner.build` from a new goal and clean disposable repository, with PM → Architect → Planner → coder/reviewer → verifier (including repair if needed).
- [ ] Independently verify all delivered files, Git receipt SHA, `python -m unittest discover -s tests -v`, clean Git and **logical** BuildResult success.

**Acceptance Criteria:**
- Given a clean repo and goal, when Full is submitted once, then two dependent issues complete autonomously, with no manual product-code edits.
- Given a reviewer objection, when the existing repair loop runs, then revised code is rechecked by an independent verifier.
- Given a native Git handoff, when integration merge is accepted, then the resulting HEAD contains actual issue code and test files.
- Given terminal Full, when independent readback is done, then `BuildResult.success=true`, executable tests PASS and source/merge SHA are recorded.

## Spec Change Log

- 2026-10-24: Existing native Git handoff `bf27224d36044b7a809f795edab344dc137aa3d3` fixed Codex sandbox/shared Git-index failure. It is validated by 144/144 targeted tests and two real AgentField `execute` issue merges.
- 2026-10-24: First autonomous Full from an entirely new goal is running as root `exec_20261024_122809_ibilsknz`, run `run_20261024_122809_fyaqc1cy`, in disposable `C:/Temp/sweaf-full-green-20261024-2` initialized at `768ecc885d5de53221d5a6579ee840aca3a061dd`; no PR. Not yet an acceptance verdict.

## Verification

- Native Git checkpoint handoff: local code commit `bf27224d36044b7a809f795edab344dc137aa3d3`, 144/144 targeted Python tests PASS.
- Real preplanned native AgentField `execute`: `exec_20261024_024051_od67wi3u`, two native issue commits, two merge receipts; final HEAD `8967bf11f53daad66a09e479605c426c282431eb`, 4/4 host unittest PASS.
- Pending autonomous full: independent terminal `BuildResult.success`, final SHA, tests, repairs and provenance must be read back.
