# R14 architecture drift analysis and proposed ADR

Date: 2026-10-14. Disposition: **CANDIDATE / UNACCEPTED**; independent QA/reducer owns acceptance.

## Evidence and local Definition of Done

FACT: The active BMad acceptance source is `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`; it keeps the native ReleaseGate scenario blocked. This bounded slice supplies three C4 views, source fingerprints, observed API/DTO contract JSON, complete comparisons of the two supplied index snapshots, and deterministic Node positive/negative checks. It does not rerun or replace the native Python runtime gate.

FACT: Before implementation, all six authorized R14 targets returned ENOENT. The first Node inventory test failed because the contract JSON was absent while successfully reporting the inspected source/index inventory. Files are newly authored, not replacements of an accepted predecessor candidate.

The measurable local DoD is: exact seven permitted input hashes; actual Python and Go public signatures and DTO call-site observations; every declared diagram node/edge anchored in inspected source; honest Go/package and checkpoint blind spots; negative mutation checks; final focused test receipt; one typed UNACCEPTED result. No imported source, generated index, live task, claim, or production state is changed.

## Index comparison

| Observation | ArchSteer | Graphify | Meaning |
|---|---:|---:|---|
| Components / nodes | 288 components | 9,256 nodes | Different units; not comparable coverage percentages |
| Edges | Not used as a parity count | 24,563 | Index snapshot only, not live call traffic |
| Distinct node `source_file` paths | Not applicable | 542 | Exact paths, not aliases normalized into canonical paths |
| `.go` file paths | 0 | 200 | ArchSteer snapshot has a verified Go blind spot |
| Go parent directories | 0 | 36 | Indexed directories, not proof of all buildable packages |
| `swe_af/issue/build.py` | Six `exported_apis` names | 10 nodes | Node count includes file/docstring/nested-function nodes |
| `go/internal/orch/build.go` | Absent | 14 nodes | Build, Handlers, buildInput and helpers are indexed |
| `.graphify-code/` alias paths | 69 components | 0 node paths | Mixed snapshot population; do not deduplicate by basename |

FACT: ArchSteer names `_checkpoint_if_dirty` and `_checkpoint_interrupted_issue`, and both exist in the inspected Python source. Graphify has no corresponding function labels in nodes whose `source_file` is exactly `swe_af/issue/build.py`. This is a bounded symbol-level omission, not a claim that no graph edge elsewhere mentions checkpointing.

FACT: The inspected Go source directly imports five internal package paths: `config`, `envelope`, `hitl`, `schemas`, and `workspace`. Graphify has `.go` nodes under each corresponding `go/internal/` directory; ArchSteer has no such Go components. Their implementations, package declarations, build tags, compile results, and runtime behavior are outside the permitted source scope. `agentfield/sdk/go/agent` is an external dependency; it is not counted as an inspected internal package.

FACT: `.archsteer/model.json` reports schema version `1.0.0`, repository `swe-af`, short commit `064a722`, and timestamp `2026-10-14T15:43:32.432127+00:00`. These are producer metadata, not DTO versions or evidence of generation from the current dirty source state. Neither a matching commit prefix nor hashing the generated file verifies each indexed source fingerprint. Graph source freshness and repository-wide completeness remain **UNVERIFIED**. Index graph-navigation MCP tools are not exposed in this worker; complete bounded JSON analysis uses the authorized Node runner instead.

## Actual API and DTO comparison

| Boundary | Observed source | Limitation |
|---|---|---|
| Python entry | `implement_issue(issue, repo_path, base_branch="", artifacts_dir=".artifacts", additional_context="", config=None) -> dict` | Async Python issue API; no planning agents in this flow |
| Python input | `IssueSpec.model_validate(issue)` and `IssueBuildConfig(**(config or {}))` | DTO definitions/defaults/version are outside scope |
| Python output | 16 named `IssueBuildResult` constructor fields, then `.model_dump()` | Constructor names observed; serializer aliases and schema version unverified |
| Go entry | `Build(ctx context.Context, deps *Deps, input map[string]any) (any, error)` | Registered by this file's `Handlers()` as `build`; whole-node registration unverified |
| Go input | Ten exact JSON tags in `buildInput` | Marshal/unmarshal errors are ignored; no strict-boundary-validation claim |
| Go output | Seven `schemas.BuildResult` constructor fields; conditional `workspace_handle` map entry | Go field names are not asserted to be JSON wire names; definitions outside scope |

FACT: Python's issue API is not the counterpart of Go's goal-based full-build API. Both orchestrate coding/verification, but input shapes, planning, success logic, side effects, and result constructors differ. The inspected Go file's comment says it ports Python `build()` in `app.py`; that file is outside scope. Graphify indexes `go/internal/issue/` paths, but inspecting them is not authorized. No issue-level Python/Go parity claim follows.

