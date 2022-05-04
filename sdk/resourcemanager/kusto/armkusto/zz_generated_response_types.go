//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armkusto

// AttachedDatabaseConfigurationsClientCheckNameAvailabilityResponse contains the response from method AttachedDatabaseConfigurationsClient.CheckNameAvailability.
type AttachedDatabaseConfigurationsClientCheckNameAvailabilityResponse struct {
	CheckNameResult
}

// AttachedDatabaseConfigurationsClientCreateOrUpdateResponse contains the response from method AttachedDatabaseConfigurationsClient.CreateOrUpdate.
type AttachedDatabaseConfigurationsClientCreateOrUpdateResponse struct {
	AttachedDatabaseConfiguration
}

// AttachedDatabaseConfigurationsClientDeleteResponse contains the response from method AttachedDatabaseConfigurationsClient.Delete.
type AttachedDatabaseConfigurationsClientDeleteResponse struct {
	// placeholder for future response values
}

// AttachedDatabaseConfigurationsClientGetResponse contains the response from method AttachedDatabaseConfigurationsClient.Get.
type AttachedDatabaseConfigurationsClientGetResponse struct {
	AttachedDatabaseConfiguration
}

// AttachedDatabaseConfigurationsClientListByClusterResponse contains the response from method AttachedDatabaseConfigurationsClient.ListByCluster.
type AttachedDatabaseConfigurationsClientListByClusterResponse struct {
	AttachedDatabaseConfigurationListResult
}

// ClusterPrincipalAssignmentsClientCheckNameAvailabilityResponse contains the response from method ClusterPrincipalAssignmentsClient.CheckNameAvailability.
type ClusterPrincipalAssignmentsClientCheckNameAvailabilityResponse struct {
	CheckNameResult
}

// ClusterPrincipalAssignmentsClientCreateOrUpdateResponse contains the response from method ClusterPrincipalAssignmentsClient.CreateOrUpdate.
type ClusterPrincipalAssignmentsClientCreateOrUpdateResponse struct {
	ClusterPrincipalAssignment
}

// ClusterPrincipalAssignmentsClientDeleteResponse contains the response from method ClusterPrincipalAssignmentsClient.Delete.
type ClusterPrincipalAssignmentsClientDeleteResponse struct {
	// placeholder for future response values
}

// ClusterPrincipalAssignmentsClientGetResponse contains the response from method ClusterPrincipalAssignmentsClient.Get.
type ClusterPrincipalAssignmentsClientGetResponse struct {
	ClusterPrincipalAssignment
}

// ClusterPrincipalAssignmentsClientListResponse contains the response from method ClusterPrincipalAssignmentsClient.List.
type ClusterPrincipalAssignmentsClientListResponse struct {
	ClusterPrincipalAssignmentListResult
}

// ClustersClientAddLanguageExtensionsResponse contains the response from method ClustersClient.AddLanguageExtensions.
type ClustersClientAddLanguageExtensionsResponse struct {
	// placeholder for future response values
}

// ClustersClientCheckNameAvailabilityResponse contains the response from method ClustersClient.CheckNameAvailability.
type ClustersClientCheckNameAvailabilityResponse struct {
	CheckNameResult
}

// ClustersClientCreateOrUpdateResponse contains the response from method ClustersClient.CreateOrUpdate.
type ClustersClientCreateOrUpdateResponse struct {
	Cluster
}

// ClustersClientDeleteResponse contains the response from method ClustersClient.Delete.
type ClustersClientDeleteResponse struct {
	// placeholder for future response values
}

// ClustersClientDetachFollowerDatabasesResponse contains the response from method ClustersClient.DetachFollowerDatabases.
type ClustersClientDetachFollowerDatabasesResponse struct {
	// placeholder for future response values
}

// ClustersClientDiagnoseVirtualNetworkResponse contains the response from method ClustersClient.DiagnoseVirtualNetwork.
type ClustersClientDiagnoseVirtualNetworkResponse struct {
	DiagnoseVirtualNetworkResult
}

// ClustersClientGetResponse contains the response from method ClustersClient.Get.
type ClustersClientGetResponse struct {
	Cluster
}

