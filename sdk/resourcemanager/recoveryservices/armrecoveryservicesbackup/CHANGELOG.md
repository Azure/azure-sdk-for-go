# Release History

## 5.0.0-beta.1 (2026-02-09)
### Breaking Changes

- Enum `ContainerType` has been removed
- Enum `HealthState` has been removed
- Enum `IntentItemType` has been removed
- Enum `JobOperationType` has been removed
- Enum `JobStatus` has been removed
- Enum `RestorePointQueryType` has been removed
- Enum `Type` has been removed
- Operation `*ProtectedItemsClient.CreateOrUpdate` has been changed to LRO, use `*ProtectedItemsClient.BeginCreateOrUpdate` instead.
- Struct `BMSBackupEngineQueryObject` has been removed
- Struct `BMSBackupEnginesQueryObject` has been removed
- Struct `BMSBackupSummariesQueryObject` has been removed
- Struct `BMSContainerQueryObject` has been removed
- Struct `BMSContainersInquiryQueryObject` has been removed
- Struct `BMSPOQueryObject` has been removed
- Struct `BMSRPQueryObject` has been removed
- Struct `BMSRefreshContainersQueryObject` has been removed
- Struct `BMSWorkloadItemQueryObject` has been removed
- Struct `ErrorAdditionalInfo` has been removed
- Struct `GetProtectedItemQueryObject` has been removed
- Struct `JobQueryObject` has been removed
- Struct `NewErrorResponse` has been removed
- Struct `NewErrorResponseError` has been removed
- Struct `OperationWorkerResponse` has been removed
- Struct `ProtectedItemQueryObject` has been removed
- Struct `ProtectionIntentQueryObject` has been removed
- Struct `ProtectionPolicyQueryObject` has been removed
- Struct `Resource` has been removed
- Struct `ResourceList` has been removed

### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `InstanceProtectionReadiness` with values `InstanceProtectionReadinessPartialProtection`, `InstanceProtectionReadinessProtectionError`, `InstanceProtectionReadinessReady`, `InstanceProtectionReadinessScheduleDisabled`, `InstanceProtectionReadinessUnknown`
- New enum type `ProtectionLevel` with values `ProtectionLevelDatabase`, `ProtectionLevelDatabaseUnderInstance`
- New enum type `SourceSideScanStatus` with values `SourceSideScanStatusConfigured`, `SourceSideScanStatusNotApplicable`, `SourceSideScanStatusNotConfigured`
- New enum type `SourceSideScanSummary` with values `SourceSideScanSummaryHealthy`, `SourceSideScanSummaryNotApplicable`, `SourceSideScanSummarySuspicious`, `SourceSideScanSummaryUnknown`
- New enum type `ThreatSeverity` with values `ThreatSeverityCritical`, `ThreatSeverityHigh`, `ThreatSeverityInformational`, `ThreatSeverityWarning`
- New enum type `ThreatState` with values `ThreatStateActive`, `ThreatStateIgnored`, `ThreatStateInProgress`, `ThreatStateResolved`
- New enum type `ThreatStatus` with values `ThreatStatusHealthy`, `ThreatStatusNotAvailable`, `ThreatStatusUnHealthy`, `ThreatStatusUnknown`, `ThreatStatusWarning`
- New enum type `VMWorkloadPolicyType` with values `VMWorkloadPolicyTypeInvalid`, `VMWorkloadPolicyTypeSnapshotV1`, `VMWorkloadPolicyTypeSnapshotV2`, `VMWorkloadPolicyTypeStreaming`
- New function `*AzureVMWorkloadSAPHanaScaleoutProtectableItem.GetAzureVMWorkloadProtectableItem() *AzureVMWorkloadProtectableItem`
- New function `*AzureVMWorkloadSAPHanaScaleoutProtectableItem.GetWorkloadProtectableItem() *WorkloadProtectableItem`
- New function `*AzureVMWorkloadSQLInstanceProtectedItem.GetAzureVMWorkloadProtectedItem() *AzureVMWorkloadProtectedItem`
- New function `*AzureVMWorkloadSQLInstanceProtectedItem.GetProtectedItem() *ProtectedItem`
- New function `*RecoveryPointsClient.Update(ctx context.Context, resourceGroupName string, vaultName string, fabricName string, containerName string, protectedItemName string, recoveryPointID string, parameters UpdateRecoveryPointRequest, options *RecoveryPointsClientUpdateOptions) (RecoveryPointsClientUpdateResponse, error)`
- New struct `AzureVMWorkloadSAPHanaScaleoutProtectableItem`
- New struct `AzureVMWorkloadSQLInstanceProtectedItem`
- New struct `DatabaseInRP`
- New struct `OkResponse`
- New struct `PatchRecoveryPointInput`
- New struct `PatchRecoveryPointPropertiesInput`
- New struct `SourceSideScanInfo`
- New struct `SystemData`
- New struct `ThreatInfo`
- New struct `UpdateRecoveryPointRequest`
- New field `ThreatInfo`, `ThreatStatus` in struct `AzureFileShareRecoveryPoint`
- New field `SourceSideScanInfo` in struct `AzureFileshareProtectedItem`
- New field `SourceSideScanInfo` in struct `AzureIaaSClassicComputeVMProtectedItem`
- New field `SourceSideScanInfo` in struct `AzureIaaSComputeVMProtectedItem`
- New field `SourceSideScanInfo` in struct `AzureIaaSVMProtectedItem`
- New field `SourceSideScanInfo` in struct `AzureSQLProtectedItem`
- New field `VMWorkloadPolicyType` in struct `AzureVMWorkloadProtectionPolicy`
- New field `SourceSideScanInfo` in struct `AzureVMWorkloadSAPAseDatabaseProtectedItem`
- New field `SourceSideScanInfo` in struct `AzureVMWorkloadSAPHanaDBInstanceProtectedItem`
- New field `SourceSideScanInfo` in struct `AzureVMWorkloadSAPHanaDatabaseProtectedItem`
- New field `ParentProtectedItem`, `ProtectionLevel`, `SourceSideScanInfo` in struct `AzureVMWorkloadSQLDatabaseProtectedItem`
- New field `ThreatInfo`, `ThreatStatus` in struct `AzureWorkloadPointInTimeRecoveryPoint`
- New field `ThreatInfo`, `ThreatStatus` in struct `AzureWorkloadRecoveryPoint`
- New field `ThreatInfo`, `ThreatStatus` in struct `AzureWorkloadSAPAsePointInTimeRecoveryPoint`
- New field `ThreatInfo`, `ThreatStatus` in struct `AzureWorkloadSAPAseRecoveryPoint`
- New field `ThreatInfo`, `ThreatStatus` in struct `AzureWorkloadSAPHanaPointInTimeRecoveryPoint`
- New field `ThreatInfo`, `ThreatStatus` in struct `AzureWorkloadSAPHanaRecoveryPoint`
- New field `ThreatInfo`, `ThreatStatus` in struct `AzureWorkloadSQLPointInTimeRecoveryPoint`
- New field `ThreatInfo`, `ThreatStatus` in struct `AzureWorkloadSQLRecoveryPoint`
- New field `IncludedDatabases` in struct `AzureWorkloadSQLRecoveryPointExtendedInfo`
- New field `SystemData` in struct `BackupEngineBaseResource`
- New field `NextLink` in struct `BackupManagementUsageList`
- New field `SystemData` in struct `BackupRequestResource`
- New field `SystemData` in struct `BackupResourceConfigResource`
- New field `SystemData` in struct `BackupResourceEncryptionConfigExtendedResource`
- New field `SystemData` in struct `BackupResourceEncryptionConfigResource`
- New field `SystemData` in struct `BackupResourceVaultConfigResource`
- New anonymous field `OkResponse` in struct `ClientBMSPrepareDataMoveResponse`
- New anonymous field `OkResponse` in struct `ClientBMSTriggerDataMoveResponse`
- New field `SourceSideScanInfo` in struct `DPMProtectedItem`
- New field `SourceSideScanInfo` in struct `GenericProtectedItem`
- New field `ThreatInfo`, `ThreatStatus` in struct `GenericRecoveryPoint`
- New field `SystemData` in struct `ILRRequestResource`
- New field `ThreatInfo`, `ThreatStatus` in struct `IaasVMRecoveryPoint`
- New field `SystemData` in struct `JobResource`
- New field `SourceSideScanInfo` in struct `MabFileFolderProtectedItem`
- New field `SystemData` in struct `PrivateEndpointConnectionResource`
- New field `SystemData` in struct `ProtectableContainerResource`
- New field `SystemData` in struct `ProtectedItemResource`
- New field `SystemData` in struct `ProtectionContainerResource`
- New field `SystemData` in struct `ProtectionIntentResource`
- New field `SystemData` in struct `ProtectionPolicyResource`
- New field `ThreatInfo`, `ThreatStatus` in struct `RecoveryPoint`
- New field `SystemData` in struct `RecoveryPointResource`
- New field `SystemData` in struct `ResourceGuardProxyBaseResource`
- New field `SystemData` in struct `RestoreRequestResource`
- New field `SystemData` in struct `WorkloadItemResource`
- New field `SystemData` in struct `WorkloadProtectableItemResource`


