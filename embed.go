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

// CloudEntitlementSchemaJSON is contracts/cloud/entitlement.schema.json.
//
//go:embed cloud/entitlement.schema.json
var CloudEntitlementSchemaJSON []byte

// CloudLicenseSchemaJSON is contracts/cloud/license.schema.json.
//
//go:embed cloud/license.schema.json
var CloudLicenseSchemaJSON []byte

// CloudPaymentSchemaJSON is contracts/cloud/payment.schema.json.
//
//go:embed cloud/payment.schema.json
var CloudPaymentSchemaJSON []byte

// CloudPackListingSchemaJSON is contracts/cloud/pack_listing.schema.json.
//
//go:embed cloud/pack_listing.schema.json
var CloudPackListingSchemaJSON []byte

// CloudSessionSchemaJSON is contracts/cloud/session.schema.json.
//
//go:embed cloud/session.schema.json
var CloudSessionSchemaJSON []byte

// CloudReleaseSchemaJSON is contracts/cloud/release.schema.json.
//
//go:embed cloud/release.schema.json
var CloudReleaseSchemaJSON []byte

// CloudWebSurfaceSchemaJSON is contracts/cloud/web_surface.schema.json.
//
//go:embed cloud/web_surface.schema.json
var CloudWebSurfaceSchemaJSON []byte
