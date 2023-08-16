# Release History

## 1.4.0-beta.2 (2023-07-19)

### Bug Fixes

- Fixed a potential panic in faked paged and long-running operations.

## 1.4.0-beta.1 (2023-06-12)

### Features Added

- Support for test fakes and OpenTelemetry trace spans.

## 1.3.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module

## 1.2.0 (2022-12-23)
### Features Added

- New type alias `ListEncryptionScopesInclude`
- New field `FailoverType` in struct `AccountsClientBeginFailoverOptions`
- New field `TierToCold` in struct `ManagementPolicyBaseBlob`
- New field `TierToHot` in struct `ManagementPolicyBaseBlob`
- New field `Filter` in struct `EncryptionScopesClientListOptions`
- New field `Include` in struct `EncryptionScopesClientListOptions`
- New field `Maxpagesize` in struct `EncryptionScopesClientListOptions`
- New field `TierToHot` in struct `ManagementPolicyVersion`
- New field `TierToCold` in struct `ManagementPolicyVersion`
- New field `TierToCold` in struct `ManagementPolicySnapShot`
- New field `TierToHot` in struct `ManagementPolicySnapShot`


## 1.1.0 (2022-08-10)
### Features Added

- New const `DirectoryServiceOptionsAADKERB`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).