package candidates_test

import (
	"encoding/json"
	"testing"

	"github.com/truzhen/contracts/candidates"
)

func TestContributionCandidatesMarshalCanonicalRefs(t *testing.T) {
	usage := candidates.PackUsageContributionCandidate{
		CandidateID:            "contribution_candidate://pack/env/1",
		SourceTransactionRef:   "transaction://case/1",
		SourcePackRef:          "pack://environmental-enforcement",
		SourcePackVersionRef:   "pack_version://environmental-enforcement/1.0.0",
		SanitizationReceiptRef: "receipt://sanitize/1",
		OwnerDecisionRef:       "owner_decision://base/1",
		SummaryRedacted:        "脱敏后的复盘摘要",
		Status:                 candidates.ContributionStatusOwnerApproved,
		CandidateOnly:          true,
		NonFormal:              true,
		CreatedAt:              "2026-07-07T00:00:00Z",
	}
	raw, err := json.Marshal(usage)
	if err != nil {
		t.Fatalf("marshal usage contribution: %v", err)
	}
	for _, want := range []string{
		`"candidate_only":true`,
		`"non_formal":true`,
		`"owner_decision_ref":"owner_decision://base/1"`,
		`"sanitization_receipt_ref":"receipt://sanitize/1"`,
	} {
		if !json.Valid(raw) || !containsJSON(raw, want) {
			t.Fatalf("usage contribution JSON missing %s: %s", want, raw)
		}
	}

	migration := candidates.PackVersionMigrationCandidate{
		CandidateID:         "migration_candidate://pack/env/1",
		TransactionRef:      "transaction://case/1",
		PackRef:             "pack://environmental-enforcement",
		FromPackVersionRef:  "pack_version://environmental-enforcement/1.0.0",
		ToPackVersionRef:    "pack_version://environmental-enforcement/1.1.0",
		RiskColor:           candidates.MigrationRiskYellow,
		RequiredOwnerAction: candidates.MigrationOwnerConfirm,
		Status:              candidates.MigrationStatusProposed,
		CandidateOnly:       true,
		NonFormal:           true,
		CreatedAt:           "2026-07-07T00:00:00Z",
	}
	raw, err = json.Marshal(migration)
	if err != nil {
		t.Fatalf("marshal migration candidate: %v", err)
	}
	for _, want := range []string{
		`"from_pack_version_ref":"pack_version://environmental-enforcement/1.0.0"`,
		`"to_pack_version_ref":"pack_version://environmental-enforcement/1.1.0"`,
		`"required_owner_action":"owner_confirm"`,
	} {
		if !json.Valid(raw) || !containsJSON(raw, want) {
			t.Fatalf("migration candidate JSON missing %s: %s", want, raw)
		}
	}
}

func containsJSON(raw []byte, want string) bool {
	return json.Valid(raw) && stringContains(string(raw), want)
}

func stringContains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
