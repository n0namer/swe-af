# Northstar DAG checkpoint fencing C4 QA — 2026-10-13

## Decision

**CANDIDATE / UNACCEPTED.**

A native same-package Go test now exercises the real `saveCheckpoint` and `loadCheckpoint` functions without mocks or an import cycle. Runtime acceptance remains blocked because the authorized Go runner invoked the package from `/workspace` instead of the nested Go module root and exited before compilation.

## Authority and source identity

| Receipt | Exact value | Status |
|---|---|---|
| Task | `swe-af:v2-c4-native-acceptance-dag-fencing-20261013-r1` | trusted binding |
| Attempt | `attempt-b0a9ba0b-7ba5-4f6b-b918-cc05749872fe` | trusted binding |
| Session | `manual-20261013-c4-native-providerbounded-refil-8cb2bc83425b09ca` | trusted binding |
| Git identity | `064a7226a500365005b5b662044fc9014a328722` | prelaunch receipt; not independently mutated |
| Production source | `go/internal/dag/checkpoint.go` | read only |
| Production source SHA-256 | `ab3d45da6b49558d9ca65c3308d69cacc55b7b10236b614958df4dd82e2f64d8` | direct scoped read |
| Focused test | `go/internal/dag/northstar_checkpoint_owner_fence_20261013_test.go` | new scoped artifact |
| Focused test SHA-256 | `cd04f5bd7c1ee0b116feb3977d24580e4c7cbce17f81aa481fee6edd6bf6091a` | final scoped artifact after self-review |
| Pre-self-review runner source identity | `7874653f8b47c1f2d40e2db48494299046df597edff0ddc4f50d6b34fadcb762` | runner exited before compilation; not current test identity |

## C4 boundary

The test remains in package `dag`, so it calls the existing unexported native boundary directly and does not introduce an import cycle:

- `saveCheckpoint(dagState *schemas.DAGState, note noteFunc)` persists or rejects a checkpoint write.
- `loadCheckpoint(artifactsDir string) *schemas.DAGState` returns the decoded checkpoint or `nil` for missing, unreadable, or corrupt input.
- The exercised DTO is `schemas.DAGState`, specifically `ArtifactsDir`, `CurrentLevel`, `GitInitialCommit`, and `AllIssues[].id/status`.
- Terminal monotonicity is evaluated by `checkpointRegressesTerminalIssue` and terminal statuses recognized by `isTerminalIssueStatus`.

## Native scenarios

| Scenario | Test assertion | Runtime status |
|---|---|---|
| Positive round trip | Saves a running checkpoint, loads it, advances to completed, and preserves revision, level, issue ID, and status. | BLOCKED before compile |
| Monotonic terminal checkpoint | Persists revision `revision-a` at level 4 with terminal status `accepted`; concurrent stale candidates cannot replace it. | BLOCKED before compile |
| Stale restart | Concurrent lower-level same-revision candidate is rejected and emits the ignored-checkpoint note. | BLOCKED before compile |
| Terminal regression | Concurrent higher-level same-revision candidate changing `accepted` to `running` is rejected. | BLOCKED before compile |
| Foreign revision | Concurrent `revision-b` candidate is rejected while `revision-a` remains persisted. | BLOCKED before compile |
| Missing checkpoint | `loadCheckpoint` returns `nil`. | BLOCKED before compile |
| Corrupt checkpoint | Invalid JSON returns `nil`, failing closed. | BLOCKED before compile |
| Crash/retry recovery | A later native save replaces corrupt content with a loadable terminal `failed` checkpoint. | BLOCKED before compile |

## Focused runner receipt

The authorized scoped runner was called exactly once for the pre-self-review test identity. The final assertion-only revision was not rerun because the unchanged module-root infrastructure failure would prevent compilation again.

```text
runner: go_wave.run_project_tests
file: go/internal/dag/northstar_checkpoint_owner_fence_20261013_test.go
engine: docker:go-test
package: ./go/internal/dag
source_identity_sha256: 7874653f8b47c1f2d40e2db48494299046df597edff0ddc4f50d6b34fadcb762
status: FAIL
exit_code: 1
stdout: <empty>
stderr: go: cannot find main module, but found .git/config in /workspace
        to create a module there, run:
        go mod init
```

Classification: **INFRASTRUCTURE / RUNNER WORKING-DIRECTORY DEFECT.** The runner selected the correct package path relative to the repository root but executed outside the nested Go module. It never compiled or ran the tests, so no scenario is reported as PASS or product FAIL.

## Runnable commands awaiting infrastructure repair

These commands are source-backed recommendations and were **not run** in this attempt:

```powershell
Set-Location go
go test ./internal/dag -run '^TestNorthstarCheckpointOwnerFence' -count=1
go test ./internal/dag -count=1
```

The focused command is the required next runtime gate; the package command is the affected regression gate. Their exit codes are currently `null` because the authorized runner cannot select the nested module root.

## Edge-case review

The new artifact covers every requested branch reachable through the public task boundary: native positive persistence, lower-level stale restart, terminal regression at a higher level, foreign revision, missing file, corrupt JSON, and retry recovery. Concurrent rejection is synchronized and the final persisted DTO is asserted after all writers complete.

One stronger overlap remains intentionally unclaimed: simultaneous first writes can pass the production read-check before either write is visible, and `checkpoint.go` documents a non-atomic write. The production file is read only in this task, and the blocked runner provides no runtime evidence establishing whether that race manifests in this repository. It must not be labeled fixed or passing.

## Root-cause-only repair route

1. **Runner repair:** resume the existing GO Wave runner/tooling owner and execute Go tests with working directory `/workspace/go`, or equivalent nested-module support. Do not initialize a new module at `/workspace`.
2. **Focused rerun:** run the exact focused command against final test SHA `cd04f5bd7c1ee0b116feb3977d24580e4c7cbce17f81aa481fee6edd6bf6091a` and record full output plus exit code.
3. **Production repair only if reproduced:** if the native test or a deterministic overlapping-write extension fails, resume the original CrewSync owner for `go/internal/dag/checkpoint.go`. Its exact task key is not present in this attempt's permitted evidence and must be resolved from existing CrewSync state rather than invented.
4. **Independent review:** an independent reviewer must inspect the real test receipt and source identities before any GREEN decision. Worker publication is transport evidence only.

## Acceptance checklist

- Actual native positive path and invariant assertions: **IMPLEMENTED / NOT RUN**.
- At least two adversarial negatives including stale/restart: **IMPLEMENTED / NOT RUN**.
- Missing/corrupt checkpoint fails closed and retry recovers: **IMPLEMENTED / NOT RUN**.
- Exact production and test SHA receipts: **PASS**.
- Runnable focused and affected-regression commands: **PASS, explicitly unrun**.
- Focused runtime exit code: **FAIL 1 due runner module root; zero tests executed**.
- C4 function/DTO boundary: **PASS by direct source inspection**.
- Root-cause-only sticky repair route: **PASS, exact foreign owner intentionally not invented**.
- Product `ACCEPTED_GREEN`: **not claimed**.
