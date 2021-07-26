# Unreleased

## Additive Changes

### New Constants

1. AppResourceProvisioningState.Deleting
1. DeploymentResourceProvisioningState.DeploymentResourceProvisioningStateDeleting
1. TrafficDirection.Inbound
1. TrafficDirection.Outbound

### New Funcs

1. PossibleTrafficDirectionValues() []TrafficDirection
1. RequiredTraffic.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. RequiredTraffic

#### New Struct Fields

1. DeploymentResource.Sku
1. MetricDimension.ToBeExportedForShoebox
1. NetworkProfile.RequiredTraffics
