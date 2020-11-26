
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Field `Details` of struct `ErrorResponse` has been removed
- Field `AdditionalInfo` of struct `ErrorResponse` has been removed
- Field `Code` of struct `ErrorResponse` has been removed
- Field `Message` of struct `ErrorResponse` has been removed
- Field `Target` of struct `ErrorResponse` has been removed

## New Content

- Function `CreatorsClient.UpdateResponder(*http.Response) (Creator,error)` is added
- Function `NewCreatorsClient(string) CreatorsClient` is added
- Function `CreatorCreateParameters.MarshalJSON() ([]byte,error)` is added
- Function `CreatorsClient.GetSender(*http.Request) (*http.Response,error)` is added
- Function `CreatorsClient.DeleteResponder(*http.Response) (autorest.Response,error)` is added
- Function `CreatorsClient.CreateOrUpdatePreparer(context.Context,string,string,string,CreatorCreateParameters) (*http.Request,error)` is added
- Function `CreatorsClient.ListByAccountSender(*http.Request) (*http.Response,error)` is added
- Function `CreatorsClient.Update(context.Context,string,string,string,CreatorUpdateParameters) (Creator,error)` is added
- Function `CreatorsClient.ListByAccountResponder(*http.Response) (CreatorList,error)` is added
- Function `NewCreatorsClientWithBaseURI(string,string) CreatorsClient` is added
- Function `CreatorsClient.Delete(context.Context,string,string,string) (autorest.Response,error)` is added
- Function `CreatorsClient.CreateOrUpdateSender(*http.Request) (*http.Response,error)` is added
- Function `Creator.MarshalJSON() ([]byte,error)` is added
- Function `CreatorsClient.Get(context.Context,string,string,string) (Creator,error)` is added
- Function `CreatorsClient.UpdateSender(*http.Request) (*http.Response,error)` is added
- Function `CreatorUpdateParameters.MarshalJSON() ([]byte,error)` is added
- Function `CreatorsClient.DeleteSender(*http.Request) (*http.Response,error)` is added
- Function `CreatorsClient.GetResponder(*http.Response) (Creator,error)` is added
- Function `CreatorsClient.UpdatePreparer(context.Context,string,string,string,CreatorUpdateParameters) (*http.Request,error)` is added
- Function `CreatorsClient.DeletePreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `CreatorsClient.CreateOrUpdate(context.Context,string,string,string,CreatorCreateParameters) (Creator,error)` is added
- Function `CreatorsClient.ListByAccountPreparer(context.Context,string,string) (*http.Request,error)` is added
- Function `CreatorsClient.CreateOrUpdateResponder(*http.Response) (Creator,error)` is added
- Function `CreatorsClient.ListByAccount(context.Context,string,string) (CreatorList,error)` is added
- Function `CreatorsClient.GetPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Struct `Creator` is added
- Struct `CreatorCreateParameters` is added
- Struct `CreatorList` is added
- Struct `CreatorProperties` is added
- Struct `CreatorUpdateParameters` is added
- Struct `CreatorsClient` is added
- Struct `ErrorDetail` is added
- Field `Error` is added to struct `ErrorResponse`

