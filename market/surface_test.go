package market_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/truzhen/contracts/market"
)

// 市场表面契约（§18-A 契约化）：本地网关（truzhenos 17 薄代理）与云端市场
// （truzhen-cloud 03/market-proxy-server）之间的表面形状唯一声明处。
// 这些测试是黄金断言：改任何一项都等于改跨仓契约，必须 bump VERSION。

func TestSessionHeaderName(t *testing.T) {
	if market.SessionHeader != "X-Truzhen-Session-Id" {
		t.Fatalf("会话头名称漂移：%q", market.SessionHeader)
	}
}

func TestLoginRequestJSONShape(t *testing.T) {
	raw, err := json.Marshal(market.LoginRequest{Phone: "p", Password: "test-not-a-real-secret"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, k := range []string{"phone", "password"} {
		if _, ok := m[k]; !ok {
			t.Fatalf("LoginRequest 缺 JSON 字段 %q：%s", k, raw)
		}
	}
}

func TestLoginResponseJSONShape(t *testing.T) {
	golden := `{"session_id":"s1","role":"buyer","display_name":"张","phone_masked":"138****0000"}`
	var resp market.LoginResponse
	if err := json.Unmarshal([]byte(golden), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.SessionID != "s1" || resp.Role != "buyer" || resp.DisplayName != "张" || resp.PhoneMasked != "138****0000" {
		t.Fatalf("LoginResponse 字段映射漂移：%+v", resp)
	}
}

func TestSurfacePathsGolden(t *testing.T) {
	cases := map[string]string{
		market.PathAuthLogin:                   "/auth/login",
		market.PathPackProxyPackages:           "/pack-proxy/packages",
		market.PathPackProxyPackTypes:          "/pack-proxy/pack-types",
		market.PathAuthorRevenue:               "/v3/market/author/revenue",
		market.PathAuthorCertification:         "/v3/market/author/certification",
		market.PathAuthorCertificationRegister: "/v3/market/author/certification/register",
		market.PathAuthorWithdrawals:           "/v3/market/author/withdrawals",
		market.PathAuthorUploads:               "/v3/market/author/uploads",
		market.PathPackUpload:                  "/v3/market/packs/upload",
		market.PathLicenseProducts:             "/v3/market/license/products",
		market.PathLicenseCheckout:             "/v3/market/license/checkout",
		market.PathLicenseEntitlements:         "/v3/market/license/entitlements",
		market.PathLicenseLocalGateCheck:       "/v3/market/license/local-gate/check",
		market.PathDemandPool:                  "/v3/market/demand-pool",
		market.PathSuggestions:                 "/v3/market/suggestions",
	}
	for got, want := range cases {
		if got != want {
			t.Fatalf("表面路径漂移：got %q want %q", got, want)
		}
	}
}

func TestPathBuilersEscapeSegments(t *testing.T) {
	if got := market.LicenseOrderPath("ord-1"); got != "/v3/market/license/orders/ord-1" {
		t.Fatalf("LicenseOrderPath: %q", got)
	}
	if got := market.WithdrawalCancelPath("wd-1"); got != "/v3/market/author/withdrawals/wd-1/cancel" {
		t.Fatalf("WithdrawalCancelPath: %q", got)
	}
	if got := market.PackDownloadPath("pack a/b"); !strings.Contains(got, "pack%20a%2Fb") {
		t.Fatalf("PackDownloadPath 必须转义路径段：%q", got)
	}
	if got := market.LicenseOrderPath("../admin"); strings.Contains(got, "../") {
		t.Fatalf("LicenseOrderPath 必须转义防路径穿越：%q", got)
	}
}

func TestAdminForwardAllowlistGolden(t *testing.T) {
	want := []string{
		"/v3/admin/stats/finance",
		"/v3/admin/withdrawals",
		"/v3/admin/finance/settings",
		"/v3/admin/packs/",
		"/v3/admin/ops/users",
		"/v3/admin/suggestions",
		"/v3/admin/certifications/",
		"/v3/market/license/policy/registration",
	}
	got := market.AdminForwardAllowlist()
	if len(got) != len(want) {
		t.Fatalf("allowlist 条目数漂移：got %d want %d（%v）", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("allowlist[%d] 漂移：got %q want %q", i, got[i], want[i])
		}
	}
	// 返回副本：调用方改写不得污染契约本体。
	got[0] = "mutated"
	if again := market.AdminForwardAllowlist(); again[0] != want[0] {
		t.Fatalf("AdminForwardAllowlist 必须返回副本，本体被污染为 %q", again[0])
	}
}

func TestAdminPathAllowed(t *testing.T) {
	allowed := []string{
		"/v3/admin/stats/finance",
		"/v3/admin/withdrawals",
		"/v3/admin/withdrawals/wd-1/approve",
		"/v3/admin/packs",       // 尾斜杠条目的无斜杠精确命中
		"/v3/admin/packs/p1/审核", // 前缀命中
		"/v3/admin/ops/users?page=2",
		"/v3/market/license/policy/registration",
	}
	for _, p := range allowed {
		if !market.AdminPathAllowed(p) {
			t.Fatalf("应放行：%q", p)
		}
	}
	denied := []string{
		"/v3/admin",
		"/v3/admin/users",
		"/v3/admin/stats",
		"/v3/admin/packsX", // 无斜杠精确命中不得放大为任意前缀
		"/v3/market/license/checkout",
		"",
	}
	for _, p := range denied {
		if market.AdminPathAllowed(p) {
			t.Fatalf("必须拒绝：%q", p)
		}
	}
}
