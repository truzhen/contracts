package candidates

type TaskCandidate struct {
	CandidateEnvelope
	TaskName string `json:"task_name"`
	// TaskKind declares the task's 07 TaskType at the cross-repo contract layer:
	// "immediate" (即时任务, atomic one-shot), "stage" (阶段任务, a milestone carrying
	// a phase gate), or "scheduled". Empty defaults to "immediate" — additive and
	// backward compatible (older packs that omit it stay valid). It mirrors the
	// truzhenos taskgovernance domain TaskType values; the domain stays the truth
	// source for task STATE, while this field is the authoring-time declaration
	// consumers (packs / cloud / client) read from the contract.
	TaskKind string `json:"task_kind,omitempty"`
	// ParentTaskRef links a child immediate task to its owning 阶段任务 (stage
	// FormalTask); PhaseRef optionally names the phase. Both empty for top-level
	// tasks. Used by the 阶段门 (phase gate) child aggregation. The value is the
	// parent stage's FormalTask id verbatim (truzhenos mints it as
	// "formal_task_<hex>"); the gate matches children by this exact id, so a
	// differently-shaped ref would silently disable the gate — always fill the
	// real parent id, never a decorative placeholder.
	ParentTaskRef string `json:"parent_task_ref,omitempty"`
	PhaseRef      string `json:"phase_ref,omitempty"`
}
