// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/stretchr/testify/require"
)

const (
	fakeServiceURL   = "https://fakeaccount.blob.core.windows.net/"
	fakeContainerURL = fakeServiceURL + "testcontainer"
)

// fakeTokenCredential returns a static token so the bearer token policy can be exercised without
// contacting an identity provider.
type fakeTokenCredential struct{}

func (fakeTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fake-bearer-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

// testClientOptions returns options that route all traffic to the mock server and disable retries
// so each operation results in exactly one transport call.
func testClientOptions(srv *mock.Server) *ClientOptions {
	return &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: srv,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
	}
}

// newTestContainerClient creates a ContainerClient backed by a mock server for testing.
func newTestContainerClient(t *testing.T, srv *mock.Server) *generated.ContainerClient {
	t.Helper()
	azClient, err := azcore.NewClient("test", "v1.0.0", runtime.PipelineOptions{}, &policy.ClientOptions{
		Transport: srv,
		Retry:     policy.RetryOptions{MaxRetries: -1},
	})
	require.NoError(t, err)
	return generated.NewContainerClient(fakeContainerURL, azClient)
}

// newTestSessionProvider creates a containerSessionProvider backed by a mock server for testing.
func newTestSessionProvider(t *testing.T, srv *mock.Server) *containerSessionProvider {
	t.Helper()
	p, err := NewContainerSessionProvider(fakeTokenCredential{}, fakeServiceURL, testClientOptions(srv))
	require.NoError(t, err)
	provider, ok := p.(*containerSessionProvider)
	require.True(t, ok)
	return provider
}

// createSessionResponseXML creates a session response XML body for testing.
func createSessionResponseXML(sessionKey, sessionToken string, expiration time.Time) []byte {
	return []byte(`<?xml version="1.0" encoding="utf-8"?>
<CreateSessionResult>
	<AuthenticationType>HMAC</AuthenticationType>
	<Id>test-session-id</Id>
	<Credentials>
		<SessionKey>` + sessionKey + `</SessionKey>
		<SessionToken>` + sessionToken + `</SessionToken>
	</Credentials>
	<Expiration>` + expiration.Format(time.RFC1123) + `</Expiration>
</CreateSessionResult>`)
}

// createErrorResponseXML creates an error response XML body for testing.
func createErrorResponseXML(code, message string) []byte {
	return []byte(`<?xml version="1.0" encoding="utf-8"?>
<Error>
	<Code>` + code + `</Code>
	<Message>` + message + `</Message>
</Error>`)
}

// appendSessionResponse queues a successful CreateSession response on the mock server.
func appendSessionResponse(srv *mock.Server, sessionKey, sessionToken string, expiration time.Time) {
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusCreated),
		mock.WithBody(createSessionResponseXML(sessionKey, sessionToken, expiration)),
	)
}

func mustParseURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()
	u, err := url.Parse(rawURL)
	require.NoError(t, err)
	return u
}

func newTestRequest(t *testing.T, method, rawURL string) *http.Request {
	t.Helper()
	req, err := http.NewRequestWithContext(context.Background(), method, rawURL, nil)
	require.NoError(t, err)
	return req
}

func TestAcquireSession_Success(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	expiration := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	appendSessionResponse(srv, "test-key", "test-token", expiration)

	client := newTestContainerClient(t, srv)

	creds, exp, err := acquireSession(client)(context.Background())
	require.NoError(t, err)
	require.False(t, creds.Fallback())
	require.Equal(t, "test-key", creds.Key())
	require.Equal(t, "test-token", creds.Token())
	require.Equal(t, expiration.Format(time.RFC1123), exp.Format(time.RFC1123))
	require.Equal(t, expiration.Format(time.RFC1123), creds.Expiry().Format(time.RFC1123))
}

func TestAcquireSession_FallbackToBearer_TransientFailure(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		errorCode  string
	}{
		{
			name:       "ServerError500",
			statusCode: http.StatusInternalServerError,
			errorCode:  "InternalError",
		},
		{
			name:       "ServiceUnavailable503",
			statusCode: http.StatusServiceUnavailable,
			errorCode:  "ServiceUnavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeFn()

			srv.SetResponse(
				mock.WithStatusCode(tt.statusCode),
				mock.WithHeader("x-ms-error-code", tt.errorCode),
				mock.WithBody(createErrorResponseXML(tt.errorCode, "error message")),
			)

			client := newTestContainerClient(t, srv)

			before := time.Now()
			creds, exp, err := acquireSession(client)(context.Background())
			require.NoError(t, err)
			require.True(t, creds.Fallback())
			require.Empty(t, creds.Token())
			require.Empty(t, creds.Key())
			// the fallback decision is cached for a short cooldown so sessions are retried soon
			require.WithinDuration(t, before.Add(transientFailureCooldown), exp, time.Minute)
		})
	}
}

