# Release History

## 0.2.0 (2026-09-22)
### Breaking Changes

- Type of `AppLinkMemberUpdateProperties.ConnectivityProfile` has been changed from `*ConnectivityProfile` to `*ConnectivityProfileUpdate`
- Type of `AppLinkMemberUpdateProperties.UpgradeProfile` has been changed from `*UpgradeProfile` to `*UpgradeProfileUpdate`
- Field `ObservabilityProfile` of struct `AppLinkMemberUpdateProperties` has been removed

### Features Added

- New struct `ConnectivityProfileUpdate`
- New struct `EastWestGatewayProfileUpdate`
- New struct `FullyManagedUpgradeProfileUpdate`
- New struct `ManagedServiceIdentityUpdate`
- New struct `SelfManagedUpgradeProfileUpdate`
- New struct `UpgradeProfileUpdate`
- New field `Identity` in struct `AppLinkUpdate`
- New field `Network` in struct `ConnectivityProfile`


## 0.1.0 (2026-03-25)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appnetwork/armappnetwork` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).