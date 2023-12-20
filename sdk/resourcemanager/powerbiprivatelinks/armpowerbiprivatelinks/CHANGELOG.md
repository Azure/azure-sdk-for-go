# Release History

## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 2.0.0 (2023-04-03)
### Breaking Changes

- Function `NewPowerBIResourcesClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*PowerBIResourcesClient.Create` parameter(s) have been changed from `(context.Context, TenantResource, *PowerBIResourcesClientCreateOptions)` to `(context.Context, string, string, TenantResource, *PowerBIResourcesClientCreateOptions)`
- Function `*PowerBIResourcesClient.Delete` parameter(s) have been changed from `(context.Context, *PowerBIResourcesClientDeleteOptions)` to `(context.Context, string, string, *PowerBIResourcesClientDeleteOptions)`
- Function `*PowerBIResourcesClient.ListByResourceName` parameter(s) have been changed from `(context.Context, *PowerBIResourcesClientListByResourceNameOptions)` to `(context.Context, string, string, *PowerBIResourcesClientListByResourceNameOptions)`
- Function `*PowerBIResourcesClient.Update` parameter(s) have been changed from `(context.Context, TenantResource, *PowerBIResourcesClientUpdateOptions)` to `(context.Context, string, string, TenantResource, *PowerBIResourcesClientUpdateOptions)`
- Function `NewPrivateEndpointConnectionsClient` parameter(s) have been changed from `(string, string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*PrivateEndpointConnectionsClient.BeginDelete` parameter(s) have been changed from `(context.Context, *PrivateEndpointConnectionsClientBeginDeleteOptions)` to `(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions)`
- Function `*PrivateEndpointConnectionsClient.Create` parameter(s) have been changed from `(context.Context, PrivateEndpointConnection, *PrivateEndpointConnectionsClientCreateOptions)` to `(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientCreateOptions)`
- Function `*PrivateEndpointConnectionsClient.Get` parameter(s) have been changed from `(context.Context, *PrivateEndpointConnectionsClientGetOptions)` to `(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions)`
- Function `NewPrivateLinkResourcesClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*PrivateLinkResourcesClient.Get` parameter(s) have been changed from `(context.Context, string, *PrivateLinkResourcesClientGetOptions)` to `(context.Context, string, string, string, *PrivateLinkResourcesClientGetOptions)`
- Function `*PrivateLinkResourcesClient.NewListByResourcePager` parameter(s) have been changed from `(*PrivateLinkResourcesClientListByResourceOptions)` to `(string, string, *PrivateLinkResourcesClientListByResourceOptions)`
- Function `NewPrivateLinkServiceResourceOperationResultsClient` parameter(s) have been changed from `(string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*PrivateLinkServiceResourceOperationResultsClient.BeginGet` parameter(s) have been changed from `(context.Context, *PrivateLinkServiceResourceOperationResultsClientBeginGetOptions)` to `(context.Context, string, *PrivateLinkServiceResourceOperationResultsClientBeginGetOptions)`
- Function `NewPrivateLinkServicesClient` parameter(s) have been changed from `(string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*PrivateLinkServicesClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, *PrivateLinkServicesClientListByResourceGroupOptions)` to `(context.Context, string, *PrivateLinkServicesClientListByResourceGroupOptions)`

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/powerbiprivatelinks/armpowerbiprivatelinks` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).