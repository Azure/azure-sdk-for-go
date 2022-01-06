# Unreleased

## Breaking Changes

### Removed Constants

1. Code.BadRequest
1. Code.Conflict
1. Code.NotFound

### Removed Funcs

1. PossibleCodeValues() []Code

### Signature Changes

#### Struct Fields

1. ExtendedErrorInfo.Code changed type from Code to *string

## Additive Changes

### Struct Changes

#### New Struct Fields

1. Error.Code
1. Error.Message
