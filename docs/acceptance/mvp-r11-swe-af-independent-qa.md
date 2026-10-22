# MVP R11 independent C4 QA receipt

## Scope and authority

This is an independent QA oracle for the frozen r10 C4 candidate at
`swe_af.issue.build._implement_issue_impl`. It changes no producer source and
does not modify the producer's r10 test. The only implementation artifact is
`tests/issue/test_northstar_mvp_r11_issue_qa.py`.

The QA test binds execution to the supplied frozen candidate revision:

```text
064a7226a500365005b5b662044fc9014a328722
```

It also computes the runtime SHA-256 of `swe_af/issue/build.py`. A mismatch of
the Git revision is a product/source `RED` condition, not a fixture pass.

## Independent oracle coverage

| Scenario | Expected proof |
| --- | --- |
| Positive C4 path | A rejected disposable issue returns `success=False`, then a separate completed issue returns `success=True`; both use the same caller base SHA and produce one commit on distinct disposable branches. |
| Negative: invalid input | A missing issue title raises before the coding loop runs and leaves the caller `HEAD` and working tree unchanged. |
| Negative: cancellation | Cancellation checkpoints dirty work only on the issue branch, propagates `CancelledError`, leaves the caller unchanged, and preserves the checkpointed file for inspection. |
| Crash/replay | A later replay runs on a fresh branch, produces exactly one commit, and has no caller-side file or status mutation. |
| Source binding | The checked-out project Git revision matches the frozen r10 receipt and the QA run computes the producer source SHA-256. |

The QA test uses a distinct temporary caller repository, filenames, issue names,
loop stubs, and assertions from r10. It is therefore a noncolliding oracle,
not a modification of the producer fixture.

## Required validation

```text
python -m pytest -q tests/issue/test_northstar_mvp_r11_issue_qa.py
python -m pytest -q tests/issue/test_northstar_mvp_r10_disposable_git_e2e.py
```

The authorized GO Wave test runner for this attempt accepts only exact
Node/Vitest/tsx tests. Python/pytest must be classified as infrastructure
blocked if that runner rejects either command; no denied command is promoted
to PASS and no repeated identical runner call is valid repair work.

## Evidence status

- Frozen r10 Git identity is embedded as a runtime assertion: **PENDING native pytest**.
- Positive red/green, two negatives, and cancellation/replay are encoded at the exact C4 boundary: **PENDING native pytest**.
- Producer source and original r10 test remain unchanged by this QA task: **PASS by write scope**.
- Product release and `ACCEPTED_GREEN`: **not claimed**.

## Reduction requirement

Only an independent reducer may combine a real focused pytest receipt, source
identity evidence, and the original producer result into a later acceptance
decision. This QA result remains `UNACCEPTED` regardless of a publishing
transport acknowledgement.
