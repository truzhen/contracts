package market_test

import (
	"encoding/json"
	"testing"

	"github.com/truzhen/contracts/market"
)

func TestBuyerMarketDTOsMarshalStableFields(t *testing.T) {
	product := market.MarketCatalogProduct{
		ProductID:          "product_env",
		Name:               "环保执法包",
		ProductKind:        market.ProductKindScenePack,
		Version:            "1.0.0",
		SellerUserID:       "seller_1",
		AuthorID:           "author_1",
		Status:             market.ProductStatusListed,
		PriceWeeklyMinor:   1000,
		PriceMonthlyMinor:  3000,
		PriceYearlyMinor:   30000,
		PriceLifetimeMinor: 90000,
		Currency:           "CNY",
		IsFree:             false,
		RegionCode:         "CN",
		RegionDisplayName:  "中国",
		ArchTags:           []string{"transaction_spine"},
	}
	raw, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("marshal product: %v", err)
	}
	for _, want := range []string{
		`"product_id":"product_env"`,
		`"price_lifetime_minor":90000`,
		`"product_kind":"scene_pack"`,
	} {
		if !json.Valid(raw) || !containsJSON(raw, want) {
			t.Fatalf("product JSON missing %s: %s", want, raw)
		}
	}

	bundle := market.PackExportBundle{
		Status:       "ready",
		PackRef:      "pack://environmental-enforcement",
		Version:      "1.0.0",
		BundleID:     "bundle_env_1",
		BundleZipURL: "/v3/capability/lifecycle/export/bundle.zip",
		SHA256:       "abc123",
		SizeBytes:    12,
		UploadPrefill: market.PackExportUploadPrefill{
			PackID: "environmental-enforcement",
			Kind:   string(market.ProductKindScenePack),
		},
	}
	raw, err = json.Marshal(bundle)
	if err != nil {
		t.Fatalf("marshal bundle: %v", err)
	}
	for _, want := range []string{
		`"bundle_zip_url":"/v3/capability/lifecycle/export/bundle.zip"`,
		`"sha256":"abc123"`,
		`"upload_prefill"`,
	} {
		if !json.Valid(raw) || !containsJSON(raw, want) {
			t.Fatalf("bundle JSON missing %s: %s", want, raw)
		}
	}
}

func containsJSON(raw []byte, want string) bool {
	s := string(raw)
	for i := 0; i+len(want) <= len(s); i++ {
		if s[i:i+len(want)] == want {
			return true
		}
	}
	return false
}
