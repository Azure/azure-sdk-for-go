# Release History

## 2.0.0 (2026-02-11)
### Breaking Changes

- Function `*ClientFactory.NewMicrosoftSerialConsoleClient` has been removed
- Function `NewMicrosoftSerialConsoleClient` has been removed
- Function `*MicrosoftSerialConsoleClient.DisableConsole` has been removed
- Function `*MicrosoftSerialConsoleClient.EnableConsole` has been removed
- Function `*MicrosoftSerialConsoleClient.GetConsoleStatus` has been removed
- Function `*MicrosoftSerialConsoleClient.ListOperations` has been removed
- Function `*SerialPortsClient.Delete` has been removed
- Struct `GetSerialConsoleSubscriptionNotFound` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Field `Disabled` of struct `DisableSerialConsoleResult` has been removed
- Field `Disabled` of struct `EnableSerialConsoleResult` has been removed
- Field `Disabled` of struct `Status` has been removed

### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `SerialPortConnectionState` with values `SerialPortConnectionStateActive`, `SerialPortConnectionStateInactive`
- New function `NewClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*Client, error)`
- New function `*Client.DisableConsole(ctx context.Context, defaultParam string, options *ClientDisableConsoleOptions) (ClientDisableConsoleResponse, error)`
- New function `*Client.EnableConsole(ctx context.Context, defaultParam string, options *ClientEnableConsoleOptions) (ClientEnableConsoleResponse, error)`
- New function `*Client.GetConsoleStatus(ctx context.Context, defaultParam string, options *ClientGetConsoleStatusOptions) (ClientGetConsoleStatusResponse, error)`
- New function `*Client.ListOperations(ctx context.Context, options *ClientListOperationsOptions) (ClientListOperationsResponse, error)`
- New function `*Client.NewSerialPortsClient() *SerialPortsClient`
- New function `*ClientFactory.NewClient() *Client`
- New struct `DisableSerialConsoleResultProperties`
- New struct `EnableSerialConsoleResultProperties`
- New struct `StatusProperties`
- New struct `SystemData`
- New field `Properties` in struct `DisableSerialConsoleResult`
- New field `Properties` in struct `EnableSerialConsoleResult`
- New field `SystemData` in struct `SerialPort`
- New field `ConnectionState` in struct `SerialPortProperties`
- New field `Properties` in struct `Status`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/serialconsole/armserialconsole` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).