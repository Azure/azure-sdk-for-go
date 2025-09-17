# Release History

## 2.0.0 (2025-09-17)
### Breaking Changes

- Function `*FirewallRulesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, FirewallRule, *FirewallRulesClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, string, string, FirewallRule, *FirewallRulesClientBeginCreateOrUpdateOptions)`
- Function `*FirewallRulesClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, string, *FirewallRulesClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, *FirewallRulesClientBeginDeleteOptions)`
- Function `*FirewallRulesClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *FirewallRulesClientGetOptions)` to `(context.Context, string, string, string, string, *FirewallRulesClientGetOptions)`
- Function `*FirewallRulesClient.NewListByMongoClusterPager` parameter(s) have been changed from `(string, string, *FirewallRulesClientListByMongoClusterOptions)` to `(string, string, string, *FirewallRulesClientListByMongoClusterOptions)`
- Function `*MongoClustersClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, MongoCluster, *MongoClustersClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, string, MongoCluster, *MongoClustersClientBeginCreateOrUpdateOptions)`
- Function `*MongoClustersClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *MongoClustersClientBeginDeleteOptions)` to `(context.Context, string, string, string, *MongoClustersClientBeginDeleteOptions)`
- Function `*MongoClustersClient.BeginPromote` parameter(s) have been changed from `(context.Context, string, string, PromoteReplicaRequest, *MongoClustersClientBeginPromoteOptions)` to `(context.Context, string, string, string, PromoteReplicaRequest, *MongoClustersClientBeginPromoteOptions)`
- Function `*MongoClustersClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, Update, *MongoClustersClientBeginUpdateOptions)` to `(context.Context, string, string, string, Update, *MongoClustersClientBeginUpdateOptions)`
- Function `*MongoClustersClient.CheckNameAvailability` parameter(s) have been changed from `(context.Context, string, CheckNameAvailabilityRequest, *MongoClustersClientCheckNameAvailabilityOptions)` to `(context.Context, string, string, CheckNameAvailabilityRequest, *MongoClustersClientCheckNameAvailabilityOptions)`
- Function `*MongoClustersClient.Get` parameter(s) have been changed from `(context.Context, string, string, *MongoClustersClientGetOptions)` to `(context.Context, string, string, string, *MongoClustersClientGetOptions)`
- Function `*MongoClustersClient.ListConnectionStrings` parameter(s) have been changed from `(context.Context, string, string, *MongoClustersClientListConnectionStringsOptions)` to `(context.Context, string, string, string, *MongoClustersClientListConnectionStringsOptions)`
- Function `*MongoClustersClient.NewListByResourceGroupPager` parameter(s) have been changed from `(string, *MongoClustersClientListByResourceGroupOptions)` to `(string, string, *MongoClustersClientListByResourceGroupOptions)`
- Function `*MongoClustersClient.NewListPager` parameter(s) have been changed from `(*MongoClustersClientListOptions)` to `(string, *MongoClustersClientListOptions)`
- Function `*OperationsClient.NewListPager` parameter(s) have been changed from `(*OperationsClientListOptions)` to `(string, *OperationsClientListOptions)`
- Function `*PrivateEndpointConnectionsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, string, PrivateEndpointConnectionResource, *PrivateEndpointConnectionsClientBeginCreateOptions)` to `(context.Context, string, string, string, string, PrivateEndpointConnectionResource, *PrivateEndpointConnectionsClientBeginCreateOptions)`
- Function `*PrivateEndpointConnectionsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions)`
- Function `*PrivateEndpointConnectionsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions)` to `(context.Context, string, string, string, string, *PrivateEndpointConnectionsClientGetOptions)`
- Function `*PrivateEndpointConnectionsClient.NewListByMongoClusterPager` parameter(s) have been changed from `(string, string, *PrivateEndpointConnectionsClientListByMongoClusterOptions)` to `(string, string, string, *PrivateEndpointConnectionsClientListByMongoClusterOptions)`
- Function `*PrivateLinksClient.NewListByMongoClusterPager` parameter(s) have been changed from `(string, string, *PrivateLinksClientListByMongoClusterOptions)` to `(string, string, string, *PrivateLinksClientListByMongoClusterOptions)`
- Function `*ReplicasClient.NewListByParentPager` parameter(s) have been changed from `(string, string, *ReplicasClientListByParentOptions)` to `(string, string, string, *ReplicasClientListByParentOptions)`

### Features Added

