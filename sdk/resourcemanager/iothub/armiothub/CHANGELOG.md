# Release History

## 2.0.0-beta.1 (2025-10-15)
### Breaking Changes

- Field `DefaultTTLAsIso8601` of struct `CloudToDeviceProperties` has been removed
- Field `LockDurationAsIso8601`, `TTLAsIso8601` of struct `FeedbackProperties` has been removed
- Field `LockDurationAsIso8601`, `TTLAsIso8601` of struct `MessagingEndpointProperties` has been removed
- Field `SasTTLAsIso8601` of struct `StorageEndpointProperties` has been removed

### Features Added

- New value `IotHubSKUGEN2` added to enum type `IotHubSKU`
- New value `IotHubSKUTierGeneration2` added to enum type `IotHubSKUTier`
- New value `RoutingSourceDigitalTwinChangeEvents`, `RoutingSourceMqttBrokerMessages` added to enum type `RoutingSource`
- New enum type `IPVersion` with values `IPVersionIPv4`, `IPVersionIPv4IPv6`, `IPVersionIPv6`
- New struct `DeviceRegistry`
- New struct `EncryptionPropertiesDescription`
- New struct `KeyVaultKeyProperties`
- New struct `PropertiesDeviceStreams`
- New struct `RootCertificateProperties`
- New field `PolicyResourceID` in struct `CertificateProperties`
- New field `PolicyResourceID` in struct `CertificatePropertiesWithNonce`
- New field `DefaultTTLAsISO8601` in struct `CloudToDeviceProperties`
- New field `LockDurationAsISO8601`, `TTLAsISO8601` in struct `FeedbackProperties`
- New field `LockDurationAsISO8601`, `TTLAsISO8601` in struct `MessagingEndpointProperties`
- New field `DeviceRegistry`, `DeviceStreams`, `Encryption`, `IPVersion`, `RootCertificate` in struct `Properties`
- New field `SasTTLAsISO8601` in struct `StorageEndpointProperties`


## 1.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-09-22)
### Features Added

- New struct `RoutingCosmosDBSQLAPIProperties`
- New field `CosmosDBSQLContainers` in struct `RoutingEndpoints`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iothub/armiothub` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).