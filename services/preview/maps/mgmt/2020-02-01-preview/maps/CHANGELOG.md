Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Field `Code` of struct `ErrorResponse` has been removed
- Field `Message` of struct `ErrorResponse` has been removed
- Field `Target` of struct `ErrorResponse` has been removed
- Field `Details` of struct `ErrorResponse` has been removed
- Field `AdditionalInfo` of struct `ErrorResponse` has been removed

## New Content

- New function `CreatorsClient.DeleteResponder(*http.Response) (autorest.Response, error)`
- New function `CreatorsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)`
- New function `CreatorsClient.ListByAccountSender(*http.Request) (*http.Response, error)`
- New function `CreatorsClient.ListByAccountResponder(*http.Response) (CreatorList, error)`
- New function `NewCreatorsClientWithBaseURI(string, string) CreatorsClient`
- New function `CreatorCreateParameters.MarshalJSON() ([]byte, error)`
- New function `CreatorsClient.CreateOrUpdate(context.Context, string, string, string, CreatorCreateParameters) (Creator, error)`
- New function `Creator.MarshalJSON() ([]byte, error)`
- New function `CreatorsClient.ListByAccount(context.Context, string, string) (CreatorList, error)`
- New function `CreatorsClient.UpdateSender(*http.Request) (*http.Response, error)`
- New function `CreatorsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)`
- New function `CreatorsClient.UpdateResponder(*http.Response) (Creator, error)`
- New function `CreatorsClient.UpdatePreparer(context.Context, string, string, string, CreatorUpdateParameters) (*http.Request, error)`
- New function `CreatorsClient.CreateOrUpdatePreparer(context.Context, string, string, string, CreatorCreateParameters) (*http.Request, error)`
- New function `CreatorsClient.Update(context.Context, string, string, string, CreatorUpdateParameters) (Creator, error)`
- New function `CreatorsClient.Delete(context.Context, string, string, string) (autorest.Response, error)`
- New function `CreatorsClient.CreateOrUpdateResponder(*http.Response) (Creator, error)`
- New function `CreatorUpdateParameters.MarshalJSON() ([]byte, error)`
- New function `CreatorsClient.GetResponder(*http.Response) (Creator, error)`
- New function `NewCreatorsClient(string) CreatorsClient`
- New function `CreatorsClient.DeleteSender(*http.Request) (*http.Response, error)`
- New function `CreatorsClient.ListByAccountPreparer(context.Context, string, string) (*http.Request, error)`
- New function `CreatorsClient.GetSender(*http.Request) (*http.Response, error)`
- New function `CreatorsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)`
- New function `CreatorsClient.Get(context.Context, string, string, string) (Creator, error)`
- New struct `Creator`
- New struct `CreatorCreateParameters`
- New struct `CreatorList`
- New struct `CreatorProperties`
- New struct `CreatorUpdateParameters`
- New struct `CreatorsClient`
- New struct `ErrorDetail`
- New field `Error` in struct `ErrorResponse`
