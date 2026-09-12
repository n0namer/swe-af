# SWE-AF Project Plan (Source of Truth)

Status: active
Last reconciled: 2026-09-12
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

## CURRENT — 2026-09-12

### Source/runtime identity

- CURRENT live product source: `/src/swe-af`.
- Git baseline HEAD: `58c4e0d19081bc52363c120b7963a34cebb1e894`.
- Working tree is intentionally ahead of baseline: fresh readback shows 58 modified/staged paths + 2 untracked paths.
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
- Current exact tested planner binary: `/tmp/swe-planner-fullbuild-route-20260912`, SHA256 `5d4404bbecea63a21ff440c9c814ec3907980f543e8cddac15ada3635645608f`; current planner PID `483648`, callback/health on port `8005`, control-plane active.
- The binary includes the proven runtime-owned FCM overlay plus a full-Build planning parity fix: when runtime/provider resolves to OpenCode, Product Manager now uses the harness/OpenCode path instead of AgentField direct AI. Deterministic RED `TestProductManagerOpenCodeUsesHarnessWhenDirectAIIsAvailable` reproduced the old direct-AI misroute for `fcm/fcm`; targeted planning tests and full `go test ./... -count=1` PASS after the fix.
- Python reference `run_product_manager` always uses the configured harness provider/model; the Go direct-AI shortcut was therefore a port/runtime-contract regression, not an FCM outage.
- Post-fix issue-level task `exec_20260912_094308_za4vz5sv` reached FCM repeatedly, edited the target repo, and completed with a bounded two-file commit. The issue execution route is recovered; full-Build acceptance remains separate and is currently 0/3.

### FCM

- FCM is not globally down.
- Fresh `/health`: version `0.5.81`, 34 active models, 8 effective providers, no in-flight request at readback.
- Small current direct telemetry for `routerai/z-ai/glm-5.3-flash`: 6/6 successful calls.
- Therefore the L3-26 failure is narrower than “FCM is broken”: current evidence points to the OpenCode/FCM tool trajectory, structured-output/provider boundary, or another execution-route integration seam.
- FCM decides model/provider routing; it does not own authoritative engineering obligation state.

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
- Current clean precondition: sacrificial repo `/tmp/swe-af-fullbuild-current-20260912` is clean on baseline branch `fullbuild-current-baseline-20260912`, exact SHA `e03ee19d1940a29318ebd9f820c92be7fe4931f6`; no live OpenCode process targets that workspace; planner `/tmp/swe-planner-fullbuild-route-20260912` is healthy and registered.

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
production acceptance has not yet exercised the top-level autonomous `orch.Build` lifecycle. Component-level `implement_issue` evidence is useful but insufficient; MICRO-2 proved a repairable verifier failure can terminate there without self-repair.

THIS BATCH:
materialize a frozen local baseline from CURRENT tested `/src/swe-af`, then run exactly one `swe-planner.build` canary with the production recovery controls enabled (replanning, issue advisor/retries, integration testing, verifier-fix cycles, deterministic Git). Keep the proven `fcm/fcm` route. Disable only external GitHub PR/CI side effects for this first canary because they are outside the already-authorized local scope; do not claim CI coverage.

NORTH-STAR DELTA:
replace micro-harness confidence with direct evidence about whether the full SWE controller can autonomously plan, implement, validate, repair, verify and finalize one real SWE-AF change without operator task-code edits.

STOP CONDITION:
either the full Build reaches terminal success and then passes independent canonical validation/oracle, or one run identifies the first evidence-backed missing contract in the full lifecycle. At that point stop, write back, and repair only that proven framework gap before retrying the same canary.

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

Run **the same FB-0 full `swe-planner.build` canary for attempt 3** from the clean frozen baseline `e03ee19d1940a29318ebd9f820c92be7fe4931f6` using the tested planner `/tmp/swe-planner-fullbuild-route-20260912` and `fcm/fcm`. Before START require zero live OpenCode process targeting the sacrificial repo. After START the operator makes zero task-code edits. Follow actual process/artifact/effect state rather than stale control-plane labels. Stop after either autonomous terminal success followed by canonical Go validation + independent oracle, or the first new evidence-backed full-lifecycle blocker; write that blocker back before any framework repair.
