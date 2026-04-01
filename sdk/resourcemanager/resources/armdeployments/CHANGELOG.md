# Release History

## 2.0.0 (2026-04-01)
### Breaking Changes

- Function `*ClientFactory.NewDeploymentsClient` has been removed
- Function `NewDeploymentsClient` has been removed
- Function `*DeploymentsClient.CalculateTemplateHash` has been removed
- Function `*DeploymentsClient.Cancel` has been removed
- Function `*DeploymentsClient.CancelAtManagementGroupScope` has been removed
- Function `*DeploymentsClient.CancelAtScope` has been removed
- Function `*DeploymentsClient.CancelAtSubscriptionScope` has been removed
- Function `*DeploymentsClient.CancelAtTenantScope` has been removed
- Function `*DeploymentsClient.CheckExistence` has been removed
- Function `*DeploymentsClient.CheckExistenceAtManagementGroupScope` has been removed
- Function `*DeploymentsClient.CheckExistenceAtScope` has been removed
- Function `*DeploymentsClient.CheckExistenceAtSubscriptionScope` has been removed
- Function `*DeploymentsClient.CheckExistenceAtTenantScope` has been removed
- Function `*DeploymentsClient.BeginCreateOrUpdate` has been removed
- Function `*DeploymentsClient.BeginCreateOrUpdateAtManagementGroupScope` has been removed
- Function `*DeploymentsClient.BeginCreateOrUpdateAtScope` has been removed
- Function `*DeploymentsClient.BeginCreateOrUpdateAtSubscriptionScope` has been removed
- Function `*DeploymentsClient.BeginCreateOrUpdateAtTenantScope` has been removed
- Function `*DeploymentsClient.BeginDelete` has been removed
- Function `*DeploymentsClient.BeginDeleteAtManagementGroupScope` has been removed
- Function `*DeploymentsClient.BeginDeleteAtScope` has been removed
- Function `*DeploymentsClient.BeginDeleteAtSubscriptionScope` has been removed
- Function `*DeploymentsClient.BeginDeleteAtTenantScope` has been removed
- Function `*DeploymentsClient.ExportTemplate` has been removed
- Function `*DeploymentsClient.ExportTemplateAtManagementGroupScope` has been removed
- Function `*DeploymentsClient.ExportTemplateAtScope` has been removed
- Function `*DeploymentsClient.ExportTemplateAtSubscriptionScope` has been removed
- Function `*DeploymentsClient.ExportTemplateAtTenantScope` has been removed
- Function `*DeploymentsClient.Get` has been removed
- Function `*DeploymentsClient.GetAtManagementGroupScope` has been removed
- Function `*DeploymentsClient.GetAtScope` has been removed
- Function `*DeploymentsClient.GetAtSubscriptionScope` has been removed
- Function `*DeploymentsClient.GetAtTenantScope` has been removed
- Function `*DeploymentsClient.NewListAtManagementGroupScopePager` has been removed
- Function `*DeploymentsClient.NewListAtScopePager` has been removed
- Function `*DeploymentsClient.NewListAtSubscriptionScopePager` has been removed
- Function `*DeploymentsClient.NewListAtTenantScopePager` has been removed
- Function `*DeploymentsClient.NewListByResourceGroupPager` has been removed
- Function `*DeploymentsClient.BeginValidate` has been removed
- Function `*DeploymentsClient.BeginValidateAtManagementGroupScope` has been removed
- Function `*DeploymentsClient.BeginValidateAtScope` has been removed
- Function `*DeploymentsClient.BeginValidateAtSubscriptionScope` has been removed
- Function `*DeploymentsClient.BeginValidateAtTenantScope` has been removed
- Function `*DeploymentsClient.BeginWhatIf` has been removed
- Function `*DeploymentsClient.BeginWhatIfAtManagementGroupScope` has been removed
- Function `*DeploymentsClient.BeginWhatIfAtSubscriptionScope` has been removed
- Function `*DeploymentsClient.BeginWhatIfAtTenantScope` has been removed
- Struct `DeploymentExtendedFilter` has been removed

### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New function `NewClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*Client, error)`
- New function `*Client.CalculateTemplateHash(ctx context.Context, templateParam any, options *ClientCalculateTemplateHashOptions) (ClientCalculateTemplateHashResponse, error)`
- New function `*Client.Cancel(ctx context.Context, resourceGroupName string, deploymentName string, options *ClientCancelOptions) (ClientCancelResponse, error)`
- New function `*Client.CancelAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, options *ClientCancelAtManagementGroupScopeOptions) (ClientCancelAtManagementGroupScopeResponse, error)`
- New function `*Client.CancelAtScope(ctx context.Context, scope string, deploymentName string, options *ClientCancelAtScopeOptions) (ClientCancelAtScopeResponse, error)`
- New function `*Client.CancelAtSubscriptionScope(ctx context.Context, deploymentName string, options *ClientCancelAtSubscriptionScopeOptions) (ClientCancelAtSubscriptionScopeResponse, error)`
- New function `*Client.CancelAtTenantScope(ctx context.Context, deploymentName string, options *ClientCancelAtTenantScopeOptions) (ClientCancelAtTenantScopeResponse, error)`
- New function `*Client.CheckExistence(ctx context.Context, resourceGroupName string, deploymentName string, options *ClientCheckExistenceOptions) (ClientCheckExistenceResponse, error)`
- New function `*Client.CheckExistenceAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, options *ClientCheckExistenceAtManagementGroupScopeOptions) (ClientCheckExistenceAtManagementGroupScopeResponse, error)`
- New function `*Client.CheckExistenceAtScope(ctx context.Context, scope string, deploymentName string, options *ClientCheckExistenceAtScopeOptions) (ClientCheckExistenceAtScopeResponse, error)`
- New function `*Client.CheckExistenceAtSubscriptionScope(ctx context.Context, deploymentName string, options *ClientCheckExistenceAtSubscriptionScopeOptions) (ClientCheckExistenceAtSubscriptionScopeResponse, error)`
- New function `*Client.CheckExistenceAtTenantScope(ctx context.Context, deploymentName string, options *ClientCheckExistenceAtTenantScopeOptions) (ClientCheckExistenceAtTenantScopeResponse, error)`
- New function `*Client.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, deploymentName string, parameters Deployment, options *ClientBeginCreateOrUpdateOptions) (*runtime.Poller[ClientCreateOrUpdateResponse], error)`
- New function `*Client.BeginCreateOrUpdateAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, parameters ScopedDeployment, options *ClientBeginCreateOrUpdateAtManagementGroupScopeOptions) (*runtime.Poller[ClientCreateOrUpdateAtManagementGroupScopeResponse], error)`
- New function `*Client.BeginCreateOrUpdateAtScope(ctx context.Context, scope string, deploymentName string, parameters Deployment, options *ClientBeginCreateOrUpdateAtScopeOptions) (*runtime.Poller[ClientCreateOrUpdateAtScopeResponse], error)`
- New function `*Client.BeginCreateOrUpdateAtSubscriptionScope(ctx context.Context, deploymentName string, parameters Deployment, options *ClientBeginCreateOrUpdateAtSubscriptionScopeOptions) (*runtime.Poller[ClientCreateOrUpdateAtSubscriptionScopeResponse], error)`
- New function `*Client.BeginCreateOrUpdateAtTenantScope(ctx context.Context, deploymentName string, parameters ScopedDeployment, options *ClientBeginCreateOrUpdateAtTenantScopeOptions) (*runtime.Poller[ClientCreateOrUpdateAtTenantScopeResponse], error)`
- New function `*Client.BeginDelete(ctx context.Context, resourceGroupName string, deploymentName string, options *ClientBeginDeleteOptions) (*runtime.Poller[ClientDeleteResponse], error)`
- New function `*Client.BeginDeleteAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, options *ClientBeginDeleteAtManagementGroupScopeOptions) (*runtime.Poller[ClientDeleteAtManagementGroupScopeResponse], error)`
- New function `*Client.BeginDeleteAtScope(ctx context.Context, scope string, deploymentName string, options *ClientBeginDeleteAtScopeOptions) (*runtime.Poller[ClientDeleteAtScopeResponse], error)`
- New function `*Client.BeginDeleteAtSubscriptionScope(ctx context.Context, deploymentName string, options *ClientBeginDeleteAtSubscriptionScopeOptions) (*runtime.Poller[ClientDeleteAtSubscriptionScopeResponse], error)`
- New function `*Client.BeginDeleteAtTenantScope(ctx context.Context, deploymentName string, options *ClientBeginDeleteAtTenantScopeOptions) (*runtime.Poller[ClientDeleteAtTenantScopeResponse], error)`
- New function `*Client.ExportTemplate(ctx context.Context, resourceGroupName string, deploymentName string, options *ClientExportTemplateOptions) (ClientExportTemplateResponse, error)`
- New function `*Client.ExportTemplateAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, options *ClientExportTemplateAtManagementGroupScopeOptions) (ClientExportTemplateAtManagementGroupScopeResponse, error)`
- New function `*Client.ExportTemplateAtScope(ctx context.Context, scope string, deploymentName string, options *ClientExportTemplateAtScopeOptions) (ClientExportTemplateAtScopeResponse, error)`
- New function `*Client.ExportTemplateAtSubscriptionScope(ctx context.Context, deploymentName string, options *ClientExportTemplateAtSubscriptionScopeOptions) (ClientExportTemplateAtSubscriptionScopeResponse, error)`
- New function `*Client.ExportTemplateAtTenantScope(ctx context.Context, deploymentName string, options *ClientExportTemplateAtTenantScopeOptions) (ClientExportTemplateAtTenantScopeResponse, error)`
- New function `*Client.Get(ctx context.Context, resourceGroupName string, deploymentName string, options *ClientGetOptions) (ClientGetResponse, error)`
- New function `*Client.GetAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, options *ClientGetAtManagementGroupScopeOptions) (ClientGetAtManagementGroupScopeResponse, error)`
- New function `*Client.GetAtScope(ctx context.Context, scope string, deploymentName string, options *ClientGetAtScopeOptions) (ClientGetAtScopeResponse, error)`
- New function `*Client.GetAtSubscriptionScope(ctx context.Context, deploymentName string, options *ClientGetAtSubscriptionScopeOptions) (ClientGetAtSubscriptionScopeResponse, error)`
- New function `*Client.GetAtTenantScope(ctx context.Context, deploymentName string, options *ClientGetAtTenantScopeOptions) (ClientGetAtTenantScopeResponse, error)`
- New function `*Client.NewListAtManagementGroupScopePager(groupID string, options *ClientListAtManagementGroupScopeOptions) *runtime.Pager[ClientListAtManagementGroupScopeResponse]`
- New function `*Client.NewListAtScopePager(scope string, options *ClientListAtScopeOptions) *runtime.Pager[ClientListAtScopeResponse]`
- New function `*Client.NewListAtSubscriptionScopePager(options *ClientListAtSubscriptionScopeOptions) *runtime.Pager[ClientListAtSubscriptionScopeResponse]`
- New function `*Client.NewListAtTenantScopePager(options *ClientListAtTenantScopeOptions) *runtime.Pager[ClientListAtTenantScopeResponse]`
- New function `*Client.NewListByResourceGroupPager(resourceGroupName string, options *ClientListByResourceGroupOptions) *runtime.Pager[ClientListByResourceGroupResponse]`
- New function `*Client.BeginValidate(ctx context.Context, resourceGroupName string, deploymentName string, parameters Deployment, options *ClientBeginValidateOptions) (*runtime.Poller[ClientValidateResponse], error)`
- New function `*Client.BeginValidateAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, parameters ScopedDeployment, options *ClientBeginValidateAtManagementGroupScopeOptions) (*runtime.Poller[ClientValidateAtManagementGroupScopeResponse], error)`
- New function `*Client.BeginValidateAtScope(ctx context.Context, scope string, deploymentName string, parameters Deployment, options *ClientBeginValidateAtScopeOptions) (*runtime.Poller[ClientValidateAtScopeResponse], error)`
- New function `*Client.BeginValidateAtSubscriptionScope(ctx context.Context, deploymentName string, parameters Deployment, options *ClientBeginValidateAtSubscriptionScopeOptions) (*runtime.Poller[ClientValidateAtSubscriptionScopeResponse], error)`
- New function `*Client.BeginValidateAtTenantScope(ctx context.Context, deploymentName string, parameters ScopedDeployment, options *ClientBeginValidateAtTenantScopeOptions) (*runtime.Poller[ClientValidateAtTenantScopeResponse], error)`
- New function `*Client.BeginWhatIf(ctx context.Context, resourceGroupName string, deploymentName string, parameters DeploymentWhatIf, options *ClientBeginWhatIfOptions) (*runtime.Poller[ClientWhatIfResponse], error)`
- New function `*Client.BeginWhatIfAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, parameters ScopedDeploymentWhatIf, options *ClientBeginWhatIfAtManagementGroupScopeOptions) (*runtime.Poller[ClientWhatIfAtManagementGroupScopeResponse], error)`
- New function `*Client.BeginWhatIfAtSubscriptionScope(ctx context.Context, deploymentName string, parameters DeploymentWhatIf, options *ClientBeginWhatIfAtSubscriptionScopeOptions) (*runtime.Poller[ClientWhatIfAtSubscriptionScopeResponse], error)`
- New function `*Client.BeginWhatIfAtTenantScope(ctx context.Context, deploymentName string, parameters ScopedDeploymentWhatIf, options *ClientBeginWhatIfAtTenantScopeOptions) (*runtime.Poller[ClientWhatIfAtTenantScopeResponse], error)`
- New function `*ClientFactory.NewClient() *Client`
- New struct `SystemData`
- New field `SystemData` in struct `DeploymentExtended`


