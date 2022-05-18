# Release History

## 0.4.0 (2022-05-18)
### Breaking Changes

- Function `*ProjectsClient.BeginCreate` return value(s) have been changed from `(*armruntime.Poller[ProjectsClientCreateResponse], error)` to `(*runtime.Poller[ProjectsClientCreateResponse], error)`
- Function `ProjectResourceListResult.MarshalJSON` has been removed
- Function `AccountResourceListResult.MarshalJSON` has been removed
- Function `ExtensionResourceListResult.MarshalJSON` has been removed
- Function `OperationListResult.MarshalJSON` has been removed


## 0.3.0 (2022-04-13)
### Breaking Changes

- Function `*ProjectsClient.BeginCreate` return value(s) have been changed from `(ProjectsClientCreatePollerResponse, error)` to `(*armruntime.Poller[ProjectsClientCreateResponse], error)`
- Function `NewExtensionsClient` return value(s) have been changed from `(*ExtensionsClient)` to `(*ExtensionsClient, error)`
- Function `NewProjectsClient` return value(s) have been changed from `(*ProjectsClient)` to `(*ProjectsClient, error)`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `NewAccountsClient` return value(s) have been changed from `(*AccountsClient)` to `(*AccountsClient, error)`
- Function `*ProjectsClientCreatePollerResponse.Resume` has been removed
- Function `AccountResourceRequestOperationType.ToPtr` has been removed
- Function `*ProjectsClientCreatePoller.Done` has been removed
- Function `*ProjectsClientCreatePoller.FinalResponse` has been removed
- Function `*ProjectsClientCreatePoller.ResumeToken` has been removed
- Function `*ProjectsClientCreatePoller.Poll` has been removed
- Function `ProjectsClientCreatePollerResponse.PollUntilDone` has been removed
- Struct `AccountsClientCheckNameAvailabilityResult` has been removed
- Struct `AccountsClientCreateOrUpdateResult` has been removed
- Struct `AccountsClientGetResult` has been removed
- Struct `AccountsClientListByResourceGroupResult` has been removed
- Struct `AccountsClientUpdateResult` has been removed
- Struct `ExtensionsClientCreateResult` has been removed
- Struct `ExtensionsClientGetResult` has been removed
- Struct `ExtensionsClientListByAccountResult` has been removed
- Struct `ExtensionsClientUpdateResult` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `ProjectsClientCreatePoller` has been removed
- Struct `ProjectsClientCreatePollerResponse` has been removed
- Struct `ProjectsClientCreateResult` has been removed
- Struct `ProjectsClientGetJobStatusResult` has been removed
- Struct `ProjectsClientGetResult` has been removed
- Struct `ProjectsClientListByResourceGroupResult` has been removed
- Struct `ProjectsClientUpdateResult` has been removed
- Field `AccountsClientGetResult` of struct `AccountsClientGetResponse` has been removed
- Field `RawResponse` of struct `AccountsClientGetResponse` has been removed
- Field `ExtensionsClientGetResult` of struct `ExtensionsClientGetResponse` has been removed
- Field `RawResponse` of struct `ExtensionsClientGetResponse` has been removed
- Field `ProjectsClientUpdateResult` of struct `ProjectsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ProjectsClientUpdateResponse` has been removed
- Field `ProjectsClientCreateResult` of struct `ProjectsClientCreateResponse` has been removed
- Field `RawResponse` of struct `ProjectsClientCreateResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `ExtensionsClientUpdateResult` of struct `ExtensionsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ExtensionsClientUpdateResponse` has been removed
- Field `AccountsClientCheckNameAvailabilityResult` of struct `AccountsClientCheckNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `AccountsClientCheckNameAvailabilityResponse` has been removed
- Field `ProjectsClientGetResult` of struct `ProjectsClientGetResponse` has been removed
- Field `RawResponse` of struct `ProjectsClientGetResponse` has been removed
- Field `ExtensionsClientCreateResult` of struct `ExtensionsClientCreateResponse` has been removed
- Field `RawResponse` of struct `ExtensionsClientCreateResponse` has been removed
- Field `ProjectsClientGetJobStatusResult` of struct `ProjectsClientGetJobStatusResponse` has been removed
- Field `RawResponse` of struct `ProjectsClientGetJobStatusResponse` has been removed
- Field `AccountsClientCreateOrUpdateResult` of struct `AccountsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `AccountsClientCreateOrUpdateResponse` has been removed
- Field `AccountsClientListByResourceGroupResult` of struct `AccountsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `AccountsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `ExtensionsClientDeleteResponse` has been removed
- Field `AccountsClientUpdateResult` of struct `AccountsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `AccountsClientUpdateResponse` has been removed
- Field `ExtensionsClientListByAccountResult` of struct `ExtensionsClientListByAccountResponse` has been removed
- Field `RawResponse` of struct `ExtensionsClientListByAccountResponse` has been removed
- Field `RawResponse` of struct `AccountsClientDeleteResponse` has been removed
- Field `ProjectsClientListByResourceGroupResult` of struct `ProjectsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `ProjectsClientListByResourceGroupResponse` has been removed

### Features Added

- New field `ResumeToken` in struct `ProjectsClientBeginCreateOptions`
- New anonymous field `ExtensionResource` in struct `ExtensionsClientCreateResponse`
- New anonymous field `AccountResourceListResult` in struct `AccountsClientListByResourceGroupResponse`
- New anonymous field `ExtensionResourceListResult` in struct `ExtensionsClientListByAccountResponse`
- New anonymous field `ProjectResource` in struct `ProjectsClientGetJobStatusResponse`
- New anonymous field `ProjectResource` in struct `ProjectsClientCreateResponse`
- New anonymous field `ProjectResourceListResult` in struct `ProjectsClientListByResourceGroupResponse`
- New anonymous field `ProjectResource` in struct `ProjectsClientUpdateResponse`
- New anonymous field `ProjectResource` in struct `ProjectsClientGetResponse`
- New anonymous field `ExtensionResource` in struct `ExtensionsClientUpdateResponse`
- New anonymous field `ExtensionResource` in struct `ExtensionsClientGetResponse`
- New anonymous field `CheckNameAvailabilityResult` in struct `AccountsClientCheckNameAvailabilityResponse`
- New anonymous field `AccountResource` in struct `AccountsClientCreateOrUpdateResponse`
- New anonymous field `AccountResource` in struct `AccountsClientUpdateResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `AccountResource` in struct `AccountsClientGetResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*ProjectsClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, string, string, *ProjectsListByResourceGroupOptions)` to `(context.Context, string, string, *ProjectsClientListByResourceGroupOptions)`
- Function `*ProjectsClient.ListByResourceGroup` return value(s) have been changed from `(ProjectsListByResourceGroupResponse, error)` to `(ProjectsClientListByResourceGroupResponse, error)`
- Function `*AccountsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *AccountsGetOptions)` to `(context.Context, string, string, *AccountsClientGetOptions)`
- Function `*AccountsClient.Get` return value(s) have been changed from `(AccountsGetResponse, error)` to `(AccountsClientGetResponse, error)`
- Function `*AccountsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *AccountsDeleteOptions)` to `(context.Context, string, string, *AccountsClientDeleteOptions)`
- Function `*AccountsClient.Delete` return value(s) have been changed from `(AccountsDeleteResponse, error)` to `(AccountsClientDeleteResponse, error)`
- Function `*ExtensionsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *ExtensionsGetOptions)` to `(context.Context, string, string, string, *ExtensionsClientGetOptions)`
- Function `*ExtensionsClient.Get` return value(s) have been changed from `(ExtensionsGetResponse, error)` to `(ExtensionsClientGetResponse, error)`
- Function `*ExtensionsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, string, *ExtensionsDeleteOptions)` to `(context.Context, string, string, string, *ExtensionsClientDeleteOptions)`
- Function `*ExtensionsClient.Delete` return value(s) have been changed from `(ExtensionsDeleteResponse, error)` to `(ExtensionsClientDeleteResponse, error)`
- Function `*ProjectsClient.GetJobStatus` parameter(s) have been changed from `(context.Context, string, string, string, string, string, *ProjectsGetJobStatusOptions)` to `(context.Context, string, string, string, string, string, *ProjectsClientGetJobStatusOptions)`
- Function `*ProjectsClient.GetJobStatus` return value(s) have been changed from `(ProjectsGetJobStatusResponse, error)` to `(ProjectsClientGetJobStatusResponse, error)`
- Function `*AccountsClient.CheckNameAvailability` parameter(s) have been changed from `(context.Context, CheckNameAvailabilityParameter, *AccountsCheckNameAvailabilityOptions)` to `(context.Context, CheckNameAvailabilityParameter, *AccountsClientCheckNameAvailabilityOptions)`
- Function `*AccountsClient.CheckNameAvailability` return value(s) have been changed from `(AccountsCheckNameAvailabilityResponse, error)` to `(AccountsClientCheckNameAvailabilityResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Function `*ExtensionsClient.Create` parameter(s) have been changed from `(context.Context, string, string, string, ExtensionResourceRequest, *ExtensionsCreateOptions)` to `(context.Context, string, string, string, ExtensionResourceRequest, *ExtensionsClientCreateOptions)`
- Function `*ExtensionsClient.Create` return value(s) have been changed from `(ExtensionsCreateResponse, error)` to `(ExtensionsClientCreateResponse, error)`
- Function `*AccountsClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, string, *AccountsListByResourceGroupOptions)` to `(context.Context, string, *AccountsClientListByResourceGroupOptions)`
- Function `*AccountsClient.ListByResourceGroup` return value(s) have been changed from `(AccountsListByResourceGroupResponse, error)` to `(AccountsClientListByResourceGroupResponse, error)`
- Function `*ExtensionsClient.Update` parameter(s) have been changed from `(context.Context, string, string, string, ExtensionResourceRequest, *ExtensionsUpdateOptions)` to `(context.Context, string, string, string, ExtensionResourceRequest, *ExtensionsClientUpdateOptions)`
- Function `*ExtensionsClient.Update` return value(s) have been changed from `(ExtensionsUpdateResponse, error)` to `(ExtensionsClientUpdateResponse, error)`
- Function `*AccountsClient.Update` parameter(s) have been changed from `(context.Context, string, string, AccountTagRequest, *AccountsUpdateOptions)` to `(context.Context, string, string, AccountTagRequest, *AccountsClientUpdateOptions)`
- Function `*AccountsClient.Update` return value(s) have been changed from `(AccountsUpdateResponse, error)` to `(AccountsClientUpdateResponse, error)`
- Function `*ExtensionsClient.ListByAccount` parameter(s) have been changed from `(context.Context, string, string, *ExtensionsListByAccountOptions)` to `(context.Context, string, string, *ExtensionsClientListByAccountOptions)`
- Function `*ExtensionsClient.ListByAccount` return value(s) have been changed from `(ExtensionsListByAccountResponse, error)` to `(ExtensionsClientListByAccountResponse, error)`
- Function `*ProjectsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *ProjectsGetOptions)` to `(context.Context, string, string, string, *ProjectsClientGetOptions)`
- Function `*ProjectsClient.Get` return value(s) have been changed from `(ProjectsGetResponse, error)` to `(ProjectsClientGetResponse, error)`
- Function `*ProjectsClient.Update` parameter(s) have been changed from `(context.Context, string, string, string, ProjectResource, *ProjectsUpdateOptions)` to `(context.Context, string, string, string, ProjectResource, *ProjectsClientUpdateOptions)`
- Function `*ProjectsClient.Update` return value(s) have been changed from `(ProjectsUpdateResponse, error)` to `(ProjectsClientUpdateResponse, error)`
- Function `*AccountsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, AccountResourceRequest, *AccountsCreateOrUpdateOptions)` to `(context.Context, string, string, AccountResourceRequest, *AccountsClientCreateOrUpdateOptions)`
- Function `*AccountsClient.CreateOrUpdate` return value(s) have been changed from `(AccountsCreateOrUpdateResponse, error)` to `(AccountsClientCreateOrUpdateResponse, error)`
- Function `*ProjectsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, string, ProjectResource, *ProjectsBeginCreateOptions)` to `(context.Context, string, string, string, ProjectResource, *ProjectsClientBeginCreateOptions)`
- Function `*ProjectsClient.BeginCreate` return value(s) have been changed from `(ProjectsCreatePollerResponse, error)` to `(ProjectsClientCreatePollerResponse, error)`
- Function `ProjectsCreatePollerResponse.PollUntilDone` has been removed
- Function `*ProjectsCreatePoller.FinalResponse` has been removed
- Function `*ProjectsCreatePollerResponse.Resume` has been removed
- Function `*ProjectsCreatePoller.Done` has been removed
- Function `*ProjectsCreatePoller.Poll` has been removed
- Function `*ProjectsCreatePoller.ResumeToken` has been removed
- Struct `AccountsCheckNameAvailabilityOptions` has been removed
- Struct `AccountsCheckNameAvailabilityResponse` has been removed
- Struct `AccountsCheckNameAvailabilityResult` has been removed
- Struct `AccountsCreateOrUpdateOptions` has been removed
- Struct `AccountsCreateOrUpdateResponse` has been removed
- Struct `AccountsCreateOrUpdateResult` has been removed
- Struct `AccountsDeleteOptions` has been removed
- Struct `AccountsDeleteResponse` has been removed
- Struct `AccountsGetOptions` has been removed
- Struct `AccountsGetResponse` has been removed
- Struct `AccountsGetResult` has been removed
- Struct `AccountsListByResourceGroupOptions` has been removed
- Struct `AccountsListByResourceGroupResponse` has been removed
- Struct `AccountsListByResourceGroupResult` has been removed
- Struct `AccountsUpdateOptions` has been removed
- Struct `AccountsUpdateResponse` has been removed
- Struct `AccountsUpdateResult` has been removed
- Struct `ExtensionsCreateOptions` has been removed
- Struct `ExtensionsCreateResponse` has been removed
- Struct `ExtensionsCreateResult` has been removed
- Struct `ExtensionsDeleteOptions` has been removed
- Struct `ExtensionsDeleteResponse` has been removed
- Struct `ExtensionsGetOptions` has been removed
- Struct `ExtensionsGetResponse` has been removed
- Struct `ExtensionsGetResult` has been removed
- Struct `ExtensionsListByAccountOptions` has been removed
- Struct `ExtensionsListByAccountResponse` has been removed
- Struct `ExtensionsListByAccountResult` has been removed
- Struct `ExtensionsUpdateOptions` has been removed
- Struct `ExtensionsUpdateResponse` has been removed
- Struct `ExtensionsUpdateResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `ProjectsBeginCreateOptions` has been removed
- Struct `ProjectsCreatePoller` has been removed
- Struct `ProjectsCreatePollerResponse` has been removed
- Struct `ProjectsCreateResponse` has been removed
- Struct `ProjectsCreateResult` has been removed
- Struct `ProjectsGetJobStatusOptions` has been removed
- Struct `ProjectsGetJobStatusResponse` has been removed
- Struct `ProjectsGetJobStatusResult` has been removed
- Struct `ProjectsGetOptions` has been removed
- Struct `ProjectsGetResponse` has been removed
- Struct `ProjectsGetResult` has been removed
- Struct `ProjectsListByResourceGroupOptions` has been removed
- Struct `ProjectsListByResourceGroupResponse` has been removed
- Struct `ProjectsListByResourceGroupResult` has been removed
- Struct `ProjectsUpdateOptions` has been removed
- Struct `ProjectsUpdateResponse` has been removed
- Struct `ProjectsUpdateResult` has been removed
- Field `Resource` of struct `AccountResource` has been removed
- Field `Resource` of struct `ProjectResource` has been removed
- Field `Resource` of struct `ExtensionResource` has been removed

