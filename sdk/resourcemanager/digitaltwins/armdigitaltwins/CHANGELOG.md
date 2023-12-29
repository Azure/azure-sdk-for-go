# Release History

## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-24)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New value `DigitalTwinsIdentityTypeSystemAssignedUserAssigned`, `DigitalTwinsIdentityTypeUserAssigned` added to enum type `DigitalTwinsIdentityType`
- New enum type `CleanupConnectionArtifacts` with values `CleanupConnectionArtifactsFalse`, `CleanupConnectionArtifactsTrue`
- New enum type `IdentityType` with values `IdentityTypeSystemAssigned`, `IdentityTypeUserAssigned`
- New enum type `RecordPropertyAndItemRemovals` with values `RecordPropertyAndItemRemovalsFalse`, `RecordPropertyAndItemRemovalsTrue`
- New struct `ManagedIdentityReference`
- New struct `UserAssignedIdentity`
- New field `AdxRelationshipLifecycleEventsTableName` in struct `AzureDataExplorerConnectionProperties`
- New field `AdxTwinLifecycleEventsTableName` in struct `AzureDataExplorerConnectionProperties`
- New field `Identity` in struct `AzureDataExplorerConnectionProperties`
- New field `RecordPropertyAndItemRemovals` in struct `AzureDataExplorerConnectionProperties`
- New field `Identity` in struct `EndpointResourceProperties`
- New field `Identity` in struct `EventGrid`
- New field `Identity` in struct `EventHub`
- New field `UserAssignedIdentities` in struct `Identity`
- New field `Identity` in struct `ServiceBus`
- New field `Identity` in struct `TimeSeriesDatabaseConnectionProperties`
- New field `CleanupConnectionArtifacts` in struct `TimeSeriesDatabaseConnectionsClientBeginDeleteOptions`


## 1.0.0 (2022-06-14)
### Features Added

- New const `TimeSeriesDatabaseConnectionStateUpdating`
- New const `EndpointProvisioningStateUpdating`


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/digitaltwins/armdigitaltwins` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).