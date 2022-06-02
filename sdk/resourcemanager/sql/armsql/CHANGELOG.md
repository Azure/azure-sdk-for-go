# Release History

## 1.0.0 (2022-06-02)
### Features Added

- New const `AdvancedThreatProtectionStateDisabled`
- New const `AdvancedThreatProtectionStateEnabled`
- New const `AdvancedThreatProtectionNameDefault`
- New const `AdvancedThreatProtectionStateNew`
- New function `PossibleAdvancedThreatProtectionStateValues() []AdvancedThreatProtectionState`
- New function `PossibleAdvancedThreatProtectionNameValues() []AdvancedThreatProtectionName`
- New function `*AdvancedThreatProtectionProperties.UnmarshalJSON([]byte) error`
- New function `AdvancedThreatProtectionProperties.MarshalJSON() ([]byte, error)`
- New struct `AdvancedThreatProtectionProperties`
- New struct `DatabaseAdvancedThreatProtection`
- New struct `DatabaseAdvancedThreatProtectionListResult`
- New struct `DatabaseAdvancedThreatProtectionSettingsClientCreateOrUpdateOptions`
- New struct `DatabaseAdvancedThreatProtectionSettingsClientCreateOrUpdateResponse`
- New struct `DatabaseAdvancedThreatProtectionSettingsClientGetOptions`
- New struct `DatabaseAdvancedThreatProtectionSettingsClientGetResponse`
- New struct `DatabaseAdvancedThreatProtectionSettingsClientListByDatabaseOptions`
- New struct `DatabaseAdvancedThreatProtectionSettingsClientListByDatabaseResponse`
- New struct `LogicalServerAdvancedThreatProtectionListResult`
- New struct `ServerAdvancedThreatProtection`
- New struct `ServerAdvancedThreatProtectionSettingsClientBeginCreateOrUpdateOptions`
- New struct `ServerAdvancedThreatProtectionSettingsClientCreateOrUpdateResponse`
- New struct `ServerAdvancedThreatProtectionSettingsClientGetOptions`
- New struct `ServerAdvancedThreatProtectionSettingsClientGetResponse`
- New struct `ServerAdvancedThreatProtectionSettingsClientListByServerOptions`
- New struct `ServerAdvancedThreatProtectionSettingsClientListByServerResponse`


## 0.6.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).