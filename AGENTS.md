# AGENTS.md — truzhen-contracts

本文件是 Agent 维护 `truzhen-contracts` 的工作纪律入口。目标是让 Agent **只读本文件即可安全上手**；需要更完整的边界解释时，再读 `CONTRACTS_GOVERNANCE.md`、`README.md`、`MODULES.md`、`CLAUDE.md`。

## 0. AI 维护速查（6 要素）

1. **本仓职责** — Truzhen 契约层 SDK（`github.com/truzhen/contracts`，开源 Apache-2.0），定义六仓之间跨边界的数据形状：Pack / Candidate / Receipt / Surface / ReadModel schema、候选信封、门控裁定、回执 / 审计、注册切片、监控事件、三主线引用，以及机器可校验的 `*.schema.json`（cloud Entitlement / License / Payment / PackListing / Session / Release / WebSurface 契约面为治理清单，见 `MODULES.md`，具体 schema 待真实消费方出现再落）。当前公开消费版本为 `v0.3.0`，基座通过 `go.mod require github.com/truzhen/contracts v0.3.0` 消费。
2. **本仓不是什么** — 不是基座 `truzhenos`，不拥有 Base Gate / Receipt Ledger / Gateway / provider / runtime 实现；不是 `truzhen-packs`，不保存行业 Pack 数据、安装脚本或 Pack 内容；不是 client repo，不实现 React/Tauri UI；不是 `truzhen-software`，不记录本机外部软件安装事实或 sidecar 运行态；不是 `truzhen-cloud`，不实现云端 server、支付 webhook 或 License / Entitlement 服务。
3. **允许内容** — Go 类型、接口、常量、枚举、JSON tag、JSON Schema、schema embed、无外部副作用的确定性校验 / ref 派生 / helper。helper 只能表达契约边界，不得访问 DB、网络、文件系统、provider、真实执行环境或用户资产。
4. **禁止内容** — 不引入 DB / 网络 / 文件 I/O / 并发运行时 / provider 调用 / 真实执行；不 import `github.com/lights314/truzhenos`、`github.com/truzhen/packs`、`github.com/truzhen/truzhen-cloud` 或 client / provider 实现仓；不把 Pack 数据、前端源码、安装脚本、provider registry、raw secret、token、terminal_sn、激活码放进本仓；不无意识删字段、改必填、改语义。
5. **验证入口** — `go build ./... && go test ./... && go vet ./...`；反向依赖检查；全部 schema JSON 解析检查。修改 schema / embed 时必须同时检查 `embed.go` 覆盖关系。
6. **必须回 Owner** — 任何破坏性契约变更、SemVer 版本策略变化、新增 / 删除子包、新增 / 删除 schema、schema 必填字段变化、跨仓边界或依赖方向调整、`git push` / tag / release。

## 1. 首读要求

新任务开工先读：

1. `AGENTS.md`
2. `CONTRACTS_GOVERNANCE.md`
3. `README.md`
4. `MODULES.md`
5. `CLAUDE.md`
6. `CONTRIBUTING.md`

按任务范围补读：

- 改 Go 契约类型：读对应子包文件与相关 JSON Schema。
- 改 `*.schema.json`：读 `embed.go`、`MODULES.md` 的 schema 清单、对应 Go struct 或 schema 消费说明。
- 改 scene-pack / scene-runtime / scene-studio schema：同时只读参考基座设计文档 `/Users/li/Documents/truzhenos/docs/design/scene-pack-vertical-profession-workbench-upgrade-20260626.md`（若存在），不得修改基座仓，除非 Owner 另行授权。
- 改 client layer schema：同时只读核对 `/Users/li/Documents/truzhen-client-web-desktop/src/contracts/CONTRACTS_VENDOR.md`（若存在）和 client vendor / codegen 消费方式；不得修改 client repo，除非 Owner 另行授权。

## 2. 五落点身份边界

