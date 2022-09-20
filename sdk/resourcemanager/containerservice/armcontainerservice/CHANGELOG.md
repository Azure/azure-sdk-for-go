# Release History

## 2.2.0-beta.1 (2022-09-17)
### Features Added

- New const `FleetMemberProvisioningStateFailed`
- New const `ControlledValuesRequestsOnly`
- New const `FleetMemberProvisioningStateCanceled`
- New const `TrustedAccessRoleBindingProvisioningStateUpdating`
- New const `LevelWarning`
- New const `UpdateModeOff`
- New const `FleetProvisioningStateCreating`
- New const `UpdateModeInitial`
- New const `FleetMemberProvisioningStateLeaving`
- New const `IpvsSchedulerRoundRobin`
- New const `IpvsSchedulerLeastConnection`
- New const `FleetMemberProvisioningStateJoining`
- New const `FleetProvisioningStateFailed`
- New const `BackendPoolTypeNodeIPConfiguration`
- New const `OSSKUMariner`
- New const `FleetProvisioningStateSucceeded`
- New const `ModeIPVS`
- New const `FleetProvisioningStateUpdating`
- New const `NetworkPluginModeOverlay`
- New const `FleetMemberProvisioningStateUpdating`
- New const `TrustedAccessRoleBindingProvisioningStateFailed`
- New const `TrustedAccessRoleBindingProvisioningStateDeleting`
- New const `ModeIPTABLES`
- New const `UpdateModeAuto`
- New const `BackendPoolTypeNodeIP`
- New const `FleetMemberProvisioningStateSucceeded`
- New const `LevelOff`
- New const `UpdateModeRecreate`
- New const `ControlledValuesRequestsAndLimits`
- New const `FleetProvisioningStateDeleting`
- New const `TrustedAccessRoleBindingProvisioningStateSucceeded`
- New const `PublicNetworkAccessSecuredByPerimeter`
- New const `FleetProvisioningStateCanceled`
- New const `SnapshotTypeManagedCluster`
- New const `LevelEnforcement`
- New type alias `NetworkPluginMode`
- New type alias `IpvsScheduler`
- New type alias `FleetMemberProvisioningState`
- New type alias `TrustedAccessRoleBindingProvisioningState`
- New type alias `FleetProvisioningState`
- New type alias `Level`
- New type alias `Mode`
- New type alias `BackendPoolType`
- New type alias `ControlledValues`
- New type alias `UpdateMode`
- New function `PossibleUpdateModeValues() []UpdateMode`
- New function `PossibleModeValues() []Mode`
- New function `*AgentPoolsClient.AbortLatestOperation(context.Context, string, string, string, *AgentPoolsClientAbortLatestOperationOptions) (AgentPoolsClientAbortLatestOperationResponse, error)`
- New function `PossibleNetworkPluginModeValues() []NetworkPluginMode`
- New function `NewTrustedAccessRoleBindingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TrustedAccessRoleBindingsClient, error)`
- New function `PossibleTrustedAccessRoleBindingProvisioningStateValues() []TrustedAccessRoleBindingProvisioningState`
- New function `*ManagedClusterSnapshotsClient.NewListPager(*ManagedClusterSnapshotsClientListOptions) *runtime.Pager[ManagedClusterSnapshotsClientListResponse]`
- New function `NewFleetMembersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FleetMembersClient, error)`
- New function `*ManagedClusterSnapshotsClient.Delete(context.Context, string, string, *ManagedClusterSnapshotsClientDeleteOptions) (ManagedClusterSnapshotsClientDeleteResponse, error)`
- New function `NewManagedClusterSnapshotsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedClusterSnapshotsClient, error)`
- New function `PossibleIpvsSchedulerValues() []IpvsScheduler`
- New function `*FleetMembersClient.Get(context.Context, string, string, string, *FleetMembersClientGetOptions) (FleetMembersClientGetResponse, error)`
- New function `NewTrustedAccessRolesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TrustedAccessRolesClient, error)`
- New function `*ManagedClustersClient.BeginRotateServiceAccountSigningKeys(context.Context, string, string, *ManagedClustersClientBeginRotateServiceAccountSigningKeysOptions) (*runtime.Poller[ManagedClustersClientRotateServiceAccountSigningKeysResponse], error)`
- New function `*TrustedAccessRolesClient.NewListPager(string, *TrustedAccessRolesClientListOptions) *runtime.Pager[TrustedAccessRolesClientListResponse]`
- New function `PossibleBackendPoolTypeValues() []BackendPoolType`
- New function `*TrustedAccessRoleBindingsClient.NewListPager(string, string, *TrustedAccessRoleBindingsClientListOptions) *runtime.Pager[TrustedAccessRoleBindingsClientListResponse]`
- New function `*ManagedClusterSnapshotsClient.NewListByResourceGroupPager(string, *ManagedClusterSnapshotsClientListByResourceGroupOptions) *runtime.Pager[ManagedClusterSnapshotsClientListByResourceGroupResponse]`
- New function `*FleetMembersClient.BeginDelete(context.Context, string, string, string, *FleetMembersClientBeginDeleteOptions) (*runtime.Poller[FleetMembersClientDeleteResponse], error)`
- New function `*ManagedClusterSnapshotsClient.Get(context.Context, string, string, *ManagedClusterSnapshotsClientGetOptions) (ManagedClusterSnapshotsClientGetResponse, error)`
- New function `PossibleFleetProvisioningStateValues() []FleetProvisioningState`
- New function `*TrustedAccessRoleBindingsClient.Delete(context.Context, string, string, string, *TrustedAccessRoleBindingsClientDeleteOptions) (TrustedAccessRoleBindingsClientDeleteResponse, error)`
- New function `*FleetsClient.NewListPager(*FleetsClientListOptions) *runtime.Pager[FleetsClientListResponse]`
- New function `PossibleFleetMemberProvisioningStateValues() []FleetMemberProvisioningState`
- New function `*FleetMembersClient.BeginCreateOrUpdate(context.Context, string, string, string, FleetMember, *FleetMembersClientBeginCreateOrUpdateOptions) (*runtime.Poller[FleetMembersClientCreateOrUpdateResponse], error)`
- New function `NewFleetsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FleetsClient, error)`
- New function `*ManagedClusterSnapshotsClient.CreateOrUpdate(context.Context, string, string, ManagedClusterSnapshot, *ManagedClusterSnapshotsClientCreateOrUpdateOptions) (ManagedClusterSnapshotsClientCreateOrUpdateResponse, error)`
- New function `*ManagedClusterSnapshotsClient.UpdateTags(context.Context, string, string, TagsObject, *ManagedClusterSnapshotsClientUpdateTagsOptions) (ManagedClusterSnapshotsClientUpdateTagsResponse, error)`
- New function `*FleetsClient.ListCredentials(context.Context, string, string, *FleetsClientListCredentialsOptions) (FleetsClientListCredentialsResponse, error)`
- New function `*FleetsClient.NewListByResourceGroupPager(string, *FleetsClientListByResourceGroupOptions) *runtime.Pager[FleetsClientListByResourceGroupResponse]`
- New function `PossibleControlledValuesValues() []ControlledValues`
- New function `*FleetsClient.Update(context.Context, string, string, *FleetsClientUpdateOptions) (FleetsClientUpdateResponse, error)`
- New function `*FleetMembersClient.NewListByFleetPager(string, string, *FleetMembersClientListByFleetOptions) *runtime.Pager[FleetMembersClientListByFleetResponse]`
- New function `*ManagedClustersClient.AbortLatestOperation(context.Context, string, string, *ManagedClustersClientAbortLatestOperationOptions) (ManagedClustersClientAbortLatestOperationResponse, error)`
- New function `*FleetsClient.BeginDelete(context.Context, string, string, *FleetsClientBeginDeleteOptions) (*runtime.Poller[FleetsClientDeleteResponse], error)`
- New function `*TrustedAccessRoleBindingsClient.Get(context.Context, string, string, string, *TrustedAccessRoleBindingsClientGetOptions) (TrustedAccessRoleBindingsClientGetResponse, error)`
- New function `*TrustedAccessRoleBindingsClient.CreateOrUpdate(context.Context, string, string, string, TrustedAccessRoleBinding, *TrustedAccessRoleBindingsClientCreateOrUpdateOptions) (TrustedAccessRoleBindingsClientCreateOrUpdateResponse, error)`
- New function `PossibleLevelValues() []Level`
- New function `*FleetsClient.BeginCreateOrUpdate(context.Context, string, string, Fleet, *FleetsClientBeginCreateOrUpdateOptions) (*runtime.Poller[FleetsClientCreateOrUpdateResponse], error)`
- New function `*FleetsClient.Get(context.Context, string, string, *FleetsClientGetOptions) (FleetsClientGetResponse, error)`
- New struct `AgentPoolWindowsProfile`
- New struct `AgentPoolsClientAbortLatestOperationOptions`
- New struct `AgentPoolsClientAbortLatestOperationResponse`
- New struct `AzureEntityResource`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `Fleet`
- New struct `FleetCredentialResult`
- New struct `FleetCredentialResults`
- New struct `FleetHubProfile`
- New struct `FleetListResult`
- New struct `FleetMember`
- New struct `FleetMemberProperties`
- New struct `FleetMembersClient`
- New struct `FleetMembersClientBeginCreateOrUpdateOptions`
- New struct `FleetMembersClientBeginDeleteOptions`
- New struct `FleetMembersClientCreateOrUpdateResponse`
- New struct `FleetMembersClientDeleteResponse`
- New struct `FleetMembersClientGetOptions`
- New struct `FleetMembersClientGetResponse`
- New struct `FleetMembersClientListByFleetOptions`
- New struct `FleetMembersClientListByFleetResponse`
- New struct `FleetMembersListResult`
- New struct `FleetPatch`
- New struct `FleetProperties`
- New struct `FleetsClient`
- New struct `FleetsClientBeginCreateOrUpdateOptions`
- New struct `FleetsClientBeginDeleteOptions`
- New struct `FleetsClientCreateOrUpdateResponse`
- New struct `FleetsClientDeleteResponse`
- New struct `FleetsClientGetOptions`
- New struct `FleetsClientGetResponse`
- New struct `FleetsClientListByResourceGroupOptions`
- New struct `FleetsClientListByResourceGroupResponse`
- New struct `FleetsClientListCredentialsOptions`
- New struct `FleetsClientListCredentialsResponse`
- New struct `FleetsClientListOptions`
- New struct `FleetsClientListResponse`
- New struct `FleetsClientUpdateOptions`
- New struct `FleetsClientUpdateResponse`
- New struct `GuardrailsProfile`
- New struct `ManagedClusterAzureMonitorProfile`
- New struct `ManagedClusterAzureMonitorProfileKubeStateMetrics`
- New struct `ManagedClusterAzureMonitorProfileMetrics`
- New struct `ManagedClusterIngressProfile`
- New struct `ManagedClusterIngressProfileWebAppRouting`
- New struct `ManagedClusterOIDCIssuerProfile`
- New struct `ManagedClusterPropertiesForSnapshot`
- New struct `ManagedClusterSecurityProfileImageCleaner`
- New struct `ManagedClusterSecurityProfileNodeRestriction`
- New struct `ManagedClusterSecurityProfileWorkloadIdentity`
- New struct `ManagedClusterSnapshot`
- New struct `ManagedClusterSnapshotListResult`
- New struct `ManagedClusterSnapshotProperties`
- New struct `ManagedClusterSnapshotsClient`
- New struct `ManagedClusterSnapshotsClientCreateOrUpdateOptions`
- New struct `ManagedClusterSnapshotsClientCreateOrUpdateResponse`
- New struct `ManagedClusterSnapshotsClientDeleteOptions`
- New struct `ManagedClusterSnapshotsClientDeleteResponse`
- New struct `ManagedClusterSnapshotsClientGetOptions`
- New struct `ManagedClusterSnapshotsClientGetResponse`
- New struct `ManagedClusterSnapshotsClientListByResourceGroupOptions`
- New struct `ManagedClusterSnapshotsClientListByResourceGroupResponse`
- New struct `ManagedClusterSnapshotsClientListOptions`
- New struct `ManagedClusterSnapshotsClientListResponse`
- New struct `ManagedClusterSnapshotsClientUpdateTagsOptions`
- New struct `ManagedClusterSnapshotsClientUpdateTagsResponse`
- New struct `ManagedClusterStorageProfileBlobCSIDriver`
- New struct `ManagedClusterWorkloadAutoScalerProfile`
- New struct `ManagedClusterWorkloadAutoScalerProfileKeda`
- New struct `ManagedClusterWorkloadAutoScalerProfileVerticalPodAutoscaler`
- New struct `ManagedClustersClientAbortLatestOperationOptions`
- New struct `ManagedClustersClientAbortLatestOperationResponse`
- New struct `ManagedClustersClientBeginRotateServiceAccountSigningKeysOptions`
- New struct `ManagedClustersClientRotateServiceAccountSigningKeysResponse`
- New struct `NetworkProfileForSnapshot`
- New struct `NetworkProfileKubeProxyConfig`
- New struct `NetworkProfileKubeProxyConfigIpvsConfig`
- New struct `TrustedAccessRole`
- New struct `TrustedAccessRoleBinding`
- New struct `TrustedAccessRoleBindingListResult`
- New struct `TrustedAccessRoleBindingProperties`
- New struct `TrustedAccessRoleBindingsClient`
- New struct `TrustedAccessRoleBindingsClientCreateOrUpdateOptions`
- New struct `TrustedAccessRoleBindingsClientCreateOrUpdateResponse`
- New struct `TrustedAccessRoleBindingsClientDeleteOptions`
- New struct `TrustedAccessRoleBindingsClientDeleteResponse`
- New struct `TrustedAccessRoleBindingsClientGetOptions`
- New struct `TrustedAccessRoleBindingsClientGetResponse`
- New struct `TrustedAccessRoleBindingsClientListOptions`
- New struct `TrustedAccessRoleBindingsClientListResponse`
- New struct `TrustedAccessRoleListResult`
- New struct `TrustedAccessRoleRule`
- New struct `TrustedAccessRolesClient`
- New struct `TrustedAccessRolesClientListOptions`
- New struct `TrustedAccessRolesClientListResponse`
- New field `ImageCleaner` in struct `ManagedClusterSecurityProfile`
- New field `NodeRestriction` in struct `ManagedClusterSecurityProfile`
- New field `WorkloadIdentity` in struct `ManagedClusterSecurityProfile`
- New field `BackendPoolType` in struct `ManagedClusterLoadBalancerProfile`
- New field `IngressProfile` in struct `ManagedClusterProperties`
- New field `OidcIssuerProfile` in struct `ManagedClusterProperties`
- New field `EnableNamespaceResources` in struct `ManagedClusterProperties`
- New field `CreationData` in struct `ManagedClusterProperties`
- New field `AzureMonitorProfile` in struct `ManagedClusterProperties`
- New field `GuardrailsProfile` in struct `ManagedClusterProperties`
- New field `WorkloadAutoScalerProfile` in struct `ManagedClusterProperties`
- New field `EffectiveNoProxy` in struct `ManagedClusterHTTPProxyConfig`
- New field `Version` in struct `ManagedClusterStorageProfileDiskCSIDriver`
- New field `IgnorePodDisruptionBudget` in struct `ManagedClustersClientBeginDeleteOptions`
- New field `BlobCSIDriver` in struct `ManagedClusterStorageProfile`
- New field `KubeProxyConfig` in struct `NetworkProfile`
- New field `NetworkPluginMode` in struct `NetworkProfile`
- New field `IgnorePodDisruptionBudget` in struct `AgentPoolsClientBeginDeleteOptions`
- New field `EnableVnetIntegration` in struct `ManagedClusterAPIServerAccessProfile`
- New field `SubnetID` in struct `ManagedClusterAPIServerAccessProfile`
- New field `CapacityReservationGroupID` in struct `ManagedClusterAgentPoolProfile`
- New field `WindowsProfile` in struct `ManagedClusterAgentPoolProfile`
- New field `EnableCustomCATrust` in struct `ManagedClusterAgentPoolProfile`
- New field `MessageOfTheDay` in struct `ManagedClusterAgentPoolProfile`
- New field `MessageOfTheDay` in struct `ManagedClusterAgentPoolProfileProperties`
- New field `CapacityReservationGroupID` in struct `ManagedClusterAgentPoolProfileProperties`
- New field `EnableCustomCATrust` in struct `ManagedClusterAgentPoolProfileProperties`
- New field `WindowsProfile` in struct `ManagedClusterAgentPoolProfileProperties`


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
