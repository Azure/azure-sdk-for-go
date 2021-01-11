Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82/specification/web/resource-manager/readme.md tag: `package-2015-08-preview`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *HostingEnvironmentsCreateOrUpdateHostingEnvironmentFuture.Result(HostingEnvironmentsClient) (HostingEnvironment, error)
1. *HostingEnvironmentsCreateOrUpdateMultiRolePoolFuture.Result(HostingEnvironmentsClient) (WorkerPool, error)
1. *HostingEnvironmentsCreateOrUpdateWorkerPoolFuture.Result(HostingEnvironmentsClient) (WorkerPool, error)
1. *HostingEnvironmentsDeleteHostingEnvironmentFuture.Result(HostingEnvironmentsClient) (SetObject, error)
1. *HostingEnvironmentsResumeHostingEnvironmentAllFuture.Result(HostingEnvironmentsClient) (SiteCollectionPage, error)
1. *HostingEnvironmentsResumeHostingEnvironmentFuture.Result(HostingEnvironmentsClient) (SiteCollectionPage, error)
1. *HostingEnvironmentsSuspendHostingEnvironmentAllFuture.Result(HostingEnvironmentsClient) (SiteCollectionPage, error)
1. *HostingEnvironmentsSuspendHostingEnvironmentFuture.Result(HostingEnvironmentsClient) (SiteCollectionPage, error)
1. *ManagedHostingEnvironmentsCreateOrUpdateManagedHostingEnvironmentFuture.Result(ManagedHostingEnvironmentsClient) (HostingEnvironment, error)
1. *ManagedHostingEnvironmentsDeleteManagedHostingEnvironmentFuture.Result(ManagedHostingEnvironmentsClient) (SetObject, error)
1. *ServerFarmsCreateOrUpdateServerFarmFuture.Result(ServerFarmsClient) (ServerFarmWithRichSku, error)
1. *SitesCreateOrUpdateSiteFuture.Result(SitesClient) (Site, error)
1. *SitesCreateOrUpdateSiteSlotFuture.Result(SitesClient) (Site, error)
1. *SitesListSitePublishingCredentialsFuture.Result(SitesClient) (User, error)
1. *SitesListSitePublishingCredentialsSlotFuture.Result(SitesClient) (User, error)
1. *SitesRecoverSiteFuture.Result(SitesClient) (Site, error)
1. *SitesRecoverSiteSlotFuture.Result(SitesClient) (Site, error)
1. *SitesRestoreSiteFuture.Result(SitesClient) (RestoreResponse, error)
1. *SitesRestoreSiteSlotFuture.Result(SitesClient) (RestoreResponse, error)
1. *SitesSwapSlotWithProductionFuture.Result(SitesClient) (SetObject, error)
1. *SitesSwapSlotsSlotFuture.Result(SitesClient) (SetObject, error)

## Struct Changes

### Removed Struct Fields

1. HostingEnvironmentsCreateOrUpdateHostingEnvironmentFuture.azure.Future
1. HostingEnvironmentsCreateOrUpdateMultiRolePoolFuture.azure.Future
1. HostingEnvironmentsCreateOrUpdateWorkerPoolFuture.azure.Future
1. HostingEnvironmentsDeleteHostingEnvironmentFuture.azure.Future
1. HostingEnvironmentsResumeHostingEnvironmentAllFuture.azure.Future
1. HostingEnvironmentsResumeHostingEnvironmentFuture.azure.Future
1. HostingEnvironmentsSuspendHostingEnvironmentAllFuture.azure.Future
1. HostingEnvironmentsSuspendHostingEnvironmentFuture.azure.Future
1. ManagedHostingEnvironmentsCreateOrUpdateManagedHostingEnvironmentFuture.azure.Future
1. ManagedHostingEnvironmentsDeleteManagedHostingEnvironmentFuture.azure.Future
1. ServerFarmsCreateOrUpdateServerFarmFuture.azure.Future
1. SitesCreateOrUpdateSiteFuture.azure.Future
1. SitesCreateOrUpdateSiteSlotFuture.azure.Future
1. SitesListSitePublishingCredentialsFuture.azure.Future
1. SitesListSitePublishingCredentialsSlotFuture.azure.Future
1. SitesRecoverSiteFuture.azure.Future
1. SitesRecoverSiteSlotFuture.azure.Future
1. SitesRestoreSiteFuture.azure.Future
1. SitesRestoreSiteSlotFuture.azure.Future
1. SitesSwapSlotWithProductionFuture.azure.Future
1. SitesSwapSlotsSlotFuture.azure.Future

