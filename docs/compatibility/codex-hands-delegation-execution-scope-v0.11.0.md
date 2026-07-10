# Codex Hands 委托执行边界兼容说明（v0.11.0）

## 变更范围

本次只扩展 `base.OwnerDelegationGrant` 既有委托契约，不新增 `EnvelopeGrant` 或任何平行授权。新增字段均为 Go JSON 可选字段：

- `DelegationScope.execution_scope`
- `DelegationSubject.execution`

旧 JSON 不包含 `execution_scope` / `execution` 时，反序列化和重新序列化行为保持不变；缺省状态不代表任何代码执行授权。

## 新增边界

`DelegationExecutionScope` 声明 Owner 预授权允许的代码执行上限：

- `capability_refs`
- `workroot_ref`
- `provider_refs`
- `sandbox_profile_ref`
- `network_policy_ceiling`
- `max_runs`
- `max_duration_seconds`

`DelegationExecutionSubject` 是服务器派生的、已经原子预留待执行本次 run 后的累计事实。为保持 JSON 向后兼容，字段名不变：

- `capability_ref`
- `workroot_ref`
- `provider_ref`
- `sandbox_profile_ref`
- `network_policy`
- `consumed_runs`
- `consumed_duration_seconds`

`consumed_runs` 和 `consumed_duration_seconds` 都必须大于 0，且分别不得超过 `max_runs` 和 `max_duration_seconds`。并发消费方必须先用 OCC（例如版本号 compare-and-swap）或等价的原子 compare-and-reserve 操作预留本次 run，再构造 subject 并校验；禁止先读、校验通过后再普通递增，否则并发 run 可能共同越过预算。

## 组合校验入口

代码执行授权裁定只能使用 `DelegationGrantWithinScope(grant, subject, evaluationTime)`。`evaluationTime` 必须由调用方显式传入；helper 不调用 `time.Now()`，避免隐式时钟导致不可复现或不可测的到期判断。该入口按固定顺序执行：

1. 复用 `ValidateOwnerDelegationGrant` 校验 `grant_id`、`owner_decision_ref`、`delegate_ref`、scope、`expires_at` 和 status 枚举；该通用结构校验保持旧 grant 的 `revocable=false` 行为兼容。
2. 代码执行完整入口额外要求 `revocable=true`，并要求正式 grant 具有非空 `receipt_ref`。为保持旧 JSON round-trip 兼容，字段形状仍为可选，但不可撤销或缺少回执证明的 grant 不能授权执行。
3. 要求 `evaluationTime` 非零、grant status 必须为 `active`，且 `expires_at` 必须严格晚于 evaluation time；`revoked`、`expired`、`suspended_by_emergency_stop` 均拒绝，恰好在 evaluation time 到期也拒绝。
4. 调用 `DelegationWithinScope` 校验 task、risk hard floor / ceiling、transaction、Pack 和 amount。
5. 仅当 `DelegationSubject.execution` 非空时，再调用 `DelegationExecutionWithinScope` 校验执行维度。

因此 inactive / expired grant、high / critical、父级 scope 越界和旧 grant 缺少 `execution_scope` 都不能被 execution helper 绕过。`DelegationExecutionWithinScope` 仍保留为执行维度的低层兼容 helper，但它不校验 grant 生命周期或完整父级 subject，不能单独授权代码执行。ref 字段按不透明字符串比较，不解析本机路径。

### 完整入口硬不变量审计

| 不变量 | 完整入口 fail-closed 位置 |
| --- | --- |
| 必要 refs | `grant_id`、`owner_decision_ref`、`delegate_ref` 由 `ValidateOwnerDelegationGrant` 拒绝空值；正式授权另要求非空 `receipt_ref`。 |
| 可撤销 | 通用 `ValidateOwnerDelegationGrant` 保持旧结构校验行为；代码执行完整入口额外拒绝 `revocable=false`。 |
| status | 结构校验拒绝未知枚举；完整入口只接受 `active`。 |
| expiry | 结构校验拒绝零 `expires_at`；完整入口要求 `expires_at > evaluationTime`。 |
| EmergencyStop | `suspended_by_emergency_stop` 作为服务器派生 grant 状态由完整入口拒绝；Base 必须在 EmergencyStop 启用时把 active grant 投射为该状态，调用方不得传入陈旧 active 快照。 |
| grant 预算 | `quota.per_day`、`max_runs`、`max_duration_seconds` 必须大于 0；`amount_limit_cents` 允许 0 表示不授权金额动作，但禁止负数。 |
| execution 累计预算 | `consumed_runs`、`consumed_duration_seconds` 必须大于 0，且不得超过 scope 上限。 |

`created_at`、`updated_at` 是可选审计元数据，不参与授权；`receipt_ref` 虽保持可选 JSON 形状，但完整授权入口将其作为正式 grant 的必要证明。

## 网络策略

`ExecutionNetworkPolicy` 枚举值：

- `none`
- `egress_model_only`
- `gated_bridge`

授权顺序为 `none < egress_model_only`：

- ceiling 为 `none` 时，只接受 subject `none`。
- ceiling 为 `egress_model_only` 时，接受 subject `none` 或 `egress_model_only`。
- `gated_bridge` 可以作为服务器派生 subject 的事实值出现，但 grant ceiling 不允许该值，任何合法 execution scope 都不会授权它。

## SemVer 与下游影响

这是 `v0.10.0 -> v0.11.0` 的 minor 兼容扩展：

- 不删除旧字段。
- 不改变旧 JSON tag。
- 不修改任何 JSON Schema。
- 不改变 high / critical 不可委托硬地板。
- 不改变旧 `OwnerDelegationGrant` 的授权语义。

下游接入顺序建议：先升级 contracts，再由 `truzhenos` Base / Execution Gateway 使用完整 grant、服务器派生 subject 和显式 evaluation time 调用 `DelegationGrantWithinScope`。前端、Pack、Provider 不应自铸 execution subject，也不得用低层 `DelegationExecutionWithinScope` 单独作授权裁定。
