package readmodels

import (
	"errors"
	"fmt"
	"strings"

	"github.com/truzhen/contracts/candidates"
	"github.com/truzhen/contracts/receipts"
)

// ExperienceAssetCandidateReadModel is the client-safe projection of one
// desensitized experience candidate. It cannot carry raw Receipt payloads,
// Owner authority, or a market-upload proof.
type ExperienceAssetCandidateReadModel struct {
	CandidateID                    string                                     `json:"candidate_id"`
	CandidateOnly                  bool                                       `json:"candidate_only"`
	NonFormal                      bool                                       `json:"non_formal"`
	Status                         string                                     `json:"status"`
	Version                        int                                        `json:"version"`
	SourceReceiptRef               string                                     `json:"source_receipt_ref"`
	TransactionRef                 string                                     `json:"transaction_ref"`
	PackRef                        string                                     `json:"pack_ref"`
	SourcePackVersionRef           string                                     `json:"source_pack_version_ref"`
	SourceModule                   string                                     `json:"source_module"`
	Provenance                     ExperienceAssetProvenance                  `json:"provenance"`
	InheritedPermissions           string                                     `json:"inherited_permissions"`
	DesensitizedExperience         string                                     `json:"desensitized_experience,omitempty"`
	DesensitizationReport          ExperienceAssetDesensitizationReport       `json:"desensitization_report"`
	ContributionCandidate          *candidates.PackUsageContributionCandidate `json:"contribution_candidate,omitempty"`
	PackVersionMigrationCandidate  *candidates.PackVersionMigrationCandidate  `json:"pack_version_migration_candidate,omitempty"`
	ContributionReceipt            *receipts.ContributionReceipt              `json:"contribution_receipt,omitempty"`
	PackVersionMigrationReceiptRef string                                     `json:"pack_version_migration_receipt_ref,omitempty"`
	CreatedAt                      string                                     `json:"created_at"`
}

type ExperienceAssetProvenance struct {
	ReceiptRef     string `json:"receipt_ref"`
	TransactionRef string `json:"transaction_ref"`
	TaskRef        string `json:"task_ref,omitempty"`
	SourceModule   string `json:"source_module"`
	SourceHash     string `json:"source_hash"`
}

type ExperienceAssetDesensitizationReport struct {
	OriginalLength int      `json:"original_length"`
	RedactedCount  int      `json:"redacted_count"`
	FieldsMasked   []string `json:"fields_masked"`
	Rejected       bool     `json:"rejected"`
	Reason         string   `json:"reason,omitempty"`
}

// ValidateExperienceAssetCandidateReadModel checks candidate/formal isolation
// and immutable provenance without treating the projection as a truth source.
func ValidateExperienceAssetCandidateReadModel(model ExperienceAssetCandidateReadModel) error {
	for name, value := range map[string]string{
		"candidate_id": model.CandidateID, "status": model.Status, "source_receipt_ref": model.SourceReceiptRef,
		"transaction_ref": model.TransactionRef, "pack_ref": model.PackRef, "source_pack_version_ref": model.SourcePackVersionRef,
		"source_module": model.SourceModule, "inherited_permissions": model.InheritedPermissions, "created_at": model.CreatedAt,
		"provenance.receipt_ref": model.Provenance.ReceiptRef, "provenance.transaction_ref": model.Provenance.TransactionRef,
		"provenance.source_module": model.Provenance.SourceModule, "provenance.source_hash": model.Provenance.SourceHash,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("experience asset read model %s is required", name)
		}
	}
	if !model.CandidateOnly || !model.NonFormal {
		return errors.New("experience asset read model must remain candidate_only and non_formal")
	}
	if model.Version < 1 || model.DesensitizationReport.OriginalLength < 0 || model.DesensitizationReport.RedactedCount < 0 {
		return errors.New("experience asset read model contains invalid numeric fields")
	}
	if model.Provenance.ReceiptRef != model.SourceReceiptRef || model.Provenance.TransactionRef != model.TransactionRef || model.Provenance.SourceModule != model.SourceModule {
		return errors.New("experience asset read model provenance must bind source receipt, transaction, and module")
	}
	if !validExperienceAssetStatus(model.Status) {
		return fmt.Errorf("experience asset read model has unknown status %q", model.Status)
	}
	if err := validateUniqueRefs("desensitization_report.fields_masked", model.DesensitizationReport.FieldsMasked); err != nil {
		return err
	}
	if model.ContributionReceipt != nil && model.ContributionCandidate == nil {
		return errors.New("experience asset contribution receipt requires a contribution candidate")
	}
	return nil
}

func validExperienceAssetStatus(status string) bool {
	switch status {
	case "EXTRACTED", "DESENSITIZED", "PENDING_REVIEW", "OWNER_APPROVED", "REJECTED_SENSITIVE":
		return true
	default:
		return false
	}
}
