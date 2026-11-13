#!/usr/bin/env bash
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words esac
# cSpell:words Eeuo

set -Eeuo pipefail

usage() {
  echo "usage: run-go-diagnostic-leg.sh --leg <go-v1|go-v2-native> --workers <8|16> --run-dir <dir> --config <json> --key-fd <n> -- <command...>" >&2
  exit 2
}

leg= workers= run_dir= config= key_fd=
while [[ $# -gt 0 ]]; do
  case "$1" in
    --leg) leg=$2; shift 2 ;;
    --workers) workers=$2; shift 2 ;;
    --run-dir) run_dir=$2; shift 2 ;;
    --config) config=$2; shift 2 ;;
    --key-fd) key_fd=$2; shift 2 ;;
    --) shift; break ;;
    *) usage ;;
  esac
done
[[ "$leg" =~ ^(go-v1|go-v2-native)$ && "$workers" =~ ^(8|16)$ ]] || usage
[[ "$key_fd" =~ ^[3-9][0-9]*$ && -r "/proc/self/fd/$key_fd" ]] || usage
[[ -f "$config" && $# -gt 0 ]] || usage

leg_dir=$run_dir/w$workers/go-diagnostic/repeat-1/$leg
raw=$leg_dir/raw
[[ ! -e "$leg_dir" ]] || { echo "refusing to overwrite $leg_dir" >&2; exit 1; }
mkdir -p "$raw" "$leg_dir/derived" "$run_dir/config"
cp "$config" "$run_dir/config/effective.json"

export COSMOS_PERF_CONFIG=$run_dir/config/effective.json
export COSMOS_PERF_WORKERS=$workers
export COSMOS_PERF_ADAPTER_RESULT=$raw/adapter-result.json
export COSMOS_PERF_KEY_FD=$key_fd
export COSMOS_PERF_CPU_PROFILE=$raw/cpu.pprof
export COSMOS_PERF_HEAP_PROFILE=$raw/heap.pprof
export COSMOS_PERF_RUNTIME_TRACE=$raw/runtime.trace
export COSMOS_PERF_BLOCK_PROFILE=$raw/block.pprof
export COSMOS_PERF_MUTEX_PROFILE=$raw/mutex.pprof

binary=$(readlink -f "$1")
sha256sum "$binary" > "$raw/binary.sha256"
readelf -n "$binary" > "$raw/binary-build-id.txt" 2>&1 || true
cp --preserve=mode,timestamps "$binary" "$raw/binary"
[[ -f "$binary.provenance.json" ]] || { echo "missing $binary.provenance.json" >&2; exit 1; }
cp "$binary.provenance.json" "$leg_dir/provenance.json"
"$@" > "$raw/stdout.txt" 2> "$raw/stderr.txt"
for evidence in cpu.pprof heap.pprof runtime.trace block.pprof mutex.pprof; do
  [[ -s "$raw/$evidence" ]] || { echo "missing $evidence" >&2; exit 1; }
done
python3 "$(dirname "$0")/convert-adapter-output.py" \
  --adapter-result "$raw/adapter-result.json" --raw-dir "$raw" \
  --leg "$leg" --run-class go-diagnostic --workers "$workers" --repeat 1 \
  --output "$leg_dir/metrics.json"
sha256sum "$raw"/* > "$raw/SHA256SUMS"
chmod -R a-w "$raw"
