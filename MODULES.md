# truzhen-contracts 模块与 schema 清单

每个子包只声明跨边界的数据形状（类型 / 接口 / 常量），或提供无外部副作用的契约校验 / ref 派生 helper。本仓是六仓协同中的 **Pack / Candidate / Receipt / Surface / ReadModel / Cloud schema** 权威源（cloud 契约面当前为治理清单，见下）。

## Truzhen 底层逻辑指向

权威总纲：`/Users/li/Documents/truzhen-contracts/TRUZHEN_PHILOSOPHY.md`（远端 `github.com/truzhen/contracts` 根目录同名文件）。本仓负责把总纲中的主权链、候选 / 正式隔离、回执、跨仓依赖方向和数据形状固化为可消费契约；其它仓应引用本文件和总纲，不另造一套哲学真相源。

### 本仓映射表

| 总纲原则 | 本仓落点 | 边界 |
| --- | --- | --- |
| AI 是工具智能，只能生成可治理候选 | `candidates/`、`base.GateCandidateEnvelope`、候选类 schema | 只声明 Candidate 形状，不产 Formal Record。 |
| Owner 才能做 Decision；Decision 可即时确认，也可来自预授权策略 | `base.OwnerDecision`、`OwnerDelegationGrant`、授权模式、GateDecision 相关类型 | 只定义决定与授权的形状，不实现弹窗、策略引擎或自动放行逻辑。 |
| 正式动作必须过闸并留下可回放回执 | `base/`、`gates/`、`receipts/`、`receipt-envelope.schema.json` | 只定义 Gate / Receipt 契约，不实现 ledger、执行、发送或存储。 |
| 六仓依赖方向单向、不可逆 | 本仓只用标准库；下游通过 Go SDK / JSON Schema 消费 | 不 import `truzhenos`、`truzhen-packs`、`truzhen-cloud`、client 或 provider 实现。 |
| Pack 是行业经验给 AI 装上的边界 | `scene-pack-spec.schema.json`、`scene-flow-spec.schema.json`、ProviderRequirement 相关字段 | schema 表达边界，不替真实客户证据设计投机 Pack。 |
| ReadModel / Surface 只是展示投影 | `readmodels/`、`visual-unit-spec.schema.json`、`transaction-object-projection.schema.json` | 不把投影、页面或 DTO 写成事实源。 |

### 禁止误读清单

- 不得把总纲理解成“可以随意新增 schema / Go 类型”；新增契约仍必须有真实消费方、SemVer 影响清单和 Owner 裁定。
- 不得把 `OwnerDecision` 理解成“所有动作都必须即时弹窗”；预授权策略也是 Owner 的决定，只是拍板时点提前。
- 不得把本仓当运行时真相源；本仓不持有 Base Gate、Receipt Ledger、Provider readiness、Pack enabled state、License state 或前端状态。
- 不得让 AI 未经 Owner 确认就使自己起草的治理规则、风险分级或 Pack manifest 生效。
- 不得把 Provider 示例、OS Agent 示例或云端例子读成近期集成承诺；本仓只收敛跨边界形状。

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
| `market/` | 市场表面契约（§18-A，2026-07-02 v0.4.0 新增；2026-07-03 v0.5.0 补 SessionProjection、session/payable/role 枚举和作者端 DTO / 枚举；2026-07-07 v0.7.0 新增 canonical PackManifest / ProviderRequirement / PackSoftwareRequirement / SoftwareResolutionLock；2026-07-08 v0.8.0 补 resolver MVP lock 结果 `install_required` / `version_conflict` / `isolation_required`）：SessionHeader、LoginRequest / LoginResponse、作者认证 / 收益 / 提现 / 上传 ReadModel、Pack manifest、软件依赖声明、resolver lock 形状、市场表面端点路径常量与路径构造器（LicenseOrderPath / WithdrawalCancelPath / PackDownloadPath）、AdminForwardAllowlist / AdminPathAllowed（admin 转发硬 allowlist 主权边界）。消费方：truzhenos 17-market / 02 registry、truzhen-cloud 03 上传链、client 软件目录投影。 | 不实现代理转发、不签发订单 / 价格 / 权益（服务端真相唯一在 truzhen-cloud），不解析本机 provider，不保存本机路径 / 端口 / secret / runtime state，不持会话状态。 |
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
| `receipt-envelope.schema.json` | `receipts.ReceiptEnvelope` 的 JSON 表达。 含可选 `actual_edits`（执行后事实，v0.10.0 additive，Owner O-1~O-4 裁定 2026-07-10）。 | client receipt card、CI |
| `pack-manifest.schema.json` | 云端上传与 Pack 分发可校验的 canonical manifest，含 `software_requirements` 与可选 `lifecycle_status`（八档中文枚举，v0.9.0 additive，Owner 2026-07-10 裁定）。 | `truzhen-cloud` 上传校验、`truzhen-packs` CI、`truzhenos` Pack loader |
| `provider-requirement.schema.json` | Pack 声明 provider 能力需求的 canonical 形状。 | `truzhen-packs`、`truzhenos` readiness / resolver |
| `software-resolution-lock.schema.json` | `truzhenos` resolver 产出的本机软件依赖解析锁。 | `truzhenos`、client 软件目录、cloud 只读投影 |
| `monitoring/monitoring-event.schema.json` | `monitoring.MonitoringEvent` 的 JSON 表达，含 additive `error_code` 稳定错误码字段。 | `truzhenos` 统一监控 / CI / support bundle |
| `monitoring/fault-incident.schema.json` | `monitoring.FaultIncident` / `FaultSignature` 的 JSON 表达，约束 `error_code` 格式。 | `truzhenos` fault classifier / support bundle / cloud symbolication |
| `spines/intent-event.schema.json` | IntentEvent JSON 表达。 | `truzhenos` 13 / 07 / 01、CI |
| `spines/intent-inbox-item.schema.json` | IntentInboxItem JSON 表达。 | `truzhenos` inbox / client projection |
| `spines/intent-classification.schema.json` | IntentClassification JSON 表达。 | `truzhenos` intent classifier / CI |
| `spines/intent-to-candidate-result.schema.json` | Intent fan-out 结果 JSON 表达。 | `truzhenos` candidate routing / CI |
| `spines/intent-receipt.schema.json` | IntentReceipt JSON 表达。 | `truzhenos` receipt candidate / CI |

