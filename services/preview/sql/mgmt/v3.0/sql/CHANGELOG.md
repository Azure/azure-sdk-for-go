Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82/specification/sql/resource-manager/readme.md tag: `package-composite-v3`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *BackupLongTermRetentionPoliciesCreateOrUpdateFuture.Result(BackupLongTermRetentionPoliciesClient) (BackupLongTermRetentionPolicy, error)
1. *BackupShortTermRetentionPoliciesCreateOrUpdateFuture.Result(BackupShortTermRetentionPoliciesClient) (BackupShortTermRetentionPolicy, error)
1. *BackupShortTermRetentionPoliciesUpdateFuture.Result(BackupShortTermRetentionPoliciesClient) (BackupShortTermRetentionPolicy, error)
1. *DatabaseVulnerabilityAssessmentScansInitiateScanFuture.Result(DatabaseVulnerabilityAssessmentScansClient) (autorest.Response, error)
1. *DatabasesCreateOrUpdateFuture.Result(DatabasesClient) (Database, error)
1. *DatabasesDeleteFuture.Result(DatabasesClient) (autorest.Response, error)
1. *DatabasesExportFuture.Result(DatabasesClient) (ImportExportOperationResult, error)
1. *DatabasesFailoverFuture.Result(DatabasesClient) (autorest.Response, error)
1. *DatabasesPauseFuture.Result(DatabasesClient) (Database, error)
1. *DatabasesResumeFuture.Result(DatabasesClient) (Database, error)
1. *DatabasesUpdateFuture.Result(DatabasesClient) (Database, error)
1. *DatabasesUpgradeDataWarehouseFuture.Result(DatabasesClient) (autorest.Response, error)
1. *ElasticPoolsCreateOrUpdateFuture.Result(ElasticPoolsClient) (ElasticPool, error)
1. *ElasticPoolsDeleteFuture.Result(ElasticPoolsClient) (autorest.Response, error)
1. *ElasticPoolsFailoverFuture.Result(ElasticPoolsClient) (autorest.Response, error)
1. *ElasticPoolsUpdateFuture.Result(ElasticPoolsClient) (ElasticPool, error)
1. *EncryptionProtectorsCreateOrUpdateFuture.Result(EncryptionProtectorsClient) (EncryptionProtector, error)
1. *EncryptionProtectorsRevalidateFuture.Result(EncryptionProtectorsClient) (autorest.Response, error)
1. *ExtendedServerBlobAuditingPoliciesCreateOrUpdateFuture.Result(ExtendedServerBlobAuditingPoliciesClient) (ExtendedServerBlobAuditingPolicy, error)
1. *FailoverGroupsCreateOrUpdateFuture.Result(FailoverGroupsClient) (FailoverGroup, error)
1. *FailoverGroupsDeleteFuture.Result(FailoverGroupsClient) (autorest.Response, error)
1. *FailoverGroupsFailoverFuture.Result(FailoverGroupsClient) (FailoverGroup, error)
1. *FailoverGroupsForceFailoverAllowDataLossFuture.Result(FailoverGroupsClient) (FailoverGroup, error)
1. *FailoverGroupsUpdateFuture.Result(FailoverGroupsClient) (FailoverGroup, error)
1. *ImportExportImportFuture.Result(ImportExportClient) (ImportExportOperationResult, error)
1. *InstanceFailoverGroupsCreateOrUpdateFuture.Result(InstanceFailoverGroupsClient) (InstanceFailoverGroup, error)
1. *InstanceFailoverGroupsDeleteFuture.Result(InstanceFailoverGroupsClient) (autorest.Response, error)
1. *InstanceFailoverGroupsFailoverFuture.Result(InstanceFailoverGroupsClient) (InstanceFailoverGroup, error)
1. *InstanceFailoverGroupsForceFailoverAllowDataLossFuture.Result(InstanceFailoverGroupsClient) (InstanceFailoverGroup, error)
1. *InstancePoolsCreateOrUpdateFuture.Result(InstancePoolsClient) (InstancePool, error)
1. *InstancePoolsDeleteFuture.Result(InstancePoolsClient) (autorest.Response, error)
1. *InstancePoolsUpdateFuture.Result(InstancePoolsClient) (InstancePool, error)
1. *JobAgentsCreateOrUpdateFuture.Result(JobAgentsClient) (JobAgent, error)
1. *JobAgentsDeleteFuture.Result(JobAgentsClient) (autorest.Response, error)
1. *JobAgentsUpdateFuture.Result(JobAgentsClient) (JobAgent, error)
1. *JobExecutionsCreateFuture.Result(JobExecutionsClient) (JobExecution, error)
1. *JobExecutionsCreateOrUpdateFuture.Result(JobExecutionsClient) (JobExecution, error)
1. *LongTermRetentionBackupsDeleteByResourceGroupFuture.Result(LongTermRetentionBackupsClient) (autorest.Response, error)
1. *LongTermRetentionBackupsDeleteFuture.Result(LongTermRetentionBackupsClient) (autorest.Response, error)
1. *LongTermRetentionManagedInstanceBackupsDeleteByResourceGroupFuture.Result(LongTermRetentionManagedInstanceBackupsClient) (autorest.Response, error)
1. *LongTermRetentionManagedInstanceBackupsDeleteFuture.Result(LongTermRetentionManagedInstanceBackupsClient) (autorest.Response, error)
1. *ManagedBackupShortTermRetentionPoliciesCreateOrUpdateFuture.Result(ManagedBackupShortTermRetentionPoliciesClient) (ManagedBackupShortTermRetentionPolicy, error)
1. *ManagedBackupShortTermRetentionPoliciesUpdateFuture.Result(ManagedBackupShortTermRetentionPoliciesClient) (ManagedBackupShortTermRetentionPolicy, error)
1. *ManagedDatabaseVulnerabilityAssessmentScansInitiateScanFuture.Result(ManagedDatabaseVulnerabilityAssessmentScansClient) (autorest.Response, error)
1. *ManagedDatabasesCompleteRestoreFuture.Result(ManagedDatabasesClient) (autorest.Response, error)
1. *ManagedDatabasesCreateOrUpdateFuture.Result(ManagedDatabasesClient) (ManagedDatabase, error)
1. *ManagedDatabasesDeleteFuture.Result(ManagedDatabasesClient) (autorest.Response, error)
1. *ManagedDatabasesUpdateFuture.Result(ManagedDatabasesClient) (ManagedDatabase, error)
1. *ManagedInstanceAdministratorsCreateOrUpdateFuture.Result(ManagedInstanceAdministratorsClient) (ManagedInstanceAdministrator, error)
1. *ManagedInstanceAdministratorsDeleteFuture.Result(ManagedInstanceAdministratorsClient) (autorest.Response, error)
1. *ManagedInstanceAzureADOnlyAuthenticationsCreateOrUpdateFuture.Result(ManagedInstanceAzureADOnlyAuthenticationsClient) (ManagedInstanceAzureADOnlyAuthentication, error)
1. *ManagedInstanceAzureADOnlyAuthenticationsDeleteFuture.Result(ManagedInstanceAzureADOnlyAuthenticationsClient) (autorest.Response, error)
1. *ManagedInstanceEncryptionProtectorsCreateOrUpdateFuture.Result(ManagedInstanceEncryptionProtectorsClient) (ManagedInstanceEncryptionProtector, error)
1. *ManagedInstanceEncryptionProtectorsRevalidateFuture.Result(ManagedInstanceEncryptionProtectorsClient) (autorest.Response, error)
1. *ManagedInstanceKeysCreateOrUpdateFuture.Result(ManagedInstanceKeysClient) (ManagedInstanceKey, error)
1. *ManagedInstanceKeysDeleteFuture.Result(ManagedInstanceKeysClient) (autorest.Response, error)
1. *ManagedInstanceLongTermRetentionPoliciesCreateOrUpdateFuture.Result(ManagedInstanceLongTermRetentionPoliciesClient) (ManagedInstanceLongTermRetentionPolicy, error)
1. *ManagedInstanceTdeCertificatesCreateFuture.Result(ManagedInstanceTdeCertificatesClient) (autorest.Response, error)
1. *ManagedInstancesCreateOrUpdateFuture.Result(ManagedInstancesClient) (ManagedInstance, error)
1. *ManagedInstancesDeleteFuture.Result(ManagedInstancesClient) (autorest.Response, error)
1. *ManagedInstancesFailoverFuture.Result(ManagedInstancesClient) (autorest.Response, error)
1. *ManagedInstancesUpdateFuture.Result(ManagedInstancesClient) (ManagedInstance, error)
1. *ManagedRestorableDroppedDatabaseBackupShortTermRetentionPoliciesCreateOrUpdateFuture.Result(ManagedRestorableDroppedDatabaseBackupShortTermRetentionPoliciesClient) (ManagedBackupShortTermRetentionPolicy, error)
1. *ManagedRestorableDroppedDatabaseBackupShortTermRetentionPoliciesUpdateFuture.Result(ManagedRestorableDroppedDatabaseBackupShortTermRetentionPoliciesClient) (ManagedBackupShortTermRetentionPolicy, error)
1. *ManagedServerSecurityAlertPoliciesCreateOrUpdateFuture.Result(ManagedServerSecurityAlertPoliciesClient) (ManagedServerSecurityAlertPolicy, error)
1. *PrivateEndpointConnectionsCreateOrUpdateFuture.Result(PrivateEndpointConnectionsClient) (PrivateEndpointConnection, error)
1. *PrivateEndpointConnectionsDeleteFuture.Result(PrivateEndpointConnectionsClient) (autorest.Response, error)
1. *ReplicationLinksFailoverAllowDataLossFuture.Result(ReplicationLinksClient) (autorest.Response, error)
1. *ReplicationLinksFailoverFuture.Result(ReplicationLinksClient) (autorest.Response, error)
1. *ReplicationLinksUnlinkFuture.Result(ReplicationLinksClient) (autorest.Response, error)
1. *RestorePointsCreateFuture.Result(RestorePointsClient) (RestorePoint, error)
1. *ServerAzureADAdministratorsCreateOrUpdateFuture.Result(ServerAzureADAdministratorsClient) (ServerAzureADAdministrator, error)
1. *ServerAzureADAdministratorsDeleteFuture.Result(ServerAzureADAdministratorsClient) (autorest.Response, error)
1. *ServerAzureADOnlyAuthenticationsCreateOrUpdateFuture.Result(ServerAzureADOnlyAuthenticationsClient) (ServerAzureADOnlyAuthentication, error)
1. *ServerAzureADOnlyAuthenticationsDeleteFuture.Result(ServerAzureADOnlyAuthenticationsClient) (autorest.Response, error)
1. *ServerBlobAuditingPoliciesCreateOrUpdateFuture.Result(ServerBlobAuditingPoliciesClient) (ServerBlobAuditingPolicy, error)
1. *ServerCommunicationLinksCreateOrUpdateFuture.Result(ServerCommunicationLinksClient) (ServerCommunicationLink, error)
1. *ServerDNSAliasesAcquireFuture.Result(ServerDNSAliasesClient) (autorest.Response, error)
1. *ServerDNSAliasesCreateOrUpdateFuture.Result(ServerDNSAliasesClient) (ServerDNSAlias, error)
1. *ServerDNSAliasesDeleteFuture.Result(ServerDNSAliasesClient) (autorest.Response, error)
1. *ServerKeysCreateOrUpdateFuture.Result(ServerKeysClient) (ServerKey, error)
1. *ServerKeysDeleteFuture.Result(ServerKeysClient) (autorest.Response, error)
1. *ServerSecurityAlertPoliciesCreateOrUpdateFuture.Result(ServerSecurityAlertPoliciesClient) (ServerSecurityAlertPolicy, error)
1. *ServerTrustGroupsCreateOrUpdateFuture.Result(ServerTrustGroupsClient) (ServerTrustGroup, error)
1. *ServerTrustGroupsDeleteFuture.Result(ServerTrustGroupsClient) (autorest.Response, error)
1. *ServersCreateOrUpdateFuture.Result(ServersClient) (Server, error)
1. *ServersDeleteFuture.Result(ServersClient) (autorest.Response, error)
1. *ServersImportDatabaseFuture.Result(ServersClient) (ImportExportOperationResult, error)
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
1. *TdeCertificatesCreateFuture.Result(TdeCertificatesClient) (autorest.Response, error)
1. *VirtualClustersDeleteFuture.Result(VirtualClustersClient) (autorest.Response, error)
1. *VirtualClustersUpdateFuture.Result(VirtualClustersClient) (VirtualCluster, error)
1. *VirtualNetworkRulesCreateOrUpdateFuture.Result(VirtualNetworkRulesClient) (VirtualNetworkRule, error)
1. *VirtualNetworkRulesDeleteFuture.Result(VirtualNetworkRulesClient) (autorest.Response, error)
1. *WorkloadClassifiersCreateOrUpdateFuture.Result(WorkloadClassifiersClient) (WorkloadClassifier, error)
1. *WorkloadClassifiersDeleteFuture.Result(WorkloadClassifiersClient) (autorest.Response, error)
1. *WorkloadGroupsCreateOrUpdateFuture.Result(WorkloadGroupsClient) (WorkloadGroup, error)
1. *WorkloadGroupsDeleteFuture.Result(WorkloadGroupsClient) (autorest.Response, error)

