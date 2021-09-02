# Change History

## Additive Changes

### New Constants

1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User

### New Funcs

1. CheckNameAvailabilityResult.MarshalJSON() ([]byte, error)
1. NewWorkspaceClient(string) WorkspaceClient
1. NewWorkspaceClientWithBaseURI(string, string) WorkspaceClient
1. PossibleCreatedByTypeValues() []CreatedByType
1. WorkspaceClient.CheckNameAvailability(context.Context, string, CheckNameAvailabilityParameters) (CheckNameAvailabilityResult, error)
1. WorkspaceClient.CheckNameAvailabilityPreparer(context.Context, string, CheckNameAvailabilityParameters) (*http.Request, error)
1. WorkspaceClient.CheckNameAvailabilityResponder(*http.Response) (CheckNameAvailabilityResult, error)
1. WorkspaceClient.CheckNameAvailabilitySender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. CheckNameAvailabilityParameters
1. CheckNameAvailabilityResult
1. SystemData
1. WorkspaceClient

#### New Struct Fields

1. SkuDescription.RestrictedAccessURI
1. SkuDescription.Version
1. Workspace.SystemData
