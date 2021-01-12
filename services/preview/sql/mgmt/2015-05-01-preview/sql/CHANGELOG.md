Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82/specification/sql/resource-manager/readme.md tag: `package-2015-05-preview`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *DatabasesCreateImportOperationFuture.Result(DatabasesClient) (ImportExportResponse, error)
1. *DatabasesCreateOrUpdateFuture.Result(DatabasesClient) (Database, error)
1. *DatabasesExportFuture.Result(DatabasesClient) (ImportExportResponse, error)
1. *DatabasesImportFuture.Result(DatabasesClient) (ImportExportResponse, error)
1. *DatabasesPauseFuture.Result(DatabasesClient) (autorest.Response, error)
1. *DatabasesResumeFuture.Result(DatabasesClient) (autorest.Response, error)
1. *DatabasesUpdateFuture.Result(DatabasesClient) (Database, error)
1. *ElasticPoolsCreateOrUpdateFuture.Result(ElasticPoolsClient) (ElasticPool, error)
1. *ElasticPoolsUpdateFuture.Result(ElasticPoolsClient) (ElasticPool, error)
1. *EncryptionProtectorsCreateOrUpdateFuture.Result(EncryptionProtectorsClient) (EncryptionProtector, error)
1. *EncryptionProtectorsRevalidateFuture.Result(EncryptionProtectorsClient) (autorest.Response, error)
1. *FailoverGroupsCreateOrUpdateFuture.Result(FailoverGroupsClient) (FailoverGroup, error)
1. *FailoverGroupsDeleteFuture.Result(FailoverGroupsClient) (autorest.Response, error)
1. *FailoverGroupsFailoverFuture.Result(FailoverGroupsClient) (FailoverGroup, error)
1. *FailoverGroupsForceFailoverAllowDataLossFuture.Result(FailoverGroupsClient) (FailoverGroup, error)
1. *FailoverGroupsUpdateFuture.Result(FailoverGroupsClient) (FailoverGroup, error)
1. *ManagedInstancesCreateOrUpdateFuture.Result(ManagedInstancesClient) (ManagedInstance, error)
1. *ManagedInstancesDeleteFuture.Result(ManagedInstancesClient) (autorest.Response, error)
1. *ManagedInstancesUpdateFuture.Result(ManagedInstancesClient) (ManagedInstance, error)
1. *ReplicationLinksFailoverAllowDataLossFuture.Result(ReplicationLinksClient) (autorest.Response, error)
1. *ReplicationLinksFailoverFuture.Result(ReplicationLinksClient) (autorest.Response, error)
1. *ReplicationLinksUnlinkFuture.Result(ReplicationLinksClient) (autorest.Response, error)
1. *ServerAzureADAdministratorsCreateOrUpdateFuture.Result(ServerAzureADAdministratorsClient) (ServerAzureADAdministrator, error)
1. *ServerAzureADAdministratorsDeleteFuture.Result(ServerAzureADAdministratorsClient) (ServerAzureADAdministrator, error)
1. *ServerCommunicationLinksCreateOrUpdateFuture.Result(ServerCommunicationLinksClient) (ServerCommunicationLink, error)
1. *ServerKeysCreateOrUpdateFuture.Result(ServerKeysClient) (ServerKey, error)
1. *ServerKeysDeleteFuture.Result(ServerKeysClient) (autorest.Response, error)
1. *ServersCreateOrUpdateFuture.Result(ServersClient) (Server, error)
1. *ServersDeleteFuture.Result(ServersClient) (autorest.Response, error)
1. *ServersUpdateFuture.Result(ServersClient) (Server, error)
1. *SyncAgentsCreateOrUpdateFuture.Result(SyncAgentsClient) (SyncAgent, error)
1. *SyncAgentsDeleteFuture.Result(SyncAgentsClient) (autorest.Response, error)
1. *SyncGroupsCreateOrUpdateFuture.Result(SyncGroupsClient) (SyncGroup, error)
1. *SyncGroupsDeleteFuture.Result(SyncGroupsClient) (autorest.Response, error)
1. *SyncGroupsRefreshHubSchemaFuture.Result(SyncGroupsClient) (autorest.Response, error)
1. *SyncGroupsUpdateFuture.Result(SyncGroupsClient) (SyncGroup, error)
1. *SyncMembersCreateOrUpdateFuture.Result(SyncMembersClient) (SyncMember, error)
1. *SyncMembersDeleteFuture.Result(SyncMembersClient) (autorest.Response, error)
1. *SyncMembersRefreshMemberSchemaFuture.Result(SyncMembersClient) (autorest.Response, error)
1. *SyncMembersUpdateFuture.Result(SyncMembersClient) (SyncMember, error)
1. *VirtualClustersDeleteFuture.Result(VirtualClustersClient) (autorest.Response, error)
1. *VirtualClustersUpdateFuture.Result(VirtualClustersClient) (VirtualCluster, error)
1. *VirtualNetworkRulesCreateOrUpdateFuture.Result(VirtualNetworkRulesClient) (VirtualNetworkRule, error)
1. *VirtualNetworkRulesDeleteFuture.Result(VirtualNetworkRulesClient) (autorest.Response, error)

