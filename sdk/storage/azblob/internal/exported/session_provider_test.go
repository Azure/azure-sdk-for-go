// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/internal/locality"
	"github.com/stretchr/testify/require"
)

const providerServiceURL = "https://testaccount.blob.core.windows.net/"

// providerClock is a manually advanced clock for the provider and its caches.
type providerClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *providerClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *providerClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// createSessionCall records one Create Session request the service received.
type createSessionCall struct {
	container string
	host      string
	query     string
	body      string
}

// sessionService is a scripted Create Session endpoint. By default it issues session
// "<container>-<n>" valid for an hour; respond overrides that per call.
type sessionService struct {
	clock *providerClock

	mu    sync.Mutex
	calls []createSessionCall
	// respond, when set, produces the response for the nth (1-based) call to a container.
	respond func(container string, n int) *http.Response
	// gate, when set, blocks every Create Session until it is closed.
	gate chan struct{}
}

func (s *sessionService) Do(req *http.Request) (*http.Response, error) {
	if s.gate != nil {
		<-s.gate
	}
	body := ""
	if req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		body = string(b)
	}
	container := strings.Trim(req.URL.Path, "/")
	s.mu.Lock()
	s.calls = append(s.calls, createSessionCall{container: container, host: req.URL.Host, query: req.URL.RawQuery, body: body})
	n := 0
	for _, c := range s.calls {
		if c.container == container {
			n++
		}
	}
	respond := s.respond
	s.mu.Unlock()

	if respond != nil {
		if resp := respond(container, n); resp != nil {
			resp.Request = req
			return resp, nil
		}
	}
	return sessionResponse(req, fmt.Sprintf("%s-%d", container, n), s.clock.Now().Add(time.Hour)), nil
}

func (s *sessionService) callsTo(container string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, c := range s.calls {
		if c.container == container {
			n++
		}
	}
	return n
}

func (s *sessionService) allCalls() []createSessionCall {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]createSessionCall(nil), s.calls...)
}

func sessionResponse(req *http.Request, token string, expires time.Time) *http.Response {
	body := `<?xml version="1.0" encoding="utf-8"?><CreateSessionResult><AuthenticationType>HMAC</AuthenticationType>` +
		`<Id>id</Id><Credentials><SessionKey>` + testSessionKey + `</SessionKey><SessionToken>` + token + `</SessionToken></Credentials>` +
		`<Expiration>` + expires.UTC().Format(time.RFC1123) + `</Expiration></CreateSessionResult>`
	return &http.Response{StatusCode: http.StatusCreated, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body)), Request: req}
}

func errorResponse(status int, code string) *http.Response {
	body := `<?xml version="1.0" encoding="utf-8"?><Error><Code>` + code + `</Code><Message>m</Message></Error>`
	h := http.Header{}
	h.Set("x-ms-error-code", code)
	return &http.Response{StatusCode: status, Header: h, Body: io.NopCloser(strings.NewReader(body))}
}

// newTestProvider returns a provider that creates sessions through svc, with the locality
// policy in its pipeline as a real client's pipeline has.
func newTestProvider(t *testing.T, svc *sessionService) *ContainerSessionProvider {
	t.Helper()
	if svc.clock == nil {
		svc.clock = &providerClock{now: time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)}
	}
	azClient, err := azcore.NewClient("test", "v0.0.0", runtime.PipelineOptions{PerCall: []policy.Policy{locality.NewPolicy()}},
		&policy.ClientOptions{Transport: svc, Retry: policy.RetryOptions{MaxRetries: -1}})
	require.NoError(t, err)
	p := NewContainerSessionProviderFromClient(generated.NewServiceClient(providerServiceURL, azClient))
	p.now = svc.clock.Now
	return p
}

func blobRequest(t *testing.T, ctx context.Context, rawURL string) *http.Request {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	require.NoError(t, err)
	return req
}

