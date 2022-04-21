# Unreleased

## Breaking Changes

### Removed Constants

1. StorageAccountType.StorageAccountTypeZRS

### Signature Changes

#### Struct Fields

1. SQLPoolResourceProperties.CreateMode changed type from *string to CreateMode

## Additive Changes

### New Constants

1. CreateMode.CreateModeDefault
1. CreateMode.CreateModePointInTimeRestore
1. CreateMode.CreateModeRecovery
1. CreateMode.CreateModeRestore

### New Funcs

1. *DedicatedSQLminimalTLSSettings.UnmarshalJSON([]byte) error
1. *DedicatedSQLminimalTLSSettingsListResultIterator.Next() error
1. *DedicatedSQLminimalTLSSettingsListResultIterator.NextWithContext(context.Context) error
1. *DedicatedSQLminimalTLSSettingsListResultPage.Next() error
1. *DedicatedSQLminimalTLSSettingsListResultPage.NextWithContext(context.Context) error
1. *WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsUpdateFuture.UnmarshalJSON([]byte) error
1. DedicatedSQLminimalTLSSettings.MarshalJSON() ([]byte, error)
1. DedicatedSQLminimalTLSSettingsListResult.IsEmpty() bool
1. DedicatedSQLminimalTLSSettingsListResult.MarshalJSON() ([]byte, error)
1. DedicatedSQLminimalTLSSettingsListResultIterator.NotDone() bool
1. DedicatedSQLminimalTLSSettingsListResultIterator.Response() DedicatedSQLminimalTLSSettingsListResult
1. DedicatedSQLminimalTLSSettingsListResultIterator.Value() DedicatedSQLminimalTLSSettings
1. DedicatedSQLminimalTLSSettingsListResultPage.NotDone() bool
1. DedicatedSQLminimalTLSSettingsListResultPage.Response() DedicatedSQLminimalTLSSettingsListResult
1. DedicatedSQLminimalTLSSettingsListResultPage.Values() []DedicatedSQLminimalTLSSettings
1. MetadataSyncConfigProperties.MarshalJSON() ([]byte, error)
1. NewDedicatedSQLminimalTLSSettingsListResultIterator(DedicatedSQLminimalTLSSettingsListResultPage) DedicatedSQLminimalTLSSettingsListResultIterator
1. NewDedicatedSQLminimalTLSSettingsListResultPage(DedicatedSQLminimalTLSSettingsListResult, func(context.Context, DedicatedSQLminimalTLSSettingsListResult) (DedicatedSQLminimalTLSSettingsListResult, error)) DedicatedSQLminimalTLSSettingsListResultPage
1. NewWorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient(string) WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient
1. NewWorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClientWithBaseURI(string, string) WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient
1. PossibleCreateModeValues() []CreateMode
1. SQLPoolResourceProperties.MarshalJSON() ([]byte, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.Get(context.Context, string, string, string) (DedicatedSQLminimalTLSSettings, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.GetResponder(*http.Response) (DedicatedSQLminimalTLSSettings, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.List(context.Context, string, string) (DedicatedSQLminimalTLSSettingsListResultPage, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.ListComplete(context.Context, string, string) (DedicatedSQLminimalTLSSettingsListResultIterator, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.ListResponder(*http.Response) (DedicatedSQLminimalTLSSettingsListResult, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.ListSender(*http.Request) (*http.Response, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.Update(context.Context, string, string, DedicatedSQLminimalTLSSettings) (WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsUpdateFuture, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.UpdatePreparer(context.Context, string, string, DedicatedSQLminimalTLSSettings) (*http.Request, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.UpdateResponder(*http.Response) (DedicatedSQLminimalTLSSettings, error)
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient.UpdateSender(*http.Request) (WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsUpdateFuture, error)

### Struct Changes

#### New Structs

1. DedicatedSQLminimalTLSSettings
1. DedicatedSQLminimalTLSSettingsListResult
1. DedicatedSQLminimalTLSSettingsListResultIterator
1. DedicatedSQLminimalTLSSettingsListResultPage
1. DedicatedSQLminimalTLSSettingsProperties
1. ManagedIntegrationRuntimeManagedVirtualNetworkReference
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient
1. WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsUpdateFuture

#### New Struct Fields

1. DynamicExecutorAllocation.MaxExecutors
1. DynamicExecutorAllocation.MinExecutors
1. IntegrationRuntimeDataFlowProperties.Cleanup
1. ManagedIntegrationRuntime.*ManagedIntegrationRuntimeManagedVirtualNetworkReference
1. SelfHostedIntegrationRuntimeStatusTypeProperties.NewerVersions
1. SelfHostedIntegrationRuntimeStatusTypeProperties.ServiceRegion
1. WorkspaceProperties.TrustedServiceBypassEnabled
