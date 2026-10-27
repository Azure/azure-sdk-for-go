// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore Wtle myaccount

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

// testAccountKey is base64, as a real account key is. NewKeyCredential validates that, so a
// placeholder like "key" would be rejected.
const testAccountKey = "dGVzdC1hY2NvdW50LWtleQ=="

// Offline tests construct local handles without claiming public constructor readiness.
func newTestTokenClient(endpoint string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	return newClient(endpoint, "", cred, options)
}

func newTestClientWithKey(endpoint string, cred KeyCredential, options *ClientOptions) (*Client, error) {
	return newClient(endpoint, cred.accountKey, nil, options)
}

func localTestContainer(client *Client, databaseID, containerID string) *ContainerClient {
	return &ContainerClient{id: containerID, database: &DatabaseClient{id: databaseID, client: client}}
}

func mustKeyCredential(t *testing.T) KeyCredential {
	t.Helper()

	cred, err := NewKeyCredential(testAccountKey)
	require.NoError(t, err)
	return cred
}

func newTestClient(t *testing.T) *Client {
	t.Helper()

	client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	return client
}

type fakeTokenCredential struct{}

func (fakeTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     "test-access-token",
		ExpiresOn: time.Now().Add(time.Hour),
	}, nil
}

// The private construction seam lets unit tests configure credentials without network I/O.
func TestLocalClientAcceptsTokenCredential(t *testing.T) {
	client, err := newTestTokenClient("https://myaccount.documents.azure.com", fakeTokenCredential{}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	require.Equal(t, "https://myaccount.documents.azure.com", client.Endpoint())
}

func TestNewClientRejectsNilCredential(t *testing.T) {
	_, err := NewClient(t.Context(), "https://myaccount.documents.azure.com", nil, nil)
	require.Error(t, err)
}

func TestInitializeRejectsNilContext(t *testing.T) {
	client := newTestClient(t)
	err := client.initialize(nil) //nolint:staticcheck // verifies the guard
	require.Error(t, err)
}

func TestInitializeReportsUnavailableDriver(t *testing.T) {
	if driverAvailable {
		t.Skip("the driver is available in this build")
	}
	client, err := NewClientWithKey(t.Context(), "https://myaccount.documents.azure.com", mustKeyCredential(t), nil)
	require.Nil(t, client)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Contains(t, cosmosErr.Message, "cannot reach the Cosmos driver")
}

func TestAccountConstructorsRejectInvalidContext(t *testing.T) {
	for _, constructor := range []func(context.Context) (*Client, error){
		func(ctx context.Context) (*Client, error) {
			return NewClient(ctx, "https://myaccount.documents.azure.com", fakeTokenCredential{}, nil)
		},
		func(ctx context.Context) (*Client, error) {
			return NewClientWithKey(ctx, "https://myaccount.documents.azure.com", mustKeyCredential(t), nil)
		},
	} {
		client, err := constructor(nil)
		require.Nil(t, client)
		require.ErrorContains(t, err, "context")
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		client, err = constructor(ctx)
		require.Nil(t, client)
		require.ErrorIs(t, err, context.Canceled)
	}
}

// The zero value is what a caller gets if they ignore the error from NewKeyCredential.
func TestNewClientWithKeyRejectsZeroCredential(t *testing.T) {
	_, err := NewClientWithKey(t.Context(), "https://myaccount.documents.azure.com", KeyCredential{}, nil)
	require.Error(t, err)
}

func TestNewKeyCredentialRejectsEmptyKey(t *testing.T) {
	_, err := NewKeyCredential("")
	require.Error(t, err)
}

// The driver takes the key as an opaque string and decodes it once per signature, so a malformed
// key that is accepted here would surface as an authentication failure on every operation rather
// than as a configuration error at startup.
func TestNewKeyCredentialRejectsMalformedKey(t *testing.T) {
	for _, tt := range []struct {
		name string
		key  string
	}{
		{"not base64", "not-a-valid-key!"},
		{"invalid character", "dGVzdC1hY2NvdW50*a2V5"},
		{"whitespace only", " "},

		// Go's decoder skips newlines while the driver rejects them, so a key wrapped across
		// lines by a config file has to be caught here or it fails on every request instead.
		{"embedded newline", "dGVzdC1hY2\nNvdW50LWtleQ=="},
		{"newline only", "\n"},
		{"trailing newline", testAccountKey + "\n"},
		{"leading space", " " + testAccountKey},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewKeyCredential(tt.key)
			require.Error(t, err, "key %q should be rejected", tt.key)
		})
	}
}

