# SWE-AF Project Plan (Source of Truth)

Status: active
Last reconciled: 2026-09-01
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

Reconciled from live readback on 2026-09-01.

- Permanent AgentField DEV Coolify application: `edshqtkwskg3lrczekhcmd71`; repo `n0namer/universal-solver`; CURRENT loaded orchestrator generation remains `dr-debug-freeze-20260901` at SHA `b7e2f00116358d78d01a73b77aa31d1c2bdfb9d5`.
- CURRENT workforce container is `c4283a751fe4cfa7fed5a13f7b26e431354f1c46adf1e1774f475f0355ab6e2c` (`workforce-edshqtkwskg3lrczekhcmd71-084728280230`), running and healthy. Current sibling generation includes control plane `02ecf09a404da56fb7961c0a05ab9444bf82ce8efc33ede5a747e7d7fa4cacb6` and runtime-capture `214a333d1c478d7228628b239616a9f85f61e5f2909586c45711cbaec6fa6def`. Its startup provenance baseline for SWE remains `58c4e0d19081bc52363c120b7963a34cebb1e894` plus the preserved live `/src/swe-af` working-tree delta.
- `PLAN.md` write-backs on `dev` can advance the branch without changing product code. Therefore branch HEAD is **not** the live product identity during this debug batch; use `58c4e0... + live working-tree delta + loaded process generation` until accepted code is canonicalized. Do not redeploy merely to synchronize PLAN-only commits.
- The exact `b7e2f001...` compose uses `preserve_or_reconcile`: when `/src/swe-af` is at baseline `58c4e0...` and dirty, restart preserves the live working-tree delta instead of resetting it. Deep Research existing `/app` is likewise preserved.
- `/src` remains the writable persistent runtime-source lane for SWE-AF. Product debugging for this phase must edit `/src/swe-af` directly; GitHub/redeploy is not the inner coding loop.
- SourceLoop SWE canary has already passed at least once: runtime → capture → Git; accepted SWE dev SHA is recorded in Universal Solver fleet lock.
- `runtime-capture` is currently running and has recent successful capture activity for Deep Research.
- `meta_deep_research` is currently `active/ready`.
- Historical standalone SWE service `universal-solver-swe-af` (`wetscrp2tj90tklmlvkcadfw`) existed during Wave 0, but canonical server-ops records it as removed during cleanup on 2026-08-21. It is superseded topology; CURRENT target is permanent DEV `edshqtkwskg3lrczekhcmd71`.
- Historical isolated B `2zciq6hujpev6dbudcdlijqq` remains healthy but is forensic comparison only, not the current mutation target.
- `swe-planner` and `swe-pro` are live in CURRENT permanent DEV. Direct `/afhome/logs/swe-planner.log` readback from workforce `c4283a...` shows `swe-planner` listening on `:8800`, `node.register.complete`, and `agent.initialize.complete` at `2026-09-01T08:56:24Z`; the embedded pro engine shows the same registration/initialization sequence for `swe-pro` on `:8801`.
- CURRENT `b7e2f001...` generation preserves the live `/src/swe-af` delta. The live package manifest admits `OPENAI_API_KEY` and defaults SWE to `open_code` with `openai/deepseek-ai/DeepSeek-V4-Flash-0731`; `node.go` preserves `OPENAI_BASE_URL`, and `opencode.json` contains the generic OpenAI-compatible provider.
- The historical package/bootstrap admission defect is no longer the active gate. Provider/bootstrap liveness is PASS and real SWE executions have been issued successfully through a native authenticated runtime path.
- CURRENT execution-control path is now proven: production VPS Terminal target `agentfield-dev-deep-research` can call the CURRENT control-plane with its existing runtime credential without exposing the secret. Discovery, async execution, execution polling, and execution-events readback all work through this path. External AgentField Actions may still be unavailable with `no available server`, but that no longer blocks SWE acceptance work in this chat.
- Workforce HTTP health remains insufficient evidence of SWE readiness: use reasoner registration plus real execution evidence instead.
- CURRENT live source is ahead of the loaded `swe-planner` binary. `/afhome/installed.yaml` shows `swe-planner` PID `7686`, started `2026-09-01T08:56:24Z`, installed from local source path `/src/swe-af/go`; the process has not been rebuilt/reloaded after the latest live-source edits. This is `LOADED_SOURCE_DRIFT`, not a reason to redeploy the container.
- CURRENT live `/src/swe-af` delta now includes: OpenAI-compatible/Gonka admission; OpenCode runtime/config support; planning Product Manager direct AgentField-AI path with OpenCode fallback; `n.App` wired as planning AI; Deep-Research-style custom base semantics; `openai/` provider-prefix stripping before direct Go-AI model calls; and a 30-minute Go AI transport timeout matching the proven Deep Research semantic timeout profile. `git diff --check` is PASS.
- Real executions exposed two distinct diagnostics: (a) OpenCode initially failed before LLM because the permanent workforce lacked the OpenCode runtime/PATH; this was repaired live and OpenCode now launches; (b) current Gonka inference is degraded independently of SWE — `/models` is healthy while valid `/chat/completions` stalls and a minimal request eventually returned HTTP 502 `context deadline exceeded`. A minimal Deep Research reasoner on the same source SHA/base/model also stalled, confirming this is not only an SWE/OpenCode defect.
- The current acceptance objective has changed from infrastructure proving to live-task debugging: use bounded real tasks as diagnostic experiments, classify the first failed gate from runtime evidence, patch only that layer in the live container, repeat the same task, and raise task complexity only after independent PASS.
- `swe-af:main` and `swe-af:dev` were observed diverged before this SoT commit: dev was 5 commits ahead and 1 commit behind main. The main-only downstream commit is `docs: add canonical error ledger`, conflicting with the clean-mirror invariant.

