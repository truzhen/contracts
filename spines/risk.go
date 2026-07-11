package spines

// RiskEscalationPath is the escalation a pack's judgment policy (判事策略)
// requests for a matched formal action. v1 keeps exactly two values: "none"
// (record-only) and "owner_gate" (raise Allow → PendingOwner). A declaration
// NEVER grants or removes authority — the ruling stays with Owner + Base Gate
// (统一决策表 #11; KWeaver RiskType 的反例=规格无消费者，故本类型与 gate floor
// 同轮接线，见 backlog-early-trigger-plan-20260711.md W2 卡).
type RiskEscalationPath string

const (
	RiskEscalationNone      RiskEscalationPath = "none"
	RiskEscalationOwnerGate RiskEscalationPath = "owner_gate"
)

// DeclaredRiskType is the gate-facing projection of one matched pack
// risk-type declaration. Producers resolve the pack's declared risk_types
// against the concrete action type and attach only the matched entries;
// like DeclaredImpact this is evidence for gate evaluation, never an
// authorization. Conservative tier: absence changes no gate behavior.
type DeclaredRiskType struct {
	RiskTypeID          string             `json:"risk_type_id"`
	TriggerActionType   string             `json:"trigger_action_type"`
	EvidenceRequirement string             `json:"evidence_requirement,omitempty"`
	EscalationPath      RiskEscalationPath `json:"escalation_path"`
	Definition          string             `json:"definition,omitempty"`
}
