# Release History

## 1.2.0 (2026-06-17)
### Features Added

- New function `*ClientFactory.NewMemberCapOverridesClient() *MemberCapOverridesClient`
- New function `*ClientFactory.NewSharedLimitCapsClient() *SharedLimitCapsClient`
- New function `NewMemberCapOverridesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MemberCapOverridesClient, error)`
- New function `*MemberCapOverridesClient.CreateOrUpdate(ctx context.Context, location string, vmFamilyName string, memberSubscriptionID string, resource MemberCapOverride, options *MemberCapOverridesClientCreateOrUpdateOptions) (MemberCapOverridesClientCreateOrUpdateResponse, error)`
- New function `*MemberCapOverridesClient.Delete(ctx context.Context, location string, vmFamilyName string, memberSubscriptionID string, options *MemberCapOverridesClientDeleteOptions) (MemberCapOverridesClientDeleteResponse, error)`
- New function `*MemberCapOverridesClient.Get(ctx context.Context, location string, vmFamilyName string, memberSubscriptionID string, options *MemberCapOverridesClientGetOptions) (MemberCapOverridesClientGetResponse, error)`
- New function `*MemberCapOverridesClient.NewListByParentPager(location string, vmFamilyName string, options *MemberCapOverridesClientListByParentOptions) *runtime.Pager[MemberCapOverridesClientListByParentResponse]`
- New function `NewSharedLimitCapsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SharedLimitCapsClient, error)`
- New function `*SharedLimitCapsClient.CreateOrUpdate(ctx context.Context, location string, vmFamilyName string, resource SharedLimitCap, options *SharedLimitCapsClientCreateOrUpdateOptions) (SharedLimitCapsClientCreateOrUpdateResponse, error)`
- New function `*SharedLimitCapsClient.Delete(ctx context.Context, location string, vmFamilyName string, options *SharedLimitCapsClientDeleteOptions) (SharedLimitCapsClientDeleteResponse, error)`
- New function `*SharedLimitCapsClient.Get(ctx context.Context, location string, vmFamilyName string, options *SharedLimitCapsClientGetOptions) (SharedLimitCapsClientGetResponse, error)`
- New function `*SharedLimitCapsClient.NewListBySubscriptionLocationResourcePager(location string, options *SharedLimitCapsClientListBySubscriptionLocationResourceOptions) *runtime.Pager[SharedLimitCapsClientListBySubscriptionLocationResourceResponse]`
- New function `*SharedLimitCapsClient.SetMemberCapOverrides(ctx context.Context, location string, vmFamilyName string, body SetMemberCapOverridesRequest, options *SharedLimitCapsClientSetMemberCapOverridesOptions) (SharedLimitCapsClientSetMemberCapOverridesResponse, error)`
- New struct `MemberCap`
- New struct `MemberCapOverride`
- New struct `MemberCapOverrideListResult`
- New struct `MemberCapOverrideProperties`
- New struct `SetMemberCapOverridesRequest`
- New struct `SetMemberCapOverridesResult`
- New struct `SharedLimitCap`
- New struct `SharedLimitCapListResult`
- New struct `SharedLimitCapProperties`


## 1.1.0 (2026-05-26)
### Features Added

- New function `*ClientFactory.NewVMFamiliesClient() *VMFamiliesClient`
- New function `*FeaturesClient.BeginDisable(ctx context.Context, location string, featureName string, options *FeaturesClientBeginDisableOptions) (*runtime.Poller[FeaturesClientDisableResponse], error)`
- New function `NewVMFamiliesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VMFamiliesClient, error)`
- New function `*VMFamiliesClient.Get(ctx context.Context, location string, vmFamilyName string, options *VMFamiliesClientGetOptions) (VMFamiliesClientGetResponse, error)`
- New function `*VMFamiliesClient.NewListBySubscriptionLocationResourcePager(location string, options *VMFamiliesClientListBySubscriptionLocationResourceOptions) *runtime.Pager[VMFamiliesClientListBySubscriptionLocationResourceResponse]`
- New struct `FeatureEnableRequest`
- New struct `VMFamily`
- New struct `VMFamilyListResult`
- New struct `VMFamilyProperties`
- New field `Body` in struct `FeaturesClientBeginEnableOptions`


## 1.0.0 (2026-04-06)
### Features Added

- New enum type `FeatureState` with values `FeatureStateDisabled`, `FeatureStateEnabled`
- New function `*ClientFactory.NewFeaturesClient() *FeaturesClient`
- New function `NewFeaturesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*FeaturesClient, error)`
- New function `*FeaturesClient.BeginEnable(ctx context.Context, location string, featureName string, options *FeaturesClientBeginEnableOptions) (*runtime.Poller[FeaturesClientEnableResponse], error)`
- New function `*FeaturesClient.Get(ctx context.Context, location string, featureName string, options *FeaturesClientGetOptions) (FeaturesClientGetResponse, error)`
- New function `*FeaturesClient.NewListBySubscriptionLocationResourcePager(location string, options *FeaturesClientListBySubscriptionLocationResourceOptions) *runtime.Pager[FeaturesClientListBySubscriptionLocationResourceResponse]`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `Feature`
- New struct `FeatureListResult`
- New struct `FeatureProperties`
- New struct `OperationStatusResult`


## 0.1.0 (2025-11-14)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/computelimit/armcomputelimit` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).