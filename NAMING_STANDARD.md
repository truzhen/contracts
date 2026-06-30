# truzhen-contracts 命名标准

## 通用规则

- 包名使用小写单词，跨边界共享契约按业务域归类。
- Go 类型使用稳定名词，不使用实现后缀，例如 `CloudEntitlement`、`PaymentWebhook`。
- JSON 字段统一使用 `snake_case`。
- 引用字段使用 `_ref` 后缀，例如 `entitlement_ref`、`receipt_ref`。
- digest 字段必须写明算法或使用带算法前缀的字符串。

## Cloud 契约

- 官方云服务共享 DTO 统一放在 `cloud/` 包。
- 云端状态类常量使用具体领域前缀，例如 `EntitlementStatusActive`、`CloudWebPublishReleaseCandidate`。
- Cloud 契约只能表达形状，不能包含 provider 实现、DB 表结构、部署路径、secret 值或生产服务器地址。
- `truzhen-cloud`、`truzhenos` 和 client layer 需要共享的 cloud 形状先进入本仓，再由实现仓消费。
