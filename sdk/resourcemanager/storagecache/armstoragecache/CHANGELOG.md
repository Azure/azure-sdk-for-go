# Release History

## 4.1.0 (2025-09-26)
### Features Added

- New enum type `AutoExportJobAdminStatus` with values `AutoExportJobAdminStatusDisable`, `AutoExportJobAdminStatusEnable`
- New enum type `AutoExportJobProvisioningStateType` with values `AutoExportJobProvisioningStateTypeCanceled`, `AutoExportJobProvisioningStateTypeCreating`, `AutoExportJobProvisioningStateTypeDeleting`, `AutoExportJobProvisioningStateTypeFailed`, `AutoExportJobProvisioningStateTypeSucceeded`, `AutoExportJobProvisioningStateTypeUpdating`
- New enum type `AutoExportStatusType` with values `AutoExportStatusTypeDisableFailed`, `AutoExportStatusTypeDisabled`, `AutoExportStatusTypeDisabling`, `AutoExportStatusTypeFailed`, `AutoExportStatusTypeInProgress`
- New enum type `AutoImportJobPropertiesAdminStatus` with values `AutoImportJobPropertiesAdminStatusDisable`, `AutoImportJobPropertiesAdminStatusEnable`
- New enum type `AutoImportJobPropertiesProvisioningState` with values `AutoImportJobPropertiesProvisioningStateCanceled`, `AutoImportJobPropertiesProvisioningStateCreating`, `AutoImportJobPropertiesProvisioningStateDeleting`, `AutoImportJobPropertiesProvisioningStateFailed`, `AutoImportJobPropertiesProvisioningStateSucceeded`, `AutoImportJobPropertiesProvisioningStateUpdating`
- New enum type `AutoImportJobState` with values `AutoImportJobStateDisabled`, `AutoImportJobStateDisabling`, `AutoImportJobStateFailed`, `AutoImportJobStateInProgress`
- New enum type `AutoImportJobUpdatePropertiesAdminStatus` with values `AutoImportJobUpdatePropertiesAdminStatusDisable`, `AutoImportJobUpdatePropertiesAdminStatusEnable`
- New enum type `ImportJobAdminStatus` with values `ImportJobAdminStatusActive`, `ImportJobAdminStatusCancel`
- New function `NewAutoExportJobsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AutoExportJobsClient, error)`
- New function `*AutoExportJobsClient.BeginCreateOrUpdate(context.Context, string, string, string, AutoExportJob, *AutoExportJobsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AutoExportJobsClientCreateOrUpdateResponse], error)`
- New function `*AutoExportJobsClient.BeginDelete(context.Context, string, string, string, *AutoExportJobsClientBeginDeleteOptions) (*runtime.Poller[AutoExportJobsClientDeleteResponse], error)`
- New function `*AutoExportJobsClient.Get(context.Context, string, string, string, *AutoExportJobsClientGetOptions) (AutoExportJobsClientGetResponse, error)`
- New function `*AutoExportJobsClient.NewListByAmlFilesystemPager(string, string, *AutoExportJobsClientListByAmlFilesystemOptions) *runtime.Pager[AutoExportJobsClientListByAmlFilesystemResponse]`
- New function `*AutoExportJobsClient.BeginUpdate(context.Context, string, string, string, AutoExportJobUpdate, *AutoExportJobsClientBeginUpdateOptions) (*runtime.Poller[AutoExportJobsClientUpdateResponse], error)`
- New function `NewAutoImportJobsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AutoImportJobsClient, error)`
- New function `*AutoImportJobsClient.BeginCreateOrUpdate(context.Context, string, string, string, AutoImportJob, *AutoImportJobsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AutoImportJobsClientCreateOrUpdateResponse], error)`
- New function `*AutoImportJobsClient.BeginDelete(context.Context, string, string, string, *AutoImportJobsClientBeginDeleteOptions) (*runtime.Poller[AutoImportJobsClientDeleteResponse], error)`
- New function `*AutoImportJobsClient.Get(context.Context, string, string, string, *AutoImportJobsClientGetOptions) (AutoImportJobsClientGetResponse, error)`
- New function `*AutoImportJobsClient.NewListByAmlFilesystemPager(string, string, *AutoImportJobsClientListByAmlFilesystemOptions) *runtime.Pager[AutoImportJobsClientListByAmlFilesystemResponse]`
- New function `*AutoImportJobsClient.BeginUpdate(context.Context, string, string, string, AutoImportJobUpdate, *AutoImportJobsClientBeginUpdateOptions) (*runtime.Poller[AutoImportJobsClientUpdateResponse], error)`
- New function `*ClientFactory.NewAutoExportJobsClient() *AutoExportJobsClient`
- New function `*ClientFactory.NewAutoImportJobsClient() *AutoImportJobsClient`
- New struct `AutoExportJob`
- New struct `AutoExportJobProperties`
- New struct `AutoExportJobPropertiesStatus`
- New struct `AutoExportJobUpdate`
- New struct `AutoExportJobUpdateProperties`
- New struct `AutoExportJobsListResult`
- New struct `AutoImportJob`
- New struct `AutoImportJobProperties`
- New struct `AutoImportJobPropertiesStatus`
- New struct `AutoImportJobPropertiesStatusBlobSyncEvents`
- New struct `AutoImportJobUpdate`
- New struct `AutoImportJobUpdateProperties`
- New struct `AutoImportJobsListResult`
- New struct `ImportJobUpdateProperties`
- New field `AdminStatus` in struct `ImportJobProperties`
- New field `ImportedDirectories`, `ImportedFiles`, `ImportedSymlinks`, `PreexistingDirectories`, `PreexistingFiles`, `PreexistingSymlinks` in struct `ImportJobPropertiesStatus`
- New field `Properties` in struct `ImportJobUpdate`


