package contracts_test

import (
	"encoding/json"
	"testing"

	contracts "github.com/truzhen/contracts"
)

func TestV071SchemasEmbeddedAndParse(t *testing.T) {
	schemas := map[string][]byte{
		"pack-usage-contribution-candidate":   contracts.PackUsageContributionCandidateSchemaJSON,
		"pack-version-migration-candidate":   contracts.PackVersionMigrationCandidateSchemaJSON,
		"contribution-receipt":                contracts.ContributionReceiptSchemaJSON,
		"market-catalog-product":              contracts.MarketCatalogProductSchemaJSON,
		"market-entitlement":                  contracts.MarketEntitlementSchemaJSON,
		"market-checkout-result":              contracts.MarketCheckoutResultSchemaJSON,
		"market-order-status":                 contracts.MarketOrderStatusSchemaJSON,
		"market-local-gate-check-result":      contracts.MarketLocalGateCheckResultSchemaJSON,
		"pack-install-result":                 contracts.PackInstallResultSchemaJSON,
		"pack-export-bundle":                  contracts.PackExportBundleSchemaJSON,
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