## Current Stage

LIVE-TASK DEBUGGING / SEMANTIC ACCEPTANCE / RECOVERY LADDER.

Infrastructure bootstrap, node registration, authenticated execution control, and real non-mutating SWE executions are already proven. The current goal is no longer to prepare infrastructure in isolation; it is to use progressively harder live engineering tasks as controlled diagnostic experiments. The first failed gate observed in each task must be localized from runtime evidence, repaired minimally in live `/src`, and the same task repeated before complexity increases.

The CURRENT source is not yet the CURRENT loaded binary: `swe-planner` PID `7686` predates the latest Deep-Research-style direct-AI/timeout source delta. Until that exact source is loaded into the same process lane, tests against PID `7686` validate the older binary only. Do not redeploy the container merely to load source. The next execution ladder is: load exact live source into the existing SWE process lane -> L0 minimal execution -> L1 calculator coding task -> L2 seeded regression repair -> L3 multi-file feature -> L4 forced recovery -> L5 bounded real-repository issue.

## Bounded Development Batches

Default batch size: about 30 minutes. Each batch closes a coherent DoD gate and writes back this file. Prefer the smallest 20% of work that removes the next 80% blocker.

### Batch 1 — Provider/bootstrap gate and SWE liveness

Status: IN PROGRESS

