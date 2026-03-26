# Release History

## 3.0.0 (2026-03-26)
### Breaking Changes

- Type of `ContainerGroupListResult.Value` has been changed from `[]*ContainerGroup` to `[]*ListResultContainerGroup`
- Type of `ContainerGroupPropertiesProperties.InstanceView` has been changed from `*ContainerGroupPropertiesInstanceView` to `*ContainerGroupPropertiesPropertiesInstanceView`
- Type of `Operation.Origin` has been changed from `*ContainerInstanceOperationsOrigin` to `*OperationsOrigin`
- Enum `ContainerInstanceOperationsOrigin` has been removed
- Struct `ContainerGroupProperties` has been removed
- Struct `ContainerGroupPropertiesInstanceView` has been removed

### Features Added

- New value `ContainerGroupSKUNotSpecified` added to enum type `ContainerGroupSKU`
- New enum type `AzureFileShareAccessTier` with values `AzureFileShareAccessTierCool`, `AzureFileShareAccessTierHot`, `AzureFileShareAccessTierPremium`, `AzureFileShareAccessTierTransactionOptimized`
- New enum type `AzureFileShareAccessType` with values `AzureFileShareAccessTypeExclusive`, `AzureFileShareAccessTypeShared`
- New enum type `ContainerGroupProvisioningState` with values `ContainerGroupProvisioningStateAccepted`, `ContainerGroupProvisioningStateCanceled`, `ContainerGroupProvisioningStateCreating`, `ContainerGroupProvisioningStateDeleting`, `ContainerGroupProvisioningStateFailed`, `ContainerGroupProvisioningStateNotAccessible`, `ContainerGroupProvisioningStateNotSpecified`, `ContainerGroupProvisioningStatePending`, `ContainerGroupProvisioningStatePreProvisioned`, `ContainerGroupProvisioningStateRepairing`, `ContainerGroupProvisioningStateSucceeded`, `ContainerGroupProvisioningStateUnhealthy`, `ContainerGroupProvisioningStateUpdating`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `IdentityAccessLevel` with values `IdentityAccessLevelAll`, `IdentityAccessLevelSystem`, `IdentityAccessLevelUser`
- New enum type `NGroupProvisioningState` with values `NGroupProvisioningStateCanceled`, `NGroupProvisioningStateCreating`, `NGroupProvisioningStateDeleting`, `NGroupProvisioningStateFailed`, `NGroupProvisioningStateMigrating`, `NGroupProvisioningStateSucceeded`, `NGroupProvisioningStateUpdating`
- New enum type `NGroupUpdateMode` with values `NGroupUpdateModeManual`, `NGroupUpdateModeRolling`
- New enum type `OperationsOrigin` with values `OperationsOriginSystem`, `OperationsOriginUser`
- New function `NewCGProfileClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CGProfileClient, error)`
- New function `*CGProfileClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, containerGroupProfileName string, containerGroupProfile ContainerGroupProfile, options *CGProfileClientCreateOrUpdateOptions) (CGProfileClientCreateOrUpdateResponse, error)`
- New function `*CGProfileClient.Delete(ctx context.Context, resourceGroupName string, containerGroupProfileName string, options *CGProfileClientDeleteOptions) (CGProfileClientDeleteResponse, error)`
- New function `*CGProfileClient.Get(ctx context.Context, resourceGroupName string, containerGroupProfileName string, options *CGProfileClientGetOptions) (CGProfileClientGetResponse, error)`
- New function `*CGProfileClient.GetByRevisionNumber(ctx context.Context, resourceGroupName string, containerGroupProfileName string, revisionNumber string, options *CGProfileClientGetByRevisionNumberOptions) (CGProfileClientGetByRevisionNumberResponse, error)`
- New function `*CGProfileClient.NewListAllRevisionsPager(resourceGroupName string, containerGroupProfileName string, options *CGProfileClientListAllRevisionsOptions) *runtime.Pager[CGProfileClientListAllRevisionsResponse]`
- New function `*CGProfileClient.Update(ctx context.Context, resourceGroupName string, containerGroupProfileName string, properties ContainerGroupProfilePatch, options *CGProfileClientUpdateOptions) (CGProfileClientUpdateResponse, error)`
- New function `NewCGProfilesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CGProfilesClient, error)`
- New function `*CGProfilesClient.NewListByResourceGroupPager(resourceGroupName string, options *CGProfilesClientListByResourceGroupOptions) *runtime.Pager[CGProfilesClientListByResourceGroupResponse]`
- New function `*CGProfilesClient.NewListBySubscriptionPager(options *CGProfilesClientListBySubscriptionOptions) *runtime.Pager[CGProfilesClientListBySubscriptionResponse]`
- New function `*ClientFactory.NewCGProfileClient() *CGProfileClient`
- New function `*ClientFactory.NewCGProfilesClient() *CGProfilesClient`
- New function `*ClientFactory.NewNGroupsClient() *NGroupsClient`
- New function `NewNGroupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*NGroupsClient, error)`
- New function `*NGroupsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, ngroupsName string, nGroup NGroup, options *NGroupsClientBeginCreateOrUpdateOptions) (*runtime.Poller[NGroupsClientCreateOrUpdateResponse], error)`
- New function `*NGroupsClient.BeginDelete(ctx context.Context, resourceGroupName string, ngroupsName string, options *NGroupsClientBeginDeleteOptions) (*runtime.Poller[NGroupsClientDeleteResponse], error)`
- New function `*NGroupsClient.Get(ctx context.Context, resourceGroupName string, ngroupsName string, options *NGroupsClientGetOptions) (NGroupsClientGetResponse, error)`
- New function `*NGroupsClient.NewListByResourceGroupPager(resourceGroupName string, options *NGroupsClientListByResourceGroupOptions) *runtime.Pager[NGroupsClientListByResourceGroupResponse]`
- New function `*NGroupsClient.NewListPager(options *NGroupsClientListOptions) *runtime.Pager[NGroupsClientListResponse]`
- New function `*NGroupsClient.BeginRestart(ctx context.Context, resourceGroupName string, ngroupsName string, options *NGroupsClientBeginRestartOptions) (*runtime.Poller[NGroupsClientRestartResponse], error)`
- New function `*NGroupsClient.BeginStart(ctx context.Context, resourceGroupName string, ngroupsName string, options *NGroupsClientBeginStartOptions) (*runtime.Poller[NGroupsClientStartResponse], error)`
- New function `*NGroupsClient.Stop(ctx context.Context, resourceGroupName string, ngroupsName string, options *NGroupsClientStopOptions) (NGroupsClientStopResponse, error)`
- New function `*NGroupsClient.BeginUpdate(ctx context.Context, resourceGroupName string, ngroupsName string, nGroup NGroupPatch, options *NGroupsClientBeginUpdateOptions) (*runtime.Poller[NGroupsClientUpdateResponse], error)`
- New struct `APIEntityReference`
- New struct `ApplicationGateway`
- New struct `ApplicationGatewayBackendAddressPool`
- New struct `ConfigMap`
- New struct `ContainerGroupProfile`
- New struct `ContainerGroupProfileListResult`
- New struct `ContainerGroupProfilePatch`
- New struct `ContainerGroupProfileProperties`
- New struct `ContainerGroupProfileReferenceDefinition`
- New struct `ContainerGroupProfileStub`
- New struct `ContainerGroupPropertiesPropertiesInstanceView`
- New struct `ElasticProfile`
- New struct `ElasticProfileContainerGroupNamingPolicy`
- New struct `ElasticProfileContainerGroupNamingPolicyGUIDNamingPolicy`
- New struct `FileShare`
- New struct `FileShareProperties`
- New struct `IdentityACLs`
- New struct `IdentityAccessControl`
- New struct `ListResultContainerGroup`
- New struct `ListResultContainerGroupPropertiesProperties`
- New struct `LoadBalancer`
- New struct `LoadBalancerBackendAddressPool`
- New struct `NGroup`
- New struct `NGroupCGPropertyContainer`
- New struct `NGroupCGPropertyContainerProperties`
- New struct `NGroupCGPropertyVolume`
- New struct `NGroupContainerGroupProperties`
- New struct `NGroupIdentity`
- New struct `NGroupPatch`
- New struct `NGroupProperties`
- New struct `NGroupsListResult`
- New struct `NetworkProfile`
- New struct `PlacementProfile`
- New struct `SecretReference`
- New struct `StandbyPoolProfileDefinition`
- New struct `StorageProfile`
- New struct `SystemData`
- New struct `UpdateProfile`
- New struct `UpdateProfileRollingUpdateProfile`
- New field `StorageAccountKeyReference` in struct `AzureFileVolume`
- New field `SystemData` in struct `ContainerGroup`
- New field `ContainerGroupProfile`, `IdentityACLs`, `IsCreatedFromStandbyPool`, `SecretReferences`, `StandbyPoolProfile` in struct `ContainerGroupPropertiesProperties`
- New field `ConfigMap` in struct `ContainerProperties`
- New field `SecureValueReference` in struct `EnvironmentVariable`
- New field `PasswordReference` in struct `ImageRegistryCredential`
- New field `NextLink` in struct `UsageListResult`
- New field `SecretReference` in struct `Volume`


