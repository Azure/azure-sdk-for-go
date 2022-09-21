# Release History

## 2.0.0 (2022-09-21)
### Breaking Changes

- Function `*TablesClient.Update` has been removed
- Struct `TablesClientUpdateOptions` has been removed

### Features Added

- New const `CreatedByTypeManagedIdentity`
- New const `ColumnDataTypeHintEnumIP`
- New const `ColumnTypeEnumDynamic`
- New const `ColumnDataTypeHintEnumArmPath`
- New const `IdentityTypeApplication`
- New const `ColumnDataTypeHintEnumGUID`
- New const `TablePlanEnumBasic`
- New const `ColumnTypeEnumReal`
- New const `TableTypeEnumSearchResults`
- New const `CreatedByTypeKey`
- New const `RetentionInDaysAsDefaultFalse`
- New const `TableSubTypeEnumClassic`
- New const `ColumnTypeEnumDateTime`
- New const `RetentionInDaysAsDefaultTrue`
- New const `ColumnTypeEnumGUID`
- New const `TablePlanEnumAnalytics`
- New const `ColumnTypeEnumInt`
- New const `TableTypeEnumCustomLog`
- New const `SourceEnumCustomer`
- New const `TableTypeEnumRestoredLogs`
- New const `IdentityTypeUser`
- New const `TableSubTypeEnumDataCollectionRuleBased`
- New const `TableSubTypeEnumAny`
- New const `ColumnTypeEnumString`
- New const `ProvisioningStateEnumSucceeded`
- New const `CreatedByTypeUser`
- New const `DataSourceTypeIngestion`
- New const `TotalRetentionInDaysAsDefaultFalse`
- New const `ProvisioningStateEnumUpdating`
- New const `IdentityTypeKey`
- New const `SourceEnumMicrosoft`
- New const `TotalRetentionInDaysAsDefaultTrue`
- New const `IdentityTypeManagedIdentity`
- New const `TableTypeEnumMicrosoft`
- New const `ColumnTypeEnumLong`
- New const `ProvisioningStateEnumInProgress`
- New const `CreatedByTypeApplication`
- New const `ColumnTypeEnumBoolean`
- New const `ColumnDataTypeHintEnumURI`
- New type alias `TableSubTypeEnum`
- New type alias `ProvisioningStateEnum`
- New type alias `RetentionInDaysAsDefault`
- New type alias `TableTypeEnum`
- New type alias `TotalRetentionInDaysAsDefault`
- New type alias `ColumnDataTypeHintEnum`
- New type alias `SourceEnum`
- New type alias `TablePlanEnum`
- New type alias `CreatedByType`
- New type alias `ColumnTypeEnum`
- New function `*TablesClient.BeginUpdate(context.Context, string, string, string, Table, *TablesClientBeginUpdateOptions) (*runtime.Poller[TablesClientUpdateResponse], error)`
- New function `*TablesClient.Migrate(context.Context, string, string, string, *TablesClientMigrateOptions) (TablesClientMigrateResponse, error)`
- New function `*QueryPacksClient.UpdateTags(context.Context, string, string, TagsResource, *QueryPacksClientUpdateTagsOptions) (QueryPacksClientUpdateTagsResponse, error)`
- New function `NewQueriesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*QueriesClient, error)`
- New function `*QueriesClient.NewSearchPager(string, string, LogAnalyticsQueryPackQuerySearchProperties, *QueriesClientSearchOptions) *runtime.Pager[QueriesClientSearchResponse]`
- New function `PossibleProvisioningStateEnumValues() []ProvisioningStateEnum`
- New function `*QueriesClient.NewListPager(string, string, *QueriesClientListOptions) *runtime.Pager[QueriesClientListResponse]`
- New function `NewQueryPacksClient(string, azcore.TokenCredential, *arm.ClientOptions) (*QueryPacksClient, error)`
- New function `*QueriesClient.Delete(context.Context, string, string, string, *QueriesClientDeleteOptions) (QueriesClientDeleteResponse, error)`
- New function `*QueriesClient.Put(context.Context, string, string, string, LogAnalyticsQueryPackQuery, *QueriesClientPutOptions) (QueriesClientPutResponse, error)`
- New function `*QueryPacksClient.CreateOrUpdate(context.Context, string, string, LogAnalyticsQueryPack, *QueryPacksClientCreateOrUpdateOptions) (QueryPacksClientCreateOrUpdateResponse, error)`
- New function `*QueriesClient.Update(context.Context, string, string, string, LogAnalyticsQueryPackQuery, *QueriesClientUpdateOptions) (QueriesClientUpdateResponse, error)`
- New function `*QueryPacksClient.NewListPager(*QueryPacksClientListOptions) *runtime.Pager[QueryPacksClientListResponse]`
- New function `*TablesClient.BeginDelete(context.Context, string, string, string, *TablesClientBeginDeleteOptions) (*runtime.Poller[TablesClientDeleteResponse], error)`
- New function `*QueriesClient.Get(context.Context, string, string, string, *QueriesClientGetOptions) (QueriesClientGetResponse, error)`
- New function `*QueryPacksClient.NewListByResourceGroupPager(string, *QueryPacksClientListByResourceGroupOptions) *runtime.Pager[QueryPacksClientListByResourceGroupResponse]`
- New function `PossibleSourceEnumValues() []SourceEnum`
- New function `PossibleTotalRetentionInDaysAsDefaultValues() []TotalRetentionInDaysAsDefault`
- New function `PossibleTablePlanEnumValues() []TablePlanEnum`
- New function `*TablesClient.BeginCreateOrUpdate(context.Context, string, string, string, Table, *TablesClientBeginCreateOrUpdateOptions) (*runtime.Poller[TablesClientCreateOrUpdateResponse], error)`
- New function `*QueryPacksClient.Delete(context.Context, string, string, *QueryPacksClientDeleteOptions) (QueryPacksClientDeleteResponse, error)`
- New function `PossibleColumnTypeEnumValues() []ColumnTypeEnum`
- New function `*QueryPacksClient.Get(context.Context, string, string, *QueryPacksClientGetOptions) (QueryPacksClientGetResponse, error)`
- New function `PossibleColumnDataTypeHintEnumValues() []ColumnDataTypeHintEnum`
- New function `PossibleTableTypeEnumValues() []TableTypeEnum`
- New function `PossibleCreatedByTypeValues() []CreatedByType`
- New function `PossibleRetentionInDaysAsDefaultValues() []RetentionInDaysAsDefault`
- New function `*TablesClient.CancelSearch(context.Context, string, string, string, *TablesClientCancelSearchOptions) (TablesClientCancelSearchResponse, error)`
- New function `PossibleTableSubTypeEnumValues() []TableSubTypeEnum`
- New function `*QueryPacksClient.CreateOrUpdateWithoutName(context.Context, string, LogAnalyticsQueryPack, *QueryPacksClientCreateOrUpdateWithoutNameOptions) (QueryPacksClientCreateOrUpdateWithoutNameResponse, error)`
- New struct `AzureResourceProperties`
- New struct `Column`
- New struct `LogAnalyticsQueryPack`
- New struct `LogAnalyticsQueryPackListResult`
- New struct `LogAnalyticsQueryPackProperties`
- New struct `LogAnalyticsQueryPackQuery`
- New struct `LogAnalyticsQueryPackQueryListResult`
- New struct `LogAnalyticsQueryPackQueryProperties`
- New struct `LogAnalyticsQueryPackQueryPropertiesRelated`
- New struct `LogAnalyticsQueryPackQuerySearchProperties`
- New struct `LogAnalyticsQueryPackQuerySearchPropertiesRelated`
- New struct `QueriesClient`
- New struct `QueriesClientDeleteOptions`
- New struct `QueriesClientDeleteResponse`
- New struct `QueriesClientGetOptions`
- New struct `QueriesClientGetResponse`
- New struct `QueriesClientListOptions`
- New struct `QueriesClientListResponse`
- New struct `QueriesClientPutOptions`
- New struct `QueriesClientPutResponse`
- New struct `QueriesClientSearchOptions`
- New struct `QueriesClientSearchResponse`
- New struct `QueriesClientUpdateOptions`
- New struct `QueriesClientUpdateResponse`
- New struct `QueryPacksClient`
- New struct `QueryPacksClientCreateOrUpdateOptions`
- New struct `QueryPacksClientCreateOrUpdateResponse`
- New struct `QueryPacksClientCreateOrUpdateWithoutNameOptions`
- New struct `QueryPacksClientCreateOrUpdateWithoutNameResponse`
- New struct `QueryPacksClientDeleteOptions`
- New struct `QueryPacksClientDeleteResponse`
- New struct `QueryPacksClientGetOptions`
- New struct `QueryPacksClientGetResponse`
- New struct `QueryPacksClientListByResourceGroupOptions`
- New struct `QueryPacksClientListByResourceGroupResponse`
- New struct `QueryPacksClientListOptions`
- New struct `QueryPacksClientListResponse`
- New struct `QueryPacksClientUpdateTagsOptions`
- New struct `QueryPacksClientUpdateTagsResponse`
- New struct `QueryPacksResource`
- New struct `RestoredLogs`
- New struct `ResultStatistics`
- New struct `Schema`
- New struct `SearchResults`
- New struct `SystemData`
- New struct `SystemDataAutoGenerated`
- New struct `TablesClientBeginCreateOrUpdateOptions`
- New struct `TablesClientBeginDeleteOptions`
- New struct `TablesClientBeginUpdateOptions`
- New struct `TablesClientCancelSearchOptions`
- New struct `TablesClientCancelSearchResponse`
- New struct `TablesClientCreateOrUpdateResponse`
- New struct `TablesClientDeleteResponse`
- New struct `TablesClientMigrateOptions`
- New struct `TablesClientMigrateResponse`
- New struct `TagsResource`
- New field `Identity` in struct `WorkspacePatch`
- New field `DefaultDataCollectionRuleResourceID` in struct `WorkspaceProperties`
- New field `SystemData` in struct `Table`
- New field `Identity` in struct `Workspace`
- New field `SystemData` in struct `Workspace`
- New field `ArchiveRetentionInDays` in struct `TableProperties`
- New field `RetentionInDaysAsDefault` in struct `TableProperties`
- New field `TotalRetentionInDays` in struct `TableProperties`
- New field `ResultStatistics` in struct `TableProperties`
- New field `RestoredLogs` in struct `TableProperties`
- New field `LastPlanModifiedDate` in struct `TableProperties`
- New field `ProvisioningState` in struct `TableProperties`
- New field `TotalRetentionInDaysAsDefault` in struct `TableProperties`
- New field `Plan` in struct `TableProperties`
- New field `Schema` in struct `TableProperties`
- New field `SearchResults` in struct `TableProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).