# SWE-AF Project Plan (Source of Truth)

Status: active
Last reconciled: 2026-09-13
Canonical GitHub owner: `n0namer/swe-af`
Canonical branch for project SoT: `dev`
Canonical project SoT: this file (`PLAN.md`)
Product-code inner loop: CURRENT `/src/swe-af` runtime/exact-source workspace; GitHub is publication/canonicalization only.
Historical evidence: retained in Git history; old checkpoints do not override the CURRENT section below.

## Global North Star

Bring SWE / SWE-AF into a genuinely working end-to-end state and prove it on real software-engineering tasks with independent executable acceptance evidence, at economically acceptable cost.

Operationally, SWE must:
1. receive a real engineering task;
2. localize the relevant source;
3. edit the actual repository/workspace;
4. run the appropriate canonical validation;
5. produce the required diff/commit/artifact;
6. pass an independent executable oracle/verifier;
7. preserve already-completed work across partial failure/timeout;
8. recover/continue instead of blindly restarting;
9. measure correctness, cost and latency together;
10. reproduce the result across a small controlled task ladder and then on a clean materialized runtime.

AgentField, FCM, OpenCode, Coding Station, SourceLoop and contract completion are supporting mechanisms, not North Stars.

## Authority and Anti-Drift

- Latest explicit project objective + this CURRENT SoT define the project route; CURRENT runtime/readback defines actual state.
- Do not use GitHub-edit -> redeploy as the product-code debugging loop.
- Runtime code edits happen in the authorized container/exact-source workspace first; test there; publish only the exact accepted delta.
- Do not infer acceptance from model prose, reviewer PASS, `success=true`, a commit existing, HTTP/container health, or tests that do not exercise the required behavior.
- On timeout/ambiguous result, inspect effect/post-state before any retry.
- One 30-minute batch = one coherent DoD gate. After every material result, re-observe and replan from fresh CURRENT.
- Before each batch write:
  - GLOBAL NORTH STAR
  - CURRENT BLOCKER
  - THIS BATCH
  - NORTH-STAR DELTA
  - STOP CONDITION
- STOP_SIDEQUEST if work does not materially shorten the path to independently accepted SWE tasks.

## CURRENT — 2026-09-13

### Source/runtime identity

- CURRENT live product source: `/src/swe-af`.
- Git baseline HEAD: `58c4e0d19081bc52363c120b7963a34cebb1e894`.
- Working tree is intentionally ahead of baseline: fresh readback shows 60 tracked modified/staged paths + 2 untracked paths.
- Treat product identity as `58c4e0d... + exact working-tree delta + loaded process generation` until accepted deltas are canonicalized.
- Full CURRENT Go suite PASS:
  - `/usr/local/go/bin/go test ./...`
  - Go runtime: `go1.25.14`.
- Do not mass-commit or reset the dirty runtime tree.

### SWE execution lane

- Root cause of the fresh pre-broker failure is proven: clean target base `2c374989...` lacked provider `fcm` in its project `opencode.json`, so OpenCode accepted `-m fcm/fcm` but failed before a broker request.
- Executable A/B evidence:
  - clean base direct OpenCode canary -> `Unexpected server error`, FCM `requestsRouted` unchanged `697 -> 697`;
  - CURRENT `/src/swe-af` direct canary -> `DIRECT_CANARY_OK`, FCM `697 -> 699`;
  - clean base + runtime-owned FCM overlay -> `OVERLAY_CANARY_OK`, FCM `699 -> 701`.
- Runtime-owned OpenCode overlay in `/src/swe-af/go/internal/harnessx/run.go` now defines provider `fcm` while preserving no-install permissions. Deterministic RED->GREEN contract test added in `harnessx`; targeted packages and full CURRENT `go test ./...` PASS.
- Current exact tested planner binary: `/tmp/swe-planner-f06-sdk-20260913`, SHA256 `581862f59722400ee0999176339dfb99f8dc67cd7901bf39d9861caa5df4d4`; current planner PID `548849`, built from frozen SWE baseline `cdc39105e92937e0085410a656346471b2dc6834` with a temporary Go module replace to the tested local AgentField SDK fix `b160245833ec51f5296905c32e49260a62c76e26`. Control-plane node is active/ready at `http://172.16.22.7:8005` with non-empty process `instance_id`; this is process-only current-runtime evidence, not a durable AgentField deployment.
- The binary includes the proven runtime-owned FCM overlay plus a full-Build planning parity fix: when runtime/provider resolves to OpenCode, Product Manager now uses the harness/OpenCode path instead of AgentField direct AI. Deterministic RED `TestProductManagerOpenCodeUsesHarnessWhenDirectAIIsAvailable` reproduced the old direct-AI misroute for `fcm/fcm`; targeted planning tests and full `go test ./... -count=1` PASS after the fix.
- Python reference `run_product_manager` always uses the configured harness provider/model; the Go direct-AI shortcut was therefore a port/runtime-contract regression, not an FCM outage.
- Post-fix issue-level task `exec_20260912_094308_za4vz5sv` reached FCM repeatedly, edited the target repo, and completed with a bounded two-file commit. The issue execution route is recovered; full-Build acceptance remains separate and is currently 0/3.

### FCM

- FCM is not globally down.
- Fresh `/health`: version `0.5.81`, 34 active models, 8 effective providers, no in-flight request at readback.
- Small current direct telemetry for `routerai/z-ai/glm-5.3-flash`: 6/6 successful calls.
- Therefore the L3-26 failure is narrower than “FCM is broken”: current evidence points to the OpenCode/FCM tool trajectory, structured-output/provider boundary, or another execution-route integration seam.
- FCM decides model/provider routing; it does not own authoritative engineering obligation state.

### AgentField compatibility / structured-output contract audit

