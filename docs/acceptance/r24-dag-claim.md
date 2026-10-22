# R24 DAG claim checkpoint acceptance — 2026-10-15

## R25b successor repair checkpoint — BLOCKED_INFRA / UNACCEPTED

This section supersedes the earlier attempt's source-access limitation, not its acceptance decision. Historical R24 evidence below is retained verbatim. No production repair, native product RED/GREEN, or independent acceptance is claimed.

### Exact binding and scope

- Task: `swe-af:v2-r25b-dag-claim-restart-20261015-r1`.
- Attempt: `attempt-ee97ec1f-f51d-49f6-8828-de7676f3fe3f`.
- Session: `manual-20261015-basefive-r25b-win32retry-swe-af-04a5b205d4ed8c83`.
- Target: `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af`, not canonical BMAD-MNNZ.
- Trusted prelaunch Git: `064a7226a500365005b5b662044fc9014a328722`.
- Manifest SHA-256: `e7dc851f626cb0e54b7808759cdc8dfd41ae4f462de068f5622e51dabed04ee8`.
- Current writes: this existing report and new `go/internal/dag/r24_claim_checkpoint_test.go`. Authorized `go/internal/dag/checkpoint.go` remains unchanged.
- Live Git branch/HEAD/status and whole dirty-tree SHA are not independently available through the scoped tools. Prelaunch Git is not represented as a clean/current tree attestation.

### Value, smallest move and measurable DoD

Value: demonstrate real native recovery boundaries without substituting a test-local claim oracle. Constraint: production writes are limited to checkpoint.go; executor.go and gates.go are read-only, and schema definitions are outside read scope.

Local DoD: inspect actual production call paths and current file identities; reproduce a product failure with native tests; minimally repair owner/source checks and crash fencing; prove positive replay, stale-owner and invalid-source rejection, controlled process kill/restart without a second effect, and affected regression; self-review and publish exact source-bound evidence for independent QA. Only source inspection, regression-probe implementation, self-review and infrastructure evidence are achieved in this attempt. Full DoD remains unmet.

Useful-work estimate before implementation: investigation 5 minutes, implementation 5, tests 5, repairs 4, compatibility review 3, evidence/publication 3 = 25 useful minutes; external coordination overhead unknown. Actual useful elapsed time was not measured. No artificial work is added to reach that estimate. Next meaningful milestone is native compilation and genuine RED evidence after runner module-root restoration.

### Source identities and C4 boundary

FACT: Source-range reads returned these current whole-file SHA-256 values:

| File | SHA-256 |
|---|---|
| `go/internal/orch/execute.go` | `2c896419e1f9d7ba83d729b8a5be1b86e9fb41e5ac31c7d5ce11cf375d3f9254` |
| `go/internal/dag/executor.go` | `f5d89e592cda54615241563c02e852afedae47f9a105e94531766ecb2f57d844` |
| `go/internal/dag/gates.go` | `591571f78dec85564309411a4e0aed6aebb4baa3464fed59f36fe43ab2c73ba5` |
| `go/internal/dag/checkpoint.go` | `ab3d45da6b49558d9ca65c3308d69cacc55b7b10236b614958df4dd82e2f64d8` |
| `go/internal/dag/helpers.go` | `362bd5e101e62c1594877fe8f9fbf019db22d1afba3cedb40c746f898a00a462` |
| `go/internal/dag/executor_test.go` | `88b33f2151ae4855f948717d6924d3f9e0a7f1fa69736b0084362174d9b0dc30` |
| New test, first runner input | `eab9dd4bc8b29fd0a6554da9224a2ee4e0153ea7034a84fda8f9960bb725bd35` |
| New test, final after source review | `58d3daadb209fe3e83f57a4e116d1634891a6122d967640b4772f0e73933bdc3` |

