# SWE-AF Project Plan (Source of Truth)

Status: active
Last reconciled: 2026-09-06
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

Historical baseline captured from live readback on 2026-09-01. The dated evidence in this section is retained for forensic chronology; it does **not** override the newer `Current Stage`, capability-audit, Governor, and 2026-09-06 reconciliation statements below. Old PIDs/provider wording in this historical block must not be treated as CURRENT without fresh runtime readback.

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

REAL-TASK QUALITY ACCEPTANCE / SYSTEM CAPABILITY CERTIFICATION / ANTI-DRIFT. OpenClaw Governor is a parallel human-control workstream, not the primary blocker for SWE capability certification.

L2 semantic/structured-output acceptance and Batch 3 failure/recovery are PASS/CLOSED. The current P0 product gate is the full `implement_issue` reviewer -> repair -> verifier -> independent-oracle loop; reviewer defect detection is proven, but the whole loop is not yet L3 PASS. After that, proceed through one real two-issue DAG, one governed `plan -> execute/build` feature, one delivery-path task, Fast/Pro comparative canaries, clean-environment/concurrency acceptance, then durability/exact-SHA closure and repeated reliability. Batch 4 durability remains open but is deliberately deferred behind these real capability gates.

2026-09-06 anti-drift readback: CURRENT loaded `swe-planner` is PID `389456`, binary SHA256 `093cdfd3adac902c1f96e1de11672e92316051d4e3c3acbb1cfcc8022ea3b521`, callback `http://workforce:8800`; the latest reload was process-only. Live source remains `/src/swe-af` ahead of durable Git and carries the reviewer/verifier workstream delta. Durable `dev` HEAD immediately before this reconciliation is `639b8b68e2d55faebb0629ff3585a929c500d923`; PLAN-only commits must not be mistaken for loaded product identity. The currently loaded reviewer contract includes the simple single-line delimiter-in-literal discriminator, non-vacuous zero-output failure, `Edit` in reviewer tools, an explicit verdict-file-only Write/Edit exception, and an OpenCode-only permission overlay denying `task` subagents plus `external_directory` access for the reviewer.

### EvalGuard #3 reviewer/verifier reconciliation — 2026-09-06
BMAD route for this batch: `bmad-help -> bmad-testarch-test-design -> bmad-quick-dev`. The authoritative P0 remains the full `implement_issue -> coder -> reviewer -> repair -> verifier -> independent oracle` path; `implement_issue` remains **L3 PARTIAL**.
- Bad coder commit `e64bd39965564c03c6f316f3dff405efe79f468a` received a false reviewer+verifier PASS. The independent oracle proved semantic corruption of `#` inside Python string literals, so that commit is explicit negative evidence, not a success baseline.
- CURRENT live reviewer contract exposes `Read, Write, Edit, Glob, Grep, Bash`; source/test files remain non-writable by reviewer policy, while Write/Edit are allowed and required only for the structured verdict file. After a reproduced BLOCKING defect the next tool action must write the verdict; further tests/reads/static checks are forbidden. For comment/token stripping, the first discriminator is a simple single-line ordinary string containing the delimiter plus a real external comment; multiline/triple-quote cases come only after the simple case.
- Targeted live source gates are PASS: `gofmt -d`, `go test ./internal/prompts/coding ./internal/roles/coding`, and `go vet ./internal/prompts/coding ./internal/roles/coding`. Reviewer regression coverage protects both the risk discriminator and structured-verdict termination contract.
- Canary `exec_20260906_091757_9ozebdu8` independently reproduced the semantic blocker (`"#FF0000"` + real `# comment`, `num muts: 0`) but then exposed the old termination/tool-contract conflict by continuing to test instead of writing the verdict. That contract is now fixed.
- Canary `exec_20260906_092738_83036e8b` is **provider-contaminated** and was stopped: FCM repeatedly returned `AI_APICallError: All routed models failed for set: fast-coding`, causing repeated read/test fallback behavior. No further SWE patch was justified from that run.
- The FCM-owned blocker has now been repaired separately and live-verified: provider/key 429 cooldown is honored by real-inference health probes, OPEN/paused routes are not repeatedly probed, startup probes serialize siblings within a provider, balanced cadence persists at 120s, full FCM live suite is 834 PASS / 0 FAIL / 2 intentional skips across 160 suites, live plain-text + `json_schema` canaries are 2/2 PASS, and post-reload recurring probes no longer show the former OpenRouter/Requesty 429 storm. This removes the known provider-noise cause but does **not** itself close SWE L3.
- Fresh retry `exec_20260906_105150_0zqvwbo1` ran after the FCM recovery and is therefore clean enough to reclassify the remaining gate back to SWE/reviewer liveness. The reviewer independently reproduced the required simple semantic blocker on commit `e64bd39965564c03c6f316f3dff405efe79f468a`: its focused `DocstringStrippingMutator` discriminator returned `ERROR: No mutations generated - this is a BLOCKING failure`. Despite that, no `.agentfield_output.json` appeared and no `Write/Edit` verdict action followed within the bounded observation window; the execution was cancelled at `2026-09-06T10:57:42Z`. This is **ACI_HARNESS / reviewer termination L3 PARTIAL**, not a provider failure and not a semantic false-pass.
- The same retry also exposed a lower-priority prompt/path inefficiency before the useful attempt: an initial OpenCode pass tried stale external `/home/user/repos/FreakyAdy/EvalGuard/...` paths, triggered auto-rejected `external_directory` permissions, and was retried in the correct worktree.
- Follow-up diagnosis proved that AgentField's OpenCode adapter does not enforce `harness.Options.Tools`; OpenCode therefore allowed a hidden `task`/`explore` subagent which then blocked on `external_directory` permission. SWE now applies a reviewer-only `OPENCODE_CONFIG_CONTENT` overlay denying both `task` and `external_directory`. Deterministic role tests PASS and canary `exec_20260906_111051_je9w0wr8` confirmed the subagent/path stall was removed: reviewer stayed in the intended worktree and used direct bash/read actions only.
- That same canary cannot close reviewer termination because FCM again became unstable during ordinary reviewer model steps: repeated `AI_APICallError: All routed models failed for set: fast-coding` occurred after the sandbox improvement, so the execution was cancelled at `2026-09-06T11:14:25Z` as provider-contaminated. CURRENT FCM `fast-coding` has drifted to 33 configured candidate routes; prior 20-route inventory is stale. Do not patch SWE around this provider failure.
- SWE product source from this reviewer/verifier batch is intentionally **not** durable yet. Do not publish it until the unchanged bad-commit reviewer canary cleanly returns `approved=false`, `blocking=true` and terminates by writing its verdict without new SWE edits.

