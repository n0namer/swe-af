# R12 orchestration lease/restart — R13 sticky resume evidence

## Disposition

**BLOCKED / UNACCEPTED — infrastructure, not a proven product failure.**

FACT: This attempt inspected the existing R9 fixture and invoked the authorized native runner once. The Go process exited before package loading because the container working directory had no main module. No positive, negative, restart or regression assertion executed. No source repair is justified by this receipt, and no acceptance is claimed.

## Exact task and source binding

- Target root: `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af` (trusted launcher root; not BMAD-MNNZ).
- Task: `swe-af:v2-r13-orch-lease-restart-20261014-r1`.
- Attempt: `attempt-1ac4184e-1ec5-40ea-837a-54f9c0da9a2b`.
- Session: `manual-20261014-basefive-r13-sol-low-swe-af-v2--c8f99b78520bd60d`.
- Git identity: `064a7226a500365005b5b662044fc9014a328722` (prelaunch receipt, not independently checked HEAD).
- Manifest SHA-256: `298b30493572dd2bd10df97368262b3c45edf3cddfa3641f01f45181bd3ac587` (prelaunch receipt).
- Project cache SHA-256: `fdc416aa29e6e605525a13d121f0dddf6303e42c6efc010b18927084392e7ae5` (prelaunch receipt).
- Task-context SHA-256: `83d677aebb287f43b99ed5914c7a5a9cc022bf1195f6c5b018e1ddf0e50cef27` (prelaunch receipt).
- Runner source identity SHA-256: `bb4feaf926564c29685fd9b4a667d75d942619528e49280bbf53369f09f5a20a` (actual runner receipt; not asserted to be an individual file hash or a complete runtime-source snapshot).
- Exclusive write scope: `go/internal/orch/northstar_mvp_r9_resume_adversarial_qa_test.go`, `go/internal/orch/northstar_mvp_r12_resume_qa_test.go`, this document.

## BMad stage and bounded Definition of Done

FACT: The active source of truth is `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`. It keeps the C4 gate UNACCEPTED/BLOCKED and explicitly lacks interface/DTO compatibility and independent producer/QA receipts. Its native runtime receipt concerns Python ReleaseGate, not successful execution of the Go orchestration tests.

Current stage: resume/verification at the blocked acceptance gate. BMad Help was loaded first; verification-before-completion was loaded to enforce fresh evidence. Configuration resolution, Git inspection, native-runner repair and broader source inspection are not authorized by this task's tools/read scope. No shell, graph tool, coordinator mutation or competing plan was used.

Local DoD remains: source-bound C4/DTO evidence; a real positive, two adversarial negatives and restart; no duplicate external effects or foreign ownership drift; focused exit plus affected regression and only proven bounded repairs; independent QA/reducer handoff. This attempt does not meet that DoD.

## Existing candidate inspection — no reimplementation

FACT: `go/internal/orch/northstar_mvp_r9_resume_adversarial_qa_test.go` exists and was read. It contains five tests covering typed result forwarding, cancellation before dispatch, invalid input before dispatch, pinned owner repository/checkpoint metadata, and repeated resume reconstruction without checkpoint mutation.

FACT: Its application fixture records calls and returns a configured map. The restart test expects two dispatches for two resume requests. It does not itself prove deduplication of external effects, durable process restart, lease fencing or rejection of a foreign owner. The repository-isolation assertion checks an unrelated foreign checkpoint is unchanged; it is not evidence of a hostile foreign-owner request being rejected. These distinctions must remain explicit during independent QA.

FACT: Reads of `go/internal/orch/northstar_mvp_r12_resume_qa_test.go` and this document before creation returned `GO_WAVE_NODE_ERROR:ENOENT`. Missing authorized targets are not permission blockers. However, neither an R12 candidate nor a proven product/fixture defect was available in the permitted evidence, so this attempt did not create substitute tests or rename the R9 suite merely to satisfy the R12 selector.

C4 interface names, DTO versions, implementation source and prior exact R12 worker receipts are absent from the allowed evidence. They are not inferred from the mock maps or from the trusted prelaunch Git identity.

## Exact native runner receipt

Invocation: `go_wave.run_project_tests({"file":"go/internal/orch/northstar_mvp_r9_resume_adversarial_qa_test.go"})`.

