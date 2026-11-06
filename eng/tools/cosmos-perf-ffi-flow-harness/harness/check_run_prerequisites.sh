#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
[[ -f "$root/provenance/GO_BUILD_COMPLETE" ]]
python3 - "$root" <<'PY'
import json, os, pathlib, shutil, sys
root = pathlib.Path(sys.argv[1])
print(json.dumps({
    "availableCPUs": sorted(os.sched_getaffinity(0)),
    "freeGiB": round(shutil.disk_usage(root).free / 2**30, 2),
    "memoryController": pathlib.Path("/sys/fs/cgroup/memory.current").exists(),
    "libcDebugSymbols": pathlib.Path("/mnt/cosmos-perf/cpu-isolation-synthetic-20261105-d/attribution-symbols/libc6-dbg/extracted/usr/lib/debug/.build-id/6a/6cb5a80fd5e5bb5bf3aeaa111add98e88843a5.debug").exists(),
}))
PY
git -C "$root/src/rust" grep -n 'pub fn request_count' 77ebe786ad1012bc9ac56105e08c05fb4c65b321 \
  -- sdk/cosmos/azure_data_cosmos/src/diagnostics sdk/cosmos/azure_data_cosmos_driver/src/diagnostics || true
perf --version
echo RUN_PREREQUISITES_INSPECTED
