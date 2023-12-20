# Release History

## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.0.0 (2023-08-25)
### Breaking Changes

- Enum `BareMetalMachineHardwareValidationCategory` has been removed
- Function `*BareMetalMachinesClient.BeginValidateHardware` has been removed
- Function `*StorageAppliancesClient.BeginRunReadCommands` has been removed
- Function `*VirtualMachinesClient.BeginAttachVolume` has been removed
- Function `*VirtualMachinesClient.BeginDetachVolume` has been removed
- Struct `BareMetalMachineValidateHardwareParameters` has been removed
- Struct `StorageApplianceCommandSpecification` has been removed
- Struct `StorageApplianceRunReadCommandsParameters` has been removed
- Struct `VirtualMachineVolumeParameters` has been removed

### Features Added

- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `OperationStatusResult`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientCordonResponse`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientPowerOffResponse`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientReimageResponse`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientReplaceResponse`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientRestartResponse`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientRunCommandResponse`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientRunDataExtractsResponse`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientRunReadCommandsResponse`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientStartResponse`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientUncordonResponse`
- New anonymous field `OperationStatusResult` in struct `ClustersClientDeployResponse`
- New anonymous field `OperationStatusResult` in struct `ClustersClientUpdateVersionResponse`
- New anonymous field `OperationStatusResult` in struct `KubernetesClustersClientRestartNodeResponse`
- New anonymous field `OperationStatusResult` in struct `StorageAppliancesClientDisableRemoteVendorManagementResponse`
- New anonymous field `OperationStatusResult` in struct `StorageAppliancesClientEnableRemoteVendorManagementResponse`
- New anonymous field `OperationStatusResult` in struct `VirtualMachinesClientPowerOffResponse`
- New anonymous field `OperationStatusResult` in struct `VirtualMachinesClientReimageResponse`
- New anonymous field `OperationStatusResult` in struct `VirtualMachinesClientRestartResponse`
- New anonymous field `OperationStatusResult` in struct `VirtualMachinesClientStartResponse`


## 0.2.0 (2023-07-28)
### Breaking Changes

- Enum `DefaultCniNetworkDetailedStatus` has been removed
- Enum `DefaultCniNetworkProvisioningState` has been removed
- Enum `HybridAksClusterDetailedStatus` has been removed
- Enum `HybridAksClusterMachinePowerState` has been removed
- Enum `HybridAksClusterProvisioningState` has been removed
- Enum `StorageApplianceHardwareValidationCategory` has been removed
- Function `*BareMetalMachineKeySetsClient.NewListByResourceGroupPager` has been removed
- Function `*BmcKeySetsClient.NewListByResourceGroupPager` has been removed
- Function `*ClientFactory.NewDefaultCniNetworksClient` has been removed
- Function `*ClientFactory.NewHybridAksClustersClient` has been removed
- Function `*ConsolesClient.NewListByResourceGroupPager` has been removed
- Function `NewDefaultCniNetworksClient` has been removed
- Function `*DefaultCniNetworksClient.BeginCreateOrUpdate` has been removed
- Function `*DefaultCniNetworksClient.BeginDelete` has been removed
- Function `*DefaultCniNetworksClient.Get` has been removed
- Function `*DefaultCniNetworksClient.NewListByResourceGroupPager` has been removed
- Function `*DefaultCniNetworksClient.NewListBySubscriptionPager` has been removed
- Function `*DefaultCniNetworksClient.Update` has been removed
- Function `NewHybridAksClustersClient` has been removed
- Function `*HybridAksClustersClient.BeginCreateOrUpdate` has been removed
- Function `*HybridAksClustersClient.BeginDelete` has been removed
- Function `*HybridAksClustersClient.Get` has been removed
- Function `*HybridAksClustersClient.NewListByResourceGroupPager` has been removed
- Function `*HybridAksClustersClient.NewListBySubscriptionPager` has been removed
- Function `*HybridAksClustersClient.BeginRestartNode` has been removed
- Function `*HybridAksClustersClient.Update` has been removed
- Function `*MetricsConfigurationsClient.NewListByResourceGroupPager` has been removed
- Function `*StorageAppliancesClient.BeginValidateHardware` has been removed
- Struct `BgpPeer` has been removed
- Struct `CniBgpConfiguration` has been removed
- Struct `CommunityAdvertisement` has been removed
- Struct `DefaultCniNetwork` has been removed
- Struct `DefaultCniNetworkList` has been removed
- Struct `DefaultCniNetworkPatchParameters` has been removed
- Struct `DefaultCniNetworkProperties` has been removed
- Struct `HybridAksCluster` has been removed
- Struct `HybridAksClusterList` has been removed
- Struct `HybridAksClusterPatchParameters` has been removed
- Struct `HybridAksClusterProperties` has been removed
- Struct `HybridAksClusterRestartNodeParameters` has been removed
- Struct `Node` has been removed
- Struct `NodeConfiguration` has been removed
- Struct `StorageApplianceValidateHardwareParameters` has been removed

### Features Added

