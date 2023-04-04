# Release History

## 2.0.0-beta.3 (2023-04-07)
### Other Changes


## 2.0.0-beta.2 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 2.0.0-beta.1 (2022-06-02)
### Breaking Changes

- Function `*WorkbooksClient.Update` parameter(s) have been changed from `(context.Context, string, string, Workbook, *WorkbooksClientUpdateOptions)` to `(context.Context, string, string, *WorkbooksClientUpdateOptions)`
- Type of `MyWorkbook.Kind` has been changed from `*SharedTypeKind` to `*Kind`
- Type of `WorkbookProperties.TimeModified` has been changed from `*string` to `*time.Time`
- Type of `Workbook.Kind` has been changed from `*SharedTypeKind` to `*WorkbookSharedTypeKind`
- Const `SharedTypeKindUser` has been removed
- Const `SharedTypeKindShared` has been removed
- Function `PossibleSharedTypeKindValues` has been removed
- Struct `ErrorFieldContract` has been removed
- Struct `LinkProperties` has been removed
- Field `Code` of struct `WorkbookError` has been removed
- Field `Details` of struct `WorkbookError` has been removed
- Field `Message` of struct `WorkbookError` has been removed
- Field `SharedTypeKind` of struct `WorkbookProperties` has been removed
- Field `SourceResourceID` of struct `WorkbookProperties` has been removed
- Field `WorkbookID` of struct `WorkbookProperties` has been removed
- Field `Name` of struct `WorkbookProperties` has been removed
- Field `Code` of struct `MyWorkbookError` has been removed
- Field `Details` of struct `MyWorkbookError` has been removed
- Field `Message` of struct `MyWorkbookError` has been removed

### Features Added

