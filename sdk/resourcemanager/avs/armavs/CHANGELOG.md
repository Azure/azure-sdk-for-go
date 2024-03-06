# Release History

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