package receipts

import "github.com/truzhen/contracts/spines"

// ReceiptEnvelope 是回执链的基础包装。
// 根据纪律，Receipt 不等于 Memory，Sensitive payload 必须 SecureStore。
type ReceiptEnvelope struct {
	ReceiptID      string   `json:"receipt_id"`
	TransactionRef string   `json:"transaction_ref,omitempty"`
	CandidateRef   string   `json:"candidate_ref,omitempty"`
	DecisionRef    string   `json:"decision_ref,omitempty"`
	EvidenceRefs   []string `json:"evidence_refs,omitempty"`
	Sequence       int64    `json:"sequence"`
	PreviousHash   string   `json:"previous_hash"`
	PayloadHash    string   `json:"payload_hash"`
	Hash           string   `json:"hash"`

	// ActualEdits 是执行后事实（实际编辑了哪些对象），与候选侧
	// declared_impacts 两端分离：Receipt 只写事实，声明不构成授权。
	ActualEdits []spines.ActualEdit `json:"actual_edits,omitempty"`
}
