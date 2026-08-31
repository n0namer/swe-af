# SWE-AF Project Plan (Source of Truth)

Status: active
Last reconciled: 2026-08-31
Canonical project SoT: this file (`PLAN.md`)
Canonical operator lifecycle: `n0namer/universal-solver/docs/runbooks/agentfield-dev-debug-test-handoff.md`
BMAD method source: `n0namer/BMAD-MNNZ/.agents/skills/bmad-help/SKILL.md`
Engineering contract: root `AGENTS.md` (FVE-adapted for SWE-AF; canonical engineering prerequisite as of commit `8f464a5d15afe628ff162bf6a3956c411873105f`)

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

Reconciled from live readback on 2026-08-31.

- Permanent AgentField DEV Coolify application: `edshqtkwskg3lrczekhcmd71`; repo `n0namer/universal-solver`; deployed orchestrator SHA `75652b4b1f0bf18dbcdd6af9abfef40bfa068cd7`.
- Current workforce container: `workforce-edshqtkwskg3lrczekhcmd71-184945561237`; created 2026-08-30T18:50:28Z; running and healthy; Docker restart_count = 0. Coolify app history records a crash-type restart on 2026-08-30.
- `/src` is a writable persistent runtime-source volume in workforce.
- Accepted SWE-AF runtime seed in that generation: `da9228f6dcaeffa2aca3cf781f04d2ea720b5294`.
- SourceLoop SWE canary has already passed at least once: runtime → capture → Git; accepted SWE dev SHA is recorded in Universal Solver fleet lock.
- `runtime-capture` is currently running and has recent successful capture activity for Deep Research.
- `meta_deep_research` is currently `active/ready`.
- Historical standalone SWE service `universal-solver-swe-af` (`wetscrp2tj90tklmlvkcadfw`) existed during Wave 0, but canonical server-ops records it as removed during cleanup on 2026-08-21. It is superseded topology; CURRENT target is permanent DEV `edshqtkwskg3lrczekhcmd71`.
- Historical isolated B `2zciq6hujpev6dbudcdlijqq` remains healthy but is forensic comparison only, not the current mutation target.
- `swe-planner` / `swe-pro` are not live in CURRENT permanent DEV. Historical persisted logs prove they previously started, registered and served requests, but that evidence predates the current container generation.
- Canonical Universal Solver handoff contains stronger CURRENT root-cause evidence for the present generation: workforce bootstrap resolved `OPENROUTER_API_KEY` and `ANTHROPIC_API_KEY` empty, intentionally skipped `swe-planner`, and nevertheless became healthy. The prior `PROCESS_EXIT` hypothesis is superseded; primary classification is `PROVIDER / BOOTSTRAP_GATE`.
- Coolify environment inventory proves the provider key names exist, not that secret values are non-empty. Secret values remain unexposed.
- SWE's own `go/agentfield-package.yaml` requires one of `ANTHROPIC_API_KEY` or `OPENROUTER_API_KEY`, so the current Universal Solver start gate matches the published SWE package contract.
- SWE Go runtime also contains `codex`/`OPENAI_API_KEY` support and `opencode.json` contains explicit OpenAI-compatible provider support. Universal Solver already maps Gonka into `OPENAI_API_KEY` + `OPENAI_BASE_URL`. Therefore a small provider-contract/bootstrap adaptation is plausible, but exact Gonka compatibility for the SWE harness is not runtime-proven yet.
- Workforce health is insufficient evidence of SWE readiness: its healthcheck validates only `/afhome/us-e2e-provenance.txt`, not required SWE node liveness/readiness.
- Current AgentField registry still exposes `build`, `plan`, `implement_issue`, `resolve` and internal schemas for `swe-planner`, but these are not a live callable capability while the node is offline.
- `swe-af:main` and `swe-af:dev` were observed diverged before this SoT commit: dev was 5 commits ahead and 1 commit behind main. The main-only downstream commit is `docs: add canonical error ledger`, conflicting with the clean-mirror invariant.

## Current Stage

BROWNFIELD RECOVERY / PROVIDER ENABLEMENT / SEMANTIC ACCEPTANCE.

Infrastructure bootstrap, exact-SHA reconciliation, discovery, and one real SWE SourceLoop capture are already proven. The first failed prerequisite in CURRENT permanent DEV is provider/bootstrap enablement: SWE is intentionally not started because its declared Anthropic/OpenRouter provider gate is unsatisfied. The evidence ladder is therefore provider contract -> SWE active/ready -> non-mutating smoke -> bounded `implement_issue`.

## Bounded Development Batches

Default batch size: about 30 minutes. Each batch closes a coherent DoD gate and writes back this file. Prefer the smallest 20% of work that removes the next 80% blocker.

### Batch 1 — Provider/bootstrap gate and SWE liveness

Status: IN PROGRESS

