// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/stretchr/testify/require"
)

// newTestServiceClient creates a ServiceClient backed by a mock server for testing.
func newTestServiceClient(t *testing.T, srv *mock.Server) *generated.ServiceClient {
	azClient, err := azcore.NewClient("test", "v1.0.0", runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	require.NoError(t, err)
	return generated.NewServiceClient(srv.URL(), azClient)
}

// mockBearerPolicy is a mock bearer token policy for testing.
type mockBearerPolicy struct {
	doFn    func(req *policy.Request) (*http.Response, error)
	doCalls int
}

func (m *mockBearerPolicy) Do(req *policy.Request) (*http.Response, error) {
	m.doCalls++
	if m.doFn != nil {
		return m.doFn(req)
	}
	return &http.Response{StatusCode: http.StatusOK}, nil
}

// newTestResource creates a temporal.Resource for testing that returns the given credentials.
func newTestResource(creds sessionCredentials) *temporal.Resource[sessionCredentials, context.Context] {
	return temporal.NewResourceWithOptions(func(_ context.Context) (sessionCredentials, time.Time, error) {
		return creds, time.Now().Add(time.Hour), nil
	}, temporal.ResourceOptions[sessionCredentials, context.Context]{})
}

// newTestResourceWithError creates a temporal.Resource for testing that returns an error.
func newTestResourceWithError(err error) *temporal.Resource[sessionCredentials, context.Context] {
	return temporal.NewResourceWithOptions(func(_ context.Context) (sessionCredentials, time.Time, error) {
		return sessionCredentials{}, time.Time{}, err
	}, temporal.ResourceOptions[sessionCredentials, context.Context]{})
}

// newSessionPolicyWithResource creates a sessionPolicy for testing with a pre-populated resource for the given container.
func newSessionPolicyWithResource(opts SessionOptions, bearerPolicy policy.Policy, containerName string, resource *temporal.Resource[sessionCredentials, context.Context]) *sessionPolicy {
	pol := &sessionPolicy{
		bearerTokenPolicy: bearerPolicy,
		opts:              opts,
	}
	if resource != nil {
		pol.resources.Store(containerName, resource)
	}
	return pol
}

// TestNewSessionPolicy_Success tests successful creation of a session policy.
func TestNewSessionPolicy_Success(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	serviceClient := newTestServiceClient(t, srv)
	bearerPolicy := &mockBearerPolicy{}

	opts := SessionOptions{
		Mode:        SessionModeEnabled,
		AccountName: "testaccount",
	}

	pol, err := NewSessionPolicy(opts, bearerPolicy, serviceClient)
	require.NoError(t, err)
	require.NotNil(t, pol)
}

// TestNewSessionPolicy_Errors tests error cases when creating a session policy.
func TestNewSessionPolicy_Errors(t *testing.T) {
	tests := []struct {
		name          string
		opts          SessionOptions
		expectedError string
	}{
		{
			name: "MissingAccountName",
			opts: SessionOptions{
				Mode:        SessionModeEnabled,
				AccountName: "",
			},
			expectedError: "account name is required",
		},
		{
			name: "UnsupportedMode",
			opts: SessionOptions{
				Mode:        SessionMode("unsupported"),
				AccountName: "testaccount",
			},
			expectedError: "unsupported session mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeFn()

			serviceClient := newTestServiceClient(t, srv)
			bearerPolicy := &mockBearerPolicy{}

			pol, err := NewSessionPolicy(tt.opts, bearerPolicy, serviceClient)
			require.Error(t, err)
			require.Nil(t, pol)
			require.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

// TestNewSessionPolicy_ReturnsBearerPolicy tests that disabled and default modes return the bearer policy directly.
func TestNewSessionPolicy_ReturnsBearerPolicy(t *testing.T) {
	for _, mode := range []SessionMode{SessionModeDisabled, SessionModeDefault} {
		t.Run(string(mode), func(t *testing.T) {
			srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeFn()

			serviceClient := newTestServiceClient(t, srv)
			bearerPolicy := &mockBearerPolicy{}

			pol, err := NewSessionPolicy(SessionOptions{Mode: mode, AccountName: "testaccount"}, bearerPolicy, serviceClient)
			require.NoError(t, err)
			require.Same(t, bearerPolicy, pol)
		})
	}
}

// TestSessionPolicy_Do_FallbackToBearer tests scenarios where the session policy falls back to bearer token authentication.
func TestSessionPolicy_Do_FallbackToBearer(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		url                string
		useFallbackCreds   bool
		expectedBearerCall int
	}{
		{
			name:               "NonGetMethod",
			method:             http.MethodPost,
			url:                "https://testaccount.blob.core.windows.net/container/blob",
			useFallbackCreds:   false,
			expectedBearerCall: 1,
		},
		{
			name:               "CompParam",
			method:             http.MethodGet,
			url:                "https://testaccount.blob.core.windows.net/container/blob?comp=metadata",
			useFallbackCreds:   false,
			expectedBearerCall: 1,
		},
		{
			name:               "ContainerOnly",
			method:             http.MethodGet,
			url:                "https://testaccount.blob.core.windows.net/container",
			useFallbackCreds:   false,
			expectedBearerCall: 1,
		},
		{
			name:               "FallbackCredentials",
			method:             http.MethodGet,
			url:                "https://testaccount.blob.core.windows.net/container/blob",
			useFallbackCreds:   true,
			expectedBearerCall: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bearerPolicy := &mockBearerPolicy{
				doFn: func(req *policy.Request) (*http.Response, error) {
					return &http.Response{StatusCode: http.StatusOK}, nil
				},
			}

			var resource *temporal.Resource[sessionCredentials, context.Context]
			if tt.useFallbackCreds {
				resource = newTestResource(sessionCredentials{fallback: true})
			} else {
				resource = newTestResource(sessionCredentials{
					key:   "dGVzdC1rZXk=",
					token: "test-token",
				})
			}

			pol := newSessionPolicyWithResource(
				SessionOptions{AccountName: "testaccount"},
				bearerPolicy,
				"container",
				resource,
			)

			req := createTestPolicyRequest(t, tt.method, tt.url)

			resp, err := pol.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Equal(t, tt.expectedBearerCall, bearerPolicy.doCalls)
		})
	}
}

// TestCanUseSession tests the supportsSession helper function.
func TestCanUseSession(t *testing.T) {
	tests := []struct {
		name              string
		method            string
		urlStr            string
		expectedContainer string
		expectedOK        bool
	}{
		{
			name:              "ValidGETBlobRequest",
			method:            http.MethodGet,
			urlStr:            "https://account.blob.core.windows.net/container/blob",
			expectedContainer: "container",
			expectedOK:        true,
		},
		{
			name:              "ValidGETBlobRequestWithPath",
			method:            http.MethodGet,
			urlStr:            "https://account.blob.core.windows.net/container/path/to/blob",
			expectedContainer: "container",
			expectedOK:        true,
		},
		{
			name:              "NonGETMethod_POST",
			method:            http.MethodPost,
			urlStr:            "https://account.blob.core.windows.net/container/blob",
			expectedContainer: "",
			expectedOK:        false,
		},
		{
			name:              "NonGETMethod_PUT",
			method:            http.MethodPut,
			urlStr:            "https://account.blob.core.windows.net/container/blob",
			expectedContainer: "",
			expectedOK:        false,
		},
		{
			name:              "NonGETMethod_DELETE",
			method:            http.MethodDelete,
			urlStr:            "https://account.blob.core.windows.net/container/blob",
			expectedContainer: "",
			expectedOK:        false,
		},
		{
			name:              "RequestWithCompParam",
			method:            http.MethodGet,
			urlStr:            "https://account.blob.core.windows.net/container/blob?comp=metadata",
			expectedContainer: "",
			expectedOK:        false,
		},
		{
			name:              "EmptyPath",
			method:            http.MethodGet,
			urlStr:            "https://account.blob.core.windows.net/",
			expectedContainer: "",
			expectedOK:        false,
		},
		{
			name:              "ContainerOnly_NoBlob",
			method:            http.MethodGet,
			urlStr:            "https://account.blob.core.windows.net/container",
			expectedContainer: "",
			expectedOK:        false,
		},
		{
			name:              "ContainerOnly_TrailingSlash",
			method:            http.MethodGet,
			urlStr:            "https://account.blob.core.windows.net/container/",
			expectedContainer: "",
			expectedOK:        false,
		},
		{
			name:              "RootPath",
			method:            http.MethodGet,
			urlStr:            "https://account.blob.core.windows.net",
			expectedContainer: "",
			expectedOK:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.urlStr, nil)
			require.NoError(t, err)

			containerName, ok := supportsSession(req)
			require.Equal(t, tt.expectedOK, ok)
			require.Equal(t, tt.expectedContainer, containerName)
		})
	}
}

