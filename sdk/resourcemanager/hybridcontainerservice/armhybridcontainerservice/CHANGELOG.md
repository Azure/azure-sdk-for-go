# Release History

## 1.0.0 (2024-01-26)
### Breaking Changes

- Function `*ProvisionedClusterInstancesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, ProvisionedClusters, *ProvisionedClusterInstancesClientBeginCreateOrUpdateOptions)` to `(context.Context, string, ProvisionedCluster, *ProvisionedClusterInstancesClientBeginCreateOrUpdateOptions)`
- Type of `ControlPlaneProfile.ControlPlaneEndpoint` has been changed from `*ControlPlaneEndpointProfileControlPlaneEndpoint` to `*ControlPlaneProfileControlPlaneEndpoint`
- Type of `VirtualNetworkExtendedLocation.Type` has been changed from `*string` to `*ExtendedLocationTypes`
- `NetworkPolicyFlannel` from enum `NetworkPolicy` has been removed
- `ProvisioningStateCreated`, `ProvisioningStateInProgress` from enum `ProvisioningState` has been removed
- `ResourceProvisioningStateCreated`, `ResourceProvisioningStateInProgress` from enum `ResourceProvisioningState` has been removed
- Function `*AgentPoolClient.BeginUpdate` has been removed
- Operation `*AgentPoolClient.ListByProvisionedCluster` has supported pagination, use `*AgentPoolClient.NewListByProvisionedClusterPager` instead.
- Struct `AgentPoolPatch` has been removed
- Struct `AgentPoolProvisioningStatusOperationStatus` has been removed
- Struct `AgentPoolProvisioningStatusOperationStatusError` has been removed
- Struct `ControlPlaneEndpointProfileControlPlaneEndpoint` has been removed
- Struct `KubernetesVersionCapabilities` has been removed
- Struct `ProvisionedClusterPropertiesStatusOperationStatus` has been removed
- Struct `ProvisionedClusterPropertiesStatusOperationStatusError` has been removed
- Struct `ProvisionedClusters` has been removed
- Struct `ProvisionedClustersListResult` has been removed
- Struct `VirtualNetworkPropertiesInfraVnetProfileVmware` has been removed
- Field `Location` of struct `AgentPool` has been removed
- Field `AvailabilityZones`, `NodeImageVersion` of struct `AgentPoolProperties` has been removed
- Field `OperationStatus` of struct `AgentPoolProvisioningStatusStatus` has been removed
- Field `AvailabilityZones`, `LinuxProfile`, `Name`, `NodeImageVersion`, `OSSKU`, `OSType` of struct `ControlPlaneProfile` has been removed
- Field `Capabilities` of struct `KubernetesVersionProperties` has been removed
- Field `AvailabilityZones`, `NodeImageVersion` of struct `NamedAgentPoolProfile` has been removed
- Field `ProvisionedClusters` of struct `ProvisionedClusterInstancesClientCreateOrUpdateResponse` has been removed
- Field `ProvisionedClusters` of struct `ProvisionedClusterInstancesClientGetResponse` has been removed
- Field `ProvisionedClustersListResult` of struct `ProvisionedClusterInstancesClientListResponse` has been removed
- Field `Name` of struct `ProvisionedClusterPoolUpgradeProfile` has been removed
- Field `OperationStatus` of struct `ProvisionedClusterPropertiesStatus` has been removed
- Field `AgentPoolProfiles` of struct `ProvisionedClusterUpgradeProfileProperties` has been removed
- Field `DhcpServers` of struct `VirtualNetworkProperties` has been removed
- Field `Vmware` of struct `VirtualNetworkPropertiesInfraVnetProfile` has been removed
- Field `Phase` of struct `VirtualNetworkPropertiesStatusOperationStatus` has been removed

### Features Added

