# Release History

## 1.0.0 (2022-05-18)
### Breaking Changes

- Function `*WorkspaceCollectionsClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[WorkspaceCollectionsClientDeleteResponse], error)` to `(*runtime.Poller[WorkspaceCollectionsClientDeleteResponse], error)`
- Function `WorkspaceCollection.MarshalJSON` has been removed
- Function `OperationList.MarshalJSON` has been removed
- Function `WorkspaceCollectionList.MarshalJSON` has been removed
- Function `Error.MarshalJSON` has been removed
- Function `WorkspaceList.MarshalJSON` has been removed


## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*WorkspacesClient.List` has been removed
- Function `*WorkspaceCollectionsClient.ListByResourceGroup` has been removed
- Function `*WorkspaceCollectionsClient.ListBySubscription` has been removed

### Features Added

- New function `*WorkspacesClient.NewListPager(string, string, *WorkspacesClientListOptions) *runtime.Pager[WorkspacesClientListResponse]`
- New function `*WorkspaceCollectionsClient.NewListByResourceGroupPager(string, *WorkspaceCollectionsClientListByResourceGroupOptions) *runtime.Pager[WorkspaceCollectionsClientListByResourceGroupResponse]`
- New function `*WorkspaceCollectionsClient.NewListBySubscriptionPager(*WorkspaceCollectionsClientListBySubscriptionOptions) *runtime.Pager[WorkspaceCollectionsClientListBySubscriptionResponse]`


## 0.3.0 (2022-04-12)
### Breaking Changes

