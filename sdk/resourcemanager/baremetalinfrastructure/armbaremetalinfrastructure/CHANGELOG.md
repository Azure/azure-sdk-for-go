# Release History

## 2.0.0-beta.2 (2026-02-25)
### Breaking Changes

- Function `*AzureBareMetalStorageInstancesClient.Update` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, azureBareMetalStorageInstanceName string, tagsParameter Tags, options *AzureBareMetalStorageInstancesClientUpdateOptions)` to `(ctx context.Context, resourceGroupName string, azureBareMetalStorageInstanceName string, azureBareMetalStorageInstanceBodyParameter AzureBareMetalStorageInstanceBody, options *AzureBareMetalStorageInstancesClientUpdateOptions)`

### Features Added

- New enum type `ResourceIdentityType` with values `ResourceIdentityTypeNone`, `ResourceIdentityTypeSystemAssigned`
- New function `*AzureBareMetalInstancesClient.Create(ctx context.Context, resourceGroupName string, azureBareMetalInstanceName string, requestBodyParameters AzureBareMetalInstance, options *AzureBareMetalInstancesClientCreateOptions) (AzureBareMetalInstancesClientCreateResponse, error)`
- New function `*AzureBareMetalInstancesClient.Delete(ctx context.Context, resourceGroupName string, azureBareMetalInstanceName string, options *AzureBareMetalInstancesClientDeleteOptions) (AzureBareMetalInstancesClientDeleteResponse, error)`
- New struct `AzureBareMetalStorageInstanceBody`
- New struct `AzureBareMetalStorageInstanceIdentity`
- New field `Identity` in struct `AzureBareMetalStorageInstance`


## 2.0.0-beta.1 (2023-12-08)
### Breaking Changes

- Type of `NetworkProfile.NetworkInterfaces` has been changed from `[]*IPAddress` to `[]*NetworkInterface`
- Type of `Operation.Display` has been changed from `*Display` to `*OperationDisplay`
- Struct `Display` has been removed
- Struct `IPAddress` has been removed
- Struct `OperationList` has been removed
- Struct `Result` has been removed
- Field `OperationList` of struct `OperationsClientListResponse` has been removed

### Features Added

- New value `AzureBareMetalHardwareTypeNamesEnumSDFLEX` added to enum type `AzureBareMetalHardwareTypeNamesEnum`
- New value `AzureBareMetalInstanceSizeNamesEnumS448Se` added to enum type `AzureBareMetalInstanceSizeNamesEnum`
- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `AsyncOperationStatus` with values `AsyncOperationStatusExecuting`, `AsyncOperationStatusFailed`, `AsyncOperationStatusRequesting`, `AsyncOperationStatusSucceeded`
- New enum type `AzureBareMetalInstanceForcePowerState` with values `AzureBareMetalInstanceForcePowerStateActive`, `AzureBareMetalInstanceForcePowerStateInactive`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `ProvisioningState` with values `ProvisioningStateAccepted`, `ProvisioningStateCanceled`, `ProvisioningStateCreating`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateMigrating`, `ProvisioningStateSucceeded`, `ProvisioningStateUpdating`
- New function `*AzureBareMetalInstancesClient.BeginRestart(context.Context, string, string, *AzureBareMetalInstancesClientBeginRestartOptions) (*runtime.Poller[AzureBareMetalInstancesClientRestartResponse], error)`
- New function `*AzureBareMetalInstancesClient.BeginShutdown(context.Context, string, string, *AzureBareMetalInstancesClientBeginShutdownOptions) (*runtime.Poller[AzureBareMetalInstancesClientShutdownResponse], error)`
- New function `*AzureBareMetalInstancesClient.BeginStart(context.Context, string, string, *AzureBareMetalInstancesClientBeginStartOptions) (*runtime.Poller[AzureBareMetalInstancesClientStartResponse], error)`
- New function `NewAzureBareMetalStorageInstancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AzureBareMetalStorageInstancesClient, error)`
- New function `*AzureBareMetalStorageInstancesClient.Create(context.Context, string, string, AzureBareMetalStorageInstance, *AzureBareMetalStorageInstancesClientCreateOptions) (AzureBareMetalStorageInstancesClientCreateResponse, error)`
- New function `*AzureBareMetalStorageInstancesClient.Delete(context.Context, string, string, *AzureBareMetalStorageInstancesClientDeleteOptions) (AzureBareMetalStorageInstancesClientDeleteResponse, error)`
- New function `*AzureBareMetalStorageInstancesClient.Get(context.Context, string, string, *AzureBareMetalStorageInstancesClientGetOptions) (AzureBareMetalStorageInstancesClientGetResponse, error)`
- New function `*AzureBareMetalStorageInstancesClient.NewListByResourceGroupPager(string, *AzureBareMetalStorageInstancesClientListByResourceGroupOptions) *runtime.Pager[AzureBareMetalStorageInstancesClientListByResourceGroupResponse]`
- New function `*AzureBareMetalStorageInstancesClient.NewListBySubscriptionPager(*AzureBareMetalStorageInstancesClientListBySubscriptionOptions) *runtime.Pager[AzureBareMetalStorageInstancesClientListBySubscriptionResponse]`
- New function `*AzureBareMetalStorageInstancesClient.Update(context.Context, string, string, Tags, *AzureBareMetalStorageInstancesClientUpdateOptions) (AzureBareMetalStorageInstancesClientUpdateResponse, error)`
- New function `*ClientFactory.NewAzureBareMetalStorageInstancesClient() *AzureBareMetalStorageInstancesClient`
- New struct `AzureBareMetalStorageInstance`
- New struct `AzureBareMetalStorageInstanceProperties`
- New struct `AzureBareMetalStorageInstancesListResult`
- New struct `ForceState`
- New struct `NetworkInterface`
- New struct `OperationDisplay`
- New struct `OperationListResult`
- New struct `OperationStatus`
- New struct `OperationStatusError`
- New struct `StorageBillingProperties`
- New struct `StorageProperties`
- New field `ActionType`, `Origin` in struct `Operation`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/baremetalinfrastructure/armbaremetalinfrastructure` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).