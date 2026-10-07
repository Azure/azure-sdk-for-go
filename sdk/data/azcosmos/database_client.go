// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
)

// DatabaseClient is a client for a database in a Cosmos DB account. Obtain one from
// [Client.NewDatabase].
type DatabaseClient struct {
	id     string
	client *Client
}

// ID returns the identifier of the database.
func (d *DatabaseClient) ID() string {
	return d.id
}

// NewContainer returns a client for a container in the database. It does not contact the service,
// so it succeeds whether or not the container exists.
func (d *DatabaseClient) NewContainer(id string) (*ContainerClient, error) {
	if id == "" {
		return nil, errors.New("azcosmos: container id must not be empty")
	}
	return &ContainerClient{id: id, database: d}, nil
}

// Read reads the database's properties. options may be nil.
//
// Azure Cosmos DB only accepts master-key authentication for this operation: a client created
// with [NewClient] (Microsoft Entra ID) returns an error rather than reaching the service. Use
// [NewClientWithKey].
func (d *DatabaseClient) Read(ctx context.Context, options *ReadDatabaseOptions) (DatabaseResponse, error) {
	return d.read(ctx, options)
}

// Delete deletes the database and every container in it. options may be nil.
//
// Azure Cosmos DB only accepts master-key authentication for this operation: a client created
// with [NewClient] (Microsoft Entra ID) returns an error rather than reaching the service. Use
// [NewClientWithKey].
func (d *DatabaseClient) Delete(ctx context.Context, options *DeleteDatabaseOptions) (DatabaseResponse, error) {
	return d.delete(ctx, options)
}

// CreateContainer creates a new container in the database.
//
// properties.ID and properties.PartitionKeyDefinition are required; the partition key definition
// is immutable once created. options may be nil.
//
// Azure Cosmos DB only accepts master-key authentication for this operation: a client created
// with [NewClient] (Microsoft Entra ID) returns an error rather than reaching the service. Use
// [NewClientWithKey].
func (d *DatabaseClient) CreateContainer(ctx context.Context, properties ContainerProperties, options *CreateContainerOptions) (ContainerResponse, error) {
	return d.createContainer(ctx, properties, options)
}