## 4.0.0 (2024-05-24)
### Breaking Changes

- Type of `AscOperation.Error` has been changed from `*ErrorResponse` to `*AscOperationErrorResponse`

### Features Added

- New enum type `AmlFilesystemSquashMode` with values `AmlFilesystemSquashModeAll`, `AmlFilesystemSquashModeNone`, `AmlFilesystemSquashModeRootOnly`
- New enum type `ConflictResolutionMode` with values `ConflictResolutionModeFail`, `ConflictResolutionModeOverwriteAlways`, `ConflictResolutionModeOverwriteIfDirty`, `ConflictResolutionModeSkip`
- New enum type `ImportJobProvisioningStateType` with values `ImportJobProvisioningStateTypeCanceled`, `ImportJobProvisioningStateTypeCreating`, `ImportJobProvisioningStateTypeDeleting`, `ImportJobProvisioningStateTypeFailed`, `ImportJobProvisioningStateTypeSucceeded`, `ImportJobProvisioningStateTypeUpdating`
- New enum type `ImportStatusType` with values `ImportStatusTypeCanceled`, `ImportStatusTypeCancelling`, `ImportStatusTypeCompleted`, `ImportStatusTypeCompletedPartial`, `ImportStatusTypeFailed`, `ImportStatusTypeInProgress`
- New function `*ClientFactory.NewImportJobsClient() *ImportJobsClient`
- New function `NewImportJobsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ImportJobsClient, error)`
- New function `*ImportJobsClient.BeginCreateOrUpdate(context.Context, string, string, string, ImportJob, *ImportJobsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ImportJobsClientCreateOrUpdateResponse], error)`
- New function `*ImportJobsClient.BeginDelete(context.Context, string, string, string, *ImportJobsClientBeginDeleteOptions) (*runtime.Poller[ImportJobsClientDeleteResponse], error)`
- New function `*ImportJobsClient.Get(context.Context, string, string, string, *ImportJobsClientGetOptions) (ImportJobsClientGetResponse, error)`
- New function `*ImportJobsClient.NewListByAmlFilesystemPager(string, string, *ImportJobsClientListByAmlFilesystemOptions) *runtime.Pager[ImportJobsClientListByAmlFilesystemResponse]`
- New function `*ImportJobsClient.BeginUpdate(context.Context, string, string, string, ImportJobUpdate, *ImportJobsClientBeginUpdateOptions) (*runtime.Poller[ImportJobsClientUpdateResponse], error)`
- New struct `AmlFilesystemRootSquashSettings`
- New struct `AscOperationErrorResponse`
- New struct `ImportJob`
- New struct `ImportJobProperties`
- New struct `ImportJobPropertiesStatus`
- New struct `ImportJobUpdate`
- New struct `ImportJobsListResult`
- New field `ImportPrefixesInitial` in struct `AmlFilesystemHsmSettings`
- New field `RootSquashSettings` in struct `AmlFilesystemProperties`
- New field `RootSquashSettings` in struct `AmlFilesystemUpdateProperties`


