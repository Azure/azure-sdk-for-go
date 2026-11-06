#!/usr/bin/env bash
set -Eeuo pipefail
umask 077
phase=${1:-all}
[[ "$phase" == all || "$phase" == clean || "$phase" == diagnostics || "$phase" == ffi ]]
gateway=/mnt/cosmos-perf/tokio234-gatewayv1-20261105-a
source "$gateway/tools/common.sh"
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
inputs=$(cd "$(dirname "$0")" && pwd)
region='West US 2'
if [[ -f "$root/private/owner.txt" ]]; then
  owner=$(cat "$root/private/owner.txt")
else
  owner=$(python3 -c 'import uuid; print("ffi-flow-20261105-a-" + uuid.uuid4().hex[:8])')
fi
[[ -f "$root/provenance/GO_BUILD_COMPLETE" && ! -e "$root/provenance/CAMPAIGN_COMPLETE" ]]
[[ "$(nproc)" == 4 ]]
mkdir -p "$root/run" "$root/wire-proof" "$root/tools/perf-home"
printf '%s\n' "$owner" > "$root/private/owner.txt"
cp "$inputs/run-spec.json" "$root/run-spec.json"
python3 - "$root" "$inputs" <<'PY'
import os, pathlib, sys
sys.path.insert(0, sys.argv[2])
from run_leg import memory_root, memory_value
active = []
for path in pathlib.Path("/proc").iterdir():
    if not path.name.isdigit():
        continue
    try:
        binary = os.readlink(path / "exe")
    except (FileNotFoundError, PermissionError, ProcessLookupError):
        continue
    if binary.startswith("/mnt/cosmos") and pathlib.Path(binary).name in {
        "go-v2-native", "rust-bench", "control", "option-b", "c0-control", "cosmos_perf_flow",
    }:
        active.append(int(path.name))
assert not active, f"another benchmark is active: {active}"
group = memory_root() / f"cosmos-flow-prerequisite-{os.getpid()}"
group.mkdir()
try:
    assert memory_value(group) >= 0
    assert memory_value(group, peak=True) >= 0
finally:
    group.rmdir()
print("DEDICATED_MEMORY_ACCOUNTING_READY")
PY

libc=/usr/lib/x86_64-linux-gnu/libc.so.6
libc_id=$(readelf -n "$libc" | awk '/Build ID:/ {print $3}')
debug=/mnt/cosmos-perf/cpu-isolation-synthetic-20261105-d/attribution-symbols/libc6-dbg/extracted/usr/lib/debug/.build-id/${libc_id:0:2}/${libc_id:2}.debug
[[ -s "$debug" && "$(readelf -n "$debug" | awk '/Build ID:/ {print $3}')" == "$libc_id" ]]
if [[ ! -e "$root/provenance/libc-symbols.sha256" ]]; then
  HOME="$root/tools/perf-home" perf buildid-cache --add "$debug" \
    > "$root/provenance/perf-symbol-cache.log" 2>&1
  sha256sum "$libc" "$debug" > "$root/provenance/libc-symbols.sha256"
else
  sha256sum --check "$root/provenance/libc-symbols.sha256" > /dev/null
fi
date -u +%FT%TZ > "$root/provenance/RUN_STARTED"

common_args=(--endpoint "$endpoint" --database "$database" --container "$container" --owner "$owner")
control=$root/build/bin/clean/go-v2-native
seeded=false
cleanup_owned() {
  if [[ "$seeded" == true ]]; then
    enable_gateway_v1
    set +e
    with_key env -u COSMOS_PERF_WIRE_PROOF -u COSMOS_PERF_FLOW_DIR \
      "$control" "${common_args[@]}" --key-fd 3 cleanup \
      > "$root/provenance/cleanup.json" 2> "$root/private/cleanup.stderr"
    status=$?
    set -e
    [[ "$status" == 0 ]] || { echo "Owner-scoped cleanup failed" >&2; return "$status"; }
    seeded=false
  fi
}
trap cleanup_owned EXIT
enable_gateway_v1
seeded=true
needs_seed=$(python3 - "$root" <<'PY'
import json, pathlib, sys
root = pathlib.Path(sys.argv[1])
seed = root / "provenance/seed.json"
cleanup = root / "provenance/cleanup.json"
deleted = 0
if cleanup.exists() and cleanup.stat().st_size:
    deleted = json.loads(cleanup.read_text()).get("deleted", 0)
print("true" if not seed.exists() or deleted > 0 else "false")
PY
)
if [[ "$needs_seed" == true ]]; then
  if [[ -f "$root/provenance/seed.json" ]]; then
    cp "$root/provenance/seed.json" "$root/provenance/seed-before-$(date -u +%Y%m%dT%H%M%SZ).json"
  fi
  with_key "$control" "${common_args[@]}" --key-fd 3 seed \
    > "$root/provenance/seed.json" 2> "$root/private/seed.stderr"
