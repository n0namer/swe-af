# R24 orchestrator resume — 2026-10-15

## R25b successor checkpoint — 2026-10-15

**BLOCKED_INFRA / UNACCEPTED.** This section supersedes the current disposition and recovery instructions below for the successor only; historical R24 evidence remains intact.

### Exact successor binding

- Task: `swe-af:v2-r25b-orchestrator-resume-20261015-r1`.
- Attempt: `attempt-d28d0999-28a8-44ed-ba12-e2b76656f8c7`.
- Session: `manual-20261015-basefive-r25b-win32retry-swe-af-1d170e9ca03c0d01`.
- Target root: `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af`.
- Trusted prelaunch Git: `064a7226a500365005b5b662044fc9014a328722`.
- Manifest SHA-256: `e7dc851f626cb0e54b7808759cdc8dfd41ae4f462de068f5622e51dabed04ee8`.
- Task-context SHA-256: `e0b6d88edfc290cd918218c92dc0e25d3c496a97ee8001d164159fb6a1baae71`.

FACT: BMad Help was loaded first. The existing C4 reducer frontier and C4 JSON contract were read directly. This is a bounded implementation/acceptance-repair stage with an unmet native runtime gate, not a new planning phase. Quick Dev and its clarify/route step were loaded, followed by verification-before-completion. Configuration scripts, interactive workflow gates, Git commands and a complete Quick Dev workflow were not executed: only scoped MCP interfaces are authorized. No competing plan was created. Graph/Graft tools are not exposed; exact source reads were used instead. Current branch, dirty checkout identity and complete checkout hash remain independently unverified.

### Local DoD and implementation decision

Value: source-backed restart safety, fail-closed owner/source rejection, and monotonic durable results without duplicate effects. The full local DoD requires an initial failing native product test; the smallest compatible production repair in the authorized files; positive, stale-owner and invalid-source cases; real process crash/restart; affected regression; source-bound evidence; and one typed UNACCEPTED publication for this successor. Tests must compile with unique package helpers. Independent QA and coordinator terminal readback remain external acceptance gates.

FACT: The mechanism is partially implemented already. `ResumeBuildHandler` loads the checkpoint, compares some owner/epoch combinations, reconstructs the plan, adds `resume_dispatch`, calls `Deps.CallRaw("execute")`, and writes `last_resume_result` after the call. Reimplementing these features under new names would not establish the missing guarantees.

FACT: Current source inspection exposes these bounded gaps, not executed product failures:

1. `Deps.CallRaw` forwards to `App.Call` using the node-local execute target. `ExecuteHandler` decodes into `executeInput`, which has no `owner_id`, `lease_epoch`, `run_id` or `resume_dispatch` fields. Its DAG options forward resume/build/workspace/git settings but not those fencing or dispatch fields. Its external executor closure forwards only issue and repository. Thus additive top-level metadata from resume is not, by itself, production effect-side deduplication. `build_id` is an execute field, whereas resume forwards a selected build identifier under `run_id`.
2. Resume persists only after the external call. Its last-result fields are not consulted before a subsequent dispatch. Dispatch keys include owner and epoch, so they are not stable across ownership transfer. No exactly-once or terminal-monotonicity claim follows from these fields.
3. A persisted owner/epoch with an omitted request owner bypasses the owner comparisons. `mapInt` also truncates float64 values rather than rejecting fractional epochs. These are source-level fail-closed concerns requiring native regression tests before repair.
4. The existing missing/malformed-checkpoint test passes `&Deps{}`. Resume returns the missing-App error before opening the checkpoint, so its arbitrary-error assertion does not prove checkpoint rejection. Do not promote this historical test to two meaningful negative native cases.

Source-backed component path: caller input -> `ResumeBuildHandler` -> checkpoint JSON and `Deps.CallRaw` -> `App.Call("<NodeID>.execute")`; the corresponding execute component is `ExecuteHandler` -> `runDAG` seam bound to `dag.RunDAG`. The actual control-plane registration/transport, DAG implementation, effect store and downstream lease behavior were not traced in this checkpoint. The C4 JSON describes build dispatch, not a versioned resume/fencing DTO; its wire compatibility and runtime remain explicitly UNVERIFIED. This is partial source-path evidence, not complete C4 compatibility acceptance.

Read receipts:

