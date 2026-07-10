package spines

// ImpactOperation is the object-edit vocabulary shared by pre-execution
// declarations (DeclaredImpact) and post-execution facts (ActualEdit).
// v1 covers the object-edit domain only: send/execute stay expressed by
// base.SideEffectClass and the candidate types (Owner ruling O-3, 2026-07-10);
// duplicating them here would create a second truth source.
type ImpactOperation string

const (
	ImpactOperationCreate ImpactOperation = "create"
	ImpactOperationModify ImpactOperation = "modify"
	ImpactOperationDelete ImpactOperation = "delete"
)

// DeclaredImpact is a Proposer-side pre-execution declaration of what a
// candidate intends to touch. It is evidence for gate evaluation and impact
// preview only — a declaration NEVER grants write authorization; that stays
// with Owner + Base Gate (impact-model-proposal-20260710.md §4).
type DeclaredImpact struct {
	ObjectType     string          `json:"object_type"`
	Operation      ImpactOperation `json:"operation"`
	ObjectRef      string          `json:"object_ref,omitempty"`
	AffectedFields []string        `json:"affected_fields,omitempty"`
	Description    string          `json:"description,omitempty"`
}

// ActualEdit is the post-execution fact recorded on a receipt: what was
// actually touched. ObjectRef is required — facts must point at an object.
// Reconciliation against DeclaredImpacts is a monitoring/audit concern and
// never rewrites either side.
type ActualEdit struct {
	ObjectType     string          `json:"object_type"`
	Operation      ImpactOperation `json:"operation"`
	ObjectRef      string          `json:"object_ref"`
	AffectedFields []string        `json:"affected_fields,omitempty"`
}
