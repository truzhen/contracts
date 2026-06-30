# 贡献指南 — truzhen-contracts

本仓是 Truzhen 的开源契约 SDK。欢迎贡献，但必须遵守契约层边界：本仓定义跨仓数据形状，不实现基座、Pack、前端或 provider。

## 你可以做什么

- 修正注释、文档、schema description。
- 在不破坏兼容的前提下新增可选字段。
- 补充 Go 类型、枚举、接口、JSON Schema。
- 补充无外部副作用的契约校验、默认信封构造、确定性 ref 派生 helper。
- 完善 `MODULES.md`、`AGENTS.md`、`CONTRACTS_GOVERNANCE.md` 中的模块和治理说明。

## 你不能做什么

- 引入 DB、网络、文件 I/O、provider、sidecar、真实执行、后台并发或运行态状态。
- 让本仓 import `github.com/lights314/truzhenos`、`github.com/truzhen/packs`、client repo 或 provider repo。
- 把 Pack 数据、安装脚本、React/Tauri 前端源码、外部软件 registry 放进本仓。
- 提交任何真凭据：API key、token、cookie、private key、terminal_sn、激活码、密码。
- 删除字段、改必填、改类型、改 JSON tag、改 enum 语义后仍按兼容变更处理。
- 用轻量 `gates.AccessDecision` / `OwnerVerdict` 冒充正式 `base.GateDecision` / `base.OwnerDecision`。

## 变更前自查

1. 这个事实归谁？本仓是不是只声明形状？
2. 下游谁消费：`truzhenos`、`truzhen-packs`、client repo、CI、codegen？
3. 是否影响 SemVer？
4. 是否需要同步 Go struct、JSON Schema、`embed.go`、`MODULES.md`、README、client vendor / codegen？
5. 是否需要 Owner 裁定？

## 提交前自检

```sh
go build ./... && go test ./... && go vet ./...
go list -deps ./... | grep -E 'lights314/truzhenos|truzhen/packs' && echo "违规" || echo "OK"
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('*.schema.json') + glob.glob('spines/*.schema.json')];print('schema JSON 合法')"
```

改 schema 或 `embed.go` 时，必须再跑 `AGENTS.md` 中的 embed 覆盖检查，并解释 missing / extra 是否符合预期。

## 版本

遵循 SemVer：

- patch：非语义修正。
- minor：新增可选字段或兼容 schema。
- major：破坏性变更。

破坏性变更必须先回 Owner，并在 PR / closeout 中说明影响面、迁移策略、下游同步顺序和回滚方案。

## License

贡献即同意以 [Apache-2.0](LICENSE) 授权你的贡献。
