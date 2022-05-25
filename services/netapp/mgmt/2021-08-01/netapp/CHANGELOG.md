# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. VolumeGroup.Tags

## Additive Changes

### New Constants

1. ResourceIdentityType.ResourceIdentityTypeSystemAssigned
1. SkuTier.SkuTierBasic
1. SkuTier.SkuTierFree
1. SkuTier.SkuTierPremium
1. SkuTier.SkuTierStandard

### New Funcs

1. Identity.MarshalJSON() ([]byte, error)
1. PossibleResourceIdentityTypeValues() []ResourceIdentityType
1. PossibleSkuTierValues() []SkuTier
1. ResourceModelWithAllowedPropertySet.MarshalJSON() ([]byte, error)
1. ResourceModelWithAllowedPropertySetIdentity.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. Identity
1. Plan
1. ResourceModelWithAllowedPropertySet
1. ResourceModelWithAllowedPropertySetIdentity
1. ResourceModelWithAllowedPropertySetPlan
1. ResourceModelWithAllowedPropertySetSku
1. Sku

#### New Struct Fields

1. AzureEntityResource.SystemData
1. BackupPolicy.SystemData
1. CapacityPool.SystemData
1. ProxyResource.SystemData
1. Resource.SystemData
1. SnapshotPolicy.SystemData
1. TrackedResource.SystemData
1. Volume.SystemData
