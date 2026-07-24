package contracts_test

import (
	"encoding/json"
	"os"
	"testing"

	contracts "github.com/truzhen/contracts"
)

func TestDeliberationSchemasEmbeddedClosedAndAuthoritySafe(t *testing.T) {
	schemas := map[string][]byte{
		"session":          contracts.DeliberationSessionReadModelSchemaJSON,
		"turn":             contracts.DeliberationTurnReadModelSchemaJSON,
		"provider-lane":    contracts.DeliberationProviderLaneReadModelSchemaJSON,
		"automation-grant": contracts.DeliberationAutomationGrantReadModelSchemaJSON,
		"synthesis":        contracts.DeliberationSynthesisCandidateSchemaJSON,
	}
	for name, raw := range schemas {
		t.Run(name, func(t *testing.T) {
			var document struct {
				AdditionalProperties bool                       `json:"additionalProperties"`
				Required             []string                   `json:"required"`
				Properties           map[string]json.RawMessage `json:"properties"`
			}
			if err := json.Unmarshal(raw, &document); err != nil {
				t.Fatalf("schema must be valid JSON: %v", err)
			}
			if document.AdditionalProperties || len(document.Required) == 0 || len(document.Properties) == 0 {
				t.Fatal("schema must be closed and declare required properties")
			}
			for _, forbidden := range []string{"question", "prompt", "response", "response_text", "approved", "owner_action_evidence_ref", "run_id", "nonce"} {
				if _, exists := document.Properties[forbidden]; exists {
					t.Fatalf("schema must not accept raw content or client authority input %q", forbidden)
				}
			}
		})
	}
}

func TestDeliberationClientCodegenFixtureIsAuthoritySafe(t *testing.T) {
	raw, err := os.ReadFile("scripts/tests/fixtures/deliberation/client-codegen-projection.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fixture); err != nil {
		t.Fatalf("fixture must be valid JSON: %v", err)
	}
	for _, objectName := range []string{"session", "automation_grant", "synthesis_candidate"} {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(fixture[objectName], &object); err != nil {
			t.Fatalf("%s fixture must be an object: %v", objectName, err)
		}
		for _, forbidden := range []string{"question", "prompt", "response", "response_text", "approved", "owner_action_evidence_ref", "run_id", "nonce"} {
			if _, exists := object[forbidden]; exists {
				t.Fatalf("%s fixture must not teach client codegen a forbidden field %q", objectName, forbidden)
			}
		}
		if string(object["candidate_only"]) != "true" || string(object["non_formal"]) != "true" {
			t.Fatalf("%s fixture must remain candidate-only and non-formal", objectName)
		}
	}
}

func TestDeliberationSchemasKeepCandidateFlagsAndCurrentTurnHashGate(t *testing.T) {
	var grant map[string]any
	if err := json.Unmarshal(contracts.DeliberationAutomationGrantReadModelSchemaJSON, &grant); err != nil {
		t.Fatal(err)
	}
	properties := grant["properties"].(map[string]any)
	for _, field := range []string{"candidate_only", "non_formal"} {
		property := properties[field].(map[string]any)
		if value, ok := property["const"].(bool); !ok || !value {
			t.Fatalf("%s must be fixed true", field)
		}
	}
	decision, exists := properties["decision_ref"].(map[string]any)
	if !exists {
		t.Fatal("read-only grant projection must expose the Base-issued decision ref")
	}
	if value, ok := decision["readOnly"].(bool); !ok || !value {
		t.Fatal("client must not mint the Base-issued decision ref")
	}
	if _, exists := grant["allOf"]; !exists {
		t.Fatal("current-turn mode must require its question hash through a schema gate")
	}
	dispatchOnConfirm := properties["dispatch_on_confirm"].(map[string]any)
	if value, ok := dispatchOnConfirm["const"].(bool); !ok || !value {
		t.Fatal("confirmed automation grant must dispatch immediately without a second start action")
	}

	var synthesis map[string]any
	if err := json.Unmarshal(contracts.DeliberationSynthesisCandidateSchemaJSON, &synthesis); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"candidate_only", "non_formal"} {
		property := synthesis["properties"].(map[string]any)[field].(map[string]any)
		if value, ok := property["const"].(bool); !ok || !value {
			t.Fatalf("synthesis %s must be fixed true", field)
		}
	}
	receipt := synthesis["properties"].(map[string]any)["receipt_ref"].(map[string]any)
	if value, ok := receipt["readOnly"].(bool); !ok || !value {
		t.Fatal("client must not mint a ledger receipt ref")
	}
}
