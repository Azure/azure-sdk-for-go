# Unreleased

## Additive Changes

### New Constants

1. ResourceTypeBasicRecoveryVirtualNetworkCustomDetails.ResourceTypeBasicRecoveryVirtualNetworkCustomDetailsResourceTypeNew

### New Funcs

1. ExistingRecoveryVirtualNetwork.AsNewRecoveryVirtualNetwork() (*NewRecoveryVirtualNetwork, bool)
1. NewRecoveryVirtualNetwork.AsBasicRecoveryVirtualNetworkCustomDetails() (BasicRecoveryVirtualNetworkCustomDetails, bool)
1. NewRecoveryVirtualNetwork.AsExistingRecoveryVirtualNetwork() (*ExistingRecoveryVirtualNetwork, bool)
1. NewRecoveryVirtualNetwork.AsNewRecoveryVirtualNetwork() (*NewRecoveryVirtualNetwork, bool)
1. NewRecoveryVirtualNetwork.AsRecoveryVirtualNetworkCustomDetails() (*RecoveryVirtualNetworkCustomDetails, bool)
1. NewRecoveryVirtualNetwork.MarshalJSON() ([]byte, error)
1. RecoveryVirtualNetworkCustomDetails.AsNewRecoveryVirtualNetwork() (*NewRecoveryVirtualNetwork, bool)

### Struct Changes

#### New Structs

1. NewRecoveryVirtualNetwork

#### New Struct Fields

1. HyperVReplicaAzureEnableProtectionInput.TargetAvailabilitySetID
1. HyperVReplicaAzureEnableProtectionInput.TargetVMSize
1. InMageAzureV2EnableProtectionInput.TargetAvailabilitySetID
1. InMageAzureV2EnableProtectionInput.TargetVMSize
