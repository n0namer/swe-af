# MVP R10 swe-af R9 QA repair — 2026-10-14

## Decision

**Candidate / UNACCEPTED.** This bounded QA repair strengthens the native R9 resume adversarial test without changing production orchestration, leases, ownership state, or terminal state. Product acceptance remains owned by the independent verifier/reducer.

## Scope and source identity

| Receipt | Exact value | Classification |
|---|---|---|
| Project | `swe-af` | PASS — trusted task binding |
| Task | `swe-af:v2-r10-r9-qa-repair-20261014-r1` | PASS — trusted task binding |
| Attempt | `attempt-c1f4812e-3fa1-460b-b8a0-6ce2c0277b3c` | PASS — trusted task binding |
| Session | `manual-20261014-basefive-r10-active-swe-af-v2-r-e6248dbfb73f1ffe` | PASS — trusted task binding |
| Active BMad artifact | `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md` | PASS — read directly; it keeps the upstream C4 gate UNACCEPTED/BLOCKED |
| Changed test SHA-256 | `bb4feaf926564c29685fd9b4a667d75d942619528e49280bbf53369f09f5a20a` | PASS — exact post-edit file SHA |
| Pre-edit runner source SHA-256 | `129e9947b2aa078ae519e10587036eb2f924e98b3bb02f0a13e3c458f65b5020` | PASS — runner receipt before this QA delta |

## QA delta

- Checkpoint fixtures now include `last_resume` with owner, resume, effect, and interrupted-status fields.
- The foreign-repository test asserts the pinned `last_resume` DTO is preserved along with the owner issue and that the foreign checkpoint remains byte-identical.
- Invalid payload input fails closed before any planner dispatch.
- Existing positive, canceled-request, foreign-owner, and restart/replay checks remain in the same focused test file.

## Acceptance mapping

| Requirement | Status | Evidence |
|---|---|---|
| C4/source SHA recorded | PASS | Changed focused test SHA above; upstream C4 artifact remains explicitly UNACCEPTED/BLOCKED |
| Real positive returns typed durable result | TEST DEFINED | `TestMVPR9ResumeNativePositivePreservesTypedPersistedResult` |
| Invalid payload/cancel fail closed | TEST DEFINED | `TestMVPR9ResumeInvalidPayloadFailsClosedBeforeDispatch`; `TestMVPR9ResumeCanceledRequestFailsClosedBeforeDispatch` |
| Foreign owner integrity | TEST DEFINED | `TestMVPR9ResumeForeignRepositoryCannotReplacePinnedCheckpoint` |
| Crash/replay no mutation or divergent plan | TEST DEFINED | `TestMVPR9ResumeRestartReplaysWithoutCheckpointMutationOrExtraDispatch` |
| Native focused regression | BLOCKED | Authorized runner could not locate the Go module before compilation |

## Focused runner receipt

```text
command: go test ./internal/orch -run TestMVPR9 -count=1
requested_file: go/internal/orch/northstar_mvp_r9_resume_adversarial_qa_test.go
status: BLOCKED
engine: docker:go-test
go_image: golang:1.24.0-alpine
source_identity_sha256: 129e9947b2aa078ae519e10587036eb2f924e98b3bb02f0a13e3c458f65b5020
exit_code: 1
stderr: go: cannot find main module, but found .git/config in /workspace; to create a module there, run: go mod init
```

The runner failed before compiling the test and therefore provides no product PASS/FAIL result. The post-edit SHA has not been executed in this environment; do not promote the pre-edit receipt to green. The exact validation command should be rerun from the repository/module root once the Go runner is configured for this checkout.

## Self-review

Only these authorized files changed:

- `go/internal/orch/northstar_mvp_r9_resume_adversarial_qa_test.go`
- `docs/acceptance/mvp-r10-swe-af-r9-qa-repair.md`

No production source, lock, lease, task, owner, remote, or terminal state was modified. No duplicate business effect is introduced by the QA fixture; replay assertions require stable reconstructed plans and byte-identical durable checkpoints.

## Remaining

- Independent QA/reducer must execute the focused Go command from the correct module root against the post-edit test SHA.
- Independent review must verify the actual C4 DTO/source compatibility and obtain the coordinator-owned terminal readback before any GREEN or acceptance claim.
