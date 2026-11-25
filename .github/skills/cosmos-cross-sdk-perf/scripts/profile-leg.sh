#!/usr/bin/env bash
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words Eeuo esac gsub rollup smaps

set -Eeuo pipefail

usage() {
  echo "usage: profile-leg.sh --leg <go-v1|go-v2-native|rust> --workers <8|16> --repeat <n> --run-dir <dir> --config <json> --key-fd <n> -- <command...>" >&2
  exit 2
}

leg=
run_dir=
config=
workers=
repeat=
key_fd=
while [[ $# -gt 0 ]]; do
  case "$1" in
    --leg) leg=$2; shift 2 ;;
    --workers) workers=$2; shift 2 ;;
    --repeat) repeat=$2; shift 2 ;;
    --run-dir) run_dir=$2; shift 2 ;;
    --config) config=$2; shift 2 ;;
    --key-fd) key_fd=$2; shift 2 ;;
    --) shift; break ;;
    *) usage ;;
  esac
done
[[ "$leg" =~ ^(go-v1|go-v2-native|rust)$ ]] || usage
[[ "$workers" =~ ^(8|16)$ && "$repeat" =~ ^[1-9][0-9]*$ ]] || usage
[[ "$key_fd" =~ ^[3-9][0-9]*$ && -r "/proc/self/fd/$key_fd" ]] || usage
[[ -n "$run_dir" && -f "$config" && $# -gt 0 ]] || usage

leg_dir=$run_dir/w$workers/perf/repeat-$repeat/$leg
raw=$leg_dir/raw
derived=$leg_dir/derived
[[ ! -e "$leg_dir" ]] || { echo "refusing to overwrite $leg_dir" >&2; exit 1; }
mkdir -p "$raw" "$derived" "$run_dir/config"
cp "$config" "$run_dir/config/effective.json"

command_path=$(readlink -f "$1")
[[ -x "$command_path" ]]
sha256sum "$command_path" > "$raw/binary.sha256"
readelf -n "$command_path" > "$raw/binary-build-id.txt" 2>&1 || true
cp --preserve=mode,timestamps "$command_path" "$raw/binary"
[[ -f "$command_path.provenance.json" ]] || { echo "missing $command_path.provenance.json" >&2; exit 1; }
cp "$command_path.provenance.json" "$leg_dir/provenance.json"

export COSMOS_PERF_CONFIG=$run_dir/config/effective.json
export COSMOS_PERF_WORKERS=$workers
export COSMOS_PERF_ADAPTER_RESULT=$raw/adapter-result.json
export COSMOS_PERF_KEY_FD=$key_fd
unset COSMOS_PERF_CPU_PROFILE COSMOS_PERF_HEAP_PROFILE COSMOS_PERF_RUNTIME_TRACE \
  COSMOS_PERF_BLOCK_PROFILE COSMOS_PERF_MUTEX_PROFILE

call_graph=$(python3 - "$config" <<'PY'
import json
import sys
value = json.load(open(sys.argv[1], encoding="utf-8"))["profiling"]["selectedCallGraph"]
if not value:
    raise SystemExit("profiling.selectedCallGraph must be chosen by calibration before perf runs")
print(value)
PY
)
case "$call_graph" in
  frame-pointer) perf_call_graph=fp ;;
  dwarf,16384) perf_call_graph=dwarf,16384 ;;
  *) echo "unsupported call graph: $call_graph" >&2; exit 2 ;;
esac

started=$(date -u +%FT%TZ)
printf '%s\n' "$started" > "$raw/started-utc.txt"

"$@" > "$raw/stdout.txt" 2> "$raw/stderr.txt" &
pid=$!
printf '%s\n' "$pid" > "$raw/pid.txt"
perf record -F 49 -g --call-graph "$perf_call_graph" -p "$pid" \
  -o "$raw/perf.data" 2> "$raw/perf-record.stderr" &
perf_pid=$!

sample_proc() {
  local target=$1
  local cgroup_rel cgroup_dir
  cgroup_rel=$(awk -F: '$1 == "0" {print $3}' "/proc/$pid/cgroup")
  cgroup_dir=/sys/fs/cgroup$cgroup_rel
  cat "/proc/$pid/status" > "$raw/proc-status-$target.txt"
  cat "/proc/$pid/smaps_rollup" > "$raw/smaps-rollup-$target.txt"
  cat "/proc/$pid/smaps" > "$raw/smaps-$target.txt"
  {
    echo "cgroup=$cgroup_rel"
    echo "memory.current=$(cat "$cgroup_dir/memory.current" 2>/dev/null || true)"
    echo "memory.peak=$(cat "$cgroup_dir/memory.peak" 2>/dev/null || true)"
    echo "memory.stat:"
    cat "$cgroup_dir/memory.stat" 2>/dev/null || true
  } > "$raw/cgroup-memory-$target.txt"
  {
    for thread in /proc/"$pid"/task/*; do
      echo "=== ${thread##*/} ==="
      cat "$thread/stack" 2>/dev/null || true
    done
  } > "$raw/thread-stacks-$target.txt"
}

sleep 1
sample_proc start
while kill -0 "$pid" 2>/dev/null; do
  sample_proc latest
  sleep 5
done
set +e
wait "$pid"
status=$?
kill -INT "$perf_pid" 2>/dev/null || true
wait "$perf_pid"
perf_status=$?
set -e
sample_proc end 2>/dev/null || true

for base in proc-status smaps-rollup smaps cgroup-memory thread-stacks; do
  cp "$raw/$base-latest.txt" "$raw/$base.txt"
done
if [[ "$status" -eq 0 ]]; then
  [[ -s "$raw/adapter-result.json" ]] || { echo "adapter did not write its result" >&2; exit 1; }
  python3 "$(dirname "$0")/convert-adapter-output.py" \
    --adapter-result "$raw/adapter-result.json" --raw-dir "$raw" \
    --leg "$leg" --run-class perf --workers "$workers" --repeat "$repeat" \
    --output "$leg_dir/metrics.json"
fi
perf report --no-inline --stdio --show-nr-samples --percent-limit 0 \
  -i "$raw/perf.data" > "$raw/perf-report.txt"
perf script --no-inline -i "$raw/perf.data" > "$raw/perf-script.txt"
python3 "$(dirname "$0")/perf-quality.py" \
  --perf-script "$raw/perf-script.txt" \
  --perf-stderr "$raw/perf-record.stderr" \
  --mode "$call_graph" \
  --output "$raw/perf-quality.json"
{
  echo "symbol,dso,self_percent"
  perf report --no-inline --stdio --field-separator '|' \
    -F overhead,symbol,dso --percent-limit 0 \
    -i "$raw/perf.data" |
    awk -F'|' '
      /^[#]/ || NF < 3 { next }
      {
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", $1)
        gsub(/%/, "", $1)
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", $2)
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", $3)
        if ($1 ~ /^[0-9.]+$/) {
          gsub(/"/, "\"\"", $2)
          gsub(/"/, "\"\"", $3)
          print "\"" $2 "\",\"" $3 "\"," $1
        }
      }'
} > "$leg_dir/perf-self.csv"
python3 "$(dirname "$0")/normalize-perf-self.py" "$leg_dir/perf-self.csv"
sha256sum "$raw"/* > "$raw/SHA256SUMS"
chmod -R a-w "$raw"
printf '%s\n' "$(date -u +%FT%TZ)" > "$leg_dir/completed-utc.txt"
if [[ "$perf_status" -ne 0 ]]; then
  echo "perf collection failed with status $perf_status" >&2
  exit "$perf_status"
fi
exit "$status"
