//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armoperationalinsights

// AvailableServiceTiersClientListByWorkspaceResponse contains the response from method AvailableServiceTiersClient.ListByWorkspace.
type AvailableServiceTiersClientListByWorkspaceResponse struct {
	// Array of AvailableServiceTier
	AvailableServiceTierArray []*AvailableServiceTier
}

// ClustersClientCreateOrUpdateResponse contains the response from method ClustersClient.CreateOrUpdate.
type ClustersClientCreateOrUpdateResponse struct {
	Cluster
}

// ClustersClientDeleteResponse contains the response from method ClustersClient.Delete.
type ClustersClientDeleteResponse struct {
	// placeholder for future response values
}

// ClustersClientGetResponse contains the response from method ClustersClient.Get.
type ClustersClientGetResponse struct {
	Cluster
}

// ClustersClientListByResourceGroupResponse contains the response from method ClustersClient.ListByResourceGroup.
type ClustersClientListByResourceGroupResponse struct {
	ClusterListResult
}

// ClustersClientListResponse contains the response from method ClustersClient.List.
type ClustersClientListResponse struct {
	ClusterListResult
}

// ClustersClientUpdateResponse contains the response from method ClustersClient.Update.
type ClustersClientUpdateResponse struct {
	Cluster
}

// DataExportsClientCreateOrUpdateResponse contains the response from method DataExportsClient.CreateOrUpdate.
type DataExportsClientCreateOrUpdateResponse struct {
	DataExport
}

// DataExportsClientDeleteResponse contains the response from method DataExportsClient.Delete.
type DataExportsClientDeleteResponse struct {
	// placeholder for future response values
}

// DataExportsClientGetResponse contains the response from method DataExportsClient.Get.
type DataExportsClientGetResponse struct {
	DataExport
}

// DataExportsClientListByWorkspaceResponse contains the response from method DataExportsClient.ListByWorkspace.
type DataExportsClientListByWorkspaceResponse struct {
	DataExportListResult
}

// DataSourcesClientCreateOrUpdateResponse contains the response from method DataSourcesClient.CreateOrUpdate.
type DataSourcesClientCreateOrUpdateResponse struct {
	DataSource
}

// DataSourcesClientDeleteResponse contains the response from method DataSourcesClient.Delete.
type DataSourcesClientDeleteResponse struct {
	// placeholder for future response values
}

// DataSourcesClientGetResponse contains the response from method DataSourcesClient.Get.
type DataSourcesClientGetResponse struct {
	DataSource
}

// DataSourcesClientListByWorkspaceResponse contains the response from method DataSourcesClient.ListByWorkspace.
type DataSourcesClientListByWorkspaceResponse struct {
	DataSourceListResult
}

// DeletedWorkspacesClientListByResourceGroupResponse contains the response from method DeletedWorkspacesClient.ListByResourceGroup.
type DeletedWorkspacesClientListByResourceGroupResponse struct {
	WorkspaceListResult
}

// DeletedWorkspacesClientListResponse contains the response from method DeletedWorkspacesClient.List.
type DeletedWorkspacesClientListResponse struct {
	WorkspaceListResult
}

// GatewaysClientDeleteResponse contains the response from method GatewaysClient.Delete.
type GatewaysClientDeleteResponse struct {
	// placeholder for future response values
}

// IntelligencePacksClientDisableResponse contains the response from method IntelligencePacksClient.Disable.
type IntelligencePacksClientDisableResponse struct {
	// placeholder for future response values
}

// IntelligencePacksClientEnableResponse contains the response from method IntelligencePacksClient.Enable.
type IntelligencePacksClientEnableResponse struct {
	// placeholder for future response values
}

// IntelligencePacksClientListResponse contains the response from method IntelligencePacksClient.List.
type IntelligencePacksClientListResponse struct {
	// Array of IntelligencePack
	IntelligencePackArray []*IntelligencePack
}

// LinkedServicesClientCreateOrUpdateResponse contains the response from method LinkedServicesClient.CreateOrUpdate.
type LinkedServicesClientCreateOrUpdateResponse struct {
	LinkedService
}

// LinkedServicesClientDeleteResponse contains the response from method LinkedServicesClient.Delete.
type LinkedServicesClientDeleteResponse struct {
	LinkedService
}

// LinkedServicesClientGetResponse contains the response from method LinkedServicesClient.Get.
type LinkedServicesClientGetResponse struct {
	LinkedService
}

