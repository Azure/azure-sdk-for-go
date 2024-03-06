# Release History

## 1.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-05-26)
### Features Added

- New value `ActionTypeAnomalyScoring` added to enum type `ActionType`
- New value `FrontDoorResourceStateMigrated`, `FrontDoorResourceStateMigrating` added to enum type `FrontDoorResourceState`
- New function `*PoliciesClient.NewListBySubscriptionPager(*PoliciesClientListBySubscriptionOptions) *runtime.Pager[PoliciesClientListBySubscriptionResponse]`
- New function `*PoliciesClient.BeginUpdate(context.Context, string, string, TagsObject, *PoliciesClientBeginUpdateOptions) (*runtime.Poller[PoliciesClientUpdateResponse], error)`
- New struct `DefaultErrorResponse`
- New struct `DefaultErrorResponseError`
- New field `ExtendedProperties` in struct `Properties`


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).