// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// Package e2e ports the shared Cosmos SDK functional E2E certification scenarios
// (Azure/azure-sdk-for-rust's sdk/cosmos/e2e_tests) into source-native Go tests against the
// public azcosmos/v2 API. See ../internal/e2e for the vendored catalog, the Go implementation
// map, and strict validation; see catalog_test.go in this package for how that map is checked
// against the Go test functions actually compiled for this build.
//
// Scope and known limitations of this port (see internal/e2e's implementation map for the exact
// per-scenario reasons): fixtures translate the sessionSingleRegion and sessionTwoRegionsDelayed
// account profiles, sdkDefault runtime/client settings, and proximity/preferred-region routing.
// Active status certifies only the selected profile, not every profile a scenario declares.
// Gateway V2 selection is supported only for the session-staleness, response-token-capture, and
// capabilities scenarios. Other axis definitions, including dynamicTopology, fail configuration
// validation rather than silently running against a different emulator setup.
package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	internale2e "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2/internal/e2e"
	"github.com/stretchr/testify/require"
)

// emulatorKey is the Azure Cosmos DB Emulator's well-known, publicly documented default key (see
// https://learn.microsoft.com/azure/cosmos-db/emulator); it is not a secret.
const emulatorKey = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

// Item is the schema-agnostic test document shape every scenario uses. Callers own item
// serialization in azcosmos/v2, so this is ordinary application code, not an SDK type.
type Item struct {
	ID    string `json:"id"`
	PK    string `json:"pk"`
	Child string `json:"child,omitempty"`
	Value int64  `json:"value"`
	Score *int64 `json:"score,omitempty"`
}

// selectedProfile validates and returns the profile to run scenarioID under. AZURE_COSMOS_E2E_PROFILE
// selects it explicitly; unset, it defaults to the scenario's own first declared profile, since
// not every scenario applies to "smokeTests" (several apply only to "coreOperations"). An unknown
// profile, or applicable-but-unsupported setup, still fails the test configuration outright. But
// an explicitly-selected profile that scenarioID simply does not declare is a valid profile
// selection per the shared catalog contract (running the whole package
// with AZURE_COSMOS_E2E_PROFILE=smokeTests is expected to skip, not fail, every coreOperations-only
// scenario, and vice versa), so that case skips the test instead of failing it.
func selectedProfile(t *testing.T, scenarioID string) fixtureProfile {
	t.Helper()

	catalog, err := internale2e.LoadCatalog()
	require.NoError(t, err)
	scenario, known := catalog.Scenarios[scenarioID]
	require.Truef(t, known, "scenario %q is not in the pinned catalog", scenarioID)
	require.NotEmptyf(t, scenario.Profiles, "scenario %q declares no profiles", scenarioID)

	explicitProfile := os.Getenv("AZURE_COSMOS_E2E_PROFILE") != ""
	profileID := os.Getenv("AZURE_COSMOS_E2E_PROFILE")
	if profileID == "" {
		profileID = scenario.Profiles[0]
	}

	profile, known := catalog.Profiles[profileID]
	require.Truef(t, known, "AZURE_COSMOS_E2E_PROFILE=%q is not a profile in the pinned catalog", profileID)

	found := false
	for _, p := range scenario.Profiles {
		if p == profileID {
			found = true
			break
		}
	}
	if !found {
		if explicitProfile {
			t.Skipf("scenario %q does not declare profile %q; selected profiles were %v", scenarioID, profileID, scenario.Profiles)
		}
		// Unreachable in practice: the default profileID is always scenario.Profiles[0], so this
		// only guards a future bug in that default-selection logic.
		t.Fatalf("scenario %q does not declare its own default profile %q; selected profiles were %v", scenarioID, profileID, scenario.Profiles)
	}

	resolved, err := resolveFixtureProfile(
		profile,
		os.Getenv("AZURE_COSMOS_E2E_ACCOUNT"),
		os.Getenv("AZURE_COSMOS_E2E_RUNTIME"),
		os.Getenv("AZURE_COSMOS_E2E_CLIENT"),
	)
	require.NoErrorf(t, err, "profile %q is not supported by the Go emulator fixture", profileID)
	backend := os.Getenv("AZURE_COSMOS_E2E_BACKEND")
	if backend == emulatorGatewayV2 {
		switch scenarioID {
		case "consistency.session-staleness", "consistency.response-token-capture", "management.capabilities":
		default:
			t.Skipf("scenario %q has not been verified with %s", scenarioID, backend)
		}
		resolved.backend = scenario.Backends.HostedEmulatorGatewayV2
		resolved.backendID = emulatorGatewayV2
	} else {
		resolved.backend = scenario.Backends.HostedEmulatorGatewayV1
		resolved.backendID = emulatorGatewayV1
	}
	resolved.scenarioID = scenarioID
	return resolved
}

