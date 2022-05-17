# Release History

## 1.0.0 (2022-05-17)
### Breaking Changes

- Function `*PrivateEndpointConnectionsClient.BeginCreate` return value(s) have been changed from `(*armruntime.Poller[PrivateEndpointConnectionsClientCreateResponse], error)` to `(*runtime.Poller[PrivateEndpointConnectionsClientCreateResponse], error)`
- Function `*PrivateEndpointConnectionsClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)` to `(*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- Function `*PrivateLinkForAzureAdClient.BeginCreate` return value(s) have been changed from `(*armruntime.Poller[PrivateLinkForAzureAdClientCreateResponse], error)` to `(*runtime.Poller[PrivateLinkForAzureAdClientCreateResponse], error)`
- Function `ErrorDefinition.MarshalJSON` has been removed
- Function `PrivateLinkResourceProperties.MarshalJSON` has been removed
- Function `PrivateLinkPolicyListResult.MarshalJSON` has been removed
- Function `PrivateEndpointConnectionListResult.MarshalJSON` has been removed
- Function `PrivateLinkResourceListResult.MarshalJSON` has been removed


## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*PrivateEndpointConnectionsClient.ListByPolicyName` has been removed
- Function `*PrivateLinkForAzureAdClient.List` has been removed
- Function `*PrivateLinkForAzureAdClient.ListBySubscription` has been removed
- Function `*PrivateLinkResourcesClient.ListByPrivateLinkPolicy` has been removed

### Features Added

