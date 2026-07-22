<!-- cspell:ignore amd64 arm64 azcosmos azcosmoscore checksums glibc GOARCH GOOS libc LDFLAGS libazurecosmos librdkafka metapackage musl onnxruntime SRCDIR -->
# Go FFI native-driver distribution design

> **Status:** Design discussion for Go Central SDK and architecture-board review.
>
> This document captures the native-driver distribution problem for a Go SDK
> backed by the Rust Cosmos driver through FFI. It focuses on customer
> experience, Go module layout, native binary placement, and the trade-offs that
> need Go Central SDK alignment.

## 1. Problem statement

The Go v2 FFI path lets Go reuse the Rust Cosmos driver instead of re-owning the
same driver logic in Go. That reduces implementation drift, but it moves a real
decision into distribution:

> Where should the platform-specific native driver binaries live, and how should
> a Go customer get the correct one?

The native driver is built once per supported OS/architecture/libc target. If the
initial supported set includes Windows, macOS, and Linux on x64/amd64 and ARM64,
the SDK needs roughly this matrix:

| Target family | Initial target examples |
|---|---|
| Windows | amd64, arm64 |
| macOS | amd64, arm64 |
| Linux glibc | amd64, arm64 |
| Optional or follow-up | Windows x86, Linux musl/Alpine, additional long-tail targets |

Windows targets also need an ABI/toolchain dimension. The implementation design
is expected to produce both GNU/MinGW-compatible `.a` artifacts for cgo linking
and MSVC-compatible `.lib` artifacts where needed, but this document does not
finalize that artifact split. Before GA, the Windows support matrix must state
which Rust target triples and link artifacts are used for each Go target.

With an optimized native driver around **~5 MB per target**, a six-target matrix
is roughly **~30 MB before compression**. Expanding to ten targets approaches the
earlier **~50 MB** mental model. The design question is not only binary size; it
is **who downloads which targets, when, and how visible the native dependency is
to the customer**.

## 2. Go module mechanics that shape the design

Go has three mechanics that matter here.

First, a GitHub repository and a Go module are not the same thing. A repository
can contain many modules, each rooted at its own `go.mod`. Go customers normally
consume module ZIPs through the Go module cache/proxy; they do not clone the
whole repository just to use a package.

Second, Go builds packages, not an entire repository by default. A cgo
requirement in `azcosmos` does not automatically mean unrelated Azure SDK for Go
modules require cgo. Customers and CI invoke `go build`, `go test`, or `go list`
against specific modules/packages.

Third, cgo link flags are collected from packages in the build graph. The cgo
documentation says `#cgo LDFLAGS` directives from any package in the program are
concatenated at link time. That means a platform-specific driver package can
contribute the native library path, but only if that package is imported by the
program. A `require` entry alone is not enough; the package must be reachable
through an import, often a blank import.

```go
//go:build linux && amd64

package azcosmos

import _ "github.com/Azure/azure-cosmos-go-native-drivers/azcosmos-driver-linux-amd64"
```

Fourth, normal builds should benefit from Go's lazy module loading. A module can
list multiple platform driver modules in `require`, while a single-platform
`go get` / `go build` only needs the driver reached by the active build-tagged
import. Full-graph workflows such as `go mod download all`, vendoring, proxy
prewarming, or CI validation that intentionally walks the complete build list
may still fetch the whole default driver matrix.

## 3. Design goals

| Goal | Why it matters |
|---|---|
| Keep the common path close to `go get` / `go build` | Cosmos Go v1 feels like a normal Go SDK; Go v2 should not surprise mainstream customers. |
| Avoid making every customer download every long-tail binary | The native matrix can grow over time, and module-cache footprint matters. |
| Keep versioning safe | The Go wrapper, C header, and native driver ABI must match. |
| Work in enterprise and offline environments | Build-time network downloads and ad-hoc install scripts are often blocked. |
| Keep the Azure SDK for Go repository manageable | Committing many binary artifacts affects contributors who clone the repo, even if customers use module ZIPs. |
| Make unsupported platforms fail clearly | Customers should see a direct "no driver configured for this platform" message, not an obscure linker failure. |