| Source | SHA-256 |
|---|---|
| `go/internal/orch/common.go` (lines 1–240 read) | `8d2f93cd901ab10d3b6a9a5b527360343c189ec065239d8107fa219a08f18609` |
| `go/internal/orch/execute.go` (complete) | `2c896419e1f9d7ba83d729b8a5be1b86e9fb41e5ac31c7d5ce11cf375d3f9254` |
| `go/go.mod` (complete) | `12884afec8b7feed6ff3c50e28f0a889ca9066fdb9f56d826ff42f8111592127` |
| C4 reducer frontier | `feeb1fe6a88e53770831de890fe3674e6425dd40327deb913fbf73f3a01f1f17` |
| C4 JSON contract | `cc0eeb34cbc3031785787caf1e72ce8393551533376d46403b71a393c3330029` |

Resume and the R24 test were read completely; their full-file read interface returned content without a file SHA. No individual SHA is invented for them. Runner source identity below is not a complete dirty-tree attestation.

### Fresh runner receipt

One authorized call was made in this successor attempt:

```text
go_wave.run_project_tests({"file":"go/internal/orch/r24_orchestrator_resume_test.go"})
```

```json
{"status":"FAIL","code":null,"test":"go/internal/orch/r24_orchestrator_resume_test.go","engine":"docker:go-test","source_identity_sha256":"0ba7225c8a03abe137666c23bfb95cb1453501de6d821ecdf6e2072b1776706c","exit_code":1,"signal":null,"error":null,"duration_limit_ms":120000,"output_truncated":false,"stdout":"","stderr":"go: cannot find main module, but found .git/config in /workspace\n\tto create a module there, run:\n\tgo mod init\n","package":"./go/internal/orch","go_image":"golang:1.24.0-alpine"}
```

FACT: Tool wall time was 2.4145 seconds. The existing module is `github.com/Agent-Field/SWE-AF/go` in `go/go.mod`. The adapter fails before package loading; this is infrastructure failure, not PRODUCT_RED. The receipt does not attest either requested `-run` selector. No identical retry or alternate execution bypass was attempted. Affected regression, compilation, real process restart, effect count, terminal monotonicity and package-helper uniqueness are unverified. No production or test edits were made after this blocker; only this existing evidence artifact was updated.

### Recovery boundary and self-review

1. Restore the scoped native runner's module-root mount/working directory, without initializing another module. Resume this same successor owner after restoration. Required checks remain `cd go && go test ./internal/orch -run TestR24OrchestratorResume -count=1` and `cd go && go test ./internal/orch -run TestR24 -count=1 -timeout 3m` through the authorized adapter.
2. Replace the historical arbitrary-error fixture with an App-backed fixture, reproduce the actual owner/source boundaries, and complete the positive/process-restart slice before asserting repair. This remains authorized in the existing test file; it is deferred because the current native engine cannot compile or execute it.
3. Inspect DAG/effect interfaces within the existing read scope. If no existing effect-side receipt/deduplication seam can consume the resume metadata, coordinator replanning is required for at least `go/internal/orch/execute.go` and the identified DAG/effect boundary. These are not authorized writes here. A resume-only persisted flag cannot safely substitute for an external effect receipt across the effect-before-checkpoint crash window.
4. Independent QA must verify the final source-bound candidate; coordinator alone owns terminalization and readback. Publication transport does not satisfy this gate.

Self-review: this checkpoint changes only the authorized existing acceptance document, preserves prior evidence, adds no dependencies/helpers/runtime behavior, and makes no PASS, fix, full source-identity or acceptance claim. No secrets, task administration, live leases, commit, deploy or send actions were used. Useful-work estimate was 8–12 minutes for bounded inspection, one runner call and evidence; a complete runtime repair is provisionally 20–30 minutes after runner restoration and effect-interface confirmation. Total elapsed useful time was not instrumented. Next meaningful milestone: an actual compiled native reproduction under the existing Go module. The typed successor report and exact publication response are separate transport evidence; they are not preclaimed in this artifact.

## Disposition

**BLOCKED / UNACCEPTED — infrastructure and source-access prerequisites.**

FACT: Two bounded checkpoint-rejection cases were added to the authorized new Go test. The scoped native runner exited before package loading. Neither case compiled or executed. Positive execution, process restart, effect deduplication, hostile foreign-owner rejection and terminal consistency are not proved. There is no product PASS or demonstrated PRODUCT_RED.

## Binding and source of truth

