package harnessx

import (
	"context"
	"os"

	"github.com/Agent-Field/agentfield/sdk/go/agent"
	"github.com/Agent-Field/agentfield/sdk/go/harness"

	"github.com/Agent-Field/SWE-AF/go/internal/hitl"
)

// OpenCodeNoInstallPermissionOverlay is the runtime-owned OpenCode config overlay.
// Target repositories must not need to define SWE's internal FCM provider in
// their own opencode.json just to execute through the broker. The same overlay
// keeps autonomous read/test/build commands available while fail-closing
// dependency installation and cross-worktree virtualenv execution. Role code
// may prepend the current worktree's virtualenv to PATH; absolute sibling-
// worktree runners remain denied.
const OpenCodeNoInstallPermissionOverlay = `{"provider":{"fcm":{"name":"FCM OpenAI-compatible","env":["LLM_BROKER_API_KEY"],"npm":"@ai-sdk/openai-compatible","options":{"apiKey":"{env:LLM_BROKER_API_KEY}","baseURL":"{env:LLM_BROKER_BASE_URL}"},"models":{"fcm":{"id":"fcm","name":"FCM"},"fcm:keyless-dev":{"id":"fcm:keyless-dev","name":"FCM Keyless Dev"}}}},"permission":{"task":"deny","external_directory":"deny","bash":{"*":"allow","pip install *":"deny","pip3 install *":"deny","python -m pip install *":"deny","python3 -m pip install *":"deny","uv pip install *":"deny","uv add *":"deny","npm install *":"deny","npm i *":"deny","pnpm install *":"deny","pnpm add *":"deny","yarn install *":"deny","yarn add *":"deny","go get *":"deny","apt install *":"deny","apt-get install *":"deny","apk add *":"deny","*/.worktrees/*/.venv/bin/*":"deny","*/.worktrees/*/venv/bin/*":"deny"}}}`

// runIDFromContext extracts the build's run ID from the reasoner execution
// context. In production this reads agent.ExecutionContextFrom(ctx).RunID, which
// the SDK populates on the handler's ctx before dispatch. It is a package var
// (not a direct call) purely so tests can inject a run ID — the SDK's context
// key is unexported, so there is no public way to seed ExecutionContext into a
// ctx from an external package.
var runIDFromContext = func(ctx context.Context) string {
	return agent.ExecutionContextFrom(ctx).RunID
}

// HarnessCaller is the minimal method set Run needs from *agent.Agent. Declaring
// it as an interface (rather than depending on the concrete *agent.Agent) lets
// tests supply a mock harness without a live subprocess — the same seam the
// Python tests get by patching router.harness. *agent.Agent satisfies it via
// its Harness method (sdk/go/agent/harness.go:84).
type HarnessCaller interface {
	Harness(ctx context.Context, prompt string, schema map[string]any, dest any, opts harness.Options) (*harness.Result, error)
}

// Run is the single generic entry point every role reasoner uses to invoke the
// harness for structured output of type T.
//
// Sequence (design §4.1, §2.3, §4.3):
//  1. Reflect T into the JSON schema the harness consumes (cached per type).
//  2. Inject the build's run-scoped credentials into opts.Env, scoped creds
//     overriding the base env — mirroring the Python precedence where a freshly
//     minted scout token beats a stale value inherited from os.environ.
//  3. Delegate the full structured-output policy to executeStructured. Keeping
//     this base call seam thin is intentional: weak-model recovery, validation,
//     watchdog salvage, and incremental schema policy live in the adjacent
//     harnessx structured-contract module rather than in role code.
//
// Returns (*T, *harness.Result, error). The Result is returned even alongside a
// non-nil error so callers can inspect diagnostics.
func Run[T any](ctx context.Context, app HarnessCaller, prompt string, opts harness.Options) (*T, *harness.Result, error) {
	schema := schemaFor[T]()
	runID := runIDFromContext(ctx)
	opts.Env = hitl.InjectCredentialsIntoEnv(opts.Env, runID)
	return executeStructured[T](ctx, app, prompt, schema, opts)
}

// RoleOptions is the role→harness parameter mapping (design §4.1). Each role
// reasoner fills it from its resolved config, then ToOptions produces the
// harness.Options passed to Run. Centralizing the mapping here keeps every role
// consistent — the field set mirrors the keyword arguments the Python role
// reasoners pass to router.harness (system_prompt, schema, model, provider,
// tools, cwd, max_turns, permission_mode).
type RoleOptions struct {
	// Provider is the harness ADAPTER string (not the provider), e.g. the output
	// of runtimex.RuntimeToHarnessAdapter — "claude-code", "opencode", "codex".
	// Python passes provider=runtime_to_harness_adapter(ai_provider).
	Provider string

	// Model is the resolved role model identifier.
	Model string

	// MaxTurns caps agent iterations (Python DEFAULT_AGENT_MAX_TURNS per role).
	MaxTurns int

	// Tools is the allowed-tool list (e.g. ["Read","Write","Glob","Grep","Bash"]).
	Tools []string

	// PermissionMode maps to Python's permission_mode; empty means the harness
	// default (Python passes `permission_mode or None`, and harness.Options
	// treats "" as "use default").
	PermissionMode string

	// SystemPrompt is the role's module-level system prompt.
	SystemPrompt string

	// Cwd is the working directory for the subprocess (repo path / worktree).
	Cwd string

	// Env is the base environment for the subprocess. Run overlays the build's
	// scoped credentials on top of this before invoking the harness.
	Env map[string]string

	// SchemaMode is an explicit per-role override for structured-output mode.
	// Empty keeps SWE's centralized policy; narrow roles may request SDK single
	// shot when incremental whole-file edits are empirically less reliable.
	SchemaMode string
}

// ToOptions converts a RoleOptions into a harness.Options. Run injects scoped
// credentials into Env afterwards, so callers leave Env as the base env only.
func (r RoleOptions) ToOptions() harness.Options {
	binPath := ""
	if r.Provider == "opencode" {
		// SWE owns its harness binary contract. Cross-component fallbacks (for
		// example SEC_AF_OPENCODE_BIN) couple independently deployable agents and
		// belong in fleet orchestration, not SWE source.
		binPath = os.Getenv("SWE_OPENCODE_BIN")
	}
	return harness.Options{
		Provider:       r.Provider,
		BinPath:        binPath,
		Model:          r.Model,
		MaxTurns:       r.MaxTurns,
		Tools:          r.Tools,
		PermissionMode: r.PermissionMode,
		SystemPrompt:   r.SystemPrompt,
		Cwd:            r.Cwd,
		ProjectDir:     r.Cwd,
		Env:            r.Env,
		SchemaMode:     r.SchemaMode,
	}
}
