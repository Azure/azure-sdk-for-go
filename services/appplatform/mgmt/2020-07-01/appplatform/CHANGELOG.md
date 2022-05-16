# Unreleased

## Additive Changes

### New Constants

1. ActionType.Internal
1. TrafficDirection.Inbound
1. TrafficDirection.Outbound

### New Funcs

1. OperationDetail.MarshalJSON() ([]byte, error)
1. PossibleActionTypeValues() []ActionType
1. PossibleTrafficDirectionValues() []TrafficDirection
1. RequiredTraffic.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. RequiredTraffic

#### New Struct Fields

1. MetricDimension.ToBeExportedForShoebox
1. MetricSpecification.SourceMdmNamespace
1. NetworkProfile.RequiredTraffics
1. OperationDetail.ActionType
