# Unreleased

## Breaking Changes

### Removed Funcs

1. *SapMonitor.UnmarshalJSON([]byte) error
1. *SapMonitorListResultIterator.Next() error
1. *SapMonitorListResultIterator.NextWithContext(context.Context) error
1. *SapMonitorListResultPage.Next() error
1. *SapMonitorListResultPage.NextWithContext(context.Context) error
1. *SapMonitorsCreateFuture.UnmarshalJSON([]byte) error
1. *SapMonitorsDeleteFuture.UnmarshalJSON([]byte) error
1. HanaInstanceProperties.MarshalJSON() ([]byte, error)
1. HardwareProfile.MarshalJSON() ([]byte, error)
1. NetworkProfile.MarshalJSON() ([]byte, error)
1. NewSapMonitorListResultIterator(SapMonitorListResultPage) SapMonitorListResultIterator
1. NewSapMonitorListResultPage(SapMonitorListResult, func(context.Context, SapMonitorListResult) (SapMonitorListResult, error)) SapMonitorListResultPage
1. NewSapMonitorsClient(string) SapMonitorsClient
1. NewSapMonitorsClientWithBaseURI(string, string) SapMonitorsClient
1. OSProfile.MarshalJSON() ([]byte, error)
1. SapMonitor.MarshalJSON() ([]byte, error)
1. SapMonitorListResult.IsEmpty() bool
1. SapMonitorListResultIterator.NotDone() bool
1. SapMonitorListResultIterator.Response() SapMonitorListResult
1. SapMonitorListResultIterator.Value() SapMonitor
1. SapMonitorListResultPage.NotDone() bool
1. SapMonitorListResultPage.Response() SapMonitorListResult
1. SapMonitorListResultPage.Values() []SapMonitor
1. SapMonitorProperties.MarshalJSON() ([]byte, error)
1. SapMonitorsClient.Create(context.Context, string, string, SapMonitor) (SapMonitorsCreateFuture, error)
1. SapMonitorsClient.CreatePreparer(context.Context, string, string, SapMonitor) (*http.Request, error)
1. SapMonitorsClient.CreateResponder(*http.Response) (SapMonitor, error)
1. SapMonitorsClient.CreateSender(*http.Request) (SapMonitorsCreateFuture, error)
1. SapMonitorsClient.Delete(context.Context, string, string) (SapMonitorsDeleteFuture, error)
1. SapMonitorsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. SapMonitorsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. SapMonitorsClient.DeleteSender(*http.Request) (SapMonitorsDeleteFuture, error)
1. SapMonitorsClient.Get(context.Context, string, string) (SapMonitor, error)
1. SapMonitorsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. SapMonitorsClient.GetResponder(*http.Response) (SapMonitor, error)
1. SapMonitorsClient.GetSender(*http.Request) (*http.Response, error)
1. SapMonitorsClient.List(context.Context) (SapMonitorListResultPage, error)
1. SapMonitorsClient.ListComplete(context.Context) (SapMonitorListResultIterator, error)
1. SapMonitorsClient.ListPreparer(context.Context) (*http.Request, error)
1. SapMonitorsClient.ListResponder(*http.Response) (SapMonitorListResult, error)
1. SapMonitorsClient.ListSender(*http.Request) (*http.Response, error)
1. SapMonitorsClient.Update(context.Context, string, string, Tags) (SapMonitor, error)
1. SapMonitorsClient.UpdatePreparer(context.Context, string, string, Tags) (*http.Request, error)
1. SapMonitorsClient.UpdateResponder(*http.Response) (SapMonitor, error)
1. SapMonitorsClient.UpdateSender(*http.Request) (*http.Response, error)
1. StorageProfile.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. SapMonitor
1. SapMonitorListResult
1. SapMonitorListResultIterator
1. SapMonitorListResultPage
1. SapMonitorProperties
1. SapMonitorsClient
1. SapMonitorsCreateFuture
1. SapMonitorsDeleteFuture

#### Removed Struct Fields

1. ErrorResponse.Code
1. ErrorResponse.Message

## Additive Changes

### New Funcs

1. ErrorResponseError.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ErrorResponseError
1. SAPSystemID

#### New Struct Fields

1. ErrorResponse.Error
1. Operation.IsDataAction
1. StorageProfile.HanaSids
