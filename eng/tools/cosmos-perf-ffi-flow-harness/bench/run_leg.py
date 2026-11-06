#!/usr/bin/env python3
"""Run one isolated Linux leg, with phase barriers and measurement-filtered perf."""
import argparse
import hashlib
import json
import os
import pathlib
import signal
import subprocess
import time


def write_json(path, value):
    with path.open("x") as handle:
        json.dump(value, handle, indent=2)
        handle.write("\n")


def digest(path):
    value = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            value.update(chunk)
    return value.hexdigest()


def perf_time(nanoseconds):
    return f"{nanoseconds // 10**9}.{nanoseconds % 10**9:09d}"


def memory_value(cgroup, peak=False):
    choices = ("memory.peak", "memory.max_usage_in_bytes") if peak else ("memory.current", "memory.usage_in_bytes")
    for name in choices:
        path = cgroup / name
        if path.exists():
            return int(path.read_text())
    raise RuntimeError("dedicated cgroup exposes no supported memory accounting")


def memory_root():
    root = pathlib.Path("/sys/fs/cgroup")
    return root if (root / "cgroup.controllers").exists() else root / "memory"


def seal(raw):
    manifest = raw / "SHA256SUMS"
    with manifest.open("x") as handle:
        for path in sorted(raw.rglob("*")):
            if path.is_file() and path != manifest:
                handle.write(f"{digest(path)}  {path.relative_to(raw).as_posix()}\n")
    for path in raw.rglob("*"):
        if path.is_file():
            path.chmod(0o555 if path.name == "binary" else 0o444)


def accepted_directory(root, workers, run_class, repeat, leg):
    key = f"w{workers}/{run_class}/repeat-{repeat}/{leg}"
    index_path = root / "provenance/accepted-legs.json"
    index = json.loads(index_path.read_text()) if index_path.exists() else {}
    relative = pathlib.PurePosixPath(index.get(key, "run/" + key))
    if relative.is_absolute() or ".." in relative.parts:
        raise ValueError("invalid accepted-leg evidence path")
    return root.joinpath(*relative.parts)


def record_accepted(root, key, directory):
    path = root / "provenance/accepted-legs.json"
    value = json.loads(path.read_text()) if path.exists() else {}
    value[key] = directory.relative_to(root).as_posix()
    temporary = path.with_suffix(".next")
    temporary.write_text(json.dumps(value, indent=2, sort_keys=True) + "\n")
    temporary.replace(path)


def snapshot(pid, raw, name, cgroup):
    started = time.monotonic_ns()
    proc = pathlib.Path("/proc") / str(pid)
    for source, target in (("status", "status"), ("smaps_rollup", "smaps-rollup"), ("smaps", "smaps")):
        (raw / f"{target}-{name}.txt").write_bytes((proc / source).read_bytes())
    stacks = []
    threads = {}
    for thread in sorted((proc / "task").iterdir(), key=lambda path: int(path.name)):
        try:
            stack = (thread / "stack").read_text()
        except PermissionError:
            stack = "UNAVAILABLE: kernel stack permission denied\n"
        except FileNotFoundError:
            stack = "UNAVAILABLE: thread exited during snapshot\n"
        stacks.append(f"=== TID {thread.name} ===\n{stack}")
        try:
            threads[thread.name] = {
                "status": (thread / "status").read_text(),
                "schedstat": (thread / "schedstat").read_text(),
            }
        except FileNotFoundError:
            threads[thread.name] = {"unavailable": "thread exited during snapshot"}
    (raw / f"thread-stacks-{name}.txt").write_text("".join(stacks))
    write_json(raw / f"snapshot-{name}.json", {
        "monotonicStartNs": started, "monotonicEndNs": time.monotonic_ns(),
        "unixNs": time.time_ns(),
        "cgroupMemoryCurrent": memory_value(cgroup),
        "cgroupMemoryPeak": memory_value(cgroup, peak=True),
        "cgroupIsDedicatedToLeg": True,
        "threadObservations": threads,
        "schedstatsEnabled": pathlib.Path("/proc/sys/kernel/sched_schedstats").read_text().strip()
            if pathlib.Path("/proc/sys/kernel/sched_schedstats").exists() else None,
    })
    (raw / f"cgroup-memory-stat-{name}.txt").write_bytes((cgroup / "memory.stat").read_bytes())


