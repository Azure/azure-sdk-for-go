# Release History

## 2.0.0 (2025-10-23)
### Breaking Changes

- Struct `ErrorAdditionalInfoInfo` has been removed

### Features Added

- Type of `ErrorAdditionalInfo.Info` has been changed from `*ErrorAdditionalInfoInfo` to `any`
- New enum type `DatasetDestinationTarget` with values `DatasetDestinationTargetBrokerStateStore`, `DatasetDestinationTargetMqtt`, `DatasetDestinationTargetStorage`
- New enum type `EventDestinationTarget` with values `EventDestinationTargetMqtt`, `EventDestinationTargetStorage`
- New enum type `Format` with values `FormatDelta10`, `FormatJSONSchemaDraft7`
- New enum type `ManagementActionType` with values `ManagementActionTypeCall`, `ManagementActionTypeRead`, `ManagementActionTypeWrite`
- New enum type `MqttDestinationQos` with values `MqttDestinationQosQos0`, `MqttDestinationQosQos1`
- New enum type `NamespaceDiscoveredManagementActionType` with values `NamespaceDiscoveredManagementActionTypeCall`, `NamespaceDiscoveredManagementActionTypeRead`, `NamespaceDiscoveredManagementActionTypeWrite`
- New enum type `SchemaType` with values `SchemaTypeMessageSchema`
- New enum type `Scope` with values `ScopeResources`
- New enum type `StreamDestinationTarget` with values `StreamDestinationTargetMqtt`, `StreamDestinationTargetStorage`
- New enum type `SystemAssignedServiceIdentityType` with values `SystemAssignedServiceIdentityTypeNone`, `SystemAssignedServiceIdentityTypeSystemAssigned`
- New function `*ClientFactory.NewNamespaceAssetsClient() *NamespaceAssetsClient`
- New function `*ClientFactory.NewNamespaceDevicesClient() *NamespaceDevicesClient`
- New function `*ClientFactory.NewNamespaceDiscoveredAssetsClient() *NamespaceDiscoveredAssetsClient`
- New function `*ClientFactory.NewNamespaceDiscoveredDevicesClient() *NamespaceDiscoveredDevicesClient`
- New function `*ClientFactory.NewNamespacesClient() *NamespacesClient`
- New function `*ClientFactory.NewSchemaRegistriesClient() *SchemaRegistriesClient`
- New function `*ClientFactory.NewSchemaVersionsClient() *SchemaVersionsClient`
- New function `*ClientFactory.NewSchemasClient() *SchemasClient`
- New function `*DatasetBrokerStateStoreDestination.GetDatasetDestination() *DatasetDestination`
- New function `*DatasetDestination.GetDatasetDestination() *DatasetDestination`
- New function `*DatasetMqttDestination.GetDatasetDestination() *DatasetDestination`
- New function `*DatasetStorageDestination.GetDatasetDestination() *DatasetDestination`
- New function `*EventDestination.GetEventDestination() *EventDestination`
- New function `*EventMqttDestination.GetEventDestination() *EventDestination`
- New function `*EventStorageDestination.GetEventDestination() *EventDestination`
- New function `NewSchemaRegistriesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SchemaRegistriesClient, error)`
- New function `*SchemaRegistriesClient.BeginCreateOrReplace(context.Context, string, string, SchemaRegistry, *SchemaRegistriesClientBeginCreateOrReplaceOptions) (*runtime.Poller[SchemaRegistriesClientCreateOrReplaceResponse], error)`
- New function `*SchemaRegistriesClient.BeginDelete(context.Context, string, string, *SchemaRegistriesClientBeginDeleteOptions) (*runtime.Poller[SchemaRegistriesClientDeleteResponse], error)`
- New function `*SchemaRegistriesClient.Get(context.Context, string, string, *SchemaRegistriesClientGetOptions) (SchemaRegistriesClientGetResponse, error)`
- New function `*SchemaRegistriesClient.NewListByResourceGroupPager(string, *SchemaRegistriesClientListByResourceGroupOptions) *runtime.Pager[SchemaRegistriesClientListByResourceGroupResponse]`
- New function `*SchemaRegistriesClient.NewListBySubscriptionPager(*SchemaRegistriesClientListBySubscriptionOptions) *runtime.Pager[SchemaRegistriesClientListBySubscriptionResponse]`
- New function `*SchemaRegistriesClient.BeginUpdate(context.Context, string, string, SchemaRegistryUpdate, *SchemaRegistriesClientBeginUpdateOptions) (*runtime.Poller[SchemaRegistriesClientUpdateResponse], error)`
- New function `NewSchemaVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SchemaVersionsClient, error)`
- New function `*SchemaVersionsClient.CreateOrReplace(context.Context, string, string, string, string, SchemaVersion, *SchemaVersionsClientCreateOrReplaceOptions) (SchemaVersionsClientCreateOrReplaceResponse, error)`
- New function `*SchemaVersionsClient.BeginDelete(context.Context, string, string, string, string, *SchemaVersionsClientBeginDeleteOptions) (*runtime.Poller[SchemaVersionsClientDeleteResponse], error)`
- New function `*SchemaVersionsClient.Get(context.Context, string, string, string, string, *SchemaVersionsClientGetOptions) (SchemaVersionsClientGetResponse, error)`
- New function `*SchemaVersionsClient.NewListBySchemaPager(string, string, string, *SchemaVersionsClientListBySchemaOptions) *runtime.Pager[SchemaVersionsClientListBySchemaResponse]`
- New function `NewSchemasClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SchemasClient, error)`
- New function `*SchemasClient.CreateOrReplace(context.Context, string, string, string, Schema, *SchemasClientCreateOrReplaceOptions) (SchemasClientCreateOrReplaceResponse, error)`
- New function `*SchemasClient.BeginDelete(context.Context, string, string, string, *SchemasClientBeginDeleteOptions) (*runtime.Poller[SchemasClientDeleteResponse], error)`
- New function `*SchemasClient.Get(context.Context, string, string, string, *SchemasClientGetOptions) (SchemasClientGetResponse, error)`
- New function `*SchemasClient.NewListBySchemaRegistryPager(string, string, *SchemasClientListBySchemaRegistryOptions) *runtime.Pager[SchemasClientListBySchemaRegistryResponse]`
- New function `*StreamDestination.GetStreamDestination() *StreamDestination`
- New function `*StreamMqttDestination.GetStreamDestination() *StreamDestination`
- New function `*StreamStorageDestination.GetStreamDestination() *StreamDestination`
- New function `NewNamespaceAssetsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NamespaceAssetsClient, error)`
- New function `*NamespaceAssetsClient.BeginCreateOrReplace(context.Context, string, string, string, NamespaceAsset, *NamespaceAssetsClientBeginCreateOrReplaceOptions) (*runtime.Poller[NamespaceAssetsClientCreateOrReplaceResponse], error)`
- New function `*NamespaceAssetsClient.BeginDelete(context.Context, string, string, string, *NamespaceAssetsClientBeginDeleteOptions) (*runtime.Poller[NamespaceAssetsClientDeleteResponse], error)`
- New function `*NamespaceAssetsClient.Get(context.Context, string, string, string, *NamespaceAssetsClientGetOptions) (NamespaceAssetsClientGetResponse, error)`
- New function `*NamespaceAssetsClient.NewListByResourceGroupPager(string, string, *NamespaceAssetsClientListByResourceGroupOptions) *runtime.Pager[NamespaceAssetsClientListByResourceGroupResponse]`
- New function `*NamespaceAssetsClient.BeginUpdate(context.Context, string, string, string, NamespaceAssetUpdate, *NamespaceAssetsClientBeginUpdateOptions) (*runtime.Poller[NamespaceAssetsClientUpdateResponse], error)`
- New function `NewNamespaceDevicesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NamespaceDevicesClient, error)`
- New function `*NamespaceDevicesClient.BeginCreateOrReplace(context.Context, string, string, string, NamespaceDevice, *NamespaceDevicesClientBeginCreateOrReplaceOptions) (*runtime.Poller[NamespaceDevicesClientCreateOrReplaceResponse], error)`
- New function `*NamespaceDevicesClient.BeginDelete(context.Context, string, string, string, *NamespaceDevicesClientBeginDeleteOptions) (*runtime.Poller[NamespaceDevicesClientDeleteResponse], error)`
- New function `*NamespaceDevicesClient.Get(context.Context, string, string, string, *NamespaceDevicesClientGetOptions) (NamespaceDevicesClientGetResponse, error)`
- New function `*NamespaceDevicesClient.NewListByResourceGroupPager(string, string, *NamespaceDevicesClientListByResourceGroupOptions) *runtime.Pager[NamespaceDevicesClientListByResourceGroupResponse]`
- New function `*NamespaceDevicesClient.BeginUpdate(context.Context, string, string, string, NamespaceDeviceUpdate, *NamespaceDevicesClientBeginUpdateOptions) (*runtime.Poller[NamespaceDevicesClientUpdateResponse], error)`
- New function `NewNamespaceDiscoveredAssetsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NamespaceDiscoveredAssetsClient, error)`
- New function `*NamespaceDiscoveredAssetsClient.BeginCreateOrReplace(context.Context, string, string, string, NamespaceDiscoveredAsset, *NamespaceDiscoveredAssetsClientBeginCreateOrReplaceOptions) (*runtime.Poller[NamespaceDiscoveredAssetsClientCreateOrReplaceResponse], error)`
- New function `*NamespaceDiscoveredAssetsClient.BeginDelete(context.Context, string, string, string, *NamespaceDiscoveredAssetsClientBeginDeleteOptions) (*runtime.Poller[NamespaceDiscoveredAssetsClientDeleteResponse], error)`
- New function `*NamespaceDiscoveredAssetsClient.Get(context.Context, string, string, string, *NamespaceDiscoveredAssetsClientGetOptions) (NamespaceDiscoveredAssetsClientGetResponse, error)`
- New function `*NamespaceDiscoveredAssetsClient.NewListByResourceGroupPager(string, string, *NamespaceDiscoveredAssetsClientListByResourceGroupOptions) *runtime.Pager[NamespaceDiscoveredAssetsClientListByResourceGroupResponse]`
- New function `*NamespaceDiscoveredAssetsClient.BeginUpdate(context.Context, string, string, string, NamespaceDiscoveredAssetUpdate, *NamespaceDiscoveredAssetsClientBeginUpdateOptions) (*runtime.Poller[NamespaceDiscoveredAssetsClientUpdateResponse], error)`
- New function `NewNamespaceDiscoveredDevicesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NamespaceDiscoveredDevicesClient, error)`
- New function `*NamespaceDiscoveredDevicesClient.BeginCreateOrReplace(context.Context, string, string, string, NamespaceDiscoveredDevice, *NamespaceDiscoveredDevicesClientBeginCreateOrReplaceOptions) (*runtime.Poller[NamespaceDiscoveredDevicesClientCreateOrReplaceResponse], error)`
- New function `*NamespaceDiscoveredDevicesClient.BeginDelete(context.Context, string, string, string, *NamespaceDiscoveredDevicesClientBeginDeleteOptions) (*runtime.Poller[NamespaceDiscoveredDevicesClientDeleteResponse], error)`
- New function `*NamespaceDiscoveredDevicesClient.Get(context.Context, string, string, string, *NamespaceDiscoveredDevicesClientGetOptions) (NamespaceDiscoveredDevicesClientGetResponse, error)`
- New function `*NamespaceDiscoveredDevicesClient.NewListByResourceGroupPager(string, string, *NamespaceDiscoveredDevicesClientListByResourceGroupOptions) *runtime.Pager[NamespaceDiscoveredDevicesClientListByResourceGroupResponse]`
- New function `*NamespaceDiscoveredDevicesClient.BeginUpdate(context.Context, string, string, string, NamespaceDiscoveredDeviceUpdate, *NamespaceDiscoveredDevicesClientBeginUpdateOptions) (*runtime.Poller[NamespaceDiscoveredDevicesClientUpdateResponse], error)`
- New function `NewNamespacesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NamespacesClient, error)`
- New function `*NamespacesClient.BeginCreateOrReplace(context.Context, string, string, Namespace, *NamespacesClientBeginCreateOrReplaceOptions) (*runtime.Poller[NamespacesClientCreateOrReplaceResponse], error)`
- New function `*NamespacesClient.BeginDelete(context.Context, string, string, *NamespacesClientBeginDeleteOptions) (*runtime.Poller[NamespacesClientDeleteResponse], error)`
- New function `*NamespacesClient.Get(context.Context, string, string, *NamespacesClientGetOptions) (NamespacesClientGetResponse, error)`
- New function `*NamespacesClient.NewListByResourceGroupPager(string, *NamespacesClientListByResourceGroupOptions) *runtime.Pager[NamespacesClientListByResourceGroupResponse]`
- New function `*NamespacesClient.NewListBySubscriptionPager(*NamespacesClientListBySubscriptionOptions) *runtime.Pager[NamespacesClientListBySubscriptionResponse]`
- New function `*NamespacesClient.BeginMigrate(context.Context, string, string, NamespaceMigrateRequest, *NamespacesClientBeginMigrateOptions) (*runtime.Poller[NamespacesClientMigrateResponse], error)`
- New function `*NamespacesClient.BeginUpdate(context.Context, string, string, NamespaceUpdate, *NamespacesClientBeginUpdateOptions) (*runtime.Poller[NamespacesClientUpdateResponse], error)`
- New struct `BrokerStateStoreDestinationConfiguration`
- New struct `DatasetBrokerStateStoreDestination`
- New struct `DatasetMqttDestination`
- New struct `DatasetStorageDestination`
- New struct `DeviceMessagingEndpoint`
- New struct `DeviceRef`
- New struct `DeviceStatus`
- New struct `DeviceStatusEndpoint`
- New struct `DeviceStatusEndpoints`
- New struct `DiscoveredInboundEndpoints`
- New struct `DiscoveredMessagingEndpoints`
- New struct `DiscoveredOutboundEndpoints`
- New struct `ErrorDetails`
- New struct `EventMqttDestination`
- New struct `EventStorageDestination`
- New struct `HostAuthentication`
- New struct `InboundEndpoints`
- New struct `ManagementAction`
- New struct `ManagementGroup`
- New struct `Messaging`
- New struct `MessagingEndpoint`
- New struct `MessagingEndpoints`
- New struct `MqttDestinationConfiguration`
- New struct `Namespace`
- New struct `NamespaceAsset`
- New struct `NamespaceAssetListResult`
- New struct `NamespaceAssetProperties`
- New struct `NamespaceAssetStatus`
- New struct `NamespaceAssetStatusDataset`
- New struct `NamespaceAssetStatusEvent`
- New struct `NamespaceAssetStatusEventGroup`
- New struct `NamespaceAssetStatusManagementAction`
- New struct `NamespaceAssetStatusManagementGroup`
- New struct `NamespaceAssetStatusStream`
- New struct `NamespaceAssetUpdate`
- New struct `NamespaceAssetUpdateProperties`
- New struct `NamespaceDataset`
- New struct `NamespaceDatasetDataPoint`
- New struct `NamespaceDevice`
- New struct `NamespaceDeviceListResult`
- New struct `NamespaceDeviceProperties`
- New struct `NamespaceDeviceUpdate`
- New struct `NamespaceDeviceUpdateProperties`
- New struct `NamespaceDiscoveredAsset`
- New struct `NamespaceDiscoveredAssetListResult`
- New struct `NamespaceDiscoveredAssetProperties`
- New struct `NamespaceDiscoveredAssetUpdate`
- New struct `NamespaceDiscoveredAssetUpdateProperties`
- New struct `NamespaceDiscoveredDataset`
- New struct `NamespaceDiscoveredDatasetDataPoint`
- New struct `NamespaceDiscoveredDevice`
- New struct `NamespaceDiscoveredDeviceListResult`
- New struct `NamespaceDiscoveredDeviceProperties`
- New struct `NamespaceDiscoveredDeviceUpdate`
- New struct `NamespaceDiscoveredDeviceUpdateProperties`
- New struct `NamespaceDiscoveredEvent`
- New struct `NamespaceDiscoveredEventGroup`
- New struct `NamespaceDiscoveredManagementAction`
- New struct `NamespaceDiscoveredManagementGroup`
- New struct `NamespaceDiscoveredStream`
- New struct `NamespaceEvent`
- New struct `NamespaceEventGroup`
- New struct `NamespaceListResult`
- New struct `NamespaceMessageSchemaReference`
- New struct `NamespaceMigrateRequest`
- New struct `NamespaceProperties`
- New struct `NamespaceStream`
- New struct `NamespaceUpdate`
- New struct `NamespaceUpdateProperties`
- New struct `OutboundEndpoints`
- New struct `Schema`
- New struct `SchemaListResult`
- New struct `SchemaProperties`
- New struct `SchemaRegistry`
- New struct `SchemaRegistryListResult`
- New struct `SchemaRegistryProperties`
- New struct `SchemaRegistryUpdate`
- New struct `SchemaRegistryUpdateProperties`
- New struct `SchemaVersion`
- New struct `SchemaVersionListResult`
- New struct `SchemaVersionProperties`
- New struct `StatusConfig`
- New struct `StatusError`
- New struct `StorageDestinationConfiguration`
- New struct `StreamMqttDestination`
- New struct `StreamStorageDestination`
- New struct `SystemAssignedServiceIdentity`
- New struct `TrustSettings`
- New struct `X509CertificateCredentials`


