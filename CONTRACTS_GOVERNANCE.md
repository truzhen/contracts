# truzhen-contracts 契约治理总纲

本文记录 `truzhen-contracts` 的仓库身份、模块职能、契约边界和变更纪律。`AGENTS.md` 是 Agent 开工入口；本文是更完整的治理说明。

## 1. 总定位

`truzhen-contracts` 是 Truzhen 的开源契约层 SDK，模块路径为 `github.com/truzhen/contracts`。它的职责是把跨仓、跨层、跨运行时的数据形状固定下来，让基座、Pack、client layer 和 CI 能围绕同一契约演进。

本仓定义：

- Go SDK 类型、接口、枚举和常量。
- JSON Schema 机器契约。
- schema embed 入口。
- 少量无外部副作用的确定性校验、默认信封构造和 ref 派生 helper。

本仓不实现：

- Base Gate、Receipt Ledger、Gateway、provider、runtime、sidecar、前端 UI、Pack 安装器、外部软件 registry、数据库、网络服务、文件读写、真实执行。

## 2. 架构位置

Truzhen 当前按六仓协同：

```text
truzhenos          implements      truzhen-contracts       faces       truzhen-packs
私有基座 / 实现  ─────────────▶  开源契约 / 形状权威源  ◀────────  开放包层 / 面向契约
```

client repo 通过 vendor / codegen 消费 JSON Schema。外部 provider / sidecar 事实归 `truzhen-software`。本仓位于六仓中央，但只拥有“形状事实”，不拥有运行事实。

另有两落点：`truzhen-cloud` 实现本仓 cloud 契约面，是官方云服务 / 网页 / 支付 / License / Entitlement / PackListing / Session / Release 事实的真相源；client repo（`truzhen-client-web-desktop`）只通过 vendor / codegen 消费本仓 schema，永不定义契约、永不充当真相源。

## 3. 真相源原则

新增类型、字段或 schema 前必须回答：

1. 这个事实归谁？
2. 本仓只是声明形状，还是误把实现事实搬进来了？
3. 下游谁消费？
4. 是否需要 schema、Go struct、README、MODULES、embed、client codegen 同步？

默认归属：

- 主权裁定事实：`truzhenos` Base。
- 回执账本事实：`truzhenos` Receipt Ledger。
- Pack 内容事实：`truzhen-packs`。
- Provider / sidecar / runtime 资源事实：`truzhen-software` 或基座 registry。
- 云端 Entitlement / License / Payment / PackListing / Session / Release 事实：`truzhen-cloud`。
- 前端展示状态：client repo ReadModel 消费；ReadModel 不是真相源。
- 跨仓数据形状：本仓。

## 4. 核心契约域

### 4.1 Base / Gate / Receipt Candidate

`base/` 定义主权链路的稳定形状：GateRequest、GateCandidateEnvelope、GateDecision、GateReceiptCandidate、PolicySnapshot、OwnerDecision、BaseFormalizationGrant、授权模式、委托授权、Artifact 留痕与过闸边界。

治理要求：

- Base 正式裁定只在基座实现，本仓只声明裁定形状。
- helper 只能校验 ref 完整性、风险硬地板、secret-ish 字段等契约边界。
- 不得加入 orchestrator、policy engine、ledger append 或 provider call。

### 4.2 Candidate

`candidates/` 定义所有 AI / Pack / 模块提出的候选类型。候选不是正式对象。

治理要求：

- 默认 candidate-only / non-formal。
- Advice 可以引用知识 ref，但不能自断法律 / 业务效力。
- ExecutionIntent、CommunicationDraft、MemoryWrite、Task 等都只是候选。

### 4.3 Spines

`spines/` 定义 Transaction / Intent / Evidence 三主线引用。Intent Spine 五件套要求所有输入先形成 IntentEvent，再分类和 fan-out 到候选。

治理要求：

- IntentClassification 只能输出候选路由，不得直接产 FormalTask、FormalMemory、SendResult、ExecutionResult。
- 八类候选目标是封闭集合；新增第九类必须先回 Owner。

### 4.4 Receipt / Audit

`receipts/` 定义回执信封与审计信封。

治理要求：

- ReceiptEnvelope 是证据链形状，不是 ledger 实现。
- Hash / sequence 字段属于契约字段；计算和 append 在基座。

### 4.5 Registry Slice

`registry/` 定义被裁剪、mask、rank、audit 后的 RegistrySlice。

治理要求：

- Agent / Model 不得消费 full registry。
- blocked refs 不得静默消失，必须可审计。

### 4.6 ReadModel / Surface

ReadModel 和 Surface schema 只面向展示、投影和 client layer codegen。

治理要求：

- ReadModel 不是真相源。
- 前端不得手写稳定 DTO 绕过 contracts。
- 视觉单元、事务对象投影、候选卡、回执卡等 schema 改动必须同步 client vendor / codegen 状态说明。

### 4.7 Monitoring

`monitoring/` 定义统一监控与诊断包形状。

治理要求：

- 不另起日志上传或诊断包格式；新诊断对象优先纳入已有 monitoring 契约。
- SupportUploadCandidate 仍是候选，不等于真实上传成功。

### 4.8 Secrets

`secrets/` 只定义 secret 引用。

治理要求：

- 明文 secret、token、API key、cookie、private key、terminal_sn、激活码永不入仓。
- schema 示例也不得包含真实凭据。

## 5. JSON Schema 治理

Schema 文件分三类：