- New enum type `AuthenticationMode` with values `AuthenticationModeMicrosoftEntraID`, `AuthenticationModeNativeAuth`
- New enum type `DataAPIMode` with values `DataAPIModeDisabled`, `DataAPIModeEnabled`
- New enum type `EntraPrincipalType` with values `EntraPrincipalTypeServicePrincipal`, `EntraPrincipalTypeUser`
- New enum type `IdentityProviderType` with values `IdentityProviderTypeMicrosoftEntraID`
- New enum type `KeyEncryptionKeyIdentityType` with values `KeyEncryptionKeyIdentityTypeUserAssignedIdentity`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `StorageType` with values `StorageTypePremiumSSD`, `StorageTypePremiumSSDv2`
- New enum type `UserRole` with values `UserRoleRoot`
- New function `*ClientFactory.NewUsersClient() *UsersClient`
- New function `*EntraIdentityProvider.GetIdentityProvider() *IdentityProvider`
- New function `*IdentityProvider.GetIdentityProvider() *IdentityProvider`
- New function `NewUsersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*UsersClient, error)`
- New function `*UsersClient.BeginCreateOrUpdate(context.Context, string, string, string, string, User, *UsersClientBeginCreateOrUpdateOptions) (*runtime.Poller[UsersClientCreateOrUpdateResponse], error)`
- New function `*UsersClient.BeginDelete(context.Context, string, string, string, string, *UsersClientBeginDeleteOptions) (*runtime.Poller[UsersClientDeleteResponse], error)`
- New function `*UsersClient.Get(context.Context, string, string, string, string, *UsersClientGetOptions) (UsersClientGetResponse, error)`
- New function `*UsersClient.NewListByMongoClusterPager(string, string, string, *UsersClientListByMongoClusterOptions) *runtime.Pager[UsersClientListByMongoClusterResponse]`
- New struct `AuthConfigProperties`
- New struct `CustomerManagedKeyEncryptionProperties`
- New struct `DataAPIProperties`
- New struct `DatabaseRole`
- New struct `EncryptionProperties`
- New struct `EntraIdentityProvider`
- New struct `EntraIdentityProviderProperties`
- New struct `KeyEncryptionKeyIdentity`
- New struct `ManagedServiceIdentity`
- New struct `User`
- New struct `UserAssignedIdentity`
- New struct `UserListResult`
- New struct `UserProperties`
- New field `Identity` in struct `MongoCluster`
- New field `AuthConfig`, `DataAPI`, `Encryption` in struct `Properties`
- New field `Iops`, `Throughput`, `Type` in struct `StorageProperties`
- New field `Identity` in struct `Update`
- New field `AuthConfig`, `DataAPI` in struct `UpdateProperties`


## 1.1.0-beta.1 (2025-07-23)
### Features Added

- New enum type `AuthenticationMode` with values `AuthenticationModeMicrosoftEntraID`, `AuthenticationModeNativeAuth`
- New enum type `DataAPIMode` with values `DataAPIModeDisabled`, `DataAPIModeEnabled`
- New enum type `EntraPrincipalType` with values `EntraPrincipalTypeServicePrincipal`, `EntraPrincipalTypeUser`
- New enum type `IdentityProviderType` with values `IdentityProviderTypeMicrosoftEntraID`
- New enum type `KeyEncryptionKeyIdentityType` with values `KeyEncryptionKeyIdentityTypeUserAssignedIdentity`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `StorageType` with values `StorageTypePremiumSSD`, `StorageTypePremiumSSDv2`
- New enum type `UserRole` with values `UserRoleRoot`
- New function `*ClientFactory.NewUsersClient() *UsersClient`
- New function `*EntraIdentityProvider.GetIdentityProvider() *IdentityProvider`
- New function `*IdentityProvider.GetIdentityProvider() *IdentityProvider`
- New function `NewUsersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*UsersClient, error)`
- New function `*UsersClient.BeginCreateOrUpdate(context.Context, string, string, string, User, *UsersClientBeginCreateOrUpdateOptions) (*runtime.Poller[UsersClientCreateOrUpdateResponse], error)`
- New function `*UsersClient.BeginDelete(context.Context, string, string, string, *UsersClientBeginDeleteOptions) (*runtime.Poller[UsersClientDeleteResponse], error)`
- New function `*UsersClient.Get(context.Context, string, string, string, *UsersClientGetOptions) (UsersClientGetResponse, error)`
- New function `*UsersClient.NewListByMongoClusterPager(string, string, *UsersClientListByMongoClusterOptions) *runtime.Pager[UsersClientListByMongoClusterResponse]`
- New struct `AuthConfigProperties`
- New struct `CustomerManagedKeyEncryptionProperties`
- New struct `DataAPIProperties`
- New struct `DatabaseRole`
- New struct `EncryptionProperties`
- New struct `EntraIdentityProvider`
- New struct `EntraIdentityProviderProperties`
- New struct `KeyEncryptionKeyIdentity`
- New struct `ManagedServiceIdentity`
- New struct `User`
- New struct `UserAssignedIdentity`
- New struct `UserListResult`
- New struct `UserProperties`
- New field `Identity` in struct `MongoCluster`
- New field `AuthConfig`, `DataAPI`, `Encryption` in struct `Properties`
- New field `Iops`, `Throughput`, `Type` in struct `StorageProperties`
- New field `Identity` in struct `Update`
- New field `AuthConfig`, `DataAPI` in struct `UpdateProperties`


