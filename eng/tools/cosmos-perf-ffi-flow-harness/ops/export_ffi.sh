#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
python3 - <<'PY'
import base64, hashlib, io, json, pathlib, tarfile
root = pathlib.Path("/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a")
source = root / "ffi-comparison"
for line in (source / "SHA256SUMS").read_text().splitlines():
    digest, name = line.split("  ", 1)
    assert hashlib.sha256((source / name).read_bytes()).hexdigest() == digest
summary = json.loads((source / "summary.json").read_text())
keys = ("workers", "success", "submissionsPerSecond", "completionBatchesPerSecond",
        "completionsPerSecond", "meanCompletionBatchSize", "batchHistogram", "nativeReadyHigh",
        "nativeAdmittedUndrainedHigh", "nativePolledUnfinishedHigh",
        "deliveryToResumeMeanUs", "deliveryToResumeMaxUs", "diagnosticCPUUsPerRead")
rows = [{key: value[key] for key in keys} for value in summary["goV2"]]
buffer = io.BytesIO()
with tarfile.open(fileobj=buffer, mode="w:xz", preset=6) as archive:
    for path in sorted(source.iterdir()):
        data = path.read_bytes()
        info = tarfile.TarInfo("ffi-comparison/" + path.name)
        info.size, info.mode, info.mtime = len(data), 0o444, 0
        archive.addfile(info, io.BytesIO(data))
data = buffer.getvalue()
encoded = base64.b64encode(data).decode()
export = root / "export"
(export / "ffi-counter-checkpoint.tar.xz").write_bytes(data)
(export / "ffi-counter-checkpoint.b64").write_text(encoded)
print("FFI_ROWS:" + json.dumps(rows, separators=(",", ":")))
print("BUNDLE_META:" + json.dumps({
    "bytes":len(data), "sha256":hashlib.sha256(data).hexdigest(), "base64Bytes":len(encoded),
    "files":sorted(path.name for path in source.iterdir()),
}, separators=(",", ":")))
print("FFI_CHECKPOINT_EXPORT_READY")
PY
