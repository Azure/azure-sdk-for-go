// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
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

func TestBearerTokenPolicy_OnChallenge(t *testing.T) {
	for _, test := range []struct {
		challenge, desc string
	}{
		{
			desc:      "no claims",
			challenge: `Bearer authorization_uri="https://login.windows.net/", error="insufficient_claims"`,
		},
		{
			desc:      "no commas",
			challenge: `Bearer authorization_uri="https://login.windows.net/" error_description="something went wrong"`,
		},
		{
			desc:      "claims with unexpected error",
			challenge: `Bearer authorization_uri="https://login.windows.net/", error="invalid_token", claims="ey=="`,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			srv, close := mock.NewTLSServer()
			defer close()
			srv.AppendResponse(mock.WithHeader(shared.HeaderWWWAuthenticate, test.challenge), mock.WithStatusCode(http.StatusUnauthorized))
			srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

			called := false
			b := NewBearerTokenPolicy(mockCredential{}, []string{scope}, &policy.BearerTokenOptions{
				AuthorizationHandler: policy.AuthorizationHandler{
					OnChallenge: func(_ *policy.Request, res *http.Response, _ func(policy.TokenRequestOptions) error) error {
						called = true
						require.EqualValues(t, test.challenge, res.Header.Get(shared.HeaderWWWAuthenticate))
						return nil
					},
				},
			})
			pipeline := newTestPipeline(&policy.ClientOptions{PerRetryPolicies: []policy.Policy{b}, Transport: srv})

			req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
			require.NoError(t, err)
			_, err = pipeline.Do(req)
			require.NoError(t, err)
			require.True(t, called, "policy should call the client's challenge handler")
		})
	}

	t.Run("errors non-retriable", func(t *testing.T) {
		srv, close := mock.NewTLSServer()
		defer close()
		srv.AppendResponse(mock.WithHeader(shared.HeaderWWWAuthenticate, `Bearer key="value"`), mock.WithStatusCode(http.StatusUnauthorized))

		expectedErr := errors.New("something went wrong")
		b := NewBearerTokenPolicy(mockCredential{}, []string{scope}, &policy.BearerTokenOptions{
			AuthorizationHandler: policy.AuthorizationHandler{
				OnChallenge: func(_ *policy.Request, _ *http.Response, _ func(policy.TokenRequestOptions) error) error {
					return expectedErr
				},
			},
		})
		pl := newTestPipeline(&policy.ClientOptions{PerRetryPolicies: []policy.Policy{b}, Transport: srv})

		req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
		require.NoError(t, err)
		_, err = pl.Do(req)
		var nre errorinfo.NonRetriable
		require.ErrorAs(t, err, &nre, "policy should ensure OnChallenge errors are NonRetriable")
		require.EqualError(t, nre, expectedErr.Error())
	})

	t.Run("CAE challenge after non-CAE challenge", func(t *testing.T) {
		cae1 := fmt.Sprintf(`Bearer error="insufficient_claims", claims=%q`, base64.StdEncoding.EncodeToString([]byte{'1'}))
		cae2 := fmt.Sprintf(`Bearer error="insufficient_claims", claims=%q`, base64.StdEncoding.EncodeToString([]byte{'2'}))
		notCAE := `Bearer authorization_uri="...", error="invalid_token"`
		for _, caeChallengeMet := range []bool{true, false} {
			srv, close := mock.NewTLSServer()
			defer close()
			srv.AppendResponse(mock.WithHeader(shared.HeaderWWWAuthenticate, notCAE), mock.WithStatusCode(http.StatusUnauthorized))
			srv.AppendResponse(mock.WithHeader(shared.HeaderWWWAuthenticate, cae1), mock.WithStatusCode(http.StatusUnauthorized))
			if caeChallengeMet {
				srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
			} else {
				srv.AppendResponse(mock.WithHeader(shared.HeaderWWWAuthenticate, cae2), mock.WithStatusCode(http.StatusUnauthorized))
			}

			onChallengeCalled := false
			tkReqs := 0
			b := NewBearerTokenPolicy(
				mockCredential{
					getTokenImpl: func(_ context.Context, actual policy.TokenRequestOptions) (exported.AccessToken, error) {
						require.Equal(t, scope, actual.Scopes[0])
						switch tkReqs {
						case 0:
						case 1, 2:
							// second and third calls should include challenge claims
							require.Equal(t, fmt.Sprint(tkReqs), actual.Claims)
						default:
							t.Fatal("unexpected token request")
						}
						tkReqs++
						return exported.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
					},
				},
				[]string{scope},
				&policy.BearerTokenOptions{
					AuthorizationHandler: policy.AuthorizationHandler{
						OnChallenge: func(_ *policy.Request, res *http.Response, _ func(policy.TokenRequestOptions) error) error {
							require.False(t, onChallengeCalled, "policy should call the client's challenge handler only once")
							onChallengeCalled = true
							actual := res.Header.Get(shared.HeaderWWWAuthenticate)
							require.Equal(t, notCAE, actual, "policy should call the client's challenge handler only for the non-CAE challenge")
							return nil
						},
					},
				})
			pl := newTestPipeline(&policy.ClientOptions{PerRetryPolicies: []policy.Policy{b}, Transport: srv})

			req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
			require.NoError(t, err)
			res, err := pl.Do(req)
			require.NoError(t, err)
			if caeChallengeMet {
				require.Equal(t, res.StatusCode, http.StatusOK)
			} else {
				require.Equal(t, res.StatusCode, http.StatusUnauthorized)
				require.Equal(t, res.Header.Get(shared.HeaderWWWAuthenticate), cae2)
			}
			require.True(t, onChallengeCalled, "policy should call the client's challenge handler for the non-CAE challenge")
		}
	})
}

