# Unreleased

## Breaking Changes

### Removed Constants

1. VulnerabilityAssessmentPolicyBaselineName.Master

### Signature Changes

#### Const Types

1. Default changed type from VulnerabilityAssessmentPolicyBaselineName to CreateMode

#### Struct Fields

1. SQLPoolResourceProperties.CreateMode changed type from *string to CreateMode

## Additive Changes

### New Constants

1. CreateMode.PointInTimeRestore
1. CreateMode.Recovery
1. CreateMode.Restore
1. NodeSizeFamily.NodeSizeFamilyHardwareAcceleratedFPGA
1. NodeSizeFamily.NodeSizeFamilyHardwareAcceleratedGPU
1. VulnerabilityAssessmentPolicyBaselineName.VulnerabilityAssessmentPolicyBaselineNameDefault
1. VulnerabilityAssessmentPolicyBaselineName.VulnerabilityAssessmentPolicyBaselineNameMaster

### New Funcs

1. IntegrationRuntimesClient.ListOutboundNetworkDependenciesEndpoints(context.Context, string, string, string) (IntegrationRuntimeOutboundNetworkDependenciesEndpointsResponse, error)
1. IntegrationRuntimesClient.ListOutboundNetworkDependenciesEndpointsPreparer(context.Context, string, string, string) (*http.Request, error)
1. IntegrationRuntimesClient.ListOutboundNetworkDependenciesEndpointsResponder(*http.Response) (IntegrationRuntimeOutboundNetworkDependenciesEndpointsResponse, error)
1. IntegrationRuntimesClient.ListOutboundNetworkDependenciesEndpointsSender(*http.Request) (*http.Response, error)
1. PossibleCreateModeValues() []CreateMode

### Struct Changes

#### New Structs

1. IntegrationRuntimeOutboundNetworkDependenciesCategoryEndpoint
1. IntegrationRuntimeOutboundNetworkDependenciesEndpoint
1. IntegrationRuntimeOutboundNetworkDependenciesEndpointDetails
1. IntegrationRuntimeOutboundNetworkDependenciesEndpointsResponse

#### New Struct Fields

1. IntegrationRuntimeVNetProperties.SubnetID
