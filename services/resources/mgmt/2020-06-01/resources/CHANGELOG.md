Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82/specification/resources/resource-manager/readme.md tag: `package-resources-2020-06`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *CreateOrUpdateByIDFuture.Result(Client) (GenericResource, error)
1. *CreateOrUpdateFuture.Result(Client) (GenericResource, error)
1. *DeleteByIDFuture.Result(Client) (autorest.Response, error)
1. *DeleteFuture.Result(Client) (autorest.Response, error)
1. *DeploymentsCreateOrUpdateAtManagementGroupScopeFuture.Result(DeploymentsClient) (DeploymentExtended, error)
1. *DeploymentsCreateOrUpdateAtScopeFuture.Result(DeploymentsClient) (DeploymentExtended, error)
1. *DeploymentsCreateOrUpdateAtSubscriptionScopeFuture.Result(DeploymentsClient) (DeploymentExtended, error)
1. *DeploymentsCreateOrUpdateAtTenantScopeFuture.Result(DeploymentsClient) (DeploymentExtended, error)
1. *DeploymentsCreateOrUpdateFuture.Result(DeploymentsClient) (DeploymentExtended, error)
1. *DeploymentsDeleteAtManagementGroupScopeFuture.Result(DeploymentsClient) (autorest.Response, error)
1. *DeploymentsDeleteAtScopeFuture.Result(DeploymentsClient) (autorest.Response, error)
1. *DeploymentsDeleteAtSubscriptionScopeFuture.Result(DeploymentsClient) (autorest.Response, error)
1. *DeploymentsDeleteAtTenantScopeFuture.Result(DeploymentsClient) (autorest.Response, error)
1. *DeploymentsDeleteFuture.Result(DeploymentsClient) (autorest.Response, error)
1. *DeploymentsValidateAtManagementGroupScopeFuture.Result(DeploymentsClient) (DeploymentValidateResult, error)
1. *DeploymentsValidateAtScopeFuture.Result(DeploymentsClient) (DeploymentValidateResult, error)
1. *DeploymentsValidateAtSubscriptionScopeFuture.Result(DeploymentsClient) (DeploymentValidateResult, error)
1. *DeploymentsValidateAtTenantScopeFuture.Result(DeploymentsClient) (DeploymentValidateResult, error)
1. *DeploymentsValidateFuture.Result(DeploymentsClient) (DeploymentValidateResult, error)
1. *DeploymentsWhatIfAtManagementGroupScopeFuture.Result(DeploymentsClient) (WhatIfOperationResult, error)
1. *DeploymentsWhatIfAtSubscriptionScopeFuture.Result(DeploymentsClient) (WhatIfOperationResult, error)
1. *DeploymentsWhatIfAtTenantScopeFuture.Result(DeploymentsClient) (WhatIfOperationResult, error)
1. *DeploymentsWhatIfFuture.Result(DeploymentsClient) (WhatIfOperationResult, error)
1. *GroupsDeleteFuture.Result(GroupsClient) (autorest.Response, error)
1. *GroupsExportTemplateFuture.Result(GroupsClient) (GroupExportResult, error)
1. *MoveResourcesFuture.Result(Client) (autorest.Response, error)
1. *UpdateByIDFuture.Result(Client) (GenericResource, error)
1. *UpdateFuture.Result(Client) (GenericResource, error)
1. *ValidateMoveResourcesFuture.Result(Client) (autorest.Response, error)

## Struct Changes

### Removed Struct Fields

1. CreateOrUpdateByIDFuture.azure.Future
1. CreateOrUpdateFuture.azure.Future
1. DeleteByIDFuture.azure.Future
1. DeleteFuture.azure.Future
1. DeploymentsCreateOrUpdateAtManagementGroupScopeFuture.azure.Future
1. DeploymentsCreateOrUpdateAtScopeFuture.azure.Future
1. DeploymentsCreateOrUpdateAtSubscriptionScopeFuture.azure.Future
1. DeploymentsCreateOrUpdateAtTenantScopeFuture.azure.Future
1. DeploymentsCreateOrUpdateFuture.azure.Future
1. DeploymentsDeleteAtManagementGroupScopeFuture.azure.Future
1. DeploymentsDeleteAtScopeFuture.azure.Future
1. DeploymentsDeleteAtSubscriptionScopeFuture.azure.Future
1. DeploymentsDeleteAtTenantScopeFuture.azure.Future
1. DeploymentsDeleteFuture.azure.Future
1. DeploymentsValidateAtManagementGroupScopeFuture.azure.Future
1. DeploymentsValidateAtScopeFuture.azure.Future
1. DeploymentsValidateAtSubscriptionScopeFuture.azure.Future
1. DeploymentsValidateAtTenantScopeFuture.azure.Future
1. DeploymentsValidateFuture.azure.Future
1. DeploymentsWhatIfAtManagementGroupScopeFuture.azure.Future
1. DeploymentsWhatIfAtSubscriptionScopeFuture.azure.Future
1. DeploymentsWhatIfAtTenantScopeFuture.azure.Future
1. DeploymentsWhatIfFuture.azure.Future
1. GroupsDeleteFuture.azure.Future
1. GroupsExportTemplateFuture.azure.Future
1. MoveResourcesFuture.azure.Future
1. UpdateByIDFuture.azure.Future
1. UpdateFuture.azure.Future
1. ValidateMoveResourcesFuture.azure.Future

