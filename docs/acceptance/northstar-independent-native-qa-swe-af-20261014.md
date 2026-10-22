# Northstar independent native QA — 2026-10-14

## Decision

**QA disposition: UNACCEPTED / INFRA BLOCKED.**

The authorized oracle and acceptance evidence are present, but the focused
pytest engine could not start because the required regular project `.venv` is
unavailable. No runtime PASS is claimed. Independent reducer/reviewer and
CrewSync terminal readback remain required.

## Scope and authority

| Item | Evidence |
|---|---|
| Project | `swe-af` |
| Task | `swe-af:v2-independent-native-longbundle-qa-20261014-r1` |
| Attempt | `attempt-f6d9dff1-54ec-4a5a-9d0a-fde5404698c4` |
| Session | `manual-20261014-base5-independent-native-qa-hyp-0a8a360589bb0201` |
| Write scope | `tests/issue/test_northstar_independent_gitreplay_oracle_20261014.py`; this report |
| Production/concurrent writes | None performed; source `swe_af/issue/git_ops.py` remained read-only |

## Independent C4 source identity

The inspected read-only source is `swe_af/issue/git_ops.py` with immutable
SHA-256:

```text
0f73d9fb3a613f1b68ccf4a8161a93c1c5e6dd6a4a392b2c2c8a6157b7c2403d
```

The oracle pins these public helper signatures from that exact source:

- `add_worktree`
- `commit_all`
- `new_commits`
- `changed_files`
- `remove_worktree`

The repaired test artifact SHA-256 is:

```text
9145667b831a04d604488a4d8d20b795ece8b494e4b9f4c6cdc854f213f3f198
```

The repair moves the positive linked worktree to pytest's external `tmp_path`.
This preserves real branch/worktree isolation while preventing `.worktrees/`
from appearing as an untracked root-repository directory.

The producer test artifact is required and contains the public
`implement_issue` entrypoint, worktree handling, and verification assertions.
No producer `worker-result` receipt, task identity, attempt identity, or
session readback is present in the permitted source scope. This comparison is
therefore explicitly **BLOCKED**, not inferred.

## Oracle coverage

The new test file contains four checks:

1. **Positive native commit** — creates a real disposable repository and
   isolated worktree, commits exactly one changed file, verifies the branch
   delta, caller `HEAD`, clean status, no remote, and cleanup.
2. **Negative branch collision** — retries a deterministic existing branch and
   requires `GitOpsError`; verifies no extra commit and no caller mutation.
3. **Restart/replay** — checkpoints a repair, removes the interrupted worktree,
   recreates a separate replay worktree from the checkpoint, verifies identical
   source content, and requires the replay to create no second commit.
4. **C4 source/receipt contract** — verifies the immutable source SHA and helper
   signatures, confirms the producer artifact exists, and fails closed when a
   producer worker-result receipt is absent.

All Git operations use disposable local state. No remote, production branch,
Ruflo Native task, lease, lock, or terminal state is mutated.

## Validation receipts

### Focused oracle

```text
command: .venv\Scripts\python.exe -m pytest -q tests/issue/test_northstar_independent_gitreplay_oracle_20261014.py
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: 9145667b831a04d604488a4d8d20b795ece8b494e4b9f4c6cdc854f213f3f198
exit_code: null
stderr: A regular local .venv Python executable is required.
```

This is an infrastructure receipt, not a test PASS. The identical denied call
must not be retried until the regular project `.venv` exists.

### Affected regression

```text
requested: tests/issue/test_git_ops.py
status: BLOCKED
code: test_outside_task_write_scope
engine: go_wave.run_project_tests
exit_code: null
```

The authorized runner rejected the regression file because it is outside this
attempt's permitted test/write scope. No fabricated regression result is
reported.

## Acceptance checklist

- One positive native case: **IMPLEMENTED; runtime BLOCKED**.
- Two fail-closed negative/restart cases: **IMPLEMENTED; runtime BLOCKED**.
- Exact command and exit receipt: **PASS for evidence capture; execution BLOCKED**.
- Independent C4 signatures and immutable source SHA: **PASS by source readback**.
- Producer receipt comparison: **BLOCKED; receipt absent from permitted scope**.
- No production or concurrent QA writes: **PASS by execution policy/evidence**.
- Typed result and separate reducer acceptance: **PENDING; worker remains UNACCEPTED**.

## Affected-only next move

Materialize the regular project `.venv`, then run exactly:

```text
.venv\Scripts\python.exe -m pytest -q tests/issue/test_northstar_independent_gitreplay_oracle_20261014.py
```

After a real pytest receipt exists, the independent reviewer should compare the
producer's original CrewSync receipt against this source SHA and test artifact,
then perform the affected-only reducer readback. Do not claim
`ACCEPTED_GREEN` from this worker publication.
