# R12 issue Git red/green — R13 sticky-resume evidence

## Disposition — 2026-10-14

**BLOCKED / UNACCEPTED.** No test executed, no native provider or C4 compatibility claim is supported, and no ACCEPTED_GREEN claim is made. Independent QA, reducer review and coordinator terminal readback remain required.

Target: `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af`, not canonical BMAD-MNNZ. Existing BMad source of truth: `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`; its reducer decision remains UNACCEPTED / BLOCKED. This file records scoped attempt evidence, not a competing plan or acceptance authority.

## Exact attempt and source binding

- Task: `swe-af:v2-r13-issue-git-redgreen-20261014-r1`
- Attempt: `attempt-64e37825-9cdc-4d64-b43d-5a5b1736c2d6`
- Session: `manual-20261014-basefive-r13-sol-low-swe-af-v2--5d20eda6488c31d2`
- Producer thread: null; role: implementer.
- Git identity from trusted prelaunch, not independently inspected: `064a7226a500365005b5b662044fc9014a328722`.
- Manifest SHA-256 from trusted prelaunch: `298b30493572dd2bd10df97368262b3c45edf3cddfa3641f01f45181bd3ac587`.
- Current `swe_af/issue/build.py` SHA-256 from scoped range receipt: `440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd`.
- New focused test SHA-256 from scoped read and runner: `aa0c1783c4621258818b5b6c0350100c07128b22e36f434beef68561c7cbc6bf`.
- C4 artifact prelaunch SHA-256: `feeb1fe6a88e53770831de890fe3674e6425dd40327deb913fbf73f3a01f1f17`.

## Inspection and smallest compatible delta

FACT: The authorized test and this evidence artifact initially returned `GO_WAVE_NODE_ERROR:ENOENT`; they were absent from the current workspace. This does not prove absence from previous attempts or other snapshots. Existing build source already implements dirty checkpointing, interrupted-work salvage, branch cleanup and verification-dependent success. No source rewrite or speculative product repair was performed.

The new 159-line test module defines four disposable-local-Git scenarios against `_implement_issue_impl`: positive tracked repair with a real Python verifier subprocess; rejected coding outcome with no verifier effect; completed outcome with a failed verifier; interrupted tracked repair checkpoint followed by a fresh invocation based on the salvaged branch. Assertions cover caller HEAD/branch/content isolation, no remote, worktree cleanup, one commit per invocation, preserved interrupted branch/blob, and one verifier call per eligible invocation. GitHub PR creation is disabled, and unexpected external-agent calls raise an assertion.

LIMITATION: The coding loop is replaced with a deterministic boundary stub. These are orchestration/Git integration scenarios, not execution of the native coder/reviewer provider or the public decorated entrypoint. The restart scenario is a second invocation using the saved branch, not durable native task-to-result replay. Foreign CrewSync ownership is not inspected or mutated; disposable Git isolation assertions are not proof of distributed ownership fencing. No DTO version or C4 interface declaration is available in the permitted source evidence. Imported schemas/git helpers cannot be independently reviewed outside the authorized read scope. Therefore even a future green run of this module alone cannot satisfy every C4 acceptance gate.

## Local measurable DoD and check status

| Requirement | Status | Evidence |
|---|---|---|
| Preserve existing implementation and exact source identity | Source inspection only | Build SHA above; no build source edit |
| Positive + two adversarial negatives + interruption/restart | NOT EXECUTED | Four focused scenarios created; no runtime assertions evaluated |
| No duplicate effects / isolation | NOT EXECUTED | Assertions exist but runtime remains blocked |
| Real focused exit and affected regression | BLOCKED | Runner returned null exit before Python execution; no broader authorized test target supplied |
| C4 interface/DTO match | BLOCKED | Current reducer explicitly lacks schema/version and producer/QA receipts |
| Independent QA/reducer | REQUIRED | No self-acceptance or terminal task mutation |

## Exact runner receipts

Initial attempt before creating the authorized missing test:

```text
tool: go_wave.run_project_tests
file: tests/issue/test_northstar_mvp_r12_git_redgreen.py
receipt: GO_WAVE_NODE_ERROR:test_file_missing
```

After creating that file (a changed precondition, not an identical denial retry):

```json
{"status":"BLOCKED","code":"PYTHON_VENV_UNAVAILABLE","test":"tests/issue/test_northstar_mvp_r12_git_redgreen.py","engine":"python:pytest","source_identity_sha256":"aa0c1783c4621258818b5b6c0350100c07128b22e36f434beef68561c7cbc6bf","exit_code":null,"signal":null,"error":null,"duration_limit_ms":120000,"output_truncated":false,"stdout":"","stderr":"A regular local .venv Python executable is required."}
```

No retry after the environment denial. No shell, environment installation, commit, push, deployment, send, lease/claim mutation or acceptance action was performed. Git status/branch history and native C4 DTO evidence were not independently available through the scoped interfaces; trusted prelaunch HEAD is not a clean-worktree attestation.

## Self-review and recovery checkpoint

Changed files are only the new focused test and this evidence document. Existing build source is unchanged. Tests use temporary repositories, local test identity, no origin remote, tracked source changes and `sys.executable` for the verifier. They intentionally retain interrupted deliverables and assert cleanup of their worktrees. Runtime import compatibility, fixture correctness and cleanup are unverified until pytest executes; do not treat source review as PASS.

