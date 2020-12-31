Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

### Removed Funcs

1. *EntityQuery.UnmarshalJSON([]byte) error

## Struct Changes

### Removed Structs

1. EntityQueryProperties

### Removed Struct Fields

1. EntityQuery.*EntityQueryProperties

## Signature Changes

### Funcs

1. EntityQueriesClient.Get
	- Returns
		- From: EntityQuery, error
		- To: EntityQueryModel, error
1. EntityQueriesClient.GetResponder
	- Returns
		- From: EntityQuery, error
		- To: EntityQueryModel, error
1. EntityQueryListIterator.Value
	- Returns
		- From: EntityQuery
		- To: BasicEntityQuery
1. EntityQueryListPage.Values
	- Returns
		- From: []EntityQuery
		- To: []BasicEntityQuery

### Struct Fields

1. EntityQueryList.Value changed type from *[]EntityQuery to *[]BasicEntityQuery

### New Constants

1. DataConnectorKind.DataConnectorKindDynamics365
1. EntityQueryKind.Expansion
1. EntityQueryKind.Insight
1. GroupingEntityType.FileHash
1. KindBasicDataConnector.KindDynamics365
1. KindBasicDataConnectorsCheckRequirements.KindBasicDataConnectorsCheckRequirementsKindDynamics365
1. KindBasicEntityQuery.KindEntityQuery
1. KindBasicEntityQuery.KindExpansion
1. KindBasicEntityQueryItem.KindEntityQueryItem
1. KindBasicEntityQueryItem.KindInsight
1. OutputType.OutputTypeDate
1. OutputType.OutputTypeEntity
1. OutputType.OutputTypeNumber
1. OutputType.OutputTypeString

### New Funcs

