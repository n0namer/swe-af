# MVP R9 native Go resume fencing producer evidence

## Scope and provenance

- Project: `swe-af`
- Task: `swe-af:v2-hyperwave-go-orchestrator-resume-fencing-20261014-r9`
- Attempt: `attempt-98ee528c-aed8-4a64-a8ad-97394fe14c8c`
- Session: `manual-20261014-basefive-hyperwave-r9-swe-af-v2-7924a8c71b1de6ca`
- Authorized files: `go/internal/orch/resume.go`, `go/internal/orch/northstar_mvp_r9_resume_contract_test.go`, this document
- Existing C4 source of truth: `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`

## Functional delta

`ResumeBuildHandler` retains the existing raw `execute` envelope and adds a bounded native resume protocol:

- checks context cancellation before checkpoint access, before dispatch, and after execution;
- rejects a supplied stale owner or lease epoch before calling `execute`;
- reconstructs dependency-ready issues from persisted `all_issues`, `levels`, and completed issue records;
- forwards owner, lease epoch, run ID, ready issues, and deterministic per-issue replay keys to the existing execute boundary;
- persists the typed execute result and dispatch receipt through an atomic checkpoint replacement;
- leaves legacy callers without owner metadata on the existing compatibility path.

The fake runner contract test models an effect ledger keyed by the replay key. Replaying the same checkpoint therefore produces no duplicate effect, while stale-owner and cancelled paths produce no runner call.

## Focused validation

Requested validation: `go test ./internal/orch -run TestMVPR9 -count=1`.

Actual authorized runner receipt:

```text
status: FAIL
code: null
engine: docker:go-test
test: go/internal/orch/northstar_mvp_r9_resume_contract_test.go
source_identity_sha256: b0048c34f42cb3208086203cf67b61723c5fca14bc656c766a20ccb671d81ebe
exit_code: 1
stderr: go: cannot find main module, but found .git/config in /workspace; to create a module there, run: go mod init
package: ./go/internal/orch
go_image: golang:1.24.0-alpine
```

Classification: **INFRASTRUCTURE BLOCKED**. The runner starts outside the Go module (`go/go.mod`) and fails before compiling or executing tests. The identical denied invocation was not retried.

## Source receipts

- `go/internal/orch/resume.go` SHA-256: `372df0a0ff1340bb8e526beb879b6c75affa51e66b5ac60a24e6ee413ade76f1`
- `go/internal/orch/northstar_mvp_r9_resume_contract_test.go` SHA-256: `b0048c34f42cb3208086203cf67b61723c5fca14bc656c766a20ccb671d81ebe`

## Acceptance state

Candidate only. Focused runtime proof remains unavailable because of the module-root runner failure. Independent QA/reviewer and reducer evidence are still required; no `ACCEPTED_GREEN` or terminal acceptance is claimed.

Next action: run the exact focused Go test from the `go` module root, then run the affected native negative suite and retain the resulting source SHA and exit/test receipt.
