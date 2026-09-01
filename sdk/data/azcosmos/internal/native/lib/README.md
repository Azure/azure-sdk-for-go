# TEMPORARY — these binaries are meant to be deleted

The driver libraries in this directory are checked into the repository. That is not the plan of
record and it is not how this module is meant to ship. They are here because CI has no other way
to get a driver library today, and they are expected to be removed.

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
2. Delete this directory, including `.gitignore` and this file.
3. Drop the `-L`/`-rpath` flags for `internal/native/lib` from `cgo_native.go`.

Nothing else here depends on these files, so the revert is self-contained.

## Integrity and provenance

Built from a clean checkout of `Azure/azure-sdk-for-rust` at commit
`2bc179cf6ff5c22e249fd10d944c8e946bb712d2`, which is on `main`, with no local modifications:

```sh
git worktree add --detach /tmp/rustclean 2bc179cf6ff5c22e249fd10d944c8e946bb712d2

# darwin/arm64. The install name has to be @rpath; see below for why.
cd /tmp/rustclean/sdk/cosmos/azure_data_cosmos_driver_native
RUSTFLAGS="-C link-arg=-Wl,-install_name,@rpath/libazurecosmosdriver.dylib \
  --remap-path-prefix=$HOME/.cargo=/cargo \
  --remap-path-prefix=$HOME/.rustup=/rustup \
  --remap-path-prefix=/tmp/rustclean=/src" cargo build --release

# linux/amd64, in a container so the artifact does not depend on the host toolchain.
docker run --rm --platform linux/amd64 -v /tmp/rustclean:/src -w /src rust:1.95-bookworm bash -c '
  R="--remap-path-prefix=/usr/local/cargo=/cargo --remap-path-prefix=/usr/local/rustup=/rustup"
  cd sdk/cosmos/azure_data_cosmos_driver_native
  RUSTFLAGS="-C link-arg=-Wl,-soname,libazurecosmosdriver.so $R" \
    cargo build --release --target-dir /src/target-linux
  cd ../azure_data_cosmos_emulator
  RUSTFLAGS="$R" cargo build --release --target-dir /src/target-linux'
```

The `--remap-path-prefix` flags are not cosmetic. Rust embeds the absolute path of every source
file it can panic in, so without them the artifacts carry several hundred references to the
building user's home directory — which both publishes that directory and makes the bytes differ
per builder. With them there are none, which is worth checking after any rebuild:

```sh
strings -a darwin_arm64/libazurecosmosdriver.dylib | grep -cE '/Users/|/home/'   # expect 0
```

`../azurecosmosdriver.h` is the header that commit generates, vendored unmodified, so the header
and the libraries cannot disagree.

`SHA256SUMS` records the bytes that were reviewed. CI verifies it before it makes the emulator
executable and runs it, so a substituted artifact fails the build rather than executing on an
agent. Regenerate it in this directory whenever a library is replaced:

```sh
shasum -a 256 darwin_arm64/libazurecosmosdriver.dylib \
  linux_amd64/libazurecosmosdriver.so linux_amd64/azure_data_cosmos_emulator > SHA256SUMS
```

Anything replacing these must come from a clean checkout of a public commit, with that commit
recorded here. A binary is not reviewable in a diff, so this section and the checksums are what a
reviewer is actually approving.

## While they are still here

Only `linux/amd64` and `darwin/arm64` are present, because those are the platforms actually built
and tested. Another platform needs its own build dropped into the matching directory — which is
the cost these files impose, and the reason not to grow the set.

The repository root ignores `*.a` and `*.so`, so `.gitignore` here re-enables them. Without it
`git add` skips a rebuilt Linux library silently and CI fails on a library that looks committed
locally. Verify from a clean clone, not from the working tree.

That exemption names the two files exactly rather than un-ignoring the extensions, so an unrelated
`.so` that lands anywhere under this directory stays ignored and cannot be swept in by `git add -A`.

The darwin library's install name is `@rpath/libazurecosmosdriver.dylib`. It must not be an
absolute path: `dyld` resolves an absolute install name directly and never consults the loading
binary's rpaths, so a build-tree path there is both a load-time dependency on a directory that will
not exist and somewhere a local user could plant a library. Check it with
`otool -D darwin_arm64/libazurecosmosdriver.dylib`; CI asserts it.
