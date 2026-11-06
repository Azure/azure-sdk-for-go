#!/usr/bin/env python3
"""Render only a complete, validated binary-disabled three-leg matrix."""
import csv
import hashlib
import html
import json
import math
import pathlib
import statistics
import shutil
import sys

from run_leg import validate_metrics, digest, accepted_directory

LEGS = ("go-v2-native", "rust-sdk", "rust-driver")
COLORS = ("#2563eb", "#10b981", "#f59e0b")


def verify_raw(raw):
    seen = set()
    for row in (raw / "SHA256SUMS").read_text().splitlines():
        expected, separator, name = row.partition("  ")
        relative = pathlib.PurePosixPath(name)
        if (not separator or len(expected) != 64 or relative.is_absolute()
                or ".." in relative.parts or name in seen):
            raise ValueError("invalid or duplicate raw manifest path/hash")
        path = raw.joinpath(*relative.parts)
        if digest(path) != expected:
            raise ValueError(f"raw evidence hash mismatch: {name}")
        if path.stat().st_mode & 0o222:
            raise ValueError(f"raw evidence is not sealed read-only: {name}")
        seen.add(name)
    actual = {p.relative_to(raw).as_posix() for p in raw.rglob("*")
              if p.is_file() and p.name != "SHA256SUMS"}
    if seen != actual:
        raise ValueError("raw manifest does not cover the complete file set")


def memory_samples(raw):
    bounds = json.loads((raw / "barriers/measurement-bounds.json").read_text())
    values = []
    for line in (raw / "samples.jsonl").read_text().splitlines():
        value = json.loads(line)
        if not bounds["monotonicStartNs"] <= value["monotonicNs"] <= bounds["monotonicEndNs"]:
            continue
        fields = {}
        for row in value["status"].splitlines():
            name, _, content = row.partition(":")
            if name in ("VmRSS", "VmHWM", "RssAnon", "RssFile", "RssShmem"):
                fields[name] = int(content.split()[0]) * 1024
        if set(fields) != {"VmRSS", "VmHWM", "RssAnon", "RssFile", "RssShmem"}:
            raise ValueError("missing process RSS split in measurement samples")
        if abs(fields["VmRSS"] - fields["RssAnon"] - fields["RssFile"] - fields["RssShmem"]) > 4096:
            raise ValueError("RSS components do not reconcile")
        fields["cgroupMemoryCurrent"] = value["cgroupMemoryCurrent"]
        fields["cgroupMemoryPeak"] = value["cgroupMemoryPeak"]
        values.append(fields)
    if not values:
        raise ValueError("no measurement-aligned process memory samples")
    return {
        "sampleCount": len(values),
        "meanRSSMiB": statistics.mean(v["VmRSS"] for v in values) / 2**20,
        "peakSampledRSSMiB": max(v["VmRSS"] for v in values) / 2**20,
        "processLifetimeHWMMiB": max(v["VmHWM"] for v in values) / 2**20,
        "meanAnonMiB": statistics.mean(v["RssAnon"] for v in values) / 2**20,
        "meanFileMiB": statistics.mean(v["RssFile"] for v in values) / 2**20,
        "meanShmemMiB": statistics.mean(v["RssShmem"] for v in values) / 2**20,
        "meanCgroupCurrentMiB": statistics.mean(v["cgroupMemoryCurrent"] for v in values) / 2**20,
        "cgroupLifetimePeakMiB": max(v["cgroupMemoryPeak"] for v in values) / 2**20,
    }


