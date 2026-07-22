<!-- cspell:ignore amd64 azcosmos cgo glibc musl -->
# Go v2 Cosmos SDK FFI decision

> **Status:** Proposed
>
> This document records the Go SDK direction to use the Rust Cosmos driver
> through FFI for the Go v2 implementation. It is intended to be reviewed with
> the Go Central SDK team and used as the decision record for the Go package.

## Context

Go v2 needs a path that can deliver in the current timeframe while staying close
to the Rust Cosmos driver's behavior. A time-boxed pure-Go spike proved that a
gateway-mode core path can be ported and validated, but it also clarified the
central trade-off:

> Is avoiding cgo and native packaging worth re-owning driver logic that Rust
> already implements?

Both paths are technically feasible:

| Path | Shape |
|---|---|
| **FFI path** | Go calls the prebuilt Rust Cosmos driver through cgo and a C ABI. |
| **Pure-Go port** | Go reimplements the driver behavior directly in Go. |

## Decision

- Go v2 proceeds with the **Rust-driver FFI path** for the current delivery
  window.
- The Go binding consumes the **prebuilt Rust native driver** through cgo and the
  C ABI surface.
- Packaging mechanics, platform matrix, signing, ABI versioning, and native
  artifact placement are handled by a separate distribution design.

## Why FFI for Go v2

FFI has the better short-term delivery shape:

- **Parity:** Go consumes the same driver implementation as Rust.
- **Velocity:** new Rust-driver behavior does not need to be re-ported before Go
  can benefit.
- **Risk containment:** the hard driver logic stays in one implementation.
- **Testing leverage:** the Go binding can validate language-level behavior
  against the Rust driver's known behavior instead of re-proving the full driver.

The main cost is not the Go wrapper itself. The main cost is **Go packaging and
customer deployment experience**.

## What the pure-Go spike proved

The pure-Go spike showed that a gateway-mode core path can be ported in
idiomatic Go:

- `CGO_ENABLED=0`
- zero third-party dependencies
- point create/read
- retry and routing behavior
- basic failover behavior
- cross-partition unordered query fan-out
- unsupported-query rejection parity

The spike was validated with a scenario-based differential harness. The harness
put an HTTP endpoint in front of the Rust driver's in-memory emulator so the Go
port and the Rust driver could run against the same emulator-backed state and
compare caller-visible outcomes.

This is valuable evidence: the Go port is credible. It is also a useful
validation reference for Go-visible behavior. It is not the selected Go v2
delivery path because every Rust-driver change would need to be tracked, ported,
and revalidated in Go.

## Packaging implications

Cosmos Go v1 is consumed like a normal Go SDK:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos
go build ./...
```

There is no Rust compiler, native shared library, or post-install step in the
customer path. If Go v2 uses FFI, prebuilt Rust artifacts can preserve that part
of the experience, but cgo still requires a C build toolchain. The packaging bar
is therefore:

> Can the Go SDK still feel like one normal Go module across the supported
> Windows, macOS, and Linux targets?

That means native Rust driver artifacts must be produced, versioned, selected,
and linked for the supported Go OS/architecture matrix.

The initial packaging discussion is scoped to mainstream targets:

| OS family | Target combinations discussed |
|---|---|
| Windows | x64/amd64, ARM64, possibly x86 as a consideration |
| macOS | ARM64/Apple Silicon, x64/Intel |
| Linux glibc | x64/amd64, ARM64 |

Linux musl/Alpine is intentionally not counted in the baseline until the Go
packaging plan decides whether it is in scope for the initial release.

With an estimated **~5 MB optimized native binary per platform**, a six-target
matrix is roughly **~30 MB before compression**. Additional targets add to that
footprint. This is manageable for an initial release, but it is a real packaging
design point and should not be treated as an implementation detail.

## Alternatives considered

### Pure-Go driver port

The pure-Go path has a strong developer-experience story:

- `CGO_ENABLED=0`
- normal Go cross-compilation
- no native runtime dependency
- Go standard library transport/TLS/runtime

It is not selected for Go v2 because it moves the cost into ownership:

- driver behavior must be reimplemented in Go
- every Rust-driver change has to be tracked, ported, and revalidated
- long-term drift becomes the main risk
- harder areas still need evidence: richer query execution,
  continuation/split-resume behavior, broader multi-region behavior, change
  feed, and bulk/batch

### Manual native-library installation as the default

Manual native-library installation is not selected for the common path because it
breaks the normal `go get` / `go build` expectation for a first-party Go SDK.
It may still be useful as an advanced mode for customers with strict native
packaging systems.

### Pure-Go downloader shim

A downloader shim is not selected as the default path. Go modules do not provide
a standard post-install hook, and build-time network downloads are often blocked
in enterprise and offline environments.

## Consequences

- Go v2 gets the same driver implementation as Rust, reducing short-term feature
  drift and implementation risk.
- Go v2 takes a dependency on `CGO_ENABLED=1` for the native driver path.
- Customers should not manually install a Rust toolchain or copy native
  libraries for the common path, but cgo still requires an available C build
  toolchain.
- Differential validation remains important so Go-visible behavior can be tested
  against the Rust driver.

## Non-goals

- This document does not define the final Go packaging shape.
- This document does not require customers to manually install native libraries
  for the common path.
- This document does not introduce direct mode; both the Rust driver and Go v2
  path discussed here are gateway-mode SDKs.
- This document does not remove the pure-Go spike as a reference or validation
  asset.