## 1.0.0 (2025-02-10)
### Features Added

- New field `ResourceID` in struct `OperationStatusResult`


## 0.2.0 (2024-12-11)
### Breaking Changes

- Type of `AssetProperties.Version` has been changed from `*int32` to `*int64`
- Type of `AssetStatus.Version` has been changed from `*int32` to `*int64`
- Type of `DataPoint.ObservabilityMode` has been changed from `*DataPointsObservabilityMode` to `*DataPointObservabilityMode`
- Type of `ErrorAdditionalInfo.Info` has been changed from `any` to `*ErrorAdditionalInfoInfo`
- Type of `Event.ObservabilityMode` has been changed from `*EventsObservabilityMode` to `*EventObservabilityMode`
- Type of `OperationStatusResult.PercentComplete` has been changed from `*float32` to `*float64`
- Enum `DataPointsObservabilityMode` has been removed
- Enum `EventsObservabilityMode` has been removed
- Enum `UserAuthenticationMode` has been removed
- Struct `OwnCertificate` has been removed
- Struct `TransportAuthentication` has been removed
- Struct `TransportAuthenticationUpdate` has been removed
- Struct `UserAuthentication` has been removed
- Struct `UserAuthenticationUpdate` has been removed
- Struct `UsernamePasswordCredentialsUpdate` has been removed
- Struct `X509CredentialsUpdate` has been removed
- Field `TransportAuthentication`, `UserAuthentication` of struct `AssetEndpointProfileProperties` has been removed
- Field `TransportAuthentication`, `UserAuthentication` of struct `AssetEndpointProfileUpdateProperties` has been removed
- Field `AssetEndpointProfileURI`, `AssetType`, `DataPoints`, `DefaultDataPointsConfiguration` of struct `AssetProperties` has been removed
- Field `AssetType`, `DataPoints`, `DefaultDataPointsConfiguration` of struct `AssetUpdateProperties` has been removed
- Field `CapabilityID` of struct `DataPoint` has been removed
- Field `CapabilityID` of struct `Event` has been removed
- Field `PasswordReference`, `UsernameReference` of struct `UsernamePasswordCredentials` has been removed
- Field `CertificateReference` of struct `X509Credentials` has been removed

