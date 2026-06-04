# Release History

## 2.0.0-beta.1 (2026-06-04)
### Breaking Changes

- Function `*ScheduledActionsClient.VirtualMachinesCancelOperations` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody CancelOperationsContent, options *ScheduledActionsClientVirtualMachinesCancelOperationsOptions)` to `(ctx context.Context, locationparameter string, requestBody CancelOperationsRequest, options *ScheduledActionsClientVirtualMachinesCancelOperationsOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesExecuteCreate` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody ExecuteCreateContent, options *ScheduledActionsClientVirtualMachinesExecuteCreateOptions)` to `(ctx context.Context, locationparameter string, requestBody ExecuteCreateRequest, options *ScheduledActionsClientVirtualMachinesExecuteCreateOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesExecuteDeallocate` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody ExecuteDeallocateContent, options *ScheduledActionsClientVirtualMachinesExecuteDeallocateOptions)` to `(ctx context.Context, locationparameter string, requestBody ExecuteDeallocateRequest, options *ScheduledActionsClientVirtualMachinesExecuteDeallocateOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesExecuteDelete` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody ExecuteDeleteContent, options *ScheduledActionsClientVirtualMachinesExecuteDeleteOptions)` to `(ctx context.Context, locationparameter string, requestBody ExecuteDeleteRequest, options *ScheduledActionsClientVirtualMachinesExecuteDeleteOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesExecuteHibernate` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody ExecuteHibernateContent, options *ScheduledActionsClientVirtualMachinesExecuteHibernateOptions)` to `(ctx context.Context, locationparameter string, requestBody ExecuteHibernateRequest, options *ScheduledActionsClientVirtualMachinesExecuteHibernateOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesExecuteStart` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody ExecuteStartContent, options *ScheduledActionsClientVirtualMachinesExecuteStartOptions)` to `(ctx context.Context, locationparameter string, requestBody ExecuteStartRequest, options *ScheduledActionsClientVirtualMachinesExecuteStartOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesGetOperationErrors` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody GetOperationErrorsContent, options *ScheduledActionsClientVirtualMachinesGetOperationErrorsOptions)` to `(ctx context.Context, locationparameter string, requestBody GetOperationErrorsRequest, options *ScheduledActionsClientVirtualMachinesGetOperationErrorsOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesGetOperationStatus` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody GetOperationStatusContent, options *ScheduledActionsClientVirtualMachinesGetOperationStatusOptions)` to `(ctx context.Context, locationparameter string, requestBody GetOperationStatusRequest, options *ScheduledActionsClientVirtualMachinesGetOperationStatusOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesSubmitDeallocate` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody SubmitDeallocateContent, options *ScheduledActionsClientVirtualMachinesSubmitDeallocateOptions)` to `(ctx context.Context, locationparameter string, requestBody SubmitDeallocateRequest, options *ScheduledActionsClientVirtualMachinesSubmitDeallocateOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesSubmitHibernate` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody SubmitHibernateContent, options *ScheduledActionsClientVirtualMachinesSubmitHibernateOptions)` to `(ctx context.Context, locationparameter string, requestBody SubmitHibernateRequest, options *ScheduledActionsClientVirtualMachinesSubmitHibernateOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesSubmitStart` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody SubmitStartContent, options *ScheduledActionsClientVirtualMachinesSubmitStartOptions)` to `(ctx context.Context, locationparameter string, requestBody SubmitStartRequest, options *ScheduledActionsClientVirtualMachinesSubmitStartOptions)`
- Struct `CancelOperationsContent` has been removed
- Struct `ExecuteCreateContent` has been removed
- Struct `ExecuteDeallocateContent` has been removed
- Struct `ExecuteDeleteContent` has been removed
- Struct `ExecuteHibernateContent` has been removed
- Struct `ExecuteStartContent` has been removed
- Struct `GetOperationErrorsContent` has been removed
- Struct `GetOperationStatusContent` has been removed
- Struct `SubmitDeallocateContent` has been removed
- Struct `SubmitHibernateContent` has been removed
- Struct `SubmitStartContent` has been removed
- Field `BaseProfile`, `ResourceOverrides` of struct `ResourceProvisionPayload` has been removed

### Features Added

- New value `ResourceOperationTypeCreate`, `ResourceOperationTypeDelete` added to enum type `ResourceOperationType`
- New enum type `AllocationStrategy` with values `AllocationStrategyCapacityOptimized`, `AllocationStrategyLowestPrice`, `AllocationStrategyPrioritized`
- New enum type `CachingTypes` with values `CachingTypesNone`, `CachingTypesReadOnly`, `CachingTypesReadWrite`
- New enum type `DeleteOptions` with values `DeleteOptionsDelete`, `DeleteOptionsDetach`
- New enum type `DiffDiskOptions` with values `DiffDiskOptionsLocal`
- New enum type `DiffDiskPlacement` with values `DiffDiskPlacementCacheDisk`, `DiffDiskPlacementNvmeDisk`, `DiffDiskPlacementResourceDisk`
- New enum type `DiskControllerTypes` with values `DiskControllerTypesNVMe`, `DiskControllerTypesSCSI`
- New enum type `DiskCreateOptionTypes` with values `DiskCreateOptionTypesAttach`, `DiskCreateOptionTypesCopy`, `DiskCreateOptionTypesEmpty`, `DiskCreateOptionTypesFromImage`, `DiskCreateOptionTypesRestore`
- New enum type `DiskDeleteOptionTypes` with values `DiskDeleteOptionTypesDelete`, `DiskDeleteOptionTypesDetach`
- New enum type `DiskDetachOptionTypes` with values `DiskDetachOptionTypesForceDetach`
- New enum type `DistributionStrategy` with values `DistributionStrategyBestEffortBalanced`, `DistributionStrategyBestEffortSingleZone`, `DistributionStrategyPrioritized`, `DistributionStrategyStrictBalanced`
- New enum type `DomainNameLabelScopeTypes` with values `DomainNameLabelScopeTypesNoReuse`, `DomainNameLabelScopeTypesResourceGroupReuse`, `DomainNameLabelScopeTypesSubscriptionReuse`, `DomainNameLabelScopeTypesTenantReuse`
- New enum type `ExtendedLocationType` with values `ExtendedLocationTypeCustomLocation`, `ExtendedLocationTypeEdgeZone`
- New enum type `IPVersions` with values `IPVersionsIPv4`, `IPVersionsIPv6`
- New enum type `LinuxPatchAssessmentMode` with values `LinuxPatchAssessmentModeAutomaticByPlatform`, `LinuxPatchAssessmentModeImageDefault`
- New enum type `LinuxVMGuestPatchAutomaticByPlatformRebootSetting` with values `LinuxVMGuestPatchAutomaticByPlatformRebootSettingAlways`, `LinuxVMGuestPatchAutomaticByPlatformRebootSettingIfRequired`, `LinuxVMGuestPatchAutomaticByPlatformRebootSettingNever`, `LinuxVMGuestPatchAutomaticByPlatformRebootSettingUnknown`
- New enum type `LinuxVMGuestPatchMode` with values `LinuxVMGuestPatchModeAutomaticByPlatform`, `LinuxVMGuestPatchModeImageDefault`
- New enum type `Mode` with values `ModeAudit`, `ModeEnforce`
- New enum type `Modes` with values `ModesAudit`, `ModesDisabled`, `ModesEnforce`
- New enum type `NetworkAPIVersion` with values `NetworkAPIVersion20201101`, `NetworkAPIVersion20221101`
- New enum type `NetworkInterfaceAuxiliaryMode` with values `NetworkInterfaceAuxiliaryModeAcceleratedConnections`, `NetworkInterfaceAuxiliaryModeFloating`, `NetworkInterfaceAuxiliaryModeNone`
- New enum type `NetworkInterfaceAuxiliarySKU` with values `NetworkInterfaceAuxiliarySKUA1`, `NetworkInterfaceAuxiliarySKUA2`, `NetworkInterfaceAuxiliarySKUA4`, `NetworkInterfaceAuxiliarySKUA8`, `NetworkInterfaceAuxiliarySKUNone`
- New enum type `OperatingSystemTypes` with values `OperatingSystemTypesLinux`, `OperatingSystemTypesWindows`
- New enum type `OsType` with values `OsTypeLinux`, `OsTypeWindows`
- New enum type `PriorityType` with values `PriorityTypeRegular`, `PriorityTypeSpot`
- New enum type `ProtocolTypes` with values `ProtocolTypesHTTP`, `ProtocolTypesHTTPS`
- New enum type `PublicIPAddressSKUName` with values `PublicIPAddressSKUNameBasic`, `PublicIPAddressSKUNameStandard`
- New enum type `PublicIPAddressSKUTier` with values `PublicIPAddressSKUTierGlobal`, `PublicIPAddressSKUTierRegional`
- New enum type `PublicIPAllocationMethod` with values `PublicIPAllocationMethodDynamic`, `PublicIPAllocationMethodStatic`
- New enum type `ResourceIdentityType` with values `ResourceIdentityTypeNone`, `ResourceIdentityTypeSystemAssigned`, `ResourceIdentityTypeSystemAssignedUserAssigned`, `ResourceIdentityTypeUserAssigned`
- New enum type `SecurityEncryptionTypes` with values `SecurityEncryptionTypesDiskWithVMGuestState`, `SecurityEncryptionTypesNonPersistedTPM`, `SecurityEncryptionTypesVMGuestStateOnly`
- New enum type `SecurityTypes` with values `SecurityTypesConfidentialVM`, `SecurityTypesTrustedLaunch`
- New enum type `SettingNames` with values `SettingNamesAutoLogon`, `SettingNamesFirstLogonCommands`
- New enum type `StorageAccountTypes` with values `StorageAccountTypesPremiumLRS`, `StorageAccountTypesPremiumV2LRS`, `StorageAccountTypesPremiumZRS`, `StorageAccountTypesStandardLRS`, `StorageAccountTypesStandardSSDLRS`, `StorageAccountTypesStandardSSDZRS`, `StorageAccountTypesUltraSSDLRS`
- New enum type `WindowsPatchAssessmentMode` with values `WindowsPatchAssessmentModeAutomaticByPlatform`, `WindowsPatchAssessmentModeImageDefault`
- New enum type `WindowsVMGuestPatchAutomaticByPlatformRebootSetting` with values `WindowsVMGuestPatchAutomaticByPlatformRebootSettingAlways`, `WindowsVMGuestPatchAutomaticByPlatformRebootSettingIfRequired`, `WindowsVMGuestPatchAutomaticByPlatformRebootSettingNever`, `WindowsVMGuestPatchAutomaticByPlatformRebootSettingUnknown`
- New enum type `WindowsVMGuestPatchMode` with values `WindowsVMGuestPatchModeAutomaticByOS`, `WindowsVMGuestPatchModeAutomaticByPlatform`, `WindowsVMGuestPatchModeManual`
- New enum type `ZonePlacementPolicyType` with values `ZonePlacementPolicyTypeAny`, `ZonePlacementPolicyTypeAuto`
- New function `PossibleModeValues() []Mode`
- New function `*ScheduledActionsClient.VirtualMachinesExecuteCreateFlex(ctx context.Context, locationparameter string, body ExecuteCreateFlexRequest, options *ScheduledActionsClientVirtualMachinesExecuteCreateFlexOptions) (ScheduledActionsClientVirtualMachinesExecuteCreateFlexResponse, error)`
- New struct `APIEntityReference`
- New struct `AdditionalCapabilities`
- New struct `AdditionalUnattendContent`
- New struct `AllInstancesDown`
- New struct `ApplicationProfile`
- New struct `BootDiagnostics`
- New struct `BulkActionVMExtension`
- New struct `BulkActionVMExtensionProperties`
- New struct `BulkActionVMProperties`
- New struct `BulkVMConfiguration`
- New struct `CancelOperationsRequest`
- New struct `CapacityReservationProfile`
- New struct `CreateFlexResourceOperationResponse`
- New struct `DataDisk`
- New struct `DiagnosticsProfile`
- New struct `DiffDiskSettings`
- New struct `DiskEncryptionSetParameters`
- New struct `DiskEncryptionSettings`
- New struct `EncryptionIdentity`
- New struct `EventGridAndResourceGraph`
- New struct `ExecuteCreateFlexRequest`
- New struct `ExecuteCreateRequest`
- New struct `ExecuteDeallocateRequest`
- New struct `ExecuteDeleteRequest`
- New struct `ExecuteHibernateRequest`
- New struct `ExecuteStartRequest`
- New struct `ExtendedLocation`
- New struct `FallbackOperationInfo`
- New struct `FlexProperties`
- New struct `GetOperationErrorsRequest`
- New struct `GetOperationStatusRequest`
- New struct `HardwareProfile`
- New struct `HostEndpointSettings`
- New struct `ImageReference`
- New struct `KeyVaultKeyReference`
- New struct `KeyVaultSecretReference`
- New struct `LinuxConfiguration`
- New struct `LinuxPatchSettings`
- New struct `LinuxVMGuestPatchAutomaticByPlatformSettings`
- New struct `ManagedDiskParameters`
- New struct `NetworkInterfaceReference`
- New struct `NetworkInterfaceReferenceProperties`
- New struct `NetworkProfile`
- New struct `OSDisk`
- New struct `OSImageNotificationProfile`
- New struct `OSProfile`
- New struct `PatchSettings`
- New struct `Placement`
- New struct `Plan`
- New struct `PriorityProfile`
- New struct `ProxyAgentSettings`
- New struct `PublicIPAddressSKU`
- New struct `ResourceProvisionFlexPayload`
- New struct `SSHConfiguration`
- New struct `SSHPublicKey`
- New struct `ScheduledEventsAdditionalPublishingTargets`
- New struct `ScheduledEventsPolicy`
- New struct `ScheduledEventsProfile`
- New struct `SecurityProfile`
- New struct `StorageProfile`
- New struct `SubResource`
- New struct `SubmitDeallocateRequest`
- New struct `SubmitHibernateRequest`
- New struct `SubmitStartRequest`
- New struct `TerminateNotificationProfile`
- New struct `UefiSettings`
- New struct `UserAssignedIdentitiesValue`
- New struct `UserInitiatedReboot`
- New struct `UserInitiatedRedeploy`
- New struct `VMDiskSecurityProfile`
- New struct `VMGalleryApplication`
- New struct `VMSizeProfile`
- New struct `VMSizeProperties`
- New struct `VaultCertificate`
- New struct `VaultSecretGroup`
- New struct `VirtualHardDisk`
- New struct `VirtualMachineIPTag`
- New struct `VirtualMachineIdentity`
- New struct `VirtualMachineNetworkInterfaceConfiguration`
- New struct `VirtualMachineNetworkInterfaceConfigurationProperties`
- New struct `VirtualMachineNetworkInterfaceDNSSettingsConfiguration`
- New struct `VirtualMachineNetworkInterfaceIPConfiguration`
- New struct `VirtualMachineNetworkInterfaceIPConfigurationProperties`
- New struct `VirtualMachinePublicIPAddressConfiguration`
- New struct `VirtualMachinePublicIPAddressConfigurationProperties`
- New struct `VirtualMachinePublicIPAddressDNSSettingsConfiguration`
- New struct `WinRMConfiguration`
- New struct `WinRMListener`
- New struct `WindowsConfiguration`
- New struct `WindowsVMGuestPatchAutomaticByPlatformSettings`
- New struct `ZoneAllocationPolicy`
- New struct `ZonePreference`
- New field `FallbackOperationInfo` in struct `ResourceOperationDetails`
- New field `VirtualMachineBaseProfile`, `VirtualMachineOverrides` in struct `ResourceProvisionPayload`
- New field `OnFailureAction` in struct `RetryPolicy`


## 1.2.0-beta.1 (2025-07-24)
### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `Language` with values `LanguageEnUs`
- New enum type `Month` with values `MonthAll`, `MonthApril`, `MonthAugust`, `MonthDecember`, `MonthFebruary`, `MonthJanuary`, `MonthJuly`, `MonthJune`, `MonthMarch`, `MonthMay`, `MonthNovember`, `MonthOctober`, `MonthSeptember`
- New enum type `NotificationType` with values `NotificationTypeEmail`
- New enum type `OccurrenceState` with values `OccurrenceStateCanceled`, `OccurrenceStateCancelling`, `OccurrenceStateCreated`, `OccurrenceStateFailed`, `OccurrenceStateRescheduling`, `OccurrenceStateScheduled`, `OccurrenceStateSucceeded`
- New enum type `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`
- New enum type `ResourceOperationStatus` with values `ResourceOperationStatusFailed`, `ResourceOperationStatusSucceeded`
- New enum type `ResourceProvisioningState` with values `ResourceProvisioningStateCanceled`, `ResourceProvisioningStateFailed`, `ResourceProvisioningStateSucceeded`
- New enum type `ResourceType` with values `ResourceTypeVirtualMachine`, `ResourceTypeVirtualMachineScaleSet`
- New enum type `ScheduledActionType` with values `ScheduledActionTypeDeallocate`, `ScheduledActionTypeHibernate`, `ScheduledActionTypeStart`
- New enum type `WeekDay` with values `WeekDayAll`, `WeekDayFriday`, `WeekDayMonday`, `WeekDaySaturday`, `WeekDaySunday`, `WeekDayThursday`, `WeekDayTuesday`, `WeekDayWednesday`
- New function `*ClientFactory.NewOccurrenceExtensionClient() *OccurrenceExtensionClient`
- New function `*ClientFactory.NewOccurrencesClient() *OccurrencesClient`
- New function `*ClientFactory.NewScheduledActionExtensionClient() *ScheduledActionExtensionClient`
- New function `NewOccurrenceExtensionClient(azcore.TokenCredential, *arm.ClientOptions) (*OccurrenceExtensionClient, error)`
- New function `*OccurrenceExtensionClient.NewListOccurrenceByVMsPager(string, *OccurrenceExtensionClientListOccurrenceByVMsOptions) *runtime.Pager[OccurrenceExtensionClientListOccurrenceByVMsResponse]`
- New function `NewOccurrencesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*OccurrencesClient, error)`
- New function `*OccurrencesClient.Cancel(context.Context, string, string, string, CancelOccurrenceRequest, *OccurrencesClientCancelOptions) (OccurrencesClientCancelResponse, error)`
- New function `*OccurrencesClient.BeginDelay(context.Context, string, string, string, DelayRequest, *OccurrencesClientBeginDelayOptions) (*runtime.Poller[OccurrencesClientDelayResponse], error)`
- New function `*OccurrencesClient.Get(context.Context, string, string, string, *OccurrencesClientGetOptions) (OccurrencesClientGetResponse, error)`
- New function `*OccurrencesClient.NewListByScheduledActionPager(string, string, *OccurrencesClientListByScheduledActionOptions) *runtime.Pager[OccurrencesClientListByScheduledActionResponse]`
- New function `*OccurrencesClient.NewListResourcesPager(string, string, string, *OccurrencesClientListResourcesOptions) *runtime.Pager[OccurrencesClientListResourcesResponse]`
- New function `NewScheduledActionExtensionClient(azcore.TokenCredential, *arm.ClientOptions) (*ScheduledActionExtensionClient, error)`
- New function `*ScheduledActionExtensionClient.NewListByVMsPager(string, *ScheduledActionExtensionClientListByVMsOptions) *runtime.Pager[ScheduledActionExtensionClientListByVMsResponse]`
- New function `*ScheduledActionsClient.AttachResources(context.Context, string, string, ResourceAttachRequest, *ScheduledActionsClientAttachResourcesOptions) (ScheduledActionsClientAttachResourcesResponse, error)`
- New function `*ScheduledActionsClient.CancelNextOccurrence(context.Context, string, string, CancelOccurrenceRequest, *ScheduledActionsClientCancelNextOccurrenceOptions) (ScheduledActionsClientCancelNextOccurrenceResponse, error)`
- New function `*ScheduledActionsClient.BeginCreateOrUpdate(context.Context, string, string, ScheduledAction, *ScheduledActionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ScheduledActionsClientCreateOrUpdateResponse], error)`
- New function `*ScheduledActionsClient.BeginDelete(context.Context, string, string, *ScheduledActionsClientBeginDeleteOptions) (*runtime.Poller[ScheduledActionsClientDeleteResponse], error)`
- New function `*ScheduledActionsClient.DetachResources(context.Context, string, string, ResourceDetachRequest, *ScheduledActionsClientDetachResourcesOptions) (ScheduledActionsClientDetachResourcesResponse, error)`
- New function `*ScheduledActionsClient.Disable(context.Context, string, string, *ScheduledActionsClientDisableOptions) (ScheduledActionsClientDisableResponse, error)`
- New function `*ScheduledActionsClient.Enable(context.Context, string, string, *ScheduledActionsClientEnableOptions) (ScheduledActionsClientEnableResponse, error)`
- New function `*ScheduledActionsClient.Get(context.Context, string, string, *ScheduledActionsClientGetOptions) (ScheduledActionsClientGetResponse, error)`
- New function `*ScheduledActionsClient.NewListByResourceGroupPager(string, *ScheduledActionsClientListByResourceGroupOptions) *runtime.Pager[ScheduledActionsClientListByResourceGroupResponse]`
- New function `*ScheduledActionsClient.NewListBySubscriptionPager(*ScheduledActionsClientListBySubscriptionOptions) *runtime.Pager[ScheduledActionsClientListBySubscriptionResponse]`
- New function `*ScheduledActionsClient.NewListResourcesPager(string, string, *ScheduledActionsClientListResourcesOptions) *runtime.Pager[ScheduledActionsClientListResourcesResponse]`
- New function `*ScheduledActionsClient.PatchResources(context.Context, string, string, ResourcePatchRequest, *ScheduledActionsClientPatchResourcesOptions) (ScheduledActionsClientPatchResourcesResponse, error)`
- New function `*ScheduledActionsClient.TriggerManualOccurrence(context.Context, string, string, *ScheduledActionsClientTriggerManualOccurrenceOptions) (ScheduledActionsClientTriggerManualOccurrenceResponse, error)`
- New function `*ScheduledActionsClient.Update(context.Context, string, string, ScheduledActionUpdate, *ScheduledActionsClientUpdateOptions) (ScheduledActionsClientUpdateResponse, error)`
- New struct `CancelOccurrenceRequest`
- New struct `DelayRequest`
- New struct `Error`
- New struct `InnerError`
- New struct `NotificationProperties`
- New struct `Occurrence`
- New struct `OccurrenceExtensionProperties`
- New struct `OccurrenceExtensionResource`
- New struct `OccurrenceExtensionResourceListResult`
- New struct `OccurrenceListResult`
- New struct `OccurrenceProperties`
- New struct `OccurrenceResource`
- New struct `OccurrenceResourceListResponse`
- New struct `OccurrenceResultSummary`
- New struct `RecurringActionsResourceOperationResult`
- New struct `ResourceAttachRequest`
- New struct `ResourceDetachRequest`
- New struct `ResourceListResponse`
- New struct `ResourcePatchRequest`
- New struct `ResourceResultSummary`
- New struct `ResourceStatus`
- New struct `ScheduledAction`
- New struct `ScheduledActionListResult`
- New struct `ScheduledActionProperties`
- New struct `ScheduledActionResource`
- New struct `ScheduledActionResources`
- New struct `ScheduledActionResourcesListResult`
- New struct `ScheduledActionUpdate`
- New struct `ScheduledActionUpdateProperties`
- New struct `ScheduledActionsSchedule`
- New struct `SystemData`


## 1.1.0 (2025-07-15)
### Features Added

- New function `*ScheduledActionsClient.VirtualMachinesExecuteCreate(context.Context, string, ExecuteCreateRequest, *ScheduledActionsClientVirtualMachinesExecuteCreateOptions) (ScheduledActionsClientVirtualMachinesExecuteCreateResponse, error)`
- New function `*ScheduledActionsClient.VirtualMachinesExecuteDelete(context.Context, string, ExecuteDeleteRequest, *ScheduledActionsClientVirtualMachinesExecuteDeleteOptions) (ScheduledActionsClientVirtualMachinesExecuteDeleteResponse, error)`
- New struct `CreateResourceOperationResponse`
- New struct `DeleteResourceOperationResponse`
- New struct `ExecuteCreateRequest`
- New struct `ExecuteDeleteRequest`
- New struct `ResourceProvisionPayload`


## 1.0.0 (2025-01-24)
### Breaking Changes

- Type of `OperationErrorDetails.ErrorDetails` has been changed from `*time.Time` to `*string`

### Features Added

- New field `AzureOperationName`, `Timestamp` in struct `OperationErrorDetails`
- New field `Timezone` in struct `ResourceOperationDetails`
- New field `Deadline`, `Timezone` in struct `Schedule`


## 0.1.0 (2024-09-27)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/computeschedule/armcomputeschedule` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
