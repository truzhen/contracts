package spines

// T12 阶段3：Intent Spine 五件套契约（IntentEvent / IntentInboxItem /
// IntentClassification / IntentToCandidateResult / IntentReceipt）。
//
// IntentEvent 本体在 spines.go；本文件收敛其余四件 + 分类 / 目标 / 状态枚举。
// 对应 JSON Schema：contracts/spines/intent-*.schema.json（contracts embed
// 模式，CI 用 jsonschema/v6 真实编译校验）。
//
// 纪律：
//   - 所有外部输入先落 IntentEvent，再分流为八类候选之一；候选必须携带
//     intent_event_id（可反查追溯）。
//   - 分类只产候选路由，不产正式输出（FormalTask / FormalMemory /
//     SendResult / ExecutionResult 永远在 BlockedFormalOutputs）。
//   - 无法分类的 IntentEvent 进 inbox 待 Owner 处置（pending_owner），
//     不得静默丢弃。

import (
	"time"

	"github.com/truzhen/contracts/candidates"
)

// IntentClassificationType is the closed classification vocabulary. unknown
// means the event could not be classified and must wait for Owner disposition.
type IntentClassificationType string

const (
	IntentClassificationTaskRequest          IntentClassificationType = "task_request"
	IntentClassificationCommunicationRequest IntentClassificationType = "communication_request"
	IntentClassificationMemoryRequest        IntentClassificationType = "memory_request"
	IntentClassificationExecutionRequest     IntentClassificationType = "execution_request"
	IntentClassificationObjectUpdateRequest  IntentClassificationType = "object_update_request"
	IntentClassificationSceneFlowTrigger     IntentClassificationType = "scene_flow_trigger"
	IntentClassificationGoalReviewTrigger    IntentClassificationType = "goal_review_trigger"
	IntentClassificationUnknown              IntentClassificationType = "unknown"
)

// IntentCandidateTarget enumerates the EIGHT candidate types an IntentEvent
// may fan out into. No ninth type may be invented; formal outputs are not
// targets.
type IntentCandidateTarget string

const (
	IntentTargetTaskCandidate                 IntentCandidateTarget = "TaskCandidate"
	IntentTargetMemoryRequestCandidate        IntentCandidateTarget = "MemoryRequestCandidate"
	IntentTargetCommunicationDraftCandidate   IntentCandidateTarget = "CommunicationDraftCandidate"
	IntentTargetExecutionIntentCandidate      IntentCandidateTarget = "ExecutionIntentCandidate"
	IntentTargetBusinessObjectCandidate       IntentCandidateTarget = "BusinessObjectCandidate"
	IntentTargetSceneFlowRunCandidate         IntentCandidateTarget = "SceneFlowRunCandidate"
	IntentTargetCapabilityInvocationCandidate IntentCandidateTarget = "CapabilityInvocationCandidate"
	IntentTargetPackCandidate                 IntentCandidateTarget = "PackCandidate"
)

// IntentInboxItemStatus is the inbox disposition state for one IntentEvent.
//
//   - routed: classified and fanned out into candidates (refs recorded).
//   - pending_owner: could not be classified; waits for Owner disposition.
//     It must never be silently dropped.
//   - suppressed: deliberately produced no candidate cards (e.g. secretary
//     small talk); still recorded for traceability.
type IntentInboxItemStatus string

const (
	IntentInboxItemRouted       IntentInboxItemStatus = "routed"
	IntentInboxItemPendingOwner IntentInboxItemStatus = "pending_owner"
	IntentInboxItemSuppressed   IntentInboxItemStatus = "suppressed"
)

// IntentInboxItem is the per-event inbox row: every intake produces exactly
// one item keyed by intent_event_id (idempotent replays return the same item).
type IntentInboxItem struct {
	IntentInboxItemID   string                     `json:"intent_inbox_item_id"`
	IntentEventID       string                     `json:"intent_event_id"`
	TransactionRef      string                     `json:"transaction_ref"`
	Source              IntentSource               `json:"source"`
	Status              IntentInboxItemStatus      `json:"status"`
	StatusReason        string                     `json:"status_reason,omitempty"`
	ClassificationTypes []IntentClassificationType `json:"classification_types,omitempty"`
	CandidateRefs       []string                   `json:"candidate_refs,omitempty"`
	ReceiptCandidateRef string                     `json:"receipt_candidate_ref,omitempty"`
	CandidateOnly       bool                       `json:"candidate_only"`
	NonFormal           bool                       `json:"non_formal"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
}

// IntentClassification is the candidate-routing decision for one IntentEvent.
// It never grants execution; requires_base_gate and receipt_required stay true
// for every downstream formal action.
type IntentClassification struct {
	IntentEventID       string                     `json:"intent_event_id"`
	TransactionRef      string                     `json:"transaction_ref"`
	ClassificationTypes []IntentClassificationType `json:"classification_types"`
	CandidateTargets    []IntentCandidateTarget    `json:"candidate_targets,omitempty"`
	CandidateOnly       bool                       `json:"candidate_only"`
	NonFormal           bool                       `json:"non_formal"`
	RequiresBaseGate    bool                       `json:"requires_base_gate"`
	ReceiptRequired     bool                       `json:"receipt_required"`
	RiskClass           string                     `json:"risk_class,omitempty"`
	ModelGatewayRef     string                     `json:"model_gateway_ref,omitempty"`
}

// IntentToCandidateResult is the fan-out product: one of the eight candidate
// types per target, each envelope carrying intent_event_id (payload key
// "intent_event_id" + SourceEventID) so every candidate is traceable back to
// its IntentEvent.
type IntentToCandidateResult struct {
	IntentEventID        string                          `json:"intent_event_id"`
	TransactionRef       string                          `json:"transaction_ref"`
	CandidateRefs        []string                        `json:"candidate_refs,omitempty"`
	Candidates           []*candidates.CandidateEnvelope `json:"candidates,omitempty"`
	CandidateOnly        bool                            `json:"candidate_only"`
	NonFormal            bool                            `json:"non_formal"`
	BlockedFormalOutputs []string                        `json:"blocked_formal_outputs,omitempty"`
}

// IntentReceipt is the candidate-domain receipt of one intake/route step. It
// is evidence material (ref-only); formal receipts are minted by Base Gate +
// Receipt Ledger when a candidate later formalizes.
type IntentReceipt struct {
	ReceiptCandidateRef string    `json:"receipt_candidate_ref"`
	IntentEventID       string    `json:"intent_event_id"`
	TransactionRef      string    `json:"transaction_ref"`
	Action              string    `json:"action"`
	Status              string    `json:"status"`
	CandidateRefs       []string  `json:"candidate_refs,omitempty"`
	EvidenceRefs        []string  `json:"evidence_refs,omitempty"`
	CandidateOnly       bool      `json:"candidate_only"`
	NonFormal           bool      `json:"non_formal"`
	CreatedAt           time.Time `json:"created_at"`
}

// IntentReceipt actions (closed set, mirrors intent-receipt.schema.json).
const (
	IntentReceiptActionRouted       = "intent_routed"
	IntentReceiptActionPendingOwner = "intent_pending_owner"
	IntentReceiptActionSuppressed   = "intent_suppressed"
)
