#!/usr/bin/env python3
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

"""Attribute captured perf samples without changing physical cost to caller cost."""

import argparse
import collections
import csv
import dataclasses
import gzip
import hashlib
import html
import io
import json
import math
import pathlib
import re
import subprocess


CATEGORIES = (
    "go-runtime-sdk",
    "cgo-completion-ffi",
    "tokio",
    "http2-hpack",
    "http-client",
    "tls",
    "serialization-body-allocation-memcopy",
    "shared-key-signing-hmac-headers",
    "cosmos-routing-session-transport",
    "kernel-scheduling-syscalls",
    "kernel-network-memory",
    "kernel-other",
    "native-synchronization-syscalls",
    "native-standard-library",
    "resolved-other",
    "residual-unknown",
)
LABELS = {
    "go-runtime-sdk": "Go runtime / SDK",
    "cgo-completion-ffi": "cgo / native completion wrapper",
    "tokio": "Tokio / async runtime",
    "http2-hpack": "HTTP/2 + HPACK",
    "http-client": "Generic HTTP client",
    "tls": "TLS",
    "serialization-body-allocation-memcopy": "Allocation / copy / formatting",
    "shared-key-signing-hmac-headers": "Hash / HMAC / header processing",
    "cosmos-routing-session-transport": "Cosmos driver",
    "kernel-scheduling-syscalls": "Kernel scheduling / locks / syscalls",
    "kernel-network-memory": "Kernel network / memory",
    "kernel-other": "Other kernel work",
    "native-synchronization-syscalls": "Native synchronization / syscalls",
    "native-standard-library": "Generic Rust / native standard library",
    "resolved-other": "Other resolved symbols",
    "residual-unknown": "Unresolved sampled instruction",
}
HEADER = re.compile(
    r"^(?P<prefix>.*?)\s+(?P<time>\d+\.\d+):\s+"
    r"(?P<period>\d+)\s+(?P<event>\S+):\s*$"
)
FRAME = re.compile(r"^\s+(?P<ip>[0-9a-fA-F]+)\s+(?P<symbol>.+)\s+\((?P<dso>.+)\)\s*$")
GO = re.compile(
    r"^(?:runtime[./]|internal/|sync[./]|syscall[./]|reflect[./]|"
    r"fmt\.|encoding/|net[./]|crypto/|strings\.|bytes\.|time\.|"
    r"context\.|math[./]|sort\.|io[./]|os[./]|main\.|"
    r"github\.com/Azure/azure-sdk-for-go/)"
)
FFI = re.compile(
    r"^(?:_cgo_|crosscall|cosmos_[A-Za-z0-9_]+(?:\.[A-Za-z0-9_.]+)?$|azurecosmosdriver::|"
    r"runtime\.(?:cgocall|asmcgocall)|.*\._C(?:func|type)_|.*\._cgo_)"
)
KERNEL_NETWORK = re.compile(
    r"(?:^|_)(?:tcp|udp|ip|ipv6|sock|socket|net|napi|skb|inet|"
    r"recv(?:from|msg|mmsg)?|send(?:to|msg|mmsg)?|page|slab|kmem|kmalloc|kfree|mmap|mprotect|"
    r"copy_to_user|copy_from_user|copy_page|clear_page|alloc_pages|"
    r"free_pages|mem_cgroup|handle_mm_fault|do_user_addr_fault)(?:_|$)"
)
KERNEL_SCHEDULING = re.compile(
    r"(?:sched|switch|futex|epoll|poll_|spin_|spinlock|wake_|wakeup|"
    r"try_to_wake|syscall|sys_|task_work|task_tick|hrtimer|"
    r"update_rq_clock|select_task_rq|enqueue_task|dequeue_task)"
)
NATIVE_SYNC = re.compile(
    r"^(?:__lll_|__pthread_|pthread_|__futex|futex|"
    r"__nptl_|syscall$|epoll_|poll$|ppoll$|read$|write$|"
    r"recv|send|nanosleep$)|std::(?:sys::sync|sync::.*(?:Mutex|Condvar))"
)
ALLOCATION = re.compile(
    r"^(?:_?int_(?:malloc|free|realloc)|(?:__libc_)?"
    r"(?:malloc|calloc|realloc|free|cfree|aligned_alloc|posix_memalign))"
    r"|(?:__rust_(?:alloc|dealloc|realloc)|"
    r"alloc::(?:alloc|raw_vec|vec|boxed|fmt|string)|"
    r"core::fmt|std::alloc|"
    r"bytes::|serde(?:_json)?::|"
    r"(?:__)?mem(?:cpy|move|set|cmp)|"
    r"hashbrown::raw::.*(?:reserve|rehash))"
)
OWNER_CATEGORIES = {
    "cgo-completion-ffi",
    "tokio",
    "http2-hpack",
    "http-client",
    "tls",
    "shared-key-signing-hmac-headers",
    "cosmos-routing-session-transport",
}


