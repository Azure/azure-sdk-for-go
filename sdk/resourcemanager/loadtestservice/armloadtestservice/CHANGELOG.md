# Release History

## 1.0.1 (2022-11-21)

- Deprecated: use github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/loadtesting/armloadtesting instead.

## 1.0.0 (2022-10-13)
### Features Added

- New function `NewQuotasClient(string, azcore.TokenCredential, *arm.ClientOptions) (*QuotasClient, error)`
- New function `*QuotasClient.CheckAvailability(context.Context, string, string, QuotaBucketRequest, *QuotasClientCheckAvailabilityOptions) (QuotasClientCheckAvailabilityResponse, error)`
- New function `*QuotasClient.NewListPager(string, *QuotasClientListOptions) *runtime.Pager[QuotasClientListResponse]`
- New function `*LoadTestsClient.NewListOutboundNetworkDependenciesEndpointsPager(string, string, *LoadTestsClientListOutboundNetworkDependenciesEndpointsOptions) *runtime.Pager[LoadTestsClientListOutboundNetworkDependenciesEndpointsResponse]`
- New function `*QuotasClient.Get(context.Context, string, string, *QuotasClientGetOptions) (QuotasClientGetResponse, error)`
- New struct `CheckQuotaAvailabilityResponse`
- New struct `CheckQuotaAvailabilityResponseProperties`
- New struct `EndpointDependency`
- New struct `EndpointDetail`
- New struct `LoadTestsClientListOutboundNetworkDependenciesEndpointsOptions`
- New struct `LoadTestsClientListOutboundNetworkDependenciesEndpointsResponse`
- New struct `OutboundEnvironmentEndpoint`
- New struct `OutboundEnvironmentEndpointCollection`
- New struct `QuotaBucketRequest`
- New struct `QuotaBucketRequestProperties`
- New struct `QuotaBucketRequestPropertiesDimensions`
- New struct `QuotaResource`
- New struct `QuotaResourceList`
- New struct `QuotaResourceProperties`
- New struct `QuotasClient`
- New struct `QuotasClientCheckAvailabilityOptions`
- New struct `QuotasClientCheckAvailabilityResponse`
- New struct `QuotasClientGetOptions`
- New struct `QuotasClientGetResponse`
- New struct `QuotasClientListOptions`
- New struct `QuotasClientListResponse`


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/loadtestservice/armloadtestservice` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).