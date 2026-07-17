package contracts_test

import (
	"encoding/json"
	"testing"

	contracts "github.com/truzhen/contracts"
)

func TestV071SchemasEmbeddedAndParse(t *testing.T) {
	schemas := map[string][]byte{
		"pack-usage-contribution-candidate":  contracts.PackUsageContributionCandidateSchemaJSON,
		"pack-version-migration-candidate":   contracts.PackVersionMigrationCandidateSchemaJSON,
		"contribution-receipt":               contracts.ContributionReceiptSchemaJSON,
		"market-catalog-product":             contracts.MarketCatalogProductSchemaJSON,
		"market-entitlement":                 contracts.MarketEntitlementSchemaJSON,
		"market-checkout-result":             contracts.MarketCheckoutResultSchemaJSON,
		"market-order-status":                contracts.MarketOrderStatusSchemaJSON,
		"market-local-gate-check-result":     contracts.MarketLocalGateCheckResultSchemaJSON,
		"pack-install-result":                contracts.PackInstallResultSchemaJSON,
		"pack-export-bundle":                 contracts.PackExportBundleSchemaJSON,
		"mobile-pairing-bootstrap-request":   contracts.MobilePairingBootstrapRequestSchemaJSON,
		"mobile-pairing-bootstrap-candidate": contracts.MobilePairingBootstrapCandidateSchemaJSON,
		"mobile-session-issue-intent":        contracts.MobileSessionIssueIntentSchemaJSON,
	}
	for name, raw := range schemas {
		if len(raw) == 0 {
			t.Fatalf("%s schema embed is empty", name)
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatalf("%s schema must be valid JSON: %v", name, err)
		}
		if doc["additionalProperties"] != false {
			t.Fatalf("%s schema must be closed with additionalProperties=false", name)
		}
		if _, ok := doc["required"].([]interface{}); !ok {
			t.Fatalf("%s schema must declare required fields", name)
		}
	}
}

func TestPackInstallResultCarriesImmutableAuthorizationAndArtifactProof(t *testing.T) {
	var doc struct {
		Required   []string                          `json:"required"`
		Properties map[string]map[string]interface{} `json:"properties"`
	}
	if err := json.Unmarshal(contracts.PackInstallResultSchemaJSON, &doc); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"decision_ref", "run_id", "nonce", "artifact_sha256", "receipt_ref", "evidence_refs"} {
		if doc.Properties[field] == nil {
			t.Fatalf("pack install proof schema missing %s", field)
		}
	}
	// Additive optional fields preserve existing v0.12 consumers. The OS HTTP
	// success operation imposes the stronger required set for new installs.
	for _, field := range []string{"decision_ref", "run_id", "nonce", "artifact_sha256"} {
		for _, required := range doc.Required {
			if required == field {
				t.Fatalf("additive proof field %s must remain optional in shared DTO", field)
			}
		}
	}
	if got := doc.Properties["artifact_sha256"]["pattern"]; got != "^[a-f0-9]{64}$" {
		t.Fatalf("artifact_sha256 must be canonical lowercase SHA-256, got %v", got)
	}
}
