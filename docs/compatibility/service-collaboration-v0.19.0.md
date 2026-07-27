# v0.19.0 服务协作契约兼容说明

本版本新增 `service-collaboration.schema.json` 与五个 `readmodels` 形状：邀请、约定披露、对手方决定、每日交付披露、参与者 scope 投影。它们均为加法；不修改旧字段、旧 schema、既有移动配对行为或 OwnerDecision 语义。

OS 是生产方，Client 是只读消费方或对手方证据提交方；正式 Agreement、scope、会话、Gate、Receipt 与执行仍归 OS。新字段不包含 token、身份明文、订单全文、视频原件、Prompt、成本或 Provider 配置。未升级的消费者不读取新 schema，因此保持兼容。

后续 W03/W04 必须固定本提交 SHA，分别完成 OS producer 与 Client vendor/codegen，再以真实集成测试证明“已接线”；本文件不声称已发布或已接线。

`nonce` 的一次性消费、proposal hash 与 transaction/participant 的绑定、防重放、会话撤销与 Receipt append 是 OS 真相源职责；contracts 只校验这些 ref/nonce 非空，不能用无状态 helper 冒充运行时防重放。
