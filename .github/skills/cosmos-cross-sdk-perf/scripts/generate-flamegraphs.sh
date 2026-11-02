#!/usr/bin/env bash
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words Eeuo difffolded flamegraph rsplit

set -Eeuo pipefail

[[ $# -eq 3 ]] || {
  echo "usage: generate-flamegraphs.sh <run-dir> <FlameGraph-checkout> <expected-commit>" >&2
  exit 2
}
run_dir=$(readlink -f "$1")
flamegraph=$(readlink -f "$2")
expected_commit=$3
[[ "$expected_commit" =~ ^[0-9a-f]{40}$ ]]
actual_commit=$(git -C "$flamegraph" rev-parse HEAD)
[[ "$actual_commit" == "$expected_commit" ]] || {
  echo "FlameGraph checkout is at $actual_commit, expected $expected_commit" >&2
  exit 1
}
[[ -z "$(git -C "$flamegraph" status --porcelain)" ]] || {
  echo "FlameGraph checkout must be clean" >&2
  exit 1
}
collapse=$flamegraph/stackcollapse-perf.pl
render=$flamegraph/flamegraph.pl
diff=$flamegraph/difffolded.pl
[[ -x "$collapse" && -x "$render" && -x "$diff" ]]
printf '%s\n' "$actual_commit" > "$run_dir/FLAMEGRAPH-COMMIT"

wrap_html() {
  local svg=$1 html=${1%.svg}.html title=$2
  {
    echo '<!doctype html><html><head><meta charset="utf-8">'
    echo "<title>$title</title></head><body><h1>$title</h1>"
    cat "$svg"
    echo '</body></html>'
  } > "$html"
}

scale_folded() {
  local input=$1 output=$2 divisor=$3 normalize=$4
  python3 - "$input" "$output" "$divisor" "$normalize" <<'PY'
import pathlib
import sys

source, target = map(pathlib.Path, sys.argv[1:3])
divisor = float(sys.argv[3])
normalize = sys.argv[4] == "true"
rows = []
for line in source.read_text(errors="replace").splitlines():
    stack, value = line.rsplit(" ", 1)
    rows.append((stack, float(value)))
total = sum(value for _, value in rows) or 1.0
with target.open("w") as handle:
    for stack, value in rows:
        scaled = value / divisor
        if normalize:
            scaled = value / total * 1_000_000.0
        handle.write(f"{stack} {scaled:.9f}\n")
PY
}

for workers in 8 16; do
  for leg in go-v1 go-v2-native rust; do
    dir=$run_dir/w$workers/perf/repeat-1/$leg
    raw=$dir/raw
    derived=$dir/derived
    mkdir -p "$derived"
    "$collapse" --all "$raw/perf-script.txt" > "$derived/whole-process.folded"
    "$render" --hash --colors hot --title "w$workers $leg whole-process CPU" \
      "$derived/whole-process.folded" > "$derived/whole-process.svg"
    successes=$(python3 -c "import json; print(json.load(open('$dir/metrics.json'))['counts']['success'])")
    scale_folded "$derived/whole-process.folded" "$derived/cpu-per-read.folded" \
      "$successes" false
    "$render" --hash --colors hot --title "w$workers $leg CPU per successful read" \
      --countname "samples/read" "$derived/cpu-per-read.folded" \
      > "$derived/cpu-per-read.svg"
    scale_folded "$derived/whole-process.folded" "$derived/composition.folded" 1 true
    wrap_html "$derived/whole-process.svg" "w$workers $leg whole-process CPU"
    wrap_html "$derived/cpu-per-read.svg" "w$workers $leg CPU per successful read"
  done

  for base in go-v1 rust; do
    pair=$run_dir/w$workers/perf/differential/$base-vs-go-v2-native
    mkdir -p "$pair"
    base_dir=$run_dir/w$workers/perf/repeat-1/$base/derived
    v2_dir=$run_dir/w$workers/perf/repeat-1/go-v2-native/derived
    "$diff" "$base_dir/cpu-per-read.folded" "$v2_dir/cpu-per-read.folded" \
      > "$pair/forward-cpu-per-read.folded"
    "$diff" "$v2_dir/cpu-per-read.folded" "$base_dir/cpu-per-read.folded" \
      > "$pair/reverse-cpu-per-read.folded"
    "$diff" "$base_dir/whole-process.folded" "$v2_dir/whole-process.folded" \
      > "$pair/forward-whole-process.folded"
    "$diff" "$v2_dir/whole-process.folded" "$base_dir/whole-process.folded" \
      > "$pair/reverse-whole-process.folded"
    "$diff" "$base_dir/composition.folded" "$v2_dir/composition.folded" \
      > "$pair/forward-composition-normalized.folded"
    "$diff" "$v2_dir/composition.folded" "$base_dir/composition.folded" \
      > "$pair/reverse-composition-normalized.folded"
    for direction in forward-whole-process reverse-whole-process \
      forward-cpu-per-read reverse-cpu-per-read \
      forward-composition-normalized reverse-composition-normalized; do
      "$render" --hash --colors hot --negate \
        --title "w$workers $base vs go-v2-native $direction" \
        "$pair/$direction.folded" > "$pair/$direction.svg"
      wrap_html "$pair/$direction.svg" "w$workers $base vs go-v2-native $direction"
    done
  done
done
