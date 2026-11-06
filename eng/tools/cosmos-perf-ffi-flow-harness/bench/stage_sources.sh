#!/usr/bin/env bash
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
go_repo=/mnt/cosmos-perf/container-cache-cbatch-w16-20261103-a/src/go-v2-combined
rust_repo=/mnt/cosmos-perf/point-read-first-pass-20261102-i/src/rust-sdk
adapter=/mnt/cosmos-perf/cpu-isolation-synthetic-20261105-d/build/work/c0-control
go_commit=b20b0fd00ac20c78bdb2d138fb6960b35e428af0
rust_commit=77ebe786ad1012bc9ac56105e08c05fb4c65b321
[[ -f "$root/tools/INPUTS_COMPLETE" && ! -e "$root/provenance/SOURCES_STAGED" ]]
[[ "$(git -C "$go_repo" rev-parse HEAD)" == "$go_commit" ]]
[[ "$(git -C "$rust_repo" cat-file -t "$rust_commit")" == commit ]]
[[ ! -e "$root/src/go/clean" && ! -e "$root/src/rust" ]]
mkdir -p "$root/src/go/clean" "$root/src/adapter-base"
git -C "$go_repo" archive "$go_commit" sdk/data/azcosmos |
  tar -x -C "$root/src/go/clean" --strip-components=3
git -C "$rust_repo" worktree add --detach "$root/src/rust" "$rust_commit" \
  > "$root/provenance/rust-checkout.log" 2>&1
cp "$adapter/"*.go "$adapter/go.mod" "$adapter/go.sum" "$root/src/adapter-base/"
python3 - "$root" <<'PY'
import base64, hashlib, json, lzma, pathlib, shutil, sys
root = pathlib.Path(sys.argv[1])
sdk = root / "src/go/clean"
expected = {
    "driver_native.go": "ca2b996dd9d880202c42096cc17880364572dc3e85182f18aaaafc2fafeb6d1e",
    "reactor_native.go": "30aa373af3d9c7dcd700aede745e0cbffa209be65a1abc7f5a92b58e3c60e2d2",
    "operation_native.go": "7d06e016b09d823ad2a390ae58f38df2ecfa2eb22d3a2b2e5a8e414a400b69a3",
}
for name, digest in expected.items():
    actual = hashlib.sha256((sdk / name).read_bytes()).hexdigest()
    assert actual == digest, f"optimized source mismatch: {name}"
manifest = {p.relative_to(sdk).as_posix(): hashlib.sha256(p.read_bytes()).hexdigest()
            for p in sorted(sdk.rglob("*")) if p.is_file()}
(root / "provenance/go-optimized-source-sha256.json").write_text(json.dumps(manifest, indent=2) + "\n")
context = {}
for name in ("driver_native.go", "reactor_native.go", "operation_native.go", "driver.go", "go.mod"):
    context["sdk/" + name] = (sdk / name).read_text()
for name in ("main.go", "sdk_v2.go", "go.mod"):
    context["adapter/" + name] = (root / "src/adapter-base" / name).read_text()
packet = lzma.compress(json.dumps(context, sort_keys=True).encode(), preset=6)
encoded = base64.b64encode(packet).decode()
(root / "tools/go-source-context.json.xz").write_bytes(packet)
(root / "tools/go-source-context.b64").write_text(encoded)
(root / "provenance/SOURCES_STAGED").write_text("exact tracked Go snapshot and pinned Rust worktree\n")
print("CONTEXT_META:" + json.dumps({
    "base64Bytes": len(encoded), "sha256": hashlib.sha256(packet).hexdigest(),
    "sdkFileCount": len(manifest), "freeGiB": round(shutil.disk_usage(root).free / 2**30, 2),
}, separators=(",", ":")))
print("DATA:" + encoded[:3200])
print("FLOW_SOURCE_STAGE_COMPLETE")
PY