FACT: The source-grounded component path is `orch.ExecuteHandler` -> `runDAG = dag.RunDAG` -> `initDAGState` / optional `loadCheckpoint` -> `saveCheckpoint` -> `executeLevel` -> `executeSingleIssue` -> `runExecuteFn` -> injected `ExecuteFn`. The production external-target closure forwards only `issue` and `repo_path` through `callFn`. It does not forward a durable claim/fencing epoch or effect idempotency key in that closure. The alternative coding-loop branch and worktree/gate calls are not covered by the new external-callback probes.

FACT: At the context/container boundary, the existing r14 contract places the Go goal/input-map entry inside the Go orchestration boundary, with agent-call and local-artifact dependencies. Its SHA remains `cc0eeb34cbc3031785787caf1e72ce8393551533376d46403b71a393c3330029`; it explicitly marks wire compatibility UNVERIFIED and different entrypoints not parity. The new directly inspected execute/DAG path refines that component evidence; it does not certify the whole C4 model, DTO definitions or Python parity. Module identity is `github.com/Agent-Field/SWE-AF/go`, with `go 1.21` in `go/go.mod`. No graph tool is exposed in this attempt; no Graft/Graphify/ArchSteer freshness or coverage claim is made.

### Duplicate-mechanism check and root-boundary findings

FACT: Existing `northstar_checkpoint_owner_fence_20261013_test.go` already exercises save/load, stale levels, revision mismatch, terminal-status regression, corrupt input and retry. Those checks do not invoke RunDAG's effect/replay boundary. No duplicate helper was added for its checkpoint fixture.

FACT: `saveCheckpoint` currently performs load/check followed by plain `os.WriteFile`; there is no atomic read-compare-write transaction in that inspected function. Its `void` API silently returns on rejection and IO failure. `RunDAG` does not check persistence success before calling `executeLevel`.

FACT: On resume, `RunDAG` replaces the initialized state with any successfully decoded checkpoint without comparing the requested source revision/build identity. A missing or corrupt checkpoint returns nil and execution continues from a newly initialized plan. An in-flight issue is not part of `doneNames`; completion is recorded only after the effect callback returns. These are source observations, not executed failures.

HYPOTHESIS: A foreign-source resume, corrupt-checkpoint resume and a process kill after a durable callback effect but before result return violate the requested fail-closed/no-duplicate behavior. The new probes target those paths, but the runner did not reach them.

BLOCKER: Checkpoint-only edits cannot safely enforce an execution admission check through the unchanged void API, validate requested resume identity before the loaded state overwrites it, or fence an arbitrary external effect callback. Panics, indefinite blocking, fake completion, silent issue removal and test-local idempotency would conceal the missing contract rather than repair it. At minimum, a compatible error-propagating admission/recovery change in read-only `go/internal/dag/executor.go` is required. A complete durable owner/epoch/effect contract may additionally affect `gates.go` and the schema/callback boundary; those edits are not authorized here. BuildID is a build namespace, not evidence of a task lease or stale-owner token. No stale-owner acceptance is claimed.

### New native probes and review

- `TestR24ClaimCheckpointCompletedReplay`: invokes actual RunDAG twice and requires one callback invocation and one durable completed issue.
- `TestR24ClaimCheckpointRejectsInvalidResumeSource`: two independently isolated negatives, foreign revision and corrupt checkpoint; each requires zero callbacks and an explicit error. It verifies fixture persistence before exercising resume.
- `TestR24ClaimCheckpointCrashRestart`: starts the current compiled Go test binary, waits for an fsynced append-only effect and an in-flight checkpoint, kills and waits for the child, resumes via actual RunDAG, then requires exactly one durable effect and non-regressing checkpoint progress. Fail-closed recovery is allowed; silently repeating the callback is not.
- `TestR24ClaimCheckpointCrashChild`: only activates with the private test environment variable; the parent uses an anchored test name. Temporary directories isolate all files; no external service, Git mutation or live task state is used.
- New helper ownership uses the `r24Claim` prefix; existing test helpers are not modified. Compilation/package-wide uniqueness is unverified because Go never compiled the package.

