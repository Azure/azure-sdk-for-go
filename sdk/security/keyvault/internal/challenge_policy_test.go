// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
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
	"github.com/stretchr/testify/assert"
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
	concurrentRequestCount := 10

	var srv *httptest.Server
	srv = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Values("Authorization")
		if len(authz) == 0 {
			// Initial request without Authorization header. Send a
			// challenge response to the client.
			resource := srv.URL
			w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Bearer authorization="https://login.microsoftonline.com/{tenant}", resource="%s"`, resource))
			w.WriteHeader(401)
		} else {
			// Authenticated request.
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
}

func TestChallengePolicy_AutomaticMTLSPoP(t *testing.T) {
	for _, bound := range []bool{false, true} {
		t.Run(fmt.Sprint(bound), func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				n := calls.Add(1)
				if n == 1 {
					assert.Empty(t, body, "the cold challenge must not send the request body")
					assert.Empty(t, r.TLS.PeerCertificates)
					assert.Empty(t, r.Header.Get("Authorization"))
					w.Header().Set("WWW-Authenticate", kvChallenge1)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				assert.Equal(t, "request body", string(body))
				if bound {
					assert.Equal(t, "mtls_pop token", r.Header.Get("Authorization"))
					assert.Equal(t, "true", r.Header.Get("x-ms-tokenboundauth"))
					assert.Len(t, r.TLS.PeerCertificates, 1)
				} else {
					assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
					assert.Empty(t, r.Header.Get("x-ms-tokenboundauth"))
					assert.Empty(t, r.TLS.PeerCertificates)
				}
				if bound && n == 2 {
					w.Header().Set("WWW-Authenticate", kvChallenge1)
					w.WriteHeader(http.StatusUnauthorized)
					_, err := io.WriteString(w, `{"error":{"message":"[MtlsCnfClaimRequestDataValidationFailed]"}}`)
					assert.NoError(t, err)
				}
			}))
			server.TLS = &tls.Config{ClientAuth: tls.RequestClientCert}
			server.StartTLS()
			defer server.Close()
			acquisitions := 0
			cred := credentialFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
				acquisitions++
				require.True(t, tro.EnableMTLSPoP)
				require.True(t, tro.EnableCAE)
				tk := azcore.AccessToken{Token: "token", ExpiresOn: time.Now().Add(time.Hour)}
				if bound {
					tk.BindingCertificate = &server.TLS.Certificates[0]
				}
				return tk, nil
			})
			pl := runtime.NewPipeline("", "", runtime.PipelineOptions{
				PerCall: []policy.Policy{NewTokenBindingPolicy()},
				PerRetry: []policy.Policy{
					NewKeyVaultChallengePolicy(cred, &KeyVaultChallengePolicyOptions{DisableChallengeResourceVerification: true}),
				},
			}, &policy.ClientOptions{
				Transport: server.Client(), Retry: policy.RetryOptions{MaxRetries: 1, RetryDelay: -1},
			})
			for range 2 {
				req, err := runtime.NewRequest(context.Background(), http.MethodPost, server.URL)
				require.NoError(t, err)
				require.NoError(t, req.SetBody(streaming.NopCloser(strings.NewReader("request body")), "text/plain"))
				res, err := pl.Do(req)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, res.StatusCode)
				require.NoError(t, res.Body.Close())
			}
			require.Equal(t, 1, acquisitions)
			wantCalls := int32(3)
			if bound {
				wantCalls++
			}
			require.Equal(t, wantCalls, calls.Load())
		})
	}
}

func TestChallengePolicy_MTLSPoPRequiresTokenBindingPolicy(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized), mock.WithHeader("WWW-Authenticate", kvChallenge1))
	srv.AppendResponse()
	acquisitions := 0
	cred := credentialFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		acquisitions++
		require.False(t, tro.EnableMTLSPoP, "a client without the token binding policy can't retry binding failures")
		return azcore.AccessToken{Token: "token", ExpiresOn: time.Now().Add(time.Hour)}, nil
	})
	pl := runtime.NewPipeline("", "", runtime.PipelineOptions{PerRetry: []policy.Policy{
		NewKeyVaultChallengePolicy(cred, &KeyVaultChallengePolicyOptions{DisableChallengeResourceVerification: true}),
	}}, &policy.ClientOptions{Transport: srv})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	res, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, 1, acquisitions)
}

func TestChallengePolicy_BindingValidationRetry(t *testing.T) {
	const markerBody = `{"error":{"code":"Unauthorized","message":"[MtlsCnfClaimRequestDataValidationFailed] binding rejected"}}`
	for _, tc := range []struct {
		name, body                                                 string
		status, attempts                                           int
		bearer, challenge, recover, disable, override, rejectRetry bool
	}{
		{name: "bounded with challenge", body: markerBody, status: 401, attempts: 3, challenge: true},
		{name: "bounded without challenge", body: markerBody, status: 401, attempts: 3},
		{name: "recovery", body: markerBody, status: 401, attempts: 2, challenge: true, recover: true},
		{name: "disabled retries", body: markerBody, status: 401, attempts: 1, challenge: true, disable: true},
		{name: "operation override", body: markerBody, status: 401, attempts: 2, challenge: true, override: true},
		{name: "custom retry veto", body: markerBody, status: 401, attempts: 1, challenge: true, rejectRetry: true},
		{name: "bearer", body: markerBody, status: 401, attempts: 1, bearer: true},
		{name: "RBAC", body: markerBody, status: 403, attempts: 1},
		{name: "ordinary unauthorized", body: `{"error":{"message":"unauthorized"}}`, status: 401, attempts: 1},
		{name: "marker outside message", body: `{"marker":"[MtlsCnfClaimRequestDataValidationFailed]"}`, status: 401, attempts: 1},
		{name: "wrong property case", body: `{"Error":{"Message":"[MtlsCnfClaimRequestDataValidationFailed]"}}`, status: 401, attempts: 1},
		{name: "last property wins", body: `{"error":{"message":"[MtlsCnfClaimRequestDataValidationFailed]"},"error":{}}`, status: 401, attempts: 1},
		{name: "invalid JSON", body: tokenBindingValidationFailure, status: 401, attempts: 1},
		{name: "ordinary service retry", body: `{}`, status: 500, attempts: 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				n := calls.Add(1)
				if !tc.bearer {
					assert.Equal(t, "mtls_pop token", r.Header.Get("Authorization"))
					assert.Len(t, r.TLS.PeerCertificates, 1)
				}
				if tc.recover && n > 1 {
					return
				}
				if tc.challenge {
					w.Header().Set("WWW-Authenticate", kvChallenge1)
				}
				w.WriteHeader(tc.status)
				_, err := io.WriteString(w, tc.body)
				assert.NoError(t, err)
			}))
			server.TLS = &tls.Config{ClientAuth: tls.RequestClientCert}
			server.StartTLS()
			defer server.Close()
			acquisitions := 0
			cred := credentialFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
				acquisitions++
				require.True(t, tro.EnableMTLSPoP, "binding failures must not downgrade")
				tk := azcore.AccessToken{Token: "token", ExpiresOn: time.Now().Add(time.Hour)}
				if !tc.bearer {
					tk.BindingCertificate = &server.TLS.Certificates[0]
				}
				return tk, nil
			})
			p := NewKeyVaultChallengePolicy(cred, nil).(*keyVaultAuthorizer)
			p.setTokenRequestOptions(policy.TokenRequestOptions{
				TenantID: "tenant", Scopes: []string{"https://vault.azure.net/.default"}, EnableMTLSPoP: true,
			})
			retry := policy.RetryOptions{MaxRetries: 2, RetryDelay: -1}
			if tc.disable {
				retry.MaxRetries = -1
			}
			if tc.rejectRetry {
				retry.ShouldRetry = func(res *http.Response, err error) bool {
					require.NoError(t, err)
					require.Equal(t, tc.status, res.StatusCode)
					return false
				}
			}
			pl := runtime.NewPipeline("", "", runtime.PipelineOptions{
				PerCall:  []policy.Policy{NewTokenBindingPolicy()},
				PerRetry: []policy.Policy{p},
			}, &policy.ClientOptions{Transport: server.Client(), Retry: retry})
			ctx := context.Background()
			if tc.override {
				ctx = policy.WithRetryOptions(ctx, policy.RetryOptions{MaxRetries: 1, RetryDelay: -1})
			}
			req, err := runtime.NewRequest(ctx, http.MethodGet, server.URL)
			require.NoError(t, err)
			res, err := pl.Do(req)
			require.NoError(t, err)
			if tc.recover {
				require.Equal(t, http.StatusOK, res.StatusCode)
			} else {
				require.Equal(t, tc.status, res.StatusCode)
				body, err := runtime.Payload(res)
				require.NoError(t, err)
				require.Equal(t, tc.body, string(body), "classification must preserve the response body")
			}
			require.Equal(t, int32(tc.attempts), calls.Load())
			require.Equal(t, 1, acquisitions, "binding retries must reuse authentication instead of resending challenges")
		})
	}
}
