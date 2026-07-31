// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/stretchr/testify/require"
)

const (
	testAccountName = "testaccount"
	// testSessionKey is the base64 encoding of "test-key"; session keys are shared keys so they
	// must be valid base64.
	testSessionKey   = "dGVzdC1rZXk="
	testSessionToken = "test-session-token"
	testBlobURL      = "https://testaccount.blob.core.windows.net/container/blob"
)

// fakeSessionProvider is a SessionProvider whose behavior is fully controlled by the test.
type fakeSessionProvider struct {
	eligible bool
	// creds are returned by GetSession in order; the last entry is returned for any subsequent call.
	creds         []SessionCredential
	getErr        error
	invalidateErr error

	getCalls        int
	invalidateCalls int
	invalidatedWith SessionCredential
}

func (f *fakeSessionProvider) GetSession(*http.Request) (SessionCredential, error) {
	f.getCalls++
	if f.getErr != nil {
		return SessionCredential{}, f.getErr
	}
	if len(f.creds) == 0 {
		return SessionCredential{}, nil
	}
	i := f.getCalls - 1
	if i >= len(f.creds) {
		i = len(f.creds) - 1
	}
	return f.creds[i], nil
}

func (f *fakeSessionProvider) InvalidateSession(_ *http.Request, current SessionCredential) error {
	f.invalidateCalls++
	f.invalidatedWith = current
	return f.invalidateErr
}

func (f *fakeSessionProvider) IsRequestEligible(*http.Request) bool {
	return f.eligible
}

// newEligibleProvider returns a provider that hands out a usable session credential.
func newEligibleProvider() *fakeSessionProvider {
	return &fakeSessionProvider{
		eligible: true,
		creds:    []SessionCredential{NewSessionCredential(testSessionToken, testSessionKey, time.Now().Add(time.Hour))},
	}
}

// mockBearerPolicy is a stand-in for the bearer token policy.
type mockBearerPolicy struct {
	doCalls int
	// bodySeen records the request body observed by the policy, which is how tests verify that
	// the body was rewound before falling back.
	bodySeen []byte
	doFn     func(req *policy.Request) (*http.Response, error)
}

func (m *mockBearerPolicy) Do(req *policy.Request) (*http.Response, error) {
	m.doCalls++
	if body := req.Body(); body != nil {
		m.bodySeen, _ = io.ReadAll(body)
	}
	if m.doFn != nil {
		return m.doFn(req)
	}
	return newTestResponse(http.StatusOK), nil
}

// recordingTransport records the last request it saw and returns a canned response.
type recordingTransport struct {
	lastRequest *http.Request
	bodySeen    []byte
	calls       int
	statusCode  int
}

func (r *recordingTransport) Do(req *http.Request) (*http.Response, error) {
	r.calls++
	r.lastRequest = req
	if req.Body != nil {
		r.bodySeen, _ = io.ReadAll(req.Body)
	}
	code := r.statusCode
	if code == 0 {
		code = http.StatusOK
	}
	resp := newTestResponse(code)
	resp.Request = req
	return resp, nil
}

// errorTransport always fails, simulating a network-level failure.
type errorTransport struct {
	err   error
	calls int
}

func (e *errorTransport) Do(*http.Request) (*http.Response, error) {
	e.calls++
	return nil, e.err
}

// faultPolicy sits below the session policy and converts the response it receives into an
// *azcore.ResponseError. The pipeline itself never turns a status code into an error; that
// happens in the generated client code above the pipeline, so tests have to simulate it to
// exercise the session policy's error handling.
type faultPolicy struct {
	statusCode int
	errorCode  string
	calls      int
}

func (f *faultPolicy) Do(req *policy.Request) (*http.Response, error) {
	f.calls++
	resp, err := req.Next()
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != f.statusCode {
		return resp, nil
	}
	return resp, &azcore.ResponseError{
		StatusCode:  f.statusCode,
		ErrorCode:   f.errorCode,
		RawResponse: resp,
	}
}

func newTestResponse(statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader("")),
	}
}

// newTestPipeline builds a pipeline whose first policy is the session policy under test. Any
// additional policies run after it, and retries are disabled so each call reaches the transport
// exactly once.
func newTestPipeline(sessionPol policy.Policy, transport policy.Transporter, extra ...policy.Policy) runtime.Pipeline {
	return runtime.NewPipeline("test", "v1.0.0", runtime.PipelineOptions{
		PerCall: append([]policy.Policy{sessionPol}, extra...),
	}, &policy.ClientOptions{
		Transport: transport,
		Retry:     policy.RetryOptions{MaxRetries: -1},
	})
}

func newTestPolicyRequest(t *testing.T, method, rawURL string) *policy.Request {
	t.Helper()
	req, err := runtime.NewRequest(context.Background(), method, rawURL)
	require.NoError(t, err)
	return req
}

func TestSessionPolicyIneligibleRequestUsesBearer(t *testing.T) {
	provider := &fakeSessionProvider{eligible: false}
	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport)

	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodPut, testBlobURL))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, 1, bearer.doCalls)
	require.Equal(t, 0, provider.getCalls, "an ineligible request must not acquire a session")
	require.Equal(t, 0, provider.invalidateCalls)
}

