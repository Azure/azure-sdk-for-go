// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// AddEntityResponse contains response fields for Client.AddEntityResponse
type AddEntityResponse struct {
	// ETag contains the information returned from the ETag header response.
	ETag azcore.ETag

	// The OData properties of the table entity in JSON format.
	Value []byte
}

// CreateTableResponse contains response fields for Client.Create and ServiceClient.CreateTable
type CreateTableResponse struct {
	// The name of the table.
	TableName *string `json:"TableName,omitempty"`
}

// DeleteEntityResponse contains response fields for Client.DeleteEntity
type DeleteEntityResponse struct {
	// placeholder for future optional response fields
}

// DeleteTableResponse contains response fields for ServiceClient.DeleteTable and Client.Delete
type DeleteTableResponse struct {
	// placeholder for future optional response fields
}

// GetAccessPolicyResponse contains response fields for Client.GetAccessPolicy
type GetAccessPolicyResponse struct {
	SignedIdentifiers []*SignedIdentifier
}

// GetEntityResponse contains response fields for Client.GetEntity
type GetEntityResponse struct {
	// ETag contains the information returned from the ETag header response.
	ETag azcore.ETag

	// The OData properties of the table entity in JSON format.
	Value []byte
}

// GetPropertiesResponse contains response fields for Client.GetProperties
type GetPropertiesResponse struct {
	ServiceProperties
}

// GetStatisticsResponse contains response fields for Client.GetStatistics
type GetStatisticsResponse struct {
	GeoReplication *GeoReplication `xml:"GeoReplication"`
}

// ListEntitiesResponse contains response fields for ListEntitiesPager.NextPage
type ListEntitiesResponse struct {
	// NextPartitionKey contains the information returned from the x-ms-continuation-NextPartitionKey header response.
	NextPartitionKey *string

	// NextRowKey contains the information returned from the x-ms-continuation-NextRowKey header response.
	NextRowKey *string

	// List of table entities.
	Entities [][]byte
}

// ListTablesResponse contains response fields for ListTablesPager.NextPage
type ListTablesResponse struct {
	// NextTableName contains the information returned from the x-ms-continuation-NextTableName header response.
	NextTableName *string

	// List of tables.
	Tables []*TableProperties `json:"value,omitempty"`
}

// SetAccessPolicyResponse contains response fields for Client.SetAccessPolicy
type SetAccessPolicyResponse struct {
	// placeholder for future optional parameters
}

// SetPropertiesResponse contains response fields for Client.SetProperties
type SetPropertiesResponse struct {
	// placeholder for future response fields
}

// TransactionResponse contains response fields for Client.TransactionResponse
type TransactionResponse struct {
	// placeholder for future response fields
}

// UpdateEntityResponse contains response fields for Client.UpdateEntity
type UpdateEntityResponse struct {
	ETag azcore.ETag
}

// UpsertEntityResponse contains response fields for Client.InsertEntity
type UpsertEntityResponse struct {
	ETag azcore.ETag
}
