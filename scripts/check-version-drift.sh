#!/usr/bin/env bash
# check-version-drift.sh — 版本漂移 gate（B1）
#
# 根本目的：在 truzhen-contracts 建立"schema 变了就必须 bump 版本"的制度守卫。
# 契约仓的版本真相靠 git tag（v0.1.0 / v0.2.0 / v0.3.0 …）+ 仓根 VERSION 文件双真相。
# 本脚本判定：自上一个语义化 tag 以来，是否改了 *.schema.json 却没有 bump VERSION。
# 若改了 schema 而 VERSION 仍等于上次发版的版本号 -> FAIL(exit 1)，要求 bump VERSION。
#
# 兼容 macOS 自带 bash 3.2：不使用关联数组、不使用 mapfile/readarray、不使用 ${var,,}。
#
# 用法：
#   bash scripts/check-version-drift.sh            # 默认对仓根运行
#   bash scripts/check-version-drift.sh <scan-root> # 指定 scan-root（便于隔离测试）
#
# 退出码：
#   0 = PASS（未改 schema，或已 bump VERSION）
#   1 = FAIL（改了 schema 但 VERSION 未 bump），或环境错误（无 VERSION / 无 tag）

set -euo pipefail

# ---- 版本号 bump 建议（best-effort，仅用于打印提示，非严格 semver 解析） ------
# 输入形如 X.Y.Z，输出 bump 后的建议值。非规范输入时原样返回避免脚本崩溃。
bump_minor() {
  _v="$1"
  _major="$(echo "$_v" | cut -d. -f1)"
  _minor="$(echo "$_v" | cut -d. -f2)"
  case "$_major" in ''|*[!0-9]*) echo "$_v"; return;; esac
  case "$_minor" in ''|*[!0-9]*) echo "$_v"; return;; esac
  echo "${_major}.$((_minor + 1)).0"
}

bump_major() {
  _v="$1"
  _major="$(echo "$_v" | cut -d. -f1)"
  case "$_major" in ''|*[!0-9]*) echo "$_v"; return;; esac
  echo "$((_major + 1)).0.0"
}

