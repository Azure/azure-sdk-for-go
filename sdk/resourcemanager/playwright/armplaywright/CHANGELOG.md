# Release History

## 1.1.0-beta.1 (2026-04-06)
### Features Added

- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New struct `ManagedServiceIdentity`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `Workspace`
- New field `Reporting`, `StorageURI` in struct `WorkspaceProperties`
- New field `Identity` in struct `WorkspaceUpdate`
- New field `Reporting`, `StorageURI` in struct `WorkspaceUpdateProperties`


## 1.0.0 (2025-08-26)
### Breaking Changes

- Type of `Quota.Name` has been changed from `*QuotaName` to `*string`
- Type of `WorkspaceQuota.Name` has been changed from `*QuotaName` to `*string`

### Features Added

- New field `WorkspaceID` in struct `WorkspaceProperties`


## 0.1.0 (2025-07-15)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/playwright/armplaywright` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
