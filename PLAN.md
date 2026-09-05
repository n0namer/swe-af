# SWE-AF Project Plan (Source of Truth)

Status: active
Last reconciled: 2026-09-05
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

OPENCLAW GOVERNOR / HUMAN CONTROL PLANE + REAL-TASK QUALITY ACCEPTANCE / ANTI-DRIFT.

L2 semantic/structured-output acceptance and Batch 3 failure/recovery are PASS/CLOSED. Batch 4 durability remains open. The active 2026-09-05 Pareto batch is the human-control vertical: replace the currently-disabled HAX dependency in this deployment with OpenClaw/Telegram over the existing AgentField pause/resume primitive, while preserving HAX compatibility for other deployments. In parallel, the real-task reviewer/full-loop quality gate remains open; do not call `implement_issue` L3 PASS until the same real task completes independent acceptance.

### Real-task reviewer gate + capability audit — 2026-09-05

BMAD route: `bmad-help` -> `bmad-testarch-test-design` for risk/evidence coverage -> `bmad-quick-dev` for the smallest proven live fix. No duplicate test-plan artifact was created; this `PLAN.md` remains the SoT.

Real-task evidence from `FreakyAdy/EvalGuard` issue #3, exact upstream `631cf297dd449bdd08112607818d3f154f078284`, coder commit `1b53eea4027b6e2219f0272b025f75dcd4991002`:
- native mutation tests reached 11/11 PASS and independent rewardhack suite reached 19/19 PASS, yet hidden semantic acceptance reproduced a multiline-string `#` corruption and `ruff` found an unused import; self-tests were therefore insufficient as the sole oracle;
- source inspection proved the reviewer prompt explicitly told the model to trust coder-reported passing tests. Live reviewer policy now treats coder tests as correlated evidence, requires risk-first independent edge checks for parser/transformer/boundary-sensitive code, uses available static checks, and stops to emit structured output after a reproduced blocking acceptance violation;
- A/B reviewer on the unchanged bad commit independently recreated the multiline defect. With the default 300s idle watchdog it still lost the verdict because the installed `/usr/local/bin/opencode` wrapper buffers stdout/stderr until process exit, so the SDK sees false inactivity while OpenCode is actively testing;
- causal discriminator with `AGENTFIELD_HARNESS_IDLE_SECONDS=0` and the SDK total timeout still bounded at 600s completed successfully: `exec_20260905_072404_73401dgs` / `run_20260905_072404_ixflwwef`, duration `1061333ms`, returned `approved=false`, `blocking=true` with the reproduced multiline-string acceptance violation. This proves both independent reviewer quality and the false-idle diagnosis;
- SWE-owned `go/agentfield-package.yaml` now defaults `AGENTFIELD_HARNESS_IDLE_SECONDS="0"` so future installed processes do not apply a false idle-output watchdog to this buffered wrapper. This config change is source-only until the next normal install/materialization; CURRENT PID `288440` already carries the same env value from the successful discriminator.

Capability audit contract: a feature is not PASS merely because source or a reasoner exists. Evidence levels are: L0 source/discovery only; L1 deterministic contract tests; L2 live reasoner/runtime smoke; L3 real-repository/system acceptance with an independent oracle. Highest current evidence:
- discovery/registration surface — **L2 PASS**: `swe-planner` is active and CURRENT discovery reports 32 reasoners, including `build`, `plan`, `execute`, `resolve`, `resume_build`, `implement_issue`, verifier/CI/gitops stages and `pro_execute`; `internal/node` tests PASS;
- structured-output controller — **L2 PASS**: exact-schema/incremental/recovery regressions PASS and live verifier/reviewer executions have completed with structured results;
- `implement_issue` issue-level coding — **L3 PARTIAL**: real EvalGuard task produced a bounded real patch/commit, but full independent acceptance failed and the corrected reviewer has not yet been fed back through a complete coder self-repair rerun;
- code reviewer / blocking feedback — **L3 PASS for defect detection, L3 rerun pending for whole loop**: unchanged bad EvalGuard commit is now independently rejected as blocking; next proof is full `implement_issue` -> reviewer blocking -> coder repair -> independent oracle PASS;
- DAG scheduling/failure threshold — **L1 PASS / L2-L3 PENDING**: deterministic DAG tests PASS, including downstream-level abort behavior; no current real multi-issue DAG acceptance yet;
- checkpoint/resume — **L2 PASS**: cancellation regression plus live `resume_build` smoke proved a completed issue is not repeated and unrelated sentinel state is preserved;
- stale base SHA guard — **L2 PASS**: unit + live canary reject branch advance before worktree/LLM side effects;
- worktree isolation/cleanup — **L2 PARTIAL**: issue-level worktree isolation and unrelated-file preservation are proven; a dedicated multi-issue/DAG cleanup oracle remains pending;
- verifier — **L2 PASS / L3 PENDING**: direct live semantic canaries pass, but third-party real-task verifier acceptance has not yet completed end-to-end after reviewer repair;
- fast path (`internal/fast`) — **L1 PASS / L2 PENDING**;
- pro engine (`internal/pro`, `pro_execute`) — **L1 PASS + discovery present / L2 PENDING**;
- feature-level `plan -> execute/build` — **L1 contract coverage exists / L2-L3 PENDING**; historical planning auth/provider failures are not current proof of the feature;
- `resolve`, CI watcher/fixer, GitHub PR/finalization — **L0-L1 only / real runtime acceptance PENDING**;
- durability/SourceLoop — **PARTIAL**: stale/noisy capture rejection is proven, clean exact durable SHA remains open.