func getToken(t *testing.T, p *ContainerSessionProvider, rawURL string) sessionCredential {
	t.Helper()
	cred, err := p.getSession(blobRequest(t, context.Background(), rawURL))
	require.NoError(t, err)
	return cred
}

func TestProviderIsRequestEligibleOnlyForGetBlob(t *testing.T) {
	p := NewContainerSessionProviderFromClient(nil)
	const blob = "https://testaccount.blob.core.windows.net/container/blob"
	tests := []struct {
		name     string
		method   string
		url      string
		header   map[string]string
		eligible bool
	}{
		{"GetBlob", http.MethodGet, blob, nil, true},
		{"GetBlobNestedName", http.MethodGet, "https://testaccount.blob.core.windows.net/container/dir/sub/blob.txt", nil, true},
		{"GetBlobSnapshot", http.MethodGet, blob + "?snapshot=2026-01-01T00:00:00Z", nil, true},
		{"GetBlobVersion", http.MethodGet, blob + "?versionid=2026-01-01T00:00:00Z", nil, true},
		{"GetBlobWithSAS", http.MethodGet, blob + "?sv=2026-12-06&sig=x", nil, true},
		{"GetBlobPathStyleIP", http.MethodGet, "http://127.0.0.1:10000/devstoreaccount1/container/blob", nil, true},
		{"GetBlobPathStyleEmulatorPort", http.MethodGet, "http://localhost:10000/devstoreaccount1/container/blob", nil, true},

		// the four operations the private drop also signed with a session
		{"GetBlobProperties", http.MethodHead, blob, nil, false},
		{"PutBlob", http.MethodPut, blob, map[string]string{"x-ms-blob-type": "BlockBlob"}, false},
		{"PutBlock", http.MethodPut, blob + "?comp=block&blockid=YQ==", nil, false},
		{"PutBlockList", http.MethodPut, blob + "?comp=blocklist", nil, false},

		// sub-resource reads: any comp, even an empty one
		{"GetLayout", http.MethodGet, blob + "?comp=layout", nil, false},
		{"GetTags", http.MethodGet, blob + "?comp=tags", nil, false},
		{"GetBlockList", http.MethodGet, blob + "?comp=blocklist", nil, false},
		{"GetPageRanges", http.MethodGet, blob + "?comp=pagelist", nil, false},
		{"EmptyComp", http.MethodGet, blob + "?comp=", nil, false},
		{"EmptyResType", http.MethodGet, blob + "?restype=", nil, false},

		// container- and service-level requests
		{"GetContainerProperties", http.MethodGet, "https://testaccount.blob.core.windows.net/container?restype=container", nil, false},
		{"ListBlobs", http.MethodGet, "https://testaccount.blob.core.windows.net/container?restype=container&comp=list", nil, false},
		{"ContainerURLWithoutQuery", http.MethodGet, "https://testaccount.blob.core.windows.net/container", nil, false},
		{"ServiceURL", http.MethodGet, "https://testaccount.blob.core.windows.net/", nil, false},
		{"ServiceProperties", http.MethodGet, "https://testaccount.blob.core.windows.net/?restype=service&comp=properties", nil, false},
		{"PathStyleContainer", http.MethodGet, "http://localhost:10000/devstoreaccount1/container", nil, false},

		// other verbs
		{"DeleteBlob", http.MethodDelete, blob, nil, false},
		{"PostBlob", http.MethodPost, blob, nil, false},
		{"CopyBlob", http.MethodPut, blob, map[string]string{"x-ms-copy-source": "https://src/c/b"}, false},

		// structured message, under its canonical or wire-cased key
		{"StructuredMessage", http.MethodGet, blob, map[string]string{"X-Ms-Structured-Body": "XSM/1.0; properties=crc64"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			require.NoError(t, err)
			for k, v := range tt.header {
				req.Header.Set(k, v)
			}
			require.Equal(t, tt.eligible, p.isRequestEligible(req))
		})
	}

	t.Run("StructuredMessageRawKey", func(t *testing.T) {
		// the generated clients write headers with wire casing directly into the map
		req, err := http.NewRequest(http.MethodGet, blob, nil)
		require.NoError(t, err)
		req.Header["x-ms-structured-body"] = []string{"XSM/1.0; properties=crc64"}
		require.False(t, p.isRequestEligible(req))
	})
	t.Run("NilRequest", func(t *testing.T) {
		require.False(t, p.isRequestEligible(nil))
		require.False(t, p.isRequestEligible(&http.Request{Method: http.MethodGet}))
	})
}

