# Longbundle issue full-cycle acceptance — 2026-10-14

## Decision

**Candidate / UNACCEPTED.** The scoped artifact adds source-backed C4 lifecycle assertions, but the authorized pytest runner is infrastructure-blocked because the regular project `.venv` is unavailable. No `ACCEPTED_GREEN` claim is made.

## B1 — Inspect and baseline

- Target Git root: `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af`.
- Task source of truth: `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`, whose prior decision remains `UNACCEPTED / BLOCKED`.
- `swe_af.issue.build.implement_issue` is the public async entrypoint with the six-parameter DTO-compatible signature: `issue`, `repo_path`, `base_branch`, `artifacts_dir`, `additional_context`, `config`.
- Source SHA for `swe_af/issue/build.py`: `aa69eb29670c2eef760a3d4becf4651cdf0af5dc4b6a14622f5caf6eccf4dc00`.
- Source SHA for `swe_af/issue/git_ops.py`: `0f73d9fb3a613f1b68ccf4a8161a93c1c5e6dd6a4a392b2c2c8a6157b7c2403d`.
- Existing adjacent evidence: `tests/issue/test_northstar_issue_worktree_repair_20261013.py` covers lower-level live Git worktree, duplicate branch, foreign base, stale branch, and retry-budget behavior.

## B2 — Reproduction and adversarial coverage

The new test uses disposable real Git repositories and the production issue-build implementation. Only `run_coding_loop` is controlled, isolating the provider boundary while exercising real branch creation, worktree cleanup, checkpoint commits, verifier gating, changed-file receipts, and parent-checkout isolation.

Covered assertions:

1. Positive completed loop with verifier pass, immutable base SHA, one branch commit, exact changed file, preserved dirty caller state, and DTO signature.
2. Negative rejected loop with no commit, no returned branch, unchanged parent HEAD, and unchanged status.
3. Negative verifier rejection with a retained branch/commit receipt but `success == False`.
4. Crash/restart: interrupted work is salvaged into one checkpoint commit; a fresh retry starts from the same base SHA, uses a distinct branch, produces one final commit, and does not inherit the interrupted file.

## B3 — Scoped implementation

Changed only the authorized files:

- `tests/issue/test_northstar_longbundle_issue_fullcycle_20261014.py`
- `docs/acceptance/longbundle-issue-fullcycle-20261014.md`

Final focused test artifact SHA: `4673e224a2a113bc196f050406ce271b421a36c35aba348a7d37a299076e109f`.

No production source change was necessary because the inspected implementation already contains the required deterministic local Git cycle and crash salvage behavior.

## B4 — Validation receipt

Requested command:

```text
python -m pytest -q -p no:cacheprovider tests/issue/test_northstar_longbundle_issue_fullcycle_20261014.py
```

Authorized runner receipt:

```text
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
test: tests/issue/test_northstar_longbundle_issue_fullcycle_20261014.py
source_identity_sha256: 334739f015cb5114c347139fc75dd335f409778cbe6e2cd41cc46a2f027bc168
exit_code: null
stderr: A regular local .venv Python executable is required.
```

The receipt predates the final small test-hardening edit and is therefore not promoted to evidence for the final SHA. The same denied runner call must not be retried until the regular project `.venv` exists.

## B5 — Handoff and remaining gate

- Focused runtime result for final test SHA: **BLOCKED**, not PASS.
- Affected regression result: **NOT RUN**, because the only authorized Python runner is unavailable.
- Product/C4 acceptance: **UNACCEPTED**.
- Independent reviewer must rerun the focused pytest command against final SHA `4673e224a2a113bc196f050406ce271b421a36c35aba348a7d37a299076e109f`, inspect the four lifecycle assertions, and then perform the reducer/terminal readback separately.
