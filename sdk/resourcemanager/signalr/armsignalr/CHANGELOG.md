# Release History

## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0-beta.2 (2023-10-27)
### Features Added

- New struct `IPRule`
- New field `IPRules` in struct `NetworkACLs`
- New field `RegionEndpointEnabled`, `ResourceStopped` in struct `Properties`
- New field `RegionEndpointEnabled`, `ResourceStopped` in struct `ReplicaProperties`


## 1.2.0-beta.1 (2023-09-22)
### Features Added

- New function `*Client.ListReplicaSKUs(context.Context, string, string, string, *ClientListReplicaSKUsOptions) (ClientListReplicaSKUsResponse, error)`
- New function `*ClientFactory.NewReplicasClient() *ReplicasClient`
- New function `NewReplicasClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ReplicasClient, error)`
- New function `*ReplicasClient.BeginCreateOrUpdate(context.Context, string, string, string, Replica, *ReplicasClientBeginCreateOrUpdateOptions) (*runtime.Poller[ReplicasClientCreateOrUpdateResponse], error)`
- New function `*ReplicasClient.Delete(context.Context, string, string, string, *ReplicasClientDeleteOptions) (ReplicasClientDeleteResponse, error)`
- New function `*ReplicasClient.Get(context.Context, string, string, string, *ReplicasClientGetOptions) (ReplicasClientGetResponse, error)`
- New function `*ReplicasClient.NewListPager(string, string, *ReplicasClientListOptions) *runtime.Pager[ReplicasClientListResponse]`
- New function `*ReplicasClient.BeginRestart(context.Context, string, string, string, *ReplicasClientBeginRestartOptions) (*runtime.Poller[ReplicasClientRestartResponse], error)`
- New function `*ReplicasClient.BeginUpdate(context.Context, string, string, string, Replica, *ReplicasClientBeginUpdateOptions) (*runtime.Poller[ReplicasClientUpdateResponse], error)`
- New struct `Replica`
- New struct `ReplicaList`
- New struct `ReplicaProperties`
- New field `SystemData` in struct `PrivateLinkResource`
- New field `SystemData` in struct `ProxyResource`
- New field `SystemData` in struct `Resource`
- New field `SystemData` in struct `TrackedResource`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-24)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New struct `ServerlessSettings`
- New field `Serverless` in struct `Properties`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/signalr/armsignalr` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).