Generated from https://github.com/Azure/azure-rest-api-specs/tree/3a3a9452f965a227ce43e6b545035b99dd175f23/specification/synapse/resource-manager/readme.md tag: `package-2019-06-01-preview`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *BigDataPoolsCreateOrUpdateFuture.Result(BigDataPoolsClient) (BigDataPoolResourceInfo, error)
1. *BigDataPoolsDeleteFuture.Result(BigDataPoolsClient) (SetObject, error)
1. *IPFirewallRulesCreateOrUpdateFuture.Result(IPFirewallRulesClient) (IPFirewallRuleInfo, error)
1. *IPFirewallRulesDeleteFuture.Result(IPFirewallRulesClient) (SetObject, error)
1. *IPFirewallRulesReplaceAllFuture.Result(IPFirewallRulesClient) (ReplaceAllFirewallRulesOperationResponse, error)
1. *IntegrationRuntimeObjectMetadataRefreshFuture.Result(IntegrationRuntimeObjectMetadataClient) (SsisObjectMetadataStatusResponse, error)
1. *IntegrationRuntimesCreateFuture.Result(IntegrationRuntimesClient) (IntegrationRuntimeResource, error)
1. *IntegrationRuntimesDeleteFuture.Result(IntegrationRuntimesClient) (autorest.Response, error)
1. *IntegrationRuntimesDisableInteractiveQueryFuture.Result(IntegrationRuntimesClient) (autorest.Response, error)
1. *IntegrationRuntimesEnableInteractiveQueryFuture.Result(IntegrationRuntimesClient) (autorest.Response, error)
1. *IntegrationRuntimesStartFuture.Result(IntegrationRuntimesClient) (IntegrationRuntimeStatusResponse, error)
1. *IntegrationRuntimesStopFuture.Result(IntegrationRuntimesClient) (autorest.Response, error)
1. *PrivateEndpointConnectionsCreateFuture.Result(PrivateEndpointConnectionsClient) (PrivateEndpointConnection, error)
1. *PrivateEndpointConnectionsDeleteFuture.Result(PrivateEndpointConnectionsClient) (OperationResource, error)
1. *PrivateLinkHubsDeleteFuture.Result(PrivateLinkHubsClient) (autorest.Response, error)
1. *SQLPoolRestorePointsCreateFuture.Result(SQLPoolRestorePointsClient) (RestorePoint, error)
1. *SQLPoolVulnerabilityAssessmentScansInitiateScanFuture.Result(SQLPoolVulnerabilityAssessmentScansClient) (autorest.Response, error)
1. *SQLPoolWorkloadClassifierCreateOrUpdateFuture.Result(SQLPoolWorkloadClassifierClient) (WorkloadClassifier, error)
1. *SQLPoolWorkloadClassifierDeleteFuture.Result(SQLPoolWorkloadClassifierClient) (autorest.Response, error)
1. *SQLPoolWorkloadGroupCreateOrUpdateFuture.Result(SQLPoolWorkloadGroupClient) (WorkloadGroup, error)
1. *SQLPoolWorkloadGroupDeleteFuture.Result(SQLPoolWorkloadGroupClient) (autorest.Response, error)
1. *SQLPoolsCreateFuture.Result(SQLPoolsClient) (SQLPool, error)
1. *SQLPoolsDeleteFuture.Result(SQLPoolsClient) (SetObject, error)
1. *SQLPoolsPauseFuture.Result(SQLPoolsClient) (SetObject, error)
1. *SQLPoolsResumeFuture.Result(SQLPoolsClient) (SetObject, error)
1. *WorkspaceAadAdminsCreateOrUpdateFuture.Result(WorkspaceAadAdminsClient) (WorkspaceAadAdminInfo, error)
1. *WorkspaceAadAdminsDeleteFuture.Result(WorkspaceAadAdminsClient) (autorest.Response, error)
1. *WorkspaceManagedIdentitySQLControlSettingsCreateOrUpdateFuture.Result(WorkspaceManagedIdentitySQLControlSettingsClient) (ManagedIdentitySQLControlSettingsModel, error)
1. *WorkspaceManagedSQLServerBlobAuditingPoliciesCreateOrUpdateFuture.Result(WorkspaceManagedSQLServerBlobAuditingPoliciesClient) (ServerBlobAuditingPolicy, error)
1. *WorkspaceManagedSQLServerExtendedBlobAuditingPoliciesCreateOrUpdateFuture.Result(WorkspaceManagedSQLServerExtendedBlobAuditingPoliciesClient) (ExtendedServerBlobAuditingPolicy, error)
1. *WorkspaceManagedSQLServerSecurityAlertPolicyCreateOrUpdateFuture.Result(WorkspaceManagedSQLServerSecurityAlertPolicyClient) (ServerSecurityAlertPolicy, error)
1. *WorkspaceSQLAadAdminsCreateOrUpdateFuture.Result(WorkspaceSQLAadAdminsClient) (WorkspaceAadAdminInfo, error)
1. *WorkspaceSQLAadAdminsDeleteFuture.Result(WorkspaceSQLAadAdminsClient) (autorest.Response, error)
1. *WorkspacesCreateOrUpdateFuture.Result(WorkspacesClient) (Workspace, error)
1. *WorkspacesDeleteFuture.Result(WorkspacesClient) (SetObject, error)
1. *WorkspacesUpdateFuture.Result(WorkspacesClient) (Workspace, error)

