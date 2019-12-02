// +build go1.9

// Copyright 2019 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package documentdb

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2019-08-01/documentdb"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type CompositePathSortOrder = original.CompositePathSortOrder

const (
	Ascending  CompositePathSortOrder = original.Ascending
	Descending CompositePathSortOrder = original.Descending
)

type ConflictResolutionMode = original.ConflictResolutionMode

const (
	Custom         ConflictResolutionMode = original.Custom
	LastWriterWins ConflictResolutionMode = original.LastWriterWins
)

type ConnectorOffer = original.ConnectorOffer

const (
	Small ConnectorOffer = original.Small
)

type DataType = original.DataType

const (
	LineString   DataType = original.LineString
	MultiPolygon DataType = original.MultiPolygon
	Number       DataType = original.Number
	Point        DataType = original.Point
	Polygon      DataType = original.Polygon
	String       DataType = original.String
)

type DatabaseAccountKind = original.DatabaseAccountKind

const (
	GlobalDocumentDB DatabaseAccountKind = original.GlobalDocumentDB
	MongoDB          DatabaseAccountKind = original.MongoDB
	Parse            DatabaseAccountKind = original.Parse
)

type DatabaseAccountOfferType = original.DatabaseAccountOfferType

const (
	Standard DatabaseAccountOfferType = original.Standard
)

type DefaultConsistencyLevel = original.DefaultConsistencyLevel

const (
	BoundedStaleness DefaultConsistencyLevel = original.BoundedStaleness
	ConsistentPrefix DefaultConsistencyLevel = original.ConsistentPrefix
	Eventual         DefaultConsistencyLevel = original.Eventual
	Session          DefaultConsistencyLevel = original.Session
	Strong           DefaultConsistencyLevel = original.Strong
)

type IndexKind = original.IndexKind

const (
	Hash    IndexKind = original.Hash
	Range   IndexKind = original.Range
	Spatial IndexKind = original.Spatial
)

type IndexingMode = original.IndexingMode

const (
	Consistent IndexingMode = original.Consistent
	Lazy       IndexingMode = original.Lazy
	None       IndexingMode = original.None
)

type KeyKind = original.KeyKind

const (
	Primary           KeyKind = original.Primary
	PrimaryReadonly   KeyKind = original.PrimaryReadonly
	Secondary         KeyKind = original.Secondary
	SecondaryReadonly KeyKind = original.SecondaryReadonly
)

type PartitionKind = original.PartitionKind

const (
	PartitionKindHash  PartitionKind = original.PartitionKindHash
	PartitionKindRange PartitionKind = original.PartitionKindRange
)

type PrimaryAggregationType = original.PrimaryAggregationType

const (
	PrimaryAggregationTypeAverage PrimaryAggregationType = original.PrimaryAggregationTypeAverage
	PrimaryAggregationTypeLast    PrimaryAggregationType = original.PrimaryAggregationTypeLast
	PrimaryAggregationTypeMaximum PrimaryAggregationType = original.PrimaryAggregationTypeMaximum
	PrimaryAggregationTypeMinimum PrimaryAggregationType = original.PrimaryAggregationTypeMinimum
	PrimaryAggregationTypeNone    PrimaryAggregationType = original.PrimaryAggregationTypeNone
	PrimaryAggregationTypeTotal   PrimaryAggregationType = original.PrimaryAggregationTypeTotal
)

type SpatialType = original.SpatialType

const (
	SpatialTypeLineString   SpatialType = original.SpatialTypeLineString
	SpatialTypeMultiPolygon SpatialType = original.SpatialTypeMultiPolygon
	SpatialTypePoint        SpatialType = original.SpatialTypePoint
	SpatialTypePolygon      SpatialType = original.SpatialTypePolygon
)

type TriggerOperation = original.TriggerOperation

const (
	All     TriggerOperation = original.All
	Create  TriggerOperation = original.Create
	Delete  TriggerOperation = original.Delete
	Replace TriggerOperation = original.Replace
	Update  TriggerOperation = original.Update
)

type TriggerType = original.TriggerType

const (
	Post TriggerType = original.Post
	Pre  TriggerType = original.Pre
)

type UnitType = original.UnitType

