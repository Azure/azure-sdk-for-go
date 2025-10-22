# Release History

## 2.2.0 (2025-10-21)
### Features Added

- New enum type `BlockedDatesConstraintCategory` with values `BlockedDatesConstraintCategoryHiPriorityEvent`, `BlockedDatesConstraintCategoryHoliday`, `BlockedDatesConstraintCategoryQuotaExhausted`
- New enum type `LicenseKind` with values `LicenseKindVmwareFirewall`
- New enum type `LicenseName` with values `LicenseNameVmwareFirewall`
- New enum type `LicenseProvisioningState` with values `LicenseProvisioningStateCanceled`, `LicenseProvisioningStateFailed`, `LicenseProvisioningStateSucceeded`
- New enum type `MaintenanceCheckType` with values `MaintenanceCheckTypePrecheck`, `MaintenanceCheckTypePreflight`
- New enum type `MaintenanceManagementOperationKind` with values `MaintenanceManagementOperationKindMaintenanceReadinessRefresh`, `MaintenanceManagementOperationKindReschedule`, `MaintenanceManagementOperationKindSchedule`
- New enum type `MaintenanceProvisioningState` with values `MaintenanceProvisioningStateCanceled`, `MaintenanceProvisioningStateFailed`, `MaintenanceProvisioningStateSucceeded`, `MaintenanceProvisioningStateUpdating`
- New enum type `MaintenanceReadinessRefreshOperationStatus` with values `MaintenanceReadinessRefreshOperationStatusFailed`, `MaintenanceReadinessRefreshOperationStatusInProgress`, `MaintenanceReadinessRefreshOperationStatusNotApplicable`, `MaintenanceReadinessRefreshOperationStatusNotStarted`
- New enum type `MaintenanceReadinessStatus` with values `MaintenanceReadinessStatusDataNotAvailable`, `MaintenanceReadinessStatusNotApplicable`, `MaintenanceReadinessStatusNotReady`, `MaintenanceReadinessStatusReady`
- New enum type `MaintenanceStateName` with values `MaintenanceStateNameCanceled`, `MaintenanceStateNameFailed`, `MaintenanceStateNameInProgress`, `MaintenanceStateNameNotScheduled`, `MaintenanceStateNameScheduled`, `MaintenanceStateNameSuccess`
- New enum type `MaintenanceStatusFilter` with values `MaintenanceStatusFilterActive`, `MaintenanceStatusFilterInactive`
- New enum type `MaintenanceType` with values `MaintenanceTypeESXI`, `MaintenanceTypeNSXT`, `MaintenanceTypeVCSA`
- New enum type `RescheduleOperationConstraintKind` with values `RescheduleOperationConstraintKindAvailableWindowForMaintenanceWhileRescheduleOperation`, `RescheduleOperationConstraintKindBlockedWhileRescheduleOperation`
- New enum type `ScheduleOperationConstraintKind` with values `ScheduleOperationConstraintKindAvailableWindowForMaintenanceWhileScheduleOperation`, `ScheduleOperationConstraintKindBlockedWhileScheduleOperation`, `ScheduleOperationConstraintKindSchedulingWindow`
- New enum type `VcfLicenseKind` with values `VcfLicenseKindVcf5`
- New function `*AvailableWindowForMaintenanceWhileRescheduleOperation.GetRescheduleOperationConstraint() *RescheduleOperationConstraint`
- New function `*AvailableWindowForMaintenanceWhileScheduleOperation.GetScheduleOperationConstraint() *ScheduleOperationConstraint`
- New function `*BlockedWhileRescheduleOperation.GetRescheduleOperationConstraint() *RescheduleOperationConstraint`
- New function `*BlockedWhileScheduleOperation.GetScheduleOperationConstraint() *ScheduleOperationConstraint`
- New function `*ClientFactory.NewLicensesClient() *LicensesClient`
- New function `*ClientFactory.NewMaintenancesClient() *MaintenancesClient`
- New function `*LicenseProperties.GetLicenseProperties() *LicenseProperties`
- New function `NewLicensesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LicensesClient, error)`
- New function `*LicensesClient.BeginCreateOrUpdate(context.Context, string, string, LicenseName, License, *LicensesClientBeginCreateOrUpdateOptions) (*runtime.Poller[LicensesClientCreateOrUpdateResponse], error)`
- New function `*LicensesClient.BeginDelete(context.Context, string, string, LicenseName, *LicensesClientBeginDeleteOptions) (*runtime.Poller[LicensesClientDeleteResponse], error)`
- New function `*LicensesClient.Get(context.Context, string, string, LicenseName, *LicensesClientGetOptions) (LicensesClientGetResponse, error)`
- New function `*LicensesClient.GetProperties(context.Context, string, string, LicenseName, *LicensesClientGetPropertiesOptions) (LicensesClientGetPropertiesResponse, error)`
- New function `*LicensesClient.NewListPager(string, string, *LicensesClientListOptions) *runtime.Pager[LicensesClientListResponse]`
- New function `*MaintenanceManagementOperation.GetMaintenanceManagementOperation() *MaintenanceManagementOperation`
- New function `*MaintenanceReadinessRefreshOperation.GetMaintenanceManagementOperation() *MaintenanceManagementOperation`
- New function `NewMaintenancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MaintenancesClient, error)`
- New function `*MaintenancesClient.Get(context.Context, string, string, string, *MaintenancesClientGetOptions) (MaintenancesClientGetResponse, error)`
- New function `*MaintenancesClient.InitiateChecks(context.Context, string, string, string, *MaintenancesClientInitiateChecksOptions) (MaintenancesClientInitiateChecksResponse, error)`
- New function `*MaintenancesClient.NewListPager(string, string, *MaintenancesClientListOptions) *runtime.Pager[MaintenancesClientListResponse]`
- New function `*MaintenancesClient.Reschedule(context.Context, string, string, string, MaintenanceReschedule, *MaintenancesClientRescheduleOptions) (MaintenancesClientRescheduleResponse, error)`
- New function `*MaintenancesClient.Schedule(context.Context, string, string, string, MaintenanceSchedule, *MaintenancesClientScheduleOptions) (MaintenancesClientScheduleResponse, error)`
- New function `*PrivateCloudsClient.GetVcfLicense(context.Context, string, string, *PrivateCloudsClientGetVcfLicenseOptions) (PrivateCloudsClientGetVcfLicenseResponse, error)`
- New function `*RescheduleOperation.GetMaintenanceManagementOperation() *MaintenanceManagementOperation`
- New function `*RescheduleOperationConstraint.GetRescheduleOperationConstraint() *RescheduleOperationConstraint`
- New function `*ScheduleOperation.GetMaintenanceManagementOperation() *MaintenanceManagementOperation`
- New function `*ScheduleOperationConstraint.GetScheduleOperationConstraint() *ScheduleOperationConstraint`
- New function `*SchedulingWindow.GetScheduleOperationConstraint() *ScheduleOperationConstraint`
- New function `*Vcf5License.GetVcfLicense() *VcfLicense`
- New function `*VcfLicense.GetVcfLicense() *VcfLicense`
- New function `*VmwareFirewallLicenseProperties.GetLicenseProperties() *LicenseProperties`
- New struct `AvailableWindowForMaintenanceWhileRescheduleOperation`
- New struct `AvailableWindowForMaintenanceWhileScheduleOperation`
- New struct `BlockedDatesConstraintTimeRange`
- New struct `BlockedWhileRescheduleOperation`
- New struct `BlockedWhileScheduleOperation`
- New struct `ImpactedMaintenanceResource`
- New struct `ImpactedMaintenanceResourceError`
- New struct `Label`
- New struct `License`
- New struct `LicenseListResult`
- New struct `Maintenance`
- New struct `MaintenanceFailedCheck`
- New struct `MaintenanceListResult`
- New struct `MaintenanceProperties`
- New struct `MaintenanceReadiness`
- New struct `MaintenanceReadinessRefreshOperation`
- New struct `MaintenanceReschedule`
- New struct `MaintenanceSchedule`
- New struct `MaintenanceState`
- New struct `RescheduleOperation`
- New struct `ScheduleOperation`
- New struct `SchedulingWindow`
- New struct `Vcf5License`
- New struct `VmwareFirewallLicenseProperties`
- New field `VcfLicense` in struct `PrivateCloudProperties`