| 落点 | 权威位置 | 负责 | 不负责 |
| --- | --- | --- | --- |
| 主权、闸门、回执、网关、运行时、装载器 | `/Users/li/Documents/truzhenos` / `github.com/lights314/truzhenos` | 实现 contracts，持有 Base / Gateway / Receipt / runtime | 不把 contracts 当内部包私改 |
| Pack / Candidate / Receipt / Surface / ReadModel schema | `/Users/li/Documents/truzhen-contracts` / `github.com/truzhen/contracts` | 本仓；跨仓数据形状与 schema 权威源 | 不实现运行时、不保存 Pack 数据、不执行 provider |
| 行业工作台、流程、角色、知识、能力引用 | `/Users/li/Documents/truzhen-packs` / `github.com/truzhen/packs` | folder pack、manifest、role/knowledge/capability refs、安装脚本 | 不持 Base 主权、不实现 provider |
| Baserow、Frappe、OCR、IM、执行 sidecar | `/Users/li/Documents/truzhen-software` | 本机外部软件 / provider / sidecar 事实 | 不是 Git 主仓、不承载契约 |
| Web、桌面、手机、小程序 | `/Users/li/Documents/truzhen-client-web-desktop` / 后续 client repo | 前端源码、Tauri 壳、client layer DTO 消费 | 不直接写正式对象、不绕过 ReadModel / Candidate / Gateway |

跨仓读取、修改、测试、提交或推送必须遵守当前用户授权范围。本仓任务默认只修改 `truzhen-contracts`；参考其它仓只能只读，且要说明原因和影响范围。

## 3. 模块地图

### Go 子包

| 子包 | 职能 | 纪律 |
| --- | --- | --- |
| `base/` | Base 主权核心契约：Owner / Actor、Policy、GateCandidateEnvelope、GateRequest、GateDecision、GateReceiptCandidate、FormalizationGrant、授权模式、委托、Artifact 留痕 / 过闸边界、Gateway adapter request。 | 只能表达形状、枚举、确定性校验和 ref 派生；不得实现 Base orchestrator、Receipt Ledger、provider 调用或真实执行。新增 helper 不得隐式引入外部副作用。 |
| `candidates/` | AI / Pack / 模块产出的候选域类型：Advice、CommunicationDraft、ExecutionIntent、CapabilityInvocation、MemoryWrite、Task、CandidateEnvelope、CitedKnowledgeRef。 | Candidate 默认 `candidate_only=true`、`non_formal=true`；不得引入 Formal 写入或真实动作结果。 |
| `spines/` | 三主线引用与 Intent Spine 五件套：TransactionRef、IntentEvent、IntentInboxItem、IntentClassification、IntentToCandidateResult、IntentReceipt、ReceiptLink、SceneFlowRunRef、DispatchPlanRef。 | 分类只产候选路由，不产正式输出；八类候选目标不得随意扩展。 |
| `receipts/` | ReceiptEnvelope、AuditEnvelope 等证据链形状。 | 不实现账本 append、哈希链计算、存储或回放。 |
| `gates/` | AccessDecision、OwnerVerdict 轻量裁定形状。 | 轻量决策不得冒充 `base.GateDecision` / `base.OwnerDecision`。正式动作必须走 `base/` 契约。 |
| `registry/` | RegistryRef、SkillRef、RegistrySlice、masked / ranked / audited context slice。 | Agent / Model 只消费 slice；不得暴露 full registry 或 raw secret。 |
| `readmodels/` | ReadModelEnvelope。 | ReadModel 只用于展示投影，永不是真相源。 |
| `monitoring/` | MonitoringRun/Event、CollectorSnapshot、RedactionFinding、FaultIncident、SupportDiagnosticBundle、SupportUploadCandidate、BuildSymbolManifest 等统一诊断契约。 | 只定义诊断数据形状；不实现日志采集、上传或符号化服务。 |
| `secrets/` | SecretRef、SensitivePayload 等 secret 引用形状。 | 永不包含明文凭据。 |
| `events/` | ModuleEvent 与 IntentEvent alias。 | 只做事件信封，不实现事件总线。 |
| `modules/` | ModuleContract 生命周期接口。 | 不定义具体模块实现。 |
| 顶层包 | `embed.go`、`pack_knowledge_mount.go`。 | embed 只暴露 canonical schema bytes；知识挂载只定义声明和 ReadModel。 |

### JSON Schema

Schema 是跨仓机器契约。新增或修改 schema 时必须说明：

- 真相源是谁：Go struct、schema 文件、还是外部协议。
- 谁消费：`truzhenos`、`truzhen-packs`、client repo、CI、codegen。
- 是否破坏兼容：新增可选字段 / 新增必填字段 / 删除字段 / 改语义 / 改 enum。
- 是否需要同步 `embed.go`、`MODULES.md`、README、下游 vendor / codegen。

## 4. 契约变更纪律

### SemVer 判定

