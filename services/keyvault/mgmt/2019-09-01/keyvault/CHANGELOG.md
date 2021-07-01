# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. PrivateLinkServiceConnectionState.ActionRequired

## Additive Changes

### New Constants

1. VaultProvisioningState.VaultProvisioningStateRegisteringDNS
1. VaultProvisioningState.VaultProvisioningStateSucceeded

### New Funcs

1. PossibleVaultProvisioningStateValues() []VaultProvisioningState

### Struct Changes

#### New Structs

1. DimensionProperties
1. MetricSpecification

#### New Struct Fields

1. DeletedVaultProperties.PurgeProtectionEnabled
1. Operation.IsDataAction
1. PrivateEndpointConnection.Etag
1. PrivateEndpointConnectionItem.Etag
1. PrivateEndpointConnectionItem.ID
1. PrivateLinkServiceConnectionState.ActionsRequired
1. ServiceSpecification.MetricSpecifications
1. VaultProperties.HsmPoolResourceID
1. VaultProperties.ProvisioningState
1. VirtualNetworkRule.IgnoreMissingVnetServiceEndpoint