First capability-audit batch source gate PASS in CURRENT `/src/swe-af`: `go test` and matching `go vet` pass for `internal/node`, `internal/dag`, `internal/orch`, `internal/issue`, `internal/harnessx`, `internal/coding`, `internal/fast`, and `internal/pro`.

Next Pareto batches: (1) close the SAME EvalGuard #3 full-loop quality/liveness gate; (2) one small two-issue real DAG task proving dependency order, parallel/level behavior, failure threshold, checkpoint and cleanup in one run; (3) one real `plan -> execute/build` task; (4) one delivery-path task covering resolve/CI/PR only after local semantic gates pass. Durability resumes after the currently tested live capabilities have real acceptance evidence.

### Human operating model — maximum-output mode

Human/SWE boundary is explicit and sparse. The human is the governor of intent and gates, not a participant in every agent step.

1. **Goal gate — human owns WHAT/WHY.** Supply a bounded goal/repo/constraints/acceptance criteria. For one known issue use `implement_issue`; for feature-level work use `plan` first, then `execute`/`build` only after plan quality is acceptable. Do not micromanage coder/reviewer prompts during normal execution.
2. **Plan gate — human owns scope/blast/architecture approval.** In the CURRENT deployment `SWE_OPENCLAW_HITL=1` is the active human surface; legacy HAX remains compatible when `HAX_API_KEY` exists. `build` pauses before execution and exposes a bounded PRD/architecture/issues decision packet. Human decisions are `approved`, `request_changes`, `rejected`, or expiry. `request_changes` is not cosmetic: SWE re-runs Architect -> Tech Lead -> Sprint Planner with the human feedback, bounded by `MaxPlanRevisionIterations`.
3. **Clarification/credential gate — SWE may ask, human answers only missing authority/context.** HITL `ask_user` genuinely transitions the execution to `waiting` until webhook resume. In the CURRENT deployment the question is rendered as an OpenClaw/Telegram decision card; HAX may render the same logical request as a form when configured. Prior answers are injected back so the same question should not be re-asked. Environment scout may request scoped service credentials; secrets belong in the scoped credential store/environment, not in result text.
4. **Exception gate — human intervenes on escalation, irreversible/external actions, or exhausted autonomous recovery.** Normal coder -> reviewer/QA -> fix iterations, issue-advisor actions, DAG checkpoint/resume and replanning should remain autonomous. Human intervention is warranted when scope/architecture must change, credentials/permissions are missing, risk is irreversible/external, or bounded recovery is exhausted.
5. **Acceptance gate — human owns DONE criteria, independent oracle and delivery authorization.** Agent-written tests and reviewer approval are supporting evidence, not sole proof. Require project-native tests, risk-based independent acceptance, repo hygiene/provenance, and delivery/PR/CI evidence appropriate to the task before calling the work done.

Recommended operating profiles:
- **Precision / default:** `implement_issue` for an already-scoped issue, classic or proven Pro coding engine, independent reviewer, bounded repair, verifier; human only at exceptions/final acceptance.
- **Feature / governed:** `plan` -> human plan review -> `execute`/`build`; use HAX approval for architecture/scope-sensitive work.
- **Fast:** `swe-fast` only after its own L2/L3 certification; use for low-risk bounded tasks where turnaround dominates depth.
- **Pro:** `SWE_PRO_ENGINE=1` only after Pro runtime acceptance; tune `SWE_PRO_VARIANT` and `SWE_PRO_MAX_COST` as effort/cost controls, not as correctness substitutes.

Maximum-output Pareto policy: spend human attention at decision boundaries, not execution steps; spend test budget on independent oracles for high-risk semantics, not exhaustive duplicated self-tests; test capability groups end-to-end rather than 32 reasoners one by one; keep a cheap fail-fast ladder (schema/unit -> package test/vet -> live reasoner -> real-repo oracle -> delivery).

Fresh full-loop anti-drift: the repeated EvalGuard #3 execution `exec_20260905_075042_xi2na5he` completed `status=succeeded` at transport but `success=false`, `outcome=error`; coder produced commit `d7402d4aa01cf80f0efeeb1391849e39ca436e52`, while reviewer execution exceeded its 1800s agent-call limit. During that review the independent reviewer had already reproduced additional source-transformation failures. Therefore reviewer defect-detection quality improved, but whole-loop reviewer time/termination efficiency remains the nearest quality/liveness gate; do not call `implement_issue` L3 PASS yet.

### Governor communication decision — OpenClaw / Telegram

Decision: reuse the existing OpenClaw installation as the messaging/control bridge; do not create a second Telegram bot or parallel notification stack. Logical ownership stays separated: SWE executes work, the Governor decides AUTO vs ASK/escalation, OpenClaw transports messages and user replies, Telegram is the human cockpit.