## 2.2.0 (2025-10-16)
### Features Added

- New enum type `BlockedDatesConstraintCategory` with values `BlockedDatesConstraintCategoryHiPriorityEvent`, `BlockedDatesConstraintCategoryHoliday`, `BlockedDatesConstraintCategoryQuotaExhausted`
- New enum type `LicenseKind` with values `LicenseKindVmwareFirewall`
- New enum type `LicenseName` with values `LicenseNameVmwareFirewall`
- New enum type `LicenseProvisioningState` with values `LicenseProvisioningStateCanceled`, `LicenseProvisioningStateFailed`, `LicenseProvisioningStateSucceeded`
- New enum type `MaintenanceCheckType` with values `MaintenanceCheckTypePrecheck`, `MaintenanceCheckTypePreflight`
- New enum type `MaintenanceManagementOperationKind` with values `MaintenanceManagementOperationKindMaintenanceReadinessRefresh`, `MaintenanceManagementOperationKindReschedule`, `MaintenanceManagementOperationKindSchedule`
- New enum type `MaintenanceProvisioningState` with values `MaintenanceProvisioningStateCanceled`, `MaintenanceProvisioningStateFailed`, `MaintenanceProvisioningStateSucceeded`, `MaintenanceProvisioningStateUpdating`
- New enum type `MaintenanceReadinessRefreshOperationStatus` with values `MaintenanceReadinessRefreshOperationStatusFailed`, `MaintenanceReadinessRefreshOperationStatusInProgress`, `MaintenanceReadinessRefreshOperationStatusNotApplicable`, `MaintenanceReadinessRefreshOperationStatusNotStarted`
- New enum type `MaintenanceReadinessStatus` with values `MaintenanceReadinessStatusDataNotAvailable`, `MaintenanceReadinessStatusNotApplicable`, `MaintenanceReadinessStatusNotReady`, `MaintenanceReadinessStatusReady`
- New enum type `MaintenanceStateName` with values `MaintenanceStateNameCanceled`, `MaintenanceStateNameFailed`, `MaintenanceStateNameInProgress`, `MaintenanceStateNameNotScheduled`, `MaintenanceStateNameScheduled`, `MaintenanceStateNameSuccess`
- New enum type `MaintenanceStatusFilter` with values `MaintenanceStatusFilterActive`, `MaintenanceStatusFilterInactive`
- New enum type `MaintenanceType` with values `MaintenanceTypeESXI`, `MaintenanceTypeNSXT`, `MaintenanceTypeVCSA`
- New enum type `RescheduleOperationConstraintKind` with values `RescheduleOperationConstraintKindAvailableWindowForMaintenanceWhileRescheduleOperation`, `RescheduleOperationConstraintKindBlockedWhileRescheduleOperation`
- New enum type `ScheduleOperationConstraintKind` with values `ScheduleOperationConstraintKindAvailableWindowForMaintenanceWhileScheduleOperation`, `ScheduleOperationConstraintKindBlockedWhileScheduleOperation`, `ScheduleOperationConstraintKindSchedulingWindow`
- New enum type `VcfLicenseKind` with values `VcfLicenseKindVcf5`
- New function `*AvailableWindowForMaintenanceWhileRescheduleOperation.GetRescheduleOperationConstraint() *RescheduleOperationConstraint`
- New function `*AvailableWindowForMaintenanceWhileScheduleOperation.GetScheduleOperationConstraint() *ScheduleOperationConstraint`
- New function `*BlockedWhileRescheduleOperation.GetRescheduleOperationConstraint() *RescheduleOperationConstraint`
- New function `*BlockedWhileScheduleOperation.GetScheduleOperationConstraint() *ScheduleOperationConstraint`
- New function `*ClientFactory.NewLicensesClient() *LicensesClient`
- New function `*ClientFactory.NewMaintenancesClient() *MaintenancesClient`
- New function `*ClientFactory.NewServiceComponentsClient() *ServiceComponentsClient`
- New function `*LicenseProperties.GetLicenseProperties() *LicenseProperties`
- New function `NewLicensesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LicensesClient, error)`
- New function `*LicensesClient.BeginCreateOrUpdate(context.Context, string, string, LicenseName, License, *LicensesClientBeginCreateOrUpdateOptions) (*runtime.Poller[LicensesClientCreateOrUpdateResponse], error)`
- New function `*LicensesClient.BeginDelete(context.Context, string, string, LicenseName, *LicensesClientBeginDeleteOptions) (*runtime.Poller[LicensesClientDeleteResponse], error)`
- New function `*LicensesClient.Get(context.Context, string, string, LicenseName, *LicensesClientGetOptions) (LicensesClientGetResponse, error)`
- New function `*LicensesClient.GetProperties(context.Context, string, string, LicenseName, *LicensesClientGetPropertiesOptions) (LicensesClientGetPropertiesResponse, error)`
- New function `*LicensesClient.NewListPager(string, string, *LicensesClientListOptions) *runtime.Pager[LicensesClientListResponse]`
- New function `*MaintenanceManagementOperation.GetMaintenanceManagementOperation() *MaintenanceManagementOperation`
- New function `*MaintenanceReadinessRefreshOperation.GetMaintenanceManagementOperation() *MaintenanceManagementOperation`
- New function `NewMaintenancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MaintenancesClient, error)`
- New function `*MaintenancesClient.Get(context.Context, string, string, string, *MaintenancesClientGetOptions) (MaintenancesClientGetResponse, error)`
- New function `*MaintenancesClient.InitiateChecks(context.Context, string, string, string, *MaintenancesClientInitiateChecksOptions) (MaintenancesClientInitiateChecksResponse, error)`
- New function `*MaintenancesClient.NewListPager(string, string, *MaintenancesClientListOptions) *runtime.Pager[MaintenancesClientListResponse]`
- New function `*MaintenancesClient.Reschedule(context.Context, string, string, string, MaintenanceReschedule, *MaintenancesClientRescheduleOptions) (MaintenancesClientRescheduleResponse, error)`
- New function `*MaintenancesClient.Schedule(context.Context, string, string, string, MaintenanceSchedule, *MaintenancesClientScheduleOptions) (MaintenancesClientScheduleResponse, error)`
- New function `*PrivateCloudsClient.GetVcfLicense(context.Context, string, string, *PrivateCloudsClientGetVcfLicenseOptions) (PrivateCloudsClientGetVcfLicenseResponse, error)`
- New function `*RescheduleOperation.GetMaintenanceManagementOperation() *MaintenanceManagementOperation`
- New function `*RescheduleOperationConstraint.GetRescheduleOperationConstraint() *RescheduleOperationConstraint`
- New function `*ScheduleOperation.GetMaintenanceManagementOperation() *MaintenanceManagementOperation`
- New function `*ScheduleOperationConstraint.GetScheduleOperationConstraint() *ScheduleOperationConstraint`
- New function `*SchedulingWindow.GetScheduleOperationConstraint() *ScheduleOperationConstraint`
- New function `NewServiceComponentsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ServiceComponentsClient, error)`
- New function `*ServiceComponentsClient.BegincheckAvailability(context.Context, string, string, *serviceComponentsClientBegincheckAvailabilityOptions) (*runtime.Poller[serviceComponentsClientcheckAvailabilityResponse], error)`
- New function `*Vcf5License.GetVcfLicense() *VcfLicense`
- New function `*VcfLicense.GetVcfLicense() *VcfLicense`
- New function `*VmwareFirewallLicenseProperties.GetLicenseProperties() *LicenseProperties`
- New struct `AvailableWindowForMaintenanceWhileRescheduleOperation`
- New struct `AvailableWindowForMaintenanceWhileScheduleOperation`
- New struct `BlockedDatesConstraintTimeRange`
- New struct `BlockedWhileRescheduleOperation`
- New struct `BlockedWhileScheduleOperation`
- New struct `ImpactedMaintenanceResource`
- New struct `ImpactedMaintenanceResourceError`
- New struct `Label`
- New struct `License`
- New struct `LicenseListResult`
- New struct `Maintenance`
- New struct `MaintenanceFailedCheck`
- New struct `MaintenanceListResult`
- New struct `MaintenanceProperties`
- New struct `MaintenanceReadiness`
- New struct `MaintenanceReadinessRefreshOperation`
- New struct `MaintenanceReschedule`
- New struct `MaintenanceSchedule`
- New struct `MaintenanceState`
- New struct `RescheduleOperation`
- New struct `ScheduleOperation`
- New struct `SchedulingWindow`
- New struct `Vcf5License`
- New struct `VmwareFirewallLicenseProperties`
- New field `VcfLicense` in struct `PrivateCloudProperties`


