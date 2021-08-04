// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "errors"

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

// GetContainer returns a CosmosContainer object for the container.
// id - The id of the container.
func (db *CosmosDatabase) GetContainer(id string) (*CosmosContainer, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return newCosmosContainer(id, db), nil
}