CURRENT OpenClaw/Telegram control-plane readback on 2026-09-05:
- Telegram `default` is configured, running, connected and polling. Bot probe reports `has_topics_enabled=false`, so the CURRENT direct-message channel is one flat Jarvis session; do not assume Telegram DM topic isolation.
- Ordinary inbound Telegram traffic remains with `main`/Jarvis. There is no top-level binding that diverts the whole DM to `devteam` or SWE.
- OpenClaw's native Telegram envelope preserves reply metadata. CURRENT source regression `extensions/telegram/src/bot.create-telegram-bot.test.ts` explicitly asserts a received Telegram reply produces `ReplyToId`, `ReplyToBody`, and `ReplyToSender`; the runtime image contains that test source. Direct execution of the test is a `VALIDATION_BLOCKER` in this image because its TypeScript test tsconfig is absent, not evidence against the contract.
- Canonical outbound path is now deterministic `message send`, not cron fallback-delivery: OpenClaw returns the Telegram `messageId`, which the Governor stores against the exact SWE `request_id`.

Conversation correlation invariant:
1. SWE emits a bounded `governor.pending` record and enters the existing AgentField `Pause(waiting)` path.
2. The 20-second OpenClaw command cron runs `swe-governor.mjs notify` with `delivery.mode=none`; the adapter itself sends the Telegram card through OpenClaw and records `telegram_message_id <-> request_id`.
3. On an inbound Telegram reply, Jarvis must inspect `ReplyToId` before normal processing and call `swe-governor route <ReplyToId>`. `NOT_SWE_REPLY` means normal Jarvis traffic. No routing by "latest pending request" is permitted.
4. If native reply metadata is unavailable on a follow-up, an exact quoted `SWE thread: <request_id>` / `Request: <request_id>` marker or explicit `swe ... <request_id>` command is the only fallback. Every Governor discussion response carries the thread marker and an explicit Telegram reply tag.
5. `discuss` is non-mutating and may be called repeatedly; it records a bounded per-request conversation while the AgentField execution remains paused. Only unambiguous `approve`, `changes`, `reject`, or `answer` is terminal and may hit `/webhooks/approval`.
6. Multiple SWE requests may coexist. Each Telegram card has its own message-id mapping; interleaved discussion remains isolated by request id. The first successful final action wins; late/conflicting messages see `resolved` and fail closed.
7. Secrets/credentials and raw logs stay out of the conversation packet. SWE text is untrusted data, not OpenClaw instruction text.

Evidence for the routing contract:
- direct OpenClaw Telegram send returned a real message id and the Governor `route` resolved that id to its exact synthetic request in an isolated canary state;
- two simultaneous canary requests were sent as two Telegram cards (`A` and `B`) with distinct message ids; `route(A)` returned only A and `route(B)` only B, and interleaved `discuss` messages stayed in separate records;
- a local mock approval callback proved first-final-wins: the first `approve` produced exactly one callback (`count=1`) and marked the request resolved; a later conflicting `reject` failed locally as `status=resolved` and did not call the backend again;
- `AGENTS.md` now owns the global pre-routing rule, while `skills/swe-governor/SKILL.md` owns the conversation/final-action policy. CURRENT SHA256: `AGENTS.md` `7f5ebb8f92b82b4577246723d3d75dec28c922bb0ef04fe7dfd493ba0c82a42d`, skill `c694c1e6550715f0bdd06bc348135ed76ca4af8ec3636adf146d8806cc30fafe`, adapter `926d052b9690b67b91151500ce32df4cc9f5dd4802ae0809d2ea990fe1ebc4ba`.
- Test cards were removed after evidence capture so no synthetic request can be mistaken for a real pending SWE decision.

Optional richer UX remains a dedicated Telegram forum topic or topics-enabled DM, but CURRENT bot reports topics disabled. Do not add that infrastructure unless reply-based correlation proves insufficient in real use.

