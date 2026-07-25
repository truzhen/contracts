package contracts_test

import (
	"encoding/json"
	"testing"

	contracts "github.com/truzhen/contracts"
)

func TestTransactionOperationSchemasAreEmbeddedAndClosed(t *testing.T) {
	for name, raw := range map[string][]byte{
		"scene-runtime-plan-candidate":     contracts.SceneRuntimePlanCandidateSchemaJSON,
		"pack-hands-requirement-readmodel": contracts.PackHandsRequirementReadModelSchemaJSON,
		"experience-asset-candidate-readmodel": contracts.ExperienceAssetCandidateReadModelSchemaJSON,
		"approved-pack-artifact-handoff-attestation": contracts.ApprovedPackArtifactHandoffAttestationSchemaJSON,
	} {
		t.Run(name, func(t *testing.T) {
			var document struct {
				AdditionalProperties bool                       `json:"additionalProperties"`
				Properties           map[string]json.RawMessage `json:"properties"`
			}
			if err := json.Unmarshal(raw, &document); err != nil {
				t.Fatalf("schema must be valid JSON: %v", err)
			}
			if document.AdditionalProperties || len(document.Properties) == 0 {
				t.Fatal("schema must be closed and expose properties")
			}
		})
	}
}

func TestSceneRuntimeNodeCarriesAtMostOneProviderRequirementRef(t *testing.T) {
	var document struct {
		Defs map[string]struct {
			Properties map[string]json.RawMessage `json:"properties"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal(contracts.SceneRuntimePlanCandidateSchemaJSON, &document); err != nil {
		t.Fatal(err)
	}
	node := document.Defs["runtime_node"]
	property, exists := node.Properties["provider_requirement_refs"]
	if !exists {
		t.Fatal("runtime node must expose provider_requirement_refs")
	}
	var shape struct {
		Type     string `json:"type"`
		MinItems int    `json:"minItems"`
		MaxItems int    `json:"maxItems"`
		Unique   bool   `json:"uniqueItems"`
	}
	if err := json.Unmarshal(property, &shape); err != nil {
		t.Fatal(err)
	}
	if shape.Type != "array" || shape.MinItems != 1 || shape.MaxItems != 1 || !shape.Unique {
		t.Fatalf("provider requirement refs must be a unique single-ref array: %+v", shape)
	}
}
