import pathlib
import tempfile
import unittest
import patch_native
import copy
import validate_flow
import run_leg
import calibrate
import report
import json
import xml.etree.ElementTree as ET
import patch_wire
import patch_go
import shutil
import os
import validate_rust_flow
import attribute
import profile_attribution


class NativePatchTests(unittest.TestCase):
    def test_exact_pinned_source_anchors_and_compile_gate(self):
        source = pathlib.Path(__file__).parent
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            crate = root / "sdk/cosmos/azure_data_cosmos_driver_native"
            (crate / "src").mkdir(parents=True)
            (crate / "Cargo.toml").write_text((source / "native-Cargo.toml").read_text())
            for name in ("lib.rs", "runtime.rs", "submit.rs", "completion.rs"):
                (crate / "src" / name).write_text((source / "reference" / name).read_text())
            patch_native.patch_native(root)
            submit = (crate / "src/submit.rs").read_text()
            self.assertEqual(submit.count("flow!(first_poll());"), 1)
            self.assertEqual(submit.count("flow!(terminal());"), 1)
            completion = (crate / "src/completion.rs").read_text()
            self.assertEqual(completion.count("flow!(event(5, 1));"), 2)
            self.assertEqual(completion.count("flow!(queue("), 6)
            lib = (crate / "src/lib.rs").read_text()
            self.assertIn('#[cfg(feature = "__internal_perf_flow")]\nmod perf_flow;', lib)
            self.assertIn(".worker_threads(2)", (crate / "src/runtime.rs").read_text())
            with self.assertRaises(ValueError):
                patch_native.patch_native(root)

    def test_anchor_mismatch_is_an_error(self):
        with self.assertRaises(ValueError):
            patch_native.once("absent", "expected", "new")
        with self.assertRaises(ValueError):
            patch_native.once("xx", "x", "new")


class GoPatchTests(unittest.TestCase):
    @staticmethod
    def add_user_agent_fixture(sdk):
        (sdk / "user_agent_native_test.go").write_text(
            "\t\t\tsuffixPattern := ` ` + regexp.QuoteMeta(tt.applicationID)\n"
            '\t\t\tif tt.applicationID == "" {\n'
            "\t\t\t\tsuffixPattern = ` rustc/[^ |]+`\n\t\t\t}\n"
            r"require.Regexp(t, `(?:^| )azsdk-rust-cosmos-driver/\S+(?: |$)`, agent)" "\n"
            r"require.Regexp(t, suffixPattern+`\|F[0-9A-Fa-f]+$`, agent)" "\n")

    def test_exact_optimized_hooks_and_adapter_binary_disablement(self):
        source = pathlib.Path(__file__).parent
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            shutil.copytree(source / "reference/go-context/sdk", root / "sdk")
            shutil.copytree(source / "reference/go-context/adapter", root / "adapter")
            self.add_user_agent_fixture(root / "sdk")
            patch_go.patch_sdk(root / "sdk", True)
            patch_go.patch_adapter(root / "adapter", True)
            reactor = (root / "sdk/reactor_native.go").read_text()
            operation = (root / "sdk/operation_native.go").read_text()
            adapter = (root / "adapter/main.go").read_text()
            self.assertEqual(reactor.count("flowWait(int(n))"), 1)
            self.assertLess(reactor.index("flowEvent(flowDelivered, 1)"), reactor.index("p.result <- result"))
            self.assertEqual(operation.count("flowEvent(flowSubmitAttempts, 1)"), 1)
            self.assertEqual(operation.count("flowEvent(flowSubmitAccepted, 1)"), 1)
            self.assertEqual(operation.count("flowEvent(flowSubmitErrors, 1)"), 1)
            self.assertEqual(operation.count("flowResumed(pending.flowDeliveredAt)"), 2)
            self.assertIn("Operation: azcosmos.OperationOptions{BinaryEncoding:", (root / "adapter/sdk_v2.go").read_text())
            self.assertIn("Disabled: true", (root / "adapter/sdk_v2.go").read_text())
            self.assertLess(adapter.index("value.warmup, false"), adapter.index("startDiagnostics()"))
            self.assertIn("measurementBounds(monotonicStart, monotonicEnd)", adapter)
            self.assertIn("SDK preflight payload validation", adapter)
            self.assertIn(r"drv=1\.1\.0-beta\.3", (root / "sdk/user_agent_native_test.go").read_text())

    def test_clean_source_has_no_flow_hot_path_hooks(self):
        source = pathlib.Path(__file__).parent
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            shutil.copytree(source / "reference/go-context/sdk", root / "sdk")
            shutil.copytree(source / "reference/go-context/adapter", root / "adapter")
            self.add_user_agent_fixture(root / "sdk")
            patch_go.patch_sdk(root / "sdk", False)
            patch_go.patch_adapter(root / "adapter", False)
            self.assertNotIn("flowWait(", (root / "sdk/reactor_native.go").read_text())
            self.assertNotIn("flowEvent(", (root / "sdk/operation_native.go").read_text())
            self.assertFalse((root / "sdk/experimental_flow.go").exists())

