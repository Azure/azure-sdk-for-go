# Release History

## 0.6.0 (2022-06-02)
### Breaking Changes

- Type of `Schedule.DayOfWeek` has been changed from `*DayOfWeek` to `*AssessmentDayOfWeek`
- Type of `AutoBackupSettings.DaysOfWeek` has been changed from `[]*DaysOfWeek` to `[]*AutoBackupDaysOfWeek`
- Const `DaysOfWeekThursday` has been removed
- Const `DaysOfWeekMonday` has been removed
- Const `DaysOfWeekWednesday` has been removed
- Const `DaysOfWeekSaturday` has been removed
- Const `DaysOfWeekTuesday` has been removed
- Const `DaysOfWeekSunday` has been removed
- Const `DaysOfWeekFriday` has been removed
- Function `PossibleDaysOfWeekValues` has been removed

### Features Added

- New const `AutoBackupDaysOfWeekTuesday`
- New const `AssessmentDayOfWeekThursday`
- New const `DayOfWeekEveryday`
- New const `AutoBackupDaysOfWeekSaturday`
- New const `AssessmentDayOfWeekFriday`
- New const `AssessmentDayOfWeekSunday`
- New const `AutoBackupDaysOfWeekMonday`
- New const `AssessmentDayOfWeekMonday`
- New const `AutoBackupDaysOfWeekWednesday`
- New const `AssessmentDayOfWeekWednesday`
- New const `AssessmentDayOfWeekTuesday`
- New const `AutoBackupDaysOfWeekFriday`
- New const `AutoBackupDaysOfWeekSunday`
- New const `AssessmentDayOfWeekSaturday`
- New const `AutoBackupDaysOfWeekThursday`
- New function `PossibleAssessmentDayOfWeekValues() []AssessmentDayOfWeek`
- New function `PossibleAutoBackupDaysOfWeekValues() []AutoBackupDaysOfWeek`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sqlvirtualmachine/armsqlvirtualmachine` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).