- Target: `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af`, explicitly supplied by the launcher, not BMAD-MNNZ.
- Task: `swe-af:v2-r24-orchestrator-resume-20261015-r2`.
- Attempt: `attempt-2bd28e8d-0220-4bab-b89a-ca94e48bf9a4`.
- Session: `manual-20261015-basefive-r24b-20owners-swe-af-v-5d6c31ff0edc5cf6`.
- Trusted Git receipt: `064a7226a500365005b5b662044fc9014a328722`; HEAD, branch, dirty state and history were not independently checked.
- Trusted manifest SHA-256: `7073d662c4128d905310e3eb6262864f4a11d5fdf475a0bd0e261f3088ae0669`.
- Trusted project-cache SHA-256: `031364d219a544dc58e28dd2b7485663378006fb8bd35314e60cc07d4217f10a`.
- Trusted task-context SHA-256: `99835b05b704254f5300f4d977c86e51d41afd7a7235522a60f627a908929d75`.
- Existing BMad gate: `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`, read in this attempt; remains UNACCEPTED/BLOCKED. Its Python runtime receipt does not validate Go resume.
- C4 contract: `docs/architecture/basefive-r14-contracts.json`, read in this attempt. It binds `go/internal/orch/build.go` to SHA-256 `05c39a2e3b41869b01a1482543c9875630a0777fbea24ea706c05727fe640cbc`, but explicitly labels itself source-contract-only, wire compatibility UNVERIFIED and runtime UNVERIFIED.

FACT: This is the acceptance-evidence stage, not authority to repair production orchestration. BMad Help was loaded first, followed by Quick Dev and verification-before-completion. Quick Dev clarify/route and plan references were read. Configuration scripts, Git inspection, generated spec creation and interactive workflow gates were not executed: shell is forbidden, configuration paths are outside read_scope and a competing planning artifact is forbidden. No claim of full BMad workflow completion is made.

## Evidence-supported path and ownership limits

The C4 JSON's component/container entries report `Handlers -> Build -> Deps` dispatch to `plan`, `execute`, `run_verifier` and `generate_fix_issues`, followed by `BuildResult`/`ReasonerFailed`. These are contract-reported edges, not call paths observed in this attempt's runtime. The contract distinguishes Python issue orchestration from Go full-build orchestration; Python commit/success invariants cannot be silently applied to Go.

The previously existing `docs/acceptance/longbundle-orch-resume-20261014.md` reports `ResumeBuildHandler(context.Context, *Deps, map[string]any) (any, error)`, checkpoint loading at `<repo>/<artifacts>/execution/checkpoint.json`, plan reconstruction and `CallRaw("execute")` with `resume: true`. It records historical `resume.go` SHA-256 `02c1a49916f19a1e489ad8a94d6f76aea2ff39869b3bec9351bd96ca33cb7f48`. This is secondary historical evidence, not a fresh implementation-source receipt.

`docs/acceptance/basefive-r12-orch-lease-restart.md` reports existing R9 mock replay tests and explains that repeated dispatch, an unchanged unrelated checkpoint and mock results do not prove durable restart, external-effect deduplication or owner fencing. Therefore this attempt does not duplicate those tests under a new name or call them native restart proof. No prior VERIFIED R24 candidate receipt was available in the inspected evidence; absence across the repository is not asserted.

Implementation, Deps, checkpoint DTOs, adjacent test sources and nested Go instructions are outside the explicit read_scope (except this new exact test target). Graph/Graft MCP tools are not exposed. A Graphify file listing confirms cached artifacts exist but establishes neither source freshness nor runtime behavior. The C4 contract records zero Go files in ArchSteer; ArchSteer is not treated as Go coverage. Structurizr availability is not established. Current implementation owners, lease schema/version and terminal-transition interfaces remain unverified. Coordinator ownership of tasks, claims, leases and terminal acceptance is preserved.

## Bounded implementation and local DoD

Value: executable rejection checks and truthful recovery evidence, without inventing successful resume/fencing behavior.

Added `go/internal/orch/r24_orchestrator_resume_test.go` with `TestR24OrchestratorResumeRejectsUnreadableCheckpoint` and two subtests:

1. Missing checkpoint: invoke the documented handler with a disposable repository and empty dependencies; require an error and continued absence of `checkpoint.json`.
2. Malformed checkpoint: write invalid JSON in the documented location; require an error and byte-for-byte preservation of the rejected file.

