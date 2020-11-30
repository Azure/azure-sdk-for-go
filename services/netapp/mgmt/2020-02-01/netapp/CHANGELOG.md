Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Const `Standard` type has been changed from `ServiceLevel` to `PatchServiceLevel`
- Const `Ultra` type has been changed from `ServiceLevel` to `PatchServiceLevel`
- Const `Premium` type has been changed from `ServiceLevel` to `PatchServiceLevel`
- Type of `PoolProperties.ServiceLevel` has been changed from `ServiceLevel` to `PoolServiceLevel`
- Type of `PoolPatchProperties.ServiceLevel` has been changed from `ServiceLevel` to `PatchServiceLevel`
- Type of `VolumeProperties.ServiceLevel` has been changed from `ServiceLevel` to `VolumeServiceLevel`
- Const `Monthly` has been removed
- Const `Weekly` has been removed

## New Content

- New const `VolumeServiceLevelStandard`
- New const `ServiceLevelStandard`
- New const `VolumeServiceLevelUltra`
- New const `ServiceLevelUltra`
- New const `PoolServiceLevelPremium`
- New const `PoolServiceLevelUltra`
- New const `ServiceLevelPremium`
- New const `PoolServiceLevelStandard`
- New const `VolumeServiceLevelPremium`
- New function `PossibleVolumeServiceLevelValues() []VolumeServiceLevel`
- New function `PossiblePatchServiceLevelValues() []PatchServiceLevel`
- New function `PossiblePoolServiceLevelValues() []PoolServiceLevel`