Next 30-minute Pareto move: rerun the unchanged bad-commit reviewer only when FCM can sustain the call. The reviewer sandbox/subagent drift is now fixed, so the remaining clean discriminator is termination after a reproduced blocker. If `approved=false`, `blocking=true` returns cleanly, run verifier on the same case, then the full `implement_issue` repair path, then the independent/hidden oracle. If FCM fails again, classify PROVIDER and continue deterministic SWE certification rather than changing reviewer logic from contaminated evidence.

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

Authoritative Pareto certification spine: (1) close the SAME EvalGuard #3 full-loop quality/liveness gate; (2) one small two-issue real DAG task proving dependency order, parallel/level behavior, failure threshold, checkpoint and cleanup in one run; (3) one real governed `plan -> execute/build` task; (4) one delivery-path task covering resolve/CI/PR only after local semantic gates pass; (5) comparative Classic/Fast/Pro canaries on the same bounded task class; (6) clean-environment/bootstrap plus multi-project/concurrency acceptance; (7) durability/exact-SHA materialization; (8) repeated reliability across the proven spine. Historical batch chronology below is evidence history and does not override this order.

### Human operating model — maximum-output mode

Human/SWE boundary is explicit and sparse. The human is the governor of intent and gates, not a participant in every agent step.

1. **Goal gate — human owns WHAT/WHY.** Supply a bounded goal/repo/constraints/acceptance criteria. For one known issue use `implement_issue`; for feature-level work use `plan` first, then `execute`/`build` only after plan quality is acceptable. Do not micromanage coder/reviewer prompts during normal execution.
2. **Plan gate — human owns scope/blast/architecture approval.** In the CURRENT deployment `SWE_OPENCLAW_HITL=1` is the active human surface; legacy HAX remains compatible when `HAX_API_KEY` exists. `build` pauses before execution and exposes a bounded PRD/architecture/issues decision packet. Human decisions are `approved`, `request_changes`, `rejected`, or expiry. `request_changes` is not cosmetic: SWE re-runs Architect -> Tech Lead -> Sprint Planner with the human feedback, bounded by `MaxPlanRevisionIterations`.
3. **Clarification/credential gate — SWE may ask, human answers only missing authority/context.** HITL `ask_user` genuinely transitions the execution to `waiting` until webhook resume. In the CURRENT deployment the question is rendered as an OpenClaw/Telegram decision card; HAX may render the same logical request as a form when configured. Prior answers are injected back so the same question should not be re-asked. Environment scout may request scoped service credentials; secrets belong in the scoped credential store/environment, not in result text.
4. **Exception gate — human intervenes on escalation, irreversible/external actions, or exhausted autonomous recovery.** Normal coder -> reviewer/QA -> fix iterations, issue-advisor actions, DAG checkpoint/resume and replanning should remain autonomous. Human intervention is warranted when scope/architecture must change, credentials/permissions are missing, risk is irreversible/external, or bounded recovery is exhausted.
5. **Acceptance gate — human owns DONE criteria, independent oracle and delivery authorization.** Agent-written tests and reviewer approval are supporting evidence, not sole proof. Require project-native tests, risk-based independent acceptance, repo hygiene/provenance, and delivery/PR/CI evidence appropriate to the task before calling the work done.

Recommended operating profiles:
- **Precision / default:** `implement_issue` for an already-scoped issue, classic or proven Pro coding engine, independent reviewer, bounded repair, verifier; human only at exceptions/final acceptance.
- **Feature / governed:** `plan` -> human plan review -> `execute`/`build`; use the CURRENT OpenClaw/Telegram approval path for architecture/scope-sensitive work. HAX remains an optional compatible form UI, not the active dependency in this deployment.
- **Fast:** `swe-fast` only after its own L2/L3 certification; use for low-risk bounded tasks where turnaround dominates depth.
- **Pro:** `SWE_PRO_ENGINE=1` only after Pro runtime acceptance; tune `SWE_PRO_VARIANT` and `SWE_PRO_MAX_COST` as effort/cost controls, not as correctness substitutes.

Maximum-output Pareto policy: spend human attention at decision boundaries, not execution steps; spend test budget on independent oracles for high-risk semantics, not exhaustive duplicated self-tests; test capability groups end-to-end rather than 32 reasoners one by one; keep a cheap fail-fast ladder (schema/unit -> package test/vet -> live reasoner -> real-repo oracle -> delivery).

Fresh full-loop anti-drift: the repeated EvalGuard #3 execution `exec_20260905_075042_xi2na5he` completed `status=succeeded` at transport but `success=false`, `outcome=error`; coder produced commit `d7402d4aa01cf80f0efeeb1391849e39ca436e52`, while reviewer execution exceeded its 1800s agent-call limit. During that review the independent reviewer had already reproduced additional source-transformation failures. Therefore reviewer defect-detection quality improved, but whole-loop reviewer time/termination efficiency remains the nearest quality/liveness gate; do not call `implement_issue` L3 PASS yet.

### Governor communication decision — OpenClaw / Telegram

Decision: reuse the existing OpenClaw installation as the messaging/control bridge, but move the SWE cockpit to **one dedicated Telegram bot/account `swe` bound to OpenClaw agent `devteam`**. Do not create one bot per project and do not create a parallel notification service/database. The default Telegram bot remains Jarvis. One SWE bot multiplexes multiple projects using explicit project identity plus execution/request correlation.

