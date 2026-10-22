# BaseFive Go orchestration replay QA — 2026-10-14

## Decision

**Candidate evidence: UNACCEPTED / INFRA BLOCKED.**

This artifact records the source-bound C4 replay oracle and the exact focused-run failure. It does not claim runtime acceptance or `ACCEPTED_GREEN`.

## Authority and source identity

| Receipt | Exact value | Classification |
|---|---|---|
| Project | `swe-af` | PASS — trusted task binding |
| Git SHA | `064a7226a500365005b5b662044fc9014a328722` | PASS — trusted GO Core prelaunch receipt |
| Task | `swe-af:v2-basefive-go-orch-fencing-independent-qa-20261014-r4` | PASS — trusted task binding |
| Attempt | `attempt-2007dd3f-d5fa-4003-96f3-8089ecc29214` | PASS — trusted task binding |
| Session | `manual-20261014-basefive-northstar-hyper-r4-swe-65d28bb30266320e` | PASS — trusted task binding |
| Resume source SHA-256 | `02c1a49916f19a1e489ad8a94d6f76aea2ff39869b3bec9351bd96ca33cb7f48` | PASS — source-bound test assertion |

## D1–D5 evidence

### D1 — C4 signature and baseline

`go/internal/orch/northstar_basefive_replay_oracle_20261014_test.go` binds the `ResumeBuildHandler` signature:

```text
func ResumeBuildHandler(ctx context.Context, deps *Deps, input map[string]any) (any, error)
```

The oracle also binds the resume source SHA above and the trusted baseline Git SHA. It checks the `execute` dependency call, the resume flag, and the result return marker.

### D2 — Positive, adversarial negatives, restart

The test contains:

- a positive checkpoint round-trip assertion for `all_issues`, `levels`, `artifacts_dir`, and `original_plan_summary`;
- a source-order assertion that checkpoint reconstruction precedes dependency execution;
- a missing-checkpoint fail-closed assertion;
- a malformed-checkpoint fail-closed assertion;
- a canceled-restart assertion returning `context.Canceled`.

### D3 — Native ordering, cancellation, stale-owner fencing

The oracle checks the existing native source markers for checkpoint reconstruction, dependency READY ordering, cancellation guards, stale level/revision fencing, terminal regression fencing, and ignored stale writes. No production source was changed.

### D4 — Focused validation receipt

```text
requested focused test: go/internal/orch/northstar_basefive_replay_oracle_20261014_test.go
runner: docker:go-test
status: FAIL
code: null
exit_code: 1
source_identity_sha256: ddd35e582e7d18a05d9a8acb89262dd10f652637894229b2b0055f9e1c54aade
go_image: golang:1.24.0-alpine
stderr: go: cannot find main module, but found .git/config in /workspace; to create a module there, run: go mod init
classification: INFRA BLOCKED — the authorized runner did not reach test execution
```

The requested native command `go test ./internal/orch -run TestBaseFive -count=1` therefore has no process exit status from this attempt. The identical runner denial must not be retried until the runner uses the project module root.

### D5 — Self-review and handoff

Changed files are limited to the two authorized write-scope paths. No production code, task state, leases, locks, branches, remotes, or deployment artifacts were modified. Independent QA/reviewer and authoritative CrewSync terminal readback remain required.

## Remaining

1. Run the focused Go test from the actual module root with `go test ./internal/orch -run TestBaseFive -count=1`.
2. If the focused test runs, execute the affected native orchestration/DAG regression selected by the coordinator.
3. Have the independent reviewer validate the source-bound signature, positive/negative/restart assertions, stale-owner fencing, and the new runtime receipt.
4. Preserve this task identity and use the published worker-result receipt for coordinator/reducer readback.
