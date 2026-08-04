// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T) *Client {
	t.Helper()

	cred, err := NewKeyCredential("key")
	require.NoError(t, err)
	client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, nil)
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

func TestNewClientWithKeyRejectsZeroCredential(t *testing.T) {
	// The zero value is what a caller gets if they ignore the error from NewKeyCredential.
	_, err := NewClientWithKey("https://myaccount.documents.azure.com", KeyCredential{}, nil)
	require.Error(t, err)
}

func TestNewKeyCredentialRejectsEmptyKey(t *testing.T) {
	_, err := NewKeyCredential("")
	require.Error(t, err)
}

func TestNewClientRejectsNonAbsoluteEndpoint(t *testing.T) {
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)

	for _, endpoint := range []string{"", "myaccount.documents.azure.com", "/relative/path", "https://"} {
		t.Run(endpoint, func(t *testing.T) {
			_, err := NewClientWithKey(endpoint, cred, nil)
			require.Error(t, err, "endpoint %q should be rejected", endpoint)
		})
	}
}

func TestNewClientAcceptsAbsoluteEndpoint(t *testing.T) {
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)

	client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, nil)
	require.NoError(t, err)
	require.Equal(t, "https://myaccount.documents.azure.com", client.Endpoint())
}

func TestNewClientAppliesCloudDefault(t *testing.T) {
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)

	// The default has to be applied rather than left as the zero Configuration, or the driver
	// receives an empty authority host and authenticates against the wrong audience.
	for _, tt := range []struct {
		name    string
		options *ClientOptions
	}{
		{"nil options", nil},
		{"options without a cloud", &ClientOptions{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, tt.options)
			require.NoError(t, err)
			require.Equal(t, cloud.AzurePublic, client.options.Cloud)
		})
	}
}

func TestNewClientKeepsExplicitCloud(t *testing.T) {
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)

	client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, &ClientOptions{Cloud: cloud.AzureGovernment})
	require.NoError(t, err)
	require.Equal(t, cloud.AzureGovernment, client.options.Cloud)
}

// The client must not alias the caller's slice, or a later append by the caller silently changes
// the routing order of a client that has already been created.
func TestNewClientCopiesPreferredRegions(t *testing.T) {
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)

	regions := []string{"West US", "East US"}
	client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, &ClientOptions{PreferredRegions: regions})
	require.NoError(t, err)

	regions[0] = "North Europe"
	require.Equal(t, []string{"West US", "East US"}, client.options.PreferredRegions)
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

func TestParseConnectionString(t *testing.T) {
	for _, tt := range []struct {
		name             string
		connectionString string
		wantEndpoint     string
		wantKey          string
	}{
		{
			name:             "endpoint and key",
			connectionString: "AccountEndpoint=https://myaccount.documents.azure.com;AccountKey=somekey;",
			wantEndpoint:     "https://myaccount.documents.azure.com",
			wantKey:          "somekey",
		},
		{
			name:             "reversed order",
			connectionString: "AccountKey=somekey;AccountEndpoint=https://myaccount.documents.azure.com",
			wantEndpoint:     "https://myaccount.documents.azure.com",
			wantKey:          "somekey",
		},
		{
			name:             "keys are matched case-insensitively",
			connectionString: "accountendpoint=https://myaccount.documents.azure.com;ACCOUNTKEY=somekey",
			wantEndpoint:     "https://myaccount.documents.azure.com",
			wantKey:          "somekey",
		},
		{
			name:             "unrecognized pairs are ignored so a string copied from the portal works",
			connectionString: "AccountEndpoint=https://myaccount.documents.azure.com;AccountKey=somekey;ApiKind=SQL",
			wantEndpoint:     "https://myaccount.documents.azure.com",
			wantKey:          "somekey",
		},
		{
			name:             "base64 keys containing '=' are not truncated",
			connectionString: "AccountEndpoint=https://myaccount.documents.azure.com;AccountKey=C2y6yDjf5R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==",
			wantEndpoint:     "https://myaccount.documents.azure.com",
			wantKey:          "C2y6yDjf5R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			endpoint, cred, err := parseConnectionString(tt.connectionString)
			require.NoError(t, err)
			require.Equal(t, tt.wantEndpoint, endpoint)
			require.Equal(t, tt.wantKey, cred.accountKey)
		})
	}
}

