# Impact Model 建议稿 + 影响清单（橙区设计稿，零 schema/代码改动）

> 日期：2026-07-10　|　生命周期：**设计中**（本稿交 Owner 裁定挂载点与分级开关后，才有实施卡）
> 来源：统一决策表 #2（三源合流）。本稿只画形状；任何字段落地前须按 §5 实施卡切法执行。

## 1. 三源设计输入（带出处）

| 来源 | 形状 | 关键语义 |
| --- | --- | --- |
| Palantir OSDK | `modifiedEntities` = **执行前声明**；`ActionResults` 的 added/modified/deleted objects+links = **执行后事实** | 两端分离：声明 ≠ 事实（palantir-archaeology/01-schema-checkup.md:36-47） |
| KWeaver bkn-backend | `ImpactContractItem{ObjectTypeID, ExpectedOperation(add\|modify\|delete), AffectedFields[], Description}` | 影响契约结构化到字段级（K1-semantic-notes §1.5，SPECIFICATION.md:439-450） |
| BKN 规格 | `Pre-conditions` 阻断表（不满足则阻止执行）+ 影响契约 `enabled` 默认 false | **默认关闸**、渐进启用（K1-semantic-notes；SPECIFICATION.md:462-469） |

01 号报告原判据（照录）：「两者均为可选 additive 字段；`delete` 必须在 Base/Gate 风险策略中单列，Receipt 只写执行后事实。**不得让 Candidate 声明直接变成正式写入授权**。」

## 2. 建议形状（Go 伪码，实施时按仓惯例定稿）

```go
// ImpactOperation：v1 只覆盖对象编辑域。send/execute 已由 SideEffectClass
// 与既有候选类型表达（gate.go:57），重复声明会造双真相源，故不纳入。
type ImpactOperation string // "create" | "modify" | "delete"

// DeclaredImpact：执行前声明（Proposer 视角的「我将碰什么」）。
type DeclaredImpact struct {
    ObjectType     string          `json:"object_type"`               // 05 业务对象类型标识
    Operation      ImpactOperation `json:"operation"`
    ObjectRef      string          `json:"object_ref,omitempty"`      // 已知对象时填
    AffectedFields []string        `json:"affected_fields,omitempty"` // 字段级（KWeaver 输入）
    Description    string          `json:"description,omitempty"`
}

// ActualEdit：执行后事实（Receipt 视角的「实际碰了什么」）。ObjectRef 必填。
type ActualEdit struct {
    ObjectType     string          `json:"object_type"`
    Operation      ImpactOperation `json:"operation"`
    ObjectRef      string          `json:"object_ref"`
    AffectedFields []string        `json:"affected_fields,omitempty"`
}
```

## 3. 挂载点三选项（R-3，各自后果）

| 选项 | 位置 | 后果 | 建议 |
| --- | --- | --- | --- |
| **a** | `base.GateCandidateEnvelope` 增 `declared_impacts,omitempty`（gate.go:46-64，与 `SideEffectClass`/`FormalWrite` 并列）；`receipts.ReceiptEnvelope` 增 `actual_edits,omitempty`（receipt.go:5-15） | 门裁定处声明、回执处对账，链路最短；只动 base+receipts 两处；candidate-envelope.schema.json 不动，receipt-envelope.schema.json additive | **建议** |
| b | `candidates.CandidateEnvelope` 声明（出生即带）+ receipts 对账 | 声明更早，但每个候选生产方都要理解该字段，波及面大；且候选→GateEnvelope 需透传一次 | 备选 |
| c | 新建 `impacts/` 子包供双方引用 | 分层最干净；但**新增子包属 AGENTS §0.7 必回 Owner 项**，且两个小 struct 撑不起一个子包 | 不建议 v1 |

## 4. 分级开关（防 KWeaver 式纸面契约，也防一步收太紧）

- **保守档（v1 默认）**：`declared_impacts` 缺省不改变任何现有裁定行为——纯 additive 留痕+前端影响预览；唯一硬语义=声明含 `delete` 的候选自动升 `RiskClass=high` 下限（进 owner_gate 路径，呼应 01 判据「delete 单列」）。
- **严格档（Base policy 可配，默认关）**：`impact_declaration_required=true` 时，`FormalWrite` 非空而 `declared_impacts` 为空 → Gate 拒绝（blocked，留痕）。
- **对账语义**：`actual_edits` ⊄ `declared_impacts`（按 object_type+operation 匹配）时产 `impact_reconciliation_mismatch` 监控事件（走 truzhen-monitor 既有链，不另起）；严格档下升审计 finding。**任何档位下声明都不构成授权**——授权仍只来自 Owner+Base Gate。

## 5. 实施卡切法（铁律：规格与运行时同轮接线，不重蹈 KWeaver「形状 land 了运行时永远没来」）

一张实施卡必须同轮包含，缺一不 land：
1. **contracts 形状**：按 R-3 裁定挂载 + `receipt-envelope.schema.json` additive 属性（不动 required）+ golden/同步测试 + VERSION minor（注意 0.9.0 已被 lifecycle_status 占用，届时为 0.9.x→0.10.0）+ contracts-check 全绿；
2. **truzhenos 最小接线**：Base gate 把 `declared_impacts` 写入裁定上下文与 GateTrace；03 账本在**至少一条真实写路径**（建议 05 事务对象正式写）产 `actual_edits`；delete 升险规则生效并有测试；
3. **对账验收断言**：单测=不匹配必产 mismatch 事件；突变自证=注释掉升险规则→测试必 FAIL。

## 6. 影响清单（消费方逐个）

| 消费方 | 影响 | 兼容策略 |
| --- | --- | --- |
| truzhenos Base Gate/01 | 读 declared_impacts（可选字段，缺省行为不变）+ delete 升险规则 | additive；老候选无字段=保守档原行为 |
| truzhenos 03 回执账本 | 写 actual_edits（仅接线的写路径） | additive；hash 链不受影响（字段进 payload hash 正常参与） |
| 06 场景流程引擎 | 场景荚可在 controlled_execute 节点声明影响（透传，不裁定） | 后续卡，可选 |
| client vendor/codegen | receipt-envelope（及 R-3=b 时 candidate-envelope）schema 变更需 vendor 同步 | additive，codegen 重跑即可 |
| truzhen-cloud | 无 | — |
| 防漂移门禁 | 新增可选属性=兼容（本轮 lifecycle_status 已实证该判定路径） | exit 0 预期 |

## 7. 待 Owner 裁定

- **O-1** 挂载点：a / b / c（建议 a）。
- **O-2** 保守档 delete 升险是否 v1 就带（建议带——这是唯一让字段「活」的最小运行时语义）。
- **O-3** operations 是否扩 send/execute（建议不扩，理由见 §2 注释）。
- **O-4** 实施时机：按节奏治理属支线小件，建议排智能家居跑穿的空档；os 侧接线量约 1-2 天，超出「≤1.5 天」原估，**如实上报**。
