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

- The recent exact temporary `swe-planner` remains live as PID `366110`, listener/health on port `8005`.
- Recent control-plane readback for the planner reports 25 executions in 24h: 18 succeeded, 7 failed.
- Latest L3-26 coder executions fail before producing a deliverable:
  - `Schema validation failed ... output file was NOT created`;
  - latest attempts also preserve upstream `Unexpected server error`.
- L3-26 created no accepted commit/files. Do not blindly repeat the same trajectory without new route evidence.

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
- Because L3-26 failed after L3-24/L3-25, those historical positives do not by themselves satisfy the fresh production streak below.

### Supporting systems / non-P0 debt

- Coding Station functional API remains unhealthy from the operator surface (`stationHealth` / `stationReady` Gateway Timeout) despite one healthy API + one healthy runtime container. This is not the current critical path because the workforce Go validator works.
- AgentField operator gateway is intermittently `Bad Gateway`; direct authenticated control-plane readback from the workforce works.
- Control-plane has shown stale executions labelled `active`; execution-state reconciliation is a production reliability debt, but do not stop the accepted-task ladder unless it blocks safe recovery.
- SourceLoop/durable Git identity does not yet represent the entire tested live product delta. Canonicalization is required after correctness baseline, not before.
- Python canonical pytest is not provisioned in the current workforce. Do not install it merely to manufacture a green historical task; prefer tasks with already-provisioned canonical validators for the next acceptance streak.
- Process/zombie/resource hygiene is an operational risk; escalate to the critical path only if fresh task evidence shows spawn/resource pressure.

## Production-Ready Milestone

First production-readiness gate requires a fresh controlled streak of **3 consecutive independently accepted engineering tasks** on a frozen/currently proven stack:

- at least one real multi-file task;
- at least one interruption/partial-result recovery task;
- zero duplicated non-idempotent effects;
- exact tested identity = delivered identity for every task;
- canonical deterministic validation where available;
- independent executable oracle/verifier;
- task/provider/model/trajectory recorded;
- total model cost and wall time recorded for every task.

After 3/3, materialize the accepted source into canonical Git/SourceLoop, create/rebuild a clean runtime from that exact SHA, and replay the controlled ladder. Only the clean-runtime replay can close the first production-ready milestone.

## Critical Path / Pareto Order

### Gate PR-1 — recover one stable execution route

Status: BLOCKED / P0 — fresh PR-1 experiment localized the failure before the FCM broker request.

Why first: CURRENT product Go tests pass; another product-code patch is not justified until a fresh execution experiment proves a source defect.

2026-09-12 PR-1 evidence: `exec_20260912_093427_ti492y5g` ran a different real Go recovery task on clean base `2c374989b39d0b53b34ef33fd2ba6289e74194ae`. Live process evidence proved OpenCode argv `-m fcm/fcm` in isolated worktree `.worktrees/1505a247-checkpoint-completed-issue-before-cancel`. The run terminated in ~30s with `failed_unrecoverable`, `commits=[]`, `files_changed=[]`, upstream-facing `Unexpected server error`, and missing structured output. The temporary worktree/branch was cleaned and no mutation effect remained. Critically, FCM `/health` showed `requestsRouted=697` both immediately before and after the failed run, and runtime telemetry gained no new call. Therefore this reproduction did not reach the FCM broker: the current critical boundary is OpenCode/provider-adapter/config/runtime handling **before broker request**, not SWE product logic and not FCM inference. Stop repeating full `implement_issue` until that pre-broker seam is diagnosed.

DoD:
1. Freeze CURRENT SWE source/contracts/runtime guards.
2. Vary one execution-route factor only.
3. Prefer a route already evidenced to complete tool trajectories (known stronger route or another FCM route that passes the same tool canary).
4. Run one smallest discriminating real-task/canonical-validator experiment.
5. Capture exact execution ID, provider/model route, actual diff/commit or no-effect state, error boundary, wall time and cost if available.
6. Stop on the first evidence-backed failed gate; do not stack prompt/schema/model/enforcement changes.

PASS:
- a bounded deliverable reaches canonical validation + independent oracle.

FAIL/BLOCKED:
- one fresh experiment localizes a single boundary with evidence and leaves no ambiguous/repeated mutation.

### Gate PR-2 — accepted task 1/3

Status: PENDING PR-1

Use a small real Go task when possible because the canonical Go runner is already available.

DoD:
- actual source edit in the task workspace;
- bounded exact commit/diff;
- affected canonical Go test PASS;
- relevant regression/full package test PASS;
- independent oracle PASS;
- no unrelated file/effect;
- exact tested=delivered identity;
- cost/latency/continuations recorded.

### Gate PR-3 — accepted task 2/3: recovery/no-duplicate-effect

Status: PENDING PR-2

DoD:
- task performs an observable effect;
- interruption/timeout or partial-result point is exercised only where post-state/effect identity is readable;
- continuation observes already-satisfied effects and does not repeat them;
- UNKNOWN state fails closed/read-only;
- canonical tests + independent oracle PASS;
- duplicate non-idempotent effects = 0.

### Gate PR-4 — accepted task 3/3: multi-file

Status: PENDING PR-3

DoD:
- real multi-file engineering change;
- exact bounded diff;
- canonical affected + regression validation PASS;
- independent oracle PASS;
- exact tested=delivered identity;
- cost/latency recorded.

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
no freshly proven stable execution lane for the next accepted task; L3-26 cheap trajectory produces no deliverable while CURRENT product Go suite passes.

THIS BATCH:
change no product code until route evidence proves a source defect. Prove one stable route with one discriminating real-task experiment.

NORTH-STAR DELTA:
establish the lane for fresh task 1/3, or localize one single critical-path boundary failure.

STOP CONDITION:
either one bounded deliverable passes canonical validation + independent oracle, or one fresh run identifies one enidence-backed blocker. Stop after that blocker and replan.

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

Run **one** PR-1 discriminating execution-route experiment on the frozen CURRENT source. No product-code mutation unless that experiment proves a source defect. If it produces a deliverable, immediately validate it canonically and independently and count it as fresh task 1/3; if it fails, stop at the first evidenced boundary and make only the smallest root-cause fix in the container.