| 变更 | 版本策略 | Owner 裁定 |
| --- | --- | --- |
| 注释、文档、schema description 修正 | patch | 一般不需要，除非改语义 |
| 新增可选字段、可选 enum、只读投影字段 | minor | 需要说明下游兼容性 |
| 新增必填字段、删除字段、改字段类型、改 enum 含义、改 JSON tag、改默认语义 | major / 破坏性 | 必须先回 Owner，列影响清单和迁移方案 |
| 新增 / 删除 Go 子包或 schema 文件 | 至少 minor，可能 major | 必须先回 Owner |
| 改 Candidate/Formal、Gate、Receipt、ReadModel、ProviderRequirement、Surface schema 边界 | 视为跨仓边界变更 | 必须先出影响清单 |

### 兼容策略

- 不删除已发布字段；废弃用 `deprecated` 说明和新字段并存。
- 新字段默认可选，除非已有跨仓计划和迁移窗口。
- Go struct 与 JSON Schema 同步时，要检查 JSON tag、required、enum、additionalProperties。
- client layer schema 改动必须说明 vendor / codegen 是否需要更新。
- 基座消费版本是事实约束：不能把破坏性变更伪装成 patch。

### VERSION 文件 + tag 双真相与版本漂移 gate

本仓版本有两个真相源，必须始终一致：

1. **git tag**（`v0.1.0` / `v0.3.0` / `v0.3.0` …）：发布事实，由 Owner 授权后打。
2. **仓根 `VERSION` 文件**（内容如 `0.3.0`，不带 `v` 前缀）：显式版本真相源，供 CI、codegen、下游 vendor 读取，无需解析 git 历史即可知道当前版本。

**双真相必须一致**：`VERSION` 文件的值必须等于最新语义化 tag 去掉 `v` 前缀后的值（发版瞬间对齐）；一旦在最新 tag 之后修改了任何 `*.schema.json`，就必须先 bump `VERSION`，发布时再打对应的 `v*` tag。

**为什么需要制度化 gate**：契约仓最危险的回潮是"改了 schema 却没发版"——下游按旧版本号消费到了不兼容的形状，且没有任何信号。历史上正是"改 schema 不 bump"导致跨仓漂移。因此本仓用 `scripts/check-version-drift.sh` 作为 CI 强制 gate，把这条纪律从"靠自觉"变成"改了 schema 不 bump 就红"。

**SemVer bump 规则（判定改了 schema 后 `VERSION` 该怎么 bump）**：

| schema 变更类型 | bump 档位 | 示例（自 0.3.0 起） |
| --- | --- | --- |
| 新增可选字段、新增可选 enum 值、新增只读投影字段 | **minor** | `0.4.0` |
| 改必填（`required` 增项）、删字段、改字段类型、改 enum 语义、改 JSON tag、改默认语义 | **major** | `1.0.0` |
| 仅改注释 / description / 文档（不改结构与语义） | **不 bump**（patch 可选，无结构变化不触发 gate） | 保持 `0.3.0` |

**版本漂移 gate 判定逻辑**（`scripts/check-version-drift.sh`）：

- 取最新 tag `lastTag`（`git tag -l 'v*' --sort=-version:refname | head -1`），`lastVer` = 去 `v`。
- 取 `curVer` = `cat VERSION`。
- 若 `git diff --name-only <lastTag>..HEAD -- '*.schema.json'` 非空（自上次发版起改了 schema）**且** `curVer == lastVer`（没 bump）→ **FAIL(exit 1)**，打印哪些 schema 变了、要求 bump `VERSION`。
- 否则 **PASS(exit 0)**（没改 schema，或已 bump）。
- best-effort 附加：schema diff 中出现删属性行或 `required` 新增 → 打印 `WARN` 提示疑似破坏性变更需 major bump（仅提示，不额外阻断）。

本地自查：

```sh
bash scripts/check-version-drift.sh
```

CI 已在 `.github/workflows/ci.yml` 的 JSON 校验之后接入本 gate（`checkout` 用 `fetch-depth: 0` 取全部 tag，否则拿不到 `lastTag`）。发版流程：改 schema → 按上表 bump `VERSION` → gate 转绿 → Owner 授权后打对应 `v*` tag（tag 值 = `v` + `VERSION`）。

## 5. 允许的 helper 边界

本仓以契约为主，但允许少量 helper 来表达契约边界：

- 构造默认 candidate-only / non-formal 信封。
- 校验必需 ref、枚举硬地板、secret-ish 字段、Receipt / Gate 绑定条件。
- 生成稳定、确定性的 ref / hash 派生值。
- 判断 RegistrySlice 是否过期等纯数据方法。