## Struct Changes

### New Struct Fields

1. CreateOrUpdateByIDFuture.Result
1. CreateOrUpdateByIDFuture.azure.FutureAPI
1. CreateOrUpdateFuture.Result
1. CreateOrUpdateFuture.azure.FutureAPI
1. DeleteByIDFuture.Result
1. DeleteByIDFuture.azure.FutureAPI
1. DeleteFuture.Result
1. DeleteFuture.azure.FutureAPI
1. DeploymentsCreateOrUpdateAtManagementGroupScopeFuture.Result
1. DeploymentsCreateOrUpdateAtManagementGroupScopeFuture.azure.FutureAPI
1. DeploymentsCreateOrUpdateAtScopeFuture.Result
1. DeploymentsCreateOrUpdateAtScopeFuture.azure.FutureAPI
1. DeploymentsCreateOrUpdateAtSubscriptionScopeFuture.Result
1. DeploymentsCreateOrUpdateAtSubscriptionScopeFuture.azure.FutureAPI
1. DeploymentsCreateOrUpdateAtTenantScopeFuture.Result
1. DeploymentsCreateOrUpdateAtTenantScopeFuture.azure.FutureAPI
1. DeploymentsCreateOrUpdateFuture.Result
1. DeploymentsCreateOrUpdateFuture.azure.FutureAPI
1. DeploymentsDeleteAtManagementGroupScopeFuture.Result
1. DeploymentsDeleteAtManagementGroupScopeFuture.azure.FutureAPI
1. DeploymentsDeleteAtScopeFuture.Result
1. DeploymentsDeleteAtScopeFuture.azure.FutureAPI
1. DeploymentsDeleteAtSubscriptionScopeFuture.Result
1. DeploymentsDeleteAtSubscriptionScopeFuture.azure.FutureAPI
1. DeploymentsDeleteAtTenantScopeFuture.Result
1. DeploymentsDeleteAtTenantScopeFuture.azure.FutureAPI
1. DeploymentsDeleteFuture.Result
1. DeploymentsDeleteFuture.azure.FutureAPI
1. DeploymentsValidateAtManagementGroupScopeFuture.Result
1. DeploymentsValidateAtManagementGroupScopeFuture.azure.FutureAPI
1. DeploymentsValidateAtScopeFuture.Result
1. DeploymentsValidateAtScopeFuture.azure.FutureAPI
1. DeploymentsValidateAtSubscriptionScopeFuture.Result
1. DeploymentsValidateAtSubscriptionScopeFuture.azure.FutureAPI
1. DeploymentsValidateAtTenantScopeFuture.Result
1. DeploymentsValidateAtTenantScopeFuture.azure.FutureAPI
1. DeploymentsValidateFuture.Result
1. DeploymentsValidateFuture.azure.FutureAPI
1. DeploymentsWhatIfAtManagementGroupScopeFuture.Result
1. DeploymentsWhatIfAtManagementGroupScopeFuture.azure.FutureAPI
1. DeploymentsWhatIfAtSubscriptionScopeFuture.Result
1. DeploymentsWhatIfAtSubscriptionScopeFuture.azure.FutureAPI
1. DeploymentsWhatIfAtTenantScopeFuture.Result
1. DeploymentsWhatIfAtTenantScopeFuture.azure.FutureAPI
1. DeploymentsWhatIfFuture.Result
1. DeploymentsWhatIfFuture.azure.FutureAPI
1. GroupsDeleteFuture.Result
1. GroupsDeleteFuture.azure.FutureAPI
1. GroupsExportTemplateFuture.Result
1. GroupsExportTemplateFuture.azure.FutureAPI
1. MoveResourcesFuture.Result
1. MoveResourcesFuture.azure.FutureAPI
1. UpdateByIDFuture.Result
1. UpdateByIDFuture.azure.FutureAPI
1. UpdateFuture.Result
1. UpdateFuture.azure.FutureAPI
1. ValidateMoveResourcesFuture.Result
1. ValidateMoveResourcesFuture.azure.FutureAPI
