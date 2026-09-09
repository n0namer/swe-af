# ERRORS.md

## 2026-09-02 — OpenCode wrapper exists but is not executed by `swe-planner`

Status: VERIFIED runtime lesson.

Symptom:
- `swe-planner.run_coder` fails in ~ 0.5s with a harness/schema error.
- The OpenCode instrumentation files are not created, even though `/afhome/bin/opencode` exists and is executable.

Root cause:
- `af run` starts the package process with a sanitized PATH that did not include `/afhome/bin`.
  Observed on PID `8152`: `/opt/af/bin:/usr/local/go/bin:/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin`.
- `SWE_OPENCODE_BIN` was not set, so the AgentField harness fell back to resolving `opencode` from PATH.

Fix:
- For the live runtime, place the wrapper at `/usr/local/bin/opencode` (or explicitly set `SWE_OPENCODE_BIN` to a persistent path).
- After the live helper was made executable in `/usr/local/bin`, the wrapper was actually invoked (`argc=8`) and FCM/OpenCode returned real reasoning and tool calls (`ls /tmp`, then `mkdir -p /tmp/fcm-canary-task`).

Prevention:
- Do not rely on an installer/wrapper directory being in inherited PATH after `af run`.
- Prefer an explicit `SWE_OPENCODE_BIN` pointing to a persistent binary path or ensure the binary is in the `af run` PATH.
- In acceptance canaries, distinguish quick `<1s` binary-resolution failures from provider inference timeouts.

Verification evidence:
- Before fix: OpenCode capture files were absent and `run_coder` failed in ~0.5s.
- After binary-resolution fix: wrapper capture showed `argc=8`, and FCM returned real model reasoning and tool_use events.

## 2026-09-02 — Unit tests passed but external CLI acceptance still failed

Status: VERIFIED runtime lesson.

Symptom:
- Calculator coder iteration returned `complete=true` / `tests_passed=true` after 13 unit tests passed.
- Independent operator acceptance `python3 /tmp/fcm-calculator-live/cli.py divide 1 0` still produced an uncaught `ValueError` traceback, violating the explicit controlled-error criterion.
- The first regression test added during recovery passed only when invoked from the project directory; the same suite from an external cwd failed because the subprocess used relative `cli.py` and exited `2`.

Root cause:
- The initial test suite covered the library `divide()` contract but not the user-facing CLI failure path.
- The first recovery test encoded an implicit cwd assumption instead of locating the CLI independently of the caller working directory.

Fix:
- Recovery iteration 2 caught `ValueError` in the CLI, emitted `Error: Cannot divide by zero`, returned controlled exit code `1`, and removed the traceback while preserving the library exception contract.
- Recovery iteration 3 made the subprocess regression test cwd-independent; operator `unittest discover` from outside the project then passed all 13 tests.

Prevention:
- Do not accept `tests_passed=true` as task acceptance when the DoD contains externally observable CLI/API behavior; run an independent behavioral oracle for those criteria.
- Regression tests for CLI entrypoints must not depend on the caller cwd; resolve the entrypoint relative to the test file or otherwise use a stable path.
- When an operator oracle contradicts the agent's completion claim, treat the oracle as fresh acceptance evidence and start a bounded recovery iteration rather than marking DONE.

Verification evidence:
- Initial operator CLI oracle: exit `1` with uncaught traceback despite agent `complete=true`.
- Product recovery oracle: exit `1`, stdout `Error: Cannot divide by zero`, no traceback.
- External-cwd suite after test recovery: 13 tests, `OK`.
- Recovery execution `exec_20260902_120338_5ac3xq2q`: `succeeded`, `complete=true`, only `test_calculator.py` changed.

## 2026-09-08 — Coder structured-output failure after real work must recover in-place

Status: VERIFIED runtime lesson.

Symptom:
- Full `implement_issue` execution `exec_20260908_165317_24ztwr07` ended `failed_unrecoverable` even though the coder had already edited and committed the issue worktree.
- The resulting coder commit `b9c2dde90623a55af12f026290d2ad9162608697` was independently invalid: exact-source Python compile failed with `SyntaxError: 'return' outside function`.
- The coder role then failed during structured-output/schema completion, so the outer coding loop exited before reviewer/repair could recover the partial work.

