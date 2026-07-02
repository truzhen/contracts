# CLAUDE.md — truzhen-contracts

本文件供 Claude Code / Agent 在 `truzhen-contracts` 根目录工作时快速加载。完整纪律以 `AGENTS.md` 和 `CONTRACTS_GOVERNANCE.md` 为准；本文件只做提炼与导航，冲突时以治理原文为准。

> 文档纪律：中文为主；英文只用于专有名词、命令、文件名、路径、协议名、API 字段、代码标识和错误原文。

## 1. 仓库定位

`truzhen-contracts` 是 Truzhen 六仓协同架构的**开源契约层 SDK**，Go module 为 `github.com/truzhen/contracts`。它定义基座、包层、cloud repo、client layer 之间共享的数据形状和机器 schema：Pack / Candidate / Receipt / Surface / ReadModel schema、候选信封、门控裁定、回执 / 审计、注册切片、监控事件、三主线引用等。

依赖方向只允许单向：

- `truzhenos` 实现本仓契约。
- `truzhen-packs` 面向本仓契约声明 Pack。
- client repo vendor / codegen 消费本仓 schema。
- cloud repo 实现本仓 cloud 契约面（治理清单先行，schema 待真实消费）。
- 本仓不得 import 基座、packs、cloud、client 或 provider 仓。

## 2. 铁律

1. **契约层只表达边界**：允许 type / const / interface / JSON Schema / schema embed / 无外部副作用的确定性校验与 ref 派生 helper；禁止 DB、网络、文件 I/O、provider、真实执行、后台并发、运行态状态。
2. **零反向依赖**：`go list -deps ./...` 必须无 `github.com/lights314/truzhenos`、`github.com/truzhen/packs` 和 `github.com/truzhen/truzhen-cloud`。
3. **改契约 = 改跨仓边界**：删字段、改必填、改类型、改 enum 语义、改 JSON tag 是破坏性变更，必须先回 Owner 并按 SemVer 处理。
4. **Candidate/Formal 隔离不可破**：候选类型默认 `candidate_only=true`、`non_formal=true`；正式写入、发送、执行、记忆、回执账本实现都不属于本仓。
5. **ReadModel 不是真相源**：本仓可以定义投影形状，但不能把投影视为权威状态。
6. **secret 只放引用形状**：`secrets/` 只能定义 `SecretRef` / `SensitivePayload` 等引用，不得出现明文凭据。
7. **版本漂移即红**：最新 tag 后改任何 `*.schema.json` 必先 bump 仓根 `VERSION`，CI `check-version-drift.sh` 强制。

## 3. 模块导航

- `base/`：Base 主权核心契约、Gate / ReceiptCandidate / FormalizationGrant、授权、委托、Artifact 留痕与过闸边界。
- `candidates/`：Advice、CommunicationDraft、ExecutionIntent、CapabilityInvocation、MemoryWrite、Task 等候选域类型。
- `spines/`：Transaction / Intent / Evidence 三主线引用与 Intent Spine 五件套。
- `receipts/`：ReceiptEnvelope、AuditEnvelope。
- `gates/`：轻量 AccessDecision / OwnerVerdict，不能替代 `base.GateDecision`。
- `registry/`：RegistryRef、SkillRef、RegistrySlice。
- `readmodels/`：ReadModelEnvelope。
- `monitoring/`：监控、诊断、故障、支持包上传候选等契约。
- `secrets/`：secret 引用形状。
- `market/`：市场表面契约（会话头 / Login DTO / 表面路径 / admin 转发硬 allowlist；服务端真相唯一在 truzhen-cloud）。
- 顶层：`embed.go` 暴露 schema bytes；`pack_knowledge_mount.go` 定义 Pack 知识挂载契约。

完整清单见 `MODULES.md`。

## 4. 常用验证

```sh
go build ./... && go test ./... && go vet ./...
go list -deps ./... | grep -E 'lights314/truzhenos|truzhen/packs|truzhen/truzhen-cloud' && echo "反向依赖!违规" || echo "零反向依赖 OK"
python3 -c "import json,glob;fs=glob.glob('**/*.schema.json',recursive=True);assert fs;[json.load(open(f)) for f in fs];print('schema JSON 合法 x%d' % len(fs))"
```

改 schema 或 `embed.go` 时，再跑 `AGENTS.md` 中的 embed 覆盖检查。

## 5. 开工提醒

- 先运行 `git status --short --branch`。
- 新任务默认独立分支 + 独立 worktree。
- 修改治理、schema 或契约字段前，先判断真相源、消费方、兼容性和 SemVer。
- 不自动 push / tag / release / merge。
