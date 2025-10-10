# Release History

## 2.0.0 (2025-10-10)
### Breaking Changes

- Operation `*GrafanaClient.Update` has been changed to LRO, use `*GrafanaClient.BeginUpdate` instead.

### Features Added

- New enum type `CreatorCanAdmin` with values `CreatorCanAdminDisabled`, `CreatorCanAdminEnabled`
- New enum type `Size` with values `SizeX1`, `SizeX2`
- New function `*ClientFactory.NewIntegrationFabricsClient() *IntegrationFabricsClient`
- New function `*ClientFactory.NewManagedDashboardsClient() *ManagedDashboardsClient`
- New function `NewIntegrationFabricsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*IntegrationFabricsClient, error)`
- New function `*IntegrationFabricsClient.BeginCreate(context.Context, string, string, string, IntegrationFabric, *IntegrationFabricsClientBeginCreateOptions) (*runtime.Poller[IntegrationFabricsClientCreateResponse], error)`
- New function `*IntegrationFabricsClient.BeginDelete(context.Context, string, string, string, *IntegrationFabricsClientBeginDeleteOptions) (*runtime.Poller[IntegrationFabricsClientDeleteResponse], error)`
- New function `*IntegrationFabricsClient.Get(context.Context, string, string, string, *IntegrationFabricsClientGetOptions) (IntegrationFabricsClientGetResponse, error)`
- New function `*IntegrationFabricsClient.NewListPager(string, string, *IntegrationFabricsClientListOptions) *runtime.Pager[IntegrationFabricsClientListResponse]`
- New function `*IntegrationFabricsClient.BeginUpdate(context.Context, string, string, string, IntegrationFabricUpdateParameters, *IntegrationFabricsClientBeginUpdateOptions) (*runtime.Poller[IntegrationFabricsClientUpdateResponse], error)`
- New function `NewManagedDashboardsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedDashboardsClient, error)`
- New function `*ManagedDashboardsClient.BeginCreate(context.Context, string, string, ManagedDashboard, *ManagedDashboardsClientBeginCreateOptions) (*runtime.Poller[ManagedDashboardsClientCreateResponse], error)`
- New function `*ManagedDashboardsClient.Delete(context.Context, string, string, *ManagedDashboardsClientDeleteOptions) (ManagedDashboardsClientDeleteResponse, error)`
- New function `*ManagedDashboardsClient.Get(context.Context, string, string, *ManagedDashboardsClientGetOptions) (ManagedDashboardsClientGetResponse, error)`
- New function `*ManagedDashboardsClient.NewListBySubscriptionPager(*ManagedDashboardsClientListBySubscriptionOptions) *runtime.Pager[ManagedDashboardsClientListBySubscriptionResponse]`
- New function `*ManagedDashboardsClient.NewListPager(string, *ManagedDashboardsClientListOptions) *runtime.Pager[ManagedDashboardsClientListResponse]`
- New function `*ManagedDashboardsClient.Update(context.Context, string, string, ManagedDashboardUpdateParameters, *ManagedDashboardsClientUpdateOptions) (ManagedDashboardsClientUpdateResponse, error)`
- New struct `IntegrationFabric`
- New struct `IntegrationFabricListResponse`
- New struct `IntegrationFabricProperties`
- New struct `IntegrationFabricPropertiesUpdateParameters`
- New struct `IntegrationFabricUpdateParameters`
- New struct `ManagedDashboard`
- New struct `ManagedDashboardListResponse`
- New struct `ManagedDashboardProperties`
- New struct `ManagedDashboardUpdateParameters`
- New struct `Security`
- New struct `Snapshots`
- New struct `UnifiedAlertingScreenshots`
- New struct `Users`
- New field `Author`, `Type` in struct `GrafanaAvailablePlugin`
- New field `Security`, `Snapshots`, `UnifiedAlertingScreenshots`, `Users` in struct `GrafanaConfigurations`
- New field `CreatorCanAdmin` in struct `ManagedGrafanaProperties`
- New field `CreatorCanAdmin` in struct `ManagedGrafanaPropertiesUpdateParameters`
- New field `Size` in struct `ResourceSKU`


## 2.0.0-beta.1 (2025-07-18)
### Breaking Changes

- Operation `*GrafanaClient.Update` has been changed to LRO, use `*GrafanaClient.BeginUpdate` instead.

