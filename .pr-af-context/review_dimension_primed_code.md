### go/internal/node/node.go
```
1: // Package node is the wiring wave (T6.2): it constructs the shared *agent.Agent
2: // from the environment and registers every reasoner by its exact Python name so
3: // the Go node is byte-compatible with the Python swe-planner / swe-fast nodes.
4: //
5: // node.go owns agent construction (env -> agent.Config, mirroring app.py:51-59 /
6: // fast/app.py:24-31) and the cross-cutting seam wiring the orchestrators need:
7: // the pause-surface provider (orch.SetPauserProvider) backed by the *agent.Agent
8: // (which pauses via agent.Pause — webhook-resumed), and the hax REST client
9: // resolved from the environment. register.go owns the per-reasoner registration.
10: package node
11: 
12: import (
13: 	"context"
14: 	"fmt"
15: 	"os"
16: 	"strings"
17: 	"time"
18: 
19: 	"github.com/Agent-Field/agentfield/sdk/go/agent"
20: 	"github.com/Agent-Field/agentfield/sdk/go/ai"
21: 
22: 	"github.com/Agent-Field/SWE-AF/go/internal/envelope"
23: 	"github.com/Agent-Field/SWE-AF/go/internal/hitl"
24: 	"github.com/Agent-Field/SWE-AF/go/internal/orch"
25: )
26: 
27: // Node bundles the constructed agent with the resolved environment config and
28: // the collaborators the registration wave threads into every reasoner's Deps.
29: type Node struct {
30: 	// App is the SDK agent. It satisfies every role/orch/fast dependency
31: 	// interface directly (Harness, AI, Note, Call), so the Deps built in
32: 	// register.go point their fields at it.
33: 	App *agent.Agent
34: 
35: 	// NodeID is the resolved node id (NODE_ID env, or the per-binary default).
36: 	NodeID string
37: 
38: 	// AgentFieldServer is the control-plane base URL (AGENTFIELD_SERVER).
39: 	AgentFieldServer string
40: 
41: 	// Token is the control-plane bearer token (AGENTFIELD_API_KEY). Mirrors the
42: 	// Go SDK's own client, which sends AGENTFIELD_API_KEY as Authorization:
43: 	// Bearer (agent.go:593), so the approval client below authenticates the same
44: 	// way the agent does.
45: 	Token string
46: 
47: 	// hax is the hax REST client, nil when HAX_API_KEY is unset (HITL disabled,
48: 	// mirroring build_hax_client_from_env() returning None).
49: 	hax *hitl.HaxClient
50: 
51: 	// registered records every reasoner name passed through regHandler — the
52: 	// single registration path — so the parity test can assert the exact
53: 	// surface. RegisterReasoner is a pure insert keyed by name, so this slice
54: 	// equals the agent's reasoner set (the test also guards against duplicates).
55: 	registered []string
56: 
57: 	// meta records the resolved discovery metadata of every reasoner passed
58: 	// through regHandler, keyed by name. The SDK keeps its own reasoner map
59: 	// unexported, so this is the only way the surface tests can assert what a
60: 	// caller actually sees on the control plane.
61: 	meta map[string]ReasonerMeta
62: }
63: 
64: // ReasonerMeta is the discovery-facing metadata a reasoner registers with: the
65: // tags and the description `af ls` and GET /api/v1/discovery/capabilities show
66: // a caller deciding which reasoner to invoke.
67: type ReasonerMeta struct {
68: 	Tags        []string
69: 	Description string
70: }
71: 
72: // RegisteredNames returns a copy of the reasoner names registered on this node,
73: // in registration order. Used by the parity test to assert the surface exactly
74: // matches the Python node.
75: func (n *Node) RegisteredNames() []string {
76: 	return append([]string(nil), n.registered...)
77: }
78: 
79: // RegisteredMeta returns a copy of the registration metadata keyed by reasoner
80: // name. Used by the surface tests to assert the entrypoint tagging and the
81: // internal-stage markers a discovering caller routes on.
82: func (n *Node) RegisteredMeta() map[string]ReasonerMeta {
83: 	out := make(map[string]ReasonerMeta, len(n.meta))
84: 	for name, m := range n.meta {
85: 		m.Tags = append([]string(nil), m.Tags...)
86: 		out[name] = m
87: 	}
88: 	return out
89: }
90: 
91: // BuildAgent constructs the SWE-AF agent from the environment exactly as the
92: // Python entry points do (app.py:51-59 / fast/app.py:24-31):
93: //
94: //   - NODE_ID           default defaultNodeID ("swe-planner" / "swe-fast")
95: //   - AGENTFIELD_SERVER default "http://localhost:8080"
96: //   - AGENTFIELD_API_KEY -> Config.Token (bearer)
97: //   - PORT              default defaultPort ("8005" / "8006") -> ListenAddress
98: //   - AGENT_CALLBACK_URL -> Config.PublicURL — the base URL the node registers
99: //     with the control plane. The Python SDK reads the same env var
100: //     (agent_server.py:758) before defaulting to localhost; without it a
101: //     containerized node registers http://localhost:<port> and the CP cannot
102: //     route execute calls back to it (found by the T7.2 functional tests).
103: //     Unset -> the SDK's localhost default, matching Python.
104: //   - Version           "1.0.0"
105: //   - description       the per-node description string
106: //
107: // It also wires the orchestrator pause seam (orch.SetPauserProvider) so the
108: // plan-approval gate can pause via the control plane. Register the reasoners
109: // with RegisterPlanner / RegisterFast.
110: func BuildAgent(defaultNodeID, defaultPort, description string) (*Node, error) {
111: 	nodeID := envOr("NODE_ID", defaultNodeID)
112: 	server := envOr("AGENTFIELD_SERVER", "http://localhost:8080")
113: 	token := os.Getenv("AGENTFIELD_API_KEY")
114: 	port := envOr("PORT", defaultPort)
115: 
116: 	cfg := agent.Config{
117: 		NodeID:        nodeID,
118: 		Version:       "1.0.0",
119: 		AgentFieldURL: server,
120: 		Token:         token,
121: 		ListenAddress: ":" + port,
122: 		// PublicURL: the callback base URL the CP uses to reach this node.
123: 		// AGENT_CALLBACK_URL mirrors the Python SDK's env var of the same name
124: 		// (agent_server.py:758); when unset the SDK defaults to
125: 		// http://localhost:<port>, which is correct only outside containers.
126: 		PublicURL: os.Getenv("AGENT_CALLBACK_URL"),
127: 		// The Go SDK sends no node-level description in its registration payload
128: 		// (unlike the Python Agent(description=...)); AppDescription is the CLI
129: 		// help string, the closest home for the Python description. It does not
130: 		// enable CLI mode (no reasoner sets WithCLI), so Run still serves.
131: 		CLIConfig: &agent.CLIConfig{AppDescription: description},
132: 	}
133: 
134: 	// Enable the direct-LLM path (run_qa_synthesizer, the Go equivalent of
135: 	// Python's router.ai) when the environment supplies a usable AI config. An
136: 	// absent or invalid key leaves AIConfig nil so the agent still constructs and
137: 	// the synthesizer falls back deterministically — without this the LLM branch
138: 	// was unreachable (agent.AI returns "AI not configured").
139: 	cfg.AIConfig = resolveAIConfig()
140: 
141: 	app, err := agent.New(cfg)
142: 	if err != nil {
143: 		return nil, fmt.Errorf("create agent %q: %w", nodeID, err)
144: 	}
145: 
146: 	n := &Node{
147: 		App:              app,
148: 		NodeID:           nodeID,
149: 		AgentFieldServer: server,
150: 		Token:            token,
151: 		hax:              hitl.BuildHaxClientFromEnv(),
152: 	}
153: 
154: 	// Wire the orchestrator pause seam once: the plan-approval gate pauses the
155: 	// execution through this provider. The *agent.Agent satisfies hitl.Pauser
156: 	// via its Pause method (webhook-resumed — no polling); the request-scoped
157: 	// ApprovalRequest is ignored, as a single process-wide agent serves every
158: 	// execution.
159: 	orch.SetPauserProvider(func(orch.ApprovalRequest) hitl.Pauser {
160: 		return n.App
161: 	})
162: 
163: 	return n, nil
164: }
165: 
166: // newCallFn returns the app.Call + envelope-unwrap closure injected into the
167: // coding loop, the DAG executor and the fast pipeline. It is structurally
168: // identical to coding.CallFn and fast.CallFn (both func(ctx, target, kwargs)
169: // (map[string]any, error)), so a single closure satisfies both. The unwrap
170: // label is the reasoner segment after the final dot — matching orch.NewCallFn.
171: func newCallFn(app *agent.Agent) func(context.Context, string, map[string]any) (map[string]any, error) {
172: 	return func(ctx context.Context, target string, kwargs map[string]any) (map[string]any, error) {
173: 		raw, err := app.Call(ctx, target, kwargs)
174: 		if err != nil {
175: 			return nil, err
176: 		}
177: 		label := target
178: 		if i := strings.LastIndex(target, "."); i >= 0 {
179: 			label = target[i+1:]
180: 		}
181: 		return envelope.UnwrapCallResult(raw, label)
182: 	}
183: }
184: 
185: // resolveAIConfig returns ai.DefaultConfig() when it validates (an OPENAI_API_KEY
186: // or OPENROUTER_API_KEY is set), else nil. Returning nil — rather than a broken
187: // config — keeps node startup working with no key: the QA-synthesizer LLM branch
188: // is simply disabled and its deterministic fallback runs instead.
189: func resolveAIConfig() *ai.Config {
190: 	c := ai.DefaultConfig()
191: 	// The AgentField Go SDK reads AI_BASE_URL, while the harness/deployment
192: 	// contract uses OPENAI_BASE_URL for OpenAI-compatible providers such as
193: 	// Gonka. Preserve that configured endpoint explicitly instead of silently
194: 	// falling back to the SDK default OpenAI base.
195: 	if strings.TrimSpace(os.Getenv("AI_BASE_URL")) == "" {
196: 		if base := strings.TrimSpace(os.Getenv("OPENAI_BASE_URL")); base != "" {
197: 			c.BaseURL = base
198: 		}
199: 	}
200: 	// Match the proven Deep Research semantic transport contract: provider calls
201: 	// may legitimately run for minutes, so the default 30s Go SDK client timeout
202: 	// is too short for Gonka. Keep a 30-minute transport safety ceiling.
203: 	c.Timeout = 30 * time.Minute
204: 	if c.Validate() == nil {
205: 		return c
206: 	}
207: 	return nil
208: }
209: 
210: // envOr returns the value of key, or def when the env var is unset or empty.
211: func envOr(key, def string) string {
212: 	if v := os.Getenv(key); v != "" {
213: 		return v
214: 	}
215: 	return def
216: }
```
_import/usage context:_ IMPORTS: import (
IMPORTED BY: none

### go/agentfield-package.yaml
```
1: config_version: v1
2: # This is THE SWE node. It deliberately shares the root manifest's name: the
3: # root declares `superseded_by` pointing here, so installing this repo installs
4: # this package, and a user who already has the Python swe-planner gets it
5: # replaced in place — same name, same node id, same triggers, secrets kept.
6: # Installing the root as a local path (the documented escape hatch) is the one
7: # way to get the Python node, and it necessarily takes this name over.
8: name: swe-planner
9: version: 0.1.0
10: description: Autonomous SWE planning/execution agent node (Go port; 5 orchestrators + 25 role reasoners)
11: author: Agent-Field
12: language: go                            # explicit (also auto-detected from go.mod)
13: 
14: entrypoint:
15:   build: ./cmd/swe-planner
16:   start: bin/swe-planner
17:   healthcheck: /health
18: 
19: agent_node:
20:   node_id: swe-planner
21:   # 8005 rather than the Python node's 8003: during the changeover both may be
22:   # running, and triggers resolve by node id, not port.
23:   default_port: 8005
24: 
25: user_environment:
26:   require_one_of:
27:     # SWE-AF accepts Anthropic, OpenRouter, or an OpenAI-compatible provider.
28:     # The permanent DEV lane uses Gonka through OPENAI_API_KEY + OPENAI_BASE_URL.
29:     - id: llm_provider
30:       description: an LLM provider key
31:       options:
32:         - name: ANTHROPIC_API_KEY
33:           description: Anthropic API key (Claude)
34:           type: secret
35:           scope: global
36:         - name: OPENROUTER_API_KEY
37:           description: OpenRouter API key (DeepSeek/Qwen/Llama/… — 200+ models)
38:           type: secret
39:           scope: global
40:         - name: OPENAI_API_KEY
41:           description: OpenAI-compatible API key (including Gonka)
42:           type: secret
43:           scope: global
44:   optional:
45:     # Optional so an LLM key is the only secret needed to get started: builds
46:     # on local/public repos run without it; it is needed to clone private
47:     # repos, push branches, and open pull requests.
48:     - name: GH_TOKEN
49:       description: GitHub token (repo scope) — needed to clone private repos, push branches, and open pull requests
50:       type: secret
51:       scope: global
52:     - name: SWE_DEFAULT_RUNTIME
53:       description: Coding runtime for every role (claude_code | open_code | codex)
54:       default: open_code
55:     - name: SWE_OPENCODE_BIN
56:       description: >-
57:         Path to the SWE-owned OpenCode wrapper/binary used when the coding runtime
58:         is open_code. Set this explicitly when the wrapper is not discoverable by
59:         the node runner.
60:     - name: SWE_DEFAULT_MODEL
61:       description: Virtual FCM model id used by every SWE role; FCM owns upstream model routing
62:       default: fcm/fcm
63:     - name: AGENTFIELD_HARNESS_IDLE_SECONDS
64:       description: >-
65:         Idle-output watchdog for harness CLI subprocesses. SWE defaults this to 0
66:         because the installed OpenCode wrapper buffers stdout/stderr until the
67:         command exits; the SDK total timeout remains the bounded liveness guard.
68:       default: "0"
69:     - name: AGENT_CALLBACK_URL
70:       description: >-
71:         Public base URL the control plane uses to call this node. Required for
72:         containerized or remote deployments where localhost is not reachable
73:         from the control plane.
74:     - name: SWE_OPENCLAW_HITL
75:       description: >-
76:         Use OpenClaw/Telegram as the human-in-the-loop surface when HAX is not
77:         configured. Opt-in: set to 1 only where OpenClaw can reach this node's
78:         AgentField logs and /webhooks/approval callback.
79:       default: "0"
80:     - name: SWE_PRO_ENGINE
81:       description: >-
82:         High-performance coding engine (beta). On by default for nodes
83:         installed this way — set it to 0 (or false) to use the classic coding
84:         loop instead.
85:       # The resolver injects this last unless the process env or a stored
86:       # secret overrides it, so an `af install` gets the engine without asking
87:       # and `SWE_PRO_ENGINE=0 af run` still opts out.
88:       default: "1"
89:     - name: SWE_PRO_VARIANT
90:       description: Engine reasoning-effort variant (low | high) — unset keeps the engine default
91:     - name: SWE_PRO_MAX_COST
92:       description: Per-run USD ceiling for the engine — unset means no per-run cap
93:     - name: AGENTFIELD_SERVER
94:       description: Control-plane URL
95:       default: http://localhost:8080
96:     - name: AGENTFIELD_API_KEY
97:       description: Control-plane API key (if auth is enabled)
98:       type: secret
99:       scope: global
100:   # NODE_ID and PORT are deliberately NOT declared here: the installer's
101:   # runner assigns the port (injected as PORT) and the binary defaults the
102:   # identity — a manifest default would override the runner's assignment,
103:   # because resolved user_environment values are appended last to the
104:   # process env.
```
_import/usage context:_ IMPORTS: from the control plane.
IMPORTED BY: none