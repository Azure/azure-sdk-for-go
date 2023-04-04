# Release History

## 1.1.0-beta.2 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.1.0-beta.1 (2022-10-31)
### Features Added

- New const `AlwaysServeDisabled`
- New const `AlwaysServeEnabled`
- New type alias `AlwaysServe`
- New function `PossibleAlwaysServeValues() []AlwaysServe`
- New field `AlwaysServe` in struct `EndpointProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/trafficmanager/armtrafficmanager` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).