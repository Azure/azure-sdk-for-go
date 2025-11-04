# Release History

## 2.0.0-beta.7 (2025-09-11)
### Breaking Changes

- Type of `DistributedAvailabilityGroupProperties.ReplicationMode` has been changed from `*ReplicationMode` to `*ReplicationModeType`
- Enum `ReplicationMode` has been removed
- Field `LastHardenedLsn`, `LinkState`, `PrimaryAvailabilityGroupName`, `SecondaryAvailabilityGroupName`, `SourceEndpoint`, `SourceReplicaID`, `TargetDatabase`, `TargetReplicaID` of struct `DistributedAvailabilityGroupProperties` has been removed
- Field `BackupStorageAccessTier`, `MakeBackupsImmutable` of struct `LongTermRetentionPolicyProperties` has been removed

### Features Added

- New value `PhaseBuildingHyperscaleComponents`, `PhaseLogTransitionInProgress` added to enum type `Phase`
- New value `StorageKeyTypeManagedIdentity` added to enum type `StorageKeyType`
- New enum type `FailoverModeType` with values `FailoverModeTypeManual`, `FailoverModeTypeNone`
- New enum type `FailoverType` with values `FailoverTypeForcedAllowDataLoss`, `FailoverTypePlanned`
- New enum type `InstanceRole` with values `InstanceRolePrimary`, `InstanceRoleSecondary`
- New enum type `LinkRole` with values `LinkRolePrimary`, `LinkRoleSecondary`
- New enum type `ReplicaConnectedState` with values `ReplicaConnectedStateCONNECTED`, `ReplicaConnectedStateDISCONNECTED`
- New enum type `ReplicaSynchronizationHealth` with values `ReplicaSynchronizationHealthHEALTHY`, `ReplicaSynchronizationHealthNOTHEALTHY`, `ReplicaSynchronizationHealthPARTIALLYHEALTHY`
- New enum type `ReplicationModeType` with values `ReplicationModeTypeAsync`, `ReplicationModeTypeSync`
- New enum type `RoleChangeType` with values `RoleChangeTypeForced`, `RoleChangeTypePlanned`
- New enum type `SeedingModeType` with values `SeedingModeTypeAutomatic`, `SeedingModeTypeManual`
- New enum type `ServerCreateMode` with values `ServerCreateModeNormal`, `ServerCreateModeRestore`
- New enum type `SetLegalHoldImmutability` with values `SetLegalHoldImmutabilityDisabled`, `SetLegalHoldImmutabilityEnabled`
- New enum type `TimeBasedImmutability` with values `TimeBasedImmutabilityDisabled`, `TimeBasedImmutabilityEnabled`
- New enum type `TimeBasedImmutabilityMode` with values `TimeBasedImmutabilityModeLocked`, `TimeBasedImmutabilityModeUnlocked`
- New function `*DistributedAvailabilityGroupsClient.BeginFailover(context.Context, string, string, string, DistributedAvailabilityGroupsFailoverRequest, *DistributedAvailabilityGroupsClientBeginFailoverOptions) (*runtime.Poller[DistributedAvailabilityGroupsClientFailoverResponse], error)`
- New function `*DistributedAvailabilityGroupsClient.BeginSetRole(context.Context, string, string, string, DistributedAvailabilityGroupSetRole, *DistributedAvailabilityGroupsClientBeginSetRoleOptions) (*runtime.Poller[DistributedAvailabilityGroupsClientSetRoleResponse], error)`
- New function `*LongTermRetentionBackupsClient.BeginLockTimeBasedImmutability(context.Context, string, string, string, string, *LongTermRetentionBackupsClientBeginLockTimeBasedImmutabilityOptions) (*runtime.Poller[LongTermRetentionBackupsClientLockTimeBasedImmutabilityResponse], error)`
- New function `*LongTermRetentionBackupsClient.BeginLockTimeBasedImmutabilityByResourceGroup(context.Context, string, string, string, string, string, *LongTermRetentionBackupsClientBeginLockTimeBasedImmutabilityByResourceGroupOptions) (*runtime.Poller[LongTermRetentionBackupsClientLockTimeBasedImmutabilityByResourceGroupResponse], error)`
- New function `*LongTermRetentionBackupsClient.BeginRemoveLegalHoldImmutability(context.Context, string, string, string, string, *LongTermRetentionBackupsClientBeginRemoveLegalHoldImmutabilityOptions) (*runtime.Poller[LongTermRetentionBackupsClientRemoveLegalHoldImmutabilityResponse], error)`
- New function `*LongTermRetentionBackupsClient.BeginRemoveLegalHoldImmutabilityByResourceGroup(context.Context, string, string, string, string, string, *LongTermRetentionBackupsClientBeginRemoveLegalHoldImmutabilityByResourceGroupOptions) (*runtime.Poller[LongTermRetentionBackupsClientRemoveLegalHoldImmutabilityByResourceGroupResponse], error)`
- New function `*LongTermRetentionBackupsClient.BeginRemoveTimeBasedImmutability(context.Context, string, string, string, string, *LongTermRetentionBackupsClientBeginRemoveTimeBasedImmutabilityOptions) (*runtime.Poller[LongTermRetentionBackupsClientRemoveTimeBasedImmutabilityResponse], error)`
- New function `*LongTermRetentionBackupsClient.BeginRemoveTimeBasedImmutabilityByResourceGroup(context.Context, string, string, string, string, string, *LongTermRetentionBackupsClientBeginRemoveTimeBasedImmutabilityByResourceGroupOptions) (*runtime.Poller[LongTermRetentionBackupsClientRemoveTimeBasedImmutabilityByResourceGroupResponse], error)`
- New function `*LongTermRetentionBackupsClient.BeginSetLegalHoldImmutability(context.Context, string, string, string, string, *LongTermRetentionBackupsClientBeginSetLegalHoldImmutabilityOptions) (*runtime.Poller[LongTermRetentionBackupsClientSetLegalHoldImmutabilityResponse], error)`
- New function `*LongTermRetentionBackupsClient.BeginSetLegalHoldImmutabilityByResourceGroup(context.Context, string, string, string, string, string, *LongTermRetentionBackupsClientBeginSetLegalHoldImmutabilityByResourceGroupOptions) (*runtime.Poller[LongTermRetentionBackupsClientSetLegalHoldImmutabilityByResourceGroupResponse], error)`
- New function `PossibleTimeBasedImmutabilityValues() []TimeBasedImmutability`
- New struct `CertificateInfo`
- New struct `DistributedAvailabilityGroupDatabase`
- New struct `DistributedAvailabilityGroupSetRole`
- New struct `DistributedAvailabilityGroupsFailoverRequest`
- New field `Databases`, `DistributedAvailabilityGroupName`, `FailoverMode`, `InstanceAvailabilityGroupName`, `InstanceLinkRole`, `PartnerAvailabilityGroupName`, `PartnerEndpoint`, `PartnerLinkRole`, `SeedingMode` in struct `DistributedAvailabilityGroupProperties`
- New field `LegalHoldImmutability`, `TimeBasedImmutability`, `TimeBasedImmutabilityMode` in struct `LongTermRetentionBackupProperties`
- New field `TimeBasedImmutability`, `TimeBasedImmutabilityMode` in struct `LongTermRetentionPolicyProperties`
- New field `CreateMode`, `RetentionDays` in struct `ServerProperties`


