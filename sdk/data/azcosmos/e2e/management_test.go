// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"context"
	"encoding/json"
	"testing"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/stretchr/testify/require"
)

// TestDatabaseContainerManagementAvailableOperations is supplementary, non-certifying coverage
// for scenario management.resource-lifecycle, whose catalog requirement is "create, read, query,
// replace, and delete" for both databases and containers. #27700 added Create/Read/Delete for
// both resource types to azcosmos/v2's public API, but no Replace and no List/Query for either,
// so the full five-operation contract still cannot be certified; see
// internal/e2e.buildImplementations for that scenario's Blocked reason. This test exercises only
// the three operations that do exist, through an isolated database and container (not the shared
// fixture resources), as regression coverage for #27700's public surface from this package.
func TestDatabaseContainerManagementAvailableOperations(t *testing.T) {
	fx := newRawFixture(t, "")
	ctx := context.Background()

	databaseID := uniqueID(t) + "-db"
	containerID := uniqueID(t) + "-container"

	dbResp, err := fx.Client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: databaseID}, nil)
	require.NoError(t, err)
	require.Equal(t, 201, dbResp.StatusCode)

	database, err := fx.Client.NewDatabase(databaseID)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = database.Delete(context.Background(), nil) })

	readDBResp, err := database.Read(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, databaseID, readDBResp.DatabaseProperties.ID)

	containerProps := azcosmos.ContainerProperties{
		ID: containerID,
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}
	containerResp, err := database.CreateContainer(ctx, containerProps, nil)
	require.NoError(t, err)
	require.Equal(t, 201, containerResp.StatusCode)

	container, err := database.NewContainer(containerID)
	require.NoError(t, err)

	readContainerResp, err := container.Read(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, containerID, readContainerResp.ContainerProperties.ID)

	// The newly-created container has to be immediately usable, not just visible to management
	// reads - the same contract item.* scenarios rely on against the pre-provisioned fixture.
	id := uniqueID(t)
	pk := azcosmos.NewPartitionKeyString(id)
	_, err = container.CreateItem(ctx, pk, id, json.RawMessage(`{"id":"`+id+`","pk":"`+id+`"}`), nil)
	require.NoError(t, err)

	deleteContainerResp, err := container.Delete(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, 204, deleteContainerResp.StatusCode)

	deleteDBResp, err := database.Delete(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, 204, deleteDBResp.StatusCode)
}