### Features Added

- New function `*ClientFactory.NewIntegrationFabricsClient() *IntegrationFabricsClient`
- New function `*ClientFactory.NewManagedDashboardsClient() *ManagedDashboardsClient`
- New function `NewIntegrationFabricsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*IntegrationFabricsClient, error)`
- New function `*IntegrationFabricsClient.BeginCreate(context.Context, string, string, string, IntegrationFabric, *IntegrationFabricsClientBeginCreateOptions) (*runtime.Poller[IntegrationFabricsClientCreateResponse], error)`
- New function `*IntegrationFabricsClient.BeginDelete(context.Context, string, string, string, *IntegrationFabricsClientBeginDeleteOptions) (*runtime.Poller[IntegrationFabricsClientDeleteResponse], error)`
- New function `*IntegrationFabricsClient.Get(context.Context, string, string, string, *IntegrationFabricsClientGetOptions) (IntegrationFabricsClientGetResponse, error)`
- New function `*IntegrationFabricsClient.NewListPager(string, string, *IntegrationFabricsClientListOptions) *runtime.Pager[IntegrationFabricsClientListResponse]`
- New function `*IntegrationFabricsClient.BeginUpdate(context.Context, string, string, string, IntegrationFabricUpdateParameters, *IntegrationFabricsClientBeginUpdateOptions) (*runtime.Poller[IntegrationFabricsClientUpdateResponse], error)`
- New function `NewManagedDashboardsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedDashboardsClient, error)`
- New function `*ManagedDashboardsClient.BeginCreate(context.Context, string, string, ManagedDashboard, *ManagedDashboardsClientBeginCreateOptions) (*runtime.Poller[ManagedDashboardsClientCreateResponse], error)`
- New function `*ManagedDashboardsClient.Delete(context.Context, string, string, *ManagedDashboardsClientDeleteOptions) (ManagedDashboardsClientDeleteResponse, error)`
- New function `*ManagedDashboardsClient.Get(context.Context, string, string, *ManagedDashboardsClientGetOptions) (ManagedDashboardsClientGetResponse, error)`
- New function `*ManagedDashboardsClient.NewListBySubscriptionPager(*ManagedDashboardsClientListBySubscriptionOptions) *runtime.Pager[ManagedDashboardsClientListBySubscriptionResponse]`
- New function `*ManagedDashboardsClient.NewListPager(string, *ManagedDashboardsClientListOptions) *runtime.Pager[ManagedDashboardsClientListResponse]`
- New function `*ManagedDashboardsClient.Update(context.Context, string, string, ManagedDashboardUpdateParameters, *ManagedDashboardsClientUpdateOptions) (ManagedDashboardsClientUpdateResponse, error)`
- New struct `IntegrationFabric`
- New struct `IntegrationFabricListResponse`
- New struct `IntegrationFabricProperties`
- New struct `IntegrationFabricPropertiesUpdateParameters`
- New struct `IntegrationFabricUpdateParameters`
- New struct `ManagedDashboard`
- New struct `ManagedDashboardListResponse`
- New struct `ManagedDashboardProperties`
- New struct `ManagedDashboardUpdateParameters`
- New struct `Security`
- New struct `Snapshots`
- New struct `UnifiedAlertingScreenshots`
- New struct `Users`
- New field `Security`, `Snapshots`, `UnifiedAlertingScreenshots`, `Users` in struct `GrafanaConfigurations`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.
- New enum type `AvailablePromotion` with values `AvailablePromotionFreeTrial`, `AvailablePromotionNone`
- New enum type `ManagedPrivateEndpointConnectionStatus` with values `ManagedPrivateEndpointConnectionStatusApproved`, `ManagedPrivateEndpointConnectionStatusDisconnected`, `ManagedPrivateEndpointConnectionStatusPending`, `ManagedPrivateEndpointConnectionStatusRejected`
- New enum type `MarketplaceAutoRenew` with values `MarketplaceAutoRenewDisabled`, `MarketplaceAutoRenewEnabled`
- New enum type `StartTLSPolicy` with values `StartTLSPolicyMandatoryStartTLS`, `StartTLSPolicyNoStartTLS`, `StartTLSPolicyOpportunisticStartTLS`
- New function `*ClientFactory.NewManagedPrivateEndpointsClient() *ManagedPrivateEndpointsClient`
- New function `*GrafanaClient.CheckEnterpriseDetails(context.Context, string, string, *GrafanaClientCheckEnterpriseDetailsOptions) (GrafanaClientCheckEnterpriseDetailsResponse, error)`
- New function `*GrafanaClient.FetchAvailablePlugins(context.Context, string, string, *GrafanaClientFetchAvailablePluginsOptions) (GrafanaClientFetchAvailablePluginsResponse, error)`
- New function `NewManagedPrivateEndpointsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedPrivateEndpointsClient, error)`
- New function `*ManagedPrivateEndpointsClient.BeginCreate(context.Context, string, string, string, ManagedPrivateEndpointModel, *ManagedPrivateEndpointsClientBeginCreateOptions) (*runtime.Poller[ManagedPrivateEndpointsClientCreateResponse], error)`
- New function `*ManagedPrivateEndpointsClient.BeginDelete(context.Context, string, string, string, *ManagedPrivateEndpointsClientBeginDeleteOptions) (*runtime.Poller[ManagedPrivateEndpointsClientDeleteResponse], error)`
- New function `*ManagedPrivateEndpointsClient.Get(context.Context, string, string, string, *ManagedPrivateEndpointsClientGetOptions) (ManagedPrivateEndpointsClientGetResponse, error)`
- New function `*ManagedPrivateEndpointsClient.NewListPager(string, string, *ManagedPrivateEndpointsClientListOptions) *runtime.Pager[ManagedPrivateEndpointsClientListResponse]`
- New function `*ManagedPrivateEndpointsClient.BeginRefresh(context.Context, string, string, *ManagedPrivateEndpointsClientBeginRefreshOptions) (*runtime.Poller[ManagedPrivateEndpointsClientRefreshResponse], error)`
- New function `*ManagedPrivateEndpointsClient.BeginUpdate(context.Context, string, string, string, ManagedPrivateEndpointUpdateParameters, *ManagedPrivateEndpointsClientBeginUpdateOptions) (*runtime.Poller[ManagedPrivateEndpointsClientUpdateResponse], error)`
- New struct `EnterpriseConfigurations`
- New struct `EnterpriseDetails`
- New struct `GrafanaAvailablePlugin`
- New struct `GrafanaAvailablePluginListResponse`
- New struct `GrafanaConfigurations`
- New struct `GrafanaPlugin`
- New struct `ManagedPrivateEndpointConnectionState`
- New struct `ManagedPrivateEndpointModel`
- New struct `ManagedPrivateEndpointModelListResponse`
- New struct `ManagedPrivateEndpointModelProperties`
- New struct `ManagedPrivateEndpointUpdateParameters`
- New struct `MarketplaceTrialQuota`
- New struct `SMTP`
- New struct `SaasSubscriptionDetails`
- New struct `SubscriptionTerm`
- New field `AzureAsyncOperation` in struct `GrafanaClientUpdateResponse`
- New field `EnterpriseConfigurations`, `GrafanaConfigurations`, `GrafanaMajorVersion`, `GrafanaPlugins` in struct `ManagedGrafanaProperties`
- New field `EnterpriseConfigurations`, `GrafanaConfigurations`, `GrafanaMajorVersion`, `GrafanaPlugins` in struct `ManagedGrafanaPropertiesUpdateParameters`
- New field `SKU` in struct `ManagedGrafanaUpdateParameters`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-08-09)
### Breaking Changes