const (
	defaultEndpointEnv                 = "AZCOSMOS_E2E_SINGLE_REGION_ENDPOINT"
	defaultManagementEndpointEnv       = "AZCOSMOS_E2E_SINGLE_REGION_MANAGEMENT_ENDPOINT"
	resilienceEndpointEnv              = "AZCOSMOS_CONFIGURATION_RESILIENCE_ENDPOINT"
	resilienceManagementEndpointEnv    = "AZCOSMOS_CONFIGURATION_RESILIENCE_MANAGEMENT_ENDPOINT"
	supportedSingleRegionAccountID     = "sessionSingleRegion"
	supportedTwoRegionDelayedAccountID = "sessionTwoRegionsDelayed"
	supportedDefaultRuntimeID          = "sdkDefault"
	supportedDefaultClientID           = "sdkDefault"
	emulatorGatewayV1                  = "hostedEmulatorGatewayV1"
	emulatorGatewayV2                  = "hostedEmulatorGatewayV2"
)

type fixtureProfile struct {
	profileID             string
	accountID             string
	regions               []string
	runtimeID             string
	clientID              string
	endpointEnv           string
	managementEndpointEnv string
	clientOptions         *azcosmos.ClientOptions
	backend               internale2e.Backend
	backendID             string
	scenarioID            string
}

func resolveFixtureProfile(
	profile internale2e.Profile,
	accountOverride string,
	runtimeOverride string,
	clientOverride string,
) (fixtureProfile, error) {
	account, accountID, err := selectProfileDefinition(profile.ID, "account", profile.Accounts, accountOverride)
	if err != nil {
		return fixtureProfile{}, err
	}
	accountSetup, err := resolveFixtureAccount(account, accountID)
	if err != nil {
		return fixtureProfile{}, err
	}

	runtime, runtimeID, err := selectProfileDefinition(profile.ID, "runtime", profile.Runtimes, runtimeOverride)
	if err != nil {
		return fixtureProfile{}, err
	}
	if err := validateFixtureRuntime(runtime); err != nil {
		return fixtureProfile{}, fmt.Errorf("runtime %q: %w", runtimeID, err)
	}

	client, clientID, err := selectProfileDefinition(profile.ID, "client", profile.Clients, clientOverride)
	if err != nil {
		return fixtureProfile{}, err
	}
	routing, err := resolveFixtureRouting(client, accountSetup.regions)
	if err != nil {
		return fixtureProfile{}, fmt.Errorf("client %q: %w", clientID, err)
	}

	return fixtureProfile{
		profileID:             profile.ID,
		accountID:             accountID,
		regions:               accountSetup.regionNames,
		runtimeID:             runtimeID,
		clientID:              clientID,
		endpointEnv:           accountSetup.endpointEnv,
		managementEndpointEnv: accountSetup.managementEndpointEnv,
		clientOptions:         &azcosmos.ClientOptions{Routing: routing},
	}, nil
}

