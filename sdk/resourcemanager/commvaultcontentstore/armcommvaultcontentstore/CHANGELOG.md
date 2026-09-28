# Release History

## 0.2.0 (2026-09-28)
### Breaking Changes

- Function `*CloudAccountsClient.BeginUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, cloudAccountName string, properties CloudAccountUpdate, options *CloudAccountsClientBeginUpdateOptions)` to `(ctx context.Context, resourceGroupName string, cloudAccountName string, properties CloudAccount, options *CloudAccountsClientBeginUpdateOptions)`
- Struct `CloudAccountUpdate` has been removed
- Struct `CloudAccountUpdateProperties` has been removed
- Field `SaaSGUID` of struct `ActivateSaaSParameterRequest` has been removed
- Field `BackupAdminOnCcaCreate`, `MultiPersonAuthorizationOnCcaCreate` of struct `CloudAccountProperties` has been removed

### Features Added

- New enum type `ComplianceLockStatus` with values `ComplianceLockStatusDisabled`, `ComplianceLockStatusDisablementPending`, `ComplianceLockStatusEnabled`
- New function `*StoragesClient.DisableComplianceLock(ctx context.Context, resourceGroupName string, cloudAccountName string, storageName string, options *StoragesClientDisableComplianceLockOptions) (StoragesClientDisableComplianceLockResponse, error)`
- New function `*StoragesClient.EnableComplianceLock(ctx context.Context, resourceGroupName string, cloudAccountName string, storageName string, options *StoragesClientEnableComplianceLockOptions) (StoragesClientEnableComplianceLockResponse, error)`
- New function `*StoragesClient.Refresh(ctx context.Context, resourceGroupName string, cloudAccountName string, storageName string, options *StoragesClientRefreshOptions) (StoragesClientRefreshResponse, error)`
- New struct `ActivateSaaSRequestParam`
- New struct `CompanyProfile`
- New field `ActivateSaaSRequestParam`, `PublisherID`, `SaasGUID` in struct `ActivateSaaSParameterRequest`
- New field `Company`, `RoleAssignmentsOnCcaCreate` in struct `CloudAccountProperties`
- New field `ComplianceLockStatus` in struct `StorageProperties`


## 0.1.0 (2026-08-20)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/commvaultcontentstore/armcommvaultcontentstore` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).