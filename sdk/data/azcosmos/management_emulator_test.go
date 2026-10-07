// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEmulatorDatabaseAndContainerLifecycle exercises the public management surface end to end
// against a real emulator: create a database, create a container in it with a partition key and a
// unique key policy, read both back, use the new container for an item operation, then delete the
// container and the database and confirm each is gone.
//
// This is the one place the native management dispatch path (scope resolution, response decoding,
// and the delete-time cache eviction) is exercised as a whole rather than in pieces; the rest of
// this package's management tests stop at validation, auth-gating, or request/response conversion
// in isolation.
func TestEmulatorDatabaseAndContainerLifecycle(t *testing.T) {
	client, _, _ := emulatorClientConfigured(t, nil, false)
	ctx := context.Background()

	databaseID := uniqueItemID(t) + "-db"
	containerID := uniqueItemID(t) + "-container"

	dbResp, err := client.CreateDatabase(ctx, DatabaseProperties{ID: databaseID}, nil)
	require.NoError(t, err)
	require.Equal(t, 201, dbResp.StatusCode)
	require.Equal(t, databaseID, dbResp.DatabaseProperties.ID)
	require.NotEmpty(t, dbResp.DatabaseProperties.ResourceID, "a create response has to carry the server-assigned fields")

	database, err := client.NewDatabase(databaseID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = database.Delete(context.Background(), nil)
	})

	readDBResp, err := database.Read(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, 200, readDBResp.StatusCode)
	require.Equal(t, databaseID, readDBResp.DatabaseProperties.ID)

	containerProps := ContainerProperties{
		ID: containerID,
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
		UniqueKeyPolicy: &UniqueKeyPolicy{
			UniqueKeys: []UniqueKey{{Paths: []string{"/email"}}},
		},
	}
	containerResp, err := database.CreateContainer(ctx, containerProps, nil)
	require.NoError(t, err)
	require.Equal(t, 201, containerResp.StatusCode)
	require.Equal(t, containerID, containerResp.ContainerProperties.ID)
	require.Equal(t, []string{"/pk"}, containerResp.ContainerProperties.PartitionKeyDefinition.Paths)
	require.Equal(t, PartitionKeyKindHash, containerResp.ContainerProperties.PartitionKeyDefinition.Kind,
		"a single path has to infer Hash on the wire, not leave the kind unset")

	container, err := database.NewContainer(containerID)
	require.NoError(t, err)

	readContainerResp, err := container.ReadContainer(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, 200, readContainerResp.StatusCode)
	require.Equal(t, containerID, readContainerResp.ContainerProperties.ID)

	// The newly-created container has to be immediately usable for ordinary item operations, not
	// just visible to management reads.
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	item := []byte(`{"id":"` + id + `","pk":"` + id + `","email":"a@example.com"}`)
	_, err = container.CreateItem(ctx, pk, id, item, nil)
	require.NoError(t, err)

	deleteContainerResp, err := container.DeleteContainer(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, 204, deleteContainerResp.StatusCode)

	// The container is gone, so creating a fresh one under the same name and using it right away
	// has to resolve the new container, not reuse the deleted one's cached resource id.
	recreatedResp, err := database.CreateContainer(ctx, containerProps, nil)
	require.NoError(t, err)
	require.Equal(t, 201, recreatedResp.StatusCode)
	recreatedContainer, err := database.NewContainer(containerID)
	require.NoError(t, err)
	_, err = recreatedContainer.CreateItem(ctx, pk, id, item, nil)
	require.NoError(t, err, "the recreated container has to accept writes under its own, fresh resource id")
	_, err = recreatedContainer.DeleteContainer(ctx, nil)
	require.NoError(t, err)

	deleteDBResp, err := database.Delete(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, 204, deleteDBResp.StatusCode)

	_, err = database.Read(ctx, nil)
	require.Error(t, err)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeNotFound, cosmosErr.Code)
}