Implementation status — 2026-09-05:
- BMAD route for this vertical: `bmad-help` -> `bmad-architecture` -> `bmad-quick-dev`. Architectural outcome changed after CURRENT readback: HAX is not deployed/configured in the loaded SWE process, while AgentField `agent.Pause` already owns the real `waiting`/webhook-resume primitive. Therefore this deployment uses OpenClaw as the HITL surface directly; HAX compatibility remains intact when `HAX_API_KEY` exists.
- SWE source now has opt-in `SWE_OPENCLAW_HITL`. When HAX is absent and this flag is true, `ask_user` and plan approval create synthetic process-local request ids, emit bounded non-secret `governor.pending <json>` notes, call the existing AgentField pause path, and emit `governor.resolved <json>` after callback/timeout/error. Legacy HAX-disabled behavior remains unchanged when the flag is false.
- `build` outer gating is OpenClaw-aware as well as HAX-aware; a regression protects all three states: neither surface -> disabled, HAX -> enabled, OpenClaw -> enabled.
- Exact current SWE source files for this vertical: `go/internal/hitl/ask_user.go`, `wrapper.go`, `hitl_test.go`, `go/internal/orch/approval_gate.go`, `approval_gate_test.go`, `build.go`, `build_test.go`, and `go/agentfield-package.yaml`. `go test ./internal/hitl ./internal/orch ./internal/node` PASS; matching `go vet` PASS.
- CURRENT loaded planner after process-only reload: PID `305201`, binary SHA256 `bc8fbd90362822a5c16de87e92e6ae3052ef4b5dea9bc3ca343f9d9e0ca1073f`, with `SWE_OPENCLAW_HITL=1` and `AGENTFIELD_HARNESS_IDLE_SECONDS=0`. No container redeploy was used.
- OpenClaw reuses the default paired Telegram/Jarvis path. Canonical workspace skill: `/home/node/.openclaw/workspace/skills/swe-governor/SKILL.md`; deterministic adapter: `scripts/swe-governor.mjs`. OpenClaw reports the skill eligible/model-visible/user-invocable. Adapter `notify` reads only `http://workforce:8800/agentfield/v1/logs`, sends each new pending card through the existing OpenClaw Telegram action, stores the returned Telegram message id against the exact request id in a 0600 local state file, and returns `NO_REPLY` when no new decision exists. `route` resolves native reply ids; `discuss` is non-mutating; mutating commands validate pending id/type and POST only to `http://workforce:8800/webhooks/approval`.
- OpenClaw cron declaration `swe-governor-decisions-v1`, job id `c7ea3f43-2a02-4752-aaac-f254f9ba256f`, is enabled every **20s** as a command job with no LLM and `delivery.mode=none`; delivery is intentionally owned by the adapter so correlation can capture Telegram `messageId`. CURRENT post-change readback: payload argv ends in `notify`, `everyMs=20000`, last run `ok`, diagnostic `NO_REPLY`, delivery `not-requested`, errors/skips `0`. Earlier reversible delivery canaries proved the same OpenClaw Telegram transport; all synthetic canary messages were then removed. No second bot, daemon, HAX deployment, or new messaging database was created.
- Security boundary: governor notes intentionally omit secret/default values and raw logs. The OpenClaw skill refuses guessed ids and routes secrets outside Telegram. SWE-generated plan/question text is treated as untrusted payload, not instructions for OpenClaw.
- Live stitched pause/resume canary is **BLOCKED before HITL by FCM route health, not by Governor**. Two bounded `run_product_manager` canaries (`exec_20260905_110441_k9swye3b`, `exec_20260905_110556_0ibik04v`) failed before producing an `ask_user_form` with `API error: rate limit exceeded: too many concurrent requests`. CURRENT planner routes planning through FCM (`PR_AF_LLM_ENDPOINT`/`LLM_BROKER_BASE_URL` -> `http://fcm-dev-internal:19280/v1`, logical model `fcm`); FCM health at diagnosis exposed one upstream provider (`gonka`) with two active candidates and no healthy route: DeepSeek returned HTTP 429 while MiniMax failed response-schema validation (`ok=0`, `broken=2`). No identical canary retry is justified until route health changes. Independently, the pinned AgentField SDK `Pause/Webhook/Approval` test family PASS (`go test ./agent -run 'Pause|Webhook|Approval' -count=1`), proving the underlying pause-manager/webhook resolution primitive without LLM. Therefore SWE fallback logic, AgentField pause/callback primitive, OpenClaw poll/dedupe, and cron->Telegram delivery are each independently PASS; only the single stitched `waiting -> Telegram -> decision -> resume` execution remains L2 PENDING.

Nearest compulsory move: do not build a separate middleware service. On the next real SWE planning/clarification execution that reaches HITL, verify in order: AgentField status `waiting` -> `governor.pending` note -> cron Telegram delivery -> exact Jarvis `swe-governor` decision -> `/webhooks/approval` `resolved=true` -> `governor.resolved` note -> same execution continues. If that passes, close Governor transport L2 and return focus to the existing real-task quality/capability ladder.

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

### BMAD observability-lane check — 2026-09-03

The existing `/stats` correlation path was pursued without runtime mutation.

- `fcm-private-dev` is an nginx container sharing the WireGuard network namespace, not the Node router process itself; its logs only show nginx startup in the inspected window.
- `ingress-dev` is also nginx-facing and does not expose the router `/stats` surface directly.
- Personal Edge has an active WireGuard peer (`10.8.0.1/32`) with a recent handshake, but neither `10.8.0.1:4019` nor the canonical FCM DEV daemon port range beginning at `29280` is reachable from the edge. No local FCM/router process is listening on the Personal Edge host either.
- Exact FCM source confirms DEV daemon behavior: `runRouterDaemon()` defaults to `127.0.0.1`, uses DEV ports `29280..29289`, and only becomes remotely reachable when `FCM_HOST` is set to another bind address. Therefore a healthy WireGuard tunnel does not imply that `/stats` is network-exposed.
- No current canonical runtime evidence identifies a remote `/stats` bind or port-forward. Guessing a host/port or introducing a debug proxy would expand scope and is rejected.

