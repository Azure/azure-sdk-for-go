
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewDomainServiceListResultPage` signature has been changed from `(func(context.Context, DomainServiceListResult) (DomainServiceListResult, error))` to `(DomainServiceListResult,func(context.Context, DomainServiceListResult) (DomainServiceListResult, error))`
- Function `NewOperationEntityListResultPage` signature has been changed from `(func(context.Context, OperationEntityListResult) (OperationEntityListResult, error))` to `(OperationEntityListResult,func(context.Context, OperationEntityListResult) (OperationEntityListResult, error))`
- Type of `DomainServiceProperties.HealthLastEvaluated` has been changed from `*date.Time` to `*date.TimeRFC1123`

## New Content

- Const `SyncKerberosPasswordsDisabled` is added
- Const `SyncKerberosPasswordsEnabled` is added
- Const `SyncOnPremPasswordsEnabled` is added
- Const `SyncOnPremPasswordsDisabled` is added
- Function `OuContainerListResult.MarshalJSON() ([]byte,error)` is added
- Function `NewOuContainerClient(string) OuContainerClient` is added
- Function `OuContainerClient.Create(context.Context,string,string,string,ContainerAccount) (OuContainerCreateFuture,error)` is added
- Function `OuContainerListResultPage.Values() []OuContainer` is added
- Function `*OuContainerUpdateFuture.Result(OuContainerClient) (OuContainer,error)` is added
- Function `OuContainerListResultIterator.NotDone() bool` is added
- Function `OuContainerListResultIterator.Response() OuContainerListResult` is added
- Function `*OuContainerCreateFuture.Result(OuContainerClient) (OuContainer,error)` is added
- Function `OuContainerClient.ListSender(*http.Request) (*http.Response,error)` is added
- Function `OuContainerClient.ListResponder(*http.Response) (OuContainerListResult,error)` is added
- Function `PossibleSyncKerberosPasswordsValues() []SyncKerberosPasswords` is added
- Function `OuContainerListResultPage.Response() OuContainerListResult` is added
- Function `OuContainerClient.GetResponder(*http.Response) (OuContainer,error)` is added
- Function `*OuContainer.UnmarshalJSON([]byte) error` is added
- Function `OuContainerOperationsClient.ListComplete(context.Context) (OperationEntityListResultIterator,error)` is added
- Function `OuContainer.MarshalJSON() ([]byte,error)` is added
- Function `*OuContainerDeleteFuture.Result(OuContainerClient) (autorest.Response,error)` is added
- Function `NewOuContainerListResultPage(OuContainerListResult,func(context.Context, OuContainerListResult) (OuContainerListResult, error)) OuContainerListResultPage` is added
- Function `OuContainerListResultIterator.Value() OuContainer` is added
- Function `OuContainerClient.GetPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `*OuContainerListResultIterator.Next() error` is added
- Function `OuContainerClient.UpdateResponder(*http.Response) (OuContainer,error)` is added
- Function `OuContainerClient.CreateResponder(*http.Response) (OuContainer,error)` is added
- Function `*OuContainerListResultPage.NextWithContext(context.Context) error` is added
- Function `OuContainerListResultPage.NotDone() bool` is added
- Function `OuContainerClient.DeletePreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `OuContainerClient.List(context.Context,string,string) (OuContainerListResultPage,error)` is added
- Function `OuContainerClient.UpdatePreparer(context.Context,string,string,string,ContainerAccount) (*http.Request,error)` is added
- Function `NewOuContainerListResultIterator(OuContainerListResultPage) OuContainerListResultIterator` is added
- Function `OuContainerClient.DeleteResponder(*http.Response) (autorest.Response,error)` is added
- Function `OuContainerClient.ListPreparer(context.Context,string,string) (*http.Request,error)` is added
- Function `OuContainerProperties.MarshalJSON() ([]byte,error)` is added
- Function `NewOuContainerOperationsClientWithBaseURI(string,string) OuContainerOperationsClient` is added
- Function `*OuContainerListResultPage.Next() error` is added
- Function `OuContainerClient.CreatePreparer(context.Context,string,string,string,ContainerAccount) (*http.Request,error)` is added
- Function `OuContainerClient.DeleteSender(*http.Request) (OuContainerDeleteFuture,error)` is added
- Function `OuContainerClient.UpdateSender(*http.Request) (OuContainerUpdateFuture,error)` is added
- Function `PossibleSyncOnPremPasswordsValues() []SyncOnPremPasswords` is added
- Function `OuContainerClient.Update(context.Context,string,string,string,ContainerAccount) (OuContainerUpdateFuture,error)` is added
- Function `DefaultErrorResponseError.MarshalJSON() ([]byte,error)` is added
- Function `NewOuContainerOperationsClient(string) OuContainerOperationsClient` is added
- Function `OuContainerOperationsClient.ListPreparer(context.Context) (*http.Request,error)` is added
- Function `*OuContainerListResultIterator.NextWithContext(context.Context) error` is added
- Function `OuContainerOperationsClient.ListResponder(*http.Response) (OperationEntityListResult,error)` is added
- Function `OuContainerClient.CreateSender(*http.Request) (OuContainerCreateFuture,error)` is added
- Function `OuContainerListResult.IsEmpty() bool` is added
- Function `OuContainerOperationsClient.List(context.Context) (OperationEntityListResultPage,error)` is added
- Function `OuContainerClient.ListComplete(context.Context,string,string) (OuContainerListResultIterator,error)` is added
- Function `OuContainerClient.Get(context.Context,string,string,string) (OuContainer,error)` is added
- Function `OuContainerClient.Delete(context.Context,string,string,string) (OuContainerDeleteFuture,error)` is added
- Function `NewOuContainerClientWithBaseURI(string,string) OuContainerClient` is added
- Function `OuContainerClient.GetSender(*http.Request) (*http.Response,error)` is added
- Function `OuContainerOperationsClient.ListSender(*http.Request) (*http.Response,error)` is added
- Struct `CloudError` is added
- Struct `CloudErrorBody` is added
- Struct `ContainerAccount` is added
- Struct `DefaultErrorResponse` is added
- Struct `DefaultErrorResponseError` is added
- Struct `DefaultErrorResponseErrorDetailsItem` is added
- Struct `ForestTrust` is added
- Struct `OuContainer` is added
- Struct `OuContainerClient` is added
- Struct `OuContainerCreateFuture` is added
- Struct `OuContainerDeleteFuture` is added
- Struct `OuContainerListResult` is added
- Struct `OuContainerListResultIterator` is added
- Struct `OuContainerListResultPage` is added
- Struct `OuContainerOperationsClient` is added
- Struct `OuContainerProperties` is added
- Struct `OuContainerUpdateFuture` is added
- Struct `ResourceForestSettings` is added
- Field `SyncKerberosPasswords` is added to struct `DomainSecuritySettings`
- Field `SyncOnPremPasswords` is added to struct `DomainSecuritySettings`
- Field `Version` is added to struct `DomainServiceProperties`
- Field `ResourceForestSettings` is added to struct `DomainServiceProperties`
- Field `DomainConfigurationType` is added to struct `DomainServiceProperties`
- Field `Sku` is added to struct `DomainServiceProperties`
- Field `DeploymentID` is added to struct `DomainServiceProperties`