- New value `ProvisioningStateCreating`, `ProvisioningStatePending` added to enum type `ProvisioningState`
- New value `ResourceProvisioningStatePending` added to enum type `ResourceProvisioningState`
- New enum type `Expander` with values `ExpanderLeastWaste`, `ExpanderMostPods`, `ExpanderPriority`, `ExpanderRandom`
- New struct `ClusterVMAccessProfile`
- New struct `ControlPlaneProfileControlPlaneEndpoint`
- New struct `ProvisionedCluster`
- New struct `ProvisionedClusterListResult`
- New struct `ProvisionedClusterPropertiesAutoScalerProfile`
- New struct `StorageProfile`
- New struct `StorageProfileNfsCSIDriver`
- New struct `StorageProfileSmbCSIDriver`
- New field `EnableAutoScaling`, `KubernetesVersion`, `MaxCount`, `MaxPods`, `MinCount`, `NodeLabels`, `NodeTaints` in struct `AgentPoolProperties`
- New field `CurrentState` in struct `AgentPoolProvisioningStatusStatus`
- New field `KubernetesVersion` in struct `AgentPoolUpdateProfile`
- New field `EnableAutoScaling`, `KubernetesVersion`, `MaxCount`, `MaxPods`, `MinCount`, `NodeLabels`, `NodeTaints` in struct `NamedAgentPoolProfile`
- New anonymous field `ProvisionedCluster` in struct `ProvisionedClusterInstancesClientCreateOrUpdateResponse`
- New anonymous field `ProvisionedCluster` in struct `ProvisionedClusterInstancesClientGetResponse`
- New anonymous field `ProvisionedClusterListResult` in struct `ProvisionedClusterInstancesClientListResponse`
- New field `AutoScalerProfile`, `ClusterVMAccessProfile`, `StorageProfile` in struct `ProvisionedClusterProperties`
- New field `CurrentState` in struct `ProvisionedClusterPropertiesStatus`


## 0.3.0 (2023-11-24)
### Breaking Changes

