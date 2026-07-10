# Codex Hands 委托执行边界兼容说明（v0.9.0）

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

`DelegationExecutionSubject` 是服务器派生的拟执行 run 后累计事实：

- `capability_ref`
- `workroot_ref`
- `provider_ref`
- `sandbox_profile_ref`
- `network_policy`
- `consumed_runs`
- `consumed_duration_seconds`

比较规则由 `DelegationExecutionWithinScope` 表达：scope 为空时拒绝 execution subject；subject 为空时也不能授权 execution scope；ref 字段按不透明字符串比较，不解析本机路径。

## 网络策略

`ExecutionNetworkPolicy` 枚举值：

- `none`
- `egress_model_only`
- `gated_bridge`

`gated_bridge` 可以作为服务器派生 subject 的事实值出现，但 grant 的 `network_policy_ceiling` 只允许 `none` / `egress_model_only`。因此 subject 为 `gated_bridge` 时不会被本次委托授权放行。

## SemVer 与下游影响

这是 `v0.8.0 -> v0.9.0` 的 minor 兼容扩展：

- 不删除旧字段。
- 不改变旧 JSON tag。
- 不修改任何 JSON Schema。
- 不改变 high / critical 不可委托硬地板。
- 不改变旧 `OwnerDelegationGrant` 的授权语义。

下游接入顺序建议：先升级 contracts，再由 `truzhenos` Base / Execution Gateway 使用服务器派生 subject 调用 `DelegationExecutionWithinScope`。前端、Pack、Provider 不应自铸 execution subject。