CURRENT OpenClaw/Telegram control-plane readback on 2026-09-06:
- Telegram `default` remains Jarvis (`@ass_nikita_bot`) and is configured/running/connected/polling with probe OK. Named account `swe` is also configured/running/connected/polling with probe OK as `@claude_code_nikita_bot` (`SWE Governor`, bot id `8787962440`). Both bots report topics disabled, so the current UX is flat DM + exact reply correlation.
- Binding readback is exactly `telegram accountId=swe -> agentId=devteam`; the default bot remains outside that binding. `TELEGRAM_BOT_TOKEN_SWE` is present in the OpenClaw runtime and `SWE_GOVERNOR_TELEGRAM_ACCOUNT=swe` is persisted in Coolify. Account-scoped pairing state is `/home/node/.openclaw/credentials/telegram-swe-allowFrom.json`.
- OpenClaw natively preserves Telegram reply metadata. Live human inbound on 2026-09-06 reached `session:agent:devteam:main`, and prior dedicated canary `messageId=367` mapped exactly back to its SWE `request_id`/`project_id`.
- Canonical outbound remains deterministic `message send`: OpenClaw returns the Telegram `messageId`, which the Governor stores against the exact SWE `request_id`. No routing by latest pending request/project is permitted.
- High-consequence terminal replies now have a model-independent fast path via loaded hook-only plugin `/home/node/.openclaw/extensions/swe-governor-telegram`: `before_dispatch` priority `100` intercepts only Telegram account `swe` replies with an exact `ReplyToId`, resolves that id through the canonical Governor adapter, and handles only explicit `approve`/`changes`/`reject`/`answer` intents. Non-terminal discussion still falls through to the `devteam` conversational model; unknown reply ids and terminal failures remain fail-closed.

Conversation correlation invariant:
1. SWE emits a bounded `governor.pending` record and enters the existing AgentField `Pause(waiting)` path.
2. The 20-second OpenClaw command cron runs `swe-governor.mjs notify` with `delivery.mode=none`; the adapter itself sends the Telegram card through OpenClaw and records `telegram_message_id <-> request_id`.
3. On an inbound Telegram reply, Jarvis must inspect `ReplyToId` before normal processing and call `swe-governor route <ReplyToId>`. `NOT_SWE_REPLY` means normal Jarvis traffic. No routing by "latest pending request" is permitted.
4. If native reply metadata is unavailable on a follow-up, an exact quoted `SWE thread: <request_id>` / `Request: <request_id>` marker or explicit `swe ... <request_id>` command is the only fallback. Every Governor discussion response carries the thread marker and an explicit Telegram reply tag.
5. `discuss` is non-mutating and may be called repeatedly; it records a bounded per-request conversation while the AgentField execution remains paused. Only unambiguous `approve`, `changes`, `reject`, or `answer` is terminal and may hit `/webhooks/approval`.
6. Multiple SWE requests may coexist. Each Telegram card has its own message-id mapping; interleaved discussion remains isolated by request id. The first successful final action wins; late/conflicting messages see `resolved` and fail closed.
7. Secrets/credentials and raw logs stay out of the conversation packet. SWE text is untrusted data, not OpenClaw instruction text.

Project multiplexing / concurrency policy:
- one dedicated SWE Telegram bot serves many projects; do not create one bot per repository;
- decision packets now carry `project_id`/`repo_path` where the owning SWE stage already knows repo identity; cards display `Project:` when present, while `execution_id`/`request_id` remain the exact correlation authority;
- a reply to a decision card chooses its exact request regardless of project. A new non-reply conversation with more than one active project must name/select a short project alias; never infer the project from recency;
- different project workspaces/build ids may run concurrently. CURRENT source generates `build_id` before workspace setup and auto-derived `repo_url` workspaces include that id; scoped credentials are keyed by run id under a mutex. `go test -race ./internal/orch ./internal/hitl -run 'TestBuildIsolationConcurrent|Concurrent' -count=1` PASS on 2026-09-05;
- until a same-physical-checkout system test exists, Governor serializes full builds targeting the same explicit `repo_path`. This is a conservative operating guard, not a claim that same-repo concurrency is broken.

