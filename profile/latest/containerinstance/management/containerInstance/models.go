package containerinstance

import (
	 original "github.com/Azure/azure-sdk-for-go/service/containerinstance/management/2017-08-01-preview/containerInstance"
)

type (
	 ManagementClient = original.ManagementClient
	 ContainerGroupsClient = original.ContainerGroupsClient
	 ContainerLogsClient = original.ContainerLogsClient
	 ContainerGroupNetworkProtocol = original.ContainerGroupNetworkProtocol
	 ContainerRestartPolicy = original.ContainerRestartPolicy
	 OperatingSystemTypes = original.OperatingSystemTypes
	 AzureFileVolume = original.AzureFileVolume
	 Container = original.Container
	 ContainerEvent = original.ContainerEvent
	 ContainerGroup = original.ContainerGroup
	 ContainerGroupProperties = original.ContainerGroupProperties
	 ContainerGroupListResult = original.ContainerGroupListResult
	 ContainerPort = original.ContainerPort
	 ContainerProperties = original.ContainerProperties
	 ContainerPropertiesInstanceView = original.ContainerPropertiesInstanceView
	 ContainerState = original.ContainerState
	 EnvironmentVariable = original.EnvironmentVariable
	 ImageRegistryCredential = original.ImageRegistryCredential
	 IPAddress = original.IPAddress
	 Logs = original.Logs
	 Port = original.Port
	 Resource = original.Resource
	 ResourceLimits = original.ResourceLimits
	 ResourceRequests = original.ResourceRequests
	 ResourceRequirements = original.ResourceRequirements
	 Volume = original.Volume
	 VolumeMount = original.VolumeMount
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 TCP = original.TCP
	 UDP = original.UDP
	 Always = original.Always
	 Linux = original.Linux
	 Windows = original.Windows
)
