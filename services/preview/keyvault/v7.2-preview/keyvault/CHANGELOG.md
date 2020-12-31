Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

### Removed Funcs

1. *FullBackupFuture.Result(BaseClient) (FullBackupOperation, error)
1. *FullRestoreOperationFuture.Result(BaseClient) (RestoreOperation, error)
1. *HSMSecurityDomainUploadFuture.Result(HSMSecurityDomainClient) (SecurityDomainOperationStatus, error)
1. *SelectiveKeyRestoreOperationMethodFuture.Result(BaseClient) (SelectiveKeyRestoreOperation, error)

## Struct Changes

### Removed Struct Fields

1. FullBackupFuture.azure.Future
1. FullRestoreOperationFuture.azure.Future
1. HSMSecurityDomainUploadFuture.azure.Future
1. SelectiveKeyRestoreOperationMethodFuture.azure.Future

### New Funcs

1. RoleDefinitionsClient.CreateOrUpdate(context.Context, string, string, string, RoleDefinitionCreateParameters) (RoleDefinition, error)
1. RoleDefinitionsClient.CreateOrUpdatePreparer(context.Context, string, string, string, RoleDefinitionCreateParameters) (*http.Request, error)
1. RoleDefinitionsClient.CreateOrUpdateResponder(*http.Response) (RoleDefinition, error)
1. RoleDefinitionsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. RoleDefinitionsClient.Delete(context.Context, string, string, string) (RoleDefinition, error)
1. RoleDefinitionsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. RoleDefinitionsClient.DeleteResponder(*http.Response) (RoleDefinition, error)
1. RoleDefinitionsClient.DeleteSender(*http.Request) (*http.Response, error)
1. RoleDefinitionsClient.Get(context.Context, string, string, string) (RoleDefinition, error)
1. RoleDefinitionsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. RoleDefinitionsClient.GetResponder(*http.Response) (RoleDefinition, error)
1. RoleDefinitionsClient.GetSender(*http.Request) (*http.Response, error)

## Struct Changes

### New Structs

1. RoleDefinitionCreateParameters

### New Struct Fields

1. FullBackupFuture.Result
1. FullBackupFuture.azure.FutureAPI
1. FullRestoreOperationFuture.Result
1. FullRestoreOperationFuture.azure.FutureAPI
1. HSMSecurityDomainUploadFuture.Result
1. HSMSecurityDomainUploadFuture.azure.FutureAPI
1. KeyOperationResult.AdditionalAuthenticatedData
1. KeyOperationResult.AuthenticationTag
1. KeyOperationResult.Iv
1. RoleDefinition.autorest.Response
1. SelectiveKeyRestoreOperationMethodFuture.Result
1. SelectiveKeyRestoreOperationMethodFuture.azure.FutureAPI
