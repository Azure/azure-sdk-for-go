# Release History

## 2.1.0 (2022-07-21)
### Features Added

- New const `EncryptionKeySourceMicrosoftKeyVault`
- New function `*VolumesClient.BeginReestablishReplication(context.Context, string, string, string, string, ReestablishReplicationRequest, *VolumesClientBeginReestablishReplicationOptions) (*runtime.Poller[VolumesClientReestablishReplicationResponse], error)`
- New struct `ReestablishReplicationRequest`
- New struct `VolumesClientBeginReestablishReplicationOptions`
- New struct `VolumesClientReestablishReplicationResponse`
- New field `CoolnessPeriod` in struct `VolumePatchProperties`
- New field `CoolAccess` in struct `VolumePatchProperties`
- New field `KeyVaultPrivateEndpointResourceID` in struct `VolumeProperties`
- New field `CoolAccess` in struct `PoolPatchProperties`


## 2.0.0 (2022-06-22)
### Breaking Changes

- Function `*BackupsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, string, string, *BackupsClientBeginUpdateOptions)` to `(context.Context, string, string, string, string, string, BackupPatch, *BackupsClientBeginUpdateOptions)`
- Type of `VolumeProperties.EncryptionKeySource` has been changed from `*string` to `*EncryptionKeySource`
- Field `Body` of struct `BackupsClientBeginUpdateOptions` has been removed
- Field `Tags` of struct `VolumeGroupDetails` has been removed
- Field `Tags` of struct `VolumeGroup` has been removed

### Features Added

- New const `ProvisioningStateDeleting`
- New const `TypeIndividualGroupQuota`
- New const `TypeDefaultUserQuota`
- New const `TypeIndividualUserQuota`
- New const `ProvisioningStateSucceeded`
- New const `ProvisioningStateAccepted`
- New const `ProvisioningStateCreating`
- New const `ProvisioningStateMoving`
- New const `ProvisioningStatePatching`
- New const `ProvisioningStateFailed`
- New const `EncryptionKeySourceMicrosoftNetApp`
- New const `TypeDefaultGroupQuota`
- New function `*VolumesClient.BeginRelocate(context.Context, string, string, string, string, *VolumesClientBeginRelocateOptions) (*runtime.Poller[VolumesClientRelocateResponse], error)`
- New function `*VolumeQuotaRulesClient.BeginDelete(context.Context, string, string, string, string, string, *VolumeQuotaRulesClientBeginDeleteOptions) (*runtime.Poller[VolumeQuotaRulesClientDeleteResponse], error)`
- New function `*VolumeQuotaRulesClient.Get(context.Context, string, string, string, string, string, *VolumeQuotaRulesClientGetOptions) (VolumeQuotaRulesClientGetResponse, error)`
- New function `*VolumesClient.NewListReplicationsPager(string, string, string, string, *VolumesClientListReplicationsOptions) *runtime.Pager[VolumesClientListReplicationsResponse]`
- New function `PossibleTypeValues() []Type`
- New function `*VolumesClient.BeginResetCifsPassword(context.Context, string, string, string, string, *VolumesClientBeginResetCifsPasswordOptions) (*runtime.Poller[VolumesClientResetCifsPasswordResponse], error)`
- New function `*VolumeQuotaRulesClient.BeginCreate(context.Context, string, string, string, string, string, VolumeQuotaRule, *VolumeQuotaRulesClientBeginCreateOptions) (*runtime.Poller[VolumeQuotaRulesClientCreateResponse], error)`
- New function `*VolumesClient.BeginFinalizeRelocation(context.Context, string, string, string, string, *VolumesClientBeginFinalizeRelocationOptions) (*runtime.Poller[VolumesClientFinalizeRelocationResponse], error)`
- New function `*VolumeQuotaRulesClient.NewListByVolumePager(string, string, string, string, *VolumeQuotaRulesClientListByVolumeOptions) *runtime.Pager[VolumeQuotaRulesClientListByVolumeResponse]`
- New function `PossibleProvisioningStateValues() []ProvisioningState`
- New function `*VolumesClient.BeginRevertRelocation(context.Context, string, string, string, string, *VolumesClientBeginRevertRelocationOptions) (*runtime.Poller[VolumesClientRevertRelocationResponse], error)`
- New function `*VolumeQuotaRulesClient.BeginUpdate(context.Context, string, string, string, string, string, VolumeQuotaRulePatch, *VolumeQuotaRulesClientBeginUpdateOptions) (*runtime.Poller[VolumeQuotaRulesClientUpdateResponse], error)`
- New function `PossibleEncryptionKeySourceValues() []EncryptionKeySource`
- New function `NewVolumeQuotaRulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VolumeQuotaRulesClient, error)`
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
- New field `SystemData` in struct `ProxyResource`
- New field `Zones` in struct `Volume`
- New field `SystemData` in struct `Resource`
- New field `Encrypted` in struct `VolumeProperties`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).