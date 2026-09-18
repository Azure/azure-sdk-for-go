# Native emulator artifact

The Linux AMD64 emulator in this directory is checked into the repository so the native-driver
tests can run without building the Rust workspace in the Go pipeline. It is separate from the
target-specific driver archives distributed by
[azure-cosmos-driver](https://github.com/Azure/azure-cosmos-driver).

The emulator was built from a clean checkout of `Azure/azure-sdk-for-rust` at commit
`a05f024d22cd0a70f6ccac5e310cbbbaf40bf20e`, which is on `main`, with no local modifications:

```sh
git worktree add --detach /tmp/rustclean a05f024d22cd0a70f6ccac5e310cbbbaf40bf20e

docker run --rm --platform linux/amd64 -v /tmp/rustclean:/src -w /src rust:1.95-bookworm bash -c '
  R="--remap-path-prefix=/usr/local/cargo=/cargo --remap-path-prefix=/usr/local/rustup=/rustup"
  cd sdk/cosmos/azure_data_cosmos_emulator
  RUSTFLAGS="$R" cargo build --release --target-dir /src/target-linux'

cd /path/to/azure-sdk-for-go/sdk/data/azcosmos
cp /tmp/rustclean/target-linux/release/azure_data_cosmos_emulator \
  internal/native/lib/linux_amd64/
```

`SHA256SUMS` records the reviewed emulator bytes. CI verifies it before making the binary
executable, detecting checkout corruption within the same code-review trust boundary. Regenerate
the checksum whenever the emulator is replaced:

```sh
shasum -a 256 internal/native/lib/linux_amd64/azure_data_cosmos_emulator \
  > internal/native/lib/SHA256SUMS
```
