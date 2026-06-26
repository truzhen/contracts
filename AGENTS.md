# AGENTS.md — truzhen-contracts

本仓 = Truzhen **契约层** SDK（`github.com/truzhen/contracts`，开源 Apache-2.0）。智能体只读本文件即可独立维护本仓。

## 1. 本仓职责
定义基座与包之间**一切跨边界的数据形状**：候选信封、门控决议、回执 / 审计、注册切片、ReadModel 投影、监控事件、三主线引用，以及机器可校验的 `*.schema.json`。纯类型 / 接口 / 常量，无实现。

## 2. 禁止事项
- **禁止写实现**：不得引入 DB、网络、文件 I/O、time / rand 副作用、并发逻辑——那些属于基座 `truzhenos`。
- **禁止反向依赖**：不得 import `github.com/lights314/truzhenos` 或 `github.com/truzhen/packs`。
- **禁止真凭据**：`secrets/` 只放引用形状；任何 API key / token / 口令 / terminal_sn / 激活码绝不进本仓。
- **禁止无意识破坏边界**：删字段 / 改必填 / 改语义不升 major。

## 3. 必读文件
- `README.md`（依赖方向、SDK 用法、SemVer）
- `MODULES.md`（子包清单）
- `CLAUDE.md`（铁律速记）
- 改 scene-pack schema 前读对应 `*.schema.json` 与基座设计文档 `docs/design/scene-pack-vertical-profession-workbench-upgrade-20260626.md`（在基座仓 truzhenos）。

## 4. 常用验证命令
```sh
go build ./... && go test ./... && go vet ./...
go list -deps ./... | grep -E 'lights314/truzhenos|truzhen/packs' && echo "违规:反向依赖" || echo "OK:零反向依赖"
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('*.schema.json')];print('schema JSON 合法')"
```

## 5. 出错时先看哪里
- **本仓 build 失败 / import 报错** → 多半误引了基座或包的包：检查是否写成了反向依赖（铁律 2）。
- **下游基座 build 失败** → 你可能做了破坏性契约变更：核对是否删 / 改了字段、是否该升 major。
- **schema 校验失败** → 改坏了 `*.schema.json` 的 JSON 结构，或 `embed.go` 的嵌入路径不匹配。
- **`go list -deps` 出现 truzhenos / packs** → 反向依赖泄漏，立即回退到不依赖实现的写法。

## 6. 哪些变更必须回 Owner
- 任何**破坏性契约变更**（删字段 / 改必填 / 改语义 / 升 major）。
- 新增或删除子包、新增 / 删除 `*.schema.json`。
- `git push` / 打 tag / 发版等外向动作（Owner 在场授权）。
- 依赖方向、仓边界相关的任何调整。