- Function `*AgentPoolClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, AgentPool, *AgentPoolClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, AgentPool, *AgentPoolClientBeginCreateOrUpdateOptions)`
- Function `*AgentPoolClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *AgentPoolClientGetOptions)` to `(context.Context, string, string, *AgentPoolClientGetOptions)`
- Function `*AgentPoolClient.ListByProvisionedCluster` parameter(s) have been changed from `(context.Context, string, string, *AgentPoolClientListByProvisionedClusterOptions)` to `(context.Context, string, *AgentPoolClientListByProvisionedClusterOptions)`
- Function `*HybridIdentityMetadataClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *HybridIdentityMetadataClientGetOptions)` to `(context.Context, string, *HybridIdentityMetadataClientGetOptions)`
- Function `*HybridIdentityMetadataClient.NewListByClusterPager` parameter(s) have been changed from `(string, string, *HybridIdentityMetadataClientListByClusterOptions)` to `(string, *HybridIdentityMetadataClientListByClusterOptions)`
- Function `*HybridIdentityMetadataClient.Put` parameter(s) have been changed from `(context.Context, string, string, string, HybridIdentityMetadata, *HybridIdentityMetadataClientPutOptions)` to `(context.Context, string, HybridIdentityMetadata, *HybridIdentityMetadataClientPutOptions)`
- Function `*VirtualNetworksClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, VirtualNetworks, *VirtualNetworksClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, VirtualNetwork, *VirtualNetworksClientBeginCreateOrUpdateOptions)`
- Type of `AgentPool.ExtendedLocation` has been changed from `*AgentPoolExtendedLocation` to `*ExtendedLocation`
- Type of `AgentPoolProperties.ProvisioningState` has been changed from `*AgentPoolProvisioningState` to `*ResourceProvisioningState`
- Type of `AgentPoolProvisioningStatusStatus.ReadyReplicas` has been changed from `*int32` to `[]*AgentPoolUpdateProfile`
- Type of `ControlPlaneEndpointProfileControlPlaneEndpoint.Port` has been changed from `*string` to `*int32`
- Type of `HybridIdentityMetadataProperties.ProvisioningState` has been changed from `*string` to `*ResourceProvisioningState`
- Type of `NetworkProfile.LoadBalancerProfile` has been changed from `*LoadBalancerProfile` to `*NetworkProfileLoadBalancerProfile`
- Type of `ProvisionedClusterUpgradeProfileProperties.ProvisioningState` has been changed from `*string` to `*ResourceProvisioningState`
- Type of `ProvisionedClusters.ExtendedLocation` has been changed from `*ProvisionedClustersExtendedLocation` to `*ExtendedLocation`
- Type of `ProvisionedClusters.Properties` has been changed from `*ProvisionedClustersAllProperties` to `*ProvisionedClusterProperties`
- Type of `VirtualNetworksListResult.Value` has been changed from `[]*VirtualNetworks` to `[]*VirtualNetwork`
- Enum `AgentPoolProvisioningState` has been removed
- Enum `AutoUpgradeOptions` has been removed
- Enum `DeploymentState` has been removed
- Enum `LicenseType` has been removed
- Enum `LoadBalancerSKU` has been removed
- Enum `Mode` has been removed
- Enum `ResourceIdentityType` has been removed
- Function `*Client.ListOrchestrators` has been removed
- Function `*Client.ListVMSKUs` has been removed
- Function `*ClientFactory.NewProvisionedClustersClient` has been removed
- Function `*ClientFactory.NewStorageSpacesClient` has been removed
- Function `NewProvisionedClustersClient` has been removed
- Function `*ProvisionedClustersClient.BeginCreateOrUpdate` has been removed
- Function `*ProvisionedClustersClient.Delete` has been removed
- Function `*ProvisionedClustersClient.Get` has been removed
- Function `*ProvisionedClustersClient.GetUpgradeProfile` has been removed
- Function `*ProvisionedClustersClient.NewListByResourceGroupPager` has been removed
- Function `*ProvisionedClustersClient.NewListBySubscriptionPager` has been removed
- Function `*ProvisionedClustersClient.BeginUpdate` has been removed
- Function `*ProvisionedClustersClient.BeginUpgradeNodeImageVersionForEntireCluster` has been removed
- Function `NewStorageSpacesClient` has been removed
- Function `*StorageSpacesClient.BeginCreateOrUpdate` has been removed
- Function `*StorageSpacesClient.Delete` has been removed
- Function `*StorageSpacesClient.NewListByResourceGroupPager` has been removed
- Function `*StorageSpacesClient.NewListBySubscriptionPager` has been removed
- Function `*StorageSpacesClient.Retrieve` has been removed
- Function `*StorageSpacesClient.BeginUpdate` has been removed
- Operation `*AgentPoolClient.Delete` has been changed to LRO, use `*AgentPoolClient.BeginDelete` instead.
- Operation `*AgentPoolClient.Update` has been changed to LRO, use `*AgentPoolClient.BeginUpdate` instead.
- Operation `*HybridIdentityMetadataClient.Delete` has been changed to LRO, use `*HybridIdentityMetadataClient.BeginDelete` instead.
- Operation `*VirtualNetworksClient.Delete` has been changed to LRO, use `*VirtualNetworksClient.BeginDelete` instead.
- Struct `AADProfile` has been removed
- Struct `AADProfileResponse` has been removed
- Struct `AddonProfiles` has been removed
- Struct `AddonStatus` has been removed
- Struct `AgentPoolExtendedLocation` has been removed
- Struct `AgentPoolProvisioningStatusError` has been removed
- Struct `AgentPoolProvisioningStatusStatusProvisioningStatus` has been removed
- Struct `ArcAgentProfile` has been removed
- Struct `ArcAgentStatus` has been removed
- Struct `CloudProviderProfileInfraStorageProfile` has been removed
- Struct `HTTPProxyConfig` has been removed
- Struct `HTTPProxyConfigResponse` has been removed
- Struct `LoadBalancerProfile` has been removed
- Struct `OrchestratorProfile` has been removed
- Struct `OrchestratorVersionProfile` has been removed
- Struct `OrchestratorVersionProfileListResult` has been removed
- Struct `ProvisionedClusterIdentity` has been removed
- Struct `ProvisionedClustersAllProperties` has been removed
- Struct `ProvisionedClustersCommonPropertiesFeatures` has been removed
- Struct `ProvisionedClustersCommonPropertiesStatus` has been removed
- Struct `ProvisionedClustersCommonPropertiesStatusFeaturesStatus` has been removed
- Struct `ProvisionedClustersCommonPropertiesStatusProvisioningStatus` has been removed
- Struct `ProvisionedClustersCommonPropertiesStatusProvisioningStatusError` has been removed
- Struct `ProvisionedClustersExtendedLocation` has been removed
- Struct `ProvisionedClustersPatch` has been removed
- Struct `ProvisionedClustersResponse` has been removed
- Struct `ProvisionedClustersResponseExtendedLocation` has been removed
- Struct `ProvisionedClustersResponseListResult` has been removed
- Struct `ProvisionedClustersResponseProperties` has been removed
- Struct `ResourceProviderOperation` has been removed
- Struct `ResourceProviderOperationDisplay` has been removed
- Struct `ResourceProviderOperationList` has been removed
- Struct `StorageSpaces` has been removed
- Struct `StorageSpacesExtendedLocation` has been removed
- Struct `StorageSpacesListResult` has been removed
- Struct `StorageSpacesPatch` has been removed
- Struct `StorageSpacesProperties` has been removed
- Struct `StorageSpacesPropertiesHciStorageProfile` has been removed
- Struct `StorageSpacesPropertiesStatus` has been removed
- Struct `StorageSpacesPropertiesStatusProvisioningStatus` has been removed
- Struct `StorageSpacesPropertiesStatusProvisioningStatusError` has been removed
- Struct `StorageSpacesPropertiesVmwareStorageProfile` has been removed
- Struct `VMSKUListResult` has been removed
- Struct `VirtualNetworks` has been removed
- Struct `VirtualNetworksExtendedLocation` has been removed
- Struct `VirtualNetworksProperties` has been removed
- Struct `VirtualNetworksPropertiesInfraVnetProfile` has been removed
- Struct `VirtualNetworksPropertiesInfraVnetProfileHci` has been removed
- Struct `VirtualNetworksPropertiesInfraVnetProfileNetworkCloud` has been removed
- Struct `VirtualNetworksPropertiesInfraVnetProfileVmware` has been removed
- Struct `VirtualNetworksPropertiesStatus` has been removed
- Struct `VirtualNetworksPropertiesStatusProvisioningStatus` has been removed
- Struct `VirtualNetworksPropertiesStatusProvisioningStatusError` has been removed
- Struct `VirtualNetworksPropertiesVipPoolItem` has been removed
- Struct `VirtualNetworksPropertiesVmipPoolItem` has been removed
- Struct `WindowsProfile` has been removed
- Struct `WindowsProfileResponse` has been removed
- Field `CloudProviderProfile`, `MaxCount`, `MaxPods`, `MinCount`, `Mode`, `NodeLabels`, `NodeTaints` of struct `AgentPoolProperties` has been removed
- Field `ProvisioningStatus`, `Replicas` of struct `AgentPoolProvisioningStatusStatus` has been removed
- Field `InfraStorageProfile` of struct `CloudProviderProfile` has been removed
- Field `CloudProviderProfile`, `MaxCount`, `MaxPods`, `MinCount`, `Mode`, `NodeLabels`, `NodeTaints` of struct `ControlPlaneProfile` has been removed
- Field `Identity` of struct `HybridIdentityMetadataProperties` has been removed
- Field `AdminUsername` of struct `LinuxProfileProperties` has been removed
- Field `CloudProviderProfile`, `MaxCount`, `MaxPods`, `MinCount`, `Mode`, `NodeLabels`, `NodeTaints` of struct `NamedAgentPoolProfile` has been removed
- Field `DNSServiceIP`, `LoadBalancerSKU`, `PodCidrs`, `ServiceCidr`, `ServiceCidrs` of struct `NetworkProfile` has been removed
- Field `ResourceProviderOperationList` of struct `OperationsClientListResponse` has been removed
- Field `Identity`, `Location`, `Tags` of struct `ProvisionedClusters` has been removed
- Field `VirtualNetworks` of struct `VirtualNetworksClientCreateOrUpdateResponse` has been removed
- Field `VirtualNetworks` of struct `VirtualNetworksClientRetrieveResponse` has been removed
- Field `VirtualNetworks` of struct `VirtualNetworksClientUpdateResponse` has been removed

