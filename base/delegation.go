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

// ExecutionNetworkPolicy declares the strongest network surface an execution
// delegation may use. gated_bridge may appear on server-derived subjects, but
// grant ceilings deliberately allow only none or egress_model_only.
type ExecutionNetworkPolicy string

const (
	ExecutionNetworkPolicyNone            ExecutionNetworkPolicy = "none"
	ExecutionNetworkPolicyEgressModelOnly ExecutionNetworkPolicy = "egress_model_only"
	ExecutionNetworkPolicyGatedBridge     ExecutionNetworkPolicy = "gated_bridge"
)

func ValidExecutionNetworkPolicy(policy ExecutionNetworkPolicy) bool {
	switch policy {
	case ExecutionNetworkPolicyNone, ExecutionNetworkPolicyEgressModelOnly, ExecutionNetworkPolicyGatedBridge:
		return true
	}
	return false
}

func DelegationExecutionScopeNetworkCeilingAllowed(policy ExecutionNetworkPolicy) bool {
	return policy == ExecutionNetworkPolicyNone || policy == ExecutionNetworkPolicyEgressModelOnly
}

// DelegationExecutionScope is the code-execution boundary of an
// OwnerDelegationGrant. It is optional and never implied for legacy grants.
type DelegationExecutionScope struct {
	CapabilityRefs       []string               `json:"capability_refs"`
	WorkrootRef          string                 `json:"workroot_ref"`
	ProviderRefs         []string               `json:"provider_refs"`
	SandboxProfileRef    string                 `json:"sandbox_profile_ref"`
	NetworkPolicyCeiling ExecutionNetworkPolicy `json:"network_policy_ceiling"`
	MaxRuns              int                    `json:"max_runs"`
	MaxDurationSeconds   int                    `json:"max_duration_seconds"`
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
	// ExecutionScope optionally whitelists bounded code execution. Absence
	// means the grant carries no execution authority.
	ExecutionScope *DelegationExecutionScope `json:"execution_scope,omitempty"`
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
	CandidateRef           string                      `json:"candidate_ref"`
	TransactionRef         string                      `json:"transaction_ref"`
	TaskType               string                      `json:"task_type"`
	RiskLevel              RiskClass                   `json:"risk_level"`
	PackRef                string                      `json:"pack_ref,omitempty"`
	AmountCents            int64                       `json:"amount_cents,omitempty"`
	SideEffectClass        SideEffectClass             `json:"side_effect_class,omitempty"`
	QuotaDate              string                      `json:"quota_date,omitempty"`
	ConsumedDecisionsToday int                         `json:"consumed_decisions_today,omitempty"`
	Execution              *DelegationExecutionSubject `json:"execution,omitempty"`
}

// DelegationExecutionSubject is the server-derived cumulative execution fact
// after atomically reserving the run currently being evaluated. ConsumedRuns
// and ConsumedDurationSeconds therefore include the proposed run and must be
// positive. Concurrent consumers must use OCC or an equivalent atomic
// compare-and-reserve operation before constructing this subject; check-then-
// increment is not sufficient. It must not be accepted from the delegate agent
// as a self-asserted claim.
type DelegationExecutionSubject struct {
	CapabilityRef           string                 `json:"capability_ref"`
	WorkrootRef             string                 `json:"workroot_ref"`
	ProviderRef             string                 `json:"provider_ref"`
	SandboxProfileRef       string                 `json:"sandbox_profile_ref"`
	NetworkPolicy           ExecutionNetworkPolicy `json:"network_policy"`
	ConsumedRuns            int                    `json:"consumed_runs"`
	ConsumedDurationSeconds int                    `json:"consumed_duration_seconds"`
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
	if scope.ExecutionScope != nil {
		if err := ValidateDelegationExecutionScope(scope.ExecutionScope); err != nil {
			return err
		}
	}
	return nil
}

