// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
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

// Token credentials are accepted without network I/O; Initialize or an operation acquires a token.
func TestNewClientAcceptsTokenCredential(t *testing.T) {
	client, err := NewClient("https://myaccount.documents.azure.com", fakeTokenCredential{}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	require.Equal(t, "https://myaccount.documents.azure.com", client.Endpoint())
}

func TestNewClientRejectsNilCredential(t *testing.T) {
	_, err := NewClient("https://myaccount.documents.azure.com", nil, nil)
	require.Error(t, err)
}

func TestInitializeRejectsNilContext(t *testing.T) {
	client := newTestClient(t)
	err := client.Initialize(nil) //nolint:staticcheck // verifies the guard
	require.Error(t, err)
}

func TestInitializeReportsUnavailableDriver(t *testing.T) {
	if driverAvailable {
		t.Skip("the driver is available in this build")
	}
	client := newTestClient(t)
	err := client.Initialize(t.Context())
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Contains(t, cosmosErr.Message, "cannot reach the Cosmos driver")
}

// The zero value is what a caller gets if they ignore the error from NewKeyCredential.
func TestNewClientWithKeyRejectsZeroCredential(t *testing.T) {
	_, err := NewClientWithKey("https://myaccount.documents.azure.com", KeyCredential{}, nil)
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
			_, err := NewClientWithKey(endpoint, cred, nil)
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

func TestRoutingStrategy(t *testing.T) {
	require.Equal(t, RegionEastUS, ProximityTo(RegionEastUS).proximityTo)
	require.Equal(t, []Region{RegionWestUS, RegionEastUS}, PreferredRegions(RegionWestUS, RegionEastUS).preferredRegions)

	// The zero value carries neither, which is what "leave the order to the account" means.
	var zero RoutingStrategy
	require.Empty(t, zero.proximityTo)
	require.Empty(t, zero.preferredRegions)
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

func TestNewDatabaseAndNewContainer(t *testing.T) {
	client := newTestClient(t)

	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	require.Equal(t, "db", database.ID())

	container, err := database.NewContainer("items")
	require.NoError(t, err)
	require.Equal(t, "items", container.ID())

	direct, err := client.NewContainer("db", "items")
	require.NoError(t, err)
	require.Equal(t, "items", direct.ID())
	require.Equal(t, "db", direct.database.ID())
}

func TestNewDatabaseAndNewContainerRejectEmptyIDs(t *testing.T) {
	client := newTestClient(t)

	_, err := client.NewDatabase("")
	require.Error(t, err)

	_, err = client.NewContainer("db", "")
	require.Error(t, err)

	_, err = client.NewContainer("", "items")
	require.Error(t, err)
}

// The scheme is validated, because a Cosmos endpoint is always reached over HTTP(S) and anything
// else is a copy-paste mistake that would otherwise surface as an opaque driver failure.
func TestNewClientRejectsNonHTTPSchemes(t *testing.T) {
	cred := mustKeyCredential(t)

	for _, endpoint := range []string{"ftp://myaccount", "file:///tmp/x", "wss://myaccount"} {
		t.Run(endpoint, func(t *testing.T) {
			_, err := NewClientWithKey(endpoint, cred, nil)
			require.Error(t, err)
		})
	}
}

// The endpoint is named in its error so the message is actionable; this also pins the shape the
// key-leak test above relies on.
func TestEndpointErrorNamesTheEndpoint(t *testing.T) {
	cred := mustKeyCredential(t)

	_, err := NewClientWithKey("notaurl", cred, nil)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "notaurl"))
}

// After Close the client must refuse work rather than hand a released driver handle to an
// operation.
func TestOperationsAfterCloseAreRejected(t *testing.T) {
	client := newTestClient(t)
	container, err := client.NewContainer("db", "items")
	require.NoError(t, err)
	require.NoError(t, client.Close())

	_, err = container.ReadItem(context.Background(), NewPartitionKeyString("pk"), "item-1", nil)
	var cosmosErr *Error
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)

	_, err = container.CreateItem(context.Background(), NewPartitionKeyString("pk"), "x", []byte(`{"id":"x"}`), nil)
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

// Options the binding cannot implement are rejected when the client is constructed.
func TestNewClientWithKeyRejectsUnusableOptions(t *testing.T) {
	credential, err := NewKeyCredential(testAccountKey)
	require.NoError(t, err)

	for _, tt := range []struct {
		name    string
		options ClientOptions
		wantIn  string
	}{
		{
			name:    "proximity routing",
			options: ClientOptions{Routing: ProximityTo(RegionEastUS)},
			wantIn:  "ProximityTo is not supported",
		},
		{
			name:    "application id with NUL",
			options: ClientOptions{ApplicationID: "order\x00service"},
			wantIn:  "must not contain a NUL byte",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClientWithKey(
				"https://myaccount.documents.azure.com",
				credential,
				&tt.options)

			require.ErrorContains(t, err, tt.wantIn)
			require.Nil(t, client)
		})
	}
}

// The values the binding can pass through have to survive local construction.
func TestNewClientWithKeyAcceptsUsableOptions(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options ClientOptions
	}{
		{"defaults", ClientOptions{}},
		{"preferred regions", ClientOptions{Routing: PreferredRegions(RegionEastUS, RegionWestUS)}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, &tt.options)
			require.NoError(t, err)
			require.NoError(t, client.Close())
		})
	}
}