- SWE-AF is pinned to AgentField Go SDK `v0.0.0-20260723130821-20955b2637b4` (commit `20955b2637b4`, 2026-07-23). Current upstream is `v0.1.139-rc.1`, commit `4aa3fe688dfa1f2437ac49f6cbe72aed43ddca07` (2026-09-10).
- Exact pinned SDK `go test ./... -count=1` PASS across `agent`, `ai`, `client`, `did`, `harness`, `inputs`, `types`. Fresh upstream `sdk/go` `go test ./... -count=1` also PASS.
- Relevant upstream drift does **not** contain an obvious fix for the observed PM schema failure: `schema.go`, `opencode_test.go`, and `parity_test.go` are byte-identical between pin and current main. `runner.go` changed provider resolution/output-dir isolation/metrics; `opencode.go` changed token accounting. The core file-write schema retry algorithm remains materially the same.
- AgentField unit coverage includes `missing output file -> retry -> success`, but the retry provider is a mock that writes the file. There is no real OpenCode integration test that drives a non-trivial schema through the actual CLI/model/file protocol; OpenCode tests use fake scripts/mocked CLI seams. Therefore green AgentField tests do not cover our production failure mode.
- AgentField-only live A/B on the exact pinned SDK + OpenCode `1.17.15` + `fcm/fcm`:
  - simple schema `{status:string}` in single mode: PASS, parsed result `status=ok`, 3 turns;
  - exact SWE `schemas.PRD` in single mode: FAIL trajectory; model writes the JSON Schema object itself (`$schema/$defs/properties/...`) instead of a PRD instance; AgentField correctly diagnoses expected-vs-actual top-level keys and starts schema recovery, but the model repeats the schema object;
  - exact `schemas.PRD` in incremental mode: PASS, `parsed=true`; model builds a real PRD instance field-by-field and completes after 22 reported turns.
- Exact reflected PRD schema is only `1403` compact bytes (~350 estimated tokens), so AgentField `SchemaMode=auto` would **not** switch to incremental (`auto` threshold = 4000 estimated tokens). Explicit per-role policy is required for this schema/model pair.
- CURRENT FCM runtime telemetry has one model with signal: `routerai/z-ai/glm-5.3-flash`, 6/6 successful backend calls. This is the strongest current model evidence but is not yet per-execution correlated.
- Separate parity debt: commit `c8ff657` (2026-08-24) reduced Go planning/issue-writer defaults from architecture/Python `150` turns to `2`. However AgentField OpenCode provider ignores `Options.MaxTurns` entirely in both pinned and current upstream, so this does **not** explain the current OpenCode PM failure. It remains a cross-provider parity/config debt, not this batch's fix.
- Decision for PM/structured output: **do not upgrade AgentField as a speculative fix**. That failure was localized to SWE's structured-output policy interacting with the routed model and was repaired independently.

### AgentField F06 restart/state convergence audit

Controlled fault injection reproduced a distinct platform defect outside SWE task code:
- exact execution: `exec_20260912_214730_j8rp9ycc` / run `run_20260912_214730_vtmgqo7a`, reasoner `swe-planner.run_git_init`;
- execution was observed `running`, planner PID `529917` then received SIGTERM, and the same tested binary/env restarted as PID `535405`;
- control-plane still reported the execution `running` at +0.5s, +2s, +5s, +20s and at later readback; old PID became a zombie, new planner was healthy/active, no OpenCode mutator remained, and the sacrificial main worktree stayed clean;
- control-plane logs showed new planner registration but no terminalization/reap event for that execution. A second writer was intentionally **not** started after the first divergence.

Deployed platform identity:
- Coolify app `universal-solver-agentfield-exact-dev-git-20260825`, repository `n0namer/universal-solver`, branch `ops/pr-af-fcm-shared-20260902`, deployment commit `ffa6c6a56814b59da19903bd56c04f8cafdb44ae`;
- compose pins AgentField control-plane/runtime source to exact SHA `4d337c1ae5104418311fcba414a1c2f85c2abb89` (2026-08-24).

Root cause contract:
- control-plane `AgentNode.InstanceID` explicitly defines a per-OS-process identity used to orphan/reap in-flight work after restart;
- the Go AgentField SDK `NodeRegistrationRequest` has **no `instance_id` field**, in both deployed pin and current upstream; search of current `sdk/go` finds no `instance_id`/`InstanceID` support;
- Python AgentField SDK is the reference implementation: it generates a fresh UUID4-hex `agent_instance_id` per `Agent` process and sends it in registration and heartbeats, with dedicated regression tests;
- control-plane restart detection in both deployed and current code requires non-empty old/new instance IDs and a change between them. Therefore the Go planner is treated as a legacy opt-out and immediate restart reap cannot fire;
- graceful Go SDK shutdown calls `/api/v1/nodes/:node_id/shutdown`, but the deployed handler updates node presence/status only; it does not terminalize accepted executions;
- Go SDK then waits up to 5s in `http.Server.Shutdown`; long-running handler/subprocess lifetime is not independently proven contained after process exit;
- stale cleanup is only a backstop: deployed defaults are `stale_execution_timeout=30m`, `cleanup_interval=1h`, insufficient as an autonomous no-overlap guarantee.

Upstream chronology:
- deployed `4d337...` already contains the basic restart reason and whole-agent reap path, so the defect is **not** absence of the status constant;
- `2d7fc7264a7aaf5dd36fa4ed05ad55f8a297e117` (2026-08-31) adds instance-scoped reap/read identity and restart safety improvements;
- `2638b9e92a1f9ee6c6f9a3cd223da231bf1ece86` (2026-08-31) closes additional dispatch/shutdown persistence holes;
- however current upstream Go SDK still lacks process instance identity, so a control-plane-only upgrade is **not sufficient** for this Go planner.

