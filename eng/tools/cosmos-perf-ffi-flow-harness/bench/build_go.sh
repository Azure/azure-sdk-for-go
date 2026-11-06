#!/usr/bin/env bash
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
inputs=$(cd "$(dirname "$0")" && pwd)
distribution=/mnt/cosmos-perf/point-read-first-pass-20261102-i/src/native-distribution/linux/amd64
go=/usr/local/go/bin/go
[[ -f "$root/provenance/RUST_BUILD_COMPLETE" && ! -e "$root/provenance/GO_BUILD_COMPLETE" ]]
export HOME=/root GOPATH=/mnt/cosmos-perf/cache/gopath GOMODCACHE=/mnt/cosmos-perf/cache/gomod
export GOCACHE="$root/build/go-cache" TMPDIR="$root/build/tmp" CGO_ENABLED=1
mkdir -p "$GOCACHE" "$TMPDIR" "$root/build/go" "$root/build/bin"

for class in clean proof flow; do
  output=$root/build/go/$class
  sdk=$output/sdk
  adapter=$output/adapter
  native=$output/native
  if [[ -f "$output/CLASS_COMPLETE" ]]; then continue; fi
  if [[ -e "$output" && ! -f "$output/SOURCE_PREPARED" ]]; then
    echo "Unrecognized partial Go preparation; preserving it without overwrite: $class" >&2
    exit 1
  fi
  if [[ ! -f "$output/SOURCE_PREPARED" ]]; then
  mkdir -p "$sdk" "$adapter" "$native" "$root/build/bin/$class"
  cp -a "$root/src/go/clean/." "$sdk/"
  cp "$root/src/adapter-base/"*.go "$root/src/adapter-base/go.mod" \
    "$root/src/adapter-base/go.sum" "$adapter/"
  find "$distribution" -maxdepth 1 -type f \( -name '*.go' -o -name go.mod -o -name go.sum \) \
    -exec cp {} "$native/" \;
  cp "$root/build/rust/$class/azurecosmosdriver.h" "$native/azurecosmosdriver.h"
  cp "$root/build/rust/$class/azurecosmosdriver.h" "$sdk/azurecosmosdriver.h"
  chmod 0444 "$root/build/rust/$class/libazurecosmosdriver.a"
  ln "$root/build/rust/$class/libazurecosmosdriver.a" "$native/libazurecosmosdriver.a"
  python3 - "$inputs" "$sdk" "$adapter" "$class" <<'PY'
import pathlib, sys
sys.path.insert(0, sys.argv[1])
from patch_go import patch_sdk, patch_adapter
diagnostic = sys.argv[4] == "flow"
patch_sdk(pathlib.Path(sys.argv[2]), diagnostic)
patch_adapter(pathlib.Path(sys.argv[3]), diagnostic)
PY
  date -u +%FT%TZ > "$output/SOURCE_PREPARED"
  fi
  if [[ -f "$output/sdk-validation.log" ]]; then
    cp "$output/sdk-validation.log" "$output/sdk-validation-attempt-$(date -u +%Y%m%dT%H%M%SZ).log"
  fi
  (
    cd "$sdk"
    "$go" mod edit -require=github.com/Azure/azure-cosmos-driver/linux/amd64@v0.3.0
    "$go" mod edit "-replace=github.com/Azure/azure-cosmos-driver/linux/amd64=$native"
    "$go" mod tidy
    "$go" fmt ./...
    if [[ "$class" == clean ]]; then
      "$go" test ./...
    elif [[ "$class" == flow ]]; then
      "$go" test -race -run 'Test.*(Pending|Reactor|Await|OperationOptionsToNative|ClientOptionsToNative)' ./...
    fi
  ) > "$output/sdk-validation.log" 2>&1 || {
    echo "GO_SDK_VALIDATION_FAILED class=$class"; tail -n 55 "$output/sdk-validation.log"; exit 1;
  }
  (
    cd "$adapter"
    "$go" mod edit "-replace=github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2=$sdk"
    "$go" mod edit -require=github.com/Azure/azure-cosmos-driver/linux/amd64@v0.3.0
    "$go" mod edit "-replace=github.com/Azure/azure-cosmos-driver/linux/amd64=$native"
    "$go" mod tidy
    "$go" fmt ./...
    "$go" test -tags go_v2 ./...
    "$go" build -tags go_v2 -buildvcs=false -ldflags=-compressdwarf=false \
      -o "$root/build/bin/$class/go-v2-native" .
  ) > "$output/adapter-build.log" 2>&1 || {
    echo "GO_ADAPTER_BUILD_FAILED class=$class"; tail -n 55 "$output/adapter-build.log"; exit 1;
  }
  chmod 0555 "$root/build/bin/$class/go-v2-native" "$root/build/rust/$class/cosmos_perf_flow"
  ln "$root/build/rust/$class/cosmos_perf_flow" "$root/build/bin/$class/rust-bench"
  date -u +%FT%TZ > "$output/CLASS_COMPLETE"
done

(
  cd "$root/build/go/flow/adapter"
  "$go" test -race -v "$inputs/flow.go" "$inputs/flow_test.go"
  "$go" test -v "$inputs/barrier_go.go" "$inputs/barrier_go_test.go"
) > "$root/build/go/helper-tests.log" 2>&1
python3 "$inputs/write_build_provenance.py" "$root"
date -u +%FT%TZ > "$root/provenance/GO_BUILD_COMPLETE"
echo GO_LATEST_NATIVE_THREE_CLASS_BUILD_COMPLETE
