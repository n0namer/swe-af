# P0 F16 Corrupt-Checkpoint Resume Evidence

- **Task key:** `swe-af:v2-f16-corrupt-checkpoint-owner-evidence-20261011-r1`
- **Project:** `swe-af`
- **Attempt:** `attempt-6a5b59ff-e540-416c-8c21-9cfbb4a5e2f1`
- **Crew session:** `manual-20261011-sevenprojects-sweaf-p0-r1-v2-f1-90b418c3ba0969c3`
- **Assessment date:** 2026-10-11
- **Outcome:** Candidate evidence only; not accepted and not a source-code fix.

## Executive finding

**FACT — one real P0 regression is uncovered:** resume is not fail-closed against a stale or transplanted checkpoint whose source/repository identity differs from the live resume target.

`go/internal/schemas/execution.go` defines `DAGState` fields for `RepoPath`, `GitInitialCommit`, `BuildID`, `ArtifactsDir`, and progress state. However, `go/internal/dag/executor.go` initializes a fresh state and, when `resume=true`, replaces it wholesale with `loadCheckpoint(artifactsDir)`; there is no comparison of the loaded checkpoint to the current repository path, requested build identity, or current source/Git identity before issue execution proceeds. `go/internal/orch/resume.go` also reads the checkpoint as an untyped map and reconstructs only the issue graph/summary, without validating identity before forwarding `resume=true` to `execute`.

This means a validly shaped but stale/corrupted checkpoint can be accepted for another repository or source generation. The resulting run may continue from old issue progress against the wrong live source instead of refusing to resume. That violates the requested P0 fail-closed resume contract.

## Current state versus stale plans

### Current checkout evidence

- The checkout contains the Go implementation and tests under `go/`, including `go/internal/dag/checkpoint.go`, `go/internal/dag/executor.go`, `go/internal/orch/resume.go`, and `go/internal/orch/resume_test.go`.
- The repository is not a Node/TypeScript package in this checkout: the file inventory has no `package.json` and no JavaScript/Node test files. The README identifies the primary public runtime as Python and separately documents the Go node under `go/`.
- `_bmad-output/` is absent from the available project-file inventory; no current BMad generated acceptance artifact was available to corroborate the P0 claim.
- The allowed MCP surface provides file reads/writes and direct Node test invocation only. It does not expose Git status, branch, diff, HEAD, container state, or AgentField runtime readback. Therefore the live branch/diff, existing dirty files, current source commit, and runtime/container health are **UNKNOWN**, not inferred from this report.

### Plan evidence and limitations

- `PLAN.md` says it is the source of truth but was last reconciled on 2026-09-15. It records a historical live runtime identity `/src/swe-af@5b54a1dd4c579fd41ddf141e20dd46259519322e` and requires exact-source/preflight identity checks before a build counts.
- That historical plan entry is not proof of the current checkout or current AgentField runtime identity on 2026-10-11. No live readback was available in this bounded attempt.
- The report therefore treats the source/test code read from this checkout as current checkout evidence, the 2026-09-15 plan identity as historical context, and Git/BMad/runtime acceptance as unverified.

## Exact regression path

1. `go/internal/schemas/execution.go` serializes checkpoint state with `RepoPath`, `GitInitialCommit`, `BuildID`, and execution progress fields.
2. `go/internal/dag/checkpoint.go` loads any JSON that unmarshals into `schemas.DAGState`; it returns `nil` only for unreadable or syntactically/schema-invalid JSON. It performs no source/repository identity check.
3. `go/internal/dag/executor.go` enters the resume branch, calls `loadCheckpoint`, and assigns `dagState = loaded`. The newly initialized state containing the current invocation's `repoPath`, `gitConfig`, `buildID`, and workspace manifest is discarded without a compatibility check.
4. `go/internal/orch/resume.go` accepts the checkpoint if the file exists and JSON parses, then forwards reconstructed plan data and `resume=true` to `execute`; it does not reject missing/mismatched identity fields.
5. The DAG can therefore continue issue execution from checkpoint progress that was produced for a different source identity.

