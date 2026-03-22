# Release History

## 4.0.0 (2026-03-12)
### Breaking Changes

- Function `NewClientFactory` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Function `NewQuotaClient` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Function `*QuotaClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaClientBeginCreateOrUpdateOptions)`
- Function `*QuotaClient.BeginUpdate` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaClientBeginUpdateOptions)` to `(ctx context.Context, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaClientBeginUpdateOptions)`
- Function `*QuotaClient.Get` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, options *QuotaClientGetOptions)` to `(ctx context.Context, providerID string, location string, resourceName string, options *QuotaClientGetOptions)`
- Function `*QuotaClient.NewListPager` parameter(s) have been changed from `(subscriptionID string, providerID string, location string, options *QuotaClientListOptions)` to `(providerID string, location string, options *QuotaClientListOptions)`
- Function `NewQuotaRequestStatusClient` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Function `*QuotaRequestStatusClient.Get` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, providerID string, location string, id string, options *QuotaRequestStatusClientGetOptions)` to `(ctx context.Context, providerID string, location string, id string, options *QuotaRequestStatusClientGetOptions)`
- Function `*QuotaRequestStatusClient.NewListPager` parameter(s) have been changed from `(subscriptionID string, providerID string, location string, options *QuotaRequestStatusClientListOptions)` to `(providerID string, location string, options *QuotaRequestStatusClientListOptions)`
- Enum `DisplayProvisioningState` has been removed
- Enum `Location` has been removed
- Enum `UserFriendlyAppliedScopeType` has been removed
- Enum `UserFriendlyRenewState` has been removed
- Function `NewAzureReservationAPIClient` has been removed
- Function `*AzureReservationAPIClient.GetAppliedReservationList` has been removed
- Function `*AzureReservationAPIClient.NewGetCatalogPager` has been removed
- Function `*ClientFactory.NewAzureReservationAPIClient` has been removed
- Struct `CreateGenericQuotaRequestParameters` has been removed
- Struct `CurrentQuotaLimit` has been removed
- Struct `Error` has been removed
- Struct `ErrorDetails` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ExceptionResponse` has been removed
- Struct `ExtendedErrorInfo` has been removed
- Struct `ProxyResource` has been removed
- Struct `QuotaLimitsResponse` has been removed
- Struct `QuotaRequestOneResourceProperties` has been removed
- Struct `QuotaRequestOneResourceSubmitResponse` has been removed
- Struct `QuotaRequestStatusDetails` has been removed
- Struct `QuotaRequestSubmitResponse` has been removed
- Struct `QuotaRequestSubmitResponse201` has been removed
- Struct `RefundResponse` has been removed
- Struct `Resource` has been removed
- Struct `ServiceError` has been removed
- Struct `ServiceErrorDetail` has been removed
- Field `ReservationResponseArray` of struct `ReservationClientMergeResponse` has been removed
- Field `ReservationResponseArray` of struct `ReservationClientSplitResponse` has been removed

### Features Added

- New function `NewClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*Client, error)`
- New function `*Client.GetAppliedReservationList(ctx context.Context, options *ClientGetAppliedReservationListOptions) (ClientGetAppliedReservationListResponse, error)`
- New function `*Client.NewGetCatalogPager(options *ClientGetCatalogOptions) *runtime.Pager[ClientGetCatalogResponse]`
- New function `*Client.NewCalculateExchangeClient() *CalculateExchangeClient`
- New function `*Client.NewCalculateRefundClient() *CalculateRefundClient`
- New function `*Client.NewExchangeClient() *ExchangeClient`
- New function `*Client.NewOperationClient() *OperationClient`
- New function `*Client.NewQuotaClient() *QuotaClient`
- New function `*Client.NewQuotaRequestStatusClient() *QuotaRequestStatusClient`
- New function `*Client.NewReservationClient() *ReservationClient`
- New function `*Client.NewReservationOrderClient() *ReservationOrderClient`
- New function `*Client.NewReturnClient() *ReturnClient`
- New function `*ClientFactory.NewClient() *Client`
- New field `SystemData` in struct `CurrentQuotaLimitBase`
- New field `SystemData` in struct `QuotaRequestDetails`


## 3.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 3.0.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 3.0.0 (2023-03-24)
### Breaking Changes

- Response of `ReturnClient.Post` has been changed from `RefundResponse` to `ReservationOrderResponse`

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New field `P3Y` in struct `CatalogMsrp`
- New field `P5Y` in struct `CatalogMsrp`
- New anonymous field `ReservationOrderResponse` in struct `ReturnClientPostResponse`