const (
	Bytes          UnitType = original.Bytes
	BytesPerSecond UnitType = original.BytesPerSecond
	Count          UnitType = original.Count
	CountPerSecond UnitType = original.CountPerSecond
	Milliseconds   UnitType = original.Milliseconds
	Percent        UnitType = original.Percent
	Seconds        UnitType = original.Seconds
)

type ARMProxyResource = original.ARMProxyResource
type ARMResourceProperties = original.ARMResourceProperties
type BaseClient = original.BaseClient
type Capability = original.Capability
type CassandraKeyspaceCreateUpdateParameters = original.CassandraKeyspaceCreateUpdateParameters
type CassandraKeyspaceCreateUpdateProperties = original.CassandraKeyspaceCreateUpdateProperties
type CassandraKeyspaceGetProperties = original.CassandraKeyspaceGetProperties
type CassandraKeyspaceGetPropertiesResource = original.CassandraKeyspaceGetPropertiesResource
type CassandraKeyspaceGetResults = original.CassandraKeyspaceGetResults
type CassandraKeyspaceListResult = original.CassandraKeyspaceListResult
type CassandraKeyspaceResource = original.CassandraKeyspaceResource
type CassandraPartitionKey = original.CassandraPartitionKey
type CassandraResourcesClient = original.CassandraResourcesClient
type CassandraResourcesCreateUpdateCassandraKeyspaceFuture = original.CassandraResourcesCreateUpdateCassandraKeyspaceFuture
type CassandraResourcesCreateUpdateCassandraTableFuture = original.CassandraResourcesCreateUpdateCassandraTableFuture
type CassandraResourcesDeleteCassandraKeyspaceFuture = original.CassandraResourcesDeleteCassandraKeyspaceFuture
type CassandraResourcesDeleteCassandraTableFuture = original.CassandraResourcesDeleteCassandraTableFuture
type CassandraResourcesUpdateCassandraKeyspaceThroughputFuture = original.CassandraResourcesUpdateCassandraKeyspaceThroughputFuture
type CassandraResourcesUpdateCassandraTableThroughputFuture = original.CassandraResourcesUpdateCassandraTableThroughputFuture
type CassandraSchema = original.CassandraSchema
type CassandraTableCreateUpdateParameters = original.CassandraTableCreateUpdateParameters
type CassandraTableCreateUpdateProperties = original.CassandraTableCreateUpdateProperties
type CassandraTableGetProperties = original.CassandraTableGetProperties
type CassandraTableGetPropertiesResource = original.CassandraTableGetPropertiesResource
type CassandraTableGetResults = original.CassandraTableGetResults
type CassandraTableListResult = original.CassandraTableListResult
type CassandraTableResource = original.CassandraTableResource
type ClusterKey = original.ClusterKey
type CollectionClient = original.CollectionClient
type CollectionPartitionClient = original.CollectionPartitionClient
type CollectionPartitionRegionClient = original.CollectionPartitionRegionClient
type CollectionRegionClient = original.CollectionRegionClient
type Column = original.Column
type CompositePath = original.CompositePath
type ConflictResolutionPolicy = original.ConflictResolutionPolicy
type ConsistencyPolicy = original.ConsistencyPolicy
type ContainerPartitionKey = original.ContainerPartitionKey
type DatabaseAccountConnectionString = original.DatabaseAccountConnectionString
type DatabaseAccountCreateUpdateParameters = original.DatabaseAccountCreateUpdateParameters
type DatabaseAccountCreateUpdateProperties = original.DatabaseAccountCreateUpdateProperties
type DatabaseAccountGetProperties = original.DatabaseAccountGetProperties
type DatabaseAccountGetResults = original.DatabaseAccountGetResults
type DatabaseAccountListConnectionStringsResult = original.DatabaseAccountListConnectionStringsResult
type DatabaseAccountListKeysResult = original.DatabaseAccountListKeysResult
type DatabaseAccountListReadOnlyKeysResult = original.DatabaseAccountListReadOnlyKeysResult
type DatabaseAccountRegenerateKeyParameters = original.DatabaseAccountRegenerateKeyParameters
type DatabaseAccountRegionClient = original.DatabaseAccountRegionClient
type DatabaseAccountUpdateParameters = original.DatabaseAccountUpdateParameters
type DatabaseAccountUpdateProperties = original.DatabaseAccountUpdateProperties
type DatabaseAccountsClient = original.DatabaseAccountsClient
type DatabaseAccountsCreateOrUpdateFuture = original.DatabaseAccountsCreateOrUpdateFuture
type DatabaseAccountsDeleteFuture = original.DatabaseAccountsDeleteFuture
type DatabaseAccountsFailoverPriorityChangeFuture = original.DatabaseAccountsFailoverPriorityChangeFuture
type DatabaseAccountsListResult = original.DatabaseAccountsListResult
type DatabaseAccountsOfflineRegionFuture = original.DatabaseAccountsOfflineRegionFuture
type DatabaseAccountsOnlineRegionFuture = original.DatabaseAccountsOnlineRegionFuture
type DatabaseAccountsRegenerateKeyFuture = original.DatabaseAccountsRegenerateKeyFuture
type DatabaseAccountsUpdateFuture = original.DatabaseAccountsUpdateFuture
type DatabaseClient = original.DatabaseClient
type ErrorResponse = original.ErrorResponse
type ExcludedPath = original.ExcludedPath
type ExtendedResourceProperties = original.ExtendedResourceProperties
type FailoverPolicies = original.FailoverPolicies
type FailoverPolicy = original.FailoverPolicy
type GremlinDatabaseCreateUpdateParameters = original.GremlinDatabaseCreateUpdateParameters
type GremlinDatabaseCreateUpdateProperties = original.GremlinDatabaseCreateUpdateProperties
type GremlinDatabaseGetProperties = original.GremlinDatabaseGetProperties
type GremlinDatabaseGetPropertiesResource = original.GremlinDatabaseGetPropertiesResource
type GremlinDatabaseGetResults = original.GremlinDatabaseGetResults
type GremlinDatabaseListResult = original.GremlinDatabaseListResult
type GremlinDatabaseResource = original.GremlinDatabaseResource
type GremlinGraphCreateUpdateParameters = original.GremlinGraphCreateUpdateParameters
type GremlinGraphCreateUpdateProperties = original.GremlinGraphCreateUpdateProperties
type GremlinGraphGetProperties = original.GremlinGraphGetProperties
type GremlinGraphGetPropertiesResource = original.GremlinGraphGetPropertiesResource
type GremlinGraphGetResults = original.GremlinGraphGetResults
type GremlinGraphListResult = original.GremlinGraphListResult
type GremlinGraphResource = original.GremlinGraphResource
type GremlinResourcesClient = original.GremlinResourcesClient
type GremlinResourcesCreateUpdateGremlinDatabaseFuture = original.GremlinResourcesCreateUpdateGremlinDatabaseFuture
type GremlinResourcesCreateUpdateGremlinGraphFuture = original.GremlinResourcesCreateUpdateGremlinGraphFuture
type GremlinResourcesDeleteGremlinDatabaseFuture = original.GremlinResourcesDeleteGremlinDatabaseFuture
type GremlinResourcesDeleteGremlinGraphFuture = original.GremlinResourcesDeleteGremlinGraphFuture
type GremlinResourcesUpdateGremlinDatabaseThroughputFuture = original.GremlinResourcesUpdateGremlinDatabaseThroughputFuture
type GremlinResourcesUpdateGremlinGraphThroughputFuture = original.GremlinResourcesUpdateGremlinGraphThroughputFuture
type IncludedPath = original.IncludedPath
type Indexes = original.Indexes
type IndexingPolicy = original.IndexingPolicy
type Location = original.Location
type Metric = original.Metric
type MetricAvailability = original.MetricAvailability
type MetricDefinition = original.MetricDefinition
type MetricDefinitionsListResult = original.MetricDefinitionsListResult
type MetricListResult = original.MetricListResult
type MetricName = original.MetricName
type MetricValue = original.MetricValue
type MongoDBCollectionCreateUpdateParameters = original.MongoDBCollectionCreateUpdateParameters
type MongoDBCollectionCreateUpdateProperties = original.MongoDBCollectionCreateUpdateProperties
type MongoDBCollectionGetProperties = original.MongoDBCollectionGetProperties
type MongoDBCollectionGetPropertiesResource = original.MongoDBCollectionGetPropertiesResource
type MongoDBCollectionGetResults = original.MongoDBCollectionGetResults
type MongoDBCollectionListResult = original.MongoDBCollectionListResult
type MongoDBCollectionResource = original.MongoDBCollectionResource
type MongoDBDatabaseCreateUpdateParameters = original.MongoDBDatabaseCreateUpdateParameters
type MongoDBDatabaseCreateUpdateProperties = original.MongoDBDatabaseCreateUpdateProperties
type MongoDBDatabaseGetProperties = original.MongoDBDatabaseGetProperties
type MongoDBDatabaseGetPropertiesResource = original.MongoDBDatabaseGetPropertiesResource
type MongoDBDatabaseGetResults = original.MongoDBDatabaseGetResults
type MongoDBDatabaseListResult = original.MongoDBDatabaseListResult
type MongoDBDatabaseResource = original.MongoDBDatabaseResource
type MongoDBResourcesClient = original.MongoDBResourcesClient
type MongoDBResourcesCreateUpdateMongoDBCollectionFuture = original.MongoDBResourcesCreateUpdateMongoDBCollectionFuture
type MongoDBResourcesCreateUpdateMongoDBDatabaseFuture = original.MongoDBResourcesCreateUpdateMongoDBDatabaseFuture
type MongoDBResourcesDeleteMongoDBCollectionFuture = original.MongoDBResourcesDeleteMongoDBCollectionFuture
type MongoDBResourcesDeleteMongoDBDatabaseFuture = original.MongoDBResourcesDeleteMongoDBDatabaseFuture
type MongoDBResourcesUpdateMongoDBCollectionThroughputFuture = original.MongoDBResourcesUpdateMongoDBCollectionThroughputFuture
type MongoDBResourcesUpdateMongoDBDatabaseThroughputFuture = original.MongoDBResourcesUpdateMongoDBDatabaseThroughputFuture
type MongoIndex = original.MongoIndex
type MongoIndexKeys = original.MongoIndexKeys
type MongoIndexOptions = original.MongoIndexOptions
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type OperationsClient = original.OperationsClient
type PartitionKeyRangeIDClient = original.PartitionKeyRangeIDClient
type PartitionKeyRangeIDRegionClient = original.PartitionKeyRangeIDRegionClient
type PartitionMetric = original.PartitionMetric
type PartitionMetricListResult = original.PartitionMetricListResult
type PartitionUsage = original.PartitionUsage
type PartitionUsagesResult = original.PartitionUsagesResult
type PercentileClient = original.PercentileClient
type PercentileMetric = original.PercentileMetric
type PercentileMetricListResult = original.PercentileMetricListResult
type PercentileMetricValue = original.PercentileMetricValue
type PercentileSourceTargetClient = original.PercentileSourceTargetClient
type PercentileTargetClient = original.PercentileTargetClient
type RegionForOnlineOffline = original.RegionForOnlineOffline
type SQLContainerCreateUpdateParameters = original.SQLContainerCreateUpdateParameters
type SQLContainerCreateUpdateProperties = original.SQLContainerCreateUpdateProperties
type SQLContainerGetProperties = original.SQLContainerGetProperties
type SQLContainerGetPropertiesResource = original.SQLContainerGetPropertiesResource
type SQLContainerGetResults = original.SQLContainerGetResults
type SQLContainerListResult = original.SQLContainerListResult
type SQLContainerResource = original.SQLContainerResource
type SQLDatabaseCreateUpdateParameters = original.SQLDatabaseCreateUpdateParameters
type SQLDatabaseCreateUpdateProperties = original.SQLDatabaseCreateUpdateProperties
type SQLDatabaseGetProperties = original.SQLDatabaseGetProperties
type SQLDatabaseGetPropertiesResource = original.SQLDatabaseGetPropertiesResource
type SQLDatabaseGetResults = original.SQLDatabaseGetResults
type SQLDatabaseListResult = original.SQLDatabaseListResult
type SQLDatabaseResource = original.SQLDatabaseResource
type SQLResourcesClient = original.SQLResourcesClient
type SQLResourcesCreateUpdateSQLContainerFuture = original.SQLResourcesCreateUpdateSQLContainerFuture
type SQLResourcesCreateUpdateSQLDatabaseFuture = original.SQLResourcesCreateUpdateSQLDatabaseFuture
type SQLResourcesCreateUpdateSQLStoredProcedureFuture = original.SQLResourcesCreateUpdateSQLStoredProcedureFuture
type SQLResourcesCreateUpdateSQLTriggerFuture = original.SQLResourcesCreateUpdateSQLTriggerFuture
type SQLResourcesCreateUpdateSQLUserDefinedFunctionFuture = original.SQLResourcesCreateUpdateSQLUserDefinedFunctionFuture
type SQLResourcesDeleteSQLContainerFuture = original.SQLResourcesDeleteSQLContainerFuture
type SQLResourcesDeleteSQLDatabaseFuture = original.SQLResourcesDeleteSQLDatabaseFuture
type SQLResourcesDeleteSQLStoredProcedureFuture = original.SQLResourcesDeleteSQLStoredProcedureFuture
type SQLResourcesDeleteSQLTriggerFuture = original.SQLResourcesDeleteSQLTriggerFuture
type SQLResourcesDeleteSQLUserDefinedFunctionFuture = original.SQLResourcesDeleteSQLUserDefinedFunctionFuture
type SQLResourcesUpdateSQLContainerThroughputFuture = original.SQLResourcesUpdateSQLContainerThroughputFuture
type SQLResourcesUpdateSQLDatabaseThroughputFuture = original.SQLResourcesUpdateSQLDatabaseThroughputFuture
type SQLStoredProcedureCreateUpdateParameters = original.SQLStoredProcedureCreateUpdateParameters
type SQLStoredProcedureCreateUpdateProperties = original.SQLStoredProcedureCreateUpdateProperties
type SQLStoredProcedureGetProperties = original.SQLStoredProcedureGetProperties
type SQLStoredProcedureGetPropertiesResource = original.SQLStoredProcedureGetPropertiesResource
type SQLStoredProcedureGetResults = original.SQLStoredProcedureGetResults
type SQLStoredProcedureListResult = original.SQLStoredProcedureListResult
type SQLStoredProcedureResource = original.SQLStoredProcedureResource
type SQLTriggerCreateUpdateParameters = original.SQLTriggerCreateUpdateParameters
type SQLTriggerCreateUpdateProperties = original.SQLTriggerCreateUpdateProperties
type SQLTriggerGetProperties = original.SQLTriggerGetProperties
type SQLTriggerGetPropertiesResource = original.SQLTriggerGetPropertiesResource
type SQLTriggerGetResults = original.SQLTriggerGetResults
type SQLTriggerListResult = original.SQLTriggerListResult
type SQLTriggerResource = original.SQLTriggerResource
type SQLUserDefinedFunctionCreateUpdateParameters = original.SQLUserDefinedFunctionCreateUpdateParameters
type SQLUserDefinedFunctionCreateUpdateProperties = original.SQLUserDefinedFunctionCreateUpdateProperties
type SQLUserDefinedFunctionGetProperties = original.SQLUserDefinedFunctionGetProperties
type SQLUserDefinedFunctionGetPropertiesResource = original.SQLUserDefinedFunctionGetPropertiesResource
type SQLUserDefinedFunctionGetResults = original.SQLUserDefinedFunctionGetResults
type SQLUserDefinedFunctionListResult = original.SQLUserDefinedFunctionListResult
type SQLUserDefinedFunctionResource = original.SQLUserDefinedFunctionResource
type SpatialSpec = original.SpatialSpec
type TableCreateUpdateParameters = original.TableCreateUpdateParameters
type TableCreateUpdateProperties = original.TableCreateUpdateProperties
type TableGetProperties = original.TableGetProperties
type TableGetPropertiesResource = original.TableGetPropertiesResource
type TableGetResults = original.TableGetResults
type TableListResult = original.TableListResult
type TableResource = original.TableResource
type TableResourcesClient = original.TableResourcesClient
type TableResourcesCreateUpdateTableFuture = original.TableResourcesCreateUpdateTableFuture
type TableResourcesDeleteTableFuture = original.TableResourcesDeleteTableFuture
type TableResourcesUpdateTableThroughputFuture = original.TableResourcesUpdateTableThroughputFuture
type ThroughputSettingsGetProperties = original.ThroughputSettingsGetProperties
type ThroughputSettingsGetResults = original.ThroughputSettingsGetResults
type ThroughputSettingsResource = original.ThroughputSettingsResource
type ThroughputSettingsUpdateParameters = original.ThroughputSettingsUpdateParameters
type ThroughputSettingsUpdateProperties = original.ThroughputSettingsUpdateProperties
type UniqueKey = original.UniqueKey
type UniqueKeyPolicy = original.UniqueKeyPolicy
type Usage = original.Usage
type UsagesResult = original.UsagesResult
type VirtualNetworkRule = original.VirtualNetworkRule

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewCassandraResourcesClient(subscriptionID string) CassandraResourcesClient {
	return original.NewCassandraResourcesClient(subscriptionID)
}
func NewCassandraResourcesClientWithBaseURI(baseURI string, subscriptionID string) CassandraResourcesClient {
	return original.NewCassandraResourcesClientWithBaseURI(baseURI, subscriptionID)
}
func NewCollectionClient(subscriptionID string) CollectionClient {
	return original.NewCollectionClient(subscriptionID)
}
func NewCollectionClientWithBaseURI(baseURI string, subscriptionID string) CollectionClient {
	return original.NewCollectionClientWithBaseURI(baseURI, subscriptionID)
}
func NewCollectionPartitionClient(subscriptionID string) CollectionPartitionClient {
	return original.NewCollectionPartitionClient(subscriptionID)
}
func NewCollectionPartitionClientWithBaseURI(baseURI string, subscriptionID string) CollectionPartitionClient {
	return original.NewCollectionPartitionClientWithBaseURI(baseURI, subscriptionID)
}
func NewCollectionPartitionRegionClient(subscriptionID string) CollectionPartitionRegionClient {
	return original.NewCollectionPartitionRegionClient(subscriptionID)
}
func NewCollectionPartitionRegionClientWithBaseURI(baseURI string, subscriptionID string) CollectionPartitionRegionClient {
	return original.NewCollectionPartitionRegionClientWithBaseURI(baseURI, subscriptionID)
}
func NewCollectionRegionClient(subscriptionID string) CollectionRegionClient {
	return original.NewCollectionRegionClient(subscriptionID)
}
func NewCollectionRegionClientWithBaseURI(baseURI string, subscriptionID string) CollectionRegionClient {
	return original.NewCollectionRegionClientWithBaseURI(baseURI, subscriptionID)
}
func NewDatabaseAccountRegionClient(subscriptionID string) DatabaseAccountRegionClient {
	return original.NewDatabaseAccountRegionClient(subscriptionID)
}
func NewDatabaseAccountRegionClientWithBaseURI(baseURI string, subscriptionID string) DatabaseAccountRegionClient {
	return original.NewDatabaseAccountRegionClientWithBaseURI(baseURI, subscriptionID)
}
func NewDatabaseAccountsClient(subscriptionID string) DatabaseAccountsClient {
	return original.NewDatabaseAccountsClient(subscriptionID)
}
func NewDatabaseAccountsClientWithBaseURI(baseURI string, subscriptionID string) DatabaseAccountsClient {
	return original.NewDatabaseAccountsClientWithBaseURI(baseURI, subscriptionID)
}
func NewDatabaseClient(subscriptionID string) DatabaseClient {
	return original.NewDatabaseClient(subscriptionID)
}
func NewDatabaseClientWithBaseURI(baseURI string, subscriptionID string) DatabaseClient {
	return original.NewDatabaseClientWithBaseURI(baseURI, subscriptionID)
}
func NewGremlinResourcesClient(subscriptionID string) GremlinResourcesClient {
	return original.NewGremlinResourcesClient(subscriptionID)
}
func NewGremlinResourcesClientWithBaseURI(baseURI string, subscriptionID string) GremlinResourcesClient {
	return original.NewGremlinResourcesClientWithBaseURI(baseURI, subscriptionID)
}
func NewMongoDBResourcesClient(subscriptionID string) MongoDBResourcesClient {
	return original.NewMongoDBResourcesClient(subscriptionID)
}
func NewMongoDBResourcesClientWithBaseURI(baseURI string, subscriptionID string) MongoDBResourcesClient {
	return original.NewMongoDBResourcesClientWithBaseURI(baseURI, subscriptionID)
}
func NewOperationListResultIterator(page OperationListResultPage) OperationListResultIterator {
	return original.NewOperationListResultIterator(page)
}
func NewOperationListResultPage(getNextPage func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage {
	return original.NewOperationListResultPage(getNextPage)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewPartitionKeyRangeIDClient(subscriptionID string) PartitionKeyRangeIDClient {
	return original.NewPartitionKeyRangeIDClient(subscriptionID)
}
func NewPartitionKeyRangeIDClientWithBaseURI(baseURI string, subscriptionID string) PartitionKeyRangeIDClient {
	return original.NewPartitionKeyRangeIDClientWithBaseURI(baseURI, subscriptionID)
}
func NewPartitionKeyRangeIDRegionClient(subscriptionID string) PartitionKeyRangeIDRegionClient {
	return original.NewPartitionKeyRangeIDRegionClient(subscriptionID)
}
func NewPartitionKeyRangeIDRegionClientWithBaseURI(baseURI string, subscriptionID string) PartitionKeyRangeIDRegionClient {
	return original.NewPartitionKeyRangeIDRegionClientWithBaseURI(baseURI, subscriptionID)
}
func NewPercentileClient(subscriptionID string) PercentileClient {
	return original.NewPercentileClient(subscriptionID)
}
func NewPercentileClientWithBaseURI(baseURI string, subscriptionID string) PercentileClient {
	return original.NewPercentileClientWithBaseURI(baseURI, subscriptionID)
}
func NewPercentileSourceTargetClient(subscriptionID string) PercentileSourceTargetClient {
	return original.NewPercentileSourceTargetClient(subscriptionID)
}
func NewPercentileSourceTargetClientWithBaseURI(baseURI string, subscriptionID string) PercentileSourceTargetClient {
	return original.NewPercentileSourceTargetClientWithBaseURI(baseURI, subscriptionID)
}
func NewPercentileTargetClient(subscriptionID string) PercentileTargetClient {
	return original.NewPercentileTargetClient(subscriptionID)
}
func NewPercentileTargetClientWithBaseURI(baseURI string, subscriptionID string) PercentileTargetClient {
	return original.NewPercentileTargetClientWithBaseURI(baseURI, subscriptionID)
}
func NewSQLResourcesClient(subscriptionID string) SQLResourcesClient {
	return original.NewSQLResourcesClient(subscriptionID)
}
func NewSQLResourcesClientWithBaseURI(baseURI string, subscriptionID string) SQLResourcesClient {
	return original.NewSQLResourcesClientWithBaseURI(baseURI, subscriptionID)
}
func NewTableResourcesClient(subscriptionID string) TableResourcesClient {
	return original.NewTableResourcesClient(subscriptionID)
}
func NewTableResourcesClientWithBaseURI(baseURI string, subscriptionID string) TableResourcesClient {
	return original.NewTableResourcesClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleCompositePathSortOrderValues() []CompositePathSortOrder {
	return original.PossibleCompositePathSortOrderValues()
}
func PossibleConflictResolutionModeValues() []ConflictResolutionMode {
	return original.PossibleConflictResolutionModeValues()
}
func PossibleConnectorOfferValues() []ConnectorOffer {
	return original.PossibleConnectorOfferValues()
}
func PossibleDataTypeValues() []DataType {
	return original.PossibleDataTypeValues()
}
func PossibleDatabaseAccountKindValues() []DatabaseAccountKind {
	return original.PossibleDatabaseAccountKindValues()
}
func PossibleDatabaseAccountOfferTypeValues() []DatabaseAccountOfferType {
	return original.PossibleDatabaseAccountOfferTypeValues()
}
func PossibleDefaultConsistencyLevelValues() []DefaultConsistencyLevel {
	return original.PossibleDefaultConsistencyLevelValues()
}
func PossibleIndexKindValues() []IndexKind {
	return original.PossibleIndexKindValues()
}
func PossibleIndexingModeValues() []IndexingMode {
	return original.PossibleIndexingModeValues()
}
func PossibleKeyKindValues() []KeyKind {
	return original.PossibleKeyKindValues()
}
func PossiblePartitionKindValues() []PartitionKind {
	return original.PossiblePartitionKindValues()
}
func PossiblePrimaryAggregationTypeValues() []PrimaryAggregationType {
	return original.PossiblePrimaryAggregationTypeValues()
}
func PossibleSpatialTypeValues() []SpatialType {
	return original.PossibleSpatialTypeValues()
}
func PossibleTriggerOperationValues() []TriggerOperation {
	return original.PossibleTriggerOperationValues()
}
func PossibleTriggerTypeValues() []TriggerType {
	return original.PossibleTriggerTypeValues()
}
func PossibleUnitTypeValues() []UnitType {
	return original.PossibleUnitTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
