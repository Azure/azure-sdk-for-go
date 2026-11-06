#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
python3 - <<'PY'
import json, pathlib
root = pathlib.Path("/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a")
index = json.loads((root / "provenance/accepted-legs.json").read_text())
for leg in ("go-v2-native", "rust-sdk", "rust-driver"):
    for repeat, mode in ((1, "fp"), (2, "dwarf")):
        key = f"w32/perf-calibration/repeat-{repeat}/{leg}"
        directory = root / index[key]
        quality = json.loads((directory / "derived/quality.json").read_text())
        print(leg, mode, json.dumps(quality, separators=(",", ":")))
        if not quality["valid"]:
            text = (directory / "raw/perf-script.txt").read_text()
            print("FIRST_RECORD", text.split("\n\n", 1)[0][:380].replace("\n", " | "))
print("CALIBRATION_QUALITY_INSPECTED")
PY