// ValidateDelegationExecutionScope enforces the optional code-execution grant
// boundary. Legacy grants omit this object and therefore grant no execution.
func ValidateDelegationExecutionScope(scope *DelegationExecutionScope) error {
	if scope == nil {
		return errors.New("delegation execution_scope is required")
	}
	if err := validateNonEmptyUniqueRefs("delegation execution_scope capability_refs", scope.CapabilityRefs); err != nil {
		return err
	}
	if scope.WorkrootRef == "" {
		return errors.New("delegation execution_scope workroot_ref is required")
	}
	if err := validateNonEmptyUniqueRefs("delegation execution_scope provider_refs", scope.ProviderRefs); err != nil {
		return err
	}
	if scope.SandboxProfileRef == "" {
		return errors.New("delegation execution_scope sandbox_profile_ref is required")
	}
	if !DelegationExecutionScopeNetworkCeilingAllowed(scope.NetworkPolicyCeiling) {
		return fmt.Errorf("delegation execution_scope network_policy_ceiling %q is invalid: only none or egress_model_only may be delegated", scope.NetworkPolicyCeiling)
	}
	if scope.MaxRuns < 1 {
		return errors.New("delegation execution_scope max_runs must be >= 1")
	}
	if scope.MaxDurationSeconds < 1 {
		return errors.New("delegation execution_scope max_duration_seconds must be >= 1")
	}
	return nil
}

// ValidateDelegationExecutionSubject validates server-derived execution facts
// before comparing them with an OwnerDelegationGrant scope.
func ValidateDelegationExecutionSubject(subject *DelegationExecutionSubject) error {
	if subject == nil {
		return errors.New("delegation execution subject is required")
	}
	if subject.CapabilityRef == "" {
		return errors.New("delegation execution subject capability_ref is required")
	}
	if subject.WorkrootRef == "" {
		return errors.New("delegation execution subject workroot_ref is required")
	}
	if subject.ProviderRef == "" {
		return errors.New("delegation execution subject provider_ref is required")
	}
	if subject.SandboxProfileRef == "" {
		return errors.New("delegation execution subject sandbox_profile_ref is required")
	}
	if !ValidExecutionNetworkPolicy(subject.NetworkPolicy) {
		return fmt.Errorf("delegation execution subject network_policy %q is invalid", subject.NetworkPolicy)
	}
	if subject.ConsumedRuns < 1 {
		return errors.New("delegation execution subject consumed_runs must be >= 1")
	}
	if subject.ConsumedDurationSeconds < 1 {
		return errors.New("delegation execution subject consumed_duration_seconds must be >= 1")
	}
	return nil
}

// ValidateDelegationSubject validates the server-derived parent facts before
// any grant boundary comparison.
func ValidateDelegationSubject(subject *DelegationSubject) error {
	if subject == nil {
		return errors.New("delegation subject is required")
	}
	if subject.CandidateRef == "" {
		return errors.New("delegation subject candidate_ref is required")
	}
	if subject.TransactionRef == "" {
		return errors.New("delegation subject transaction_ref is required")
	}
	if subject.TaskType == "" {
		return errors.New("delegation subject task_type is required")
	}
	if !DelegationRiskWithinHardFloor(subject.RiskLevel) {
		return fmt.Errorf("delegation subject risk_level %q violates the Base hard floor: only low or medium may be delegated", subject.RiskLevel)
	}
	if subject.AmountCents < 0 {
		return errors.New("delegation subject amount_cents must not be negative")
	}
	return nil
}

// DelegationWithinScope checks all non-execution dimensions of a server-
// derived subject against a delegation scope. The Base risk hard floor is
// evaluated independently of the configurable risk ceiling.
func DelegationWithinScope(scope *DelegationScope, subject *DelegationSubject) error {
	if err := ValidateDelegationScope(scope); err != nil {
		return err
	}
	if err := ValidateDelegationSubject(subject); err != nil {
		return err
	}
	if !stringInSet(subject.TaskType, scope.TaskTypes) {
		return fmt.Errorf("delegation subject task_type %q is outside scope", subject.TaskType)
	}
	if !delegationRiskWithinCeiling(scope.RiskCeiling, subject.RiskLevel) {
		return fmt.Errorf("delegation subject risk_level %q exceeds ceiling %q", subject.RiskLevel, scope.RiskCeiling)
	}
	if len(scope.TransactionRefs) > 0 && !stringInSet(subject.TransactionRef, scope.TransactionRefs) {
		return fmt.Errorf("delegation subject transaction_ref %q is outside scope", subject.TransactionRef)
	}
	if len(scope.PackRefs) > 0 && !stringInSet(subject.PackRef, scope.PackRefs) {
		return fmt.Errorf("delegation subject pack_ref %q is outside scope", subject.PackRef)
	}
	if subject.AmountCents > scope.AmountLimitCents {
		return fmt.Errorf("delegation subject amount_cents %d exceeds amount_limit_cents %d", subject.AmountCents, scope.AmountLimitCents)
	}
	return nil
}

