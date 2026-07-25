package market

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ApprovedPackArtifactHandoffAttestation is the canonical, signed proof that
// a paired local host is handing one Owner-approved immutable Pack artifact to
// the cloud market. It is transport-neutral and never stores key material.
type ApprovedPackArtifactHandoffAttestation struct {
	HandoffID                  string   `json:"handoff_id"`
	IssuerHostRef              string   `json:"issuer_host_ref"`
	Audience                   string   `json:"audience"`
	CandidateRef               string   `json:"candidate_ref"`
	ApprovalReceiptRef         string   `json:"approval_receipt_ref"`
	ArtifactRef                string   `json:"artifact_ref"`
	PackRef                    string   `json:"pack_ref"`
	PackVersionRef             string   `json:"pack_version_ref"`
	ArtifactSHA256             string   `json:"artifact_sha256"`
	SanitizationManifestSHA256 string   `json:"sanitization_manifest_sha256"`
	EvidenceRefs               []string `json:"evidence_refs"`
	IssuedAt                   string   `json:"issued_at"`
	ExpiresAt                  string   `json:"expires_at"`
	Nonce                      string   `json:"nonce"`
	Signature                  string   `json:"signature"`
}

type approvedPackArtifactHandoffPayload struct {
	HandoffID                  string   `json:"handoff_id"`
	IssuerHostRef              string   `json:"issuer_host_ref"`
	Audience                   string   `json:"audience"`
	CandidateRef               string   `json:"candidate_ref"`
	ApprovalReceiptRef         string   `json:"approval_receipt_ref"`
	ArtifactRef                string   `json:"artifact_ref"`
	PackRef                    string   `json:"pack_ref"`
	PackVersionRef             string   `json:"pack_version_ref"`
	ArtifactSHA256             string   `json:"artifact_sha256"`
	SanitizationManifestSHA256 string   `json:"sanitization_manifest_sha256"`
	EvidenceRefs               []string `json:"evidence_refs"`
	IssuedAt                   string   `json:"issued_at"`
	ExpiresAt                  string   `json:"expires_at"`
	Nonce                      string   `json:"nonce"`
}

// ValidateApprovedPackArtifactHandoffAttestation checks signed-envelope shape.
// Cloud must additionally verify with the public key bound to the authenticated
// account and issuer_host_ref before accepting the handoff.
func ValidateApprovedPackArtifactHandoffAttestation(attestation ApprovedPackArtifactHandoffAttestation) error {
	if err := validateUnsignedApprovedPackArtifactHandoff(attestation); err != nil {
		return err
	}
	if strings.TrimSpace(attestation.Signature) == "" {
		return errors.New("approved pack artifact handoff attestation signature is required")
	}
	signature, err := base64.RawStdEncoding.DecodeString(attestation.Signature)
	if err != nil || len(signature) != ed25519.SignatureSize {
		return errors.New("approved pack artifact handoff attestation signature must be a base64 raw Ed25519 signature")
	}
	return nil
}

// CanonicalApprovedPackArtifactHandoffBytes returns exactly the JSON bytes
// signed by paired hosts. Signature is deliberately excluded.
func CanonicalApprovedPackArtifactHandoffBytes(attestation ApprovedPackArtifactHandoffAttestation) ([]byte, error) {
	if err := validateUnsignedApprovedPackArtifactHandoff(attestation); err != nil {
		return nil, err
	}
	return json.Marshal(approvedPackArtifactHandoffPayload{
		HandoffID: attestation.HandoffID, IssuerHostRef: attestation.IssuerHostRef, Audience: attestation.Audience,
		CandidateRef: attestation.CandidateRef, ApprovalReceiptRef: attestation.ApprovalReceiptRef, ArtifactRef: attestation.ArtifactRef,
		PackRef: attestation.PackRef, PackVersionRef: attestation.PackVersionRef, ArtifactSHA256: strings.ToLower(attestation.ArtifactSHA256), SanitizationManifestSHA256: strings.ToLower(attestation.SanitizationManifestSHA256),
		EvidenceRefs: append([]string(nil), attestation.EvidenceRefs...), IssuedAt: attestation.IssuedAt, ExpiresAt: attestation.ExpiresAt, Nonce: attestation.Nonce,
	})
}

