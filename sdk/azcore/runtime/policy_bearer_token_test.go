// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"

	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
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
	getTokenImpl func(ctx context.Context, options policy.TokenRequestOptions) (exported.AccessToken, error)
}

func (mc mockCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (exported.AccessToken, error) {
	if mc.getTokenImpl != nil {
		return mc.getTokenImpl(ctx, options)
	}
	return exported.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
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
	failCredential.getTokenImpl = func(ctx context.Context, options policy.TokenRequestOptions) (exported.AccessToken, error) {
		return exported.AccessToken{}, expectedErr
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
	require.NoError(t, err)
	resp, err := pipeline.Do(req)
	require.EqualError(t, err, expectedErr.Error())
	require.Nil(t, resp)
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
	srv.AppendResponse(mock.WithStatusCode(401), mock.WithHeader(shared.HeaderWWWAuthenticate, challenge))
	srv.AppendResponse(mock.WithStatusCode(200))

	req, err := NewRequest(context.Background(), "GET", "https://localhost")
	require.NoError(t, err)

	handler := struct {
		policy.AuthorizationHandler
		onChallengeCalls, onReqCalls int
	}{}
	handler.OnRequest = func(r *policy.Request, f func(policy.TokenRequestOptions) error) error {
		require.Equal(t, req.Raw().URL, r.Raw().URL)
		handler.onReqCalls++
		return nil
	}
	handler.OnChallenge = func(r *policy.Request, res *http.Response, f func(policy.TokenRequestOptions) error) error {
		require.Equal(t, req.Raw().URL, r.Raw().URL)
		handler.onChallengeCalls++
		require.Equal(t, challenge, res.Header.Get(shared.HeaderWWWAuthenticate))
		return nil
	}

	b := NewBearerTokenPolicy(mockCredential{}, nil, &policy.BearerTokenOptions{AuthorizationHandler: handler.AuthorizationHandler})
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv, PerRetryPolicies: []policy.Policy{b}})

	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, handler.onChallengeCalls)
	require.Equal(t, 1, handler.onReqCalls)
	// handler functions didn't return errors, so the policy should have sent a request after calling each
	require.Equal(t, 2, srv.Requests())

	t.Run("SupportsCAE", func(t *testing.T) {
		for _, supportsCAE := range []bool{true, false} {
			srv, close := mock.NewTLSServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(mock.WithStatusCode(200))

			calls := 0
			cred := mockCredential{
				getTokenImpl: func(_ context.Context, actual policy.TokenRequestOptions) (exported.AccessToken, error) {
					calls += 1
					require.Equal(t, supportsCAE, actual.EnableCAE, "policy should request CAE tokens only when AuthorizationHandler.SupportsCAE is true")
					return exported.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour)}, nil
				},
			}
			o := policy.BearerTokenOptions{AuthorizationHandler: policy.AuthorizationHandler{
				OnChallenge: func(*policy.Request, *http.Response, func(policy.TokenRequestOptions) error) error {
					return nil
				},
				SupportsCAE: supportsCAE,
			}}
			b = NewBearerTokenPolicy(cred, nil, &o)
			pl = newTestPipeline(&policy.ClientOptions{PerRetryPolicies: []policy.Policy{b}, Transport: srv})
			_, err = pl.Do(req)
			require.NoError(t, err)
			require.Equal(t, 1, calls, "policy should have called GetToken once")
		}
	})
}

