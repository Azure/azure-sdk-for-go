// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"testing"

	internale2e "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2/internal/e2e"
	"github.com/stretchr/testify/require"
)

func TestResolveFixtureProfile(t *testing.T) {
	catalog, err := internale2e.LoadCatalog()
	require.NoError(t, err)

	tests := []struct {
		profileID     string
		accountID     string
		endpointEnv   string
		managementEnv string
	}{
		{
			profileID:     "smokeTests",
			accountID:     supportedSingleRegionAccountID,
			endpointEnv:   defaultEndpointEnv,
			managementEnv: defaultManagementEndpointEnv,
		},
		{
			profileID:     "coreOperations",
			accountID:     supportedSingleRegionAccountID,
			endpointEnv:   defaultEndpointEnv,
			managementEnv: defaultManagementEndpointEnv,
		},
		{
			profileID:     "configurationResilience",
			accountID:     supportedTwoRegionDelayedAccountID,
			endpointEnv:   resilienceEndpointEnv,
			managementEnv: resilienceManagementEndpointEnv,
		},
	}

	for _, test := range tests {
		t.Run(test.profileID, func(t *testing.T) {
			profile, ok := catalog.Profiles[test.profileID]
			require.True(t, ok)

			resolved, err := resolveFixtureProfile(profile, "", "", "")
			require.NoError(t, err)
			require.Equal(t, test.accountID, resolved.accountID)
			require.Equal(t, supportedDefaultRuntimeID, resolved.runtimeID)
			require.Equal(t, supportedDefaultClientID, resolved.clientID)
			require.Equal(t, test.endpointEnv, resolved.endpointEnv)
			require.Equal(t, test.managementEnv, resolved.managementEndpointEnv)
			require.NotNil(t, resolved.clientOptions)
		})
	}
}

func TestResolveFixtureProfileRejectsUnsupportedSetup(t *testing.T) {
	t.Run("unknown axis override", func(t *testing.T) {
		catalog := loadFixtureTestCatalog(t)
		_, err := resolveFixtureProfile(catalog.Profiles["smokeTests"], "otherAccount", "", "")
		require.ErrorContains(t, err, "does not select a definition in the account axis")
	})

	t.Run("dynamic topology account", func(t *testing.T) {
		catalog := loadFixtureTestCatalog(t)
		_, err := resolveFixtureProfile(catalog.Profiles["dynamicTopology"], "", "", "")
		require.ErrorContains(t, err, "unsupported write, consistency, or failover settings")
	})

	t.Run("replication mismatch", func(t *testing.T) {
		catalog := loadFixtureTestCatalog(t)
		profile := catalog.Profiles["configurationResilience"]
		replication := profile.Accounts[0]["replication"].(map[string]any)
		replication["maxDelayMs"] = float64(0)
		_, err := resolveFixtureProfile(profile, "", "", "")
		require.ErrorContains(t, err, "fixed 750ms replication delay")
	})

	t.Run("runtime override", func(t *testing.T) {
		catalog := loadFixtureTestCatalog(t)
		profile := catalog.Profiles["smokeTests"]
		profile.Runtimes[0]["ppcb"] = "enabled"
		_, err := resolveFixtureProfile(profile, "", "", "")
		require.ErrorContains(t, err, "only sdkDefault runtime settings")
	})

	t.Run("client override", func(t *testing.T) {
		catalog := loadFixtureTestCatalog(t)
		profile := catalog.Profiles["smokeTests"]
		profile.Clients[0]["binaryEncoding"] = "enabled"
		_, err := resolveFixtureProfile(profile, "", "", "")
		require.ErrorContains(t, err, "only sdkDefault client binary encoding")
	})
}

