package cloud

import "time"

type PaymentProvider string

const (
	PaymentProviderWechat     PaymentProvider = "wechat"
	PaymentProviderLakala     PaymentProvider = "lakala"
	PaymentProviderShouqianba PaymentProvider = "shouqianba"
	PaymentProviderManual     PaymentProvider = "manual"
)

type PaymentOrderStatus string

const (
	PaymentOrderPending   PaymentOrderStatus = "pending"
	PaymentOrderPaid      PaymentOrderStatus = "paid"
	PaymentOrderFailed    PaymentOrderStatus = "failed"
	PaymentOrderCancelled PaymentOrderStatus = "cancelled"
	PaymentOrderRefunded  PaymentOrderStatus = "refunded"
)

type PaymentOrder struct {
	OrderRef       string             `json:"order_ref"`
	Provider       PaymentProvider    `json:"provider"`
	AccountRef     CloudAccountRef    `json:"account_ref,omitempty"`
	EntitlementRef string             `json:"entitlement_ref,omitempty"`
	AmountCents    int64              `json:"amount_cents"`
	Currency       string             `json:"currency"`
	Status         PaymentOrderStatus `json:"status"`
	ProviderRef    string             `json:"provider_ref,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      *time.Time         `json:"updated_at,omitempty"`
}

type PaymentWebhook struct {
	WebhookRef     string          `json:"webhook_ref"`
	Provider       PaymentProvider `json:"provider"`
	OrderRef       string          `json:"order_ref"`
	EventID        string          `json:"event_id"`
	SignatureRef   string          `json:"signature_ref,omitempty"`
	PayloadDigest  string          `json:"payload_digest,omitempty"`
	ReceivedAt     time.Time       `json:"received_at"`
	ReceiptRef     string          `json:"receipt_ref,omitempty"`
	ProcessingNote string          `json:"processing_note,omitempty"`
}

type PaymentReceiptRef struct {
	ReceiptRef string          `json:"receipt_ref"`
	OrderRef   string          `json:"order_ref"`
	Provider   PaymentProvider `json:"provider"`
}
