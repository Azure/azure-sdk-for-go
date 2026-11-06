#!/usr/bin/env python3
import json
import os
import pathlib
import subprocess
import sys

from run_leg import accepted_directory

root = pathlib.Path(sys.argv[1])
go = "/usr/local/go/bin/go"
environment = os.environ.copy()
environment["HOME"] = "/root"
environment["GOCACHE"] = str(root / "build/go-cache")
for workers in (32, 64):
    directory = accepted_directory(root, workers, "go-diagnostic", 1, "go-v2-native")
    raw, derived = directory / "raw", directory / "derived"
    binary = raw / "binary"
    for kind in ("sched", "sync", "syscall", "net"):
        profile = derived / f"trace-{kind}.pprof"
        with profile.open("wb") as output:
            subprocess.run([go, "tool", "trace", "-pprof=" + kind, str(raw / "runtime.trace")],
                           stdout=output, check=True, env=environment)
        with (derived / f"trace-{kind}-top.txt").open("w") as output:
            subprocess.run([go, "tool", "pprof", "-top", "-nodecount=30", str(binary), str(profile)],
                           stdout=output, check=True, env=environment)
    for name, profile, index in (
        ("cpu", "cpu.pprof", None),
        ("heap-live", "heap.pprof", "inuse_space"),
        ("heap-alloc-bytes", "heap.pprof", "alloc_space"),
        ("heap-alloc-objects", "heap.pprof", "alloc_objects"),
        ("block", "block.pprof", None),
        ("mutex", "mutex.pprof", None),
    ):
        command = [go, "tool", "pprof", "-top", "-nodecount=30"]
        if index:
            command.append("-sample_index=" + index)
        command += [str(binary), str(raw / profile)]
        with (derived / f"{name}-top.txt").open("w") as output:
            subprocess.run(command, stdout=output, check=True, env=environment)
    (derived / "go-diagnostics-summary.json").write_text(json.dumps({
        "workers": workers, "traceWindowSeconds": 5,
        "capture": "separate Go diagnostic leg, not a clean score",
        "schedulerProfile": "trace-sched.pprof",
        "blockingProfiles": ["trace-sync.pprof", "trace-syscall.pprof", "block.pprof", "mutex.pprof"],
        "allocationHotspots": ["heap-alloc-bytes-top.txt", "heap-alloc-objects-top.txt"],
        "gcEvidence": "runtime.trace for GC pauses/mark assists and cpu-top.txt for sampled Go runtime CPU",
        "limits": [
            "Heap allocation totals include process startup/warmup; not measurement-only bytes per operation.",
            "Go heap profiles do not include native Rust/C allocations.",
            "Go pprof, perf, and process CPU accounting are separate denominators.",
        ],
    }, indent=2) + "\n")
print("GO_SCHEDULER_ALLOCATION_DIAGNOSTICS_ANALYZED")