## Struct Changes

### Removed Struct Fields

1. BackupLongTermRetentionPoliciesCreateOrUpdateFuture.azure.Future
1. BackupShortTermRetentionPoliciesCreateOrUpdateFuture.azure.Future
1. BackupShortTermRetentionPoliciesUpdateFuture.azure.Future
1. DatabaseVulnerabilityAssessmentScansInitiateScanFuture.azure.Future
1. DatabasesCreateOrUpdateFuture.azure.Future
1. DatabasesDeleteFuture.azure.Future
1. DatabasesExportFuture.azure.Future
1. DatabasesFailoverFuture.azure.Future
1. DatabasesPauseFuture.azure.Future
1. DatabasesResumeFuture.azure.Future
1. DatabasesUpdateFuture.azure.Future
1. DatabasesUpgradeDataWarehouseFuture.azure.Future
1. ElasticPoolsCreateOrUpdateFuture.azure.Future
1. ElasticPoolsDeleteFuture.azure.Future
1. ElasticPoolsFailoverFuture.azure.Future
1. ElasticPoolsUpdateFuture.azure.Future
1. EncryptionProtectorsCreateOrUpdateFuture.azure.Future
1. EncryptionProtectorsRevalidateFuture.azure.Future
1. ExtendedServerBlobAuditingPoliciesCreateOrUpdateFuture.azure.Future
1. FailoverGroupsCreateOrUpdateFuture.azure.Future
1. FailoverGroupsDeleteFuture.azure.Future
1. FailoverGroupsFailoverFuture.azure.Future
1. FailoverGroupsForceFailoverAllowDataLossFuture.azure.Future
1. FailoverGroupsUpdateFuture.azure.Future
1. ImportExportImportFuture.azure.Future
1. InstanceFailoverGroupsCreateOrUpdateFuture.azure.Future
1. InstanceFailoverGroupsDeleteFuture.azure.Future
1. InstanceFailoverGroupsFailoverFuture.azure.Future
1. InstanceFailoverGroupsForceFailoverAllowDataLossFuture.azure.Future
1. InstancePoolsCreateOrUpdateFuture.azure.Future
1. InstancePoolsDeleteFuture.azure.Future
1. InstancePoolsUpdateFuture.azure.Future
1. JobAgentsCreateOrUpdateFuture.azure.Future
1. JobAgentsDeleteFuture.azure.Future
1. JobAgentsUpdateFuture.azure.Future
1. JobExecutionsCreateFuture.azure.Future
1. JobExecutionsCreateOrUpdateFuture.azure.Future
1. LongTermRetentionBackupsDeleteByResourceGroupFuture.azure.Future
1. LongTermRetentionBackupsDeleteFuture.azure.Future
1. LongTermRetentionManagedInstanceBackupsDeleteByResourceGroupFuture.azure.Future
1. LongTermRetentionManagedInstanceBackupsDeleteFuture.azure.Future
1. ManagedBackupShortTermRetentionPoliciesCreateOrUpdateFuture.azure.Future
1. ManagedBackupShortTermRetentionPoliciesUpdateFuture.azure.Future
1. ManagedDatabaseVulnerabilityAssessmentScansInitiateScanFuture.azure.Future
1. ManagedDatabasesCompleteRestoreFuture.azure.Future
1. ManagedDatabasesCreateOrUpdateFuture.azure.Future
1. ManagedDatabasesDeleteFuture.azure.Future
1. ManagedDatabasesUpdateFuture.azure.Future
1. ManagedInstanceAdministratorsCreateOrUpdateFuture.azure.Future
1. ManagedInstanceAdministratorsDeleteFuture.azure.Future
1. ManagedInstanceAzureADOnlyAuthenticationsCreateOrUpdateFuture.azure.Future
1. ManagedInstanceAzureADOnlyAuthenticationsDeleteFuture.azure.Future
1. ManagedInstanceEncryptionProtectorsCreateOrUpdateFuture.azure.Future
1. ManagedInstanceEncryptionProtectorsRevalidateFuture.azure.Future
1. ManagedInstanceKeysCreateOrUpdateFuture.azure.Future
1. ManagedInstanceKeysDeleteFuture.azure.Future
1. ManagedInstanceLongTermRetentionPoliciesCreateOrUpdateFuture.azure.Future
1. ManagedInstanceTdeCertificatesCreateFuture.azure.Future
1. ManagedInstancesCreateOrUpdateFuture.azure.Future
1. ManagedInstancesDeleteFuture.azure.Future
1. ManagedInstancesFailoverFuture.azure.Future
1. ManagedInstancesUpdateFuture.azure.Future
1. ManagedRestorableDroppedDatabaseBackupShortTermRetentionPoliciesCreateOrUpdateFuture.azure.Future
1. ManagedRestorableDroppedDatabaseBackupShortTermRetentionPoliciesUpdateFuture.azure.Future
1. ManagedServerSecurityAlertPoliciesCreateOrUpdateFuture.azure.Future
1. PrivateEndpointConnectionsCreateOrUpdateFuture.azure.Future
1. PrivateEndpointConnectionsDeleteFuture.azure.Future
1. ReplicationLinksFailoverAllowDataLossFuture.azure.Future
1. ReplicationLinksFailoverFuture.azure.Future
1. ReplicationLinksUnlinkFuture.azure.Future
1. RestorePointsCreateFuture.azure.Future
1. ServerAzureADAdministratorsCreateOrUpdateFuture.azure.Future
1. ServerAzureADAdministratorsDeleteFuture.azure.Future
1. ServerAzureADOnlyAuthenticationsCreateOrUpdateFuture.azure.Future
1. ServerAzureADOnlyAuthenticationsDeleteFuture.azure.Future
1. ServerBlobAuditingPoliciesCreateOrUpdateFuture.azure.Future
1. ServerCommunicationLinksCreateOrUpdateFuture.azure.Future
1. ServerDNSAliasesAcquireFuture.azure.Future
1. ServerDNSAliasesCreateOrUpdateFuture.azure.Future
1. ServerDNSAliasesDeleteFuture.azure.Future
1. ServerKeysCreateOrUpdateFuture.azure.Future
1. ServerKeysDeleteFuture.azure.Future
1. ServerSecurityAlertPoliciesCreateOrUpdateFuture.azure.Future
1. ServerTrustGroupsCreateOrUpdateFuture.azure.Future
1. ServerTrustGroupsDeleteFuture.azure.Future
1. ServersCreateOrUpdateFuture.azure.Future
1. ServersDeleteFuture.azure.Future
1. ServersImportDatabaseFuture.azure.Future
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
1. TdeCertificatesCreateFuture.azure.Future
1. VirtualClustersDeleteFuture.azure.Future
1. VirtualClustersUpdateFuture.azure.Future
1. VirtualNetworkRulesCreateOrUpdateFuture.azure.Future
1. VirtualNetworkRulesDeleteFuture.azure.Future
1. WorkloadClassifiersCreateOrUpdateFuture.azure.Future
1. WorkloadClassifiersDeleteFuture.azure.Future
1. WorkloadGroupsCreateOrUpdateFuture.azure.Future
1. WorkloadGroupsDeleteFuture.azure.Future

