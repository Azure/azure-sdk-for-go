package iothub

import (
	 original "github.com/Azure/azure-sdk-for-go/service/iothub/management/2017-01-19/iothub"
)

type (
	 ManagementClient = original.ManagementClient
	 AccessRights = original.AccessRights
	 Capabilities = original.Capabilities
	 IPFilterActionType = original.IPFilterActionType
	 JobStatus = original.JobStatus
	 JobType = original.JobType
	 NameUnavailabilityReason = original.NameUnavailabilityReason
	 OperationMonitoringLevel = original.OperationMonitoringLevel
	 RoutingSource = original.RoutingSource
	 ScaleType = original.ScaleType
	 Sku = original.Sku
	 SkuTier = original.SkuTier
	 Capacity = original.Capacity
	 CloudToDeviceProperties = original.CloudToDeviceProperties
	 Description = original.Description
	 DescriptionListResult = original.DescriptionListResult
	 ErrorDetails = original.ErrorDetails
	 EventHubConsumerGroupInfo = original.EventHubConsumerGroupInfo
	 EventHubConsumerGroupsListResult = original.EventHubConsumerGroupsListResult
	 EventHubProperties = original.EventHubProperties
	 ExportDevicesRequest = original.ExportDevicesRequest
	 FallbackRouteProperties = original.FallbackRouteProperties
	 FeedbackProperties = original.FeedbackProperties
	 ImportDevicesRequest = original.ImportDevicesRequest
	 IPFilterRule = original.IPFilterRule
	 JobResponse = original.JobResponse
	 JobResponseListResult = original.JobResponseListResult
	 MessagingEndpointProperties = original.MessagingEndpointProperties
	 NameAvailabilityInfo = original.NameAvailabilityInfo
	 OperationInputs = original.OperationInputs
	 OperationsMonitoringProperties = original.OperationsMonitoringProperties
	 Properties = original.Properties
	 QuotaMetricInfo = original.QuotaMetricInfo
	 QuotaMetricInfoListResult = original.QuotaMetricInfoListResult
	 RegistryStatistics = original.RegistryStatistics
	 Resource = original.Resource
	 RouteProperties = original.RouteProperties
	 RoutingEndpoints = original.RoutingEndpoints
	 RoutingEventHubProperties = original.RoutingEventHubProperties
	 RoutingProperties = original.RoutingProperties
	 RoutingServiceBusQueueEndpointProperties = original.RoutingServiceBusQueueEndpointProperties
	 RoutingServiceBusTopicEndpointProperties = original.RoutingServiceBusTopicEndpointProperties
	 SetObject = original.SetObject
	 SharedAccessSignatureAuthorizationRule = original.SharedAccessSignatureAuthorizationRule
	 SharedAccessSignatureAuthorizationRuleListResult = original.SharedAccessSignatureAuthorizationRuleListResult
	 SkuDescription = original.SkuDescription
	 SkuDescriptionListResult = original.SkuDescriptionListResult
	 SkuInfo = original.SkuInfo
	 StorageEndpointProperties = original.StorageEndpointProperties
	 ResourceClient = original.ResourceClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 DeviceConnect = original.DeviceConnect
	 RegistryRead = original.RegistryRead
	 RegistryReadDeviceConnect = original.RegistryReadDeviceConnect
	 RegistryReadRegistryWrite = original.RegistryReadRegistryWrite
	 RegistryReadRegistryWriteDeviceConnect = original.RegistryReadRegistryWriteDeviceConnect
	 RegistryReadRegistryWriteServiceConnect = original.RegistryReadRegistryWriteServiceConnect
	 RegistryReadRegistryWriteServiceConnectDeviceConnect = original.RegistryReadRegistryWriteServiceConnectDeviceConnect
	 RegistryReadServiceConnect = original.RegistryReadServiceConnect
	 RegistryReadServiceConnectDeviceConnect = original.RegistryReadServiceConnectDeviceConnect
	 RegistryWrite = original.RegistryWrite
	 RegistryWriteDeviceConnect = original.RegistryWriteDeviceConnect
	 RegistryWriteServiceConnect = original.RegistryWriteServiceConnect
	 RegistryWriteServiceConnectDeviceConnect = original.RegistryWriteServiceConnectDeviceConnect
	 ServiceConnect = original.ServiceConnect
	 ServiceConnectDeviceConnect = original.ServiceConnectDeviceConnect
	 DeviceManagement = original.DeviceManagement
	 None = original.None
	 Accept = original.Accept
	 Reject = original.Reject
	 Cancelled = original.Cancelled
	 Completed = original.Completed
	 Enqueued = original.Enqueued
	 Failed = original.Failed
	 Running = original.Running
	 Unknown = original.Unknown
	 JobTypeBackup = original.JobTypeBackup
	 JobTypeExport = original.JobTypeExport
	 JobTypeFactoryResetDevice = original.JobTypeFactoryResetDevice
	 JobTypeFirmwareUpdate = original.JobTypeFirmwareUpdate
	 JobTypeImport = original.JobTypeImport
	 JobTypeReadDeviceProperties = original.JobTypeReadDeviceProperties
	 JobTypeRebootDevice = original.JobTypeRebootDevice
	 JobTypeUnknown = original.JobTypeUnknown
	 JobTypeUpdateDeviceConfiguration = original.JobTypeUpdateDeviceConfiguration
	 JobTypeWriteDeviceProperties = original.JobTypeWriteDeviceProperties
	 AlreadyExists = original.AlreadyExists
	 Invalid = original.Invalid
	 OperationMonitoringLevelError = original.OperationMonitoringLevelError
	 OperationMonitoringLevelErrorInformation = original.OperationMonitoringLevelErrorInformation
	 OperationMonitoringLevelInformation = original.OperationMonitoringLevelInformation
	 OperationMonitoringLevelNone = original.OperationMonitoringLevelNone
	 DeviceJobLifecycleEvents = original.DeviceJobLifecycleEvents
	 DeviceLifecycleEvents = original.DeviceLifecycleEvents
	 DeviceMessages = original.DeviceMessages
	 TwinChangeEvents = original.TwinChangeEvents
	 ScaleTypeAutomatic = original.ScaleTypeAutomatic
	 ScaleTypeManual = original.ScaleTypeManual
	 ScaleTypeNone = original.ScaleTypeNone
	 F1 = original.F1
	 S1 = original.S1
	 S2 = original.S2
	 S3 = original.S3
	 Free = original.Free
	 Standard = original.Standard
)
