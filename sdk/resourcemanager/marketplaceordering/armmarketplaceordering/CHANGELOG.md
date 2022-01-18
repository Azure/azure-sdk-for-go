# Release History

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*MarketplaceAgreementsClient.List` parameter(s) have been changed from `(context.Context, *MarketplaceAgreementsListOptions)` to `(context.Context, *MarketplaceAgreementsClientListOptions)`
- Function `*MarketplaceAgreementsClient.List` return value(s) have been changed from `(MarketplaceAgreementsListResponse, error)` to `(MarketplaceAgreementsClientListResponse, error)`
- Function `*MarketplaceAgreementsClient.GetAgreement` parameter(s) have been changed from `(context.Context, string, string, string, *MarketplaceAgreementsGetAgreementOptions)` to `(context.Context, string, string, string, *MarketplaceAgreementsClientGetAgreementOptions)`
- Function `*MarketplaceAgreementsClient.GetAgreement` return value(s) have been changed from `(MarketplaceAgreementsGetAgreementResponse, error)` to `(MarketplaceAgreementsClientGetAgreementResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*MarketplaceAgreementsClient.Sign` parameter(s) have been changed from `(context.Context, string, string, string, *MarketplaceAgreementsSignOptions)` to `(context.Context, string, string, string, *MarketplaceAgreementsClientSignOptions)`
- Function `*MarketplaceAgreementsClient.Sign` return value(s) have been changed from `(MarketplaceAgreementsSignResponse, error)` to `(MarketplaceAgreementsClientSignResponse, error)`
- Function `*MarketplaceAgreementsClient.Cancel` parameter(s) have been changed from `(context.Context, string, string, string, *MarketplaceAgreementsCancelOptions)` to `(context.Context, string, string, string, *MarketplaceAgreementsClientCancelOptions)`
- Function `*MarketplaceAgreementsClient.Cancel` return value(s) have been changed from `(MarketplaceAgreementsCancelResponse, error)` to `(MarketplaceAgreementsClientCancelResponse, error)`
- Function `*MarketplaceAgreementsClient.Create` parameter(s) have been changed from `(context.Context, OfferType, string, string, string, AgreementTerms, *MarketplaceAgreementsCreateOptions)` to `(context.Context, OfferType, string, string, string, AgreementTerms, *MarketplaceAgreementsClientCreateOptions)`
- Function `*MarketplaceAgreementsClient.Create` return value(s) have been changed from `(MarketplaceAgreementsCreateResponse, error)` to `(MarketplaceAgreementsClientCreateResponse, error)`
- Function `*MarketplaceAgreementsClient.Get` parameter(s) have been changed from `(context.Context, OfferType, string, string, string, *MarketplaceAgreementsGetOptions)` to `(context.Context, OfferType, string, string, string, *MarketplaceAgreementsClientGetOptions)`
- Function `*MarketplaceAgreementsClient.Get` return value(s) have been changed from `(MarketplaceAgreementsGetResponse, error)` to `(MarketplaceAgreementsClientGetResponse, error)`
- Function `*OperationsListPager.PageResponse` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Struct `MarketplaceAgreementsCancelOptions` has been removed
- Struct `MarketplaceAgreementsCancelResponse` has been removed
- Struct `MarketplaceAgreementsCancelResult` has been removed
- Struct `MarketplaceAgreementsCreateOptions` has been removed
- Struct `MarketplaceAgreementsCreateResponse` has been removed
- Struct `MarketplaceAgreementsCreateResult` has been removed
- Struct `MarketplaceAgreementsGetAgreementOptions` has been removed
- Struct `MarketplaceAgreementsGetAgreementResponse` has been removed
- Struct `MarketplaceAgreementsGetAgreementResult` has been removed
- Struct `MarketplaceAgreementsGetOptions` has been removed
- Struct `MarketplaceAgreementsGetResponse` has been removed
- Struct `MarketplaceAgreementsGetResult` has been removed
- Struct `MarketplaceAgreementsListOptions` has been removed
- Struct `MarketplaceAgreementsListResponse` has been removed
- Struct `MarketplaceAgreementsListResult` has been removed
- Struct `MarketplaceAgreementsSignOptions` has been removed
- Struct `MarketplaceAgreementsSignResponse` has been removed
- Struct `MarketplaceAgreementsSignResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `Resource` of struct `AgreementTerms` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `*OperationsClientListPager.Err() error`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New struct `MarketplaceAgreementsClientCancelOptions`
- New struct `MarketplaceAgreementsClientCancelResponse`
- New struct `MarketplaceAgreementsClientCancelResult`
- New struct `MarketplaceAgreementsClientCreateOptions`
- New struct `MarketplaceAgreementsClientCreateResponse`
- New struct `MarketplaceAgreementsClientCreateResult`
- New struct `MarketplaceAgreementsClientGetAgreementOptions`
- New struct `MarketplaceAgreementsClientGetAgreementResponse`
- New struct `MarketplaceAgreementsClientGetAgreementResult`
- New struct `MarketplaceAgreementsClientGetOptions`
- New struct `MarketplaceAgreementsClientGetResponse`
- New struct `MarketplaceAgreementsClientGetResult`
- New struct `MarketplaceAgreementsClientListOptions`
- New struct `MarketplaceAgreementsClientListResponse`
- New struct `MarketplaceAgreementsClientListResult`
- New struct `MarketplaceAgreementsClientSignOptions`
- New struct `MarketplaceAgreementsClientSignResponse`
- New struct `MarketplaceAgreementsClientSignResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `ID` in struct `AgreementTerms`
- New field `Name` in struct `AgreementTerms`
- New field `Type` in struct `AgreementTerms`
- New field `Error` in struct `ErrorResponse`


## 0.1.0 (2021-12-07)

- Init release.