package contracts

type KnowledgeMountStatus string

const (
	KnowledgeMountPending  KnowledgeMountStatus = "pending"
	KnowledgeMountActive   KnowledgeMountStatus = "active"
	KnowledgeMountDisabled KnowledgeMountStatus = "disabled"
	KnowledgeMountBlocked  KnowledgeMountStatus = "blocked"
)

type KnowledgeScopeDeclaration struct {
	ScopeRef       string   `json:"scope_ref"`
	DisplayName    string   `json:"display_name"`
	Description    string   `json:"description,omitempty"`
	SceneRef       string   `json:"scene_ref,omitempty"`
	MountPolicy    string   `json:"mount_policy"`
	KnowledgeKinds []string `json:"knowledge_kinds"`
	Tags           []string `json:"tags,omitempty"`
	Required       bool     `json:"required,omitempty"`
}

type KnowledgeMountReadModel struct {
	MountRef           string               `json:"mount_ref"`
	OwnerID            string               `json:"owner_id"`
	PackRef            string               `json:"pack_ref"`
	PackVersionRef     string               `json:"pack_version_ref"`
	SceneRef           string               `json:"scene_ref,omitempty"`
	KnowledgeScopeRef  string               `json:"knowledge_scope_ref"`
	DisplayName        string               `json:"display_name"`
	Status             KnowledgeMountStatus `json:"status"`
	KnowledgeKinds     []string             `json:"knowledge_kinds"`
	KnowledgeRefs      []string             `json:"knowledge_refs,omitempty"`
	EnabledReceiptRef  string               `json:"enabled_receipt_ref,omitempty"`
	DisabledReceiptRef string               `json:"disabled_receipt_ref,omitempty"`
	LastReceiptRef     string               `json:"last_receipt_ref,omitempty"`
	BlockedReason      string               `json:"blocked_reason,omitempty"`
}
