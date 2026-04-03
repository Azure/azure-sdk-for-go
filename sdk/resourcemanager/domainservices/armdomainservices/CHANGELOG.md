# Release History

## 1.3.0-beta.1 (2026-04-03)
### Features Added

- New enum type `ChannelBinding` with values `ChannelBindingDisabled`, `ChannelBindingEnabled`
- New enum type `LdapSigning` with values `LdapSigningDisabled`, `LdapSigningEnabled`
- New enum type `SyncOnPremSamAccountName` with values `SyncOnPremSamAccountNameDisabled`, `SyncOnPremSamAccountNameEnabled`
- New enum type `SyncScope` with values `SyncScopeAll`, `SyncScopeCloudOnly`
- New function `*Client.Unsuspend(ctx context.Context, resourceGroupName string, domainServiceName string, options *ClientUnsuspendOptions) (ClientUnsuspendResponse, error)`
- New struct `UnsuspendDomainServiceResponse`
- New field `ChannelBinding`, `LdapSigning`, `SyncOnPremSamAccountName` in struct `DomainSecuritySettings`
- New field `SyncApplicationID`, `SyncScope` in struct `DomainServiceProperties`
- New field `SelfUnsuspendCounter` in struct `ReplicaSet`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/domainservices/armdomainservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).