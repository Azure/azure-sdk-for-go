# Release History

## 1.0.0 (2024-02-23)
### Breaking Changes

- Function `*ServicesClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, Service, *ServicesClientCreateOrUpdateOptions)` to `(context.Context, string, string, Service, *ServicesClientCreateOrUpdateOptions)`
- Function `*ServicesClient.Delete` parameter(s) have been changed from `(context.Context, string, *ServicesClientDeleteOptions)` to `(context.Context, string, string, *ServicesClientDeleteOptions)`
- Function `*ServicesClient.Get` parameter(s) have been changed from `(context.Context, string, *ServicesClientGetOptions)` to `(context.Context, string, string, *ServicesClientGetOptions)`
- Function `*ServicesClient.Update` parameter(s) have been changed from `(context.Context, string, ServiceUpdate, *ServicesClientUpdateOptions)` to `(context.Context, string, string, ServiceUpdate, *ServicesClientUpdateOptions)`
- Struct `ServiceCollection` has been removed
- Field `Properties` of struct `ServiceUpdate` has been removed
- Field `ServiceCollection` of struct `ServicesClientListByResourceGroupResponse` has been removed
- Field `ServiceCollection` of struct `ServicesClientListBySubscriptionResponse` has been removed

### Features Added

- New enum type `APIKind` with values `APIKindGraphql`, `APIKindGrpc`, `APIKindRest`, `APIKindSoap`, `APIKindWebhook`, `APIKindWebsocket`
- New enum type `APISpecExportResultFormat` with values `APISpecExportResultFormatInline`, `APISpecExportResultFormatLink`
- New enum type `APISpecImportSourceFormat` with values `APISpecImportSourceFormatInline`, `APISpecImportSourceFormatLink`
- New enum type `DeploymentState` with values `DeploymentStateActive`, `DeploymentStateInactive`
- New enum type `EnvironmentKind` with values `EnvironmentKindDevelopment`, `EnvironmentKindProduction`, `EnvironmentKindStaging`, `EnvironmentKindTesting`
- New enum type `EnvironmentServerType` with values `EnvironmentServerTypeAWSAPIGateway`, `EnvironmentServerTypeApigeeAPIManagement`, `EnvironmentServerTypeAzureAPIManagement`, `EnvironmentServerTypeAzureComputeService`, `EnvironmentServerTypeKongAPIGateway`, `EnvironmentServerTypeKubernetes`, `EnvironmentServerTypeMuleSoftAPIManagement`
- New enum type `LifecycleStage` with values `LifecycleStageDeprecated`, `LifecycleStageDesign`, `LifecycleStageDevelopment`, `LifecycleStagePreview`, `LifecycleStageProduction`, `LifecycleStageRetired`, `LifecycleStageTesting`
- New enum type `MetadataAssignmentEntity` with values `MetadataAssignmentEntityAPI`, `MetadataAssignmentEntityDeployment`, `MetadataAssignmentEntityEnvironment`
- New enum type `MetadataSchemaExportFormat` with values `MetadataSchemaExportFormatInline`, `MetadataSchemaExportFormatLink`
- New function `NewAPIDefinitionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*APIDefinitionsClient, error)`
- New function `*APIDefinitionsClient.CreateOrUpdate(context.Context, string, string, string, string, string, string, APIDefinition, *APIDefinitionsClientCreateOrUpdateOptions) (APIDefinitionsClientCreateOrUpdateResponse, error)`
- New function `*APIDefinitionsClient.Delete(context.Context, string, string, string, string, string, string, *APIDefinitionsClientDeleteOptions) (APIDefinitionsClientDeleteResponse, error)`
- New function `*APIDefinitionsClient.BeginExportSpecification(context.Context, string, string, string, string, string, string, *APIDefinitionsClientBeginExportSpecificationOptions) (*runtime.Poller[APIDefinitionsClientExportSpecificationResponse], error)`
- New function `*APIDefinitionsClient.Get(context.Context, string, string, string, string, string, string, *APIDefinitionsClientGetOptions) (APIDefinitionsClientGetResponse, error)`
- New function `*APIDefinitionsClient.Head(context.Context, string, string, string, string, string, string, *APIDefinitionsClientHeadOptions) (APIDefinitionsClientHeadResponse, error)`
- New function `*APIDefinitionsClient.BeginImportSpecification(context.Context, string, string, string, string, string, string, APISpecImportRequest, *APIDefinitionsClientBeginImportSpecificationOptions) (*runtime.Poller[APIDefinitionsClientImportSpecificationResponse], error)`
- New function `*APIDefinitionsClient.NewListPager(string, string, string, string, string, *APIDefinitionsClientListOptions) *runtime.Pager[APIDefinitionsClientListResponse]`
- New function `NewAPIVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*APIVersionsClient, error)`
- New function `*APIVersionsClient.CreateOrUpdate(context.Context, string, string, string, string, string, APIVersion, *APIVersionsClientCreateOrUpdateOptions) (APIVersionsClientCreateOrUpdateResponse, error)`
- New function `*APIVersionsClient.Delete(context.Context, string, string, string, string, string, *APIVersionsClientDeleteOptions) (APIVersionsClientDeleteResponse, error)`
- New function `*APIVersionsClient.Get(context.Context, string, string, string, string, string, *APIVersionsClientGetOptions) (APIVersionsClientGetResponse, error)`
- New function `*APIVersionsClient.Head(context.Context, string, string, string, string, string, *APIVersionsClientHeadOptions) (APIVersionsClientHeadResponse, error)`
- New function `*APIVersionsClient.NewListPager(string, string, string, string, *APIVersionsClientListOptions) *runtime.Pager[APIVersionsClientListResponse]`
- New function `NewApisClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ApisClient, error)`
- New function `*ApisClient.CreateOrUpdate(context.Context, string, string, string, string, API, *ApisClientCreateOrUpdateOptions) (ApisClientCreateOrUpdateResponse, error)`
- New function `*ApisClient.Delete(context.Context, string, string, string, string, *ApisClientDeleteOptions) (ApisClientDeleteResponse, error)`
- New function `*ApisClient.Get(context.Context, string, string, string, string, *ApisClientGetOptions) (ApisClientGetResponse, error)`
- New function `*ApisClient.Head(context.Context, string, string, string, string, *ApisClientHeadOptions) (ApisClientHeadResponse, error)`
- New function `*ApisClient.NewListPager(string, string, string, *ApisClientListOptions) *runtime.Pager[ApisClientListResponse]`
- New function `*ClientFactory.NewAPIDefinitionsClient() *APIDefinitionsClient`
- New function `*ClientFactory.NewAPIVersionsClient() *APIVersionsClient`
- New function `*ClientFactory.NewApisClient() *ApisClient`
- New function `*ClientFactory.NewDeploymentsClient() *DeploymentsClient`
- New function `*ClientFactory.NewEnvironmentsClient() *EnvironmentsClient`
- New function `*ClientFactory.NewMetadataSchemasClient() *MetadataSchemasClient`
- New function `*ClientFactory.NewWorkspacesClient() *WorkspacesClient`
- New function `NewDeploymentsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DeploymentsClient, error)`
- New function `*DeploymentsClient.CreateOrUpdate(context.Context, string, string, string, string, string, Deployment, *DeploymentsClientCreateOrUpdateOptions) (DeploymentsClientCreateOrUpdateResponse, error)`
- New function `*DeploymentsClient.Delete(context.Context, string, string, string, string, string, *DeploymentsClientDeleteOptions) (DeploymentsClientDeleteResponse, error)`
- New function `*DeploymentsClient.Get(context.Context, string, string, string, string, string, *DeploymentsClientGetOptions) (DeploymentsClientGetResponse, error)`
- New function `*DeploymentsClient.Head(context.Context, string, string, string, string, string, *DeploymentsClientHeadOptions) (DeploymentsClientHeadResponse, error)`
- New function `*DeploymentsClient.NewListPager(string, string, string, string, *DeploymentsClientListOptions) *runtime.Pager[DeploymentsClientListResponse]`
- New function `NewEnvironmentsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*EnvironmentsClient, error)`
- New function `*EnvironmentsClient.CreateOrUpdate(context.Context, string, string, string, string, Environment, *EnvironmentsClientCreateOrUpdateOptions) (EnvironmentsClientCreateOrUpdateResponse, error)`
- New function `*EnvironmentsClient.Delete(context.Context, string, string, string, string, *EnvironmentsClientDeleteOptions) (EnvironmentsClientDeleteResponse, error)`
- New function `*EnvironmentsClient.Get(context.Context, string, string, string, string, *EnvironmentsClientGetOptions) (EnvironmentsClientGetResponse, error)`
- New function `*EnvironmentsClient.Head(context.Context, string, string, string, string, *EnvironmentsClientHeadOptions) (EnvironmentsClientHeadResponse, error)`
- New function `*EnvironmentsClient.NewListPager(string, string, string, *EnvironmentsClientListOptions) *runtime.Pager[EnvironmentsClientListResponse]`
- New function `NewMetadataSchemasClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MetadataSchemasClient, error)`
- New function `*MetadataSchemasClient.CreateOrUpdate(context.Context, string, string, string, MetadataSchema, *MetadataSchemasClientCreateOrUpdateOptions) (MetadataSchemasClientCreateOrUpdateResponse, error)`
- New function `*MetadataSchemasClient.Delete(context.Context, string, string, string, *MetadataSchemasClientDeleteOptions) (MetadataSchemasClientDeleteResponse, error)`
- New function `*MetadataSchemasClient.Get(context.Context, string, string, string, *MetadataSchemasClientGetOptions) (MetadataSchemasClientGetResponse, error)`
- New function `*MetadataSchemasClient.Head(context.Context, string, string, string, *MetadataSchemasClientHeadOptions) (MetadataSchemasClientHeadResponse, error)`
- New function `*MetadataSchemasClient.NewListPager(string, string, *MetadataSchemasClientListOptions) *runtime.Pager[MetadataSchemasClientListResponse]`
- New function `*ServicesClient.BeginExportMetadataSchema(context.Context, string, string, MetadataSchemaExportRequest, *ServicesClientBeginExportMetadataSchemaOptions) (*runtime.Poller[ServicesClientExportMetadataSchemaResponse], error)`
- New function `NewWorkspacesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*WorkspacesClient, error)`
- New function `*WorkspacesClient.CreateOrUpdate(context.Context, string, string, string, Workspace, *WorkspacesClientCreateOrUpdateOptions) (WorkspacesClientCreateOrUpdateResponse, error)`
- New function `*WorkspacesClient.Delete(context.Context, string, string, string, *WorkspacesClientDeleteOptions) (WorkspacesClientDeleteResponse, error)`
- New function `*WorkspacesClient.Get(context.Context, string, string, string, *WorkspacesClientGetOptions) (WorkspacesClientGetResponse, error)`
- New function `*WorkspacesClient.Head(context.Context, string, string, string, *WorkspacesClientHeadOptions) (WorkspacesClientHeadResponse, error)`
- New function `*WorkspacesClient.NewListPager(string, string, *WorkspacesClientListOptions) *runtime.Pager[WorkspacesClientListResponse]`
- New struct `API`
- New struct `APIDefinition`
- New struct `APIDefinitionListResult`
- New struct `APIDefinitionProperties`
- New struct `APIDefinitionPropertiesSpecification`
- New struct `APIListResult`
- New struct `APIProperties`
- New struct `APISpecExportResult`
- New struct `APISpecImportRequest`
- New struct `APISpecImportRequestSpecification`
- New struct `APIVersion`
- New struct `APIVersionListResult`
- New struct `APIVersionProperties`
- New struct `Contact`
- New struct `Deployment`
- New struct `DeploymentListResult`
- New struct `DeploymentProperties`
- New struct `DeploymentServer`
- New struct `Environment`
- New struct `EnvironmentListResult`
- New struct `EnvironmentProperties`
- New struct `EnvironmentServer`
- New struct `ExternalDocumentation`
- New struct `License`
- New struct `MetadataAssignment`
- New struct `MetadataSchema`
- New struct `MetadataSchemaExportRequest`
- New struct `MetadataSchemaExportResult`
- New struct `MetadataSchemaListResult`
- New struct `MetadataSchemaProperties`
- New struct `Onboarding`
- New struct `ServiceListResult`
- New struct `TermsOfService`
- New struct `Workspace`
- New struct `WorkspaceListResult`
- New struct `WorkspaceProperties`
- New field `Identity`, `Tags` in struct `ServiceUpdate`
- New anonymous field `ServiceListResult` in struct `ServicesClientListByResourceGroupResponse`
- New anonymous field `ServiceListResult` in struct `ServicesClientListBySubscriptionResponse`


## 0.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.1.0 (2023-08-25)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apicenter/armapicenter` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).