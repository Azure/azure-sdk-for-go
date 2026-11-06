#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
python3 - <<'PY'
import json, pathlib
root = pathlib.Path("/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a")
index = json.loads((root / "provenance/accepted-legs.json").read_text())
print("VM_ROOT", root)
for workers in (32, 64):
    for category, names in (
        ("go-diagnostic", ("runtime.trace", "cpu.pprof", "heap.pprof", "block.pprof", "mutex.pprof")),
        ("scheduler", ("perf-sched.data", "sched-timehist.txt", "sched-latency.txt",
                       "snapshot-measurement-start.json", "snapshot-measurement-end.json")),
    ):
        key = f"w{workers}/{category}/repeat-1/go-v2-native"
        directory = root / index[key]
        assert json.loads((directory / "VALID.json").read_text())["valid"]
        print("LEG", directory.relative_to(root))
        for name in names:
            path = directory / "raw" / name
            assert path.is_file() and path.stat().st_size > 0, name
            print(name, path.stat().st_size)
print("DIAGNOSTIC_FILES_VERIFIED")
PY
