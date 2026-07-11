# truzhen-contracts 交付清理清单

更新时间：2026-07-11（本文件由 2026-07-11 交付债务集中处理轮创建——此前本仓无交付账本，登记缺口本身即债务）

> **账本制度**：本文件是 contracts 仓**分账**；六仓交付债务**总账** = truzhenos 仓 `DELIVERY_CLEANUP_CHECKLIST.md`。集中处理计划：truzhenos `docs/plans/six-repo-delivery-cleanup-debt-consolidation-plan-20260711.md`。
> 凭据值绝不入本文件；处理完把 `[ ]` 改 `[x]` + 日期 + 谁。

## 交付前遗留

- [ ] 🟡 **3 个 scene-studio 系 schema 未入 `embed.go`**（2026-07-06 审计发现，2026-07-11 核验仍成立）：`scene-studio-node-info.schema.json`、`scene-studio-workflow.schema.json`、`scene-runtime-plan-candidate.schema.json` 在仓根但无对应 embed 变量。**当前无 Go 侧消费方**——按「无消费方不投机施工」先登记；待真实消费方出现（Studio/引擎侧校验需要）时补 embed 变量 + `go build/test` + 覆盖关系检查（AGENTS §0.6）。
- [ ] 🟡 **schema 目录组织**：仓根 `*.schema.json` 与子目录（base/candidates/gates/market/…）并存，交付文档应说明布局约定，防新 schema 落错层。
- [x] 🟡 **README 陈旧包列表**（审计原条目「仍列已删的 events/、modules/ 包」）：✅ 2026-07-11 核验已被后续 README 重写自然修复（全文零 events/modules 包引用），无需再改。
- [x] 🟠 **AGENTS.md 自述消费版本漂移**（审计 R-1：自述 v0.3.0 落后实况）：✅ 2026-07-11 修正为「以 `VERSION` / 下游 `go.mod` 为准」，不再手抄数字。

## 跨仓关联（登记在总账，此处只指针）

- truzhenos require 版本对齐（总账「六仓治理审计遗留」contracts↔truzhenos 条）：contracts 侧 `VERSION=0.11.0` + tag `v0.11.0` 已发布无缺口；对齐动作在 truzhenos 侧 go.mod。
