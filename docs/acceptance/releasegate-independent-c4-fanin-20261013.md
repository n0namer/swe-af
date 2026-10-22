# Releasegate independent C4 fan-in QA — 2026-10-13

## Decision

**BLOCKED (static review is RED).** The bounded release-gate test contains useful positive and adversarial assertions, but the authorized runner stopped before pytest execution. No runtime PASS, native AgentField/provider end-to-end result, terminal reducer decision, or authoritative CrewSync terminal readback is available. This packet is reviewer evidence only and does not accept, complete, or mutate any task.

## Authority and gate

- Task: `swe-af:v2-releasegate-c4-independent-fanin-20261013-r1`
- Attempt: `attempt-3e4a8ccb-bb80-4719-87b9-f41e92feff0b`
- Crew session: `manual-20261013-releasegate-c4-fanin-independen-a27372b9b2ca7630`
- Reviewed release wave: `20261013-native-releasegate-hyperwave-r1-swe-af`
- Target Git root: `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af`
- BMad stage/gate: independent release verification after producer execution. Existing acceptance packets are the available project evidence; no competing plan was created.
- Write scope: this QA document only. Production source, tests, scheduler state, leases, provider state, and CrewSync terminal state remain unchanged.

## Frozen evidence and receipt binding

| Evidence | SHA-256 / status | Independent use |
| --- | --- | --- |
| `tests/issue/test_releasegate_git_20261013.py` | `8e9c356a188e43f003286fc07d85b4d4fcc1e6b252f0faca6b02ef2956adf0eb` | Current focused release-gate source and runner identity. |
| `go/internal/dag/checkpoint.go` | `ab3d45da6b49558d9ca65c3308d69cacc55b7b10236b614958df4dd82e2f64d8` | Current checkpoint persistence and recovery contract. |
| `docs/acceptance/releasegate-git-20261013.md` | Existing producer acceptance packet | Declares scripted role calls, real disposable Git effects, and an earlier blocked Python receipt. |
| Focused runner receipt | `BLOCKED / PYTHON_VENV_UNAVAILABLE` | Fresh independent attempt against the current test SHA; pytest did not start. |

The available evidence does not include a physical provider transcript, provider request/response identity, token/cost receipt, native AgentField trace, independently verified worker runtime receipt, or authoritative CrewSync terminal reducer readback for the named hyperwave. Their absence is a release-gate limitation, not proof that they failed or succeeded. The launcher-provided `THREAD_PENDING` binding is not a terminal status.

## C4 context ledger

| Actor / system | Relationship | Trust and acceptance boundary |
| --- | --- | --- |
| Coordinator / independent reviewer | Requests bounded verification and receives a typed candidate report. | Owns neither source implementation nor terminal acceptance in this worker. |
| SWE-AF issue implementation system | Exposes the Python `swe_af.issue.build.implement_issue` entrypoint used by the focused test. | A returned `success` value is producer data; it is not independent acceptance. |
| AgentField role runtime | Supplies coder, reviewer, and verifier behavior through `issue_router.app`. | The focused test substitutes `ScriptedAgent`; therefore provider/runtime behavior is not exercised. |
| Local Git and filesystem | Create worktrees, branches, commits, blobs, artifacts, and cleanup state. | These effects are real inside a disposable repository, but no remote, push, PR, or deployment is involved. |
| Python executable verifier | Imports the branch artifact and asserts the normalizer result. | This is an executable product oracle when pytest runs; the current runner never reached it. |
| Go DAG checkpoint subsystem | Persists and reloads `schemas.DAGState` for recovery. | It is source-related recovery infrastructure, not invoked by the Python release-gate test. |
| CrewSync / terminal reducer | Owns task state and final acceptance. | No administration or terminal readback was performed by this worker. |

## C4 container ledger

| Container | Actual entrypoint / dependency | Observed boundary |
| --- | --- | --- |
| Python issue orchestrator | `issue_build.implement_issue(issue=..., repo_path=..., base_branch=..., artifacts_dir=..., config=...)` | Public asynchronous orchestration is called through `asyncio.run`; the test observes a result mapping. |
| Agent application adapter | `issue_build.issue_router.app = ScriptedAgent`; `await app.call(target, **kwargs)` and `app.note(message, tags)` | The test double preserves target dispatch and notes but replaces native AgentField/provider execution. |
| Disposable Git repository | `run_git(...)` plus subprocess Git object checks | Branch, commit, blob, main-HEAD, clean-tree, remote, and worktree-cleanup assertions are real if the test executes. |
| Executable branch probe | `sys.executable -c "from releasegate import normalize; ..."` | Uses the active test interpreter and fails unless `[3, 1, 3]` becomes `[1, 3, 3]`. |
| Go DAG persistence | `saveCheckpoint`, `loadCheckpoint`, `checkpointRegressesTerminalIssue`, `recoverableIssueIDs` | Uses JSON `DAGState` at `<artifacts>/execution/checkpoint.json`; imported schema path is `github.com/Agent-Field/SWE-AF/go/internal/schemas`. |

