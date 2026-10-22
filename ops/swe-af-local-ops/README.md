# SWE-AF local multi-project operator — A55

## Purpose
Use **native** SWE-AF `implement_issue` via AgentField for independent tasks in multiple local repositories. The shell scripts are thin REST/observability wrappers, not another coding engine, DAG scheduler, or acceptance authority.

AgentField: `http://127.0.0.1:18080` (A55 localhost only).
SWE native factory: `swe-factory` (Windows Task Scheduler `SWE-AF-Native-Factory`).
Next-start code version: `de97d0e444cadfe15a33f25d9881a985b23edfab` in `D:\Users\NIKITA\Documents\DEV\swe-af-native-runtime-30-20261020` (native file-conflict serialization; independent QA 36 PASS). Existing running process and Full runs remain on their originally loaded version until the service restarts.

## Operations
From PowerShell **on A55**:

```powershell
$ops = 'D:\Users\NIKITA\Documents\DEV\swe-af-local-ops'
& "$ops\status.ps1"
& "$ops\status.ps1" -Json
& "$ops\submit-issue.ps1" -Project af-skill-intelligence -IssueFile "$ops\issues\skill-hot-agentfield.json"
# Only after checking no active duplicate, use -Submit.
& "$ops\submit-issue.ps1" -Project af-skill-intelligence -IssueFile "$ops\issues\skill-hot-agentfield.json" -Submit
```

Read-only Windows monitor: `SWE-AF-AgentField-Local-Monitor`, every **1 minute**.
Files: `logs/latest.json`, `logs/events.ndjson` (rotates), `logs/factory-node.log`.
Monitor never dispatches/accepts tasks or edits repositories.

## Current project bindings
| Project | Isolated repo | First gate |
| --- | --- | --- |
| CodeAF | `D:\Temp\codeaf-swe-full-mission-20261019` | Real Windows CLI/provider/edit/Go test/SQLite journal |
| PR-AF | `D:\Temp\swe-pr-af-mission-20261020` | G0: permanent DEV runtime implementation, source and callback identity; no unproven changes |
| AF Skill Intelligence | `D:\Temp\swe-skill-intelligence-mission-20261020` | WIP Go hot-path adapter with independent QA; `fdfa3d0` is not accepted |

**Always inspect active executions before submitting.** `submit-issue.ps1` checks duplicate issue ownership and declared-file conflicts within a project, and limits PR-AF to the G0 evidence task until runtime identity is verified. Each issue has its own Git worktree and source-owned files. Existing original checkouts remain dirty and read-only.

## Multi-project concurrency
Multiple calls to the *same* `swe-factory.implement_issue` can run concurrently if repository and declared file ownership do not conflict. Configured operator limit: **30 concurrent independent issue roots** (user override, 2026-10-20). For future native Full DAG launches, the intended width is **10 PR-AF + 10 AF Skill Intelligence + 10 CodeAF = 30 issue workers** across three projects, subject to DAG readiness. The root limit is not an AgentField server-side semaphore, and current Full runs retain their originally submitted 3+3+4 limits. These caps do not automatically enqueue or guarantee 30 available/working provider processes. Keep independent file scopes, existing-issue deduplication, observed provider throttling and independent QA gates; AgentField capacity may require later tuning. An internal SWE execution heartbeat is every 60 seconds; the AgentField node heartbeat is separate.

## Native Full architecture for 30-slot readiness

- Future Full DAG configuration: 10 concurrent issues per project (PR-AF, AF Skill Intelligence, CodeAF) = 30 potential issue workers. `project-profiles.json` holds both the operator root cap (30) and `full_dag_issue_limits`; `launch-full-20261020.py` verifies them without dispatching duplicates.
- Full `plan` produces `PlannedIssue(name, depends_on, provides, target_repo, files_to_create, files_to_modify, acceptance_criteria, testing_strategy, guidance.touches_interfaces)` and topological parallel levels. Declare exact file owners and dependent interface signatures *before* implementation. Public method/event/wire signatures and schema migrations need focused contract tests, with dependent implementation gated on accepted signatures.
- New pinned SWE-AF source `de97d0e` serializes declared same-level file overlaps using additional `depends_on` edges, preserving width for independent tasks. **It is not loaded into existing running Full jobs.** Unlisted file edits, generated files, directory-level side effects and semantic API conflicts are not provably conflict-free; QA must detect them.
- Built-in Full execution is **level-barriered** (`asyncio.gather`), not a 30-slot continuous refill queue. It does not guarantee 30 READY tasks. `max_concurrent_issues` is per DAG (not AgentField-global). AgentField currently recommends substantially fewer active executions on the 12-core A55; monitor resource use and 429.
- Each issue uses isolated Git worktrees; deterministic `--no-ff` Git merging and merger-agent fallback happen after issue execution. Do not count merges or `result.success` as CrewSync `SERVER_DONE` without independent acceptance.
- No automatic dispatch/refill has been installed. Minute Windows monitor is read-only. Re-launch Full only after reconciling existing active jobs and latest accepted candidate SHA; no restart of live Full just to raise limits.

## Add another project later
1. Inspect original Git, BMad, CrewSync and active owners; create an isolated `git worktree add -b mission/... D:\Temp\...` from an approved exact SHA. Preserve all original dirty/untracked work without assuming the branch contains it.
2. Add one `projects[]` object to `project-profiles.json` with `id`, `safe_repo_path`, `base_branch`, `source_of_truth`, `readiness` and `branch_prefix`.
3. Write a **single bounded** IssueSpec JSON under `issues\`, with exact `files_to_modify/create`, tests and acceptance criteria. Dry-run `submit-issue.ps1`, then `-Submit`. Recheck AgentField execution logs and Git.
4. A successful transport status is not business DONE: independently verify test evidence and accepted result; CrewSync `SERVER_DONE` requires its own authorized settlement/readback. Do not mark completion by a local Git commit alone.

## Safety/authority
Never `git clean` or reset the original checkouts; no automatic push, PR publication, deploy, external sends, credential printing or paid API calls. PR-AF SourceLoop requires current permanent DEV identity. Skill Intelligence owns no executable Registry promotion or new workflow scheduler. CodeAF must prove real native E2E. CrewSync tracks cross-project claims and terminal acceptance; AgentField executes tasks, not CrewSync.
