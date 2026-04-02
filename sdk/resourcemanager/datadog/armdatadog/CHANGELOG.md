# Release History

## 2.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Function `*MarketplaceAgreementsClient.CreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, options *MarketplaceAgreementsClientCreateOrUpdateOptions)` to `(ctx context.Context, body AgreementResource, options *MarketplaceAgreementsClientCreateOrUpdateOptions)`
- Function `*MonitoredSubscriptionsClient.BeginCreateorUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, monitorName string, configurationName string, options *MonitoredSubscriptionsClientBeginCreateorUpdateOptions)` to `(ctx context.Context, resourceGroupName string, monitorName string, configurationName string, body MonitoredSubscriptionProperties, options *MonitoredSubscriptionsClientBeginCreateorUpdateOptions)`
- Function `*MonitoredSubscriptionsClient.BeginUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, monitorName string, configurationName string, options *MonitoredSubscriptionsClientBeginUpdateOptions)` to `(ctx context.Context, resourceGroupName string, monitorName string, configurationName string, body MonitoredSubscriptionProperties, options *MonitoredSubscriptionsClientBeginUpdateOptions)`
- Function `*MonitorsClient.BeginCreate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, monitorName string, options *MonitorsClientBeginCreateOptions)` to `(ctx context.Context, resourceGroupName string, monitorName string, body MonitorResource, options *MonitorsClientBeginCreateOptions)`
- Function `*MonitorsClient.BeginUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, monitorName string, options *MonitorsClientBeginUpdateOptions)` to `(ctx context.Context, resourceGroupName string, monitorName string, body MonitorResourceUpdateParameters, options *MonitorsClientBeginUpdateOptions)`
- Function `*SingleSignOnConfigurationsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, monitorName string, configurationName string, options *SingleSignOnConfigurationsClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, resourceGroupName string, monitorName string, configurationName string, body SingleSignOnResource, options *SingleSignOnConfigurationsClientBeginCreateOrUpdateOptions)`
- Function `*TagRulesClient.CreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, monitorName string, ruleSetName string, options *TagRulesClientCreateOrUpdateOptions)` to `(ctx context.Context, resourceGroupName string, monitorName string, ruleSetName string, body MonitoringTagRules, options *TagRulesClientCreateOrUpdateOptions)`
- Field `Body` of struct `MarketplaceAgreementsClientCreateOrUpdateOptions` has been removed
- Field `Body` of struct `MonitoredSubscriptionsClientBeginCreateorUpdateOptions` has been removed
- Field `Body` of struct `MonitoredSubscriptionsClientBeginUpdateOptions` has been removed
- Field `Body` of struct `MonitorsClientBeginCreateOptions` has been removed
- Field `Body` of struct `MonitorsClientBeginUpdateOptions` has been removed
- Field `Body` of struct `SingleSignOnConfigurationsClientBeginCreateOrUpdateOptions` has been removed
- Field `Body` of struct `TagRulesClientCreateOrUpdateOptions` has been removed

### Features Added