Implementation / verification status:
1. SWE-AF was **not** patched for F06 and stale timeouts were not reduced. Exact deployed AgentField source `4d337c1...` was cloned container-locally to `/tmp/agentfield-f06-fix`; GitHub was not used as a programming loop and the control-plane was not redeployed.
2. Deterministic RED `TestAgentInstanceIDPropagatesAndChangesPerProcess` proved the Go wire payload omitted `instance_id`. Minimal GREEN added one UUIDv4-compatible 32-char process ID per `Agent`, stable across registration/reconnect/lease heartbeat and unique across new Agent processes. Review also added shutdown cancellation of every tracked in-flight reasoner context so context-aware children cannot survive planner shutdown as mutators.
3. Exact local AgentField commit `b160245833ec51f5296905c32e49260a62c76e26` (`fix(go-sdk): identify agent process instances`) changes only `sdk/go/agent/agent.go`, `agent_lifecycle.go`, `agent_lifecycle_test.go`, `cancel.go`, and `sdk/go/types/types.go`. Targeted tests, targeted `-race`, full `sdk/go` suite and `git diff --check` PASS. Existing control-plane restart/orphan suite, including `TestRegisterNodeHandler_ReapsOrphansOnInstanceChange`, also PASS on the exact source.
4. Frozen SWE baseline `cdc3910...` was validated against that SDK through a temporary local Go module replace: full SWE `go test ./... -count=1` PASS; planner `/tmp/swe-planner-f06-sdk-20260913` was built and loaded process-only without redeploy.
5. The original F06 fault was then re-injected on a real full `build` execution `exec_20260912_223956_dfmcy6va` / run `run_20260912_223956_a87y880t`. It was observed `running`, planner PID `548717` received SIGTERM, and the patched planner restarted as PID `548849`. The execution became terminal `failed: context canceled` after ~101 ms and remained terminal at +0.5s/+2s/+5s; old PID was zombie, new PID running, no OpenCode mutator targeted the workspace, and the sacrificial repo stayed clean. The current node record carries a non-empty process `instance_id`.
6. F06 is therefore **VERIFIED/CLOSED for the current process-only runtime**: restart no longer leaves accepted work stale `running`, and shutdown cancels the active reasoner before a second writer can overlap it. Durable AgentField source publication / normal DEV deployment is still pending and must not be confused with this runtime proof.

Deferred neighboring F07/identity risks exposed by review, not folded into this gate: reject/ignore stale heartbeat/status updates from an old `instance_id`; prevent duplicate execution-ID cancel registration from overwriting an older cancel func; close the narrow shutdown race where a new request could arrive between cancel-all and listener shutdown.

### Test coverage / risk-based gap audit

Fresh coverage was measured on CURRENT exact source with `go test ./... -count=1 -covermode=atomic -coverprofile=...` and `go tool cover -func`.

Coverage snapshot:
- SWE-AF overall Go statement coverage: **74.1%** (pre-batch baseline 73.9%).
- critical packages: `internal/coding` 54.1%, `internal/dag` 62.7%, `internal/orch` **72.6%** (70.7% before this batch), `internal/issue` 78.9%, `internal/roles/planning` 79.0%, `internal/harnessx` 86.4%, `internal/node` 92.2%.
- critical functions: `orch.Build` **69.8%** (61.4% before this batch), `orch.Plan` 87.3%, `dag.RunDAG` 77.6%, `coding.RunCodingLoop` 80.9%, `roles/planning.RunProductManager` 74.2%, `harnessx.executeStructured` 100%.
- exact pinned AgentField SDK overall suites PASS; AgentField `harness` coverage is **92.1%**. Relevant functions: `Runner.Run` 92.7%, `handleSchemaWithRetry` 92.3%, `OpenCodeProvider.Execute` 87.1%, `DiagnoseOutputFailure` 95.5%.

High-risk contracts added/strengthened in this batch:
1. `TestBuildVerifierFailureGeneratesFixAndReverifies`: proves full Build self-healing path `verifier RED -> generate_fix_issues -> execute fixes -> verifier GREEN`. This moved `Build` coverage from 61.4% to 69.8% and `failedCriteriaOf` to 100%.
2. `TestPlanForwardsExplicitOpenCodeRouteToProductManager`: protects top-level `Plan` propagation of explicit `ai_provider=open_code` + `model=fcm/fcm`, the boundary previously implicated in PM misrouting.
3. `TestAgentFieldHarnessArtifactsNeverLandOnBranch`: reproduces the MICRO-2 delivery failure where `.agentfield-out-*/.agentfield_output.json` / `.agentfield_schema.json` became product delivery. Deterministic RED reproduced `GIT_DELIVERY`; minimal runtime-junk policy now scrubs/excludes/ignores only reserved AgentField harness artifacts. Existing `TestScopedIssueRejectsUnexpectedDeliveryFiles` remains GREEN, proving ordinary unexpected files still fail closed.
4. Mutation adequacy spot-check: in a disposable copy, the verifier-fix exit condition was mutated to skip the fix cycle. `TestBuildVerifierFailureGeneratesFixAndReverifies` failed exactly as intended. The test therefore distinguishes the working and broken self-healing behavior; it is not coverage-only theater.

Remaining material gaps (do not confuse with missing line coverage):
- AgentField has no deterministic **real OpenCode** integration test for complex structured-output file protocol; existing retry tests use mocks/fake CLI. Current live A/B canaries therefore remain required for this boundary.
- control-plane stale `running` state / orphan child after agent restart is an integration/fault-injection concern outside SWE unit coverage; it requires execution-state reconciliation tests against the real control plane, not more Go unit mocks.
- `dag.runExecuteFn` 25.8% and legacy/multi-repo worktree functions (`setupWorktrees` 11.5%, `runIntegrationTests` 45.3%, cleanup paths ~23-39%) remain lower-covered. Existing advisor/replan/resume contracts are already directly tested; raise these only when the corresponding external-execute or multi-repo path enters the active acceptance ladder.
- `cmd/*` 0% is startup plumbing, not a current P0; node registration functions are already ~92-100% covered.

