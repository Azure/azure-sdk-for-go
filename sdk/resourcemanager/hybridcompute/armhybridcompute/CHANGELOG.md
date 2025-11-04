# Release History

## 2.1.0-beta.2 (2025-06-10)
### Breaking Changes

- Operation `*MachinesClient.Delete` has been changed to LRO, use `*MachinesClient.BeginDelete` instead.

### Features Added

- New value `StatusTypesAwaitingConnection` added to enum type `StatusTypes`
- New enum type `IdentityKeyStore` with values `IdentityKeyStoreDefault`, `IdentityKeyStoreTPM`
- New function `*ClientFactory.NewExtensionMetadataV2Client() *ExtensionMetadataV2Client`
- New function `*ClientFactory.NewExtensionPublisherClient() *ExtensionPublisherClient`
- New function `*ClientFactory.NewExtensionTypeClient() *ExtensionTypeClient`
- New function `NewExtensionMetadataV2Client(azcore.TokenCredential, *arm.ClientOptions) (*ExtensionMetadataV2Client, error)`
- New function `*ExtensionMetadataV2Client.Get(context.Context, string, string, string, string, *ExtensionMetadataV2ClientGetOptions) (ExtensionMetadataV2ClientGetResponse, error)`
- New function `*ExtensionMetadataV2Client.NewListPager(string, string, string, *ExtensionMetadataV2ClientListOptions) *runtime.Pager[ExtensionMetadataV2ClientListResponse]`
- New function `NewExtensionPublisherClient(azcore.TokenCredential, *arm.ClientOptions) (*ExtensionPublisherClient, error)`
- New function `*ExtensionPublisherClient.NewListPager(string, *ExtensionPublisherClientListOptions) *runtime.Pager[ExtensionPublisherClientListResponse]`
- New function `NewExtensionTypeClient(azcore.TokenCredential, *arm.ClientOptions) (*ExtensionTypeClient, error)`
- New function `*ExtensionTypeClient.NewListPager(string, string, *ExtensionTypeClientListOptions) *runtime.Pager[ExtensionTypeClientListResponse]`
- New function `*ManagementClient.BeginSetupExtensions(context.Context, string, string, SetupExtensionRequest, *ManagementClientBeginSetupExtensionsOptions) (*runtime.Poller[ManagementClientSetupExtensionsResponse], error)`
- New struct `ExtensionPublisher`
- New struct `ExtensionPublisherListResult`
- New struct `ExtensionType`
- New struct `ExtensionTypeListResult`
- New struct `ExtensionValueListResultV2`
- New struct `ExtensionValueV2`
- New struct `ExtensionValueV2Properties`
- New struct `SetupExtensionRequest`
- New field `HardwareResourceID`, `IdentityKeyStore`, `TpmEkCertificate` in struct `MachineProperties`
- New field `IdentityKeyStore`, `TpmEkCertificate` in struct `MachineUpdateProperties`
- New field `PatchNameMasksToExclude`, `PatchNameMasksToInclude` in struct `WindowsParameters`


## 2.1.0-beta.1 (2024-11-14)
### Features Added