// DelegationGrantWithinScope is the complete authorization-boundary helper for
// an OwnerDelegationGrant and a server-derived DelegationSubject at an explicit
// evaluation time. It requires an active, unexpired grant, then checks the full
// parent scope, including the Base hard floor, before checking the optional
// execution boundary. Callers must use this entry point for execution
// authorization rather than treating the lower-level execution helper as a
// complete grant decision.
func DelegationGrantWithinScope(grant *OwnerDelegationGrant, subject *DelegationSubject, evaluationTime time.Time) error {
	if err := ValidateOwnerDelegationGrant(grant); err != nil {
		return err
	}
	if !grant.Revocable {
		return errors.New("delegation grant revocable must be true for code execution authorization")
	}
	if grant.ReceiptRef == "" {
		return errors.New("delegation grant receipt_ref is required for authorization")
	}
	if evaluationTime.IsZero() {
		return errors.New("delegation evaluation time is required")
	}
	if grant.Status != DelegationGrantActive {
		return fmt.Errorf("delegation grant status %q is not active", grant.Status)
	}
	if !grant.ExpiresAt.After(evaluationTime) {
		return fmt.Errorf("delegation grant expired at %s before or at evaluation time %s", grant.ExpiresAt.Format(time.RFC3339Nano), evaluationTime.Format(time.RFC3339Nano))
	}
	if err := DelegationWithinScope(&grant.Scope, subject); err != nil {
		return err
	}
	if err := validateDelegationAuthorizationFacts(&grant.Scope, subject, evaluationTime); err != nil {
		return err
	}
	if subject.Execution == nil {
		return nil
	}
	return DelegationExecutionWithinScope(&grant.Scope, subject.Execution)
}

func validateDelegationAuthorizationFacts(scope *DelegationScope, subject *DelegationSubject, evaluationTime time.Time) error {
	if !validDelegationSideEffectClass(subject.SideEffectClass) {
		return errors.New("delegation subject side_effect_class is required and must be known")
	}
	for _, denied := range AuthorizationHardDenies() {
		if subject.SideEffectClass == denied {
			return fmt.Errorf("delegation subject side_effect_class %q is never delegable", subject.SideEffectClass)
		}
	}
	if subject.SideEffectClass == SideEffectExternalSend {
		return errors.New("delegation subject external_send is never delegable")
	}
	expectedQuotaDate := evaluationTime.UTC().Format("2006-01-02")
	if subject.QuotaDate != expectedQuotaDate {
		return fmt.Errorf("delegation subject quota_date %q does not match evaluation date %q", subject.QuotaDate, expectedQuotaDate)
	}
	if subject.ConsumedDecisionsToday < 1 {
		return errors.New("delegation subject consumed_decisions_today must include the reserved decision")
	}
	if subject.ConsumedDecisionsToday > scope.Quota.PerDay {
		return fmt.Errorf("delegation subject consumed_decisions_today %d exceeds quota.per_day %d", subject.ConsumedDecisionsToday, scope.Quota.PerDay)
	}
	return nil
}

func validDelegationSideEffectClass(sideEffect SideEffectClass) bool {
	switch sideEffect {
	case SideEffectReadOnly,
		SideEffectLocalDraft,
		SideEffectFormalWrite,
		SideEffectExternalSend,
		SideEffectLocalFileWrite,
		SideEffectGuiControl,
		SideEffectNetworkCall,
		SideEffectPayment,
		SideEffectDelete,
		SideEffectCredentialAccess,
		SideEffectRealExecution:
		return true
	}
	return false
}

