package account

import (
	 original "github.com/Azure/azure-sdk-for-go/service/datalake-store/management/2016-11-01/account"
)

type (
	 GroupClient = original.GroupClient
	 ManagementClient = original.ManagementClient
	 FirewallRulesClient = original.FirewallRulesClient
	 DataLakeStoreAccountState = original.DataLakeStoreAccountState
	 DataLakeStoreAccountStatus = original.DataLakeStoreAccountStatus
	 EncryptionConfigType = original.EncryptionConfigType
	 EncryptionProvisioningState = original.EncryptionProvisioningState
	 EncryptionState = original.EncryptionState
	 FirewallAllowAzureIpsState = original.FirewallAllowAzureIpsState
	 FirewallState = original.FirewallState
	 TierType = original.TierType
	 TrustedIDProviderState = original.TrustedIDProviderState
	 DataLakeStoreAccount = original.DataLakeStoreAccount
	 DataLakeStoreAccountBasic = original.DataLakeStoreAccountBasic
	 DataLakeStoreAccountListResult = original.DataLakeStoreAccountListResult
	 DataLakeStoreAccountProperties = original.DataLakeStoreAccountProperties
	 DataLakeStoreAccountPropertiesBasic = original.DataLakeStoreAccountPropertiesBasic
	 DataLakeStoreAccountUpdateParameters = original.DataLakeStoreAccountUpdateParameters
	 DataLakeStoreFirewallRuleListResult = original.DataLakeStoreFirewallRuleListResult
	 DataLakeStoreTrustedIDProviderListResult = original.DataLakeStoreTrustedIDProviderListResult
	 EncryptionConfig = original.EncryptionConfig
	 EncryptionIdentity = original.EncryptionIdentity
	 ErrorDetails = original.ErrorDetails
	 FirewallRule = original.FirewallRule
	 FirewallRuleProperties = original.FirewallRuleProperties
	 KeyVaultMetaInfo = original.KeyVaultMetaInfo
	 Resource = original.Resource
	 SubResource = original.SubResource
	 TrustedIDProvider = original.TrustedIDProvider
	 TrustedIDProviderProperties = original.TrustedIDProviderProperties
	 UpdateDataLakeStoreAccountProperties = original.UpdateDataLakeStoreAccountProperties
	 UpdateEncryptionConfig = original.UpdateEncryptionConfig
	 UpdateFirewallRuleParameters = original.UpdateFirewallRuleParameters
	 UpdateFirewallRuleProperties = original.UpdateFirewallRuleProperties
	 UpdateKeyVaultMetaInfo = original.UpdateKeyVaultMetaInfo
	 UpdateTrustedIDProviderParameters = original.UpdateTrustedIDProviderParameters
	 UpdateTrustedIDProviderProperties = original.UpdateTrustedIDProviderProperties
	 TrustedIDProvidersClient = original.TrustedIDProvidersClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
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
	 ServiceManaged = original.ServiceManaged
	 UserManaged = original.UserManaged
	 EncryptionProvisioningStateCreating = original.EncryptionProvisioningStateCreating
	 EncryptionProvisioningStateSucceeded = original.EncryptionProvisioningStateSucceeded
	 Disabled = original.Disabled
	 Enabled = original.Enabled
	 FirewallAllowAzureIpsStateDisabled = original.FirewallAllowAzureIpsStateDisabled
	 FirewallAllowAzureIpsStateEnabled = original.FirewallAllowAzureIpsStateEnabled
	 FirewallStateDisabled = original.FirewallStateDisabled
	 FirewallStateEnabled = original.FirewallStateEnabled
	 Commitment100TB = original.Commitment100TB
	 Commitment10TB = original.Commitment10TB
	 Commitment1PB = original.Commitment1PB
	 Commitment1TB = original.Commitment1TB
	 Commitment500TB = original.Commitment500TB
	 Commitment5PB = original.Commitment5PB
	 Consumption = original.Consumption
	 TrustedIDProviderStateDisabled = original.TrustedIDProviderStateDisabled
	 TrustedIDProviderStateEnabled = original.TrustedIDProviderStateEnabled
)