@dataclasses.dataclass(frozen=True)
class Frame:
    symbol: str
    dso: str
    ip: str = ""
    original_symbol: str = ""
    source: str = ""


@dataclasses.dataclass(frozen=True)
class Sample:
    index: int
    timestamp: float
    period: int
    event: str
    frames: tuple


def clean_symbol(symbol):
    value = re.sub(r"^\[[.kguH]\]\s*", "", symbol.strip())
    value = re.sub(r"\+0x[0-9a-fA-F]+$", "", value)
    value = re.sub(r"(?<=\w)\[[0-9a-f]+\](?=::)", "", value)
    if "$" in value and not value.startswith("_ZN"):
        escapes = {
            "$LT$": "<", "$GT$": ">", "$u20$": " ", "$C$": ",",
            "$RF$": "&", "$LP$": "(", "$RP$": ")", "$BP$": "*",
            "$u7b$": "{", "$u7d$": "}", "$u5b$": "[", "$u5d$": "]",
        }
        for encoded, decoded in escapes.items():
            value = value.replace(encoded, decoded)
        value = value.replace("..", "::")
        if value.startswith("_<"):
            value = value[1:]
    return value


def is_unknown(symbol):
    return (
        not symbol
        or symbol in {"[unknown]", "unknown", "??", "[.]"}
        or re.fullmatch(r"(?:0x)?[0-9a-fA-F]{6,}", symbol) is not None
    )


def classify(frame):
    """Only the sampled function name and DSO kind may affect physical cost."""
    symbol = clean_symbol(frame.symbol)
    if is_unknown(symbol):
        return "residual-unknown"
    dso_name = frame.dso.replace("\\", "/").rsplit("/", 1)[-1]
    kernel = dso_name in {"[kernel.kallsyms]", "[kernel.vmlinux]"} or dso_name.endswith(".ko")
    if kernel:
        if KERNEL_NETWORK.search(symbol):
            return "kernel-network-memory"
        if KERNEL_SCHEDULING.search(symbol):
            return "kernel-scheduling-syscalls"
        return "kernel-other"
    if FFI.search(symbol):
        return "cgo-completion-ffi"
    if GO.search(symbol):
        return "go-runtime-sdk"
    if re.search(r"^<?(?:tokio|tokio_util|mio|futures|futures_util|futures_core)::", symbol):
        return "tokio"
    if re.search(r"^<?(?:h2|hpack|hyper|hyper_util)::", symbol):
        return "http2-hpack"
    if re.search(r"^<?(?:rustls|rustls_webpki|webpki|ring|openssl|tokio_rustls|hyper_rustls)::", symbol):
        return "tls"
    if re.search(r"^<?(?:reqwest|tower|tower_http|http_body|http_body_util)::", symbol):
        return "http-client"
    if re.search(r"^<?(?:azure_data_cosmos_driver|azure_data_cosmos)::", symbol):
        return "cosmos-routing-session-transport"
    if ALLOCATION.search(symbol):
        return "serialization-body-allocation-memcopy"
    if re.search(r"^<?(?:hmac|sha2|digest|base64|http::header)::|authorization|shared_key", symbol):
        return "shared-key-signing-hmac-headers"
    if NATIVE_SYNC.search(symbol):
        return "native-synchronization-syscalls"
    if re.match(r"^<?(?:core|std|alloc|hashbrown|url|percent_encoding)::", symbol):
        return "native-standard-library"
    return "resolved-other"


