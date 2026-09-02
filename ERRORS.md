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
