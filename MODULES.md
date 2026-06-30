# truzhen-contracts 模块与 schema 清单

本仓是 Truzhen 五落点中的 **Pack / Candidate / Receipt / Surface / ReadModel schema** 权威源。每个子包只声明跨边界的数据形状，或提供无外部副作用的契约校验 / ref 派生 helper。

## Go 子包

| 子包 | 职责 | 不负责 |
| --- | --- | --- |
| `base/` | Base 主权核心契约：ActorContext / OwnerIdentityContext、PolicySet / PolicyKernel / PolicySnapshot、GateCandidateEnvelope、GateRequest、GateDecision、GateDecisionTrace、GateReceiptCandidate、OwnerDecision、BaseFormalizationGrant、授权模式、OwnerDelegationGrant、AgentDecision、ArtifactBinding / ArtifactUseIntent、Gateway adapter request。 | 不实现 Base orchestrator、policy engine、Receipt Ledger、provider、真实执行或存储。 |
| `candidates/` | 候选域类型：CandidateEnvelope、AdviceCandidate、CommunicationDraftCandidate、ExecutionIntentCandidate、CapabilityInvocationCandidate、MemoryWriteCandidate、TaskCandidate、CitedKnowledgeRef。 | 不产正式对象，不写 FormalMemory / FormalTask / SendResult / ExecutionResult。 |
| `spines/` | 三主线引用与 Intent Spine：TransactionRef、IntentEvent、IntentInboxItem、IntentClassification、IntentToCandidateResult、IntentReceipt、ReceiptLink、SceneFlowRunRef、DispatchPlanRef。 | 不做意图分类实现，不调用模型，不直接生成正式输出。 |
| `receipts/` | ReceiptEnvelope、AuditEnvelope。 | 不实现 append-only ledger、哈希链计算、回放查询或审计存储。 |
| `gates/` | 轻量 AccessDecision、OwnerVerdict。 | 不替代 `base.GateDecision` / `base.OwnerDecision`；不能用于正式动作主权裁定。 |
| `registry/` | RegistryRef、SkillRef、RegistrySlice、RegistrySliceItem、RegistrySliceBlockedRef、slice TTL / context refs helper。 | 不暴露 full registry，不保存 provider 实现，不提供真实解析服务。 |
| `readmodels/` | ReadModelEnvelope。 | 不持有真相源，不实现前端状态管理。 |
| `monitoring/` | MonitoringRun、MonitoringEvent、CollectorSnapshot、RedactionFinding、FaultIncident、SupportDiagnosticBundle、SupportUploadCandidate、SupportUploadAck、BuildSymbolManifest 等。 | 不采集日志、不上传诊断包、不符号化、不实现监控服务。 |
| `secrets/` | SecretRef、SensitivePayload。 | 不保存明文 secret、token、API key、cookie、private key。 |
| `events/` | ModuleEvent、IntentEvent alias。 | 不实现事件总线。 |
| `modules/` | ModuleContract 生命周期接口。 | 不实现模块启动 / 停止。 |
| 顶层 `contracts` 包 | `embed.go` 嵌入 schema bytes；`pack_knowledge_mount.go` 定义 KnowledgeScopeDeclaration / KnowledgeMountReadModel。 | 不实现 schema 校验器、知识挂载服务或 Pack lifecycle。 |

## JSON Schema

