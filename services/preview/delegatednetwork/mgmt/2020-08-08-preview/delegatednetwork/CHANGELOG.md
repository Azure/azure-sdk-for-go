# Unreleased

## Breaking Changes

### Removed Funcs

1. ErrorDefinition.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. ControllerDetailsModel
1. ErrorDefinition

#### Removed Struct Fields

1. Operation.Properties

### Signature Changes

#### Funcs

1. DelegatedSubnetServiceClient.DeleteDetails
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, *bool
1. DelegatedSubnetServiceClient.DeleteDetailsPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, *bool

#### Struct Fields

1. DelegatedSubnetProperties.ControllerDetails changed type from *ControllerDetailsModel to *ControllerDetails
1. ErrorResponse.Error changed type from *ErrorDefinition to *ErrorDetail
1. Operation.Origin changed type from *string to Origin

## Additive Changes

### New Constants

1. ActionType.Internal
1. Origin.System
1. Origin.User
1. Origin.Usersystem

### New Funcs

1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)
1. PossibleActionTypeValues() []ActionType
1. PossibleOriginValues() []Origin

### Struct Changes

#### New Structs

1. ErrorAdditionalInfo
1. ErrorDetail

#### New Struct Fields

1. Operation.ActionType
