# Unreleased

## Additive Changes

### New Constants

1. DebugDataAccessLevel.All
1. DebugDataAccessLevel.Customer
1. DebugDataAccessLevel.None
1. NestedResourceProvisioningState.NestedResourceProvisioningStateCanceled
1. NestedResourceProvisioningState.NestedResourceProvisioningStateFailed
1. NestedResourceProvisioningState.NestedResourceProvisioningStateSucceeded
1. VirtualNetworkRuleState.VirtualNetworkRuleStateActive
1. VirtualNetworkRuleState.VirtualNetworkRuleStateFailed
1. VirtualNetworkRuleState.VirtualNetworkRuleStateNetworkSourceDeleted

### New Funcs

1. *HiveMetastore.UnmarshalJSON([]byte) error
1. *VirtualNetworkRule.UnmarshalJSON([]byte) error
1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)
1. HiveMetastore.MarshalJSON() ([]byte, error)
1. HiveMetastoreListResult.MarshalJSON() ([]byte, error)
1. HiveMetastoreProperties.MarshalJSON() ([]byte, error)
1. PossibleDebugDataAccessLevelValues() []DebugDataAccessLevel
1. PossibleNestedResourceProvisioningStateValues() []NestedResourceProvisioningState
1. PossibleVirtualNetworkRuleStateValues() []VirtualNetworkRuleState
1. VirtualNetworkRule.MarshalJSON() ([]byte, error)
1. VirtualNetworkRuleListResult.MarshalJSON() ([]byte, error)
1. VirtualNetworkRuleProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ErrorAdditionalInfo
1. ErrorDetail
1. ErrorResponse
1. HiveMetastore
1. HiveMetastoreListResult
1. HiveMetastoreProperties
1. OperationMetaLogSpecification
1. OperationMetaMetricAvailabilitiesSpecification
1. OperationMetaMetricSpecification
1. OperationMetaPropertyInfo
1. OperationMetaServiceSpecification
1. VirtualNetworkRule
1. VirtualNetworkRuleListResult
1. VirtualNetworkRuleProperties

#### New Struct Fields

1. DataLakeAnalyticsAccountProperties.DebugDataAccessLevel
1. DataLakeAnalyticsAccountProperties.HierarchicalQueueState
1. DataLakeAnalyticsAccountProperties.HiveMetastores
1. DataLakeAnalyticsAccountProperties.MaxQueuedJobCountPerUser
1. DataLakeAnalyticsAccountProperties.PublicDataLakeStoreAccounts
1. DataLakeAnalyticsAccountProperties.VirtualNetworkRules
1. Operation.Properties