// ClustersClientListByResourceGroupResponse contains the response from method ClustersClient.ListByResourceGroup.
type ClustersClientListByResourceGroupResponse struct {
	ClusterListResult
}

// ClustersClientListFollowerDatabasesResponse contains the response from method ClustersClient.ListFollowerDatabases.
type ClustersClientListFollowerDatabasesResponse struct {
	FollowerDatabaseListResult
}

// ClustersClientListLanguageExtensionsResponse contains the response from method ClustersClient.ListLanguageExtensions.
type ClustersClientListLanguageExtensionsResponse struct {
	LanguageExtensionsList
}

// ClustersClientListOutboundNetworkDependenciesEndpointsResponse contains the response from method ClustersClient.ListOutboundNetworkDependenciesEndpoints.
type ClustersClientListOutboundNetworkDependenciesEndpointsResponse struct {
	OutboundNetworkDependenciesEndpointListResult
}

// ClustersClientListResponse contains the response from method ClustersClient.List.
type ClustersClientListResponse struct {
	ClusterListResult
}

// ClustersClientListSKUsByResourceResponse contains the response from method ClustersClient.ListSKUsByResource.
type ClustersClientListSKUsByResourceResponse struct {
	ListResourceSKUsResult
}

// ClustersClientListSKUsResponse contains the response from method ClustersClient.ListSKUs.
type ClustersClientListSKUsResponse struct {
	SKUDescriptionList
}

// ClustersClientRemoveLanguageExtensionsResponse contains the response from method ClustersClient.RemoveLanguageExtensions.
type ClustersClientRemoveLanguageExtensionsResponse struct {
	// placeholder for future response values
}

// ClustersClientStartResponse contains the response from method ClustersClient.Start.
type ClustersClientStartResponse struct {
	// placeholder for future response values
}

// ClustersClientStopResponse contains the response from method ClustersClient.Stop.
type ClustersClientStopResponse struct {
	// placeholder for future response values
}

// ClustersClientUpdateResponse contains the response from method ClustersClient.Update.
type ClustersClientUpdateResponse struct {
	Cluster
}

// DataConnectionsClientCheckNameAvailabilityResponse contains the response from method DataConnectionsClient.CheckNameAvailability.
type DataConnectionsClientCheckNameAvailabilityResponse struct {
	CheckNameResult
}