## Struct Changes

### New Struct Fields

1. BackupLongTermRetentionPoliciesCreateOrUpdateFuture.Result
1. BackupLongTermRetentionPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. BackupShortTermRetentionPoliciesCreateOrUpdateFuture.Result
1. BackupShortTermRetentionPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. BackupShortTermRetentionPoliciesUpdateFuture.Result
1. BackupShortTermRetentionPoliciesUpdateFuture.azure.FutureAPI
1. DatabaseVulnerabilityAssessmentScansInitiateScanFuture.Result
1. DatabaseVulnerabilityAssessmentScansInitiateScanFuture.azure.FutureAPI
1. DatabasesCreateOrUpdateFuture.Result
1. DatabasesCreateOrUpdateFuture.azure.FutureAPI
1. DatabasesDeleteFuture.Result
1. DatabasesDeleteFuture.azure.FutureAPI
1. DatabasesExportFuture.Result
1. DatabasesExportFuture.azure.FutureAPI
1. DatabasesFailoverFuture.Result
1. DatabasesFailoverFuture.azure.FutureAPI
1. DatabasesPauseFuture.Result
1. DatabasesPauseFuture.azure.FutureAPI
1. DatabasesResumeFuture.Result
1. DatabasesResumeFuture.azure.FutureAPI
1. DatabasesUpdateFuture.Result
1. DatabasesUpdateFuture.azure.FutureAPI
1. DatabasesUpgradeDataWarehouseFuture.Result
1. DatabasesUpgradeDataWarehouseFuture.azure.FutureAPI
1. ElasticPoolsCreateOrUpdateFuture.Result
1. ElasticPoolsCreateOrUpdateFuture.azure.FutureAPI
1. ElasticPoolsDeleteFuture.Result
1. ElasticPoolsDeleteFuture.azure.FutureAPI
1. ElasticPoolsFailoverFuture.Result
1. ElasticPoolsFailoverFuture.azure.FutureAPI
1. ElasticPoolsUpdateFuture.Result
1. ElasticPoolsUpdateFuture.azure.FutureAPI
1. EncryptionProtectorsCreateOrUpdateFuture.Result
1. EncryptionProtectorsCreateOrUpdateFuture.azure.FutureAPI
1. EncryptionProtectorsRevalidateFuture.Result
1. EncryptionProtectorsRevalidateFuture.azure.FutureAPI
1. ExtendedServerBlobAuditingPoliciesCreateOrUpdateFuture.Result
1. ExtendedServerBlobAuditingPoliciesCreateOrUpdateFuture.azure.FutureAPI
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
1. ImportExportImportFuture.Result
1. ImportExportImportFuture.azure.FutureAPI
1. InstanceFailoverGroupsCreateOrUpdateFuture.Result
1. InstanceFailoverGroupsCreateOrUpdateFuture.azure.FutureAPI
1. InstanceFailoverGroupsDeleteFuture.Result
1. InstanceFailoverGroupsDeleteFuture.azure.FutureAPI
1. InstanceFailoverGroupsFailoverFuture.Result
1. InstanceFailoverGroupsFailoverFuture.azure.FutureAPI
1. InstanceFailoverGroupsForceFailoverAllowDataLossFuture.Result
1. InstanceFailoverGroupsForceFailoverAllowDataLossFuture.azure.FutureAPI
1. InstancePoolsCreateOrUpdateFuture.Result
1. InstancePoolsCreateOrUpdateFuture.azure.FutureAPI
1. InstancePoolsDeleteFuture.Result
1. InstancePoolsDeleteFuture.azure.FutureAPI
1. InstancePoolsUpdateFuture.Result
1. InstancePoolsUpdateFuture.azure.FutureAPI
1. JobAgentsCreateOrUpdateFuture.Result
1. JobAgentsCreateOrUpdateFuture.azure.FutureAPI
1. JobAgentsDeleteFuture.Result
1. JobAgentsDeleteFuture.azure.FutureAPI
1. JobAgentsUpdateFuture.Result
1. JobAgentsUpdateFuture.azure.FutureAPI
1. JobExecutionsCreateFuture.Result
1. JobExecutionsCreateFuture.azure.FutureAPI
1. JobExecutionsCreateOrUpdateFuture.Result
1. JobExecutionsCreateOrUpdateFuture.azure.FutureAPI
1. LongTermRetentionBackupsDeleteByResourceGroupFuture.Result
1. LongTermRetentionBackupsDeleteByResourceGroupFuture.azure.FutureAPI
1. LongTermRetentionBackupsDeleteFuture.Result
1. LongTermRetentionBackupsDeleteFuture.azure.FutureAPI
1. LongTermRetentionManagedInstanceBackupsDeleteByResourceGroupFuture.Result
1. LongTermRetentionManagedInstanceBackupsDeleteByResourceGroupFuture.azure.FutureAPI
1. LongTermRetentionManagedInstanceBackupsDeleteFuture.Result
1. LongTermRetentionManagedInstanceBackupsDeleteFuture.azure.FutureAPI
1. ManagedBackupShortTermRetentionPoliciesCreateOrUpdateFuture.Result
1. ManagedBackupShortTermRetentionPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. ManagedBackupShortTermRetentionPoliciesUpdateFuture.Result
1. ManagedBackupShortTermRetentionPoliciesUpdateFuture.azure.FutureAPI
1. ManagedDatabaseVulnerabilityAssessmentScansInitiateScanFuture.Result
1. ManagedDatabaseVulnerabilityAssessmentScansInitiateScanFuture.azure.FutureAPI
1. ManagedDatabasesCompleteRestoreFuture.Result
1. ManagedDatabasesCompleteRestoreFuture.azure.FutureAPI
1. ManagedDatabasesCreateOrUpdateFuture.Result
1. ManagedDatabasesCreateOrUpdateFuture.azure.FutureAPI
1. ManagedDatabasesDeleteFuture.Result
1. ManagedDatabasesDeleteFuture.azure.FutureAPI
1. ManagedDatabasesUpdateFuture.Result
1. ManagedDatabasesUpdateFuture.azure.FutureAPI
1. ManagedInstanceAdministratorsCreateOrUpdateFuture.Result
1. ManagedInstanceAdministratorsCreateOrUpdateFuture.azure.FutureAPI
1. ManagedInstanceAdministratorsDeleteFuture.Result
1. ManagedInstanceAdministratorsDeleteFuture.azure.FutureAPI
1. ManagedInstanceAzureADOnlyAuthenticationsCreateOrUpdateFuture.Result
1. ManagedInstanceAzureADOnlyAuthenticationsCreateOrUpdateFuture.azure.FutureAPI
1. ManagedInstanceAzureADOnlyAuthenticationsDeleteFuture.Result
1. ManagedInstanceAzureADOnlyAuthenticationsDeleteFuture.azure.FutureAPI
1. ManagedInstanceEncryptionProtectorsCreateOrUpdateFuture.Result
1. ManagedInstanceEncryptionProtectorsCreateOrUpdateFuture.azure.FutureAPI
1. ManagedInstanceEncryptionProtectorsRevalidateFuture.Result
1. ManagedInstanceEncryptionProtectorsRevalidateFuture.azure.FutureAPI
1. ManagedInstanceKeysCreateOrUpdateFuture.Result
1. ManagedInstanceKeysCreateOrUpdateFuture.azure.FutureAPI
1. ManagedInstanceKeysDeleteFuture.Result
1. ManagedInstanceKeysDeleteFuture.azure.FutureAPI
1. ManagedInstanceLongTermRetentionPoliciesCreateOrUpdateFuture.Result
1. ManagedInstanceLongTermRetentionPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. ManagedInstanceTdeCertificatesCreateFuture.Result
1. ManagedInstanceTdeCertificatesCreateFuture.azure.FutureAPI
1. ManagedInstancesCreateOrUpdateFuture.Result
1. ManagedInstancesCreateOrUpdateFuture.azure.FutureAPI
1. ManagedInstancesDeleteFuture.Result
1. ManagedInstancesDeleteFuture.azure.FutureAPI
1. ManagedInstancesFailoverFuture.Result
1. ManagedInstancesFailoverFuture.azure.FutureAPI
1. ManagedInstancesUpdateFuture.Result
1. ManagedInstancesUpdateFuture.azure.FutureAPI
1. ManagedRestorableDroppedDatabaseBackupShortTermRetentionPoliciesCreateOrUpdateFuture.Result
1. ManagedRestorableDroppedDatabaseBackupShortTermRetentionPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. ManagedRestorableDroppedDatabaseBackupShortTermRetentionPoliciesUpdateFuture.Result
1. ManagedRestorableDroppedDatabaseBackupShortTermRetentionPoliciesUpdateFuture.azure.FutureAPI
1. ManagedServerSecurityAlertPoliciesCreateOrUpdateFuture.Result
1. ManagedServerSecurityAlertPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. PrivateEndpointConnectionsCreateOrUpdateFuture.Result
1. PrivateEndpointConnectionsCreateOrUpdateFuture.azure.FutureAPI
1. PrivateEndpointConnectionsDeleteFuture.Result
1. PrivateEndpointConnectionsDeleteFuture.azure.FutureAPI
1. ReplicationLinksFailoverAllowDataLossFuture.Result
1. ReplicationLinksFailoverAllowDataLossFuture.azure.FutureAPI
1. ReplicationLinksFailoverFuture.Result
1. ReplicationLinksFailoverFuture.azure.FutureAPI
1. ReplicationLinksUnlinkFuture.Result
1. ReplicationLinksUnlinkFuture.azure.FutureAPI
1. RestorePointsCreateFuture.Result
1. RestorePointsCreateFuture.azure.FutureAPI
1. ServerAzureADAdministratorsCreateOrUpdateFuture.Result
1. ServerAzureADAdministratorsCreateOrUpdateFuture.azure.FutureAPI
1. ServerAzureADAdministratorsDeleteFuture.Result
1. ServerAzureADAdministratorsDeleteFuture.azure.FutureAPI
1. ServerAzureADOnlyAuthenticationsCreateOrUpdateFuture.Result
1. ServerAzureADOnlyAuthenticationsCreateOrUpdateFuture.azure.FutureAPI
1. ServerAzureADOnlyAuthenticationsDeleteFuture.Result
1. ServerAzureADOnlyAuthenticationsDeleteFuture.azure.FutureAPI
1. ServerBlobAuditingPoliciesCreateOrUpdateFuture.Result
1. ServerBlobAuditingPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. ServerCommunicationLinksCreateOrUpdateFuture.Result
1. ServerCommunicationLinksCreateOrUpdateFuture.azure.FutureAPI
1. ServerDNSAliasesAcquireFuture.Result
1. ServerDNSAliasesAcquireFuture.azure.FutureAPI
1. ServerDNSAliasesCreateOrUpdateFuture.Result
1. ServerDNSAliasesCreateOrUpdateFuture.azure.FutureAPI
1. ServerDNSAliasesDeleteFuture.Result
1. ServerDNSAliasesDeleteFuture.azure.FutureAPI
1. ServerKeysCreateOrUpdateFuture.Result
1. ServerKeysCreateOrUpdateFuture.azure.FutureAPI
1. ServerKeysDeleteFuture.Result
1. ServerKeysDeleteFuture.azure.FutureAPI
1. ServerSecurityAlertPoliciesCreateOrUpdateFuture.Result
1. ServerSecurityAlertPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. ServerTrustGroupsCreateOrUpdateFuture.Result
1. ServerTrustGroupsCreateOrUpdateFuture.azure.FutureAPI
1. ServerTrustGroupsDeleteFuture.Result
1. ServerTrustGroupsDeleteFuture.azure.FutureAPI
1. ServersCreateOrUpdateFuture.Result
1. ServersCreateOrUpdateFuture.azure.FutureAPI
1. ServersDeleteFuture.Result
1. ServersDeleteFuture.azure.FutureAPI
1. ServersImportDatabaseFuture.Result
1. ServersImportDatabaseFuture.azure.FutureAPI
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
1. TdeCertificatesCreateFuture.Result
1. TdeCertificatesCreateFuture.azure.FutureAPI
1. VirtualClustersDeleteFuture.Result
1. VirtualClustersDeleteFuture.azure.FutureAPI
1. VirtualClustersUpdateFuture.Result
1. VirtualClustersUpdateFuture.azure.FutureAPI
1. VirtualNetworkRulesCreateOrUpdateFuture.Result
1. VirtualNetworkRulesCreateOrUpdateFuture.azure.FutureAPI
1. VirtualNetworkRulesDeleteFuture.Result
1. VirtualNetworkRulesDeleteFuture.azure.FutureAPI
1. WorkloadClassifiersCreateOrUpdateFuture.Result
1. WorkloadClassifiersCreateOrUpdateFuture.azure.FutureAPI
1. WorkloadClassifiersDeleteFuture.Result
1. WorkloadClassifiersDeleteFuture.azure.FutureAPI
1. WorkloadGroupsCreateOrUpdateFuture.Result
1. WorkloadGroupsCreateOrUpdateFuture.azure.FutureAPI
1. WorkloadGroupsDeleteFuture.Result
1. WorkloadGroupsDeleteFuture.azure.FutureAPI
