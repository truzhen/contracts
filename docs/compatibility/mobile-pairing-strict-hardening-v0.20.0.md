# v0.20.0 移动首配对严格解码兼容说明

状态：`契约已定（开发线）`，未发布。仓根 `VERSION` 为 `0.20.0`；本次没有 tag、push、module 发布或下游运行时部署。

## 变更范围与 SemVer

这是从 `v0.19.0` 到 `v0.20.0` 的加性 Go API 变更。既有 JSON Schema、字段、required 列表、枚举和 `embed.go` 暴露均未修改：

- `mobile-pairing-bootstrap-request.schema.json`
- `mobile-pairing-bootstrap-candidate.schema.json`
- `mobile-session-issue-intent.schema.json`

新增 `readmodels` 的移动平台 / 候选状态常量，以及这三种 shape 的严格 `Decode` 与 `Validate` helper。helper 仅做闭合 JSON 和 schema 级形状校验，不访问网络、文件、数据库、Provider、Base Gate 或 Receipt Ledger。

## 主权边界

严格 decoder 拒绝 unknown field，因此下列字段不能经移动 JSON body 走私：

- `owner_decision_ref`、任何调用方自铸的批准事实；
- `receipt_ref`、正式回执或伪造正式引用；
- `bootstrap_proof`、bearer、token、raw credential；
- 多顶层 JSON 值，以及候选中被省略的 required `false` 字段。

helper 不会创建配对、会话、OwnerDecision、FormalRecord 或 Receipt。正式动作仍归 `truzhenos` 的 Host / Base Gate / Receipt 链路；bootstrap proof 仍只能由受控 header 传递。

## 消费影响清单（manifest）

| 消费方 | 影响 | 本轮要求 | 不做的事 |
| --- | --- | --- | --- |
| `truzhen-contracts`（真相源） | 新增可选纯 Go 防御 API | 保持三份 schema canonical bytes 不变；用 Go 单测覆盖未知字段、authority smuggling、缺失 required false 和 trailing JSON。 | 不实现运行时、登录、Gate、Receipt 或存储。 |
| `truzhenos`（Host / 16-auth） | 后续可在 HTTP JSON 解析边界显式调用 helper | 采用前须做独立接线与完整 Host 生命周期验收；继续把 proof 限在 header、正式化限在 Base Gate。 | 不得把 helper 当作审批、会话签发或 Receipt 的替代。 |
| `truzhen-client-web-desktop`（PC / mobile consumer） | schema 未变，现有 vendor/codegen 不需要更新 | 如消费 Go 侧返回值或增加 API 防御测试，应以本仓常量和 schema 语义对齐。 | 不得在客户端伪造 OwnerDecision、Receipt 或 bootstrap proof。 |
| `truzhen-cloud` | 无 schema / API 消费变化 | 无需因本次变更修改云端。 | 不保存或转发移动 proof / credential。 |
| `truzhen-packs` / `truzhen-software` | 无直接依赖 | 无动作。 | 不引入移动首配对事实源。 |

## 下游验收与发布门

下游只有在独立接线、对应仓验证和 Owner 授权后才能把运行时行为标为“已接线”或“已验收”。发布 `v0.20.0` tag、更新下游 Go module 版本、push 或部署均不在本次范围内。
