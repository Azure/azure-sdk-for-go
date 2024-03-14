# Release History

## 0.8.0 (2024-03-22)
### Features Added

- New enum type `KeyType` with values `KeyTypePrimary`, `KeyTypeSecondary`
- New function `*WorkspaceClient.ListKeys(context.Context, string, string, *WorkspaceClientListKeysOptions) (WorkspaceClientListKeysResponse, error)`
- New function `*WorkspaceClient.RegenerateKeys(context.Context, string, string, APIKeys, *WorkspaceClientRegenerateKeysOptions) (WorkspaceClientRegenerateKeysResponse, error)`
- New struct `APIKey`
- New struct `APIKeys`
- New struct `ListKeysResult`
- New field `APIKeyEnabled` in struct `WorkspaceResourceProperties`


## 0.7.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.6.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.6.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quantum/armquantum` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).