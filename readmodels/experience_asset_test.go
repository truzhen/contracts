package readmodels_test

import (
	"encoding/json"
	"testing"

	contracts "github.com/truzhen/contracts"
	"github.com/truzhen/contracts/readmodels"
)

func validExperienceAssetReadModel() readmodels.ExperienceAssetCandidateReadModel {
	return readmodels.ExperienceAssetCandidateReadModel{
		CandidateID: "experience_asset_candidate://18/001", CandidateOnly: true, NonFormal: true, Status: "PENDING_REVIEW", Version: 1,
		SourceReceiptRef: "receipt://03/001", TransactionRef: "transaction://smart-home/001", PackRef: "scene-pack://smart-home", SourcePackVersionRef: "scene-pack://smart-home@1.0.0", SourceModule: "11-execution-gateway",
		Provenance:           readmodels.ExperienceAssetProvenance{ReceiptRef: "receipt://03/001", TransactionRef: "transaction://smart-home/001", SourceModule: "11-execution-gateway", SourceHash: "sha256:001"},
		InheritedPermissions: "owner_only", DesensitizedExperience: "客户要求已脱敏", DesensitizationReport: readmodels.ExperienceAssetDesensitizationReport{OriginalLength: 12, RedactedCount: 1, FieldsMasked: []string{"phone"}}, CreatedAt: "2026-07-25T10:00:00Z",
	}
}

func TestValidateExperienceAssetCandidateReadModel(t *testing.T) {
	model := validExperienceAssetReadModel()
	if err := readmodels.ValidateExperienceAssetCandidateReadModel(model); err != nil {
		t.Fatal(err)
	}
	model.CandidateOnly = false
	if err := readmodels.ValidateExperienceAssetCandidateReadModel(model); err == nil {
		t.Fatal("formal-looking experience read model accepted")
	}
}

func TestExperienceAssetCandidateReadModelSchemaIsClosedAndSafe(t *testing.T) {
	var document struct {
		AdditionalProperties bool                       `json:"additionalProperties"`
		Properties           map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(contracts.ExperienceAssetCandidateReadModelSchemaJSON, &document); err != nil {
		t.Fatal(err)
	}
	if document.AdditionalProperties || len(document.Properties) == 0 {
		t.Fatal("experience asset schema must be closed")
	}
	for _, forbidden := range []string{"raw_payload", "owner_decision", "approval_signature", "access_token"} {
		if _, exists := document.Properties[forbidden]; exists {
			t.Fatalf("experience asset schema exposes forbidden field %q", forbidden)
		}
	}
}
