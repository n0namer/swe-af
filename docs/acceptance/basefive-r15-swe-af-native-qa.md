# Basefive R15 — independent native Git QA

## Disposition

**FACT: INFRASTRUCTURE BLOCKED / UNACCEPTED.** The independent test artifact was created, but the authorized runner did not start Python or pytest. There is no runtime PASS, product acceptance, or terminal readback. This report supplements, and does not replace, `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`.

## Authority and source receipts

- Target: `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af`; not canonical BMAD-MNNZ.
- Task: `swe-af:v2-r15-independent-native-mvp-qa-20261014-r1`.
- Attempt: `attempt-af73684a-ff76-4094-91ba-5f2f07c91b65`.
- Crew session: `manual-20261014-basefive-r15-native-qa-swe-af-v-62c9239c66f687fd`.
- Coordinator-supplied Git receipt: `064a7226a500365005b5b662044fc9014a328722`; not independently refreshed with Git.
- Manifest receipt: `eee02607379292b659c65ae2045e1e4b8e8ac8601d0f15138a1d0f57eefe2315`.
- Task-data receipt: `9bf0e0172d31c7465f456b4e1c97dce71ec7df5d950f7eee8b765644df4a41a4`.
- Direct current source SHA-256, `swe_af/issue/build.py`: `440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd`.
- Direct contract SHA-256, `docs/architecture/basefive-r14-contracts.json`: `cc0eeb34cbc3031785787caf1e72ce8393551533376d46403b71a393c3330029`.
- Final test SHA-256, `tests/issue/test_northstar_r15_independent_git_qa.py`: `09aaab258ceb95e28569e19cb342f4bfdfba885243debd987e375dbeaaaccd6a`.
- Read `AGENTS.md` SHA-256: `4a8876425a07de89597f3b56c6ba918b8dd89ddbce1abb72b4f9db9ec1782ae9`.
- Active reducer SHA-256 supplied in the trusted binding: `feeb1fe6a88e53770831de890fe3674e6425dd40327deb913fbf73f3a01f1f17`; its content was read directly.

These identify the inspected source and submitted test bytes, not an immutable snapshot of the entire dirty checkout. Branch, whole-checkout dirty state, dependency files, baseline test identities and their hashes were not available through this exact read scope. No broader source-state identity is claimed.

## Existing mechanism and C4 boundary

**FACT:** Source inspection found the mechanism already implemented. Production code needs no speculative replacement for this QA attempt:

- `_checkpoint_if_dirty` scrubs tracked junk and commits only when Git reports dirty work.
- `_checkpoint_interrupted_issue` preserves checkpoint errors as notes rather than masking the original interruption.
- `_implement_issue_impl` creates a UUID-scoped branch/worktree from the committed base, checkpoints successful or interrupted work, salvages commits on exceptions, and cleans the worktree in `finally`.
- Cancellation has a distinct `asyncio.CancelledError` path, checkpoints, salvages commits, then re-raises.
- Success requires completed coding, at least one commit, and a passing verifier when verification is present.

Actual public Python signature, read from `swe_af/issue/build.py:425`:

```python
async def implement_issue(
    issue: dict,
    repo_path: str,
    base_branch: str = "",
    artifacts_dir: str = ".artifacts",
    additional_context: str = "",
    config: dict | None = None,
) -> dict:
```

The test uses this public entrypoint, not a replacement issue orchestrator. Its independent AST oracle checks parameter names, annotations, defaults and return type; its hash oracle rejects a changed source instead of silently testing another implementation. The source checks result construction through `IssueBuildResult(...).model_dump()` and input validation through `IssueSpec.model_validate(issue)`.

**FACT:** R14 describes Go `func Build(ctx context.Context, deps *Deps, input map[string]any) (any, error)` as a different entrypoint, not Python parity. It records no wire schema version and labels compatibility `UNVERIFIED`. This is secondary contract evidence, not a fresh Go source inspection. `go/internal/orch/build.go`, its checkpoint dependencies, and DTO definitions are outside the authorized read scope. Go checkpoint, stale-owner fencing and cross-language wire compatibility gates remain BLOCKED; no interface/version parity is inferred.

## Local measurable DoD and implementation

Value: an independent, source-bound regression oracle that can distinguish a real Git repair from rejected, interrupted or cancelled issue execution while protecting the caller checkout.

Constraint: only the new test and this acceptance report may be written; no production changes, live claims, leases, terminal state, network or external sends.

Smallest design: retain the native public issue orchestrator and native Git/checkpoint/verification/result paths. Substitute only the coder dispatch boundary and router application with deterministic in-process fixtures. Run an independent `sys.executable -I -c` assertion against real disposable Git worktree content. All fixture Git commits and worktree operations target pytest temporary repositories, never the project checkout. Unexpected agent targets fail the test; GitHub PR creation is disabled. No remote is configured.

The new test has eight collected-case definitions including the three-case rejection parametrization:

| Case | Independent assertions specified | Execution |
|---|---|---|
| Source/C4 contract | Source SHA, R14 source binding, exact AST signature/defaults/annotations, honest no-wire-version/different-entrypoint status; prints source/test/contract hashes | BLOCKED |
| Native red → repair → green | Independent Python assertion fails on committed broken input, passes after repair; native checkpoint produces exactly one deliverable commit; native verifier succeeds | BLOCKED |
| Coding rejection | Even a repaired, committed artifact cannot make a rejected coding outcome successful; verifier is not called | BLOCKED |
| Verification rejection | Coding completion plus a commit with still-invalid content cannot produce success | BLOCKED |
| Verifier unavailable | Exception from verifier is represented as failed verification, never success | BLOCKED |
| Interruption salvage and fresh replay | RuntimeError after a tracked edit produces a failed result with one salvage commit; fresh invocation from main repairs the same blob on a distinct branch and leaves salvage unchanged | BLOCKED |
| Cancellation | Event-synchronized real asyncio task cancellation propagates only after tracked edit salvage; no verifier call; temporary worktree removed | BLOCKED |
| Abandoned worktree isolation | Existing unrelated dirty worktree and branch remain unchanged while the new issue uses another worktree | BLOCKED |

Shared safety assertions protect caller main branch, HEAD, dirty tracked content, empty remotes and expected worktree count. Git blob IDs, commit graph/count and independent subprocess exit codes are the oracles, rather than copying native success computation.

**Limits:** The coder/reviewer service itself is substituted; this is native issue orchestration/Git QA, not live agent-quality E2E proof. The interruption is a caught Python exception and replay is a fresh invocation, not an OS process kill or durable checkpoint resume. The abandoned worktree test is isolation, not stale-owner lease fencing. These stronger scenarios require additional authorized source/test evidence and a working native runner. They must not be promoted to PASS from these assertions.

## Exact focused-run receipt

Requested validation: `python -m pytest -q tests/issue/test_northstar_r15_independent_git_qa.py`.

Actual authorized invocation: `go_wave.run_project_tests({"file":"tests/issue/test_northstar_r15_independent_git_qa.py"})`.

```json
{
  "status": "BLOCKED",
  "code": "PYTHON_VENV_UNAVAILABLE",
  "test": "tests/issue/test_northstar_r15_independent_git_qa.py",
  "engine": "python:pytest",
  "source_identity_sha256": "09aaab258ceb95e28569e19cb342f4bfdfba885243debd987e375dbeaaaccd6a",
  "exit_code": null,
  "signal": null,
  "error": null,
  "duration_limit_ms": 120000,
  "output_truncated": false,
  "stdout": "",
  "stderr": "A regular local .venv Python executable is required."
}
```

No interpreter started, no collection occurred, no scenario ran, and no failed product assertion was observed. Syntax, imports, runtime callback-envelope compatibility and all assertions remain unverified. The final test hash equals the denied runner's source identity; it was not changed after this receipt. There was one runner call, with no identical retry, global-interpreter workaround, dependency installation or shell bypass.

Affected baseline regression is NOT RUN: baseline test paths are not in this task's read/test scope, and the same Python infrastructure is unavailable. New scenarios cover affected checkpoint/verification/cleanup behaviors by design but do not substitute for baseline execution evidence.

## Source review and method boundaries

Source-only self-review checked that edits are limited to the two authorized new artifacts, production checkpoint helpers are reused, success/failure assertions are independently grounded in Git/process output, callback sends are fenced, fixture subprocesses use argument lists without a shell and finite timeouts, and caller dirty content is distinct from the committed fixture input. Both target files were absent on initial exact reads. No existing foreign edits were replaced.

BMad Help resolved the target and current gate from the existing C4 reducer: implementation exists, native runtime QA and reduction remain blocked. Quick Dev supplied a bounded implementation/self-review method; Verification Before Completion requires fresh execution before any success claim. Config resolution, unrestricted Git status/history, graph navigation, sprint-status changes, review subagents, new workflow spec files, commits and editor commands were unavailable or outside task authority; those workflow steps were not performed. No complete BMad workflow or independent review is claimed.

Initial useful-work sizing: investigation 5m + implementation 8m + testing 4m + repair allowance 4m + compatibility/self-review 3m + evidence 3m = 27m; external coordination allowance 2m, estimated efficiency 27/29. The infrastructure denial short-circuits meaningful runtime/repair work. Actual elapsed duration was not measured, and no time was padded to meet the estimate.

## Acceptance reduction and same-attempt recovery

- D1: direct Python signature and scoped source/test fingerprints AVAILABLE; exact full dirty-checkout snapshot and fresh Go/DTO verification BLOCKED.
- D2: positive, three negative, cancellation and exception-replay assertions WRITTEN, not executed; real process-crash/restart and stale-owner fencing NOT PROVEN.
- D3: independent new oracle WRITTEN inside the authorized test; runtime independence/callback compatibility UNVERIFIED.
- D4: focused native runner INFRASTRUCTURE BLOCKED; affected baseline NOT RUN; no product failures available for bounded repair.
- D5: typed BLOCKED/UNACCEPTED publication required; independent reducer and authoritative terminal readback remain coordinator-owned and pending.