// The driver's decoder treats padding as optional, so rejecting an unpadded key here would refuse
// one it would have signed with.
func TestNewKeyCredentialAcceptsBase64Key(t *testing.T) {
	for _, tt := range []struct {
		name string
		key  string
	}{
		{"padded", testAccountKey},
		{"unpadded", strings.TrimRight(testAccountKey, "=")},
		{"real key shape", strings.Repeat("A", 86) + "=="},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cred, err := NewKeyCredential(tt.key)
			require.NoError(t, err)

			// The driver wants the key as it was given, not the decoded bytes.
			require.Equal(t, tt.key, cred.key())
		})
	}
}

func TestNewClientRejectsNonAbsoluteEndpoint(t *testing.T) {
	cred := mustKeyCredential(t)

	// "https://:8081" has a non-empty Host (":8081") but no hostname, which is why the check is
	// on Hostname rather than Host.
	for _, endpoint := range []string{"", "myaccount.documents.azure.com", "/relative/path", "https://", "https://:8081"} {
		t.Run(endpoint, func(t *testing.T) {
			_, err := newTestClientWithKey(endpoint, cred, nil)
			require.Error(t, err, "endpoint %q should be rejected", endpoint)
		})
	}
}

func TestNewClientAcceptsAbsoluteEndpoint(t *testing.T) {
	client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	require.Equal(t, "https://myaccount.documents.azure.com", client.Endpoint())
}

// The client must not alias the caller's slice, or a later append by the caller silently changes
// the routing order of a client that has already been created.
func TestNewClientCopiesRoutingRegions(t *testing.T) {
	regions := []Region{RegionWestUS, RegionEastUS}
	client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, &ClientOptions{
		Routing: PreferredRegions(regions...),
	})

	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	regions[0] = RegionNorthEurope
	require.Equal(t, []Region{RegionWestUS, RegionEastUS}, client.options.Routing.preferredRegions)
}

func TestNewClientCopiesContentResponseOption(t *testing.T) {
	enabled := true
	client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, &ClientOptions{
		Operation: OperationOptions{EnableContentResponseOnWrite: &enabled},
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	enabled = false
	require.NotNil(t, client.options.Operation.EnableContentResponseOnWrite)
	require.True(t, *client.options.Operation.EnableContentResponseOnWrite)
}

// Close is documented as idempotent, safe to call concurrently, and as reporting the same result
// to every caller. Storing the teardown error in a local rather than on the client would give it
// to the first caller alone, which is what this exercises under -race.
func TestCloseIsIdempotentAndConcurrencySafe(t *testing.T) {
	client := newTestClient(t)

	const callers = 8
	results := make([]error, callers)
	var wg sync.WaitGroup
	wg.Add(callers)
	for i := range results {
		go func() {
			defer wg.Done()
			results[i] = client.Close()
		}()
	}
	wg.Wait()

	for _, err := range results {
		require.Equal(t, results[0], err, "every caller must observe the same result")
	}
	require.NoError(t, client.Close())
}

func TestNewDatabaseIsLocalAndContainerFactoriesHonorContext(t *testing.T) {
	client := newTestClient(t)

	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	require.Equal(t, "db", database.ID())

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	for _, factory := range []func(context.Context) (*ContainerClient, error){
		func(ctx context.Context) (*ContainerClient, error) { return database.NewContainer(ctx, "items") },
		func(ctx context.Context) (*ContainerClient, error) { return client.NewContainer(ctx, "db", "items") },
	} {
		container, err := factory(ctx)
		require.Nil(t, container)
		require.ErrorIs(t, err, context.Canceled)
		container, err = factory(nil)
		require.Nil(t, container)
		require.ErrorContains(t, err, "context")
	}
	require.NoError(t, client.Close())
	container, err := client.NewContainer(ctx, "db", "items")
	require.Nil(t, container)
	var failure *Error
	require.ErrorAs(t, err, &failure)
	require.Equal(t, CodeClientClosed, failure.Code)
}

func TestNewDatabaseAndNewContainerRejectEmptyIDs(t *testing.T) {
	client := newTestClient(t)

	_, err := client.NewDatabase("")
	require.Error(t, err)

	_, err = client.NewContainer(t.Context(), "db", "")
	require.Error(t, err)

	_, err = client.NewContainer(t.Context(), "", "items")
	require.Error(t, err)
}

// The scheme is validated, because a Cosmos endpoint is always reached over HTTP(S) and anything
// else is a copy-paste mistake that would otherwise surface as an opaque driver failure.
func TestNewClientRejectsNonHTTPSchemes(t *testing.T) {
	cred := mustKeyCredential(t)

	for _, endpoint := range []string{"ftp://myaccount", "file:///tmp/x", "wss://myaccount"} {
		t.Run(endpoint, func(t *testing.T) {
			_, err := newTestClientWithKey(endpoint, cred, nil)
			require.Error(t, err)
		})
	}
}

// The endpoint is named in its error so the message is actionable; this also pins the shape the
// key-leak test above relies on.
func TestEndpointErrorNamesTheEndpoint(t *testing.T) {
	cred := mustKeyCredential(t)

	_, err := newTestClientWithKey("notaurl", cred, nil)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "notaurl"))
}

