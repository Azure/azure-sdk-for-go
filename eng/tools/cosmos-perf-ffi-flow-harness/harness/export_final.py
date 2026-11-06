#!/usr/bin/env python3
import base64
import hashlib
import html
import io
import json
import pathlib
import shutil
import tarfile

from run_leg import accepted_directory

root = pathlib.Path("/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a")
assert (root / "provenance/CAMPAIGN_COMPLETE").exists()
share = root / "export/share-v3"
share.mkdir(exist_ok=False)
for name in ("clean-comparison", "ffi-comparison", "comparison-v1"):
    source = root / ("comparison-v2" if name == "comparison-v1" else name)
    for row in (source / "SHA256SUMS").read_text().splitlines():
        digest, relative = row.split("  ", 1)
        assert hashlib.sha256((source / relative).read_bytes()).hexdigest() == digest
    shutil.copytree(source, share / name)
provenance = share / "provenance"
provenance.mkdir()
for category in ("clean", "proof", "flow"):
    for kind in ("go-v2-native", "rust-bench"):
        source = root / "build/bin" / category / f"{kind}.provenance.json"
        shutil.copyfile(source, provenance / f"{category}-{kind}.json")
shutil.copyfile(root / "run-spec.json", provenance / "run-spec.json")

links = []
for workers in (32, 64):
    directory = accepted_directory(root, workers, "go-diagnostic", 1, "go-v2-native")
    output = share / f"go-diagnostics/w{workers}"
    output.mkdir(parents=True)
    for path in sorted((directory / "derived").glob("*-top.txt")):
        shutil.copyfile(path, output / path.name)
    shutil.copyfile(directory / "derived/go-diagnostics-summary.json", output / "summary.json")
    scheduler = accepted_directory(root, workers, "scheduler", 1, "go-v2-native")
    shutil.copyfile(scheduler / "raw/sched-latency.txt", output / "os-scheduler-latency.txt")
    shutil.copyfile(scheduler / "derived/scheduler-scope.json", output / "os-scheduler-scope.json")
    items = "".join(f'<li><a href="{path.relative_to(share).as_posix()}">{html.escape(path.name)}</a></li>'
                    for path in sorted(output.iterdir()))
    links.append(f"<h2>W{workers}: Go scheduler, allocation, and GC evidence</h2><ul>{items}</ul>")

(share / "index.html").write_text(
    '<!doctype html><meta charset="utf-8"><title>Cosmos W32/W64 results</title>'
    '<style>body{font:16px system-ui;max-width:1060px;margin:32px auto;line-height:1.55}.notice{background:#fff3cd;padding:16px}li{margin:5px 0}</style>'
    '<h1>Cosmos W32/W64: clean scores and diagnostics</h1>'
    '<p class="notice">Binary encoding explicitly disabled and pre-conversion text bodies verified. '
    'Clean-score, instrumented flow, perf, Go pprof/trace, and OS scheduler evidence remain separate.</p>'
    '<ul><li><a href="clean-comparison/dashboard.html">Clean throughput, CPU/read, latency and RSS</a></li>'
    '<li><a href="ffi-comparison/dashboard.html">Actual completion-queue depths, rates and batch sizes</a></li>'
    '<li><a href="comparison-v1/dashboard.html">Combined comparison and CPU-attribution flywheels</a></li>'
    '<li><a href="comparison-v1/summary.json">Full decoded comparison metadata and diagnostic ledgers</a></li></ul>'
    + "".join(links) +
    '<p>Native flamegraphs, raw traces, exact executables and debug symbols are preserved privately on the VM. '
    'They are deliberately excluded from this teammate-facing package.</p>')
