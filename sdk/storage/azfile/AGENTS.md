# Azure Files AGENTS.md

This file provides guidance for AI agents working in `sdk/storage/azfile`.

Start with the shared [storage AGENTS.md](../AGENTS.md) and the repository root [AGENTS.md](../../../AGENTS.md).

## Package Purpose

`azfile` is the Go SDK for Azure Files, Azure's managed file share service for SMB and NFS workloads.

Common scenarios:

- lift-and-shift applications that expect shared file storage
- shared application/configuration files
- directory and file workflows that resemble network file shares
- SMB handle management and share snapshots

Key docs:

- [README.md](./README.md)
- [doc.go](./doc.go)
- [`service/client.go`](./service/client.go)
- [`share/client.go`](./share/client.go)
- [`directory/client.go`](./directory/client.go)
- [`file/client.go`](./file/client.go)
- [`tsp-location.yaml`](./tsp-location.yaml)

## Client Hierarchy

Primary hierarchy:

1. `service.Client` - account-scoped
2. `share.Client` - share-scoped
3. `directory.Client` - directory-scoped
4. `file.Client` - file-scoped

Common creation flow:

- `serviceClient.NewShareClient(name)`
- `shareClient.NewDirectoryClient(path)`
- `shareClient.NewRootDirectoryClient()` for root-level traversal
- `directoryClient.NewSubdirectoryClient(name)`
- `directoryClient.NewFileClient(name)`

## Important Subpackages

- [`service/`](./service/) - service properties, share creation/listing, restore, user delegation keys
- [`share/`](./share/) - share CRUD, snapshots, metadata, ACL/policy, root directory access, permissions
- [`directory/`](./directory/) - directory CRUD, rename, metadata, file/directory listing, and handle operations
- [`file/`](./file/) - file CRUD, upload/download helpers, range operations, rename, and handle operations
- [`lease/`](./lease/) - dedicated lease clients for shares and files
- [`sas/`](./sas/) - SAS helpers
- [`fileerror/`](./fileerror/) - typed file service error codes
- `internal/` - generated layer, shared pipeline code, and test helpers

## Azure Files Concepts That Matter Here

### Shares, directories, and files

This module models a traditional file share hierarchy. Directory deletion requires the directory to be empty; file uploads often involve explicit create + range upload semantics under the hood.

### Leases

The [`lease/`](./lease/) package provides wrappers around share and file lease operations. The public methods are named `Acquire`, `Break`, `Change`, `Release`, and `Renew` on dedicated lease client types.

## Authentication Caveats

Azure Files has stricter token-auth behavior than the other storage modules.

Important details from the current clients:

- `service.Client` and `share.Client` support token credential construction mainly to enable lower-level clients, but not every service/share operation supports token auth
- `ClientOptions.FileRequestIntent` is currently required for token-authenticated Azure Files clients
- shared key, SAS, and connection string flows are the most straightforward paths for broad API coverage

When changing auth or constructors, keep these caveats intact.

## Common Usage Patterns

### Create a share

Start with `service.Client`, create a `share.Client`, then call `Create(...)`.

### Upload or download files

Common file helpers live on `file.Client`, including `UploadFile(...)` and `DownloadFile(...)`.

### Traverse directories

Use `share.Client.NewRootDirectoryClient()` and `directory.Client.NewListFilesAndDirectoriesPager(...)` for tree walking.

### Rename paths

Directory and file rename logic is implemented on the corresponding clients. If a change touches rename behavior, also check handle interactions and trailing-dot options in client options.

## Testing Notes

Key live-test variables commonly used here:

- `AZURE_STORAGE_ACCOUNT_NAME`
- `AZURE_STORAGE_ACCOUNT_KEY`
- `SECONDARY_AZURE_STORAGE_ACCOUNT_NAME` / `SECONDARY_AZURE_STORAGE_ACCOUNT_KEY`
- `SOFT_DELETE_AZURE_STORAGE_ACCOUNT_NAME` / `SOFT_DELETE_AZURE_STORAGE_ACCOUNT_KEY`
- `PREMIUM_FILE_STORAGE_ACCOUNT_NAME` / `PREMIUM_FILE_STORAGE_ACCOUNT_KEY`
- `AZURE_STORAGE_ENCRYPTION_SCOPE`

Recorded tests use the `sdk/storage/azfile/testdata` recording directory configured in `internal/testcommon`, and helpers live under [`internal/testcommon/`](./internal/testcommon/). Tests sanitize file endpoint URLs and rename-source headers before recording.

Targeted validation command:

```bash
cd sdk/storage/azfile && go test ./...
```

## Generated Code and Editing Guidance

This module includes [`tsp-location.yaml`](./tsp-location.yaml) and generated code under `internal/generated/`. Do not hand-edit generated files unless the task is explicitly about regeneration.

When making functional changes, prefer the handwritten wrappers in:

- `service/`
- `share/`
- `directory/`
- `file/`
- `lease/`

Keep token-auth constraints, range/rename helpers, and handle-management behavior aligned with the existing wrappers.