| Schema | 用途 | 主要消费方 |
| --- | --- | --- |
| `scene-pack-spec.schema.json` | 场景荚 / Domain Work Pack 规格：work_modes、allowed_mode_transitions、transaction_flow、workbench_surface、capability_requirements、receipt_rules、knowledge_scopes、export_policy。 | `truzhenos`、`truzhen-packs`、Pack Studio / CI |
| `scene-flow-spec.schema.json` | 场景流程图 / GateFlowSpec。 | `truzhenos` 06 Scene Flow、Pack Studio |
| `flow-view-spec.schema.json` | 流程视图投影。 | client layer、Pack Studio |
| `scene-runtime-plan-candidate.schema.json` | 场景运行时计划候选。 | `truzhenos` Scene Runtime / CI |
| `scene-studio-node-info.schema.json` | 制作台节点信息。 | Pack Studio / client layer |
| `scene-studio-workflow.schema.json` | 制作台工作流。 | Pack Studio / client layer |
| `visual-unit-spec.schema.json` | client layer 七类主权视觉单元（pod/object/capsule/candidate/execution/receipt/setting）封顶规格。 | client repo vendor / codegen / consistency test |
| `transaction-object-projection.schema.json` | 事务对象前端只读投影 DTO。 | client repo codegen / transaction object UI |
| `candidate-envelope.schema.json` | `candidates.CandidateEnvelope` 的 JSON 表达。 | client candidate card、CI |
| `receipt-envelope.schema.json` | `receipts.ReceiptEnvelope` 的 JSON 表达。 | client receipt card、CI |
| `spines/intent-event.schema.json` | IntentEvent JSON 表达。 | `truzhenos` 13 / 07 / 01、CI |
| `spines/intent-inbox-item.schema.json` | IntentInboxItem JSON 表达。 | `truzhenos` inbox / client projection |
| `spines/intent-classification.schema.json` | IntentClassification JSON 表达。 | `truzhenos` intent classifier / CI |
| `spines/intent-to-candidate-result.schema.json` | Intent fan-out 结果 JSON 表达。 | `truzhenos` candidate routing / CI |
| `spines/intent-receipt.schema.json` | IntentReceipt JSON 表达。 | `truzhenos` receipt candidate / CI |

## client layer 契约面

client layer（Web / Desktop / 后续移动端）面向本仓 schema 收敛跨边界 DTO，前端不手写后端稳定形状。

| 契约 | 前端对齐状态 |
| --- | --- |
| `visual-unit-spec.schema.json` | 已对齐：client 仓 vendor 副本 + codegen 类型 + 一致性测试。 |
| `transaction-object-projection.schema.json` | 已对齐：client 仓 codegen 到 generated 类型，业务类型消费生成物。 |
| `candidate-envelope.schema.json` / `receipt-envelope.schema.json` | 契约已就绪：对齐 Go struct 真相源，候选卡 / 回执卡面向；待 client 仓 vendor / codegen。 |
| Intent Spine 五件套 schema | 契约已就绪：面向 IntentEvent / inbox / classification / fan-out / receipt；下游接线状态按 `truzhenos` 与 client repo 记录为准。 |
| ReadModel 具体形状 / 秘书动作 / 其它 candidate 子类型 | 待契约化 + 对齐，前端形态稳定后统一推进。 |

消费机制：

- 有运行时常量数据的 DTO 用 vendor 副本 + 运行时一致性校验。
- 纯 type DTO 用 JSON Schema → TypeScript codegen；生成物只读勿手改。
- schema 改动必须同步说明 client vendor / codegen 是否需要更新。

## embed.go 覆盖纪律

`embed.go` 用 `//go:embed` 暴露 Go 服务和 CI 可直接消费的 canonical schema bytes。

当前状态：

- 已通过 `embed.go` 暴露：`scene-flow-spec.schema.json`、`scene-pack-spec.schema.json`、`flow-view-spec.schema.json`、`visual-unit-spec.schema.json`、`transaction-object-projection.schema.json`、`candidate-envelope.schema.json`、`receipt-envelope.schema.json`、Intent Spine 五件套 schema。
- 当前未通过 `embed.go` 暴露：`scene-runtime-plan-candidate.schema.json`、`scene-studio-node-info.schema.json`、`scene-studio-workflow.schema.json`。后续如裁定它们需要 Go API 直接消费，应补 embed 变量；如裁定只作为文件级 schema，应在 `embed.go` 附近留下说明。

改动规则：

- 新增 schema 时，默认同步新增 embed 变量。
- 若某 schema 被裁定不通过 Go API 暴露，必须在本文件和 `embed.go` 附近说明原因。
- 删除 schema 前必须先回 Owner，列下游消费影响。
- 改 schema 路径必须同步 `embed.go`、README、MODULES 和下游引用。

检查命令见 `AGENTS.md` 的“schema embed 覆盖检查”。

## 状态口径

- `契约已定`：本仓 Go type / schema 已落地。
- `已接线`：下游真实消费并通过对应验证。
- `已验收`：下游消费、schema 校验、反向依赖、兼容测试均有证据。
- `已发布`：Owner 授权 tag / module 版本发布。
