package orch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var r16BuildSignature func(context.Context, *Deps, map[string]any) (any, error) = Build

func TestR16ResumeSourceContract(t *testing.T) {
	source, err := os.ReadFile("build.go")
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(source)
	if actual := hex.EncodeToString(digest[:]); actual != "05c39a2e3b41869b01a1482543c9875630a0777fbea24ea706c05727fe640cbc" {
		t.Fatalf("Build source identity changed: %s", actual)
	}
	if r16BuildSignature == nil || Handlers()["build"] == nil {
		t.Fatal("missing Build registration")
	}
	expected := []struct{ name, kind, tag string }{
		{"Goal", "string", "goal"},
		{"RepoPath", "string", "repo_path"},
		{"RepoURL", "string", "repo_url"},
		{"ArtifactsDir", "string", "artifacts_dir"},
		{"AdditionalContext", "string", "additional_context"},
		{"Config", "map[string]interface {}", "config"},
		{"ExecuteFnTarget", "string", "execute_fn_target"},
		{"MaxTurns", "int", "max_turns"},
		{"PermissionMode", "string", "permission_mode"},
		{"EnableLearning", "bool", "enable_learning"},
	}
	actual := reflect.TypeOf(buildInput{})
	if actual.NumField() != len(expected) {
		t.Fatalf("input field count: %d", actual.NumField())
	}
	for index, field := range expected {
		got := actual.Field(index)
		if got.Name != field.name || got.Type.String() != field.kind || got.Tag.Get("json") != field.tag {
			t.Fatalf("input field %d changed: %+v", index, got)
		}
	}
}

type r16OracleState struct {
	Revision int               `json:"revision"`
	Epoch    int               `json:"epoch"`
	Owner    string            `json:"owner"`
	Results  map[string]string `json:"results"`
}

type r16OracleToken struct {
	revision int
	epoch    int
	owner    string
}

type r16ResumeOracle struct {
	directory string
}

func (oracle r16ResumeOracle) load() (r16OracleState, error) {
	state := r16OracleState{Results: map[string]string{}}
	for revision := 1; ; revision++ {
		data, err := os.ReadFile(filepath.Join(oracle.directory, fmt.Sprintf("%06d.json", revision)))
		if errors.Is(err, os.ErrNotExist) {
			return state, nil
		}
		if err != nil {
			return state, err
		}
		var next r16OracleState
		if err := json.Unmarshal(data, &next); err != nil {
			return state, err
		}
		if next.Revision != revision || next.Epoch < state.Epoch || next.Owner == "" || next.Results == nil {
			return state, errors.New("invalid checkpoint")
		}
		state = next
	}
}