func TestValidateEmulatorCapabilitiesChecksOnlySelectedRequirements(t *testing.T) {
	backend := internale2e.Backend{
		Applicability: internale2e.ApplicabilityRequired,
		Fidelity:      "full",
		Requires:      []string{"item"},
	}
	var capabilities emulatorCapabilities
	capabilities.APIVersion = 1
	capabilities.Protocols.GatewayV1 = true
	capabilities.DataPlane = []string{"item"}

	require.NoError(t, validateEmulatorCapabilities(backend, capabilities))

	backend.Requires = []string{"transactionalBatch"}
	require.ErrorContains(t, validateEmulatorCapabilities(backend, capabilities), `required capability "transactionalBatch"`)
}

func TestValidateEmulatorCapabilitiesRejectsMismatchedBackend(t *testing.T) {
	var capabilities emulatorCapabilities
	capabilities.APIVersion = 1

	err := validateEmulatorCapabilities(internale2e.Backend{
		Applicability: internale2e.ApplicabilityRequired,
		Fidelity:      "full",
	}, capabilities)
	require.ErrorContains(t, err, "hostedEmulatorGatewayV1 backend is not available")
}

func TestValidateEmulatorCapabilitiesAcceptsSimulatedFidelity(t *testing.T) {
	var capabilities emulatorCapabilities
	capabilities.APIVersion = 1
	capabilities.Protocols.GatewayV1 = true
	capabilities.DataPlane = []string{"item"}
	backend := internale2e.Backend{
		Applicability: internale2e.ApplicabilitySimulated,
		Fidelity:      "simulated",
		Requires:      []string{"item"},
	}
	require.NoError(t, validateEmulatorCapabilities(backend, capabilities))
	backend.Fidelity = "full"
	require.ErrorContains(t, validateEmulatorCapabilities(backend, capabilities), "does not match applicability")
}

func TestSelectedProfileUsesRequestedGateway(t *testing.T) {
	t.Setenv("AZURE_COSMOS_E2E_BACKEND", emulatorGatewayV2)
	profile := selectedProfile(t, "consistency.response-token-capture")
	require.Equal(t, emulatorGatewayV2, profile.backendID)
	require.Contains(t, profile.backend.Requires, "gatewayV2")
	var capabilities emulatorCapabilities
	capabilities.APIVersion = 1
	capabilities.Protocols.GatewayV1 = true
	capabilities.DataPlane = []string{"item", "query"}
	require.ErrorContains(t, validateEmulatorCapabilities(profile.backend, capabilities), `required capability "gatewayV2"`)
	capabilities.Protocols.GatewayV2 = true
	require.NoError(t, validateEmulatorCapabilities(profile.backend, capabilities))
}

func TestValidateEmulatorAccountMatchesSelectedProfile(t *testing.T) {
	profile, err := resolveFixtureProfile(loadFixtureTestCatalog(t).Profiles["configurationResilience"], "", "", "")
	require.NoError(t, err)
	endpoint := "http://127.0.0.1:8081/"
	account := emulatorAccount{
		ID:          profile.accountID,
		WriteMode:   "single",
		Consistency: "Session",
		Regions: []emulatorAccountRegion{
			{Name: "East US", GatewayEndpoint: endpoint},
			{Name: "West US", GatewayEndpoint: "http://127.0.0.1:8082/"},
		},
	}

	require.NoError(t, validateEmulatorAccount(profile, endpoint, account))

	account.Consistency = "Eventual"
	require.ErrorContains(t, validateEmulatorAccount(profile, endpoint, account), "does not match profile")
	account.Consistency = "Session"

	account.Regions = account.Regions[:1]
	require.ErrorContains(t, validateEmulatorAccount(profile, endpoint, account), "regions")
}

func TestValidateEmulatorBackend(t *testing.T) {
	require.NoError(t, validateEmulatorBackend(""))
	require.NoError(t, validateEmulatorBackend(emulatorGatewayV1))
	require.NoError(t, validateEmulatorBackend(emulatorGatewayV2))
	require.ErrorContains(t, validateEmulatorBackend("unknown"), "not \"unknown\"")
}

func loadFixtureTestCatalog(t *testing.T) internale2e.Catalog {
	t.Helper()
	catalog, err := internale2e.LoadCatalog()
	require.NoError(t, err)
	return catalog
}
