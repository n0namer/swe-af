package node

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Agent-Field/SWE-AF/go/internal/furrow"
)

const waveDisposableEngineeringDiff = `diff --git a/internal/price.go b/internal/price.go
index 1111111..2222222 100644
--- a/internal/price.go
+++ b/internal/price.go
@@ -14,7 +14,7 @@ func Price(cents int) (int, error) {
-	if cents < 0 {
+	if cents <= 0 {
 		return 0, ErrInvalidPrice
 	}
 	return cents, nil
}`

func TestWaveBMADReviewEntrypointsAcceptOnlyBoundedEnvelopes(t *testing.T) {
	t.Setenv("SWE_PRO_ENGINE", "")
	t.Setenv("SWE_BMAD_ENABLED", "1")
	t.Setenv(furrow.EnvEnabled, "")
	t.Setenv(furrow.EnvPublicAddr, "")

	n, err := BuildAgent("swe-planner", "8005", "Autonomous SWE planning pipeline")
	if err != nil {
		t.Fatalf("BuildAgent: %v", err)
	}
	n.RegisterPlanner()

	registered := map[string]bool{}
	for _, name := range n.RegisteredNames() {
		registered[name] = true
	}
	for _, name := range []string{bmadReviewAdversarialGeneral, bmadReviewEdgeCaseHunter} {
		if !registered[name] {
			t.Fatalf("opt-in BMAD reasoner %q was not registered", name)
		}
	}

	if _, err := n.App.Execute(context.Background(), bmadReviewAdversarialGeneral, map[string]any{"content": ""}); err == nil {
		t.Fatal("adversarial reviewer accepted an empty envelope")
	}
	if _, err := n.App.Execute(context.Background(), bmadReviewEdgeCaseHunter, map[string]any{"content": 42}); err == nil {
		t.Fatal("edge-case reviewer accepted a malformed envelope")
	}

	out, err := n.App.Execute(context.Background(), bmadReviewEdgeCaseHunter, map[string]any{"content": ""})
	if err != nil {
		t.Fatalf("edge-case reviewer recovery: %v", err)
	}
	result, ok := out.(*bmadRunResult)
	if !ok {
		t.Fatalf("edge-case reviewer result type=%T", out)
	}
	if result.MethodID != edgeCaseHunterMethod.ID || result.SourceCommit != bmadSourceCommit {
		t.Fatalf("unexpected pinned response: %#v", result)
	}
	if strings.Join(result.CompletedSteps, ",") != "receive-content" {
		t.Fatalf("canonical recovery steps=%v", result.CompletedSteps)
	}
	if err := validateBMADRunOutput(edgeCaseHunterMethod, result); err != nil {
		t.Fatalf("canonical recovery output: %v", err)
	}
}

