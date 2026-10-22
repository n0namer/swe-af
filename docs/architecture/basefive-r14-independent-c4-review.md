# Independent R14 C4 source-bound QA

Date: 2026-10-15. Disposition: **CANDIDATE / UNACCEPTED**, pending native test evidence below. This is an independent source-contract review, not product acceptance.

## Binding, gate and measurable DoD

- FACT: Target is `D:/Users/NIKITA/Documents/DEV/universal-solver-site/swe-af`; BMAD-MNNZ supplies workflow instructions only.
- FACT: Task `swe-af:v2-r14b-independent-c4-qa-20261014-r1`, attempt `attempt-5970ef82-dcac-4ef1-a3f9-a09258761703`, session `manual-20261015-basefive-goal20-swe-af-v2-r14b--604f2591bbcccc34`; producer thread is null under the trusted binding.
- FACT: Supplied Git receipt is `064a7226a500365005b5b662044fc9014a328722`; manifest SHA is `8a2abf4b31d4111628fb8aa718da720fd6c95ce701e3a3920a5948a97739268b`. These are coordinator receipts, not an independent clean-worktree inspection.
- FACT: BMad source of truth remains `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md`. Stage: post-producer architecture QA. Current gate: source-bound C4 edges/DTO evidence plus native negative tests; native ReleaseGate runtime acceptance remains a separate blocked gate.
- FACT: Both exact QA write targets returned ENOENT before implementation. The existing producer suite already checks signatures, diagram syntax, SHA drift and DTO call sites; it is reused in-process rather than copied or replaced.
- DoD: independently bind all seven permitted input SHA values; parse complete supplied index snapshots for Go coverage; verify all 20 C4 edges against the two source files and the three diagrams; corroborate DTO declarations without asserting imported schema compatibility; execute at least two negative drift cases in native Node with a real exit; preserve limitations and publish exactly one typed UNACCEPTED result.
- Constraint: only this review and `basefive-r14-independent-c4.test.mjs` may change. No product, producer, generated-index, acceptance-source, dependency, task, lease or terminal-state edits.

## Source identities

The independent suite pins these values from direct scoped read receipts, separately from the producer contract, and hashes the actual bytes before comparing them. LF normalization is used only for text extraction, never for hashing.

| Source | SHA-256 |
|---|---|
| `AGENTS.md` | `4a8876425a07de89597f3b56c6ba918b8dd89ddbce1abb72b4f9db9ec1782ae9` |
| `docs/acceptance/c4-releasegate-reducer-frontier-20261013.md` | `feeb1fe6a88e53770831de890fe3674e6425dd40327deb913fbf73f3a01f1f17` |
| `.archsteer/model.json` | `ed964b0a700a566bd9e42caa55551baee873706aa500fbd57c446a28b7971260` |
| `.archsteer/architecture.md` | `7d7ee7046812a9c3c5da1d3b7484015432c4cb3b3cb106a6d34e8bffeb65a3d2` |
| `swe_af/issue/build.py` | `440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd` |
| `go/internal/orch/build.go` | `05c39a2e3b41869b01a1482543c9875630a0777fbea24ea706c05727fe640cbc` |
| `graphify-out/graph.json` | `e18cd58bd72e8975c07f305a988192227c529a39ab35752388ad0db706f5c277` |

FACT: Producer contract SHA from direct read is `cc0eeb34cbc3031785787caf1e72ce8393551533376d46403b71a393c3330029`; producer test SHA is `6b584644ced0edc8dfcd7cad01e526a79ddc0828fc52d12ca6ee7a8d22bf65ba`. The final native receipt prints exact hashes of all seven inputs, six producer artifacts and both QA outputs, without attempting a self-referential review hash in this document.

## Cross-checks and interpretation

- FACT: Python `implement_issue` delegates to `_implement_issue_impl` (`build.py:450`); issue validation occurs in the implementation (`build.py:196`). It does not call the Go Build entrypoint. Its input is a scoped issue, not a whole-build goal.
- FACT: Go `Build` and `Handlers` (`build.go:28`, `build.go:51`) expose the separate whole-build API. `buildInput` (`build.go:36`) contains ten fields with explicit JSON tags; ignored marshal/unmarshal errors do not establish strict boundary validation.
- FACT: Python constructs 16 `IssueBuildResult` fields and serializes via `.model_dump()` (`build.py:393`). Go constructs seven `schemas.BuildResult` fields and conditionally appends `workspace_handle` (`build.go:495`). Go constructor identifiers are not wire field names.
- FACT: The C4 views have 5/6/10 nodes and 4/8/8 edges. The independent checker freezes endpoint/source ownership and source witnesses, corroborates actual diagram edges, and checks conditional verification, PR gating, cancellation propagation, cleanup and failed-result transport. Reused producer checks additionally validate the restricted Mermaid grammar and exact node/edge labels.
- FACT: Python verification is optional (`build.py:308`); its PR gate does not test verifier success (`build.py:321`). Source success requires coding completion, commits, and verification success when present (`build.py:370`). These source observations are not interruption/restart runtime proof.
- FACT: Go success is derived from `verification["passed"]` (`build.go:448`); `agent.ReasonerFailed` can carry `buildResultMap` (`build.go:522`). Static Node tests do not exercise SDK transport.
- FACT: The supplied ArchSteer snapshot has zero `.go` component paths; Graphify has 200 distinct `.go` source paths in 36 parent directories, including all five internal directories imported by Build. These are index directory observations, not compiled package or repository-completeness evidence.
- FACT: Graphify records 10 nodes for the exact Python path and 14 for the Go path. ArchSteer lists both checkpoint helpers, but Graphify lacks their exact function labels at the Python path. Complete JSON parsing by Node provides the bounded inventory; no graph-navigation MCP is exposed, so no graph freshness/coverage-tool claim is made.

