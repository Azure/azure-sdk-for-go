# Release History

## 4.4.0 (2023-10-27)
### Features Added

- New enum type `IstioIngressGatewayMode` with values `IstioIngressGatewayModeExternal`, `IstioIngressGatewayModeInternal`
- New enum type `ServiceMeshMode` with values `ServiceMeshModeDisabled`, `ServiceMeshModeIstio`
- New function `*ManagedClustersClient.GetMeshRevisionProfile(context.Context, string, string, *ManagedClustersClientGetMeshRevisionProfileOptions) (ManagedClustersClientGetMeshRevisionProfileResponse, error)`
- New function `*ManagedClustersClient.GetMeshUpgradeProfile(context.Context, string, string, string, *ManagedClustersClientGetMeshUpgradeProfileOptions) (ManagedClustersClientGetMeshUpgradeProfileResponse, error)`
- New function `*ManagedClustersClient.NewListMeshRevisionProfilesPager(string, *ManagedClustersClientListMeshRevisionProfilesOptions) *runtime.Pager[ManagedClustersClientListMeshRevisionProfilesResponse]`
- New function `*ManagedClustersClient.NewListMeshUpgradeProfilesPager(string, string, *ManagedClustersClientListMeshUpgradeProfilesOptions) *runtime.Pager[ManagedClustersClientListMeshUpgradeProfilesResponse]`
- New struct `CompatibleVersions`
- New struct `IstioCertificateAuthority`
- New struct `IstioComponents`
- New struct `IstioEgressGateway`
- New struct `IstioIngressGateway`
- New struct `IstioPluginCertificateAuthority`
- New struct `IstioServiceMesh`
- New struct `MeshRevision`
- New struct `MeshRevisionProfile`
- New struct `MeshRevisionProfileList`
- New struct `MeshRevisionProfileProperties`
- New struct `MeshUpgradeProfile`
- New struct `MeshUpgradeProfileList`
- New struct `MeshUpgradeProfileProperties`
- New struct `ServiceMeshProfile`
- New field `ResourceUID`, `ServiceMeshProfile` in struct `ManagedClusterProperties`


## 4.4.0-beta.2 (2023-10-09)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.

## 4.4.0-beta.1 (2023-09-22)
### Features Added

