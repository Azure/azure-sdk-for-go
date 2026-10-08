// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadCatalogParsesEveryVendoredScenarioAndProfile(t *testing.T) {
	catalog, err := LoadCatalog()
	require.NoError(t, err)
	require.NotEmpty(t, catalog.Scenarios)
	require.NotEmpty(t, catalog.Profiles)

	for _, id := range []string{"smokeTests", "coreOperations", "configurationResilience"} {
		_, ok := catalog.Profiles[id]
		require.Truef(t, ok, "required profile %q must be present in the pinned catalog", id)
	}

	for id, scenario := range catalog.Scenarios {
		require.Equal(t, id, scenario.ID)
		require.NotEmpty(t, scenario.Profiles, "scenario %q must declare at least one profile", id)
		for _, profileID := range scenario.Profiles {
			_, ok := catalog.Profiles[profileID]
			require.Truef(t, ok, "scenario %q references unknown profile %q", id, profileID)
		}
	}
}

func TestLoadCatalogScenarioCountMatchesVendoredSnapshot(t *testing.T) {
	// Pins the expected count so a future catalog refresh that silently drops files is caught
	// here instead of only showing up as missing coverage.
	catalog, err := LoadCatalog()
	require.NoError(t, err)
	require.Len(t, catalog.Scenarios, 44)
	require.Len(t, catalog.Profiles, 6)
}
