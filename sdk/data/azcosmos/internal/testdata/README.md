# Building the test emulator

The Go native-driver tests use `azure_data_cosmos_emulator` from
[Azure/azure-sdk-for-rust](https://github.com/Azure/azure-sdk-for-rust).
This crate is not published to crates.io. Build it from source; it is separate from the
target-specific driver libraries consumed as Go modules.

The `DriverEmulator` stage in [`ci.yml`](../../ci.yml) pins:

- Rust source revision: `a05f024d22cd0a70f6ccac5e310cbbbaf40bf20e`.
- Rust toolchain: `1.95.0`.
- Dependency versions: the source revision's `Cargo.lock`, enforced with `--locked`.

Building requires Git, rustup, a C compiler/linker, and network access to GitHub,
Rust toolchain downloads, and Cargo dependencies. CI builds on glibc Linux/amd64.
For local use, build on a platform supported by the Go native binding (glibc Linux/amd64
or macOS/arm64).

From a directory outside the Go checkout, with rustup installed:

```sh
rustup toolchain install 1.95.0 --profile minimal
git init cosmos-emulator-source
git -C cosmos-emulator-source fetch --depth=1 \
  https://github.com/Azure/azure-sdk-for-rust.git \
  a05f024d22cd0a70f6ccac5e310cbbbaf40bf20e
git -C cosmos-emulator-source checkout --detach FETCH_HEAD
cd cosmos-emulator-source
cargo +1.95.0 build --locked --release -p azure_data_cosmos_emulator \
  --target-dir ../emulator-target
```

Start `emulator-target/release/azure_data_cosmos_emulator` with `--config` pointing to
[`emulator-config.json`](emulator-config.json). Use the `accountEndpoint` in the stdout
JSON ready record as `AZCOSMOS_ENDPOINT`, and set `EMULATOR=true` for the Go tests.
When running in Docker, keep the emulator and Go tests in the same network namespace:
the configured listeners bind loopback.

When updating the emulator, change the pins in CI and this document together, then rerun
the Go emulator suite. A successful source build alone is not a compatibility check.
