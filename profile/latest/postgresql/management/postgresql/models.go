package postgresql

import (
	 original "github.com/Azure/azure-sdk-for-go/service/postgresql/management/2017-04-30-preview/postgresql"
)

type (
	 ManagementClient = original.ManagementClient
	 OperationsClient = original.OperationsClient
	 ServersClient = original.ServersClient
	 OperationOrigin = original.OperationOrigin
	 ServerState = original.ServerState
	 ServerVersion = original.ServerVersion
	 SkuTier = original.SkuTier
	 SslEnforcementEnum = original.SslEnforcementEnum
	 Configuration = original.Configuration
	 ConfigurationListResult = original.ConfigurationListResult
	 ConfigurationProperties = original.ConfigurationProperties
	 Database = original.Database
	 DatabaseListResult = original.DatabaseListResult
	 DatabaseProperties = original.DatabaseProperties
	 FirewallRule = original.FirewallRule
	 FirewallRuleListResult = original.FirewallRuleListResult
	 FirewallRuleProperties = original.FirewallRuleProperties
	 LogFile = original.LogFile
	 LogFileListResult = original.LogFileListResult
	 LogFileProperties = original.LogFileProperties
	 Operation = original.Operation
	 OperationDisplay = original.OperationDisplay
	 OperationListResult = original.OperationListResult
	 ProxyResource = original.ProxyResource
	 Server = original.Server
	 ServerForCreate = original.ServerForCreate
	 ServerListResult = original.ServerListResult
	 ServerProperties = original.ServerProperties
	 ServerPropertiesForCreate = original.ServerPropertiesForCreate
	 ServerPropertiesForDefaultCreate = original.ServerPropertiesForDefaultCreate
	 ServerPropertiesForRestore = original.ServerPropertiesForRestore
	 ServerUpdateParameters = original.ServerUpdateParameters
	 ServerUpdateParametersProperties = original.ServerUpdateParametersProperties
	 Sku = original.Sku
	 TrackedResource = original.TrackedResource
	 ConfigurationsClient = original.ConfigurationsClient
	 DatabasesClient = original.DatabasesClient
	 FirewallRulesClient = original.FirewallRulesClient
	 LogFilesClient = original.LogFilesClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 NotSpecified = original.NotSpecified
	 System = original.System
	 User = original.User
	 Disabled = original.Disabled
	 Dropping = original.Dropping
	 Ready = original.Ready
	 NineFullStopFive = original.NineFullStopFive
	 NineFullStopSix = original.NineFullStopSix
	 Basic = original.Basic
	 Standard = original.Standard
	 SslEnforcementEnumDisabled = original.SslEnforcementEnumDisabled
	 SslEnforcementEnumEnabled = original.SslEnforcementEnumEnabled
)
