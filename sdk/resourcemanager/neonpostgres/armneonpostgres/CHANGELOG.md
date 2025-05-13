# Release History

## 2.0.0 (2025-05-13)
### Breaking Changes

- Function `*BranchesClient.BeginUpdate` has been removed
- Function `*ComputesClient.BeginCreateOrUpdate` has been removed
- Function `*ComputesClient.Delete` has been removed
- Function `*ComputesClient.Get` has been removed
- Function `*ComputesClient.BeginUpdate` has been removed
- Function `*EndpointsClient.BeginCreateOrUpdate` has been removed
- Function `*EndpointsClient.Delete` has been removed
- Function `*EndpointsClient.Get` has been removed
- Function `*EndpointsClient.BeginUpdate` has been removed
- Function `*ProjectsClient.BeginUpdate` has been removed
- Function `*NeonDatabasesClient.BeginCreateOrUpdate` has been removed
- Function `*NeonDatabasesClient.Delete` has been removed
- Function `*NeonDatabasesClient.Get` has been removed
- Function `*NeonDatabasesClient.BeginUpdate` has been removed
- Function `*NeonRolesClient.BeginCreateOrUpdate` has been removed
- Function `*NeonRolesClient.Delete` has been removed
- Function `*NeonRolesClient.Get` has been removed
- Function `*NeonRolesClient.BeginUpdate` has been removed


## 1.0.0 (2025-04-07)
### Features Added

