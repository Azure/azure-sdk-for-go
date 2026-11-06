#!/usr/bin/env bash
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
inputs=$(cd "$(dirname "$0")" && pwd)
[[ -f "$root/provenance/GO_BUILD_COMPLETE" && -f "$root/provenance/CLEAN_SCORE_COMPLETE" ]]
[[ ! -e "$root/provenance/CAMPAIGN_COMPLETE" ]]
export HOME=/root GOCACHE="$root/build/go-cache"
python3 - "$root" <<'PY'
import json, pathlib, sys
root = pathlib.Path(sys.argv[1])
value = json.loads((root / "provenance/cleanup.json").read_text())
assert value["deleted"] == 1024 and value["nonOwnedSkipped"] == 0
PY
python3 "$inputs/analyze_go_diagnostics.py" "$root"
python3 "$inputs/profile_attribution.py" "$root"
python3 "$inputs/report.py" "$root"
date -u +%FT%TZ > "$root/provenance/CAMPAIGN_COMPLETE"
echo MATCHED_W32_W64_THREE_LEG_CAMPAIGN_COMPLETE
