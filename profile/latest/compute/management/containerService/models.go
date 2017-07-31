package containerservice

import (
	 original "github.com/Azure/azure-sdk-for-go/service/compute/management/2017-01-31/containerService"
)

type (
	 ManagementClient = original.ManagementClient
	 ContainerServicesClient = original.ContainerServicesClient
	 OrchestratorTypes = original.OrchestratorTypes
	 VMSizeTypes = original.VMSizeTypes
	 AgentPoolProfile = original.AgentPoolProfile
	 CustomProfile = original.CustomProfile
	 DiagnosticsProfile = original.DiagnosticsProfile
	 LinuxProfile = original.LinuxProfile
	 ListResult = original.ListResult
	 MasterProfile = original.MasterProfile
	 Model = original.Model
	 OrchestratorProfile = original.OrchestratorProfile
	 Properties = original.Properties
	 Resource = original.Resource
	 ServicePrincipalProfile = original.ServicePrincipalProfile
	 SSHConfiguration = original.SSHConfiguration
	 SSHPublicKey = original.SSHPublicKey
	 VMDiagnostics = original.VMDiagnostics
	 WindowsProfile = original.WindowsProfile
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Custom = original.Custom
	 DCOS = original.DCOS
	 Kubernetes = original.Kubernetes
	 Swarm = original.Swarm
	 StandardA0 = original.StandardA0
	 StandardA1 = original.StandardA1
	 StandardA10 = original.StandardA10
	 StandardA11 = original.StandardA11
	 StandardA2 = original.StandardA2
	 StandardA3 = original.StandardA3
	 StandardA4 = original.StandardA4
	 StandardA5 = original.StandardA5
	 StandardA6 = original.StandardA6
	 StandardA7 = original.StandardA7
	 StandardA8 = original.StandardA8
	 StandardA9 = original.StandardA9
	 StandardD1 = original.StandardD1
	 StandardD11 = original.StandardD11
	 StandardD11V2 = original.StandardD11V2
	 StandardD12 = original.StandardD12
	 StandardD12V2 = original.StandardD12V2
	 StandardD13 = original.StandardD13
	 StandardD13V2 = original.StandardD13V2
	 StandardD14 = original.StandardD14
	 StandardD14V2 = original.StandardD14V2
	 StandardD1V2 = original.StandardD1V2
	 StandardD2 = original.StandardD2
	 StandardD2V2 = original.StandardD2V2
	 StandardD3 = original.StandardD3
	 StandardD3V2 = original.StandardD3V2
	 StandardD4 = original.StandardD4
	 StandardD4V2 = original.StandardD4V2
	 StandardD5V2 = original.StandardD5V2
	 StandardDS1 = original.StandardDS1
	 StandardDS11 = original.StandardDS11
	 StandardDS12 = original.StandardDS12
	 StandardDS13 = original.StandardDS13
	 StandardDS14 = original.StandardDS14
	 StandardDS2 = original.StandardDS2
	 StandardDS3 = original.StandardDS3
	 StandardDS4 = original.StandardDS4
	 StandardG1 = original.StandardG1
	 StandardG2 = original.StandardG2
	 StandardG3 = original.StandardG3
	 StandardG4 = original.StandardG4
	 StandardG5 = original.StandardG5
	 StandardGS1 = original.StandardGS1
	 StandardGS2 = original.StandardGS2
	 StandardGS3 = original.StandardGS3
	 StandardGS4 = original.StandardGS4
	 StandardGS5 = original.StandardGS5
)
