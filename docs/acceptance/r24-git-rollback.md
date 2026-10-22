# R24 Git rollback/restart — blocked partial evidence

## Disposition — 2026-10-15

**FACT: INFRASTRUCTURE BLOCKED / UNACCEPTED. The requested native acceptance evidence is not delivered.** The scoped Python runner did not start Python. No product failure, rollback success, restart success, exactly-once property, foreign-owner fencing, or PASS is claimed.

The new test is explicitly a **declared C4 contract preflight only**, not an integration test. Its filename is the exact authorized target, not evidence of behavioral coverage. Restoring Python and passing this preflight would still leave the native acceptance criteria unmet.

## Authority and source of truth

- Project/root: `swe-af`, `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af` (trusted launcher declaration).
- Task: `swe-af:v2-r24-git-worktree-rollback-20261015-r2`.
- Attempt: `attempt-84d3dde9-ccbd-4cdb-b3bb-82c6d4cab5d8`.
- Session: `manual-20261015-basefive-r24b-20owners-swe-af-v-c520266100b528cf`.
- Producer: implementer; thread ID null as required by the physical footer.
- Trusted prelaunch Git identity: `064a7226a500365005b5b662044fc9014a328722`. Live HEAD, branch, dirty state and history were not independently queried; shell access is forbidden.
- Manifest receipt: `7073d662c4128d905310e3eb6262864f4a11d5fdf475a0bd0e261f3088ae0669`.
- Existing BMad source of truth/current gate: `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`, still UNACCEPTED / BLOCKED. This report does not replace or promote it.
- C4 contract directly read: `docs/architecture/basefive-r14-contracts.json`; scoped range receipt SHA-256 `cc0eeb34cbc3031785787caf1e72ce8393551533376d46403b71a393c3330029`.
- The C4 contract declares `swe_af/issue/build.py` SHA-256 `440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd`. This attempt did not read or hash production source; the value is a contract declaration, not a fresh source attestation.

## Inspection and bounded design

Value: independent runtime evidence that failed verification cannot become success and restart cannot duplicate commits or affect another owner. Local DoD: inspect actual source/contract and existing mechanism; exercise native positive execution, two negative conditions and a separate-process crash/restart; preserve exact commands/exits and identities; run affected regression; repair only owned artifacts; self-review; publish once under the original binding with independent QA required.

**FACT:** Both exact writable targets initially returned `GO_WAVE_NODE_ERROR:ENOENT`. They are new authorized targets, not an infrastructure blocker. No pre-existing contents were overwritten.

**FACT, document-backed only:** `docs/acceptance/releasegate-git-20261013.md` already describes real disposable Git worktree/branch/commit checks around `implement_issue`, rejected verification and interrupted repair. `docs/acceptance/basefive-r16-git-issue-cycle.md` describes successful coding, coding rejection, verifier rejection/unavailability and cooperative cancellation plus a second invocation based on a salvaged branch. These documents explicitly do not prove operating-system process-kill recovery or durable exactly-once semantics. Existing fixture execution is blocked in their recorded receipts. This does not establish that rollback or process restart is absent from all product source.

**FACT, scope limitation:** `read_scope` includes AGENTS.md, the C4 contract, acceptance documents, Graphify output and `.archsteer`; it excludes `swe_af/issue/build.py`, Git helpers, schemas and existing native test fixtures. Exact production signatures beyond the public C4 declaration and fixture construction cannot be safely established from this permitted evidence. No speculative production stub, alternate rollback subsystem, copied unrelated fixture or duplicate native implementation was introduced. The new preflight is the bounded source-of-truth check that can be authored without inventing those APIs. It does not satisfy the requested full functional slice.

## Declared call paths and owners

The directly read C4 `diagrams` component and container entries cite the following declarations. They are contract-backed paths, **not freshly observed runtime traces**:

| Call/path | Declared owner and anchor | Proof boundary |
|---|---|---|
| Public issue entry to implementation | `swe_af/issue/build.py`: `implement_issue` -> `_implement_issue_impl`; `return await _implement_issue_impl(` | Contract only |
| Implementation to coding boundary | `swe_af/issue/build.py`: `loop_result = await run_coding_loop(`; imported from `swe_af.execution.coding_loop` | Imported owner's source not inspected |
| Checkpoint/salvage | `swe_af/issue/build.py`: `_checkpoint_if_dirty`, `await _checkpoint_interrupted_issue(` | No OS-crash durability claim |
| Verification | `swe_af/issue/build.py`: `_run_verification`, `verification = await _run_verification(` | Conditional on coding/commits |
| Result | `swe_af/issue/build.py`: `return IssueBuildResult(` | DTO definitions not inspected |
| Git boundary | `swe_af/issue/build.py`: `git_ops.add_worktree`, `git_ops.remove_worktree` | Actual helper implementation not inspected |

C4 invariants at `docs/architecture/basefive-r14-contracts.json:68` declare success as `coding_ok and bool(commits) and verify_ok`, conditional verification, and cancellation checkpoint/re-raise. They do not declare automatic deletion/reset of rejected deliverables. Prior reports describe retaining rejected branches. Therefore the requested word "rollback" must not silently become destructive deletion of those branches. Independent QA must distinguish caller/worktree isolation, retained failed deliverables and any actual product rollback contract.