Decision: do not chase an arbitrary global coverage target. Use statement coverage to locate weak areas, but require behavior/oracle evidence on North-Star paths. Academic mutation-testing evidence supports focusing on the oracle gap and changed critical code rather than whole-repo mutation volume; future mutation checks should remain incremental and risk-targeted.

### System-level fault-model test design (BMAD)

Mode: **System-Level** (`bmad-help` -> `bmad-testarch-test-design`). This section is the canonical test-design output; no separate BMAD test-design documents are created because `PLAN.md` is the existing project SoT.

Testability assessment:
- strong: deterministic Go core, injectable `CallFn` seams, explicit checkpoints, isolated sacrificial repo, canonical Go validator, executable verifier/fix/replan state machines;
- actionable gaps: no first-class ambiguous-effect state, no deterministic real control-plane restart/orphan harness, no real OpenCode complex-schema integration test, incomplete causal telemetry tying child execution/model/effect to parent;
- reliability ASR: **UNKNOWN/ambiguous effect forbids autonomous mutation until reconciled**;
- reliability ASR: **completed non-idempotent effect executes at most once across timeout/restart/resume**;
- observability ASR: every async child must be correlatable to parent execution + workspace + tested/delivered identity.

Risk register and fault families (P=probability 1-3, I=impact 1-3):

| ID | Fault family | Cat | P | I | Score | Priority | Owner layer | Required evidence |
|---|---|---:|---:|---:|---:|---|---|---|
| F01 | config / routing / precedence | TECH | 2 | 3 | 6 | P0 | SWE/Platform | deterministic contract + live route canary |
| F02 | version / provenance / tested!=delivered | TECH | 2 | 3 | 6 | P0 | SWE/SourceLoop | exact SHA/config manifest + differential canary |
| F03 | omission (missing output/file/field/checkpoint/ack) | TECH | 2 | 3 | 6 | P0 | SWE/AgentField | unit + integration omission injection |
| F04 | invalid value / schema / corrupted structured output | TECH | 3 | 2 | 6 | P0 | SWE/AgentField | schema RED/recovery + real OpenCode canary |
| F05 | crash / dependency unavailable / 5xx / disconnect | OPS | 2 | 3 | 6 | P0 | AgentField/Platform | process/network fault injection |
| F06 | timing / stale async state / delayed completion | OPS | 3 | 3 | 9 | P0 | AgentField/Platform | control-plane convergence test |
| F07 | concurrency / ordering / parallel writers | TECH | 2 | 3 | 6 | P0 | SWE/AgentField | deterministic race/barrier + isolated-workspace test |
| F08 | persistence / idempotency / ambiguous effect after timeout | DATA | 2 | 3 | 6 | P0 | SWE + execution substrate | timeout/effect fault injection + resume oracle |
| F09 | recovery-policy failure (retry/advisor/replan/fix skipped or loops) | TECH | 2 | 3 | 6 | P0 | SWE | state-machine contract + mutation adequacy |
| F10 | scope / delivery contamination / dirty worktree | DATA | 2 | 2 | 4 | P1 | SWE | fail-closed delivery tests |
| F11 | oracle / validation false PASS | TECH | 3 | 3 | 9 | P0 | SWE QA | independent oracle + mutation test |
| F12 | arbitrary model/tool behavior (ignores contract, wrong tool/action) | TECH | 3 | 2 | 6 | P0 | SWE/AgentField | adversarial structured-output/tool-call matrix |
| F13 | resource / budget / rate-limit / runaway turns | PERF | 2 | 2 | 4 | P1 | FCM/SWE | bounded-budget + 429/timeout tests, cost telemetry |
| F14 | observability / evidence loss / missing causal IDs or cost | OPS | 2 | 2 | 4 | P1 | AgentField/FCM/SWE | provenance completeness assertions |
| F15 | permission / secret / external-mutation violation | SEC | 1 | 3 | 3 | P1 | SWE/AgentField | permission-deny + secret-redaction tests |
| F16 | partial write / corrupt checkpoint / artifact-state mismatch | DATA | 2 | 3 | 6 | P0 | SWE | corruption injection + fail-closed resume |

Coverage design (priority != execution timing):
- **Unit/component:** F01, F02, F03, F04, F08, F09, F10, F11, F15, F16 where invariants are locally decidable.
- **Real integration:** F03/F04 at OpenCode file protocol; F05/F06/F07/F08/F12/F14 at AgentField/control-plane/workspace boundaries.
- **Fault injection / chaos:** restart, delayed/duplicated completion, timeout-after-effect, provider 429/5xx, corrupt checkpoint, parallel writer collision.
- **Metamorphic/differential:** equivalent config forms -> same route; clean run vs resume -> same final state with no duplicate effect; pinned vs candidate SDK -> same harness semantics.
- **Mutation adequacy:** only critical gates/recovery/oracles; a test is accepted only if a plausible mutation makes it RED.

Current family status:
- covered/strong: F01 routing, F02 exact identity discipline, **F08 ambiguous-effect fail-closed across built-in coder, external ExecuteFn, standalone implement_issue, DAG checkpoint and top-level Build**, F09 verifier-fix + advisor/replan/resume, F10 delivery contamination, F11 one proven mutation oracle;
- partial: F03/F04 (unit + live canaries, but no deterministic real-OpenCode CI integration), F07 (concurrency limits but not workspace collision fault injection), F12, F14, F15, F16;
- covered/strong on CURRENT runtime: **F06 restart/state convergence** — deterministic SDK RED->GREEN + race/full suites + existing control-plane orphan tests + live same-fault re-injection prove terminal convergence and zero overlapping mutator after planner restart. Durable AgentField publication remains pending.
- active gaps: **F05 dependency/crash handling** and partial F07 concurrency/identity neighbors (stale old-instance heartbeat/status, duplicate execution-ID cancel registration, narrow shutdown accept race).

