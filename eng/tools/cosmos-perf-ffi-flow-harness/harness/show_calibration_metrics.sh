#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
python3 - <<'PY'
import json, pathlib
root = pathlib.Path("/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a")
index = json.loads((root / "provenance/accepted-legs.json").read_text())
rows = []
for leg in ("go-v2-native", "rust-sdk", "rust-driver"):
    key = f"w32/perf-calibration/repeat-2/{leg}"
    directory = root / index[key]
    value = json.loads((directory / "metrics.json").read_text())
    cpu, wall, success = value["timing"]["processCPUSeconds"], value["timing"]["wallSeconds"], value["counts"]["success"]
    rows.append({
        "leg": leg, "runClass": "5-second profiled calibration, not clean-score",
        "workers": 32, "processCPUSeconds": cpu, "wallSeconds": wall,
        "success": success, "throughput": success / wall, "cpuUsPerRead": cpu * 1e6 / success,
        "busyCores": cpu / wall, "latencyMs": value["latencyMilliseconds"], "counts": value["counts"],
    })
print(json.dumps(rows, separators=(",", ":")))
print("CALIBRATION_CPU_LATENCY_VERIFIED")
PY
