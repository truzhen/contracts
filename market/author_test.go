package market_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/truzhen/contracts/market"
)

func TestAuthorCertificationStatusGolden(t *testing.T) {
	cases := []struct {
		got  string
		want string
	}{
		{string(market.AuthorCertificationStatusNone), "none"},
		{string(market.AuthorCertificationStatusPendingPayment), "pending_payment"},
		{string(market.AuthorCertificationStatusPendingReview), "pending_review"},
		{string(market.AuthorCertificationStatusCertified), "certified"},
		{string(market.AuthorCertificationStatusRejected), "rejected"},
		{string(market.AuthorCertificationStatusRevoked), "revoked"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Fatalf("作者认证状态漂移：got %q want %q", tc.got, tc.want)
		}
	}
}

func TestWithdrawalStatusGolden(t *testing.T) {
	cases := []struct {
		got  string
		want string
	}{
		{string(market.WithdrawalStatusRequested), "REQUESTED"},
		{string(market.WithdrawalStatusUnderReview), "UNDER_REVIEW"},
		{string(market.WithdrawalStatusApproved), "APPROVED"},
		{string(market.WithdrawalStatusTransferRegistered), "TRANSFER_REGISTERED"},
		{string(market.WithdrawalStatusCompleted), "COMPLETED"},
		{string(market.WithdrawalStatusRejected), "REJECTED"},
		{string(market.WithdrawalStatusCancelled), "CANCELLED"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Fatalf("提现状态漂移：got %q want %q", tc.got, tc.want)
		}
	}
}

func TestProductAndTrustEnumsGolden(t *testing.T) {
	cases := []struct {
		got  string
		want string
	}{
		{string(market.ProductKindScenePack), "scene_pack"},
		{string(market.ProductKindRolePack), "role_pack"},
		{string(market.ProductKindCapabilityPack), "capability_pack"},
		{string(market.ProductKindSkillBundle), "skill_bundle"},
		{string(market.ProductStatusPendingReview), "pending_review"},
		{string(market.ProductStatusListed), "listed"},
		{string(market.ProductStatusDelisted), "delisted"},
		{string(market.TrustVerifyStatusVerified), "verified"},
		{string(market.TrustVerifyStatusUnverified), "unverified"},
		{string(market.TrustVerifyStatusFailed), "failed"},
		{string(market.TrustVerifyStatusProviderMissing), "provider_missing"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Fatalf("商品 / 信任枚举漂移：got %q want %q", tc.got, tc.want)
		}
	}
}

func TestAuthorRevenueJSONShape(t *testing.T) {
	rm := market.AuthorRevenueReadModel{
		AuthorID:               "author-1",
		Currency:               "CNY",
		GrossAuthorShareMinor:  10000,
		TaxRateBps:             300,
		TaxMinor:               300,
		NetAccruedMinor:        9700,
		WithdrawnMinor:         1000,
		PendingWithdrawalMinor: 2000,
		WithdrawableMinor:      6700,
		WithdrawFeeBps:         50,
		WithdrawMinAmountMinor: 1000,
		PerPack: []market.PackRevenueLineReadModel{{
			ProductID:        "pack-1",
			ProductName:      "墅学家装修 Pack",
			RegionCode:       "CN-SH",
			ArchTags:         []string{"renovation"},
			ShareRatioBps:    7000,
			PaidOrderCount:   2,
			GrossMinor:       12000,
			AuthorShareMinor: 8400,
		}},
		GeneratedAt: time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC),
	}
	raw, err := json.Marshal(rm)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"author_id", "currency", "gross_author_share_minor", "tax_rate_bps", "tax_minor", "net_accrued_minor", "withdrawn_minor", "pending_withdrawal_minor", "withdrawable_minor", "withdraw_fee_bps", "withdraw_min_amount_minor", "per_pack", "generated_at"} {
		if _, ok := m[key]; !ok {
			t.Fatalf("AuthorRevenueReadModel 缺 JSON 字段 %q：%s", key, raw)
		}
	}
	if lines, ok := m["per_pack"].([]any); !ok || len(lines) != 1 {
		t.Fatalf("per_pack 应包含逐包收益行：%s", raw)
	}
}

