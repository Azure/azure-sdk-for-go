# Release History

## 1.2.0-beta.2 (2023-09-22)
### Breaking Changes

- Field `CollectionName` of struct `RoutingCosmosDBSQLAPIProperties` has been removed
- Field `CosmosDBSQLCollections` of struct `RoutingEndpoints` has been removed

### Features Added

- New field `ContainerName` in struct `RoutingCosmosDBSQLAPIProperties`
- New field `CosmosDBSQLContainers` in struct `RoutingEndpoints`


## 1.2.0-beta.1 (2023-04-28)
### Features Added

- New value `RoutingSourceDigitalTwinChangeEvents`, `RoutingSourceMqttBrokerMessages` added to enum type `RoutingSource`
- New enum type `IPVersion` with values `IPVersionIPv4`, `IPVersionIPv4IPv6`, `IPVersionIPv6`
- New struct `EncryptionPropertiesDescription`
- New struct `KeyVaultKeyProperties`
- New struct `PropertiesDeviceStreams`
- New struct `RootCertificateProperties`
- New struct `RoutingCosmosDBSQLAPIProperties`
- New field `DeviceStreams` in struct `Properties`
- New field `Encryption` in struct `Properties`
- New field `IPVersion` in struct `Properties`
- New field `RootCertificate` in struct `Properties`
- New field `CosmosDBSQLCollections` in struct `RoutingEndpoints`


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