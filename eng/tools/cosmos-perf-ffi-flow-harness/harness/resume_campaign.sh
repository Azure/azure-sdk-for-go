#!/usr/bin/env bash
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
inputs=$(cd "$(dirname "$0")" && pwd)
[[ -f "$root/provenance/GO_BUILD_COMPLETE" && ! -e "$root/provenance/CAMPAIGN_COMPLETE" ]]
python3 - "$root" "$inputs" <<'PY'
import hashlib, json, os, pathlib, stat, sys, tempfile
root = pathlib.Path(sys.argv[1])
sys.path.insert(0, sys.argv[2])
from run_leg import seal
for path in (root / "build/bin").glob("*/*"):
    if path.name not in ("go-v2-native", "rust-bench"):
        continue
    value = json.loads(pathlib.Path(str(path) + ".provenance.json").read_text())
    assert hashlib.sha256(path.read_bytes()).hexdigest() == value["binary"]["sha256"]
    path.chmod(0o555)
with tempfile.TemporaryDirectory(dir=root / "build/tmp") as temporary:
    temporary = pathlib.Path(temporary)
    executable = temporary / "staged"
    executable.write_text("#!/bin/sh\nexit 0\n")
    executable.chmod(0o555)
    raw = temporary / "raw"
    raw.mkdir()
    os.link(executable, raw / "binary")
    seal(raw)
    assert executable.stat().st_mode & 0o111, "hard-link sealing removed execution"
    assert executable.stat().st_mode & 0o222 == 0, "sealed binary became writable"
for source in (root / "private/cleanup.stderr", root / "provenance/cleanup.json"):
    if source.exists():
        digest = hashlib.sha256(source.read_bytes()).hexdigest()[:12]
        backup = source.with_name(source.name + ".superseded-" + digest)
        if not backup.exists():
            backup.write_bytes(source.read_bytes())
print("EXECUTABLE_HASHES_AND_HARDLINK_SEALING_VERIFIED")
PY
exec /bin/bash "$inputs/run_campaign.sh" "$@"