// After Close the client must refuse work rather than hand a released driver handle to an
// operation.
func TestOperationsAfterCloseAreRejected(t *testing.T) {
	client := newTestClient(t)
	container := localTestContainer(client, "db", "items")
	require.NoError(t, client.Close())

	_, err := container.ReadItem(context.Background(), NewPartitionKeyString("pk"), "item-1", nil)
	var cosmosErr *Error
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)

	_, err = container.CreateItem(context.Background(), NewPartitionKeyString("pk"), "x", json.RawMessage(`{"id":"x"}`), nil)
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
}

// Close waits for in-flight operations before releasing anything. This is the safety-critical part
// of the lifetime contract: an operation still running when the driver's handles are freed is a
// use-after-free in Rust-owned memory, not a nil panic.
//
// Asserting only that operations fail after Close returns would not test this — that passes even
// with no lock at all. So this holds an operation open, checks Close is blocked, then releases it.
func TestCloseWaitsForInFlightOperations(t *testing.T) {
	client := newTestClient(t)

	release, err := client.acquire()
	require.NoError(t, err)

	closed := make(chan error, 1)
	go func() { closed <- client.Close() }()

	select {
	case <-closed:
		t.Fatal("Close returned while an operation was still in flight")
	case <-time.After(50 * time.Millisecond):
	}

	release()

	select {
	case err := <-closed:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("Close did not return after the in-flight operation finished")
	}
}

func TestNewClientWithKeyRejectsInvalidOptions(t *testing.T) {
	credential, err := NewKeyCredential(testAccountKey)
	require.NoError(t, err)

	client, err := newTestClientWithKey(
		"https://myaccount.documents.azure.com",
		credential,
		&ClientOptions{Operation: OperationOptions{CustomHeaders: map[string]string{"x-test": "\x00"}}})

	require.ErrorContains(t, err, "custom header")
	require.Nil(t, client)
}

// TestClientOptionsValidateDiagnosticsVerbosity walks every defined DiagnosticsVerbosity constant
// plus an out-of-range value, verifying the enum's zero value and named members are all accepted
// and anything else is rejected before a client is ever constructed.
func TestClientOptionsValidateDiagnosticsVerbosity(t *testing.T) {
	for _, tt := range []struct {
		name      string
		verbosity DiagnosticsVerbosity
		wantErr   string
	}{
		{"zero value / Default", DiagnosticsVerbosityDefault, ""},
		{"Summary", DiagnosticsVerbositySummary, ""},
		{"Detailed", DiagnosticsVerbosityDetailed, ""},
		{"out of range", DiagnosticsVerbosity(99), "unrecognized DiagnosticsVerbosity(99)"},
		{"negative", DiagnosticsVerbosity(-1), "unrecognized DiagnosticsVerbosity(-1)"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := ClientOptions{DiagnosticsVerbosity: tt.verbosity}.validate()
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestNewClientWithKeyAcceptsOptions(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options ClientOptions
	}{
		{"defaults", ClientOptions{}},
		{"preferred regions", ClientOptions{Routing: PreferredRegions(RegionEastUS, RegionWestUS)}},
		{"proximity to known region", ClientOptions{Routing: ProximityTo(RegionEastUS)}},
		{"proximity to unknown region", ClientOptions{Routing: ProximityTo("not-a-real-region")}},
		{"diagnostics verbosity summary", ClientOptions{DiagnosticsVerbosity: DiagnosticsVerbositySummary}},
		{"diagnostics verbosity detailed", ClientOptions{DiagnosticsVerbosity: DiagnosticsVerbosityDetailed}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, &tt.options)
			require.NoError(t, err)
			require.NoError(t, client.Close())
		})
	}
}