def perf_armed(process):
    if process.poll() is not None:
        raise RuntimeError(f"perf exited before arming: {process.returncode}")
    try:
        paths = list((pathlib.Path("/proc") / str(process.pid) / "fd").iterdir())
    except FileNotFoundError:
        return False
    for path in paths:
        try:
            if "perf_event" in os.readlink(path):
                return True
        except FileNotFoundError:
            continue
    return False


def validate_metrics(value, workers):
    counts = value["counts"]
    if counts["success"] <= 0 or any(counts[key] for key in ("error", "retry", "throttle")):
        raise ValueError(f"invalid measured counts: {counts}")
    workload = value["workload"]
    if workload["workers"] != workers or not workload["closedLoop"] or workload["binaryEncoding"] is not False:
        raise ValueError("worker, closed-loop, or encoding configuration mismatch")
    for key, expected in (
        ("partitionKeyPath", "/pk"), ("seedCount", 1024), ("documentBytes", 1024),
        ("operation", "point-read"), ("keySequence", "round-robin-0-through-1023"),
        ("rawResponseMode", True), ("validationOutsideTimedPath", True),
    ):
        if workload[key] != expected:
            raise ValueError(f"workload mismatch: {key}")
    proof = value["transportProof"]
    if proof["transportKind"] != "gateway-v1" or proof["negotiatedHTTPVersion"] != "HTTP/2":
        raise ValueError("transport proof does not match Gateway V1 HTTP/2")
    if min(value["timing"]["wallSeconds"], value["timing"]["processCPUSeconds"]) <= 0:
        raise ValueError("nonpositive timing/CPU")
    if value["timing"]["processCPUSeconds"] / value["timing"]["wallSeconds"] > 4.02:
        raise ValueError("process CPU exceeds four-vCPU capacity")
    latency = value["latencyMilliseconds"]
    values = [latency[name] for name in ("p50", "p90", "p95", "p99", "p999")]
    if any(not isinstance(value, (int, float)) or not 0 < value < float("inf") for value in values) or values != sorted(values):
        raise ValueError("latency percentiles are missing, nonpositive, or inconsistent")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=pathlib.Path, required=True)
    parser.add_argument("--leg", choices=("go-v2-native", "rust-sdk", "rust-driver"), required=True)
    parser.add_argument("--workers", type=int, choices=(32, 64), required=True)
    parser.add_argument("--repeat", type=int, required=True)
    parser.add_argument("--class", dest="run_class", choices=("clean-score", "perf", "perf-calibration", "ffi-flow", "rust-flow", "go-diagnostic", "scheduler"), required=True)
    parser.add_argument("--calibration-mode", choices=("fp", "dwarf"))
    parser.add_argument("--key-fd", type=int, default=3)
    parser.add_argument("command", nargs=argparse.REMAINDER)
    args = parser.parse_args()
    if args.command and args.command[0] == "--":
        args.command.pop(0)
    if not args.command or args.key_fd < 3 or args.repeat < 1:
        raise ValueError("missing command, invalid key descriptor, or invalid repeat")
    os.fstat(args.key_fd)
    if args.run_class in ("ffi-flow", "go-diagnostic", "scheduler") and args.leg != "go-v2-native":
        raise ValueError("Go diagnostic classes require the Go leg")
    if args.run_class == "rust-flow" and args.leg == "go-v2-native":
        raise ValueError("Rust operation-rate diagnostics require a Rust leg")
    if (args.run_class == "perf-calibration") != (args.calibration_mode is not None):
        raise ValueError("calibration mode is required only for perf-calibration legs")
    is_perf = args.run_class in ("perf", "perf-calibration")
    is_scheduler = args.run_class == "scheduler"
    if len(os.sched_getaffinity(0)) != 4:
        raise ValueError("expected exactly four available CPUs")
    wire = json.loads((args.root / "provenance/wire-proof.json").read_text())
    proof = wire["legs"][args.leg]
    if not (proof["binaryDisabledExplicitly"] and proof["eligibleResponses"] > 0
            and proof["binaryResponsesBeforeConversion"] == 0
            and proof["textResponsesBeforeConversion"] == proof["eligibleResponses"]
            and wire["rustSourceCommit"] == "77ebe786ad1012bc9ac56105e08c05fb4c65b321"):
        raise ValueError("missing or invalid source-pinned pre-conversion wire proof")
    key = f"w{args.workers}/{args.run_class}/repeat-{args.repeat}/{args.leg}"
    existing = accepted_directory(args.root, args.workers, args.run_class, args.repeat, args.leg)
    if (existing / "VALID.json").exists() and not (existing / "REJECTED.json").exists():
        valid = json.loads((existing / "VALID.json").read_text())
        if valid["binarySHA256"] == digest(pathlib.Path(args.command[0])):
            record_accepted(args.root, key, existing)
            print("REUSING_VERIFIED_ACCEPTED_LEG " + key)
            return
        print("PRESERVING_PRIOR_EXECUTABLE_LEG " + key)
    base = args.root / "run" / key
    leg = base
    attempt = 1
    while leg.exists():
        attempt += 1
        leg = base.with_name(base.name + f"-attempt-{attempt}")
    leg.mkdir(parents=True, exist_ok=False)
    raw = leg / "raw"
    raw.mkdir()
    (leg / "derived").mkdir()
    barrier = raw / "barriers"
    barrier.mkdir()
    binary = pathlib.Path(args.command[0]).resolve(strict=True)
    expected = json.loads(pathlib.Path(str(binary) + ".provenance.json").read_text())
    binary_sha = digest(binary)
    if expected["binary"]["sha256"] != binary_sha:
        raise ValueError("binary hash differs from sealed build provenance")
    callgraph = "fp"
    if args.run_class == "perf-calibration":
        callgraph = "fp" if args.calibration_mode == "fp" else "dwarf,16384"
    if args.run_class == "perf":
        calibration = json.loads((args.root / "provenance/callgraph-calibration.json").read_text())
        if not calibration["valid"] or binary_sha not in calibration["binarySHA256"]:
            raise ValueError("passing call-chain calibration is required for this exact binary")
        callgraph = calibration["callgraph"]
        if callgraph not in ("fp", "dwarf,16384"):
            raise ValueError("unsupported calibrated callgraph")
    if proof["buildClassSourceHash"] != expected["wireSourceHash"]:
        raise ValueError("wire proof and measured binary are not from the same wire-option/source configuration")
    write_json(leg / "provenance.json", expected)
    # Hard links are safe only to sealed read-only staging files on the same filesystem.
    if binary.stat().st_mode & 0o222:
        raise ValueError("staged binary must be read-only")
    os.link(binary, raw / "binary")
    write_json(raw / "binary.json", {"sha256": binary_sha, "preservedBy": "read-only hard link"})
    cgroup = memory_root() / f"cosmos-flow-{os.getpid()}-{time.time_ns()}"
    cgroup.mkdir()
    try:
        memory_value(cgroup)
    except RuntimeError:
        cgroup.rmdir()
        raise RuntimeError("memory controller unavailable for dedicated leg cgroup")
    perf_environment = os.environ.copy()
    perf_environment["HOME"] = str(args.root / "tools/perf-home")
    environment = os.environ.copy()
    for name in ("COSMOS_PERF_CPU_PROFILE", "COSMOS_PERF_HEAP_PROFILE", "COSMOS_PERF_RUNTIME_TRACE",
                 "COSMOS_PERF_BLOCK_PROFILE", "COSMOS_PERF_MUTEX_PROFILE", "COSMOS_PERF_FLOW_DIR"):
        environment.pop(name, None)
    environment.update(
        COSMOS_PERF_ADAPTER_RESULT=str(raw / "adapter-result.json"),
        COSMOS_PERF_BARRIER_DIR=str(barrier), COSMOS_PERF_CLK_TCK=str(os.sysconf("SC_CLK_TCK")),
        GOMAXPROCS="4",
    )
    if args.run_class == "ffi-flow":
        environment["COSMOS_PERF_FLOW_DIR"] = str(raw)
    if args.run_class == "go-diagnostic":
        for variable, name in (
            ("COSMOS_PERF_CPU_PROFILE", "cpu.pprof"), ("COSMOS_PERF_HEAP_PROFILE", "heap.pprof"),
            ("COSMOS_PERF_RUNTIME_TRACE", "runtime.trace"), ("COSMOS_PERF_BLOCK_PROFILE", "block.pprof"),
            ("COSMOS_PERF_MUTEX_PROFILE", "mutex.pprof"),
        ):
            environment[variable] = str(raw / name)
        environment["COSMOS_PERF_TRACE_SECONDS"] = "5"

    def isolate():
        (cgroup / "cgroup.procs").write_text(str(os.getpid()))

    process = profiler = None
    failure = None
    try:
        with (raw / "stdout.txt").open("xb") as stdout, (raw / "stderr.txt").open("xb") as stderr, \
             (raw / "samples.jsonl").open("x") as samples, \
             (raw / "perf-record.stdout").open("xb") as perf_stdout, \
             (raw / "perf-record.stderr").open("xb") as perf_stderr:
            process = subprocess.Popen(args.command, stdout=stdout, stderr=stderr, env=environment,
                                       pass_fds=(args.key_fd,), preexec_fn=isolate)
            write_json(raw / "pid.json", {"pid": process.pid, "affinity": sorted(os.sched_getaffinity(process.pid))})
            started = time.monotonic()
            next_sample = started
            next_full = started + 1
            sample_index = 0
            start_released = end_released = False
            arm_started = None
            while process.poll() is None:
                now = time.monotonic()
                if now - started > 240:
                    raise TimeoutError("leg exceeded bounded 240-second process lifetime")
                if now >= next_sample:
                    proc = pathlib.Path("/proc") / str(process.pid)
                    try:
                        status = (proc / "status").read_text()
                        stat = (proc / "stat").read_text()
                        row = {
                            "monotonicNs": time.monotonic_ns(), "unixNs": time.time_ns(),
                            "status": status, "stat": stat,
                            "cgroupMemoryCurrent": memory_value(cgroup),
                            "cgroupMemoryPeak": memory_value(cgroup, peak=True),
                        }
                        samples.write(json.dumps(row) + "\n")
                        samples.flush()
                    except (FileNotFoundError, ProcessLookupError):
                        if process.poll() is None:
                            raise
                    next_sample = now + 0.25
                if now >= next_full and not end_released:
                    try:
                        snapshot(process.pid, raw, f"periodic-{sample_index:03}", cgroup)
                    except (FileNotFoundError, ProcessLookupError):
                        if process.poll() is None:
                            raise
                    sample_index += 1
                    next_full = now + 15
                if (barrier / "request-start").exists() and not start_released:
                    if is_perf or is_scheduler:
                        if profiler is None:
                            command = (
                                ["perf", "sched", "record", "-p", str(process.pid),
                                 "-o", str(raw / "perf-sched.data")]
                                if is_scheduler else
                                ["perf", "record", "-e", "cpu-clock", "-F", "49", "-g",
                                 "--call-graph", callgraph, "--clockid", "mono", "-p",
                                 str(process.pid), "-o", str(raw / "perf.data")]
                            )
                            profiler = subprocess.Popen(
                                command,
                                stdout=perf_stdout, stderr=perf_stderr, env=perf_environment)
                            arm_started = now
                        if not perf_armed(profiler):
                            if now - arm_started > 10:
                                raise TimeoutError("perf did not arm in 10 seconds")
                            time.sleep(0.01)
                            continue
                    snapshot(process.pid, raw, "measurement-start", cgroup)
                    (barrier / "ready-start").touch(exist_ok=False)
                    start_released = True
                if (barrier / "request-end").exists() and not end_released:
                    if not start_released:
                        raise RuntimeError("measurement ended before start barrier")
                    if profiler is not None:
                        profiler.send_signal(signal.SIGINT)
                        status = profiler.wait(timeout=20)
                        write_json(raw / "profiler-stop.json", {
                            "exitCode": status, "requestedSignal": "SIGINT",
                            "controlledStopAtMeasurementEnd": True,
                        })
                        if status not in (0, -signal.SIGINT, 128 + signal.SIGINT):
                            raise RuntimeError(f"perf record failed: exit {status}")
                    snapshot(process.pid, raw, "measurement-end", cgroup)
                    (barrier / "ready-end").touch(exist_ok=False)
                    end_released = True
                if profiler is not None and start_released and not end_released and profiler.poll() is not None:
                    raise RuntimeError("profiler exited before the measurement ended")
                time.sleep(0.01)
            if process.returncode:
                raise RuntimeError(f"adapter failed: exit {process.returncode}; raw error evidence preserved")
            if not (start_released and end_released):
                raise RuntimeError("adapter did not implement both measurement barriers")
        value = json.loads((raw / "adapter-result.json").read_text())
        validate_metrics(value, args.workers)
        bounds = json.loads((barrier / "measurement-bounds.json").read_text())
        begin, end = bounds["monotonicStartNs"], bounds["monotonicEndNs"]
        if end <= begin:
            raise ValueError("invalid monotonic measurement bounds")
        write_json(leg / "metrics.json", value)
        if is_perf:
            interval = f"{perf_time(begin)},{perf_time(end)}"
            for name, command in (
                ("perf-script.txt", ["perf", "script", "--no-inline", "--ns", "--show-lost-events", "--time", interval, "-i", str(raw / "perf.data")]),
                ("perf-self.txt", ["perf", "report", "--stdio", "--no-inline", "--no-children", "--percent-limit", "0",
                                   "--show-nr-samples", "--sort", "dso,symbol", "--time", interval,
                                   "-i", str(raw / "perf.data")]),
            ):
                with (raw / name).open("xb") as output, (raw / (name + ".stderr")).open("xb") as errors:
                    subprocess.run(command, stdout=output, stderr=errors, check=True, env=perf_environment)
            if (raw / "perf-script.txt").stat().st_size == 0:
                raise RuntimeError("measurement-filtered perf has no samples")
            from calibrate import calculate
            quality = calculate((raw / "perf-script.txt").read_text(),
                                (raw / "perf-record.stderr").read_text())
            write_json(leg / "derived/quality.json", quality)
            if args.run_class == "perf" and (not quality["valid"] or quality["lostSamples"] != 0):
                raise RuntimeError("profile failed call-chain quality gates")
        if args.run_class == "go-diagnostic":
            for name in ("cpu.pprof", "heap.pprof", "runtime.trace", "block.pprof", "mutex.pprof"):
                if not (raw / name).is_file() or (raw / name).stat().st_size == 0:
                    raise RuntimeError(f"missing Go diagnostic evidence: {name}")
        if args.run_class == "ffi-flow":
            from validate_flow import validate
            result = validate(json.loads((raw / "go-flow.json").read_text()),
                              json.loads((raw / "native-flow.json").read_text()), value, args.workers)
            write_json(leg / "derived/flow.json", result)
        if args.run_class == "rust-flow":
            from validate_rust_flow import validate
            write_json(leg / "derived/operation-rates.json", validate(value))
        if is_scheduler:
            for name, command in (
                ("sched-timehist.txt", ["perf", "sched", "timehist", "-i", str(raw / "perf-sched.data")]),
                ("sched-latency.txt", ["perf", "sched", "latency", "-i", str(raw / "perf-sched.data")]),
            ):
                with (raw / name).open("xb") as output, (raw / (name + ".stderr")).open("xb") as errors:
                    subprocess.run(command, stdout=output, stderr=errors, check=True, env=perf_environment)
            write_json(leg / "derived/scheduler-scope.json", {
                "scope": "PID-attached perf sched recording of the benchmark",
                "authoritativeScore": False,
                "limitation": "Not a system-wide scheduler trace; do not infer exhaustive host wakeup/runqueue delays.",
                "threadCounters": "per-thread status/schedstat snapshots; disabled schedstats are marked explicitly",
            })
        write_json(leg / "VALID.json", {"valid": True, "runClass": args.run_class, "binarySHA256": binary_sha})
        record_accepted(args.root, key, leg)
    except BaseException as error:
        failure = error
        write_json(leg / "REJECTED.json", {"valid": False, "reasonType": type(error).__name__,
                                         "reason": str(error)})
        raise
    finally:
        if profiler is not None and profiler.poll() is None:
            profiler.send_signal(signal.SIGINT)
            profiler.wait(timeout=20)
        if process is not None and process.poll() is None:
            process.terminate()
            try:
                process.wait(timeout=15)
            except subprocess.TimeoutExpired:
                process.kill()
                process.wait(timeout=10)
        if cgroup.exists():
            if (cgroup / "cgroup.procs").read_text().strip():
                if failure is None:
                    raise RuntimeError("dedicated cgroup is not empty after leg exit")
            else:
                cgroup.rmdir()
        seal(raw)


if __name__ == "__main__":
    main()
