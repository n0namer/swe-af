package harnessx

import (
	"encoding/json"
	"testing"

	tekjsonschema "github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
	"strings"
)

func TestProbeCoderDefaultValidates(t *testing.T) {
	schema := schemaFor[schemas.CoderResult]()
	sb, _ := json.Marshal(schema)
	t.Logf("SCHEMA: %s", sb)

	// Compile exactly as recoverStructuredText does (checks $schema draft acceptance).
	sbb, _ := json.Marshal(schema)
	compiler := tekjsonschema.NewCompiler()
	if err := compiler.AddResource("mem://probe.json", bytesReader(sbb)); err != nil {
		t.Fatalf("add: %v", err)
	}
	compiled, err := compiler.Compile("mem://probe.json")
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	// 1. Default-seeded CoderResult (what UnmarshalJSON materializes) re-marshaled.
	var seeded schemas.CoderResult
	if err := json.Unmarshal([]byte(`{}`), &seeded); err != nil {
		t.Fatalf("seed: %v", err)
	}
	nb, _ := json.Marshal(seeded)
	t.Logf("DEFAULT-SEEDED MARSHAL: %s", nb)
	var data any
	_ = json.Unmarshal(nb, &data)
	if err := compiled.Validate(data); err != nil {
		t.Logf("DEFAULT-SEEDED VALIDATION FAIL: %v", err)
	} else {
		t.Logf("DEFAULT-SEEDED VALIDATION OK")
	}

	// 2. Full recoverStructuredText path with a weak-model candidate.
	var dest schemas.CoderResult
	err = recoverStructuredText(`{"complete": true, "agent_retro": "plain string retro"}`, schema, &dest)
	t.Logf("RECOVER err=%v dest=%+v", err, dest)
}

func TestProbeRecoverMissingFields(t *testing.T) {
	schema := schemaFor[schemas.CoderResult]()
	var dest schemas.CoderResult
	err := recoverStructuredText(`{"summary":"did the thing"}`, schema, &dest)
	t.Logf("RECOVER-MISSING err=%v dest=%+v", err, dest)
}
