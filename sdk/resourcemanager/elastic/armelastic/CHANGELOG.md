# Release History

## 2.0.0 (2025-11-20)
### Breaking Changes

- Operation `*MonitorsClient.Update` has been changed to LRO, use `*MonitorsClient.BeginUpdate` instead.

### Features Added

- New enum type `ConfigurationType` with values `ConfigurationTypeGeneralPurpose`, `ConfigurationTypeNotApplicable`, `ConfigurationTypeTimeSeries`, `ConfigurationTypeVector`
- New enum type `HostingType` with values `HostingTypeHosted`, `HostingTypeServerless`
- New enum type `Operation` with values `OperationActive`, `OperationAddBegin`, `OperationAddComplete`, `OperationDeleteBegin`, `OperationDeleteComplete`
- New enum type `ProjectType` with values `ProjectTypeElasticsearch`, `ProjectTypeNotApplicable`, `ProjectTypeObservability`, `ProjectTypeSecurity`
- New enum type `Status` with values `StatusActive`, `StatusDeleting`, `StatusFailed`, `StatusInProgress`
- New function `*ClientFactory.NewMonitoredSubscriptionsClient() *MonitoredSubscriptionsClient`
- New function `NewMonitoredSubscriptionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MonitoredSubscriptionsClient, error)`
- New function `*MonitoredSubscriptionsClient.BeginCreateorUpdate(context.Context, string, string, string, *MonitoredSubscriptionsClientBeginCreateorUpdateOptions) (*runtime.Poller[MonitoredSubscriptionsClientCreateorUpdateResponse], error)`
- New function `*MonitoredSubscriptionsClient.BeginDelete(context.Context, string, string, string, *MonitoredSubscriptionsClientBeginDeleteOptions) (*runtime.Poller[MonitoredSubscriptionsClientDeleteResponse], error)`
- New function `*MonitoredSubscriptionsClient.Get(context.Context, string, string, string, *MonitoredSubscriptionsClientGetOptions) (MonitoredSubscriptionsClientGetResponse, error)`
- New function `*MonitoredSubscriptionsClient.NewListPager(string, string, *MonitoredSubscriptionsClientListOptions) *runtime.Pager[MonitoredSubscriptionsClientListResponse]`
- New function `*MonitoredSubscriptionsClient.BeginUpdate(context.Context, string, string, string, *MonitoredSubscriptionsClientBeginUpdateOptions) (*runtime.Poller[MonitoredSubscriptionsClientUpdateResponse], error)`
- New function `*OrganizationsClient.BeginResubscribe(context.Context, string, string, *OrganizationsClientBeginResubscribeOptions) (*runtime.Poller[OrganizationsClientResubscribeResponse], error)`
- New struct `MonitoredSubscription`
- New struct `MonitoredSubscriptionProperties`
- New struct `MonitoredSubscriptionPropertiesList`
- New struct `ProjectDetails`
- New struct `ResubscribeProperties`
- New struct `SubscriptionList`
- New field `Type` in struct `ConnectedPartnerResourceProperties`
- New field `ConfigurationType`, `ProjectType` in struct `DeploymentInfoResponse`
- New field `OfferID`, `PublisherID` in struct `MarketplaceSaaSInfoMarketplaceSubscription`
- New field `HostingType`, `ProjectDetails` in struct `MonitorProperties`
- New field `Kind` in struct `MonitorResource`
- New field `OpenAIConnectorID` in struct `OpenAIIntegrationProperties`


## 1.0.0 (2024-10-24)
### Features Added

