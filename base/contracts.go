package base

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/truzhen/contracts/spines"
)

const (
	PolicySnapshotIssuerBaseKernel = "base-policy-kernel"
	AdapterKindMemoryWrite         = "memory_write"
	AdapterKindSend                = "send"
	AdapterKindExecution           = "execution"
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

// OwnerDirective represents an owner's policy directive
type OwnerDirective struct {
	DirectiveID string `json:"directive_id"`
	Content     string `json:"content"`
}

// PolicySet holds a set of active policies
type PolicySet struct {
	SetID      string           `json:"set_id"`
	Version    string           `json:"version,omitempty"`
	Directives []OwnerDirective `json:"directives"`
}

// GatePolicySet is the replayable policy set emitted by the Base policy kernel.
type GatePolicySet struct {
	PolicySetRef         string    `json:"policy_set_ref"`
	Version              string    `json:"version"`
	OwnerPolicyRef       string    `json:"owner_policy_ref,omitempty"`
	RiskPolicyRef        string    `json:"risk_policy_ref,omitempty"`
	SideEffectPolicyRef  string    `json:"side_effect_policy_ref,omitempty"`
	FormalWritePolicyRef string    `json:"formal_write_policy_ref,omitempty"`
	CapabilityPolicyRef  string    `json:"capability_policy_ref,omitempty"`
	ModuleBoundaryRef    string    `json:"module_boundary_ref,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

// PolicyKernel identifies the Base-owned kernel that can mint policy snapshots.
type PolicyKernel struct {
	KernelRef string `json:"kernel_ref"`
	Version   string `json:"version"`
}

// PolicySnapshot represents an immutable snapshot of policies used during gate evaluation
type PolicySnapshot struct {
	SnapshotRef   string         `json:"snapshot_ref,omitempty"`
	SnapshotHash  string         `json:"snapshot_hash"`
	PolicySet     PolicySet      `json:"policy_set"`
	GatePolicySet *GatePolicySet `json:"gate_policy_set,omitempty"`
	KernelRef     string         `json:"kernel_ref,omitempty"`
	Version       string         `json:"version,omitempty"`
	IssuedBy      string         `json:"issued_by,omitempty"`
	CreatedAt     time.Time      `json:"created_at,omitempty"`
}

// SideEffectClass defines the category of side effect
type SideEffectClass string

const (
	SideEffectReadOnly         SideEffectClass = "read_only"
	SideEffectLocalDraft       SideEffectClass = "local_draft"
	SideEffectFormalWrite      SideEffectClass = "formal_write"
	SideEffectExternalSend     SideEffectClass = "external_send"
	SideEffectLocalFileWrite   SideEffectClass = "local_file_write"
	SideEffectGuiControl       SideEffectClass = "gui_control"
	SideEffectNetworkCall      SideEffectClass = "network_call"
	SideEffectPayment          SideEffectClass = "payment"
	SideEffectDelete           SideEffectClass = "delete"
	SideEffectCredentialAccess SideEffectClass = "credential_access"
	SideEffectRealExecution    SideEffectClass = "real_execution"
	SideEffectUnknown          SideEffectClass = "unknown"
)

// RiskClass defines the risk level of an action
type RiskClass string

const (
	RiskLow      RiskClass = "low"
	RiskMedium   RiskClass = "medium"
	RiskHigh     RiskClass = "high"
	RiskCritical RiskClass = "critical"
)

// IdempotencyKey prevents duplicate processing
type IdempotencyKey string

// GateCandidateEnvelope wraps candidate data for gate evaluation
type GateCandidateEnvelope struct {
	EnvelopeID           string                         `json:"envelope_id"`
	CandidateRef         string                         `json:"candidate_ref,omitempty"`
	TransactionRef       string                         `json:"transaction_ref,omitempty"`
	SourceEventID        string                         `json:"source_event_id,omitempty"`
	SourceModule         string                         `json:"source_module,omitempty"`
	IntentEventRef       string                         `json:"intent_event_ref,omitempty"`
	EvidenceRefs         []string                       `json:"evidence_refs,omitempty"`
	CandidateOnly        bool                           `json:"candidate_only"`
	NonFormal            bool                           `json:"non_formal"`
	RiskClass            RiskClass                      `json:"risk_class"`
	SideEffectClass      SideEffectClass                `json:"side_effect_class"`
	Payload              interface{}                    `json:"payload"`
	PayloadHash          string                         `json:"payload_hash,omitempty"`
	CapabilityInvocation *CapabilityInvocationCandidate `json:"capability_invocation,omitempty"`
	FormalWrite          *FormalWriteRequest            `json:"formal_write,omitempty"`
	ReceiptLink          *spines.ReceiptLink            `json:"receipt_link,omitempty"`
	IdempotencyKey       IdempotencyKey                 `json:"idempotency_key,omitempty"`
}

// GateRequest represents an evaluation request sent to the Gate Orchestrator
type GateRequest struct {
	Candidate          *GateCandidateEnvelope `json:"candidate"`
	ActorContext       *ActorContext      `json:"actor_context"`
	TransactionRef     string             `json:"transaction_ref,omitempty"`
	CandidateRef       string             `json:"candidate_ref,omitempty"`
	IntentEventRef     string             `json:"intent_event_ref,omitempty"`
	EvidenceRefs       []string           `json:"evidence_refs,omitempty"`
	ReceiptRequirement string             `json:"receipt_requirement,omitempty"`
	PolicySnapshotRef  string             `json:"policy_snapshot_ref,omitempty"`

	// IssuedDecisionRef / IssuedRunID / IssuedNonce carry a Base-issued,
	// already-verified gated-action decision (the T06 prepare/confirm trio).
	// A gateway base-gate adapter must BIND its decision_ref and receipt to
	// these refs rather than minting its own dev refs; their absence means no
	// verified owner authorization is present and the adapter must fail closed
	// (no auto-allow). The real Base orchestrator ignores these (it evaluates
	// the gate chain itself); only thin gateway adapters downstream of a
	// verified T06 grant consume them.
	IssuedDecisionRef string `json:"issued_decision_ref,omitempty"`
	IssuedRunID       string `json:"issued_run_id,omitempty"`
	IssuedNonce       string `json:"issued_nonce,omitempty"`
}

// GateDecisionStatus is the status of the gate decision
type GateDecisionStatus string

const (
	GateDecisionAllow              GateDecisionStatus = "allow"
	GateDecisionDeny               GateDecisionStatus = "deny"
	GateDecisionPendingOwner       GateDecisionStatus = "pending_owner"
	GateDecisionRequiresRevision   GateDecisionStatus = "requires_revision"
	GateDecisionDefer              GateDecisionStatus = "defer"
	GateDecisionBlockedByPolicy    GateDecisionStatus = "blocked_by_policy"
	GateDecisionBlockedByEmergency GateDecisionStatus = "blocked_by_emergency_stop"
)

// GateDecisionTrace records the deterministic per-gate explanation for a decision.
type GateDecisionTrace struct {
	GateName           string             `json:"gate_name"`
	Status             GateDecisionStatus `json:"status"`
	Reason             string             `json:"reason"`
	PolicySnapshotHash string             `json:"policy_snapshot_hash"`
	EvidenceRefs       []string           `json:"evidence_refs,omitempty"`
	CheckedAt          time.Time          `json:"checked_at"`
}

// GateReceiptCandidate represents the evidence of a gate decision, required for formalization
type GateReceiptCandidate struct {
	ReceiptID          string         `json:"receipt_id"`
	TransactionRef     string         `json:"transaction_ref,omitempty"`
	CandidateRef       string         `json:"candidate_ref"`
	DecisionRef        string         `json:"decision_ref,omitempty"`
	PolicySnapshot     string         `json:"policy_snapshot_hash"`
	SourceModule       string         `json:"source_module,omitempty"`
	PayloadHash        string         `json:"payload_hash,omitempty"`
	ActorContext       *ActorContext  `json:"actor_context"`
	DecisionReason     string         `json:"decision_reason"`
	ReceiptRequirement string         `json:"receipt_requirement,omitempty"`
	EvidenceRefs       []string       `json:"evidence_refs,omitempty"`
	IdempotencyKey     IdempotencyKey `json:"idempotency_key,omitempty"`
}

// GateDecision represents the outcome of a gate evaluation
type GateDecision struct {
	DecisionRef         string                `json:"decision_ref,omitempty"`
	Status              GateDecisionStatus    `json:"status"`
	Reason              string                `json:"reason"`
	GateTrace           []GateDecisionTrace   `json:"gate_trace,omitempty"`
	DecidedAt           time.Time             `json:"decided_at,omitempty"`
	ActorContext        *ActorContext         `json:"actor_context,omitempty"`
	CandidateRef        string                `json:"candidate_ref,omitempty"`
	TransactionRef      string                `json:"transaction_ref,omitempty"`
	PolicySnapshotHash  string                `json:"policy_snapshot_hash,omitempty"`
	OwnerDecisionRef    string                `json:"owner_decision_ref,omitempty"`
	ReceiptCandidateRef string                `json:"receipt_candidate_ref,omitempty"`
	PolicySnapshot      *PolicySnapshot       `json:"policy_snapshot,omitempty"`
	ReceiptCandidate    *GateReceiptCandidate `json:"receipt_candidate,omitempty"`
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

// FormalWriteRequest represents a request to formally write data, requiring a receipt
type FormalWriteRequest struct {
	Payload            interface{}             `json:"payload"`
	GateDecision       *GateDecision           `json:"gate_decision,omitempty"`
	ReceiptCandidate   *GateReceiptCandidate   `json:"receipt_candidate"`
	FormalizationGrant *BaseFormalizationGrant `json:"formalization_grant,omitempty"`
}

// BaseFormalizationGrant is the only Base-issued grant that can authorize a formal write.
type BaseFormalizationGrant struct {
	GrantRef            string    `json:"grant_ref"`
	DecisionRef         string    `json:"decision_ref"`
	CandidateRef        string    `json:"candidate_ref"`
	TransactionRef      string    `json:"transaction_ref"`
	ReceiptCandidateRef string    `json:"receipt_candidate_ref"`
	PolicySnapshotHash  string    `json:"policy_snapshot_hash"`
	TargetObjectType    string    `json:"target_object_type"`
	WriteScope          string    `json:"write_scope"`
	GrantedAt           time.Time `json:"granted_at"`
	ExpiresAt           time.Time `json:"expires_at"`
}

// CapabilityInvocationCandidate is the structured capability request; Payload is not the protocol.
type CapabilityInvocationCandidate struct {
	CapabilityID      string   `json:"capability_id"`
	CapabilityVersion string   `json:"capability_version,omitempty"`
	InvocationRef     string   `json:"invocation_ref"`
	TransactionRef    string   `json:"transaction_ref"`
	CandidateRef      string   `json:"candidate_ref"`
	RequiredScopes    []string `json:"required_scopes,omitempty"`
	DependencyRefs    []string `json:"dependency_refs,omitempty"`
}

// CapabilityInvocationGateRequest represents a request to invoke a capability
type CapabilityInvocationGateRequest struct {
	CapabilityID string                         `json:"capability_id"`
	Candidate    *GateCandidateEnvelope             `json:"candidate"`
	ActorContext *ActorContext                  `json:"actor_context"`
	Invocation   *CapabilityInvocationCandidate `json:"invocation,omitempty"`
}

// EmergencyStopState defines the emergency stop status
type EmergencyStopState string

const (
	EmergencyStopDisabled EmergencyStopState = "emergency_stop_disabled"
	EmergencyStopEnabled  EmergencyStopState = "emergency_stop_enabled"
)

// GateAuditEvent is an audit log record of a gate decision
type GateAuditEvent struct {
	EventID       string         `json:"event_id"`
	GateDecision  *GateDecision  `json:"gate_decision"`
	OwnerDecision *OwnerDecision `json:"owner_decision,omitempty"`
	Timestamp     int64          `json:"timestamp"`
}

// CandidateIntakeEnvelope records the Base-owned intake result before gate orchestration.
type CandidateIntakeEnvelope struct {
	IntakeRef          string          `json:"intake_ref"`
	CandidateRef       string          `json:"candidate_ref"`
	TransactionRef     string          `json:"transaction_ref"`
	IntentEventRef     string          `json:"intent_event_ref"`
	EvidenceRefs       []string        `json:"evidence_refs"`
	RiskClass          RiskClass       `json:"risk_class"`
	SideEffectClass    SideEffectClass `json:"side_effect_class"`
	PolicySnapshotRef  string          `json:"policy_snapshot_ref"`
	ReceiptRequirement string          `json:"receipt_requirement"`
	ReceivedAt         time.Time       `json:"received_at"`
}

// GateOrchestrationResult is the stable envelope returned by the gate engine.
type GateOrchestrationResult struct {
	ResultRef         string              `json:"result_ref"`
	CandidateRef      string              `json:"candidate_ref"`
	TransactionRef    string              `json:"transaction_ref"`
	Decision          *GateDecision       `json:"decision"`
	PolicySnapshotRef string              `json:"policy_snapshot_ref"`
	Trace             []GateDecisionTrace `json:"trace"`
	CreatedAt         time.Time           `json:"created_at"`
}

// ModuleGateAdoptionStatus is the read-model contract showing whether a module accepts Base governance.
type ModuleGateAdoptionStatus struct {
	ModuleID                string    `json:"module_id"`
	ModuleName              string    `json:"module_name"`
	RequiresGateDecision    bool      `json:"requires_gate_decision"`
	LastGatedRequestRef     string    `json:"last_gated_request_ref,omitempty"`
	LastReceiptCandidateRef string    `json:"last_receipt_candidate_ref,omitempty"`
	LastFailureReason       string    `json:"last_failure_reason,omitempty"`
	AdoptionState           string    `json:"adoption_state"`
	CheckedAt               time.Time `json:"checked_at"`
}

type RegistrySliceContextRef struct {
	SliceRef            string `json:"slice_ref"`
	TransactionRef      string `json:"transaction_ref"`
	IntentEventRef      string `json:"intent_event_ref"`
	ActorRef            string `json:"actor_ref"`
	ScopePolicyRef      string `json:"scope_policy_ref"`
	AuditRef            string `json:"audit_ref"`
	ReceiptCandidateRef string `json:"receipt_candidate_ref,omitempty"`
}

type FormalMemoryWriteRef struct {
	MemoryWriteCandidateRef string `json:"memory_write_candidate_ref"`
	FormalMemoryStoreRef    string `json:"formal_memory_store_ref"`
	TransactionRef          string `json:"transaction_ref"`
}

type ObjectDirectoryBoundaryRef struct {
	ObjectDirectoryEntryRef string `json:"object_directory_entry_ref"`
	ObjectTruthSourceModule string `json:"object_truth_source_module"`
	TransactionRef          string `json:"transaction_ref"`
}

// BaseGateAdapterRequest carries only references needed for Base裁定; it never carries provider objects.
type BaseGateAdapterRequest struct {
	AdapterRequestRef   string        `json:"adapter_request_ref"`
	AdapterKind         string        `json:"adapter_kind"`
	CandidateRef        string        `json:"candidate_ref"`
	TransactionRef      string        `json:"transaction_ref"`
	DecisionRef         string        `json:"decision_ref"`
	ReceiptCandidateRef string        `json:"receipt_candidate_ref"`
	PolicySnapshotHash  string        `json:"policy_snapshot_hash"`
	EvidenceRefs        []string      `json:"evidence_refs"`
	ActorContext        *ActorContext `json:"actor_context"`
	RequestedAt         time.Time     `json:"requested_at"`
}

type MemoryWriteGateAdapterRequest struct {
	Gate                 BaseGateAdapterRequest `json:"gate"`
	FormalMemoryWriteRef FormalMemoryWriteRef   `json:"formal_memory_write_ref"`
	MemoryCandidateRef   string                 `json:"memory_candidate_ref"`
	DerivedContextRefs   []string               `json:"derived_context_refs,omitempty"`
}

type SendGateAdapterRequest struct {
	Gate          BaseGateAdapterRequest `json:"gate"`
	SendIntentRef string                 `json:"send_intent_ref"`
	ChannelRef    string                 `json:"channel_ref"`
	RecipientRef  string                 `json:"recipient_ref,omitempty"`
}

type ExecutionGateAdapterRequest struct {
	Gate               BaseGateAdapterRequest `json:"gate"`
	ExecutionIntentRef string                 `json:"execution_intent_ref"`
	RunID              string                 `json:"run_id"`
	Nonce              string                 `json:"nonce"`
	IdempotencyKey     IdempotencyKey         `json:"idempotency_key"`
	StopConditionRef   string                 `json:"stop_condition_ref"`
	ArtifactRefs       []string               `json:"artifact_refs,omitempty"`
}

func NewGateCandidateEnvelope(candidateRef, transactionRef, intentEventRef string, risk RiskClass, sideEffect SideEffectClass) *GateCandidateEnvelope {
	return &GateCandidateEnvelope{
		EnvelopeID:      candidateRef,
		CandidateRef:    candidateRef,
		TransactionRef:  transactionRef,
		IntentEventRef:  intentEventRef,
		SourceEventID:   intentEventRef,
		SourceModule:    "base-governance",
		PayloadHash:     stableRef("payload", candidateRef, transactionRef, intentEventRef),
		EvidenceRefs:    []string{intentEventRef},
		CandidateOnly:   true,
		NonFormal:       true,
		RiskClass:       risk,
		SideEffectClass: sideEffect,
	}
}

func NewPolicySnapshot(kernel PolicyKernel, policySet GatePolicySet, createdAt time.Time) (*PolicySnapshot, error) {
	if kernel.KernelRef == "" || kernel.Version == "" {
		return nil, errors.New("policy kernel ref and version are required")
	}
	if err := ValidateGatePolicySet(&policySet); err != nil {
		return nil, err
	}
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	raw := fmt.Sprintf("%s|%s|%s|%s|%d", kernel.KernelRef, kernel.Version, policySet.PolicySetRef, policySet.Version, createdAt.UnixNano())
	sum := sha256.Sum256([]byte(raw))
	hash := hex.EncodeToString(sum[:])
	return &PolicySnapshot{
		SnapshotRef:  fmt.Sprintf("policy_snapshot_%s", hash[:16]),
		SnapshotHash: hash,
		PolicySet: PolicySet{
			SetID:   policySet.PolicySetRef,
			Version: policySet.Version,
		},
		GatePolicySet: &policySet,
		KernelRef:     kernel.KernelRef,
		Version:       kernel.Version,
		IssuedBy:      PolicySnapshotIssuerBaseKernel,
		CreatedAt:     createdAt,
	}, nil
}

func NewCandidateIntakeEnvelope(req *GateRequest, policySnapshotRef string, receivedAt time.Time) (*CandidateIntakeEnvelope, error) {
	if err := ValidateGateRequest(req); err != nil {
		return nil, err
	}
	if policySnapshotRef == "" {
		return nil, errors.New("policy_snapshot_ref is required")
	}
	if receivedAt.IsZero() {
		receivedAt = time.Now().UTC()
	}
	return &CandidateIntakeEnvelope{
		IntakeRef:          stableRef("candidate_intake", effectiveGateCandidateRef(req), effectiveGateTransactionRef(req), effectiveIntentEventRef(req), policySnapshotRef, req.ReceiptRequirement),
		CandidateRef:       effectiveGateCandidateRef(req),
		TransactionRef:     effectiveGateTransactionRef(req),
		IntentEventRef:     effectiveIntentEventRef(req),
		EvidenceRefs:       effectiveEvidenceRefs(req),
		RiskClass:          req.Candidate.RiskClass,
		SideEffectClass:    req.Candidate.SideEffectClass,
		PolicySnapshotRef:  policySnapshotRef,
		ReceiptRequirement: req.ReceiptRequirement,
		ReceivedAt:         receivedAt,
	}, nil
}

func NewGateDecision(status GateDecisionStatus, reason string, req *GateRequest, snapshot *PolicySnapshot, trace []GateDecisionTrace, decidedAt time.Time) (*GateDecision, error) {
	if err := ValidateGateRequest(req); err != nil {
		return nil, err
	}
	if err := ValidatePolicySnapshot(snapshot); err != nil {
		return nil, err
	}
	if status == "" || reason == "" {
		return nil, errors.New("decision status and reason are required")
	}
	if decidedAt.IsZero() {
		decidedAt = time.Now().UTC()
	}
	if len(trace) == 0 {
		trace = []GateDecisionTrace{
			{
				GateName:           "GateOrchestrator",
				Status:             status,
				Reason:             reason,
				PolicySnapshotHash: snapshot.SnapshotHash,
				EvidenceRefs:       effectiveEvidenceRefs(req),
				CheckedAt:          decidedAt,
			},
		}
	}
	decisionRef := stableRef("gate_decision", effectiveGateCandidateRef(req), effectiveGateTransactionRef(req), snapshot.SnapshotHash, string(status), reason, req.ReceiptRequirement)
	return &GateDecision{
		DecisionRef:        decisionRef,
		Status:             status,
		Reason:             reason,
		GateTrace:          trace,
		DecidedAt:          decidedAt,
		ActorContext:       req.ActorContext,
		CandidateRef:       effectiveGateCandidateRef(req),
		TransactionRef:     effectiveGateTransactionRef(req),
		PolicySnapshotHash: snapshot.SnapshotHash,
		PolicySnapshot:     snapshot,
	}, nil
}

func NewGateReceiptCandidate(decision *GateDecision, req *GateRequest) (*GateReceiptCandidate, error) {
	if err := ValidateGateDecision(decision); err != nil {
		return nil, err
	}
	if err := ValidateGateRequest(req); err != nil {
		return nil, err
	}
	receiptRef := stableRef("gate_receipt_candidate", decision.DecisionRef, decision.CandidateRef, decision.TransactionRef)
	return &GateReceiptCandidate{
		ReceiptID:          receiptRef,
		TransactionRef:     decision.TransactionRef,
		CandidateRef:       decision.CandidateRef,
		DecisionRef:        decision.DecisionRef,
		PolicySnapshot:     decisionPolicyHash(decision),
		SourceModule:       req.Candidate.SourceModule,
		PayloadHash:        req.Candidate.PayloadHash,
		ActorContext:       decision.ActorContext,
		DecisionReason:     decision.Reason,
		ReceiptRequirement: req.ReceiptRequirement,
		EvidenceRefs:       effectiveEvidenceRefs(req),
		IdempotencyKey:     req.Candidate.IdempotencyKey,
	}, nil
}

func NewGateOrchestrationResult(decision *GateDecision, policySnapshotRef string, createdAt time.Time) (*GateOrchestrationResult, error) {
	if err := ValidateGateDecision(decision); err != nil {
		return nil, err
	}
	if policySnapshotRef == "" {
		return nil, errors.New("policy_snapshot_ref is required")
	}
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	return &GateOrchestrationResult{
		ResultRef:         stableRef("gate_orchestration_result", decision.DecisionRef, decision.CandidateRef, decision.TransactionRef, policySnapshotRef),
		CandidateRef:      decision.CandidateRef,
		TransactionRef:    decision.TransactionRef,
		Decision:          decision,
		PolicySnapshotRef: policySnapshotRef,
		Trace:             decision.GateTrace,
		CreatedAt:         createdAt,
	}, nil
}

func ValidateGateCandidateEnvelope(env *GateCandidateEnvelope) error {
	if env == nil {
		return errors.New("candidate envelope is required")
	}
	if effectiveCandidateRef(env) == "" {
		return errors.New("candidate_ref is required")
	}
	if env.TransactionRef == "" {
		return errors.New("transaction_ref is required")
	}
	if !env.CandidateOnly || !env.NonFormal {
		return errors.New("candidate_only and non_formal must both be true before formalization")
	}
	if env.RiskClass == "" {
		return errors.New("risk_class is required")
	}
	if env.SideEffectClass == "" {
		return errors.New("side_effect_class is required")
	}
	return nil
}

func ValidateGateRequest(req *GateRequest) error {
	if req == nil {
		return errors.New("gate request is required")
	}
	if err := ValidateGateCandidateEnvelope(req.Candidate); err != nil {
		return err
	}
	if req.ActorContext == nil || !hasActorIdentity(req.ActorContext) {
		return errors.New("actor context is required")
	}
	if effectiveGateTransactionRef(req) == "" {
		return errors.New("transaction spine transaction_ref is required")
	}
	if effectiveGateCandidateRef(req) == "" {
		return errors.New("candidate_ref is required")
	}
	if effectiveIntentEventRef(req) == "" {
		return errors.New("intent spine intent_event_ref is required")
	}
	if len(effectiveEvidenceRefs(req)) == 0 {
		return errors.New("evidence spine evidence_refs are required")
	}
	if req.ReceiptRequirement == "" {
		return errors.New("receipt_requirement is required")
	}
	if req.PolicySnapshotRef == "" {
		return errors.New("policy_snapshot_ref is required")
	}
	return nil
}

func ValidatePolicySnapshot(snapshot *PolicySnapshot) error {
	if snapshot == nil {
		return errors.New("policy snapshot is required")
	}
	if snapshot.SnapshotRef == "" || snapshot.SnapshotHash == "" {
		return errors.New("policy snapshot ref and hash are required")
	}
	if snapshot.IssuedBy != PolicySnapshotIssuerBaseKernel {
		return errors.New("policy snapshot must be issued by Base policy kernel")
	}
	if snapshot.KernelRef == "" || snapshot.Version == "" {
		return errors.New("policy snapshot kernel ref and version are required")
	}
	if snapshot.CreatedAt.IsZero() {
		return errors.New("policy snapshot created_at is required")
	}
	return nil
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

func ValidateGatePolicySet(set *GatePolicySet) error {
	if set == nil {
		return errors.New("gate policy set is required")
	}
	if set.PolicySetRef == "" || set.Version == "" {
		return errors.New("policy set ref and version are required")
	}
	if set.CreatedAt.IsZero() {
		return errors.New("policy set created_at is required")
	}
	return nil
}

func ValidateCandidateIntakeEnvelope(env *CandidateIntakeEnvelope) error {
	if env == nil {
		return errors.New("candidate intake envelope is required")
	}
	if env.IntakeRef == "" || env.CandidateRef == "" || env.TransactionRef == "" || env.IntentEventRef == "" || env.PolicySnapshotRef == "" {
		return errors.New("intake, candidate, transaction, intent, and policy refs are required")
	}
	if len(env.EvidenceRefs) == 0 {
		return errors.New("candidate intake evidence refs are required")
	}
	if env.RiskClass == "" || env.SideEffectClass == "" || env.ReceiptRequirement == "" {
		return errors.New("risk, side effect, and receipt requirement are required")
	}
	if env.ReceivedAt.IsZero() {
		return errors.New("received_at is required")
	}
	return nil
}

func ValidateGateDecision(decision *GateDecision) error {
	if decision == nil {
		return errors.New("gate decision is required")
	}
	if decision.DecisionRef == "" {
		return errors.New("decision_ref is required")
	}
	if decision.Status == "" {
		return errors.New("decision status is required")
	}
	if decision.CandidateRef == "" || decision.TransactionRef == "" {
		return errors.New("candidate_ref and transaction_ref are required")
	}
	if decision.DecidedAt.IsZero() {
		return errors.New("decided_at is required")
	}
	if decision.ActorContext == nil || !hasActorIdentity(decision.ActorContext) {
		return errors.New("actor context is required")
	}
	if decision.PolicySnapshotHash == "" {
		if decision.PolicySnapshot == nil || decision.PolicySnapshot.SnapshotHash == "" {
			return errors.New("policy snapshot hash is required")
		}
	}
	if decision.PolicySnapshot != nil {
		if err := ValidatePolicySnapshot(decision.PolicySnapshot); err != nil {
			return err
		}
	}
	if len(decision.GateTrace) == 0 {
		return errors.New("gate trace is required")
	}
	for _, trace := range decision.GateTrace {
		if trace.GateName == "" || trace.Status == "" || trace.PolicySnapshotHash == "" || trace.CheckedAt.IsZero() {
			return errors.New("gate trace entries require gate, status, policy hash, and checked_at")
		}
	}
	return nil
}

func ValidateGateReceiptCandidate(candidate *GateReceiptCandidate) error {
	if candidate == nil {
		return errors.New("gate receipt candidate is required")
	}
	if candidate.ReceiptID == "" || candidate.DecisionRef == "" || candidate.CandidateRef == "" || candidate.TransactionRef == "" {
		return errors.New("receipt_id, decision_ref, candidate_ref, and transaction_ref are required")
	}
	if candidate.PolicySnapshot == "" {
		return errors.New("policy_snapshot_hash is required")
	}
	if candidate.SourceModule == "" || candidate.PayloadHash == "" {
		return errors.New("source_module and payload_hash are required")
	}
	if candidate.ActorContext == nil || !hasActorIdentity(candidate.ActorContext) {
		return errors.New("actor context is required")
	}
	if candidate.ReceiptRequirement == "" {
		return errors.New("receipt_requirement is required")
	}
	if len(candidate.EvidenceRefs) == 0 {
		return errors.New("evidence_refs are required")
	}
	return nil
}

func ValidateGateOrchestrationResult(result *GateOrchestrationResult) error {
	if result == nil {
		return errors.New("gate orchestration result is required")
	}
	if result.ResultRef == "" || result.CandidateRef == "" || result.TransactionRef == "" || result.PolicySnapshotRef == "" {
		return errors.New("result, candidate, transaction, and policy refs are required")
	}
	if result.CreatedAt.IsZero() {
		return errors.New("created_at is required")
	}
	if len(result.Trace) == 0 {
		return errors.New("orchestration trace is required")
	}
	return ValidateGateDecision(result.Decision)
}

func ValidateModuleGateAdoptionStatus(status *ModuleGateAdoptionStatus) error {
	if status == nil {
		return errors.New("module gate adoption status is required")
	}
	if status.ModuleID == "" || status.ModuleName == "" || status.AdoptionState == "" {
		return errors.New("module id, module name, and adoption state are required")
	}
	if status.CheckedAt.IsZero() {
		return errors.New("checked_at is required")
	}
	if status.RequiresGateDecision && status.LastGatedRequestRef == "" && status.LastFailureReason == "" {
		return errors.New("module requiring gate decision must expose last gated request or failure reason")
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

func NewBaseFormalizationGrant(decision *GateDecision, targetObjectType, writeScope string, expiresAt time.Time) (*BaseFormalizationGrant, error) {
	if err := ValidateGateDecision(decision); err != nil {
		return nil, err
	}
	if decision.Status != GateDecisionAllow {
		return nil, errors.New("formalization grant requires allow decision")
	}
	if decision.ReceiptCandidate == nil {
		return nil, errors.New("formalization grant requires receipt candidate")
	}
	if err := ValidateGateReceiptCandidate(decision.ReceiptCandidate); err != nil {
		return nil, err
	}
	if targetObjectType == "" || writeScope == "" {
		return nil, errors.New("target object type and write scope are required")
	}
	if expiresAt.IsZero() || !expiresAt.After(decision.DecidedAt) {
		return nil, errors.New("grant expiration must be after decision time")
	}
	return &BaseFormalizationGrant{
		GrantRef:            fmt.Sprintf("base_grant_%s", decision.DecisionRef),
		DecisionRef:         decision.DecisionRef,
		CandidateRef:        decision.CandidateRef,
		TransactionRef:      decision.TransactionRef,
		ReceiptCandidateRef: decision.ReceiptCandidate.ReceiptID,
		PolicySnapshotHash:  decisionPolicyHash(decision),
		TargetObjectType:    targetObjectType,
		WriteScope:          writeScope,
		GrantedAt:           decision.DecidedAt,
		ExpiresAt:           expiresAt,
	}, nil
}

func ValidateCapabilityInvocationGateRequest(req *CapabilityInvocationGateRequest) error {
	if req == nil {
		return errors.New("capability invocation gate request is required")
	}
	if err := ValidateGateCandidateEnvelope(req.Candidate); err != nil {
		return err
	}
	if req.ActorContext == nil || !hasActorIdentity(req.ActorContext) {
		return errors.New("actor context is required")
	}
	if req.Invocation == nil {
		return errors.New("structured capability invocation candidate is required")
	}
	if req.Invocation.CapabilityID == "" || req.Invocation.InvocationRef == "" || req.Invocation.CandidateRef == "" || req.Invocation.TransactionRef == "" {
		return errors.New("capability invocation requires capability_id, invocation_ref, candidate_ref, and transaction_ref")
	}
	return nil
}

func ValidateRegistrySliceContextRef(ref *RegistrySliceContextRef) error {
	if ref == nil {
		return errors.New("registry slice context ref is required")
	}
	if ref.SliceRef == "" || ref.TransactionRef == "" || ref.IntentEventRef == "" || ref.ActorRef == "" || ref.ScopePolicyRef == "" || ref.AuditRef == "" {
		return errors.New("registry slice, transaction, intent, actor, scope policy, and audit refs are required")
	}
	return nil
}

func ValidateFormalMemoryWriteRef(ref *FormalMemoryWriteRef) error {
	if ref == nil {
		return errors.New("formal memory write ref is required")
	}
	if ref.MemoryWriteCandidateRef == "" || ref.FormalMemoryStoreRef == "" || ref.TransactionRef == "" {
		return errors.New("memory write candidate, formal memory store, and transaction refs are required")
	}
	return nil
}

func ValidateObjectDirectoryBoundaryRef(ref *ObjectDirectoryBoundaryRef) error {
	if ref == nil {
		return errors.New("object directory boundary ref is required")
	}
	if ref.ObjectDirectoryEntryRef == "" || ref.ObjectTruthSourceModule == "" || ref.TransactionRef == "" {
		return errors.New("object directory entry, truth source module, and transaction refs are required")
	}
	if ref.ObjectTruthSourceModule != "05-business-object-workbench" {
		return errors.New("object truth source must be 05-business-object-workbench")
	}
	return nil
}

func ValidateBaseAdapterRequest(req BaseGateAdapterRequest, expectedKind string) error {
	if req.AdapterRequestRef == "" || req.AdapterKind == "" {
		return errors.New("adapter request ref and kind are required")
	}
	if expectedKind != "" && req.AdapterKind != expectedKind {
		return fmt.Errorf("adapter kind mismatch: expected %s", expectedKind)
	}
	if req.CandidateRef == "" || req.TransactionRef == "" || req.DecisionRef == "" || req.ReceiptCandidateRef == "" {
		return errors.New("adapter request requires candidate, transaction, decision, and receipt refs")
	}
	if req.PolicySnapshotHash == "" || len(req.EvidenceRefs) == 0 || req.RequestedAt.IsZero() {
		return errors.New("adapter request requires policy snapshot, evidence refs, and requested_at")
	}
	if req.ActorContext == nil || !hasActorIdentity(req.ActorContext) {
		return errors.New("adapter request actor context is required")
	}
	encoded, _ := json.Marshal(req)
	if containsSecretish(string(encoded)) {
		return errors.New("adapter request must not expose raw secrets")
	}
	return nil
}

func ValidateMemoryWriteGateAdapterRequest(req *MemoryWriteGateAdapterRequest) error {
	if req == nil {
		return errors.New("memory write adapter request is required")
	}
	if err := ValidateBaseAdapterRequest(req.Gate, AdapterKindMemoryWrite); err != nil {
		return err
	}
	if err := ValidateFormalMemoryWriteRef(&req.FormalMemoryWriteRef); err != nil {
		return err
	}
	encoded, _ := json.Marshal(req)
	if containsSecretish(string(encoded)) {
		return errors.New("memory write adapter request must not expose raw secrets")
	}
	return nil
}

func ValidateSendGateAdapterRequest(req *SendGateAdapterRequest) error {
	if req == nil {
		return errors.New("send adapter request is required")
	}
	if err := ValidateBaseAdapterRequest(req.Gate, AdapterKindSend); err != nil {
		return err
	}
	if req.SendIntentRef == "" || req.ChannelRef == "" {
		return errors.New("send intent and channel refs are required")
	}
	encoded, _ := json.Marshal(req)
	if containsSecretish(string(encoded)) {
		return errors.New("send adapter request must not expose raw secrets")
	}
	return nil
}

func ValidateExecutionGateAdapterRequest(req *ExecutionGateAdapterRequest) error {
	if req == nil {
		return errors.New("execution adapter request is required")
	}
	if err := ValidateBaseAdapterRequest(req.Gate, AdapterKindExecution); err != nil {
		return err
	}
	if req.ExecutionIntentRef == "" || req.RunID == "" || req.Nonce == "" || req.IdempotencyKey == "" || req.StopConditionRef == "" {
		return errors.New("execution intent, run_id, nonce, idempotency_key, and stop condition are required")
	}
	encoded, _ := json.Marshal(req)
	if containsSecretish(string(encoded)) {
		return errors.New("execution adapter request must not expose raw secrets")
	}
	return nil
}

func hasActorIdentity(actor *ActorContext) bool {
	return actor.OwnerProfile != "" || actor.LocalIdentityRef != "" || actor.DeviceIdentityRef != ""
}

func effectiveCandidateRef(env *GateCandidateEnvelope) string {
	if env == nil {
		return ""
	}
	if env.CandidateRef != "" {
		return env.CandidateRef
	}
	return env.EnvelopeID
}

func effectiveGateCandidateRef(req *GateRequest) string {
	if req.CandidateRef != "" {
		return req.CandidateRef
	}
	return effectiveCandidateRef(req.Candidate)
}

func effectiveGateTransactionRef(req *GateRequest) string {
	if req.TransactionRef != "" {
		return req.TransactionRef
	}
	if req.Candidate != nil {
		return req.Candidate.TransactionRef
	}
	return ""
}

func effectiveIntentEventRef(req *GateRequest) string {
	if req.IntentEventRef != "" {
		return req.IntentEventRef
	}
	if req.Candidate != nil {
		if req.Candidate.IntentEventRef != "" {
			return req.Candidate.IntentEventRef
		}
		return req.Candidate.SourceEventID
	}
	return ""
}

func effectiveEvidenceRefs(req *GateRequest) []string {
	refs := append([]string{}, req.EvidenceRefs...)
	if req.Candidate != nil {
		refs = append(refs, req.Candidate.EvidenceRefs...)
		if req.Candidate.SourceEventID != "" {
			refs = append(refs, req.Candidate.SourceEventID)
		}
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(refs))
	for _, ref := range refs {
		if ref == "" {
			continue
		}
		if _, ok := seen[ref]; ok {
			continue
		}
		seen[ref] = struct{}{}
		out = append(out, ref)
	}
	return out
}

func decisionPolicyHash(decision *GateDecision) string {
	if decision.PolicySnapshotHash != "" {
		return decision.PolicySnapshotHash
	}
	if decision.PolicySnapshot != nil {
		return decision.PolicySnapshot.SnapshotHash
	}
	return ""
}

func containsSecretish(s string) bool {
	low := strings.ToLower(s)
	for _, marker := range []string{"raw_secret", "password", "token", "cookie", "private_key", "api_key"} {
		if strings.Contains(low, marker) {
			return true
		}
	}
	return false
}

func stableRef(prefix string, parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(sum[:])[:16])
}
