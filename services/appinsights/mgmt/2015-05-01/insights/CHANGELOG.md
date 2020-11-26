
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewWebTestListResultPage` signature has been changed from `(func(context.Context, WebTestListResult) (WebTestListResult, error))` to `(WebTestListResult,func(context.Context, WebTestListResult) (WebTestListResult, error))`
- Function `NewApplicationInsightsComponentListResultPage` signature has been changed from `(func(context.Context, ApplicationInsightsComponentListResult) (ApplicationInsightsComponentListResult, error))` to `(ApplicationInsightsComponentListResult,func(context.Context, ApplicationInsightsComponentListResult) (ApplicationInsightsComponentListResult, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`

## New Content

- Function `MyWorkbooksClient.ListByResourceGroupResponder(*http.Response) (MyWorkbooksListResult,error)` is added
- Function `MyWorkbookProperties.MarshalJSON() ([]byte,error)` is added
- Function `MyWorkbooksClient.DeleteSender(*http.Request) (*http.Response,error)` is added
- Function `MyWorkbooksClient.UpdatePreparer(context.Context,string,string,MyWorkbook) (*http.Request,error)` is added
- Function `MyWorkbooksClient.Delete(context.Context,string,string) (autorest.Response,error)` is added
- Function `MyWorkbooksClient.ListBySubscriptionResponder(*http.Response) (MyWorkbooksListResult,error)` is added
- Function `MyWorkbooksClient.GetPreparer(context.Context,string,string) (*http.Request,error)` is added
- Function `MyWorkbooksClient.ListBySubscriptionPreparer(context.Context,CategoryType,[]string,*bool) (*http.Request,error)` is added
- Function `*MyWorkbook.UnmarshalJSON([]byte) error` is added
- Function `MyWorkbooksClient.DeleteResponder(*http.Response) (autorest.Response,error)` is added
- Function `MyWorkbooksClient.GetResponder(*http.Response) (MyWorkbook,error)` is added
- Function `MyWorkbook.MarshalJSON() ([]byte,error)` is added
- Function `MyWorkbooksClient.UpdateSender(*http.Request) (*http.Response,error)` is added
- Function `MyWorkbooksClient.CreateOrUpdatePreparer(context.Context,string,string,MyWorkbook) (*http.Request,error)` is added
- Function `MyWorkbooksClient.ListByResourceGroupPreparer(context.Context,string,CategoryType,[]string,*bool) (*http.Request,error)` is added
- Function `MyWorkbooksClient.ListBySubscription(context.Context,CategoryType,[]string,*bool) (MyWorkbooksListResult,error)` is added
- Function `MyWorkbooksClient.CreateOrUpdate(context.Context,string,string,MyWorkbook) (MyWorkbook,error)` is added
- Function `MyWorkbooksClient.CreateOrUpdateSender(*http.Request) (*http.Response,error)` is added
- Function `MyWorkbooksClient.Get(context.Context,string,string) (MyWorkbook,error)` is added
- Function `MyWorkbookResource.MarshalJSON() ([]byte,error)` is added
- Function `MyWorkbooksClient.GetSender(*http.Request) (*http.Response,error)` is added
- Function `MyWorkbooksClient.CreateOrUpdateResponder(*http.Response) (MyWorkbook,error)` is added
- Function `MyWorkbooksClient.ListByResourceGroupSender(*http.Request) (*http.Response,error)` is added
- Function `MyWorkbooksClient.UpdateResponder(*http.Response) (MyWorkbook,error)` is added
- Function `NewMyWorkbooksClientWithBaseURI(string,string) MyWorkbooksClient` is added
- Function `MyWorkbooksClient.Update(context.Context,string,string,MyWorkbook) (MyWorkbook,error)` is added
- Function `NewMyWorkbooksClient(string) MyWorkbooksClient` is added
- Function `MyWorkbooksClient.ListBySubscriptionSender(*http.Request) (*http.Response,error)` is added
- Function `MyWorkbooksClient.ListByResourceGroup(context.Context,string,CategoryType,[]string,*bool) (MyWorkbooksListResult,error)` is added
- Function `MyWorkbooksClient.DeletePreparer(context.Context,string,string) (*http.Request,error)` is added
- Struct `MyWorkbook` is added
- Struct `MyWorkbookError` is added
- Struct `MyWorkbookProperties` is added
- Struct `MyWorkbookResource` is added
- Struct `MyWorkbooksClient` is added
- Struct `MyWorkbooksListResult` is added

