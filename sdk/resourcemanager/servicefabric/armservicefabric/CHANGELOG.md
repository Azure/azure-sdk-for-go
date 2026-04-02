# Release History

## 3.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Type of `ApplicationResource.Identity` has been changed from `*ManagedIdentity` to `*ManagedServiceIdentity`
- Type of `OperationListResult.Value` has been changed from `[]*OperationResult` to `[]*Operation`
- Type of `ServiceResourceProperties.ServicePlacementPolicies` has been changed from `[]ServicePlacementPolicyDescriptionClassification` to `[]*ServicePlacementPolicyDescription`
- Type of `ServiceResourceUpdateProperties.ServicePlacementPolicies` has been changed from `[]ServicePlacementPolicyDescriptionClassification` to `[]*ServicePlacementPolicyDescription`
- Type of `StatefulServiceProperties.ServicePlacementPolicies` has been changed from `[]ServicePlacementPolicyDescriptionClassification` to `[]*ServicePlacementPolicyDescription`
- Type of `StatefulServiceUpdateProperties.ServicePlacementPolicies` has been changed from `[]ServicePlacementPolicyDescriptionClassification` to `[]*ServicePlacementPolicyDescription`
- Type of `StatelessServiceProperties.ServicePlacementPolicies` has been changed from `[]ServicePlacementPolicyDescriptionClassification` to `[]*ServicePlacementPolicyDescription`
- Type of `StatelessServiceUpdateProperties.ServicePlacementPolicies` has been changed from `[]ServicePlacementPolicyDescriptionClassification` to `[]*ServicePlacementPolicyDescription`
- Type of `SystemData.CreatedByType` has been changed from `*string` to `*CreatedByType`
- Type of `SystemData.LastModifiedByType` has been changed from `*string` to `*CreatedByType`
- Enum `ManagedIdentityType` has been removed
- Function `*ServicePlacementPolicyDescription.GetServicePlacementPolicyDescription` has been removed
- Struct `AvailableOperationDisplay` has been removed
- Struct `ErrorModel` has been removed
- Struct `ErrorModelError` has been removed
- Struct `ManagedIdentity` has been removed
- Struct `OperationResult` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `ServiceResourcePropertiesBase` has been removed
- Field `Etag`, `Location` of struct `ApplicationResourceUpdate` has been removed
- Field `ApplicationResource` of struct `ApplicationsClientUpdateResponse` has been removed
- Field `Cluster` of struct `ClustersClientUpdateResponse` has been removed
- Field `Etag` of struct `ServiceResource` has been removed
- Field `Etag` of struct `ServiceResourceUpdate` has been removed
- Field `ServiceResource` of struct `ServicesClientUpdateResponse` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New function `*ClientFactory.NewUnsupportedVMSizesClient() *UnsupportedVMSizesClient`
- New function `NewUnsupportedVMSizesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*UnsupportedVMSizesClient, error)`
- New function `*UnsupportedVMSizesClient.Get(ctx context.Context, location string, vmSize string, options *UnsupportedVMSizesClientGetOptions) (UnsupportedVMSizesClientGetResponse, error)`
- New function `*UnsupportedVMSizesClient.NewListPager(location string, options *UnsupportedVMSizesClientListOptions) *runtime.Pager[UnsupportedVMSizesClientListResponse]`
- New struct `ManagedServiceIdentity`
- New struct `Operation`
- New struct `OperationDisplay`
- New struct `VMSize`
- New struct `VMSizeResource`
- New struct `VMSizesResult`
- New field `SystemData` in struct `ClusterCodeVersionsResult`
- New field `EnableHTTPGatewayExclusiveAuthMode` in struct `ClusterProperties`
- New field `EnableHTTPGatewayExclusiveAuthMode` in struct `ClusterPropertiesUpdateParameters`
- New field `HTTPGatewayTokenAuthEndpointPort` in struct `NodeTypeDescription`
- New field `ETag` in struct `ServiceResource`
- New field `ETag` in struct `ServiceResourceUpdate`
- New field `MinInstanceCount`, `MinInstancePercentage` in struct `StatelessServiceProperties`


## 2.0.0 (2023-12-22)
### Breaking Changes

- Operation `*ApplicationTypeVersionsClient.List` has supported pagination, use `*ApplicationTypeVersionsClient.NewListPager` instead.
- Operation `*ApplicationTypesClient.List` has supported pagination, use `*ApplicationTypesClient.NewListPager` instead.
- Operation `*ApplicationsClient.List` has supported pagination, use `*ApplicationsClient.NewListPager` instead.
- Operation `*ClustersClient.List` has supported pagination, use `*ClustersClient.NewListPager` instead.
- Operation `*ClustersClient.ListByResourceGroup` has supported pagination, use `*ClustersClient.NewListByResourceGroupPager` instead.
- Operation `*ServicesClient.List` has supported pagination, use `*ServicesClient.NewListPager` instead.


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicefabric/armservicefabric` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).