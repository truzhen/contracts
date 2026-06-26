package base

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// ─────────────────── OwnerDelegationGrant（反向授权契约，T06） ───────────────────
//
// 红线（proposal-20260611-addendum §2.2）：委托不是转让主权。Owner 预先签发
// 一张有边界、有期限、可撤销的授权书，由 Base 在每次裁定时强制校验边界。
// 高风险硬地板写死在 Base policy：high / critical 永远不可委托，grant 内容
// 无法覆盖（创建期与校验期双重强制）。

// DelegationGrantStatus is the grant lifecycle status.
type DelegationGrantStatus string

const (
	DelegationGrantActive                   DelegationGrantStatus = "active"
	DelegationGrantRevoked                  DelegationGrantStatus = "revoked"
	DelegationGrantExpired                  DelegationGrantStatus = "expired"
	DelegationGrantSuspendedByEmergencyStop DelegationGrantStatus = "suspended_by_emergency_stop"
)

func ValidDelegationGrantStatus(s DelegationGrantStatus) bool {
	switch s {
	case DelegationGrantActive, DelegationGrantRevoked, DelegationGrantExpired, DelegationGrantSuspendedByEmergencyStop:
		return true
	}
	return false
}

// DelegationQuota caps how many agent decisions a grant covers per day.
type DelegationQuota struct {
	PerDay int `json:"per_day"`
}

// DelegationScope is the boundary of an OwnerDelegationGrant. Every dimension
// is validated by Base on every single agent decision (fail-closed).
type DelegationScope struct {
	// TaskTypes whitelists the task types the delegate may sign
	// (e.g. stage / immediate; scheduled timetable creation is never implied).
	TaskTypes []string `json:"task_types"`
	// RiskCeiling is at most medium. high/critical is a Base hard floor and
	// can never be granted (ValidateOwnerDelegationGrant enforces it again).
	RiskCeiling RiskClass `json:"risk_ceiling"`
	// TransactionRefs, when non-empty, restricts the grant to these
	// transactions. Empty means all transactions.
	TransactionRefs []string `json:"transaction_refs,omitempty"`
	// PackRefs, when non-empty, restricts the grant to candidates produced by
	// these packs. Empty means no pack restriction.
	PackRefs []string `json:"pack_refs,omitempty"`
	// Quota caps signatures per day. Required (per_day >= 1).
	Quota DelegationQuota `json:"quota"`
	// AmountLimitCents caps money-related actions. 0 means the grant covers
	// no money actions at all: a candidate carrying any amount is denied.
	AmountLimitCents int64 `json:"amount_limit_cents,omitempty"`
}

// OwnerDelegationGrant is the Owner-issued, bounded, expiring, revocable
// authorization for a delegate agent to sign candidates on the Owner's behalf.
// Creating and revoking a grant are themselves gated actions (confirm card +
// OwnerDecision + GrantReceipt in the 03 ledger).
type OwnerDelegationGrant struct {
	GrantID          string                `json:"grant_id"`
	OwnerDecisionRef string                `json:"owner_decision_ref"`
	DelegateRef      string                `json:"delegate_ref"` // e.g. agent://secretary_chief
	Scope            DelegationScope       `json:"scope"`
	ExpiresAt        time.Time             `json:"expires_at"`
	Revocable        bool                  `json:"revocable"`
	Status           DelegationGrantStatus `json:"status"`
	ReceiptRef       string                `json:"receipt_ref,omitempty"`
	CreatedAt        time.Time             `json:"created_at,omitempty"`
	UpdatedAt        time.Time             `json:"updated_at,omitempty"`
}

// AgentDecision is the only new power a delegate agent gains: signing a
// candidate inside the boundary of an Owner grant. It never carries formal
// adjudication power of its own; Base validates the grant on every decision.
type AgentDecision struct {
	AgentDecisionRef   string    `json:"agent_decision_ref,omitempty"`
	CandidateRef       string    `json:"candidate_ref"`
	DelegationGrantRef string    `json:"delegation_grant_ref"`
	ActorRef           string    `json:"actor_ref"` // must equal grant.delegate_ref
	Reasoning          string    `json:"reasoning"` // required; recorded into the receipt
	DecidedAt          time.Time `json:"decided_at"`
}

