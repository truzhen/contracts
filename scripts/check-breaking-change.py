#!/usr/bin/env python3
"""检测 JSON Schema 的破坏性变更。

仓库模式将最新 v* tag 中的真实 *.schema.json 与工作树比较；fixture
模式比较两个目录中的普通 *.json。退出码：0=兼容，1=破坏，2=无法判定。
"""

import argparse
import json
import os
import subprocess
import sys


COMPOSITION_KEYS = {
    "anyOf",
    "oneOf",
    "allOf",
    "not",
    "if",
    "then",
    "else",
    "patternProperties",
    "dependentRequired",
    "dependentSchemas",
}
TIGHTEN_UP = ("minimum", "exclusiveMinimum", "minLength", "minItems", "minProperties")
TIGHTEN_DOWN = ("maximum", "exclusiveMaximum", "maxLength", "maxItems", "maxProperties")
EXCLUDED_REPO_PREFIXES = (".git/", "scripts/tests/")


class Findings:
    def __init__(self):
        self.breaking = []
        self.notes = []
        self.unsupported = []

    def add(self, bucket, filename, pointer, message):
        getattr(self, bucket).append((filename, pointer, message))


def pointer_join(pointer, token):
    escaped = token.replace("~", "~0").replace("/", "~1")
    return "%s/%s" % (pointer, escaped)


def json_value_set(values):
    return {json.dumps(value, sort_keys=True, separators=(",", ":")) for value in values}


def json_type_set(node, filename, pointer, findings):
    if "type" not in node:
        return None
    value = node["type"]
    if isinstance(value, str):
        return {value}
    if isinstance(value, list) and value and all(isinstance(item, str) for item in value):
        return set(value)
    findings.add("unsupported", filename, pointer_join(pointer, "type"), "type 不是 string 或非空 string 数组")
    return False


def is_number(value):
    return isinstance(value, (int, float)) and not isinstance(value, bool)


def local_ref(node, root, filename, pointer, findings):
    if not isinstance(node, dict) or "$ref" not in node:
        return node
    reference = node["$ref"]
    if not isinstance(reference, str) or not reference.startswith("#/"):
        findings.add("unsupported", filename, pointer, "外部或无效 $ref: %r" % reference)
        return None
    current = root
    for part in reference[2:].split("/"):
        key = part.replace("~1", "/").replace("~0", "~")
        if not isinstance(current, dict) or key not in current:
            findings.add("unsupported", filename, pointer, "$ref 目标不存在: %s" % reference)
            return None
        current = current[key]
    if not isinstance(current, dict):
        findings.add("unsupported", filename, pointer, "$ref 目标不是 schema object: %s" % reference)
        return None
    return current


def has_changed_composition(old, new):
    if not isinstance(old, dict) or not isinstance(new, dict):
        return None
    for key in COMPOSITION_KEYS:
        if key in old or key in new:
            return key
    return None


