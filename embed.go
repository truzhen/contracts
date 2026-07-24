// Package contracts embeds the machine-checkable JSON Schemas so Go services
// and CI tests can validate produced artifacts against the canonical contract
// files instead of re-declaring the shapes in code (T08 阶段2：契约收敛，
// contracts 校验进 CI).
package contracts

import _ "embed"

// SceneFlowSpecSchemaJSON is contracts/scene-flow-spec.schema.json.
//
//go:embed scene-flow-spec.schema.json
var SceneFlowSpecSchemaJSON []byte

// ScenePackSpecSchemaJSON is contracts/scene-pack-spec.schema.json.
//
//go:embed scene-pack-spec.schema.json
var ScenePackSpecSchemaJSON []byte

// FlowViewSpecSchemaJSON is contracts/flow-view-spec.schema.json.
//
//go:embed flow-view-spec.schema.json
var FlowViewSpecSchemaJSON []byte

// VisualUnitSpecSchemaJSON is contracts/visual-unit-spec.schema.json
// （前端 7 类主权视觉单元封顶规格契约，client layer 多端单一来源；
// T12 阶段A：client layer 契约收敛起步）。
//
//go:embed visual-unit-spec.schema.json
var VisualUnitSpecSchemaJSON []byte

// TransactionObjectProjectionSchemaJSON is contracts/transaction-object-projection.schema.json
// （事务对象 05 BusinessObject 的前端只读投影 DTO 契约，client layer 渲染事务对象卡的单一来源；
// T12 阶段A：client layer 契约面补全）。
//
//go:embed transaction-object-projection.schema.json
var TransactionObjectProjectionSchemaJSON []byte

// MobilePairingBootstrapRequestSchemaJSON is the authority-free first-pairing
// request shape. It intentionally excludes every authority and credential.
//
//go:embed mobile-pairing-bootstrap-request.schema.json
var MobilePairingBootstrapRequestSchemaJSON []byte

// MobilePairingBootstrapCandidateSchemaJSON is the Host-owned candidate
// projection shown to phone and PC clients before a mobile session exists.
//
//go:embed mobile-pairing-bootstrap-candidate.schema.json
var MobilePairingBootstrapCandidateSchemaJSON []byte

// MobileSessionIssueIntentSchemaJSON is the body shape for post-approval
// session issuance. Bootstrap proof stays header-only and is not embedded.
//
//go:embed mobile-session-issue-intent.schema.json
var MobileSessionIssueIntentSchemaJSON []byte

// CandidateEnvelopeSchemaJSON is contracts/candidate-envelope.schema.json
// （candidates.CandidateEnvelope 的 JSON 表达，client layer 候选卡面向、CI 校验；
// T12 阶段A：client layer 契约面补全）。
//
//go:embed candidate-envelope.schema.json
var CandidateEnvelopeSchemaJSON []byte

// ReceiptEnvelopeSchemaJSON is contracts/receipt-envelope.schema.json
// （receipts.ReceiptEnvelope 的 JSON 表达，client layer 回执卡面向、CI 校验；
// T12 阶段A：client layer 契约面补全）。
//
//go:embed receipt-envelope.schema.json
var ReceiptEnvelopeSchemaJSON []byte

// PackManifestSchemaJSON is contracts/pack-manifest.schema.json.
//
//go:embed pack-manifest.schema.json
var PackManifestSchemaJSON []byte

// ProviderRequirementSchemaJSON is contracts/provider-requirement.schema.json.
//
//go:embed provider-requirement.schema.json
var ProviderRequirementSchemaJSON []byte

// SoftwareResolutionLockSchemaJSON is contracts/software-resolution-lock.schema.json.
//
//go:embed software-resolution-lock.schema.json
var SoftwareResolutionLockSchemaJSON []byte

// PackUsageContributionCandidateSchemaJSON is contracts/pack-usage-contribution-candidate.schema.json.
//
//go:embed pack-usage-contribution-candidate.schema.json
var PackUsageContributionCandidateSchemaJSON []byte

// PackVersionMigrationCandidateSchemaJSON is contracts/pack-version-migration-candidate.schema.json.
//
//go:embed pack-version-migration-candidate.schema.json
var PackVersionMigrationCandidateSchemaJSON []byte

