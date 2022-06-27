# Release History

## 3.0.0 (2022-06-24)
### Breaking Changes

- Function `*OperationResultsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, string, string, *OperationResultsClientGetOptions)` to `(context.Context, string, string, *OperationResultsClientGetOptions)`
- Function `*OperationStatusesClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, string, string, *OperationStatusesClientGetOptions)` to `(context.Context, string, string, *OperationStatusesClientGetOptions)`
- Function `*MediaServiceOperationStatusesClient.Get` has been removed
- Function `NewMediaServiceOperationResultsClient` has been removed
- Function `*MediaServiceOperationResultsClient.Get` has been removed
- Function `NewMediaServiceOperationStatusesClient` has been removed
- Struct `MediaServiceOperationResultsClient` has been removed
- Struct `MediaServiceOperationResultsClientGetOptions` has been removed
- Struct `MediaServiceOperationResultsClientGetResponse` has been removed
- Struct `MediaServiceOperationStatusesClient` has been removed
- Struct `MediaServiceOperationStatusesClientGetOptions` has been removed
- Struct `MediaServiceOperationStatusesClientGetResponse` has been removed
- Field `AssetTrackOperationStatus` of struct `OperationStatusesClientGetResponse` has been removed
- Field `AssetTrack` of struct `OperationResultsClientGetResponse` has been removed

### Features Added

- New function `NewAssetTrackOperationStatusesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AssetTrackOperationStatusesClient, error)`
- New function `NewAssetTrackOperationResultsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AssetTrackOperationResultsClient, error)`
- New function `*AssetTrackOperationResultsClient.Get(context.Context, string, string, string, string, string, *AssetTrackOperationResultsClientGetOptions) (AssetTrackOperationResultsClientGetResponse, error)`
- New function `*AssetTrackOperationStatusesClient.Get(context.Context, string, string, string, string, string, *AssetTrackOperationStatusesClientGetOptions) (AssetTrackOperationStatusesClientGetResponse, error)`
- New struct `AssetTrackOperationResultsClient`
- New struct `AssetTrackOperationResultsClientGetOptions`
- New struct `AssetTrackOperationResultsClientGetResponse`
- New struct `AssetTrackOperationStatusesClient`
- New struct `AssetTrackOperationStatusesClientGetOptions`
- New struct `AssetTrackOperationStatusesClientGetResponse`
- New anonymous field `MediaServiceOperationStatus` in struct `OperationStatusesClientGetResponse`
- New anonymous field `MediaService` in struct `OperationResultsClientGetResponse`


## 2.0.0 (2022-06-22)
### Breaking Changes

- Function `*Client.Update` has been removed
- Function `*Client.CreateOrUpdate` has been removed
- Struct `ClientCreateOrUpdateOptions` has been removed
- Struct `ClientUpdateOptions` has been removed

### Features Added

- New function `*MediaServiceOperationResultsClient.Get(context.Context, string, string, *MediaServiceOperationResultsClientGetOptions) (MediaServiceOperationResultsClientGetResponse, error)`
- New function `*Client.BeginCreateOrUpdate(context.Context, string, string, MediaService, *ClientBeginCreateOrUpdateOptions) (*runtime.Poller[ClientCreateOrUpdateResponse], error)`
- New function `*MediaServiceOperationStatusesClient.Get(context.Context, string, string, *MediaServiceOperationStatusesClientGetOptions) (MediaServiceOperationStatusesClientGetResponse, error)`
- New function `*Client.BeginUpdate(context.Context, string, string, MediaServiceUpdate, *ClientBeginUpdateOptions) (*runtime.Poller[ClientUpdateResponse], error)`
- New function `NewMediaServiceOperationStatusesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MediaServiceOperationStatusesClient, error)`
- New function `NewMediaServiceOperationResultsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MediaServiceOperationResultsClient, error)`
- New struct `ClientBeginCreateOrUpdateOptions`
- New struct `ClientBeginUpdateOptions`
- New struct `MediaServiceOperationResultsClient`
- New struct `MediaServiceOperationResultsClientGetOptions`
- New struct `MediaServiceOperationResultsClientGetResponse`
- New struct `MediaServiceOperationStatus`
- New struct `MediaServiceOperationStatusesClient`
- New struct `MediaServiceOperationStatusesClientGetOptions`
- New struct `MediaServiceOperationStatusesClientGetResponse`
- New field `PrivateEndpointConnections` in struct `MediaServiceProperties`
- New field `ProvisioningState` in struct `MediaServiceProperties`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mediaservices/armmediaservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).