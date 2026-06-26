package receipts

type AuditEnvelope struct {
	AuditID   string `json:"audit_id"`
	ReceiptID string `json:"receipt_id"`
}
