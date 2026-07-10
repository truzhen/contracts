# KWeaver 借鉴「现在补」两项落地计划（待 Owner 批「开工」）

> 日期：2026-07-10　|　作者：Claude Fable 5　|　生命周期：**设计中**（本文件写完即停，等 Owner 逐项裁定 + 明确「开工」）
> 两项 = 统一决策表 #2（Impact Model 建议稿）+ #3（场景荚 manifest `lifecycle_status`）。

## 0. 计划纪律必答项

| 栏 | 回答 |
| --- | --- |
| 版本/优先级 | 当前主线的**支线小件**（合计施工 ≤1.5 天）。节奏治理阶段 1 唯一主线=智能家居跑穿；本两项经 Owner 2026-07-10 明示「现在就落地」纳入，**施工时机建议=不挤占主线的空档** |
| 真实客户/场景证据 | **诚实标注：两项均缺直接客户原话证据**。A 的证据=三个独立来源（Palantir osdk `modifiedEntities`/ActionLogicRule、KWeaver ImpactContractItem、BKN impact_contracts）收敛到同一形状=行业共识 + Palantir 侧 01 决策表已裁「下张任务卡」；B 的证据=自家八档生命周期纪律已成文而 manifest 无字段（治理自洽修补，非新功能） |
| 最小可交付 | A=一份建议稿+影响清单文件（**零 schema/代码改动**）；B=contracts additive 字段+VERSION minor+packs 4+1 manifest 补值+门禁全绿 |
| 真相源 | 形状=truzhen-contracts；manifest 内容=truzhen-packs 各 pack 目录；A 涉及的运行事实（actual_edits）=03 回执账本（本轮不碰，只画形状） |
| 仓库/层归属 | A：contracts `docs/plans/`（建议稿）；B：contracts（schema+Go struct）→ packs（manifest 值），依赖方向=先契约后消费方 |
| 风险颜色 | A=**橙**（契约设计，只出稿不实施，天然合规）；B=**黄**（additive optional 字段，走既有防漂移门禁验证；required 数组零变化） |
| 契约影响 | A=零；B=minor（`VERSION` 0.8.0→0.9.0），`check-breaking-change.py` 对「新增可选属性」应 exit 0——若门禁误判则**修门禁根因，不绕过、不降级** |
| 上下文允许 | contracts 全仓、packs 各 manifest.json、考古笔记目录；只读参考 truzhenos manifest 消费点（如需） |
| 禁止边界 | 不动任何 schema `required` 数组；不删/改现有字段；不实施 A 的任何 runtime；KWeaver 代码零进仓；不碰 os 冻结模块；不 push/tag（归 Owner） |
| 验收设计 | 见 §3；独立验收=另派子代理复跑门禁 + 突变自证（详见任务步骤），执行方自报不算数 |
| 变更影响 | A：无运行影响；B：contracts 消费方（truzhenos pack 装载、cloud market、client vendor）对未知字段的容忍性——`PackManifest` 新字段带 `omitempty`，JSON 反序列化天然向后兼容；schema `additionalProperties:false` 意味着**必须先改 schema 再给 manifest 加值**（顺序锁死） |
| 生命周期档位 | A：设计中 →（Owner 批准建议稿后）契约已定；B：想法 → 已实现 → 已接线（校验读到）→ 已验收（独立复核） |

## 1. 待 Owner 裁定项（开工前逐项）

- **R-1 `lifecycle_status` 落点**：a) `pack-manifest.schema.json` + `market.PackManifest`（通用，三类 Pack 都受益）——**建议**；b) 仅 `scene-pack-spec.schema.json`（场景荚专属）。
- **R-2 枚举 token**：八档建议英文值 `idea / designing / contract_fixed / implemented / wired / accepted / released / deprecated`（中文八档一一对应）；若既有仓里已有生命周期 token 约定则以其为准（开工第一步核验）。
- **R-3 A 建议稿的挂载点候选**：建议稿内并列给出 1-3 选项（`candidates.ExecutionIntent.declared_impacts` / `base.GateCandidateEnvelope` 扩展 / 独立 `impacts` 子包）+ 各自后果，不预设结论。
- **R-4 B 是否顺手把 `market.PackManifest ↔ pack-manifest.schema.json` 登记为 go-schema-map 第 5 对**（把该结构纳入 Go↔Schema 一致性门禁，+0.5 小时）——**建议做**（现有 4 对不含它；required 语义经核对匹配：5 个必填字段均无 omitempty）。

## 2. 任务分解

