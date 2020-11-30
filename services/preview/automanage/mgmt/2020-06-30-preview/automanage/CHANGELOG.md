Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Const `Failed` type has been changed from `ProvisioningStatus` to `ProvisioningState`
- Const `Succeeded` type has been changed from `ProvisioningStatus` to `ProvisioningState`
- Const `Created` type has been changed from `ProvisioningStatus` to `ProvisioningState`
- Function `AccountsClient.UpdatePreparer` parameter(s) have been changed from `(context.Context, string, string, Account)` to `(context.Context, string, string, AccountUpdate)`
- Function `ConfigurationProfilePreferencesClient.Update` parameter(s) have been changed from `(context.Context, string, string, ConfigurationProfilePreference)` to `(context.Context, string, string, ConfigurationProfilePreferenceUpdate)`
- Function `AccountsClient.Update` parameter(s) have been changed from `(context.Context, string, string, Account)` to `(context.Context, string, string, AccountUpdate)`
- Function `ConfigurationProfilePreferencesClient.UpdatePreparer` parameter(s) have been changed from `(context.Context, string, string, ConfigurationProfilePreference)` to `(context.Context, string, string, ConfigurationProfilePreferenceUpdate)`
- Type of `ErrorResponse.Error` has been changed from `*ErrorResponseBody` to `*ErrorDetail`
- Const `AzureBestPracticesProd` has been removed
- Const `AzureBestPracticesTestDev` has been removed
- Function `PossibleProvisioningStatusValues` has been removed
- Struct `ErrorResponseBody` has been removed
- Field `ProvisioningStatus` of struct `ConfigurationProfileAssignmentProperties` has been removed
- Field `Location` of struct `Resource` has been removed
- Field `Location` of struct `ConfigurationProfileAssignment` has been removed
- Field `Location` of struct `ProxyResource` has been removed

## New Content

- New const `AzurevirtualmachinebestpracticesProduction`
- New const `AzurevirtualmachinebestpracticesDevTest`
- New function `ConfigurationProfilePreferenceUpdate.MarshalJSON() ([]byte, error)`
- New function `UpdateResource.MarshalJSON() ([]byte, error)`
- New function `AccountUpdate.MarshalJSON() ([]byte, error)`
- New function `PossibleProvisioningStateValues() []ProvisioningState`
- New struct `AccountUpdate`
- New struct `AzureEntityResource`
- New struct `ConfigurationProfilePreferenceUpdate`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `UpdateResource`
- New field `ProvisioningState` in struct `ConfigurationProfileAssignmentProperties`
