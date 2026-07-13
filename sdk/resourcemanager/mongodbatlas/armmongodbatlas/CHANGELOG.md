# Release History

## 1.1.0-beta.1 (2026-06-20)
### Features Added

- New enum type `ClusterTier` with values `ClusterTierFLEX`, `ClusterTierFREE`, `ClusterTierM10`, `ClusterTierM30`
- New function `*ClientFactory.NewClustersClient() *ClustersClient`
- New function `*ClientFactory.NewProjectsClient() *ProjectsClient`
- New function `NewClustersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClustersClient, error)`
- New function `*ClustersClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, organizationName string, projectName string, clusterName string, resource Cluster, options *ClustersClientBeginCreateOrUpdateOptions) (*runtime.Poller[ClustersClientCreateOrUpdateResponse], error)`
- New function `*ClustersClient.BeginDelete(ctx context.Context, resourceGroupName string, organizationName string, projectName string, clusterName string, options *ClustersClientBeginDeleteOptions) (*runtime.Poller[ClustersClientDeleteResponse], error)`
- New function `*ClustersClient.Get(ctx context.Context, resourceGroupName string, organizationName string, projectName string, clusterName string, options *ClustersClientGetOptions) (ClustersClientGetResponse, error)`
- New function `*ClustersClient.NewListPager(resourceGroupName string, organizationName string, projectName string, options *ClustersClientListOptions) *runtime.Pager[ClustersClientListResponse]`
- New function `NewProjectsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ProjectsClient, error)`
- New function `*ProjectsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, organizationName string, projectName string, resource Project, options *ProjectsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ProjectsClientCreateOrUpdateResponse], error)`
- New function `*ProjectsClient.BeginDelete(ctx context.Context, resourceGroupName string, organizationName string, projectName string, options *ProjectsClientBeginDeleteOptions) (*runtime.Poller[ProjectsClientDeleteResponse], error)`
- New function `*ProjectsClient.Get(ctx context.Context, resourceGroupName string, organizationName string, projectName string, options *ProjectsClientGetOptions) (ProjectsClientGetResponse, error)`
- New function `*ProjectsClient.ListClusterTierRegions(ctx context.Context, resourceGroupName string, organizationName string, projectName string, options *ProjectsClientListClusterTierRegionsOptions) (ProjectsClientListClusterTierRegionsResponse, error)`
- New function `*ProjectsClient.NewListPager(resourceGroupName string, organizationName string, options *ProjectsClientListOptions) *runtime.Pager[ProjectsClientListResponse]`
- New function `*ProjectsClient.TierLimitReached(ctx context.Context, resourceGroupName string, organizationName string, projectName string, options *ProjectsClientTierLimitReachedOptions) (ProjectsClientTierLimitReachedResponse, error)`
- New struct `Cluster`
- New struct `ClusterListResult`
- New struct `ClusterProperties`
- New struct `Project`
- New struct `ProjectLimitStatus`
- New struct `ProjectListResult`
- New struct `ProjectProperties`
- New struct `RegionsByTierResponse`
- New struct `TierLimitReachedResponse`
- New struct `TierRegions`


## 1.0.0 (2025-07-02)
### Breaking Changes

- Function `*OrganizationsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, OrganizationResource, *OrganizationsClientBeginUpdateOptions)` to `(context.Context, string, string, OrganizationResourceUpdate, *OrganizationsClientBeginUpdateOptions)`

### Features Added

- New struct `OrganizationResourceUpdate`
- New struct `OrganizationResourceUpdateProperties`


## 0.1.0 (2025-05-07)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mongodbatlas/armmongodbatlas` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).