class WirePatchTests(unittest.TestCase):
    def test_resume_accepts_rustfmt_multiline_call_and_is_idempotent(self):
        source = """crate::perf_wire::point_response(
            response.body(),
            plan.transcodes_response_to_text(),
        );"""
        value = patch_wire.normalize_wire_call(source)
        self.assertIn("point_response(response,", value)
        self.assertEqual(value, patch_wire.normalize_wire_call(value))
        with self.assertRaises(ValueError):
            patch_wire.normalize_wire_call("no call")

    def test_wire_hook_is_feature_gated_and_before_original_conversion(self):
        source = pathlib.Path(__file__).parent
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            driver = root / "sdk/cosmos/azure_data_cosmos_driver"
            native = root / "sdk/cosmos/azure_data_cosmos_driver_native"
            (driver / "src/driver").mkdir(parents=True)
            (native / "src").mkdir(parents=True)
            (driver / "Cargo.toml").write_text((source / "reference/driver-Cargo.toml").read_text())
            (driver / "src/lib.rs").write_text((source / "reference/driver_lib.rs").read_text())
            (driver / "src/driver/cosmos_driver.rs").write_text((source / "reference/cosmos_driver.rs").read_text())
            (native / "Cargo.toml").write_text((source / "native-Cargo.toml").read_text())
            (native / "src/runtime.rs").write_text((source / "reference/runtime.rs").read_text())
            patch_wire.patch_wire(root)
            text = (driver / "src/driver/cosmos_driver.rs").read_text()
            begin = text.index("    fn finalize_plan_response(")
            end = text.index("    async fn uses_cache_only_topology(", begin)
            method = text[begin:end]
            self.assertLess(method.index("crate::perf_wire::point_response("),
                            method.index("if plan.transcodes_response_to_text()"))
            self.assertIn("transcode_body_to_text_with_utf8_mode", method)
            self.assertIn("plan.poison_continuation()", method)
            self.assertIn('#[cfg(feature = "__internal_perf_wire")]', method)
            with self.assertRaises(ValueError):
                patch_wire.patch_wire(root)

