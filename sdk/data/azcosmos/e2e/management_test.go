// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/stretchr/testify/require"
)

func TestEmulatorCapabilities(t *testing.T) {
	profile := selectedProfile(t, "management.capabilities")
	endpoint, managementEndpoint, _, _ := requireProfileEmulator(t, profile)
	require.NoError(t, validateEmulatorBackend(os.Getenv("AZURE_COSMOS_E2E_BACKEND")))
	probeEmulatorCapabilities(t, profile, endpoint, managementEndpoint)

	var raw json.RawMessage
	require.NoError(t, readManagementDocument(t, managementEndpoint, "capabilities", &raw))
	var document map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &document))
	for _, field := range []string{"apiVersion", "emulatorVersion", "protocols", "dataPlane", "managementActions", "limitations"} {
		require.Contains(t, document, field)
		require.NotEqual(t, "null", string(document[field]))
	}
	var protocols map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(document["protocols"], &protocols))
	require.Equal(t, "true", string(protocols["gatewayV1"]))
	require.Equal(t, "false", string(protocols["gatewayV2"]))
	var capabilities emulatorCapabilities
	require.NoError(t, json.Unmarshal(raw, &capabilities))
	require.NoError(t, validateEmulatorCapabilities(profile.backend, capabilities))
	require.NotEmpty(t, capabilities.EmulatorVersion)
	require.True(t, capabilities.Protocols.GatewayV1)
	require.False(t, capabilities.Protocols.GatewayV2)
	require.ElementsMatch(t, []string{"database", "container", "offer", "item", "query", "changeFeed", "transactionalBatch", "patch"}, capabilities.DataPlane)
	require.ElementsMatch(t, []string{"partitionSplit", "partitionMerge", "perPartitionFailover", "replicationPauseResume", "regionLifecycle", "writeRegionFailover"}, capabilities.ManagementActions)
	require.ElementsMatch(t, []string{"authenticationNotEnforced", "loopbackOnly", "volatileStorage"}, capabilities.Limitations)
}

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
	container, containerResp := createTestContainer(t, database, containerProps)
	require.Equal(t, 201, containerResp.StatusCode)

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