// LinkedServicesClientListByWorkspaceResponse contains the response from method LinkedServicesClient.ListByWorkspace.
type LinkedServicesClientListByWorkspaceResponse struct {
	LinkedServiceListResult
}

// LinkedStorageAccountsClientCreateOrUpdateResponse contains the response from method LinkedStorageAccountsClient.CreateOrUpdate.
type LinkedStorageAccountsClientCreateOrUpdateResponse struct {
	LinkedStorageAccountsResource
}

// LinkedStorageAccountsClientDeleteResponse contains the response from method LinkedStorageAccountsClient.Delete.
type LinkedStorageAccountsClientDeleteResponse struct {
	// placeholder for future response values
}

// LinkedStorageAccountsClientGetResponse contains the response from method LinkedStorageAccountsClient.Get.
type LinkedStorageAccountsClientGetResponse struct {
	LinkedStorageAccountsResource
}

// LinkedStorageAccountsClientListByWorkspaceResponse contains the response from method LinkedStorageAccountsClient.ListByWorkspace.
type LinkedStorageAccountsClientListByWorkspaceResponse struct {
	LinkedStorageAccountsListResult
}

// ManagementGroupsClientListResponse contains the response from method ManagementGroupsClient.List.
type ManagementGroupsClientListResponse struct {
	WorkspaceListManagementGroupsResult
}

// OperationStatusesClientGetResponse contains the response from method OperationStatusesClient.Get.
type OperationStatusesClientGetResponse struct {
	OperationStatus
}

// OperationsClientListResponse contains the response from method OperationsClient.List.
type OperationsClientListResponse struct {
	OperationListResult
}

// QueriesClientDeleteResponse contains the response from method QueriesClient.Delete.
type QueriesClientDeleteResponse struct {
	// placeholder for future response values
}

// QueriesClientGetResponse contains the response from method QueriesClient.Get.
type QueriesClientGetResponse struct {
	LogAnalyticsQueryPackQuery
}

// QueriesClientListResponse contains the response from method QueriesClient.List.
type QueriesClientListResponse struct {
	LogAnalyticsQueryPackQueryListResult
}

// QueriesClientPutResponse contains the response from method QueriesClient.Put.
type QueriesClientPutResponse struct {
	LogAnalyticsQueryPackQuery
}

// QueriesClientSearchResponse contains the response from method QueriesClient.Search.
type QueriesClientSearchResponse struct {
	LogAnalyticsQueryPackQueryListResult
}

// QueriesClientUpdateResponse contains the response from method QueriesClient.Update.
type QueriesClientUpdateResponse struct {
	LogAnalyticsQueryPackQuery
}

// QueryPacksClientCreateOrUpdateResponse contains the response from method QueryPacksClient.CreateOrUpdate.
type QueryPacksClientCreateOrUpdateResponse struct {
	LogAnalyticsQueryPack
}

// QueryPacksClientDeleteResponse contains the response from method QueryPacksClient.Delete.
type QueryPacksClientDeleteResponse struct {
	// placeholder for future response values
}

// QueryPacksClientGetResponse contains the response from method QueryPacksClient.Get.
type QueryPacksClientGetResponse struct {
	LogAnalyticsQueryPack
}

// QueryPacksClientListByResourceGroupResponse contains the response from method QueryPacksClient.ListByResourceGroup.
type QueryPacksClientListByResourceGroupResponse struct {
	LogAnalyticsQueryPackListResult
}

// QueryPacksClientListResponse contains the response from method QueryPacksClient.List.
type QueryPacksClientListResponse struct {
	LogAnalyticsQueryPackListResult
}

// QueryPacksClientUpdateTagsResponse contains the response from method QueryPacksClient.UpdateTags.
type QueryPacksClientUpdateTagsResponse struct {
	LogAnalyticsQueryPack
}

// SavedSearchesClientCreateOrUpdateResponse contains the response from method SavedSearchesClient.CreateOrUpdate.
type SavedSearchesClientCreateOrUpdateResponse struct {
	SavedSearch
}

// SavedSearchesClientDeleteResponse contains the response from method SavedSearchesClient.Delete.
type SavedSearchesClientDeleteResponse struct {
	// placeholder for future response values
}