func selectProfileDefinition(profileID, axis string, definitions []map[string]any, override string) (map[string]any, string, error) {
	if len(definitions) == 0 {
		return nil, "", fmt.Errorf("profile %q declares no %s definitions", profileID, axis)
	}
	if override == "" {
		if len(definitions) != 1 {
			return nil, "", fmt.Errorf("profile %q declares %d %s definitions; set AZURE_COSMOS_E2E_%s to select one", profileID, len(definitions), axis, strings.ToUpper(axis))
		}
		definition := definitions[0]
		id, ok := definition["id"].(string)
		if !ok || id == "" {
			return nil, "", fmt.Errorf("profile %q has a %s definition without a string id", profileID, axis)
		}
		return definition, id, nil
	}

	var selected map[string]any
	for _, definition := range definitions {
		id, ok := definition["id"].(string)
		if !ok || id == "" {
			return nil, "", fmt.Errorf("profile %q has a %s definition without a string id", profileID, axis)
		}
		if id == override {
			if selected != nil {
				return nil, "", fmt.Errorf("profile %q declares duplicate %s id %q", profileID, axis, override)
			}
			selected = definition
		}
	}
	if selected == nil {
		return nil, "", fmt.Errorf("AZURE_COSMOS_E2E_%s=%q does not select a definition in the %s axis of profile %q", strings.ToUpper(axis), override, axis, profileID)
	}
	return selected, override, nil
}

type fixtureAccountSetup struct {
	regions               []azcosmos.Region
	regionNames           []string
	endpointEnv           string
	managementEndpointEnv string
}

func resolveFixtureAccount(account map[string]any, id string) (fixtureAccountSetup, error) {
	if err := requireDefinitionKeys("account", account, "id", "writeMode", "consistency", "regions", "replication", "perPartitionFailover"); err != nil {
		return fixtureAccountSetup{}, err
	}
	if account["writeMode"] != "single" || account["consistency"] != "session" || account["perPartitionFailover"] != false {
		return fixtureAccountSetup{}, fmt.Errorf("account %q requires unsupported write, consistency, or failover settings", id)
	}
	regionsValue, ok := account["regions"].([]any)
	if !ok {
		return fixtureAccountSetup{}, fmt.Errorf("account %q has invalid regions", id)
	}
	regionNames := make([]string, 0, len(regionsValue))
	for _, value := range regionsValue {
		region, ok := value.(map[string]any)
		if !ok || len(region) != 1 {
			return fixtureAccountSetup{}, fmt.Errorf("account %q has an unsupported region definition", id)
		}
		name, ok := region["name"].(string)
		if !ok {
			return fixtureAccountSetup{}, fmt.Errorf("account %q has a region without a string name", id)
		}
		regionNames = append(regionNames, name)
	}

	replication, ok := account["replication"].(map[string]any)
	if !ok || len(replication) != 2 {
		return fixtureAccountSetup{}, fmt.Errorf("account %q has an unsupported replication definition", id)
	}
	minDelay, minOK := numberValue(replication["minDelayMs"])
	maxDelay, maxOK := numberValue(replication["maxDelayMs"])
	if !minOK || !maxOK {
		return fixtureAccountSetup{}, fmt.Errorf("account %q has invalid replication delays", id)
	}

	switch id {
	case supportedSingleRegionAccountID:
		if !sameStrings(regionNames, []string{"East US"}) || minDelay != 0 || maxDelay != 0 {
			return fixtureAccountSetup{}, fmt.Errorf("account %q must be one East US region with zero replication delay", id)
		}
		return fixtureAccountSetup{
			regions:               []azcosmos.Region{azcosmos.RegionEastUS},
			regionNames:           regionNames,
			endpointEnv:           defaultEndpointEnv,
			managementEndpointEnv: defaultManagementEndpointEnv,
		}, nil
	case supportedTwoRegionDelayedAccountID:
		if !sameStrings(regionNames, []string{"East US", "West US"}) || minDelay != 750 || maxDelay != 750 {
			return fixtureAccountSetup{}, fmt.Errorf("account %q must be East US and West US with a fixed 750ms replication delay", id)
		}
		return fixtureAccountSetup{
			regions:               []azcosmos.Region{azcosmos.RegionEastUS, azcosmos.RegionWestUS},
			regionNames:           regionNames,
			endpointEnv:           resilienceEndpointEnv,
			managementEndpointEnv: resilienceManagementEndpointEnv,
		}, nil
	default:
		return fixtureAccountSetup{}, fmt.Errorf("account %q has no matching emulator configuration", id)
	}
}

func validateFixtureRuntime(runtime map[string]any) error {
	if err := requireDefinitionKeys("runtime", runtime, "id", "gatewayV2", "ppcb"); err != nil {
		return err
	}
	if runtime["id"] != supportedDefaultRuntimeID || runtime["gatewayV2"] != "backendDefault" || runtime["ppcb"] != "sdkDefault" {
		return fmt.Errorf("only sdkDefault runtime settings are supported")
	}
	return nil
}

