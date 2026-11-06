#!/usr/bin/env bash
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
source=$root/src/rust
inputs=$(cd "$(dirname "$0")" && pwd)
tools=/mnt/cosmos-perf/point-read-first-pass-20261102-a/tools
sysroot=/mnt/cosmos-perf/session-token-stream-w64-20261105-a/tools/build-prerequisites/sysroot
[[ -f "$root/provenance/SOURCES_STAGED" && ! -e "$root/provenance/RUST_BUILD_COMPLETE" ]]
[[ "$(git -C "$source" rev-parse HEAD)" == 77ebe786ad1012bc9ac56105e08c05fb4c65b321 ]]
mkdir -p "$root/build/rust" "$root/build/cargo-target"
export PATH="$tools/rustup/toolchains/1.95-x86_64-unknown-linux-gnu/bin:$tools/cargo/bin:$PATH"
export RUSTUP_HOME="$tools/rustup" CARGO_HOME="$tools/cargo"
export CARGO_TARGET_DIR="$root/build/cargo-target"
export CARGO_INCREMENTAL=0 CARGO_BUILD_JOBS=4
export RUSTFLAGS="-C force-frame-pointers=yes -C debuginfo=2 -C link-arg=-Wl,--build-id=sha1"
export CARGO_PROFILE_RELEASE_DEBUG=2 CARGO_PROFILE_RELEASE_STRIP=false
export PKG_CONFIG="$sysroot/usr/bin/pkg-config" PKG_CONFIG_SYSROOT_DIR="$sysroot"
export PKG_CONFIG_LIBDIR="$sysroot/usr/lib/x86_64-linux-gnu/pkgconfig"
export OPENSSL_LIB_DIR="$sysroot/usr/lib/x86_64-linux-gnu" OPENSSL_INCLUDE_DIR="$sysroot/usr/include"
export LD_LIBRARY_PATH="$sysroot/usr/lib/x86_64-linux-gnu"
export CFLAGS="-I$sysroot/usr/include/x86_64-linux-gnu"

if [[ ! -e "$root/provenance/RUST_SOURCE_PREPARED" ]]; then
  python3 "$inputs/prepare_rust.py" "$source" > "$root/build/rust/prepare.log" 2>&1
  date -u +%FT%TZ > "$root/provenance/RUST_SOURCE_PREPARED"
fi

cd "$source"
cargo fmt -p azure_data_cosmos_driver -p azure_data_cosmos_driver_native -p cosmos_perf_flow \
  > "$root/build/rust/format.log" 2>&1
cargo fmt -p azure_data_cosmos_driver -p azure_data_cosmos_driver_native -p cosmos_perf_flow -- --check \
  >> "$root/build/rust/format.log" 2>&1

for class in clean proof flow; do
  output=$root/build/rust/$class
  mkdir -p "$output"
  features=()
  locked=(--locked)
  if [[ "$class" == clean ]]; then
    locked=()
  elif [[ "$class" == proof ]]; then
    features=(--features azure_data_cosmos_driver_native/__internal_perf_wire,cosmos_perf_flow/wire-proof)
  else
    features=(--features azure_data_cosmos_driver_native/__internal_perf_flow,cosmos_perf_flow/flow-diagnostic)
  fi
  if ! cargo build "${locked[@]}" --release -p azure_data_cosmos_driver_native -p cosmos_perf_flow \
      "${features[@]}" > "$output/build.log" 2>&1; then
    echo "RUST_BUILD_FAILED class=$class"
    tail -n 55 "$output/build.log"
    exit 1
  fi
  cp "$CARGO_TARGET_DIR/release/libazurecosmosdriver.a" "$output/"
  cp "$CARGO_TARGET_DIR/release/cosmos_perf_flow" "$output/"
  cp "$source/sdk/cosmos/azure_data_cosmos_driver_native/include/azurecosmosdriver.h" "$output/"
  sha256sum "$output/libazurecosmosdriver.a" "$output/cosmos_perf_flow" "$output/azurecosmosdriver.h" \
    > "$output/SHA256SUMS"
  readelf -n "$output/cosmos_perf_flow" > "$output/elf-notes.txt"
done

cargo test --locked -p azure_data_cosmos_driver_native --lib --features __internal_perf_flow \
  perf_flow::tests:: -- --test-threads=1 > "$root/build/rust/flow-tests.log" 2>&1
git diff --binary > "$root/provenance/rust-source.patch"
git ls-files --others --exclude-standard > "$root/provenance/rust-added-source-files.txt"
cp Cargo.lock "$root/provenance/Cargo.lock"
rustc -Vv > "$root/provenance/rustc-version.txt"
cargo -V > "$root/provenance/cargo-version.txt"
printf '%s\n' "$RUSTFLAGS" > "$root/provenance/rustflags.txt"
date -u +%FT%TZ > "$root/provenance/RUST_BUILD_COMPLETE"
echo "RUST_THREE_CLASS_BUILD_COMPLETE"
