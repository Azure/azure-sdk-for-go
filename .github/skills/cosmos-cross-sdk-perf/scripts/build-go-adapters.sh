#!/usr/bin/env bash
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words Eeuo esac gcflags ldflags tmpl

set -Eeuo pipefail

usage() {
  echo "usage: build-go-adapters.sh --go-v2-checkout <dir> --go-v2-commit <sha> --native-module <dir> --native-commit <sha> --output <dir>" >&2
  exit 2
}

go_v2_checkout= go_v2_commit= native_module= native_commit= output=
while [[ $# -gt 0 ]]; do
  case "$1" in
    --go-v2-checkout) go_v2_checkout=$(readlink -f "$2"); shift 2 ;;
    --go-v2-commit) go_v2_commit=$2; shift 2 ;;
    --native-module) native_module=$(readlink -f "$2"); shift 2 ;;
    --native-commit) native_commit=$2; shift 2 ;;
    --output) output=$(readlink -m "$2"); shift 2 ;;
    *) usage ;;
  esac
done
[[ "$go_v2_commit" =~ ^[0-9a-f]{40}$ && "$native_commit" =~ ^[0-9a-f]{40}$ ]] || usage
[[ -d "$go_v2_checkout" && -d "$native_module" && -n "$output" ]] || usage
[[ ! -e "$output" ]] || { echo "refusing to overwrite $output" >&2; exit 1; }

skill=$(cd "$(dirname "$0")/.." && pwd)
adapter=$skill/adapters/go
for checkout in "$go_v2_checkout" "$native_module"; do
  [[ -z "$(git -C "$checkout" status --short)" ]] || {
    echo "$checkout must be clean" >&2
    exit 1
  }
done
[[ "$(git -C "$go_v2_checkout" rev-parse HEAD)" == "$go_v2_commit" ]] || {
  echo "Go v2 checkout commit mismatch" >&2
  exit 1
}
[[ "$(git -C "$native_module" rev-parse HEAD)" == "$native_commit" ]] || {
  echo "native distribution checkout commit mismatch" >&2
  exit 1
}

mkdir -p "$output/work/go-v1" "$output/work/go-v2" "$output/bin"
cp "$adapter"/*.go "$output/work/go-v1/"
cp "$adapter"/*.go "$output/work/go-v2/"
cp "$adapter/sdk_v2.go.tmpl" "$output/work/go-v2/sdk_v2.go"

cat > "$output/work/go-v1/go.mod" <<'EOF'
module cosmos-perf-go-v1

go 1.23.0

require github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos v1.5.0
EOF
(
  cd "$output/work/go-v1"
  go mod download
  CGO_ENABLED=0 go build -tags go_v1 -buildvcs=false \
    -ldflags=-compressdwarf=false -o "$output/bin/cosmos-perf-go-v1" .
)

cat > "$output/work/go-v2/go.mod" <<EOF
module cosmos-perf-go-v2

go 1.27.0

require github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2 v2.0.0

replace github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2 => $go_v2_checkout/sdk/data/azcosmos
replace github.com/Azure/azure-cosmos-driver/linux/amd64 => $native_module/linux/amd64
EOF
(
  cd "$output/work/go-v2"
  go mod download
  CGO_ENABLED=1 go build -tags go_v2 -buildvcs=false \
    -ldflags=-compressdwarf=false -o "$output/bin/cosmos-perf-go-v2-native" .
)

write_provenance() {
  local leg=$1 binary=$2 repository=$3 commit=$4 tree=$5 dependencies=$6 native=$7
  local build_id
  build_id=$(go tool buildid "$binary")
  python3 - "$leg" "$binary" "$repository" "$commit" "$tree" "$dependencies" "$native" "$build_id" <<'PY'
import hashlib
import json
import pathlib
import subprocess
import sys

leg, binary, repository, commit, tree, dependencies, native, build_id = sys.argv[1:]
binary_path = pathlib.Path(binary)
value = {
    "leg": leg,
    "source": {
        "repository": repository,
        "commit": commit,
        "tree": tree,
        "dirty": False,
    },
    "compiler": subprocess.check_output(["go", "version"], text=True).strip(),
    "buildCommand": "CGO_ENABLED=0 go build -tags go_v1 -ldflags=-compressdwarf=false"
    if leg == "go-v1"
    else "CGO_ENABLED=1 go build -tags go_v2 -ldflags=-compressdwarf=false",
    "dependencyVersions": json.loads(dependencies),
    "binary": {
        "sha256": hashlib.sha256(binary_path.read_bytes()).hexdigest(),
        "buildId": build_id,
        "debugSymbols": "embedded",
    },
}
if native:
    value["nativeDriver"] = json.loads(native)
binary_path.with_suffix(binary_path.suffix + ".provenance.json").write_text(
    json.dumps(value, indent=2, sort_keys=True) + "\n", encoding="utf-8"
)
PY
}

adapter_commit=$(git -C "$skill" rev-parse HEAD)
adapter_tree=$(git -C "$skill" rev-parse HEAD^{tree})
write_provenance go-v1 "$output/bin/cosmos-perf-go-v1" \
  https://github.com/Azure/azure-sdk-for-go.git "$adapter_commit" "$adapter_tree" \
  '{"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos":"v1.5.0"}' ''
write_provenance go-v2-native "$output/bin/cosmos-perf-go-v2-native" \
  https://github.com/Azure/azure-sdk-for-go.git "$go_v2_commit" \
  "$(git -C "$go_v2_checkout" rev-parse HEAD^{tree})" \
  '{"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2":"local-pinned-commit"}' \
  "{\"sourceCommit\":\"62192bf69dac1a3384505971d7deef19eb8282df\",\"distributionCommit\":\"$native_commit\"}"

for binary in "$output"/bin/cosmos-perf-go-*; do
  [[ "$binary" == *.json ]] && continue
  sha256sum "$binary"
  go tool buildid "$binary"
done > "$output/SHA256SUMS"
rm -rf "$output/work"
chmod -R a-w "$output"