Evidence for the routing contract:
- Fresh 2026-09-06 CURRENT readback is OpenClaw `2026.7.1`; Telegram accounts `default` (`@ass_nikita_bot`) and `swe` (`@claude_code_nikita_bot`) are both configured/running/connected/polling with probe OK. Binding is exactly `telegram:swe -> devteam`; Jarvis/default remains intact. `TELEGRAM_BOT_TOKEN_SWE` is present and `SWE_GOVERNOR_TELEGRAM_ACCOUNT=swe`; dedicated activation and live inbound channel isolation are no longer blocked.
- Fresh compatibility canary sent two synthetic cards through `default`: A `tgcanary-A-20260905` -> Telegram message `14077`, B `tgcanary-B-20260905` -> `14078`. `route(14077)` returned only A/project-A, `route(14078)` only B/project-B, unknown id returned `NOT_SWE_REPLY`; interleaved A1/B1/A2/B2 `discuss` turns stayed isolated and both requests remained pending with mock callback count `0`.
- A concurrent terminal-action canary exposed a real Governor race: before the fix, simultaneous `approve` + `reject` on the same pending request both succeeded and produced callback count `2`. `bmad-quick-dev` therefore changed only `/home/node/.openclaw/workspace/skills/swe-governor/scripts/swe-governor.mjs`, adding an exclusive stale-recoverable terminal-submit lock. Repeating the same race after the fix produced callback count `1`: the winner resolved the request and the competing action failed closed with `status=resolved`; lock cleanup was verified.
- Post-fix regression preserved exact A/B routing and multi-turn discussion with callback count `0`. `node --check` PASS. No OpenClaw reload/redeploy was needed because cron invokes a fresh Node process.
- Cron `swe-governor-decisions-v1`, id `c7ea3f43-2a02-4752-aaac-f254f9ba256f`, remains enabled every 20s, command `node /home/node/.openclaw/workspace/skills/swe-governor/scripts/swe-governor.mjs notify`, `delivery.mode=none`, `timeoutSeconds=15`; post-change runs remain `ok` with diagnostic `NO_REPLY`, and the final 2026-09-05 readback sampled `109ms` with errors/skips `0`.
- A dedicated-agent readiness defect was found before bot activation: `openclaw skills info swe-governor --agent devteam --json` returned `not found` even though the canonical Governor skill was eligible for `main`. A plain cross-workspace symlink was tested and rejected by the CURRENT resolver, then rolled back. The supported OpenClaw mechanism was applied instead: `skills.load.allowSymlinkTargets=["/home/node/.openclaw/workspace/skills/swe-governor"]` plus `/home/node/.openclaw/workspace-devteam/skills/swe-governor -> /home/node/.openclaw/workspace/skills/swe-governor`. `openclaw config validate` PASS and `skills info ... --agent devteam` now reports `eligible=true`, `modelVisible=true`, `userInvocable=true`; OpenClaw reported that no gateway restart was required.
- `/home/node/.openclaw/workspace-devteam/AGENTS.md` now owns the dedicated SWE-chat routing invariant: exact `ReplyToId` -> exact Governor request, `discuss` remains non-terminal, ambiguous/unknown messages cannot commit, first successful final action wins, and multi-project non-reply traffic must select a project/request rather than infer recency. CURRENT SHA256 is `90f26301d29c31c7188e51de1f559239c6dc78ceaac5bac3447048b2a6db3f4f`.
- The original `notify` behavior could send up to three pending cards sequentially; a three-card canary delivered all three but the outer command exceeded 30s, which is incompatible with the cron's 15s command budget. A concurrent-send experiment was then rejected: all three OpenClaw CLI sends timed out at ~15.2s with no message-id mappings. The accepted Pareto fix sends exactly the oldest one pending card per cron tick and persists its returned Telegram message id before exit. Repeating the timing canary succeeded in `10503ms`, mapping `canary-notify-P-20260905 -> 14082` while leaving the other two requests unsent for later ticks. This bounds one notification inside the 15s cron budget and avoids concurrent CLI-send races; it deliberately trades simultaneous backlog drain for one-card-per-tick reliability.
- Follow-up isolated canaries independently reconfirmed exact routing and discussion: `canary-gov-A-20260905 -> 14079`, `canary-gov-B-20260905 -> 14080`, `canary-gov-C-20260905 -> 14081`; A/B exact route PASS, unknown reply id -> `NOT_SWE_REPLY`, interleaved A1/B1/A2/B2 stayed isolated and pending with callback count `0`, and a local mock plan-review final action produced exactly one callback while a late conflicting action failed closed. A separate concurrent approve/reject race after the existing terminal lock also produced exactly one callback and final `resolved` state.
- `AGENTS.md` owns the shared-Jarvis compatibility pre-routing rule; `skills/swe-governor/SKILL.md` remains the single canonical Governor skill owner; the `devteam` workspace references it rather than copying it. CURRENT SHA256: main `AGENTS.md` `7f5ebb8f92b82b4577246723d3d75dec28c922bb0ef04fe7dfd493ba0c82a42d`, Governor skill `03ee024c1683410887571186c612d19e40d8b760b63431a6b880c2acbd7a20c0`, adapter `cd151545a1fc851f94ade73c056298f7bc5b0290d0fcb13a01c08c788d4f8637`, OpenClaw config `8d32430b4ad5441eccc870e97225f4c3b42439e7d833ec3628220d2b3039befb`.
- Canary Telegram messages `14077`, `14078`, `14079`, `14080`, `14081`, and `14082` are absent after cleanup/readback; isolated canary state/mock files were removed. No real SWE approval callback was used for these synthetic proofs.
- 2026-09-06 dedicated Telegram activation is now live. Coolify `AI Platform / production / ai-agent-stack` exposes `TELEGRAM_BOT_TOKEN_SWE` to the recreated OpenClaw container and `SWE_GOVERNOR_TELEGRAM_ACCOUNT=swe`. OpenClaw named account `swe` is configured as display name `SWE Assistant`, probes successfully as Telegram bot `@claude_code_nikita_bot` (bot id `8787962440`), and is `running=true` with `tokenStatus=available`; default/Jarvis `@ass_nikita_bot` remains `running=true` and probe OK. Binding readback is exactly `telegram accountId=swe -> agentId=devteam`. Account-scoped pairing state exists at `/home/node/.openclaw/credentials/telegram-swe-allowFrom.json` and was validated without exposing the user id.
- Dedicated-account Governor canary `canary-dedicated-swe-20260906` sent through account `swe` in `9963ms`, persisted Telegram message id `367`, and `route 367` returned exactly that request with `project_id=swe-af`; the Telegram canary message was then deleted successfully and local canary state removed. This proves dedicated outbound + exact reply correlation on the intended account.
- 2026-09-06 live human inbound closes the remaining channel-isolation proof. User message `Allo` at 14:22 reached the dedicated SWE Telegram bot and OpenClaw runtime logs show the inbound run under `session:agent:devteam:main`, proving `telegram:swe -> devteam` on a real human turn. The bot response then failed at the model/provider layer with external `gonka` HTTP `429` (`too many concurrent requests`) and surfaced the Telegram fallback `All models are temporarily rate-limited`. This is outside the Telegram/OpenClaw Governor ownership boundary and does not invalidate transport routing; no FCM/provider repair was attempted.

Optional richer UX remains a dedicated Telegram forum topic or topics-enabled DM, but CURRENT bot reports topics disabled. Do not add that infrastructure unless reply-based correlation proves insufficient in real use.