## Struct Changes

### Removed Struct Fields

1. DatabasesCreateImportOperationFuture.azure.Future
1. DatabasesCreateOrUpdateFuture.azure.Future
1. DatabasesExportFuture.azure.Future
1. DatabasesImportFuture.azure.Future
1. DatabasesPauseFuture.azure.Future
1. DatabasesResumeFuture.azure.Future
1. DatabasesUpdateFuture.azure.Future
1. ElasticPoolsCreateOrUpdateFuture.azure.Future
1. ElasticPoolsUpdateFuture.azure.Future
1. EncryptionProtectorsCreateOrUpdateFuture.azure.Future
1. EncryptionProtectorsRevalidateFuture.azure.Future
1. FailoverGroupsCreateOrUpdateFuture.azure.Future
1. FailoverGroupsDeleteFuture.azure.Future
1. FailoverGroupsFailoverFuture.azure.Future
1. FailoverGroupsForceFailoverAllowDataLossFuture.azure.Future
1. FailoverGroupsUpdateFuture.azure.Future
1. ManagedInstancesCreateOrUpdateFuture.azure.Future
1. ManagedInstancesDeleteFuture.azure.Future
1. ManagedInstancesUpdateFuture.azure.Future
1. ReplicationLinksFailoverAllowDataLossFuture.azure.Future
1. ReplicationLinksFailoverFuture.azure.Future
1. ReplicationLinksUnlinkFuture.azure.Future
1. ServerAzureADAdministratorsCreateOrUpdateFuture.azure.Future
1. ServerAzureADAdministratorsDeleteFuture.azure.Future
1. ServerCommunicationLinksCreateOrUpdateFuture.azure.Future
1. ServerKeysCreateOrUpdateFuture.azure.Future
1. ServerKeysDeleteFuture.azure.Future
1. ServersCreateOrUpdateFuture.azure.Future
1. ServersDeleteFuture.azure.Future
1. ServersUpdateFuture.azure.Future
1. SyncAgentsCreateOrUpdateFuture.azure.Future
1. SyncAgentsDeleteFuture.azure.Future
1. SyncGroupsCreateOrUpdateFuture.azure.Future
1. SyncGroupsDeleteFuture.azure.Future
1. SyncGroupsRefreshHubSchemaFuture.azure.Future
1. SyncGroupsUpdateFuture.azure.Future
1. SyncMembersCreateOrUpdateFuture.azure.Future
1. SyncMembersDeleteFuture.azure.Future
1. SyncMembersRefreshMemberSchemaFuture.azure.Future
1. SyncMembersUpdateFuture.azure.Future
1. VirtualClustersDeleteFuture.azure.Future
1. VirtualClustersUpdateFuture.azure.Future
1. VirtualNetworkRulesCreateOrUpdateFuture.azure.Future
1. VirtualNetworkRulesDeleteFuture.azure.Future

## Struct Changes

### New Struct Fields