def caller_owner(frames):
    """Keep caller ownership separate; never promote an ancestor into Self CPU."""
    for frame in frames[1:]:
        category = classify(frame)
        if category in OWNER_CATEGORIES:
            return category, frame.symbol
        symbol = clean_symbol(frame.symbol)
        if symbol.startswith("github.com/Azure/azure-sdk-for-go/"):
            return "go-runtime-sdk", frame.symbol
    return "unattributed-caller", ""


def parse_perf(path):
    samples = []
    current = None
    frames = []

    def finish():
        if current is None:
            return
        if not frames:
            raise ValueError(f"{path}: sample {len(samples) + 1} has no frames")
        if int(current["period"]) <= 0:
            raise ValueError(f"{path}: nonpositive period")
        samples.append(
            Sample(
                len(samples) + 1,
                float(current["time"]),
                int(current["period"]),
                current["event"],
                tuple(frames),
            )
        )

    with path.open(encoding="utf-8", errors="strict") as stream:
        for line_number, line in enumerate(stream, 1):
            header = HEADER.fullmatch(line.rstrip("\n"))
            if header:
                finish()
                current = header.groupdict()
                frames = []
                continue
            if not line.strip() or line.startswith("#"):
                continue
            frame = FRAME.fullmatch(line.rstrip("\n"))
            if frame and current is not None:
                frames.append(Frame(clean_symbol(frame["symbol"]), frame["dso"], frame["ip"]))
                continue
            raise ValueError(f"{path}:{line_number}: unparsed perf record")
    finish()
    if not samples:
        raise ValueError(f"{path}: no samples")
    events = {sample.event for sample in samples}
    if len(events) != 1 or not next(iter(events)).startswith(("task-clock", "cpu-clock")):
        raise ValueError(f"{path}: unexpected or mixed event types: {sorted(events)}")
    return samples


def demangle_map(samples_by_variant):
    symbols = sorted(
        {
            frame.symbol
            for samples in samples_by_variant.values()
            for sample in samples
            for frame in sample.frames
        }
    )
    inputs = [re.sub(r"\.llvm\.\d+$", "", symbol) for symbol in symbols]
    result = subprocess.run(
        ["c++filt", "-n", "-s", "rust"],
        input="\n".join(inputs) + "\n",
        capture_output=True,
        text=True,
        check=True,
    )
    outputs = result.stdout.splitlines()
    if len(outputs) != len(symbols):
        raise ValueError("Rust demangler changed symbol count")
    return dict(zip(symbols, outputs))


def parse_mappings(text):
    mappings = []
    for line in text.splitlines():
        match = re.match(
            r"^([0-9a-f]+)-([0-9a-f]+)\s+[rwxps-]{4}\s+"
            r"([0-9a-f]+)\s+\S+\s+\d+\s*(.*)$",
            line,
        )
        if match:
            mappings.append((int(match[1], 16), int(match[2], 16), int(match[3], 16), match[4]))
    if not mappings:
        raise ValueError("No process address mappings parsed")
    return mappings


def file_offset(ip, dso, mappings):
    address = int(ip, 16)
    matches = [
        address - start + offset
        for start, end, offset, path in mappings
        if start <= address < end and path == dso
    ]
    if len(matches) > 1:
        raise ValueError("Ambiguous process mapping for sampled instruction")
    return matches[0] if matches else None


