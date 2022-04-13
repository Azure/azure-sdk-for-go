# Release History

## 0.3.0 (2022-04-12)
### Breaking Changes

- Function `*OperationClient.List` return value(s) have been changed from `(*OperationClientListPager)` to `(*runtime.Pager[OperationClientListResponse])`
- Function `NewPartnersClient` return value(s) have been changed from `(*PartnersClient)` to `(*PartnersClient, error)`
- Function `NewPartnerClient` return value(s) have been changed from `(*PartnerClient)` to `(*PartnerClient, error)`
- Function `NewOperationClient` return value(s) have been changed from `(*OperationClient)` to `(*OperationClient, error)`
- Function `*OperationClientListPager.Err` has been removed
- Function `*OperationClientListPager.NextPage` has been removed
- Function `ManagementPartnerState.ToPtr` has been removed
- Function `*OperationClientListPager.PageResponse` has been removed
- Struct `OperationClientListPager` has been removed
- Struct `OperationClientListResult` has been removed
- Struct `PartnerClientCreateResult` has been removed
- Struct `PartnerClientGetResult` has been removed
- Struct `PartnerClientUpdateResult` has been removed
- Struct `PartnersClientGetResult` has been removed
- Field `PartnerClientCreateResult` of struct `PartnerClientCreateResponse` has been removed
- Field `RawResponse` of struct `PartnerClientCreateResponse` has been removed
- Field `PartnersClientGetResult` of struct `PartnersClientGetResponse` has been removed
- Field `RawResponse` of struct `PartnersClientGetResponse` has been removed
- Field `PartnerClientGetResult` of struct `PartnerClientGetResponse` has been removed
- Field `RawResponse` of struct `PartnerClientGetResponse` has been removed
- Field `RawResponse` of struct `PartnerClientDeleteResponse` has been removed
- Field `PartnerClientUpdateResult` of struct `PartnerClientUpdateResponse` has been removed
- Field `RawResponse` of struct `PartnerClientUpdateResponse` has been removed
- Field `OperationClientListResult` of struct `OperationClientListResponse` has been removed
- Field `RawResponse` of struct `OperationClientListResponse` has been removed

### Features Added

- New anonymous field `OperationList` in struct `OperationClientListResponse`
- New anonymous field `PartnerResponse` in struct `PartnerClientUpdateResponse`
- New anonymous field `PartnerResponse` in struct `PartnersClientGetResponse`
- New anonymous field `PartnerResponse` in struct `PartnerClientCreateResponse`
- New anonymous field `PartnerResponse` in struct `PartnerClientGetResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-17)
### Breaking Changes

- Function `*PartnersClient.Get` parameter(s) have been changed from `(context.Context, *PartnersGetOptions)` to `(context.Context, *PartnersClientGetOptions)`
- Function `*PartnersClient.Get` return value(s) have been changed from `(PartnersGetResponse, error)` to `(PartnersClientGetResponse, error)`
- Function `*OperationClient.List` parameter(s) have been changed from `(*OperationListOptions)` to `(*OperationClientListOptions)`
- Function `*OperationClient.List` return value(s) have been changed from `(*OperationListPager)` to `(*OperationClientListPager)`
- Function `*PartnerClient.Get` parameter(s) have been changed from `(context.Context, string, *PartnerGetOptions)` to `(context.Context, string, *PartnerClientGetOptions)`
- Function `*PartnerClient.Get` return value(s) have been changed from `(PartnerGetResponse, error)` to `(PartnerClientGetResponse, error)`
- Function `*PartnerClient.Create` parameter(s) have been changed from `(context.Context, string, *PartnerCreateOptions)` to `(context.Context, string, *PartnerClientCreateOptions)`
- Function `*PartnerClient.Create` return value(s) have been changed from `(PartnerCreateResponse, error)` to `(PartnerClientCreateResponse, error)`
- Function `*PartnerClient.Update` parameter(s) have been changed from `(context.Context, string, *PartnerUpdateOptions)` to `(context.Context, string, *PartnerClientUpdateOptions)`
- Function `*PartnerClient.Update` return value(s) have been changed from `(PartnerUpdateResponse, error)` to `(PartnerClientUpdateResponse, error)`
- Function `*PartnerClient.Delete` parameter(s) have been changed from `(context.Context, string, *PartnerDeleteOptions)` to `(context.Context, string, *PartnerClientDeleteOptions)`
- Function `*PartnerClient.Delete` return value(s) have been changed from `(PartnerDeleteResponse, error)` to `(PartnerClientDeleteResponse, error)`
- Type of `ExtendedErrorInfo.Code` has been changed from `*ErrorResponseCode` to `*string`
- Const `ErrorResponseCodeConflict` has been removed
- Const `ErrorResponseCodeNotFound` has been removed
- Const `ErrorResponseCodeBadRequest` has been removed
- Function `*OperationListPager.PageResponse` has been removed
- Function `ErrorResponseCode.ToPtr` has been removed
- Function `*OperationListPager.NextPage` has been removed
- Function `*OperationListPager.Err` has been removed
- Function `Error.Error` has been removed
- Function `PossibleErrorResponseCodeValues` has been removed
- Struct `OperationListOptions` has been removed
- Struct `OperationListPager` has been removed
- Struct `OperationListResponse` has been removed
- Struct `OperationListResult` has been removed
- Struct `PartnerCreateOptions` has been removed
- Struct `PartnerCreateResponse` has been removed
- Struct `PartnerCreateResult` has been removed
- Struct `PartnerDeleteOptions` has been removed
- Struct `PartnerDeleteResponse` has been removed
- Struct `PartnerGetOptions` has been removed
- Struct `PartnerGetResponse` has been removed
- Struct `PartnerGetResult` has been removed
- Struct `PartnerUpdateOptions` has been removed
- Struct `PartnerUpdateResponse` has been removed
- Struct `PartnerUpdateResult` has been removed
- Struct `PartnersGetOptions` has been removed
- Struct `PartnersGetResponse` has been removed
- Struct `PartnersGetResult` has been removed
- Field `InnerError` of struct `Error` has been removed

### Features Added

- New function `*OperationClientListPager.Err() error`
- New function `*OperationClientListPager.NextPage(context.Context) bool`
- New function `*OperationClientListPager.PageResponse() OperationClientListResponse`
- New struct `OperationClientListOptions`
- New struct `OperationClientListPager`
- New struct `OperationClientListResponse`
- New struct `OperationClientListResult`
- New struct `PartnerClientCreateOptions`
- New struct `PartnerClientCreateResponse`
- New struct `PartnerClientCreateResult`
- New struct `PartnerClientDeleteOptions`
- New struct `PartnerClientDeleteResponse`
- New struct `PartnerClientGetOptions`
- New struct `PartnerClientGetResponse`
- New struct `PartnerClientGetResult`
- New struct `PartnerClientUpdateOptions`
- New struct `PartnerClientUpdateResponse`
- New struct `PartnerClientUpdateResult`
- New struct `PartnersClientGetOptions`
- New struct `PartnersClientGetResponse`
- New struct `PartnersClientGetResult`
- New field `Code` in struct `Error`
- New field `Error` in struct `Error`
- New field `Message` in struct `Error`


## 0.1.0 (2021-12-07)

- Init release.