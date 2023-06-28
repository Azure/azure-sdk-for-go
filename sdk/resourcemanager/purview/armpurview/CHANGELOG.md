# Release History

## 2.0.0 (2023-06-28)
### Breaking Changes

- Field `Count` of struct `OperationList` has been removed
- Field `Count` of struct `PrivateEndpointConnectionList` has been removed
- Field `Count` of struct `PrivateLinkResourceList` has been removed

### Features Added

- New enum type `AccountProvisioningState` with values `AccountProvisioningStateCanceled`, `AccountProvisioningStateCreating`, `AccountProvisioningStateDeleting`, `AccountProvisioningStateFailed`, `AccountProvisioningStateMoving`, `AccountProvisioningStateSoftDeleted`, `AccountProvisioningStateSoftDeleting`, `AccountProvisioningStateSucceeded`, `AccountProvisioningStateUnknown`, `AccountProvisioningStateUpdating`
- New enum type `CredentialsType` with values `CredentialsTypeNone`, `CredentialsTypeSystemAssigned`, `CredentialsTypeUserAssigned`
- New enum type `EventHubType` with values `EventHubTypeHook`, `EventHubTypeNotification`
- New enum type `EventStreamingState` with values `EventStreamingStateDisabled`, `EventStreamingStateEnabled`
- New enum type `EventStreamingType` with values `EventStreamingTypeAzure`, `EventStreamingTypeManaged`, `EventStreamingTypeNone`
- New enum type `ManagedEventHubState` with values `ManagedEventHubStateDisabled`, `ManagedEventHubStateEnabled`, `ManagedEventHubStateNotSpecified`
- New enum type `ManagedResourcesPublicNetworkAccess` with values `ManagedResourcesPublicNetworkAccessDisabled`, `ManagedResourcesPublicNetworkAccessEnabled`, `ManagedResourcesPublicNetworkAccessNotSpecified`
- New function `*ClientFactory.NewFeaturesClient() *FeaturesClient`
- New function `*ClientFactory.NewKafkaConfigurationsClient() *KafkaConfigurationsClient`
- New function `*ClientFactory.NewUsagesClient() *UsagesClient`
- New function `NewFeaturesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FeaturesClient, error)`
- New function `*FeaturesClient.AccountGet(context.Context, string, string, BatchFeatureRequest, *FeaturesClientAccountGetOptions) (FeaturesClientAccountGetResponse, error)`
- New function `*FeaturesClient.SubscriptionGet(context.Context, string, BatchFeatureRequest, *FeaturesClientSubscriptionGetOptions) (FeaturesClientSubscriptionGetResponse, error)`
- New function `NewKafkaConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*KafkaConfigurationsClient, error)`
- New function `*KafkaConfigurationsClient.CreateOrUpdate(context.Context, string, string, string, KafkaConfiguration, *KafkaConfigurationsClientCreateOrUpdateOptions) (KafkaConfigurationsClientCreateOrUpdateResponse, error)`
- New function `*KafkaConfigurationsClient.Delete(context.Context, string, string, string, *KafkaConfigurationsClientDeleteOptions) (KafkaConfigurationsClientDeleteResponse, error)`
- New function `*KafkaConfigurationsClient.Get(context.Context, string, string, string, *KafkaConfigurationsClientGetOptions) (KafkaConfigurationsClientGetResponse, error)`
- New function `*KafkaConfigurationsClient.NewListByAccountPager(string, string, *KafkaConfigurationsClientListByAccountOptions) *runtime.Pager[KafkaConfigurationsClientListByAccountResponse]`
- New function `NewUsagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*UsagesClient, error)`
- New function `*UsagesClient.Get(context.Context, string, *UsagesClientGetOptions) (UsagesClientGetResponse, error)`
- New struct `AccountPropertiesAccountStatus`
- New struct `AccountStatus`
- New struct `AccountStatusErrorDetails`
- New struct `BatchFeatureRequest`
- New struct `BatchFeatureStatus`
- New struct `Credentials`
- New struct `KafkaConfiguration`
- New struct `KafkaConfigurationList`
- New struct `KafkaConfigurationProperties`
- New struct `ProxyResourceSystemData`
- New struct `QuotaName`
- New struct `Usage`
- New struct `UsageList`
- New struct `UsageName`
- New field `AccountStatus`, `ManagedEventHubState`, `ManagedResourcesPublicNetworkAccess` in struct `AccountProperties`
- New field `SystemData` in struct `PrivateEndpointConnection`
- New field `SystemData` in struct `ProxyResource`


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/purview/armpurview` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).