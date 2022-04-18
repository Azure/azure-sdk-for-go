# Release History

## 0.5.1 (2022-04-18)
### Other Changes


## 0.5.0 (2022-04-18)
### Breaking Changes

- Function `*Client.ListByResourceGroup` has been removed
- Function `*Client.ListBySubscription` has been removed

### Features Added

- New function `*Client.NewListBySubscriptionPager(*ClientListBySubscriptionOptions) *runtime.Pager[ClientListBySubscriptionResponse]`
- New function `*Client.NewListByResourceGroupPager(string, *ClientListByResourceGroupOptions) *runtime.Pager[ClientListByResourceGroupResponse]`


## 0.4.0 (2022-04-14)
### Breaking Changes

- Function `*Client.ListByResourceGroup` return value(s) have been changed from `(*ClientListByResourceGroupPager)` to `(*runtime.Pager[ClientListByResourceGroupResponse])`
- Function `NewClient` return value(s) have been changed from `(*Client)` to `(*Client, error)`
- Function `*Client.BeginCreate` return value(s) have been changed from `(ClientCreatePollerResponse, error)` to `(*armruntime.Poller[ClientCreateResponse], error)`
- Function `*Client.ListBySubscription` return value(s) have been changed from `(*ClientListBySubscriptionPager)` to `(*runtime.Pager[ClientListBySubscriptionResponse])`
- Function `*ClientUpdateResult.UnmarshalJSON` has been removed
- Function `*ClientCreatePoller.Done` has been removed
- Function `CleanupOptions.ToPtr` has been removed
- Function `ScriptProvisioningState.ToPtr` has been removed
- Function `*ClientListBySubscriptionPager.Err` has been removed
- Function `*ClientCreatePoller.Poll` has been removed
- Function `*ClientCreatePollerResponse.Resume` has been removed
- Function `*ClientGetResult.UnmarshalJSON` has been removed
- Function `ManagedServiceIdentityType.ToPtr` has been removed
- Function `*ClientListByResourceGroupPager.PageResponse` has been removed
- Function `ClientCreatePollerResponse.PollUntilDone` has been removed
- Function `*ClientListBySubscriptionPager.PageResponse` has been removed
- Function `*ClientCreatePoller.FinalResponse` has been removed
- Function `*ClientListBySubscriptionPager.NextPage` has been removed
- Function `*ClientCreateResult.UnmarshalJSON` has been removed
- Function `*ClientListByResourceGroupPager.Err` has been removed
- Function `*ClientCreatePoller.ResumeToken` has been removed
- Function `*ClientListByResourceGroupPager.NextPage` has been removed
- Function `CreatedByType.ToPtr` has been removed
- Function `ScriptType.ToPtr` has been removed
- Struct `ClientCreatePoller` has been removed
- Struct `ClientCreatePollerResponse` has been removed
- Struct `ClientCreateResult` has been removed
- Struct `ClientGetLogsDefaultResult` has been removed
- Struct `ClientGetLogsResult` has been removed
- Struct `ClientGetResult` has been removed
- Struct `ClientListByResourceGroupPager` has been removed
- Struct `ClientListByResourceGroupResult` has been removed
- Struct `ClientListBySubscriptionPager` has been removed
- Struct `ClientListBySubscriptionResult` has been removed
- Struct `ClientUpdateResult` has been removed
- Field `ClientGetLogsDefaultResult` of struct `ClientGetLogsDefaultResponse` has been removed
- Field `RawResponse` of struct `ClientGetLogsDefaultResponse` has been removed
- Field `RawResponse` of struct `ClientDeleteResponse` has been removed
- Field `ClientListByResourceGroupResult` of struct `ClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `ClientListByResourceGroupResponse` has been removed
- Field `ClientGetResult` of struct `ClientGetResponse` has been removed
- Field `RawResponse` of struct `ClientGetResponse` has been removed
- Field `ClientCreateResult` of struct `ClientCreateResponse` has been removed
- Field `RawResponse` of struct `ClientCreateResponse` has been removed
- Field `ClientListBySubscriptionResult` of struct `ClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `ClientListBySubscriptionResponse` has been removed
- Field `ClientGetLogsResult` of struct `ClientGetLogsResponse` has been removed
- Field `RawResponse` of struct `ClientGetLogsResponse` has been removed
- Field `ClientUpdateResult` of struct `ClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ClientUpdateResponse` has been removed

### Features Added

