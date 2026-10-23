# Epic 1: SWE-AF — independently accepted MVP

Source of product goal and invariants: repository `PLAN.md` (Global North Star) and `AGENTS.md`. This epic defines only the **current VP49 MVP delivery gates**. It does not replace those project-wide contracts.

Frozen starting evidence (2026-10-23): `integration/swe-af-vp49-green-20261023@52b3fddc844fc3b5d7e5f0943bf1b19aaf78201b` (clean on VP49). Existing Go DAG/orch/issue/node tests and targeted Python tests passed. A real `swe-fast.build` execution previously produced the expected file but no Git commit and returned `output_data.success=false`. Current launcher file points to the integration worktree; **loaded-process identity and new E2E are not yet verified**.

## Story 1.1: Exact runtime identity
Reuse Windows Scheduled Tasks, `ops/windows/SWE.ps1`, and existing AgentField nodes, not a new controller. Given a quiescent system, when inspecting the running node, prove loaded source SHA, environment, model route, reachable callback, single owner and clean target; switch/restart only the existing owner if needed. CrewSync key: `sweaf:mvp:01-runtime`.

## Story 1.2: Fast E2E Git receipt
Reuse `swe_af/fast/*`, `swe_af/issue/git_ops.commit_all`, native git init/finalize and the already integrated fail-closed `52b3fdd` change. Given one disposable Git repository, when a bounded `swe-fast.build` task finishes, the required file must be included in HEAD and independently validated; accept only `output_data.success=true`, matching SHA and no unintended file changes. On 429 or ambiguous POST, inspect execution/effects before retry. CrewSync key: `sweaf:mvp:02-fast-commit`.

## Story 1.3: Native implement_issue E2E
Given a safe existing repository, when `swe-planner.implement_issue` runs, prove one bounded code change, canonical tests, independent executable oracle, resulting commit and execution ID; preserve the original repo outside the intended worktree. CrewSync key: `sweaf:mvp:03-issue-e2e`.

## Story 1.4: Recovery invariants
Given a timeout/restart and foreign/corrupt checkpoint, when resuming, prove no duplicate external effects, fail-closed invalid state and durable progress using existing Go `dag/checkpoint` and `orch/resume`. Cherry-pick only actually needed and verified A55 rescue changes. CrewSync key: `sweaf:mvp:04-recovery`.

## Story 1.5: One two-issue Full Build
Given two dependent scoped changes, when running existing `swe-planner.build`, prove ordered execution, valid merge, passing independent tests, clean terminal result and exact Git SHA. Avoid expanding to a large multi-agent load test. CrewSync key: `sweaf:mvp:05-full-build`.

## Story 1.6: Repeated clean-runtime MVP acceptance
Given an exact-SHA clean materialization, when running a small task ladder and three bounded Full Builds, prove semantic PASS (not just HTTP or agent status), recovery safety, immutable Git artifacts and cost/latency/provenance receipts; synchronize the authoritative BMAD status and CrewSync after independent readback. CrewSync key: `sweaf:mvp:06-acceptance`.

No extra orchestration, framework installation, scheduler, blanket Furrow Windows port or expanded benchmarking until these six gates pass.
