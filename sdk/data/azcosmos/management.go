// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
)

// Management operation kinds, mirroring the driver's CosmosOperationKind discriminants used for
// database and container lifecycle operations. See cgo_native.go for the static asserts that pin
// these against the vendored header.
const (
	operationKindCreateDatabase  operationKind = 1
	operationKindReadDatabase    operationKind = 7
	operationKindDeleteDatabase  operationKind = 8
	operationKindCreateContainer operationKind = 9
	operationKindReadContainer   operationKind = 12
	operationKindDeleteContainer operationKind = 14
)

// CreateDatabaseOptions configures [Client.CreateDatabase]. A nil *CreateDatabaseOptions selects
// the defaults for every field.
type CreateDatabaseOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions
}

// ReadDatabaseOptions configures [DatabaseClient.Read]. A nil *ReadDatabaseOptions selects the
// defaults for every field.
type ReadDatabaseOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions
}

// DeleteDatabaseOptions configures [DatabaseClient.Delete]. A nil *DeleteDatabaseOptions selects
// the defaults for every field.
type DeleteDatabaseOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions
}

// CreateContainerOptions configures [DatabaseClient.CreateContainer]. A nil *CreateContainerOptions
// selects the defaults for every field.
type CreateContainerOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions
}

// ReadContainerOptions configures [ContainerClient.ReadContainer]. A nil *ReadContainerOptions
// selects the defaults for every field.
type ReadContainerOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions
}

// DeleteContainerOptions configures [ContainerClient.DeleteContainer]. A nil
// *DeleteContainerOptions selects the defaults for every field.
type DeleteContainerOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions
}

// errKeyCredentialRequired reports that a management operation requires a client created with
// [NewClientWithKey]. Azure Cosmos DB only accepts master-key authentication for database and
// container management operations issued through the data-plane SQL API; a client authenticated
// with Microsoft Entra ID (NewClient) cannot perform them.
func errKeyCredentialRequired() error {
	return &Error{
		Code: CodeClientError,
		Message: "azcosmos: database and container management operations require a client " +
			"created with NewClientWithKey; Microsoft Entra ID authentication (NewClient) is not " +
			"supported for these operations",
	}
}

// createDatabase implements [Client.CreateDatabase]; kept free of the public method's doc comment
// so both it and tests can call the same validated path.
func (c *Client) createDatabase(ctx context.Context, properties DatabaseProperties, options *CreateDatabaseOptions) (DatabaseResponse, error) {
	if properties.ID == "" {
		return DatabaseResponse{}, errors.New("azcosmos: DatabaseProperties.ID must not be empty")
	}
	if !c.usesKeyAuth {
		return DatabaseResponse{}, errKeyCredentialRequired()
	}
	release, err := c.acquire()
	if err != nil {
		return DatabaseResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return DatabaseResponse{}, err
	}

	body, err := json.Marshal(databaseCreateBody{ID: properties.ID})
	if err != nil {
		return DatabaseResponse{}, err
	}

	operation := OperationOptions{}
	if options != nil {
		operation = options.Operation
	}
	// A database create has no purpose without its server-assigned properties (resource ID,
	// self-link, etag) coming back, so this always asks for the response body regardless of
	// any write-content preference the caller configured for item operations.
	forceContentResponseOnWrite(&operation)

	response, respBody, err := c.executeManagement(ctx, operationKindCreateDatabase, "", "", body, operation)
	if err != nil {
		return DatabaseResponse{}, err
	}
	return newDatabaseResponse(response, respBody)
}

// forceContentResponseOnWrite overrides the operation's write-content preference to always
// return the response body, matching the v1 SDK's behavior for database/container creates.
func forceContentResponseOnWrite(operation *OperationOptions) {
	alwaysReturnContent := true
	operation.EnableContentResponseOnWrite = &alwaysReturnContent
}

// databaseCreateBody is the wire body for creating a database: only ID is settable.
type databaseCreateBody struct {
	ID string `json:"id"`
}

func newDatabaseResponse(response Response, body []byte) (DatabaseResponse, error) {
	result := DatabaseResponse{Response: response}
	if len(body) == 0 {
		return result, nil
	}
	if err := json.Unmarshal(body, &result.DatabaseProperties); err != nil {
		return DatabaseResponse{}, err
	}
	return result, nil
}

func newContainerResponse(response Response, body []byte) (ContainerResponse, error) {
	result := ContainerResponse{Response: response}
	if len(body) == 0 {
		return result, nil
	}
	if err := json.Unmarshal(body, &result.ContainerProperties); err != nil {
		return ContainerResponse{}, err
	}
	return result, nil
}

