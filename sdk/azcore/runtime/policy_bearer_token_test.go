// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"

	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

const (
	tokenValue                = "***"
	accessTokenRespSuccess    = `{"access_token": "` + tokenValue + `", "expires_in": 3600}`
	accessTokenRespShortLived = `{"access_token": "` + tokenValue + `", "expires_in": 0}`
	scope                     = "scope"
)

type mockCredential struct {
	getTokenImpl func(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error)
}

func (mc mockCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if mc.getTokenImpl != nil {
		return mc.getTokenImpl(ctx, options)
	}
	return azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func (mc mockCredential) Do(req *policy.Request) (*http.Response, error) {
	return nil, nil
}

func defaultTestPipeline(srv policy.Transporter, scope string) Pipeline {
	retryOpts := policy.RetryOptions{
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    time.Millisecond,
	}
	b := NewBearerTokenPolicy(mockCredential{}, []string{scope}, nil)
	return NewPipeline(
		"testmodule",
		"v0.1.0",
		PipelineOptions{PerRetry: []policy.Policy{b}},
		&policy.ClientOptions{Retry: retryOpts, Transport: srv},
	)
}

func TestBearerPolicy_SuccessGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	pipeline := defaultTestPipeline(srv, scope)
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
	const expectedToken = shared.BearerTokenPrefix + tokenValue
	if token := resp.Request.Header.Get(shared.HeaderAuthorization); token != expectedToken {
		t.Fatalf("expected token '%s', got '%s'", expectedToken, token)
	}
}

func TestBearerPolicy_CredentialFailGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	expectedErr := errors.New("oops")
	failCredential := mockCredential{}
	failCredential.getTokenImpl = func(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
		return azcore.AccessToken{}, expectedErr
	}
	b := NewBearerTokenPolicy(failCredential, nil, nil)
	pipeline := newTestPipeline(&policy.ClientOptions{
		Transport: srv,
		Retry: policy.RetryOptions{
			RetryDelay: 10 * time.Millisecond,
		},
		PerRetryPolicies: []policy.Policy{b},
	})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pipeline.Do(req)
	if err != expectedErr {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestBearerTokenPolicy_TokenExpired(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespShortLived)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	pipeline := defaultTestPipeline(srv, scope)
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestBearerPolicy_GetTokenFailsNoDeadlock(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	retryOpts := policy.RetryOptions{
		// use a negative try timeout to trigger a deadline exceeded error causing GetToken() to fail
		TryTimeout:    -1 * time.Nanosecond,
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    50 * time.Millisecond,
		MaxRetries:    3,
	}
	b := NewBearerTokenPolicy(mockCredential{}, nil, nil)
	pipeline := newTestPipeline(&policy.ClientOptions{Transport: srv, Retry: retryOpts, PerRetryPolicies: []policy.Policy{b}})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pipeline.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestBearerTokenPolicy_AuthZHandler(t *testing.T) {
	challenge := "Scheme parameters..."
	srv, close := mock.NewTLSServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(mock.WithStatusCode(401), mock.WithHeader("WWW-Authenticate", challenge))
	srv.AppendResponse(mock.WithStatusCode(200))

	req, err := NewRequest(context.Background(), "GET", "https://localhost")
	require.NoError(t, err)

	type fakeHandler struct {
		policy.AuthorizationHandler
		onChallengeCalls, onReqCalls int
		onChallengeErr, onReqErr     error
	}
	handler := fakeHandler{}
	handler.OnRequest = func(r *policy.Request, f func(policy.TokenRequestOptions) error) error {
		require.Equal(t, req.Raw().URL, r.Raw().URL)
		handler.onReqCalls++
		return handler.onReqErr
	}
	handler.OnChallenge = func(r *policy.Request, res *http.Response, f func(policy.TokenRequestOptions) error) error {
		require.Equal(t, req.Raw().URL, r.Raw().URL)
		handler.onChallengeCalls++
		require.Equal(t, challenge, res.Header.Get("WWW-Authenticate"))
		return handler.onChallengeErr
	}

	b := NewBearerTokenPolicy(mockCredential{}, nil, &policy.BearerTokenOptions{AuthorizationHandler: handler.AuthorizationHandler})
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv, PerRetryPolicies: []policy.Policy{b}})

	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, handler.onChallengeCalls)
	require.Equal(t, 1, handler.onReqCalls)
	// handler functions didn't return errors, so the policy should have sent a request after calling each
	require.Equal(t, 2, srv.Requests())

	// policy should propagate the handler's errors
	handler.onReqErr = &nonRetriableError{}
	_, err = pl.Do(req)
	require.ErrorIs(t, err, handler.onReqErr)
	require.Equal(t, 1, handler.onChallengeCalls)
	require.Equal(t, 2, handler.onReqCalls)
	// OnRequest returned an error, so the policy shouldn't have sent the request
	require.Equal(t, 2, srv.Requests())

	srv.AppendResponse(mock.WithStatusCode(401), mock.WithHeader("WWW-Authenticate", challenge))
	handler.onReqErr = nil
	handler.onChallengeErr = &nonRetriableError{}
	_, err = pl.Do(req)
	require.ErrorIs(t, err, handler.onChallengeErr)
	// OnChallenge returned an error, so the policy shouldn't have sent the request again
	require.Equal(t, 3, srv.Requests())
}