Root cause:
- `RunCodingLoop` treated any non-fatal coder call error as immediately unrecoverable, without distinguishing an untouched worktree from a worktree that the failed coder had already modified or committed.
- Therefore a transport/schema/output failure could terminate the orchestration after real product changes, bypassing the existing multi-iteration repair budget.

Fix:
- Fingerprint git HEAD + porcelain worktree state before each coder call.
- If a non-fatal coder error leaves changed git state and coding iterations remain, preserve the same worktree, record `coder_retry` feedback requiring exact syntax/build/tests, and continue to the next coder iteration.
- If the coder error leaves the worktree unchanged, keep the previous fail-fast behavior.

Prevention:
- Do not equate `coder call returned error` with `no work was produced`.
- Before classifying a coder failure as unrecoverable, compare pre/post git state and preserve partial work only when that state changed.
- The recovery iteration must distrust prior test claims and re-run validation on the exact worktree before completing.

Verification evidence:
- New regression `TestCoderExceptionAfterWorktreeChangeRetries` passes: iteration 1 commits partial work and returns a structured-output error; iteration 2 receives repair feedback and the loop completes instead of returning `failed_unrecoverable`.
- Fresh real-repo execution `exec_20260908_180746_5l94yhja` subsequently completed successfully on EvalGuard #3 and produced final commit `c1442dd2878bbd09f43da5d6696e8326d7232a21`.

## 2026-09-08 — Exact-source acceptance must bind the repository-local Python environment

Status: VERIFIED runtime lesson.

Symptom:
- A test command invoked with a virtualenv executable from another EvalGuard clone reported passing tests even though the intended final worktree contained invalid code in an earlier run.
- Later verifier runs also failed with `python: command not found` despite a valid repository-local Python environment being present as `venv/bin` rather than `.venv/bin`.

Root cause:
- Editable Python installs and absolute virtualenv executables can silently bind imports to a different checkout than the worktree being certified.
- SWE role environment discovery recognized only `.venv/bin`, while real repositories may use either `.venv/bin` or `venv/bin`.

Fix:
- Exact acceptance commands bind `PYTHONPATH` to the intended worktree and put that worktree's repository-local virtualenv first in `PATH`.
- Coder, reviewer, and verifier OpenCode role environments now recognize both `.venv/bin` and `venv/bin`, preferring `.venv` when both exist.

Prevention:
- Never accept a test PASS unless tested-source identity is proven; editable-install provenance from another checkout is insufficient.
- For Python repo acceptance, verify the interpreter path and source import root together.
- Support both common repository-local virtualenv layouts instead of assuming one naming convention.

Verification evidence:
- On exact final EvalGuard commit `c1442dd2878bbd09f43da5d6696e8326d7232a21`, Python compile PASS, focused tests 12/12 PASS, full suite 71/71 PASS, and hidden delimiter oracle 3/3 PASS when `PATH` and `PYTHONPATH` are bound to that worktree.
- Standalone verifier `exec_20260908_191034_fqswvej5` used `/tmp/evalguard-issue-recovery/.worktrees/794aeb7d-evalguard-3/venv/bin/python` and completed terminally with `passed=true`, all 5 acceptance criteria PASS.

## 2026-09-09 — Reviewer structured output must not assemble verdict files in the product worktree

Status: VERIFIED runtime lesson.

Symptom:
- EvalGuard repeat-5 reached coder repair, reviewer approval, and verifier 5/5 PASS, but final delivery correctly failed closed because reviewer OpenCode left `?? .review_verdict.json` in the product worktree.
- The branch diff itself was otherwise scoped to exactly the three declared product files.

Root cause:
- Reviewer used incremental structured-output assembly, so its transport-level verdict JSON was materialized in the role cwd, which is the product worktree.
- Delivery hygiene therefore treated reviewer transport state as a product mutation.

Fix:
- Set reviewer-only `RoleOptions.SchemaMode="single"` while leaving coder incremental and verifier single-shot.
- Preserve reviewer OpenCode permission overlay (`task=deny`, `external_directory=deny`) and keep the final Git cleanliness guard unchanged.

