# Unreleased

## Breaking Changes

### Signature Changes

#### Struct Fields

1. QuotaProperties.ResourceType changed type from interface{} to ResourceType
1. QuotaRequestOneResourceProperties.ProvisioningState changed type from interface{} to QuotaRequestState
1. QuotaRequestProperties.ProvisioningState changed type from interface{} to QuotaRequestState
1. QuotaRequestStatusDetails.ProvisioningState changed type from interface{} to QuotaRequestState
1. SubRequest.ProvisioningState changed type from interface{} to QuotaRequestState

## Additive Changes

### New Constants

1. QuotaRequestState.QuotaRequestStateAccepted
1. QuotaRequestState.QuotaRequestStateFailed
1. QuotaRequestState.QuotaRequestStateInProgress
1. QuotaRequestState.QuotaRequestStateInvalid
1. QuotaRequestState.QuotaRequestStateSucceeded
1. ResourceType.ResourceTypeDedicated
1. ResourceType.ResourceTypeLowPriority
1. ResourceType.ResourceTypeServiceSpecific
1. ResourceType.ResourceTypeShared
1. ResourceType.ResourceTypeStandard

### New Funcs

1. CurrentQuotaLimitBase.MarshalJSON() ([]byte, error)
1. PossibleQuotaRequestStateValues() []QuotaRequestState
1. PossibleResourceTypeValues() []ResourceType

### Struct Changes

#### New Struct Fields

1. CurrentQuotaLimitBase.ID
1. CurrentQuotaLimitBase.Name
1. CurrentQuotaLimitBase.Type
