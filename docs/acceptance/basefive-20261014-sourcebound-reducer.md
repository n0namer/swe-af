# Basefive source-bound reducer — 2026-10-14

## Decision

**CANDIDATE / UNACCEPTED.** This source-bound reducer records the C4 interface/signature evidence and the available positive, negative, and restart assertions without promoting unexecuted Python or Go checks to PASS. No `ACCEPTED_GREEN` claim is made.

This artifact changes no producer, QA, runtime, task, lease, or terminal state.

## D1 — Source identity and C4 signatures

| Item | Evidence | Status |
|---|---|---|
| Project | `swe-af` | PASS — trusted task binding |
| Baseline Git/source SHA | `064a7226a500365005b5b662044fc9014a328722` | PASS — trusted GO Core prelaunch receipt |
| Authorized write scope | `docs/acceptance/basefive-20261014-sourcebound-reducer.md` | PASS — exact manifest scope |
| Python C4 helper source | `swe_af/issue/git_ops.py`, pinned by `tests/issue/test_northstar_independent_gitreplay_oracle_20261014.py` | SOURCE-BOUND |
| Python signature markers | `add_worktree`, `commit_all`, `new_commits`, `changed_files`, `remove_worktree` | SOURCE-BOUND |
| Go resume interface | `ResumeBuildHandler(ctx context.Context, deps *Deps, input map[string]any) (any, error)` | SOURCE-BOUND |
| Go dependency seam | `deps.CallRaw(ctx, "execute", ...)` | SOURCE-BOUND |
| Go restart marker | `"resume": true` and `return result, nil` | SOURCE-BOUND |
| Python helper source digest | `0f73d9fb3a613f1b68ccf4a8161a93c1c5e6dd6a4a392b2c2c8a6157b7c2403d` | SOURCE-BOUND by the authorized Python oracle |
| Go resume source digest | `02c1a49916f19a1e489ad8a94d6f76aea2ff39869b3bec9351bd96ca33cb7f48` | SOURCE-BOUND by the authorized Go test |

The available source evidence demonstrates the interface seam and its version-pinning checks. It does not provide producer/QA worker-result receipts or prove runtime compatibility across those workers.

## D2 — Positive, negative, and restart assertions

The authorized source contains the following executable assertions:

1. **Positive native Git path:** creates an isolated worktree and one commit, checks the exact changed file, preserves caller `main`/`HEAD`, rejects remotes, and cleans up the worktree.
2. **Negative branch-collision path:** a second `add_worktree` on the same branch must raise `GitOpsError`; no second commit is created and caller `HEAD` remains unchanged.
3. **Negative restart-input paths:** missing checkpoint must fail with `No checkpoint found`; malformed checkpoint must fail closed; canceled restart must return `context.Canceled`.
4. **Restart replay path:** checkpoint state is reconstructed before the execute seam, the unchanged repair produces no duplicate commit, blob identity matches across replay branches, and caller state remains clean.
5. **Checkpoint safety path:** duplicate/regressive edits are rejected through stale-level, revision-mismatch, and terminal-regression guards.

These are source-backed test assertions, not execution receipts.

## D3 — Native SWE-AF seam reduction

The Python oracle compares the native Git boundary through `swe_af.issue.git_ops` and deliberately avoids producer fixture import shadowing. The Go oracle checks the long-bundle resume boundary, checkpoint reconstruction ordering, cancellation behavior, and checkpoint regression guards.

The permitted evidence does **not** include exact producer/QA task identities, DTO schema/version names beyond the source markers above, or runtime worker-result readbacks. Those facts remain blocked rather than inferred.

## D4 — Focused and affected-only validation

### Exact attempted commands

```text
runner: go_wave.run_project_tests
file: tests/issue/test_northstar_independent_gitreplay_oracle_20261014.py
result: GO_WAVE_NODE_ERROR:test_outside_task_write_scope
status: INFRA BLOCKED
exit_code: unavailable
```

```text
runner: go_wave.run_project_tests
file: go/internal/orch/northstar_longbundle_resume_20261014_test.go
result: GO_WAVE_NODE_ERROR:test_outside_task_write_scope
status: INFRA BLOCKED
exit_code: unavailable
```

The scoped runner does not permit either Python or Go test file for this documentation-only write scope. No identical retry was made. The native pytest command and Go test command therefore remain unexecuted in this attempt; neither is promoted to PASS.

Affected regression: **not run** because the only permitted write target is this acceptance document and the available runner rejected both native test paths before execution.

## D5 — Self-review and handoff

- Changed file is exactly inside the manifest write scope.
- No source, runtime, task, lease, lock, or terminal state was modified.
- No duplicate producer, QA, reducer, or repair task is proposed.
- Positive, two negative, and restart assertions are recorded with their source locations and limitations.
- Runtime status is separated as `INFRA BLOCKED`; no false PASS is claimed.
- Independent reviewer/reducer and authoritative terminal readback remain required before any acceptance decision.

## Remaining blockers

1. Restore or authorize the native Python/Go test environment and run the exact focused tests against the recorded baseline/source state.
2. Obtain the original producer and independent-QA worker-result receipts from the coordinator; do not invent replacements.
3. Verify producer/QA C4 DTO/interface identifiers and runtime rejected/interrupted/restarted behavior using those receipts.
4. Perform independent review and coordinator terminal readback; this worker cannot self-accept.

## Reducer conclusion

The strongest source-bound result for this attempt is **CANDIDATE / UNACCEPTED**: the C4 seams and required positive/negative/restart assertions are present in the authorized evidence, while native execution and cross-worker DTO/receipt compatibility remain infrastructure- or scope-blocked.
