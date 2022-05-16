# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. BuildServiceAgentPoolClient.UpdatePut
	- Params
		- From: context.Context, string, string, string, string, BuildServiceAgentPoolSizeProperties
		- To: context.Context, string, string, string, string, BuildServiceAgentPoolResource
1. BuildServiceAgentPoolClient.UpdatePutPreparer
	- Params
		- From: context.Context, string, string, string, string, BuildServiceAgentPoolSizeProperties
		- To: context.Context, string, string, string, string, BuildServiceAgentPoolResource

## Additive Changes

### New Constants

1. ActionType.ActionTypeInternal
1. ProvisioningState.ProvisioningStateStarting
1. ProvisioningState.ProvisioningStateStopping

### New Funcs

1. OperationDetail.MarshalJSON() ([]byte, error)
1. PossibleActionTypeValues() []ActionType

### Struct Changes

#### New Struct Fields

1. OperationDetail.ActionType