// ContributionReceiptSchemaJSON is contracts/contribution-receipt.schema.json.
//
//go:embed contribution-receipt.schema.json
var ContributionReceiptSchemaJSON []byte

// MarketCatalogProductSchemaJSON is contracts/market-catalog-product.schema.json.
//
//go:embed market-catalog-product.schema.json
var MarketCatalogProductSchemaJSON []byte

// MarketEntitlementSchemaJSON is contracts/market-entitlement.schema.json.
//
//go:embed market-entitlement.schema.json
var MarketEntitlementSchemaJSON []byte

// MarketCheckoutResultSchemaJSON is contracts/market-checkout-result.schema.json.
//
//go:embed market-checkout-result.schema.json
var MarketCheckoutResultSchemaJSON []byte

// MarketOrderStatusSchemaJSON is contracts/market-order-status.schema.json.
//
//go:embed market-order-status.schema.json
var MarketOrderStatusSchemaJSON []byte

// MarketLocalGateCheckResultSchemaJSON is contracts/market-local-gate-check-result.schema.json.
//
//go:embed market-local-gate-check-result.schema.json
var MarketLocalGateCheckResultSchemaJSON []byte

// PackInstallResultSchemaJSON is contracts/pack-install-result.schema.json.
//
//go:embed pack-install-result.schema.json
var PackInstallResultSchemaJSON []byte

// PackExportBundleSchemaJSON is contracts/pack-export-bundle.schema.json.
//
//go:embed pack-export-bundle.schema.json
var PackExportBundleSchemaJSON []byte

// Intent Spine 五件套（T12 阶段3：契约收敛进 contracts/spines/）。

// IntentEventSchemaJSON is contracts/spines/intent-event.schema.json.
//
//go:embed spines/intent-event.schema.json
var IntentEventSchemaJSON []byte

// IntentInboxItemSchemaJSON is contracts/spines/intent-inbox-item.schema.json.
//
//go:embed spines/intent-inbox-item.schema.json
var IntentInboxItemSchemaJSON []byte

// IntentClassificationSchemaJSON is contracts/spines/intent-classification.schema.json.
//
//go:embed spines/intent-classification.schema.json
var IntentClassificationSchemaJSON []byte

// IntentToCandidateResultSchemaJSON is contracts/spines/intent-to-candidate-result.schema.json.
//
//go:embed spines/intent-to-candidate-result.schema.json
var IntentToCandidateResultSchemaJSON []byte

// IntentReceiptSchemaJSON is contracts/spines/intent-receipt.schema.json.
//
//go:embed spines/intent-receipt.schema.json
var IntentReceiptSchemaJSON []byte

// MonitoringEventSchemaJSON is contracts/monitoring/monitoring-event.schema.json.
//
//go:embed monitoring/monitoring-event.schema.json
var MonitoringEventSchemaJSON []byte

// FaultIncidentSchemaJSON is contracts/monitoring/fault-incident.schema.json.
//
//go:embed monitoring/fault-incident.schema.json
var FaultIncidentSchemaJSON []byte

// DeliberationSessionReadModelSchemaJSON is the client-safe session projection
// for the governed multi-provider deliberation flow.
//
//go:embed deliberation-session-readmodel.schema.json
var DeliberationSessionReadModelSchemaJSON []byte

// DeliberationTurnReadModelSchemaJSON is the client-safe turn projection. It
// carries a question artifact reference and SHA-256, never question content.
//
//go:embed deliberation-turn-readmodel.schema.json
var DeliberationTurnReadModelSchemaJSON []byte

// DeliberationProviderLaneReadModelSchemaJSON is the adapter-lane projection
// with separately declared release eligibility and runtime readiness.
//
//go:embed deliberation-provider-lane-readmodel.schema.json
var DeliberationProviderLaneReadModelSchemaJSON []byte

// DeliberationAutomationGrantReadModelSchemaJSON is the projection of a
// Base-issued bounded automation grant; clients cannot self-mint its refs.
//
//go:embed deliberation-automation-grant-readmodel.schema.json
var DeliberationAutomationGrantReadModelSchemaJSON []byte

// DeliberationSynthesisCandidateSchemaJSON describes an AI synthesis output
// that is permanently candidate_only and non_formal.
//
//go:embed deliberation-synthesis-candidate.schema.json
var DeliberationSynthesisCandidateSchemaJSON []byte