func TestBearerTokenPolicy_AuthZHandlerErrors(t *testing.T) {
	srv, close := mock.NewTLSServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.SetResponse(mock.WithStatusCode(401), mock.WithHeader(shared.HeaderWWWAuthenticate, "..."))

	req, err := NewRequest(context.Background(), "GET", "https://localhost")
	require.NoError(t, err)

	handler := struct {
		policy.AuthorizationHandler
		onChallengeErr, onReqErr error
	}{}
	handler.OnRequest = func(r *policy.Request, f func(policy.TokenRequestOptions) error) error {
		return handler.onReqErr
	}
	handler.OnChallenge = func(r *policy.Request, res *http.Response, f func(policy.TokenRequestOptions) error) error {
		return handler.onChallengeErr
	}

	b := NewBearerTokenPolicy(mockCredential{}, nil, &policy.BearerTokenOptions{AuthorizationHandler: handler.AuthorizationHandler})
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv, PerRetryPolicies: []policy.Policy{b}})

	// the policy should propagate the handler's errors, wrapping them to make them nonretriable, if necessary
	fatalErr := errors.New("something went wrong")
	var nre errorinfo.NonRetriable
	for i, e := range []error{fatalErr, errorinfo.NonRetriableError(fatalErr)} {
		handler.onReqErr = e
		_, err = pl.Do(req)
		require.ErrorAs(t, err, &nre)
		require.EqualError(t, nre, fatalErr.Error())
		// the policy shouldn't have sent a request, because OnRequest returned an error
		require.Equal(t, i, srv.Requests())

		handler.onReqErr = nil
		handler.onChallengeErr = e
		_, err = pl.Do(req)
		require.ErrorAs(t, err, &nre)
		require.EqualError(t, nre, fatalErr.Error())
		handler.onChallengeErr = nil
		// the policy should have sent one request, because OnRequest returned nil but OnChallenge returned an error
		require.Equal(t, i+1, srv.Requests())
	}
}

func TestBearerTokenPolicyChallengeHandling(t *testing.T) {
	for _, test := range []struct {
		challenge, desc, expectedClaims string
		err                             error
	}{
		{
			desc: "no challenge",
		},
		{
			desc:      "no claims",
			challenge: `Bearer authorization_uri="https://login.windows.net/", error="invalid_token", error_description="The authentication failed because of missing 'Authorization' header."`,
			err:       (*exported.ResponseError)(nil),
		},
		{
			desc:      "parsing error",
			challenge: `Bearer claims="not base64"`,
			err:       (*exported.ResponseError)(nil),
		},

		// CAE claims challenges. Position of the "claims" parameter within the challenge shouldn't affect parsing.
		{
			desc:           "insufficient claims",
			challenge:      `Bearer realm="", authorization_uri="https://login.microsoftonline.com/common/oauth2/authorize", client_id="00000003-0000-0000-c000-000000000000", error="insufficient_claims", claims="eyJhY2Nlc3NfdG9rZW4iOiB7ImZvbyI6ICJiYXIifX0="`,
			expectedClaims: `{"access_token": {"foo": "bar"}}`,
		},
		{
			desc:           "insufficient claims",
			challenge:      `Bearer claims="eyJhY2Nlc3NfdG9rZW4iOiB7ImZvbyI6ICJiYXIifX0=", realm="", authorization_uri="https://login.microsoftonline.com/common/oauth2/authorize", client_id="00000003-0000-0000-c000-000000000000", error="insufficient_claims"`,
			expectedClaims: `{"access_token": {"foo": "bar"}}`,
		},
		{
			desc:           "sessions revoked",
			challenge:      `Bearer authorization_uri="https://login.windows.net/", error="invalid_token", error_description="User session has been revoked", claims="eyJhY2Nlc3NfdG9rZW4iOnsibmJmIjp7ImVzc2VudGlhbCI6dHJ1ZSwgInZhbHVlIjoiMTYwMzc0MjgwMCJ9fX0="`,
			expectedClaims: `{"access_token":{"nbf":{"essential":true, "value":"1603742800"}}}`,
		},
		{
			desc:           "sessions revoked",
			challenge:      `Bearer authorization_uri="https://login.windows.net/", claims="eyJhY2Nlc3NfdG9rZW4iOnsibmJmIjp7ImVzc2VudGlhbCI6dHJ1ZSwgInZhbHVlIjoiMTYwMzc0MjgwMCJ9fX0=", error="invalid_token", error_description="User session has been revoked"`,
			expectedClaims: `{"access_token":{"nbf":{"essential":true, "value":"1603742800"}}}`,
		},
		{
			desc:           "IP policy",
			challenge:      `Bearer authorization_uri="https://login.windows.net/", error="invalid_token", error_description="Tenant IP Policy validate failed.", claims="eyJhY2Nlc3NfdG9rZW4iOnsibmJmIjp7ImVzc2VudGlhbCI6dHJ1ZSwidmFsdWUiOiIxNjEwNTYzMDA2In0sInhtc19ycF9pcGFkZHIiOnsidmFsdWUiOiIxLjIuMy40In19fQ"`,
			expectedClaims: `{"access_token":{"nbf":{"essential":true,"value":"1610563006"},"xms_rp_ipaddr":{"value":"1.2.3.4"}}}`,
		},
		{
			desc:           "IP policy",
			challenge:      `Bearer authorization_uri="https://login.windows.net/", error="invalid_token", claims="eyJhY2Nlc3NfdG9rZW4iOnsibmJmIjp7ImVzc2VudGlhbCI6dHJ1ZSwidmFsdWUiOiIxNjEwNTYzMDA2In0sInhtc19ycF9pcGFkZHIiOnsidmFsdWUiOiIxLjIuMy40In19fQ", error_description="Tenant IP Policy validate failed."`,
			expectedClaims: `{"access_token":{"nbf":{"essential":true,"value":"1610563006"},"xms_rp_ipaddr":{"value":"1.2.3.4"}}}`,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			srv, close := mock.NewTLSServer()
			defer close()
			srv.SetResponse(mock.WithHeader(shared.HeaderWWWAuthenticate, test.challenge), mock.WithStatusCode(http.StatusUnauthorized))
			tkReqs := 0
			cred := mockCredential{
				getTokenImpl: func(_ context.Context, actual policy.TokenRequestOptions) (exported.AccessToken, error) {
					require.True(t, actual.EnableCAE)
					tkReqs += 1
					switch tkReqs {
					case 1:
						require.Empty(t, actual.Claims)
					case 2:
						require.Equal(t, test.expectedClaims, actual.Claims)
					default:
						t.Fatalf("unexpected token request")
					}
					return exported.AccessToken{Token: "...", ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
				},
			}
			b := NewBearerTokenPolicy(cred, []string{scope}, nil)
			pipeline := newTestPipeline(&policy.ClientOptions{PerRetryPolicies: []policy.Policy{b}, Transport: srv})
			req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
			require.NoError(t, err)
			_, err = pipeline.Do(req)
			if test.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, &test.err)
			}
			if test.expectedClaims != "" {
				require.Equal(t, 2, tkReqs, "policy should have requested a new token upon receiving the challenge")
			}
		})
	}
}