F08 executable evidence:
- pre-fix RED `TestCoderTimeoutFailsClosedAsAmbiguousEffect`: mutation-capable coder ignored cancellation, remained in-flight after local timeout, while the coding loop returned ordinary failure;
- pre-fix RED `TestCoderTimeoutAbortsDAGWithInFlightCheckpoint`: DAG continued recovery and even accepted the issue with debt while the coder call was still live;
- GREEN `TestCoderTimeoutFailsClosedAsAmbiguousEffect`: exactly one coder call, typed/stable `AMBIGUOUS_EFFECT`, no automatic coder retry;
- GREEN `TestCoderTimeoutAbortsDAGWithInFlightCheckpoint`: no advisor/replanner, RunDAG returns error, checkpoint preserves `in_flight_issues=[a]` for reconciliation;
- GREEN `TestExternalExecuteFnAmbiguousEffectSkipsRetryAdvisor`: cross-process marker from remote ExecuteFn also fails closed; no retry-advisor/issue-advisor/replanner;
- GREEN `TestAmbiguousCoderTimeoutPreservesWorktreeAndSkipsDelivery`: standalone `implement_issue` propagates UNKNOWN, preserves worktree, skips verifier/PR/delivery cleanup;
- GREEN `TestBuildStopsImmediatelyOnAmbiguousExecuteEffect`: full Build stops before verifier/finalize when execute returns UNKNOWN;
- regression guards `TestAdvisorTimeoutFailsNotHang` and `TestCoderExceptionFailsUnrecoverable` remain GREEN, proving read-only/ordinary failures did not become ambiguous-effect aborts;
- targeted `coding/dag/issue/orch` fault tests PASS, full `/usr/local/go/bin/go test ./... -count=1` PASS, `git diff --check` PASS.

NFR evidence plan:
- reliability: P0 fault families 100% pass; duplicate non-idempotent effects = 0; UNKNOWN forbids further mutation; evidence = Go tests + fault-injection runlogs/checkpoints;
- maintainability: canonical Go suite + targeted mutation checks; evidence = commands/coverage/mutation RED;
- security: no secret disclosure and permission-denied external mutation; thresholds beyond existing permission contract remain UNKNOWN until explicitly specified;
- performance/cost: no correctness gate based on guessed latency/cost; record wall time/model/cost where available and treat missing telemetry as evidence gap.

Execution strategy:
- PR/inner loop: deterministic unit/component/contract tests (<15 min), full Go suite, targeted mutation checks;
- nightly/controlled: real OpenCode structured-output matrix and bounded provider failure tests;
- weekly/pre-release: control-plane restart/orphan/stale-state chaos and full-Build recovery ladder.

Quality gates:
- P0 fault-family invariants: **100% PASS**;
- P1: >=95% PASS, no unresolved high-risk regression on active production path;
- no open score >=6 fault without an explicit fail-closed mitigation/evidence plan;
- overall line coverage is secondary; risk/fault-family coverage and oracle adequacy are release evidence;
- full NFR PASS/CONCERNS/FAIL remains deferred until executable evidence exists.

Entry criteria for resilience testing: exact tested source/runtime identity, isolated sacrificial workspace, zero pre-existing mutating child on target, planner/control-plane reachable. Exit criteria: all P0 families have a discriminating contract at every applicable critical boundary, plus at least one real fault-injection proof for each external async boundary.

Immediate mandatory gate: **FB-0 full-Build attempt 5** on frozen baseline `cdc39105e92937e0085410a656346471b2dc6834` with process-only planner `/tmp/swe-planner-f06-sdk-20260913` (SHA256 `581862f59722400ee0999176339dfb99f8dc67cd7901bf39d9861caa5df4d4`). F06 is verified on CURRENT runtime; do not divert into durable AgentField publication before the autonomous full-Build critical path unless a fresh runtime defect requires it.

### Acceptance evidence

- L3-24 (`qa-synthesizer-fcm-smart-l3-24`): historical strong positive evidence. Full issue reached coder -> reviewer block -> repair -> second review -> verifier; bounded two-file delivery; independently inspected; canonical pytest was unavailable, so acceptance had an explicit validation limitation.
- L3-25 (`dag-unknown-dependency-fcm-smart-l3-25`):
  - exact commit `00517ce677e96b17bcd462bd46b9e5f9fa620674`;
  - base `2c374989b39d0b53b34ef33fd2ba6289e74194ae`;
  - changed only `swe_af/execution/dag_utils.py` and `tests/test_dag_utils.py`;
  - independent exact-SHA verification on 2026-09-12: `git diff --check` PASS; deterministic acceptance/adversarial oracle 6/6 PASS; all 4 committed regression test functions 4/4 PASS; `compileall` PASS;
  - canonical pytest remains unavailable and no dependency was installed.
  - Verdict: this bounded task is independently executable-accepted with an explicit canonical-pytest environment limitation. It is not broad L3 PASS.
