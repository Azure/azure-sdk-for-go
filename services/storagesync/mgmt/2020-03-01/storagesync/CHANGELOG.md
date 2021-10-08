# Unreleased

## Additive Changes

### New Constants

1. RegisteredServerAgentVersionStatus.Blocked
1. RegisteredServerAgentVersionStatus.Expired
1. RegisteredServerAgentVersionStatus.NearExpiry
1. RegisteredServerAgentVersionStatus.Ok
1. ServerEndpointSyncMode.InitialFullDownload
1. ServerEndpointSyncMode.InitialUpload
1. ServerEndpointSyncMode.NamespaceDownload
1. ServerEndpointSyncMode.Regular
1. ServerEndpointSyncMode.SnapshotUpload

### New Funcs

1. BaseClient.LocationOperationStatusMethod(context.Context, string, string) (LocationOperationStatus, error)
1. BaseClient.LocationOperationStatusMethodPreparer(context.Context, string, string) (*http.Request, error)
1. BaseClient.LocationOperationStatusMethodResponder(*http.Response) (LocationOperationStatus, error)
1. BaseClient.LocationOperationStatusMethodSender(*http.Request) (*http.Response, error)
1. LocationOperationStatus.MarshalJSON() ([]byte, error)
1. PossibleRegisteredServerAgentVersionStatusValues() []RegisteredServerAgentVersionStatus
1. PossibleServerEndpointSyncModeValues() []ServerEndpointSyncMode
1. RegisteredServerProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. InnerErrorDetails
1. LocationOperationStatus
1. OperationProperties
1. OperationResourceMetricSpecification
1. OperationResourceMetricSpecificationDimension
1. OperationResourceServiceSpecification

#### New Struct Fields

1. APIError.InnerError
1. ErrorDetails.ExceptionType
1. ErrorDetails.HTTPErrorCode
1. ErrorDetails.HTTPMethod
1. ErrorDetails.HashedMessage
1. ErrorDetails.RequestURI
1. OperationEntity.Properties
1. RegisteredServerProperties.AgentVersionExpirationDate
1. RegisteredServerProperties.AgentVersionStatus
1. ServerEndpointSyncActivityStatus.SyncMode
1. ServerEndpointSyncSessionStatus.LastSyncMode