// SavedSearchesClientGetResponse contains the response from method SavedSearchesClient.Get.
type SavedSearchesClientGetResponse struct {
	SavedSearch
}

// SavedSearchesClientListByWorkspaceResponse contains the response from method SavedSearchesClient.ListByWorkspace.
type SavedSearchesClientListByWorkspaceResponse struct {
	SavedSearchesListResult
}

// SchemaClientGetResponse contains the response from method SchemaClient.Get.
type SchemaClientGetResponse struct {
	SearchGetSchemaResponse
}

// SharedKeysClientGetSharedKeysResponse contains the response from method SharedKeysClient.GetSharedKeys.
type SharedKeysClientGetSharedKeysResponse struct {
	SharedKeys
}

// SharedKeysClientRegenerateResponse contains the response from method SharedKeysClient.Regenerate.
type SharedKeysClientRegenerateResponse struct {
	SharedKeys
}

// StorageInsightConfigsClientCreateOrUpdateResponse contains the response from method StorageInsightConfigsClient.CreateOrUpdate.
type StorageInsightConfigsClientCreateOrUpdateResponse struct {
	StorageInsight
}

// StorageInsightConfigsClientDeleteResponse contains the response from method StorageInsightConfigsClient.Delete.
type StorageInsightConfigsClientDeleteResponse struct {
	// placeholder for future response values
}

// StorageInsightConfigsClientGetResponse contains the response from method StorageInsightConfigsClient.Get.
type StorageInsightConfigsClientGetResponse struct {
	StorageInsight
}

// StorageInsightConfigsClientListByWorkspaceResponse contains the response from method StorageInsightConfigsClient.ListByWorkspace.
type StorageInsightConfigsClientListByWorkspaceResponse struct {
	StorageInsightListResult
}

// TablesClientCreateOrUpdateResponse contains the response from method TablesClient.CreateOrUpdate.
type TablesClientCreateOrUpdateResponse struct {
	Table
}

// TablesClientDeleteResponse contains the response from method TablesClient.Delete.
type TablesClientDeleteResponse struct {
	// placeholder for future response values
}

// TablesClientGetResponse contains the response from method TablesClient.Get.
type TablesClientGetResponse struct {
	Table
}

// TablesClientListByWorkspaceResponse contains the response from method TablesClient.ListByWorkspace.
type TablesClientListByWorkspaceResponse struct {
	TablesListResult
}

// TablesClientMigrateResponse contains the response from method TablesClient.Migrate.
type TablesClientMigrateResponse struct {
	// placeholder for future response values
}

// TablesClientUpdateResponse contains the response from method TablesClient.Update.
type TablesClientUpdateResponse struct {
	Table
}

// UsagesClientListResponse contains the response from method UsagesClient.List.
type UsagesClientListResponse struct {
	WorkspaceListUsagesResult
}

// WorkspacePurgeClientGetPurgeStatusResponse contains the response from method WorkspacePurgeClient.GetPurgeStatus.
type WorkspacePurgeClientGetPurgeStatusResponse struct {
	WorkspacePurgeStatusResponse
}

// WorkspacePurgeClientPurgeResponse contains the response from method WorkspacePurgeClient.Purge.
type WorkspacePurgeClientPurgeResponse struct {
	WorkspacePurgeResponse
	// XMSStatusLocation contains the information returned from the x-ms-status-location header response.
	XMSStatusLocation *string
}

// WorkspacesClientCreateOrUpdateResponse contains the response from method WorkspacesClient.CreateOrUpdate.
type WorkspacesClientCreateOrUpdateResponse struct {
	Workspace
}

// WorkspacesClientDeleteResponse contains the response from method WorkspacesClient.Delete.
type WorkspacesClientDeleteResponse struct {
	// placeholder for future response values
}

// WorkspacesClientGetResponse contains the response from method WorkspacesClient.Get.
type WorkspacesClientGetResponse struct {
	Workspace
}

// WorkspacesClientListByResourceGroupResponse contains the response from method WorkspacesClient.ListByResourceGroup.
type WorkspacesClientListByResourceGroupResponse struct {
	WorkspaceListResult
}

// WorkspacesClientListResponse contains the response from method WorkspacesClient.List.
type WorkspacesClientListResponse struct {
	WorkspaceListResult
}

// WorkspacesClientUpdateResponse contains the response from method WorkspacesClient.Update.
type WorkspacesClientUpdateResponse struct {
	Workspace
}
