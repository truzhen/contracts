package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	AdapterKindMemoryWrite = "memory_write"
	AdapterKindSend        = "send"
	AdapterKindExecution   = "execution"
)

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

func containsSecretish(s string) bool {
	low := strings.ToLower(s)
	for _, marker := range []string{"raw_secret", "password", "token", "cookie", "private_key", "api_key"} {
		if strings.Contains(low, marker) {
			return true
		}
	}
	return false
}