## 2.1.0 (2025-07-29)
### Features Added

- New enum type `HostKind` with values `HostKindGeneral`, `HostKindSpecialized`
- New enum type `HostMaintenance` with values `HostMaintenanceReplacement`, `HostMaintenanceUpgrade`
- New enum type `HostProvisioningState` with values `HostProvisioningStateCanceled`, `HostProvisioningStateFailed`, `HostProvisioningStateSucceeded`
- New enum type `ProvisionedNetworkProvisioningState` with values `ProvisionedNetworkProvisioningStateCanceled`, `ProvisionedNetworkProvisioningStateFailed`, `ProvisionedNetworkProvisioningStateSucceeded`
- New enum type `ProvisionedNetworkTypes` with values `ProvisionedNetworkTypesEsxManagement`, `ProvisionedNetworkTypesEsxReplication`, `ProvisionedNetworkTypesHcxManagement`, `ProvisionedNetworkTypesHcxUplink`, `ProvisionedNetworkTypesVcenterManagement`, `ProvisionedNetworkTypesVmotion`, `ProvisionedNetworkTypesVsan`
- New enum type `PureStoragePolicyProvisioningState` with values `PureStoragePolicyProvisioningStateCanceled`, `PureStoragePolicyProvisioningStateDeleting`, `PureStoragePolicyProvisioningStateFailed`, `PureStoragePolicyProvisioningStateSucceeded`, `PureStoragePolicyProvisioningStateUpdating`
- New enum type `ResourceSKUResourceType` with values `ResourceSKUResourceTypePrivateClouds`, `ResourceSKUResourceTypePrivateCloudsClusters`
- New enum type `ResourceSKURestrictionsReasonCode` with values `ResourceSKURestrictionsReasonCodeNotAvailableForSubscription`, `ResourceSKURestrictionsReasonCodeQuotaID`
- New enum type `ResourceSKURestrictionsType` with values `ResourceSKURestrictionsTypeLocation`, `ResourceSKURestrictionsTypeZone`
- New function `*ClientFactory.NewHostsClient() *HostsClient`
- New function `*ClientFactory.NewProvisionedNetworksClient() *ProvisionedNetworksClient`
- New function `*ClientFactory.NewPureStoragePoliciesClient() *PureStoragePoliciesClient`
- New function `*ClientFactory.NewSKUsClient() *SKUsClient`
- New function `*GeneralHostProperties.GetHostProperties() *HostProperties`
- New function `*HostProperties.GetHostProperties() *HostProperties`
- New function `NewHostsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*HostsClient, error)`
- New function `*HostsClient.Get(context.Context, string, string, string, string, *HostsClientGetOptions) (HostsClientGetResponse, error)`
- New function `*HostsClient.NewListPager(string, string, string, *HostsClientListOptions) *runtime.Pager[HostsClientListResponse]`
- New function `NewProvisionedNetworksClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProvisionedNetworksClient, error)`
- New function `*ProvisionedNetworksClient.Get(context.Context, string, string, string, *ProvisionedNetworksClientGetOptions) (ProvisionedNetworksClientGetResponse, error)`
- New function `*ProvisionedNetworksClient.NewListPager(string, string, *ProvisionedNetworksClientListOptions) *runtime.Pager[ProvisionedNetworksClientListResponse]`
- New function `NewPureStoragePoliciesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PureStoragePoliciesClient, error)`
- New function `*PureStoragePoliciesClient.BeginCreateOrUpdate(context.Context, string, string, string, PureStoragePolicy, *PureStoragePoliciesClientBeginCreateOrUpdateOptions) (*runtime.Poller[PureStoragePoliciesClientCreateOrUpdateResponse], error)`
- New function `*PureStoragePoliciesClient.BeginDelete(context.Context, string, string, string, *PureStoragePoliciesClientBeginDeleteOptions) (*runtime.Poller[PureStoragePoliciesClientDeleteResponse], error)`
- New function `*PureStoragePoliciesClient.Get(context.Context, string, string, string, *PureStoragePoliciesClientGetOptions) (PureStoragePoliciesClientGetResponse, error)`
- New function `*PureStoragePoliciesClient.NewListPager(string, string, *PureStoragePoliciesClientListOptions) *runtime.Pager[PureStoragePoliciesClientListResponse]`
- New function `NewSKUsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SKUsClient, error)`
- New function `*SKUsClient.NewListPager(*SKUsClientListOptions) *runtime.Pager[SKUsClientListResponse]`
- New function `*SpecializedHostProperties.GetHostProperties() *HostProperties`
- New struct `GeneralHostProperties`
- New struct `Host`
- New struct `HostListResult`
- New struct `PagedResourceSKU`
- New struct `ProvisionedNetwork`
- New struct `ProvisionedNetworkListResult`
- New struct `ProvisionedNetworkProperties`
- New struct `PureStoragePolicy`
- New struct `PureStoragePolicyListResult`
- New struct `PureStoragePolicyProperties`
- New struct `PureStorageVolume`
- New struct `ResourceSKU`
- New struct `ResourceSKUCapabilities`
- New struct `ResourceSKULocationInfo`
- New struct `ResourceSKURestrictionInfo`
- New struct `ResourceSKURestrictions`
- New struct `ResourceSKUZoneDetails`
- New struct `SpecializedHostProperties`
- New field `ManagementNetwork`, `UplinkNetwork` in struct `AddonHcxProperties`
- New field `PureStorageVolume` in struct `DatastoreProperties`
- New field `Zones` in struct `PrivateCloud`


