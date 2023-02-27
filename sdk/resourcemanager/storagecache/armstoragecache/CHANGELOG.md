# Release History

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