// TestHandleSessionError_PassThrough tests that errors that are not session-specific are passed through unchanged.
func TestHandleSessionError_PassThrough(t *testing.T) {
	tests := []struct {
		name string
		err  error
		resp *http.Response
	}{
		{
			name: "NonResponseError",
			err:  errors.New("some random error"),
			resp: &http.Response{StatusCode: http.StatusOK},
		},
		{
			name: "NotFound",
			err:  &azcore.ResponseError{StatusCode: http.StatusNotFound, ErrorCode: "BlobNotFound"},
			resp: &http.Response{StatusCode: http.StatusNotFound},
		},
		{
			name: "503_NonSessionErrorCode",
			err:  &azcore.ResponseError{StatusCode: http.StatusServiceUnavailable, ErrorCode: "ServerBusy"},
			resp: &http.Response{StatusCode: http.StatusServiceUnavailable},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pol := &sessionPolicy{
				opts: SessionOptions{AccountName: "testaccount"},
			}
			resource := newTestResource(sessionCredentials{})

			retResp, retErr := pol.handleSessionError(nil, tt.resp, tt.err, resource)
			require.Equal(t, tt.resp, retResp)
			require.Equal(t, tt.err, retErr)
		})
	}
}

// TestHandleSessionError_FallbackToBearer tests that specific error responses trigger fallback to bearer token auth.
func TestHandleSessionError_FallbackToBearer(t *testing.T) {
	tests := []struct {
		name string
		err  *azcore.ResponseError
		resp *http.Response
	}{
		{
			name: "503_SessionUnavailable",
			err: &azcore.ResponseError{
				StatusCode: http.StatusServiceUnavailable,
				ErrorCode:  sessionUnavailable,
			},
			resp: &http.Response{StatusCode: http.StatusServiceUnavailable},
		},
		{
			name: "403_SessionSchemeNotSupported",
			err: &azcore.ResponseError{
				StatusCode: http.StatusForbidden,
				ErrorCode:  "AuthenticationFailed",
				RawResponse: &http.Response{
					StatusCode: http.StatusForbidden,
					Status:     "403 Forbidden",
					Header:     http.Header{"Content-Type": []string{"application/xml"}},
					Body:       io.NopCloser(strings.NewReader("<?xml version=\"1.0\" encoding=\"utf-8\"?><Error><Code>AuthenticationFailed</Code><Message>Authentication scheme Session is not supported.</Message></Error>")),
					Request:    &http.Request{Method: http.MethodGet, URL: &url.URL{Scheme: "https", Host: "testaccount.blob.core.windows.net", Path: "/container/blob"}},
				},
			},
			resp: &http.Response{StatusCode: http.StatusForbidden},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pol := &sessionPolicy{
				opts: SessionOptions{AccountName: "testaccount"},
			}
			resource := newTestResource(sessionCredentials{})

			retResp, retErr := pol.handleSessionError(nil, tt.resp, tt.err, resource)
			require.Nil(t, retResp)
			require.ErrorIs(t, retErr, errFallbackToBearer)
		})
	}
}

