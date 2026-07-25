package market

import (
	"bytes"
	"crypto/ed25519"
	"testing"
)

func validApprovedPackArtifactHandoffAttestation() ApprovedPackArtifactHandoffAttestation {
	return ApprovedPackArtifactHandoffAttestation{
		HandoffID: "handoff://p13/001", IssuerHostRef: "paired-host://author/001", Audience: "cloud-market", CandidateRef: "candidate://18/001",
		ApprovalReceiptRef: "receipt://03/approval-001", ArtifactRef: "artifact://14/001", PackRef: "scene-pack://smart-home", PackVersionRef: "scene-pack://smart-home@1.2.0",
		ArtifactSHA256: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", SanitizationManifestSHA256: "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
		EvidenceRefs: []string{"receipt://03/approval-001", "receipt://03/source-001"},
		IssuedAt:     "2026-07-25T10:00:00Z", ExpiresAt: "2026-07-25T10:10:00Z", Nonce: "nonce-p13-001",
	}
}

func TestApprovedPackArtifactHandoffAttestationCanonicalSignature(t *testing.T) {
	privateKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{7}, ed25519.SeedSize))
	attestation := validApprovedPackArtifactHandoffAttestation()
	if err := SignApprovedPackArtifactHandoff(&attestation, privateKey); err != nil {
		t.Fatal(err)
	}
	if err := VerifyApprovedPackArtifactHandoff(attestation, privateKey.Public().(ed25519.PublicKey)); err != nil {
		t.Fatalf("signed attestation rejected: %v", err)
	}
	canonical, err := CanonicalApprovedPackArtifactHandoffBytes(attestation)
	if err != nil {
		t.Fatal(err)
	}
	const golden = `{"handoff_id":"handoff://p13/001","issuer_host_ref":"paired-host://author/001","audience":"cloud-market","candidate_ref":"candidate://18/001","approval_receipt_ref":"receipt://03/approval-001","artifact_ref":"artifact://14/001","pack_ref":"scene-pack://smart-home","pack_version_ref":"scene-pack://smart-home@1.2.0","artifact_sha256":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","sanitization_manifest_sha256":"abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789","evidence_refs":["receipt://03/approval-001","receipt://03/source-001"],"issued_at":"2026-07-25T10:00:00Z","expires_at":"2026-07-25T10:10:00Z","nonce":"nonce-p13-001"}`
	if string(canonical) != golden {
		t.Fatalf("canonical bytes drift\n got: %s\nwant: %s", canonical, golden)
	}
}

func TestApprovedPackArtifactHandoffAttestationRejectsTampering(t *testing.T) {
	privateKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{8}, ed25519.SeedSize))
	attestation := validApprovedPackArtifactHandoffAttestation()
	if err := SignApprovedPackArtifactHandoff(&attestation, privateKey); err != nil {
		t.Fatal(err)
	}
	for name, mutate := range map[string]func(*ApprovedPackArtifactHandoffAttestation){
		"pack": func(a *ApprovedPackArtifactHandoffAttestation) { a.PackRef = "scene-pack://other" },
		"digest": func(a *ApprovedPackArtifactHandoffAttestation) {
			a.ArtifactSHA256 = "abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		},
		"sanitization_manifest": func(a *ApprovedPackArtifactHandoffAttestation) {
			a.SanitizationManifestSHA256 = "abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		},
		"expiry": func(a *ApprovedPackArtifactHandoffAttestation) { a.ExpiresAt = a.IssuedAt },
		"evidence": func(a *ApprovedPackArtifactHandoffAttestation) {
			a.EvidenceRefs = []string{"receipt://same", "receipt://same"}
		},
	} {
		t.Run(name, func(t *testing.T) {
			candidate := attestation
			mutate(&candidate)
			if err := VerifyApprovedPackArtifactHandoff(candidate, privateKey.Public().(ed25519.PublicKey)); err == nil {
				t.Fatal("invalid or tampered attestation accepted")
			}
		})
	}
}
