# Release History

## 3.0.0-beta.1 (2026-03-27)
### Breaking Changes

- Type of `SystemData.CreatedByType` has been changed from `*string` to `*CreatedByType`
- Type of `SystemData.LastModifiedByType` has been changed from `*string` to `*CreatedByType`
- Struct `ErrorModel` has been removed
- Struct `ErrorModelError` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `ServiceResourcePropertiesBase` has been removed

### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New function `*ClientFactory.NewUnsupportedVMSizesClient() *UnsupportedVMSizesClient`
- New function `NewUnsupportedVMSizesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*UnsupportedVMSizesClient, error)`
- New function `*UnsupportedVMSizesClient.Get(ctx context.Context, location string, vmSize string, options *UnsupportedVMSizesClientGetOptions) (UnsupportedVMSizesClientGetResponse, error)`
- New function `*UnsupportedVMSizesClient.NewListPager(location string, options *UnsupportedVMSizesClientListOptions) *runtime.Pager[UnsupportedVMSizesClientListResponse]`
- New struct `VMSize`
- New struct `VMSizeResource`
- New struct `VMSizesResult`
- New field `EnableHTTPGatewayExclusiveAuthMode` in struct `ClusterProperties`
- New field `EnableHTTPGatewayExclusiveAuthMode` in struct `ClusterPropertiesUpdateParameters`
- New field `HTTPGatewayTokenAuthEndpointPort` in struct `NodeTypeDescription`
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