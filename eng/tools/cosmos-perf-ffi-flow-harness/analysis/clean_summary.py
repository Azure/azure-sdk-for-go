#!/usr/bin/env python3
import csv
import hashlib
import html
import json
import pathlib
import sys

from report import aggregate, bars, memory_samples, verify_raw
from run_leg import accepted_directory, validate_metrics

root = pathlib.Path(sys.argv[1])
rows, results = [], []
for workers in (32, 64):
    for leg in ("go-v2-native", "rust-sdk", "rust-driver"):
        selected = []
        for repeat in (1, 2):
            directory = accepted_directory(root, workers, "clean-score", repeat, leg)
            valid = json.loads((directory / "VALID.json").read_text())
            assert valid["valid"] and not (directory / "REJECTED.json").exists()
            verify_raw(directory / "raw")
            metrics = json.loads((directory / "metrics.json").read_text())
            validate_metrics(metrics, workers)
            row = {"workers": workers, "leg": leg, "repeat": repeat,
                   "metrics": metrics, "memory": memory_samples(directory / "raw")}
            rows.append(row)
            selected.append(row)
        results.append(aggregate(selected, workers, leg))
output = root / "clean-comparison"
output.mkdir(exist_ok=False)
value = {
    "runClass": "clean-score", "repeatsPerLeg": 2,
    "rustSourceCommit": "77ebe786ad1012bc9ac56105e08c05fb4c65b321",
    "goOptimizedBaseCommit": "b20b0fd00ac20c78bdb2d138fb6960b35e428af0",
    "binaryEncodingDisabled": True,
    "wireProof": json.loads((root / "provenance/wire-proof.json").read_text()),
    "aggregates": results, "repeats": rows,
    "latencyDefinition": "mean of per-repeat percentiles, not pooled percentiles",
    "scope": "clean performance checkpoint; separate diagnostic report is not complete",
}
(output / "summary.json").write_text(json.dumps(value, indent=2, sort_keys=True) + "\n")
columns = [name for name, item in results[0].items() if not isinstance(item, list)]
with (output / "summary.csv").open("w", newline="") as handle:
    writer = csv.DictWriter(handle, fieldnames=columns, extrasaction="ignore")
    writer.writeheader()
    writer.writerows(results)
table = "| Workers | SDK | Reads/s | CPU us/read | Busy cores | Mean repeat p50 ms | Mean repeat p99 ms | Mean RSS MiB |\n|---|---|---:|---:|---:|---:|---:|---:|\n"
for row in results:
    table += (f'| {row["workers"]} | {row["leg"]} | {row["throughput"]:.2f} | '
              f'{row["cpuMicrosecondsPerSuccess"]:.2f} | {row["busyCores"]:.3f} | '
              f'{row["meanRepeatP50Ms"]:.3f} | {row["meanRepeatP99Ms"]:.3f} | {row["meanRSSMiB"]:.2f} |\n')
(output / "report.md").write_text(
    "# Binary-disabled W32/W64 clean-score checkpoint\n\n" + table +
    "\nTwo mirrored clean repeats per SDK/worker count; 10 s excluded warmup and 60 s measurement. "
    "No perf, pprof, trace, or new hot-path diagnostic counters in these legs. "
    "Latency entries are means of repeat percentiles, not pooled values. "
    "Historical binary-wire results are not matched controls; this does not isolate one session-token commit.\n")
for name, field, title, unit in (
    ("throughput", "throughput", "Clean successful point reads/s", "Reads/s"),
    ("cpu-per-read", "cpuMicrosecondsPerSuccess", "Clean process CPU per successful read", "CPU microseconds/read"),
    ("latency-p99", "meanRepeatP99Ms", "Mean repeat p99 latency", "Milliseconds, not pooled p99"),
    ("rss", "meanRSSMiB", "Mean measurement process RSS", "MiB"),
):
    (output / (name + ".svg")).write_text(bars(results, field, title, unit))
body = "<table><tr><th>Workers</th><th>SDK</th><th>Reads/s</th><th>CPU us/read</th><th>Busy cores</th><th>p50 ms</th><th>p99 ms</th></tr>"
for row in results:
    body += (f'<tr><td>{row["workers"]}</td><td>{row["leg"]}</td><td>{row["throughput"]:.2f}</td>'
             f'<td>{row["cpuMicrosecondsPerSuccess"]:.2f}</td><td>{row["busyCores"]:.3f}</td>'
             f'<td>{row["meanRepeatP50Ms"]:.3f}</td><td>{row["meanRepeatP99Ms"]:.3f}</td></tr>')
body += "</table>" + "".join(f'<img src="{name}.svg" alt="{name}">' for name in ("throughput", "cpu-per-read", "latency-p99", "rss"))
(output / "dashboard.html").write_text(
    '<!doctype html><meta charset="utf-8"><title>Cosmos W32/W64 clean scores</title>'
    '<style>body{font:16px system-ui;max-width:1120px;margin:32px auto}table{width:100%;border-collapse:collapse}th,td{padding:8px;border-bottom:1px solid #ccd3df;text-align:right}img{width:100%}</style>'
    '<h1>Cosmos W32/W64: binary-disabled clean scores</h1><p>Two mirrored repeats per leg. '
    'CPU excludes warmup and profiling overhead; latency is mean repeat percentiles. Diagnostic captures are separate.</p>'
    + body)
with (output / "SHA256SUMS").open("x") as handle:
    for path in sorted(output.iterdir()):
        if path.name != "SHA256SUMS":
            handle.write(f"{hashlib.sha256(path.read_bytes()).hexdigest()}  {path.name}\n")
print("CLEAN_SCORE_SUMMARY_READY")
