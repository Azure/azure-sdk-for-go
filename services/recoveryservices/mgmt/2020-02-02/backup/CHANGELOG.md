# Unreleased

## Breaking Changes

### Removed Constants

1. AzureFileShareType.XSMB
1. AzureFileShareType.XSync
1. ContainerType.ContainerTypeAzureBackupServerContainer
1. ContainerType.ContainerTypeAzureSQLContainer
1. ContainerType.ContainerTypeCluster
1. ContainerType.ContainerTypeDPMContainer
1. ContainerType.ContainerTypeGenericContainer
1. ContainerType.ContainerTypeIaasVMContainer
1. ContainerType.ContainerTypeIaasVMServiceContainer
1. ContainerType.ContainerTypeInvalid
1. ContainerType.ContainerTypeMABContainer
1. ContainerType.ContainerTypeSQLAGWorkLoadContainer
1. ContainerType.ContainerTypeStorageContainer
1. ContainerType.ContainerTypeUnknown
1. ContainerType.ContainerTypeVCenter
1. ContainerType.ContainerTypeVMAppContainer
1. ContainerType.ContainerTypeWindows
1. ContainerTypeBasicProtectionContainer.ContainerTypeAzureBackupServerContainer1
1. ContainerTypeBasicProtectionContainer.ContainerTypeAzureSQLContainer1
1. ContainerTypeBasicProtectionContainer.ContainerTypeAzureWorkloadContainer
1. ContainerTypeBasicProtectionContainer.ContainerTypeDPMContainer1
1. ContainerTypeBasicProtectionContainer.ContainerTypeGenericContainer1
1. ContainerTypeBasicProtectionContainer.ContainerTypeIaaSVMContainer
1. ContainerTypeBasicProtectionContainer.ContainerTypeMicrosoftClassicComputevirtualMachines
1. ContainerTypeBasicProtectionContainer.ContainerTypeMicrosoftComputevirtualMachines
1. ContainerTypeBasicProtectionContainer.ContainerTypeProtectionContainer
1. ContainerTypeBasicProtectionContainer.ContainerTypeSQLAGWorkLoadContainer1
1. ContainerTypeBasicProtectionContainer.ContainerTypeStorageContainer1
1. ContainerTypeBasicProtectionContainer.ContainerTypeVMAppContainer1
1. ContainerTypeBasicProtectionContainer.ContainerTypeWindows1
1. CopyOptions.CopyOptionsCreateCopy
1. CopyOptions.CopyOptionsFailOnConflict
1. CopyOptions.CopyOptionsInvalid
1. CopyOptions.CopyOptionsOverwrite
1. CopyOptions.CopyOptionsSkip
1. CreateMode.CreateModeDefault
1. CreateMode.CreateModeInvalid
1. CreateMode.CreateModeRecover
1. DataSourceType.DataSourceTypeAzureFileShare
1. DataSourceType.DataSourceTypeAzureSQLDb
1. DataSourceType.DataSourceTypeClient
1. DataSourceType.DataSourceTypeExchange
1. DataSourceType.DataSourceTypeFileFolder
1. DataSourceType.DataSourceTypeGenericDataSource
1. DataSourceType.DataSourceTypeInvalid
1. DataSourceType.DataSourceTypeSAPAseDatabase
1. DataSourceType.DataSourceTypeSAPHanaDatabase
1. DataSourceType.DataSourceTypeSQLDB
1. DataSourceType.DataSourceTypeSQLDataBase
1. DataSourceType.DataSourceTypeSharepoint
1. DataSourceType.DataSourceTypeSystemState
1. DataSourceType.DataSourceTypeVM
1. DataSourceType.DataSourceTypeVMwareVM
1. DayOfWeek.Friday
1. DayOfWeek.Monday
1. DayOfWeek.Saturday
1. DayOfWeek.Sunday
1. DayOfWeek.Thursday
1. DayOfWeek.Tuesday
1. DayOfWeek.Wednesday
1. EngineType.BackupEngineTypeAzureBackupServerEngine
1. EngineType.BackupEngineTypeBackupEngineBase
1. EngineType.BackupEngineTypeDpmBackupEngine
1. EnhancedSecurityState.EnhancedSecurityStateDisabled
1. EnhancedSecurityState.EnhancedSecurityStateEnabled
1. EnhancedSecurityState.EnhancedSecurityStateInvalid
1. FabricName.FabricNameAzure
1. FabricName.FabricNameInvalid
1. FeatureType.FeatureTypeAzureBackupGoals
1. FeatureType.FeatureTypeAzureVMResourceBackup
1. FeatureType.FeatureTypeFeatureSupportRequest
1. HTTPStatusCode.Accepted
1. HTTPStatusCode.Ambiguous
1. HTTPStatusCode.BadGateway
1. HTTPStatusCode.BadRequest
1. HTTPStatusCode.Conflict
1. HTTPStatusCode.Continue
1. HTTPStatusCode.Created
1. HTTPStatusCode.ExpectationFailed
1. HTTPStatusCode.Forbidden
1. HTTPStatusCode.Found
1. HTTPStatusCode.GatewayTimeout
1. HTTPStatusCode.Gone
1. HTTPStatusCode.HTTPVersionNotSupported
1. HTTPStatusCode.InternalServerError
1. HTTPStatusCode.LengthRequired
1. HTTPStatusCode.MethodNotAllowed
1. HTTPStatusCode.Moved
1. HTTPStatusCode.MovedPermanently
1. HTTPStatusCode.MultipleChoices
1. HTTPStatusCode.NoContent
1. HTTPStatusCode.NonAuthoritativeInformation
1. HTTPStatusCode.NotAcceptable
1. HTTPStatusCode.NotFound
1. HTTPStatusCode.NotImplemented
1. HTTPStatusCode.NotModified
1. HTTPStatusCode.OK
1. HTTPStatusCode.PartialContent
1. HTTPStatusCode.PaymentRequired
1. HTTPStatusCode.PreconditionFailed
1. HTTPStatusCode.ProxyAuthenticationRequired
1. HTTPStatusCode.Redirect
1. HTTPStatusCode.RedirectKeepVerb
1. HTTPStatusCode.RedirectMethod
1. HTTPStatusCode.RequestEntityTooLarge
1. HTTPStatusCode.RequestTimeout
1. HTTPStatusCode.RequestURITooLong
1. HTTPStatusCode.RequestedRangeNotSatisfiable
1. HTTPStatusCode.ResetContent
1. HTTPStatusCode.SeeOther
1. HTTPStatusCode.ServiceUnavailable
1. HTTPStatusCode.SwitchingProtocols
1. HTTPStatusCode.TemporaryRedirect
1. HTTPStatusCode.Unauthorized
1. HTTPStatusCode.UnsupportedMediaType
1. HTTPStatusCode.Unused
1. HTTPStatusCode.UpgradeRequired
1. HTTPStatusCode.UseProxy
1. HealthState.HealthStateActionRequired
1. HealthState.HealthStateActionSuggested
1. HealthState.HealthStateInvalid
1. HealthState.HealthStatePassed
1. HealthStatus.HealthStatusActionRequired
1. HealthStatus.HealthStatusActionSuggested
1. HealthStatus.HealthStatusInvalid
1. HealthStatus.HealthStatusPassed
1. InquiryStatus.InquiryStatusFailed
1. InquiryStatus.InquiryStatusInvalid
1. InquiryStatus.InquiryStatusSuccess
1. IntentItemType.IntentItemTypeInvalid
1. IntentItemType.IntentItemTypeSQLAvailabilityGroupContainer
1. IntentItemType.IntentItemTypeSQLInstance
1. ItemType.ItemTypeAzureFileShare
1. ItemType.ItemTypeAzureSQLDb
1. ItemType.ItemTypeClient
1. ItemType.ItemTypeExchange
1. ItemType.ItemTypeFileFolder
1. ItemType.ItemTypeGenericDataSource
1. ItemType.ItemTypeInvalid
1. ItemType.ItemTypeSAPAseDatabase
1. ItemType.ItemTypeSAPHanaDatabase
1. ItemType.ItemTypeSQLDB
1. ItemType.ItemTypeSQLDataBase
1. ItemType.ItemTypeSharepoint
1. ItemType.ItemTypeSystemState
1. ItemType.ItemTypeVM
1. ItemType.ItemTypeVMwareVM
1. JobOperationType.JobOperationTypeBackup
1. JobOperationType.JobOperationTypeConfigureBackup
1. JobOperationType.JobOperationTypeCrossRegionRestore
1. JobOperationType.JobOperationTypeDeleteBackupData
1. JobOperationType.JobOperationTypeDisableBackup
1. JobOperationType.JobOperationTypeInvalid
1. JobOperationType.JobOperationTypeRegister
1. JobOperationType.JobOperationTypeRestore
1. JobOperationType.JobOperationTypeUnRegister
1. JobOperationType.JobOperationTypeUndelete
1. JobStatus.JobStatusCancelled
1. JobStatus.JobStatusCancelling
1. JobStatus.JobStatusCompleted
1. JobStatus.JobStatusCompletedWithWarnings
1. JobStatus.JobStatusFailed
1. JobStatus.JobStatusInProgress
1. JobStatus.JobStatusInvalid
1. JobSupportedAction.JobSupportedActionCancellable
1. JobSupportedAction.JobSupportedActionInvalid
1. JobSupportedAction.JobSupportedActionRetriable
1. JobType.JobTypeAzureIaaSVMJob
1. JobType.JobTypeAzureStorageJob
1. JobType.JobTypeAzureWorkloadJob
1. JobType.JobTypeDpmJob
1. JobType.JobTypeJob
1. JobType.JobTypeMabJob
1. LastBackupStatus.LastBackupStatusHealthy
1. LastBackupStatus.LastBackupStatusIRPending
1. LastBackupStatus.LastBackupStatusInvalid
1. LastBackupStatus.LastBackupStatusUnhealthy
1. MabServerType.MabServerTypeAzureBackupServerContainer
1. MabServerType.MabServerTypeAzureSQLContainer
1. MabServerType.MabServerTypeCluster
1. MabServerType.MabServerTypeDPMContainer
1. MabServerType.MabServerTypeGenericContainer
1. MabServerType.MabServerTypeIaasVMContainer
1. MabServerType.MabServerTypeIaasVMServiceContainer
1. MabServerType.MabServerTypeInvalid
1. MabServerType.MabServerTypeMABContainer
1. MabServerType.MabServerTypeSQLAGWorkLoadContainer
1. MabServerType.MabServerTypeStorageContainer
1. MabServerType.MabServerTypeUnknown
1. MabServerType.MabServerTypeVCenter
1. MabServerType.MabServerTypeVMAppContainer
1. MabServerType.MabServerTypeWindows
1. ManagementType.ManagementTypeAzureBackupServer
1. ManagementType.ManagementTypeAzureIaasVM
1. ManagementType.ManagementTypeAzureSQL
1. ManagementType.ManagementTypeAzureStorage
1. ManagementType.ManagementTypeAzureWorkload
1. ManagementType.ManagementTypeDPM
1. ManagementType.ManagementTypeDefaultBackup
1. ManagementType.ManagementTypeInvalid
1. ManagementType.ManagementTypeMAB
1. ManagementTypeBasicProtectionPolicy.BackupManagementTypeAzureIaasVM
1. ManagementTypeBasicProtectionPolicy.BackupManagementTypeAzureSQL
1. ManagementTypeBasicProtectionPolicy.BackupManagementTypeAzureStorage
1. ManagementTypeBasicProtectionPolicy.BackupManagementTypeAzureWorkload
1. ManagementTypeBasicProtectionPolicy.BackupManagementTypeGenericProtectionPolicy
1. ManagementTypeBasicProtectionPolicy.BackupManagementTypeMAB
1. ManagementTypeBasicProtectionPolicy.BackupManagementTypeProtectionPolicy
1. MonthOfYear.MonthOfYearApril
1. MonthOfYear.MonthOfYearAugust
1. MonthOfYear.MonthOfYearDecember
1. MonthOfYear.MonthOfYearFebruary
1. MonthOfYear.MonthOfYearInvalid
1. MonthOfYear.MonthOfYearJanuary
1. MonthOfYear.MonthOfYearJuly
1. MonthOfYear.MonthOfYearJune
1. MonthOfYear.MonthOfYearMarch
1. MonthOfYear.MonthOfYearMay
1. MonthOfYear.MonthOfYearNovember
1. MonthOfYear.MonthOfYearOctober
1. MonthOfYear.MonthOfYearSeptember
1. ObjectType.ObjectTypeOperationStatusExtendedInfo
1. ObjectType.ObjectTypeOperationStatusJobExtendedInfo
1. ObjectType.ObjectTypeOperationStatusJobsExtendedInfo
1. ObjectType.ObjectTypeOperationStatusProvisionILRExtendedInfo
1. ObjectTypeBasicILRRequest.ObjectTypeAzureFileShareProvisionILRRequest
1. ObjectTypeBasicILRRequest.ObjectTypeILRRequest
1. ObjectTypeBasicILRRequest.ObjectTypeIaasVMILRRegistrationRequest
1. ObjectTypeBasicOperationResultInfoBase.ObjectTypeExportJobsOperationResultInfo
1. ObjectTypeBasicOperationResultInfoBase.ObjectTypeOperationResultInfo
1. ObjectTypeBasicOperationResultInfoBase.ObjectTypeOperationResultInfoBase
1. ObjectTypeBasicRecoveryPoint.ObjectTypeAzureFileShareRecoveryPoint
1. ObjectTypeBasicRecoveryPoint.ObjectTypeAzureWorkloadPointInTimeRecoveryPoint
1. ObjectTypeBasicRecoveryPoint.ObjectTypeAzureWorkloadRecoveryPoint
1. ObjectTypeBasicRecoveryPoint.ObjectTypeAzureWorkloadSAPHanaPointInTimeRecoveryPoint
1. ObjectTypeBasicRecoveryPoint.ObjectTypeAzureWorkloadSAPHanaRecoveryPoint
1. ObjectTypeBasicRecoveryPoint.ObjectTypeAzureWorkloadSQLPointInTimeRecoveryPoint
1. ObjectTypeBasicRecoveryPoint.ObjectTypeAzureWorkloadSQLRecoveryPoint
1. ObjectTypeBasicRecoveryPoint.ObjectTypeGenericRecoveryPoint
1. ObjectTypeBasicRecoveryPoint.ObjectTypeIaasVMRecoveryPoint
1. ObjectTypeBasicRecoveryPoint.ObjectTypeRecoveryPoint
1. ObjectTypeBasicRequest.ObjectTypeAzureFileShareBackupRequest
1. ObjectTypeBasicRequest.ObjectTypeAzureWorkloadBackupRequest
1. ObjectTypeBasicRequest.ObjectTypeBackupRequest
1. ObjectTypeBasicRequest.ObjectTypeIaasVMBackupRequest
1. ObjectTypeBasicRestoreRequest.ObjectTypeAzureFileShareRestoreRequest
1. ObjectTypeBasicRestoreRequest.ObjectTypeAzureWorkloadPointInTimeRestoreRequest
1. ObjectTypeBasicRestoreRequest.ObjectTypeAzureWorkloadRestoreRequest
1. ObjectTypeBasicRestoreRequest.ObjectTypeAzureWorkloadSAPHanaPointInTimeRestoreRequest
1. ObjectTypeBasicRestoreRequest.ObjectTypeAzureWorkloadSAPHanaRestoreRequest
1. ObjectTypeBasicRestoreRequest.ObjectTypeAzureWorkloadSQLPointInTimeRestoreRequest
1. ObjectTypeBasicRestoreRequest.ObjectTypeAzureWorkloadSQLRestoreRequest
1. ObjectTypeBasicRestoreRequest.ObjectTypeIaasVMRestoreRequest
1. ObjectTypeBasicRestoreRequest.ObjectTypeRestoreRequest
1. ObjectTypeBasicValidateOperationRequest.ObjectTypeValidateIaasVMRestoreOperationRequest
1. ObjectTypeBasicValidateOperationRequest.ObjectTypeValidateOperationRequest
1. ObjectTypeBasicValidateOperationRequest.ObjectTypeValidateRestoreOperationRequest
1. OperationStatusValues.OperationStatusValuesCanceled
1. OperationStatusValues.OperationStatusValuesFailed
1. OperationStatusValues.OperationStatusValuesInProgress
1. OperationStatusValues.OperationStatusValuesInvalid
1. OperationStatusValues.OperationStatusValuesSucceeded
1. OperationType.OperationTypeInvalid
1. OperationType.OperationTypeRegister
1. OperationType.OperationTypeReregister
1. OverwriteOptions.OverwriteOptionsFailOnConflict
1. OverwriteOptions.OverwriteOptionsInvalid
1. OverwriteOptions.OverwriteOptionsOverwrite
1. PolicyType.PolicyTypeCopyOnlyFull
1. PolicyType.PolicyTypeDifferential
1. PolicyType.PolicyTypeFull
1. PolicyType.PolicyTypeInvalid
1. PolicyType.PolicyTypeLog
1. ProtectableContainerType.ProtectableContainerTypeProtectableContainer
1. ProtectableContainerType.ProtectableContainerTypeStorageContainer
1. ProtectableContainerType.ProtectableContainerTypeVMAppContainer
1. ProtectableItemType.ProtectableItemTypeAzureFileShare
1. ProtectableItemType.ProtectableItemTypeAzureVMWorkloadProtectableItem
1. ProtectableItemType.ProtectableItemTypeIaaSVMProtectableItem
1. ProtectableItemType.ProtectableItemTypeMicrosoftClassicComputevirtualMachines
1. ProtectableItemType.ProtectableItemTypeMicrosoftComputevirtualMachines
1. ProtectableItemType.ProtectableItemTypeSAPAseSystem
1. ProtectableItemType.ProtectableItemTypeSAPHanaDatabase
1. ProtectableItemType.ProtectableItemTypeSAPHanaSystem
1. ProtectableItemType.ProtectableItemTypeSQLAvailabilityGroupContainer
1. ProtectableItemType.ProtectableItemTypeSQLDataBase
1. ProtectableItemType.ProtectableItemTypeSQLInstance
1. ProtectableItemType.ProtectableItemTypeWorkloadProtectableItem
1. ProtectedItemHealthStatus.ProtectedItemHealthStatusHealthy
1. ProtectedItemHealthStatus.ProtectedItemHealthStatusIRPending
1. ProtectedItemHealthStatus.ProtectedItemHealthStatusInvalid
1. ProtectedItemHealthStatus.ProtectedItemHealthStatusNotReachable
1. ProtectedItemHealthStatus.ProtectedItemHealthStatusUnhealthy
1. ProtectedItemState.ProtectedItemStateIRPending
1. ProtectedItemState.ProtectedItemStateInvalid
1. ProtectedItemState.ProtectedItemStateProtected
1. ProtectedItemState.ProtectedItemStateProtectionError
1. ProtectedItemState.ProtectedItemStateProtectionPaused
1. ProtectedItemState.ProtectedItemStateProtectionStopped
1. ProtectedItemType.ProtectedItemTypeAzureFileShareProtectedItem
1. ProtectedItemType.ProtectedItemTypeAzureIaaSVMProtectedItem
1. ProtectedItemType.ProtectedItemTypeAzureVMWorkloadProtectedItem
1. ProtectedItemType.ProtectedItemTypeAzureVMWorkloadSAPAseDatabase
1. ProtectedItemType.ProtectedItemTypeAzureVMWorkloadSAPHanaDatabase
1. ProtectedItemType.ProtectedItemTypeAzureVMWorkloadSQLDatabase
1. ProtectedItemType.ProtectedItemTypeDPMProtectedItem
1. ProtectedItemType.ProtectedItemTypeGenericProtectedItem
1. ProtectedItemType.ProtectedItemTypeMabFileFolderProtectedItem
1. ProtectedItemType.ProtectedItemTypeMicrosoftClassicComputevirtualMachines
1. ProtectedItemType.ProtectedItemTypeMicrosoftComputevirtualMachines
1. ProtectedItemType.ProtectedItemTypeMicrosoftSqlserversdatabases
1. ProtectedItemType.ProtectedItemTypeProtectedItem
1. ProtectionIntentItemType.ProtectionIntentItemTypeAzureResourceItem
1. ProtectionIntentItemType.ProtectionIntentItemTypeAzureWorkloadAutoProtectionIntent
1. ProtectionIntentItemType.ProtectionIntentItemTypeAzureWorkloadSQLAutoProtectionIntent
1. ProtectionIntentItemType.ProtectionIntentItemTypeProtectionIntent
1. ProtectionIntentItemType.ProtectionIntentItemTypeRecoveryServiceVaultItem
1. ProtectionState.ProtectionStateIRPending
1. ProtectionState.ProtectionStateInvalid
1. ProtectionState.ProtectionStateProtected
1. ProtectionState.ProtectionStateProtectionError
1. ProtectionState.ProtectionStateProtectionPaused
1. ProtectionState.ProtectionStateProtectionStopped
1. ProtectionStatus.ProtectionStatusInvalid
1. ProtectionStatus.ProtectionStatusNotProtected
1. ProtectionStatus.ProtectionStatusProtected
1. ProtectionStatus.ProtectionStatusProtecting
1. ProtectionStatus.ProtectionStatusProtectionFailed
1. RecoveryMode.RecoveryModeFileRecovery
1. RecoveryMode.RecoveryModeInvalid
1. RecoveryMode.RecoveryModeWorkloadRecovery
1. RecoveryPointTierStatus.RecoveryPointTierStatusDeleted
1. RecoveryPointTierStatus.RecoveryPointTierStatusDisabled
1. RecoveryPointTierStatus.RecoveryPointTierStatusInvalid
1. RecoveryPointTierStatus.RecoveryPointTierStatusValid
1. RecoveryPointTierType.RecoveryPointTierTypeHardenedRP
1. RecoveryPointTierType.RecoveryPointTierTypeInstantRP
1. RecoveryPointTierType.RecoveryPointTierTypeInvalid
1. RecoveryType.RecoveryTypeAlternateLocation
1. RecoveryType.RecoveryTypeInvalid
1. RecoveryType.RecoveryTypeOffline
1. RecoveryType.RecoveryTypeOriginalLocation
1. RecoveryType.RecoveryTypeRestoreDisks
1. ResourceHealthStatus.ResourceHealthStatusHealthy
1. ResourceHealthStatus.ResourceHealthStatusInvalid
1. ResourceHealthStatus.ResourceHealthStatusPersistentDegraded
1. ResourceHealthStatus.ResourceHealthStatusPersistentUnhealthy
1. ResourceHealthStatus.ResourceHealthStatusTransientDegraded
1. ResourceHealthStatus.ResourceHealthStatusTransientUnhealthy
1. RestorePointQueryType.RestorePointQueryTypeAll
1. RestorePointQueryType.RestorePointQueryTypeDifferential
1. RestorePointQueryType.RestorePointQueryTypeFull
1. RestorePointQueryType.RestorePointQueryTypeFullAndDifferential
1. RestorePointQueryType.RestorePointQueryTypeInvalid
1. RestorePointQueryType.RestorePointQueryTypeLog
1. RestorePointType.RestorePointTypeDifferential
1. RestorePointType.RestorePointTypeFull
1. RestorePointType.RestorePointTypeInvalid
1. RestorePointType.RestorePointTypeLog
1. RestoreRequestType.RestoreRequestTypeFullShareRestore
1. RestoreRequestType.RestoreRequestTypeInvalid
1. RestoreRequestType.RestoreRequestTypeItemLevelRestore
1. RetentionDurationType.RetentionDurationTypeDays
1. RetentionDurationType.RetentionDurationTypeInvalid
1. RetentionDurationType.RetentionDurationTypeMonths
1. RetentionDurationType.RetentionDurationTypeWeeks
1. RetentionDurationType.RetentionDurationTypeYears
1. RetentionPolicyType.RetentionPolicyTypeLongTermRetentionPolicy
1. RetentionPolicyType.RetentionPolicyTypeRetentionPolicy
1. RetentionPolicyType.RetentionPolicyTypeSimpleRetentionPolicy
1. RetentionScheduleFormat.RetentionScheduleFormatDaily
1. RetentionScheduleFormat.RetentionScheduleFormatInvalid
1. RetentionScheduleFormat.RetentionScheduleFormatWeekly
1. SQLDataDirectoryType.SQLDataDirectoryTypeData
1. SQLDataDirectoryType.SQLDataDirectoryTypeInvalid
1. SQLDataDirectoryType.SQLDataDirectoryTypeLog
1. SchedulePolicyType.SchedulePolicyTypeLogSchedulePolicy
1. SchedulePolicyType.SchedulePolicyTypeLongTermSchedulePolicy
1. SchedulePolicyType.SchedulePolicyTypeSchedulePolicy
1. SchedulePolicyType.SchedulePolicyTypeSimpleSchedulePolicy
1. ScheduleRunType.ScheduleRunTypeDaily
1. ScheduleRunType.ScheduleRunTypeInvalid
1. ScheduleRunType.ScheduleRunTypeWeekly
1. SoftDeleteFeatureState.SoftDeleteFeatureStateDisabled
1. SoftDeleteFeatureState.SoftDeleteFeatureStateEnabled
1. SoftDeleteFeatureState.SoftDeleteFeatureStateInvalid
1. StorageType.StorageTypeGeoRedundant
1. StorageType.StorageTypeInvalid
1. StorageType.StorageTypeLocallyRedundant
1. StorageTypeState.StorageTypeStateInvalid
1. StorageTypeState.StorageTypeStateLocked
1. StorageTypeState.StorageTypeStateUnlocked
1. SupportStatus.SupportStatusDefaultOFF
1. SupportStatus.SupportStatusDefaultON
1. SupportStatus.SupportStatusInvalid
1. SupportStatus.SupportStatusNotSupported
1. SupportStatus.SupportStatusSupported
1. Type.TypeBackupProtectedItemCountSummary
1. Type.TypeBackupProtectionContainerCountSummary
1. Type.TypeInvalid
1. TypeEnum.TypeEnumCopyOnlyFull
1. TypeEnum.TypeEnumDifferential
1. TypeEnum.TypeEnumFull
1. TypeEnum.TypeEnumInvalid
1. TypeEnum.TypeEnumLog
1. UsagesUnit.Bytes
1. UsagesUnit.BytesPerSecond
1. UsagesUnit.Count
1. UsagesUnit.CountPerSecond
1. UsagesUnit.Percent
1. UsagesUnit.Seconds
1. ValidationStatus.ValidationStatusFailed
1. ValidationStatus.ValidationStatusInvalid
1. ValidationStatus.ValidationStatusSucceeded
1. WeekOfMonth.WeekOfMonthFirst
1. WeekOfMonth.WeekOfMonthFourth
1. WeekOfMonth.WeekOfMonthInvalid
1. WeekOfMonth.WeekOfMonthLast
1. WeekOfMonth.WeekOfMonthSecond
1. WeekOfMonth.WeekOfMonthThird
1. WorkloadItemType.WorkloadItemTypeInvalid
1. WorkloadItemType.WorkloadItemTypeSAPAseDatabase
1. WorkloadItemType.WorkloadItemTypeSAPAseSystem
1. WorkloadItemType.WorkloadItemTypeSAPHanaDatabase
1. WorkloadItemType.WorkloadItemTypeSAPHanaSystem
1. WorkloadItemType.WorkloadItemTypeSQLDataBase
1. WorkloadItemType.WorkloadItemTypeSQLInstance
1. WorkloadItemTypeBasicWorkloadItem.WorkloadItemTypeAzureVMWorkloadItem
1. WorkloadItemTypeBasicWorkloadItem.WorkloadItemTypeSAPAseDatabase1
1. WorkloadItemTypeBasicWorkloadItem.WorkloadItemTypeSAPAseSystem1
1. WorkloadItemTypeBasicWorkloadItem.WorkloadItemTypeSAPHanaDatabase1
1. WorkloadItemTypeBasicWorkloadItem.WorkloadItemTypeSAPHanaSystem1
1. WorkloadItemTypeBasicWorkloadItem.WorkloadItemTypeSQLDataBase1
1. WorkloadItemTypeBasicWorkloadItem.WorkloadItemTypeSQLInstance1
1. WorkloadItemTypeBasicWorkloadItem.WorkloadItemTypeWorkloadItem
1. WorkloadType.WorkloadTypeAzureFileShare
1. WorkloadType.WorkloadTypeAzureSQLDb
1. WorkloadType.WorkloadTypeClient
1. WorkloadType.WorkloadTypeExchange
1. WorkloadType.WorkloadTypeFileFolder
1. WorkloadType.WorkloadTypeGenericDataSource
1. WorkloadType.WorkloadTypeInvalid
1. WorkloadType.WorkloadTypeSAPAseDatabase
1. WorkloadType.WorkloadTypeSAPHanaDatabase
1. WorkloadType.WorkloadTypeSQLDB
1. WorkloadType.WorkloadTypeSQLDataBase
1. WorkloadType.WorkloadTypeSharepoint
1. WorkloadType.WorkloadTypeSystemState
1. WorkloadType.WorkloadTypeVM
1. WorkloadType.WorkloadTypeVMwareVM

