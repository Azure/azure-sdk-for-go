# Release History

## 3.1.0 (2022-09-02)
### Features Added

- New const `SecurityLevelSL150`
- New const `AsyncOperationStatusInProgress`
- New const `AsyncOperationStatusSucceeded`
- New const `SecurityLevelSL3000`
- New const `AsyncOperationStatusFailed`
- New const `SecurityLevelUnknown`
- New const `SecurityLevelSL2000`
- New type alias `AsyncOperationStatus`
- New type alias `SecurityLevel`
- New function `PossibleSecurityLevelValues() []SecurityLevel`
- New function `*LiveEventsClient.OperationLocation(context.Context, string, string, string, string, *LiveEventsClientOperationLocationOptions) (LiveEventsClientOperationLocationResponse, error)`
- New function `*LiveOutputsClient.AsyncOperation(context.Context, string, string, string, *LiveOutputsClientAsyncOperationOptions) (LiveOutputsClientAsyncOperationResponse, error)`
- New function `PossibleAsyncOperationStatusValues() []AsyncOperationStatus`
- New function `*LiveEventsClient.AsyncOperation(context.Context, string, string, string, *LiveEventsClientAsyncOperationOptions) (LiveEventsClientAsyncOperationResponse, error)`
- New function `*StreamingEndpointsClient.AsyncOperation(context.Context, string, string, string, *StreamingEndpointsClientAsyncOperationOptions) (StreamingEndpointsClientAsyncOperationResponse, error)`
- New function `*LiveOutputsClient.OperationLocation(context.Context, string, string, string, string, string, *LiveOutputsClientOperationLocationOptions) (LiveOutputsClientOperationLocationResponse, error)`
- New function `*StreamingEndpointsClient.OperationLocation(context.Context, string, string, string, string, *StreamingEndpointsClientOperationLocationOptions) (StreamingEndpointsClientOperationLocationResponse, error)`
- New struct `AsyncOperationResult`
- New struct `ClearKeyEncryptionConfiguration`
- New struct `DashSettings`
- New struct `LiveEventsClientAsyncOperationOptions`
- New struct `LiveEventsClientAsyncOperationResponse`
- New struct `LiveEventsClientOperationLocationOptions`
- New struct `LiveEventsClientOperationLocationResponse`
- New struct `LiveOutputsClientAsyncOperationOptions`
- New struct `LiveOutputsClientAsyncOperationResponse`
- New struct `LiveOutputsClientOperationLocationOptions`
- New struct `LiveOutputsClientOperationLocationResponse`
- New struct `StreamingEndpointsClientAsyncOperationOptions`
- New struct `StreamingEndpointsClientAsyncOperationResponse`
- New struct `StreamingEndpointsClientOperationLocationOptions`
- New struct `StreamingEndpointsClientOperationLocationResponse`
- New field `SecurityLevel` in struct `ContentKeyPolicyPlayReadyLicense`
- New field `ClearKeyEncryptionConfiguration` in struct `CommonEncryptionCenc`
- New field `RewindWindowLength` in struct `LiveOutputProperties`
- New field `ClearKeyEncryptionConfiguration` in struct `CommonEncryptionCbcs`
- New field `DashSettings` in struct `AudioTrack`
- New field `DisplayName` in struct `AudioTrack`
- New field `FileName` in struct `AudioTrack`
- New field `HlsSettings` in struct `AudioTrack`
- New field `LanguageCode` in struct `AudioTrack`
- New field `Mpeg4TrackID` in struct `AudioTrack`
- New field `BitRate` in struct `AudioTrack`


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