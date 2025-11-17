# Release History

## 2.0.0 (2025-11-17)
### Breaking Changes

- Operation `*LocalRulestacksClient.ListAppIDs` has supported pagination, use `*LocalRulestacksClient.NewListAppIDsPager` instead.
- Operation `*LocalRulestacksClient.ListCountries` has supported pagination, use `*LocalRulestacksClient.NewListCountriesPager` instead.
- Operation `*LocalRulestacksClient.ListPredefinedURLCategories` has supported pagination, use `*LocalRulestacksClient.NewListPredefinedURLCategoriesPager` instead.

### Features Added

- New enum type `EnableStatus` with values `EnableStatusDisabled`, `EnableStatusEnabled`
- New enum type `ProductSerialStatusValues` with values `ProductSerialStatusValuesAllocated`, `ProductSerialStatusValuesInProgress`
- New enum type `RegistrationStatus` with values `RegistrationStatusNotRegistered`, `RegistrationStatusRegistered`
- New function `*ClientFactory.NewMetricsObjectFirewallClient() *MetricsObjectFirewallClient`
- New function `*ClientFactory.NewPaloAltoNetworksCloudngfwOperationsClient() *PaloAltoNetworksCloudngfwOperationsClient`
- New function `NewMetricsObjectFirewallClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MetricsObjectFirewallClient, error)`
- New function `*MetricsObjectFirewallClient.BeginCreateOrUpdate(context.Context, string, string, MetricsObjectFirewallResource, *MetricsObjectFirewallClientBeginCreateOrUpdateOptions) (*runtime.Poller[MetricsObjectFirewallClientCreateOrUpdateResponse], error)`
- New function `*MetricsObjectFirewallClient.BeginDelete(context.Context, string, string, *MetricsObjectFirewallClientBeginDeleteOptions) (*runtime.Poller[MetricsObjectFirewallClientDeleteResponse], error)`
- New function `*MetricsObjectFirewallClient.Get(context.Context, string, string, *MetricsObjectFirewallClientGetOptions) (MetricsObjectFirewallClientGetResponse, error)`
- New function `*MetricsObjectFirewallClient.NewListByFirewallsPager(string, string, *MetricsObjectFirewallClientListByFirewallsOptions) *runtime.Pager[MetricsObjectFirewallClientListByFirewallsResponse]`
- New function `NewPaloAltoNetworksCloudngfwOperationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PaloAltoNetworksCloudngfwOperationsClient, error)`
- New function `*PaloAltoNetworksCloudngfwOperationsClient.CreateProductSerialNumber(context.Context, *PaloAltoNetworksCloudngfwOperationsClientCreateProductSerialNumberOptions) (PaloAltoNetworksCloudngfwOperationsClientCreateProductSerialNumberResponse, error)`
- New function `*PaloAltoNetworksCloudngfwOperationsClient.ListCloudManagerTenants(context.Context, *PaloAltoNetworksCloudngfwOperationsClientListCloudManagerTenantsOptions) (PaloAltoNetworksCloudngfwOperationsClientListCloudManagerTenantsResponse, error)`
- New function `*PaloAltoNetworksCloudngfwOperationsClient.ListProductSerialNumberStatus(context.Context, *PaloAltoNetworksCloudngfwOperationsClientListProductSerialNumberStatusOptions) (PaloAltoNetworksCloudngfwOperationsClientListProductSerialNumberStatusResponse, error)`
- New function `*PaloAltoNetworksCloudngfwOperationsClient.ListSupportInfo(context.Context, *PaloAltoNetworksCloudngfwOperationsClientListSupportInfoOptions) (PaloAltoNetworksCloudngfwOperationsClientListSupportInfoResponse, error)`
- New struct `CloudManagerTenantList`
- New struct `MetricsObject`
- New struct `MetricsObjectFirewallResource`
- New struct `MetricsObjectFirewallResourceListResult`
- New struct `ProductSerialNumberRequestStatus`
- New struct `ProductSerialNumberStatus`
- New struct `StrataCloudManagerConfig`
- New struct `StrataCloudManagerInfo`
- New struct `SupportInfoModel`
- New field `IsStrataCloudManaged`, `StrataCloudManagerConfig` in struct `FirewallDeploymentProperties`
- New field `IsStrataCloudManaged`, `StrataCloudManagerConfig` in struct `FirewallResourceUpdateProperties`
- New field `IsStrataCloudManaged`, `StrataCloudManagerInfo` in struct `FirewallStatusProperty`
- New field `PrivateSourceNatRulesDestination` in struct `NetworkProfile`


## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.
- New field `TrustedRanges` in struct `NetworkProfile`


## 1.0.0 (2023-07-14)
### Other Changes

- Release stable version.

## 0.1.0 (2023-04-28)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/paloaltonetworksngfw/armpanngfw` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).