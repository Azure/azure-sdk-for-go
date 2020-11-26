
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Const `Failed` type has been changed from `ProvisioningStatus` to `ProvisioningState`
- Const `Created` type has been changed from `ProvisioningStatus` to `ProvisioningState`
- Const `Succeeded` type has been changed from `ProvisioningStatus` to `ProvisioningState`
- Function `AccountsClient.Update` signature has been changed from `(context.Context,string,string,Account)` to `(context.Context,string,string,AccountUpdate)`
- Function `ConfigurationProfilePreferencesClient.Update` signature has been changed from `(context.Context,string,string,ConfigurationProfilePreference)` to `(context.Context,string,string,ConfigurationProfilePreferenceUpdate)`
- Function `ConfigurationProfilePreferencesClient.UpdatePreparer` signature has been changed from `(context.Context,string,string,ConfigurationProfilePreference)` to `(context.Context,string,string,ConfigurationProfilePreferenceUpdate)`
- Function `AccountsClient.UpdatePreparer` signature has been changed from `(context.Context,string,string,Account)` to `(context.Context,string,string,AccountUpdate)`
- Type of `ErrorResponse.Error` has been changed from `*ErrorResponseBody` to `*ErrorDetail`
- Const `AzureBestPracticesTestDev` has been removed
- Const `AzureBestPracticesProd` has been removed
- Function `PossibleProvisioningStatusValues` has been removed
- Struct `ErrorResponseBody` has been removed
- Field `Location` of struct `ProxyResource` has been removed
- Field `ProvisioningStatus` of struct `ConfigurationProfileAssignmentProperties` has been removed
- Field `Location` of struct `ConfigurationProfileAssignment` has been removed
- Field `Location` of struct `Resource` has been removed

## New Content

- Const `AzurevirtualmachinebestpracticesDevTest` is added
- Const `AzurevirtualmachinebestpracticesProduction` is added
- Function `UpdateResource.MarshalJSON() ([]byte,error)` is added
- Function `PossibleProvisioningStateValues() []ProvisioningState` is added
- Function `AccountUpdate.MarshalJSON() ([]byte,error)` is added
- Function `ConfigurationProfilePreferenceUpdate.MarshalJSON() ([]byte,error)` is added
- Struct `AccountUpdate` is added
- Struct `AzureEntityResource` is added
- Struct `ConfigurationProfilePreferenceUpdate` is added
- Struct `ErrorAdditionalInfo` is added
- Struct `ErrorDetail` is added
- Struct `UpdateResource` is added
- Field `ProvisioningState` is added to struct `ConfigurationProfileAssignmentProperties`

