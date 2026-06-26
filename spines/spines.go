package spines

import "time"

// TransactionRef is the cross-module anchor for one owner-visible "thing".
// It may point to a formal TransactionObject, or to a candidate transaction
// while the object is still under governance review.
type TransactionRef struct {
	TransactionRef  string `json:"transaction_ref"`
	ObjectRef       string `json:"object_ref,omitempty"`
	PersonRef       string `json:"person_ref,omitempty"`
	SourceAuthority string `json:"source_authority,omitempty"`
	IdempotencyKey  string `json:"idempotency_key,omitempty"`
}

// IntentSource classifies every user, external, agent, pack, time, or system
// input before it fans out into module-specific candidates.
type IntentSource string

const (
	IntentSourceAgentSuggestionInput IntentSource = "agent_suggestion_input"
	IntentSourceUserMessageInput     IntentSource = "user_message_input"
	IntentSourceExternalMessageInput IntentSource = "external_message_input"
	IntentSourceExternalFileInput    IntentSource = "external_file_input"
	IntentSourceTaskFailureInput     IntentSource = "task_failure_input"
	IntentSourceTimeTriggerInput     IntentSource = "time_trigger_input"
	IntentSourcePackRunEventInput    IntentSource = "pack_run_event_input"
	IntentSourceSceneFlowEventInput  IntentSource = "scene_flow_event_input"
	IntentSourceGoalReviewInput      IntentSource = "goal_review_input"

	IntentSourceUserInput           IntentSource = IntentSourceUserMessageInput
	IntentSourceExternalMessage     IntentSource = IntentSourceExternalMessageInput
	IntentSourceTaskFailure         IntentSource = IntentSourceTaskFailureInput
	IntentSourcePackNodeTrigger     IntentSource = IntentSourcePackRunEventInput
	IntentSourceGoalReview          IntentSource = IntentSourceGoalReviewInput
	IntentSourceAgentSuggestion     IntentSource = IntentSourceAgentSuggestionInput
	IntentSourceTimeTrigger         IntentSource = IntentSourceTimeTriggerInput
	IntentSourceExternalSystemEvent IntentSource = IntentSourceExternalMessageInput
)

// IntentEvent is the normalized inbox event. It is non-formal until downstream
// candidates pass their module gates and receive receipts.
type IntentEvent struct {
	IntentEventID     string       `json:"intent_event_id"`
	Source            IntentSource `json:"source"`
	SourceRef         string       `json:"source_ref,omitempty"`
	RawInputSummary   string       `json:"raw_input_summary"`
	TransactionRef    string       `json:"transaction_ref"`
	ObjectRef         string       `json:"object_ref,omitempty"`
	PersonRef         string       `json:"person_ref,omitempty"`
	RelatedPersonRefs []string     `json:"related_person_refs,omitempty"`
	RelatedObjectRefs []string     `json:"related_object_refs,omitempty"`
	ActorRef          string       `json:"actor_ref,omitempty"`
	RiskHint          string       `json:"risk_hint,omitempty"`
	DesiredOutcome    string       `json:"desired_outcome,omitempty"`
	PayloadHash       string       `json:"payload_hash,omitempty"`
	CandidateTargets  []string     `json:"candidate_targets,omitempty"`
	CandidateOnly     bool         `json:"candidate_only"`
	NonFormal         bool         `json:"non_formal"`
	CreatedAt         time.Time    `json:"created_at"`
}

// SceneFlowRunRef binds a SceneFlowRun to the transaction, originating intent,
// pack version, current step, and any receipt evidence it emits.
type SceneFlowRunRef struct {
	SceneFlowRunRef string `json:"scene_flow_run_ref"`
	PackVersionRef  string `json:"pack_version_ref,omitempty"`
	TransactionRef  string `json:"transaction_ref"`
	StepRunRef      string `json:"step_run_ref,omitempty"`
	SourceEventID   string `json:"source_event_id,omitempty"`
	ReceiptRef      string `json:"receipt_ref,omitempty"`
}

// DispatchPlanRef is the handoff object from Task Governance to controlled
// gateways. It remains a candidate plan and never executes by itself.
type DispatchPlanRef struct {
	DispatchPlanRef string `json:"dispatch_plan_ref"`
	TaskRef         string `json:"task_ref"`
	TransactionRef  string `json:"transaction_ref"`
	DispatchType    string `json:"dispatch_type"`
	RequiresBase    bool   `json:"requires_base"`
	BaseDecisionRef string `json:"base_decision_ref,omitempty"`
	IdempotencyKey  string `json:"idempotency_key,omitempty"`
}

// ReceiptLink is the evidence edge that makes candidate -> formal transitions
// replayable through Receipt Ledger.
type ReceiptLink struct {
	ReceiptRef     string   `json:"receipt_ref"`
	TransactionRef string   `json:"transaction_ref"`
	CandidateRef   string   `json:"candidate_ref,omitempty"`
	DecisionRef    string   `json:"decision_ref,omitempty"`
	ActorRef       string   `json:"actor_ref,omitempty"`
	SourceModule   string   `json:"source_module,omitempty"`
	PayloadHash    string   `json:"payload_hash,omitempty"`
	EvidenceRefs   []string `json:"evidence_refs,omitempty"`
}