## Struct Changes

### Removed Struct Fields

1. BigDataPoolsCreateOrUpdateFuture.azure.Future
1. BigDataPoolsDeleteFuture.azure.Future
1. IPFirewallRulesCreateOrUpdateFuture.azure.Future
1. IPFirewallRulesDeleteFuture.azure.Future
1. IPFirewallRulesReplaceAllFuture.azure.Future
1. IntegrationRuntimeObjectMetadataRefreshFuture.azure.Future
1. IntegrationRuntimesCreateFuture.azure.Future
1. IntegrationRuntimesDeleteFuture.azure.Future
1. IntegrationRuntimesDisableInteractiveQueryFuture.azure.Future
1. IntegrationRuntimesEnableInteractiveQueryFuture.azure.Future
1. IntegrationRuntimesStartFuture.azure.Future
1. IntegrationRuntimesStopFuture.azure.Future
1. PrivateEndpointConnectionsCreateFuture.azure.Future
1. PrivateEndpointConnectionsDeleteFuture.azure.Future
1. PrivateLinkHubsDeleteFuture.azure.Future
1. SQLPoolRestorePointsCreateFuture.azure.Future
1. SQLPoolVulnerabilityAssessmentScansInitiateScanFuture.azure.Future
1. SQLPoolWorkloadClassifierCreateOrUpdateFuture.azure.Future
1. SQLPoolWorkloadClassifierDeleteFuture.azure.Future
1. SQLPoolWorkloadGroupCreateOrUpdateFuture.azure.Future
1. SQLPoolWorkloadGroupDeleteFuture.azure.Future
1. SQLPoolsCreateFuture.azure.Future
1. SQLPoolsDeleteFuture.azure.Future
1. SQLPoolsPauseFuture.azure.Future
1. SQLPoolsResumeFuture.azure.Future
1. WorkspaceAadAdminsCreateOrUpdateFuture.azure.Future
1. WorkspaceAadAdminsDeleteFuture.azure.Future
1. WorkspaceManagedIdentitySQLControlSettingsCreateOrUpdateFuture.azure.Future
1. WorkspaceManagedSQLServerBlobAuditingPoliciesCreateOrUpdateFuture.azure.Future
1. WorkspaceManagedSQLServerExtendedBlobAuditingPoliciesCreateOrUpdateFuture.azure.Future
1. WorkspaceManagedSQLServerSecurityAlertPolicyCreateOrUpdateFuture.azure.Future
1. WorkspaceSQLAadAdminsCreateOrUpdateFuture.azure.Future
1. WorkspaceSQLAadAdminsDeleteFuture.azure.Future
1. WorkspacesCreateOrUpdateFuture.azure.Future
1. WorkspacesDeleteFuture.azure.Future
1. WorkspacesUpdateFuture.azure.Future

## Struct Changes

### New Struct Fields