- New value `VirtualMachineDetailedStatusRunning`, `VirtualMachineDetailedStatusScheduling`, `VirtualMachineDetailedStatusStopped`, `VirtualMachineDetailedStatusTerminating`, `VirtualMachineDetailedStatusUnknown` added to enum type `VirtualMachineDetailedStatus`
- New value `VirtualMachinePowerStateUnknown` added to enum type `VirtualMachinePowerState`
- New enum type `AdvertiseToFabric` with values `AdvertiseToFabricFalse`, `AdvertiseToFabricTrue`
- New enum type `AgentPoolDetailedStatus` with values `AgentPoolDetailedStatusAvailable`, `AgentPoolDetailedStatusError`, `AgentPoolDetailedStatusProvisioning`
- New enum type `AgentPoolMode` with values `AgentPoolModeNotApplicable`, `AgentPoolModeSystem`, `AgentPoolModeUser`
- New enum type `AgentPoolProvisioningState` with values `AgentPoolProvisioningStateAccepted`, `AgentPoolProvisioningStateCanceled`, `AgentPoolProvisioningStateDeleting`, `AgentPoolProvisioningStateFailed`, `AgentPoolProvisioningStateInProgress`, `AgentPoolProvisioningStateSucceeded`, `AgentPoolProvisioningStateUpdating`
- New enum type `AvailabilityLifecycle` with values `AvailabilityLifecycleGenerallyAvailable`, `AvailabilityLifecyclePreview`
- New enum type `BfdEnabled` with values `BfdEnabledFalse`, `BfdEnabledTrue`
- New enum type `BgpMultiHop` with values `BgpMultiHopFalse`, `BgpMultiHopTrue`
- New enum type `FabricPeeringEnabled` with values `FabricPeeringEnabledFalse`, `FabricPeeringEnabledTrue`
- New enum type `FeatureDetailedStatus` with values `FeatureDetailedStatusFailed`, `FeatureDetailedStatusRunning`, `FeatureDetailedStatusUnknown`
- New enum type `HugepagesSize` with values `HugepagesSizeOneG`, `HugepagesSizeTwoM`
- New enum type `KubernetesClusterDetailedStatus` with values `KubernetesClusterDetailedStatusAvailable`, `KubernetesClusterDetailedStatusError`, `KubernetesClusterDetailedStatusProvisioning`
- New enum type `KubernetesClusterNodeDetailedStatus` with values `KubernetesClusterNodeDetailedStatusAvailable`, `KubernetesClusterNodeDetailedStatusError`, `KubernetesClusterNodeDetailedStatusProvisioning`, `KubernetesClusterNodeDetailedStatusRunning`, `KubernetesClusterNodeDetailedStatusScheduling`, `KubernetesClusterNodeDetailedStatusStopped`, `KubernetesClusterNodeDetailedStatusTerminating`, `KubernetesClusterNodeDetailedStatusUnknown`
- New enum type `KubernetesClusterProvisioningState` with values `KubernetesClusterProvisioningStateAccepted`, `KubernetesClusterProvisioningStateCanceled`, `KubernetesClusterProvisioningStateCreated`, `KubernetesClusterProvisioningStateDeleting`, `KubernetesClusterProvisioningStateFailed`, `KubernetesClusterProvisioningStateInProgress`, `KubernetesClusterProvisioningStateSucceeded`, `KubernetesClusterProvisioningStateUpdating`
- New enum type `KubernetesNodePowerState` with values `KubernetesNodePowerStateOff`, `KubernetesNodePowerStateOn`, `KubernetesNodePowerStateUnknown`
- New enum type `KubernetesNodeRole` with values `KubernetesNodeRoleControlPlane`, `KubernetesNodeRoleWorker`
- New enum type `KubernetesPluginType` with values `KubernetesPluginTypeDPDK`, `KubernetesPluginTypeIPVLAN`, `KubernetesPluginTypeMACVLAN`, `KubernetesPluginTypeOSDevice`, `KubernetesPluginTypeSRIOV`
- New enum type `L3NetworkConfigurationIpamEnabled` with values `L3NetworkConfigurationIpamEnabledFalse`, `L3NetworkConfigurationIpamEnabledTrue`
- New function `NewAgentPoolsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AgentPoolsClient, error)`
- New function `*AgentPoolsClient.BeginCreateOrUpdate(context.Context, string, string, string, AgentPool, *AgentPoolsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AgentPoolsClientCreateOrUpdateResponse], error)`
- New function `*AgentPoolsClient.BeginDelete(context.Context, string, string, string, *AgentPoolsClientBeginDeleteOptions) (*runtime.Poller[AgentPoolsClientDeleteResponse], error)`
- New function `*AgentPoolsClient.Get(context.Context, string, string, string, *AgentPoolsClientGetOptions) (AgentPoolsClientGetResponse, error)`
- New function `*AgentPoolsClient.NewListByKubernetesClusterPager(string, string, *AgentPoolsClientListByKubernetesClusterOptions) *runtime.Pager[AgentPoolsClientListByKubernetesClusterResponse]`
- New function `*AgentPoolsClient.BeginUpdate(context.Context, string, string, string, AgentPoolPatchParameters, *AgentPoolsClientBeginUpdateOptions) (*runtime.Poller[AgentPoolsClientUpdateResponse], error)`
- New function `*BareMetalMachineKeySetsClient.NewListByClusterPager(string, string, *BareMetalMachineKeySetsClientListByClusterOptions) *runtime.Pager[BareMetalMachineKeySetsClientListByClusterResponse]`
- New function `*BmcKeySetsClient.NewListByClusterPager(string, string, *BmcKeySetsClientListByClusterOptions) *runtime.Pager[BmcKeySetsClientListByClusterResponse]`
- New function `*ClientFactory.NewAgentPoolsClient() *AgentPoolsClient`
- New function `*ClientFactory.NewKubernetesClustersClient() *KubernetesClustersClient`
- New function `*ConsolesClient.NewListByVirtualMachinePager(string, string, *ConsolesClientListByVirtualMachineOptions) *runtime.Pager[ConsolesClientListByVirtualMachineResponse]`
- New function `NewKubernetesClustersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*KubernetesClustersClient, error)`
- New function `*KubernetesClustersClient.BeginCreateOrUpdate(context.Context, string, string, KubernetesCluster, *KubernetesClustersClientBeginCreateOrUpdateOptions) (*runtime.Poller[KubernetesClustersClientCreateOrUpdateResponse], error)`
- New function `*KubernetesClustersClient.BeginDelete(context.Context, string, string, *KubernetesClustersClientBeginDeleteOptions) (*runtime.Poller[KubernetesClustersClientDeleteResponse], error)`
- New function `*KubernetesClustersClient.Get(context.Context, string, string, *KubernetesClustersClientGetOptions) (KubernetesClustersClientGetResponse, error)`
- New function `*KubernetesClustersClient.NewListByResourceGroupPager(string, *KubernetesClustersClientListByResourceGroupOptions) *runtime.Pager[KubernetesClustersClientListByResourceGroupResponse]`
- New function `*KubernetesClustersClient.NewListBySubscriptionPager(*KubernetesClustersClientListBySubscriptionOptions) *runtime.Pager[KubernetesClustersClientListBySubscriptionResponse]`
- New function `*KubernetesClustersClient.BeginRestartNode(context.Context, string, string, KubernetesClusterRestartNodeParameters, *KubernetesClustersClientBeginRestartNodeOptions) (*runtime.Poller[KubernetesClustersClientRestartNodeResponse], error)`
- New function `*KubernetesClustersClient.BeginUpdate(context.Context, string, string, KubernetesClusterPatchParameters, *KubernetesClustersClientBeginUpdateOptions) (*runtime.Poller[KubernetesClustersClientUpdateResponse], error)`
- New function `*MetricsConfigurationsClient.NewListByClusterPager(string, string, *MetricsConfigurationsClientListByClusterOptions) *runtime.Pager[MetricsConfigurationsClientListByClusterResponse]`
- New struct `AADConfiguration`
- New struct `AdministratorConfiguration`
- New struct `AgentOptions`
- New struct `AgentPool`
- New struct `AgentPoolList`
- New struct `AgentPoolPatchParameters`
- New struct `AgentPoolPatchProperties`
- New struct `AgentPoolProperties`
- New struct `AgentPoolUpgradeSettings`
- New struct `AttachedNetworkConfiguration`
- New struct `AvailableUpgrade`
- New struct `BgpAdvertisement`
- New struct `BgpServiceLoadBalancerConfiguration`
- New struct `ControlPlaneNodeConfiguration`
- New struct `ControlPlaneNodePatchConfiguration`
- New struct `FeatureStatus`
- New struct `IPAddressPool`
- New struct `InitialAgentPoolConfiguration`
- New struct `KubernetesCluster`
- New struct `KubernetesClusterList`
- New struct `KubernetesClusterNode`
- New struct `KubernetesClusterPatchParameters`
- New struct `KubernetesClusterPatchProperties`
- New struct `KubernetesClusterProperties`
- New struct `KubernetesClusterRestartNodeParameters`
- New struct `KubernetesLabel`
- New struct `L2NetworkAttachmentConfiguration`
- New struct `L3NetworkAttachmentConfiguration`
- New struct `NetworkConfiguration`
- New struct `ServiceLoadBalancerBgpPeer`
- New struct `TrunkedNetworkAttachmentConfiguration`
- New field `AssociatedResourceIDs` in struct `BareMetalMachineProperties`
- New field `AssociatedResourceIDs` in struct `CloudServicesNetworkProperties`
- New field `AssociatedResourceIDs` in struct `L2NetworkProperties`
- New field `AssociatedResourceIDs` in struct `L3NetworkProperties`
- New field `AssociatedResourceIDs` in struct `TrunkedNetworkProperties`
- New field `AvailabilityZone` in struct `VirtualMachineProperties`


## 0.1.0 (2023-05-26)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/networkcloud/armnetworkcloud` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).