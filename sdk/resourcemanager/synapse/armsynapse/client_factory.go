//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsynapse

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

// ClientFactory is a client factory used to create any client in this module.
// Don't use this type directly, use NewClientFactory instead.
type ClientFactory struct {
	subscriptionID string
	credential azcore.TokenCredential
	options *arm.ClientOptions
}

// NewClientFactory creates a new instance of ClientFactory with the specified values.
// The parameter values will be propagated to any client created from this factory.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewClientFactory(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClientFactory, error) {
	_, err := arm.NewClient(moduleName+".ClientFactory", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	return &ClientFactory{
		subscriptionID: 	subscriptionID,		credential: credential,
		options: options.Clone(),
	}, nil
}

func (c *ClientFactory) NewAzureADOnlyAuthenticationsClient() *AzureADOnlyAuthenticationsClient {
	subClient, _ := NewAzureADOnlyAuthenticationsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewBigDataPoolsClient() *BigDataPoolsClient {
	subClient, _ := NewBigDataPoolsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewDataMaskingPoliciesClient() *DataMaskingPoliciesClient {
	subClient, _ := NewDataMaskingPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewDataMaskingRulesClient() *DataMaskingRulesClient {
	subClient, _ := NewDataMaskingRulesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewExtendedSQLPoolBlobAuditingPoliciesClient() *ExtendedSQLPoolBlobAuditingPoliciesClient {
	subClient, _ := NewExtendedSQLPoolBlobAuditingPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewGetClient() *GetClient {
	subClient, _ := NewGetClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewIPFirewallRulesClient() *IPFirewallRulesClient {
	subClient, _ := NewIPFirewallRulesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewIntegrationRuntimeAuthKeysClient() *IntegrationRuntimeAuthKeysClient {
	subClient, _ := NewIntegrationRuntimeAuthKeysClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewIntegrationRuntimeConnectionInfosClient() *IntegrationRuntimeConnectionInfosClient {
	subClient, _ := NewIntegrationRuntimeConnectionInfosClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewIntegrationRuntimeCredentialsClient() *IntegrationRuntimeCredentialsClient {
	subClient, _ := NewIntegrationRuntimeCredentialsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewIntegrationRuntimeMonitoringDataClient() *IntegrationRuntimeMonitoringDataClient {
	subClient, _ := NewIntegrationRuntimeMonitoringDataClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewIntegrationRuntimeNodeIPAddressClient() *IntegrationRuntimeNodeIPAddressClient {
	subClient, _ := NewIntegrationRuntimeNodeIPAddressClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewIntegrationRuntimeNodesClient() *IntegrationRuntimeNodesClient {
	subClient, _ := NewIntegrationRuntimeNodesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewIntegrationRuntimeObjectMetadataClient() *IntegrationRuntimeObjectMetadataClient {
	subClient, _ := NewIntegrationRuntimeObjectMetadataClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewIntegrationRuntimeStatusClient() *IntegrationRuntimeStatusClient {
	subClient, _ := NewIntegrationRuntimeStatusClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewIntegrationRuntimesClient() *IntegrationRuntimesClient {
	subClient, _ := NewIntegrationRuntimesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewKeysClient() *KeysClient {
	subClient, _ := NewKeysClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewKustoOperationsClient() *KustoOperationsClient {
	subClient, _ := NewKustoOperationsClient(c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewKustoPoolAttachedDatabaseConfigurationsClient() *KustoPoolAttachedDatabaseConfigurationsClient {
	subClient, _ := NewKustoPoolAttachedDatabaseConfigurationsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewKustoPoolChildResourceClient() *KustoPoolChildResourceClient {
	subClient, _ := NewKustoPoolChildResourceClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewKustoPoolDataConnectionsClient() *KustoPoolDataConnectionsClient {
	subClient, _ := NewKustoPoolDataConnectionsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewKustoPoolDatabasePrincipalAssignmentsClient() *KustoPoolDatabasePrincipalAssignmentsClient {
	subClient, _ := NewKustoPoolDatabasePrincipalAssignmentsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewKustoPoolDatabasesClient() *KustoPoolDatabasesClient {
	subClient, _ := NewKustoPoolDatabasesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewKustoPoolPrincipalAssignmentsClient() *KustoPoolPrincipalAssignmentsClient {
	subClient, _ := NewKustoPoolPrincipalAssignmentsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewKustoPoolPrivateLinkResourcesClient() *KustoPoolPrivateLinkResourcesClient {
	subClient, _ := NewKustoPoolPrivateLinkResourcesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewKustoPoolsClient() *KustoPoolsClient {
	subClient, _ := NewKustoPoolsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewLibrariesClient() *LibrariesClient {
	subClient, _ := NewLibrariesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewLibraryClient() *LibraryClient {
	subClient, _ := NewLibraryClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewOperationsClient() *OperationsClient {
	subClient, _ := NewOperationsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient {
	subClient, _ := NewPrivateEndpointConnectionsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewPrivateEndpointConnectionsPrivateLinkHubClient() *PrivateEndpointConnectionsPrivateLinkHubClient {
	subClient, _ := NewPrivateEndpointConnectionsPrivateLinkHubClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewPrivateLinkHubPrivateLinkResourcesClient() *PrivateLinkHubPrivateLinkResourcesClient {
	subClient, _ := NewPrivateLinkHubPrivateLinkResourcesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewPrivateLinkHubsClient() *PrivateLinkHubsClient {
	subClient, _ := NewPrivateLinkHubsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewPrivateLinkResourcesClient() *PrivateLinkResourcesClient {
	subClient, _ := NewPrivateLinkResourcesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewRestorableDroppedSQLPoolsClient() *RestorableDroppedSQLPoolsClient {
	subClient, _ := NewRestorableDroppedSQLPoolsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolBlobAuditingPoliciesClient() *SQLPoolBlobAuditingPoliciesClient {
	subClient, _ := NewSQLPoolBlobAuditingPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolColumnsClient() *SQLPoolColumnsClient {
	subClient, _ := NewSQLPoolColumnsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolConnectionPoliciesClient() *SQLPoolConnectionPoliciesClient {
	subClient, _ := NewSQLPoolConnectionPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolDataWarehouseUserActivitiesClient() *SQLPoolDataWarehouseUserActivitiesClient {
	subClient, _ := NewSQLPoolDataWarehouseUserActivitiesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolGeoBackupPoliciesClient() *SQLPoolGeoBackupPoliciesClient {
	subClient, _ := NewSQLPoolGeoBackupPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolMaintenanceWindowOptionsClient() *SQLPoolMaintenanceWindowOptionsClient {
	subClient, _ := NewSQLPoolMaintenanceWindowOptionsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolMaintenanceWindowsClient() *SQLPoolMaintenanceWindowsClient {
	subClient, _ := NewSQLPoolMaintenanceWindowsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolMetadataSyncConfigsClient() *SQLPoolMetadataSyncConfigsClient {
	subClient, _ := NewSQLPoolMetadataSyncConfigsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolOperationResultsClient() *SQLPoolOperationResultsClient {
	subClient, _ := NewSQLPoolOperationResultsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolOperationsClient() *SQLPoolOperationsClient {
	subClient, _ := NewSQLPoolOperationsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolRecommendedSensitivityLabelsClient() *SQLPoolRecommendedSensitivityLabelsClient {
	subClient, _ := NewSQLPoolRecommendedSensitivityLabelsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolReplicationLinksClient() *SQLPoolReplicationLinksClient {
	subClient, _ := NewSQLPoolReplicationLinksClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolRestorePointsClient() *SQLPoolRestorePointsClient {
	subClient, _ := NewSQLPoolRestorePointsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolSchemasClient() *SQLPoolSchemasClient {
	subClient, _ := NewSQLPoolSchemasClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolSecurityAlertPoliciesClient() *SQLPoolSecurityAlertPoliciesClient {
	subClient, _ := NewSQLPoolSecurityAlertPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolSensitivityLabelsClient() *SQLPoolSensitivityLabelsClient {
	subClient, _ := NewSQLPoolSensitivityLabelsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolTableColumnsClient() *SQLPoolTableColumnsClient {
	subClient, _ := NewSQLPoolTableColumnsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolTablesClient() *SQLPoolTablesClient {
	subClient, _ := NewSQLPoolTablesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolTransparentDataEncryptionsClient() *SQLPoolTransparentDataEncryptionsClient {
	subClient, _ := NewSQLPoolTransparentDataEncryptionsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolUsagesClient() *SQLPoolUsagesClient {
	subClient, _ := NewSQLPoolUsagesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolVulnerabilityAssessmentRuleBaselinesClient() *SQLPoolVulnerabilityAssessmentRuleBaselinesClient {
	subClient, _ := NewSQLPoolVulnerabilityAssessmentRuleBaselinesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolVulnerabilityAssessmentScansClient() *SQLPoolVulnerabilityAssessmentScansClient {
	subClient, _ := NewSQLPoolVulnerabilityAssessmentScansClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolVulnerabilityAssessmentsClient() *SQLPoolVulnerabilityAssessmentsClient {
	subClient, _ := NewSQLPoolVulnerabilityAssessmentsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolWorkloadClassifierClient() *SQLPoolWorkloadClassifierClient {
	subClient, _ := NewSQLPoolWorkloadClassifierClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolWorkloadGroupClient() *SQLPoolWorkloadGroupClient {
	subClient, _ := NewSQLPoolWorkloadGroupClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSQLPoolsClient() *SQLPoolsClient {
	subClient, _ := NewSQLPoolsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSparkConfigurationClient() *SparkConfigurationClient {
	subClient, _ := NewSparkConfigurationClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSparkConfigurationsClient() *SparkConfigurationsClient {
	subClient, _ := NewSparkConfigurationsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceAADAdminsClient() *WorkspaceAADAdminsClient {
	subClient, _ := NewWorkspaceAADAdminsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceManagedIdentitySQLControlSettingsClient() *WorkspaceManagedIdentitySQLControlSettingsClient {
	subClient, _ := NewWorkspaceManagedIdentitySQLControlSettingsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceManagedSQLServerBlobAuditingPoliciesClient() *WorkspaceManagedSQLServerBlobAuditingPoliciesClient {
	subClient, _ := NewWorkspaceManagedSQLServerBlobAuditingPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient() *WorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient {
	subClient, _ := NewWorkspaceManagedSQLServerDedicatedSQLMinimalTLSSettingsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceManagedSQLServerEncryptionProtectorClient() *WorkspaceManagedSQLServerEncryptionProtectorClient {
	subClient, _ := NewWorkspaceManagedSQLServerEncryptionProtectorClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceManagedSQLServerExtendedBlobAuditingPoliciesClient() *WorkspaceManagedSQLServerExtendedBlobAuditingPoliciesClient {
	subClient, _ := NewWorkspaceManagedSQLServerExtendedBlobAuditingPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceManagedSQLServerRecoverableSQLPoolsClient() *WorkspaceManagedSQLServerRecoverableSQLPoolsClient {
	subClient, _ := NewWorkspaceManagedSQLServerRecoverableSQLPoolsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceManagedSQLServerSecurityAlertPolicyClient() *WorkspaceManagedSQLServerSecurityAlertPolicyClient {
	subClient, _ := NewWorkspaceManagedSQLServerSecurityAlertPolicyClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceManagedSQLServerUsagesClient() *WorkspaceManagedSQLServerUsagesClient {
	subClient, _ := NewWorkspaceManagedSQLServerUsagesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceManagedSQLServerVulnerabilityAssessmentsClient() *WorkspaceManagedSQLServerVulnerabilityAssessmentsClient {
	subClient, _ := NewWorkspaceManagedSQLServerVulnerabilityAssessmentsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspaceSQLAADAdminsClient() *WorkspaceSQLAADAdminsClient {
	subClient, _ := NewWorkspaceSQLAADAdminsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewWorkspacesClient() *WorkspacesClient {
	subClient, _ := NewWorkspacesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