helper 禁止：

- 访问 DB、网络、文件系统、环境变量、provider、sidecar、外部服务。
- 执行真实发送、真实写入、真实记忆、真实执行、真实上传。
- 持有全局状态、后台 goroutine、锁、缓存刷新或并发调度。
- 引入随机数、真实凭据、用户本地路径或运行态事实。
- 新增隐式 `time.Now()` 作为契约事实来源；需要时间时优先由调用方显式传入。历史兼容 helper 如需保留，应在注释中明确其只用于默认填充，不代表 Base 真实签发时间。

## 6. 工作纪律

- 修改任何文件前先运行 `git status --short --branch`。
- 新任务默认独立分支 + 独立 worktree；不要在主仓 `main` 直接开发。
- 不覆盖、不回滚、不删除他人 WIP。
- 不自动 merge、push、tag、release。
- 所有项目文档默认中文；英文只用于专有名词、命令、路径、代码标识、协议名、API 字段、错误原文。
- 不提交 `node_modules`、`dist`、`build`、`.vite`、日志、数据库、密钥、真实凭据、运行态文件或本地临时产物。
- 计划、阶段汇报、closeout、审计报告超过 500 字时写入单独 Markdown 文件，对话只给路径和短结论。
- 发现错误或测试失败必须定位根因；禁止吞异常、假成功、fallback 占位或绕过报错冒充完成。

## 7. 验证命令

基础验证：

```sh
go build ./... && go test ./... && go vet ./...
```

反向依赖：

```sh
go list -deps ./... | grep -E 'lights314/truzhenos|truzhen/packs' && echo "违规:反向依赖" || echo "OK:零反向依赖"
```

schema JSON 合法性：

```sh
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('*.schema.json') + glob.glob('spines/*.schema.json')];print('schema JSON 合法')"
```

schema embed 覆盖检查（改 schema / `embed.go` 时必跑）：

```sh
python3 - <<'PY'
import glob, pathlib, re
schemas=set(glob.glob('*.schema.json')+glob.glob('spines/*.schema.json'))
embeds=set(re.findall(r'//go:embed\s+([^\n]+)', pathlib.Path('embed.go').read_text()))
print('missing:', sorted(schemas-embeds))
print('extra:', sorted(embeds-schemas))
PY
```

说明：若某 schema 被裁定只作为文件契约、不通过 Go embed 暴露，必须在 `MODULES.md` 与 `embed.go` 附近写清楚，不得让“漏 embed”和“故意不 embed”混在一起。

## 8. 出错先看哪里

- build / import 报错：检查是否误引 `truzhenos`、`truzhen-packs` 或其它实现仓；检查 Go 标准库版本与 `go.mod`。
- 下游基座 build 失败：优先怀疑破坏性契约变更、JSON tag 漂移、必填字段变化、enum 语义变化。
- schema 校验失败：检查 JSON 格式、required、enum、additionalProperties、`embed.go` 路径。
- client codegen 失败：检查 schema draft、字段命名、required、nullable 表达、vendor 副本是否同步。
- 反向依赖检查失败：立即移除实现仓 import，改为 ref / interface / schema 形状。

## 9. 生命周期口径

本仓描述功能 / 契约状态时统一使用：

`想法 -> 设计中 -> 契约已定 -> 已实现 -> 已接线 -> 已验收 -> 已发布 -> 已弃用`

契约仓自身常用口径：

- `契约已定`：Go type / schema 已落本仓，但下游未必接线。
- `已接线`：至少一个下游仓真实消费并通过对应验证。
- `已验收`：下游消费、schema 校验、反向依赖、兼容测试均有证据。
- `已发布`：tag / module 版本已由 Owner 授权发布。

只写计划、mock、demo、readmodel-only 或 candidate-only 不能冒充 `已验收` / `已发布`。

## 10. 哪些变更必须回 Owner

- 任何破坏性契约变更：删字段、改必填、改类型、改 enum 语义、改 JSON tag、改主权链语义。
- 新增 / 删除子包、新增 / 删除 schema。
- 修改 Candidate/Formal 隔离、Gate/Receipt/Formalization、ReadModel 真相源边界。
- 引入任何非标准库依赖。
- 调整 module path、版本策略、开源许可、发布流程。
- 跨仓边界、依赖方向、五落点归属相关调整。
- `git push`、打 tag、发版、改远端设置。
