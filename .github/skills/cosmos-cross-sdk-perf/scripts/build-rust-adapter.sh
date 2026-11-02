#!/usr/bin/env bash
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words Eeuo esac rustc rustflags serde

set -Eeuo pipefail

usage() {
  echo "usage: build-rust-adapter.sh --checkout <azure-sdk-for-rust> --commit <sha> --output <dir>" >&2
  exit 2
}

checkout= commit= output=
while [[ $# -gt 0 ]]; do
  case "$1" in
    --checkout) checkout=$(readlink -f "$2"); shift 2 ;;
    --commit) commit=$2; shift 2 ;;
    --output) output=$(readlink -m "$2"); shift 2 ;;
    *) usage ;;
  esac
done
[[ "$commit" =~ ^[0-9a-f]{40}$ && -d "$checkout" && -n "$output" ]] || usage
[[ ! -e "$output" ]] || { echo "refusing to overwrite $output" >&2; exit 1; }
[[ "$(git -C "$checkout" rev-parse HEAD)" == "$commit" ]] || {
  echo "Rust checkout commit mismatch" >&2
  exit 1
}
[[ -z "$(git -C "$checkout" status --short)" ]] || {
  echo "Rust checkout must be clean" >&2
  exit 1
}

skill=$(cd "$(dirname "$0")/.." && pwd)
mkdir -p "$output/work" "$output/bin"
git -C "$checkout" worktree add --detach "$output/work/sdk" "$commit"
trap 'git -C "$checkout" worktree remove --force "$output/work/sdk" 2>/dev/null || true' EXIT

crate=$output/work/sdk/cosmos-perf-adapter
mkdir -p "$crate/src"
cp "$skill/adapters/rust/src/"*.rs "$crate/src/"
cat > "$crate/Cargo.toml" <<'EOF'
[package]
name = "cosmos-perf-adapter"
version = "0.0.0"
edition = "2021"
publish = false

[dependencies]
azure_core = { path = "../sdk/core/azure_core" }
azure_data_cosmos = { path = "../sdk/cosmos/azure_data_cosmos", features = ["key_auth"] }
serde = { workspace = true, features = ["derive"] }
serde_json = { workspace = true }
tokio = { workspace = true, features = ["rt-multi-thread", "macros", "time"] }
EOF

python3 - "$output/work/sdk/Cargo.toml" <<'PY'
import pathlib
import sys

path = pathlib.Path(sys.argv[1])
text = path.read_text(encoding="utf-8")
marker = "members = ["
index = text.index(marker) + len(marker)
text = text[:index] + '\n  "cosmos-perf-adapter",' + text[index:]
path.write_text(text, encoding="utf-8")
PY

export CARGO_TARGET_DIR=$output/work/target
export RUSTFLAGS="-C force-frame-pointers=yes -C debuginfo=2 -C link-arg=-Wl,--build-id=sha1"
export CARGO_PROFILE_RELEASE_DEBUG=2
export CARGO_PROFILE_RELEASE_STRIP=false
cp "$output/work/sdk/Cargo.lock" "$output/Cargo.lock.original"
(
  cd "$output/work/sdk"
  cargo test -p cosmos-perf-adapter
  git diff -- Cargo.lock > "$output/Cargo.lock.diff"
  python3 - "$output/Cargo.lock.diff" <<'PY'
import pathlib
import sys

diff = pathlib.Path(sys.argv[1]).read_text(encoding="utf-8")
if 'name = "cosmos-perf-adapter"' not in diff:
    raise SystemExit("Cargo.lock did not add the local adapter package")
removed = [
    line for line in diff.splitlines()
    if line.startswith("-") and not line.startswith("---")
]
if removed:
    raise SystemExit("Cargo.lock changed existing pinned dependency entries")
PY
  cp Cargo.lock "$output/Cargo.lock.resolved"
  cargo build --locked --release -p cosmos-perf-adapter
)
cp "$CARGO_TARGET_DIR/release/cosmos-perf-adapter" "$output/bin/cosmos-perf-rust"

build_id=$(readelf -n "$output/bin/cosmos-perf-rust" | awk '/Build ID:/ {print $3; exit}')
python3 - "$output/bin/cosmos-perf-rust" "$commit" \
  "$(git -C "$checkout" rev-parse HEAD^{tree})" "$build_id" \
  "$output/Cargo.lock.original" "$output/Cargo.lock.resolved" <<'PY'
import hashlib
import json
import pathlib
import subprocess
import sys

binary, commit, tree, build_id, original_lock, resolved_lock = sys.argv[1:]
path = pathlib.Path(binary)
original_lock_path = pathlib.Path(original_lock)
resolved_lock_path = pathlib.Path(resolved_lock)
value = {
    "leg": "rust",
    "source": {
        "repository": "https://github.com/Azure/azure-sdk-for-rust.git",
        "commit": commit,
        "tree": tree,
        "dirty": False,
    },
    "compiler": {
        "rustc": subprocess.check_output(
            ["rustc", "--version", "--verbose"], text=True
        ).strip(),
        "cargo": subprocess.check_output(["cargo", "--version"], text=True).strip(),
    },
    "buildCommand": [
        "cargo test -p cosmos-perf-adapter",
        "cargo build --locked --release -p cosmos-perf-adapter",
    ],
    "dependencyVersions": {
        "azure_data_cosmos": "pinned-workspace",
        "azure_data_cosmos_driver": "pinned-workspace",
        "cargoLockOriginalSha256": hashlib.sha256(
            original_lock_path.read_bytes()
        ).hexdigest(),
        "cargoLockResolvedSha256": hashlib.sha256(
            resolved_lock_path.read_bytes()
        ).hexdigest(),
    },
    "binary": {
        "sha256": hashlib.sha256(path.read_bytes()).hexdigest(),
        "buildId": build_id,
        "debugSymbols": "embedded",
    },
}
path.with_suffix(path.suffix + ".provenance.json").write_text(
    json.dumps(value, indent=2, sort_keys=True) + "\n", encoding="utf-8"
)
PY
sha256sum "$output/bin/cosmos-perf-rust" > "$output/SHA256SUMS"
git -C "$checkout" worktree remove --force "$output/work/sdk"
trap - EXIT
rm -rf "$output/work"
chmod -R a-w "$output"
