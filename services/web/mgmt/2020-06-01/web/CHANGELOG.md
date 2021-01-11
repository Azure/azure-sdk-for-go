Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82/specification/web/resource-manager/readme.md tag: `package-2020-06`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *AppServiceCertificateOrdersCreateOrUpdateCertificateFuture.Result(AppServiceCertificateOrdersClient) (AppServiceCertificateResource, error)
1. *AppServiceCertificateOrdersCreateOrUpdateFuture.Result(AppServiceCertificateOrdersClient) (AppServiceCertificateOrder, error)
1. *AppServiceEnvironmentsChangeVnetAllFuture.Result(AppServiceEnvironmentsClient) (AppCollectionPage, error)
1. *AppServiceEnvironmentsChangeVnetFuture.Result(AppServiceEnvironmentsClient) (AppCollectionPage, error)
1. *AppServiceEnvironmentsCreateOrUpdateFuture.Result(AppServiceEnvironmentsClient) (AppServiceEnvironmentResource, error)
1. *AppServiceEnvironmentsCreateOrUpdateMultiRolePoolFuture.Result(AppServiceEnvironmentsClient) (WorkerPoolResource, error)
1. *AppServiceEnvironmentsCreateOrUpdateWorkerPoolFuture.Result(AppServiceEnvironmentsClient) (WorkerPoolResource, error)
1. *AppServiceEnvironmentsDeleteFuture.Result(AppServiceEnvironmentsClient) (autorest.Response, error)
1. *AppServiceEnvironmentsResumeAllFuture.Result(AppServiceEnvironmentsClient) (AppCollectionPage, error)
1. *AppServiceEnvironmentsResumeFuture.Result(AppServiceEnvironmentsClient) (AppCollectionPage, error)
1. *AppServiceEnvironmentsSuspendAllFuture.Result(AppServiceEnvironmentsClient) (AppCollectionPage, error)
1. *AppServiceEnvironmentsSuspendFuture.Result(AppServiceEnvironmentsClient) (AppCollectionPage, error)
1. *AppServicePlansCreateOrUpdateFuture.Result(AppServicePlansClient) (AppServicePlan, error)
1. *AppsApproveOrRejectPrivateEndpointConnectionFuture.Result(AppsClient) (PrivateEndpointConnectionResource, error)
1. *AppsCopyProductionSlotFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsCopySlotSlotFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsCreateFunctionFuture.Result(AppsClient) (FunctionEnvelope, error)
1. *AppsCreateInstanceFunctionSlotFuture.Result(AppsClient) (FunctionEnvelope, error)
1. *AppsCreateInstanceMSDeployOperationFuture.Result(AppsClient) (MSDeployStatus, error)
1. *AppsCreateInstanceMSDeployOperationSlotFuture.Result(AppsClient) (MSDeployStatus, error)
1. *AppsCreateMSDeployOperationFuture.Result(AppsClient) (MSDeployStatus, error)
1. *AppsCreateMSDeployOperationSlotFuture.Result(AppsClient) (MSDeployStatus, error)
1. *AppsCreateOrUpdateFuture.Result(AppsClient) (Site, error)
1. *AppsCreateOrUpdateSlotFuture.Result(AppsClient) (Site, error)
1. *AppsCreateOrUpdateSourceControlFuture.Result(AppsClient) (SiteSourceControl, error)
1. *AppsCreateOrUpdateSourceControlSlotFuture.Result(AppsClient) (SiteSourceControl, error)
1. *AppsDeletePrivateEndpointConnectionFuture.Result(AppsClient) (SetObject, error)
1. *AppsInstallSiteExtensionFuture.Result(AppsClient) (SiteExtensionInfo, error)
1. *AppsInstallSiteExtensionSlotFuture.Result(AppsClient) (SiteExtensionInfo, error)
1. *AppsListPublishingCredentialsFuture.Result(AppsClient) (User, error)
1. *AppsListPublishingCredentialsSlotFuture.Result(AppsClient) (User, error)
1. *AppsMigrateMySQLFuture.Result(AppsClient) (Operation, error)
1. *AppsMigrateStorageFuture.Result(AppsClient) (StorageMigrationResponse, error)
1. *AppsRestoreFromBackupBlobFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsRestoreFromBackupBlobSlotFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsRestoreFromDeletedAppFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsRestoreFromDeletedAppSlotFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsRestoreFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsRestoreSlotFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsRestoreSnapshotFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsRestoreSnapshotSlotFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsStartNetworkTraceFuture.Result(AppsClient) (ListNetworkTrace, error)
1. *AppsStartNetworkTraceSlotFuture.Result(AppsClient) (ListNetworkTrace, error)
1. *AppsStartWebSiteNetworkTraceOperationFuture.Result(AppsClient) (ListNetworkTrace, error)
1. *AppsStartWebSiteNetworkTraceOperationSlotFuture.Result(AppsClient) (ListNetworkTrace, error)
1. *AppsSwapSlotSlotFuture.Result(AppsClient) (autorest.Response, error)
1. *AppsSwapSlotWithProductionFuture.Result(AppsClient) (autorest.Response, error)
1. *DomainsCreateOrUpdateFuture.Result(DomainsClient) (Domain, error)

