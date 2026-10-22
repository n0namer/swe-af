# Base20 Issue-Build Repair Evidence — 2026-10-13

## Scope and source of truth

- Task: `swe-af:v2-base20-issue-build-repair-ladder-20261013-r1`
- Project source of truth: `PLAN.md` (active; current build goal requires a real task, deterministic oracle, preserved interrupted work, and recovery without blind repeated effects).
- Bounded implementation: reuse `_implement_issue_impl`; no new orchestrator, no push, no PR, and no external provider calls.
- Source examined: `swe_af/issue/build.py`, `swe_af/issue/git_ops.py`, `tests/issue/test_implement_issue.py`, and `tests/issue/test_implement_issue_wave_full_cycle.py`.

## Demonstrated defect and repair

`git_ops.add_worktree()` retried every failed `git worktree add`, although only transient repository-lock failures are retryable. A deterministic branch collision therefore repeated the same mutation attempt and sleep/backoff.

The repair adds `_is_retryable_worktree_add_failure()` and makes `add_worktree()` stop immediately for `reference already exists` and `branch named` failures. Lock-contention signatures retain bounded retry/backoff. The receipt now records the actual attempt count, including zero attempts.

## Executable contract added

`tests/issue/test_build_recovery_goal_wave.py` uses disposable local Git repositories and `_implement_issue_impl` for:

1. A real bounded retry-helper task with a deliberate first nonblocking RED review, same-task second-pass repair, subprocess oracle, result branch, and preserved foreign dirty caller file.
2. A second task interrupted after artifact write; the checkpoint branch must retain the artifact for recovery.
3. An empty-worktree `commit_all()` call that must return no revision and preserve `HEAD`.
4. A deterministic pre-existing branch collision that must raise without retry/sleep and leave no worktree or branch mutation.
5. Exhausted nonblocking review budget with one checkpoint commit and no duplicate branch effect.

These are offline scripted-agent contracts; Git setup, worktrees, commits, branch inspection, and verifier subprocesses are real. They are not live-provider or independently accepted product proof.

## Validation receipt

Requested command attempted once:

```text
python -m pytest -q tests/issue/test_build_recovery_goal_wave.py
```

The scoped test runner returned:

```text
status=BLOCKED
code=PYTHON_VENV_UNAVAILABLE
engine=python:pytest
stderr=A regular local .venv Python executable is required.
source_identity_sha256=c7d66aa97ed5ca1edee5fb3a04b57d382040e05dd3e96c6225908bf8519937f2
```

No identical rerun was attempted. Therefore no RED/GREEN execution, Git revision receipt, measured token cost, or wall-time claim is asserted by this worker. Independent QA must provide a valid Python environment, run the command above, then run the affected Python issue regression suite and inspect the resulting disposable Git branches.

## Acceptance boundary

- Candidate implementation only; independent verifier/reducer owns acceptance.
- No push, remote, PR, deployment, task administration, or foreign worktree mutation occurred.
- Remaining validation requirement: provision a regular `.venv` Python executable and execute the focused suite before any product-pass claim.
