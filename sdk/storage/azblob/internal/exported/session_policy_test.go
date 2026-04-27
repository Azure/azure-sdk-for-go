// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"context"
	"errors"
	"io"
	"net/http"
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

// TestNewSessionPolicy_Success tests successful creation of a session policy.
func TestNewSessionPolicy_Success(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	serviceClient := newTestServiceClient(t, srv)
	bearerPolicy := &mockBearerPolicy{}

	opts := SessionOptions{
		Mode:          SessionModeSingleSpecifiedContainer,
		AccountName:   "testaccount",
		ContainerName: "testcontainer",
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
				Mode:          SessionModeSingleSpecifiedContainer,
				AccountName:   "",
				ContainerName: "testcontainer",
			},
			expectedError: "account name is required",
		},
		{
			name: "MissingContainerName",
			opts: SessionOptions{
				Mode:          SessionModeSingleSpecifiedContainer,
				AccountName:   "testaccount",
				ContainerName: "",
			},
			expectedError: "container name is required",
		},
		{
			name: "UnsupportedMode",
			opts: SessionOptions{
				Mode:          SessionMode("unsupported"),
				AccountName:   "testaccount",
				ContainerName: "testcontainer",
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

			pol := &sessionPolicy{
				bearerTokenPolicy: bearerPolicy,
				opts: SessionOptions{
					AccountName:   "testaccount",
					ContainerName: "container",
				},
				resource: resource,
			}

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

// TestHandleSessionError_NonResponseError tests that non-ResponseError errors are passed through.
func TestHandleSessionError_NonResponseError(t *testing.T) {
	pol := &sessionPolicy{
		opts: SessionOptions{AccountName: "testaccount"},
	}

	originalErr := errors.New("some random error")
	resp := &http.Response{StatusCode: http.StatusOK}

	retResp, retErr := pol.handleSessionError(nil, resp, originalErr)
	require.Equal(t, resp, retResp)
	require.Equal(t, originalErr, retErr)
}

// TestHandleSessionError_ServiceUnavailable tests fallback to bearer on 503 with SessionOperationsTemporarilyUnavailable.
func TestHandleSessionError_ServiceUnavailable(t *testing.T) {
	pol := &sessionPolicy{
		opts: SessionOptions{AccountName: "testaccount"},
	}

	originalErr := &azcore.ResponseError{
		StatusCode: http.StatusServiceUnavailable,
		ErrorCode:  sessionUnavailable,
	}
	resp := &http.Response{StatusCode: http.StatusServiceUnavailable}

	retResp, retErr := pol.handleSessionError(nil, resp, originalErr)
	require.Nil(t, retResp)
	require.ErrorIs(t, retErr, errFallbackToBearer)
}

// TestHandleSessionError_OtherError tests that other errors are passed through.
func TestHandleSessionError_OtherError(t *testing.T) {
	pol := &sessionPolicy{
		opts: SessionOptions{AccountName: "testaccount"},
	}

	originalErr := &azcore.ResponseError{
		StatusCode: http.StatusNotFound,
		ErrorCode:  "BlobNotFound",
	}
	resp := &http.Response{StatusCode: http.StatusNotFound}

	retResp, retErr := pol.handleSessionError(nil, resp, originalErr)
	require.Equal(t, resp, retResp)
	require.Equal(t, originalErr, retErr)
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

	pol := &sessionPolicy{
		opts: SessionOptions{
			AccountName: "testaccount",
		},
		resource: resource,
	}

	// Create a helper policy to pass the request through handleSessionError
	testPolicy := &testRetryPolicy{
		pol: pol,
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
	originalErr  error
	originalResp *http.Response
}

func (p *testRetryPolicy) Do(req *policy.Request) (*http.Response, error) {
	return p.pol.handleSessionError(req, p.originalResp, p.originalErr)
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

// TestApplySessionReq_SetsAuthorizationHeader tests that applySessionReq sets the authorization header correctly.
func TestApplySessionReq_SetsAuthorizationHeader(t *testing.T) {
	sessionKey := "dGVzdC1rZXk=" // base64 encoded "test-key"
	sessionToken := "test-token"

	transport := &recordingTransport{}

	resource := newTestResource(sessionCredentials{
		key:   sessionKey,
		token: sessionToken,
	})

	pol := &sessionPolicy{
		opts: SessionOptions{
			AccountName: "testaccount",
		},
		resource: resource,
	}

	// Create a pipeline with our policy that will call applySessionReq
	testPolicy := &testApplyPolicy{pol: pol}

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

	// Verify the Authorization header was set
	authHeader := transport.lastRequest.Header.Get(shared.HeaderAuthorization)
	require.True(t, strings.HasPrefix(authHeader, "Session "))
	require.Contains(t, authHeader, sessionToken)
}

// testApplyPolicy is a helper policy that calls applySessionReq for testing.
type testApplyPolicy struct {
	pol *sessionPolicy
}

func (p *testApplyPolicy) Do(req *policy.Request) (*http.Response, error) {
	return p.pol.applySessionReq(req)
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

// TestIntegration_SessionPolicy_SuccessfulRequest tests the full flow of a successful session request.
func TestIntegration_SessionPolicy_SuccessfulRequest(t *testing.T) {
	sessionKey := "dGVzdC1rZXk=" // base64 encoded "test-key"
	sessionToken := "test-session-token"

	resource := newTestResource(sessionCredentials{
		key:   sessionKey,
		token: sessionToken,
	})

	transport := &recordingTransport{}

	bearerPolicy := &mockBearerPolicy{}

	pol := &sessionPolicy{
		bearerTokenPolicy: bearerPolicy,
		opts: SessionOptions{
			AccountName:   "testaccount",
			ContainerName: "testcontainer",
		},
		resource: resource,
	}

	// Create request through runtime to get proper Next() support
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

	// Verify session credentials were used
	require.Equal(t, 0, bearerPolicy.doCalls)

	// Verify Authorization header was set with Session prefix
	authHeader := transport.lastRequest.Header.Get(shared.HeaderAuthorization)
	require.True(t, strings.HasPrefix(authHeader, "Session "))
}

// TestIntegration_SessionPolicy_FallbackToBearer tests that non-session requests fallback to bearer.
func TestIntegration_SessionPolicy_FallbackToBearer(t *testing.T) {
	resource := newTestResource(sessionCredentials{
		key:   "dGVzdC1rZXk=",
		token: "test-token",
	})

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
		bearerTokenPolicy: bearerPolicy,
		opts: SessionOptions{
			AccountName:   "testaccount",
			ContainerName: "testcontainer",
		},
		resource: resource,
	}

	// POST request should not use session
	req := createTestPolicyRequest(t, http.MethodPost, "https://testaccount.blob.core.windows.net/testcontainer/blob")

	resp, err := pol.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify bearer was used, not session
	require.Equal(t, 1, bearerPolicy.doCalls)
}

// TestDoWithSession_ResourceError tests doWithSession when resource returns an error.
func TestDoWithSession_ResourceError(t *testing.T) {
	expectedErr := errors.New("resource error")
	resource := newTestResourceWithError(expectedErr)

	pol := &sessionPolicy{
		opts:     SessionOptions{AccountName: "testaccount"},
		resource: resource,
	}

	req := createTestPolicyRequest(t, http.MethodGet, "https://testaccount.blob.core.windows.net/container/blob")

	resp, err := pol.doWithSession(req)
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

	pol := &sessionPolicy{
		opts: SessionOptions{
			AccountName: "testaccount",
		},
		resource: resource,
	}

	// Create a helper policy to call doWithSession
	testPolicy := &testDoWithSessionPolicy{pol: pol}

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
	pol *sessionPolicy
}

func (p *testDoWithSessionPolicy) Do(req *policy.Request) (*http.Response, error) {
	return p.pol.doWithSession(req)
}

// TestApplySessionReq_EmptySessionKey tests applySessionReq with empty session key.
func TestApplySessionReq_EmptySessionKey(t *testing.T) {
	sessionToken := "test-token"

	transport := &recordingTransport{}

	resource := newTestResource(sessionCredentials{
		token: sessionToken,
	})

	pol := &sessionPolicy{
		opts: SessionOptions{
			AccountName: "testaccount",
		},
		resource: resource,
	}

	// Create a pipeline with applySessionReq policy
	testPolicy := &testApplyPolicy{pol: pol}

	pl := runtime.NewPipeline("test", "v1.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{testPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://testaccount.blob.core.windows.net/container/blob")
	require.NoError(t, err)

	// Empty key is valid base64 (decodes to empty bytes), so no error is expected
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify the Authorization header was set with Session prefix and the token
	authHeader := transport.lastRequest.Header.Get(shared.HeaderAuthorization)
	require.True(t, strings.HasPrefix(authHeader, "Session "))
	require.Contains(t, authHeader, sessionToken)
}

// TestApplySessionReq_NilSessionToken tests applySessionReq with nil session token.
func TestApplySessionReq_NilSessionToken(t *testing.T) {
	sessionKey := "dGVzdC1rZXk=" // base64 encoded "test-key"

	resource := newTestResource(sessionCredentials{
		key: sessionKey,
	})

	transport := &mockTransport{}

	pol := &sessionPolicy{
		opts: SessionOptions{
			AccountName: "testaccount",
		},
		resource: resource,
	}

	// Create a pipeline with our policy that will call applySessionReq
	testPolicy := &testApplyPolicy{pol: pol}

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
}

// TestDoWithSession_FallbackCreds tests doWithSession when resource returns fallback credentials.
func TestDoWithSession_FallbackCreds(t *testing.T) {
	resource := newTestResource(sessionCredentials{fallback: true})

	pol := &sessionPolicy{
		opts:     SessionOptions{AccountName: "testaccount"},
		resource: resource,
	}

	req := createTestPolicyRequest(t, http.MethodGet, "https://testaccount.blob.core.windows.net/container/blob")

	resp, err := pol.doWithSession(req)
	require.Nil(t, resp)
	require.ErrorIs(t, err, errFallbackToBearer)
}
