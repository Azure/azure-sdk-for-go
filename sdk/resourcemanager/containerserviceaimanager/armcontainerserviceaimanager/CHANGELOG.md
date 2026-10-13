# Release History

## 0.2.0 (2026-09-16)
### Breaking Changes

- Function `*AIModelsClient.CalculateCost` parameter(s) have been changed from `(ctx context.Context, location string, aiModelName string, body CalculateCostRequest, options *AIModelsClientCalculateCostOptions)` to `(ctx context.Context, location string, aiModelName string, options *AIModelsClientCalculateCostOptions)`
- Struct `CalculateCostRequest` has been removed

### Features Added

- New value `ModelSourceTypeMicrosoftFoundry` added to enum type `ModelSourceType`
- New enum type `CustomAIModelProvisioningState` with values `CustomAIModelProvisioningStateCanceled`, `CustomAIModelProvisioningStateCreating`, `CustomAIModelProvisioningStateDeleting`, `CustomAIModelProvisioningStateFailed`, `CustomAIModelProvisioningStateSucceeded`, `CustomAIModelProvisioningStateUpdating`
- New function `*ClientFactory.NewCustomAIModelsClient() *CustomAIModelsClient`
- New function `NewCustomAIModelsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CustomAIModelsClient, error)`
- New function `*CustomAIModelsClient.CalculateCost(ctx context.Context, resourceGroupName string, aiManagerName string, customAIModelName string, options *CustomAIModelsClientCalculateCostOptions) (CustomAIModelsClientCalculateCostResponse, error)`
- New function `*CustomAIModelsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, aiManagerName string, customAIModelName string, resource CustomAIModel, options *CustomAIModelsClientBeginCreateOrUpdateOptions) (*runtime.Poller[CustomAIModelsClientCreateOrUpdateResponse], error)`
- New function `*CustomAIModelsClient.BeginDelete(ctx context.Context, resourceGroupName string, aiManagerName string, customAIModelName string, options *CustomAIModelsClientBeginDeleteOptions) (*runtime.Poller[CustomAIModelsClientDeleteResponse], error)`
- New function `*CustomAIModelsClient.Get(ctx context.Context, resourceGroupName string, aiManagerName string, customAIModelName string, options *CustomAIModelsClientGetOptions) (CustomAIModelsClientGetResponse, error)`
- New function `*CustomAIModelsClient.NewListPager(resourceGroupName string, aiManagerName string, options *CustomAIModelsClientListOptions) *runtime.Pager[CustomAIModelsClientListResponse]`
- New struct `BaseModelReference`
- New struct `CustomAIModel`
- New struct `CustomAIModelListResult`
- New struct `CustomAIModelProperties`
- New struct `CustomAIModelSpec`
- New struct `ManagedIdentityCredential`
- New struct `MicrosoftFoundrySource`
- New field `ClusterResourceID` in struct `AIManagerProperties`
- New field `ManagedIdentity` in struct `CredentialValue`
- New field `MicrosoftFoundry` in struct `ModelSourceProperties`


## 0.1.0 (2026-08-05)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerserviceaimanager/armcontainerserviceaimanager` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).