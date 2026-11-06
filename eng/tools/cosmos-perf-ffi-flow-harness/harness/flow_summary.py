#!/usr/bin/env python3
import hashlib
import json
import pathlib
import sys

from run_leg import accepted_directory
from report import flow_chart, verify_raw

root = pathlib.Path(sys.argv[1])
output = root / "ffi-comparison"
output.mkdir(exist_ok=False)
values = []
for workers in (32, 64):
    directory = accepted_directory(root, workers, "ffi-flow", 1, "go-v2-native")
    verify_raw(directory / "raw")
    value = json.loads((directory / "derived/flow.json").read_text())
    assert value["valid"]
    value["workers"] = workers
    metrics = json.loads((directory / "metrics.json").read_text())
    value["diagnosticCPUUsPerRead"] = metrics["timing"]["processCPUSeconds"] * 1e6 / metrics["counts"]["success"]
    values.append(value)
    (output / f"w{workers}-flow.json").write_text(json.dumps(value, indent=2) + "\n")
    (output / f"w{workers}-flow-rates.svg").write_text(flow_chart(value))
    (output / f"w{workers}-queue-gauges.svg").write_text(flow_chart(value, True))
(output / "summary.json").write_text(json.dumps({
    "runClass": "ffi-flow diagnostic, not clean-score", "goV2": values,
    "scope": "actual completion deque depth; admitted-undrained includes work awaiting I/O",
    "histogram": "keys are completion count per nonempty wait; values are number of such waits",
}, indent=2) + "\n")
table = "| Workers | Submissions/s | Batches/s | Completions/s | Mean batch | Ready queue max | Admitted-undrained max | Mean delivery-resume us |\n|---|---:|---:|---:|---:|---:|---:|---:|\n"
for value in values:
    table += (f'| {value["workers"]} | {value["submissionsPerSecond"]:.2f} | {value["completionBatchesPerSecond"]:.2f} | '
              f'{value["completionsPerSecond"]:.2f} | {value["meanCompletionBatchSize"]:.3f} | '
              f'{value["nativeReadyHigh"]} | {value["nativeAdmittedUndrainedHigh"]} | {value["deliveryToResumeMeanUs"]:.3f} |\n')
(output / "report.md").write_text("# Go v2 FFI-flow diagnostics\n\n" + table +
    "\nSeparate instrumented runs. Batch size is not backlog depth. Delivery-to-resume includes publication and Go scheduling, not pure scheduler time.\n")
body = '<!doctype html><meta charset="utf-8"><title>Go v2 completion queue diagnostics</title><style>body{font:16px system-ui;max-width:1100px;margin:32px auto}img{width:100%}</style><h1>Go v2 FFI-flow diagnostics</h1><p>Separate instrumented runs, not authoritative clean throughput scores. Batch counts and actual queue depth are different.</p>'
for value in values:
    workers = value["workers"]
    body += f'<h2>W{workers}</h2><a href="w{workers}-flow.json">Exact decoded counters and batch histogram</a>'
    body += f'<img src="w{workers}-flow-rates.svg" alt="Flow rates"><img src="w{workers}-queue-gauges.svg" alt="Actual queue gauges">'
(output / "dashboard.html").write_text(body)
with (output / "SHA256SUMS").open("x") as handle:
    for path in sorted(output.iterdir()):
        if path.name != "SHA256SUMS":
            handle.write(f"{hashlib.sha256(path.read_bytes()).hexdigest()}  {path.name}\n")
print("FFI_COUNTER_CHECKPOINT_READY")
