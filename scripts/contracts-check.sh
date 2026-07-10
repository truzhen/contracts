#!/usr/bin/env bash
# 聚合 contracts 的 schema 版本、破坏性与 Go 对照门禁。
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
skip_version_drift=false

if [[ "${1:-}" == "--skip-version-drift" ]]; then
  skip_version_drift=true
  shift
fi
if [[ $# -ne 0 ]]; then
  echo "用法：bash scripts/contracts-check.sh [--skip-version-drift]" >&2
  exit 2
fi

if [[ "$skip_version_drift" == true ]]; then
  echo "跳过 version-drift（由调用方已执行）"
else
  bash "$root/scripts/check-version-drift.sh"
fi

python3 "$root/scripts/check-breaking-change.py" --repo-root "$root"

checker_bin="$(mktemp "${TMPDIR:-/tmp}/truzhen-contracts-go-schema.XXXXXX")"
trap 'rm -f "$checker_bin"' EXIT
go build -o "$checker_bin" "$root/scripts/check_go_schema_consistency.go"
"$checker_bin" --repo-root "$root" --map "$root/scripts/go-schema-map.json"

python3 -m unittest discover -s "$root/scripts/tests" -p 'test_check_breaking_change.py'
python3 -m unittest discover -s "$root/scripts/tests" -p 'test_go_schema_consistency.py'

echo "contracts-check: PASS"