Initial useful-work estimate: investigation 4 minutes, bounded test/artifact implementation 5, focused validation/repair 5, review/evidence 3 (17 total; external environment restoration unbounded). Execution stopped at a concrete environment prerequisite, not to consume an invented duration. Tool receipts do not provide a reliable total elapsed-work measurement.

Recovery: restore a regular project `.venv` through the authorized infrastructure owner, then resume this same task/attempt/thread and execute `python -m pytest -q tests/issue/test_northstar_mvp_r12_git_redgreen.py` through the scoped runner against the recorded final SHA. Diagnose only observed fixture/product failures and rerun after each in-scope repair. Supply exact C4 DTO/version and original producer/QA receipts through coordinator-authorized read scope; native task/provider/replay verification remains a separate required gate. Independent QA and reducer must evaluate real exits before any acceptance claim.

The immutable worker publication response is recorded in the conversation/tool receipt after this document is written; no transport status is predicted here.

## 2026-10-15 exact-task resumption — infrastructure checkpoint

**FACT: BLOCKED / UNACCEPTED.** This resumption inspected the existing candidate rather than recreating it. The historical attempt evidence above is preserved and is not the current attempt binding.

### Current trusted binding and source receipts

- Task: `swe-af:v2-r13-issue-git-redgreen-20261014-r1`.
- Attempt: `attempt-ad6524b1-fa15-4618-81ad-f139c52727fb`.
- Session: `manual-20261015-basefive-goal20-swe-af-v2-r13-i-7c24ec4eeb9bfe37`.
- Producer thread: null; role: implementer.
- Target root: `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af` (launcher-declared; no shell/Git inspection authorized).
- Trusted prelaunch Git: `064a7226a500365005b5b662044fc9014a328722`; manifest SHA-256: `8a2abf4b31d4111628fb8aa718da720fd6c95ce701e3a3920a5948a97739268b`.
- Build source scoped-read SHA-256: `440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd`.
- Existing focused test scoped-read and runner SHA-256: `aa0c1783c4621258818b5b6c0350100c07128b22e36f434beef68561c7cbc6bf` (159 lines; unchanged this attempt).
- Governing C4 artifact: `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`; trusted prelaunch SHA-256: `feeb1fe6a88e53770831de890fe3674e6425dd40327deb913fbf73f3a01f1f17`. Its contents were read; its own hash was not independently computed.

### Bounded DoD and observed result

Value: obtain real execution evidence for the existing candidate without duplicating implementation. The smallest compatible move was one authorized focused run on the existing test SHA, followed only by repairs justified by actual failures. Required proof remains: positive repair, two adversarial negatives, interruption/restart, isolation/no duplicate external effects, affected regression, exact C4 compatibility evidence, and independent QA/reduction. Passing deterministic boundary tests would not alone prove native provider/replay or distributed ownership behavior.

BMad Help was loaded first and resolved the target from the trusted launcher; the existing C4 reducer artifact is the current gate, not canonical BMAD-MNNZ. Verification-before-completion was then loaded. Shell-based configuration resolution, Git status/history, graph navigation and native environment installation were unavailable under the scoped interfaces. No unsupported workflow step is claimed as executed.

Exactly one current-attempt runner call was made:

```json
{"status":"BLOCKED","code":"PYTHON_VENV_UNAVAILABLE","test":"tests/issue/test_northstar_mvp_r12_git_redgreen.py","engine":"python:pytest","source_identity_sha256":"aa0c1783c4621258818b5b6c0350100c07128b22e36f434beef68561c7cbc6bf","exit_code":null,"signal":null,"error":null,"duration_limit_ms":120000,"output_truncated":false,"stdout":"","stderr":"A regular local .venv Python executable is required."}
```

This is an infrastructure failure before pytest execution, not a failed product assertion and not PASS. No identical retry, substitute Node test, environment workaround or speculative product/fixture repair was attempted. Positive, negative, restart and affected-regression runtime checks remain unexecuted. Exact C4 DTO/schema versions and original producer/QA receipts remain unavailable in authorized evidence; no compatibility or foreign-owner-fencing claim is supported.

### Self-review, sizing and recovery

Only this evidence document changes in the current attempt. Build source and all four existing test scenarios remain untouched; no task/lease/claim, foreign state, remote, commit, deploy or send operation was performed. Source inspection confirms the intended scenarios exist, not that their assertions pass. No broader test file is authorized, and no native runtime success or clean-worktree attestation is asserted.

Estimated useful scope before the runner gate: investigation 3 minutes, implementation 0 unless a defect is reproduced, focused execution/repair 4–8, compatibility/evidence review 3 (10–14 minutes total; infrastructure restoration overhead unknown). The concrete prerequisite failure prevents that runtime cycle; tools do not expose a reliable measured total duration. No work was added merely to meet a duration target.

Next milestone: the infrastructure owner must supply a regular project `.venv` compatible with the scoped runner, then resume this exact sticky task/owner and run `python -m pytest -q tests/issue/test_northstar_mvp_r12_git_redgreen.py` through `go_wave.run_project_tests` against the recorded SHA. Repair only demonstrated in-scope defects, rerun after repairs, and obtain coordinator-authorized C4 DTO/producer/QA evidence and affected-regression scope. Independent QA and reducer remain mandatory before any acceptance. The single current-attempt publication response will remain in the tool/conversation receipt; no transport status is predicted here.