fi
with_key "$control" "${common_args[@]}" --key-fd 3 verify \
  > "$root/provenance/verify-$(date -u +%Y%m%dT%H%M%SZ).json" 2> "$root/private/verify.stderr"

reusable_wire=$(python3 - "$root" <<'PY'
import json, pathlib, sys
root = pathlib.Path(sys.argv[1])
path = root / "provenance/wire-proof.json"
valid = path.exists()
if valid:
    proof = json.loads(path.read_text())
    for leg in ("go-v2-native", "rust-sdk", "rust-driver"):
        name = "go-v2-native" if leg == "go-v2-native" else "rust-bench"
        current = json.loads((root / f"build/bin/proof/{name}.provenance.json").read_text())
        prior = proof["legs"][leg]
        valid = valid and prior["proofBinarySHA256"] == current["binary"]["sha256"]
        valid = valid and prior["buildClassSourceHash"] == current["wireSourceHash"]
print("true" if valid else "false")
PY
)
proof_id=$(sha256sum "$root/build/bin/proof/rust-bench.provenance.json" | cut -c 1-12)
proof_root=$root/wire-proof/$proof_id
if [[ "$reusable_wire" != true ]]; then
if [[ -f "$root/provenance/wire-proof.json" ]]; then
  cp "$root/provenance/wire-proof.json" "$root/provenance/wire-proof-before-$proof_id.json"
fi
for leg in go-v2-native rust-sdk rust-driver; do
  directory=$proof_root/$leg
  mkdir -p "$directory"
  if [[ "$leg" == go-v2-native ]]; then
    binary=$root/build/bin/proof/go-v2-native
    extra=()
  else
    binary=$root/build/bin/proof/rust-bench
    extra=(--leg "$leg" --region "$region")
  fi
  enable_gateway_v1
  with_key env COSMOS_PERF_WIRE_PROOF="$directory/representation.json" \
    "$binary" "${common_args[@]}" "${extra[@]}" --key-fd 3 preflight \
    > "$directory/caller.json" 2> "$directory/stderr.txt"
done
python3 - "$root" "$proof_root" <<'PY'
import hashlib, json, pathlib, sys
root = pathlib.Path(sys.argv[1])
proof_root = pathlib.Path(sys.argv[2])
output = {
    "rustSourceCommit": "77ebe786ad1012bc9ac56105e08c05fb4c65b321",
    "goOptimizedBaseCommit": "b20b0fd00ac20c78bdb2d138fb6960b35e428af0",
    "nativeVersion": "0.3.0", "legs": {},
}
for leg in ("go-v2-native", "rust-sdk", "rust-driver"):
    directory = proof_root / leg
    value = json.loads((directory / "representation.json").read_text())
    count = value["eligibleResponses"]
    assert count > 0 and value["binaryResponsesBeforeConversion"] == 0
    assert value["textResponsesBeforeConversion"] == count
    assert value["invalidResponsesBeforeConversion"] == 0
    assert value["gatewayV1HTTP2Responses"] == count
    assert value["textOutputRequests"] == 0
    name = "go-v2-native" if leg == "go-v2-native" else "rust-bench"
    provenance = json.loads((root / "build/bin/proof" / f"{name}.provenance.json").read_text())
    output["legs"][leg] = dict(value, binaryDisabledExplicitly=True,
                              buildClassSourceHash=provenance["wireSourceHash"],
                              proofBinarySHA256=provenance["binary"]["sha256"],
                              evidenceDirectory=directory.relative_to(root).as_posix())
    with (directory / "SHA256SUMS").open("x") as handle:
        for path in sorted(directory.iterdir()):
            if path.name != "SHA256SUMS":
                handle.write(f"{hashlib.sha256(path.read_bytes()).hexdigest()}  {path.name}\n")
(root / "provenance/wire-proof.json").write_text(json.dumps(output, indent=2) + "\n")
print("ALL_THREE_BINARY_DISABLED_GATEWAY_V1_HTTP2_PROOFS_PASS")
PY
else
  python3 - "$root" <<'PY'
