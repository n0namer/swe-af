# Affected repair — native Git branch fixtures — 2026-10-14

## Scope and source of truth

- Project: `swe-af`
- Task: `swe-af:v2-affected-git-branch-fixtures-repair-20261014-r1`
- Attempt: `attempt-8edb2b87-3023-4720-8462-7fab7a3b726b`
- Session: `manual-20261014-base5-affected-qa-repair-wave-r-847d6952a81ff4fc`
- Governing artifact: `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`
- C4 fixture contract: `tests/issue/test_native_git_branch_contract_20261013.py`

## Root cause and bounded change

The fixture contained two stale expectations, not a product implementation defect:

1. The branch contract already matched `issue/<hex>-restart-oracle`, but one assertion still required the obsolete `/restart-oracle` substring.
2. A branch checked out in a linked worktree is rendered by `git branch --list` with a `+` marker; the fixture expected the unmarked branch name.

Only `tests/issue/test_native_git_branch_contract_20261013.py` was changed:

```diff
-    assert "/restart-oracle" in result["branch"]
+    assert result["branch"].endswith("-restart-oracle")

-        assert _git(repo, "branch", "--list", branch).stdout.strip() == branch
+        assert _git(repo, "branch", "--list", branch).stdout.strip() == f"+ {branch}"
```

No product source, ownership/fencing behavior, path validation, restart logic, or remote behavior was changed.

## Source and changed-file evidence

- `swe_af/issue/git_ops.py` SHA-256: `0f73d9fb3a613f1b68ccf4a8161a93c1c5e6dd6a4a392b2c2c8a6157b7c2403d`
- `tests/issue/test_native_git_branch_contract_20261013.py` SHA-256 after repair: `bda35109965759d0d2e79a5120e706f74522be9dd608fe2b92f2acd1b971a2b3`
- Changed files are limited to the authorized write scope:
  - `tests/issue/test_native_git_branch_contract_20261013.py`
  - `docs/acceptance/affected-repair-git-fixtures-20261014.md`

## R1–R5 evidence

- **R1 Reproduce:** authorized native pytest runner invoked for `tests/issue/test_native_git_branch_contract_20261013.py`; runner returned `PYTHON_VENV_UNAVAILABLE`, `exit_code: null`, before test execution.
- **R2 Isolate:** direct fixture/source inspection identified assertion drift; `git_ops.add_worktree` explicitly preserves deterministic collision failure and linked-worktree semantics are represented by the `+` marker.
- **R3 Patch:** corrected only the two fixture expectations above; no weakening of negative, owner/fencing, dirty-parent, verifier, or restart assertions.
- **R4 Verify:** post-edit bounded source read confirms both corrections and all affected safety/restart assertions remain present. A second identical pytest retry was not made because the runner denial is infrastructure, not a code failure.
- **R5 Handoff:** candidate remains `UNACCEPTED`; independent QA must run the required pytest command after a regular project `.venv` is available.

## Focused validation receipt

```text
command: python -m pytest -q -p no:cacheprovider tests/issue/test_native_git_branch_contract_20261013.py
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: 0998535604e429e0885dc2fa8e82d3efeed30071fdd7cc9bc351fe7bba576d5a
exit_code: null
stderr: A regular local .venv Python executable is required.
```

## Remaining gate

The implementation candidate is ready for independent QA, but product acceptance is not claimed. The exact focused pytest command must execute with exit code `0`; the current environment lacks the runner-required regular `.venv`.
