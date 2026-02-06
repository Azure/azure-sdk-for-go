# Release History

## 2.5.0 (2026-02-06)
### Features Added

- New value `AgentUpgradeBlockedReasonReInstallRequired` added to enum type `AgentUpgradeBlockedReason`
- New enum type `AgentReinstallBlockedReason` with values `AgentReinstallBlockedReasonAgentNoHeartbeat`, `AgentReinstallBlockedReasonDistroNotSupported`, `AgentReinstallBlockedReasonUnknown`
- New enum type `MobilityAgentReinstallType` with values `MobilityAgentReinstallTypeAutoTriggered`, `MobilityAgentReinstallTypeUserTriggered`
- New function `*ReplicationProtectedItemsClient.BeginReinstallMobilityService(ctx context.Context, resourceGroupName string, resourceName string, fabricName string, protectionContainerName string, replicatedProtectedItemName string, updateMobilityServiceRequest ReinstallMobilityServiceRequest, options *ReplicationProtectedItemsClientBeginReinstallMobilityServiceOptions) (*runtime.Poller[ReplicationProtectedItemsClientReinstallMobilityServiceResponse], error)`
- New struct `A2AAgentReinstallBlockingErrorDetails`
- New struct `InMageRcmAgentReinstallBlockingErrorDetails`
- New struct `ReinstallMobilityServiceRequest`
- New struct `ReinstallMobilityServiceRequestProperties`
- New field `PlatformFaultDomain` in struct `A2AEnableProtectionInput`
- New field `AgentReinstallAttemptToVersion`, `AutoAgentUpgradeRetryCount`, `DistroName`, `DistroNameForWhichAgentIsInstalled`, `IsAgentReinstallRequired`, `IsAgentUpgradeInProgress`, `IsAgentUpgradeRetryThresholdExhausted`, `IsAgentUpgradeable`, `OSFamilyName`, `PlatformFaultDomain`, `ReasonsBlockingReInstall`, `ReasonsBlockingReinstallDetails` in struct `A2AReplicationDetails`
- New field `PlatformFaultDomain` in struct `A2ASwitchProtectionInput`
- New field `PlatformFaultDomain`, `RecoveryAvailabilityZone` in struct `A2AUpdateReplicationProtectedItemInput`
- New field `DiskSizeInGB`, `Iops`, `ThroughputInMbps` in struct `HyperVReplicaAzureDiskInputDetails`
- New field `TargetCapacityReservationGroupID` in struct `HyperVReplicaAzureEnableProtectionInput`
- New field `DiskSizeInGB`, `Iops`, `ThroughputInMbps` in struct `HyperVReplicaAzureManagedDiskDetails`
- New field `TargetCapacityReservationGroupID` in struct `HyperVReplicaAzurePlannedFailoverProviderInput`
- New field `TargetCapacityReservationGroupID` in struct `HyperVReplicaAzureReplicationDetails`
- New field `TargetCapacityReservationGroupID` in struct `HyperVReplicaAzureUpdateReplicationProtectedItemInput`
- New field `DiskSizeInGB`, `Iops`, `ThroughputInMbps` in struct `InMageRcmDiskInput`
- New field `DiskSizeInGB`, `Iops`, `ThroughputInMbps` in struct `InMageRcmDisksDefaultInput`
- New field `TargetCapacityReservationGroupID` in struct `InMageRcmEnableProtectionInput`
- New field `AgentReinstallAttemptToVersion`, `AgentReinstallJobID`, `AgentReinstallState`, `DistroName`, `DistroNameForWhichAgentIsInstalled`, `IsAgentReinstallRequired`, `IsAgentUpgradeable`, `IsLastReinstallSuccessful`, `LastAgentReinstallType`, `OSFamilyName`, `ReasonsBlockingReinstall`, `ReasonsBlockingReinstallDetails` in struct `InMageRcmMobilityAgentDetails`
- New field `DiskSizeInGB`, `Iops`, `ThroughputInMbps` in struct `InMageRcmProtectedDiskDetails`
- New field `TargetCapacityReservationGroupID` in struct `InMageRcmReplicationDetails`
- New field `TargetCapacityReservationGroupID` in struct `InMageRcmUnplannedFailoverInput`
- New field `TargetCapacityReservationGroupID`, `VMDisks` in struct `InMageRcmUpdateReplicationProtectedItemInput`
- New field `DiskSizeInGB`, `Iops`, `ThroughputInMbps` in struct `UpdateDiskInput`
- New field `DiskSizeInGB`, `Iops`, `ThroughputInMbps` in struct `VMwareCbtDiskInput`
- New field `TargetCapacityReservationGroupID` in struct `VMwareCbtEnableMigrationInput`
- New field `TargetCapacityReservationGroupID` in struct `VMwareCbtMigrateInput`
- New field `TargetCapacityReservationGroupID` in struct `VMwareCbtMigrationDetails`
- New field `DiskSizeInGB`, `Iops`, `ThroughputInMbps` in struct `VMwareCbtProtectedDiskDetails`
- New field `DiskSizeInGB`, `Iops`, `ThroughputInMbps` in struct `VMwareCbtUpdateDiskInput`
- New field `TargetCapacityReservationGroupID` in struct `VMwareCbtUpdateMigrationItemInput`


## 2.4.0 (2025-04-25)
### Features Added