Updated gate: this is now an `OBSERVABILITY_GAP` with a source-backed explanation — the diagnostic endpoint exists but its CURRENT deployment location/exposure is not evidenced. The next safe discriminator must therefore come from an already-authorized runtime-local read of the router port/log/requestLog, or from the same bounded SWE L2 retry capturing the FCM `x-request-id` through the existing OpenCode wrapper artifacts. Do not mutate FCM streaming behavior until one of those paths correlates the failure.

### BMAD L2 retry + quick-dev result — 2026-09-03

The compulsory same-level retry was executed through the CURRENT authenticated AgentField control plane with no provider/config/redeploy mutation.

- Parent: `exec_20260903_120353_r0nor5re`, run `run_20260903_120353_fvi6k9nh`, target `swe-planner.implement_issue`, repo `/tmp/fcm-calculator-live`, base `master`.
- The run completed in ~102s with a valid structured top-level result. The previous FCM/OpenCode `Stream error occurred` / missing-output failure did **not** reproduce on this retry, so FCM streaming must not be patched from the earlier correlation alone.
- Coder produced a bounded 3-file change and commit `662ba6a5533cf9c597a3b925930037c17b1e65bf`; reviewer then correctly found unmet CLI acceptance: power/modulo were not wired into parser choices and JSON output was not implemented. Result was `success=false`, `outcome=failed_unrecoverable`, one iteration.
- This moved the first failed gate from provider transport to SWE semantic self-repair. Source inspection on canonical SWE `dev` showed a contract mismatch: the reviewer prompt defines `blocking=true` as "must be fixed before merge" including missing core functionality, while `runDefaultPath` converted any blocking review directly to action `block`, and the outer loop terminates `block` immediately instead of feeding feedback into the next bounded coder iteration.
- `bmad-quick-dev` applied the smallest durable source fix on `dev`: commit `2310dbd418499a9a230ffaed39f613d0b0c129cf` changes default-path blocking review to `fix` so bounded iteration/stuck/exhaustion guards own recoverability; commit `a21f1ce8f345cc7e294fbb5821d7db44f441b954` replaces the immediate-fail regression with `TestBlockingReviewRetriesAndCanRecover`, requiring blocking iteration 1 followed by successful iteration 2.
- Stale-safe reread and secret scans passed for both source edits. Exact-source test execution is currently **VALIDATION_BLOCKER**: Coding Station, the existing runtime target, and Personal Edge do not provide `go`; repository CI has no `workflow_dispatch` and does not auto-run on `dev`.
- A bounded validation ref `copilot/bmad-validate-a21f1ce` was created at exact SHA `a21f1ce8...` to reuse the existing CI branch trigger without changing code. Creating the ref did not produce a workflow run, so no CI PASS is claimed and no further release/PR action was taken.

Current gate is therefore **SOURCE FIX DURABLE / VALIDATION + RUNTIME MATERIALIZATION PENDING**. Do not advance L2 yet. Next mandatory move: run the canonical Go regression on exact source `a21f1ce8...` using an existing authorized Go-capable runner; if PASS, materialize this exact accepted delta into CURRENT `/src/swe-af`, reload only `swe-planner` if required, and repeat the identical L2 canary. If no Go-capable authorized runner exists, keep status `VALIDATION_BLOCKER` rather than treating the code review as PASS.

### Fresh functional check — 2026-09-03 13:03Z

A new identical L2 canary was executed to test whether CURRENT runtime behavior had already converged without further materialization:

- Parent execution `exec_20260903_130335_s8486jav`, run `run_20260903_130335_p7riw9h0`, duration ~91.6s.
- AgentField transport status was `succeeded`, but semantic result was `success=false`, `outcome=failed_unrecoverable`.
- The failure returned at coder iteration 1: `Schema validation failed after 2 retry attempt(s). Last error: The output file was NOT created.`
- The run still produced commit `04905bd385090a465bbf606235edb4c3046c8dda` and a non-empty diff before losing the structured result; this again proves transport/file mutation is not semantic acceptance.
- Because this attempt failed before reviewer execution, it does not prove whether the durable blocking-review retry fix is loaded in CURRENT `swe-planner`.

Verdict: CURRENT L2 is **not reliably working yet**. The provider/structured-output failure is intermittent: the immediately preceding L2 reached a valid reviewer result, while this fresh retry regressed to missing structured output. Keep the FCM/provider boundary open and keep runtime materialization of the self-repair fix unproven. No L3 advancement.

### Structured-output root cause and recovery design — 2026-09-03

BMAD routing for this batch: `bmad-help` → `bmad-testarch-test-design` for the evidence/risk gate → `bmad-quick-dev` for the smallest source fix. No duplicate spec/test-plan artifact was created; this `PLAN.md` remains the single project SoT.

Fresh forensic evidence for `exec_20260903_130335_s8486jav` changes the owning failure classification from generic `PROVIDER / ACI_HARNESS` suspicion to a concrete structured-result recovery defect for this execution:

- The first OpenCode attempt completed repository work and emitted a useful near-schema `CoderResult` JSON object in an assistant text event. It reported `complete=true`, `tests_passed=true`, a non-empty `files_changed` list, summary/test summary/codebase learnings, and `agent_retro` as a plain string.
- `.agentfield_output.json` was not created. The weak-model shape `agent_retro:string` conflicts with the typed Go field `map[string]any`; other omitted CoderResult fields are defaulted by the canonical Pydantic model and are not independently invalid.
- AgentField then executed schema retries. Later retries emitted unrelated generic OpenCode text. The pinned SDK exposes only the latest retry text as `Result.Result`, but preserves prior OpenCode events in `Result.Messages`. The useful first-attempt JSON therefore remained available but SWE's wrapper did not inspect it.
- The pinned SDK text extractor also counts `{`/`}` without JSON-string awareness, so valid coder summaries/code fragments containing unmatched braces inside quoted strings can evade recovery.

Decision / safety contract:

1. Do **not** weaken semantic validation and do **not** convert repository side effects/commit existence into success. A recovered result must still satisfy the exact generated JSON Schema.
2. For `FailureSchema` / `FailureNoOutput` only, search candidate structured results in assistant text surfaces from the latest result plus prior `Result.Messages`; exclude tool outputs so arbitrary repository/file JSON cannot be mistaken for the orchestration envelope.
3. Extract candidate JSON objects with a string/escape-aware scanner. Validate each candidate against the exact schema before accepting it.
4. Apply one narrow weak-model compatibility normalization for `CoderResult`: `agent_retro:"text"` → `agent_retro:{"summary":"text"}`; leave unrelated type errors strict. Materialize normal CoderResult defaults before final schema validation because the canonical Pydantic model defines those fields with defaults.
5. Provider/API/transport failures remain fail-closed even when text happens to resemble JSON; this recovery path must not mask `Stream error occurred` or other upstream failures.
6. Longer-term reliability direction after L2 PASS: prefer provider/OpenCode schema-constrained finalization for the packaging phase, while keeping domain/semantic validation separate. Constrained decoding improves structural compliance but is not itself semantic correctness.

The container-first implementation was expanded and closed in CURRENT `/src/swe-af` across the exact owning slice below; no GitHub-first coding or whole-workforce redeploy was used:

- `go/internal/harnessx/run.go`: assistant-event recovery from `Result` + prior `Messages`, string/escape-aware JSON extraction, exact JSON-Schema revalidation, watchdog salvage for already-valid structured output, and OpenCode `SchemaMode="incremental"` so schema retries retain task context.
- `go/internal/harnessx/harnessx_test.go`: regressions for quoted braces, schema-invalid rejection, provider-failure fail-closed behavior, observed retry-overwrite, watchdog salvage, generic transport fail-closed behavior, and incremental OpenCode schema mode.
- `go/internal/schemas/defaults.go`: narrow `agent_retro:string` → `agent_retro:{"summary":...}` normalization in `CoderResult.UnmarshalJSON`; unrelated type errors remain strict.
- `go/internal/roles/coding/coding.go`: coder schema/transport failures fail closed; one bounded same-worktree retry is allowed only for recoverable no-progress/stream conditions so completed filesystem work is not discarded; reviewer output loss also fails closed instead of auto-approving.
- `go/internal/roles/coding/coding_test.go`: fail-closed coder/reviewer regressions, provider-error preservation, and same-worktree no-progress retry coverage.
- `go/internal/prompts/coding/verifier.go`: verifier uses the PRD/acceptance criteria already embedded in the prompt and no longer instructs OpenCode to read artifact paths outside the current worktree.
- `go/internal/prompts/coding/testdata/task_verifier_a.txt`, `task_verifier_b.txt`: canonical prompt goldens updated for the bounded verifier contract.

Validation and runtime evidence:

- Exact live source targeted validation PASS in the permanent workforce: `go test ./internal/harnessx ./internal/schemas ./internal/roles/coding ./internal/roles/advisor ./internal/prompts/coding`; matching `go vet` PASS.
- The no-progress behavior is now **progress/evidence-driven**, not governed by a hard total workflow budget: a no-progress watchdog remains as a liveness guard, but if assistant text already contains an exact-schema-valid structured result it is salvaged; otherwise coder may retry once in the same worktree and generic/provider errors remain fail closed.
- CURRENT loaded planner is `/afhome/packages/swe-planner/bin/swe-planner` SHA256 `ff662a38022d29ef641506def62fce3cf6ecdc02f7ef478bcaaef3d6a3aa12d0`, PID `167730`, started `2026-09-03T19:24:35Z`; `swe-pro` is PID `167737`.
- Final same L2 canary result `/tmp/l2-canary-verifier-fix.json`: `success=true`, `outcome=completed`, build `9a7e6fb6`, one bounded 3-file task commit `52fc4a500a6a3e6e4218a47382ffe53159320e0e`, reviewer `approved=true` / `blocking=false`, verifier `passed=true`, **5/5 acceptance criteria PASS**, and independent `python3 -m unittest test_calculator.py` **22 tests PASS**.
- Exact runtime → Git durability check PASS for all eight source files: Git blob OIDs on `dev` match `git hash-object` of CURRENT `/src/swe-af` byte-for-byte. Pre-PLAN code head after canonicalization is `52bdb640fa9f2a796dc423a96cdec971d883172d`; this PLAN-only write-back may advance branch HEAD without changing those product blobs.

