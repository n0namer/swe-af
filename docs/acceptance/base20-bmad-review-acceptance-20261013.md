# Base20 BMAD Reviewer Acceptance Evidence — 2026-10-13

## Scope

- Task: `swe-af:v2-base20-bmad-review-acceptance-20261013-r1`
- Target: the opt-in `bmad_review_adversarial_general` and `bmad_review_edge_case_hunter` AgentField reasoners in `go/internal/node/register.go`.
- Source of truth: `PLAN.md`, which records the reviewer-runtime fix as done and requires fail-closed structured reviewer output, pinned source identity, ordered steps, read-only behavior, and AgentField tracing.

## Contract Exercised

`go/internal/node/bmad_review_goal_wave_test.go` codifies the following deterministic coverage against the disposable `internal/price.go` diff where a boundary change rejects zero. Execution remains blocked until the authorized Go image is restored:

1. Both opt-in names register on the real `Node.RegisterPlanner()` AgentField surface when `SWE_BMAD_ENABLED=1`.
2. The adversarial entrypoint refuses empty input; the edge-case entrypoint refuses malformed input and returns the canonical structured empty-input finding.
3. The native executor preserves the pinned method order for both reviewer methods and carries the final reviewer output as the candidate repair payload.
4. The real registered entrypoints receive the disposable diff but reject `SWE_DEFAULT_RUNTIME=codex` before a provider call because the runtime cannot enforce the required read-only tool allowlist.
5. A canceled edge-case run returns `context.Canceled` after only `receive-content`; no partial result is routed to repair. A fresh retry runs the complete five-step method, validates its structured output, and produces exactly one candidate repair payload.

## Acceptance Boundary

- These are deterministic offline contract tests plus real AgentField entrypoint dispatch. They are not a paid-provider or live model end-to-end acceptance run.
- Reviewer results remain candidate review data. Any repair decision is human or independently authorized; no reviewer result self-accepts product work.
- The review path is read-only by contract. Existing package coverage in `node_test.go` separately asserts the OpenCode read-only policy (`Read` tool only, plan permission mode, deny-by-default policy) and rejects a Codex runtime that cannot enforce it.

## Validation Receipt

Attempted command through the only authorized scoped runner:

```text
go test ./internal/node -count=1
```

Actual runner result on 2026-10-13:

```text
status=BLOCKED
code=GO_DOCKER_IMAGE_UNAVAILABLE
engine=docker:go-test
required_image=golang:1.24.0-alpine
source_identity_sha256=f0f227f64310e94ba22045366bc663caa8208cdebbbcd64a5d72eb60f1b62d28
```

The image failure occurred before Go compilation or test execution. It is infrastructure evidence, not a passing receipt. Restore the authorized runner image, then rerun `go test ./internal/node -count=1`; independent QA/reducer must review the result before product acceptance.
