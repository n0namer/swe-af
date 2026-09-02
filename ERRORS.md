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
- For the live runtime, place the wrapper at `/usr/local/bin/opencode` (or explicitly set `S[E_OPENCODE_BIN> to a persistent path).
- After the live helper was made executable in `/usr/local/bin`, the wrapper was actually invoked (`argc=8`) and FCM /OpenCode returned real reasoning and tool calls (`ls /tmp`, then `mkdir -p /tmp/fcm-canary-task`).

Prevention:
- Do not rely on an installer/wrapper directory being in herited PATH after `af run`.
 - Prefer an explicit `SWE_OPENCODE_BIN` pointing to a persistent binary path or ensure the binary is in the `af run` PATH.
- In acceptance canaries, distinguish quick`<1s` binary-resolution failures from provider inference timeouts.

Verification evidence:
- Before fix: OpenCode capture files were absent and `run_coder` failed in ~0.5s.
- After binary-resolution fix: wrapper capture showed `argc=8`, and FoM returned real model reasoning and tool_use events.