func TestSessionPolicyGetSessionErrorPropagates(t *testing.T) {
	expectedErr := errors.New("failed to acquire session")
	provider := &fakeSessionProvider{eligible: true, getErr: expectedErr}
	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport)

	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
	require.Equal(t, 0, bearer.doCalls)
	require.Equal(t, 0, transport.calls)
}

func TestSessionPolicyFallbackCredentialUsesBearer(t *testing.T) {
	provider := &fakeSessionProvider{
		eligible: true,
		creds:    []SessionCredential{NewSessionCredentialFallback(time.Now().Add(time.Hour))},
	}
	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport)

	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, 1, provider.getCalls)
	require.Equal(t, 1, bearer.doCalls)
	require.Equal(t, 0, provider.invalidateCalls, "a fallback credential is not a rejected session")
}

func TestSessionPolicySignsRequest(t *testing.T) {
	provider := newEligibleProvider()
	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport)

	before := time.Now().UTC().Add(-time.Second).Truncate(time.Second)
	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, 1, provider.getCalls)
	require.Equal(t, 0, bearer.doCalls, "an eligible request must not use the bearer token policy")
	require.Equal(t, 1, transport.calls)

	authHeader := transport.lastRequest.Header.Get(shared.HeaderAuthorization)
	require.True(t, strings.HasPrefix(authHeader, "Session "), "unexpected auth header %q", authHeader)

	token, signature, found := strings.Cut(strings.TrimPrefix(authHeader, "Session "), ":")
	require.True(t, found, "auth header must be of the form 'Session <token>:<signature>'")
	require.Equal(t, testSessionToken, token)
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	require.NoError(t, err, "the signature must be base64 encoded")
	require.Len(t, sigBytes, 32, "the signature must be an HMAC-SHA256")

	dateHeader := transport.lastRequest.Header.Get(shared.HeaderXmsDate)
	require.NotEmpty(t, dateHeader)
	parsedDate, err := time.Parse(http.TimeFormat, dateHeader)
	require.NoError(t, err)
	require.False(t, parsedDate.Before(before))
}

func TestSessionPolicySignatureMatchesSharedKeySignature(t *testing.T) {
	provider := newEligibleProvider()
	transport := &recordingTransport{}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, &mockBearerPolicy{}), transport)

	_, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.NoError(t, err)

	// recompute the signature over the request the transport observed; it must match the one the
	// policy produced using the session key
	cred, err := NewSharedKeyCredential(testAccountName, testSessionKey)
	require.NoError(t, err)
	stringToSign, err := cred.buildStringToSign(transport.lastRequest)
	require.NoError(t, err)
	expectedSignature, err := cred.computeHMACSHA256(stringToSign)
	require.NoError(t, err)

	require.Equal(t, "Session "+testSessionToken+":"+expectedSignature, transport.lastRequest.Header.Get(shared.HeaderAuthorization))
}

// The configured account name is used verbatim when signing; it is never re-derived from the
// request URL. Configuring the wrong account therefore produces a signature the service rejects.
func TestSessionPolicyUsesConfiguredAccountNameForSignature(t *testing.T) {
	const wrongAccountName = "wrongaccount"

	provider := newEligibleProvider()
	transport := &recordingTransport{}

	// the request targets testaccount, but the policy is configured with a different account
	pl := newTestPipeline(NewSessionPolicy(wrongAccountName, provider, &mockBearerPolicy{}), transport)

	_, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.NoError(t, err)

	authHeader := transport.lastRequest.Header.Get(shared.HeaderAuthorization)

	// the signature matches the configured account name...
	wrongCred, err := NewSharedKeyCredential(wrongAccountName, testSessionKey)
	require.NoError(t, err)
	stringToSign, err := wrongCred.buildStringToSign(transport.lastRequest)
	require.NoError(t, err)
	require.Contains(t, stringToSign, "/"+wrongAccountName+"/", "the canonicalized resource uses the configured account name")
	wrongSignature, err := wrongCred.computeHMACSHA256(stringToSign)
	require.NoError(t, err)
	require.Equal(t, "Session "+testSessionToken+":"+wrongSignature, authHeader)

	// ...and therefore does not match the account the request is actually addressed to
	rightCred, err := NewSharedKeyCredential(testAccountName, testSessionKey)
	require.NoError(t, err)
	rightStringToSign, err := rightCred.buildStringToSign(transport.lastRequest)
	require.NoError(t, err)
	rightSignature, err := rightCred.computeHMACSHA256(rightStringToSign)
	require.NoError(t, err)
	require.NotEqual(t, rightSignature, wrongSignature, "a mismatched account name yields a different signature")
}

