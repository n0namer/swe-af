# Base20 independent contract repair QA — 2026-10-13

## Gate and authority

- Task: `swe-af:v2-base20-independent-test-contract-repair-20261013-r1`
- Attempt: `attempt-2608c463-fef6-417c-bea2-57466c0c2bd5`
- Crew session: `manual-20261013-five-modules-targeted-repair-qa-73e66ce900e3a642`
- Project source of truth: the active repository instructions and existing Base20 acceptance packets; no competing plan or production source change was created.
- Scope: reviewer-owned QA artifacts only. No production, migration, deployment, branch, lease, scheduler, or CrewSync state was changed.

## Defect classification

### Python receipt helper

The prior QA source SHA was `dcab3e13ec20b67691ccde2a03c472efcd23da08c08762362f09e83d2da89a48`.

The recorded producer run had eight passing and two failing cases. Both failures used `verify: false`, for which the implementation result contains `verification: None`. The helper evaluated `result.get("verification", {}).get("passed")`; the default only applies when the key is absent, so an explicit `None` caused `AttributeError` before the intended RED receipt could be evaluated.

The repair treats verifier evidence as passing only when it is a mapping whose `passed` member is exactly `True`. Missing, null, malformed, false, or truthy non-boolean evidence remains fail-closed and cannot be converted into acceptance.

Final QA source SHA: `316292be89f16e2d5367ef5562d10e9867eb73b23e7b9e96fa34d1758ac83e52`.

### Go DAG/orchestration independent test

The prior QA source SHA was `4b8d3e5567171b4b3a031378c5a5442af10301b3ad2f3336676a996d501f2512`.

The file declared `package dag` to exercise unexported checkpoint persistence functions, while importing `internal/orch` to call `ResumeBuildHandler`. Production `orch` imports production `dag`, so the test introduced the test-only cycle `dag -> orch -> dag` even though the production packages remain separately testable.

The repair removes the direct `orch` import. DAG checkpoint assertions remain same-package white-box tests. The orchestration assertion is injected as a temporary same-package `orch` test through the Go tool's overlay mechanism and executed by a nested focused `go test`. This retains the real `ResumeBuildHandler`, exact duplicate-dispatch observation, raw terminal receipt equality, `resume=true` input assertion, checkpoint fixture, and fail-closed terminal evidence check without mutating production files or the checkout.

Final QA source SHA: `845fd4e25943fef3dad6b7437e64325f4e001c24d88e4ac2c059c2411c93b94c`.

## Focused execution receipts

### Python

Authorized runner request:

```text
go_wave.run_project_tests({file:"tests/issue/test_base20_independent_issue_review_qa.py"})
```

Observed result before the scoped repair:

```text
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: dcab3e13ec20b67691ccde2a03c472efcd23da08c08762362f09e83d2da89a48
exit_code: null
stderr: A regular local .venv Python executable is required.
```

The identical denied runner was not retried after editing. Required deterministic command remains:

```text
python -m pytest -q tests/issue/test_base20_independent_issue_review_qa.py tests/issue/test_implement_issue.py
```

### Go

Authorized runner requests were made before and after the scoped repair:

```text
go_wave.run_project_tests({file:"go/internal/dag/base20_independent_checkpoint_resume_qa_test.go"})
```

Both reached the configured `golang:1.24.0-alpine` engine but invoked the package from `/workspace`, above the module rooted at `/workspace/go`:

```text
status: FAIL
engine: docker:go-test
exit_code: 1
stderr: go: cannot find main module, but found .git/config in /workspace
package: ./go/internal/dag
```

Source identities were:

- before repair: `4b8d3e5567171b4b3a031378c5a5442af10301b3ad2f3336676a996d501f2512`;
- after repair: `845fd4e25943fef3dad6b7437e64325f4e001c24d88e4ac2c059c2411c93b94c`.

The failure occurs before Go parses or executes the package, so it is runner/module-root infrastructure evidence, not a product or test PASS. Required deterministic commands remain:

```text
cd go && go test ./internal/dag -run Base20Independent -count=1
cd go && go test ./internal/dag ./internal/orch ./internal/node -count=1
```

## Positive and adversarial contract retained

- Two safe Python tasks still require deliberate RED, bounded same-task repair, reviewer rerun, executable verifier success, frozen source SHA, branch receipt, and a clean repository.
- Reviewer approval cannot override verifier failure.
- A moving source branch remains RED even when the producer reports success.
- Interrupted or budget-exhausted work remains recoverable RED evidence and cannot duplicate commits.
- DAG recovery excludes an accepted source effect, selects the ambiguous verification step, rejects stale/lower-level and moving-source checkpoint writes, and retains the original source revision.
- Orchestration still exposes duplicate resume dispatch, missing checkpoint-version enforcement, and missing writer-fencing enforcement as producer-owned repair boundaries.
- Version mismatch, absent version, moving-source terminal evidence, and timeout verification remain fail-closed.

## Self-review and remaining owner gap

The changed files are limited to the exact manifest write scope. No acceptance is inferred from a process exit, mock response, transport status, or publication result.

Validation is infrastructure-blocked until the authorized Python runner provides a regular `.venv` and the Go runner executes from the `go` module root. The original producer still owns any real product defect exposed after those commands execute; this repair changes only independent QA fixture mechanics. Independent reviewer, reducer, and exact CrewSync terminal readback remain required before acceptance.
