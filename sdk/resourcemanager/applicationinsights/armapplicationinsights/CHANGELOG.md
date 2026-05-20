# Release History

## 2.0.0 (2026-05-20)
### Breaking Changes

- Function `*WorkbooksClient.Update` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties Workbook, options *WorkbooksClientUpdateOptions)` to `(ctx context.Context, resourceGroupName string, resourceName string, workbookUpdateParameters WorkbookUpdateParameters, options *WorkbooksClientUpdateOptions)`
- Type of `Workbook.Kind` has been changed from `*SharedTypeKind` to `*WorkbookSharedTypeKind`
- Type of `WorkbookProperties.TimeModified` has been changed from `*string` to `*time.Time`
- Enum `SharedTypeKind` has been removed
- Function `*ClientFactory.NewMyWorkbooksClient` has been removed
- Function `NewMyWorkbooksClient` has been removed
- Function `*MyWorkbooksClient.CreateOrUpdate` has been removed
- Function `*MyWorkbooksClient.Delete` has been removed
- Function `*MyWorkbooksClient.Get` has been removed
- Function `*MyWorkbooksClient.NewListByResourceGroupPager` has been removed
- Function `*MyWorkbooksClient.NewListBySubscriptionPager` has been removed
- Function `*MyWorkbooksClient.Update` has been removed
- Struct `AnnotationError` has been removed
- Struct `ComponentsResource` has been removed
- Struct `ErrorFieldContract` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ErrorResponseComponents` has been removed
- Struct `ErrorResponseComponentsError` has been removed
- Struct `InnerError` has been removed
- Struct `LinkProperties` has been removed
- Struct `MyWorkbook` has been removed
- Struct `MyWorkbookError` has been removed
- Struct `MyWorkbookProperties` has been removed
- Struct `MyWorkbookResource` has been removed
- Struct `MyWorkbooksListResult` has been removed
- Struct `WebtestsResource` has been removed
- Struct `WorkItemConfigurationError` has been removed
- Struct `WorkbookError` has been removed
- Struct `WorkbookResource` has been removed
- Field `Name`, `SharedTypeKind`, `SourceResourceID`, `WorkbookID` of struct `WorkbookProperties` has been removed

### Features Added