func TestAuthorCertificationJSONShape(t *testing.T) {
	rm := market.AuthorCertificationReadModel{
		Status:                 market.AuthorCertificationStatusPendingReview,
		FeeMinor:               0,
		Currency:               "CNY",
		OrderID:                "order-cert-1",
		RegionCode:             "CN-SH",
		RegionDisplay:          "上海",
		ApplicantName:          "李*",
		ContactMasked:          "138****0000",
		PayoutAccountMasked:    "****1234",
		IDCardMasked:           "3101**********1234",
		ResidenceAddressMasked: "上海市***",
		BankName:               "招商银行",
		ReviewedAt:             nil,
		RejectReason:           "",
	}
	raw, err := json.Marshal(rm)
	if err != nil {
		t.Fatalf("marshal readmodel: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal readmodel: %v", err)
	}
	for _, key := range []string{"status", "fee_minor", "currency", "order_id", "region_code", "region_display", "applicant_name", "contact_masked", "payout_account_masked", "id_card_masked", "residence_address_masked", "bank_name", "reviewed_at", "reject_reason"} {
		if _, ok := m[key]; !ok {
			t.Fatalf("AuthorCertificationReadModel 缺 JSON 字段 %q：%s", key, raw)
		}
	}

	registration := market.AuthorCertificationRegistration{
		ApplicantName:    "李雷",
		IDCard:           "310101199001011234",
		Contact:          "13800000000",
		PayoutAccount:    "6222000000001234",
		ResidenceAddress: "上海市浦东新区",
		BankName:         "招商银行",
		RegionCode:       "CN-SH",
		RegionDisplay:    "上海",
	}
	raw, err = json.Marshal(registration)
	if err != nil {
		t.Fatalf("marshal registration: %v", err)
	}
	m = map[string]any{}
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal registration: %v", err)
	}
	for _, key := range []string{"applicant_name", "id_card", "contact", "payout_account", "residence_address", "bank_name", "region_code", "region_display"} {
		if _, ok := m[key]; !ok {
			t.Fatalf("AuthorCertificationRegistration 缺 JSON 字段 %q：%s", key, raw)
		}
	}
}

func TestWithdrawalJSONShape(t *testing.T) {
	rm := market.WithdrawalRequestReadModel{
		WithdrawalID:         "wd-1",
		AuthorID:             "author-1",
		GrossAmountMinor:     10000,
		FeeMinor:             50,
		NetAmountMinor:       9950,
		Currency:             "CNY",
		Status:               market.WithdrawalStatusApproved,
		BankAccountMasked:    "****1234",
		RequestedAt:          time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC),
		ReviewedBy:           "ops-1",
		ReviewedAt:           nil,
		ReviewReason:         "ok",
		TransferReference:    "TRF-20260703-001",
		TransferRegisteredAt: nil,
		CompletedAt:          nil,
		ReceiptRef:           "receipt_candidate:withdrawal:wd-1",
	}
	raw, err := json.Marshal(rm)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"withdrawal_id", "author_id", "gross_amount_minor", "fee_minor", "net_amount_minor", "currency", "status", "bank_account_masked", "requested_at", "reviewed_by", "reviewed_at", "review_reason", "transfer_reference", "transfer_registered_at", "completed_at", "receipt_ref"} {
		if _, ok := m[key]; !ok {
			t.Fatalf("WithdrawalRequestReadModel 缺 JSON 字段 %q：%s", key, raw)
		}
	}
}

func TestPackUploadContractsJSONShape(t *testing.T) {
	manifest := market.PackUploadManifest{
		PackID:            "villa-renovation",
		Name:              "墅学家装修 Pack",
		Version:           "1.0.0",
		Kind:              market.ProductKindScenePack,
		MinTruzhenVersion: "3.0.0",
		Description:       "装修咨询场景包",
		RegionCode:        "CN-SH",
		ArchTags:          []string{"renovation", "villa"},
		PriceWeeklyMinor:  0,
		PriceMonthlyMinor: 0,
		PriceYearlyMinor:  0,
		IsFree:            true,
	}
	raw, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal manifest: %v", err)
	}
	for _, key := range []string{"pack_id", "name", "version", "kind", "min_truzhen_version", "description", "region_code", "arch_tags", "price_weekly_minor", "price_monthly_minor", "price_yearly_minor", "is_free"} {
		if _, ok := m[key]; !ok {
			t.Fatalf("PackUploadManifest 缺 JSON 字段 %q：%s", key, raw)
		}
	}

	rm := market.PackUploadReadModel{
		UploadID:          "upload-1",
		ProductID:         "villa-renovation",
		Version:           "1.0.0",
		FileName:          "villa-renovation.zip",
		SizeBytes:         1024,
		SHA256:            "sha256:abc",
		TrustVerifyStatus: market.TrustVerifyStatusProviderMissing,
		Status:            market.ProductStatusPendingReview,
		CreatedAt:         time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC),
	}
	raw, err = json.Marshal(rm)
	if err != nil {
		t.Fatalf("marshal readmodel: %v", err)
	}
	m = map[string]any{}
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal readmodel: %v", err)
	}
	for _, key := range []string{"upload_id", "product_id", "version", "file_name", "size_bytes", "sha256", "trust_verify_status", "status", "created_at"} {
		if _, ok := m[key]; !ok {
			t.Fatalf("PackUploadReadModel 缺 JSON 字段 %q：%s", key, raw)
		}
	}
}
