# Unreleased

## Additive Changes

### New Constants

1. IPRuleAction.IPRuleActionAllow
1. ResourceIdentityType.ResourceIdentityTypeSystemAssigned
1. SkuTier.SkuTierBasic
1. SkuTier.SkuTierFree
1. SkuTier.SkuTierPremium
1. SkuTier.SkuTierStandard

### New Funcs

1. Identity.MarshalJSON() ([]byte, error)
1. NetworkRuleSetIPRule.MarshalJSON() ([]byte, error)
1. PossibleIPRuleActionValues() []IPRuleAction
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

1. NetworkRuleSetIPRule.Action
