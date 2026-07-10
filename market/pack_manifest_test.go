package market_test

import (
	"encoding/json"
	"strings"
	"testing"

	contracts "github.com/truzhen/contracts"
	"github.com/truzhen/contracts/market"
)

func TestPackSoftwareRequirementGoldenValues(t *testing.T) {
	req := market.PackSoftwareRequirement{
		RequirementID:        "frappe-runtime",
		SoftwareFamily:       "frappe-erpnext-stack",
		ProviderFamily:       "frappe",
		VersionRange:         ">=14.0.0,<16.0.0",
		RequiredCapabilities: []string{"project.customer", "project.milestone"},
		LicensePolicy:        market.LicensePolicyReviewRequired,
		IsolationPolicy:      market.SoftwareIsolationReusePreferred,
		FallbackPolicy:       market.SoftwareFallbackProviderMissing,
		GatewayClass:         market.GatewayClassExecution,
		RiskClass:            market.RiskClassHigh,
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal requirement: %v", err)
	}
	var roundTrip market.PackSoftwareRequirement
	if err := json.Unmarshal(b, &roundTrip); err != nil {
		t.Fatalf("unmarshal requirement: %v", err)
	}
	if roundTrip.VersionRange != ">=14.0.0,<16.0.0" {
		t.Fatalf("version range drift: got %q", roundTrip.VersionRange)
	}
	s := string(b)
	for _, want := range []string{
		`"software_family":"frappe-erpnext-stack"`,
		`"isolation_policy":"reuse_preferred"`,
		`"fallback_policy":"provider_missing"`,
		`"gateway_class":"execution"`,
		`"risk_class":"high"`,
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %s in %s", want, s)
		}
	}
}

func TestSoftwareResolutionLockGoldenValues(t *testing.T) {
	lock := market.SoftwareResolutionLock{
		LockID:              "software_lock://frappe-runtime",
		PackRef:             "scene_pack://smart-home-owner@1.0.0",
		RequirementID:       "frappe-runtime",
		SoftwareFamily:      "frappe-erpnext-stack",
		ResolvedSoftwareRef: "software://frappe-erpnext-v16",
		ResolvedVersion:     "15.0.0",
		ProviderResourceRef: "provider_resource://frappe-erpnext-v16",
		Resolution:          market.SoftwareResolutionReused,
		ReceiptRef:          "receipt://software-lock-001",
		ResolvedAt:          "2026-07-07T00:00:00Z",
	}
	b, err := json.Marshal(lock)
	if err != nil {
		t.Fatalf("marshal lock: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		`"resolution":"reused"`,
		`"receipt_ref":"receipt://software-lock-001"`,
		`"provider_resource_ref":"provider_resource://frappe-erpnext-v16"`,
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %s in %s", want, s)
		}
	}
}

func TestPackManifestLifecycleStatusGoldenValues(t *testing.T) {
	m := market.PackManifest{
		PackID:            "smart-home-owner",
		Name:              "智能家居业主服务包",
		Version:           "1.0.0",
		Kind:              market.ProductKindScenePack,
		MinTruzhenVersion: "3.0.0",
		LifecycleStatus:   market.PackLifecycleAccepted,
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	if !strings.Contains(string(b), `"lifecycle_status":"已验收"`) {
		t.Fatalf("missing lifecycle_status in %s", string(b))
	}

	// omitempty：未声明生命周期的历史 manifest 序列化后不得出现该键。
	legacy := market.PackManifest{
		PackID:            "legacy-pack",
		Name:              "legacy",
		Version:           "0.1.0",
		Kind:              market.ProductKindCapabilityPack,
		MinTruzhenVersion: "3.0.0",
	}
	lb, err := json.Marshal(legacy)
	if err != nil {
		t.Fatalf("marshal legacy manifest: %v", err)
	}
	if strings.Contains(string(lb), "lifecycle_status") {
		t.Fatalf("legacy manifest must omit lifecycle_status: %s", string(lb))
	}
}

func TestPackLifecycleStatusCoversEightTiers(t *testing.T) {
	// 八档统一生命周期（治理字典中文值为准，见 candidate-set.json 既有约定）。
	tiers := []market.PackLifecycleStatus{
		market.PackLifecycleIdea,
		market.PackLifecycleDesigning,
		market.PackLifecycleContractFixed,
		market.PackLifecycleImplemented,
		market.PackLifecycleWired,
		market.PackLifecycleAccepted,
		market.PackLifecycleReleased,
		market.PackLifecycleDeprecated,
	}
	want := []string{"想法", "设计中", "契约已定", "已实现", "已接线", "已验收", "已发布", "已弃用"}
	if len(tiers) != len(want) {
		t.Fatalf("tier count %d != %d", len(tiers), len(want))
	}
	for i, tier := range tiers {
		if string(tier) != want[i] {
			t.Fatalf("tier %d: got %q want %q", i, tier, want[i])
		}
	}
}

// Go 常量与 embed schema 的 enum 必须逐值一致（防同义词漂移）。
// 背景：go-schema-map 第 5 对因 checker v1 不支持 $ref 属性暂不可登记
// （backlog：嵌套 $ref 展开），本测试承担该字段的同步守卫。
func TestPackManifestLifecycleEnumMatchesSchema(t *testing.T) {
	var schema struct {
		Properties struct {
			LifecycleStatus struct {
				Enum []string `json:"enum"`
			} `json:"lifecycle_status"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(contracts.PackManifestSchemaJSON, &schema); err != nil {
		t.Fatalf("parse embedded pack-manifest schema: %v", err)
	}
	want := []string{
		string(market.PackLifecycleIdea),
		string(market.PackLifecycleDesigning),
		string(market.PackLifecycleContractFixed),
		string(market.PackLifecycleImplemented),
		string(market.PackLifecycleWired),
		string(market.PackLifecycleAccepted),
		string(market.PackLifecycleReleased),
		string(market.PackLifecycleDeprecated),
	}
	got := schema.Properties.LifecycleStatus.Enum
	if len(got) != len(want) {
		t.Fatalf("schema enum count %d != Go constants %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("enum[%d]: schema %q != Go %q", i, got[i], want[i])
		}
	}
}

func TestSoftwareResolutionLockSupportsResolverMVPOutcomes(t *testing.T) {
	outcomes := []market.SoftwareResolution{
		market.SoftwareResolutionReused,
		market.SoftwareResolutionInstallRequired,
		market.SoftwareResolutionVersionConflict,
		market.SoftwareResolutionIsolationRequired,
		market.SoftwareResolutionBlocked,
		market.SoftwareResolutionNotReady,
		market.SoftwareResolutionProviderMissing,
	}
	for _, outcome := range outcomes {
		lock := market.SoftwareResolutionLock{
			LockID:         "software_lock://baserow-runtime",
			PackRef:        "scene_pack://baserow-pack-b@1.0.0",
			RequirementID:  "baserow-runtime",
			SoftwareFamily: "baserow-family",
			Resolution:     outcome,
			ResolvedAt:     "2026-07-08T00:00:00Z",
		}
		b, err := json.Marshal(lock)
		if err != nil {
			t.Fatalf("marshal %s: %v", outcome, err)
		}
		if !strings.Contains(string(b), `"resolution":"`+string(outcome)+`"`) {
			t.Fatalf("missing resolution %s in %s", outcome, string(b))
		}
	}
}
