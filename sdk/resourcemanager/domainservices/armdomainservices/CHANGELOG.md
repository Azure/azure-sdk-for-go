# Release History

## 2.0.0 (2022-11-10)
### Breaking Changes

- Struct `CloudError` has been removed
- Struct `CloudErrorBody` has been removed

### Features Added

- New const `ChannelBindingDisabled`
- New const `SyncScopeCloudOnly`
- New const `LdapSigningEnabled`
- New const `ChannelBindingEnabled`
- New const `SyncScopeAll`
- New const `LdapSigningDisabled`
- New type alias `SyncScope`
- New type alias `ChannelBinding`
- New type alias `LdapSigning`
- New function `PossibleLdapSigningValues() []LdapSigning`
- New function `PossibleChannelBindingValues() []ChannelBinding`
- New function `PossibleSyncScopeValues() []SyncScope`
- New field `SyncScope` in struct `DomainServiceProperties`
- New field `ChannelBinding` in struct `DomainSecuritySettings`
- New field `LdapSigning` in struct `DomainSecuritySettings`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/domainservices/armdomainservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).