### Features Added

- Support for test fakes and OpenTelemetry trace spans.
- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `AddonPhase` with values `AddonPhaseDeleting`, `AddonPhaseFailed`, `AddonPhasePending`, `AddonPhaseProvisioned`, `AddonPhaseProvisioning`, `AddonPhaseProvisioningHelmChartInstalled`, `AddonPhaseProvisioningMSICertificateDownloaded`, `AddonPhaseUpgrading`
- New enum type `AzureHybridBenefit` with values `AzureHybridBenefitFalse`, `AzureHybridBenefitNotApplicable`, `AzureHybridBenefitTrue`
- New enum type `ExtendedLocationTypes` with values `ExtendedLocationTypesCustomLocation`
- New enum type `OSSKU` with values `OSSKUCBLMariner`, `OSSKUWindows2019`, `OSSKUWindows2022`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `ResourceProvisioningState` with values `ResourceProvisioningStateAccepted`, `ResourceProvisioningStateCanceled`, `ResourceProvisioningStateCreated`, `ResourceProvisioningStateCreating`, `ResourceProvisioningStateDeleting`, `ResourceProvisioningStateFailed`, `ResourceProvisioningStateInProgress`, `ResourceProvisioningStateSucceeded`, `ResourceProvisioningStateUpdating`, `ResourceProvisioningStateUpgrading`
- New function `*Client.BeginDeleteKubernetesVersions(context.Context, string, *ClientBeginDeleteKubernetesVersionsOptions) (*runtime.Poller[ClientDeleteKubernetesVersionsResponse], error)`
- New function `*Client.BeginDeleteVMSKUs(context.Context, string, *ClientBeginDeleteVMSKUsOptions) (*runtime.Poller[ClientDeleteVMSKUsResponse], error)`
- New function `*Client.GetKubernetesVersions(context.Context, string, *ClientGetKubernetesVersionsOptions) (ClientGetKubernetesVersionsResponse, error)`
- New function `*Client.GetVMSKUs(context.Context, string, *ClientGetVMSKUsOptions) (ClientGetVMSKUsResponse, error)`
- New function `*Client.BeginPutKubernetesVersions(context.Context, string, KubernetesVersionProfile, *ClientBeginPutKubernetesVersionsOptions) (*runtime.Poller[ClientPutKubernetesVersionsResponse], error)`
- New function `*Client.BeginPutVMSKUs(context.Context, string, VMSKUProfile, *ClientBeginPutVMSKUsOptions) (*runtime.Poller[ClientPutVMSKUsResponse], error)`
- New function `*ClientFactory.NewKubernetesVersionsClient() *KubernetesVersionsClient`
- New function `*ClientFactory.NewProvisionedClusterInstancesClient() *ProvisionedClusterInstancesClient`
- New function `*ClientFactory.NewVMSKUsClient() *VMSKUsClient`
- New function `NewKubernetesVersionsClient(azcore.TokenCredential, *arm.ClientOptions) (*KubernetesVersionsClient, error)`
- New function `*KubernetesVersionsClient.NewListPager(string, *KubernetesVersionsClientListOptions) *runtime.Pager[KubernetesVersionsClientListResponse]`
- New function `NewProvisionedClusterInstancesClient(azcore.TokenCredential, *arm.ClientOptions) (*ProvisionedClusterInstancesClient, error)`
- New function `*ProvisionedClusterInstancesClient.BeginCreateOrUpdate(context.Context, string, ProvisionedClusters, *ProvisionedClusterInstancesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ProvisionedClusterInstancesClientCreateOrUpdateResponse], error)`
- New function `*ProvisionedClusterInstancesClient.BeginDelete(context.Context, string, *ProvisionedClusterInstancesClientBeginDeleteOptions) (*runtime.Poller[ProvisionedClusterInstancesClientDeleteResponse], error)`
- New function `*ProvisionedClusterInstancesClient.Get(context.Context, string, *ProvisionedClusterInstancesClientGetOptions) (ProvisionedClusterInstancesClientGetResponse, error)`
- New function `*ProvisionedClusterInstancesClient.GetUpgradeProfile(context.Context, string, *ProvisionedClusterInstancesClientGetUpgradeProfileOptions) (ProvisionedClusterInstancesClientGetUpgradeProfileResponse, error)`
- New function `*ProvisionedClusterInstancesClient.BeginListAdminKubeconfig(context.Context, string, *ProvisionedClusterInstancesClientBeginListAdminKubeconfigOptions) (*runtime.Poller[ProvisionedClusterInstancesClientListAdminKubeconfigResponse], error)`
- New function `*ProvisionedClusterInstancesClient.NewListPager(string, *ProvisionedClusterInstancesClientListOptions) *runtime.Pager[ProvisionedClusterInstancesClientListResponse]`
- New function `*ProvisionedClusterInstancesClient.BeginListUserKubeconfig(context.Context, string, *ProvisionedClusterInstancesClientBeginListUserKubeconfigOptions) (*runtime.Poller[ProvisionedClusterInstancesClientListUserKubeconfigResponse], error)`
- New function `NewVMSKUsClient(azcore.TokenCredential, *arm.ClientOptions) (*VMSKUsClient, error)`
- New function `*VMSKUsClient.NewListPager(string, *VMSKUsClientListOptions) *runtime.Pager[VMSKUsClientListResponse]`
- New struct `AddonStatusProfile`
- New struct `AgentPoolPatch`
- New struct `AgentPoolProvisioningStatusOperationStatus`
- New struct `AgentPoolProvisioningStatusOperationStatusError`
- New struct `AgentPoolUpdateProfile`
- New struct `CredentialResult`
- New struct `ExtendedLocation`
- New struct `KubernetesPatchVersions`
- New struct `KubernetesVersionCapabilities`
- New struct `KubernetesVersionProfile`
- New struct `KubernetesVersionProfileList`
- New struct `KubernetesVersionProfileProperties`
- New struct `KubernetesVersionProperties`
- New struct `KubernetesVersionReadiness`
- New struct `ListCredentialResponse`
- New struct `ListCredentialResponseError`
- New struct `ListCredentialResponseProperties`
- New struct `NetworkProfileLoadBalancerProfile`
- New struct `Operation`
- New struct `OperationDisplay`
- New struct `OperationListResult`
- New struct `ProvisionedClusterLicenseProfile`
- New struct `ProvisionedClusterProperties`
- New struct `ProvisionedClusterPropertiesStatus`
- New struct `ProvisionedClusterPropertiesStatusOperationStatus`
- New struct `ProvisionedClusterPropertiesStatusOperationStatusError`
- New struct `ProvisionedClustersListResult`
- New struct `VMSKUCapabilities`
- New struct `VMSKUProfile`
- New struct `VMSKUProfileList`
- New struct `VMSKUProfileProperties`
- New struct `VMSKUProperties`
- New struct `VirtualNetwork`
- New struct `VirtualNetworkExtendedLocation`
- New struct `VirtualNetworkProperties`
- New struct `VirtualNetworkPropertiesInfraVnetProfile`
- New struct `VirtualNetworkPropertiesInfraVnetProfileHci`
- New struct `VirtualNetworkPropertiesInfraVnetProfileVmware`
- New struct `VirtualNetworkPropertiesStatus`
- New struct `VirtualNetworkPropertiesStatusOperationStatus`
- New struct `VirtualNetworkPropertiesStatusOperationStatusError`
- New struct `VirtualNetworkPropertiesVipPoolItem`
- New struct `VirtualNetworkPropertiesVmipPoolItem`
- New field `OSSKU` in struct `AgentPoolProperties`
- New field `OperationStatus` in struct `AgentPoolProvisioningStatusStatus`
- New field `OSSKU` in struct `ControlPlaneProfile`
- New field `OSSKU` in struct `NamedAgentPoolProfile`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New field `SystemData` in struct `ProvisionedClusterUpgradeProfile`
- New anonymous field `VirtualNetwork` in struct `VirtualNetworksClientCreateOrUpdateResponse`
- New anonymous field `VirtualNetwork` in struct `VirtualNetworksClientRetrieveResponse`
- New anonymous field `VirtualNetwork` in struct `VirtualNetworksClientUpdateResponse`