func TestBearerTokenPolicy_CAEChallengeHandling(t *testing.T) {
	// requireToken is a mock.Response predicate that checks a request for the expected token
	requireToken := func(t *testing.T, want string) func(req *http.Request) bool {
		return func(r *http.Request) bool {
			_, actual, _ := strings.Cut(r.Header.Get(shared.HeaderAuthorization), " ")
			require.Equal(t, want, actual)
			return true
		}
	}
	for _, test := range []struct {
		challenge, desc, expectedClaims string
		err                             error
	}{
		{
			desc: "no challenge",
		},
		{
			desc:      "invalid claims",
			challenge: `Bearer claims="not base64", error="insufficient_claims"`,
			err:       (*exported.ResponseError)(nil),
		},
		{
			desc:           "standard",
			challenge:      `Bearer realm="", authorization_uri="http://localhost", error="insufficient_claims", claims="eyJhY2Nlc3NfdG9rZW4iOnsibmJmIjp7ImVzc2VudGlhbCI6dHJ1ZSwidmFsdWUiOiIxNzI2MDc3NTk1In0sInhtc19jYWVlcnJvciI6eyJ2YWx1ZSI6IjEwMDEyIn19fQ=="`,
			expectedClaims: `{"access_token":{"nbf":{"essential":true,"value":"1726077595"},"xms_caeerror":{"value":"10012"}}}`,
		},
		{
			desc:           "multiple challenges",
			challenge:      `PoP realm="", authorization_uri="http://localhost", client_id="...", nonce="ey==", Bearer realm="", error="insufficient_claims", authorization_uri="http://localhost", client_id="...", error_description="Continuous access evaluation resulted in challenge with result: InteractionRequired and code: TokenIssuedBeforeRevocationTimestamp", claims="eyJhY2Nlc3NfdG9rZW4iOnsibmJmIjp7ImVzc2VudGlhbCI6dHJ1ZSwgInZhbHVlIjoiMTcyNjI1ODEyMiJ9fX0="`,
			expectedClaims: `{"access_token":{"nbf":{"essential":true, "value":"1726258122"}}}`,
		},
		{
			desc:           "CAE+unparseable challenge",
			challenge:      `Foo bar=can't parse this, error=my bad, Bearer claims="ey==", error="insufficient_claims"`,
			expectedClaims: "{",
		},
	} {
		for _, customOnRequest := range []bool{false, true} {
			expectedTRO := policy.TokenRequestOptions{
				Claims:    test.expectedClaims,
				EnableCAE: true,
				Scopes:    []string{scope},
			}
			var (
				name      string
				onRequest func(*policy.Request, func(policy.TokenRequestOptions) error) error
			)
			if customOnRequest {
				name = "/custom OnRequest"
				expectedTRO.Scopes = []string{"scope set by OnRequest"}
				expectedTRO.TenantID = "tenant set by OnRequest"
				onRequest = func(_ *policy.Request, authNZ func(policy.TokenRequestOptions) error) error {
					tro := expectedTRO
					// zero fields the policy should set so the test fails when it doesn't set them
					tro.Claims = ""
					tro.EnableCAE = false
					return authNZ(tro)
				}
			}
			t.Run(test.desc+name, func(t *testing.T) {
				challengedToken := "needs more claims"
				tokenWithClaims := "all the claims"

				srv, close := mock.NewTLSServer()
				defer close()
				srv.AppendResponse(
					mock.WithHeader(shared.HeaderWWWAuthenticate, test.challenge),
					mock.WithPredicate(requireToken(t, challengedToken)),
					mock.WithStatusCode(http.StatusUnauthorized),
				)
				srv.AppendResponse() // when a response's predicate returns true, srv pops the following one
				srv.AppendResponse(mock.WithPredicate(requireToken(t, tokenWithClaims)))
				srv.AppendResponse()
				srv.AppendResponse(mock.WithPredicate(requireToken(t, tokenWithClaims)))
				srv.AppendResponse()

				tkReqs := 0
				cred := mockCredential{
					getTokenImpl: func(_ context.Context, actual policy.TokenRequestOptions) (exported.AccessToken, error) {
						require.True(t, actual.EnableCAE, "policy should always request CAE-enabled tokens")
						tkReqs += 1
						tk := challengedToken
						switch tkReqs {
						case 1:
							require.Empty(t, actual.Claims, "policy should specify claims only when handling a CAE challenge")
						case 2:
							tk = tokenWithClaims
							require.Equal(t, expectedTRO, actual)
						default:
							t.Fatal("unexpected token request")
						}
						return exported.AccessToken{Token: tk, ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
					},
				}
				var scopes []string
				if !customOnRequest {
					cp := make([]string, len(expectedTRO.Scopes))
					copy(cp, expectedTRO.Scopes)
					scopes = cp
				}
				b := NewBearerTokenPolicy(cred, scopes, &policy.BearerTokenOptions{
					AuthorizationHandler: policy.AuthorizationHandler{
						OnChallenge: func(*policy.Request, *http.Response, func(policy.TokenRequestOptions) error) error {
							t.Fatal("policy shouldn't call a client's challenge handler")
							return nil
						},
						OnRequest: onRequest,
					},
				})
				pipeline := newTestPipeline(&policy.ClientOptions{PerRetryPolicies: []policy.Policy{b}, Transport: srv})
				req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
				require.NoError(t, err)
				_, err = pipeline.Do(req)
				if test.err == nil {
					require.NoError(t, err)
					// send another request to verify the policy cached the token it acquired to satisfy the challenge
					_, err = pipeline.Do(req)
					require.NoError(t, err)
				} else {
					require.ErrorAs(t, err, &test.err)
				}
				if test.expectedClaims != "" {
					require.Equal(t, 2, tkReqs, "policy should request a new token upon receiving the challenge")
				}
			})
		}
	}

	t.Run("consecutive challenges", func(t *testing.T) {
		srv, close := mock.NewTLSServer()
		defer close()
		srv.SetResponse(
			mock.WithHeader(shared.HeaderWWWAuthenticate, `Bearer error="insufficient_claims", claims="ey=="`),
			mock.WithStatusCode(http.StatusUnauthorized),
		)

		tkReqs := 0
		cred := mockCredential{
			getTokenImpl: func(_ context.Context, actual policy.TokenRequestOptions) (exported.AccessToken, error) {
				tkReqs++
				return exported.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
			},
		}
		b := NewBearerTokenPolicy(cred, []string{scope}, &policy.BearerTokenOptions{
			AuthorizationHandler: policy.AuthorizationHandler{
				OnChallenge: func(*policy.Request, *http.Response, func(policy.TokenRequestOptions) error) error {
					t.Fatal("policy shouldn't call a client's challenge handler")
					return nil
				},
			},
		})
		pipeline := newTestPipeline(&policy.ClientOptions{PerRetryPolicies: []policy.Policy{b}, Transport: srv})

		req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
		require.NoError(t, err)
		_, err = pipeline.Do(req)
		require.NoError(t, err)
		require.Equal(t, 2, tkReqs, "policy shouldn't handle a second CAE challenge for the same request")
		require.Equal(t, 2, srv.Requests(), "policy shouldn't handle a second CAE challenge for the same request")
	})

	t.Run("errors non-retriable", func(t *testing.T) {
		srv, close := mock.NewTLSServer()
		defer close()
		srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
		srv.AppendResponse(
			mock.WithHeader(shared.HeaderWWWAuthenticate, `Bearer error="insufficient_claims", claims="ey=="`),
			mock.WithStatusCode(http.StatusUnauthorized),
		)

		called := false
		expectedErr := errors.New("something went wrong")
		cred := mockCredential{
			getTokenImpl: func(context.Context, policy.TokenRequestOptions) (exported.AccessToken, error) {
				if called {
					return exported.AccessToken{}, expectedErr
				}
				called = true
				return exported.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
			},
		}
		counter := &countingPolicy{}
		btp := NewBearerTokenPolicy(cred, []string{scope}, nil)
		pl := newTestPipeline(&policy.ClientOptions{PerRetryPolicies: []policy.Policy{counter, btp}, Transport: srv})

		req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
		require.NoError(t, err)
		_, err = pl.Do(req)
		require.NoError(t, err)

		req, err = NewRequest(context.Background(), http.MethodGet, srv.URL())
		require.NoError(t, err)
		_, err = pl.Do(req)
		require.EqualError(t, err, expectedErr.Error())
		require.ErrorAs(t, err, new(errorinfo.NonRetriable))
		// this is the crucial assertion; the retry policy would have retried the request
		// if BearerTokenPolicy didn't make the credential's error NonRetriable
		require.Equal(t, 2, counter.count, "BearerTokenPolicy should make the authentication error NonRetriable")
	})
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

func TestBearerTokenPolicy_RewindsBeforeRetry(t *testing.T) {
	const expected = "expected"
	for _, test := range []struct {
		challenge, desc string
		onChallenge     bool
	}{
		{
			desc:      "CAE challenge",
			challenge: `Bearer error="insufficient_claims", claims="ey=="`,
		},
		{
			desc:        "non-CAE challenge",
			challenge:   `Bearer authorization_uri="https://login.windows.net/", error="invalid_token"`,
			onChallenge: true,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			read := func(r *http.Request) bool {
				actual, err := io.ReadAll(r.Body)
				require.NoError(t, err, "request should have body content")
				require.EqualValues(t, expected, actual)
				return true
			}
			srv, close := mock.NewTLSServer()
			defer close()
			srv.AppendResponse(
				mock.WithHeader(shared.HeaderWWWAuthenticate, test.challenge),
				mock.WithPredicate(read),
				mock.WithStatusCode(http.StatusUnauthorized),
			)
			srv.AppendResponse()
			srv.AppendResponse(mock.WithPredicate(read))
			srv.AppendResponse()

			called := false
			o := &policy.BearerTokenOptions{}
			if test.onChallenge {
				o.AuthorizationHandler.OnChallenge = func(*policy.Request, *http.Response, func(policy.TokenRequestOptions) error) error {
					called = true
					return nil
				}
			}
			b := NewBearerTokenPolicy(mockCredential{}, []string{scope}, o)
			pl := newTestPipeline(&policy.ClientOptions{PerRetryPolicies: []policy.Policy{b}, Transport: srv})
			req, err := NewRequest(context.Background(), http.MethodPost, srv.URL())
			require.NoError(t, err)
			require.NoError(t, req.SetBody(streaming.NopCloser(strings.NewReader(expected)), "text/plain"))

			_, err = pl.Do(req)
			require.NoError(t, err)
			require.Equal(t, test.onChallenge, called, "policy should call OnChallenge when set")
		})
	}
}

func TestBearerTokenPolicy_ShouldRefresh(t *testing.T) {
	now := time.Now()
	for _, test := range []struct {
		desc     string
		expected bool
		tk       exported.AccessToken
	}{
		{
			desc: "distant RefreshOn/distant ExpiresOn",
			tk: exported.AccessToken{
				ExpiresOn: now.Add(2 * time.Hour).UTC(),
				RefreshOn: now.Add(time.Hour).UTC(),
			},
		},
		{
			desc: "zero RefreshOn/distant ExpiresOn",
			tk: exported.AccessToken{
				ExpiresOn: now.Add(time.Hour).UTC(),
			},
		},
		{
			desc: "zero RefreshOn/imminent ExpiresOn",
			tk: exported.AccessToken{
				ExpiresOn: now.Add(4 * time.Minute).UTC(),
			},
			expected: true,
		},
		{
			desc: "zero RefreshOn/past ExpiresOn",
			tk: exported.AccessToken{
				ExpiresOn: now.Add(-time.Minute).UTC(),
			},
			expected: true,
		},
		{
			desc: "past RefreshOn",
			tk: exported.AccessToken{
				ExpiresOn: now.Add(time.Hour).UTC(),
				RefreshOn: now.Add(-time.Minute).UTC(),
			},
			expected: true,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			actual := shouldRefresh(test.tk, acquiringResourceState{})
			require.Equal(t, test.expected, actual)
		})
	}
	t.Run("called", func(t *testing.T) {
		expected := exported.AccessToken{Token: "***", ExpiresOn: now.Add(time.Hour).UTC(), RefreshOn: now.Add(-time.Minute).UTC()}
		called := false
		before := shouldRefresh
		defer func() { shouldRefresh = before }()
		shouldRefresh = func(tk exported.AccessToken, state acquiringResourceState) bool {
			require.Equal(t, expected, tk)
			called = true
			return false
		}
		c := mockCredential{
			getTokenImpl: func(context.Context, policy.TokenRequestOptions) (exported.AccessToken, error) {
				return expected, nil
			},
		}
		p := NewBearerTokenPolicy(c, []string{scope}, nil)
		srv, close := mock.NewTLSServer()
		defer close()
		pl := newTestPipeline(&policy.ClientOptions{PerRetryPolicies: []policy.Policy{p}, Transport: srv})
		for range 2 {
			srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
			req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
			require.NoError(t, err)
			res, err := pl.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
		}
		require.True(t, called, "temporal.Resource should have called shouldRefresh")
	})
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
