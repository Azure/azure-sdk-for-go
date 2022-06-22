# Release History

## 2.0.0 (2022-06-22)
### Breaking Changes

- Function `*Client.Update` has been removed
- Function `*Client.CreateOrUpdate` has been removed
- Struct `ClientCreateOrUpdateOptions` has been removed
- Struct `ClientUpdateOptions` has been removed

### Features Added

- New field `PrivateEndpointConnections` in struct `MediaServiceProperties`
- New field `ProvisioningState` in struct `MediaServiceProperties`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mediaservices/armmediaservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).