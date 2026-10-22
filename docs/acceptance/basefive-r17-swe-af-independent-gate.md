# Basefive R17 swe-af independent gate — 2026-10-15

FACT: the original 2026-10-14 attempt evidence is preserved below as history. The current authorized attempt is `attempt-da738521-8f2d-4907-a76d-83274345011a`; its binding, repairs and fresh tool receipts are in the 2026-10-15 addendum. Historical claims about artifact contents or execution must not be attributed to the repaired final SHA.

## Historical reducer verdict — 2026-10-14

**BLOCKED / UNACCEPTED. No ACCEPTANCE_RECOMMENDED or ACCEPTED_GREEN recommendation.**

FACT: source identity can be cross-checked against the permitted R14 contract. Neither original R16 typed receipt is supplied in the permitted evidence. Native producer checks were denied by the scoped runner, and the new independent QA check was blocked by missing Python infrastructure. Therefore this attempt cannot prove the full local DoD. No task, lease, terminal state, production source, remote, send, or deployment was changed.

The source-backed trust-boundary finding below is not an executed QA_RED receipt. Product, fixture, evidence-contract and infrastructure gaps remain distinct.

## Binding and Source of Truth

- Project/root: `swe-af`, `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af`.
- Original task: `swe-af:v2-r17-acceptance-reducer-20261014-r1`.
- Attempt: `attempt-dbffb585-aa82-4c09-9b98-820cbff924a1`.
- Session: `manual-20261014-basefive-r17-review-swe-af-v2-r-71d4ca1d8e6df591`.
- Git receipt: `064a7226a500365005b5b662044fc9014a328722` (trusted prelaunch identity, not a fresh Git status).
- Manifest SHA: `1b4db7837eb4744228ee2f2695ef9133c02c7797522b50aa1e1a070bf22401d4`.
- Active BMad SoT: `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`.
- Stage/current gate: independent QA/reduction following R16 producers; native runtime proof, exact original receipt binding and coordinator terminal CAS/readback remain required.
- BMad Help, adversarial review and verification-before-completion were loaded. Canonical BMAD-MNNZ is the workflow source, not this project. BMad config resolver, Git inspection and graph MCP methods were not available through this attempt's authorized interfaces/read scope; no execution or graph freshness claim is made.

## Bounded design and DoD

Value: prevent acceptance based on an oracle, stale SHA, truthy malformed verifier output or transport-only receipt. Smallest compatible move: reuse R16 disposable Git fixtures and the public Python entrypoint, add independent assertions in the one authorized QA file, and record this reducer verdict without producer edits.

Required DoD: both original typed R16 receipts bound to immutable source/test SHA; independent real-native positive, two negatives and restart with no duplicate effects; exact code/test/log identity; architecture-anchor review; classified reds; coordinator terminal CAS before any positive terminal recommendation. This is not met.

Estimated useful scope: approximately 15–20 minutes for investigation, bounded QA implementation, focused checks, source review and evidence; failed native admission reduces executable validation work. Actual elapsed work was not instrumented. Next meaningful milestone is an admitted native check with immutable stdout/stderr and original receipt binding, not another implementation launch.

## Independently observed file identity

Hashes below are exact `read_project_file` SHA receipts; they are not original producer result receipts.

| File | SHA-256 | Reduction |
|---|---|---|
| `swe_af/issue/build.py` | `440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd` | Matches R14 contract and Python R16 test pin |
| `go/internal/orch/build.go` | `05c39a2e3b41869b01a1482543c9875630a0777fbea24ea706c05727fe640cbc` | Matches R14 contract and Go R16 test pin |
| `tests/issue/test_northstar_r16_git_cycle.py` | `3a550d5d87bd68bb71d70a8d7a86cde97687726f06a6a651f3660347fbb1e648` | Recorded independently; original typed receipt unavailable |
| `go/internal/orch/northstar_r16_resume_canary_test.go` | `8356d7644b233c422eed1bd215d7bcfe8e0c315e37dc10132cd9d9f1a0c1182c` | Recorded independently; original typed receipt unavailable |
| `tests/issue/test_northstar_r17_acceptance_qa.py` | `c0966e900640276d9f002f9d1da88cd255c3dc1a48816f555c5b9b8d6e471473` | Exact runner source identity; not executed |

The R14 contract records Graphify SHA `e18cd58bd72e8975c07f305a988192227c529a39ab35752388ad0db706f5c277`, ArchSteer model SHA `ed964b0a700a566bd9e42caa55551baee873706aa500fbd57c446a28b7971260`, and ArchSteer architecture SHA `7d7ee7046812a9c3c5da1d3b7484015432c4cb3b3cb106a6d34e8bffeb65a3d2`. FACT: these are inherited contract entries, not independently hashed raw artifacts in this attempt. Raw artifacts are outside read_scope. The contract explicitly records graph freshness UNVERIFIED, zero ArchSteer Go file coverage, missing Python graph symbols and DTO wire compatibility UNVERIFIED. No graph-completeness or cross-language parity claim follows.

## Independent QA artifact

`tests/issue/test_northstar_r17_acceptance_qa.py` adds:

- Four immutable source/test pins and recursive checks of C4 source anchors for the two permitted production sources, without modifying the contract.
- Positive public-entrypoint Git delivery with an actual isolated Python verifier subprocess, one coding invocation, one verifier invocation, one checkpoint, no PR, caller isolation and cleanup assertions.
- Negative missing/rejected verifier output cases and a separate malformed non-boolean approval boundary (`"false"`, `"true"`, integer, list and object values).
- Cancellation salvage followed by two no-op restarts; stable branch head/blob, one total repair commit, no duplicate verifier call and deletion of empty replay branches.