def resolve_libc(samples_by_variant, root):
    build_id = "6a6cb5a80fd5e5bb5bf3aeaa111add98e88843a5"
    debug = (
        root / "attribution-symbols/libc6-dbg/extracted/usr/lib/debug/.build-id"
        / build_id[:2] / (build_id[2:] + ".debug")
    )
    if not debug.is_file():
        raise ValueError("Exact-build libc debug symbols are missing; refusing inferred symbol ownership")
    notes = subprocess.check_output(["readelf", "-n", str(debug)], text=True)
    if f"Build ID: {build_id}" not in notes:
        raise ValueError("libc debug-symbol build ID mismatch")
    offsets = {}
    for variant, samples in samples_by_variant.items():
        leg = root / "perf-v4" / variant / "w64/perf/repeat-1/go-v2-native"
        captured = subprocess.check_output(
            ["perf", "buildid-list", "-i", str(leg / "raw/perf.data")], text=True
        )
        if not re.search(rf"^{build_id}\s+.*libc\.so\.6$", captured, re.M):
            raise ValueError(f"{variant}: recorded libc build ID does not match debug symbols")
        mappings = parse_mappings((leg / "raw/smaps.txt").read_text())
        for sample in samples:
            for frame in sample.frames:
                if frame.dso.endswith("/libc.so.6") and is_unknown(frame.symbol):
                    offset = file_offset(frame.ip, frame.dso, mappings)
                    if offset is not None:
                        offsets[(variant, frame.ip)] = offset
    unique = sorted(set(offsets.values()))
    result = subprocess.run(
        ["addr2line", "-f", "-C", "-e", str(debug)],
        input="".join(f"0x{offset:x}\n" for offset in unique),
        capture_output=True,
        text=True,
        check=True,
    )
    lines = result.stdout.splitlines()
    if len(lines) != 2 * len(unique):
        raise ValueError("addr2line did not return exactly one function/source pair per address")
    resolved = {
        offset: (lines[index * 2], lines[index * 2 + 1])
        for index, offset in enumerate(unique)
    }
    remapped = {}
    resolved_count = 0
    for variant, samples in samples_by_variant.items():
        output = []
        for sample in samples:
            frames = []
            for frame in sample.frames:
                offset = offsets.get((variant, frame.ip))
                if offset is not None:
                    symbol, source = resolved[offset]
                    if not is_unknown(symbol):
                        frame = dataclasses.replace(
                            frame, symbol=symbol, original_symbol=frame.symbol, source=source
                        )
                        resolved_count += 1
                frames.append(frame)
            output.append(dataclasses.replace(sample, frames=tuple(frames)))
        remapped[variant] = output
    return remapped, {
        "library": "libc.so.6",
        "capturedAndDebugBuildId": build_id,
        "debugFileSHA256": hashlib.sha256(debug.read_bytes()).hexdigest(),
        "uniqueInstructionOffsets": len(unique),
        "resolvedFrameOccurrences": resolved_count,
        "method": "captured smaps address/file offsets plus exact-build addr2line",
    }


def remap_samples(samples, symbols):
    return [
        dataclasses.replace(
            sample,
            frames=tuple(
                dataclasses.replace(frame, symbol=symbols[frame.symbol])
                for frame in sample.frames
            ),
        )
        for sample in samples
    ]


def summarize(samples, metrics, quality):
    if len(samples) != quality["samples"]:
        raise ValueError("Parsed samples do not match perf capture sample count")
    if quality.get("lostSamples", 0) != 0 or quality["lostEventPercent"] != 0:
        raise ValueError("Capture reports lost samples")
    counts = metrics["counts"]
    if counts["success"] <= 0 or any(counts[key] for key in ("error", "retry", "throttle")):
        raise ValueError("Profiled operations are not clean")
    periods = collections.Counter()
    callers = collections.Counter()
    leaves = collections.Counter()
    caller_leaves = collections.Counter()
    mapping = []
    for sample in samples:
        leaf = sample.frames[0]
        category = classify(leaf)
        owner, caller = caller_owner(sample.frames)
        periods[category] += sample.period
        callers[owner] += sample.period
        leaves[(category, leaf.symbol, leaf.dso)] += sample.period
        caller_leaves[(category, leaf.symbol, owner, caller)] += sample.period
        mapping.append({
            "sample": sample.index,
            "timestamp": sample.timestamp,
            "period": sample.period,
            "selfCategory": category,
            "leaf": leaf.symbol,
            "originalLeaf": leaf.original_symbol,
            "instructionAddress": leaf.ip,
            "source": leaf.source,
            "dso": leaf.dso,
            "callerCategory": owner,
            "caller": caller,
        })
    total = sum(sample.period for sample in samples)
    if sum(periods.values()) != total or sum(callers.values()) != total:
        raise ValueError("Sample accounting drift")
    if len(mapping) != len(samples):
        raise ValueError("Not every sample received exactly one mapping")
    percentages = {category: 100.0 * periods[category] / total for category in CATEGORIES}
    if not math.isclose(sum(percentages.values()), 100.0, abs_tol=1e-10):
        raise ValueError("Self percentages do not normalize to 100%")
    cpu = metrics["timing"]["processCPUSeconds"]
    wall = metrics["timing"]["wallSeconds"]
    result = {
        "samples": len(samples),
        "event": samples[0].event,
        "totalPeriodNs": total,
        "firstSampleTimestamp": min(sample.timestamp for sample in samples),
        "lastSampleTimestamp": max(sample.timestamp for sample in samples),
        "sampleWindowSeconds": max(sample.timestamp for sample in samples) - min(sample.timestamp for sample in samples),
        "profileCPUUsPerRead": cpu * 1_000_000.0 / counts["success"],
        "profileThroughput": counts["success"] / wall,
        "profileBusyCores": cpu / wall,
        "sampledCPUSeconds": total / 1_000_000_000.0,
        "processCPUSeconds": cpu,
        "sampledToProcessCPUPercent": total / 1_000_000_000.0 / cpu * 100.0,
        "selfPercentTotal": sum(percentages.values()),
        "unresolvedLeafPercent": percentages["residual-unknown"],
        "captureQuality": quality,
        "selfPercent": percentages,
        "callerPercent": {key: period / total * 100.0 for key, period in sorted(callers.items())},
        "topLeaves": [
            {"category": category, "symbol": symbol, "dso": dso, "selfPercent": period / total * 100.0}
            for (category, symbol, dso), period in leaves.most_common(50)
        ],
        "topCallerSites": [
            {
                "selfCategory": category,
                "leaf": symbol,
                "callerCategory": owner,
                "caller": caller,
                "selfPercent": period / total * 100.0,
            }
            for (category, symbol, owner, caller), period in caller_leaves.most_common(50)
        ],
    }
    return result, mapping


