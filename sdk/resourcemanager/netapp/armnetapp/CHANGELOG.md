# Release History

## 2.0.0 (2022-06-10)
### Breaking Changes

- Type of `VolumeProperties.EncryptionKeySource` has been changed from `*string` to `*EncryptionKeySource`
- Function `VolumeGroupDetails.MarshalJSON` has been removed
- Field `Tags` of struct `VolumeGroupDetails` has been removed
- Field `Tags` of struct `VolumeGroup` has been removed

### Features Added

- New const `ProvisioningStateFailed`
- New const `ProvisioningStateSucceeded`
- New const `TypeDefaultGroupQuota`
- New const `ProvisioningStateDeleting`
- New const `ProvisioningStatePatching`
- New const `ProvisioningStateCreating`
- New const `TypeDefaultUserQuota`
- New const `TypeIndividualGroupQuota`
- New const `ProvisioningStateMoving`
- New const `ProvisioningStateAccepted`
- New const `EncryptionKeySourceMicrosoftNetApp`
- New const `TypeIndividualUserQuota`
- New function `NewVolumeQuotaRulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VolumeQuotaRulesClient, error)`
- New function `*VolumeQuotaRulesClient.BeginUpdate(context.Context, string, string, string, string, string, VolumeQuotaRulePatch, *VolumeQuotaRulesClientBeginUpdateOptions) (*runtime.Poller[VolumeQuotaRulesClientUpdateResponse], error)`
- New function `*VolumesClient.BeginRevertRelocation(context.Context, string, string, string, string, *VolumesClientBeginRevertRelocationOptions) (*runtime.Poller[VolumesClientRevertRelocationResponse], error)`
- New function `*VolumesClient.BeginFinalizeRelocation(context.Context, string, string, string, string, *VolumesClientBeginFinalizeRelocationOptions) (*runtime.Poller[VolumesClientFinalizeRelocationResponse], error)`
- New function `VolumeQuotaRulePatch.MarshalJSON() ([]byte, error)`
- New function `PossibleEncryptionKeySourceValues() []EncryptionKeySource`
- New function `*VolumesClient.BeginResetCifsPassword(context.Context, string, string, string, string, *VolumesClientBeginResetCifsPasswordOptions) (*runtime.Poller[VolumesClientResetCifsPasswordResponse], error)`
- New function `TrackedResource.MarshalJSON() ([]byte, error)`
- New function `PossibleProvisioningStateValues() []ProvisioningState`
- New function `PossibleTypeValues() []Type`
- New function `*VolumesClient.BeginRelocate(context.Context, string, string, string, string, *VolumesClientBeginRelocateOptions) (*runtime.Poller[VolumesClientRelocateResponse], error)`
- New function `*VolumeQuotaRulesClient.NewListByVolumePager(string, string, string, string, *VolumeQuotaRulesClientListByVolumeOptions) *runtime.Pager[VolumeQuotaRulesClientListByVolumeResponse]`
- New function `*VolumeQuotaRulesClient.Get(context.Context, string, string, string, string, string, *VolumeQuotaRulesClientGetOptions) (VolumeQuotaRulesClientGetResponse, error)`
- New function `*VolumeQuotaRulesClient.BeginCreate(context.Context, string, string, string, string, string, VolumeQuotaRule, *VolumeQuotaRulesClientBeginCreateOptions) (*runtime.Poller[VolumeQuotaRulesClientCreateResponse], error)`
- New function `VolumeQuotaRule.MarshalJSON() ([]byte, error)`
- New function `*VolumesClient.NewListReplicationsPager(string, string, string, string, *VolumesClientListReplicationsOptions) *runtime.Pager[VolumesClientListReplicationsResponse]`
- New function `*VolumeQuotaRulesClient.BeginDelete(context.Context, string, string, string, string, string, *VolumeQuotaRulesClientBeginDeleteOptions) (*runtime.Poller[VolumeQuotaRulesClientDeleteResponse], error)`
- New struct `ListReplications`
- New struct `Replication`
- New struct `TrackedResource`
- New struct `VolumeQuotaRule`
- New struct `VolumeQuotaRulePatch`
- New struct `VolumeQuotaRulesClient`
- New struct `VolumeQuotaRulesClientBeginCreateOptions`
- New struct `VolumeQuotaRulesClientBeginDeleteOptions`
- New struct `VolumeQuotaRulesClientBeginUpdateOptions`
- New struct `VolumeQuotaRulesClientCreateResponse`
- New struct `VolumeQuotaRulesClientDeleteResponse`
- New struct `VolumeQuotaRulesClientGetOptions`
- New struct `VolumeQuotaRulesClientGetResponse`
- New struct `VolumeQuotaRulesClientListByVolumeOptions`
- New struct `VolumeQuotaRulesClientListByVolumeResponse`
- New struct `VolumeQuotaRulesClientUpdateResponse`
- New struct `VolumeQuotaRulesList`
- New struct `VolumeQuotaRulesProperties`
- New struct `VolumeRelocationProperties`
- New struct `VolumesClientBeginFinalizeRelocationOptions`
- New struct `VolumesClientBeginRelocateOptions`
- New struct `VolumesClientBeginResetCifsPasswordOptions`
- New struct `VolumesClientBeginRevertRelocationOptions`
- New struct `VolumesClientFinalizeRelocationResponse`
- New struct `VolumesClientListReplicationsOptions`
- New struct `VolumesClientListReplicationsResponse`
- New struct `VolumesClientRelocateResponse`
- New struct `VolumesClientResetCifsPasswordResponse`
- New struct `VolumesClientRevertRelocationResponse`
- New field `Encrypted` in struct `VolumeProperties`
- New field `SystemData` in struct `Resource`
- New field `Zones` in struct `Volume`
- New field `SystemData` in struct `ProxyResource`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).