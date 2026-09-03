# SWE-AF Project Plan (Source of Truth)

Status: active
Last reconciled: 2026-09-02
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

LIVE-TASK DEBUGGING / SEMANTIC ACCEPTANCE / STRUCTURED-OUTPUT RELIABILITY.

Fresh CURRENT readback on 2026-09-02 supersedes the older PID/provider-gate wording below for the active debug batch:

- Permanent DEV workforce is container `1437789b5c4debc061992bec718132dcbccc66ac67e0b9e45ce1eea5b651c9e7` (`workforce-edshqtkwskg3lrczekhcmd71-105906341777`), image `edshqtkwskg3lrczekhcmd71_workforce:ffa6c6a56814b59da19903bd56c04f8cafdb44ae`, healthy, restart count `0`, OOM=`false`.
- `/afhome/installed.yaml` is the loaded-process authority for this generation: `swe-planner` was installed from local `/src/swe-af/go`, PID `8059`, started `2026-09-02T11:09:42Z`, port `8800`. `swe-planner.log` proves both `swe-planner` and `swe-pro` registered and initialized in this generation. Provider/bootstrap/liveness is therefore no longer the blocking gate.
- Direct `run_coder` behavior is not uniformly broken. After three early structured-output failures (`exec_20260902_112138_30mrlcvk`, `exec_20260902_113131_ju4aqfom`, `exec_20260902_113412_3jvaz47h`), four later direct coder executions completed (`exec_20260902_113619_5298bek0`, `exec_20260902_115708_d5rur0zd`, `exec_20260902_120030_8m4rpg5e`, `exec_20260902_120338_5ac3xq2q`). The last one is also independently recorded in `ERRORS.md` with 13/13 external-cwd tests PASS after bounded recovery.
- Two full `implement_issue` L2 attempts then failed at the child-coder boundary: parent `exec_20260902_120757_zb70eml7` -> child `exec_20260902_120757_yk42tvru`, and parent `exec_20260902_121608_24prhnbd` -> child `exec_20260902_121608_ezcvem44`. Both parents emitted `call.outbound.failed`; no false PASS was observed. Both child failures were surfaced as `Schema validation failed after 2 retry attempt(s). Last error: The output file was NOT created.`
- Container-first wrapper evidence localizes the second L2 failure more precisely. The three OpenCode attempts for `/tmp/fcm-calculator-live/.worktrees/40e731f4-calculator-l2-retry-power-modulo-json-cl` are captured as `/tmp/opencode-run-20451-1788351368.*`, `opencode-run-20974-1788351484.*`, and `opencode-run-21477-1788351564.*`. The agent successfully inspected/edited calculator files and repeatedly ran `python3 test_calculator.py`; the final retry reported 13 tests `OK`. Each attempt then terminated with OpenCode `type:error`, `UnknownError`, message `"Stream error occurred"` before the AgentField output file was written. Therefore `output file was NOT created` is a downstream schema symptom for this L2 failure, not yet the first causal error.
- The first causal boundary currently evidenced for L2 is `OpenCode -> fcm stream -> AgentField structured-output file protocol`. Failure class is `PROVIDER / ACI_HARNESS boundary` until the upstream stream error is made more specific. Do not patch coder business logic merely to silence the schema symptom.
- A separate planning defect remains isolated: `plan` parent `exec_20260902_113433_etlccqv6` / `run_product_manager` child `exec_20260902_113433_82kd1vcp` failed in ~16 ms with `API error: Invalid or missing broker client token`. This is an authentication/configuration gap in the planning path and is not the current L2 coding blocker.
- PR-AF is separately active in the same workforce (PID `31633`, started `2026-09-02T13:18:24Z`) and continues high-latency parallel OpenCode/LLM activity. Its evidence may help classify shared provider/schema behavior, but correlation alone must not be promoted to SWE causation.
- The connected AgentField Action currently reports healthy gateway status while execution-list/node-list reads return `Bad Gateway`; use direct permanent-DEV component logs and wrapper captures as the CURRENT execution evidence until that control-plane read path is restored.

BMAD `bmad-testarch-test-design` revalidation verdict for this stage:

1. **P0 semantic gate:** a bounded full `implement_issue` must terminate with a valid structured child result and no `call.outbound.failed`; transport status or successful task-file edits alone are insufficient.
2. **P0 failure-classification gate:** preserve the underlying OpenCode/provider error when schema output is absent; a generic schema error must not mask the first causal provider/transport event.
3. **P1 concurrency/provider gate:** before attributing failure to PR-AF or provider contention, compare one quiet/bounded SWE retry with runtime/provider evidence; do not stop or mutate unrelated agents merely to simplify the test.
4. **P1 planning-auth debt:** track `broker client token` independently and repair only after the L2 coding gate or if it becomes a hard prerequisite.
5. **Nearest compulsory move:** use the existing failed-task artifacts to expose/classify the upstream FCM stream failure, then rerun the SAME bounded L2 `implement_issue` once. Patch only the proven owning layer; after any source patch run its canonical regression, reload only the owning same-runtime target if required, then repeat the identical canary.

The recovery ladder remains L2-first until this gate passes. Do not advance to L3/L4/L5 and do not redeploy the whole workforce as a debugging primitive.

### BMAD test-design gate update — 2026-09-03

`bmad-help` was reread from canonical `n0namer/BMAD-MNNZ`, and the existing brownfield test/debug stage was revalidated with `bmad-testarch-test-design` in Edit mode. No duplicate test-plan artifact was created.

CURRENT evidence for the compulsory L2 failure-classification gate:

- AgentField gateway health is reachable, while capability discovery and recent-execution reads still return `Bad Gateway`; execution absence cannot be inferred from that control-plane API.
- Permanent workforce container `1437789b5c4debc061992bec718132dcbccc66ac67e0b9e45ce1eea5b651c9e7` remains the current healthy SWE canary runtime from the 2026-09-02 batch.
- DEV FCM topology is live without restart/OOM evidence at the container layer. `broker-dev-wgifzaww64jjnhazzed2nrrz` is healthy with restart count `0`, and `fcm-private-dev-wgifzaww64jjnhazzed2nrrz` is running behind nginx.
- Canonical FCM service SoT (`n0namer/server-ops/services/fcm-llm-gateway.md`) identifies DEV `keyless-dev` as primary `kilo/kilo-auto/free`, secondary `llm7/minimax-m2.7`, and explicitly leaves real request-time 429/5xx/timeout failover acceptance open. A health-aware preselection proof is not equivalent to runtime failover proof.
- Canonical FCM source is `n0namer/free-coding-models:dev`; the current source tree contains dedicated router/failover/telemetry components (`src/core/router-daemon.js`, `provider-cooldown.js`, `runtime-telemetry.js`). This establishes the likely owning layer but does not prove the failed request's candidate-attempt sequence.
- Broker logs available through the current bounded container-log surface are dominated by successful `/healthz` reads and do not expose the required per-request candidate chain. Read-only `/work` inspection through generic container exec is currently fail-closed by DEV operator mediation with `OBSERVE_REQUIRED: scope_unknown`, even after container identity/mount observation. This is an `OBSERVABILITY_GAP`, not evidence that failover did or did not occur.

Risk/test-design verdict remains **P0 BLOCKED on evidence, not implementation**: do not patch SWE coder/schema logic, do not stop PR-AF, do not change provider/model credentials, and do not redeploy. The next mandatory test action is one request-correlated `keyless-dev` trace proving `candidate #1 -> exact upstream result -> failover decision -> candidate #2 attempt/result or explicit absence -> returned FCM error`. Only after that evidence may ownership move to an FCM implementation/config fix; otherwise retain the same SWE L2 canary and investigate stream/concurrency lifecycle.

L2 semantic DoD remains unchanged and is not PASS: the same bounded `implement_issue` must produce a valid structured child result, no `call.outbound.failed`, expected bounded diff, independent tests PASS, and no unresolved provider/stream failure.

### BMAD source-level discriminator — 2026-09-03

Following the same `bmad-testarch-test-design` evidence gate, canonical FCM source was inspected at exact `n0namer/free-coding-models:dev` SHA `574f300458249252a5616406184c7ef3395e4f3f` before any implementation mutation.

Source evidence in `src/core/router-daemon.js` narrows the candidate defect:

