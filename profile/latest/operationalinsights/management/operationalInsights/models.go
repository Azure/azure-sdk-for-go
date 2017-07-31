package operationalinsights

import (
	 original "github.com/Azure/azure-sdk-for-go/service/operationalinsights/management/2015-11-01-preview/operationalInsights"
)

type (
	 DataSourcesClient = original.DataSourcesClient
	 LinkedServicesClient = original.LinkedServicesClient
	 DataSourceKind = original.DataSourceKind
	 EntityStatus = original.EntityStatus
	 SkuNameEnum = original.SkuNameEnum
	 DataSource = original.DataSource
	 DataSourceFilter = original.DataSourceFilter
	 DataSourceListResult = original.DataSourceListResult
	 IntelligencePack = original.IntelligencePack
	 LinkedService = original.LinkedService
	 LinkedServiceListResult = original.LinkedServiceListResult
	 LinkedServiceProperties = original.LinkedServiceProperties
	 ListIntelligencePack = original.ListIntelligencePack
	 ManagementGroup = original.ManagementGroup
	 ManagementGroupProperties = original.ManagementGroupProperties
	 MetricName = original.MetricName
	 ProxyResource = original.ProxyResource
	 Resource = original.Resource
	 SharedKeys = original.SharedKeys
	 Sku = original.Sku
	 UsageMetric = original.UsageMetric
	 Workspace = original.Workspace
	 WorkspaceListManagementGroupsResult = original.WorkspaceListManagementGroupsResult
	 WorkspaceListResult = original.WorkspaceListResult
	 WorkspaceListUsagesResult = original.WorkspaceListUsagesResult
	 WorkspaceProperties = original.WorkspaceProperties
	 WorkspacesClient = original.WorkspacesClient
	 ManagementClient = original.ManagementClient
)

const (
	 AzureActivityLog = original.AzureActivityLog
	 ChangeTrackingCustomRegistry = original.ChangeTrackingCustomRegistry
	 ChangeTrackingDefaultPath = original.ChangeTrackingDefaultPath
	 ChangeTrackingDefaultRegistry = original.ChangeTrackingDefaultRegistry
	 ChangeTrackingPath = original.ChangeTrackingPath
	 CustomLog = original.CustomLog
	 CustomLogCollection = original.CustomLogCollection
	 GenericDataSource = original.GenericDataSource
	 IISLogs = original.IISLogs
	 LinuxPerformanceCollection = original.LinuxPerformanceCollection
	 LinuxPerformanceObject = original.LinuxPerformanceObject
	 LinuxSyslog = original.LinuxSyslog
	 LinuxSyslogCollection = original.LinuxSyslogCollection
	 WindowsEvent = original.WindowsEvent
	 WindowsPerformanceCounter = original.WindowsPerformanceCounter
	 Canceled = original.Canceled
	 Creating = original.Creating
	 Deleting = original.Deleting
	 Failed = original.Failed
	 ProvisioningAccount = original.ProvisioningAccount
	 Succeeded = original.Succeeded
	 Free = original.Free
	 PerNode = original.PerNode
	 Premium = original.Premium
	 Standalone = original.Standalone
	 Standard = original.Standard
	 Unlimited = original.Unlimited
	 DefaultBaseURI = original.DefaultBaseURI
)