## 1.0.0 (2025-12-23)
### Other Changes

- Release stable version.

## 0.2.0 (2025-12-16)
### Breaking Changes

- Function `NewClient` has been removed
- Function `*Client.CalculateTemplateHash` has been removed
- Function `*Client.Cancel` has been removed
- Function `*Client.CancelAtManagementGroupScope` has been removed
- Function `*Client.CancelAtScope` has been removed
- Function `*Client.CancelAtSubscriptionScope` has been removed
- Function `*Client.CancelAtTenantScope` has been removed
- Function `*Client.CheckExistence` has been removed
- Function `*Client.CheckExistenceAtManagementGroupScope` has been removed
- Function `*Client.CheckExistenceAtScope` has been removed
- Function `*Client.CheckExistenceAtSubscriptionScope` has been removed
- Function `*Client.CheckExistenceAtTenantScope` has been removed
- Function `*Client.BeginCreateOrUpdate` has been removed
- Function `*Client.BeginCreateOrUpdateAtManagementGroupScope` has been removed
- Function `*Client.BeginCreateOrUpdateAtScope` has been removed
- Function `*Client.BeginCreateOrUpdateAtSubscriptionScope` has been removed
- Function `*Client.BeginCreateOrUpdateAtTenantScope` has been removed
- Function `*Client.BeginDelete` has been removed
- Function `*Client.BeginDeleteAtManagementGroupScope` has been removed
- Function `*Client.BeginDeleteAtScope` has been removed
- Function `*Client.BeginDeleteAtSubscriptionScope` has been removed
- Function `*Client.BeginDeleteAtTenantScope` has been removed
- Function `*Client.ExportTemplate` has been removed
- Function `*Client.ExportTemplateAtManagementGroupScope` has been removed
- Function `*Client.ExportTemplateAtScope` has been removed
- Function `*Client.ExportTemplateAtSubscriptionScope` has been removed
- Function `*Client.ExportTemplateAtTenantScope` has been removed
- Function `*Client.Get` has been removed
- Function `*Client.GetAtManagementGroupScope` has been removed
- Function `*Client.GetAtScope` has been removed
- Function `*Client.GetAtSubscriptionScope` has been removed
- Function `*Client.GetAtTenantScope` has been removed
- Function `*Client.NewListAtManagementGroupScopePager` has been removed
- Function `*Client.NewListAtScopePager` has been removed
- Function `*Client.NewListAtSubscriptionScopePager` has been removed
- Function `*Client.NewListAtTenantScopePager` has been removed
- Function `*Client.NewListByResourceGroupPager` has been removed
- Function `*Client.BeginValidate` has been removed
- Function `*Client.BeginValidateAtManagementGroupScope` has been removed
- Function `*Client.BeginValidateAtScope` has been removed
- Function `*Client.BeginValidateAtSubscriptionScope` has been removed
- Function `*Client.BeginValidateAtTenantScope` has been removed
- Function `*Client.BeginWhatIf` has been removed
- Function `*Client.BeginWhatIfAtManagementGroupScope` has been removed
- Function `*Client.BeginWhatIfAtSubscriptionScope` has been removed
- Function `*Client.BeginWhatIfAtTenantScope` has been removed
- Function `*ClientFactory.NewClient` has been removed