- New function `*PrivateEndpointConnectionsClient.NewListByPolicyNamePager(string, string, *PrivateEndpointConnectionsClientListByPolicyNameOptions) *runtime.Pager[PrivateEndpointConnectionsClientListByPolicyNameResponse]`
- New function `*PrivateLinkForAzureAdClient.NewListPager(string, *PrivateLinkForAzureAdClientListOptions) *runtime.Pager[PrivateLinkForAzureAdClientListResponse]`
- New function `*PrivateLinkResourcesClient.NewListByPrivateLinkPolicyPager(string, string, *PrivateLinkResourcesClientListByPrivateLinkPolicyOptions) *runtime.Pager[PrivateLinkResourcesClientListByPrivateLinkPolicyResponse]`
- New function `*PrivateLinkForAzureAdClient.NewListBySubscriptionPager(*PrivateLinkForAzureAdClientListBySubscriptionOptions) *runtime.Pager[PrivateLinkForAzureAdClientListBySubscriptionResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Const `CategoryAuditLogs` has been removed
- Const `CategorySignInLogs` has been removed
- Const `CategoryTypeLogs` has been removed
- Function `NewDiagnosticSettingsCategoryClient` has been removed
- Function `*DiagnosticSettingsCategoryClient.List` has been removed
- Function `Category.ToPtr` has been removed
- Function `DiagnosticSettingsCategoryResourceCollection.MarshalJSON` has been removed
- Function `OperationsDiscoveryCollection.MarshalJSON` has been removed
- Function `PossibleCategoryValues` has been removed
- Function `*OperationsClient.List` has been removed
- Function `CategoryType.ToPtr` has been removed
- Function `NewDiagnosticSettingsClient` has been removed
- Function `*DiagnosticSettingsClient.List` has been removed
- Function `NewOperationsClient` has been removed
- Function `PossibleCategoryTypeValues` has been removed
- Function `DiagnosticSettingsResourceCollection.MarshalJSON` has been removed
- Function `*DiagnosticSettingsClient.Delete` has been removed
- Function `*DiagnosticSettingsClient.CreateOrUpdate` has been removed
- Function `DiagnosticSettings.MarshalJSON` has been removed
- Function `*DiagnosticSettingsClient.Get` has been removed
- Struct `DiagnosticSettings` has been removed
- Struct `DiagnosticSettingsCategory` has been removed
- Struct `DiagnosticSettingsCategoryClient` has been removed
- Struct `DiagnosticSettingsCategoryClientListOptions` has been removed
- Struct `DiagnosticSettingsCategoryClientListResponse` has been removed
- Struct `DiagnosticSettingsCategoryClientListResult` has been removed
- Struct `DiagnosticSettingsCategoryResource` has been removed
- Struct `DiagnosticSettingsCategoryResourceCollection` has been removed
- Struct `DiagnosticSettingsClient` has been removed
- Struct `DiagnosticSettingsClientCreateOrUpdateOptions` has been removed
- Struct `DiagnosticSettingsClientCreateOrUpdateResponse` has been removed
- Struct `DiagnosticSettingsClientCreateOrUpdateResult` has been removed
- Struct `DiagnosticSettingsClientDeleteOptions` has been removed
- Struct `DiagnosticSettingsClientDeleteResponse` has been removed
- Struct `DiagnosticSettingsClientGetOptions` has been removed
- Struct `DiagnosticSettingsClientGetResponse` has been removed
- Struct `DiagnosticSettingsClientGetResult` has been removed
- Struct `DiagnosticSettingsClientListOptions` has been removed
- Struct `DiagnosticSettingsClientListResponse` has been removed
- Struct `DiagnosticSettingsClientListResult` has been removed
- Struct `DiagnosticSettingsResource` has been removed
- Struct `DiagnosticSettingsResourceCollection` has been removed
- Struct `Display` has been removed
- Struct `LogSettings` has been removed
- Struct `OperationsClient` has been removed
- Struct `OperationsClientListOptions` has been removed
- Struct `OperationsClientListResponse` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `OperationsDiscovery` has been removed
- Struct `OperationsDiscoveryCollection` has been removed
- Struct `ProxyOnlyResource` has been removed
- Struct `RetentionPolicy` has been removed

### Features Added

- New const `PrivateEndpointServiceConnectionStatusPending`
- New const `PrivateEndpointServiceConnectionStatusRejected`
- New const `PrivateEndpointConnectionProvisioningStateFailed`
- New const `PrivateEndpointConnectionProvisioningStateProvisioning`
- New const `PrivateEndpointConnectionProvisioningStateSucceeded`
- New const `PrivateEndpointServiceConnectionStatusApproved`
- New const `PrivateEndpointServiceConnectionStatusDisconnected`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `PrivateLinkResourceProperties.MarshalJSON() ([]byte, error)`
- New function `*PrivateLinkForAzureAdClient.Update(context.Context, string, string, *PrivateLinkForAzureAdClientUpdateOptions) (PrivateLinkForAzureAdClientUpdateResponse, error)`
- New function `*PrivateLinkForAzureAdClient.Get(context.Context, string, string, *PrivateLinkForAzureAdClientGetOptions) (PrivateLinkForAzureAdClientGetResponse, error)`
- New function `PrivateLinkPolicyListResult.MarshalJSON() ([]byte, error)`
- New function `*PrivateLinkForAzureAdClient.ListBySubscription(*PrivateLinkForAzureAdClientListBySubscriptionOptions) *runtime.Pager[PrivateLinkForAzureAdClientListBySubscriptionResponse]`
- New function `*PrivateEndpointConnectionsClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions) (*armruntime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `NewPrivateLinkForAzureAdClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateLinkForAzureAdClient, error)`
- New function `*PrivateEndpointConnectionsClient.BeginCreate(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientBeginCreateOptions) (*armruntime.Poller[PrivateEndpointConnectionsClientCreateResponse], error)`
- New function `*PrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.ListByPolicyName(string, string, *PrivateEndpointConnectionsClientListByPolicyNameOptions) *runtime.Pager[PrivateEndpointConnectionsClientListByPolicyNameResponse]`
- New function `PrivateLinkPolicy.MarshalJSON() ([]byte, error)`
- New function `PrivateLinkPolicyUpdateParameter.MarshalJSON() ([]byte, error)`
- New function `*PrivateLinkForAzureAdClient.Delete(context.Context, string, string, *PrivateLinkForAzureAdClientDeleteOptions) (PrivateLinkForAzureAdClientDeleteResponse, error)`
- New function `PossiblePrivateEndpointServiceConnectionStatusValues() []PrivateEndpointServiceConnectionStatus`
- New function `PrivateEndpointConnectionListResult.MarshalJSON() ([]byte, error)`
- New function `*PrivateLinkForAzureAdClient.BeginCreate(context.Context, string, string, PrivateLinkPolicy, *PrivateLinkForAzureAdClientBeginCreateOptions) (*armruntime.Poller[PrivateLinkForAzureAdClientCreateResponse], error)`
- New function `PrivateLinkResourceListResult.MarshalJSON() ([]byte, error)`
- New function `NewPrivateLinkResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateLinkResourcesClient, error)`
- New function `*PrivateLinkResourcesClient.Get(context.Context, string, string, string, *PrivateLinkResourcesClientGetOptions) (PrivateLinkResourcesClientGetResponse, error)`
- New function `*PrivateLinkForAzureAdClient.List(string, *PrivateLinkForAzureAdClientListOptions) *runtime.Pager[PrivateLinkForAzureAdClientListResponse]`
- New function `PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState`
- New function `*PrivateLinkResourcesClient.ListByPrivateLinkPolicy(string, string, *PrivateLinkResourcesClientListByPrivateLinkPolicyOptions) *runtime.Pager[PrivateLinkResourcesClientListByPrivateLinkPolicyResponse]`
- New function `TagsResource.MarshalJSON() ([]byte, error)`
- New struct `ARMProxyResource`
- New struct `AzureResourceBase`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateEndpointConnectionsClient`
- New struct `PrivateEndpointConnectionsClientBeginCreateOptions`
- New struct `PrivateEndpointConnectionsClientBeginDeleteOptions`
- New struct `PrivateEndpointConnectionsClientCreateResponse`
- New struct `PrivateEndpointConnectionsClientDeleteResponse`
- New struct `PrivateEndpointConnectionsClientGetOptions`
- New struct `PrivateEndpointConnectionsClientGetResponse`
- New struct `PrivateEndpointConnectionsClientListByPolicyNameOptions`
- New struct `PrivateEndpointConnectionsClientListByPolicyNameResponse`
- New struct `PrivateLinkForAzureAdClient`
- New struct `PrivateLinkForAzureAdClientBeginCreateOptions`
- New struct `PrivateLinkForAzureAdClientCreateResponse`
- New struct `PrivateLinkForAzureAdClientDeleteOptions`
- New struct `PrivateLinkForAzureAdClientDeleteResponse`
- New struct `PrivateLinkForAzureAdClientGetOptions`
- New struct `PrivateLinkForAzureAdClientGetResponse`
- New struct `PrivateLinkForAzureAdClientListBySubscriptionOptions`
- New struct `PrivateLinkForAzureAdClientListBySubscriptionResponse`
- New struct `PrivateLinkForAzureAdClientListOptions`
- New struct `PrivateLinkForAzureAdClientListResponse`
- New struct `PrivateLinkForAzureAdClientUpdateOptions`
- New struct `PrivateLinkForAzureAdClientUpdateResponse`
- New struct `PrivateLinkPolicy`
- New struct `PrivateLinkPolicyListResult`
- New struct `PrivateLinkPolicyUpdateParameter`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkResourcesClient`
- New struct `PrivateLinkResourcesClientGetOptions`
- New struct `PrivateLinkResourcesClientGetResponse`
- New struct `PrivateLinkResourcesClientListByPrivateLinkPolicyOptions`
- New struct `PrivateLinkResourcesClientListByPrivateLinkPolicyResponse`
- New struct `PrivateLinkServiceConnectionState`
- New struct `ProxyResource`
- New struct `Resource`
- New struct `TagsResource`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*DiagnosticSettingsClient.Get` parameter(s) have been changed from `(context.Context, string, *DiagnosticSettingsGetOptions)` to `(context.Context, string, *DiagnosticSettingsClientGetOptions)`
- Function `*DiagnosticSettingsClient.Get` return value(s) have been changed from `(DiagnosticSettingsGetResponse, error)` to `(DiagnosticSettingsClientGetResponse, error)`
- Function `*DiagnosticSettingsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, DiagnosticSettingsResource, *DiagnosticSettingsCreateOrUpdateOptions)` to `(context.Context, string, DiagnosticSettingsResource, *DiagnosticSettingsClientCreateOrUpdateOptions)`
- Function `*DiagnosticSettingsClient.CreateOrUpdate` return value(s) have been changed from `(DiagnosticSettingsCreateOrUpdateResponse, error)` to `(DiagnosticSettingsClientCreateOrUpdateResponse, error)`
- Function `*DiagnosticSettingsCategoryClient.List` parameter(s) have been changed from `(context.Context, *DiagnosticSettingsCategoryListOptions)` to `(context.Context, *DiagnosticSettingsCategoryClientListOptions)`
- Function `*DiagnosticSettingsCategoryClient.List` return value(s) have been changed from `(DiagnosticSettingsCategoryListResponse, error)` to `(DiagnosticSettingsCategoryClientListResponse, error)`
- Function `*DiagnosticSettingsClient.Delete` parameter(s) have been changed from `(context.Context, string, *DiagnosticSettingsDeleteOptions)` to `(context.Context, string, *DiagnosticSettingsClientDeleteOptions)`
- Function `*DiagnosticSettingsClient.Delete` return value(s) have been changed from `(DiagnosticSettingsDeleteResponse, error)` to `(DiagnosticSettingsClientDeleteResponse, error)`
- Function `*DiagnosticSettingsClient.List` parameter(s) have been changed from `(context.Context, *DiagnosticSettingsListOptions)` to `(context.Context, *DiagnosticSettingsClientListOptions)`
- Function `*DiagnosticSettingsClient.List` return value(s) have been changed from `(DiagnosticSettingsListResponse, error)` to `(DiagnosticSettingsClientListResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Function `ErrorResponse.Error` has been removed
- Struct `DiagnosticSettingsCategoryListOptions` has been removed
- Struct `DiagnosticSettingsCategoryListResponse` has been removed
- Struct `DiagnosticSettingsCategoryListResult` has been removed
- Struct `DiagnosticSettingsCreateOrUpdateOptions` has been removed
- Struct `DiagnosticSettingsCreateOrUpdateResponse` has been removed
- Struct `DiagnosticSettingsCreateOrUpdateResult` has been removed
- Struct `DiagnosticSettingsDeleteOptions` has been removed
- Struct `DiagnosticSettingsDeleteResponse` has been removed
- Struct `DiagnosticSettingsGetOptions` has been removed
- Struct `DiagnosticSettingsGetResponse` has been removed
- Struct `DiagnosticSettingsGetResult` has been removed
- Struct `DiagnosticSettingsListOptions` has been removed
- Struct `DiagnosticSettingsListResponse` has been removed
- Struct `DiagnosticSettingsListResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `ProxyOnlyResource` of struct `DiagnosticSettingsCategoryResource` has been removed
- Field `ProxyOnlyResource` of struct `DiagnosticSettingsResource` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New struct `DiagnosticSettingsCategoryClientListOptions`
- New struct `DiagnosticSettingsCategoryClientListResponse`
- New struct `DiagnosticSettingsCategoryClientListResult`
- New struct `DiagnosticSettingsClientCreateOrUpdateOptions`
- New struct `DiagnosticSettingsClientCreateOrUpdateResponse`
- New struct `DiagnosticSettingsClientCreateOrUpdateResult`
- New struct `DiagnosticSettingsClientDeleteOptions`
- New struct `DiagnosticSettingsClientDeleteResponse`
- New struct `DiagnosticSettingsClientGetOptions`
- New struct `DiagnosticSettingsClientGetResponse`
- New struct `DiagnosticSettingsClientGetResult`
- New struct `DiagnosticSettingsClientListOptions`
- New struct `DiagnosticSettingsClientListResponse`
- New struct `DiagnosticSettingsClientListResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `Error` in struct `ErrorResponse`
- New field `ID` in struct `DiagnosticSettingsCategoryResource`
- New field `Name` in struct `DiagnosticSettingsCategoryResource`
- New field `Type` in struct `DiagnosticSettingsCategoryResource`
- New field `Name` in struct `DiagnosticSettingsResource`
- New field `Type` in struct `DiagnosticSettingsResource`
- New field `ID` in struct `DiagnosticSettingsResource`


## 0.1.0 (2021-11-30)

- Initial preview release.
