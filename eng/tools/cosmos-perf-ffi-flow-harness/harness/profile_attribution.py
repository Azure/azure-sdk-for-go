#!/usr/bin/env python3
import collections
import csv
import html
import json
import math
import pathlib
import subprocess
import sys

import attribute
from run_leg import accepted_directory

COLORS = ("#2563eb", "#38bdf8", "#8b5cf6", "#10b981", "#14b8a6", "#f59e0b",
          "#f97316", "#ef4444", "#6366f1", "#a855f7", "#eab308", "#84cc16",
          "#64748b", "#475569", "#94a3b8", "#0f172a")


def flywheel(summary, title):
    shares = {category: summary["selfPercent"][category] * summary["profileBusyCores"] / 4
              for category in attribute.CATEGORIES}
    total = sum(shares.values())
    if not 0 < total <= 100:
        raise ValueError("invalid four-vCPU capacity bridge")
    radius, circumference = 112, 2 * math.pi * 112
    parts = [
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1080 520" role="img">',
        f'<title>{html.escape(title)}</title><rect width="100%" height="100%" fill="white"/>',
        f'<text x="24" y="32" font-family="sans-serif" font-size="20">{html.escape(title)}</text>',
        '<circle cx="170" cy="230" r="112" fill="none" stroke="#e2e8f0" stroke-width="48"/>',
    ]
    offset = 0
    for index, category in enumerate(attribute.CATEGORIES):
        value = shares[category]
        length = value / 100 * circumference
        parts.append(f'<circle cx="170" cy="230" r="{radius}" fill="none" stroke="{COLORS[index]}" stroke-width="48" '
                     f'stroke-dasharray="{length:.9f} {circumference-length:.9f}" stroke-dashoffset="{-offset:.9f}" '
                     'transform="rotate(-90 170 230)"/>')
        offset += length
        cpu_per_read = summary["selfPercent"][category] / 100 * summary["profileCPUUsPerRead"]
        y = 66 + index * 23
        parts.extend([
            f'<rect x="348" y="{y-10}" width="10" height="10" fill="{COLORS[index]}"/>',
            f'<text x="370" y="{y}" font-family="sans-serif" font-size="12">{html.escape(attribute.LABELS[category])}</text>',
            f'<text x="805" y="{y}" font-family="monospace" font-size="12">{value:.3f}% capacity</text>',
            f'<text x="962" y="{y}" font-family="monospace" font-size="12">{cpu_per_read:.2f} us/read</text>',
        ])
    parts.extend([
        f'<text x="170" y="227" text-anchor="middle" font-family="sans-serif" font-size="25">{total:.2f}%</text>',
        '<text x="170" y="252" text-anchor="middle" font-family="sans-serif" font-size="12">benchmark process / 4 vCPUs</text>',
        f'<text x="24" y="408" font-family="sans-serif" font-size="12">Gray remainder: {100-total:.2f}% capacity not charged to this process; includes idle or other work.</text>',
        f'<text x="24" y="432" font-family="sans-serif" font-size="12">Captured/process CPU coverage: {summary["sampledToProcessCPUPercent"]:.2f}%; physical Self totals normalize to 100%.</text>',
        '<text x="24" y="456" font-family="sans-serif" font-size="12">Profile-run estimates only, not clean-score timers or latency. Do not add perf and pprof denominators.</text>',
        '</svg>',
    ])
    return "\n".join(parts)


def main(root, output_directory="derived"):
    flamegraph = pathlib.Path("/mnt/cosmos-perf/point-read-first-pass-20261102-l/tools/FlameGraph")
    commit = subprocess.check_output(["git", "-C", str(flamegraph), "rev-parse", "HEAD"], text=True).strip()
    if subprocess.check_output(["git", "-C", str(flamegraph), "status", "--porcelain"], text=True).strip():
        raise ValueError("FlameGraph source is dirty")
    (root / "provenance/FlameGraph-commit.txt").write_text(commit + "\n")
    for workers in (32, 64):
        for leg in ("go-v2-native", "rust-sdk", "rust-driver"):
            directory = accepted_directory(root, workers, "perf", 1, leg)
            raw, derived = directory / "raw", directory / output_directory
            derived.mkdir(exist_ok=True)
            metrics = json.loads((directory / "metrics.json").read_text())
            quality = json.loads((directory / "derived/quality.json").read_text())
            if not quality["valid"] or quality["lostSamples"] != 0:
                raise ValueError("profile loss/quality gates failed")
            samples = attribute.parse_perf(raw / "perf-script.txt")
            symbols = attribute.demangle_map({"profile": samples})
            samples = attribute.remap_samples(samples, symbols)
            normalized_quality = dict(quality, lostEventPercent=quality["lostPercent"])
            summary, mapping = attribute.summarize(samples, metrics, normalized_quality)
            if not 90 <= summary["sampledToProcessCPUPercent"] <= 105:
                raise ValueError("captured CPU-clock/process CPU coverage is outside 90..105%")
            if not math.isclose(summary["selfPercentTotal"], 100, abs_tol=1e-8):
                raise ValueError("Self mapping drift")
            (derived / "attribution.json").write_text(json.dumps(summary, indent=2) + "\n")
            attribute.write_csv(derived / "sample-mapping.csv", mapping)
            leaves = collections.Counter()
            for sample in samples:
                leaf = sample.frames[0]
                leaves[(leaf.symbol, leaf.dso, attribute.classify(leaf))] += sample.period
            self_rows = [
                {"symbol": symbol, "dso": dso, "category": category,
                 "self_percent": period / summary["totalPeriodNs"] * 100}
                for (symbol, dso, category), period in leaves.most_common()
            ]
            if not math.isclose(sum(row["self_percent"] for row in self_rows), 100, abs_tol=1e-8):
                raise ValueError("full perf Self row accounting drift")
            attribute.write_csv(derived / "perf-self.csv", self_rows)
            title = f"W{workers} {leg}: physical CPU / four-vCPU capacity"
            (derived / "cpu-flywheel.svg").write_text(flywheel(summary, title))
            (derived / "cpu-flywheel.html").write_text(
                f'<!doctype html><meta charset="utf-8"><title>{html.escape(title)}</title>' +
                (derived / "cpu-flywheel.svg").read_text())
            stacks = collections.Counter()
            for sample in samples:
                stack = ";".join(frame.symbol.replace(";", ":") for frame in reversed(sample.frames))
                stacks[stack] += sample.period
            for name, divisor, countname in (
                ("whole-process", 1e6, "sampled CPU ms"),
                ("cpu-per-read", summary["totalPeriodNs"] / summary["profileCPUUsPerRead"], "estimated CPU us/read"),
            ):
                folded = derived / (name + ".folded")
                folded.write_text("".join(f"{stack} {period/divisor:.9f}\n" for stack, period in sorted(stacks.items())))
                with (derived / (name + ".svg")).open("w") as handle:
                    subprocess.run(
                        ["perl", str(flamegraph / "flamegraph.pl"), "--hash", "--title",
                         f"W{workers} {leg} {name}", "--countname", countname, str(folded)],
                        stdout=handle, check=True)
                (derived / (name + ".html")).write_text(
                    f'<!doctype html><meta charset="utf-8"><title>W{workers} {leg} {name}</title>' +
                    (derived / (name + ".svg")).read_text())
    print("PROFILE_SELF_ACCOUNTING_AND_FLYWHEELS_COMPLETE")


if __name__ == "__main__":
    main(pathlib.Path(sys.argv[1]), sys.argv[2] if len(sys.argv) > 2 else "derived")