### Features Added

- New function `*ClientFactory.NewDeploymentsClient() *DeploymentsClient`
- New function `NewDeploymentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DeploymentsClient, error)`
- New function `*DeploymentsClient.CalculateTemplateHash(ctx context.Context, templateParam any, options *DeploymentsClientCalculateTemplateHashOptions) (DeploymentsClientCalculateTemplateHashResponse, error)`
- New function `*DeploymentsClient.Cancel(ctx context.Context, resourceGroupName string, deploymentName string, options *DeploymentsClientCancelOptions) (DeploymentsClientCancelResponse, error)`
- New function `*DeploymentsClient.CancelAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, options *DeploymentsClientCancelAtManagementGroupScopeOptions) (DeploymentsClientCancelAtManagementGroupScopeResponse, error)`
- New function `*DeploymentsClient.CancelAtScope(ctx context.Context, scope string, deploymentName string, options *DeploymentsClientCancelAtScopeOptions) (DeploymentsClientCancelAtScopeResponse, error)`
- New function `*DeploymentsClient.CancelAtSubscriptionScope(ctx context.Context, deploymentName string, options *DeploymentsClientCancelAtSubscriptionScopeOptions) (DeploymentsClientCancelAtSubscriptionScopeResponse, error)`
- New function `*DeploymentsClient.CancelAtTenantScope(ctx context.Context, deploymentName string, options *DeploymentsClientCancelAtTenantScopeOptions) (DeploymentsClientCancelAtTenantScopeResponse, error)`
- New function `*DeploymentsClient.CheckExistence(ctx context.Context, resourceGroupName string, deploymentName string, options *DeploymentsClientCheckExistenceOptions) (DeploymentsClientCheckExistenceResponse, error)`
- New function `*DeploymentsClient.CheckExistenceAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, options *DeploymentsClientCheckExistenceAtManagementGroupScopeOptions) (DeploymentsClientCheckExistenceAtManagementGroupScopeResponse, error)`
- New function `*DeploymentsClient.CheckExistenceAtScope(ctx context.Context, scope string, deploymentName string, options *DeploymentsClientCheckExistenceAtScopeOptions) (DeploymentsClientCheckExistenceAtScopeResponse, error)`
- New function `*DeploymentsClient.CheckExistenceAtSubscriptionScope(ctx context.Context, deploymentName string, options *DeploymentsClientCheckExistenceAtSubscriptionScopeOptions) (DeploymentsClientCheckExistenceAtSubscriptionScopeResponse, error)`
- New function `*DeploymentsClient.CheckExistenceAtTenantScope(ctx context.Context, deploymentName string, options *DeploymentsClientCheckExistenceAtTenantScopeOptions) (DeploymentsClientCheckExistenceAtTenantScopeResponse, error)`
- New function `*DeploymentsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, deploymentName string, parameters Deployment, options *DeploymentsClientBeginCreateOrUpdateOptions) (*runtime.Poller[DeploymentsClientCreateOrUpdateResponse], error)`
- New function `*DeploymentsClient.BeginCreateOrUpdateAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, parameters ScopedDeployment, options *DeploymentsClientBeginCreateOrUpdateAtManagementGroupScopeOptions) (*runtime.Poller[DeploymentsClientCreateOrUpdateAtManagementGroupScopeResponse], error)`
- New function `*DeploymentsClient.BeginCreateOrUpdateAtScope(ctx context.Context, scope string, deploymentName string, parameters Deployment, options *DeploymentsClientBeginCreateOrUpdateAtScopeOptions) (*runtime.Poller[DeploymentsClientCreateOrUpdateAtScopeResponse], error)`
- New function `*DeploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(ctx context.Context, deploymentName string, parameters Deployment, options *DeploymentsClientBeginCreateOrUpdateAtSubscriptionScopeOptions) (*runtime.Poller[DeploymentsClientCreateOrUpdateAtSubscriptionScopeResponse], error)`
- New function `*DeploymentsClient.BeginCreateOrUpdateAtTenantScope(ctx context.Context, deploymentName string, parameters ScopedDeployment, options *DeploymentsClientBeginCreateOrUpdateAtTenantScopeOptions) (*runtime.Poller[DeploymentsClientCreateOrUpdateAtTenantScopeResponse], error)`
- New function `*DeploymentsClient.BeginDelete(ctx context.Context, resourceGroupName string, deploymentName string, options *DeploymentsClientBeginDeleteOptions) (*runtime.Poller[DeploymentsClientDeleteResponse], error)`
- New function `*DeploymentsClient.BeginDeleteAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, options *DeploymentsClientBeginDeleteAtManagementGroupScopeOptions) (*runtime.Poller[DeploymentsClientDeleteAtManagementGroupScopeResponse], error)`
- New function `*DeploymentsClient.BeginDeleteAtScope(ctx context.Context, scope string, deploymentName string, options *DeploymentsClientBeginDeleteAtScopeOptions) (*runtime.Poller[DeploymentsClientDeleteAtScopeResponse], error)`
- New function `*DeploymentsClient.BeginDeleteAtSubscriptionScope(ctx context.Context, deploymentName string, options *DeploymentsClientBeginDeleteAtSubscriptionScopeOptions) (*runtime.Poller[DeploymentsClientDeleteAtSubscriptionScopeResponse], error)`
- New function `*DeploymentsClient.BeginDeleteAtTenantScope(ctx context.Context, deploymentName string, options *DeploymentsClientBeginDeleteAtTenantScopeOptions) (*runtime.Poller[DeploymentsClientDeleteAtTenantScopeResponse], error)`
- New function `*DeploymentsClient.ExportTemplate(ctx context.Context, resourceGroupName string, deploymentName string, options *DeploymentsClientExportTemplateOptions) (DeploymentsClientExportTemplateResponse, error)`
- New function `*DeploymentsClient.ExportTemplateAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, options *DeploymentsClientExportTemplateAtManagementGroupScopeOptions) (DeploymentsClientExportTemplateAtManagementGroupScopeResponse, error)`
- New function `*DeploymentsClient.ExportTemplateAtScope(ctx context.Context, scope string, deploymentName string, options *DeploymentsClientExportTemplateAtScopeOptions) (DeploymentsClientExportTemplateAtScopeResponse, error)`
- New function `*DeploymentsClient.ExportTemplateAtSubscriptionScope(ctx context.Context, deploymentName string, options *DeploymentsClientExportTemplateAtSubscriptionScopeOptions) (DeploymentsClientExportTemplateAtSubscriptionScopeResponse, error)`
- New function `*DeploymentsClient.ExportTemplateAtTenantScope(ctx context.Context, deploymentName string, options *DeploymentsClientExportTemplateAtTenantScopeOptions) (DeploymentsClientExportTemplateAtTenantScopeResponse, error)`
- New function `*DeploymentsClient.Get(ctx context.Context, resourceGroupName string, deploymentName string, options *DeploymentsClientGetOptions) (DeploymentsClientGetResponse, error)`
- New function `*DeploymentsClient.GetAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, options *DeploymentsClientGetAtManagementGroupScopeOptions) (DeploymentsClientGetAtManagementGroupScopeResponse, error)`
- New function `*DeploymentsClient.GetAtScope(ctx context.Context, scope string, deploymentName string, options *DeploymentsClientGetAtScopeOptions) (DeploymentsClientGetAtScopeResponse, error)`
- New function `*DeploymentsClient.GetAtSubscriptionScope(ctx context.Context, deploymentName string, options *DeploymentsClientGetAtSubscriptionScopeOptions) (DeploymentsClientGetAtSubscriptionScopeResponse, error)`
- New function `*DeploymentsClient.GetAtTenantScope(ctx context.Context, deploymentName string, options *DeploymentsClientGetAtTenantScopeOptions) (DeploymentsClientGetAtTenantScopeResponse, error)`
- New function `*DeploymentsClient.NewListAtManagementGroupScopePager(groupID string, options *DeploymentsClientListAtManagementGroupScopeOptions) *runtime.Pager[DeploymentsClientListAtManagementGroupScopeResponse]`
- New function `*DeploymentsClient.NewListAtScopePager(scope string, options *DeploymentsClientListAtScopeOptions) *runtime.Pager[DeploymentsClientListAtScopeResponse]`
- New function `*DeploymentsClient.NewListAtSubscriptionScopePager(options *DeploymentsClientListAtSubscriptionScopeOptions) *runtime.Pager[DeploymentsClientListAtSubscriptionScopeResponse]`
- New function `*DeploymentsClient.NewListAtTenantScopePager(options *DeploymentsClientListAtTenantScopeOptions) *runtime.Pager[DeploymentsClientListAtTenantScopeResponse]`
- New function `*DeploymentsClient.NewListByResourceGroupPager(resourceGroupName string, options *DeploymentsClientListByResourceGroupOptions) *runtime.Pager[DeploymentsClientListByResourceGroupResponse]`
- New function `*DeploymentsClient.BeginValidate(ctx context.Context, resourceGroupName string, deploymentName string, parameters Deployment, options *DeploymentsClientBeginValidateOptions) (*runtime.Poller[DeploymentsClientValidateResponse], error)`
- New function `*DeploymentsClient.BeginValidateAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, parameters ScopedDeployment, options *DeploymentsClientBeginValidateAtManagementGroupScopeOptions) (*runtime.Poller[DeploymentsClientValidateAtManagementGroupScopeResponse], error)`
- New function `*DeploymentsClient.BeginValidateAtScope(ctx context.Context, scope string, deploymentName string, parameters Deployment, options *DeploymentsClientBeginValidateAtScopeOptions) (*runtime.Poller[DeploymentsClientValidateAtScopeResponse], error)`
- New function `*DeploymentsClient.BeginValidateAtSubscriptionScope(ctx context.Context, deploymentName string, parameters Deployment, options *DeploymentsClientBeginValidateAtSubscriptionScopeOptions) (*runtime.Poller[DeploymentsClientValidateAtSubscriptionScopeResponse], error)`
- New function `*DeploymentsClient.BeginValidateAtTenantScope(ctx context.Context, deploymentName string, parameters ScopedDeployment, options *DeploymentsClientBeginValidateAtTenantScopeOptions) (*runtime.Poller[DeploymentsClientValidateAtTenantScopeResponse], error)`
- New function `*DeploymentsClient.BeginWhatIf(ctx context.Context, resourceGroupName string, deploymentName string, parameters DeploymentWhatIf, options *DeploymentsClientBeginWhatIfOptions) (*runtime.Poller[DeploymentsClientWhatIfResponse], error)`
- New function `*DeploymentsClient.BeginWhatIfAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, parameters ScopedDeploymentWhatIf, options *DeploymentsClientBeginWhatIfAtManagementGroupScopeOptions) (*runtime.Poller[DeploymentsClientWhatIfAtManagementGroupScopeResponse], error)`
- New function `*DeploymentsClient.BeginWhatIfAtSubscriptionScope(ctx context.Context, deploymentName string, parameters DeploymentWhatIf, options *DeploymentsClientBeginWhatIfAtSubscriptionScopeOptions) (*runtime.Poller[DeploymentsClientWhatIfAtSubscriptionScopeResponse], error)`
- New function `*DeploymentsClient.BeginWhatIfAtTenantScope(ctx context.Context, deploymentName string, parameters ScopedDeploymentWhatIf, options *DeploymentsClientBeginWhatIfAtTenantScopeOptions) (*runtime.Poller[DeploymentsClientWhatIfAtTenantScopeResponse], error)`


## 0.1.0 (2025-07-24)
### Other Changes

This package is split from `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources`.

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).