## cloud 契约面（truzhen-cloud / truzhenos / client 协同）

cloud 契约只定义跨仓形状，不实现云服务。`truzhen-cloud` 是官方云端真相源，`truzhenos` 是本地 Cloud proxy / License Gate 消费端，client repo 只消费 ReadModel / DTO。七类 cloud 契约必须在本仓收敛后再被下游实现或 codegen：

| 契约类 | 职责边界 |
|---|---|
| `Entitlement` | 用户 / 组织 / 设备对 Pack、能力、Release 或服务的授权权益形状；不保存支付流水实现。 |
| `License` | License 状态、有效期、席位、激活与撤销形状；真实核验服务归 `truzhen-cloud`，本地只消费裁定结果。 |
| `Payment` | 支付订单、支付结果、退款 / 取消、webhook 事件形状；支付网关实现与密钥不进本仓。 |
| `PackListing` | 云市场商品、版本、价格、作者、分发状态与审核状态形状；Pack manifest 不是商品真相源。 |
| `Session` | 云端登录态、Session ID、续期、退出、设备绑定相关形状；raw token / password 不进本仓。 |
| `Release` | 云端版本、下载、灰度、校验、回滚与客户端升级提示形状；二进制产物不进本仓。 |
| `WebSurface` | 官网、市场页、作者后台、运营后台、支付结果页等官方云端网页的路由、状态与展示 DTO 形状；页面实现归 `truzhen-cloud`。 |

> 当前本节是治理契约清单。新增具体 schema / Go 包时必须单独列兼容策略、SemVer 影响、下游旧路径迁移清单和反向依赖检查。

## client layer 契约面（前端面向收敛）

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

- 已通过 `embed.go` 暴露：`scene-flow-spec.schema.json`、`scene-pack-spec.schema.json`、`flow-view-spec.schema.json`、`visual-unit-spec.schema.json`、`transaction-object-projection.schema.json`、`candidate-envelope.schema.json`、`receipt-envelope.schema.json`、`pack-manifest.schema.json`、`provider-requirement.schema.json`、`software-resolution-lock.schema.json`、`monitoring/monitoring-event.schema.json`、`monitoring/fault-incident.schema.json`、Intent Spine 五件套 schema。
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

