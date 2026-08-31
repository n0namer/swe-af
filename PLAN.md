# SWE-AF Project Plan (Source of Truth)

Status: active
Last reconciled: 2026-08-31
Canonical project SoT: this file (`PLAN.md`)
Canonical operator lifecycle: `n0namer/universal-solver/docs/runbooks/agentfield-dev-debug-test-handoff.md`
BMAD method source: `n0namer/BMAD,MNNZ h.agents/skills/bmad-help/SKILL.md`

## North Star

Deliver SWE-AF as a reliable AgentField software-engineering capability that can take a bounded real-repository task, produce the correct change, pass canonical tests/acceptance checks, recover from expected failures, and produce a durable exact-SHA Git result with complete provenance.

For well-scoped work, `swe-planner.implement_issue` is the preferred acceptance entrypoint. `swe-planner.build` is reserved for feature-level decomposition/execution after lower gates pass.

## Operating Invariants (Anti-Drift)

- Universal Solver runbook owns the permanent AgentField DEV development lifecycle; SWE-AF docs own SWE-AF product architecture.
- Lane A = fast live-container debug: observe current source → bounded live patch → targeted test → same-target reload only if needed → functional canary → iterate.
- Lane B = capture/durability/promotion: only accepted runtime deltas become durable `fork/dev` state and later exact-SHA materialization.
- Do NOT use GitHub-edit → redeploy as a normal inner debug loop.
- Runtime may temporarily be ahead of Git during debugging; release/handoff may not.
- Effective live-source identity = Git HEAD + working-tree delta + loaded process/container generation.
- SourceLoop/runtime-capture is a stale-safe proposal/delivery layer, not an automatic filesystem-to-main synchronizer.
- Generated files, caches, logs, test artifacts and temporary instrumentation are not candidate source.
- `fork/main` should remain a clean upstream mirror; downstream project deltas belong on `fork/dev` or their canonical operational owner.

## Current Factual State

Cast at 2026-08-31 from live readback.

- Permanent AgentField DEV Coolify application: `edshqtkwskg3lrczekhcmd71`; repo `n0namer/universal-solver`; deployed orchestrator SHA `75652b4b1f0bf18dbcdd6af9abfef40bfa068cd7`.
- Workforce container: `workforce-edshqtkwskg3lrczekhcmd71-184945561237`; running and healthy; restart_count `3`? Note: container inspection readback showed `restart_count=0`, while Coolify application metadata reports `max_restart_count=10` and last restart type `crash` on 2026-08-30. Use Docker container readback for current container restart_count, Coolify for app history.
- `/src` is a writable persistent runtime-source volume in workforce.
- Accepted SWE-AF runtime seed in that generation: `da9228f6dcaeffa2aca3cf781f04d2ea720b5294` (`swe-af:dev` at capture acceptance).
- SourceLoop swee-canary has already passed at least once: runtime → capture → Git, and the accepted SWE dev SHA is recorded in Universal Solver fleet lock.
- `runtime-capture` is currently running and has recent successful capture activity for Deep Research.
- `meta_deep_research` node is currently `active/ready`.
- `swe-planner` is currently `inactive/offline`; last heartbeat 2026-08-25T13:30:41Z.
- `swe-pro` is currently `inactive`; last heartbeat 2026-08-25T13:30:41Z.
- Persisted `swe-planner.log` shows that `swe-planner` did successfully start as a background process (PID 6119, port 8005) and registered with AgentField, but it is no longer running.
- Current workforce container logs show only the provenance HTTP server activity, not live SYE process output.
- Coolify env inventory shows both `OPENROUTER_API_KEY` and `ANTHROPIC_API_KEY` are configured as this application runtime keys; values are intentionally not exposed here.
- Current AgentField registry still exposes the `build`, `plan`, `implement_issue`, `resolve` and internal reasoner schemas for `swe-planner`, but they are not callable as a live SWE capability while the node is offline.
- `swe-af:main` and `swe-af:dev` are currently diverged: `dev` is 5 commits ahead and 1 commit behind main. The main-only downstream commit is `docs: add canonical error ledger`, conflicting with the clean-mirror invariant.

## Current Stage

BROWNFIELD RECOVERY / SEMANTIC ACCEPTANCE.

Infrastructure bootstrap, exact SHA reconciliation, discovery, and one real SWE SourceLoop capture are already proven. The next mandatory gate is to restore a live `swe-planner` node in the existing permanent DEV generation and then run the smallest discriminating `implement_issue` canary.