## Behavioral drift and safety boundaries

- FACT: Python success is completed coding (including debt), nonempty commits, and verifier pass when verification is present. Verification is optional and gated by `cfg.verify`, coding completion, and commits.
- FACT: Python catches cancellation after execution starts, checkpoints interrupted work, collects commits, and re-raises; generic errors salvage committed work. Cleanup is in `finally`. This describes control flow, not tested interruption/restart safety.
- FACT: Python's PR gate checks `enable_github_pr`, coding completion, and commits, but not verifier success. Do not draw a verifier-gated PR edge or assume verifier rejection prevents a configured PR.
- FACT: Go success derives from `verification["passed"]`. The empty-build guard returns `agent.ReasonerFailed` with the result map when available. Node static checks do not execute the SDK's failure transport.
- FACT: The reducer's prior native test receipt is `PYTHON_VENV_UNAVAILABLE`, `exit_code: null`; the documented final Python test SHA was not executed. This attempt does not restore infrastructure or promote that receipt to PASS.

## Proposed ADR R14 — source-bound, non-promoting architecture contracts

Status: **PROPOSED, not architecture acceptance**.

Context: Generated indexes have materially different populations and cannot by themselves establish API/DTO compatibility or restart safety.

Decision proposal: Keep three readable C4 views and one versioned documentation contract. Pin the exact permitted source/index bytes. Use actual public declarations, JSON tags, result-constructor names, and source anchors as a deterministic drift gate. Treat missing languages/symbols as visible gaps; keep schema versions, deployment topology, complete package coverage, and runtime parity unverified until the appropriate owner supplies evidence.

Consequences: Source changes fail the focused gate until a reviewed rebaseline updates the contract. Generated graph refreshes also require review. A source-contract PASS is deliberately weaker than runtime acceptance. The Mermaid checker validates the documented restricted flowchart syntax and exact declared nodes/edges; no renderer execution is claimed. This adds no production dependency or new scheduling/acceptance subsystem.

## Validation and workflow receipt

Focused command: `node --test docs/architecture/basefive-r14-contract.test.mjs` via `go_wave.run_project_tests`.

Focused result: **PASS, 37/37 tests, 0 failures, exit_code 0, node:test**, receipt `2026-10-14T16:05:23.508Z`, TAP duration `663.6197 ms`. The suite includes 34 negative checks and three positive/source-artifact checks. Initial inventory runs failed with exit 1 because the contract JSON was absent, not because of infrastructure. The first integrated run was 29/31 PASS: a mistaken `deps.CallRaw(ctx, "plan"` source anchor failed the positive check and masked one diagram-label negative check. Inspection of `go/internal/orch/build.go:183` showed `deps.App.Call(ctx, deps.target("plan"), planKwargs)`. The anchor and dispatch inventory were repaired; retest passed 31/31. Self-review then clarified delegation-before-validation wording, required a valid pristine contract before mutation checks, expanded drift rejection to all seven inputs, and added six-file SHA/line-limit receipts; the expanded suite passed 37/37. Publication evidence must use the final rerun's artifact fingerprints after this documentation update.

Self-review: all changes are the six new authorized documentation/test targets; no production code, imported DTO definition, generated index, acceptance source, dependency, lease, or task state is modified. All six artifacts are below 500 lines. The three views contain 21 total nodes and 20 total edges checked against source anchors and exact contract labels. No Python/Go parity, wire-version, renderer, repository-completeness, Git-cleanliness, or runtime success claim is made.

Method limitation: BMad Help and Quick Dev instructions were loaded. Config resolvers, Git status/history, out-of-scope spec/sprint writes, interactive approval, and delegated review were not executed because the immutable worker contract permits only scoped MCP reads/writes/tests and forbids agents/shell/task mutation. The explicit task contract supplies the local DoD. No full BMad workflow-completion claim is made.

Useful-work estimate: investigation 5 min + artifacts 8 min + positive/negative tests 5 min + repair allowance 3 min + self-review/evidence 3 min = 24 min; external coordination overhead estimated 0–1 min. This is planning, not a claim of measured elapsed time or a duration floor.

## Remaining gate and next move

Independent QA should rerun the exact Node command and review source anchors, especially optional verification and PR gating. A Mermaid renderer check remains unexecuted. For runtime acceptance, coordinator must resume the original native ReleaseGate producer/QA owners, restore the regular local Python `.venv`, execute the final native test SHA, inspect actual DTO definitions/versions and matching counterpart APIs under expanded read authority, then obtain independent reducer and terminal readback. This architecture candidate makes none of those acceptance decisions.