## 2.0.0 (2024-09-26)
### Breaking Changes

- Function `*WorkloadNetworksClient.Get` parameter(s) have been changed from `(context.Context, string, string, WorkloadNetworkName, *WorkloadNetworksClientGetOptions)` to `(context.Context, string, string, *WorkloadNetworksClientGetOptions)`
- Type of `Operation.Origin` has been changed from `*string` to `*Origin`
- Enum `WorkloadNetworkName` has been removed
- Struct `LogSpecification` has been removed
- Struct `MetricDimension` has been removed
- Struct `MetricSpecification` has been removed
- Struct `OperationList` has been removed
- Struct `OperationProperties` has been removed
- Struct `ServiceSpecification` has been removed
- Field `Properties` of struct `Operation` has been removed
- Field `OperationList` of struct `OperationsClientListResponse` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `CloudLinkProvisioningState` with values `CloudLinkProvisioningStateCanceled`, `CloudLinkProvisioningStateFailed`, `CloudLinkProvisioningStateSucceeded`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `DNSZoneType` with values `DNSZoneTypePrivate`, `DNSZoneTypePublic`
- New enum type `HcxEnterpriseSiteProvisioningState` with values `HcxEnterpriseSiteProvisioningStateCanceled`, `HcxEnterpriseSiteProvisioningStateFailed`, `HcxEnterpriseSiteProvisioningStateSucceeded`
- New enum type `IscsiPathProvisioningState` with values `IscsiPathProvisioningStateBuilding`, `IscsiPathProvisioningStateCanceled`, `IscsiPathProvisioningStateDeleting`, `IscsiPathProvisioningStateFailed`, `IscsiPathProvisioningStatePending`, `IscsiPathProvisioningStateSucceeded`, `IscsiPathProvisioningStateUpdating`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `SKUTier` with values `SKUTierBasic`, `SKUTierFree`, `SKUTierPremium`, `SKUTierStandard`
- New enum type `ScriptCmdletAudience` with values `ScriptCmdletAudienceAny`, `ScriptCmdletAudienceAutomation`
- New enum type `ScriptCmdletProvisioningState` with values `ScriptCmdletProvisioningStateCanceled`, `ScriptCmdletProvisioningStateFailed`, `ScriptCmdletProvisioningStateSucceeded`
- New enum type `ScriptPackageProvisioningState` with values `ScriptPackageProvisioningStateCanceled`, `ScriptPackageProvisioningStateFailed`, `ScriptPackageProvisioningStateSucceeded`
- New enum type `VirtualMachineProvisioningState` with values `VirtualMachineProvisioningStateCanceled`, `VirtualMachineProvisioningStateFailed`, `VirtualMachineProvisioningStateSucceeded`
- New enum type `WorkloadNetworkProvisioningState` with values `WorkloadNetworkProvisioningStateBuilding`, `WorkloadNetworkProvisioningStateCanceled`, `WorkloadNetworkProvisioningStateDeleting`, `WorkloadNetworkProvisioningStateFailed`, `WorkloadNetworkProvisioningStateSucceeded`, `WorkloadNetworkProvisioningStateUpdating`
- New function `*ClientFactory.NewIscsiPathsClient() *IscsiPathsClient`
- New function `NewIscsiPathsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*IscsiPathsClient, error)`
- New function `*IscsiPathsClient.BeginCreateOrUpdate(context.Context, string, string, IscsiPath, *IscsiPathsClientBeginCreateOrUpdateOptions) (*runtime.Poller[IscsiPathsClientCreateOrUpdateResponse], error)`
- New function `*IscsiPathsClient.BeginDelete(context.Context, string, string, *IscsiPathsClientBeginDeleteOptions) (*runtime.Poller[IscsiPathsClientDeleteResponse], error)`
- New function `*IscsiPathsClient.Get(context.Context, string, string, *IscsiPathsClientGetOptions) (IscsiPathsClientGetResponse, error)`
- New function `*IscsiPathsClient.NewListByPrivateCloudPager(string, string, *IscsiPathsClientListByPrivateCloudOptions) *runtime.Pager[IscsiPathsClientListByPrivateCloudResponse]`
- New struct `ElasticSanVolume`
- New struct `IscsiPath`
- New struct `IscsiPathListResult`
- New struct `IscsiPathProperties`
- New struct `OperationListResult`
- New struct `SystemData`
- New struct `WorkloadNetworkProperties`
- New field `SystemData` in struct `Addon`
- New field `SystemData` in struct `CloudLink`
- New field `ProvisioningState` in struct `CloudLinkProperties`
- New field `SystemData` in struct `Cluster`
- New field `VsanDatastoreName` in struct `ClusterProperties`
- New field `SKU` in struct `ClusterUpdate`
- New field `SystemData` in struct `Datastore`
- New field `ElasticSanVolume` in struct `DatastoreProperties`
- New field `HcxCloudManagerIP`, `NsxtManagerIP`, `VcenterIP` in struct `Endpoints`
- New field `SystemData` in struct `ExpressRouteAuthorization`
- New field `SystemData` in struct `GlobalReachConnection`
- New field `SystemData` in struct `HcxEnterpriseSite`
- New field `ProvisioningState` in struct `HcxEnterpriseSiteProperties`
- New field `VsanDatastoreName` in struct `ManagementCluster`
- New field `ActionType` in struct `Operation`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New field `SystemData` in struct `PlacementPolicy`
- New field `SystemData` in struct `PrivateCloud`
- New field `DNSZoneType`, `VirtualNetworkID` in struct `PrivateCloudProperties`
- New field `SKU` in struct `PrivateCloudUpdate`
- New field `DNSZoneType` in struct `PrivateCloudUpdateProperties`
- New field `Capacity`, `Family`, `Size`, `Tier` in struct `SKU`
- New field `SystemData` in struct `ScriptCmdlet`
- New field `Audience`, `ProvisioningState` in struct `ScriptCmdletProperties`
- New field `SystemData` in struct `ScriptExecution`
- New field `SystemData` in struct `ScriptPackage`
- New field `ProvisioningState` in struct `ScriptPackageProperties`
- New field `SystemData` in struct `VirtualMachine`
- New field `ProvisioningState` in struct `VirtualMachineProperties`
- New field `Properties`, `SystemData` in struct `WorkloadNetwork`
- New field `SystemData` in struct `WorkloadNetworkDNSService`
- New field `SystemData` in struct `WorkloadNetworkDNSZone`
- New field `SystemData` in struct `WorkloadNetworkDhcp`
- New field `SystemData` in struct `WorkloadNetworkGateway`
- New field `ProvisioningState` in struct `WorkloadNetworkGatewayProperties`
- New field `SystemData` in struct `WorkloadNetworkPortMirroring`
- New field `SystemData` in struct `WorkloadNetworkPublicIP`
- New field `SystemData` in struct `WorkloadNetworkSegment`
- New field `SystemData` in struct `WorkloadNetworkVMGroup`
- New field `SystemData` in struct `WorkloadNetworkVirtualMachine`
- New field `ProvisioningState` in struct `WorkloadNetworkVirtualMachineProperties`


