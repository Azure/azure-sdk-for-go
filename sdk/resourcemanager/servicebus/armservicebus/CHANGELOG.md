# Release History

## 2.0.0-beta.4 (2026-03-19)
### Breaking Changes

- Type of `Identity.UserAssignedIdentities` has been changed from `map[string]*UserAssignedIdentity` to `map[string]*DictionaryValue`
- Type of `SBNamespaceUpdateParameters.Properties` has been changed from `*SBNamespaceUpdateProperties` to `*SBNamespaceProperties`
- Enum `PublicNetworkAccess` has been removed
- Enum `PublicNetworkAccessFlag` has been removed
- Enum `TLSVersion` has been removed
- Struct `ProxyResource` has been removed
- Struct `SBClientAffineProperties` has been removed
- Struct `SBNamespaceUpdateProperties` has been removed
- Struct `UserAssignedIdentity` has been removed
- Field `Location` of struct `ArmDisasterRecovery` has been removed
- Field `Location` of struct `MigrationConfigProperties` has been removed
- Field `Location` of struct `NetworkRuleSet` has been removed
- Field `PublicNetworkAccess`, `TrustedServiceAccessEnabled` of struct `NetworkRuleSetProperties` has been removed
- Field `IsDataAction`, `Origin`, `Properties` of struct `Operation` has been removed
- Field `Description` of struct `OperationDisplay` has been removed
- Field `Location` of struct `PrivateEndpointConnection` has been removed
- Field `Location` of struct `Rule` has been removed
- Field `Location` of struct `SBAuthorizationRule` has been removed
- Field `AlternateName`, `DisableLocalAuth`, `MinimumTLSVersion`, `PremiumMessagingPartitions`, `PublicNetworkAccess` of struct `SBNamespaceProperties` has been removed
- Field `Location` of struct `SBQueue` has been removed
- Field `MaxMessageSizeInKilobytes` of struct `SBQueueProperties` has been removed
- Field `Location` of struct `SBSubscription` has been removed
- Field `ClientAffineProperties`, `IsClientAffine` of struct `SBSubscriptionProperties` has been removed
- Field `Location` of struct `SBTopic` has been removed
- Field `MaxMessageSizeInKilobytes` of struct `SBTopicProperties` has been removed

### Features Added

- New struct `DictionaryValue`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-04-07)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).