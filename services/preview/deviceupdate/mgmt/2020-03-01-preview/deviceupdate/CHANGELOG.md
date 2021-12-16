# Change History

## Breaking Changes

### Removed Constants

1. Origin.System
1. Origin.Usersystem
1. ProvisioningState.Accepted
1. ProvisioningState.Creating
1. ProvisioningState.Deleted

### Removed Funcs

1. ErrorDefinition.MarshalJSON() ([]byte, error)
1. ErrorResponse.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. ErrorDefinition

#### Removed Struct Fields

1. ErrorResponse.AdditionalInfo
1. ErrorResponse.Code
1. ErrorResponse.Details
1. ErrorResponse.Message
1. ErrorResponse.Target

### Signature Changes

#### Const Types

1. Canceled changed type from ProvisioningState to GroupIDProvisioningState
1. Failed changed type from ProvisioningState to GroupIDProvisioningState
1. Succeeded changed type from ProvisioningState to GroupIDProvisioningState
1. User changed type from Origin to CreatedByType

## Additive Changes

### New Constants

1. CheckNameAvailabilityReason.AlreadyExists
1. CheckNameAvailabilityReason.Invalid
1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. ManagedServiceIdentityType.None
1. ManagedServiceIdentityType.SystemAssigned
1. ManagedServiceIdentityType.SystemAssignedUserAssigned
1. ManagedServiceIdentityType.UserAssigned
1. Origin.OriginSystem
1. Origin.OriginUser
1. Origin.OriginUsersystem
1. PrivateEndpointConnectionProvisioningState.PrivateEndpointConnectionProvisioningStateCreating
1. PrivateEndpointConnectionProvisioningState.PrivateEndpointConnectionProvisioningStateDeleting
1. PrivateEndpointConnectionProvisioningState.PrivateEndpointConnectionProvisioningStateFailed
1. PrivateEndpointConnectionProvisioningState.PrivateEndpointConnectionProvisioningStateSucceeded
1. PrivateEndpointServiceConnectionStatus.Approved
1. PrivateEndpointServiceConnectionStatus.Pending
1. PrivateEndpointServiceConnectionStatus.Rejected
1. ProvisioningState.ProvisioningStateAccepted
1. ProvisioningState.ProvisioningStateCanceled
1. ProvisioningState.ProvisioningStateCreating
1. ProvisioningState.ProvisioningStateDeleted
1. ProvisioningState.ProvisioningStateFailed
1. ProvisioningState.ProvisioningStateSucceeded
1. PublicNetworkAccess.Disabled
1. PublicNetworkAccess.Enabled

### New Funcs