func TestProviderCreatesSessionWithHMAC(t *testing.T) {
	svc := &sessionService{}
	p := newTestProvider(t, svc)
	cred := getToken(t, p, "https://testaccount.blob.core.windows.net/container/blob")
	require.Equal(t, "container-1", cred.token)
	require.Equal(t, testSessionKey, cred.key)
	require.False(t, cred.fallback)

	calls := svc.allCalls()
	require.Len(t, calls, 1)
	query, err := url.ParseQuery(calls[0].query)
	require.NoError(t, err)
	require.Equal(t, url.Values{"restype": {"container"}, "comp": {"session"}}, query)
	require.Contains(t, calls[0].body, "<AuthenticationType>HMAC</AuthenticationType>")
}

func TestProviderCachesOneSessionPerContainer(t *testing.T) {
	svc := &sessionService{}
	p := newTestProvider(t, svc)

	a1 := getToken(t, p, "https://testaccount.blob.core.windows.net/a/blob1")
	a2 := getToken(t, p, "https://testaccount.blob.core.windows.net/a/dir/blob2")
	b1 := getToken(t, p, "https://testaccount.blob.core.windows.net/b/blob1")

	require.Equal(t, a1, a2, "blobs in one container share its session")
	require.NotEqual(t, a1.token, b1.token, "containers never share a session")
	require.Equal(t, 1, svc.callsTo("a"))
	require.Equal(t, 1, svc.callsTo("b"))
}

func TestProviderConcurrentColdRequestsCreateOneSession(t *testing.T) {
	svc := &sessionService{gate: make(chan struct{})}
	p := newTestProvider(t, svc)

	const callers = 50
	tokens := make(chan string, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cred, err := p.getSession(blobRequest(t, context.Background(), "https://testaccount.blob.core.windows.net/container/blob"))
			require.NoError(t, err)
			tokens <- cred.token
		}()
	}
	time.Sleep(20 * time.Millisecond)
	close(svc.gate)
	wg.Wait()
	close(tokens)

	require.Equal(t, 1, svc.callsTo("container"))
	for token := range tokens {
		require.Equal(t, "container-1", token)
	}
}

func TestProviderFallbackCachedForFiveMinutesPerContainer(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status int
		code   string
	}{
		{"InternalServerError", http.StatusInternalServerError, "InternalError"},
		{"ServiceUnavailable", http.StatusServiceUnavailable, "ServerBusy"},
		{"Forbidden", http.StatusForbidden, "AuthorizationPermissionMismatch"},
		{"FeatureNotEnabled", http.StatusBadRequest, "FeatureNotEnabled"},
		{"FeatureNotEnabledAnyCase", http.StatusBadRequest, "featurenotenabled"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			svc := &sessionService{respond: func(container string, n int) *http.Response {
				if container == "unavailable" && n == 1 {
					return errorResponse(tt.status, tt.code)
				}
				return nil
			}}
			p := newTestProvider(t, svc)
			const url = "https://testaccount.blob.core.windows.net/unavailable/blob"

			cred := getToken(t, p, url)
			require.True(t, cred.fallback, "the failure is a fallback decision, not an error")

			// Honored for the whole cooldown: no refresh before it ends.
			svc.clock.Advance(5*time.Minute - time.Second)
			for range 5 {
				require.True(t, getToken(t, p, url).fallback)
			}
			require.Equal(t, 1, svc.callsTo("unavailable"), "the service isn't hammered during the cooldown")

			// Other containers aren't affected.
			require.False(t, getToken(t, p, "https://testaccount.blob.core.windows.net/other/blob").fallback)

			// One foreground attempt when the cooldown ends.
			svc.clock.Advance(time.Second)
			cred = getToken(t, p, url)
			require.False(t, cred.fallback)
			require.Equal(t, "unavailable-2", cred.token)
			require.Equal(t, 2, svc.callsTo("unavailable"))
		})
	}
}