```json
{
  "status": "FAIL",
  "code": null,
  "test": "go/internal/orch/northstar_mvp_r9_resume_adversarial_qa_test.go",
  "engine": "docker:go-test",
  "source_identity_sha256": "bb4feaf926564c29685fd9b4a667d75d942619528e49280bbf53369f09f5a20a",
  "exit_code": 1,
  "signal": null,
  "error": null,
  "duration_limit_ms": 120000,
  "output_truncated": false,
  "stdout": "",
  "stderr": "go: cannot find main module, but found .git/config in /workspace\n\tto create a module there, run:\n\tgo mod init\n",
  "package": "./go/internal/orch",
  "go_image": "golang:1.24.0-alpine"
}
```

Classification: **INFRA BLOCKED** despite the runner's raw FAIL status. Package loading did not occur. Do not follow the diagnostic suggestion to initialize a new module. Do not repeat this unchanged call or bypass the scoped runner with a shell/Node subprocess.

The requested `go test ./internal/orch -run TestMVPR12 -count=1` was not executed. Affected regression is also blocked by the same module-working-directory prerequisite. No 429 occurred during these tool calls; the previous provider failure is trusted task context only.

## Self-review and acceptance checklist

- Available trusted identity and actual runner SHA/exit preserved: yes; complete source/C4 DTO proof remains blocked.
- Positive and two negative scenarios: present in the R9 source, unexecuted.
- Restart: mock replay assertions present, unexecuted; real interrupted/durable restart unproved.
- Duplicate external effects and foreign-owner rejection: unproved.
- Product/fixture repair: none; no product failure reproduced.
- Changed files: this evidence document only; existing R9 source retained, absent R12 target not invented.
- Independent QA, reducer and terminal readback: still required; no worker transport receipt constitutes acceptance.

## Recovery on the same authorized attempt/thread

1. Coordinator/operator restores the scoped native Go runner's verified module-root working directory and cache configuration outside this worker's write authority. The `/workspace` failure must be resolved without adding a competing module or weakening tests.
2. Restore/locate the original R12 candidate and exact producer/C4 DTO evidence, and authorize the minimal implementation/contract read paths needed to verify them. Do not dispatch a duplicate implementation or infer compatibility from R9 mock data.
3. Resume this sticky owner, run the requested focused R12 command through the repaired native adapter, and run the existing R9 affected regression. Verify the test-selection receipt actually executes the expected tests rather than reporting no tests to run.
4. Repair only reproduced in-scope defects, retest, then hand source-bound receipts to independent QA and the reducer.

Useful-work sizing: the initial bounded inspection/runner/evidence slice was estimated at 5–10 minutes; the implementation/repair estimate cannot be grounded until the original candidate and native module root are available. No idle duration or speculative code was added to meet a time target.

## 2026-10-15 sticky attempt — current receipt

**Disposition: INFRA BLOCKED / UNACCEPTED.** The sections above remain historical evidence for the previous attempt; the binding below identifies this attempt only.

### Binding and inspected evidence

- Task: `swe-af:v2-r13-orch-lease-restart-20261014-r1`.
- Attempt: `attempt-a1a77a61-43ce-4c21-8b21-6b8c5bbf5ee6`.
- Session: `manual-20261015-basefive-goal20-swe-af-v2-r13-o-2ee25aafb075e0f4`.
- Target: `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af`, as declared by the launcher; current branch and HEAD were not independently inspected because shell access is forbidden.
- Trusted prelaunch Git: `064a7226a500365005b5b662044fc9014a328722`.
- Trusted manifest SHA-256: `8a2abf4b31d4111628fb8aa718da720fd6c95ce701e3a3920a5948a97739268b`.
- Trusted project-cache SHA-256: `d9d1ecc2603d6c6b235bc81cfa5a4d6084cd9b1b5d9ee39c121b069b76fa2269`.
- Trusted task-context SHA-256: `3fea1b2577ba58fb2d83e83fc898a087adf527ee5ce976cad08effd49a14cf00`.

FACT: BMad Help was loaded first, followed by verification-before-completion. The target remains swe-af, not the canonical BMad repository. The permitted C4 source of truth was read and still records UNACCEPTED/BLOCKED, missing C4 interface/DTO compatibility evidence and independent receipts. This is an acceptance/verification stage, not authority to introduce a new implementation. Configuration resolution, Git status/history, graph discovery and implementation-contract inspection remain unavailable under the exact tool/read scope; no claim is made that those workflow steps ran.