// TestHandleSessionError_Unauthorized_TriggersRetry tests that a 401 response triggers retry with a new session.
func TestHandleSessionError_Unauthorized_TriggersRetry(t *testing.T) {
	sessionKey := "dGVzdC1rZXk=" // base64 encoded "test-key"
	sessionToken := "new-token"

	resource := newTestResource(sessionCredentials{
		key:   sessionKey,
		token: sessionToken,
	})

	transport := &mockTransport{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
		},
	}

	pol := newSessionPolicyWithResource(
		SessionOptions{AccountName: "testaccount"},
		nil,
		"container",
		resource,
	)

	// Create a helper policy to pass the request through handleSessionError
	testPolicy := &testRetryPolicy{
		pol:      pol,
		resource: resource,
	}

	pl := runtime.NewPipeline("test", "v1.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{testPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://testaccount.blob.core.windows.net/container/blob")
	require.NoError(t, err)

	originalErr := &azcore.ResponseError{
		StatusCode: http.StatusUnauthorized,
		ErrorCode:  "AuthenticationFailed",
	}
	unauthorizedResp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Header:     make(http.Header),
	}

	testPolicy.originalErr = originalErr
	testPolicy.originalResp = unauthorizedResp

	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

// testRetryPolicy is a helper policy for testing handleSessionError with 401.
type testRetryPolicy struct {
	pol          *sessionPolicy
	resource     *temporal.Resource[sessionCredentials, context.Context]
	originalErr  error
	originalResp *http.Response
}

