---
title: 'Restore AgentField to SWE-AF callback delivery'
type: 'bugfix'
created: '2026-10-18'
status: 'draft'
review_loop_iteration: 0
context:
  - '{project-root}/AGENTS.md'
---
<frozen-after-approval reason="human-owned intent — do not modify unless human renegotiates">
## Intent
**Problem:** AgentField receives SWE Fast and Full heartbeats, but cannot deliver /reasoners/build requests. Real Golden Fast executions return agent_unreachable before calling an LLM.
**Approach:** Resolve the control-plane network boundary, configure a reachable callback address with bounded proxy exclusions, prove transport from the actual control plane, then execute Golden Fast via live FCM.

## Boundaries & Constraints
**Always:** Preserve existing dirty worktrees and production services. Work in isolated fixtures. Record real execution IDs, source SHA, tests, exact network path and provider receipts. Control-plane-side probes are mandatory.
**Ask First:** Global firewall/proxy changes, production restart, opening public endpoints or paid provider changes.
**Never:** Present direct Codex output as SWE-AF E2E, claim mock LLM as Golden, or weaken verification.

## I/O & Edge-Case Matrix
| Scenario | Input / State | Expected Output / Behavior | Error Handling |
|---|---|---|---|
| Fast | SWE Fast active, callback reachable | Real POST /reasoners/build delivered | Exact execution receipt |
| Full | SWE Full active, callback reachable | Real POST delivered | Exact execution receipt |
| Loopback | Callback 127.0.0.1 from isolated runtime | Fail preflight before dispatch | No silent timeout |
| Proxy | HTTP_PROXY applies to local call | Narrow bypass for trusted local endpoints | No global proxy disable |
| Unreachable | Bridge/firewall blocks callback | Bounded network failure | Correct BLOCKED_RUNTIME |
</frozen-after-approval>

## Code Map
- `swe_af/fast/app.py` — Fast AgentField entrypoint.
- `swe_af/app.py` — Full AgentField entrypoint.
- `sdk/python/agentfield/agent.py` (AgentField dependency) — callback URL resolver.
- `control-plane/internal/handlers/nodes_register.go` (AgentField dependency) — registration discovery.
- `_bmad-output/implementation-artifacts/` — current isolated candidate spec.

## Tasks & Acceptance
**Execution:**
- [ ] Identify actual control-plane network namespace, bridge and its proxy variables; prove candidate callback GET /health from this runtime.
- [ ] Implement minimum restricted callback transport/config, preserving existing live processes and global security.
- [ ] Add regression for reachable callback, loopback collision and proxy interception.
- [ ] Run live Golden Fast on disposable repo through FCM, with tests, Git SHA, execution ID and honest provider telemetry.
- [ ] Only after Fast GREEN run Full; do adversarial and edge-case review; write evidence to BMad and CrewSync.

**Acceptance Criteria:**
- Given connected SWE nodes, when the control plane invokes /reasoners/build, then Fast and Full receive real POST without agent_unreachable.
- Given a disposable fixture, when Fast runs via FCM, then real code is committed and independent tests pass.
- Given an unreachable endpoint, when preflight runs, then it fails closed with exact transport cause.
- Given a terminal execution, when verified, then receipts match code, test and Git artifacts.

## Verification
**Commands:**
- `curl http://127.0.0.1:8004/health` — Windows-only Fast health (not callback proof).
- `curl http://127.0.0.1:8080/api/v1/nodes/swe-fast/status` — registration only (not callback proof).
- Control-plane-namespace HTTP GET /health and real POST /reasoners/build — required proof.
- Independent fixture tests and `git rev-parse HEAD` — Golden evidence.

## Evidence (2026-10-18)
- FCM model list reachable at 127.0.0.1:19280; OpenCodex reachable at 127.0.0.1:10110.
- Control plane /health and SWE Fast /health: HTTP 200 on Windows.
- Executions exec_20261018_152445_5kw5h7yz, exec_20261018_152651_ww56buhu, exec_20261018_153246_5itn7t1p failed with agent_unreachable.
- Docker WSL with proxy interception returned 503 for host.docker.internal; direct access to Windows 172.20.20.187:8004 was refused; other gateway probes not accepted. No verified control-plane-to-Fast route yet.
- No live Golden GREEN.

## Evidence update — runtime topology
- Windows `Get-NetTCPConnection` confirms SWE Fast 0.0.0.0:8004 (PID 159464), Full 0.0.0.0:8003 (PID 152068).
- Windows netstat shows established outbound SWE agent connections to 127.0.0.1:8080; inbound control-plane callback still unproven.
- Docker Desktop WSL has running dockerd/containerd and container processes; Windows `docker ps` does not return promptly. Do not classify Docker itself as stopped.
- docker-desktop WSL `wget` to 172.18.0.2:8080/health did not return a valid HTTP reply; this address is not yet verified as the control-plane host.
- Current AgentField `swe-fast` node status: active/ready with heartbeats. This is registration proof, not callback delivery.
- Acceptance remains RED: no control-plane-namespace POST /reasoners/build and no Golden Fast/Full.

