# Release History

## 2.0.0 (2024-05-09)
### Breaking Changes

- Function `*WorkloadNetworksClient.BeginUpdateDNSService` parameter(s) have been changed from `(context.Context, string, string, string, WorkloadNetworkDNSService, *WorkloadNetworksClientBeginUpdateDNSServiceOptions)` to `(context.Context, string, string, string, WorkloadNetworkDNSServiceUpdate, *WorkloadNetworksClientBeginUpdateDNSServiceOptions)`
- Function `*WorkloadNetworksClient.BeginUpdateDNSZone` parameter(s) have been changed from `(context.Context, string, string, string, WorkloadNetworkDNSZone, *WorkloadNetworksClientBeginUpdateDNSZoneOptions)` to `(context.Context, string, string, string, WorkloadNetworkDNSZoneUpdate, *WorkloadNetworksClientBeginUpdateDNSZoneOptions)`
- Function `*WorkloadNetworksClient.BeginUpdateDhcp` parameter(s) have been changed from `(context.Context, string, string, string, WorkloadNetworkDhcp, *WorkloadNetworksClientBeginUpdateDhcpOptions)` to `(context.Context, string, string, string, WorkloadNetworkDhcpUpdate, *WorkloadNetworksClientBeginUpdateDhcpOptions)`
- Function `*WorkloadNetworksClient.BeginUpdatePortMirroring` parameter(s) have been changed from `(context.Context, string, string, string, WorkloadNetworkPortMirroring, *WorkloadNetworksClientBeginUpdatePortMirroringOptions)` to `(context.Context, string, string, string, WorkloadNetworkPortMirroringUpdate, *WorkloadNetworksClientBeginUpdatePortMirroringOptions)`
- Function `*WorkloadNetworksClient.BeginUpdateSegments` parameter(s) have been changed from `(context.Context, string, string, string, WorkloadNetworkSegment, *WorkloadNetworksClientBeginUpdateSegmentsOptions)` to `(context.Context, string, string, string, WorkloadNetworkSegmentUpdate, *WorkloadNetworksClientBeginUpdateSegmentsOptions)`
- Function `*WorkloadNetworksClient.BeginUpdateVMGroup` parameter(s) have been changed from `(context.Context, string, string, string, WorkloadNetworkVMGroup, *WorkloadNetworksClientBeginUpdateVMGroupOptions)` to `(context.Context, string, string, string, WorkloadNetworkVMGroupUpdate, *WorkloadNetworksClientBeginUpdateVMGroupOptions)`
- Function `*WorkloadNetworksClient.Get` parameter(s) have been changed from `(context.Context, string, string, WorkloadNetworkName, *WorkloadNetworksClientGetOptions)` to `(context.Context, string, string, *WorkloadNetworksClientGetOptions)`
- Type of `Operation.Origin` has been changed from `*string` to `*Origin`
- Type of `PrivateCloud.Identity` has been changed from `*PrivateCloudIdentity` to `*SystemAssignedServiceIdentity`
- Type of `PrivateCloudUpdate.Identity` has been changed from `*PrivateCloudIdentity` to `*SystemAssignedServiceIdentity`
- Type of `WorkloadNetworkDNSZoneProperties.DNSServices` has been changed from `*int64` to `*int32`
- Type of `WorkloadNetworkDhcpServer.LeaseTime` has been changed from `*int64` to `*int32`
- Enum `ResourceIdentityType` has been removed
- Enum `WorkloadNetworkName` has been removed
- Struct `AddonList` has been removed
- Struct `CloudLinkList` has been removed
- Struct `ClusterList` has been removed
- Struct `DatastoreList` has been removed
- Struct `ExpressRouteAuthorizationList` has been removed
- Struct `GlobalReachConnectionList` has been removed
- Struct `HcxEnterpriseSiteList` has been removed
- Struct `LogSpecification` has been removed
- Struct `MetricDimension` has been removed
- Struct `MetricSpecification` has been removed
- Struct `OperationList` has been removed
- Struct `OperationProperties` has been removed
- Struct `PlacementPoliciesList` has been removed
- Struct `PrivateCloudIdentity` has been removed
- Struct `PrivateCloudList` has been removed
- Struct `ScriptCmdletsList` has been removed
- Struct `ScriptExecutionsList` has been removed
- Struct `ScriptPackagesList` has been removed
- Struct `ServiceSpecification` has been removed
- Struct `VirtualMachinesList` has been removed
- Struct `WorkloadNetworkDNSServicesList` has been removed
- Struct `WorkloadNetworkDNSZonesList` has been removed
- Struct `WorkloadNetworkDhcpList` has been removed
- Struct `WorkloadNetworkGatewayList` has been removed
- Struct `WorkloadNetworkList` has been removed
- Struct `WorkloadNetworkPortMirroringList` has been removed
- Struct `WorkloadNetworkPublicIPsList` has been removed
- Struct `WorkloadNetworkSegmentsList` has been removed
- Struct `WorkloadNetworkVMGroupsList` has been removed
- Struct `WorkloadNetworkVirtualMachinesList` has been removed
- Field `AddonList` of struct `AddonsClientListResponse` has been removed
- Field `ExpressRouteAuthorizationList` of struct `AuthorizationsClientListResponse` has been removed
- Field `CloudLinkList` of struct `CloudLinksClientListResponse` has been removed
- Field `ClusterList` of struct `ClustersClientListResponse` has been removed
- Field `DatastoreList` of struct `DatastoresClientListResponse` has been removed
- Field `GlobalReachConnectionList` of struct `GlobalReachConnectionsClientListResponse` has been removed
- Field `HcxEnterpriseSiteList` of struct `HcxEnterpriseSitesClientListResponse` has been removed
- Field `Properties` of struct `Operation` has been removed
- Field `OperationList` of struct `OperationsClientListResponse` has been removed
- Field `PlacementPoliciesList` of struct `PlacementPoliciesClientListResponse` has been removed
- Field `PrivateCloudList` of struct `PrivateCloudsClientListInSubscriptionResponse` has been removed
- Field `PrivateCloudList` of struct `PrivateCloudsClientListResponse` has been removed
- Field `ScriptCmdletsList` of struct `ScriptCmdletsClientListResponse` has been removed
- Field `ScriptExecutionsList` of struct `ScriptExecutionsClientListResponse` has been removed
- Field `ScriptPackagesList` of struct `ScriptPackagesClientListResponse` has been removed
- Field `VirtualMachinesList` of struct `VirtualMachinesClientListResponse` has been removed
- Field `WorkloadNetworkDNSServicesList` of struct `WorkloadNetworksClientListDNSServicesResponse` has been removed
- Field `WorkloadNetworkDNSZonesList` of struct `WorkloadNetworksClientListDNSZonesResponse` has been removed
- Field `WorkloadNetworkDhcpList` of struct `WorkloadNetworksClientListDhcpResponse` has been removed
- Field `WorkloadNetworkGatewayList` of struct `WorkloadNetworksClientListGatewaysResponse` has been removed
- Field `WorkloadNetworkPortMirroringList` of struct `WorkloadNetworksClientListPortMirroringResponse` has been removed
- Field `WorkloadNetworkPublicIPsList` of struct `WorkloadNetworksClientListPublicIPsResponse` has been removed
- Field `WorkloadNetworkList` of struct `WorkloadNetworksClientListResponse` has been removed
- Field `WorkloadNetworkSegmentsList` of struct `WorkloadNetworksClientListSegmentsResponse` has been removed
- Field `WorkloadNetworkVMGroupsList` of struct `WorkloadNetworksClientListVMGroupsResponse` has been removed
- Field `WorkloadNetworkVirtualMachinesList` of struct `WorkloadNetworksClientListVirtualMachinesResponse` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `CloudLinkProvisioningState` with values `CloudLinkProvisioningStateCanceled`, `CloudLinkProvisioningStateFailed`, `CloudLinkProvisioningStateSucceeded`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `DNSZoneType` with values `DNSZoneTypePrivate`, `DNSZoneTypePublic`
- New enum type `HcxEnterpriseSiteProvisioningState` with values `HcxEnterpriseSiteProvisioningStateCanceled`, `HcxEnterpriseSiteProvisioningStateFailed`, `HcxEnterpriseSiteProvisioningStateSucceeded`
- New enum type `IscsiPathProvisioningState` with values `IscsiPathProvisioningStateBuilding`, `IscsiPathProvisioningStateCanceled`, `IscsiPathProvisioningStateDeleting`, `IscsiPathProvisioningStateFailed`, `IscsiPathProvisioningStatePending`, `IscsiPathProvisioningStateSucceeded`, `IscsiPathProvisioningStateUpdating`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `SKUTier` with values `SKUTierBasic`, `SKUTierFree`, `SKUTierPremium`, `SKUTierStandard`
- New enum type `ScriptCmdletAudience` with values `ScriptCmdletAudienceAny`, `ScriptCmdletAudienceAutomation`
- New enum type `ScriptCmdletProvisioningState` with values `ScriptCmdletProvisioningStateCanceled`, `ScriptCmdletProvisioningStateFailed`, `ScriptCmdletProvisioningStateSucceeded`
- New enum type `ScriptPackageProvisioningState` with values `ScriptPackageProvisioningStateCanceled`, `ScriptPackageProvisioningStateFailed`, `ScriptPackageProvisioningStateSucceeded`
- New enum type `SystemAssignedServiceIdentityType` with values `SystemAssignedServiceIdentityTypeNone`, `SystemAssignedServiceIdentityTypeSystemAssigned`
- New enum type `VirtualMachineProvisioningState` with values `VirtualMachineProvisioningStateCanceled`, `VirtualMachineProvisioningStateFailed`, `VirtualMachineProvisioningStateSucceeded`
- New enum type `WorkloadNetworkProvisioningState` with values `WorkloadNetworkProvisioningStateBuilding`, `WorkloadNetworkProvisioningStateCanceled`, `WorkloadNetworkProvisioningStateDeleting`, `WorkloadNetworkProvisioningStateFailed`, `WorkloadNetworkProvisioningStateSucceeded`, `WorkloadNetworkProvisioningStateUpdating`
- New function `*ClientFactory.NewIscsiPathsClient() *IscsiPathsClient`
- New function `NewIscsiPathsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*IscsiPathsClient, error)`
- New function `*IscsiPathsClient.BeginCreateOrUpdate(context.Context, string, string, IscsiPath, *IscsiPathsClientBeginCreateOrUpdateOptions) (*runtime.Poller[IscsiPathsClientCreateOrUpdateResponse], error)`
- New function `*IscsiPathsClient.BeginDelete(context.Context, string, string, *IscsiPathsClientBeginDeleteOptions) (*runtime.Poller[IscsiPathsClientDeleteResponse], error)`
- New function `*IscsiPathsClient.Get(context.Context, string, string, *IscsiPathsClientGetOptions) (IscsiPathsClientGetResponse, error)`
- New function `*IscsiPathsClient.NewListByPrivateCloudPager(string, string, *IscsiPathsClientListByPrivateCloudOptions) *runtime.Pager[IscsiPathsClientListByPrivateCloudResponse]`
- New function `*WorkloadNetworkDhcpEntityUpdate.GetWorkloadNetworkDhcpEntityUpdate() *WorkloadNetworkDhcpEntityUpdate`
- New function `*WorkloadNetworkDhcpRelayUpdate.GetWorkloadNetworkDhcpEntityUpdate() *WorkloadNetworkDhcpEntityUpdate`
- New function `*WorkloadNetworkDhcpServerUpdate.GetWorkloadNetworkDhcpEntityUpdate() *WorkloadNetworkDhcpEntityUpdate`
- New struct `AddonListResult`
- New struct `CloudLinkListResult`
- New struct `ClusterListResult`
- New struct `DatastoreListResult`
- New struct `ElasticSanVolume`
- New struct `ExpressRouteAuthorizationListResult`
- New struct `GlobalReachConnectionListResult`
- New struct `HcxEnterpriseSiteListResult`
- New struct `IscsiPath`
- New struct `IscsiPathListResult`
- New struct `IscsiPathProperties`
- New struct `OperationListResult`
- New struct `PlacementPolicyListResult`
- New struct `PrivateCloudListResult`
- New struct `ScriptCmdletListResult`
- New struct `ScriptExecutionListResult`
- New struct `ScriptPackageListResult`
- New struct `SystemAssignedServiceIdentity`
- New struct `SystemData`
- New struct `VirtualMachineListResult`
- New struct `WorkloadNetworkDNSServiceListResult`
- New struct `WorkloadNetworkDNSServiceUpdate`
- New struct `WorkloadNetworkDNSZoneListResult`
- New struct `WorkloadNetworkDNSZoneUpdate`
- New struct `WorkloadNetworkDhcpListResult`
- New struct `WorkloadNetworkDhcpRelayUpdate`
- New struct `WorkloadNetworkDhcpServerUpdate`
- New struct `WorkloadNetworkDhcpUpdate`
- New struct `WorkloadNetworkGatewayListResult`
- New struct `WorkloadNetworkListResult`
- New struct `WorkloadNetworkPortMirroringListResult`
- New struct `WorkloadNetworkPortMirroringUpdate`
- New struct `WorkloadNetworkProperties`
- New struct `WorkloadNetworkPublicIPListResult`
- New struct `WorkloadNetworkSegmentListResult`
- New struct `WorkloadNetworkSegmentUpdate`
- New struct `WorkloadNetworkVMGroupListResult`
- New struct `WorkloadNetworkVMGroupUpdate`
- New struct `WorkloadNetworkVirtualMachineListResult`
- New field `SystemData` in struct `Addon`
- New anonymous field `AddonListResult` in struct `AddonsClientListResponse`
- New anonymous field `ExpressRouteAuthorizationListResult` in struct `AuthorizationsClientListResponse`
- New field `SystemData` in struct `CloudLink`
- New field `ProvisioningState` in struct `CloudLinkProperties`
- New anonymous field `CloudLinkListResult` in struct `CloudLinksClientListResponse`
- New field `SystemData` in struct `Cluster`
- New field `VsanDatastoreName` in struct `ClusterProperties`
- New field `SKU` in struct `ClusterUpdate`
- New anonymous field `ClusterListResult` in struct `ClustersClientListResponse`
- New field `SystemData` in struct `Datastore`
- New field `ElasticSanVolume` in struct `DatastoreProperties`
- New anonymous field `DatastoreListResult` in struct `DatastoresClientListResponse`
- New field `HcxCloudManagerIP`, `NsxtManagerIP`, `VcenterIP` in struct `Endpoints`
- New field `SystemData` in struct `ExpressRouteAuthorization`
- New field `SystemData` in struct `GlobalReachConnection`
- New anonymous field `GlobalReachConnectionListResult` in struct `GlobalReachConnectionsClientListResponse`
- New field `SystemData` in struct `HcxEnterpriseSite`
- New field `ProvisioningState` in struct `HcxEnterpriseSiteProperties`
- New anonymous field `HcxEnterpriseSiteListResult` in struct `HcxEnterpriseSitesClientListResponse`
- New field `VsanDatastoreName` in struct `ManagementCluster`
- New field `ActionType` in struct `Operation`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `PlacementPolicyListResult` in struct `PlacementPoliciesClientListResponse`
- New field `SystemData` in struct `PlacementPolicy`
- New field `SystemData` in struct `PrivateCloud`
- New field `DNSZoneType`, `VirtualNetworkID` in struct `PrivateCloudProperties`
- New field `SKU` in struct `PrivateCloudUpdate`
- New field `DNSZoneType` in struct `PrivateCloudUpdateProperties`
- New anonymous field `PrivateCloudListResult` in struct `PrivateCloudsClientListInSubscriptionResponse`
- New anonymous field `PrivateCloudListResult` in struct `PrivateCloudsClientListResponse`
- New field `Capacity`, `Family`, `Size`, `Tier` in struct `SKU`
- New field `SystemData` in struct `ScriptCmdlet`
- New field `Audience`, `ProvisioningState` in struct `ScriptCmdletProperties`
- New anonymous field `ScriptCmdletListResult` in struct `ScriptCmdletsClientListResponse`
- New field `SystemData` in struct `ScriptExecution`
- New anonymous field `ScriptExecutionListResult` in struct `ScriptExecutionsClientListResponse`
- New field `SystemData` in struct `ScriptPackage`
- New field `ProvisioningState` in struct `ScriptPackageProperties`
- New anonymous field `ScriptPackageListResult` in struct `ScriptPackagesClientListResponse`
- New field `SystemData` in struct `VirtualMachine`
- New field `ProvisioningState` in struct `VirtualMachineProperties`
- New anonymous field `VirtualMachineListResult` in struct `VirtualMachinesClientListResponse`
- New field `Properties`, `SystemData` in struct `WorkloadNetwork`
- New field `SystemData` in struct `WorkloadNetworkDNSService`
- New field `SystemData` in struct `WorkloadNetworkDNSZone`
- New field `SystemData` in struct `WorkloadNetworkDhcp`
- New field `SystemData` in struct `WorkloadNetworkGateway`
- New field `ProvisioningState` in struct `WorkloadNetworkGatewayProperties`
- New field `SystemData` in struct `WorkloadNetworkPortMirroring`
- New field `SystemData` in struct `WorkloadNetworkPublicIP`
- New field `SystemData` in struct `WorkloadNetworkSegment`
- New field `SystemData` in struct `WorkloadNetworkVMGroup`
- New field `SystemData` in struct `WorkloadNetworkVirtualMachine`
- New field `ProvisioningState` in struct `WorkloadNetworkVirtualMachineProperties`
- New anonymous field `WorkloadNetworkDNSServiceListResult` in struct `WorkloadNetworksClientListDNSServicesResponse`
- New anonymous field `WorkloadNetworkDNSZoneListResult` in struct `WorkloadNetworksClientListDNSZonesResponse`
- New anonymous field `WorkloadNetworkDhcpListResult` in struct `WorkloadNetworksClientListDhcpResponse`
- New anonymous field `WorkloadNetworkGatewayListResult` in struct `WorkloadNetworksClientListGatewaysResponse`
- New anonymous field `WorkloadNetworkPortMirroringListResult` in struct `WorkloadNetworksClientListPortMirroringResponse`
- New anonymous field `WorkloadNetworkPublicIPListResult` in struct `WorkloadNetworksClientListPublicIPsResponse`
- New anonymous field `WorkloadNetworkListResult` in struct `WorkloadNetworksClientListResponse`
- New anonymous field `WorkloadNetworkSegmentListResult` in struct `WorkloadNetworksClientListSegmentsResponse`
- New anonymous field `WorkloadNetworkVMGroupListResult` in struct `WorkloadNetworksClientListVMGroupsResponse`
- New anonymous field `WorkloadNetworkVirtualMachineListResult` in struct `WorkloadNetworksClientListVirtualMachinesResponse`


