# Unreleased

## Breaking Changes

### Removed Constants

1. FabricName.FabricNameAzure
1. FabricName.FabricNameInvalid
1. FeatureType.FeatureTypeAzureBackupGoals
1. FeatureType.FeatureTypeAzureVMResourceBackup
1. FeatureType.FeatureTypeFeatureSupportRequest
1. IntentItemType.IntentItemTypeInvalid
1. IntentItemType.IntentItemTypeSQLAvailabilityGroupContainer
1. IntentItemType.IntentItemTypeSQLInstance
1. ObjectTypeBasicCrrAccessToken.ObjectTypeBasicCrrAccessTokenObjectTypeCrrAccessToken
1. ObjectTypeBasicCrrAccessToken.ObjectTypeBasicCrrAccessTokenObjectTypeWorkloadCrrAccessToken
1. ObjectTypeBasicOperationStatusExtendedInfo.ObjectTypeBasicOperationStatusExtendedInfoObjectTypeOperationStatusRecoveryPointExtendedInfo
1. ProtectionIntentItemType.ProtectionIntentItemTypeAzureResourceItem
1. ProtectionIntentItemType.ProtectionIntentItemTypeAzureWorkloadAutoProtectionIntent
1. ProtectionIntentItemType.ProtectionIntentItemTypeAzureWorkloadSQLAutoProtectionIntent
1. ProtectionIntentItemType.ProtectionIntentItemTypeProtectionIntent
1. ProtectionIntentItemType.ProtectionIntentItemTypeRecoveryServiceVaultItem
1. SupportStatus.SupportStatusDefaultOFF
1. SupportStatus.SupportStatusDefaultON
1. SupportStatus.SupportStatusInvalid
1. SupportStatus.SupportStatusNotSupported
1. SupportStatus.SupportStatusSupported
1. Type.TypeBackupProtectedItemCountSummary
1. Type.TypeBackupProtectionContainerCountSummary
1. TypeEnum.TypeEnumCopyOnlyFull
1. TypeEnum.TypeEnumDifferential
1. TypeEnum.TypeEnumFull
1. TypeEnum.TypeEnumIncremental
1. TypeEnum.TypeEnumInvalid
1. TypeEnum.TypeEnumLog
1. UsagesUnit.UsagesUnitBytes
1. UsagesUnit.UsagesUnitBytesPerSecond
1. UsagesUnit.UsagesUnitCount
1. UsagesUnit.UsagesUnitCountPerSecond
1. UsagesUnit.UsagesUnitPercent
1. UsagesUnit.UsagesUnitSeconds
1. ValidationStatus.ValidationStatusFailed
1. ValidationStatus.ValidationStatusInvalid
1. ValidationStatus.ValidationStatusSucceeded

### Removed Funcs