## 3.4.0-beta.1 (2024-02-23)
### Features Added

- New enum type `AmlFilesystemSquashMode` with values `AmlFilesystemSquashModeAll`, `AmlFilesystemSquashModeNone`, `AmlFilesystemSquashModeRootOnly`
- New struct `AmlFilesystemRootSquashSettings`
- New field `RootSquashSettings` in struct `AmlFilesystemProperties`
- New field `RootSquashSettings` in struct `AmlFilesystemUpdateProperties`


## 3.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 3.2.1 (2023-06-23)
### Bugs Fixed

- Change `ProvisioningStateTypeCancelled` value to `Canceled`


## 3.2.0 (2023-06-23)
### Features Added

- New enum type `AmlFilesystemHealthStateType` with values `AmlFilesystemHealthStateTypeAvailable`, `AmlFilesystemHealthStateTypeDegraded`, `AmlFilesystemHealthStateTypeMaintenance`, `AmlFilesystemHealthStateTypeTransitioning`, `AmlFilesystemHealthStateTypeUnavailable`
- New enum type `AmlFilesystemIdentityType` with values `AmlFilesystemIdentityTypeNone`, `AmlFilesystemIdentityTypeUserAssigned`
- New enum type `AmlFilesystemProvisioningStateType` with values `AmlFilesystemProvisioningStateTypeCanceled`, `AmlFilesystemProvisioningStateTypeCreating`, `AmlFilesystemProvisioningStateTypeDeleting`, `AmlFilesystemProvisioningStateTypeFailed`, `AmlFilesystemProvisioningStateTypeSucceeded`, `AmlFilesystemProvisioningStateTypeUpdating`
- New enum type `ArchiveStatusType` with values `ArchiveStatusTypeCanceled`, `ArchiveStatusTypeCancelling`, `ArchiveStatusTypeCompleted`, `ArchiveStatusTypeFSScanInProgress`, `ArchiveStatusTypeFailed`, `ArchiveStatusTypeIdle`, `ArchiveStatusTypeInProgress`, `ArchiveStatusTypeNotConfigured`
- New enum type `FilesystemSubnetStatusType` with values `FilesystemSubnetStatusTypeInvalid`, `FilesystemSubnetStatusTypeOk`
- New enum type `MaintenanceDayOfWeekType` with values `MaintenanceDayOfWeekTypeFriday`, `MaintenanceDayOfWeekTypeMonday`, `MaintenanceDayOfWeekTypeSaturday`, `MaintenanceDayOfWeekTypeSunday`, `MaintenanceDayOfWeekTypeThursday`, `MaintenanceDayOfWeekTypeTuesday`, `MaintenanceDayOfWeekTypeWednesday`
- New function `NewAmlFilesystemsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AmlFilesystemsClient, error)`
- New function `*AmlFilesystemsClient.Archive(context.Context, string, string, *AmlFilesystemsClientArchiveOptions) (AmlFilesystemsClientArchiveResponse, error)`
- New function `*AmlFilesystemsClient.CancelArchive(context.Context, string, string, *AmlFilesystemsClientCancelArchiveOptions) (AmlFilesystemsClientCancelArchiveResponse, error)`
- New function `*AmlFilesystemsClient.BeginCreateOrUpdate(context.Context, string, string, AmlFilesystem, *AmlFilesystemsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AmlFilesystemsClientCreateOrUpdateResponse], error)`
- New function `*AmlFilesystemsClient.BeginDelete(context.Context, string, string, *AmlFilesystemsClientBeginDeleteOptions) (*runtime.Poller[AmlFilesystemsClientDeleteResponse], error)`
- New function `*AmlFilesystemsClient.Get(context.Context, string, string, *AmlFilesystemsClientGetOptions) (AmlFilesystemsClientGetResponse, error)`
- New function `*AmlFilesystemsClient.NewListByResourceGroupPager(string, *AmlFilesystemsClientListByResourceGroupOptions) *runtime.Pager[AmlFilesystemsClientListByResourceGroupResponse]`
- New function `*AmlFilesystemsClient.NewListPager(*AmlFilesystemsClientListOptions) *runtime.Pager[AmlFilesystemsClientListResponse]`
- New function `*AmlFilesystemsClient.BeginUpdate(context.Context, string, string, AmlFilesystemUpdate, *AmlFilesystemsClientBeginUpdateOptions) (*runtime.Poller[AmlFilesystemsClientUpdateResponse], error)`
- New function `*ClientFactory.NewAmlFilesystemsClient() *AmlFilesystemsClient`
- New function `*ClientFactory.NewManagementClient() *ManagementClient`
- New function `NewManagementClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagementClient, error)`
- New function `*ManagementClient.CheckAmlFSSubnets(context.Context, *ManagementClientCheckAmlFSSubnetsOptions) (ManagementClientCheckAmlFSSubnetsResponse, error)`
- New function `*ManagementClient.GetRequiredAmlFSSubnetsSize(context.Context, *ManagementClientGetRequiredAmlFSSubnetsSizeOptions) (ManagementClientGetRequiredAmlFSSubnetsSizeResponse, error)`
- New struct `AmlFilesystem`
- New struct `AmlFilesystemArchive`
- New struct `AmlFilesystemArchiveInfo`
- New struct `AmlFilesystemArchiveStatus`
- New struct `AmlFilesystemCheckSubnetError`
- New struct `AmlFilesystemCheckSubnetErrorFilesystemSubnet`
- New struct `AmlFilesystemClientInfo`
- New struct `AmlFilesystemContainerStorageInterface`
- New struct `AmlFilesystemEncryptionSettings`
- New struct `AmlFilesystemHealth`
- New struct `AmlFilesystemHsmSettings`
- New struct `AmlFilesystemIdentity`
- New struct `AmlFilesystemProperties`
- New struct `AmlFilesystemPropertiesHsm`
- New struct `AmlFilesystemPropertiesMaintenanceWindow`
- New struct `AmlFilesystemSubnetInfo`
- New struct `AmlFilesystemUpdate`
- New struct `AmlFilesystemUpdateProperties`
- New struct `AmlFilesystemUpdatePropertiesMaintenanceWindow`
- New struct `AmlFilesystemsListResult`
- New struct `RequiredAmlFilesystemSubnetsSize`
- New struct `RequiredAmlFilesystemSubnetsSizeInfo`
- New struct `Resource`
- New struct `SKUName`
- New struct `TrackedResource`


