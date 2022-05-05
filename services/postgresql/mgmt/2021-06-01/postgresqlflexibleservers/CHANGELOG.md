# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. ServerProperties.Tags

## Additive Changes

### New Constants

1. Reason.ReasonAlreadyExists
1. Reason.ReasonInvalid

### New Funcs

1. NameAvailability.MarshalJSON() ([]byte, error)
1. PossibleReasonValues() []Reason

### Struct Changes

#### New Struct Fields

1. ConfigurationProperties.DocumentationLink
1. ConfigurationProperties.IsConfigPendingRestart
1. ConfigurationProperties.IsDynamicConfig
1. ConfigurationProperties.IsReadOnly
1. ConfigurationProperties.Unit
1. NameAvailability.Reason