Implementation status — 2026-09-05:
- BMAD route for this vertical: `bmad-help` -> `bmad-architecture` -> `bmad-quick-dev`. Architectural outcome changed after CURRENT readback: HAX is not deployed/configured in the loaded SWE process, while AgentField `agent.Pause` already owns the real `waiting`/webhook-resume primitive. Therefore this deployment uses OpenClaw as the HITL surface directly; HAX compatibility remains intact when `HAX_API_KEY` exists.
- SWE source now has opt-in `SWE_OPENCLAW_HITL`. When HAX is absent and this flag is true, `ask_user` and plan approval create synthetic process-local request ids, emit bounded non-secret `governor.pending <json>` notes, call the existing AgentField pause path, and emit `governor.resolved <json>` after callback/timeout/error. Legacy HAX-disabled behavior remains unchanged when the flag is false.
- `build` outer gating is OpenClaw-aware as well as HAX-aware; a regression protects all three states: neither surface -> disabled, HAX -> enabled, OpenClaw -> enabled.
- Exact current SWE source slice for Governor/HITL includes `go/internal/hitl/ask_user.go`, `wrapper.go`, `hitl_test.go`, `go/internal/orch/approval_gate.go`, `approval_gate_test.go`, `build.go`, `build_test.go`, `go/internal/roles/planning/planning.go`, and `go/agentfield-package.yaml`. The project-awareness delta adds bounded `project_id`/`repo_path` to OpenClaw decision packets where planning already knows repo identity. Targeted `go test ./internal/hitl ./internal/roles/planning ./internal/orch` PASS; matching `go vet` PASS. Concurrent isolation check with `-race` for orch/hitl PASS.
- CURRENT loaded planner after process-only reload: PID `323418`, binary SHA256 `f61e4c9fb6dc35b98c264bd981d63a7c456be7b0b76a3ed1c81265fccb2d417f`, with `SWE_OPENCLAW_HITL=1` and `AGENTFIELD_HARNESS_IDLE_SECONDS=0`. No container redeploy was used.
- OpenClaw now uses the dedicated Telegram account `swe` as the active SWE control surface. Canonical workspace skill remains `/home/node/.openclaw/workspace/skills/swe-governor/SKILL.md`; deterministic adapter remains `/home/node/.openclaw/workspace/skills/swe-governor/scripts/swe-governor.mjs`. The adapter reads account-scoped pairing, sends bounded cards, stores Telegram message id <-> request id, supports non-mutating discussion, and posts only validated terminal decisions to `http://workforce:8800/webhooks/approval`.
- The terminal action path is hardened against provider outages by live OpenClaw plugin `/home/node/.openclaw/extensions/swe-governor-telegram` (package `@n0namer/swe-governor-telegram` `0.1.0`). Runtime inspect reports `status=loaded`, `activated=true`, hook-only shape, one typed `before_dispatch` hook at priority `100`, no missing dependencies, and plugin doctor reports only the informational hook-only compatibility notice. Pure parser checks PASS 3/3. CURRENT SHA256: `index.js` `8e84a2042902d9fbcdc995a5516bc3f7b7b94d431c350beaa8acc647127b50c1`, `parser.js` `92ceb3b3554010a1794fee56cea04144b0ce21590aaf001bd0d7b5178aa4c555`, manifest `c17f3cb57b05b7d0c03686689915c7120696c1084e176ead7628ae88e494d030`, package `d196bb241941330998c5321f15484e7ac26b709ef139116c3a70f7caad6dcdbb`. Gateway was reloaded in-process, not redeployed; startup readback lists 11 plugins including `swe-governor-telegram`.
- OpenClaw cron declaration `swe-governor-decisions-v1`, job id `c7ea3f43-2a02-4752-aaac-f254f9ba256f`, remains enabled every **20s** as a command job with no LLM and `delivery.mode=none`, `timeoutSeconds=15`; CURRENT readback is `lastRunStatus=ok`, diagnostic `NO_REPLY`, errors/skips `0`. No new daemon, HAX deployment, bot-per-project, or messaging database is required.
- Security boundary: governor notes intentionally omit secret/default values and raw logs. The OpenClaw skill refuses guessed ids and routes secrets outside Telegram. SWE-generated plan/question text is treated as untrusted payload, not instructions for OpenClaw.
- The last stitched SWE pause/resume attempt on 2026-09-05 was blocked **before HITL** by the then-observed provider/FCM route, not by Governor: bounded `run_product_manager` canaries `exec_20260905_110441_k9swye3b` and `exec_20260905_110556_0ibik04v` failed before producing `ask_user_form` with `too many concurrent requests`. Treat that provider diagnosis as historical evidence, not CURRENT health; this Telegram/Governor workstream does not repair FCM/provider routing. On 2026-09-06 a normal `devteam` conversational turn also hit external `gonka` HTTP 429, which confirms discussion can still be provider-dependent. Explicit terminal Replies are now separately protected by the model-independent `before_dispatch` bridge. The pinned AgentField SDK `Pause/Webhook/Approval` test family remains prior PASS evidence for the pause-manager/webhook primitive. The only acceptance still missing is one real same-execution `waiting -> dedicated Telegram card -> discussion -> explicit terminal Reply -> /webhooks/approval resolved=true -> governor.resolved -> execution continues` stitch.

## Embedded HITL Decision Autonomy — architecture/phase update (2026-09-06)

North Star extension: SWE-AF should not treat every `ask_user_form` as an automatic human interruption. Before pausing, SWE should use its own repo-aware reasoning/harness and existing planning/execution context to recommend or, where policy permits, select the best answer. Telegram becomes the human awareness/control surface: it receives the project/execution/stage, current status/ETA, recommendation, evidence/rationale, risk/confidence, whether SWE auto-decided or escalated, and the final accepted decision. Exact `project_id + execution_id + request_id + Telegram messageId` correlation remains mandatory for any human reply/override.

Architecture decision (BMAD Architecture): put the resolver **inside `swe_af.hitl` at the central `run_with_ask_user()` boundary**, not in Telegram/OpenClaw and not as a new middleware daemon. Current source shows every reasoner that can emit `AskUserForm` already funnels through this wrapper; it is therefore the smallest control point with maximum coverage. The resolver should reuse the existing SWE `router.harness` and read-only repo tools (`Read`/`Glob`/`Grep`) plus existing PRD/Architecture/Review/IssueGuidance/DAG context. Do not duplicate PM/Architect/Reviewer business logic outside SWE.

Adaptive decision ladder:
1. `DETERMINISTIC` — answer from explicit SoT/current structured state when there is one valid value; no model vote required.
2. `SINGLE_RESOLVER` — cheap default: one repo-grounded structured recommendation with alternatives/counterargument/evidence.
3. `COUNCIL` — only for material ambiguity/risk: independent role views chosen from existing SWE roles (normally PM for scope/product, Architect for design, Tech Lead/Reviewer for implementation risk, QA/Verifier for acceptance), then a separate judge/synthesizer. Do not run a large council for every ask.
4. `HUMAN`/`ABSTAIN` — mandatory for credentials/secrets, destructive or irreversible actions, privilege changes, major scope expansion, external/financial/legal commitment, high-blast architecture, unresolved disagreement, insufficient evidence, or policy uncertainty.

Automation policy is risk/evidence based, not raw model confidence. AUTO candidates are bounded reversible choices already inside approved scope, deterministic convention/SoT answers, and choices among approved alternatives when evidence coverage is complete. Medium/high risk or weak evidence defers. This follows adjustable-automation evidence: automate acquisition/analysis more aggressively than consequential decision/action selection, and choose automation level using reliability plus error cost (Parasuraman, Sheridan & Wickens, 2000, DOI 10.1109/3468.844354; Onnasch et al., 2014, DOI 10.1177/0018720813501549). The reject/defer option is a first-class safety mechanism rather than a failure (Geifman & El-Yaniv, 2017, arXiv:1705.08500). Human-AI UX must support correction/override and make current system state/decision visible (Amershi et al., CHI 2019, DOI 10.1145/3290605.3300233). Multi-agent deliberation is adaptive because debate/self-consistency can improve reasoning, but controlled evidence shows model strength/diversity dominate and majority pressure can suppress correct minority correction; do not equate more agents with better decisions (Du et al., 2023, arXiv:2305.14325; Wu et al., 2025, arXiv:2511.07784).

