# Release History

## 1.1.0 (2022-09-19)
### Features Added

- New const `CloudTieringLowDiskModeStateDisabled`
- New const `CloudTieringLowDiskModeStateEnabled`
- New type alias `CloudTieringLowDiskModeState`
- New function `PossibleCloudTieringLowDiskModeStateValues() []CloudTieringLowDiskModeState`
- New function `*CloudEndpointsClient.AfsShareMetadataCertificatePublicKeys(context.Context, string, string, string, string, *CloudEndpointsClientAfsShareMetadataCertificatePublicKeysOptions) (CloudEndpointsClientAfsShareMetadataCertificatePublicKeysResponse, error)`
- New struct `CloudEndpointAfsShareMetadataCertificatePublicKeys`
- New struct `CloudEndpointsClientAfsShareMetadataCertificatePublicKeysOptions`
- New struct `CloudEndpointsClientAfsShareMetadataCertificatePublicKeysResponse`
- New struct `CloudTieringLowDiskMode`
- New field `LowDiskMode` in struct `ServerEndpointCloudTieringStatus`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagesync/armstoragesync` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).