## 2.0.0-beta.6 (2024-08-30)
### Breaking Changes

- Type of `DistributedAvailabilityGroupProperties.ReplicationMode` has been changed from `*ReplicationModeType` to `*ReplicationMode`
- Enum `FailoverModeType` has been removed
- Enum `FailoverType` has been removed
- Enum `InstanceRole` has been removed
- Enum `LinkRole` has been removed
- Enum `ReplicaConnectedState` has been removed
- Enum `ReplicaSynchronizationHealth` has been removed
- Enum `ReplicationModeType` has been removed
- Enum `RoleChangeType` has been removed
- Enum `SeedingModeType` has been removed
- Function `*DistributedAvailabilityGroupsClient.BeginFailover` has been removed
- Function `*DistributedAvailabilityGroupsClient.BeginSetRole` has been removed
- Struct `CertificateInfo` has been removed
- Struct `DistributedAvailabilityGroupDatabase` has been removed
- Struct `DistributedAvailabilityGroupSetRole` has been removed
- Struct `DistributedAvailabilityGroupsFailoverRequest` has been removed
- Field `Databases`, `DistributedAvailabilityGroupName`, `FailoverMode`, `InstanceAvailabilityGroupName`, `InstanceLinkRole`, `PartnerAvailabilityGroupName`, `PartnerEndpoint`, `PartnerLinkRole`, `SeedingMode` of struct `DistributedAvailabilityGroupProperties` has been removed

### Features Added

- New enum type `FailoverGroupDatabasesSecondaryType` with values `FailoverGroupDatabasesSecondaryTypeGeo`, `FailoverGroupDatabasesSecondaryTypeStandby`
- New enum type `ReplicationMode` with values `ReplicationModeAsync`, `ReplicationModeSync`
- New function `*ReplicationLinksClient.BeginCreateOrUpdate(context.Context, string, string, string, string, ReplicationLink, *ReplicationLinksClientBeginCreateOrUpdateOptions) (*runtime.Poller[ReplicationLinksClientCreateOrUpdateResponse], error)`
- New function `*ReplicationLinksClient.BeginUpdate(context.Context, string, string, string, string, ReplicationLinkUpdate, *ReplicationLinksClientBeginUpdateOptions) (*runtime.Poller[ReplicationLinksClientUpdateResponse], error)`
- New struct `ReplicationLinkUpdate`
- New struct `ReplicationLinkUpdateProperties`
- New field `LastHardenedLsn`, `LinkState`, `PrimaryAvailabilityGroupName`, `SecondaryAvailabilityGroupName`, `SourceEndpoint`, `SourceReplicaID`, `TargetDatabase`, `TargetReplicaID` in struct `DistributedAvailabilityGroupProperties`
- New field `SecondaryType` in struct `FailoverGroupProperties`
- New field `SecondaryType` in struct `FailoverGroupUpdateProperties`
- New field `PartnerDatabaseID` in struct `ReplicationLinkProperties`


## 2.0.0-beta.5 (2024-05-24)
### Breaking Changes

- Type of `DistributedAvailabilityGroupProperties.ReplicationMode` has been changed from `*ReplicationMode` to `*ReplicationModeType`
- Type of `ManagedInstanceProperties.ProvisioningState` has been changed from `*ManagedInstancePropertiesProvisioningState` to `*ProvisioningState`
- Type of `TopQueries.Queries` has been changed from `[]*QueryStatisticsProperties` to `[]*QueryStatisticsPropertiesAutoGenerated`
- Enum `ManagedInstancePropertiesProvisioningState` has been removed
- Enum `ReplicationMode` has been removed
- Field `LastHardenedLsn`, `LinkState`, `PrimaryAvailabilityGroupName`, `SecondaryAvailabilityGroupName`, `SourceEndpoint`, `SourceReplicaID`, `TargetDatabase`, `TargetReplicaID` of struct `DistributedAvailabilityGroupProperties` has been removed

### Features Added

- New enum type `AuthMetadataLookupModes` with values `AuthMetadataLookupModesAzureAD`, `AuthMetadataLookupModesPaired`, `AuthMetadataLookupModesWindows`
- New enum type `FailoverModeType` with values `FailoverModeTypeManual`, `FailoverModeTypeNone`
- New enum type `FailoverType` with values `FailoverTypeForcedAllowDataLoss`, `FailoverTypePlanned`
- New enum type `FreemiumType` with values `FreemiumTypeFreemium`, `FreemiumTypeRegular`
- New enum type `HybridSecondaryUsage` with values `HybridSecondaryUsageActive`, `HybridSecondaryUsagePassive`
- New enum type `HybridSecondaryUsageDetected` with values `HybridSecondaryUsageDetectedActive`, `HybridSecondaryUsageDetectedPassive`
- New enum type `InstanceRole` with values `InstanceRolePrimary`, `InstanceRoleSecondary`
- New enum type `LinkRole` with values `LinkRolePrimary`, `LinkRoleSecondary`
- New enum type `ManagedInstanceDatabaseFormat` with values `ManagedInstanceDatabaseFormatAlwaysUpToDate`, `ManagedInstanceDatabaseFormatSQLServer2022`
- New enum type `Phase` with values `PhaseCatchup`, `PhaseCopying`, `PhaseCutoverInProgress`, `PhaseWaitingForCutover`
- New enum type `ReplicaConnectedState` with values `ReplicaConnectedStateCONNECTED`, `ReplicaConnectedStateDISCONNECTED`
- New enum type `ReplicaSynchronizationHealth` with values `ReplicaSynchronizationHealthHEALTHY`, `ReplicaSynchronizationHealthNOTHEALTHY`, `ReplicaSynchronizationHealthPARTIALLYHEALTHY`
- New enum type `ReplicationModeType` with values `ReplicationModeTypeAsync`, `ReplicationModeTypeSync`
- New enum type `RoleChangeType` with values `RoleChangeTypeForced`, `RoleChangeTypePlanned`
- New enum type `SeedingModeType` with values `SeedingModeTypeAutomatic`, `SeedingModeTypeManual`
- New function `*DistributedAvailabilityGroupsClient.BeginFailover(context.Context, string, string, string, DistributedAvailabilityGroupsFailoverRequest, *DistributedAvailabilityGroupsClientBeginFailoverOptions) (*runtime.Poller[DistributedAvailabilityGroupsClientFailoverResponse], error)`
- New function `*DistributedAvailabilityGroupsClient.BeginSetRole(context.Context, string, string, string, DistributedAvailabilityGroupSetRole, *DistributedAvailabilityGroupsClientBeginSetRoleOptions) (*runtime.Poller[DistributedAvailabilityGroupsClientSetRoleResponse], error)`
- New function `*ManagedInstancesClient.BeginRefreshStatus(context.Context, string, string, *ManagedInstancesClientBeginRefreshStatusOptions) (*runtime.Poller[ManagedInstancesClientRefreshStatusResponse], error)`
- New function `PossibleHybridSecondaryUsageValues() []HybridSecondaryUsage`
- New struct `CertificateInfo`
- New struct `DistributedAvailabilityGroupDatabase`
- New struct `DistributedAvailabilityGroupSetRole`
- New struct `DistributedAvailabilityGroupsFailoverRequest`
- New struct `PhaseDetails`
- New struct `QueryMetricIntervalAutoGenerated`
- New struct `QueryStatisticsPropertiesAutoGenerated`
- New struct `RefreshExternalGovernanceStatusOperationResultMI`
- New struct `RefreshExternalGovernanceStatusOperationResultPropertiesMI`
- New field `OperationPhaseDetails` in struct `DatabaseOperationProperties`
- New field `Databases`, `DistributedAvailabilityGroupName`, `FailoverMode`, `InstanceAvailabilityGroupName`, `InstanceLinkRole`, `PartnerAvailabilityGroupName`, `PartnerEndpoint`, `PartnerLinkRole`, `SeedingMode` in struct `DistributedAvailabilityGroupProperties`
- New field `AuthenticationMetadata`, `CreateTime`, `DatabaseFormat`, `ExternalGovernanceStatus`, `HybridSecondaryUsage`, `HybridSecondaryUsageDetected`, `IsGeneralPurposeV2`, `PricingModel`, `StorageIOps`, `StorageThroughputMBps`, `VirtualClusterID` in struct `ManagedInstanceProperties`
- New anonymous field `ManagedInstance` in struct `ManagedInstancesClientStartResponse`
- New anonymous field `ManagedInstance` in struct `ManagedInstancesClientStopResponse`


