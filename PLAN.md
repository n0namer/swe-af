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

### 2026-09-11 L3 semantic acceptance checkpoint

Current gate: **FCM-only full L3 acceptance on `qa-synthesizer-fcm-l3-19`**.

Fresh runtime readback shows live source identity remains `58c4e0d19081bc52363c120b7963a34cebb1e894` plus a large intentional dirty delta. The l3-19 issue worktree is terminal enough to contain bounded commit `53c2cca` on base `2c37498`, touching only `swe_af/reasoners/execution_agents.py` and `tests/test_qa_synthesizer_direct_schema.py`; its only remaining untracked product-tree item is `.agentfield-out-1432774580/` containing the coder completion envelope.

Independent operator acceptance rejects `53c2cca` as final L3 evidence. Its new tests now *do* invoke the production function `run_qa_synthesizer`, closing the l3-18 vacuous-test defect, but the committed implementation still uses `hasattr(...)`/duck typing instead of strict `QASynthesisResult` validation. By contrast, the current live `/src/swe-af` delta already contains the minimal stricter repair: `isinstance(result, QASynthesisResult)` plus an explicit error note before fail-safe fallback, and a production-path regression test that exercises both a real BLOCK result and a non-schema object.

Validation status is **PARTIAL / VALIDATION_BLOCKER**: system `python3` is available, but `pytest` is not installed and there is no repo `.venv`; installing test dependencies is explicitly forbidden for this acceptance lane, so no PASS may be claimed from this container alone. Previous prompt/role hardening and `make check` evidence remain historical evidence only until the exact current live delta is revalidated in an environment that already has the canonical test runner.

Closed DoD so far: bounded candidate commit exists; regression now exercises the changed production entrypoint; live repair is strict-schema rather than duck-typed. Still open: prove actual OpenCode argv `-m fcm/fcm` for l3-19 from preserved runtime evidence; prove reviewer did not install dependencies and correctly graded test adequacy; prove repair/verifier/delivery behavior; remove/archive `.agentfield-out-*` outside the product worktree; run canonical tests without installing missing tooling; verify tested identity equals delivered identity.

One next move: recover the l3-19 reviewer/verifier/wrapper evidence and exact process/model provenance from preserved runtime artifacts; then validate the existing strict live repair with the smallest already-provisioned canonical runner. Do not add broader prompt rules or install dependencies merely to manufacture a green test.

### 2026-09-11 L3-20 checkpoint

Fresh same-case run `qa-synthesizer-fcm-l3-20` was executed from a planner rebuilt from exact live source. Actual process readback closed the FCM routing gate: both `/afhome/bin/opencode ... -m fcm/fcm` and the real OpenCode child argv contained `-m fcm/fcm`. The coder produced bounded commit `9e2877f` touching only `swe_af/reasoners/execution_agents.py` and `tests/test_qa_synthesizer_direct_schema.py`; it did not run dependency-install commands and did not manufacture a green test result when pytest was unavailable (`tests_passed` remained non-PASS).

Independent operator inspection found the candidate tests still invalid: they passed a nonexistent `router=` keyword to `run_qa_synthesizer`, and one fallback assertion contradicted the real fail-safe summary. Hardened reviewer independently reproduced these defects, used no dependency installs, and correctly returned `approved=false, blocking=true`. Therefore the previous reviewer semantic fail-open is not reproduced on l3-20.

The run did not reach repair/verifier/delivery because the temporary process-only planner disappeared after the blocking review; no live planner PID or listener on `127.0.0.1:8005` remained. Classify this as a current `AGENTFIELD_RUNTIME` / operator-session lifecycle blocker until a concrete process-exit cause is proven; do not relabel it as FCM/provider failure. The Go reviewer prompt already contains the stop-after-first-reproducible-blocker rule, so no broader duplicate reviewer prompt rule is justified from this run.

Current DoD: FCM argv PASS; coder no-install behavior PASS; reviewer fail-closed semantic grading PASS; repair/verifier/delivery still OPEN. Fresh source readback supersedes the earlier coder-prompt hypothesis: current Python/Go/golden coder prompts already forbid installing missing validation dependencies, so no additional coder prompt patch is justified. The reviewer runtime prompt also already has an explicit stop-after-first-reproducible-blocker rule; l3-20 therefore shows model overrun, not a missing reviewer rule. One next move: make the same-case planner lifetime independent of the managed terminal session using the smallest existing target-native process route, then rerun the unchanged issue through repair -> verifier -> delivery. Continue to require exact worktree/commit/operator evidence before L3 PASS.

