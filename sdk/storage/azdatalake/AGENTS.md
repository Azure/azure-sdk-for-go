# Azure Data Lake Storage Gen2 AGENTS.md

This file provides guidance for AI agents working in `sdk/storage/azdatalake`.

Start with the shared [storage AGENTS.md](../AGENTS.md) and the repository root [AGENTS.md](../../../AGENTS.md).

## Package Purpose

`azdatalake` is the Go SDK for Azure Data Lake Storage Gen2 (ADLS Gen2). It exposes hierarchical namespace (HNS) storage semantics on top of Azure Blob Storage.

Key characteristics:

- true directory and file paths
- POSIX-style ACLs and permissions
- path-based APIs and path rename/move operations
- Blob-backed storage with both DFS and Blob endpoints involved internally

Key docs:

- [README.md](./README.md)
- [doc.go](./doc.go)
- [`service/client.go`](./service/client.go)
- [`filesystem/client.go`](./filesystem/client.go)
- [`directory/client.go`](./directory/client.go)
- [`file/client.go`](./file/client.go)

## Client Hierarchy

Primary hierarchy:

1. `service.Client` - account-scoped
2. `filesystem.Client` - filesystem-scoped
3. `directory.Client` / `file.Client` - path-scoped

Unlike `azblob`, this module does not center on a single top-level convenience client in the root package. The main entry points are the `service`, `filesystem`, `directory`, and `file` subpackages.

Common creation flow:

- `serviceClient.NewFileSystemClient(name)`
- `filesystemClient.NewDirectoryClient(path)`
- `filesystemClient.NewFileClient(path)`
- `directoryClient.NewSubdirectoryClient(name)` or `directoryClient.NewFileClient(name)`

## Important Subpackages

- [`service/`](./service/) - account-level operations and filesystem creation/listing
- [`filesystem/`](./filesystem/) - filesystem operations plus creation of file/directory clients
- [`directory/`](./directory/) - directory CRUD, rename, ACLs, recursive ACL updates, and child path creation
- [`file/`](./file/) - file CRUD, append/flush/upload/download helpers, rename, and ACLs
- [`lease/`](./lease/) - filesystem and path lease helpers
- [`sas/`](./sas/) - SAS construction/parsing
- [`datalakeerror/`](./datalakeerror/) - typed storage error codes
- `internal/` - DFS-specific generated code and helpers
- `internal/generated_blob/` - Blob-endpoint generated layer used because the SDK bridges DFS and Blob behaviors

## Differences from `azblob`

Important conceptual differences from Blob Storage:

- HNS gives you **real directories**, not just blob name prefixes
- directory/file **rename and move** are first-class path operations
- POSIX-like access control is part of the API surface
- many clients track both **DFSURL** and **BlobURL** internally

When changing path behavior, be careful not to break the dual-endpoint design. Service/filesystem/path clients often call both DFS and Blob layers.

## Common Usage Patterns

### Create a filesystem

Start with `service.Client`, then call `NewFileSystemClient(name)` and `Create(...)`.

### Upload or download files

The main high-level helpers are on `file.Client`, including:

- `UploadFile(...)`
- `UploadStream(...)`
- `DownloadFile(...)`

### Work with directories

Use `filesystem.Client.NewDirectoryClient(path)` or `directory.Client.NewSubdirectoryClient(name)` for hierarchical operations.

### Set ACLs or permissions

Key methods include:

- `directory.Client.SetAccessControl(...)`
- `directory.Client.SetAccessControlRecursive(...)`
- `directory.Client.GetAccessControl(...)`
- `file.Client.SetAccessControl(...)`
- `file.Client.GetAccessControl(...)`

If a task involves permissions, check both file and directory variants because the recursive ACL surface exists only on directories.

### Rename or move paths

Both `directory.Client` and `file.Client` have `Rename(...)`. These operations are a major reason to use ADLS Gen2 instead of plain blob prefixes.

## Testing Notes

Key live-test variables commonly used here:

- `AZURE_STORAGE_ACCOUNT_NAME`
- `AZURE_STORAGE_ACCOUNT_KEY`
- `DATALAKE_AZURE_STORAGE_ACCOUNT_NAME`
- `DATALAKE_AZURE_STORAGE_ACCOUNT_KEY`
- `DATALAKE_AZURE_STORAGE_ENCRYPTION_SCOPE`
- `SECONDARY_AZURE_STORAGE_ACCOUNT_NAME` / `SECONDARY_AZURE_STORAGE_ACCOUNT_KEY`
- `PREMIUM_AZURE_STORAGE_ACCOUNT_NAME` / `PREMIUM_AZURE_STORAGE_ACCOUNT_KEY`
- `SOFT_DELETE_AZURE_STORAGE_ACCOUNT_NAME` / `SOFT_DELETE_AZURE_STORAGE_ACCOUNT_KEY`

Test helpers live under [`internal/testcommon/`](./internal/testcommon/). Recorded tests sanitize both Blob and DFS URLs because this module talks to both endpoint shapes.

Targeted validation command:

```bash
cd sdk/storage/azdatalake && go test ./...
```

## Editing Guidance

- Preserve the DFS/Blob dual-layer behavior in `service`, `filesystem`, `directory`, and `file`
- Prefer handwritten wrappers over `internal/generated/` and `internal/generated_blob/`
- Keep ACL logic in the path-level clients that already own it
- Use existing rename and transfer helpers instead of duplicating path or upload logic
