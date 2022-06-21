# Release History

## 2.0.0 (2022-06-21)
### Breaking Changes

- Type of `VolumeProperties.EncryptionKeySource` has been changed from `*string` to `*EncryptionKeySource`
- Field `Body` of struct `BackupsClientBeginUpdateOptions` has been removed
- Field `Tags` of struct `VolumeGroupDetails` has been removed
- Field `Tags` of struct `VolumeGroup` has been removed

### Features Added

- New field `SystemData` in struct `Resource`
- New field `SystemData` in struct `ProxyResource`
- New field `Zones` in struct `Volume`
- New field `Encrypted` in struct `VolumeProperties`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).