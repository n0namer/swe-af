# Disposable Git release gate — 2026-10-13

## Scope

This acceptance artifact covers the public `swe_af.issue.build.implement_issue` entrypoint using only a disposable local Git repository and scripted role responses. It does not contact a remote, push, create a pull request, install dependencies, or modify production source.

The executable artifact is:

- `tests/issue/test_releasegate_git_20261013.py`

## Bounded scenario

The test defines exactly two logical issues:

1. `add-stable-normalizer` creates a real source file on an isolated issue branch. The reviewer permits executable verification, and the verifier rejects the branch because the implementation drops duplicates.
2. `repair-stable-normalizer` replaces the faulty implementation. Its first run is interrupted after the source write, proving checkpoint salvage. The same bounded repair is then restarted from the rejected branch and passes the executable verifier.

The interruption and restart are two attempts of the same repair issue, not a third logical issue.

## Positive evidence asserted

- The public `implement_issue` wrapper dispatches coder, reviewer, and verifier roles.
- Real Git worktrees, branches, commits, source diffs, and cleanup are exercised.
- The rejected branch remains independently inspectable with `success == false`.
- Interrupted uncommitted repair work is checkpointed on a retained branch.
- The restarted repair passes an executable Python behavior probe.
- Rejected, interrupted, and repaired branch commit IDs are asserted to be 40-character Git SHAs.
- Source blob IDs are asserted to be 40-character Git SHAs.
- The rejected blob differs from the repaired blob.
- The interrupted repair blob equals the successful restarted repair blob, proving the replay produced the same source state.

## Negative and safety evidence asserted

- The rejected implementation is never reported successful.
- The interrupted attempt never reaches reviewer or verifier.
- The repair contains one `return sorted(items)` effect and no residual `set(items)` effect.
- The three attempt branches are distinct.
- The caller remains on `main` at its original HEAD with a clean working tree.
- `main` does not contain the generated source file.
- The disposable repository has no Git remote.
- Temporary issue worktrees are removed.
- No existing producer or QA file is overwritten.

## Reproduction

Run the focused test from the project root with the repository's regular virtual environment:

```text
.venv\Scripts\python.exe -m pytest -q tests/issue/test_releasegate_git_20261013.py
```

The test is deterministic apart from generated branch/build identifiers. It validates exact Git object identities at runtime rather than embedding machine-specific SHAs in this document.

## Current runner receipt

The authorized scoped runner was invoked once for `tests/issue/test_releasegate_git_20261013.py` and returned:

```text
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: 0ddb1e772d3fac42a7c0ce8fa03029871485c8002869b4e8be81a5391e981b00
exit_code: null
stderr: A regular local .venv Python executable is required.
```

Per the task policy, the identical denied runner call was not retried. This is infrastructure-blocked validation, not a passing test claim.

A subsequent source-only adversarial review replaced the verifier subprocess's generic `python` command with the active `sys.executable` and strengthened base-SHA/one-commit assertions. The final test artifact SHA-256 is:

```text
8e9c356a188e43f003286fc07d85b4d4fcc1e6b252f0faca6b02ef2956adf0eb
```

That final SHA has not been executed because the same required `.venv` remains unavailable. The earlier blocked receipt is retained exactly and is not represented as evidence for the post-review SHA.

## Handoff

1. Materialize the project's regular `.venv` without changing the test artifact.
2. Run the focused command above.
3. Preserve the pytest exit code and failure output, if any.
4. If the focused test passes, independently review the three generated branch/source SHA assertions and the no-remote/no-main-mutation checks.
5. Keep acceptance `UNACCEPTED` until independent reducer/reviewer and exact CrewSync terminal readback complete.
