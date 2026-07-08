package market_test

import (
	"encoding/json"
	"strings"
	"testing"

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
