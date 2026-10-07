// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDatabaseRejectsEmptyID(t *testing.T) {
	client := newTestClient(t)
	_, err := client.CreateDatabase(context.Background(), DatabaseProperties{}, nil)
	require.ErrorContains(t, err, "DatabaseProperties.ID")
}

func TestCreateContainerRejectsEmptyID(t *testing.T) {
	client := newTestClient(t)
	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	_, err = database.CreateContainer(context.Background(), ContainerProperties{
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}},
	}, nil)
	require.ErrorContains(t, err, "ContainerProperties.ID")
}

func TestCreateContainerRejectsMissingPartitionKeyPaths(t *testing.T) {
	client := newTestClient(t)
	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	_, err = database.CreateContainer(context.Background(), ContainerProperties{ID: "items"}, nil)
	require.ErrorContains(t, err, "PartitionKeyDefinition.Paths")
}

// Management operations are only valid against a master-key client: Azure Cosmos DB does not
// accept Microsoft Entra ID authentication for database/container lifecycle operations issued
// through the data-plane SQL API.
func TestManagementOperationsRejectTokenCredentialClient(t *testing.T) {
	client, err := NewClient("https://myaccount.documents.azure.com", fakeTokenCredential{}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	container, err := database.NewContainer("items")
	require.NoError(t, err)

	ctx := context.Background()
	validContainer := ContainerProperties{ID: "items", PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}}}

	_, err = client.CreateDatabase(ctx, DatabaseProperties{ID: "db"}, nil)
	require.ErrorContains(t, err, "NewClientWithKey")

	_, err = database.Read(ctx, nil)
	require.ErrorContains(t, err, "NewClientWithKey")

	_, err = database.Delete(ctx, nil)
	require.ErrorContains(t, err, "NewClientWithKey")

	_, err = database.CreateContainer(ctx, validContainer, nil)
	require.ErrorContains(t, err, "NewClientWithKey")

	_, err = container.ReadContainer(ctx, nil)
	require.ErrorContains(t, err, "NewClientWithKey")

	_, err = container.DeleteContainer(ctx, nil)
	require.ErrorContains(t, err, "NewClientWithKey")
}

func TestDatabaseResponseParsesProperties(t *testing.T) {
	response, err := newDatabaseResponse(Response{}, []byte(`{"id":"db","_etag":"\"abc\"","_rid":"rid","_self":"self"}`))
	require.NoError(t, err)
	require.Equal(t, "db", response.DatabaseProperties.ID)
	require.Equal(t, "rid", response.DatabaseProperties.ResourceID)
	require.NotNil(t, response.DatabaseProperties.ETag)
	require.Equal(t, "\"abc\"", string(*response.DatabaseProperties.ETag))
}

func TestDatabaseResponseEmptyBodyIsZeroValue(t *testing.T) {
	response, err := newDatabaseResponse(Response{StatusCode: 204}, nil)
	require.NoError(t, err)
	require.Equal(t, DatabaseProperties{}, response.DatabaseProperties)
	require.Equal(t, 204, response.StatusCode)
}

func TestContainerResponseParsesProperties(t *testing.T) {
	body := []byte(`{
		"id":"items",
		"partitionKey":{"paths":["/pk"],"kind":"Hash"},
		"indexingPolicy":{"automatic":true,"indexingMode":"consistent"},
		"uniqueKeyPolicy":{"uniqueKeys":[{"paths":["/email"]}]}
	}`)
	response, err := newContainerResponse(Response{}, body)
	require.NoError(t, err)
	require.Equal(t, "items", response.ContainerProperties.ID)
	require.Equal(t, []string{"/pk"}, response.ContainerProperties.PartitionKeyDefinition.Paths)
	require.NotNil(t, response.ContainerProperties.IndexingPolicy)
	require.True(t, response.ContainerProperties.IndexingPolicy.Automatic)
	require.NotNil(t, response.ContainerProperties.UniqueKeyPolicy)
	require.Equal(t, []string{"/email"}, response.ContainerProperties.UniqueKeyPolicy.UniqueKeys[0].Paths)
}
