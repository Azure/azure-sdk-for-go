# Unreleased

## Additive Changes

### New Funcs

1. BaseClient.LocationGet(context.Context, string) (LocationGetResult, error)
1. BaseClient.LocationGetPreparer(context.Context, string) (*http.Request, error)
1. BaseClient.LocationGetResponder(*http.Response) (LocationGetResult, error)
1. BaseClient.LocationGetSender(*http.Request) (*http.Response, error)
1. BaseClient.LocationList(context.Context) (LocationListResult, error)
1. BaseClient.LocationListPreparer(context.Context) (*http.Request, error)
1. BaseClient.LocationListResponder(*http.Response) (LocationListResult, error)
1. BaseClient.LocationListSender(*http.Request) (*http.Response, error)
1. LocationGetResult.MarshalJSON() ([]byte, error)
1. LocationListResult.MarshalJSON() ([]byte, error)
1. LocationProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. LocationGetResult
1. LocationListResult
1. LocationProperties