## 2.0.0-beta.1 (2024-06-28)
### Breaking Changes

- Function `*WorkloadNetworksClient.Get` parameter(s) have been changed from `(context.Context, string, string, WorkloadNetworkName, *WorkloadNetworksClientGetOptions)` to `(context.Context, string, string, *WorkloadNetworksClientGetOptions)`
- Type of `Operation.Origin` has been changed from `*string` to `*Origin`
- Enum `WorkloadNetworkName` has been removed
- Struct `LogSpecification` has been removed
- Struct `MetricDimension` has been removed
- Struct `MetricSpecification` has been removed
- Struct `OperationList` has been removed
- Struct `OperationProperties` has been removed
- Struct `ServiceSpecification` has been removed
- Field `Properties` of struct `Operation` has been removed
- Field `OperationList` of struct `OperationsClientListResponse` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `CloudLinkProvisioningState` with values `CloudLinkProvisioningStateCanceled`, `CloudLinkProvisioningStateFailed`, `CloudLinkProvisioningStateSucceeded`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `DNSZoneType` with values `DNSZoneTypePrivate`, `DNSZoneTypePublic`
- New enum type `HcxEnterpriseSiteProvisioningState` with values `HcxEnterpriseSiteProvisioningStateCanceled`, `HcxEnterpriseSiteProvisioningStateFailed`, `HcxEnterpriseSiteProvisioningStateSucceeded`
- New enum type `IscsiPathProvisioningState` with values `IscsiPathProvisioningStateBuilding`, `IscsiPathProvisioningStateCanceled`, `IscsiPathProvisioningStateDeleting`, `IscsiPathProvisioningStateFailed`, `IscsiPathProvisioningStatePending`, `IscsiPathProvisioningStateSucceeded`, `IscsiPathProvisioningStateUpdating`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `SKUTier` with values `SKUTierBasic`, `SKUTierFree`, `SKUTierPremium`, `SKUTierStandard`
- New enum type `ScriptCmdletAudience` with values `ScriptCmdletAudienceAny`, `ScriptCmdletAudienceAutomation`
- New enum type `ScriptCmdletProvisioningState` with values `ScriptCmdletProvisioningStateCanceled`, `ScriptCmdletProvisioningStateFailed`, `ScriptCmdletProvisioningStateSucceeded`
- New enum type `ScriptPackageProvisioningState` with values `ScriptPackageProvisioningStateCanceled`, `ScriptPackageProvisioningStateFailed`, `ScriptPackageProvisioningStateSucceeded`
- New enum type `VirtualMachineProvisioningState` with values `VirtualMachineProvisioningStateCanceled`, `VirtualMachineProvisioningStateFailed`, `VirtualMachineProvisioningStateSucceeded`
- New enum type `WorkloadNetworkProvisioningState` with values `WorkloadNetworkProvisioningStateBuilding`, `WorkloadNetworkProvisioningStateCanceled`, `WorkloadNetworkProvisioningStateDeleting`, `WorkloadNetworkProvisioningStateFailed`, `WorkloadNetworkProvisioningStateSucceeded`, `WorkloadNetworkProvisioningStateUpdating`
- New function `*ClientFactory.NewIscsiPathsClient() *IscsiPathsClient`
- New function `NewIscsiPathsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*IscsiPathsClient, error)`
- New function `*IscsiPathsClient.BeginCreateOrUpdate(context.Context, string, string, IscsiPath, *IscsiPathsClientBeginCreateOrUpdateOptions) (*runtime.Poller[IscsiPathsClientCreateOrUpdateResponse], error)`
- New function `*IscsiPathsClient.BeginDelete(context.Context, string, string, *IscsiPathsClientBeginDeleteOptions) (*runtime.Poller[IscsiPathsClientDeleteResponse], error)`
- New function `*IscsiPathsClient.Get(context.Context, string, string, *IscsiPathsClientGetOptions) (IscsiPathsClientGetResponse, error)`
- New function `*IscsiPathsClient.NewListByPrivateCloudPager(string, string, *IscsiPathsClientListByPrivateCloudOptions) *runtime.Pager[IscsiPathsClientListByPrivateCloudResponse]`
- New struct `ElasticSanVolume`
- New struct `IscsiPath`
- New struct `IscsiPathListResult`
- New struct `IscsiPathProperties`
- New struct `OperationListResult`
- New struct `SystemData`
- New struct `WorkloadNetworkProperties`
- New field `SystemData` in struct `Addon`
- New field `SystemData` in struct `CloudLink`
- New field `ProvisioningState` in struct `CloudLinkProperties`
- New field `SystemData` in struct `Cluster`
- New field `VsanDatastoreName` in struct `ClusterProperties`
- New field `SKU` in struct `ClusterUpdate`
- New field `SystemData` in struct `Datastore`
- New field `ElasticSanVolume` in struct `DatastoreProperties`
- New field `HcxCloudManagerIP`, `NsxtManagerIP`, `VcenterIP` in struct `Endpoints`
- New field `SystemData` in struct `ExpressRouteAuthorization`
- New field `SystemData` in struct `GlobalReachConnection`
- New field `SystemData` in struct `HcxEnterpriseSite`
- New field `ProvisioningState` in struct `HcxEnterpriseSiteProperties`
- New field `VsanDatastoreName` in struct `ManagementCluster`
- New field `ActionType` in struct `Operation`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New field `SystemData` in struct `PlacementPolicy`
- New field `SystemData` in struct `PrivateCloud`
- New field `DNSZoneType`, `VirtualNetworkID` in struct `PrivateCloudProperties`
- New field `SKU` in struct `PrivateCloudUpdate`
- New field `DNSZoneType` in struct `PrivateCloudUpdateProperties`
- New field `Capacity`, `Family`, `Size`, `Tier` in struct `SKU`
- New field `SystemData` in struct `ScriptCmdlet`
- New field `Audience`, `ProvisioningState` in struct `ScriptCmdletProperties`
- New field `SystemData` in struct `ScriptExecution`
- New field `SystemData` in struct `ScriptPackage`
- New field `ProvisioningState` in struct `ScriptPackageProperties`
- New field `SystemData` in struct `VirtualMachine`
- New field `ProvisioningState` in struct `VirtualMachineProperties`
- New field `Properties`, `SystemData` in struct `WorkloadNetwork`
- New field `SystemData` in struct `WorkloadNetworkDNSService`
- New field `SystemData` in struct `WorkloadNetworkDNSZone`
- New field `SystemData` in struct `WorkloadNetworkDhcp`
- New field `SystemData` in struct `WorkloadNetworkGateway`
- New field `ProvisioningState` in struct `WorkloadNetworkGatewayProperties`
- New field `SystemData` in struct `WorkloadNetworkPortMirroring`
- New field `SystemData` in struct `WorkloadNetworkPublicIP`
- New field `SystemData` in struct `WorkloadNetworkSegment`
- New field `SystemData` in struct `WorkloadNetworkVMGroup`
- New field `SystemData` in struct `WorkloadNetworkVirtualMachine`
- New field `ProvisioningState` in struct `WorkloadNetworkVirtualMachineProperties`


