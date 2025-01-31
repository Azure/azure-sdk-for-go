//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

const (
	challengedToken = "needs more claims"
	claimsToken     = "all the claims"
	kvChallenge1    = `Bearer authorization="https://login.microsoftonline.com/tenant", resource="https://vault.azure.net"`
	kvChallenge2    = `Bearer authorization="https://login.microsoftonline.com/tenant2", resource="https://vault.azure.net"`
	caeChallenge1   = `Bearer realm="", authorization_uri="https://login.microsoftonline.com/common/oauth2/authorize", error="insufficient_claims", claims="dGVzdGluZzE="`
	caeChallenge2   = `Bearer realm="", authorization_uri="https://login.microsoftonline.com/common/oauth2/authorize", error="insufficient_claims", claims="dGVzdGluZzI="`
)

// requireToken is a mock.Response predicate that checks a request for the expected token
var requireToken = func(t *testing.T, want string) func(req *http.Request) bool {
	return func(r *http.Request) bool {
		_, actual, _ := strings.Cut(r.Header.Get("Authorization"), " ")
		require.Equal(t, want, actual)
		return true
	}
}

type credentialFunc func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error)

func (cf credentialFunc) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return cf(ctx, options)
}

func TestChallengePolicy(t *testing.T) {
	accessToken := "***"
	resource := "https://vault.azure.net"
	scope := "https://vault.azure.net/.default"
	challengeResource := `Bearer authorization="https://login.microsoftonline.com/{tenant}", resource="{resource}"`
	challengeScope := `Bearer authorization="https://login.microsoftonline.com/{tenant}", scope="{resource}"`

	for _, test := range []struct {
		expectedScope, format, resource string
		disableVerify, err              bool
	}{
		// happy path: resource matches requested vault's host (vault.azure.net)
		{format: challengeResource, resource: resource, expectedScope: scope},
		{format: challengeResource, resource: resource, disableVerify: true, expectedScope: scope},
		{format: challengeScope, resource: scope, expectedScope: scope},
		{format: challengeScope, resource: scope, disableVerify: true, expectedScope: scope},
		// the policy should prefer scope to resource when a challenge specifies both
		{format: fmt.Sprintf(`%s scope="%s"`, challengeResource, scope), resource: resource, expectedScope: scope},
		{format: challengeScope + ` resource="ignore me"`, resource: scope, expectedScope: scope},

		// error cases: resource/scope doesn't match the requested vault's host (vault.azure.net)
		{format: challengeResource, resource: "https://vault.azure.cn", err: true},
		{format: challengeResource, resource: "https://myvault.azure.net", err: true},
		{format: challengeScope, resource: "https://vault.azure.cn/.default", err: true},
		{format: challengeScope, resource: "https://myvault.azure.net/.default", err: true},

		// the policy shouldn't return errors for the above cases when verification is disabled
		{format: challengeResource, resource: "https://vault.azure.cn", disableVerify: true, expectedScope: "https://vault.azure.cn/.default"},
		{format: challengeResource, resource: "https://myvault.azure.net", disableVerify: true, expectedScope: "https://myvault.azure.net/.default"},
		{format: challengeScope, resource: "https://vault.azure.cn/.default", disableVerify: true, expectedScope: "https://vault.azure.cn/.default"},
		{format: challengeScope, resource: "https://myvault.azure.net/.default", disableVerify: true, expectedScope: "https://myvault.azure.net/.default"},
	} {
		t.Run("", func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(
				mock.WithHeader("WWW-Authenticate", strings.ReplaceAll(test.format, "{resource}", test.resource)),
				mock.WithStatusCode(401),
			)
			srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
				if authz := r.Header.Values("Authorization"); len(authz) != 1 || authz[0] != "Bearer "+accessToken {
					t.Errorf(`unexpected Authorization "%s"`, authz)
				}
				return true
			}))
			srv.AppendResponse()
			authenticated := false
			cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
				authenticated = true
				require.Equal(t, []string{test.expectedScope}, tro.Scopes)
				require.Equal(t, "{tenant}", tro.TenantID)
				return azcore.AccessToken{Token: accessToken, ExpiresOn: time.Now().Add(time.Hour)}, nil
			})
			p := NewKeyVaultChallengePolicy(cred, &KeyVaultChallengePolicyOptions{DisableChallengeResourceVerification: test.disableVerify})
			pl := runtime.NewPipeline("", "",
				runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
				&policy.ClientOptions{Transport: srv},
			)
			req, err := runtime.NewRequest(context.Background(), "GET", "https://42.vault.azure.net")
			require.NoError(t, err)
			_, err = pl.Do(req)
			if test.err {
				expected := fmt.Sprintf(challengeMatchError, test.resource)
				require.EqualError(t, err, expected)
				var nre errorinfo.NonRetriable
				require.ErrorAs(t, err, &nre)
			} else {
				require.True(t, authenticated, "policy should have authenticated")
			}
		})
	}
}