### Features Added

- New value `ProvisioningStateDeleting` added to enum type `ProvisioningState`
- New enum type `AuthenticationMethod` with values `AuthenticationMethodAnonymous`, `AuthenticationMethodCertificate`, `AuthenticationMethodUsernamePassword`
- New enum type `DataPointObservabilityMode` with values `DataPointObservabilityModeCounter`, `DataPointObservabilityModeGauge`, `DataPointObservabilityModeHistogram`, `DataPointObservabilityModeLog`, `DataPointObservabilityModeNone`
- New enum type `EventObservabilityMode` with values `EventObservabilityModeLog`, `EventObservabilityModeNone`
- New enum type `TopicRetainType` with values `TopicRetainTypeKeep`, `TopicRetainTypeNever`
- New function `NewBillingContainersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BillingContainersClient, error)`
- New function `*BillingContainersClient.Get(context.Context, string, *BillingContainersClientGetOptions) (BillingContainersClientGetResponse, error)`
- New function `*BillingContainersClient.NewListBySubscriptionPager(*BillingContainersClientListBySubscriptionOptions) *runtime.Pager[BillingContainersClientListBySubscriptionResponse]`
- New function `*ClientFactory.NewBillingContainersClient() *BillingContainersClient`
- New struct `AssetEndpointProfileStatus`
- New struct `AssetEndpointProfileStatusError`
- New struct `AssetStatusDataset`
- New struct `AssetStatusEvent`
- New struct `Authentication`
- New struct `BillingContainer`
- New struct `BillingContainerListResult`
- New struct `BillingContainerProperties`
- New struct `Dataset`
- New struct `ErrorAdditionalInfoInfo`
- New struct `MessageSchemaReference`
- New struct `Topic`
- New field `Authentication`, `DiscoveredAssetEndpointProfileRef`, `EndpointProfileType`, `Status` in struct `AssetEndpointProfileProperties`
- New field `Authentication`, `EndpointProfileType` in struct `AssetEndpointProfileUpdateProperties`
- New field `AssetEndpointProfileRef`, `Datasets`, `DefaultDatasetsConfiguration`, `DefaultTopic`, `DiscoveredAssetRefs` in struct `AssetProperties`
- New field `Datasets`, `Events` in struct `AssetStatus`
- New field `Datasets`, `DefaultDatasetsConfiguration`, `DefaultTopic` in struct `AssetUpdateProperties`
- New field `Topic` in struct `Event`
- New field `PasswordSecretName`, `UsernameSecretName` in struct `UsernamePasswordCredentials`
- New field `CertificateSecretName` in struct `X509Credentials`


## 0.1.0 (2024-04-26)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/deviceregistry/armdeviceregistry` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).