import hashlib, json, pathlib, sys
root = pathlib.Path(sys.argv[1])
value = json.loads((root / "provenance/wire-proof.json").read_text())
for leg, proof in value["legs"].items():
    assert proof["binaryResponsesBeforeConversion"] == 0
    assert proof["textResponsesBeforeConversion"] == proof["eligibleResponses"] > 0
    assert proof["gatewayV1HTTP2Responses"] == proof["eligibleResponses"]
    directory = root / proof.get("evidenceDirectory", "wire-proof/" + leg)
    for row in (directory / "SHA256SUMS").read_text().splitlines():
        digest, name = row.split("  ", 1)
        assert hashlib.sha256((directory / name).read_bytes()).hexdigest() == digest
print("REUSING_HASH_VERIFIED_THREE_LEG_WIRE_PROOF")
PY
fi

run_leg() {
  local workers=$1 class=$2 repeat=$3 leg=$4 warmup=$5 measure=$6 mode=${7:-}
  local build_class=clean
  [[ "$class" == ffi-flow || "$class" == rust-flow ]] && build_class=flow
  local binary=$root/build/bin/$build_class/go-v2-native
  local extra=() calibration=()
  if [[ "$leg" != go-v2-native ]]; then
    binary=$root/build/bin/$build_class/rust-bench
    extra=(--leg "$leg" --region "$region")
  fi
  [[ -n "$mode" ]] && calibration=(--calibration-mode "$mode")
  enable_gateway_v1
  printf 'RUN class=%s workers=%s repeat=%s leg=%s warmup=%s measurement=%s\n' \
    "$class" "$workers" "$repeat" "$leg" "$warmup" "$measure"
  with_key python3 "$inputs/run_leg.py" --root "$root" --leg "$leg" \
    --workers "$workers" --class "$class" --repeat "$repeat" --key-fd 3 "${calibration[@]}" \
    -- "$binary" "${common_args[@]}" "${extra[@]}" --key-fd 3 \
    --workers "$workers" --warmup-seconds "$warmup" --measurement-seconds "$measure" run
}

for leg in go-v2-native rust-sdk rust-driver; do
  run_leg 32 perf-calibration 1 "$leg" 2 5 fp
  run_leg 32 perf-calibration 2 "$leg" 2 5 dwarf
done
python3 "$inputs/calibrate.py" "$root"
for workers in 32 64; do
  if [[ "$workers" == 32 ]]; then
    first=(go-v2-native rust-sdk rust-driver)
    second=(rust-driver rust-sdk go-v2-native)
  else
    first=(rust-driver rust-sdk go-v2-native)
    second=(go-v2-native rust-sdk rust-driver)
  fi
  if [[ "$phase" != diagnostics && "$phase" != ffi ]]; then
    for leg in "${first[@]}"; do run_leg "$workers" clean-score 1 "$leg" 10 60; done
    for leg in "${second[@]}"; do run_leg "$workers" clean-score 2 "$leg" 10 60; done
  fi
done
if [[ "$phase" != diagnostics && "$phase" != ffi ]]; then
  python3 "$inputs/clean_summary.py" "$root"
  date -u +%FT%TZ > "$root/provenance/CLEAN_SCORE_COMPLETE"
fi
if [[ "$phase" == clean ]]; then
  cleanup_owned
  trap - EXIT
  echo MATCHED_CLEAN_W32_W64_CHECKPOINT_COMPLETE
  exit 0
fi
if [[ "$phase" == ffi ]]; then
  for workers in 32 64; do run_leg "$workers" ffi-flow 1 go-v2-native 10 60; done
  python3 "$inputs/flow_summary.py" "$root"
  cleanup_owned
  trap - EXIT
  echo MATCHED_FFI_W32_W64_CHECKPOINT_COMPLETE
  exit 0
fi
for workers in 32 64; do
  first=(go-v2-native rust-sdk rust-driver)
  for leg in "${first[@]}"; do run_leg "$workers" perf 1 "$leg" 10 60; done
  run_leg "$workers" ffi-flow 1 go-v2-native 10 60
  run_leg "$workers" rust-flow 1 rust-sdk 10 60
  run_leg "$workers" rust-flow 1 rust-driver 10 60
  run_leg "$workers" go-diagnostic 1 go-v2-native 2 15
  run_leg "$workers" scheduler 1 go-v2-native 2 10
done
cleanup_owned
trap - EXIT
python3 - "$root/provenance/cleanup.json" <<'PY'
import json, pathlib, sys
value = json.loads(pathlib.Path(sys.argv[1]).read_text())
assert value["deleted"] == 1024 and value["nonOwnedSkipped"] == 0
PY
python3 "$inputs/analyze_go_diagnostics.py" "$root"
python3 "$inputs/profile_attribution.py" "$root"
python3 "$inputs/report.py" "$root"
date -u +%FT%TZ > "$root/provenance/CAMPAIGN_COMPLETE"
echo MATCHED_W32_W64_THREE_LEG_CAMPAIGN_COMPLETE
