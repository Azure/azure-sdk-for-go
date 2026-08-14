// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "errors"

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
