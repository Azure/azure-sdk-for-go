# Release History

## 1.0.0 (2024-06-28)
### Other Changes

- Release stable version.


## 0.2.0 (2024-06-26)
### Breaking Changes

- Type of `CloudExadataInfrastructureProperties.DataStorageSizeInTbs` has been changed from `*int32` to `*float64`
- Type of `CloudVMClusterProperties.NsgCidrs` has been changed from `[]*NSGCidr` to `[]*NsgCidr`
- Type of `OracleSubscriptionUpdate.Plan` has been changed from `*ResourcePlanTypeUpdate` to `*PlanUpdate`
- Struct `NSGCidr` has been removed
- Struct `ResourcePlanTypeUpdate` has been removed
- Field `AutonomousDatabaseID`, `DatabaseSizeInTBs`, `SizeInTBs`, `Type` of struct `AutonomousDatabaseBackupProperties` has been removed

### Features Added

- New enum type `RepeatCadenceType` with values `RepeatCadenceTypeMonthly`, `RepeatCadenceTypeOneTime`, `RepeatCadenceTypeWeekly`, `RepeatCadenceTypeYearly`
- New function `*AutonomousDatabasesClient.BeginRestore(context.Context, string, string, RestoreAutonomousDatabaseDetails, *AutonomousDatabasesClientBeginRestoreOptions) (*runtime.Poller[AutonomousDatabasesClientRestoreResponse], error)`
- New function `*AutonomousDatabasesClient.BeginShrink(context.Context, string, string, *AutonomousDatabasesClientBeginShrinkOptions) (*runtime.Poller[AutonomousDatabasesClientShrinkResponse], error)`
- New function `*ClientFactory.NewSystemVersionsClient() *SystemVersionsClient`
- New function `NewSystemVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SystemVersionsClient, error)`
- New function `*SystemVersionsClient.Get(context.Context, string, string, *SystemVersionsClientGetOptions) (SystemVersionsClientGetResponse, error)`
- New function `*SystemVersionsClient.NewListByLocationPager(string, *SystemVersionsClientListByLocationOptions) *runtime.Pager[SystemVersionsClientListByLocationResponse]`
- New struct `LongTermBackUpScheduleDetails`
- New struct `NsgCidr`
- New struct `PlanUpdate`
- New struct `RestoreAutonomousDatabaseDetails`
- New struct `SystemVersion`
- New struct `SystemVersionListResult`
- New struct `SystemVersionProperties`
- New field `AutonomousDatabaseOcid`, `BackupType`, `DatabaseSizeInTbs`, `SizeInTbs`, `TimeStarted` in struct `AutonomousDatabaseBackupProperties`
- New field `LongTermBackupSchedule`, `NextLongTermBackupTimeStamp` in struct `AutonomousDatabaseBaseProperties`
- New field `LongTermBackupSchedule`, `NextLongTermBackupTimeStamp` in struct `AutonomousDatabaseCloneProperties`
- New field `LongTermBackupSchedule`, `NextLongTermBackupTimeStamp` in struct `AutonomousDatabaseProperties`
- New field `LongTermBackupSchedule` in struct `AutonomousDatabaseUpdateProperties`


## 0.1.0 (2024-05-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/oracledatabase/armoracledatabase` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
