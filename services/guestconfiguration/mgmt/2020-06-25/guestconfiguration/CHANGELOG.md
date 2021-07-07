# Unreleased

## Breaking Changes

### Removed Constants

1. ConfigurationMode.ApplyOnly

### Removed Funcs

1. Navigation.MarshalJSON() ([]byte, error)

### Signature Changes

#### Const Types

1. ApplyAndAutoCorrect changed type from ConfigurationMode to AssignmentType
1. ApplyAndMonitor changed type from ConfigurationMode to AssignmentType

## Additive Changes

### New Constants

1. AssignmentType.Audit
1. AssignmentType.DeployAndAutoCorrect
1. ConfigurationMode.ConfigurationModeApplyAndAutoCorrect
1. ConfigurationMode.ConfigurationModeApplyAndMonitor
1. ConfigurationMode.ConfigurationModeApplyOnly

### New Funcs

1. PossibleAssignmentTypeValues() []AssignmentType

### Struct Changes

#### New Struct Fields

1. Navigation.AssignmentType
