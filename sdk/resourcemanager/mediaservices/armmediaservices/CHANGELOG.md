# Release History

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