- L3-26: failed cheap baseline; zero accepted deliverable.
- Fresh PR-1 / accepted task 1/3 (`checkpoint-completed-issue-before-cancel`):
  - execution `exec_20260912_094308_za4vz5sv`, clean target base `2c374989b39d0b53b34ef33fd2ba6289e74194ae`, route `fcm/fcm`;
  - runtime duration `2,092,499 ms` (~34.9 min); coder artifact reports `fallback_count=1`, `peak_models_routed=3`; numeric token/RUB cost is not recorded -> `EVIDENCE_MISSING`, not zero;
  - model-produced commit `8d28dcd49176d640853f333023ac842935b3a528` changed `go/internal/dag/executor.go` + `executor_test.go`, but reviewer explicitly noted the resume test did not actually exercise resume; reviewer PASS was therefore insufficient;
  - operator BMAD test-design / TDD review replaced the weak test with a deterministic downstream-interruption -> real `WithResume(true)` oracle. RED reproduced lost completed-state; minimal GREEN added immediate checkpoint after completed/failed/skipped results;
  - BMAD test-design exposed a false-positive reviewer gate: the model-authored resume test did not execute the real resume path. A deterministic RED reproduced the gap; the bounded repair aligned the task branch with the already-proven CURRENT recovery semantics and corrected the stale baseline threshold assertion to test the real invariant (downstream issue `c` never executes despite bounded repair attempts);
  - final exact task branch head `0cfe48c` is clean and changes only `go/internal/dag/executor.go` + `executor_test.go` relative to base;
  - exact-commit targeted tests `TestResumeAfterCancellationDoesNotRepeatCompletedIssue|TestLevelFailureThresholdAborts` PASS;
  - exact-commit full `/usr/local/go/bin/go test ./... -count=1` PASS;
  - independent temporary oracle (not committed into the candidate) exercised interruption -> checkpoint -> real resume and observed effect counts `a=1`, `b=1`; PASS, then the oracle file was removed before commit;
  - `git diff --check` PASS; duplicate effects observed: 0;
  - verdict: independently accepted fresh task **1/3** with exact tested identity = delivered identity `0cfe48c`; it also satisfies the milestone's recovery/no-duplicate and multi-file task requirements.
- Because L3-26 failed after L3-24/L3-25, those historical positives alone do not satisfy the fresh production streak; PR-1 starts the new controlled streak at 1/3.

### Supporting systems / non-P0 debt

- Coding Station functional API remains unhealthy from the operator surface (`stationHealth` / `stationReady` Gateway Timeout) despite one healthy API + one healthy runtime container. This is not the current critical path because the workforce Go validator works.
- AgentField operator gateway is intermittently `Bad Gateway`; direct authenticated control-plane readback from the workforce works.
- Control-plane stale execution state is now **active P0 reliability evidence for full Build**, not deferred debt. FB-0 attempt 1 showed `run_git_init` with a finished OpenCode process and an output artifact while both child and parent remained `running`; the failed parent later spawned orphan continuation processes. Any retry must therefore use effect/process readback, not control-plane status alone, and must prove no orphan process shares the target workspace.
- SourceLoop/durable Git identity does not yet represent the entire tested live product delta. Canonicalization is required after correctness baseline, not before.
- Python canonical pytest is not provisioned in the current workforce. Do not install it merely to manufacture a green historical task; prefer tasks with already-provisioned canonical validators for the next acceptance streak.
- Process/zombie/resource hygiene is an operational risk; escalate to the critical path only if fresh task evidence shows spawn/resource pressure.

## Production-Ready Milestone

Production acceptance is defined at the **full `orch.Build` lifecycle**, not at the `implement_issue` sub-harness. `implement_issue` remains a useful micro-diagnostic primitive but its PASS does not advance the autonomous-production milestone.

First autonomous production-readiness gate requires a fresh controlled streak of **3 consecutive independently accepted full-Build engineering tasks** on a frozen/currently proven stack, with **zero operator edits to task product code between Build START and terminal state**:

- planning / plan review exercised by the Build lifecycle;
- issue execution through the normal DAG path;
- coder/reviewer repair loop available and exercised where needed;
- issue advisor + bounded retry available;
- replanning available for unrecoverable issue failure;
- integration validation enabled;
- verifier -> generated fix issue -> re-execute -> reverify available;
- CI/finalization path enabled where a safe local/ephemeral CI surface exists; external PR/CI side effects require explicit scope and are not manufactured merely to claim coverage;
- at least one real multi-file task;
- at least one interruption/partial-result recovery task;
- zero duplicated non-idempotent effects;
- exact tested identity = delivered identity for every task;
- canonical deterministic validation before semantic reviewer acceptance where technically available;
- independent executable oracle/verifier;
- task/provider/model/trajectory recorded;
- total model cost and wall time recorded for every task.

If a repairable internal task defect occurs (compile/test failure, reviewer finding, unmet verifier criterion), the Build must repair or exhaust its own bounded recovery policy. Operator repair of task code invalidates that run as autonomous acceptance and instead becomes framework-debug evidence.

After 3/3 full-Build tasks, materialize the accepted source into canonical Git/SourceLoop, create/rebuild a clean runtime from that exact SHA, and replay the controlled ladder. Only the clean-runtime replay can close the first production-ready milestone.

## Critical Path / Pareto Order

### Gate PR-1 — recover one stable execution route

Status: DONE

Evidence:
- pre-fix execution `exec_20260912_093427_ti492y5g` failed before any FCM request;
- A/B canaries proved the clean target repo lacked runtime-owned `fcm` provider config;
- runtime overlay now owns the FCM provider contract; deterministic RED->GREEN test and full CURRENT Go suite PASS;
- rebuilt planner `/tmp/swe-planner-fcm-overlay-20260912` (SHA256 `0224447d...`) is active on PID `424603`;
- post-fix real task `exec_20260912_094308_za4vz5sv` reached FCM, edited source and completed.

### Gate MICRO-1 — issue-level diagnostic evidence

Status: DONE / DOES NOT COUNT TOWARD FULL-BUILD STREAK

Task: `checkpoint-completed-issue-before-cancel` via `implement_issue`.

Evidence:
- exact target base `2c374989b39d0b53b34ef33fd2ba6289e74194ae`;
- final local task branch head `0cfe48c`;
- two-file bounded diff;
- exact-commit targeted recovery tests PASS;
- full exact-commit `go test ./... -count=1` PASS;
- independent real interruption -> checkpoint -> resume oracle observed effects `a=1`, `b=1`; PASS;
- duplicate effects = 0;
- wall time ~34.9 min; numeric token/RUB cost remains `EVIDENCE_MISSING`.

