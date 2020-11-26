
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewRegistrationDefinitionListPage` signature has been changed from `(func(context.Context, RegistrationDefinitionList) (RegistrationDefinitionList, error))` to `(RegistrationDefinitionList,func(context.Context, RegistrationDefinitionList) (RegistrationDefinitionList, error))`
- Function `NewRegistrationAssignmentListPage` signature has been changed from `(func(context.Context, RegistrationAssignmentList) (RegistrationAssignmentList, error))` to `(RegistrationAssignmentList,func(context.Context, RegistrationAssignmentList) (RegistrationAssignmentList, error))`
- Type of `ErrorResponse.Error` has been changed from `*ErrorResponseError` to `*ErrorDefinition`
- Struct `ErrorResponseError` has been removed

## New Content

- Struct `ErrorDefinition` is added
- Field `DelegatedRoleDefinitionIds` is added to struct `Authorization`
- Field `PrincipalIDDisplayName` is added to struct `Authorization`

