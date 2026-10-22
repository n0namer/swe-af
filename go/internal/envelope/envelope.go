// Package envelope unwraps execution envelopes returned by agent.Call.
//
// It is a 1:1 port of swe_af/execution/envelope.py. The AgentField SDK's Call
// has two execution paths: the async path unwraps the result and surfaces
// failures, while the sync fallback path may return the full execution envelope
// (with the inner result nil) on failed executions. UnwrapCallResult normalises
// both cases so pipeline code can always expect the raw reasoner output.
package envelope

import (
	"fmt"
	"strings"

	"github.com/Agent-Field/SWE-AF/go/internal/fatal"
)

// Binding is the optional G8 provenance and correlation binding carried by
// execution envelopes and checkpoints. All fields are optional so envelopes
// produced by older callers remain valid.
type Binding struct {
	CapabilityID  string `json:"capability_id,omitempty"`
	DefinitionID  string `json:"definition_id,omitempty"`
	SourceRevision string `json:"source_revision,omitempty"`
	ParentID      string `json:"parent_id,omitempty"`
	RootID        string `json:"root_id,omitempty"`
	ThreadID      string `json:"thread_id,omitempty"`
}

var bindingKeys = []string{"capability_id", "definition_id", "source_revision", "parent_id", "root_id", "thread_id"}

// BindingFrom extracts the canonical binding from a map. AgentField's native
// names are accepted as aliases for parent/root correlation, and nested
// provenance/correlation maps are accepted for compatibility with new clients.
func BindingFrom(values map[string]any) Binding {
	var b Binding
	if values == nil {
		return b
	}
	read := func(key string) string {
		if v, ok := values[key].(string); ok {
			return strings.TrimSpace(v)
		}
		return ""
	}
	b.CapabilityID, b.DefinitionID, b.SourceRevision = read("capability_id"), read("definition_id"), read("source_revision")
	b.ParentID, b.RootID, b.ThreadID = read("parent_id"), read("root_id"), read("thread_id")
	if b.ParentID == "" { b.ParentID = read("parent_execution_id") }
	if b.RootID == "" { b.RootID = read("root_workflow_id") }
	for key := range map[string]struct{}{"provenance": {}, "correlation": {}} {
		nested, ok := values[key].(map[string]any)
		if !ok { continue }
		if key == "provenance" {
			if b.CapabilityID == "" { b.CapabilityID = mapString(nested, "capability_id") }
			if b.DefinitionID == "" { b.DefinitionID = mapString(nested, "definition_id") }
			if b.SourceRevision == "" { b.SourceRevision = mapString(nested, "source_revision") }
		} else {
			if b.ParentID == "" { b.ParentID = mapString(nested, "parent_id") }
			if b.RootID == "" { b.RootID = mapString(nested, "root_id") }
			if b.ThreadID == "" { b.ThreadID = mapString(nested, "thread_id") }
		}
	}
	return b
}

func mapString(m map[string]any, key string) string {
	if s, ok := m[key].(string); ok { return strings.TrimSpace(s) }
	return ""
}

// HasBinding reports whether values carry at least one binding field.
func HasBinding(values map[string]any) bool {
	if values == nil { return false }
	for _, key := range bindingKeys { if _, ok := values[key]; ok { return true } }
	_, p := values["provenance"]
	_, c := values["correlation"]
	return p || c || values["parent_execution_id"] != nil || values["root_workflow_id"] != nil
}

// ValidateBinding checks only fields present in expected. Empty expected fields
// are ignored, preserving compatibility with legacy calls and partial bindings.
func ValidateBinding(actual, expected map[string]any) error {
	if !HasBinding(expected) { return nil }
	a, e := BindingFrom(actual), BindingFrom(expected)
	for name, pair := range map[string][2]string{
		"capability_id": {a.CapabilityID, e.CapabilityID}, "definition_id": {a.DefinitionID, e.DefinitionID},
		"source_revision": {a.SourceRevision, e.SourceRevision}, "parent_id": {a.ParentID, e.ParentID},
		"root_id": {a.RootID, e.RootID}, "thread_id": {a.ThreadID, e.ThreadID},
	} {
		av, ev := pair[0], pair[1]
		if ev != "" && av != ev { return fmt.Errorf("provenance/correlation mismatch for %s: got %q, want %q", name, av, ev) }
	}
	return nil
}