1. DatabasesCreateImportOperationFuture.Result
1. DatabasesCreateImportOperationFuture.azure.FutureAPI
1. DatabasesCreateOrUpdateFuture.Result
1. DatabasesCreateOrUpdateFuture.azure.FutureAPI
1. DatabasesExportFuture.Result
1. DatabasesExportFuture.azure.FutureAPI
1. DatabasesImportFuture.Result
1. DatabasesImportFuture.azure.FutureAPI
1. DatabasesPauseFuture.Result
1. DatabasesPauseFuture.azure.FutureAPI
1. DatabasesResumeFuture.Result
1. DatabasesResumeFuture.azure.FutureAPI
1. DatabasesUpdateFuture.Result
1. DatabasesUpdateFuture.azure.FutureAPI
1. ElasticPoolsCreateOrUpdateFuture.Result
1. ElasticPoolsCreateOrUpdateFuture.azure.FutureAPI
1. ElasticPoolsUpdateFuture.Result
1. ElasticPoolsUpdateFuture.azure.FutureAPI
1. EncryptionProtectorsCreateOrUpdateFuture.Result
1. EncryptionProtectorsCreateOrUpdateFuture.azure.FutureAPI
1. EncryptionProtectorsRevalidateFuture.Result
1. EncryptionProtectorsRevalidateFuture.azure.FutureAPI
1. FailoverGroupsCreateOrUpdateFuture.Result
1. FailoverGroupsCreateOrUpdateFuture.azure.FutureAPI
1. FailoverGroupsDeleteFuture.Result
1. FailoverGroupsDeleteFuture.azure.FutureAPI
1. FailoverGroupsFailoverFuture.Result
1. FailoverGroupsFailoverFuture.azure.FutureAPI
1. FailoverGroupsForceFailoverAllowDataLossFuture.Result
1. FailoverGroupsForceFailoverAllowDataLossFuture.azure.FutureAPI
1. FailoverGroupsUpdateFuture.Result
1. FailoverGroupsUpdateFuture.azure.FutureAPI
1. ManagedInstancesCreateOrUpdateFuture.Result
1. ManagedInstancesCreateOrUpdateFuture.azure.FutureAPI
1. ManagedInstancesDeleteFuture.Result
1. ManagedInstancesDeleteFuture.azure.FutureAPI
1. ManagedInstancesUpdateFuture.Result
1. ManagedInstancesUpdateFuture.azure.FutureAPI
1. ReplicationLinksFailoverAllowDataLossFuture.Result
1. ReplicationLinksFailoverAllowDataLossFuture.azure.FutureAPI
1. ReplicationLinksFailoverFuture.Result
1. ReplicationLinksFailoverFuture.azure.FutureAPI
1. ReplicationLinksUnlinkFuture.Result
1. ReplicationLinksUnlinkFuture.azure.FutureAPI
1. ServerAzureADAdministratorsCreateOrUpdateFuture.Result
1. ServerAzureADAdministratorsCreateOrUpdateFuture.azure.FutureAPI
1. ServerAzureADAdministratorsDeleteFuture.Result
1. ServerAzureADAdministratorsDeleteFuture.azure.FutureAPI
1. ServerCommunicationLinksCreateOrUpdateFuture.Result
1. ServerCommunicationLinksCreateOrUpdateFuture.azure.FutureAPI
1. ServerKeysCreateOrUpdateFuture.Result
1. ServerKeysCreateOrUpdateFuture.azure.FutureAPI
1. ServerKeysDeleteFuture.Result
1. ServerKeysDeleteFuture.azure.FutureAPI
1. ServersCreateOrUpdateFuture.Result
1. ServersCreateOrUpdateFuture.azure.FutureAPI
1. ServersDeleteFuture.Result
1. ServersDeleteFuture.azure.FutureAPI
1. ServersUpdateFuture.Result
1. ServersUpdateFuture.azure.FutureAPI
1. SyncAgentsCreateOrUpdateFuture.Result
1. SyncAgentsCreateOrUpdateFuture.azure.FutureAPI
1. SyncAgentsDeleteFuture.Result
1. SyncAgentsDeleteFuture.azure.FutureAPI
1. SyncGroupsCreateOrUpdateFuture.Result
1. SyncGroupsCreateOrUpdateFuture.azure.FutureAPI
1. SyncGroupsDeleteFuture.Result
1. SyncGroupsDeleteFuture.azure.FutureAPI
1. SyncGroupsRefreshHubSchemaFuture.Result
1. SyncGroupsRefreshHubSchemaFuture.azure.FutureAPI
1. SyncGroupsUpdateFuture.Result
1. SyncGroupsUpdateFuture.azure.FutureAPI
1. SyncMembersCreateOrUpdateFuture.Result
1. SyncMembersCreateOrUpdateFuture.azure.FutureAPI
1. SyncMembersDeleteFuture.Result
1. SyncMembersDeleteFuture.azure.FutureAPI
1. SyncMembersRefreshMemberSchemaFuture.Result
1. SyncMembersRefreshMemberSchemaFuture.azure.FutureAPI
1. SyncMembersUpdateFuture.Result
1. SyncMembersUpdateFuture.azure.FutureAPI
1. VirtualClustersDeleteFuture.Result
1. VirtualClustersDeleteFuture.azure.FutureAPI
1. VirtualClustersUpdateFuture.Result
1. VirtualClustersUpdateFuture.azure.FutureAPI
1. VirtualNetworkRulesCreateOrUpdateFuture.Result
1. VirtualNetworkRulesCreateOrUpdateFuture.azure.FutureAPI
1. VirtualNetworkRulesDeleteFuture.Result
1. VirtualNetworkRulesDeleteFuture.azure.FutureAPI