## 1.4.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.3.0 (2023-08-25)

### Features Added

- New field `ExtendedNetworkBlocks` in struct `PrivateCloudProperties`
- New field `ExtendedNetworkBlocks` in struct `PrivateCloudUpdateProperties`


## 1.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.2.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.1.0 (2022-10-13)

### Features Added

- New const `ExpressRouteAuthorizationProvisioningStateCanceled`
- New const `AffinityStrengthShould`
- New const `AddonTypeArc`
- New const `PrivateCloudProvisioningStateCanceled`
- New const `NsxPublicIPQuotaRaisedEnumEnabled`
- New const `AzureHybridBenefitTypeNone`
- New const `WorkloadNetworkPublicIPProvisioningStateCanceled`
- New const `WorkloadNetworkDNSServiceProvisioningStateCanceled`
- New const `WorkloadNetworkSegmentProvisioningStateCanceled`
- New const `WorkloadNetworkDNSZoneProvisioningStateCanceled`
- New const `WorkloadNetworkNameDefault`
- New const `PlacementPolicyProvisioningStateCanceled`
- New const `WorkloadNetworkDhcpProvisioningStateCanceled`
- New const `WorkloadNetworkPortMirroringProvisioningStateCanceled`
- New const `WorkloadNetworkVMGroupProvisioningStateCanceled`
- New const `NsxPublicIPQuotaRaisedEnumDisabled`
- New const `DatastoreProvisioningStateCanceled`
- New const `AzureHybridBenefitTypeSQLHost`
- New const `AddonProvisioningStateCanceled`
- New const `ClusterProvisioningStateCanceled`
- New const `AffinityStrengthMust`
- New const `GlobalReachConnectionProvisioningStateCanceled`
- New const `ScriptExecutionProvisioningStateCanceled`
- New type alias `NsxPublicIPQuotaRaisedEnum`
- New type alias `AzureHybridBenefitType`
- New type alias `AffinityStrength`
- New type alias `WorkloadNetworkName`
- New function `PossibleAzureHybridBenefitTypeValues() []AzureHybridBenefitType`
- New function `*WorkloadNetworksClient.Get(context.Context, string, string, WorkloadNetworkName, *WorkloadNetworksClientGetOptions) (WorkloadNetworksClientGetResponse, error)`
- New function `*ClustersClient.ListZones(context.Context, string, string, string, *ClustersClientListZonesOptions) (ClustersClientListZonesResponse, error)`
- New function `PossibleNsxPublicIPQuotaRaisedEnumValues() []NsxPublicIPQuotaRaisedEnum`
- New function `PossibleWorkloadNetworkNameValues() []WorkloadNetworkName`
- New function `*WorkloadNetworksClient.NewListPager(string, string, *WorkloadNetworksClientListOptions) *runtime.Pager[WorkloadNetworksClientListResponse]`
- New function `*AddonArcProperties.GetAddonProperties() *AddonProperties`
- New function `PossibleAffinityStrengthValues() []AffinityStrength`
- New struct `AddonArcProperties`
- New struct `ClusterZone`
- New struct `ClusterZoneList`
- New struct `ClustersClientListZonesOptions`
- New struct `ClustersClientListZonesResponse`
- New struct `WorkloadNetwork`
- New struct `WorkloadNetworkList`
- New struct `WorkloadNetworksClientGetOptions`
- New struct `WorkloadNetworksClientGetResponse`
- New struct `WorkloadNetworksClientListOptions`
- New struct `WorkloadNetworksClientListResponse`
- New field `AffinityStrength` in struct `PlacementPolicyUpdateProperties`
- New field `AzureHybridBenefitType` in struct `PlacementPolicyUpdateProperties`
- New field `AutoDetectedKeyVersion` in struct `EncryptionKeyVaultProperties`
- New field `SKU` in struct `LocationsClientCheckTrialAvailabilityOptions`
- New field `AzureHybridBenefitType` in struct `VMHostPlacementPolicyProperties`
- New field `AffinityStrength` in struct `VMHostPlacementPolicyProperties`
- New field `NsxPublicIPQuotaRaised` in struct `PrivateCloudProperties`
- New field `Company` in struct `ScriptPackageProperties`
- New field `URI` in struct `ScriptPackageProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/avs/armavs` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).