- Type of `ManagedGrafanaUpdateParameters.Identity` has been changed from `*ManagedIdentity` to `*ManagedServiceIdentity`
- Type of `SystemData.LastModifiedByType` has been changed from `*LastModifiedByType` to `*CreatedByType`
- Type of `OperationListResult.Value` has been changed from `[]*OperationResult` to `[]*Operation`
- Type of `ManagedGrafana.Identity` has been changed from `*ManagedIdentity` to `*ManagedServiceIdentity`
- Const `LastModifiedByTypeUser` has been removed
- Const `LastModifiedByTypeApplication` has been removed
- Const `LastModifiedByTypeKey` has been removed
- Const `LastModifiedByTypeManagedIdentity` has been removed
- Const `IdentityTypeNone` has been removed
- Const `IdentityTypeSystemAssigned` has been removed
- Function `PossibleLastModifiedByTypeValues` has been removed
- Function `PossibleIdentityTypeValues` has been removed
- Struct `ManagedIdentity` has been removed
- Struct `OperationResult` has been removed

### Features Added

- New const `APIKeyDisabled`
- New const `PrivateEndpointConnectionProvisioningStateDeleting`
- New const `DeterministicOutboundIPDisabled`
- New const `PublicNetworkAccessEnabled`
- New const `PrivateEndpointConnectionProvisioningStateSucceeded`
- New const `PrivateEndpointConnectionProvisioningStateFailed`
- New const `ManagedServiceIdentityTypeSystemAssignedUserAssigned`
- New const `PrivateEndpointServiceConnectionStatusPending`
- New const `PublicNetworkAccessDisabled`
- New const `ManagedServiceIdentityTypeUserAssigned`
- New const `PrivateEndpointServiceConnectionStatusRejected`
- New const `DeterministicOutboundIPEnabled`
- New const `ManagedServiceIdentityTypeSystemAssigned`
- New const `ManagedServiceIdentityTypeNone`
- New const `PrivateEndpointConnectionProvisioningStateCreating`
- New const `APIKeyEnabled`
- New const `PrivateEndpointServiceConnectionStatusApproved`
- New function `NewPrivateLinkResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateLinkResourcesClient, error)`
- New function `PossibleAPIKeyValues() []APIKey`
- New function `*PrivateEndpointConnectionsClient.BeginApprove(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientBeginApproveOptions) (*runtime.Poller[PrivateEndpointConnectionsClientApproveResponse], error)`
- New function `PossibleDeterministicOutboundIPValues() []DeterministicOutboundIP`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionsClient.NewListPager(string, string, *PrivateEndpointConnectionsClientListOptions) *runtime.Pager[PrivateEndpointConnectionsClientListResponse]`
- New function `PossiblePublicNetworkAccessValues() []PublicNetworkAccess`
- New function `*PrivateLinkResourcesClient.Get(context.Context, string, string, string, *PrivateLinkResourcesClientGetOptions) (PrivateLinkResourcesClientGetResponse, error)`
- New function `PossibleManagedServiceIdentityTypeValues() []ManagedServiceIdentityType`
- New function `PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState`
- New function `*PrivateLinkResourcesClient.NewListPager(string, string, *PrivateLinkResourcesClientListOptions) *runtime.Pager[PrivateLinkResourcesClientListResponse]`
- New function `*PrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `PossiblePrivateEndpointServiceConnectionStatusValues() []PrivateEndpointServiceConnectionStatus`
- New struct `AzureMonitorWorkspaceIntegration`
- New struct `GrafanaIntegrations`
- New struct `ManagedGrafanaPropertiesUpdateParameters`
- New struct `ManagedServiceIdentity`
- New struct `Operation`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateEndpointConnectionsClient`
- New struct `PrivateEndpointConnectionsClientApproveResponse`
- New struct `PrivateEndpointConnectionsClientBeginApproveOptions`
- New struct `PrivateEndpointConnectionsClientBeginDeleteOptions`
- New struct `PrivateEndpointConnectionsClientDeleteResponse`
- New struct `PrivateEndpointConnectionsClientGetOptions`
- New struct `PrivateEndpointConnectionsClientGetResponse`
- New struct `PrivateEndpointConnectionsClientListOptions`
- New struct `PrivateEndpointConnectionsClientListResponse`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkResourcesClient`
- New struct `PrivateLinkResourcesClientGetOptions`
- New struct `PrivateLinkResourcesClientGetResponse`
- New struct `PrivateLinkResourcesClientListOptions`
- New struct `PrivateLinkResourcesClientListResponse`
- New struct `PrivateLinkServiceConnectionState`
- New struct `Resource`
- New field `OutboundIPs` in struct `ManagedGrafanaProperties`
- New field `APIKey` in struct `ManagedGrafanaProperties`
- New field `GrafanaIntegrations` in struct `ManagedGrafanaProperties`
- New field `PublicNetworkAccess` in struct `ManagedGrafanaProperties`
- New field `DeterministicOutboundIP` in struct `ManagedGrafanaProperties`
- New field `PrivateEndpointConnections` in struct `ManagedGrafanaProperties`
- New field `Properties` in struct `ManagedGrafanaUpdateParameters`


## 0.3.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dashboard/armdashboard` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.3.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