def aggregate(rows, workers, leg):
    if len(rows) != 2 or any(row["workers"] != workers or row["leg"] != leg for row in rows):
        raise ValueError("exactly two matched clean repeats are required")
    wall = sum(row["metrics"]["timing"]["wallSeconds"] for row in rows)
    cpu = sum(row["metrics"]["timing"]["processCPUSeconds"] for row in rows)
    success = sum(row["metrics"]["counts"]["success"] for row in rows)
    if min(wall, cpu, success) <= 0 or cpu / wall > 4.02:
        raise ValueError("invalid aggregate timing or four-vCPU accounting")
    throughput = success / wall
    result = {
        "workers": workers, "leg": leg, "success": success,
        "wallSeconds": wall, "processCPUSeconds": cpu,
        "throughput": throughput, "busyCores": cpu / wall,
        "fourVCPUCapacityPercent": cpu / wall / 4 * 100,
        "cpuMicrosecondsPerSuccess": cpu * 1e6 / success,
        "meanRepeatP50Ms": statistics.mean(row["metrics"]["latencyMilliseconds"]["p50"] for row in rows),
        "meanRepeatP95Ms": statistics.mean(row["metrics"]["latencyMilliseconds"]["p95"] for row in rows),
        "meanRepeatP99Ms": statistics.mean(row["metrics"]["latencyMilliseconds"]["p99"] for row in rows),
        "repeatThroughput": [row["metrics"]["counts"]["success"] / row["metrics"]["timing"]["wallSeconds"] for row in rows],
    }
    for name in rows[0]["memory"]:
        if name != "sampleCount":
            result[name] = (max(row["memory"][name] for row in rows)
                            if name.startswith(("peak", "processLifetime", "cgroupLifetime"))
                            else statistics.mean(row["memory"][name] for row in rows))
    result["repeatThroughputSpreadPercent"] = (
        max(result["repeatThroughput"]) - min(result["repeatThroughput"])
    ) / throughput * 100
    if not all(math.isfinite(value) for value in result.values() if isinstance(value, float)):
        raise ValueError("nonfinite report metric")
    return result


def bars(rows, field, title, unit):
    width, height = 960, 440
    maximum = max(row[field] for row in rows) * 1.15
    if maximum <= 0:
        raise ValueError("cannot render nonpositive chart scale")
    parts = [
        f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {width} {height}" role="img">',
        f'<title>{html.escape(title)}</title><rect width="100%" height="100%" fill="#fff"/>',
        f'<text x="24" y="34" font-size="21" font-family="sans-serif">{html.escape(title)}</text>',
    ]
    for index, row in enumerate(rows):
        x = 70 + index * 145
        size = row[field] / maximum * 290
        color = COLORS[LEGS.index(row["leg"])]
        parts.extend([
            f'<rect x="{x}" y="{350-size:.3f}" width="100" height="{size:.3f}" fill="{color}"/>',
            f'<text x="{x+50}" y="{338-size:.3f}" text-anchor="middle" font-family="sans-serif" font-size="13">{row[field]:,.2f}</text>',
            f'<text x="{x+50}" y="376" text-anchor="middle" font-family="sans-serif" font-size="12">{row["leg"]}</text>',
            f'<text x="{x+50}" y="395" text-anchor="middle" font-family="sans-serif" font-size="12">W{row["workers"]}</text>',
        ])
    parts.append(f'<text x="24" y="425" font-size="12" font-family="sans-serif">{html.escape(unit)}; clean scores only; binary disabled; latest pinned Rust main</text></svg>')
    return "\n".join(parts)


def flow_chart(flow, queue=False):
    buckets = flow["buckets"]
    if not buckets:
        raise ValueError("empty flow timeline")
    if queue:
        series = [
            ("Ready-queue high-water", [b["nativeReadyHigh"] or 0 for b in buckets], "#2563eb"),
            ("Admitted-undrained high-water", [b["nativeAdmittedUndrainedHigh"] or 0 for b in buckets], "#10b981"),
            ("First-polled-unfinished high-water", [b["nativePolledUnfinishedHigh"] or 0 for b in buckets], "#f59e0b"),
        ]
        title = f'W{flow["workers"]}: actual native queue/work gauges'
    else:
        series = [
            ("Go submissions/s", [b["goRates"].get("submitAttempts", 0) for b in buckets], "#2563eb"),
            ("Rust terminal/s", [b["nativeRates"].get("terminal", 0) for b in buckets], "#10b981"),
            ("Go completions/s", [b["goRates"].get("drained", 0) for b in buckets], "#f59e0b"),
            ("Go waiter resumptions/s", [b["goRates"].get("resumed", 0) for b in buckets], "#8b5cf6"),
        ]
        title = f'W{flow["workers"]}: measurement-only FFI flow rates'
    maximum = max(1, max(max(values) for _, values, _ in series)) * 1.1
    parts = [
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1000 420" role="img">',
        f'<title>{html.escape(title)}</title><rect width="100%" height="100%" fill="white"/>',
        f'<text x="24" y="28" font-family="sans-serif" font-size="20">{html.escape(title)}</text>',
    ]
    for index, (name, values, color) in enumerate(series):
        points = " ".join(f"{60+i/max(1,len(values)-1)*880:.3f},{320-value/maximum*230:.3f}"
                          for i, value in enumerate(values))
        parts.append(f'<polyline fill="none" stroke="{color}" stroke-width="2" points="{points}"/>')
        parts.append(f'<text x="{24+index*245}" y="62" font-size="12" font-family="sans-serif" fill="{color}">{html.escape(name)}</text>')
    parts.extend([
        '<text x="60" y="346" font-size="12" font-family="sans-serif">Measurement second buckets, chronological order</text>',
        '<text x="24" y="376" font-size="12" font-family="sans-serif">Partial first/last seconds use actual duration; per-bucket arrivals and completions need not balance.</text>',
        '<text x="24" y="399" font-size="12" font-family="sans-serif">Diagnostic run only. Polled-unfinished includes awaiting I/O; queue/batch size is not CPU utilization.</text>',
        '</svg>',
    ])
    return "\n".join(parts)


