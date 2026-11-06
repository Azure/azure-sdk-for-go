#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
source /mnt/cosmos-perf/tokio234-gatewayv1-20261105-a/tools/common.sh
python3 - "$database" "$container" <<'PY'
import json, sys
print(json.dumps({"database": sys.argv[1], "container": sys.argv[2]}))
print("DATABASE_CONTAINER_NAMES_VERIFIED")
PY
