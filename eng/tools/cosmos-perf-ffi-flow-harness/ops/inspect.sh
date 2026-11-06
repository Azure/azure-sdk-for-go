#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
python3 - <<'PY'
import hashlib, json, os, pathlib, shutil, subprocess
root = pathlib.Path("/mnt/cosmos-perf/cpu-isolation-synthetic-20261105-d")
active = []
for p in pathlib.Path("/proc").iterdir():
    if not p.name.isdigit():
        continue
    try:
        exe = os.readlink(p / "exe")
    except (FileNotFoundError, PermissionError, ProcessLookupError):
        continue
    if exe.startswith("/mnt/cosmos-perf/") and ("cosmos-perf" in pathlib.Path(exe).name or pathlib.Path(exe).name in ("control", "option-b", "c0-control")):
        active.append({"pid": int(p.name), "name": pathlib.Path(exe).name})
sdk = root / "src/go/control"
print(json.dumps({
    "freeGiB": round(shutil.disk_usage(root).free / 2**30, 3),
    "cpus": os.cpu_count(),
    "active": active,
    "sdkFiles": {p.name: hashlib.sha256(p.read_bytes()).hexdigest() for p in sdk.glob("*.go") if p.name in ("driver_native.go", "reactor_native.go", "operation_native.go", "native_cgo.go")},
    "sdkExists": sdk.is_dir(),
    "mainKnown": subprocess.run(["git", "-C", "/mnt/cosmos-perf/point-read-first-pass-20261102-i/src/rust-sdk", "cat-file", "-t", "77ebe786ad1012bc9ac56105e08c05fb4c65b321"], capture_output=True, text=True).stdout.strip(),
}, sort_keys=True))
PY
echo FLOW_INSPECTION_COMPLETE
