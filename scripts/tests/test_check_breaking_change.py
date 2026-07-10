import json
import os
import subprocess
import sys
import tempfile
import unittest


ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
SCRIPT = os.path.join(ROOT, "scripts", "check-breaking-change.py")
FIXTURES = os.path.join(ROOT, "scripts", "tests", "fixtures", "breaking-change")


def run_dirs(old_dir, new_dir):
    return subprocess.run(
        [sys.executable, SCRIPT, "--old-dir", old_dir, "--new-dir", new_dir],
        capture_output=True,
        text=True,
    )


class BreakingChangeTest(unittest.TestCase):
    def assert_case(self, case, expected_code, expected_text):
        result = run_dirs(
            os.path.join(FIXTURES, case, "old"),
            os.path.join(FIXTURES, case, "new"),
        )
        output = result.stdout + result.stderr
        self.assertEqual(result.returncode, expected_code, msg=output)
        self.assertIn(expected_text, output)

    def test_compatible_addition_passes(self):
        self.assert_case("compatible-add", 0, "PASS")

    def test_deleted_property_fails(self):
        self.assert_case("delete-field", 1, "properties/amount")

    def test_added_required_fails(self):
        self.assert_case("add-required", 1, "required/status")

    def test_changed_type_fails(self):
        self.assert_case("change-type", 1, "properties/amount/type")

    def test_widened_type_passes(self):
        with tempfile.TemporaryDirectory() as directory:
            old_dir = os.path.join(directory, "old")
            new_dir = os.path.join(directory, "new")
            os.makedirs(old_dir)
            os.makedirs(new_dir)
            with open(os.path.join(old_dir, "x.json"), "w", encoding="utf-8") as handle:
                json.dump({"type": "object", "properties": {"value": {"type": "string"}}}, handle)
            with open(os.path.join(new_dir, "x.json"), "w", encoding="utf-8") as handle:
                json.dump({"type": "object", "properties": {"value": {"type": ["string", "null"]}}}, handle)
            result = run_dirs(old_dir, new_dir)
            self.assertEqual(result.returncode, 0, msg=result.stdout + result.stderr)

    def test_major_migration_requires_explicit_owner_approval(self):
        with tempfile.TemporaryDirectory() as root:
            def git(*args):
                result = subprocess.run(["git", *args], cwd=root, capture_output=True, text=True)
                self.assertEqual(result.returncode, 0, msg=result.stdout + result.stderr)

            git("init", "-q")
            git("config", "user.email", "tests@example.invalid")
            git("config", "user.name", "contracts-check test")
            with open(os.path.join(root, "x.schema.json"), "w", encoding="utf-8") as handle:
                json.dump({"type": "object", "properties": {"value": {"type": "string"}}}, handle)
            with open(os.path.join(root, "VERSION"), "w", encoding="utf-8") as handle:
                handle.write("1.0.0\n")
            git("add", ".")
            git("commit", "-qm", "baseline")
            git("tag", "v1.0.0")
            with open(os.path.join(root, "x.schema.json"), "w", encoding="utf-8") as handle:
                json.dump({"type": "object", "properties": {"value": {"type": "integer"}}}, handle)
            with open(os.path.join(root, "VERSION"), "w", encoding="utf-8") as handle:
                handle.write("2.0.0\n")
            os.makedirs(os.path.join(root, "MIGRATIONS"))
            migration = os.path.join(root, "MIGRATIONS", "2.0.0.md")
            with open(migration, "w", encoding="utf-8") as handle:
                handle.write("# 2.0.0 migration\n")

            result = subprocess.run([sys.executable, SCRIPT, "--repo-root", root], capture_output=True, text=True)
            self.assertEqual(result.returncode, 1, msg=result.stdout + result.stderr)
            self.assertIn("Owner-Approval", result.stdout + result.stderr)

            with open(migration, "a", encoding="utf-8") as handle:
                handle.write("\nOwner-Approval: APPROVED\n")
            result = subprocess.run([sys.executable, SCRIPT, "--repo-root", root], capture_output=True, text=True)
            self.assertEqual(result.returncode, 0, msg=result.stdout + result.stderr)

    def test_removed_enum_value_fails(self):
        self.assert_case("remove-enum-value", 1, "enum")

    def test_changed_const_fails(self):
        with tempfile.TemporaryDirectory() as directory:
            old_dir = os.path.join(directory, "old")
            new_dir = os.path.join(directory, "new")
            os.makedirs(old_dir)
            os.makedirs(new_dir)
            with open(os.path.join(old_dir, "x.json"), "w", encoding="utf-8") as handle:
                json.dump({"type": "object", "properties": {"state": {"const": "draft"}}}, handle)
            with open(os.path.join(new_dir, "x.json"), "w", encoding="utf-8") as handle:
                json.dump({"type": "object", "properties": {"state": {"const": "formal"}}}, handle)
            result = run_dirs(old_dir, new_dir)
            self.assertEqual(result.returncode, 1, msg=result.stdout + result.stderr)
            self.assertIn("const", result.stdout + result.stderr)

    def test_tightened_length_constraint_fails(self):
        with tempfile.TemporaryDirectory() as directory:
            old_dir = os.path.join(directory, "old")
            new_dir = os.path.join(directory, "new")
            os.makedirs(old_dir)
            os.makedirs(new_dir)
            with open(os.path.join(old_dir, "x.json"), "w", encoding="utf-8") as handle:
                json.dump({"type": "object", "properties": {"name": {"type": "string"}}}, handle)
            with open(os.path.join(new_dir, "x.json"), "w", encoding="utf-8") as handle:
                json.dump({"type": "object", "properties": {"name": {"type": "string", "minLength": 2}}}, handle)
            result = run_dirs(old_dir, new_dir)
            self.assertEqual(result.returncode, 1, msg=result.stdout + result.stderr)
            self.assertIn("minLength", result.stdout + result.stderr)

    def test_tightened_additional_properties_fails(self):
        self.assert_case("tighten-additional-properties", 1, "additionalProperties")

    def test_changed_composition_fails_closed(self):
        self.assert_case("unsupported-composition", 2, "anyOf")

    def test_external_ref_fails_closed(self):
        with tempfile.TemporaryDirectory() as directory:
            old_dir = os.path.join(directory, "old")
            new_dir = os.path.join(directory, "new")
            os.makedirs(old_dir)
            os.makedirs(new_dir)
            with open(os.path.join(old_dir, "x.json"), "w", encoding="utf-8") as handle:
                json.dump({"type": "object", "properties": {"value": {"type": "string"}}}, handle)
            with open(os.path.join(new_dir, "x.json"), "w", encoding="utf-8") as handle:
                json.dump({"type": "object", "properties": {"value": {"$ref": "https://example.invalid/value"}}}, handle)
            result = run_dirs(old_dir, new_dir)
            self.assertEqual(result.returncode, 2, msg=result.stdout + result.stderr)
            self.assertIn("$ref", result.stdout + result.stderr)

    def test_invalid_json_is_tool_error(self):
        with tempfile.TemporaryDirectory() as directory:
            old_dir = os.path.join(directory, "old")
            new_dir = os.path.join(directory, "new")
            os.makedirs(old_dir)
            os.makedirs(new_dir)
            with open(os.path.join(old_dir, "x.json"), "w", encoding="utf-8") as handle:
                handle.write("{invalid json")
            with open(os.path.join(new_dir, "x.json"), "w", encoding="utf-8") as handle:
                handle.write("{}")
            result = run_dirs(old_dir, new_dir)
            self.assertEqual(result.returncode, 2, msg=result.stdout + result.stderr)
            self.assertIn("JSON", result.stdout + result.stderr)

    def test_repository_mode_passes_current_baseline(self):
        result = subprocess.run(
            [sys.executable, SCRIPT, "--repo-root", ROOT],
            capture_output=True,
            text=True,
        )
        self.assertEqual(result.returncode, 0, msg=result.stdout + result.stderr)
        self.assertIn("RESULT: PASS", result.stdout + result.stderr)


if __name__ == "__main__":
    unittest.main()
