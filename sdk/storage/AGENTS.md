# Storage SDK AGENTS.md

This file provides guidance for AI agents working in `sdk/storage/` of the Azure SDK for Go repository.

For repository-wide rules, see the root [AGENTS.md](../../AGENTS.md). For Copilot-specific guidance, see [`.github/copilot-instructions.md`](../../.github/copilot-instructions.md).

## Folder Overview

The `sdk/storage` folder contains the Go data-plane SDKs for Azure Storage:

- [Azure Blob Storage (`azblob`)](./azblob/README.md) - object storage for containers and blobs
- [Azure Data Lake Storage Gen2 (`azdatalake`)](./azdatalake/README.md) - hierarchical namespace storage built on Blob Storage
- [Azure Files (`azfile`)](./azfile/README.md) - SMB/NFS file shares
- [Azure Queue Storage (`azqueue`)](./azqueue/README.md) - message queues

Each top-level storage package under `sdk/storage/` is its own Go module with its own `go.mod`. Subpackages such as `blockblob`, `filesystem`, `directory`, or `share` live inside those modules and do not have separate `go.mod` files.

## Shared Client Hierarchy Pattern

Most storage SDKs follow a consistent resource hierarchy:

1. **Service client** - scoped to the storage account endpoint
2. **Intermediate container/share/filesystem client** - scoped to a named child resource
3. **Leaf resource client** - scoped to a blob, file, or directory/path

Examples:

- `azblob`: `Client`/`service.Client` -> `container.Client` -> `blob.Client` / `blockblob.Client` / `appendblob.Client` / `pageblob.Client`
- `azdatalake`: `service.Client` -> `filesystem.Client` -> `directory.Client` / `file.Client`
- `azfile`: `service.Client` -> `share.Client` -> `directory.Client` -> `file.Client`
- `azqueue`: `ServiceClient` -> `QueueClient`

When changing client creation flows, preserve this pattern and prefer adding behavior at the most specific layer that owns the operation.

## Authentication Patterns

Shared auth options across storage modules:

- **Shared key**: account name + account key via `NewSharedKeyCredential(...)`
- **SAS**: resource URLs with signed query parameters, with helpers in each package's `sas/` subpackage
- **Azure AD**: `azcore.TokenCredential`, typically from `sdk/azidentity`
- **Connection string**: convenience constructors such as `NewClientFromConnectionString(...)`

Storage packages commonly expose three constructor families:

- `NewClient(...)` or equivalent for Azure AD
- `NewClientWithSharedKeyCredential(...)`
- `NewClientWithNoCredential(...)` for anonymous/SAS URLs

Package-specific auth caveats belong in the package-level `AGENTS.md` files.

## SAS Token Generation

Each storage module has a `sas/` subpackage with resource-specific signature types and query parameter helpers.

Typical flow:

1. Build a package-specific signature value struct (for example `sas.BlobSignatureValues`)
2. Sign it with a shared key or user delegation credential
3. Append the returned SAS query parameters to the resource URL

Some service/resource clients also expose `GetSASURL(...)` helpers. Prefer existing helpers over manual string concatenation when they already exist.

## Generated Code

Do not manually edit generated clients unless the task is explicitly about regeneration.

Generated code patterns in storage packages:

- `azblob/internal/generated/**` is generated and has `//go:generate autorest ./autorest.md` entrypoints
- `azdatalake/internal/generated/**` and `azdatalake/internal/generated_blob/**` are generated and bridge DFS and Blob endpoints
- `azfile/internal/generated/**` is generated from service specs; this module also includes [`tsp-location.yaml`](./azfile/tsp-location.yaml)
- `azqueue/internal/generated/**` is generated from service specs; this module also includes [`tsp-location.yaml`](./azqueue/tsp-location.yaml)

In all four storage modules, prefer changing handwritten wrapper layers first (for example `client.go`, `service/`, `filesystem/`, `share/`, `blob/`, `file/`, `directory/`) and defer regeneration workflow details to the root [AGENTS.md](../../AGENTS.md).

## Service Version

Each storage module pins its service version in a hand-written `internal/generated/constants.go`:

```go
const ServiceVersion = "2026-06-06"
```

This file has no `zz_` prefix and is not overwritten by regeneration despite living under `internal/generated/`. It is manually updated when bumping the API version.

The `x-ms-version` header sent on every request is derived from this constant (each generated client is initialized with `version: ServiceVersion`). These modules do not currently wire up an azcore `APIVersion` policy, so the header is **not** overridable per client via `ClientOptions`; changing the service version means editing the constant (and regenerating).