def category_comparison(variants, normalizations):
    rows = []
    for category in CATEGORIES:
        control_percent = variants["control"]["selfPercent"][category]
        synthetic_percent = variants["s1-3ms"]["selfPercent"][category]
        row = {
            "category": category,
            "label": LABELS[category],
            "controlSelfPercent": control_percent,
            "syntheticSelfPercent": synthetic_percent,
        }
        for kind in ("profile", "clean"):
            control_us = control_percent * normalizations["control"][kind] / 100.0
            synthetic_us = synthetic_percent * normalizations["s1-3ms"][kind] / 100.0
            row[kind + "ControlCPUUsPerRead"] = control_us
            row[kind + "SyntheticCPUUsPerRead"] = synthetic_us
            row[kind + "DifferenceCPUUsPerRead"] = control_us - synthetic_us
        rows.append(row)
    for variant, prefix in (("control", "Control"), ("s1-3ms", "Synthetic")):
        for kind in ("profile", "clean"):
            total = sum(row[kind + prefix + "CPUUsPerRead"] for row in rows)
            if not math.isclose(total, normalizations[variant][kind], abs_tol=1e-8):
                raise ValueError("CPU-per-read bridge does not reconcile")
    return sorted(rows, key=lambda row: -row["cleanDifferenceCPUUsPerRead"])


def write_csv(path, rows):
    if not rows:
        raise ValueError("Cannot write an empty accounting table")
    with path.open("w", newline="", encoding="utf-8") as stream:
        writer = csv.DictWriter(stream, fieldnames=list(rows[0]))
        writer.writeheader()
        writer.writerows(rows)


