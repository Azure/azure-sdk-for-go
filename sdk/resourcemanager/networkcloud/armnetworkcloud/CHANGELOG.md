# Release History

## 1.2.0-beta.1 (2024-11-15)
### Features Added

- New value `ClusterConnectionStatusDisconnected` added to enum type `ClusterConnectionStatus`
- New value `ClusterDetailedStatusUpdatePaused` added to enum type `ClusterDetailedStatus`
- New value `RackSKUProvisioningStateCanceled`, `RackSKUProvisioningStateFailed` added to enum type `RackSKUProvisioningState`
- New enum type `ClusterContinueUpdateVersionMachineGroupTargetingMode` with values `ClusterContinueUpdateVersionMachineGroupTargetingModeAlphaByRack`
- New enum type `ClusterScanRuntimeParametersScanActivity` with values `ClusterScanRuntimeParametersScanActivityScan`, `ClusterScanRuntimeParametersScanActivitySkip`
- New enum type `ClusterSecretArchiveEnabled` with values `ClusterSecretArchiveEnabledFalse`, `ClusterSecretArchiveEnabledTrue`
- New enum type `ClusterUpdateStrategyType` with values `ClusterUpdateStrategyTypePauseAfterRack`, `ClusterUpdateStrategyTypeRack`
- New enum type `KubernetesClusterFeatureAvailabilityLifecycle` with values `KubernetesClusterFeatureAvailabilityLifecycleGenerallyAvailable`, `KubernetesClusterFeatureAvailabilityLifecyclePreview`
- New enum type `KubernetesClusterFeatureDetailedStatus` with values `KubernetesClusterFeatureDetailedStatusError`, `KubernetesClusterFeatureDetailedStatusInstalled`, `KubernetesClusterFeatureDetailedStatusProvisioning`
- New enum type `KubernetesClusterFeatureProvisioningState` with values `KubernetesClusterFeatureProvisioningStateAccepted`, `KubernetesClusterFeatureProvisioningStateCanceled`, `KubernetesClusterFeatureProvisioningStateDeleting`, `KubernetesClusterFeatureProvisioningStateFailed`, `KubernetesClusterFeatureProvisioningStateSucceeded`, `KubernetesClusterFeatureProvisioningStateUpdating`
- New enum type `KubernetesClusterFeatureRequired` with values `KubernetesClusterFeatureRequiredFalse`, `KubernetesClusterFeatureRequiredTrue`
- New enum type `ManagedServiceIdentitySelectorType` with values `ManagedServiceIdentitySelectorTypeSystemAssignedIdentity`, `ManagedServiceIdentitySelectorTypeUserAssignedIdentity`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `RuntimeProtectionEnforcementLevel` with values `RuntimeProtectionEnforcementLevelAudit`, `RuntimeProtectionEnforcementLevelDisabled`, `RuntimeProtectionEnforcementLevelOnDemand`, `RuntimeProtectionEnforcementLevelPassive`, `RuntimeProtectionEnforcementLevelRealTime`
- New function `*ClientFactory.NewKubernetesClusterFeaturesClient() *KubernetesClusterFeaturesClient`
- New function `*ClustersClient.BeginContinueUpdateVersion(context.Context, string, string, ClusterContinueUpdateVersionParameters, *ClustersClientBeginContinueUpdateVersionOptions) (*runtime.Poller[ClustersClientContinueUpdateVersionResponse], error)`
- New function `*ClustersClient.BeginScanRuntime(context.Context, string, string, *ClustersClientBeginScanRuntimeOptions) (*runtime.Poller[ClustersClientScanRuntimeResponse], error)`
- New function `NewKubernetesClusterFeaturesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*KubernetesClusterFeaturesClient, error)`
- New function `*KubernetesClusterFeaturesClient.BeginCreateOrUpdate(context.Context, string, string, string, KubernetesClusterFeature, *KubernetesClusterFeaturesClientBeginCreateOrUpdateOptions) (*runtime.Poller[KubernetesClusterFeaturesClientCreateOrUpdateResponse], error)`
- New function `*KubernetesClusterFeaturesClient.BeginDelete(context.Context, string, string, string, *KubernetesClusterFeaturesClientBeginDeleteOptions) (*runtime.Poller[KubernetesClusterFeaturesClientDeleteResponse], error)`
- New function `*KubernetesClusterFeaturesClient.Get(context.Context, string, string, string, *KubernetesClusterFeaturesClientGetOptions) (KubernetesClusterFeaturesClientGetResponse, error)`
- New function `*KubernetesClusterFeaturesClient.NewListByKubernetesClusterPager(string, string, *KubernetesClusterFeaturesClientListByKubernetesClusterOptions) *runtime.Pager[KubernetesClusterFeaturesClientListByKubernetesClusterResponse]`
- New function `*KubernetesClusterFeaturesClient.BeginUpdate(context.Context, string, string, string, KubernetesClusterFeaturePatchParameters, *KubernetesClusterFeaturesClientBeginUpdateOptions) (*runtime.Poller[KubernetesClusterFeaturesClientUpdateResponse], error)`
- New struct `AdministratorConfigurationPatch`
- New struct `ClusterContinueUpdateVersionParameters`
- New struct `ClusterScanRuntimeParameters`
- New struct `ClusterSecretArchive`
- New struct `ClusterUpdateStrategy`
- New struct `CommandOutputSettings`
- New struct `IdentitySelector`
- New struct `KubernetesClusterFeature`
- New struct `KubernetesClusterFeatureList`
- New struct `KubernetesClusterFeaturePatchParameters`
- New struct `KubernetesClusterFeaturePatchProperties`
- New struct `KubernetesClusterFeatureProperties`
- New struct `L2ServiceLoadBalancerConfiguration`
- New struct `ManagedServiceIdentity`
- New struct `NodePoolAdministratorConfigurationPatch`
- New struct `OperationStatusResultProperties`
- New struct `RuntimeProtectionConfiguration`
- New struct `RuntimeProtectionStatus`
- New struct `SecretArchiveReference`
- New struct `SecretRotationStatus`
- New struct `StringKeyValuePair`
- New struct `UserAssignedIdentity`
- New field `AdministratorConfiguration` in struct `AgentPoolPatchProperties`
- New field `DrainTimeout`, `MaxUnavailable` in struct `AgentPoolUpgradeSettings`
- New anonymous field `OperationStatusResult` in struct `AgentPoolsClientDeleteResponse`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachineKeySetsClientDeleteResponse`
- New field `MachineClusterVersion`, `MachineRoles`, `RuntimeProtectionStatus`, `SecretRotationStatus` in struct `BareMetalMachineProperties`
- New anonymous field `OperationStatusResult` in struct `BareMetalMachinesClientDeleteResponse`
- New anonymous field `OperationStatusResult` in struct `BmcKeySetsClientDeleteResponse`
- New anonymous field `OperationStatusResult` in struct `CloudServicesNetworksClientDeleteResponse`
- New field `Identity` in struct `Cluster`
- New field `Identity` in struct `ClusterManager`
- New field `Identity` in struct `ClusterManagerPatchParameters`
- New anonymous field `OperationStatusResult` in struct `ClusterManagersClientDeleteResponse`
- New field `Identity` in struct `ClusterPatchParameters`
- New field `CommandOutputSettings`, `RuntimeProtectionConfiguration`, `SecretArchive`, `UpdateStrategy` in struct `ClusterPatchProperties`
- New field `CommandOutputSettings`, `RuntimeProtectionConfiguration`, `SecretArchive`, `UpdateStrategy` in struct `ClusterProperties`
- New anonymous field `OperationStatusResult` in struct `ClustersClientDeleteResponse`
- New anonymous field `OperationStatusResult` in struct `ConsolesClientDeleteResponse`
- New field `AdministratorConfiguration` in struct `ControlPlaneNodePatchConfiguration`
- New field `UserPrincipalName` in struct `KeySetUser`
- New field `AdministratorConfiguration` in struct `KubernetesClusterPatchProperties`
- New anonymous field `OperationStatusResult` in struct `KubernetesClustersClientDeleteResponse`
- New anonymous field `OperationStatusResult` in struct `L2NetworksClientDeleteResponse`
- New anonymous field `OperationStatusResult` in struct `L3NetworksClientDeleteResponse`
- New anonymous field `OperationStatusResult` in struct `MetricsConfigurationsClientDeleteResponse`
- New field `L2ServiceLoadBalancerConfiguration` in struct `NetworkConfiguration`
- New field `Properties` in struct `OperationStatusResult`
- New anonymous field `OperationStatusResult` in struct `RacksClientDeleteResponse`
- New field `Manufacturer`, `Model`, `SecretRotationStatus`, `Version` in struct `StorageApplianceProperties`
- New anonymous field `OperationStatusResult` in struct `StorageAppliancesClientDeleteResponse`
- New anonymous field `OperationStatusResult` in struct `TrunkedNetworksClientDeleteResponse`
- New anonymous field `OperationStatusResult` in struct `VirtualMachinesClientDeleteResponse`
- New anonymous field `OperationStatusResult` in struct `VolumesClientDeleteResponse`


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