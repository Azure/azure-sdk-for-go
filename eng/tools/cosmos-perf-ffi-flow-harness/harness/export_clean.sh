#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
python3 - <<'PY'
import base64, hashlib, io, json, pathlib, tarfile
root = pathlib.Path("/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a")
assert (root / "provenance/CLEAN_SCORE_COMPLETE").exists()
cleanup = json.loads((root / "provenance/cleanup.json").read_text())
assert cleanup["deleted"] == 1024 and cleanup["nonOwnedSkipped"] == 0
source = root / "clean-comparison"
for line in (source / "SHA256SUMS").read_text().splitlines():
    digest, name = line.split("  ", 1)
    assert hashlib.sha256((source / name).read_bytes()).hexdigest() == digest
buffer = io.BytesIO()
with tarfile.open(fileobj=buffer, mode="w:xz", preset=6) as archive:
    for path in sorted(source.iterdir()):
        assert path.is_file() and not path.is_symlink()
        data = path.read_bytes()
        info = tarfile.TarInfo("clean-comparison/" + path.name)
        info.size = len(data)
        info.mode = 0o444
        info.mtime = 0
        archive.addfile(info, io.BytesIO(data))
data = buffer.getvalue()
export = root / "export"
export.mkdir(exist_ok=True)
path = export / "clean-score-checkpoint.tar.xz"
path.write_bytes(data)
encoded = base64.b64encode(data).decode()
(export / "clean-score-checkpoint.b64").write_text(encoded)
summary = json.loads((source / "summary.json").read_text())
rows = [{
    "workers": row["workers"], "leg": row["leg"],
    "readsPerSecond": round(row["throughput"], 2),
    "cpuUsPerRead": round(row["cpuMicrosecondsPerSuccess"], 2),
    "busyCores": round(row["busyCores"], 3),
    "meanRepeatP50Ms": round(row["meanRepeatP50Ms"], 3),
    "meanRepeatP99Ms": round(row["meanRepeatP99Ms"], 3),
    "meanRSSMiB": round(row["meanRSSMiB"], 2),
    "repeatThroughputSpreadPercent": round(row["repeatThroughputSpreadPercent"], 2),
} for row in summary["aggregates"]]
print("CLEAN_ROWS:" + json.dumps(rows, separators=(",", ":")))
print("BUNDLE_META:" + json.dumps({
    "bytes": len(data), "sha256": hashlib.sha256(data).hexdigest(), "base64Bytes": len(encoded),
    "files": sorted(path.name for path in source.iterdir()), "cleanupDeleted": cleanup["deleted"],
}, separators=(",", ":")))
print("CLEAN_CHECKPOINT_EXPORT_READY")
PY
