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

---

# 第二阶段（Owner 2026-07-11 追加裁定：「推动全部完成」）

逐项「缺什么 → 完成路径」与波次。事实底座：四路扫描（其中 05 §6 与 hands 两路误扫冻结仓 truzhenv3，其结构性结论采信、行号与现状需在 truzhenos 复核后才施工）。

| # | 缺什么 | 完成路径 | 波次/量级 |
| --- | --- | --- | --- |
| #13 版本三态 | **几乎不缺**：lifecycle 端点（draft/promote/confirm/disable/new-version/reactivate/history/packs）、version_events 存储、前端 mergeScenePackLifecycleReceipts 投影、管理页版本展示全部已存在 | 复核后定性「已存在」，决策表收口；如有小 UI 缺口（三态视觉统一/new-version 入口）单独小补 | W1（复核收口，0-0.5 天） |
| #6 双层裁剪 | 治理条款与两层实现（02 五步管线 / 09 KnowledgeMount 六维键）大概率已存在；冻结仓证据需 truzhenos 复核；可能残缺口=05 工作区记忆两分接缝（planned） | truzhenos 复核 → 有缺口补缺口，无缺口定性「准入已满足/已实现」收口 | W1（复核，0.5 天） |
| #18 sandbox 双层隔离 | 缺定性归档：无本地沙箱需求（11 走 provider 模式）+ 源仓无 LICENSE（思路引用须注明出处） | 写 V4 概念档案文档（clean-room、出处、触发条件），档位=想法→设计中归档 | W1（0.5 天） |
| #7 Receipt context | 缺的是**推翻既有裁定的确认**：原裁定不补理由=「无回放缺口证据 + generic decision_context 重复存事实」（palantir 01-schema-checkup L67-84 原文） | 建议以 **typed context 按需加** 条款收口（#14 的 knowledge_recall 已是第一个 typed context 实例=此项已以正确形态开始完成）；不做 generic 大袋子字段。条款入治理 | W1（条款化，0.5 天） |
| #11 RiskType 五件套 | 缺 contracts additive `risk_types[]` + **Base Gate 真实消费点**（同轮接线铁律）+ 至少 1 个 pack 声明样例 + 双向突变 | 按 Impact Model 模式三仓实装：contracts 形状 → os gate 消费（声明 escalation=owner_gate 且匹配时 Allow→PendingOwner）→ packs 样例声明 | W2（1-2 天） |
| #15 per-agent 身份 | 缺身份模型与消费者；完整 token 身份=16/03 红区且无场景 | 两段完成：①现在实装「归因收口」＝候选/回执统一带 role_pack_ref+slot_ref+task_purpose 三元组（补缺处+守卫测试）；②独立 token 身份出实施卡挂触发条件。若 Owner 坚持 token 现在做，先出 16 红区实施卡再施工 | W2.5（①0.5-1 天） |
| #12 制作台引导 | 前置物判定推翻：三制作台（场景画布/能力四幕/角色三步+StudioShell）已上线。真缺三件=填空脚手架、@ 联想、1 必填渐进展开 | client 仓独立分支改造三台表单容器（渐进展开容器+联想组件），562 测试回归+build | W3（2-3 天，client 仓） |
| #5 LinkType | 缺形状与施工。**双真相源风险已有解**：LinkType 注册表（05 契约对象）+ RelationEdge 增 `link_type_ref` additive（不建第二张边表；RelationType 字符串与 ref 并存=方案A 冗余向后兼容）；写路径已过 Base gate 验真三件套 | contracts schema+Go → 05 注册表存储+端点 → RelationEdge additive+校验 → 测试；实施前在 truzhenos 复核冻结仓证据 | W4（2-3 天） |
| #17 exec resume | 缺的是**在途线程收尾，不是新施工**：codex-hands 主线正在执行（worktree w0-w1-execution + w2-delegation-session 均在跑，contracts main 今日 ac80e72 即其 land）；我方抢做=撞车重复造 | 挂到 codex-hands 线程：把 exec resume 验收点（token 幂等/门禁三证/断点恢复）交叉核对进其 W1 验收；该线收尾时本项随之完成 | 协调项（不施工） |

波次执行：W1 全部（复核收口+文档+条款）→ W2 #11 → W2.5 #15① → W3 #12 → W4 #5。#17 只做协调核对。每波独立分支、独立验证、停等 land。

## W2 实施卡：#11 RiskType 五件套（侦察后定稿，2026-07-11）

侦察事实（file:line 见侦察报告）：GateCandidateEnvelope 无 pack 回溯字段、DeclaredImpacts 当前**零生产者**（floor 已 live）；缺口 A=scenepack/spec.go 无 risk_types、B=confirm_flow buildGateRequest 未提取、C=sceneflowdev handlers.go:786 buildSceneFlowGateRequest 有 pack_version_ref 但未查 lifecycle record（GetRecord 接口已存在 ports.go:61）。

- **contracts**（本 campaign 分支）：①`spines/risk.go` 新增 `DeclaredRiskType{RiskTypeID, TriggerActionType, EvidenceRequirement, EscalationPath("none"|"owner_gate"), Definition}`（envelope 用，避免 base→market 依赖）；②`market.PackRiskType` 五件套（risk_type_id/definition/trigger_action_types[]/evidence_requirement/escalation_path/fallback）+ PackManifest.RiskTypes additive + pack-manifest.schema.json additive + enum 同步守卫测试；③`base.GateCandidateEnvelope.DeclaredRiskTypes []spines.DeclaredRiskType` additive；④VERSION minor bump（以分支基线 ac80e72 的 VERSION 为准）。
- **os**（wave2 分支 claude/backlog-wave2-os-20260711）：①scenepack spec 增 risk_types（跟随 ProviderRequirements 既有类型风格）+ lifecycle draft intake additive；②base/gate 新 `risk_type_gate.go`（镜像 impact_risk_gate：declared 含 escalation=owner_gate → Allow→PendingOwner + trace "RiskTypeGate"）+ orchestrator 两处 Evaluate 调用点接 applyRiskTypeFloor（紧跟 applyImpactRiskFloor 后）；③生产者=06：buildSceneFlowGateRequest 按 pack_version_ref 查 lifecycle record → 匹配 trigger_action_types → 附 DeclaredRiskTypes（KWeaver 规格-only 反例的解药就在这一跳）；④TDD：floor 三态单测 + 06 run-start 带声明必 pending_owner 的行为测 + 突变闭环。
- **packs**（campaign 分支）：env pack manifest 增 1 条真实 risk_types 样例 + 四 install.py draft payload 透传 risk_types（additive）。
- 顺序：contracts 形状+守卫绿 → os replace 本地 contracts 实装（land 时按 G2a 翻 require）→ packs 样例。验收=改哪证哪 + 全量回归。