def main(root, output_name="comparison-v1", attribution_directory="derived"):
    spec = json.loads((root / "run-spec.json").read_text())
    wire = json.loads((root / "provenance/wire-proof.json").read_text())
    if wire["rustSourceCommit"] != spec["rustSourceCommit"]:
        raise ValueError("wire proof source differs from campaign")
    rows, aggregates, flows, profiles, rust_flows, schedulers, go_diagnostics = [], [], [], [], [], [], []
    expected_workload = None
    for workers in (32, 64):
        for leg in LEGS:
            selected = []
            for repeat in (1, 2):
                directory = accepted_directory(root, workers, "clean-score", repeat, leg)
                valid = json.loads((directory / "VALID.json").read_text())
                if not valid["valid"] or (directory / "REJECTED.json").exists():
                    raise ValueError("clean-score leg is rejected")
                metrics = json.loads((directory / "metrics.json").read_text())
                verify_raw(directory / "raw")
                if digest(directory / "raw/binary") != valid["binarySHA256"]:
                    raise ValueError("clean binary hash no longer matches validation record")
                validate_metrics(metrics, workers)
                if metrics["workload"]["warmupSeconds"] != spec["warmupSeconds"] or metrics["workload"]["measurementSeconds"] != spec["measurementSeconds"]:
                    raise ValueError("warmup/measurement mismatch")
                contract = {name: metrics["workload"][name] for name in (
                    "partitionKeyPath", "seedCount", "documentBytes", "operation", "keySequence",
                    "closedLoop", "binaryEncoding", "authentication")}
                if expected_workload is None:
                    expected_workload = contract
                if contract != expected_workload:
                    raise ValueError("cross-leg effective workload/authentication mismatch")
                row = {"workers": workers, "leg": leg, "repeat": repeat, "metrics": metrics,
                       "memory": memory_samples(directory / "raw")}
                rows.append(row)
                selected.append(row)
            aggregates.append(aggregate(selected, workers, leg))
            directory = accepted_directory(root, workers, "perf", 1, leg)
            valid = json.loads((directory / "VALID.json").read_text())
            quality = json.loads((directory / "derived/quality.json").read_text())
            if not valid["valid"] or not quality["valid"]:
                raise ValueError("missing or failed profile-quality gate")
            verify_raw(directory / "raw")
            profiles.append({
                "workers": workers, "leg": leg, "quality": quality,
                "attribution": json.loads((directory / attribution_directory / "attribution.json").read_text()),
                "attributionDirectory": attribution_directory,
                "path": str(directory.relative_to(root).as_posix()),
            })
            if not math.isclose(profiles[-1]["attribution"]["selfPercentTotal"], 100, abs_tol=1e-8):
                raise ValueError("profile Self attribution does not normalize")
            if leg != "go-v2-native":
                rust_dir = accepted_directory(root, workers, "rust-flow", 1, leg)
                verify_raw(rust_dir / "raw")
                rates = json.loads((rust_dir / "derived/operation-rates.json").read_text())
                if not rates["valid"]:
                    raise ValueError("invalid standalone Rust operation-rate ledger")
                rust_flows.append(dict(rates, workers=workers, leg=leg))
        flow_dir = accepted_directory(root, workers, "ffi-flow", 1, "go-v2-native")
        flow = json.loads((flow_dir / "derived/flow.json").read_text())
        if not flow["valid"]:
            raise ValueError("invalid flow ledger")
        verify_raw(flow_dir / "raw")
        diagnostic = json.loads((flow_dir / "metrics.json").read_text())
        reference = next(row for row in aggregates if row["workers"] == workers and row["leg"] == "go-v2-native")
        throughput = diagnostic["counts"]["success"] / diagnostic["timing"]["wallSeconds"]
        cpu_per_read = diagnostic["timing"]["processCPUSeconds"] * 1e6 / diagnostic["counts"]["success"]
        flow.update(
            workers=workers, diagnosticThroughput=throughput,
            diagnosticCPUPerReadUs=cpu_per_read,
            throughputDifferenceFromCleanPercent=(throughput / reference["throughput"] - 1) * 100,
            cpuPerReadDifferenceFromCleanPercent=(cpu_per_read / reference["cpuMicrosecondsPerSuccess"] - 1) * 100,
            taxCaveat="One diagnostic leg versus two clean repeats; differences include temporal variation.",
        )
        flows.append(flow)
        diagnostic_dir = accepted_directory(root, workers, "go-diagnostic", 1, "go-v2-native")
        diagnostic_valid = json.loads((diagnostic_dir / "VALID.json").read_text())
        if not diagnostic_valid["valid"]:
            raise ValueError("Go pprof/trace diagnostic leg is invalid")
        verify_raw(diagnostic_dir / "raw")
        for name in ("cpu.pprof", "heap.pprof", "runtime.trace", "block.pprof", "mutex.pprof"):
            if (diagnostic_dir / "raw" / name).stat().st_size == 0:
                raise ValueError(f"empty Go diagnostic profile: {name}")
        go_diagnostics.append(json.loads((diagnostic_dir / "derived/go-diagnostics-summary.json").read_text()))
        scheduler_dir = accepted_directory(root, workers, "scheduler", 1, "go-v2-native")
        verify_raw(scheduler_dir / "raw")
        schedulers.append(dict(json.loads((scheduler_dir / "derived/scheduler-scope.json").read_text()),
                               workers=workers, path=scheduler_dir.relative_to(root).as_posix()))
    cleanup = json.loads((root / "provenance/cleanup.json").read_text())
    if cleanup["deleted"] != 1024 or cleanup["nonOwnedSkipped"] != 0:
        raise ValueError("campaign-owned cleanup did not reconcile")
    output = root / output_name
    output.mkdir(exist_ok=False)
    summary = {
        "spec": spec, "wireProof": wire, "cleanAggregates": aggregates,
        "cleanRepeats": rows, "flowDiagnostics": flows, "profiles": profiles,
        "rustFlowDiagnostics": rust_flows, "schedulerEvidence": schedulers,
        "goDiagnostics": go_diagnostics,
        "limits": [
            "Historical binary-wire/text-output scores are not matched binary-disabled controls.",
            "A full Rust-main upgrade cannot establish the isolated effect of one session-token change.",
            "Latency columns are means of repeat percentiles, not pooled percentiles.",
            "Completion batch size and actual ready-queue depth are distinct measurements.",
            "Go pprof and Linux perf retain separate CPU denominators.",
            "RSS/HWM includes native, Go, mapped libraries, and stacks; heap profiles are not process RSS.",
        ],
    }
    (output / "summary.json").write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n")
    for profile in profiles:
        derived = root / profile["path"] / profile["attributionDirectory"]
        stem = f'w{profile["workers"]}-{profile["leg"]}'
        shutil.copyfile(derived / "cpu-flywheel.svg", output / (stem + "-cpu-flywheel.svg"))
    for flow in flows:
        for queue, name in ((False, "flow-rates"), (True, "queue-gauges")):
            (output / f'w{flow["workers"]}-{name}.svg').write_text(flow_chart(flow, queue))
    columns = [name for name, value in aggregates[0].items() if not isinstance(value, list)]
    with (output / "summary.csv").open("w", newline="") as handle:
        writer = csv.DictWriter(handle, fieldnames=columns, extrasaction="ignore")
        writer.writeheader()
        writer.writerows(aggregates)
    for name, field, title, unit in (
        ("throughput", "throughput", "Successful point reads per second", "Reads/s"),
        ("cpu-per-read", "cpuMicrosecondsPerSuccess", "Process CPU per successful read", "Microseconds of process CPU/read"),
        ("latency-p99", "meanRepeatP99Ms", "Mean repeat p99 latency", "Milliseconds, not pooled p99"),
        ("memory-rss", "meanRSSMiB", "Measurement-aligned mean process RSS", "MiB; total process memory, not heap-only"),
    ):
        (output / (name + ".svg")).write_text(bars(aggregates, field, title, unit))
    heading = "| Workers | Leg | Reads/s | CPU us/read | Busy cores | Mean repeat p99 ms | Mean RSS MiB |\n|---|---|---:|---:|---:|---:|---:|"
    table = "\n".join(
        f'| {r["workers"]} | {r["leg"]} | {r["throughput"]:,.2f} | {r["cpuMicrosecondsPerSuccess"]:.2f} | {r["busyCores"]:.3f} | {r["meanRepeatP99Ms"]:.3f} | {r["meanRSSMiB"]:.2f} |'
        for r in aggregates)
    report = ("# Binary-disabled Cosmos W32/W64 comparison\n\n"
              f'Rust source: `{spec["rustSourceCommit"]}`. Binary encoding explicitly disabled; '
              "pre-conversion text-wire evidence is required for every leg.\n\n"
              + heading + "\n" + table + "\n\n"
              + "\n".join("- " + limit for limit in summary["limits"]) + "\n")
    (output / "report.md").write_text(report)
    cells = "".join(
        f'<tr><td>{r["workers"]}</td><td>{r["leg"]}</td><td>{r["throughput"]:,.2f}</td>'
        f'<td>{r["cpuMicrosecondsPerSuccess"]:.2f}</td><td>{r["busyCores"]:.3f}</td>'
        f'<td>{r["meanRepeatP99Ms"]:.3f}</td><td>{r["meanRSSMiB"]:.2f}</td></tr>'
        for r in aggregates)
    charts = "".join(f'<img src="{name}.svg" alt="{name}">' for name in
                     ("throughput", "cpu-per-read", "latency-p99", "memory-rss"))
    charts += "<h2>Physical CPU attribution: separate profile runs</h2>" + "".join(
        f'<img src="w{profile["workers"]}-{profile["leg"]}-cpu-flywheel.svg" alt="CPU attribution W{profile["workers"]} {profile["leg"]}">'
        for profile in profiles)
    charts += "<h2>FFI flow: separate instrumented runs</h2>" + "".join(
        f'<img src="w{flow["workers"]}-{name}.svg" alt="W{flow["workers"]} {name}">'
        for flow in flows for name in ("flow-rates", "queue-gauges"))
    page = (
        '<!doctype html><meta charset="utf-8"><title>Cosmos binary-disabled W32/W64</title>'
        '<style>body{font:16px system-ui;max-width:1180px;margin:32px auto;color:#17233b}table{border-collapse:collapse;width:100%}td,th{padding:10px;border-bottom:1px solid #ccd3df;text-align:right}img{width:100%;margin:12px 0}.notice{padding:16px;background:#fff2c9}code{overflow-wrap:anywhere}</style>'
        '<h1>Cosmos: binary-disabled W32/W64</h1>'
        '<p class="notice">All three legs require verified text wire bodies before conversion. '
        'Older binary-wire campaigns are not matched controls. This does not isolate a single session-token commit.</p>'
        f'<p>Rust source <code>{spec["rustSourceCommit"]}</code>; Go optimizations <code>{spec["goOptimizedSourceCommit"]}</code>.</p>'
        '<table><tr><th>Workers</th><th>Leg</th><th>Reads/s</th><th>CPU us/read</th><th>Busy cores</th><th>Mean repeat p99 ms</th><th>Mean RSS MiB</th></tr>'
        + cells + "</table>" + charts
        + "<h2>Interpretation limits</h2><ul>"
        + "".join("<li>" + html.escape(limit) + "</li>" for limit in summary["limits"])
        + '</ul><p>Exact repeats, wire proof, profile links, and flow ledgers: <a href="summary.json">summary.json</a>.</p>')
    (output / "dashboard.html").write_text(page)
    with (output / "SHA256SUMS").open("x") as handle:
        for path in sorted(output.iterdir()):
            if path.name != "SHA256SUMS":
                handle.write(f"{hashlib.sha256(path.read_bytes()).hexdigest()}  {path.name}\n")
    print("MATCHED_THREE_LEG_REPORT_COMPLETE")


if __name__ == "__main__":
    main(pathlib.Path(sys.argv[1]), sys.argv[2] if len(sys.argv) > 2 else "comparison-v1",
         sys.argv[3] if len(sys.argv) > 3 else "derived")
