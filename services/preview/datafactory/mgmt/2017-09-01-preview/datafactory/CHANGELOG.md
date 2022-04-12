# Unreleased

## Breaking Changes

### Removed Constants

1. IntegrationRuntimeState.Initial
1. IntegrationRuntimeState.Limited
1. IntegrationRuntimeState.NeedRegistration
1. IntegrationRuntimeState.Offline
1. IntegrationRuntimeState.Stopping
1. ManagedIntegrationRuntimeNodeStatus.ManagedIntegrationRuntimeNodeStatusAvailable
1. ManagedIntegrationRuntimeNodeStatus.ManagedIntegrationRuntimeNodeStatusRecycling
1. ManagedIntegrationRuntimeNodeStatus.ManagedIntegrationRuntimeNodeStatusStarting
1. ManagedIntegrationRuntimeNodeStatus.ManagedIntegrationRuntimeNodeStatusUnavailable
1. TriggerRuntimeState.TriggerRuntimeStateDisabled
1. TriggerRuntimeState.TriggerRuntimeStateStarted
1. TriggerRuntimeState.TriggerRuntimeStateStopped

### Signature Changes

#### Const Types

1. Online changed type from IntegrationRuntimeState to DynamicsDeploymentType
1. Started changed type from IntegrationRuntimeState to TriggerRuntimeState
1. Starting changed type from IntegrationRuntimeState to ManagedIntegrationRuntimeNodeStatus
1. Stopped changed type from IntegrationRuntimeState to TriggerRuntimeState

#### Struct Fields

1. AzureSearchIndexSink.WriteBehavior changed type from interface{} to AzureSearchIndexWriteBehaviorType
1. CassandraSource.ConsistencyLevel changed type from interface{} to CassandraSourceReadConsistencyLevels
1. DatasetDeflateCompression.Level changed type from interface{} to DatasetCompressionLevel
1. DatasetGZipCompression.Level changed type from interface{} to DatasetCompressionLevel
1. DatasetZipDeflateCompression.Level changed type from interface{} to DatasetCompressionLevel
1. DynamicsLinkedServiceTypeProperties.AuthenticationType changed type from interface{} to DynamicsAuthenticationType
1. DynamicsLinkedServiceTypeProperties.DeploymentType changed type from interface{} to DynamicsDeploymentType
1. DynamicsSink.WriteBehavior changed type from interface{} to *string
1. PolybaseSettings.RejectType changed type from interface{} to PolybaseSettingsRejectType
1. SalesforceSink.WriteBehavior changed type from interface{} to SalesforceSinkWriteBehavior
1. SalesforceSource.ReadBehavior changed type from interface{} to SalesforceSourceReadBehavior
1. SapCloudForCustomerSink.WriteBehavior changed type from interface{} to SapCloudForCustomerSinkWriteBehavior
1. StoredProcedureParameter.Type changed type from interface{} to StoredProcedureParameterType

## Additive Changes

### New Constants

1. AzureSearchIndexWriteBehaviorType.Merge
1. AzureSearchIndexWriteBehaviorType.Upload
1. CassandraSourceReadConsistencyLevels.ALL
1. CassandraSourceReadConsistencyLevels.EACHQUORUM
1. CassandraSourceReadConsistencyLevels.LOCALONE
1. CassandraSourceReadConsistencyLevels.LOCALQUORUM
1. CassandraSourceReadConsistencyLevels.LOCALSERIAL
1. CassandraSourceReadConsistencyLevels.ONE
1. CassandraSourceReadConsistencyLevels.QUORUM
1. CassandraSourceReadConsistencyLevels.SERIAL
1. CassandraSourceReadConsistencyLevels.THREE
1. CassandraSourceReadConsistencyLevels.TWO
1. DatasetCompressionLevel.Fastest
1. DatasetCompressionLevel.Optimal
1. DynamicsAuthenticationType.Ifd
1. DynamicsAuthenticationType.Office365
1. DynamicsDeploymentType.OnPremisesWithIfd
1. IntegrationRuntimeState.IntegrationRuntimeStateInitial
1. IntegrationRuntimeState.IntegrationRuntimeStateLimited
1. IntegrationRuntimeState.IntegrationRuntimeStateNeedRegistration
1. IntegrationRuntimeState.IntegrationRuntimeStateOffline
1. IntegrationRuntimeState.IntegrationRuntimeStateOnline
1. IntegrationRuntimeState.IntegrationRuntimeStateStarted
1. IntegrationRuntimeState.IntegrationRuntimeStateStarting
1. IntegrationRuntimeState.IntegrationRuntimeStateStopped
1. IntegrationRuntimeState.IntegrationRuntimeStateStopping
1. ManagedIntegrationRuntimeNodeStatus.Available
1. ManagedIntegrationRuntimeNodeStatus.Recycling
1. ManagedIntegrationRuntimeNodeStatus.Unavailable
1. PolybaseSettingsRejectType.Percentage
1. PolybaseSettingsRejectType.Value
1. SalesforceSinkWriteBehavior.Insert
1. SalesforceSinkWriteBehavior.Upsert
1. SalesforceSourceReadBehavior.Query
1. SalesforceSourceReadBehavior.QueryAll
1. SapCloudForCustomerSinkWriteBehavior.SapCloudForCustomerSinkWriteBehaviorInsert
1. SapCloudForCustomerSinkWriteBehavior.SapCloudForCustomerSinkWriteBehaviorUpdate
1. StoredProcedureParameterType.Boolean
1. StoredProcedureParameterType.Date
1. StoredProcedureParameterType.Decimal
1. StoredProcedureParameterType.GUID
1. StoredProcedureParameterType.Int
1. StoredProcedureParameterType.Int64
1. StoredProcedureParameterType.String
1. TriggerRuntimeState.Disabled

### New Funcs

1. PossibleAzureSearchIndexWriteBehaviorTypeValues() []AzureSearchIndexWriteBehaviorType
1. PossibleCassandraSourceReadConsistencyLevelsValues() []CassandraSourceReadConsistencyLevels
1. PossibleDatasetCompressionLevelValues() []DatasetCompressionLevel
1. PossibleDynamicsAuthenticationTypeValues() []DynamicsAuthenticationType
1. PossibleDynamicsDeploymentTypeValues() []DynamicsDeploymentType
1. PossiblePolybaseSettingsRejectTypeValues() []PolybaseSettingsRejectType
1. PossibleSalesforceSinkWriteBehaviorValues() []SalesforceSinkWriteBehavior
1. PossibleSalesforceSourceReadBehaviorValues() []SalesforceSourceReadBehavior
1. PossibleSapCloudForCustomerSinkWriteBehaviorValues() []SapCloudForCustomerSinkWriteBehavior
1. PossibleStoredProcedureParameterTypeValues() []StoredProcedureParameterType

### Struct Changes

#### New Struct Fields

1. WebActivityTypeProperties.DisableCertValidation
