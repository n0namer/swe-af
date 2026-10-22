# MVP independent audit — native issue disposable repair (2026-10-14)

## Verdict

**CANDIDATE / UNACCEPTED.** Independent QA supplied a fresh host-native runtime receipt with the focused MVP restart test and the independent Git replay oracle both passing. This artifact records that evidence without converting producer exit 0 or test success into `ACCEPTED_GREEN`; coordinator reduction and authoritative terminal readback remain independent.

The authorized `go_wave` runner in this worker remains infrastructure-blocked because its required regular project `.venv` is unavailable. That local runner limitation is distinct from the supplied host-native receipt and is not promoted to PASS.

## Task binding

| Field | Value |
|---|---|
| Project | `swe-af` |
| Task key | `swe-af:v2-mvp-independent-audit-native-issue-disposable-repair-20261014-r8` |
| Attempt | `attempt-a2b17d51-a1ab-42dc-bbd6-d36bf504609c` |
| Crew session | `manual-20261014-basefive-mvp-10owners-r8-swe-af-295166bd0013c556` |
| Sticky thread | `01a137d5-c074-7e03-b095-1d8093cd6b00` |
| Producer thread | `null` in the trusted binding |
| Producer physical result | Exit 0 reported in r8 launch-results; not equivalent to acceptance |
| Source under audit | `swe_af/issue/build.py` |
| Source SHA-256 | `2F727B2D10C47888368D007A76DA03B1BAD98CC295052B447E3E50F2F00418BF` |
| MVP test SHA-256 | `CA3779FC30F9FBFC40980801D5D0E3951FC2CB601BA09475B9AB99E3DCAC2044` |
| Active BMad source of truth | `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md` |

## Scope and safety

- Only this reserved acceptance document is written by this task.
- Production code and producer tests were not modified.
- No native task, lease, claim, terminal state, remote, branch, or deployment state was mutated.
- This is a QA candidate only; it does not self-accept or claim `ACCEPTED_GREEN`.

## Evidence inspected

| Item | Result |
|---|---|
| C4 reducer artifact | Read previously; it requires runtime proof, C4 compatibility proof, original producer/QA receipts, and coordinator terminal readback. |
| `swe_af/issue/build.py` | Previously inspected. `implement_issue` uses isolated worktrees, checkpoint salvage on interruption, optional verification, and cleanup; verifier exceptions are fail-closed. |
| MVP restart test | Independent QA reports exact test SHA `CA3779FC30F9FBFC40980801D5D0E3951FC2CB601BA09475B9AB99E3DCAC2044`. |
| Independent oracle | `tests/issue/test_northstar_independent_gitreplay_oracle_20261014.py`; source identity from the authorized runner is `9145667b831a04d604488a4d8d20b795ece8b494e4b9f4c6cdc854f213f3f198`. |
| Review coverage | Independent QA reports source/C4 signature inspection without test overfitting, positive behavior, two negative paths, crash/restart, Git status, and external-effect checks. |

## D1–D5 audit reduction

| Gate | Status | Evidence-backed result |
|---|---|---|
| D1 exact source and C4 contract | CANDIDATE | Independent QA reports source SHA `2F727B2D...00418BF` and C4 signature inspection. Exact producer/QA typed receipts still require coordinator-side verification. |
| D2 positive durable scenario | CANDIDATE | Host-native receipt reports the combined focused command completed with 8 passed in 8.80 seconds and exit 0; QA reports one durable isolated Git effect and no caller/main mutation. |
| D3 wrong identity / malformed output / cancellation fail-closed | CANDIDATE | QA reports two negative paths and source-level fail-closed behavior. The authorized worker runner did not execute due missing `.venv`, so this remains supplied independent evidence. |
| D4 crash/restart no duplicate external effect | CANDIDATE | QA reports crash-restart replay and no duplicate persisted external effect; the independent oracle covers checkpoint replay with no second commit. |
| D5 focused and affected regression evidence | CANDIDATE | Supplied host-native command passed both named tests. Exact log SHA is recorded below. Coordinator still owns final reduction and terminal readback. |

## Supplied host-native runtime receipt

```text
host: fresh host-native Python 3.12+AgentField SDK environment
command: python -m pytest -q -p no:cacheprovider tests/issue/test_mvp_native_issue_restart_20261014.py tests/issue/test_northstar_independent_gitreplay_oracle_20261014.py
status: PASS (supplied independent QA receipt)
summary: 8 passed in 8.80 s
exit_code: 0
log: D:/Temp/basefive-mvp10-r8-swe-native-independent-20261014.log
log_sha256: DA24B8D8FE70BDE9DDC69CFF5552244551E61BDFEA7F05BA0801D5CC2BE2D7DA
```

This receipt is independent QA evidence supplied to the sticky repair. It is not a claim that the restricted `go_wave` runner executed successfully.

## Authorized worker-runner receipts

### MVP restart test

```text
runner: go_wave.run_project_tests
file: tests/issue/test_mvp_native_issue_restart_20261014.py
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: ca3779fc30f9fbfc40980801d5d0e3951fc2cb601ba09475b9ab99e3dcac2044
exit_code: null
stderr: A regular local .venv Python executable is required.
```

### Independent Git replay oracle

```text
runner: go_wave.run_project_tests
file: tests/issue/test_northstar_independent_gitreplay_oracle_20261014.py
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: 9145667b831a04d604488a4d8d20b795ece8b494e4b9f4c6cdc854f213f3f198
exit_code: null
stderr: A regular local .venv Python executable is required.
```

The unavailable worker runner is not retried or relabeled as PASS. The supplied host-native receipt is preserved separately.

## Self-review

- Producer physical exit 0 is explicitly not treated as `ACCEPTED_GREEN`.
- Host-native test success is labeled supplied independent QA evidence, not local execution evidence.
- The prior invalid publication is not described as verified; it returned `INVALID_REPORT` because `recorded_at` was not an ISO timestamp.
- This repair uses the new attempt binding and retains the exact project/task/session identities.
- No production or producer-test files were changed.
- Product acceptance remains unclaimed.

## Remaining blockers and next move

1. Coordinator/reducer must independently verify the supplied producer and QA receipts, C4 DTO/revision compatibility, and source/test SHA bindings.
2. The coordinator must perform authoritative CrewSync terminal readback; worker publication is transport evidence only.
3. If local reproduction is required, restore the regular project `.venv` before rerunning the exact pytest command through the authorized runner.