1. *GroupInformation.UnmarshalJSON([]byte) error
1. *PrivateEndpointConnection.UnmarshalJSON([]byte) error
1. *PrivateEndpointConnectionsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *PrivateEndpointConnectionsDeleteFuture.UnmarshalJSON([]byte) error
1. AccountsClient.Head(context.Context, string, string) (autorest.Response, error)
1. AccountsClient.HeadPreparer(context.Context, string, string) (*http.Request, error)
1. AccountsClient.HeadResponder(*http.Response) (autorest.Response, error)
1. AccountsClient.HeadSender(*http.Request) (*http.Response, error)
1. BaseClient.CheckNameAvailability(context.Context, CheckNameAvailabilityRequest) (CheckNameAvailabilityResponse, error)
1. BaseClient.CheckNameAvailabilityPreparer(context.Context, CheckNameAvailabilityRequest) (*http.Request, error)
1. BaseClient.CheckNameAvailabilityResponder(*http.Response) (CheckNameAvailabilityResponse, error)
1. BaseClient.CheckNameAvailabilitySender(*http.Request) (*http.Response, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)
1. GroupInformation.MarshalJSON() ([]byte, error)
1. GroupInformationProperties.MarshalJSON() ([]byte, error)
1. InstancesClient.Head(context.Context, string, string, string) (autorest.Response, error)
1. InstancesClient.HeadPreparer(context.Context, string, string, string) (*http.Request, error)
1. InstancesClient.HeadResponder(*http.Response) (autorest.Response, error)
1. InstancesClient.HeadSender(*http.Request) (*http.Response, error)
1. ManagedServiceIdentity.MarshalJSON() ([]byte, error)
1. NewPrivateEndpointConnectionsClient(string) PrivateEndpointConnectionsClient
1. NewPrivateEndpointConnectionsClientWithBaseURI(string, string) PrivateEndpointConnectionsClient
1. NewPrivateLinkResourcesClient(string) PrivateLinkResourcesClient
1. NewPrivateLinkResourcesClientWithBaseURI(string, string) PrivateLinkResourcesClient
1. PossibleCheckNameAvailabilityReasonValues() []CheckNameAvailabilityReason
1. PossibleCreatedByTypeValues() []CreatedByType
1. PossibleGroupIDProvisioningStateValues() []GroupIDProvisioningState
1. PossibleManagedServiceIdentityTypeValues() []ManagedServiceIdentityType
1. PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState
1. PossiblePrivateEndpointServiceConnectionStatusValues() []PrivateEndpointServiceConnectionStatus
1. PossiblePublicNetworkAccessValues() []PublicNetworkAccess
1. PrivateEndpoint.MarshalJSON() ([]byte, error)
1. PrivateEndpointConnection.MarshalJSON() ([]byte, error)
1. PrivateEndpointConnectionsClient.CreateOrUpdate(context.Context, string, string, string, PrivateEndpointConnection) (PrivateEndpointConnectionsCreateOrUpdateFuture, error)
1. PrivateEndpointConnectionsClient.CreateOrUpdatePreparer(context.Context, string, string, string, PrivateEndpointConnection) (*http.Request, error)
1. PrivateEndpointConnectionsClient.CreateOrUpdateResponder(*http.Response) (PrivateEndpointConnection, error)
1. PrivateEndpointConnectionsClient.CreateOrUpdateSender(*http.Request) (PrivateEndpointConnectionsCreateOrUpdateFuture, error)
1. PrivateEndpointConnectionsClient.Delete(context.Context, string, string, string) (PrivateEndpointConnectionsDeleteFuture, error)
1. PrivateEndpointConnectionsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateEndpointConnectionsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. PrivateEndpointConnectionsClient.DeleteSender(*http.Request) (PrivateEndpointConnectionsDeleteFuture, error)
1. PrivateEndpointConnectionsClient.Get(context.Context, string, string, string) (PrivateEndpointConnection, error)
1. PrivateEndpointConnectionsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateEndpointConnectionsClient.GetResponder(*http.Response) (PrivateEndpointConnection, error)
1. PrivateEndpointConnectionsClient.GetSender(*http.Request) (*http.Response, error)
1. PrivateEndpointConnectionsClient.ListByAccount(context.Context, string, string) (PrivateEndpointConnectionListResult, error)
1. PrivateEndpointConnectionsClient.ListByAccountPreparer(context.Context, string, string) (*http.Request, error)
1. PrivateEndpointConnectionsClient.ListByAccountResponder(*http.Response) (PrivateEndpointConnectionListResult, error)
1. PrivateEndpointConnectionsClient.ListByAccountSender(*http.Request) (*http.Response, error)
1. PrivateLinkResourceProperties.MarshalJSON() ([]byte, error)
1. PrivateLinkResourcesClient.Get(context.Context, string, string, string) (GroupInformation, error)
1. PrivateLinkResourcesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateLinkResourcesClient.GetResponder(*http.Response) (GroupInformation, error)
1. PrivateLinkResourcesClient.GetSender(*http.Request) (*http.Response, error)
1. PrivateLinkResourcesClient.ListByAccount(context.Context, string, string) (PrivateLinkResourceListResult, error)
1. PrivateLinkResourcesClient.ListByAccountPreparer(context.Context, string, string) (*http.Request, error)
1. PrivateLinkResourcesClient.ListByAccountResponder(*http.Response) (PrivateLinkResourceListResult, error)
1. PrivateLinkResourcesClient.ListByAccountSender(*http.Request) (*http.Response, error)
1. UserAssignedIdentity.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. CheckNameAvailabilityRequest
1. CheckNameAvailabilityResponse
1. DiagnosticStorageProperties
1. ErrorDetail
1. GroupInformation
1. GroupInformationProperties
1. ManagedServiceIdentity
1. PrivateEndpoint
1. PrivateEndpointConnection
1. PrivateEndpointConnectionListResult
1. PrivateEndpointConnectionProperties
1. PrivateEndpointConnectionsClient
1. PrivateEndpointConnectionsCreateOrUpdateFuture
1. PrivateEndpointConnectionsDeleteFuture
1. PrivateLinkResourceListResult
1. PrivateLinkResourceProperties
1. PrivateLinkResourcesClient
1. PrivateLinkServiceConnectionState
1. SystemData
1. UserAssignedIdentity

#### New Struct Fields

1. Account.Identity
1. Account.SystemData
1. AccountProperties.PublicNetworkAccess
1. AccountUpdate.Identity
1. AzureEntityResource.SystemData
1. ErrorResponse.Error
1. Instance.SystemData
1. InstanceProperties.DiagnosticStorageProperties
1. InstanceProperties.EnableDiagnostics
1. ProxyResource.SystemData
1. Resource.SystemData
1. TrackedResource.SystemData