Source review added a bounded one-minute timer while the crash child waits so it does not rely on an all-goroutines-asleep deadlock to terminate. That refinement changed the test SHA after the blocked runner input. The final SHA has not executed. The process probe requires actual child-process support from the restored native runner. No assertion, skip, production contract or fixture was weakened to obtain GREEN. This is an uncompiled regression-probe artifact, not a QA candidate or a proven failing regression test.

### Exact current runner receipt

Invocation: `go_wave.run_project_tests({file:"go/internal/dag/r24_claim_checkpoint_test.go"})`.

```json
{"status":"FAIL","code":null,"test":"go/internal/dag/r24_claim_checkpoint_test.go","engine":"docker:go-test","source_identity_sha256":"eab9dd4bc8b29fd0a6554da9224a2ee4e0153ea7034a84fda8f9960bb725bd35","exit_code":1,"signal":null,"error":null,"duration_limit_ms":120000,"output_truncated":false,"stdout":"","stderr":"go: cannot find main module, but found .git/config in /workspace\n\tto create a module there, run:\n\tgo mod init\n","package":"./go/internal/dag","go_image":"golang:1.24.0-alpine"}
```

Classification: **BLOCKED_INFRA**, not a product assertion failure. The runner selected package `./go/internal/dag` from `/workspace`; the existing module is `/workspace/go/go.mod`. No root module is created as a workaround. The exact argv is not exposed by the receipt; no claim is made that either requested validation command ran. No identical retry or broader run was attempted after this infrastructure failure.

Required after infrastructure restoration, NOT VERIFIED here:

```text
cd go && go test ./internal/dag -run TestR24ClaimCheckpoint -count=1
cd go && go test ./internal/dag -run TestR24 -count=1 -timeout 3m
```

### Method, remaining gates and sticky recovery

BMad Help was first; the existing C4 reducer frontier is the target SoT and remains UNACCEPTED. The current phase is bounded implementation/native acceptance, not release. Quick Dev activation defaults and clarify/plan steps were read. Project BMad config/scripts, Git commands, new spec paths, interactive approval workflow and delegation are not available under this manifest; those method steps are not claimed executed, and no competing PLAN/spec was created. Verification-before-completion was applied to the actual runner output. Systematic product debugging/TDD RED-to-GREEN is not claimed because infrastructure prevented any product reproduction.

All positive, negative, real-crash, stale-owner, package-compile, CAS-concurrency and affected-regression gates remain UNVERIFIED. No production files changed; no acceptance, deployment, lease, claim or remote operation occurred. The main remaining risk is the source-observed unfenced execution/replay path; the uncompiled tests may also need repair once the compiler runs.

Coordinator recovery: preserve this original successor key and sticky owner; restore the Go runner's module cwd/package mapping; run the final test SHA; then authorize only the necessary executor/callback/schema scope after reviewing the source-bound contract gap. Do not rerun or duplicate the earlier R24 implementation. Diagnose actual product RED, repair the authorized boundary, retest both requested commands and obtain independent QA/reducer evidence. Worker transport verification is not product acceptance or terminal readback. The single typed publication and its exact response are recorded by the publication receipt/final handoff, not predicted in this file.

---

## Historical R24 attempt (retained)

## Disposition

**BLOCKED / UNACCEPTED.** This attempt produced a scoped evidence report, not a compiled QA candidate. No native test, positive scenario, negative scenario, process restart or affected regression is reported as passing. No production defect is established.

## Original binding and source receipts

- Project/root: `swe-af`, `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af` (trusted launcher).
- Task: `swe-af:v2-r24-dag-checkpoint-claim-20261015-r2`.
- Attempt: `attempt-4a860159-335e-4b43-96da-b47f2d9baeca`.
- Session: `manual-20261015-basefive-r24b-20owners-swe-af-v-48412608a51bac22`.
- Prelaunch Git receipt: `064a7226a500365005b5b662044fc9014a328722`; live HEAD, branch and dirty state were not independently inspected.
- Manifest SHA-256: `7073d662c4128d905310e3eb6262864f4a11d5fdf475a0bd0e261f3088ae0669` (trusted receipt).
- Task-context SHA-256: `76af60c25a12f7d6db355e4a6c094fb06938b3afe3a80fe39787c8681f5098b2` (trusted receipt).
- Contract: `docs/architecture/basefive-r14-contracts.json`, schema `swe-af.basefive-r14-contracts.v1`; prelaunch SHA-256 `cc0eeb34cbc3031785787caf1e72ce8393551533376d46403b71a393c3330029`.

