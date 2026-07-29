package readmodels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
	"unicode/utf8"
)

const (
	MobilePlatformAndroid = "android"
	MobilePlatformIOS     = "ios"

	MobilePairingStatusPendingPCApproval = "pending_pc_approval"
	MobilePairingStatusPCApproved        = "pc_approved"

	MobilePairingCandidateKind             = "mobile_pairing_bootstrap_candidate"
	MobilePairingCredentialStateNotExposed = "not_exposed_to_mobile"
)

// MobilePairingBootstrapRequest is an authority-free description supplied by a
// phone before it holds a mobile session. The Host, not the phone, creates all
// binding and pairing refs after the PC owner approves the candidate.
type MobilePairingBootstrapRequest struct {
	DeviceLabel    string `json:"device_label"`
	Platform       string `json:"platform"`
	AppInstanceRef string `json:"app_instance_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

// MobilePairingBootstrapCandidate is a pre-approval ReadModel/Candidate
// projection. It carries no raw bootstrap proof, bearer, formal identity,
// binding, pairing, OwnerDecision, or Receipt truth. The Host creates formal
// binding and pairing facts only after the PC owner approves this candidate.
type MobilePairingBootstrapCandidate struct {
	OK                    bool   `json:"ok"`
	CandidateCreated      bool   `json:"candidate_created"`
	Duplicate             bool   `json:"duplicate"`
	CandidateOnly         bool   `json:"candidate_only"`
	MobileTruthSource     bool   `json:"mobile_truth_source"`
	Status                string `json:"status"`
	CandidateRef          string `json:"candidate_ref"`
	CandidateKind         string `json:"candidate_kind"`
	DeviceLabel           string `json:"device_label"`
	Platform              string `json:"platform"`
	IdempotencyKey        string `json:"idempotency_key"`
	CreatedAt             string `json:"created_at"`
	ProducesOwnerDecision bool   `json:"produces_owner_decision"`
	CredentialState       string `json:"credential_state"`
}

// MobileSessionIssueIntent is the JSON body for a post-approval session issue.
// The one-time bootstrap proof travels only in a controlled request header and
// is intentionally not represented here.
type MobileSessionIssueIntent struct {
	CandidateRef   string `json:"candidate_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

// DecodeMobilePairingBootstrapRequest accepts only the closed request schema.
// The mobile caller cannot smuggle formal references or approval facts through
// fields that the schema does not own.
func DecodeMobilePairingBootstrapRequest(raw []byte) (MobilePairingBootstrapRequest, error) {
	var request MobilePairingBootstrapRequest
	if err := decodeStrictMobileJSON(raw, &request, "device_label", "platform", "app_instance_ref", "idempotency_key"); err != nil {
		return MobilePairingBootstrapRequest{}, err
	}
	return request, ValidateMobilePairingBootstrapRequest(request)
}

// ValidateMobilePairingBootstrapRequest validates the schema-level shape only.
// It does not issue a pairing fact, an OwnerDecision, a session, or a Receipt.
func ValidateMobilePairingBootstrapRequest(request MobilePairingBootstrapRequest) error {
	if err := validateMobileString("device_label", request.DeviceLabel, 128); err != nil {
		return err
	}
	if !isMobilePlatform(request.Platform) {
		return fmt.Errorf("platform must be %q or %q", MobilePlatformAndroid, MobilePlatformIOS)
	}
	if err := validateMobileString("app_instance_ref", request.AppInstanceRef, 256); err != nil {
		return err
	}
	return validateMobileString("idempotency_key", request.IdempotencyKey, 128)
}

// DecodeMobilePairingBootstrapCandidate accepts only the closed candidate
// projection. Formal facts remain Host-owned until a separate Base Gate path
// has approved them.
func DecodeMobilePairingBootstrapCandidate(raw []byte) (MobilePairingBootstrapCandidate, error) {
	var candidate MobilePairingBootstrapCandidate
	if err := decodeStrictMobileJSON(raw, &candidate,
		"ok", "candidate_created", "duplicate", "candidate_only", "mobile_truth_source",
		"status", "candidate_ref", "candidate_kind", "device_label", "platform",
		"idempotency_key", "created_at", "produces_owner_decision", "credential_state"); err != nil {
		return MobilePairingBootstrapCandidate{}, err
	}
	return candidate, ValidateMobilePairingBootstrapCandidate(candidate)
}

// ValidateMobilePairingBootstrapCandidate ensures that this contract cannot be
// presented as a formal mobile truth source, approval, or credential container.
func ValidateMobilePairingBootstrapCandidate(candidate MobilePairingBootstrapCandidate) error {
	if !candidate.CandidateOnly {
		return fmt.Errorf("candidate_only must be true")
	}
	if candidate.MobileTruthSource {
		return fmt.Errorf("mobile_truth_source must be false")
	}
	if candidate.ProducesOwnerDecision {
		return fmt.Errorf("produces_owner_decision must be false")
	}
	if candidate.Status != MobilePairingStatusPendingPCApproval && candidate.Status != MobilePairingStatusPCApproved {
		return fmt.Errorf("unknown mobile pairing status %q", candidate.Status)
	}
	if err := validateMobileString("candidate_ref", candidate.CandidateRef, 0); err != nil {
		return err
	}
	if candidate.CandidateKind != MobilePairingCandidateKind {
		return fmt.Errorf("candidate_kind must be %q", MobilePairingCandidateKind)
	}
	if err := validateMobileString("device_label", candidate.DeviceLabel, 128); err != nil {
		return err
	}
	if !isMobilePlatform(candidate.Platform) {
		return fmt.Errorf("platform must be %q or %q", MobilePlatformAndroid, MobilePlatformIOS)
	}
	if err := validateMobileString("idempotency_key", candidate.IdempotencyKey, 0); err != nil {
		return err
	}
	if _, err := time.Parse(time.RFC3339, candidate.CreatedAt); err != nil {
		return fmt.Errorf("created_at must be RFC3339: %w", err)
	}
	if candidate.CredentialState != MobilePairingCredentialStateNotExposed {
		return fmt.Errorf("credential_state must be %q", MobilePairingCredentialStateNotExposed)
	}
	return nil
}

// DecodeMobileSessionIssueIntent accepts the header-free body only. The
// one-time bootstrap proof remains in the controlled request header and cannot
// be replayed through JSON.
func DecodeMobileSessionIssueIntent(raw []byte) (MobileSessionIssueIntent, error) {
	var intent MobileSessionIssueIntent
	if err := decodeStrictMobileJSON(raw, &intent, "candidate_ref", "idempotency_key"); err != nil {
		return MobileSessionIssueIntent{}, err
	}
	return intent, ValidateMobileSessionIssueIntent(intent)
}

// ValidateMobileSessionIssueIntent validates the body shape without issuing a
// session or accepting a bootstrap proof.
func ValidateMobileSessionIssueIntent(intent MobileSessionIssueIntent) error {
	if err := validateMobileString("candidate_ref", intent.CandidateRef, 0); err != nil {
		return err
	}
	return validateMobileString("idempotency_key", intent.IdempotencyKey, 128)
}

func decodeStrictMobileJSON(raw []byte, target any, required ...string) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode mobile pairing JSON: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode mobile pairing JSON: multiple top-level values")
		}
		return fmt.Errorf("decode mobile pairing JSON trailing data: %w", err)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return fmt.Errorf("decode mobile pairing JSON object: %w", err)
	}
	if fields == nil {
		return fmt.Errorf("mobile pairing JSON must be an object")
	}
	for _, name := range required {
		if _, ok := fields[name]; !ok {
			return fmt.Errorf("mobile pairing JSON missing required field %q", name)
		}
	}
	return nil
}

func validateMobileString(name, value string, maximumRunes int) error {
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	if !utf8.ValidString(value) {
		return fmt.Errorf("%s must be valid UTF-8", name)
	}
	if maximumRunes > 0 && utf8.RuneCountInString(value) > maximumRunes {
		return fmt.Errorf("%s exceeds %d characters", name, maximumRunes)
	}
	return nil
}

func isMobilePlatform(platform string) bool {
	return platform == MobilePlatformAndroid || platform == MobilePlatformIOS
}