func TestAcquireSession_FallbackToBearer_FeatureUnavailable(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		errorCode  string
	}{
		{
			name:       "FeatureNotEnabled",
			statusCode: http.StatusBadRequest,
			errorCode:  featureNotEnabled,
		},
		{
			name:       "Forbidden",
			statusCode: http.StatusForbidden,
			errorCode:  "AuthorizationFailure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeFn()

			srv.AppendResponse(
				mock.WithStatusCode(tt.statusCode),
				mock.WithHeader("x-ms-error-code", tt.errorCode),
				mock.WithBody(createErrorResponseXML(tt.errorCode, "error message")),
			)

			client := newTestContainerClient(t, srv)

			before := time.Now()
			creds, exp, err := acquireSession(client)(context.Background())
			require.NoError(t, err)
			require.True(t, creds.Fallback())
			// the feature is unavailable, so the decision is cached for much longer
			require.WithinDuration(t, before.Add(featureUnavailableCooldown), exp, time.Minute)
		})
	}
}

func TestAcquireSession_BadRequestOtherErrorCodeReturnsError(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	srv.AppendResponse(
		mock.WithStatusCode(http.StatusBadRequest),
		mock.WithHeader("x-ms-error-code", "InvalidQueryParameterValue"),
		mock.WithBody(createErrorResponseXML("InvalidQueryParameterValue", "bad parameter")),
	)

	client := newTestContainerClient(t, srv)

	creds, _, err := acquireSession(client)(context.Background())
	require.Error(t, err)
	require.False(t, creds.Fallback())
}

func TestAcquireSession_Error(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	srv.AppendResponse(
		mock.WithStatusCode(http.StatusNotFound),
		mock.WithHeader("x-ms-error-code", "ContainerNotFound"),
		mock.WithBody(createErrorResponseXML("ContainerNotFound", "Container not found")),
	)

	client := newTestContainerClient(t, srv)

	creds, _, err := acquireSession(client)(context.Background())
	require.Error(t, err)
	// Should NOT be a fallback - this is a real error that should propagate
	require.False(t, creds.Fallback())

	var respErr *azcore.ResponseError
	require.True(t, errors.As(err, &respErr))
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
}

func TestShouldRefreshSession(t *testing.T) {
	tests := []struct {
		name     string
		creds    exported.SessionCredential
		expected bool
	}{
		{
			name:     "NotExpiringSoon",
			creds:    exported.NewSessionCredential("token", "key", time.Now().Add(5*time.Minute)),
			expected: false,
		},
		{
			name:     "ExpiringSoon",
			creds:    exported.NewSessionCredential("token", "key", time.Now().Add(10*time.Second)),
			expected: true,
		},
		{
			name:     "AlreadyExpired",
			creds:    exported.NewSessionCredential("token", "key", time.Now().Add(-1*time.Minute)),
			expected: true,
		},
		{
			name:     "FallbackIsNeverRefreshedEarly",
			creds:    exported.NewSessionCredentialFallback(time.Now().Add(10 * time.Second)),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, shouldRefreshSession(tt.creds, context.Background()))
		})
	}
}

func TestGetContainerClient(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	provider := newTestSessionProvider(t, srv)

	containerClient := provider.getContainerClient("mycontainer")
	require.NotNil(t, containerClient)
	require.Equal(t, fakeServiceURL+"mycontainer", containerClient.Endpoint())
}

func TestIsRequestEligible(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	provider := newTestSessionProvider(t, srv)

	tests := []struct {
		name     string
		req      *http.Request
		expected bool
	}{
		{
			name:     "NilRequest",
			req:      nil,
			expected: false,
		},
		{
			name:     "GetBlob",
			req:      newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob"),
			expected: true,
		},
		{
			name:     "GetBlobInVirtualDirectory",
			req:      newTestRequest(t, http.MethodGet, fakeContainerURL+"/dir/sub/myblob"),
			expected: true,
		},
		{
			name:     "PutBlob",
			req:      newTestRequest(t, http.MethodPut, fakeContainerURL+"/myblob"),
			expected: false,
		},
		{
			name:     "HeadBlob",
			req:      newTestRequest(t, http.MethodHead, fakeContainerURL+"/myblob"),
			expected: false,
		},
		{
			name:     "GetWithCompQueryParameter",
			req:      newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob?comp=tags"),
			expected: false,
		},
		{
			name:     "ContainerLevelGet",
			req:      newTestRequest(t, http.MethodGet, fakeContainerURL),
			expected: false,
		},
		{
			name:     "ServiceLevelGet",
			req:      newTestRequest(t, http.MethodGet, fakeServiceURL),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, provider.IsRequestEligible(tt.req))
		})
	}
}