1. BigDataPoolsCreateOrUpdateFuture.Result
1. BigDataPoolsCreateOrUpdateFuture.azure.FutureAPI
1. BigDataPoolsDeleteFuture.Result
1. BigDataPoolsDeleteFuture.azure.FutureAPI
1. IPFirewallRulesCreateOrUpdateFuture.Result
1. IPFirewallRulesCreateOrUpdateFuture.azure.FutureAPI
1. IPFirewallRulesDeleteFuture.Result
1. IPFirewallRulesDeleteFuture.azure.FutureAPI
1. IPFirewallRulesReplaceAllFuture.Result
1. IPFirewallRulesReplaceAllFuture.azure.FutureAPI
1. IntegrationRuntimeObjectMetadataRefreshFuture.Result
1. IntegrationRuntimeObjectMetadataRefreshFuture.azure.FutureAPI
1. IntegrationRuntimesCreateFuture.Result
1. IntegrationRuntimesCreateFuture.azure.FutureAPI
1. IntegrationRuntimesDeleteFuture.Result
1. IntegrationRuntimesDeleteFuture.azure.FutureAPI
1. IntegrationRuntimesDisableInteractiveQueryFuture.Result
1. IntegrationRuntimesDisableInteractiveQueryFuture.azure.FutureAPI
1. IntegrationRuntimesEnableInteractiveQueryFuture.Result
1. IntegrationRuntimesEnableInteractiveQueryFuture.azure.FutureAPI
1. IntegrationRuntimesStartFuture.Result
1. IntegrationRuntimesStartFuture.azure.FutureAPI
1. IntegrationRuntimesStopFuture.Result
1. IntegrationRuntimesStopFuture.azure.FutureAPI
1. PrivateEndpointConnectionsCreateFuture.Result
1. PrivateEndpointConnectionsCreateFuture.azure.FutureAPI
1. PrivateEndpointConnectionsDeleteFuture.Result
1. PrivateEndpointConnectionsDeleteFuture.azure.FutureAPI
1. PrivateLinkHubsDeleteFuture.Result
1. PrivateLinkHubsDeleteFuture.azure.FutureAPI
1. SQLPoolRestorePointsCreateFuture.Result
1. SQLPoolRestorePointsCreateFuture.azure.FutureAPI
1. SQLPoolVulnerabilityAssessmentScansInitiateScanFuture.Result
1. SQLPoolVulnerabilityAssessmentScansInitiateScanFuture.azure.FutureAPI
1. SQLPoolWorkloadClassifierCreateOrUpdateFuture.Result
1. SQLPoolWorkloadClassifierCreateOrUpdateFuture.azure.FutureAPI
1. SQLPoolWorkloadClassifierDeleteFuture.Result
1. SQLPoolWorkloadClassifierDeleteFuture.azure.FutureAPI
1. SQLPoolWorkloadGroupCreateOrUpdateFuture.Result
1. SQLPoolWorkloadGroupCreateOrUpdateFuture.azure.FutureAPI
1. SQLPoolWorkloadGroupDeleteFuture.Result
1. SQLPoolWorkloadGroupDeleteFuture.azure.FutureAPI
1. SQLPoolsCreateFuture.Result
1. SQLPoolsCreateFuture.azure.FutureAPI
1. SQLPoolsDeleteFuture.Result
1. SQLPoolsDeleteFuture.azure.FutureAPI
1. SQLPoolsPauseFuture.Result
1. SQLPoolsPauseFuture.azure.FutureAPI
1. SQLPoolsResumeFuture.Result
1. SQLPoolsResumeFuture.azure.FutureAPI
1. WorkspaceAadAdminsCreateOrUpdateFuture.Result
1. WorkspaceAadAdminsCreateOrUpdateFuture.azure.FutureAPI
1. WorkspaceAadAdminsDeleteFuture.Result
1. WorkspaceAadAdminsDeleteFuture.azure.FutureAPI
1. WorkspaceManagedIdentitySQLControlSettingsCreateOrUpdateFuture.Result
1. WorkspaceManagedIdentitySQLControlSettingsCreateOrUpdateFuture.azure.FutureAPI
1. WorkspaceManagedSQLServerBlobAuditingPoliciesCreateOrUpdateFuture.Result
1. WorkspaceManagedSQLServerBlobAuditingPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. WorkspaceManagedSQLServerExtendedBlobAuditingPoliciesCreateOrUpdateFuture.Result
1. WorkspaceManagedSQLServerExtendedBlobAuditingPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. WorkspaceManagedSQLServerSecurityAlertPolicyCreateOrUpdateFuture.Result
1. WorkspaceManagedSQLServerSecurityAlertPolicyCreateOrUpdateFuture.azure.FutureAPI
1. WorkspaceSQLAadAdminsCreateOrUpdateFuture.Result
1. WorkspaceSQLAadAdminsCreateOrUpdateFuture.azure.FutureAPI
1. WorkspaceSQLAadAdminsDeleteFuture.Result
1. WorkspaceSQLAadAdminsDeleteFuture.azure.FutureAPI
1. WorkspacesCreateOrUpdateFuture.Result
1. WorkspacesCreateOrUpdateFuture.azure.FutureAPI
1. WorkspacesDeleteFuture.Result
1. WorkspacesDeleteFuture.azure.FutureAPI
1. WorkspacesUpdateFuture.Result
1. WorkspacesUpdateFuture.azure.FutureAPI