## 2.0.0 (2023-02-24)
### Breaking Changes

- Operation `*AzureReservationAPIClient.GetCatalog` has supported pagination, use `*AzureReservationAPIClient.NewGetCatalogPager` instead.
- Operation `*ReturnClient.Post` has been changed to LRO, use `*ReturnClient.BeginPost` instead.

### Features Added

- New value `AppliedScopeTypeManagementGroup` added to type alias `AppliedScopeType`
- New value `DisplayProvisioningStateNoBenefit`, `DisplayProvisioningStateWarning` added to type alias `DisplayProvisioningState`
- New type alias `BillingPlan` with values `BillingPlanP1M`
- New type alias `CommitmentGrain` with values `CommitmentGrainHourly`
- New type alias `SavingsPlanTerm` with values `SavingsPlanTermP1Y`, `SavingsPlanTermP3Y`
- New struct `AppliedScopeProperties`
- New struct `CatalogsResult`
- New struct `Commitment`
- New struct `ProxyResource`
- New struct `ReservationSwapProperties`
- New struct `Resource`
- New struct `SavingsPlanPurchaseRequest`
- New struct `SavingsPlanPurchaseRequestProperties`
- New struct `SavingsPlanToPurchaseCalculateExchange`
- New struct `SavingsPlanToPurchaseExchange`
- New anonymous field `CatalogsResult` in struct `AzureReservationAPIClientGetCatalogResponse`
- New field `SavingsPlansToPurchase` in struct `CalculateExchangeRequestProperties`
- New field `SavingsPlansToPurchase` in struct `CalculateExchangeResponseProperties`
- New field `SavingsPlansToPurchase` in struct `ExchangeResponseProperties`
- New field `AppliedScopeProperties` in struct `PatchProperties`
- New field `ReviewDateTime` in struct `PatchProperties`
- New field `AppliedScopeProperties` in struct `Properties`
- New field `ExpiryDateTime` in struct `Properties`
- New field `PurchaseDateTime` in struct `Properties`
- New field `ReviewDateTime` in struct `Properties`
- New field `SwapProperties` in struct `Properties`
- New field `AppliedScopeProperties` in struct `PurchaseRequestProperties`
- New field `ReviewDateTime` in struct `PurchaseRequestProperties`
- New field `ExpiryDateTime` in struct `ReservationOrderProperties`
- New field `ReviewDateTime` in struct `ReservationOrderProperties`
- New field `NoBenefitCount` in struct `ReservationSummary`
- New field `WarningCount` in struct `ReservationSummary`


## 1.1.0 (2022-09-16)
### Features Added

- New const `ErrorResponseCodeRefundLimitExceeded`
- New const `ErrorResponseCodeSelfServiceRefundNotSupported`
- New function `*ReservationClient.Archive(context.Context, string, string, *ReservationClientArchiveOptions) (ReservationClientArchiveResponse, error)`
- New function `NewCalculateRefundClient(azcore.TokenCredential, *arm.ClientOptions) (*CalculateRefundClient, error)`
- New function `*CalculateRefundClient.Post(context.Context, string, CalculateRefundRequest, *CalculateRefundClientPostOptions) (CalculateRefundClientPostResponse, error)`
- New function `*ReservationClient.Unarchive(context.Context, string, string, *ReservationClientUnarchiveOptions) (ReservationClientUnarchiveResponse, error)`
- New function `NewReturnClient(azcore.TokenCredential, *arm.ClientOptions) (*ReturnClient, error)`
- New function `*ReturnClient.Post(context.Context, string, RefundRequest, *ReturnClientPostOptions) (ReturnClientPostResponse, error)`
- New struct `CalculateRefundClient`
- New struct `CalculateRefundClientPostOptions`
- New struct `CalculateRefundClientPostResponse`
- New struct `CalculateRefundRequest`
- New struct `CalculateRefundRequestProperties`
- New struct `CalculateRefundResponse`
- New struct `RefundBillingInformation`
- New struct `RefundPolicyError`
- New struct `RefundPolicyResult`
- New struct `RefundPolicyResultProperty`
- New struct `RefundRequest`
- New struct `RefundRequestProperties`
- New struct `RefundResponse`
- New struct `RefundResponseProperties`
- New struct `ReservationClientArchiveOptions`
- New struct `ReservationClientArchiveResponse`
- New struct `ReservationClientUnarchiveOptions`
- New struct `ReservationClientUnarchiveResponse`
- New struct `ReturnClient`
- New struct `ReturnClientPostOptions`
- New struct `ReturnClientPostResponse`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/reservations/armreservations` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).