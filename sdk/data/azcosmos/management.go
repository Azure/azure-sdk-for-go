// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	// Operation holds the settings every operation accepts. CreateDatabase always requests the
	// response body, so Operation.EnableContentResponseOnWrite is ignored.
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
	// Operation holds the settings every operation accepts. CreateContainer always requests the
	// response body, so Operation.EnableContentResponseOnWrite is ignored.
	Operation OperationOptions
}

// ReadContainerOptions configures [ContainerClient.Read]. A nil *ReadContainerOptions
// selects the defaults for every field.
type ReadContainerOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions
}

// DeleteContainerOptions configures [ContainerClient.Delete]. A nil
// *DeleteContainerOptions selects the defaults for every field.
type DeleteContainerOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions
}

// errKeyCredentialRequired reports that a management mutation requires a client created with
// [NewClientWithKey]. Microsoft Entra ID can read metadata with the appropriate data-plane role,
// but the service requires master-key authentication for database and container mutations.
func errKeyCredentialRequired() error {
	return &Error{
		Code: CodeClientError,
		Message: "azcosmos: database and container management mutations require a client " +
			"created with NewClientWithKey; Microsoft Entra ID authentication (NewClient) is not " +
			"supported for create or delete operations",
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

func newDatabaseResponse(response ItemResponse, body []byte) (DatabaseResponse, error) {
	result := DatabaseResponse{Response: response.Response}
	if len(body) == 0 {
		return result, nil
	}
	var properties databasePropertiesWire
	if err := json.Unmarshal(body, &properties); err != nil {
		return DatabaseResponse{}, newManagementSerializationError(response, body, "database", err)
	}
	result.DatabaseProperties = properties.properties()
	return result, nil
}

func newContainerResponse(response ItemResponse, body []byte) (ContainerResponse, error) {
	result := ContainerResponse{Response: response.Response}
	if len(body) == 0 {
		return result, nil
	}
	var properties containerPropertiesWire
	if err := json.Unmarshal(body, &properties); err != nil {
		return ContainerResponse{}, newManagementSerializationError(response, body, "container", err)
	}
	result.ContainerProperties = properties.properties()
	return result, nil
}

func newManagementSerializationError(response ItemResponse, body []byte, resource string, cause error) *Error {
	return &Error{
		Code:          CodeSerializationFailed,
		Diagnostics:   response.Diagnostics,
		StatusCode:    response.StatusCode,
		SubStatus:     response.SubStatus,
		AttemptCount:  response.AttemptCount,
		RequestCharge: response.RequestCharge,
		ActivityID:    response.ActivityID,
		SessionToken:  response.SessionToken,
		ETag:          response.ETag,
		Message:       fmt.Sprintf("azcosmos: decoding %s response", resource),
		Body:          append([]byte(nil), body...),
		FromWire:      true,
		cause:         cause,
	}
}

// read implements [DatabaseClient.Read].
func (d *DatabaseClient) read(ctx context.Context, options *ReadDatabaseOptions) (DatabaseResponse, error) {
	if options != nil {
		if err := options.Operation.ConsistencyStrategy.validate(); err != nil {
			return DatabaseResponse{}, err
		}
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

	response, respBody, err := d.client.executeManagement(ctx, operationKindCreateContainer, d.id, properties.ID, body, operation)
	if err != nil {
		return ContainerResponse{}, err
	}
	return newContainerResponse(response, respBody)
}

// readContainer implements [ContainerClient.Read].
func (c *ContainerClient) readContainer(ctx context.Context, options *ReadContainerOptions) (ContainerResponse, error) {
	if options != nil {
		if err := options.Operation.ConsistencyStrategy.validate(); err != nil {
			return ContainerResponse{}, err
		}
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

// deleteContainer implements [ContainerClient.Delete].
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
