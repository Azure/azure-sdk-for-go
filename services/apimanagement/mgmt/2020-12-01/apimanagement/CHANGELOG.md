# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. OperationResultContract.ActionLog
1. OperationResultContract.Error
1. OperationResultContract.ResultInfo
1. OperationResultContract.Started
1. OperationResultContract.Status
1. OperationResultContract.Updated
1. TenantConfigurationSyncStateContract.Branch
1. TenantConfigurationSyncStateContract.CommitID
1. TenantConfigurationSyncStateContract.ConfigurationChangeDate
1. TenantConfigurationSyncStateContract.IsExport
1. TenantConfigurationSyncStateContract.IsGitEnabled
1. TenantConfigurationSyncStateContract.IsSynced
1. TenantConfigurationSyncStateContract.SyncDate

## Additive Changes

### New Funcs

1. *OperationResultContract.UnmarshalJSON([]byte) error
1. *TenantConfigurationSyncStateContract.UnmarshalJSON([]byte) error
1. OperationResultContractProperties.MarshalJSON() ([]byte, error)
1. TenantConfigurationSyncStateContract.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. OperationResultContractProperties
1. TenantConfigurationSyncStateContractProperties

#### New Struct Fields

1. OperationResultContract.*OperationResultContractProperties
1. OperationResultContract.Name
1. OperationResultContract.Type
1. TenantConfigurationSyncStateContract.*TenantConfigurationSyncStateContractProperties
