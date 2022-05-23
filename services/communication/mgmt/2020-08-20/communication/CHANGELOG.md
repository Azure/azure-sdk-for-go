# Unreleased

## Breaking Changes

### Removed Constants

1. Status.StatusCanceled
1. Status.StatusCreating
1. Status.StatusDeleting
1. Status.StatusFailed
1. Status.StatusMoving
1. Status.StatusSucceeded

### Removed Funcs

1. *OperationStatus.UnmarshalJSON([]byte) error
1. NewOperationStatusesClient(string) OperationStatusesClient
1. NewOperationStatusesClientWithBaseURI(string, string) OperationStatusesClient
1. OperationStatus.MarshalJSON() ([]byte, error)
1. OperationStatusesClient.Get(context.Context, string, string) (OperationStatus, error)
1. OperationStatusesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. OperationStatusesClient.GetResponder(*http.Response) (OperationStatus, error)
1. OperationStatusesClient.GetSender(*http.Request) (*http.Response, error)
1. PossibleStatusValues() []Status

### Struct Changes

#### Removed Structs

1. OperationStatus
1. OperationStatusesClient
