package market

import "time"

// ProductKind 是云市场商品类型枚举。Pack 上传只允许其中的 Pack 类型；
// 服务端商品目录仍可承载 truzhen_software / feature 等非 Pack 商品。
type ProductKind string

const (
	ProductKindTruzhenSoftware ProductKind = "truzhen_software"
	ProductKindScenePack       ProductKind = "scene_pack"
	ProductKindRolePack        ProductKind = "role_pack"
	ProductKindCapabilityPack  ProductKind = "capability_pack"
	ProductKindSkillBundle     ProductKind = "skill_bundle"
	ProductKindFeature         ProductKind = "feature"
)

// ProductStatus 是商品和作者上传列表可投射的商品状态。
type ProductStatus string

const (
	ProductStatusDraft              ProductStatus = "draft"
	ProductStatusPendingReview      ProductStatus = "pending_review"
	ProductStatusListed             ProductStatus = "listed"
	ProductStatusDelisted           ProductStatus = "delisted"
	ProductStatusRejected           ProductStatus = "rejected"
	ProductStatusPriceChangePending ProductStatus = "price_change_pending_review"
)

// TrustVerifyStatus 是 Pack 上传签名 / 信任验证结果。
type TrustVerifyStatus string

const (
	TrustVerifyStatusVerified        TrustVerifyStatus = "verified"
	TrustVerifyStatusUnverified      TrustVerifyStatus = "unverified"
	TrustVerifyStatusFailed          TrustVerifyStatus = "failed"
	TrustVerifyStatusProviderMissing TrustVerifyStatus = "provider_missing"
)

// AuthorCertificationStatus 是作者实名认证生命周期状态。
type AuthorCertificationStatus string

const (
	AuthorCertificationStatusNone           AuthorCertificationStatus = "none"
	AuthorCertificationStatusPendingPayment AuthorCertificationStatus = "pending_payment"
	AuthorCertificationStatusPendingReview  AuthorCertificationStatus = "pending_review"
	AuthorCertificationStatusCertified      AuthorCertificationStatus = "certified"
	AuthorCertificationStatusRejected       AuthorCertificationStatus = "rejected"
	AuthorCertificationStatusRevoked        AuthorCertificationStatus = "revoked"
)

// AuthorCertificationReadModel 是作者 / 运营可读认证投射；只承载脱敏字段。
type AuthorCertificationReadModel struct {
	Status                 AuthorCertificationStatus `json:"status"`
	FeeMinor               int64                     `json:"fee_minor"`
	Currency               string                    `json:"currency"`
	OrderID                string                    `json:"order_id"`
	RegionCode             string                    `json:"region_code"`
	RegionDisplay          string                    `json:"region_display"`
	GrantedAt              *time.Time                `json:"granted_at"`
	RequestedAt            time.Time                 `json:"requested_at"`
	ApplicantName          string                    `json:"applicant_name"`
	ContactMasked          string                    `json:"contact_masked"`
	PayoutAccountMasked    string                    `json:"payout_account_masked"`
	IDCardMasked           string                    `json:"id_card_masked"`
	Address                string                    `json:"address"`
	Gender                 string                    `json:"gender"`
	ValidUntil             string                    `json:"valid_until"`
	ResidenceAddressMasked string                    `json:"residence_address_masked"`
	BankName               string                    `json:"bank_name"`
	ReviewedAt             *time.Time                `json:"reviewed_at"`
	RejectReason           string                    `json:"reject_reason"`
	ReceiptRef             string                    `json:"receipt_ref"`
}

// AuthorCertificationRegistration 是作者实名认证登记请求。照片不上云；OCR 明文字段
// 只作为提交输入，云端应加密 / 脱敏后持久化。
type AuthorCertificationRegistration struct {
	ApplicantName    string `json:"applicant_name"`
	IDCard           string `json:"id_card"`
	Gender           string `json:"gender"`
	IDAddress        string `json:"id_address"`
	ValidUntil       string `json:"valid_until"`
	Contact          string `json:"contact"`
	ResidenceAddress string `json:"residence_address"`
	PayoutAccount    string `json:"payout_account"`
	BankName         string `json:"bank_name"`
	RegionCode       string `json:"region_code"`
	RegionDisplay    string `json:"region_display"`
}

// IDCardOCRFields 是本地 OCR 候选结果字段；照片不得由本契约要求上传云端。
type IDCardOCRFields struct {
	IDCard     string `json:"id_card"`
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	Nation     string `json:"nation"`
	Birth      string `json:"birth"`
	Address    string `json:"address"`
	ValidUntil string `json:"valid_until"`
}

type IDCardOCRResult struct {
	Status         string           `json:"status"`
	Fields         *IDCardOCRFields `json:"fields"`
	PhotoPersisted bool             `json:"photo_persisted"`
	Detail         string           `json:"detail"`
	Hint           string           `json:"hint"`
}