- New const `MyWorkbookManagedIdentityTypeNone`
- New const `WorkbookUpdateSharedTypeKindShared`
- New const `CreatedByTypeManagedIdentity`
- New const `ManagedServiceIdentityTypeUserAssigned`
- New const `StorageTypeServiceProfiler`
- New const `KindUser`
- New const `CreatedByTypeApplication`
- New const `CreatedByTypeUser`
- New const `ManagedServiceIdentityTypeSystemAssignedUserAssigned`
- New const `MyWorkbookManagedIdentityTypeUserAssigned`
- New const `WorkbookSharedTypeKindShared`
- New const `CreatedByTypeKey`
- New const `KindShared`
- New const `ManagedServiceIdentityTypeNone`
- New const `ManagedServiceIdentityTypeSystemAssigned`
- New function `TrackedResource.MarshalJSON() ([]byte, error)`
- New function `PossibleWorkbookUpdateSharedTypeKindValues() []WorkbookUpdateSharedTypeKind`
- New function `WorkbookTemplate.MarshalJSON() ([]byte, error)`
- New function `*WorkbooksClient.RevisionGet(context.Context, string, string, string, *WorkbooksClientRevisionGetOptions) (WorkbooksClientRevisionGetResponse, error)`
- New function `WorkbookTemplateProperties.MarshalJSON() ([]byte, error)`
- New function `WorkbookTemplateLocalizedGallery.MarshalJSON() ([]byte, error)`
- New function `WorkbookTemplateResource.MarshalJSON() ([]byte, error)`
- New function `PossibleCreatedByTypeValues() []CreatedByType`
- New function `PossibleStorageTypeValues() []StorageType`
- New function `*WorkbookProperties.UnmarshalJSON([]byte) error`
- New function `WorkbookPropertiesUpdateParameters.MarshalJSON() ([]byte, error)`
- New function `PossibleManagedServiceIdentityTypeValues() []ManagedServiceIdentityType`
- New function `*WorkbooksClient.NewListBySubscriptionPager(CategoryType, *WorkbooksClientListBySubscriptionOptions) *runtime.Pager[WorkbooksClientListBySubscriptionResponse]`
- New function `WorkbookTemplateUpdateParameters.MarshalJSON() ([]byte, error)`
- New function `PossibleWorkbookSharedTypeKindValues() []WorkbookSharedTypeKind`
- New function `PossibleMyWorkbookManagedIdentityTypeValues() []MyWorkbookManagedIdentityType`
- New function `ComponentLinkedStorageAccountsPatch.MarshalJSON() ([]byte, error)`
- New function `WorkbookUpdateParameters.MarshalJSON() ([]byte, error)`
- New function `*WorkbooksClient.NewRevisionsListPager(string, string, *WorkbooksClientRevisionsListOptions) *runtime.Pager[WorkbooksClientRevisionsListResponse]`
- New function `PossibleKindValues() []Kind`
- New function `ManagedServiceIdentity.MarshalJSON() ([]byte, error)`
- New function `WorkbookResourceIdentity.MarshalJSON() ([]byte, error)`
- New function `*SystemData.UnmarshalJSON([]byte) error`
- New function `SystemData.MarshalJSON() ([]byte, error)`
- New struct `ComponentLinkedStorageAccounts`
- New struct `ComponentLinkedStorageAccountsClientCreateAndUpdateOptions`
- New struct `ComponentLinkedStorageAccountsClientCreateAndUpdateResponse`
- New struct `ComponentLinkedStorageAccountsClientDeleteOptions`
- New struct `ComponentLinkedStorageAccountsClientDeleteResponse`
- New struct `ComponentLinkedStorageAccountsClientGetOptions`
- New struct `ComponentLinkedStorageAccountsClientGetResponse`
- New struct `ComponentLinkedStorageAccountsClientUpdateOptions`
- New struct `ComponentLinkedStorageAccountsClientUpdateResponse`
- New struct `ComponentLinkedStorageAccountsPatch`
- New struct `ErrorDefinition`
- New struct `ErrorResponseLinkedStorage`
- New struct `ErrorResponseLinkedStorageError`
- New struct `InnerErrorTrace`
- New struct `LinkedStorageAccountsProperties`
- New struct `LiveTokenClientGetOptions`
- New struct `LiveTokenClientGetResponse`
- New struct `LiveTokenResponse`
- New struct `ManagedServiceIdentity`
- New struct `MyWorkbookManagedIdentity`
- New struct `MyWorkbookUserAssignedIdentities`
- New struct `OperationInfo`
- New struct `OperationLive`
- New struct `OperationsListResult`
- New struct `ProxyResource`
- New struct `Resource`
- New struct `SystemData`
- New struct `TrackedResource`
- New struct `UserAssignedIdentity`
- New struct `WorkbookErrorDefinition`
- New struct `WorkbookInnerErrorTrace`
- New struct `WorkbookPropertiesUpdateParameters`
- New struct `WorkbookResourceIdentity`
- New struct `WorkbookTemplate`
- New struct `WorkbookTemplateError`
- New struct `WorkbookTemplateErrorBody`
- New struct `WorkbookTemplateErrorFieldContract`
- New struct `WorkbookTemplateGallery`
- New struct `WorkbookTemplateLocalizedGallery`
- New struct `WorkbookTemplateProperties`
- New struct `WorkbookTemplateResource`
- New struct `WorkbookTemplateUpdateParameters`
- New struct `WorkbookTemplatesClientCreateOrUpdateOptions`
- New struct `WorkbookTemplatesClientCreateOrUpdateResponse`
- New struct `WorkbookTemplatesClientDeleteOptions`
- New struct `WorkbookTemplatesClientDeleteResponse`
- New struct `WorkbookTemplatesClientGetOptions`
- New struct `WorkbookTemplatesClientGetResponse`
- New struct `WorkbookTemplatesClientListByResourceGroupOptions`
- New struct `WorkbookTemplatesClientListByResourceGroupResponse`
- New struct `WorkbookTemplatesClientUpdateOptions`
- New struct `WorkbookTemplatesClientUpdateResponse`
- New struct `WorkbookTemplatesListResult`
- New struct `WorkbookUpdateParameters`
- New struct `WorkbooksClientListBySubscriptionOptions`
- New struct `WorkbooksClientListBySubscriptionResponse`
- New struct `WorkbooksClientRevisionGetOptions`
- New struct `WorkbooksClientRevisionGetResponse`
- New struct `WorkbooksClientRevisionsListOptions`
- New struct `WorkbooksClientRevisionsListResponse`
- New field `CanFetchContent` in struct `WorkbooksClientGetOptions`
- New field `SourceID` in struct `MyWorkbooksClientListByResourceGroupOptions`
- New field `NextLink` in struct `WorkbooksListResult`
- New field `Etag` in struct `WorkbookResource`
- New field `Identity` in struct `WorkbookResource`
- New field `Kind` in struct `WorkbookResource`
- New field `SystemData` in struct `Workbook`
- New field `Etag` in struct `Workbook`
- New field `Identity` in struct `Workbook`
- New field `Etag` in struct `MyWorkbook`
- New field `Identity` in struct `MyWorkbook`
- New field `SystemData` in struct `MyWorkbook`
- New field `Error` in struct `WorkbookError`
- New field `SourceID` in struct `WorkbooksClientUpdateOptions`
- New field `WorkbookUpdateParameters` in struct `WorkbooksClientUpdateOptions`
- New field `StorageURI` in struct `MyWorkbookProperties`
- New field `SourceID` in struct `MyWorkbooksClientUpdateOptions`
- New field `NextLink` in struct `MyWorkbooksListResult`
- New field `Error` in struct `MyWorkbookError`
- New field `SourceID` in struct `MyWorkbooksClientCreateOrUpdateOptions`
- New field `SourceID` in struct `WorkbooksClientListByResourceGroupOptions`
- New field `SourceID` in struct `WorkbookProperties`
- New field `DisplayName` in struct `WorkbookProperties`
- New field `Revision` in struct `WorkbookProperties`
- New field `StorageURI` in struct `WorkbookProperties`
- New field `Description` in struct `WorkbookProperties`
- New field `Identity` in struct `MyWorkbookResource`
- New field `Etag` in struct `MyWorkbookResource`
- New field `SourceID` in struct `WorkbooksClientCreateOrUpdateOptions`


## 1.0.0 (2022-06-02)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).