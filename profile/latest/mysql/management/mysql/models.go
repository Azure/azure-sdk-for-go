package mysql

import (
	 original "github.com/Azure/azure-sdk-for-go/service/mysql/management/2017-04-30-preview/mysql"
)

type (
	 ManagementClient = original.ManagementClient
	 FirewallRulesClient = original.FirewallRulesClient
	 LogFilesClient = original.LogFilesClient
	 OperationsClient = original.OperationsClient
	 ConfigurationsClient = original.ConfigurationsClient
	 DatabasesClient = original.DatabasesClient
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
	 ServersClient = original.ServersClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 NotSpecified = original.NotSpecified
	 System = original.System
	 User = original.User
	 Disabled = original.Disabled
	 Dropping = original.Dropping
	 Ready = original.Ready
	 FiveFullStopSeven = original.FiveFullStopSeven
	 FiveFullStopSix = original.FiveFullStopSix
	 Basic = original.Basic
	 Standard = original.Standard
	 SslEnforcementEnumDisabled = original.SslEnforcementEnumDisabled
	 SslEnforcementEnumEnabled = original.SslEnforcementEnumEnabled
)
