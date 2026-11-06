#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
python3 - <<'PY'
import json, pathlib, re
root = pathlib.Path("/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a")
rejected = sorted((root / "run").rglob("REJECTED.json"), key=lambda p: p.stat().st_mtime)
for record in rejected[-2:]:
    directory = record.parent
    print("REJECTED_LEG", directory.relative_to(root))
    print(record.read_text())
    for name in ("stderr.txt", "perf-record.stderr"):
        path = directory / "raw" / name
        if path.exists():
            text = path.read_text(errors="replace")[-1500:]
            text = re.sub(r"https?://[^\s\"')]+", "<url>", text)
            text = re.sub(r"(?i)(authorization|access_token|sig|accountkey)\s*[:=]\s*\S+", r"\1=<redacted>", text)
            print(name, text)
    barrier = directory / "raw/barriers"
    print("barriers", sorted(p.name for p in barrier.iterdir()) if barrier.exists() else [])
print("cleanup", (root / "provenance/cleanup.json").read_text() if (root / "provenance/cleanup.json").exists() else "absent")
print("REJECTED_CONTEXT_INSPECTED")
PY
