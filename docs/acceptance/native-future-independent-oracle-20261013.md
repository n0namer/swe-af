# Native future independent oracle — 2026-10-13

## Decision

**CANDIDATE / UNACCEPTED / runtime BLOCKED.**

This attempt adds an independent, fail-closed future oracle for the native ReleaseGate Git scenario. The oracle does not use a fake request handler, mocked receipt, in-memory persistence substitute, skip, or xfail as acceptance evidence. Its authorized execution was blocked before pytest started because the regular project `.venv` Python executable is unavailable.

No `PASS`, `ACCEPTED`, or `ACCEPTED_GREEN` claim is made.

## Authority and source identity

| Receipt | Exact value | Status |
|---|---|---|
| Project | `swe-af` | PASS — trusted binding |
| Task | `swe-af:v2-native-future-independent-oracle-20261013-r1` | PASS — trusted binding |
| Attempt | `attempt-42575b33-4a01-4a94-ba0b-a20c8fb997f4` | PASS — trusted binding |
| Crew session | `manual-20261013-native-future-independent-qa-wa-35a1cd22c32367b5` | PASS — trusted binding |
| Git identity | `064a7226a500365005b5b662044fc9014a328722` | PASS — prelaunch receipt |
| Manifest SHA-256 | `c3effa82c6988b63d35355465538bd2856c2ea4782a2c439d93f9378295b2899` | PASS — prelaunch receipt |
| Accepted evidence cache SHA-256 | `93970bb38e9c1d5168cebac857d8fe9397da628c3a7ee276b9eba11b702106d4` | PASS — trusted task data; contents not expanded beyond read scope |
| New oracle source identity SHA-256 | `03119b752aadff57cbc18a45bef58cf52afc7fcfcc3c44f14322425dd548b6a8` | PASS — final authorized runner receipt |
| Pinned native scenario SHA-256 | `8e9c356a188e43f003286fc07d85b4d4fcc1e6b252f0faca6b02ef2956adf0eb` | PASS — prior C4 reducer evidence; runtime remains blocked |

## Independent oracle contract

The new test is `tests/issue/test_releasegate_future_git_qa_20261013.py`.

It establishes these measurable gates:

1. Import the real public boundary `swe_af.issue.build.implement_issue`.
2. Resolve its source file and require that it is under the product `swe_af/` tree, rejecting test or fake adapters.
3. Capture the callable's exact runtime signature using `inspect.signature`; this is the C4 request-signature boundary and avoids inventing parameter or DTO names that were not exposed in this attempt's read scope.
4. Require the existing native scenario `tests/issue/test_releasegate_git_20261013.py` to match reviewed SHA-256 `8e9c356a188e43f003286fc07d85b4d4fcc1e6b252f0faca6b02ef2956adf0eb` before execution.
5. Execute that pinned scenario with the regular project `.venv` and real pytest engine. The prior C4 reducer records that this scenario covers native request, rejected result, interruption before review/verifier, persisted repair state, restart/retry, isolation, and cleanup.
6. Require a non-zero pytest pass count and reject skip, xfail, or xpass output. Exit zero without an actual pass count is not GREEN.
7. On failure, include the exact inspected public signature and subprocess output in the assertion receipt.

### C4 DTO boundary

The exact producer DTO class/version identifiers and exact `implement_issue` parameter list are **BLOCKED** by the authorized read scope and unavailable Python runtime. They are not invented here.

The boundary is instead source-grounded and executable:

```text
request signature := inspect.signature(swe_af.issue.build.implement_issue)
implementation origin := a file below <project>/swe_af/
scenario identity := sha256(tests/issue/test_releasegate_git_20261013.py)
receipt := real pytest subprocess exit + stdout + stderr + positive pass count
acceptance exclusions := skip, xfail, xpass, fake/test adapter, SHA drift, missing .venv
```

A future real run exposes the exact callable signature in any failure receipt. Independent QA must then compare that observed signature and the producer's concrete return/persisted receipt objects with the versioned C4 DTO declarations from the producer source. Until that comparison exists, DTO compatibility remains BLOCKED.

## Positive, negative, and recovery evidence

