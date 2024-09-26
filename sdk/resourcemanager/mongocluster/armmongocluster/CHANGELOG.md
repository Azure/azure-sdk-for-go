# Release History

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