(share / "README.md").write_text(
    "# Cosmos W32/W64 result package\n\n"
    "Extract the complete ZIP and open index.html. Relative links require keeping the folders together.\n\n"
    "Three legs: Go v2 over native C ABI 0.3.0, standalone Rust SDK, and direct Rust driver. "
    "Rust source 77ebe786ad1012bc9ac56105e08c05fb4c65b321; optimized Go base b20b0fd00ac20c78bdb2d138fb6960b35e428af0, "
    "with the explicit experimental Linux native-0.3.0 header/version port. Tokio=2, GOMAXPROCS=4, four-vCPU VM. "
    "Two mirrored clean repeats at 32 and 64 workers, 10 s excluded warmup and 60 s measurement. "
    "SDK diagnostics Summary; not disabled. Gateway V1 HTTP/2, binary encoding disabled and live-verified. "
    "Point reads against 1024 exact-1024-byte seeded documents before service metadata; shared-key auth passed through inherited private FD. "
    "All accepted measured errors/retries/throttles are zero. No Go v1 leg in this new comparison.\n\n"
    "CPU/read is process user+system CPU divided by successful reads, not wall latency or hardware cycles. "
    "Latency aggregates are means of repeat percentiles, not pooled percentiles. "
    "Do not add pprof/perf/process CPU denominators or interpret batch size as queue depth. "
    "Publication-to-resume delay includes publication plus scheduling, not pure scheduler time. "
    "First-polled-unfinished includes tasks awaiting I/O. "
    "OS scheduler recording is PID-attached, not an exhaustive system-wide wakeup trace. "
    "Go heap allocation totals include startup/warmup and exclude native allocations; RSS is whole-process. "
    "Go CPU pprof does not provide complete attribution for foreign Rust threads; use native perf separately. "
    "Historical binary-wire/text-output scores are not matched controls, and a full Rust-main upgrade does not isolate a session-token commit.\n\n"
    "Teammate-facing package excludes private Azure configuration, resource identifiers, credentials, raw perf stack dumps, "
    "raw execution traces, build caches and executables. Raw evidence and exact symbols remain private on the VM. "
    "Do not publish the raw evidence bundle without a separate privacy review.\n")

files = sorted(path for path in share.rglob("*") if path.is_file())
for path in files:
    assert path.suffix in {".html", ".svg", ".md", ".json", ".csv", ".txt"} or path.name == "SHA256SUMS"
    text = path.read_text(errors="strict")
    for forbidden in ("go-perf-20261020204009", "go-perf-ngw-d4-20261027-a",
                      "AccountKey=", "SharedAccessSignature=", "Authorization: Bearer "):
        assert forbidden not in text, f"forbidden private data marker in {path.name}"
with (share / "SHA256SUMS").open("x") as handle:
    for path in files:
        handle.write(f"{hashlib.sha256(path.read_bytes()).hexdigest()}  {path.relative_to(share).as_posix()}\n")
files = sorted(path for path in share.rglob("*") if path.is_file())
buffer = io.BytesIO()
with tarfile.open(fileobj=buffer, mode="w:xz", preset=9) as archive:
    for path in files:
        data = path.read_bytes()
        info = tarfile.TarInfo("share-v3/" + path.relative_to(share).as_posix())
        info.size, info.mode, info.mtime = len(data), 0o444, 0
        archive.addfile(info, io.BytesIO(data))
data = buffer.getvalue()
encoded = base64.b64encode(data).decode()
prior = root / "export/final-share.tar.xz"
if prior.exists():
    backup = prior.with_name("final-share-superseded-" + hashlib.sha256(prior.read_bytes()).hexdigest()[:12] + ".tar.xz")
    if not backup.exists():
        backup.write_bytes(prior.read_bytes())
(root / "export/final-share.tar.xz").write_bytes(data)
(root / "export/final-share.b64").write_text(encoded)
metadata = {"bytes": len(data), "sha256": hashlib.sha256(data).hexdigest(),
            "base64Bytes": len(encoded), "fileCount": len(files), "rootFolder": "share-v3"}
(root / "export/final-share-meta.json").write_text(json.dumps(metadata, indent=2) + "\n")
print("FINAL_META:" + json.dumps(metadata, separators=(",", ":")))
print("FINAL_SHARE_EXPORT_READY")