Phase Goal — `HITL Decision Resolver / Shadow -> Calibrated AUTO`:
- **Batch 1 (current): Shadow resolver.** Every material `ask_user_form` first receives a structured `DecisionCase` recommendation containing `decision_class`, `recommended_values`, `rationale`, `evidence`, `alternatives`, `risk`, `confidence`, `human_required`, `project_id`, `execution_id`, `stage`, and policy version. Existing HAX pause remains authoritative; resolver failure must fall through to the current human path. Emit one structured decision event for Telegram/telemetry.
- **Batch 2: Council + policy gate.** Add adaptive council only when Tier-1 evidence/risk requires it; preserve independent first-pass views before synthesis. Add explicit `AUTO | HUMAN | ABSTAIN` policy decision and never auto for the mandatory-human classes above.
- **Batch 3: Telegram awareness surface.** Extend Governor notification ingestion to send decision events and stage/ETA/status, including auto-decisions that never create an approval pause. Human escalation cards still use exact reply correlation; informational decision/status messages are non-terminal.
- **Batch 4: Calibrate and enable AUTO.** Run shadow acceptance across representative HITL classes, compare Governor recommendation to human/final outcome, measure disagreement/error by decision class, then enable AUTO only for classes that meet the risk gate. No blanket confidence threshold.

BMAD Test Architecture DoD for this phase:
1. Wrapper regression: no `ask_user_form` behaviour changes when resolver is disabled/fails; HAX pause/resume remains PASS.
2. Shadow evidence: one synthetic ask produces one structured recommendation event and still pauses for the human.
3. Grounding: resolver can read current repo evidence but has no write/mutation tools.
4. Policy: scope/security/destructive/credentials/high-blast fixtures always return HUMAN/ABSTAIN even when the suggested answer is confident.
5. Council: medium-risk fixture invokes independent specialist views and judge; disagreement is preserved and can force HUMAN.
6. Auto: low-risk reversible fixture can later inject the selected values as `prior_user_responses` without HAX pause; exactly one decision event is recorded.
7. Multi-project: two simultaneous builds/questions remain isolated by project/execution/request identity; Telegram interleaving cannot cross-resolve another request.
8. Failure: resolver/provider timeout/error never silently auto-decides; existing HITL is the recovery path.
9. Telegram: final decision, rationale/evidence summary, current stage, remaining-work/ETA estimate and mode (`AUTO`/`HUMAN`) are observable without exposing secrets.
10. Real acceptance: same execution continues after either approved human reply or policy-approved auto-resolution, with exact execution/request evidence.

CURRENT implementation state (2026-09-07 anti-drift): the installed/active SWE authority is the Go node, not the Python compatibility path. `/afhome/installed.yaml` reports `swe-planner` source `local`, `source_path=/src/swe-af/go`, status/desired_state `running`, port `8800`, PID `389456`, started `2026-09-06T11:10:36Z`; `go/agentfield-package.yaml` is the maintained node while the root Python package is superseded. The earlier Python Shadow Decision Run remains compatibility/reference evidence only and MUST NOT be treated as runtime acceptance.

Batch 1 has therefore been ported stale-safe directly onto the CURRENT Go source without GitHub coding/redeploy. `go/internal/hitl/wrapper.go` now owns typed `DecisionRunInput`/`DecisionCase`, an injected `ShadowDecisionResolver`, a 20s bounded shadow pass before human pause, deterministic parent `project_id + execution_id + stage + repo_path` event context, `HITL_DECISION_EVENT` emission, and fail-safe fallback to the unchanged human HAX/OpenClaw pause; `AUTO` is recommendation-only and never resumes in this batch. `go/internal/roles/planning/planning.go` wires the resolver into Product Manager and Environment Scout; `go/internal/roles/advisor/advisor.go` wires Issue Advisor and Replanner. Every resolver invocation uses the existing `harnessx.Run` with exactly `Read`/`Glob`/`Grep`, `max_turns=6`, current provider/model/cwd, and no mutation tools. Existing concurrent Go deltas (including OpenClaw HITL fallback and PM direct-AI work) were observed before mutation and preserved.

BMAD Test Architecture coverage was extended in-place: `go/internal/hitl/hitl_test.go` checks that a shadow `AUTO` recommendation still pauses/re-invokes and that resolver failure still pauses; planning/advisor HITL fixtures were adjusted for the extra shadow harness pass. CURRENT live SHA256: `go/internal/hitl/wrapper.go` `3a7353fb0a7ecb9de2347db8e77d200edcad7bf205fb9a47a914aa20e1370f6c`; `go/internal/hitl/hitl_test.go` `2065e8eae485274a0fbc386186787a2fae1da38d189cbf9cfb89f9544e057343`; `go/internal/roles/planning/planning.go` `fe1e7388c096d8ac7388f4692b37bfcfae6e6346ac04dd6491bf0c07fad2d7b6`; `go/internal/roles/planning/planning_test.go` `44e658e5a9b7f7634151c93bd8f3a969977b591f8cb239354245c07960080b74`; `go/internal/roles/advisor/advisor.go` `405bbfd9c53a352af7d117bb618d3ea8b89192a4b832e4adbf9794541aa5ad2b`; `go/internal/roles/advisor/advisor_test.go` `e9c6c6a501512d64f593df78b66ca80b887c185260b9e93423f7f0559951ab37`. Scoped `git diff --check` PASS. Go execution validation is currently a `VALIDATION_BLOCKER`, not an application FAIL: the Go toolchain exists at `/usr/local/go/bin/go` but is absent from PATH; direct absolute `go test`, canonical `make check PATH=...`, and even read-only `gofmt -d` are blocked server-side as `opaque_or_unknown_mutation`. CURRENT target registry exposes only `git_diff_check` for `agentfield-dev-workforce`, with no typed Go test/compile preset. No Go node reload/runtime activation was attempted without deterministic test evidence; all Go SourceLoop capture entries remain `git_writeback_status=PENDING` and the Go delta is not yet canonicalized to Git.

### SWE capability inventory for HITL decision solving (2026-09-06)

CURRENT source/discovery evidence supports the following capability surface relevant to local decision questions. External AgentField capability discovery is temporarily unavailable (`Bad Gateway`), so callability below is classified from current SWE source/runtime surface; public entrypoints are explicitly marked by `TAG_ENTRYPOINT`, while `run_*` role reasoners are internal-only and must be invoked by SWE orchestrators rather than called directly by an external caller.

