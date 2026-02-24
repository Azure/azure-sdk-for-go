# Release History

## 3.1.0-beta.1 (2026-01-16)
### Features Added

- New struct `AzureFrontDoorProperties`
- New struct `ManagedOnBehalfOfConfiguration`
- New struct `MoboBrokerResource`
- New struct `TelemetryProperties`
- New field `AzureFrontDoor`, `ManagedOnBehalfOfConfiguration`, `Telemetry` in struct `ConfigurationStoreProperties`
- New field `AzureFrontDoor`, `Telemetry` in struct `ConfigurationStorePropertiesUpdateParameters`
- New field `SystemData` in struct `DeletedConfigurationStore`
- New field `SystemData` in struct `KeyValue`
- New field `SystemData` in struct `PrivateEndpointConnection`
- New field `SystemData` in struct `PrivateLinkResource`
- New field `SystemData` in struct `Snapshot`


## 3.0.0 (2025-06-10)
### Breaking Changes

- Function `*KeyValuesClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, *KeyValuesClientCreateOrUpdateOptions)` to `(context.Context, string, string, string, KeyValue, *KeyValuesClientCreateOrUpdateOptions)`
- Field `KeyValueParameters` of struct `KeyValuesClientCreateOrUpdateOptions` has been removed

### Features Added

- New field `DefaultKeyValueRevisionRetentionPeriodInSeconds` in struct `ConfigurationStoreProperties`
- New field `DefaultKeyValueRevisionRetentionPeriodInSeconds` in struct `ConfigurationStorePropertiesUpdateParameters`


## 2.2.0 (2024-11-20)
### Features Added

- New enum type `AuthenticationMode` with values `AuthenticationModeLocal`, `AuthenticationModePassThrough`
- New enum type `CompositionType` with values `CompositionTypeKey`, `CompositionTypeKeyLabel`
- New enum type `PrivateLinkDelegation` with values `PrivateLinkDelegationDisabled`, `PrivateLinkDelegationEnabled`
- New enum type `SnapshotStatus` with values `SnapshotStatusArchived`, `SnapshotStatusFailed`, `SnapshotStatusProvisioning`, `SnapshotStatusReady`
- New function `*ClientFactory.NewSnapshotsClient() *SnapshotsClient`
- New function `NewSnapshotsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SnapshotsClient, error)`
- New function `*SnapshotsClient.BeginCreate(context.Context, string, string, string, Snapshot, *SnapshotsClientBeginCreateOptions) (*runtime.Poller[SnapshotsClientCreateResponse], error)`
- New function `*SnapshotsClient.Get(context.Context, string, string, string, *SnapshotsClientGetOptions) (SnapshotsClientGetResponse, error)`
- New struct `DataPlaneProxyProperties`
- New struct `KeyValueFilter`
- New struct `Snapshot`
- New struct `SnapshotProperties`
- New field `DataPlaneProxy` in struct `ConfigurationStoreProperties`
- New field `DataPlaneProxy` in struct `ConfigurationStorePropertiesUpdateParameters`


## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-04-28)
### Breaking Changes

- Function `*KeyValuesClient.NewListByConfigurationStorePager` has been removed

### Features Added

- New enum type `ReplicaProvisioningState` with values `ReplicaProvisioningStateCanceled`, `ReplicaProvisioningStateCreating`, `ReplicaProvisioningStateDeleting`, `ReplicaProvisioningStateFailed`, `ReplicaProvisioningStateSucceeded`
- New function `*ClientFactory.NewReplicasClient() *ReplicasClient`
- New function `NewReplicasClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ReplicasClient, error)`
- New function `*ReplicasClient.BeginCreate(context.Context, string, string, string, Replica, *ReplicasClientBeginCreateOptions) (*runtime.Poller[ReplicasClientCreateResponse], error)`
- New function `*ReplicasClient.BeginDelete(context.Context, string, string, string, *ReplicasClientBeginDeleteOptions) (*runtime.Poller[ReplicasClientDeleteResponse], error)`
- New function `*ReplicasClient.Get(context.Context, string, string, string, *ReplicasClientGetOptions) (ReplicasClientGetResponse, error)`
- New function `*ReplicasClient.NewListByConfigurationStorePager(string, string, *ReplicasClientListByConfigurationStoreOptions) *runtime.Pager[ReplicasClientListByConfigurationStoreResponse]`
- New struct `Replica`
- New struct `ReplicaListResult`
- New struct `ReplicaProperties`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appconfiguration/armappconfiguration` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).