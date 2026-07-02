package base

import (
	"errors"
	"time"
)

// ActorContext represents the identity context of the caller
type ActorContext struct {
	OwnerProfile      string `json:"owner_profile,omitempty"`
	LocalIdentityRef  string `json:"local_identity_ref,omitempty"`
	DeviceIdentityRef string `json:"device_identity_ref,omitempty"`
}

// OwnerIdentityContext is the stable owner identity contract used by the policy kernel.
type OwnerIdentityContext struct {
	OwnerRef          string    `json:"owner_ref"`
	LocalIdentityRef  string    `json:"local_identity_ref"`
	DeviceIdentityRef string    `json:"device_identity_ref"`
	ActorRef          string    `json:"actor_ref"`
	AuthenticatedAt   time.Time `json:"authenticated_at"`
}

// OwnerDecision represents the owner's manual override or approval
type OwnerDecision struct {
	DecisionRef        string             `json:"decision_ref,omitempty"`
	OwnerDecisionRef   string             `json:"owner_decision_ref,omitempty"`
	TargetCandidateID  string             `json:"target_candidate_id"`
	TargetCandidateRef string             `json:"target_candidate_ref,omitempty"`
	Approved           bool               `json:"approved"`
	Comment            string             `json:"comment"`
	Status             GateDecisionStatus `json:"status"` // allow, deny, defer, requires_revision
	ActorContext       *ActorContext      `json:"actor_context,omitempty"`
	PolicySnapshotRef  string             `json:"policy_snapshot_ref,omitempty"`
	DecidedAt          time.Time          `json:"decided_at,omitempty"`
}

func ValidateOwnerIdentityContext(ctx *OwnerIdentityContext) error {
	if ctx == nil {
		return errors.New("owner identity context is required")
	}
	if ctx.OwnerRef == "" || ctx.LocalIdentityRef == "" || ctx.DeviceIdentityRef == "" || ctx.ActorRef == "" {
		return errors.New("owner_ref, local_identity_ref, device_identity_ref, and actor_ref are required")
	}
	if ctx.AuthenticatedAt.IsZero() {
		return errors.New("authenticated_at is required")
	}
	return nil
}

func ValidateOwnerDecisionForCandidate(ownerDecision *OwnerDecision, candidateRef string) error {
	if ownerDecision == nil {
		return errors.New("owner decision is required")
	}
	if candidateRef == "" {
		return errors.New("candidate_ref is required")
	}
	target := ownerDecision.TargetCandidateRef
	if target == "" {
		target = ownerDecision.TargetCandidateID
	}
	if target == "" {
		return errors.New("owner decision target candidate is required")
	}
	if target != candidateRef {
		return errors.New("owner decision target mismatch")
	}
	if ownerDecision.ActorContext == nil || !hasActorIdentity(ownerDecision.ActorContext) {
		return errors.New("owner decision actor context is required")
	}
	if ownerDecision.PolicySnapshotRef == "" {
		return errors.New("owner decision policy_snapshot_ref is required")
	}
	if ownerDecision.DecidedAt.IsZero() {
		return errors.New("owner decision decided_at is required")
	}
	return nil
}

func hasActorIdentity(actor *ActorContext) bool {
	return actor.OwnerProfile != "" || actor.LocalIdentityRef != "" || actor.DeviceIdentityRef != ""
}
