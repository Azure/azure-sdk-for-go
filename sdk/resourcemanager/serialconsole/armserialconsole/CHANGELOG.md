# Release History

## 1.0.0 (2022-05-17)
### Breaking Changes

- Function `Operations.MarshalJSON` has been removed
- Function `SerialPortListResult.MarshalJSON` has been removed
- Function `CloudErrorBody.MarshalJSON` has been removed


## 0.3.0 (2022-04-13)
### Breaking Changes

- Function `NewSerialPortsClient` return value(s) have been changed from `(*SerialPortsClient)` to `(*SerialPortsClient, error)`
- Function `NewMicrosoftSerialConsoleClient` return value(s) have been changed from `(*MicrosoftSerialConsoleClient)` to `(*MicrosoftSerialConsoleClient, error)`
- Function `SerialPortState.ToPtr` has been removed
- Struct `MicrosoftSerialConsoleClientListOperationsResult` has been removed
- Struct `SerialPortsClientConnectResult` has been removed
- Struct `SerialPortsClientCreateResult` has been removed
- Struct `SerialPortsClientGetResult` has been removed
- Struct `SerialPortsClientListBySubscriptionsResult` has been removed
- Struct `SerialPortsClientListResult` has been removed
- Field `SerialPortsClientListBySubscriptionsResult` of struct `SerialPortsClientListBySubscriptionsResponse` has been removed
- Field `RawResponse` of struct `SerialPortsClientListBySubscriptionsResponse` has been removed
- Field `SerialPortsClientGetResult` of struct `SerialPortsClientGetResponse` has been removed
- Field `RawResponse` of struct `SerialPortsClientGetResponse` has been removed
- Field `RawResponse` of struct `MicrosoftSerialConsoleClientGetConsoleStatusResponse` has been removed
- Field `SerialPortsClientListResult` of struct `SerialPortsClientListResponse` has been removed
- Field `RawResponse` of struct `SerialPortsClientListResponse` has been removed
- Field `MicrosoftSerialConsoleClientListOperationsResult` of struct `MicrosoftSerialConsoleClientListOperationsResponse` has been removed
- Field `RawResponse` of struct `MicrosoftSerialConsoleClientListOperationsResponse` has been removed
- Field `SerialPortsClientCreateResult` of struct `SerialPortsClientCreateResponse` has been removed
- Field `RawResponse` of struct `SerialPortsClientCreateResponse` has been removed
- Field `RawResponse` of struct `SerialPortsClientDeleteResponse` has been removed
- Field `SerialPortsClientConnectResult` of struct `SerialPortsClientConnectResponse` has been removed
- Field `RawResponse` of struct `SerialPortsClientConnectResponse` has been removed
- Field `RawResponse` of struct `MicrosoftSerialConsoleClientEnableConsoleResponse` has been removed
- Field `RawResponse` of struct `MicrosoftSerialConsoleClientDisableConsoleResponse` has been removed

### Features Added