These identities do not prove the current native DAG implementation or its runtime behavior.

## Value, present state and bounded DoD

Value: independent native evidence that checkpoint recovery and competing resumes do not repeat effects or admit a foreign owner. Constraint: writes are limited to this report and `go/internal/dag/r24_claim_checkpoint_test.go`; the explicit read scope does not include production DAG sources, Go schema definitions, module files or existing Go tests.

The local DoD requires current source-backed call paths and owner/fencing interfaces; a test using those actual interfaces rather than a replacement oracle; native positive execution; two independent negative conditions; real child-process crash/restart; affected native regression; exact command/exit/source receipts; source-backed self-review; and one original-task UNACCEPTED publication for independent review. That DoD is not met.

Smallest safe move completed: inspect the permitted contract and existing acceptance evidence, check both exact output paths, distinguish historical checkpoint coverage from claim fencing, and preserve the authorization gap without fabricating a test.

## Investigation and duplicate-mechanism check

FACT: Both authorized output paths were absent when read. The exact test read returned:

```text
GO_WAVE_NODE_ERROR:ENOENT: no such file or directory, stat 'D:\Users\NIKITA\Documents\DEV\universal-solver-site\swe-af\go\internal\dag\r24_claim_checkpoint_test.go'
```

The report path also returned ENOENT before creation. Missing authorized targets are not themselves blockers; the report was subsequently created through the scoped writer. No placeholder, skipped, always-red or replacement-oracle test was created merely to make a runner invocation possible.

FACT: `docs/acceptance/northstar-dag-fencing-c4-qa-20261013.md` already documents same-package native checkpoint tests for positive persistence, stale level, terminal regression, foreign revision, missing/corrupt checkpoints and retry. It attributes these boundaries to `go/internal/dag/checkpoint.go`:

- `saveCheckpoint(dagState *schemas.DAGState, note noteFunc)`.
- `loadCheckpoint(artifactsDir string) *schemas.DAGState`.
- `checkpointRegressesTerminalIssue` and `isTerminalIssueStatus`.
- DTO fields `ArtifactsDir`, `CurrentLevel`, `GitInitialCommit`, and `AllIssues[].id/status`.

FACT: `docs/acceptance/longbundle-dag-monotonic-20261014.md` additionally documents `recoverableIssueIDs(*schemas.DAGState) []string`. These are historical report attributions, not direct inspection of the current native source. A Git revision guard is not evidence of a task-owner claim token or durable effect fencing. No exact production claim acquisition/renewal/release/CAS API or authoritative owner call chain was established in the permitted evidence.

FACT: `docs/acceptance/basefive-r16-dag-resume-fence.md` explicitly distinguishes its test-local oracle from the production checkpoint engine and reconstructed-object restart from a killed-process test. Duplicating that oracle would not deliver this task's requested native acceptance evidence.

FACT: The r14 C4 contract describes the Go Build boundary and states `wireCompatibility: UNVERIFIED`, `DIFFERENT_ENTRYPOINTS_NOT_PARITY`, `archsteerGoFileCount: 0`, and unverified graph-source freshness. It does not supply a verified DAG owner/fencing contract. ArchSteer is not used as Go source coverage proof. No Graphify/Graft live graph inspection or source completeness claim is made.

## Actual blocker and validation status

**Authorization/source blocker:** current production DAG code, schema definitions, module identity and existing native tests cannot be inspected under the explicit read scope. The write scope authorizes creating the new test, but does not authorize reading its native dependencies. Inferring claim APIs and crash durability from report summaries would risk compiling against invented interfaces and asserting the wrong safety contract. No out-of-scope read or write was attempted.