func TestChallengePolicy_Tenant(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", kvChallenge1),
		mock.WithStatusCode(401),
	)
	srv.AppendResponse()
	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", kvChallenge2),
		mock.WithStatusCode(401),
	)
	srv.AppendResponse()

	tkReqs := 0
	cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		require.True(t, tro.EnableCAE)
		tkReqs += 1
		switch tkReqs {
		case 1:
			require.Equal(t, "tenant", tro.TenantID)
		case 2:
			require.Equal(t, "tenant2", tro.TenantID)
		default:
			t.Fatal("unexpected token request")
		}
		return azcore.AccessToken{Token: "token", ExpiresOn: time.Now().Add(time.Hour)}, nil
	})
	p := NewKeyVaultChallengePolicy(cred, nil)
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)
	req, err := runtime.NewRequest(context.Background(), "GET", "https://42.vault.azure.net")
	require.NoError(t, err)
	res, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, tkReqs, 1)

	req, err = runtime.NewRequest(context.Background(), "GET", "https://42.vault.azure.net")
	require.NoError(t, err)
	res, err = pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, tkReqs, 2)
}

func TestChallengePolicy_CAE(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", kvChallenge1),
		mock.WithStatusCode(401),
		mock.WithPredicate(requireToken(t, "")),
	)
	srv.AppendResponse() // when a response's predicate returns true, srv pops the following one
	srv.AppendResponse()

	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", caeChallenge1),
		mock.WithStatusCode(401),
		mock.WithPredicate(requireToken(t, challengedToken)),
	)
	srv.AppendResponse() // when a response's predicate returns true, srv pops the following one
	srv.AppendResponse(
		mock.WithPredicate(requireToken(t, claimsToken)),
	)
	srv.AppendResponse() // when a response's predicate returns true, srv pops the following one

	tkReqs := 0
	cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		require.True(t, tro.EnableCAE)
		tkReqs += 1
		tk := challengedToken
		switch tkReqs {
		case 1:
			require.Empty(t, tro.Claims)
		case 2:
			tk = claimsToken
			require.Equal(t, "testing1", tro.Claims)
		default:
			t.Fatal("unexpected token request")
		}
		return azcore.AccessToken{Token: tk, ExpiresOn: time.Now().Add(time.Hour)}, nil
	})
	p := NewKeyVaultChallengePolicy(cred, nil)
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)

	// req 1 kv then regular
	req, err := runtime.NewRequest(context.Background(), "POST", "https://42.vault.azure.net")
	require.NoError(t, err)
	err = req.SetBody(streaming.NopCloser(bytes.NewReader([]byte("test"))), "text/plain")
	require.NoError(t, err)
	res, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, 1, tkReqs)

	// req 2 cae
	req, err = runtime.NewRequest(context.Background(), "POST", "https://42.vault.azure.net")
	require.NoError(t, err)
	err = req.SetBody(streaming.NopCloser(bytes.NewReader([]byte("test2"))), "text/plain")
	require.NoError(t, err)
	res, err = pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, 2, tkReqs)
}

func TestChallengePolicy_KVThenCAE(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", kvChallenge1),
		mock.WithStatusCode(401),
		mock.WithPredicate(requireToken(t, "")),
	)
	srv.AppendResponse() // when a response's predicate returns true, srv pops the following one
	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", caeChallenge1),
		mock.WithStatusCode(401),
		mock.WithPredicate(requireToken(t, challengedToken)),
	)
	srv.AppendResponse() // when a response's predicate returns true, srv pops the following one
	srv.AppendResponse(
		mock.WithPredicate(requireToken(t, claimsToken)),
	)
	srv.AppendResponse() // when a response's predicate returns true, srv pops the following one

	tkReqs := 0
	cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		require.True(t, tro.EnableCAE)
		tkReqs += 1
		tk := challengedToken
		switch tkReqs {
		case 1:
			require.Empty(t, tro.Claims)
		case 2:
			tk = claimsToken
			require.Equal(t, "testing1", tro.Claims)
		default:
			t.Fatal("unexpected token request")
		}
		return azcore.AccessToken{Token: tk, ExpiresOn: time.Now().Add(time.Hour)}, nil
	})
	p := NewKeyVaultChallengePolicy(cred, nil)
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)
	req, err := runtime.NewRequest(context.Background(), "GET", "https://42.vault.azure.net")
	require.NoError(t, err)
	res, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, tkReqs, 2)
}

