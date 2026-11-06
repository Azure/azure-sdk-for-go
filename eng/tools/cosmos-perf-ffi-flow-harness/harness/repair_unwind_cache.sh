#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
libc=/usr/lib/x86_64-linux-gnu/libc.so.6
id=$(readelf -n "$libc" | awk '/Build ID:/ {print $3}')
debug=/mnt/cosmos-perf/cpu-isolation-synthetic-20261105-d/attribution-symbols/libc6-dbg/extracted/usr/lib/debug/.build-id/${id:0:2}/${id:2}.debug
cache=$(readlink -f "$root/tools/perf-home/.debug/.build-id/${id:0:2}/${id:2}")
[[ -d "$cache" && -f "$cache/elf" && -s "$debug" ]]
[[ "$cache" == "$root/tools/perf-home/.debug/"* ]]
[[ "$(readelf -n "$debug" 2>/dev/null | awk '/Build ID:/ {print $3}')" == "$id" ]]
if [[ ! -e "$cache/elf-debug-only-superseded" ]]; then
  cp "$cache/elf" "$cache/elf-debug-only-superseded"
fi
cp "$libc" "$cache/elf.runtime-next"
mv "$cache/elf.runtime-next" "$cache/elf"
cp "$debug" "$cache/debug"
readelf -S "$cache/elf" | grep -E 'eh_frame.*PROGBITS'
sha256sum "$libc" "$cache/elf" "$debug" "$cache/debug" > "$root/provenance/unwind-cache-sha256.txt"
python3 - "$root" <<'PY'
import hashlib, json, os, pathlib, subprocess, sys
root = pathlib.Path(sys.argv[1])
inputs = root / "tools/inputs-resume-run-bf69585884b0"
sys.path.insert(0, str(inputs))
from calibrate import calculate
index = json.loads((root / "provenance/accepted-legs.json").read_text())
environment = os.environ.copy()
environment["HOME"] = str(root / "tools/perf-home")
environment["DEBUGINFOD_URLS"] = ""
for leg in ("go-v2-native", "rust-sdk", "rust-driver"):
    for repeat, mode in ((1, "fp"), (2, "dwarf,16384")):
        key = f"w32/perf-calibration/repeat-{repeat}/{leg}"
        directory = root / index[key]
        raw = directory / "raw"
        derived = directory / "derived/unwind-v2"
        derived.mkdir(exist_ok=False)
        bounds = json.loads((raw / "barriers/measurement-bounds.json").read_text())
        stamp = lambda ns: f"{ns//10**9}.{ns%10**9:09d}"
        interval = stamp(bounds["monotonicStartNs"]) + "," + stamp(bounds["monotonicEndNs"])
        with (derived / "perf-script.txt").open("wb") as output, (derived / "stderr.txt").open("wb") as errors:
            subprocess.run(["perf", "script", "--no-inline", "--ns", "--show-lost-events",
                            "--time", interval, "-i", str(raw / "perf.data")],
                           stdout=output, stderr=errors, env=environment, check=True)
        quality = calculate((derived / "perf-script.txt").read_text(),
                            (raw / "perf-record.stderr").read_text())
        (derived / "quality.json").write_text(json.dumps(quality, indent=2) + "\n")
        metadata = {
            "perfScript": "derived/unwind-v2/perf-script.txt",
            "quality": "derived/unwind-v2/quality.json",
            "rawPerfSHA256": hashlib.sha256((raw / "perf.data").read_bytes()).hexdigest(),
            "rawEvidenceUnmodified": True,
            "correction": "runtime libc ELF restored; debug-only ELF no longer used for unwind",
        }
        (directory / "derived/analysis-source.json").write_text(json.dumps(metadata, indent=2) + "\n")
        print(leg, mode, json.dumps(quality, separators=(",", ":")))
print("UNWIND_CACHE_REPAIR_AND_REPROCESS_COMPLETE")
PY