## 3.1.0 (2023-04-07)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 3.0.0 (2023-02-24)
### Breaking Changes

- Operation `*CachesClient.Update` has been changed to LRO, use `*CachesClient.BeginUpdate` instead.

### Features Added

- New function `*StorageTargetsClient.BeginRestoreDefaults(context.Context, string, string, string, *StorageTargetsClientBeginRestoreDefaultsOptions) (*runtime.Poller[StorageTargetsClientRestoreDefaultsResponse], error)`
- New field `VerificationTimer` in struct `BlobNfsTarget`
- New field `WriteBackTimer` in struct `BlobNfsTarget`
- New field `VerificationTimer` in struct `Nfs3Target`
- New field `WriteBackTimer` in struct `Nfs3Target`


## 2.0.0 (2022-07-06)
### Breaking Changes

- Function `*StorageTargetsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, *StorageTargetsClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, string, StorageTarget, *StorageTargetsClientBeginCreateOrUpdateOptions)`
- Function `*CachesClient.Update` parameter(s) have been changed from `(context.Context, string, string, *CachesClientUpdateOptions)` to `(context.Context, string, string, Cache, *CachesClientUpdateOptions)`
- Function `*CachesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, *CachesClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, Cache, *CachesClientBeginCreateOrUpdateOptions)`
- Field `Cache` of struct `CachesClientUpdateOptions` has been removed
- Field `Cache` of struct `CachesClientBeginCreateOrUpdateOptions` has been removed
- Field `Storagetarget` of struct `StorageTargetsClientBeginCreateOrUpdateOptions` has been removed

