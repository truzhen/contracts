# Codex Hands 执行事实 fail-closed 兼容报告（v0.13.0）

## 变更

`DelegationGrantWithinScope` 现在要求执行授权的 grant 和服务器派生 subject 在执行维度上对称：

- grant 含 `scope.execution_scope` 时，subject 必须含 `execution`；缺失时返回 `delegation subject execution facts are required by execution_scope`。
- subject 含 `execution` 但 grant 不含 `execution_scope` 时，继续按既有逻辑拒绝。
- grant 与 subject 均不含执行维度时，普通非执行委托行为保持不变。

## 风险与兼容性

这是 additive 契约字段发布后的 fail-closed 行为修复，不改变 JSON 形状、字段类型或枚举。正确填写服务器执行事实的 `v0.11.0`/`v0.12.0` 消费方无需修改；遗漏 `DelegationSubject.Execution` 的执行授权调用将从错误放行改为明确拒绝。

`truzhenos` 的 Codex Hands 授权器已经从服务器派生 capability、workroot、provider、sandbox、network 和预算事实；升级后应运行 delegation、T06、EmergencyStop 和 Docker 三轮续跑回归，证明没有调用点依赖缺失执行事实的旧行为。

## 发布顺序

1. contracts 功能分支通过 build/test/vet 和反向依赖检查。
2. Owner 单独批准并发布 `v0.13.0` tag/release。
3. `truzhenos` 删除旧版本依赖缓存后升级 `github.com/truzhen/contracts v0.13.0`，禁止本地 `replace`。
4. client 只消费 `truzhenos` 投影，不直接解释缺失执行事实。