class LedgerTests(unittest.TestCase):
    def setUp(self):
        boundaries = {"phaseStartUnixNs": [0, 0, 1_250_000_000],
                      "phaseEndUnixNs": [0, 0, 2_500_000_000]}
        self.go = {
            "overflow": 0, **boundaries,
            "rows": [
                {"phase": 1, "unixSecond": 0, "counts": {"submitAttempts": 900}, "batchHistogram": {}},
                {"phase": 2, "unixSecond": 1,
                 "counts": {"submitAttempts": 3, "submitAccepted": 3, "submitErrors": 0,
                            "drained": 3, "delivered": 3, "resumed": 3, "waitReturns": 2, "emptyWaits": 0},
                 "batchHistogram": {"1": 1, "2": 1}, "deliveryToResumeTotalNs": 300, "deliveryToResumeMaxNs": 100},
                {"phase": 2, "unixSecond": 2,
                 "counts": {"submitAttempts": 1, "submitAccepted": 1, "submitErrors": 0,
                            "drained": 1, "delivered": 1, "resumed": 1, "waitReturns": 2, "emptyWaits": 1},
                 "batchHistogram": {"1": 1}, "deliveryToResumeTotalNs": 100, "deliveryToResumeMaxNs": 100},
            ],
        }
        counts = {name: 4 for name in ("admitted", "task_submitted", "first_poll", "terminal", "enqueued", "drained")}
        counts.update(admission_aborted=0, enqueue_rejected=0)
        self.native = {
            "overflow": 0, "activeAtSnapshot": 0, "readyAtSnapshot": 0,
            "admittedUndrainedAtSnapshot": 0, **boundaries,
            "rows": [{"phase": 2, "unixSecond": 1, "counts": counts, "readyHigh": 2,
                      "admittedUndrainedHigh": 4, "polledUnfinishedHigh": 3}],
        }
        self.metrics = {"counts": {"success": 4, "error": 0, "retry": 0, "throttle": 0},
                        "timing": {"wallSeconds": 1.25}}

    def test_warmup_excluded_and_partial_seconds_normalized(self):
        result = validate_flow.validate(self.go, self.native, self.metrics, 32)
        self.assertEqual(result["success"], 4)
        self.assertEqual(result["submissionsPerSecond"], 3.2)
        self.assertEqual(result["completionBatchesPerSecond"], 2.4)
        self.assertAlmostEqual(result["meanCompletionBatchSize"], 4 / 3)
        self.assertEqual([row["durationSeconds"] for row in result["buckets"]], [0.75, 0.5])
        self.assertEqual(result["buckets"][0]["goRates"]["submitAttempts"], 4)
        self.assertEqual(result["buckets"][1]["goRates"]["submitAttempts"], 2)

    def test_counter_drift_rejected(self):
        for field in ("admitted", "terminal", "enqueued", "drained"):
            with self.subTest(field=field):
                native = copy.deepcopy(self.native)
                native["rows"][0]["counts"][field] -= 1
                with self.assertRaises(ValueError):
                    validate_flow.validate(self.go, native, self.metrics, 32)

    def test_retries_clock_and_batch_drift_rejected(self):
        metrics = copy.deepcopy(self.metrics)
        metrics["counts"]["retry"] = 1
        with self.assertRaises(ValueError):
            validate_flow.validate(self.go, self.native, metrics, 32)
        native = copy.deepcopy(self.native)
        native["overflow"] = 1
        with self.assertRaises(ValueError):
            validate_flow.validate(self.go, native, self.metrics, 32)
        go = copy.deepcopy(self.go)
        go["rows"][1]["batchHistogram"]["2"] = 2
        with self.assertRaises(ValueError):
            validate_flow.validate(go, self.native, self.metrics, 32)

    def test_closed_loop_high_water_rejected(self):
        native = copy.deepcopy(self.native)
        native["rows"][0]["readyHigh"] = 33
        with self.assertRaises(ValueError):
            validate_flow.validate(self.go, native, self.metrics, 32)

    def test_phase_boundary_mismatch_rejected(self):
        native = copy.deepcopy(self.native)
        native["phaseEndUnixNs"] = [0, 0, 2_600_000_000]
        with self.assertRaises(ValueError):
            validate_flow.validate(self.go, native, self.metrics, 32)

    def test_wall_clock_drift_from_monotonic_timing_rejected(self):
        metrics = copy.deepcopy(self.metrics)
        metrics["timing"]["wallSeconds"] = 1.5
        with self.assertRaises(ValueError):
            validate_flow.validate(self.go, self.native, metrics, 32)


