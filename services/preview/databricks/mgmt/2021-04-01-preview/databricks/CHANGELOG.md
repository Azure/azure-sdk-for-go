# Unreleased

## Additive Changes

### New Funcs

1. NewOutboundNetworkDependenciesEndpointsClient(string) OutboundNetworkDependenciesEndpointsClient
1. NewOutboundNetworkDependenciesEndpointsClientWithBaseURI(string, string) OutboundNetworkDependenciesEndpointsClient
1. OutboundNetworkDependenciesEndpointsClient.List(context.Context, string, string) (ListOutboundEnvironmentEndpoint, error)
1. OutboundNetworkDependenciesEndpointsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. OutboundNetworkDependenciesEndpointsClient.ListResponder(*http.Response) (ListOutboundEnvironmentEndpoint, error)
1. OutboundNetworkDependenciesEndpointsClient.ListSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. EndpointDependency
1. EndpointDetail
1. ListOutboundEnvironmentEndpoint
1. OutboundEnvironmentEndpoint
1. OutboundNetworkDependenciesEndpointsClient
