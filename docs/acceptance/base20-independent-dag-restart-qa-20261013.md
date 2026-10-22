# Base20 independent DAG restart QA — 2026-10-13

## Decision

**ACCEPT_SCOPED/RED** — this packet independently audits the two finished candidate surfaces as one recovery boundary, but it does not self-accept the product. The existing candidates demonstrate offline checkpoint persistence and cancellation-aware resume forwarding; they do not yet demonstrate a fail-closed, idempotent native recovery lifecycle.

## Bounded sources and exact revisions

| Surface | SHA-256 | Independent finding |
| --- | --- | --- |
| `go/internal/dag/checkpoint.go` | `a65104dcc74b4aed5b7975cb4edc7c90eec877f5c899b9e5d70ca375475fc19c` | Persists/loads the DAG state and rejects a lower-level or changed-source checkpoint write. It has no checkpoint schema/version field or atomic writer/fencing epoch. |
| `go/internal/dag/checkpoint_goal_wave_test.go` | `c1e060b13368dd043e6a3e4e54401f579ccbc7da14b59397bc50b6d7768fe080` | Producer evidence for completed-effect exclusion, ambiguous test state, and stale-write rejection. |
| `go/internal/orch/resume.go` | `02c1a49916f19a1e489ad8a94d6f76aea2ff39869b3bec9351bd96ca33cb7f48` | Stops a pre-cancelled request, reconstructs persisted graph data, and calls `execute` with `resume=true`; it does not validate checkpoint version, source revision, fencing epoch, duplicate-resume state, or terminal receipt. |
| `go/internal/orch/resume_goal_wave_test.go` | `cf9cfc41d4e5bfc035e56f3ee4bbc3b7a785d9ce73e33d10fa183a923ca84008` | Producer evidence for cancellation, raw-result forwarding, and an offline next-ready-node oracle. |
| `go/internal/dag/base20_independent_checkpoint_resume_qa_test.go` | `4b8d3e5567171b4b3a031378c5a5442af10301b3ad2f3336676a996d501f2512` | New independent QA bundle. |

The project source of truth is `PLAN.md` (`452ebfff07df68ba55feaac214190a3ab5c0a3296406bcef19e6b45192ed4cc7`), which requires recovery without blindly restarting and says that timeout/ambiguous results must inspect post-state before retry.

## Executable QA scenarios

The new offline QA test deliberately combines both candidates through the actual exported `orch.ResumeBuildHandler` and the actual unexported DAG checkpoint serializer/loader.

1. **Crash between source edit and verification:** checkpoint records an accepted `source-edit` effect, then `verify-edit` as `unknown` / `timeout-after-edit`; its independent ready oracle selects only `verify-edit`, never repeats the accepted edit.
2. **Stale writer and moving source:** actual `saveCheckpoint` rejects both an older level and a same-level different-source revision, retaining the original persisted recovery point.
3. **Duplicate resume:** two actual calls to `ResumeBuildHandler` make two actual `execute` dispatches. This is an observed repair boundary, not a successful idempotency claim.
4. **Positive terminal evidence:** a typed raw execution result with matching source revision and `verification: passed` satisfies the independent receipt audit, while the handler itself continues to return the raw typed result rather than self-accepting it.
5. **Negative fail-closed audit:** mismatched/missing checkpoint version, moving-source terminal evidence, and timeout verification receipts are rejected by the independent QA receipt validator.

## Validation receipt

Requested focused command:

```text
go test ./internal/dag -run Base20Independent -count=1
```

Required affected-package command:

```text
go test ./internal/dag ./internal/orch -count=1
```

Actual scoped execution attempt:

```text
go_wave.run_project_tests({file:"go/internal/dag/base20_independent_checkpoint_resume_qa_test.go"})
=> GO_WAVE_NODE_ERROR:test_file_not_allowed
```

This runner supports Node/Vitest/tsx files only, so neither Go command was runnable in this worker. The denial is an **INFRA BLOCKED validation gap**, not a passing test result. Per the task policy, the identical denied request was not retried.

## Required same-owner repair boundary

Keep this QA packet RED until the owner of the producer sources supplies all of the following, with focused Go evidence:

1. A typed checkpoint version contract, fail-closed compatibility handling, and a source revision/fencing epoch check before resume dispatch.
2. A durable idempotency key or compare-and-swap state transition that prevents duplicate `ResumeBuildHandler` calls from dispatching the same unresolved effect twice.
3. A recovery classifier that treats `unknown`/timeout as an inspection-required state and only schedules affected nonterminal work after source/effect receipt reconciliation.
4. A typed terminal receipt bound to checkpoint source revision, effect identity, verification result, and writer fencing epoch; orchestration must not infer product acceptance from transport `status: succeeded`.
5. A focused Go runner execution of both requested commands, including negative mismatch, stale-fencing, duplicate-resume, and accepted-effect non-replay cases.

## Native integration limitation

No native AgentField/control-plane integration was executed or implied. The mock app only makes the resume handler's `execute` boundary observable. Live AgentField idempotency, durable ownership/fencing, real worker crash recovery, and source/effect post-state reconciliation remain separately unproven.

## Independent conclusion

The candidate pair has useful local safety evidence: checkpoint stale-write protection, resume cancellation before dispatch, and raw typed-result forwarding. It is insufficient for an independent recovery-green decision because required versioning, fencing, duplicate prevention, and native AgentField lifecycle proof are absent or not evidenced. This QA packet is therefore **ACCEPT_SCOPED/RED** with the repair scope limited to the two producer surfaces and their focused tests; it makes no scheduler, deployment, or acceptance-state change.
