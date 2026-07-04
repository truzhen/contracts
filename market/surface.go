// Package market 声明本地网关与云端市场之间的「市场表面契约」（§18-A 契约化）。
//
// 背景：六仓拆分后，市场/许可的服务端真相唯一归属 truzhen-cloud（03-market-license
// + 04-payment-settlement）；truzhenos 17-market 只保留薄代理（受控 ReadModel 代理，
// 订单/价格/权益全部由云端签发，本地绝不自铸）。本包是两侧共享形状的唯一声明处：
//   - 会话头名称（本地代理转发云端请求时携带的登录态头）；
//   - 市场表面端点路径（本地代理允许触达的云端表面）;
//   - 管理面转发硬 allowlist（主权边界：哪些 admin 路径允许离开本机）。
//
// 改动纪律：本包任何常量/名单变化 = 跨仓契约变化，必须 bump VERSION 并三仓对齐。
package market

import (
	"net/url"
	"strings"
)

// SessionHeader 是本地代理向云端市场转发请求时携带云端会话 ID 的 HTTP 头。
const SessionHeader = "X-Truzhen-Session-Id"

// LoginRequest 是云端市场登录请求体（本地代理透传，不落盘口令）。
type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// LoginResponse 是云端市场登录成功后的会话形状。
type LoginResponse struct {
	SessionID   string `json:"session_id"`
	Role        string `json:"role"`
	DisplayName string `json:"display_name"`
	PhoneMasked string `json:"phone_masked"`
}

// SessionStatus 是云市场会话投射状态的封闭字符串集合。
type SessionStatus string

const (
	SessionStatusLoggedIn              SessionStatus = "logged_in"
	SessionStatusLoggedOut             SessionStatus = "logged_out"
	SessionStatusAuthSessionRequired   SessionStatus = "auth_session_required"
	SessionStatusRequiresMarketRelogin SessionStatus = "requires_market_relogin"
	SessionStatusCloudSessionInvalid   SessionStatus = "cloud_session_invalid"
	SessionStatusNotReady              SessionStatus = "not_ready"
	SessionStatusBlocked               SessionStatus = "blocked"
)

// PayableStatus 表达当前会话对云市场下单 / 授权的可支付性。
type PayableStatus string

const (
	PayableStatusUnknown                PayableStatus = "unknown"
	PayableStatusPayable                PayableStatus = "payable"
	PayableStatusNotPayable             PayableStatus = "not_payable"
	PayableStatusRequiresMarketRelogin  PayableStatus = "requires_market_relogin"
	PayableStatusPaymentProviderMissing PayableStatus = "payment_provider_missing"
	PayableStatusPaymentBlocked         PayableStatus = "payment_blocked"
)

// Role 是市场表面可投射的账户角色枚举；服务端仍是真相源。
type Role string

const (
	RoleBuyer    Role = "buyer"
	RoleAuthor   Role = "author"
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
)

// SessionProjection 是本地网关 / client 消费的云市场会话投射形状。
// requires_market_relogin 是商业动作硬信号；旧 requires_relogin 字段可由消费方
// 兼容保留，但新契约统一读 requires_market_relogin。
type SessionProjection struct {
	LoggedIn               bool          `json:"logged_in"`
	SessionID              string        `json:"session_id,omitempty"`
	SessionStatus          SessionStatus `json:"session_status,omitempty"`
	RefreshStatus          SessionStatus `json:"refresh_status,omitempty"`
	SessionRefStatus       SessionStatus `json:"session_ref_status,omitempty"`
	RequiresMarketRelogin  bool          `json:"requires_market_relogin,omitempty"`
	Payable                PayableStatus `json:"payable,omitempty"`
	Role                   Role          `json:"role,omitempty"`
	DisplayName            string        `json:"display_name,omitempty"`
	PhoneMasked            string        `json:"phone_masked,omitempty"`
	CloudMarketAuthStatus  SessionStatus `json:"cloud_market_authorization_status,omitempty"`
	CloudMarketAuthMessage string        `json:"cloud_market_authorization_message,omitempty"`
	GatewayProxyRequired   bool          `json:"gateway_proxy_required,omitempty"`
	NoRawPasswordSaved     bool          `json:"no_raw_password_saved,omitempty"`
	NoRawTokenReturned     bool          `json:"no_raw_token_returned,omitempty"`
	NoRawCookieReturned    bool          `json:"no_raw_cookie_returned,omitempty"`
	AuthIntentEvidenceRef  string        `json:"auth_intent_evidence_ref,omitempty"`
	ReceiptRef             string        `json:"receipt_ref,omitempty"`
	OwnerActionRef         string        `json:"owner_action_ref,omitempty"`
}

