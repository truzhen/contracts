import os
import unittest


ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
CI_FILE = os.path.join(ROOT, ".github", "workflows", "ci.yml")


class CIContractsGateTest(unittest.TestCase):
    def test_ci_keeps_version_drift_then_runs_aggregate_without_duplication(self):
        with open(CI_FILE, encoding="utf-8") as handle:
            content = handle.read()
        version = content.index("bash scripts/check-version-drift.sh")
        aggregate = content.index("bash scripts/contracts-check.sh --skip-version-drift")
        self.assertLess(version, aggregate)


if __name__ == "__main__":
    unittest.main()
