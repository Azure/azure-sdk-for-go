# TEMPORARY — these binaries are meant to be deleted

The driver archives in the target-specific packages beside this directory and the emulator here
are checked into the repository. That is not the plan of record and it is not how this module is
meant to ship. They are here because CI has no other way to get a driver today, and they are
expected to be removed.

## Why they are here

The driver is built from
[azure_data_cosmos_driver_native](https://github.com/Azure/azure-sdk-for-rust/tree/main/sdk/cosmos/azure_data_cosmos_driver_native),
which is Rust. Building it in CI would mean a Rust toolchain and a `cargo` fetch from crates.io in
the Go pipeline, for a dependency the shipped module is not supposed to build from source in the
first place.

## What replaces them

The distribution design puts each target's library in its own Go module in the
[azure-cosmos-driver](https://github.com/Azure/azure-cosmos-driver) repository, selected by
`GOOS`/`GOARCH` build constraints. That repository exists and already carries a `darwin/arm64`
module, but not yet one for `linux/amd64`, which is the platform CI runs on.

## When to delete this directory

**As soon as `azure-cosmos-driver` carries a module for every platform this module builds on.**
At that point:

1. Add the driver modules to `go.mod`.
2. Delete the temporary `internal/native/{darwinarm64,linuxamd64}` packages and this directory.
3. Move the archive link directives from `cgo_native.go` into the per-target modules.

Nothing else here depends on these files, so the revert is self-contained.

## Integrity and provenance

Built from a clean checkout of `Azure/azure-sdk-for-rust` at commit
`2bc179cf6ff5c22e249fd10d944c8e946bb712d2`, which is on `main`, with no local modifications:

```sh
git worktree add --detach /tmp/rustclean 2bc179cf6ff5c22e249fd10d944c8e946bb712d2

# darwin/arm64.
cd /tmp/rustclean/sdk/cosmos/azure_data_cosmos_driver_native
MACOSX_DEPLOYMENT_TARGET=12.0 \
CFLAGS="-mmacosx-version-min=12.0" \
CXXFLAGS="-mmacosx-version-min=12.0" \
RUSTFLAGS="-C link-arg=-mmacosx-version-min=12.0 \
  --remap-path-prefix=$HOME/.cargo=/cargo \
  --remap-path-prefix=$HOME/.rustup=/rustup \
  --remap-path-prefix=/tmp/rustclean=/src" cargo build --release

# linux/amd64, built against the same glibc ABI used by the supported Go build.
docker run --rm --platform linux/amd64 -v /tmp/rustclean:/src -w /src rust:1.95-bookworm bash -c '
  R="--remap-path-prefix=/usr/local/cargo=/cargo --remap-path-prefix=/usr/local/rustup=/rustup"
  cd sdk/cosmos/azure_data_cosmos_driver_native
  RUSTFLAGS="$R" cargo build --release --target-dir /src/target-linux
  cd ../azure_data_cosmos_emulator
  RUSTFLAGS="$R" cargo build --release --target-dir /src/target-linux'

# Install from the azcosmos module root.
cd /path/to/azure-sdk-for-go/sdk/data/azcosmos
cp /tmp/rustclean/target/release/libazurecosmosdriver.a \
  internal/native/darwinarm64/libazurecosmosdriver.syso
cp /tmp/rustclean/target-linux/release/libazurecosmosdriver.a \
  internal/native/linuxamd64/libazurecosmosdriver.syso
cp /tmp/rustclean/target-linux/release/azure_data_cosmos_emulator \
  internal/native/lib/linux_amd64/
cp /tmp/rustclean/sdk/cosmos/azure_data_cosmos_driver_native/include/azurecosmosdriver.h .
```

The `--remap-path-prefix` flags are not cosmetic. Rust embeds the absolute path of every source
file it can panic in, so without them the artifacts carry several hundred references to the
building user's home directory — which both publishes that directory and makes the bytes differ
per builder. With them there are none, which is worth checking after any rebuild:

```sh
strings -a internal/native/darwinarm64/libazurecosmosdriver.syso |
  grep -c "$HOME/"   # expect 0
```

`azurecosmosdriver.h` is the header that commit generates, vendored unmodified, so the header and
the libraries cannot disagree.

`SHA256SUMS` records the bytes that were reviewed. CI verifies it before it makes the emulator
executable and runs it, so a substituted artifact fails the build rather than executing on an
agent. Regenerate it in this directory whenever a library is replaced:

```sh
shasum -a 256 internal/native/darwinarm64/libazurecosmosdriver.syso \
  internal/native/linuxamd64/libazurecosmosdriver.syso \
  internal/native/lib/linux_amd64/azure_data_cosmos_emulator > internal/native/lib/SHA256SUMS
```

Anything replacing these must come from a clean checkout of a public commit, with that commit
recorded here. A binary is not reviewable in a diff, so this section and the checksums are what a
reviewer is actually approving.

## While they are still here

Only `linux/amd64` and `darwin/arm64` are present, because those are the platforms actually built
and tested. Another platform needs its own `.syso` file in the module root — which is the cost these
files impose, and the reason not to grow the set.

The `.syso` suffix is required rather than merely convenient. Go preserves platform-specific
`.syso` files in module zips and `go mod vendor`, then links the matching file automatically.
Keeping an ordinary `.a` under this subdirectory fails for vendored consumers because Go drops the
entire native subtree.
