# Windows A55: AgentField / SWE-AF no-console deployment
Date: 2026-10-20. Host: DESKTOP-A55K2JN. Scope: console windows only, not new orchestration.

## Source of truth
- AgentField repository: D:\Users\NIKITA\Documents\DEV\universal-solver-site\agentfield
- Shared SDK launcher: sdk/python/agentfield/harness/_cli.py
- Applied AgentField dev commit: 2d44bf89fffa03e372b7f3b44cee3a2984bd1700 (2 added lines)
- Isolated fix commit: 6f52089d5edaf4944eedd9b3c4102c1ed46f770f
- Reproducible patch: agentfield-windows-cli-2d44bf89.patch
- Runtime scripts: D:\Users\NIKITA\Documents\DEV\swe-af-local-ops
- Current Task Scheduler exports: scheduled-tasks-installed/
- Prechange Task Scheduler exports: scheduled-tasks-before/
- SDK prechange backup: D:\Temp\af-window-cli-predeploy-20261020-143616.py
- Runtime test snapshot: deployment-evidence.json (QA limitations remain recorded there)

## What was changed
1. Python AgentField harness run_cli uses subprocess.CREATE_NO_WINDOW on Windows; taskkill fallback suppresses its own console.
2. SWE-AF Native Factory and SWE-AF Native Full30 scheduled starts use PowerShell -WindowStyle Hidden.
3. Periodic AgentField monitor uses hidden PowerShell; Python recovery and diagnostics use pythonw.exe + scheduled-pythonw.py, logging stdout/stderr to the normal ops logs.
4. Recovery watch's Git subprocess uses CREATE_NO_WINDOW. All original triggers, principals and schedules remain with existing tasks.

## Recovery
Check AgentField /health on 18080 and SWE-AF /health on 8009 / 8010; check /api/v1/executions/active?limit=100 before any restart.
For a task configuration rollback, inspect scheduled-tasks-before/*.xml and restore only the specific task using штатный Task Scheduler; do not mass-import or restart active tasks. Do not terminate other users' cmd/python/codex processes.
For AgentField code rollback, reverse only commit 2d44bf89, after coordinating the owner of the dev branch and confirming no writers. All active long-lived Python processes may retain prior imported code until they are restarted.
Four legacy nodes 8003..8006 were not recycled or verified to load the changed SDK: they have separate owners and live process configurations. Independent QA, full live Codex reasoning E2E and post-reboot acceptance are NOT complete.

## Safety
Do not stage logs, tokens, private .env, monitor state, or unowned worktrees. This source-control snapshot creates no extra process manager and should not run any additional tasks.