# ---- 疑似破坏性变更探测（best-effort WARN，不阻断） --------------------------
# 扫描自 lastTag 以来 schema diff 的删除行：
#   - `- "xxx":`   删除了某个属性 -> 疑似 major
#   - required 数组新增（+ 行出现在 required 上下文）-> 疑似 major
# 只提示，不改变调用方退出码。
detect_breaking() {
  _root="$1"
  _lastTag="$2"
  _diff="$(git -C "$_root" diff "$_lastTag"..HEAD -- '*.schema.json' 2>/dev/null || true)"
  [ -z "$_diff" ] && return 0

  # 删除的属性键（形如   - "someKey": ）
  _removedProps="$(printf '%s\n' "$_diff" | grep -E '^-[[:space:]]*"[^"]+"[[:space:]]*:' || true)"
  # 新增到 required 数组的项（形如   + "someKey" ，且不带冒号，粗略近似）
  _addedRequired="$(printf '%s\n' "$_diff" | grep -E '^\+[[:space:]]*"[^"]+"[[:space:]]*,?[[:space:]]*$' || true)"

  if [ -n "$_removedProps" ] || [ -n "$_addedRequired" ]; then
    echo "" >&2
    echo "[version-drift][WARN] 疑似破坏性变更（best-effort 提示，需人工确认，可能是 major bump）：" >&2
    if [ -n "$_removedProps" ]; then
      echo "  删除/改动的属性行：" >&2
      printf '%s\n' "$_removedProps" | sed 's/^/    /' >&2
    fi
    if [ -n "$_addedRequired" ]; then
      echo "  新增的字符串项（疑似 required / enum 追加，请确认是否破坏兼容）：" >&2
      printf '%s\n' "$_addedRequired" | sed 's/^/    /' >&2
    fi
    echo "  破坏性变更（删字段 / 改必填 / 改类型 / 改 enum 语义）必须 major bump 并回 Owner。" >&2
  fi
  return 0
}

# ---- 定位 scan-root ----------------------------------------------------------
# 优先用第一个参数；否则用脚本所在目录的上一级（scripts/ 的父目录 = 仓根）。
if [ "$#" -ge 1 ] && [ -n "${1:-}" ]; then
  ROOT="$1"
else
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
  ROOT="$(cd "$SCRIPT_DIR/.." >/dev/null 2>&1 && pwd)"
fi

if ! git -C "$ROOT" rev-parse --git-dir >/dev/null 2>&1; then
  echo "[version-drift][ERROR] scan-root 不是 git 仓：$ROOT" >&2
  exit 1
fi

# ---- 读取 VERSION 文件 -------------------------------------------------------
if [ ! -f "$ROOT/VERSION" ]; then
  echo "[version-drift][FAIL] 仓根缺少 VERSION 文件：$ROOT/VERSION" >&2
  echo "  VERSION 文件是显式版本真相源，必须与最新语义化 tag 对齐。" >&2
  exit 1
fi
# 去掉首尾空白与可能的 'v' 前缀
curVer="$(tr -d '[:space:]' < "$ROOT/VERSION")"
curVer="${curVer#v}"
if [ -z "$curVer" ]; then
  echo "[version-drift][FAIL] VERSION 文件为空：$ROOT/VERSION" >&2
  exit 1
fi

# ---- 找最新语义化 tag --------------------------------------------------------
lastTag="$(git -C "$ROOT" tag -l 'v*' --sort=-version:refname | head -1)"
if [ -z "$lastTag" ]; then
  echo "[version-drift][FAIL] 找不到任何 v* 语义化 tag，无法判定漂移。" >&2
  echo "  首个发布需先打 tag（如 v0.1.0）并写入 VERSION 文件。" >&2
  exit 1
fi
lastVer="${lastTag#v}"

echo "[version-drift] 最新 tag = ${lastTag}（版本 ${lastVer}）"
echo "[version-drift] VERSION 文件 = ${curVer}"

# ---- 计算自上次 tag 以来变化的 schema ----------------------------------------
# 发布 gate 既要看 tag..HEAD，也要看当前工作区：本地验收若忽略暂存、未暂存或
# 未跟踪 schema，会错误宣称“无需 bump”。pathspec 会匹配任意层级（含 spines/）。
changedSchemas="$({
  git -C "$ROOT" diff --name-only "$lastTag"..HEAD -- '*.schema.json'
  git -C "$ROOT" diff --name-only -- '*.schema.json'
  git -C "$ROOT" diff --cached --name-only -- '*.schema.json'
  git -C "$ROOT" ls-files --others --exclude-standard -- '*.schema.json'
} | sed '/^$/d' | sort -u)"

if [ -z "$changedSchemas" ]; then
  echo "[version-drift][PASS] 自 $lastTag 起未修改任何 *.schema.json，无需 bump。"
  exit 0
fi

echo "[version-drift] 自 $lastTag 起发生变化的 schema："
printf '%s\n' "$changedSchemas" | sed 's/^/    - /'

# ---- 核心判定：改了 schema 但 VERSION 未 bump ---------------------------------
if [ "$curVer" = "$lastVer" ]; then
  echo "" >&2
  echo "[version-drift][FAIL] 检测到版本漂移：自 $lastTag 起修改了上述 schema，" >&2
  echo "  但 VERSION 仍为 ${curVer}（== 上次发版版本 ${lastVer}），未 bump。" >&2
  echo "  修复：按 SemVer bump 仓根 VERSION 文件（并在发布时打对应 v* tag）。" >&2
  echo "    新增可选字段 -> minor（如 $(bump_minor "$curVer")）" >&2
  echo "    删字段 / 改必填 / 改类型 / 改 enum 语义 -> major（如 $(bump_major "$curVer")）" >&2

  detect_breaking "$ROOT" "$lastTag"

  exit 1
fi

echo "[version-drift][PASS] schema 有变化，且 VERSION 已从 ${lastVer} bump 到 ${curVer}，放行。"

# 即便已 bump，也给出破坏性提示（仅 WARN，不改变 PASS 结果）。
detect_breaking "$ROOT" "$lastTag"

exit 0
