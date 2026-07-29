package readmodels_test

import (
	"encoding/json"
	"testing"

	contracts "github.com/truzhen/contracts"
	"github.com/truzhen/contracts/readmodels"
)

func validMobilePairingBootstrapRequest() readmodels.MobilePairingBootstrapRequest {
	return readmodels.MobilePairingBootstrapRequest{
		DeviceLabel:    "Owner phone",
		Platform:       readmodels.MobilePlatformAndroid,
		AppInstanceRef: "app_instance://mobile/owner-phone",
		IdempotencyKey: "mobile-bootstrap-001",
	}
}

func validMobilePairingBootstrapCandidate() readmodels.MobilePairingBootstrapCandidate {
	return readmodels.MobilePairingBootstrapCandidate{
		OK:                    true,
		CandidateCreated:      true,
		Duplicate:             false,
		CandidateOnly:         true,
		MobileTruthSource:     false,
		Status:                readmodels.MobilePairingStatusPendingPCApproval,
		CandidateRef:          "candidate://mobile-pairing/001",
		CandidateKind:         readmodels.MobilePairingCandidateKind,
		DeviceLabel:           "Owner phone",
		Platform:              readmodels.MobilePlatformAndroid,
		IdempotencyKey:        "mobile-bootstrap-001",
		CreatedAt:             "2026-07-29T00:00:00Z",
		ProducesOwnerDecision: false,
		CredentialState:       readmodels.MobilePairingCredentialStateNotExposed,
	}
}

func TestMobilePairingStrictDecodeAndValidateRoundTrips(t *testing.T) {
	request := validMobilePairingBootstrapRequest()
	requestRaw, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	decodedRequest, err := readmodels.DecodeMobilePairingBootstrapRequest(requestRaw)
	if err != nil {
		t.Fatalf("decode request: %v", err)
	}
	if decodedRequest != request {
		t.Fatalf("request changed during strict decode: %#v", decodedRequest)
	}
	if err := readmodels.ValidateMobilePairingBootstrapRequest(decodedRequest); err != nil {
		t.Fatalf("validate request: %v", err)
	}

	candidate := validMobilePairingBootstrapCandidate()
	candidateRaw, err := json.Marshal(candidate)
	if err != nil {
		t.Fatal(err)
	}
	decodedCandidate, err := readmodels.DecodeMobilePairingBootstrapCandidate(candidateRaw)
	if err != nil {
		t.Fatalf("decode candidate: %v", err)
	}
	if decodedCandidate != candidate {
		t.Fatalf("candidate changed during strict decode: %#v", decodedCandidate)
	}
	if err := readmodels.ValidateMobilePairingBootstrapCandidate(decodedCandidate); err != nil {
		t.Fatalf("validate candidate: %v", err)
	}

	issue := readmodels.MobileSessionIssueIntent{CandidateRef: candidate.CandidateRef, IdempotencyKey: "mobile-session-001"}
	issueRaw, err := json.Marshal(issue)
	if err != nil {
		t.Fatal(err)
	}
	decodedIssue, err := readmodels.DecodeMobileSessionIssueIntent(issueRaw)
	if err != nil {
		t.Fatalf("decode issue intent: %v", err)
	}
	if decodedIssue != issue {
		t.Fatalf("issue intent changed during strict decode: %#v", decodedIssue)
	}
	if err := readmodels.ValidateMobileSessionIssueIntent(decodedIssue); err != nil {
		t.Fatalf("validate issue intent: %v", err)
	}
}

func TestMobilePairingStrictDecodeRejectsAuthoritySmuggling(t *testing.T) {
	requestRaw := []byte(`{"device_label":"Owner phone","platform":"android","app_instance_ref":"app_instance://mobile/owner-phone","idempotency_key":"mobile-bootstrap-001","owner_decision_ref":"decision://forged"}`)
	if _, err := readmodels.DecodeMobilePairingBootstrapRequest(requestRaw); err == nil {
		t.Fatal("request carrying an owner decision ref was accepted")
	}

	candidateRaw := []byte(`{"ok":true,"candidate_created":true,"duplicate":false,"candidate_only":true,"mobile_truth_source":false,"status":"pending_pc_approval","candidate_ref":"candidate://mobile-pairing/001","candidate_kind":"mobile_pairing_bootstrap_candidate","device_label":"Owner phone","platform":"android","idempotency_key":"mobile-bootstrap-001","created_at":"2026-07-29T00:00:00Z","produces_owner_decision":false,"credential_state":"not_exposed_to_mobile","receipt_ref":"receipt://forged"}`)
	if _, err := readmodels.DecodeMobilePairingBootstrapCandidate(candidateRaw); err == nil {
		t.Fatal("candidate carrying a receipt ref was accepted")
	}

	issueRaw := []byte(`{"candidate_ref":"candidate://mobile-pairing/001","idempotency_key":"mobile-session-001","bootstrap_proof":"bearer"}`)
	if _, err := readmodels.DecodeMobileSessionIssueIntent(issueRaw); err == nil {
		t.Fatal("session issue body carrying a bootstrap proof was accepted")
	}
}