def svg_chart(rows, kind):
    left = 370
    width = 1260
    height = 145 + 46 * len(rows)
    maximum = max(row[kind + "ControlCPUUsPerRead"] for row in rows) * 1.15
    scale = (width - left - 140) / maximum
    parts = [
        f'<svg xmlns="http://www.w3.org/2000/svg" width="{width}" height="{height}" viewBox="0 0 {width} {height}">',
        '<rect width="100%" height="100%" fill="white"/>',
        '<style>text{font-family:Segoe UI,Arial,sans-serif;fill:#0f172a}.label{font-size:13px}.value{font-size:12px}</style>',
        f'<text x="24" y="32" font-size="22" font-weight="700">Corrected Self CPU/read - {kind}-run normalization</text>',
        '<text x="24" y="56" font-size="13">Estimated from sampled instructions; callers do not reassign physical Self cost.</text>',
        '<rect x="24" y="72" width="16" height="12" fill="#2563eb"/><text x="48" y="83" class="label">Real transport</text>',
        '<rect x="190" y="72" width="16" height="12" fill="#f59e0b"/><text x="214" y="83" class="label">Synthetic 3 ms</text>',
    ]
    for index, row in enumerate(rows):
        y = 105 + index * 46
        control = row[kind + "ControlCPUUsPerRead"]
        synthetic = row[kind + "SyntheticCPUUsPerRead"]
        parts.append(f'<text x="24" y="{y + 17}" class="label">{html.escape(row["label"])}</text>')
        for offset, value, color in ((0, control, "#2563eb"), (18, synthetic, "#f59e0b")):
            bar = value * scale
            parts.append(f'<rect x="{left}" y="{y + offset}" width="{bar:.3f}" height="15" fill="{color}"/>')
            parts.append(f'<text x="{left + bar + 6:.3f}" y="{y + offset + 12}" class="value">{value:.2f} us/read</text>')
    parts.append(
        f'<text x="24" y="{height - 15}" class="label">All samples counted once, including unresolved instructions. W64, Tokio=2, one Go reactor. No new benchmark.</text>'
    )
    parts.append("</svg>")
    return "\n".join(parts) + "\n"


def verify_raw(leg):
    entries = {}
    for line in (leg / "raw/SHA256SUMS").read_text().splitlines():
        match = re.fullmatch(r"([0-9a-f]{64})\s+[*]?(.+)", line)
        if match:
            entries[pathlib.Path(match[2]).name] = match[1]
    checked = {}
    for name in ("perf-script.txt", "perf-quality.json", "perf-record.stderr", "perf-stat.txt", "adapter-result.json", "smaps.txt"):
        path = leg / "raw" / name
        digest = hashlib.sha256(path.read_bytes()).hexdigest()
        if entries.get(name) != digest:
            raise ValueError(f"Sealed raw checksum mismatch: {name}")
        checked[name] = digest
    return checked


