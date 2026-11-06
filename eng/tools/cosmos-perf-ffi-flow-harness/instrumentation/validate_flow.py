#!/usr/bin/env python3
"""Reconcile measurement-only Go/native flow ledgers, without adding profiling denominators."""
import json
import pathlib
import sys


def require(condition, message):
    if not condition:
        raise ValueError(message)


def totals(rows):
    result = {}
    for row in rows:
        for name, value in row["counts"].items():
            require(isinstance(value, int) and value >= 0, f"invalid counter {name}")
            result[name] = result.get(name, 0) + value
    return result


def validate(go, native, metrics, workers):
    require(go["overflow"] == native["overflow"] == 0, "flow clock moved backwards or capacity overflowed")
    require(native["activeAtSnapshot"] == 0, "unfinished Rust tasks at final snapshot")
    require(native["readyAtSnapshot"] == native["admittedUndrainedAtSnapshot"] == 0,
            "native queue/outstanding work not drained at final snapshot")
    gr = [r for r in go["rows"] if r["phase"] == 2]
    nr = [r for r in native["rows"] if r["phase"] == 2]
    g, n = totals(gr), totals(nr)
    require(gr and nr, "missing measured flow buckets")
    require(metrics["counts"]["error"] == metrics["counts"]["retry"] == metrics["counts"]["throttle"] == 0,
            "invalid measured adapter counts")
    success = metrics["counts"]["success"]
    require(success > 0, "no measured operations")
    require(g["submitAttempts"] == g["submitAccepted"] == success,
            "Go submit totals do not reconcile to successful measured operations")
    require(g.get("submitErrors", 0) == 0, "Go submission errors")
    require(g["drained"] == g["delivered"] == g["resumed"] == success,
            "Go drained/delivered/resumed completions drift")
    require(n["admitted"] == n["task_submitted"] == n["first_poll"] == n["terminal"] == n["enqueued"] == n["drained"] == success,
            "native lifecycle does not reconcile to successful measured operations")
    require(n.get("admission_aborted", 0) == n.get("enqueue_rejected", 0) == 0,
            "native aborted admission or lost completion")
    hist = {}
    for row in gr:
        for size, batches in row["batchHistogram"].items():
            size = int(size)
            require(1 <= size <= 64 and batches >= 0, "invalid batch histogram")
            hist[size] = hist.get(size, 0) + batches
    nonempty = sum(hist.values())
    require(sum(size * count for size, count in hist.items()) == success, "batch/completion ledger drift")
    require(g["waitReturns"] == nonempty + g.get("emptyWaits", 0), "wait/nonempty/empty ledger drift")
    ready_high = max(r["readyHigh"] for r in nr)
    outstanding_high = max(r["admittedUndrainedHigh"] for r in nr)
    active_high = max(r["polledUnfinishedHigh"] for r in nr)
    require(ready_high <= workers and outstanding_high <= workers and active_high <= workers,
            "native high-water exceeds strict closed-loop worker cap")
    start = native["phaseStartUnixNs"][2]
    end = native["phaseEndUnixNs"][2]
    require(end > start > 0, "missing native measurement phase boundaries")
    require(go["phaseStartUnixNs"][2] == start and go["phaseEndUnixNs"][2] == end,
            "Go/native measurement phases do not share exact boundaries")
    # Events crossing a Unix-second boundary need not balance within one row.
    # Whole measurement must balance; per-bucket rates retain the actual event timestamps.
    elapsed = (end - start) / 1e9
    adapter_wall = metrics["timing"]["wallSeconds"]
    require(abs(elapsed - adapter_wall) <= max(0.02, adapter_wall * 0.001),
            "shared wall-clock phase duration drifted from adapter monotonic timing")
    rows = []
    by_g = {r["unixSecond"]: r for r in gr}
    by_n = {r["unixSecond"]: r for r in nr}
    for second in sorted(set(by_g) | set(by_n)):
        duration = max(0, min(end, (second + 1) * 10**9) - max(start, second * 10**9)) / 1e9
        require(duration > 0, "bucket event outside measurement boundaries")
        gg = by_g.get(second, {"counts": {}, "batchHistogram": {}})
        nn = by_n.get(second, {"counts": {}})
        rows.append({
            "unixSecond": second, "durationSeconds": duration,
            "goCounts": gg["counts"], "nativeCounts": nn["counts"],
            "goRates": {k: v / duration for k, v in gg["counts"].items()},
            "nativeRates": {k: v / duration for k, v in nn["counts"].items()},
            "batchHistogram": gg["batchHistogram"],
            "nativeReadyHigh": nn.get("readyHigh"),
            "nativeAdmittedUndrainedHigh": nn.get("admittedUndrainedHigh"),
            "nativePolledUnfinishedHigh": nn.get("polledUnfinishedHigh"),
            "nativeReadyLast": nn.get("readyLast"),
            "nativeAdmittedUndrainedLast": nn.get("admittedUndrainedLast"),
        })
    return {
        "valid": True, "success": success, "flowWallSeconds": elapsed,
        "adapterWallSeconds": adapter_wall,
        "goTotals": g, "nativeTotals": n,
        "submissionsPerSecond": success / elapsed,
        "completionBatchesPerSecond": nonempty / elapsed,
        "completionsPerSecond": success / elapsed,
        "meanCompletionBatchSize": success / nonempty,
        "batchHistogram": hist, "nativeReadyHigh": ready_high,
        "deliveryToResumeMeanUs": sum(row["deliveryToResumeTotalNs"] for row in gr) / success / 1000,
        "deliveryToResumeMaxUs": max(row["deliveryToResumeMaxNs"] for row in gr) / 1000,
        "nativeAdmittedUndrainedHigh": outstanding_high, "nativePolledUnfinishedHigh": active_high,
        "buckets": rows,
        "interpretation": "batch size is not backlog depth; closed-loop workers bound caller-visible concurrency",
    }


if __name__ == "__main__":
    root = pathlib.Path(sys.argv[1])
    value = validate(*(json.loads((root / name).read_text()) for name in
                       ("go-flow.json", "native-flow.json", "adapter-result.json")), int(sys.argv[2]))
    print(json.dumps(value, indent=2))