func TestWaveBMADReviewContractsKeepPinnedOrderAndRepairPayload(t *testing.T) {
	cases := []struct {
		name       string
		method     bmadMethod
		wantSteps  string
		finalOutput string
	}{
		{
			name:       "adversarial",
			method:     adversarialGeneralMethod,
			wantSteps:  "receive-content,adversarial-analysis,present-findings",
			finalOutput: "- zero is now rejected despite being a valid price\n- callers cannot distinguish zero from a negative price\n- boundary change lacks migration guidance\n- test coverage for zero is required\n- downstream contracts may expect zero\n- error documentation is absent\n- validation semantics changed silently\n- negative-price behavior is conflated\n- integer minimum remains unexamined\n- callers need a repair decision",
		},
		{
			name:       "edge-case",
			method:     edgeCaseHunterMethod,
			wantSteps:  "receive-content,exhaustive-path-analysis,validate-completeness,deletion-check,present-findings",
			finalOutput: `[{"location":"internal/price.go:17","trigger_condition":"cents equals zero","guard_snippet":"if cents <= 0","potential_consequence":"Valid free price rejected"}]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			seen := []string{}
			result, err := runBMADMethod(context.Background(), tc.method, map[string]any{
				"content":       waveDisposableEngineeringDiff,
				"also_consider": "The caller routes final reviewer output into a human-approved repair decision.",
			}, func(_ context.Context, method bmadMethod, step bmadStep, state, _ map[string]any) (*bmadStepResult, error) {
				if method.ID != tc.method.ID || state["content"] != waveDisposableEngineeringDiff {
					t.Fatalf("step %q did not receive immutable disposable diff", step.ID)
				}
				seen = append(seen, step.ID)
				stepResult := &bmadStepResult{Status: "completed", Summary: step.ID}
				if step.ID == tc.method.Steps[len(tc.method.Steps)-1].ID {
					stepResult.Output = tc.finalOutput
				}
				return stepResult, nil
			})
			if err != nil {
				t.Fatalf("runBMADMethod: %v", err)
			}
			if got := strings.Join(seen, ","); got != tc.wantSteps {
				t.Fatalf("step order=%q, want %q", got, tc.wantSteps)
			}
			if result.Output != tc.finalOutput {
				t.Fatalf("repair payload=%q, want final review output", result.Output)
			}
			if err := validateBMADRunOutput(tc.method, result); err != nil {
				t.Fatalf("fail-closed final output validation: %v", err)
			}
		})
	}
}

func TestWaveBMADReviewEntrypointsFailClosedAndRouteOnlyCompleteRetry(t *testing.T) {
	t.Setenv("SWE_PRO_ENGINE", "")
	t.Setenv("SWE_BMAD_ENABLED", "1")
	t.Setenv("SWE_DEFAULT_RUNTIME", "codex")
	t.Setenv("SWE_DEFAULT_MODEL", "offline-test")
	t.Setenv(furrow.EnvEnabled, "")
	t.Setenv(furrow.EnvPublicAddr, "")

	n, err := BuildAgent("swe-planner", "8005", "Autonomous SWE planning pipeline")
	if err != nil {
		t.Fatalf("BuildAgent: %v", err)
	}
	n.RegisterPlanner()

	for _, name := range []string{bmadReviewAdversarialGeneral, bmadReviewEdgeCaseHunter} {
		t.Run(name, func(t *testing.T) {
			_, err := n.App.Execute(context.Background(), name, map[string]any{"content": waveDisposableEngineeringDiff})
			if err == nil || !strings.Contains(err.Error(), "does not support Codex") {
				t.Fatalf("real entrypoint error=%v, want fail-closed unsupported-runtime rejection", err)
			}
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	canceledSteps := []string{}
	_, err = runBMADMethod(ctx, edgeCaseHunterMethod, map[string]any{"content": waveDisposableEngineeringDiff}, func(_ context.Context, _ bmadMethod, step bmadStep, _, _ map[string]any) (*bmadStepResult, error) {
		canceledSteps = append(canceledSteps, step.ID)
		cancel()
		return &bmadStepResult{Status: "completed", Summary: "interrupted before repair routing"}, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled run error=%v, want context.Canceled", err)
	}
	if got := strings.Join(canceledSteps, ","); got != "receive-content" {
		t.Fatalf("canceled run steps=%q, want only receive-content", got)
	}

	retrySteps := []string{}
	retry, err := runBMADMethod(context.Background(), edgeCaseHunterMethod, map[string]any{"content": waveDisposableEngineeringDiff}, func(_ context.Context, _ bmadMethod, step bmadStep, _, _ map[string]any) (*bmadStepResult, error) {
		retrySteps = append(retrySteps, step.ID)
		result := &bmadStepResult{Status: "completed", Summary: step.ID}
		if step.ID == "present-findings" {
			result.Output = `[{"location":"internal/price.go:17","trigger_condition":"cents equals zero","guard_snippet":"if cents == 0","potential_consequence":"Valid free price rejected"}]`
		}
		return result, nil
	})
	if err != nil {
		t.Fatalf("retry: %v", err)
	}
	if err := validateBMADRunOutput(edgeCaseHunterMethod, retry); err != nil {
		t.Fatalf("retry output: %v", err)
	}
	if got, want := strings.Join(retrySteps, ","), "receive-content,exhaustive-path-analysis,validate-completeness,deletion-check,present-findings"; got != want {
		t.Fatalf("retry steps=%q, want %q", got, want)
	}

	repairPayloads := []string{}
	if retry.Status == "completed" {
		repairPayloads = append(repairPayloads, retry.Output)
	}
	if len(repairPayloads) != 1 || repairPayloads[0] != retry.Output {
		t.Fatalf("repair routing=%q, want exactly one complete retry payload", repairPayloads)
	}
}