- Function `*WorkspacesClient.List` parameter(s) have been changed from `(context.Context, string, string, *WorkspacesClientListOptions)` to `(string, string, *WorkspacesClientListOptions)`
- Function `*WorkspacesClient.List` return value(s) have been changed from `(WorkspacesClientListResponse, error)` to `(*runtime.Pager[WorkspacesClientListResponse])`
- Function `NewWorkspaceCollectionsClient` return value(s) have been changed from `(*WorkspaceCollectionsClient)` to `(*WorkspaceCollectionsClient, error)`
- Function `NewWorkspacesClient` return value(s) have been changed from `(*WorkspacesClient)` to `(*WorkspacesClient, error)`
- Function `NewManagementClient` return value(s) have been changed from `(*ManagementClient)` to `(*ManagementClient, error)`
- Function `*WorkspaceCollectionsClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, string, *WorkspaceCollectionsClientListByResourceGroupOptions)` to `(string, *WorkspaceCollectionsClientListByResourceGroupOptions)`
- Function `*WorkspaceCollectionsClient.ListByResourceGroup` return value(s) have been changed from `(WorkspaceCollectionsClientListByResourceGroupResponse, error)` to `(*runtime.Pager[WorkspaceCollectionsClientListByResourceGroupResponse])`
- Function `*WorkspaceCollectionsClient.BeginDelete` return value(s) have been changed from `(WorkspaceCollectionsClientDeletePollerResponse, error)` to `(*armruntime.Poller[WorkspaceCollectionsClientDeleteResponse], error)`
- Function `*WorkspaceCollectionsClient.ListBySubscription` parameter(s) have been changed from `(context.Context, *WorkspaceCollectionsClientListBySubscriptionOptions)` to `(*WorkspaceCollectionsClientListBySubscriptionOptions)`
- Function `*WorkspaceCollectionsClient.ListBySubscription` return value(s) have been changed from `(WorkspaceCollectionsClientListBySubscriptionResponse, error)` to `(*runtime.Pager[WorkspaceCollectionsClientListBySubscriptionResponse])`
- Type of `Workspace.Properties` has been changed from `map[string]interface{}` to `interface{}`
- Type of `WorkspaceCollection.Properties` has been changed from `map[string]interface{}` to `interface{}`
- Function `*WorkspaceCollectionsClientDeletePoller.ResumeToken` has been removed
- Function `*WorkspaceCollectionsClientDeletePoller.FinalResponse` has been removed
- Function `AzureSKUName.ToPtr` has been removed
- Function `*WorkspaceCollectionsClientDeletePoller.Done` has been removed
- Function `*WorkspaceCollectionsClientDeletePollerResponse.Resume` has been removed
- Function `*WorkspaceCollectionsClientDeletePoller.Poll` has been removed
- Function `AccessKeyName.ToPtr` has been removed
- Function `WorkspaceCollectionsClientDeletePollerResponse.PollUntilDone` has been removed
- Function `CheckNameReason.ToPtr` has been removed
- Function `AzureSKUTier.ToPtr` has been removed
- Struct `ManagementClientGetAvailableOperationsResult` has been removed
- Struct `WorkspaceCollectionsClientCheckNameAvailabilityResult` has been removed
- Struct `WorkspaceCollectionsClientCreateResult` has been removed
- Struct `WorkspaceCollectionsClientDeletePoller` has been removed
- Struct `WorkspaceCollectionsClientDeletePollerResponse` has been removed
- Struct `WorkspaceCollectionsClientGetAccessKeysResult` has been removed
- Struct `WorkspaceCollectionsClientGetByNameResult` has been removed
- Struct `WorkspaceCollectionsClientListByResourceGroupResult` has been removed
- Struct `WorkspaceCollectionsClientListBySubscriptionResult` has been removed
- Struct `WorkspaceCollectionsClientRegenerateKeyResult` has been removed
- Struct `WorkspaceCollectionsClientUpdateResult` has been removed
- Struct `WorkspacesClientListResult` has been removed
- Field `WorkspacesClientListResult` of struct `WorkspacesClientListResponse` has been removed
- Field `RawResponse` of struct `WorkspacesClientListResponse` has been removed
- Field `ManagementClientGetAvailableOperationsResult` of struct `ManagementClientGetAvailableOperationsResponse` has been removed
- Field `RawResponse` of struct `ManagementClientGetAvailableOperationsResponse` has been removed
- Field `RawResponse` of struct `WorkspaceCollectionsClientDeleteResponse` has been removed
- Field `WorkspaceCollectionsClientGetByNameResult` of struct `WorkspaceCollectionsClientGetByNameResponse` has been removed
- Field `RawResponse` of struct `WorkspaceCollectionsClientGetByNameResponse` has been removed
- Field `WorkspaceCollectionsClientRegenerateKeyResult` of struct `WorkspaceCollectionsClientRegenerateKeyResponse` has been removed
- Field `RawResponse` of struct `WorkspaceCollectionsClientRegenerateKeyResponse` has been removed
- Field `WorkspaceCollectionsClientCreateResult` of struct `WorkspaceCollectionsClientCreateResponse` has been removed
- Field `RawResponse` of struct `WorkspaceCollectionsClientCreateResponse` has been removed
- Field `WorkspaceCollectionsClientListByResourceGroupResult` of struct `WorkspaceCollectionsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `WorkspaceCollectionsClientListByResourceGroupResponse` has been removed
- Field `WorkspaceCollectionsClientUpdateResult` of struct `WorkspaceCollectionsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `WorkspaceCollectionsClientUpdateResponse` has been removed
- Field `WorkspaceCollectionsClientCheckNameAvailabilityResult` of struct `WorkspaceCollectionsClientCheckNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `WorkspaceCollectionsClientCheckNameAvailabilityResponse` has been removed
- Field `WorkspaceCollectionsClientListBySubscriptionResult` of struct `WorkspaceCollectionsClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `WorkspaceCollectionsClientListBySubscriptionResponse` has been removed
- Field `WorkspaceCollectionsClientGetAccessKeysResult` of struct `WorkspaceCollectionsClientGetAccessKeysResponse` has been removed
- Field `RawResponse` of struct `WorkspaceCollectionsClientGetAccessKeysResponse` has been removed
- Field `RawResponse` of struct `WorkspaceCollectionsClientMigrateResponse` has been removed

### Features Added

