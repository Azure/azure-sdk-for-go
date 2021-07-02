# Unreleased

## Breaking Changes

### Removed Funcs

1. *BackupInstancesValidateRestoreFuture.UnmarshalJSON([]byte) error
1. *BackupVaultsPatchFuture.UnmarshalJSON([]byte) error
1. BackupInstancesClient.ValidateRestore(context.Context, string, string, string, ValidateRestoreRequestObject) (BackupInstancesValidateRestoreFuture, error)
1. BackupInstancesClient.ValidateRestorePreparer(context.Context, string, string, string, ValidateRestoreRequestObject) (*http.Request, error)
1. BackupInstancesClient.ValidateRestoreResponder(*http.Response) (OperationJobExtendedInfo, error)
1. BackupInstancesClient.ValidateRestoreSender(*http.Request) (BackupInstancesValidateRestoreFuture, error)
1. BackupVaultsClient.GetResourcesInResourceGroup(context.Context, string) (BackupVaultResourceListPage, error)
1. BackupVaultsClient.GetResourcesInResourceGroupComplete(context.Context, string) (BackupVaultResourceListIterator, error)
1. BackupVaultsClient.GetResourcesInResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. BackupVaultsClient.GetResourcesInResourceGroupResponder(*http.Response) (BackupVaultResourceList, error)
1. BackupVaultsClient.GetResourcesInResourceGroupSender(*http.Request) (*http.Response, error)
1. BackupVaultsClient.GetResourcesInSubscription(context.Context) (BackupVaultResourceListPage, error)
1. BackupVaultsClient.GetResourcesInSubscriptionComplete(context.Context) (BackupVaultResourceListIterator, error)
1. BackupVaultsClient.GetResourcesInSubscriptionPreparer(context.Context) (*http.Request, error)
1. BackupVaultsClient.GetResourcesInSubscriptionResponder(*http.Response) (BackupVaultResourceList, error)
1. BackupVaultsClient.GetResourcesInSubscriptionSender(*http.Request) (*http.Response, error)
1. BackupVaultsClient.Patch(context.Context, string, string, PatchResourceRequestInput) (BackupVaultsPatchFuture, error)
1. BackupVaultsClient.PatchPreparer(context.Context, string, string, PatchResourceRequestInput) (*http.Request, error)
1. BackupVaultsClient.PatchResponder(*http.Response) (BackupVaultResource, error)
1. BackupVaultsClient.PatchSender(*http.Request) (BackupVaultsPatchFuture, error)
1. BaseClient.CheckFeatureSupport(context.Context, string, BasicFeatureValidationRequestBase) (FeatureValidationResponseBaseModel, error)
1. BaseClient.CheckFeatureSupportPreparer(context.Context, string, BasicFeatureValidationRequestBase) (*http.Request, error)
1. BaseClient.CheckFeatureSupportResponder(*http.Response) (FeatureValidationResponseBaseModel, error)
1. BaseClient.CheckFeatureSupportSender(*http.Request) (*http.Response, error)
1. BaseClient.GetOperationResultPatch(context.Context, string, string, string) (BackupVaultResource, error)
1. BaseClient.GetOperationResultPatchPreparer(context.Context, string, string, string) (*http.Request, error)
1. BaseClient.GetOperationResultPatchResponder(*http.Response) (BackupVaultResource, error)
1. BaseClient.GetOperationResultPatchSender(*http.Request) (*http.Response, error)
1. BaseClient.GetOperationStatus(context.Context, string, string) (OperationResource, error)
1. BaseClient.GetOperationStatusPreparer(context.Context, string, string) (*http.Request, error)
1. BaseClient.GetOperationStatusResponder(*http.Response) (OperationResource, error)
1. BaseClient.GetOperationStatusSender(*http.Request) (*http.Response, error)
1. FindRestorableTimeRangesClient.Post(context.Context, string, string, string, AzureBackupFindRestorableTimeRangesRequest) (AzureBackupFindRestorableTimeRangesResponseResource, error)
1. FindRestorableTimeRangesClient.PostPreparer(context.Context, string, string, string, AzureBackupFindRestorableTimeRangesRequest) (*http.Request, error)
1. FindRestorableTimeRangesClient.PostResponder(*http.Response) (AzureBackupFindRestorableTimeRangesResponseResource, error)
1. FindRestorableTimeRangesClient.PostSender(*http.Request) (*http.Response, error)
1. JobClient.Get(context.Context, string, string, string) (AzureBackupJobResource, error)
1. JobClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. JobClient.GetResponder(*http.Response) (AzureBackupJobResource, error)
1. JobClient.GetSender(*http.Request) (*http.Response, error)
1. NewFindRestorableTimeRangesClient(string) FindRestorableTimeRangesClient
1. NewFindRestorableTimeRangesClientWithBaseURI(string, string) FindRestorableTimeRangesClient
1. NewJobClient(string) JobClient
1. NewJobClientWithBaseURI(string, string) JobClient
1. NewRecoveryPointClient(string) RecoveryPointClient
1. NewRecoveryPointClientWithBaseURI(string, string) RecoveryPointClient
1. RecoveryPointClient.Get(context.Context, string, string, string, string) (AzureBackupRecoveryPointResource, error)
1. RecoveryPointClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. RecoveryPointClient.GetResponder(*http.Response) (AzureBackupRecoveryPointResource, error)
1. RecoveryPointClient.GetSender(*http.Request) (*http.Response, error)
1. RecoveryPointsClient.GetList(context.Context, string, string, string, string, string) (AzureBackupRecoveryPointResourceListPage, error)
1. RecoveryPointsClient.GetListComplete(context.Context, string, string, string, string, string) (AzureBackupRecoveryPointResourceListIterator, error)
1. RecoveryPointsClient.GetListPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. RecoveryPointsClient.GetListResponder(*http.Response) (AzureBackupRecoveryPointResourceList, error)
1. RecoveryPointsClient.GetListSender(*http.Request) (*http.Response, error)

