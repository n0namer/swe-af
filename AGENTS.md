# SWE-AF Engineering Contract

Applies to the entire repository unless a future nested `AGENTS.md` adds stricter subtree rules.

## Repository role and authority

- Target repository: `n0namer/swe-af`.
- Active downstream development branch: `dev`.
- `fork/main` is intended to remain a clean upstream mirror. Do not put downstream project work on `main`.
- Project/design SoT: root `PLAN.md` for North Star, current phase, bounded batches, DoD, current runtime facts and known drift.
- Product architecture: `README.md` + `docs/ARCHITECTURE.md`.
- Canonical permanent-DEV operator lifecycle: `n0namer/universal-solver/docs/runbooks/agentfield-dev-debug-test-handoff.md`.
- Cross-repo engineering standard: `n0namer/server-ops/docs/standards/FAST_VERIFIED_ENGINEERING.md` on `main`.
- BMAD is the project-management method source only; read relevant skills from `n0namer/BMAD-MNNZ`. Do not create parallel BMAD artifacts when `PLAN.md` can own the state.

## North Star

Deliver SWE-AF as a reliable AgentField software-engineering capability that can take a bounded real-repository task, produce the correct change, pass canonical tests/acceptance checks, recover from expected failures, and produce a durable exact-SHA Git result with complete provenance.

For bounded well-scoped work prefer `swe-planner.implement_issue`. Use `swe-planner.build` only when feature-level decomposition is actually required and lower gates have passed.

## Required entry read

Before a material mutation:

1. Read this file.
2. Read root `PLAN.md` and resolve North Star -> Phase Goal -> current DoD -> next bounded move.
3. Read the relevant parts of `README.md` and `docs/ARCHITECTURE.md`.
4. Read root `ERRORS.md` if it exists in the active branch. It does not currently exist on `dev`; do not treat the `main`-only ledger as active `dev` state.
5. For permanent AgentField DEV work, read the Universal Solver handoff above.
6. For complex implementation/debug work, use `bmad-help` to route to the smallest fitting BMAD skill. For the current permanent-DEV acceptance workstream, `bmad-testarch-test-design` owns risk-based gate design; `bmad-quick-dev` is the implementation method once a source/config defect is proven.

## Fast Verified Engineering (FVE)

Optimize for **time-to-verified-running-change**, not time-to-patch or time-to-merge.

Use the smallest evidence-complete route:

```text
OBSERVE -> LOCALIZE -> ROUTE -> PATCH -> TARGETED VERIFY -> ITERATE -> FULL VERIFY -> RUNTIME PROOF -> CANONICALIZE -> DEPLOY -> POST-DEPLOY VERIFY -> WRITE-BACK
```

Do not mechanically execute irrelevant stages. The smallest route that closes the DoD with evidence wins.

## Evidence layers — KEEP SEPARATE

Always distinguish:

- Project/design SoT — what should exist.
- Source-on-disk — files in the intended repo/workspace and their exact identity/dirty state.
- Loaded runtime — process/container/image/config actually loaded.
- Concrete execution — a specific AgentField execution/request/job and its execution/correlation identity.
- Deterministic validation — exact checks/tests passed on intended source.
- Functional/semantic outcome — whether the system actually did the required engineering work.

Never substitute one for another:

- code-on-disk != loaded runtime;
- test PASS != deploy/load proof;
- HTTP 200 or healthy != functional acceptance;
- AgentField `execution succeeded` != semantic correctness.

## Coding lanes

### Runtime-first lane — preferred for CURRENT work

Use the already-authorized permanent DEV runtime when the defect depends on real provider/config, AgentField registration, process loading, network behaviour, integration or semantic E2E.

CURRENT known DEV topology is owned by `PLAN.md`; do not hard-code stale container generations here. Stable identifiers:

- Coolify application: `edshqtkwskg3lrczekhcmd71`.
- Persistent runtime source root: `/src`.
- SWE-AF live workspace: `/src/swe-af`.
- SourceLoop/runtime-capture bridges accepted runtime source deltas toward durable Git.

Runtime debug loop:

```text
CURRENT runtime/source
-> bounded stale-safe patch
-> smallest affected check
-> same-target reload/restart only if needed
-> functional canary
-> bounded logs/execution evidence
-> iterate
```

Rules:

- Preserve unrelated dirty state.
- Prefer typed/scoped live-patch capabilities over opaque broad mutation.
- Restart only the same owning process/target when code loading requires it.
- Permanent DEV may temporarily be ahead of Git during active debugging; accepted changes may not remain container-only.
- Production live editing is forbidden by default.

