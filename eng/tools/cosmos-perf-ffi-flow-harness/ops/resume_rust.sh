#!/usr/bin/env bash
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
inputs=$(cd "$(dirname "$0")" && pwd)
[[ -f "$root/provenance/RUST_SOURCE_PREPARED" && ! -e "$root/provenance/RUST_BUILD_COMPLETE" ]]
cp "$inputs/rust_main.rs" "$root/src/rust/sdk/cosmos/cosmos_perf_flow/src/main.rs"
cp "$inputs/wire_proof.rs" "$root/src/rust/sdk/cosmos/azure_data_cosmos_driver/src/perf_wire.rs"
cp "$inputs/flow.rs" "$root/src/rust/sdk/cosmos/azure_data_cosmos_driver_native/src/perf_flow.rs"
python3 - "$root/src/rust/sdk/cosmos/azure_data_cosmos_driver/src/driver/cosmos_driver.rs" "$inputs" <<'PY'
import pathlib, sys
sys.path.insert(0, sys.argv[2])
from patch_wire import normalize_wire_call
path = pathlib.Path(sys.argv[1])
path.write_text(normalize_wire_call(path.read_text()))
PY
for class in clean proof flow; do
  if [[ -f "$root/build/rust/$class/build.log" ]]; then
    mv "$root/build/rust/$class/build.log" \
      "$root/build/rust/$class/build-attempt-$(date -u +%Y%m%dT%H%M%SZ).log"
  fi
done
exec /bin/bash "$inputs/build_rust.sh"