## Adversarial findings and residual limitations

- The Python and Go entries are different API levels; shared purpose cannot establish parity. The candidate explicitly labels this distinction, and the independent edge topology prevents a fabricated Python-to-Go call.
- Imported DTO definitions are outside scope. Versions, serializer aliases, required-field semantics and cross-runtime compatibility remain UNVERIFIED; no version may be borrowed from ArchSteer's `schema_version`.
- The producer's anchor-existence test alone cannot prove that an edge is semantically supported. The independent suite adds fixed reviewed endpoint/source/witness mappings and rejects substitution of an unrelated but existing anchor.
- Matching a diagram to a contract alone permits coordinated drift in both artifacts. The independent suite rejects a jointly changed diagram/contract edge against the source-reviewed topology.
- A hash-only drift check would not demonstrate DTO extraction. An in-memory Go source-tag mutation is rejected by the DTO checker without consulting hashes.
- Finding a verifier call does not establish its guard. An in-memory Python guard mutation leaves that call intact and is rejected by the independent control-flow witness.
- Zero ArchSteer Go paths is a snapshot blind spot, not proof that Go is absent from the project. The independent inventory preserves the gap and rejects fabricated Go coverage.
- Directory presence in Graphify is not Go package compilation or execution. No `go test`, build tags, imported implementation or complete package coverage was inspected.
- Index file fingerprints and matching short commit metadata do not establish source freshness. GraphSourceFreshness, repositoryCompleteness and packageRuntimeCoverage remain UNVERIFIED.
- Mermaid parsing is restricted grammar validation, not renderer execution or proof of deployment topology. No actual Python/Go process boundary, rendering or live call traffic is inferred.
- The reducer records `PYTHON_VENV_UNAVAILABLE`, `exit_code: null`, for a prior native Python test identity, not execution of the final Python artifact. This QA does not rerun, repair or promote that gate.
- Exact producer CrewSync publication/acceptance receipts and authoritative terminal readback are not in the read scope. Passing Node tests and VERIFIED publication transport must not be treated as ACCEPTED_GREEN.

## Native validation receipt

Command: `node --test docs/architecture/basefive-r14-independent-c4.test.mjs`, through `go_wave.run_project_tests`.

Status: NOT RUN YET. No runtime success is claimed before the actual runner receipt. The suite imports the existing producer test module through native ESM (no shell, child process or package script), then adds independent positive, negative and artifact-evidence tests.

Required negative cases: Go input JSON tag rename and fabricated Python-to-Go edge. Additional boundaries: unrelated existing source witness, fabricated ArchSteer Go count, fabricated DTO version, mutated Go source tag, jointly mutated diagram/contract edge and mutated Python verifier guard. Each checker validates the pristine input before expecting mutation rejection; all mutations remain in memory.

## Self-review, methodology and handoff

FACT: BMad Help selected this source-contract review gate, not a new planning workflow. Adversarial review and verification-before-completion methods were loaded. Config resolution, Git status/history, additional instruction discovery outside the read scope, graph-navigation tools, shell commands, interactive workflows and delegated review are unavailable or forbidden under the immutable scope. No complete BMad workflow or clean-Git claim is made.

Self-review: only two new authorized QA files; no new dependency; fixed allowlisted reads; no candidate-controlled filesystem lookup; native Node assertions and SHA-256; no acceptance or task-state mutation; both QA outputs must be below 500 lines. Existing producer tests run unchanged inside the independent test entrypoint. The remaining runtime/DTO/deployment limitations are intentional evidence boundaries, not claims repaired by changing product code.

Useful-work estimate: investigation 5 min + independent tests 6 min + meaningful validation/repair allowance 4 min + review/evidence 4 min = 19 min, external coordination overhead 0–1 min. This is a sizing estimate, not measured duration or a time floor.

Next move: independent reviewer/reducer may consume the exact final native Node receipt for this source-bound slice only. Keep product UNACCEPTED; coordinator must use the existing native ReleaseGate owners for Python environment recovery/runtime evidence, expanded DTO authority and authoritative terminal readback. No duplicate producer or scheduling action is requested.
