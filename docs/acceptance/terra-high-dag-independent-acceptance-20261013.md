# Terra High native DAG independent checkpoint acceptance — 2026-10-13

## Scope and disposition

**Candidate disposition: UNACCEPTED.** This independent QA artifact adds only the authorized same-package native Go conformance test. It does not modify production code, existing tests, CrewSync state, Ruflo Native tasks, providers, or deployment state. Product `ACCEPTED_GREEN` remains reserved for the independent reducer/reviewer.

## Immutable source binding

| Receipt | Exact value |
|---|---|
| Project | `swe-af` |
| Git identity | `064a7226a500365005b5b662044fc9014a328722` |
| Prelaunch manifest SHA-256 | `f9ea540859050c520e05bc0dc2570a188ae9b6c2c056c4f698ef530301d81a9b` |
| Task context SHA-256 | `78fa9180d56f0b34aa94ab3a5eaf1c7b1696305f5335820e5c2c7557d909366d` |
| Active C4 source | `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md` |
| Active C4 source SHA-256 | `feeb1fe6a88e53770831de890fe3674e6425dd40327deb913fbf73f3a01f1f17` |
| New native test | `go/internal/dag/terra_high_native_checkpoint_conformance_20261013_test.go` |
| New native test SHA-256 | `9623256ad22a8259ffd4940d74e64aedcaee93a02fb67259346da333b79a80e7` |

The Git identity binds the read-only `go/internal/dag/checkpoint.go` source and its C4/DTO surface. A standalone blob SHA for that read-only source was not exposed by the authorized scoped reader; none is invented here.

## C4 and DTO surface verified from source

The new test is same-package `dag` code and exercises the implementation directly, without mocks:

- `checkpointPath(dagState *schemas.DAGState) string` resolves the persisted checkpoint location.
- `saveCheckpoint(dagState *schemas.DAGState, note noteFunc)` serializes state and rejects stale, revision-mismatched, or terminal-regressive writes.
- `loadCheckpoint(artifactsDir string) *schemas.DAGState` returns `nil` for unreadable or invalid JSON, providing fail-closed recovery.
- `recoverableIssueIDs(dagState *schemas.DAGState) []string` excludes terminal issue statuses and sorts recoverable issue IDs.
- The observed DTO fields used by this conformance slice are `ArtifactsDir`, `CurrentLevel`, `GitInitialCommit`, and `AllIssues` on `schemas.DAGState`.

`checkpoint.go` imports `github.com/Agent-Field/SWE-AF/go/internal/schemas`. Pinning or validating the future native AgentField dependency version is outside this task's read and write scope; independent review must confirm that dependency against the authoritative module lock/source before any green acceptance.

## Executable conformance cases

| Case | Expected native behavior |
|---|---|
| Persist and select recoverable issues | A valid checkpoint reloads with its source revision; completed effects are excluded while interrupted/pending issues remain selected. |
| Reject stale overwrite | A lower-level checkpoint cannot replace a newer persisted recovery state. |
| Fail closed for corrupt checkpoint | Invalid JSON yields `nil` rather than a replayable state. |
| Interrupted retry preserves accepted effect | Only the interrupted issue is recoverable; the completed effect is never selected again, and a completed retry leaves no recovery work. |

The test uses `t.TempDir`, direct checkpoint persistence, `loadCheckpoint`, and `recoverableIssueIDs`; it does not emulate live behavior with mocks.

## Focused validation receipt

Required command under the nested `go` module with Go `1.26.5`:

```text
go test ./internal/dag -run '^TestTerraHighNativeCheckpoint' -count=1
```

**Status: BLOCKED (runner environment).** The scoped runner attempted the test as `docker:go-test`, but executed from `/workspace` instead of the nested `go` module and selected `golang:1.24.0-alpine`, not required Go `1.26.5`. It exited `1` before compilation with:

```text
go: cannot find main module, but found .git/config in /workspace
	to create a module there, run:
	go mod init
```

Receipt: test `go/internal/dag/terra_high_native_checkpoint_conformance_20261013_test.go`; package `./go/internal/dag`; source identity SHA-256 `78eba1d749e1db324d83fcc39cf370971be2c700f1cfc648cb26eb28d219a255`; exit code `1`. This receipt predates the final post-review test SHA above (`9623256ad22a8259ffd4940d74e64aedcaee93a02fb67259346da333b79a80e7`), which adds a nil assertion after the completed-retry save and remains unexecuted. This is an infrastructure/runner configuration blocker, not a product test failure or PASS. The identical scoped call is not retried.

## Review boundary

- The source, test, and report are candidate evidence only.
- No production/test files outside the two authorized new paths were written.
- An independent reviewer/reducer must inspect the focused test receipt, source binding, future `schemas.DAGState`/AgentField compatibility, and authoritative CrewSync readback before considering any green decision.
