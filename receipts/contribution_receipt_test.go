package receipts_test

import (
	"encoding/json"
	"testing"

	"github.com/truzhen/contracts/receipts"
)

func TestContributionReceiptMarshalCanonicalRefs(t *testing.T) {
	receipt := receipts.ContributionReceipt{
		ReceiptRef:               "receipt://contribution/1",
		Kind:                     receipts.ContributionReceiptKindAuthorAdoption,
		ContributionCandidateRef: "contribution_candidate://pack/env/1",
		SourcePackVersionRef:     "pack_version://environmental-enforcement/1.0.0",
		ResultPackVersionRef:     "pack_version://environmental-enforcement/1.1.0",
		CreatedAt:                "2026-07-07T00:00:00Z",
	}
	raw, err := json.Marshal(receipt)
	if err != nil {
		t.Fatalf("marshal contribution receipt: %v", err)
	}
	for _, want := range []string{
		`"receipt_ref":"receipt://contribution/1"`,
		`"kind":"author_adoption"`,
		`"contribution_candidate_ref":"contribution_candidate://pack/env/1"`,
		`"result_pack_version_ref":"pack_version://environmental-enforcement/1.1.0"`,
	} {
		if !json.Valid(raw) || !contains(raw, want) {
			t.Fatalf("contribution receipt JSON missing %s: %s", want, raw)
		}
	}
}

func contains(raw []byte, want string) bool {
	s := string(raw)
	for i := 0; i+len(want) <= len(s); i++ {
		if s[i:i+len(want)] == want {
			return true
		}
	}
	return false
}