class SamplerGateTests(unittest.TestCase):
    @unittest.skipIf(os.name == "nt", "POSIX execute bits are verified on the Linux VM")
    def test_hardlinked_raw_binary_keeps_execute_permission(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            staged = root / "staged"
            staged.write_text("#!/bin/sh\nexit 0\n")
            staged.chmod(0o555)
            raw = root / "raw"
            raw.mkdir()
            os.link(staged, raw / "binary")
            run_leg.seal(raw)
            self.assertTrue(staged.stat().st_mode & 0o111)
            self.assertEqual(staged.stat().st_mode & 0o222, 0)

    def test_accepted_leg_index_is_explicit_and_rejects_traversal(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            (root / "provenance").mkdir()
            path = root / "run/w32/perf/repeat-1/go-v2-native-attempt-2"
            run_leg.record_accepted(root, "w32/perf/repeat-1/go-v2-native", path)
            self.assertEqual(run_leg.accepted_directory(root, 32, "perf", 1, "go-v2-native"), path)
            (root / "provenance/accepted-legs.json").write_text(
                json.dumps({"w32/perf/repeat-1/go-v2-native": "../outside"}))
            with self.assertRaises(ValueError):
                run_leg.accepted_directory(root, 32, "perf", 1, "go-v2-native")

    def setUp(self):
        self.value = {
            "counts": {"success": 20, "error": 0, "retry": 0, "throttle": 0},
            "timing": {"wallSeconds": 10, "processCPUSeconds": 20},
            "workload": {
                "workers": 32, "closedLoop": True, "binaryEncoding": False,
                "partitionKeyPath": "/pk", "seedCount": 1024, "documentBytes": 1024,
                "operation": "point-read", "keySequence": "round-robin-0-through-1023",
                "rawResponseMode": True, "validationOutsideTimedPath": True,
            },
            "transportProof": {"transportKind": "gateway-v1", "negotiatedHTTPVersion": "HTTP/2"},
            "latencyMilliseconds": {"p50": 1, "p90": 2, "p95": 3, "p99": 4, "p999": 5},
        }

    def test_exact_nanosecond_perf_time(self):
        self.assertEqual(run_leg.perf_time(1234567890123456), "1234567.890123456")
        self.assertEqual(run_leg.perf_time(1), "0.000000001")

    def test_memory_accounting_accepts_v1_and_v2_without_silent_zero(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            with self.assertRaises(RuntimeError):
                run_leg.memory_value(root)
            (root / "memory.usage_in_bytes").write_text("100")
            (root / "memory.max_usage_in_bytes").write_text("200")
            self.assertEqual(run_leg.memory_value(root), 100)
            self.assertEqual(run_leg.memory_value(root, peak=True), 200)
            (root / "memory.current").write_text("120")
            (root / "memory.peak").write_text("220")
            self.assertEqual(run_leg.memory_value(root), 120)
            self.assertEqual(run_leg.memory_value(root, peak=True), 220)

    def test_expected_leg_contract(self):
        run_leg.validate_metrics(self.value, 32)

    def test_wrong_workload_is_rejected(self):
        for field, value in (("workers", 64), ("binaryEncoding", True), ("seedCount", 100),
                             ("rawResponseMode", False), ("partitionKeyPath", "/id")):
            with self.subTest(field=field):
                leg = copy.deepcopy(self.value)
                leg["workload"][field] = value
                with self.assertRaises(ValueError):
                    run_leg.validate_metrics(leg, 32)

    def test_errors_and_cpu_capacity_drift_rejected(self):
        for name in ("error", "retry", "throttle"):
            leg = copy.deepcopy(self.value)
            leg["counts"][name] = 1
            with self.assertRaises(ValueError):
                run_leg.validate_metrics(leg, 32)
        leg = copy.deepcopy(self.value)
        leg["timing"]["processCPUSeconds"] = 50
        with self.assertRaises(ValueError):
            run_leg.validate_metrics(leg, 32)

    def test_raw_manifest_matches_and_is_not_silently_overwritten(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            (root / "value.txt").write_text("evidence\n")
            run_leg.seal(root)
            line = (root / "SHA256SUMS").read_text().strip()
            self.assertEqual(line, f"{run_leg.digest(root / 'value.txt')}  value.txt")
            self.assertEqual((root / "value.txt").stat().st_mode & 0o222, 0)
            with self.assertRaises(FileExistsError):
                run_leg.seal(root)
            # Windows TemporaryDirectory needs writable files for deterministic cleanup.
            for path in root.iterdir():
                path.chmod(0o600)


class CallchainTests(unittest.TestCase):
    @staticmethod
    def script(leaf="azure_data_cosmos_driver::execute", count=100):
        sample = (f"worker 101 [000] 12.345678901: 12345 cpu-clock:\n"
                  f"\t7f00001 {leaf} (/mnt/cosmos-perf/cosmos-test)\n"
                  "\t7f00002 future::poll (/mnt/cosmos-perf/cosmos-test)\n"
                  "\t7f00003 runtime::worker (/mnt/cosmos-perf/cosmos-test)\n"
                  "\t7f00004 thread::start (/mnt/cosmos-perf/cosmos-test)\n\n")
        return sample * count

    def test_quality_requires_resolved_leaf_and_real_native_symbol(self):
        quality = calibrate.calculate(self.script(), "")
        self.assertTrue(quality["valid"])
        self.assertEqual(quality["samples"], 100)
        self.assertEqual(quality["medianFrameDepth"], 4)
        false_path_match = calibrate.calculate(self.script("malloc"), "")
        self.assertEqual(false_path_match["nativeFrameSamplePercent"], 0)
        self.assertFalse(false_path_match["valid"])
        unresolved = calibrate.calculate(self.script("[unknown]"), "")
        self.assertEqual(unresolved["resolvedLeafPercent"], 0)
        self.assertFalse(unresolved["valid"])

    def test_lost_events_fail_calibration(self):
        quality = calibrate.calculate(self.script() + "LOST 2 events\n", "")
        self.assertFalse(quality["valid"])
        self.assertEqual(quality["lostSamples"], 2)

    def test_better_dwarf_depth_wins_and_no_passing_mode_fails(self):
        fp = calibrate.calculate(self.script(), "")
        dwarf = dict(fp, medianFrameDepth=10)
        qualities = {leg: {"fp": fp, "dwarf,16384": dwarf}
                     for leg in ("go-v2-native", "rust-sdk", "rust-driver")}
        self.assertEqual(calibrate.choose(qualities), "dwarf,16384")
        qualities["rust-driver"] = {"fp": dict(fp, valid=False), "dwarf,16384": dict(dwarf, valid=False)}
        with self.assertRaises(ValueError):
            calibrate.choose(qualities)


class ReportTests(unittest.TestCase):
    def test_rust_benchmark_namespace_is_not_a_c_abi_export(self):
        self.assertEqual(attribute.classify(attribute.Frame("cosmos_perf_flow::phase", "/tmp/rust-bench")),
                         "resolved-other")
        self.assertEqual(attribute.classify(attribute.Frame("cosmos_submit_singleton_operation", "/tmp/go-v2")),
                         "cgo-completion-ffi")
    def test_profile_cpu_clock_parser_and_four_core_flywheel(self):
        with tempfile.TemporaryDirectory() as temporary:
            path = pathlib.Path(temporary) / "perf.txt"
            path.write_text("worker 123 1.123: 20000000 cpu-clock:\n"
                            "\t7f001 malloc (/usr/lib/libc.so.6)\n"
                            "\t7f002 azure_data_cosmos_driver::execute (/tmp/cosmos-test)\n")
            samples = attribute.parse_perf(path)
            self.assertEqual(len(samples), 1)
            metrics = {"counts": {"success": 100, "error": 0, "retry": 0, "throttle": 0},
                       "timing": {"processCPUSeconds": 1, "wallSeconds": 2}}
            summary, rows = attribute.summarize(samples, metrics,
                {"samples": 1, "lostSamples": 0, "lostEventPercent": 0})
            self.assertEqual(summary["selfPercentTotal"], 100)
            self.assertEqual(rows[0]["selfCategory"], "serialization-body-allocation-memcopy")
            svg = profile_attribution.flywheel(summary, "Test profile")
            ET.fromstring(svg)
            self.assertIn("capacity not charged", svg)

    def test_pooled_cpu_throughput_and_percentile_labeling(self):
        memory = {"sampleCount": 1, "meanRSSMiB": 4.0, "peakSampledRSSMiB": 6.0}
        rows = [
            {"workers": 32, "leg": "go-v2-native", "memory": memory,
             "metrics": {"counts": {"success": 3}, "timing": {"wallSeconds": 1, "processCPUSeconds": 1},
                         "latencyMilliseconds": {"p50": 2, "p95": 4, "p99": 6}}},
            {"workers": 32, "leg": "go-v2-native", "memory": memory,
             "metrics": {"counts": {"success": 1}, "timing": {"wallSeconds": 3, "processCPUSeconds": 1},
                         "latencyMilliseconds": {"p50": 4, "p95": 6, "p99": 8}}},
        ]
        result = report.aggregate(rows, 32, "go-v2-native")
        self.assertEqual(result["throughput"], 1)
        self.assertEqual(result["busyCores"], 0.5)
        self.assertEqual(result["fourVCPUCapacityPercent"], 12.5)
        self.assertEqual(result["cpuMicrosecondsPerSuccess"], 500_000)
        self.assertEqual(result["meanRepeatP99Ms"], 7)
        rows_for_chart = [dict(result, leg=leg, workers=workers) for workers in (32, 64) for leg in report.LEGS]
        svg = report.bars(rows_for_chart, "throughput", "Throughput", "Reads/s")
        parsed = ET.fromstring(svg)
        self.assertEqual(parsed.tag, "{http://www.w3.org/2000/svg}svg")
        with self.assertRaises(ValueError):
            report.aggregate(rows[:1], 32, "go-v2-native")

    def test_raw_tamper_and_unmanifested_files_rejected(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            (root / "value.txt").write_text("evidence\n")
            run_leg.seal(root)
            report.verify_raw(root)
            (root / "extra.txt").write_text("not in manifest")
            with self.assertRaises(ValueError):
                report.verify_raw(root)
            (root / "extra.txt").unlink()
            (root / "value.txt").chmod(0o600)
            (root / "value.txt").write_text("tampered\n")
            with self.assertRaises(ValueError):
                report.verify_raw(root)
            for path in root.iterdir():
                path.chmod(0o600)

    def test_memory_window_excludes_warmup_and_checks_split(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            (root / "barriers").mkdir()
            (root / "barriers/measurement-bounds.json").write_text(
                json.dumps({"monotonicStartNs": 100, "monotonicEndNs": 200}))
            status = "VmRSS: 16 kB\nVmHWM: 20 kB\nRssAnon: 12 kB\nRssFile: 4 kB\nRssShmem: 0 kB\n"
            samples = [{"monotonicNs": stamp, "status": status, "cgroupMemoryCurrent": 32768,
                        "cgroupMemoryPeak": 65536} for stamp in (50, 150, 250)]
            (root / "samples.jsonl").write_text("\n".join(map(json.dumps, samples)))
            result = report.memory_samples(root)
            self.assertEqual(result["sampleCount"], 1)
            self.assertEqual(result["meanRSSMiB"], 16 / 1024)
            self.assertEqual(result["meanAnonMiB"] + result["meanFileMiB"], result["meanRSSMiB"])


class RustRateTests(unittest.TestCase):
    def test_rust_completion_buckets_exclude_warmup_and_reconcile(self):
        metrics = {
            "operationBuckets": {"1": 3, "2": 1},
            "counts": {"success": 4},
            "timing": {"measurementStartUnixNs": 1_250_000_000,
                       "measurementEndUnixNs": 2_500_000_000, "wallSeconds": 1.25},
        }
        result = validate_rust_flow.validate(metrics)
        self.assertEqual(result["operationsPerSecond"], 3.2)
        self.assertEqual(result["buckets"][0]["durationSeconds"], 0.75)
        metrics["counts"]["success"] = 5
        with self.assertRaises(ValueError):
            validate_rust_flow.validate(metrics)
if __name__ == "__main__":
    unittest.main()