FACT: these are authored assertions, not passing results. Coding dispatch remains a fixture replacement reused from R16; the local verifier subprocess would be real if executed, but this does not establish native provider/task/result behavior. No test is weakened, skipped or xfailed to conceal the malformed-verdict finding.

## Exact test/error receipts

### Original Python R16 validation

Requested manifest command: `python -m pytest -q tests/issue/test_northstar_r16_git_cycle.py`.

Actual authorized tool request: `run_project_tests({"file":"tests/issue/test_northstar_r16_git_cycle.py"})`.

```text
GO_WAVE_NODE_ERROR:test_outside_task_write_scope
```

Classification: **INFRA/admission BLOCKED**. No command execution, exit code, test count or native log SHA was returned. read_scope is not test_scope.

### Original Go R16 validation

Requested manifest command from the Go module: `go test ./internal/orch -run TestR16Resume -count=1`.

Actual authorized tool request: `run_project_tests({"file":"go/internal/orch/northstar_r16_resume_canary_test.go"})`.

```text
GO_WAVE_NODE_ERROR:test_outside_task_write_scope
```

Classification: **INFRA/admission BLOCKED**. No Go execution, exit code, test count or native log SHA was returned. This error does not establish whether the underlying runner supports Go once admitted.

### Independent R17 QA validation

Actual authorized request: `run_project_tests({"file":"tests/issue/test_northstar_r17_acceptance_qa.py"})`.

```json
{"status":"BLOCKED","code":"PYTHON_VENV_UNAVAILABLE","test":"tests/issue/test_northstar_r17_acceptance_qa.py","engine":"python:pytest","source_identity_sha256":"c0966e900640276d9f002f9d1da88cd255c3dc1a48816f555c5b9b8d6e471473","exit_code":null,"signal":null,"error":null,"duration_limit_ms":120000,"output_truncated":false,"stdout":"","stderr":"A regular local .venv Python executable is required."}
```

Classification: **INFRA BLOCKED**, not PASS or product FAIL. Identical blocked calls were not retried. No shell fallback, environment installation or runner bypass is authorized.

## Adversarial findings and classifications

- Neither original R16 typed worker receipt, exact original task/attempt/session binding nor publication readback is present in the permitted evidence; original producer receipt cross-verification is evidence-contract BLOCKED.
- Hashing the current R16 tests proves file identity, not their historical execution or binding to a producer publication.
- Python R16 coding and verifier boundaries are patched fixtures; real disposable Git alone does not prove native agent/provider task-to-result execution.
- Go `TestR16ResumeOracle*` methods exercise `r16ResumeOracle`, a test-local JSON-file state machine; they do not invoke production `Build`, `Deps`, native task dispatch or a production restart API. This is a fixture coverage gap, not a demonstrated production restart failure.
- Go `TestR16ResumeSourceContract` verifies SHA/signature/registration, not production Build execution; a green oracle suite would still not satisfy the real-native requirement.
- The Python R16 validation was not admitted, so its positive/negative/restart assertions have no fresh runtime receipt here.
- The Go R16 validation was not admitted, so no native or oracle runtime result can be reported here.
- Independent QA cannot collect because the runner requires a regular project `.venv`; no passing QA or red/green cycle can be inferred.
- Python `_run_verification` accepts a dictionary unchanged and final success applies `bool(verification.get("passed"))`. FACT: nonempty `"false"` is truthy under this source logic. HYPOTHESIS for executable QA: malformed non-boolean approvals authorize success when coding and commits pass. This is a source-backed product trust-boundary concern, not a runtime failure receipt; any production repair belongs outside this task's write_scope.
- R14 DTO wire version is null and compatibility is UNVERIFIED; differing Python issue and Go full-build entrypoints are not parity proof.
- Raw C4/Graphify/ArchSteer artifacts are not available in scope; inherited graph hashes and recorded coverage counts cannot establish current graph freshness or complete Go source anchors.
- No exact CrewSync terminal CAS/readback is available; neither future passing checks nor VERIFIED publication transport could substitute for coordinator terminal authority.

## Self-review and affected-only recovery

Only the two authorized files were created. Existing production sources, R16 tests, architecture contract and prior reducer SoT were left unchanged. The QA uses existing fixture setup, no new scheduler or replacement native subsystem, disposable local Git, no remote, and no send/deploy action. Source review identified no reason to repair or weaken the producer tests in this scope. Git dirty-state inspection and actual Python syntax/collection/runtime checks remain unverified because their authorized execution capability is absent or blocked.

The coordinator should restore the regular project `.venv`, authorize the existing R16 files as test_scope, expose the two existing original typed receipts plus immutable log/SHA evidence and terminal CAS/readback, and resume this same original sticky attempt. Then run exactly the two manifest validation commands and `python -m pytest -q tests/issue/test_northstar_r17_acceptance_qa.py`. Do not create duplicate R16 producers or reinterpret the Go oracle as native proof. If malformed-verdict QA fails, preserve the exact failure and route only the affected Python production source repair to its existing authorized owner; this QA task cannot edit it. A production-native Go canary requires the existing Go owner's approved scope and must be re-evaluated before any positive recommendation.

Publication is a separate typed UNACCEPTED evidence receipt. The exact publisher response/readback is retained in the tool transcript and final handoff, not preclaimed by this document. No terminal acceptance is performed.