**L2 structured-output reliability gate: PASS / CLOSED.** The original intermittent `The output file was NOT created` failure is no longer an open product gate for this bounded path. Advance only to the next documented recovery/failure gate; do not regress to GitHub-first debugging or reintroduce hard total wall-clock budgets as an acceptance primitive.

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

Status: PASS / CLOSED — FAILURE/RECOVERY CONTRACT PROVEN

Prerequisite: L2 structured-output reliability gate PASS/CLOSED above.

Fresh CURRENT anti-drift readback — 2026-09-04:
- Permanent workforce remains `1437789b5c4debc061992bec718132dcbccc66ac67e0b9e45ce1eea5b651c9e7`; live `/src/swe-af` identity remains baseline `58c4e0d19081bc52363c120b7963a34cebb1e894` plus the preserved dirty source delta.
- `/afhome/installed.yaml` is stale execution metadata and must not be used as CURRENT PID authority. Historical PIDs `8059`, `167730`, `230927`, `239824`, and `241940` are superseded/zombie generations. CURRENT loaded `swe-planner` after the resume-idempotency batch is PID `243734`; binary SHA256 is `74501716485eedc14964086cde655ac77cb505f6bfc8986a0987f398e7ec6832`. Logs prove fresh `node.register.complete` / `agent.initialize.complete` for `swe-planner` and embedded `swe-pro` after the process-only reload on 2026-09-04.
- Root `ERRORS.md` exists in durable project history but not in the preserved runtime baseline; do not synthesize a duplicate inside `/src` merely to remove metadata drift.

Structured-output architecture spine (BMAD `bmad-architecture` + `bmad-testarch-test-design` + `bmad-quick-dev`):
1. All role structured calls continue through the existing `harnessx.Run[T]` choke point. `Run[T]` is now deliberately thin: reflect schema, inject run credentials, call `executeStructured[T]`.
2. `executeStructured` is the single SWE-owned policy boundary for weak-model structured output: OpenCode incremental schema mode, exact-schema validation, assistant-event recovery, exact output-file watchdog salvage, default seeding, and provider/fatal fail-closed behavior. Role code no longer owns JSON retry strategy.
3. Machine contract and business semantics stay separate. The centralized layer guarantees structural/schema contract only; coder/reviewer/verifier semantic acceptance remains with the owning role/application checks.
4. Pareto decision: use the pinned SDK's existing incremental schema mode as the first field-by-field repair engine instead of re-implementing a second generic repair loop. Add an SWE-owned partial/group repair loop only when a real failure survives incremental mode with recoverable partial state.
5. Desired long-term file layout is a dedicated adjacent structured-output module. DEV typed `file.create` currently denies new source files as `out_of_scope`; this was not bypassed with shell writes. Therefore the accepted live refactor temporarily houses the controller in existing adjacent `internal/harnessx/schema.go`, keeping `run.go` merge-friendly. A later SourceLoop/canonicalization step may split it to `structured.go` without changing the seam or behavior.

Centralization batch evidence:
- Live source changes are confined to `go/internal/harnessx/run.go`, `schema.go`, and the existing `harnessx_test.go` regression surface on top of the prior accepted runtime delta; no GitHub product-code edit and no container redeploy were used.
- `gofmt -d` clean.
- `go test ./internal/harnessx ./internal/schemas ./internal/roles/coding ./internal/roles/advisor ./internal/prompts/coding` PASS; matching `go vet` PASS.
- Contract regression proves `RoleOptions` no longer owns `SchemaMode`; the central controller injects `SchemaMode="incremental"` for OpenCode immediately before the base harness call.
- Functional canary `exec_20260904_171701_d2vmapuh` / `run_20260904_171701_qfs8mep3` on loaded PID `239824` completed `succeeded` in `153831ms`, returned `passed=true`, 2/2 criteria PASS, and produced a complete exact-schema verifier result. During execution the output file was observed being built incrementally before terminal success.

Batch 3 DoD:
- nonexistent/invalid task fails closed or abstains without repo damage. **PASS** — `exec_20260904_134916_1cbcid8b`; root HEAD/status unchanged.
- unrelated worktree files are preserved. **PASS** — `exec_20260904_135618_p78m3o2b`; unrelated `.pyc` SHA256 unchanged.
- centralized weak-model structured-output policy has one role-agnostic seam and real runtime proof. **PASS** — source tests/vet + `exec_20260904_171701_d2vmapuh`.
- bounded run interruption + resume is idempotent. **PASS 2026-09-04**: source inspection found a crash window where `executeLevel` completed but the checkpoint was written before `CompletedIssues/FailedIssues/SkippedIssues` were recorded. Live fix adds an immediate checkpoint after recording the level-result barrier. Deterministic regression `TestResumeAfterCancellationDoesNotRepeatCompletedIssue` cancels after issue `a` completes, proves the interrupted checkpoint contains `a`, then resumes and invokes only `b`; `go test ./internal/dag`, `go test ./internal/orch`, and matching `go vet ./internal/dag ./internal/orch` PASS. Runtime smoke `exec_20260904_174531_yv9d2ym9` / `run_20260904_174531_rjp812gl` resumed a synthetic interrupted checkpoint in 282ms: only child `swe-planner.execute` ran (~6ms), no coder/issue reasoner call occurred, `completed_issues` remained exactly `[a]`, `current_level` advanced 0→1, and sentinel SHA256 stayed `d4383fc33176254b2abac7cb3c14bcaae5701681537db3d69adf740b1feb1b9c`.
- stale SHA/branch advance is detected and fails closed. **PASS 2026-09-04**: `implement_issue` now accepts optional `expected_base_sha` and compares it immediately after resolving `base_branch`, before artifacts/worktree/LLM side effects. Source regressions `go test ./internal/issue ./internal/node` and matching `go vet` PASS. Runtime canary `exec_20260904_173012_503eoaro` used disposable `/tmp/swe-stale-sha-canary`: expected `f498992e8eaaba72595494ecbd30974c63cbda64`, current advanced `52ae4c29cd9b5604c07cbea8988709bdf2e064b1`; execution failed in ~11ms with explicit stale-base error. Independent readback proved HEAD remained `52ae4c29...`, status stayed clean, and no `issue/*` branch/worktree was created.

