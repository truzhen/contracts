# CLAUDE.md — truzhen-contracts

本仓是 Truzhen 五落点架构的**契约层**（`github.com/truzhen/contracts`），开源、Apache-2.0。供 Claude Code / Agent 在本仓工作时加载。**完整纪律见 `AGENTS.md`。**

> 文档纪律：中文为主，英文只用于专有名词、命令、文件名、路径、代码标识。

## 这是什么

纯接口 / 类型 / JSON Schema 的 SDK，定义基座与包之间一切跨边界的数据形状。**零反向依赖**：只依赖 Go 标准库；基座实现它、包面向它，它谁都不依赖。

## 铁律（改任何文件前先记住）

1. **只声明形状，不写实现**：本仓只有 `type` / `const` / 接口与 `*.schema.json`。任何带 DB、网络、文件、并发、时间副作用的代码**不属于这里**——属于基座 `truzhenos`。
2. **禁止反向依赖**：契约**不得 import** `github.com/lights314/truzhenos`（基座）或 `github.com/truzhen/packs`（包）的任何内容。`go list -deps ./...` 必须无这两者。
3. **改契约 = 改跨仓边界**：删字段 / 改语义 / 改必填 → 破坏性，必须升 major（SemVer）；新增可选字段升 minor。改动前想清楚谁在依赖。
4. **secret 只放引用形状**：`secrets/` 只定义 `SecretRef` / `SensitivePayload` 等「对密文的引用」类型，**真凭据值绝不进本开源仓**。

## 验证命令

```sh
go build ./...
go test ./...
go vet ./...
go list -deps ./... | grep -E 'lights314/truzhenos|truzhen/packs' && echo "反向依赖!违规" || echo "零反向依赖 OK"
```

## 子包导航

见 `MODULES.md`。
