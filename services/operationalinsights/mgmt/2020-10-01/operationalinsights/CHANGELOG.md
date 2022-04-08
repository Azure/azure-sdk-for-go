# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. CapacityReservationProperties.MaxCapacity
1. WorkspaceSku.MaxCapacityReservationLevel

## Additive Changes

### New Constants

1. DataSourceKind.ApplicationInsights
1. DataSourceKind.AzureActivityLog
1. DataSourceKind.AzureAuditLog
1. DataSourceKind.ChangeTrackingContentLocation
1. DataSourceKind.ChangeTrackingCustomPath
1. DataSourceKind.ChangeTrackingDataTypeConfiguration
1. DataSourceKind.ChangeTrackingDefaultRegistry
1. DataSourceKind.ChangeTrackingLinuxPath
1. DataSourceKind.ChangeTrackingPath
1. DataSourceKind.ChangeTrackingRegistry
1. DataSourceKind.ChangeTrackingServices
1. DataSourceKind.CustomLog
1. DataSourceKind.CustomLogCollection
1. DataSourceKind.DNSAnalytics
1. DataSourceKind.GenericDataSource
1. DataSourceKind.IISLogs
1. DataSourceKind.ImportComputerGroup
1. DataSourceKind.Itsm
1. DataSourceKind.LinuxChangeTrackingPath
1. DataSourceKind.LinuxPerformanceCollection
1. DataSourceKind.LinuxPerformanceObject
1. DataSourceKind.LinuxSyslog
1. DataSourceKind.LinuxSyslogCollection
1. DataSourceKind.NetworkMonitoring
1. DataSourceKind.Office365
1. DataSourceKind.SQLDataClassification
1. DataSourceKind.SecurityCenterSecurityWindowsBaselineConfiguration
1. DataSourceKind.SecurityEventCollectionConfiguration
1. DataSourceKind.SecurityInsightsSecurityEventCollectionConfiguration
1. DataSourceKind.SecurityWindowsBaselineConfiguration
1. DataSourceKind.WindowsEvent
1. DataSourceKind.WindowsPerformanceCounter
1. DataSourceKind.WindowsTelemetry
1. DataSourceType.Alerts
1. DataSourceType.AzureWatson
1. DataSourceType.CustomLogs
1. DataSourceType.Query
1. LinkedServiceEntityStatus.LinkedServiceEntityStatusDeleting
1. LinkedServiceEntityStatus.LinkedServiceEntityStatusProvisioningAccount
1. LinkedServiceEntityStatus.LinkedServiceEntityStatusSucceeded
1. LinkedServiceEntityStatus.LinkedServiceEntityStatusUpdating
1. PurgeState.Completed
1. PurgeState.Pending
1. SearchSortEnum.Asc
1. SearchSortEnum.Desc
1. SkuNameEnum.SkuNameEnumCapacityReservation
1. SkuNameEnum.SkuNameEnumFree
1. SkuNameEnum.SkuNameEnumPerGB2018
1. SkuNameEnum.SkuNameEnumPerNode
1. SkuNameEnum.SkuNameEnumPremium
1. SkuNameEnum.SkuNameEnumStandalone
1. SkuNameEnum.SkuNameEnumStandard
1. StorageInsightState.ERROR
1. StorageInsightState.OK
1. Type.TypeEventHub
1. Type.TypeStorageAccount

### New Funcs

