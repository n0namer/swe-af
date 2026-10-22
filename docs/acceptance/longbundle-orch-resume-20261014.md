# Longbundle ResumeBuildHandler — 2026-10-14

## Decision

**UNACCEPTED / INFRASTRUCTURE-BLOCKED.** The scoped test and source-backed contract evidence are present, but the authorized Go runner could not locate the module and therefore did not execute the focused assertions. No `ACCEPTED_GREEN` claim is made.

## Scope and source identity

| Item | Exact value | Status |
|---|---|---|
| Task | `swe-af:v2-longbundle-orch-typed-resume-20261014-r1` | PASS — trusted binding |
| Project | `swe-af` | PASS — trusted binding |
| Attempt | `attempt-51ad76d9-4c8f-4107-9e13-29667b71af19` | PASS — trusted binding |
| Session | `manual-20261014-base5-longbundles-northstar-r1--9484499db5dbb6c6` | PASS — trusted binding |
| `resume.go` SHA-256 | `02c1a49916f19a1e489ad8a94d6f76aea2ff39869b3bec9351bd96ca33cb7f48` | PASS — read receipt |
| `checkpoint.go` SHA-256 | `ab3d45da6b49558d9ca65c3308d69cacc55b7b10236b614958df4dd82e2f64d8` | PASS — read receipt |
| Added test SHA-256 | `c1e4f2e695c4cb5cc78b0d0c25f8214416afc0f3759cf6f200821f4fe474a7bf` | PASS — write receipt |

## B1–B5 evidence

### B1 — Inspect and map

- `ResumeBuildHandler` has the stable signature `func(context.Context, *Deps, map[string]any) (any, error)`.
- The existing implementation loads `<repo>/<artifacts>/execution/checkpoint.json`, reconstructs `plan_result`, and calls `execute` through `CallRaw` with `resume: true`.
- Existing checkpoint persistence rejects stale levels, changed initial revisions, and terminal-status regression before writing.
- The permitted source set does not expose owner/lease DTOs or a typed resume-receipt interface; those claims are not invented.

### B2 — Reproduction and classification

- Source-backed positive assertions cover signature, immutable source SHA, checkpoint reconstruction, dependency ordering, raw result preservation, and `resume: true`.
- Negative assertions cover missing checkpoint and malformed checkpoint fail-closed behavior.
- Crash/restart assertions cover cancellation guards before and after reconstruction and the required checkpoint fields used to rebuild the resume plan.
- Runtime classification: **INFRASTRUCTURE-BLOCKED**, not product PASS.

### B3 — Scoped implementation

- Added `go/internal/orch/northstar_longbundle_resume_20261014_test.go` only.
- Added this acceptance receipt only.
- No production source, leases, tasks, locks, generated files, or unrelated dirty files were changed.

### B4 — Focused validation

```text
runner: go_wave.run_project_tests
requested: go/internal/orch/northstar_longbundle_resume_20261014_test.go
engine: docker:go-test
package: ./go/internal/orch
authorized source_identity_sha256: c1e4f2e695c4cb5cc78b0d0c25f8214416afc0f3759cf6f200821f4fe474a7bf
status: BLOCKED
exit_code: 1
stderr: go: cannot find main module, but found .git/config in /workspace
       to create a module there, run: go mod init
```

The runner failed before compilation. The identical denied runner call must not be retried until its working directory/module mounting is corrected.

### B5 — Self-review and handoff

- The test file remains within the exact write scope and is under the project’s 500-line guidance.
- Assertions are deterministic and use real temporary filesystem checkpoints for the two negative paths; no weak dependency mock was introduced.
- Independent QA/reviewer must run `go test ./internal/orch -run TestNorthstarLongBundle -count=1` from the repository root after the Go module is mounted.
- Independent QA must separately verify owner/lease mismatch fail-closed behavior and typed receipt compatibility against the producer/DTO sources, which were outside this attempt’s read scope.

## Acceptance matrix

| Requirement | Status | Evidence |
|---|---|---|
| C4 function/DTO signature and immutable source SHA | PARTIAL / PASS | Handler signature and `resume.go` SHA are recorded; DTO/owner/lease compatibility remains unverified. |
| Positive resume assertion | BLOCKED | Test authored; runner failed before execution. |
| Two negative assertions | BLOCKED | Missing and malformed checkpoint tests authored; runner failed before execution. |
| Crash/restart assertions | BLOCKED | Cancellation and reconstruction-order tests authored; runner failed before execution. |
| Focused command exit 0 | BLOCKED | Go runner module-root failure; no exit-0 evidence. |
| Affected regression | SKIP | Same infrastructure blocker prevented focused execution; no broad retry issued. |
| Independent review and reducer | PENDING | Coordinator must route original task evidence to independent QA/reducer. |

## Next move

Restore the Go runner’s repository-root/module mount, run the exact focused command once, then reuse this task’s test and acceptance artifact for sticky independent QA. Do not claim `ACCEPTED_GREEN` from this candidate or from publisher transport status.