- New value `DiskAccountTypePremiumV2LRS`, `DiskAccountTypePremiumZRS`, `DiskAccountTypeStandardSSDZRS`, `DiskAccountTypeUltraSSDLRS` added to enum type `DiskAccountType`
- New enum type `ClusterRecoveryPointType` with values `ClusterRecoveryPointTypeApplicationConsistent`, `ClusterRecoveryPointTypeCrashConsistent`, `ClusterRecoveryPointTypeNotSpecified`
- New enum type `DiskState` with values `DiskStateInitialReplicationFailed`, `DiskStateInitialReplicationPending`, `DiskStateProtected`, `DiskStateUnavailable`
- New enum type `FailoverDirection` with values `FailoverDirectionPrimaryToRecovery`, `FailoverDirectionRecoveryToPrimary`
- New enum type `LinuxLicenseType` with values `LinuxLicenseTypeLinuxServer`, `LinuxLicenseTypeNoLicenseType`, `LinuxLicenseTypeNotSpecified`
- New enum type `SecurityConfiguration` with values `SecurityConfigurationDisabled`, `SecurityConfigurationEnabled`
- New function `*A2AApplyClusterRecoveryPointInput.GetApplyClusterRecoveryPointProviderSpecificInput() *ApplyClusterRecoveryPointProviderSpecificInput`
- New function `*A2AClusterRecoveryPointDetails.GetClusterProviderSpecificRecoveryPointDetails() *ClusterProviderSpecificRecoveryPointDetails`
- New function `*A2AClusterTestFailoverInput.GetClusterTestFailoverProviderSpecificInput() *ClusterTestFailoverProviderSpecificInput`
- New function `*A2AClusterUnplannedFailoverInput.GetClusterUnplannedFailoverProviderSpecificInput() *ClusterUnplannedFailoverProviderSpecificInput`
- New function `*A2AReplicationProtectionClusterDetails.GetReplicationClusterProviderSpecificSettings() *ReplicationClusterProviderSpecificSettings`
- New function `*A2ASharedDiskReplicationDetails.GetSharedDiskReplicationProviderSpecificSettings() *SharedDiskReplicationProviderSpecificSettings`
- New function `*A2ASwitchClusterProtectionInput.GetSwitchClusterProtectionProviderSpecificInput() *SwitchClusterProtectionProviderSpecificInput`
- New function `*ApplyClusterRecoveryPointProviderSpecificInput.GetApplyClusterRecoveryPointProviderSpecificInput() *ApplyClusterRecoveryPointProviderSpecificInput`
- New function `*ClientFactory.NewClusterRecoveryPointClient() *ClusterRecoveryPointClient`
- New function `*ClientFactory.NewClusterRecoveryPointsClient() *ClusterRecoveryPointsClient`
- New function `*ClientFactory.NewReplicationProtectionClustersClient() *ReplicationProtectionClustersClient`
- New function `*ClusterFailoverJobDetails.GetJobDetails() *JobDetails`
- New function `*ClusterProviderSpecificRecoveryPointDetails.GetClusterProviderSpecificRecoveryPointDetails() *ClusterProviderSpecificRecoveryPointDetails`
- New function `NewClusterRecoveryPointClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ClusterRecoveryPointClient, error)`
- New function `*ClusterRecoveryPointClient.Get(context.Context, string, string, string, string, string, string, *ClusterRecoveryPointClientGetOptions) (ClusterRecoveryPointClientGetResponse, error)`
- New function `NewClusterRecoveryPointsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ClusterRecoveryPointsClient, error)`
- New function `*ClusterRecoveryPointsClient.NewListByReplicationProtectionClusterPager(string, string, string, string, string, *ClusterRecoveryPointsClientListByReplicationProtectionClusterOptions) *runtime.Pager[ClusterRecoveryPointsClientListByReplicationProtectionClusterResponse]`
- New function `*ClusterSwitchProtectionJobDetails.GetJobDetails() *JobDetails`
- New function `*ClusterTestFailoverJobDetails.GetJobDetails() *JobDetails`
- New function `*ClusterTestFailoverProviderSpecificInput.GetClusterTestFailoverProviderSpecificInput() *ClusterTestFailoverProviderSpecificInput`
- New function `*ClusterUnplannedFailoverProviderSpecificInput.GetClusterUnplannedFailoverProviderSpecificInput() *ClusterUnplannedFailoverProviderSpecificInput`
- New function `*InMageRcmAddDisksInput.GetAddDisksProviderSpecificInput() *AddDisksProviderSpecificInput`
- New function `*ReplicationClusterProviderSpecificSettings.GetReplicationClusterProviderSpecificSettings() *ReplicationClusterProviderSpecificSettings`
- New function `NewReplicationProtectionClustersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ReplicationProtectionClustersClient, error)`
- New function `*ReplicationProtectionClustersClient.BeginApplyRecoveryPoint(context.Context, string, string, string, string, string, ApplyClusterRecoveryPointInput, *ReplicationProtectionClustersClientBeginApplyRecoveryPointOptions) (*runtime.Poller[ReplicationProtectionClustersClientApplyRecoveryPointResponse], error)`
- New function `*ReplicationProtectionClustersClient.BeginCreate(context.Context, string, string, string, string, string, ReplicationProtectionCluster, *ReplicationProtectionClustersClientBeginCreateOptions) (*runtime.Poller[ReplicationProtectionClustersClientCreateResponse], error)`
- New function `*ReplicationProtectionClustersClient.BeginFailoverCommit(context.Context, string, string, string, string, string, *ReplicationProtectionClustersClientBeginFailoverCommitOptions) (*runtime.Poller[ReplicationProtectionClustersClientFailoverCommitResponse], error)`
- New function `*ReplicationProtectionClustersClient.Get(context.Context, string, string, string, string, string, *ReplicationProtectionClustersClientGetOptions) (ReplicationProtectionClustersClientGetResponse, error)`
- New function `*ReplicationProtectionClustersClient.GetOperationResults(context.Context, string, string, string, string, string, string, *ReplicationProtectionClustersClientGetOperationResultsOptions) (ReplicationProtectionClustersClientGetOperationResultsResponse, error)`
- New function `*ReplicationProtectionClustersClient.NewListByReplicationProtectionContainersPager(string, string, string, string, *ReplicationProtectionClustersClientListByReplicationProtectionContainersOptions) *runtime.Pager[ReplicationProtectionClustersClientListByReplicationProtectionContainersResponse]`
- New function `*ReplicationProtectionClustersClient.NewListPager(string, string, *ReplicationProtectionClustersClientListOptions) *runtime.Pager[ReplicationProtectionClustersClientListResponse]`
- New function `*ReplicationProtectionClustersClient.BeginPurge(context.Context, string, string, string, string, string, *ReplicationProtectionClustersClientBeginPurgeOptions) (*runtime.Poller[ReplicationProtectionClustersClientPurgeResponse], error)`
- New function `*ReplicationProtectionClustersClient.BeginRepairReplication(context.Context, string, string, string, string, string, *ReplicationProtectionClustersClientBeginRepairReplicationOptions) (*runtime.Poller[ReplicationProtectionClustersClientRepairReplicationResponse], error)`
- New function `*ReplicationProtectionClustersClient.BeginTestFailover(context.Context, string, string, string, string, string, ClusterTestFailoverInput, *ReplicationProtectionClustersClientBeginTestFailoverOptions) (*runtime.Poller[ReplicationProtectionClustersClientTestFailoverResponse], error)`
- New function `*ReplicationProtectionClustersClient.BeginTestFailoverCleanup(context.Context, string, string, string, string, string, ClusterTestFailoverCleanupInput, *ReplicationProtectionClustersClientBeginTestFailoverCleanupOptions) (*runtime.Poller[ReplicationProtectionClustersClientTestFailoverCleanupResponse], error)`
- New function `*ReplicationProtectionClustersClient.BeginUnplannedFailover(context.Context, string, string, string, string, string, ClusterUnplannedFailoverInput, *ReplicationProtectionClustersClientBeginUnplannedFailoverOptions) (*runtime.Poller[ReplicationProtectionClustersClientUnplannedFailoverResponse], error)`
- New function `*ReplicationProtectionContainersClient.BeginSwitchClusterProtection(context.Context, string, string, string, string, SwitchClusterProtectionInput, *ReplicationProtectionContainersClientBeginSwitchClusterProtectionOptions) (*runtime.Poller[ReplicationProtectionContainersClientSwitchClusterProtectionResponse], error)`
- New function `*SharedDiskReplicationProviderSpecificSettings.GetSharedDiskReplicationProviderSpecificSettings() *SharedDiskReplicationProviderSpecificSettings`
- New function `*SwitchClusterProtectionProviderSpecificInput.GetSwitchClusterProtectionProviderSpecificInput() *SwitchClusterProtectionProviderSpecificInput`
- New struct `A2AApplyClusterRecoveryPointInput`
- New struct `A2AClusterRecoveryPointDetails`
- New struct `A2AClusterTestFailoverInput`
- New struct `A2AClusterUnplannedFailoverInput`
- New struct `A2AProtectedItemDetail`
- New struct `A2AReplicationProtectionClusterDetails`
- New struct `A2ASharedDiskIRErrorDetails`
- New struct `A2ASharedDiskReplicationDetails`
- New struct `A2ASwitchClusterProtectionInput`
- New struct `ApplyClusterRecoveryPointInput`
- New struct `ApplyClusterRecoveryPointInputProperties`
- New struct `ClusterFailoverJobDetails`
- New struct `ClusterRecoveryPoint`
- New struct `ClusterRecoveryPointCollection`
- New struct `ClusterRecoveryPointProperties`
- New struct `ClusterSwitchProtectionJobDetails`
- New struct `ClusterTestFailoverCleanupInput`
- New struct `ClusterTestFailoverCleanupInputProperties`
- New struct `ClusterTestFailoverInput`
- New struct `ClusterTestFailoverInputProperties`
- New struct `ClusterTestFailoverJobDetails`
- New struct `ClusterUnplannedFailoverInput`
- New struct `ClusterUnplannedFailoverInputProperties`
- New struct `InMageRcmAddDisksInput`
- New struct `InMageRcmUnProtectedDiskDetails`
- New struct `ManagedRunCommandScriptInput`
- New struct `ProtectedClustersQueryParameter`
- New struct `RegisteredClusterNodes`
- New struct `ReplicationProtectionCluster`
- New struct `ReplicationProtectionClusterCollection`
- New struct `ReplicationProtectionClusterProperties`
- New struct `SecurityProfileProperties`
- New struct `SharedDiskReplicationItemProperties`
- New struct `SwitchClusterProtectionInput`
- New struct `SwitchClusterProtectionInputProperties`
- New struct `UserCreatedResourceTag`
- New field `ProtectionClusterID` in struct `A2AEnableProtectionInput`
- New field `IsClusterInfraReady`, `ProtectionClusterID` in struct `A2AReplicationDetails`
- New field `SectorSizeInBytes` in struct `HyperVReplicaAzureDiskInputDetails`
- New field `LinuxLicenseType`, `TargetVMSecurityProfile`, `UserSelectedOSName` in struct `HyperVReplicaAzureEnableProtectionInput`
- New field `SectorSizeInBytes`, `TargetDiskAccountType` in struct `HyperVReplicaAzureManagedDiskDetails`
- New field `LinuxLicenseType`, `TargetVMSecurityProfile` in struct `HyperVReplicaAzureReplicationDetails`
- New field `LinuxLicenseType`, `UserSelectedOSName` in struct `HyperVReplicaAzureUpdateReplicationProtectedItemInput`
- New field `SectorSizeInBytes` in struct `InMageRcmDiskInput`
- New field `SectorSizeInBytes` in struct `InMageRcmDisksDefaultInput`
- New field `LinuxLicenseType`, `SQLServerLicenseType`, `SeedManagedDiskTags`, `TargetManagedDiskTags`, `TargetNicTags`, `TargetVMSecurityProfile`, `TargetVMTags`, `UserSelectedOSName` in struct `InMageRcmEnableProtectionInput`
- New field `TargetNicName` in struct `InMageRcmNicDetails`
- New field `TargetNicName` in struct `InMageRcmNicInput`
- New field `CustomTargetDiskName`, `DiskState`, `SectorSizeInBytes` in struct `InMageRcmProtectedDiskDetails`
- New field `LinuxLicenseType`, `OSName`, `SQLServerLicenseType`, `SeedManagedDiskTags`, `SupportedOSVersions`, `TargetManagedDiskTags`, `TargetNicTags`, `TargetVMSecurityProfile`, `TargetVMTags`, `UnprotectedDisks` in struct `InMageRcmReplicationDetails`
- New field `OSUpgradeVersion` in struct `InMageRcmTestFailoverInput`
- New field `OSUpgradeVersion` in struct `InMageRcmUnplannedFailoverInput`
- New field `LinuxLicenseType`, `SQLServerLicenseType`, `TargetManagedDiskTags`, `TargetNicTags`, `TargetVMTags`, `UserSelectedOSName` in struct `InMageRcmUpdateReplicationProtectedItemInput`
- New field `UserSelectedOSName` in struct `OSDetails`
- New field `SectorSizeInBytes` in struct `VMwareCbtDiskInput`
- New field `LinuxLicenseType`, `UserSelectedOSName` in struct `VMwareCbtEnableMigrationInput`
- New field `PostMigrationSteps` in struct `VMwareCbtMigrateInput`
- New field `LinuxLicenseType` in struct `VMwareCbtMigrationDetails`
- New field `SectorSizeInBytes` in struct `VMwareCbtProtectedDiskDetails`
- New field `PostMigrationSteps` in struct `VMwareCbtTestMigrateInput`
- New field `LinuxLicenseType`, `UserSelectedOSName` in struct `VMwareCbtUpdateMigrationItemInput`


