#!/usr/bin/env python3
import argparse
import json
import pathlib
import re
import statistics
import hashlib
from run_leg import accepted_directory


HEADER = re.compile(r"^\S.*?\s+\d+(?:/\d+)?\s+(?:\[\d+\]\s+)?\d+\.\d+:.*$", re.M)
FRAME = re.compile(r"^\s+[0-9a-fA-F]+\s+(.+)$", re.M)


def calculate(script, errors):
    headers = list(HEADER.finditer(script))
    depths, resolved, native = [], 0, 0
    for index, match in enumerate(headers):
        end = headers[index + 1].start() if index + 1 < len(headers) else len(script)
        frames = FRAME.findall(script[match.end():end])
        depths.append(len(frames))
        if frames:
            leaf = frames[0].rpartition(" (")[0] or frames[0]
            if not ("unknown" in leaf.lower() or re.fullmatch(r"(?:0x)?[0-9a-fA-F]+", leaf)):
                resolved += 1
        if any(any(name in frame.rpartition(" (")[0] for name in
                   ("azure_data_cosmos", "azurecosmosdriver", "tokio", "hyper", "rustls")) for frame in frames):
            native += 1
    lost = sum(int(value) for value in re.findall(r"\bLOST\s+(\d+)\s+events", script, re.I))
    recorder_lost = re.findall(r"\blost\s+(\d+)\s+(?:events|samples)", errors, re.I)
    if recorder_lost:
        lost = max(lost, sum(map(int, recorder_lost)))
    count = len(headers)
    value = {
        "samples": count, "lostSamples": lost,
        "lostPercent": lost / max(1, count + lost) * 100,
        "resolvedLeafPercent": resolved / max(1, count) * 100,
        "multiFramePercent": sum(depth >= 2 for depth in depths) / max(1, count) * 100,
        "medianFrameDepth": statistics.median(depths) if depths else 0,
        "nativeFrameSamplePercent": native / max(1, count) * 100,
    }
    value["valid"] = (
        count >= 100 and value["lostPercent"] <= 1
        and value["resolvedLeafPercent"] >= 95 and value["multiFramePercent"] >= 90
        and value["medianFrameDepth"] >= 4 and value["nativeFrameSamplePercent"] >= 5
    )
    return value


def choose(qualities):
    legs = ("go-v2-native", "rust-sdk", "rust-driver")
    passing = [mode for mode in ("fp", "dwarf,16384")
               if all(qualities[leg][mode]["valid"] for leg in legs)]
    if not passing:
        raise ValueError("no call-graph mode passes on every requested leg")
    if len(passing) == 1:
        return passing[0]
    for leg in legs:
        fp, dwarf = qualities[leg]["fp"], qualities[leg]["dwarf,16384"]
        if (fp["resolvedLeafPercent"] + 2 < dwarf["resolvedLeafPercent"]
            or fp["multiFramePercent"] + 2 < dwarf["multiFramePercent"]
            or fp["medianFrameDepth"] < dwarf["medianFrameDepth"] * 0.8):
            return "dwarf,16384"
    return "fp"


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("root", type=pathlib.Path)
    args = parser.parse_args()
    values, binaries = {}, set()
    for leg in ("go-v2-native", "rust-sdk", "rust-driver"):
        values[leg] = {}
        for repeat, mode in ((1, "fp"), (2, "dwarf,16384")):
            directory = accepted_directory(args.root, 32, "perf-calibration", repeat, leg)
            valid = json.loads((directory / "VALID.json").read_text())
            if not valid["valid"]:
                raise ValueError("rejected calibration leg")
            binaries.add(valid["binarySHA256"])
            raw = directory / "raw"
            script = raw / "perf-script.txt"
            metadata_path = directory / "derived/analysis-source.json"
            if metadata_path.exists():
                metadata = json.loads(metadata_path.read_text())
                relative = pathlib.PurePosixPath(metadata["perfScript"])
                if relative.is_absolute() or ".." in relative.parts:
                    raise ValueError("invalid corrected analysis path")
                digest = hashlib.sha256()
                with (raw / "perf.data").open("rb") as handle:
                    for chunk in iter(lambda: handle.read(1024 * 1024), b""):
                        digest.update(chunk)
                if digest.hexdigest() != metadata["rawPerfSHA256"]:
                    raise ValueError("corrected analysis raw-evidence hash mismatch")
                script = directory.joinpath(*relative.parts)
            values[leg][mode] = calculate(script.read_text(),
                                         (raw / "perf-record.stderr").read_text())
    mode = choose(values)
    target = args.root / "provenance/callgraph-calibration.json"
    if target.exists():
        value = json.loads(target.read_text())
        if not value["valid"] or value["binarySHA256"] != sorted(binaries) or value["callgraph"] != mode:
            raise ValueError("existing calibration differs from current accepted evidence")
        return
    with target.open("x") as handle:
        json.dump({"valid": True, "callgraph": mode, "binarySHA256": sorted(binaries), "qualities": values,
                   "sampleEvent": "cpu-clock", "frequencyHz": 49,
                   "thresholds": {"minimumSamples": 100, "maximumLostPercent": 1,
                                  "minimumResolvedLeafPercent": 95, "minimumMultiFramePercent": 90,
                                  "minimumMedianDepth": 4, "minimumNativeFrameSamplePercent": 5}}, handle, indent=2)
        handle.write("\n")


if __name__ == "__main__":
    main()
