# 本地优先 AI 事务操作层 contracts v0.16.0 兼容说明

状态：契约已定，未发布。`VERSION` 已更新为 `0.16.0`；tag、push 和下游接线仍需 Owner 单独授权。

## 变更范围

- `market.ProviderRequirement` 新增可选 `software_requirement_refs[]`，用于表达一个 Hands 对多个软件组件的复合依赖。
- `market.SoftwareResolutionLock` 新增四个可选投影字段：`provider_requirement_ref`、`provider_family`、`resolved_control_method_ref`、`resolved_execution_mode`。
- `scene-runtime-plan-candidate.schema.json` 的 runtime node 新增可选 `provider_requirement_refs[]`；出现时必须是固定 Pack 版本内唯一的单一引用。
- 新增 `pack-hands-requirement-readmodel.schema.json` 和 `readmodels.PackHandsRequirementReadModel`，供服务端向 Owner 展示解析和阻断事实。

## 兼容性

以上均为加性变更。旧的 `0.15.0` payload 不需要新增字段即可继续通过原有 required 约束；旧消费者忽略新增可选字段即可保持兼容。Provider manifest 中一旦声明 `software_requirement_refs[]`，数组必须非空、去重；ReadModel 的依赖和 lock 投影数组允许为空但拒绝空 ref 和重复 ref。跨 Pack 引用和 family 一致性由纯数据校验 helper 在固定 Pack 版本上下文中拒绝。

ReadModel 中的 provider/resource/control method ref、解析状态和安装候选是服务端只读投影，不能被 client 回传用来强制选择实例，也不等价于 Owner/Base 授权或真实执行结果。

## 下游顺序

按 `contracts → truzhen-software → truzhenos → truzhen-packs → truzhen-cloud → client` 顺序消费。未发布 tag 前，不得用本地 `replace`、手抄 schema 或临时 vendor 副本宣称接线完成。
