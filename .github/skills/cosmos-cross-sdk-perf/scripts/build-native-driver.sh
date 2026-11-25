#!/usr/bin/env bash
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words Eeuo rustc

set -Eeuo pipefail

usage() {
  echo "usage: build-native-driver.sh <azure-sdk-for-rust-checkout> <commit> <output-dir>" >&2
  exit 2
}

[[ $# -eq 3 ]] || usage
checkout=$(readlink -f "$1")
commit=$2
output=$(readlink -m "$3")
[[ "$commit" =~ ^[0-9a-f]{40}$ ]] || usage
[[ -d "$checkout/.git" || -f "$checkout/.git" ]]
[[ ! -e "$output" ]] || { echo "refusing to overwrite $output" >&2; exit 1; }

actual_commit=$(git -C "$checkout" rev-parse HEAD)
[[ "$actual_commit" == "$commit" ]] || {
  echo "checkout is at $actual_commit, expected $commit" >&2
  exit 1
}
[[ -z "$(git -C "$checkout" status --short)" ]] || {
  echo "checkout must be clean" >&2
  exit 1
}

mkdir -p "$output"
target=$output/target
export CARGO_TARGET_DIR=$target
export RUSTFLAGS="-C force-frame-pointers=yes -C debuginfo=2 -C link-arg=-Wl,--build-id=sha1"
export CARGO_PROFILE_RELEASE_DEBUG=2
export CARGO_PROFILE_RELEASE_STRIP=false

{
  echo "source_commit=$actual_commit"
  echo "source_tree=$(git -C "$checkout" rev-parse HEAD^{tree})"
  echo "cargo_lock_blob=$(git -C "$checkout" rev-parse HEAD:Cargo.lock)"
  echo "cargo_lock_sha256=$(sha256sum "$checkout/Cargo.lock" | awk '{print $1}')"
  echo "rustc=$(rustc --version --verbose | tr '\n' ';')"
  echo "cargo=$(cargo --version --verbose | tr '\n' ';')"
  echo "target=$(rustc -vV | awk '/^host:/{print $2}')"
  echo "rustflags=$RUSTFLAGS"
  echo "cargo_profile_release_debug=$CARGO_PROFILE_RELEASE_DEBUG"
  echo "cargo_profile_release_strip=$CARGO_PROFILE_RELEASE_STRIP"
  echo "test_command=cargo test --locked -p azure_data_cosmos_driver_native --all-features"
  echo "build_command=cargo build --locked --release -p azure_data_cosmos_driver_native"
} > "$output/BUILD-PROVENANCE.txt"

(
  cd "$checkout"
  cargo test --locked -p azure_data_cosmos_driver_native --all-features
  cargo build --locked --release -p azure_data_cosmos_driver_native
) 2>&1 | tee "$output/build.log"

archive=$target/release/libazurecosmosdriver.a
[[ -s "$archive" ]]
cp "$archive" "$output/"
sha256sum "$output/libazurecosmosdriver.a" > "$output/SHA256SUMS"
readelf -n "$output/libazurecosmosdriver.a" > "$output/BUILD-ID.txt" 2>&1 || true
{
  echo "archive_sha256=$(sha256sum "$output/libazurecosmosdriver.a" | awk '{print $1}')"
  echo "archive_size=$(stat -c %s "$output/libazurecosmosdriver.a")"
} >> "$output/BUILD-PROVENANCE.txt"
chmod -R a-w "$output"
