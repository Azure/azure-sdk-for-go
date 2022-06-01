# Unreleased

## Breaking Changes

### Removed Constants

1. PrincipalType.Application
1. PrincipalType.DirectoryObjectOrGroup
1. PrincipalType.DirectoryRoleTemplate
1. PrincipalType.Everyone
1. PrincipalType.MSI
1. PrincipalType.Unknown

### Removed Funcs

1. Principal.MarshalJSON() ([]byte, error)

## Additive Changes

### New Funcs

1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ErrorAdditionalInfo
1. ErrorDetail
1. ErrorResponse

#### New Struct Fields

1. Principal.DisplayName
1. Principal.Email
