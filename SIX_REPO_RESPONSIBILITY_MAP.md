# 六仓职责矩阵（单一权威来源）

> **制度（2026-07-11 收敛）**：六仓职责表此前 ≥6 份手抄、已实测三处漂移（版本号、client 范围、mobile 归属）。本文件是六仓职责矩阵的**唯一权威来源**；各仓治理文件中的同类表为便携副本，须在表前注明指针，与本文件冲突时以本文件为准。修改职责边界必须先改本文件，再同步便携副本。
>
> 哲学总纲见 `TRUZHEN_PHILOSOPHY.md`（六仓公认权威）；本文件只管「哪个仓负责什么」。

| 落点 | 权威位置 | 职责 | 不负责 |
| --- | --- | --- | --- |
| 主权、闸门、回执、网关、运行时、装载器 | `github.com/truzhen/truzhenos`（私有） / `/Users/li/Documents/truzhenos` | Base Gate、Owner 主权链、Receipt Ledger、四大网关（08 模型 / 09 记忆 / 10 沟通 / 11 执行）、场景流程引擎、Agent 编排、`/v3` API、本地 Cloud proxy / Session / License Gate 消费端 | 云端 server / 部署脚本、官网 / 市场页、支付云端实现、Pack 行业数据、前端 UI 源码 |
| 跨仓数据形状 / 契约 SDK | `github.com/truzhen/contracts`（开源 Apache-2.0） / `/Users/li/Documents/truzhen-contracts` | Pack / Candidate / Receipt / Surface / ReadModel schema、候选信封、门控裁定、监控事件、三主线引用、`*.schema.json` + embed；版本以 `VERSION` + 下游 `go.mod` 为准 | 任何实现（DB / 网络 / provider / runtime）、Pack 数据、UI、云端服务 |
| 行业工作台、流程、角色、知识、能力引用 | `github.com/truzhen/packs`（开放） / `/Users/li/Documents/truzhen-packs` | 场景荚 / 角色荚 / 能力荚引用（三类 Pack 封顶）、结构化知识、装入 / 卸载脚本、作者端模板；Pack 不持 Base 主权 | Base / Gateway / provider 实现、市场商品与支付真相（归 cloud）、前端渲染 |
| 本机外部软件 / provider / sidecar registry | `/Users/li/Documents/truzhen-software`（本机目录 + Git 治理文件） | Baserow、Frappe、OCR、IM、执行 sidecar、Ory 身份栈等外部软件登记（source / descriptor 可提交，runtime local-only）；由基座 02 解析为 ProviderResource | 业务主权、任何 Base 事实、云服务；providers 运行态不入 git |
| 官方云服务与云端网页 | `github.com/truzhen/truzhen-cloud`（私有） / `/Users/li/Documents/truzhen-cloud` | 云账号 / 会话 / 作者身份（Kratos / Keto）、商品 / 订单 / 支付 / License / Entitlement、Pack 文件分发（Forgejo）、官网与云市场页、rendezvous 中继、云端 release 治理 | Owner / Base / Receipt 主权、本地执行、Pack 业务逻辑、客户端源码 |
| 客户端（Web / Desktop / 移动） | `github.com/truzhen/truzhen-client-web-desktop`（私有） / `/Users/li/Documents/truzhen-client-web-desktop` | React + Vite + TS、Tauri 桌面壳、**移动遥控端唯一实现落点 `apps/mobile-remote/`（Owner 2026-07-06 裁定；基座 `modules/24` 仅治理 overlay）**；只做 ReadModel / UI Projection / 候选卡投影 | 真相源（Frontend ≠ Truth Source）、主权裁定、直接写库 / 发送 / 执行 |

## 漂移防线

- 版本号：任何治理文件不得手抄 contracts 版本数字，一律「以 `VERSION` / `go.mod` 为准」。
- 便携副本改动：先改本文件 → 再同步各仓副本；发现副本与本文件冲突按本文件纠偏并回查漂移原因。
- 登记点：truzhenos `V3_GOVERNANCE.md` / `README.md`（2026-07-11 已加指针）；其余仓副本随后续治理轮补指针。