// DataConnectionsClientCreateOrUpdateResponse contains the response from method DataConnectionsClient.CreateOrUpdate.
type DataConnectionsClientCreateOrUpdateResponse struct {
	DataConnectionClassification
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DataConnectionsClientCreateOrUpdateResponse.
func (d *DataConnectionsClientCreateOrUpdateResponse) UnmarshalJSON(data []byte) error {
	res, err := unmarshalDataConnectionClassification(data)
	if err != nil {
		return err
	}
	d.DataConnectionClassification = res
	return nil
}

// DataConnectionsClientDataConnectionValidationResponse contains the response from method DataConnectionsClient.DataConnectionValidation.
type DataConnectionsClientDataConnectionValidationResponse struct {
	DataConnectionValidationListResult
}

// DataConnectionsClientDeleteResponse contains the response from method DataConnectionsClient.Delete.
type DataConnectionsClientDeleteResponse struct {
	// placeholder for future response values
}

// DataConnectionsClientGetResponse contains the response from method DataConnectionsClient.Get.
type DataConnectionsClientGetResponse struct {
	DataConnectionClassification
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DataConnectionsClientGetResponse.
func (d *DataConnectionsClientGetResponse) UnmarshalJSON(data []byte) error {
	res, err := unmarshalDataConnectionClassification(data)
	if err != nil {
		return err
	}
	d.DataConnectionClassification = res
	return nil
}

// DataConnectionsClientListByDatabaseResponse contains the response from method DataConnectionsClient.ListByDatabase.
type DataConnectionsClientListByDatabaseResponse struct {
	DataConnectionListResult
}

// DataConnectionsClientUpdateResponse contains the response from method DataConnectionsClient.Update.
type DataConnectionsClientUpdateResponse struct {
	DataConnectionClassification
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DataConnectionsClientUpdateResponse.
func (d *DataConnectionsClientUpdateResponse) UnmarshalJSON(data []byte) error {
	res, err := unmarshalDataConnectionClassification(data)
	if err != nil {
		return err
	}
	d.DataConnectionClassification = res
	return nil
}

// DatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse contains the response from method DatabasePrincipalAssignmentsClient.CheckNameAvailability.
type DatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse struct {
	CheckNameResult
}

// DatabasePrincipalAssignmentsClientCreateOrUpdateResponse contains the response from method DatabasePrincipalAssignmentsClient.CreateOrUpdate.
type DatabasePrincipalAssignmentsClientCreateOrUpdateResponse struct {
	DatabasePrincipalAssignment
}

// DatabasePrincipalAssignmentsClientDeleteResponse contains the response from method DatabasePrincipalAssignmentsClient.Delete.
type DatabasePrincipalAssignmentsClientDeleteResponse struct {
	// placeholder for future response values
}

// DatabasePrincipalAssignmentsClientGetResponse contains the response from method DatabasePrincipalAssignmentsClient.Get.
type DatabasePrincipalAssignmentsClientGetResponse struct {
	DatabasePrincipalAssignment
}

// DatabasePrincipalAssignmentsClientListResponse contains the response from method DatabasePrincipalAssignmentsClient.List.
type DatabasePrincipalAssignmentsClientListResponse struct {
	DatabasePrincipalAssignmentListResult
}

// DatabasesClientAddPrincipalsResponse contains the response from method DatabasesClient.AddPrincipals.
type DatabasesClientAddPrincipalsResponse struct {
	DatabasePrincipalListResult
}

// DatabasesClientCheckNameAvailabilityResponse contains the response from method DatabasesClient.CheckNameAvailability.
type DatabasesClientCheckNameAvailabilityResponse struct {
	CheckNameResult
}

// DatabasesClientCreateOrUpdateResponse contains the response from method DatabasesClient.CreateOrUpdate.
type DatabasesClientCreateOrUpdateResponse struct {
	DatabaseClassification
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DatabasesClientCreateOrUpdateResponse.
func (d *DatabasesClientCreateOrUpdateResponse) UnmarshalJSON(data []byte) error {
	res, err := unmarshalDatabaseClassification(data)
	if err != nil {
		return err
	}
	d.DatabaseClassification = res
	return nil
}

// DatabasesClientDeleteResponse contains the response from method DatabasesClient.Delete.
type DatabasesClientDeleteResponse struct {
	// placeholder for future response values
}

// DatabasesClientGetResponse contains the response from method DatabasesClient.Get.
type DatabasesClientGetResponse struct {
	DatabaseClassification
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DatabasesClientGetResponse.
func (d *DatabasesClientGetResponse) UnmarshalJSON(data []byte) error {
	res, err := unmarshalDatabaseClassification(data)
	if err != nil {
		return err
	}
	d.DatabaseClassification = res
	return nil
}

// DatabasesClientListByClusterResponse contains the response from method DatabasesClient.ListByCluster.
type DatabasesClientListByClusterResponse struct {
	DatabaseListResult
}

// DatabasesClientListPrincipalsResponse contains the response from method DatabasesClient.ListPrincipals.
type DatabasesClientListPrincipalsResponse struct {
	DatabasePrincipalListResult
}

// DatabasesClientRemovePrincipalsResponse contains the response from method DatabasesClient.RemovePrincipals.
type DatabasesClientRemovePrincipalsResponse struct {
	DatabasePrincipalListResult
}

// DatabasesClientUpdateResponse contains the response from method DatabasesClient.Update.
type DatabasesClientUpdateResponse struct {
	DatabaseClassification
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DatabasesClientUpdateResponse.
func (d *DatabasesClientUpdateResponse) UnmarshalJSON(data []byte) error {
	res, err := unmarshalDatabaseClassification(data)
	if err != nil {
		return err
	}
	d.DatabaseClassification = res
	return nil
}

// ManagedPrivateEndpointsClientCheckNameAvailabilityResponse contains the response from method ManagedPrivateEndpointsClient.CheckNameAvailability.
type ManagedPrivateEndpointsClientCheckNameAvailabilityResponse struct {
	CheckNameResult
}

// ManagedPrivateEndpointsClientCreateOrUpdateResponse contains the response from method ManagedPrivateEndpointsClient.CreateOrUpdate.
type ManagedPrivateEndpointsClientCreateOrUpdateResponse struct {
	ManagedPrivateEndpoint
}

// ManagedPrivateEndpointsClientDeleteResponse contains the response from method ManagedPrivateEndpointsClient.Delete.
type ManagedPrivateEndpointsClientDeleteResponse struct {
	// placeholder for future response values
}

// ManagedPrivateEndpointsClientGetResponse contains the response from method ManagedPrivateEndpointsClient.Get.
type ManagedPrivateEndpointsClientGetResponse struct {
	ManagedPrivateEndpoint
}

// ManagedPrivateEndpointsClientListResponse contains the response from method ManagedPrivateEndpointsClient.List.
type ManagedPrivateEndpointsClientListResponse struct {
	ManagedPrivateEndpointListResult
}

// ManagedPrivateEndpointsClientUpdateResponse contains the response from method ManagedPrivateEndpointsClient.Update.
type ManagedPrivateEndpointsClientUpdateResponse struct {
	ManagedPrivateEndpoint
}

// OperationsClientListResponse contains the response from method OperationsClient.List.
type OperationsClientListResponse struct {
	OperationListResult
}

// OperationsResultsClientGetResponse contains the response from method OperationsResultsClient.Get.
type OperationsResultsClientGetResponse struct {
	OperationResult
}

// OperationsResultsLocationClientGetResponse contains the response from method OperationsResultsLocationClient.Get.
type OperationsResultsLocationClientGetResponse struct {
	// placeholder for future response values
}

// PrivateEndpointConnectionsClientCreateOrUpdateResponse contains the response from method PrivateEndpointConnectionsClient.CreateOrUpdate.
type PrivateEndpointConnectionsClientCreateOrUpdateResponse struct {
	PrivateEndpointConnection
}

// PrivateEndpointConnectionsClientDeleteResponse contains the response from method PrivateEndpointConnectionsClient.Delete.
type PrivateEndpointConnectionsClientDeleteResponse struct {
	// placeholder for future response values
}

// PrivateEndpointConnectionsClientGetResponse contains the response from method PrivateEndpointConnectionsClient.Get.
type PrivateEndpointConnectionsClientGetResponse struct {
	PrivateEndpointConnection
}

// PrivateEndpointConnectionsClientListResponse contains the response from method PrivateEndpointConnectionsClient.List.
type PrivateEndpointConnectionsClientListResponse struct {
	PrivateEndpointConnectionListResult
}

// PrivateLinkResourcesClientGetResponse contains the response from method PrivateLinkResourcesClient.Get.
type PrivateLinkResourcesClientGetResponse struct {
	PrivateLinkResource
}

// PrivateLinkResourcesClientListResponse contains the response from method PrivateLinkResourcesClient.List.
type PrivateLinkResourcesClientListResponse struct {
	PrivateLinkResourceListResult
}

// ScriptsClientCheckNameAvailabilityResponse contains the response from method ScriptsClient.CheckNameAvailability.
type ScriptsClientCheckNameAvailabilityResponse struct {
	CheckNameResult
}

// ScriptsClientCreateOrUpdateResponse contains the response from method ScriptsClient.CreateOrUpdate.
type ScriptsClientCreateOrUpdateResponse struct {
	Script
}

// ScriptsClientDeleteResponse contains the response from method ScriptsClient.Delete.
type ScriptsClientDeleteResponse struct {
	// placeholder for future response values
}

// ScriptsClientGetResponse contains the response from method ScriptsClient.Get.
type ScriptsClientGetResponse struct {
	Script
}

// ScriptsClientListByDatabaseResponse contains the response from method ScriptsClient.ListByDatabase.
type ScriptsClientListByDatabaseResponse struct {
	ScriptListResult
}

// ScriptsClientUpdateResponse contains the response from method ScriptsClient.Update.
type ScriptsClientUpdateResponse struct {
	Script
}