- New value `WebTestKindStandard` added to enum type `WebTestKind`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `StorageType` with values `StorageTypeServiceProfiler`
- New enum type `WorkbookSharedTypeKind` with values `WorkbookSharedTypeKindShared`
- New enum type `WorkbookUpdateSharedTypeKind` with values `WorkbookUpdateSharedTypeKindShared`
- New function `*ClientFactory.NewComponentLinkedStorageAccountsClient() *ComponentLinkedStorageAccountsClient`
- New function `*ClientFactory.NewDeletedWorkbooksClient() *DeletedWorkbooksClient`
- New function `*ClientFactory.NewLiveTokenClient() *LiveTokenClient`
- New function `*ClientFactory.NewOperationsClient() *OperationsClient`
- New function `*ClientFactory.NewWorkbookTemplatesClient() *WorkbookTemplatesClient`
- New function `NewComponentLinkedStorageAccountsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ComponentLinkedStorageAccountsClient, error)`
- New function `*ComponentLinkedStorageAccountsClient.CreateAndUpdate(ctx context.Context, resourceGroupName string, resourceName string, storageType StorageType, linkedStorageAccountsProperties ComponentLinkedStorageAccounts, options *ComponentLinkedStorageAccountsClientCreateAndUpdateOptions) (ComponentLinkedStorageAccountsClientCreateAndUpdateResponse, error)`
- New function `*ComponentLinkedStorageAccountsClient.Delete(ctx context.Context, resourceGroupName string, resourceName string, storageType StorageType, options *ComponentLinkedStorageAccountsClientDeleteOptions) (ComponentLinkedStorageAccountsClientDeleteResponse, error)`
- New function `*ComponentLinkedStorageAccountsClient.Get(ctx context.Context, resourceGroupName string, resourceName string, storageType StorageType, options *ComponentLinkedStorageAccountsClientGetOptions) (ComponentLinkedStorageAccountsClientGetResponse, error)`
- New function `*ComponentLinkedStorageAccountsClient.Update(ctx context.Context, resourceGroupName string, resourceName string, storageType StorageType, linkedStorageAccountsProperties ComponentLinkedStorageAccountsPatch, options *ComponentLinkedStorageAccountsClientUpdateOptions) (ComponentLinkedStorageAccountsClientUpdateResponse, error)`
- New function `NewDeletedWorkbooksClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DeletedWorkbooksClient, error)`
- New function `*DeletedWorkbooksClient.NewListBySubscriptionPager(options *DeletedWorkbooksClientListBySubscriptionOptions) *runtime.Pager[DeletedWorkbooksClientListBySubscriptionResponse]`
- New function `NewLiveTokenClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*LiveTokenClient, error)`
- New function `*LiveTokenClient.Get(ctx context.Context, resourceURI string, options *LiveTokenClientGetOptions) (LiveTokenClientGetResponse, error)`
- New function `NewOperationsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationsClient, error)`
- New function `*OperationsClient.NewListPager(options *OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `NewWorkbookTemplatesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkbookTemplatesClient, error)`
- New function `*WorkbookTemplatesClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, workbookTemplateProperties WorkbookTemplate, options *WorkbookTemplatesClientCreateOrUpdateOptions) (WorkbookTemplatesClientCreateOrUpdateResponse, error)`
- New function `*WorkbookTemplatesClient.Delete(ctx context.Context, resourceGroupName string, resourceName string, options *WorkbookTemplatesClientDeleteOptions) (WorkbookTemplatesClientDeleteResponse, error)`
- New function `*WorkbookTemplatesClient.Get(ctx context.Context, resourceGroupName string, resourceName string, options *WorkbookTemplatesClientGetOptions) (WorkbookTemplatesClientGetResponse, error)`
- New function `*WorkbookTemplatesClient.NewListByResourceGroupPager(resourceGroupName string, options *WorkbookTemplatesClientListByResourceGroupOptions) *runtime.Pager[WorkbookTemplatesClientListByResourceGroupResponse]`
- New function `*WorkbookTemplatesClient.Update(ctx context.Context, resourceGroupName string, resourceName string, workbookTemplateUpdateParameters WorkbookTemplateUpdateParameters, options *WorkbookTemplatesClientUpdateOptions) (WorkbookTemplatesClientUpdateResponse, error)`
- New function `*WorkbooksClient.NewListBySubscriptionPager(category CategoryType, options *WorkbooksClientListBySubscriptionOptions) *runtime.Pager[WorkbooksClientListBySubscriptionResponse]`
- New function `*WorkbooksClient.RevisionGet(ctx context.Context, resourceGroupName string, resourceName string, revisionID string, options *WorkbooksClientRevisionGetOptions) (WorkbooksClientRevisionGetResponse, error)`
- New function `*WorkbooksClient.NewRevisionsListPager(resourceGroupName string, resourceName string, options *WorkbooksClientRevisionsListOptions) *runtime.Pager[WorkbooksClientRevisionsListResponse]`
- New struct `ComponentLinkedStorageAccounts`
- New struct `ComponentLinkedStorageAccountsPatch`
- New struct `DeletedWorkbook`
- New struct `DeletedWorkbookProperties`
- New struct `DeletedWorkbooksListResult`
- New struct `HeaderField`
- New struct `LinkedStorageAccountsProperties`
- New struct `LiveTokenResponse`
- New struct `SystemData`
- New struct `UserAssignedIdentity`
- New struct `WebTestPropertiesRequest`
- New struct `WebTestPropertiesValidationRules`
- New struct `WebTestPropertiesValidationRulesContentValidation`
- New struct `WorkbookPropertiesUpdateParameters`
- New struct `WorkbookResourceIdentity`
- New struct `WorkbookTemplate`
- New struct `WorkbookTemplateGallery`
- New struct `WorkbookTemplateLocalizedGallery`
- New struct `WorkbookTemplateProperties`
- New struct `WorkbookTemplateUpdateParameters`
- New struct `WorkbookTemplatesListResult`
- New struct `WorkbookUpdateParameters`
- New field `NextLink` in struct `AnnotationsListResult`
- New field `NextLink` in struct `ComponentAPIKeyListResult`
- New field `NextLink` in struct `WebTestLocationsListResult`
- New field `Request`, `ValidationRules` in struct `WebTestProperties`
- New field `NextLink` in struct `WorkItemConfigurationsListResult`
- New field `Etag`, `Identity`, `SystemData` in struct `Workbook`
- New field `Description`, `DisplayName`, `Revision`, `SourceID`, `StorageURI` in struct `WorkbookProperties`
- New field `SourceID` in struct `WorkbooksClientCreateOrUpdateOptions`
- New field `CanFetchContent` in struct `WorkbooksClientGetOptions`
- New field `SourceID` in struct `WorkbooksClientListByResourceGroupOptions`
- New field `SourceID` in struct `WorkbooksClientUpdateOptions`
- New field `NextLink` in struct `WorkbooksListResult`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.1.0 (2023-04-06)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-06-02)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).