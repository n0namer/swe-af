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
- **Native architecture maximization:** for any active fault family, prefer and exercise the full relevant SWE-AF + AgentField mechanisms already present in the architecture before introducing narrower local substitutes. A local fix is insufficient if the native lifecycle/identity/persistence/recovery/validation primitive that owns the risk is available but bypassed. “Maximum use” means maximum *relevant* native capability, not enabling unrelated features for ceremony. End-to-end acceptance must prove the relevant native mechanisms compose correctly across SWE-AF, AgentField, SourceLoop/FVE, canonical validators and independent oracles.
- **Acceptance Runtime Equivalence Gate:** diagnostic/process-only runtimes may reproduce or isolate a fault, but they never count toward the North-Star full-Build streak unless they satisfy both (1) **dependency/bootstrap equivalence** — the same required OpenCode wrapper/binary, provider config, broker env contract, runtime defaults and filesystem/package assets as the native SWE runtime — and (2) **lifecycle equivalence** — the planner/control-plane process is owned by the normal durable runtime/package/service lifecycle, survives operator-session expiry, has stable reachable callback identity, and leaves independently readable terminal/effect state after the initiating operator session disappears. A managed terminal session, ad-hoc localhost control-plane, or binary launched in another service's container is diagnostic evidence only unless these equivalence criteria are explicitly proven before START.
- **Preflight-before-Build rule:** before any run can count, executable preflight must prove exact source/binary identity, clean target worktree, zero pre-existing mutators, real OpenCode invocation through the native wrapper with the frozen model route, reachable control-plane callback, and durable post-session execution readback. If any of these are missing, stop with `RUNTIME_MATERIALIZATION_GAP`; do not start a Build and reinterpret the resulting failure as product/framework evidence.

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

### AgentField takeover reconciliation — 2026-09-13 fresh BMAD/graph evidence

The F06 **process-only** proof above remains valid for exact local commit `b160245833ec51f5296905c32e49260a62c76e26`, but it no longer closes the broader restart / generation-ownership / stale-state / late-effect fault family. Fresh takeover readback and independent upstream graph analysis reopened that family before FB-0.

CURRENT owner-source readback:
- Exact AgentField owner clone `/tmp/agentfield-f06-fix` is clean at local-only commit `c0923acdfca043c2c07e3d34daaa09e2a7e41d38` (`fix(agentfield): harden restart generation and stale-state recovery`), 27 changed files including migration `035_execution_instance_id.sql`; unsupported Python ABA experiment and generated `.archsteer/` are absent.
- Frozen SWE repo `/tmp/swe-af-fullbuild-current-20260912` is independently clean at `cdc39105e92937e0085410a656346471b2dc6834`. Full SWE Go suite with a temporary module replace to `/tmp/agentfield-f06-fix/sdk/go` (exact c092 HEAD) PASS after one localized baseline concurrency-test flake; `TestUnlimitedConcurrencyWhenZero` then passed 5/5 on both candidate and pinned SDK, and the evidence-changed full rerun PASS.
- Exact planner built from frozen SWE + c092 SDK: `/tmp/swe-planner-agentfield-c092-20260913`, SHA256 `c923a3653cceddbbc0672f6c6b8b673f3b819ef15751878e255014a9d0732e79`. BMAD mode remains System-Level `bmad-testarch-test-design`; no separate test-plan document exists because this `PLAN.md` is the canonical owner.

Independent GitHub GraphQL/code search against `Agent-Field/agentfield` proved that the deployed/frozen AgentField base `4d337c1ae5104418311fcba414a1c2f85c2abb89` (2026-08-24) predates a cluster of merged fixes for the same fault family:
- PR #1000 / `93994097fabc8999259afe6c10431a64e0437fca`: Go/TS graceful shutdown drains control-plane-dispatched executions under `AGENTFIELD_SHUTDOWN_TIMEOUT`; accepted async work is tracked, drained, then cancelled+settled on deadline.
- PR #1001 / `9a517debd0eac5b0e01c07328a1b1c3c827cc583`: **async control-plane lane** rejects before persistence and drains/fails queued async jobs on control-plane shutdown. This is directly relevant because SWE-AF uses the async execution API.
- PR #1004 / `c477f6d3263542ff872805918bfb40941edd5bf7`: executions/workflow executions are stamped with serving `instance_id`; re-registration orphan cleanup becomes instance-scoped and deferred through a drain window. The frozen base stores only `agent_node_id`, so it cannot distinguish old vs replacement execution generations during reap.
- PR #1011 / `b01e8315dee590f9dd6d50d939e66c08c80a5766`: adversarial follow-up restores Go SDK **notify-then-drain** ordering; setting shutdown admission too early could reject work arriving during the control-plane notify and turn it into a non-retryable failure.
- PR #1046 / `78215f17ac762d12729177d17292aeb6ad960800`: stale workflow cleanup must consult the paired execution activity clock. Upstream reproduced live heartbeats followed by premature workflow reap and late completion HTTP 409. SWE runs have already lasted ~35–45 minutes, so this boundary is P0-risk until CURRENT cleanup applicability is executable-read back.

