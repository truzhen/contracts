# 徒真（truzhen）Contracts

> 让社区作者的 Pack 能被企业安全安装的公开契约。

徒真（truzhen）希望让开源作者、工具开发者和行业专家，把自己的软件、工具链和行业经验做成企业可以安装、试用、购买和回放的 Pack。企业愿意安装陌生作者的能力，前提是边界清楚：Pack 能声明什么、能请求什么、哪些动作必须确认、执行后留下什么回执。

`github.com/truzhen/contracts` 就是这套公开边界。它提供 Go 类型、JSON Schema、候选信封、门控和回执形状，让 Pack 作者面向稳定规则创作，而不是依赖徒真私有基座内部实现。

如果你想创作能力包、角色包或场景包，请从 [`github.com/truzhen/packs`](https://github.com/truzhen/packs) 开始。本仓回答的是：你的 Pack 为什么不能越权、平台如何识别候选、企业如何知道动作留下了什么证据。

## 为什么 Pack 作者需要 contracts

contracts 让作者和企业对同一件事有共同理解：

- Pack 可以声明能力、角色、场景、候选、证据和回执要求。
- Pack 不能直接绕过主人确认。
- Pack 不能直接写正式数据。
- Pack 不能直接读取 secret、token、cookie、private key 或用户凭据。
- Pack 不能跳过 Gateway 和 Receipt 自行完成真实发送或真实执行。
- ReadModel 和 Surface 只用于展示，不是真相源。

这意味着作者可以把能力开放给企业使用，同时不用把自己变成企业权限系统、支付系统、审计系统或桌面客户端的维护者。

## 这里有什么

| 内容 | 给作者的意义 |
| --- | --- |
| Candidate / 候选形状 | AI、角色和 Pack 先提出建议或草稿，不直接变成正式结果。 |
| Gate / 门控形状 | 高风险动作必须回到主人确认和平台裁定。 |
| Receipt / 回执形状 | 重要动作要留下可回放证据，企业可以复核发生过什么。 |
| Registry / Provider 引用 | Pack 只声明需要什么外部能力，不把 provider 实现塞进 Pack。 |
| ReadModel / Surface schema | 前端展示有统一形状，但展示不等于事实。 |
| Market 契约面 | 支付、授权、下载和商品信息有边界；真实云端状态不在 Pack 里伪造。 |

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

当前公开消费版本以仓根 `VERSION` 和已发布 tag 为准。破坏性变更必须按 SemVer 处理，不能把删除字段、改必填、改语义伪装成兼容更新。

## 对作者的安全承诺

- contracts 只定义可声明、可请求、可校验的形状，不给 Pack 执行权。
- contracts 不实现 Base Gate、Receipt Ledger、Gateway、provider、云端支付或前端运行时。
- 候选成为正式任务、正式记忆、真实发送或真实执行前，必须回到徒真受控链路完成确认、调用和回执。
- Secret 只在本仓表现为引用形状；明文凭据永不进入 contracts。
- 本仓只依赖 Go 标准库，不反向依赖基座、packs、cloud、client 或 provider 实现。

## 贡献前自查

1. 这个字段或 schema 是给真实消费方用的吗？
2. 它只是声明形状，还是把运行实现搬进 contracts 了？
3. 是否影响 Pack 作者、平台基座、云端市场或客户端展示？
4. 是否是破坏性契约变更，是否需要版本升级和迁移说明？

## 验证

```sh
go build ./... && go test ./... && go vet ./...
go list -deps ./... | grep -E 'lights314/truzhenos|truzhen/packs|truzhen/truzhen-cloud' && echo "违规:反向依赖" || echo "OK:零反向依赖"
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('*.schema.json') + glob.glob('spines/*.schema.json')];print('schema JSON 合法')"
```

## 治理入口

- Agent 开工纪律：[AGENTS.md](AGENTS.md)
- 契约治理总纲：[CONTRACTS_GOVERNANCE.md](CONTRACTS_GOVERNANCE.md)
- 子包和 schema 清单：[MODULES.md](MODULES.md)
- 贡献指南：[CONTRIBUTING.md](CONTRIBUTING.md)

## License

[Apache-2.0](LICENSE)。
