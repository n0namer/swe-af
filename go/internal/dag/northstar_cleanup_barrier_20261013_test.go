package dag

import (
	"context"
	"testing"
)

func nthRoleIndex(log []string, role string, occurrence int) int {
	seen := 0
	for index, entry := range log {
		if entry != role {
			continue
		}
		if seen == occurrence {
			return index
		}
		seen++
	}
	return -1
}

func TestCleanupAwaitedBeforeAdvanceRegression(t *testing.T) {
	mock := newMock()
	plan := makePlan([]map[string]any{issue("setup"), issue("advance")}, [][]string{{"setup"}, {"advance"}})
	git := map[string]any{"integration_branch": "integ/main", "original_branch": "main", "mode": "existing"}

	state, err := RunDAG(context.Background(), plan, "/repo", mock.fn, "swe-planner", testCfg(t, nil), WithGitConfig(git))
	if err != nil {
		t.Fatalf("RunDAG: %v", err)
	}
	if len(state.CompletedIssues) != 2 {
		t.Fatalf("expected both levels to complete, got %d", len(state.CompletedIssues))
	}

	mock.mu.Lock()
	log := append([]string(nil), mock.log...)
	mock.mu.Unlock()

	firstCleanup := nthRoleIndex(log, "run_workspace_cleanup", 0)
	secondSetup := nthRoleIndex(log, "run_workspace_setup", 1)
	if firstCleanup < 0 {
		t.Fatalf("expected cleanup after setup, got cleanup=%d; calls=%v", firstCleanup, log)
	}
	if secondSetup < 0 {
		t.Fatalf("expected second-level setup, got setup=%d; calls=%v", secondSetup, log)
	}
	if firstCleanup > secondSetup {
		t.Fatalf("cleanup was not awaited before advancing: cleanup=%d setup=%d calls=%v", firstCleanup, secondSetup, log)
	}
}

func TestCleanupBarrierRetriesNonFatalFailureBeforeAdvance(t *testing.T) {
	mock := newMock()
	cleanupCalls := 0
	mock.on("run_workspace_cleanup", func(map[string]any) (map[string]any, error) {
		cleanupCalls++
		return map[string]any{"success": false, "cleaned": []any{}}, nil
	})
	plan := makePlan([]map[string]any{issue("setup"), issue("advance")}, [][]string{{"setup"}, {"advance"}})
	git := map[string]any{"integration_branch": "integ/main", "original_branch": "main", "mode": "existing"}

	state, err := RunDAG(context.Background(), plan, "/repo", mock.fn, "swe-planner", testCfg(t, nil), WithGitConfig(git))
	if err != nil {
		t.Fatalf("RunDAG: %v", err)
	}
	if cleanupCalls != 2 {
		t.Fatalf("expected one retry before advancing, got %d cleanup calls", cleanupCalls)
	}
	if len(state.CompletedIssues) != 2 {
		t.Fatalf("expected level advancement after non-fatal cleanup failure, got %d completed", len(state.CompletedIssues))
	}
}

func TestCleanupBarrierHonorsPreCancelledContext(t *testing.T) {
	mock := newMock()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	plan := makePlan([]map[string]any{issue("setup")}, [][]string{{"setup"}})
	git := map[string]any{"integration_branch": "integ/main", "original_branch": "main", "mode": "existing"}

	_, err := RunDAG(ctx, plan, "/repo", mock.fn, "swe-planner", testCfg(t, nil), WithGitConfig(git))
	if err == nil {
		t.Fatal("expected pre-cancelled context error")
	}
	if mock.count("run_workspace_setup") != 0 || mock.count("run_workspace_cleanup") != 0 {
		t.Fatalf("expected no git barrier calls after cancellation, setup=%d cleanup=%d", mock.count("run_workspace_setup"), mock.count("run_workspace_cleanup"))
	}
}