func resolveFixtureRouting(client map[string]any, regions []azcosmos.Region) (azcosmos.RoutingStrategy, error) {
	if err := requireDefinitionKeys("client", client, "id", "binaryEncoding", "routing"); err != nil {
		return azcosmos.RoutingStrategy{}, err
	}
	if client["id"] != supportedDefaultClientID || client["binaryEncoding"] != "sdkDefault" {
		return azcosmos.RoutingStrategy{}, fmt.Errorf("only sdkDefault client binary encoding is supported")
	}
	switch client["routing"] {
	case "proximity":
		return azcosmos.ProximityTo(azcosmos.RegionEastUS), nil
	case "preferredRegions":
		return azcosmos.PreferredRegions(regions...), nil
	case "accountOrder":
		return azcosmos.RoutingStrategy{}, nil
	default:
		return azcosmos.RoutingStrategy{}, fmt.Errorf("routing %q is unsupported", client["routing"])
	}
}

func requireDefinitionKeys(axis string, definition map[string]any, keys ...string) error {
	if len(definition) != len(keys) {
		return fmt.Errorf("%s definition has unsupported fields", axis)
	}
	for _, key := range keys {
		if _, ok := definition[key]; !ok {
			return fmt.Errorf("%s definition is missing %q", axis, key)
		}
	}
	return nil
}

func numberValue(value any) (float64, bool) {
	number, ok := value.(float64)
	return number, ok
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

// requireEmulator skips (not fails) when EMULATOR is unset, matching the rest of the package's
// emulator-gated tests; a missing emulator is an environment precondition, not a scenario defect.
func requireEmulator(t *testing.T) (endpoint, databaseID, containerID string) {
	t.Helper()
	if os.Getenv("EMULATOR") == "" {
		t.Skip("set EMULATOR to run E2E certification scenarios against the Cosmos DB emulator")
	}
	endpoint = os.Getenv("AZCOSMOS_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://localhost:8081/"
	}
	databaseID = os.Getenv("AZCOSMOS_DATABASE")
	if databaseID == "" {
		databaseID = "itemdb"
	}
	containerID = os.Getenv("AZCOSMOS_CONTAINER")
	if containerID == "" {
		containerID = "items"
	}
	return
}

func requireProfileEmulator(t *testing.T, profile fixtureProfile) (endpoint, managementEndpoint, databaseID, containerID string) {
	t.Helper()
	if os.Getenv("EMULATOR") == "" {
		t.Skip("set EMULATOR to run E2E certification scenarios against the Cosmos DB emulator")
	}
	endpoint = os.Getenv(profile.endpointEnv)
	if endpoint == "" {
		t.Fatalf("%s must be set for profile %q; refusing to use a different account configuration", profile.endpointEnv, profile.profileID)
	}
	managementEndpoint = os.Getenv(profile.managementEndpointEnv)
	if managementEndpoint == "" {
		t.Fatalf("%s must be set to probe the selected emulator profile", profile.managementEndpointEnv)
	}
	databaseID = os.Getenv("AZCOSMOS_DATABASE")
	if databaseID == "" {
		databaseID = "itemdb"
	}
	containerID = os.Getenv("AZCOSMOS_CONTAINER")
	if containerID == "" {
		containerID = "items"
	}
	return
}

func validateEmulatorBackend(backend string) error {
	if backend != "" && backend != emulatorGatewayV1 && backend != emulatorGatewayV2 {
		return fmt.Errorf("the Go emulator fixture supports %q and %q, not %q", emulatorGatewayV1, emulatorGatewayV2, backend)
	}
	return nil
}

type emulatorCapabilities struct {
	APIVersion      int    `json:"apiVersion"`
	EmulatorVersion string `json:"emulatorVersion"`
	Protocols       struct {
		GatewayV1 bool `json:"gatewayV1"`
		GatewayV2 bool `json:"gatewayV2"`
	} `json:"protocols"`
	DataPlane         []string `json:"dataPlane"`
	ManagementActions []string `json:"managementActions"`
	Limitations       []string `json:"limitations"`
}

func probeEmulatorCapabilities(t *testing.T, profile fixtureProfile, endpoint, managementEndpoint string) {
	t.Helper()
	var capabilities emulatorCapabilities
	require.NoError(t, readManagementDocument(t, managementEndpoint, "capabilities", &capabilities))
	require.NoError(t, validateEmulatorCapabilities(profile.backend, capabilities))

	var account emulatorAccount
	require.NoError(t, readManagementDocument(t, managementEndpoint, "account", &account))
	require.NoError(t, validateEmulatorAccount(profile, endpoint, account))
}

type emulatorGatewayRequests struct {
	Gateway   uint64 `json:"gatewayRequests"`
	GatewayV2 uint64 `json:"gateway20Requests"`
}

func readGatewayRequests(t *testing.T, managementEndpoint string) emulatorGatewayRequests {
	t.Helper()
	var requests emulatorGatewayRequests
	require.NoError(t, readManagementDocument(t, managementEndpoint, "health", &requests))
	return requests
}

func assertGatewayRequests(t *testing.T, profile fixtureProfile, before, after emulatorGatewayRequests) {
	t.Helper()
	if profile.backendID == emulatorGatewayV2 {
		require.Greater(t, after.GatewayV2, before.GatewayV2, "scenario must use Gateway V2")
	} else {
		require.Greater(t, after.Gateway, before.Gateway, "scenario must use Gateway V1")
		require.Equal(t, before.GatewayV2, after.GatewayV2, "Gateway V1 scenario must not use Gateway V2")
	}
}

func readManagementDocument(t *testing.T, managementEndpoint, path string, target any) error {
	t.Helper()
	endpoint, err := url.JoinPath(managementEndpoint, path)
	if err != nil {
		return fmt.Errorf("invalid emulator management endpoint %q: %w", managementEndpoint, err)
	}
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("creating request for emulator management endpoint %q: %w", endpoint, err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("reading emulator management endpoint %q: %w", endpoint, err)
	}
	defer func() { require.NoError(t, response.Body.Close()) }()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("emulator management endpoint %q returned %s", endpoint, response.Status)
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("decoding emulator management endpoint %q: %w", endpoint, err)
	}
	return nil
}