| Gate | Status | Evidence |
|---|---|---|
| Real product entrypoint boundary | BLOCKED | Test implemented; import/signature check did not execute because `.venv` is unavailable. |
| Positive native request to actual receipt | BLOCKED | Delegated only to the pinned real-Git scenario; no runtime receipt exists. |
| Rejected implementation remains unsuccessful | BLOCKED | Covered by the pinned scenario contract recorded by the prior reducer; not executed here. |
| Interruption before reviewer/verifier | BLOCKED | Covered by the pinned scenario contract; not executed here. |
| Persisted repair survives restart | BLOCKED | Covered by the pinned scenario contract; not executed here. |
| Retry/restart is idempotent and reaches verifier | BLOCKED | Covered by the pinned scenario contract; not executed here. |
| Caller/main isolation and cleanup | BLOCKED | Covered by the pinned scenario contract; not executed here. |
| Fake or skip-only GREEN rejected | PASS — source review only | New oracle explicitly rejects non-product adapter origin, SHA drift, absent pass count, skip, xfail, and xpass. Runtime behavior remains unexecuted. |
| Broad regression | SKIP | Out of scope; this task authorizes one focused oracle only. |

## Reproducible command and actual receipt

Intended focused command:

```text
.venv\Scripts\python.exe -m pytest -q tests/issue/test_releasegate_future_git_qa_20261013.py
```

Authorized runner receipt:

```text
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
test: tests/issue/test_releasegate_future_git_qa_20261013.py
engine: python:pytest
source_identity_sha256: 03119b752aadff57cbc18a45bef58cf52afc7fcfcc3c44f14322425dd548b6a8
exit_code: null
signal: null
duration_limit_ms: 120000
stdout: <empty>
stderr: A regular local .venv Python executable is required.
```

This denied call was not retried.

## Future-ready affected-only dependency graph

No duplicate producer, QA, repair, or reducer is proposed.

1. **Original ReleaseGate repair/producer owner — coordinator-resolved existing task ID.** Restore/materialize the regular project `.venv`; do not create a replacement owner. Run the pinned original scenario first if its source remains at the reviewed SHA.
2. **This exact independent-oracle owner/task** `swe-af:v2-native-future-independent-oracle-20261013-r1`. Resume the same sticky task after infrastructure restoration and run the command above. If the pinned scenario SHA changed, stop on `NATIVE_SCENARIO_DRIFT` and require independent source review before updating the pin.
3. **Original independent-QA owner — coordinator-resolved existing task ID.** Compare the runtime-observed `implement_issue` signature, concrete request/return DTOs, persisted receipt schema/version, and restart receipt with producer source. Review positive and all failure/recovery assertions from actual output.
4. **Existing C4 reducer** `swe-af:v2-c4-releasegate-independent-reducer-20261013-r6`. Resume only after real producer and independent-QA receipts exist.
5. **Coordinator terminal readback.** Required after independent review; worker publication transport is not product acceptance.

The exact original producer and independent-QA task IDs were not present in the authorized evidence. The coordinator must recover those existing IDs from CrewSync. Inventing replacement IDs would violate affected-only repair ownership.

## Measurable Definition of Done for the next READY run

- Regular project `.venv` Python exists and launches pytest.
- New oracle source identity is recorded for the exact run.
- Both tests in the new oracle execute; no skip/xfail/xpass occurs.
- Public C4 entrypoint resolves to product source and its exact signature is captured.
- Pinned native scenario SHA matches or receives a new independent source review before pin change.
- Nested native scenario returns exit code `0` with a positive pass count.
- Actual output proves the scenario's positive, rejection, interruption, persistence, restart/retry, isolation, and cleanup assertions.
- Independent QA verifies exact request, return, and persisted receipt DTO names/schema versions against producer source.
- Existing reducer and coordinator perform separate acceptance and terminal readback.

## Self-review

- Changes are limited to the two authorized paths.
- No production code, current owned test, task state, lease, claim, dependency, or remote was changed.
- The oracle is fail-closed on missing infrastructure, source drift, fake adapter origin, nested pytest failure, zero-pass output, and non-executed pytest outcomes.
- Remaining risk: the exact producer signature and DTO identifiers could not be observed or executed in this attempt; compatibility is explicitly BLOCKED rather than guessed.