## 4. Industry reference points

### Confluent Kafka: bundle-first hybrid

`confluent-kafka-go` is the closest market precedent. It is a Go library backed
by the native `librdkafka` implementation. Its developer documentation describes
bundled platform-specific static builds as the default, with dynamic/manual
linking as an escape hatch for special cases.

References:
[`kafka/README.md#build-tags`](https://github.com/confluentinc/confluent-kafka-go/blob/master/kafka/README.md#build-tags)
and
[`librdkafka_vendor`](https://pkg.go.dev/github.com/confluentinc/confluent-kafka-go/v2/kafka/librdkafka_vendor).

The useful lesson is not "copy Confluent exactly." The useful lesson is that a
Go library can be native-backed and still preserve a mostly normal customer
experience by bundling common platform binaries and documenting escape hatches.

### ONNX Runtime Go: wrapper-only, user supplies native runtime

`onnxruntime_go` keeps the Go wrapper separate from the platform runtime. The
customer supplies the matching `onnxruntime.dll`, `.so`, or `.dylib` and points
the wrapper at it before initialization.

Reference:
[`onnxruntime_go` requirements](https://github.com/yalue/onnxruntime_go#requirements).

The useful lesson is that wrapper-only distribution keeps the Go module small,
but it makes native acquisition part of the customer setup. That is often
acceptable in ML/runtime scenarios; it is a harder fit for a first-party Azure
data SDK's default path.

### go-sqlite3: cgo toolchain requirement is explicit

`go-sqlite3` is a widely used Go package that directly documents the cgo and
compiler requirement.

Reference:
[`go-sqlite3` installation](https://github.com/mattn/go-sqlite3#installation).

The useful lesson is that the Go ecosystem accepts cgo in some packages, but it
does not hide the compiler requirement. Cosmos should be equally explicit:
prebuilt Rust artifacts remove the need for a Rust toolchain and manual native
library copying, but **cgo still requires a C build toolchain**.

### Go community discussion: no wheel-like standard

The Go community has discussed prebuilt cgo dependencies and the absence of a
Python-wheel-like mechanism for `go.mod`.

Reference:
[`golang-nuts` discussion](https://groups.google.com/g/golang-nuts/c/ahZXdoClBGg).

The useful lesson is that Go does not provide a standard native-binary packaging
solution. Cosmos has to choose and own a distribution model.

## 5. Option A: one public module bundles the full native matrix

This is the simplest Confluent-like shape.

```text
github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos
├── go.mod
├── client.go
├── internal/core/
├── internal/native/
│   ├── windows-amd64/{libazurecosmos.a or azurecosmos.lib}
│   ├── windows-arm64/{libazurecosmos.a or azurecosmos.lib}
│   ├── darwin-amd64/libazurecosmos.a
│   ├── darwin-arm64/libazurecosmos.a
│   ├── linux-amd64-gnu/libazurecosmos.a
│   └── linux-arm64-gnu/libazurecosmos.a
├── driver_windows_amd64.go
├── driver_darwin_arm64.go
└── driver_linux_amd64.go
```

Selection happens through Go build tags:

```go
//go:build darwin && arm64

package azcosmos

/*
#cgo LDFLAGS: -L${SRCDIR}/internal/native/darwin-arm64 -lazurecosmos
*/
import "C"
```

Customer flow:

```text
customer app imports azcosmos
        │
        ▼
Go downloads one azcosmos module ZIP containing the default native matrix
        │
        ▼
Go build tags select the one matching native library
        │
        ▼
application links with the Rust driver for the current platform
```

**Size gist:** customers download/cache the whole default native matrix as part
of one `azcosmos` module. A six-target default set is roughly ~30 MB before
compression; a ten-target set approaches ~50 MB. The final app links only the
current platform's native driver.

| Area | Effect |
|---|---|
| Customer experience | Best. One `go get`; one module; no platform package choices. |
| Download footprint | Worst. Every `azcosmos` customer downloads the default native matrix. |
| Version safety | Strong. Go wrapper and native bits are versioned together. |
| Repository impact | High if the module lives in `azure-sdk-for-go`; binaries are committed with SDK source. |
| Long-tail targets | Adding one target increases the default module for everyone. |

## 6. Option B: split modules in the Azure SDK for Go repository

This is a multi-module model where the public SDK stays zero-touch for default
platforms, but each native binary lives in its own Go module.

```text
github.com/Azure/azure-sdk-for-go
└── sdk/data/azcosmos-core/
    ├── go.mod
    ├── core.go
    └── include/azurecosmos.h

└── sdk/data/azcosmos-driver-linux-amd64-gnu/
    ├── go.mod
    ├── link_linux_amd64.go
    └── native/libazurecosmos.a

└── sdk/data/azcosmos-driver-darwin-arm64/
    ├── go.mod
    ├── link_darwin_arm64.go
    └── native/libazurecosmos.a

└── sdk/data/azcosmos/
    ├── go.mod
    ├── client.go
    ├── default_driver_linux_amd64.go
    ├── default_driver_darwin_arm64.go
    └── unsupported_driver.go
```

`azcosmos-core` contains the shared Go wrapper and C ABI declarations:

```go
// sdk/data/azcosmos-core/go.mod
module github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos-core

go 1.25.0
```

```c
/* sdk/data/azcosmos-core/include/azurecosmos.h */
#pragma once

const char* azurecosmos_abi_version(void);
void* azurecosmos_client_new(const char* endpoint, const char* key);
void azurecosmos_client_free(void* client);
```

```go
// sdk/data/azcosmos-core/core.go
package azcosmoscore

/*
#cgo CFLAGS: -I${SRCDIR}/include
#include "azurecosmos.h"
*/
import "C"

func ABIVersion() string {
    return C.GoString(C.azurecosmos_abi_version())
}
```

Each driver module contains exactly one native library and the link flags for
that target:

```go
// sdk/data/azcosmos-driver-linux-amd64-gnu/go.mod
module github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos-driver-linux-amd64-gnu

go 1.25.0
```

```go
// sdk/data/azcosmos-driver-linux-amd64-gnu/link_linux_amd64.go
//go:build linux && amd64 && !musl

package azcosmosdriverlinuxamd64gnu

/*
#cgo LDFLAGS: -L${SRCDIR}/native -lazurecosmos
*/
import "C"
```

The public `azcosmos` module depends on the core module and the default driver
modules:

```go
// sdk/data/azcosmos/go.mod
module github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos

go 1.25.0

require (
    github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos-core v1.2.3
    github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos-driver-darwin-arm64 v1.2.3
    github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos-driver-linux-amd64-gnu v1.2.3
    github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos-driver-windows-amd64 v1.2.3
)
```

`azcosmos` imports the active default driver through build-tagged blank imports.
This import is important: a `require` entry alone does not make the driver's
cgo link flags participate in the final link.

```go
// sdk/data/azcosmos/default_driver_linux_amd64.go
//go:build linux && amd64 && !musl

package azcosmos

import _ "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos-driver-linux-amd64-gnu"
```

```go
// sdk/data/azcosmos/default_driver_darwin_arm64.go
//go:build darwin && arm64

package azcosmos

import _ "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos-driver-darwin-arm64"
```

A customer on a default platform writes ordinary Go application code. They do
not import the driver package directly:

```go
// customer-app/go.mod
module example.com/customer-app

go 1.25.0

require github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos v1.2.3
```

```go
// customer-app/main.go
package main

import "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"

func main() {
    _ = azcosmos.NativeDriverVersion()
}
```

For `GOOS=linux GOARCH=amd64`, Go includes `default_driver_linux_amd64.go`.
That file blank-imports `azcosmos-driver-linux-amd64-gnu`, so the driver module
contributes `#cgo LDFLAGS` and the app links against
`native/libazurecosmos.a`.

For libc variants, Go does not provide a built-in `glibc` versus `musl` target
dimension. If both are supported, the design needs an explicit convention such
as a custom `musl` build tag, a separate opt-in package, or a default Linux
choice with the other variant documented as optional.

**Size gist:** for ordinary `go get` / `go build` workflows, customers should
fetch the active platform driver module rather than the full default matrix. A
full-graph command such as `go mod download all`, vendoring, proxy prewarming, or
CI validation can still fetch all default driver modules. The final app links
only the current platform's native driver.

| Area | Effect |
|---|---|
| Customer experience | Good. Common platforms still use `go get azcosmos` / `go build`. |
| Download footprint | Normal builds fetch the active driver module; full-graph workflows can fetch the default driver matrix. |
| Version safety | Manageable if all modules are released together and exact versions are pinned. |
| Repository impact | Still high for `azure-sdk-for-go` contributors because the repo contains all native modules. |
| Long-tail targets | Better. Optional targets can be extra modules not included in default `azcosmos`. |

## 7. Option C: split modules, native drivers in a separate repository

This keeps the customer-facing SDK in `azure-sdk-for-go`, but moves the binary
payload modules to a separate Azure-owned repository. This is the preferred
direction for further design because it keeps the large binaries out of the
Azure SDK for Go repository without changing the common customer flow.

```text
github.com/Azure/azure-sdk-for-go
└── sdk/data/azcosmos/
└── sdk/data/azcosmos-core/

github.com/Azure/azure-cosmos-go-native-drivers
└── azcosmos-driver-windows-amd64/
└── azcosmos-driver-windows-arm64/
└── azcosmos-driver-darwin-amd64/
└── azcosmos-driver-darwin-arm64/
└── azcosmos-driver-linux-amd64-gnu/
└── azcosmos-driver-linux-arm64-gnu/
```

Customer flow is still Go-native:

```text
customer runs:
  go get github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos

Go module graph includes:
  azcosmos
  azcosmos-core
  github.com/Azure/azure-cosmos-go-native-drivers/azcosmos-driver-darwin-arm64

customer does not clone the native-driver repo manually
```

The public SDK can hide the separate repository behind build-tagged blank
imports:

```go
//go:build linux && amd64

package azcosmos

import _ "github.com/Azure/azure-cosmos-go-native-drivers/azcosmos-driver-linux-amd64-gnu"
```

**Size gist:** for customers, this is similar to Option B. Ordinary builds fetch
the active driver module from the separate native-driver repository; full-graph
workflows can fetch the full default matrix. The final app links only one target.

| Area | Effect |
|---|---|
| Customer experience | Good for default platforms if `azcosmos` wires driver imports automatically. |
| Download footprint | Same customer behavior as Option B, but the active driver module comes from a separate repository. |
| Version safety | Requires stronger release discipline across repositories. |
| Repository impact | Better for Azure SDK for Go contributors; binary churn lives elsewhere. |
| Governance | Needs clear ownership, release, security, signing, and support boundaries. |

This model separates two concerns:

```text
Go customers care about modules.
Azure SDK contributors care about repositories.
```

Separate modules help customer download granularity. A separate repository helps
the Azure SDK for Go repository avoid carrying binary payloads in its Git
history.

## 8. Advanced mode: customer-managed dynamic library

This is not the default path. It is an advanced mode for customers who want or
need to manage the native library through their own deployment system.

If the Go wrapper links dynamically, a customer could choose not to use an
`azcosmos-driver-*` package and instead use `azcosmos-core` directly. As long as
the matching `libazurecosmos` is available on the platform's library load path,
the application can link/load the native driver through the customer's packaging
system.

```text
customer app
  imports azcosmos-core or an advanced azcosmos configuration path
  places libazurecosmos on the library load path
  builds/runs with the customer's native package manager or image layout
```

**Size gist:** the Go module stays small because it does not carry the native
matrix. Customers still need one native driver for their platform, but they
acquire, place, and update it themselves.

| Area | Effect |
|---|---|
| Customer experience | Good only for customers that already manage native dependencies. |
| Download footprint | Minimal Go module footprint. |
| Version safety | Must validate SDK/native ABI compatibility at startup. |
| Enterprise/offline | Useful for strict packaging environments. |
| Supportability | Native setup becomes a customer-owned prerequisite. |

This should be documented as an escape hatch, not the default Azure SDK
experience.

## 9. Options not selected for the default path

### Small SDK module plus GitHub release assets

This model keeps the Go SDK module small and publishes native drivers as GitHub
release assets:

```text
GitHub release v1.2.3
├── azurecosmos-native-windows-amd64.zip
├── azurecosmos-native-darwin-arm64.tar.gz
├── azurecosmos-native-linux-amd64-gnu.tar.gz
└── checksums.txt
```

It gives customers a smaller per-platform download, but requires an explicit
download/install step or helper command before `go build`. That is too
disruptive for the default Go SDK experience and introduces extra proxy, cache,
checksum, install-path, and linker-path failure modes.

### Standalone native-backed Go SDK

A separate public SDK identity, such as `azcosmosnative`, would avoid impacting
existing `azcosmos` users, but it confuses the support and migration matrix. If
Go v2 is not shipped as the normal `azcosmos` evolution, customers must
understand which package is the supported future and how long the old package
continues. That is a product-level split rather than a packaging solution, so it
is not selected here.

## 10. Default driver set selection

If a split-module model is used, there are two distinct decisions:

1. **Which driver modules are published?**
2. **Which of those are included in the default `azcosmos` experience?**

For example:

```text
Published modules:
  windows-amd64
  windows-arm64
  darwin-amd64
  darwin-arm64
  linux-amd64-gnu
  linux-arm64-gnu
  windows-386
  linux-amd64-musl

Default azcosmos imports:
  windows-amd64
  windows-arm64
  darwin-amd64
  darwin-arm64
  linux-amd64-gnu
  linux-arm64-gnu

Optional:
  windows-386
  linux-amd64-musl
```

Build-time selection then happens through build tags:

```text
GOOS=darwin GOARCH=arm64
       │
       ▼
default_driver_darwin_arm64.go is included
       │
       ▼
azcosmos-driver-darwin-arm64 is imported
       │
       ▼
darwin/arm64 native library contributes link flags
```

For a platform outside the default set, the SDK should fail clearly:

```text
azcosmos: no native Cosmos driver package is configured for linux/amd64/musl.
Add the linux/amd64/musl driver package or use a supported default platform.
```

Optional targets need a documented activation story. Without an import, Go may
know about the module version in `go.mod`, but the driver's cgo link flags will
not participate in the final link.

## 11. Versioning and ABI safety

All models need a version contract between:

- Go wrapper package
- `azurecosmos.h`
- native Rust driver library
- C ABI version

**Go module versions cannot enforce this contract.** A `require` line sets a
*minimum* version, not a pin. Under Minimal Version Selection (MVS), any other
module in the customer's build graph can raise a driver module to a newer
version than the wrapper declared. The wrapper has no way to demand exact
module-version equality, so correctness must rest on an **ABI-compatibility
contract**, not on matching module version strings.

The rule is: the wrapper declares the native ABI it needs, and any
MVS-selected driver is acceptable as long as it stays ABI-compatible.

```text
azcosmos v1.2.3
  requires native ABI major 1, minimum minor 2
  accepts any driver module whose native ABI is 1.x with x >= 2
  (MVS may select azcosmos-driver-* v1.2.3, v1.3.0, or newer —
   all valid while the ABI stays 1.x compatible)
```

ABI compatibility follows the usual rule: same major, minor greater-than-or-equal
to what the wrapper needs. A major bump means incompatible; the wrapper must
reject it.

At minimum, the native driver should expose an ABI/version function:

```c
const char* azurecosmos_abi_version(void);
```

The Go wrapper validates it during initialization against a compatibility
*range*, not an exact string, and fails only on genuine incompatibility:

```text
azcosmos: native driver ABI incompatible:
  Go wrapper requires ABI major 1, minor >= 2
  linked native driver reports ABI 2.0  (major mismatch)
```

A newer-but-compatible driver (for example ABI 1.4 when the wrapper needs 1.2)
must pass. This validation matters most for manual, dynamic, or optional-driver
paths, and still catches packaging mistakes on bundled paths.

## 12. Customer experience comparison

| Model | Common-platform customer flow | Customer-visible native setup | Download behavior |
|---|---|---|---|
| A. Single bundled module | `go get azcosmos`; `go build` | cgo toolchain only | Downloads whole default matrix in one module |
| B. Split modules, same repo | `go get azcosmos`; `go build` | cgo toolchain only | Normal builds fetch active driver module; full-graph workflows can fetch default matrix |
| C. Split modules, separate repo | `go get azcosmos`; `go build` | cgo toolchain only | Same as Option B, but the active driver module comes from a separate repo |
| Advanced dynamic library | imports/configures advanced path | customer-managed native library | Downloads only Go wrapper; customer supplies native payload |

## 13. Repository and release comparison

| Model | Azure SDK for Go repo impact | Release coordination | Reviewability |
|---|---|---|---|
| A. Single bundled module | Highest; all binaries live under one module | Simple, one module version | Large binary diffs in SDK PRs |
| B. Split modules, same repo | High; binaries still live in repo, but module payloads are separated | Moderate; many modules, same repo | Binary PRs still affect SDK repo |
| C. Split modules, separate repo | Lower; native payload outside main SDK repo | Higher; cross-repo version alignment | Cleaner SDK PRs, separate native PRs |
| Advanced dynamic library | Low source impact | Higher customer responsibility | Native setup mostly outside SDK PRs |

## 14. Discussion checkpoints for Go Central SDK review

These are the questions that likely need board-level alignment:

1. Is the default Azure SDK experience allowed to require cgo and a C toolchain
   when prebuilt native driver artifacts are included?
2. Should the default Cosmos Go module include all mainstream platform drivers,
   or should some platforms be opt-in?
3. Should native driver packages live in a separate Azure-owned repository to
   keep binary payload out of `azure-sdk-for-go`?
4. Under MVS, the customer graph can select a newer driver module than the
   wrapper declared. What is the ABI-compatibility policy that keeps this safe —
   same-major / minimum-minor, and who owns bumping the major?
5. Does the Azure SDK release system support synchronized versioning across a
   public SDK module, core wrapper module, and several driver modules?
6. How should optional/long-tail platforms be activated: direct blank import,
   Azure SDK metapackage, dynamic/manual native path, or unsupported?
7. What is the signing, checksum, and provenance story for native artifacts?
8. What exact error should customers see when no driver is configured for their
   target?

## 15. Current read

The split-module design with native driver modules in a **separate Azure-owned
repository** is the strongest candidate for discussion. It preserves the
common-platform customer experience while keeping large binary artifacts out of
the Azure SDK for Go repository.

```text
Common path:
  go get azcosmos
  go build
  default driver selected through build tags

Repository placement:
  azcosmos Go code in azure-sdk-for-go
  native driver modules in a separate Azure-owned repository
```

For customer experience, the separate-repo option can still be nearly invisible
if the driver packages are normal Go modules and the public `azcosmos` package
imports the default drivers automatically.
