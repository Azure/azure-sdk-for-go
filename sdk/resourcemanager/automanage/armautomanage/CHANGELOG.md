# Release History

## 1.0.0 (2022-08-01)
### Breaking Changes

- Function `*ConfigurationProfilesVersionsClient.Update` has been removed
- Struct `ConfigurationProfilesVersionsClientUpdateOptions` has been removed
- Struct `ConfigurationProfilesVersionsClientUpdateResponse` has been removed
- Field `ProfileOverrides` of struct `ConfigurationProfileAssignmentProperties` has been removed
- Field `Overrides` of struct `ConfigurationProfileProperties` has been removed

### Features Added

- New field `ManagedBy` in struct `ConfigurationProfileAssignment`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automanage/armautomanage` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).