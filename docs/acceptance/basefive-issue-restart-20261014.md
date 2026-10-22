# Basefive issue Git restart — 2026-10-14

## Scope

This bounded artifact records the source-bound C4 helper surface and the
focused native Git lifecycle test authorized for the Basefive issue restart
attempt. It does not claim product acceptance or mutate task, lease, or
terminal state.

## D1 — Source and interface identity

- Project: `swe-af`
- Source: `swe_af/issue/git_ops.py`
- Source SHA-256: `0f73d9fb3a613f1b68ccf4a8161a93c1c5e6dd6a4a392b2c2c8a6157b7c2403d`
- Pinned C4 helpers: `add_worktree`, `commit_all`, `new_commits`,
  `changed_files`, `remove_worktree`
- Public error type: `GitOpsError`

The source was already implemented before this attempt. No production change
was made without a reproduced defect.

## D2/D3 — Focused assertions

`tests/issue/test_basefive_git_issue_restart_20261014.py` covers:

1. Positive disposable-repository issue commit: one new branch commit, exact
   changed file, caller `main` unchanged, no remote, and clean caller status.
2. Negative branch collision: the second worktree creation raises
   `GitOpsError`, creates no additional commit, and leaves the caller HEAD
   unchanged.
3. Restart/replay: an interrupted checkpoint is reused from its exact commit;
   replaying the unchanged repair produces no duplicate commit, while the
   original branch retains exactly one new commit.
4. C4 signature and source-SHA comparison.

## D4 — Runner evidence

The authorized native runner could not produce a Python test exit status in
this attempt. The preserved receipts are:

```text
command: go_wave.run_project_tests(tests/issue/test_northstar_independent_gitreplay_oracle_20261014.py)
status: BLOCKED
code: test_outside_task_write_scope
exit_code: null
```

```text
command: go_wave.run_project_tests(tests/issue/test_basefive_git_issue_restart_20261014.py)
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: 3aada0fd020b71e692795feba36216230de4819a425a84dfe5d88bbe34ba44ae
exit_code: null
stderr: A regular local .venv Python executable is required.
```

The requested validation command remains:

```text
python -m pytest -q tests/issue/test_basefive_git_issue_restart_20261014.py
```

It was not promoted to PASS because the scoped runner requires the regular
project `.venv`, which is unavailable. No identical denied call was retried.

## D5 — Self-review

- Changed files remain inside the authorized write scope.
- No source defect was proven, so `swe_af/issue/git_ops.py` was left unchanged.
- No commit, push, deploy, remote mutation, or task administration occurred.
- Independent reviewer/reducer and authoritative terminal readback remain
  required before any `ACCEPTED_GREEN` decision.
