#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
python3 - <<'PY'
import json, os, pathlib, re, subprocess
root = pathlib.Path("/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a")
index = json.loads((root / "provenance/accepted-legs.json").read_text())
environment = os.environ.copy()
environment.update(HOME="/root", GOCACHE=str(root / "build/go-cache"))
for workers in (32, 64):
    directory = root / index[f"w{workers}/go-diagnostic/repeat-1/go-v2-native"]
    derived = directory / "derived"
    text = subprocess.check_output(
        ["/usr/local/go/bin/go", "tool", "pprof", "-raw", str(derived / "trace-sched.pprof")],
        text=True, env=environment)
    sample_section = text.split("Samples:", 1)[1].split("Locations", 1)[0]
    samples = re.findall(r"^\s*(\d+)\s+(\d+):", sample_section, re.M)
    count = sum(int(row[0]) for row in samples)
    delay = sum(int(row[1]) for row in samples)
    print("SCHEDULER", json.dumps({
        "workers": workers, "traceWindowSeconds": 5,
        "runnableTransitions": count, "cumulativeRunnableDelayMs": delay / 1e6,
        "meanRunnableDelayUs": delay / max(1, count) / 1000,
        "population": "all trace-profiled Go goroutines, not only SDK callers",
    }, separators=(",", ":")))
    for name in ("trace-sched-top.txt", "cpu-top.txt", "heap-alloc-bytes-top.txt"):
        rows = []
        for line in (derived / name).read_text().splitlines():
            if re.match(r"^\s*[\d.]+(?:ns|us|µs|ms|s|kB|MB|GB)\s+[\d.]+%", line):
                rows.append(line.strip())
        print(name, " | ".join(rows[:5]))
print("GO_SCHEDULER_ANALYSIS_INSPECTED")
PY
