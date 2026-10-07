// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "context"

// ContainerClient is a client for a container in a Cosmos DB database, and is where item
// operations live. Obtain one from [Client.NewContainer] or [DatabaseClient.NewContainer].
type ContainerClient struct {
	id       string
	database *DatabaseClient
}

// ID returns the identifier of the container.
func (c *ContainerClient) ID() string {
	return c.id
}

// ReadContainer reads the container's properties. options may be nil.
//
// Azure Cosmos DB only accepts master-key authentication for this operation: a client created
// with [NewClient] (Microsoft Entra ID) returns an error rather than reaching the service. Use
// [NewClientWithKey].
func (c *ContainerClient) ReadContainer(ctx context.Context, options *ReadContainerOptions) (ContainerResponse, error) {
	return c.readContainer(ctx, options)
}

// DeleteContainer deletes the container and every item in it. options may be nil.
//
// Azure Cosmos DB only accepts master-key authentication for this operation: a client created
// with [NewClient] (Microsoft Entra ID) returns an error rather than reaching the service. Use
// [NewClientWithKey].
func (c *ContainerClient) DeleteContainer(ctx context.Context, options *DeleteContainerOptions) (ContainerResponse, error) {
	return c.deleteContainer(ctx, options)
}
