package readmodels_test

import (
	"encoding/json"
	"testing"

	contracts "github.com/truzhen/contracts"
	"github.com/truzhen/contracts/readmodels"
)

func validPackHandsRequirementReadModel() readmodels.PackHandsRequirementReadModel {
	return readmodels.PackHandsRequirementReadModel{
		PackRef:                      "scene_pack://smart-home-owner",
		PackVersionRef:               "scene_pack://smart-home-owner@1.0.0",
		RequirementRef:               "provider_requirement://frappe-read",
		RequiredCapabilities:         []string{"project.customer", "project.milestone"},
		ProviderFamily:               "frappe",
		SoftwareRequirementRefs:      []string{"frappe-suite-runtime", "frappe-mcp-runtime"},
		SoftwareResolutionLockRefs:   []string{"software_lock://frappe-suite", "software_lock://frappe-mcp"},
		ResolutionStatus:             "bound",
		ResolvedProviderResourceRef:  "provider_resource://frappe-local",
		ResolvedProviderDisplayLabel: "本地 Frappe",
		ResolvedControlMethodRef:     "control_method://frappe-mcp-read",
		ResolvedExecutionMode:        "mcp",
		Installable:                  false,
		RiskClass:                    "high",
		Optional:                     false,
		LastProbeReceiptRef:          "receipt://probe-001",
		LastResolutionReceiptRef:     "receipt://resolution-001",
		Stale:                        false,
		ExpiresAt:                    "2026-07-24T01:00:00Z",
	}
}

func TestValidatePackHandsRequirementReadModel(t *testing.T) {
	model := validPackHandsRequirementReadModel()
	if err := readmodels.ValidatePackHandsRequirementReadModel(model); err != nil {
		t.Fatalf("valid read model rejected: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*readmodels.PackHandsRequirementReadModel)
	}{
		{name: "empty software ref", mutate: func(m *readmodels.PackHandsRequirementReadModel) {
			m.SoftwareRequirementRefs = []string{""}
		}},
		{name: "duplicate lock ref", mutate: func(m *readmodels.PackHandsRequirementReadModel) {
			m.SoftwareResolutionLockRefs = []string{"software_lock://same", "software_lock://same"}
		}},
		{name: "unknown resolution status", mutate: func(m *readmodels.PackHandsRequirementReadModel) {
			m.ResolutionStatus = "ready_to_execute"
		}},
		{name: "installable without candidate", mutate: func(m *readmodels.PackHandsRequirementReadModel) {
			m.Installable = true
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidate := validPackHandsRequirementReadModel()
			tt.mutate(&candidate)
			if err := readmodels.ValidatePackHandsRequirementReadModel(candidate); err == nil {
				t.Fatalf("expected %s to be rejected", tt.name)
			}
		})
	}
}

func TestPackHandsRequirementReadModelSchemaIsClosedAndAuthoritySafe(t *testing.T) {
	var document struct {
		AdditionalProperties bool                       `json:"additionalProperties"`
		Required             []string                   `json:"required"`
		Properties           map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(contracts.PackHandsRequirementReadModelSchemaJSON, &document); err != nil {
		t.Fatalf("schema must be valid JSON: %v", err)
	}
	if document.AdditionalProperties || len(document.Required) == 0 || len(document.Properties) == 0 {
		t.Fatal("schema must be closed and declare required properties")
	}
	for _, forbidden := range []string{"raw_endpoint", "raw_credential", "access_token", "secret", "ready_to_execute", "selected_provider_resource_ref"} {
		if _, exists := document.Properties[forbidden]; exists {
			t.Fatalf("schema must not expose forbidden field %q", forbidden)
		}
	}
}
