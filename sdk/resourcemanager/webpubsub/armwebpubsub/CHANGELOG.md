# Release History

## 1.1.0-beta.1 (2022-11-03)
### Features Added

- New const `EventListenerFilterDiscriminatorEventName`
- New const `EventListenerEndpointDiscriminatorEventHub`
- New type alias `EventListenerEndpointDiscriminator`
- New type alias `EventListenerFilterDiscriminator`
- New function `PossibleEventListenerFilterDiscriminatorValues() []EventListenerFilterDiscriminator`
- New function `*EventHubEndpoint.GetEventListenerEndpoint() *EventListenerEndpoint`
- New function `*CustomCertificatesClient.BeginCreateOrUpdate(context.Context, string, string, string, CustomCertificate, *CustomCertificatesClientBeginCreateOrUpdateOptions) (*runtime.Poller[CustomCertificatesClientCreateOrUpdateResponse], error)`
- New function `*CustomCertificatesClient.NewListPager(string, string, *CustomCertificatesClientListOptions) *runtime.Pager[CustomCertificatesClientListResponse]`
- New function `*CustomCertificatesClient.Get(context.Context, string, string, string, *CustomCertificatesClientGetOptions) (CustomCertificatesClientGetResponse, error)`
- New function `*CustomDomainsClient.BeginDelete(context.Context, string, string, string, *CustomDomainsClientBeginDeleteOptions) (*runtime.Poller[CustomDomainsClientDeleteResponse], error)`
- New function `*CustomDomainsClient.Get(context.Context, string, string, string, *CustomDomainsClientGetOptions) (CustomDomainsClientGetResponse, error)`
- New function `NewCustomCertificatesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CustomCertificatesClient, error)`
- New function `PossibleEventListenerEndpointDiscriminatorValues() []EventListenerEndpointDiscriminator`
- New function `*CustomCertificatesClient.Delete(context.Context, string, string, string, *CustomCertificatesClientDeleteOptions) (CustomCertificatesClientDeleteResponse, error)`
- New function `*EventListenerEndpoint.GetEventListenerEndpoint() *EventListenerEndpoint`
- New function `*CustomDomainsClient.NewListPager(string, string, *CustomDomainsClientListOptions) *runtime.Pager[CustomDomainsClientListResponse]`
- New function `*CustomDomainsClient.BeginCreateOrUpdate(context.Context, string, string, string, CustomDomain, *CustomDomainsClientBeginCreateOrUpdateOptions) (*runtime.Poller[CustomDomainsClientCreateOrUpdateResponse], error)`
- New function `*EventListenerFilter.GetEventListenerFilter() *EventListenerFilter`
- New function `NewCustomDomainsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CustomDomainsClient, error)`
- New function `*EventNameFilter.GetEventListenerFilter() *EventListenerFilter`
- New struct `CustomCertificate`
- New struct `CustomCertificateList`
- New struct `CustomCertificateProperties`
- New struct `CustomCertificatesClient`
- New struct `CustomCertificatesClientBeginCreateOrUpdateOptions`
- New struct `CustomCertificatesClientCreateOrUpdateResponse`
- New struct `CustomCertificatesClientDeleteOptions`
- New struct `CustomCertificatesClientDeleteResponse`
- New struct `CustomCertificatesClientGetOptions`
- New struct `CustomCertificatesClientGetResponse`
- New struct `CustomCertificatesClientListOptions`
- New struct `CustomCertificatesClientListResponse`
- New struct `CustomDomain`
- New struct `CustomDomainList`
- New struct `CustomDomainProperties`
- New struct `CustomDomainsClient`
- New struct `CustomDomainsClientBeginCreateOrUpdateOptions`
- New struct `CustomDomainsClientBeginDeleteOptions`
- New struct `CustomDomainsClientCreateOrUpdateResponse`
- New struct `CustomDomainsClientDeleteResponse`
- New struct `CustomDomainsClientGetOptions`
- New struct `CustomDomainsClientGetResponse`
- New struct `CustomDomainsClientListOptions`
- New struct `CustomDomainsClientListResponse`
- New struct `EventHubEndpoint`
- New struct `EventListener`
- New struct `EventListenerEndpoint`
- New struct `EventListenerFilter`
- New struct `EventNameFilter`
- New struct `ResourceReference`
- New field `EventListeners` in struct `HubProperties`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/webpubsub/armwebpubsub` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).