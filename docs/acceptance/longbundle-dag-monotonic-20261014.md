# Longbundle DAG monotonic recovery — 2026-10-14

## Status

**CANDIDATE / UNACCEPTED.** The authorized Go test artifact is written and source-reviewed, but the scoped runner could not execute it because it started outside the nested Go module. No `ACCEPTED_GREEN` claim is made.

## Binding and source identity

| Field | Value |
|---|---|
| Project | `swe-af` |
| Task | `swe-af:v2-longbundle-dag-monotonic-recovery-20261014-r1` |
| Attempt | `attempt-7e72c9e9-8a35-4cb5-afdc-d07c7ce18dd3` |
| Session | `manual-20261014-base5-longbundles-northstar-r1--940479ca8176e08c` |
| Git baseline | `064a7226a500365005b5b662044fc9014a328722` |
| Focused source identity | `d90b2f023cb7e43e33f6039efa08cc2be7e6cc0d60d2559e5918d4c54791c00d` |
| Authorized test | `go/internal/dag/northstar_longbundle_dag_monotonic_20261014_test.go` |

The focused source identity is the runner-provided SHA for the written test source and is retained as immutable evidence for this attempt.

## B1 — Inspect and contract map

- `saveCheckpoint(*schemas.DAGState, noteFunc)` persists `execution/checkpoint.json` and fences stale level, initial revision, and terminal-status regression.
- `loadCheckpoint(string) *schemas.DAGState` restores a serialized `schemas.DAGState` for resume.
- `checkpointRegressesTerminalIssue(*schemas.DAGState, *schemas.DAGState) bool` protects terminal issue effects.
- `recoverableIssueIDs(*schemas.DAGState) []string` excludes terminal statuses and sorts remaining issue IDs.
- Compile-time assignments in the focused test pin all four function/DTO signatures.

## B2 — Reproduction and classification

The first authorized test execution was attempted through the scoped project runner. It failed before compilation:

```text
status: FAIL
engine: docker:go-test
source_identity_sha256: d90b2f023cb7e43e33f6039efa08cc2be7e6cc0d60d2559e5918d4c54791c00d
exit_code: 1
package: ./go/internal/dag
go: cannot find main module, but found .git/config in /workspace
to create a module there, run: go mod init
go_image: golang:1.24.0-alpine
```

Classification: **INFRA BLOCKED**. The task requires execution from the nested Go module with Go 1.26.5, while the authorized runner starts at `/workspace` and selects no module. The identical denied call was not retried.

## B3 — Scoped implementation

Added only the authorized test artifact. It uses `t.TempDir()` and real JSON checkpoint I/O; it does not mock persistence or alter production code.

Assertions cover:

1. Positive checkpoint write, load, and crash/restart recovery.
2. Stale lower-level checkpoint rejection.
3. Initial-revision mismatch rejection.
4. Terminal `accepted` effect idempotence and regression fencing.
5. Terminal statuses excluded and recoverable IDs sorted deterministically.

## B4 — Focused validation

Requested validation:

```text
go test ./internal/dag -run TestNorthstarLongBundle -count=1
```

Execution through `go_wave.run_project_tests` is blocked by the runner root/module mismatch recorded in B2. Therefore there is no truthful exit-0 test claim and no affected regression claim for this attempt.

## B5 — Self-review and handoff

- Changed paths are limited to the manifest write scope.
- No production source, task state, lease, lock, commit, deploy, or remote state was changed.
- The test source is source-backed against `checkpoint.go` and the existing DAG state contract visible in the permitted read scope.
- Independent reviewer/reducer must rerun the focused test from the nested Go module using the required Go 1.26.5 environment, then decide acceptance.

## Remaining blocker

Provide a scoped runner invocation rooted at the nested `go` module (or materialize the authorized Go 1.26.5 module environment) and rerun `TestNorthstarLongBundle`. Until that receipt has exit 0, this artifact remains `UNACCEPTED`.
