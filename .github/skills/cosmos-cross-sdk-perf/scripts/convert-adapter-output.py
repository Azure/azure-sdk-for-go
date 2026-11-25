#!/usr/bin/env python3
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

"""Convert portable adapter output and launcher samples to leg-result.schema.json."""
# cSpell:words rollup smaps

from __future__ import annotations

import argparse
import json
import pathlib
import re
from typing import Any


def read_kib(path: pathlib.Path, field: str) -> int:
    if not path.is_file():
        return 0
    match = re.search(
        rf"^{re.escape(field)}:\s+(\d+)\s+kB$", path.read_text(errors="replace"), re.MULTILINE
    )
    return int(match.group(1)) * 1024 if match else 0


def read_cgroup(path: pathlib.Path, field: str) -> int:
    if not path.is_file():
        return 0
    match = re.search(rf"^{re.escape(field)}=(\d+)$", path.read_text(errors="replace"), re.MULTILINE)
    return int(match.group(1)) if match else 0


def convert(
    adapter_result: dict[str, Any],
    raw_dir: pathlib.Path,
    leg: str,
    run_class: str,
    workers: int,
    repeat: int,
) -> dict[str, Any]:
    result = {
        "leg": leg,
        "runClass": run_class,
        "workers": workers,
        "repeat": repeat,
    }
    for field in (
        "workload",
        "transportProof",
        "counts",
        "timing",
        "latencyMilliseconds",
        "queueDelayMilliseconds",
        "gauges",
    ):
        result[field] = adapter_result[field]
    result["memoryBytes"] = {
        "rss": read_kib(raw_dir / "proc-status.txt", "VmRSS"),
        "hwm": read_kib(raw_dir / "proc-status.txt", "VmHWM"),
        "anonymous": read_kib(raw_dir / "smaps-rollup.txt", "Anonymous"),
        "file": read_kib(raw_dir / "smaps-rollup.txt", "File"),
        "cgroupCurrent": read_cgroup(raw_dir / "cgroup-memory.txt", "memory.current"),
        "cgroupPeak": read_cgroup(raw_dir / "cgroup-memory.txt", "memory.peak"),
    }
    return result


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--adapter-result", required=True, type=pathlib.Path)
    parser.add_argument("--raw-dir", required=True, type=pathlib.Path)
    parser.add_argument("--leg", required=True, choices=("go-v1", "go-v2-native", "rust"))
    parser.add_argument(
        "--run-class", required=True, choices=("clean-score", "perf", "go-diagnostic")
    )
    parser.add_argument("--workers", required=True, type=int, choices=(8, 16))
    parser.add_argument("--repeat", required=True, type=int)
    parser.add_argument("--output", required=True, type=pathlib.Path)
    arguments = parser.parse_args()
    adapter_result = json.loads(arguments.adapter_result.read_text(encoding="utf-8"))
    result = convert(
        adapter_result,
        arguments.raw_dir,
        arguments.leg,
        arguments.run_class,
        arguments.workers,
        arguments.repeat,
    )
    arguments.output.write_text(
        json.dumps(result, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
