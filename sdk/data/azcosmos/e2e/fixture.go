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
// per-scenario reasons): every active scenario here uses only the "sdkDefault"/single-definition
// account, runtime, and client axis of the smokeTests/coreOperations profiles - and, for a
// scenario that declares additional profiles beyond the one it is actually certified under (for
// example item.lifecycle also declares lifecycleConsistencyMatrix/readConsistencyOverrideMatrix),
// Active status here certifies only that one resolved default profile, not every profile the
// scenario declares; see selectedProfile's doc comment for how the default is picked. Go has no
// database/container management API, so fixtures reuse the pre-provisioned emulator-config.json
// fixture (../internal/testdata/emulator-config.json) rather than provisioning isolated
// resources per scenario; tests instead use unique item IDs/partition keys per case to avoid
// interference. Profile axes beyond the single default definition
// (AZURE_COSMOS_E2E_ACCOUNT/RUNTIME/CLIENT) are not yet translatable through the public Go API
// and are rejected as a configuration failure rather than silently ignored.
package e2e

import (
	"context"
	"encoding/json"
	"fmt"
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
// not every scenario applies to "smokeTests" (several apply only to "coreOperations"). Either way
// this fails the test configuration - not silently skips - if the resolved profile is unknown to
// the pinned catalog or is not one scenarioID actually declares.
func selectedProfile(t *testing.T, scenarioID string) string {
	t.Helper()

	catalog, err := internale2e.LoadCatalog()
	require.NoError(t, err)
	scenario, known := catalog.Scenarios[scenarioID]
	require.Truef(t, known, "scenario %q is not in the pinned catalog", scenarioID)
	require.NotEmptyf(t, scenario.Profiles, "scenario %q declares no profiles", scenarioID)

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
	require.Truef(t, found, "scenario %q does not declare profile %q; selected profiles were %v", scenarioID, profileID, scenario.Profiles)

	// Every active scenario in this pass only has a single ("sdkDefault"-equivalent) account,
	// runtime, and client definition in its applicable profiles. Go cannot yet translate any
	// other axis value through the public API, so reject one explicitly rather than silently
	// running against defaults while claiming to honor the override.
	for _, axis := range []string{"AZURE_COSMOS_E2E_ACCOUNT", "AZURE_COSMOS_E2E_RUNTIME", "AZURE_COSMOS_E2E_CLIENT"} {
		if v := os.Getenv(axis); v != "" {
			t.Fatalf("%s=%q is not supported: this Go port only implements each active scenario's single default axis definition", axis, v)
		}
	}

	// The profile itself - not just an env override - can also select a non-default account,
	// runtime, or client axis (for example configurationResilience's two-region, 750ms-delay
	// account). This fixture only ever constructs the single pre-provisioned sessionSingleRegion
	// account with default client/runtime options, so fail loudly here instead of silently
	// running that mismatched infrastructure while a scenario claims the stricter profile.
	requireSingleAxisID(t, profileID, "account", profile.Accounts, supportedAccountID)
	requireSingleAxisID(t, profileID, "runtime", profile.Runtimes, supportedRuntimeID)
	requireSingleAxisID(t, profileID, "client", profile.Clients, supportedClientID)

	return profileID
}

// The only account/runtime/client axis definitions this fixture actually realizes: the
// pre-provisioned emulator account constructed with default azcosmos/v2 client options (see
// ../internal/testdata/emulator-config.json).
//
// That underlying account does not literally match the pinned catalog's sessionSingleRegion
// definition: the catalog declares exactly one region (East US, zero replication delay), but
// emulator-config.json is shared, non-e2e-specific test infrastructure this port does not own or
// vary per scenario, and it always provisions two region endpoints (East US, West US; zero
// replication delay). This fixture never sets PreferredRegions, so every client built here only
// ever talks to the first-listed (East US) endpoint - observably single-region behavior for the
// CRUD/query scenarios this pass certifies - but that is an accepted approximation of
// sessionSingleRegion, not a byte-for-byte match. A scenario whose contract depends on the
// account genuinely having only one reachable region endpoint (rather than one endpoint this
// fixture happens to use) would need isolated per-scenario provisioning this fixture does not
// yet have.
const (
	supportedAccountID = "sessionSingleRegion"
	supportedRuntimeID = "sdkDefault"
	supportedClientID  = "sdkDefault"
)

// requireSingleAxisID fails the test configuration unless axis declares exactly one definition
// and it is wantID, matching the single account/runtime/client shape this fixture builds.
func requireSingleAxisID(t *testing.T, profileID, axisName string, axis []map[string]any, wantID string) {
	t.Helper()
	require.Lenf(t, axis, 1, "profile %q declares %d %s definitions; this fixture only realizes a single default definition", profileID, len(axis), axisName)
	gotID, _ := axis[0]["id"].(string)
	if gotID != wantID {
		t.Fatalf("profile %q selects %s %q, which this fixture cannot construct through the public azcosmos/v2 API; only %q is realized by the pinned emulator config, so scenarios needing this axis must be reclassified as blocked rather than run against mismatched infrastructure", profileID, axisName, gotID, wantID)
	}
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

// fixture is one scenario's isolated handle to the shared emulator-backed container.
type fixture struct {
	Client    *azcosmos.Client
	Container *azcosmos.ContainerClient
}

// newFixture validates the selected profile for scenarioID, then builds a client/container pair.
// containerIDOverride lets hierarchical-partition-key scenarios select "query-hierarchical"
// instead of the default "items" (hash partition key "/pk") container.
func newFixture(t *testing.T, scenarioID string, containerIDOverride string) fixture {
	t.Helper()
	selectedProfile(t, scenarioID)
	return newRawFixture(t, containerIDOverride)
}

// newRawFixture builds a client/container pair against the pinned single-region emulator account
// without validating any catalog scenario ID or profile. Use this only for a supplementary
// regression test that is not registered as scenario certification evidence in
// internal/e2e.Implementations (for example, TestSessionTokenExplicitManagement, whose declared
// scenario profile needs infrastructure this fixture cannot build); every scenario-certifying
// test must go through newFixture instead so its profile is validated.
func newRawFixture(t *testing.T, containerIDOverride string) fixture {
	t.Helper()
	endpoint, databaseID, containerID := requireEmulator(t)
	if containerIDOverride != "" {
		containerID = containerIDOverride
	}

	cred, err := azcosmos.NewKeyCredential(emulatorKey)
	require.NoError(t, err)
	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	require.NoError(t, client.Initialize(context.Background()))

	container, err := client.NewContainer(databaseID, containerID)
	require.NoError(t, err)

	return fixture{Client: client, Container: container}
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

func mustMarshalItem(t *testing.T, value Item) []byte {
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

func jsonMarshalTenantItem(id, pk, child string, value int64) ([]byte, error) {
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
	require.NotEmpty(t, diagnostics.JSON, "diagnostics must render a JSON summary")
}