Requested command (NOT EXECUTED):

```text
cd go && go test ./internal/dag -run TestR24ClaimCheckpoint -count=1
```

Exit code: **not available**. No test source was created, so calling the runner on the absent path would only reproduce the already-observed ENOENT rather than exercise the native engine. This attempt has no fresh runner denial, compile failure, assertion failure or product RED receipt.

Historical infrastructure evidence, NOT a current run: the permitted earlier DAG reports record `docker:go-test`, package `./go/internal/dag`, exit 1, and `go: cannot find main module, but found .git/config in /workspace`, using `golang:1.24.0-alpine`. The longbundle report calls for Go 1.26.5. Neither current runner recovery nor the current required toolchain was independently verified. Do not promote those historical receipts into a current failure or PASS, and do not initialize a root module as a workaround.

| Acceptance gate | Current result |
|---|---|
| Current source SHA, actual call paths and owner contract | BLOCKED; baseline/contract receipts available, native source not authorized |
| Native positive execution with exact command and exit | NOT RUN |
| Two negative claim/fencing conditions | NOT IMPLEMENTED / NOT RUN |
| Real process crash and competing resumes | NOT IMPLEMENTED / NOT RUN |
| No repeated effect or foreign-owner success | UNVERIFIED |
| Affected native regression | NOT RUN |
| Preserve RED evidence without weakening fixtures | No product RED observed; historical infrastructure receipts kept distinct |
| Scoped report and original-task typed publication | Report created; publication/readback must be recorded by publisher response |
| Independent reviewer/reducer | REQUIRED; no self-acceptance |

## Method and self-review

BMad Help was loaded first. The target is swe-af, not canonical BMAD-MNNZ. The current phase is blocked native acceptance/QA, governed by the existing C4 reducer frontier, not release. BMad Quick Dev activation defaults and clarify/plan steps were read; planning stops at its verified-source/readiness gap. Configuration scripts, project customization/config discovery and Git inspection are unavailable within this attempt's tool/read authority; none is claimed executed. No competing spec or PLAN.md was created. Verification-before-completion was loaded and applied: no fresh runtime evidence means no completion/PASS claim. Systematic debugging was not invoked because no current product failure was reproduced.

Only this authorized report changed. Native production sources, existing tests, fixtures, module files, leases, claims and task terminal state were not changed. No source compatibility, clean-diff or runtime regression claim is made. No secrets or external sends were introduced by the report.

Useful-work estimate for the complete slice, once source access is authorized: investigation 4 minutes, implementation 6, focused positive/negative/crash tests 6, repair allowance 4, compatibility/self-review 3, evidence/publication 2; total 25 minutes plus unknown coordination overhead. Actual total elapsed useful work was not measured. No artificial runtime or extra code was added to fill that estimate.

## Same-task recovery and independent QA

1. Coordinator should amend this original task's read authorization to include current `go/internal/dag` source/tests, its relevant ancestor instructions, actual `schemas.DAGState` definitions and Go module/workspace files. Keep write ownership restricted to the original two output paths. No new claim or replacement task is requested.
2. Resume this same Goal Owner/attempt if supported, inspect exact native dependencies and existing coverage, then add only missing native claim-fencing/crash-process assertions. Identify the actual durable owner/effect contract; do not substitute Git revision checks for claim ownership or a test-local state machine for production behavior.
3. Use the authorized native runner from the verified Go module root for the exact focused command and then authorized affected package regression. Capture current test/source identities, actual spawned command, working directory, engine, exit, stdout/stderr and scenario results. Stop identical retries on infrastructure denial. If production repair is needed, report the precise foreign-owned source boundary without editing it.
4. Request independent reviewer/reducer QA on the original task after runtime evidence exists. Worker publication transport, including VERIFIED, does not establish product acceptance or coordinator terminal readback.