def compare(old, new, old_root, new_root, filename, pointer, findings, depth=0):
    if old == new:
        return
    if depth > 64:
        findings.add("unsupported", filename, pointer, "schema 嵌套超过 64 层")
        return

    composition = has_changed_composition(old, new)
    if composition:
        findings.add(
            "unsupported",
            filename,
            pointer_join(pointer, composition),
            "变化子树或祖先使用组合关键字，本工具不判定",
        )
        return

    old = local_ref(old, old_root, filename, pointer, findings)
    new = local_ref(new, new_root, filename, pointer, findings)
    if old is None or new is None:
        return
    if not isinstance(old, dict) or not isinstance(new, dict):
        findings.add("unsupported", filename, pointer, "schema 节点不是 object")
        return

    old_types = json_type_set(old, filename, pointer, findings)
    new_types = json_type_set(new, filename, pointer, findings)
    if old_types is not False and new_types is not False and old_types != new_types:
        if new_types is None:
            findings.add("notes", filename, pointer_join(pointer, "type"), "type 移除（放宽）")
        elif old_types is None:
            findings.add("breaking", filename, pointer_join(pointer, "type"), "新增 type 限制: %r" % sorted(new_types))
        else:
            removed_types = old_types - new_types
            if removed_types:
                findings.add("breaking", filename, pointer_join(pointer, "type"), "type 删值（收紧）: %s" % sorted(removed_types))
            else:
                findings.add("notes", filename, pointer_join(pointer, "type"), "type 加值（兼容）")

    old_enum, new_enum = old.get("enum"), new.get("enum")
    if old_enum != new_enum:
        if new_enum is None:
            findings.add("notes", filename, pointer_join(pointer, "enum"), "enum 移除（放宽）")
        elif old_enum is None:
            findings.add("breaking", filename, pointer_join(pointer, "enum"), "新增 enum 限制")
        else:
            removed = json_value_set(old_enum) - json_value_set(new_enum)
            if removed:
                findings.add("breaking", filename, pointer_join(pointer, "enum"), "enum 删值: %s" % sorted(removed))
            else:
                findings.add("notes", filename, pointer_join(pointer, "enum"), "enum 加值（兼容）")

    old_const, new_const = old.get("const"), new.get("const")
    if old_const != new_const:
        if new_const is None:
            findings.add("notes", filename, pointer_join(pointer, "const"), "const 移除（放宽）")
        else:
            findings.add("breaking", filename, pointer_join(pointer, "const"), "const 新增或改值")

    old_ap, new_ap = old.get("additionalProperties", True), new.get("additionalProperties", True)
    if old_ap != new_ap:
        if new_ap is False and old_ap is not False:
            findings.add("breaking", filename, pointer_join(pointer, "additionalProperties"), "additionalProperties 收紧为 false")
        elif old_ap is False and new_ap is not False:
            findings.add("notes", filename, pointer_join(pointer, "additionalProperties"), "additionalProperties 放宽")
        else:
            findings.add("unsupported", filename, pointer_join(pointer, "additionalProperties"), "复杂 additionalProperties 变更")

    for key in TIGHTEN_UP:
        old_value, new_value = old.get(key), new.get(key)
        if new_value is None:
            continue
        if not is_number(new_value) or (old_value is not None and not is_number(old_value)):
            findings.add("unsupported", filename, pointer_join(pointer, key), "非数值约束")
        elif old_value is None or new_value > old_value:
            findings.add("breaking", filename, pointer_join(pointer, key), "约束收紧: %r -> %r" % (old_value, new_value))

    for key in TIGHTEN_DOWN:
        old_value, new_value = old.get(key), new.get(key)
        if new_value is None:
            continue
        if not is_number(new_value) or (old_value is not None and not is_number(old_value)):
            findings.add("unsupported", filename, pointer_join(pointer, key), "非数值约束")
        elif old_value is None or new_value < old_value:
            findings.add("breaking", filename, pointer_join(pointer, key), "约束收紧: %r -> %r" % (old_value, new_value))

    if old.get("uniqueItems") is not True and new.get("uniqueItems") is True:
        findings.add("breaking", filename, pointer_join(pointer, "uniqueItems"), "uniqueItems 收紧为 true")

    old_required = set(old.get("required") or [])
    new_required = set(new.get("required") or [])
    for field in sorted(new_required - old_required):
        findings.add("breaking", filename, pointer_join(pointer_join(pointer, "required"), field), "required 新增")
    for field in sorted(old_required - new_required):
        findings.add("notes", filename, pointer_join(pointer_join(pointer, "required"), field), "required 移除（放宽）")

    compare_object_members(
        old.get("properties") or {}, new.get("properties") or {}, old_root, new_root,
        filename, pointer_join(pointer, "properties"), findings, "property", depth,
    )
    compare_object_members(
        old.get("$defs") or {}, new.get("$defs") or {}, old_root, new_root,
        filename, pointer_join(pointer, "$defs"), findings, "definition", depth,
    )

    old_items, new_items = old.get("items"), new.get("items")
    if old_items != new_items:
        if isinstance(old_items, dict) and isinstance(new_items, dict):
            compare(old_items, new_items, old_root, new_root, filename, pointer_join(pointer, "items"), findings, depth + 1)
        elif old_items is None and isinstance(new_items, dict):
            findings.add("breaking", filename, pointer_join(pointer, "items"), "items 新增限制")
        elif isinstance(old_items, dict) and new_items is None:
            findings.add("notes", filename, pointer_join(pointer, "items"), "items 移除（放宽）")
        else:
            findings.add("unsupported", filename, pointer_join(pointer, "items"), "items 形状变更")


def compare_object_members(old_members, new_members, old_root, new_root, filename, pointer, findings, noun, depth):
    if not isinstance(old_members, dict) or not isinstance(new_members, dict):
        findings.add("unsupported", filename, pointer, "%s 集合不是 object" % noun)
        return
    for name in sorted(set(old_members) - set(new_members)):
        findings.add("breaking", filename, pointer_join(pointer, name), "%s 删除" % noun)
    for name in sorted(set(new_members) - set(old_members)):
        findings.add("notes", filename, pointer_join(pointer, name), "新增可选 %s（兼容）" % noun)
    for name in sorted(set(old_members) & set(new_members)):
        compare(old_members[name], new_members[name], old_root, new_root, filename, pointer_join(pointer, name), findings, depth + 1)


def load_json(content, filename, findings):
    try:
        value = json.loads(content)
    except (TypeError, json.JSONDecodeError) as error:
        findings.add("unsupported", filename, "#", "JSON 解析失败: %s" % error)
        return None
    if not isinstance(value, dict):
        findings.add("unsupported", filename, "#", "根节点不是 object")
        return None
    return value


def list_files(root, suffix, exclude_repo_dirs=False):
    results = []
    for directory, dirnames, filenames in os.walk(root):
        relative = os.path.relpath(directory, root).replace(os.sep, "/")
        if relative == ".git" or relative.startswith(".git/"):
            dirnames[:] = []
            continue
        if exclude_repo_dirs and any(relative == prefix[:-1] or relative.startswith(prefix) for prefix in EXCLUDED_REPO_PREFIXES):
            dirnames[:] = []
            continue
        for filename in filenames:
            if filename.endswith(suffix):
                results.append(filename if relative == "." else os.path.join(relative, filename))
    return sorted(results)


