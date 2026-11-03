# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words rollup smaps

import importlib.util
import pathlib
import tempfile
import unittest


SKILL = pathlib.Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location(
    "convert_adapter_output", SKILL / "scripts" / "convert-adapter-output.py"
)
assert SPEC and SPEC.loader
CONVERTER = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(CONVERTER)
PERF_SPEC = importlib.util.spec_from_file_location(
    "perf_quality", SKILL / "scripts" / "perf-quality.py"
)
assert PERF_SPEC and PERF_SPEC.loader
PERF_QUALITY = importlib.util.module_from_spec(PERF_SPEC)
PERF_SPEC.loader.exec_module(PERF_QUALITY)


class AdapterOutputTests(unittest.TestCase):
    def test_perf_quality_accepts_script_headers_with_and_without_cpu(self):
        script = """\
cosmos-perf-go-  2670 521.716053: 20408163 task-clock:ppp:
          4aaf23 runtime.futex.abi0+0x23 (binary)

cosmos-perf-go-  2678 [002] 522.186463: 20408163 task-clock:ppp:
          0 [unknown] ([unknown])

"""
        result = PERF_QUALITY.calculate(script, "", "frame-pointer")
        self.assertEqual(2, result["samples"])
        self.assertEqual(100, result["callchainSamplePercent"])
        self.assertEqual(50, result["resolvedSamplePercent"])

    def test_conversion_is_deterministic_and_reads_memory(self):
        adapter = {
            "workload": {"workers": 8},
            "transportProof": {
                "transportKind": "gateway-v1",
                "negotiatedHTTPVersion": "HTTP/2",
                "binaryEncoding": False,
            },
            "counts": {"success": 1, "error": 0, "retry": 0, "throttle": 0},
            "timing": {"wallSeconds": 1, "processCPUSeconds": 1},
            "latencyMilliseconds": {"p50": 1, "p90": 1, "p95": 1, "p99": 1, "p999": 1},
            "queueDelayMilliseconds": None,
            "gauges": {
                "queueHighWater": None,
                "inFlightHighWater": 8,
                "completionHighWater": None,
            },
        }
        with tempfile.TemporaryDirectory() as directory:
            raw = pathlib.Path(directory)
            (raw / "proc-status.txt").write_text("VmRSS:\t10 kB\nVmHWM:\t12 kB\n")
            (raw / "smaps-rollup.txt").write_text("Anonymous: 4 kB\nFile: 3 kB\n")
            (raw / "cgroup-memory.txt").write_text("memory.current=99\nmemory.peak=101\n")
            result = CONVERTER.convert(adapter, raw, "go-v1", "clean-score", 8, 1)
        self.assertEqual(result["leg"], "go-v1")
        self.assertEqual(
            result["memoryBytes"],
            {
                "rss": 10 * 1024,
                "hwm": 12 * 1024,
                "anonymous": 4 * 1024,
                "file": 3 * 1024,
                "cgroupCurrent": 99,
                "cgroupPeak": 101,
            },
        )


if __name__ == "__main__":
    unittest.main()
