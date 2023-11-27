# Release History

## 1.3.0 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.3.0-beta.1 (2023-10-27)
### Features Added

- New enum type `EventSubTypeValues` with values `EventSubTypeValuesRetirement`
- New field `MaintenanceEndTime`, `MaintenanceStartTime`, `ResourceGroup`, `ResourceName`, `Status` in struct `EventImpactedResourceProperties`
- New field `ArgQuery`, `EventSubType`, `MaintenanceID`, `MaintenanceType` in struct `EventProperties`


## 1.2.0 (2023-05-26)
### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `EventLevelValues` with values `EventLevelValuesCritical`, `EventLevelValuesError`, `EventLevelValuesInformational`, `EventLevelValuesWarning`
- New enum type `EventSourceValues` with values `EventSourceValuesResourceHealth`, `EventSourceValuesServiceHealth`
- New enum type `EventStatusValues` with values `EventStatusValuesActive`, `EventStatusValuesResolved`
- New enum type `EventTypeValues` with values `EventTypeValuesEmergingIssues`, `EventTypeValuesHealthAdvisory`, `EventTypeValuesPlannedMaintenance`, `EventTypeValuesRCA`, `EventTypeValuesSecurityAdvisory`, `EventTypeValuesServiceIssue`
- New enum type `IssueNameParameter` with values `IssueNameParameterDefault`
- New enum type `LevelValues` with values `LevelValuesCritical`, `LevelValuesWarning`
- New enum type `LinkTypeValues` with values `LinkTypeValuesButton`, `LinkTypeValuesHyperlink`
- New enum type `Scenario` with values `ScenarioAlerts`
- New enum type `SeverityValues` with values `SeverityValuesError`, `SeverityValuesInformation`, `SeverityValuesWarning`
- New enum type `StageValues` with values `StageValuesActive`, `StageValuesArchived`, `StageValuesResolve`
- New function `NewChildAvailabilityStatusesClient(azcore.TokenCredential, *arm.ClientOptions) (*ChildAvailabilityStatusesClient, error)`
- New function `*ChildAvailabilityStatusesClient.GetByResource(context.Context, string, *ChildAvailabilityStatusesClientGetByResourceOptions) (ChildAvailabilityStatusesClientGetByResourceResponse, error)`
- New function `*ChildAvailabilityStatusesClient.NewListPager(string, *ChildAvailabilityStatusesClientListOptions) *runtime.Pager[ChildAvailabilityStatusesClientListResponse]`
- New function `NewChildResourcesClient(azcore.TokenCredential, *arm.ClientOptions) (*ChildResourcesClient, error)`
- New function `*ChildResourcesClient.NewListPager(string, *ChildResourcesClientListOptions) *runtime.Pager[ChildResourcesClientListResponse]`
- New function `*ClientFactory.NewChildAvailabilityStatusesClient() *ChildAvailabilityStatusesClient`
- New function `*ClientFactory.NewChildResourcesClient() *ChildResourcesClient`
- New function `*ClientFactory.NewEmergingIssuesClient() *EmergingIssuesClient`
- New function `*ClientFactory.NewEventClient() *EventClient`
- New function `*ClientFactory.NewEventsClient() *EventsClient`
- New function `*ClientFactory.NewImpactedResourcesClient() *ImpactedResourcesClient`
- New function `*ClientFactory.NewMetadataClient() *MetadataClient`
- New function `*ClientFactory.NewSecurityAdvisoryImpactedResourcesClient() *SecurityAdvisoryImpactedResourcesClient`
- New function `NewEmergingIssuesClient(azcore.TokenCredential, *arm.ClientOptions) (*EmergingIssuesClient, error)`
- New function `*EmergingIssuesClient.Get(context.Context, IssueNameParameter, *EmergingIssuesClientGetOptions) (EmergingIssuesClientGetResponse, error)`
- New function `*EmergingIssuesClient.NewListPager(*EmergingIssuesClientListOptions) *runtime.Pager[EmergingIssuesClientListResponse]`
- New function `NewEventClient(string, azcore.TokenCredential, *arm.ClientOptions) (*EventClient, error)`
- New function `*EventClient.FetchDetailsBySubscriptionIDAndTrackingID(context.Context, string, *EventClientFetchDetailsBySubscriptionIDAndTrackingIDOptions) (EventClientFetchDetailsBySubscriptionIDAndTrackingIDResponse, error)`
- New function `*EventClient.FetchDetailsByTenantIDAndTrackingID(context.Context, string, *EventClientFetchDetailsByTenantIDAndTrackingIDOptions) (EventClientFetchDetailsByTenantIDAndTrackingIDResponse, error)`
- New function `*EventClient.GetBySubscriptionIDAndTrackingID(context.Context, string, *EventClientGetBySubscriptionIDAndTrackingIDOptions) (EventClientGetBySubscriptionIDAndTrackingIDResponse, error)`
- New function `*EventClient.GetByTenantIDAndTrackingID(context.Context, string, *EventClientGetByTenantIDAndTrackingIDOptions) (EventClientGetByTenantIDAndTrackingIDResponse, error)`
- New function `NewEventsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*EventsClient, error)`
- New function `*EventsClient.NewListBySingleResourcePager(string, *EventsClientListBySingleResourceOptions) *runtime.Pager[EventsClientListBySingleResourceResponse]`
- New function `*EventsClient.NewListBySubscriptionIDPager(*EventsClientListBySubscriptionIDOptions) *runtime.Pager[EventsClientListBySubscriptionIDResponse]`
- New function `*EventsClient.NewListByTenantIDPager(*EventsClientListByTenantIDOptions) *runtime.Pager[EventsClientListByTenantIDResponse]`
- New function `NewImpactedResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ImpactedResourcesClient, error)`
- New function `*ImpactedResourcesClient.Get(context.Context, string, string, *ImpactedResourcesClientGetOptions) (ImpactedResourcesClientGetResponse, error)`
- New function `*ImpactedResourcesClient.GetByTenantID(context.Context, string, string, *ImpactedResourcesClientGetByTenantIDOptions) (ImpactedResourcesClientGetByTenantIDResponse, error)`
- New function `*ImpactedResourcesClient.NewListBySubscriptionIDAndEventIDPager(string, *ImpactedResourcesClientListBySubscriptionIDAndEventIDOptions) *runtime.Pager[ImpactedResourcesClientListBySubscriptionIDAndEventIDResponse]`
- New function `*ImpactedResourcesClient.NewListByTenantIDAndEventIDPager(string, *ImpactedResourcesClientListByTenantIDAndEventIDOptions) *runtime.Pager[ImpactedResourcesClientListByTenantIDAndEventIDResponse]`
- New function `NewMetadataClient(azcore.TokenCredential, *arm.ClientOptions) (*MetadataClient, error)`
- New function `*MetadataClient.GetEntity(context.Context, string, *MetadataClientGetEntityOptions) (MetadataClientGetEntityResponse, error)`
- New function `*MetadataClient.NewListPager(*MetadataClientListOptions) *runtime.Pager[MetadataClientListResponse]`
- New function `NewSecurityAdvisoryImpactedResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SecurityAdvisoryImpactedResourcesClient, error)`
- New function `*SecurityAdvisoryImpactedResourcesClient.NewListBySubscriptionIDAndEventIDPager(string, *SecurityAdvisoryImpactedResourcesClientListBySubscriptionIDAndEventIDOptions) *runtime.Pager[SecurityAdvisoryImpactedResourcesClientListBySubscriptionIDAndEventIDResponse]`
- New function `*SecurityAdvisoryImpactedResourcesClient.NewListByTenantIDAndEventIDPager(string, *SecurityAdvisoryImpactedResourcesClientListByTenantIDAndEventIDOptions) *runtime.Pager[SecurityAdvisoryImpactedResourcesClientListByTenantIDAndEventIDResponse]`
- New struct `EmergingIssue`
- New struct `EmergingIssueImpact`
- New struct `EmergingIssueListResult`
- New struct `EmergingIssuesGetResult`
- New struct `Event`
- New struct `EventImpactedResource`
- New struct `EventImpactedResourceListResult`
- New struct `EventImpactedResourceProperties`
- New struct `EventProperties`
- New struct `EventPropertiesAdditionalInformation`
- New struct `EventPropertiesArticle`
- New struct `EventPropertiesRecommendedActions`
- New struct `EventPropertiesRecommendedActionsItem`
- New struct `Events`
- New struct `Faq`
- New struct `Impact`
- New struct `ImpactedServiceRegion`
- New struct `KeyValueItem`
- New struct `Link`
- New struct `LinkDisplayText`
- New struct `MetadataEntity`
- New struct `MetadataEntityListResult`
- New struct `MetadataEntityProperties`
- New struct `MetadataSupportedValueDetail`
- New struct `StatusActiveEvent`
- New struct `SystemData`
- New struct `Update`
- New field `ArticleID`, `Category`, `Context` in struct `AvailabilityStatusProperties`
- New field `ActionURLComment` in struct `RecommendedAction`


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcehealth/armresourcehealth` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).