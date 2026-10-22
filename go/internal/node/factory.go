package node

import (
	"context"
	"fmt"
	"strings"

	"github.com/Agent-Field/SWE-AF/go/internal/fast"
	"github.com/Agent-Field/SWE-AF/go/internal/orch"
)

// factoryBuildHandler presents Full and Fast as execution modes of one AgentField
// entrypoint. Full is the safe default; the mode selector is not forwarded to the
// underlying Python-compatible build handlers.
func factoryBuildHandler(full orch.Handler, quick fast.Handler, quickDeps *fast.Deps) orch.Handler {
	return func(ctx context.Context, fullDeps *orch.Deps, input map[string]any) (any, error) {
		mode := "full"
		if raw, exists := input["mode"]; exists {
			value, ok := raw.(string)
			if !ok {
				return nil, fmt.Errorf("mode must be a string (full or fast)")
			}
			if strings.TrimSpace(value) != "" {
				mode = strings.ToLower(strings.TrimSpace(value))
			}
		}

		forward := make(map[string]any, len(input))
		for key, value := range input {
			if key != "mode" {
				forward[key] = value
			}
		}

		switch mode {
		case "full":
			return full(ctx, fullDeps, forward)
		case "fast":
			return quick(ctx, quickDeps, forward)
		default:
			return nil, fmt.Errorf("unsupported mode %q: expected full or fast", mode)
		}
	}
}
