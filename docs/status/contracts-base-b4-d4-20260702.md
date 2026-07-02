# contracts base B4/D4 进展（2026-07-02）

## 派活卡

- 我要做的事：完成 WP-D4 `base/contracts.go` 主题拆分，并完成 WP-B4 `base` 首批 helper 测试。
- 真相源：`truzhen-contracts` 当前 Git 工作树、Go 编译 / 测试输出、schema / 版本漂移检查、`truzhenos` 交叉 build。
- 仓库 / 层归属：只修改 `truzhen-contracts`；交叉验证只读 `truzhenos` 隔离 worktree。
- 风险颜色：黄色；不改 schema、不改 VERSION、不改导出签名、不发布 tag。
- 契约影响：零导出面变化；`github.com/truzhen/contracts/base` import 路径不变。
- 禁止边界：不处理 D5 `events/` / `modules/`，该项仍需 Owner 裁定；不改下游仓 `go.mod` / `go.sum`。

## 本轮完成

| 项 | 状态 | 说明 |
| --- | --- | --- |
| D4 拆分 | 已完成 | 删除 `base/contracts.go`，按主题拆为 `base/policy.go`、`base/gate.go`、`base/owner.go`、`base/formalization.go` |
| B4 测试 | 已完成 | 新增 `base/helpers_test.go`，覆盖 policy snapshot 稳定 ref、Gate intake / decision / receipt / grant 链、secret-ish adapter 拦截、object truth source 约束、`stableRef` 确定性 |
| schema / version | 已验证 | 未修改任何 `*.schema.json`，`VERSION` 仍为 `0.3.0` |
| 下游交叉 build | 已验证 | 用临时 `go.work` 让 `truzhenos` 消费当前 contracts 工作树，`go build ./backend/...` 通过；未写入 replace |

## 验证证据

```sh
go test ./base -count=1
bash scripts/check-version-drift.sh
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('*.schema.json') + glob.glob('spines/*.schema.json')];print('schema JSON 合法')"
go list -deps ./... | grep -E 'lights314/truzhenos|truzhen/packs|truzhen/truzhen-cloud' && echo '违规:反向依赖' || echo 'OK:零反向依赖'
go build ./... && go test ./... && go vet ./...
git diff --check
GOWORK=<临时 go.work> go -C /Users/li/Documents/truzhenv3worktree/truzhenos-wp-a-doc-honesty-20260702 build ./backend/...
```

全部通过。

## 剩余事项

- WP-D5 `events/` / `modules/` 仍等待 Owner 裁定：删除还是标注占位。
- 本分支未 commit、未 push、未 tag、未发布。
