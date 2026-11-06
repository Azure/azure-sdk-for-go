#!/usr/bin/env bash
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
inputs=$(cd "$(dirname "$0")" && pwd)
[[ -f "$root/provenance/RUST_BUILD_COMPLETE" && ! -e "$root/provenance/GO_BUILD_COMPLETE" ]]
sdk=$root/build/go/clean/sdk
[[ -f "$sdk/driver.go" && -f "$root/build/go/clean/sdk-validation.log" ]]
grep -q 'const nativeDriverVersion = "0.3.0"' "$sdk/driver.go"
cmp "$sdk/azurecosmosdriver.h" "$root/build/rust/clean/azurecosmosdriver.h"
python3 - "$inputs" "$sdk" <<'PY'
import pathlib, sys
sys.path.insert(0, sys.argv[1])
from patch_go import patch_user_agent_test
sdk = pathlib.Path(sys.argv[2])
text = (sdk / "user_agent_native_test.go").read_text()
if "azsdk-rust-cosmos-driver/" in text:
    patch_user_agent_test(sdk)
else:
    assert r"drv=1\.1\.0-beta\.3" in text, "unexpected user-agent test state"
PY
date -u +%FT%TZ > "$root/build/go/clean/SOURCE_PREPARED"
exec /bin/bash "$inputs/build_go.sh"
