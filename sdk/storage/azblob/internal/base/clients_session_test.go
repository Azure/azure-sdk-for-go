// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/stretchr/testify/require"
)

// stubSessionProvider is a SessionProvider that records calls, letting tests assert that a
// caller-supplied provider is wired into the pipeline instead of the default one.
type stubSessionProvider struct {
	eligible bool
	cred     exported.SessionCredential
	getCalls int
}

func (s *stubSessionProvider) GetSession(*http.Request) (exported.SessionCredential, error) {
	s.getCalls++
	return s.cred, nil
}

func (s *stubSessionProvider) InvalidateSession(*http.Request, exported.SessionCredential) error {
	return nil
}

func (s *stubSessionProvider) IsRequestEligible(*http.Request) bool {
	return s.eligible
}

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

func TestGetAzClientSessionModeBearerOnly(t *testing.T) {
	modes := []struct {
		name string
		mode exported.SessionMode
	}{
		{"Default", exported.SessionModeDefault},
		{"Disabled", exported.SessionModeDisabled},
	}

	for _, tt := range modes {
		t.Run(tt.name, func(t *testing.T) {
			srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeFn()
			srv.SetResponse(mock.WithStatusCode(http.StatusOK))

			opts := testClientOptions(srv)
			opts.Session = exported.SessionOptions{Mode: tt.mode}

			azClient, err := GetAzClient(fakeServiceURL, fakeTokenCredential{}, nil, opts)
			require.NoError(t, err)

			scheme := doTestRequest(t, azClient, fakeContainerURL+"/myblob")
			require.Equal(t, "Bearer", scheme)
			require.Equal(t, 1, srv.Requests(), "no session should be created")
		})
	}
}

func TestGetAzClientSessionModeEnabledUsesDefaultProvider(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	// first the CreateSession call, then the blob GET
	appendSessionResponse(srv, "dGVzdC1rZXk=", "test-token", time.Now().Add(time.Hour))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

	opts := testClientOptions(srv)
	opts.Session = exported.SessionOptions{Mode: exported.SessionModeEnabled}

	azClient, err := GetAzClient(fakeServiceURL, fakeTokenCredential{}, nil, opts)
	require.NoError(t, err)

	scheme := doTestRequest(t, azClient, fakeContainerURL+"/myblob")
	require.Equal(t, "Session", scheme)
	require.Equal(t, 2, srv.Requests(), "one CreateSession plus the blob GET")
}

func TestGetAzClientSessionModeEnabledUsesSuppliedProvider(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	provider := &stubSessionProvider{
		eligible: true,
		cred:     exported.NewSessionCredential("supplied-token", "dGVzdC1rZXk=", time.Now().Add(time.Hour)),
	}

	opts := testClientOptions(srv)
	opts.Session = exported.SessionOptions{
		Mode:        exported.SessionModeEnabled,
		AccountName: "fakeaccount",
		Provider:    provider,
	}

	azClient, err := GetAzClient(fakeServiceURL, fakeTokenCredential{}, nil, opts)
	require.NoError(t, err)

	scheme := doTestRequest(t, azClient, fakeContainerURL+"/myblob")
	require.Equal(t, "Session", scheme)
	require.Equal(t, 1, provider.getCalls, "the supplied provider must be used")
	require.Equal(t, 1, srv.Requests(), "the supplied provider does not call CreateSession")
}

func TestGetAzClientSessionModeEnabledFallsBackWhenProviderIneligible(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	provider := &stubSessionProvider{eligible: false}

	opts := testClientOptions(srv)
	opts.Session = exported.SessionOptions{
		Mode:        exported.SessionModeEnabled,
		AccountName: "fakeaccount",
		Provider:    provider,
	}

	azClient, err := GetAzClient(fakeServiceURL, fakeTokenCredential{}, nil, opts)
	require.NoError(t, err)

	scheme := doTestRequest(t, azClient, fakeContainerURL+"/myblob")
	require.Equal(t, "Bearer", scheme)
	require.Equal(t, 0, provider.getCalls)
}

func TestGetAzClientSessionModeEnabledDerivesAccountName(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	provider := &stubSessionProvider{
		eligible: true,
		cred:     exported.NewSessionCredential("token", "dGVzdC1rZXk=", time.Now().Add(time.Hour)),
	}

	// AccountName is omitted, so it must be derived from the service URL
	opts := testClientOptions(srv)
	opts.Session = exported.SessionOptions{
		Mode:     exported.SessionModeEnabled,
		Provider: provider,
	}

	azClient, err := GetAzClient(fakeServiceURL, fakeTokenCredential{}, nil, opts)
	require.NoError(t, err)

	scheme := doTestRequest(t, azClient, fakeContainerURL+"/myblob")
	require.Equal(t, "Session", scheme)
}

func TestGetAzClientSessionModeEnabledAccountNameError(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	opts := testClientOptions(srv)
	opts.Session = exported.SessionOptions{
		Mode:     exported.SessionModeEnabled,
		Provider: &stubSessionProvider{},
	}

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

func TestGetAzClientSessionIgnoredWithoutTokenCredential(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	sharedKey, err := exported.NewSharedKeyCredential("fakeaccount", "dGVzdC1rZXk=")
	require.NoError(t, err)

	provider := &stubSessionProvider{eligible: true}
	opts := testClientOptions(srv)
	opts.Session = exported.SessionOptions{
		Mode:        exported.SessionModeEnabled,
		AccountName: "fakeaccount",
		Provider:    provider,
	}

	// session auth requires a token credential; shared key auth is unaffected by session options
	azClient, err := GetAzClient(fakeServiceURL, nil, sharedKey, opts)
	require.NoError(t, err)

	scheme := doTestRequest(t, azClient, fakeContainerURL+"/myblob")
	require.Equal(t, "SharedKey", scheme)
	require.Equal(t, 0, provider.getCalls)
}

func TestGetAzClientNoCredentialHasNoAuthPolicy(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	opts := testClientOptions(srv)
	opts.Session = exported.SessionOptions{Mode: exported.SessionModeEnabled, AccountName: "fakeaccount"}

	azClient, err := GetAzClient(fakeServiceURL, nil, nil, opts)
	require.NoError(t, err)

	scheme := doTestRequest(t, azClient, fakeContainerURL+"/myblob")
	require.Empty(t, scheme, "anonymous/SAS access must not be authenticated")
}

// errTokenCredential always fails to produce a token.
type errTokenCredential struct{ err error }

func (e errTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{}, e.err
}

func TestGetAzClientSessionFallsBackToBearerAndSurfacesTokenError(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	tokenErr := errors.New("no token for you")
	provider := &stubSessionProvider{
		eligible: true,
		cred:     exported.NewSessionCredentialFallback(time.Now().Add(time.Hour)),
	}

	opts := testClientOptions(srv)
	opts.Session = exported.SessionOptions{
		Mode:        exported.SessionModeEnabled,
		AccountName: "fakeaccount",
		Provider:    provider,
	}

	azClient, err := GetAzClient(fakeServiceURL, errTokenCredential{err: tokenErr}, nil, opts)
	require.NoError(t, err)

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, fakeContainerURL+"/myblob")
	require.NoError(t, err)

	// a fallback credential routes to the bearer token policy, whose failure must surface
	_, err = azClient.Pipeline().Do(req)
	require.ErrorIs(t, err, tokenErr)
	require.Equal(t, 1, provider.getCalls)
}