- New enum type `EndpointType` with values `EndpointTypeReadOnly`, `EndpointTypeReadWrite`
- New function `NewBranchesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BranchesClient, error)`
- New function `*BranchesClient.BeginCreateOrUpdate(context.Context, string, string, string, string, Branch, *BranchesClientBeginCreateOrUpdateOptions) (*runtime.Poller[BranchesClientCreateOrUpdateResponse], error)`
- New function `*BranchesClient.Delete(context.Context, string, string, string, string, *BranchesClientDeleteOptions) (BranchesClientDeleteResponse, error)`
- New function `*BranchesClient.Get(context.Context, string, string, string, string, *BranchesClientGetOptions) (BranchesClientGetResponse, error)`
- New function `*BranchesClient.NewListPager(string, string, string, *BranchesClientListOptions) *runtime.Pager[BranchesClientListResponse]`
- New function `*BranchesClient.BeginUpdate(context.Context, string, string, string, string, Branch, *BranchesClientBeginUpdateOptions) (*runtime.Poller[BranchesClientUpdateResponse], error)`
- New function `*ClientFactory.NewBranchesClient() *BranchesClient`
- New function `*ClientFactory.NewComputesClient() *ComputesClient`
- New function `*ClientFactory.NewEndpointsClient() *EndpointsClient`
- New function `*ClientFactory.NewNeonDatabasesClient() *NeonDatabasesClient`
- New function `*ClientFactory.NewNeonRolesClient() *NeonRolesClient`
- New function `*ClientFactory.NewProjectsClient() *ProjectsClient`
- New function `NewComputesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ComputesClient, error)`
- New function `*ComputesClient.BeginCreateOrUpdate(context.Context, string, string, string, string, string, Compute, *ComputesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ComputesClientCreateOrUpdateResponse], error)`
- New function `*ComputesClient.Delete(context.Context, string, string, string, string, string, *ComputesClientDeleteOptions) (ComputesClientDeleteResponse, error)`
- New function `*ComputesClient.Get(context.Context, string, string, string, string, string, *ComputesClientGetOptions) (ComputesClientGetResponse, error)`
- New function `*ComputesClient.NewListPager(string, string, string, string, *ComputesClientListOptions) *runtime.Pager[ComputesClientListResponse]`
- New function `*ComputesClient.BeginUpdate(context.Context, string, string, string, string, string, Compute, *ComputesClientBeginUpdateOptions) (*runtime.Poller[ComputesClientUpdateResponse], error)`
- New function `NewEndpointsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*EndpointsClient, error)`
- New function `*EndpointsClient.BeginCreateOrUpdate(context.Context, string, string, string, string, string, Endpoint, *EndpointsClientBeginCreateOrUpdateOptions) (*runtime.Poller[EndpointsClientCreateOrUpdateResponse], error)`
- New function `*EndpointsClient.Delete(context.Context, string, string, string, string, string, *EndpointsClientDeleteOptions) (EndpointsClientDeleteResponse, error)`
- New function `*EndpointsClient.Get(context.Context, string, string, string, string, string, *EndpointsClientGetOptions) (EndpointsClientGetResponse, error)`
- New function `*EndpointsClient.NewListPager(string, string, string, string, *EndpointsClientListOptions) *runtime.Pager[EndpointsClientListResponse]`
- New function `*EndpointsClient.BeginUpdate(context.Context, string, string, string, string, string, Endpoint, *EndpointsClientBeginUpdateOptions) (*runtime.Poller[EndpointsClientUpdateResponse], error)`
- New function `*OrganizationsClient.GetPostgresVersions(context.Context, string, *OrganizationsClientGetPostgresVersionsOptions) (OrganizationsClientGetPostgresVersionsResponse, error)`
- New function `NewProjectsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProjectsClient, error)`
- New function `*ProjectsClient.BeginCreateOrUpdate(context.Context, string, string, string, Project, *ProjectsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ProjectsClientCreateOrUpdateResponse], error)`
- New function `*ProjectsClient.Delete(context.Context, string, string, string, *ProjectsClientDeleteOptions) (ProjectsClientDeleteResponse, error)`
- New function `*ProjectsClient.Get(context.Context, string, string, string, *ProjectsClientGetOptions) (ProjectsClientGetResponse, error)`
- New function `*ProjectsClient.GetConnectionURI(context.Context, string, string, string, ConnectionURIProperties, *ProjectsClientGetConnectionURIOptions) (ProjectsClientGetConnectionURIResponse, error)`
- New function `*ProjectsClient.NewListPager(string, string, *ProjectsClientListOptions) *runtime.Pager[ProjectsClientListResponse]`
- New function `*ProjectsClient.BeginUpdate(context.Context, string, string, string, Project, *ProjectsClientBeginUpdateOptions) (*runtime.Poller[ProjectsClientUpdateResponse], error)`
- New function `NewNeonDatabasesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NeonDatabasesClient, error)`
- New function `*NeonDatabasesClient.BeginCreateOrUpdate(context.Context, string, string, string, string, string, NeonDatabase, *NeonDatabasesClientBeginCreateOrUpdateOptions) (*runtime.Poller[NeonDatabasesClientCreateOrUpdateResponse], error)`
- New function `*NeonDatabasesClient.Delete(context.Context, string, string, string, string, string, *NeonDatabasesClientDeleteOptions) (NeonDatabasesClientDeleteResponse, error)`
- New function `*NeonDatabasesClient.Get(context.Context, string, string, string, string, string, *NeonDatabasesClientGetOptions) (NeonDatabasesClientGetResponse, error)`
- New function `*NeonDatabasesClient.NewListPager(string, string, string, string, *NeonDatabasesClientListOptions) *runtime.Pager[NeonDatabasesClientListResponse]`
- New function `*NeonDatabasesClient.BeginUpdate(context.Context, string, string, string, string, string, NeonDatabase, *NeonDatabasesClientBeginUpdateOptions) (*runtime.Poller[NeonDatabasesClientUpdateResponse], error)`
- New function `NewNeonRolesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NeonRolesClient, error)`
- New function `*NeonRolesClient.BeginCreateOrUpdate(context.Context, string, string, string, string, string, NeonRole, *NeonRolesClientBeginCreateOrUpdateOptions) (*runtime.Poller[NeonRolesClientCreateOrUpdateResponse], error)`
- New function `*NeonRolesClient.Delete(context.Context, string, string, string, string, string, *NeonRolesClientDeleteOptions) (NeonRolesClientDeleteResponse, error)`
- New function `*NeonRolesClient.Get(context.Context, string, string, string, string, string, *NeonRolesClientGetOptions) (NeonRolesClientGetResponse, error)`
- New function `*NeonRolesClient.NewListPager(string, string, string, string, *NeonRolesClientListOptions) *runtime.Pager[NeonRolesClientListResponse]`
- New function `*NeonRolesClient.BeginUpdate(context.Context, string, string, string, string, string, NeonRole, *NeonRolesClientBeginUpdateOptions) (*runtime.Poller[NeonRolesClientUpdateResponse], error)`
- New struct `Attributes`
- New struct `Branch`
- New struct `BranchListResult`
- New struct `BranchProperties`
- New struct `Compute`
- New struct `ComputeListResult`
- New struct `ComputeProperties`
- New struct `ConnectionURIProperties`
- New struct `DefaultEndpointSettings`
- New struct `Endpoint`
- New struct `EndpointListResult`
- New struct `EndpointProperties`
- New struct `NeonDatabase`
- New struct `NeonDatabaseListResult`
- New struct `NeonDatabaseProperties`
- New struct `NeonRole`
- New struct `NeonRoleListResult`
- New struct `NeonRoleProperties`
- New struct `PgVersion`
- New struct `PgVersionsResult`
- New struct `Project`
- New struct `ProjectListResult`
- New struct `ProjectProperties`
- New field `ProjectProperties` in struct `OrganizationProperties`


## 0.1.0 (2024-11-20)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/neonpostgres/armneonpostgres` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).