func (oracle r16ResumeOracle) persist(state r16OracleState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(oracle.directory, fmt.Sprintf("%06d.json", state.Revision)), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func (oracle r16ResumeOracle) takeover(expected int, owner string) (r16OracleToken, error) {
	state, err := oracle.load()
	if err != nil {
		return r16OracleToken{}, err
	}
	if owner == "" || state.Revision != expected {
		return r16OracleToken{}, errors.New("takeover CAS rejected")
	}
	state.Revision++
	state.Epoch++
	state.Owner = owner
	if err := oracle.persist(state); err != nil {
		return r16OracleToken{}, err
	}
	return r16OracleToken{state.Revision, state.Epoch, state.Owner}, nil
}

func (oracle r16ResumeOracle) execute(ctx context.Context, token r16OracleToken, issue, result string, dependencies ...string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	state, err := oracle.load()
	if err != nil {
		return "", err
	}
	if token.owner != state.Owner || token.epoch != state.Epoch {
		return "", errors.New("stale owner fenced")
	}
	if issue == "" || result == "" {
		return "", errors.New("empty issue or result")
	}
	if previous, exists := state.Results[issue]; exists {
		if previous != result {
			return "", errors.New("conflicting replay")
		}
		return previous, nil
	}
	if token.revision != state.Revision {
		return "", errors.New("checkpoint CAS rejected")
	}
	for _, dependency := range dependencies {
		if _, exists := state.Results[dependency]; !exists {
			return "", errors.New("dependency incomplete")
		}
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	state.Revision++
	state.Results[issue] = result
	if err := oracle.persist(state); err != nil {
		return "", err
	}
	return result, nil
}

func r16ReadState(t *testing.T, oracle r16ResumeOracle) r16OracleState {
	t.Helper()
	state, err := oracle.load()
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func r16AssertUnchanged(t *testing.T, oracle r16ResumeOracle, before r16OracleState) {
	t.Helper()
	if after := r16ReadState(t, oracle); !reflect.DeepEqual(before, after) {
		t.Fatalf("rejected operation changed checkpoint: before=%+v after=%+v", before, after)
	}
	entries, err := os.ReadDir(oracle.directory)
	if err != nil || len(entries) != before.Revision {
		t.Fatalf("unexpected checkpoint files: count=%d err=%v", len(entries), err)
	}
}

func TestR16ResumeOracleDependenciesCancellationAndCAS(t *testing.T) {
	oracle := r16ResumeOracle{t.TempDir()}
	token, err := oracle.takeover(0, "owner-a")
	if err != nil {
		t.Fatal(err)
	}
	before := r16ReadState(t, oracle)
	if _, err := oracle.execute(context.Background(), token, "child", "child-receipt", "parent"); err == nil {
		t.Fatal("incomplete dependency allowed")
	}
	r16AssertUnchanged(t, oracle, before)
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := oracle.execute(cancelled, token, "parent", "parent-receipt"); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation not preserved: %v", err)
	}
	r16AssertUnchanged(t, oracle, before)
	if result, err := oracle.execute(context.Background(), token, "parent", "parent-receipt"); err != nil || result != "parent-receipt" {
		t.Fatalf("positive parent: result=%q err=%v", result, err)
	}
	before = r16ReadState(t, oracle)
	if _, err := oracle.execute(context.Background(), token, "child", "child-receipt", "parent"); err == nil {
		t.Fatal("stale checkpoint CAS allowed")
	}
	r16AssertUnchanged(t, oracle, before)
	token.revision = before.Revision
	if _, err := oracle.execute(context.Background(), token, "child", "child-receipt", "parent"); err != nil {
		t.Fatal(err)
	}
	if state := r16ReadState(t, oracle); len(state.Results) != 2 || state.Results["child"] != "child-receipt" {
		t.Fatalf("durable child receipt missing: %+v", state)
	}
}

func TestR16ResumeOracleRestartFencingAndIdempotentResult(t *testing.T) {
	oracle := r16ResumeOracle{t.TempDir()}
	old, err := oracle.takeover(0, "owner-a")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := oracle.execute(context.Background(), old, "parent", "durable-result"); err != nil {
		t.Fatal(err)
	}
	restarted := r16ResumeOracle{oracle.directory}
	checkpoint := r16ReadState(t, restarted)
	current, err := restarted.takeover(checkpoint.Revision, "owner-b")
	if err != nil {
		t.Fatal(err)
	}
	before := r16ReadState(t, restarted)
	if _, err := oracle.execute(context.Background(), old, "parent", "durable-result"); err == nil {
		t.Fatal("stale owner replay allowed")
	}
	r16AssertUnchanged(t, restarted, before)
	if _, err := restarted.takeover(checkpoint.Revision, "owner-c"); err == nil {
		t.Fatal("stale takeover CAS allowed")
	}
	r16AssertUnchanged(t, restarted, before)
	for replay := 0; replay < 2; replay++ {
		if result, err := restarted.execute(context.Background(), current, "parent", "durable-result"); err != nil || result != "durable-result" {
			t.Fatalf("restart replay: result=%q err=%v", result, err)
		}
		r16AssertUnchanged(t, restarted, before)
	}
	if _, err := restarted.execute(context.Background(), current, "parent", "different-result"); err == nil {
		t.Fatal("conflicting idempotency result allowed")
	}
	r16AssertUnchanged(t, restarted, before)
	if _, err := restarted.execute(context.Background(), current, "child", "child-result", "parent"); err != nil {
		t.Fatal(err)
	}
}

func TestR16ResumeOraclePartialCheckpointFailsClosed(t *testing.T) {
	oracle := r16ResumeOracle{t.TempDir()}
	if err := os.WriteFile(filepath.Join(oracle.directory, "000001.json"), []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := oracle.takeover(0, "owner-a"); err == nil {
		t.Fatal("partial checkpoint accepted")
	}
	entries, err := os.ReadDir(oracle.directory)
	if err != nil || len(entries) != 1 {
		t.Fatalf("partial checkpoint caused additional writes: %d %v", len(entries), err)
	}
}