func (p *testRetryPolicy) Do(req *policy.Request) (*http.Response, error) {
	return p.pol.handleSessionError(req, p.originalResp, p.originalErr, p.resource)
}

// createTestPolicyRequest creates a policy.Request for testing with Next() support.
func createTestPolicyRequest(t *testing.T, method, urlStr string) *policy.Request {
	httpReq, err := http.NewRequestWithContext(context.Background(), method, urlStr, nil)
	require.NoError(t, err)

	// Create a minimal pipeline for testing
	_ = runtime.NewPipeline("test", "v1.0.0", runtime.PipelineOptions{}, &policy.ClientOptions{
		Transport: &mockTransport{},
	})

	req, err := runtime.NewRequest(context.Background(), method, urlStr)
	require.NoError(t, err)
	req.Raw().Header = httpReq.Header

	return req
}

// mockTransport is a mock HTTP transport for testing.
type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) Do(_ *http.Request) (*http.Response, error) {
	if m.response != nil {
		return m.response, m.err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}, nil
}

// TestApplySessionReq_SetsHeaders tests that applySessionReq sets the authorization and x-ms-date headers correctly.
func TestApplySessionReq_SetsHeaders(t *testing.T) {
	sessionKey := "dGVzdC1rZXk=" // base64 encoded "test-key"
	sessionToken := "test-token"

	transport := &recordingTransport{}

	resource := newTestResource(sessionCredentials{
		key:   sessionKey,
		token: sessionToken,
	})

	pol := newSessionPolicyWithResource(
		SessionOptions{AccountName: "testaccount"},
		nil,
		"container",
		resource,
	)

	// Create a pipeline with our policy that will call applySessionReq
	testPolicy := &testApplyPolicy{pol: pol, resource: resource}

	pl := runtime.NewPipeline("test", "v1.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{testPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	before := time.Now().UTC()
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://testaccount.blob.core.windows.net/container/blob")
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify the Authorization header was set
	authHeader := transport.lastRequest.Header.Get(shared.HeaderAuthorization)
	require.True(t, strings.HasPrefix(authHeader, "Session "))
	require.Contains(t, authHeader, sessionToken)

	// Verify x-ms-date header was set to a recent time
	dateStr := transport.lastRequest.Header.Get(shared.HeaderXmsDate)
	require.NotEmpty(t, dateStr)
	parsedDate, err := time.Parse(http.TimeFormat, dateStr)
	require.NoError(t, err)
	require.False(t, parsedDate.Before(before.Add(-1*time.Second)))
}

// testApplyPolicy is a helper policy that calls applySessionReq for testing.
type testApplyPolicy struct {
	pol      *sessionPolicy
	resource *temporal.Resource[sessionCredentials, context.Context]
}

func (p *testApplyPolicy) Do(req *policy.Request) (*http.Response, error) {
	return p.pol.applySessionReq(req, p.resource)
}

// recordingTransport records the last request for verification.
type recordingTransport struct {
	lastRequest *http.Request
}

func (r *recordingTransport) Do(req *http.Request) (*http.Response, error) {
	r.lastRequest = req
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}, nil
}

// TestDoWithSession_ResourceError tests doWithSession when resource returns an error.
func TestDoWithSession_ResourceError(t *testing.T) {
	expectedErr := errors.New("resource error")
	resource := newTestResourceWithError(expectedErr)

	pol := newSessionPolicyWithResource(
		SessionOptions{AccountName: "testaccount"},
		nil,
		"container",
		resource,
	)

	req := createTestPolicyRequest(t, http.MethodGet, "https://testaccount.blob.core.windows.net/container/blob")

	resp, err := pol.doWithSession(req, "container")
	require.Nil(t, resp)
	require.Equal(t, expectedErr, err)
}

