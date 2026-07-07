package market

type ProviderReadiness struct {
	Status       string `json:"status"`
	PackRegistry string `json:"pack_registry,omitempty"`
	Commerce     string `json:"commerce,omitempty"`
	Payment      string `json:"payment,omitempty"`
	License      string `json:"license,omitempty"`
}

type MarketCatalogProduct struct {
	ProductID           string            `json:"product_id"`
	Name                string            `json:"name"`
	ProductKind         ProductKind       `json:"product_kind"`
	Version             string            `json:"version"`
	SellerUserID        string            `json:"seller_user_id"`
	AuthorID            string            `json:"author_id"`
	Status              ProductStatus     `json:"status"`
	PriceWeeklyMinor    int64             `json:"price_weekly_minor"`
	PriceMonthlyMinor   int64             `json:"price_monthly_minor"`
	PriceYearlyMinor    int64             `json:"price_yearly_minor"`
	PriceLifetimeMinor  int64             `json:"price_lifetime_minor,omitempty"`
	Currency            string            `json:"currency"`
	IsFree              bool              `json:"is_free"`
	RegionCode          string            `json:"region_code"`
	RegionDisplayName   string            `json:"region_display_name"`
	ArchTags            []string          `json:"arch_tags,omitempty"`
	AuthorShareRatioBps int64             `json:"author_share_ratio_bps,omitempty"`
	Source              string            `json:"source,omitempty"`
	HasCommerce         bool              `json:"has_commerce,omitempty"`
	HasForgejo          bool              `json:"has_forgejo,omitempty"`
	DownloadPath        string            `json:"download_path,omitempty"`
	PackType            string            `json:"pack_type,omitempty"`
	Owner               string            `json:"owner,omitempty"`
	CreatedAt           string            `json:"created_at,omitempty"`
	ProviderReadiness   ProviderReadiness `json:"provider_readiness,omitempty"`
}

type MarketEntitlement struct {
	EntitlementID  string   `json:"entitlement_id"`
	EntitlementRef string   `json:"entitlement_ref,omitempty"`
	OrderRef       string   `json:"order_ref,omitempty"`
	ProductID      string   `json:"product_id"`
	ProductKind    string   `json:"product_kind"`
	LicenseType    string   `json:"license_type"`
	IsPerpetual    bool     `json:"is_perpetual,omitempty"`
	ValidFrom      string   `json:"valid_from"`
	ValidTo        string   `json:"valid_to"`
	Status         string   `json:"status"`
	SourceOrderID  string   `json:"source_order_id"`
	DeviceLimit    int64    `json:"device_limit"`
	ReceiptRef     string   `json:"receipt_ref,omitempty"`
	EvidenceRefs   []string `json:"evidence_refs,omitempty"`
}

type MarketCheckoutResult struct {
	OrderID              string                 `json:"order_id"`
	OrderRef             string                 `json:"order_ref,omitempty"`
	PaymentRef           string                 `json:"payment_ref,omitempty"`
	EntitlementRef       string                 `json:"entitlement_ref,omitempty"`
	AmountMinor          int64                  `json:"amount_minor,omitempty"`
	Currency             string                 `json:"currency,omitempty"`
	Channel              string                 `json:"channel,omitempty"`
	CodeURL              string                 `json:"code_url"`
	PaymentMode          string                 `json:"payment_mode"`
	NoRealPayment        bool                   `json:"no_real_payment,omitempty"`
	TransactionRef       string                 `json:"transaction_ref,omitempty"`
	Detail               string                 `json:"detail,omitempty"`
	PriceSnapshot        map[string]interface{} `json:"price_snapshot,omitempty"`
	OrderCandidate       map[string]interface{} `json:"order_candidate,omitempty"`
	PaymentCandidate     map[string]interface{} `json:"payment_candidate,omitempty"`
	EntitlementCandidate map[string]interface{} `json:"entitlement_candidate,omitempty"`
	ReceiptRef           string                 `json:"receipt_ref,omitempty"`
	EvidenceRefs         []string               `json:"evidence_refs,omitempty"`
}

type MarketOrderStatus struct {
	OrderID        string   `json:"order_id"`
	OrderRef       string   `json:"order_ref,omitempty"`
	PaymentRef     string   `json:"payment_ref,omitempty"`
	EntitlementRef string   `json:"entitlement_ref,omitempty"`
	Status         string   `json:"status"`
	AmountMinor    int64    `json:"amount_minor,omitempty"`
	Currency       string   `json:"currency,omitempty"`
	EntitlementID  string   `json:"entitlement_id,omitempty"`
	PaidAt         string   `json:"paid_at,omitempty"`
	ReceiptRef     string   `json:"receipt_ref,omitempty"`
	EvidenceRefs   []string `json:"evidence_refs,omitempty"`
}

type MarketLocalGateCheckResult struct {
	DecisionID      string   `json:"decision_id,omitempty"`
	Subject         string   `json:"subject,omitempty"`
	BuyerID         string   `json:"buyer_id,omitempty"`
	ProductID       string   `json:"product_id,omitempty"`
	Status          string   `json:"status"`
	Allowed         bool     `json:"allowed,omitempty"`
	Reason          string   `json:"reason,omitempty"`
	EntitlementID   string   `json:"entitlement_id,omitempty"`
	ValidUntil      string   `json:"valid_until,omitempty"`
	ReadOnlyAllowed bool     `json:"read_only_allowed,omitempty"`
	ReceiptRef      string   `json:"receipt_ref,omitempty"`
	EvidenceRefs    []string `json:"evidence_refs,omitempty"`
	CheckedAt       string   `json:"checked_at,omitempty"`
}

type PackInstallResult struct {
	OK            bool     `json:"ok,omitempty"`
	ProductID     string   `json:"product_id,omitempty"`
	PackRef       string   `json:"pack_ref,omitempty"`
	Version       string   `json:"version,omitempty"`
	InstallStatus string   `json:"install_status,omitempty"`
	ReceiptRef    string   `json:"receipt_ref,omitempty"`
	EvidenceRefs  []string `json:"evidence_refs,omitempty"`
	Error         string   `json:"error,omitempty"`
	Detail        string   `json:"detail,omitempty"`
}

type PackExportUploadPrefill struct {
	PackID       string `json:"pack_id,omitempty"`
	PackRef      string `json:"pack_ref,omitempty"`
	Version      string `json:"version,omitempty"`
	Kind         string `json:"kind,omitempty"`
	BundleZipURL string `json:"bundle_zip_url,omitempty"`
}

type PackExportBundle struct {
	Status        string                  `json:"status"`
	PackRef       string                  `json:"pack_ref,omitempty"`
	Version       string                  `json:"version,omitempty"`
	BundleID      string                  `json:"bundle_id,omitempty"`
	BundleZipURL  string                  `json:"bundle_zip_url"`
	SHA256        string                  `json:"sha256"`
	SizeBytes     int64                   `json:"size_bytes,omitempty"`
	UploadPrefill PackExportUploadPrefill `json:"upload_prefill,omitempty"`
}