### Features Added

- New function `*ProjectsClientCreatePoller.ResumeToken() (string, error)`
- New function `*ProjectsClientCreatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*ProjectsClientCreatePollerResponse.Resume(context.Context, *ProjectsClient, string) error`
- New function `ProjectsClientCreatePollerResponse.PollUntilDone(context.Context, time.Duration) (ProjectsClientCreateResponse, error)`
- New function `*ProjectsClientCreatePoller.Done() bool`
- New function `*ProjectsClientCreatePoller.FinalResponse(context.Context) (ProjectsClientCreateResponse, error)`
- New struct `AccountsClientCheckNameAvailabilityOptions`
- New struct `AccountsClientCheckNameAvailabilityResponse`
- New struct `AccountsClientCheckNameAvailabilityResult`
- New struct `AccountsClientCreateOrUpdateOptions`
- New struct `AccountsClientCreateOrUpdateResponse`
- New struct `AccountsClientCreateOrUpdateResult`
- New struct `AccountsClientDeleteOptions`
- New struct `AccountsClientDeleteResponse`
- New struct `AccountsClientGetOptions`
- New struct `AccountsClientGetResponse`
- New struct `AccountsClientGetResult`
- New struct `AccountsClientListByResourceGroupOptions`
- New struct `AccountsClientListByResourceGroupResponse`
- New struct `AccountsClientListByResourceGroupResult`
- New struct `AccountsClientUpdateOptions`
- New struct `AccountsClientUpdateResponse`
- New struct `AccountsClientUpdateResult`
- New struct `ExtensionsClientCreateOptions`
- New struct `ExtensionsClientCreateResponse`
- New struct `ExtensionsClientCreateResult`
- New struct `ExtensionsClientDeleteOptions`
- New struct `ExtensionsClientDeleteResponse`
- New struct `ExtensionsClientGetOptions`
- New struct `ExtensionsClientGetResponse`
- New struct `ExtensionsClientGetResult`
- New struct `ExtensionsClientListByAccountOptions`
- New struct `ExtensionsClientListByAccountResponse`
- New struct `ExtensionsClientListByAccountResult`
- New struct `ExtensionsClientUpdateOptions`
- New struct `ExtensionsClientUpdateResponse`
- New struct `ExtensionsClientUpdateResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `ProjectsClientBeginCreateOptions`
- New struct `ProjectsClientCreatePoller`
- New struct `ProjectsClientCreatePollerResponse`
- New struct `ProjectsClientCreateResponse`
- New struct `ProjectsClientCreateResult`
- New struct `ProjectsClientGetJobStatusOptions`
- New struct `ProjectsClientGetJobStatusResponse`
- New struct `ProjectsClientGetJobStatusResult`
- New struct `ProjectsClientGetOptions`
- New struct `ProjectsClientGetResponse`
- New struct `ProjectsClientGetResult`
- New struct `ProjectsClientListByResourceGroupOptions`
- New struct `ProjectsClientListByResourceGroupResponse`
- New struct `ProjectsClientListByResourceGroupResult`
- New struct `ProjectsClientUpdateOptions`
- New struct `ProjectsClientUpdateResponse`
- New struct `ProjectsClientUpdateResult`
- New field `Tags` in struct `ProjectResource`
- New field `ID` in struct `ProjectResource`
- New field `Name` in struct `ProjectResource`
- New field `Type` in struct `ProjectResource`
- New field `Location` in struct `ProjectResource`
- New field `Name` in struct `ExtensionResource`
- New field `Type` in struct `ExtensionResource`
- New field `Location` in struct `ExtensionResource`
- New field `Tags` in struct `ExtensionResource`
- New field `ID` in struct `ExtensionResource`
- New field `Tags` in struct `AccountResource`
- New field `ID` in struct `AccountResource`
- New field `Name` in struct `AccountResource`
- New field `Type` in struct `AccountResource`
- New field `Location` in struct `AccountResource`


## 0.1.0 (2021-12-23)

- Init release.