// SignApprovedPackArtifactHandoff sets a raw-base64 Ed25519 signature. Key
// ownership and storage are intentionally outside the contracts repository.
func SignApprovedPackArtifactHandoff(attestation *ApprovedPackArtifactHandoffAttestation, privateKey ed25519.PrivateKey) error {
	if attestation == nil {
		return errors.New("approved pack artifact handoff attestation is required")
	}
	if len(privateKey) != ed25519.PrivateKeySize {
		return errors.New("paired host Ed25519 private key has invalid length")
	}
	payload, err := CanonicalApprovedPackArtifactHandoffBytes(*attestation)
	if err != nil {
		return err
	}
	attestation.Signature = base64.RawStdEncoding.EncodeToString(ed25519.Sign(privateKey, payload))
	return nil
}

// VerifyApprovedPackArtifactHandoff verifies canonical bytes. It does not
// establish account-to-host binding, which remains the cloud identity truth.
func VerifyApprovedPackArtifactHandoff(attestation ApprovedPackArtifactHandoffAttestation, publicKey ed25519.PublicKey) error {
	if err := ValidateApprovedPackArtifactHandoffAttestation(attestation); err != nil {
		return err
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return errors.New("paired host Ed25519 public key has invalid length")
	}
	payload, err := CanonicalApprovedPackArtifactHandoffBytes(attestation)
	if err != nil {
		return err
	}
	signature, _ := base64.RawStdEncoding.DecodeString(attestation.Signature)
	if !ed25519.Verify(publicKey, payload, signature) {
		return errors.New("approved pack artifact handoff attestation signature is invalid")
	}
	return nil
}

func validateUnsignedApprovedPackArtifactHandoff(attestation ApprovedPackArtifactHandoffAttestation) error {
	for name, value := range map[string]string{
		"handoff_id": attestation.HandoffID, "issuer_host_ref": attestation.IssuerHostRef, "audience": attestation.Audience,
		"candidate_ref": attestation.CandidateRef, "approval_receipt_ref": attestation.ApprovalReceiptRef, "artifact_ref": attestation.ArtifactRef,
		"pack_ref": attestation.PackRef, "pack_version_ref": attestation.PackVersionRef, "artifact_sha256": attestation.ArtifactSHA256,
		"sanitization_manifest_sha256": attestation.SanitizationManifestSHA256,
		"issued_at":                    attestation.IssuedAt, "expires_at": attestation.ExpiresAt, "nonce": attestation.Nonce,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("approved pack artifact handoff attestation %s is required", name)
		}
	}
	if err := validateAttestationTimes(attestation.IssuedAt, attestation.ExpiresAt); err != nil {
		return err
	}
	if _, err := decodeSHA256(attestation.ArtifactSHA256); err != nil {
		return fmt.Errorf("approved pack artifact handoff attestation artifact_sha256: %w", err)
	}
	if _, err := decodeSHA256(attestation.SanitizationManifestSHA256); err != nil {
		return fmt.Errorf("approved pack artifact handoff attestation sanitization_manifest_sha256: %w", err)
	}
	return validateUniqueHandoffEvidence(attestation.EvidenceRefs)
}

func validateAttestationTimes(issuedAtRaw, expiresAtRaw string) error {
	issuedAt, err := time.Parse(time.RFC3339, issuedAtRaw)
	if err != nil {
		return errors.New("approved pack artifact handoff attestation issued_at must be RFC3339")
	}
	expiresAt, err := time.Parse(time.RFC3339, expiresAtRaw)
	if err != nil {
		return errors.New("approved pack artifact handoff attestation expires_at must be RFC3339")
	}
	if !expiresAt.After(issuedAt) {
		return errors.New("approved pack artifact handoff attestation expires_at must be after issued_at")
	}
	return nil
}

func decodeSHA256(value string) ([]byte, error) {
	if len(value) != 64 {
		return nil, errors.New("must contain 64 hexadecimal characters")
	}
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != 32 {
		return nil, errors.New("must contain 64 hexadecimal characters")
	}
	return decoded, nil
}

func validateUniqueHandoffEvidence(refs []string) error {
	if len(refs) == 0 {
		return errors.New("approved pack artifact handoff attestation evidence_refs is required")
	}
	seen := make(map[string]struct{}, len(refs))
	for _, ref := range refs {
		ref = strings.TrimSpace(ref)
		if ref == "" {
			return errors.New("approved pack artifact handoff attestation evidence_refs contains an empty ref")
		}
		if _, exists := seen[ref]; exists {
			return fmt.Errorf("approved pack artifact handoff attestation evidence_refs contains duplicate ref %q", ref)
		}
		seen[ref] = struct{}{}
	}
	return nil
}