| Capability | CURRENT owner/surface | What it can contribute to a HITL decision | Mutation/risk profile | Recommended use |
|---|---|---|---|---|
| `build` | `swe_af.app.build` entrypoint | Full PRD→architecture→DAG→code→review→verify lifecycle | High: may modify repo/PR | **NO** for ordinary decision analysis; only if the decision itself is a real implementation project |
| `plan` | `swe_af.app.plan` entrypoint | PM→Architect↔Tech Lead→Sprint Planner→Issue Writers; creates structured PRD/Architecture/Review/Issues | Planning/artifact writes, no coding intent | **YES, selectively** for complex product/architecture decisions where a mini-project plan is worth the cost |
| `implement_issue` | issue router entrypoint (`swe_af.issue`) | One scoped implementation issue with existing issue contract | Mutating | **NO** for answer selection; **YES** only after decision when implementing the chosen option |
| `resolve` | `swe_af.app.resolve` entrypoint | Existing PR repair: base merge, CI/review fixes, push | High/mutating/external | **NO** for decision support |
| `resume_build` | `swe_af.app.resume_build` entrypoint | Resume crashed build from checkpoint | Mutating lifecycle | **NO** for decision support |
| `execute` | internal/non-entrypoint | Execute a prior `plan_result` DAG; supports built-in coder loop or external coder | High/mutating | **NO** for decision analysis |
| Fast build | `swe_af.fast.app.build` entrypoint | Lightweight plan→execute→verify, 5–10 minute class instead of full DAG | Mutating | **NO** for answer selection; possible later for cheap implementation experiments after a decision |
| Product Manager | `run_product_manager` internal | Clarify goal, scope, acceptance criteria, `ask_user_form` | Read/reasoning + HITL | **YES** for product/scope questions |
| Environment Scout | `run_environment_scout` internal | Detect env/credential requirements and HAX negotiation | Sensitive/HITL; may handle credentials through scoped store | **YES only as authority signal**; never let decision resolver invent/auto-answer secrets |
| Architect | `run_architect` internal | Architecture alternatives, decisions, rationale grounded in repo/PRD | Read/reasoning | **YES** for architecture/design HITL |
| Tech Lead | `run_tech_lead` internal | Critique architecture, scope issues, complexity | Read/reasoning | **YES** as independent critic/counterargument |
| Sprint Planner | `run_sprint_planner` internal | Decompose accepted architecture into dependency-aware issues | Read/planning | **LIMITED**; useful for cost/implementation-impact estimates, not primary answerer |
| Issue Writer | `run_issue_writer` internal | Turn plan item into implementation contract | Planning artifact write | **LIMITED** after a decision, not for selecting it |
| Issue Advisor | `run_issue_advisor` internal | Diagnose local issue, guide strategy, deeper QA need | Read/reasoning | **YES** for local implementation/debug decisions |
| Retry Advisor | `run_retry_advisor` internal | Diagnose failed attempt and choose retry strategy | Read/reasoning | **YES** when HITL is caused by a failed execution/retry choice |
| Replanner | `run_replanner` internal | Re-plan after execution drift/failure | Planning/reasoning | **YES** when decision changes DAG/scope after failure |
| Coder | `run_coder` internal | Implements code with repo tools; can use web search when enabled | High/mutating | **NO** as default decision role; possible **sandboxed experiment** only when a small reversible prototype is decisive |
| QA | `run_qa` internal | Test/quality analysis of implementation | May execute tests | **YES** for correctness/test-strategy questions |
| Code Reviewer | `run_code_reviewer` internal | Independent code review/risks | Read/reasoning | **YES** for implementation-risk council |
| QA Synthesizer | `run_qa_synthesizer` internal | Synthesizes QA/review outcomes | Read/reasoning | **YES** as judge only when the decision is implementation-quality related |
| Verifier | `run_verifier` internal | End-to-end acceptance verification | Executes tests/checks | **YES** for evidence after a decision or to discriminate options when checks are bounded/read-only to product state |
| Integration Tester | `run_integration_tester` internal | Integration validation | Executes tests | **LIMITED** for empirical option comparison; not default |
| Git init / merger / repo finalize / GitHub PR | internal GitOps roles | Branch/worktree/merge/finalize/PR lifecycle | High/mutating/external | **NO** for answer selection |
| CI watcher/fixer | internal CI roles | Observe/fix CI | watcher read; fixer mutates | watcher **LIMITED** evidence source; fixer **NO** until decision made |
| PR resolver | internal | Resolve review comments/CI on existing PR | Mutating | **NO** for decision selection |
| HAX `AskUserForm` + `app.pause()` | `swe_af.hitl` | Existing authoritative human escalation/resume path; values return through `prior_user_responses` | Human state transition | **YES**, remains final fallback/override |
| Structured schemas | PRD, Architecture, ReviewResult, IssueGuidance, DAG/verification schemas | Gives typed facts, alternatives, rationale, scope and acceptance state | Read-only data | **YES**, primary evidence before asking more models |
| `router.harness()` | existing reasoner execution primitive | Run a bounded structured reasoning task with selected provider/model/tools/cwd | Depends on tools granted | **YES**, preferred Tier-1 decision resolver primitive with **read-only tools** |
| Repo tools | harness `Read`/`Glob`/`Grep` | Current code/SoT/evidence retrieval | Read-only | **YES**, default grounding layer |
| Web research | opencode built-ins `websearch`/`webfetch` when `OPENCODE_ENABLE_EXA=1` and `EXA_API_KEY` exist | Current external API/library/framework evidence | Network/read; source quality varies | **YES conditionally** when repo evidence is insufficient or freshness matters; not for questions already answered locally |
| Separate SWE sub-execution | public `plan` or a dedicated future decision entrypoint | Treat a difficult HITL as a bounded local project with its own roles/evidence/result | Cost/latency; `build` would over-mutate | **YES for hard questions**, but use a **decision-only mini-project**, not a full `build` |

