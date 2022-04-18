# Release History

## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*PrivateEndpointConnectionsClient.List` has been removed

### Features Added

- New function `*PrivateEndpointConnectionsClient.NewListPager(string, string, *PrivateEndpointConnectionsClientListOptions) *runtime.Pager[PrivateEndpointConnectionsClientListResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `NewPrivateEndpointConnectionsClient` return value(s) have been changed from `(*PrivateEndpointConnectionsClient)` to `(*PrivateEndpointConnectionsClient, error)`
- Function `*PrivateEndpointConnectionsClient.List` parameter(s) have been changed from `(context.Context, string, string, *PrivateEndpointConnectionsClientListOptions)` to `(string, string, *PrivateEndpointConnectionsClientListOptions)`
- Function `*PrivateEndpointConnectionsClient.List` return value(s) have been changed from `(PrivateEndpointConnectionsClientListResponse, error)` to `(*runtime.Pager[PrivateEndpointConnectionsClientListResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `NewProvidersClient` return value(s) have been changed from `(*ProvidersClient)` to `(*ProvidersClient, error)`
- Function `CreatedByType.ToPtr` has been removed
- Function `PrivateEndpointConnectionProvisioningState.ToPtr` has been removed
- Function `AttestationServiceStatus.ToPtr` has been removed
- Function `PrivateEndpointServiceConnectionStatus.ToPtr` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `PrivateEndpointConnectionsClientCreateResult` has been removed
- Struct `PrivateEndpointConnectionsClientGetResult` has been removed
- Struct `PrivateEndpointConnectionsClientListResult` has been removed
- Struct `ProvidersClientCreateResult` has been removed
- Struct `ProvidersClientGetDefaultByLocationResult` has been removed
- Struct `ProvidersClientGetResult` has been removed
- Struct `ProvidersClientListByResourceGroupResult` has been removed
- Struct `ProvidersClientListDefaultResult` has been removed
- Struct `ProvidersClientListResult` has been removed
- Struct `ProvidersClientUpdateResult` has been removed
- Field `RawResponse` of struct `PrivateEndpointConnectionsClientDeleteResponse` has been removed
- Field `ProvidersClientGetDefaultByLocationResult` of struct `ProvidersClientGetDefaultByLocationResponse` has been removed
- Field `RawResponse` of struct `ProvidersClientGetDefaultByLocationResponse` has been removed
- Field `ProvidersClientListDefaultResult` of struct `ProvidersClientListDefaultResponse` has been removed
- Field `RawResponse` of struct `ProvidersClientListDefaultResponse` has been removed
- Field `PrivateEndpointConnectionsClientCreateResult` of struct `PrivateEndpointConnectionsClientCreateResponse` has been removed
- Field `RawResponse` of struct `PrivateEndpointConnectionsClientCreateResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `PrivateEndpointConnectionsClientListResult` of struct `PrivateEndpointConnectionsClientListResponse` has been removed
- Field `RawResponse` of struct `PrivateEndpointConnectionsClientListResponse` has been removed
- Field `ProvidersClientUpdateResult` of struct `ProvidersClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ProvidersClientUpdateResponse` has been removed
- Field `ProvidersClientCreateResult` of struct `ProvidersClientCreateResponse` has been removed
- Field `RawResponse` of struct `ProvidersClientCreateResponse` has been removed
- Field `RawResponse` of struct `ProvidersClientDeleteResponse` has been removed
- Field `ProvidersClientListByResourceGroupResult` of struct `ProvidersClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `ProvidersClientListByResourceGroupResponse` has been removed
- Field `ProvidersClientGetResult` of struct `ProvidersClientGetResponse` has been removed
- Field `RawResponse` of struct `ProvidersClientGetResponse` has been removed
- Field `ProvidersClientListResult` of struct `ProvidersClientListResponse` has been removed
- Field `RawResponse` of struct `ProvidersClientListResponse` has been removed
- Field `PrivateEndpointConnectionsClientGetResult` of struct `PrivateEndpointConnectionsClientGetResponse` has been removed
- Field `RawResponse` of struct `PrivateEndpointConnectionsClientGetResponse` has been removed

### Features Added

- New anonymous field `OperationList` in struct `OperationsClientListResponse`
- New anonymous field `PrivateEndpointConnection` in struct `PrivateEndpointConnectionsClientGetResponse`
- New anonymous field `PrivateEndpointConnection` in struct `PrivateEndpointConnectionsClientCreateResponse`
- New anonymous field `Provider` in struct `ProvidersClientGetDefaultByLocationResponse`
- New anonymous field `Provider` in struct `ProvidersClientCreateResponse`
- New anonymous field `ProviderListResult` in struct `ProvidersClientListByResourceGroupResponse`
- New anonymous field `ProviderListResult` in struct `ProvidersClientListResponse`
- New anonymous field `ProviderListResult` in struct `ProvidersClientListDefaultResponse`
- New anonymous field `PrivateEndpointConnectionListResult` in struct `PrivateEndpointConnectionsClientListResponse`
- New anonymous field `Provider` in struct `ProvidersClientGetResponse`
- New anonymous field `Provider` in struct `ProvidersClientUpdateResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*PrivateEndpointConnectionsClient.Create` parameter(s) have been changed from `(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsCreateOptions)` to `(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientCreateOptions)`
- Function `*PrivateEndpointConnectionsClient.Create` return value(s) have been changed from `(PrivateEndpointConnectionsCreateResponse, error)` to `(PrivateEndpointConnectionsClientCreateResponse, error)`
- Function `*PrivateEndpointConnectionsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, string, *PrivateEndpointConnectionsDeleteOptions)` to `(context.Context, string, string, string, *PrivateEndpointConnectionsClientDeleteOptions)`
- Function `*PrivateEndpointConnectionsClient.Delete` return value(s) have been changed from `(PrivateEndpointConnectionsDeleteResponse, error)` to `(PrivateEndpointConnectionsClientDeleteResponse, error)`
- Function `*PrivateEndpointConnectionsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *PrivateEndpointConnectionsGetOptions)` to `(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions)`
- Function `*PrivateEndpointConnectionsClient.Get` return value(s) have been changed from `(PrivateEndpointConnectionsGetResponse, error)` to `(PrivateEndpointConnectionsClientGetResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Function `*PrivateEndpointConnectionsClient.List` parameter(s) have been changed from `(context.Context, string, string, *PrivateEndpointConnectionsListOptions)` to `(context.Context, string, string, *PrivateEndpointConnectionsClientListOptions)`
- Function `*PrivateEndpointConnectionsClient.List` return value(s) have been changed from `(PrivateEndpointConnectionsListResponse, error)` to `(PrivateEndpointConnectionsClientListResponse, error)`
- Function `*AttestationProvidersClient.Update` has been removed
- Function `AttestationProvider.MarshalJSON` has been removed
- Function `AttestationProviderListResult.MarshalJSON` has been removed
- Function `NewAttestationProvidersClient` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*AttestationProvidersClient.List` has been removed
- Function `*AttestationProvidersClient.Delete` has been removed
- Function `*AttestationProvidersClient.Create` has been removed
- Function `AttestationServiceCreationParams.MarshalJSON` has been removed
- Function `PrivateEndpointConnection.MarshalJSON` has been removed
- Function `*AttestationProvidersClient.ListDefault` has been removed
- Function `*AttestationProvidersClient.GetDefaultByLocation` has been removed
- Function `CloudError.Error` has been removed
- Function `AttestationServicePatchParams.MarshalJSON` has been removed
- Function `*AttestationProvidersClient.Get` has been removed
- Function `*AttestationProvidersClient.ListByResourceGroup` has been removed
- Struct `AttestationProvider` has been removed
- Struct `AttestationProviderListResult` has been removed
- Struct `AttestationProvidersClient` has been removed
- Struct `AttestationProvidersCreateOptions` has been removed
- Struct `AttestationProvidersCreateResponse` has been removed
- Struct `AttestationProvidersCreateResult` has been removed
- Struct `AttestationProvidersDeleteOptions` has been removed
- Struct `AttestationProvidersDeleteResponse` has been removed
- Struct `AttestationProvidersGetDefaultByLocationOptions` has been removed
- Struct `AttestationProvidersGetDefaultByLocationResponse` has been removed
- Struct `AttestationProvidersGetDefaultByLocationResult` has been removed
- Struct `AttestationProvidersGetOptions` has been removed
- Struct `AttestationProvidersGetResponse` has been removed
- Struct `AttestationProvidersGetResult` has been removed
- Struct `AttestationProvidersListByResourceGroupOptions` has been removed
- Struct `AttestationProvidersListByResourceGroupResponse` has been removed
- Struct `AttestationProvidersListByResourceGroupResult` has been removed
- Struct `AttestationProvidersListDefaultOptions` has been removed
- Struct `AttestationProvidersListDefaultResponse` has been removed
- Struct `AttestationProvidersListDefaultResult` has been removed
- Struct `AttestationProvidersListOptions` has been removed
- Struct `AttestationProvidersListResponse` has been removed
- Struct `AttestationProvidersListResult` has been removed
- Struct `AttestationProvidersUpdateOptions` has been removed
- Struct `AttestationProvidersUpdateResponse` has been removed
- Struct `AttestationProvidersUpdateResult` has been removed
- Struct `AttestationServiceCreationParams` has been removed
- Struct `AttestationServiceCreationSpecificParams` has been removed
- Struct `AttestationServicePatchParams` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `PrivateEndpointConnectionsCreateOptions` has been removed
- Struct `PrivateEndpointConnectionsCreateResponse` has been removed
- Struct `PrivateEndpointConnectionsCreateResult` has been removed
- Struct `PrivateEndpointConnectionsDeleteOptions` has been removed
- Struct `PrivateEndpointConnectionsDeleteResponse` has been removed
- Struct `PrivateEndpointConnectionsGetOptions` has been removed
- Struct `PrivateEndpointConnectionsGetResponse` has been removed
- Struct `PrivateEndpointConnectionsGetResult` has been removed
- Struct `PrivateEndpointConnectionsListOptions` has been removed
- Struct `PrivateEndpointConnectionsListResponse` has been removed
- Struct `PrivateEndpointConnectionsListResult` has been removed
- Field `InnerError` of struct `CloudError` has been removed
- Field `Resource` of struct `PrivateEndpointConnection` has been removed
- Field `Resource` of struct `TrackedResource` has been removed

### Features Added

- New function `*ProvidersClient.GetDefaultByLocation(context.Context, string, *ProvidersClientGetDefaultByLocationOptions) (ProvidersClientGetDefaultByLocationResponse, error)`
- New function `*ProvidersClient.List(context.Context, *ProvidersClientListOptions) (ProvidersClientListResponse, error)`
- New function `ServiceCreationParams.MarshalJSON() ([]byte, error)`
- New function `*ProvidersClient.Create(context.Context, string, string, ServiceCreationParams, *ProvidersClientCreateOptions) (ProvidersClientCreateResponse, error)`
- New function `*ProvidersClient.ListByResourceGroup(context.Context, string, *ProvidersClientListByResourceGroupOptions) (ProvidersClientListByResourceGroupResponse, error)`
- New function `*ProvidersClient.ListDefault(context.Context, *ProvidersClientListDefaultOptions) (ProvidersClientListDefaultResponse, error)`
- New function `*ProvidersClient.Update(context.Context, string, string, ServicePatchParams, *ProvidersClientUpdateOptions) (ProvidersClientUpdateResponse, error)`
- New function `ServicePatchParams.MarshalJSON() ([]byte, error)`
- New function `*ProvidersClient.Delete(context.Context, string, string, *ProvidersClientDeleteOptions) (ProvidersClientDeleteResponse, error)`
- New function `Provider.MarshalJSON() ([]byte, error)`
- New function `ProviderListResult.MarshalJSON() ([]byte, error)`
- New function `*ProvidersClient.Get(context.Context, string, string, *ProvidersClientGetOptions) (ProvidersClientGetResponse, error)`
- New function `NewProvidersClient(string, azcore.TokenCredential, *arm.ClientOptions) *ProvidersClient`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `PrivateEndpointConnectionsClientCreateOptions`
- New struct `PrivateEndpointConnectionsClientCreateResponse`
- New struct `PrivateEndpointConnectionsClientCreateResult`
- New struct `PrivateEndpointConnectionsClientDeleteOptions`
- New struct `PrivateEndpointConnectionsClientDeleteResponse`
- New struct `PrivateEndpointConnectionsClientGetOptions`
- New struct `PrivateEndpointConnectionsClientGetResponse`
- New struct `PrivateEndpointConnectionsClientGetResult`
- New struct `PrivateEndpointConnectionsClientListOptions`
- New struct `PrivateEndpointConnectionsClientListResponse`
- New struct `PrivateEndpointConnectionsClientListResult`
- New struct `Provider`
- New struct `ProviderListResult`
- New struct `ProvidersClient`
- New struct `ProvidersClientCreateOptions`
- New struct `ProvidersClientCreateResponse`
- New struct `ProvidersClientCreateResult`
- New struct `ProvidersClientDeleteOptions`
- New struct `ProvidersClientDeleteResponse`
- New struct `ProvidersClientGetDefaultByLocationOptions`
- New struct `ProvidersClientGetDefaultByLocationResponse`
- New struct `ProvidersClientGetDefaultByLocationResult`
- New struct `ProvidersClientGetOptions`
- New struct `ProvidersClientGetResponse`
- New struct `ProvidersClientGetResult`
- New struct `ProvidersClientListByResourceGroupOptions`
- New struct `ProvidersClientListByResourceGroupResponse`
- New struct `ProvidersClientListByResourceGroupResult`
- New struct `ProvidersClientListDefaultOptions`
- New struct `ProvidersClientListDefaultResponse`
- New struct `ProvidersClientListDefaultResult`
- New struct `ProvidersClientListOptions`
- New struct `ProvidersClientListResponse`
- New struct `ProvidersClientListResult`
- New struct `ProvidersClientUpdateOptions`
- New struct `ProvidersClientUpdateResponse`
- New struct `ProvidersClientUpdateResult`
- New struct `ServiceCreationParams`
- New struct `ServiceCreationSpecificParams`
- New struct `ServicePatchParams`
- New field `ID` in struct `PrivateEndpointConnection`
- New field `Name` in struct `PrivateEndpointConnection`
- New field `Type` in struct `PrivateEndpointConnection`
- New field `Error` in struct `CloudError`
- New field `ID` in struct `TrackedResource`
- New field `Name` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`


## 0.1.0 (2021-11-30)

- Initial preview release.