- New enum type `ExecutionState` with values `ExecutionStateCanceled`, `ExecutionStateFailed`, `ExecutionStatePending`, `ExecutionStateRunning`, `ExecutionStateSucceeded`, `ExecutionStateTimedOut`, `ExecutionStateUnknown`
- New enum type `ExtensionsStatusLevelTypes` with values `ExtensionsStatusLevelTypesError`, `ExtensionsStatusLevelTypesInfo`, `ExtensionsStatusLevelTypesWarning`
- New enum type `GatewayType` with values `GatewayTypePublic`
- New function `*ClientFactory.NewGatewaysClient() *GatewaysClient`
- New function `*ClientFactory.NewMachineRunCommandsClient() *MachineRunCommandsClient`
- New function `*ClientFactory.NewSettingsClient() *SettingsClient`
- New function `NewGatewaysClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GatewaysClient, error)`
- New function `*GatewaysClient.BeginCreateOrUpdate(context.Context, string, string, Gateway, *GatewaysClientBeginCreateOrUpdateOptions) (*runtime.Poller[GatewaysClientCreateOrUpdateResponse], error)`
- New function `*GatewaysClient.BeginDelete(context.Context, string, string, *GatewaysClientBeginDeleteOptions) (*runtime.Poller[GatewaysClientDeleteResponse], error)`
- New function `*GatewaysClient.Get(context.Context, string, string, *GatewaysClientGetOptions) (GatewaysClientGetResponse, error)`
- New function `*GatewaysClient.NewListByResourceGroupPager(string, *GatewaysClientListByResourceGroupOptions) *runtime.Pager[GatewaysClientListByResourceGroupResponse]`
- New function `*GatewaysClient.NewListBySubscriptionPager(*GatewaysClientListBySubscriptionOptions) *runtime.Pager[GatewaysClientListBySubscriptionResponse]`
- New function `*GatewaysClient.Update(context.Context, string, string, GatewayUpdate, *GatewaysClientUpdateOptions) (GatewaysClientUpdateResponse, error)`
- New function `NewMachineRunCommandsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MachineRunCommandsClient, error)`
- New function `*MachineRunCommandsClient.BeginCreateOrUpdate(context.Context, string, string, string, MachineRunCommand, *MachineRunCommandsClientBeginCreateOrUpdateOptions) (*runtime.Poller[MachineRunCommandsClientCreateOrUpdateResponse], error)`
- New function `*MachineRunCommandsClient.BeginDelete(context.Context, string, string, string, *MachineRunCommandsClientBeginDeleteOptions) (*runtime.Poller[MachineRunCommandsClientDeleteResponse], error)`
- New function `*MachineRunCommandsClient.Get(context.Context, string, string, string, *MachineRunCommandsClientGetOptions) (MachineRunCommandsClientGetResponse, error)`
- New function `*MachineRunCommandsClient.NewListPager(string, string, *MachineRunCommandsClientListOptions) *runtime.Pager[MachineRunCommandsClientListResponse]`
- New function `NewSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SettingsClient, error)`
- New function `*SettingsClient.Get(context.Context, string, string, string, string, string, *SettingsClientGetOptions) (SettingsClientGetResponse, error)`
- New function `*SettingsClient.Patch(context.Context, string, string, string, string, string, Settings, *SettingsClientPatchOptions) (SettingsClientPatchResponse, error)`
- New function `*SettingsClient.Update(context.Context, string, string, string, string, string, Settings, *SettingsClientUpdateOptions) (SettingsClientUpdateResponse, error)`
- New struct `Disk`
- New struct `ExtensionsResourceStatus`
- New struct `FirmwareProfile`
- New struct `Gateway`
- New struct `GatewayProperties`
- New struct `GatewayUpdate`
- New struct `GatewayUpdateProperties`
- New struct `GatewaysListResult`
- New struct `HardwareProfile`
- New struct `MachineRunCommand`
- New struct `MachineRunCommandInstanceView`
- New struct `MachineRunCommandProperties`
- New struct `MachineRunCommandScriptSource`
- New struct `MachineRunCommandsListResult`
- New struct `Processor`
- New struct `RunCommandInputParameter`
- New struct `RunCommandManagedIdentity`
- New struct `Settings`
- New struct `SettingsGatewayProperties`
- New struct `SettingsProperties`
- New struct `StorageProfile`
- New field `FirmwareProfile`, `HardwareProfile`, `StorageProfile` in struct `MachineProperties`
- New field `ID`, `MacAddress`, `Name` in struct `NetworkInterface`


## 2.0.0 (2024-10-11)
### Breaking Changes

- Type of `MachineExtensionProperties.ProtectedSettings` has been changed from `any` to `map[string]any`
- Type of `MachineExtensionProperties.Settings` has been changed from `any` to `map[string]any`
- Type of `MachineExtensionUpdateProperties.ProtectedSettings` has been changed from `any` to `map[string]any`
- Type of `MachineExtensionUpdateProperties.Settings` has been changed from `any` to `map[string]any`