## Existing test coverage and gap

- `go/internal/orch/resume_test.go` verifies missing-checkpoint failure, plan reconstruction, raw envelope passthrough, and the default artifacts directory. It does **not** test repository-path mismatch, build-ID mismatch, Git/source-commit mismatch, or missing identity fields.
- `go/internal/dag/checkpoint_normalize_test.go` verifies that checkpoint list fields serialize as `[]` rather than `null` and that a normalized checkpoint round-trips. It does **not** verify fail-closed identity validation.
- The current tests prove shape compatibility, not P0 source-identity safety. This is the uncovered negative regression.

## Pareto-prioritized next coding/test task

**Task:** add a single fail-closed checkpoint validation seam before `RunDAG` accepts a loaded checkpoint, then cover it with mismatch and valid-resume tests.

**Likely source files:**

- `go/internal/dag/checkpoint.go`: add a focused validator that compares the loaded checkpoint to the current invocation's canonical repository/source identity and requested build identity, reusing existing Git identity helpers rather than creating a second Git mechanism.
- `go/internal/dag/executor.go`: invoke the validator before assigning `dagState = loaded`; return a deterministic error and do not execute issues on mismatch or missing required identity.
- `go/internal/orch/resume.go`: if this handler remains the public entry point, reject malformed/missing identity before forwarding to `execute`, or ensure the DAG validator is the sole authoritative gate and add handler-level evidence that the error is propagated unchanged.
- `go/internal/dag/*_test.go` and `go/internal/orch/resume_test.go`: add negative and positive regression tests.

**File-level Definition of Done:**

1. A checkpoint with invalid JSON, absent required source identity, repository-path mismatch, build-ID mismatch, or source/Git identity mismatch returns a deterministic error before any issue/coder call is made.
2. A checkpoint with matching canonical repository path, source/Git identity, and build identity resumes and preserves existing progress exactly.
3. Tests prove both the fail-closed negative path and the valid-resume path, including a guard that the coder/call seam is not invoked after rejection.
4. The validator handles path canonicalization consistently and documents the policy for legacy checkpoints; no silent fallback to trusting an unbound checkpoint is allowed.
5. Focused Go tests for DAG/checkpoint and orchestrator resume pass, followed by the relevant package regression tests; independent review must verify the source identity is obtained from the live target rather than copied from the checkpoint.

## Validation performed in this attempt

- **Read-only source inspection:** PASS — inspected `PLAN.md`, `README.md`, `AGENTS.md`, `go/internal/dag/checkpoint.go`, `go/internal/dag/checkpoint_normalize_test.go`, `go/internal/dag/executor.go`, `go/internal/orch/execute.go`, `go/internal/orch/resume.go`, `go/internal/orch/resume_test.go`, and the `DAGState` definition in `go/internal/schemas/execution.go`.
- **Node test execution:** NOT AVAILABLE — no JavaScript test file or `package.json` exists in the available checkout inventory; invoking Node against the Go/Python tests would not be a meaningful project check. No test result is claimed.
- **Go test execution:** NOT RUN — the task authorizes the scoped Node test runner only, and the available MCP surface does not provide a Go test runner.
- **Git branch/diff/status:** UNKNOWN — no Git inspection tool is available in the authorized scoped interface; no stale conversation state was treated as current.
- **AgentField/container/runtime readback:** UNKNOWN/BLOCKED — runtime/container mutation or inspection is outside the authorized interface; no source-loaded or deployment proof is claimed.

## Scope and safety

Only this report was created in the authorized write scope:

- `docs/go-wave-v2-f16-corrupt-checkpoint-evidence-20261011.md`

No runtime/container, Git, task, lease, lock, deployment, outbound message, or foreign/dirty file was modified. Candidate evidence remains **UNACCEPTED** pending independent source/runtime verification and review.
