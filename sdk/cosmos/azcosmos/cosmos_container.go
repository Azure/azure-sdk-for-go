// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// A CosmosContainer lets you perform read, update, change throughput, and delete container operations.
// It also lets you perform read, update, change throughput, and delete item operations.
type CosmosContainer struct {
	// The Id of the Cosmos container
	Id string
	// The database that contains the container
	Database *CosmosDatabase
}

func newCosmosContainer(id string, database *CosmosDatabase) *CosmosContainer {
	return &CosmosContainer{Id: id, Database: database}
}
