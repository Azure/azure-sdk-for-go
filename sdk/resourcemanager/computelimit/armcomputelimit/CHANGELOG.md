# Release History

## 1.0.0 (2026-04-17)
### Features Added

- New enum type `FeatureState` with values `FeatureStateDisabled`, `FeatureStateEnabled`
- New function `*ClientFactory.NewFeaturesClient() *FeaturesClient`
- New function `*ClientFactory.NewVMFamiliesClient() *VMFamiliesClient`
- New function `NewFeaturesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*FeaturesClient, error)`
- New function `*FeaturesClient.BeginDisable(ctx context.Context, location string, featureName string, options *FeaturesClientBeginDisableOptions) (*runtime.Poller[FeaturesClientDisableResponse], error)`
- New function `*FeaturesClient.BeginEnable(ctx context.Context, location string, featureName string, options *FeaturesClientBeginEnableOptions) (*runtime.Poller[FeaturesClientEnableResponse], error)`
- New function `*FeaturesClient.Get(ctx context.Context, location string, featureName string, options *FeaturesClientGetOptions) (FeaturesClientGetResponse, error)`
- New function `*FeaturesClient.NewListBySubscriptionLocationResourcePager(location string, options *FeaturesClientListBySubscriptionLocationResourceOptions) *runtime.Pager[FeaturesClientListBySubscriptionLocationResourceResponse]`
- New function `NewVMFamiliesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VMFamiliesClient, error)`
- New function `*VMFamiliesClient.Get(ctx context.Context, location string, vmFamilyName string, options *VMFamiliesClientGetOptions) (VMFamiliesClientGetResponse, error)`
- New function `*VMFamiliesClient.NewListBySubscriptionLocationResourcePager(location string, options *VMFamiliesClientListBySubscriptionLocationResourceOptions) *runtime.Pager[VMFamiliesClientListBySubscriptionLocationResourceResponse]`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `Feature`
- New struct `FeatureListResult`
- New struct `FeatureProperties`
- New struct `OperationStatusResult`
- New struct `VMFamily`
- New struct `VMFamilyListResult`
- New struct `VMFamilyProperties`


## 0.1.0 (2025-11-14)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/computelimit/armcomputelimit` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).