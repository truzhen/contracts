# Truzhen Contracts

> Truzhen 主权事务操作层的**契约 SDK**：Go 类型、接口、常量、JSON Schema、schema embed，以及少量无外部副作用的契约校验 / ref 派生 helper。

`github.com/truzhen/contracts` 是 Truzhen 六仓协同架构的**契约层**，定义基座、包层、cloud repo 和 client layer 之间一切跨边界的数据形状（Pack / Candidate / Receipt / Surface / ReadModel schema、cloud 契约面（Entitlement / License / Payment / PackListing / Session / Release / WebSurface，当前为治理清单、具体 schema 待真实消费方）、候选信封、门控决议、回执、注册切片、监控事件、三主线引用等）以及机器可校验的 JSON Schema。它**只声明形状，不含任何实现**（无 DB、无网络、无副作用），因此谁都可以安全依赖它，而它谁都不依赖。

## 依赖方向（单向不可逆）

```text
truzhenos (私有基座)  implements  truzhen-contracts  faces  truzhen-packs
实现契约 / 持主权                  契约形状权威源              面向契约声明包
```

- **基座**（`github.com/lights314/truzhenos`，私有）**实现**这些契约。
- **包**（`github.com/truzhen/packs`，开放）**面向**这些契约编写，物理上 import 不到基座内部。
- **云仓**（`github.com/truzhen/truzhen-cloud`，私有）**实现**官方云端服务、WebSurface、支付、License / Entitlement、Pack listing / 分发和 Release / Session 相关契约。
- **客户端**（`github.com/truzhen/truzhen-client-web-desktop`，私有）通过 vendor / codegen 消费 schema，不私造稳定 DTO。
- **契约本身零反向依赖**：只依赖 Go 标准库，不依赖基座、包、云仓、client 或 provider 实现。可用 `go list -deps ./...` 验证。

## 作为 SDK 使用

```go
import (
    "github.com/truzhen/contracts/base"
    "github.com/truzhen/contracts/candidates"
    "github.com/truzhen/contracts/receipts"
)
```

```sh
go get github.com/truzhen/contracts@latest
```

当前基座消费版本：`github.com/truzhen/contracts v0.3.0`。破坏性变更不得落到 `v0.x` patch 假装兼容；跨仓边界变化必须同步更新 `truzhenos`、`truzhen-packs`、`truzhen-cloud`、`truzhen-software` 和 client repo 的治理说明。

## 子包总览

完整清单见 [MODULES.md](MODULES.md)。核心：`base/`（主权门控核心类型）、`candidates/`（AI 候选域）、`gates/`（门控裁定）、`receipts/`（回执/审计）、`spines/`（事务/意图/证据三主线）、`registry/`、`readmodels/`、`monitoring/`、`secrets/`（secret **引用**契约，不含真凭据）、`events/`、`modules/`。cloud 契约组按 `Entitlement`、`License`、`Payment`、`PackListing`、`Session`、`Release`、`WebSurface` 七类收敛，落地为具体 schema / Go 包前不得宣称下游已接线。

- `base/`：Base 主权核心契约、Gate / ReceiptCandidate / FormalizationGrant、授权、委托、Artifact 留痕与过闸边界。
- `candidates/`：AI / Pack / 模块候选域类型。
- `spines/`：Transaction / Intent / Evidence 三主线引用与 Intent Spine 五件套。
- `receipts/`：回执与审计信封。
- `gates/`：轻量 AccessDecision / OwnerVerdict。
- `registry/`：注册引用和 RegistrySlice。
- `readmodels/`：前端只读投影信封。
- `monitoring/`：统一监控、诊断、故障和支持包契约。
- `secrets/`：secret 引用形状，不含明文凭据。
- `events/`、`modules/`：模块事件和生命周期接口。
- 顶层：schema embed 与 Pack 知识挂载契约。

## 设计原则

1. **形状权威源**：跨仓 DTO、schema、ref、信封、枚举在本仓收敛。
2. **零反向依赖**：本仓不得依赖任何实现仓。
3. **Candidate/Formal 隔离**：候选默认 `candidate_only=true`、`non_formal=true`；正式化必须在基座完成。
4. **ReadModel 不是真相源**：前端投影只能展示，不能决定事实。
5. **机器可校验**：JSON Schema 与 `embed.go` 为 Go 服务、client codegen 和 CI 提供统一契约入口。
6. **helper 克制**：只允许无外部副作用的契约校验和确定性 ref 派生；不写运行时实现。

## 治理入口

- Agent 开工纪律：[AGENTS.md](AGENTS.md)
- 契约治理总纲：[CONTRACTS_GOVERNANCE.md](CONTRACTS_GOVERNANCE.md)
- 子包和 schema 清单：[MODULES.md](MODULES.md)
- 快速加载说明：[CLAUDE.md](CLAUDE.md)
- 贡献指南：[CONTRIBUTING.md](CONTRIBUTING.md)

## 验证

```sh
go build ./... && go test ./... && go vet ./...
go list -deps ./... | grep -E 'lights314/truzhenos|truzhen/packs' && echo "违规:反向依赖" || echo "OK:零反向依赖"
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('*.schema.json') + glob.glob('spines/*.schema.json')];print('schema JSON 合法')"
```

## License

[Apache-2.0](LICENSE)。
