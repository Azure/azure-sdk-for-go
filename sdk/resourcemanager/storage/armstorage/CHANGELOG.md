# Release History

## 2.0.0 (2022-11-14)
### Breaking Changes

- Struct `CloudError` has been removed
- Struct `CloudErrorBody` has been removed

### Features Added

- New const `ListEncryptionScopesIncludeEnabled`
- New const `ListEncryptionScopesIncludeAll`
- New const `ListEncryptionScopesIncludeDisabled`
- New type alias `ListEncryptionScopesInclude`
- New function `PossibleListEncryptionScopesIncludeValues() []ListEncryptionScopesInclude`
- New field `TierToCold` in struct `ManagementPolicyBaseBlob`
- New field `TierToHot` in struct `ManagementPolicyBaseBlob`
- New field `FailoverType` in struct `AccountsClientBeginFailoverOptions`
- New field `TierToCold` in struct `ManagementPolicyVersion`
- New field `TierToHot` in struct `ManagementPolicyVersion`
- New field `Filter` in struct `EncryptionScopesClientListOptions`
- New field `Include` in struct `EncryptionScopesClientListOptions`
- New field `Maxpagesize` in struct `EncryptionScopesClientListOptions`
- New field `TierToCold` in struct `ManagementPolicySnapShot`
- New field `TierToHot` in struct `ManagementPolicySnapShot`


## 1.1.0 (2022-08-10)
### Features Added

- New const `DirectoryServiceOptionsAADKERB`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).