## 2.5.0-beta.1 (2024-10-23)
### Features Added

- New function `*ClientFactory.NewContainerGroupProfileClient() *ContainerGroupProfileClient`
- New function `*ClientFactory.NewContainerGroupProfilesClient() *ContainerGroupProfilesClient`
- New function `NewContainerGroupProfileClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ContainerGroupProfileClient, error)`
- New function `*ContainerGroupProfileClient.GetByRevisionNumber(context.Context, string, string, string, *ContainerGroupProfileClientGetByRevisionNumberOptions) (ContainerGroupProfileClientGetByRevisionNumberResponse, error)`
- New function `*ContainerGroupProfileClient.NewListAllRevisionsPager(string, string, *ContainerGroupProfileClientListAllRevisionsOptions) *runtime.Pager[ContainerGroupProfileClientListAllRevisionsResponse]`
- New function `NewContainerGroupProfilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ContainerGroupProfilesClient, error)`
- New function `*ContainerGroupProfilesClient.CreateOrUpdate(context.Context, string, string, ContainerGroupProfile, *ContainerGroupProfilesClientCreateOrUpdateOptions) (ContainerGroupProfilesClientCreateOrUpdateResponse, error)`
- New function `*ContainerGroupProfilesClient.Delete(context.Context, string, string, *ContainerGroupProfilesClientDeleteOptions) (ContainerGroupProfilesClientDeleteResponse, error)`
- New function `*ContainerGroupProfilesClient.Get(context.Context, string, string, *ContainerGroupProfilesClientGetOptions) (ContainerGroupProfilesClientGetResponse, error)`
- New function `*ContainerGroupProfilesClient.NewListByResourceGroupPager(string, *ContainerGroupProfilesClientListByResourceGroupOptions) *runtime.Pager[ContainerGroupProfilesClientListByResourceGroupResponse]`
- New function `*ContainerGroupProfilesClient.NewListPager(*ContainerGroupProfilesClientListOptions) *runtime.Pager[ContainerGroupProfilesClientListResponse]`
- New function `*ContainerGroupProfilesClient.Patch(context.Context, string, string, ContainerGroupProfilePatch, *ContainerGroupProfilesClientPatchOptions) (ContainerGroupProfilesClientPatchResponse, error)`
- New struct `ConfigMap`
- New struct `ContainerGroupProfile`
- New struct `ContainerGroupProfileListResult`
- New struct `ContainerGroupProfilePatch`
- New struct `ContainerGroupProfileProperties`
- New struct `ContainerGroupProfilePropertiesProperties`
- New struct `ContainerGroupProfileReferenceDefinition`
- New struct `StandbyPoolProfileDefinition`
- New field `ContainerGroupProfile`, `IsCreatedFromStandbyPool`, `StandbyPoolProfile` in struct `ContainerGroupPropertiesProperties`
- New field `ConfigMap` in struct `ContainerProperties`


