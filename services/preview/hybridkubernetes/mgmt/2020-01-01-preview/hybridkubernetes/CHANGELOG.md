Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewConnectedClusterListPage` parameter(s) have been changed from `(func(context.Context, ConnectedClusterList) (ConnectedClusterList, error))` to `(ConnectedClusterList, func(context.Context, ConnectedClusterList) (ConnectedClusterList, error))`
- Function `ConnectedClusterClient.ListClusterUserCredentials` parameter(s) have been changed from `(context.Context, string, string, *AuthenticationDetails)` to `(context.Context, string, string, *bool, *AuthenticationDetails)`
- Function `NewOperationListPage` parameter(s) have been changed from `(func(context.Context, OperationList) (OperationList, error))` to `(OperationList, func(context.Context, OperationList) (OperationList, error))`
- Function `ConnectedClusterClient.ListClusterUserCredentialsPreparer` parameter(s) have been changed from `(context.Context, string, string, *AuthenticationDetails)` to `(context.Context, string, string, *bool, *AuthenticationDetails)`
- Type of `AuthenticationDetails.AuthenticationMethod` has been changed from `AuthenticationMethod` to `*string`
- Type of `ErrorResponse.Error` has been changed from `*ErrorDetails` to `*ErrorDetail`
- Const `ClientCertificate` has been removed
- Const `Token` has been removed
- Function `PossibleAuthenticationMethodValues` has been removed
- Struct `AuthenticationCertificateDetails` has been removed
- Struct `ErrorDetails` has been removed
- Field `ClientCertificate` of struct `AuthenticationDetailsValue` has been removed

## New Content

- New const `Offline`
- New const `Connecting`
- New const `Expired`
- New const `Connected`
- New function `PossibleConnectivityStatusValues() []ConnectivityStatus`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `HybridConnectionConfig`
- New field `HybridConnectionConfig` in struct `CredentialResults`
- New field `ConnectivityStatus` in struct `ConnectedClusterProperties`
- New field `Distribution` in struct `ConnectedClusterProperties`
- New field `LastConnectivityTime` in struct `ConnectedClusterProperties`
- New field `TotalCoreCount` in struct `ConnectedClusterProperties`
- New field `Infrastructure` in struct `ConnectedClusterProperties`
- New field `Offering` in struct `ConnectedClusterProperties`
- New field `ManagedIdentityCertificateExpirationTime` in struct `ConnectedClusterProperties`
