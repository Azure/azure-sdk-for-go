#!/usr/bin/env bash
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words Eeuo esac

set -Eeuo pipefail

usage() {
  echo "usage: stage-adapters.sh --go-build <dir> --rust-build <dir> --output <dir>" >&2
  exit 2
}

go_build= rust_build= output=
while [[ $# -gt 0 ]]; do
  case "$1" in
    --go-build) go_build=$(readlink -f "$2"); shift 2 ;;
    --rust-build) rust_build=$(readlink -f "$2"); shift 2 ;;
    --output) output=$(readlink -m "$2"); shift 2 ;;
    *) usage ;;
  esac
done
[[ -d "$go_build/bin" && -d "$rust_build/bin" && -n "$output" ]] || usage
[[ ! -e "$output" ]] || { echo "refusing to overwrite $output" >&2; exit 1; }

mkdir -p "$output/bin"
for name in cosmos-perf-go-v1 cosmos-perf-go-v2-native; do
  install -m 0755 "$go_build/bin/$name" "$output/bin/$name"
  install -m 0444 "$go_build/bin/$name.provenance.json" "$output/bin/$name.provenance.json"
done
install -m 0755 "$rust_build/bin/cosmos-perf-rust" "$output/bin/cosmos-perf-rust"
install -m 0444 "$rust_build/bin/cosmos-perf-rust.provenance.json" \
  "$output/bin/cosmos-perf-rust.provenance.json"
(
  cd "$output"
  sha256sum bin/* > SHA256SUMS
)
chmod -R a-w "$output"
