# Release History

## 2.0.0 (2025-06-16)
### Breaking Changes

- Struct `ErrorResponse` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed

### Features Added

- New enum type `CloudNativeType` with values `CloudNativeTypeEc2`
- New enum type `HostType` with values `HostTypeAWS`
- New enum type `ResourceProvisioningState` with values `ResourceProvisioningStateCanceled`, `ResourceProvisioningStateFailed`, `ResourceProvisioningStateSucceeded`
- New enum type `SolutionConfigurationStatus` with values `SolutionConfigurationStatusCompleted`, `SolutionConfigurationStatusFailed`, `SolutionConfigurationStatusInProgress`, `SolutionConfigurationStatusNew`
- New function `*ClientFactory.NewGenerateAwsTemplateClient(string) *GenerateAwsTemplateClient`
- New function `*ClientFactory.NewInventoryClient() *InventoryClient`
- New function `*ClientFactory.NewPublicCloudConnectorsClient(string) *PublicCloudConnectorsClient`
- New function `*ClientFactory.NewSolutionConfigurationsClient() *SolutionConfigurationsClient`
- New function `*ClientFactory.NewSolutionTypesClient(string) *SolutionTypesClient`
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
- New function `*PublicCloudConnectorsClient.Update(context.Context, string, string, PublicCloudConnectorUpdate, *PublicCloudConnectorsClientUpdateOptions) (PublicCloudConnectorsClientUpdateResponse, error)`
- New function `NewSolutionConfigurationsClient(azcore.TokenCredential, *arm.ClientOptions) (*SolutionConfigurationsClient, error)`
- New function `*SolutionConfigurationsClient.CreateOrUpdate(context.Context, string, string, SolutionConfiguration, *SolutionConfigurationsClientCreateOrUpdateOptions) (SolutionConfigurationsClientCreateOrUpdateResponse, error)`
- New function `*SolutionConfigurationsClient.Delete(context.Context, string, string, *SolutionConfigurationsClientDeleteOptions) (SolutionConfigurationsClientDeleteResponse, error)`
- New function `*SolutionConfigurationsClient.Get(context.Context, string, string, *SolutionConfigurationsClientGetOptions) (SolutionConfigurationsClientGetResponse, error)`
- New function `*SolutionConfigurationsClient.NewListPager(string, *SolutionConfigurationsClientListOptions) *runtime.Pager[SolutionConfigurationsClientListResponse]`
- New function `*SolutionConfigurationsClient.BeginSyncNow(context.Context, string, string, *SolutionConfigurationsClientBeginSyncNowOptions) (*runtime.Poller[SolutionConfigurationsClientSyncNowResponse], error)`
- New function `*SolutionConfigurationsClient.Update(context.Context, string, string, SolutionConfigurationUpdate, *SolutionConfigurationsClientUpdateOptions) (SolutionConfigurationsClientUpdateResponse, error)`
- New function `NewSolutionTypesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SolutionTypesClient, error)`
- New function `*SolutionTypesClient.Get(context.Context, string, string, *SolutionTypesClientGetOptions) (SolutionTypesClientGetResponse, error)`
- New function `*SolutionTypesClient.NewListByResourceGroupPager(string, *SolutionTypesClientListByResourceGroupOptions) *runtime.Pager[SolutionTypesClientListByResourceGroupResponse]`
- New function `*SolutionTypesClient.NewListBySubscriptionPager(*SolutionTypesClientListBySubscriptionOptions) *runtime.Pager[SolutionTypesClientListBySubscriptionResponse]`
- New struct `AwsCloudProfile`
- New struct `AwsCloudProfileUpdate`
- New struct `GenerateAwsTemplateRequest`
- New struct `GenerateAwsTemplateResponse`
- New struct `InventoryProperties`
- New struct `InventoryResource`
- New struct `InventoryResourceListResult`
- New struct `OperationStatusResult`
- New struct `PublicCloudConnector`
- New struct `PublicCloudConnectorListResult`
- New struct `PublicCloudConnectorProperties`
- New struct `PublicCloudConnectorPropertiesUpdate`
- New struct `PublicCloudConnectorUpdate`
- New struct `SolutionConfiguration`
- New struct `SolutionConfigurationListResult`
- New struct `SolutionConfigurationProperties`
- New struct `SolutionConfigurationPropertiesUpdate`
- New struct `SolutionConfigurationUpdate`
- New struct `SolutionSettings`
- New struct `SolutionTypeProperties`
- New struct `SolutionTypeResource`
- New struct `SolutionTypeResourceListResult`
- New struct `SolutionTypeSettings`
- New struct `SolutionTypeSettingsProperties`


## 1.2.0-beta.1 (2025-02-27)

### Features Added

- New enum type `CloudNativeType` with values `CloudNativeTypeEc2`
- New enum type `HostType` with values `HostTypeAWS`
- New enum type `ResourceProvisioningState` with values `ResourceProvisioningStateCanceled`, `ResourceProvisioningStateFailed`, `ResourceProvisioningStateSucceeded`
- New enum type `SolutionConfigurationStatus` with values `SolutionConfigurationStatusCompleted`, `SolutionConfigurationStatusFailed`, `SolutionConfigurationStatusInProgress`, `SolutionConfigurationStatusNew`
- New function `*ClientFactory.NewGenerateAwsTemplateClient(string) *GenerateAwsTemplateClient`
- New function `*ClientFactory.NewInventoryClient() *InventoryClient`
- New function `*ClientFactory.NewPublicCloudConnectorsClient(string) *PublicCloudConnectorsClient`
- New function `*ClientFactory.NewSolutionConfigurationsClient() *SolutionConfigurationsClient`
- New function `*ClientFactory.NewSolutionTypesClient(string) *SolutionTypesClient`
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