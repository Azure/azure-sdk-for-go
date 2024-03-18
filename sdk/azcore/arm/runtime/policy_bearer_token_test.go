// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"errors"
	"strings"

	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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
	getTokenImpl func(ctx context.Context, options azpolicy.TokenRequestOptions) (azcore.AccessToken, error)
}

func (mc mockCredential) GetToken(ctx context.Context, options azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
	if !options.EnableCAE {
		return azcore.AccessToken{}, errors.New("ARM clients should set EnableCAE to true")
	}
	if mc.getTokenImpl != nil {
		return mc.getTokenImpl(ctx, options)
	}
	return azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func (mc mockCredential) Do(req *azpolicy.Request) (*http.Response, error) {
	return nil, nil
}

func newTestPipeline(opts *azpolicy.ClientOptions) runtime.Pipeline {
	return runtime.NewPipeline("testmodule", "v0.1.0", runtime.PipelineOptions{}, opts)
}

func defaultTestPipeline(srv azpolicy.Transporter) (runtime.Pipeline, error) {
	retryOpts := azpolicy.RetryOptions{
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    time.Millisecond,
	}
	return NewPipeline(
		"testmodule",
		"v0.1.0",
		mockCredential{},
		runtime.PipelineOptions{},
		&armpolicy.ClientOptions{
			ClientOptions: azpolicy.ClientOptions{
				Retry:     retryOpts,
				Transport: srv,
			},
		})
}

func TestBearerPolicy_SuccessGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	pipeline, err := defaultTestPipeline(srv)
	if err != nil {
		t.Fatal(err)
	}
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
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
	expectedErr := "oops"
	failCredential := mockCredential{}
	failCredential.getTokenImpl = func(ctx context.Context, options azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
		return azcore.AccessToken{}, errors.New(expectedErr)
	}
	b := NewBearerTokenPolicy(failCredential, nil)
	pipeline := newTestPipeline(&azpolicy.ClientOptions{
		Transport: srv,
		Retry: azpolicy.RetryOptions{
			RetryDelay: 10 * time.Millisecond,
		},
		PerRetryPolicies: []azpolicy.Policy{b},
	})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	resp, err := pipeline.Do(req)
	require.EqualError(t, err, expectedErr)
	require.Nil(t, resp)
	require.Implements(t, (*errorinfo.NonRetriable)(nil), err)
}

func TestBearerTokenPolicy_TokenExpired(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespShortLived)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	pipeline, err := defaultTestPipeline(srv)
	if err != nil {
		t.Fatal(err)
	}
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
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
	retryOpts := azpolicy.RetryOptions{
		// use a negative try timeout to trigger a deadline exceeded error causing GetToken() to fail
		TryTimeout:    -1 * time.Nanosecond,
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    50 * time.Millisecond,
		MaxRetries:    3,
	}
	b := NewBearerTokenPolicy(mockCredential{}, nil)
	pipeline := newTestPipeline(&azpolicy.ClientOptions{Transport: srv, Retry: retryOpts, PerRetryPolicies: []azpolicy.Policy{b}})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
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

func TestAuxiliaryTenants(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	primary := "primary"
	auxTenants := []string{"aux1", "aux2", "aux3"}
	expectCache := false
	b := NewBearerTokenPolicy(
		mockCredential{
			// getTokenImpl returns a token whose value equals the requested tenant so the test can validate how the policy handles tenants
			// i.e., primary tenant token goes in Authorization header and aux tenant tokens go in x-ms-authorization-auxiliary
			getTokenImpl: func(ctx context.Context, options azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
				require.False(t, expectCache, "client should have used a cached token instead of requesting another")
				tenant := primary
				if options.TenantID != "" {
					tenant = options.TenantID
				}
				return azcore.AccessToken{Token: tenant, ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
			},
		},
		&armpolicy.BearerTokenOptions{AuxiliaryTenants: auxTenants, Scopes: []string{scope}},
	)
	pipeline := newTestPipeline(&azpolicy.ClientOptions{Transport: srv, PerRetryPolicies: []azpolicy.Policy{b}})
	expected := strings.Split(shared.BearerTokenPrefix+strings.Join(auxTenants, ","+shared.BearerTokenPrefix), ",")
	for i := 0; i < 3; i++ {
		if i == 1 {
			// policy should have a cached token after the first iteration
			expectCache = true
		}
		req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
		require.NoError(t, err)
		resp, err := pipeline.Do(req)
		require.NoError(t, err)
		require.Equal(t, shared.BearerTokenPrefix+primary, resp.Request.Header.Get(shared.HeaderAuthorization), "Authorization header must contain primary tenant token")
		actual := strings.Split(resp.Request.Header.Get(headerAuxiliaryAuthorization), ", ")
		// auxiliary tokens may appear in arbitrary order
		require.ElementsMatch(t, expected, actual)
	}
}

func TestBearerTokenPolicyChallengeParsing(t *testing.T) {
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
			err:       (*azcore.ResponseError)(nil),
		},
		{
			desc:      "parsing error",
			challenge: `Bearer claims="invalid"`,
			// the specific error type isn't important but it must be nonretriable
			err: (errorinfo.NonRetriable)(nil),
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
			calls := 0
			cred := mockCredential{
				getTokenImpl: func(ctx context.Context, actual azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
					calls += 1
					if calls == 2 && test.expectedClaims != "" {
						require.Equal(t, test.expectedClaims, actual.Claims)
					}
					return azcore.AccessToken{Token: "...", ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
				},
			}
			b := NewBearerTokenPolicy(cred, &armpolicy.BearerTokenOptions{Scopes: []string{scope}})
			pipeline := newTestPipeline(&azpolicy.ClientOptions{Transport: srv, PerRetryPolicies: []azpolicy.Policy{b}})
			req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
			require.NoError(t, err)
			_, err = pipeline.Do(req)
			if test.err != nil {
				require.ErrorAs(t, err, &test.err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedClaims != "" {
				require.Equal(t, 2, calls, "policy should have requested a new token upon receiving the challenge")
			}
		})
	}
}

func TestBearerTokenPolicyRequiresHTTPS(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	b := NewBearerTokenPolicy(mockCredential{}, nil)
	pl := newTestPipeline(&azpolicy.ClientOptions{Transport: srv, PerRetryPolicies: []azpolicy.Policy{b}})
	req, err := runtime.NewRequest(context.Background(), "GET", srv.URL())
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.Error(t, err)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
}

func TestBearerTokenPolicyAllowHTTP(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	b := NewBearerTokenPolicy(mockCredential{}, &armpolicy.BearerTokenOptions{
		InsecureAllowCredentialWithHTTP: true,
	})
	pl := newTestPipeline(&azpolicy.ClientOptions{Transport: srv, PerRetryPolicies: []azpolicy.Policy{b}})
	req, err := runtime.NewRequest(context.Background(), "GET", srv.URL())
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
}
