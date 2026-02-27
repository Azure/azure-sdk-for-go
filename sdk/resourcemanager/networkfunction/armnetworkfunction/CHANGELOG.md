# Release History

## 3.0.0 (2026-02-05)
### Breaking Changes

- Type of `AzureTrafficCollector.SystemData` has been changed from `*TrackedResourceSystemData` to `*SystemData`
- Type of `CollectorPolicy.SystemData` has been changed from `*TrackedResourceSystemData` to `*SystemData`
- Enum `APIVersionParameter` has been removed
- Struct `ProxyResource` has been removed
- Struct `TrackedResource` has been removed
- Struct `TrackedResourceSystemData` has been removed

### Features Added

- New field `LastModifiedAt` in struct `SystemData`


## 2.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 2.0.0 (2022-12-23)
### Breaking Changes

- Type of `AzureTrafficCollectorPropertiesFormat.CollectorPolicies` has been changed from `[]*CollectorPolicy` to `[]*ResourceReference`
- Type of `CollectorPolicy.SystemData` has been changed from `*CollectorPolicySystemData` to `*TrackedResourceSystemData`
- Struct `CollectorPolicySystemData` has been removed

### Features Added

- New type alias `APIVersionParameter`
- New function `*CollectorPoliciesClient.UpdateTags(context.Context, string, string, string, TagsObject, *CollectorPoliciesClientUpdateTagsOptions) (CollectorPoliciesClientUpdateTagsResponse, error)`
- New field `Location` in struct `CollectorPolicy`
- New field `Tags` in struct `CollectorPolicy`


## 1.0.0 (2022-07-07)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/networkfunction/armnetworkfunction` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).