# 学习与探讨合议契约兼容说明（v0.15.0）

## 变更范围

`v0.15.0` 新增加性、只读的跨边界形状：

- `readmodels.DeliberationSessionReadModel`
- `readmodels.DeliberationTurnReadModel`
- `readmodels.DeliberationProviderLaneReadModel`
- `readmodels.DeliberationAutomationGrantReadModel`
- `candidates.DeliberationSynthesisCandidate`
- 上述五个类型的 JSON Schema、`embed.go` 导出与安全 client codegen 样例。

这不是一套新的主权或执行协议。Base 是 `decision_ref` 的签发与反查权威；`OwnerDelegationGrant` 仍是委托授权的既有契约；Receipt Ledger 仍是 `receipt_ref` 的唯一真相源。本次不修改它们，也不新增平行 grant、Receipt 或 Intent target。

## 不变量与数据边界

- 所有会话、轮次、通道和授权投影固定 `candidate_only=true`、`non_formal=true`；它们不能直接形成 Formal Record。
- 问题、提示词、供应商回答、网页 DOM、cookie、凭据和原始网络载荷不进入本仓类型或 schema。跨边界只传 `*_artifact_ref`、SHA-256、受控引用、状态和脱敏失败摘要。
- `decision_ref`、`policy_snapshot_ref` 和 `evidence_refs` 只是服务端签发事实的只读投影。client 不能以字符串拼接、`Date.now()` 或随机值自铸它们。
- `current_turn_auto` 必须同时带 `turn_ref` 与小写 SHA-256 `question_sha256`；`max_dispatches_per_lane` 固定为 `1`，`dispatch_on_confirm` 固定为 `true`。Base 确认成功后必须立即进入既定 lane 的 dispatch，不存在第二个开始动作。
- 合议输出只能是 `DeliberationSynthesisCandidate`。每个结论项必须引用同一个候选所列 `material_refs`；`receipt_ref` 必须来自真实账本，不能以 `receipt://` 模板假造。
- 通道发布资格（`release_eligibility`）与运行时就绪度（`runtime_readiness`）分开表示；`not_ready`、`blocked`、登录、验证码、adapter drift 与恢复需求都必须如实投影，不能伪装成功。
- Session、Turn、Lane 与 Grant 的校验器对未知枚举、跨 Session/Turn 混入、错误 hash、缺失已导入 lane 的 Artifact/Receipt 引用均 fail-closed；其校验只验证形状和引用关系，不替代 Base 或 os-03 的真实反查。

## SemVer 与下游迁移

这是 `v0.14.0 -> v0.15.0` 的 minor 加性变更：未删除或改名既有字段、未收紧既有 schema、未改变 `OwnerDelegationGrant` 的默认语义，旧消费者无需修改即可继续工作。

有真实消费需求的下游按以下顺序接线：

1. `truzhenos` 以新版本 module 实现会话、受控 artifact、Base 决策反查、Gateway 及 Receipt 写入；不得把本契约当真相源。
2. client repo vendor 五个 schema 或由 schema codegen 生成只读类型；可使用 `scripts/tests/fixtures/deliberation/client-codegen-projection.json` 校验字段对齐，但不得把该 fixture 当运行时 fallback 数据。
3. client 只能展示候选、请求审批或提交已有的 `idempotency_key`；不得构造 `decision_ref`、`owner_action_evidence_ref`、`run_id`、`nonce` 或正式回执。
4. `truzhen-packs`、`truzhen-software` 和外部 provider 若后续成为消费者，必须先登记真实消费点和反向依赖，再新增其专用字段；不得从本契约推断 provider 可真实执行。

发布前须由 Owner 选择 tag 并发布 module；发布后，OS 与 client 应各自做 schema/vendor 或 codegen 兼容验证。未发布前生命周期仍为“设计中（实现待发布）”；只有 Owner 发布兼容版本后才进入“契约已定”，更不是“已接线”“已验收”或“已发布”。
