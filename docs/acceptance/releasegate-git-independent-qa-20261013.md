# Independent disposable Git release-gate QA — 2026-10-13

## Authority and verdict

- Independent QA task: `swe-af:v2-releasegate-independent-qa-20261013-r1`
- Producer task under review: `swe-af:v2-releasegate-disposable-git-releasegate-20261013-r1`
- Review scope is limited to the separate QA test and this evidence record.
- Producer files, shared fixtures, production source, CrewSync state, and Ruflo state were not edited.
- Verdict: **CANDIDATE / UNACCEPTED / REPAIR REQUIRED**.
- No product-terminal or scoped-acceptance recommendation is made because executable pytest evidence and authoritative original CrewSync readback are both unavailable.
- The restricted no-shell tool surface exposes file SHA receipts but no project Git `HEAD` or CrewSync status readback operation; neither identity is inferred from documents.

## Source-bound evidence

The producer executable artifact is `tests/issue/test_releasegate_git_20261013.py`.

- Producer test SHA-256: `8e9c356a188e43f003286fc07d85b4d4fcc1e6b252f0faca6b02ef2956adf0eb`
- Producer acceptance artifact: `docs/acceptance/releasegate-git-20261013.md`
- Producer acceptance artifact explicitly records that the final producer test SHA was never executed.
- The only recorded producer runner receipt is `PYTHON_VENV_UNAVAILABLE`, with no pytest exit code.
- Independent QA test: `tests/issue/test_releasegate_git_independent_qa_20261013.py`
- Independent QA test SHA-256: `3d56a509596e44eab7a4dd8f63543e4b18880e93119b8812f6af928cc0bd16b3`

The source SHA values above bind the reviewed producer artifact and the independent test submitted to the scoped runner. They do not substitute for a Git commit identity or a passing process result.

## Factual collection defect

The producer test imports a fixture helper with:

```text
from tests.issue.conftest import run_git
```

The repository has `tests/issue/__init__.py` but no `tests/__init__.py`. Therefore `tests` is not a regular local package and can be shadowed by an installed regular package named `tests`. The independent QA test contains a deterministic reproduction that supplies such a shadow package and expects the import to fail with `No module named 'tests.issue'`.

This defect belongs to the producer artifact. Independent QA did not patch the producer test or add `tests/__init__.py`.

## Independent executable challenge

The separate QA file defines its own Git helper and fixture instead of importing producer test support. Its real-Git scenario checks:

- a misleading reviewer approval cannot override an executable verifier rejection;
- the rejected branch has one commit and preserves the duplicate-dropping source;
- a simulated crash after the repair write is checkpointed without reaching reviewer or verifier;
- restart from the rejected branch produces one successful repair commit;
- rejected and repaired source blob SHAs differ;
- interrupted and restarted repair blob SHAs match;
- repaired source contains exactly one repair effect and no residual `set(items)` effect;
- all attempt branches are distinct;
- caller `HEAD`, current branch, and clean status remain unchanged;
- `main` does not gain the generated source;
- no remote is configured and temporary worktrees are removed.

This supplies an adversarial negative path and a crash/retry boundary without relying on the producer fixture import.

## Focused runner receipt

The authorized scoped runner was called once for `tests/issue/test_releasegate_git_independent_qa_20261013.py` and returned:

```text
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: 3d56a509596e44eab7a4dd8f63543e4b18880e93119b8812f6af928cc0bd16b3
exit_code: null
stderr: A regular local .venv Python executable is required.
```

The identical denied call was not retried. This is an infrastructure-blocked check, not PASS and not evidence that either independent test function collected or executed.

## Adversarial findings

1. The final producer SHA has no executable pytest receipt.
2. The producer's absolute `tests.issue.conftest` import is shadowable during collection.
3. The producer evidence documents an infrastructure block but cannot support a runtime success claim.
4. No authoritative original CrewSync terminal status or attempt readback exists in the permitted repository evidence.
5. No immutable producer `go-core.worker-result.v1` payload is available in the permitted read scope for field-by-field verification.
6. No process exit receipt is bound to the final producer test SHA.
7. Forty-character Git identity assertions assume SHA-1 repositories and do not generalize to SHA-256 object format.
8. Equal interrupted/restarted blob IDs prove equal content, not identical attempt ancestry or absence of all external duplicate effects.
9. Branch distinctness does not itself prove the repair attempt was idempotent outside the generated file.
10. Source-content assertions do not assert an exact Git diff against the rejected parent.
11. Temporary worktree cleanup is checked, but artifact-directory cleanup is not established.
12. Publication transport and typed-result verification are external to the producer test and are not proven by its Git assertions.

## Same-producer sticky repair request

Return the original producer task to the same producer/thread and require a bounded repair of `tests/issue/test_releasegate_git_20261013.py`:

1. Remove the collection dependency on `from tests.issue.conftest import run_git` by using a producer-local helper or another repository pattern that cannot be shadowed by a third-party `tests` package.
2. Run the repaired producer test in a materialized regular `.venv` and preserve the exact pytest exit code and output.
3. Run this independent QA test unchanged in the same environment.
4. Preserve Git `HEAD`, producer source SHA, process result, and the original typed worker-result readback in one immutable evidence chain.
5. Obtain authoritative coordinator/CrewSync status readback for the exact original task and attempt; do not infer status from process completion or publisher transport.

## Reducer next action

A separate reducer/reviewer should remain fail-closed until:

- the producer collection defect is repaired by the original producer;
- producer and independent focused tests execute successfully;
- the repaired producer artifact receives a new source SHA and exact test receipt;
- the original producer typed result and exact CrewSync status/attempt are independently read back;
- process success, publication status, and product acceptance remain explicitly distinct.