## 2.4.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.3.0 (2023-04-28)
### Features Added

- New value `ContainerGroupSKUConfidential` added to enum type `ContainerGroupSKU`
- New enum type `ContainerGroupPriority` with values `ContainerGroupPriorityRegular`, `ContainerGroupPrioritySpot`
- New struct `ConfidentialComputeProperties`
- New struct `SecurityContextCapabilitiesDefinition`
- New struct `SecurityContextDefinition`
- New field `ConfidentialComputeProperties`, `Priority` in struct `ContainerGroupPropertiesProperties`
- New field `SecurityContext` in struct `ContainerProperties`
- New field `SecurityContext` in struct `InitContainerPropertiesDefinition`


## 2.2.1 (2023-04-18)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 2.2.0 (2023-04-07)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 2.1.0 (2022-11-08)
### Features Added

- New struct `DeploymentExtensionSpec`
- New struct `DeploymentExtensionSpecProperties`
- New field `Extensions` in struct `ContainerGroupPropertiesProperties`
- New field `Identity` in struct `EncryptionProperties`


## 2.0.0 (2022-08-26)
### Breaking Changes

- Type of `ContainerGroup.Properties` has been changed from `*ContainerGroupProperties` to `*ContainerGroupPropertiesProperties`
- Type of `ContainerGroupIdentity.UserAssignedIdentities` has been changed from `map[string]*Components10Wh5UdSchemasContainergroupidentityPropertiesUserassignedidentitiesAdditionalproperties` to `map[string]*UserAssignedIdentities`
- Const `AutoGeneratedDomainNameLabelScopeResourceGroupReuse` has been removed
- Const `AutoGeneratedDomainNameLabelScopeSubscriptionReuse` has been removed
- Const `AutoGeneratedDomainNameLabelScopeNoreuse` has been removed
- Const `AutoGeneratedDomainNameLabelScopeTenantReuse` has been removed
- Const `AutoGeneratedDomainNameLabelScopeUnsecure` has been removed
- Type alias `AutoGeneratedDomainNameLabelScope` has been removed
- Function `PossibleAutoGeneratedDomainNameLabelScopeValues` has been removed
- Struct `Components10Wh5UdSchemasContainergroupidentityPropertiesUserassignedidentitiesAdditionalproperties` has been removed
- Field `DNSNameLabelReusePolicy` of struct `IPAddress` has been removed
- Field `Containers` of struct `ContainerGroupProperties` has been removed
- Field `SubnetIDs` of struct `ContainerGroupProperties` has been removed
- Field `RestartPolicy` of struct `ContainerGroupProperties` has been removed
- Field `Volumes` of struct `ContainerGroupProperties` has been removed
- Field `ProvisioningState` of struct `ContainerGroupProperties` has been removed
- Field `Diagnostics` of struct `ContainerGroupProperties` has been removed
- Field `SKU` of struct `ContainerGroupProperties` has been removed
- Field `InstanceView` of struct `ContainerGroupProperties` has been removed
- Field `OSType` of struct `ContainerGroupProperties` has been removed
- Field `EncryptionProperties` of struct `ContainerGroupProperties` has been removed
- Field `InitContainers` of struct `ContainerGroupProperties` has been removed
- Field `DNSConfig` of struct `ContainerGroupProperties` has been removed
- Field `IPAddress` of struct `ContainerGroupProperties` has been removed
- Field `ImageRegistryCredentials` of struct `ContainerGroupProperties` has been removed

