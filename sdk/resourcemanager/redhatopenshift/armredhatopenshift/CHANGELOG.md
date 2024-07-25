# Release History

## 1.6.0 (2024-07-26)
### Features Added

- New value `ProvisioningStateCanceled` added to enum type `ProvisioningState`
- New struct `EffectiveOutboundIP`
- New struct `LoadBalancerProfile`
- New struct `ManagedOutboundIPs`
- New field `LoadBalancerProfile` in struct `NetworkProfile`


## 1.5.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.4.0 (2023-10-27)
### Features Added

- New enum type `PreconfiguredNSG` with values `PreconfiguredNSGDisabled`, `PreconfiguredNSGEnabled`
- New field `PreconfiguredNSG` in struct `NetworkProfile`
- New field `WorkerProfilesStatus` in struct `OpenShiftClusterProperties`


## 1.3.0 (2023-08-25)
### Features Added

- New enum type `OutboundType` with values `OutboundTypeLoadbalancer`, `OutboundTypeUserDefinedRouting`
- New field `OutboundType` in struct `NetworkProfile`


## 1.2.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.1.0 (2023-01-27)
### Features Added

- New function `NewMachinePoolsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MachinePoolsClient, error)`
- New function `*MachinePoolsClient.CreateOrUpdate(context.Context, string, string, string, MachinePool, *MachinePoolsClientCreateOrUpdateOptions) (MachinePoolsClientCreateOrUpdateResponse, error)`
- New function `*MachinePoolsClient.Delete(context.Context, string, string, string, *MachinePoolsClientDeleteOptions) (MachinePoolsClientDeleteResponse, error)`
- New function `*MachinePoolsClient.Get(context.Context, string, string, string, *MachinePoolsClientGetOptions) (MachinePoolsClientGetResponse, error)`
- New function `*MachinePoolsClient.NewListPager(string, string, *MachinePoolsClientListOptions) *runtime.Pager[MachinePoolsClientListResponse]`
- New function `*MachinePoolsClient.Update(context.Context, string, string, string, MachinePoolUpdate, *MachinePoolsClientUpdateOptions) (MachinePoolsClientUpdateResponse, error)`
- New function `NewOpenShiftVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*OpenShiftVersionsClient, error)`
- New function `*OpenShiftVersionsClient.NewListPager(string, *OpenShiftVersionsClientListOptions) *runtime.Pager[OpenShiftVersionsClientListResponse]`
- New function `NewSecretsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SecretsClient, error)`
- New function `*SecretsClient.CreateOrUpdate(context.Context, string, string, string, Secret, *SecretsClientCreateOrUpdateOptions) (SecretsClientCreateOrUpdateResponse, error)`
- New function `*SecretsClient.Delete(context.Context, string, string, string, *SecretsClientDeleteOptions) (SecretsClientDeleteResponse, error)`
- New function `*SecretsClient.Get(context.Context, string, string, string, *SecretsClientGetOptions) (SecretsClientGetResponse, error)`
- New function `*SecretsClient.NewListPager(string, string, *SecretsClientListOptions) *runtime.Pager[SecretsClientListResponse]`
- New function `*SecretsClient.Update(context.Context, string, string, string, SecretUpdate, *SecretsClientUpdateOptions) (SecretsClientUpdateResponse, error)`
- New function `NewSyncIdentityProvidersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SyncIdentityProvidersClient, error)`
- New function `*SyncIdentityProvidersClient.CreateOrUpdate(context.Context, string, string, string, SyncIdentityProvider, *SyncIdentityProvidersClientCreateOrUpdateOptions) (SyncIdentityProvidersClientCreateOrUpdateResponse, error)`
- New function `*SyncIdentityProvidersClient.Delete(context.Context, string, string, string, *SyncIdentityProvidersClientDeleteOptions) (SyncIdentityProvidersClientDeleteResponse, error)`
- New function `*SyncIdentityProvidersClient.Get(context.Context, string, string, string, *SyncIdentityProvidersClientGetOptions) (SyncIdentityProvidersClientGetResponse, error)`
- New function `*SyncIdentityProvidersClient.NewListPager(string, string, *SyncIdentityProvidersClientListOptions) *runtime.Pager[SyncIdentityProvidersClientListResponse]`
- New function `*SyncIdentityProvidersClient.Update(context.Context, string, string, string, SyncIdentityProviderUpdate, *SyncIdentityProvidersClientUpdateOptions) (SyncIdentityProvidersClientUpdateResponse, error)`
- New function `NewSyncSetsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SyncSetsClient, error)`
- New function `*SyncSetsClient.CreateOrUpdate(context.Context, string, string, string, SyncSet, *SyncSetsClientCreateOrUpdateOptions) (SyncSetsClientCreateOrUpdateResponse, error)`
- New function `*SyncSetsClient.Delete(context.Context, string, string, string, *SyncSetsClientDeleteOptions) (SyncSetsClientDeleteResponse, error)`
- New function `*SyncSetsClient.Get(context.Context, string, string, string, *SyncSetsClientGetOptions) (SyncSetsClientGetResponse, error)`
- New function `*SyncSetsClient.NewListPager(string, string, *SyncSetsClientListOptions) *runtime.Pager[SyncSetsClientListResponse]`
- New function `*SyncSetsClient.Update(context.Context, string, string, string, SyncSetUpdate, *SyncSetsClientUpdateOptions) (SyncSetsClientUpdateResponse, error)`
- New struct `MachinePool`
- New struct `MachinePoolList`
- New struct `MachinePoolProperties`
- New struct `MachinePoolUpdate`
- New struct `MachinePoolsClient`
- New struct `MachinePoolsClientListResponse`
- New struct `OpenShiftVersion`
- New struct `OpenShiftVersionList`
- New struct `OpenShiftVersionProperties`
- New struct `OpenShiftVersionsClient`
- New struct `OpenShiftVersionsClientListResponse`
- New struct `ProxyResource`
- New struct `Secret`
- New struct `SecretList`
- New struct `SecretProperties`
- New struct `SecretUpdate`
- New struct `SecretsClient`
- New struct `SecretsClientListResponse`
- New struct `SyncIdentityProvider`
- New struct `SyncIdentityProviderList`
- New struct `SyncIdentityProviderProperties`
- New struct `SyncIdentityProviderUpdate`
- New struct `SyncIdentityProvidersClient`
- New struct `SyncIdentityProvidersClientListResponse`
- New struct `SyncSet`
- New struct `SyncSetList`
- New struct `SyncSetProperties`
- New struct `SyncSetUpdate`
- New struct `SyncSetsClient`
- New struct `SyncSetsClientListResponse`
- New field `SystemData` in struct `Resource`
- New field `SystemData` in struct `TrackedResource`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redhatopenshift/armredhatopenshift` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).