## 2.3.0 (2024-02-23)
### Features Added

- New function `*ReplicationFabricsClient.BeginRemoveInfra(context.Context, string, string, string, *ReplicationFabricsClientBeginRemoveInfraOptions) (*runtime.Poller[ReplicationFabricsClientRemoveInfraResponse], error)`
- New field `AutoProtectionOfDataDisk` in struct `A2AEnableProtectionInput`


## 2.2.0 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.1.0 (2023-09-22)
### Features Added

- New enum type `ChurnOptionSelected` with values `ChurnOptionSelectedHigh`, `ChurnOptionSelectedNormal`
- New enum type `SecurityType` with values `SecurityTypeConfidentialVM`, `SecurityTypeNone`, `SecurityTypeTrustedLaunch`
- New struct `A2AFabricSpecificLocationDetails`
- New struct `ApplianceMonitoringDetails`
- New struct `ApplianceResourceDetails`
- New struct `DataStoreUtilizationDetails`
- New struct `GatewayOperationDetails`
- New struct `OSUpgradeSupportedVersions`
- New struct `VMwareCbtSecurityProfileProperties`
- New field `ChurnOptionSelected` in struct `A2AReplicationDetails`
- New field `LocationDetails` in struct `AzureFabricSpecificDetails`
- New field `ExtendedLocationMappings`, `LocationDetails` in struct `FabricQueryParameter`
- New field `OSUpgradeVersion` in struct `HyperVReplicaAzurePlannedFailoverProviderInput`
- New field `AllAvailableOSUpgradeConfigurations` in struct `HyperVReplicaAzureReplicationDetails`
- New field `OSUpgradeVersion` in struct `HyperVReplicaAzureTestFailoverInput`
- New field `AllAvailableOSUpgradeConfigurations`, `OSName`, `SupportedOSVersions` in struct `InMageAzureV2ReplicationDetails`
- New field `OSUpgradeVersion` in struct `InMageAzureV2TestFailoverInput`
- New field `OSUpgradeVersion` in struct `InMageAzureV2UnplannedFailoverInput`
- New field `ConfidentialVMKeyVaultID`, `TargetVMSecurityProfile` in struct `VMwareCbtEnableMigrationInput`
- New field `OSUpgradeVersion` in struct `VMwareCbtMigrateInput`
- New field `ApplianceMonitoringDetails`, `ConfidentialVMKeyVaultID`, `DeltaSyncProgressPercentage`, `DeltaSyncRetryCount`, `GatewayOperationDetails`, `IsCheckSumResyncCycle`, `OSName`, `OperationName`, `SupportedOSVersions`, `TargetVMSecurityProfile` in struct `VMwareCbtMigrationDetails`
- New field `GatewayOperationDetails` in struct `VMwareCbtProtectedDiskDetails`
- New field `ExcludedSKUs` in struct `VMwareCbtProtectionContainerMappingDetails`
- New field `OSUpgradeVersion` in struct `VMwareCbtTestMigrateInput`


## 2.0.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 2.0.0 (2023-04-03)
### Breaking Changes

- Function `NewMigrationRecoveryPointsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*MigrationRecoveryPointsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, string, *MigrationRecoveryPointsClientGetOptions)` to `(context.Context, string, string, string, string, string, string, *MigrationRecoveryPointsClientGetOptions)`
- Function `*MigrationRecoveryPointsClient.NewListByReplicationMigrationItemsPager` parameter(s) have been changed from `(string, string, string, *MigrationRecoveryPointsClientListByReplicationMigrationItemsOptions)` to `(string, string, string, string, string, *MigrationRecoveryPointsClientListByReplicationMigrationItemsOptions)`
- Function `NewOperationsClient` parameter(s) have been changed from `(string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*OperationsClient.NewListPager` parameter(s) have been changed from `(*OperationsClientListOptions)` to `(string, *OperationsClientListOptions)`
- Function `NewRecoveryPointsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*RecoveryPointsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, string, *RecoveryPointsClientGetOptions)` to `(context.Context, string, string, string, string, string, string, *RecoveryPointsClientGetOptions)`
- Function `*RecoveryPointsClient.NewListByReplicationProtectedItemsPager` parameter(s) have been changed from `(string, string, string, *RecoveryPointsClientListByReplicationProtectedItemsOptions)` to `(string, string, string, string, string, *RecoveryPointsClientListByReplicationProtectedItemsOptions)`
- Function `NewReplicationAlertSettingsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationAlertSettingsClient.Create` parameter(s) have been changed from `(context.Context, string, ConfigureAlertRequest, *ReplicationAlertSettingsClientCreateOptions)` to `(context.Context, string, string, string, ConfigureAlertRequest, *ReplicationAlertSettingsClientCreateOptions)`
- Function `*ReplicationAlertSettingsClient.Get` parameter(s) have been changed from `(context.Context, string, *ReplicationAlertSettingsClientGetOptions)` to `(context.Context, string, string, string, *ReplicationAlertSettingsClientGetOptions)`
- Function `*ReplicationAlertSettingsClient.NewListPager` parameter(s) have been changed from `(*ReplicationAlertSettingsClientListOptions)` to `(string, string, *ReplicationAlertSettingsClientListOptions)`
- Function `NewReplicationAppliancesClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationAppliancesClient.NewListPager` parameter(s) have been changed from `(*ReplicationAppliancesClientListOptions)` to `(string, string, *ReplicationAppliancesClientListOptions)`
- Function `NewReplicationEligibilityResultsClient` parameter(s) have been changed from `(string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationEligibilityResultsClient.Get` parameter(s) have been changed from `(context.Context, string, *ReplicationEligibilityResultsClientGetOptions)` to `(context.Context, string, string, *ReplicationEligibilityResultsClientGetOptions)`
- Function `*ReplicationEligibilityResultsClient.List` parameter(s) have been changed from `(context.Context, string, *ReplicationEligibilityResultsClientListOptions)` to `(context.Context, string, string, *ReplicationEligibilityResultsClientListOptions)`
- Function `NewReplicationEventsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationEventsClient.Get` parameter(s) have been changed from `(context.Context, string, *ReplicationEventsClientGetOptions)` to `(context.Context, string, string, string, *ReplicationEventsClientGetOptions)`
- Function `*ReplicationEventsClient.NewListPager` parameter(s) have been changed from `(*ReplicationEventsClientListOptions)` to `(string, string, *ReplicationEventsClientListOptions)`
- Function `NewReplicationFabricsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationFabricsClient.BeginCheckConsistency` parameter(s) have been changed from `(context.Context, string, *ReplicationFabricsClientBeginCheckConsistencyOptions)` to `(context.Context, string, string, string, *ReplicationFabricsClientBeginCheckConsistencyOptions)`
- Function `*ReplicationFabricsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, FabricCreationInput, *ReplicationFabricsClientBeginCreateOptions)` to `(context.Context, string, string, string, FabricCreationInput, *ReplicationFabricsClientBeginCreateOptions)`
- Function `*ReplicationFabricsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, *ReplicationFabricsClientBeginDeleteOptions)` to `(context.Context, string, string, string, *ReplicationFabricsClientBeginDeleteOptions)`
- Function `*ReplicationFabricsClient.BeginMigrateToAAD` parameter(s) have been changed from `(context.Context, string, *ReplicationFabricsClientBeginMigrateToAADOptions)` to `(context.Context, string, string, string, *ReplicationFabricsClientBeginMigrateToAADOptions)`
- Function `*ReplicationFabricsClient.BeginPurge` parameter(s) have been changed from `(context.Context, string, *ReplicationFabricsClientBeginPurgeOptions)` to `(context.Context, string, string, string, *ReplicationFabricsClientBeginPurgeOptions)`
- Function `*ReplicationFabricsClient.BeginReassociateGateway` parameter(s) have been changed from `(context.Context, string, FailoverProcessServerRequest, *ReplicationFabricsClientBeginReassociateGatewayOptions)` to `(context.Context, string, string, string, FailoverProcessServerRequest, *ReplicationFabricsClientBeginReassociateGatewayOptions)`
- Function `*ReplicationFabricsClient.BeginRenewCertificate` parameter(s) have been changed from `(context.Context, string, RenewCertificateInput, *ReplicationFabricsClientBeginRenewCertificateOptions)` to `(context.Context, string, string, string, RenewCertificateInput, *ReplicationFabricsClientBeginRenewCertificateOptions)`
- Function `*ReplicationFabricsClient.Get` parameter(s) have been changed from `(context.Context, string, *ReplicationFabricsClientGetOptions)` to `(context.Context, string, string, string, *ReplicationFabricsClientGetOptions)`
- Function `*ReplicationFabricsClient.NewListPager` parameter(s) have been changed from `(*ReplicationFabricsClientListOptions)` to `(string, string, *ReplicationFabricsClientListOptions)`
- Function `NewReplicationJobsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationJobsClient.BeginCancel` parameter(s) have been changed from `(context.Context, string, *ReplicationJobsClientBeginCancelOptions)` to `(context.Context, string, string, string, *ReplicationJobsClientBeginCancelOptions)`
- Function `*ReplicationJobsClient.BeginExport` parameter(s) have been changed from `(context.Context, JobQueryParameter, *ReplicationJobsClientBeginExportOptions)` to `(context.Context, string, string, JobQueryParameter, *ReplicationJobsClientBeginExportOptions)`
- Function `*ReplicationJobsClient.BeginRestart` parameter(s) have been changed from `(context.Context, string, *ReplicationJobsClientBeginRestartOptions)` to `(context.Context, string, string, string, *ReplicationJobsClientBeginRestartOptions)`
- Function `*ReplicationJobsClient.BeginResume` parameter(s) have been changed from `(context.Context, string, ResumeJobParams, *ReplicationJobsClientBeginResumeOptions)` to `(context.Context, string, string, string, ResumeJobParams, *ReplicationJobsClientBeginResumeOptions)`
- Function `*ReplicationJobsClient.Get` parameter(s) have been changed from `(context.Context, string, *ReplicationJobsClientGetOptions)` to `(context.Context, string, string, string, *ReplicationJobsClientGetOptions)`
- Function `*ReplicationJobsClient.NewListPager` parameter(s) have been changed from `(*ReplicationJobsClientListOptions)` to `(string, string, *ReplicationJobsClientListOptions)`
- Function `NewReplicationLogicalNetworksClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationLogicalNetworksClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ReplicationLogicalNetworksClientGetOptions)` to `(context.Context, string, string, string, string, *ReplicationLogicalNetworksClientGetOptions)`
- Function `*ReplicationLogicalNetworksClient.NewListByReplicationFabricsPager` parameter(s) have been changed from `(string, *ReplicationLogicalNetworksClientListByReplicationFabricsOptions)` to `(string, string, string, *ReplicationLogicalNetworksClientListByReplicationFabricsOptions)`
- Function `NewReplicationMigrationItemsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationMigrationItemsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, string, EnableMigrationInput, *ReplicationMigrationItemsClientBeginCreateOptions)` to `(context.Context, string, string, string, string, string, EnableMigrationInput, *ReplicationMigrationItemsClientBeginCreateOptions)`
- Function `*ReplicationMigrationItemsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationMigrationItemsClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, string, *ReplicationMigrationItemsClientBeginDeleteOptions)`
- Function `*ReplicationMigrationItemsClient.BeginMigrate` parameter(s) have been changed from `(context.Context, string, string, string, MigrateInput, *ReplicationMigrationItemsClientBeginMigrateOptions)` to `(context.Context, string, string, string, string, string, MigrateInput, *ReplicationMigrationItemsClientBeginMigrateOptions)`
- Function `*ReplicationMigrationItemsClient.BeginPauseReplication` parameter(s) have been changed from `(context.Context, string, string, string, PauseReplicationInput, *ReplicationMigrationItemsClientBeginPauseReplicationOptions)` to `(context.Context, string, string, string, string, string, PauseReplicationInput, *ReplicationMigrationItemsClientBeginPauseReplicationOptions)`
- Function `*ReplicationMigrationItemsClient.BeginResumeReplication` parameter(s) have been changed from `(context.Context, string, string, string, ResumeReplicationInput, *ReplicationMigrationItemsClientBeginResumeReplicationOptions)` to `(context.Context, string, string, string, string, string, ResumeReplicationInput, *ReplicationMigrationItemsClientBeginResumeReplicationOptions)`
- Function `*ReplicationMigrationItemsClient.BeginResync` parameter(s) have been changed from `(context.Context, string, string, string, ResyncInput, *ReplicationMigrationItemsClientBeginResyncOptions)` to `(context.Context, string, string, string, string, string, ResyncInput, *ReplicationMigrationItemsClientBeginResyncOptions)`
- Function `*ReplicationMigrationItemsClient.BeginTestMigrate` parameter(s) have been changed from `(context.Context, string, string, string, TestMigrateInput, *ReplicationMigrationItemsClientBeginTestMigrateOptions)` to `(context.Context, string, string, string, string, string, TestMigrateInput, *ReplicationMigrationItemsClientBeginTestMigrateOptions)`
- Function `*ReplicationMigrationItemsClient.BeginTestMigrateCleanup` parameter(s) have been changed from `(context.Context, string, string, string, TestMigrateCleanupInput, *ReplicationMigrationItemsClientBeginTestMigrateCleanupOptions)` to `(context.Context, string, string, string, string, string, TestMigrateCleanupInput, *ReplicationMigrationItemsClientBeginTestMigrateCleanupOptions)`
- Function `*ReplicationMigrationItemsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, UpdateMigrationItemInput, *ReplicationMigrationItemsClientBeginUpdateOptions)` to `(context.Context, string, string, string, string, string, UpdateMigrationItemInput, *ReplicationMigrationItemsClientBeginUpdateOptions)`
- Function `*ReplicationMigrationItemsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationMigrationItemsClientGetOptions)` to `(context.Context, string, string, string, string, string, *ReplicationMigrationItemsClientGetOptions)`
- Function `*ReplicationMigrationItemsClient.NewListByReplicationProtectionContainersPager` parameter(s) have been changed from `(string, string, *ReplicationMigrationItemsClientListByReplicationProtectionContainersOptions)` to `(string, string, string, string, *ReplicationMigrationItemsClientListByReplicationProtectionContainersOptions)`
- Function `*ReplicationMigrationItemsClient.NewListPager` parameter(s) have been changed from `(*ReplicationMigrationItemsClientListOptions)` to `(string, string, *ReplicationMigrationItemsClientListOptions)`
- Function `NewReplicationNetworkMappingsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationNetworkMappingsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, string, CreateNetworkMappingInput, *ReplicationNetworkMappingsClientBeginCreateOptions)` to `(context.Context, string, string, string, string, string, CreateNetworkMappingInput, *ReplicationNetworkMappingsClientBeginCreateOptions)`
- Function `*ReplicationNetworkMappingsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationNetworkMappingsClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, string, *ReplicationNetworkMappingsClientBeginDeleteOptions)`
- Function `*ReplicationNetworkMappingsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, UpdateNetworkMappingInput, *ReplicationNetworkMappingsClientBeginUpdateOptions)` to `(context.Context, string, string, string, string, string, UpdateNetworkMappingInput, *ReplicationNetworkMappingsClientBeginUpdateOptions)`
- Function `*ReplicationNetworkMappingsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationNetworkMappingsClientGetOptions)` to `(context.Context, string, string, string, string, string, *ReplicationNetworkMappingsClientGetOptions)`
- Function `*ReplicationNetworkMappingsClient.NewListByReplicationNetworksPager` parameter(s) have been changed from `(string, string, *ReplicationNetworkMappingsClientListByReplicationNetworksOptions)` to `(string, string, string, string, *ReplicationNetworkMappingsClientListByReplicationNetworksOptions)`
- Function `*ReplicationNetworkMappingsClient.NewListPager` parameter(s) have been changed from `(*ReplicationNetworkMappingsClientListOptions)` to `(string, string, *ReplicationNetworkMappingsClientListOptions)`
- Function `NewReplicationNetworksClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationNetworksClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ReplicationNetworksClientGetOptions)` to `(context.Context, string, string, string, string, *ReplicationNetworksClientGetOptions)`
- Function `*ReplicationNetworksClient.NewListByReplicationFabricsPager` parameter(s) have been changed from `(string, *ReplicationNetworksClientListByReplicationFabricsOptions)` to `(string, string, string, *ReplicationNetworksClientListByReplicationFabricsOptions)`
- Function `*ReplicationNetworksClient.NewListPager` parameter(s) have been changed from `(*ReplicationNetworksClientListOptions)` to `(string, string, *ReplicationNetworksClientListOptions)`
- Function `NewReplicationPoliciesClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationPoliciesClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, CreatePolicyInput, *ReplicationPoliciesClientBeginCreateOptions)` to `(context.Context, string, string, string, CreatePolicyInput, *ReplicationPoliciesClientBeginCreateOptions)`
- Function `*ReplicationPoliciesClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, *ReplicationPoliciesClientBeginDeleteOptions)` to `(context.Context, string, string, string, *ReplicationPoliciesClientBeginDeleteOptions)`
- Function `*ReplicationPoliciesClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, UpdatePolicyInput, *ReplicationPoliciesClientBeginUpdateOptions)` to `(context.Context, string, string, string, UpdatePolicyInput, *ReplicationPoliciesClientBeginUpdateOptions)`
- Function `*ReplicationPoliciesClient.Get` parameter(s) have been changed from `(context.Context, string, *ReplicationPoliciesClientGetOptions)` to `(context.Context, string, string, string, *ReplicationPoliciesClientGetOptions)`
- Function `*ReplicationPoliciesClient.NewListPager` parameter(s) have been changed from `(*ReplicationPoliciesClientListOptions)` to `(string, string, *ReplicationPoliciesClientListOptions)`
- Function `NewReplicationProtectableItemsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationProtectableItemsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationProtectableItemsClientGetOptions)` to `(context.Context, string, string, string, string, string, *ReplicationProtectableItemsClientGetOptions)`
- Function `*ReplicationProtectableItemsClient.NewListByReplicationProtectionContainersPager` parameter(s) have been changed from `(string, string, *ReplicationProtectableItemsClientListByReplicationProtectionContainersOptions)` to `(string, string, string, string, *ReplicationProtectableItemsClientListByReplicationProtectionContainersOptions)`
- Function `NewReplicationProtectedItemsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationProtectedItemsClient.BeginAddDisks` parameter(s) have been changed from `(context.Context, string, string, string, AddDisksInput, *ReplicationProtectedItemsClientBeginAddDisksOptions)` to `(context.Context, string, string, string, string, string, AddDisksInput, *ReplicationProtectedItemsClientBeginAddDisksOptions)`
- Function `*ReplicationProtectedItemsClient.BeginApplyRecoveryPoint` parameter(s) have been changed from `(context.Context, string, string, string, ApplyRecoveryPointInput, *ReplicationProtectedItemsClientBeginApplyRecoveryPointOptions)` to `(context.Context, string, string, string, string, string, ApplyRecoveryPointInput, *ReplicationProtectedItemsClientBeginApplyRecoveryPointOptions)`
- Function `*ReplicationProtectedItemsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, string, EnableProtectionInput, *ReplicationProtectedItemsClientBeginCreateOptions)` to `(context.Context, string, string, string, string, string, EnableProtectionInput, *ReplicationProtectedItemsClientBeginCreateOptions)`
- Function `*ReplicationProtectedItemsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, string, DisableProtectionInput, *ReplicationProtectedItemsClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, string, DisableProtectionInput, *ReplicationProtectedItemsClientBeginDeleteOptions)`
- Function `*ReplicationProtectedItemsClient.BeginFailoverCancel` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationProtectedItemsClientBeginFailoverCancelOptions)` to `(context.Context, string, string, string, string, string, *ReplicationProtectedItemsClientBeginFailoverCancelOptions)`
- Function `*ReplicationProtectedItemsClient.BeginFailoverCommit` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationProtectedItemsClientBeginFailoverCommitOptions)` to `(context.Context, string, string, string, string, string, *ReplicationProtectedItemsClientBeginFailoverCommitOptions)`
- Function `*ReplicationProtectedItemsClient.BeginPlannedFailover` parameter(s) have been changed from `(context.Context, string, string, string, PlannedFailoverInput, *ReplicationProtectedItemsClientBeginPlannedFailoverOptions)` to `(context.Context, string, string, string, string, string, PlannedFailoverInput, *ReplicationProtectedItemsClientBeginPlannedFailoverOptions)`
- Function `*ReplicationProtectedItemsClient.BeginPurge` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationProtectedItemsClientBeginPurgeOptions)` to `(context.Context, string, string, string, string, string, *ReplicationProtectedItemsClientBeginPurgeOptions)`
- Function `*ReplicationProtectedItemsClient.BeginRemoveDisks` parameter(s) have been changed from `(context.Context, string, string, string, RemoveDisksInput, *ReplicationProtectedItemsClientBeginRemoveDisksOptions)` to `(context.Context, string, string, string, string, string, RemoveDisksInput, *ReplicationProtectedItemsClientBeginRemoveDisksOptions)`
- Function `*ReplicationProtectedItemsClient.BeginRepairReplication` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationProtectedItemsClientBeginRepairReplicationOptions)` to `(context.Context, string, string, string, string, string, *ReplicationProtectedItemsClientBeginRepairReplicationOptions)`
- Function `*ReplicationProtectedItemsClient.BeginReprotect` parameter(s) have been changed from `(context.Context, string, string, string, ReverseReplicationInput, *ReplicationProtectedItemsClientBeginReprotectOptions)` to `(context.Context, string, string, string, string, string, ReverseReplicationInput, *ReplicationProtectedItemsClientBeginReprotectOptions)`
- Function `*ReplicationProtectedItemsClient.BeginResolveHealthErrors` parameter(s) have been changed from `(context.Context, string, string, string, ResolveHealthInput, *ReplicationProtectedItemsClientBeginResolveHealthErrorsOptions)` to `(context.Context, string, string, string, string, string, ResolveHealthInput, *ReplicationProtectedItemsClientBeginResolveHealthErrorsOptions)`
- Function `*ReplicationProtectedItemsClient.BeginSwitchProvider` parameter(s) have been changed from `(context.Context, string, string, string, SwitchProviderInput, *ReplicationProtectedItemsClientBeginSwitchProviderOptions)` to `(context.Context, string, string, string, string, string, SwitchProviderInput, *ReplicationProtectedItemsClientBeginSwitchProviderOptions)`
- Function `*ReplicationProtectedItemsClient.BeginTestFailover` parameter(s) have been changed from `(context.Context, string, string, string, TestFailoverInput, *ReplicationProtectedItemsClientBeginTestFailoverOptions)` to `(context.Context, string, string, string, string, string, TestFailoverInput, *ReplicationProtectedItemsClientBeginTestFailoverOptions)`
- Function `*ReplicationProtectedItemsClient.BeginTestFailoverCleanup` parameter(s) have been changed from `(context.Context, string, string, string, TestFailoverCleanupInput, *ReplicationProtectedItemsClientBeginTestFailoverCleanupOptions)` to `(context.Context, string, string, string, string, string, TestFailoverCleanupInput, *ReplicationProtectedItemsClientBeginTestFailoverCleanupOptions)`
- Function `*ReplicationProtectedItemsClient.BeginUnplannedFailover` parameter(s) have been changed from `(context.Context, string, string, string, UnplannedFailoverInput, *ReplicationProtectedItemsClientBeginUnplannedFailoverOptions)` to `(context.Context, string, string, string, string, string, UnplannedFailoverInput, *ReplicationProtectedItemsClientBeginUnplannedFailoverOptions)`
- Function `*ReplicationProtectedItemsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, UpdateReplicationProtectedItemInput, *ReplicationProtectedItemsClientBeginUpdateOptions)` to `(context.Context, string, string, string, string, string, UpdateReplicationProtectedItemInput, *ReplicationProtectedItemsClientBeginUpdateOptions)`
- Function `*ReplicationProtectedItemsClient.BeginUpdateAppliance` parameter(s) have been changed from `(context.Context, string, string, string, UpdateApplianceForReplicationProtectedItemInput, *ReplicationProtectedItemsClientBeginUpdateApplianceOptions)` to `(context.Context, string, string, string, string, string, UpdateApplianceForReplicationProtectedItemInput, *ReplicationProtectedItemsClientBeginUpdateApplianceOptions)`
- Function `*ReplicationProtectedItemsClient.BeginUpdateMobilityService` parameter(s) have been changed from `(context.Context, string, string, string, UpdateMobilityServiceRequest, *ReplicationProtectedItemsClientBeginUpdateMobilityServiceOptions)` to `(context.Context, string, string, string, string, string, UpdateMobilityServiceRequest, *ReplicationProtectedItemsClientBeginUpdateMobilityServiceOptions)`
- Function `*ReplicationProtectedItemsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationProtectedItemsClientGetOptions)` to `(context.Context, string, string, string, string, string, *ReplicationProtectedItemsClientGetOptions)`
- Function `*ReplicationProtectedItemsClient.NewListByReplicationProtectionContainersPager` parameter(s) have been changed from `(string, string, *ReplicationProtectedItemsClientListByReplicationProtectionContainersOptions)` to `(string, string, string, string, *ReplicationProtectedItemsClientListByReplicationProtectionContainersOptions)`
- Function `*ReplicationProtectedItemsClient.NewListPager` parameter(s) have been changed from `(*ReplicationProtectedItemsClientListOptions)` to `(string, string, *ReplicationProtectedItemsClientListOptions)`
- Function `NewReplicationProtectionContainerMappingsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationProtectionContainerMappingsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, string, CreateProtectionContainerMappingInput, *ReplicationProtectionContainerMappingsClientBeginCreateOptions)` to `(context.Context, string, string, string, string, string, CreateProtectionContainerMappingInput, *ReplicationProtectionContainerMappingsClientBeginCreateOptions)`
- Function `*ReplicationProtectionContainerMappingsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, string, RemoveProtectionContainerMappingInput, *ReplicationProtectionContainerMappingsClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, string, RemoveProtectionContainerMappingInput, *ReplicationProtectionContainerMappingsClientBeginDeleteOptions)`
- Function `*ReplicationProtectionContainerMappingsClient.BeginPurge` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationProtectionContainerMappingsClientBeginPurgeOptions)` to `(context.Context, string, string, string, string, string, *ReplicationProtectionContainerMappingsClientBeginPurgeOptions)`
- Function `*ReplicationProtectionContainerMappingsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, UpdateProtectionContainerMappingInput, *ReplicationProtectionContainerMappingsClientBeginUpdateOptions)` to `(context.Context, string, string, string, string, string, UpdateProtectionContainerMappingInput, *ReplicationProtectionContainerMappingsClientBeginUpdateOptions)`
- Function `*ReplicationProtectionContainerMappingsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationProtectionContainerMappingsClientGetOptions)` to `(context.Context, string, string, string, string, string, *ReplicationProtectionContainerMappingsClientGetOptions)`
- Function `*ReplicationProtectionContainerMappingsClient.NewListByReplicationProtectionContainersPager` parameter(s) have been changed from `(string, string, *ReplicationProtectionContainerMappingsClientListByReplicationProtectionContainersOptions)` to `(string, string, string, string, *ReplicationProtectionContainerMappingsClientListByReplicationProtectionContainersOptions)`
- Function `*ReplicationProtectionContainerMappingsClient.NewListPager` parameter(s) have been changed from `(*ReplicationProtectionContainerMappingsClientListOptions)` to `(string, string, *ReplicationProtectionContainerMappingsClientListOptions)`
- Function `NewReplicationProtectionContainersClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationProtectionContainersClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, CreateProtectionContainerInput, *ReplicationProtectionContainersClientBeginCreateOptions)` to `(context.Context, string, string, string, string, CreateProtectionContainerInput, *ReplicationProtectionContainersClientBeginCreateOptions)`
- Function `*ReplicationProtectionContainersClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *ReplicationProtectionContainersClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, *ReplicationProtectionContainersClientBeginDeleteOptions)`
- Function `*ReplicationProtectionContainersClient.BeginDiscoverProtectableItem` parameter(s) have been changed from `(context.Context, string, string, DiscoverProtectableItemRequest, *ReplicationProtectionContainersClientBeginDiscoverProtectableItemOptions)` to `(context.Context, string, string, string, string, DiscoverProtectableItemRequest, *ReplicationProtectionContainersClientBeginDiscoverProtectableItemOptions)`
- Function `*ReplicationProtectionContainersClient.BeginSwitchProtection` parameter(s) have been changed from `(context.Context, string, string, SwitchProtectionInput, *ReplicationProtectionContainersClientBeginSwitchProtectionOptions)` to `(context.Context, string, string, string, string, SwitchProtectionInput, *ReplicationProtectionContainersClientBeginSwitchProtectionOptions)`
- Function `*ReplicationProtectionContainersClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ReplicationProtectionContainersClientGetOptions)` to `(context.Context, string, string, string, string, *ReplicationProtectionContainersClientGetOptions)`
- Function `*ReplicationProtectionContainersClient.NewListByReplicationFabricsPager` parameter(s) have been changed from `(string, *ReplicationProtectionContainersClientListByReplicationFabricsOptions)` to `(string, string, string, *ReplicationProtectionContainersClientListByReplicationFabricsOptions)`
- Function `*ReplicationProtectionContainersClient.NewListPager` parameter(s) have been changed from `(*ReplicationProtectionContainersClientListOptions)` to `(string, string, *ReplicationProtectionContainersClientListOptions)`
- Function `NewReplicationProtectionIntentsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationProtectionIntentsClient.Create` parameter(s) have been changed from `(context.Context, string, CreateProtectionIntentInput, *ReplicationProtectionIntentsClientCreateOptions)` to `(context.Context, string, string, string, CreateProtectionIntentInput, *ReplicationProtectionIntentsClientCreateOptions)`
- Function `*ReplicationProtectionIntentsClient.Get` parameter(s) have been changed from `(context.Context, string, *ReplicationProtectionIntentsClientGetOptions)` to `(context.Context, string, string, string, *ReplicationProtectionIntentsClientGetOptions)`
- Function `*ReplicationProtectionIntentsClient.NewListPager` parameter(s) have been changed from `(*ReplicationProtectionIntentsClientListOptions)` to `(string, string, *ReplicationProtectionIntentsClientListOptions)`
- Function `NewReplicationRecoveryPlansClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationRecoveryPlansClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, CreateRecoveryPlanInput, *ReplicationRecoveryPlansClientBeginCreateOptions)` to `(context.Context, string, string, string, CreateRecoveryPlanInput, *ReplicationRecoveryPlansClientBeginCreateOptions)`
- Function `*ReplicationRecoveryPlansClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, *ReplicationRecoveryPlansClientBeginDeleteOptions)` to `(context.Context, string, string, string, *ReplicationRecoveryPlansClientBeginDeleteOptions)`
- Function `*ReplicationRecoveryPlansClient.BeginFailoverCancel` parameter(s) have been changed from `(context.Context, string, *ReplicationRecoveryPlansClientBeginFailoverCancelOptions)` to `(context.Context, string, string, string, *ReplicationRecoveryPlansClientBeginFailoverCancelOptions)`
- Function `*ReplicationRecoveryPlansClient.BeginFailoverCommit` parameter(s) have been changed from `(context.Context, string, *ReplicationRecoveryPlansClientBeginFailoverCommitOptions)` to `(context.Context, string, string, string, *ReplicationRecoveryPlansClientBeginFailoverCommitOptions)`
- Function `*ReplicationRecoveryPlansClient.BeginPlannedFailover` parameter(s) have been changed from `(context.Context, string, RecoveryPlanPlannedFailoverInput, *ReplicationRecoveryPlansClientBeginPlannedFailoverOptions)` to `(context.Context, string, string, string, RecoveryPlanPlannedFailoverInput, *ReplicationRecoveryPlansClientBeginPlannedFailoverOptions)`
- Function `*ReplicationRecoveryPlansClient.BeginReprotect` parameter(s) have been changed from `(context.Context, string, *ReplicationRecoveryPlansClientBeginReprotectOptions)` to `(context.Context, string, string, string, *ReplicationRecoveryPlansClientBeginReprotectOptions)`
- Function `*ReplicationRecoveryPlansClient.BeginTestFailover` parameter(s) have been changed from `(context.Context, string, RecoveryPlanTestFailoverInput, *ReplicationRecoveryPlansClientBeginTestFailoverOptions)` to `(context.Context, string, string, string, RecoveryPlanTestFailoverInput, *ReplicationRecoveryPlansClientBeginTestFailoverOptions)`
- Function `*ReplicationRecoveryPlansClient.BeginTestFailoverCleanup` parameter(s) have been changed from `(context.Context, string, RecoveryPlanTestFailoverCleanupInput, *ReplicationRecoveryPlansClientBeginTestFailoverCleanupOptions)` to `(context.Context, string, string, string, RecoveryPlanTestFailoverCleanupInput, *ReplicationRecoveryPlansClientBeginTestFailoverCleanupOptions)`
- Function `*ReplicationRecoveryPlansClient.BeginUnplannedFailover` parameter(s) have been changed from `(context.Context, string, RecoveryPlanUnplannedFailoverInput, *ReplicationRecoveryPlansClientBeginUnplannedFailoverOptions)` to `(context.Context, string, string, string, RecoveryPlanUnplannedFailoverInput, *ReplicationRecoveryPlansClientBeginUnplannedFailoverOptions)`
- Function `*ReplicationRecoveryPlansClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, UpdateRecoveryPlanInput, *ReplicationRecoveryPlansClientBeginUpdateOptions)` to `(context.Context, string, string, string, UpdateRecoveryPlanInput, *ReplicationRecoveryPlansClientBeginUpdateOptions)`
- Function `*ReplicationRecoveryPlansClient.Get` parameter(s) have been changed from `(context.Context, string, *ReplicationRecoveryPlansClientGetOptions)` to `(context.Context, string, string, string, *ReplicationRecoveryPlansClientGetOptions)`
- Function `*ReplicationRecoveryPlansClient.NewListPager` parameter(s) have been changed from `(*ReplicationRecoveryPlansClientListOptions)` to `(string, string, *ReplicationRecoveryPlansClientListOptions)`
- Function `NewReplicationRecoveryServicesProvidersClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationRecoveryServicesProvidersClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, AddRecoveryServicesProviderInput, *ReplicationRecoveryServicesProvidersClientBeginCreateOptions)` to `(context.Context, string, string, string, string, AddRecoveryServicesProviderInput, *ReplicationRecoveryServicesProvidersClientBeginCreateOptions)`
- Function `*ReplicationRecoveryServicesProvidersClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *ReplicationRecoveryServicesProvidersClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, *ReplicationRecoveryServicesProvidersClientBeginDeleteOptions)`
- Function `*ReplicationRecoveryServicesProvidersClient.BeginPurge` parameter(s) have been changed from `(context.Context, string, string, *ReplicationRecoveryServicesProvidersClientBeginPurgeOptions)` to `(context.Context, string, string, string, string, *ReplicationRecoveryServicesProvidersClientBeginPurgeOptions)`
- Function `*ReplicationRecoveryServicesProvidersClient.BeginRefreshProvider` parameter(s) have been changed from `(context.Context, string, string, *ReplicationRecoveryServicesProvidersClientBeginRefreshProviderOptions)` to `(context.Context, string, string, string, string, *ReplicationRecoveryServicesProvidersClientBeginRefreshProviderOptions)`
- Function `*ReplicationRecoveryServicesProvidersClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ReplicationRecoveryServicesProvidersClientGetOptions)` to `(context.Context, string, string, string, string, *ReplicationRecoveryServicesProvidersClientGetOptions)`
- Function `*ReplicationRecoveryServicesProvidersClient.NewListByReplicationFabricsPager` parameter(s) have been changed from `(string, *ReplicationRecoveryServicesProvidersClientListByReplicationFabricsOptions)` to `(string, string, string, *ReplicationRecoveryServicesProvidersClientListByReplicationFabricsOptions)`
- Function `*ReplicationRecoveryServicesProvidersClient.NewListPager` parameter(s) have been changed from `(*ReplicationRecoveryServicesProvidersClientListOptions)` to `(string, string, *ReplicationRecoveryServicesProvidersClientListOptions)`
- Function `NewReplicationStorageClassificationMappingsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationStorageClassificationMappingsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, string, StorageClassificationMappingInput, *ReplicationStorageClassificationMappingsClientBeginCreateOptions)` to `(context.Context, string, string, string, string, string, StorageClassificationMappingInput, *ReplicationStorageClassificationMappingsClientBeginCreateOptions)`
- Function `*ReplicationStorageClassificationMappingsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationStorageClassificationMappingsClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, string, *ReplicationStorageClassificationMappingsClientBeginDeleteOptions)`
- Function `*ReplicationStorageClassificationMappingsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *ReplicationStorageClassificationMappingsClientGetOptions)` to `(context.Context, string, string, string, string, string, *ReplicationStorageClassificationMappingsClientGetOptions)`
- Function `*ReplicationStorageClassificationMappingsClient.NewListByReplicationStorageClassificationsPager` parameter(s) have been changed from `(string, string, *ReplicationStorageClassificationMappingsClientListByReplicationStorageClassificationsOptions)` to `(string, string, string, string, *ReplicationStorageClassificationMappingsClientListByReplicationStorageClassificationsOptions)`
- Function `*ReplicationStorageClassificationMappingsClient.NewListPager` parameter(s) have been changed from `(*ReplicationStorageClassificationMappingsClientListOptions)` to `(string, string, *ReplicationStorageClassificationMappingsClientListOptions)`
- Function `NewReplicationStorageClassificationsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationStorageClassificationsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ReplicationStorageClassificationsClientGetOptions)` to `(context.Context, string, string, string, string, *ReplicationStorageClassificationsClientGetOptions)`
- Function `*ReplicationStorageClassificationsClient.NewListByReplicationFabricsPager` parameter(s) have been changed from `(string, *ReplicationStorageClassificationsClientListByReplicationFabricsOptions)` to `(string, string, string, *ReplicationStorageClassificationsClientListByReplicationFabricsOptions)`
- Function `*ReplicationStorageClassificationsClient.NewListPager` parameter(s) have been changed from `(*ReplicationStorageClassificationsClientListOptions)` to `(string, string, *ReplicationStorageClassificationsClientListOptions)`
- Function `NewReplicationVaultHealthClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationVaultHealthClient.BeginRefresh` parameter(s) have been changed from `(context.Context, *ReplicationVaultHealthClientBeginRefreshOptions)` to `(context.Context, string, string, *ReplicationVaultHealthClientBeginRefreshOptions)`
- Function `*ReplicationVaultHealthClient.Get` parameter(s) have been changed from `(context.Context, *ReplicationVaultHealthClientGetOptions)` to `(context.Context, string, string, *ReplicationVaultHealthClientGetOptions)`
- Function `NewReplicationVaultSettingClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationVaultSettingClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, VaultSettingCreationInput, *ReplicationVaultSettingClientBeginCreateOptions)` to `(context.Context, string, string, string, VaultSettingCreationInput, *ReplicationVaultSettingClientBeginCreateOptions)`
- Function `*ReplicationVaultSettingClient.Get` parameter(s) have been changed from `(context.Context, string, *ReplicationVaultSettingClientGetOptions)` to `(context.Context, string, string, string, *ReplicationVaultSettingClientGetOptions)`
- Function `*ReplicationVaultSettingClient.NewListPager` parameter(s) have been changed from `(*ReplicationVaultSettingClientListOptions)` to `(string, string, *ReplicationVaultSettingClientListOptions)`
- Function `NewReplicationvCentersClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ReplicationvCentersClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, AddVCenterRequest, *ReplicationvCentersClientBeginCreateOptions)` to `(context.Context, string, string, string, string, AddVCenterRequest, *ReplicationvCentersClientBeginCreateOptions)`
- Function `*ReplicationvCentersClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *ReplicationvCentersClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, *ReplicationvCentersClientBeginDeleteOptions)`
- Function `*ReplicationvCentersClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, UpdateVCenterRequest, *ReplicationvCentersClientBeginUpdateOptions)` to `(context.Context, string, string, string, string, UpdateVCenterRequest, *ReplicationvCentersClientBeginUpdateOptions)`
- Function `*ReplicationvCentersClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ReplicationvCentersClientGetOptions)` to `(context.Context, string, string, string, string, *ReplicationvCentersClientGetOptions)`
- Function `*ReplicationvCentersClient.NewListByReplicationFabricsPager` parameter(s) have been changed from `(string, *ReplicationvCentersClientListByReplicationFabricsOptions)` to `(string, string, string, *ReplicationvCentersClientListByReplicationFabricsOptions)`
- Function `*ReplicationvCentersClient.NewListPager` parameter(s) have been changed from `(*ReplicationvCentersClientListOptions)` to `(string, string, *ReplicationvCentersClientListOptions)`
- Function `NewSupportedOperatingSystemsClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*SupportedOperatingSystemsClient.Get` parameter(s) have been changed from `(context.Context, *SupportedOperatingSystemsClientGetOptions)` to `(context.Context, string, string, *SupportedOperatingSystemsClientGetOptions)`
- Function `NewTargetComputeSizesClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*TargetComputeSizesClient.NewListByReplicationProtectedItemsPager` parameter(s) have been changed from `(string, string, string, *TargetComputeSizesClientListByReplicationProtectedItemsOptions)` to `(string, string, string, string, string, *TargetComputeSizesClientListByReplicationProtectedItemsOptions)`

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.1.0 (2022-12-23)
### Features Added

- New value `MigrationItemOperationPauseReplication`, `MigrationItemOperationResumeReplication` added to type alias `MigrationItemOperation`
- New value `MigrationStateMigrationCompletedWithInformation`, `MigrationStateMigrationPartiallySucceeded`, `MigrationStateProtectionSuspended`, `MigrationStateResumeInProgress`, `MigrationStateResumeInitiated`, `MigrationStateSuspendingProtection` added to type alias `MigrationState`
- New value `TestMigrationStateTestMigrationCompletedWithInformation`, `TestMigrationStateTestMigrationPartiallySucceeded` added to type alias `TestMigrationState`
- New function `*ReplicationMigrationItemsClient.BeginPauseReplication(context.Context, string, string, string, PauseReplicationInput, *ReplicationMigrationItemsClientBeginPauseReplicationOptions) (*runtime.Poller[ReplicationMigrationItemsClientPauseReplicationResponse], error)`
- New function `*ReplicationMigrationItemsClient.BeginResumeReplication(context.Context, string, string, string, ResumeReplicationInput, *ReplicationMigrationItemsClientBeginResumeReplicationOptions) (*runtime.Poller[ReplicationMigrationItemsClientResumeReplicationResponse], error)`
- New function `*ResumeReplicationProviderSpecificInput.GetResumeReplicationProviderSpecificInput() *ResumeReplicationProviderSpecificInput`
- New function `*VMwareCbtResumeReplicationInput.GetResumeReplicationProviderSpecificInput() *ResumeReplicationProviderSpecificInput`
- New struct `A2AExtendedLocationDetails`
- New struct `CriticalJobHistoryDetails`
- New struct `PauseReplicationInput`
- New struct `PauseReplicationInputProperties`
- New struct `ReplicationMigrationItemsClientPauseReplicationResponse`
- New struct `ReplicationMigrationItemsClientResumeReplicationResponse`
- New struct `ResumeReplicationInput`
- New struct `ResumeReplicationInputProperties`
- New struct `VMwareCbtResumeReplicationInput`
- New field `ExtendedLocations` in struct `AzureFabricSpecificDetails`
- New field `SeedBlobURI` in struct `InMageRcmProtectedDiskDetails`
- New field `StorageAccountID` in struct `InMageRcmReplicationDetails`
- New field `CriticalJobHistory` in struct `MigrationItemProperties`
- New field `LastMigrationStatus` in struct `MigrationItemProperties`
- New field `LastMigrationTime` in struct `MigrationItemProperties`
- New field `RecoveryServicesProviderID` in struct `MigrationItemProperties`
- New field `ReplicationStatus` in struct `MigrationItemProperties`
- New field `PrimaryExtendedLocation` in struct `RecoveryPlanA2ADetails`
- New field `RecoveryExtendedLocation` in struct `RecoveryPlanA2ADetails`
- New field `PerformSQLBulkRegistration` in struct `VMwareCbtEnableMigrationInput`
- New field `ResumeProgressPercentage` in struct `VMwareCbtMigrationDetails`
- New field `ResumeRetryCount` in struct `VMwareCbtMigrationDetails`
- New field `StorageAccountID` in struct `VMwareCbtMigrationDetails`
- New field `TestNetworkID` in struct `VMwareCbtMigrationDetails`
- New field `SeedBlobURI` in struct `VMwareCbtProtectedDiskDetails`
- New field `TargetBlobURI` in struct `VMwareCbtProtectedDiskDetails`
- New field `RoleSizeToNicCountMap` in struct `VMwareCbtProtectionContainerMappingDetails`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicessiterecovery` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).