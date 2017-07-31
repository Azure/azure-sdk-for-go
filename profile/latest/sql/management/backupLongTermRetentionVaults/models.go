package backuplongtermretentionvaults

import (
	 original "github.com/Azure/azure-sdk-for-go/service/sql/management/2014-04-01/backupLongTermRetentionVaults"
)

type (
	 GroupClient = original.GroupClient
	 ManagementClient = original.ManagementClient
	 BackupLongTermRetentionVault = original.BackupLongTermRetentionVault
	 Properties = original.Properties
	 ProxyResource = original.ProxyResource
	 Resource = original.Resource
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
)
