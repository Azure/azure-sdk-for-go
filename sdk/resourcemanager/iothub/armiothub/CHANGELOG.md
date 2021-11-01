# Release History

## 0.2.0 (2021-11-01)
### Breaking Changes

- Function `NewPrivateEndpointConnectionsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOperationsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewResourceProviderCommonClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewIotHubResourceClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateLinkResourcesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCertificatesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewIotHubClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`

### New Content

- New const `RoutingSourceMqttBrokerMessages`
- New const `RoutingSourceDigitalTwinChangeEvents`
- New function `IotHubPropertiesDeviceStreams.MarshalJSON() ([]byte, error)`
- New function `EncryptionPropertiesDescription.MarshalJSON() ([]byte, error)`
- New struct `EncryptionPropertiesDescription`
- New struct `IotHubPropertiesDeviceStreams`
- New struct `KeyVaultKeyProperties`
- New field `DeviceStreams` in struct `IotHubProperties`
- New field `Encryption` in struct `IotHubProperties`

Total 7 breaking change(s), 11 additive change(s).


## 0.1.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.1.0 (2021-10-18)

- Initial preview release.
