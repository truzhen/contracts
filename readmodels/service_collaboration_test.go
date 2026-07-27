package readmodels

import "testing"

func TestServiceCollaborationContractsFailClosed(t *testing.T) {
	invite := ServiceCollaborationInvite{"invite", "tx-a", "participant-a", "nonce-a", "2026-07-27T00:00:00Z", true, false}
	if err := invite.Validate(); err != nil {
		t.Fatal(err)
	}
	invite.ProducesOwnerDecision = true
	if err := invite.Validate(); err == nil {
		t.Fatal("invite must not produce an OwnerDecision")
	}

	disclosure := ServiceAgreementDisclosure{"disclosure", "tx-a", "v1", "hash-v1", "2026-07-27T00:00:00Z", true, false}
	if err := disclosure.Validate(); err != nil {
		t.Fatal(err)
	}
	disclosure.FormalAgreement = true
	if err := disclosure.Validate(); err == nil {
		t.Fatal("disclosure must not claim a FormalAgreement")
	}

	decision := CounterpartyDecisionEnvelope{"decision", "tx-a", "participant-a", "v1", "hash-v1", "nonce-a", "accepted", true, false}
	if err := decision.Validate(); err != nil {
		t.Fatal(err)
	}
	decision.Decision = "owner_approved"
	if err := decision.Validate(); err == nil {
		t.Fatal("counterparty decision must reject OwnerDecision semantics")
	}

	delivery := DailyDeliveryDisclosure{"disclosure", "tx-a", "v1", []string{"i1", "i2", "i3", "i4", "i5"}, []string{"artifact-a"}, true}
	if err := delivery.Validate(); err != nil {
		t.Fatal(err)
	}
	delivery.ItemRefs[4] = "i4"
	if err := delivery.Validate(); err == nil {
		t.Fatal("delivery item refs must be unique")
	}

	scope := ParticipantScopeReadModel{"scope", "participant-a", []string{"tx-a"}, []string{"DocumentCommentCandidate"}, "2026-07-27T00:00:00Z", true, false}
	if err := scope.Validate(); err != nil {
		t.Fatal(err)
	}
	scope.AllowedTransactionRefs = []string{"tx-a", "tx-b"}
	if err := scope.Validate(); err == nil {
		t.Fatal("participant scope must stay transaction-bounded")
	}
}
