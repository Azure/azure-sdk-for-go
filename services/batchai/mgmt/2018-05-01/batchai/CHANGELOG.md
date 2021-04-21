# Unreleased Content

## Breaking Changes

### Removed Funcs

1. *FileServersDeleteFuture.UnmarshalJSON([]byte) error
1. FileServersClient.Delete(context.Context, string, string, string) (FileServersDeleteFuture, error)
1. FileServersClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. FileServersClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. FileServersClient.DeleteSender(*http.Request) (FileServersDeleteFuture, error)
1. FileServersClient.Get(context.Context, string, string, string) (FileServer, error)
1. FileServersClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. FileServersClient.GetResponder(*http.Response) (FileServer, error)
1. FileServersClient.GetSender(*http.Request) (*http.Response, error)

### Struct Changes

#### Removed Structs

1. FileServersDeleteFuture