DoD:
1. Reconcile CURRENT topology/provider contract. **PASS**: legacy standalone SWE is removed; permanent DEV `edsh...` is current; isolated B is comparison-only.
2. Prove why CURRENT permanent DEV has no live `swe-planner`. **PASS**: canonical bootstrap evidence shows Anthropic/OpenRouter resolved empty and SWE was intentionally skipped.
3. Determine whether already-configured Gonka/OpenAI-compatible credentials can satisfy SWE without adding a new external secret. **PARTIAL / STRONGER EVIDENCE**: Universal Solver now injects Gonka as `OPENAI_API_KEY` + `OPENAI_BASE_URL`, and SWE runtime is explicitly pinned to `SWE_DEFAULT_RUNTIME=open_code` with `SWE_DEFAULT_MODEL=openai/<Gonka model>`. AgentField OpenCode harness passes caller env through to the subprocess, so the configured OpenAI-compatible base can survive to the coding tool. A live semantic call is still required to prove no fallback.
4. Choose the smallest safe enablement path. **DECIDED / COMPATIBILITY SHIM PENDING PROOF**: CURRENT `af run swe-planner` still fails before process launch because the installed package preflight accepts only `ANTHROPIC_API_KEY` or `OPENROUTER_API_KEY`, even with Gonka present. Direct live-file mutation of `/src/swe-af` is blocked by the current DEV mediator. Universal Solver therefore uses a bounded dummy `OPENROUTER_API_KEY` only to satisfy the legacy preflight while keeping actual SWE routing pinned to `open_code + openai/<Gonka>`. This shim is accepted only if execution evidence proves no OpenRouter fallback; otherwise it must be reverted and the SWE manifest fixed when live write capability is available.
5. Apply SWE product code/config changes only in CURRENT `/src/swe-af`; run targeted tests and `go/ make check` before restart. Do not program through GitHub/redeploy.
6. Start/reload only SWE or the smallest owning target; verify `active/ready` independently of workforce provenance health.
7. Run one non-mutating reasoner/schema smoke with execution evidence and confirm live `implement_issue` callability.
8. Write accepted evidence here; only after acceptance materialize the exact live delta durably to `fork/dev`/SourceLoop.

Stop conditions: no real coding canary while SWE is offline; no provider-secret mutation without explicit authorization; no whole-fleet or operator redeploy merely to work around the inner-loop tooling gap.

### Batch 2 — First bounded `implement_issue` canary

Status: PENDING

Prerequisite: Batch 1 DoD PASS.

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
- exact `WORKING_DEV_SHA` is recorded;
- repo tests pass on that durable SHA.

### Batch 5 — Materialization and regression

Status: PENDING

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
- `MODEL_REASONING`
- `PROMPT_OR_PLANNING`
- `ACI_HARNESS`
- `REPO_ENVIRONMENT`
- `GIT_DELIVERY`
- `AGENTFIELD_RUNTIME`
- `PROVIDER`
- `SOURCELOOP_RECONCILIATION`
- `OBSERVABILITY_GAP`

## Anti-Drift Checklist per Batch

- Reread `PLAN.md` before mutation.
- Observe CURRENT runtime before mutation.
- Verify exact target/SHA/process generation.
- Use the smallest discriminating test first.
- Reload only the same target if needed.
- Do not advance on health alone; require functional/semantic evidence.
- After acceptance, materialize exact source to Git and record SHA.
- When runtime and Git disagree, label `DESIGN_RUNTIME_DRIFT` or `PROVENANCE_GAP` as appropriate; do not silently reconcile.

## Known Drift / Debt

- [OPEN] `swe-af:main` contains a downstream `ERRORS.md` commit and therefore is not a clean upstream mirror. Resolve before the next upstream rebase/reconciliation cycle.
- [OPEN] DEV VPS Terminal source/deployed target registry contains `agentfield-dev-workforce`, but the CURRENT callable ChatGPT DEV surface lacks typed File ACI/process-start operations and generic `execContainer`/`startSession` fail closed as `REVIEW_REQUIRED: opaque_or_unknown_mutation`. `getOperatorGuidance` explicitly returned `ALLOW` for the scoped reversible SWE start, but execution remained blocked by mediation. Generic approval is unavailable (`approval_capability_gap`). This is `OPERATOR_PLANE_CAPABILITY_DRIFT`, not a SWE defect; do not redeploy the operator merely to bypass it.
- [OPEN] Current workforce healthcheck proves provenance HTTP only and can be green while `swe-planner`/`swe-pro` are offline.
- [OPEN] Current SWE absence is explained by the provider/bootstrap gate, not by a proven current-generation process crash: Anthropic/OpenRouter resolved empty and bootstrap intentionally skipped SWE. Historical process-exit logs remain forensic only.
- [OPEN] Provider-contract gap: SWE code contains `codex`/`OPENAI_API_KEY` and OpenAI-compatible OpenCode support, while `agentfield-package.yaml` and Universal Solver bootstrap admit only Anthropic/OpenRouter. Exact Gonka execution compatibility must be proven before adapting the contract.

## BMAD Workflow Used for Current Batch

- Entry: `bmad-help` from `BMAD-MNNZ`.
- Initial classification: brownfield recovery / quick implementation; `bmad-correct-course` rejected because its PRD + Epics prerequisites do not exist here.
- After CURRENT reconciliation exposed an implementation/debugging evidence gate, the canonical Universal Solver handoff explicitly routes this workstream to `bmad-testarch-test-design`.
- Active BMAD mode: risk-based system-level test/debug architecture embedded in this existing SoT (no duplicate BMAD artifact). Gate order: A0 CURRENT state -> A1 targeted regression -> A2 reasoner/pipeline -> A3 same-target reload if needed -> A4 functional canary -> A5 semantic E2E -> A6 durable accepted SHA.
- `bmad-quick-dev` remains the implementation method once a source/config defect is proven and live-edit capability is available.
- Write-back target: this `PLAN.md`.