## 2.0.0-beta.4 (2023-12-22)
### Breaking Changes

- Type of `LongTermRetentionPolicy.Properties` has been changed from `*BaseLongTermRetentionPolicyProperties` to `*LongTermRetentionPolicyProperties`
- Type of `ServerProperties.MinimalTLSVersion` has been changed from `*string` to `*MinimalTLSVersion`

### Features Added

- New enum type `BackupStorageAccessTier` with values `BackupStorageAccessTierArchive`, `BackupStorageAccessTierHot`
- New enum type `MinimalTLSVersion` with values `MinimalTLSVersionNone`, `MinimalTLSVersionOne0`, `MinimalTLSVersionOne1`, `MinimalTLSVersionOne2`, `MinimalTLSVersionOne3`
- New function `*ClientFactory.NewJobPrivateEndpointsClient() *JobPrivateEndpointsClient`
- New function `NewJobPrivateEndpointsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*JobPrivateEndpointsClient, error)`
- New function `*JobPrivateEndpointsClient.BeginCreateOrUpdate(context.Context, string, string, string, string, JobPrivateEndpoint, *JobPrivateEndpointsClientBeginCreateOrUpdateOptions) (*runtime.Poller[JobPrivateEndpointsClientCreateOrUpdateResponse], error)`
- New function `*JobPrivateEndpointsClient.BeginDelete(context.Context, string, string, string, string, *JobPrivateEndpointsClientBeginDeleteOptions) (*runtime.Poller[JobPrivateEndpointsClientDeleteResponse], error)`
- New function `*JobPrivateEndpointsClient.Get(context.Context, string, string, string, string, *JobPrivateEndpointsClientGetOptions) (JobPrivateEndpointsClientGetResponse, error)`
- New function `*JobPrivateEndpointsClient.NewListByAgentPager(string, string, string, *JobPrivateEndpointsClientListByAgentOptions) *runtime.Pager[JobPrivateEndpointsClientListByAgentResponse]`
- New function `*LongTermRetentionBackupsClient.BeginChangeAccessTier(context.Context, string, string, string, string, ChangeLongTermRetentionBackupAccessTierParameters, *LongTermRetentionBackupsClientBeginChangeAccessTierOptions) (*runtime.Poller[LongTermRetentionBackupsClientChangeAccessTierResponse], error)`
- New function `*LongTermRetentionBackupsClient.BeginChangeAccessTierByResourceGroup(context.Context, string, string, string, string, string, ChangeLongTermRetentionBackupAccessTierParameters, *LongTermRetentionBackupsClientBeginChangeAccessTierByResourceGroupOptions) (*runtime.Poller[LongTermRetentionBackupsClientChangeAccessTierByResourceGroupResponse], error)`
- New struct `ChangeLongTermRetentionBackupAccessTierParameters`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `JobPrivateEndpoint`
- New struct `JobPrivateEndpointListResult`
- New struct `JobPrivateEndpointProperties`
- New struct `LongTermRetentionPolicyProperties`
- New field `DNSZone`, `MaintenanceConfigurationID` in struct `InstancePoolProperties`
- New field `Properties`, `SKU` in struct `InstancePoolUpdate`
- New field `BackupStorageAccessTier`, `IsBackupImmutable` in struct `LongTermRetentionBackupProperties`


## 2.0.0-beta.3 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0-beta.2 (2023-09-22)
### Features Added

- New enum type `FreeLimitExhaustionBehavior` with values `FreeLimitExhaustionBehaviorAutoPause`, `FreeLimitExhaustionBehaviorBillOverUsage`
- New field `EncryptionProtectorAutoRotation`, `FreeLimitExhaustionBehavior`, `UseFreeLimit` in struct `DatabaseProperties`
- New field `EncryptionProtectorAutoRotation`, `FreeLimitExhaustionBehavior`, `UseFreeLimit` in struct `DatabaseUpdateProperties`
- New field `TargetServer` in struct `FailoverGroupReadOnlyEndpoint`
- New field `PartnerServers` in struct `FailoverGroupUpdateProperties`
- New field `IsIPv6Enabled` in struct `ServerProperties`


## 2.0.0-beta.1 (2023-07-28)
### Breaking Changes

- Function `*ServerDevOpsAuditSettingsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, ServerDevOpsAuditingSettings, *ServerDevOpsAuditSettingsClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, DevOpsAuditingSettingsName, ServerDevOpsAuditingSettings, *ServerDevOpsAuditSettingsClientBeginCreateOrUpdateOptions)`
- Function `*ServerDevOpsAuditSettingsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *ServerDevOpsAuditSettingsClientGetOptions)` to `(context.Context, string, string, DevOpsAuditingSettingsName, *ServerDevOpsAuditSettingsClientGetOptions)`
- Type of `ManagedDatabaseRestoreDetailsProperties.NumberOfFilesDetected` has been changed from `*int64` to `*int32`
- Type of `ManagedDatabaseRestoreDetailsProperties.PercentCompleted` has been changed from `*float64` to `*int32`
- Type of `ManagedDatabaseRestoreDetailsProperties.UnrestorableFiles` has been changed from `[]*string` to `[]*ManagedDatabaseRestoreDetailsUnrestorableFileProperties`
- Type of `ServerProperties.PublicNetworkAccess` has been changed from `*ServerNetworkAccessFlag` to `*ServerPublicNetworkAccessFlag`
- Enum `DNSRefreshConfigurationPropertiesStatus` has been removed
- Operation `*ReplicationLinksClient.Delete` has been changed to LRO, use `*ReplicationLinksClient.BeginDelete` instead.
- Operation `*TransparentDataEncryptionsClient.CreateOrUpdate` has been changed to LRO, use `*TransparentDataEncryptionsClient.BeginCreateOrUpdate` instead.
- Operation `*VirtualClustersClient.UpdateDNSServers` has been changed to LRO, use `*VirtualClustersClient.BeginUpdateDNSServers` instead.
- Struct `DNSRefreshConfigurationProperties` has been removed
- Struct `UpdateManagedInstanceDNSServersOperation` has been removed
- Field `Family`, `MaintenanceConfigurationID` of struct `VirtualClusterProperties` has been removed

### Features Added