## Struct Changes

### New Struct Fields

1. HostingEnvironmentsCreateOrUpdateHostingEnvironmentFuture.Result
1. HostingEnvironmentsCreateOrUpdateHostingEnvironmentFuture.azure.FutureAPI
1. HostingEnvironmentsCreateOrUpdateMultiRolePoolFuture.Result
1. HostingEnvironmentsCreateOrUpdateMultiRolePoolFuture.azure.FutureAPI
1. HostingEnvironmentsCreateOrUpdateWorkerPoolFuture.Result
1. HostingEnvironmentsCreateOrUpdateWorkerPoolFuture.azure.FutureAPI
1. HostingEnvironmentsDeleteHostingEnvironmentFuture.Result
1. HostingEnvironmentsDeleteHostingEnvironmentFuture.azure.FutureAPI
1. HostingEnvironmentsResumeHostingEnvironmentAllFuture.Result
1. HostingEnvironmentsResumeHostingEnvironmentAllFuture.azure.FutureAPI
1. HostingEnvironmentsResumeHostingEnvironmentFuture.Result
1. HostingEnvironmentsResumeHostingEnvironmentFuture.azure.FutureAPI
1. HostingEnvironmentsSuspendHostingEnvironmentAllFuture.Result
1. HostingEnvironmentsSuspendHostingEnvironmentAllFuture.azure.FutureAPI
1. HostingEnvironmentsSuspendHostingEnvironmentFuture.Result
1. HostingEnvironmentsSuspendHostingEnvironmentFuture.azure.FutureAPI
1. ManagedHostingEnvironmentsCreateOrUpdateManagedHostingEnvironmentFuture.Result
1. ManagedHostingEnvironmentsCreateOrUpdateManagedHostingEnvironmentFuture.azure.FutureAPI
1. ManagedHostingEnvironmentsDeleteManagedHostingEnvironmentFuture.Result
1. ManagedHostingEnvironmentsDeleteManagedHostingEnvironmentFuture.azure.FutureAPI
1. ServerFarmsCreateOrUpdateServerFarmFuture.Result
1. ServerFarmsCreateOrUpdateServerFarmFuture.azure.FutureAPI
1. SitesCreateOrUpdateSiteFuture.Result
1. SitesCreateOrUpdateSiteFuture.azure.FutureAPI
1. SitesCreateOrUpdateSiteSlotFuture.Result
1. SitesCreateOrUpdateSiteSlotFuture.azure.FutureAPI
1. SitesListSitePublishingCredentialsFuture.Result
1. SitesListSitePublishingCredentialsFuture.azure.FutureAPI
1. SitesListSitePublishingCredentialsSlotFuture.Result
1. SitesListSitePublishingCredentialsSlotFuture.azure.FutureAPI
1. SitesRecoverSiteFuture.Result
1. SitesRecoverSiteFuture.azure.FutureAPI
1. SitesRecoverSiteSlotFuture.Result
1. SitesRecoverSiteSlotFuture.azure.FutureAPI
1. SitesRestoreSiteFuture.Result
1. SitesRestoreSiteFuture.azure.FutureAPI
1. SitesRestoreSiteSlotFuture.Result
1. SitesRestoreSiteSlotFuture.azure.FutureAPI
1. SitesSwapSlotWithProductionFuture.Result
1. SitesSwapSlotWithProductionFuture.azure.FutureAPI
1. SitesSwapSlotsSlotFuture.Result
1. SitesSwapSlotsSlotFuture.azure.FutureAPI
