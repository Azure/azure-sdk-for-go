#!/usr/bin/env python3
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

"""Calculate sample-level quality metrics from Linux perf output."""

from __future__ import annotations

import argparse
import json
import pathlib
import re


EVENT_HEADER = re.compile(
    r"^\S.*?[ \t]+\d+[ \t]+(?:\[\d+\][ \t]+)?\d+\.\d+:.*$", re.MULTILINE
)
FRAME = re.compile(r"^[ \t]+[0-9a-f]+[ \t]+.+$", re.MULTILINE)


def calculate(script_text: str, stderr_text: str, mode: str) -> dict[str, float | int | str]:
    matches = list(EVENT_HEADER.finditer(script_text))
    samples = len(matches)
    resolved = 0
    callchains = 0
    for index, match in enumerate(matches):
        end = matches[index + 1].start() if index + 1 < samples else len(script_text)
        frames = FRAME.findall(script_text[match.end() : end])
        if frames:
            callchains += 1
        if any("unknown" not in frame.lower() for frame in frames):
            resolved += 1
    lost_match = re.search(r"lost samples:\s*(\d+)", stderr_text, re.IGNORECASE)
    lost = int(lost_match.group(1)) if lost_match else 0
    denominator = max(samples + lost, 1)
    return {
        "callGraphMode": mode,
        "samples": samples,
        "lostSamples": lost,
        "lostEventPercent": lost / denominator * 100,
        "resolvedSamplePercent": resolved / max(samples, 1) * 100,
        "callchainSamplePercent": callchains / max(samples, 1) * 100,
    }


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--perf-script", required=True, type=pathlib.Path)
    parser.add_argument("--perf-stderr", required=True, type=pathlib.Path)
    parser.add_argument("--mode", required=True, choices=("frame-pointer", "dwarf,16384"))
    parser.add_argument("--output", required=True, type=pathlib.Path)
    arguments = parser.parse_args()
    value = calculate(
        arguments.perf_script.read_text(errors="replace"),
        arguments.perf_stderr.read_text(errors="replace"),
        arguments.mode,
    )
    arguments.output.write_text(json.dumps(value, indent=2) + "\n", encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
