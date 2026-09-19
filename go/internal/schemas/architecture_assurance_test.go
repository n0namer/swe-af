package schemas

import "testing"

func baseArchitecture() Architecture {
	return Architecture{
		Summary:    "x",
		Components: []ArchitectureComponent{{Name: "a", Responsibility: "a"}},
		Interfaces: []string{}, Decisions: []ArchitectureDecision{}, FileChangesOverview: "x",
		Assurance: ArchitectureAssurance{Required: false, Rationale: "single component", ImpactClasses: []string{"local"}, Rules: []ArchitectureAssuranceRule{}, VerificationCommands: []string{}, RuntimeEvidence: []string{}},
	}
}

func TestArchitectureAssuranceSingleComponentMayOptOutWithRationale(t *testing.T) {
	if err := ValidateArchitectureAssurance(baseArchitecture()); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestArchitectureAssuranceSingleComponentStructuralImpactCannotOptOut(t *testing.T) {
	a := baseArchitecture()
	a.Assurance.ImpactClasses = []string{"authority"}
	if err := ValidateArchitectureAssurance(a); err == nil {
		t.Fatal("expected structural impact assurance requirement")
	}
}

func TestArchitectureAssuranceMultiComponentCannotOptOut(t *testing.T) {
	a := baseArchitecture()
	a.Components = append(a.Components, ArchitectureComponent{Name: "b", Responsibility: "b"})
	if err := ValidateArchitectureAssurance(a); err == nil {
		t.Fatal("expected multi-component assurance requirement")
	}
}

func TestArchitectureAssuranceRequiredNeedsProvidersRulesAndCommands(t *testing.T) {
	a := baseArchitecture()
	a.Assurance.Required = true
	if err := ValidateArchitectureAssurance(a); err == nil {
		t.Fatal("expected incomplete required assurance to fail")
	}
}

func TestArchitectureAssuranceRequiredPathNeedsVia(t *testing.T) {
	a := baseArchitecture()
	a.Assurance = ArchitectureAssurance{Required: true, Rationale: "structural", ImpactClasses: []string{"dependency"}, CanonicalModel: "docs/c4/workspace.dsl", CodeRealityProvider: "archsteer", GraphProvider: "codeql", Rules: []ArchitectureAssuranceRule{{ID: "p", Kind: "required_path", Source: "a", Target: "c"}}, VerificationCommands: []string{"make architecture-check"}}
	if err := ValidateArchitectureAssurance(a); err == nil {
		t.Fatal("expected required_path without via to fail")
	}
}

func TestArchitectureAssuranceDuplicateRuleIDsFail(t *testing.T) {
	a := baseArchitecture()
	rule := ArchitectureAssuranceRule{ID: "x", Kind: "required_edge", Source: "a", Target: "b"}
	a.Assurance = ArchitectureAssurance{Required: true, Rationale: "structural", ImpactClasses: []string{"dependency"}, CanonicalModel: "docs/c4/workspace.dsl", CodeRealityProvider: "archsteer", GraphProvider: "codeql", Rules: []ArchitectureAssuranceRule{rule, rule}, VerificationCommands: []string{"make architecture-check"}}
	if err := ValidateArchitectureAssurance(a); err == nil {
		t.Fatal("expected duplicate IDs to fail")
	}
}

func TestArchitectureAssuranceValidRequiredContract(t *testing.T) {
	a := baseArchitecture()
	a.Assurance = ArchitectureAssurance{Required: true, Rationale: "structural", ImpactClasses: []string{"dependency"}, CanonicalModel: "docs/c4/workspace.dsl", CodeRealityProvider: "archsteer", GraphProvider: "codeql", Rules: []ArchitectureAssuranceRule{{ID: "p", Kind: "required_path", Source: "a", Target: "c", Via: []string{"b"}}, {ID: "f", Kind: "forbidden_edge", Source: "a", Target: "c"}}, VerificationCommands: []string{"make architecture-check"}, RuntimeEvidence: []string{"e2e trace"}}
	if err := ValidateArchitectureAssurance(a); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestArchitectureAssuranceRejectsLoopingPathAndCompoundCommand(t *testing.T) {
	a := baseArchitecture()
	a.Assurance = ArchitectureAssurance{
		Required: true, Rationale: "structural", ImpactClasses: []string{"dependency"},
		CanonicalModel: "docs/c4/workspace.dsl", CodeRealityProvider: "archsteer", GraphProvider: "codeql",
		Rules:                []ArchitectureAssuranceRule{{ID: "p", Kind: "required_path", Source: "a", Target: "c", Via: []string{"a"}}},
		VerificationCommands: []string{"archsteer check && rm -rf x"},
	}
	if err := ValidateArchitectureAssurance(a); err == nil {
		t.Fatal("expected looping path/compound command rejection")
	}
	a.Assurance.Rules = []ArchitectureAssuranceRule{{ID: "p", Kind: "required_path", Source: "a", Target: "c", Via: []string{"b"}}}
	if err := ValidateArchitectureAssurance(a); err == nil {
		t.Fatal("expected compound command rejection")
	}
}

func TestArchitectureAssuranceRejectsDuplicateCommands(t *testing.T) {
	a := baseArchitecture()
	a.Assurance = ArchitectureAssurance{Required: true, Rationale: "structural", ImpactClasses: []string{"dependency"}, CanonicalModel: "docs/c4/workspace.dsl", CodeRealityProvider: "archsteer", GraphProvider: "codeql", Rules: []ArchitectureAssuranceRule{{ID: "e", Kind: "required_edge", Source: "a", Target: "b"}}, VerificationCommands: []string{"archsteer check", "archsteer check"}}
	if err := ValidateArchitectureAssurance(a); err == nil {
		t.Fatal("expected duplicate command rejection")
	}
}

func TestArchitectureAssuranceRejectsExternalModelAndMutationCommand(t *testing.T) {
	a := baseArchitecture()
	a.Assurance = ArchitectureAssurance{Required: true, Rationale: "structural", ImpactClasses: []string{"dependency"}, CanonicalModel: "https://example.com/model", CodeRealityProvider: "archsteer", GraphProvider: "codeql", Rules: []ArchitectureAssuranceRule{{ID: "e", Kind: "required_edge", Source: "a", Target: "b"}}, VerificationCommands: []string{"kubectl apply -f x"}}
	if err := ValidateArchitectureAssurance(a); err == nil {
		t.Fatal("expected external canonical model rejection")
	}
	a.Assurance.CanonicalModel = "docs/c4/workspace.dsl"
	if err := ValidateArchitectureAssurance(a); err == nil {
		t.Fatal("expected mutating command rejection")
	}
}
