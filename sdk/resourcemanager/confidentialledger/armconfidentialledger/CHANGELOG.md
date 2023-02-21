# Release History

## 2.0.0-beta.1 (2023-02-21)
### Breaking Changes

- Struct `ResourceLocation` has been removed

### Features Added

- New type alias `LanguageRuntime` with values `LanguageRuntimeCPP`, `LanguageRuntimeJS`
- New type alias `RunningState` with values `RunningStateActive`, `RunningStatePaused`, `RunningStatePausing`, `RunningStateResuming`, `RunningStateUnknown`
- New function `NewManagedCCFClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedCCFClient, error)`
- New function `*ManagedCCFClient.BeginCreate(context.Context, string, string, ManagedCCF, *ManagedCCFClientBeginCreateOptions) (*runtime.Poller[ManagedCCFClientCreateResponse], error)`
- New function `*ManagedCCFClient.BeginDelete(context.Context, string, string, *ManagedCCFClientBeginDeleteOptions) (*runtime.Poller[ManagedCCFClientDeleteResponse], error)`
- New function `*ManagedCCFClient.Get(context.Context, string, string, *ManagedCCFClientGetOptions) (ManagedCCFClientGetResponse, error)`
- New function `*ManagedCCFClient.NewListByResourceGroupPager(string, *ManagedCCFClientListByResourceGroupOptions) *runtime.Pager[ManagedCCFClientListByResourceGroupResponse]`
- New function `*ManagedCCFClient.NewListBySubscriptionPager(*ManagedCCFClientListBySubscriptionOptions) *runtime.Pager[ManagedCCFClientListBySubscriptionResponse]`
- New function `*ManagedCCFClient.BeginUpdate(context.Context, string, string, ManagedCCF, *ManagedCCFClientBeginUpdateOptions) (*runtime.Poller[ManagedCCFClientUpdateResponse], error)`
- New struct `CertificateTags`
- New struct `DeploymentType`
- New struct `ManagedCCF`
- New struct `ManagedCCFClient`
- New struct `ManagedCCFList`
- New struct `ManagedCCFProperties`
- New struct `MemberIdentityCertificate`
- New struct `TrackedResource`
- New field `RunningState` in struct `ConfidentialLedger`


## 1.0.0 (2022-05-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/confidentialledger/armconfidentialledger` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).