Neighboring but not automatically current-path P0:
- PR #1031 / `2d7fc7264a7aaf5dd36fa4ed05ad55f8a297e117` adds a multi-replica orphan-reap kill switch because a new `instance_id` may be a sibling replica rather than a replacement. CURRENT workforce compose readback shows `container_number=1`, so keep this P1 unless replica topology changes.
- PR #1033 / `2638b9e92a1f9ee6c6f9a3cd223da231bf1ece86` extends admit-before-persist / shutdown terminalization to sync, restart and MCP lanes. Its own contract states the async lane was already protected by #1001; retain as sibling regression coverage, not an excuse for a broad upgrade.

P0 test-design matrix for reconciliation (a test must discriminate broken vs working behavior):
1. Go process `instance_id`: stable within one Agent process, unique across processes, present on registration + heartbeat wire payload.
2. Go in-process cancel ownership: duplicate execution ID cannot let an older release remove the newest generation owner.
3. stale process messages: old/missing instance heartbeat/status cannot mutate a modern replacement after ownership changes; legacy compatibility must be explicit rather than accidental.
4. Go shutdown: accepted async reasoner/skill work is tracked; notify/admission/drain ordering is race-safe; deadline cancellation settles terminal status; no post-shutdown untracked owner is admitted.
5. control-plane async shutdown/admission: queue/capacity rejection persists zero execution/effect state; pool shutdown terminalizes already-persisted accepted work rather than abandoning `running` rows.
6. persisted execution generation: execution + workflow rows carry serving `instance_id`; orphan reap targets only the departing generation (plus explicitly supported legacy rows), not the replacement.
7. stale cleanup: recent activity on the paired execution prevents workflow reap; genuinely stale pairs still converge terminally.
8. real restart fault: SIGTERM/restart on frozen SWE produces one terminal outcome, zero overlapping OpenCode mutator and clean workspace; tested identity = exercised identity.

CURRENT validation/capability state:
- Exact owner clone `/tmp/agentfield-f06-fix` is now an authorized and exercised container-first route. DEV target `agentfield-dev-workforce` covers the clone for stale-safe patch/test work; ordinary `vps-terminal` was independently proven to provide the bounded local Git write path needed for cleanup/staging/commit on the same target. Coding Station is **not** on the critical path.
- Unsupported Python ABA experiment files were restored to HEAD; generated `.archsteer/` and the in-repo temporary patch artifact were removed; the exact candidate contains only the proven Go/control-plane delta plus migration `035_execution_instance_id.sql`.
- Fresh exact-tree validators before commit: `go test ./internal/storage ./internal/handlers ./internal/server -count=1` PASS, full `sdk/go` `go test ./... -count=1` PASS, `git diff --check` PASS, staged `git diff --cached --check` PASS.
- One local-only exact AgentField commit now exists: `c0923acdfca043c2c07e3d34daaa09e2a7e41d38` (`fix(agentfield): harden restart generation and stale-state recovery`), 27 files, 1,163 insertions / 133 deletions, including migration `035_execution_instance_id.sql`; post-commit `git status --short` is clean.
- Targeted `-race` remains an environment validator limitation rather than a product failure: legacy `boltdb/bolt@v1.3.1` trips Go checkptr before the target test, while `-race` with checkptr disabled exceeds the 120s execution ceiling without verdict.

Takeover gate verdict: **BEHAVIOR MATRIX CLOSED LOCALLY; EXACT CANDIDATE CREATED**. FB-0 remains blocked only until this exact AgentField commit is validated against frozen SWE and the real F06 SIGTERM/restart fault is re-injected with tested identity = exercised identity.

#### BMAD quick-dev + adversarial/edge reconciliation — 2026-09-13

BMAD usage in this takeover is evidence-driven, not ceremonial:
- `bmad-help` remains the entrypoint and `bmad-testarch-test-design` owns the system-level risk matrix above.
- `bmad-quick-dev` was loaded. Its full workflow requires project `_bmad/bmm/config.yaml` plus project customization state, but `/src/swe-af/_bmad` does not exist. No BMAD project/config/plan files were created because `PLAN.md` already owns project state; use the Quick Dev Ready-for-Development rules (actionable file-level tasks, ordered dependencies, explicit executable ACs) as the implementation standard once the exact owner source is writable/testable.
- `bmad-review-adversarial-general` and `bmad-review-edge-case-hunter` were applied to the exact 12-file uncommitted AgentField diff and compared against the post-base upstream fixes.

