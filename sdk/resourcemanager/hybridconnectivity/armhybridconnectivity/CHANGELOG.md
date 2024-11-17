# Release History

## 2.0.0 (2024-11-17)
### Breaking Changes

- Function `NewClientFactory` parameter(s) have been changed from `(azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Type of `ErrorAdditionalInfo.Info` has been changed from `any` to `*ErrorAdditionalInfoInfo`
- Enum `ActionType` has been removed
- Enum `Origin` has been removed
- Enum `ProvisioningState` has been removed
- Enum `ServiceName` has been removed
- Enum `Type` has been removed
- Function `*ClientFactory.NewEndpointsClient` has been removed
- Function `*ClientFactory.NewOperationsClient` has been removed
- Function `*ClientFactory.NewServiceConfigurationsClient` has been removed
- Function `NewEndpointsClient` has been removed
- Function `*EndpointsClient.CreateOrUpdate` has been removed
- Function `*EndpointsClient.Delete` has been removed
- Function `*EndpointsClient.Get` has been removed
- Function `*EndpointsClient.ListCredentials` has been removed
- Function `*EndpointsClient.ListIngressGatewayCredentials` has been removed
- Function `*EndpointsClient.ListManagedProxyDetails` has been removed
- Function `*EndpointsClient.NewListPager` has been removed
- Function `*EndpointsClient.Update` has been removed
- Function `NewOperationsClient` has been removed
- Function `*OperationsClient.NewListPager` has been removed
- Function `NewServiceConfigurationsClient` has been removed
- Function `*ServiceConfigurationsClient.CreateOrupdate` has been removed
- Function `*ServiceConfigurationsClient.Delete` has been removed
- Function `*ServiceConfigurationsClient.Get` has been removed
- Function `*ServiceConfigurationsClient.NewListByEndpointResourcePager` has been removed
- Function `*ServiceConfigurationsClient.Update` has been removed
- Struct `AADProfileProperties` has been removed
- Struct `EndpointAccessResource` has been removed
- Struct `EndpointProperties` has been removed
- Struct `EndpointResource` has been removed
- Struct `EndpointsList` has been removed
- Struct `ErrorResponse` has been removed
- Struct `IngressGatewayResource` has been removed
- Struct `IngressProfileProperties` has been removed
- Struct `ListCredentialsRequest` has been removed
- Struct `ListIngressGatewayCredentialsRequest` has been removed
- Struct `ManagedProxyRequest` has been removed
- Struct `ManagedProxyResource` has been removed
- Struct `Operation` has been removed
- Struct `OperationDisplay` has been removed
- Struct `OperationListResult` has been removed
- Struct `ProxyResource` has been removed
- Struct `RelayNamespaceAccessProperties` has been removed
- Struct `Resource` has been removed
- Struct `ServiceConfigurationList` has been removed
- Struct `ServiceConfigurationProperties` has been removed
- Struct `ServiceConfigurationPropertiesPatch` has been removed
- Struct `ServiceConfigurationResource` has been removed
- Struct `ServiceConfigurationResourcePatch` has been removed

### Features Added

- New enum type `CloudNativeType` with values `CloudNativeTypeEc2`
- New enum type `HostType` with values `HostTypeAWS`
- New enum type `ResourceProvisioningState` with values `ResourceProvisioningStateCanceled`, `ResourceProvisioningStateFailed`, `ResourceProvisioningStateSucceeded`
- New enum type `SolutionConfigurationStatus` with values `SolutionConfigurationStatusCompleted`, `SolutionConfigurationStatusFailed`, `SolutionConfigurationStatusInProgress`, `SolutionConfigurationStatusNew`
- New function `*ClientFactory.NewGenerateAwsTemplateClient() *GenerateAwsTemplateClient`
- New function `*ClientFactory.NewInventoryClient() *InventoryClient`
- New function `*ClientFactory.NewPublicCloudConnectorsClient() *PublicCloudConnectorsClient`
- New function `*ClientFactory.NewSolutionConfigurationsClient() *SolutionConfigurationsClient`
- New function `*ClientFactory.NewSolutionTypesClient() *SolutionTypesClient`
- New function `NewGenerateAwsTemplateClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GenerateAwsTemplateClient, error)`
- New function `*GenerateAwsTemplateClient.Post(context.Context, GenerateAwsTemplateRequest, *GenerateAwsTemplateClientPostOptions) (GenerateAwsTemplateClientPostResponse, error)`
- New function `NewInventoryClient(azcore.TokenCredential, *arm.ClientOptions) (*InventoryClient, error)`
- New function `*InventoryClient.Get(context.Context, string, string, string, *InventoryClientGetOptions) (InventoryClientGetResponse, error)`
- New function `*InventoryClient.NewListBySolutionConfigurationPager(string, string, *InventoryClientListBySolutionConfigurationOptions) *runtime.Pager[InventoryClientListBySolutionConfigurationResponse]`
- New function `NewPublicCloudConnectorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PublicCloudConnectorsClient, error)`
- New function `*PublicCloudConnectorsClient.BeginCreateOrUpdate(context.Context, string, string, PublicCloudConnector, *PublicCloudConnectorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[PublicCloudConnectorsClientCreateOrUpdateResponse], error)`
- New function `*PublicCloudConnectorsClient.Delete(context.Context, string, string, *PublicCloudConnectorsClientDeleteOptions) (PublicCloudConnectorsClientDeleteResponse, error)`
- New function `*PublicCloudConnectorsClient.Get(context.Context, string, string, *PublicCloudConnectorsClientGetOptions) (PublicCloudConnectorsClientGetResponse, error)`
- New function `*PublicCloudConnectorsClient.NewListByResourceGroupPager(string, *PublicCloudConnectorsClientListByResourceGroupOptions) *runtime.Pager[PublicCloudConnectorsClientListByResourceGroupResponse]`
- New function `*PublicCloudConnectorsClient.NewListBySubscriptionPager(*PublicCloudConnectorsClientListBySubscriptionOptions) *runtime.Pager[PublicCloudConnectorsClientListBySubscriptionResponse]`
- New function `*PublicCloudConnectorsClient.BeginTestPermissions(context.Context, string, string, *PublicCloudConnectorsClientBeginTestPermissionsOptions) (*runtime.Poller[PublicCloudConnectorsClientTestPermissionsResponse], error)`
- New function `*PublicCloudConnectorsClient.Update(context.Context, string, string, PublicCloudConnector, *PublicCloudConnectorsClientUpdateOptions) (PublicCloudConnectorsClientUpdateResponse, error)`
- New function `NewSolutionConfigurationsClient(azcore.TokenCredential, *arm.ClientOptions) (*SolutionConfigurationsClient, error)`
- New function `*SolutionConfigurationsClient.CreateOrUpdate(context.Context, string, string, SolutionConfiguration, *SolutionConfigurationsClientCreateOrUpdateOptions) (SolutionConfigurationsClientCreateOrUpdateResponse, error)`
- New function `*SolutionConfigurationsClient.Delete(context.Context, string, string, *SolutionConfigurationsClientDeleteOptions) (SolutionConfigurationsClientDeleteResponse, error)`
- New function `*SolutionConfigurationsClient.Get(context.Context, string, string, *SolutionConfigurationsClientGetOptions) (SolutionConfigurationsClientGetResponse, error)`
- New function `*SolutionConfigurationsClient.NewListPager(string, *SolutionConfigurationsClientListOptions) *runtime.Pager[SolutionConfigurationsClientListResponse]`
- New function `*SolutionConfigurationsClient.BeginSyncNow(context.Context, string, string, *SolutionConfigurationsClientBeginSyncNowOptions) (*runtime.Poller[SolutionConfigurationsClientSyncNowResponse], error)`
- New function `*SolutionConfigurationsClient.Update(context.Context, string, string, SolutionConfiguration, *SolutionConfigurationsClientUpdateOptions) (SolutionConfigurationsClientUpdateResponse, error)`
- New function `NewSolutionTypesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SolutionTypesClient, error)`
- New function `*SolutionTypesClient.Get(context.Context, string, string, *SolutionTypesClientGetOptions) (SolutionTypesClientGetResponse, error)`
- New function `*SolutionTypesClient.NewListByResourceGroupPager(string, *SolutionTypesClientListByResourceGroupOptions) *runtime.Pager[SolutionTypesClientListByResourceGroupResponse]`
- New function `*SolutionTypesClient.NewListBySubscriptionPager(*SolutionTypesClientListBySubscriptionOptions) *runtime.Pager[SolutionTypesClientListBySubscriptionResponse]`
- New struct `AwsCloudProfile`
- New struct `ErrorAdditionalInfoInfo`
- New struct `GenerateAwsTemplateRequest`
- New struct `InventoryProperties`
- New struct `InventoryResource`
- New struct `InventoryResourceListResult`
- New struct `OperationStatusResult`
- New struct `PostResponse`
- New struct `PublicCloudConnector`
- New struct `PublicCloudConnectorListResult`
- New struct `PublicCloudConnectorProperties`
- New struct `SolutionConfiguration`
- New struct `SolutionConfigurationListResult`
- New struct `SolutionConfigurationProperties`
- New struct `SolutionSettings`
- New struct `SolutionTypeProperties`
- New struct `SolutionTypeResource`
- New struct `SolutionTypeResourceListResult`
- New struct `SolutionTypeSettings`
- New struct `SolutionTypeSettingsProperties`


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