type emulatorAccount struct {
	ID                   string                  `json:"id"`
	WriteMode            string                  `json:"writeMode"`
	Consistency          string                  `json:"consistency"`
	PerPartitionFailover bool                    `json:"perPartitionFailover"`
	Regions              []emulatorAccountRegion `json:"regions"`
}

type emulatorAccountRegion struct {
	Name            string `json:"name"`
	GatewayEndpoint string `json:"gatewayEndpoint"`
}

func validateEmulatorAccount(profile fixtureProfile, endpoint string, account emulatorAccount) error {
	if account.ID != profile.accountID || account.WriteMode != "single" || !strings.EqualFold(account.Consistency, "session") || account.PerPartitionFailover {
		return fmt.Errorf("emulator account does not match profile %q", profile.profileID)
	}
	regionNames := make([]string, len(account.Regions))
	for i, region := range account.Regions {
		regionNames[i] = region.Name
	}
	if !sameStrings(regionNames, profile.regions) {
		return fmt.Errorf("emulator account regions %v do not match profile %q regions %v", regionNames, profile.profileID, profile.regions)
	}
	if len(account.Regions) == 0 || strings.TrimRight(account.Regions[0].GatewayEndpoint, "/") != strings.TrimRight(endpoint, "/") {
		return fmt.Errorf("emulator endpoint %q does not match the first region of profile %q", endpoint, profile.profileID)
	}
	return nil
}