- The router loops through candidates, records each `candidate.key`, and invokes `proxyStreamingRequest` for streaming requests. When a result returns `failoverToNext`, the next eligible candidate is selected and an `app_router_failover` event is emitted.
- Retryable upstream HTTP responses, auth failures, provider URL failures, empty streams, first-chunk failures, and transport errors before any client data is sent all return `failoverToNext: true`.
- After partial streaming data has already been sent, only `timeout` / `stream_stall_timeout` returns `failoverToNext: true`.
- A non-stall mid-stream error after partial output explicitly logs `Streaming failure after partial response`, closes the response, and returns `{ done: true }`; the outer router therefore does **not** attempt the secondary candidate on that branch.
- Final `All routed models failed for set: <set>` is produced only after the candidate loop exhausts its allowed attempts without a `done` result. The response already carries `models_tried` and quota metadata.
- The router has an existing unauthenticated `/stats` payload containing `requestLog`, `activeRequests`, routing order, and circuit-breaker state. This is the preferred bounded correlation surface; no new observability service or file is justified.

This is a **source-proven behavioral gap / runtime-correlation pending** result. It is consistent with the observed OpenCode `Stream error occurred`, but the failed SWE request has not yet been correlated to the post-partial non-stall branch by request id. Therefore `bmad-quick-dev` implementation is not yet authorized by evidence: changing streaming failover now would promote an inference into a product fix.

Nearest compulsory move is now sharper: obtain the router `/stats` / request-log entry for the same failed SWE/OpenCode request (or reproduce the same bounded L2 once while capturing its FCM `x-request-id`). If the entry shows a post-partial non-stall stream error with no second candidate, ownership moves to FCM and `bmad-quick-dev` should implement the smallest failover/stream-lifecycle fix with regression coverage. If it instead shows both `keyless-dev` candidates attempted, classify their exact upstream errors and patch that owning provider/config path instead.

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
- Coder-only diagnostics then isolated the next layer. `exec_20260901_130824_ktc5o8be` reproduced the same schema failure and raw OpenCode returned `Unexpected server error` (`err_67786025`) with correct CLI argv/repo/model. Renaming the custom provider from built-in `openai` to explicit `gonka` did not remove the error (`err_bf30da7d`).
- Adding an explicit `gonka.models` catalog for `deepseek-ai/DeepSeek-V4-Flash-0731` changed the failure signature materially: `exec_20260901_131316_4naadw0k` no longer failed immediately and instead reached real inference until the harness 300s no-progress timeout. This is accepted evidence that explicit model registration is required for the custom OpenCode provider and that the remaining Gonka blocker is inference availability/latency, not local provider resolution.
- Free OpenRouter canary `exec_20260901_131920_y14gwukz` reached the OpenRouter endpoint but returned HTTP 401 `Missing Authentication header`; the CURRENT workforce does not have a usable OpenRouter runtime credential. Do not inject or copy a secret merely to continue this diagnostic ladder without explicit authorization.
- Historical Universal Solver evidence identified the previously accepted SWE OpenCode runtime contract: OpenCode `1.17.15`, `SWE_OPENCODE_BIN=/afhome/.opencode/bin/opencode`, `SWE_DEFAULT_RUNTIME=open_code`, model `gonka/deepseek-ai/DeepSeek-V4-Flash-0731`, `AI_BASE_URL`/`OPENAI_BASE_URL` pointed at Gonka, explicit `gonka` provider/model catalog, and `ProjectDir=Cwd`. Current `af-deep-research` does not use OpenCode; it uses direct LiteLLM/AgentField AI.
- The live runtime was aligned to that historical OpenCode contract without rebuilding Go: wrapper pinned OpenCode `1.17.15`, `opencode.json` returned Gonka `baseURL` to `AI_BASE_URL`, and the wrapper preserved the loaded PID's `openai/...` -> `gonka/...` model translation. Calculator rerun `exec_20260901_133823_bvptiqkq` / coder `exec_20260901_133825_rm0tpmjp` still reached exactly the harness `300s` no-progress timeout with zero files changed. Therefore OpenCode version/provider/base/model wiring is no longer the leading hypothesis; CURRENT Gonka inference availability/latency remains the first failed gate.
- Independent filesystem oracle confirmed no `canary/calculator` directory was created. The disposable canary copy inherited the pre-existing live SWE source delta from `/src`, but the calculator task itself produced no task files.

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
