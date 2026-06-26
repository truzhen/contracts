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
