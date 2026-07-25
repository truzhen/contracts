# P13 Experience / Pack / Market contracts v0.17.0 兼容说明

状态：契约已定，未发布。`VERSION` 已更新为 `0.17.0`；tag、push、下游 module 更新和真实 paired-host 注册仍需 Owner 单独授权。

## 变更范围

- 新增 `experience-asset-candidate-readmodel.schema.json` 与 `readmodels.ExperienceAssetCandidateReadModel`。它仅投影已脱敏经验候选、来源 Receipt/事务/固定 Pack 版本、candidate-only 状态及已有候选/回执引用；明确不含 raw Receipt payload、OwnerDecision、token 或市场上传证明。
- 新增 `approved-pack-artifact-handoff-attestation.schema.json` 与 `market.ApprovedPackArtifactHandoffAttestation`。它要求 paired host 对固定 `candidate_ref`、`approval_receipt_ref`、artifact digest、已签名清洗清单摘要、Pack/version、受众、时效、nonce 和证据引用的 canonical JSON 进行 Ed25519 签名。
- `market.CanonicalApprovedPackArtifactHandoffBytes`、`SignApprovedPackArtifactHandoff`、`VerifyApprovedPackArtifactHandoff` 均为无副作用的纯 helper；不保存私钥、不绑定账号、不上传制品。账号↔host 公钥绑定与幂等持久化由 cloud 真相源实现。

## 兼容性

这是 `v0.16.0 -> v0.17.0` 的 minor 加性变更：未删除、改名或收紧任何既有 schema 字段。旧消费者不需要读取新 schema 即可继续使用 v0.16.0 的形状。

新 handoff consumer 必须 fail-closed：opaque Receipt/candidate ref 不能替代签名；签名校验前必须同时确认 authenticated account 与 `issuer_host_ref` 的绑定、`audience`、时效、artifact SHA-256、Pack/version 与重放 nonce。ReadModel consumer 不得把投影字段回传成 Owner authority 或正式市场结果。

## 下游顺序

1. Owner 授权发布 `v0.17.0` tag 后，cloud 先实现 account-bound paired-host 校验、durable handoff operation 与 pending-review 事实。
2. `truzhenos` 再实现 Base 审批后的 handoff candidate、Cloud Market Hands 与外部动作 Receipt；os-18 不直连 cloud。
3. client vendor 两份 schema，运行 `npm run gen:contracts-types`，只消费生成类型与后端 ReadModel。

未发布 tag 前，不得用本地 `replace`、手抄 DTO 或临时 vendor 副本宣称三仓已接线。