Historical adversarial/edge findings — resolved or explicitly dispositioned before exact commit `c0923acdfca043c2c07e3d34daaa09e2a7e41d38`:
1. The dirty Go patch sets `shuttingDown=true` **before** notifying the control plane. Upstream #1011 explicitly reversed this after adversarial testing: work arriving while shutdown notification is in flight must remain admissible; admission closes after notify. The current local order can manufacture non-retryable 503 failures.
2. The dirty patch cancels current registrations but has no `executionWG`/bounded graceful drain/post-cancel settlement contract from #1000. Accepted async work is therefore not proven to finish or terminalize before shutdown returns.
3. The local shutdown tests only prove “after shutdown, new work is rejected”; they do not cover the critical notify-window interleaving that #1011 demonstrated.
4. The deployed base still lacks #1001 async admission-before-persistence. Queue/concurrency rejection may persist execution/workflow/payload state that should not exist.
5. The deployed base still lacks #1001 worker-pool shutdown terminalization. Accepted/queued async jobs can be abandoned as non-terminal `running` state.
6. The local delta does not add persisted execution/workflow `instance_id`, migration `035_execution_instance_id.sql`, or storage read/write propagation from #1004. Node process identity alone cannot make orphan reap generation-safe.
7. The local delta does not implement #1004 instance-scoped/deferred orphan reap. A replacement generation can still be affected by node-wide cleanup.
8. The local delta does not implement #1046 paired execution-activity protection for stale workflow cleanup. Long-running SWE work can still be falsely reaped while the paired execution is active.
9. `rejectStaleAgentInstance` currently rejects a **missing** incoming `instance_id` after a modern node has registered. That is stricter than the original takeover contract (“reject when both IDs are non-empty and differ”) and needs an explicit compatibility decision plus discriminating tests; do not ship the stricter rule accidentally.
10. Heartbeat handling now performs an unconditional authoritative `GetAgent` before the presence/cache fast path. That changes hot-path storage/error behavior and requires positive current-instance, legacy-empty-instance, cache-hit and storage-error regressions, not only stale-negative tests.
11. The dirty heartbeat/status tests cover stale and missing-ID rejection but do not pin the positive paths: current instance accepted and legacy stored-empty instance accepted.
12. Skill shutdown tracking is keyed by execution ID; if skill calls without `X-Execution-ID` are valid, that path is not proven cancellable/drain-safe.
13. The fallback/internal async reasoner path can observe shutdown after registration, but the diff does not prove terminal-status propagation for every non-HTTP invocation path.
14. Python `owner_task` deregistration is a plausible sibling ABA guard, but its current test exercises the registry API directly. Upstream main still uses unconditional `pop`; retain the Python production change only if a real caller-level interleaving is RED on the exact base and GREEN with the fix.
15. `.archsteer/` is generated analysis output and must be excluded from the exact product commit.
16. No fresh race/full-suite executable evidence exists for the dirty delta; the earlier GREEN for `b1602458...` cannot be inherited by these later changes.

Historical local-delta disposition before c092 closure (superseded by the exact tested commit and runtime proof):
- **candidate to keep after re-test:** Go cancel-registration pointer ownership / duplicate-ID ABA fix;
- **must redesign before acceptance:** Go shutdown/admission code, using `notify -> close admission -> drain -> deadline cancel -> settlement` rather than the current early-close behavior;
- **hold for explicit compatibility oracle:** stale heartbeat/status handling when incoming `instance_id` is missing;
- **experiment only:** Python owner-aware deregistration until caller-level RED exists;
- **missing P0 behavior:** #1001 async admission/shutdown, #1004 persisted generation/reap, #1046 stale workflow activity; #1000/#1011 define the Go drain/order contract.

Exact-source capability attempt after user authorization:
- durable current target entry is `agentfield-dev-workforce`, revision `1`, selector `com.docker.compose.project=edshqtkwskg3lrczekhcmd71`, capabilities unchanged, `live_patch_roots=["/src"]`;
- a minimal typed registry upsert was attempted with revision `2`, preserving selector/capabilities/checks and adding only `/tmp/agentfield-f06-fix` to `live_patch_roots`;
- the server rejected it with `target_registry_scope_denied: target_id is outside the server-owned mutation allowlist`; no registry state changed;
- a readback of all registered `live_patch_roots` found no existing target that lawfully covers `/tmp/agentfield-f06-fix`;
- direct editing of `/var/lib/vps-terminal/targets.runtime.json`, operator redeploy, or creation of a duplicate workspace/runtime is intentionally not used as a bypass.

This sharpens the blocker from “approval missing” to **server-owner capability gap**: the user authorized the bounded scope widening, but the current typed mutation service is not callable for this target. Product mutation remains blocked until the owner of the DEV terminal mutation allowlist exposes this target (or an equivalent exact-source typed route) through supported configuration/deployment.

#### DEV terminal owner-route recovery — 2026-09-13

The owner-layer fix was then executed through the canonical `n0namer/vps-terminal` source rather than by overwriting a hidden Coolify env value:
- owner rules `AGENTS.md` + `ERRORS.md` were re-read before mutation; they require exact source identity, container-first/debug evidence and independent activation readback;
- Coolify exposes the `TARGET_REGISTRY_MUTATION_ALLOWLIST` key but masks its value, so direct env overwrite was rejected as unsafe because it could erase unrelated allowed targets;
- exact deployed owner base is `n0namer/vps-terminal@7ed1daa2824096e2624c0025f7ee9571fa78ff87`, app `vps-terminal-dev` / `tsnhqqr60rcacdv8kfiw4eqf`, canonical DEV branch `archops/k4b-container-cleanup`;
- isolated branch `archops/agentfield-registry-scope` was created from that exact SHA. Candidate head `e1b671fd279eee6c03a640ae37fe8ebba27ce501` changes only `gateway/live-aci.mjs`, `gateway/server-live.mjs`, and existing test owner `gateway/test/live-aci.test.mjs`;
- candidate adds pure `mergeAllowedTargetIds(...)`: configured env entries are preserved/deduplicated and required `agentfield-dev-workforce` is added exactly once. `server-live.mjs` composes the required target additively instead of replacing the hidden env value. Focused regression asserts both preservation and deduplication;
- PR `n0namer/vps-terminal#160` targets the exact currently configured DEV branch. It remains **open/unmerged** pending executable validation;
- canonical GitHub Actions `ci` cannot currently validate the candidate: both push and PR runs terminate `startup_failure` with **0 jobs**, and the same startup failure is observable on the branch baseline. Classify this as `VALIDATION_BLOCKER`, not application test FAIL;
- Coding Station remains unavailable (`Gateway Timeout`), so it did not provide a substitute exact-source validator;
- registered target `vps-terminal-dev-gateway` exposes `/app/gateway` with `node_check`, but stale-safe apply of the exact candidate failed `EROFS`; the image source is read-only. No direct registry/runtime file bypass was used;
- with user-authorized DEV deploy scope, Coolify desired source was temporarily changed to candidate branch/head and a forced deployment was requested. Initial request stayed queued/not-applied; one safe retry with `instant_deploy=true` was issued after independent runtime readback still showed `7ed1daa...`. No further retry is permitted without new evidence;
- Coolify logs then showed `ApplicationDeploymentJob RUNNING`; host inventory observed the old gateway removed and the DEV action endpoint temporarily unavailable, proving the deployment reached runtime replacement rather than remaining a mere queue acknowledgement. At the latest write-back point the replacement had **not yet reached a verified healthy candidate identity**. Do not count health/config desired SHA as candidate PASS until an image/source fingerprint proves `e1b671f...` is loaded and canonical validation executes on that identity.