This is strong component evidence for the issue-level coding/recovery seam, but because the run used `implement_issue` and required operator repair after a false-positive model test, it is **not** an autonomous full-Build acceptance task.

### Gate MICRO-2 — issue-level self-repair gap

Status: DONE / BLOCKING EVIDENCE

Task: `openclaw-hitl-enables-build-approval-pr3` via `implement_issue`, execution `exec_20260912_141025_thblod44`.

Fresh terminal evidence:
- control plane status `succeeded`, but result `success=false`;
- duration `2,679,677 ms` (~44.7 min);
- commits `0d12147de8fa4109f9046150f51e0f7066fedf8d`, `f53048294d6cb1c3dd91092ed66d9bbf20cd1b4d`;
- verifier correctly reported `existing internal/orch tests pass = false` because `build_test.go` used `os` without importing it;
- reviewer still marked the semantic code path approved and explicitly noted the compile blocker;
- delivery also committed generated `.agentfield-out-598283605/.agentfield_output.json`, causing `GIT_DELIVERY` unexpected-file / dirty-worktree failure;
- `implement_issue` did **not** run a verifier-fix-reverify cycle for the repairable missing import; it terminated with `success=false`.

Decision: do not manually repair this task and do not count it toward production readiness. This is direct evidence that the issue-level harness is intentionally insufficient as the top-level autonomous acceptance boundary.

### Gate FB-0 — first autonomous full-Build canary

Status: ACTIVE / P0

Fresh attempts:
- Attempt 1: `exec_20260912_163419_wcpjjtdr` on frozen local baseline `e03ee19d1940a29318ebd9f820c92be7fe4931f6`. `build -> plan -> run_product_manager` failed immediately because the Go PM took AgentField direct AI and rejected `fcm/fcm` as unsupported. Root cause was a Go-port parity/runtime-contract defect: `deps.AI` was always non-nil, so `runtime=open_code` was ignored for PM. RED->GREEN fix now routes OpenCode PM through the harness; targeted planning tests and full Go suite PASS.
- Attempt 1 also proved structured-output self-repair exists in git-init: the model wrote a JSON Schema instead of an instance, the harness detected it and launched an incremental continuation with the exact validation contract. That continuation repeated the same schema mistake, so recovery capability exists but is not yet sufficient on this trajectory.
- Attempt 2: `exec_20260912_164532_i9s249f1` proved the PM routing fix functionally — a live Product Manager OpenCode process started with `-m fcm/fcm`. The attempt is **invalid acceptance evidence** because an orphan git-init OpenCode process from attempt 1 was concurrently mutating the same sacrificial repo. It was stopped and the repo was restored to exact baseline before the next attempt.
- Attempt 3: `exec_20260912_164957_ikbjlc8x` did not reach the agent; it failed after 15s with `agent_unreachable` because the manually restarted planner had registered callback `http://localhost:8005`. The node contract explicitly requires `AGENT_CALLBACK_URL` in containers. Planner was restarted with `AGENT_CALLBACK_URL=http://172.16.22.7:8005`; control-plane now reports that callback active. This attempt contains no task-code effect and is infrastructure evidence only.
- Attempt 4: `exec_20260912_165200_an8f8ym5` reached the correct PM OpenCode/FCM harness path but failed structured output: `Schema validation failed after 2 retry attempt(s)` / output file missing. Compatibility A/B then proved the base AgentField/OpenCode/FCM contract works for a simple schema, while exact `schemas.PRD` in single mode makes the routed model write the JSON Schema itself instead of a PRD instance. The same exact PRD contract in incremental mode PASS (`parsed=true`) after field-by-field construction. PM policy was therefore changed only for OpenCode: `SchemaMode=incremental`; non-OpenCode PM harness paths retain default policy. RED/edge-case tests, full planning package, full `go test ./... -count=1`, and `git diff --check` PASS.
- Current clean precondition: sacrificial repo `/tmp/swe-af-fullbuild-current-20260912` is clean on baseline branch `fullbuild-current-baseline-20260912`, refreshed exact local SHA `cdc3910`; full target `go test ./... -count=1` PASS after the F08 fault gate; no live OpenCode process targeted that workspace at planner reload; planner `/tmp/swe-planner-fullbuild-f08-20260912` (SHA256 `f5adc31c29b3383f463f34df6863d3e06767e7904930b45b5f06463dd4a0d0de`) is healthy and registered at the routable callback.

Goal: exercise the actual `swe-planner.build` / `orch.Build` lifecycle on a frozen materialization of CURRENT tested SWE-AF source, with zero operator edits to task code during the run.

Canary engineering objective: add a deterministic pre-review validation contract so trivial compile/test failures are fed back into the coding repair loop before semantic reviewer spend. This directly addresses the observed missing-import waste without replacing the existing reviewer/verifier loops.

Configuration policy:
- use the proven planner/runtime route and `fcm/fcm`;
- `enable_replanning=true`;
- `enable_issue_advisor=true`;
- `enable_integration_testing=true`;
- bounded retries/replans/verify-fix cycles remain enabled;
- deterministic Git enabled;
- `enable_learning=false` for a reproducible baseline;
- external GitHub PR/CI side effects disabled for this first local canary only; CI semantics remain a later safe-surface gate, not silently claimed as covered.

DoD:
- Build performs planning -> DAG execution -> coding/review/repair -> integration -> verifier/fix lifecycle as applicable;
- operator performs zero task-code edits between Build START and terminal state;
- any repairable compile/test/reviewer/verifier failure is repaired by Build itself or the run fails with evidence;
- exact final source identity and artifacts are recoverable;
- canonical Go validation + independent oracle run after terminal state;
- no generated harness artifact is accepted as product source;
- cost/latency/repair/replan counts recorded.