- New value `NodeOSUpgradeChannelSecurityPatch` added to enum type `NodeOSUpgradeChannel`
- New value `OSSKUMariner` added to enum type `OSSKU`
- New value `PublicNetworkAccessSecuredByPerimeter` added to enum type `PublicNetworkAccess`
- New value `SnapshotTypeManagedCluster` added to enum type `SnapshotType`
- New value `WorkloadRuntimeKataMshvVMIsolation` added to enum type `WorkloadRuntime`
- New enum type `AgentPoolSSHAccess` with values `AgentPoolSSHAccessDisabled`, `AgentPoolSSHAccessLocalUser`
- New enum type `BackendPoolType` with values `BackendPoolTypeNodeIP`, `BackendPoolTypeNodeIPConfiguration`
- New enum type `GuardrailsSupport` with values `GuardrailsSupportPreview`, `GuardrailsSupportStable`
- New enum type `IpvsScheduler` with values `IpvsSchedulerLeastConnection`, `IpvsSchedulerRoundRobin`
- New enum type `IstioIngressGatewayMode` with values `IstioIngressGatewayModeExternal`, `IstioIngressGatewayModeInternal`
- New enum type `Level` with values `LevelEnforcement`, `LevelOff`, `LevelWarning`
- New enum type `Mode` with values `ModeIPTABLES`, `ModeIPVS`
- New enum type `Protocol` with values `ProtocolTCP`, `ProtocolUDP`
- New enum type `RestrictionLevel` with values `RestrictionLevelReadOnly`, `RestrictionLevelUnrestricted`
- New enum type `ServiceMeshMode` with values `ServiceMeshModeDisabled`, `ServiceMeshModeIstio`
- New enum type `TrustedAccessRoleBindingProvisioningState` with values `TrustedAccessRoleBindingProvisioningStateCanceled`, `TrustedAccessRoleBindingProvisioningStateDeleting`, `TrustedAccessRoleBindingProvisioningStateFailed`, `TrustedAccessRoleBindingProvisioningStateSucceeded`, `TrustedAccessRoleBindingProvisioningStateUpdating`
- New function `*ClientFactory.NewMachinesClient() *MachinesClient`
- New function `*ClientFactory.NewManagedClusterSnapshotsClient() *ManagedClusterSnapshotsClient`
- New function `*ClientFactory.NewTrustedAccessRoleBindingsClient() *TrustedAccessRoleBindingsClient`
- New function `*ClientFactory.NewTrustedAccessRolesClient() *TrustedAccessRolesClient`
- New function `NewMachinesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MachinesClient, error)`
- New function `*MachinesClient.Get(context.Context, string, string, string, string, *MachinesClientGetOptions) (MachinesClientGetResponse, error)`
- New function `*MachinesClient.NewListPager(string, string, string, *MachinesClientListOptions) *runtime.Pager[MachinesClientListResponse]`
- New function `NewManagedClusterSnapshotsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedClusterSnapshotsClient, error)`
- New function `*ManagedClusterSnapshotsClient.CreateOrUpdate(context.Context, string, string, ManagedClusterSnapshot, *ManagedClusterSnapshotsClientCreateOrUpdateOptions) (ManagedClusterSnapshotsClientCreateOrUpdateResponse, error)`
- New function `*ManagedClusterSnapshotsClient.Delete(context.Context, string, string, *ManagedClusterSnapshotsClientDeleteOptions) (ManagedClusterSnapshotsClientDeleteResponse, error)`
- New function `*ManagedClusterSnapshotsClient.Get(context.Context, string, string, *ManagedClusterSnapshotsClientGetOptions) (ManagedClusterSnapshotsClientGetResponse, error)`
- New function `*ManagedClusterSnapshotsClient.NewListByResourceGroupPager(string, *ManagedClusterSnapshotsClientListByResourceGroupOptions) *runtime.Pager[ManagedClusterSnapshotsClientListByResourceGroupResponse]`
- New function `*ManagedClusterSnapshotsClient.NewListPager(*ManagedClusterSnapshotsClientListOptions) *runtime.Pager[ManagedClusterSnapshotsClientListResponse]`
- New function `*ManagedClusterSnapshotsClient.UpdateTags(context.Context, string, string, TagsObject, *ManagedClusterSnapshotsClientUpdateTagsOptions) (ManagedClusterSnapshotsClientUpdateTagsResponse, error)`
- New function `*ManagedClustersClient.GetGuardrailsVersions(context.Context, string, string, *ManagedClustersClientGetGuardrailsVersionsOptions) (ManagedClustersClientGetGuardrailsVersionsResponse, error)`
- New function `*ManagedClustersClient.GetMeshRevisionProfile(context.Context, string, string, *ManagedClustersClientGetMeshRevisionProfileOptions) (ManagedClustersClientGetMeshRevisionProfileResponse, error)`
- New function `*ManagedClustersClient.GetMeshUpgradeProfile(context.Context, string, string, string, *ManagedClustersClientGetMeshUpgradeProfileOptions) (ManagedClustersClientGetMeshUpgradeProfileResponse, error)`
- New function `*ManagedClustersClient.NewListGuardrailsVersionsPager(string, *ManagedClustersClientListGuardrailsVersionsOptions) *runtime.Pager[ManagedClustersClientListGuardrailsVersionsResponse]`
- New function `*ManagedClustersClient.NewListMeshRevisionProfilesPager(string, *ManagedClustersClientListMeshRevisionProfilesOptions) *runtime.Pager[ManagedClustersClientListMeshRevisionProfilesResponse]`
- New function `*ManagedClustersClient.NewListMeshUpgradeProfilesPager(string, string, *ManagedClustersClientListMeshUpgradeProfilesOptions) *runtime.Pager[ManagedClustersClientListMeshUpgradeProfilesResponse]`
- New function `NewTrustedAccessRoleBindingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TrustedAccessRoleBindingsClient, error)`
- New function `*TrustedAccessRoleBindingsClient.CreateOrUpdate(context.Context, string, string, string, TrustedAccessRoleBinding, *TrustedAccessRoleBindingsClientCreateOrUpdateOptions) (TrustedAccessRoleBindingsClientCreateOrUpdateResponse, error)`
- New function `*TrustedAccessRoleBindingsClient.Delete(context.Context, string, string, string, *TrustedAccessRoleBindingsClientDeleteOptions) (TrustedAccessRoleBindingsClientDeleteResponse, error)`
- New function `*TrustedAccessRoleBindingsClient.Get(context.Context, string, string, string, *TrustedAccessRoleBindingsClientGetOptions) (TrustedAccessRoleBindingsClientGetResponse, error)`
- New function `*TrustedAccessRoleBindingsClient.NewListPager(string, string, *TrustedAccessRoleBindingsClientListOptions) *runtime.Pager[TrustedAccessRoleBindingsClientListResponse]`
- New function `NewTrustedAccessRolesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TrustedAccessRolesClient, error)`
- New function `*TrustedAccessRolesClient.NewListPager(string, *TrustedAccessRolesClientListOptions) *runtime.Pager[TrustedAccessRolesClientListResponse]`
- New struct `AgentPoolNetworkProfile`
- New struct `AgentPoolSecurityProfile`
- New struct `AgentPoolWindowsProfile`
- New struct `CompatibleVersions`
- New struct `GuardrailsAvailableVersion`
- New struct `GuardrailsAvailableVersionsList`
- New struct `GuardrailsAvailableVersionsProperties`
- New struct `GuardrailsProfile`
- New struct `IPTag`
- New struct `IstioCertificateAuthority`
- New struct `IstioComponents`
- New struct `IstioIngressGateway`
- New struct `IstioPluginCertificateAuthority`
- New struct `IstioServiceMesh`
- New struct `Machine`
- New struct `MachineIPAddress`
- New struct `MachineListResult`
- New struct `MachineNetworkProperties`
- New struct `MachineProperties`
- New struct `ManagedClusterAzureMonitorProfileAppMonitoring`
- New struct `ManagedClusterAzureMonitorProfileAppMonitoringOpenTelemetryMetrics`
- New struct `ManagedClusterAzureMonitorProfileContainerInsights`
- New struct `ManagedClusterAzureMonitorProfileLogs`
- New struct `ManagedClusterAzureMonitorProfileWindowsHostLogs`
- New struct `ManagedClusterCostAnalysis`
- New struct `ManagedClusterIngressProfile`
- New struct `ManagedClusterIngressProfileWebAppRouting`
- New struct `ManagedClusterMetricsProfile`
- New struct `ManagedClusterNodeResourceGroupProfile`
- New struct `ManagedClusterPropertiesForSnapshot`
- New struct `ManagedClusterSecurityProfileImageIntegrity`
- New struct `ManagedClusterSecurityProfileNodeRestriction`
- New struct `ManagedClusterSnapshot`
- New struct `ManagedClusterSnapshotListResult`
- New struct `ManagedClusterSnapshotProperties`
- New struct `MeshRevision`
- New struct `MeshRevisionProfile`
- New struct `MeshRevisionProfileList`
- New struct `MeshRevisionProfileProperties`
- New struct `MeshUpgradeProfile`
- New struct `MeshUpgradeProfileList`
- New struct `MeshUpgradeProfileProperties`
- New struct `NetworkMonitoring`
- New struct `NetworkProfileForSnapshot`
- New struct `NetworkProfileKubeProxyConfig`
- New struct `NetworkProfileKubeProxyConfigIpvsConfig`
- New struct `PortRange`
- New struct `ServiceMeshProfile`
- New struct `TrustedAccessRole`
- New struct `TrustedAccessRoleBinding`
- New struct `TrustedAccessRoleBindingListResult`
- New struct `TrustedAccessRoleBindingProperties`
- New struct `TrustedAccessRoleListResult`
- New struct `TrustedAccessRoleRule`
- New field `IgnorePodDisruptionBudget` in struct `AgentPoolsClientBeginDeleteOptions`
- New field `EnableVnetIntegration`, `SubnetID` in struct `ManagedClusterAPIServerAccessProfile`
- New field `CapacityReservationGroupID`, `EnableCustomCATrust`, `MessageOfTheDay`, `NetworkProfile`, `SecurityProfile`, `WindowsProfile` in struct `ManagedClusterAgentPoolProfile`
- New field `CapacityReservationGroupID`, `EnableCustomCATrust`, `MessageOfTheDay`, `NetworkProfile`, `SecurityProfile`, `WindowsProfile` in struct `ManagedClusterAgentPoolProfileProperties`
- New field `Logs` in struct `ManagedClusterAzureMonitorProfile`
- New field `AppMonitoringOpenTelemetryMetrics` in struct `ManagedClusterAzureMonitorProfileMetrics`
- New field `EffectiveNoProxy` in struct `ManagedClusterHTTPProxyConfig`
- New field `BackendPoolType` in struct `ManagedClusterLoadBalancerProfile`
- New field `CreationData`, `EnableNamespaceResources`, `GuardrailsProfile`, `IngressProfile`, `MetricsProfile`, `NodeResourceGroupProfile`, `ResourceUID`, `ServiceMeshProfile` in struct `ManagedClusterProperties`
- New field `CustomCATrustCertificates`, `ImageIntegrity`, `NodeRestriction` in struct `ManagedClusterSecurityProfile`
- New field `Version` in struct `ManagedClusterStorageProfileDiskCSIDriver`
- New field `IgnorePodDisruptionBudget` in struct `ManagedClustersClientBeginDeleteOptions`
- New field `KubeProxyConfig`, `Monitoring` in struct `NetworkProfile`


