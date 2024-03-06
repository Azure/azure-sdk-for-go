# Release History

## 1.3.0-beta.1 (2023-11-30)
### Features Added

- New enum type `ServiceKind` with values `ServiceKindSocketIO`, `ServiceKindWebPubSub`
- New function `*Client.ListReplicaSKUs(context.Context, string, string, string, *ClientListReplicaSKUsOptions) (ClientListReplicaSKUsResponse, error)`
- New function `*ClientFactory.NewReplicasClient() *ReplicasClient`
- New function `NewReplicasClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ReplicasClient, error)`
- New function `*ReplicasClient.BeginCreateOrUpdate(context.Context, string, string, string, Replica, *ReplicasClientBeginCreateOrUpdateOptions) (*runtime.Poller[ReplicasClientCreateOrUpdateResponse], error)`
- New function `*ReplicasClient.Delete(context.Context, string, string, string, *ReplicasClientDeleteOptions) (ReplicasClientDeleteResponse, error)`
- New function `*ReplicasClient.Get(context.Context, string, string, string, *ReplicasClientGetOptions) (ReplicasClientGetResponse, error)`
- New function `*ReplicasClient.NewListPager(string, string, *ReplicasClientListOptions) *runtime.Pager[ReplicasClientListResponse]`
- New function `*ReplicasClient.BeginRestart(context.Context, string, string, string, *ReplicasClientBeginRestartOptions) (*runtime.Poller[ReplicasClientRestartResponse], error)`
- New function `*ReplicasClient.BeginUpdate(context.Context, string, string, string, Replica, *ReplicasClientBeginUpdateOptions) (*runtime.Poller[ReplicasClientUpdateResponse], error)`
- New struct `IPRule`
- New struct `Replica`
- New struct `ReplicaList`
- New struct `ReplicaProperties`
- New field `IPRules` in struct `NetworkACLs`
- New field `SystemData` in struct `PrivateLinkResource`
- New field `RegionEndpointEnabled`, `ResourceStopped` in struct `Properties`
- New field `Kind` in struct `ResourceInfo`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0-beta.2 (2023-10-27)
### Features Added

- New struct `IPRule`
- New field `IPRules` in struct `NetworkACLs`
- New field `RegionEndpointEnabled`, `ResourceStopped` in struct `Properties`
- New field `RegionEndpointEnabled`, `ResourceStopped` in struct `ReplicaProperties`


## 1.2.0-beta.1 (2023-07-28)
### Features Added

- New enum type `ServiceKind` with values `ServiceKindSocketIO`, `ServiceKindWebPubSub`
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
- New field `Kind` in struct `ResourceInfo`
- New field `SystemData` in struct `TrackedResource`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-24)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New enum type `EventListenerEndpointDiscriminator` with values `EventListenerEndpointDiscriminatorEventHub`
- New enum type `EventListenerFilterDiscriminator` with values `EventListenerFilterDiscriminatorEventName`
- New function `NewCustomCertificatesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CustomCertificatesClient, error)`
- New function `*CustomCertificatesClient.BeginCreateOrUpdate(context.Context, string, string, string, CustomCertificate, *CustomCertificatesClientBeginCreateOrUpdateOptions) (*runtime.Poller[CustomCertificatesClientCreateOrUpdateResponse], error)`
- New function `*CustomCertificatesClient.Delete(context.Context, string, string, string, *CustomCertificatesClientDeleteOptions) (CustomCertificatesClientDeleteResponse, error)`
- New function `*CustomCertificatesClient.Get(context.Context, string, string, string, *CustomCertificatesClientGetOptions) (CustomCertificatesClientGetResponse, error)`
- New function `*CustomCertificatesClient.NewListPager(string, string, *CustomCertificatesClientListOptions) *runtime.Pager[CustomCertificatesClientListResponse]`
- New function `NewCustomDomainsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CustomDomainsClient, error)`
- New function `*CustomDomainsClient.BeginCreateOrUpdate(context.Context, string, string, string, CustomDomain, *CustomDomainsClientBeginCreateOrUpdateOptions) (*runtime.Poller[CustomDomainsClientCreateOrUpdateResponse], error)`
- New function `*CustomDomainsClient.BeginDelete(context.Context, string, string, string, *CustomDomainsClientBeginDeleteOptions) (*runtime.Poller[CustomDomainsClientDeleteResponse], error)`
- New function `*CustomDomainsClient.Get(context.Context, string, string, string, *CustomDomainsClientGetOptions) (CustomDomainsClientGetResponse, error)`
- New function `*CustomDomainsClient.NewListPager(string, string, *CustomDomainsClientListOptions) *runtime.Pager[CustomDomainsClientListResponse]`
- New function `*EventHubEndpoint.GetEventListenerEndpoint() *EventListenerEndpoint`
- New function `*EventListenerEndpoint.GetEventListenerEndpoint() *EventListenerEndpoint`
- New function `*EventListenerFilter.GetEventListenerFilter() *EventListenerFilter`
- New function `*EventNameFilter.GetEventListenerFilter() *EventListenerFilter`
- New struct `CustomCertificate`
- New struct `CustomCertificateList`
- New struct `CustomCertificateProperties`
- New struct `CustomDomain`
- New struct `CustomDomainList`
- New struct `CustomDomainProperties`
- New struct `EventHubEndpoint`
- New struct `EventListener`
- New struct `EventNameFilter`
- New struct `ResourceReference`
- New field `EventListeners` in struct `HubProperties`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/webpubsub/armwebpubsub` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).