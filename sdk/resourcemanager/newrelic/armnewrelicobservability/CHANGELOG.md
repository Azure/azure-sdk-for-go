# Release History

## 1.2.0 (2024-03-22)
### Features Added

- New enum type `ConfigurationName` with values `ConfigurationNameDefault`
- New enum type `PatchOperation` with values `PatchOperationActive`, `PatchOperationAddBegin`, `PatchOperationAddComplete`, `PatchOperationDeleteBegin`, `PatchOperationDeleteComplete`
- New enum type `Status` with values `StatusActive`, `StatusDeleting`, `StatusFailed`, `StatusInProgress`
- New function `NewBillingInfoClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BillingInfoClient, error)`
- New function `*BillingInfoClient.Get(context.Context, string, string, *BillingInfoClientGetOptions) (BillingInfoClientGetResponse, error)`
- New function `*ClientFactory.NewBillingInfoClient() *BillingInfoClient`
- New function `*ClientFactory.NewConnectedPartnerResourcesClient() *ConnectedPartnerResourcesClient`
- New function `*ClientFactory.NewMonitoredSubscriptionsClient() *MonitoredSubscriptionsClient`
- New function `NewConnectedPartnerResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConnectedPartnerResourcesClient, error)`
- New function `*ConnectedPartnerResourcesClient.NewListPager(string, string, *ConnectedPartnerResourcesClientListOptions) *runtime.Pager[ConnectedPartnerResourcesClientListResponse]`
- New function `NewMonitoredSubscriptionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MonitoredSubscriptionsClient, error)`
- New function `*MonitoredSubscriptionsClient.BeginCreateorUpdate(context.Context, string, string, ConfigurationName, MonitoredSubscriptionProperties, *MonitoredSubscriptionsClientBeginCreateorUpdateOptions) (*runtime.Poller[MonitoredSubscriptionsClientCreateorUpdateResponse], error)`
- New function `*MonitoredSubscriptionsClient.BeginDelete(context.Context, string, string, ConfigurationName, *MonitoredSubscriptionsClientBeginDeleteOptions) (*runtime.Poller[MonitoredSubscriptionsClientDeleteResponse], error)`
- New function `*MonitoredSubscriptionsClient.Get(context.Context, string, string, ConfigurationName, *MonitoredSubscriptionsClientGetOptions) (MonitoredSubscriptionsClientGetResponse, error)`
- New function `*MonitoredSubscriptionsClient.NewListPager(string, string, *MonitoredSubscriptionsClientListOptions) *runtime.Pager[MonitoredSubscriptionsClientListResponse]`
- New function `*MonitoredSubscriptionsClient.BeginUpdate(context.Context, string, string, ConfigurationName, MonitoredSubscriptionProperties, *MonitoredSubscriptionsClientBeginUpdateOptions) (*runtime.Poller[MonitoredSubscriptionsClientUpdateResponse], error)`
- New function `*MonitorsClient.NewListLinkedResourcesPager(string, string, *MonitorsClientListLinkedResourcesOptions) *runtime.Pager[MonitorsClientListLinkedResourcesResponse]`
- New struct `BillingInfoResponse`
- New struct `ConnectedPartnerResourceProperties`
- New struct `ConnectedPartnerResourcesListFormat`
- New struct `ConnectedPartnerResourcesListResponse`
- New struct `LinkedResource`
- New struct `LinkedResourceListResponse`
- New struct `MarketplaceSaaSInfo`
- New struct `MonitoredSubscription`
- New struct `MonitoredSubscriptionProperties`
- New struct `MonitoredSubscriptionPropertiesList`
- New struct `PartnerBillingEntity`
- New struct `SubscriptionList`
- New field `SaaSAzureSubscriptionStatus`, `SubscriptionState` in struct `MonitorProperties`


## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.0.0 (2023-05-26)
### Features Added

- New anonymous field `NewRelicMonitorResource` in struct `MonitorsClientSwitchBillingResponse`
- New field `RetryAfter` in struct `MonitorsClientSwitchBillingResponse`


## 0.1.0 (2023-03-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/newrelic/armnewrelicobservability` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).