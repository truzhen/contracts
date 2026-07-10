import os
import subprocess
import unittest


ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
SCRIPT = os.path.join(ROOT, "scripts", "contracts-check.sh")


class ContractsCheckTest(unittest.TestCase):
    def test_aggregate_gate_passes_without_duplicate_version_check(self):
        result = subprocess.run(
            ["bash", SCRIPT, "--skip-version-drift"],
            capture_output=True,
            text=True,
            cwd=ROOT,
        )
        self.assertEqual(result.returncode, 0, msg=result.stdout + result.stderr)
        self.assertIn("跳过 version-drift", result.stdout)
        self.assertGreaterEqual(result.stdout.count("RESULT: PASS"), 2)


if __name__ == "__main__":
    unittest.main()