// read implements [DatabaseClient.Read].
func (d *DatabaseClient) read(ctx context.Context, options *ReadDatabaseOptions) (DatabaseResponse, error) {
	if !d.client.usesKeyAuth {
		return DatabaseResponse{}, errKeyCredentialRequired()
	}
	release, err := d.client.acquire()
	if err != nil {
		return DatabaseResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return DatabaseResponse{}, err
	}

	operation := OperationOptions{}
	if options != nil {
		operation = options.Operation
	}

	response, body, err := d.client.executeManagement(ctx, operationKindReadDatabase, d.id, "", nil, operation)
	if err != nil {
		return DatabaseResponse{}, err
	}
	return newDatabaseResponse(response, body)
}

// delete implements [DatabaseClient.Delete].
func (d *DatabaseClient) delete(ctx context.Context, options *DeleteDatabaseOptions) (DatabaseResponse, error) {
	if !d.client.usesKeyAuth {
		return DatabaseResponse{}, errKeyCredentialRequired()
	}
	release, err := d.client.acquire()
	if err != nil {
		return DatabaseResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return DatabaseResponse{}, err
	}

	operation := OperationOptions{}
	if options != nil {
		operation = options.Operation
	}

	response, body, err := d.client.executeManagement(ctx, operationKindDeleteDatabase, d.id, "", nil, operation)
	if err != nil {
		return DatabaseResponse{}, err
	}
	return newDatabaseResponse(response, body)
}

// createContainer implements [DatabaseClient.CreateContainer].
func (d *DatabaseClient) createContainer(ctx context.Context, properties ContainerProperties, options *CreateContainerOptions) (ContainerResponse, error) {
	if properties.ID == "" {
		return ContainerResponse{}, errors.New("azcosmos: ContainerProperties.ID must not be empty")
	}
	if len(properties.PartitionKeyDefinition.Paths) == 0 {
		return ContainerResponse{}, errors.New("azcosmos: ContainerProperties.PartitionKeyDefinition.Paths must not be empty")
	}
	if !d.client.usesKeyAuth {
		return ContainerResponse{}, errKeyCredentialRequired()
	}
	release, err := d.client.acquire()
	if err != nil {
		return ContainerResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return ContainerResponse{}, err
	}

	body, err := json.Marshal(properties)
	if err != nil {
		return ContainerResponse{}, err
	}

	operation := OperationOptions{}
	if options != nil {
		operation = options.Operation
	}
	// Same rationale as [Client.createDatabase]: the caller always needs the server-assigned
	// properties back from a container create.
	forceContentResponseOnWrite(&operation)

	response, respBody, err := d.client.executeManagement(ctx, operationKindCreateContainer, d.id, "", body, operation)
	if err != nil {
		return ContainerResponse{}, err
	}
	return newContainerResponse(response, respBody)
}

// readContainer implements [ContainerClient.ReadContainer].
func (c *ContainerClient) readContainer(ctx context.Context, options *ReadContainerOptions) (ContainerResponse, error) {
	if !c.database.client.usesKeyAuth {
		return ContainerResponse{}, errKeyCredentialRequired()
	}
	release, err := c.database.client.acquire()
	if err != nil {
		return ContainerResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return ContainerResponse{}, err
	}

	operation := OperationOptions{}
	if options != nil {
		operation = options.Operation
	}

	response, body, err := c.database.client.executeManagement(ctx, operationKindReadContainer, c.database.id, c.id, nil, operation)
	if err != nil {
		return ContainerResponse{}, err
	}
	return newContainerResponse(response, body)
}

// deleteContainer implements [ContainerClient.DeleteContainer].
func (c *ContainerClient) deleteContainer(ctx context.Context, options *DeleteContainerOptions) (ContainerResponse, error) {
	if !c.database.client.usesKeyAuth {
		return ContainerResponse{}, errKeyCredentialRequired()
	}
	release, err := c.database.client.acquire()
	if err != nil {
		return ContainerResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return ContainerResponse{}, err
	}

	operation := OperationOptions{}
	if options != nil {
		operation = options.Operation
	}

	response, body, err := c.database.client.executeManagement(ctx, operationKindDeleteContainer, c.database.id, c.id, nil, operation)
	if err != nil {
		return ContainerResponse{}, err
	}
	return newContainerResponse(response, body)
}