- New function `NewBillingInfoClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BillingInfoClient, error)`
- New function `*BillingInfoClient.Get(context.Context, string, string, *BillingInfoClientGetOptions) (BillingInfoClientGetResponse, error)`
- New function `*ClientFactory.NewBillingInfoClient() *BillingInfoClient`
- New function `*ClientFactory.NewConnectedPartnerResourcesClient() *ConnectedPartnerResourcesClient`
- New function `*ClientFactory.NewOpenAIClient() *OpenAIClient`
- New function `NewConnectedPartnerResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConnectedPartnerResourcesClient, error)`
- New function `*ConnectedPartnerResourcesClient.NewListPager(string, string, *ConnectedPartnerResourcesClientListOptions) *runtime.Pager[ConnectedPartnerResourcesClientListResponse]`
- New function `NewOpenAIClient(string, azcore.TokenCredential, *arm.ClientOptions) (*OpenAIClient, error)`
- New function `*OpenAIClient.CreateOrUpdate(context.Context, string, string, string, *OpenAIClientCreateOrUpdateOptions) (OpenAIClientCreateOrUpdateResponse, error)`
- New function `*OpenAIClient.Delete(context.Context, string, string, string, *OpenAIClientDeleteOptions) (OpenAIClientDeleteResponse, error)`
- New function `*OpenAIClient.Get(context.Context, string, string, string, *OpenAIClientGetOptions) (OpenAIClientGetResponse, error)`
- New function `*OpenAIClient.GetStatus(context.Context, string, string, string, *OpenAIClientGetStatusOptions) (OpenAIClientGetStatusResponse, error)`
- New function `*OpenAIClient.NewListPager(string, string, *OpenAIClientListOptions) *runtime.Pager[OpenAIClientListResponse]`
- New function `*OrganizationsClient.GetElasticToAzureSubscriptionMapping(context.Context, *OrganizationsClientGetElasticToAzureSubscriptionMappingOptions) (OrganizationsClientGetElasticToAzureSubscriptionMappingResponse, error)`
- New struct `BillingInfoResponse`
- New struct `ConnectedPartnerResourceProperties`
- New struct `ConnectedPartnerResourcesListFormat`
- New struct `ConnectedPartnerResourcesListResponse`
- New struct `OpenAIIntegrationProperties`
- New struct `OpenAIIntegrationRPModel`
- New struct `OpenAIIntegrationRPModelListResponse`
- New struct `OpenAIIntegrationStatusResponse`
- New struct `OpenAIIntegrationStatusResponseProperties`
- New struct `OrganizationToAzureSubscriptionMappingResponse`
- New struct `OrganizationToAzureSubscriptionMappingResponseProperties`
- New struct `PartnerBillingEntity`
- New struct `PlanDetails`
- New field `ElasticsearchEndPoint` in struct `DeploymentInfoResponse`
- New field `BilledAzureSubscriptionID`, `MarketplaceStatus`, `Subscribed` in struct `MarketplaceSaaSInfo`
- New field `PlanDetails`, `SaaSAzureSubscriptionStatus`, `SourceCampaignID`, `SourceCampaignName`, `SubscriptionState` in struct `MonitorProperties`


## 0.10.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.9.0 (2023-05-26)
### Breaking Changes

- Function `*OrganizationsClient.GetAPIKey` parameter(s) have been changed from `(context.Context, string, *OrganizationsClientGetAPIKeyOptions)` to `(context.Context, *OrganizationsClientGetAPIKeyOptions)`
- Field `GenerateAPIKey` of struct `MonitorResource` has been removed
- Field `APIKey` of struct `UserAPIKeyResponse` has been removed

### Features Added

- New function `*ClientFactory.NewVersionsClient() *VersionsClient`
- New function `NewVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VersionsClient, error)`
- New function `*VersionsClient.NewListPager(string, *VersionsClientListOptions) *runtime.Pager[VersionsClientListResponse]`
- New struct `UserAPIKeyResponseProperties`
- New struct `VersionListFormat`
- New struct `VersionListProperties`
- New struct `VersionsListResponse`
- New field `GenerateAPIKey` in struct `MonitorProperties`
- New field `Properties` in struct `UserAPIKeyResponse`


## 0.8.0 (2023-04-28)
### Features Added

- New function `*ClientFactory.NewOrganizationsClient() *OrganizationsClient`
- New function `NewOrganizationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*OrganizationsClient, error)`
- New function `*OrganizationsClient.GetAPIKey(context.Context, string, *OrganizationsClientGetAPIKeyOptions) (OrganizationsClientGetAPIKeyResponse, error)`
- New struct `MarketplaceSaaSInfo`
- New struct `MarketplaceSaaSInfoMarketplaceSubscription`
- New struct `UserAPIKeyResponse`
- New struct `UserEmailID`
- New field `DeploymentURL` in struct `DeploymentInfoResponse`
- New field `MarketplaceSaasInfo` in struct `DeploymentInfoResponse`
- New field `GenerateAPIKey` in struct `MonitorResource`


## 0.7.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.6.0 (2022-11-10)
### Features Added