// TestDoWithSession_Success tests successful doWithSession flow.
func TestDoWithSession_Success(t *testing.T) {
	sessionKey := "dGVzdC1rZXk=" // base64 encoded "test-key"
	sessionToken := "test-token"

	resource := newTestResource(sessionCredentials{
		key:   sessionKey,
		token: sessionToken,
	})

	transport := &mockTransport{}

	pol := newSessionPolicyWithResource(
		SessionOptions{AccountName: "testaccount"},
		nil,
		"container",
		resource,
	)

	// Create a helper policy to call doWithSession
	testPolicy := &testDoWithSessionPolicy{pol: pol, containerName: "container"}

	pl := runtime.NewPipeline("test", "v1.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{testPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://testaccount.blob.core.windows.net/container/blob")
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

// testDoWithSessionPolicy is a helper policy that calls doWithSession for testing.
type testDoWithSessionPolicy struct {
	pol           *sessionPolicy
	containerName string
}

func (p *testDoWithSessionPolicy) Do(req *policy.Request) (*http.Response, error) {
	return p.pol.doWithSession(req, p.containerName)
}

// TestDoWithSession_FallbackCreds tests doWithSession when resource returns fallback credentials.
func TestDoWithSession_FallbackCreds(t *testing.T) {
	resource := newTestResource(sessionCredentials{fallback: true})

	pol := newSessionPolicyWithResource(
		SessionOptions{AccountName: "testaccount"},
		nil,
		"container",
		resource,
	)

	req := createTestPolicyRequest(t, http.MethodGet, "https://testaccount.blob.core.windows.net/container/blob")

	resp, err := pol.doWithSession(req, "container")
	require.Nil(t, resp)
	require.ErrorIs(t, err, errFallbackToBearer)
}

// TestGetOrCreateResource_CreatesNewResource tests that the policy creates a resource via CreateSession
// when no pre-populated resource exists for the container.
func TestGetOrCreateResource_CreatesNewResource(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	expiration := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	sessionKey := "dGVzdC1rZXk="
	sessionToken := "test-session-token"

	// The mock server will handle the CreateSession POST request
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusCreated),
		mock.WithBody(createSessionResponseXML(sessionKey, sessionToken, expiration)),
	)

	serviceClient := newTestServiceClient(t, srv)

	pol := &sessionPolicy{
		bearerTokenPolicy:  &mockBearerPolicy{},
		opts:               SessionOptions{AccountName: "testaccount"},
		oauthServiceClient: serviceClient,
	}

	// No resource pre-populated — getOrCreateResource should create one
	resource := pol.getOrCreateResource("testcontainer")
	require.NotNil(t, resource)

	// Verify the resource can be used to get session credentials
	creds, err := resource.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, sessionKey, creds.key)
	require.Equal(t, sessionToken, creds.token)
	require.False(t, creds.fallback)
}

// TestGetOrCreateResource_Caching tests that getOrCreateResource reuses resources for the same
// container and creates separate resources for different containers.
func TestGetOrCreateResource_Caching(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	expiration := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	srv.SetResponse(
		mock.WithStatusCode(http.StatusCreated),
		mock.WithBody(createSessionResponseXML("dGVzdC1rZXk=", "test-token", expiration)),
	)

	serviceClient := newTestServiceClient(t, srv)

	pol := &sessionPolicy{
		bearerTokenPolicy:  &mockBearerPolicy{},
		opts:               SessionOptions{AccountName: "testaccount"},
		oauthServiceClient: serviceClient,
	}

	t.Run("SameContainerReusesResource", func(t *testing.T) {
		r1 := pol.getOrCreateResource("container1")
		r2 := pol.getOrCreateResource("container1")
		require.Same(t, r1, r2)
	})

	t.Run("DifferentContainersSeparateResources", func(t *testing.T) {
		rA := pol.getOrCreateResource("containerA")
		rB := pol.getOrCreateResource("containerB")
		require.NotSame(t, rA, rB)
	})
}

