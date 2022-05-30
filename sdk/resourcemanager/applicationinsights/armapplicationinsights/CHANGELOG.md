# Release History

## 0.6.0 (2022-05-30)
### Breaking Changes

- Type of `WorkbookUpdateParameters.Kind` has been changed from `*SharedTypeKind` to `*WorkbookUpdateSharedTypeKind`
- Type of `WorkbookResource.Kind` has been changed from `*Kind` to `*WorkbookSharedTypeKind`
- Type of `Workbook.Kind` has been changed from `*Kind` to `*WorkbookSharedTypeKind`
- Const `SharedTypeKindUser` has been removed
- Const `SharedTypeKindShared` has been removed
- Function `PossibleSharedTypeKindValues` has been removed

### Features Added

- New const `WorkbookSharedTypeKindShared`
- New const `WorkbookUpdateSharedTypeKindShared`
- New function `PossibleWorkbookSharedTypeKindValues() []WorkbookSharedTypeKind`
- New function `PossibleWorkbookUpdateSharedTypeKindValues() []WorkbookUpdateSharedTypeKind`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).