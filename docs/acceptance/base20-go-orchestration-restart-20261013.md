# Base20 Go orchestration restart evidence — 2026-10-13

## Scope

- Task: `swe-af:v2-base20-go-orchestration-restart-20261013-r1`
- Source of truth: `PLAN.md` (active; canonical project branch `dev`).
- Candidate scope: `go/internal/orch/resume.go` and its dedicated offline contract test.
- Non-goal: this candidate does not claim live AgentField execution, independent QA, reducer acceptance, or product `GREEN`.

## Reproduction and bounded result

Source inspection found that `ResumeBuildHandler` already checks `ctx.Err()` before checkpoint access and before its single `execute` control-plane call. No narrowly-scoped defect in `go/internal/orch/resume.go` was proven, so this candidate does not alter production resume behavior.

The dedicated offline contract test now traces the handler against a dependency-aware independent oracle. A fresh request still reconstructs the persisted plan and delegates exactly once with `resume=true`; final effect, stale-worker classification, bounded retry, and acceptance remain the executor's typed-result responsibility, never an inference from a transport envelope's `status` field.

## Deterministic contract scenarios

| Scenario | Input | Expected result |
| --- | --- | --- |
| Multi-step restart | Persisted three-level graph: `prepare -> implement -> verify`, with accepted `prepare` result | The handler forwards the original graph and persisted artifacts directory exactly once with `resume=true`; the dependency-aware independent oracle identifies only `implement` as ready, even if `implement` and `verify` are placed in the same candidate level. |
| Stale/lost worker | Same persisted checkpoint, but cancelled request context | The handler returns `context.Canceled` and makes zero `execute` calls. |
| Ambiguous worker result | Raw execute envelope reports `status: timed_out`, `worker_state: lost`, and a typed dispatch record | The handler returns that raw result unchanged and does not create a terminal acceptance decision or retry from the envelope. |

## Evidence boundary

The focused package test is an offline contract check using a local mock control-plane app and an independent dependency-aware ready-node oracle. It proves handler cancellation, raw-result preservation, and forwarding behavior only. The actual DAG's dependency scheduler, worker staleness classification, bounded repair policy, and idempotent effect persistence are owned by the executor/checkpoint implementation outside this task's write scope and require independent integration QA before acceptance.

## Validation

- Requested command: `cd go && go test ./internal/orch -count=1`
- Focused test set: `TestResumeGoalWaveCancelledRequestDoesNotDispatch`, `TestResumeGoalWaveReturnsAmbiguousWorkerResultWithoutAcceptance`, and `TestResumeGoalWaveRestartsFromCheckpointWithoutSynthesizingAcceptance`.
- Actual scoped attempt: `go_wave.run_project_tests({file:"go/internal/orch/resume_goal_wave_test.go"})` returned `BLOCKED` / `GO_DOCKER_IMAGE_UNAVAILABLE`; its `docker:go-test` engine could not find `golang:1.24.0-alpine`.
- This is an infrastructure block, not a passing test result. The exact requested Go command remains the next validation after the configured image is available.