- New anonymous field `SerialPort` in struct `SerialPortsClientGetResponse`
- New anonymous field `SerialPortListResult` in struct `SerialPortsClientListBySubscriptionsResponse`
- New anonymous field `SerialPortListResult` in struct `SerialPortsClientListResponse`
- New anonymous field `SerialPort` in struct `SerialPortsClientCreateResponse`
- New anonymous field `SerialPortConnectResult` in struct `SerialPortsClientConnectResponse`
- New anonymous field `Operations` in struct `MicrosoftSerialConsoleClientListOperationsResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*SerialPortsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, string, string, *SerialPortsGetOptions)` to `(context.Context, string, string, string, string, string, *SerialPortsClientGetOptions)`
- Function `*SerialPortsClient.Get` return value(s) have been changed from `(SerialPortsGetResponse, error)` to `(SerialPortsClientGetResponse, error)`
- Function `*SerialPortsClient.List` parameter(s) have been changed from `(context.Context, string, string, string, string, *SerialPortsListOptions)` to `(context.Context, string, string, string, string, *SerialPortsClientListOptions)`
- Function `*SerialPortsClient.List` return value(s) have been changed from `(SerialPortsListResponse, error)` to `(SerialPortsClientListResponse, error)`
- Function `*SerialPortsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, string, string, string, *SerialPortsDeleteOptions)` to `(context.Context, string, string, string, string, string, *SerialPortsClientDeleteOptions)`
- Function `*SerialPortsClient.Delete` return value(s) have been changed from `(SerialPortsDeleteResponse, error)` to `(SerialPortsClientDeleteResponse, error)`
- Function `*SerialPortsClient.Connect` parameter(s) have been changed from `(context.Context, string, string, string, string, string, *SerialPortsConnectOptions)` to `(context.Context, string, string, string, string, string, *SerialPortsClientConnectOptions)`
- Function `*SerialPortsClient.Connect` return value(s) have been changed from `(SerialPortsConnectResponse, error)` to `(SerialPortsClientConnectResponse, error)`
- Function `*SerialPortsClient.ListBySubscriptions` parameter(s) have been changed from `(context.Context, *SerialPortsListBySubscriptionsOptions)` to `(context.Context, *SerialPortsClientListBySubscriptionsOptions)`
- Function `*SerialPortsClient.ListBySubscriptions` return value(s) have been changed from `(SerialPortsListBySubscriptionsResponse, error)` to `(SerialPortsClientListBySubscriptionsResponse, error)`
- Function `*SerialPortsClient.Create` parameter(s) have been changed from `(context.Context, string, string, string, string, string, SerialPort, *SerialPortsCreateOptions)` to `(context.Context, string, string, string, string, string, SerialPort, *SerialPortsClientCreateOptions)`
- Function `*SerialPortsClient.Create` return value(s) have been changed from `(SerialPortsCreateResponse, error)` to `(SerialPortsClientCreateResponse, error)`
- Function `CloudError.Error` has been removed
- Function `SerialConsoleOperations.MarshalJSON` has been removed
- Struct `SerialConsoleOperations` has been removed
- Struct `SerialConsoleOperationsValueItem` has been removed
- Struct `SerialConsoleOperationsValueItemDisplay` has been removed
- Struct `SerialConsoleStatus` has been removed
- Struct `SerialPortsConnectOptions` has been removed
- Struct `SerialPortsConnectResponse` has been removed
- Struct `SerialPortsConnectResult` has been removed
- Struct `SerialPortsCreateOptions` has been removed
- Struct `SerialPortsCreateResponse` has been removed
- Struct `SerialPortsCreateResult` has been removed
- Struct `SerialPortsDeleteOptions` has been removed
- Struct `SerialPortsDeleteResponse` has been removed
- Struct `SerialPortsGetOptions` has been removed
- Struct `SerialPortsGetResponse` has been removed
- Struct `SerialPortsGetResult` has been removed
- Struct `SerialPortsListBySubscriptionsOptions` has been removed
- Struct `SerialPortsListBySubscriptionsResponse` has been removed
- Struct `SerialPortsListBySubscriptionsResult` has been removed
- Struct `SerialPortsListOptions` has been removed
- Struct `SerialPortsListResponse` has been removed
- Struct `SerialPortsListResult` has been removed
- Field `Resource` of struct `ProxyResource` has been removed
- Field `ProxyResource` of struct `SerialPort` has been removed
- Field `InnerError` of struct `CloudError` has been removed
- Field `SerialConsoleOperations` of struct `MicrosoftSerialConsoleClientListOperationsResult` has been removed

### Features Added

- New function `Operations.MarshalJSON() ([]byte, error)`
- New struct `Operations`
- New struct `OperationsValueItem`
- New struct `OperationsValueItemDisplay`
- New struct `SerialPortsClientConnectOptions`
- New struct `SerialPortsClientConnectResponse`
- New struct `SerialPortsClientConnectResult`
- New struct `SerialPortsClientCreateOptions`
- New struct `SerialPortsClientCreateResponse`
- New struct `SerialPortsClientCreateResult`
- New struct `SerialPortsClientDeleteOptions`
- New struct `SerialPortsClientDeleteResponse`
- New struct `SerialPortsClientGetOptions`
- New struct `SerialPortsClientGetResponse`
- New struct `SerialPortsClientGetResult`
- New struct `SerialPortsClientListBySubscriptionsOptions`
- New struct `SerialPortsClientListBySubscriptionsResponse`
- New struct `SerialPortsClientListBySubscriptionsResult`
- New struct `SerialPortsClientListOptions`
- New struct `SerialPortsClientListResponse`
- New struct `SerialPortsClientListResult`
- New struct `Status`
- New anonymous field `Operations` in struct `MicrosoftSerialConsoleClientListOperationsResult`
- New field `Error` in struct `CloudError`
- New field `ID` in struct `SerialPort`
- New field `Name` in struct `SerialPort`
- New field `Type` in struct `SerialPort`
- New field `Type` in struct `ProxyResource`
- New field `ID` in struct `ProxyResource`
- New field `Name` in struct `ProxyResource`


## 0.1.0 (2021-12-09)

- Init release.
