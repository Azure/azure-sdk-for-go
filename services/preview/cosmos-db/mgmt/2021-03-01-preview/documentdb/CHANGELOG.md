# Change History

## Breaking Changes

### Removed Constants

1. APIType.Cassandra
1. APIType.Gremlin
1. APIType.GremlinV2
1. APIType.MongoDB
1. APIType.SQL
1. APIType.Table
1. BackupPolicyType.Continuous
1. BackupPolicyType.Periodic
1. BackupStorageRedundancy.Geo
1. BackupStorageRedundancy.Local
1. BackupStorageRedundancy.Zone
1. CompositePathSortOrder.Ascending
1. CompositePathSortOrder.Descending
1. ConflictResolutionMode.Custom
1. ConflictResolutionMode.LastWriterWins
1. ConnectorOffer.Small
1. CreateMode.Default
1. CreateMode.Restore
1. CreateModeBasicDatabaseAccountCreateUpdateProperties.CreateModeDatabaseAccountCreateUpdateProperties
1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User
1. DataType.LineString
1. DataType.MultiPolygon
1. DataType.Number
1. DataType.Point
1. DataType.Polygon
1. DataType.String
1. DatabaseAccountOfferType.Standard
1. DefaultConsistencyLevel.BoundedStaleness
1. DefaultConsistencyLevel.ConsistentPrefix
1. DefaultConsistencyLevel.Eventual
1. DefaultConsistencyLevel.Session
1. DefaultConsistencyLevel.Strong
1. IndexKind.Hash
1. IndexKind.Range
1. IndexKind.Spatial
1. IndexingMode.Consistent
1. IndexingMode.Lazy
1. IndexingMode.None
1. KeyKind.Primary
1. KeyKind.PrimaryReadonly
1. KeyKind.Secondary
1. KeyKind.SecondaryReadonly
1. ManagedCassandraProvisioningState.Canceled
1. ManagedCassandraProvisioningState.Creating
1. ManagedCassandraProvisioningState.Deleting
1. ManagedCassandraProvisioningState.Failed
1. ManagedCassandraProvisioningState.Succeeded
1. ManagedCassandraProvisioningState.Updating
1. NodeState.Joining
1. NodeState.Leaving
1. NodeState.Moving
1. NodeState.Normal
1. NodeState.Stopped
1. NodeStatus.Down
1. NodeStatus.Up
1. OperationType.Create
1. OperationType.Delete
1. OperationType.Replace
1. OperationType.SystemOperation
1. PublicNetworkAccess.Disabled
1. PublicNetworkAccess.Enabled
1. RestoreMode.PointInTime
1. RoleDefinitionType.BuiltInRole
1. RoleDefinitionType.CustomRole
1. ServerVersion.FourFullStopZero
1. ServerVersion.ThreeFullStopSix
1. ServerVersion.ThreeFullStopTwo
1. TriggerType.Post
1. TriggerType.Pre
1. UnitType.Bytes
1. UnitType.BytesPerSecond
1. UnitType.Count
1. UnitType.CountPerSecond
1. UnitType.Milliseconds
1. UnitType.Percent
1. UnitType.Seconds

### Signature Changes

#### Const Types

1. CreateModeDefault changed type from CreateModeBasicDatabaseAccountCreateUpdateProperties to CreateMode
1. CreateModeRestore changed type from CreateModeBasicDatabaseAccountCreateUpdateProperties to CreateMode

## Additive Changes

### New Constants

