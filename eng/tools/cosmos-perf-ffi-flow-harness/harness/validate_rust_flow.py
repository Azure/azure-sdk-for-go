#!/usr/bin/env python3


def validate(value):
    source = value["operationBuckets"]
    if not isinstance(source, dict) or not source:
        raise ValueError("missing Rust per-second completed-operation counters")
    start = value["timing"]["measurementStartUnixNs"]
    end = value["timing"]["measurementEndUnixNs"]
    if not end > start > 0:
        raise ValueError("invalid Rust measurement boundaries")
    rows = []
    for key, count in sorted(source.items(), key=lambda row: int(row[0])):
        second = int(key)
        if not isinstance(count, int) or count < 0:
            raise ValueError("invalid Rust completion count")
        duration = (min(end, (second + 1) * 10**9) - max(start, second * 10**9)) / 1e9
        if duration <= 0:
            raise ValueError("Rust counter bucket outside measured phase")
        rows.append({"unixSecond": second, "durationSeconds": duration,
                     "completed": count, "completedPerSecond": count / duration})
    total = sum(row["completed"] for row in rows)
    if total != value["counts"]["success"]:
        raise ValueError("Rust per-second completions do not reconcile to successful operations")
    wall = (end - start) / 1e9
    expected = value["timing"]["wallSeconds"]
    if abs(wall - expected) > max(0.02, expected * 0.001):
        raise ValueError("Rust wall-clock bucket duration drifted from monotonic measurement")
    return {"valid": True, "success": total, "wallSeconds": wall,
            "operationsPerSecond": total / wall, "buckets": rows,
            "definition": "completed SDK/direct-driver point reads; no FFI or completion queue in these Rust legs"}
