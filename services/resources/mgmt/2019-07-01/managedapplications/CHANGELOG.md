# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. ErrorResponse.ErrorCode
1. ErrorResponse.ErrorMessage
1. ErrorResponse.HTTPStatus

## Additive Changes

### New Funcs

1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)
1. OperationListResult.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ErrorAdditionalInfo
1. ErrorDetail

#### New Struct Fields

1. ApplicationPackageLockingPolicyDefinition.AllowedDataActions
1. ErrorResponse.Error