// PackRevenueLineReadModel 是作者收益投射中的逐包收益行。
type PackRevenueLineReadModel struct {
	ProductID        string   `json:"product_id"`
	ProductName      string   `json:"product_name"`
	RegionCode       string   `json:"region_code"`
	ArchTags         []string `json:"arch_tags"`
	ShareRatioBps    int64    `json:"share_ratio_bps"`
	PaidOrderCount   int64    `json:"paid_order_count"`
	GrossMinor       int64    `json:"gross_minor"`
	AuthorShareMinor int64    `json:"author_share_minor"`
}

// AuthorRevenueReadModel 是作者收益只读投射；真相来自云端已支付订单、
// 不可变分成快照、税费快照和提现状态机。
type AuthorRevenueReadModel struct {
	AuthorID               string                     `json:"author_id"`
	Currency               string                     `json:"currency"`
	GrossAuthorShareMinor  int64                      `json:"gross_author_share_minor"`
	TaxRateBps             int64                      `json:"tax_rate_bps"`
	TaxMinor               int64                      `json:"tax_minor"`
	NetAccruedMinor        int64                      `json:"net_accrued_minor"`
	WithdrawnMinor         int64                      `json:"withdrawn_minor"`
	PendingWithdrawalMinor int64                      `json:"pending_withdrawal_minor"`
	WithdrawableMinor      int64                      `json:"withdrawable_minor"`
	WithdrawFeeBps         int64                      `json:"withdraw_fee_bps"`
	WithdrawMinAmountMinor int64                      `json:"withdraw_min_amount_minor"`
	PerPack                []PackRevenueLineReadModel `json:"per_pack"`
	GeneratedAt            time.Time                  `json:"generated_at"`
}

// WithdrawalStatus 是作者提现状态机的封闭字符串集合。
type WithdrawalStatus string

const (
	WithdrawalStatusRequested          WithdrawalStatus = "REQUESTED"
	WithdrawalStatusUnderReview        WithdrawalStatus = "UNDER_REVIEW"
	WithdrawalStatusApproved           WithdrawalStatus = "APPROVED"
	WithdrawalStatusTransferRegistered WithdrawalStatus = "TRANSFER_REGISTERED"
	WithdrawalStatusCompleted          WithdrawalStatus = "COMPLETED"
	WithdrawalStatusRejected           WithdrawalStatus = "REJECTED"
	WithdrawalStatusCancelled          WithdrawalStatus = "CANCELLED"
)

type WithdrawalRequestReadModel struct {
	WithdrawalID         string           `json:"withdrawal_id"`
	AuthorID             string           `json:"author_id"`
	GrossAmountMinor     int64            `json:"gross_amount_minor"`
	FeeMinor             int64            `json:"fee_minor"`
	NetAmountMinor       int64            `json:"net_amount_minor"`
	Currency             string           `json:"currency"`
	Status               WithdrawalStatus `json:"status"`
	BankAccountMasked    string           `json:"bank_account_masked"`
	RequestedAt          time.Time        `json:"requested_at"`
	ReviewedBy           string           `json:"reviewed_by"`
	ReviewedAt           *time.Time       `json:"reviewed_at"`
	ReviewReason         string           `json:"review_reason"`
	TransferReference    string           `json:"transfer_reference"`
	TransferRegisteredAt *time.Time       `json:"transfer_registered_at"`
	CompletedAt          *time.Time       `json:"completed_at"`
	ReceiptRef           string           `json:"receipt_ref"`
}

// PackUploadManifest 是作者上传包内 manifest 与上传表单的共享契约形状。
type PackUploadManifest struct {
	PackID            string      `json:"pack_id"`
	Name              string      `json:"name"`
	Version           string      `json:"version"`
	Kind              ProductKind `json:"kind"`
	MinTruzhenVersion string      `json:"min_truzhen_version"`
	Description       string      `json:"description"`
	RegionCode        string      `json:"region_code"`
	ArchTags          []string    `json:"arch_tags"`
	PriceWeeklyMinor  int64       `json:"price_weekly_minor"`
	PriceMonthlyMinor int64       `json:"price_monthly_minor"`
	PriceYearlyMinor  int64       `json:"price_yearly_minor"`
	IsFree            bool        `json:"is_free"`
}

type PackUploadReadModel struct {
	UploadID          string            `json:"upload_id"`
	ProductID         string            `json:"product_id"`
	Version           string            `json:"version"`
	FileName          string            `json:"file_name"`
	SizeBytes         int64             `json:"size_bytes"`
	SHA256            string            `json:"sha256"`
	TrustVerifyStatus TrustVerifyStatus `json:"trust_verify_status"`
	Status            ProductStatus     `json:"status"`
	CreatedAt         time.Time         `json:"created_at"`
}