### 2026-09-11 L3-21 checkpoint

`qa-synthesizer-fcm-l3-21` closed the prior planner-session hypothesis: a planner rebuilt from exact live source was launched detached with PPID 1 and health 200, remained alive through coder -> reviewer -> repair, and again exposed actual OpenCode argv with `-m fcm/fcm`. The candidate worktree is currently clean on bounded commit `77ae2fc` over base `2c37498`, touching the intended production/test surface only.

The candidate is still NOT accepted. Coder corrected the l3-20 unsupported `router=` test pattern by monkeypatching the production module and calling `run_qa_synthesizer`, but implementation/tests still treated an arbitrary wrapper object carrying `.parsed=QASynthesisResult` as happy-path compatibility. That conflicts with this L3 gate's strict direct-schema contract: an unexpected non-schema object must fail safely rather than be promoted by duck/wrapper typing.

Reviewer did block a separate reproducible test-adequacy defect and the pipeline entered repair, proving repair-path reachability. However reviewer also produced invalid acceptance evidence by executing pytest/Python from an older sibling l3-18 worktree against l3-21, and incorrectly characterized the `.parsed` wrapper fallback as acceptable compatibility. During repair, coder then executed dependency installation (`pip install pytest`, `pip install pydantic`) despite the explicit no-install contract; OpenCode runtime permissions allowed those commands. The invalid repair was stopped and no verifier/delivery acceptance may be inferred from this run.

Root-cause split is now explicit: prompt/model adherence and runtime enforcement are separate gates. A textual no-install rule is insufficient by itself. Critical acceptance invariants that invalidate evidence when broken (dependency installation during validation, cross-worktree test-runner provenance, direct-upstream model override, writes outside the bounded worktree) must be independently observable and, where the existing OpenCode permission layer supports it, fail-closed at runtime. This is configuration of the existing enforcement mechanism, not a new sandbox.

30-minute Pareto batch **COMPLETE**. A single shared owner now exists at `harnessx.OpenCodeNoInstallPermissionOverlay`; coder, reviewer, and verifier all receive the same OpenCode policy. Deterministic role tests assert that shared value rather than duplicating JSON. Exact live validation PASS: Python prompt compile PASS; targeted `go test ./internal/harnessx ./internal/prompts/coding ./internal/roles/coding ./internal/roles/advisor -count=1` PASS; `go vet ./...` PASS; `go test ./...` PASS.

Runtime enforcement was then proven against the real OpenCode/FCM path, not inferred from unit tests. Canary `pip install --help .` matched `pip install *` and was denied before execution. A second canary using the absolute l3-18 sibling-worktree `.venv/bin/python --version` matched `*/.worktrees/*/.venv/bin/*` and was denied. A normal `python --version` Bash call matched the catch-all allow rule and executed (then failed only because system `python` is absent), proving the overlay is selective rather than a blanket Bash shutdown. The same canaries logged provider/model `fcm/fcm`.

Batch DoD: **PASS** for runtime policy wiring and live enforcement. This closes the observed package-install/cross-worktree-runner escape route; it does not prove model adherence or full L3 correctness. L3 remains OPEN until a fresh same-case trajectory reaches repair -> verifier -> delivery and independent operator acceptance on the delivered SHA.

Next 30-minute Pareto batch: run a tiny deterministic policy micro-probe before another full issue. Measure separately: attempted forbidden action, executed forbidden action, allowed normal test action, tool-call count to stop, and final `VALIDATION_BLOCKER`/status. Use current full role prompt first; only if attempts remain high, compare a shortened hard-invariants-first prompt on the same FCM route. Do not change model, prompt, and enforcement simultaneously. This follows BMAD TEA test-design isolation and trajectory-eval practice: one factor per comparison, deterministic policy oracle first, model judge only where no programmatic check exists.

### 2026-09-12 handoff recovery / L3-25 independent acceptance checkpoint

Global project objective remains working SWE/SWE-AF on real engineering tasks with independent executable acceptance; AgentField/FCM/Coding Station remain supporting mechanisms only when they are on that critical path.