## 4.2.0 (2025-04-24)
### Features Added

- New value `OperationTypeRehydrate` added to enum type `OperationType`
- New function `*AzureVMWorkloadSAPAseDatabaseProtectableItem.GetAzureVMWorkloadProtectableItem() *AzureVMWorkloadProtectableItem`
- New function `*AzureVMWorkloadSAPAseDatabaseProtectableItem.GetWorkloadProtectableItem() *WorkloadProtectableItem`
- New function `*AzureWorkloadSAPAsePointInTimeRecoveryPoint.GetAzureWorkloadPointInTimeRecoveryPoint() *AzureWorkloadPointInTimeRecoveryPoint`
- New function `*AzureWorkloadSAPAsePointInTimeRecoveryPoint.GetAzureWorkloadRecoveryPoint() *AzureWorkloadRecoveryPoint`
- New function `*AzureWorkloadSAPAsePointInTimeRecoveryPoint.GetRecoveryPoint() *RecoveryPoint`
- New function `*AzureWorkloadSAPAsePointInTimeRestoreRequest.GetAzureWorkloadRestoreRequest() *AzureWorkloadRestoreRequest`
- New function `*AzureWorkloadSAPAsePointInTimeRestoreRequest.GetAzureWorkloadSAPAseRestoreRequest() *AzureWorkloadSAPAseRestoreRequest`
- New function `*AzureWorkloadSAPAsePointInTimeRestoreRequest.GetRestoreRequest() *RestoreRequest`
- New function `*AzureWorkloadSAPAseRecoveryPoint.GetAzureWorkloadRecoveryPoint() *AzureWorkloadRecoveryPoint`
- New function `*AzureWorkloadSAPAseRecoveryPoint.GetRecoveryPoint() *RecoveryPoint`
- New function `*AzureWorkloadSAPAseRestoreRequest.GetAzureWorkloadRestoreRequest() *AzureWorkloadRestoreRequest`
- New function `*AzureWorkloadSAPAseRestoreRequest.GetAzureWorkloadSAPAseRestoreRequest() *AzureWorkloadSAPAseRestoreRequest`
- New function `*AzureWorkloadSAPAseRestoreRequest.GetRestoreRequest() *RestoreRequest`
- New struct `AzureVMWorkloadSAPAseDatabaseProtectableItem`
- New struct `AzureWorkloadSAPAsePointInTimeRecoveryPoint`
- New struct `AzureWorkloadSAPAsePointInTimeRestoreRequest`
- New struct `AzureWorkloadSAPAseRecoveryPoint`
- New field `RecoveryPointTierDetails` in struct `AzureFileShareRecoveryPoint`
- New field `PolicyType` in struct `AzureIaaSClassicComputeVMProtectedItem`
- New field `PolicyType` in struct `AzureIaaSComputeVMProtectedItem`
- New field `PolicyType` in struct `AzureIaaSVMProtectedItem`
- New field `OperationType` in struct `AzureStorageContainer`


## 4.1.0 (2024-05-24)
### Features Added

- New enum type `IaasVMSnapshotConsistencyType` with values `IaasVMSnapshotConsistencyTypeOnlyCrashConsistent`
- New field `ResourceGuardOperationRequests` in struct `AzureFileShareRestoreRequest`
- New field `SnapshotConsistencyType` in struct `AzureIaaSVMProtectionPolicy`
- New field `ResourceGuardOperationRequests` in struct `AzureWorkloadPointInTimeRestoreRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureWorkloadRestoreRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureWorkloadSAPHanaPointInTimeRestoreRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureWorkloadSAPHanaPointInTimeRestoreWithRehydrateRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureWorkloadSAPHanaRestoreRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureWorkloadSAPHanaRestoreWithRehydrateRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureWorkloadSQLPointInTimeRestoreRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureWorkloadSQLPointInTimeRestoreWithRehydrateRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureWorkloadSQLRestoreRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureWorkloadSQLRestoreWithRehydrateRequest`
- New field `XMSAuthorizationAuxiliary` in struct `BackupResourceVaultConfigsClientPutOptions`
- New field `XMSAuthorizationAuxiliary` in struct `BackupResourceVaultConfigsClientUpdateOptions`
- New field `ResourceGuardOperationRequests` in struct `IaasVMRestoreRequest`
- New field `ResourceGuardOperationRequests` in struct `IaasVMRestoreWithRehydrationRequest`
- New field `XMSAuthorizationAuxiliary` in struct `ProtectedItemsClientCreateOrUpdateOptions`
- New field `XMSAuthorizationAuxiliary` in struct `ProtectionPoliciesClientCreateOrUpdateOptions`
- New field `ResourceGuardOperationRequests` in struct `RestoreRequest`
- New field `XMSAuthorizationAuxiliary` in struct `RestoresClientBeginTriggerOptions`
- New field `XMSAuthorizationAuxiliary` in struct `SecurityPINsClientGetOptions`