FACT: The existing R9 test source and this evidence document were inspected. The R12 target read returned exactly:

```text
GO_WAVE_NODE_ERROR:ENOENT: no such file or directory, stat 'D:\Users\NIKITA\Documents\DEV\universal-solver-site\swe-af\go\internal\orch\northstar_mvp_r12_resume_qa_test.go'
```

The absent file is an authorized creation target, not a permission blocker. However, the requested original R12 candidate is not present there, implementation/DTO sources are outside read scope, and no product or fixture failure has been reproduced. Creating a substitute suite or merely aliasing R9 as R12 would not prove the requested native restart, fencing or external-effect properties. Existing tests are therefore preserved unchanged.

### One actual runner invocation

Invocation: `go_wave.run_project_tests({"file":"go/internal/orch/northstar_mvp_r9_resume_adversarial_qa_test.go"})`.

```json
{
  "status": "FAIL",
  "code": null,
  "test": "go/internal/orch/northstar_mvp_r9_resume_adversarial_qa_test.go",
  "engine": "docker:go-test",
  "source_identity_sha256": "bb4feaf926564c29685fd9b4a667d75d942619528e49280bbf53369f09f5a20a",
  "exit_code": 1,
  "signal": null,
  "error": null,
  "duration_limit_ms": 120000,
  "output_truncated": false,
  "stdout": "",
  "stderr": "go: cannot find main module, but found .git/config in /workspace\n\tto create a module there, run:\n\tgo mod init\n",
  "package": "./go/internal/orch",
  "go_image": "golang:1.24.0-alpine"
}
```

FACT: The current receipt matches the historical module-root failure and source identity. No package or test assertion ran. Classification is INFRA BLOCKED, not product FAIL or PASS. The adapter reports package `./go/internal/orch`; it does not attest execution of the requested module-root command `go test ./internal/orch -run TestMVPR12 -count=1`. No retry, alternate runner, shell bypass or module initialization was attempted after this failure. The SHA is the runner's source identity, not an independently computed individual-file hash, clean Git snapshot, or C4 compatibility proof.

### Current DoD and self-review

- Value and smallest move: resume the existing candidate and obtain truthful source-bound execution evidence, without duplicate implementation or foreign state mutation.
- Positive, cancellation/invalid-input negatives and replay: source assertions exist in R9, but all remain unexecuted in this attempt.
- Durable interruption/restart, lease fencing, no duplicate external effects and hostile foreign-owner rejection: unproved; the mock recorder and unrelated foreign checkpoint do not establish these properties.
- Requested R12 focused run and affected regression: BLOCKED before package loading; no green or no-tests-to-run result is claimed.
- C4 interface/DTO versions and exact original R12 producer receipt: unavailable in permitted evidence.
- Product/fixture repair: none; no proven defect justifies one.
- Change scope: only this existing acceptance document is appended; both Go targets are unchanged, and no runtime tasks, leases, claims, terminal state or external effects were mutated.
- Review: current receipt, exact binding, historical/current distinction, scope and unsupported claims checked before publication. No independent QA or acceptance is implied.

### Required recovery milestone

Coordinator/operator must restore the scoped native adapter's module-root working directory and cache prerequisites outside this write scope, locate the original R12 candidate/producer receipt, and grant the minimal implementation/C4 DTO read paths required for verification. Do not initialize a competing module or dispatch a duplicate implementation. Resume this same owner after concrete infrastructure restoration, execute `go test ./internal/orch -run TestMVPR12 -count=1` from the verified module root through the authorized adapter, then the R9 affected regression; preserve test selection, source identity and actual exit. Only then repair a reproduced in-scope defect and hand evidence to independent QA and the reducer.

Useful-work estimate: 5–10 minutes for this bounded inspection, one runner check and evidence preservation; further implementation/testing cannot be responsibly sized until the native environment and original candidate are available. The runner returned in approximately 3.6 seconds. Total elapsed work was not instrumented. No filler changes were made to satisfy a duration target. One typed BLOCKED/UNACCEPTED report is required for this exact binding; publication transport is not product acceptance.

