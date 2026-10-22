# MVP r9 Go orchestrator resume/fencing QA — 2026-10-14

## Disposition

**QA_READY_NOT_FINAL / UNACCEPTED.** This is an independent QA artifact for the existing `resume.go` implementation. It changes no product source, task state, lease, owner, or terminal state. Independent reviewer and reducer acceptance remain required.

## Pinned provenance

| Receipt | Exact value |
|---|---|
| Project | `swe-af` |
| Trusted task | `swe-af:v2-hyperwave-independent-go-orchestrator-resume-fencing-20261014-r9` |
| Attempt | `attempt-f5b67b15-aabc-411a-889f-b7ef5e566f66` |
| Crew session | `manual-20261014-basefive-hyperwave-r9-swe-af-v2-d3b5c0c5741019bd` |
| Prelaunch Git SHA | `064a7226a500365005b5b662044fc9014a328722` |
| C4 artifact | `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md` |
| C4 artifact SHA-256 | `feeb1fe6a88e53770831de890fe3674e6425dd40327deb913fbf73f3a01f1f17` |
| Product source | `go/internal/orch/resume.go` |
| Product source SHA-256 | `7961bbc5e15e9a79a0ef93253ee6dd1880eb6208f5107c77ebcbb3dcb2f04ccc` |
| QA test | `go/internal/orch/northstar_mvp_r9_resume_adversarial_qa_test.go` |
| QA test SHA-256 | `129e9947b2aa078ae519e10587036eb2f924e98b3bb02f0a13e3c458f65b5020` |

## Source-backed boundary

`ResumeBuildHandler` checks `ctx.Err()` before filesystem access and again before dispatch, resolves the checkpoint under the requested repository/artifacts directory, reconstructs the persisted issue/level plan, and calls the node-local `execute` reasoner with `resume: true` through `CallRaw`. The handler returns that raw typed result without synthesizing acceptance.

## Independent checks

The authorized QA test adds four focused cases:

- Native positive: preserves a typed persisted execute result and the `resume: true` dispatch contract.
- Cancellation: returns `context.Canceled`, no result, and zero execute calls.
- Foreign repository: uses only the pinned owner repository checkpoint and leaves the foreign checkpoint unchanged.
- Restart/replay: reconstructs the same plan and result on a second invocation without mutating the durable checkpoint or adding calls within either invocation.

## Validation receipt

Requested native command:

```text
go test ./internal/orch -run TestMVPR9 -count=1
```

Status: `INFRA_BLOCKED` before compilation. The scoped runner used `docker:go-test` but searched `/workspace` rather than the `go/` module root and returned:

```text
source_identity_sha256: 129e9947b2aa078ae519e10587036eb2f924e98b3bb02f0a13e3c458f65b5020
exit_code: 1
stderr: go: cannot find main module, but found .git/config in /workspace
        to create a module there, run:
        go mod init
package: ./go/internal/orch
go_image: golang:1.24.0-alpine
```

This is a runner/workdir failure, not a product test result. Do not promote it to PASS and do not repeat the identical denied invocation. The repository does contain the module at `go/go.mod`; the authorized interface did not expose a module-root override.

Affected regression: existing `go/internal/orch` resume tests were inspected but not broadened beyond the authorized native command.

## Acceptance boundary

- Product source modified: **NO**.
- QA test and evidence artifact added: **YES**.
- Positive and negative source-level test contract: **DEFINED**.
- Native Go execution: **INFRA_BLOCKED** before compilation.
- `ACCEPTED_GREEN`: **not claimed**.
- Independent verifier/reviewer, reducer, and coordinator terminal readback: **required**.