## 4.0.0 (2023-12-08)
### Breaking Changes

- Function `*OperationClient.Validate` parameter(s) have been changed from `(context.Context, string, string, ValidateOperationRequestClassification, *OperationClientValidateOptions)` to `(context.Context, string, string, ValidateOperationRequestResource, *OperationClientValidateOptions)`
- Function `*ValidateOperationClient.BeginTrigger` parameter(s) have been changed from `(context.Context, string, string, ValidateOperationRequestClassification, *ValidateOperationClientBeginTriggerOptions)` to `(context.Context, string, string, ValidateOperationRequestResource, *ValidateOperationClientBeginTriggerOptions)`
- Operation `*ProtectionContainersClient.Register` has been changed to LRO, use `*ProtectionContainersClient.BeginRegister` instead.

### Features Added

- New value `RecoveryModeRecoveryUsingSnapshot`, `RecoveryModeSnapshotAttach`, `RecoveryModeSnapshotAttachAndRecover` added to enum type `RecoveryMode`
- New function `*ClientFactory.NewFetchTieringCostClient() *FetchTieringCostClient`
- New function `*ClientFactory.NewGetTieringCostOperationResultClient() *GetTieringCostOperationResultClient`
- New function `*ClientFactory.NewTieringCostOperationStatusClient() *TieringCostOperationStatusClient`
- New function `NewFetchTieringCostClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FetchTieringCostClient, error)`
- New function `*FetchTieringCostClient.BeginPost(context.Context, string, string, FetchTieringCostInfoRequestClassification, *FetchTieringCostClientBeginPostOptions) (*runtime.Poller[FetchTieringCostClientPostResponse], error)`
- New function `*FetchTieringCostInfoForRehydrationRequest.GetFetchTieringCostInfoRequest() *FetchTieringCostInfoRequest`
- New function `*FetchTieringCostInfoRequest.GetFetchTieringCostInfoRequest() *FetchTieringCostInfoRequest`
- New function `*FetchTieringCostSavingsInfoForPolicyRequest.GetFetchTieringCostInfoRequest() *FetchTieringCostInfoRequest`
- New function `*FetchTieringCostSavingsInfoForProtectedItemRequest.GetFetchTieringCostInfoRequest() *FetchTieringCostInfoRequest`
- New function `*FetchTieringCostSavingsInfoForVaultRequest.GetFetchTieringCostInfoRequest() *FetchTieringCostInfoRequest`
- New function `NewGetTieringCostOperationResultClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GetTieringCostOperationResultClient, error)`
- New function `*GetTieringCostOperationResultClient.Get(context.Context, string, string, string, *GetTieringCostOperationResultClientGetOptions) (GetTieringCostOperationResultClientGetResponse, error)`
- New function `*TieringCostInfo.GetTieringCostInfo() *TieringCostInfo`
- New function `NewTieringCostOperationStatusClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TieringCostOperationStatusClient, error)`
- New function `*TieringCostOperationStatusClient.Get(context.Context, string, string, string, *TieringCostOperationStatusClientGetOptions) (TieringCostOperationStatusClientGetResponse, error)`
- New function `*TieringCostRehydrationInfo.GetTieringCostInfo() *TieringCostInfo`
- New function `*TieringCostSavingInfo.GetTieringCostInfo() *TieringCostInfo`
- New struct `FetchTieringCostInfoForRehydrationRequest`
- New struct `FetchTieringCostSavingsInfoForPolicyRequest`
- New struct `FetchTieringCostSavingsInfoForProtectedItemRequest`
- New struct `FetchTieringCostSavingsInfoForVaultRequest`
- New struct `SnapshotBackupAdditionalDetails`
- New struct `SnapshotRestoreParameters`
- New struct `TieringCostRehydrationInfo`
- New struct `TieringCostSavingInfo`
- New struct `UserAssignedIdentityProperties`
- New struct `UserAssignedManagedIdentityDetails`
- New struct `ValidateOperationRequestResource`
- New struct `VaultRetentionPolicy`
- New field `VaultRetentionPolicy` in struct `AzureFileShareProtectionPolicy`
- New field `VaultID` in struct `AzureFileshareProtectedItem`
- New field `VaultID` in struct `AzureIaaSClassicComputeVMProtectedItem`
- New field `VaultID` in struct `AzureIaaSComputeVMProtectedItem`
- New field `VaultID` in struct `AzureIaaSVMProtectedItem`
- New field `VaultID` in struct `AzureSQLProtectedItem`
- New field `VaultID` in struct `AzureVMWorkloadProtectedItem`
- New field `VaultID` in struct `AzureVMWorkloadSAPAseDatabaseProtectedItem`
- New field `VaultID` in struct `AzureVMWorkloadSAPHanaDBInstanceProtectedItem`
- New field `VaultID` in struct `AzureVMWorkloadSAPHanaDatabaseProtectedItem`
- New field `VaultID` in struct `AzureVMWorkloadSQLDatabaseProtectedItem`
- New field `SnapshotRestoreParameters`, `TargetResourceGroupName`, `UserAssignedManagedIdentityDetails` in struct `AzureWorkloadPointInTimeRestoreRequest`
- New field `SnapshotRestoreParameters`, `TargetResourceGroupName`, `UserAssignedManagedIdentityDetails` in struct `AzureWorkloadRestoreRequest`
- New field `SnapshotRestoreParameters`, `TargetResourceGroupName`, `UserAssignedManagedIdentityDetails` in struct `AzureWorkloadSAPHanaPointInTimeRestoreRequest`
- New field `SnapshotRestoreParameters`, `TargetResourceGroupName`, `UserAssignedManagedIdentityDetails` in struct `AzureWorkloadSAPHanaPointInTimeRestoreWithRehydrateRequest`
- New field `SnapshotRestoreParameters`, `TargetResourceGroupName`, `UserAssignedManagedIdentityDetails` in struct `AzureWorkloadSAPHanaRestoreRequest`
- New field `SnapshotRestoreParameters`, `TargetResourceGroupName`, `UserAssignedManagedIdentityDetails` in struct `AzureWorkloadSAPHanaRestoreWithRehydrateRequest`
- New field `SnapshotRestoreParameters`, `TargetResourceGroupName`, `UserAssignedManagedIdentityDetails` in struct `AzureWorkloadSQLPointInTimeRestoreRequest`
- New field `SnapshotRestoreParameters`, `TargetResourceGroupName`, `UserAssignedManagedIdentityDetails` in struct `AzureWorkloadSQLPointInTimeRestoreWithRehydrateRequest`
- New field `SnapshotRestoreParameters`, `TargetResourceGroupName`, `UserAssignedManagedIdentityDetails` in struct `AzureWorkloadSQLRestoreRequest`
- New field `SnapshotRestoreParameters`, `TargetResourceGroupName`, `UserAssignedManagedIdentityDetails` in struct `AzureWorkloadSQLRestoreWithRehydrateRequest`
- New field `VaultID` in struct `DPMProtectedItem`
- New field `VaultID` in struct `GenericProtectedItem`
- New field `ExtendedLocation` in struct `IaasVMRecoveryPoint`
- New field `VaultID` in struct `MabFileFolderProtectedItem`
- New field `VaultID` in struct `ProtectedItem`
- New field `SnapshotBackupAdditionalDetails` in struct `SubProtectionPolicy`