### Features Added

- New value `PublicNetworkAccessTypeSecuredByPerimeter` added to enum type `PublicNetworkAccessType`
- New enum type `AccessMode` with values `AccessModeAudit`, `AccessModeEnforced`, `AccessModeLearning`
- New enum type `AccessRuleDirection` with values `AccessRuleDirectionInbound`, `AccessRuleDirectionOutbound`
- New enum type `AgentConfigurationMode` with values `AgentConfigurationModeFull`, `AgentConfigurationModeMonitor`
- New enum type `ArcKindEnum` with values `ArcKindEnumAVS`, `ArcKindEnumAWS`, `ArcKindEnumEPS`, `ArcKindEnumGCP`, `ArcKindEnumHCI`, `ArcKindEnumSCVMM`, `ArcKindEnumVMware`
- New enum type `EsuEligibility` with values `EsuEligibilityEligible`, `EsuEligibilityIneligible`, `EsuEligibilityUnknown`
- New enum type `EsuKeyState` with values `EsuKeyStateActive`, `EsuKeyStateInactive`
- New enum type `EsuServerType` with values `EsuServerTypeDatacenter`, `EsuServerTypeStandard`
- New enum type `HotpatchEnablementStatus` with values `HotpatchEnablementStatusActionRequired`, `HotpatchEnablementStatusDisabled`, `HotpatchEnablementStatusEnabled`, `HotpatchEnablementStatusPendingEvaluation`, `HotpatchEnablementStatusUnknown`
- New enum type `LastAttemptStatusEnum` with values `LastAttemptStatusEnumFailed`, `LastAttemptStatusEnumSuccess`
- New enum type `LicenseAssignmentState` with values `LicenseAssignmentStateAssigned`, `LicenseAssignmentStateNotAssigned`
- New enum type `LicenseCoreType` with values `LicenseCoreTypePCore`, `LicenseCoreTypeVCore`
- New enum type `LicenseEdition` with values `LicenseEditionDatacenter`, `LicenseEditionStandard`
- New enum type `LicenseProfileProductType` with values `LicenseProfileProductTypeWindowsIoTEnterprise`, `LicenseProfileProductTypeWindowsServer`
- New enum type `LicenseProfileSubscriptionStatus` with values `LicenseProfileSubscriptionStatusDisabled`, `LicenseProfileSubscriptionStatusDisabling`, `LicenseProfileSubscriptionStatusEnabled`, `LicenseProfileSubscriptionStatusEnabling`, `LicenseProfileSubscriptionStatusFailed`, `LicenseProfileSubscriptionStatusUnknown`
- New enum type `LicenseProfileSubscriptionStatusUpdate` with values `LicenseProfileSubscriptionStatusUpdateDisable`, `LicenseProfileSubscriptionStatusUpdateEnable`
- New enum type `LicenseState` with values `LicenseStateActivated`, `LicenseStateDeactivated`
- New enum type `LicenseStatus` with values `LicenseStatusExtendedGrace`, `LicenseStatusLicensed`, `LicenseStatusNonGenuineGrace`, `LicenseStatusNotification`, `LicenseStatusOOBGrace`, `LicenseStatusOOTGrace`, `LicenseStatusUnlicensed`
- New enum type `LicenseTarget` with values `LicenseTargetWindowsServer2012`, `LicenseTargetWindowsServer2012R2`
- New enum type `LicenseType` with values `LicenseTypeESU`
- New enum type `OsType` with values `OsTypeLinux`, `OsTypeWindows`
- New enum type `PatchOperationStartedBy` with values `PatchOperationStartedByPlatform`, `PatchOperationStartedByUser`
- New enum type `PatchOperationStatus` with values `PatchOperationStatusCompletedWithWarnings`, `PatchOperationStatusFailed`, `PatchOperationStatusInProgress`, `PatchOperationStatusSucceeded`, `PatchOperationStatusUnknown`
- New enum type `PatchServiceUsed` with values `PatchServiceUsedAPT`, `PatchServiceUsedUnknown`, `PatchServiceUsedWU`, `PatchServiceUsedWUWSUS`, `PatchServiceUsedYUM`, `PatchServiceUsedZypper`
- New enum type `ProgramYear` with values `ProgramYearYear1`, `ProgramYearYear2`, `ProgramYearYear3`
- New enum type `ProvisioningIssueSeverity` with values `ProvisioningIssueSeverityError`, `ProvisioningIssueSeverityWarning`
- New enum type `ProvisioningIssueType` with values `ProvisioningIssueTypeConfigurationPropagationFailure`, `ProvisioningIssueTypeMissingIdentityConfiguration`, `ProvisioningIssueTypeMissingPerimeterConfiguration`, `ProvisioningIssueTypeOther`
- New enum type `ProvisioningState` with values `ProvisioningStateAccepted`, `ProvisioningStateCanceled`, `ProvisioningStateCreating`, `ProvisioningStateDeleted`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`, `ProvisioningStateUpdating`
- New enum type `VMGuestPatchClassificationLinux` with values `VMGuestPatchClassificationLinuxCritical`, `VMGuestPatchClassificationLinuxOther`, `VMGuestPatchClassificationLinuxSecurity`
- New enum type `VMGuestPatchClassificationWindows` with values `VMGuestPatchClassificationWindowsCritical`, `VMGuestPatchClassificationWindowsDefinition`, `VMGuestPatchClassificationWindowsFeaturePack`, `VMGuestPatchClassificationWindowsSecurity`, `VMGuestPatchClassificationWindowsServicePack`, `VMGuestPatchClassificationWindowsTools`, `VMGuestPatchClassificationWindowsUpdateRollUp`, `VMGuestPatchClassificationWindowsUpdates`
- New enum type `VMGuestPatchRebootSetting` with values `VMGuestPatchRebootSettingAlways`, `VMGuestPatchRebootSettingIfRequired`, `VMGuestPatchRebootSettingNever`
- New enum type `VMGuestPatchRebootStatus` with values `VMGuestPatchRebootStatusCompleted`, `VMGuestPatchRebootStatusFailed`, `VMGuestPatchRebootStatusNotNeeded`, `VMGuestPatchRebootStatusRequired`, `VMGuestPatchRebootStatusStarted`, `VMGuestPatchRebootStatusUnknown`
- New function `*ClientFactory.NewExtensionMetadataClient() *ExtensionMetadataClient`
- New function `*ClientFactory.NewLicenseProfilesClient() *LicenseProfilesClient`
- New function `*ClientFactory.NewLicensesClient() *LicensesClient`
- New function `*ClientFactory.NewNetworkProfileClient() *NetworkProfileClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `NewExtensionMetadataClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ExtensionMetadataClient, error)`
- New function `*ExtensionMetadataClient.Get(context.Context, string, string, string, string, *ExtensionMetadataClientGetOptions) (ExtensionMetadataClientGetResponse, error)`
- New function `*ExtensionMetadataClient.NewListPager(string, string, string, *ExtensionMetadataClientListOptions) *runtime.Pager[ExtensionMetadataClientListResponse]`
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
- New function `PossibleLicenseProfileSubscriptionStatusValues() []LicenseProfileSubscriptionStatus`
- New function `NewNetworkProfileClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkProfileClient, error)`
- New function `*NetworkProfileClient.Get(context.Context, string, string, *NetworkProfileClientGetOptions) (NetworkProfileClientGetResponse, error)`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.GetByPrivateLinkScope(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientGetByPrivateLinkScopeOptions) (NetworkSecurityPerimeterConfigurationsClientGetByPrivateLinkScopeResponse, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.NewListByPrivateLinkScopePager(string, string, *NetworkSecurityPerimeterConfigurationsClientListByPrivateLinkScopeOptions) *runtime.Pager[NetworkSecurityPerimeterConfigurationsClientListByPrivateLinkScopeResponse]`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginReconcileForPrivateLinkScope(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginReconcileForPrivateLinkScopeOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientReconcileForPrivateLinkScopeResponse], error)`
- New struct `AccessRule`
- New struct `AccessRuleProperties`
- New struct `AgentUpgrade`
- New struct `AvailablePatchCountByClassification`
- New struct `EsuKey`
- New struct `EsuProfileUpdateProperties`
- New struct `ExtensionValue`
- New struct `ExtensionValueListResult`
- New struct `ExtensionValueProperties`
- New struct `IPAddress`
- New struct `License`
- New struct `LicenseDetails`
- New struct `LicenseProfile`
- New struct `LicenseProfileArmEsuProperties`
- New struct `LicenseProfileArmProductProfileProperties`
- New struct `LicenseProfileMachineInstanceView`
- New struct `LicenseProfileMachineInstanceViewEsuProperties`
- New struct `LicenseProfileMachineInstanceViewSoftwareAssurance`
- New struct `LicenseProfileProperties`
- New struct `LicenseProfilePropertiesSoftwareAssurance`
- New struct `LicenseProfileUpdate`
- New struct `LicenseProfileUpdateProperties`
- New struct `LicenseProfileUpdatePropertiesSoftwareAssurance`
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
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationListResult`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityPerimeterConfigurationReconcileResult`
- New struct `NetworkSecurityPerimeterProfile`
- New struct `PatchSettingsStatus`
- New struct `ProductFeature`
- New struct `ProductFeatureUpdate`
- New struct `ProductProfileUpdateProperties`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `ResourceAssociation`
- New struct `Subnet`
- New struct `VolumeLicenseDetails`
- New struct `WindowsParameters`
- New field `ConfigMode` in struct `AgentConfiguration`
- New field `Kind`, `Resources` in struct `Machine`
- New field `EnableAutomaticUpgrade` in struct `MachineExtensionUpdateProperties`
- New field `AgentUpgrade`, `LicenseProfile`, `NetworkProfile`, `OSEdition` in struct `MachineProperties`
- New field `Kind` in struct `MachineUpdate`
- New field `AgentUpgrade` in struct `MachineUpdateProperties`
- New field `Expand` in struct `MachinesClientCreateOrUpdateOptions`
- New field `Expand` in struct `MachinesClientListByResourceGroupOptions`
- New field `EnableHotpatching`, `Status` in struct `PatchSettings`


## 2.0.0-beta.4 (2024-07-23)
### Breaking Changes

- Function `*MachineRunCommandsClient.BeginUpdate` has been removed
- Struct `MachineRunCommandUpdate` has been removed

### Features Added

- New value `LicenseProfileSubscriptionStatusDisabling`, `LicenseProfileSubscriptionStatusFailed` added to enum type `LicenseProfileSubscriptionStatus`
- New enum type `HotpatchEnablementStatus` with values `HotpatchEnablementStatusActionRequired`, `HotpatchEnablementStatusDisabled`, `HotpatchEnablementStatusEnabled`, `HotpatchEnablementStatusPendingEvaluation`, `HotpatchEnablementStatusUnknown`
- New function `*MachinesClient.CreateOrUpdate(context.Context, string, string, Machine, *MachinesClientCreateOrUpdateOptions) (MachinesClientCreateOrUpdateResponse, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginReconcileForPrivateLinkScope(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginReconcileForPrivateLinkScopeOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientReconcileForPrivateLinkScopeResponse], error)`
- New struct `PatchSettingsStatus`
- New field `BillingEndDate`, `Error` in struct `LicenseProfileArmProductProfileProperties`
- New field `EnableHotpatching`, `Status` in struct `PatchSettings`
- New field `BillingEndDate`, `Error` in struct `ProductFeature`


## 2.0.0-beta.3 (2024-06-21)
### Breaking Changes

- Type of `EsuKey.LicenseStatus` has been changed from `*string` to `*int32`

### Features Added

- New value `PublicNetworkAccessTypeSecuredByPerimeter` added to enum type `PublicNetworkAccessType`
- New enum type `AccessMode` with values `AccessModeAudit`, `AccessModeEnforced`, `AccessModeLearning`
- New enum type `AccessRuleDirection` with values `AccessRuleDirectionInbound`, `AccessRuleDirectionOutbound`
- New enum type `GatewayType` with values `GatewayTypePublic`
- New enum type `ProgramYear` with values `ProgramYearYear1`, `ProgramYearYear2`, `ProgramYearYear3`
- New enum type `ProvisioningIssueSeverity` with values `ProvisioningIssueSeverityError`, `ProvisioningIssueSeverityWarning`
- New enum type `ProvisioningIssueType` with values `ProvisioningIssueTypeConfigurationPropagationFailure`, `ProvisioningIssueTypeMissingIdentityConfiguration`, `ProvisioningIssueTypeMissingPerimeterConfiguration`, `ProvisioningIssueTypeOther`
- New function `*ClientFactory.NewGatewaysClient() *GatewaysClient`
- New function `*ClientFactory.NewLicensesClient() *LicensesClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `*ClientFactory.NewSettingsClient() *SettingsClient`
- New function `NewGatewaysClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GatewaysClient, error)`
- New function `*GatewaysClient.BeginCreateOrUpdate(context.Context, string, string, Gateway, *GatewaysClientBeginCreateOrUpdateOptions) (*runtime.Poller[GatewaysClientCreateOrUpdateResponse], error)`
- New function `*GatewaysClient.BeginDelete(context.Context, string, string, *GatewaysClientBeginDeleteOptions) (*runtime.Poller[GatewaysClientDeleteResponse], error)`
- New function `*GatewaysClient.Get(context.Context, string, string, *GatewaysClientGetOptions) (GatewaysClientGetResponse, error)`
- New function `*GatewaysClient.NewListByResourceGroupPager(string, *GatewaysClientListByResourceGroupOptions) *runtime.Pager[GatewaysClientListByResourceGroupResponse]`
- New function `*GatewaysClient.NewListBySubscriptionPager(*GatewaysClientListBySubscriptionOptions) *runtime.Pager[GatewaysClientListBySubscriptionResponse]`
- New function `*GatewaysClient.Update(context.Context, string, string, GatewayUpdate, *GatewaysClientUpdateOptions) (GatewaysClientUpdateResponse, error)`
- New function `NewLicensesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LicensesClient, error)`
- New function `*LicensesClient.BeginCreateOrUpdate(context.Context, string, string, License, *LicensesClientBeginCreateOrUpdateOptions) (*runtime.Poller[LicensesClientCreateOrUpdateResponse], error)`
- New function `*LicensesClient.BeginDelete(context.Context, string, string, *LicensesClientBeginDeleteOptions) (*runtime.Poller[LicensesClientDeleteResponse], error)`
- New function `*LicensesClient.Get(context.Context, string, string, *LicensesClientGetOptions) (LicensesClientGetResponse, error)`
- New function `*LicensesClient.NewListByResourceGroupPager(string, *LicensesClientListByResourceGroupOptions) *runtime.Pager[LicensesClientListByResourceGroupResponse]`
- New function `*LicensesClient.NewListBySubscriptionPager(*LicensesClientListBySubscriptionOptions) *runtime.Pager[LicensesClientListBySubscriptionResponse]`
- New function `*LicensesClient.BeginUpdate(context.Context, string, string, LicenseUpdate, *LicensesClientBeginUpdateOptions) (*runtime.Poller[LicensesClientUpdateResponse], error)`
- New function `*MachineRunCommandsClient.BeginUpdate(context.Context, string, string, string, MachineRunCommandUpdate, *MachineRunCommandsClientBeginUpdateOptions) (*runtime.Poller[MachineRunCommandsClientUpdateResponse], error)`
- New function `NewSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SettingsClient, error)`
- New function `*SettingsClient.Get(context.Context, string, string, string, string, string, *SettingsClientGetOptions) (SettingsClientGetResponse, error)`
- New function `*SettingsClient.Patch(context.Context, string, string, string, string, string, Settings, *SettingsClientPatchOptions) (SettingsClientPatchResponse, error)`
- New function `*SettingsClient.Update(context.Context, string, string, string, string, string, Settings, *SettingsClientUpdateOptions) (SettingsClientUpdateResponse, error)`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.GetByPrivateLinkScope(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientGetByPrivateLinkScopeOptions) (NetworkSecurityPerimeterConfigurationsClientGetByPrivateLinkScopeResponse, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.NewListByPrivateLinkScopePager(string, string, *NetworkSecurityPerimeterConfigurationsClientListByPrivateLinkScopeOptions) *runtime.Pager[NetworkSecurityPerimeterConfigurationsClientListByPrivateLinkScopeResponse]`
- New struct `AccessRule`
- New struct `AccessRuleProperties`
- New struct `Gateway`
- New struct `GatewayProperties`
- New struct `GatewayUpdate`
- New struct `GatewayUpdateProperties`
- New struct `GatewaysListResult`
- New struct `LicenseUpdate`
- New struct `LicenseUpdateProperties`
- New struct `LicenseUpdatePropertiesLicenseDetails`
- New struct `LicensesListResult`
- New struct `MachineRunCommandUpdate`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationListResult`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityPerimeterProfile`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `ResourceAssociation`
- New struct `Settings`
- New struct `SettingsGatewayProperties`
- New struct `SettingsProperties`
- New struct `VolumeLicenseDetails`
- New field `VolumeLicenseDetails` in struct `LicenseDetails`


