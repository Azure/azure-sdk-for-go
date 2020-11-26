
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Const `Premium` type has been changed from `ServiceLevel` to `PatchServiceLevel`
- Const `Standard` type has been changed from `ServiceLevel` to `PatchServiceLevel`
- Const `Ultra` type has been changed from `ServiceLevel` to `PatchServiceLevel`
- Type of `VolumeProperties.ServiceLevel` has been changed from `ServiceLevel` to `VolumeServiceLevel`
- Type of `PoolPatchProperties.ServiceLevel` has been changed from `ServiceLevel` to `PatchServiceLevel`
- Type of `PoolProperties.ServiceLevel` has been changed from `ServiceLevel` to `PoolServiceLevel`
- Const `Monthly` has been removed
- Const `Weekly` has been removed

## New Content

- Const `ServiceLevelStandard` is added
- Const `VolumeServiceLevelStandard` is added
- Const `ServiceLevelPremium` is added
- Const `VolumeServiceLevelPremium` is added
- Const `PoolServiceLevelUltra` is added
- Const `PoolServiceLevelStandard` is added
- Const `ServiceLevelUltra` is added
- Const `VolumeServiceLevelUltra` is added
- Const `PoolServiceLevelPremium` is added
- Function `PossiblePoolServiceLevelValues() []PoolServiceLevel` is added
- Function `PossibleVolumeServiceLevelValues() []VolumeServiceLevel` is added
- Function `PossiblePatchServiceLevelValues() []PatchServiceLevel` is added

