# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Structs

1. RestorePointProvisioningDetails

#### Removed Struct Fields

1. RestorePointProperties.ProvisioningDetails

### Signature Changes

#### Struct Fields

1. OrchestrationServiceStateInput.ServiceName changed type from *string to OrchestrationServiceNames

## Additive Changes

### New Constants

1. DiskCreateOption.DiskCreateOptionCopyStart
1. DiskState.DiskStateActiveSASFrozen
1. DiskState.DiskStateFrozen
1. PublicNetworkAccess.PublicNetworkAccessDisabled
1. PublicNetworkAccess.PublicNetworkAccessEnabled

### New Funcs

1. PossiblePublicNetworkAccessValues() []PublicNetworkAccess

### Struct Changes

#### New Structs

1. SupportedCapabilities

#### New Struct Fields

1. DiskAccess.ExtendedLocation
1. DiskProperties.CompletionPercent
1. DiskProperties.PublicNetworkAccess
1. DiskProperties.SupportedCapabilities
1. DiskRestorePointProperties.CompletionPercent
1. DiskRestorePointProperties.DiskAccessID
1. DiskRestorePointProperties.NetworkAccessPolicy
1. DiskRestorePointProperties.PublicNetworkAccess
1. DiskRestorePointProperties.SupportedCapabilities
1. DiskUpdateProperties.PublicNetworkAccess
1. DiskUpdateProperties.SupportedCapabilities
1. EncryptionSetProperties.AutoKeyRotationError
1. RestorePointProperties.TimeCreated
1. SnapshotProperties.CompletionPercent
1. SnapshotProperties.PublicNetworkAccess
1. SnapshotProperties.SupportedCapabilities
1. SnapshotUpdateProperties.PublicNetworkAccess