1. Go struct 的 JSON 表达：如 `candidate-envelope.schema.json`、`receipt-envelope.schema.json`。
2. client layer / Surface 契约：如 `visual-unit-spec.schema.json`、`transaction-object-projection.schema.json`。
3. Pack / Scene / Studio / Intent 机器契约：如 `scene-pack-spec.schema.json`、`scene-runtime-plan-candidate.schema.json`、`spines/intent-*.schema.json`。

修改 schema 必须：

- 保持合法 JSON。
- 明确 `$schema`、`title`、`type`、`required`、`additionalProperties`。
- 谨慎改 enum 和 required。
- 同步 `MODULES.md` schema 清单。
- 判断是否应通过 `embed.go` 暴露。
- 若 schema 供 client codegen，说明 vendor / generated 类型是否需要更新。

软件依赖相关 schema 的附加边界：`pack-manifest.schema.json` 只允许 Pack 声明底层软件需求；`software-resolution-lock.schema.json` 只允许表达 `truzhenos` resolver 产出的用户侧 lock。schema 中不得加入 raw local path、raw endpoint、credential、DB、模型权重、镜像、端口或 runtime state 字段；这些事实分别归 `truzhen-software` registry 与 `truzhenos` 用户侧状态。

## 6. 版本与兼容

本仓遵守 SemVer。契约层没有“内部小改”豁免；一旦下游依赖，字段就是跨仓边界。

- patch：注释、文档、description、非语义修正。
- minor：新增可选字段、新增非破坏 schema、新增可选 helper。
- major：删字段、改字段类型、改必填、改 enum 语义、改 JSON tag、改 Candidate/Formal/Gate/Receipt 主权语义。

`v0.x` 不得被用来掩盖破坏性变更。破坏性变更必须列影响面、迁移策略、下游同步顺序和回滚方案。

改 `*.schema.json` 必先 bump 仓根 `VERSION`，CI 以 `scripts/check-version-drift.sh` 强制（版本漂移即红）。

### 6.1 语义与知识资产治理原则（2026-07-10）

「领域语义与隐性业务知识」（行业概念、口径、关系、判断约束，俗称"暗知识"）的跨仓分工固定为：**contracts 只定形状，packs 声明内容，truzhenos os-09 持正式知识真相与挂载 / 裁剪，前端只投影**。

- 契约层对语义只做形状声明，不做业务规则解释：`scene-pack-spec.schema.json` 的 `knowledge_scopes[]`（`knowledge_kinds` 枚举 `law_article / sop / case / glossary / checklist / index`）、`business_object_schema_refs` 与 `pack_knowledge_mount.go`（`KnowledgeScopeDeclaration` / `KnowledgeMountReadModel`）已是语义资产的既有承载，不为"语义"另立平行 schema、`semantic_model` / `ontology` 字段或独立契约域。
- `knowledge_kinds` 枚举扩展按本节 §6 既有 SemVer 纪律走 minor bump 并通知下游 codegen，不为"语义"开豁免；新增语义相关字段必须由真实消费方驱动（无消费方 = 投机，按 v0.3.1 删零消费包先例）。
- 语义资产是三类 Pack 的横向内容维度，不改变三类 Pack 封顶（见 `TRUZHEN_PHILOSOPHY.md` §7）；不新增第四种 Pack、中央语义层或 Semantic Runtime 契约。

## 7. 依赖纪律

允许：

- Go 标准库。
- 本仓内部包。

需要 Owner 裁定：

- 新增第三方依赖。
- 引入 codegen 工具链。
- 调整 module path。

禁止：

- import `github.com/lights314/truzhenos`。
- import `github.com/truzhen/packs`。
- import `github.com/truzhen/truzhen-cloud`。
- import client repo。
- 依赖 provider / sidecar / external software package。

## 8. Helper 纪律

契约 helper 的边界是“帮助消费方正确理解契约”，不是“替消费方实现业务”。

允许：

- `Validate*` 检查必填字段、枚举、引用绑定。
- `New*` 生成默认 candidate-only/non-formal 信封或稳定 ref。
- 纯函数计算稳定 hash / ref。
- 不访问外部世界的状态判断。

禁止：

- 调用 DB、HTTP、文件系统、provider、sidecar、shell。
- 读写环境变量或用户本地路径。
- 开 goroutine、持久缓存、全局调度。
- 生成真实 OwnerDecision / GateDecision 的业务裁定。
- 用随机数或隐式时间制造正式事实。

## 9. 验收口径

不同改动对应不同证明：

- 改文档：`git diff --check`，必要时跑基础 Go / schema 检查。
- 改 Go 类型：`go build ./... && go test ./... && go vet ./...` + 反向依赖检查。
- 改 schema：schema JSON 合法性 + embed 覆盖检查 + 下游消费说明。
- 改 Candidate/Gate/Receipt/ReadModel：必须列契约影响和兼容策略。
- 改版本 / 发布：必须有 Owner 授权。

本仓当前没有统一测试文件时，`go test` 通过只说明编译通过；不能把“无测试文件”说成“行为已覆盖”。

## 10. 生命周期术语

统一使用：

`想法 -> 设计中 -> 契约已定 -> 已实现 -> 已接线 -> 已验收 -> 已发布 -> 已弃用`

契约仓常见状态：

- `契约已定`：本仓 Go type / schema 已落地。
- `已接线`：下游真实消费并通过验证。
- `已验收`：跨仓消费、兼容、schema、反向依赖均有证据。
- `已发布`：Owner 授权发版并生成 tag / module 版本。
