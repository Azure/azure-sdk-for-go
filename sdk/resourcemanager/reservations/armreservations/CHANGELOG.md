# Release History

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