func TestChallengePolicy_TwoCAEChallenges(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", kvChallenge1),
		mock.WithStatusCode(401),
		mock.WithPredicate(requireToken(t, "")),
	)
	srv.AppendResponse() // when a response's predicate returns true, srv pops the following one
	srv.AppendResponse()

	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", caeChallenge1),
		mock.WithStatusCode(401),
		mock.WithPredicate(requireToken(t, challengedToken)),
	)
	srv.AppendResponse() // when a response's predicate returns true, srv pops the following one
	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", caeChallenge2),
		mock.WithStatusCode(401),
		mock.WithPredicate(requireToken(t, claimsToken)),
	)
	srv.AppendResponse() // when a response's predicate returns true, srv pops the following one
	tkReqs := 0
	cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		require.True(t, tro.EnableCAE)
		tk := challengedToken
		tkReqs += 1
		switch tkReqs {
		case 1:
			require.Empty(t, tro.Claims)
		case 2:
			tk = claimsToken
			require.Equal(t, "testing1", tro.Claims)
		default:
			t.Fatal("unexpected token request")
		}
		return azcore.AccessToken{Token: tk, ExpiresOn: time.Now().Add(time.Hour)}, nil
	})
	p := NewKeyVaultChallengePolicy(cred, nil)
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)

	// req 1 kv then regular
	req, err := runtime.NewRequest(context.Background(), "GET", "https://42.vault.azure.net")
	require.NoError(t, err)
	res, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, tkReqs, 1)

	// req 2 cae twice
	req, err = runtime.NewRequest(context.Background(), "GET", "https://42.vault.azure.net")
	require.NoError(t, err)
	res, err = pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, 401, res.StatusCode)
	require.Equal(t, caeChallenge2, res.Header.Get("WWW-Authenticate"))
	require.Equal(t, tkReqs, 2)
}

func TestParseTenant(t *testing.T) {
	actual := parseTenant("")
	require.Empty(t, actual)

	expected := "00000000-0000-0000-0000-000000000000"
	sampleURL := "https://login.microsoftonline.com/" + expected
	actual = parseTenant(sampleURL)
	require.Equal(t, expected, actual, "tenant was not properly parsed, got %s, expected %s", actual, expected)
}

func TestChallengePolicy_ConcurrentRequests(t *testing.T) {
	concurrentRequestCount := 3

	serverAuthenticateRequests := atomic.Int32{}
	serverAuthenticatedRequests := atomic.Int32{}
	var srv *httptest.Server
	srv = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Values("Authorization")
		if len(authz) == 0 {
			// Initial request without Authorization header. Send a
			// challenge response to the client.
			serverAuthenticateRequests.Add(1)
			resource := srv.URL
			w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Bearer authorization="https://login.microsoftonline.com/{tenant}", resource="%s"`, resource))
			w.WriteHeader(401)
		} else {
			// Authenticated request.
			serverAuthenticatedRequests.Add(1)
			if len(authz) != 1 || authz[0] != "Bearer ***" {
				t.Errorf(`unexpected Authorization "%s"`, authz)
			}
			// Return nothing.
			w.WriteHeader(200)
		}
	}))
	defer srv.Close()
	srv.StartTLS()

	cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		return azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
	})
	p := NewKeyVaultChallengePolicy(cred, &KeyVaultChallengePolicyOptions{
		// Challenge resource verification will always fail because we
		// use local IPs instead of domain names and subdomains in this
		// test.
		DisableChallengeResourceVerification: true,
	})
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv.Client()},
	)

	wg := sync.WaitGroup{}
	for i := 0; i < concurrentRequestCount; i += 1 {
		go (func() {
			defer wg.Done()
			req, err := runtime.NewRequest(context.Background(), "GET", srv.URL)
			require.NoError(t, err)
			res, err := pl.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()
		})()
		wg.Add(1)
	}
	wg.Wait()

	require.GreaterOrEqual(t, int(serverAuthenticateRequests.Load()), 1, "client should have sent at least one preflight request")
	require.LessOrEqual(t, int(serverAuthenticateRequests.Load()), concurrentRequestCount, "client should have sent no more preflight requests than client requests")
	require.EqualValues(t, concurrentRequestCount, serverAuthenticatedRequests.Load(), "client preflight request count should equal server preflight request count")
}