## Bounded Development Batches

Default batch size: ~30 minutes. Each batch closes a coherent DoD gate and writes back this file.

### Batch 1 — Restore SWE liveness and establish baseline

Status: IN PROGRESS

DoD:
1. Identify why the previously started `swe-planner` died or became unreachable.
2. Restore the same `swe-planner` process/node without whole-fleet redeploy/recreation.
3. Verify AgentField discovery shows `active/ready` for `swe-planner`.
4. Run a non-mutating reasoner smoke or minimal schema/health call and record execution_id/status.
5. Confirm `mplement_issue` schema/callability ready for the next batch.
6. Update `PLAN.md` with exact eavidence and next move.

Stop condition: don't run a real coding canary until the node is active/ready.

### Batch 2 — First bounded `implement_issue` canary

Status: PENDING

Prereduisite: Batch 1 DoD PASS.

DoD:
1. Use an existing safe/sacrificial local repo/workspace; do not create persistent infrastructure unless required.
2. Run ONE well-scoped issue with machine-checkable acceptance criteria.
3. Capture: execution_id, run/reasoner status, diff, test results, tool calls, wall time, retries, unintended files changed.
4. Require targeted and canonical repo tests to pass.
5. Verify the resulting branch/diff is bounded to the issue scope.
6. Update `PLAN.md` with result and next gate.

### Batch 3 — Failure/recovery gate

Status: PENDING

Prerequisite: Batch 2 DoD PASS.

DoD:
- nonexistent/invalid task fails closed or abstains without repo damage;
- bounded run interruption + resume is idempotent;
- stale SHA/branch advance is detected and fails closed;
- unrelated worktree files are preserved.

### Batch 4 — Durability / SourceLoop gate

Status: PENDING

Prerequisite: accepted runtime fix or accepted canary delta.

DoD:
- only source delta is captured; noise excluded;
- stale capture fails closed;
- accepted delta reaches `fork/dev`;
- exact WORKING_DEV_SHA is recorded;
- repo tests pass on that durable SHA.

### Batch 5 — Materialization and regression

DoD:
- exact accepted SHA → materialized runtime;
- functional canary repeats;
- container-only required deltas = 0;
- provenance chain complete.

## Acceptance Metrics

For each canary record at least:
- task_success;
- tests_passed;
- unintended_files_changed;
- human_interventions;
- LLM calls;
- tool calls;
- wall time;
- cost if available;
- retries;
- recovery_success;
- provenance_complete.

## Failure Classification

Assign each primary failure to exactly one class:
- `M81ER_REASONING`
P pROMPT_OR_PLANNING@
- `ACI_HARNESS`
- REPO_ENVIRONMENT`
- `GIT_DELIVERY`
KH `AGENTFIELD_RUNTIME`
KHp`PROVIDER`
- SOURCELOOP_RECONCILIATION@
- `OBSERVABILITY_GAP`

## Anti-Drift Checklist per Batch

- Reread `PLAN.md` before mutation.
- Observe CURRENT runtime before mutation.
- Verify exact target/SHA/process generation.
- Use the smallest discriminating test first.
- Reload only the same target if needed.
- Do not advance on health alone; require functional/semantic evidence.
- After acceptance, materialize exact source to Git and record SHA.
- When runtime and Git disagree, label `DESIGN_RUNTIME_DRIFT` or `FROVENANCE_GAP` as appropriate; do not silently reconcile.

## Known Drift/Debt

- [OPEN] `swe-af:main` contains a downstream `ERRORS.md` commit and therefore is not a clean upstream mirror. Resolve before the next upstream rebase/reconciliation cycle.
- [OPEN] Current non-DEV terminal runtime does not yet register the `lagentfield-dev-workforce` target even though current `vps-terminal` Git config does; the DEV terminal knows the target but mediates generic exec as review-required. This is an operator-plane capability drift, NOT a product defect.
- [OPEN] `swe-planner` process died after successful bootstrap; exact cause not yet proven.

## BMad Workflow Used for Current Batch

- Entry: `bmad-help` from `BMAD-MNNZ`.
- Classification: brownfield recovery/quick implementation.
- `bmad-correct-course` rejected as a structural fit because it requires PRD + Epics artifacts that do not exist here; creating them would be method theater.
- Selected: `bmad-quick-dev` for bounded recovery + first semantic canary.
- Write-back target: this `PLAN.md`.
