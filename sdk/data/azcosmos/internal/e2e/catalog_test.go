// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package e2e

import (
	"encoding/json"
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

// validScenarioJSON and validProfileJSON are minimal documents satisfying every field this
// package's parseScenario/parseProfile require, used as a base for the drift-detection tests
// below so each test removes exactly one field instead of hand-maintaining a full fixture twice.
const (
	validScenarioJSON = `{
		"$schema": "../../schema/scenario.v1.json",
		"specVersion": "1.0",
		"id": "test.drift",
		"title": "t",
		"requirement": "r",
		"maturity": "stable",
		"precedents": [{"sdk": "rust", "path": "p", "test": "t"}],
		"profiles": ["smokeTests"],
		"tags": ["t"],
		"backends": {
			"hostedEmulatorGatewayV1": {"applicability": "required", "fidelity": "exact"},
			"hostedEmulatorGatewayV2": {"applicability": "required", "fidelity": "exact"},
			"azureLive": {"applicability": "required", "fidelity": "exact"}
		}
	}`
	validProfileJSON = `{
		"$schema": "../schema/profile.v1.json",
		"specVersion": "1.0",
		"id": "test.profile",
		"accounts": [{"id": "a"}],
		"runtimes": [{"id": "r"}],
		"clients": [{"id": "c"}]
	}`
)

// TestLoadCatalogRejectsScenarioMissingSchemaFields verifies that parseScenario genuinely fails
// -- not silently succeeds -- when a required shared-schema field (requirement, precedents,
// specVersion, $schema) is absent, so a vendored-snapshot refresh that drops one of these fields
// is caught here instead of only showing up as a quieter downstream behavior change.
func TestLoadCatalogRejectsScenarioMissingSchemaFields(t *testing.T) {
	_, err := parseScenario([]byte(validScenarioJSON), "fixture.json")
	require.NoError(t, err, "the base fixture itself must be valid")

	for _, field := range []string{"$schema", "specVersion", "requirement", "precedents"} {
		var doc map[string]any
		require.NoError(t, json.Unmarshal([]byte(validScenarioJSON), &doc))
		delete(doc, field)
		data, err := json.Marshal(doc)
		require.NoError(t, err)

		_, err = parseScenario(data, "fixture.json")
		require.Errorf(t, err, "removing required field %q must fail scenario parsing", field)
	}
}

// TestLoadCatalogRejectsProfileMissingSchemaFields is TestLoadCatalogRejectsScenarioMissingSchemaFields's
// counterpart for Profile's required $schema/specVersion fields.
func TestLoadCatalogRejectsProfileMissingSchemaFields(t *testing.T) {
	_, err := parseProfile([]byte(validProfileJSON), "fixture.json")
	require.NoError(t, err, "the base fixture itself must be valid")

	for _, field := range []string{"$schema", "specVersion"} {
		var doc map[string]any
		require.NoError(t, json.Unmarshal([]byte(validProfileJSON), &doc))
		delete(doc, field)
		data, err := json.Marshal(doc)
		require.NoError(t, err)

		_, err = parseProfile(data, "fixture.json")
		require.Errorf(t, err, "removing required field %q must fail profile parsing", field)
	}
}
