# truzhen-contracts 子包清单

每个子包只声明跨边界的数据形状（类型 / 接口 / 常量），无实现、无副作用。

| 子包 | 职责（一句话） |
|---|---|
| `base/` | Base 主权核心契约：授权模式语义、Owner 身份与指令（ActorContext/OwnerDirective）、PolicySet/PolicyKernel/PolicySnapshot、GateCandidateEnvelope/GateRequest/GateDecisionTrace/GateReceiptCandidate 等门控信封。 |
| `candidates/` | 候选域类型：AI Proposer 产出的各类候选（Advice/CommunicationDraft/ExecutionIntent/CapabilityInvocation/MemoryWrite/Task）与 CandidateEnvelope、CitedKnowledgeRef。 |
| `gates/` | 门控裁定结果类型：AccessDecision、OwnerVerdict。 |
| `receipts/` | 回执与审计信封：ReceiptEnvelope、AuditEnvelope（Evidence Spine 落点）。 |
| `spines/` | 三主线引用：IntentEvent/IntentInboxItem/IntentClassification（意图主线）、TransactionRef（事务主线）、ReceiptLink（证据主线）、SceneFlowRunRef、DispatchPlanRef。 |
| `registry/` | 注册中心契约切片：RegistrySlice/RegistrySliceItem、RegistryRef、SkillRef。 |
| `readmodels/` | 前端投影信封：ReadModelEnvelope（ReadModel ≠ 真相源）。 |
| `monitoring/` | 统一监控契约：MonitoringRun/MonitoringEvent、CollectorSnapshot、RedactionFinding、SupportDiagnosticBundle、FaultIncident 等。 |
| `secrets/` | secret **引用**契约：SecretRef、SensitivePayload。只声明「对密文的引用」形状，**不含任何真凭据值**。 |
| `events/` | 模块事件类型：ModuleEvent。 |
| `modules/` | 模块契约描述：ModuleContract。 |
| 顶层 | `embed.go`（嵌入下列 `*.schema.json` 供 Go 服务与 CI 校验）、`pack_knowledge_mount.go`（Pack 知识挂载契约）。 |

## JSON Schema

| Schema | 用途 |
|---|---|
| `scene-pack-spec.schema.json` | 场景荚（Domain Work Pack）规格 v2：work_modes / allowed_mode_transitions / transaction_flow / workbench_surface / surface_slot 等。 |
| `scene-flow-spec.schema.json` | 场景流程图（GateFlowSpec）。 |
| `flow-view-spec.schema.json` | 流程视图投影。 |
| `scene-runtime-plan-candidate.schema.json` | 场景运行时计划候选。 |
| `scene-studio-node-info.schema.json` | 制作台节点信息。 |
| `scene-studio-workflow.schema.json` | 制作台工作流。 |
| `visual-unit-spec.schema.json` | **client layer**：前端 7 类主权视觉单元（pod/object/capsule/candidate/execution/receipt/setting）封顶规格契约。 |
| `transaction-object-projection.schema.json` | **client layer**：事务对象（05 BusinessObject）前端只读投影 DTO 契约。 |
| `candidate-envelope.schema.json` | **client layer**：候选统一包装（`candidates.CandidateEnvelope` 的 JSON 表达，候选卡面向）。 |
| `receipt-envelope.schema.json` | **client layer**：回执链包装（`receipts.ReceiptEnvelope` 的 JSON 表达，回执卡面向）。 |

## client layer 契约面（前端面向收敛）

client layer（Web / 桌面等多端前端）面向以下 contracts 契约收敛跨边界 DTO（前端不手写后端形状，面向契约单一来源）：

| 契约 | 前端对齐状态 |
|---|---|
| `visual-unit-spec.schema.json` | ✅ 已对齐（前端 vendor 副本 + 一致性测试） |
| `transaction-object-projection.schema.json` | ⏳ 契约已就绪，待前端对齐（codegen 单元） |
| `candidate-envelope.schema.json` / `receipt-envelope.schema.json` | ✅ 契约就绪（对齐 contracts Go struct 真相源，候选卡/回执卡面向），待前端对齐 |
| ReadModel 具体形状 / 秘书动作 / 其他 candidate 子类型等 | ⏳ 待契约化 + 对齐（前端形态稳定后统一推） |

> 前端消费机制：有运行时常量数据的 DTO（如视觉单元规格表）用 vendor 副本 + 运行时一致性校验；纯 type DTO 待 codegen（JSON Schema → TS type 生成）。
