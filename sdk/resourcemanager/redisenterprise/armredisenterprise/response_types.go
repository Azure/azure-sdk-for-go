//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armredisenterprise

// ClientCreateResponse contains the response from method Client.BeginCreate.
type ClientCreateResponse struct {
	// Describes the RedisEnterprise cluster
	Cluster
}

// ClientDeleteResponse contains the response from method Client.BeginDelete.
type ClientDeleteResponse struct {
	// placeholder for future response values
}

// ClientGetResponse contains the response from method Client.Get.
type ClientGetResponse struct {
	// Describes the RedisEnterprise cluster
	Cluster
}

// ClientListByResourceGroupResponse contains the response from method Client.NewListByResourceGroupPager.
type ClientListByResourceGroupResponse struct {
	// The response of a list-all operation.
	ClusterList
}

// ClientListResponse contains the response from method Client.NewListPager.
type ClientListResponse struct {
	// The response of a list-all operation.
	ClusterList
}

// ClientUpdateResponse contains the response from method Client.BeginUpdate.
type ClientUpdateResponse struct {
	// Describes the RedisEnterprise cluster
	Cluster
}

// DatabasesClientCreateResponse contains the response from method DatabasesClient.BeginCreate.
type DatabasesClientCreateResponse struct {
	// Describes a database on the RedisEnterprise cluster
	Database
}

// DatabasesClientDeleteResponse contains the response from method DatabasesClient.BeginDelete.
type DatabasesClientDeleteResponse struct {
	// placeholder for future response values
}

// DatabasesClientExportResponse contains the response from method DatabasesClient.BeginExport.
type DatabasesClientExportResponse struct {
	// placeholder for future response values
}

// DatabasesClientForceUnlinkResponse contains the response from method DatabasesClient.BeginForceUnlink.
type DatabasesClientForceUnlinkResponse struct {
	// placeholder for future response values
}

// DatabasesClientGetResponse contains the response from method DatabasesClient.Get.
type DatabasesClientGetResponse struct {
	// Describes a database on the RedisEnterprise cluster
	Database
}

// DatabasesClientImportResponse contains the response from method DatabasesClient.BeginImport.
type DatabasesClientImportResponse struct {
	// placeholder for future response values
}

// DatabasesClientListByClusterResponse contains the response from method DatabasesClient.NewListByClusterPager.
type DatabasesClientListByClusterResponse struct {
	// The response of a list-all operation.
	DatabaseList
}

// DatabasesClientListKeysResponse contains the response from method DatabasesClient.ListKeys.
type DatabasesClientListKeysResponse struct {
	// The secret access keys used for authenticating connections to redis
	AccessKeys
}

// DatabasesClientRegenerateKeyResponse contains the response from method DatabasesClient.BeginRegenerateKey.
type DatabasesClientRegenerateKeyResponse struct {
	// The secret access keys used for authenticating connections to redis
	AccessKeys
}

// DatabasesClientUpdateResponse contains the response from method DatabasesClient.BeginUpdate.
type DatabasesClientUpdateResponse struct {
	// Describes a database on the RedisEnterprise cluster
	Database
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to get the next set of results.
	OperationListResult
}

// OperationsStatusClientGetResponse contains the response from method OperationsStatusClient.Get.
type OperationsStatusClientGetResponse struct {
	// The status of a long-running operation.
	OperationStatus
}

// PrivateEndpointConnectionsClientDeleteResponse contains the response from method PrivateEndpointConnectionsClient.Delete.
type PrivateEndpointConnectionsClientDeleteResponse struct {
	// placeholder for future response values
}

// PrivateEndpointConnectionsClientGetResponse contains the response from method PrivateEndpointConnectionsClient.Get.
type PrivateEndpointConnectionsClientGetResponse struct {
	// The Private Endpoint Connection resource.
	PrivateEndpointConnection
}

// PrivateEndpointConnectionsClientListResponse contains the response from method PrivateEndpointConnectionsClient.NewListPager.
type PrivateEndpointConnectionsClientListResponse struct {
	// List of private endpoint connection associated with the specified storage account
	PrivateEndpointConnectionListResult
}

// PrivateEndpointConnectionsClientPutResponse contains the response from method PrivateEndpointConnectionsClient.BeginPut.
type PrivateEndpointConnectionsClientPutResponse struct {
	// The Private Endpoint Connection resource.
	PrivateEndpointConnection
}

// PrivateLinkResourcesClientListByClusterResponse contains the response from method PrivateLinkResourcesClient.NewListByClusterPager.
type PrivateLinkResourcesClientListByClusterResponse struct {
	// A list of private link resources
	PrivateLinkResourceListResult
}

