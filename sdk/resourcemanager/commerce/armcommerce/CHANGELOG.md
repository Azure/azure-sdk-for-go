# Release History

## 0.4.0 (2026-03-11)
### Breaking Changes

- Struct `InfoField` has been removed
- Field `UsageAggregationListResult` of struct `UsageAggregatesClientListResponse` has been removed

### Features Added

- Type of `UsageSample.InfoFields` has been changed from `*InfoField` to `any`
- New struct `ErrorObjectResponse`
- New field `Value` in struct `UsageAggregatesClientListResponse`


## 0.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.2.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.1.0 (2022-06-10)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/commerce/armcommerce` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).