ASSUMPTION: The handler signature, `Deps` construction and `repo_path`/`artifacts_dir` inputs remain compatible with the historical acceptance evidence. Compilation and the exact error path are unverified. An arbitrary error is not proof that checkpoint parsing was reached. Empty dependencies intentionally supply no configured dispatcher, but do not prove no external effect occurred. Panics are not recovered or converted to success. Tests have no skip, network call, production fixture, success stub or permissive fallback.

The full local DoD remains source-backed current path/owners, positive execution, two meaningful rejection conditions, genuine crash/restart, no repeated effects/foreign-owner drift/false terminal success, focused native execution, affected regression, self-review and independent QA/reducer handoff. Only a partial negative test artifact and exact blocked-run evidence were delivered. No repair of production behavior is justified by this infrastructure failure.

## Actual runner receipt

Requested validation from the launcher:

```text
cd go && go test ./internal/orch -run TestR24OrchestratorResume -count=1
```

Actual authorized invocation, once:

```text
go_wave.run_project_tests({"file":"go/internal/orch/r24_orchestrator_resume_test.go"})
```

Complete returned receipt:

```json
{
  "status": "FAIL",
  "code": null,
  "test": "go/internal/orch/r24_orchestrator_resume_test.go",
  "engine": "docker:go-test",
  "source_identity_sha256": "0ba7225c8a03abe137666c23bfb95cb1453501de6d821ecdf6e2072b1776706c",
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

Classification: INFRA BLOCKED, despite raw runner status FAIL. The runner used package `./go/internal/orch`; it does not attest the requested working directory or `-run` selector. The source identity is the runner's receipt, not a separately calculated test-file hash or immutable complete checkout snapshot. Tool wall time was 4.2824 seconds. No unchanged retry, shell/Node bypass, module initialization or runner modification was attempted. The same module-root symptom appears in historical evidence but this is the current attempt's own receipt.

## Acceptance matrix and self-review

| Requirement | Current evidence | Disposition |
|---|---|---|
| Current source/C4 path and owners | Trusted Git and read C4 contract; historical resume description only | PARTIAL; current source/owners unverified |
| Native positive execution, exact command/exit | Actual adapter invocation and exit 1 recorded; no package loaded | BLOCKED |
| Two negative conditions | Missing/malformed cases authored; no assertions executed | BLOCKED |
| Genuine process crash/restart | No current API/source basis for safe realistic fixture | BLOCKED; not substituted with mock replay |
| No repeated effects, foreign owner or false success | No execution/owner contract available | UNVERIFIED |
| Preserve failing evidence without weakening fixture | Exact adapter output retained; no skips or softened assertions | Receipt preserved |
| Affected regression | Same native infrastructure prerequisite | BLOCKED; not attempted through another runner |
| Independent reviewer/reducer | Required after source and runner restoration | PENDING |

Self-review: both new artifacts are within the two exact authorized paths and below 500 lines. Tests use only standard-library imports and the documented existing handler; no dependency, runtime source, module file, coordinator state or foreign file was changed. Their intended scope is strictly checkpoint rejection, not resume acceptance. The suite remains uncompiled and may need same-owner fixture repair after source access and runner restoration. Broader Git/diff, formatting, type checks and nested-instruction validation are not asserted to have run.

## Recovery and handoff

1. Operator/coordinator restores the scoped native adapter's module-root mount/working directory and verifies it can execute from the existing `go` module. Do not follow the diagnostic suggestion to create a competing module.
2. Coordinator grants minimal read access to `go/internal/orch/resume.go`, the actual Deps definition, checkpoint/owner DTOs, nearby relevant tests and applicable nested instructions. Do not broaden writes: production defects require affected-owner repair.
3. Resume this exact task/attempt owner, inspect those sources, complete the positive and real process-restart/effect/fencing fixture, and run the requested focused command through the corrected adapter. Repair only these owned artifacts against real evidence, then run explicitly authorized affected native regression.
4. Independent reviewer/reducer evaluates the executable evidence; the coordinator owns authoritative terminal readback. No replacement implementation task, new claim or self-acceptance is requested.

Useful-work estimate: 8–12 minutes for permitted evidence inspection, bounded negative fixture, one runner call and receipt/self-review. Remaining positive/restart work cannot be responsibly sized without current source and a usable native runner. Total useful elapsed time was not instrumented; no duration target was padded. The next meaningful milestone is compilation and actual execution under the existing Go module, not another unchanged failed launch.

One typed BLOCKED/UNACCEPTED result is to be published for the exact binding above. Its publication response is transport evidence only and is reported separately; no VERIFIED transport status is preclaimed here.
