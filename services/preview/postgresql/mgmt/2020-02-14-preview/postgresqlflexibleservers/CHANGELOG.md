# Unreleased

## Breaking Changes

### Signature Changes

#### Const Types

1. Disabled changed type from HAEnabledEnum to GeoRedundantBackupEnum
1. Enabled changed type from HAEnabledEnum to GeoRedundantBackupEnum

## Additive Changes

### New Constants

1. HAEnabledEnum.HAEnabledEnumDisabled
1. HAEnabledEnum.HAEnabledEnumEnabled

### New Funcs

1. PossibleGeoRedundantBackupEnumValues() []GeoRedundantBackupEnum

### Struct Changes

#### New Struct Fields

1. ServerProperties.EarliestRestoreDate
1. ServerProperties.LogBackupStorageSku
1. ServerProperties.MinorVersion
1. ServerProperties.StandbyCount
1. ServerPropertiesForUpdate.StandbyCount
1. StorageProfile.GeoRedundantBackup