### Features Added

- New const `DNSNameLabelReusePolicyNoreuse`
- New const `DNSNameLabelReusePolicyUnsecure`
- New const `DNSNameLabelReusePolicySubscriptionReuse`
- New const `DNSNameLabelReusePolicyTenantReuse`
- New const `DNSNameLabelReusePolicyResourceGroupReuse`
- New type alias `DNSNameLabelReusePolicy`
- New function `*SubnetServiceAssociationLinkClient.BeginDelete(context.Context, string, string, string, *SubnetServiceAssociationLinkClientBeginDeleteOptions) (*runtime.Poller[SubnetServiceAssociationLinkClientDeleteResponse], error)`
- New function `NewSubnetServiceAssociationLinkClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SubnetServiceAssociationLinkClient, error)`
- New function `PossibleDNSNameLabelReusePolicyValues() []DNSNameLabelReusePolicy`
- New struct `ContainerGroupPropertiesProperties`
- New struct `SubnetServiceAssociationLinkClient`
- New struct `SubnetServiceAssociationLinkClientBeginDeleteOptions`
- New struct `SubnetServiceAssociationLinkClientDeleteResponse`
- New struct `UserAssignedIdentities`
- New field `Properties` in struct `ContainerGroupProperties`
- New field `Identity` in struct `ContainerGroupProperties`
- New field `ID` in struct `Usage`
- New field `AutoGeneratedDomainNameLabelScope` in struct `IPAddress`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).