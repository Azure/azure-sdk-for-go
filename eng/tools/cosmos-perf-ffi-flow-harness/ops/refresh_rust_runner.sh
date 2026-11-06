#!/usr/bin/env bash
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
source=$root/src/rust
inputs=$(cd "$(dirname "$0")" && pwd)
tools=/mnt/cosmos-perf/point-read-first-pass-20261102-a/tools
[[ -f "$root/provenance/GO_BUILD_COMPLETE" && ! -e "$root/provenance/CAMPAIGN_COMPLETE" ]]
cp "$inputs/rust_main.rs" "$source/sdk/cosmos/cosmos_perf_flow/src/main.rs"
export PATH="$tools/rustup/toolchains/1.95-x86_64-unknown-linux-gnu/bin:$tools/cargo/bin:$PATH"
export RUSTUP_HOME="$tools/rustup" CARGO_HOME="$tools/cargo"
export CARGO_TARGET_DIR="$root/build/cargo-target" CARGO_INCREMENTAL=0 CARGO_BUILD_JOBS=4
export RUSTFLAGS="-C force-frame-pointers=yes -C debuginfo=2 -C link-arg=-Wl,--build-id=sha1"
export CARGO_PROFILE_RELEASE_DEBUG=2 CARGO_PROFILE_RELEASE_STRIP=false
export TMPDIR="$root/build/tmp"
cd "$source"
cargo fmt -p cosmos_perf_flow
for class in clean proof flow; do
  features=()
  [[ "$class" == proof ]] && features=(--features wire-proof)
  [[ "$class" == flow ]] && features=(--features flow-diagnostic)
  cargo build --locked --release -p cosmos_perf_flow "${features[@]}" \
    > "$root/build/rust/$class/runner-refresh.log" 2>&1 || {
      tail -n 45 "$root/build/rust/$class/runner-refresh.log"; exit 1;
    }
  cp "$CARGO_TARGET_DIR/release/cosmos_perf_flow" "$root/build/bin/$class/rust-bench.next"
  chmod 0555 "$root/build/bin/$class/rust-bench.next"
  mv "$root/build/bin/$class/rust-bench.next" "$root/build/bin/$class/rust-bench"
done
python3 "$inputs/write_build_provenance.py" "$root"
echo RUST_WARMUP_ACCOUNTING_REFRESH_COMPLETE
