# Unreleased

## Breaking Changes

### Removed Constants

1. CaseSeverity.CaseSeverityCritical
1. CaseSeverity.CaseSeverityHigh
1. CaseSeverity.CaseSeverityInformational
1. CaseSeverity.CaseSeverityLow
1. CaseSeverity.CaseSeverityMedium
1. KindBasicSettings.KindBasicSettingsKindSettings
1. KindBasicSettings.KindBasicSettingsKindToggleSettings
1. KindBasicSettings.KindBasicSettingsKindUebaSettings
1. LicenseStatus.LicenseStatusDisabled
1. LicenseStatus.LicenseStatusEnabled
1. SettingKind.SettingKindToggleSettings
1. SettingKind.SettingKindUebaSettings
1. StatusInMcas.StatusInMcasDisabled
1. StatusInMcas.StatusInMcasEnabled

### Removed Funcs

1. *ToggleSettings.UnmarshalJSON([]byte) error
1. *UebaSettings.UnmarshalJSON([]byte) error
1. PossibleCaseSeverityValues() []CaseSeverity
1. PossibleKindBasicSettingsValues() []KindBasicSettings
1. PossibleLicenseStatusValues() []LicenseStatus
1. PossibleSettingKindValues() []SettingKind
1. PossibleStatusInMcasValues() []StatusInMcas
1. Settings.AsBasicSettings() (BasicSettings, bool)
1. Settings.AsSettings() (*Settings, bool)
1. Settings.AsToggleSettings() (*ToggleSettings, bool)
1. Settings.AsUebaSettings() (*UebaSettings, bool)
1. Settings.MarshalJSON() ([]byte, error)
1. ToggleSettings.AsBasicSettings() (BasicSettings, bool)
1. ToggleSettings.AsSettings() (*Settings, bool)
1. ToggleSettings.AsToggleSettings() (*ToggleSettings, bool)
1. ToggleSettings.AsUebaSettings() (*UebaSettings, bool)
1. ToggleSettings.MarshalJSON() ([]byte, error)
1. UebaSettings.AsBasicSettings() (BasicSettings, bool)
1. UebaSettings.AsSettings() (*Settings, bool)
1. UebaSettings.AsToggleSettings() (*ToggleSettings, bool)
1. UebaSettings.AsUebaSettings() (*UebaSettings, bool)
1. UebaSettings.MarshalJSON() ([]byte, error)
1. UebaSettingsProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. Settings
1. ToggleSettings
1. ToggleSettingsProperties
1. UebaSettings
1. UebaSettingsProperties

### Signature Changes

#### Struct Fields

1. IncidentInfo.Severity changed type from CaseSeverity to IncidentSeverity
