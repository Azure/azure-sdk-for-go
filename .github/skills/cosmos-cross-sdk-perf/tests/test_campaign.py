# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words kallsyms mallocgc rollup smaps vcpu westus

import csv
import hashlib
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
    def test_default_config_has_governed_point_read_pilot(self):
        config = load_default()
        config["azure"]["subscriptionId"] = "00000000-0000-0000-0000-000000000000"
        config["azure"]["resourceGroup"] = "perf-rg"
        config["azure"]["existingCosmosAccountResourceId"] = (
            "/subscriptions/00000000-0000-0000-0000-000000000000/"
            "resourceGroups/cosmos-rg/providers/Microsoft.DocumentDB/databaseAccounts/account"
        )
        campaign.validate_config(config)
        self.assertEqual(
            "https://github.com/brendangregg/FlameGraph.git",
            config["sources"]["flameGraphRepository"],
        )
        self.assertRegex(config["sources"]["flameGraphCommit"], r"^[0-9a-f]{40}$")

    def test_flamegraph_script_uses_canonical_upstream_tools(self):
        script = (ROOT / "scripts" / "generate-flamegraphs.sh").read_text(
            encoding="utf-8"
        )
        self.assertIn('"$collapse" --all', script)
        self.assertIn("flamegraph.pl", script)
        self.assertIn("difffolded.pl", script)
        self.assertIn("FLAMEGRAPH-COMMIT", script)
        self.assertNotIn("custom flamegraph", script.lower())

    def test_config_rejects_non_read_operation(self):
        config = load_default()
        config["workload"]["operation"] = "mixed"
        with self.assertRaisesRegex(campaign.CampaignError, "workload.operation"):
            campaign.validate_config(config, allow_incomplete_azure=True)

    def test_config_rejects_second_pass_or_larger_matrix(self):
        config = load_default()
        config["workload"]["workers"].append(32)
        with self.assertRaisesRegex(campaign.CampaignError, "workload.workers"):
            campaign.validate_config(config, allow_incomplete_azure=True)

    def test_classification_maps_shared_key_and_totals_100(self):
        rows = [
            campaign.PerfRow("runtime.mallocgc", "go-v2", 20.0),
            campaign.PerfRow("tokio::runtime::park", "rust", 20.0),
            campaign.PerfRow("hmac_sha256", "binary", 20.0),
            campaign.PerfRow("tcp_sendmsg", "[kernel.kallsyms]", 20.0),
            campaign.PerfRow("custom_resolved_symbol", "binary", 20.0),
        ]
        totals, unmapped = campaign.classify(rows)
        self.assertEqual([], unmapped)
        self.assertAlmostEqual(100.0, sum(totals.values()))
        self.assertEqual(20.0, totals["shared-key-signing-hmac-headers"])

    def test_call_graph_calibration_selects_best_passing_mode(self):
        config = load_default()
        with tempfile.TemporaryDirectory() as temp:
            root = pathlib.Path(temp)
            config_path = root / "config.json"
            frame_path = root / "frame.json"
            dwarf_path = root / "dwarf.json"
            config_path.write_text(json.dumps(config), encoding="utf-8")
            frame_path.write_text(
                json.dumps(
                    {
                        "callGraphMode": "frame-pointer",
                        "lostEventPercent": 0,
                        "resolvedSamplePercent": 99,
                        "callchainSamplePercent": 99,
                    }
                ),
                encoding="utf-8",
            )
            dwarf_path.write_text(
                json.dumps(
                    {
                        "callGraphMode": "dwarf,16384",
                        "lostEventPercent": 0,
                        "resolvedSamplePercent": 98,
                        "callchainSamplePercent": 98,
                    }
                ),
                encoding="utf-8",
            )
            campaign.calibrate_call_graph(config_path, frame_path, dwarf_path)
            selected = json.loads(config_path.read_text(encoding="utf-8"))
            self.assertEqual("frame-pointer", selected["profiling"]["selectedCallGraph"])

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

    def test_clean_run_rejects_profile_artifacts(self):
        config = load_default()
        config["profiling"]["selectedCallGraph"] = "frame-pointer"
        with tempfile.TemporaryDirectory() as temp:
            leg_dir = pathlib.Path(temp) / "go-v1"
            self._write_leg(leg_dir, "go-v1", "clean-score", 8, 1, config)
            (leg_dir / "raw" / "cpu.pprof").write_bytes(b"forbidden")
            with self.assertRaisesRegex(campaign.CampaignError, "forbidden"):
                campaign.validate_leg(leg_dir, config)

    def test_complete_pilot_report_uses_clean_scores(self):
        config = load_default()
        config["profiling"]["selectedCallGraph"] = "frame-pointer"
        with tempfile.TemporaryDirectory() as temp:
            run_dir = pathlib.Path(temp)
            for workers in (8, 16):
                for repeat in (1, 2):
                    for leg in config["execution"]["order"][str(workers)]:
                        self._write_leg(
                            campaign.run_path(
                                run_dir, workers, "clean-score", repeat, leg
                            ),
                            leg,
                            "clean-score",
                            workers,
                            repeat,
                            config,
                        )
                for leg in config["execution"]["order"][str(workers)]:
                    self._write_leg(
                        campaign.run_path(run_dir, workers, "perf", 1, leg),
                        leg,
                        "perf",
                        workers,
                        1,
                        config,
                    )
                for leg in ("go-v1", "go-v2-native"):
                    self._write_leg(
                        campaign.run_path(run_dir, workers, "go-diagnostic", 1, leg),
                        leg,
                        "go-diagnostic",
                        workers,
                        1,
                        config,
                    )
            campaign.report(run_dir, config)
            summary = json.loads(
                (run_dir / "comparison" / "summary.json").read_text(encoding="utf-8")
            )
            self.assertEqual("clean-score", summary["authoritativeClass"])
            self.assertEqual(6, len(summary["authoritative"]))
            self.assertEqual(10, len(summary["observerEffect"]))
            self.assertEqual(22, len(summary["allRuns"]))
            self.assertTrue((run_dir / "dashboard.html").is_file())
            self.assertTrue(
                (run_dir / "comparison" / "observer-effect.csv").is_file()
            )
            self.assertTrue((run_dir / "comparison" / "setup-equivalence.csv").is_file())
            self.assertIn(
                "## Observer effect",
                (run_dir / "comparison" / "report.md").read_text(encoding="utf-8"),
            )
            self.assertTrue((run_dir / "SHA256SUMS").is_file())

    def test_sanitize_redacts_sensitive_json_and_includes_dashboard(self):
        with tempfile.TemporaryDirectory() as temp:
            run_dir = pathlib.Path(temp) / "run"
            run_dir.mkdir()
            (run_dir / "summary.json").write_text(
                json.dumps({"endpoint": "https://account", "success": 10}), encoding="utf-8"
            )
            (run_dir / "dashboard.html").write_text("<html></html>", encoding="utf-8")
            (run_dir / "observer-effect.csv").write_text(
                "workers,leg,run_class\n8,go-v1,perf\n", encoding="utf-8"
            )
            differential = run_dir / "w8" / "perf" / "differential" / "go-v1-vs-go-v2"
            differential.mkdir(parents=True)
            (differential / "forward.svg").write_text(
                "<svg></svg>", encoding="utf-8"
            )
            output = pathlib.Path(temp) / "bundle.tar.gz"
            campaign.sanitize(run_dir, output)
            with tarfile.open(output, "r:gz") as archive:
                member = archive.extractfile("cosmos-perf-result/summary.json")
                assert member is not None
                value = json.load(member)
                self.assertIsNotNone(
                    archive.getmember("cosmos-perf-result/dashboard.html")
                )
                self.assertIsNotNone(
                    archive.getmember("cosmos-perf-result/observer-effect.csv")
                )
                self.assertIsNotNone(
                    archive.getmember(
                        "cosmos-perf-result/w8/perf/differential/"
                        "go-v1-vs-go-v2/forward.svg"
                    )
                )
            self.assertTrue(value["endpoint"].startswith("sha256:"))
            self.assertEqual(10, value["success"])

    def _write_leg(self, leg_dir, leg, run_class, workers, repeat, config):
        raw = leg_dir / "raw"
        raw.mkdir(parents=True)
        binary = b"test-binary-" + leg.encode()
        (raw / "binary").write_bytes(binary)
        if run_class in ("clean-score", "perf"):
            for path in (
                "proc-status.txt",
                "smaps-rollup.txt",
                "smaps.txt",
                "cgroup-memory.txt",
                "thread-stacks.txt",
            ):
                (raw / path).write_text("evidence\n", encoding="utf-8")
        if run_class == "perf":
            (raw / "perf.data").write_bytes(b"perf")
            (raw / "perf-quality.json").write_text(
                json.dumps(
                    {
                        "callGraphMode": config["profiling"]["selectedCallGraph"],
                        "lostEventPercent": 0,
                        "resolvedSamplePercent": 100,
                        "callchainSamplePercent": 100,
                    }
                ),
                encoding="utf-8",
            )
        if run_class == "go-diagnostic":
            for path in (
                "cpu.pprof",
                "heap.pprof",
                "runtime.trace",
                "block.pprof",
                "mutex.pprof",
            ):
                (raw / path).write_bytes(path.encode())

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
                "sha256": hashlib.sha256(binary).hexdigest(),
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

        workload = dict(config["workload"])
        workload["workers"] = workers
        metrics = {
            "leg": leg,
            "runClass": run_class,
            "workers": workers,
            "repeat": repeat,
            "workload": workload,
            "transportProof": {
                "transportKind": "gateway-v1",
                "negotiatedHTTPVersion": "HTTP/2",
                "binaryEncoding": False,
            },
            "counts": {"success": 3000 + repeat, "error": 0, "retry": 0, "throttle": 0},
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
                "inFlightHighWater": workers,
                "completionHighWater": None,
            },
        }
        (leg_dir / "metrics.json").write_text(
            json.dumps(metrics, indent=2) + "\n", encoding="utf-8"
        )
        if run_class == "perf":
            with (leg_dir / "perf-self.csv").open(
                "w", newline="", encoding="utf-8"
            ) as handle:
                writer = csv.writer(handle)
                writer.writerow(["symbol", "dso", "self_percent"])
                writer.writerow(["runtime.mallocgc", "binary", "40"])
                writer.writerow(["tokio::runtime::park", "binary", "20"])
                writer.writerow(["tcp_sendmsg", "[kernel.kallsyms]", "20"])
                writer.writerow(["resolved_symbol", "binary", "20"])


if __name__ == "__main__":
    unittest.main()
