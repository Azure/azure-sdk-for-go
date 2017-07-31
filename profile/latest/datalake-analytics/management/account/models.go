package account

import (
	 original "github.com/Azure/azure-sdk-for-go/service/datalake-analytics/management/2016-11-01/account"
)

type (
	 DataLakeStoreAccountsClient = original.DataLakeStoreAccountsClient
	 FirewallRulesClient = original.FirewallRulesClient
	 AADObjectType = original.AADObjectType
	 DataLakeAnalyticsAccountState = original.DataLakeAnalyticsAccountState
	 DataLakeAnalyticsAccountStatus = original.DataLakeAnalyticsAccountStatus
	 FirewallAllowAzureIpsState = original.FirewallAllowAzureIpsState
	 FirewallState = original.FirewallState
	 TierType = original.TierType
	 AddDataLakeStoreParameters = original.AddDataLakeStoreParameters
	 AddStorageAccountParameters = original.AddStorageAccountParameters
	 ComputePolicy = original.ComputePolicy
	 ComputePolicyAccountCreateParameters = original.ComputePolicyAccountCreateParameters
	 ComputePolicyCreateOrUpdateParameters = original.ComputePolicyCreateOrUpdateParameters
	 ComputePolicyListResult = original.ComputePolicyListResult
	 ComputePolicyProperties = original.ComputePolicyProperties
	 ComputePolicyPropertiesCreateParameters = original.ComputePolicyPropertiesCreateParameters
	 DataLakeAnalyticsAccount = original.DataLakeAnalyticsAccount
	 DataLakeAnalyticsAccountBasic = original.DataLakeAnalyticsAccountBasic
	 DataLakeAnalyticsAccountListDataLakeStoreResult = original.DataLakeAnalyticsAccountListDataLakeStoreResult
	 DataLakeAnalyticsAccountListResult = original.DataLakeAnalyticsAccountListResult
	 DataLakeAnalyticsAccountListStorageAccountsResult = original.DataLakeAnalyticsAccountListStorageAccountsResult
	 DataLakeAnalyticsAccountProperties = original.DataLakeAnalyticsAccountProperties
	 DataLakeAnalyticsAccountPropertiesBasic = original.DataLakeAnalyticsAccountPropertiesBasic
	 DataLakeAnalyticsAccountUpdateParameters = original.DataLakeAnalyticsAccountUpdateParameters
	 DataLakeAnalyticsFirewallRuleListResult = original.DataLakeAnalyticsFirewallRuleListResult
	 DataLakeStoreAccountInfo = original.DataLakeStoreAccountInfo
	 DataLakeStoreAccountInfoProperties = original.DataLakeStoreAccountInfoProperties
	 FirewallRule = original.FirewallRule
	 FirewallRuleProperties = original.FirewallRuleProperties
	 ListSasTokensResult = original.ListSasTokensResult
	 ListStorageContainersResult = original.ListStorageContainersResult
	 OptionalSubResource = original.OptionalSubResource
	 Resource = original.Resource
	 SasTokenInfo = original.SasTokenInfo
	 StorageAccountInfo = original.StorageAccountInfo
	 StorageAccountProperties = original.StorageAccountProperties
	 StorageContainer = original.StorageContainer
	 StorageContainerProperties = original.StorageContainerProperties
	 SubResource = original.SubResource
	 UpdateDataLakeAnalyticsAccountProperties = original.UpdateDataLakeAnalyticsAccountProperties
	 UpdateFirewallRuleParameters = original.UpdateFirewallRuleParameters
	 UpdateFirewallRuleProperties = original.UpdateFirewallRuleProperties
	 UpdateStorageAccountParameters = original.UpdateStorageAccountParameters
	 UpdateStorageAccountProperties = original.UpdateStorageAccountProperties
	 StorageAccountsClient = original.StorageAccountsClient
	 GroupClient = original.GroupClient
	 ManagementClient = original.ManagementClient
	 ComputePoliciesClient = original.ComputePoliciesClient
)

const (
	 Group = original.Group
	 ServicePrincipal = original.ServicePrincipal
	 User = original.User
	 Active = original.Active
	 Suspended = original.Suspended
	 Creating = original.Creating
	 Deleted = original.Deleted
	 Deleting = original.Deleting
	 Failed = original.Failed
	 Patching = original.Patching
	 Resuming = original.Resuming
	 Running = original.Running
	 Succeeded = original.Succeeded
	 Suspending = original.Suspending
	 Disabled = original.Disabled
	 Enabled = original.Enabled
	 FirewallStateDisabled = original.FirewallStateDisabled
	 FirewallStateEnabled = original.FirewallStateEnabled
	 Commitment100000AUHours = original.Commitment100000AUHours
	 Commitment10000AUHours = original.Commitment10000AUHours
	 Commitment1000AUHours = original.Commitment1000AUHours
	 Commitment100AUHours = original.Commitment100AUHours
	 Commitment500000AUHours = original.Commitment500000AUHours
	 Commitment50000AUHours = original.Commitment50000AUHours
	 Commitment5000AUHours = original.Commitment5000AUHours
	 Commitment500AUHours = original.Commitment500AUHours
	 Consumption = original.Consumption
	 DefaultBaseURI = original.DefaultBaseURI
)
