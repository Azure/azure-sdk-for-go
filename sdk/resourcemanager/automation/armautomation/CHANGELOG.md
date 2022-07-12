# Release History

## 0.7.0 (2022-07-12)
### Breaking Changes

- Function `*DscConfigurationClient.UpdateWithJSON` parameter(s) have been changed from `(context.Context, string, string, string, *DscConfigurationClientUpdateWithJSONOptions)` to `(context.Context, string, string, string, DscConfigurationUpdateParameters, *DscConfigurationClientUpdateWithJSONOptions)`
- Function `*DscConfigurationClient.UpdateWithText` parameter(s) have been changed from `(context.Context, string, string, string, *DscConfigurationClientUpdateWithTextOptions)` to `(context.Context, string, string, string, string, *DscConfigurationClientUpdateWithTextOptions)`
- Struct `HybridRunbookWorkerGroupUpdateParameters` has been removed
- Struct `HybridRunbookWorkerLegacy` has been removed


## 0.6.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automation/armautomation` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).