- New enum type `ConnectorAction` with values `ConnectorActionAdd`, `ConnectorActionRemove`
- New function `NewBillingInfoClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*BillingInfoClient, error)`
- New function `*BillingInfoClient.Get(ctx context.Context, resourceGroupName string, monitorName string, options *BillingInfoClientGetOptions) (BillingInfoClientGetResponse, error)`
- New function `*ClientFactory.NewBillingInfoClient() *BillingInfoClient`
- New function `*ClientFactory.NewMonitorResourcesClient() *MonitorResourcesClient`
- New function `*ClientFactory.NewOrganizationsClient() *OrganizationsClient`
- New function `*ClientFactory.NewSaaSOperationGroupClient() *SaaSOperationGroupClient`
- New function `NewMonitorResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MonitorResourcesClient, error)`
- New function `*MonitorResourcesClient.LatestLinkedSaaS(ctx context.Context, resourceGroupName string, monitorName string, options *MonitorResourcesClientLatestLinkedSaaSOptions) (MonitorResourcesClientLatestLinkedSaaSResponse, error)`
- New function `*MonitorResourcesClient.BeginLinkSaaS(ctx context.Context, resourceGroupName string, monitorName string, body SaaSData, options *MonitorResourcesClientBeginLinkSaaSOptions) (*runtime.Poller[MonitorResourcesClientLinkSaaSResponse], error)`
- New function `*MonitorsClient.GetDefaultApplicationKey(ctx context.Context, resourceGroupName string, monitorName string, options *MonitorsClientGetDefaultApplicationKeyOptions) (MonitorsClientGetDefaultApplicationKeyResponse, error)`
- New function `*MonitorsClient.ManageSreAgentConnectors(ctx context.Context, resourceGroupName string, monitorName string, request SreAgentConnectorRequest, options *MonitorsClientManageSreAgentConnectorsOptions) (MonitorsClientManageSreAgentConnectorsResponse, error)`
- New function `NewOrganizationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OrganizationsClient, error)`
- New function `*OrganizationsClient.BeginResubscribe(ctx context.Context, resourceGroupName string, monitorName string, options *OrganizationsClientBeginResubscribeOptions) (*runtime.Poller[OrganizationsClientResubscribeResponse], error)`
- New function `NewSaaSOperationGroupClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SaaSOperationGroupClient, error)`
- New function `*SaaSOperationGroupClient.ActivateResource(ctx context.Context, body ActivateSaaSParameterRequest, options *SaaSOperationGroupClientActivateResourceOptions) (SaaSOperationGroupClientActivateResourceResponse, error)`
- New struct `ActivateSaaSParameterRequest`
- New struct `AgentRules`
- New struct `ApplicationKey`
- New struct `BillingInfoResponse`
- New struct `LatestLinkedSaaSResponse`
- New struct `MarketplaceOfferDetails`
- New struct `MarketplaceSaaSInfo`
- New struct `PartnerBillingEntity`
- New struct `ResubscribeProperties`
- New struct `SaaSData`
- New struct `SaaSResourceDetailsResponse`
- New struct `SreAgentConfiguration`
- New struct `SreAgentConfigurationListResponse`
- New struct `SreAgentConnectorRequest`
- New field `NextLink` in struct `CreateResourceSupportedResponseList`
- New field `Location` in struct `LinkedResource`
- New field `MarketplaceOfferDetails`, `SaaSData`, `SreAgentConfiguration` in struct `MonitorProperties`
- New field `ResourceCollection` in struct `MonitorUpdateProperties`
- New field `SystemData` in struct `MonitoredSubscriptionProperties`
- New field `NextLink` in struct `MonitoredSubscriptionPropertiesList`
- New field `AgentRules`, `CustomMetrics` in struct `MonitoringTagRulesProperties`
- New field `ResourceCollection` in struct `OrganizationProperties`


## 1.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-10-27)
### Features Added

- New enum type `Operation` with values `OperationActive`, `OperationAddBegin`, `OperationAddComplete`, `OperationDeleteBegin`, `OperationDeleteComplete`
- New enum type `Status` with values `StatusActive`, `StatusDeleting`, `StatusFailed`, `StatusInProgress`
- New function `*ClientFactory.NewCreationSupportedClient() *CreationSupportedClient`
- New function `*ClientFactory.NewMonitoredSubscriptionsClient() *MonitoredSubscriptionsClient`
- New function `NewCreationSupportedClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CreationSupportedClient, error)`
- New function `*CreationSupportedClient.Get(context.Context, string, *CreationSupportedClientGetOptions) (CreationSupportedClientGetResponse, error)`
- New function `*CreationSupportedClient.NewListPager(string, *CreationSupportedClientListOptions) *runtime.Pager[CreationSupportedClientListResponse]`
- New function `NewMonitoredSubscriptionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MonitoredSubscriptionsClient, error)`
- New function `*MonitoredSubscriptionsClient.BeginCreateorUpdate(context.Context, string, string, string, *MonitoredSubscriptionsClientBeginCreateorUpdateOptions) (*runtime.Poller[MonitoredSubscriptionsClientCreateorUpdateResponse], error)`
- New function `*MonitoredSubscriptionsClient.BeginDelete(context.Context, string, string, string, *MonitoredSubscriptionsClientBeginDeleteOptions) (*runtime.Poller[MonitoredSubscriptionsClientDeleteResponse], error)`
- New function `*MonitoredSubscriptionsClient.Get(context.Context, string, string, string, *MonitoredSubscriptionsClientGetOptions) (MonitoredSubscriptionsClientGetResponse, error)`
- New function `*MonitoredSubscriptionsClient.NewListPager(string, string, *MonitoredSubscriptionsClientListOptions) *runtime.Pager[MonitoredSubscriptionsClientListResponse]`
- New function `*MonitoredSubscriptionsClient.BeginUpdate(context.Context, string, string, string, *MonitoredSubscriptionsClientBeginUpdateOptions) (*runtime.Poller[MonitoredSubscriptionsClientUpdateResponse], error)`
- New struct `CreateResourceSupportedProperties`
- New struct `CreateResourceSupportedResponse`
- New struct `CreateResourceSupportedResponseList`
- New struct `MonitoredSubscription`
- New struct `MonitoredSubscriptionProperties`
- New struct `MonitoredSubscriptionPropertiesList`
- New struct `SubscriptionList`
- New field `Cspm` in struct `MonitorUpdateProperties`
- New field `Automuting` in struct `MonitoringTagRulesProperties`
- New field `Cspm` in struct `OrganizationProperties`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datadog/armdatadog` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).