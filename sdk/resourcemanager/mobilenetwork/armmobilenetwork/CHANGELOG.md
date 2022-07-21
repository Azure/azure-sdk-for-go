# Release History

## 0.6.0 (2022-07-21)
### Breaking Changes

- Function `*SimsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *SimsClientGetOptions)` to `(context.Context, string, string, string, *SimsClientGetOptions)`
- Function `*SimsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *SimsClientBeginDeleteOptions)` to `(context.Context, string, string, string, *SimsClientBeginDeleteOptions)`
- Function `*SimsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, Sim, *SimsClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, string, Sim, *SimsClientBeginCreateOrUpdateOptions)`
- Function `*SimsClient.NewListByResourceGroupPager` has been removed
- Function `*SimsClient.NewListBySubscriptionPager` has been removed
- Function `*SimsClient.UpdateTags` has been removed
- Struct `SimsClientListByResourceGroupOptions` has been removed
- Struct `SimsClientListByResourceGroupResponse` has been removed
- Struct `SimsClientListBySubscriptionOptions` has been removed
- Struct `SimsClientListBySubscriptionResponse` has been removed
- Struct `SimsClientUpdateTagsOptions` has been removed
- Struct `SimsClientUpdateTagsResponse` has been removed
- Field `MobileNetwork` of struct `SimPropertiesFormat` has been removed
- Field `Tags` of struct `Sim` has been removed
- Field `Location` of struct `Sim` has been removed
- Field `CustomLocation` of struct `PacketCoreControlPlanePropertiesFormat` has been removed

### Features Added

- New const `ManagedServiceIdentityTypeSystemAssigned`
- New const `VersionStateValidationFailed`
- New const `VersionStateUnknown`
- New const `ManagedServiceIdentityTypeNone`
- New const `RecommendedVersionNotRecommended`
- New const `BillingSKUEdgeSite2GBPS`
- New const `BillingSKULargePackage`
- New const `VersionStateActive`
- New const `PlatformTypeAKSHCI`
- New const `BillingSKUFlagshipStarterPackage`
- New const `VersionStateDeprecated`
- New const `BillingSKUEdgeSite4GBPS`
- New const `BillingSKUEdgeSite3GBPS`
- New const `PlatformTypeBaseVM`
- New const `BillingSKUEvaluationPackage`
- New const `ManagedServiceIdentityTypeSystemAssignedUserAssigned`
- New const `BillingSKUMediumPackage`
- New const `ManagedServiceIdentityTypeUserAssigned`
- New const `RecommendedVersionRecommended`
- New const `VersionStatePreview`
- New const `VersionStateValidating`
- New function `PossibleRecommendedVersionValues() []RecommendedVersion`
- New function `PossibleBillingSKUValues() []BillingSKU`
- New function `*SimGroupsClient.UpdateTags(context.Context, string, string, TagsObject, *SimGroupsClientUpdateTagsOptions) (SimGroupsClientUpdateTagsResponse, error)`
- New function `*PacketCoreControlPlaneVersionsClient.NewListByResourceGroupPager(*PacketCoreControlPlaneVersionsClientListByResourceGroupOptions) *runtime.Pager[PacketCoreControlPlaneVersionsClientListByResourceGroupResponse]`
- New function `*SimGroupsClient.NewListByResourceGroupPager(string, *SimGroupsClientListByResourceGroupOptions) *runtime.Pager[SimGroupsClientListByResourceGroupResponse]`
- New function `*PacketCoreControlPlaneVersionsClient.Get(context.Context, string, *PacketCoreControlPlaneVersionsClientGetOptions) (PacketCoreControlPlaneVersionsClientGetResponse, error)`
- New function `*SimGroupsClient.Get(context.Context, string, string, *SimGroupsClientGetOptions) (SimGroupsClientGetResponse, error)`
- New function `PossibleManagedServiceIdentityTypeValues() []ManagedServiceIdentityType`
- New function `*SimGroupsClient.NewListBySubscriptionPager(*SimGroupsClientListBySubscriptionOptions) *runtime.Pager[SimGroupsClientListBySubscriptionResponse]`
- New function `*SimGroupsClient.BeginDelete(context.Context, string, string, *SimGroupsClientBeginDeleteOptions) (*runtime.Poller[SimGroupsClientDeleteResponse], error)`
- New function `NewPacketCoreControlPlaneVersionsClient(azcore.TokenCredential, *arm.ClientOptions) (*PacketCoreControlPlaneVersionsClient, error)`
- New function `PossibleVersionStateValues() []VersionState`
- New function `*SimsClient.NewListBySimGroupPager(string, string, *SimsClientListBySimGroupOptions) *runtime.Pager[SimsClientListBySimGroupResponse]`
- New function `NewSimGroupsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SimGroupsClient, error)`
- New function `PossiblePlatformTypeValues() []PlatformType`
- New function `*SimGroupsClient.BeginCreateOrUpdate(context.Context, string, string, SimGroup, *SimGroupsClientBeginCreateOrUpdateOptions) (*runtime.Poller[SimGroupsClientCreateOrUpdateResponse], error)`
- New struct `AzureStackEdgeDeviceResourceID`
- New struct `ConnectedClusterResourceID`
- New struct `KeyVaultCertificate`
- New struct `KeyVaultKey`
- New struct `LocalDiagnosticsAccessConfiguration`
- New struct `ManagedServiceIdentity`
- New struct `PacketCoreControlPlaneVersion`
- New struct `PacketCoreControlPlaneVersionListResult`
- New struct `PacketCoreControlPlaneVersionPropertiesFormat`
- New struct `PacketCoreControlPlaneVersionsClient`
- New struct `PacketCoreControlPlaneVersionsClientGetOptions`
- New struct `PacketCoreControlPlaneVersionsClientGetResponse`
- New struct `PacketCoreControlPlaneVersionsClientListByResourceGroupOptions`
- New struct `PacketCoreControlPlaneVersionsClientListByResourceGroupResponse`
- New struct `PlatformConfiguration`
- New struct `ProxyResource`
- New struct `SimGroup`
- New struct `SimGroupListResult`
- New struct `SimGroupPropertiesFormat`
- New struct `SimGroupResourceID`
- New struct `SimGroupsClient`
- New struct `SimGroupsClientBeginCreateOrUpdateOptions`
- New struct `SimGroupsClientBeginDeleteOptions`
- New struct `SimGroupsClientCreateOrUpdateResponse`
- New struct `SimGroupsClientDeleteResponse`
- New struct `SimGroupsClientGetOptions`
- New struct `SimGroupsClientGetResponse`
- New struct `SimGroupsClientListByResourceGroupOptions`
- New struct `SimGroupsClientListByResourceGroupResponse`
- New struct `SimGroupsClientListBySubscriptionOptions`
- New struct `SimGroupsClientListBySubscriptionResponse`
- New struct `SimGroupsClientUpdateTagsOptions`
- New struct `SimGroupsClientUpdateTagsResponse`
- New struct `SimsClientListBySimGroupOptions`
- New struct `SimsClientListBySimGroupResponse`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `PacketCoreControlPlane`
- New field `InteropSettings` in struct `PacketCoreControlPlanePropertiesFormat`
- New field `SKU` in struct `PacketCoreControlPlanePropertiesFormat`
- New field `LocalDiagnosticsAccess` in struct `PacketCoreControlPlanePropertiesFormat`
- New field `Platform` in struct `PacketCoreControlPlanePropertiesFormat`
- New field `DNSAddresses` in struct `AttachedDataNetworkPropertiesFormat`


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mobilenetwork/armmobilenetwork` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).