### Removed Funcs

1. *AzureFileShareProtectionPolicy.UnmarshalJSON([]byte) error
1. *AzureIaaSVMProtectionPolicy.UnmarshalJSON([]byte) error
1. *AzureSQLProtectionPolicy.UnmarshalJSON([]byte) error
1. *ClientDiscoveryResponseIterator.Next() error
1. *ClientDiscoveryResponseIterator.NextWithContext(context.Context) error
1. *ClientDiscoveryResponsePage.Next() error
1. *ClientDiscoveryResponsePage.NextWithContext(context.Context) error
1. *EngineBaseResource.UnmarshalJSON([]byte) error
1. *EngineBaseResourceListIterator.Next() error
1. *EngineBaseResourceListIterator.NextWithContext(context.Context) error
1. *EngineBaseResourceListPage.Next() error
1. *EngineBaseResourceListPage.NextWithContext(context.Context) error
1. *ILRRequestResource.UnmarshalJSON([]byte) error
1. *JobResource.UnmarshalJSON([]byte) error
1. *JobResourceListIterator.Next() error
1. *JobResourceListIterator.NextWithContext(context.Context) error
1. *JobResourceListPage.Next() error
1. *JobResourceListPage.NextWithContext(context.Context) error
1. *MabProtectionPolicy.UnmarshalJSON([]byte) error
1. *OperationResultInfoBaseResource.UnmarshalJSON([]byte) error
1. *OperationStatus.UnmarshalJSON([]byte) error
1. *ProtectableContainerResource.UnmarshalJSON([]byte) error
1. *ProtectableContainerResourceListIterator.Next() error
1. *ProtectableContainerResourceListIterator.NextWithContext(context.Context) error
1. *ProtectableContainerResourceListPage.Next() error
1. *ProtectableContainerResourceListPage.NextWithContext(context.Context) error
1. *ProtectedItemResource.UnmarshalJSON([]byte) error
1. *ProtectedItemResourceListIterator.Next() error
1. *ProtectedItemResourceListIterator.NextWithContext(context.Context) error
1. *ProtectedItemResourceListPage.Next() error
1. *ProtectedItemResourceListPage.NextWithContext(context.Context) error
1. *ProtectionContainerResource.UnmarshalJSON([]byte) error
1. *ProtectionContainerResourceListIterator.Next() error
1. *ProtectionContainerResourceListIterator.NextWithContext(context.Context) error
1. *ProtectionContainerResourceListPage.Next() error
1. *ProtectionContainerResourceListPage.NextWithContext(context.Context) error
1. *ProtectionIntentResource.UnmarshalJSON([]byte) error
1. *ProtectionIntentResourceListIterator.Next() error
1. *ProtectionIntentResourceListIterator.NextWithContext(context.Context) error
1. *ProtectionIntentResourceListPage.Next() error
1. *ProtectionIntentResourceListPage.NextWithContext(context.Context) error
1. *ProtectionPolicyResource.UnmarshalJSON([]byte) error
1. *ProtectionPolicyResourceListIterator.Next() error
1. *ProtectionPolicyResourceListIterator.NextWithContext(context.Context) error
1. *ProtectionPolicyResourceListPage.Next() error
1. *ProtectionPolicyResourceListPage.NextWithContext(context.Context) error
1. *RecoveryPointResource.UnmarshalJSON([]byte) error
1. *RecoveryPointResourceListIterator.Next() error
1. *RecoveryPointResourceListIterator.NextWithContext(context.Context) error
1. *RecoveryPointResourceListPage.Next() error
1. *RecoveryPointResourceListPage.NextWithContext(context.Context) error
1. *RequestResource.UnmarshalJSON([]byte) error
1. *RestoreRequestResource.UnmarshalJSON([]byte) error
1. *SubProtectionPolicy.UnmarshalJSON([]byte) error
1. *ValidateIaasVMRestoreOperationRequest.UnmarshalJSON([]byte) error
1. *ValidateRestoreOperationRequest.UnmarshalJSON([]byte) error
1. *WorkloadItemResource.UnmarshalJSON([]byte) error
1. *WorkloadItemResourceListIterator.Next() error
1. *WorkloadItemResourceListIterator.NextWithContext(context.Context) error
1. *WorkloadItemResourceListPage.Next() error
1. *WorkloadItemResourceListPage.NextWithContext(context.Context) error
1. *WorkloadProtectableItemResource.UnmarshalJSON([]byte) error
1. *WorkloadProtectableItemResourceListIterator.Next() error
1. *WorkloadProtectableItemResourceListIterator.NextWithContext(context.Context) error
1. *WorkloadProtectableItemResourceListPage.Next() error
1. *WorkloadProtectableItemResourceListPage.NextWithContext(context.Context) error
1. AzureBackupGoalFeatureSupportRequest.AsAzureBackupGoalFeatureSupportRequest() (*AzureBackupGoalFeatureSupportRequest, bool)
1. AzureBackupGoalFeatureSupportRequest.AsAzureVMResourceFeatureSupportRequest() (*AzureVMResourceFeatureSupportRequest, bool)
1. AzureBackupGoalFeatureSupportRequest.AsBasicFeatureSupportRequest() (BasicFeatureSupportRequest, bool)
1. AzureBackupGoalFeatureSupportRequest.AsFeatureSupportRequest() (*FeatureSupportRequest, bool)
1. AzureBackupGoalFeatureSupportRequest.MarshalJSON() ([]byte, error)
1. AzureBackupServerContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. AzureBackupServerContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. AzureBackupServerContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. AzureBackupServerContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. AzureBackupServerContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. AzureBackupServerContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. AzureBackupServerContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. AzureBackupServerContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. AzureBackupServerContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. AzureBackupServerContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. AzureBackupServerContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. AzureBackupServerContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. AzureBackupServerContainer.AsDpmContainer() (*DpmContainer, bool)
1. AzureBackupServerContainer.AsGenericContainer() (*GenericContainer, bool)
1. AzureBackupServerContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. AzureBackupServerContainer.AsMabContainer() (*MabContainer, bool)
1. AzureBackupServerContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. AzureBackupServerContainer.MarshalJSON() ([]byte, error)
1. AzureBackupServerEngine.AsAzureBackupServerEngine() (*AzureBackupServerEngine, bool)
1. AzureBackupServerEngine.AsBasicEngineBase() (BasicEngineBase, bool)
1. AzureBackupServerEngine.AsDpmBackupEngine() (*DpmBackupEngine, bool)
1. AzureBackupServerEngine.AsEngineBase() (*EngineBase, bool)
1. AzureBackupServerEngine.MarshalJSON() ([]byte, error)
1. AzureFileShareBackupRequest.AsAzureFileShareBackupRequest() (*AzureFileShareBackupRequest, bool)
1. AzureFileShareBackupRequest.AsAzureWorkloadBackupRequest() (*AzureWorkloadBackupRequest, bool)
1. AzureFileShareBackupRequest.AsBasicRequest() (BasicRequest, bool)
1. AzureFileShareBackupRequest.AsIaasVMBackupRequest() (*IaasVMBackupRequest, bool)
1. AzureFileShareBackupRequest.AsRequest() (*Request, bool)
1. AzureFileShareBackupRequest.MarshalJSON() ([]byte, error)
1. AzureFileShareProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. AzureFileShareProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. AzureFileShareProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. AzureFileShareProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. AzureFileShareProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. AzureFileShareProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. AzureFileShareProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. AzureFileShareProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. AzureFileShareProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. AzureFileShareProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. AzureFileShareProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. AzureFileShareProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. AzureFileShareProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. AzureFileShareProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. AzureFileShareProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. AzureFileShareProtectableItem.MarshalJSON() ([]byte, error)
1. AzureFileShareProtectionPolicy.AsAzureFileShareProtectionPolicy() (*AzureFileShareProtectionPolicy, bool)
1. AzureFileShareProtectionPolicy.AsAzureIaaSVMProtectionPolicy() (*AzureIaaSVMProtectionPolicy, bool)
1. AzureFileShareProtectionPolicy.AsAzureSQLProtectionPolicy() (*AzureSQLProtectionPolicy, bool)
1. AzureFileShareProtectionPolicy.AsAzureVMWorkloadProtectionPolicy() (*AzureVMWorkloadProtectionPolicy, bool)
1. AzureFileShareProtectionPolicy.AsBasicProtectionPolicy() (BasicProtectionPolicy, bool)
1. AzureFileShareProtectionPolicy.AsGenericProtectionPolicy() (*GenericProtectionPolicy, bool)
1. AzureFileShareProtectionPolicy.AsMabProtectionPolicy() (*MabProtectionPolicy, bool)
1. AzureFileShareProtectionPolicy.AsProtectionPolicy() (*ProtectionPolicy, bool)
1. AzureFileShareProtectionPolicy.MarshalJSON() ([]byte, error)
1. AzureFileShareProvisionILRRequest.AsAzureFileShareProvisionILRRequest() (*AzureFileShareProvisionILRRequest, bool)
1. AzureFileShareProvisionILRRequest.AsBasicILRRequest() (BasicILRRequest, bool)
1. AzureFileShareProvisionILRRequest.AsILRRequest() (*ILRRequest, bool)
1. AzureFileShareProvisionILRRequest.AsIaasVMILRRegistrationRequest() (*IaasVMILRRegistrationRequest, bool)
1. AzureFileShareProvisionILRRequest.MarshalJSON() ([]byte, error)
1. AzureFileShareRecoveryPoint.AsAzureFileShareRecoveryPoint() (*AzureFileShareRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsAzureWorkloadPointInTimeRecoveryPoint() (*AzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsAzureWorkloadRecoveryPoint() (*AzureWorkloadRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsAzureWorkloadSAPHanaPointInTimeRecoveryPoint() (*AzureWorkloadSAPHanaPointInTimeRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsAzureWorkloadSAPHanaRecoveryPoint() (*AzureWorkloadSAPHanaRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsAzureWorkloadSQLPointInTimeRecoveryPoint() (*AzureWorkloadSQLPointInTimeRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsAzureWorkloadSQLRecoveryPoint() (*AzureWorkloadSQLRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsBasicAzureWorkloadPointInTimeRecoveryPoint() (BasicAzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsBasicAzureWorkloadRecoveryPoint() (BasicAzureWorkloadRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsBasicAzureWorkloadSQLRecoveryPoint() (BasicAzureWorkloadSQLRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsBasicRecoveryPoint() (BasicRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsGenericRecoveryPoint() (*GenericRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsIaasVMRecoveryPoint() (*IaasVMRecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.AsRecoveryPoint() (*RecoveryPoint, bool)
1. AzureFileShareRecoveryPoint.MarshalJSON() ([]byte, error)
1. AzureFileShareRestoreRequest.AsAzureFileShareRestoreRequest() (*AzureFileShareRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsAzureWorkloadPointInTimeRestoreRequest() (*AzureWorkloadPointInTimeRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsAzureWorkloadRestoreRequest() (*AzureWorkloadRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsAzureWorkloadSAPHanaPointInTimeRestoreRequest() (*AzureWorkloadSAPHanaPointInTimeRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsAzureWorkloadSAPHanaRestoreRequest() (*AzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsAzureWorkloadSQLPointInTimeRestoreRequest() (*AzureWorkloadSQLPointInTimeRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsAzureWorkloadSQLRestoreRequest() (*AzureWorkloadSQLRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsBasicAzureWorkloadRestoreRequest() (BasicAzureWorkloadRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsBasicAzureWorkloadSAPHanaRestoreRequest() (BasicAzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsBasicAzureWorkloadSQLRestoreRequest() (BasicAzureWorkloadSQLRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsBasicRestoreRequest() (BasicRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsIaasVMRestoreRequest() (*IaasVMRestoreRequest, bool)
1. AzureFileShareRestoreRequest.AsRestoreRequest() (*RestoreRequest, bool)
1. AzureFileShareRestoreRequest.MarshalJSON() ([]byte, error)
1. AzureFileshareProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. AzureFileshareProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. AzureFileshareProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. AzureFileshareProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. AzureFileshareProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. AzureFileshareProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. AzureFileshareProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. AzureFileshareProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. AzureFileshareProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. AzureFileshareProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. AzureFileshareProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. AzureFileshareProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. AzureFileshareProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. AzureFileshareProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. AzureFileshareProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. AzureFileshareProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. AzureFileshareProtectedItem.MarshalJSON() ([]byte, error)
1. AzureFileshareProtectedItemExtendedInfo.MarshalJSON() ([]byte, error)
1. AzureIaaSClassicComputeVMContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsDpmContainer() (*DpmContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsGenericContainer() (*GenericContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsMabContainer() (*MabContainer, bool)
1. AzureIaaSClassicComputeVMContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. AzureIaaSClassicComputeVMContainer.MarshalJSON() ([]byte, error)
1. AzureIaaSClassicComputeVMProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. AzureIaaSClassicComputeVMProtectableItem.MarshalJSON() ([]byte, error)
1. AzureIaaSClassicComputeVMProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. AzureIaaSClassicComputeVMProtectedItem.MarshalJSON() ([]byte, error)
1. AzureIaaSComputeVMContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. AzureIaaSComputeVMContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. AzureIaaSComputeVMContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. AzureIaaSComputeVMContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. AzureIaaSComputeVMContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. AzureIaaSComputeVMContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. AzureIaaSComputeVMContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. AzureIaaSComputeVMContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. AzureIaaSComputeVMContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. AzureIaaSComputeVMContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. AzureIaaSComputeVMContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. AzureIaaSComputeVMContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. AzureIaaSComputeVMContainer.AsDpmContainer() (*DpmContainer, bool)
1. AzureIaaSComputeVMContainer.AsGenericContainer() (*GenericContainer, bool)
1. AzureIaaSComputeVMContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. AzureIaaSComputeVMContainer.AsMabContainer() (*MabContainer, bool)
1. AzureIaaSComputeVMContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. AzureIaaSComputeVMContainer.MarshalJSON() ([]byte, error)
1. AzureIaaSComputeVMProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. AzureIaaSComputeVMProtectableItem.MarshalJSON() ([]byte, error)
1. AzureIaaSComputeVMProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. AzureIaaSComputeVMProtectedItem.MarshalJSON() ([]byte, error)
1. AzureIaaSVMErrorInfo.MarshalJSON() ([]byte, error)
1. AzureIaaSVMHealthDetails.MarshalJSON() ([]byte, error)
1. AzureIaaSVMJob.AsAzureIaaSVMJob() (*AzureIaaSVMJob, bool)
1. AzureIaaSVMJob.AsAzureStorageJob() (*AzureStorageJob, bool)
1. AzureIaaSVMJob.AsAzureWorkloadJob() (*AzureWorkloadJob, bool)
1. AzureIaaSVMJob.AsBasicJob() (BasicJob, bool)
1. AzureIaaSVMJob.AsDpmJob() (*DpmJob, bool)
1. AzureIaaSVMJob.AsJob() (*Job, bool)
1. AzureIaaSVMJob.AsMabJob() (*MabJob, bool)
1. AzureIaaSVMJob.MarshalJSON() ([]byte, error)
1. AzureIaaSVMJobExtendedInfo.MarshalJSON() ([]byte, error)
1. AzureIaaSVMProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. AzureIaaSVMProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. AzureIaaSVMProtectedItem.MarshalJSON() ([]byte, error)
1. AzureIaaSVMProtectionPolicy.AsAzureFileShareProtectionPolicy() (*AzureFileShareProtectionPolicy, bool)
1. AzureIaaSVMProtectionPolicy.AsAzureIaaSVMProtectionPolicy() (*AzureIaaSVMProtectionPolicy, bool)
1. AzureIaaSVMProtectionPolicy.AsAzureSQLProtectionPolicy() (*AzureSQLProtectionPolicy, bool)
1. AzureIaaSVMProtectionPolicy.AsAzureVMWorkloadProtectionPolicy() (*AzureVMWorkloadProtectionPolicy, bool)
1. AzureIaaSVMProtectionPolicy.AsBasicProtectionPolicy() (BasicProtectionPolicy, bool)
1. AzureIaaSVMProtectionPolicy.AsGenericProtectionPolicy() (*GenericProtectionPolicy, bool)
1. AzureIaaSVMProtectionPolicy.AsMabProtectionPolicy() (*MabProtectionPolicy, bool)
1. AzureIaaSVMProtectionPolicy.AsProtectionPolicy() (*ProtectionPolicy, bool)
1. AzureIaaSVMProtectionPolicy.MarshalJSON() ([]byte, error)
1. AzureRecoveryServiceVaultProtectionIntent.AsAzureRecoveryServiceVaultProtectionIntent() (*AzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureRecoveryServiceVaultProtectionIntent.AsAzureResourceProtectionIntent() (*AzureResourceProtectionIntent, bool)
1. AzureRecoveryServiceVaultProtectionIntent.AsAzureWorkloadAutoProtectionIntent() (*AzureWorkloadAutoProtectionIntent, bool)
1. AzureRecoveryServiceVaultProtectionIntent.AsAzureWorkloadSQLAutoProtectionIntent() (*AzureWorkloadSQLAutoProtectionIntent, bool)
1. AzureRecoveryServiceVaultProtectionIntent.AsBasicAzureRecoveryServiceVaultProtectionIntent() (BasicAzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureRecoveryServiceVaultProtectionIntent.AsBasicAzureWorkloadAutoProtectionIntent() (BasicAzureWorkloadAutoProtectionIntent, bool)
1. AzureRecoveryServiceVaultProtectionIntent.AsBasicProtectionIntent() (BasicProtectionIntent, bool)
1. AzureRecoveryServiceVaultProtectionIntent.AsProtectionIntent() (*ProtectionIntent, bool)
1. AzureRecoveryServiceVaultProtectionIntent.MarshalJSON() ([]byte, error)
1. AzureResourceProtectionIntent.AsAzureRecoveryServiceVaultProtectionIntent() (*AzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureResourceProtectionIntent.AsAzureResourceProtectionIntent() (*AzureResourceProtectionIntent, bool)
1. AzureResourceProtectionIntent.AsAzureWorkloadAutoProtectionIntent() (*AzureWorkloadAutoProtectionIntent, bool)
1. AzureResourceProtectionIntent.AsAzureWorkloadSQLAutoProtectionIntent() (*AzureWorkloadSQLAutoProtectionIntent, bool)
1. AzureResourceProtectionIntent.AsBasicAzureRecoveryServiceVaultProtectionIntent() (BasicAzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureResourceProtectionIntent.AsBasicAzureWorkloadAutoProtectionIntent() (BasicAzureWorkloadAutoProtectionIntent, bool)
1. AzureResourceProtectionIntent.AsBasicProtectionIntent() (BasicProtectionIntent, bool)
1. AzureResourceProtectionIntent.AsProtectionIntent() (*ProtectionIntent, bool)
1. AzureResourceProtectionIntent.MarshalJSON() ([]byte, error)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsDpmContainer() (*DpmContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsGenericContainer() (*GenericContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsMabContainer() (*MabContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. AzureSQLAGWorkloadContainerProtectionContainer.MarshalJSON() ([]byte, error)
1. AzureSQLContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. AzureSQLContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. AzureSQLContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. AzureSQLContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. AzureSQLContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. AzureSQLContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. AzureSQLContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. AzureSQLContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. AzureSQLContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. AzureSQLContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. AzureSQLContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. AzureSQLContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. AzureSQLContainer.AsDpmContainer() (*DpmContainer, bool)
1. AzureSQLContainer.AsGenericContainer() (*GenericContainer, bool)
1. AzureSQLContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. AzureSQLContainer.AsMabContainer() (*MabContainer, bool)
1. AzureSQLContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. AzureSQLContainer.MarshalJSON() ([]byte, error)
1. AzureSQLProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. AzureSQLProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. AzureSQLProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. AzureSQLProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. AzureSQLProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. AzureSQLProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. AzureSQLProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. AzureSQLProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. AzureSQLProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. AzureSQLProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. AzureSQLProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. AzureSQLProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. AzureSQLProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. AzureSQLProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. AzureSQLProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. AzureSQLProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. AzureSQLProtectedItem.MarshalJSON() ([]byte, error)
1. AzureSQLProtectionPolicy.AsAzureFileShareProtectionPolicy() (*AzureFileShareProtectionPolicy, bool)
1. AzureSQLProtectionPolicy.AsAzureIaaSVMProtectionPolicy() (*AzureIaaSVMProtectionPolicy, bool)
1. AzureSQLProtectionPolicy.AsAzureSQLProtectionPolicy() (*AzureSQLProtectionPolicy, bool)
1. AzureSQLProtectionPolicy.AsAzureVMWorkloadProtectionPolicy() (*AzureVMWorkloadProtectionPolicy, bool)
1. AzureSQLProtectionPolicy.AsBasicProtectionPolicy() (BasicProtectionPolicy, bool)
1. AzureSQLProtectionPolicy.AsGenericProtectionPolicy() (*GenericProtectionPolicy, bool)
1. AzureSQLProtectionPolicy.AsMabProtectionPolicy() (*MabProtectionPolicy, bool)
1. AzureSQLProtectionPolicy.AsProtectionPolicy() (*ProtectionPolicy, bool)
1. AzureSQLProtectionPolicy.MarshalJSON() ([]byte, error)
1. AzureStorageContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. AzureStorageContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. AzureStorageContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. AzureStorageContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. AzureStorageContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. AzureStorageContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. AzureStorageContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. AzureStorageContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. AzureStorageContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. AzureStorageContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. AzureStorageContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. AzureStorageContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. AzureStorageContainer.AsDpmContainer() (*DpmContainer, bool)
1. AzureStorageContainer.AsGenericContainer() (*GenericContainer, bool)
1. AzureStorageContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. AzureStorageContainer.AsMabContainer() (*MabContainer, bool)
1. AzureStorageContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. AzureStorageContainer.MarshalJSON() ([]byte, error)
1. AzureStorageJob.AsAzureIaaSVMJob() (*AzureIaaSVMJob, bool)
1. AzureStorageJob.AsAzureStorageJob() (*AzureStorageJob, bool)
1. AzureStorageJob.AsAzureWorkloadJob() (*AzureWorkloadJob, bool)
1. AzureStorageJob.AsBasicJob() (BasicJob, bool)
1. AzureStorageJob.AsDpmJob() (*DpmJob, bool)
1. AzureStorageJob.AsJob() (*Job, bool)
1. AzureStorageJob.AsMabJob() (*MabJob, bool)
1. AzureStorageJob.MarshalJSON() ([]byte, error)
1. AzureStorageJobExtendedInfo.MarshalJSON() ([]byte, error)
1. AzureStorageProtectableContainer.AsAzureStorageProtectableContainer() (*AzureStorageProtectableContainer, bool)
1. AzureStorageProtectableContainer.AsAzureVMAppContainerProtectableContainer() (*AzureVMAppContainerProtectableContainer, bool)
1. AzureStorageProtectableContainer.AsBasicProtectableContainer() (BasicProtectableContainer, bool)
1. AzureStorageProtectableContainer.AsProtectableContainer() (*ProtectableContainer, bool)
1. AzureStorageProtectableContainer.MarshalJSON() ([]byte, error)
1. AzureVMAppContainerProtectableContainer.AsAzureStorageProtectableContainer() (*AzureStorageProtectableContainer, bool)
1. AzureVMAppContainerProtectableContainer.AsAzureVMAppContainerProtectableContainer() (*AzureVMAppContainerProtectableContainer, bool)
1. AzureVMAppContainerProtectableContainer.AsBasicProtectableContainer() (BasicProtectableContainer, bool)
1. AzureVMAppContainerProtectableContainer.AsProtectableContainer() (*ProtectableContainer, bool)
1. AzureVMAppContainerProtectableContainer.MarshalJSON() ([]byte, error)
1. AzureVMAppContainerProtectionContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsDpmContainer() (*DpmContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsGenericContainer() (*GenericContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsMabContainer() (*MabContainer, bool)
1. AzureVMAppContainerProtectionContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. AzureVMAppContainerProtectionContainer.MarshalJSON() ([]byte, error)
1. AzureVMResourceFeatureSupportRequest.AsAzureBackupGoalFeatureSupportRequest() (*AzureBackupGoalFeatureSupportRequest, bool)
1. AzureVMResourceFeatureSupportRequest.AsAzureVMResourceFeatureSupportRequest() (*AzureVMResourceFeatureSupportRequest, bool)
1. AzureVMResourceFeatureSupportRequest.AsBasicFeatureSupportRequest() (BasicFeatureSupportRequest, bool)
1. AzureVMResourceFeatureSupportRequest.AsFeatureSupportRequest() (*FeatureSupportRequest, bool)
1. AzureVMResourceFeatureSupportRequest.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadItem.AsAzureVMWorkloadItem() (*AzureVMWorkloadItem, bool)
1. AzureVMWorkloadItem.AsAzureVMWorkloadSAPAseDatabaseWorkloadItem() (*AzureVMWorkloadSAPAseDatabaseWorkloadItem, bool)
1. AzureVMWorkloadItem.AsAzureVMWorkloadSAPAseSystemWorkloadItem() (*AzureVMWorkloadSAPAseSystemWorkloadItem, bool)
1. AzureVMWorkloadItem.AsAzureVMWorkloadSAPHanaDatabaseWorkloadItem() (*AzureVMWorkloadSAPHanaDatabaseWorkloadItem, bool)
1. AzureVMWorkloadItem.AsAzureVMWorkloadSAPHanaSystemWorkloadItem() (*AzureVMWorkloadSAPHanaSystemWorkloadItem, bool)
1. AzureVMWorkloadItem.AsAzureVMWorkloadSQLDatabaseWorkloadItem() (*AzureVMWorkloadSQLDatabaseWorkloadItem, bool)
1. AzureVMWorkloadItem.AsAzureVMWorkloadSQLInstanceWorkloadItem() (*AzureVMWorkloadSQLInstanceWorkloadItem, bool)
1. AzureVMWorkloadItem.AsBasicAzureVMWorkloadItem() (BasicAzureVMWorkloadItem, bool)
1. AzureVMWorkloadItem.AsBasicWorkloadItem() (BasicWorkloadItem, bool)
1. AzureVMWorkloadItem.AsWorkloadItem() (*WorkloadItem, bool)
1. AzureVMWorkloadItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. AzureVMWorkloadProtectableItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. AzureVMWorkloadProtectedItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadProtectionPolicy.AsAzureFileShareProtectionPolicy() (*AzureFileShareProtectionPolicy, bool)
1. AzureVMWorkloadProtectionPolicy.AsAzureIaaSVMProtectionPolicy() (*AzureIaaSVMProtectionPolicy, bool)
1. AzureVMWorkloadProtectionPolicy.AsAzureSQLProtectionPolicy() (*AzureSQLProtectionPolicy, bool)
1. AzureVMWorkloadProtectionPolicy.AsAzureVMWorkloadProtectionPolicy() (*AzureVMWorkloadProtectionPolicy, bool)
1. AzureVMWorkloadProtectionPolicy.AsBasicProtectionPolicy() (BasicProtectionPolicy, bool)
1. AzureVMWorkloadProtectionPolicy.AsGenericProtectionPolicy() (*GenericProtectionPolicy, bool)
1. AzureVMWorkloadProtectionPolicy.AsMabProtectionPolicy() (*MabProtectionPolicy, bool)
1. AzureVMWorkloadProtectionPolicy.AsProtectionPolicy() (*ProtectionPolicy, bool)
1. AzureVMWorkloadProtectionPolicy.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. AzureVMWorkloadSAPAseDatabaseProtectedItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.AsAzureVMWorkloadItem() (*AzureVMWorkloadItem, bool)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.AsAzureVMWorkloadSAPAseDatabaseWorkloadItem() (*AzureVMWorkloadSAPAseDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.AsAzureVMWorkloadSAPAseSystemWorkloadItem() (*AzureVMWorkloadSAPAseSystemWorkloadItem, bool)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.AsAzureVMWorkloadSAPHanaDatabaseWorkloadItem() (*AzureVMWorkloadSAPHanaDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.AsAzureVMWorkloadSAPHanaSystemWorkloadItem() (*AzureVMWorkloadSAPHanaSystemWorkloadItem, bool)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.AsAzureVMWorkloadSQLDatabaseWorkloadItem() (*AzureVMWorkloadSQLDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.AsAzureVMWorkloadSQLInstanceWorkloadItem() (*AzureVMWorkloadSQLInstanceWorkloadItem, bool)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.AsBasicAzureVMWorkloadItem() (BasicAzureVMWorkloadItem, bool)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.AsBasicWorkloadItem() (BasicWorkloadItem, bool)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.AsWorkloadItem() (*WorkloadItem, bool)
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPAseSystemProtectableItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.AsAzureVMWorkloadItem() (*AzureVMWorkloadItem, bool)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.AsAzureVMWorkloadSAPAseDatabaseWorkloadItem() (*AzureVMWorkloadSAPAseDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.AsAzureVMWorkloadSAPAseSystemWorkloadItem() (*AzureVMWorkloadSAPAseSystemWorkloadItem, bool)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.AsAzureVMWorkloadSAPHanaDatabaseWorkloadItem() (*AzureVMWorkloadSAPHanaDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.AsAzureVMWorkloadSAPHanaSystemWorkloadItem() (*AzureVMWorkloadSAPHanaSystemWorkloadItem, bool)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.AsAzureVMWorkloadSQLDatabaseWorkloadItem() (*AzureVMWorkloadSQLDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.AsAzureVMWorkloadSQLInstanceWorkloadItem() (*AzureVMWorkloadSQLInstanceWorkloadItem, bool)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.AsBasicAzureVMWorkloadItem() (BasicAzureVMWorkloadItem, bool)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.AsBasicWorkloadItem() (BasicWorkloadItem, bool)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.AsWorkloadItem() (*WorkloadItem, bool)
1. AzureVMWorkloadSAPAseSystemWorkloadItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.AsAzureVMWorkloadItem() (*AzureVMWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.AsAzureVMWorkloadSAPAseDatabaseWorkloadItem() (*AzureVMWorkloadSAPAseDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.AsAzureVMWorkloadSAPAseSystemWorkloadItem() (*AzureVMWorkloadSAPAseSystemWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.AsAzureVMWorkloadSAPHanaDatabaseWorkloadItem() (*AzureVMWorkloadSAPHanaDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.AsAzureVMWorkloadSAPHanaSystemWorkloadItem() (*AzureVMWorkloadSAPHanaSystemWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.AsAzureVMWorkloadSQLDatabaseWorkloadItem() (*AzureVMWorkloadSQLDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.AsAzureVMWorkloadSQLInstanceWorkloadItem() (*AzureVMWorkloadSQLInstanceWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.AsBasicAzureVMWorkloadItem() (BasicAzureVMWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.AsBasicWorkloadItem() (BasicWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.AsWorkloadItem() (*WorkloadItem, bool)
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. AzureVMWorkloadSAPHanaSystemProtectableItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.AsAzureVMWorkloadItem() (*AzureVMWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.AsAzureVMWorkloadSAPAseDatabaseWorkloadItem() (*AzureVMWorkloadSAPAseDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.AsAzureVMWorkloadSAPAseSystemWorkloadItem() (*AzureVMWorkloadSAPAseSystemWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.AsAzureVMWorkloadSAPHanaDatabaseWorkloadItem() (*AzureVMWorkloadSAPHanaDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.AsAzureVMWorkloadSAPHanaSystemWorkloadItem() (*AzureVMWorkloadSAPHanaSystemWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.AsAzureVMWorkloadSQLDatabaseWorkloadItem() (*AzureVMWorkloadSQLDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.AsAzureVMWorkloadSQLInstanceWorkloadItem() (*AzureVMWorkloadSQLInstanceWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.AsBasicAzureVMWorkloadItem() (BasicAzureVMWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.AsBasicWorkloadItem() (BasicWorkloadItem, bool)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.AsWorkloadItem() (*WorkloadItem, bool)
1. AzureVMWorkloadSAPHanaSystemWorkloadItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectableItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. AzureVMWorkloadSQLDatabaseProtectedItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.AsAzureVMWorkloadItem() (*AzureVMWorkloadItem, bool)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.AsAzureVMWorkloadSAPAseDatabaseWorkloadItem() (*AzureVMWorkloadSAPAseDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.AsAzureVMWorkloadSAPAseSystemWorkloadItem() (*AzureVMWorkloadSAPAseSystemWorkloadItem, bool)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.AsAzureVMWorkloadSAPHanaDatabaseWorkloadItem() (*AzureVMWorkloadSAPHanaDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.AsAzureVMWorkloadSAPHanaSystemWorkloadItem() (*AzureVMWorkloadSAPHanaSystemWorkloadItem, bool)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.AsAzureVMWorkloadSQLDatabaseWorkloadItem() (*AzureVMWorkloadSQLDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.AsAzureVMWorkloadSQLInstanceWorkloadItem() (*AzureVMWorkloadSQLInstanceWorkloadItem, bool)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.AsBasicAzureVMWorkloadItem() (BasicAzureVMWorkloadItem, bool)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.AsBasicWorkloadItem() (BasicWorkloadItem, bool)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.AsWorkloadItem() (*WorkloadItem, bool)
1. AzureVMWorkloadSQLDatabaseWorkloadItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. AzureVMWorkloadSQLInstanceProtectableItem.MarshalJSON() ([]byte, error)
1. AzureVMWorkloadSQLInstanceWorkloadItem.AsAzureVMWorkloadItem() (*AzureVMWorkloadItem, bool)
1. AzureVMWorkloadSQLInstanceWorkloadItem.AsAzureVMWorkloadSAPAseDatabaseWorkloadItem() (*AzureVMWorkloadSAPAseDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSQLInstanceWorkloadItem.AsAzureVMWorkloadSAPAseSystemWorkloadItem() (*AzureVMWorkloadSAPAseSystemWorkloadItem, bool)
1. AzureVMWorkloadSQLInstanceWorkloadItem.AsAzureVMWorkloadSAPHanaDatabaseWorkloadItem() (*AzureVMWorkloadSAPHanaDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSQLInstanceWorkloadItem.AsAzureVMWorkloadSAPHanaSystemWorkloadItem() (*AzureVMWorkloadSAPHanaSystemWorkloadItem, bool)
1. AzureVMWorkloadSQLInstanceWorkloadItem.AsAzureVMWorkloadSQLDatabaseWorkloadItem() (*AzureVMWorkloadSQLDatabaseWorkloadItem, bool)
1. AzureVMWorkloadSQLInstanceWorkloadItem.AsAzureVMWorkloadSQLInstanceWorkloadItem() (*AzureVMWorkloadSQLInstanceWorkloadItem, bool)
1. AzureVMWorkloadSQLInstanceWorkloadItem.AsBasicAzureVMWorkloadItem() (BasicAzureVMWorkloadItem, bool)
1. AzureVMWorkloadSQLInstanceWorkloadItem.AsBasicWorkloadItem() (BasicWorkloadItem, bool)
1. AzureVMWorkloadSQLInstanceWorkloadItem.AsWorkloadItem() (*WorkloadItem, bool)
1. AzureVMWorkloadSQLInstanceWorkloadItem.MarshalJSON() ([]byte, error)
1. AzureWorkloadAutoProtectionIntent.AsAzureRecoveryServiceVaultProtectionIntent() (*AzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsAzureResourceProtectionIntent() (*AzureResourceProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsAzureWorkloadAutoProtectionIntent() (*AzureWorkloadAutoProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsAzureWorkloadSQLAutoProtectionIntent() (*AzureWorkloadSQLAutoProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsBasicAzureRecoveryServiceVaultProtectionIntent() (BasicAzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsBasicAzureWorkloadAutoProtectionIntent() (BasicAzureWorkloadAutoProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsBasicProtectionIntent() (BasicProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsProtectionIntent() (*ProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.MarshalJSON() ([]byte, error)
1. AzureWorkloadBackupRequest.AsAzureFileShareBackupRequest() (*AzureFileShareBackupRequest, bool)
1. AzureWorkloadBackupRequest.AsAzureWorkloadBackupRequest() (*AzureWorkloadBackupRequest, bool)
1. AzureWorkloadBackupRequest.AsBasicRequest() (BasicRequest, bool)
1. AzureWorkloadBackupRequest.AsIaasVMBackupRequest() (*IaasVMBackupRequest, bool)
1. AzureWorkloadBackupRequest.AsRequest() (*Request, bool)
1. AzureWorkloadBackupRequest.MarshalJSON() ([]byte, error)
1. AzureWorkloadContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. AzureWorkloadContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. AzureWorkloadContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. AzureWorkloadContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. AzureWorkloadContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. AzureWorkloadContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. AzureWorkloadContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. AzureWorkloadContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. AzureWorkloadContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. AzureWorkloadContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. AzureWorkloadContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. AzureWorkloadContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. AzureWorkloadContainer.AsDpmContainer() (*DpmContainer, bool)
1. AzureWorkloadContainer.AsGenericContainer() (*GenericContainer, bool)
1. AzureWorkloadContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. AzureWorkloadContainer.AsMabContainer() (*MabContainer, bool)
1. AzureWorkloadContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. AzureWorkloadContainer.MarshalJSON() ([]byte, error)
1. AzureWorkloadJob.AsAzureIaaSVMJob() (*AzureIaaSVMJob, bool)
1. AzureWorkloadJob.AsAzureStorageJob() (*AzureStorageJob, bool)
1. AzureWorkloadJob.AsAzureWorkloadJob() (*AzureWorkloadJob, bool)
1. AzureWorkloadJob.AsBasicJob() (BasicJob, bool)
1. AzureWorkloadJob.AsDpmJob() (*DpmJob, bool)
1. AzureWorkloadJob.AsJob() (*Job, bool)
1. AzureWorkloadJob.AsMabJob() (*MabJob, bool)
1. AzureWorkloadJob.MarshalJSON() ([]byte, error)
1. AzureWorkloadJobExtendedInfo.MarshalJSON() ([]byte, error)
1. AzureWorkloadPointInTimeRecoveryPoint.AsAzureFileShareRecoveryPoint() (*AzureFileShareRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsAzureWorkloadPointInTimeRecoveryPoint() (*AzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsAzureWorkloadRecoveryPoint() (*AzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsAzureWorkloadSAPHanaPointInTimeRecoveryPoint() (*AzureWorkloadSAPHanaPointInTimeRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsAzureWorkloadSAPHanaRecoveryPoint() (*AzureWorkloadSAPHanaRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsAzureWorkloadSQLPointInTimeRecoveryPoint() (*AzureWorkloadSQLPointInTimeRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsAzureWorkloadSQLRecoveryPoint() (*AzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsBasicAzureWorkloadPointInTimeRecoveryPoint() (BasicAzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsBasicAzureWorkloadRecoveryPoint() (BasicAzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsBasicAzureWorkloadSQLRecoveryPoint() (BasicAzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsBasicRecoveryPoint() (BasicRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsGenericRecoveryPoint() (*GenericRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsIaasVMRecoveryPoint() (*IaasVMRecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.AsRecoveryPoint() (*RecoveryPoint, bool)
1. AzureWorkloadPointInTimeRecoveryPoint.MarshalJSON() ([]byte, error)
1. AzureWorkloadPointInTimeRestoreRequest.AsAzureFileShareRestoreRequest() (*AzureFileShareRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsAzureWorkloadPointInTimeRestoreRequest() (*AzureWorkloadPointInTimeRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsAzureWorkloadRestoreRequest() (*AzureWorkloadRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsAzureWorkloadSAPHanaPointInTimeRestoreRequest() (*AzureWorkloadSAPHanaPointInTimeRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsAzureWorkloadSAPHanaRestoreRequest() (*AzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsAzureWorkloadSQLPointInTimeRestoreRequest() (*AzureWorkloadSQLPointInTimeRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsAzureWorkloadSQLRestoreRequest() (*AzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsBasicAzureWorkloadRestoreRequest() (BasicAzureWorkloadRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsBasicAzureWorkloadSAPHanaRestoreRequest() (BasicAzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsBasicAzureWorkloadSQLRestoreRequest() (BasicAzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsBasicRestoreRequest() (BasicRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsIaasVMRestoreRequest() (*IaasVMRestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.AsRestoreRequest() (*RestoreRequest, bool)
1. AzureWorkloadPointInTimeRestoreRequest.MarshalJSON() ([]byte, error)
1. AzureWorkloadRecoveryPoint.AsAzureFileShareRecoveryPoint() (*AzureFileShareRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsAzureWorkloadPointInTimeRecoveryPoint() (*AzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsAzureWorkloadRecoveryPoint() (*AzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsAzureWorkloadSAPHanaPointInTimeRecoveryPoint() (*AzureWorkloadSAPHanaPointInTimeRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsAzureWorkloadSAPHanaRecoveryPoint() (*AzureWorkloadSAPHanaRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsAzureWorkloadSQLPointInTimeRecoveryPoint() (*AzureWorkloadSQLPointInTimeRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsAzureWorkloadSQLRecoveryPoint() (*AzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsBasicAzureWorkloadPointInTimeRecoveryPoint() (BasicAzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsBasicAzureWorkloadRecoveryPoint() (BasicAzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsBasicAzureWorkloadSQLRecoveryPoint() (BasicAzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsBasicRecoveryPoint() (BasicRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsGenericRecoveryPoint() (*GenericRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsIaasVMRecoveryPoint() (*IaasVMRecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.AsRecoveryPoint() (*RecoveryPoint, bool)
1. AzureWorkloadRecoveryPoint.MarshalJSON() ([]byte, error)
1. AzureWorkloadRestoreRequest.AsAzureFileShareRestoreRequest() (*AzureFileShareRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsAzureWorkloadPointInTimeRestoreRequest() (*AzureWorkloadPointInTimeRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsAzureWorkloadRestoreRequest() (*AzureWorkloadRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsAzureWorkloadSAPHanaPointInTimeRestoreRequest() (*AzureWorkloadSAPHanaPointInTimeRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsAzureWorkloadSAPHanaRestoreRequest() (*AzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsAzureWorkloadSQLPointInTimeRestoreRequest() (*AzureWorkloadSQLPointInTimeRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsAzureWorkloadSQLRestoreRequest() (*AzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsBasicAzureWorkloadRestoreRequest() (BasicAzureWorkloadRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsBasicAzureWorkloadSAPHanaRestoreRequest() (BasicAzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsBasicAzureWorkloadSQLRestoreRequest() (BasicAzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsBasicRestoreRequest() (BasicRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsIaasVMRestoreRequest() (*IaasVMRestoreRequest, bool)
1. AzureWorkloadRestoreRequest.AsRestoreRequest() (*RestoreRequest, bool)
1. AzureWorkloadRestoreRequest.MarshalJSON() ([]byte, error)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsAzureFileShareRecoveryPoint() (*AzureFileShareRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsAzureWorkloadPointInTimeRecoveryPoint() (*AzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsAzureWorkloadRecoveryPoint() (*AzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsAzureWorkloadSAPHanaPointInTimeRecoveryPoint() (*AzureWorkloadSAPHanaPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsAzureWorkloadSAPHanaRecoveryPoint() (*AzureWorkloadSAPHanaRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsAzureWorkloadSQLPointInTimeRecoveryPoint() (*AzureWorkloadSQLPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsAzureWorkloadSQLRecoveryPoint() (*AzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsBasicAzureWorkloadPointInTimeRecoveryPoint() (BasicAzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsBasicAzureWorkloadRecoveryPoint() (BasicAzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsBasicAzureWorkloadSQLRecoveryPoint() (BasicAzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsBasicRecoveryPoint() (BasicRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsGenericRecoveryPoint() (*GenericRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsIaasVMRecoveryPoint() (*IaasVMRecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.AsRecoveryPoint() (*RecoveryPoint, bool)
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint.MarshalJSON() ([]byte, error)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsAzureFileShareRestoreRequest() (*AzureFileShareRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsAzureWorkloadPointInTimeRestoreRequest() (*AzureWorkloadPointInTimeRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsAzureWorkloadRestoreRequest() (*AzureWorkloadRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsAzureWorkloadSAPHanaPointInTimeRestoreRequest() (*AzureWorkloadSAPHanaPointInTimeRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsAzureWorkloadSAPHanaRestoreRequest() (*AzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsAzureWorkloadSQLPointInTimeRestoreRequest() (*AzureWorkloadSQLPointInTimeRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsAzureWorkloadSQLRestoreRequest() (*AzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsBasicAzureWorkloadRestoreRequest() (BasicAzureWorkloadRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsBasicAzureWorkloadSAPHanaRestoreRequest() (BasicAzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsBasicAzureWorkloadSQLRestoreRequest() (BasicAzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsBasicRestoreRequest() (BasicRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsIaasVMRestoreRequest() (*IaasVMRestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.AsRestoreRequest() (*RestoreRequest, bool)
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.MarshalJSON() ([]byte, error)
1. AzureWorkloadSAPHanaRecoveryPoint.AsAzureFileShareRecoveryPoint() (*AzureFileShareRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsAzureWorkloadPointInTimeRecoveryPoint() (*AzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsAzureWorkloadRecoveryPoint() (*AzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsAzureWorkloadSAPHanaPointInTimeRecoveryPoint() (*AzureWorkloadSAPHanaPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsAzureWorkloadSAPHanaRecoveryPoint() (*AzureWorkloadSAPHanaRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsAzureWorkloadSQLPointInTimeRecoveryPoint() (*AzureWorkloadSQLPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsAzureWorkloadSQLRecoveryPoint() (*AzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsBasicAzureWorkloadPointInTimeRecoveryPoint() (BasicAzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsBasicAzureWorkloadRecoveryPoint() (BasicAzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsBasicAzureWorkloadSQLRecoveryPoint() (BasicAzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsBasicRecoveryPoint() (BasicRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsGenericRecoveryPoint() (*GenericRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsIaasVMRecoveryPoint() (*IaasVMRecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.AsRecoveryPoint() (*RecoveryPoint, bool)
1. AzureWorkloadSAPHanaRecoveryPoint.MarshalJSON() ([]byte, error)
1. AzureWorkloadSAPHanaRestoreRequest.AsAzureFileShareRestoreRequest() (*AzureFileShareRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsAzureWorkloadPointInTimeRestoreRequest() (*AzureWorkloadPointInTimeRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsAzureWorkloadRestoreRequest() (*AzureWorkloadRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsAzureWorkloadSAPHanaPointInTimeRestoreRequest() (*AzureWorkloadSAPHanaPointInTimeRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsAzureWorkloadSAPHanaRestoreRequest() (*AzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsAzureWorkloadSQLPointInTimeRestoreRequest() (*AzureWorkloadSQLPointInTimeRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsAzureWorkloadSQLRestoreRequest() (*AzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsBasicAzureWorkloadRestoreRequest() (BasicAzureWorkloadRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsBasicAzureWorkloadSAPHanaRestoreRequest() (BasicAzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsBasicAzureWorkloadSQLRestoreRequest() (BasicAzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsBasicRestoreRequest() (BasicRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsIaasVMRestoreRequest() (*IaasVMRestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.AsRestoreRequest() (*RestoreRequest, bool)
1. AzureWorkloadSAPHanaRestoreRequest.MarshalJSON() ([]byte, error)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureRecoveryServiceVaultProtectionIntent() (*AzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureResourceProtectionIntent() (*AzureResourceProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureWorkloadAutoProtectionIntent() (*AzureWorkloadAutoProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureWorkloadSQLAutoProtectionIntent() (*AzureWorkloadSQLAutoProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsBasicAzureRecoveryServiceVaultProtectionIntent() (BasicAzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsBasicAzureWorkloadAutoProtectionIntent() (BasicAzureWorkloadAutoProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsBasicProtectionIntent() (BasicProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsProtectionIntent() (*ProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.MarshalJSON() ([]byte, error)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsAzureFileShareRecoveryPoint() (*AzureFileShareRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsAzureWorkloadPointInTimeRecoveryPoint() (*AzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsAzureWorkloadRecoveryPoint() (*AzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsAzureWorkloadSAPHanaPointInTimeRecoveryPoint() (*AzureWorkloadSAPHanaPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsAzureWorkloadSAPHanaRecoveryPoint() (*AzureWorkloadSAPHanaRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsAzureWorkloadSQLPointInTimeRecoveryPoint() (*AzureWorkloadSQLPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsAzureWorkloadSQLRecoveryPoint() (*AzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsBasicAzureWorkloadPointInTimeRecoveryPoint() (BasicAzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsBasicAzureWorkloadRecoveryPoint() (BasicAzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsBasicAzureWorkloadSQLRecoveryPoint() (BasicAzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsBasicRecoveryPoint() (BasicRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsGenericRecoveryPoint() (*GenericRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsIaasVMRecoveryPoint() (*IaasVMRecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.AsRecoveryPoint() (*RecoveryPoint, bool)
1. AzureWorkloadSQLPointInTimeRecoveryPoint.MarshalJSON() ([]byte, error)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsAzureFileShareRestoreRequest() (*AzureFileShareRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsAzureWorkloadPointInTimeRestoreRequest() (*AzureWorkloadPointInTimeRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsAzureWorkloadRestoreRequest() (*AzureWorkloadRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsAzureWorkloadSAPHanaPointInTimeRestoreRequest() (*AzureWorkloadSAPHanaPointInTimeRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsAzureWorkloadSAPHanaRestoreRequest() (*AzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsAzureWorkloadSQLPointInTimeRestoreRequest() (*AzureWorkloadSQLPointInTimeRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsAzureWorkloadSQLRestoreRequest() (*AzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsBasicAzureWorkloadRestoreRequest() (BasicAzureWorkloadRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsBasicAzureWorkloadSAPHanaRestoreRequest() (BasicAzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsBasicAzureWorkloadSQLRestoreRequest() (BasicAzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsBasicRestoreRequest() (BasicRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsIaasVMRestoreRequest() (*IaasVMRestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.AsRestoreRequest() (*RestoreRequest, bool)
1. AzureWorkloadSQLPointInTimeRestoreRequest.MarshalJSON() ([]byte, error)
1. AzureWorkloadSQLRecoveryPoint.AsAzureFileShareRecoveryPoint() (*AzureFileShareRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsAzureWorkloadPointInTimeRecoveryPoint() (*AzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsAzureWorkloadRecoveryPoint() (*AzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsAzureWorkloadSAPHanaPointInTimeRecoveryPoint() (*AzureWorkloadSAPHanaPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsAzureWorkloadSAPHanaRecoveryPoint() (*AzureWorkloadSAPHanaRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsAzureWorkloadSQLPointInTimeRecoveryPoint() (*AzureWorkloadSQLPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsAzureWorkloadSQLRecoveryPoint() (*AzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsBasicAzureWorkloadPointInTimeRecoveryPoint() (BasicAzureWorkloadPointInTimeRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsBasicAzureWorkloadRecoveryPoint() (BasicAzureWorkloadRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsBasicAzureWorkloadSQLRecoveryPoint() (BasicAzureWorkloadSQLRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsBasicRecoveryPoint() (BasicRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsGenericRecoveryPoint() (*GenericRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsIaasVMRecoveryPoint() (*IaasVMRecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.AsRecoveryPoint() (*RecoveryPoint, bool)
1. AzureWorkloadSQLRecoveryPoint.MarshalJSON() ([]byte, error)
1. AzureWorkloadSQLRecoveryPointExtendedInfo.MarshalJSON() ([]byte, error)
1. AzureWorkloadSQLRestoreRequest.AsAzureFileShareRestoreRequest() (*AzureFileShareRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsAzureWorkloadPointInTimeRestoreRequest() (*AzureWorkloadPointInTimeRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsAzureWorkloadRestoreRequest() (*AzureWorkloadRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsAzureWorkloadSAPHanaPointInTimeRestoreRequest() (*AzureWorkloadSAPHanaPointInTimeRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsAzureWorkloadSAPHanaRestoreRequest() (*AzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsAzureWorkloadSQLPointInTimeRestoreRequest() (*AzureWorkloadSQLPointInTimeRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsAzureWorkloadSQLRestoreRequest() (*AzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsBasicAzureWorkloadRestoreRequest() (BasicAzureWorkloadRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsBasicAzureWorkloadSAPHanaRestoreRequest() (BasicAzureWorkloadSAPHanaRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsBasicAzureWorkloadSQLRestoreRequest() (BasicAzureWorkloadSQLRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsBasicRestoreRequest() (BasicRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsIaasVMRestoreRequest() (*IaasVMRestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.AsRestoreRequest() (*RestoreRequest, bool)
1. AzureWorkloadSQLRestoreRequest.MarshalJSON() ([]byte, error)
1. BackupsClient.Trigger(context.Context, string, string, string, string, string, RequestResource) (autorest.Response, error)
1. BackupsClient.TriggerPreparer(context.Context, string, string, string, string, string, RequestResource) (*http.Request, error)
1. BackupsClient.TriggerResponder(*http.Response) (autorest.Response, error)
1. BackupsClient.TriggerSender(*http.Request) (*http.Response, error)
1. ClientDiscoveryResponse.IsEmpty() bool
1. ClientDiscoveryResponseIterator.NotDone() bool
1. ClientDiscoveryResponseIterator.Response() ClientDiscoveryResponse
1. ClientDiscoveryResponseIterator.Value() ClientDiscoveryValueForSingleAPI
1. ClientDiscoveryResponsePage.NotDone() bool
1. ClientDiscoveryResponsePage.Response() ClientDiscoveryResponse
1. ClientDiscoveryResponsePage.Values() []ClientDiscoveryValueForSingleAPI
1. DPMProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. DPMProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. DPMProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. DPMProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. DPMProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. DPMProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. DPMProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. DPMProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. DPMProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. DPMProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. DPMProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. DPMProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. DPMProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. DPMProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. DPMProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. DPMProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. DPMProtectedItem.MarshalJSON() ([]byte, error)
1. DPMProtectedItemExtendedInfo.MarshalJSON() ([]byte, error)
1. DpmBackupEngine.AsAzureBackupServerEngine() (*AzureBackupServerEngine, bool)
1. DpmBackupEngine.AsBasicEngineBase() (BasicEngineBase, bool)
1. DpmBackupEngine.AsDpmBackupEngine() (*DpmBackupEngine, bool)
1. DpmBackupEngine.AsEngineBase() (*EngineBase, bool)
1. DpmBackupEngine.MarshalJSON() ([]byte, error)
1. DpmContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. DpmContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. DpmContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. DpmContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. DpmContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. DpmContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. DpmContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. DpmContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. DpmContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. DpmContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. DpmContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. DpmContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. DpmContainer.AsDpmContainer() (*DpmContainer, bool)
1. DpmContainer.AsGenericContainer() (*GenericContainer, bool)
1. DpmContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. DpmContainer.AsMabContainer() (*MabContainer, bool)
1. DpmContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. DpmContainer.MarshalJSON() ([]byte, error)
1. DpmJob.AsAzureIaaSVMJob() (*AzureIaaSVMJob, bool)
1. DpmJob.AsAzureStorageJob() (*AzureStorageJob, bool)
1. DpmJob.AsAzureWorkloadJob() (*AzureWorkloadJob, bool)
1. DpmJob.AsBasicJob() (BasicJob, bool)
1. DpmJob.AsDpmJob() (*DpmJob, bool)
1. DpmJob.AsJob() (*Job, bool)
1. DpmJob.AsMabJob() (*MabJob, bool)
1. DpmJob.MarshalJSON() ([]byte, error)
1. DpmJobExtendedInfo.MarshalJSON() ([]byte, error)
1. EngineBase.AsAzureBackupServerEngine() (*AzureBackupServerEngine, bool)
1. EngineBase.AsBasicEngineBase() (BasicEngineBase, bool)
1. EngineBase.AsDpmBackupEngine() (*DpmBackupEngine, bool)
1. EngineBase.AsEngineBase() (*EngineBase, bool)
1. EngineBase.MarshalJSON() ([]byte, error)
1. EngineBaseResource.MarshalJSON() ([]byte, error)
1. EngineBaseResourceList.IsEmpty() bool
1. EngineBaseResourceListIterator.NotDone() bool
1. EngineBaseResourceListIterator.Response() EngineBaseResourceList
1. EngineBaseResourceListIterator.Value() EngineBaseResource
1. EngineBaseResourceListPage.NotDone() bool
1. EngineBaseResourceListPage.Response() EngineBaseResourceList
1. EngineBaseResourceListPage.Values() []EngineBaseResource
1. EnginesClient.Get(context.Context, string, string, string, string, string) (EngineBaseResource, error)
1. EnginesClient.GetPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. EnginesClient.GetResponder(*http.Response) (EngineBaseResource, error)
1. EnginesClient.GetSender(*http.Request) (*http.Response, error)
1. EnginesClient.List(context.Context, string, string, string, string) (EngineBaseResourceListPage, error)
1. EnginesClient.ListComplete(context.Context, string, string, string, string) (EngineBaseResourceListIterator, error)
1. EnginesClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. EnginesClient.ListResponder(*http.Response) (EngineBaseResourceList, error)
1. EnginesClient.ListSender(*http.Request) (*http.Response, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)
1. ExportJobsOperationResultInfo.AsBasicOperationResultInfoBase() (BasicOperationResultInfoBase, bool)
1. ExportJobsOperationResultInfo.AsExportJobsOperationResultInfo() (*ExportJobsOperationResultInfo, bool)
1. ExportJobsOperationResultInfo.AsOperationResultInfo() (*OperationResultInfo, bool)
1. ExportJobsOperationResultInfo.AsOperationResultInfoBase() (*OperationResultInfoBase, bool)
1. ExportJobsOperationResultInfo.MarshalJSON() ([]byte, error)
1. ExportJobsOperationResultsClient.Get(context.Context, string, string, string) (OperationResultInfoBaseResource, error)
1. ExportJobsOperationResultsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ExportJobsOperationResultsClient.GetResponder(*http.Response) (OperationResultInfoBaseResource, error)
1. ExportJobsOperationResultsClient.GetSender(*http.Request) (*http.Response, error)
1. FeatureSupportClient.Validate(context.Context, string, BasicFeatureSupportRequest) (AzureVMResourceFeatureSupportResponse, error)
1. FeatureSupportClient.ValidatePreparer(context.Context, string, BasicFeatureSupportRequest) (*http.Request, error)
1. FeatureSupportClient.ValidateResponder(*http.Response) (AzureVMResourceFeatureSupportResponse, error)
1. FeatureSupportClient.ValidateSender(*http.Request) (*http.Response, error)
1. FeatureSupportRequest.AsAzureBackupGoalFeatureSupportRequest() (*AzureBackupGoalFeatureSupportRequest, bool)
1. FeatureSupportRequest.AsAzureVMResourceFeatureSupportRequest() (*AzureVMResourceFeatureSupportRequest, bool)
1. FeatureSupportRequest.AsBasicFeatureSupportRequest() (BasicFeatureSupportRequest, bool)
1. FeatureSupportRequest.AsFeatureSupportRequest() (*FeatureSupportRequest, bool)
1. FeatureSupportRequest.MarshalJSON() ([]byte, error)
1. GenericContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. GenericContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. GenericContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. GenericContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. GenericContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. GenericContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. GenericContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. GenericContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. GenericContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. GenericContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. GenericContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. GenericContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. GenericContainer.AsDpmContainer() (*DpmContainer, bool)
1. GenericContainer.AsGenericContainer() (*GenericContainer, bool)
1. GenericContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. GenericContainer.AsMabContainer() (*MabContainer, bool)
1. GenericContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. GenericContainer.MarshalJSON() ([]byte, error)
1. GenericContainerExtendedInfo.MarshalJSON() ([]byte, error)
1. GenericProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. GenericProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. GenericProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. GenericProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. GenericProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. GenericProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. GenericProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. GenericProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. GenericProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. GenericProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. GenericProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. GenericProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. GenericProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. GenericProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. GenericProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. GenericProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. GenericProtectedItem.MarshalJSON() ([]byte, error)
1. GenericProtectionPolicy.AsAzureFileShareProtectionPolicy() (*AzureFileShareProtectionPolicy, bool)
1. GenericProtectionPolicy.AsAzureIaaSVMProtectionPolicy() (*AzureIaaSVMProtectionPolicy, bool)
1. GenericProtectionPolicy.AsAzureSQLProtectionPolicy() (*AzureSQLProtectionPolicy, bool)
1. GenericProtectionPolicy.AsAzureVMWorkloadProtectionPolicy() (*AzureVMWorkloadProtectionPolicy, bool)
1. GenericProtectionPolicy.AsBasicProtectionPolicy() (BasicProtectionPolicy, bool)
1. GenericProtectionPolicy.AsGenericProtectionPolicy() (*GenericProtectionPolicy, bool)
1. GenericProtectionPolicy.AsMabProtectionPolicy() (*MabProtectionPolicy, bool)
1. GenericProtectionPolicy.AsProtectionPolicy() (*ProtectionPolicy, bool)
1. GenericProtectionPolicy.MarshalJSON() ([]byte, error)
1. GenericRecoveryPoint.AsAzureFileShareRecoveryPoint() (*AzureFileShareRecoveryPoint, bool)
1. GenericRecoveryPoint.AsAzureWorkloadPointInTimeRecoveryPoint() (*AzureWorkloadPointInTimeRecoveryPoint, bool)
1. GenericRecoveryPoint.AsAzureWorkloadRecoveryPoint() (*AzureWorkloadRecoveryPoint, bool)
1. GenericRecoveryPoint.AsAzureWorkloadSAPHanaPointInTimeRecoveryPoint() (*AzureWorkloadSAPHanaPointInTimeRecoveryPoint, bool)
1. GenericRecoveryPoint.AsAzureWorkloadSAPHanaRecoveryPoint() (*AzureWorkloadSAPHanaRecoveryPoint, bool)
1. GenericRecoveryPoint.AsAzureWorkloadSQLPointInTimeRecoveryPoint() (*AzureWorkloadSQLPointInTimeRecoveryPoint, bool)
1. GenericRecoveryPoint.AsAzureWorkloadSQLRecoveryPoint() (*AzureWorkloadSQLRecoveryPoint, bool)
1. GenericRecoveryPoint.AsBasicAzureWorkloadPointInTimeRecoveryPoint() (BasicAzureWorkloadPointInTimeRecoveryPoint, bool)
1. GenericRecoveryPoint.AsBasicAzureWorkloadRecoveryPoint() (BasicAzureWorkloadRecoveryPoint, bool)
1. GenericRecoveryPoint.AsBasicAzureWorkloadSQLRecoveryPoint() (BasicAzureWorkloadSQLRecoveryPoint, bool)
1. GenericRecoveryPoint.AsBasicRecoveryPoint() (BasicRecoveryPoint, bool)
1. GenericRecoveryPoint.AsGenericRecoveryPoint() (*GenericRecoveryPoint, bool)
1. GenericRecoveryPoint.AsIaasVMRecoveryPoint() (*IaasVMRecoveryPoint, bool)
1. GenericRecoveryPoint.AsRecoveryPoint() (*RecoveryPoint, bool)
1. GenericRecoveryPoint.MarshalJSON() ([]byte, error)
1. ILRRequest.AsAzureFileShareProvisionILRRequest() (*AzureFileShareProvisionILRRequest, bool)
1. ILRRequest.AsBasicILRRequest() (BasicILRRequest, bool)
1. ILRRequest.AsILRRequest() (*ILRRequest, bool)
1. ILRRequest.AsIaasVMILRRegistrationRequest() (*IaasVMILRRegistrationRequest, bool)
1. ILRRequest.MarshalJSON() ([]byte, error)
1. ILRRequestResource.MarshalJSON() ([]byte, error)
1. IaaSVMContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. IaaSVMContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. IaaSVMContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. IaaSVMContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. IaaSVMContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. IaaSVMContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. IaaSVMContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. IaaSVMContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. IaaSVMContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. IaaSVMContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. IaaSVMContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. IaaSVMContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. IaaSVMContainer.AsDpmContainer() (*DpmContainer, bool)
1. IaaSVMContainer.AsGenericContainer() (*GenericContainer, bool)
1. IaaSVMContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. IaaSVMContainer.AsMabContainer() (*MabContainer, bool)
1. IaaSVMContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. IaaSVMContainer.MarshalJSON() ([]byte, error)
1. IaaSVMProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. IaaSVMProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. IaaSVMProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. IaaSVMProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. IaaSVMProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. IaaSVMProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. IaaSVMProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. IaaSVMProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. IaaSVMProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. IaaSVMProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. IaaSVMProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. IaaSVMProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. IaaSVMProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. IaaSVMProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. IaaSVMProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. IaaSVMProtectableItem.MarshalJSON() ([]byte, error)
1. IaasVMBackupRequest.AsAzureFileShareBackupRequest() (*AzureFileShareBackupRequest, bool)
1. IaasVMBackupRequest.AsAzureWorkloadBackupRequest() (*AzureWorkloadBackupRequest, bool)
1. IaasVMBackupRequest.AsBasicRequest() (BasicRequest, bool)
1. IaasVMBackupRequest.AsIaasVMBackupRequest() (*IaasVMBackupRequest, bool)
1. IaasVMBackupRequest.AsRequest() (*Request, bool)
1. IaasVMBackupRequest.MarshalJSON() ([]byte, error)
1. IaasVMILRRegistrationRequest.AsAzureFileShareProvisionILRRequest() (*AzureFileShareProvisionILRRequest, bool)
1. IaasVMILRRegistrationRequest.AsBasicILRRequest() (BasicILRRequest, bool)
1. IaasVMILRRegistrationRequest.AsILRRequest() (*ILRRequest, bool)
1. IaasVMILRRegistrationRequest.AsIaasVMILRRegistrationRequest() (*IaasVMILRRegistrationRequest, bool)
1. IaasVMILRRegistrationRequest.MarshalJSON() ([]byte, error)
1. IaasVMRecoveryPoint.AsAzureFileShareRecoveryPoint() (*AzureFileShareRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsAzureWorkloadPointInTimeRecoveryPoint() (*AzureWorkloadPointInTimeRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsAzureWorkloadRecoveryPoint() (*AzureWorkloadRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsAzureWorkloadSAPHanaPointInTimeRecoveryPoint() (*AzureWorkloadSAPHanaPointInTimeRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsAzureWorkloadSAPHanaRecoveryPoint() (*AzureWorkloadSAPHanaRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsAzureWorkloadSQLPointInTimeRecoveryPoint() (*AzureWorkloadSQLPointInTimeRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsAzureWorkloadSQLRecoveryPoint() (*AzureWorkloadSQLRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsBasicAzureWorkloadPointInTimeRecoveryPoint() (BasicAzureWorkloadPointInTimeRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsBasicAzureWorkloadRecoveryPoint() (BasicAzureWorkloadRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsBasicAzureWorkloadSQLRecoveryPoint() (BasicAzureWorkloadSQLRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsBasicRecoveryPoint() (BasicRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsGenericRecoveryPoint() (*GenericRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsIaasVMRecoveryPoint() (*IaasVMRecoveryPoint, bool)
1. IaasVMRecoveryPoint.AsRecoveryPoint() (*RecoveryPoint, bool)
1. IaasVMRecoveryPoint.MarshalJSON() ([]byte, error)
1. IaasVMRestoreRequest.AsAzureFileShareRestoreRequest() (*AzureFileShareRestoreRequest, bool)
1. IaasVMRestoreRequest.AsAzureWorkloadPointInTimeRestoreRequest() (*AzureWorkloadPointInTimeRestoreRequest, bool)
1. IaasVMRestoreRequest.AsAzureWorkloadRestoreRequest() (*AzureWorkloadRestoreRequest, bool)
1. IaasVMRestoreRequest.AsAzureWorkloadSAPHanaPointInTimeRestoreRequest() (*AzureWorkloadSAPHanaPointInTimeRestoreRequest, bool)
1. IaasVMRestoreRequest.AsAzureWorkloadSAPHanaRestoreRequest() (*AzureWorkloadSAPHanaRestoreRequest, bool)
1. IaasVMRestoreRequest.AsAzureWorkloadSQLPointInTimeRestoreRequest() (*AzureWorkloadSQLPointInTimeRestoreRequest, bool)
1. IaasVMRestoreRequest.AsAzureWorkloadSQLRestoreRequest() (*AzureWorkloadSQLRestoreRequest, bool)
1. IaasVMRestoreRequest.AsBasicAzureWorkloadRestoreRequest() (BasicAzureWorkloadRestoreRequest, bool)
1. IaasVMRestoreRequest.AsBasicAzureWorkloadSAPHanaRestoreRequest() (BasicAzureWorkloadSAPHanaRestoreRequest, bool)
1. IaasVMRestoreRequest.AsBasicAzureWorkloadSQLRestoreRequest() (BasicAzureWorkloadSQLRestoreRequest, bool)
1. IaasVMRestoreRequest.AsBasicRestoreRequest() (BasicRestoreRequest, bool)
1. IaasVMRestoreRequest.AsIaasVMRestoreRequest() (*IaasVMRestoreRequest, bool)
1. IaasVMRestoreRequest.AsRestoreRequest() (*RestoreRequest, bool)
1. IaasVMRestoreRequest.MarshalJSON() ([]byte, error)
1. InquiryValidation.MarshalJSON() ([]byte, error)
1. ItemLevelRecoveryConnectionsClient.Provision(context.Context, string, string, string, string, string, string, ILRRequestResource) (autorest.Response, error)
1. ItemLevelRecoveryConnectionsClient.ProvisionPreparer(context.Context, string, string, string, string, string, string, ILRRequestResource) (*http.Request, error)
1. ItemLevelRecoveryConnectionsClient.ProvisionResponder(*http.Response) (autorest.Response, error)
1. ItemLevelRecoveryConnectionsClient.ProvisionSender(*http.Request) (*http.Response, error)
1. ItemLevelRecoveryConnectionsClient.Revoke(context.Context, string, string, string, string, string, string) (autorest.Response, error)
1. ItemLevelRecoveryConnectionsClient.RevokePreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. ItemLevelRecoveryConnectionsClient.RevokeResponder(*http.Response) (autorest.Response, error)
1. ItemLevelRecoveryConnectionsClient.RevokeSender(*http.Request) (*http.Response, error)
1. Job.AsAzureIaaSVMJob() (*AzureIaaSVMJob, bool)
1. Job.AsAzureStorageJob() (*AzureStorageJob, bool)
1. Job.AsAzureWorkloadJob() (*AzureWorkloadJob, bool)
1. Job.AsBasicJob() (BasicJob, bool)
1. Job.AsDpmJob() (*DpmJob, bool)
1. Job.AsJob() (*Job, bool)
1. Job.AsMabJob() (*MabJob, bool)
1. Job.MarshalJSON() ([]byte, error)
1. JobCancellationsClient.Trigger(context.Context, string, string, string) (autorest.Response, error)
1. JobCancellationsClient.TriggerPreparer(context.Context, string, string, string) (*http.Request, error)
1. JobCancellationsClient.TriggerResponder(*http.Response) (autorest.Response, error)
1. JobCancellationsClient.TriggerSender(*http.Request) (*http.Response, error)
1. JobDetailsClient.Get(context.Context, string, string, string) (JobResource, error)
1. JobDetailsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. JobDetailsClient.GetResponder(*http.Response) (JobResource, error)
1. JobDetailsClient.GetSender(*http.Request) (*http.Response, error)
1. JobOperationResultsClient.Get(context.Context, string, string, string, string) (autorest.Response, error)
1. JobOperationResultsClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. JobOperationResultsClient.GetResponder(*http.Response) (autorest.Response, error)
1. JobOperationResultsClient.GetSender(*http.Request) (*http.Response, error)
1. JobResource.MarshalJSON() ([]byte, error)
1. JobResourceList.IsEmpty() bool
1. JobResourceListIterator.NotDone() bool
1. JobResourceListIterator.Response() JobResourceList
1. JobResourceListIterator.Value() JobResource
1. JobResourceListPage.NotDone() bool
1. JobResourceListPage.Response() JobResourceList
1. JobResourceListPage.Values() []JobResource
1. JobsClient.List(context.Context, string, string, string, string) (JobResourceListPage, error)
1. JobsClient.ListComplete(context.Context, string, string, string, string) (JobResourceListIterator, error)
1. JobsClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. JobsClient.ListResponder(*http.Response) (JobResourceList, error)
1. JobsClient.ListSender(*http.Request) (*http.Response, error)
1. JobsGroupClient.Export(context.Context, string, string, string) (autorest.Response, error)
1. JobsGroupClient.ExportPreparer(context.Context, string, string, string) (*http.Request, error)
1. JobsGroupClient.ExportResponder(*http.Response) (autorest.Response, error)
1. JobsGroupClient.ExportSender(*http.Request) (*http.Response, error)
1. LogSchedulePolicy.AsBasicSchedulePolicy() (BasicSchedulePolicy, bool)
1. LogSchedulePolicy.AsLogSchedulePolicy() (*LogSchedulePolicy, bool)
1. LogSchedulePolicy.AsLongTermSchedulePolicy() (*LongTermSchedulePolicy, bool)
1. LogSchedulePolicy.AsSchedulePolicy() (*SchedulePolicy, bool)
1. LogSchedulePolicy.AsSimpleSchedulePolicy() (*SimpleSchedulePolicy, bool)
1. LogSchedulePolicy.MarshalJSON() ([]byte, error)
1. LongTermRetentionPolicy.AsBasicRetentionPolicy() (BasicRetentionPolicy, bool)
1. LongTermRetentionPolicy.AsLongTermRetentionPolicy() (*LongTermRetentionPolicy, bool)
1. LongTermRetentionPolicy.AsRetentionPolicy() (*RetentionPolicy, bool)
1. LongTermRetentionPolicy.AsSimpleRetentionPolicy() (*SimpleRetentionPolicy, bool)
1. LongTermRetentionPolicy.MarshalJSON() ([]byte, error)
1. LongTermSchedulePolicy.AsBasicSchedulePolicy() (BasicSchedulePolicy, bool)
1. LongTermSchedulePolicy.AsLogSchedulePolicy() (*LogSchedulePolicy, bool)
1. LongTermSchedulePolicy.AsLongTermSchedulePolicy() (*LongTermSchedulePolicy, bool)
1. LongTermSchedulePolicy.AsSchedulePolicy() (*SchedulePolicy, bool)
1. LongTermSchedulePolicy.AsSimpleSchedulePolicy() (*SimpleSchedulePolicy, bool)
1. LongTermSchedulePolicy.MarshalJSON() ([]byte, error)
1. MabContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. MabContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. MabContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. MabContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. MabContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. MabContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. MabContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. MabContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. MabContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. MabContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. MabContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. MabContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. MabContainer.AsDpmContainer() (*DpmContainer, bool)
1. MabContainer.AsGenericContainer() (*GenericContainer, bool)
1. MabContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. MabContainer.AsMabContainer() (*MabContainer, bool)
1. MabContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. MabContainer.MarshalJSON() ([]byte, error)
1. MabErrorInfo.MarshalJSON() ([]byte, error)
1. MabFileFolderProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. MabFileFolderProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. MabFileFolderProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. MabFileFolderProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. MabFileFolderProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. MabFileFolderProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. MabFileFolderProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. MabFileFolderProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. MabFileFolderProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. MabFileFolderProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. MabFileFolderProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. MabFileFolderProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. MabFileFolderProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. MabFileFolderProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. MabFileFolderProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. MabFileFolderProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. MabFileFolderProtectedItem.MarshalJSON() ([]byte, error)
1. MabJob.AsAzureIaaSVMJob() (*AzureIaaSVMJob, bool)
1. MabJob.AsAzureStorageJob() (*AzureStorageJob, bool)
1. MabJob.AsAzureWorkloadJob() (*AzureWorkloadJob, bool)
1. MabJob.AsBasicJob() (BasicJob, bool)
1. MabJob.AsDpmJob() (*DpmJob, bool)
1. MabJob.AsJob() (*Job, bool)
1. MabJob.AsMabJob() (*MabJob, bool)
1. MabJob.MarshalJSON() ([]byte, error)
1. MabJobExtendedInfo.MarshalJSON() ([]byte, error)
1. MabProtectionPolicy.AsAzureFileShareProtectionPolicy() (*AzureFileShareProtectionPolicy, bool)
1. MabProtectionPolicy.AsAzureIaaSVMProtectionPolicy() (*AzureIaaSVMProtectionPolicy, bool)
1. MabProtectionPolicy.AsAzureSQLProtectionPolicy() (*AzureSQLProtectionPolicy, bool)
1. MabProtectionPolicy.AsAzureVMWorkloadProtectionPolicy() (*AzureVMWorkloadProtectionPolicy, bool)
1. MabProtectionPolicy.AsBasicProtectionPolicy() (BasicProtectionPolicy, bool)
1. MabProtectionPolicy.AsGenericProtectionPolicy() (*GenericProtectionPolicy, bool)
1. MabProtectionPolicy.AsMabProtectionPolicy() (*MabProtectionPolicy, bool)
1. MabProtectionPolicy.AsProtectionPolicy() (*ProtectionPolicy, bool)
1. MabProtectionPolicy.MarshalJSON() ([]byte, error)
1. NewBackupsClient(string) BackupsClient
1. NewBackupsClientWithBaseURI(string, string) BackupsClient
1. NewClientDiscoveryResponseIterator(ClientDiscoveryResponsePage) ClientDiscoveryResponseIterator
1. NewClientDiscoveryResponsePage(ClientDiscoveryResponse, func(context.Context, ClientDiscoveryResponse) (ClientDiscoveryResponse, error)) ClientDiscoveryResponsePage
1. NewEngineBaseResourceListIterator(EngineBaseResourceListPage) EngineBaseResourceListIterator
1. NewEngineBaseResourceListPage(EngineBaseResourceList, func(context.Context, EngineBaseResourceList) (EngineBaseResourceList, error)) EngineBaseResourceListPage
1. NewEnginesClient(string) EnginesClient
1. NewEnginesClientWithBaseURI(string, string) EnginesClient
1. NewExportJobsOperationResultsClient(string) ExportJobsOperationResultsClient
1. NewExportJobsOperationResultsClientWithBaseURI(string, string) ExportJobsOperationResultsClient
1. NewFeatureSupportClient(string) FeatureSupportClient
1. NewFeatureSupportClientWithBaseURI(string, string) FeatureSupportClient
1. NewItemLevelRecoveryConnectionsClient(string) ItemLevelRecoveryConnectionsClient
1. NewItemLevelRecoveryConnectionsClientWithBaseURI(string, string) ItemLevelRecoveryConnectionsClient
1. NewJobCancellationsClient(string) JobCancellationsClient
1. NewJobCancellationsClientWithBaseURI(string, string) JobCancellationsClient
1. NewJobDetailsClient(string) JobDetailsClient
1. NewJobDetailsClientWithBaseURI(string, string) JobDetailsClient
1. NewJobOperationResultsClient(string) JobOperationResultsClient
1. NewJobOperationResultsClientWithBaseURI(string, string) JobOperationResultsClient
1. NewJobResourceListIterator(JobResourceListPage) JobResourceListIterator
1. NewJobResourceListPage(JobResourceList, func(context.Context, JobResourceList) (JobResourceList, error)) JobResourceListPage
1. NewJobsClient(string) JobsClient
1. NewJobsClientWithBaseURI(string, string) JobsClient
1. NewJobsGroupClient(string) JobsGroupClient
1. NewJobsGroupClientWithBaseURI(string, string) JobsGroupClient
1. NewOperationClient(string) OperationClient
1. NewOperationClientWithBaseURI(string, string) OperationClient
1. NewOperationResultsClient(string) OperationResultsClient
1. NewOperationResultsClientWithBaseURI(string, string) OperationResultsClient
1. NewOperationStatusesClient(string) OperationStatusesClient
1. NewOperationStatusesClientWithBaseURI(string, string) OperationStatusesClient
1. NewOperationsClient(string) OperationsClient
1. NewOperationsClientWithBaseURI(string, string) OperationsClient
1. NewPoliciesClient(string) PoliciesClient
1. NewPoliciesClientWithBaseURI(string, string) PoliciesClient
1. NewProtectableContainerResourceListIterator(ProtectableContainerResourceListPage) ProtectableContainerResourceListIterator
1. NewProtectableContainerResourceListPage(ProtectableContainerResourceList, func(context.Context, ProtectableContainerResourceList) (ProtectableContainerResourceList, error)) ProtectableContainerResourceListPage
1. NewProtectableContainersClient(string) ProtectableContainersClient
1. NewProtectableContainersClientWithBaseURI(string, string) ProtectableContainersClient
1. NewProtectableItemsClient(string) ProtectableItemsClient
1. NewProtectableItemsClientWithBaseURI(string, string) ProtectableItemsClient
1. NewProtectedItemOperationResultsClient(string) ProtectedItemOperationResultsClient
1. NewProtectedItemOperationResultsClientWithBaseURI(string, string) ProtectedItemOperationResultsClient
1. NewProtectedItemOperationStatusesClient(string) ProtectedItemOperationStatusesClient
1. NewProtectedItemOperationStatusesClientWithBaseURI(string, string) ProtectedItemOperationStatusesClient
1. NewProtectedItemResourceListIterator(ProtectedItemResourceListPage) ProtectedItemResourceListIterator
1. NewProtectedItemResourceListPage(ProtectedItemResourceList, func(context.Context, ProtectedItemResourceList) (ProtectedItemResourceList, error)) ProtectedItemResourceListPage
1. NewProtectedItemsClient(string) ProtectedItemsClient
1. NewProtectedItemsClientWithBaseURI(string, string) ProtectedItemsClient
1. NewProtectedItemsGroupClient(string) ProtectedItemsGroupClient
1. NewProtectedItemsGroupClientWithBaseURI(string, string) ProtectedItemsGroupClient
1. NewProtectionContainerOperationResultsClient(string) ProtectionContainerOperationResultsClient
1. NewProtectionContainerOperationResultsClientWithBaseURI(string, string) ProtectionContainerOperationResultsClient
1. NewProtectionContainerRefreshOperationResultsClient(string) ProtectionContainerRefreshOperationResultsClient
1. NewProtectionContainerRefreshOperationResultsClientWithBaseURI(string, string) ProtectionContainerRefreshOperationResultsClient
1. NewProtectionContainerResourceListIterator(ProtectionContainerResourceListPage) ProtectionContainerResourceListIterator
1. NewProtectionContainerResourceListPage(ProtectionContainerResourceList, func(context.Context, ProtectionContainerResourceList) (ProtectionContainerResourceList, error)) ProtectionContainerResourceListPage
1. NewProtectionContainersClient(string) ProtectionContainersClient
1. NewProtectionContainersClientWithBaseURI(string, string) ProtectionContainersClient
1. NewProtectionContainersGroupClient(string) ProtectionContainersGroupClient
1. NewProtectionContainersGroupClientWithBaseURI(string, string) ProtectionContainersGroupClient
1. NewProtectionIntentClient(string) ProtectionIntentClient
1. NewProtectionIntentClientWithBaseURI(string, string) ProtectionIntentClient
1. NewProtectionIntentGroupClient(string) ProtectionIntentGroupClient
1. NewProtectionIntentGroupClientWithBaseURI(string, string) ProtectionIntentGroupClient
1. NewProtectionIntentResourceListIterator(ProtectionIntentResourceListPage) ProtectionIntentResourceListIterator
1. NewProtectionIntentResourceListPage(ProtectionIntentResourceList, func(context.Context, ProtectionIntentResourceList) (ProtectionIntentResourceList, error)) ProtectionIntentResourceListPage
1. NewProtectionPoliciesClient(string) ProtectionPoliciesClient
1. NewProtectionPoliciesClientWithBaseURI(string, string) ProtectionPoliciesClient
1. NewProtectionPolicyOperationResultsClient(string) ProtectionPolicyOperationResultsClient
1. NewProtectionPolicyOperationResultsClientWithBaseURI(string, string) ProtectionPolicyOperationResultsClient
1. NewProtectionPolicyOperationStatusesClient(string) ProtectionPolicyOperationStatusesClient
1. NewProtectionPolicyOperationStatusesClientWithBaseURI(string, string) ProtectionPolicyOperationStatusesClient
1. NewProtectionPolicyResourceListIterator(ProtectionPolicyResourceListPage) ProtectionPolicyResourceListIterator
1. NewProtectionPolicyResourceListPage(ProtectionPolicyResourceList, func(context.Context, ProtectionPolicyResourceList) (ProtectionPolicyResourceList, error)) ProtectionPolicyResourceListPage
1. NewRecoveryPointResourceListIterator(RecoveryPointResourceListPage) RecoveryPointResourceListIterator
1. NewRecoveryPointResourceListPage(RecoveryPointResourceList, func(context.Context, RecoveryPointResourceList) (RecoveryPointResourceList, error)) RecoveryPointResourceListPage
1. NewRecoveryPointsClient(string) RecoveryPointsClient
1. NewRecoveryPointsClientWithBaseURI(string, string) RecoveryPointsClient
1. NewResourceStorageConfigsClient(string) ResourceStorageConfigsClient
1. NewResourceStorageConfigsClientWithBaseURI(string, string) ResourceStorageConfigsClient
1. NewResourceVaultConfigsClient(string) ResourceVaultConfigsClient
1. NewResourceVaultConfigsClientWithBaseURI(string, string) ResourceVaultConfigsClient
1. NewRestoresClient(string) RestoresClient
1. NewRestoresClientWithBaseURI(string, string) RestoresClient
1. NewSecurityPINsClient(string) SecurityPINsClient
1. NewSecurityPINsClientWithBaseURI(string, string) SecurityPINsClient
1. NewStatusClient(string) StatusClient
1. NewStatusClientWithBaseURI(string, string) StatusClient
1. NewUsageSummariesClient(string) UsageSummariesClient
1. NewUsageSummariesClientWithBaseURI(string, string) UsageSummariesClient
1. NewWorkloadItemResourceListIterator(WorkloadItemResourceListPage) WorkloadItemResourceListIterator
1. NewWorkloadItemResourceListPage(WorkloadItemResourceList, func(context.Context, WorkloadItemResourceList) (WorkloadItemResourceList, error)) WorkloadItemResourceListPage
1. NewWorkloadItemsClient(string) WorkloadItemsClient
1. NewWorkloadItemsClientWithBaseURI(string, string) WorkloadItemsClient
1. NewWorkloadProtectableItemResourceListIterator(WorkloadProtectableItemResourceListPage) WorkloadProtectableItemResourceListIterator
1. NewWorkloadProtectableItemResourceListPage(WorkloadProtectableItemResourceList, func(context.Context, WorkloadProtectableItemResourceList) (WorkloadProtectableItemResourceList, error)) WorkloadProtectableItemResourceListPage
1. OperationClient.Validate(context.Context, string, string, BasicValidateOperationRequest) (ValidateOperationsResponse, error)
1. OperationClient.ValidatePreparer(context.Context, string, string, BasicValidateOperationRequest) (*http.Request, error)
1. OperationClient.ValidateResponder(*http.Response) (ValidateOperationsResponse, error)
1. OperationClient.ValidateSender(*http.Request) (*http.Response, error)
1. OperationResultInfo.AsBasicOperationResultInfoBase() (BasicOperationResultInfoBase, bool)
1. OperationResultInfo.AsExportJobsOperationResultInfo() (*ExportJobsOperationResultInfo, bool)
1. OperationResultInfo.AsOperationResultInfo() (*OperationResultInfo, bool)
1. OperationResultInfo.AsOperationResultInfoBase() (*OperationResultInfoBase, bool)
1. OperationResultInfo.MarshalJSON() ([]byte, error)
1. OperationResultInfoBase.AsBasicOperationResultInfoBase() (BasicOperationResultInfoBase, bool)
1. OperationResultInfoBase.AsExportJobsOperationResultInfo() (*ExportJobsOperationResultInfo, bool)
1. OperationResultInfoBase.AsOperationResultInfo() (*OperationResultInfo, bool)
1. OperationResultInfoBase.AsOperationResultInfoBase() (*OperationResultInfoBase, bool)
1. OperationResultInfoBase.MarshalJSON() ([]byte, error)
1. OperationResultInfoBaseResource.MarshalJSON() ([]byte, error)
1. OperationResultsClient.Get(context.Context, string, string, string) (autorest.Response, error)
1. OperationResultsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. OperationResultsClient.GetResponder(*http.Response) (autorest.Response, error)
1. OperationResultsClient.GetSender(*http.Request) (*http.Response, error)
1. OperationStatusExtendedInfo.AsBasicOperationStatusExtendedInfo() (BasicOperationStatusExtendedInfo, bool)
1. OperationStatusExtendedInfo.AsOperationStatusExtendedInfo() (*OperationStatusExtendedInfo, bool)
1. OperationStatusExtendedInfo.AsOperationStatusJobExtendedInfo() (*OperationStatusJobExtendedInfo, bool)
1. OperationStatusExtendedInfo.AsOperationStatusJobsExtendedInfo() (*OperationStatusJobsExtendedInfo, bool)
1. OperationStatusExtendedInfo.AsOperationStatusProvisionILRExtendedInfo() (*OperationStatusProvisionILRExtendedInfo, bool)
1. OperationStatusExtendedInfo.MarshalJSON() ([]byte, error)
1. OperationStatusJobExtendedInfo.AsBasicOperationStatusExtendedInfo() (BasicOperationStatusExtendedInfo, bool)
1. OperationStatusJobExtendedInfo.AsOperationStatusExtendedInfo() (*OperationStatusExtendedInfo, bool)
1. OperationStatusJobExtendedInfo.AsOperationStatusJobExtendedInfo() (*OperationStatusJobExtendedInfo, bool)
1. OperationStatusJobExtendedInfo.AsOperationStatusJobsExtendedInfo() (*OperationStatusJobsExtendedInfo, bool)
1. OperationStatusJobExtendedInfo.AsOperationStatusProvisionILRExtendedInfo() (*OperationStatusProvisionILRExtendedInfo, bool)
1. OperationStatusJobExtendedInfo.MarshalJSON() ([]byte, error)
1. OperationStatusJobsExtendedInfo.AsBasicOperationStatusExtendedInfo() (BasicOperationStatusExtendedInfo, bool)
1. OperationStatusJobsExtendedInfo.AsOperationStatusExtendedInfo() (*OperationStatusExtendedInfo, bool)
1. OperationStatusJobsExtendedInfo.AsOperationStatusJobExtendedInfo() (*OperationStatusJobExtendedInfo, bool)
1. OperationStatusJobsExtendedInfo.AsOperationStatusJobsExtendedInfo() (*OperationStatusJobsExtendedInfo, bool)
1. OperationStatusJobsExtendedInfo.AsOperationStatusProvisionILRExtendedInfo() (*OperationStatusProvisionILRExtendedInfo, bool)
1. OperationStatusJobsExtendedInfo.MarshalJSON() ([]byte, error)
1. OperationStatusProvisionILRExtendedInfo.AsBasicOperationStatusExtendedInfo() (BasicOperationStatusExtendedInfo, bool)
1. OperationStatusProvisionILRExtendedInfo.AsOperationStatusExtendedInfo() (*OperationStatusExtendedInfo, bool)
1. OperationStatusProvisionILRExtendedInfo.AsOperationStatusJobExtendedInfo() (*OperationStatusJobExtendedInfo, bool)
1. OperationStatusProvisionILRExtendedInfo.AsOperationStatusJobsExtendedInfo() (*OperationStatusJobsExtendedInfo, bool)
1. OperationStatusProvisionILRExtendedInfo.AsOperationStatusProvisionILRExtendedInfo() (*OperationStatusProvisionILRExtendedInfo, bool)
1. OperationStatusProvisionILRExtendedInfo.MarshalJSON() ([]byte, error)
1. OperationStatusesClient.Get(context.Context, string, string, string) (OperationStatus, error)
1. OperationStatusesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. OperationStatusesClient.GetResponder(*http.Response) (OperationStatus, error)
1. OperationStatusesClient.GetSender(*http.Request) (*http.Response, error)
1. OperationWorkerResponse.MarshalJSON() ([]byte, error)
1. OperationsClient.List(context.Context) (ClientDiscoveryResponsePage, error)
1. OperationsClient.ListComplete(context.Context) (ClientDiscoveryResponseIterator, error)
1. OperationsClient.ListPreparer(context.Context) (*http.Request, error)
1. OperationsClient.ListResponder(*http.Response) (ClientDiscoveryResponse, error)
1. OperationsClient.ListSender(*http.Request) (*http.Response, error)
1. PoliciesClient.List(context.Context, string, string, string) (ProtectionPolicyResourceListPage, error)
1. PoliciesClient.ListComplete(context.Context, string, string, string) (ProtectionPolicyResourceListIterator, error)
1. PoliciesClient.ListPreparer(context.Context, string, string, string) (*http.Request, error)
1. PoliciesClient.ListResponder(*http.Response) (ProtectionPolicyResourceList, error)
1. PoliciesClient.ListSender(*http.Request) (*http.Response, error)
1. PossibleAzureFileShareTypeValues() []AzureFileShareType
1. PossibleContainerTypeBasicProtectionContainerValues() []ContainerTypeBasicProtectionContainer
1. PossibleContainerTypeValues() []ContainerType
1. PossibleCopyOptionsValues() []CopyOptions
1. PossibleCreateModeValues() []CreateMode
1. PossibleDataSourceTypeValues() []DataSourceType
1. PossibleDayOfWeekValues() []DayOfWeek
1. PossibleEngineTypeValues() []EngineType
1. PossibleEnhancedSecurityStateValues() []EnhancedSecurityState
1. PossibleFabricNameValues() []FabricName
1. PossibleFeatureTypeValues() []FeatureType
1. PossibleHTTPStatusCodeValues() []HTTPStatusCode
1. PossibleHealthStateValues() []HealthState
1. PossibleHealthStatusValues() []HealthStatus
1. PossibleInquiryStatusValues() []InquiryStatus
1. PossibleIntentItemTypeValues() []IntentItemType
1. PossibleItemTypeValues() []ItemType
1. PossibleJobOperationTypeValues() []JobOperationType
1. PossibleJobStatusValues() []JobStatus
1. PossibleJobSupportedActionValues() []JobSupportedAction
1. PossibleJobTypeValues() []JobType
1. PossibleLastBackupStatusValues() []LastBackupStatus
1. PossibleMabServerTypeValues() []MabServerType
1. PossibleManagementTypeBasicProtectionPolicyValues() []ManagementTypeBasicProtectionPolicy
1. PossibleManagementTypeValues() []ManagementType
1. PossibleMonthOfYearValues() []MonthOfYear
1. PossibleObjectTypeBasicILRRequestValues() []ObjectTypeBasicILRRequest
1. PossibleObjectTypeBasicOperationResultInfoBaseValues() []ObjectTypeBasicOperationResultInfoBase
1. PossibleObjectTypeBasicRecoveryPointValues() []ObjectTypeBasicRecoveryPoint
1. PossibleObjectTypeBasicRequestValues() []ObjectTypeBasicRequest
1. PossibleObjectTypeBasicRestoreRequestValues() []ObjectTypeBasicRestoreRequest
1. PossibleObjectTypeBasicValidateOperationRequestValues() []ObjectTypeBasicValidateOperationRequest
1. PossibleObjectTypeValues() []ObjectType
1. PossibleOperationTypeValues() []OperationType
1. PossibleOverwriteOptionsValues() []OverwriteOptions
1. PossiblePolicyTypeValues() []PolicyType
1. PossibleProtectableContainerTypeValues() []ProtectableContainerType
1. PossibleProtectableItemTypeValues() []ProtectableItemType
1. PossibleProtectedItemHealthStatusValues() []ProtectedItemHealthStatus
1. PossibleProtectedItemStateValues() []ProtectedItemState
1. PossibleProtectedItemTypeValues() []ProtectedItemType
1. PossibleProtectionIntentItemTypeValues() []ProtectionIntentItemType
1. PossibleProtectionStateValues() []ProtectionState
1. PossibleProtectionStatusValues() []ProtectionStatus
1. PossibleRecoveryModeValues() []RecoveryMode
1. PossibleRecoveryPointTierStatusValues() []RecoveryPointTierStatus
1. PossibleRecoveryPointTierTypeValues() []RecoveryPointTierType
1. PossibleRecoveryTypeValues() []RecoveryType
1. PossibleResourceHealthStatusValues() []ResourceHealthStatus
1. PossibleRestorePointQueryTypeValues() []RestorePointQueryType
1. PossibleRestorePointTypeValues() []RestorePointType
1. PossibleRestoreRequestTypeValues() []RestoreRequestType
1. PossibleRetentionDurationTypeValues() []RetentionDurationType
1. PossibleRetentionPolicyTypeValues() []RetentionPolicyType
1. PossibleRetentionScheduleFormatValues() []RetentionScheduleFormat
1. PossibleSQLDataDirectoryTypeValues() []SQLDataDirectoryType
1. PossibleSchedulePolicyTypeValues() []SchedulePolicyType
1. PossibleScheduleRunTypeValues() []ScheduleRunType
1. PossibleSoftDeleteFeatureStateValues() []SoftDeleteFeatureState
1. PossibleStorageTypeStateValues() []StorageTypeState
1. PossibleStorageTypeValues() []StorageType
1. PossibleSupportStatusValues() []SupportStatus
1. PossibleTypeEnumValues() []TypeEnum
1. PossibleTypeValues() []Type
1. PossibleUsagesUnitValues() []UsagesUnit
1. PossibleValidationStatusValues() []ValidationStatus
1. PossibleWeekOfMonthValues() []WeekOfMonth
1. PossibleWorkloadItemTypeBasicWorkloadItemValues() []WorkloadItemTypeBasicWorkloadItem
1. PossibleWorkloadItemTypeValues() []WorkloadItemType
1. PossibleWorkloadTypeValues() []WorkloadType
1. ProtectableContainer.AsAzureStorageProtectableContainer() (*AzureStorageProtectableContainer, bool)
1. ProtectableContainer.AsAzureVMAppContainerProtectableContainer() (*AzureVMAppContainerProtectableContainer, bool)
1. ProtectableContainer.AsBasicProtectableContainer() (BasicProtectableContainer, bool)
1. ProtectableContainer.AsProtectableContainer() (*ProtectableContainer, bool)
1. ProtectableContainer.MarshalJSON() ([]byte, error)
1. ProtectableContainerResource.MarshalJSON() ([]byte, error)
1. ProtectableContainerResourceList.IsEmpty() bool
1. ProtectableContainerResourceListIterator.NotDone() bool
1. ProtectableContainerResourceListIterator.Response() ProtectableContainerResourceList
1. ProtectableContainerResourceListIterator.Value() ProtectableContainerResource
1. ProtectableContainerResourceListPage.NotDone() bool
1. ProtectableContainerResourceListPage.Response() ProtectableContainerResourceList
1. ProtectableContainerResourceListPage.Values() []ProtectableContainerResource
1. ProtectableContainersClient.List(context.Context, string, string, string, string) (ProtectableContainerResourceListPage, error)
1. ProtectableContainersClient.ListComplete(context.Context, string, string, string, string) (ProtectableContainerResourceListIterator, error)
1. ProtectableContainersClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectableContainersClient.ListResponder(*http.Response) (ProtectableContainerResourceList, error)
1. ProtectableContainersClient.ListSender(*http.Request) (*http.Response, error)
1. ProtectableItemsClient.List(context.Context, string, string, string, string) (WorkloadProtectableItemResourceListPage, error)
1. ProtectableItemsClient.ListComplete(context.Context, string, string, string, string) (WorkloadProtectableItemResourceListIterator, error)
1. ProtectableItemsClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectableItemsClient.ListResponder(*http.Response) (WorkloadProtectableItemResourceList, error)
1. ProtectableItemsClient.ListSender(*http.Request) (*http.Response, error)
1. ProtectedItem.AsAzureFileshareProtectedItem() (*AzureFileshareProtectedItem, bool)
1. ProtectedItem.AsAzureIaaSClassicComputeVMProtectedItem() (*AzureIaaSClassicComputeVMProtectedItem, bool)
1. ProtectedItem.AsAzureIaaSComputeVMProtectedItem() (*AzureIaaSComputeVMProtectedItem, bool)
1. ProtectedItem.AsAzureIaaSVMProtectedItem() (*AzureIaaSVMProtectedItem, bool)
1. ProtectedItem.AsAzureSQLProtectedItem() (*AzureSQLProtectedItem, bool)
1. ProtectedItem.AsAzureVMWorkloadProtectedItem() (*AzureVMWorkloadProtectedItem, bool)
1. ProtectedItem.AsAzureVMWorkloadSAPAseDatabaseProtectedItem() (*AzureVMWorkloadSAPAseDatabaseProtectedItem, bool)
1. ProtectedItem.AsAzureVMWorkloadSAPHanaDatabaseProtectedItem() (*AzureVMWorkloadSAPHanaDatabaseProtectedItem, bool)
1. ProtectedItem.AsAzureVMWorkloadSQLDatabaseProtectedItem() (*AzureVMWorkloadSQLDatabaseProtectedItem, bool)
1. ProtectedItem.AsBasicAzureIaaSVMProtectedItem() (BasicAzureIaaSVMProtectedItem, bool)
1. ProtectedItem.AsBasicAzureVMWorkloadProtectedItem() (BasicAzureVMWorkloadProtectedItem, bool)
1. ProtectedItem.AsBasicProtectedItem() (BasicProtectedItem, bool)
1. ProtectedItem.AsDPMProtectedItem() (*DPMProtectedItem, bool)
1. ProtectedItem.AsGenericProtectedItem() (*GenericProtectedItem, bool)
1. ProtectedItem.AsMabFileFolderProtectedItem() (*MabFileFolderProtectedItem, bool)
1. ProtectedItem.AsProtectedItem() (*ProtectedItem, bool)
1. ProtectedItem.MarshalJSON() ([]byte, error)
1. ProtectedItemOperationResultsClient.Get(context.Context, string, string, string, string, string, string) (ProtectedItemResource, error)
1. ProtectedItemOperationResultsClient.GetPreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. ProtectedItemOperationResultsClient.GetResponder(*http.Response) (ProtectedItemResource, error)
1. ProtectedItemOperationResultsClient.GetSender(*http.Request) (*http.Response, error)
1. ProtectedItemOperationStatusesClient.Get(context.Context, string, string, string, string, string, string) (OperationStatus, error)
1. ProtectedItemOperationStatusesClient.GetPreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. ProtectedItemOperationStatusesClient.GetResponder(*http.Response) (OperationStatus, error)
1. ProtectedItemOperationStatusesClient.GetSender(*http.Request) (*http.Response, error)
1. ProtectedItemResource.MarshalJSON() ([]byte, error)
1. ProtectedItemResourceList.IsEmpty() bool
1. ProtectedItemResourceListIterator.NotDone() bool
1. ProtectedItemResourceListIterator.Response() ProtectedItemResourceList
1. ProtectedItemResourceListIterator.Value() ProtectedItemResource
1. ProtectedItemResourceListPage.NotDone() bool
1. ProtectedItemResourceListPage.Response() ProtectedItemResourceList
1. ProtectedItemResourceListPage.Values() []ProtectedItemResource
1. ProtectedItemsClient.CreateOrUpdate(context.Context, string, string, string, string, string, ProtectedItemResource) (ProtectedItemResource, error)
1. ProtectedItemsClient.CreateOrUpdatePreparer(context.Context, string, string, string, string, string, ProtectedItemResource) (*http.Request, error)
1. ProtectedItemsClient.CreateOrUpdateResponder(*http.Response) (ProtectedItemResource, error)
1. ProtectedItemsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ProtectedItemsClient.Delete(context.Context, string, string, string, string, string) (autorest.Response, error)
1. ProtectedItemsClient.DeletePreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. ProtectedItemsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ProtectedItemsClient.DeleteSender(*http.Request) (*http.Response, error)
1. ProtectedItemsClient.Get(context.Context, string, string, string, string, string, string) (ProtectedItemResource, error)
1. ProtectedItemsClient.GetPreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. ProtectedItemsClient.GetResponder(*http.Response) (ProtectedItemResource, error)
1. ProtectedItemsClient.GetSender(*http.Request) (*http.Response, error)
1. ProtectedItemsGroupClient.List(context.Context, string, string, string, string) (ProtectedItemResourceListPage, error)
1. ProtectedItemsGroupClient.ListComplete(context.Context, string, string, string, string) (ProtectedItemResourceListIterator, error)
1. ProtectedItemsGroupClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectedItemsGroupClient.ListResponder(*http.Response) (ProtectedItemResourceList, error)
1. ProtectedItemsGroupClient.ListSender(*http.Request) (*http.Response, error)
1. ProtectionContainer.AsAzureBackupServerContainer() (*AzureBackupServerContainer, bool)
1. ProtectionContainer.AsAzureIaaSClassicComputeVMContainer() (*AzureIaaSClassicComputeVMContainer, bool)
1. ProtectionContainer.AsAzureIaaSComputeVMContainer() (*AzureIaaSComputeVMContainer, bool)
1. ProtectionContainer.AsAzureSQLAGWorkloadContainerProtectionContainer() (*AzureSQLAGWorkloadContainerProtectionContainer, bool)
1. ProtectionContainer.AsAzureSQLContainer() (*AzureSQLContainer, bool)
1. ProtectionContainer.AsAzureStorageContainer() (*AzureStorageContainer, bool)
1. ProtectionContainer.AsAzureVMAppContainerProtectionContainer() (*AzureVMAppContainerProtectionContainer, bool)
1. ProtectionContainer.AsAzureWorkloadContainer() (*AzureWorkloadContainer, bool)
1. ProtectionContainer.AsBasicAzureWorkloadContainer() (BasicAzureWorkloadContainer, bool)
1. ProtectionContainer.AsBasicDpmContainer() (BasicDpmContainer, bool)
1. ProtectionContainer.AsBasicIaaSVMContainer() (BasicIaaSVMContainer, bool)
1. ProtectionContainer.AsBasicProtectionContainer() (BasicProtectionContainer, bool)
1. ProtectionContainer.AsDpmContainer() (*DpmContainer, bool)
1. ProtectionContainer.AsGenericContainer() (*GenericContainer, bool)
1. ProtectionContainer.AsIaaSVMContainer() (*IaaSVMContainer, bool)
1. ProtectionContainer.AsMabContainer() (*MabContainer, bool)
1. ProtectionContainer.AsProtectionContainer() (*ProtectionContainer, bool)
1. ProtectionContainer.MarshalJSON() ([]byte, error)
1. ProtectionContainerOperationResultsClient.Get(context.Context, string, string, string, string, string) (ProtectionContainerResource, error)
1. ProtectionContainerOperationResultsClient.GetPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. ProtectionContainerOperationResultsClient.GetResponder(*http.Response) (ProtectionContainerResource, error)
1. ProtectionContainerOperationResultsClient.GetSender(*http.Request) (*http.Response, error)
1. ProtectionContainerRefreshOperationResultsClient.Get(context.Context, string, string, string, string) (autorest.Response, error)
1. ProtectionContainerRefreshOperationResultsClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectionContainerRefreshOperationResultsClient.GetResponder(*http.Response) (autorest.Response, error)
1. ProtectionContainerRefreshOperationResultsClient.GetSender(*http.Request) (*http.Response, error)
1. ProtectionContainerResource.MarshalJSON() ([]byte, error)
1. ProtectionContainerResourceList.IsEmpty() bool
1. ProtectionContainerResourceListIterator.NotDone() bool
1. ProtectionContainerResourceListIterator.Response() ProtectionContainerResourceList
1. ProtectionContainerResourceListIterator.Value() ProtectionContainerResource
1. ProtectionContainerResourceListPage.NotDone() bool
1. ProtectionContainerResourceListPage.Response() ProtectionContainerResourceList
1. ProtectionContainerResourceListPage.Values() []ProtectionContainerResource
1. ProtectionContainersClient.Get(context.Context, string, string, string, string) (ProtectionContainerResource, error)
1. ProtectionContainersClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectionContainersClient.GetResponder(*http.Response) (ProtectionContainerResource, error)
1. ProtectionContainersClient.GetSender(*http.Request) (*http.Response, error)
1. ProtectionContainersClient.Inquire(context.Context, string, string, string, string, string) (autorest.Response, error)
1. ProtectionContainersClient.InquirePreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. ProtectionContainersClient.InquireResponder(*http.Response) (autorest.Response, error)
1. ProtectionContainersClient.InquireSender(*http.Request) (*http.Response, error)
1. ProtectionContainersClient.Refresh(context.Context, string, string, string, string) (autorest.Response, error)
1. ProtectionContainersClient.RefreshPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectionContainersClient.RefreshResponder(*http.Response) (autorest.Response, error)
1. ProtectionContainersClient.RefreshSender(*http.Request) (*http.Response, error)
1. ProtectionContainersClient.Register(context.Context, string, string, string, string, ProtectionContainerResource) (ProtectionContainerResource, error)
1. ProtectionContainersClient.RegisterPreparer(context.Context, string, string, string, string, ProtectionContainerResource) (*http.Request, error)
1. ProtectionContainersClient.RegisterResponder(*http.Response) (ProtectionContainerResource, error)
1. ProtectionContainersClient.RegisterSender(*http.Request) (*http.Response, error)
1. ProtectionContainersClient.Unregister(context.Context, string, string, string, string) (autorest.Response, error)
1. ProtectionContainersClient.UnregisterPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectionContainersClient.UnregisterResponder(*http.Response) (autorest.Response, error)
1. ProtectionContainersClient.UnregisterSender(*http.Request) (*http.Response, error)
1. ProtectionContainersGroupClient.List(context.Context, string, string, string) (ProtectionContainerResourceListPage, error)
1. ProtectionContainersGroupClient.ListComplete(context.Context, string, string, string) (ProtectionContainerResourceListIterator, error)
1. ProtectionContainersGroupClient.ListPreparer(context.Context, string, string, string) (*http.Request, error)
1. ProtectionContainersGroupClient.ListResponder(*http.Response) (ProtectionContainerResourceList, error)
1. ProtectionContainersGroupClient.ListSender(*http.Request) (*http.Response, error)
1. ProtectionIntent.AsAzureRecoveryServiceVaultProtectionIntent() (*AzureRecoveryServiceVaultProtectionIntent, bool)
1. ProtectionIntent.AsAzureResourceProtectionIntent() (*AzureResourceProtectionIntent, bool)
1. ProtectionIntent.AsAzureWorkloadAutoProtectionIntent() (*AzureWorkloadAutoProtectionIntent, bool)
1. ProtectionIntent.AsAzureWorkloadSQLAutoProtectionIntent() (*AzureWorkloadSQLAutoProtectionIntent, bool)
1. ProtectionIntent.AsBasicAzureRecoveryServiceVaultProtectionIntent() (BasicAzureRecoveryServiceVaultProtectionIntent, bool)
1. ProtectionIntent.AsBasicAzureWorkloadAutoProtectionIntent() (BasicAzureWorkloadAutoProtectionIntent, bool)
1. ProtectionIntent.AsBasicProtectionIntent() (BasicProtectionIntent, bool)
1. ProtectionIntent.AsProtectionIntent() (*ProtectionIntent, bool)
1. ProtectionIntent.MarshalJSON() ([]byte, error)
1. ProtectionIntentClient.CreateOrUpdate(context.Context, string, string, string, string, ProtectionIntentResource) (ProtectionIntentResource, error)
1. ProtectionIntentClient.CreateOrUpdatePreparer(context.Context, string, string, string, string, ProtectionIntentResource) (*http.Request, error)
1. ProtectionIntentClient.CreateOrUpdateResponder(*http.Response) (ProtectionIntentResource, error)
1. ProtectionIntentClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ProtectionIntentClient.Delete(context.Context, string, string, string, string) (autorest.Response, error)
1. ProtectionIntentClient.DeletePreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectionIntentClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ProtectionIntentClient.DeleteSender(*http.Request) (*http.Response, error)
1. ProtectionIntentClient.Get(context.Context, string, string, string, string) (ProtectionIntentResource, error)
1. ProtectionIntentClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectionIntentClient.GetResponder(*http.Response) (ProtectionIntentResource, error)
1. ProtectionIntentClient.GetSender(*http.Request) (*http.Response, error)
1. ProtectionIntentClient.Validate(context.Context, string, PreValidateEnableBackupRequest) (PreValidateEnableBackupResponse, error)
1. ProtectionIntentClient.ValidatePreparer(context.Context, string, PreValidateEnableBackupRequest) (*http.Request, error)
1. ProtectionIntentClient.ValidateResponder(*http.Response) (PreValidateEnableBackupResponse, error)
1. ProtectionIntentClient.ValidateSender(*http.Request) (*http.Response, error)
1. ProtectionIntentGroupClient.List(context.Context, string, string, string, string) (ProtectionIntentResourceListPage, error)
1. ProtectionIntentGroupClient.ListComplete(context.Context, string, string, string, string) (ProtectionIntentResourceListIterator, error)
1. ProtectionIntentGroupClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectionIntentGroupClient.ListResponder(*http.Response) (ProtectionIntentResourceList, error)
1. ProtectionIntentGroupClient.ListSender(*http.Request) (*http.Response, error)
1. ProtectionIntentResource.MarshalJSON() ([]byte, error)
1. ProtectionIntentResourceList.IsEmpty() bool
1. ProtectionIntentResourceListIterator.NotDone() bool
1. ProtectionIntentResourceListIterator.Response() ProtectionIntentResourceList
1. ProtectionIntentResourceListIterator.Value() ProtectionIntentResource
1. ProtectionIntentResourceListPage.NotDone() bool
1. ProtectionIntentResourceListPage.Response() ProtectionIntentResourceList
1. ProtectionIntentResourceListPage.Values() []ProtectionIntentResource
1. ProtectionPoliciesClient.CreateOrUpdate(context.Context, string, string, string, ProtectionPolicyResource) (ProtectionPolicyResource, error)
1. ProtectionPoliciesClient.CreateOrUpdatePreparer(context.Context, string, string, string, ProtectionPolicyResource) (*http.Request, error)
1. ProtectionPoliciesClient.CreateOrUpdateResponder(*http.Response) (ProtectionPolicyResource, error)
1. ProtectionPoliciesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ProtectionPoliciesClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. ProtectionPoliciesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ProtectionPoliciesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ProtectionPoliciesClient.DeleteSender(*http.Request) (*http.Response, error)
1. ProtectionPoliciesClient.Get(context.Context, string, string, string) (ProtectionPolicyResource, error)
1. ProtectionPoliciesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ProtectionPoliciesClient.GetResponder(*http.Response) (ProtectionPolicyResource, error)
1. ProtectionPoliciesClient.GetSender(*http.Request) (*http.Response, error)
1. ProtectionPolicy.AsAzureFileShareProtectionPolicy() (*AzureFileShareProtectionPolicy, bool)
1. ProtectionPolicy.AsAzureIaaSVMProtectionPolicy() (*AzureIaaSVMProtectionPolicy, bool)
1. ProtectionPolicy.AsAzureSQLProtectionPolicy() (*AzureSQLProtectionPolicy, bool)
1. ProtectionPolicy.AsAzureVMWorkloadProtectionPolicy() (*AzureVMWorkloadProtectionPolicy, bool)
1. ProtectionPolicy.AsBasicProtectionPolicy() (BasicProtectionPolicy, bool)
1. ProtectionPolicy.AsGenericProtectionPolicy() (*GenericProtectionPolicy, bool)
1. ProtectionPolicy.AsMabProtectionPolicy() (*MabProtectionPolicy, bool)
1. ProtectionPolicy.AsProtectionPolicy() (*ProtectionPolicy, bool)
1. ProtectionPolicy.MarshalJSON() ([]byte, error)
1. ProtectionPolicyOperationResultsClient.Get(context.Context, string, string, string, string) (ProtectionPolicyResource, error)
1. ProtectionPolicyOperationResultsClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectionPolicyOperationResultsClient.GetResponder(*http.Response) (ProtectionPolicyResource, error)
1. ProtectionPolicyOperationResultsClient.GetSender(*http.Request) (*http.Response, error)
1. ProtectionPolicyOperationStatusesClient.Get(context.Context, string, string, string, string) (OperationStatus, error)
1. ProtectionPolicyOperationStatusesClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectionPolicyOperationStatusesClient.GetResponder(*http.Response) (OperationStatus, error)
1. ProtectionPolicyOperationStatusesClient.GetSender(*http.Request) (*http.Response, error)
1. ProtectionPolicyResource.MarshalJSON() ([]byte, error)
1. ProtectionPolicyResourceList.IsEmpty() bool
1. ProtectionPolicyResourceListIterator.NotDone() bool
1. ProtectionPolicyResourceListIterator.Response() ProtectionPolicyResourceList
1. ProtectionPolicyResourceListIterator.Value() ProtectionPolicyResource
1. ProtectionPolicyResourceListPage.NotDone() bool
1. ProtectionPolicyResourceListPage.Response() ProtectionPolicyResourceList
1. ProtectionPolicyResourceListPage.Values() []ProtectionPolicyResource
1. RecoveryPoint.AsAzureFileShareRecoveryPoint() (*AzureFileShareRecoveryPoint, bool)
1. RecoveryPoint.AsAzureWorkloadPointInTimeRecoveryPoint() (*AzureWorkloadPointInTimeRecoveryPoint, bool)
1. RecoveryPoint.AsAzureWorkloadRecoveryPoint() (*AzureWorkloadRecoveryPoint, bool)
1. RecoveryPoint.AsAzureWorkloadSAPHanaPointInTimeRecoveryPoint() (*AzureWorkloadSAPHanaPointInTimeRecoveryPoint, bool)
1. RecoveryPoint.AsAzureWorkloadSAPHanaRecoveryPoint() (*AzureWorkloadSAPHanaRecoveryPoint, bool)
1. RecoveryPoint.AsAzureWorkloadSQLPointInTimeRecoveryPoint() (*AzureWorkloadSQLPointInTimeRecoveryPoint, bool)
1. RecoveryPoint.AsAzureWorkloadSQLRecoveryPoint() (*AzureWorkloadSQLRecoveryPoint, bool)
1. RecoveryPoint.AsBasicAzureWorkloadPointInTimeRecoveryPoint() (BasicAzureWorkloadPointInTimeRecoveryPoint, bool)
1. RecoveryPoint.AsBasicAzureWorkloadRecoveryPoint() (BasicAzureWorkloadRecoveryPoint, bool)
1. RecoveryPoint.AsBasicAzureWorkloadSQLRecoveryPoint() (BasicAzureWorkloadSQLRecoveryPoint, bool)
1. RecoveryPoint.AsBasicRecoveryPoint() (BasicRecoveryPoint, bool)
1. RecoveryPoint.AsGenericRecoveryPoint() (*GenericRecoveryPoint, bool)
1. RecoveryPoint.AsIaasVMRecoveryPoint() (*IaasVMRecoveryPoint, bool)
1. RecoveryPoint.AsRecoveryPoint() (*RecoveryPoint, bool)
1. RecoveryPoint.MarshalJSON() ([]byte, error)
1. RecoveryPointResource.MarshalJSON() ([]byte, error)
1. RecoveryPointResourceList.IsEmpty() bool
1. RecoveryPointResourceListIterator.NotDone() bool
1. RecoveryPointResourceListIterator.Response() RecoveryPointResourceList
1. RecoveryPointResourceListIterator.Value() RecoveryPointResource
1. RecoveryPointResourceListPage.NotDone() bool
1. RecoveryPointResourceListPage.Response() RecoveryPointResourceList
1. RecoveryPointResourceListPage.Values() []RecoveryPointResource
1. RecoveryPointsClient.Get(context.Context, string, string, string, string, string, string) (RecoveryPointResource, error)
1. RecoveryPointsClient.GetPreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. RecoveryPointsClient.GetResponder(*http.Response) (RecoveryPointResource, error)
1. RecoveryPointsClient.GetSender(*http.Request) (*http.Response, error)
1. RecoveryPointsClient.List(context.Context, string, string, string, string, string, string) (RecoveryPointResourceListPage, error)
1. RecoveryPointsClient.ListComplete(context.Context, string, string, string, string, string, string) (RecoveryPointResourceListIterator, error)
1. RecoveryPointsClient.ListPreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. RecoveryPointsClient.ListResponder(*http.Response) (RecoveryPointResourceList, error)
1. RecoveryPointsClient.ListSender(*http.Request) (*http.Response, error)
1. Request.AsAzureFileShareBackupRequest() (*AzureFileShareBackupRequest, bool)
1. Request.AsAzureWorkloadBackupRequest() (*AzureWorkloadBackupRequest, bool)
1. Request.AsBasicRequest() (BasicRequest, bool)
1. Request.AsIaasVMBackupRequest() (*IaasVMBackupRequest, bool)
1. Request.AsRequest() (*Request, bool)
1. Request.MarshalJSON() ([]byte, error)
1. RequestResource.MarshalJSON() ([]byte, error)
1. ResourceConfigResource.MarshalJSON() ([]byte, error)
1. ResourceHealthDetails.MarshalJSON() ([]byte, error)
1. ResourceStorageConfigsClient.Get(context.Context, string, string) (ResourceConfigResource, error)
1. ResourceStorageConfigsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. ResourceStorageConfigsClient.GetResponder(*http.Response) (ResourceConfigResource, error)
1. ResourceStorageConfigsClient.GetSender(*http.Request) (*http.Response, error)
1. ResourceStorageConfigsClient.Patch(context.Context, string, string, ResourceConfigResource) (autorest.Response, error)
1. ResourceStorageConfigsClient.PatchPreparer(context.Context, string, string, ResourceConfigResource) (*http.Request, error)
1. ResourceStorageConfigsClient.PatchResponder(*http.Response) (autorest.Response, error)
1. ResourceStorageConfigsClient.PatchSender(*http.Request) (*http.Response, error)
1. ResourceStorageConfigsClient.Update(context.Context, string, string, ResourceConfigResource) (ResourceConfigResource, error)
1. ResourceStorageConfigsClient.UpdatePreparer(context.Context, string, string, ResourceConfigResource) (*http.Request, error)
1. ResourceStorageConfigsClient.UpdateResponder(*http.Response) (ResourceConfigResource, error)
1. ResourceStorageConfigsClient.UpdateSender(*http.Request) (*http.Response, error)
1. ResourceVaultConfigResource.MarshalJSON() ([]byte, error)
1. ResourceVaultConfigsClient.Get(context.Context, string, string) (ResourceVaultConfigResource, error)
1. ResourceVaultConfigsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. ResourceVaultConfigsClient.GetResponder(*http.Response) (ResourceVaultConfigResource, error)
1. ResourceVaultConfigsClient.GetSender(*http.Request) (*http.Response, error)
1. ResourceVaultConfigsClient.Put(context.Context, string, string, ResourceVaultConfigResource) (ResourceVaultConfigResource, error)
1. ResourceVaultConfigsClient.PutPreparer(context.Context, string, string, ResourceVaultConfigResource) (*http.Request, error)
1. ResourceVaultConfigsClient.PutResponder(*http.Response) (ResourceVaultConfigResource, error)
1. ResourceVaultConfigsClient.PutSender(*http.Request) (*http.Response, error)
1. ResourceVaultConfigsClient.Update(context.Context, string, string, ResourceVaultConfigResource) (ResourceVaultConfigResource, error)
1. ResourceVaultConfigsClient.UpdatePreparer(context.Context, string, string, ResourceVaultConfigResource) (*http.Request, error)
1. ResourceVaultConfigsClient.UpdateResponder(*http.Response) (ResourceVaultConfigResource, error)
1. ResourceVaultConfigsClient.UpdateSender(*http.Request) (*http.Response, error)
1. RestoreRequest.AsAzureFileShareRestoreRequest() (*AzureFileShareRestoreRequest, bool)
1. RestoreRequest.AsAzureWorkloadPointInTimeRestoreRequest() (*AzureWorkloadPointInTimeRestoreRequest, bool)
1. RestoreRequest.AsAzureWorkloadRestoreRequest() (*AzureWorkloadRestoreRequest, bool)
1. RestoreRequest.AsAzureWorkloadSAPHanaPointInTimeRestoreRequest() (*AzureWorkloadSAPHanaPointInTimeRestoreRequest, bool)
1. RestoreRequest.AsAzureWorkloadSAPHanaRestoreRequest() (*AzureWorkloadSAPHanaRestoreRequest, bool)
1. RestoreRequest.AsAzureWorkloadSQLPointInTimeRestoreRequest() (*AzureWorkloadSQLPointInTimeRestoreRequest, bool)
1. RestoreRequest.AsAzureWorkloadSQLRestoreRequest() (*AzureWorkloadSQLRestoreRequest, bool)
1. RestoreRequest.AsBasicAzureWorkloadRestoreRequest() (BasicAzureWorkloadRestoreRequest, bool)
1. RestoreRequest.AsBasicAzureWorkloadSAPHanaRestoreRequest() (BasicAzureWorkloadSAPHanaRestoreRequest, bool)
1. RestoreRequest.AsBasicAzureWorkloadSQLRestoreRequest() (BasicAzureWorkloadSQLRestoreRequest, bool)
1. RestoreRequest.AsBasicRestoreRequest() (BasicRestoreRequest, bool)
1. RestoreRequest.AsIaasVMRestoreRequest() (*IaasVMRestoreRequest, bool)
1. RestoreRequest.AsRestoreRequest() (*RestoreRequest, bool)
1. RestoreRequest.MarshalJSON() ([]byte, error)
1. RestoreRequestResource.MarshalJSON() ([]byte, error)
1. RestoresClient.Trigger(context.Context, string, string, string, string, string, string, RestoreRequestResource) (autorest.Response, error)
1. RestoresClient.TriggerPreparer(context.Context, string, string, string, string, string, string, RestoreRequestResource) (*http.Request, error)
1. RestoresClient.TriggerResponder(*http.Response) (autorest.Response, error)
1. RestoresClient.TriggerSender(*http.Request) (*http.Response, error)
1. RetentionPolicy.AsBasicRetentionPolicy() (BasicRetentionPolicy, bool)
1. RetentionPolicy.AsLongTermRetentionPolicy() (*LongTermRetentionPolicy, bool)
1. RetentionPolicy.AsRetentionPolicy() (*RetentionPolicy, bool)
1. RetentionPolicy.AsSimpleRetentionPolicy() (*SimpleRetentionPolicy, bool)
1. RetentionPolicy.MarshalJSON() ([]byte, error)
1. SchedulePolicy.AsBasicSchedulePolicy() (BasicSchedulePolicy, bool)
1. SchedulePolicy.AsLogSchedulePolicy() (*LogSchedulePolicy, bool)
1. SchedulePolicy.AsLongTermSchedulePolicy() (*LongTermSchedulePolicy, bool)
1. SchedulePolicy.AsSchedulePolicy() (*SchedulePolicy, bool)
1. SchedulePolicy.AsSimpleSchedulePolicy() (*SimpleSchedulePolicy, bool)
1. SchedulePolicy.MarshalJSON() ([]byte, error)
1. SecurityPINsClient.Get(context.Context, string, string) (TokenInformation, error)
1. SecurityPINsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. SecurityPINsClient.GetResponder(*http.Response) (TokenInformation, error)
1. SecurityPINsClient.GetSender(*http.Request) (*http.Response, error)
1. SimpleRetentionPolicy.AsBasicRetentionPolicy() (BasicRetentionPolicy, bool)
1. SimpleRetentionPolicy.AsLongTermRetentionPolicy() (*LongTermRetentionPolicy, bool)
1. SimpleRetentionPolicy.AsRetentionPolicy() (*RetentionPolicy, bool)
1. SimpleRetentionPolicy.AsSimpleRetentionPolicy() (*SimpleRetentionPolicy, bool)
1. SimpleRetentionPolicy.MarshalJSON() ([]byte, error)
1. SimpleSchedulePolicy.AsBasicSchedulePolicy() (BasicSchedulePolicy, bool)
1. SimpleSchedulePolicy.AsLogSchedulePolicy() (*LogSchedulePolicy, bool)
1. SimpleSchedulePolicy.AsLongTermSchedulePolicy() (*LongTermSchedulePolicy, bool)
1. SimpleSchedulePolicy.AsSchedulePolicy() (*SchedulePolicy, bool)
1. SimpleSchedulePolicy.AsSimpleSchedulePolicy() (*SimpleSchedulePolicy, bool)
1. SimpleSchedulePolicy.MarshalJSON() ([]byte, error)
1. StatusClient.Get(context.Context, string, StatusRequest) (StatusResponse, error)
1. StatusClient.GetPreparer(context.Context, string, StatusRequest) (*http.Request, error)
1. StatusClient.GetResponder(*http.Response) (StatusResponse, error)
1. StatusClient.GetSender(*http.Request) (*http.Response, error)
1. UsageSummariesClient.List(context.Context, string, string, string, string) (ManagementUsageList, error)
1. UsageSummariesClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. UsageSummariesClient.ListResponder(*http.Response) (ManagementUsageList, error)
1. UsageSummariesClient.ListSender(*http.Request) (*http.Response, error)
1. ValidateIaasVMRestoreOperationRequest.AsBasicValidateOperationRequest() (BasicValidateOperationRequest, bool)
1. ValidateIaasVMRestoreOperationRequest.AsBasicValidateRestoreOperationRequest() (BasicValidateRestoreOperationRequest, bool)
1. ValidateIaasVMRestoreOperationRequest.AsValidateIaasVMRestoreOperationRequest() (*ValidateIaasVMRestoreOperationRequest, bool)
1. ValidateIaasVMRestoreOperationRequest.AsValidateOperationRequest() (*ValidateOperationRequest, bool)
1. ValidateIaasVMRestoreOperationRequest.AsValidateRestoreOperationRequest() (*ValidateRestoreOperationRequest, bool)
1. ValidateIaasVMRestoreOperationRequest.MarshalJSON() ([]byte, error)
1. ValidateOperationRequest.AsBasicValidateOperationRequest() (BasicValidateOperationRequest, bool)
1. ValidateOperationRequest.AsBasicValidateRestoreOperationRequest() (BasicValidateRestoreOperationRequest, bool)
1. ValidateOperationRequest.AsValidateIaasVMRestoreOperationRequest() (*ValidateIaasVMRestoreOperationRequest, bool)
1. ValidateOperationRequest.AsValidateOperationRequest() (*ValidateOperationRequest, bool)
1. ValidateOperationRequest.AsValidateRestoreOperationRequest() (*ValidateRestoreOperationRequest, bool)
1. ValidateOperationRequest.MarshalJSON() ([]byte, error)
1. ValidateRestoreOperationRequest.AsBasicValidateOperationRequest() (BasicValidateOperationRequest, bool)
1. ValidateRestoreOperationRequest.AsBasicValidateRestoreOperationRequest() (BasicValidateRestoreOperationRequest, bool)
1. ValidateRestoreOperationRequest.AsValidateIaasVMRestoreOperationRequest() (*ValidateIaasVMRestoreOperationRequest, bool)
1. ValidateRestoreOperationRequest.AsValidateOperationRequest() (*ValidateOperationRequest, bool)
1. ValidateRestoreOperationRequest.AsValidateRestoreOperationRequest() (*ValidateRestoreOperationRequest, bool)
1. ValidateRestoreOperationRequest.MarshalJSON() ([]byte, error)
1. WorkloadItem.AsAzureVMWorkloadItem() (*AzureVMWorkloadItem, bool)
1. WorkloadItem.AsAzureVMWorkloadSAPAseDatabaseWorkloadItem() (*AzureVMWorkloadSAPAseDatabaseWorkloadItem, bool)
1. WorkloadItem.AsAzureVMWorkloadSAPAseSystemWorkloadItem() (*AzureVMWorkloadSAPAseSystemWorkloadItem, bool)
1. WorkloadItem.AsAzureVMWorkloadSAPHanaDatabaseWorkloadItem() (*AzureVMWorkloadSAPHanaDatabaseWorkloadItem, bool)
1. WorkloadItem.AsAzureVMWorkloadSAPHanaSystemWorkloadItem() (*AzureVMWorkloadSAPHanaSystemWorkloadItem, bool)
1. WorkloadItem.AsAzureVMWorkloadSQLDatabaseWorkloadItem() (*AzureVMWorkloadSQLDatabaseWorkloadItem, bool)
1. WorkloadItem.AsAzureVMWorkloadSQLInstanceWorkloadItem() (*AzureVMWorkloadSQLInstanceWorkloadItem, bool)
1. WorkloadItem.AsBasicAzureVMWorkloadItem() (BasicAzureVMWorkloadItem, bool)
1. WorkloadItem.AsBasicWorkloadItem() (BasicWorkloadItem, bool)
1. WorkloadItem.AsWorkloadItem() (*WorkloadItem, bool)
1. WorkloadItem.MarshalJSON() ([]byte, error)
1. WorkloadItemResource.MarshalJSON() ([]byte, error)
1. WorkloadItemResourceList.IsEmpty() bool
1. WorkloadItemResourceListIterator.NotDone() bool
1. WorkloadItemResourceListIterator.Response() WorkloadItemResourceList
1. WorkloadItemResourceListIterator.Value() WorkloadItemResource
1. WorkloadItemResourceListPage.NotDone() bool
1. WorkloadItemResourceListPage.Response() WorkloadItemResourceList
1. WorkloadItemResourceListPage.Values() []WorkloadItemResource
1. WorkloadItemsClient.List(context.Context, string, string, string, string, string, string) (WorkloadItemResourceListPage, error)
1. WorkloadItemsClient.ListComplete(context.Context, string, string, string, string, string, string) (WorkloadItemResourceListIterator, error)
1. WorkloadItemsClient.ListPreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. WorkloadItemsClient.ListResponder(*http.Response) (WorkloadItemResourceList, error)
1. WorkloadItemsClient.ListSender(*http.Request) (*http.Response, error)
1. WorkloadProtectableItem.AsAzureFileShareProtectableItem() (*AzureFileShareProtectableItem, bool)
1. WorkloadProtectableItem.AsAzureIaaSClassicComputeVMProtectableItem() (*AzureIaaSClassicComputeVMProtectableItem, bool)
1. WorkloadProtectableItem.AsAzureIaaSComputeVMProtectableItem() (*AzureIaaSComputeVMProtectableItem, bool)
1. WorkloadProtectableItem.AsAzureVMWorkloadProtectableItem() (*AzureVMWorkloadProtectableItem, bool)
1. WorkloadProtectableItem.AsAzureVMWorkloadSAPAseSystemProtectableItem() (*AzureVMWorkloadSAPAseSystemProtectableItem, bool)
1. WorkloadProtectableItem.AsAzureVMWorkloadSAPHanaDatabaseProtectableItem() (*AzureVMWorkloadSAPHanaDatabaseProtectableItem, bool)
1. WorkloadProtectableItem.AsAzureVMWorkloadSAPHanaSystemProtectableItem() (*AzureVMWorkloadSAPHanaSystemProtectableItem, bool)
1. WorkloadProtectableItem.AsAzureVMWorkloadSQLAvailabilityGroupProtectableItem() (*AzureVMWorkloadSQLAvailabilityGroupProtectableItem, bool)
1. WorkloadProtectableItem.AsAzureVMWorkloadSQLDatabaseProtectableItem() (*AzureVMWorkloadSQLDatabaseProtectableItem, bool)
1. WorkloadProtectableItem.AsAzureVMWorkloadSQLInstanceProtectableItem() (*AzureVMWorkloadSQLInstanceProtectableItem, bool)
1. WorkloadProtectableItem.AsBasicAzureVMWorkloadProtectableItem() (BasicAzureVMWorkloadProtectableItem, bool)
1. WorkloadProtectableItem.AsBasicIaaSVMProtectableItem() (BasicIaaSVMProtectableItem, bool)
1. WorkloadProtectableItem.AsBasicWorkloadProtectableItem() (BasicWorkloadProtectableItem, bool)
1. WorkloadProtectableItem.AsIaaSVMProtectableItem() (*IaaSVMProtectableItem, bool)
1. WorkloadProtectableItem.AsWorkloadProtectableItem() (*WorkloadProtectableItem, bool)
1. WorkloadProtectableItem.MarshalJSON() ([]byte, error)
1. WorkloadProtectableItemResource.MarshalJSON() ([]byte, error)
1. WorkloadProtectableItemResourceList.IsEmpty() bool
1. WorkloadProtectableItemResourceListIterator.NotDone() bool
1. WorkloadProtectableItemResourceListIterator.Response() WorkloadProtectableItemResourceList
1. WorkloadProtectableItemResourceListIterator.Value() WorkloadProtectableItemResource
1. WorkloadProtectableItemResourceListPage.NotDone() bool
1. WorkloadProtectableItemResourceListPage.Response() WorkloadProtectableItemResourceList
1. WorkloadProtectableItemResourceListPage.Values() []WorkloadProtectableItemResource

### Struct Changes

#### Removed Structs

1. AzureBackupGoalFeatureSupportRequest
1. AzureBackupServerContainer
1. AzureBackupServerEngine
1. AzureFileShareBackupRequest
1. AzureFileShareProtectableItem
1. AzureFileShareProtectionPolicy
1. AzureFileShareProvisionILRRequest
1. AzureFileShareRecoveryPoint
1. AzureFileShareRestoreRequest
1. AzureFileshareProtectedItem
1. AzureFileshareProtectedItemExtendedInfo
1. AzureIaaSClassicComputeVMContainer
1. AzureIaaSClassicComputeVMProtectableItem
1. AzureIaaSClassicComputeVMProtectedItem
1. AzureIaaSComputeVMContainer
1. AzureIaaSComputeVMProtectableItem
1. AzureIaaSComputeVMProtectedItem
1. AzureIaaSVMErrorInfo
1. AzureIaaSVMHealthDetails
1. AzureIaaSVMJob
1. AzureIaaSVMJobExtendedInfo
1. AzureIaaSVMJobTaskDetails
1. AzureIaaSVMProtectedItem
1. AzureIaaSVMProtectedItemExtendedInfo
1. AzureIaaSVMProtectionPolicy
1. AzureRecoveryServiceVaultProtectionIntent
1. AzureResourceProtectionIntent
1. AzureSQLAGWorkloadContainerProtectionContainer
1. AzureSQLContainer
1. AzureSQLProtectedItem
1. AzureSQLProtectedItemExtendedInfo
1. AzureSQLProtectionPolicy
1. AzureStorageContainer
1. AzureStorageErrorInfo
1. AzureStorageJob
1. AzureStorageJobExtendedInfo
1. AzureStorageJobTaskDetails
1. AzureStorageProtectableContainer
1. AzureVMAppContainerProtectableContainer
1. AzureVMAppContainerProtectionContainer
1. AzureVMResourceFeatureSupportRequest
1. AzureVMResourceFeatureSupportResponse
1. AzureVMWorkloadItem
1. AzureVMWorkloadProtectableItem
1. AzureVMWorkloadProtectedItem
1. AzureVMWorkloadProtectedItemExtendedInfo
1. AzureVMWorkloadProtectionPolicy
1. AzureVMWorkloadSAPAseDatabaseProtectedItem
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem
1. AzureVMWorkloadSAPAseSystemProtectableItem
1. AzureVMWorkloadSAPAseSystemWorkloadItem
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem
1. AzureVMWorkloadSAPHanaDatabaseProtectedItem
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem
1. AzureVMWorkloadSAPHanaSystemProtectableItem
1. AzureVMWorkloadSAPHanaSystemWorkloadItem
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem
1. AzureVMWorkloadSQLDatabaseProtectableItem
1. AzureVMWorkloadSQLDatabaseProtectedItem
1. AzureVMWorkloadSQLDatabaseWorkloadItem
1. AzureVMWorkloadSQLInstanceProtectableItem
1. AzureVMWorkloadSQLInstanceWorkloadItem
1. AzureWorkloadAutoProtectionIntent
1. AzureWorkloadBackupRequest
1. AzureWorkloadContainer
1. AzureWorkloadContainerExtendedInfo
1. AzureWorkloadErrorInfo
1. AzureWorkloadJob
1. AzureWorkloadJobExtendedInfo
1. AzureWorkloadJobTaskDetails
1. AzureWorkloadPointInTimeRecoveryPoint
1. AzureWorkloadPointInTimeRestoreRequest
1. AzureWorkloadRecoveryPoint
1. AzureWorkloadRestoreRequest
1. AzureWorkloadSAPHanaPointInTimeRecoveryPoint
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest
1. AzureWorkloadSAPHanaRecoveryPoint
1. AzureWorkloadSAPHanaRestoreRequest
1. AzureWorkloadSQLAutoProtectionIntent
1. AzureWorkloadSQLPointInTimeRecoveryPoint
1. AzureWorkloadSQLPointInTimeRestoreRequest
1. AzureWorkloadSQLRecoveryPoint
1. AzureWorkloadSQLRecoveryPointExtendedInfo
1. AzureWorkloadSQLRestoreRequest
1. BEKDetails
1. BMSBackupEngineQueryObject
1. BMSBackupEnginesQueryObject
1. BMSBackupSummariesQueryObject
1. BMSContainerQueryObject
1. BMSContainersInquiryQueryObject
1. BMSPOQueryObject
1. BMSRPQueryObject
1. BMSRefreshContainersQueryObject
1. BMSWorkloadItemQueryObject
1. BackupsClient
1. ClientDiscoveryDisplay
1. ClientDiscoveryForLogSpecification
1. ClientDiscoveryForProperties
1. ClientDiscoveryForServiceSpecification
1. ClientDiscoveryResponse
1. ClientDiscoveryResponseIterator
1. ClientDiscoveryResponsePage
1. ClientDiscoveryValueForSingleAPI
1. ClientScriptForConnect
1. ContainerIdentityInfo
1. DPMContainerExtendedInfo
1. DPMProtectedItem
1. DPMProtectedItemExtendedInfo
1. DailyRetentionFormat
1. DailyRetentionSchedule
1. Day
1. DiskExclusionProperties
1. DiskInformation
1. DistributedNodesInfo
1. DpmBackupEngine
1. DpmContainer
1. DpmErrorInfo
1. DpmJob
1. DpmJobExtendedInfo
1. DpmJobTaskDetails
1. EncryptionDetails
1. EngineBase
1. EngineBaseResource
1. EngineBaseResourceList
1. EngineBaseResourceListIterator
1. EngineBaseResourceListPage
1. EngineExtendedInfo
1. EnginesClient
1. ErrorDetail
1. ExportJobsOperationResultInfo
1. ExportJobsOperationResultsClient
1. ExtendedProperties
1. FeatureSupportClient
1. FeatureSupportRequest
1. GenericContainer
1. GenericContainerExtendedInfo
1. GenericProtectedItem
1. GenericProtectionPolicy
1. GenericRecoveryPoint
1. GetProtectedItemQueryObject
1. ILRRequest
1. ILRRequestResource
1. IaaSVMContainer
1. IaaSVMProtectableItem
1. IaasVMBackupRequest
1. IaasVMILRRegistrationRequest
1. IaasVMRecoveryPoint
1. IaasVMRestoreRequest
1. InquiryInfo
1. InquiryValidation
1. InstantItemRecoveryTarget
1. InstantRPAdditionalDetails
1. ItemLevelRecoveryConnectionsClient
1. Job
1. JobCancellationsClient
1. JobDetailsClient
1. JobOperationResultsClient
1. JobQueryObject
1. JobResource
1. JobResourceList
1. JobResourceListIterator
1. JobResourceListPage
1. JobsClient
1. JobsGroupClient
1. KEKDetails
1. KPIResourceHealthDetails
1. KeyAndSecretDetails
1. LogSchedulePolicy
1. LongTermRetentionPolicy
1. LongTermSchedulePolicy
1. MABContainerHealthDetails
1. MabContainer
1. MabContainerExtendedInfo
1. MabErrorInfo
1. MabFileFolderProtectedItem
1. MabFileFolderProtectedItemExtendedInfo
1. MabJob
1. MabJobExtendedInfo
1. MabJobTaskDetails
1. MabProtectionPolicy
1. ManagementUsage
1. ManagementUsageList
1. MonthlyRetentionSchedule
1. NameInfo
1. OperationClient
1. OperationResultInfo
1. OperationResultInfoBase
1. OperationResultInfoBaseResource
1. OperationResultsClient
1. OperationStatusExtendedInfo
1. OperationStatusJobExtendedInfo
1. OperationStatusJobsExtendedInfo
1. OperationStatusProvisionILRExtendedInfo
1. OperationStatusesClient
1. OperationWorkerResponse
1. OperationsClient
1. PointInTimeRange
1. PoliciesClient
1. PreBackupValidation
1. PreValidateEnableBackupRequest
1. PreValidateEnableBackupResponse
1. ProtectableContainer
1. ProtectableContainerResource
1. ProtectableContainerResourceList
1. ProtectableContainerResourceListIterator
1. ProtectableContainerResourceListPage
1. ProtectableContainersClient
1. ProtectableItemsClient
1. ProtectedItem
1. ProtectedItemOperationResultsClient
1. ProtectedItemOperationStatusesClient
1. ProtectedItemQueryObject
1. ProtectedItemResource
1. ProtectedItemResourceList
1. ProtectedItemResourceListIterator
1. ProtectedItemResourceListPage
1. ProtectedItemsClient
1. ProtectedItemsGroupClient
1. ProtectionContainer
1. ProtectionContainerOperationResultsClient
1. ProtectionContainerRefreshOperationResultsClient
1. ProtectionContainerResource
1. ProtectionContainerResourceList
1. ProtectionContainerResourceListIterator
1. ProtectionContainerResourceListPage
1. ProtectionContainersClient
1. ProtectionContainersGroupClient
1. ProtectionIntent
1. ProtectionIntentClient
1. ProtectionIntentGroupClient
1. ProtectionIntentQueryObject
1. ProtectionIntentResource
1. ProtectionIntentResourceList
1. ProtectionIntentResourceListIterator
1. ProtectionIntentResourceListPage
1. ProtectionPoliciesClient
1. ProtectionPolicy
1. ProtectionPolicyOperationResultsClient
1. ProtectionPolicyOperationStatusesClient
1. ProtectionPolicyQueryObject
1. ProtectionPolicyResource
1. ProtectionPolicyResourceList
1. ProtectionPolicyResourceListIterator
1. ProtectionPolicyResourceListPage
1. RecoveryPoint
1. RecoveryPointDiskConfiguration
1. RecoveryPointResource
1. RecoveryPointResourceList
1. RecoveryPointResourceListIterator
1. RecoveryPointResourceListPage
1. RecoveryPointTierInformation
1. RecoveryPointsClient
1. Request
1. RequestResource
1. ResourceConfig
1. ResourceConfigResource
1. ResourceHealthDetails
1. ResourceList
1. ResourceStorageConfigsClient
1. ResourceVaultConfig
1. ResourceVaultConfigResource
1. ResourceVaultConfigsClient
1. RestoreFileSpecs
1. RestoreRequest
1. RestoreRequestResource
1. RestoresClient
1. RetentionDuration
1. RetentionPolicy
1. SQLDataDirectory
1. SQLDataDirectoryMapping
1. SchedulePolicy
1. SecurityPINsClient
1. Settings
1. SimpleRetentionPolicy
1. SimpleSchedulePolicy
1. StatusClient
1. StatusRequest
1. StatusResponse
1. SubProtectionPolicy
1. TargetAFSRestoreInfo
1. TargetRestoreInfo
1. TokenInformation
1. UsageSummariesClient
1. ValidateIaasVMRestoreOperationRequest
1. ValidateOperationRequest
1. ValidateOperationResponse
1. ValidateOperationsResponse
1. ValidateRestoreOperationRequest
1. WeeklyRetentionFormat
1. WeeklyRetentionSchedule
1. WorkloadInquiryDetails
1. WorkloadItem
1. WorkloadItemResource
1. WorkloadItemResourceList
1. WorkloadItemResourceListIterator
1. WorkloadItemResourceListPage
1. WorkloadItemsClient
1. WorkloadProtectableItem
1. WorkloadProtectableItemResource
1. WorkloadProtectableItemResourceList
1. WorkloadProtectableItemResourceListIterator
1. WorkloadProtectableItemResourceListPage
1. YearlyRetentionSchedule

#### Removed Struct Fields

1. OperationStatus.Properties

### Signature Changes

#### Const Types

1. Invalid changed type from AzureFileShareType to OperationStatusValues

## Additive Changes

### New Constants

1. OperationStatusValues.Canceled
1. OperationStatusValues.Failed
1. OperationStatusValues.InProgress
1. OperationStatusValues.Succeeded
