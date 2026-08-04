// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

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
