# Luna High disposable Git restart oracle — 2026-10-13

## Scope and decision

**Candidate evidence only: UNACCEPTED.** This artifact adds an independent,
local, disposable-repository oracle for the existing issue-level Git lifecycle.
It does not change production code, invoke a live AgentField task, create a
remote, push, merge, or claim `ACCEPTED_GREEN`.

## Source-backed C4 boundary

The tested component boundary is:

```text
implement_issue(issue, repo_path, base_branch, artifacts_dir, config)
  -> _implement_issue_impl(...)
     -> IssueSpec / IssueBuildConfig DTO validation
     -> git_ops.ensure_issue_ready_repo / resolve_base / ensure_local_excludes
     -> git_ops.add_worktree(branch off base SHA)
     -> run_coding_loop(issue, DAGState, ...)
     -> git_ops.scrub_tracked_junk / commit_all
     -> git_ops.new_commits / changed_files / diff_stat
     -> optional _run_verification(...)
     -> git_ops.remove_worktree (branch remains deliverable)
     -> IssueBuildResult.model_dump()
```

Relevant source evidence:

- `swe_af/issue/build.py` — SHA-256
  `aa69eb29670c2eef760a3d4becf4651cdf0af5dc4b6a14622f5caf6eccf4dc00`.
  The public `implement_issue` contract returns a structured result whose
  `branch` is the deliverable; the implementation isolates the caller checkout
  in a worktree and performs verifier work only after a completed loop and a
  commit.
- `swe_af/issue/git_ops.py` — SHA-256
  `0f73d9fb3a613f1b68ccf4a8161a93c1c5e6dd6a4a392b2c2c8a6157b7c2403d`.
  `add_worktree` retries only transient lock failures; deterministic branch or
  base failures stop without repeating the mutation. `commit_all` is
  append-only and returns an empty SHA when there is nothing new to commit.
- DTO boundary: `IssueSpec` and `IssueBuildConfig` are validated before setup;
  `IssueBuildResult` exposes `success`, `base_sha`, `branch`, `commits`,
  `files_changed`, `verification`, and `error_message`.
- Native AgentField dependency: the public wrapper obtains `issue_router.app`
  and calls its `agent.call` adapter; the oracle replaces only the internal
  coding-loop seam and supplies a deterministic verifier callback. No live
  provider, native task, lease, or CrewSync state is exercised or inferred.

## Oracle coverage

`tests/issue/test_luna_high_issue_git_restart_20261013.py` — SHA-256
`e6afbdb35a2ab3147bd8dc42c6e4f485ba0984ad63824798b1d04274b59edb29`.

The test file uses real local Git repositories and covers:

1. **Positive:** edit → checkpoint commit → verifier → persistent branch
   artifact, with exact base SHA, one commit, changed-file assertion, no
   remote, and worktree cleanup.
2. **Dirty parent isolation:** uncommitted caller content remains unchanged and
   outside the issue branch.
3. **Path boundary:** a malicious issue name is slugified into a safe branch
   name without `..` traversal segments.
4. **Interruption/restart:** a coder exception salvages one checkpoint commit;
   a later clean run produces one independent final commit and does not duplicate
   accepted history.
5. **Fail closed:** a non-completed coding outcome never invokes verification
   and remains unsuccessful even when checkpoint work exists.
6. **Idempotence:** a second `commit_all` with no new changes returns no SHA and
   leaves exactly one accepted commit.

## Validation receipt

The authorized task requires a focused Python test execution. The GO Wave
runner available to this worker exposes scoped Node/Vitest/tsx execution only;
no Python/pytest runner is authorized in this attempt. Therefore the following
is intentionally **BLOCKED**, not promoted to PASS:

```text
command: .venv\Scripts\python.exe -m pytest -q tests/issue/test_luna_high_issue_git_restart_20261013.py
status: BLOCKED (runner capability not available in scoped interface)
engine: python:pytest
source_identity_sha256: d6cfd240e6abae0bf1500515b791ecfc6ce2491936991e9bfc7ceaee92d66670
exit_code: null
```

The required next validation is one real local pytest run against this exact
test SHA, followed by independent QA review of the source-backed C4 boundary.
No live end-to-end claim is made here.

## Acceptance checklist

- Real positive scenario is specified with exact source and test SHA: **PASS**.
- At least two fail-closed/restart regressions are specified: **PASS**.
- C4 API/DTO boundary and AgentField dependency are traced to source: **PASS**.
- Actual focused Python test exit: **BLOCKED** by authorized runner capability.
- Worker result and independent QA required for `ACCEPTED_GREEN`: **PENDING**.