// 会话建立与目录同步表面。
const (
	PathAuthLogin          = "/auth/login"
	PathPackProxyPackages  = "/pack-proxy/packages"
	PathPackProxyPackTypes = "/pack-proxy/pack-types"
)

// 市场业务表面（作者侧 + 买家侧）。
const (
	PathAuthorRevenue               = "/v3/market/author/revenue"
	PathAuthorCertification         = "/v3/market/author/certification"
	PathAuthorCertificationRegister = "/v3/market/author/certification/register"
	PathAuthorWithdrawals           = "/v3/market/author/withdrawals"
	PathAuthorUploads               = "/v3/market/author/uploads"
	PathAuthorProducts              = "/v3/market/author/products"
	PathPackUpload                  = "/v3/market/packs/upload"
	PathLicenseProducts             = "/v3/market/license/products"
	PathLicenseCheckout             = "/v3/market/license/checkout"
	PathLicenseEntitlements         = "/v3/market/license/entitlements"
	PathLicenseLocalGateCheck       = "/v3/market/license/local-gate/check"
	PathDemandPool                  = "/v3/market/demand-pool"
	PathSuggestions                 = "/v3/market/suggestions"
)

// LicenseOrderPath 返回订单状态查询路径；订单 ID 作路径段转义，防路径穿越。
func LicenseOrderPath(orderID string) string {
	return "/v3/market/license/orders/" + url.PathEscape(orderID)
}

// WithdrawalCancelPath 返回提现撤销路径；提现 ID 作路径段转义。
func WithdrawalCancelPath(withdrawalID string) string {
	return "/v3/market/author/withdrawals/" + url.PathEscape(withdrawalID) + "/cancel"
}

// AuthorProductPricingPath 返回作者商品改价路径；商品 ID 作路径段转义。
func AuthorProductPricingPath(productID string) string {
	return PathAuthorProducts + "/" + url.PathEscape(productID) + "/pricing"
}

// AuthorProductDelistPath 返回作者商品下架路径；商品 ID 作路径段转义。
func AuthorProductDelistPath(productID string) string {
	return PathAuthorProducts + "/" + url.PathEscape(productID) + "/delist"
}

// AuthorProductRelistPath 返回作者商品重新提交上架路径；商品 ID 作路径段转义。
func AuthorProductRelistPath(productID string) string {
	return PathAuthorProducts + "/" + url.PathEscape(productID) + "/relist"
}

// PackDownloadPath 返回包下载路径；商品 ID 作路径段转义。
func PackDownloadPath(productID string) string {
	return "/v3/market/packs/" + url.PathEscape(productID) + "/download"
}

// adminForwardAllowlist 是本地网关允许转发的云端 admin 路径硬 allowlist。
// 主权边界：不在名单内的 admin 请求在离开本机前即被拒绝。
// 尾斜杠条目按前缀匹配（含去尾斜杠的精确命中）；无尾斜杠条目按精确或前缀匹配。
var adminForwardAllowlist = []string{
	"/v3/admin/stats/finance",
	"/v3/admin/withdrawals",
	"/v3/admin/finance/settings",
	"/v3/admin/packs/",
	"/v3/admin/ops/users",
	"/v3/admin/suggestions",
	"/v3/admin/certifications/",
	"/v3/market/license/policy/registration",
}

// AdminForwardAllowlist 返回 allowlist 副本；调用方改写不会污染契约本体。
func AdminForwardAllowlist() []string {
	return append([]string(nil), adminForwardAllowlist...)
}

// AdminPathAllowed 判定 admin 转发路径是否在硬 allowlist 内（query 部分不参与匹配）。
func AdminPathAllowed(path string) bool {
	if idx := strings.Index(path, "?"); idx >= 0 {
		path = path[:idx]
	}
	if path == "" {
		return false
	}
	for _, prefix := range adminForwardAllowlist {
		if path == strings.TrimSuffix(prefix, "/") || strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
