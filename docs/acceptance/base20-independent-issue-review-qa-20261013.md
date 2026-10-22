# Base20 Independent Issue Review QA — 2026-10-13

## Gate and scope

- **Task:** `swe-af:v2-base20-qa-issue-review-full-cycle-20261013-r1`
- **Source of truth:** `PLAN.md`, currently active. It requires an executable independent oracle, durable recovery evidence, and forbids inferring acceptance from a reviewer response, process exit, or green-looking branch.
- **Reviewed producer evidence:** `tests/issue/test_build_recovery_goal_wave.py`, `docs/acceptance/base20-issue-build-repair-20261013.md`, `go/internal/node/bmad_review_goal_wave_test.go`, and `docs/acceptance/base20-bmad-review-acceptance-20261013.md`.
- **Frozen mechanism:** `swe_af.issue.build._implement_issue_impl`, exercised only with disposable local Git repositories. No producer source, Go runtime source, branch, or live AgentField state is changed.

## Independent lifecycle review

The new pytest module is reviewer-owned evidence. It performs two independent, safe issue runs (`alpha_qa.py` and `beta_qa.py`) that each:

1. Freeze the issue result base SHA, issue branch SHA, current `main` SHA, verifier outcome, and executable artifact.
2. Use a deliberate first incorrect implementation (`answer() == 0`), followed by the bounded same-task repair that writes the requested value.
3. Require the independent Python oracle to import the branch artifact and prove the requested result.
4. Require the exact producer method sequence: coder, reviewer RED, coder repair, reviewer, verifier.
5. Treat only the frozen base/current source match plus a passing oracle and deliverable branch as `ACCEPT_SCOPED` in this limited fixture scope.

This is not a live AgentField/provider exercise. The Python tests replace only the model/AgentField boundary; no configured provider, network endpoint, or runtime deployment is represented.

## Fail-closed adversarial checks

The QA suite also proves that the reviewer gate issues `RED`, never an inferred acceptance, for:

- a reviewer-approved result whose independent verifier returns `passed: false`;
- an isolated issue branch that completes after `main` advances, making its frozen base stale;
- a simulated interruption after an uncommitted artifact is written, which must become recoverable RED evidence on the retained branch;
- a one-iteration blocked review where a partial artifact is checkpointed but no duplicate commit is allowed.

The stale-source check deliberately distinguishes a successful producer process from the independent reviewer decision: the producer can return `success: true` while the reviewer still returns `RED` because the source snapshot moved.

## BMad method and trust boundary

The published BMad tests establish method ordering and structured-output validation separately:

- Adversarial General: `receive-content → adversarial-analysis → present-findings`.
- Edge Case Hunter: `receive-content → exhaustive-path-analysis → validate-completeness → deletion-check → present-findings`.
- Both preserve the immutable submitted content in ordered local execution and return review output for a separate repair decision; neither method mutates a reviewed diff or confers product acceptance.
- The native entry-point coverage proves malformed envelopes are rejected and canonical recovery output remains pinned/validated. This QA bundle relies on that contract but does not misrepresent the mock Python boundary as a live AgentField run.

## Review result

**Status: BLOCKED pending scoped runner evidence.** The one allowed focused runner invocation for `tests/issue/test_base20_independent_issue_review_qa.py` returned `GO_WAVE_NODE_ERROR:test_file_not_allowed`. The available runner admits Node-family tests, not Python pytest files; no identical retry was attempted, and the affected Python sibling was not retried through the same unsupported mechanism. This worker has not converted a process result into a product acceptance decision. After fresh focused evidence is available, the independent reviewer may issue only `ACCEPT_SCOPED` or `RED` for this bounded fixture; live-provider/runtime acceptance remains outside this scope.

## Required verification

```text
python -m pytest -q tests/issue/test_base20_independent_issue_review_qa.py
python -m pytest -q tests/issue/test_build_recovery_goal_wave.py
```

The expected evidence is the complete pytest output and the exact test-source SHA from the verification environment. Cost/token telemetry is unavailable because these tests do not call a provider.

## Narrow repair request if RED

Repair only the producer behavior demonstrated by a failing focused assertion. Preserve the frozen-source rule, independent oracle authority, reviewed-diff trust boundary, and no-duplicate-effect expectation. Do not modify this reviewer artifact to reinterpret a failing verifier, moving source, budget exhaustion, or infrastructure denial as acceptance.
