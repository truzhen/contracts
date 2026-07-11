# Backlog 提前触发战役计划（统一决策表 #5-#18）

- 日期：2026-07-11
- 触发：Owner 指令「这些先做了」——统一决策表（`systong/truzhen-notes/kweaver-archaeology/unified-decision-table-20260710.md`）#5-#18 提前触发。
- **裁定记录**：这些条目原设计为「backlog 带触发条件」；Owner 2026-07-11 裁定提前触发，本计划按此执行并入档。执行 Agent 已履行反向提醒义务（节奏闸门/缺证据项如实标注）。

## 版本/优先级

Owner 裁定提前触发（原为 backlog）。执行中仍按风险颜色分层：绿/黄实装、橙只出影响评估、前置物缺失项不做空壳。

## 真实客户/场景证据

- #8/#10：**有间接证据**——pack 交付质量项（install 半装、知识漂移是真实交付风险，四个 pack 已实测过安装链路）。
- #14：**有间接证据**——ENV-R4-01 有界召回透明性已是既定规范，缺的是可回放落点。
- #9/#16：设计/文档级，消费材料已备（K2 笔记）。
- #5/#11/#15：**缺证据**（尤其 #15 per-agent 身份），故只出影响评估不施工。

## 最小可交付（本轮范围）

| 组 | 条目 | 交付物 | 仓 |
| --- | --- | --- | --- |
| 实装 | #10 knowledge 内容 checksum 防漂移 | knowledge-index 每 entry 增 `checksum`（sha256 原始文件字节）+ 生成器脚本 + CI 校验步 + install.py 前置校验（新错误码）| packs |
| 实装 | #8 install 事务日志+断点续装 | install.py 安装 journal（步骤+refs 落盘）、失败时半装状态显式化报告、重跑幂等续装；**不自动反做正式域对象**（撤销走 Owner 侧既有禁用状态机） | packs |
| 实装 | #14 召回预算参数化入 Receipt | ContextCompileResult 增 total/truncated 参数化字段并落入 09 compile receipt（additive） | truzhenos |
| 设计 | #9 .skill 三段式设计稿 | proposal 文档（references/assets 分区进技能包规格 + SKILL.md 负面清单），停等裁定 | contracts docs/plans |
| 设计 | #16 GenAI span 字段对照 | truzhen-monitor 字段字典对照文档 | truzhenos docs |
| 评估 | #5 LinkType / #11 RiskType 五件套 / #15 per-agent 身份 | 各一份影响评估（契约影响+兼容策略+成本），停等裁定 | contracts docs/plans |

**明确不做**：#6（准入条件 05 §6 未满足）、#12/#13（14 制作台/版本表面不存在，做了即空壳，违反成品门禁）、#7（已裁不补）、#17（已并入 codex-hands 计划）、#18（V4 概念，源仓无 LICENSE）。

## 真相源

- 知识内容真相源=pack 仓文件本体；checksum 是完整性声明非第二真相源。
- 安装进度真相源=os 侧各状态机（06/14 lifecycle/09 candidates）；journal 是**本机恢复辅助记录**，不是真相源，冲突时以 os 端点实况为准。
- 召回统计真相源=09 compile 时刻事实，落 receipt 即不可变回放事实。

## 仓库/层归属

packs（#8/#10）、truzhenos 09（#14）、contracts docs（#9/#5/#11/#15 文档）。**contracts 代码/schema 本轮零改动**。

## 风险颜色

#10 绿；#8 黄（改安装链路，land 前需 Owner 确认）；#14 黄（additive+受 09 红区棘轮约束，撞棘轮即停报）；#9/#16 绿（文档）；#5/#11/#15 橙（只评估）。

## 契约影响

无。#14 若 ContextCompileResult 在 contracts 有镜像 schema 则停下上报（预扫未见）；#5/#11/#15 的契约影响写在各自评估文档内，本轮不改。

## 上下文与禁止边界

- 允许读：三仓全量、systong 考古笔记。
- 禁止：动 Base Gate/主权链语义；install.py 自动撤销 formal 对象；新增平行监控/账本；碰他会话 worktree 与未跟踪文档；09 红区加行为代码超出 additive 必要量；merge/push（停等 Owner）。

## 验收设计（改了什么证明什么）

- #10：突变测试——改一个知识文件字节→CI 校验必 FAIL→重新生成 checksum→PASS；install.py 前置校验对 mismatch 出新错误码。
- #8：单元级模拟分步失败（monkeypatch HTTP）→journal 内容断言→重跑续装跳过已完成步→完成态 journal 标记；半装报告含悬空 candidate_refs。
- #14：TDD——compiler 测试断言新字段与 total 一致；compile receipt 反查含召回统计；恢复旧行为的突变必红。
- 两仓各自全量验证：packs CI 等价本地跑 + truzhenos `bash scripts/verify.sh` VERIFY_EXIT=0 + go test -race（范围按 go-test-packages.txt）。
- 独立验收：实装完成后由独立子代理复跑关键断言。

## 变更影响

packs：5 个 pack 的 knowledge-index/install.py 同构更新 + CI workflow 增一步；os：09 contextcompiler+compile receipt 路径 additive；无 UI/契约/部署影响。

## 生命周期档位

#10/#8/#14 目标：已实现→已接线→（Owner land 后）已验收；#9/#16：设计中；#5/#11/#15：想法→设计中（评估文档）。