func TestParseConnectionStringErrors(t *testing.T) {
	for _, tt := range []struct {
		name             string
		connectionString string
	}{
		{"empty", ""},
		{"missing key", "AccountEndpoint=https://myaccount.documents.azure.com"},
		{"missing endpoint", "AccountKey=somekey"},
		{"empty key value", "AccountEndpoint=https://myaccount.documents.azure.com;AccountKey="},
		{"not key=value pairs", "AccountEndpoint=https://myaccount.documents.azure.com;garbage"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseConnectionString(tt.connectionString)
			require.Error(t, err)
		})
	}
}

func TestNewClientFromConnectionString(t *testing.T) {
	client, err := NewClientFromConnectionString("AccountEndpoint=https://myaccount.documents.azure.com;AccountKey=somekey;", nil)
	require.NoError(t, err)
	require.Equal(t, "https://myaccount.documents.azure.com", client.Endpoint())
}

func TestNewClientFromConnectionStringPropagatesParseErrors(t *testing.T) {
	_, err := NewClientFromConnectionString("AccountEndpoint=https://myaccount.documents.azure.com", nil)
	require.Error(t, err)
}

// The endpoint embedded in a connection string is validated the same way as one passed directly.
func TestNewClientFromConnectionStringValidatesEndpoint(t *testing.T) {
	_, err := NewClientFromConnectionString("AccountEndpoint=myaccount.documents.azure.com;AccountKey=somekey", nil)
	require.Error(t, err)
}

func TestNewDatabaseAndNewContainer(t *testing.T) {
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)
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
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)
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
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)

	for _, endpoint := range []string{"ftp://myaccount", "file:///tmp/x", "wss://myaccount"} {
		t.Run(endpoint, func(t *testing.T) {
			_, err := NewClientWithKey(endpoint, cred, nil)
			require.Error(t, err)
		})
	}
}

// The emulator's connection string is the one every Cosmos developer pastes, so it is the highest
// value single case for the parser.
func TestParseConnectionStringAcceptsTheEmulatorString(t *testing.T) {
	const emulator = "AccountEndpoint=https://localhost:8081/;AccountKey=C2y6yDjf5R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

	endpoint, cred, err := parseConnectionString(emulator)
	require.NoError(t, err)
	require.Equal(t, "https://localhost:8081/", endpoint)
	require.Equal(t, "C2y6yDjf5R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==", cred.accountKey)

	client, err := NewClientFromConnectionString(emulator, nil)
	require.NoError(t, err)
	require.NoError(t, client.Close())
}

// A connection string carries a credential, so nothing derived from it may reach an error message
// or a log.
func TestParseConnectionStringNeverLeaksTheAccountKey(t *testing.T) {
	const key = "supersecretaccountkey"

	for _, connectionString := range []string{
		"AccountKey=" + key,
		"AccountKey=" + key + ";garbage",
		"AccountEndpoint=;AccountKey=" + key,
	} {
		_, _, err := parseConnectionString(connectionString)
		require.Error(t, err)
		require.NotContains(t, err.Error(), key)
	}

	_, err := NewClientFromConnectionString("AccountEndpoint=notaurl;AccountKey="+key, nil)
	require.Error(t, err)
	require.NotContains(t, err.Error(), key)
}

// The endpoint is named in its error so the message is actionable; this also pins the shape the
// key-leak test above relies on.
func TestEndpointErrorNamesTheEndpoint(t *testing.T) {
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)

	_, err = NewClientWithKey("notaurl", cred, nil)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "notaurl"))
}