## 3.1.0 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 3.0.0 (2023-09-22)
### Breaking Changes

- Function `*AzureVMWorkloadSAPHanaHSR.GetAzureVMWorkloadProtectableItem` has been removed
- Function `*AzureVMWorkloadSAPHanaHSR.GetWorkloadProtectableItem` has been removed
- Struct `AzureVMWorkloadSAPHanaHSR` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `AzureFileshareProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `AzureIaaSClassicComputeVMProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `AzureIaaSComputeVMProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `AzureIaaSVMProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `AzureSQLProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `AzureVMWorkloadProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `AzureVMWorkloadSAPAseDatabaseProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `AzureVMWorkloadSAPHanaDBInstanceProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `AzureVMWorkloadSAPHanaDatabaseProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `AzureVMWorkloadSQLDatabaseProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `DPMProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `GenericProtectedItem` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `MabFileFolderProtectedItem` has been removed
- Field `ActionRequired` of struct `PrivateLinkServiceConnectionState` has been removed
- Field `SoftDeleteRetentionPeriod` of struct `ProtectedItem` has been removed

### Features Added

- New value `SoftDeleteFeatureStateAlwaysON` added to enum type `SoftDeleteFeatureState`
- New enum type `VaultSubResourceType` with values `VaultSubResourceTypeAzureBackup`, `VaultSubResourceTypeAzureBackupSecondary`, `VaultSubResourceTypeAzureSiteRecovery`
- New function `*AzureVMWorkloadSAPHanaHSRProtectableItem.GetAzureVMWorkloadProtectableItem() *AzureVMWorkloadProtectableItem`
- New function `*AzureVMWorkloadSAPHanaHSRProtectableItem.GetWorkloadProtectableItem() *WorkloadProtectableItem`
- New struct `AzureVMWorkloadSAPHanaHSRProtectableItem`
- New field `SoftDeleteRetentionPeriodInDays` in struct `AzureFileshareProtectedItem`
- New field `SoftDeleteRetentionPeriodInDays` in struct `AzureIaaSClassicComputeVMProtectedItem`
- New field `SoftDeleteRetentionPeriodInDays` in struct `AzureIaaSComputeVMProtectedItem`
- New field `SoftDeleteRetentionPeriodInDays` in struct `AzureIaaSVMProtectedItem`
- New field `SoftDeleteRetentionPeriodInDays` in struct `AzureSQLProtectedItem`
- New field `NodesList`, `SoftDeleteRetentionPeriodInDays` in struct `AzureVMWorkloadProtectedItem`
- New field `NodesList`, `SoftDeleteRetentionPeriodInDays` in struct `AzureVMWorkloadSAPAseDatabaseProtectedItem`
- New field `IsProtectable` in struct `AzureVMWorkloadSAPAseSystemProtectableItem`
- New field `IsProtectable` in struct `AzureVMWorkloadSAPHanaDBInstance`
- New field `NodesList`, `SoftDeleteRetentionPeriodInDays` in struct `AzureVMWorkloadSAPHanaDBInstanceProtectedItem`
- New field `IsProtectable` in struct `AzureVMWorkloadSAPHanaDatabaseProtectableItem`
- New field `NodesList`, `SoftDeleteRetentionPeriodInDays` in struct `AzureVMWorkloadSAPHanaDatabaseProtectedItem`
- New field `IsProtectable` in struct `AzureVMWorkloadSAPHanaSystemProtectableItem`
- New field `IsProtectable`, `NodesList` in struct `AzureVMWorkloadSQLAvailabilityGroupProtectableItem`
- New field `IsProtectable` in struct `AzureVMWorkloadSQLDatabaseProtectableItem`
- New field `NodesList`, `SoftDeleteRetentionPeriodInDays` in struct `AzureVMWorkloadSQLDatabaseProtectedItem`
- New field `IsProtectable` in struct `AzureVMWorkloadSQLInstanceProtectableItem`
- New field `SoftDeleteRetentionPeriodInDays` in struct `BackupResourceVaultConfig`
- New field `AcquireStorageAccountLock`, `ProtectedItemsCount` in struct `BackupStatusResponse`
- New field `SoftDeleteRetentionPeriodInDays` in struct `DPMProtectedItem`
- New field `SourceResourceID` in struct `DistributedNodesInfo`
- New field `SoftDeleteRetentionPeriodInDays` in struct `GenericProtectedItem`
- New field `ProtectableItemCount` in struct `InquiryValidation`
- New field `SoftDeleteRetentionPeriodInDays` in struct `MabFileFolderProtectedItem`
- New field `GroupIDs` in struct `PrivateEndpointConnection`
- New field `ActionsRequired` in struct `PrivateLinkServiceConnectionState`
- New field `SoftDeleteRetentionPeriodInDays` in struct `ProtectedItem`


## 2.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 2.1.0 (2023-03-24)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New enum type `TargetDiskNetworkAccessOption` with values `TargetDiskNetworkAccessOptionEnablePrivateAccessForAllDisks`, `TargetDiskNetworkAccessOptionEnablePublicAccessForAllDisks`, `TargetDiskNetworkAccessOptionSameAsOnSourceDisks`
- New struct `ExtendedLocation`
- New struct `SecuredVMDetails`
- New struct `TargetDiskNetworkAccessSettings`
- New field `IncludeSoftDeletedRP` in struct `BMSRPQueryObject`
- New field `IsPrivateAccessEnabledOnAnyDisk` in struct `IaasVMRecoveryPoint`
- New field `SecurityType` in struct `IaasVMRecoveryPoint`
- New field `ExtendedLocation` in struct `IaasVMRestoreRequest`
- New field `SecuredVMDetails` in struct `IaasVMRestoreRequest`
- New field `TargetDiskNetworkAccessSettings` in struct `IaasVMRestoreRequest`
- New field `ExtendedLocation` in struct `IaasVMRestoreWithRehydrationRequest`
- New field `SecuredVMDetails` in struct `IaasVMRestoreWithRehydrationRequest`
- New field `TargetDiskNetworkAccessSettings` in struct `IaasVMRestoreWithRehydrationRequest`
- New field `IsSoftDeleted` in struct `RecoveryPointProperties`


## 2.0.0 (2023-01-19)
### Breaking Changes

- Type of `AzureBackupServerContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `AzureIaaSClassicComputeVMContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `AzureIaaSComputeVMContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `AzureSQLAGWorkloadContainerProtectionContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `AzureSQLContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `AzureStorageContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `AzureStorageProtectableContainer.ProtectableContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `AzureVMAppContainerProtectableContainer.ProtectableContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `AzureVMAppContainerProtectionContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `AzureWorkloadContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `DpmContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `GenericContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `IaaSVMContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `MabContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `ProtectableContainer.ProtectableContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Type of `ProtectionContainer.ContainerType` has been changed from `*ContainerType` to `*ProtectableContainerType`
- Const `ContainerTypeAzureWorkloadContainer`, `ContainerTypeMicrosoftClassicComputeVirtualMachines`, `ContainerTypeMicrosoftComputeVirtualMachines` from type alias `ContainerType` has been removed

