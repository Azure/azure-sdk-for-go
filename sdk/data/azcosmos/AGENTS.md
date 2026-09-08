---
applyTo: "sdk/data/azcosmos/**/*.*"
---

# Azure Cosmos DB SDK for Go v2 Instructions

This file contains guidance specific to the Azure Cosmos DB SDK for Go v2
(`sdk/data/azcosmos`). Follow the repository-level `AGENTS.md` and
[Azure SDK Design Guidelines for Go](https://azure.github.io/azure-sdk/golang_introduction.html)
unless a Cosmos-specific rule below is more restrictive.

The v2 module is under construction. Preserve the incremental API-review approach: do not imply
that an operation works until its native-driver binding and tests have landed.

## Architecture

v2 is an idiomatic Go API over the shared Rust Cosmos driver. Keep responsibilities separated:

- The Go SDK owns its public API, Go types, credentials, argument validation, and translation
  between Go values and the driver ABI. Callers own application-item serialization.
- The Rust driver owns transport, routing, retries, session handling, failover, query fan-out, and
  Cosmos wire encoding.
- Do not duplicate driver policy in Go. Validate stable Go API contracts at the boundary, but leave
  service topology and protocol decisions to the driver.
- Do not expose native-driver or FFI types through the public Go API.
- Keep item payloads schema-agnostic. The Go API accepts and returns encoded `[]byte` values, and the
  binding passes them through unchanged. Neither the binding nor the driver should inspect
  application-defined item fields unless an operation explicitly requires a value such as an item
  ID.

The Rust-side architectural source of truth is
[`sdk/cosmos/AGENTS.md`](https://github.com/Azure/azure-sdk-for-rust/blob/main/sdk/cosmos/AGENTS.md).
When changing behavior that crosses the language boundary, check the driver contract there and in
the corresponding Rust implementation.

## Public API Design

- Use Go conventions: `context.Context` first, required parameters next, and an options pointer
  last. A nil options pointer must select defaults.
- Prefer useful zero values. When zero has domain meaning, document it and preserve it across the
  ABI. Do not collapse "unset", "default", and explicit false when the driver distinguishes them.
- Use immutable value semantics for domain values such as `PartitionKey`. Copy slices or maps when
  retaining caller-owned data.
- Model databases, containers, partition keys, session tokens, routing strategies, and responses as
  distinct types.
- Keep partition keys explicit and type-safe. Preserve the distinction between null and undefined,
  and support hierarchical partition keys without parsing item bodies.
- Use `ID`, `URL`, `HTTP`, and other initialisms consistently in exported names.
- Add concise GoDoc to every exported symbol. Document nil behavior, zero values, error conditions,
  partition scope, and request-unit implications when they affect callers.
- Do not embed `azcore.ClientOptions` unless the native driver actually honors every exposed field.
  Advertising azcore HTTP-pipeline controls that v2 bypasses is misleading.
- Prefer Microsoft Entra ID authentication. Account-key APIs must not log, format, or include keys
  in errors.

## Client Lifetime and Concurrency

- `Client` is long-lived, safe for concurrent use, and owns native resources and caches. Do not
  introduce per-operation client construction.
- Every operation using native state must acquire the client lifetime guard and release it on every
  path. `Close` must wait for in-flight operations before releasing native memory.
- `Close` remains idempotent and concurrency-safe, and all callers observe the same teardown result.
- An operation started after `Close` reports `CodeClientClosed`. Preserve the established ordering:
  deterministic argument errors first, then closed-client state, then context cancellation.
- Never allow a Go call to use a native handle after its owner has been closed.

## FFI Boundary

- Pin the native driver version deliberately. Update `nativeDriverVersion`, native artifacts, ABI
  declarations, and compatibility tests together.
- Keep ownership explicit for every pointer, buffer, string, callback, and handle crossing the ABI.
  Document who allocates, who frees, and how long borrowed memory remains valid.
- Do not let native code retain a Go pointer after a cgo call. Copy or allocate boundary data when
  lifetimes do not naturally fit the call.
- Convert native failures into Go errors at the boundary. A panic, unwinding exception, or
  language-specific error must never cross the C ABI.
- Minimize cgo crossings in hot paths and avoid unnecessary conversions or allocations, but do not
  trade away ownership safety for fewer copies.
- Keep FFI implementation details unexported. Public APIs should remain idiomatic Go and should not
  require callers to understand native resource management.

## Errors and Responses

- Return `*Error` for classified Cosmos operation failures. Callers use `errors.As` and branch on
  `Error.Code`; the text from `Error.Error()` is not a stable contract.
- Preserve standard Go error behavior. Wrapped errors must support `errors.Is`, especially for
  context cancellation.
- Populate status code, substatus, request charge, activity ID, retry delay, session token, ETag,
  response body, and wire-origin information when the driver provides them.
- Return the zero response value when an operation fails. Metadata for a failed request belongs on
  `*Error`, including request charge because failed requests can consume RUs.
- Never infer whether the service received a request from the status code alone; preserve the
  driver's explicit wire-origin signal.
- Return errors instead of panicking for caller-controlled inputs. Panic only for impossible
  internal invariants.

## Code Organization

The package currently uses a flat layout. Keep related responsibilities recognizable:

- `client.go`, `database_client.go`, and `container_client.go`: resource clients and lifetime
- `*_item.go`: operation APIs and operation-specific options
- `driver.go`: native-driver versioning and binding support
- `errors.go`: public error classification and translation
- `response.go`: successful operation metadata and bodies
- `partition_key.go`, `session_token.go`, and strategy files: domain values
- `valuelist/`: the generator for closed sets of string values

Do not manually edit files marked `Code generated ... DO NOT EDIT`; update the generator input and
regenerate them.

## Testing

- Use `github.com/stretchr/testify/require` for assertions.
- Unit tests should emphasize boundary contracts: argument validation, zero and nil behavior,
  ownership and copying, error translation, ABI conversions, client lifetime, and concurrency.
- Use table-driven tests for equivalent input classes and exact assertions for complete contracts.
- Test cancellation, `Close` races, and native-resource lifetimes with the race detector.
- Do not add tests for trivial field assignment or getters unless they protect an API, ownership, or
  ABI invariant.
- Add integration coverage for CRUD, partition keys, query pagination, continuation tokens,
  session consistency, retries, and failover as those operations become available.
- Keep examples in `example_test.go`, make them compile, and avoid credentials or customer data.
- Emulator tests use the repository's Cosmos emulator pipeline and `EMULATOR=true`; do not make
  ordinary unit tests depend on an emulator or live account.

## Validation

Run checks from `sdk/data/azcosmos` and start with the narrowest command that covers the change:

```bash
gofmt -w <changed-go-files>
go test ./...
go test -race ./...
go vet ./...
```

`valuelist/` is a nested Go module and is not covered by those commands. When it changes, validate
it separately:

```bash
cd valuelist
go test ./...
go vet ./...
```

For documentation-only changes, run the repository spelling check when available. Stage newly
created files first because the script discovers files through `git diff`:

```bash
pwsh ../../../eng/common/scripts/check-spelling-in-changed-files.ps1 \
  -TargetCommittish "<target-committish>" -ExitWithError
```

## Documentation and Comments

- Keep inline comments to one or two lines unless they capture a non-obvious ABI, ownership, or
  concurrency invariant.
- Explain why code exists, not what the code already says.
- Put multi-paragraph rationale in package documentation or a decision record and link to it.
- Never include real or customer-derived account keys, tokens, item bodies, identifiers, or
  personally identifiable information in logs, diagnostics, examples, or test fixtures. Clearly
  synthetic values are appropriate for tests.

## Additional Resources

- [Azure Cosmos DB REST API](https://learn.microsoft.com/rest/api/cosmos-db/)
- [Azure Cosmos DB partitioning](https://learn.microsoft.com/azure/cosmos-db/partitioning-overview)
- [Azure SDK Design Guidelines for Go](https://azure.github.io/azure-sdk/golang_introduction.html)
- [Azure SDK for Go contributing guide](https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md)
