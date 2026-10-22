# Northstar Git repair C4 QA — 2026-10-13

## Decision

**BLOCKED / UNACCEPTED / INFRASTRUCTURE BLOCKED.**

The scoped acceptance suite is implemented, but the authorized runner could not execute it because the regular project `.venv` Python executable is unavailable. No runtime PASS or product acceptance is claimed.

## Authority and identities

| Receipt | Exact value |
|---|---|
| Project | `swe-af` |
| Task | `swe-af:v2-c4-native-acceptance-git-repair-20261013-r1` |
| Attempt | `attempt-41f1689d-f59f-4e51-a735-ac81d2bd032e` |
| Crew session | `manual-20261013-c4-acceptance-provider-retryfro-f2f0e360e2b398e8` |
| Prelaunch Git identity | `064a7226a500365005b5b662044fc9014a328722` |
| Focused test source identity SHA-256 | `a7bf079ec1b31ec63e77d7d98ae212c7469f4c34f700c7e22c99df678e6c4413` |
| Production source boundary | `swe_af/issue/git_ops.py` (read-only in this attempt) |
| Test artifact | `tests/issue/test_northstar_issue_worktree_repair_20261013.py` |

## C4 function boundary

**FACT:** The exercised native boundary is the existing synchronous `swe_af.issue.git_ops` API:

- `resolve_base(repo_path, base_branch="") -> tuple[str, str]`
- `add_worktree(repo_path, worktree_path, branch, base_sha, attempts=3) -> None`
- `commit_all(worktree_path, message) -> str`
- `new_commits(repo_path, base_sha, branch) -> list[str]`
- `changed_files(repo_path, base_sha, branch) -> list[str]`
- `remove_worktree(repo_path, worktree_path) -> None`
- `delete_branch(repo_path, branch) -> None`
- failure DTO: `GitOpsError` with an attempt-count-bearing message from `add_worktree`

**FACT:** No versioned DTO object exists at this boundary; values are Python primitives plus `GitOpsError`. The test preserves that existing signature and does not introduce a competing C4 interface.

## Acceptance scenarios

| Scenario | Native or seam | Required invariant | Runtime status |
|---|---|---|---|
| Positive worktree and commit | Native Git subprocesses in a disposable repository | Branch starts at exact base SHA; exactly one final commit; only `solution.txt` changes; no remote; dirty parent status/content preserved; worktree removed | BLOCKED — not executed |
| Duplicate branch/worktree | Native Git subprocesses | Deterministic collision fails after one attempt despite a larger budget; second worktree absent; dirty parent preserved | BLOCKED — not executed |
| Foreign base SHA | Native Git subprocesses across two disposable repositories | Unknown foreign commit is rejected after one attempt; no branch/worktree created; dirty parent preserved | BLOCKED — not executed |
| Stale interrupted branch and restart | Native Git subprocesses | Stale branch prevents blind retry; explicit cleanup permits restart; interrupted file is absent; restart produces exactly one final commit | BLOCKED — not executed |
| Transient retry budget | Explicit monkeypatched unit seam, not labelled live | Three retryable failures consume exactly three attempts and only two backoff sleeps (`0.5`, `1.0`) | BLOCKED — not executed |

## Runner evidence

Requested focused validation:

```text
.venv\Scripts\python.exe -m pytest -q tests/issue/test_northstar_issue_worktree_repair_20261013.py
```

Authorized runner receipt:

```text
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: a7bf079ec1b31ec63e77d7d98ae212c7469f4c34f700c7e22c99df678e6c4413
exit_code: null
stderr: A regular local .venv Python executable is required.
```

**FACT:** The focused suite did not run. The identical denied call was not retried.

Affected regression command, intentionally not invoked after the same engine prerequisite was proven absent:

```text
.venv\Scripts\python.exe -m pytest -q tests/issue/test_git_ops.py tests/issue/test_northstar_issue_worktree_repair_20261013.py
```

Status: **BLOCKED by the same `PYTHON_VENV_UNAVAILABLE` prerequisite; no exit code and no PASS claim.**

## Edge-case and self-review

**FACT:** All four integration scenarios create disposable repositories with local commit identity and no remote. Tests never commit in the project checkout, contact a provider, restart a service, or mutate CrewSync/Ruflo state.

**FACT:** Parent dirt is captured before native operations and compared byte-for-byte afterward in the positive, duplicate, and foreign-base paths.

**FACT:** Deterministic duplicate/stale failures assert `failed after 1 attempts`, while the retryable lock-contention seam asserts exhaustion of the configured budget. This distinguishes non-retryable intent collisions from transient retry behavior.

**FACT:** The restart path removes the interrupted worktree, observes the stale branch collision, performs explicit branch cleanup, recreates from the original base SHA, and asserts one final commit.

**ASSUMPTION:** Git is available to pytest because the production implementation itself shells out to `git`; this remains unverified until the blocked command runs.

## Root-cause-only repair route

No production defect is proven while execution is blocked. Do not modify `swe_af/issue/git_ops.py` preemptively.

1. Restore the regular project `.venv` outside this worker's write scope.
2. Resume this exact task/attempt lineage for the focused command rather than creating a duplicate acceptance task.
3. If the focused test reports a semantic failure in `git_ops.py`, route only that reproduced defect to the original coordinator-owned CrewSync production-repair owner for `swe_af/issue/git_ops.py`; this worker has no authority to claim or edit that file.
4. Return the production patch to this same acceptance lineage, rerun the focused suite, then run the affected regression command.
5. Require independent reviewer evidence and authoritative coordinator readback before any GREEN decision.

The permitted evidence does not expose a separate historical production-owner task key, so none is invented here. The coordinator must resolve and reuse the existing owner from CrewSync.

## Measurable DoD

- Actual native positive path defined with exact base, one-commit, changed-file, isolation, dirty-parent, and cleanup assertions: **IMPLEMENTED / RUNTIME BLOCKED**.
- At least two adversarial native negatives including stale/restart: **IMPLEMENTED / RUNTIME BLOCKED**.
- Interrupted retry budget: **IMPLEMENTED as an explicitly labelled unit seam / RUNTIME BLOCKED**.
- Exact source SHA, runnable command, and honest exit status: **PASS** (`exit_code: null`, not run).
- Existing C4 function/error boundary documented without source refactor: **PASS**.
- Root-cause-only repair route through original CrewSync ownership: **PASS**.
- Focused and affected regression runtime evidence: **BLOCKED**.
- Independent reviewer and coordinator acceptance: **REQUIRED; not performed or claimed**.
