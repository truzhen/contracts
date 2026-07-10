<div align="center">

# 徒真（truzhen）Contracts

**信任不靠口头承诺，靠白纸黑字的公开契约。**

想创作 Pack？从 [truzhen/packs](https://github.com/truzhen/packs) 出发 · 本仓回答的是：这套生意为什么做得成

简体中文 · [繁體中文](docs/i18n.md#zh-tw) · [English](docs/i18n.md#en) · [日本語](docs/i18n.md#ja) · [Русский](docs/i18n.md#ru) · [Deutsch](docs/i18n.md#de) · [Français](docs/i18n.md#fr) · [Español](docs/i18n.md#es)

</div>

---

## 一个市场最难的问题

徒真（truzhen）让开源作者和行业专家把能力做成 Pack，卖给中国的中小企业。这里有个绕不开的问题：**企业凭什么敢安装一个陌生人做的东西？**

答案不是「相信作者人品」，而是一套公开、稳定、机器可校验的规则：Pack 能声明什么、必须经过谁的确认、做完之后留下什么凭证——全部写在这个仓库里，作者看得见，企业也看得见。

`github.com/truzhen/contracts` 就是这套规则的权威源。它只有 Go 类型和 JSON Schema，零实现、零副作用：谁都可以放心依赖它，它不依赖任何实现仓。

## 三句话说清这套规则

1. **AI 和 Pack 只能提建议**——所有产出先是「候选」，不直接变成正式结果。
2. **要紧的动作必须过门**——发送、执行、写正式数据，都要经过企业负责人确认和平台裁定。
3. **做过的事必须留回执**——企业随时可以回放核对发生了什么。

这三条对作者是保护而不是束缚：你的 Pack 不用碰权限、支付和审计，出了问题记录说话，责任分得清。

## 你什么时候需要看本仓

- **只想做 Pack**：先去 [truzhen/packs](https://github.com/truzhen/packs)。你只需要知道候选、门控和回执这三条底线，不必先读完整 SDK。
- **要写工具或校验器**：用本仓的 Go 类型和 JSON Schema 检查 Pack、候选、回执或市场表面数据。
- **要改跨仓字段**：先在本仓定清形状、版本和兼容策略，再让基座、Pack、云端或客户端消费。

## 这里有什么

| 契约 | 它保证的事 |
| --- | --- |
| Candidate（候选） | AI 的建议和草稿有统一形状，永远和正式结果隔离 |
| Gate（门控） | 高风险动作必须回到企业负责人确认和平台裁定 |
| Receipt（回执） | 重要动作留下可回放的凭证 |
| Registry / Provider 引用 | Pack 只声明需要什么外部能力，不夹带实现 |
| Delegation（委托） | Owner 预授权只能在明确边界内生效；代码执行委托必须额外声明 execution_scope |
| ReadModel / Surface | 界面展示有统一形状，但展示不等于事实 |
| Market 契约面 | 支付、授权、下载各归其位，云端状态无法在 Pack 里伪造 |

## 作为 SDK 使用

如果你只是提交 Pack，通常不需要写这里的 Go 代码；这一段主要给基座、CI、工具作者和高级 Pack 作者使用。

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

依赖方向是单向的：徒真基座**实现**这些契约，Pack **面向**这些契约声明，本仓不依赖任何实现仓（只依赖 Go 标准库，可用 `go list -deps ./...` 验证）。

当前版本以仓根 `VERSION` 和已发布 tag 为准。破坏性变更按 SemVer 处理，不把删字段、改必填、改语义伪装成兼容更新。

`market.PackSoftwareRequirement` 用于声明 Baserow / OCR 等底层软件需求；`market.SoftwareResolutionLock` 用于消费 `truzhenos` resolver 产出的复用、需安装、版本冲突、需隔离、blocked、not_ready 等结果。contracts 不解析用户本机环境，不保存本机软件事实。

## 子包速览

核心：`base/`（主权门控核心类型，含 `OwnerDelegationGrant` 与可选代码执行委托边界）、`candidates/`（候选域）、`gates/`（门控裁定）、`receipts/`（回执 / 审计）、`spines/`（事务 / 意图 / 证据三主线）、`registry/`、`readmodels/`、`monitoring/`、`secrets/`（只有 secret 的**引用**形状，永无明文凭据）、`market/`。完整清单见 [MODULES.md](MODULES.md)。

## 我们的承诺

- contracts 只定义形状，不给任何 Pack 执行权。
- 明文凭据永不进入本仓。
- 契约稳定可依赖：破坏性变更必须升版本、给迁移说明，不搞突然袭击。
- 本仓零反向依赖，你引用它不会被拖进私有实现。

## 验证

```sh
go build ./... && go test ./... && go vet ./...
go list -deps ./... | grep -E 'lights314/truzhenos|truzhen/packs|truzhen/truzhen-cloud' && echo "违规:反向依赖" || echo "OK:零反向依赖"
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('*.schema.json') + glob.glob('spines/*.schema.json')];print('schema JSON 合法')"
```

## 想改契约？

先想清四个问题：给谁用、是不是只声明形状、影响哪些消费方、算不算破坏性变更。然后读：

- 贡献指南：[CONTRIBUTING.md](CONTRIBUTING.md)
- 契约治理总纲：[CONTRACTS_GOVERNANCE.md](CONTRACTS_GOVERNANCE.md)
- Agent 开工纪律：[AGENTS.md](AGENTS.md)

## License

[Apache-2.0](LICENSE)。徒真（truzhen）是对外品牌与商标。
