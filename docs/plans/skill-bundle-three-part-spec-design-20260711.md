# 技能包三段式内容规格 · 设计稿（决策表 #9，设计-only 不施工）

- 日期：2026-07-11　|　生命周期档位：**设计中**（停等装修 Pack 作者轮验证后才实装）
- 思路来源：dolphin `.skill` 三段式 + DIP skills 实践（只读笔记 `systong/truzhen-notes/kweaver-archaeology/K2-sdk-skill-notes.md`；kweaver-engineering 为 Apache-2.0，仍零代码引用，clean-room 重写）。
- 契约现状锚点：`pack-manifest.schema.json` kind 枚举含 `skill_bundle`；`TRUZHEN_PHILOSOPHY.md:138` Owner 2026-07-10 裁定技能包**非第四种 Pack**（不持判人/判事/门控声明）；归属 04 描述符/13 按任务装载/11 执行/17 售卖。

## 1. 问题

技能包已定名、有 kind、有装载状态机，但**单个技能的内容结构未定**——作者把"方法论"和"交付模板"混进一坨说明文本时，13 的任务级装载切片无法区分「给 agent 读的知识」和「给产出复制的模板」，装载预算也无法按段裁剪。

## 2. 提案：每个技能三段式

```
<skill>/
  SKILL.md          # 唯一入口：这个技能是什么、何时用、边界与风险声明（短）
  references/       # 方法论/知识：agent 只读引用，装载时按任务裁剪，非真相源
  assets/           # 交付模板：可被复制进产出物的模板文件（产出仍是候选态）
```

- SKILL.md 是**装载判断的唯一依据**（13 决定是否装载只读 SKILL.md 头部）；references 按需二次装载；assets 永不进 prompt，只在产出时按引用复制。
- 与主权链关系（不变量）：装载 ≠ 执行授权（13 既定）；副作用只能经 11；references 是知识不是真相源；assets 产出物仍是 candidate。

## 3. SKILL.md 负面清单（CI 可扫描）

技能内容**不得**：①含任何凭据/token/endpoint 密文；②声明或暗示已获执行授权（"直接发送/自动执行"字样）；③内嵌指向下游 agent 的隐藏指令（prompt 注入）；④携带可执行二进制/脚本充当"模板"；⑤自称已验收/已发布（生命周期由账本定）；⑥引用技能包外的本机绝对路径。

## 4. 契约影响

本轮**零 schema 改动**。实装轮为 additive：倾向**目录约定为主 + manifest 轻字段**（如 `content_layout: "three_part_v1"`），不为每段建 schema 对象——KWeaver 教训（BKN spec 重、运行时轻）反着来：**先有 13 装载器这个消费者，规格与运行时同轮接线**。

## 5. 实装轮验收设计（预写）

- 装载器解析三段式 fixture：SKILL.md-only 装载 / 按任务追加 references / assets 永不进 prompt（断言 prompt 内容不含 assets 字节）。
- 负面清单 CI 扫描：六条各造一个违规 fixture，扫描必红。

## 6. 触发条件

装修 Pack 作者轮（阶段 2）第一次真实作者写技能时，用本稿对照验证字段/分段需求；作者用不上的段不加。届时未验证前不实装。