- New const `TypeAzurePrivateEndpoint`
- New const `TypeIP`
- New type alias `Type`
- New function `*TrafficFiltersClient.Delete(context.Context, string, string, *TrafficFiltersClientDeleteOptions) (TrafficFiltersClientDeleteResponse, error)`
- New function `*AssociateTrafficFilterClient.BeginAssociate(context.Context, string, string, *AssociateTrafficFilterClientBeginAssociateOptions) (*runtime.Poller[AssociateTrafficFilterClientAssociateResponse], error)`
- New function `NewCreateAndAssociateIPFilterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CreateAndAssociateIPFilterClient, error)`
- New function `NewUpgradableVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*UpgradableVersionsClient, error)`
- New function `NewTrafficFiltersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TrafficFiltersClient, error)`
- New function `NewDetachTrafficFilterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DetachTrafficFilterClient, error)`
- New function `NewAllTrafficFiltersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AllTrafficFiltersClient, error)`
- New function `NewDetachAndDeleteTrafficFilterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DetachAndDeleteTrafficFilterClient, error)`
- New function `*MonitorClient.BeginUpgrade(context.Context, string, string, *MonitorClientBeginUpgradeOptions) (*runtime.Poller[MonitorClientUpgradeResponse], error)`
- New function `*AllTrafficFiltersClient.List(context.Context, string, string, *AllTrafficFiltersClientListOptions) (AllTrafficFiltersClientListResponse, error)`
- New function `*ExternalUserClient.CreateOrUpdate(context.Context, string, string, *ExternalUserClientCreateOrUpdateOptions) (ExternalUserClientCreateOrUpdateResponse, error)`
- New function `NewListAssociatedTrafficFiltersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ListAssociatedTrafficFiltersClient, error)`
- New function `*CreateAndAssociateIPFilterClient.BeginCreate(context.Context, string, string, *CreateAndAssociateIPFilterClientBeginCreateOptions) (*runtime.Poller[CreateAndAssociateIPFilterClientCreateResponse], error)`
- New function `*DetachAndDeleteTrafficFilterClient.Delete(context.Context, string, string, *DetachAndDeleteTrafficFilterClientDeleteOptions) (DetachAndDeleteTrafficFilterClientDeleteResponse, error)`
- New function `NewAssociateTrafficFilterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AssociateTrafficFilterClient, error)`
- New function `*DetachTrafficFilterClient.BeginUpdate(context.Context, string, string, *DetachTrafficFilterClientBeginUpdateOptions) (*runtime.Poller[DetachTrafficFilterClientUpdateResponse], error)`
- New function `*CreateAndAssociatePLFilterClient.BeginCreate(context.Context, string, string, *CreateAndAssociatePLFilterClientBeginCreateOptions) (*runtime.Poller[CreateAndAssociatePLFilterClientCreateResponse], error)`
- New function `*ListAssociatedTrafficFiltersClient.List(context.Context, string, string, *ListAssociatedTrafficFiltersClientListOptions) (ListAssociatedTrafficFiltersClientListResponse, error)`
- New function `PossibleTypeValues() []Type`
- New function `NewCreateAndAssociatePLFilterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CreateAndAssociatePLFilterClient, error)`
- New function `NewMonitorClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MonitorClient, error)`
- New function `NewExternalUserClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ExternalUserClient, error)`
- New function `*UpgradableVersionsClient.Details(context.Context, string, string, *UpgradableVersionsClientDetailsOptions) (UpgradableVersionsClientDetailsResponse, error)`
- New struct `AllTrafficFiltersClient`
- New struct `AllTrafficFiltersClientListOptions`
- New struct `AllTrafficFiltersClientListResponse`
- New struct `AssociateTrafficFilterClient`
- New struct `AssociateTrafficFilterClientAssociateResponse`
- New struct `AssociateTrafficFilterClientBeginAssociateOptions`
- New struct `CreateAndAssociateIPFilterClient`
- New struct `CreateAndAssociateIPFilterClientBeginCreateOptions`
- New struct `CreateAndAssociateIPFilterClientCreateResponse`
- New struct `CreateAndAssociatePLFilterClient`
- New struct `CreateAndAssociatePLFilterClientBeginCreateOptions`
- New struct `CreateAndAssociatePLFilterClientCreateResponse`
- New struct `DetachAndDeleteTrafficFilterClient`
- New struct `DetachAndDeleteTrafficFilterClientDeleteOptions`
- New struct `DetachAndDeleteTrafficFilterClientDeleteResponse`
- New struct `DetachTrafficFilterClient`
- New struct `DetachTrafficFilterClientBeginUpdateOptions`
- New struct `DetachTrafficFilterClientUpdateResponse`
- New struct `ExternalUserClient`
- New struct `ExternalUserClientCreateOrUpdateOptions`
- New struct `ExternalUserClientCreateOrUpdateResponse`
- New struct `ExternalUserCreationResponse`
- New struct `ExternalUserInfo`
- New struct `ListAssociatedTrafficFiltersClient`
- New struct `ListAssociatedTrafficFiltersClientListOptions`
- New struct `ListAssociatedTrafficFiltersClientListResponse`
- New struct `MonitorClient`
- New struct `MonitorClientBeginUpgradeOptions`
- New struct `MonitorClientUpgradeResponse`
- New struct `MonitorUpgrade`
- New struct `TrafficFilter`
- New struct `TrafficFilterResponse`
- New struct `TrafficFilterRule`
- New struct `TrafficFiltersClient`
- New struct `TrafficFiltersClientDeleteOptions`
- New struct `TrafficFiltersClientDeleteResponse`
- New struct `UpgradableVersionsClient`
- New struct `UpgradableVersionsClientDetailsOptions`
- New struct `UpgradableVersionsClientDetailsResponse`
- New struct `UpgradableVersionsList`
- New field `Version` in struct `MonitorProperties`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/elastic/armelastic` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).