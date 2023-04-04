# Release History

## 2.0.0 (2023-04-07)
### Breaking Changes

- Struct `CloudError` has been removed
- Struct `CloudErrorBody` has been removed

### Features Added

- New function `NewClientFactory(string, azcore.TokenCredential, *arm.ClientOptions) (*ClientFactory, error)`
- New function `*ClientFactory.NewEndpointsClient() *EndpointsClient`
- New function `*ClientFactory.NewGeographicHierarchiesClient() *GeographicHierarchiesClient`
- New function `*ClientFactory.NewHeatMapClient() *HeatMapClient`
- New function `*ClientFactory.NewProfilesClient() *ProfilesClient`
- New function `*ClientFactory.NewUserMetricsKeysClient() *UserMetricsKeysClient`
- New struct `ClientFactory`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/trafficmanager/armtrafficmanager` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).