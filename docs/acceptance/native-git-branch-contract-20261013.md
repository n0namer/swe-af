# Native Git branch contract — 2026-10-13

## Decision

**Candidate evidence: UNACCEPTED / BLOCKED.**

The native issue path already contains the relevant branch/worktree safety behavior in `swe_af/issue/git_ops.py`: transient worktree-add failures may retry, while deterministic branch collisions fail without repeating the mutation. This attempt adds an independent disposable-Git oracle for the contract, but the authorized Python runner could not execute because the regular project `.venv` is unavailable.

No production source, foreign QA test, remote, live task, lease, or terminal state was modified.

## Source and C4 identity

| Receipt | Exact value | Classification |
|---|---|---|
| Project | `swe-af` | PASS — trusted task binding |
| GO source Git SHA | `064a7226a500365005b5b662044fc9014a328722` | PASS — trusted prelaunch receipt |
| C4 acceptance signature | `feeb1fe6a88e53770831de890fe3674e6425dd40327deb913fbf73f3a01f1f17` | PASS — trusted reference receipt for `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md` |
| Added test source identity | `0998535604e429e0885dc2fa8e82d3efeed30071fdd7cc9bc351fe7bba576d5a` | PASS — runner-bound source identity |

## Fixture mismatch classification

The original QA assertion checks that the returned branch contains the literal `/restart-oracle`. The issue input is path-like (`../../restart oracle`), while the observed contract described by the task is a namespaced branch of the form `issue/<hash>-restart-oracle`. The independent oracle therefore checks the safe contract with `issue/[0-9a-f]+-restart-oracle`, rejects `..`, and verifies no leading-root path form.

**Classification:** fixture expectation is inconsistent with the intended safe branch naming contract; no production correction is proposed in this attempt. Runtime reproduction of the original pytest assertion remains BLOCKED because the runner denied the original file as outside task write scope before pytest execution.

## Added independent coverage

`tests/issue/test_native_git_branch_contract_20261013.py` covers:

- positive disposable Git edit, one commit, verifier call, safe `issue/<hash>-restart-oracle` branch, no remote, and dirty-parent preservation;
- deterministic branch collision rejection with no duplicate worktree or commit effect;
- traversal-shaped issue input that stays inside the issue branch namespace and fail-closes verification for a rejected outcome;
- coder interruption salvage followed by restart with exactly one commit per independent branch and no duplicate effect.

All Git assertions use real local subprocess commands and exit statuses. The test replaces only the documented coding/verifier seam and never contacts a remote.

## Validation receipts

### Original QA file

```text
command: pytest tests/issue/test_luna_high_issue_git_restart_20261013.py
status: BLOCKED
code: GO_WAVE_NODE_ERROR:test_outside_task_write_scope
exit_code: null
note: The scoped runner rejected the read-only regression file before pytest execution; the identical denied call was not retried.
```

### Independent oracle

```text
command: .venv\Scripts\python.exe -m pytest -q tests/issue/test_native_git_branch_contract_20261013.py
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: 0998535604e429e0885dc2fa8e82d3efeed30071fdd7cc9bc351fe7bba576d5a
exit_code: null
stderr: A regular local .venv Python executable is required.
```

The test result is not promoted to PASS. No executable positive or negative pytest result is available in this attempt.

## Next bounded move

Materialize the regular project `.venv`, then run exactly:

```text
.venv\Scripts\python.exe -m pytest -q tests/issue/test_native_git_branch_contract_20261013.py
```

After a real pytest receipt exists, the original QA owner should update its branch assertion to the safe namespaced contract and independently review the new oracle. Product acceptance remains with the independent verifier/reviewer and coordinator; this artifact claims neither `ACCEPTED` nor `ACCEPTED_GREEN`.
