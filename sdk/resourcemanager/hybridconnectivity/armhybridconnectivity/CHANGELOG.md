# Release History

## 1.1.1 (2024-11-15)
### Other Changes


## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.0.0 (2023-09-22)
### Features Added

- New enum type `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateCreating`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`, `ProvisioningStateUpdating`
- New enum type `ServiceName` with values `ServiceNameSSH`, `ServiceNameWAC`
- New function `*ClientFactory.NewServiceConfigurationsClient() *ServiceConfigurationsClient`
- New function `*EndpointsClient.ListIngressGatewayCredentials(context.Context, string, string, *EndpointsClientListIngressGatewayCredentialsOptions) (EndpointsClientListIngressGatewayCredentialsResponse, error)`
- New function `NewServiceConfigurationsClient(azcore.TokenCredential, *arm.ClientOptions) (*ServiceConfigurationsClient, error)`
- New function `*ServiceConfigurationsClient.CreateOrupdate(context.Context, string, string, string, ServiceConfigurationResource, *ServiceConfigurationsClientCreateOrupdateOptions) (ServiceConfigurationsClientCreateOrupdateResponse, error)`
- New function `*ServiceConfigurationsClient.Delete(context.Context, string, string, string, *ServiceConfigurationsClientDeleteOptions) (ServiceConfigurationsClientDeleteResponse, error)`
- New function `*ServiceConfigurationsClient.Get(context.Context, string, string, string, *ServiceConfigurationsClientGetOptions) (ServiceConfigurationsClientGetResponse, error)`
- New function `*ServiceConfigurationsClient.NewListByEndpointResourcePager(string, string, *ServiceConfigurationsClientListByEndpointResourceOptions) *runtime.Pager[ServiceConfigurationsClientListByEndpointResourceResponse]`
- New function `*ServiceConfigurationsClient.Update(context.Context, string, string, string, ServiceConfigurationResourcePatch, *ServiceConfigurationsClientUpdateOptions) (ServiceConfigurationsClientUpdateResponse, error)`
- New struct `ListCredentialsRequest`
- New struct `ListIngressGatewayCredentialsRequest`
- New struct `ServiceConfigurationList`
- New struct `ServiceConfigurationProperties`
- New struct `ServiceConfigurationPropertiesPatch`
- New struct `ServiceConfigurationResource`
- New struct `ServiceConfigurationResourcePatch`
- New field `ListCredentialsRequest` in struct `EndpointsClientListCredentialsOptions`
- New field `ServiceName` in struct `ManagedProxyRequest`
- New field `SystemData` in struct `ProxyResource`
- New field `ServiceConfigurationToken` in struct `RelayNamespaceAccessProperties`
- New field `SystemData` in struct `Resource`


## 0.6.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.6.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridconnectivity/armhybridconnectivity` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).