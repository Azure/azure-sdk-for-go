// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// A CosmosDatabase lets you perform read, update, change throughput, and delete database operations.
type CosmosDatabase struct {
	// The Id of the Cosmos database
	Id string
	// The client associated with the Cosmos database
	client *CosmosClient
}

func newCosmosDatabase(id string, client *CosmosClient) *CosmosDatabase {
	return &CosmosDatabase{Id: id, client: client}
}