## Component and DTO contract ledger

### Python issue boundary

The focused test demonstrates this callable contract by direct invocation:

```text
implement_issue(
  issue: mapping,
  repo_path: string,
  base_branch: string,
  artifacts_dir: string,
  config: {max_coding_iterations: integer}
) -> result mapping
```

Observed issue DTO members are `name`, `title`, `description`, `acceptance_criteria`, and either `files_to_create` or `files_to_modify`.

Observed role request/response contracts are:

- Coder request includes `worktree_path`; response contains `files_changed`, `summary`, and `complete`.
- Reviewer response contains `approved`, `blocking`, and `summary`.
- Verifier request includes `repo_path`; response contains `passed`, `summary`, and `criteria_results[]` with `criterion` and `passed`.
- The agent adapter receives target names ending in `run_coder`, `run_code_reviewer`, or `run_verifier`.

Observed result DTO members asserted by the test are `success`, `outcome`, `verification`, `branch`, `base_sha`, `commits`, `files_changed`, and `iterations`. The test expects `outcome == "completed"` both when verification rejects the source and when it passes; consumers must therefore evaluate typed verification evidence and source identity rather than treating `outcome` alone as acceptance.

### Go checkpoint boundary

The current source exposes package-private functions:

```text
checkpointPath(*schemas.DAGState) string
saveCheckpoint(*schemas.DAGState, noteFunc)
loadCheckpoint(artifactsDir string) *schemas.DAGState
checkpointRegressesTerminalIssue(previous, candidate *schemas.DAGState) bool
recoverableIssueIDs(*schemas.DAGState) []string
```

`saveCheckpoint` rejects a lower `CurrentLevel`, a changed non-empty `GitInitialCommit`, or regression of a terminal issue status. Terminal statuses are case-insensitive `completed`, `failed`, `skipped`, and `accepted`. Writes use `os.WriteFile` after `json.MarshalIndent`; the source explicitly states that the write is not atomic.

## Versioned handshake and cross-module dependency

- No explicit schema/version discriminator is present in the observed Python issue, role response, or result mappings. Compatibility is implicit in dictionary keys.
- No explicit checkpoint format version or writer-fencing epoch is checked by `saveCheckpoint` or `loadCheckpoint`. `DAGState.UnmarshalJSON` defaulting is mentioned for `max_replans`, but that is not a negotiated version handshake.
- The `go-core.worker-result.v1` schema governs this review report only; it does not version the product runtime DTOs above.
- Python crosses into AgentField through the mutable `issue_router.app` adapter. The focused test replaces that boundary, so it proves orchestration sequencing against a test double, not native AgentField transport, policy, tracing, provider execution, or failure semantics.
- Go crosses into AgentField/SWE-AF shared schemas through `github.com/Agent-Field/SWE-AF/go/internal/schemas`. The reviewed checkpoint code persists that shared DTO directly, coupling recovery compatibility to the schema's JSON representation.

## Positive and adversarial checks

### Static positive evidence, not runtime PASS

The current test source requires all of the following when executable:

- a rejected buggy implementation with `success is False` and verifier `passed is False`;
- a salvaged interrupted repair with `failed_unrecoverable`, no verifier call, and one retained commit;
- a restarted repair with `success is True`, verifier `passed is True`, and exactly one iteration;
- distinct 40-character branch commit identities and 40-character blob identities;
- equality of interrupted and restarted repaired blobs, and inequality from the rejected blob;
- unchanged clean `main`, no remote, no generated file on `main`, and removed temporary worktrees.

### Static negative evidence, not runtime PASS

The current test source also rejects or constrains:

- duplicate-dropping `sorted(set(items))` behavior;
- reviewer approval overriding an executable verifier failure;
- reviewer/verifier calls after the simulated interruption;
- duplicate repaired effects (`return sorted(items)` must occur once);
- residual `set(items)` in the repaired source;
- branch reuse across rejected, interrupted, and restarted attempts;
- mutation of the caller's original HEAD or branch;
- leakage of a generated artifact onto `main`;
- presence of a Git remote in the disposable test repository;
- leftover issue worktrees after completion.

### Fresh execution receipt

Requested once through the authorized runner:

```text
go_wave.run_project_tests({file:"tests/issue/test_releasegate_git_20261013.py"})
```

Observed result:

```text
status: BLOCKED
code: PYTHON_VENV_UNAVAILABLE
engine: python:pytest
source_identity_sha256: 8e9c356a188e43f003286fc07d85b4d4fcc1e6b252f0faca6b02ef2956adf0eb
exit_code: null
stdout: ""
stderr: A regular local .venv Python executable is required.
```

The denial occurred before pytest collection or execution. Runtime positive count is zero, runtime negative count is zero, and no PASS is recorded. The identical denied call was not retried.

## Adversarial findings

1. The final focused test SHA has no executable receipt because the required regular `.venv` is unavailable.
2. Scripted role calls replace native AgentField/provider behavior, so provider authentication, transport, policy, tracing, cancellation, timeout, and malformed-response behavior remain unproven.
3. No physical provider or worker receipt is available in the permitted evidence to bind a native run to the reviewed source SHA.
4. No authoritative CrewSync terminal reducer readback is available; `THREAD_PENDING` cannot support a terminal release claim.
5. The result uses `outcome == "completed"` for both rejected and accepted verification states, creating a fail-open risk for consumers that ignore `verification.passed`.
6. Product DTOs are unversioned dictionaries, so missing, renamed, null, or type-changed fields have no explicit compatibility negotiation.
7. Checkpoint JSON has no explicit format version and no migration or fail-closed compatibility gate.
8. Checkpoint persistence is explicitly non-atomic, leaving torn-write/crash behavior outside the demonstrated release-gate path.
9. Checkpoint source protection uses `GitInitialCommit` but no writer fencing epoch or durable compare-and-swap ownership proof is visible in the reviewed component.
10. The Python release-gate test and Go checkpoint subsystem are not executed together, so the claimed recovery fan-in remains a document-level correlation rather than a native cross-module integration proof.
11. The executable verifier runs inside the test-controlled branch environment; it proves the normalizer criterion only, not arbitrary issue acceptance criteria or isolation from malicious branch code.
12. Static assertions for clean `main`, no remote, and worktree cleanup are valuable but have no fresh runtime evidence for the current SHA.

## Scoped repair and affected-only replan

Original sticky owners should make only evidence-producing or defect-specific repairs:

1. **Runner owner:** materialize the repository's regular `.venv`, then run `.venv\Scripts\python.exe -m pytest -q tests/issue/test_releasegate_git_20261013.py` unchanged and preserve full output, exit code, duration, and source SHA.
2. **Issue-orchestrator owner, only if the focused test fails:** repair the exact failing producer assertion in `swe_af.issue.build` or its existing Git helper; do not weaken verifier, source-identity, cleanup, or no-duplicate assertions.
3. **Runtime/AgentField owner:** execute a separate native, policy-compliant AgentField/provider scenario bound to immutable source and physical worker/provider receipts. Do not relabel `ScriptedAgent` evidence as native E2E.
4. **Checkpoint owner:** retain the already identified narrow boundaries for explicit checkpoint versioning, atomic/fenced persistence, duplicate-resume prevention, and source/effect receipt reconciliation; validate them in focused Go tests from the actual `go` module root.
5. **Reducer/coordinator:** correlate the focused Python receipt, any native AgentField receipt, source identities, worker publication, and terminal CrewSync readback before issuing any full release decision.

If the Python test passes unchanged, this packet may advance only to **ACCEPT_SCOPED** for the disposable-Git scripted-role boundary. Full GREEN still requires separately proven native AgentField/provider E2E and terminal reducer acceptance. If it fails, route the exact failure to the original sticky source owner and rerun only affected checks before broader regression.

## Self-review

- Only the authorized QA document was created.
- No production source or prior producer/QA artifact was edited.
- Source hashes are bound to direct scoped reads or the fresh runner receipt.
- Static source observations are labeled separately from runtime evidence.
- No mock, process state, publication transport, or producer `success` field is treated as product acceptance.
- Verdict remains **BLOCKED**, acceptance remains external, and the repair routes are scoped to original owners.
