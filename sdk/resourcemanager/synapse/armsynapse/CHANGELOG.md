# Release History

## 0.6.0 (2022-08-29)
### Breaking Changes

- Type of `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentity.ActualState` has been changed from `*ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityActualState` to `*ActualState`
- Type of `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentity.DesiredState` has been changed from `*ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityDesiredState` to `*DesiredState`
- Const `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityActualStateDisabled` has been removed
- Const `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityDesiredStateEnabled` has been removed
- Const `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityActualStateUnknown` has been removed
- Const `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityActualStateEnabled` has been removed
- Const `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityDesiredStateDisabled` has been removed
- Const `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityActualStateDisabling` has been removed
- Const `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityActualStateEnabling` has been removed
- Type alias `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityActualState` has been removed
- Type alias `ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityDesiredState` has been removed
- Function `PossibleManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityActualStateValues` has been removed
- Function `PossibleManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentityDesiredStateValues` has been removed

### Features Added

- New const `ActualStateEnabled`
- New const `ActualStateDisabled`
- New const `ActualStateDisabling`
- New const `ActualStateEnabling`
- New const `DesiredStateEnabled`
- New const `ActualStateUnknown`
- New const `DesiredStateDisabled`
- New type alias `ActualState`
- New type alias `DesiredState`
- New function `PossibleActualStateValues() []ActualState`
- New function `PossibleDesiredStateValues() []DesiredState`
- New field `AADObjectID` in struct `ClusterPrincipalProperties`
- New field `AADObjectID` in struct `DatabasePrincipalProperties`
- New field `IsAutotuneEnabled` in struct `BigDataPoolResourceProperties`
- New field `ConfigMergeRule` in struct `SparkConfigurationInfo`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).