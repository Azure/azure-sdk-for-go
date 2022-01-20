# Release History

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*ADCOperationsClient.List` parameter(s) have been changed from `(context.Context, *ADCOperationsListOptions)` to `(context.Context, *ADCOperationsClientListOptions)`
- Function `*ADCOperationsClient.List` return value(s) have been changed from `(ADCOperationsListResponse, error)` to `(ADCOperationsClientListResponse, error)`
- Function `*ADCCatalogsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, *ADCCatalogsBeginDeleteOptions)` to `(context.Context, string, *ADCCatalogsClientBeginDeleteOptions)`
- Function `*ADCCatalogsClient.BeginDelete` return value(s) have been changed from `(ADCCatalogsDeletePollerResponse, error)` to `(ADCCatalogsClientDeletePollerResponse, error)`
- Function `*ADCCatalogsClient.ListtByResourceGroup` parameter(s) have been changed from `(context.Context, string, *ADCCatalogsListtByResourceGroupOptions)` to `(context.Context, string, *ADCCatalogsClientListtByResourceGroupOptions)`
- Function `*ADCCatalogsClient.ListtByResourceGroup` return value(s) have been changed from `(ADCCatalogsListtByResourceGroupResponse, error)` to `(ADCCatalogsClientListtByResourceGroupResponse, error)`
- Function `*ADCCatalogsClient.Update` parameter(s) have been changed from `(context.Context, string, ADCCatalog, *ADCCatalogsUpdateOptions)` to `(context.Context, string, ADCCatalog, *ADCCatalogsClientUpdateOptions)`
- Function `*ADCCatalogsClient.Update` return value(s) have been changed from `(ADCCatalogsUpdateResponse, error)` to `(ADCCatalogsClientUpdateResponse, error)`
- Function `*ADCCatalogsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, ADCCatalog, *ADCCatalogsCreateOrUpdateOptions)` to `(context.Context, string, ADCCatalog, *ADCCatalogsClientCreateOrUpdateOptions)`
- Function `*ADCCatalogsClient.CreateOrUpdate` return value(s) have been changed from `(ADCCatalogsCreateOrUpdateResponse, error)` to `(ADCCatalogsClientCreateOrUpdateResponse, error)`
- Function `*ADCCatalogsClient.Get` parameter(s) have been changed from `(context.Context, string, *ADCCatalogsGetOptions)` to `(context.Context, string, *ADCCatalogsClientGetOptions)`
- Function `*ADCCatalogsClient.Get` return value(s) have been changed from `(ADCCatalogsGetResponse, error)` to `(ADCCatalogsClientGetResponse, error)`
- Function `*ADCCatalogsDeletePollerResponse.Resume` has been removed
- Function `*ADCCatalogsDeletePoller.FinalResponse` has been removed
- Function `*ADCCatalogsDeletePoller.Done` has been removed
- Function `*ADCCatalogsDeletePoller.ResumeToken` has been removed
- Function `ADCCatalogsDeletePollerResponse.PollUntilDone` has been removed
- Function `*ADCCatalogsDeletePoller.Poll` has been removed
- Struct `ADCCatalogsBeginDeleteOptions` has been removed
- Struct `ADCCatalogsCreateOrUpdateOptions` has been removed
- Struct `ADCCatalogsCreateOrUpdateResponse` has been removed
- Struct `ADCCatalogsCreateOrUpdateResult` has been removed
- Struct `ADCCatalogsDeletePoller` has been removed
- Struct `ADCCatalogsDeletePollerResponse` has been removed
- Struct `ADCCatalogsDeleteResponse` has been removed
- Struct `ADCCatalogsGetOptions` has been removed
- Struct `ADCCatalogsGetResponse` has been removed
- Struct `ADCCatalogsGetResult` has been removed
- Struct `ADCCatalogsListtByResourceGroupOptions` has been removed
- Struct `ADCCatalogsListtByResourceGroupResponse` has been removed
- Struct `ADCCatalogsListtByResourceGroupResult` has been removed
- Struct `ADCCatalogsUpdateOptions` has been removed
- Struct `ADCCatalogsUpdateResponse` has been removed
- Struct `ADCCatalogsUpdateResult` has been removed
- Struct `ADCOperationsListOptions` has been removed
- Struct `ADCOperationsListResponse` has been removed
- Struct `ADCOperationsListResult` has been removed
- Field `Resource` of struct `ADCCatalog` has been removed

### Features Added

- New function `*ADCCatalogsClientDeletePoller.ResumeToken() (string, error)`
- New function `*ADCCatalogsClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*ADCCatalogsClientDeletePoller.Done() bool`
- New function `*ADCCatalogsClientDeletePoller.FinalResponse(context.Context) (ADCCatalogsClientDeleteResponse, error)`
- New function `*ADCCatalogsClientDeletePollerResponse.Resume(context.Context, *ADCCatalogsClient, string) error`
- New function `ADCCatalogsClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (ADCCatalogsClientDeleteResponse, error)`
- New struct `ADCCatalogsClientBeginDeleteOptions`
- New struct `ADCCatalogsClientCreateOrUpdateOptions`
- New struct `ADCCatalogsClientCreateOrUpdateResponse`
- New struct `ADCCatalogsClientCreateOrUpdateResult`
- New struct `ADCCatalogsClientDeletePoller`
- New struct `ADCCatalogsClientDeletePollerResponse`
- New struct `ADCCatalogsClientDeleteResponse`
- New struct `ADCCatalogsClientGetOptions`
- New struct `ADCCatalogsClientGetResponse`
- New struct `ADCCatalogsClientGetResult`
- New struct `ADCCatalogsClientListtByResourceGroupOptions`
- New struct `ADCCatalogsClientListtByResourceGroupResponse`
- New struct `ADCCatalogsClientListtByResourceGroupResult`
- New struct `ADCCatalogsClientUpdateOptions`
- New struct `ADCCatalogsClientUpdateResponse`
- New struct `ADCCatalogsClientUpdateResult`
- New struct `ADCOperationsClientListOptions`
- New struct `ADCOperationsClientListResponse`
- New struct `ADCOperationsClientListResult`
- New field `Location` in struct `ADCCatalog`
- New field `Tags` in struct `ADCCatalog`
- New field `ID` in struct `ADCCatalog`
- New field `Name` in struct `ADCCatalog`
- New field `Type` in struct `ADCCatalog`
- New field `Etag` in struct `ADCCatalog`


## 0.1.0 (2021-12-07)

- Initial preview release.
