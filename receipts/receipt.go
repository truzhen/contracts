package receipts

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
}