DoD:
1. Reconcile CURRENT topology/provider contract. **PASS**: legacy standalone SWE is removed; permanent DEV `edsh...` is current; isolated B is comparison-only.
2. Prove why CURRENT permanent DEV has no live `swe-planner`. **PASS**: canonical bootstrap evidence shows Anthropic/OpenRouter resolved empty and SWE was intentionally skipped.
3. Determine whether already-configured Gonka/OpenAI-compatible credentials can satisfy SWE without adding a new external secret. **PARTIAL / SOURCE SUPPORT PROVEN, RUNTIME PROOF MISSING**: CURRENT `b78866...` already maps Gonka as `OPENAI_API_KEY` + `OPENAI_BASE_URL`; SWE source contains OpenAI-capable runtime/provider paths. The remaining proof is an actual `swe-planner` execution routed to Gonka with no OpenRouter fallback.
4. Choose the smallest safe enablement path. **DECIDED / CURRENT 30-MINUTE BATCH**: first patch SWE's live `/src/swe-af/go/agentfield-package.yaml` so `llm_provider.require_one_of` also admits `OPENAI_API_KEY`; run the smallest manifest/config regression plus `cd go && make check` when the live runner is available; then start/reload only `swe-planner` with explicit `open_code + openai/<Gonka>` runtime/model and preserve the configured OpenAI-compatible base. Do not modify provider secrets and do not use GitHub-edit -> redeploy for implementation. If the live typed file/process lane is unavailable, classify `OPERATOR_PLANE_CAPABILITY_DRIFT` and stop mutation rather than bypassing mediation.
5. Apply SWE product code/config changes only in CURRENT `/src/swe-af`; run targeted tests and `cd go && make check` before restart. **LIVE DELTA APPLIED / SOURCE-INTEGRITY PASS / FULL GO VALIDATION BLOCKED**: shared `/src` was edited container-first. Current live delta: (a) `go/agentfield-package.yaml` admits `OPENAI_API_KEY` and now defaults installed SWE to `SWE_DEFAULT_RUNTIME=open_code` plus `SWE_DEFAULT_MODEL=openai/deepseek-ai/DeepSeek-V4-Flash-0731`; (b) `go/internal/node/node.go` preserves `OPENAI_BASE_URL` into AgentField Go `ai.Config.BaseURL` when `AI_BASE_URL` is unset; (c) `go/internal/node/aiconfig_test.go` adds a regression for that custom base; (d) `opencode.json` adds a generic `openai` provider via `@ai-sdk/openai-compatible` using `OPENAI_API_KEY` + `OPENAI_BASE_URL`. CURRENT workforce readback sees the manifest defaults and `git diff --check` PASS. Runtime-capture automatically materialized the container delta as provenance proposal PR #7; latest capture commit observed `0dc63fcdaab907c6c131d194001546d44be4cd38` (do not merge/use as debug source). Exact-source Coding Station validation was attempted but its runner lacks `go` (`FileNotFoundError`), so this is `VALIDATION_BLOCKER`, not a test failure. No GitHub product-code edit or redeploy was used.
6. Start/reload only SWE or the smallest owning target; verify liveness independently of workforce provenance health. **PASS**: CURRENT workforce `c4283a...` consumed the preserved live source. `/afhome/logs/swe-planner.log` proves `swe-planner` listening on `:8800` with `node.register.complete` + `agent.initialize.complete`, and `swe-pro` listening on `:8801` with the same registration/initialization sequence at `2026-09-01T08:56:24Z`.
7. Run one non-mutating reasoner/schema smoke with execution evidence and confirm live reasoner callability. **PASS FOR EXECUTION CONTROL / PARTIAL FOR SEMANTIC ROUTE**: authenticated discovery, async reasoner execution, polling and event readback now work through the existing `agentfield-dev-deep-research` VPS-Terminal target into CURRENT control-plane. Multiple `swe-planner.plan` and internal reasoner executions were observed. The remaining semantic issue is provider/runtime behavior, not execution-control availability.
8. Load the latest live `/src/swe-af` source into the existing `swe-planner` process lane without GitHub-first redeploy, then repeat the exact discriminating smoke. **PENDING / LOADED_SOURCE_DRIFT**: installed local package points to `/src/swe-af/go`, but PID `7686` predates the newest direct-AI/timeout source delta. Do not claim the source fix tested until loaded-process identity matches the edited source.

Stop conditions: no real coding canary while SWE is offline; no provider-secret mutation without explicit authorization; no whole-fleet or operator redeploy merely to work around the inner-loop tooling gap.

### Batch 2 — Live-task diagnostic ladder: calculator → seeded bug

Status: IN PROGRESS — L1 FIRST RUN EXPOSED FALSE SUCCESS

Prerequisite note: the first L1 run was intentionally executed against CURRENT loaded PID `7686` to expose actual runtime behavior before latest source reload. Results below therefore diagnose the loaded binary, not the newest `/src` delta.

