# Release History

## 2.0.0-beta.1 (2024-03-22)
### Breaking Changes

- Function `*Client.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, NotificationHubCreateOrUpdateParameters, *ClientCreateOrUpdateOptions)` to `(context.Context, string, string, string, NotificationHubResource, *ClientCreateOrUpdateOptions)`
- Function `*Client.CreateOrUpdateAuthorizationRule` parameter(s) have been changed from `(context.Context, string, string, string, string, SharedAccessAuthorizationRuleCreateOrUpdateParameters, *ClientCreateOrUpdateAuthorizationRuleOptions)` to `(context.Context, string, string, string, string, SharedAccessAuthorizationRuleResource, *ClientCreateOrUpdateAuthorizationRuleOptions)`
- Function `*Client.RegenerateKeys` parameter(s) have been changed from `(context.Context, string, string, string, string, PolicykeyResource, *ClientRegenerateKeysOptions)` to `(context.Context, string, string, string, string, PolicyKeyResource, *ClientRegenerateKeysOptions)`
- Function `*NamespacesClient.CreateOrUpdateAuthorizationRule` parameter(s) have been changed from `(context.Context, string, string, string, SharedAccessAuthorizationRuleCreateOrUpdateParameters, *NamespacesClientCreateOrUpdateAuthorizationRuleOptions)` to `(context.Context, string, string, string, SharedAccessAuthorizationRuleResource, *NamespacesClientCreateOrUpdateAuthorizationRuleOptions)`
- Function `*NamespacesClient.RegenerateKeys` parameter(s) have been changed from `(context.Context, string, string, string, PolicykeyResource, *NamespacesClientRegenerateKeysOptions)` to `(context.Context, string, string, string, PolicyKeyResource, *NamespacesClientRegenerateKeysOptions)`
- Type of `DebugSendResult.Failure` has been changed from `*float32` to `*int64`
- Type of `DebugSendResult.Results` has been changed from `any` to `[]*RegistrationResult`
- Type of `DebugSendResult.Success` has been changed from `*float32` to `*int64`
- Type of `NamespaceProperties.ProvisioningState` has been changed from `*string` to `*OperationProvisioningState`
- Type of `NamespaceProperties.Status` has been changed from `*string` to `*NamespaceStatus`
- Type of `PnsCredentialsResource.Properties` has been changed from `*PnsCredentialsProperties` to `*PnsCredentials`
- Type of `SharedAccessAuthorizationRuleProperties.CreatedTime` has been changed from `*string` to `*time.Time`
- Type of `SharedAccessAuthorizationRuleProperties.ModifiedTime` has been changed from `*string` to `*time.Time`
- Function `*Client.Patch` has been removed
- Function `*NamespacesClient.Patch` has been removed
- Operation `*NamespacesClient.CreateOrUpdate` has been changed to LRO, use `*NamespacesClient.BeginCreateOrUpdate` instead.
- Operation `*NamespacesClient.BeginDelete` has been changed to non-LRO, use `*NamespacesClient.Delete` instead.
- Struct `NamespaceCreateOrUpdateParameters` has been removed
- Struct `NotificationHubCreateOrUpdateParameters` has been removed
- Struct `PnsCredentialsProperties` has been removed
- Struct `PolicykeyResource` has been removed
- Struct `SharedAccessAuthorizationRuleCreateOrUpdateParameters` has been removed
- Struct `SubResource` has been removed
- Field `Parameters` of struct `ClientDebugSendOptions` has been removed
- Field `SKU` of struct `DebugSendResponse` has been removed
- Field `Code`, `Message` of struct `ErrorResponse` has been removed
- Field `ID`, `Location`, `Name`, `Type` of struct `NotificationHubPatchParameters` has been removed
- Field `SKU` of struct `PnsCredentialsResource` has been removed
- Field `Location`, `SKU`, `Tags` of struct `Resource` has been removed
- Field `SKU` of struct `SharedAccessAuthorizationRuleResource` has been removed

### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `NamespaceStatus` with values `NamespaceStatusCreated`, `NamespaceStatusCreating`, `NamespaceStatusDeleting`, `NamespaceStatusSuspended`
- New enum type `OperationProvisioningState` with values `OperationProvisioningStateCanceled`, `OperationProvisioningStateDisabled`, `OperationProvisioningStateFailed`, `OperationProvisioningStateInProgress`, `OperationProvisioningStatePending`, `OperationProvisioningStateSucceeded`, `OperationProvisioningStateUnknown`
- New enum type `PolicyKeyType` with values `PolicyKeyTypePrimaryKey`, `PolicyKeyTypeSecondaryKey`
- New enum type `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleted`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateDeletingByProxy`, `PrivateEndpointConnectionProvisioningStateSucceeded`, `PrivateEndpointConnectionProvisioningStateUnknown`, `PrivateEndpointConnectionProvisioningStateUpdating`, `PrivateEndpointConnectionProvisioningStateUpdatingByProxy`
- New enum type `PrivateLinkConnectionStatus` with values `PrivateLinkConnectionStatusApproved`, `PrivateLinkConnectionStatusDisconnected`, `PrivateLinkConnectionStatusPending`, `PrivateLinkConnectionStatusRejected`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`
- New enum type `ReplicationRegion` with values `ReplicationRegionAustraliaEast`, `ReplicationRegionBrazilSouth`, `ReplicationRegionDefault`, `ReplicationRegionNone`, `ReplicationRegionNorthEurope`, `ReplicationRegionSouthAfricaNorth`, `ReplicationRegionSouthEastAsia`, `ReplicationRegionWestUs2`
- New enum type `ZoneRedundancyPreference` with values `ZoneRedundancyPreferenceDisabled`, `ZoneRedundancyPreferenceEnabled`
- New function `*Client.Update(context.Context, string, string, string, NotificationHubPatchParameters, *ClientUpdateOptions) (ClientUpdateResponse, error)`
- New function `*ClientFactory.NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.GetGroupID(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetGroupIDOptions) (PrivateEndpointConnectionsClientGetGroupIDResponse, error)`
- New function `*PrivateEndpointConnectionsClient.NewListGroupIDsPager(string, string, *PrivateEndpointConnectionsClientListGroupIDsOptions) *runtime.Pager[PrivateEndpointConnectionsClientListGroupIDsResponse]`
- New function `*PrivateEndpointConnectionsClient.NewListPager(string, string, *PrivateEndpointConnectionsClientListOptions) *runtime.Pager[PrivateEndpointConnectionsClientListResponse]`
- New function `*PrivateEndpointConnectionsClient.BeginUpdate(context.Context, string, string, string, PrivateEndpointConnectionResource, *PrivateEndpointConnectionsClientBeginUpdateOptions) (*runtime.Poller[PrivateEndpointConnectionsClientUpdateResponse], error)`
- New function `*NamespacesClient.GetPnsCredentials(context.Context, string, string, *NamespacesClientGetPnsCredentialsOptions) (NamespacesClientGetPnsCredentialsResponse, error)`
- New function `*NamespacesClient.Update(context.Context, string, string, NamespacePatchParameters, *NamespacesClientUpdateOptions) (NamespacesClientUpdateResponse, error)`
- New struct `Availability`
- New struct `BrowserCredential`
- New struct `BrowserCredentialProperties`
- New struct `ConnectionDetails`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `FcmV1Credential`
- New struct `FcmV1CredentialProperties`
- New struct `GroupConnectivityInformation`
- New struct `IPRule`
- New struct `LogSpecification`
- New struct `MetricSpecification`
- New struct `NetworkACLs`
- New struct `OperationProperties`
- New struct `PnsCredentials`
- New struct `PolicyKeyResource`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateEndpointConnectionResource`
- New struct `PrivateEndpointConnectionResourceListResult`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnection`
- New struct `ProxyResource`
- New struct `PublicInternetAuthorizationRule`
- New struct `RegistrationResult`
- New struct `RemotePrivateEndpointConnection`
- New struct `RemotePrivateLinkServiceConnectionState`
- New struct `ServiceSpecification`
- New struct `SystemData`
- New struct `TrackedResource`
- New struct `XiaomiCredential`
- New struct `XiaomiCredentialProperties`
- New field `SystemData` in struct `CheckAvailabilityResult`
- New field `SkipToken`, `Top` in struct `ClientListOptions`
- New field `SystemData` in struct `DebugSendResponse`
- New field `Error` in struct `ErrorResponse`
- New field `Properties` in struct `NamespacePatchParameters`
- New field `NetworkACLs`, `PnsCredentials`, `PrivateEndpointConnections`, `PublicNetworkAccess`, `ReplicationRegion`, `ZoneRedundancy` in struct `NamespaceProperties`
- New field `SystemData` in struct `NamespaceResource`
- New field `SkipToken`, `Top` in struct `NamespacesClientListAllOptions`
- New field `SkipToken`, `Top` in struct `NamespacesClientListOptions`
- New field `BrowserCredential`, `DailyMaxActiveDevices`, `FcmV1Credential`, `XiaomiCredential` in struct `NotificationHubProperties`
- New field `SystemData` in struct `NotificationHubResource`
- New field `IsDataAction`, `Properties` in struct `Operation`
- New field `Description` in struct `OperationDisplay`
- New field `SystemData` in struct `PnsCredentialsResource`
- New field `SystemData` in struct `Resource`
- New field `SystemData` in struct `SharedAccessAuthorizationRuleResource`
- New field `CertificateKey`, `WnsCertificate` in struct `WnsCredentialProperties`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/notificationhubs/armnotificationhubs` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).