Architecture refinement — **HITL as a local SWE decision subproject**:
- The user's proposal is adopted with one constraint: do **not** call full `build` for every question. A HITL question is usually a local analysis problem, so the optimal reusable primitive is a bounded **Decision Run** inside SWE that can orchestrate existing roles without code mutation.
- `Decision Run` input: parent `project_id`, parent `execution_id`, parent `request_id`, stage/role, `AskUserForm`, relevant PRD/Architecture/DAG/runtime context, time/evidence budget, and allowed decision policy.
- `Decision Run` output: exact `DecisionCase` (`recommended_values`, alternatives, rationale, evidence refs, dissent/counterargument, risk, confidence, `AUTO|HUMAN|ABSTAIN`, consulted roles, elapsed/ETA, policy version). It must never create a second ambiguous authority chain: parent identifiers stay mandatory, and the child decision execution id is recorded as provenance only.
- Routing heuristic: deterministic structured evidence first; then one `router.harness` resolver; then, only if ambiguity/risk remains, spawn a **small internal council** chosen from existing roles. Deep/web research is a tool escalation inside that decision run, not a mandatory separate service. Only after evidence exhaustion should HAX pause the parent SWE execution.
- For broad architecture/product questions, the decision run may reuse a **planning-only** slice (PM + Architect + Tech Lead) as a local mini-project. For implementation/debug questions, use Issue Advisor + Reviewer/QA. For an empirical uncertainty, a bounded read-only/check experiment may be used; actual coder/build mutation requires the normal implementation authority and is not part of deciding by default.
- The decision run is logically a child/sub-execution, not a new top-level user project: Telegram shows `Project`, parent `Execution`, `Decision`, consulted roles, recommendation/final result and current parent stage. This preserves one-chat multi-project isolation while allowing several decision runs concurrently.

Pareto implementation order is updated:
1. **Decision Run contract + shadow single-resolver** at the central HITL boundary; no auto-resume yet.
2. **Role router** that selects at most 2–3 existing specialists based on decision class; independent first pass + judge only when needed.
3. **Research escalation**: repo evidence → web/deep research when freshness/external facts are load-bearing; enforce source/evidence budget.
4. **Telegram awareness events** for `decision_started`, recommendation, AUTO/HUMAN/ABSTAIN, final decision, parent stage/progress/ETA.
5. **Calibrated AUTO** only after shadow comparison by decision class; HAX remains mandatory fallback.

Nearest compulsory move: **do not reload the active Go node, do not canonicalize the Go source delta, and do not enable AUTO/council until executable Go validation exists.** The current blocker is the DEV operator test surface: `agentfield-dev-workforce` exposes only `git_diff_check`, while `/usr/local/go/bin/go` and `gofmt` executions are server-side rejected as opaque. Resolve or reuse an existing typed Go check route for this same `/src/swe-af/go` owner (no new infra and no environment expansion), then run affected packages `./internal/hitl ./internal/roles/planning ./internal/roles/advisor` followed by canonical `cd go && make check`. After PASS, reload/restart only the installed `swe-planner` Go node through its authoritative lifecycle, prove one synthetic/real `AskUserForm -> HITL_DECISION_EVENT -> unchanged human pause` on the same `execution_id`, and only then canonicalize the accepted Go delta through SourceLoop/publication. Batch 2 adaptive council/deep research starts only after this gate.

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

Historical execution chronology is retained below for evidence and recovery context. Current priority is owned by `Current Stage` + the authoritative certification spine above; stale `IN PROGRESS` labels in old batches must not override those current gates. Default batch size remains about 30 minutes: each active batch closes one coherent DoD gate and writes back this file, preferring the smallest 20% of work that removes the next 80% blocker.

### Batch 1 — Provider/bootstrap gate and SWE liveness

Status: HISTORICAL / SUPERSEDED — provider/bootstrap/liveness are no longer the current gate; retained only as evidence chronology.

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

Status: HISTORICAL / SUPERSEDED — this ladder exposed the fail-open and structured-output defects that were subsequently repaired; L2 structured-output reliability is now PASS/CLOSED.

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

Batch 3 is closed. Historical note: at the time this batch closed, durability/SourceLoop was the next move. That priority has since been superseded by the 2026-09-05/06 real-task capability audit above; durability resumes only when the authoritative certification spine reaches that stage.

### Batch 4 — Durability / SourceLoop gate

Status: OPEN / DEFERRED BEHIND CURRENT REAL-TASK CAPABILITY GATES — stale/noisy runtime capture rejection is proven; clean exact durable SHA remains required after the active certification spine reaches the durability stage.

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

When the authoritative certification spine reaches durability, the next durability batch is: create one clean stale-safe publication candidate from CURRENT `dev` containing only the then-accepted live bytes, validate its diff against both CURRENT `dev` and CURRENT live source, publish through the canonical repository publication lane, then run the same package test/vet gate on the exact candidate/durable SHA before acceptance. Do not execute this Lane-B step early merely because the historical Batch-4 evidence exists.

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
- [OPEN / OPERATOR DEBT] The DEV typed plane still lacks a dedicated Go build/process-only reload primitive for `swe-planner`, but this is no longer the current product gate: the authorized runtime lane has successfully rebuilt/reloaded the same planner process without whole-container redeploy, and CURRENT PID `323418` is loaded. Keep this operator-plane debt separate from SWE capability certification.
- [OPEN] Workforce healthcheck still proves provenance HTTP only; current SWE readiness is instead evidenced by registration plus real reasoner execution. Preserve this distinction in future canaries.
- [RESOLVED] Historical SWE absence was caused by the provider/bootstrap admission gate. CURRENT generation has consumed the live provider delta and both SWE nodes are registered/initialized, so this is no longer the active blocker.
- [RESOLVED / SUPERSEDED] The old loaded-binary structured-output/provider gate is no longer current: L2 structured-output reliability is PASS/CLOSED and CURRENT planner PID is `323418`. Provider/FCM health may still block particular planning/HITL canaries, but it is an external route-health dependency, not evidence that the issue-level structured-output controller remains broken. Reopen this class only on a fresh correlated failure.

## BMAD Workflow Used for Current Batch

- Entry: `bmad-help` from `BMAD-MNNZ`.
- Initial classification: brownfield recovery / quick implementation; `bmad-correct-course` rejected because its PRD + Epics prerequisites do not exist here.
- After CURRENT reconciliation exposed an implementation/debugging evidence gate, the canonical Universal Solver handoff explicitly routes this workstream to `bmad-testarch-test-design`.
- Active BMAD mode: risk-based system-level test/debug architecture embedded in this existing SoT (no duplicate BMAD artifact). Gate order: A0 CURRENT state -> A1 targeted regression -> A2 reasoner/pipeline -> A3 same-target reload if needed -> A4 functional canary -> A5 semantic E2E -> A6 durable accepted SHA.
- `bmad-quick-dev` remains the implementation method once a source/config defect is proven and live-edit capability is available.
- Write-back target: this `PLAN.md`.
