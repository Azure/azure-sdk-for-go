# Release History

## 1.0.2 (2025-09-09)

### Other Changes

- This module is now DEPRECATED. The latest supported version of this module is at [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/playwright/armplaywright](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/playwright/armplaywright)

## 1.0.1 (2025-09-08)

### Other Changes

- This module is now DEPRECATED. The latest supported version of this module is at [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/playwright/armplaywright](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/playwright/armplaywright)

## 1.0.0 (2024-12-26)
### Breaking Changes

- Field `AllocatedValue`, `CreatedAt`, `ExpiryAt`, `PercentageUsed`, `UsedValue` of struct `FreeTrialProperties` has been removed

### Features Added

- New value `FreeTrialStateNotEligible`, `FreeTrialStateNotRegistered` added to enum type `FreeTrialState`
- New value `ProvisioningStateCreating` added to enum type `ProvisioningState`
- New value `QuotaNamesReporting` added to enum type `QuotaNames`
- New enum type `CheckNameAvailabilityReason` with values `CheckNameAvailabilityReasonAlreadyExists`, `CheckNameAvailabilityReasonInvalid`
- New enum type `OfferingType` with values `OfferingTypeGeneralAvailability`, `OfferingTypeNotApplicable`, `OfferingTypePrivatePreview`, `OfferingTypePublicPreview`
- New function `NewAccountQuotasClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccountQuotasClient, error)`
- New function `*AccountQuotasClient.Get(context.Context, string, string, QuotaNames, *AccountQuotasClientGetOptions) (AccountQuotasClientGetResponse, error)`
- New function `*AccountQuotasClient.NewListByAccountPager(string, string, *AccountQuotasClientListByAccountOptions) *runtime.Pager[AccountQuotasClientListByAccountResponse]`
- New function `*AccountsClient.CheckNameAvailability(context.Context, CheckNameAvailabilityRequest, *AccountsClientCheckNameAvailabilityOptions) (AccountsClientCheckNameAvailabilityResponse, error)`
- New function `*ClientFactory.NewAccountQuotasClient() *AccountQuotasClient`
- New struct `AccountFreeTrialProperties`
- New struct `AccountQuota`
- New struct `AccountQuotaListResult`
- New struct `AccountQuotaProperties`
- New struct `CheckNameAvailabilityRequest`
- New struct `CheckNameAvailabilityResponse`
- New field `LocalAuth` in struct `AccountProperties`
- New field `LocalAuth` in struct `AccountUpdateProperties`
- New field `OfferingType` in struct `QuotaProperties`


## 0.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.1.0 (2023-09-27)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/playwrighttesting/armplaywrighttesting` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).