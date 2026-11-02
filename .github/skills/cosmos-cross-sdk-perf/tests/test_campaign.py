# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

import csv
import importlib.util
import json
import pathlib
import sys
import tarfile
import tempfile
import unittest
from unittest import mock


ROOT = pathlib.Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location("campaign", ROOT / "scripts" / "campaign.py")
campaign = importlib.util.module_from_spec(SPEC)
assert SPEC.loader is not None
sys.modules[SPEC.name] = campaign
SPEC.loader.exec_module(campaign)


def load_default():
    return json.loads((ROOT / "config" / "default.json").read_text(encoding="utf-8"))


class CampaignTests(unittest.TestCase):
    def test_default_config_has_governed_workload(self):
        config = load_default()
        config["azure"]["subscriptionId"] = "00000000-0000-0000-0000-000000000000"
        config["azure"]["resourceGroup"] = "perf-rg"
        config["azure"]["existingCosmosAccountResourceId"] = (
            "/subscriptions/00000000-0000-0000-0000-000000000000/"
            "resourceGroups/cosmos-rg/providers/Microsoft.DocumentDB/databaseAccounts/account"
        )
        campaign.validate_config(config)

    def test_config_rejects_unmatched_cycle(self):
        config = load_default()
        config["workload"]["cycle"] = ["read", "upsert"]
        with self.assertRaisesRegex(campaign.CampaignError, "workload.cycle"):
            campaign.validate_config(config, allow_incomplete_azure=True)

    def test_classification_maps_every_row_and_totals_100(self):
        rows = [
            campaign.PerfRow("runtime.mallocgc", "go-v2", 25.0),
            campaign.PerfRow("tokio::runtime::park", "rust", 25.0),
            campaign.PerfRow("tcp_sendmsg", "[kernel.kallsyms]", 25.0),
            campaign.PerfRow("custom_resolved_symbol", "binary", 25.0),
        ]
        totals, unmapped = campaign.classify(rows)
        self.assertEqual([], unmapped)
        self.assertAlmostEqual(100.0, sum(totals.values()))
        self.assertEqual(25.0, totals["go-runtime-sdk"])
        self.assertEqual(25.0, totals["tokio"])

    @mock.patch.object(campaign, "az_json")
    @mock.patch.object(campaign, "run")
    def test_dry_run_plan_verifies_region_and_four_vcpu_sku(self, run_mock, az_mock):
        run_mock.side_effect = ["", "false"]
        config = load_default()
        config["azure"].update(
            {
                "subscriptionId": "sub",
                "resourceGroup": "perf-rg",
                "existingCosmosAccountResourceId": (
                    "/subscriptions/sub/resourceGroups/cosmos-rg/providers/"
                    "Microsoft.DocumentDB/databaseAccounts/account"
                ),
            }
        )
        az_mock.side_effect = [
            {
                "location": "westus2",
                "locations": [{"locationName": "West US 2", "failoverPriority": 0}],
            },
            [
                {
                    "name": "Standard_D4s_v5",
                    "restrictions": [],
                    "capabilities": [{"name": "vCPUs", "value": "4"}],
                }
            ],
        ]
        plan = campaign.provision_plan(config)
        self.assertEqual("West US 2", plan["location"])
        self.assertEqual("Standard_D4s_v5", plan["vmSku"])
        self.assertEqual(
            [
                mock.call(["az", "account", "set", "--subscription", "sub"]),
                mock.call(
                    [
                        "az",
                        "group",
                        "exists",
                        "--name",
                        "perf-rg",
                        "--output",
                        "tsv",
                    ]
                ),
            ],
            run_mock.call_args_list,
        )

    def test_validate_leg_and_report(self):
        config = load_default()
        with tempfile.TemporaryDirectory() as temp:
            run_dir = pathlib.Path(temp)
            for leg in campaign.LEGS:
                self._write_leg(run_dir / leg, leg, config)
            campaign.report(run_dir, config)
            comparison = json.loads(
                (run_dir / "comparison" / "comparison.json").read_text(encoding="utf-8")
            )
            self.assertEqual(3, len(comparison["legs"]))
            self.assertTrue((run_dir / "comparison" / "cpu-flywheel.svg").is_file())
            self.assertTrue((run_dir / "SHA256SUMS").is_file())

    def test_sanitize_redacts_sensitive_json_fields(self):
        with tempfile.TemporaryDirectory() as temp:
            run_dir = pathlib.Path(temp) / "run"
            run_dir.mkdir()
            (run_dir / "metrics.json").write_text(
                json.dumps({"endpoint": "https://account", "success": 10}), encoding="utf-8"
            )
            output = pathlib.Path(temp) / "bundle.tar.gz"
            campaign.sanitize(run_dir, output)
            with tarfile.open(output, "r:gz") as archive:
                member = archive.extractfile("cosmos-perf-result/metrics.json")
                assert member is not None
                value = json.load(member)
            self.assertTrue(value["endpoint"].startswith("sha256:"))
            self.assertEqual(10, value["success"])

    def _write_leg(self, leg_dir, leg, config):
        (leg_dir / "raw").mkdir(parents=True)
        binary = b"test-binary-" + leg.encode()
        (leg_dir / "raw" / "binary").write_bytes(binary)
        for path in (
            "perf.data",
            "proc-status.txt",
            "smaps-rollup.txt",
            "smaps.txt",
            "cgroup-memory.txt",
            "thread-stacks.txt",
        ):
            (leg_dir / "raw" / path).write_text("evidence\n", encoding="utf-8")
        if leg.startswith("go-"):
            (leg_dir / "raw" / "cpu.pprof").write_bytes(b"cpu")
            (leg_dir / "raw" / "heap.pprof").write_bytes(b"heap")
        source_commit = {
            "go-v1": "1" * 40,
            "go-v2-native": config["sources"]["goV2Commit"],
            "rust": config["sources"]["rustCommit"],
        }[leg]
        provenance = {
            "leg": leg,
            "source": {
                "repository": "https://example.invalid/repository.git",
                "commit": source_commit,
                "tree": "2" * 40,
                "dirty": False,
            },
            "compiler": {"version": "test"},
            "buildCommand": ["build"],
            "dependencyVersions": {
                config["sources"]["goV1Module"]: config["sources"]["goV1Version"]
            },
            "binary": {
                "sha256": __import__("hashlib").sha256(binary).hexdigest(),
                "buildId": "test-build-id",
                "debugSymbols": "embedded",
            },
        }
        if leg == "go-v2-native":
            provenance["nativeDriver"] = {
                "sourceCommit": config["sources"]["nativeDriverSourceCommit"],
                "distributionCommit": config["sources"]["nativeDriverDistributionCommit"],
                "archiveSha256": "3" * 64,
                "buildId": "native-build-id",
            }
        (leg_dir / "provenance.json").write_text(
            json.dumps(provenance) + "\n", encoding="utf-8"
        )
        metrics = {
            "leg": leg,
            "workload": config["workload"],
            "counts": {"success": 3000, "error": 0, "retry": 0, "throttle": 0},
            "timing": {"wallSeconds": 300, "processCPUSeconds": 600},
            "latencyMilliseconds": {"p50": 1, "p90": 2, "p95": 3, "p99": 4, "p999": 5},
            "queueDelayMilliseconds": None,
            "memoryBytes": {
                "rss": 100,
                "hwm": 120,
                "anonymous": 80,
                "file": 20,
                "cgroupCurrent": 110,
                "cgroupPeak": 130,
            },
            "gauges": {
                "queueHighWater": None,
                "inFlightHighWater": 8,
                "completionHighWater": None,
            },
        }
        (leg_dir / "metrics.json").write_text(
            json.dumps(metrics, indent=2) + "\n", encoding="utf-8"
        )
        with (leg_dir / "perf-self.csv").open("w", newline="", encoding="utf-8") as handle:
            writer = csv.writer(handle)
            writer.writerow(["symbol", "dso", "self_percent"])
            writer.writerow(["runtime.mallocgc", "binary", "40"])
            writer.writerow(["tokio::runtime::park", "binary", "20"])
            writer.writerow(["tcp_sendmsg", "[kernel.kallsyms]", "20"])
            writer.writerow(["resolved_symbol", "binary", "20"])


if __name__ == "__main__":
    unittest.main()