// DelegationExecutionWithinScope checks a server-derived execution subject
// against an Owner grant boundary. It compares refs as opaque strings and never
// parses local paths or provider-specific identifiers. This is a lower-level
// dimensional check; use DelegationGrantWithinScope for an authorization
// decision so parent scope and hard-floor checks cannot be skipped.
func DelegationExecutionWithinScope(scope *DelegationScope, subject *DelegationExecutionSubject) error {
	if scope == nil {
		return errors.New("delegation scope is required")
	}
	if subject == nil {
		return errors.New("delegation execution subject is required")
	}
	if scope.ExecutionScope == nil {
		return errors.New("delegation scope has no execution_scope: legacy grants do not authorize code execution")
	}
	if err := ValidateDelegationScope(scope); err != nil {
		return err
	}
	if err := ValidateDelegationExecutionSubject(subject); err != nil {
		return err
	}

	execScope := scope.ExecutionScope
	if !stringInSet(subject.CapabilityRef, execScope.CapabilityRefs) {
		return fmt.Errorf("delegation execution subject capability_ref %q is outside scope", subject.CapabilityRef)
	}
	if subject.WorkrootRef != execScope.WorkrootRef {
		return fmt.Errorf("delegation execution subject workroot_ref %q is outside scope", subject.WorkrootRef)
	}
	if !stringInSet(subject.ProviderRef, execScope.ProviderRefs) {
		return fmt.Errorf("delegation execution subject provider_ref %q is outside scope", subject.ProviderRef)
	}
	if subject.SandboxProfileRef != execScope.SandboxProfileRef {
		return fmt.Errorf("delegation execution subject sandbox_profile_ref %q is outside scope", subject.SandboxProfileRef)
	}
	if !executionNetworkPolicyWithinCeiling(execScope.NetworkPolicyCeiling, subject.NetworkPolicy) {
		return fmt.Errorf("delegation execution subject network_policy %q exceeds ceiling %q", subject.NetworkPolicy, execScope.NetworkPolicyCeiling)
	}
	if subject.ConsumedRuns > execScope.MaxRuns {
		return fmt.Errorf("delegation execution subject consumed_runs %d exceeds max_runs %d", subject.ConsumedRuns, execScope.MaxRuns)
	}
	if subject.ConsumedDurationSeconds > execScope.MaxDurationSeconds {
		return fmt.Errorf("delegation execution subject consumed_duration_seconds %d exceeds max_duration_seconds %d", subject.ConsumedDurationSeconds, execScope.MaxDurationSeconds)
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

func validateNonEmptyUniqueRefs(field string, refs []string) error {
	if len(refs) == 0 {
		return fmt.Errorf("%s are required", field)
	}
	seen := make(map[string]struct{}, len(refs))
	for _, ref := range refs {
		if ref == "" {
			return fmt.Errorf("%s must not contain empty entries", field)
		}
		if _, ok := seen[ref]; ok {
			return fmt.Errorf("%s must not contain duplicate entries: %q", field, ref)
		}
		seen[ref] = struct{}{}
	}
	return nil
}

func stringInSet(value string, allowed []string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}

func delegationRiskWithinCeiling(ceiling, risk RiskClass) bool {
	switch ceiling {
	case RiskLow:
		return risk == RiskLow
	case RiskMedium:
		return risk == RiskLow || risk == RiskMedium
	default:
		return false
	}
}

func executionNetworkPolicyWithinCeiling(ceiling, policy ExecutionNetworkPolicy) bool {
	switch ceiling {
	case ExecutionNetworkPolicyNone:
		return policy == ExecutionNetworkPolicyNone
	case ExecutionNetworkPolicyEgressModelOnly:
		return policy == ExecutionNetworkPolicyNone || policy == ExecutionNetworkPolicyEgressModelOnly
	default:
		return false
	}
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