func TestProviderConcurrentFallbackCallsCreateOnce(t *testing.T) {
	svc := &sessionService{gate: make(chan struct{}), respond: func(string, int) *http.Response {
		return errorResponse(http.StatusServiceUnavailable, "ServerBusy")
	}}
	p := newTestProvider(t, svc)

	var fallbacks atomic.Int32
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cred, err := p.getSession(blobRequest(t, context.Background(), "https://testaccount.blob.core.windows.net/container/blob"))
			require.NoError(t, err)
			if cred.fallback {
				fallbacks.Add(1)
			}
		}()
	}
	time.Sleep(20 * time.Millisecond)
	close(svc.gate)
	wg.Wait()
	require.Equal(t, int32(20), fallbacks.Load())
	require.Equal(t, 1, svc.callsTo("container"))
}

func TestProviderOtherErrorsPropagateAndAreNotCached(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status int
		code   string
	}{
		{"BadRequestOtherCode", http.StatusBadRequest, "InvalidHeaderValue"},
		{"NotFound", http.StatusNotFound, "ContainerNotFound"},
		{"Unauthorized", http.StatusUnauthorized, "InvalidAuthenticationInfo"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			svc := &sessionService{respond: func(_ string, n int) *http.Response {
				if n == 1 {
					return errorResponse(tt.status, tt.code)
				}
				return nil
			}}
			p := newTestProvider(t, svc)
			const url = "https://testaccount.blob.core.windows.net/container/blob"

			_, err := p.getSession(blobRequest(t, context.Background(), url))
			var respErr *azcore.ResponseError
			require.ErrorAs(t, err, &respErr)
			require.Equal(t, tt.status, respErr.StatusCode)

			require.Equal(t, "container-2", getToken(t, p, url).token, "an error isn't cached")
		})
	}
}

func TestProviderIncompleteResponseIsAnError(t *testing.T) {
	svc := &sessionService{respond: func(string, int) *http.Response {
		body := `<?xml version="1.0" encoding="utf-8"?><CreateSessionResult><Id>id</Id></CreateSessionResult>`
		return &http.Response{StatusCode: http.StatusCreated, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body))}
	}}
	p := newTestProvider(t, svc)
	_, err := p.getSession(blobRequest(t, context.Background(), "https://testaccount.blob.core.windows.net/container/blob"))
	require.ErrorContains(t, err, "missing the session credentials")
}

func TestProviderRefreshesInBackgroundBeforeExpiry(t *testing.T) {
	svc := &sessionService{}
	p := newTestProvider(t, svc)
	const url = "https://testaccount.blob.core.windows.net/container/blob"
	require.Equal(t, "container-1", getToken(t, p, url).token)

	// Not yet within 30 seconds of expiry: no refresh.
	svc.clock.Advance(time.Hour - 31*time.Second)
	require.Equal(t, "container-1", getToken(t, p, url).token)
	require.Equal(t, 1, svc.callsTo("container"))

	// Within the window, the caller still gets the current session without waiting while the
	// refresh runs in the background.
	svc.mu.Lock()
	svc.gate = make(chan struct{})
	svc.mu.Unlock()
	svc.clock.Advance(time.Second)
	require.Equal(t, "container-1", getToken(t, p, url).token)
	require.Equal(t, "container-1", getToken(t, p, url).token)

	close(svc.gate)
	require.Eventually(t, func() bool { return getToken(t, p, url).token == "container-2" }, 5*time.Second, time.Millisecond)
	require.Equal(t, 2, svc.callsTo("container"))
}

