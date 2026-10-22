package node

import (
	"context"
	"testing"

	"github.com/Agent-Field/SWE-AF/go/internal/fast"
	"github.com/Agent-Field/SWE-AF/go/internal/orch"
)

func TestFactoryBuildDefaultsToFull(t *testing.T) {
	fullCalled, fastCalled := false, false
	full := func(_ context.Context, _ *orch.Deps, input map[string]any) (any, error) {
		fullCalled = true
		if _, exists := input["mode"]; exists {
			t.Fatal("mode must not leak into the full build handler")
		}
		return map[string]any{"selected": "full"}, nil
	}
	quick := func(_ context.Context, _ *fast.Deps, _ map[string]any) (any, error) {
		fastCalled = true
		return map[string]any{"selected": "fast"}, nil
	}
	h := factoryBuildHandler(full, quick, &fast.Deps{})
	out, err := h(context.Background(), &orch.Deps{}, map[string]any{"goal": "test"})
	if err != nil {
		t.Fatal(err)
	}
	if out.(map[string]any)["selected"] != "full" || !fullCalled || fastCalled {
		t.Fatalf("dispatch result=%v full=%v fast=%v", out, fullCalled, fastCalled)
	}
}

func TestFactoryBuildFastIsExplicitMode(t *testing.T) {
	fullCalled, fastCalled := false, false
	full := func(_ context.Context, _ *orch.Deps, _ map[string]any) (any, error) {
		fullCalled = true
		return map[string]any{"selected": "full"}, nil
	}
	quick := func(_ context.Context, _ *fast.Deps, input map[string]any) (any, error) {
		fastCalled = true
		if _, exists := input["mode"]; exists {
			t.Fatal("mode must not leak into the fast build handler")
		}
		if input["goal"] != "test" {
			t.Fatalf("goal not forwarded: %v", input)
		}
		return map[string]any{"selected": "fast"}, nil
	}
	h := factoryBuildHandler(full, quick, &fast.Deps{})
	out, err := h(context.Background(), &orch.Deps{}, map[string]any{"goal": "test", "mode": " FAST "})
	if err != nil {
		t.Fatal(err)
	}
	if out.(map[string]any)["selected"] != "fast" || !fastCalled || fullCalled {
		t.Fatalf("dispatch result=%v full=%v fast=%v", out, fullCalled, fastCalled)
	}
}

func TestFactoryBuildRejectsUnknownMode(t *testing.T) {
	called := false
	full := func(context.Context, *orch.Deps, map[string]any) (any, error) {
		called = true
		return nil, nil
	}
	quick := func(context.Context, *fast.Deps, map[string]any) (any, error) {
		called = true
		return nil, nil
	}
	h := factoryBuildHandler(full, quick, &fast.Deps{})
	if _, err := h(context.Background(), &orch.Deps{}, map[string]any{"goal": "test", "mode": "turbo"}); err == nil {
		t.Fatal("unknown mode must fail")
	}
	if called {
		t.Fatal("unknown mode must be rejected before invoking a build handler")
	}
}