### Gate FB-1 / FB-2 — full-Build streak 2/3 and 3/3

Status: PENDING FB-0

Repeat on different real engineering tasks with the same frozen/proven control stack. At least one full-Build run must exercise interruption/recovery with zero duplicate effects. No cost optimization before 3/3 full-Build correctness.

### Gate PR-5 — state truth + durable source

Status: PENDING 3/3

DoD:
- reconcile stale `active` execution state against real liveness/artifacts/effects;
- capture only accepted source deltas; exclude generated/noise files;
- materialize exact accepted Git SHA;
- canonical tests PASS on durable SHA;
- clean runtime built/materialized from that SHA;
- controlled 3-task ladder replays successfully;
- container-only required product deltas = 0.

### Gate PR-6 — cost optimization

Status: PENDING reproducible correctness

Metric: `total inference cost / independently accepted task`.

Then, and only then:
- compare cheap vs stronger model on paired frozen tasks;
- optimize route/continuation count/latency;
- investigate L3-26 cheap structured-output/provider path if it still dominates cost;
- generalize AgentField contract completion only if fresh acceptance evidence proves unresolved-obligation recovery is the dominant blocker.

## Contract Completion Decision

Contract completion is a hypothesis for reliability, not an architectural mandate.

Reuse before build:
- SWE already has finish-only/same-worktree continuation semantics.
- AgentField upstream already has incremental field recovery (`DiagnoseFieldFailures` -> `BuildIncrementalFollowup` with session continuation).
- Do not create a parallel workflow/obligation engine.

If this becomes the proven blocker, first deterministic RED tests must cover:
1. observed SATISFIED overrides model MISSING;
2. observed MISSING overrides model SATISFIED;
3. UNKNOWN forbids mutation;
4. mutation observed after timeout is not repeated.

No implementation closes without a deterministic test for every guarantee.

## Engineering Method for 30-Minute Batches

BMAD:
- Local project entry: `bmad-help`.
- Current risk/evidence planning: `bmad-testarch-test-design`.
- Implementation only after a source defect is proven: use the locally available BMAD implementation skill; public BMAD v6.12 calls the official implementation workflow **Build** (`bmad-build`) and explicitly chooses ceremony after investigation. Do not add ceremony before evidence.
- Review once after GRESN; triage findings with verdict + evidence; do not repeat settled reviews.

External skill practices used as method, not installed dependencies:
- systematic-debugging: reproduce -> observe -> isolate layer -> one falsifiable hypothesis -> one minimal experiment -> update diagnosis;
- test-driven-development: deterministic RED before behavior-changing implementation where practical;
- verification-before-completion: run fresh executable proof before PASS/DONE;
- condition-based-waiting/effect readback for async/timeout work;
- root-cause tracing across component boundaries before mutation.

Do not install BMAD/skill frameworks merely to say they were used.

## Batch Anti-Drift — CURRENT

GLOBAL NORTH STAR:
working SWE/SWE-AF with independently accepted real engineering tasks.

CURRENT BLOCKER:
F06 restart/state convergence is **VERIFIED/CLOSED for the current process-only runtime**. The production milestone remains blocked because autonomous full-Build acceptance is still 0/3; the next unknown must come from FB-0 attempt 5, not more infrastructure speculation.

THIS BATCH:
start FB-0 attempt 5 from clean frozen baseline `cdc39105e92937e0085410a656346471b2dc6834` using planner `/tmp/swe-planner-f06-sdk-20260913`, route `fcm/fcm`, and the already-proven recovery configuration. After START perform zero operator task-code edits; observe process/artifact/effect truth and stop at the first new evidence-backed lifecycle blocker.

NORTH-STAR DELTA:
F06 removed the stale-running/orphan safety blocker: Go Agent processes now have stable per-process identity, shutdown cancels in-flight reasoners, and the same restart fault converges terminally with zero overlapping mutator. The critical path returns to autonomous Build correctness.

STOP CONDITION:
either FB-0 attempt 5 reaches autonomous terminal success and then passes canonical Go validation + independent oracle, or it identifies the first new full-lifecycle blocker. Write that blocker back before any further framework mutation. Durable publication of local AgentField commit `b1602458...` is required before a clean deployed-runtime milestone, but is not allowed to displace the current correctness gate.

## Acceptance Metrics Per Task

Record:
- task ID/input and base SHA;
- provider/model/route;
- actual changed files/diff/commits;
- canonical tests and exact commands;
- independent oracle result;
- final acceptance verdict;
- model calls/tokens/cost where available;
- tool calls;
- wall time;
- continuations/repairs;
- duplicated effects;
- UNKNOWN states;
- recovery success;
- tested SHA and delivered SHA;
- provenance completeness.

## Write-Back Rules

- `dev/PLAN.md` owns current project plan/state/decisions.
- Product code: container/exact-source first -> tests/oracle -> exact delta -> SourceLoop/Git publication.
- PLAN-only Git commits must never be mistaken for loaded product identity.
- Do not create a second plan file.
- Do not overwrite unrelated dirty source.
- Historical checkpoints remain recoverable from Git history; keep this SoT concise and CURRENT.

## Current Next Move

After explicit supporting-platform scope authorization, implement **AgentField Go SDK process-instance identity parity** in the AgentField owner source, not in SWE-AF: fresh UUID4-like ID once per Go Agent process; stable within that process; propagate it through registration and heartbeat; add deterministic parity tests matching the Python SDK contract. Validate the exact AgentField SDK/control-plane source, update the existing DEV stack to exact tested pins, then rerun `exec`-style F06 restart convergence with planner `/tmp/swe-planner-fullbuild-f08-20260912` semantics. Do not resume FB-0 until the pre-restart execution leaves `running` deterministically and process/workspace readback proves zero overlapping mutators.
