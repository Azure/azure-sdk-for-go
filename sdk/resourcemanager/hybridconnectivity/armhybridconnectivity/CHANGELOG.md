# Release History

## 1.0.0 (2022-12-22)
### Features Added

- New type alias `ServiceName` with values `ServiceNameSSH`, `ServiceNameWAC`
- New function `*EndpointsClient.ListIngressGatewayCredentials(context.Context, string, string, *EndpointsClientListIngressGatewayCredentialsOptions) (EndpointsClientListIngressGatewayCredentialsResponse, error)`
- New struct `ListCredentialsRequest`
- New struct `ServiceConfiguration`
- New field `ServiceConfigurations` in struct `EndpointProperties`
- New field `ListCredentialsRequest` in struct `EndpointsClientListCredentialsOptions`
- New field `ServiceName` in struct `ManagedProxyRequest`
- New field `ServiceConfigurationToken` in struct `RelayNamespaceAccessProperties`


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridconnectivity/armhybridconnectivity` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).