func TestBearerTokenPolicy_RequiresHTTPS(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	b := NewBearerTokenPolicy(mockCredential{}, nil, nil)
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv, PerRetryPolicies: []policy.Policy{b}})
	req, err := NewRequest(context.Background(), "GET", srv.URL())
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.Error(t, err)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
}

func TestCheckHTTPSForAuth(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodGet, "http://contoso.com")
	require.NoError(t, err)
	require.Error(t, checkHTTPSForAuth(req, false))
	req, err = NewRequest(context.Background(), http.MethodGet, "https://contoso.com")
	require.NoError(t, err)
	require.NoError(t, checkHTTPSForAuth(req, false))
	req, err = NewRequest(context.Background(), http.MethodGet, "http://contoso.com")
	require.NoError(t, err)
	require.NoError(t, checkHTTPSForAuth(req, true))
	req, err = NewRequest(context.Background(), http.MethodGet, "https://contoso.com")
	require.NoError(t, err)
	require.NoError(t, checkHTTPSForAuth(req, true))
}

func TestBearerTokenPolicy_NilCredential(t *testing.T) {
	policy := NewBearerTokenPolicy(nil, nil, nil)
	pl := exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		require.Zero(t, req.Header.Get(shared.HeaderAuthorization))
		return &http.Response{}, nil
	}), policy)
	req, err := NewRequest(context.Background(), "GET", "http://contoso.com")
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
}