## 2.0.0-beta.2 (2024-04-26)
### Breaking Changes

- Type of `AgentUpgrade.LastAttemptTimestamp` has been changed from `*string` to `*time.Time`
- Type of `MachinesClientGetOptions.Expand` has been changed from `*InstanceViewTypes` to `*string`
- Enum `InstanceViewTypes` has been removed
- Function `NewAgentVersionClient` has been removed
- Function `*AgentVersionClient.Get` has been removed
- Function `*AgentVersionClient.List` has been removed
- Function `*ClientFactory.NewAgentVersionClient` has been removed
- Function `*ClientFactory.NewHybridIdentityMetadataClient` has been removed
- Function `*ClientFactory.NewLicenseProfilesClient` has been removed
- Function `*ClientFactory.NewLicensesClient` has been removed
- Function `NewHybridIdentityMetadataClient` has been removed
- Function `*HybridIdentityMetadataClient.Get` has been removed
- Function `*HybridIdentityMetadataClient.NewListByMachinesPager` has been removed
- Function `NewLicenseProfilesClient` has been removed
- Function `*LicenseProfilesClient.BeginCreateOrUpdate` has been removed
- Function `*LicenseProfilesClient.BeginDelete` has been removed
- Function `*LicenseProfilesClient.Get` has been removed
- Function `*LicenseProfilesClient.NewListPager` has been removed
- Function `*LicenseProfilesClient.BeginUpdate` has been removed
- Function `NewLicensesClient` has been removed
- Function `*LicensesClient.BeginCreateOrUpdate` has been removed
- Function `*LicensesClient.BeginDelete` has been removed
- Function `*LicensesClient.Get` has been removed
- Function `*LicensesClient.NewListByResourceGroupPager` has been removed
- Function `*LicensesClient.NewListBySubscriptionPager` has been removed
- Function `*LicensesClient.BeginUpdate` has been removed
- Function `*LicensesClient.BeginValidateLicense` has been removed
- Function `*MachinesClient.CreateOrUpdate` has been removed
- Struct `AgentVersion` has been removed
- Struct `AgentVersionsList` has been removed
- Struct `EsuProfileUpdateProperties` has been removed
- Struct `HybridIdentityMetadata` has been removed
- Struct `HybridIdentityMetadataList` has been removed
- Struct `HybridIdentityMetadataProperties` has been removed
- Struct `LicenseProfile` has been removed
- Struct `LicenseProfileArmEsuProperties` has been removed
- Struct `LicenseProfileProperties` has been removed
- Struct `LicenseProfileUpdate` has been removed
- Struct `LicenseProfileUpdateProperties` has been removed
- Struct `LicenseProfilesListResult` has been removed
- Struct `LicenseUpdate` has been removed
- Struct `LicenseUpdateProperties` has been removed
- Struct `LicenseUpdatePropertiesLicenseDetails` has been removed
- Struct `LicensesListResult` has been removed