func TestMobilePairingStrictDecodeRejectsMissingRequiredFalseAndTrailingValues(t *testing.T) {
	missingFalse := []byte(`{"ok":true,"candidate_created":true,"duplicate":false,"candidate_only":true,"status":"pending_pc_approval","candidate_ref":"candidate://mobile-pairing/001","candidate_kind":"mobile_pairing_bootstrap_candidate","device_label":"Owner phone","platform":"android","idempotency_key":"mobile-bootstrap-001","created_at":"2026-07-29T00:00:00Z","produces_owner_decision":false,"credential_state":"not_exposed_to_mobile"}`)
	if _, err := readmodels.DecodeMobilePairingBootstrapCandidate(missingFalse); err == nil {
		t.Fatal("candidate missing required false mobile_truth_source was accepted")
	}

	trailing := []byte(`{"candidate_ref":"candidate://mobile-pairing/001","idempotency_key":"mobile-session-001"} {}`)
	if _, err := readmodels.DecodeMobileSessionIssueIntent(trailing); err == nil {
		t.Fatal("multiple top-level JSON values were accepted")
	}
}

func TestValidateMobilePairingRejectsFormalLookingValues(t *testing.T) {
	candidate := validMobilePairingBootstrapCandidate()
	candidate.CandidateOnly = false
	if err := readmodels.ValidateMobilePairingBootstrapCandidate(candidate); err == nil {
		t.Fatal("non-candidate mobile pairing record was accepted")
	}
	candidate = validMobilePairingBootstrapCandidate()
	candidate.MobileTruthSource = true
	if err := readmodels.ValidateMobilePairingBootstrapCandidate(candidate); err == nil {
		t.Fatal("mobile truth source was accepted")
	}
	candidate = validMobilePairingBootstrapCandidate()
	candidate.ProducesOwnerDecision = true
	if err := readmodels.ValidateMobilePairingBootstrapCandidate(candidate); err == nil {
		t.Fatal("owner-decision-producing candidate was accepted")
	}
}

func TestMobilePairingSchemasRemainClosedAndAuthorityFree(t *testing.T) {
	for name, raw := range map[string][]byte{
		"bootstrap request":   contracts.MobilePairingBootstrapRequestSchemaJSON,
		"bootstrap candidate": contracts.MobilePairingBootstrapCandidateSchemaJSON,
		"session issue":       contracts.MobileSessionIssueIntentSchemaJSON,
	} {
		var document struct {
			AdditionalProperties bool                       `json:"additionalProperties"`
			Required             []string                   `json:"required"`
			Properties           map[string]json.RawMessage `json:"properties"`
		}
		if err := json.Unmarshal(raw, &document); err != nil {
			t.Fatalf("%s schema must be JSON: %v", name, err)
		}
		if document.AdditionalProperties || len(document.Required) == 0 || len(document.Properties) == 0 {
			t.Fatalf("%s schema must be closed and required", name)
		}
		for _, forbidden := range []string{"owner_decision_ref", "receipt_ref", "bootstrap_proof", "access_token", "raw_credential"} {
			if _, exists := document.Properties[forbidden]; exists {
				t.Fatalf("%s schema exposes forbidden field %q", name, forbidden)
			}
		}
	}
}

func TestMobilePairingConstantsMatchCanonicalSchema(t *testing.T) {
	var request struct {
		Properties map[string]struct {
			Enum []string `json:"enum"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(contracts.MobilePairingBootstrapRequestSchemaJSON, &request); err != nil {
		t.Fatal(err)
	}
	platforms := request.Properties["platform"].Enum
	if len(platforms) != 2 || platforms[0] != readmodels.MobilePlatformAndroid || platforms[1] != readmodels.MobilePlatformIOS {
		t.Fatalf("platform constants diverge from schema: %#v", platforms)
	}

	var candidate struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(contracts.MobilePairingBootstrapCandidateSchemaJSON, &candidate); err != nil {
		t.Fatal(err)
	}
	var statusProperty struct {
		Enum []string `json:"enum"`
	}
	if err := json.Unmarshal(candidate.Properties["status"], &statusProperty); err != nil {
		t.Fatal(err)
	}
	statuses := statusProperty.Enum
	if len(statuses) != 2 || statuses[0] != readmodels.MobilePairingStatusPendingPCApproval || statuses[1] != readmodels.MobilePairingStatusPCApproved {
		t.Fatalf("status constants diverge from schema: %#v", statuses)
	}
	var candidateKind struct {
		Const string `json:"const"`
	}
	if err := json.Unmarshal(candidate.Properties["candidate_kind"], &candidateKind); err != nil {
		t.Fatal(err)
	}
	if got := candidateKind.Const; got != readmodels.MobilePairingCandidateKind {
		t.Fatalf("candidate kind constant diverges from schema: %q", got)
	}
	var credentialState struct {
		Const string `json:"const"`
	}
	if err := json.Unmarshal(candidate.Properties["credential_state"], &credentialState); err != nil {
		t.Fatal(err)
	}
	if got := credentialState.Const; got != readmodels.MobilePairingCredentialStateNotExposed {
		t.Fatalf("credential state constant diverges from schema: %q", got)
	}
}
