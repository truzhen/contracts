# 贡献指南 — truzhen-contracts

本仓是 Truzhen 的开源契约 SDK。欢迎贡献，但请遵守契约层的边界。

## 你可以做什么
- 修正类型注释、文档、schema 描述。
- 在**不破坏现有形状**的前提下，新增可选字段（升 minor）。
- 补充 `MODULES.md` 子包说明。

## 你不能做什么
- 引入任何实现（DB / 网络 / 文件 / 副作用）——那属于基座 `truzhenos`。
- 让契约 import 基座（`github.com/lights314/truzhenos`）或包（`github.com/truzhen/packs`）——零反向依赖是本仓的命根。
- 提交任何真凭据（API key / token / 口令 / terminal_sn / 激活码）。`secrets/` 只放引用形状。
- 无意识的破坏性变更（删字段 / 改必填 / 改语义）而不升 major。

## 提交前自检
```sh
go build ./... && go test ./... && go vet ./...
go list -deps ./... | grep -E 'lights314/truzhenos|truzhen/packs' && echo "违规" || echo "OK"
```

## 版本
遵循 SemVer。破坏性变更必须升 major，并在 PR 说明影响面。

## License
贡献即同意以 [Apache-2.0](LICENSE) 授权你的贡献。