def git(root, *args):
    return subprocess.run(["git", *args], cwd=root, capture_output=True, text=True)


def repository_pairs(root, findings):
    tags = [tag for tag in git(root, "tag", "-l", "v*", "--sort=-version:refname").stdout.splitlines() if tag]
    if not tags:
        findings.add("unsupported", "repository", "#", "找不到 v* tag")
        return [], None
    base_tag = tags[0]
    base_files = [
        path for path in git(root, "ls-tree", "-r", "--name-only", base_tag).stdout.splitlines()
        if path.endswith(".schema.json") and not any(path.startswith(prefix) for prefix in EXCLUDED_REPO_PREFIXES)
    ]
    current_files = list_files(root, ".schema.json", exclude_repo_dirs=True)
    current_set = {path.replace(os.sep, "/") for path in current_files}
    pairs = []
    for path in sorted(base_files):
        if path not in current_set:
            findings.add("breaking", path, "#", "schema 文件删除")
            continue
        old_result = git(root, "show", "%s:%s" % (base_tag, path))
        if old_result.returncode != 0:
            findings.add("unsupported", path, "#", "无法读取 tag 基线")
            continue
        with open(os.path.join(root, path), encoding="utf-8") as handle:
            pairs.append((path, old_result.stdout, handle.read()))
    for path in sorted(current_set - set(base_files)):
        findings.add("notes", path, "#", "新增 schema（仍须 Owner 与 SemVer 裁定）")
    return pairs, base_tag


def fixture_pairs(old_dir, new_dir, findings):
    old_files = list_files(old_dir, ".json")
    new_files = list_files(new_dir, ".json")
    pairs = []
    for path in sorted(set(old_files) - set(new_files)):
        findings.add("breaking", path, "#", "fixture 文件删除")
    for path in sorted(set(new_files) - set(old_files)):
        findings.add("notes", path, "#", "fixture 文件新增")
    for path in sorted(set(old_files) & set(new_files)):
        with open(os.path.join(old_dir, path), encoding="utf-8") as old_handle:
            old_content = old_handle.read()
        with open(os.path.join(new_dir, path), encoding="utf-8") as new_handle:
            new_content = new_handle.read()
        pairs.append((path, old_content, new_content))
    return pairs


def major_migration_allows(root, base_tag):
    try:
        base_major = int(base_tag.removeprefix("v").split(".")[0])
        with open(os.path.join(root, "VERSION"), encoding="utf-8") as handle:
            current_version = handle.read().strip().removeprefix("v")
        current_major = int(current_version.split(".")[0])
    except (OSError, ValueError, IndexError):
        return False, "VERSION 或 tag 不是有效 SemVer"
    migration = os.path.join(root, "MIGRATIONS", current_version + ".md")
    if current_major <= base_major or not os.path.isfile(migration):
        return False, "未同时满足主版本升级与 %s" % migration
    try:
        with open(migration, encoding="utf-8") as handle:
            approved = any(line.strip() == "Owner-Approval: APPROVED" for line in handle)
    except OSError as error:
        return False, "无法读取迁移说明 %s: %s" % (migration, error)
    if not approved:
        return False, "%s 缺少精确的 Owner-Approval: APPROVED 记录" % migration
    return True, migration


def print_findings(findings):
    for category, rows in (("NOTE", findings.notes), ("BREAKING", findings.breaking), ("UNSUPPORTED", findings.unsupported)):
        for filename, pointer, message in rows:
            print("[%s] %s %s — %s" % (category, filename, pointer, message))


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--repo-root", default=".")
    parser.add_argument("--old-dir")
    parser.add_argument("--new-dir")
    args = parser.parse_args()
    if bool(args.old_dir) != bool(args.new_dir):
        print("[UNSUPPORTED] arguments # — --old-dir 与 --new-dir 必须同时提供")
        return 2

    findings = Findings()
    base_tag = None
    if args.old_dir:
        pairs = fixture_pairs(os.path.abspath(args.old_dir), os.path.abspath(args.new_dir), findings)
    else:
        root = os.path.abspath(args.repo_root)
        pairs, base_tag = repository_pairs(root, findings)

    for filename, old_content, new_content in pairs:
        old = load_json(old_content, filename, findings)
        new = load_json(new_content, filename, findings)
        if old is not None and new is not None:
            compare(old, new, old, new, filename, "#", findings)

    print_findings(findings)
    if findings.unsupported:
        print("RESULT: UNSUPPORTED")
        return 2
    if findings.breaking:
        if base_tag:
            allowed, detail = major_migration_allows(os.path.abspath(args.repo_root), base_tag)
            if allowed:
                print("RESULT: BREAKING documented by major migration: %s" % detail)
                return 0
            print("[NOTE] %s" % detail)
        print("RESULT: BREAKING")
        return 1
    print("RESULT: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