DoD:
1. Use an ephemeral calculator workspace only; no GitHub delivery and no persistent infrastructure.
2. L1: ask SWE to implement a tiny CLI calculator with add/subtract/multiply/divide, controlled divide-by-zero behavior and agent-written unit tests.
3. Independently verify deterministic examples, simple algebraic properties and unintended files changed; agent-written tests are not the sole oracle.
4. Capture execution_id, reasoner/node, status, wall time, runtime-log causal timeline, before/after file set, actual diff, tests, retries and primary failure class.
5. On FAIL, stop at the first causal failed gate, patch only that layer live in the container and rerun the same task.
6. After L1 PASS, seed one known regression in the same disposable workspace and require SWE to diagnose and repair it minimally.
7. Require two independent PASS observations across the level boundary before advancing.
8. Write exact evidence and the next failed gate back here.

Observed L1 run evidence — 2026-09-01:
- Canary repo: ephemeral `/tmp/swe-af-calculator-canary`; product `/src/swe-af` was not the task target.
- Top execution: `exec_20260901_124949_s0zbmlbq`, duration ~60.1s, AgentField status `succeeded`.
- Semantic result contradicted transport status: `success=false`, `0 commit(s)`, `0 file(s) changed`, no verification. This is a `FALSE_SUCCESS` acceptance defect.
- Child coder `exec_20260901_124951_rrem07v9` returned AgentField status `succeeded` but payload `complete=false`, `files_changed=[]`, summary `Coder agent failed ... Schema validation failed ... output file was NOT created`.
- Child reviewer `exec_20260901_125021_77nrhrs2` also returned AgentField status `succeeded` while synthesizing `approved=true`, `blocking=false`, summary `reviewer agent failed ... not blocking`.
- Root cause: fail-open role/orchestration contracts convert coder/reviewer harness failure into ordinary result payloads and allow top-level completion with zero diff.
- Live `/src` repair applied container-first: coder schema/transport failure now returns error; reviewer failure/no structured output now returns error; default and flagged coding loops no longer convert reviewer error into approval. Existing role tests were changed from fail-open expectations to fail-closed regression expectations. `git diff --check` PASS.
- This repair is source-only until loaded-process identity is refreshed; do not claim runtime PASS on PID `7686` yet.

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
- [OPEN] DEV VPS Terminal can edit/read CURRENT live source, but it still lacks a typed Go build/process-only restart primitive for `swe-planner`; generic `go`, `kill`, and similar process commands fail closed under mediation. This is now specifically a `LOADED_SOURCE_DRIFT` / process-hot-load gap, not an execution-control gap. Do not redeploy the whole container merely to load source.
- [OPEN] Workforce healthcheck still proves provenance HTTP only; current SWE readiness is instead evidenced by registration plus real reasoner execution. Preserve this distinction in future canaries.
- [RESOLVED] Historical SWE absence was caused by the provider/bootstrap admission gate. CURRENT generation has consumed the live provider delta and both SWE nodes are registered/initialized, so this is no longer the active blocker.
- [OPEN] Provider/runtime semantic acceptance remains incomplete. Real SWE executions now reach coder/reviewer stages, but structured-output harness failures persist on the loaded binary. Independent Gonka probes show healthy `/models` with stalled valid `/chat/completions`, and the same stall reproduces through the previously working Deep Research path. The latest live `/src` delta aligns planning with Deep Research-style direct AI + custom base + 30-minute timeout, but it is not yet loaded into PID `7686`.

## BMAD Workflow Used for Current Batch

- Entry: `bmad-help` from `BMAD-MNNZ`.
- Initial classification: brownfield recovery / quick implementation; `bmad-correct-course` rejected because its PRD + Epics prerequisites do not exist here.
- After CURRENT reconciliation exposed an implementation/debugging evidence gate, the canonical Universal Solver handoff explicitly routes this workstream to `bmad-testarch-test-design`.
- Active BMAD mode: risk-based system-level test/debug architecture embedded in this existing SoT (no duplicate BMAD artifact). Gate order: A0 CURRENT state -> A1 targeted regression -> A2 reasoner/pipeline -> A3 same-target reload if needed -> A4 functional canary -> A5 semantic E2E -> A6 durable accepted SHA.
- `bmad-quick-dev` remains the implementation method once a source/config defect is proven and live-edit capability is available.
- Write-back target: this `PLAN.md`.