## Struct Changes

### Removed Struct Fields

1. AppServiceCertificateOrdersCreateOrUpdateCertificateFuture.azure.Future
1. AppServiceCertificateOrdersCreateOrUpdateFuture.azure.Future
1. AppServiceEnvironmentsChangeVnetAllFuture.azure.Future
1. AppServiceEnvironmentsChangeVnetFuture.azure.Future
1. AppServiceEnvironmentsCreateOrUpdateFuture.azure.Future
1. AppServiceEnvironmentsCreateOrUpdateMultiRolePoolFuture.azure.Future
1. AppServiceEnvironmentsCreateOrUpdateWorkerPoolFuture.azure.Future
1. AppServiceEnvironmentsDeleteFuture.azure.Future
1. AppServiceEnvironmentsResumeAllFuture.azure.Future
1. AppServiceEnvironmentsResumeFuture.azure.Future
1. AppServiceEnvironmentsSuspendAllFuture.azure.Future
1. AppServiceEnvironmentsSuspendFuture.azure.Future
1. AppServicePlansCreateOrUpdateFuture.azure.Future
1. AppsApproveOrRejectPrivateEndpointConnectionFuture.azure.Future
1. AppsCopyProductionSlotFuture.azure.Future
1. AppsCopySlotSlotFuture.azure.Future
1. AppsCreateFunctionFuture.azure.Future
1. AppsCreateInstanceFunctionSlotFuture.azure.Future
1. AppsCreateInstanceMSDeployOperationFuture.azure.Future
1. AppsCreateInstanceMSDeployOperationSlotFuture.azure.Future
1. AppsCreateMSDeployOperationFuture.azure.Future
1. AppsCreateMSDeployOperationSlotFuture.azure.Future
1. AppsCreateOrUpdateFuture.azure.Future
1. AppsCreateOrUpdateSlotFuture.azure.Future
1. AppsCreateOrUpdateSourceControlFuture.azure.Future
1. AppsCreateOrUpdateSourceControlSlotFuture.azure.Future
1. AppsDeletePrivateEndpointConnectionFuture.azure.Future
1. AppsInstallSiteExtensionFuture.azure.Future
1. AppsInstallSiteExtensionSlotFuture.azure.Future
1. AppsListPublishingCredentialsFuture.azure.Future
1. AppsListPublishingCredentialsSlotFuture.azure.Future
1. AppsMigrateMySQLFuture.azure.Future
1. AppsMigrateStorageFuture.azure.Future
1. AppsRestoreFromBackupBlobFuture.azure.Future
1. AppsRestoreFromBackupBlobSlotFuture.azure.Future
1. AppsRestoreFromDeletedAppFuture.azure.Future
1. AppsRestoreFromDeletedAppSlotFuture.azure.Future
1. AppsRestoreFuture.azure.Future
1. AppsRestoreSlotFuture.azure.Future
1. AppsRestoreSnapshotFuture.azure.Future
1. AppsRestoreSnapshotSlotFuture.azure.Future
1. AppsStartNetworkTraceFuture.azure.Future
1. AppsStartNetworkTraceSlotFuture.azure.Future
1. AppsStartWebSiteNetworkTraceOperationFuture.azure.Future
1. AppsStartWebSiteNetworkTraceOperationSlotFuture.azure.Future
1. AppsSwapSlotSlotFuture.azure.Future
1. AppsSwapSlotWithProductionFuture.azure.Future
1. DomainsCreateOrUpdateFuture.azure.Future