Graph limitations: no graph navigation/query tools were exposed. The Graphify file listing was truncated and was not treated as exhaustive; manifest metadata was read only. `.archsteer/architecture.md` is generated architecture documentation, not live proof. C4 records ArchSteer Go file count zero and unverified graph/source freshness. No Go coverage, Structurizr runtime or complete call graph is claimed. Source-only C4 declares different Python/Go entrypoints and UNVERIFIED wire compatibility; no cross-language parity is inferred.

## Authorized artifact and self-review

`tests/issue/test_r24_git_rollback_restart.py` defines one contract-only preflight. It checks declared project/Git/build receipts, public signature, result-field names, component owner/edge anchors and retained UNACCEPTED/UNVERIFIED trust boundaries. It reads only the authorized C4 JSON. It does not import/run the product, create repositories, dispatch providers, mutate Git, contact remotes, or exercise rollback, two runtime negatives, process restart or owner fencing.

A source-only review corrected the initial diagram-key assumption from `c4` to the directly observed `diagrams` **before** the single runner call. The runner's test identity matches that final correction. No post-run test edit, fixture weakening or product repair occurred.

Final test SHA-256: `1d68fafd4baadce19a82c07b2a0d026556b67908c53cb2e37c05da51b1329b01` (scoped edit and runner receipts). Test collection, syntax/import behavior and assertions remain unverified because Python never started.

## Exact current runner evidence

Requested validation (not claimed executed):

```text
python -m pytest -q -p no:cacheprovider tests/issue/test_r24_git_rollback_restart.py
```

Actual authorized invocation, exactly once:

```text
go_wave.run_project_tests({"file":"tests/issue/test_r24_git_rollback_restart.py"})
```

Exact response:

```json
{"status":"BLOCKED","code":"PYTHON_VENV_UNAVAILABLE","test":"tests/issue/test_r24_git_rollback_restart.py","engine":"python:pytest","source_identity_sha256":"1d68fafd4baadce19a82c07b2a0d026556b67908c53cb2e37c05da51b1329b01","exit_code":null,"signal":null,"error":null,"duration_limit_ms":120000,"output_truncated":false,"stdout":"","stderr":"A regular local .venv Python executable is required."}
```

The tool exposes no argument for the requested `-p no:cacheprovider`; it supplies no executed command here because Python did not start. No exact-command equivalence or successful process exit can be claimed. No identical retry, Node substitute, shell fallback, environment installation or unauthorized broader test occurred.

## Acceptance and remaining gates

| Required gate | Current evidence |
|---|---|
| Actual source SHA, observed paths and owners | PARTIAL: exact declared C4 receipts/anchors; no direct source/runtime observation |
| Native positive command and exit | BLOCKED: exit null; native positive scenario not implemented in this partial artifact |
| Two runtime negatives | NOT IMPLEMENTED/NOT RUN: rejecting verifier and unavailable verifier/coding rejection remain required |
| Process crash/restart without duplicate effects | NOT IMPLEMENTED/NOT RUN; cooperative cancellation in earlier documents is not a substitute |
| No foreign owner or false success | UNVERIFIED; no ownership state was touched or inspected |
| Affected native regression | NOT RUN; no additional exact native test target is authorized |
| Preserve RED | Exact infrastructure receipt retained; no reproduced PRODUCT_RED |
| Scoped files and typed publication | Two authorized new files; single BLOCKED/UNACCEPTED publication follows this report |
| Independent QA/reducer | REQUIRED, not performed or self-accepted |

Only the two exact write-scope files were created. Existing source, fixtures, acceptance SoT, schemas, dependencies, live task/lease/lock/terminal state and foreign changes were not modified. Git status/diff could not be run; review is limited to scoped reads and edit/write receipts. No commit, push, send or deployment occurred.

## Recovery checkpoint and next milestone

Initial useful-work estimate: investigation 3–4 minutes, test/artifact implementation 3–4, meaningful native tests/conditional repair 4–7, compatibility/evidence 2–3 (12–18 useful minutes); infrastructure overhead unknown. Reliable total measured wall time is unavailable from these per-call receipts. No extra work was manufactured to meet a duration target.

1. Infrastructure owner restores a regular local project `.venv` Python executable and pytest/project dependencies required by the authorized runner. This worker has no environment write authority.
2. Coordinator authorizes read-only access to `swe_af/issue/build.py`, its Git/schema/coding-loop dependencies, applicable nested instructions, and the existing nearby native fixtures; supply an exact affected-regression test scope. Clarify rollback semantics only if direct source/fixtures do not settle retained-branch behavior.
3. Resume this same original sticky task/owner, rather than starting a duplicate implementation. Inspect the existing mechanism, extend this exact test into native integration coverage using verified APIs, and exercise positive + two negatives + a real separate-process interruption/restart. Do not treat the present preflight as an acceptance fixture.
4. Run the exact requested pytest command through an authorized Python runner, preserve source/test identities and actual exits, repair only demonstrated owned-fixture defects and retest. Product defects outside write scope require coordinator replanning.
5. Obtain independent source-bound reviewer/reducer QA and authoritative coordinator readback. The immutable publication response is preserved in the tool/conversation receipt after this document; no transport status is predicted here and no acceptance is inferred from it.