// DelegationSubject is the server-derived description of the candidate an
// AgentDecision targets. It must be built from stored governance state by the
// module that owns the candidate (07), never from agent-supplied claims.
type DelegationSubject struct {
	CandidateRef   string    `json:"candidate_ref"`
	TransactionRef string    `json:"transaction_ref"`
	TaskType       string    `json:"task_type"`
	RiskLevel      RiskClass `json:"risk_level"`
	PackRef        string    `json:"pack_ref,omitempty"`
	AmountCents    int64     `json:"amount_cents,omitempty"`
}

// DelegationRiskWithinHardFloor is the Base policy hard floor: only low and
// medium risk may ever be delegated. This is not configurable by grants.
func DelegationRiskWithinHardFloor(risk RiskClass) bool {
	return risk == RiskLow || risk == RiskMedium
}

// ValidateDelegationScope enforces the scope shape, including the creation-time
// hard floor on risk_ceiling.
func ValidateDelegationScope(scope *DelegationScope) error {
	if scope == nil {
		return errors.New("delegation scope is required")
	}
	if len(scope.TaskTypes) == 0 {
		return errors.New("delegation scope task_types are required")
	}
	for _, t := range scope.TaskTypes {
		if t == "" {
			return errors.New("delegation scope task_types must not contain empty entries")
		}
	}
	if !DelegationRiskWithinHardFloor(scope.RiskCeiling) {
		return fmt.Errorf("delegation risk_ceiling %q violates the Base hard floor: only low or medium may be delegated", scope.RiskCeiling)
	}
	if scope.Quota.PerDay < 1 {
		return errors.New("delegation scope quota.per_day must be >= 1")
	}
	if scope.AmountLimitCents < 0 {
		return errors.New("delegation scope amount_limit_cents must not be negative")
	}
	return nil
}

// ValidateOwnerDelegationGrant validates a full grant object.
func ValidateOwnerDelegationGrant(grant *OwnerDelegationGrant) error {
	if grant == nil {
		return errors.New("owner delegation grant is required")
	}
	if grant.GrantID == "" {
		return errors.New("grant_id is required")
	}
	if grant.OwnerDecisionRef == "" {
		return errors.New("grant owner_decision_ref is required: a grant is itself a gated action")
	}
	if grant.DelegateRef == "" {
		return errors.New("grant delegate_ref is required")
	}
	if err := ValidateDelegationScope(&grant.Scope); err != nil {
		return err
	}
	if grant.ExpiresAt.IsZero() {
		return errors.New("grant expires_at is required")
	}
	if !ValidDelegationGrantStatus(grant.Status) {
		return fmt.Errorf("grant status %q is invalid", grant.Status)
	}
	return nil
}

// ValidateAgentDecisionForCandidate checks the agent decision shape against
// the candidate it claims to sign.
func ValidateAgentDecisionForCandidate(dec *AgentDecision, candidateRef string) error {
	if dec == nil {
		return errors.New("agent decision is required")
	}
	if candidateRef == "" {
		return errors.New("candidate_ref is required")
	}
	if dec.CandidateRef == "" || dec.CandidateRef != candidateRef {
		return errors.New("agent decision target candidate mismatch")
	}
	if dec.DelegationGrantRef == "" {
		return errors.New("agent decision delegation_grant_ref is required: agent power derives only from an owner grant")
	}
	if dec.ActorRef == "" {
		return errors.New("agent decision actor_ref is required")
	}
	if dec.Reasoning == "" {
		return errors.New("agent decision reasoning is required and is recorded in the receipt")
	}
	if dec.DecidedAt.IsZero() {
		return errors.New("agent decision decided_at is required")
	}
	return nil
}

// DelegationGrantCandidateRef is the public candidate-ref formula for grant
// creation; the OwnerDecision on the confirm card must target this ref.
func DelegationGrantCandidateRef(idempotencyKey string) string {
	sum := sha256.Sum256([]byte(idempotencyKey + "\x00delegation_grant"))
	return "delegation_grant_candidate_" + hex.EncodeToString(sum[:])[:16]
}

// DelegationGrantActionCandidateRef is the candidate-ref formula for gated
// grant lifecycle actions (revoke / resume).
func DelegationGrantActionCandidateRef(action, grantID string) string {
	sum := sha256.Sum256([]byte(action + "\x00" + grantID + "\x00delegation_grant_action"))
	return "delegation_grant_" + action + "_candidate_" + hex.EncodeToString(sum[:])[:16]
}
