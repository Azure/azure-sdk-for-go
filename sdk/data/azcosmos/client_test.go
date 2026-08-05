// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T) *Client {
	t.Helper()

	client, err := NewClientWithKey("https://myaccount.documents.azure.com", azcore.NewKeyCredential("key"), nil)
	require.NoError(t, err)
	return client
}

type fakeTokenCredential struct{}

func (fakeTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{}, nil
}

func TestNewClientAcceptsTokenCredential(t *testing.T) {
	client, err := NewClient("https://myaccount.documents.azure.com", fakeTokenCredential{}, nil)
	require.NoError(t, err)
	require.Equal(t, "https://myaccount.documents.azure.com", client.Endpoint())
}

func TestNewClientRejectsNilCredential(t *testing.T) {
	_, err := NewClient("https://myaccount.documents.azure.com", nil, nil)
	require.Error(t, err)
}
func TestNewClientWithKeyRejectsNilCredential(t *testing.T) {
	_, err := NewClientWithKey("https://myaccount.documents.azure.com", nil, nil)
	require.Error(t, err)
}

func TestNewClientRejectsNonAbsoluteEndpoint(t *testing.T) {
	cred := azcore.NewKeyCredential("key")

	for _, endpoint := range []string{"", "myaccount.documents.azure.com", "/relative/path", "https://"} {
		t.Run(endpoint, func(t *testing.T) {
			_, err := NewClientWithKey(endpoint, cred, nil)
			require.Error(t, err, "endpoint %q should be rejected", endpoint)
		})
	}
}

func TestNewClientAcceptsAbsoluteEndpoint(t *testing.T) {
	cred := azcore.NewKeyCredential("key")

	client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, nil)
	require.NoError(t, err)
	require.Equal(t, "https://myaccount.documents.azure.com", client.Endpoint())
}

// The client must not alias the caller's slice, or a later append by the caller silently changes
// the routing order of a client that has already been created.
func TestNewClientCopiesRoutingRegions(t *testing.T) {
	regions := []Region{RegionWestUS, RegionEastUS}
	client, err := NewClientWithKey("https://myaccount.documents.azure.com", azcore.NewKeyCredential("key"), &ClientOptions{
		Routing: PreferredRegions(regions...),
	})
	require.NoError(t, err)

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
	cred := azcore.NewKeyCredential("key")
	client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, nil)
	require.NoError(t, err)

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
	cred := azcore.NewKeyCredential("key")
	client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, nil)
	require.NoError(t, err)

	_, err = client.NewDatabase("")
	require.Error(t, err)

	_, err = client.NewContainer("db", "")
	require.Error(t, err)

	_, err = client.NewContainer("", "items")
	require.Error(t, err)
}

// The scheme is validated, because a Cosmos endpoint is always reached over HTTP(S) and anything
// else is a copy-paste mistake that would otherwise surface as an opaque driver failure.
func TestNewClientRejectsNonHTTPSchemes(t *testing.T) {
	cred := azcore.NewKeyCredential("key")

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
	cred := azcore.NewKeyCredential("key")

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

	_, err = container.CreateItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"x"}`), nil)
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
}