1. *ClientDiscoveryResponseIterator.Next() error
1. *ClientDiscoveryResponseIterator.NextWithContext(context.Context) error
1. *ClientDiscoveryResponsePage.Next() error
1. *ClientDiscoveryResponsePage.NextWithContext(context.Context) error
1. *CrossRegionRestoreRequest.UnmarshalJSON([]byte) error
1. *CrossRegionRestoreTriggerFuture.UnmarshalJSON([]byte) error
1. *CrrAccessTokenResource.UnmarshalJSON([]byte) error
1. *OperationStatusRecoveryPointExtendedInfo.UnmarshalJSON([]byte) error
1. *ProtectionIntentResource.UnmarshalJSON([]byte) error
1. *ProtectionIntentResourceListIterator.Next() error
1. *ProtectionIntentResourceListIterator.NextWithContext(context.Context) error
1. *ProtectionIntentResourceListPage.Next() error
1. *ProtectionIntentResourceListPage.NextWithContext(context.Context) error
1. AADPropertiesResource.MarshalJSON() ([]byte, error)
1. AadPropertiesClient.Get(context.Context, string, string) (AADPropertiesResource, error)
1. AadPropertiesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. AadPropertiesClient.GetResponder(*http.Response) (AADPropertiesResource, error)
1. AadPropertiesClient.GetSender(*http.Request) (*http.Response, error)
1. AzureBackupGoalFeatureSupportRequest.AsAzureBackupGoalFeatureSupportRequest() (*AzureBackupGoalFeatureSupportRequest, bool)
1. AzureBackupGoalFeatureSupportRequest.AsAzureVMResourceFeatureSupportRequest() (*AzureVMResourceFeatureSupportRequest, bool)
1. AzureBackupGoalFeatureSupportRequest.AsBasicFeatureSupportRequest() (BasicFeatureSupportRequest, bool)
1. AzureBackupGoalFeatureSupportRequest.AsFeatureSupportRequest() (*FeatureSupportRequest, bool)
1. AzureBackupGoalFeatureSupportRequest.MarshalJSON() ([]byte, error)
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
1. AzureVMResourceFeatureSupportRequest.AsAzureBackupGoalFeatureSupportRequest() (*AzureBackupGoalFeatureSupportRequest, bool)
1. AzureVMResourceFeatureSupportRequest.AsAzureVMResourceFeatureSupportRequest() (*AzureVMResourceFeatureSupportRequest, bool)
1. AzureVMResourceFeatureSupportRequest.AsBasicFeatureSupportRequest() (BasicFeatureSupportRequest, bool)
1. AzureVMResourceFeatureSupportRequest.AsFeatureSupportRequest() (*FeatureSupportRequest, bool)
1. AzureVMResourceFeatureSupportRequest.MarshalJSON() ([]byte, error)
1. AzureWorkloadAutoProtectionIntent.AsAzureRecoveryServiceVaultProtectionIntent() (*AzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsAzureResourceProtectionIntent() (*AzureResourceProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsAzureWorkloadAutoProtectionIntent() (*AzureWorkloadAutoProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsAzureWorkloadSQLAutoProtectionIntent() (*AzureWorkloadSQLAutoProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsBasicAzureRecoveryServiceVaultProtectionIntent() (BasicAzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsBasicAzureWorkloadAutoProtectionIntent() (BasicAzureWorkloadAutoProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsBasicProtectionIntent() (BasicProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.AsProtectionIntent() (*ProtectionIntent, bool)
1. AzureWorkloadAutoProtectionIntent.MarshalJSON() ([]byte, error)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureRecoveryServiceVaultProtectionIntent() (*AzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureResourceProtectionIntent() (*AzureResourceProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureWorkloadAutoProtectionIntent() (*AzureWorkloadAutoProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsAzureWorkloadSQLAutoProtectionIntent() (*AzureWorkloadSQLAutoProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsBasicAzureRecoveryServiceVaultProtectionIntent() (BasicAzureRecoveryServiceVaultProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsBasicAzureWorkloadAutoProtectionIntent() (BasicAzureWorkloadAutoProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsBasicProtectionIntent() (BasicProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.AsProtectionIntent() (*ProtectionIntent, bool)
1. AzureWorkloadSQLAutoProtectionIntent.MarshalJSON() ([]byte, error)
1. AzureWorkloadSQLRecoveryPointExtendedInfo.MarshalJSON() ([]byte, error)
1. ClientDiscoveryResponse.IsEmpty() bool
1. ClientDiscoveryResponseIterator.NotDone() bool
1. ClientDiscoveryResponseIterator.Response() ClientDiscoveryResponse
1. ClientDiscoveryResponseIterator.Value() ClientDiscoveryValueForSingleAPI
1. ClientDiscoveryResponsePage.NotDone() bool
1. ClientDiscoveryResponsePage.Response() ClientDiscoveryResponse
1. ClientDiscoveryResponsePage.Values() []ClientDiscoveryValueForSingleAPI
1. CrossRegionRestoreClient.Trigger(context.Context, string, CrossRegionRestoreRequest) (CrossRegionRestoreTriggerFuture, error)
1. CrossRegionRestoreClient.TriggerPreparer(context.Context, string, CrossRegionRestoreRequest) (*http.Request, error)
1. CrossRegionRestoreClient.TriggerResponder(*http.Response) (autorest.Response, error)
1. CrossRegionRestoreClient.TriggerSender(*http.Request) (CrossRegionRestoreTriggerFuture, error)
1. CrossRegionRestoreRequestResource.MarshalJSON() ([]byte, error)
1. CrrAccessToken.AsBasicCrrAccessToken() (BasicCrrAccessToken, bool)
1. CrrAccessToken.AsCrrAccessToken() (*CrrAccessToken, bool)
1. CrrAccessToken.AsWorkloadCrrAccessToken() (*WorkloadCrrAccessToken, bool)
1. CrrAccessToken.MarshalJSON() ([]byte, error)
1. CrrAccessTokenResource.MarshalJSON() ([]byte, error)
1. CrrJobDetailsClient.Get(context.Context, string, CrrJobRequest) (JobResource, error)
1. CrrJobDetailsClient.GetPreparer(context.Context, string, CrrJobRequest) (*http.Request, error)
1. CrrJobDetailsClient.GetResponder(*http.Response) (JobResource, error)
1. CrrJobDetailsClient.GetSender(*http.Request) (*http.Response, error)
1. CrrJobRequestResource.MarshalJSON() ([]byte, error)
1. CrrJobsClient.List(context.Context, string, CrrJobRequest, string, string) (JobResourceListPage, error)
1. CrrJobsClient.ListComplete(context.Context, string, CrrJobRequest, string, string) (JobResourceListIterator, error)
1. CrrJobsClient.ListPreparer(context.Context, string, CrrJobRequest, string, string) (*http.Request, error)
1. CrrJobsClient.ListResponder(*http.Response) (JobResourceList, error)
1. CrrJobsClient.ListSender(*http.Request) (*http.Response, error)
1. CrrOperationResultsClient.Get(context.Context, string, string) (autorest.Response, error)
1. CrrOperationResultsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. CrrOperationResultsClient.GetResponder(*http.Response) (autorest.Response, error)
1. CrrOperationResultsClient.GetSender(*http.Request) (*http.Response, error)
1. CrrOperationStatusClient.Get(context.Context, string, string) (OperationStatus, error)
1. CrrOperationStatusClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. CrrOperationStatusClient.GetResponder(*http.Response) (OperationStatus, error)
1. CrrOperationStatusClient.GetSender(*http.Request) (*http.Response, error)
1. FeatureSupportClient.Validate(context.Context, string, BasicFeatureSupportRequest) (AzureVMResourceFeatureSupportResponse, error)
1. FeatureSupportClient.ValidatePreparer(context.Context, string, BasicFeatureSupportRequest) (*http.Request, error)
1. FeatureSupportClient.ValidateResponder(*http.Response) (AzureVMResourceFeatureSupportResponse, error)
1. FeatureSupportClient.ValidateSender(*http.Request) (*http.Response, error)
1. FeatureSupportRequest.AsAzureBackupGoalFeatureSupportRequest() (*AzureBackupGoalFeatureSupportRequest, bool)
1. FeatureSupportRequest.AsAzureVMResourceFeatureSupportRequest() (*AzureVMResourceFeatureSupportRequest, bool)
1. FeatureSupportRequest.AsBasicFeatureSupportRequest() (BasicFeatureSupportRequest, bool)
1. FeatureSupportRequest.AsFeatureSupportRequest() (*FeatureSupportRequest, bool)
1. FeatureSupportRequest.MarshalJSON() ([]byte, error)
1. NewAadPropertiesClient(string) AadPropertiesClient
1. NewAadPropertiesClientWithBaseURI(string, string) AadPropertiesClient
1. NewClientDiscoveryResponseIterator(ClientDiscoveryResponsePage) ClientDiscoveryResponseIterator
1. NewClientDiscoveryResponsePage(ClientDiscoveryResponse, func(context.Context, ClientDiscoveryResponse) (ClientDiscoveryResponse, error)) ClientDiscoveryResponsePage
1. NewCrossRegionRestoreClient(string) CrossRegionRestoreClient
1. NewCrossRegionRestoreClientWithBaseURI(string, string) CrossRegionRestoreClient
1. NewCrrJobDetailsClient(string) CrrJobDetailsClient
1. NewCrrJobDetailsClientWithBaseURI(string, string) CrrJobDetailsClient
1. NewCrrJobsClient(string) CrrJobsClient
1. NewCrrJobsClientWithBaseURI(string, string) CrrJobsClient
1. NewCrrOperationResultsClient(string) CrrOperationResultsClient
1. NewCrrOperationResultsClientWithBaseURI(string, string) CrrOperationResultsClient
1. NewCrrOperationStatusClient(string) CrrOperationStatusClient
1. NewCrrOperationStatusClientWithBaseURI(string, string) CrrOperationStatusClient
1. NewFeatureSupportClient(string) FeatureSupportClient
1. NewFeatureSupportClientWithBaseURI(string, string) FeatureSupportClient
1. NewOperationsClient(string) OperationsClient
1. NewOperationsClientWithBaseURI(string, string) OperationsClient
1. NewProtectedItemsCrrClient(string) ProtectedItemsCrrClient
1. NewProtectedItemsCrrClientWithBaseURI(string, string) ProtectedItemsCrrClient
1. NewProtectionIntentClient(string) ProtectionIntentClient
1. NewProtectionIntentClientWithBaseURI(string, string) ProtectionIntentClient
1. NewProtectionIntentGroupClient(string) ProtectionIntentGroupClient
1. NewProtectionIntentGroupClientWithBaseURI(string, string) ProtectionIntentGroupClient
1. NewProtectionIntentResourceListIterator(ProtectionIntentResourceListPage) ProtectionIntentResourceListIterator
1. NewProtectionIntentResourceListPage(ProtectionIntentResourceList, func(context.Context, ProtectionIntentResourceList) (ProtectionIntentResourceList, error)) ProtectionIntentResourceListPage
1. NewRecoveryPointsCrrClient(string) RecoveryPointsCrrClient
1. NewRecoveryPointsCrrClientWithBaseURI(string, string) RecoveryPointsCrrClient
1. NewResourceStorageConfigsClient(string) ResourceStorageConfigsClient
1. NewResourceStorageConfigsClientWithBaseURI(string, string) ResourceStorageConfigsClient
1. NewStatusClient(string) StatusClient
1. NewStatusClientWithBaseURI(string, string) StatusClient
1. NewUsageSummariesCRRClient(string) UsageSummariesCRRClient
1. NewUsageSummariesCRRClientWithBaseURI(string, string) UsageSummariesCRRClient
1. NewUsageSummariesClient(string) UsageSummariesClient
1. NewUsageSummariesClientWithBaseURI(string, string) UsageSummariesClient
1. OperationStatusExtendedInfo.AsOperationStatusRecoveryPointExtendedInfo() (*OperationStatusRecoveryPointExtendedInfo, bool)
1. OperationStatusJobExtendedInfo.AsOperationStatusRecoveryPointExtendedInfo() (*OperationStatusRecoveryPointExtendedInfo, bool)
1. OperationStatusJobsExtendedInfo.AsOperationStatusRecoveryPointExtendedInfo() (*OperationStatusRecoveryPointExtendedInfo, bool)
1. OperationStatusProvisionILRExtendedInfo.AsOperationStatusRecoveryPointExtendedInfo() (*OperationStatusRecoveryPointExtendedInfo, bool)
1. OperationStatusRecoveryPointExtendedInfo.AsBasicOperationStatusExtendedInfo() (BasicOperationStatusExtendedInfo, bool)
1. OperationStatusRecoveryPointExtendedInfo.AsOperationStatusExtendedInfo() (*OperationStatusExtendedInfo, bool)
1. OperationStatusRecoveryPointExtendedInfo.AsOperationStatusJobExtendedInfo() (*OperationStatusJobExtendedInfo, bool)
1. OperationStatusRecoveryPointExtendedInfo.AsOperationStatusJobsExtendedInfo() (*OperationStatusJobsExtendedInfo, bool)
1. OperationStatusRecoveryPointExtendedInfo.AsOperationStatusProvisionILRExtendedInfo() (*OperationStatusProvisionILRExtendedInfo, bool)
1. OperationStatusRecoveryPointExtendedInfo.AsOperationStatusRecoveryPointExtendedInfo() (*OperationStatusRecoveryPointExtendedInfo, bool)
1. OperationStatusRecoveryPointExtendedInfo.MarshalJSON() ([]byte, error)
1. OperationsClient.List(context.Context) (ClientDiscoveryResponsePage, error)
1. OperationsClient.ListComplete(context.Context) (ClientDiscoveryResponseIterator, error)
1. OperationsClient.ListPreparer(context.Context) (*http.Request, error)
1. OperationsClient.ListResponder(*http.Response) (ClientDiscoveryResponse, error)
1. OperationsClient.ListSender(*http.Request) (*http.Response, error)
1. PossibleFabricNameValues() []FabricName
1. PossibleFeatureTypeValues() []FeatureType
1. PossibleIntentItemTypeValues() []IntentItemType
1. PossibleObjectTypeBasicCrrAccessTokenValues() []ObjectTypeBasicCrrAccessToken
1. PossibleProtectionIntentItemTypeValues() []ProtectionIntentItemType
1. PossibleSupportStatusValues() []SupportStatus
1. PossibleTypeEnumValues() []TypeEnum
1. PossibleUsagesUnitValues() []UsagesUnit
1. PossibleValidationStatusValues() []ValidationStatus
1. ProtectedItemsCrrClient.List(context.Context, string, string, string, string) (ProtectedItemResourceListPage, error)
1. ProtectedItemsCrrClient.ListComplete(context.Context, string, string, string, string) (ProtectedItemResourceListIterator, error)
1. ProtectedItemsCrrClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectedItemsCrrClient.ListResponder(*http.Response) (ProtectedItemResourceList, error)
1. ProtectedItemsCrrClient.ListSender(*http.Request) (*http.Response, error)
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
1. RecoveryPointsClient.GetAccessToken(context.Context, string, string, string, string, string, string, AADPropertiesResource) (CrrAccessTokenResource, error)
1. RecoveryPointsClient.GetAccessTokenPreparer(context.Context, string, string, string, string, string, string, AADPropertiesResource) (*http.Request, error)
1. RecoveryPointsClient.GetAccessTokenResponder(*http.Response) (CrrAccessTokenResource, error)
1. RecoveryPointsClient.GetAccessTokenSender(*http.Request) (*http.Response, error)
1. RecoveryPointsCrrClient.List(context.Context, string, string, string, string, string, string) (RecoveryPointResourceListPage, error)
1. RecoveryPointsCrrClient.ListComplete(context.Context, string, string, string, string, string, string) (RecoveryPointResourceListIterator, error)
1. RecoveryPointsCrrClient.ListPreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. RecoveryPointsCrrClient.ListResponder(*http.Response) (RecoveryPointResourceList, error)
1. RecoveryPointsCrrClient.ListSender(*http.Request) (*http.Response, error)
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
1. StatusClient.Get(context.Context, string, StatusRequest) (StatusResponse, error)
1. StatusClient.GetPreparer(context.Context, string, StatusRequest) (*http.Request, error)
1. StatusClient.GetResponder(*http.Response) (StatusResponse, error)
1. StatusClient.GetSender(*http.Request) (*http.Response, error)
1. UsageSummariesCRRClient.List(context.Context, string, string, string, string) (ManagementUsageList, error)
1. UsageSummariesCRRClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. UsageSummariesCRRClient.ListResponder(*http.Response) (ManagementUsageList, error)
1. UsageSummariesCRRClient.ListSender(*http.Request) (*http.Response, error)
1. UsageSummariesClient.List(context.Context, string, string, string, string) (ManagementUsageList, error)
1. UsageSummariesClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. UsageSummariesClient.ListResponder(*http.Response) (ManagementUsageList, error)
1. UsageSummariesClient.ListSender(*http.Request) (*http.Response, error)
1. WorkloadCrrAccessToken.AsBasicCrrAccessToken() (BasicCrrAccessToken, bool)
1. WorkloadCrrAccessToken.AsCrrAccessToken() (*CrrAccessToken, bool)
1. WorkloadCrrAccessToken.AsWorkloadCrrAccessToken() (*WorkloadCrrAccessToken, bool)
1. WorkloadCrrAccessToken.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. AADProperties
1. AADPropertiesResource
1. AadPropertiesClient
1. AzureBackupGoalFeatureSupportRequest
1. AzureRecoveryServiceVaultProtectionIntent
1. AzureResourceProtectionIntent
1. AzureVMResourceFeatureSupportRequest
1. AzureVMResourceFeatureSupportResponse
1. AzureWorkloadAutoProtectionIntent
1. AzureWorkloadSQLAutoProtectionIntent
1. BMSAADPropertiesQueryObject
1. BMSBackupSummariesQueryObject
1. ClientDiscoveryDisplay
1. ClientDiscoveryForLogSpecification
1. ClientDiscoveryForProperties
1. ClientDiscoveryForServiceSpecification
1. ClientDiscoveryResponse
1. ClientDiscoveryResponseIterator
1. ClientDiscoveryResponsePage
1. ClientDiscoveryValueForSingleAPI
1. CrossRegionRestoreClient
1. CrossRegionRestoreRequest
1. CrossRegionRestoreRequestResource
1. CrossRegionRestoreTriggerFuture
1. CrrAccessToken
1. CrrAccessTokenResource
1. CrrJobDetailsClient
1. CrrJobRequest
1. CrrJobRequestResource
1. CrrJobsClient
1. CrrOperationResultsClient
1. CrrOperationStatusClient
1. FeatureSupportClient
1. FeatureSupportRequest
1. ManagementUsage
1. ManagementUsageList
1. NameInfo
1. OperationStatusRecoveryPointExtendedInfo
1. OperationsClient
1. PreValidateEnableBackupRequest
1. PreValidateEnableBackupResponse
1. ProtectedItemsCrrClient
1. ProtectionIntent
1. ProtectionIntentClient
1. ProtectionIntentGroupClient
1. ProtectionIntentQueryObject
1. ProtectionIntentResource
1. ProtectionIntentResourceList
1. ProtectionIntentResourceListIterator
1. ProtectionIntentResourceListPage
1. RecoveryPointsCrrClient
1. ResourceConfig
1. ResourceConfigResource
1. ResourceStorageConfigsClient
1. StatusClient
1. StatusRequest
1. StatusResponse
1. UsageSummariesCRRClient
1. UsageSummariesClient
1. WorkloadCrrAccessToken

#### Removed Struct Fields

1. AzureFileshareProtectedItem.HealthStatus
1. AzureWorkloadPointInTimeRestoreRequest.TargetVirtualMachineID
1. AzureWorkloadRestoreRequest.TargetVirtualMachineID
1. AzureWorkloadSAPHanaPointInTimeRestoreRequest.TargetVirtualMachineID
1. AzureWorkloadSAPHanaPointInTimeRestoreWithRehydrateRequest.TargetVirtualMachineID
1. AzureWorkloadSAPHanaRestoreRequest.TargetVirtualMachineID
1. AzureWorkloadSAPHanaRestoreWithRehydrateRequest.TargetVirtualMachineID
1. AzureWorkloadSQLPointInTimeRestoreRequest.TargetVirtualMachineID
1. AzureWorkloadSQLPointInTimeRestoreWithRehydrateRequest.TargetVirtualMachineID
1. AzureWorkloadSQLRestoreRequest.TargetVirtualMachineID
1. AzureWorkloadSQLRestoreWithRehydrateRequest.TargetVirtualMachineID

### Signature Changes

#### Struct Fields

1. AzureWorkloadBackupRequest.BackupType changed type from TypeEnum to Type

## Additive Changes

### New Constants

1. Type.TypeCopyOnlyFull
1. Type.TypeDifferential
1. Type.TypeFull
1. Type.TypeIncremental
1. Type.TypeLog
