// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

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

func TestNewClientWithNilOptionsUsesDefaults(t *testing.T) {
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)

	client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, nil)
	require.NoError(t, err)
	require.Empty(t, client.options.PreferredRegions)
	require.False(t, client.options.EnableContentResponseOnWrite)
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

func TestCloseIsIdempotent(t *testing.T) {
	cred, err := NewKeyCredential("key")
	require.NoError(t, err)

	client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, nil)
	require.NoError(t, err)

	require.NoError(t, client.Close())
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