func validateEmulatorCapabilities(backend internale2e.Backend, capabilities emulatorCapabilities) error {
	if capabilities.APIVersion != 1 {
		return fmt.Errorf("unsupported emulator capabilities API version %d", capabilities.APIVersion)
	}
	if backend.Applicability != internale2e.ApplicabilityRequired && backend.Applicability != internale2e.ApplicabilitySupported && backend.Applicability != internale2e.ApplicabilitySimulated {
		return fmt.Errorf("selected hosted emulator backend applicability is %q", backend.Applicability)
	}
	if (backend.Applicability == internale2e.ApplicabilitySimulated && backend.Fidelity != "simulated") ||
		(backend.Applicability != internale2e.ApplicabilitySimulated && backend.Fidelity != "full") {
		return fmt.Errorf("selected hosted emulator backend fidelity %q does not match applicability %q", backend.Fidelity, backend.Applicability)
	}
	if !capabilities.Protocols.GatewayV1 {
		return fmt.Errorf("selected hostedEmulatorGatewayV1 backend is not available")
	}
	available := make(map[string]bool, len(capabilities.DataPlane)+len(capabilities.ManagementActions)+2)
	// A successfully decoded versioned document establishes the capabilities endpoint itself.
	available["capabilities"] = true
	for _, capability := range capabilities.DataPlane {
		available[capability] = true
	}
	for _, capability := range capabilities.ManagementActions {
		available[capability] = true
	}
	if capabilities.Protocols.GatewayV1 {
		available["gatewayV1"] = true
	}
	if capabilities.Protocols.GatewayV2 {
		available["gatewayV2"] = true
	}
	for _, required := range backend.Requires {
		if !available[required] {
			return fmt.Errorf("selected hosted emulator backend does not provide required capability %q", required)
		}
	}
	return nil
}

// fixture is one scenario's isolated handle to the shared emulator-backed container.
type fixture struct {
	Client    *azcosmos.Client
	Database  *azcosmos.DatabaseClient
	Container *azcosmos.ContainerClient
}

// newFixture validates the selected profile for scenarioID, then builds a client/container pair.
// containerIDOverride lets hierarchical-partition-key scenarios select "query-hierarchical"
// instead of the default "items" (hash partition key "/pk") container.
func newFixture(t *testing.T, scenarioID string, containerIDOverride string) fixture {
	t.Helper()
	profile := selectedProfile(t, scenarioID)
	endpoint, managementEndpoint, databaseID, containerID := requireProfileEmulator(t, profile)
	require.NoErrorf(t, validateEmulatorBackend(os.Getenv("AZURE_COSMOS_E2E_BACKEND")), "scenario %q", scenarioID)
	probeEmulatorCapabilities(t, profile, endpoint, managementEndpoint)
	return newFixtureAt(t, endpoint, databaseID, containerID, containerIDOverride, profile.clientOptions)
}

// newRawFixture builds a client/container pair against the default single-region emulator account
// without validating any catalog scenario ID or profile. Use this only for a supplementary
// regression test that is not registered as scenario certification evidence in
// internal/e2e.Implementations; every scenario-certifying test must validate its profile instead.
func newRawFixture(t *testing.T, containerIDOverride string) fixture {
	t.Helper()
	endpoint, databaseID, containerID := requireEmulator(t)
	return newFixtureAt(t, endpoint, databaseID, containerID, containerIDOverride, nil)
}

func newFixtureAt(t *testing.T, endpoint, databaseID, containerID, containerIDOverride string, options *azcosmos.ClientOptions) fixture {
	t.Helper()
	if containerIDOverride != "" {
		containerID = containerIDOverride
	}

	cred, err := azcosmos.NewKeyCredential(emulatorKey)
	require.NoError(t, err)
	client, err := azcosmos.NewClientWithKey(t.Context(), endpoint, cred, options)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	database, err := client.NewDatabase(databaseID)
	require.NoError(t, err)
	container, err := client.NewContainer(t.Context(), databaseID, containerID)
	require.NoError(t, err)

	return fixture{Client: client, Database: database, Container: container}
}

func createTestContainer(t *testing.T, database *azcosmos.DatabaseClient, properties azcosmos.ContainerProperties) (*azcosmos.ContainerClient, azcosmos.ContainerResponse) {
	t.Helper()
	if properties.ID == "" {
		properties.ID = uniqueID(t) + "-container"
	}
	response, err := database.CreateContainer(t.Context(), properties, nil)
	require.NoError(t, err)
	container, err := database.NewContainer(t.Context(), properties.ID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := container.Delete(context.Background(), nil)
		if err == nil {
			return
		}
		var cosmosErr *azcosmos.Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, azcosmos.CodeNotFound, cosmosErr.Code)
	})
	return container, response
}