## 1.0.1 (2024-10-14)
### Other Changes
- Add examples

## 1.0.0 (2024-09-27)
### Breaking Changes

- Enum `NodeKind` has been removed
- Struct `NodeGroupSpec` has been removed
- Field `AdministratorLogin`, `AdministratorLoginPassword`, `EarliestRestoreTime`, `NodeGroupSpecs` of struct `Properties` has been removed
- Field `AdministratorLogin`, `AdministratorLoginPassword`, `NodeGroupSpecs` of struct `UpdateProperties` has been removed

### Features Added

- New enum type `HighAvailabilityMode` with values `HighAvailabilityModeDisabled`, `HighAvailabilityModeSameZone`, `HighAvailabilityModeZoneRedundantPreferred`
- New struct `AdministratorProperties`
- New struct `BackupProperties`
- New struct `ComputeProperties`
- New struct `HighAvailabilityProperties`
- New struct `ShardingProperties`
- New struct `StorageProperties`
- New field `Name` in struct `ConnectionString`
- New field `Administrator`, `Backup`, `Compute`, `HighAvailability`, `Sharding`, `Storage` in struct `Properties`
- New field `Administrator`, `Backup`, `Compute`, `HighAvailability`, `Sharding`, `Storage` in struct `UpdateProperties`


## 0.2.0 (2024-09-26)
### Breaking Changes

- Type of `Properties.ClusterStatus` has been changed from `*MongoClusterStatus` to `*Status`
- Enum `MongoClusterStatus` has been removed

### Features Added

- New value `CreateModeGeoReplica`, `CreateModeReplica` added to enum type `CreateMode`
- New enum type `PreviewFeature` with values `PreviewFeatureGeoReplicas`
- New enum type `PromoteMode` with values `PromoteModeSwitchover`
- New enum type `PromoteOption` with values `PromoteOptionForced`
- New enum type `ReplicationRole` with values `ReplicationRoleAsyncReplica`, `ReplicationRoleGeoAsyncReplica`, `ReplicationRolePrimary`
- New enum type `ReplicationState` with values `ReplicationStateActive`, `ReplicationStateBroken`, `ReplicationStateCatchup`, `ReplicationStateProvisioning`, `ReplicationStateReconfiguring`, `ReplicationStateUpdating`
- New enum type `Status` with values `StatusDropping`, `StatusProvisioning`, `StatusReady`, `StatusStarting`, `StatusStopped`, `StatusStopping`, `StatusUpdating`
- New function `*ClientFactory.NewReplicasClient() *ReplicasClient`
- New function `*MongoClustersClient.BeginPromote(context.Context, string, string, PromoteReplicaRequest, *MongoClustersClientBeginPromoteOptions) (*runtime.Poller[MongoClustersClientPromoteResponse], error)`
- New function `NewReplicasClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ReplicasClient, error)`
- New function `*ReplicasClient.NewListByParentPager(string, string, *ReplicasClientListByParentOptions) *runtime.Pager[ReplicasClientListByParentResponse]`
- New struct `PromoteReplicaRequest`
- New struct `Replica`
- New struct `ReplicaListResult`
- New struct `ReplicaParameters`
- New struct `ReplicationProperties`
- New field `InfrastructureVersion`, `PreviewFeatures`, `Replica`, `ReplicaParameters` in struct `Properties`
- New field `PreviewFeatures` in struct `UpdateProperties`


## 0.1.0 (2024-07-05)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mongocluster/armmongocluster` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
