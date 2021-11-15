# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Structs

1. DefaultErrorResponse

#### Removed Struct Fields

1. RestorableDatabaseAccountGetResult.Identity
1. RestorableDatabaseAccountGetResult.Tags
1. RestorableMongodbCollectionGetResult.Identity
1. RestorableMongodbCollectionGetResult.Location
1. RestorableMongodbCollectionGetResult.Tags
1. RestorableMongodbDatabaseGetResult.Identity
1. RestorableMongodbDatabaseGetResult.Location
1. RestorableMongodbDatabaseGetResult.Tags
1. RestorableSQLContainerGetResult.Identity
1. RestorableSQLContainerGetResult.Location
1. RestorableSQLContainerGetResult.Tags
1. RestorableSQLDatabaseGetResult.Identity
1. RestorableSQLDatabaseGetResult.Location
1. RestorableSQLDatabaseGetResult.Tags

## Additive Changes

### New Constants

1. BackupStorageRedundancy.Geo
1. BackupStorageRedundancy.Local
1. BackupStorageRedundancy.Zone

### New Funcs

1. PossibleBackupStorageRedundancyValues() []BackupStorageRedundancy

### Struct Changes

#### New Structs

1. CloudError

#### New Struct Fields

1. PeriodicModeProperties.BackupStorageRedundancy
