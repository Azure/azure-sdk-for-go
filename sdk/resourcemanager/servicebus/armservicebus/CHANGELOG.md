# Release History

## 2.0.0-beta.2 (2023-02-24)
### Breaking Changes

- Type of `Encryption.KeySource` has been changed from `*KeySource` to `*string`
- Type alias `KeySource` has been removed

### Features Added

- New field `PremiumMessagingPartitions` in struct `SBNamespaceProperties`


## 2.0.0-beta.1 (2022-05-24)
### Breaking Changes

- Type of `Encryption.KeySource` has been changed from `*string` to `*KeySource`

### Features Added

- New const `TLSVersionOne2`
- New const `KeySourceMicrosoftKeyVault`
- New const `TLSVersionOne1`
- New const `PublicNetworkAccessDisabled`
- New const `TLSVersionOne0`
- New const `PublicNetworkAccessSecuredByPerimeter`
- New const `PublicNetworkAccessEnabled`
- New function `PossibleKeySourceValues() []KeySource`
- New function `PossiblePublicNetworkAccessValues() []PublicNetworkAccess`
- New function `PossibleTLSVersionValues() []TLSVersion`
- New field `MinimumTLSVersion` in struct `SBNamespaceProperties`
- New field `PublicNetworkAccess` in struct `SBNamespaceProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).