BMAD trace gate for this operator fix:
1. candidate source identity = `e1b671fd279eee6c03a640ae37fe8ebba27ce501`;
2. executable requirement = additive allowlist preserves arbitrary configured entries and contains `agentfield-dev-workforce` once;
3. required validation = focused helper test + repository `npm run check` (or exact equivalent on the candidate image), syntax/runtime readiness, then typed target-registry upsert readback;
4. activation proof = DEV gateway loaded source/image corresponds to candidate/merged identity; `/health` alone is insufficient;
5. functional oracle = revision-guarded upsert of `agentfield-dev-workforce` succeeds and readback shows only `/tmp/agentfield-f06-fix` added to its `live_patch_roots` while every other target field is unchanged;
6. rollback = restore app branch/SHA to `archops/k4b-container-cleanup@7ed1daa...` if candidate validation or activation fails. Do not merge PR #160 until this gate is GREEN.

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
- **F06/F07 are VERIFIED/CLOSED for the CURRENT single-workforce async/restart path at exact AgentField commit `c0923acdfca043c2c07e3d34daaa09e2a7e41d38`.** Deterministic local contracts cover #1000/#1001/#1004/#1011/#1046 plus process identity, stale-message fencing and cancel-owner ABA. Full storage/handlers/server + Go SDK suites PASS; frozen SWE full suite with the c092 SDK PASS. Real candidate-runtime fault injection also passed: hard-killing planner instance `8de0300c...` left the accepted execution temporarily running, replacement instance `e4357c65...` registered, and after the native 60s drain window the old-generation execution became terminal `failed` with `agent_restart_orphaned`; a separate real SIGTERM on replacement PID while Build execution `exec_20260913_201344_93msxfav` was `running` converged terminal `failed(node_unavailable)` in `19.676s`, before the 30s shutdown deadline, with no live OpenCode mutator and a clean sacrificial repo. A subsequent replacement registered unique instance `c51b3a1a...`.
- F05 broader dependency/network-provider crash coverage remains partial, but the process hard-crash/restart boundary required by the takeover is executable-proven. P1 siblings remain multi-replica false-reap (#1031) if topology exceeds one workforce replica and non-async dispatch-lane parity (#1033). No AgentField P0 from the takeover matrix remains open on the CURRENT path.

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

Immediate mandatory gate: **AgentField takeover reconciliation is VERIFIED/CLOSED at exact local commit `c0923acdfca043c2c07e3d34daaa09e2a7e41d38`; FB-0 attempt 5 is now the active P0.** Frozen SWE baseline `cdc39105e92937e0085410a656346471b2dc6834` is clean and full-suite compatible with the c092 Go SDK; the exact planner artifact SHA is `c923a3653cceddbbc0672f6c6b8b673f3b819ef15751878e255014a9d0732e79`. Candidate control-plane + planner process-only fault injection proved both hard-crash generation reap and graceful SIGTERM convergence with no overlapping live mutator. The next batch must rehydrate that exact candidate runtime identity, verify the sacrificial repo is clean and has zero live mutators, then run FB-0 attempt 5 with zero operator task-code edits.

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
- Historical control-plane stale execution/orphan evidence remains important, but the specific F06 restart/stale-running contract is **VERIFIED/CLOSED on the current process-only runtime** by AgentField local commit `b1602458...` plus live same-fault re-injection. Continue to require effect/process readback for timeouts/restarts; durable publication of that AgentField delta is pending before a clean deployed-runtime milestone.
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

### Frozen autonomy experiment manifest — preregistered before FB-0 attempt 5

Purpose: prevent task/config/oracle/operator drift while measuring the first autonomous full-Build streak. This is an operational readiness experiment, not a statistical claim about general SWE reliability.

Frozen stack identity for the streak unless a run exposes a proven framework blocker. Attempt 6 exposed one such blocker, so the stack was amended **before Attempt 7 START** and the streak remains/reset to 0/3:
- SWE target/base: `/tmp/swe-af-fullbuild-current-20260912` at `6f5b4382e6231721f60be7045b9d91fd85e34fb5` before each task START. This commit is exactly prior baseline `cdc39105e92937e0085410a656346471b2dc6834` plus the two-file OpenCode PM FCM-overlay framework fix and its regression assertion;
- AgentField exact candidate: `/tmp/agentfield-f06-fix` at `c0923acdfca043c2c07e3d34daaa09e2a7e41d38`;
- planner artifact: `/tmp/swe-planner-fb0-6f5b438-c092`, SHA256 `3dc66b5d7f727f1b3e973e8e2fe6b7b29004853d7079f26fa599e13bacab5cd3`;
- route/model is explicit in every Build request, not inherited from ambient env: `config.runtime="open_code"` and `config.models.default="fcm/fcm"`; this is required because the current workforce has `ANTHROPIC_API_KEY` present while `SWE_DEFAULT_RUNTIME`/model vars and `OPENROUTER_API_KEY` are absent, which otherwise resolves to `claude_code -> sonnet`;
- frozen Build config for the local streak: `{"runtime":"open_code","models":{"default":"fcm/fcm"},"enable_github_pr":false,"check_ci":false}`. All omitted fields keep exact BuildConfig defaults: replanning ON, issue advisor ON, integration testing ON, deterministic Git ON, learning OFF and bounded native retry/replan/verifier-fix policies;
- external PR/CI side effects remain OFF for the local streak unless a later task explicitly requires a safe isolated CI surface;
- acceptance oracle is frozen before START: terminal full `swe-planner.build` result + canonical deterministic validation on the resulting exact source + an independent executable task-specific oracle. Semantic reviewer prose alone never accepts a task.

Intervention freeze from Build START to terminal state:
- forbidden: operator task-code edits, prompt/model/provider/config/retry/oracle changes for the active task, manual continuation that changes task semantics, or starting a second mutator on the same workspace;
- allowed observation: execution/control-plane state, process/effect/workspace readback, logs/telemetry and condition-based waiting;
- infrastructure failure may be diagnosed after post-state readback, but any framework/source/config mutation invalidates that task as acceptance evidence; repair uses deterministic RED->GREEN where practical and restarts the streak from 0/3 on the repaired frozen stack.

Preregistered task ladder:
1. **FB-0 / pre-review deterministic validation** — exact objective already specified below: add a deterministic pre-review validation contract so trivial compile/test failures are fed into the coding repair loop before semantic reviewer spend. Must be a real full Build and independently validated after terminal state.
2. **FB-1 / different multi-file behavior task** — select the first eligible real repository task from the canonical project backlog/task source in stable ascending task order that (a) is not semantically the FB-0 change, (b) requires at least two product/test files, (c) has a provisioned deterministic validator, and (d) does not require unavailable external credentials/services. Do not skip an eligible earlier task because it appears harder. If no eligible canonical task exists, record `TASK_POOL_GAP` before creating/choosing a new benchmark task; do not improvise after seeing FB-0 outcome.
3. **FB-2 / interruption-recovery task** — select the first eligible task by the same stable-order rule whose execution can be subjected to one preregistered interruption after a completed non-idempotent sub-effect. Acceptance additionally requires resume/recovery with that effect count remaining exactly one. If the canonical pool has no such eligible task, record `TASK_POOL_GAP` and define the task before running it, not after observing FB-1.

Streak rule: only consecutive independently accepted full Builds on the same frozen stack count. A task-level autonomous repair inside Build is allowed and measured. Operator/framework repair, invalid environment, or failed independent oracle does not count and resets the operational streak after the repaired/frozen stack is revalidated. `3/3` closes the first operational autonomy milestone only; broader reliability/generalization requires a later held-out evaluation with no framework changes between tasks.

Per-task experiment record must include: task text/ID/source, base SHA, AgentField SHA, planner SHA, route/model, config policy, execution/run IDs, final source SHA/diff, canonical validator command/result, independent oracle/result, wall time, model calls/tokens/cost when available, retries/continuations/repairs/replans, UNKNOWN states, duplicate effects, tested-vs-delivered identity and any evidence gap. Cost is measured now; economic PASS/FAIL remains undefined until an explicit business SLO exists.

### Gate FB-0 — first autonomous full-Build canary

Status: ACTIVE / P0

Fresh attempts:
- Attempt 1: `exec_20260912_163419_wcpjjtdr` on frozen local baseline `e03ee19d1940a29318ebd9f820c92be7fe4931f6`. `build -> plan -> run_product_manager` failed immediately because the Go PM took AgentField direct AI and rejected `fcm/fcm` as unsupported. Root cause was a Go-port parity/runtime-contract defect: `deps.AI` was always non-nil, so `runtime=open_code` was ignored for PM. RED->GREEN fix now routes OpenCode PM through the harness; targeted planning tests and full Go suite PASS.
- Attempt 1 also proved structured-output self-repair exists in git-init: the model wrote a JSON Schema instead of an instance, the harness detected it and launched an incremental continuation with the exact validation contract. That continuation repeated the same schema mistake, so recovery capability exists but is not yet sufficient on this trajectory.
- Attempt 2: `exec_20260912_164532_i9s249f1` proved the PM routing fix functionally — a live Product Manager OpenCode process started with `-m fcm/fcm`. The attempt is **invalid acceptance evidence** because an orphan git-init OpenCode process from attempt 1 was concurrently mutating the same sacrificial repo. It was stopped and the repo was restored to exact baseline before the next attempt.
- Attempt 3: `exec_20260912_164957_ikbjlc8x` did not reach the agent; it failed after 15s with `agent_unreachable` because the manually restarted planner had registered callback `http://localhost:8005`. The node contract explicitly requires `AGENT_CALLBACK_URL` in containers. Planner was restarted with `AGENT_CALLBACK_URL=http://172.16.22.7:8005`; control-plane now reports that callback active. This attempt contains no task-code effect and is infrastructure evidence only.
- Attempt 4: `exec_20260912_165200_an8f8ym5` reached the correct PM OpenCode/FCM harness path but failed structured output: `Schema validation failed after 2 retry attempt(s)` / output file missing. Compatibility A/B then proved the base AgentField/OpenCode/FCM contract works for a simple schema, while exact `schemas.PRD` in single mode makes the routed model write the JSON Schema itself instead of a PRD instance. The same exact PRD contract in incremental mode PASS (`parsed=true`) after field-by-field construction. PM policy was therefore changed only for OpenCode: `SchemaMode=incremental`; non-OpenCode PM harness paths retain default policy. RED/edge-case tests, full planning package, full `go test ./... -count=1`, and `git diff --check` PASS.
- Attempt 5: `exec_20260913_232654_qj96b61g` / run `run_20260913_232654_ayw1aces` on exact c092 control-plane/planner stack failed fail-fast before task mutation: PM child `exec_20260913_232654_p55p1nqi` returned `API error (400): unsupported model "sonnet"`. Root cause is experiment-config drift, not product code: runtime env has no explicit SWE runtime/model vars, `ANTHROPIC_API_KEY` is present and OpenRouter key absent, so `DefaultRuntime()` resolves `claude_code` and role defaults resolve `sonnet`. The preregistered intent said `fcm/fcm`, but the Attempt-5 request omitted explicit `runtime/models`, so the run is invalid acceptance evidence and streak remains 0/3. Exact config tests `TestBuildConfig_OpenCodeProvider|TestResolveRuntimeModels_Precedence` PASS; the frozen manifest now requires explicit `runtime=open_code` + `models.default=fcm/fcm` in every Build request.
- Attempt 6: `exec_20260913_232855_bg4bjnrm` / run `run_20260913_232855_xm5pa9z1` used the corrected explicit `open_code + fcm/fcm` Build config and reached the PM harness, but PM child `exec_20260913_232855_6egehnrk` failed in ~0.5s with `Schema validation failed ... output file was NOT created`, `turns=0`, `messages=0`. Root cause was a framework parity gap: SWE owns `OpenCodeNoInstallPermissionOverlay` with the FCM provider and injects it for coding/advisor OpenCode calls, but `RunProductManager` built OpenCode `RoleOptions` without `Env["OPENCODE_CONFIG_CONTENT"]`; ambient config is absent. Deterministic RED extended `TestProductManagerOpenCodeUsesHarnessWhenDirectAIIsAvailable` to require that overlay and failed exactly on the missing contract. Minimal GREEN injects the existing overlay only for OpenCode PM; non-OpenCode Env remains unchanged. Targeted planning suite and full CURRENT `/src/swe-af/go` suite PASS; the exact same two-file delta was propagated stale-safely to the frozen repo, whose targeted and full Go suites also PASS, `git diff --check` PASS, and the framework fix was committed locally as `6f5b4382e6231721f60be7045b9d91fd85e34fb5` (`fix(planning): provide runtime FCM config to OpenCode PM`).
- Attempt 7: `exec_20260913_233822_sesklhiu` / run `run_20260913_233822_uddk9m55` started on frozen SWE `6f5b4382...` + AgentField c092 + planner SHA `3dc66b5d...` and again failed before model progress with the PM no-output diagnostic. Post-terminal root-cause tracing proved this run was **invalid runtime-materialization evidence**, not a new schema defect: the manually launched planner was running inside the AgentField workforce where `PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin`, `SWE_OPENCODE_BIN` is empty, `opencode` is not executable/discoverable (`/usr/bin/env opencode` -> exit 127), and filesystem search finds no installed OpenCode binary outside temporary test fakes. By contrast, canonical SWE `go/Dockerfile` installs OpenCode via `https://opencode.ai/install`, puts `/root/.opencode/bin` on PATH and ships the runtime OpenCode config; `go/agentfield-package.yaml` also defines `SWE_OPENCODE_BIN` for runners where discovery is unavailable. Therefore process-only planner execution in the AgentField workforce cannot serve as a valid OpenCode acceptance runtime.
- Attempt 8 preflight is now GREEN without creating a new container/service. Main AgentField topology exposed an existing native SWE package runtime in the same workforce: `/afhome/packages/swe-planner/bin/swe-planner` with `SWE_OPENCODE_BIN=/afhome/bin/opencode`; the wrapper pins OpenCode runtime `v1.17.15`, SWE config `/src/swe-af/opencode.json`, and broker model remapping. Direct non-mutating wrapper canary from the frozen repo returned `NATIVE_OPENCODE_CANARY_OK` with real `provider=fcm/model=fcm` turns and token telemetry. The invalid manual planner was stopped; the exact repaired planner `/tmp/swe-planner-fb0-6f5b438-c092` is now running on port `8801` with the native wrapper contract against isolated exact c092 control-plane, registered as `swe-planner-fb0-native`, instance `9e9b64c8a3de4730905a4786b8bdaada`, ready at `http://workforce:8801`. Frozen source remains `6f5b4382e6231721f60be7045b9d91fd85e34fb5`, AgentField remains `c0923acd...`, frozen Build config unchanged, streak 0/3. This closes `RUNTIME_MATERIALIZATION_GAP` for the canary path; Attempt 8 may START only after immediate clean-worktree/zero-live-mutator readback.

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

## AI Handoff / Bootstrap Protocol

This section is the mandatory bootstrap for any new AI/operator taking over the project. The handoff text or chat history is **not** authoritative by itself: first reconcile this `PLAN.md` against CURRENT runtime/readback and only then continue.

### Authority / reading order

Read only what is needed, in this order:
1. **This file first:** `n0namer/swe-af`, branch `dev`, `PLAN.md`. Extract Global North Star, CURRENT, active gate, DoD, anti-drift and Current Next Move.
2. **Local execution rules before mutation:** `/src/swe-af/AGENTS.md`; if touching another owner repo/workspace, read its nearest `AGENTS.md` and root `ERRORS.md` when present.
3. **Architecture only for decisions that need it:** `/src/swe-af/docs/ARCHITECTURE.md`. Do not let runtime accidents redefine architecture; runtime readback owns actual state, architecture/SoT owns intended design.
4. **BMAD entrypoint:** canonical `n0namer/BMAD-MNNZ/.agents/skills/bmad-help/SKILL.md`. BMAD-MNNZ is rules/skills only, never project SoT. From `bmad-help`, load exactly the specialized skill needed for the current gate (typically `bmad-testarch-test-design` for fault/risk design, `bmad-quick-dev` for a proven bounded source fix, then adversarial/edge review once).
5. **CURRENT runtime/readback:** verify live source identity, dirty state, exact loaded planner binary/PID, planner registration/`instance_id`, target workspace cleanliness, live OpenCode mutators, control-plane execution state, and exact tested artifacts. Never assume PIDs, `/tmp` artifacts, health, or execution state survived from this document.
6. **Supporting owner source only if the active blocker requires it:** current AgentField work has been container-local at `/tmp/agentfield-f06-fix`, based on deployed AgentField SHA `4d337c1ae5104418311fcba414a1c2f85c2abb89`. If that path is absent, reconstruct from the exact deployed/source identity; do not guess from latest upstream.
7. **External code-graph tools if present:** Codebase Index `/tmp/tools-codebase-index`, Graphify `/tmp/tools-graphify`, ArchSteer `/tmp/tools-archsteer`. Reuse them before installing alternatives. Codebase Index + Graphify are the primary Go fault-boundary discovery tools; ArchSteer is useful as an architecture/governance view but is not authoritative for Go coverage.

### North Star to preserve

Bring SWE/SWE-AF to a genuinely autonomous end-to-end engineering system that can take a real task, localize and edit the correct source, validate it canonically, repair bounded failures by itself, recover across interruption without duplicate effects, produce an exact deliverable, and pass an independent executable oracle. Production evidence is the **full `orch.Build` lifecycle**, not `implement_issue`, reviewer prose, HTTP health, or a commit existing. First milestone is a fresh **3/3 consecutive independently accepted full-Build streak**, then materialize the exact accepted source and replay on a clean runtime.

### Mandatory takeover reconciliation before any new mutation

The next AI must independently verify these facts instead of trusting this prose:
- canonical SoT is still this `dev/PLAN.md` and no newer explicit project decision supersedes it;
- `/src/swe-af` actual source/dirty state and frozen sacrificial baseline are what CURRENT says;
- current planner/process/control-plane identity is fresh and the target workspace has zero pre-existing mutators;
- F06 process-identity/restart fix evidence remains reproducible and no later runtime change invalidated it;
- `/tmp/agentfield-f06-fix` CURRENT HEAD and working-tree delta are inspected before deciding whether FB-0 can resume.

CURRENT handoff readback at the end of the takeover batch:
- `/tmp/agentfield-f06-fix` HEAD = exact clean local candidate `c0923acdfca043c2c07e3d34daaa09e2a7e41d38`; unsupported Python experiment and generated `.archsteer/` are absent;
- `/tmp/swe-af-fullbuild-current-20260912` is clean at exact baseline `cdc39105e92937e0085410a656346471b2dc6834` and its full Go suite PASS with the c092 SDK through a temporary external modfile;
- exact planner artifact is `/tmp/swe-planner-agentfield-c092-20260913`, SHA256 `c923a3653cceddbbc0672f6c6b8b673f3b819ef15751878e255014a9d0732e79`; the isolated c092 control-plane/planner fault runtime was stopped and its temporary DB/config removed after hard-crash + SIGTERM verification, so no old PID should be reused as truth.

The AgentField takeover is complete. The next operator must re-observe clean repo/mutator/runtime preconditions and then execute FB-0 attempt 5; do not reopen the takeover without new evidence from FB-0 or a changed topology/runtime.

### Takeover batch — COMPLETED DoD (historical checklist)

Use one bounded reconciliation batch:
1. inspect `git diff` in `/tmp/agentfield-f06-fix`; classify each delta by fault-family/boundary and discard only unsupported/accidental changes, preserving proven work;
2. rerun the discriminating tests for the graph-found identity/concurrency family: process `instance_id` propagation/stability, restart reap, stale old-instance status/heartbeat rejection, duplicate execution-ID cancel ownership, shutdown admission/cancellation including skills/async reasoners, plus any cross-language sibling test that actually has a RED oracle;
3. run targeted `-race` where applicable, full AgentField Go SDK suite, relevant control-plane handler suite, Python targeted suite if changed, and `git diff --check`; do not call missing dependencies an application failure;
4. re-index the final exact AgentField source with Codebase Index + Graphify and repeat the same fault-family search. Do not add duplicate tests where terminal-state/row-lock/PID/owner guards already have a discriminating oracle; add RED->GREEN only for a genuinely uncovered critical boundary;
5. run one BMAD adversarial/edge review on the final delta, triage findings once, then create a **local-only exact tested commit** in the owner clone. Do not publish/deploy merely to debug;
6. validate the exact AgentField candidate against frozen SWE via local module replace/full Go suite and, if the change affects restart/process containment, repeat the same real F06 fault injection before calling it VERIFIED;
7. write exact files/tests/commands/commit/runtime evidence back into this `PLAN.md`, refresh anti-drift, and only then choose the next gate from fresh evidence.

Takeover STOP CONDITION: **SATISFIED** at exact clean AgentField commit `c0923acdfca043c2c07e3d34daaa09e2a7e41d38`, with applicable P0 identity/concurrency/restart invariants GREEN, frozen SWE compatibility PASS, and real hard-crash + SIGTERM runtime proof complete. FB-0 remains intentionally a separate next batch.

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
The active product P0 remains **F16 fail-closed resume on corrupted checkpoint** in exact `/src/swe-af`. Fresh readback on 2026-09-14 reconfirmed the defect: `loadCheckpoint` maps both read errors and JSON decode errors to `nil`, while `RunDAG(..., WithResume(true))` then continues with fresh state and immediately saves a checkpoint. The previous operator blocker has moved. `DEV_VALIDATION_CAPABILITY_GAP` is now **closed**: `vps-terminal-dev` is loaded from canonical image source `9cf1f189e02df1827440494bc536c2deb23d3ea7`, `/v1/target-registry/action` is callable, runtime registry revision 7 for `agentfield-dev-workforce` includes `debug_clone`, and activation was proven by `prepareDebugClone` advancing past capability policy. The new exact blocker is **`DEV_WORKFORCE_SOURCE_ABSENT`**: `prepareDebugClone(agentfield-dev-workforce,/src)` now fails `source_target_unavailable: Expected one source target, found 0`, and CURRENT Docker inventory contains no workforce container for compose project `edshqtkwskg3lrczekhcmd71` although `workforce` remains declared in the canonical compose. The prior PLAN statement that the workforce merely exists as `Exited 31` is stale actual-state evidence.

THIS BATCH:
used BMAD test-design/trace + bounded quick-dev discipline, systematic-debugging, condition-based readback, and verification-before-completion. Reconciled `dev/PLAN.md`, `/src/swe-af/AGENTS.md`, live F16 code, current Docker/Coolify state, and the `vps-terminal` owner source/ERRORS ledger. Found that the DEV gateway was actually loading stale persistent hotfix source rather than its canonical image code; restored the already-canonical target-registry control path directly in the container, `node_check` passed, then restarted the same DEV app and independently proved the new gateway loads `/app/gateway/server-live.mjs` from the existing `9cf1f189...` image with on-disk=loaded identity. Runtime registry was stale-safely advanced from workforce revision 6 to 7 by adding only `debug_clone`; readback verified SHA `9a69919c2f5512190ef80807cbf24c6afe3ef55ef2ae2b6b0423b100951ed215`. A second same-app restart loaded that registry. The discriminating capability probe then changed from `debug_clone_not_allowed` to `source_target_unavailable`, proving activation and isolating the next layer. No SWE product source changed. Coolify owner-lifecycle recovery for the existing AgentField DEV app is queued; CURRENT must be re-read before any additional lifecycle action.

NORTH-STAR DELTA:
the operator-side debug-clone capability is now executable rather than hypothetical. We eliminated stale hotfix activation as the gate and reduced the remaining prerequisite to one concrete state fact: restore the already-declared DEV workforce source generation, then use the typed TTL clone to get the proven Go image/toolchain without weakening provenance or inventing another validator. F16 itself remains unchanged and has not been falsely marked RED/GREEN.

STOP CONDITION:
do not weaken `/opt/reconcile_dev_workspace.sh`, do not create a substitute validator/service, do not mutate SWE semantics before a real F16 RED, and do not repeat lifecycle calls without fresh post-state. First re-observe the already-queued AgentField DEV lifecycle action. If the declared workforce source appears, immediately call `prepareDebugClone(agentfield-dev-workforce,/src)`, execute the approved TTL clone, prove exact `/src/swe-af` + `/usr/local/go/bin/go`, and run the corrupt-checkpoint discriminator. If the owner lifecycle completes yet workforce remains absent, classify that exact compose/source-generation recovery boundary before changing operator code.

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

The sole next move is to close **`DEV_WORKFORCE_SOURCE_ABSENT`** through the existing AgentField DEV owner lifecycle, then immediately continue the already-proven typed validator path. `debug_clone` activation is complete; do not reopen it unless fresh evidence regresses to `debug_clone_not_allowed`. Re-observe the queued Coolify lifecycle action for `edshqtkwskg3lrczekhcmd71`; if its declared `workforce` container appears, call `prepareDebugClone(agentfield-dev-workforce,/src)` at once, execute one bounded TTL clone, prove `/usr/local/go/bin/go` and exact `/src/swe-af`, then run the F16 corrupt-checkpoint regression to RED. Only after RED may SWE checkpoint/resume semantics change. If the lifecycle action reaches terminal state and workforce is still absent, localize that compose/source-generation boundary before changing `vps-terminal` or SWE code. Do not weaken provenance, create a substitute validator/service, or return to GitHub-edit -> redeploy as the programming loop.
