// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/stretchr/testify/require"
)

// doTestRequest sends a GET for a blob through the client's pipeline and returns the authorization
// scheme (the first token of the Authorization header) that reached the transport.
func doTestRequest(t *testing.T, azClient *azcore.Client, rawURL string) string {
	t.Helper()
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, rawURL)
	require.NoError(t, err)
	resp, err := azClient.Pipeline().Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Request, "the transport must report the request it sent")
	scheme, _, _ := strings.Cut(resp.Request.Header.Get("Authorization"), " ")
	return scheme
}

func TestGetAzClientSessionModeEnabledAccountNameError(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	opts := testClientOptions(srv)
	opts.Session = exported.SessionOptions{Mode: exported.SessionModeEnabled}

	// a host without a subdomain has no account name to derive
	_, err := GetAzClient("https://localhost/", fakeTokenCredential{}, nil, opts)
	require.Error(t, err)
	require.Contains(t, err.Error(), "account name could not be determined")
}

func TestGetAzClientUnsupportedSessionMode(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	opts := testClientOptions(srv)
	opts.Session = exported.SessionOptions{Mode: exported.SessionMode("bogus"), AccountName: "fakeaccount"}

	_, err := GetAzClient(fakeServiceURL, fakeTokenCredential{}, nil, opts)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported session mode")
}

// Session-based authentication signs with a session key obtained via a token credential, so
// enabling it without one is a configuration error rather than something to silently ignore.
func TestGetAzClientSessionRequiresTokenCredential(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	sharedKey, err := exported.NewSharedKeyCredential("fakeaccount", "dGVzdC1rZXk=")
	require.NoError(t, err)

	creds := []struct {
		name      string
		sharedKey *exported.SharedKeyCredential
	}{
		{"SharedKey", sharedKey},
		{"NoCredential", nil},
	}

	for _, tt := range creds {
		t.Run(tt.name, func(t *testing.T) {
			opts := testClientOptions(srv)
			opts.Session = exported.SessionOptions{
				Mode:        exported.SessionModeEnabled,
				AccountName: "fakeaccount",
			}

			_, err := GetAzClient(fakeServiceURL, nil, tt.sharedKey, opts)
			require.Error(t, err)
			require.Contains(t, err.Error(), "requires a TokenCredential")
		})
	}
}

func TestGetAzClientSharedKeyUnaffectedBySessionDefaults(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	sharedKey, err := exported.NewSharedKeyCredential("fakeaccount", "dGVzdC1rZXk=")
	require.NoError(t, err)

	azClient, err := GetAzClient(fakeServiceURL, nil, sharedKey, testClientOptions(srv))
	require.NoError(t, err)

	scheme := doTestRequest(t, azClient, fakeContainerURL+"/myblob")
	require.Equal(t, "SharedKey", scheme)
}

func TestGetAzClientNoCredentialHasNoAuthPolicy(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	azClient, err := GetAzClient(fakeServiceURL, nil, nil, testClientOptions(srv))
	require.NoError(t, err)

	scheme := doTestRequest(t, azClient, fakeContainerURL+"/myblob")
	require.Empty(t, scheme, "anonymous/SAS access must not be authenticated")
}
