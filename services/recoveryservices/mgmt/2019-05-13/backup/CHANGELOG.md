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
1. EngineType.BackupEngineTypeAzureBackupServerEngine
1. EngineType.BackupEngineTypeBackupEngineBase
1. EngineType.BackupEngineTypeDpmBackupEngine
1. FabricName.FabricNameAzure
1. FabricName.FabricNameInvalid
1. FeatureType.FeatureTypeAzureBackupGoals
1. FeatureType.FeatureTypeAzureVMResourceBackup
1. FeatureType.FeatureTypeFeatureSupportRequest
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
1. ObjectTypeBasicILRRequest.ObjectTypeAzureFileShareProvisionILRRequest
1. ObjectTypeBasicILRRequest.ObjectTypeILRRequest
1. ObjectTypeBasicILRRequest.ObjectTypeIaasVMILRRegistrationRequest
1. ObjectTypeBasicOperationStatusExtendedInfo.ObjectTypeOperationStatusExtendedInfo
1. ObjectTypeBasicOperationStatusExtendedInfo.ObjectTypeOperationStatusJobExtendedInfo
1. ObjectTypeBasicOperationStatusExtendedInfo.ObjectTypeOperationStatusJobsExtendedInfo
1. ObjectTypeBasicOperationStatusExtendedInfo.ObjectTypeOperationStatusProvisionILRExtendedInfo
1. ObjectTypeBasicRequest.ObjectTypeAzureFileShareBackupRequest
1. ObjectTypeBasicRequest.ObjectTypeAzureWorkloadBackupRequest
1. ObjectTypeBasicRequest.ObjectTypeBackupRequest
1. ObjectTypeBasicRequest.ObjectTypeIaasVMBackupRequest
1. OperationStatusValues.OperationStatusValuesCanceled
1. OperationStatusValues.OperationStatusValuesFailed
1. OperationStatusValues.OperationStatusValuesInProgress
1. OperationStatusValues.OperationStatusValuesInvalid
1. OperationStatusValues.OperationStatusValuesSucceeded
1. OperationType.OperationTypeInvalid
1. OperationType.OperationTypeRegister
1. OperationType.OperationTypeReregister
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
1. ProtectionIntentItemType.ProtectionIntentItemTypeAzureResourceItem
1. ProtectionIntentItemType.ProtectionIntentItemTypeAzureWorkloadAutoProtectionIntent
1. ProtectionIntentItemType.ProtectionIntentItemTypeAzureWorkloadSQLAutoProtectionIntent
1. ProtectionIntentItemType.ProtectionIntentItemTypeProtectionIntent
1. ProtectionIntentItemType.ProtectionIntentItemTypeRecoveryServiceVaultItem
1. ProtectionStatus.ProtectionStatusInvalid
1. ProtectionStatus.ProtectionStatusNotProtected
1. ProtectionStatus.ProtectionStatusProtected
1. ProtectionStatus.ProtectionStatusProtecting
1. ProtectionStatus.ProtectionStatusProtectionFailed
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

