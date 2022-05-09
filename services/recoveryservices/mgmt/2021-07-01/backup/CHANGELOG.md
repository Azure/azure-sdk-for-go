# Unreleased

## Breaking Changes

### Removed Constants

1. ObjectTypeBasicCrrAccessToken.ObjectTypeBasicCrrAccessTokenObjectTypeCrrAccessToken
1. ObjectTypeBasicCrrAccessToken.ObjectTypeBasicCrrAccessTokenObjectTypeWorkloadCrrAccessToken
1. ObjectTypeBasicOperationStatusExtendedInfo.ObjectTypeBasicOperationStatusExtendedInfoObjectTypeOperationStatusRecoveryPointExtendedInfo

### Removed Funcs

1. *CrossRegionRestoreRequest.UnmarshalJSON([]byte) error
1. *CrossRegionRestoreTriggerFuture.UnmarshalJSON([]byte) error
1. *CrrAccessTokenResource.UnmarshalJSON([]byte) error
1. *OperationStatusRecoveryPointExtendedInfo.UnmarshalJSON([]byte) error
1. AADPropertiesResource.MarshalJSON() ([]byte, error)
1. AadPropertiesClient.Get(context.Context, string, string) (AADPropertiesResource, error)
1. AadPropertiesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. AadPropertiesClient.GetResponder(*http.Response) (AADPropertiesResource, error)
1. AadPropertiesClient.GetSender(*http.Request) (*http.Response, error)
1. AzureWorkloadSQLRecoveryPointExtendedInfo.MarshalJSON() ([]byte, error)
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
1. NewAadPropertiesClient(string) AadPropertiesClient
1. NewAadPropertiesClientWithBaseURI(string, string) AadPropertiesClient
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
1. NewProtectedItemsCrrClient(string) ProtectedItemsCrrClient
1. NewProtectedItemsCrrClientWithBaseURI(string, string) ProtectedItemsCrrClient
1. NewRecoveryPointsCrrClient(string) RecoveryPointsCrrClient
1. NewRecoveryPointsCrrClientWithBaseURI(string, string) RecoveryPointsCrrClient
1. NewResourceStorageConfigsClient(string) ResourceStorageConfigsClient
1. NewResourceStorageConfigsClientWithBaseURI(string, string) ResourceStorageConfigsClient
1. NewUsageSummariesCRRClient(string) UsageSummariesCRRClient
1. NewUsageSummariesCRRClientWithBaseURI(string, string) UsageSummariesCRRClient
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
1. PossibleObjectTypeBasicCrrAccessTokenValues() []ObjectTypeBasicCrrAccessToken
1. ProtectedItemsCrrClient.List(context.Context, string, string, string, string) (ProtectedItemResourceListPage, error)
1. ProtectedItemsCrrClient.ListComplete(context.Context, string, string, string, string) (ProtectedItemResourceListIterator, error)
1. ProtectedItemsCrrClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ProtectedItemsCrrClient.ListResponder(*http.Response) (ProtectedItemResourceList, error)
1. ProtectedItemsCrrClient.ListSender(*http.Request) (*http.Response, error)
1. RecoveryPointsClient.GetAccessToken(context.Context, string, string, string, string, string, string, AADPropertiesResource) (CrrAccessTokenResource, error)
1. RecoveryPointsClient.GetAccessTokenPreparer(context.Context, string, string, string, string, string, string, AADPropertiesResource) (*http.Request, error)
1. RecoveryPointsClient.GetAccessTokenResponder(*http.Response) (CrrAccessTokenResource, error)
1. RecoveryPointsClient.GetAccessTokenSender(*http.Request) (*http.Response, error)
1. RecoveryPointsCrrClient.List(context.Context, string, string, string, string, string, string) (RecoveryPointResourceListPage, error)
1. RecoveryPointsCrrClient.ListComplete(context.Context, string, string, string, string, string, string) (RecoveryPointResourceListIterator, error)
1. RecoveryPointsCrrClient.ListPreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. RecoveryPointsCrrClient.ListResponder(*http.Response) (RecoveryPointResourceList, error)
1. RecoveryPointsCrrClient.ListSender(*http.Request) (*http.Response, error)
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
1. UsageSummariesCRRClient.List(context.Context, string, string, string, string) (ManagementUsageList, error)
1. UsageSummariesCRRClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. UsageSummariesCRRClient.ListResponder(*http.Response) (ManagementUsageList, error)
1. UsageSummariesCRRClient.ListSender(*http.Request) (*http.Response, error)
1. WorkloadCrrAccessToken.AsBasicCrrAccessToken() (BasicCrrAccessToken, bool)
1. WorkloadCrrAccessToken.AsCrrAccessToken() (*CrrAccessToken, bool)
1. WorkloadCrrAccessToken.AsWorkloadCrrAccessToken() (*WorkloadCrrAccessToken, bool)
1. WorkloadCrrAccessToken.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. AADProperties
1. AADPropertiesResource
1. AadPropertiesClient
1. BMSAADPropertiesQueryObject
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
1. OperationStatusRecoveryPointExtendedInfo
1. ProtectedItemsCrrClient
1. RecoveryPointsCrrClient
1. ResourceStorageConfigsClient
1. UsageSummariesCRRClient
1. WorkloadCrrAccessToken

#### Removed Struct Fields

1. AzureFileshareProtectedItem.HealthStatus