## 4.3.0 (2023-08-25)
### Features Added

- New struct `ClusterUpgradeSettings`
- New struct `UpgradeOverrideSettings`
- New field `UpgradeSettings` in struct `ManagedClusterProperties`


## 4.2.0 (2023-08-25)
### Features Added

- New enum type `NodeOSUpgradeChannel` with values `NodeOSUpgradeChannelNodeImage`, `NodeOSUpgradeChannelNone`, `NodeOSUpgradeChannelUnmanaged`
- New struct `DelegatedResource`
- New struct `ManagedClusterWorkloadAutoScalerProfileVerticalPodAutoscaler`
- New field `DrainTimeoutInMinutes` in struct `AgentPoolUpgradeSettings`
- New field `NodeOSUpgradeChannel` in struct `ManagedClusterAutoUpgradeProfile`
- New field `DelegatedResources` in struct `ManagedClusterIdentity`
- New field `VerticalPodAutoscaler` in struct `ManagedClusterWorkloadAutoScalerProfile`


## 4.1.0 (2023-07-28)
### Features Added

- New enum type `Type` with values `TypeFirst`, `TypeFourth`, `TypeLast`, `TypeSecond`, `TypeThird`
- New struct `AbsoluteMonthlySchedule`
- New struct `DailySchedule`
- New struct `DateSpan`
- New struct `MaintenanceWindow`
- New struct `RelativeMonthlySchedule`
- New struct `Schedule`
- New struct `WeeklySchedule`
- New field `MaintenanceWindow` in struct `MaintenanceConfigurationProperties`


