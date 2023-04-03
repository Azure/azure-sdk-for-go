# Release History

## 0.6.0 (2023-04-03)
### Breaking Changes

- Function `NewDeviceGroupsClient` parameter(s) have been changed from `(string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*DeviceGroupsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, DeviceGroupModel, *DeviceGroupsClientCreateOrUpdateOptions)` to `(context.Context, string, string, DeviceGroupModel, *DeviceGroupsClientCreateOrUpdateOptions)`
- Function `*DeviceGroupsClient.Delete` parameter(s) have been changed from `(context.Context, string, *DeviceGroupsClientDeleteOptions)` to `(context.Context, string, string, *DeviceGroupsClientDeleteOptions)`
- Function `*DeviceGroupsClient.Get` parameter(s) have been changed from `(context.Context, string, *DeviceGroupsClientGetOptions)` to `(context.Context, string, string, *DeviceGroupsClientGetOptions)`
- Function `*DeviceGroupsClient.NewListPager` parameter(s) have been changed from `(*DeviceGroupsClientListOptions)` to `(string, *DeviceGroupsClientListOptions)`
- Function `NewDevicesClient` parameter(s) have been changed from `(string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*DevicesClient.Get` parameter(s) have been changed from `(context.Context, string, string, *DevicesClientGetOptions)` to `(context.Context, string, string, string, *DevicesClientGetOptions)`
- Function `*DevicesClient.NewListPager` parameter(s) have been changed from `(string, *DevicesClientListOptions)` to `(string, string, *DevicesClientListOptions)`
- Function `NewLocationsClient` parameter(s) have been changed from `(string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*LocationsClient.Get` parameter(s) have been changed from `(context.Context, *LocationsClientGetOptions)` to `(context.Context, string, *LocationsClientGetOptions)`

### Features Added

- New function `NewClientFactory(string, azcore.TokenCredential, *arm.ClientOptions) (*ClientFactory, error)`
- New function `*ClientFactory.NewDefenderSettingsClient() *DefenderSettingsClient`
- New function `*ClientFactory.NewDeviceGroupsClient() *DeviceGroupsClient`
- New function `*ClientFactory.NewDevicesClient() *DevicesClient`
- New function `*ClientFactory.NewLocationsClient() *LocationsClient`
- New function `*ClientFactory.NewOnPremiseSensorsClient() *OnPremiseSensorsClient`
- New function `*ClientFactory.NewOperationsClient() *OperationsClient`
- New function `*ClientFactory.NewSensorsClient() *SensorsClient`
- New function `*ClientFactory.NewSitesClient() *SitesClient`
- New struct `ClientFactory`


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotsecurity/armiotsecurity` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).