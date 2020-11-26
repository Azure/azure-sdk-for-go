
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Type of `SnapshotPolicyProperties.HourlySchedule` has been changed from `interface{}` to `*HourlySchedule`
- Type of `SnapshotPolicyProperties.DailySchedule` has been changed from `interface{}` to `*DailySchedule`
- Type of `SnapshotPolicyProperties.WeeklySchedule` has been changed from `interface{}` to `*WeeklySchedule`
- Type of `SnapshotPolicyProperties.MonthlySchedule` has been changed from `interface{}` to `*MonthlySchedule`
- Const `Monthly` has been removed
- Const `Weekly` has been removed

## New Content

- Function `SnapshotPolicyProperties.MarshalJSON() ([]byte,error)` is added
- Field `BackupID` is added to struct `BackupProperties`
- Field `ProvisioningState` is added to struct `SnapshotPolicyProperties`
- Field `Name` is added to struct `SnapshotPolicyProperties`