Prevention:
- Structured role outputs are runtime artifacts, not product files; their ownership must be explicit at the role boundary.
- Never whitelist transport scratch by filename in the delivery guard.
- Prefer role-local transport configuration over weakening final Git cleanliness checks.

Verification evidence:
- Deterministic reviewer regression, targeted/full Go test/build/vet gates passed on the exact live source.
- Loaded planner generation SHA256 `6260583d91662f242d9471872a35385791cb64b350b952e95a2eedf38c3016ae` ran repeat-7 reviewer `exec_20260909_100454_xxwg4zkc` to `approved=true`, `blocking=false` with a clean worktree and no `.review_verdict.json`.

## 2026-09-09 — Final verifier must execute an independent boundary-negative variant

Status: VERIFIED runtime lesson.

Symptom:
- Repeat-9 completed `success=true`; reviewer approved and verifier reported 5/5 PASS.
- An independent post-run oracle on the exact worktree found that `CommentDocstringStripMutator.mutate("x = 1")` returned a mutation whose only change was an added trailing newline, violating the explicit no-formatting-only-mutation criterion.

Root cause:
- Final verification reused repository/coder happy-path evidence and did not independently vary a representation boundary relevant to no-op/normalization behavior.
- A clean input with a trailing newline passed while the equivalent clean input without the trailing newline failed, so ordinary tests were insufficient to certify the criterion.

Fix:
- Add a verifier `Boundary-negative discriminator`: for no-op, idempotency, normalization, preservation, parser/serializer, or boundary-sensitive criteria, execute at least one independent variant not copied from coder tests.
- For source/text transforms, explicitly vary trailing-newline presence/absence, whitespace/minimal input, and relevant delimiter/quote form; formatting-only changes fail the criterion conservatively.

Prevention:
- A verifier criterion cannot pass solely from agent-written tests when the criterion has a meaningful representation boundary.
- Keep an independent acceptance layer that varies the input shape, not just the expected business behavior.
- When an independent boundary discriminator contradicts prior PASS evidence, the discriminator wins and the criterion must fail.

Verification evidence:
- Prompt/golden regressions plus full Go test/build/vet passed and were loaded in planner SHA256 `4e07d56ad23f6df9d2759bb7937029be5a75263b47174dd949ffd07061451794`.
- Verifier canary `exec_20260909_173118_xncpgrv7` explicitly compared clean source with and without a trailing newline, reproduced the formatting-only mutation, changed criterion 4 to FAIL, and produced overall FAIL 4/5 with the correct owning fix in `python_mutator.py`.

## 2026-09-09 — OpenCode verifier must be noninteractive and repository-local

Status: VERIFIED runtime lesson.

Symptom:
- The first isolated hardened-verifier canary delegated to `@explore` and then stalled on `external_directory=ask` despite final verification being intended as an unattended, repository-local role.
- No semantic verdict could be trusted while the hidden subagent was waiting for a permission decision.

Root cause:
- OpenCode does not currently enforce the verifier's harness tool list strongly enough to prevent hidden task/subagent delegation.
- Without a role-level permission overlay, a noninteractive verifier could enter a permission-prompt path outside the intended repository boundary.

Fix:
- For OpenCode verifier runs, set `OPENCODE_CONFIG_CONTENT` permission overlay to `task=deny` and `external_directory=deny`.
- Keep verifier `SchemaMode="single"` and repository-local cwd/environment binding.

Prevention:
- Noninteractive acceptance roles must not depend on permission prompts or hidden subagents.
- Enforce repository locality at the runtime permission layer, not only in natural-language prompts.
- Treat a verifier that waits for interactive permission as inconclusive, not as a semantic application failure.

Verification evidence:
- Regression coverage asserts the exact OpenCode permission overlay and full Go test/build/vet gates pass.
- On the next loaded generation, verifier canary `exec_20260909_173118_xncpgrv7` stayed in the primary repo-local agent path, did not enter `@explore` / `external_directory=ask`, and proceeded to the independent boundary test that correctly failed criterion 4.