### Removed Funcs

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
1. *OperationStatus.UnmarshalJSON([]byte) error
1. *ProtectableContainerResource.UnmarshalJSON([]byte) error
1. *ProtectableContainerResourceListIterator.Next() error
1. *ProtectableContainerResourceListIterator.NextWithContext(context.Context) error
1. *ProtectableContainerResourceListPage.Next() error
1. *ProtectableContainerResourceListPage.NextWithContext(context.Context) error
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
1. *RequestResource.UnmarshalJSON([]byte) error
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
1. AzureFileShareProvisionILRRequest.AsAzureFileShareProvisionILRRequest() (*AzureFileShareProvisionILRRequest, bool)
1. AzureFileShareProvisionILRRequest.AsBasicILRRequest() (BasicILRRequest, bool)
1. AzureFileShareProvisionILRRequest.AsILRRequest() (*ILRRequest, bool)
1. AzureFileShareProvisionILRRequest.AsIaasVMILRRegistrationRequest() (*IaasVMILRRegistrationRequest, bool)
1. AzureFileShareProvisionILRRequest.MarshalJSON() ([]byte, error)
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
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureRecoveryServiceVaultProtectionIntent() (*AzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureResourceProtectionIntent() (*AzureResourceProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureWorkloadAutoProtectionIntent() (*AzureWorkloadAutoProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureWorkloadSQLAutoProtectionIntent() (*AzureWorkloadSQLAutoProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsBasicAzureRecoveryServiceVaultProtectionIntent() (BasicAzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsBasicAzureWorkloadAutoProtectionIntent() (BasicAzureWorkloadAutoProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsBasicProtectionIntent() (BasicProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsProtectionIntent() (*ProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.MarshalJSON() ([]byte, error)
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
1. InquiryValidation.MarshalJSON() ([]byte, error)
1. ItemLevelRecoveryConnectionsClient.Provision(context.Context, string, string, string, string, string, string, ILRRequestResource) (autorest.Response, error)
1. ItemLevelRecoveryConnectionsClient.ProvisionPreparer(context.Context, string, string, string, string, string, string, ILRRequestResource) (*http.Request, error)
1. ItemLevelRecoveryConnectionsClient.ProvisionResponder(*http.Response) (autorest.Response, error)
1. ItemLevelRecoveryConnectionsClient.ProvisionSender(*http.Request) (*http.Response, error)
1. ItemLevelRecoveryConnectionsClient.Revoke(context.Context, string, string, string, string, string, string) (autorest.Response, error)
1. ItemLevelRecoveryConnectionsClient.RevokePreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. ItemLevelRecoveryConnectionsClient.RevokeResponder(*http.Response) (autorest.Response, error)
1. ItemLevelRecoveryConnectionsClient.RevokeSender(*http.Request) (*http.Response, error)
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
1. NewBackupsClient(string) BackupsClient
1. NewBackupsClientWithBaseURI(string, string) BackupsClient
1. NewClientDiscoveryResponseIterator(ClientDiscoveryResponsePage) ClientDiscoveryResponseIterator
1. NewClientDiscoveryResponsePage(ClientDiscoveryResponse, func(context.Context, ClientDiscoveryResponse) (ClientDiscoveryResponse, error)) ClientDiscoveryResponsePage
1. NewEngineBaseResourceListIterator(EngineBaseResourceListPage) EngineBaseResourceListIterator
1. NewEngineBaseResourceListPage(EngineBaseResourceList, func(context.Context, EngineBaseResourceList) (EngineBaseResourceList, error)) EngineBaseResourceListPage
1. NewEnginesClient(string) EnginesClient
1. NewEnginesClientWithBaseURI(string, string) EnginesClient
1. NewFeatureSupportClient(string) FeatureSupportClient
1. NewFeatureSupportClientWithBaseURI(string, string) FeatureSupportClient
1. NewItemLevelRecoveryConnectionsClient(string) ItemLevelRecoveryConnectionsClient
1. NewItemLevelRecoveryConnectionsClientWithBaseURI(string, string) ItemLevelRecoveryConnectionsClient
1. NewOperationResultsClient(string) OperationResultsClient
1. NewOperationResultsClientWithBaseURI(string, string) OperationResultsClient
1. NewOperationStatusesClient(string) OperationStatusesClient
1. NewOperationStatusesClientWithBaseURI(string, string) OperationStatusesClient
1. NewOperationsClient(string) OperationsClient
1. NewOperationsClientWithBaseURI(string, string) OperationsClient
1. NewProtectableContainerResourceListIterator(ProtectableContainerResourceListPage) ProtectableContainerResourceListIterator
1. NewProtectableContainerResourceListPage(ProtectableContainerResourceList, func(context.Context, ProtectableContainerResourceList) (ProtectableContainerResourceList, error)) ProtectableContainerResourceListPage
1. NewProtectableContainersClient(string) ProtectableContainersClient
1. NewProtectableContainersClientWithBaseURI(string, string) ProtectableContainersClient
1. NewProtectableItemsClient(string) ProtectableItemsClient
1. NewProtectableItemsClientWithBaseURI(string, string) ProtectableItemsClient
1. NewProtectedItemOperationStatusesClient(string) ProtectedItemOperationStatusesClient
1. NewProtectedItemOperationStatusesClientWithBaseURI(string, string) ProtectedItemOperationStatusesClient
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
1. NewProtectionPolicyOperationStatusesClient(string) ProtectionPolicyOperationStatusesClient
1. NewProtectionPolicyOperationStatusesClientWithBaseURI(string, string) ProtectionPolicyOperationStatusesClient
1. NewResourceStorageConfigsClient(string) ResourceStorageConfigsClient
1. NewResourceStorageConfigsClientWithBaseURI(string, string) ResourceStorageConfigsClient
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
1. OperationsClient.List(context.Context) (ClientDiscoveryResponsePage, error)
1. OperationsClient.ListComplete(context.Context) (ClientDiscoveryResponseIterator, error)
1. OperationsClient.ListPreparer(context.Context) (*http.Request, error)
1. OperationsClient.ListResponder(*http.Response) (ClientDiscoveryResponse, error)
1. OperationsClient.ListSender(*http.Request) (*http.Response, error)
1. PossibleAzureFileShareTypeValues() []AzureFileShareType
1. PossibleContainerTypeBasicProtectionContainerValues() []ContainerTypeBasicProtectionContainer
1. PossibleContainerTypeValues() []ContainerType
1. PossibleEngineTypeValues() []EngineType
1. PossibleFabricNameValues() []FabricName
1. PossibleFeatureTypeValues() []FeatureType
1. PossibleInquiryStatusValues() []InquiryStatus
1. PossibleIntentItemTypeValues() []IntentItemType
1. PossibleItemTypeValues() []ItemType
1. PossibleObjectTypeBasicILRRequestValues() []ObjectTypeBasicILRRequest
1. PossibleObjectTypeBasicOperationStatusExtendedInfoValues() []ObjectTypeBasicOperationStatusExtendedInfo
1. PossibleObjectTypeBasicRequestValues() []ObjectTypeBasicRequest
1. PossibleOperationStatusValuesValues() []OperationStatusValues
1. PossibleOperationTypeValues() []OperationType
1. PossibleProtectableContainerTypeValues() []ProtectableContainerType
1. PossibleProtectableItemTypeValues() []ProtectableItemType
1. PossibleProtectionIntentItemTypeValues() []ProtectionIntentItemType
1. PossibleProtectionStatusValues() []ProtectionStatus
1. PossibleSupportStatusValues() []SupportStatus
1. PossibleTypeEnumValues() []TypeEnum
1. PossibleTypeValues() []Type
1. PossibleUsagesUnitValues() []UsagesUnit
1. PossibleValidationStatusValues() []ValidationStatus
1. PossibleWorkloadItemTypeBasicWorkloadItemValues() []WorkloadItemTypeBasicWorkloadItem
1. PossibleWorkloadItemTypeValues() []WorkloadItemType
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
1. ProtectedItemOperationStatusesClient.Get(context.Context, string, string, string, string, string, string) (OperationStatus, error)
1. ProtectedItemOperationStatusesClient.GetPreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. ProtectedItemOperationStatusesClient.GetResponder(*http.Response) (OperationStatus, error)
1. ProtectedItemOperationStatusesClient.GetSender(*http.Request) (*http.Response, error)
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
1. ProtectionPoliciesClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. ProtectionPoliciesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ProtectionPoliciesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ProtectionPoliciesClient.DeleteSender(*http.Request) (*http.Response, error)
1. ProtectionPolicyOperationStatusesClient.Get(context.Context, string, string, string, string) (OperationStatus, error)
1. ProtectionPolicyOperationStatusesClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectionPolicyOperationStatusesClient.GetResponder(*http.Response) (OperationStatus, error)
1. ProtectionPolicyOperationStatusesClient.GetSender(*http.Request) (*http.Response, error)
1. Request.AsAzureFileShareBackupRequest() (*AzureFileShareBackupRequest, bool)
1. Request.AsAzureWorkloadBackupRequest() (*AzureWorkloadBackupRequest, bool)
1. Request.AsBasicRequest() (BasicRequest, bool)
1. Request.AsIaasVMBackupRequest() (*IaasVMBackupRequest, bool)
1. Request.AsRequest() (*Request, bool)
1. Request.MarshalJSON() ([]byte, error)
1. RequestResource.MarshalJSON() ([]byte, error)
1. ResourceConfigResource.MarshalJSON() ([]byte, error)
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
1. SecurityPINsClient.Get(context.Context, string, string) (TokenInformation, error)
1. SecurityPINsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. SecurityPINsClient.GetResponder(*http.Response) (TokenInformation, error)
1. SecurityPINsClient.GetSender(*http.Request) (*http.Response, error)
1. StatusClient.Get(context.Context, string, StatusRequest) (StatusResponse, error)
1. StatusClient.GetPreparer(context.Context, string, StatusRequest) (*http.Request, error)
1. StatusClient.GetResponder(*http.Response) (StatusResponse, error)
1. StatusClient.GetSender(*http.Request) (*http.Response, error)
1. UsageSummariesClient.List(context.Context, string, string, string, string) (ManagementUsageList, error)
1. UsageSummariesClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. UsageSummariesClient.ListResponder(*http.Response) (ManagementUsageList, error)
1. UsageSummariesClient.ListSender(*http.Request) (*http.Response, error)
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
1. AzureFileShareProvisionILRRequest
1. AzureIaaSClassicComputeVMContainer
1. AzureIaaSClassicComputeVMProtectableItem
1. AzureIaaSComputeVMContainer
1. AzureIaaSComputeVMProtectableItem
1. AzureRecoveryServiceVaultProtectionIntent
1. AzureResourceProtectionIntent
1. AzureSQLAGWorkloadContainerProtectionContainer
1. AzureSQLContainer
1. AzureStorageContainer
1. AzureStorageProtectableContainer
1. AzureVMAppContainerProtectableContainer
1. AzureVMAppContainerProtectionContainer
1. AzureVMResourceFeatureSupportRequest
1. AzureVMResourceFeatureSupportResponse
1. AzureVMWorkloadItem
1. AzureVMWorkloadProtectableItem
1. AzureVMWorkloadSAPAseDatabaseWorkloadItem
1. AzureVMWorkloadSAPAseSystemProtectableItem
1. AzureVMWorkloadSAPAseSystemWorkloadItem
1. AzureVMWorkloadSAPHanaDatabaseProtectableItem
1. AzureVMWorkloadSAPHanaDatabaseWorkloadItem
1. AzureVMWorkloadSAPHanaSystemProtectableItem
1. AzureVMWorkloadSAPHanaSystemWorkloadItem
1. AzureVMWorkloadSQLAvailabilityGroupProtectableItem
1. AzureVMWorkloadSQLDatabaseProtectableItem
1. AzureVMWorkloadSQLDatabaseWorkloadItem
1. AzureVMWorkloadSQLInstanceProtectableItem
1. AzureVMWorkloadSQLInstanceWorkloadItem
1. AzureWorkloadAutoProtectionIntent
1. AzureWorkloadBackupRequest
1. AzureWorkloadContainer
1. AzureWorkloadContainerExtendedInfo
1. AzureWorkloadSQLAutoProtectionIntent
1. BMSBackupEngineQueryObject
1. BMSBackupEnginesQueryObject
1. BMSBackupSummariesQueryObject
1. BMSContainerQueryObject
1. BMSContainersInquiryQueryObject
1. BMSPOQueryObject
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
1. DistributedNodesInfo
1. DpmBackupEngine
1. DpmContainer
1. EngineBase
1. EngineBaseResource
1. EngineBaseResourceList
1. EngineBaseResourceListIterator
1. EngineBaseResourceListPage
1. EngineExtendedInfo
1. EnginesClient
1. FeatureSupportClient
1. FeatureSupportRequest
1. GenericContainer
1. GenericContainerExtendedInfo
1. ILRRequest
1. ILRRequestResource
1. IaaSVMContainer
1. IaaSVMProtectableItem
1. IaasVMBackupRequest
1. IaasVMILRRegistrationRequest
1. InquiryInfo
1. InquiryValidation
1. InstantItemRecoveryTarget
1. InstantRPAdditionalDetails
1. ItemLevelRecoveryConnectionsClient
1. MABContainerHealthDetails
1. MabContainer
1. MabContainerExtendedInfo
1. ManagementUsage
1. ManagementUsageList
1. NameInfo
1. OperationResultsClient
1. OperationStatus
1. OperationStatusError
1. OperationStatusExtendedInfo
1. OperationStatusJobExtendedInfo
1. OperationStatusJobsExtendedInfo
1. OperationStatusProvisionILRExtendedInfo
1. OperationStatusesClient
1. OperationsClient
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
1. ProtectedItemOperationStatusesClient
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
1. ProtectionPolicyOperationStatusesClient
1. Request
1. RequestResource
1. ResourceConfig
1. ResourceConfigResource
1. ResourceStorageConfigsClient
1. SecurityPINsClient
1. StatusClient
1. StatusRequest
1. StatusResponse
1. TokenInformation
1. UsageSummariesClient
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

### Signature Changes

#### Const Types

1. Invalid changed type from AzureFileShareType to CopyOptions

## Additive Changes

### New Constants

1. CopyOptions.CreateCopy
1. CopyOptions.FailOnConflict
1. CopyOptions.Overwrite
1. CopyOptions.Skip