## 0.2.0 (2023-03-24)
### Breaking Changes

- Struct `VirtualNetworksPropertiesInfraVnetProfileKubevirt` has been removed
- Field `Kubevirt` of struct `VirtualNetworksPropertiesInfraVnetProfile` has been removed

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New function `*ProvisionedClustersClient.GetUpgradeProfile(context.Context, string, string, *ProvisionedClustersClientGetUpgradeProfileOptions) (ProvisionedClustersClientGetUpgradeProfileResponse, error)`
- New function `*ProvisionedClustersClient.BeginUpgradeNodeImageVersionForEntireCluster(context.Context, string, string, *ProvisionedClustersClientBeginUpgradeNodeImageVersionForEntireClusterOptions) (*runtime.Poller[ProvisionedClustersClientUpgradeNodeImageVersionForEntireClusterResponse], error)`
- New struct `ProvisionedClusterPoolUpgradeProfile`
- New struct `ProvisionedClusterPoolUpgradeProfileProperties`
- New struct `ProvisionedClusterUpgradeProfile`
- New struct `ProvisionedClusterUpgradeProfileProperties`
- New struct `VirtualNetworksPropertiesInfraVnetProfileNetworkCloud`
- New field `NetworkCloud` in struct `VirtualNetworksPropertiesInfraVnetProfile`


## 0.1.1 (2022-10-12)
### Other Changes
- Loosen Go version requirement.

## 0.1.0 (2022-09-13)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridcontainerservice/armhybridcontainerservice` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