### Features Added

- New value `BackupItemTypeSAPHanaDBInstance` added to type alias `BackupItemType`
- New value `BackupTypeSnapshotCopyOnlyFull`, `BackupTypeSnapshotFull` added to type alias `BackupType`
- New value `ContainerTypeHanaHSRContainer` added to type alias `ContainerType`
- New value `DataSourceTypeSAPHanaDBInstance` added to type alias `DataSourceType`
- New value `PolicyTypeSnapshotCopyOnlyFull`, `PolicyTypeSnapshotFull` added to type alias `PolicyType`
- New value `ProtectedItemStateBackupsSuspended` added to type alias `ProtectedItemState`
- New value `ProtectionStateBackupsSuspended` added to type alias `ProtectionState`
- New value `RestorePointQueryTypeSnapshotCopyOnlyFull`, `RestorePointQueryTypeSnapshotFull` added to type alias `RestorePointQueryType`
- New value `RestorePointTypeSnapshotCopyOnlyFull`, `RestorePointTypeSnapshotFull` added to type alias `RestorePointType`
- New value `WorkloadItemTypeSAPHanaDBInstance` added to type alias `WorkloadItemType`
- New value `WorkloadTypeSAPHanaDBInstance` added to type alias `WorkloadType`
- New type alias `ProtectableContainerType` with values `ProtectableContainerTypeAzureBackupServerContainer`, `ProtectableContainerTypeAzureSQLContainer`, `ProtectableContainerTypeAzureWorkloadContainer`, `ProtectableContainerTypeCluster`, `ProtectableContainerTypeDPMContainer`, `ProtectableContainerTypeGenericContainer`, `ProtectableContainerTypeIaasVMContainer`, `ProtectableContainerTypeIaasVMServiceContainer`, `ProtectableContainerTypeInvalid`, `ProtectableContainerTypeMABContainer`, `ProtectableContainerTypeMicrosoftClassicComputeVirtualMachines`, `ProtectableContainerTypeMicrosoftComputeVirtualMachines`, `ProtectableContainerTypeSQLAGWorkLoadContainer`, `ProtectableContainerTypeStorageContainer`, `ProtectableContainerTypeUnknown`, `ProtectableContainerTypeVCenter`, `ProtectableContainerTypeVMAppContainer`, `ProtectableContainerTypeWindows`
- New type alias `TieringMode` with values `TieringModeDoNotTier`, `TieringModeInvalid`, `TieringModeTierAfter`, `TieringModeTierRecommended`
- New function `*AzureVMWorkloadSAPHanaDBInstance.GetAzureVMWorkloadProtectableItem() *AzureVMWorkloadProtectableItem`
- New function `*AzureVMWorkloadSAPHanaDBInstance.GetWorkloadProtectableItem() *WorkloadProtectableItem`
- New function `*AzureVMWorkloadSAPHanaDBInstanceProtectedItem.GetAzureVMWorkloadProtectedItem() *AzureVMWorkloadProtectedItem`
- New function `*AzureVMWorkloadSAPHanaDBInstanceProtectedItem.GetProtectedItem() *ProtectedItem`
- New function `*AzureVMWorkloadSAPHanaHSR.GetAzureVMWorkloadProtectableItem() *AzureVMWorkloadProtectableItem`
- New function `*AzureVMWorkloadSAPHanaHSR.GetWorkloadProtectableItem() *WorkloadProtectableItem`
- New function `NewDeletedProtectionContainersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DeletedProtectionContainersClient, error)`
- New function `*DeletedProtectionContainersClient.NewListPager(string, string, *DeletedProtectionContainersClientListOptions) *runtime.Pager[DeletedProtectionContainersClientListResponse]`
- New struct `AzureVMWorkloadSAPHanaDBInstance`
- New struct `AzureVMWorkloadSAPHanaDBInstanceProtectedItem`
- New struct `AzureVMWorkloadSAPHanaHSR`
- New struct `DeletedProtectionContainersClient`
- New struct `DeletedProtectionContainersClientListResponse`
- New struct `RecoveryPointProperties`
- New struct `TieringPolicy`
- New field `RecoveryPointProperties` in struct `AzureFileShareRecoveryPoint`
- New field `SoftDeleteRetentionPeriod` in struct `AzureFileshareProtectedItem`
- New field `SoftDeleteRetentionPeriod` in struct `AzureIaaSClassicComputeVMProtectedItem`
- New field `SoftDeleteRetentionPeriod` in struct `AzureIaaSComputeVMProtectedItem`
- New field `SoftDeleteRetentionPeriod` in struct `AzureIaaSVMProtectedItem`
- New field `NewestRecoveryPointInArchive` in struct `AzureIaaSVMProtectedItemExtendedInfo`
- New field `OldestRecoveryPointInArchive` in struct `AzureIaaSVMProtectedItemExtendedInfo`
- New field `OldestRecoveryPointInVault` in struct `AzureIaaSVMProtectedItemExtendedInfo`
- New field `TieringPolicy` in struct `AzureIaaSVMProtectionPolicy`
- New field `SoftDeleteRetentionPeriod` in struct `AzureSQLProtectedItem`
- New field `NewestRecoveryPointInArchive` in struct `AzureVMWorkloadProtectedItemExtendedInfo`
- New field `OldestRecoveryPointInArchive` in struct `AzureVMWorkloadProtectedItemExtendedInfo`
- New field `OldestRecoveryPointInVault` in struct `AzureVMWorkloadProtectedItemExtendedInfo`
- New field `SoftDeleteRetentionPeriod` in struct `AzureVMWorkloadSAPAseDatabaseProtectedItem`
- New field `SoftDeleteRetentionPeriod` in struct `AzureVMWorkloadSAPHanaDatabaseProtectedItem`
- New field `SoftDeleteRetentionPeriod` in struct `AzureVMWorkloadSQLDatabaseProtectedItem`
- New field `RecoveryPointProperties` in struct `AzureWorkloadPointInTimeRecoveryPoint`
- New field `RecoveryPointProperties` in struct `AzureWorkloadRecoveryPoint`
- New field `RecoveryPointProperties` in struct `AzureWorkloadSAPHanaPointInTimeRecoveryPoint`
- New field `RecoveryPointProperties` in struct `AzureWorkloadSAPHanaRecoveryPoint`
- New field `RecoveryPointProperties` in struct `AzureWorkloadSQLPointInTimeRecoveryPoint`
- New field `RecoveryPointProperties` in struct `AzureWorkloadSQLRecoveryPoint`
- New field `SoftDeleteRetentionPeriod` in struct `DPMProtectedItem`
- New field `SoftDeleteRetentionPeriod` in struct `GenericProtectedItem`
- New field `RecoveryPointProperties` in struct `GenericRecoveryPoint`
- New field `RecoveryPointProperties` in struct `IaasVMRecoveryPoint`
- New field `SoftDeleteRetentionPeriod` in struct `MabFileFolderProtectedItem`
- New field `TieringPolicy` in struct `SubProtectionPolicy`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