func TestSessionPolicyRefreshesDateHeader(t *testing.T) {
	provider := newEligibleProvider()
	transport := &recordingTransport{}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, &mockBearerPolicy{}), transport)

	req := newTestPolicyRequest(t, http.MethodGet, testBlobURL)
	stale := time.Now().UTC().Add(-time.Hour).Format(http.TimeFormat)
	req.Raw().Header.Set(shared.HeaderXmsDate, stale)

	_, err := pl.Do(req)
	require.NoError(t, err)

	require.NotEqual(t, stale, transport.lastRequest.Header.Get(shared.HeaderXmsDate), "the date must be refreshed on every attempt")
}

func TestSessionPolicyInvalidSessionKeyReturnsError(t *testing.T) {
	provider := &fakeSessionProvider{
		eligible: true,
		creds:    []SessionCredential{NewSessionCredential(testSessionToken, "not-base64!", time.Now().Add(time.Hour))},
	}
	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport)

	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, 0, transport.calls)
	require.Equal(t, 0, bearer.doCalls)
}

// A 401 response is the service rejecting the session, even though the pipeline surfaces it
// without an error. The session must be discarded and the request retried with bearer auth.
func TestSessionPolicyInvalidatesSessionOnUnauthorizedResponse(t *testing.T) {
	provider := newEligibleProvider()
	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{statusCode: http.StatusUnauthorized}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport)

	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "the bearer token response is returned")

	require.Equal(t, 1, provider.invalidateCalls, "the rejected session must be discarded")
	require.Equal(t, testSessionToken, provider.invalidatedWith.Token())
	require.Equal(t, 1, bearer.doCalls)
}

// A non-401 response without an error is returned unchanged.
func TestSessionPolicyNonUnauthorizedResponseIsReturned(t *testing.T) {
	provider := newEligibleProvider()
	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{statusCode: http.StatusNotFound}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport)

	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	require.Equal(t, 0, provider.invalidateCalls)
	require.Equal(t, 0, bearer.doCalls)
}

// A transport error is not a rejected session, so it propagates untouched.
func TestSessionPolicyTransportErrorPropagates(t *testing.T) {
	expectedErr := errors.New("connection reset")
	provider := newEligibleProvider()
	bearer := &mockBearerPolicy{}
	transport := &errorTransport{err: expectedErr}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport)

	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
	require.Equal(t, 0, provider.invalidateCalls)
	require.Equal(t, 0, bearer.doCalls)
}

// An error takes precedence over the status code: the request already failed, so the session is
// left alone and the error is returned to the caller.
func TestSessionPolicyUnauthorizedWithErrorIsReturned(t *testing.T) {
	provider := newEligibleProvider()
	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{statusCode: http.StatusUnauthorized}
	fault := &faultPolicy{statusCode: http.StatusUnauthorized, errorCode: "AuthenticationFailed"}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport, fault)

	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.Error(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	require.Equal(t, 0, provider.invalidateCalls)
	require.Equal(t, 0, bearer.doCalls)
}

func TestSessionPolicyRewindsBodyBeforeBearerFallback(t *testing.T) {
	body := []byte("request payload")

	provider := newEligibleProvider()
	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{statusCode: http.StatusUnauthorized}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport)

	req := newTestPolicyRequest(t, http.MethodGet, testBlobURL)
	require.NoError(t, req.SetBody(streaming.NopCloser(bytes.NewReader(body)), "application/octet-stream"))

	_, err := pl.Do(req)
	require.NoError(t, err)

	require.Equal(t, body, transport.bodySeen, "the session-authenticated attempt sends the body")
	require.Equal(t, body, bearer.bodySeen, "the body must be rewound before falling back to bearer auth")
}

func TestSessionPolicyNonUnauthorizedErrorPropagates(t *testing.T) {
	provider := newEligibleProvider()
	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{statusCode: http.StatusForbidden}
	fault := &faultPolicy{statusCode: http.StatusForbidden, errorCode: "AuthorizationFailure"}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport, fault)

	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.Error(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	var respErr *azcore.ResponseError
	require.True(t, errors.As(err, &respErr))
	require.Equal(t, http.StatusForbidden, respErr.StatusCode)

	require.Equal(t, 0, provider.invalidateCalls, "only a rejected session is invalidated")
	require.Equal(t, 0, bearer.doCalls)
}

func TestSessionPolicyInvalidateSessionErrorPropagates(t *testing.T) {
	expectedErr := errors.New("invalidate failed")
	provider := newEligibleProvider()
	provider.invalidateErr = expectedErr

	bearer := &mockBearerPolicy{}
	transport := &recordingTransport{statusCode: http.StatusUnauthorized}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, bearer), transport)

	resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
	require.Equal(t, 0, bearer.doCalls)
}

func TestSessionPolicyAcquiresSessionPerRequest(t *testing.T) {
	// the provider owns caching, so the policy asks for a session on every eligible request
	provider := newEligibleProvider()
	transport := &recordingTransport{}

	pl := newTestPipeline(NewSessionPolicy(testAccountName, provider, &mockBearerPolicy{}), transport)

	for i := 0; i < 3; i++ {
		resp, err := pl.Do(newTestPolicyRequest(t, http.MethodGet, testBlobURL))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	}

	require.Equal(t, 3, provider.getCalls)
	require.Equal(t, 3, transport.calls)
}
