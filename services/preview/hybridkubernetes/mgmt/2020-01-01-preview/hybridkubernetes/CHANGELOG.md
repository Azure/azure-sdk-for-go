
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewConnectedClusterListPage` signature has been changed from `(func(context.Context, ConnectedClusterList) (ConnectedClusterList, error))` to `(ConnectedClusterList,func(context.Context, ConnectedClusterList) (ConnectedClusterList, error))`
- Function `NewOperationListPage` signature has been changed from `(func(context.Context, OperationList) (OperationList, error))` to `(OperationList,func(context.Context, OperationList) (OperationList, error))`
- Type of `ErrorResponse.Error` has been changed from `*ErrorDetails` to `*ErrorDetail`
- Type of `AuthenticationDetails.AuthenticationMethod` has been changed from `AuthenticationMethod` to `*string`
- Const `ClientCertificate` has been removed
- Const `Token` has been removed
- Function `PossibleAuthenticationMethodValues` has been removed
- Struct `AuthenticationCertificateDetails` has been removed
- Struct `ErrorDetails` has been removed
- Field `ClientCertificate` of struct `AuthenticationDetailsValue` has been removed

## New Content

- Struct `ErrorAdditionalInfo` is added
- Struct `ErrorDetail` is added