1. APIType.APITypeCassandra
1. APIType.APITypeGremlin
1. APIType.APITypeGremlinV2
1. APIType.APITypeMongoDB
1. APIType.APITypeSQL
1. APIType.APITypeTable
1. BackupPolicyType.BackupPolicyTypeContinuous
1. BackupPolicyType.BackupPolicyTypePeriodic
1. BackupStorageRedundancy.BackupStorageRedundancyGeo
1. BackupStorageRedundancy.BackupStorageRedundancyLocal
1. BackupStorageRedundancy.BackupStorageRedundancyZone
1. CompositePathSortOrder.CompositePathSortOrderAscending
1. CompositePathSortOrder.CompositePathSortOrderDescending
1. ConflictResolutionMode.ConflictResolutionModeCustom
1. ConflictResolutionMode.ConflictResolutionModeLastWriterWins
1. ConnectorOffer.ConnectorOfferSmall
1. CreateModeBasicDatabaseAccountCreateUpdateProperties.CreateModeBasicDatabaseAccountCreateUpdatePropertiesCreateModeDatabaseAccountCreateUpdateProperties
1. CreateModeBasicDatabaseAccountCreateUpdateProperties.CreateModeBasicDatabaseAccountCreateUpdatePropertiesCreateModeDefault
1. CreateModeBasicDatabaseAccountCreateUpdateProperties.CreateModeBasicDatabaseAccountCreateUpdatePropertiesCreateModeRestore
1. CreatedByType.CreatedByTypeApplication
1. CreatedByType.CreatedByTypeKey
1. CreatedByType.CreatedByTypeManagedIdentity
1. CreatedByType.CreatedByTypeUser
1. DataType.DataTypeLineString
1. DataType.DataTypeMultiPolygon
1. DataType.DataTypeNumber
1. DataType.DataTypePoint
1. DataType.DataTypePolygon
1. DataType.DataTypeString
1. DatabaseAccountOfferType.DatabaseAccountOfferTypeStandard
1. DefaultConsistencyLevel.DefaultConsistencyLevelBoundedStaleness
1. DefaultConsistencyLevel.DefaultConsistencyLevelConsistentPrefix
1. DefaultConsistencyLevel.DefaultConsistencyLevelEventual
1. DefaultConsistencyLevel.DefaultConsistencyLevelSession
1. DefaultConsistencyLevel.DefaultConsistencyLevelStrong
1. IndexKind.IndexKindHash
1. IndexKind.IndexKindRange
1. IndexKind.IndexKindSpatial
1. IndexingMode.IndexingModeConsistent
1. IndexingMode.IndexingModeLazy
1. IndexingMode.IndexingModeNone
1. KeyKind.KeyKindPrimary
1. KeyKind.KeyKindPrimaryReadonly
1. KeyKind.KeyKindSecondary
1. KeyKind.KeyKindSecondaryReadonly
1. ManagedCassandraProvisioningState.ManagedCassandraProvisioningStateCanceled
1. ManagedCassandraProvisioningState.ManagedCassandraProvisioningStateCreating
1. ManagedCassandraProvisioningState.ManagedCassandraProvisioningStateDeleting
1. ManagedCassandraProvisioningState.ManagedCassandraProvisioningStateFailed
1. ManagedCassandraProvisioningState.ManagedCassandraProvisioningStateSucceeded
1. ManagedCassandraProvisioningState.ManagedCassandraProvisioningStateUpdating
1. NodeState.NodeStateJoining
1. NodeState.NodeStateLeaving
1. NodeState.NodeStateMoving
1. NodeState.NodeStateNormal
1. NodeState.NodeStateStopped
1. NodeStatus.NodeStatusDown
1. NodeStatus.NodeStatusUp
1. OperationType.OperationTypeCreate
1. OperationType.OperationTypeDelete
1. OperationType.OperationTypeReplace
1. OperationType.OperationTypeSystemOperation
1. PublicNetworkAccess.PublicNetworkAccessDisabled
1. PublicNetworkAccess.PublicNetworkAccessEnabled
1. RestoreMode.RestoreModePointInTime
1. RoleDefinitionType.RoleDefinitionTypeBuiltInRole
1. RoleDefinitionType.RoleDefinitionTypeCustomRole
1. ServerVersion.ServerVersionFourFullStopZero
1. ServerVersion.ServerVersionThreeFullStopSix
1. ServerVersion.ServerVersionThreeFullStopTwo
1. TriggerType.TriggerTypePost
1. TriggerType.TriggerTypePre
1. UnitType.UnitTypeBytes
1. UnitType.UnitTypeBytesPerSecond
1. UnitType.UnitTypeCount
1. UnitType.UnitTypeCountPerSecond
1. UnitType.UnitTypeMilliseconds
1. UnitType.UnitTypePercent
1. UnitType.UnitTypeSeconds

### New Funcs

1. BaseClient.LocationGet(context.Context, string) (LocationGetResult, error)
1. BaseClient.LocationGetPreparer(context.Context, string) (*http.Request, error)
1. BaseClient.LocationGetResponder(*http.Response) (LocationGetResult, error)
1. BaseClient.LocationGetSender(*http.Request) (*http.Response, error)
1. BaseClient.LocationList(context.Context) (LocationListResult, error)
1. BaseClient.LocationListPreparer(context.Context) (*http.Request, error)
1. BaseClient.LocationListResponder(*http.Response) (LocationListResult, error)
1. BaseClient.LocationListSender(*http.Request) (*http.Response, error)
1. LocationGetResult.MarshalJSON() ([]byte, error)
1. LocationListResult.MarshalJSON() ([]byte, error)
1. LocationProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. LocationGetResult
1. LocationListResult
1. LocationProperties