### Repository-first lane

Use an exact-SHA isolated Coding Station workspace for source-bound logic, multi-file changes, refactors, dependency/build work, or work that does not require real runtime state for each iteration.

Canonicalize the **same exact verified delta**. Do not re-implement a proven runtime fix from memory in a second workspace.

## Canonicalization and release

GitHub/CI/redeploy are canonicalization/release boundaries, not the inner debug loop for every hypothesis.

For permanent DEV:

- Lane A = live debugging and fast hypothesis -> evidence.
- Lane B = accepted delta -> SourceLoop/capture -> durable `fork/dev` -> exact SHA -> materialization.

Debug rule: runtime may be ahead of Git.
Handoff/release rule: runtime may not remain ahead of Git.

Generated files, caches, logs, test artifacts and temporary instrumentation are not source write-back.

## Validation ladder

Use progressively more expensive evidence:

1. Parse/syntax/static check.
2. Directly affected tests.
3. Related module/component regression.
4. Full required canonical suite.
5. Runtime smoke/integration.
6. Semantic/business/E2E acceptance when required by DoD.

Do not skip final full/runtime/semantic gates merely because an early targeted test passed.

## Project-specific build and test commands

Do not invent commands. The following are verified from the repository.

### Python / repo root

From repo root:

```bash
make test      # python -m pytest tests/ -x -q
make check     # test + python -m compileall -q swe_af/
```

### Go / installed AgentField node path

From `go/`:

```bash
make build     # go build ./...
make vet       # go vet ./...
make test      # go test ./...
make check     # vet + test
```

Canonical Go check: `cd go && make check`.

The installed AgentField node path is the Go port. For permanent-DEV node/bootstrap/provider changes, validate Go first.

The Go module depends on the AgentField Go SDK through a relative replace/dev workspace. If the required sibling SDK checkout/runner is missing, classify `VALIDATION_BLOCKER`; do not call the product broken merely because the validator environment is incomplete.

## Runtime proof and observability

For current permanent DEV acceptance, correlate at minimum:

- AgentField node id and `active/ready` status;
- execution_id for reasoner/pipeline runs;
- execution-scoped logs when available;
- bounded workforce/node logs only after narrower evidence;
- target process/container generation and source/WORKING_DEV_SHA when applicable;
- functional diff/tests for `swe-planner.implement_issue`.

Diagnose narrow-to-broad:

1. Structured AgentField execution/status evidence.
2. Execution-scoped logs/notes/debug details.
3. Bounded runtime logs for exact node/container.
4. Broader service/node logs only if needed.

The current workforce HTTP healthcheck proves only the provenance helper service. It does **not** prove `swe-planner` readiness.

## Current semantic acceptance order

Follow the current `PLAN.md`, but the present gate order is:

1. Provider/bootstrap gate.
2. `swe-planner` becomes `active/ready`.
3. Non-mutating reasoner/schema smoke with execution evidence.
4. ONE bounded `swe-planner.implement_issue` canary on an existing safe workspace.
5. Affected + canonical repo tests on the result.
6. Failure/recovery canaries.
7. SourceLoop durability.
8. Exact-SHA materialization and repeated semantic canary.

Do not use full `build` as a substitute for these lower gates.

## Change safety and recovery

- Preserve unrelated state.
- Keep a rollback/recovery path for operational changes.
- Before destructive/loss-inducing actions, preview exactly what will and will not change.
- After timeout, ambiguous mutation or tool error, read post-state before retrying.
- Retry an identical failed mutation at most once unless new evidence changes the diagnosis.
- Never expose provider keys, tokens or credentials in logs, docs or reports.
- Production live editing is forbidden by default.

## Write-back owners

- Project state, current gates, DoD, evidence and debt: `PLAN.md`.
- Product source: accepted exact delta on `dev`, preferably the same delta proven in permanent DEV and then captured/canonicalized stale-safely.
- General AgentField DEV lifecycle/fleet materialization: Universal Solver runbook/fleet owners.
- Server/operator topology and cross-repo decisions: `n0namer/server-ops`.
- Generated files, caches, logs and test artifacts are not source write-back.

## DONE contract

Final status is exactly one of:

```text
DONE | PARTIAL | BLOCKED | FAILED | EVIDENCE_MISSING
```

`DONE` requires every task-specific DoD criterion to have the required evidence. A merge, test PASS, healthy container or AgentField success acknowledgement is not sufficient when the DoD requires runtime or semantic proof.

Every handoff must distinguish: what changed; what was tested; which runtime executed; what functional/semantic evidence passed; what exact source/runtime identity was proven; and what remains.
