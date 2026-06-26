package candidates

// CitedKnowledgeRef 是建议 / 质询候选引用的一条结构化、可反查的正式知识引用。
//
// 它来自真实知识检索（FormalKnowledge slice 命中），不是模型自然语言文本里
// 抽取的字符串：每条都能在 03 回执账本经 KnowledgeRef / ReceiptRef 反查到具体
// 法条版本（Version）+ 来源（Source）+ 检索时点（AsOf）。复用 09 知识库 D2 的
// 回溯字段语义（authority / verification / effectivity / source）。
//
// VerificationStatus 继承知识库的核验状态：pending_human_review 的法条引用必须
// 在建议里继承「待核验，不得直接作为正式文书依据」caveat（Proposer 不得自断
// 法律效力）。
type CitedKnowledgeRef struct {
	KnowledgeRef       string `json:"knowledge_ref"`
	Title              string `json:"title,omitempty"`
	Version            int64  `json:"version"`
	Source             string `json:"source,omitempty"`
	AsOf               string `json:"as_of,omitempty"`
	VerificationStatus string `json:"verification_status"`
	ReceiptRef         string `json:"receipt_ref,omitempty"`
}

// AdviceCandidate 是多角色场景流程中的建议 / 质询 / 对照候选。
// 它只表达可审计的角色观点，不执行动作、不生成正式结论。
//
// CitedKnowledgeRefs 承载该意见结构化引用的法条 / 知识（来自真实检索，非 NL
// 抽取），使「处罚意见 → 引用法条 → 法条版本 → 来源」在 03 可反查（依赖 C3）。
type AdviceCandidate struct {
	CandidateEnvelope
	SlotRef            string              `json:"slot_ref"`
	AgentRef           string              `json:"agent_ref"`
	RolePackRef        string              `json:"role_pack_ref"`
	AdviceKind         string              `json:"advice_kind"`
	Reasoning          string              `json:"reasoning"`
	Summary            string              `json:"summary"`
	EvidenceRefs       []string            `json:"evidence_refs,omitempty"`
	CitedKnowledgeRefs []CitedKnowledgeRef `json:"cited_knowledge_refs,omitempty"`
	KnowledgeCaveat    string              `json:"knowledge_caveat,omitempty"`
	// 知识挂载审计（Task 6，前端 / 03 回执展示用）：KnowledgeMountStatus="disabled"
	// 表示该场景包知识 scope 无 active 挂载（停用后引用被阻断）、"active" 表示来自
	// active 挂载；非 pack-scoped 时留空。KnowledgeMountRef 可空。
	KnowledgeMountRef    string `json:"knowledge_mount_ref,omitempty"`
	KnowledgeMountStatus string `json:"knowledge_mount_status,omitempty"`
	TransactionRef       string `json:"transaction_ref,omitempty"`
	SourceEventID        string `json:"source_event_id,omitempty"`
	PackVersionRef       string `json:"pack_version_ref,omitempty"`
	SceneFlowRunRef      string `json:"scene_flow_run_ref,omitempty"`
	StepRunRef           string `json:"step_run_ref,omitempty"`
}

func NewAdviceCandidate(slotRef, agentRef, rolePackRef, kind, summary string) *AdviceCandidate {
	return &AdviceCandidate{
		CandidateEnvelope: *NewCandidateEnvelope(nil),
		SlotRef:           slotRef,
		AgentRef:          agentRef,
		RolePackRef:       rolePackRef,
		AdviceKind:        kind,
		Summary:           summary,
	}
}