## 4.0.0 (2023-05-26)
### Breaking Changes

- Field `DockerBridgeCidr` of struct `NetworkProfile` has been removed

### Features Added

- New value `OSSKUAzureLinux` added to enum type `OSSKU`


## 3.0.0 (2023-04-28)
### Breaking Changes

- Const `ManagedClusterSKUNameBasic` from type alias `ManagedClusterSKUName` has been removed
- Const `ManagedClusterSKUTierPaid` from type alias `ManagedClusterSKUTier` has been removed

### Features Added

- New value `ManagedClusterSKUTierPremium` added to enum type `ManagedClusterSKUTier`
- New value `NetworkPolicyCilium` added to enum type `NetworkPolicy`
- New enum type `KubernetesSupportPlan` with values `KubernetesSupportPlanAKSLongTermSupport`, `KubernetesSupportPlanKubernetesOfficial`
- New enum type `NetworkDataplane` with values `NetworkDataplaneAzure`, `NetworkDataplaneCilium`
- New enum type `NetworkPluginMode` with values `NetworkPluginModeOverlay`
- New function `*ManagedClustersClient.ListKubernetesVersions(context.Context, string, *ManagedClustersClientListKubernetesVersionsOptions) (ManagedClustersClientListKubernetesVersionsResponse, error)`
- New struct `KubernetesPatchVersion`
- New struct `KubernetesVersion`
- New struct `KubernetesVersionCapabilities`
- New struct `KubernetesVersionListResult`
- New struct `ManagedClusterSecurityProfileImageCleaner`
- New struct `ManagedClusterSecurityProfileWorkloadIdentity`
- New field `SupportPlan` in struct `ManagedClusterProperties`
- New field `ImageCleaner` in struct `ManagedClusterSecurityProfile`
- New field `WorkloadIdentity` in struct `ManagedClusterSecurityProfile`
- New field `NetworkDataplane` in struct `NetworkProfile`
- New field `NetworkPluginMode` in struct `NetworkProfile`


