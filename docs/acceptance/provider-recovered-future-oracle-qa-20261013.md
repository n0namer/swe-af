# Provider-recovered future-oracle QA — 2026-10-13

## Decision

**BLOCKED / UNACCEPTED.**

A new independent regression wrapper was added, but the authorized runner stopped before pytest because the regular project `.venv` Python executable is unavailable. The reported provider health recovery is therefore not treated as native test evidence. No `PASS`, `CANDIDATE`, `ACCEPTED`, or `ACCEPTED_GREEN` claim is made.

## Authority and source identity

| Receipt | Exact value | Status |
|---|---|---|
| Project | `swe-af` | trusted binding |
| Current task | `swe-af:v2-provider-recovery-postfailure-independent-qa-20261013-r1` | trusted binding |
| Current attempt | `attempt-01543fea-96c6-4d67-8ba2-90a76d9a27eb` | trusted binding |
| Crew session | `manual-20261013-provider-recovery-postfailure-w-369fb454e43b2106` | trusted binding |
| Git/source SHA | `064a7226a500365005b5b662044fc9014a328722` | prelaunch receipt; runtime comparison did not execute |
| Manifest SHA-256 | `f475ae55588cf7d13a9597eeab8a1e423301afbcf56fd3ca48dedb0c74eac79a` | prelaunch receipt |
| Current QA test SHA-256 | `704fb1d99fb83bd8e9c416be39c8b01ee80380039826cc828fac0388b70c2907` | authorized runner source identity |
| Prior future-oracle SHA-256 pin | `03119b752aadff57cbc18a45bef58cf52afc7fcfcc3c44f14322425dd548b6a8` | prior acceptance artifact |
| Prior native scenario SHA-256 pin | `8e9c356a188e43f003286fc07d85b4d4fcc1e6b252f0faca6b02ef2956adf0eb` | prior acceptance artifact |

## Prior failure and recovery context

The coordinator identified the failed predecessor as `swe-af:v2-native-future-oracle-runtime-independent-qa-20261013-r1`: physical Codex exit `1` after OpenCodex at `127.0.0.1:10110` returned `503` or timed out. Its attempt identifier was not supplied in the authorized evidence and is not invented here.

The active BMad source of truth, `docs/acceptance/native-future-independent-oracle-20261013.md`, records the earlier oracle owner `swe-af:v2-native-future-independent-oracle-20261013-r1`, attempt `attempt-42575b33-4a01-4a94-ba0b-a20c8fb997f4`, as runtime blocked by the same missing regular `.venv`.

The current task states that provider health has recovered to HTTP `200`. The new test checks that health receipt, but it also requires the native oracle to pass. A health response alone cannot satisfy the regression gate.

## New independent QA contract

The new test is `tests/issue/test_future_oracle_provider_recovery_20261013.py` and contains two focused tests:

1. `test_provider_health_is_recorded_but_not_treated_as_native_pass`
   - requests `http://127.0.0.1:10110/health` and requires HTTP `200`;
   - deliberately supplies an invalid prior-oracle SHA and requires the fail-closed drift assertion;
   - prints a structured health receipt only after both checks execute.
2. `test_native_future_oracle_positive_and_restart_replay`
   - requires Git `HEAD` to equal `064a7226a500365005b5b662044fc9014a328722`;
   - requires the prior future-oracle source SHA-256 to equal `03119b752aadff57cbc18a45bef58cf52afc7fcfcc3c44f14322425dd548b6a8`;
   - imports the real `swe_af.issue.build.implement_issue` boundary;
   - requires its implementation source to resolve below the product `swe_af/` tree;
   - records the exact runtime callable signature, parameter kinds, annotations, defaults, return annotation, source path, and source SHA-256;
   - executes the pinned prior future oracle once as the positive run and a second time as restart/replay;
   - requires exit code `0`, a positive pytest pass marker, and no skip, xfail, or xpass marker on both runs.

The prior future oracle is itself pinned to the reviewed native Git scenario, which the previous acceptance artifact says exercises the rejected result, interruption before review/verifier, persisted repair state, retry/restart, isolation, and cleanup. This attempt did not execute that chain, so those behaviors remain unverified here.

## C4 DTO and API compatibility

The exact real C4 callable and DTO evidence is designed to be emitted from the product runtime rather than guessed:

```text
qualified callable := swe_af.issue.build.implement_issue
signature := inspect.signature(implement_issue)
parameters := runtime names + kinds + annotation repr + default repr
return DTO boundary := runtime return annotation repr
implementation identity := product source path + SHA-256
source state := git rev-parse HEAD
```

Because pytest never started, the actual signature, annotation/DTO identifiers, and product source SHA-256 were not observed in this attempt. C4 compatibility is **INFRA BLOCKED**, not PASS.

## Exact validation command and result

Requested focused command:

```text
.venv\Scripts\python.exe -m pytest -q tests/issue/test_future_oracle_provider_recovery_20261013.py
```

Authorized scoped-runner receipt:

```text
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
test: tests/issue/test_future_oracle_provider_recovery_20261013.py
engine: python:pytest
source_identity_sha256: 704fb1d99fb83bd8e9c416be39c8b01ee80380039826cc828fac0388b70c2907
exit_code: null
signal: null
duration_limit_ms: 120000
stdout: <empty>
stderr: A regular local .venv Python executable is required.
```

The denied call was not retried. Provider HTTP `200` was not freshly observed because the runner blocked before test collection.

## Gate status

| Gate | Status | Evidence |
|---|---|---|
| Health is separate from native PASS | PASS — source contract only | The first test records health and independently requires a fail-closed negative gate. Runtime did not execute. |
| Real positive native path | INFRA BLOCKED | Prior oracle invocation did not start. |
| Real negative path | INFRA BLOCKED | Source-drift rejection and pinned native rejected-result path did not execute. |
| Restart/replay path | INFRA BLOCKED | The prior oracle was not run twice. |
| Exact Git/source SHA | PARTIAL | Trusted Git SHA and current test SHA are recorded; runtime Git and product source checks did not execute. |
| Real C4 DTO/API signatures | INFRA BLOCKED | Runtime introspection did not execute. |
| Writes limited to new QA paths | PASS | Only the two authorized paths were created. |
| Product acceptance | UNACCEPTED | Independent reducer and coordinator terminal readback remain required. |

## Recovery and independent QA

1. Restore or materialize the regular project `.venv` without changing this task's source ownership.
2. Resume this exact sticky task and run the focused command above once.
3. Preserve the complete `PROVIDER_RECOVERY_HEALTH_RECEIPT` and `PROVIDER_RECOVERY_NATIVE_RECEIPT` output, including both nested pytest exits and the observed C4 signature/source identity.
4. If a SHA pin differs, stop as source drift and require independent review before changing the pin.
5. Independent QA compares the observed parameter and return annotations, concrete request/return objects, and persisted receipt schema/version against producer source.
6. Resume the existing reducer only after real positive, negative, and replay receipts exist; coordinator terminal readback remains the acceptance authority.

## Self-review

- No production source, previous claimed test, prior acceptance artifact, task, claim, lease, lock, dependency, or remote state was modified.
- The new wrapper fails closed on missing health, Git drift, oracle drift, fake/non-product entrypoint origin, nested pytest failure, zero-pass output, skip, xfail, and xpass.
- The same native oracle is executed twice to make restart/replay evidence explicit when infrastructure is available.
- Remaining risk is fully runtime-bound: the regular `.venv` is absent, so no native behavior or C4 signature was observed.
