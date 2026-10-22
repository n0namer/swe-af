# Affected repair — DAG cleanup barrier test name — 2026-10-14

## Scope and outcome

- Project: `swe-af`
- Task: `swe-af:v2-affected-duplicate-dag-test-repair-20261014-r1`
- Attempt: `attempt-c564efea-6fc4-43c3-a5fc-6af6149f10f7`
- Session: `manual-20261014-base5-affected-qa-repair-wave-r-c50943fe5bd24de1`
- Outcome: `CANDIDATE` / `UNACCEPTED`
- Change scope: rename only the newly isolated duplicate regression test; no production code changed.

## Root-cause evidence

`go/internal/dag/executor_test.go` already defines the package-level test symbol `TestCleanupAwaitedBeforeAdvance`. The isolated regression file also defined that exact symbol, causing Go package compilation to fail before test execution. The compatible repair is to rename only the isolated test to `TestCleanupAwaitedBeforeAdvanceRegression`.

## R1-R5 receipts

### R1 — reproduce and inspect

- Authorized source inspection confirmed the duplicate function name and preserved the original `executor_test.go` test.
- First authorized native attempt:

```text
engine: docker:go-test
test: go/internal/dag/northstar_cleanup_barrier_20261013_test.go
source_identity_sha256: 113bc50a72adc082b3ad22ca9cc21b385a4bc6a67deb9aef7b1a869cb83fa370
exit_code: 1
stderr: go: cannot find main module, but found .git/config in /workspace
        to create a module there, run:
        go mod init
package: ./go/internal/dag
go_image: golang:1.24.0-alpine
```

The runner did not reach the requested nested `/go` module, so this receipt is infrastructure-blocked evidence rather than a product test result.

### R2 — isolate fixture versus product cause

The observed source defect is test-fixture namespace collision: both files declare the same package-level test function. No production implementation or cleanup barrier contract requires modification. The isolated file's assertions remain unchanged:

- cleanup is observed before second-level setup;
- a non-fatal cleanup failure is retried before advancement;
- a pre-cancelled context makes no setup or cleanup calls.

### R3 — minimal repair

Changed exactly one identifier in `go/internal/dag/northstar_cleanup_barrier_20261013_test.go`:

```diff
-func TestCleanupAwaitedBeforeAdvance(t *testing.T) {
+func TestCleanupAwaitedBeforeAdvanceRegression(t *testing.T) {
```

Previous file SHA-256: `113bc50a72adc082b3ad22ca9cc21b385a4bc6a67deb9aef7b1a869cb83fa370`

Current file SHA-256: `4109d386d7c276c24c6c94cc49e0b743d5883a7ede97eea59df5e437ae56e48d`

### R4 — focused validation and boundary review

Second authorized native attempt used the changed source and produced the same environment failure:

```text
engine: docker:go-test
test: go/internal/dag/northstar_cleanup_barrier_20261013_test.go
source_identity_sha256: 4109d386d7c276c24c6c94cc49e0b743d5883a7ede97eea59df5e437ae56e48d
exit_code: 1
stderr: go: cannot find main module, but found .git/config in /workspace
        to create a module there, run:
        go mod init
package: ./go/internal/dag
go_image: golang:1.24.0-alpine
```

Status: `INFRA BLOCKED`; no green claim is made. The requested command `go test ./internal/dag -count=1` could not be executed against the nested `/go` module by the authorized runner. The retry and cancellation negative boundaries remain source-identical apart from the test symbol rename. Owner/fencing and restart behavior were not broadened or weakened; runtime proof remains pending the correct Go-module working directory.

### R5 — self-review and QA handoff

- Changed files are inside the exact write scope.
- Original production test and its symbol were not modified.
- No production files, locks, leases, tasks, commits, or deployment state were touched.
- Independent QA must rerun the package command from the nested `/go` module under Go 1.26.5 and verify duplicate-symbol compilation, cleanup ordering, retry, cancellation, owner/fencing, and restart boundaries.

## Acceptance status

- Original failure eliminated at source level: candidate evidence supports the duplicate-symbol repair, but native compilation is not independently proven here.
- Focused native test exit 0: not verified; runner is blocked by incorrect module root.
- Negative, owner/fencing, and restart boundaries: assertions preserved by the one-identifier diff; runtime execution pending.
- C4/source identity: changed test SHA and exact diff recorded above.
- Independent QA before GREEN: required and not claimed.
