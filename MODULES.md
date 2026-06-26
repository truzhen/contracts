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