def build_report(root, output):
    if output.exists():
        raise ValueError(f"Refusing to overwrite derived report: {output}")
    samples = {}
    evidence = {}
    for variant in ("control", "s1-3ms"):
        leg = root / "perf-v4" / variant / "w64/perf/repeat-1/go-v2-native"
        evidence[variant] = verify_raw(leg)
        samples[variant] = parse_perf(leg / "raw/perf-script.txt")
    samples, symbol_resolution = resolve_libc(samples, root)
    symbol_map = demangle_map(samples)
    clean = json.loads((root / "comparison/summary.json").read_text())
    baselines = {
        "control": next(row for row in clean["aggregates"] if row["workers"] == 64 and row["variant"] == "c0-control"),
        "s1-3ms": next(row for row in clean["aggregates"] if row["workers"] == 64 and row["variant"] == "s1-3ms"),
    }
    variants = {}
    mappings = {}
    normalizations = {}
    for variant in samples:
        leg = root / "perf-v4" / variant / "w64/perf/repeat-1/go-v2-native"
        metrics = json.loads((leg / "metrics.json").read_text())
        quality = json.loads((leg / "raw/perf-quality.json").read_text())
        variants[variant], mappings[variant] = summarize(
            remap_samples(samples[variant], symbol_map), metrics, quality
        )
        normalizations[variant] = {
            "profile": variants[variant]["profileCPUUsPerRead"],
            "clean": baselines[variant]["cpu_us_per_read"],
        }
    rows = category_comparison(variants, normalizations)
    output.mkdir(parents=True)
    (output / "tools").mkdir()
    (output / "tools/attribute.py").write_bytes(pathlib.Path(__file__).read_bytes())
    write_csv(output / "self-attribution.csv", rows)
    leaf_rows = []
    caller_rows = []
    for variant, result in variants.items():
        for row in result["topLeaves"]:
            leaf_rows.append({"variant": variant, **row})
        for row in result["topCallerSites"]:
            caller_rows.append({"variant": variant, **row})
        with (output / f"{variant}-sample-mapping.csv.gz").open("wb") as stream:
            with gzip.GzipFile(fileobj=stream, mode="wb", mtime=0, filename="") as zipped:
                with io.TextIOWrapper(zipped, encoding="utf-8", newline="") as text:
                    writer = csv.DictWriter(text, fieldnames=list(mappings[variant][0]))
                    writer.writeheader()
                    writer.writerows(mappings[variant])
    write_csv(output / "top-leaves.csv", leaf_rows)
    write_csv(output / "caller-sites.csv", caller_rows)
    result = {
        "version": 2,
        "status": "corrected-self-accounting",
        "sourceCampaign": root.name,
        "sourceRawSHA256": evidence,
        "nativeSymbolResolution": symbol_resolution,
        "normalizationsCPUUsPerRead": normalizations,
        "cleanScoreBaselines": baselines,
        "variants": variants,
        "selfAttribution": rows,
        "limitations": [
            "Only one selected perf capture per variant, separate from the two-repeat clean scores.",
            "CPU/read category values estimate shares from finite on-CPU sampling; they are not direct category timers.",
            "Clean-normalized values assume the profiled composition represents the unprofiled clean run.",
            "Caller ownership is a separate table, never added to or substituted for physical Self CPU.",
            "An unresolved instruction remains unresolved even when a caller is known.",
            "Captured task-clock period and measured process CPU differ; coverage is reported rather than silently equated.",
            "The synthetic path changes network, response, and timing behavior; subtraction is not a precise additive causal decomposition.",
            "Attribution-v1 is withdrawn because binary-path matching and caller promotion corrupted its categories.",
        ],
        "benchmarkStatus": "No new benchmark was run; dispatcher A/B remains on hold.",
    }
    for row in rows:
        for variant, prefix in (("control", "Control"), ("s1-3ms", "Synthetic")):
            share = variants[variant]["selfPercent"][row["category"]] / 100.0
            row["profile" + prefix + "FourVCPUPercent"] = share * variants[variant]["profileBusyCores"] / 4 * 100.0
            row["clean" + prefix + "FourVCPUPercent"] = share * baselines[variant]["busy_cores"] / 4 * 100.0
    write_csv(output / "self-attribution.csv", rows)
    (output / "attribution.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    for kind in ("profile", "clean"):
        (output / f"cpu-per-read-{kind}.svg").write_text(svg_chart(rows, kind), encoding="utf-8")
    lines = [
        "# Corrected control-versus-synthetic CPU attribution",
        "",
        "This is derived analysis of the sealed CPU-isolation campaign, not a new benchmark.",
        "It is separate from the authoritative first-pass and earlier batch-size campaigns.",
        "**Attribution-v1 is withdrawn.** Category matching now excludes DSO paths and preserves",
        "Go, native, kernel, and unresolved physical cost independently of callers.",
        "",
        "| Variant | Samples | Profile CPU/read | Clean CPU/read | Unresolved leaf | Sampled/process CPU |",
        "|---|---:|---:|---:|---:|---:|",
    ]
    for variant, value in variants.items():
        lines.append(
            f"| {variant} | {value['samples']} | {value['profileCPUUsPerRead']:.2f} us | "
            f"{normalizations[variant]['clean']:.2f} us | {value['unresolvedLeafPercent']:.2f}% | "
            f"{value['sampledToProcessCPUPercent']:.2f}% |"
        )
    lines += [
        "",
        "## Physical Self cost",
        "",
        "Values below are **estimates normalized to the clean-score CPU/read**, not additional",
        "measurements. The profile-normalized view is available separately. No denominators are added.",
        "",
        "| Category | Real us/read | Synthetic us/read | Difference us/read |",
        "|---|---:|---:|---:|",
    ]
    for row in rows:
        lines.append(
            f"| {row['label']} | {row['cleanControlCPUUsPerRead']:.2f} | "
            f"{row['cleanSyntheticCPUUsPerRead']:.2f} | {row['cleanDifferenceCPUUsPerRead']:.2f} |"
        )
    lines += ["", "## Most expensive resolved control leaves", ""]
    for row in [row for row in variants["control"]["topLeaves"] if row["category"] != "residual-unknown"][:15]:
        lines.append(f"- {row['selfPercent']:.2f}%: `{row['symbol']}` ({row['category']})")
    lines += ["", "## Interpretation limits", ""]
    lines += ["- " + limitation for limitation in result["limitations"]]
    lines += [
        "",
        "A known caller can explain an allocation or lock, but does not prove that the caller",
        "is removable. Select a targeted A/B only after checking concrete call sites and any",
        "material unresolved or truncated stack coverage.",
        "",
        "Reference: https://man7.org/linux/man-pages/man1/perf-report.1.html",
    ]
    (output / "report.md").write_text("\n".join(lines) + "\n", encoding="utf-8")
    table = "".join(
        f"<tr><td>{html.escape(row['label'])}</td><td>{row['cleanControlCPUUsPerRead']:.2f}</td>"
        f"<td>{row['cleanSyntheticCPUUsPerRead']:.2f}</td><td>{row['cleanDifferenceCPUUsPerRead']:.2f}</td></tr>"
        for row in rows
    )
    limitations = "".join(f"<li>{html.escape(text)}</li>" for text in result["limitations"])
    page = (
        '<!doctype html><html lang="en"><meta charset="utf-8"><title>Corrected Cosmos CPU attribution</title>'
        '<style>body{font:15px Segoe UI,Arial,sans-serif;margin:32px;color:#0f172a;max-width:1400px}'
        '.warning{background:#fff7ed;border:1px solid #fb923c;padding:16px}'
        'table{border-collapse:collapse}td,th{padding:9px 16px;border:1px solid #cbd5e1;text-align:right}'
        'td:first-child,th:first-child{text-align:left}svg{max-width:100%;height:auto}'
        'p,li{line-height:1.55}</style>'
        '<h1>Corrected control-versus-synthetic CPU attribution</h1>'
        '<p class="warning"><b>Attribution-v1 is withdrawn.</b> This report separates sampled Self cost from caller ownership. '
        'No new benchmark was run. The dispatcher A/B is on hold.</p>'
        '<p>W64, Tokio=2, one Go reactor. Existing real-transport and synthetic-3-ms captures. '
        'Separate from the authoritative first-pass and batch-size A/B.</p>'
        '<h2>Clean-normalized estimates</h2>'
        + svg_chart(rows, "clean")
        + '<table><tr><th>Physical Self category</th><th>Real us/read</th><th>Synthetic us/read</th><th>Difference us/read</th></tr>'
        + table + '</table><h2>Profile-run-normalized estimates</h2>'
        + svg_chart(rows, "profile")
        + '<h2>Evidence and limitations</h2><ul>' + limitations + '</ul>'
        '<p><a href="report.md">Detailed report</a> | <a href="attribution.json">JSON</a> | '
        '<a href="self-attribution.csv">Self CSV</a> | <a href="top-leaves.csv">Leaf symbols</a> | '
        '<a href="caller-sites.csv">Separate caller ownership</a></p></html>'
    )
    (output / "dashboard.html").write_text(page, encoding="utf-8")
    manifest = []
    for path in sorted(output.rglob("*")):
        if path.is_file() and path.name != "SHA256SUMS":
            manifest.append(f"{hashlib.sha256(path.read_bytes()).hexdigest()}  {path.relative_to(output).as_posix()}")
    (output / "SHA256SUMS").write_text("\n".join(manifest) + "\n", encoding="ascii")
    print(json.dumps({
        "status": result["status"],
        "samples": {key: value["samples"] for key, value in variants.items()},
        "normalizations": normalizations,
        "unresolvedSelfPercent": {key: value["unresolvedLeafPercent"] for key, value in variants.items()},
        "nativeSymbolResolution": symbol_resolution,
        "largestSelfDifferences": [
            {"category": row["category"], "cleanDifferenceCPUUsPerRead": row["cleanDifferenceCPUUsPerRead"]}
            for row in rows[:8]
        ],
    }, sort_keys=True))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--campaign-root", type=pathlib.Path)
    parser.add_argument("--output", type=pathlib.Path)
    args = parser.parse_args()
    root = args.campaign_root or pathlib.Path("/mnt/cosmos-perf/cpu-isolation-synthetic-20261105-d")
    output = args.output or root / "attribution-v2-final"
    build_report(root, output)


if __name__ == "__main__":
    main()