1. *DataExport.UnmarshalJSON([]byte) error
1. *DataExportProperties.UnmarshalJSON([]byte) error
1. *DataSourceListResultIterator.Next() error
1. *DataSourceListResultIterator.NextWithContext(context.Context) error
1. *DataSourceListResultPage.Next() error
1. *DataSourceListResultPage.NextWithContext(context.Context) error
1. *Destination.UnmarshalJSON([]byte) error
1. *LinkedService.UnmarshalJSON([]byte) error
1. *LinkedServicesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *LinkedServicesDeleteFuture.UnmarshalJSON([]byte) error
1. *LinkedStorageAccountsResource.UnmarshalJSON([]byte) error
1. *ManagementGroup.UnmarshalJSON([]byte) error
1. *SavedSearch.UnmarshalJSON([]byte) error
1. *StorageInsight.UnmarshalJSON([]byte) error
1. *StorageInsightListResultIterator.Next() error
1. *StorageInsightListResultIterator.NextWithContext(context.Context) error
1. *StorageInsightListResultPage.Next() error
1. *StorageInsightListResultPage.NextWithContext(context.Context) error
1. *Table.UnmarshalJSON([]byte) error
1. AvailableServiceTier.MarshalJSON() ([]byte, error)
1. AvailableServiceTiersClient.ListByWorkspace(context.Context, string, string) (ListAvailableServiceTier, error)
1. AvailableServiceTiersClient.ListByWorkspacePreparer(context.Context, string, string) (*http.Request, error)
1. AvailableServiceTiersClient.ListByWorkspaceResponder(*http.Response) (ListAvailableServiceTier, error)
1. AvailableServiceTiersClient.ListByWorkspaceSender(*http.Request) (*http.Response, error)
1. DataExport.MarshalJSON() ([]byte, error)
1. DataExportProperties.MarshalJSON() ([]byte, error)
1. DataExportsClient.CreateOrUpdate(context.Context, string, string, string, DataExport) (DataExport, error)
1. DataExportsClient.CreateOrUpdatePreparer(context.Context, string, string, string, DataExport) (*http.Request, error)
1. DataExportsClient.CreateOrUpdateResponder(*http.Response) (DataExport, error)
1. DataExportsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. DataExportsClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. DataExportsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. DataExportsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. DataExportsClient.DeleteSender(*http.Request) (*http.Response, error)
1. DataExportsClient.Get(context.Context, string, string, string) (DataExport, error)
1. DataExportsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. DataExportsClient.GetResponder(*http.Response) (DataExport, error)
1. DataExportsClient.GetSender(*http.Request) (*http.Response, error)
1. DataExportsClient.ListByWorkspace(context.Context, string, string) (DataExportListResult, error)
1. DataExportsClient.ListByWorkspacePreparer(context.Context, string, string) (*http.Request, error)
1. DataExportsClient.ListByWorkspaceResponder(*http.Response) (DataExportListResult, error)
1. DataExportsClient.ListByWorkspaceSender(*http.Request) (*http.Response, error)
1. DataSource.MarshalJSON() ([]byte, error)
1. DataSourceListResult.IsEmpty() bool
1. DataSourceListResultIterator.NotDone() bool
1. DataSourceListResultIterator.Response() DataSourceListResult
1. DataSourceListResultIterator.Value() DataSource
1. DataSourceListResultPage.NotDone() bool
1. DataSourceListResultPage.Response() DataSourceListResult
1. DataSourceListResultPage.Values() []DataSource
1. DataSourcesClient.CreateOrUpdate(context.Context, string, string, string, DataSource) (DataSource, error)
1. DataSourcesClient.CreateOrUpdatePreparer(context.Context, string, string, string, DataSource) (*http.Request, error)
1. DataSourcesClient.CreateOrUpdateResponder(*http.Response) (DataSource, error)
1. DataSourcesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. DataSourcesClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. DataSourcesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. DataSourcesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. DataSourcesClient.DeleteSender(*http.Request) (*http.Response, error)
1. DataSourcesClient.Get(context.Context, string, string, string) (DataSource, error)
1. DataSourcesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. DataSourcesClient.GetResponder(*http.Response) (DataSource, error)
1. DataSourcesClient.GetSender(*http.Request) (*http.Response, error)
1. DataSourcesClient.ListByWorkspace(context.Context, string, string, string, string) (DataSourceListResultPage, error)
1. DataSourcesClient.ListByWorkspaceComplete(context.Context, string, string, string, string) (DataSourceListResultIterator, error)
1. DataSourcesClient.ListByWorkspacePreparer(context.Context, string, string, string, string) (*http.Request, error)
1. DataSourcesClient.ListByWorkspaceResponder(*http.Response) (DataSourceListResult, error)
1. DataSourcesClient.ListByWorkspaceSender(*http.Request) (*http.Response, error)
1. Destination.MarshalJSON() ([]byte, error)
1. GatewaysClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. GatewaysClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. GatewaysClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. GatewaysClient.DeleteSender(*http.Request) (*http.Response, error)
1. IntelligencePacksClient.Disable(context.Context, string, string, string) (autorest.Response, error)
1. IntelligencePacksClient.DisablePreparer(context.Context, string, string, string) (*http.Request, error)
1. IntelligencePacksClient.DisableResponder(*http.Response) (autorest.Response, error)
1. IntelligencePacksClient.DisableSender(*http.Request) (*http.Response, error)
1. IntelligencePacksClient.Enable(context.Context, string, string, string) (autorest.Response, error)
1. IntelligencePacksClient.EnablePreparer(context.Context, string, string, string) (*http.Request, error)
1. IntelligencePacksClient.EnableResponder(*http.Response) (autorest.Response, error)
1. IntelligencePacksClient.EnableSender(*http.Request) (*http.Response, error)
1. IntelligencePacksClient.List(context.Context, string, string) (ListIntelligencePack, error)
1. IntelligencePacksClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. IntelligencePacksClient.ListResponder(*http.Response) (ListIntelligencePack, error)
1. IntelligencePacksClient.ListSender(*http.Request) (*http.Response, error)
1. LinkedService.MarshalJSON() ([]byte, error)
1. LinkedServicesClient.CreateOrUpdate(context.Context, string, string, string, LinkedService) (LinkedServicesCreateOrUpdateFuture, error)
1. LinkedServicesClient.CreateOrUpdatePreparer(context.Context, string, string, string, LinkedService) (*http.Request, error)
1. LinkedServicesClient.CreateOrUpdateResponder(*http.Response) (LinkedService, error)
1. LinkedServicesClient.CreateOrUpdateSender(*http.Request) (LinkedServicesCreateOrUpdateFuture, error)
1. LinkedServicesClient.Delete(context.Context, string, string, string) (LinkedServicesDeleteFuture, error)
1. LinkedServicesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. LinkedServicesClient.DeleteResponder(*http.Response) (LinkedService, error)
1. LinkedServicesClient.DeleteSender(*http.Request) (LinkedServicesDeleteFuture, error)
1. LinkedServicesClient.Get(context.Context, string, string, string) (LinkedService, error)
1. LinkedServicesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. LinkedServicesClient.GetResponder(*http.Response) (LinkedService, error)
1. LinkedServicesClient.GetSender(*http.Request) (*http.Response, error)
1. LinkedServicesClient.ListByWorkspace(context.Context, string, string) (LinkedServiceListResult, error)
1. LinkedServicesClient.ListByWorkspacePreparer(context.Context, string, string) (*http.Request, error)
1. LinkedServicesClient.ListByWorkspaceResponder(*http.Response) (LinkedServiceListResult, error)
1. LinkedServicesClient.ListByWorkspaceSender(*http.Request) (*http.Response, error)
1. LinkedStorageAccountsClient.CreateOrUpdate(context.Context, string, string, DataSourceType, LinkedStorageAccountsResource) (LinkedStorageAccountsResource, error)
1. LinkedStorageAccountsClient.CreateOrUpdatePreparer(context.Context, string, string, DataSourceType, LinkedStorageAccountsResource) (*http.Request, error)
1. LinkedStorageAccountsClient.CreateOrUpdateResponder(*http.Response) (LinkedStorageAccountsResource, error)
1. LinkedStorageAccountsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. LinkedStorageAccountsClient.Delete(context.Context, string, string, DataSourceType) (autorest.Response, error)
1. LinkedStorageAccountsClient.DeletePreparer(context.Context, string, string, DataSourceType) (*http.Request, error)
1. LinkedStorageAccountsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. LinkedStorageAccountsClient.DeleteSender(*http.Request) (*http.Response, error)
1. LinkedStorageAccountsClient.Get(context.Context, string, string, DataSourceType) (LinkedStorageAccountsResource, error)
1. LinkedStorageAccountsClient.GetPreparer(context.Context, string, string, DataSourceType) (*http.Request, error)
1. LinkedStorageAccountsClient.GetResponder(*http.Response) (LinkedStorageAccountsResource, error)
1. LinkedStorageAccountsClient.GetSender(*http.Request) (*http.Response, error)
1. LinkedStorageAccountsClient.ListByWorkspace(context.Context, string, string) (LinkedStorageAccountsListResult, error)
1. LinkedStorageAccountsClient.ListByWorkspacePreparer(context.Context, string, string) (*http.Request, error)
1. LinkedStorageAccountsClient.ListByWorkspaceResponder(*http.Response) (LinkedStorageAccountsListResult, error)
1. LinkedStorageAccountsClient.ListByWorkspaceSender(*http.Request) (*http.Response, error)
1. LinkedStorageAccountsProperties.MarshalJSON() ([]byte, error)
1. LinkedStorageAccountsResource.MarshalJSON() ([]byte, error)
1. ManagementGroup.MarshalJSON() ([]byte, error)
1. ManagementGroupsClient.List(context.Context, string, string) (WorkspaceListManagementGroupsResult, error)
1. ManagementGroupsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. ManagementGroupsClient.ListResponder(*http.Response) (WorkspaceListManagementGroupsResult, error)
1. ManagementGroupsClient.ListSender(*http.Request) (*http.Response, error)
1. NewAvailableServiceTiersClient(string) AvailableServiceTiersClient
1. NewAvailableServiceTiersClientWithBaseURI(string, string) AvailableServiceTiersClient
1. NewDataExportsClient(string) DataExportsClient
1. NewDataExportsClientWithBaseURI(string, string) DataExportsClient
1. NewDataSourceListResultIterator(DataSourceListResultPage) DataSourceListResultIterator
1. NewDataSourceListResultPage(DataSourceListResult, func(context.Context, DataSourceListResult) (DataSourceListResult, error)) DataSourceListResultPage
1. NewDataSourcesClient(string) DataSourcesClient
1. NewDataSourcesClientWithBaseURI(string, string) DataSourcesClient
1. NewGatewaysClient(string) GatewaysClient
1. NewGatewaysClientWithBaseURI(string, string) GatewaysClient
1. NewIntelligencePacksClient(string) IntelligencePacksClient
1. NewIntelligencePacksClientWithBaseURI(string, string) IntelligencePacksClient
1. NewLinkedServicesClient(string) LinkedServicesClient
1. NewLinkedServicesClientWithBaseURI(string, string) LinkedServicesClient
1. NewLinkedStorageAccountsClient(string) LinkedStorageAccountsClient
1. NewLinkedStorageAccountsClientWithBaseURI(string, string) LinkedStorageAccountsClient
1. NewManagementGroupsClient(string) ManagementGroupsClient
1. NewManagementGroupsClientWithBaseURI(string, string) ManagementGroupsClient
1. NewOperationStatusesClient(string) OperationStatusesClient
1. NewOperationStatusesClientWithBaseURI(string, string) OperationStatusesClient
1. NewSavedSearchesClient(string) SavedSearchesClient
1. NewSavedSearchesClientWithBaseURI(string, string) SavedSearchesClient
1. NewSchemaClient(string) SchemaClient
1. NewSchemaClientWithBaseURI(string, string) SchemaClient
1. NewSharedKeysClient(string) SharedKeysClient
1. NewSharedKeysClientWithBaseURI(string, string) SharedKeysClient
1. NewStorageInsightConfigsClient(string) StorageInsightConfigsClient
1. NewStorageInsightConfigsClientWithBaseURI(string, string) StorageInsightConfigsClient
1. NewStorageInsightListResultIterator(StorageInsightListResultPage) StorageInsightListResultIterator
1. NewStorageInsightListResultPage(StorageInsightListResult, func(context.Context, StorageInsightListResult) (StorageInsightListResult, error)) StorageInsightListResultPage
1. NewTablesClient(string) TablesClient
1. NewTablesClientWithBaseURI(string, string) TablesClient
1. NewUsagesClient(string) UsagesClient
1. NewUsagesClientWithBaseURI(string, string) UsagesClient
1. NewWorkspacePurgeClient(string) WorkspacePurgeClient
1. NewWorkspacePurgeClientWithBaseURI(string, string) WorkspacePurgeClient
1. OperationStatusesClient.Get(context.Context, string, string) (OperationStatus, error)
1. OperationStatusesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. OperationStatusesClient.GetResponder(*http.Response) (OperationStatus, error)
1. OperationStatusesClient.GetSender(*http.Request) (*http.Response, error)
1. PossibleDataSourceKindValues() []DataSourceKind
1. PossibleDataSourceTypeValues() []DataSourceType
1. PossibleLinkedServiceEntityStatusValues() []LinkedServiceEntityStatus
1. PossiblePurgeStateValues() []PurgeState
1. PossibleSearchSortEnumValues() []SearchSortEnum
1. PossibleSkuNameEnumValues() []SkuNameEnum
1. PossibleStorageInsightStateValues() []StorageInsightState
1. PossibleTypeValues() []Type
1. SavedSearch.MarshalJSON() ([]byte, error)
1. SavedSearchesClient.CreateOrUpdate(context.Context, string, string, string, SavedSearch) (SavedSearch, error)
1. SavedSearchesClient.CreateOrUpdatePreparer(context.Context, string, string, string, SavedSearch) (*http.Request, error)
1. SavedSearchesClient.CreateOrUpdateResponder(*http.Response) (SavedSearch, error)
1. SavedSearchesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. SavedSearchesClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. SavedSearchesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. SavedSearchesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. SavedSearchesClient.DeleteSender(*http.Request) (*http.Response, error)
1. SavedSearchesClient.Get(context.Context, string, string, string) (SavedSearch, error)
1. SavedSearchesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. SavedSearchesClient.GetResponder(*http.Response) (SavedSearch, error)
1. SavedSearchesClient.GetSender(*http.Request) (*http.Response, error)
1. SavedSearchesClient.ListByWorkspace(context.Context, string, string) (SavedSearchesListResult, error)
1. SavedSearchesClient.ListByWorkspacePreparer(context.Context, string, string) (*http.Request, error)
1. SavedSearchesClient.ListByWorkspaceResponder(*http.Response) (SavedSearchesListResult, error)
1. SavedSearchesClient.ListByWorkspaceSender(*http.Request) (*http.Response, error)
1. SchemaClient.Get(context.Context, string, string) (SearchGetSchemaResponse, error)
1. SchemaClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. SchemaClient.GetResponder(*http.Response) (SearchGetSchemaResponse, error)
1. SchemaClient.GetSender(*http.Request) (*http.Response, error)
1. SharedKeysClient.GetSharedKeys(context.Context, string, string) (SharedKeys, error)
1. SharedKeysClient.GetSharedKeysPreparer(context.Context, string, string) (*http.Request, error)
1. SharedKeysClient.GetSharedKeysResponder(*http.Response) (SharedKeys, error)
1. SharedKeysClient.GetSharedKeysSender(*http.Request) (*http.Response, error)
1. SharedKeysClient.Regenerate(context.Context, string, string) (SharedKeys, error)
1. SharedKeysClient.RegeneratePreparer(context.Context, string, string) (*http.Request, error)
1. SharedKeysClient.RegenerateResponder(*http.Response) (SharedKeys, error)
1. SharedKeysClient.RegenerateSender(*http.Request) (*http.Response, error)
1. StorageInsight.MarshalJSON() ([]byte, error)
1. StorageInsightConfigsClient.CreateOrUpdate(context.Context, string, string, string, StorageInsight) (StorageInsight, error)
1. StorageInsightConfigsClient.CreateOrUpdatePreparer(context.Context, string, string, string, StorageInsight) (*http.Request, error)
1. StorageInsightConfigsClient.CreateOrUpdateResponder(*http.Response) (StorageInsight, error)
1. StorageInsightConfigsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. StorageInsightConfigsClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. StorageInsightConfigsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. StorageInsightConfigsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. StorageInsightConfigsClient.DeleteSender(*http.Request) (*http.Response, error)
1. StorageInsightConfigsClient.Get(context.Context, string, string, string) (StorageInsight, error)
1. StorageInsightConfigsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. StorageInsightConfigsClient.GetResponder(*http.Response) (StorageInsight, error)
1. StorageInsightConfigsClient.GetSender(*http.Request) (*http.Response, error)
1. StorageInsightConfigsClient.ListByWorkspace(context.Context, string, string) (StorageInsightListResultPage, error)
1. StorageInsightConfigsClient.ListByWorkspaceComplete(context.Context, string, string) (StorageInsightListResultIterator, error)
1. StorageInsightConfigsClient.ListByWorkspacePreparer(context.Context, string, string) (*http.Request, error)
1. StorageInsightConfigsClient.ListByWorkspaceResponder(*http.Response) (StorageInsightListResult, error)
1. StorageInsightConfigsClient.ListByWorkspaceSender(*http.Request) (*http.Response, error)
1. StorageInsightListResult.IsEmpty() bool
1. StorageInsightListResultIterator.NotDone() bool
1. StorageInsightListResultIterator.Response() StorageInsightListResult
1. StorageInsightListResultIterator.Value() StorageInsight
1. StorageInsightListResultPage.NotDone() bool
1. StorageInsightListResultPage.Response() StorageInsightListResult
1. StorageInsightListResultPage.Values() []StorageInsight
1. StorageInsightProperties.MarshalJSON() ([]byte, error)
1. Table.MarshalJSON() ([]byte, error)
1. TablesClient.Get(context.Context, string, string, string) (Table, error)
1. TablesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. TablesClient.GetResponder(*http.Response) (Table, error)
1. TablesClient.GetSender(*http.Request) (*http.Response, error)
1. TablesClient.ListByWorkspace(context.Context, string, string) (TablesListResult, error)
1. TablesClient.ListByWorkspacePreparer(context.Context, string, string) (*http.Request, error)
1. TablesClient.ListByWorkspaceResponder(*http.Response) (TablesListResult, error)
1. TablesClient.ListByWorkspaceSender(*http.Request) (*http.Response, error)
1. TablesClient.Update(context.Context, string, string, string, Table) (Table, error)
1. TablesClient.UpdatePreparer(context.Context, string, string, string, Table) (*http.Request, error)
1. TablesClient.UpdateResponder(*http.Response) (Table, error)
1. TablesClient.UpdateSender(*http.Request) (*http.Response, error)
1. UsagesClient.List(context.Context, string, string) (WorkspaceListUsagesResult, error)
1. UsagesClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. UsagesClient.ListResponder(*http.Response) (WorkspaceListUsagesResult, error)
1. UsagesClient.ListSender(*http.Request) (*http.Response, error)
1. WorkspacePurgeClient.GetPurgeStatus(context.Context, string, string, string) (WorkspacePurgeStatusResponse, error)
1. WorkspacePurgeClient.GetPurgeStatusPreparer(context.Context, string, string, string) (*http.Request, error)
1. WorkspacePurgeClient.GetPurgeStatusResponder(*http.Response) (WorkspacePurgeStatusResponse, error)
1. WorkspacePurgeClient.GetPurgeStatusSender(*http.Request) (*http.Response, error)
1. WorkspacePurgeClient.Purge(context.Context, string, string, WorkspacePurgeBody) (WorkspacePurgeResponse, error)
1. WorkspacePurgeClient.PurgePreparer(context.Context, string, string, WorkspacePurgeBody) (*http.Request, error)
1. WorkspacePurgeClient.PurgeResponder(*http.Response) (WorkspacePurgeResponse, error)
1. WorkspacePurgeClient.PurgeSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. AvailableServiceTier
1. AvailableServiceTiersClient
1. CoreSummary
1. DataExport
1. DataExportListResult
1. DataExportProperties
1. DataExportsClient
1. DataSource
1. DataSourceFilter
1. DataSourceListResult
1. DataSourceListResultIterator
1. DataSourceListResultPage
1. DataSourcesClient
1. Destination
1. DestinationMetaData
1. GatewaysClient
1. IntelligencePack
1. IntelligencePacksClient
1. LinkedService
1. LinkedServiceListResult
1. LinkedServiceProperties
1. LinkedServicesClient
1. LinkedServicesCreateOrUpdateFuture
1. LinkedServicesDeleteFuture
1. LinkedStorageAccountsClient
1. LinkedStorageAccountsListResult
1. LinkedStorageAccountsProperties
1. LinkedStorageAccountsResource
1. ListAvailableServiceTier
1. ListIntelligencePack
1. ManagementGroup
1. ManagementGroupProperties
1. ManagementGroupsClient
1. MetricName
1. OperationStatus
1. OperationStatusesClient
1. SavedSearch
1. SavedSearchProperties
1. SavedSearchesClient
1. SavedSearchesListResult
1. SchemaClient
1. SearchGetSchemaResponse
1. SearchMetadata
1. SearchMetadataSchema
1. SearchSchemaValue
1. SearchSort
1. SharedKeys
1. SharedKeysClient
1. StorageAccount
1. StorageInsight
1. StorageInsightConfigsClient
1. StorageInsightListResult
1. StorageInsightListResultIterator
1. StorageInsightListResultPage
1. StorageInsightProperties
1. StorageInsightStatus
1. Table
1. TableProperties
1. TablesClient
1. TablesListResult
1. Tag
1. UsageMetric
1. UsagesClient
1. WorkspaceFeatures
1. WorkspaceListManagementGroupsResult
1. WorkspaceListUsagesResult
1. WorkspacePurgeBody
1. WorkspacePurgeBodyFilters
1. WorkspacePurgeClient
1. WorkspacePurgeResponse
1. WorkspacePurgeStatusResponse

#### New Struct Fields

1. ClusterPatchProperties.BillingType
1. WorkspaceProperties.CreatedDate
1. WorkspaceProperties.Features
1. WorkspaceProperties.ModifiedDate