func TestIsRequestEligibleNilURL(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	provider := newTestSessionProvider(t, srv)

	req := newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob")
	req.URL = nil
	require.False(t, provider.IsRequestEligible(req))
}

func TestGetSessionCachesPerContainer(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	expiration := time.Now().Add(time.Hour)
	appendSessionResponse(srv, "key-one", "token-one", expiration)
	appendSessionResponse(srv, "key-two", "token-two", expiration)

	provider := newTestSessionProvider(t, srv)

	req := newTestRequest(t, http.MethodGet, fakeServiceURL+"container-one/myblob")
	first, err := provider.GetSession(req)
	require.NoError(t, err)
	require.Equal(t, "token-one", first.Token())
	require.Equal(t, 1, srv.Requests())

	// a second request for the same container reuses the cached session
	second, err := provider.GetSession(newTestRequest(t, http.MethodGet, fakeServiceURL+"container-one/otherblob"))
	require.NoError(t, err)
	require.Equal(t, "token-one", second.Token())
	require.Equal(t, 1, srv.Requests())

	// a different container gets its own session
	other, err := provider.GetSession(newTestRequest(t, http.MethodGet, fakeServiceURL+"container-two/myblob"))
	require.NoError(t, err)
	require.Equal(t, "token-two", other.Token())
	require.Equal(t, 2, srv.Requests())
}

func TestInvalidateSession(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	expiration := time.Now().Add(time.Hour)
	appendSessionResponse(srv, "key-one", "token-one", expiration)
	appendSessionResponse(srv, "key-two", "token-two", expiration)

	provider := newTestSessionProvider(t, srv)

	req := newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob")
	first, err := provider.GetSession(req)
	require.NoError(t, err)
	require.Equal(t, "token-one", first.Token())
	require.Equal(t, 1, srv.Requests())

	require.NoError(t, provider.InvalidateSession(req, first))

	// the cached session was discarded, so a new one is acquired
	second, err := provider.GetSession(req)
	require.NoError(t, err)
	require.Equal(t, "token-two", second.Token())
	require.Equal(t, 2, srv.Requests())
}

func TestGetSessionFallbackIsCached(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	srv.SetResponse(
		mock.WithStatusCode(http.StatusBadRequest),
		mock.WithHeader("x-ms-error-code", featureNotEnabled),
		mock.WithBody(createErrorResponseXML(featureNotEnabled, "feature not enabled")),
	)

	provider := newTestSessionProvider(t, srv)

	req := newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob")
	creds, err := provider.GetSession(req)
	require.NoError(t, err)
	require.True(t, creds.Fallback())
	require.Equal(t, 1, srv.Requests())

	// the fallback decision is cached, so no further CreateSession calls are made
	creds, err = provider.GetSession(req)
	require.NoError(t, err)
	require.True(t, creds.Fallback())
	require.Equal(t, 1, srv.Requests())
}

func TestGetSessionRequiresContainer(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	provider := newTestSessionProvider(t, srv)

	_, err := provider.GetSession(newTestRequest(t, http.MethodGet, fakeServiceURL))
	require.Error(t, err)

	err = provider.InvalidateSession(newTestRequest(t, http.MethodGet, fakeServiceURL), exported.SessionCredential{})
	require.Error(t, err)
}

func TestResourceForRequestNilRequest(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	provider := newTestSessionProvider(t, srv)

	_, err := provider.resourceForRequest(nil)
	require.Error(t, err)
}

func TestNewContainerSessionProviderInvalidURL(t *testing.T) {
	_, err := NewContainerSessionProvider(fakeTokenCredential{}, "://not a url", nil)
	require.Error(t, err)
}

func TestNewContainerSessionProviderTrimsPath(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	p, err := NewContainerSessionProvider(fakeTokenCredential{}, fakeContainerURL+"/myblob", testClientOptions(srv))
	require.NoError(t, err)
	provider, ok := p.(*containerSessionProvider)
	require.True(t, ok)
	require.Equal(t, fakeServiceURL, provider.svcURL)

	// the session scope is still resolved from the request URL, not from the URL the provider
	// was constructed with
	container, blob, err := shared.GetContainerAndBlobName(mustParseURL(t, fakeServiceURL+"othercontainer/myblob"))
	require.NoError(t, err)
	require.Equal(t, "othercontainer", container)
	require.Equal(t, "myblob", blob)
}