## 2.4.0 (2023-03-24)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New value `ManagedClusterSKUNameBase` added to enum type `ManagedClusterSKUName`
- New value `ManagedClusterSKUTierStandard` added to enum type `ManagedClusterSKUTier`
- New function `*AgentPoolsClient.BeginAbortLatestOperation(context.Context, string, string, string, *AgentPoolsClientBeginAbortLatestOperationOptions) (*runtime.Poller[AgentPoolsClientAbortLatestOperationResponse], error)`
- New function `*ManagedClustersClient.BeginAbortLatestOperation(context.Context, string, string, *ManagedClustersClientBeginAbortLatestOperationOptions) (*runtime.Poller[ManagedClustersClientAbortLatestOperationResponse], error)`
- New struct `ManagedClusterAzureMonitorProfile`
- New struct `ManagedClusterAzureMonitorProfileKubeStateMetrics`
- New struct `ManagedClusterAzureMonitorProfileMetrics`
- New field `AzureMonitorProfile` in struct `ManagedClusterProperties`


## 2.3.0 (2023-01-27)
### Features Added

- New value `ManagedClusterPodIdentityProvisioningStateCanceled`, `ManagedClusterPodIdentityProvisioningStateSucceeded` added to type alias `ManagedClusterPodIdentityProvisioningState`
- New value `PrivateEndpointConnectionProvisioningStateCanceled` added to type alias `PrivateEndpointConnectionProvisioningState`
- New struct `ManagedClusterWorkloadAutoScalerProfile`
- New struct `ManagedClusterWorkloadAutoScalerProfileKeda`
- New field `WorkloadAutoScalerProfile` in struct `ManagedClusterProperties`
- New field `Location` in struct `ManagedClustersClientGetCommandResultResponse`


## 2.2.0 (2022-10-26)
### Features Added

- New function `*ManagedClustersClient.BeginRotateServiceAccountSigningKeys(context.Context, string, string, *ManagedClustersClientBeginRotateServiceAccountSigningKeysOptions) (*runtime.Poller[ManagedClustersClientRotateServiceAccountSigningKeysResponse], error)`
- New struct `ManagedClusterOIDCIssuerProfile`
- New struct `ManagedClusterStorageProfileBlobCSIDriver`
- New struct `ManagedClustersClientBeginRotateServiceAccountSigningKeysOptions`
- New struct `ManagedClustersClientRotateServiceAccountSigningKeysResponse`
- New field `BlobCSIDriver` in struct `ManagedClusterStorageProfile`
- New field `OidcIssuerProfile` in struct `ManagedClusterProperties`


## 2.1.0 (2022-08-25)
### Features Added

- New const `OSSKUWindows2019`
- New const `OSSKUWindows2022`


## 2.0.0 (2022-07-22)
### Breaking Changes

- Struct `ManagedClusterSecurityProfileAzureDefender` has been removed
- Field `AzureDefender` of struct `ManagedClusterSecurityProfile` has been removed

### Features Added

- New const `KeyVaultNetworkAccessTypesPrivate`
- New const `NetworkPluginNone`
- New const `KeyVaultNetworkAccessTypesPublic`
- New function `PossibleKeyVaultNetworkAccessTypesValues() []KeyVaultNetworkAccessTypes`
- New struct `AzureKeyVaultKms`
- New struct `ManagedClusterSecurityProfileDefender`
- New struct `ManagedClusterSecurityProfileDefenderSecurityMonitoring`
- New field `HostGroupID` in struct `ManagedClusterAgentPoolProfileProperties`
- New field `HostGroupID` in struct `ManagedClusterAgentPoolProfile`
- New field `AzureKeyVaultKms` in struct `ManagedClusterSecurityProfile`
- New field `Defender` in struct `ManagedClusterSecurityProfile`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
