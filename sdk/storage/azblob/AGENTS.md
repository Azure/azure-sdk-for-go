# Azure Blob Storage AGENTS.md

This file provides guidance for AI agents working in `sdk/storage/azblob`.

Start with the shared [storage AGENTS.md](../AGENTS.md) and the repository root [AGENTS.md](../../../AGENTS.md).

## Package Purpose

`azblob` is the Go SDK for Azure Blob Storage, Azure's object storage service for unstructured data. Common use cases in this module include:

- general file/object storage
- browser or service download scenarios
- large file upload/download flows
- archival, backup, and media content
- append-style logging and page-backed disk/VHD scenarios

Key docs:

- [README.md](./README.md)
- [doc.go](./doc.go)
- [examples_test.go](./examples_test.go)
- [client.go](./client.go)

## Client Hierarchy

`azblob` has both a convenience top-level client and specialized subpackage clients.

Primary hierarchy:

1. `azblob.Client` or `service.Client` - account-scoped
2. `container.Client` - container-scoped
3. `blob.Client` / `blockblob.Client` / `appendblob.Client` / `pageblob.Client` - blob-scoped

Common creation flow:

- `client.ServiceClient().NewContainerClient(name)` from the top-level `azblob.Client`
- `containerClient.NewBlobClient(name)` for generic blob operations
- `containerClient.NewBlockBlobClient(name)`, `NewAppendBlobClient(name)`, or `NewPageBlobClient(name)` for blob-type-specific operations

The top-level `azblob.Client` in [`client.go`](./client.go) is mostly a convenience wrapper that forwards to `service.Client` and creates lower-level clients on demand.

## Important Subpackages

- [`blob/`](./blob/) - operations common to all blob types, including metadata, tags, properties, delete/undelete, and downloads
- [`blockblob/`](./blockblob/) - block blob uploads, staging/committing blocks, and high-level file/stream transfer helpers
- [`appendblob/`](./appendblob/) - append-only blobs used for logging-style workloads
- [`pageblob/`](./pageblob/) - page-oriented blobs used for random-write/page workloads such as VHD scenarios
- [`container/`](./container/) - container creation, deletion, metadata, access policy, and blob listing
- [`service/`](./service/) - service-level operations such as listing containers and account/service settings
- [`lease/`](./lease/) - blob and container lease helpers for concurrency and ownership control
- [`sas/`](./sas/) - service SAS and user delegation SAS helpers
- [`bloberror/`](./bloberror/) - typed blob service error codes
- `internal/` - generated layers, shared helpers, test helpers, and pipeline internals

## Blob Type Guidance

Choose the blob type that matches the scenario:

- **Block blobs** - default choice for most uploads and downloads, especially general files and large payloads
- **Append blobs** - optimized for append-only writes such as log accumulation
- **Page blobs** - optimized for page-aligned random access such as disks and VHD-like workloads

If a change only affects generic blob behavior, prefer `blob/`. If it depends on block staging, append semantics, or page ranges, change the specialized package instead.

## Large File Transfer Helpers

For high-level transfer behavior, start in [`blockblob/client.go`](./blockblob/client.go).

Important helpers:

- `UploadFile(...)`
- `UploadStream(...)`
- `DownloadFile(...)`

These helpers already implement chunked/parallel transfer behavior. Prefer using or extending them instead of re-implementing block management in callers.

## Lease Pattern

The [`lease/`](./lease/) package contains dedicated clients for blob and container leases.

Common operations:

- acquire
- renew
- change
- release
- break

Use leases when a change involves concurrency control, exclusive writers, or coordination between clients. Blob leases and container leases are separate client types.

## Common Usage Patterns

### Get a container client

Start from `azblob.Client` or `service.Client`, then create a `container.Client` with `NewContainerClient(name)`.

### List blobs

Use `container.Client.NewListBlobsFlatPager(...)` for flat listing, or the corresponding hierarchy-aware pager when needed.

### Set metadata or tags

- Blob metadata/tags usually live on `blob.Client` or specialized blob clients such as `blockblob.Client`
- Container metadata lives on `container.Client`

If a method exists on both a specialized client and the generic `blob.Client`, prefer following the existing package pattern instead of duplicating logic.

## Testing Notes

Key live-test variables commonly used here:

- `AZURE_STORAGE_ACCOUNT_NAME`
- `AZURE_STORAGE_ACCOUNT_KEY`
- `SECONDARY_AZURE_STORAGE_ACCOUNT_NAME` / `SECONDARY_AZURE_STORAGE_ACCOUNT_KEY`
- `PREMIUM_AZURE_STORAGE_ACCOUNT_NAME` / `PREMIUM_AZURE_STORAGE_ACCOUNT_KEY`
- `SOFT_DELETE_AZURE_STORAGE_ACCOUNT_NAME` / `SOFT_DELETE_AZURE_STORAGE_ACCOUNT_KEY`
- `DATALAKE_AZURE_STORAGE_ACCOUNT_NAME` / `DATALAKE_AZURE_STORAGE_ACCOUNT_KEY`
- `IMMUTABLE_AZURE_STORAGE_ACCOUNT_NAME` / `IMMUTABLE_AZURE_STORAGE_ACCOUNT_KEY`

Test helpers live under [`internal/testcommon/`](./internal/testcommon/). Recorded tests use `sdk/internal/recording`, start the proxy in suite setup, and sanitize storage URLs plus auth headers before recording.

Targeted validation command:

```bash
cd sdk/storage/azblob && go test ./...
```

## Editing Guidance

- Prefer handwritten wrappers in `client.go`, `service/`, `container/`, `blob/`, or specialized blob packages over `internal/generated/`
- Keep the convenience-client forwarding pattern intact when changing `azblob.Client`
- Reuse existing pagers, transfer helpers, and lease clients instead of adding parallel abstractions
- Keep blob-type-specific behavior in the correct specialized package
