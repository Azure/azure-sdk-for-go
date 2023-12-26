# Release History

## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-04-03)
### Breaking Changes

- Function `NewSavingsPlanClient` parameter(s) have been changed from `(*string, azcore.TokenCredential, *arm.ClientOptions)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSavingsPlanOrderClient` parameter(s) have been changed from `(*string, azcore.TokenCredential, *arm.ClientOptions)` to `(azcore.TokenCredential, *arm.ClientOptions)`

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New field `Expand` in struct `SavingsPlanClientGetOptions`
- New field `Expand` in struct `SavingsPlanOrderClientGetOptions`


## 1.0.0 (2022-12-23)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billingbenefits/armbillingbenefits` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).