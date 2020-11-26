
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewSpatialAnchorsAccountListPage` signature has been changed from `(func(context.Context, SpatialAnchorsAccountList) (SpatialAnchorsAccountList, error))` to `(SpatialAnchorsAccountList,func(context.Context, SpatialAnchorsAccountList) (SpatialAnchorsAccountList, error))`
- Function `NewOperationListPage` signature has been changed from `(func(context.Context, OperationList) (OperationList, error))` to `(OperationList,func(context.Context, OperationList) (OperationList, error))`
- Type of `CheckNameAvailabilityResponse.NameAvailable` has been changed from `NameAvailability` to `*bool`
- Const `False` has been removed
- Const `True` has been removed
- Function `PossibleNameAvailabilityValues` has been removed

## New Content

- Const `Free` is added
- Const `Premium` is added
- Const `Standard` is added
- Const `SystemAssigned` is added
- Const `Basic` is added
- Function `PossibleSkuTierValues() []SkuTier` is added
- Function `PossibleResourceIdentityTypeValues() []ResourceIdentityType` is added
- Function `ResourceModelWithAllowedPropertySet.MarshalJSON() ([]byte,error)` is added
- Function `ResourceModelWithAllowedPropertySetIdentity.MarshalJSON() ([]byte,error)` is added
- Function `Identity.MarshalJSON() ([]byte,error)` is added
- Struct `Identity` is added
- Struct `Plan` is added
- Struct `ResourceModelWithAllowedPropertySet` is added
- Struct `ResourceModelWithAllowedPropertySetIdentity` is added
- Struct `ResourceModelWithAllowedPropertySetPlan` is added
- Struct `ResourceModelWithAllowedPropertySetSku` is added
- Struct `Sku` is added
- Field `Identity` is added to struct `SpatialAnchorsAccount`
- Field `IsDataAction` is added to struct `Operation`