- New function `*ClientUpdateResponse.UnmarshalJSON([]byte) error`
- New function `*ClientCreateResponse.UnmarshalJSON([]byte) error`
- New function `*ClientGetResponse.UnmarshalJSON([]byte) error`
- New struct `Error`
- New anonymous field `DeploymentScriptClassification` in struct `ClientGetResponse`
- New anonymous field `ScriptLog` in struct `ClientGetLogsDefaultResponse`
- New anonymous field `DeploymentScriptListResult` in struct `ClientListByResourceGroupResponse`
- New anonymous field `ScriptLogsList` in struct `ClientGetLogsResponse`
- New anonymous field `DeploymentScriptListResult` in struct `ClientListBySubscriptionResponse`
- New anonymous field `DeploymentScriptClassification` in struct `ClientCreateResponse`
- New anonymous field `DeploymentScriptClassification` in struct `ClientUpdateResponse`
- New field `ResumeToken` in struct `ClientBeginCreateOptions`


## 0.3.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.3.0 (2022-02-16)
### Breaking Changes

- Type of `AzurePowerShellScriptProperties.Outputs` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Type of `DeploymentScriptPropertiesBase.Outputs` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Type of `AzureCliScriptProperties.Outputs` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Struct `Error` has been removed


## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*DeploymentScriptsClient.GetLogsDefault` has been removed
- Function `*DeploymentScriptsCreatePoller.Done` has been removed
- Function `*DeploymentScriptsClient.ListBySubscription` has been removed
- Function `*DeploymentScriptsCreatePoller.Poll` has been removed
- Function `*DeploymentScript.UnmarshalJSON` has been removed
- Function `*DeploymentScriptsClient.BeginCreate` has been removed
- Function `DeploymentScriptsCreatePollerResponse.PollUntilDone` has been removed
- Function `*DeploymentScriptsListBySubscriptionPager.PageResponse` has been removed
- Function `*DeploymentScriptsListByResourceGroupPager.Err` has been removed
- Function `ScriptLog.MarshalJSON` has been removed
- Function `*DeploymentScriptsGetResult.UnmarshalJSON` has been removed
- Function `NewDeploymentScriptsClient` has been removed
- Function `*DeploymentScriptsCreatePoller.ResumeToken` has been removed
- Function `*ScriptLog.UnmarshalJSON` has been removed
- Function `AzureResourceBase.MarshalJSON` has been removed
- Function `*DeploymentScriptUpdateParameter.UnmarshalJSON` has been removed
- Function `DeploymentScriptsError.Error` has been removed
- Function `*DeploymentScriptsListBySubscriptionPager.NextPage` has been removed
- Function `*DeploymentScriptsListByResourceGroupPager.PageResponse` has been removed
- Function `*DeploymentScriptsClient.Get` has been removed
- Function `*DeploymentScriptsCreatePollerResponse.Resume` has been removed
- Function `*DeploymentScriptsListBySubscriptionPager.Err` has been removed
- Function `*DeploymentScriptsUpdateResult.UnmarshalJSON` has been removed
- Function `*AzureResourceBase.UnmarshalJSON` has been removed
- Function `*DeploymentScriptsClient.Delete` has been removed
- Function `*DeploymentScriptsClient.Update` has been removed
- Function `*DeploymentScriptsCreateResult.UnmarshalJSON` has been removed
- Function `*DeploymentScriptsCreatePoller.FinalResponse` has been removed
- Function `*DeploymentScriptsListByResourceGroupPager.NextPage` has been removed
- Function `*DeploymentScriptsClient.GetLogs` has been removed
- Function `*DeploymentScriptsClient.ListByResourceGroup` has been removed
- Struct `DeploymentScriptsBeginCreateOptions` has been removed
- Struct `DeploymentScriptsClient` has been removed
- Struct `DeploymentScriptsCreatePoller` has been removed
- Struct `DeploymentScriptsCreatePollerResponse` has been removed
- Struct `DeploymentScriptsCreateResponse` has been removed
- Struct `DeploymentScriptsCreateResult` has been removed
- Struct `DeploymentScriptsDeleteOptions` has been removed
- Struct `DeploymentScriptsDeleteResponse` has been removed
- Struct `DeploymentScriptsError` has been removed
- Struct `DeploymentScriptsGetLogsDefaultOptions` has been removed
- Struct `DeploymentScriptsGetLogsDefaultResponse` has been removed
- Struct `DeploymentScriptsGetLogsDefaultResult` has been removed
- Struct `DeploymentScriptsGetLogsOptions` has been removed
- Struct `DeploymentScriptsGetLogsResponse` has been removed
- Struct `DeploymentScriptsGetLogsResult` has been removed
- Struct `DeploymentScriptsGetOptions` has been removed
- Struct `DeploymentScriptsGetResponse` has been removed
- Struct `DeploymentScriptsGetResult` has been removed
- Struct `DeploymentScriptsListByResourceGroupOptions` has been removed
- Struct `DeploymentScriptsListByResourceGroupPager` has been removed
- Struct `DeploymentScriptsListByResourceGroupResponse` has been removed
- Struct `DeploymentScriptsListByResourceGroupResult` has been removed
- Struct `DeploymentScriptsListBySubscriptionOptions` has been removed
- Struct `DeploymentScriptsListBySubscriptionPager` has been removed
- Struct `DeploymentScriptsListBySubscriptionResponse` has been removed
- Struct `DeploymentScriptsListBySubscriptionResult` has been removed
- Struct `DeploymentScriptsUpdateOptions` has been removed
- Struct `DeploymentScriptsUpdateResponse` has been removed
- Struct `DeploymentScriptsUpdateResult` has been removed
- Field `DeploymentScriptPropertiesBase` of struct `AzurePowerShellScriptProperties` has been removed
- Field `ScriptConfigurationBase` of struct `AzurePowerShellScriptProperties` has been removed
- Field `AzureResourceBase` of struct `DeploymentScriptUpdateParameter` has been removed
- Field `DeploymentScript` of struct `AzureCliScript` has been removed
- Field `AzureResourceBase` of struct `DeploymentScript` has been removed
- Field `AzureResourceBase` of struct `ScriptLog` has been removed
- Field `DeploymentScriptPropertiesBase` of struct `AzureCliScriptProperties` has been removed
- Field `ScriptConfigurationBase` of struct `AzureCliScriptProperties` has been removed
- Field `DeploymentScript` of struct `AzurePowerShellScript` has been removed