### SAS version coupling

The `sv=` parameter in SAS tokens defaults to the same `ServiceVersion` constant, but through a separate code path:

- `sas/query_params.go` -> `var Version = generated.ServiceVersion`
- `sas/service.go` -> `if v.Version == "" { v.Version = Version }`

`sas.Version` is a package-level `var` (not a `const`), so it can be reassigned globally, but this is not per-client. The header path and the SAS path both start from `ServiceVersion` but are otherwise independent: a task that changes the effective service version must address both.

## Common Subpackage Patterns

Across storage modules, watch for these recurring package roles:

- `internal/` - internal-only helpers, generated layers, shared pipeline code, and test helpers; do not suggest external consumers import these packages
- `sas/` - SAS construction and parsing helpers
- error packages such as `bloberror`, `datalakeerror`, `fileerror`, `queueerror` - typed service error codes and helpers
- `lease/` - lease management helpers for resources that support leases

## Test Layout and Setup

Storage modules use the shared Azure SDK recording framework in `sdk/internal/recording`.

Common test assets:

- `test-resources.json` - ARM template describing live test resources, account variants, and role assignments
- `assets.json` - pins the Azure SDK assets repository tag used for recorded test assets
- `testdata/` - recorded session files consumed by playback tests
- `internal/testcommon/` - package-specific helpers for test accounts, sanitizers, and client construction

### Record modes

Tests default to playback when `AZURE_RECORD_MODE` is unset.

Supported modes from `sdk/internal/recording/recording.go`:

- `AZURE_RECORD_MODE=playback`
- `AZURE_RECORD_MODE=record`
- `AZURE_RECORD_MODE=live`

Recorded suites typically start the test proxy in `SetupSuite()` and call `recording.Start()`/`recording.Stop()` per test. When updating tests, keep sanitizers and proxy setup aligned with existing `internal/testcommon` helpers.

### Environment variables

The most common shared-key live test variables are:

- `AZURE_STORAGE_ACCOUNT_NAME`
- `AZURE_STORAGE_ACCOUNT_KEY`

Storage tests also use prefixed variants for specialized accounts, including:

- `SECONDARY_AZURE_STORAGE_ACCOUNT_NAME` / `SECONDARY_AZURE_STORAGE_ACCOUNT_KEY`
- `PREMIUM_AZURE_STORAGE_ACCOUNT_NAME` / `PREMIUM_AZURE_STORAGE_ACCOUNT_KEY`
- `SOFT_DELETE_AZURE_STORAGE_ACCOUNT_NAME` / `SOFT_DELETE_AZURE_STORAGE_ACCOUNT_KEY`
- `DATALAKE_AZURE_STORAGE_ACCOUNT_NAME` / `DATALAKE_AZURE_STORAGE_ACCOUNT_KEY`
- `IMMUTABLE_AZURE_STORAGE_ACCOUNT_NAME` / `IMMUTABLE_AZURE_STORAGE_ACCOUNT_KEY`

AAD-based tests typically rely on the usual `azidentity` environment variables such as `AZURE_TENANT_ID`, `AZURE_CLIENT_ID`, and `AZURE_CLIENT_SECRET`.

Playback usually does not require real storage credentials because test helpers substitute fake account values.

## Live vs. Playback Testing

Typical targeted commands:

```bash
cd sdk/storage/azblob && go test ./...
cd sdk/storage/azdatalake && go test ./...
cd sdk/storage/azfile && go test ./...
cd sdk/storage/azqueue && go test ./...
```

Guidance:

- Use **playback** for most validation of test changes
- Use **record** when you intentionally update recordings
- Use **live** for scenarios that are marked live-only or when validating behavior against actual service state
- For record/live, ensure the package's required storage env vars are present before running tests

## How to Work Safely in Storage Modules

1. Read the package-level `AGENTS.md` before editing `azblob`, `azdatalake`, or `azfile`
2. Check whether the behavior lives in a handwritten wrapper or a generated layer
3. Preserve the existing client hierarchy and constructor patterns
4. Reuse `internal/testcommon` helpers instead of inventing ad hoc test setup
5. Keep SAS, lease, and error-handling changes in their dedicated subpackages when possible

## Note on `azqueue`

`azqueue` does not have its own package-level `AGENTS.md` because of its relatively flat structure. Agents working in `azqueue` should use this shared storage `AGENTS.md` for guidance. The client hierarchy is simply `ServiceClient` -> `QueueClient`, with generated code in `internal/generated/` following the same patterns as the other storage modules (including a hand-written `internal/generated/constants.go` `ServiceVersion` constant).
