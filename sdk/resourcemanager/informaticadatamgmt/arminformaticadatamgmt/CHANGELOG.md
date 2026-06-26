# Release History

## 2.0.0 (2026-06-26)
### Breaking Changes

- Enum `ActionType` has been removed
- Enum `Origin` has been removed
- Function `*OperationsClient.NewListPager` has been removed
- Struct `Operation` has been removed
- Struct `OperationDisplay` has been removed
- Struct `OperationListResult` has been removed

### Features Added

- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `MarketplaceSubscriptionStatus` with values `MarketplaceSubscriptionStatusPendingFulfillmentStart`, `MarketplaceSubscriptionStatusSubscribe`, `MarketplaceSubscriptionStatusSuspend`, `MarketplaceSubscriptionStatusUnsubscribe`
- New struct `ManagedServiceIdentity`
- New struct `ServerlessRuntimeDataDisk`
- New struct `UserAssignedIdentity`
- New field `ServerlessRuntimeDataDisks` in struct `InfaServerlessFetchConfigProperties`
- New field `Identity` in struct `InformaticaOrganizationResource`
- New field `Identity` in struct `InformaticaOrganizationResourceUpdate`
- New field `ServerlessRuntimeDataDisks` in struct `InformaticaServerlessRuntimeProperties`
- New field `Identity` in struct `InformaticaServerlessRuntimeResource`
- New field `Identity` in struct `InformaticaServerlessRuntimeResourceUpdate`
- New field `MarketplaceSubscriptionStatus` in struct `MarketplaceDetails`
- New field `MarketplaceSubscriptionStatus` in struct `MarketplaceDetailsUpdate`
- New field `ServerlessRuntimeDataDisks` in struct `ServerlessRuntimePropertiesCustomUpdate`


## 1.0.0 (2024-07-15)
### Other Changes

- Release stable version.


## 0.1.0 (2024-05-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/informaticadatamgmt/arminformaticadatamgmt` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
