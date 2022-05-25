# Unreleased

## Breaking Changes

### Removed Constants

1. DaysOfWeek.DaysOfWeekFriday
1. DaysOfWeek.DaysOfWeekMonday
1. DaysOfWeek.DaysOfWeekSaturday
1. DaysOfWeek.DaysOfWeekSunday
1. DaysOfWeek.DaysOfWeekThursday
1. DaysOfWeek.DaysOfWeekTuesday
1. DaysOfWeek.DaysOfWeekWednesday

### Removed Funcs

1. PossibleDaysOfWeekValues() []DaysOfWeek

### Signature Changes

#### Struct Fields

1. AutoBackupSettings.DaysOfWeek changed type from *[]DaysOfWeek to *[]AutoBackupDaysOfWeek
1. Schedule.DayOfWeek changed type from DayOfWeek to AssessmentDayOfWeek

## Additive Changes

### New Constants

1. AssessmentDayOfWeek.AssessmentDayOfWeekFriday
1. AssessmentDayOfWeek.AssessmentDayOfWeekMonday
1. AssessmentDayOfWeek.AssessmentDayOfWeekSaturday
1. AssessmentDayOfWeek.AssessmentDayOfWeekSunday
1. AssessmentDayOfWeek.AssessmentDayOfWeekThursday
1. AssessmentDayOfWeek.AssessmentDayOfWeekTuesday
1. AssessmentDayOfWeek.AssessmentDayOfWeekWednesday
1. AutoBackupDaysOfWeek.AutoBackupDaysOfWeekFriday
1. AutoBackupDaysOfWeek.AutoBackupDaysOfWeekMonday
1. AutoBackupDaysOfWeek.AutoBackupDaysOfWeekSaturday
1. AutoBackupDaysOfWeek.AutoBackupDaysOfWeekSunday
1. AutoBackupDaysOfWeek.AutoBackupDaysOfWeekThursday
1. AutoBackupDaysOfWeek.AutoBackupDaysOfWeekTuesday
1. AutoBackupDaysOfWeek.AutoBackupDaysOfWeekWednesday
1. DayOfWeek.DayOfWeekEveryday

### New Funcs

1. PossibleAssessmentDayOfWeekValues() []AssessmentDayOfWeek
1. PossibleAutoBackupDaysOfWeekValues() []AutoBackupDaysOfWeek