// item builds the schema-agnostic test document used across scenarios.
func item(id, pk string, value int64) Item {
	return Item{ID: id, PK: pk, Value: value}
}

// itemWithScore builds an Item with the additional "score" property the parameterized-filter
// query scenario orders and filters on.
func itemWithScore(id, pk string, score int64) Item {
	return Item{ID: id, PK: pk, Value: score, Score: &score}
}

func mustMarshalItem(t *testing.T, value Item) json.RawMessage {
	t.Helper()
	body, err := json.Marshal(value)
	require.NoError(t, err)
	return body
}

func mustUnmarshalItem(t *testing.T, body []byte) Item {
	t.Helper()
	var value Item
	require.NoError(t, json.Unmarshal(body, &value))
	return value
}

// mustUnmarshalInto is the untyped counterpart used by scenarios that decode into a map or a
// different document shape (TenantItem, a plain id string, ...).
func mustUnmarshalInto(t *testing.T, body []byte, target any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(body, target))
}

// TenantItem is the schema-agnostic hierarchical-partition-key test document shape, used by
// scenarios against the pre-provisioned "query-hierarchical" container (partition key paths
// "/pk", "/child").
type TenantItem struct {
	ID    string `json:"id"`
	PK    string `json:"pk"`
	Child string `json:"child"`
	Value int64  `json:"value"`
}

func tenantItem(id, pk, child string, value int64) TenantItem {
	return TenantItem{ID: id, PK: pk, Child: child, Value: value}
}

func jsonMarshalTenantItem(id, pk, child string, value int64) (json.RawMessage, error) {
	return json.Marshal(tenantItem(id, pk, child, value))
}

// uniqueID keeps scenarios that share the pre-provisioned container from colliding with each
// other or with a previous run, mirroring emulator_test.go's uniqueItemID in package azcosmos.
// to returns a pointer to v, for the *bool-valued options fields (e.g. EnableContentResponseOnWrite)
// that default to "unset" rather than false.
func to[T any](v T) *T {
	return &v
}

func uniqueID(t *testing.T) string {
	t.Helper()
	name := strings.NewReplacer("/", "-", "\\", "-", "?", "-", "#", "-").Replace(t.Name())
	return fmt.Sprintf("%s-%d", name, time.Now().UnixNano())
}

// trackItem deletes the item at cleanup, tolerating it already being gone (deletion is itself
// under test in some scenarios).
func trackItem(t *testing.T, container *azcosmos.ContainerClient, pk azcosmos.PartitionKey, id string) {
	t.Helper()
	t.Cleanup(func() {
		_, err := container.DeleteItem(context.Background(), pk, id, nil)
		if err == nil {
			return
		}
		var cosmosErr *azcosmos.Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, azcosmos.CodeNotFound, cosmosErr.Code)
	})
}

// requireCode asserts err is an *azcosmos.Error with the given Code, failing with the actual
// error otherwise so a mismatch is diagnosable without a second run.
func requireCode(t *testing.T, err error, code azcosmos.Code) *azcosmos.Error {
	t.Helper()
	var cosmosErr *azcosmos.Error
	require.ErrorAsf(t, err, &cosmosErr, "expected an *azcosmos.Error with Code %s, got: %v", code, err)
	require.Equalf(t, code, cosmosErr.Code, "unexpected error: %v", cosmosErr)
	return cosmosErr
}

// requireCriticalDiagnostics checks the diagnostics fields the shared scenarios treat as
// contractual (status/sub-status presence and a populated attempt count), without asserting
// exact diagnostic JSON text, timings, or RU values, per this port's assertion boundary.
func requireCriticalDiagnostics(t *testing.T, diagnostics *azcosmos.Diagnostics, expectStatusCode int) {
	t.Helper()
	require.NotNilf(t, diagnostics, "operation must carry diagnostics")
	require.Equal(t, expectStatusCode, diagnostics.StatusCode)
	require.NotZero(t, diagnostics.AttemptCount, "diagnostics must record at least one attempt")
	rendered, err := diagnostics.JSON()
	require.NoError(t, err)
	require.NotEmpty(t, rendered, "diagnostics must render a JSON summary")
}