### Features Added

- New const `PrimingJobStateQueued`
- New const `PrimingJobStateComplete`
- New const `PrimingJobStateRunning`
- New const `HealthStateTypeStartFailed`
- New const `HealthStateTypeUpgradeFailed`
- New const `PrimingJobStatePaused`
- New const `HealthStateTypeWaitingForKey`
- New function `*CachesClient.BeginStartPrimingJob(context.Context, string, string, *CachesClientBeginStartPrimingJobOptions) (*runtime.Poller[CachesClientStartPrimingJobResponse], error)`
- New function `PossiblePrimingJobStateValues() []PrimingJobState`
- New function `*CachesClient.BeginStopPrimingJob(context.Context, string, string, *CachesClientBeginStopPrimingJobOptions) (*runtime.Poller[CachesClientStopPrimingJobResponse], error)`
- New function `*CachesClient.BeginResumePrimingJob(context.Context, string, string, *CachesClientBeginResumePrimingJobOptions) (*runtime.Poller[CachesClientResumePrimingJobResponse], error)`
- New function `*CachesClient.BeginPausePrimingJob(context.Context, string, string, *CachesClientBeginPausePrimingJobOptions) (*runtime.Poller[CachesClientPausePrimingJobResponse], error)`
- New function `*CachesClient.BeginSpaceAllocation(context.Context, string, string, *CachesClientBeginSpaceAllocationOptions) (*runtime.Poller[CachesClientSpaceAllocationResponse], error)`
- New struct `CacheUpgradeSettings`
- New struct `CachesClientBeginPausePrimingJobOptions`
- New struct `CachesClientBeginResumePrimingJobOptions`
- New struct `CachesClientBeginSpaceAllocationOptions`
- New struct `CachesClientBeginStartPrimingJobOptions`
- New struct `CachesClientBeginStopPrimingJobOptions`
- New struct `CachesClientPausePrimingJobResponse`
- New struct `CachesClientResumePrimingJobResponse`
- New struct `CachesClientSpaceAllocationResponse`
- New struct `CachesClientStartPrimingJobResponse`
- New struct `CachesClientStopPrimingJobResponse`
- New struct `LogSpecification`
- New struct `PrimingJob`
- New struct `PrimingJobIDParameter`
- New struct `StorageTargetSpaceAllocation`
- New field `LogSpecifications` in struct `APIOperationPropertiesServiceSpecification`
- New field `AllocationPercentage` in struct `StorageTargetProperties`
- New field `UpgradeSettings` in struct `CacheProperties`
- New field `PrimingJobs` in struct `CacheProperties`
- New field `SpaceAllocation` in struct `CacheProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagecache/armstoragecache` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).