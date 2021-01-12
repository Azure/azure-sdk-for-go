Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9/specification/netapp/resource-manager/readme.md tag: `package-netapp-2020-09-01`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *AccountBackupsDeleteFuture.Result(AccountBackupsClient) (autorest.Response, error)
1. *AccountsCreateOrUpdateFuture.Result(AccountsClient) (Account, error)
1. *AccountsDeleteFuture.Result(AccountsClient) (autorest.Response, error)
1. *AccountsUpdateFuture.Result(AccountsClient) (Account, error)
1. *BackupPoliciesCreateFuture.Result(BackupPoliciesClient) (BackupPolicy, error)
1. *BackupPoliciesDeleteFuture.Result(BackupPoliciesClient) (autorest.Response, error)
1. *BackupsCreateFuture.Result(BackupsClient) (Backup, error)
1. *BackupsDeleteFuture.Result(BackupsClient) (autorest.Response, error)
1. *PoolsCreateOrUpdateFuture.Result(PoolsClient) (CapacityPool, error)
1. *PoolsDeleteFuture.Result(PoolsClient) (autorest.Response, error)
1. *PoolsUpdateFuture.Result(PoolsClient) (CapacityPool, error)
1. *SnapshotPoliciesDeleteFuture.Result(SnapshotPoliciesClient) (autorest.Response, error)
1. *SnapshotsCreateFuture.Result(SnapshotsClient) (Snapshot, error)
1. *SnapshotsDeleteFuture.Result(SnapshotsClient) (autorest.Response, error)
1. *SnapshotsUpdateFuture.Result(SnapshotsClient) (Snapshot, error)
1. *VolumesAuthorizeReplicationFuture.Result(VolumesClient) (autorest.Response, error)
1. *VolumesBreakReplicationFuture.Result(VolumesClient) (autorest.Response, error)
1. *VolumesCreateOrUpdateFuture.Result(VolumesClient) (Volume, error)
1. *VolumesDeleteFuture.Result(VolumesClient) (autorest.Response, error)
1. *VolumesDeleteReplicationFuture.Result(VolumesClient) (autorest.Response, error)
1. *VolumesPoolChangeFuture.Result(VolumesClient) (autorest.Response, error)
1. *VolumesReInitializeReplicationFuture.Result(VolumesClient) (autorest.Response, error)
1. *VolumesResyncReplicationFuture.Result(VolumesClient) (autorest.Response, error)
1. *VolumesRevertFuture.Result(VolumesClient) (autorest.Response, error)
1. *VolumesUpdateFuture.Result(VolumesClient) (Volume, error)

## Struct Changes

### Removed Struct Fields

1. AccountBackupsDeleteFuture.azure.Future
1. AccountsCreateOrUpdateFuture.azure.Future
1. AccountsDeleteFuture.azure.Future
1. AccountsUpdateFuture.azure.Future
1. BackupPoliciesCreateFuture.azure.Future
1. BackupPoliciesDeleteFuture.azure.Future
1. BackupsCreateFuture.azure.Future
1. BackupsDeleteFuture.azure.Future
1. PoolsCreateOrUpdateFuture.azure.Future
1. PoolsDeleteFuture.azure.Future
1. PoolsUpdateFuture.azure.Future
1. SnapshotPoliciesDeleteFuture.azure.Future
1. SnapshotsCreateFuture.azure.Future
1. SnapshotsDeleteFuture.azure.Future
1. SnapshotsUpdateFuture.azure.Future
1. VolumesAuthorizeReplicationFuture.azure.Future
1. VolumesBreakReplicationFuture.azure.Future
1. VolumesCreateOrUpdateFuture.azure.Future
1. VolumesDeleteFuture.azure.Future
1. VolumesDeleteReplicationFuture.azure.Future
1. VolumesPoolChangeFuture.azure.Future
1. VolumesReInitializeReplicationFuture.azure.Future
1. VolumesResyncReplicationFuture.azure.Future
1. VolumesRevertFuture.azure.Future
1. VolumesUpdateFuture.azure.Future

## Struct Changes

### New Struct Fields

1. AccountBackupsDeleteFuture.Result
1. AccountBackupsDeleteFuture.azure.FutureAPI
1. AccountsCreateOrUpdateFuture.Result
1. AccountsCreateOrUpdateFuture.azure.FutureAPI
1. AccountsDeleteFuture.Result
1. AccountsDeleteFuture.azure.FutureAPI
1. AccountsUpdateFuture.Result
1. AccountsUpdateFuture.azure.FutureAPI
1. BackupPoliciesCreateFuture.Result
1. BackupPoliciesCreateFuture.azure.FutureAPI
1. BackupPoliciesDeleteFuture.Result
1. BackupPoliciesDeleteFuture.azure.FutureAPI
1. BackupsCreateFuture.Result
1. BackupsCreateFuture.azure.FutureAPI
1. BackupsDeleteFuture.Result
1. BackupsDeleteFuture.azure.FutureAPI
1. PoolsCreateOrUpdateFuture.Result
1. PoolsCreateOrUpdateFuture.azure.FutureAPI
1. PoolsDeleteFuture.Result
1. PoolsDeleteFuture.azure.FutureAPI
1. PoolsUpdateFuture.Result
1. PoolsUpdateFuture.azure.FutureAPI
1. SnapshotPoliciesDeleteFuture.Result
1. SnapshotPoliciesDeleteFuture.azure.FutureAPI
1. SnapshotsCreateFuture.Result
1. SnapshotsCreateFuture.azure.FutureAPI
1. SnapshotsDeleteFuture.Result
1. SnapshotsDeleteFuture.azure.FutureAPI
1. SnapshotsUpdateFuture.Result
1. SnapshotsUpdateFuture.azure.FutureAPI
1. VolumesAuthorizeReplicationFuture.Result
1. VolumesAuthorizeReplicationFuture.azure.FutureAPI
1. VolumesBreakReplicationFuture.Result
1. VolumesBreakReplicationFuture.azure.FutureAPI
1. VolumesCreateOrUpdateFuture.Result
1. VolumesCreateOrUpdateFuture.azure.FutureAPI
1. VolumesDeleteFuture.Result
1. VolumesDeleteFuture.azure.FutureAPI
1. VolumesDeleteReplicationFuture.Result
1. VolumesDeleteReplicationFuture.azure.FutureAPI
1. VolumesPoolChangeFuture.Result
1. VolumesPoolChangeFuture.azure.FutureAPI
1. VolumesReInitializeReplicationFuture.Result
1. VolumesReInitializeReplicationFuture.azure.FutureAPI
1. VolumesResyncReplicationFuture.Result
1. VolumesResyncReplicationFuture.azure.FutureAPI
1. VolumesRevertFuture.Result
1. VolumesRevertFuture.azure.FutureAPI
1. VolumesUpdateFuture.Result
1. VolumesUpdateFuture.azure.FutureAPI
