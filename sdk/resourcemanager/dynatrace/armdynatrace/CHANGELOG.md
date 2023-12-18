# Release History

## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-08-25)
### Breaking Changes

- Function `*MonitorsClient.GetAccountCredentials` has been removed
- Function `*TagRulesClient.Update` has been removed
- Struct `AccountInfoSecure` has been removed
- Struct `TagRuleUpdate` has been removed
- Field `DynatraceEnvironmentProperties`, `MarketplaceSubscriptionStatus`, `MonitoringStatus`, `PlanData`, `UserInfo` of struct `MonitorResourceUpdate` has been removed

### Features Added

- New function `*MonitorsClient.GetMarketplaceSaaSResourceDetails(context.Context, MarketplaceSaaSResourceDetailsRequest, *MonitorsClientGetMarketplaceSaaSResourceDetailsOptions) (MonitorsClientGetMarketplaceSaaSResourceDetailsResponse, error)`
- New function `*MonitorsClient.GetMetricStatus(context.Context, string, string, *MonitorsClientGetMetricStatusOptions) (MonitorsClientGetMetricStatusResponse, error)`
- New struct `MarketplaceSaaSResourceDetailsRequest`
- New struct `MarketplaceSaaSResourceDetailsResponse`
- New struct `MetricsStatusResponse`
- New field `SendingMetrics` in struct `MetricRules`


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-09-20)
### Other Changes

- Release stable version.

## 0.1.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dynatrace/armdynatrace` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).