Next move: coordinator/operator restores a regular project `.venv` with the project's native pytest dependencies, then resumes this same task/owner and runs the unchanged test through `go_wave.run_project_tests`. Do not create a duplicate implementation attempt or rerun the identical denied call before restoration. The operator-side equivalent is `.venv\Scripts\python.exe -m pytest -q tests/issue/test_northstar_r15_independent_git_qa.py`; the managed runner remains authoritative here. Grant exact affected baseline test paths and Go/checkpoint/owner/DTO read evidence before extending those proofs. Diagnose any resulting test failure in this same scoped test, without production writes. After runtime evidence, independent QA/reducer and coordinator terminal readback decide acceptance.

The publication response is preserved in the final tool receipt and user handoff, not predicted here. Transport verification alone cannot change this UNACCEPTED disposition.

## Resumed verification checkpoint — 2026-10-15

**FACT: INFRASTRUCTURE BLOCKED / UNACCEPTED remains the disposition.** This section records the current attempt without replacing the historical receipt above.

- Task: `swe-af:v2-r15-independent-native-mvp-qa-20261014-r1`.
- Current attempt: `attempt-9bfe5fcb-7dd9-48f8-bd38-5cdbc64039e5`.
- Current session: `manual-20261015-basefive-goal20-swe-af-v2-r15-i-313bada68cdbe753`.
- Prelaunch Git receipt: `064a7226a500365005b5b662044fc9014a328722`; no independent Git status, branch or dirty snapshot was available.
- Current manifest receipt: `8a2abf4b31d4111628fb8aa718da720fd6c95ce701e3a3920a5948a97739268b`.
- Current task-context receipt: `2f8e1fcbad336d2d332f4dd74508c1cecca11c7d07167849ec9485ca1c99ead2`.
- Current project-cache receipt: `d9d1ecc2603d6c6b235bc81cfa5a4d6084cd9b1b5d9ee39c121b069b76fa2269`.

BMad Help was loaded first. The explicit target remains swe-af, and the existing C4 reducer remains the source of truth: native QA and independent reduction are blocked, not accepted. Verification Before Completion was selected because an existing QA implementation was discovered in the exact write scope. No duplicate test implementation or speculative production repair was justified. The permitted source read directly returned SHA-256 `440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd`; its public signature at line 425 matches the Python signature recorded above. R14 still explicitly describes different Python/Go entrypoints and unverified wire compatibility. Actual Go checkpoint source and baseline test paths remain outside the supplied read/test scope. No graph tools or unrestricted Git/config workflow were used.

Local DoD for this resumption: reuse the existing oracle; obtain a real native focused receipt for the unchanged test; repair only concrete in-scope test defects if execution exposes them; review compatibility and source fingerprints; preserve exact blockers and publish one UNACCEPTED result. Full acceptance still requires actual positive/negative/restart execution, an exact dirty-checkout identity, affected baseline regression, authorized Go/checkpoint/owner evidence, independent reduction and terminal readback. A missing target was not inferred; both authorized artifacts already existed.

One current-attempt invocation of `go_wave.run_project_tests({"file":"tests/issue/test_northstar_r15_independent_git_qa.py"})` returned:

```json
{
  "status": "BLOCKED",
  "code": "PYTHON_VENV_UNAVAILABLE",
  "test": "tests/issue/test_northstar_r15_independent_git_qa.py",
  "engine": "python:pytest",
  "source_identity_sha256": "09aaab258ceb95e28569e19cb342f4bfdfba885243debd987e375dbeaaaccd6a",
  "exit_code": null,
  "signal": null,
  "error": null,
  "duration_limit_ms": 120000,
  "output_truncated": false,
  "stdout": "",
  "stderr": "A regular local .venv Python executable is required."
}
```

**FACT:** No interpreter started, no tests collected or ran, and no product failure was observed. The runner's test fingerprint exactly matches the prior artifact receipt. The test was not edited in this attempt. There was no second runner call, fallback interpreter, shell bypass, installation, production write, remote/send, task/lease mutation or self-acceptance. Baseline regression remains NOT RUN, not PASS. Exception replay is not process-crash recovery; abandoned-worktree isolation is not stale-owner fencing. The source/test pair is identified, not the entire dirty checkout.

Self-review: only this checkpoint was appended to the authorized acceptance document; historical evidence and existing test bytes were preserved. Initial resumed-work estimate was 5–8 useful minutes for scoped investigation, one verification attempt, source review and evidence if blocked; execution/repair effort cannot be usefully estimated until the native runner works. Actual elapsed time was not measured. No artificial work was added to meet a duration target.

Recovery: operator restores the regular local project `.venv` and its pytest dependencies, then resumes this same task/thread for the unchanged focused test via the managed runner. Operator equivalent: `.venv\Scripts\python.exe -m pytest -q tests/issue/test_northstar_r15_independent_git_qa.py`. Do not repeat the denied invocation until restoration is confirmed. Coordinator must supply exact affected baseline test authorization and Go/checkpoint/owner read evidence for the remaining gates; the worker must not broaden scope. After executable evidence, independent reviewer/reducer and coordinator terminal readback remain required. The single current publication response is retained in the tool receipt and final handoff; it cannot establish acceptance.
