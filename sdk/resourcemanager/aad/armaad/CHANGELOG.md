# Release History

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*DiagnosticSettingsClient.Get` parameter(s) have been changed from `(context.Context, string, *DiagnosticSettingsGetOptions)` to `(context.Context, string, *DiagnosticSettingsClientGetOptions)`
- Function `*DiagnosticSettingsClient.Get` return value(s) have been changed from `(DiagnosticSettingsGetResponse, error)` to `(DiagnosticSettingsClientGetResponse, error)`
- Function `*DiagnosticSettingsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, DiagnosticSettingsResource, *DiagnosticSettingsCreateOrUpdateOptions)` to `(context.Context, string, DiagnosticSettingsResource, *DiagnosticSettingsClientCreateOrUpdateOptions)`
- Function `*DiagnosticSettingsClient.CreateOrUpdate` return value(s) have been changed from `(DiagnosticSettingsCreateOrUpdateResponse, error)` to `(DiagnosticSettingsClientCreateOrUpdateResponse, error)`
- Function `*DiagnosticSettingsCategoryClient.List` parameter(s) have been changed from `(context.Context, *DiagnosticSettingsCategoryListOptions)` to `(context.Context, *DiagnosticSettingsCategoryClientListOptions)`
- Function `*DiagnosticSettingsCategoryClient.List` return value(s) have been changed from `(DiagnosticSettingsCategoryListResponse, error)` to `(DiagnosticSettingsCategoryClientListResponse, error)`
- Function `*DiagnosticSettingsClient.Delete` parameter(s) have been changed from `(context.Context, string, *DiagnosticSettingsDeleteOptions)` to `(context.Context, string, *DiagnosticSettingsClientDeleteOptions)`
- Function `*DiagnosticSettingsClient.Delete` return value(s) have been changed from `(DiagnosticSettingsDeleteResponse, error)` to `(DiagnosticSettingsClientDeleteResponse, error)`
- Function `*DiagnosticSettingsClient.List` parameter(s) have been changed from `(context.Context, *DiagnosticSettingsListOptions)` to `(context.Context, *DiagnosticSettingsClientListOptions)`
- Function `*DiagnosticSettingsClient.List` return value(s) have been changed from `(DiagnosticSettingsListResponse, error)` to `(DiagnosticSettingsClientListResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Function `ErrorResponse.Error` has been removed
- Struct `DiagnosticSettingsCategoryListOptions` has been removed
- Struct `DiagnosticSettingsCategoryListResponse` has been removed
- Struct `DiagnosticSettingsCategoryListResult` has been removed
- Struct `DiagnosticSettingsCreateOrUpdateOptions` has been removed
- Struct `DiagnosticSettingsCreateOrUpdateResponse` has been removed
- Struct `DiagnosticSettingsCreateOrUpdateResult` has been removed
- Struct `DiagnosticSettingsDeleteOptions` has been removed
- Struct `DiagnosticSettingsDeleteResponse` has been removed
- Struct `DiagnosticSettingsGetOptions` has been removed
- Struct `DiagnosticSettingsGetResponse` has been removed
- Struct `DiagnosticSettingsGetResult` has been removed
- Struct `DiagnosticSettingsListOptions` has been removed
- Struct `DiagnosticSettingsListResponse` has been removed
- Struct `DiagnosticSettingsListResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `ProxyOnlyResource` of struct `DiagnosticSettingsCategoryResource` has been removed
- Field `ProxyOnlyResource` of struct `DiagnosticSettingsResource` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New struct `DiagnosticSettingsCategoryClientListOptions`
- New struct `DiagnosticSettingsCategoryClientListResponse`
- New struct `DiagnosticSettingsCategoryClientListResult`
- New struct `DiagnosticSettingsClientCreateOrUpdateOptions`
- New struct `DiagnosticSettingsClientCreateOrUpdateResponse`
- New struct `DiagnosticSettingsClientCreateOrUpdateResult`
- New struct `DiagnosticSettingsClientDeleteOptions`
- New struct `DiagnosticSettingsClientDeleteResponse`
- New struct `DiagnosticSettingsClientGetOptions`
- New struct `DiagnosticSettingsClientGetResponse`
- New struct `DiagnosticSettingsClientGetResult`
- New struct `DiagnosticSettingsClientListOptions`
- New struct `DiagnosticSettingsClientListResponse`
- New struct `DiagnosticSettingsClientListResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `Error` in struct `ErrorResponse`
- New field `ID` in struct `DiagnosticSettingsCategoryResource`
- New field `Name` in struct `DiagnosticSettingsCategoryResource`
- New field `Type` in struct `DiagnosticSettingsCategoryResource`
- New field `Name` in struct `DiagnosticSettingsResource`
- New field `Type` in struct `DiagnosticSettingsResource`
- New field `ID` in struct `DiagnosticSettingsResource`


## 0.1.0 (2021-11-30)

- Initial preview release.
