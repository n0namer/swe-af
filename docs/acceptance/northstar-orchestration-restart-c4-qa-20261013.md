# Northstar orchestration restart C4 QA — 2026-10-13

## Decision

**UNACCEPTED / INFRASTRUCTURE-BLOCKED.**

This attempt adds a focused Go test candidate for `ResumeBuildHandler` and executes the authorized GO Wave runner once. The runner starts its Docker Go engine but fails before package compilation because it invokes the nested Go package from repository root without a root module. No test is reported PASS, and no native positive-path claim is made.

## Authority and exact source identity

| Receipt | Value | Status |
|---|---|---|
| Project Git identity | `064a7226a500365005b5b662044fc9014a328722` | trusted prelaunch receipt |
| Production source | `go/internal/orch/resume.go` | read-only |
| Production source SHA-256 | `02c1a49916f19a1e489ad8a94d6f76aea2ff39869b3bec9351bd96ca33cb7f48` | read directly |
| Focused test | `go/internal/orch/northstar_native_resume_20261013_test.go` | new scoped artifact |
| Focused test SHA-256 | `f9dd39706f537e36be13a7b5f6b5aceea741b2a1e39333508b02939638cdb530` | read directly after write |

## C4 function and DTO boundary

`ResumeBuildHandler(ctx context.Context, deps *Deps, input map[string]any) (any, error)` accepts `repo_path`, optional `artifacts_dir`, `config`, and `git_config`. It reads `<absolute repo>/<artifacts_dir>/execution/checkpoint.json`, decodes a map DTO, reconstructs `plan_result`, and calls `deps.CallRaw(ctx, "execute", ...)` with `resume: true`.

The checkpoint-to-plan DTO mapping visible in the scoped source is:

- `all_issues` -> `plan_result.issues`
- `levels` -> `plan_result.levels`
- `artifacts_dir` -> `plan_result.artifacts_dir`
- `original_plan_summary` -> `plan_result.rationale`
- empty `prd`, `architecture`, `review`, and `file_conflicts`

The function returns the raw execute envelope unchanged. Dependency ordering is forwarded through `levels`; the resume boundary does not itself execute or reorder issues.

## Test candidate coverage

The focused artifact contains direct native handler checks for:

1. canceled context rejection before checkpoint access;
2. absent checkpoint rejection with the exact resolved checkpoint path;
3. malformed checkpoint rejection while preserving the durable file for diagnosis;
4. stale checkpoint removal before restart;
5. repeat malformed-checkpoint attempts producing the same error without mutating the checkpoint;
6. source-contract assertions that the durable checkpoint path, `all_issues`, `levels`, raw `execute` call, and `resume: true` remain connected.

The source-contract test is explicitly not labeled a live positive path.

## Authorized runner receipt

```text
invocation: go_wave.run_project_tests({"file":"go/internal/orch/northstar_native_resume_20261013_test.go"})
status: FAIL
engine: docker:go-test
package: ./go/internal/orch
go_image: golang:1.24.0-alpine
source_identity_sha256: f9dd39706f537e36be13a7b5f6b5aceea741b2a1e39333508b02939638cdb530
exit_code: 1
stderr: go: cannot find main module, but found .git/config in /workspace
        to create a module there, run:
        go mod init
```

Classification: **INFRASTRUCTURE-BLOCKED before compilation**. Creating a new root module would be an unauthorized and incorrect workaround. The identical runner call was not repeated.

The likely operator command, subject to confirming the nested module root, is `go test ./internal/orch -run NorthstarNativeResume` with working directory `<project>/go`. It was not run in this attempt and is not a PASS receipt.

## Edge-case review and production gaps

The scoped production function validates checkpoint existence and JSON syntax, but the visible C4 boundary contains no owner/fencing identity field or guard. A syntactically valid foreign-owner checkpoint therefore reaches the downstream `execute` boundary unless another component rejects it.

The resume boundary also has no replay key, completed-issue filter, or duplicate-attempt guard. Each valid invocation calls `execute` again with `resume: true`; no-duplicate-after-retry behavior must be proved in the downstream execute implementation or added by its existing owner.

`all_issues` and `levels` default through `any0`, so missing checkpoint fields are forwarded rather than rejected here. Whether downstream validation safely rejects or normalizes them is outside the permitted read scope.

These are source-backed boundary observations, not runtime defect confirmations. The exact downstream `Deps`/`execute` implementation and its existing tests were outside this task's read scope.

## Acceptance reduction

| Acceptance item | Result | Evidence |
|---|---|---|
| Actual native positive path and invariant assertion | BLOCKED | runner failed before compile; positive dependency seam is outside read scope |
| At least two adversarial negatives including stale/restart | CANDIDATE, unexecuted | canceled, missing, malformed, stale removal, and repeat-failure tests are present |
| Durable checkpoint and replay idempotence | PARTIAL | failed retries preserve checkpoint; successful replay dedupe remains unproved |
| Rejected foreign owner | BLOCKED / production gap | no owner field or guard is visible in `ResumeBuildHandler` |
| Dependency ordering | SOURCE-CONTRACT ONLY | `levels` is forwarded unchanged to `execute`; runtime ordering is unproved |
| Exact source SHA, command, and exit code | PASS | recorded above; no unrun PASS claimed |
| Independent reviewer required for GREEN | REQUIRED | this worker remains `UNACCEPTED` |

## Affected-only repair route

1. Resume the existing GO Wave/test-runner owner and make the Docker Go runner select the nearest nested `go.mod` (or set the working directory to the Go module) instead of creating a repository-root module.
2. Re-run this exact focused test artifact at SHA `f9dd39706f537e36be13a7b5f6b5aceea741b2a1e39333508b02939638cdb530`.
3. If compilation exposes fixture/API mismatch, repair only this test under the original task identity.
4. If the focused test compiles, resume the existing production owner for the downstream `Deps`/`execute` boundary to prove or repair owner rejection, dependency ordering, and successful retry deduplication. The exact CrewSync production-owner task ID is not present in the permitted evidence and must be resolved by the coordinator; no replacement task should be invented.
5. Run independent QA and authoritative terminal readback before any GREEN claim.

## Remaining risk

The candidate has not compiled or executed. The native positive path, foreign-owner rejection, successful crash/restart replay, and no-duplicate-after-retry invariant remain unverified. Product acceptance remains blocked.
