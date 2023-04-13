# Release History

## 2.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 2.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 2.0.0 (2022-09-26)
### Breaking Changes

- Function `*RoleAssignmentsClient.NewListForResourcePager` parameter(s) have been changed from `(string, string, string, string, string, *RoleAssignmentsClientListForResourceOptions)` to `(string, string, string, string, *RoleAssignmentsClientListForResourceOptions)`
- Type of `RoleAssignment.Properties` has been changed from `*RoleAssignmentPropertiesWithScope` to `*RoleAssignmentProperties`
- Function `*RoleAssignmentsClient.NewListPager` has been renamed to `*RoleAssignmentsClient.NewListForSubscriptionPager`

### Features Added

- New function `*DenyAssignmentsClient.Get(context.Context, string, string, *DenyAssignmentsClientGetOptions) (DenyAssignmentsClientGetResponse, error)`
- New function `*DenyAssignmentsClient.NewListForScopePager(string, *DenyAssignmentsClientListForScopeOptions) *runtime.Pager[DenyAssignmentsClientListForScopeResponse]`
- New function `*DenyAssignmentsClient.NewListForResourcePager(string, string, string, string, string, *DenyAssignmentsClientListForResourceOptions) *runtime.Pager[DenyAssignmentsClientListForResourceResponse]`
- New function `*DenyAssignmentsClient.NewListForResourceGroupPager(string, *DenyAssignmentsClientListForResourceGroupOptions) *runtime.Pager[DenyAssignmentsClientListForResourceGroupResponse]`
- New function `*DenyAssignmentsClient.GetByID(context.Context, string, *DenyAssignmentsClientGetByIDOptions) (DenyAssignmentsClientGetByIDResponse, error)`
- New function `*DenyAssignmentsClient.NewListPager(*DenyAssignmentsClientListOptions) *runtime.Pager[DenyAssignmentsClientListResponse]`
- New function `NewDenyAssignmentsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DenyAssignmentsClient, error)`
- New struct `DenyAssignment`
- New struct `DenyAssignmentFilter`
- New struct `DenyAssignmentListResult`
- New struct `DenyAssignmentPermission`
- New struct `DenyAssignmentProperties`
- New struct `DenyAssignmentsClient`
- New struct `DenyAssignmentsClientGetByIDOptions`
- New struct `DenyAssignmentsClientGetByIDResponse`
- New struct `DenyAssignmentsClientGetOptions`
- New struct `DenyAssignmentsClientGetResponse`
- New struct `DenyAssignmentsClientListForResourceGroupOptions`
- New struct `DenyAssignmentsClientListForResourceGroupResponse`
- New struct `DenyAssignmentsClientListForResourceOptions`
- New struct `DenyAssignmentsClientListForResourceResponse`
- New struct `DenyAssignmentsClientListForScopeOptions`
- New struct `DenyAssignmentsClientListForScopeResponse`
- New struct `DenyAssignmentsClientListOptions`
- New struct `DenyAssignmentsClientListResponse`
- New struct `ValidationResponse`
- New struct `ValidationResponseErrorInfo`
- New field `TenantID` in struct `RoleAssignmentsClientGetByIDOptions`
- New field `DataActions` in struct `Permission`
- New field `NotDataActions` in struct `Permission`
- New field `TenantID` in struct `RoleAssignmentsClientListForResourceOptions`
- New field `UpdatedBy` in struct `RoleAssignmentProperties`
- New field `Condition` in struct `RoleAssignmentProperties`
- New field `CreatedOn` in struct `RoleAssignmentProperties`
- New field `UpdatedOn` in struct `RoleAssignmentProperties`
- New field `CreatedBy` in struct `RoleAssignmentProperties`
- New field `ConditionVersion` in struct `RoleAssignmentProperties`
- New field `DelegatedManagedIdentityResourceID` in struct `RoleAssignmentProperties`
- New field `Description` in struct `RoleAssignmentProperties`
- New field `PrincipalType` in struct `RoleAssignmentProperties`
- New field `Scope` in struct `RoleAssignmentProperties`
- New field `TenantID` in struct `RoleAssignmentsClientDeleteByIDOptions`
- New field `IsDataAction` in struct `ProviderOperation`
- New field `TenantID` in struct `RoleAssignmentsClientDeleteOptions`
- New field `Type` in struct `RoleDefinitionFilter`
- New field `TenantID` in struct `RoleAssignmentsClientListForResourceGroupOptions`
- New field `TenantID` in struct `RoleAssignmentsClientGetOptions`
- New field `SkipToken` in struct `RoleAssignmentsClientListForScopeOptions`
- New field `TenantID` in struct `RoleAssignmentsClientListForScopeOptions`


## 1.0.0 (2022-06-02)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
