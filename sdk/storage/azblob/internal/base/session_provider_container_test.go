// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"
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

func TestAcquireSession_FallbackToBearer(t *testing.T) {
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

func TestGetSessionConcurrentCallsCreateOneSession(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	// only one CreateSession response is queued; a second call would panic the mock server
	appendSessionResponse(srv, "key-one", "token-one", time.Now().Add(time.Hour))

	provider := newTestSessionProvider(t, srv)

	const goroutines = 8
	var wg sync.WaitGroup
	results := make([]exported.SessionCredential, goroutines)
	errs := make([]error, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i], errs[i] = provider.GetSession(newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob"))
		}(i)
	}
	wg.Wait()

	for i := 0; i < goroutines; i++ {
		require.NoError(t, errs[i])
		require.Equal(t, "token-one", results[i].Token())
	}
	require.Equal(t, 1, srv.Requests(), "concurrent callers must share a single CreateSession call")
}

func TestGetSessionReacquiresExpiredSession(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	// The first session is already expired, so the next call re-acquires unconditionally. Note
	// that a session merely *nearing* expiry is refreshed eagerly at most once every 30 seconds,
	// so that path can't be driven deterministically from here.
	appendSessionResponse(srv, "key-one", "token-one", time.Now().Add(-time.Minute))
	appendSessionResponse(srv, "key-two", "token-two", time.Now().Add(time.Hour))

	provider := newTestSessionProvider(t, srv)
	req := newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob")

	first, err := provider.GetSession(req)
	require.NoError(t, err)
	require.Equal(t, "token-one", first.Token())

	second, err := provider.GetSession(req)
	require.NoError(t, err)
	require.Equal(t, "token-two", second.Token(), "an expired session must be re-acquired")
	require.Equal(t, 2, srv.Requests())
}

func TestGetSessionErrorIsNotCached(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	// a hard failure is not a fallback decision, so it must not be cached
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusNotFound),
		mock.WithHeader("x-ms-error-code", "ContainerNotFound"),
		mock.WithBody(createErrorResponseXML("ContainerNotFound", "Container not found")),
	)
	appendSessionResponse(srv, "key-one", "token-one", time.Now().Add(time.Hour))

	provider := newTestSessionProvider(t, srv)
	req := newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob")

	_, err := provider.GetSession(req)
	require.Error(t, err)

	creds, err := provider.GetSession(req)
	require.NoError(t, err)
	require.Equal(t, "token-one", creds.Token())
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

// A request that reports a rejected session after it has already been replaced must not discard
// the replacement.
func TestInvalidateSessionIgnoresStaleCredential(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	// only one CreateSession response is queued; a second acquisition would panic the mock server
	appendSessionResponse(srv, "key-one", "token-one", time.Now().Add(time.Hour))

	provider := newTestSessionProvider(t, srv)

	req := newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob")
	current, err := provider.GetSession(req)
	require.NoError(t, err)
	require.Equal(t, "token-one", current.Token())
	require.Equal(t, 1, srv.Requests())

	// report a session that is no longer the cached one
	stale := exported.NewSessionCredential("stale-token", "key-zero", time.Now().Add(time.Hour))
	require.NoError(t, provider.InvalidateSession(req, stale))

	// the cached session survives, so no new one is acquired
	after, err := provider.GetSession(req)
	require.NoError(t, err)
	require.Equal(t, "token-one", after.Token())
	require.Equal(t, 1, srv.Requests())
}

// Many requests can fail with the same rejected session at once, but they must replace it once.
func TestInvalidateSessionConcurrentCallsReplaceSessionOnce(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	// only two sessions are queued; a second replacement would panic the mock server
	expiration := time.Now().Add(time.Hour)
	appendSessionResponse(srv, "key-one", "token-one", expiration)
	appendSessionResponse(srv, "key-two", "token-two", expiration)

	provider := newTestSessionProvider(t, srv)
	req := newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob")

	rejected, err := provider.GetSession(req)
	require.NoError(t, err)
	require.Equal(t, "token-one", rejected.Token())

	// every goroutine reports the same rejected session
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := provider.InvalidateSession(newTestRequest(t, http.MethodGet, fakeContainerURL+"/myblob"), rejected); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	replacement, err := provider.GetSession(req)
	require.NoError(t, err)
	require.Equal(t, "token-two", replacement.Token())
	require.Equal(t, 2, srv.Requests(), "concurrent invalidations must cause a single replacement")
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
	require.Equal(t, fakeServiceURL, provider.genClient.Endpoint())

	// the session scope is still resolved from the request URL, not from the URL the provider
	// was constructed with
	container, blob, err := shared.GetContainerAndBlobName(mustParseURL(t, fakeServiceURL+"othercontainer/myblob"))
	require.NoError(t, err)
	require.Equal(t, "othercontainer", container)
	require.Equal(t, "myblob", blob)
}