- New anonymous field `WorkspaceCollectionAccessKeys` in struct `WorkspaceCollectionsClientGetAccessKeysResponse`
- New anonymous field `OperationList` in struct `ManagementClientGetAvailableOperationsResponse`
- New anonymous field `WorkspaceCollection` in struct `WorkspaceCollectionsClientGetByNameResponse`
- New anonymous field `WorkspaceList` in struct `WorkspacesClientListResponse`
- New anonymous field `WorkspaceCollection` in struct `WorkspaceCollectionsClientUpdateResponse`
- New anonymous field `WorkspaceCollectionAccessKeys` in struct `WorkspaceCollectionsClientRegenerateKeyResponse`
- New field `ResumeToken` in struct `WorkspaceCollectionsClientBeginDeleteOptions`
- New anonymous field `WorkspaceCollectionList` in struct `WorkspaceCollectionsClientListBySubscriptionResponse`
- New anonymous field `CheckNameResponse` in struct `WorkspaceCollectionsClientCheckNameAvailabilityResponse`
- New anonymous field `WorkspaceCollectionList` in struct `WorkspaceCollectionsClientListByResourceGroupResponse`
- New anonymous field `WorkspaceCollection` in struct `WorkspaceCollectionsClientCreateResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*WorkspaceCollectionsClient.ListBySubscription` parameter(s) have been changed from `(context.Context, *WorkspaceCollectionsListBySubscriptionOptions)` to `(context.Context, *WorkspaceCollectionsClientListBySubscriptionOptions)`
- Function `*WorkspaceCollectionsClient.ListBySubscription` return value(s) have been changed from `(WorkspaceCollectionsListBySubscriptionResponse, error)` to `(WorkspaceCollectionsClientListBySubscriptionResponse, error)`
- Function `*WorkspaceCollectionsClient.Migrate` parameter(s) have been changed from `(context.Context, string, MigrateWorkspaceCollectionRequest, *WorkspaceCollectionsMigrateOptions)` to `(context.Context, string, MigrateWorkspaceCollectionRequest, *WorkspaceCollectionsClientMigrateOptions)`
- Function `*WorkspaceCollectionsClient.Migrate` return value(s) have been changed from `(WorkspaceCollectionsMigrateResponse, error)` to `(WorkspaceCollectionsClientMigrateResponse, error)`
- Function `*WorkspaceCollectionsClient.Create` parameter(s) have been changed from `(context.Context, string, string, CreateWorkspaceCollectionRequest, *WorkspaceCollectionsCreateOptions)` to `(context.Context, string, string, CreateWorkspaceCollectionRequest, *WorkspaceCollectionsClientCreateOptions)`
- Function `*WorkspaceCollectionsClient.Create` return value(s) have been changed from `(WorkspaceCollectionsCreateResponse, error)` to `(WorkspaceCollectionsClientCreateResponse, error)`
- Function `*WorkspaceCollectionsClient.GetByName` parameter(s) have been changed from `(context.Context, string, string, *WorkspaceCollectionsGetByNameOptions)` to `(context.Context, string, string, *WorkspaceCollectionsClientGetByNameOptions)`
- Function `*WorkspaceCollectionsClient.GetByName` return value(s) have been changed from `(WorkspaceCollectionsGetByNameResponse, error)` to `(WorkspaceCollectionsClientGetByNameResponse, error)`
- Function `*WorkspaceCollectionsClient.RegenerateKey` parameter(s) have been changed from `(context.Context, string, string, WorkspaceCollectionAccessKey, *WorkspaceCollectionsRegenerateKeyOptions)` to `(context.Context, string, string, WorkspaceCollectionAccessKey, *WorkspaceCollectionsClientRegenerateKeyOptions)`
- Function `*WorkspaceCollectionsClient.RegenerateKey` return value(s) have been changed from `(WorkspaceCollectionsRegenerateKeyResponse, error)` to `(WorkspaceCollectionsClientRegenerateKeyResponse, error)`
- Function `*WorkspaceCollectionsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *WorkspaceCollectionsBeginDeleteOptions)` to `(context.Context, string, string, *WorkspaceCollectionsClientBeginDeleteOptions)`
- Function `*WorkspaceCollectionsClient.BeginDelete` return value(s) have been changed from `(WorkspaceCollectionsDeletePollerResponse, error)` to `(WorkspaceCollectionsClientDeletePollerResponse, error)`
- Function `*WorkspaceCollectionsClient.GetAccessKeys` parameter(s) have been changed from `(context.Context, string, string, *WorkspaceCollectionsGetAccessKeysOptions)` to `(context.Context, string, string, *WorkspaceCollectionsClientGetAccessKeysOptions)`
- Function `*WorkspaceCollectionsClient.GetAccessKeys` return value(s) have been changed from `(WorkspaceCollectionsGetAccessKeysResponse, error)` to `(WorkspaceCollectionsClientGetAccessKeysResponse, error)`
- Function `*WorkspaceCollectionsClient.CheckNameAvailability` parameter(s) have been changed from `(context.Context, string, CheckNameRequest, *WorkspaceCollectionsCheckNameAvailabilityOptions)` to `(context.Context, string, CheckNameRequest, *WorkspaceCollectionsClientCheckNameAvailabilityOptions)`
- Function `*WorkspaceCollectionsClient.CheckNameAvailability` return value(s) have been changed from `(WorkspaceCollectionsCheckNameAvailabilityResponse, error)` to `(WorkspaceCollectionsClientCheckNameAvailabilityResponse, error)`
- Function `*WorkspacesClient.List` parameter(s) have been changed from `(context.Context, string, string, *WorkspacesListOptions)` to `(context.Context, string, string, *WorkspacesClientListOptions)`
- Function `*WorkspacesClient.List` return value(s) have been changed from `(WorkspacesListResponse, error)` to `(WorkspacesClientListResponse, error)`
- Function `*WorkspaceCollectionsClient.Update` parameter(s) have been changed from `(context.Context, string, string, UpdateWorkspaceCollectionRequest, *WorkspaceCollectionsUpdateOptions)` to `(context.Context, string, string, UpdateWorkspaceCollectionRequest, *WorkspaceCollectionsClientUpdateOptions)`
- Function `*WorkspaceCollectionsClient.Update` return value(s) have been changed from `(WorkspaceCollectionsUpdateResponse, error)` to `(WorkspaceCollectionsClientUpdateResponse, error)`
- Function `*WorkspaceCollectionsClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, string, *WorkspaceCollectionsListByResourceGroupOptions)` to `(context.Context, string, *WorkspaceCollectionsClientListByResourceGroupOptions)`
- Function `*WorkspaceCollectionsClient.ListByResourceGroup` return value(s) have been changed from `(WorkspaceCollectionsListByResourceGroupResponse, error)` to `(WorkspaceCollectionsClientListByResourceGroupResponse, error)`
- Function `*PowerBIEmbeddedManagementClient.GetAvailableOperations` has been removed
- Function `*WorkspaceCollectionsDeletePoller.FinalResponse` has been removed
- Function `*WorkspaceCollectionsDeletePoller.ResumeToken` has been removed
- Function `*WorkspaceCollectionsDeletePoller.Done` has been removed
- Function `NewPowerBIEmbeddedManagementClient` has been removed
- Function `*WorkspaceCollectionsDeletePoller.Poll` has been removed
- Function `WorkspaceCollectionsDeletePollerResponse.PollUntilDone` has been removed
- Function `*WorkspaceCollectionsDeletePollerResponse.Resume` has been removed
- Function `Error.Error` has been removed
- Struct `PowerBIEmbeddedManagementClient` has been removed
- Struct `PowerBIEmbeddedManagementClientGetAvailableOperationsOptions` has been removed
- Struct `PowerBIEmbeddedManagementClientGetAvailableOperationsResponse` has been removed
- Struct `PowerBIEmbeddedManagementClientGetAvailableOperationsResult` has been removed
- Struct `WorkspaceCollectionsBeginDeleteOptions` has been removed
- Struct `WorkspaceCollectionsCheckNameAvailabilityOptions` has been removed
- Struct `WorkspaceCollectionsCheckNameAvailabilityResponse` has been removed
- Struct `WorkspaceCollectionsCheckNameAvailabilityResult` has been removed
- Struct `WorkspaceCollectionsCreateOptions` has been removed
- Struct `WorkspaceCollectionsCreateResponse` has been removed
- Struct `WorkspaceCollectionsCreateResult` has been removed
- Struct `WorkspaceCollectionsDeletePoller` has been removed
- Struct `WorkspaceCollectionsDeletePollerResponse` has been removed
- Struct `WorkspaceCollectionsDeleteResponse` has been removed
- Struct `WorkspaceCollectionsGetAccessKeysOptions` has been removed
- Struct `WorkspaceCollectionsGetAccessKeysResponse` has been removed
- Struct `WorkspaceCollectionsGetAccessKeysResult` has been removed
- Struct `WorkspaceCollectionsGetByNameOptions` has been removed
- Struct `WorkspaceCollectionsGetByNameResponse` has been removed
- Struct `WorkspaceCollectionsGetByNameResult` has been removed
- Struct `WorkspaceCollectionsListByResourceGroupOptions` has been removed
- Struct `WorkspaceCollectionsListByResourceGroupResponse` has been removed
- Struct `WorkspaceCollectionsListByResourceGroupResult` has been removed
- Struct `WorkspaceCollectionsListBySubscriptionOptions` has been removed
- Struct `WorkspaceCollectionsListBySubscriptionResponse` has been removed
- Struct `WorkspaceCollectionsListBySubscriptionResult` has been removed
- Struct `WorkspaceCollectionsMigrateOptions` has been removed
- Struct `WorkspaceCollectionsMigrateResponse` has been removed
- Struct `WorkspaceCollectionsRegenerateKeyOptions` has been removed
- Struct `WorkspaceCollectionsRegenerateKeyResponse` has been removed
- Struct `WorkspaceCollectionsRegenerateKeyResult` has been removed
- Struct `WorkspaceCollectionsUpdateOptions` has been removed
- Struct `WorkspaceCollectionsUpdateResponse` has been removed
- Struct `WorkspaceCollectionsUpdateResult` has been removed
- Struct `WorkspacesListOptions` has been removed
- Struct `WorkspacesListResponse` has been removed
- Struct `WorkspacesListResult` has been removed

