# Release History

## 2.0.0 (2025-11-26)
### Breaking Changes

- Type of `Operation.Display` has been changed from `*OperationDisplayProperties` to `*OperationDisplay`
- Operation `*HierarchySettingsClient.List` has supported pagination, use `*HierarchySettingsClient.NewListPager` instead.
- Struct `CreateManagementGroupChildInfo` has been removed
- Struct `EntityHierarchyItem` has been removed
- Struct `EntityHierarchyItemProperties` has been removed
- Struct `ErrorDetails` has been removed
- Struct `ErrorResponse` has been removed
- Struct `OperationDisplayProperties` has been removed
- Struct `OperationResults` has been removed
- Field `Count` of struct `EntityListResult` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New struct `OperationDisplay`
- New struct `SystemData`
- New field `SystemData` in struct `HierarchySettings`
- New field `SystemData` in struct `ManagementGroup`
- New field `ActionType`, `IsDataAction`, `Origin` in struct `Operation`
- New field `SystemData` in struct `SubscriptionUnderManagementGroup`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).