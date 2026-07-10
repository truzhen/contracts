import json
import os
import subprocess
import tempfile
import unittest


ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
TOOL = os.path.join(ROOT, "scripts", "check_go_schema_consistency.go")
FIXTURES = os.path.join(ROOT, "scripts", "tests", "fixtures", "go-schema")


def write_map(pairs):
    handle = tempfile.NamedTemporaryFile("w", suffix=".json", delete=False, encoding="utf-8")
    json.dump({"pairs": pairs}, handle)
    handle.close()
    return handle.name


def run_map(path):
    return subprocess.run(
        [GoSchemaConsistencyTest.binary, "--map", path, "--repo-root", ROOT],
        capture_output=True,
        text=True,
        cwd=ROOT,
    )


class GoSchemaConsistencyTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.binary_dir = tempfile.TemporaryDirectory()
        cls.binary = os.path.join(cls.binary_dir.name, "check-go-schema-consistency")
        result = subprocess.run(
            ["go", "build", "-o", cls.binary, TOOL],
            capture_output=True,
            text=True,
            cwd=ROOT,
        )
        if result.returncode:
            raise RuntimeError(result.stdout + result.stderr)

    @classmethod
    def tearDownClass(cls):
        cls.binary_dir.cleanup()

    def test_real_pairs_pass(self):
        result = run_map(os.path.join(ROOT, "scripts", "go-schema-map.json"))
        self.assertEqual(result.returncode, 0, msg=result.stdout + result.stderr)
        self.assertIn("mapped_pairs=4", result.stdout)
        self.assertIn("passed_pairs=4", result.stdout)

    def test_schema_only_field_fails(self):
        path = write_map([{
            "go_dir": "candidates",
            "go_type": "CandidateEnvelope",
            "schema": "scripts/tests/fixtures/go-schema/candidate-envelope-drift.json",
        }])
        self.addCleanup(os.unlink, path)
        result = run_map(path)
        self.assertEqual(result.returncode, 1, msg=result.stdout + result.stderr)
        self.assertIn("bogus_field", result.stdout)

    def test_named_string_alias_kind_fails(self):
        path = write_map([{
            "go_dir": "market",
            "go_type": "ProviderRequirement",
            "schema": "scripts/tests/fixtures/go-schema/provider-requirement-kind-drift.json",
        }])
        self.addCleanup(os.unlink, path)
        result = run_map(path)
        self.assertEqual(result.returncode, 1, msg=result.stdout + result.stderr)
        self.assertIn("gateway_class", result.stdout)

    def test_missing_schema_is_tool_error(self):
        path = write_map([{
            "go_dir": "candidates",
            "go_type": "CandidateEnvelope",
            "schema": "scripts/tests/fixtures/go-schema/missing.json",
        }])
        self.addCleanup(os.unlink, path)
        result = run_map(path)
        self.assertEqual(result.returncode, 2, msg=result.stdout + result.stderr)

    def test_known_gaps_are_not_an_allowed_exemption(self):
        path = write_map([{
            "go_dir": "candidates",
            "go_type": "CandidateEnvelope",
            "schema": "candidate-envelope.schema.json",
            "known_gaps": ["payload"],
        }])
        self.addCleanup(os.unlink, path)
        result = run_map(path)
        self.assertEqual(result.returncode, 2, msg=result.stdout + result.stderr)
        self.assertIn("known_gaps", result.stdout + result.stderr)

    def test_composed_property_is_not_silently_treated_as_any(self):
        with open(os.path.join(ROOT, "candidate-envelope.schema.json"), encoding="utf-8") as handle:
            schema = json.load(handle)
        schema["properties"]["payload"] = {"oneOf": [{"type": "string"}, {"type": "object"}]}
        handle = tempfile.NamedTemporaryFile("w", suffix=".json", delete=False, encoding="utf-8")
        json.dump(schema, handle)
        handle.close()
        self.addCleanup(os.unlink, handle.name)
        path = write_map([{
            "go_dir": "candidates",
            "go_type": "CandidateEnvelope",
            "schema": handle.name,
        }])
        self.addCleanup(os.unlink, path)
        result = run_map(path)
        self.assertEqual(result.returncode, 2, msg=result.stdout + result.stderr)
        self.assertIn("oneOf", result.stdout + result.stderr)


if __name__ == "__main__":
    unittest.main()
