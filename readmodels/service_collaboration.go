package readmodels

import (
	"errors"
	"strings"
)

// ServiceCollaborationInvite is a host-created, single-use invitation shape.
// It carries no credential and never grants authority by itself.
type ServiceCollaborationInvite struct {
	InviteRef              string `json:"invite_ref"`
	TransactionRef         string `json:"transaction_ref"`
	ExpectedParticipantRef string `json:"expected_participant_ref"`
	Nonce                  string `json:"nonce"`
	ExpiresAt              string `json:"expires_at"`
	CandidateOnly          bool   `json:"candidate_only"`
	ProducesOwnerDecision  bool   `json:"produces_owner_decision"`
}

// ServiceAgreementDisclosure is a read-only, version-bound disclosure.
type ServiceAgreementDisclosure struct {
	DisclosureRef   string `json:"disclosure_ref"`
	TransactionRef  string `json:"transaction_ref"`
	ProposalVersion string `json:"proposal_version"`
	ProposalHash    string `json:"proposal_hash"`
	ExpiresAt       string `json:"expires_at"`
	ReadModelOnly   bool   `json:"read_model_only"`
	FormalAgreement bool   `json:"formal_agreement"`
}

// CounterpartyDecisionEnvelope records counterparty evidence, never an OwnerDecision.
type CounterpartyDecisionEnvelope struct {
	DecisionRef           string `json:"decision_ref"`
	TransactionRef        string `json:"transaction_ref"`
	ParticipantRef        string `json:"participant_ref"`
	ProposalVersion       string `json:"proposal_version"`
	ProposalHash          string `json:"proposal_hash"`
	Nonce                 string `json:"nonce"`
	Decision              string `json:"decision"`
	CandidateOnly         bool   `json:"candidate_only"`
	ProducesOwnerDecision bool   `json:"produces_owner_decision"`
}

// DailyDeliveryDisclosure exposes only approved delivery references.
type DailyDeliveryDisclosure struct {
	DisclosureRef    string   `json:"disclosure_ref"`
	TransactionRef   string   `json:"transaction_ref"`
	AgreementVersion string   `json:"agreement_version"`
	ItemRefs         []string `json:"item_refs"`
	ArtifactRefs     []string `json:"artifact_refs"`
	ReadModelOnly    bool     `json:"read_model_only"`
}

// ParticipantScopeReadModel is a projection of host-owned scope, never the scope truth.
type ParticipantScopeReadModel struct {
	ScopeRef                string   `json:"scope_ref"`
	ParticipantRef          string   `json:"participant_ref"`
	AllowedTransactionRefs  []string `json:"allowed_transaction_refs"`
	AllowedCandidateTypes   []string `json:"allowed_candidate_types"`
	ExpiresAt               string   `json:"expires_at"`
	ReadModelOnly           bool     `json:"read_model_only"`
	CanProduceOwnerDecision bool     `json:"can_produce_owner_decision"`
}

func (v ServiceCollaborationInvite) Validate() error {
	if blank(v.InviteRef, v.TransactionRef, v.ExpectedParticipantRef, v.Nonce, v.ExpiresAt) || !v.CandidateOnly || v.ProducesOwnerDecision {
		return errors.New("invalid_service_collaboration_invite")
	}
	return nil
}

func (v ServiceAgreementDisclosure) Validate() error {
	if blank(v.DisclosureRef, v.TransactionRef, v.ProposalVersion, v.ProposalHash, v.ExpiresAt) || !v.ReadModelOnly || v.FormalAgreement {
		return errors.New("invalid_service_agreement_disclosure")
	}
	return nil
}

func (v CounterpartyDecisionEnvelope) Validate() error {
	if blank(v.DecisionRef, v.TransactionRef, v.ParticipantRef, v.ProposalVersion, v.ProposalHash, v.Nonce) || !v.CandidateOnly || v.ProducesOwnerDecision {
		return errors.New("invalid_counterparty_decision")
	}
	switch v.Decision {
	case "accepted", "revision_requested", "rejected":
		return nil
	default:
		return errors.New("invalid_counterparty_decision_kind")
	}
}

func (v DailyDeliveryDisclosure) Validate() error {
	if blank(v.DisclosureRef, v.TransactionRef, v.AgreementVersion) || !v.ReadModelOnly || len(v.ItemRefs) != 5 || hasBlankOrDuplicate(v.ItemRefs) || hasBlankOrDuplicate(v.ArtifactRefs) {
		return errors.New("invalid_daily_delivery_disclosure")
	}
	return nil
}

func (v ParticipantScopeReadModel) Validate() error {
	if blank(v.ScopeRef, v.ParticipantRef, v.ExpiresAt) || !v.ReadModelOnly || v.CanProduceOwnerDecision || len(v.AllowedTransactionRefs) != 1 || hasBlankOrDuplicate(v.AllowedTransactionRefs) || hasBlankOrDuplicate(v.AllowedCandidateTypes) {
		return errors.New("invalid_participant_scope_readmodel")
	}
	return nil
}

func blank(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	return false
}

func hasBlankOrDuplicate(values []string) bool {
	seen := map[string]struct{}{}
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
		if _, exists := seen[value]; exists {
			return true
		}
		seen[value] = struct{}{}
	}
	return false
}
