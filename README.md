# Truzhen Contracts

> Truzhen 主权事务操作层的**契约 SDK**：纯接口、类型与 JSON Schema，零反向依赖。

`github.com/truzhen/contracts` 是 Truzhen 五落点架构的**契约层**，定义基座、包层和 client layer 之间一切跨边界的数据形状（Pack / Candidate / Receipt / Surface / ReadModel schema、候选信封、门控决议、回执、注册切片、监控事件、三主线引用等）以及机器可校验的 JSON Schema。它**只声明形状，不含任何实现**（无 DB、无网络、无副作用），因此谁都可以安全依赖它，而它谁都不依赖。

## 依赖方向（单向不可逆）

```
┌─────────────────────┐   implements   ┌──────────────────────┐   faces   ┌────────────────────┐
│ truzhenos (基座·私有) │ ─────────────▶ │ truzhen-contracts    │ ◀──────── │ truzhen-packs      │
│ 实现契约             │                │ 纯接口/类型/Schema    │           │ 面向契约           │
└─────────────────────┘                └──────────────────────┘           └────────────────────┘
                          契约零反向依赖：不 import 基座、不 import 包
```

- **基座**（`github.com/lights314/truzhenos`，私有）**实现**这些契约。
- **包**（`github.com/truzhen/packs`，开放）**面向**这些契约编写，物理上 import 不到基座内部。
- **客户端**（`github.com/truzhen/truzhen-client-web-desktop`，私有）通过 vendor / codegen 消费 schema，不私造稳定 DTO。
- **契约本身零反向依赖**：只依赖 Go 标准库，不依赖基座或包。可用 `go list -deps ./...` 验证。

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

当前基座消费版本：`github.com/truzhen/contracts v0.2.0`。破坏性变更不得落到 `v0.x` patch 假装兼容；跨仓边界变化必须同步更新 `truzhenos`、`truzhen-packs` 和 client repo 的治理说明。

## 子包总览

完整清单见 [MODULES.md](MODULES.md)。核心：`base/`（主权门控核心类型）、`candidates/`（AI 候选域）、`gates/`（门控裁定）、`receipts/`（回执/审计）、`spines/`（事务/意图/证据三主线）、`registry/`、`readmodels/`、`monitoring/`、`secrets/`（secret **引用**契约，不含真凭据）、`cloud/`（官方云服务共享 DTO）、`events/`、`modules/`。

## Cloud 契约

`cloud/` 只放官方云服务共享形状，供 `truzhen-cloud`、`truzhenos` 本地 Cloud proxy / License Gate 消费端和 client layer 共同对齐。当前包含：

- `CloudEntitlement`
- `LicenseToken` / `LocalActivationToken` / `LicenseValidationResult`
- `PaymentOrder` / `PaymentWebhook`
- `CloudPackListing`
- `CloudSession`
- `CloudReleaseCandidate` / `CloudReleaseReceipt`
- `CloudWebSurface` / `CloudWebRoute`

这些类型不包含 DB、HTTP client、支付 provider、部署脚本、secret 值或本地主权链实现。

## 版本策略

遵循 [SemVer](https://semver.org/lang/zh-CN/)。契约是跨仓边界：**破坏性变更（删字段、改语义、改必填）必须升 major**；新增可选字段升 minor。基座与包通过各自 `go.mod` 的 `require` 钉版本协同演进。

## 设计原则

1. **纯形状**：只有 `type`、`const`、接口与 JSON Schema；任何带 DB / 网络 / 文件 / 时间副作用的代码不属于本仓。
2. **零反向依赖**：CI 与 `go list -deps` 双重把守，契约不得 import 任何实现仓。
3. **机器可校验**：`*.schema.json` 经 `embed.go` 嵌入，基座与 CI 用它校验产物，而非在代码里重复声明形状。

## License

[Apache-2.0](LICENSE)。