### Features Added

- New function `*Client.ListByResourceGroup(string, *ClientListByResourceGroupOptions) *ClientListByResourceGroupPager`
- New function `*ClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `NewClient(string, azcore.TokenCredential, *arm.ClientOptions) *Client`
- New function `ClientCreatePollerResponse.PollUntilDone(context.Context, time.Duration) (ClientCreateResponse, error)`
- New function `*ClientCreatePollerResponse.Resume(context.Context, *Client, string) error`
- New function `*Client.GetLogsDefault(context.Context, string, string, *ClientGetLogsDefaultOptions) (ClientGetLogsDefaultResponse, error)`
- New function `*ClientCreatePoller.FinalResponse(context.Context) (ClientCreateResponse, error)`
- New function `*ClientCreateResult.UnmarshalJSON([]byte) error`
- New function `*ClientCreatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*ClientListBySubscriptionPager.Err() error`
- New function `*ClientCreatePoller.Done() bool`
- New function `*ClientGetResult.UnmarshalJSON([]byte) error`
- New function `*Client.Update(context.Context, string, string, *ClientUpdateOptions) (ClientUpdateResponse, error)`
- New function `*ClientListBySubscriptionPager.PageResponse() ClientListBySubscriptionResponse`
- New function `DeploymentScript.MarshalJSON() ([]byte, error)`
- New function `*AzurePowerShellScript.GetDeploymentScript() *DeploymentScript`
- New function `*ClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*ClientCreatePoller.ResumeToken() (string, error)`
- New function `*ClientListByResourceGroupPager.Err() error`
- New function `*Client.Get(context.Context, string, string, *ClientGetOptions) (ClientGetResponse, error)`
- New function `*ClientUpdateResult.UnmarshalJSON([]byte) error`
- New function `*Client.ListBySubscription(*ClientListBySubscriptionOptions) *ClientListBySubscriptionPager`
- New function `*Client.Delete(context.Context, string, string, *ClientDeleteOptions) (ClientDeleteResponse, error)`
- New function `*ClientListByResourceGroupPager.PageResponse() ClientListByResourceGroupResponse`
- New function `*Client.GetLogs(context.Context, string, string, *ClientGetLogsOptions) (ClientGetLogsResponse, error)`
- New function `*Client.BeginCreate(context.Context, string, string, DeploymentScriptClassification, *ClientBeginCreateOptions) (ClientCreatePollerResponse, error)`
- New function `*AzureCliScript.GetDeploymentScript() *DeploymentScript`
- New struct `Client`
- New struct `ClientBeginCreateOptions`
- New struct `ClientCreatePoller`
- New struct `ClientCreatePollerResponse`
- New struct `ClientCreateResponse`
- New struct `ClientCreateResult`
- New struct `ClientDeleteOptions`
- New struct `ClientDeleteResponse`
- New struct `ClientGetLogsDefaultOptions`
- New struct `ClientGetLogsDefaultResponse`
- New struct `ClientGetLogsDefaultResult`
- New struct `ClientGetLogsOptions`
- New struct `ClientGetLogsResponse`
- New struct `ClientGetLogsResult`
- New struct `ClientGetOptions`
- New struct `ClientGetResponse`
- New struct `ClientGetResult`
- New struct `ClientListByResourceGroupOptions`
- New struct `ClientListByResourceGroupPager`
- New struct `ClientListByResourceGroupResponse`
- New struct `ClientListByResourceGroupResult`
- New struct `ClientListBySubscriptionOptions`
- New struct `ClientListBySubscriptionPager`
- New struct `ClientListBySubscriptionResponse`
- New struct `ClientListBySubscriptionResult`
- New struct `ClientUpdateOptions`
- New struct `ClientUpdateResponse`
- New struct `ClientUpdateResult`
- New struct `Error`
- New field `Tags` in struct `AzureCliScript`
- New field `SystemData` in struct `AzureCliScript`
- New field `Location` in struct `AzureCliScript`
- New field `ID` in struct `AzureCliScript`
- New field `Name` in struct `AzureCliScript`
- New field `Type` in struct `AzureCliScript`
- New field `Kind` in struct `AzureCliScript`
- New field `Identity` in struct `AzureCliScript`
- New field `Location` in struct `AzurePowerShellScript`
- New field `Kind` in struct `AzurePowerShellScript`
- New field `SystemData` in struct `AzurePowerShellScript`
- New field `Identity` in struct `AzurePowerShellScript`
- New field `Tags` in struct `AzurePowerShellScript`
- New field `ID` in struct `AzurePowerShellScript`
- New field `Name` in struct `AzurePowerShellScript`
- New field `Type` in struct `AzurePowerShellScript`
- New field `RetentionInterval` in struct `AzureCliScriptProperties`
- New field `Arguments` in struct `AzureCliScriptProperties`
- New field `ScriptContent` in struct `AzureCliScriptProperties`
- New field `EnvironmentVariables` in struct `AzureCliScriptProperties`
- New field `ForceUpdateTag` in struct `AzureCliScriptProperties`
- New field `Timeout` in struct `AzureCliScriptProperties`
- New field `ContainerSettings` in struct `AzureCliScriptProperties`
- New field `ProvisioningState` in struct `AzureCliScriptProperties`
- New field `Outputs` in struct `AzureCliScriptProperties`
- New field `CleanupPreference` in struct `AzureCliScriptProperties`
- New field `Status` in struct `AzureCliScriptProperties`
- New field `SupportingScriptUris` in struct `AzureCliScriptProperties`
- New field `PrimaryScriptURI` in struct `AzureCliScriptProperties`
- New field `StorageAccountSettings` in struct `AzureCliScriptProperties`
- New field `StorageAccountSettings` in struct `AzurePowerShellScriptProperties`
- New field `Arguments` in struct `AzurePowerShellScriptProperties`
- New field `CleanupPreference` in struct `AzurePowerShellScriptProperties`
- New field `RetentionInterval` in struct `AzurePowerShellScriptProperties`
- New field `ForceUpdateTag` in struct `AzurePowerShellScriptProperties`
- New field `Outputs` in struct `AzurePowerShellScriptProperties`
- New field `ContainerSettings` in struct `AzurePowerShellScriptProperties`
- New field `SupportingScriptUris` in struct `AzurePowerShellScriptProperties`
- New field `Status` in struct `AzurePowerShellScriptProperties`
- New field `EnvironmentVariables` in struct `AzurePowerShellScriptProperties`
- New field `Timeout` in struct `AzurePowerShellScriptProperties`
- New field `ProvisioningState` in struct `AzurePowerShellScriptProperties`
- New field `ScriptContent` in struct `AzurePowerShellScriptProperties`
- New field `PrimaryScriptURI` in struct `AzurePowerShellScriptProperties`
- New field `ID` in struct `DeploymentScriptUpdateParameter`
- New field `Name` in struct `DeploymentScriptUpdateParameter`
- New field `Type` in struct `DeploymentScriptUpdateParameter`
- New field `Type` in struct `ScriptLog`
- New field `ID` in struct `ScriptLog`
- New field `Name` in struct `ScriptLog`
- New field `ID` in struct `DeploymentScript`
- New field `Name` in struct `DeploymentScript`
- New field `Type` in struct `DeploymentScript`


## 0.1.1 (2021-12-13)

### Other Changes

- Fix the go minimum version to `1.16`

## 0.1.0 (2021-11-16)

- Initial preview release.
