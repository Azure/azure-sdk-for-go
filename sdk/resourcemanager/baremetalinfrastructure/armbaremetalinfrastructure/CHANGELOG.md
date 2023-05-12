# Release History

## 2.0.0 (2023-05-12)
### Breaking Changes

- Type of `ErrorResponse.Error` has been changed from `*ErrorDefinition` to `*ErrorDetail`
- Struct `ErrorDefinition` has been removed

### Features Added

- New enum type `ProvisioningState` with values `ProvisioningStateAccepted`, `ProvisioningStateCanceled`, `ProvisioningStateCreating`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateMigrating`, `ProvisioningStateSucceeded`, `ProvisioningStateUpdating`
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
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `StorageBillingProperties`
- New struct `StorageProperties`


## 1.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/baremetalinfrastructure/armbaremetalinfrastructure` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).