Batch 3 is closed. Next 30-minute Pareto move is Batch 4 durability/SourceLoop: capture only accepted source files from CURRENT `/src/swe-af`, reject generated/noise/stale state, obtain a stale-safe candidate against `fork/dev`, publish/accept through the canonical SourceLoop lane, record exact durable SHA, then rerun repository tests on that exact durable source before any materialization step.

### Batch 4 — Durability / SourceLoop gate

Status: IN PROGRESS — STALE/NOISY RUNTIME CAPTURE REJECTED; CLEAN CANDIDATE REQUIRED

Prerequisite: Batch 3 PASS/CLOSED above.

Fresh CURRENT anti-drift — 2026-09-04:
- CURRENT loaded planner is PID `246239`, binary SHA256 `952748c1b92c60a1dc60541bc225a2c0444a682b60d57e05e30d024ff51a03c2`, built from the live source after restoring the durable `blocking review -> fix/retry` contract and retaining the newer fail-closed reviewer behavior.
- Whole-worktree comparison must be against CURRENT durable `dev`, not runtime Git status against baseline `58c4e0...`; the latter misses files that advanced in Git but were never materialized into the preserved runtime baseline.
- That anti-drift comparison exposed and repaired a split-brain in `internal/coding/loop.go` / `loop_test.go`, then exposed a stale level-failure test oracle. After merging the intended semantics, combined `go test` and matching `go vet` PASS across coding, DAG, orch, harnessx, issue, node, schemas, coding/advisor roles, and coding prompts.
- Runtime-capture is healthy and automatically captured CURRENT `/src/swe-af` to PR #7, head `556d910369158c2e664deb51a84000a68240711f`, but PR #7 is **not an acceptable durability candidate**: it is a 21-file full runtime snapshot rooted at stale runtime base `58c4e0...`, its PR base metadata predates CURRENT `dev`, and it mixes accepted recovery/structured-output changes with older provider/planning deltas that have no current Batch-4 acceptance verdict.
- SourceLoop per-change journal is also incomplete for aggregate durability because some accepted patches were applied through the production typed live-patch lane and are not present in the DEV journal; requesting an artifact for a recent DEV change returned `capture_artifact_reference_invalid`. Therefore neither “latest journal event” nor PR #7 may be promoted blindly.
- Durable-only `ERRORS.md`, this `PLAN.md`, and `docs/LLM_PROVIDER_SECURITY_CONTRACT.md` are absent from the preserved runtime baseline and must be preserved; their apparent deletion in a raw runtime snapshot is noise/drift, not accepted source delta.

Accepted durability allowlist for the current recovery batch is the merged SWE slice only: `go/internal/coding/loop.go`, `go/internal/dag/executor.go`, `go/internal/dag/executor_test.go`, `go/internal/harnessx/harnessx_test.go`, `go/internal/harnessx/run.go`, `go/internal/harnessx/schema.go`, `go/internal/issue/build.go`, `go/internal/issue/build_test.go`, and `go/internal/node/register.go`. `loop_test.go` already matches durable `dev` after anti-drift repair and needs no publication delta. Older provider/planning files (`agentfield-package.yaml`, node AI config/source, `roles/planning/planning.go`, `opencode.json`) stay out until separately accepted.

DoD:
- only accepted source delta captured; runtime/provider/planning noise excluded. **IN PROGRESS**
- stale/noisy capture fails closed. **PASS** — PR #7 explicitly rejected as promotion input.
- accepted delta reaches `fork/dev`. **PENDING**
- exact `WORKING_DEV_SHA` is recorded. **PENDING**
- repo tests pass on that exact durable SHA. **PENDING**

Next 30-minute Pareto move: create one clean stale-safe publication candidate from CURRENT `dev` containing only the allowlisted live bytes above, validate its diff against both CURRENT `dev` and CURRENT live source, publish through the canonical repository publication lane, then run the same package test/vet gate on the exact candidate/durable SHA before acceptance. No product logic is to be edited in GitHub during this step; this is Lane-B canonicalization of already-tested live bytes.

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
