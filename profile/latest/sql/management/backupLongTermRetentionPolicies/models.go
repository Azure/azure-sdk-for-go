package backuplongtermretentionpolicies

import (
	 original "github.com/Azure/azure-sdk-for-go/service/sql/management/2014-04-01/backupLongTermRetentionPolicies"
)

type (
	 ManagementClient = original.ManagementClient
	 BackupLongTermRetentionPolicyState = original.BackupLongTermRetentionPolicyState
	 BackupLongTermRetentionPolicy = original.BackupLongTermRetentionPolicy
	 BackupLongTermRetentionPolicyProperties = original.BackupLongTermRetentionPolicyProperties
	 ProxyResource = original.ProxyResource
	 Resource = original.Resource
	 GroupClient = original.GroupClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Disabled = original.Disabled
	 Enabled = original.Enabled
)