// TestDo_CreatesResourceAndSignsRequest tests the full Do path where the policy creates a resource
// on the fly (no pre-populated resource) and uses it to sign the request.
func TestDo_CreatesResourceAndSignsRequest(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	expiration := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	sessionKey := "dGVzdC1rZXk="
	sessionToken := "test-session-token"

	// CreateSession response
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusCreated),
		mock.WithBody(createSessionResponseXML(sessionKey, sessionToken, expiration)),
	)

	serviceClient := newTestServiceClient(t, srv)
	bearerPolicy := &mockBearerPolicy{}

	pol := &sessionPolicy{
		bearerTokenPolicy:  bearerPolicy,
		opts:               SessionOptions{AccountName: "testaccount"},
		oauthServiceClient: serviceClient,
	}

	transport := &recordingTransport{}

	pl := runtime.NewPipeline("test", "v1.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{pol},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://testaccount.blob.core.windows.net/testcontainer/blob")
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Bearer policy should NOT have been called — session was used
	require.Equal(t, 0, bearerPolicy.doCalls)

	// Verify the Authorization header has the Session prefix with the token
	authHeader := transport.lastRequest.Header.Get(shared.HeaderAuthorization)
	require.True(t, strings.HasPrefix(authHeader, "Session "))
	require.Contains(t, authHeader, sessionToken)
}

// TestDo_CreateSessionError_FallsBackToBearer tests the full Do path where CreateSession fails
// and the policy falls back to bearer token authentication.
func TestDo_CreateSessionError_FallsBackToBearer(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		errorCode  string
	}{
		{
			name:       "InternalServerError",
			statusCode: http.StatusInternalServerError,
			errorCode:  "InternalError",
		},
		{
			name:       "FeatureNotEnabled",
			statusCode: http.StatusBadRequest,
			errorCode:  featureNotEnabled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeFn()

			srv.SetResponse(
				mock.WithStatusCode(tt.statusCode),
				mock.WithHeader("x-ms-error-code", tt.errorCode),
				mock.WithBody(createErrorResponseXML(tt.errorCode, "error")),
			)

			serviceClient := newTestServiceClient(t, srv)

			bearerPolicy := &mockBearerPolicy{
				doFn: func(req *policy.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("")),
						Header:     make(http.Header),
					}, nil
				},
			}

			pol := &sessionPolicy{
				bearerTokenPolicy:  bearerPolicy,
				opts:               SessionOptions{AccountName: "testaccount"},
				oauthServiceClient: serviceClient,
			}

			pl := runtime.NewPipeline("test", "v1.0.0", runtime.PipelineOptions{
				PerCall: []policy.Policy{pol},
			}, &policy.ClientOptions{
				Transport: srv,
			})

			req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://testaccount.blob.core.windows.net/testcontainer/blob")
			require.NoError(t, err)

			resp, err := pl.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Equal(t, 1, bearerPolicy.doCalls)
		})
	}
}

// TestDo_MultipleContainers_UsesSeparateSessions tests that requests to different containers
// use separate session credentials.
func TestDo_MultipleContainers_UsesSeparateSessions(t *testing.T) {
	bearerPolicy := &mockBearerPolicy{}

	resourceA := newTestResource(sessionCredentials{
		key:   "dGVzdC1rZXk=",
		token: "token-A",
	})
	resourceB := newTestResource(sessionCredentials{
		key:   "dGVzdC1rZXk=",
		token: "token-B",
	})

	pol := &sessionPolicy{
		bearerTokenPolicy: bearerPolicy,
		opts:              SessionOptions{AccountName: "testaccount"},
	}
	pol.resources.Store("containerA", resourceA)
	pol.resources.Store("containerB", resourceB)

	transportA := &recordingTransport{}
	plA := runtime.NewPipeline("test", "v1.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{pol},
	}, &policy.ClientOptions{
		Transport: transportA,
	})

	reqA, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://testaccount.blob.core.windows.net/containerA/blob")
	require.NoError(t, err)
	resp, err := plA.Do(reqA)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	authA := transportA.lastRequest.Header.Get(shared.HeaderAuthorization)
	require.Contains(t, authA, "token-A")

	transportB := &recordingTransport{}
	plB := runtime.NewPipeline("test", "v1.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{pol},
	}, &policy.ClientOptions{
		Transport: transportB,
	})

	reqB, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://testaccount.blob.core.windows.net/containerB/blob")
	require.NoError(t, err)
	resp, err = plB.Do(reqB)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	authB := transportB.lastRequest.Header.Get(shared.HeaderAuthorization)
	require.Contains(t, authB, "token-B")

	// Verify bearer was never called
	require.Equal(t, 0, bearerPolicy.doCalls)
}
