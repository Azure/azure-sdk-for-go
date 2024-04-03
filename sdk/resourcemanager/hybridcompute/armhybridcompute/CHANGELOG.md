# Release History

## 2.0.0-beta.2 (2024-04-26)
### Features Added

- New enum type `AccessMode` with values `AccessModeAudit`, `AccessModeEnforced`, `AccessModeLearning`
- New enum type `AccessRuleDirection` with values `AccessRuleDirectionInbound`, `AccessRuleDirectionOutbound`
- New enum type `ExecutionState` with values `ExecutionStateCanceled`, `ExecutionStateFailed`, `ExecutionStatePending`, `ExecutionStateRunning`, `ExecutionStateSucceeded`, `ExecutionStateTimedOut`, `ExecutionStateUnknown`
- New enum type `ExtensionsStatusLevelTypes` with values `ExtensionsStatusLevelTypesError`, `ExtensionsStatusLevelTypesInfo`, `ExtensionsStatusLevelTypesWarning`
- New enum type `LicenseProfileProductType` with values `LicenseProfileProductTypeWindowsIoTEnterprise`, `LicenseProfileProductTypeWindowsServer`
- New enum type `LicenseProfileSubscriptionStatus` with values `LicenseProfileSubscriptionStatusDisabled`, `LicenseProfileSubscriptionStatusEnabled`, `LicenseProfileSubscriptionStatusEnabling`, `LicenseProfileSubscriptionStatusUnknown`
- New enum type `LicenseProfileSubscriptionStatusUpdate` with values `LicenseProfileSubscriptionStatusUpdateDisable`, `LicenseProfileSubscriptionStatusUpdateEnable`
- New enum type `LicenseStatus` with values `LicenseStatusExtendedGrace`, `LicenseStatusLicensed`, `LicenseStatusNonGenuineGrace`, `LicenseStatusNotification`, `LicenseStatusOOBGrace`, `LicenseStatusOOTGrace`, `LicenseStatusUnlicensed`
- New enum type `ProvisioningIssueSeverity` with values `ProvisioningIssueSeverityError`, `ProvisioningIssueSeverityWarning`
- New enum type `ProvisioningIssueType` with values `ProvisioningIssueTypeConfigurationPropagationFailure`, `ProvisioningIssueTypeMissingIdentityConfiguration`, `ProvisioningIssueTypeMissingPerimeterConfiguration`, `ProvisioningIssueTypeOther`
- New function `*ClientFactory.NewMachineRunCommandsClient() *MachineRunCommandsClient`
- New function `*ClientFactory.NewNetworkConfigurationsClient() *NetworkConfigurationsClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `NewMachineRunCommandsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MachineRunCommandsClient, error)`
- New function `*MachineRunCommandsClient.BeginCreateOrUpdate(context.Context, string, string, string, MachineRunCommand, *MachineRunCommandsClientBeginCreateOrUpdateOptions) (*runtime.Poller[MachineRunCommandsClientCreateOrUpdateResponse], error)`
- New function `*MachineRunCommandsClient.BeginDelete(context.Context, string, string, string, *MachineRunCommandsClientBeginDeleteOptions) (*runtime.Poller[MachineRunCommandsClientDeleteResponse], error)`
- New function `*MachineRunCommandsClient.Get(context.Context, string, string, string, *MachineRunCommandsClientGetOptions) (MachineRunCommandsClientGetResponse, error)`
- New function `*MachineRunCommandsClient.NewListPager(string, string, *MachineRunCommandsClientListOptions) *runtime.Pager[MachineRunCommandsClientListResponse]`
- New function `*MachineRunCommandsClient.BeginUpdate(context.Context, string, string, string, MachineRunCommandUpdate, *MachineRunCommandsClientBeginUpdateOptions) (*runtime.Poller[MachineRunCommandsClientUpdateResponse], error)`
- New function `NewNetworkConfigurationsClient(azcore.TokenCredential, *arm.ClientOptions) (*NetworkConfigurationsClient, error)`
- New function `*NetworkConfigurationsClient.CreateOrUpdate(context.Context, string, NetworkConfiguration, *NetworkConfigurationsClientCreateOrUpdateOptions) (NetworkConfigurationsClientCreateOrUpdateResponse, error)`
- New function `*NetworkConfigurationsClient.Get(context.Context, string, *NetworkConfigurationsClientGetOptions) (NetworkConfigurationsClientGetResponse, error)`
- New function `*NetworkConfigurationsClient.Update(context.Context, string, NetworkConfiguration, *NetworkConfigurationsClientUpdateOptions) (NetworkConfigurationsClientUpdateResponse, error)`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.GetByPrivateLinkScope(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientGetByPrivateLinkScopeOptions) (NetworkSecurityPerimeterConfigurationsClientGetByPrivateLinkScopeResponse, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.NewListByPrivateLinkScopePager(string, string, *NetworkSecurityPerimeterConfigurationsClientListByPrivateLinkScopeOptions) *runtime.Pager[NetworkSecurityPerimeterConfigurationsClientListByPrivateLinkScopeResponse]`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginReconcileForPrivateLinkScope(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginReconcileForPrivateLinkScopeOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientReconcileForPrivateLinkScopeResponse], error)`
- New struct `AccessRule`
- New struct `AccessRuleProperties`
- New struct `ExtensionsResourceStatus`
- New struct `KeyDetails`
- New struct `KeyProperties`
- New struct `LicenseProfileArmProductProfileProperties`
- New struct `LicenseProfileMachineInstanceViewSoftwareAssurance`
- New struct `LicenseProfilePropertiesSoftwareAssurance`
- New struct `LicenseProfileUpdatePropertiesSoftwareAssurance`
- New struct `MachineRunCommand`
- New struct `MachineRunCommandInstanceView`
- New struct `MachineRunCommandProperties`
- New struct `MachineRunCommandScriptSource`
- New struct `MachineRunCommandUpdate`
- New struct `MachineRunCommandsListResult`
- New struct `NetworkConfiguration`
- New struct `NetworkConfigurationProperties`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationListResult`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityPerimeterProfile`
- New struct `ProductFeature`
- New struct `ProductFeatureUpdate`
- New struct `ProductProfileUpdateProperties`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `ResourceAssociation`
- New struct `RunCommandInputParameter`
- New struct `RunCommandManagedIdentity`
- New field `LicenseChannel`, `LicenseStatus`, `ProductProfile`, `SoftwareAssurance` in struct `LicenseProfileMachineInstanceView`
- New field `ProductProfile`, `SoftwareAssurance` in struct `LicenseProfileProperties`
- New field `ProductProfile`, `SoftwareAssurance` in struct `LicenseProfileUpdateProperties`
- New field `OSEdition` in struct `MachineProperties`


## 2.0.0-beta.1 (2023-11-24)
### Breaking Changes

- Type of `MachineExtensionProperties.ProtectedSettings` has been changed from `any` to `map[string]any`
- Type of `MachineExtensionProperties.Settings` has been changed from `any` to `map[string]any`
- Type of `MachineExtensionUpdateProperties.ProtectedSettings` has been changed from `any` to `map[string]any`
- Type of `MachineExtensionUpdateProperties.Settings` has been changed from `any` to `map[string]any`

### Features Added

- New enum type `AgentConfigurationMode` with values `AgentConfigurationModeFull`, `AgentConfigurationModeMonitor`
- New enum type `ArcKindEnum` with values `ArcKindEnumAVS`, `ArcKindEnumAWS`, `ArcKindEnumEPS`, `ArcKindEnumGCP`, `ArcKindEnumHCI`, `ArcKindEnumSCVMM`, `ArcKindEnumVMware`
- New enum type `EsuEligibility` with values `EsuEligibilityEligible`, `EsuEligibilityIneligible`, `EsuEligibilityUnknown`
- New enum type `EsuKeyState` with values `EsuKeyStateActive`, `EsuKeyStateInactive`
- New enum type `EsuServerType` with values `EsuServerTypeDatacenter`, `EsuServerTypeStandard`
- New enum type `LastAttemptStatusEnum` with values `LastAttemptStatusEnumFailed`, `LastAttemptStatusEnumSuccess`
- New enum type `LicenseAssignmentState` with values `LicenseAssignmentStateAssigned`, `LicenseAssignmentStateNotAssigned`
- New enum type `LicenseCoreType` with values `LicenseCoreTypePCore`, `LicenseCoreTypeVCore`
- New enum type `LicenseEdition` with values `LicenseEditionDatacenter`, `LicenseEditionStandard`
- New enum type `LicenseState` with values `LicenseStateActivated`, `LicenseStateDeactivated`
- New enum type `LicenseTarget` with values `LicenseTargetWindowsServer2012`, `LicenseTargetWindowsServer2012R2`
- New enum type `LicenseType` with values `LicenseTypeESU`
- New enum type `OsType` with values `OsTypeLinux`, `OsTypeWindows`
- New enum type `PatchOperationStartedBy` with values `PatchOperationStartedByPlatform`, `PatchOperationStartedByUser`
- New enum type `PatchOperationStatus` with values `PatchOperationStatusCompletedWithWarnings`, `PatchOperationStatusFailed`, `PatchOperationStatusInProgress`, `PatchOperationStatusSucceeded`, `PatchOperationStatusUnknown`
- New enum type `PatchServiceUsed` with values `PatchServiceUsedAPT`, `PatchServiceUsedUnknown`, `PatchServiceUsedWU`, `PatchServiceUsedWUWSUS`, `PatchServiceUsedYUM`, `PatchServiceUsedZypper`
- New enum type `ProvisioningState` with values `ProvisioningStateAccepted`, `ProvisioningStateCanceled`, `ProvisioningStateCreating`, `ProvisioningStateDeleted`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`, `ProvisioningStateUpdating`
- New enum type `VMGuestPatchClassificationLinux` with values `VMGuestPatchClassificationLinuxCritical`, `VMGuestPatchClassificationLinuxOther`, `VMGuestPatchClassificationLinuxSecurity`
- New enum type `VMGuestPatchClassificationWindows` with values `VMGuestPatchClassificationWindowsCritical`, `VMGuestPatchClassificationWindowsDefinition`, `VMGuestPatchClassificationWindowsFeaturePack`, `VMGuestPatchClassificationWindowsSecurity`, `VMGuestPatchClassificationWindowsServicePack`, `VMGuestPatchClassificationWindowsTools`, `VMGuestPatchClassificationWindowsUpdateRollUp`, `VMGuestPatchClassificationWindowsUpdates`
- New enum type `VMGuestPatchRebootSetting` with values `VMGuestPatchRebootSettingAlways`, `VMGuestPatchRebootSettingIfRequired`, `VMGuestPatchRebootSettingNever`
- New enum type `VMGuestPatchRebootStatus` with values `VMGuestPatchRebootStatusCompleted`, `VMGuestPatchRebootStatusFailed`, `VMGuestPatchRebootStatusNotNeeded`, `VMGuestPatchRebootStatusRequired`, `VMGuestPatchRebootStatusStarted`, `VMGuestPatchRebootStatusUnknown`
- New function `NewAgentVersionClient(azcore.TokenCredential, *arm.ClientOptions) (*AgentVersionClient, error)`
- New function `*AgentVersionClient.Get(context.Context, string, string, *AgentVersionClientGetOptions) (AgentVersionClientGetResponse, error)`
- New function `*AgentVersionClient.List(context.Context, string, *AgentVersionClientListOptions) (AgentVersionClientListResponse, error)`
- New function `*ClientFactory.NewAgentVersionClient() *AgentVersionClient`
- New function `*ClientFactory.NewExtensionMetadataClient() *ExtensionMetadataClient`
- New function `*ClientFactory.NewHybridIdentityMetadataClient() *HybridIdentityMetadataClient`
- New function `*ClientFactory.NewLicenseProfilesClient() *LicenseProfilesClient`
- New function `*ClientFactory.NewLicensesClient() *LicensesClient`
- New function `*ClientFactory.NewNetworkProfileClient() *NetworkProfileClient`
- New function `NewExtensionMetadataClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ExtensionMetadataClient, error)`
- New function `*ExtensionMetadataClient.Get(context.Context, string, string, string, string, *ExtensionMetadataClientGetOptions) (ExtensionMetadataClientGetResponse, error)`
- New function `*ExtensionMetadataClient.NewListPager(string, string, string, *ExtensionMetadataClientListOptions) *runtime.Pager[ExtensionMetadataClientListResponse]`
- New function `NewHybridIdentityMetadataClient(string, azcore.TokenCredential, *arm.ClientOptions) (*HybridIdentityMetadataClient, error)`
- New function `*HybridIdentityMetadataClient.Get(context.Context, string, string, string, *HybridIdentityMetadataClientGetOptions) (HybridIdentityMetadataClientGetResponse, error)`
- New function `*HybridIdentityMetadataClient.NewListByMachinesPager(string, string, *HybridIdentityMetadataClientListByMachinesOptions) *runtime.Pager[HybridIdentityMetadataClientListByMachinesResponse]`
- New function `NewLicenseProfilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LicenseProfilesClient, error)`
- New function `*LicenseProfilesClient.BeginCreateOrUpdate(context.Context, string, string, LicenseProfile, *LicenseProfilesClientBeginCreateOrUpdateOptions) (*runtime.Poller[LicenseProfilesClientCreateOrUpdateResponse], error)`
- New function `*LicenseProfilesClient.BeginDelete(context.Context, string, string, *LicenseProfilesClientBeginDeleteOptions) (*runtime.Poller[LicenseProfilesClientDeleteResponse], error)`
- New function `*LicenseProfilesClient.Get(context.Context, string, string, *LicenseProfilesClientGetOptions) (LicenseProfilesClientGetResponse, error)`
- New function `*LicenseProfilesClient.NewListPager(string, string, *LicenseProfilesClientListOptions) *runtime.Pager[LicenseProfilesClientListResponse]`
- New function `*LicenseProfilesClient.BeginUpdate(context.Context, string, string, LicenseProfileUpdate, *LicenseProfilesClientBeginUpdateOptions) (*runtime.Poller[LicenseProfilesClientUpdateResponse], error)`
- New function `NewLicensesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LicensesClient, error)`
- New function `*LicensesClient.BeginCreateOrUpdate(context.Context, string, string, License, *LicensesClientBeginCreateOrUpdateOptions) (*runtime.Poller[LicensesClientCreateOrUpdateResponse], error)`
- New function `*LicensesClient.BeginDelete(context.Context, string, string, *LicensesClientBeginDeleteOptions) (*runtime.Poller[LicensesClientDeleteResponse], error)`
- New function `*LicensesClient.Get(context.Context, string, string, *LicensesClientGetOptions) (LicensesClientGetResponse, error)`
- New function `*LicensesClient.NewListByResourceGroupPager(string, *LicensesClientListByResourceGroupOptions) *runtime.Pager[LicensesClientListByResourceGroupResponse]`
- New function `*LicensesClient.NewListBySubscriptionPager(*LicensesClientListBySubscriptionOptions) *runtime.Pager[LicensesClientListBySubscriptionResponse]`
- New function `*LicensesClient.BeginUpdate(context.Context, string, string, LicenseUpdate, *LicensesClientBeginUpdateOptions) (*runtime.Poller[LicensesClientUpdateResponse], error)`
- New function `*LicensesClient.BeginValidateLicense(context.Context, License, *LicensesClientBeginValidateLicenseOptions) (*runtime.Poller[LicensesClientValidateLicenseResponse], error)`
- New function `*MachinesClient.BeginAssessPatches(context.Context, string, string, *MachinesClientBeginAssessPatchesOptions) (*runtime.Poller[MachinesClientAssessPatchesResponse], error)`
- New function `*MachinesClient.BeginInstallPatches(context.Context, string, string, MachineInstallPatchesParameters, *MachinesClientBeginInstallPatchesOptions) (*runtime.Poller[MachinesClientInstallPatchesResponse], error)`
- New function `NewNetworkProfileClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkProfileClient, error)`
- New function `*NetworkProfileClient.Get(context.Context, string, string, *NetworkProfileClientGetOptions) (NetworkProfileClientGetResponse, error)`
- New struct `AgentUpgrade`
- New struct `AgentVersion`
- New struct `AgentVersionsList`
- New struct `AvailablePatchCountByClassification`
- New struct `EsuKey`
- New struct `EsuProfileUpdateProperties`
- New struct `ExtensionValue`
- New struct `ExtensionValueListResult`
- New struct `ExtensionValueProperties`
- New struct `HybridIdentityMetadata`
- New struct `HybridIdentityMetadataList`
- New struct `HybridIdentityMetadataProperties`
- New struct `IPAddress`
- New struct `License`
- New struct `LicenseDetails`
- New struct `LicenseProfile`
- New struct `LicenseProfileArmEsuProperties`
- New struct `LicenseProfileMachineInstanceView`
- New struct `LicenseProfileMachineInstanceViewEsuProperties`
- New struct `LicenseProfileProperties`
- New struct `LicenseProfileUpdate`
- New struct `LicenseProfileUpdateProperties`
- New struct `LicenseProfilesListResult`
- New struct `LicenseProperties`
- New struct `LicenseUpdate`
- New struct `LicenseUpdateProperties`
- New struct `LicenseUpdatePropertiesLicenseDetails`
- New struct `LicensesListResult`
- New struct `LinuxParameters`
- New struct `MachineAssessPatchesResult`
- New struct `MachineInstallPatchesParameters`
- New struct `MachineInstallPatchesResult`
- New struct `NetworkInterface`
- New struct `NetworkProfile`
- New struct `Subnet`
- New struct `WindowsParameters`
- New field `ConfigMode` in struct `AgentConfiguration`
- New field `Kind`, `Resources` in struct `Machine`
- New field `EnableAutomaticUpgrade` in struct `MachineExtensionUpdateProperties`
- New field `AgentUpgrade`, `LicenseProfile`, `NetworkProfile` in struct `MachineProperties`
- New field `Kind` in struct `MachineUpdate`
- New field `AgentUpgrade` in struct `MachineUpdateProperties`
- New field `Expand` in struct `MachinesClientCreateOrUpdateOptions`
- New field `Expand` in struct `MachinesClientListByResourceGroupOptions`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridcompute/armhybridcompute` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).