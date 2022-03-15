# Unreleased

## Breaking Changes

### Removed Constants

1. QosType.Auto

### Signature Changes

#### Const Types

1. Manual changed type from QosType to BackupType

#### Struct Fields

1. BackupProperties.BackupType changed type from *string to BackupType

## Additive Changes

### New Constants

1. BackupType.Scheduled
1. QosType.QosTypeAuto
1. QosType.QosTypeManual

### New Funcs

1. PossibleBackupTypeValues() []BackupType