## Struct Changes

### New Struct Fields

1. AppServiceCertificateOrdersCreateOrUpdateCertificateFuture.Result
1. AppServiceCertificateOrdersCreateOrUpdateCertificateFuture.azure.FutureAPI
1. AppServiceCertificateOrdersCreateOrUpdateFuture.Result
1. AppServiceCertificateOrdersCreateOrUpdateFuture.azure.FutureAPI
1. AppServiceEnvironmentsChangeVnetAllFuture.Result
1. AppServiceEnvironmentsChangeVnetAllFuture.azure.FutureAPI
1. AppServiceEnvironmentsChangeVnetFuture.Result
1. AppServiceEnvironmentsChangeVnetFuture.azure.FutureAPI
1. AppServiceEnvironmentsCreateOrUpdateFuture.Result
1. AppServiceEnvironmentsCreateOrUpdateFuture.azure.FutureAPI
1. AppServiceEnvironmentsCreateOrUpdateMultiRolePoolFuture.Result
1. AppServiceEnvironmentsCreateOrUpdateMultiRolePoolFuture.azure.FutureAPI
1. AppServiceEnvironmentsCreateOrUpdateWorkerPoolFuture.Result
1. AppServiceEnvironmentsCreateOrUpdateWorkerPoolFuture.azure.FutureAPI
1. AppServiceEnvironmentsDeleteFuture.Result
1. AppServiceEnvironmentsDeleteFuture.azure.FutureAPI
1. AppServiceEnvironmentsResumeAllFuture.Result
1. AppServiceEnvironmentsResumeAllFuture.azure.FutureAPI
1. AppServiceEnvironmentsResumeFuture.Result
1. AppServiceEnvironmentsResumeFuture.azure.FutureAPI
1. AppServiceEnvironmentsSuspendAllFuture.Result
1. AppServiceEnvironmentsSuspendAllFuture.azure.FutureAPI
1. AppServiceEnvironmentsSuspendFuture.Result
1. AppServiceEnvironmentsSuspendFuture.azure.FutureAPI
1. AppServicePlansCreateOrUpdateFuture.Result
1. AppServicePlansCreateOrUpdateFuture.azure.FutureAPI
1. AppsApproveOrRejectPrivateEndpointConnectionFuture.Result
1. AppsApproveOrRejectPrivateEndpointConnectionFuture.azure.FutureAPI
1. AppsCopyProductionSlotFuture.Result
1. AppsCopyProductionSlotFuture.azure.FutureAPI
1. AppsCopySlotSlotFuture.Result
1. AppsCopySlotSlotFuture.azure.FutureAPI
1. AppsCreateFunctionFuture.Result
1. AppsCreateFunctionFuture.azure.FutureAPI
1. AppsCreateInstanceFunctionSlotFuture.Result
1. AppsCreateInstanceFunctionSlotFuture.azure.FutureAPI
1. AppsCreateInstanceMSDeployOperationFuture.Result
1. AppsCreateInstanceMSDeployOperationFuture.azure.FutureAPI
1. AppsCreateInstanceMSDeployOperationSlotFuture.Result
1. AppsCreateInstanceMSDeployOperationSlotFuture.azure.FutureAPI
1. AppsCreateMSDeployOperationFuture.Result
1. AppsCreateMSDeployOperationFuture.azure.FutureAPI
1. AppsCreateMSDeployOperationSlotFuture.Result
1. AppsCreateMSDeployOperationSlotFuture.azure.FutureAPI
1. AppsCreateOrUpdateFuture.Result
1. AppsCreateOrUpdateFuture.azure.FutureAPI
1. AppsCreateOrUpdateSlotFuture.Result
1. AppsCreateOrUpdateSlotFuture.azure.FutureAPI
1. AppsCreateOrUpdateSourceControlFuture.Result
1. AppsCreateOrUpdateSourceControlFuture.azure.FutureAPI
1. AppsCreateOrUpdateSourceControlSlotFuture.Result
1. AppsCreateOrUpdateSourceControlSlotFuture.azure.FutureAPI
1. AppsDeletePrivateEndpointConnectionFuture.Result
1. AppsDeletePrivateEndpointConnectionFuture.azure.FutureAPI
1. AppsInstallSiteExtensionFuture.Result
1. AppsInstallSiteExtensionFuture.azure.FutureAPI
1. AppsInstallSiteExtensionSlotFuture.Result
1. AppsInstallSiteExtensionSlotFuture.azure.FutureAPI
1. AppsListPublishingCredentialsFuture.Result
1. AppsListPublishingCredentialsFuture.azure.FutureAPI
1. AppsListPublishingCredentialsSlotFuture.Result
1. AppsListPublishingCredentialsSlotFuture.azure.FutureAPI
1. AppsMigrateMySQLFuture.Result
1. AppsMigrateMySQLFuture.azure.FutureAPI
1. AppsMigrateStorageFuture.Result
1. AppsMigrateStorageFuture.azure.FutureAPI
1. AppsRestoreFromBackupBlobFuture.Result
1. AppsRestoreFromBackupBlobFuture.azure.FutureAPI
1. AppsRestoreFromBackupBlobSlotFuture.Result
1. AppsRestoreFromBackupBlobSlotFuture.azure.FutureAPI
1. AppsRestoreFromDeletedAppFuture.Result
1. AppsRestoreFromDeletedAppFuture.azure.FutureAPI
1. AppsRestoreFromDeletedAppSlotFuture.Result
1. AppsRestoreFromDeletedAppSlotFuture.azure.FutureAPI
1. AppsRestoreFuture.Result
1. AppsRestoreFuture.azure.FutureAPI
1. AppsRestoreSlotFuture.Result
1. AppsRestoreSlotFuture.azure.FutureAPI
1. AppsRestoreSnapshotFuture.Result
1. AppsRestoreSnapshotFuture.azure.FutureAPI
1. AppsRestoreSnapshotSlotFuture.Result
1. AppsRestoreSnapshotSlotFuture.azure.FutureAPI
1. AppsStartNetworkTraceFuture.Result
1. AppsStartNetworkTraceFuture.azure.FutureAPI
1. AppsStartNetworkTraceSlotFuture.Result
1. AppsStartNetworkTraceSlotFuture.azure.FutureAPI
1. AppsStartWebSiteNetworkTraceOperationFuture.Result
1. AppsStartWebSiteNetworkTraceOperationFuture.azure.FutureAPI
1. AppsStartWebSiteNetworkTraceOperationSlotFuture.Result
1. AppsStartWebSiteNetworkTraceOperationSlotFuture.azure.FutureAPI
1. AppsSwapSlotSlotFuture.Result
1. AppsSwapSlotSlotFuture.azure.FutureAPI
1. AppsSwapSlotWithProductionFuture.Result
1. AppsSwapSlotWithProductionFuture.azure.FutureAPI
1. DomainsCreateOrUpdateFuture.Result
1. DomainsCreateOrUpdateFuture.azure.FutureAPI