## Evidence update — final callback root-cause check
- Actual agentfield-server runs inside Docker network 172.26.0.6/16, default gateway 172.26.0.1 (verified with nsenter on PID 4248).
- Control-plane-network GET http://172.26.0.1:8004/health: connection refused.
- Windows has enabled inbound block rule codex_sandbox_offline_block_inbound, Action=Block, Profile=Any, Protocol=Any, LocalPort=Any, with remote address ranges covering non-loopback IPs.
- This is an intentional isolation control. Do not disable, weaken, or tunnel around it. Reassess supported topology with the responsible operator.
- Alternative Ubuntu WSL has Python 3.12.3 and source-tree visibility, but no installed agentfield Python package; no live Linux SWE-AF callback or FCM Golden proof.
- Status: BLOCKED_SECURITY_BOUNDARY (callback) and GOLDEN_UNACCEPTED; do not claim GREEN.
- The existing three failed AgentField execution IDs remain definitive negative evidence.

## MVP Fast Git preflight update — 2026-10-19
- Native AgentField callback to SWE Fast is working: execution exec_20261019_151810_u0xllzqw reached run_git_init, but returned result.success=false and timed_out=true.
- Replaced model-driven Fast Git-init with existing swe_af.issue.git_ops read-only preflight for ready Git repositories.
- Removed the ungrounded 60-second Git-init timeout. Unit regression verifies no run_git_init LLM call for an initialized repository.
- tests/fast/test_app.py: 27 passed, exit 0.
- Live execution exec_20261019_165548_pygl1ny6 submitted through native AgentField; last observed status running, Golden not accepted.
- Current gate: LIVE_CODING_RED. Do not mark CrewSync DONE or claim a successful full Golden until independent code/test/commit evidence exists.

## Golden Fast S evidence — 2026-10-19 17:08 local
- Native AgentField execution exec_20261019_165548_pygl1ny6 invoked swe-fast.build with live Codex gpt-6.1-sol.
- Git preflight used existing repository; planner completed one task and coder produced calculator.py and test_calculator.py.
- Coder returned complete=true, tests_passed=true, with commit e4398a075b1a2358e49dd10fd18884e274fd7867.
- Independent python -m unittest discover -q: 5 tests PASS, exit 0.
- Last checked root execution remains running at fast_verify; terminal business success and complete Golden acceptance still unconfirmed.
- Keep LIVE_CODING_RED until AgentField returns terminal result.success=true and corresponding verification receipts.

## 2026-10-19 quality gate: Fast false-positive / Full RED
- Golden Fast exec_20261019_165548_pygl1ny6 terminated AgentField succeeded, result.success=true, but execution_result.completed_count=0, failed_count=1 and task outcome timeout. Its independent calculator unittest suite 5/5 PASS and Git commit e4398a0 do not prove successful autonomous completion.
- Fixed swe_af/fast/app.py acceptance: verifier passed is necessary but insufficient; all planned task_results must be completed, with no failed_count or timed_out.
- Added regression test for verifier PASS plus task timeout; local tests/fast/test_app.py 28/28 PASS; related 45/45 tests PASS.
- Golden Full exec_20261019_170907_b572kk1l failed with execution timeout; initial phase remained planning + model-driven git init.
- Reused existing swe_af.issue.git_ops deterministic git preflight for initialized repos in Full when deterministic_git=true; keep LLM fallback for fresh/uninitialized repos and explicit deterministic_git=false.
- New Full exec_20261019_174256_v7bjqcip is still running at latest readback, not accepted.
- FCM /v1/models HTTP 200, but /v1/chat/completions smoke timed out after 30s; inference readiness remains RED.
- Still OPEN: Full live model plan/code/test/finalize terminal GREEN; Fast genuine task success receipt; independent multi-file repair scenario; truthful terminal outcomes.
- Never mark MVP DONE based solely on top-level transport status succeeded.

## 2026-10-19 adversarial QA readback
- Fast exec_20261019_165548_pygl1ny6: transport status succeeded, result.success=true but completed_count=0, failed_count=1, timed-out task. Corrected acceptance to reject any failed or timed-out tasks even if verifier passes.
- Regression test demonstrates failed task plus positive verifier must return result.success=false. 28/28 Fast tests PASS.
- Full deterministic Git preflight reuses existing swe_af.issue.git_ops for ready committed repositories when deterministic_git=true; LLM fallback for uninitialized repositories.
- Full E2E exec_20261019_170907_b572kk1l failed. Replacement exec_20261019_174256_v7bjqcip entered planning, still not terminal on last query.
- FCM model enumeration 200 but real completions 30s timeout. Do not accept inference readiness.
- Additional selected suite: 45/45 PASS. Broad pytest -x: 221 PASS and 1 check failed because test_ac_14_no_existing_swe_af_files_modified introspects uncommitted git diff and Full source is intentionally changed. Re-run after commit.
- OPEN: independent live Full acceptance and reliable Fast terminal success; do not declare MVP complete.
