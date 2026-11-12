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

// NewContainer resolves existing container metadata through the native cache before returning.
// Cache misses contact the service; this method does not create a container.
func (d *DatabaseClient) NewContainer(ctx context.Context, id string) (*ContainerClient, error) {
	if id == "" {
		return nil, errors.New("azcosmos: container id must not be empty")
	}
	if ctx == nil {
		return nil, errors.New("azcosmos: context must not be nil")
	}
	release, err := d.client.acquire()
	if err != nil {
		return nil, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := d.client.driver.resolveContainer(ctx, d.id, id); err != nil {
		return nil, err
	}
	return &ContainerClient{id: id, database: d}, nil
}

// Read reads the database's properties. options may be nil.
//
// Microsoft Entra ID authentication through [NewClient] is supported when the identity has the
// Microsoft.DocumentDB/databaseAccounts/readMetadata data action.
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
// is immutable once created. An omitted partition key kind is inferred from the number of paths,
// and an omitted version creates a version 2 partition key. options may be nil. The operation
// always requests response content, so
// CreateContainerOptions.Operation.EnableContentResponseOnWrite is ignored.
//
// Azure Cosmos DB only accepts master-key authentication for this operation: a client created
// with [NewClient] (Microsoft Entra ID) returns an error rather than reaching the service. Use
// [NewClientWithKey].
func (d *DatabaseClient) CreateContainer(ctx context.Context, properties ContainerProperties, options *CreateContainerOptions) (ContainerResponse, error) {
	return d.createContainer(ctx, properties, options)
}
