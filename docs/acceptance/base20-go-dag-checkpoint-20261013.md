# Base20 Go DAG checkpoint fencing — 2026-10-13

## Scope and source of truth

- Task: `swe-af:v2-base20-go-dag-checkpoint-fencing-20261013-r1`
- Project SoT: `PLAN.md` (active; normal container-first SWE programming)
- Bounded source: `go/internal/dag/checkpoint.go`
- Regression: `go/internal/dag/checkpoint_goal_wave_test.go`
- No scheduler, queue, CrewSync state, task lease, or live runtime state was changed.

## Exercised checkpoint contract

The regression constructs one bounded DAG with ordered node IDs:

1. `edit-file` is completed with immutable source revision `9f8e7d6c5b4a3210` and a file-write effect ID.
2. `run-tests` is left `unknown` after a simulated crash with stale owner `stale-worker` and ambiguous result.
3. `accept-result` remains pending.

It writes and reloads the actual Go checkpoint serializer, verifies the node IDs, revision, and checkpoint level, then uses the production recovery selector. Only `run-tests` and `accept-result` are recoverable; the completed edit is excluded.

Negative paths retain the original checkpoint when a writer attempts to:

- replace it with a different immutable revision;
- replace it from a lower checkpoint level; or
- use the same level and revision but regress `edit-file` from `completed` to `pending`; or
- use the same level and revision but omit the completed `edit-file` node entirely.

The source repair makes terminal node status immutable across a checkpoint write, rejects omission of an existing terminal node, and treats an empty incoming revision as mismatched when an existing checkpoint has a revision. This prevents a stale snapshot from making a completed file effect eligible for recovery again.

## Validation receipt

Requested command: `cd go && go test ./internal/dag -count=1`

Executed through the authorized scoped runner on `go/internal/dag/checkpoint_goal_wave_test.go`:

- Status: `BLOCKED`
- Engine: `docker:go-test`
- Blocker: required local image `golang:1.24.0-alpine` is unavailable.
- Runner source identity before the edit: `c1e060b13368dd043e6a3e4e54401f579ccbc7da14b59397bc50b6d7768fe080`.

No PASS claim is made. The next deterministic validation is the requested package command after the authorized runner/image is available; independent QA/reducer still owns product acceptance.
