# Native DAG cleanup barrier — 2026-10-13

## Decision

**Candidate result: BLOCKED on the authorized Go test runner.** The source review finds the requested cleanup-await mechanism already implemented in `RunDAG`; no executor change is justified from the available evidence. The focused regression file was added, but the permitted runner failed before compilation because it launched from the repository root instead of the nested Go module.

## Scope and source identity

- Project: `swe-af`
- Nested module: `go/go.mod`
- Module SHA-256: `12884afec8b7feed6ff3c50e28f0a889ca9066fdb9f56d826ff42f8111592127`
- Executor SHA-256: `f5d89e592cda54615241563c02e852afedae47f9a105e94531766ecb2f57d844`
- Regression test SHA-256: `113bc50a72adc082b3ad22ca9cc21b385a4bc6a67deb9aef7b1a869cb83fa370`
- C4 source of truth: `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`
- C4 decision remains `UNACCEPTED / BLOCKED`; this task does not promote it.

## C4 API signature

The reviewed executor entry point is:

```go
func RunDAG(
    ctx context.Context,
    planResult map[string]any,
    repoPath string,
    callFn CallFn,
    nodeID string,
    cfg *config.ExecutionConfig,
    opts ...Option,
) (*schemas.DAGState, error)
```

The cleanup contract is source-backed: after each Git level, `RunDAG` starts cleanup asynchronously, awaits the cleanup completion channel, invokes `WithLevelCompleteFn`, and only then increments `CurrentLevel`. Cleanup failures from the non-fatal workspace-cleanup path are retried once and do not bypass the await barrier.

## Reproduction and classification

The original symptom was reported as `TestCleanupAwaitedBeforeAdvance: expected cleanup after setup, got cleanup=-1`.

Source inspection shows the barrier is gated by the fixture's Git path: `callFn` must be non-nil and `GitIntegrationBranch` must be non-empty, and the active level must produce a branch to clean. The regression fixture explicitly supplies `WithGitConfig` with `integration_branch`, two sequential levels, and a scripted call function. Therefore, the `cleanup=-1` symptom is classified as **fixture/path activation evidence pending runtime confirmation**, not a proven executor defect.

## Added checks

`go/internal/dag/northstar_cleanup_barrier_20261013_test.go` adds:

1. `TestCleanupAwaitedBeforeAdvance` — positive two-level ordering check: first cleanup must occur before second-level setup.
2. `TestCleanupBarrierRetriesNonFatalFailureBeforeAdvance` — adversarial cleanup `success=false` path must retry once and still advance only after the barrier.
3. `TestCleanupBarrierHonorsPreCancelledContext` — adversarial pre-cancel path must make no Git setup or cleanup calls.

## Runner receipt

```text
command: go test ./internal/dag
status: BLOCKED
code: go: cannot find main module, but found .git/config in /workspace
engine: docker:go-test
source_identity_sha256: 113bc50a72adc082b3ad22ca9cc21b385a4bc6a67deb9aef7b1a869cb83fa370
exit_code: 1
stderr: go: cannot find main module, but found .git/config in /workspace; to create a module there, run: go mod init
```

The authorized runner did not enter `go/`, so none of the three tests has a valid PASS or FAIL result. The same denied invocation must not be reclassified as product evidence or retried identically. The required recovery is to run the focused test from `go/` (for example, `go test ./internal/dag`) once the runner supports the nested module working directory.

## Remaining

- Runtime reproduction and positive/adversarial exit evidence remain BLOCKED by runner working-directory configuration.
- No change was made to `go/internal/dag/executor.go`; source evidence did not establish a genuine runtime bug.
- Independent review and final acceptance remain with the coordinator.
