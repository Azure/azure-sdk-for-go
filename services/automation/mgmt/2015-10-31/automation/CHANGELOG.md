# Unreleased

## Additive Changes

### New Funcs

1. *PrivateEndpointConnection.UnmarshalJSON([]byte) error
1. *Watcher.UnmarshalJSON([]byte) error
1. *WatcherListResultIterator.Next() error
1. *WatcherListResultIterator.NextWithContext(context.Context) error
1. *WatcherListResultPage.Next() error
1. *WatcherListResultPage.NextWithContext(context.Context) error
1. *WatcherUpdateParameters.UnmarshalJSON([]byte) error
1. NewWatcherClient(string) WatcherClient
1. NewWatcherClientWithBaseURI(string, string) WatcherClient
1. NewWatcherListResultIterator(WatcherListResultPage) WatcherListResultIterator
1. NewWatcherListResultPage(WatcherListResult, func(context.Context, WatcherListResult) (WatcherListResult, error)) WatcherListResultPage
1. PrivateEndpointConnection.MarshalJSON() ([]byte, error)
1. PrivateLinkServiceConnectionStateProperty.MarshalJSON() ([]byte, error)
1. Watcher.MarshalJSON() ([]byte, error)
1. WatcherClient.CreateOrUpdate(context.Context, string, string, string, Watcher) (Watcher, error)
1. WatcherClient.CreateOrUpdatePreparer(context.Context, string, string, string, Watcher) (*http.Request, error)
1. WatcherClient.CreateOrUpdateResponder(*http.Response) (Watcher, error)
1. WatcherClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. WatcherClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. WatcherClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. WatcherClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. WatcherClient.DeleteSender(*http.Request) (*http.Response, error)
1. WatcherClient.Get(context.Context, string, string, string) (Watcher, error)
1. WatcherClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. WatcherClient.GetResponder(*http.Response) (Watcher, error)
1. WatcherClient.GetSender(*http.Request) (*http.Response, error)
1. WatcherClient.ListByAutomationAccount(context.Context, string, string, string) (WatcherListResultPage, error)
1. WatcherClient.ListByAutomationAccountComplete(context.Context, string, string, string) (WatcherListResultIterator, error)
1. WatcherClient.ListByAutomationAccountPreparer(context.Context, string, string, string) (*http.Request, error)
1. WatcherClient.ListByAutomationAccountResponder(*http.Response) (WatcherListResult, error)
1. WatcherClient.ListByAutomationAccountSender(*http.Request) (*http.Response, error)
1. WatcherClient.Start(context.Context, string, string, string) (autorest.Response, error)
1. WatcherClient.StartPreparer(context.Context, string, string, string) (*http.Request, error)
1. WatcherClient.StartResponder(*http.Response) (autorest.Response, error)
1. WatcherClient.StartSender(*http.Request) (*http.Response, error)
1. WatcherClient.Stop(context.Context, string, string, string) (autorest.Response, error)
1. WatcherClient.StopPreparer(context.Context, string, string, string) (*http.Request, error)
1. WatcherClient.StopResponder(*http.Response) (autorest.Response, error)
1. WatcherClient.StopSender(*http.Request) (*http.Response, error)
1. WatcherClient.Update(context.Context, string, string, string, WatcherUpdateParameters) (Watcher, error)
1. WatcherClient.UpdatePreparer(context.Context, string, string, string, WatcherUpdateParameters) (*http.Request, error)
1. WatcherClient.UpdateResponder(*http.Response) (Watcher, error)
1. WatcherClient.UpdateSender(*http.Request) (*http.Response, error)
1. WatcherListResult.IsEmpty() bool
1. WatcherListResultIterator.NotDone() bool
1. WatcherListResultIterator.Response() WatcherListResult
1. WatcherListResultIterator.Value() Watcher
1. WatcherListResultPage.NotDone() bool
1. WatcherListResultPage.Response() WatcherListResult
1. WatcherListResultPage.Values() []Watcher
1. WatcherProperties.MarshalJSON() ([]byte, error)
1. WatcherUpdateParameters.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. PrivateEndpointConnection
1. PrivateEndpointConnectionProperties
1. PrivateEndpointProperty
1. PrivateLinkServiceConnectionStateProperty
1. Watcher
1. WatcherClient
1. WatcherListResult
1. WatcherListResultIterator
1. WatcherListResultPage
1. WatcherProperties
1. WatcherUpdateParameters
1. WatcherUpdateProperties