- New value `ManagedDatabaseStatusDbCopying`, `ManagedDatabaseStatusDbMoving`, `ManagedDatabaseStatusStarting`, `ManagedDatabaseStatusStopped`, `ManagedDatabaseStatusStopping` added to enum type `ManagedDatabaseStatus`
- New value `ReplicationLinkTypeSTANDBY` added to enum type `ReplicationLinkType`
- New value `SecondaryTypeStandby` added to enum type `SecondaryType`
- New enum type `AlwaysEncryptedEnclaveType` with values `AlwaysEncryptedEnclaveTypeDefault`, `AlwaysEncryptedEnclaveTypeVBS`
- New enum type `AvailabilityZoneType` with values `AvailabilityZoneTypeNoPreference`, `AvailabilityZoneTypeOne`, `AvailabilityZoneTypeThree`, `AvailabilityZoneTypeTwo`
- New enum type `BaselineName` with values `BaselineNameDefault`
- New enum type `DNSRefreshOperationStatus` with values `DNSRefreshOperationStatusFailed`, `DNSRefreshOperationStatusInProgress`, `DNSRefreshOperationStatusSucceeded`
- New enum type `DatabaseKeyType` with values `DatabaseKeyTypeAzureKeyVault`
- New enum type `DevOpsAuditingSettingsName` with values `DevOpsAuditingSettingsNameDefault`
- New enum type `DtcName` with values `DtcNameCurrent`
- New enum type `ExternalGovernanceStatus` with values `ExternalGovernanceStatusDisabled`, `ExternalGovernanceStatusEnabled`
- New enum type `ManagedLedgerDigestUploadsName` with values `ManagedLedgerDigestUploadsNameCurrent`
- New enum type `ManagedLedgerDigestUploadsState` with values `ManagedLedgerDigestUploadsStateDisabled`, `ManagedLedgerDigestUploadsStateEnabled`
- New enum type `MoveOperationMode` with values `MoveOperationModeCopy`, `MoveOperationModeMove`
- New enum type `RuleSeverity` with values `RuleSeverityHigh`, `RuleSeverityInformational`, `RuleSeverityLow`, `RuleSeverityMedium`, `RuleSeverityObsolete`
- New enum type `RuleStatus` with values `RuleStatusFinding`, `RuleStatusInternalError`, `RuleStatusNonFinding`
- New enum type `RuleType` with values `RuleTypeBaselineExpected`, `RuleTypeBinary`, `RuleTypeNegativeList`, `RuleTypePositiveList`
- New enum type `SQLVulnerabilityAssessmentName` with values `SQLVulnerabilityAssessmentNameDefault`
- New enum type `SQLVulnerabilityAssessmentState` with values `SQLVulnerabilityAssessmentStateDisabled`, `SQLVulnerabilityAssessmentStateEnabled`
- New enum type `SecondaryInstanceType` with values `SecondaryInstanceTypeGeo`, `SecondaryInstanceTypeStandby`
- New enum type `ServerConfigurationOptionName` with values `ServerConfigurationOptionNameAllowPolybaseExport`
- New enum type `ServerPublicNetworkAccessFlag` with values `ServerPublicNetworkAccessFlagDisabled`, `ServerPublicNetworkAccessFlagEnabled`, `ServerPublicNetworkAccessFlagSecuredByPerimeter`
- New enum type `StartStopScheduleName` with values `StartStopScheduleNameDefault`
- New function `*ClientFactory.NewDatabaseEncryptionProtectorsClient() *DatabaseEncryptionProtectorsClient`
- New function `*ClientFactory.NewDatabaseSQLVulnerabilityAssessmentBaselinesClient() *DatabaseSQLVulnerabilityAssessmentBaselinesClient`
- New function `*ClientFactory.NewDatabaseSQLVulnerabilityAssessmentExecuteScanClient() *DatabaseSQLVulnerabilityAssessmentExecuteScanClient`
- New function `*ClientFactory.NewDatabaseSQLVulnerabilityAssessmentRuleBaselinesClient() *DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient`
- New function `*ClientFactory.NewDatabaseSQLVulnerabilityAssessmentScanResultClient() *DatabaseSQLVulnerabilityAssessmentScanResultClient`
- New function `*ClientFactory.NewDatabaseSQLVulnerabilityAssessmentScansClient() *DatabaseSQLVulnerabilityAssessmentScansClient`
- New function `*ClientFactory.NewDatabaseSQLVulnerabilityAssessmentsSettingsClient() *DatabaseSQLVulnerabilityAssessmentsSettingsClient`
- New function `*ClientFactory.NewManagedDatabaseAdvancedThreatProtectionSettingsClient() *ManagedDatabaseAdvancedThreatProtectionSettingsClient`
- New function `*ClientFactory.NewManagedDatabaseMoveOperationsClient() *ManagedDatabaseMoveOperationsClient`
- New function `*ClientFactory.NewManagedInstanceAdvancedThreatProtectionSettingsClient() *ManagedInstanceAdvancedThreatProtectionSettingsClient`
- New function `*ClientFactory.NewManagedInstanceDtcsClient() *ManagedInstanceDtcsClient`
- New function `*ClientFactory.NewManagedLedgerDigestUploadsClient() *ManagedLedgerDigestUploadsClient`
- New function `*ClientFactory.NewManagedServerDNSAliasesClient() *ManagedServerDNSAliasesClient`
- New function `*ClientFactory.NewServerConfigurationOptionsClient() *ServerConfigurationOptionsClient`
- New function `*ClientFactory.NewStartStopManagedInstanceSchedulesClient() *StartStopManagedInstanceSchedulesClient`
- New function `*ClientFactory.NewSynapseLinkWorkspacesClient() *SynapseLinkWorkspacesClient`
- New function `*ClientFactory.NewVulnerabilityAssessmentBaselineClient() *VulnerabilityAssessmentBaselineClient`
- New function `*ClientFactory.NewVulnerabilityAssessmentBaselinesClient() *VulnerabilityAssessmentBaselinesClient`
- New function `*ClientFactory.NewVulnerabilityAssessmentExecuteScanClient() *VulnerabilityAssessmentExecuteScanClient`
- New function `*ClientFactory.NewVulnerabilityAssessmentRuleBaselineClient() *VulnerabilityAssessmentRuleBaselineClient`
- New function `*ClientFactory.NewVulnerabilityAssessmentRuleBaselinesClient() *VulnerabilityAssessmentRuleBaselinesClient`
- New function `*ClientFactory.NewVulnerabilityAssessmentScanResultClient() *VulnerabilityAssessmentScanResultClient`
- New function `*ClientFactory.NewVulnerabilityAssessmentScansClient() *VulnerabilityAssessmentScansClient`
- New function `*ClientFactory.NewVulnerabilityAssessmentsClient() *VulnerabilityAssessmentsClient`
- New function `*ClientFactory.NewVulnerabilityAssessmentsSettingsClient() *VulnerabilityAssessmentsSettingsClient`
- New function `NewDatabaseEncryptionProtectorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DatabaseEncryptionProtectorsClient, error)`
- New function `*DatabaseEncryptionProtectorsClient.BeginRevalidate(context.Context, string, string, string, EncryptionProtectorName, *DatabaseEncryptionProtectorsClientBeginRevalidateOptions) (*runtime.Poller[DatabaseEncryptionProtectorsClientRevalidateResponse], error)`
- New function `*DatabaseEncryptionProtectorsClient.BeginRevert(context.Context, string, string, string, EncryptionProtectorName, *DatabaseEncryptionProtectorsClientBeginRevertOptions) (*runtime.Poller[DatabaseEncryptionProtectorsClientRevertResponse], error)`
- New function `NewDatabaseSQLVulnerabilityAssessmentBaselinesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DatabaseSQLVulnerabilityAssessmentBaselinesClient, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentBaselinesClient.CreateOrUpdate(context.Context, string, string, string, VulnerabilityAssessmentName, BaselineName, DatabaseSQLVulnerabilityAssessmentRuleBaselineListInput, *DatabaseSQLVulnerabilityAssessmentBaselinesClientCreateOrUpdateOptions) (DatabaseSQLVulnerabilityAssessmentBaselinesClientCreateOrUpdateResponse, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentBaselinesClient.Get(context.Context, string, string, string, VulnerabilityAssessmentName, BaselineName, *DatabaseSQLVulnerabilityAssessmentBaselinesClientGetOptions) (DatabaseSQLVulnerabilityAssessmentBaselinesClientGetResponse, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentBaselinesClient.NewListBySQLVulnerabilityAssessmentPager(string, string, string, VulnerabilityAssessmentName, *DatabaseSQLVulnerabilityAssessmentBaselinesClientListBySQLVulnerabilityAssessmentOptions) *runtime.Pager[DatabaseSQLVulnerabilityAssessmentBaselinesClientListBySQLVulnerabilityAssessmentResponse]`
- New function `NewDatabaseSQLVulnerabilityAssessmentExecuteScanClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DatabaseSQLVulnerabilityAssessmentExecuteScanClient, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentExecuteScanClient.BeginExecute(context.Context, string, string, string, VulnerabilityAssessmentName, *DatabaseSQLVulnerabilityAssessmentExecuteScanClientBeginExecuteOptions) (*runtime.Poller[DatabaseSQLVulnerabilityAssessmentExecuteScanClientExecuteResponse], error)`
- New function `NewDatabaseSQLVulnerabilityAssessmentRuleBaselinesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.CreateOrUpdate(context.Context, string, string, string, VulnerabilityAssessmentName, BaselineName, string, DatabaseSQLVulnerabilityAssessmentRuleBaselineInput, *DatabaseSQLVulnerabilityAssessmentRuleBaselinesClientCreateOrUpdateOptions) (DatabaseSQLVulnerabilityAssessmentRuleBaselinesClientCreateOrUpdateResponse, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.Delete(context.Context, string, string, string, VulnerabilityAssessmentName, BaselineName, string, *DatabaseSQLVulnerabilityAssessmentRuleBaselinesClientDeleteOptions) (DatabaseSQLVulnerabilityAssessmentRuleBaselinesClientDeleteResponse, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.Get(context.Context, string, string, string, VulnerabilityAssessmentName, BaselineName, string, *DatabaseSQLVulnerabilityAssessmentRuleBaselinesClientGetOptions) (DatabaseSQLVulnerabilityAssessmentRuleBaselinesClientGetResponse, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.NewListByBaselinePager(string, string, string, VulnerabilityAssessmentName, BaselineName, *DatabaseSQLVulnerabilityAssessmentRuleBaselinesClientListByBaselineOptions) *runtime.Pager[DatabaseSQLVulnerabilityAssessmentRuleBaselinesClientListByBaselineResponse]`
- New function `NewDatabaseSQLVulnerabilityAssessmentScanResultClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DatabaseSQLVulnerabilityAssessmentScanResultClient, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentScanResultClient.Get(context.Context, string, string, string, SQLVulnerabilityAssessmentName, string, string, *DatabaseSQLVulnerabilityAssessmentScanResultClientGetOptions) (DatabaseSQLVulnerabilityAssessmentScanResultClientGetResponse, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentScanResultClient.NewListByScanPager(string, string, string, SQLVulnerabilityAssessmentName, string, *DatabaseSQLVulnerabilityAssessmentScanResultClientListByScanOptions) *runtime.Pager[DatabaseSQLVulnerabilityAssessmentScanResultClientListByScanResponse]`
- New function `NewDatabaseSQLVulnerabilityAssessmentScansClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DatabaseSQLVulnerabilityAssessmentScansClient, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentScansClient.Get(context.Context, string, string, string, VulnerabilityAssessmentName, string, *DatabaseSQLVulnerabilityAssessmentScansClientGetOptions) (DatabaseSQLVulnerabilityAssessmentScansClientGetResponse, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentScansClient.NewListBySQLVulnerabilityAssessmentsPager(string, string, string, VulnerabilityAssessmentName, *DatabaseSQLVulnerabilityAssessmentScansClientListBySQLVulnerabilityAssessmentsOptions) *runtime.Pager[DatabaseSQLVulnerabilityAssessmentScansClientListBySQLVulnerabilityAssessmentsResponse]`
- New function `NewDatabaseSQLVulnerabilityAssessmentsSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DatabaseSQLVulnerabilityAssessmentsSettingsClient, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentsSettingsClient.Get(context.Context, string, string, string, SQLVulnerabilityAssessmentName, *DatabaseSQLVulnerabilityAssessmentsSettingsClientGetOptions) (DatabaseSQLVulnerabilityAssessmentsSettingsClientGetResponse, error)`
- New function `*DatabaseSQLVulnerabilityAssessmentsSettingsClient.NewListByDatabasePager(string, string, string, *DatabaseSQLVulnerabilityAssessmentsSettingsClientListByDatabaseOptions) *runtime.Pager[DatabaseSQLVulnerabilityAssessmentsSettingsClientListByDatabaseResponse]`
- New function `*FailoverGroupsClient.BeginTryPlannedBeforeForcedFailover(context.Context, string, string, string, *FailoverGroupsClientBeginTryPlannedBeforeForcedFailoverOptions) (*runtime.Poller[FailoverGroupsClientTryPlannedBeforeForcedFailoverResponse], error)`
- New function `NewManagedDatabaseAdvancedThreatProtectionSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedDatabaseAdvancedThreatProtectionSettingsClient, error)`
- New function `*ManagedDatabaseAdvancedThreatProtectionSettingsClient.CreateOrUpdate(context.Context, string, string, string, AdvancedThreatProtectionName, ManagedDatabaseAdvancedThreatProtection, *ManagedDatabaseAdvancedThreatProtectionSettingsClientCreateOrUpdateOptions) (ManagedDatabaseAdvancedThreatProtectionSettingsClientCreateOrUpdateResponse, error)`
- New function `*ManagedDatabaseAdvancedThreatProtectionSettingsClient.Get(context.Context, string, string, string, AdvancedThreatProtectionName, *ManagedDatabaseAdvancedThreatProtectionSettingsClientGetOptions) (ManagedDatabaseAdvancedThreatProtectionSettingsClientGetResponse, error)`
- New function `*ManagedDatabaseAdvancedThreatProtectionSettingsClient.NewListByDatabasePager(string, string, string, *ManagedDatabaseAdvancedThreatProtectionSettingsClientListByDatabaseOptions) *runtime.Pager[ManagedDatabaseAdvancedThreatProtectionSettingsClientListByDatabaseResponse]`
- New function `NewManagedDatabaseMoveOperationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedDatabaseMoveOperationsClient, error)`
- New function `*ManagedDatabaseMoveOperationsClient.Get(context.Context, string, string, string, *ManagedDatabaseMoveOperationsClientGetOptions) (ManagedDatabaseMoveOperationsClientGetResponse, error)`
- New function `*ManagedDatabaseMoveOperationsClient.NewListByLocationPager(string, string, *ManagedDatabaseMoveOperationsClientListByLocationOptions) *runtime.Pager[ManagedDatabaseMoveOperationsClientListByLocationResponse]`
- New function `*ManagedDatabasesClient.BeginCancelMove(context.Context, string, string, string, ManagedDatabaseMoveDefinition, *ManagedDatabasesClientBeginCancelMoveOptions) (*runtime.Poller[ManagedDatabasesClientCancelMoveResponse], error)`
- New function `*ManagedDatabasesClient.BeginCompleteMove(context.Context, string, string, string, ManagedDatabaseMoveDefinition, *ManagedDatabasesClientBeginCompleteMoveOptions) (*runtime.Poller[ManagedDatabasesClientCompleteMoveResponse], error)`
- New function `*ManagedDatabasesClient.BeginStartMove(context.Context, string, string, string, ManagedDatabaseStartMoveDefinition, *ManagedDatabasesClientBeginStartMoveOptions) (*runtime.Poller[ManagedDatabasesClientStartMoveResponse], error)`
- New function `NewManagedInstanceAdvancedThreatProtectionSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedInstanceAdvancedThreatProtectionSettingsClient, error)`
- New function `*ManagedInstanceAdvancedThreatProtectionSettingsClient.BeginCreateOrUpdate(context.Context, string, string, AdvancedThreatProtectionName, ManagedInstanceAdvancedThreatProtection, *ManagedInstanceAdvancedThreatProtectionSettingsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ManagedInstanceAdvancedThreatProtectionSettingsClientCreateOrUpdateResponse], error)`
- New function `*ManagedInstanceAdvancedThreatProtectionSettingsClient.Get(context.Context, string, string, AdvancedThreatProtectionName, *ManagedInstanceAdvancedThreatProtectionSettingsClientGetOptions) (ManagedInstanceAdvancedThreatProtectionSettingsClientGetResponse, error)`
- New function `*ManagedInstanceAdvancedThreatProtectionSettingsClient.NewListByInstancePager(string, string, *ManagedInstanceAdvancedThreatProtectionSettingsClientListByInstanceOptions) *runtime.Pager[ManagedInstanceAdvancedThreatProtectionSettingsClientListByInstanceResponse]`
- New function `NewManagedInstanceDtcsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedInstanceDtcsClient, error)`
- New function `*ManagedInstanceDtcsClient.BeginCreateOrUpdate(context.Context, string, string, DtcName, ManagedInstanceDtc, *ManagedInstanceDtcsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ManagedInstanceDtcsClientCreateOrUpdateResponse], error)`
- New function `*ManagedInstanceDtcsClient.Get(context.Context, string, string, DtcName, *ManagedInstanceDtcsClientGetOptions) (ManagedInstanceDtcsClientGetResponse, error)`
- New function `*ManagedInstanceDtcsClient.NewListByManagedInstancePager(string, string, *ManagedInstanceDtcsClientListByManagedInstanceOptions) *runtime.Pager[ManagedInstanceDtcsClientListByManagedInstanceResponse]`
- New function `*ManagedInstancesClient.NewListOutboundNetworkDependenciesByManagedInstancePager(string, string, *ManagedInstancesClientListOutboundNetworkDependenciesByManagedInstanceOptions) *runtime.Pager[ManagedInstancesClientListOutboundNetworkDependenciesByManagedInstanceResponse]`
- New function `*ManagedInstancesClient.BeginStart(context.Context, string, string, *ManagedInstancesClientBeginStartOptions) (*runtime.Poller[ManagedInstancesClientStartResponse], error)`
- New function `*ManagedInstancesClient.BeginStop(context.Context, string, string, *ManagedInstancesClientBeginStopOptions) (*runtime.Poller[ManagedInstancesClientStopResponse], error)`
- New function `NewManagedLedgerDigestUploadsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedLedgerDigestUploadsClient, error)`
- New function `*ManagedLedgerDigestUploadsClient.BeginCreateOrUpdate(context.Context, string, string, string, ManagedLedgerDigestUploadsName, ManagedLedgerDigestUploads, *ManagedLedgerDigestUploadsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ManagedLedgerDigestUploadsClientCreateOrUpdateResponse], error)`
- New function `*ManagedLedgerDigestUploadsClient.BeginDisable(context.Context, string, string, string, ManagedLedgerDigestUploadsName, *ManagedLedgerDigestUploadsClientBeginDisableOptions) (*runtime.Poller[ManagedLedgerDigestUploadsClientDisableResponse], error)`
- New function `*ManagedLedgerDigestUploadsClient.Get(context.Context, string, string, string, ManagedLedgerDigestUploadsName, *ManagedLedgerDigestUploadsClientGetOptions) (ManagedLedgerDigestUploadsClientGetResponse, error)`
- New function `*ManagedLedgerDigestUploadsClient.NewListByDatabasePager(string, string, string, *ManagedLedgerDigestUploadsClientListByDatabaseOptions) *runtime.Pager[ManagedLedgerDigestUploadsClientListByDatabaseResponse]`
- New function `NewManagedServerDNSAliasesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedServerDNSAliasesClient, error)`
- New function `*ManagedServerDNSAliasesClient.BeginAcquire(context.Context, string, string, string, ManagedServerDNSAliasAcquisition, *ManagedServerDNSAliasesClientBeginAcquireOptions) (*runtime.Poller[ManagedServerDNSAliasesClientAcquireResponse], error)`
- New function `*ManagedServerDNSAliasesClient.BeginCreateOrUpdate(context.Context, string, string, string, ManagedServerDNSAliasCreation, *ManagedServerDNSAliasesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ManagedServerDNSAliasesClientCreateOrUpdateResponse], error)`
- New function `*ManagedServerDNSAliasesClient.BeginDelete(context.Context, string, string, string, *ManagedServerDNSAliasesClientBeginDeleteOptions) (*runtime.Poller[ManagedServerDNSAliasesClientDeleteResponse], error)`
- New function `*ManagedServerDNSAliasesClient.Get(context.Context, string, string, string, *ManagedServerDNSAliasesClientGetOptions) (ManagedServerDNSAliasesClientGetResponse, error)`
- New function `*ManagedServerDNSAliasesClient.NewListByManagedInstancePager(string, string, *ManagedServerDNSAliasesClientListByManagedInstanceOptions) *runtime.Pager[ManagedServerDNSAliasesClientListByManagedInstanceResponse]`
- New function `NewServerConfigurationOptionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ServerConfigurationOptionsClient, error)`
- New function `*ServerConfigurationOptionsClient.BeginCreateOrUpdate(context.Context, string, string, ServerConfigurationOptionName, ServerConfigurationOption, *ServerConfigurationOptionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ServerConfigurationOptionsClientCreateOrUpdateResponse], error)`
- New function `*ServerConfigurationOptionsClient.Get(context.Context, string, string, ServerConfigurationOptionName, *ServerConfigurationOptionsClientGetOptions) (ServerConfigurationOptionsClientGetResponse, error)`
- New function `*ServerConfigurationOptionsClient.NewListByManagedInstancePager(string, string, *ServerConfigurationOptionsClientListByManagedInstanceOptions) *runtime.Pager[ServerConfigurationOptionsClientListByManagedInstanceResponse]`
- New function `*ServersClient.BeginRefreshStatus(context.Context, string, string, *ServersClientBeginRefreshStatusOptions) (*runtime.Poller[ServersClientRefreshStatusResponse], error)`
- New function `NewStartStopManagedInstanceSchedulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*StartStopManagedInstanceSchedulesClient, error)`
- New function `*StartStopManagedInstanceSchedulesClient.CreateOrUpdate(context.Context, string, string, StartStopScheduleName, StartStopManagedInstanceSchedule, *StartStopManagedInstanceSchedulesClientCreateOrUpdateOptions) (StartStopManagedInstanceSchedulesClientCreateOrUpdateResponse, error)`
- New function `*StartStopManagedInstanceSchedulesClient.Delete(context.Context, string, string, StartStopScheduleName, *StartStopManagedInstanceSchedulesClientDeleteOptions) (StartStopManagedInstanceSchedulesClientDeleteResponse, error)`
- New function `*StartStopManagedInstanceSchedulesClient.Get(context.Context, string, string, StartStopScheduleName, *StartStopManagedInstanceSchedulesClientGetOptions) (StartStopManagedInstanceSchedulesClientGetResponse, error)`
- New function `*StartStopManagedInstanceSchedulesClient.NewListByInstancePager(string, string, *StartStopManagedInstanceSchedulesClientListByInstanceOptions) *runtime.Pager[StartStopManagedInstanceSchedulesClientListByInstanceResponse]`
- New function `NewSynapseLinkWorkspacesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SynapseLinkWorkspacesClient, error)`
- New function `*SynapseLinkWorkspacesClient.NewListByDatabasePager(string, string, string, *SynapseLinkWorkspacesClientListByDatabaseOptions) *runtime.Pager[SynapseLinkWorkspacesClientListByDatabaseResponse]`
- New function `NewVulnerabilityAssessmentBaselineClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VulnerabilityAssessmentBaselineClient, error)`
- New function `*VulnerabilityAssessmentBaselineClient.Get(context.Context, string, string, VulnerabilityAssessmentName, BaselineName, *VulnerabilityAssessmentBaselineClientGetOptions) (VulnerabilityAssessmentBaselineClientGetResponse, error)`
- New function `*VulnerabilityAssessmentBaselineClient.NewListBySQLVulnerabilityAssessmentPager(string, string, VulnerabilityAssessmentName, *VulnerabilityAssessmentBaselineClientListBySQLVulnerabilityAssessmentOptions) *runtime.Pager[VulnerabilityAssessmentBaselineClientListBySQLVulnerabilityAssessmentResponse]`
- New function `NewVulnerabilityAssessmentBaselinesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VulnerabilityAssessmentBaselinesClient, error)`
- New function `*VulnerabilityAssessmentBaselinesClient.CreateOrUpdate(context.Context, string, string, VulnerabilityAssessmentName, BaselineName, DatabaseSQLVulnerabilityAssessmentRuleBaselineListInput, *VulnerabilityAssessmentBaselinesClientCreateOrUpdateOptions) (VulnerabilityAssessmentBaselinesClientCreateOrUpdateResponse, error)`
- New function `NewVulnerabilityAssessmentExecuteScanClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VulnerabilityAssessmentExecuteScanClient, error)`
- New function `*VulnerabilityAssessmentExecuteScanClient.BeginExecute(context.Context, string, string, VulnerabilityAssessmentName, *VulnerabilityAssessmentExecuteScanClientBeginExecuteOptions) (*runtime.Poller[VulnerabilityAssessmentExecuteScanClientExecuteResponse], error)`
- New function `NewVulnerabilityAssessmentRuleBaselineClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VulnerabilityAssessmentRuleBaselineClient, error)`
- New function `*VulnerabilityAssessmentRuleBaselineClient.CreateOrUpdate(context.Context, string, string, VulnerabilityAssessmentName, BaselineName, string, DatabaseSQLVulnerabilityAssessmentRuleBaselineInput, *VulnerabilityAssessmentRuleBaselineClientCreateOrUpdateOptions) (VulnerabilityAssessmentRuleBaselineClientCreateOrUpdateResponse, error)`
- New function `*VulnerabilityAssessmentRuleBaselineClient.Get(context.Context, string, string, VulnerabilityAssessmentName, BaselineName, string, *VulnerabilityAssessmentRuleBaselineClientGetOptions) (VulnerabilityAssessmentRuleBaselineClientGetResponse, error)`
- New function `*VulnerabilityAssessmentRuleBaselineClient.NewListByBaselinePager(string, string, VulnerabilityAssessmentName, BaselineName, *VulnerabilityAssessmentRuleBaselineClientListByBaselineOptions) *runtime.Pager[VulnerabilityAssessmentRuleBaselineClientListByBaselineResponse]`
- New function `NewVulnerabilityAssessmentRuleBaselinesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VulnerabilityAssessmentRuleBaselinesClient, error)`
- New function `*VulnerabilityAssessmentRuleBaselinesClient.Delete(context.Context, string, string, VulnerabilityAssessmentName, BaselineName, string, *VulnerabilityAssessmentRuleBaselinesClientDeleteOptions) (VulnerabilityAssessmentRuleBaselinesClientDeleteResponse, error)`
- New function `NewVulnerabilityAssessmentScanResultClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VulnerabilityAssessmentScanResultClient, error)`
- New function `*VulnerabilityAssessmentScanResultClient.Get(context.Context, string, string, SQLVulnerabilityAssessmentName, string, string, *VulnerabilityAssessmentScanResultClientGetOptions) (VulnerabilityAssessmentScanResultClientGetResponse, error)`
- New function `*VulnerabilityAssessmentScanResultClient.NewListByScanPager(string, string, SQLVulnerabilityAssessmentName, string, *VulnerabilityAssessmentScanResultClientListByScanOptions) *runtime.Pager[VulnerabilityAssessmentScanResultClientListByScanResponse]`
- New function `NewVulnerabilityAssessmentScansClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VulnerabilityAssessmentScansClient, error)`
- New function `*VulnerabilityAssessmentScansClient.Get(context.Context, string, string, VulnerabilityAssessmentName, string, *VulnerabilityAssessmentScansClientGetOptions) (VulnerabilityAssessmentScansClientGetResponse, error)`
- New function `*VulnerabilityAssessmentScansClient.NewListBySQLVulnerabilityAssessmentsPager(string, string, VulnerabilityAssessmentName, *VulnerabilityAssessmentScansClientListBySQLVulnerabilityAssessmentsOptions) *runtime.Pager[VulnerabilityAssessmentScansClientListBySQLVulnerabilityAssessmentsResponse]`
- New function `NewVulnerabilityAssessmentsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VulnerabilityAssessmentsClient, error)`
- New function `*VulnerabilityAssessmentsClient.Delete(context.Context, string, string, VulnerabilityAssessmentName, *VulnerabilityAssessmentsClientDeleteOptions) (VulnerabilityAssessmentsClientDeleteResponse, error)`
- New function `NewVulnerabilityAssessmentsSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VulnerabilityAssessmentsSettingsClient, error)`
- New function `*VulnerabilityAssessmentsSettingsClient.CreateOrUpdate(context.Context, string, string, VulnerabilityAssessmentName, VulnerabilityAssessment, *VulnerabilityAssessmentsSettingsClientCreateOrUpdateOptions) (VulnerabilityAssessmentsSettingsClientCreateOrUpdateResponse, error)`
- New function `*VulnerabilityAssessmentsSettingsClient.Get(context.Context, string, string, SQLVulnerabilityAssessmentName, *VulnerabilityAssessmentsSettingsClientGetOptions) (VulnerabilityAssessmentsSettingsClientGetResponse, error)`
- New function `*VulnerabilityAssessmentsSettingsClient.NewListByServerPager(string, string, *VulnerabilityAssessmentsSettingsClientListByServerOptions) *runtime.Pager[VulnerabilityAssessmentsSettingsClientListByServerResponse]`
- New struct `Baseline`
- New struct `BaselineAdjustedResult`
- New struct `BenchmarkReference`
- New struct `DatabaseKey`
- New struct `DatabaseSQLVulnerabilityAssessmentBaselineSet`
- New struct `DatabaseSQLVulnerabilityAssessmentBaselineSetListResult`
- New struct `DatabaseSQLVulnerabilityAssessmentBaselineSetProperties`
- New struct `DatabaseSQLVulnerabilityAssessmentRuleBaseline`
- New struct `DatabaseSQLVulnerabilityAssessmentRuleBaselineInput`
- New struct `DatabaseSQLVulnerabilityAssessmentRuleBaselineInputProperties`
- New struct `DatabaseSQLVulnerabilityAssessmentRuleBaselineListInput`
- New struct `DatabaseSQLVulnerabilityAssessmentRuleBaselineListInputProperties`
- New struct `DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult`
- New struct `DatabaseSQLVulnerabilityAssessmentRuleBaselineProperties`
- New struct `EndpointDependency`
- New struct `EndpointDetail`
- New struct `ManagedDatabaseAdvancedThreatProtection`
- New struct `ManagedDatabaseAdvancedThreatProtectionListResult`
- New struct `ManagedDatabaseMoveDefinition`
- New struct `ManagedDatabaseMoveOperationListResult`
- New struct `ManagedDatabaseMoveOperationResult`
- New struct `ManagedDatabaseMoveOperationResultProperties`
- New struct `ManagedDatabaseRestoreDetailsBackupSetProperties`
- New struct `ManagedDatabaseRestoreDetailsUnrestorableFileProperties`
- New struct `ManagedDatabaseStartMoveDefinition`
- New struct `ManagedInstanceAdvancedThreatProtection`
- New struct `ManagedInstanceAdvancedThreatProtectionListResult`
- New struct `ManagedInstanceDtc`
- New struct `ManagedInstanceDtcListResult`
- New struct `ManagedInstanceDtcProperties`
- New struct `ManagedInstanceDtcSecuritySettings`
- New struct `ManagedInstanceDtcTransactionManagerCommunicationSettings`
- New struct `ManagedLedgerDigestUploads`
- New struct `ManagedLedgerDigestUploadsListResult`
- New struct `ManagedLedgerDigestUploadsProperties`
- New struct `ManagedServerDNSAlias`
- New struct `ManagedServerDNSAliasAcquisition`
- New struct `ManagedServerDNSAliasCreation`
- New struct `ManagedServerDNSAliasListResult`
- New struct `ManagedServerDNSAliasProperties`
- New struct `OutboundEnvironmentEndpoint`
- New struct `OutboundEnvironmentEndpointCollection`
- New struct `QueryCheck`
- New struct `RefreshExternalGovernanceStatusOperationResult`
- New struct `RefreshExternalGovernanceStatusOperationResultProperties`
- New struct `Remediation`
- New struct `ScheduleItem`
- New struct `ServerConfigurationOption`
- New struct `ServerConfigurationOptionListResult`
- New struct `ServerConfigurationOptionProperties`
- New struct `StartStopManagedInstanceSchedule`
- New struct `StartStopManagedInstanceScheduleListResult`
- New struct `StartStopManagedInstanceScheduleProperties`
- New struct `SynapseLinkWorkspace`
- New struct `SynapseLinkWorkspaceInfoProperties`
- New struct `SynapseLinkWorkspaceListResult`
- New struct `SynapseLinkWorkspaceProperties`
- New struct `UpdateVirtualClusterDNSServersOperation`
- New struct `VaRule`
- New struct `VirtualClusterDNSServersProperties`
- New struct `VulnerabilityAssessment`
- New struct `VulnerabilityAssessmentListResult`
- New struct `VulnerabilityAssessmentPolicyProperties`
- New struct `VulnerabilityAssessmentScanForSQLError`
- New struct `VulnerabilityAssessmentScanListResult`
- New struct `VulnerabilityAssessmentScanRecordForSQL`
- New struct `VulnerabilityAssessmentScanRecordForSQLListResult`
- New struct `VulnerabilityAssessmentScanRecordForSQLProperties`
- New struct `VulnerabilityAssessmentScanResultProperties`
- New struct `VulnerabilityAssessmentScanResults`
- New field `AvailabilityZone`, `EncryptionProtector`, `Keys`, `ManualCutover`, `PerformCutover`, `PreferredEnclaveType` in struct `DatabaseProperties`
- New field `EncryptionProtector`, `Keys`, `ManualCutover`, `PerformCutover`, `PreferredEnclaveType` in struct `DatabaseUpdateProperties`
- New field `Expand`, `Filter` in struct `DatabasesClientGetOptions`
- New field `AvailabilityZone`, `MinCapacity`, `PreferredEnclaveType` in struct `ElasticPoolProperties`
- New field `AvailabilityZone`, `MinCapacity`, `PreferredEnclaveType` in struct `ElasticPoolUpdateProperties`
- New field `SecondaryType` in struct `InstanceFailoverGroupProperties`
- New field `CrossSubscriptionRestorableDroppedDatabaseID`, `CrossSubscriptionSourceDatabaseID`, `CrossSubscriptionTargetManagedInstanceID`, `IsLedgerOn`, `StorageContainerIdentity` in struct `ManagedDatabaseProperties`
- New field `CurrentBackupType`, `CurrentRestorePlanSizeMB`, `CurrentRestoredSizeMB`, `DiffBackupSets`, `FullBackupSets`, `LogBackupSets`, `NumberOfFilesQueued`, `NumberOfFilesRestored`, `NumberOfFilesRestoring`, `NumberOfFilesSkipped`, `NumberOfFilesUnrestorable`, `Type` in struct `ManagedDatabaseRestoreDetailsProperties`
- New field `GroupIDs` in struct `PrivateEndpointConnectionProperties`
- New field `NextLink` in struct `RecoverableDatabaseListResult`
- New field `Keys` in struct `RecoverableDatabaseProperties`
- New field `Expand`, `Filter` in struct `RecoverableDatabasesClientGetOptions`
- New field `Keys` in struct `RestorableDroppedDatabaseProperties`
- New field `Expand`, `Filter` in struct `RestorableDroppedDatabasesClientGetOptions`
- New field `IsManagedIdentityInUse` in struct `ServerDevOpsAuditSettingsProperties`
- New field `ExternalGovernanceStatus` in struct `ServerProperties`
- New field `Version` in struct `VirtualClusterProperties`


## 1.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-06-02)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).