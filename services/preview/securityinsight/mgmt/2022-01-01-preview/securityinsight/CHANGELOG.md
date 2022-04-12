# Unreleased

## Breaking Changes

### Signature Changes

#### Struct Fields

1. WatchlistItemProperties.EntityMapping changed type from interface{} to map[string]interface{}
1. WatchlistItemProperties.ItemsKeyValue changed type from interface{} to map[string]interface{}

## Additive Changes

### New Constants

1. ProvisioningState.ProvisioningStateCanceled
1. ProvisioningState.ProvisioningStateFailed
1. ProvisioningState.ProvisioningStateInProgress
1. ProvisioningState.ProvisioningStateSucceeded

### New Funcs

1. PossibleProvisioningStateValues() []ProvisioningState
1. WatchlistItemProperties.MarshalJSON() ([]byte, error)
1. WatchlistProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Struct Fields

1. WatchlistProperties.ProvisioningState
1. WatchlistProperties.SasURI