1. *Dynamics365CheckRequirements.UnmarshalJSON([]byte) error
1. *Dynamics365DataConnector.UnmarshalJSON([]byte) error
1. *EntityQueryList.UnmarshalJSON([]byte) error
1. *EntityQueryModel.UnmarshalJSON([]byte) error
1. *ExpansionEntityQuery.UnmarshalJSON([]byte) error
1. *GetQueriesResponse.UnmarshalJSON([]byte) error
1. AADCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. AADDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. AATPCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. AATPDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. ASCCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. ASCDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. AwsCloudTrailCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. AwsCloudTrailDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. DataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. DataConnectorsCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. Dynamics365CheckRequirements.AsAADCheckRequirements() (*AADCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsAATPCheckRequirements() (*AATPCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsASCCheckRequirements() (*ASCCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsAwsCloudTrailCheckRequirements() (*AwsCloudTrailCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsBasicDataConnectorsCheckRequirements() (BasicDataConnectorsCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsDataConnectorsCheckRequirements() (*DataConnectorsCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. Dynamics365CheckRequirements.AsMCASCheckRequirements() (*MCASCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsMDATPCheckRequirements() (*MDATPCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsOfficeATPCheckRequirements() (*OfficeATPCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsTICheckRequirements() (*TICheckRequirements, bool)
1. Dynamics365CheckRequirements.AsTiTaxiiCheckRequirements() (*TiTaxiiCheckRequirements, bool)
1. Dynamics365CheckRequirements.MarshalJSON() ([]byte, error)
1. Dynamics365DataConnector.AsAADDataConnector() (*AADDataConnector, bool)
1. Dynamics365DataConnector.AsAATPDataConnector() (*AATPDataConnector, bool)
1. Dynamics365DataConnector.AsASCDataConnector() (*ASCDataConnector, bool)
1. Dynamics365DataConnector.AsAwsCloudTrailDataConnector() (*AwsCloudTrailDataConnector, bool)
1. Dynamics365DataConnector.AsBasicDataConnector() (BasicDataConnector, bool)
1. Dynamics365DataConnector.AsDataConnector() (*DataConnector, bool)
1. Dynamics365DataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. Dynamics365DataConnector.AsMCASDataConnector() (*MCASDataConnector, bool)
1. Dynamics365DataConnector.AsMDATPDataConnector() (*MDATPDataConnector, bool)
1. Dynamics365DataConnector.AsOfficeATPDataConnector() (*OfficeATPDataConnector, bool)
1. Dynamics365DataConnector.AsOfficeDataConnector() (*OfficeDataConnector, bool)
1. Dynamics365DataConnector.AsTIDataConnector() (*TIDataConnector, bool)
1. Dynamics365DataConnector.AsTiTaxiiDataConnector() (*TiTaxiiDataConnector, bool)
1. Dynamics365DataConnector.MarshalJSON() ([]byte, error)
1. EntitiesClient.GetInsights(context.Context, string, string, string, string, EntityGetInsightsParameters) (EntityGetInsightsResponse, error)
1. EntitiesClient.GetInsightsPreparer(context.Context, string, string, string, string, EntityGetInsightsParameters) (*http.Request, error)
1. EntitiesClient.GetInsightsResponder(*http.Response) (EntityGetInsightsResponse, error)
1. EntitiesClient.GetInsightsSender(*http.Request) (*http.Response, error)
1. EntitiesClient.Queries(context.Context, string, string, string, string) (GetQueriesResponse, error)
1. EntitiesClient.QueriesPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. EntitiesClient.QueriesResponder(*http.Response) (GetQueriesResponse, error)
1. EntitiesClient.QueriesSender(*http.Request) (*http.Response, error)
1. EntityQuery.AsBasicEntityQuery() (BasicEntityQuery, bool)
1. EntityQuery.AsEntityQuery() (*EntityQuery, bool)
1. EntityQuery.AsExpansionEntityQuery() (*ExpansionEntityQuery, bool)
1. EntityQueryItem.AsBasicEntityQueryItem() (BasicEntityQueryItem, bool)
1. EntityQueryItem.AsEntityQueryItem() (*EntityQueryItem, bool)
1. EntityQueryItem.AsInsightQueryItem() (*InsightQueryItem, bool)
1. EntityQueryItem.MarshalJSON() ([]byte, error)
1. ExpansionEntityQuery.AsBasicEntityQuery() (BasicEntityQuery, bool)
1. ExpansionEntityQuery.AsEntityQuery() (*EntityQuery, bool)
1. ExpansionEntityQuery.AsExpansionEntityQuery() (*ExpansionEntityQuery, bool)
1. ExpansionEntityQuery.MarshalJSON() ([]byte, error)
1. InsightQueryItem.AsBasicEntityQueryItem() (BasicEntityQueryItem, bool)
1. InsightQueryItem.AsEntityQueryItem() (*EntityQueryItem, bool)
1. InsightQueryItem.AsInsightQueryItem() (*InsightQueryItem, bool)
1. InsightQueryItem.MarshalJSON() ([]byte, error)
1. MCASCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. MCASDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. MDATPCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. MDATPDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. OfficeATPCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. OfficeATPDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. OfficeDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. PossibleEntityQueryKindValues() []EntityQueryKind
1. PossibleKindBasicEntityQueryItemValues() []KindBasicEntityQueryItem
1. PossibleKindBasicEntityQueryValues() []KindBasicEntityQuery
1. PossibleOutputTypeValues() []OutputType
1. TICheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. TIDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. TiTaxiiCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. TiTaxiiDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)

## Struct Changes

### New Structs

1. Dynamics365CheckRequirements
1. Dynamics365CheckRequirementsProperties
1. Dynamics365DataConnector
1. Dynamics365DataConnectorDataTypes
1. Dynamics365DataConnectorDataTypesDynamics365CdsActivities
1. Dynamics365DataConnectorProperties
1. EntityGetInsightsParameters
1. EntityGetInsightsResponse
1. EntityInsightItem
1. EntityInsightItemQueryTimeInterval
1. EntityQueryItem
1. EntityQueryItemProperties
1. EntityQueryItemPropertiesDataTypesItem
1. EntityQueryKind1
1. EntityQueryModel
1. ExpansionEntityQueriesProperties
1. ExpansionEntityQuery
1. GetInsightsError
1. GetInsightsResultsMetadata
1. GetQueriesResponse
1. InsightQueryItem
1. InsightQueryItemProperties
1. InsightQueryItemPropertiesAdditionalQuery
1. InsightQueryItemPropertiesDefaultTimeRange
1. InsightQueryItemPropertiesReferenceTimeRange
1. InsightQueryItemPropertiesTableQuery
1. InsightQueryItemPropertiesTableQueryColumnsDefinitionsItem
1. InsightQueryItemPropertiesTableQueryQueriesDefinitionsItem
1. InsightQueryItemPropertiesTableQueryQueriesDefinitionsItemLinkColumnsDefinitionsItem
1. InsightsTableResult
1. InsightsTableResultColumnsItem

### New Struct Fields

1. AlertRuleTemplatePropertiesBase.LastUpdatedDateUTC
1. EntityQuery.Etag
1. EntityQuery.Kind
1. FusionAlertRuleTemplateProperties.LastUpdatedDateUTC
1. MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties.LastUpdatedDateUTC
1. ScheduledAlertRuleTemplateProperties.LastUpdatedDateUTC
