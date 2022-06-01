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

### Struct Changes

#### New Struct Fields

1. Principal.DisplayName
1. Principal.Email
