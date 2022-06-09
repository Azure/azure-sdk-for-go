# Release History

## 2.0.0 (2022-06-09)
### Breaking Changes

- Type of `VolumeProperties.EncryptionKeySource` has been changed from `*string` to `*EncryptionKeySource`
- Function `VolumeGroupDetails.MarshalJSON` has been removed
- Field `Tags` of struct `VolumeGroupDetails` has been removed
- Field `Tags` of struct `VolumeGroup` has been removed

### Features Added

- New const `TypeIndividualGroupQuota`
- New const `TypeIndividualUserQuota`
- New const `ProvisioningStateFailed`
- New const `TypeDefaultGroupQuota`
- New const `TypeDefaultUserQuota`
- New const `EncryptionKeySourceMicrosoftNetApp`
- New const `ProvisioningStateMoving`
- New const `ProvisioningStateDeleting`
- New const `ProvisioningStateCreating`
- New const `ProvisioningStatePatching`
- New const `ProvisioningStateSucceeded`
- New const `ProvisioningStateAccepted`
- New function `*VolumesClient.ListReplications(context.Context, string, string, string, string, *VolumesClientListReplicationsOptions) (VolumesClientListReplicationsResponse, error)`
- New function `PossibleTypeValues() []Type`
- New function `PossibleEncryptionKeySourceValues() []EncryptionKeySource`
- New function `TrackedResource.MarshalJSON() ([]byte, error)`
- New function `*VolumesClient.BeginRevertRelocation(context.Context, string, string, string, string, *VolumesClientBeginRevertRelocationOptions) (*runtime.Poller[VolumesClientRevertRelocationResponse], error)`
- New function `PossibleProvisioningStateValues() []ProvisioningState`
- New function `VolumeQuotaRule.MarshalJSON() ([]byte, error)`
- New function `*VolumesClient.BeginResetCifsPassword(context.Context, string, string, string, string, *VolumesClientBeginResetCifsPasswordOptions) (*runtime.Poller[VolumesClientResetCifsPasswordResponse], error)`
- New function `*VolumesClient.BeginRelocate(context.Context, string, string, string, string, *VolumesClientBeginRelocateOptions) (*runtime.Poller[VolumesClientRelocateResponse], error)`
- New function `VolumeQuotaRulePatch.MarshalJSON() ([]byte, error)`
- New function `*VolumesClient.BeginFinalizeRelocation(context.Context, string, string, string, string, *VolumesClientBeginFinalizeRelocationOptions) (*runtime.Poller[VolumesClientFinalizeRelocationResponse], error)`
- New struct `ListReplications`
- New struct `Replication`
- New struct `TrackedResource`
- New struct `VolumeQuotaRule`
- New struct `VolumeQuotaRulePatch`
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
- New field `Zones` in struct `Volume`
- New field `SystemData` in struct `Resource`
- New field `SystemData` in struct `ProxyResource`
- New field `Encrypted` in struct `VolumeProperties`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).