### Features Added

- New enum type `ExecutionState` with values `ExecutionStateCanceled`, `ExecutionStateFailed`, `ExecutionStatePending`, `ExecutionStateRunning`, `ExecutionStateSucceeded`, `ExecutionStateTimedOut`, `ExecutionStateUnknown`
- New enum type `ExtensionsStatusLevelTypes` with values `ExtensionsStatusLevelTypesError`, `ExtensionsStatusLevelTypesInfo`, `ExtensionsStatusLevelTypesWarning`
- New enum type `LicenseProfileProductType` with values `LicenseProfileProductTypeWindowsIoTEnterprise`, `LicenseProfileProductTypeWindowsServer`
- New enum type `LicenseProfileSubscriptionStatus` with values `LicenseProfileSubscriptionStatusDisabled`, `LicenseProfileSubscriptionStatusEnabled`, `LicenseProfileSubscriptionStatusEnabling`, `LicenseProfileSubscriptionStatusUnknown`
- New enum type `LicenseStatus` with values `LicenseStatusExtendedGrace`, `LicenseStatusLicensed`, `LicenseStatusNonGenuineGrace`, `LicenseStatusNotification`, `LicenseStatusOOBGrace`, `LicenseStatusOOTGrace`, `LicenseStatusUnlicensed`
- New function `*ClientFactory.NewMachineRunCommandsClient() *MachineRunCommandsClient`
- New function `NewMachineRunCommandsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MachineRunCommandsClient, error)`
- New function `*MachineRunCommandsClient.BeginCreateOrUpdate(context.Context, string, string, string, MachineRunCommand, *MachineRunCommandsClientBeginCreateOrUpdateOptions) (*runtime.Poller[MachineRunCommandsClientCreateOrUpdateResponse], error)`
- New function `*MachineRunCommandsClient.BeginDelete(context.Context, string, string, string, *MachineRunCommandsClientBeginDeleteOptions) (*runtime.Poller[MachineRunCommandsClientDeleteResponse], error)`
- New function `*MachineRunCommandsClient.Get(context.Context, string, string, string, *MachineRunCommandsClientGetOptions) (MachineRunCommandsClientGetResponse, error)`
- New function `*MachineRunCommandsClient.NewListPager(string, string, *MachineRunCommandsClientListOptions) *runtime.Pager[MachineRunCommandsClientListResponse]`
- New struct `ExtensionsResourceStatus`
- New struct `LicenseProfileArmProductProfileProperties`
- New struct `LicenseProfileMachineInstanceViewSoftwareAssurance`
- New struct `MachineRunCommand`
- New struct `MachineRunCommandInstanceView`
- New struct `MachineRunCommandProperties`
- New struct `MachineRunCommandScriptSource`
- New struct `MachineRunCommandsListResult`
- New struct `ProductFeature`
- New struct `RunCommandInputParameter`
- New struct `RunCommandManagedIdentity`
- New field `LicenseChannel`, `LicenseStatus`, `ProductProfile`, `SoftwareAssurance` in struct `LicenseProfileMachineInstanceView`
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
