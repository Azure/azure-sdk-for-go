Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Type of `SnapshotPolicyProperties.HourlySchedule` has been changed from `interface{}` to `*HourlySchedule`
- Type of `SnapshotPolicyProperties.DailySchedule` has been changed from `interface{}` to `*DailySchedule`
- Type of `SnapshotPolicyProperties.WeeklySchedule` has been changed from `interface{}` to `*WeeklySchedule`
- Type of `SnapshotPolicyProperties.MonthlySchedule` has been changed from `interface{}` to `*MonthlySchedule`
- Const `Weekly` has been removed
- Const `Monthly` has been removed

## New Content

- New function `SnapshotPolicyProperties.MarshalJSON() ([]byte, error)`
- New field `BackupID` in struct `BackupProperties`
- New field `ProvisioningState` in struct `SnapshotPolicyProperties`
- New field `Name` in struct `SnapshotPolicyProperties`