// envelopeKeys is the set of keys present in the execution envelope returned by
// the control plane's _build_execute_response. Verbatim port of _ENVELOPE_KEYS.
var envelopeKeys = map[string]struct{}{
	"execution_id":  {},
	"run_id":        {},
	"node_id":       {},
	"type":          {},
	"target":        {},
	"status":        {},
	"duration_ms":   {},
	"timestamp":     {},
	"result":        {},
	"error_message": {},
	"cost":          {},
	"capability_id": {},
	"definition_id": {},
	"source_revision": {},
	"parent_id":     {},
	"root_id":       {},
	"thread_id":     {},
	"parent_execution_id": {},
	"root_workflow_id": {},
	"provenance":    {},
	"correlation":   {},
}

// failureStatuses are the terminal-failure envelope statuses.
var failureStatuses = map[string]struct{}{
	"failed":    {},
	"error":     {},
	"cancelled": {},
	"timeout":   {},
}

// UnwrapCallResult extracts the actual reasoner output from an agent.Call
// response.
//
// raw is the value returned by agent.Call (or the equivalent call_fn); label is
// a human-readable name (e.g. the target) used in error messages.
//
// Behaviour (mirrors unwrap_call_result):
//   - No envelope keys present -> already unwrapped, return raw unchanged.
//   - Envelope with a terminal-failure status -> a *fatal.FatalHarnessError when
//     the error is fatal, otherwise a plain error mirroring Python's
//     RuntimeError: "<label> failed (status=<status>): <err>".
//   - Otherwise return raw["result"] when present (and non-nil), else raw.
func UnwrapCallResult(raw map[string]any, label string) (map[string]any, error) {
	if raw == nil {
		return raw, nil
	}

	// Fast path: already unwrapped (no envelope keys present).
	if !hasEnvelopeKey(raw) {
		return raw, nil
	}

	// Looks like the execution envelope — check for errors first.
	status := ""
	if s, ok := raw["status"]; ok && s != nil {
		status = strings.ToLower(fmt.Sprintf("%v", s))
	}
	if _, isFailure := failureStatuses[status]; isFailure {
		err := errorText(raw)
		if fatal.IsFatalError(err) {
			return nil, &fatal.FatalHarnessError{OriginalMessage: err}
		}
		return nil, fmt.Errorf("%s failed (status=%s): %s", label, status, err)
	}

	inner, ok := raw["result"]
	if ok && inner != nil {
		if m, isMap := inner.(map[string]any); isMap {
			return m, nil
		}
		// Non-object result (e.g. a scalar or list). Python returns it as-is
		// via an Any return, but reasoner outputs are always JSON objects
		// (model_dump dicts); fall through to return the envelope so the caller
		// can validate rather than losing information.
	}

	// Envelope present but result is nil (or non-object) and status isn't a
	// known failure — return as-is (caller should validate).
	return raw, nil
}

// hasEnvelopeKey reports whether raw contains any envelope key (the Go
// equivalent of _ENVELOPE_KEYS.intersection(result)).
func hasEnvelopeKey(raw map[string]any) bool {
	for k := range raw {
		if _, ok := envelopeKeys[k]; ok {
			return true
		}
	}
	return false
}

// errorText reproduces Python's `result.get("error_message") or
// result.get("error") or "unknown"` truthiness chain, then str()-ifies the
// chosen value.
func errorText(raw map[string]any) string {
	if s, ok := truthy(raw["error_message"]); ok {
		return s
	}
	if s, ok := truthy(raw["error"]); ok {
		return s
	}
	return "unknown"
}

// truthy returns the string form of v and whether v is Python-truthy for the
// purpose of the `or` chain (nil and the empty string are falsy).
func truthy(v any) (string, bool) {
	if v == nil {
		return "", false
	}
	if s, ok := v.(string); ok {
		return s, s != ""
	}
	return fmt.Sprintf("%v", v), true
}