### Features Added

- New function `Error.MarshalJSON() ([]byte, error)`
- New function `*WorkspaceCollectionsClientDeletePoller.FinalResponse(context.Context) (WorkspaceCollectionsClientDeleteResponse, error)`
- New function `*WorkspaceCollectionsClientDeletePollerResponse.Resume(context.Context, *WorkspaceCollectionsClient, string) error`
- New function `NewManagementClient(azcore.TokenCredential, *arm.ClientOptions) *ManagementClient`
- New function `*WorkspaceCollectionsClientDeletePoller.Done() bool`
- New function `*WorkspaceCollectionsClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*WorkspaceCollectionsClientDeletePoller.ResumeToken() (string, error)`
- New function `*ManagementClient.GetAvailableOperations(context.Context, *ManagementClientGetAvailableOperationsOptions) (ManagementClientGetAvailableOperationsResponse, error)`
- New function `WorkspaceCollectionsClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (WorkspaceCollectionsClientDeleteResponse, error)`
- New struct `ManagementClient`
- New struct `ManagementClientGetAvailableOperationsOptions`
- New struct `ManagementClientGetAvailableOperationsResponse`
- New struct `ManagementClientGetAvailableOperationsResult`
- New struct `WorkspaceCollectionsClientBeginDeleteOptions`
- New struct `WorkspaceCollectionsClientCheckNameAvailabilityOptions`
- New struct `WorkspaceCollectionsClientCheckNameAvailabilityResponse`
- New struct `WorkspaceCollectionsClientCheckNameAvailabilityResult`
- New struct `WorkspaceCollectionsClientCreateOptions`
- New struct `WorkspaceCollectionsClientCreateResponse`
- New struct `WorkspaceCollectionsClientCreateResult`
- New struct `WorkspaceCollectionsClientDeletePoller`
- New struct `WorkspaceCollectionsClientDeletePollerResponse`
- New struct `WorkspaceCollectionsClientDeleteResponse`
- New struct `WorkspaceCollectionsClientGetAccessKeysOptions`
- New struct `WorkspaceCollectionsClientGetAccessKeysResponse`
- New struct `WorkspaceCollectionsClientGetAccessKeysResult`
- New struct `WorkspaceCollectionsClientGetByNameOptions`
- New struct `WorkspaceCollectionsClientGetByNameResponse`
- New struct `WorkspaceCollectionsClientGetByNameResult`
- New struct `WorkspaceCollectionsClientListByResourceGroupOptions`
- New struct `WorkspaceCollectionsClientListByResourceGroupResponse`
- New struct `WorkspaceCollectionsClientListByResourceGroupResult`
- New struct `WorkspaceCollectionsClientListBySubscriptionOptions`
- New struct `WorkspaceCollectionsClientListBySubscriptionResponse`
- New struct `WorkspaceCollectionsClientListBySubscriptionResult`
- New struct `WorkspaceCollectionsClientMigrateOptions`
- New struct `WorkspaceCollectionsClientMigrateResponse`
- New struct `WorkspaceCollectionsClientRegenerateKeyOptions`
- New struct `WorkspaceCollectionsClientRegenerateKeyResponse`
- New struct `WorkspaceCollectionsClientRegenerateKeyResult`
- New struct `WorkspaceCollectionsClientUpdateOptions`
- New struct `WorkspaceCollectionsClientUpdateResponse`
- New struct `WorkspaceCollectionsClientUpdateResult`
- New struct `WorkspacesClientListOptions`
- New struct `WorkspacesClientListResponse`
- New struct `WorkspacesClientListResult`


## 0.1.0 (2021-12-22)

- Init release.