func TestProviderInvalidateIsConditionalAndPerContainer(t *testing.T) {
	svc := &sessionService{}
	p := newTestProvider(t, svc)
	const a, b = "https://testaccount.blob.core.windows.net/a/blob", "https://testaccount.blob.core.windows.net/b/blob"
	used := getToken(t, p, a)
	getToken(t, p, b)

	// A stale credential doesn't discard the current session.
	p.invalidateSession(blobRequest(t, context.Background(), a), sessionCredential{token: "stale"})
	require.Equal(t, used, getToken(t, p, a))

	// Invalidating never creates a session itself.
	p.invalidateSession(blobRequest(t, context.Background(), a), used)
	require.Equal(t, 1, svc.callsTo("a"))

	// The next request creates a new session; the other container keeps its own.
	require.Equal(t, "a-2", getToken(t, p, a).token)
	require.Equal(t, "b-1", getToken(t, p, b).token)
	require.Equal(t, 1, svc.callsTo("b"))
}

func TestProviderConcurrentInvalidationsReplaceSessionOnce(t *testing.T) {
	svc := &sessionService{}
	p := newTestProvider(t, svc)
	const url = "https://testaccount.blob.core.windows.net/container/blob"
	rejected := getToken(t, p, url)

	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.invalidateSession(blobRequest(t, context.Background(), url), rejected)
			cred, err := p.getSession(blobRequest(t, context.Background(), url))
			require.NoError(t, err)
			require.NotEqual(t, rejected.token, cred.token)
		}()
	}
	wg.Wait()
	require.Equal(t, 2, svc.callsTo("container"), "one replacement, not one per rejected request")
}

func TestProviderCreatesSessionAtAccountEndpointForRoutedRequest(t *testing.T) {
	svc := &sessionService{}
	p := newTestProvider(t, svc)
	ctx := locality.WithEndpoint(context.Background(), "https://layout.blob.core.windows.net")
	_, err := p.getSession(blobRequest(t, ctx, "https://testaccount.blob.core.windows.net/container/blob"))
	require.NoError(t, err)

	calls := svc.allCalls()
	require.Len(t, calls, 1)
	require.Equal(t, "testaccount.blob.core.windows.net", calls[0].host, "Create Session must not follow the chunk's layout endpoint")
}

func TestProviderRequiresContainer(t *testing.T) {
	p := newTestProvider(t, &sessionService{})
	_, err := p.getSession(blobRequest(t, context.Background(), "https://testaccount.blob.core.windows.net/"))
	require.Error(t, err)
	_, err = p.getSession(nil)
	require.Error(t, err)
	// and invalidation of such a request is a harmless no-op
	p.invalidateSession(blobRequest(t, context.Background(), "https://testaccount.blob.core.windows.net/"), sessionCredential{})
}

func TestSessionModeResolution(t *testing.T) {
	require.Equal(t, SessionModeDisabled, ResolveSessionMode(SessionModeAuto), "Auto currently resolves to Disabled")
	require.Equal(t, SessionModeDisabled, ResolveSessionMode(SessionModeDisabled))
	require.Equal(t, SessionModeEnabled, ResolveSessionMode(SessionModeEnabled))
	require.Equal(t, SessionModeAuto, SessionMode(""), "Auto is the zero value")
	require.ElementsMatch(t, []SessionMode{SessionModeAuto, SessionModeDisabled, SessionModeEnabled}, PossibleSessionModeValues())

	for _, m := range PossibleSessionModeValues() {
		require.NoError(t, ValidateSessionMode(m))
	}
	require.Error(t, ValidateSessionMode("enabled"), "values are case-sensitive")
}
