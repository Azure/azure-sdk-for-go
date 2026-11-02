#!/usr/bin/env bash
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words Eeuo esac rollup smaps

set -Eeuo pipefail

usage() {
  echo "usage: run-clean-leg.sh --leg <go-v1|go-v2-native|rust> --workers <8|16> --repeat <n> --run-dir <dir> --config <json> -- <command...>" >&2
  exit 2
}

leg= workers= repeat= run_dir= config=
while [[ $# -gt 0 ]]; do
  case "$1" in
    --leg) leg=$2; shift 2 ;;
    --workers) workers=$2; shift 2 ;;
    --repeat) repeat=$2; shift 2 ;;
    --run-dir) run_dir=$2; shift 2 ;;
    --config) config=$2; shift 2 ;;
    --) shift; break ;;
    *) usage ;;
  esac
done
[[ "$leg" =~ ^(go-v1|go-v2-native|rust)$ ]] || usage
[[ "$workers" =~ ^(8|16)$ && "$repeat" =~ ^[1-9][0-9]*$ ]] || usage
[[ -f "$config" && $# -gt 0 ]] || usage

leg_dir=$run_dir/w$workers/clean-score/repeat-$repeat/$leg
raw=$leg_dir/raw
[[ ! -e "$leg_dir" ]] || { echo "refusing to overwrite $leg_dir" >&2; exit 1; }
mkdir -p "$raw" "$leg_dir/derived" "$run_dir/config"
cp "$config" "$run_dir/config/effective.json"

unset COSMOS_PERF_CPU_PROFILE COSMOS_PERF_HEAP_PROFILE COSMOS_PERF_RUNTIME_TRACE \
  COSMOS_PERF_BLOCK_PROFILE COSMOS_PERF_MUTEX_PROFILE
export COSMOS_PERF_CONFIG=$run_dir/config/effective.json

binary=$(readlink -f "$1")
sha256sum "$binary" > "$raw/binary.sha256"
readelf -n "$binary" > "$raw/binary-build-id.txt" 2>&1 || true
cp --preserve=mode,timestamps "$binary" "$raw/binary"

"$@" > "$raw/stdout.txt" 2> "$raw/stderr.txt" &
pid=$!
echo "$pid" > "$raw/pid.txt"

sample() {
  local name=$1 cgroup_rel cgroup_dir
  cgroup_rel=$(awk -F: '$1 == "0" {print $3}' "/proc/$pid/cgroup")
  cgroup_dir=/sys/fs/cgroup$cgroup_rel
  cat "/proc/$pid/status" > "$raw/proc-status-$name.txt"
  cat "/proc/$pid/smaps_rollup" > "$raw/smaps-rollup-$name.txt"
  cat "/proc/$pid/smaps" > "$raw/smaps-$name.txt"
  {
    echo "cgroup=$cgroup_rel"
    echo "memory.current=$(cat "$cgroup_dir/memory.current" 2>/dev/null || true)"
    echo "memory.peak=$(cat "$cgroup_dir/memory.peak" 2>/dev/null || true)"
    cat "$cgroup_dir/memory.stat" 2>/dev/null || true
  } > "$raw/cgroup-memory-$name.txt"
  for thread in /proc/"$pid"/task/*; do
    echo "=== ${thread##*/} ==="
    cat "$thread/stack" 2>/dev/null || true
  done > "$raw/thread-stacks-$name.txt"
}

sleep 1
sample start
while kill -0 "$pid" 2>/dev/null; do sample latest; sleep 5; done
set +e
wait "$pid"
status=$?
set -e
sample end 2>/dev/null || true
for base in proc-status smaps-rollup smaps cgroup-memory thread-stacks; do
  cp "$raw/$base-latest.txt" "$raw/$base.txt"
done
sha256sum "$raw"/* > "$raw/SHA256SUMS"
chmod -R a-w "$raw"
exit "$status"