## 1.4.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.3.0 (2023-08-25)

### Features Added

- New field `ExtendedNetworkBlocks` in struct `PrivateCloudProperties`
- New field `ExtendedNetworkBlocks` in struct `PrivateCloudUpdateProperties`


## 1.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.2.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.1.0 (2022-10-13)

### Features Added

- New const `ExpressRouteAuthorizationProvisioningStateCanceled`
- New const `AffinityStrengthShould`
- New const `AddonTypeArc`
- New const `PrivateCloudProvisioningStateCanceled`
- New const `NsxPublicIPQuotaRaisedEnumEnabled`
- New const `AzureHybridBenefitTypeNone`
- New const `WorkloadNetworkPublicIPProvisioningStateCanceled`
- New const `WorkloadNetworkDNSServiceProvisioningStateCanceled`
- New const `WorkloadNetworkSegmentProvisioningStateCanceled`
- New const `WorkloadNetworkDNSZoneProvisioningStateCanceled`
- New const `WorkloadNetworkNameDefault`
- New const `PlacementPolicyProvisioningStateCanceled`
- New const `WorkloadNetworkDhcpProvisioningStateCanceled`
- New const `WorkloadNetworkPortMirroringProvisioningStateCanceled`
- New const `WorkloadNetworkVMGroupProvisioningStateCanceled`
- New const `NsxPublicIPQuotaRaisedEnumDisabled`
- New const `DatastoreProvisioningStateCanceled`
- New const `AzureHybridBenefitTypeSQLHost`
- New const `AddonProvisioningStateCanceled`
- New const `ClusterProvisioningStateCanceled`
- New const `AffinityStrengthMust`
- New const `GlobalReachConnectionProvisioningStateCanceled`
- New const `ScriptExecutionProvisioningStateCanceled`
- New type alias `NsxPublicIPQuotaRaisedEnum`
- New type alias `AzureHybridBenefitType`
- New type alias `AffinityStrength`
- New type alias `WorkloadNetworkName`
- New function `PossibleAzureHybridBenefitTypeValues() []AzureHybridBenefitType`
- New function `*WorkloadNetworksClient.Get(context.Context, string, string, WorkloadNetworkName, *WorkloadNetworksClientGetOptions) (WorkloadNetworksClientGetResponse, error)`
- New function `*ClustersClient.ListZones(context.Context, string, string, string, *ClustersClientListZonesOptions) (ClustersClientListZonesResponse, error)`
- New function `PossibleNsxPublicIPQuotaRaisedEnumValues() []NsxPublicIPQuotaRaisedEnum`
- New function `PossibleWorkloadNetworkNameValues() []WorkloadNetworkName`
- New function `*WorkloadNetworksClient.NewListPager(string, string, *WorkloadNetworksClientListOptions) *runtime.Pager[WorkloadNetworksClientListResponse]`
- New function `*AddonArcProperties.GetAddonProperties() *AddonProperties`
- New function `PossibleAffinityStrengthValues() []AffinityStrength`
- New struct `AddonArcProperties`
- New struct `ClusterZone`
- New struct `ClusterZoneList`
- New struct `ClustersClientListZonesOptions`
- New struct `ClustersClientListZonesResponse`
- New struct `WorkloadNetwork`
- New struct `WorkloadNetworkList`
- New struct `WorkloadNetworksClientGetOptions`
- New struct `WorkloadNetworksClientGetResponse`
- New struct `WorkloadNetworksClientListOptions`
- New struct `WorkloadNetworksClientListResponse`
- New field `AffinityStrength` in struct `PlacementPolicyUpdateProperties`
- New field `AzureHybridBenefitType` in struct `PlacementPolicyUpdateProperties`
- New field `AutoDetectedKeyVersion` in struct `EncryptionKeyVaultProperties`
- New field `SKU` in struct `LocationsClientCheckTrialAvailabilityOptions`
- New field `AzureHybridBenefitType` in struct `VMHostPlacementPolicyProperties`
- New field `AffinityStrength` in struct `VMHostPlacementPolicyProperties`
- New field `NsxPublicIPQuotaRaised` in struct `PrivateCloudProperties`
- New field `Company` in struct `ScriptPackageProperties`
- New field `URI` in struct `ScriptPackageProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/avs/armavs` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).