> 2026-07-02（v0.3.1）：删除 `events/` 与 `modules/` 两个零消费占位包（三仓 grep 复核零 import；按「无消费方=投机」先例）。若未来需要事件信封/模块生命周期契约，由真实消费方驱动重建。
> 2026-07-02（v0.4.0）：新增 `market/` 市场表面契约包（§18-A，tag `ff042d0`）；消费方 truzhenos 17-market 薄代理已接，cloud 03 对齐为后续小卡。
> 2026-07-03（v0.5.0）：`market/` 补云市场 SessionProjection、session/payable/role 枚举，以及作者认证、收益、提现、Pack 上传 manifest / ReadModel、product_kind、trust_verify_status 契约形状；用于登录态云端校验、可支付性投射和作者端商品化闭环。兼容新增，不删除旧 Login DTO。
> 2026-07-07（v0.7.0）：新增 canonical `PackManifest`、`ProviderRequirement`、`PackSoftwareRequirement`、`SoftwareResolutionLock` 与对应 schema；`scene-pack-spec.schema.json` additive 增 `software_requirements`，旧 `external_software_refs` 保留并标 deprecated。消费方按 contracts -> software -> truzhenos -> packs -> cloud -> client 顺序接线。
> 2026-07-08（v0.8.0）：`SoftwareResolutionLock` enum additive 增 `install_required`、`version_conflict`、`isolation_required`，用于本地 resolver + 用户侧 lock 的最小闭环；仍禁止把用户本机事实写入 contracts。
> 2026-07-10（登记待裁，无版本变更）：`pack-manifest.schema.json` 的 `kind` 枚举含 `skill_bundle`，与 `TRUZHEN_PHILOSOPHY.md` §7「三类 Pack 封顶」存在表述张力，且该值未登记于 truzhenos `NAMING_STANDARD.md`。Owner 2026-07-10 裁定本轮只登记、另立任务裁定其地位（合法子类需补哲学 / 命名表述；历史遗留则标 deprecated——删枚举值属 major）。本轮语义治理条文（CONTRACTS_GOVERNANCE §6.1 / PHILOSOPHY §7.1）为纯文档，未改 schema，不 bump VERSION。
> 2026-07-10 晚（**已裁定闭环**，无版本变更）：Owner 裁定 `skill_bundle` = 智能体执行过程的工具（技能包，truzhenos NAMING_STANDARD 已定名），**不是第四种 Pack**——是智能体工具集资产的分发封装 kind，枚举保留、语义澄清见 `TRUZHEN_PHILOSOPHY.md` §7 澄清段。落地：truzhenos 接线 13 任务级技能切片端点（按任务装载）、client 管理面迁团队设置「技能包」tab。表述张力解除，纯文档，不 bump VERSION。

## 包体积与完善状态（2026-07-03 集成分支实测，v0.5.0）

> 核心源码 LOC（非测试）；本仓无 `FEATURE_LEDGER.md`，模块进度以本表 + 上文子包职责为准（契约仓账本即 MODULES.md）。详见 truzhenos `docs/status/six-repo-module-audit-20260702.md` §17.2。

| 包 | 核心源码 LOC | 功能 | 待解决问题 |
| --- | ---: | --- | --- |
| `base/` | 1,695 | 主权核心契约（Owner / Policy / Gate / OwnerDecision / Formalization） | ✅ |
| `candidates/` | 140 | AI / Pack / 模块候选域类型 | 🔧 缺单测 |
| `spines/` | 236 | 三主线引用 + Intent 五件套（含 5 schema） | 🔧 缺单测 |
| `receipts/` | 21 | 回执 / 审计信封 | 🔧 缺 schema JSON 侧 |
| `registry/` | 91 | RegistryRef / RegistrySlice（类型强制） | ✅ |
| `monitoring/` | 统一监控 / 诊断 / 故障包契约 + 2 schema | `MonitoringEvent.error_code` 已契约化；schema / embed 已暴露 |
| `secrets/` | 12 | SecretRef / SensitivePayload | ✅ |
| `gates/` | 15 | 轻量 AccessDecision / OwnerVerdict | 🔧 |
| `readmodels/` | 7 | ReadModelEnvelope | 🔧 具体形状待契约化 |
| `market/` | 约 330（+430 测试） | 市场表面契约：SessionHeader / Login DTO / SessionProjection / session-payable-role 枚举 / 作者端 DTO / 表面路径 / admin 硬 allowlist，黄金断言守护 | 🔧 v0.5.0 契约已落地；下游需在集成分支先吸收 / 发布 contracts 后再编译期引用 |
| 根级 | 全仓 30 个 `*.schema.json`；其中 27 个由 `embed.go` 暴露，3 个 scene schema 未 embed | schema 嵌入 + 版本漂移门禁（`check-version-drift.sh`）+ 破坏性变更与 Go↔Schema 配对门禁（`contracts-check.sh`，type 变更一律判 breaking，Owner 2026-07-10 R-a 裁定） | ✅ |