CURRENT source/runtime readback: `/src/swe-af` is still detached at `58c4e0d19081bc52363c120b7963a34cebb1e894` with a large intentional dirty delta. The live `swe-planner` is currently `active/ready`, PID 366110, callback/listener on `:8005`, and `GET /health` returns `{"status":"ok"}`. Control-plane readback reports 29 SWE-planner executions in the last 24h: 20 succeeded, 8 failed, 1 still running.

The cheapest L3-26 path is NOT the immediate North-Star blocker. Its recent `run_coder` executions fail before producing a structured coder artifact: repeated `Schema validation failed ... output file was NOT created`, and later attempts also include `provider error: Unexpected server error`. No code commit/file result exists for L3-26. Do not retry that path blindly; treat it as a later cost/reliability optimization after accepted-task correctness is established.

The shorter North-Star path was preserved L3-25 work. Branch `issue/58417244-dag-unknown-dependency-fcm-smart-l3-25` still points to exact commit `00517ce677e96b17bcd462bd46b9e5f9fa620674` on base `2c374989b39d0b53b34ef33fd2ba6289e74194ae`, changing only `swe_af/execution/dag_utils.py` and `tests/test_dag_utils.py`; the branch/worktree diff is clean and `git diff --check` passes. Existing coder/reviewer evidence was not repeated.

Independent operator verification was run against that exact SHA in a temporary detached validation worktree and then cleaned up. Canonical `python3 -m pytest --version` remains unavailable (`No module named pytest`), and no dependency was installed. Instead, a deterministic direct oracle exercised six acceptance/adversarial cases against the exact candidate `recompute_levels` implementation and passed 6/6; a minimal runner then executed all four committed regression test functions from `tests/test_dag_utils.py` and passed 4/4. `python3 -m compileall -q swe_af/execution/dag_utils.py tests/test_dag_utils.py` also passed. The tested branch/commit was preserved unchanged after cleanup.

Result: L3-25 is now independently executable-accepted for its bounded pure-function task, with an explicit canonical-pytest environment limitation. This is one accepted controlled engineering task, not a broad SWE L3 PASS and not yet the three-consecutive-task milestone.

Next bounded move: choose one additional real task whose canonical validator is already provisioned in the CURRENT runtime (prefer a small Go task if it has a real existing regression seam), and require exact diff + canonical test + independent oracle on the delivered SHA. Do not spend the next batch debugging the cheap L3-26 provider/schema path unless that becomes the blocker to accepted-task progression.

## Bounded Development Batches

Default batch size: about 30 minutes. Each batch closes a coherent DoD gate and writes back this file. Prefer the smallest 20% of work that removes the next 80% blocker.

**Current pointer (2026-09-11):** the L3-19/L3-20/L3-21 checkpoints and the next policy micro-probe above are the active execution plan. The legacy provider/bootstrap batches below are retained as historical evidence only; their old `IN PROGRESS` labels do not override the current FCM-only L3 gate or runtime readback.

### Batch 1 — Provider/bootstrap gate and SWE liveness

Status: IN PROGRESS

DoD:
1. Reconcile CURRENT topology/provider contract. **PASS**: legacy standalone SWE is removed; permanent DEV `edsh...` is current; isolated B is comparison-only.
2. Prove why CURRENT permanent DEV has no live `swe-planner`. **PASS**: canonical bootstrap evidence shows Anthropic/OpenRouter resolved empty and SWE was intentionally skipped.
3. Determine whether already-configured Gonka/OpenAI-compatible credentials can satisfy SWE without adding a new external secret. **PARTIAL**: repository code has `codex`/OpenAI auth; CURRENT Universal Solver compose already maps Gonka to workforce `OPENAI_API_KEY` + `OPENAI_BASE_URL`. Comparison with the proven Deep Research fix shows the required pattern is first-class OpenAI-compatible config plus preservation of the configured API base through dynamic/runtime overrides. Exact Gonka endpoint/model/tool-call compatibility inside the installed SWE harness still requires one live runtime canary.
4. Choose the smallest safe enablement path. **DECIDED / TWO-OWNER DELTA**: (a) SWE product contract must admit the existing OpenAI-compatible lane instead of declaring only Anthropic/OpenRouter; (b) Universal Solver runtime wiring must preserve that lane end-to-end, including the AgentField Go SDK's `AI_BASE_URL`/model contract for direct-AI calls, while harness/codex continues to receive `OPENAI_API_KEY` + `OPENAI_BASE_URL`. The bootstrap admission condition must accept this lane. Do not add/rewire Anthropic/OpenRouter secrets unless this canary disproves the existing lane.
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
