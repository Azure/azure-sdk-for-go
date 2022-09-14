# Release History

## 1.1.0 (2022-09-14)
### Features Added

- New const `ProvisioningAccepted`
- New const `ProvisioningSucceeded`
- New const `ProvisioningPending`
- New type alias `Provisioning`
- New function `PossibleProvisioningValues() []Provisioning`
- New field `ProvisioningState` in struct `AcceptOwnershipStatusResponse`
- New field `TenantID` in struct `Subscription`
- New field `Tags` in struct `Subscription`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).