### Struct Changes

#### Removed Structs

1. BackupInstancesValidateRestoreFuture
1. BackupVaultsPatchFuture
1. FindRestorableTimeRangesClient
1. JobClient
1. RecoveryPointClient

## Additive Changes

### New Funcs

1. *BackupInstancesValidateForRestoreFuture.UnmarshalJSON([]byte) error
1. *BackupVaultsUpdateFuture.UnmarshalJSON([]byte) error
1. BackupInstancesClient.ValidateForRestore(context.Context, string, string, string, ValidateRestoreRequestObject) (BackupInstancesValidateForRestoreFuture, error)
1. BackupInstancesClient.ValidateForRestorePreparer(context.Context, string, string, string, ValidateRestoreRequestObject) (*http.Request, error)
1. BackupInstancesClient.ValidateForRestoreResponder(*http.Response) (OperationJobExtendedInfo, error)
1. BackupInstancesClient.ValidateForRestoreSender(*http.Request) (BackupInstancesValidateForRestoreFuture, error)
1. BackupVaultOperationResultsClient.Get(context.Context, string, string, string) (BackupVaultResource, error)
1. BackupVaultOperationResultsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. BackupVaultOperationResultsClient.GetResponder(*http.Response) (BackupVaultResource, error)
1. BackupVaultOperationResultsClient.GetSender(*http.Request) (*http.Response, error)
1. BackupVaultsClient.GetInResourceGroup(context.Context, string) (BackupVaultResourceListPage, error)
1. BackupVaultsClient.GetInResourceGroupComplete(context.Context, string) (BackupVaultResourceListIterator, error)
1. BackupVaultsClient.GetInResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. BackupVaultsClient.GetInResourceGroupResponder(*http.Response) (BackupVaultResourceList, error)
1. BackupVaultsClient.GetInResourceGroupSender(*http.Request) (*http.Response, error)
1. BackupVaultsClient.GetInSubscription(context.Context) (BackupVaultResourceListPage, error)
1. BackupVaultsClient.GetInSubscriptionComplete(context.Context) (BackupVaultResourceListIterator, error)
1. BackupVaultsClient.GetInSubscriptionPreparer(context.Context) (*http.Request, error)
1. BackupVaultsClient.GetInSubscriptionResponder(*http.Response) (BackupVaultResourceList, error)
1. BackupVaultsClient.GetInSubscriptionSender(*http.Request) (*http.Response, error)
1. BackupVaultsClient.Update(context.Context, string, string, PatchResourceRequestInput) (BackupVaultsUpdateFuture, error)
1. BackupVaultsClient.UpdatePreparer(context.Context, string, string, PatchResourceRequestInput) (*http.Request, error)
1. BackupVaultsClient.UpdateResponder(*http.Response) (BackupVaultResource, error)
1. BackupVaultsClient.UpdateSender(*http.Request) (BackupVaultsUpdateFuture, error)
1. Client.CheckFeatureSupport(context.Context, string, BasicFeatureValidationRequestBase) (FeatureValidationResponseBaseModel, error)
1. Client.CheckFeatureSupportPreparer(context.Context, string, BasicFeatureValidationRequestBase) (*http.Request, error)
1. Client.CheckFeatureSupportResponder(*http.Response) (FeatureValidationResponseBaseModel, error)
1. Client.CheckFeatureSupportSender(*http.Request) (*http.Response, error)
1. JobsClient.Get(context.Context, string, string, string) (AzureBackupJobResource, error)
1. JobsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. JobsClient.GetResponder(*http.Response) (AzureBackupJobResource, error)
1. JobsClient.GetSender(*http.Request) (*http.Response, error)
1. NewBackupVaultOperationResultsClient(string) BackupVaultOperationResultsClient
1. NewBackupVaultOperationResultsClientWithBaseURI(string, string) BackupVaultOperationResultsClient
1. NewClient(string) Client
1. NewClientWithBaseURI(string, string) Client
1. NewOperationStatusClient(string) OperationStatusClient
1. NewOperationStatusClientWithBaseURI(string, string) OperationStatusClient
1. NewRestorableTimeRangesClient(string) RestorableTimeRangesClient
1. NewRestorableTimeRangesClientWithBaseURI(string, string) RestorableTimeRangesClient
1. OperationStatusClient.Get(context.Context, string, string) (OperationResource, error)
1. OperationStatusClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. OperationStatusClient.GetResponder(*http.Response) (OperationResource, error)
1. OperationStatusClient.GetSender(*http.Request) (*http.Response, error)
1. RecoveryPointsClient.Get(context.Context, string, string, string, string) (AzureBackupRecoveryPointResource, error)
1. RecoveryPointsClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. RecoveryPointsClient.GetResponder(*http.Response) (AzureBackupRecoveryPointResource, error)
1. RecoveryPointsClient.GetSender(*http.Request) (*http.Response, error)
1. RecoveryPointsClient.List(context.Context, string, string, string, string, string) (AzureBackupRecoveryPointResourceListPage, error)
1. RecoveryPointsClient.ListComplete(context.Context, string, string, string, string, string) (AzureBackupRecoveryPointResourceListIterator, error)
1. RecoveryPointsClient.ListPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. RecoveryPointsClient.ListResponder(*http.Response) (AzureBackupRecoveryPointResourceList, error)
1. RecoveryPointsClient.ListSender(*http.Request) (*http.Response, error)
1. RestorableTimeRangesClient.Find(context.Context, string, string, string, AzureBackupFindRestorableTimeRangesRequest) (AzureBackupFindRestorableTimeRangesResponseResource, error)
1. RestorableTimeRangesClient.FindPreparer(context.Context, string, string, string, AzureBackupFindRestorableTimeRangesRequest) (*http.Request, error)
1. RestorableTimeRangesClient.FindResponder(*http.Response) (AzureBackupFindRestorableTimeRangesResponseResource, error)
1. RestorableTimeRangesClient.FindSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. BackupInstancesValidateForRestoreFuture
1. BackupVaultOperationResultsClient
1. BackupVaultsUpdateFuture
1. Client
1. OperationStatusClient
1. RestorableTimeRangesClient