### 任务 A：Impact Model 建议稿 + 影响清单（橙区，纯文档，~0.5-1 天）

产出文件：`docs/plans/impact-model-proposal-20260711.md`（本仓）

- [ ] A1 读三源设计输入并摘录字段：Palantir `modifiedEntities`/ActionLogicRule 14 变体（`/Users/li/Documents/systong/truzhen-notes/palantir-archaeology/01-schema-checkup.md`、`05-*.md`）；KWeaver `ImpactContractItem{ObjectTypeID,ExpectedOperation,AffectedFields[],Description}`；BKN impact_contracts 的 pre-conditions 阻断表 + `enabled` 默认 false（`kweaver-archaeology/K1-semantic-notes.md`）
- [ ] A2 起草 `declared_impacts[]` 形状：`{object_type, operation(create|update|delete|send|execute), object_ref?, affected_fields?[], description?}` + 默认关闸语义（未声明影响的正式动作是否拒绝=分级开关，给保守/严格两档方案）
- [ ] A3 起草 `receipt.actual_edits[]` 对账形状（与 `receipts.ReceiptEnvelope` 现有字段并列摆放，标注不改 required）
- [ ] A4 影响清单：逐个列消费方（truzhenos Base Gate/03 账本/06 引擎/client vendor/cloud）+ 兼容策略（全部 additive optional）+ 验收断言草案（declared vs actual 对账测试形状）
- [ ] A5 停。建议稿交 Owner 裁定挂载点与分级开关后，才有下一张实施卡。

### 任务 B：`lifecycle_status` 字段（黄区，~0.5 天，TDD）

前置：R-1/R-2 已裁。以下按 R-1=a（通用 manifest）写步骤：

- [ ] B1 全仓核验既有生命周期 token：`grep -rn "lifecycle" /Users/li/Documents/truzhen-contracts /Users/li/Documents/truzhen-packs --include='*.go' --include='*.json' --include='*.md' -l`；若有既有约定，token 以其为准并回改 R-2
- [ ] B2 **先写失败测试**：`market/pack_manifest_test.go` 增测——manifest JSON 带 `"lifecycle_status":"accepted"` 能解析入 struct 且非法值被校验 helper 拒绝（若本仓惯例无枚举校验 helper 则只测解析+常量存在）；跑 `go test ./market/` 确认 FAIL
- [ ] B3 `market/pack_manifest.go`：`PackManifest` 增 `LifecycleStatus PackLifecycleStatus \`json:"lifecycle_status,omitempty"\``+类型与 8 常量；`pack-manifest.schema.json` `properties` 增 `lifecycle_status: {type:"string", enum:[八值]}`（**不进 required**）
- [ ] B4 跑门禁：`go build ./... && go test ./... && go vet ./...`；`bash scripts/contracts-check.sh`（含 R-a 收紧后的 breaking-change 门：新增可选属性应 exit 0）；`VERSION` 0.8.0→0.9.0；检查 `embed.go` 覆盖
- [ ] B5（R-4 若批）`scripts/go-schema-map.json` 增第 5 对 `market.PackManifest ↔ pack-manifest.schema.json`，复跑 `contracts-check.sh` 确认 `mapped_pairs=5 passed_pairs=5`
- [ ] B6 packs 仓（先等 contracts 步完成）：4 个真实 pack + 1 个模板的 `manifest.json` 增 `lifecycle_status`（各 pack 现状档位由 Owner 或按 FEATURE_LEDGER 事实填，建议 env/housekeeping/smart-home=`accepted`、shuxuejia=`designing`、模板=`idea`——**填值属事实声明，开工时和 Owner 确认一遍**）；跑 packs 既有校验（`go test ./...` + pack_diagnostics）
- [ ] B7 突变自证：临时删 schema 中 `lifecycle_status` 属性 → 若 R-4 已做，consistency 门必 FAIL（证门禁真在看这个字段）→ 恢复 → 全绿
- [ ] B8 独立验收：派子代理在干净 checkout 复跑 contracts-check + packs 校验；两仓各自登记 FEATURE_LEDGER/账本；**不 merge 不 push，报 Owner 裁 land**

## 3. 验收断言汇总（改了什么证明什么）

- contracts：`contracts-check.sh` exit 0（22+ 工具测试、breaking-change、go↔schema 门全绿）；`git diff --check` 干净；VERSION bump 存在
- packs：manifest 校验绿；无 schema 之外的字段发明
- 反伪造：B7 突变必 FAIL 后恢复必 PASS
- 文档：本计划各任务勾选框如实更新；两项完成后统一决策表 #2/#3 状态回写
