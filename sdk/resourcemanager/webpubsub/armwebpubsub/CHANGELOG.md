# Release History

## 1.1.1 (2023-04-14)
### Other Changes


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