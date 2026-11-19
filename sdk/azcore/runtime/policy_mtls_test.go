// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMTLSPoPResourcePipeline(t *testing.T) {
	var calls atomic.Int32
	var cert tls.Certificate
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "mtls_pop "+tokenValue, r.Header.Get("Authorization"))
		assert.Empty(t, r.Header.Get("x-ms-tokenboundauth"), "service-specific binding headers belong to service policies")
		assert.Equal(t, uint16(tls.VersionTLS12), r.TLS.Version)
		assert.Len(t, r.TLS.PeerCertificates, 1)
		if len(r.TLS.PeerCertificates) == 1 {
			assert.Equal(t, cert.Certificate[0], r.TLS.PeerCertificates[0].Raw)
		}
		if calls.Add(1) == 1 {
			w.Header().Set("WWW-Authenticate", `Bearer error="insufficient_claims", claims="e30="`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	server.TLS = &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAnyClientCert,
	}
	server.StartTLS()
	defer server.Close()
	cert = server.TLS.Certificates[0]
	var acquisitions atomic.Int32
	cred := mockCredential{getTokenImpl: func(_ context.Context, opts policy.TokenRequestOptions) (exported.AccessToken, error) {
		assert.True(t, opts.EnableMTLSPoP)
		assert.True(t, opts.EnableCAE)
		assert.Equal(t, []string{scope}, opts.Scopes)
		if acquisitions.Add(1) == 2 {
			assert.Equal(t, "{}", opts.Claims)
		}
		return exported.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour), BindingCertificate: &cert}, nil
	}}
	base := server.Client()
	tr, ok := base.Transport.(*http.Transport)
	require.True(t, ok)
	p := NewPipeline("test", "v1.0.0", PipelineOptions{
		PerRetry: []policy.Policy{NewBearerTokenPolicy(cred, []string{scope}, &policy.BearerTokenOptions{
			AuthorizationHandler: policy.AuthorizationHandler{OnRequest: requestMTLSPoP},
		})},
	}, &policy.ClientOptions{Transport: base})
	req, err := NewRequest(context.Background(), http.MethodGet, server.URL)
	require.NoError(t, err)
	resp, err := p.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(2), acquisitions.Load())
	assert.Nil(t, tr.TLSClientConfig.GetClientCertificate, "the caller's transport must not be modified")

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			req, err := NewRequest(context.Background(), http.MethodGet, server.URL)
			assert.NoError(t, err)
			if err != nil {
				return
			}
			resp, err := p.Do(req)
			assert.NoError(t, err)
			if resp != nil {
				assert.NoError(t, resp.Body.Close())
			}
		})
	}
	wg.Wait()
	assert.Equal(t, int32(2), acquisitions.Load(), "bound tokens should be cached")
}

func TestMTLSPoPAdaptiveAuthentication(t *testing.T) {
	for _, tc := range []struct {
		name                  string
		requestEnabled, bound bool
	}{
		{name: "service requests bound token", requestEnabled: true, bound: true},
		{name: "credential returns bound token", bound: true},
		{name: "service accepts bearer token", requestEnabled: true},
		{name: "ordinary bearer request"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var requests atomic.Int32
			var cert tls.Certificate
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				assert.Equal(t, "preserved", r.Header.Get("X-Test-Middleware"))
				if tc.bound {
					assert.Equal(t, "mtls_pop "+tokenValue, r.Header.Get("Authorization"))
					assert.Equal(t, uint16(tls.VersionTLS12), r.TLS.Version)
					if assert.Len(t, r.TLS.PeerCertificates, 1) {
						assert.Equal(t, cert.Certificate[0], r.TLS.PeerCertificates[0].Raw)
					}
				} else {
					assert.Equal(t, "Bearer "+tokenValue, r.Header.Get("Authorization"))
					assert.Empty(t, r.Header.Get("x-ms-tokenboundauth"))
					assert.Equal(t, uint16(tls.VersionTLS13), r.TLS.Version,
						"bearer requests must retain the caller's TLS configuration")
					assert.Empty(t, r.TLS.PeerCertificates)
				}
				w.WriteHeader(http.StatusOK)
			}))
			server.TLS = &tls.Config{MinVersion: tls.VersionTLS12, ClientAuth: tls.RequestClientCert}
			server.StartTLS()
			defer server.Close()
			cert = server.TLS.Certificates[0]
			base := server.Client()
			tr, ok := base.Transport.(*http.Transport)
			require.True(t, ok)
			if !tc.bound {
				tr.TLSClientConfig.MinVersion = tls.VersionTLS13
			}
			var acquisitions, middlewareCalls int
			auth := NewBearerTokenPolicy(mockCredential{getTokenImpl: func(_ context.Context, opts policy.TokenRequestOptions) (exported.AccessToken, error) {
				acquisitions++
				assert.Equal(t, tc.requestEnabled, opts.EnableMTLSPoP)
				tk := exported.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour)}
				if tc.bound {
					tk.BindingCertificate = &cert
				}
				return tk, nil
			}}, nil, &policy.BearerTokenOptions{AuthorizationHandler: policy.AuthorizationHandler{
				OnRequest: func(_ *policy.Request, authNZ func(policy.TokenRequestOptions) error) error {
					return authNZ(policy.TokenRequestOptions{Scopes: []string{scope}, EnableMTLSPoP: tc.requestEnabled})
				},
			}})
			p := NewPipeline("test", "v1.0.0", PipelineOptions{
				PerRetry: []policy.Policy{auth},
			}, &policy.ClientOptions{
				Transport: base,
				PerRetryPolicies: []policy.Policy{exported.PolicyFunc(func(req *policy.Request) (*http.Response, error) {
					middlewareCalls++
					req.Raw().Header.Set("X-Test-Middleware", "preserved")
					return req.Next()
				})},
			})
			for range 2 {
				req, err := NewRequest(context.Background(), http.MethodGet, server.URL)
				require.NoError(t, err)
				resp, err := p.Do(req)
				require.NoError(t, err)
				require.NoError(t, resp.Body.Close())
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}
			assert.Equal(t, 1, acquisitions)
			assert.Equal(t, 2, middlewareCalls)
			assert.Equal(t, int32(2), requests.Load())
			assert.Nil(t, tr.TLSClientConfig.GetClientCertificate)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			req, err := NewRequest(ctx, http.MethodGet, server.URL)
			require.NoError(t, err)
			resp, err := p.Do(req)
			require.ErrorIs(t, err, context.Canceled)
			assert.Nil(t, resp)
			assert.Equal(t, int32(2), requests.Load())
		})
	}
}

func TestMTLSPoPBoundToBearerChallenge(t *testing.T) {
	var requests atomic.Int32
	var cert tls.Certificate
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call := requests.Add(1)
		if call == 1 {
			assert.Equal(t, "mtls_pop "+tokenValue, r.Header.Get("Authorization"))
			assert.Equal(t, uint16(tls.VersionTLS12), r.TLS.Version)
			if assert.Len(t, r.TLS.PeerCertificates, 1) {
				assert.Equal(t, cert.Certificate[0], r.TLS.PeerCertificates[0].Raw)
			}
			w.Header().Set("WWW-Authenticate", `Bearer resource="https://example.com"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		assert.Equal(t, "Bearer bearer-token", r.Header.Get("Authorization"))
		assert.Empty(t, r.Header.Get("x-ms-tokenboundauth"))
		assert.Equal(t, uint16(tls.VersionTLS13), r.TLS.Version)
		assert.Empty(t, r.TLS.PeerCertificates, "bearer requests must not reuse the bound connection")
		if call == 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	server.TLS = &tls.Config{MinVersion: tls.VersionTLS12, ClientAuth: tls.RequestClientCert}
	server.StartTLS()
	defer server.Close()
	cert = server.TLS.Certificates[0]
	var acquisitions int
	challenged := false
	auth := NewBearerTokenPolicy(mockCredential{getTokenImpl: func(_ context.Context, opts policy.TokenRequestOptions) (exported.AccessToken, error) {
		acquisitions++
		tk := exported.AccessToken{Token: "bearer-token", ExpiresOn: time.Now().Add(time.Hour)}
		if opts.EnableMTLSPoP {
			tk.Token = tokenValue
			tk.BindingCertificate = &cert
		}
		return tk, nil
	}}, nil, &policy.BearerTokenOptions{AuthorizationHandler: policy.AuthorizationHandler{
		OnRequest: func(_ *policy.Request, authNZ func(policy.TokenRequestOptions) error) error {
			return authNZ(policy.TokenRequestOptions{Scopes: []string{scope}, EnableMTLSPoP: !challenged})
		},
		OnChallenge: func(_ *policy.Request, _ *http.Response, authNZ func(policy.TokenRequestOptions) error) error {
			challenged = true
			return authNZ(policy.TokenRequestOptions{Scopes: []string{scope}})
		},
	}})
	p := NewPipeline("test", "v1.0.0", PipelineOptions{
		PerRetry: []policy.Policy{auth},
	}, &policy.ClientOptions{
		Transport: server.Client(),
		Retry:     policy.RetryOptions{MaxRetries: 1, RetryDelay: time.Millisecond, MaxRetryDelay: time.Millisecond},
	})
	for range 2 {
		req, err := NewRequest(context.Background(), http.MethodGet, server.URL)
		require.NoError(t, err)
		resp, err := p.Do(req)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	assert.Equal(t, 2, acquisitions, "the challenge acquires a new bearer token; retries and later reads reuse it")
	assert.Equal(t, int32(4), requests.Load())
}

func TestMTLSPoPTokenRenewal(t *testing.T) {
	var calls atomic.Int32
	var certs [2]tls.Certificate
	tokens := [2]string{tokenValue, "renewed-token"}
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call := calls.Add(1)
		index := 0
		if call > 1 {
			index = 1
		}
		assert.Equal(t, "mtls_pop "+tokens[index], r.Header.Get("Authorization"))
		assert.Equal(t, uint16(tls.VersionTLS12), r.TLS.Version)
		if assert.Len(t, r.TLS.PeerCertificates, 1) {
			assert.Equal(t, certs[index].Certificate[0], r.TLS.PeerCertificates[0].Raw,
				"renewal must not reuse a connection authenticated with the old certificate")
		}
		if call == 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	server.TLS = &tls.Config{MinVersion: tls.VersionTLS12, ClientAuth: tls.RequireAnyClientCert}
	server.StartTLS()
	defer server.Close()
	certs[0] = server.TLS.Certificates[0]
	leaf, err := x509.ParseCertificate(certs[0].Certificate[0])
	require.NoError(t, err)
	leaf.SerialNumber = new(big.Int).Add(leaf.SerialNumber, big.NewInt(1))
	der, err := x509.CreateCertificate(rand.Reader, leaf, leaf, leaf.PublicKey, certs[0].PrivateKey)
	require.NoError(t, err)
	certs[1] = tls.Certificate{Certificate: [][]byte{der}, PrivateKey: certs[0].PrivateKey}

	var acquisitions int
	auth := NewBearerTokenPolicy(mockCredential{getTokenImpl: func(_ context.Context, opts policy.TokenRequestOptions) (exported.AccessToken, error) {
		assert.True(t, opts.EnableMTLSPoP)
		assert.Equal(t, []string{scope}, opts.Scopes)
		index := 0
		if acquisitions > 0 {
			index = 1
		}
		acquisitions++
		return exported.AccessToken{
			Token: tokens[index], ExpiresOn: time.Now().Add(time.Hour), BindingCertificate: &certs[index],
		}, nil
	}}, []string{scope}, &policy.BearerTokenOptions{
		AuthorizationHandler: policy.AuthorizationHandler{OnRequest: requestMTLSPoP},
	})
	p := NewPipeline("test", "v1.0.0", PipelineOptions{
		PerRetry: []policy.Policy{auth},
	}, &policy.ClientOptions{
		Transport: server.Client(),
		Retry:     policy.RetryOptions{MaxRetries: 1, RetryDelay: time.Millisecond, MaxRetryDelay: time.Millisecond},
	})
	for i := range 3 {
		if i == 1 {
			auth.mtlsResource.Expire()
		}
		req, err := NewRequest(context.Background(), http.MethodGet, server.URL)
		require.NoError(t, err)
		resp, err := p.Do(req)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	assert.Equal(t, 2, acquisitions, "the retry and subsequent request must reuse the renewed bound token")
	assert.Equal(t, int32(4), calls.Load(), "three successful reads and one retried 503")
}

func TestMTLSPoPReacquiredCertificateReusesConnection(t *testing.T) {
	var connections atomic.Int32
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Len(t, r.TLS.PeerCertificates, 1)
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			connections.Add(1)
		}
	}
	server.TLS = &tls.Config{MinVersion: tls.VersionTLS12, ClientAuth: tls.RequireAnyClientCert}
	server.StartTLS()
	defer server.Close()
	var acquisitions int
	auth := NewBearerTokenPolicy(mockCredential{getTokenImpl: func(context.Context, policy.TokenRequestOptions) (exported.AccessToken, error) {
		acquisitions++
		// credentials such as MSAL return a new *tls.Certificate for the same certificate on each acquisition
		cert := server.TLS.Certificates[0]
		return exported.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour), BindingCertificate: &cert}, nil
	}}, []string{scope}, &policy.BearerTokenOptions{
		AuthorizationHandler: policy.AuthorizationHandler{OnRequest: requestMTLSPoP},
	})
	p := NewPipeline("test", "v1.0.0", PipelineOptions{PerRetry: []policy.Policy{auth}}, &policy.ClientOptions{Transport: server.Client()})
	for range 3 {
		auth.mtlsResource.Expire()
		req, err := NewRequest(context.Background(), http.MethodGet, server.URL)
		require.NoError(t, err)
		resp, err := p.Do(req)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	assert.Equal(t, 3, acquisitions)
	assert.Equal(t, int32(1), connections.Load(), "reacquiring the same certificate must reuse its connection")
}

func TestMTLSPoPReacquiredCertificateUsesLatestSigner(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Len(t, r.TLS.PeerCertificates, 1)
	}))
	server.TLS = &tls.Config{MinVersion: tls.VersionTLS12, ClientAuth: tls.RequireAnyClientCert}
	server.StartTLS()
	defer server.Close()
	signer, ok := server.TLS.Certificates[0].PrivateKey.(crypto.Signer)
	require.True(t, ok)
	var acquisitions int
	auth := NewBearerTokenPolicy(mockCredential{getTokenImpl: func(context.Context, policy.TokenRequestOptions) (exported.AccessToken, error) {
		acquisitions++
		cert := server.TLS.Certificates[0]
		if acquisitions == 1 {
			cert.PrivateKey = failingSigner{signer}
		}
		return exported.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour), BindingCertificate: &cert}, nil
	}}, []string{scope}, &policy.BearerTokenOptions{
		AuthorizationHandler: policy.AuthorizationHandler{OnRequest: requestMTLSPoP},
	})
	p := NewPipeline("test", "v1.0.0", PipelineOptions{PerRetry: []policy.Policy{auth}}, &policy.ClientOptions{
		Retry:     policy.RetryOptions{MaxRetries: -1},
		Transport: server.Client(),
	})
	req, err := NewRequest(context.Background(), http.MethodGet, server.URL)
	require.NoError(t, err)
	_, err = p.Do(req)
	require.ErrorContains(t, err, "key handle closed")

	auth.mtlsResource.Expire()
	req, err = NewRequest(context.Background(), http.MethodGet, server.URL)
	require.NoError(t, err)
	resp, err := p.Do(req)
	require.NoError(t, err, "a new handshake must present the most recently acquired signer")
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, acquisitions)
}

// failingSigner simulates a key handle that can no longer sign
type failingSigner struct{ crypto.Signer }

func (failingSigner) Sign(io.Reader, []byte, crypto.SignerOpts) ([]byte, error) {
	return nil, errors.New("key handle closed")
}

func TestMTLSPoPChallengeAuthentication(t *testing.T) {
	for _, bound := range []bool{true, false} {
		name := "bound token"
		if !bound {
			name = "unbound token"
		}
		t.Run(name, func(t *testing.T) {
			var requests atomic.Int32
			var cert tls.Certificate
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if requests.Add(1) == 1 {
					assert.Equal(t, uint16(tls.VersionTLS13), r.TLS.Version)
					assert.Empty(t, r.Header.Get("Authorization"))
					assert.Empty(t, r.Header.Get("x-ms-tokenboundauth"))
					assert.Len(t, r.TLS.PeerCertificates, 1, "the cold challenge must preserve the caller's TLS configuration")
					w.Header().Set("WWW-Authenticate", `Bearer resource="https://example.com"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if bound {
					assert.Equal(t, uint16(tls.VersionTLS12), r.TLS.Version)
					assert.Equal(t, "mtls_pop "+tokenValue, r.Header.Get("Authorization"))
				} else {
					assert.Equal(t, uint16(tls.VersionTLS13), r.TLS.Version)
					assert.Equal(t, "Bearer "+tokenValue, r.Header.Get("Authorization"))
					assert.Empty(t, r.Header.Get("x-ms-tokenboundauth"))
				}
				if assert.Len(t, r.TLS.PeerCertificates, 1) {
					assert.Equal(t, cert.Certificate[0], r.TLS.PeerCertificates[0].Raw)
				}
				w.WriteHeader(http.StatusOK)
			}))
			server.TLS = &tls.Config{MinVersion: tls.VersionTLS12, ClientAuth: tls.RequestClientCert}
			server.StartTLS()
			defer server.Close()
			cert = server.TLS.Certificates[0]
			base := server.Client()
			tr, ok := base.Transport.(*http.Transport)
			require.True(t, ok)
			tr.TLSClientConfig.Certificates = []tls.Certificate{cert}
			tr.TLSClientConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
				return &cert, nil
			}
			var acquisitions int
			cred := mockCredential{getTokenImpl: func(_ context.Context, opts policy.TokenRequestOptions) (exported.AccessToken, error) {
				acquisitions++
				assert.True(t, opts.EnableMTLSPoP)
				assert.True(t, opts.EnableCAE)
				assert.Equal(t, []string{scope}, opts.Scopes)
				token := exported.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour)}
				if bound {
					token.BindingCertificate = &cert
				}
				return token, nil
			}}
			challenged := false
			p := NewPipeline("test", "v1.0.0", PipelineOptions{
				PerRetry: []policy.Policy{NewBearerTokenPolicy(cred, nil, &policy.BearerTokenOptions{
					AuthorizationHandler: policy.AuthorizationHandler{
						OnRequest: func(_ *policy.Request, authNZ func(policy.TokenRequestOptions) error) error {
							if !challenged {
								return nil
							}
							return requestMTLSPoP(nil, authNZ)
						},
						OnChallenge: func(_ *policy.Request, _ *http.Response, authNZ func(policy.TokenRequestOptions) error) error {
							challenged = true
							return requestMTLSPoP(nil, authNZ)
						},
					},
				})},
			}, &policy.ClientOptions{Transport: base})
			for range 2 {
				req, err := NewRequest(context.Background(), http.MethodGet, server.URL)
				require.NoError(t, err)
				resp, err := p.Do(req)
				require.NoError(t, err)
				require.NoError(t, resp.Body.Close())
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}
			assert.Equal(t, int32(3), requests.Load(), "one challenge followed by two authenticated requests")
			assert.Equal(t, 1, acquisitions, "the token should be cached")
			assert.Len(t, tr.TLSClientConfig.Certificates, 1, "the caller's transport must not be modified")
		})
	}
}

func TestMTLSPoPDoesNotFallbackOnCredentialFailure(t *testing.T) {
	for _, expected := range []error{errors.New("attestation failed"), context.Canceled} {
		t.Run(expected.Error(), func(t *testing.T) {
			var acquisitions int
			cred := mockCredential{getTokenImpl: func(_ context.Context, opts policy.TokenRequestOptions) (exported.AccessToken, error) {
				acquisitions++
				assert.True(t, opts.EnableMTLSPoP)
				return exported.AccessToken{}, expected
			}}
			srv, close := mock.NewTLSServer()
			defer close()
			p := NewPipeline("test", "v1.0.0", PipelineOptions{
				PerRetry: []policy.Policy{NewBearerTokenPolicy(cred, []string{scope}, &policy.BearerTokenOptions{
					AuthorizationHandler: policy.AuthorizationHandler{OnRequest: requestMTLSPoP},
				})},
			}, &policy.ClientOptions{Transport: srv})
			req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
			require.NoError(t, err)
			resp, err := p.Do(req)
			require.ErrorIs(t, err, expected)
			require.Nil(t, resp)
			assert.Equal(t, 1, acquisitions, "a failed PoP acquisition must not trigger bearer acquisition")
			assert.Zero(t, srv.Requests())
		})
	}
}

func TestMTLSPoPSeparateTokenCaches(t *testing.T) {
	for _, bound := range []bool{false, true} {
		name := "bearer response"
		if bound {
			name = "bound response"
		}
		t.Run(name, func(t *testing.T) {
			var bearerRequests, popRequests int
			b := NewBearerTokenPolicy(mockCredential{getTokenImpl: func(_ context.Context, opts policy.TokenRequestOptions) (exported.AccessToken, error) {
				tk := exported.AccessToken{Token: "bearer-request", ExpiresOn: time.Now().Add(time.Hour)}
				if opts.EnableMTLSPoP {
					popRequests++
					tk.Token = "pop-request"
					if bound {
						tk.BindingCertificate = &tls.Certificate{}
					}
				} else {
					bearerRequests++
				}
				return tk, nil
			}}, []string{scope}, nil)
			for _, enabled := range []bool{false, true, false, true} {
				req, err := NewRequest(context.Background(), http.MethodGet, "https://example.com")
				require.NoError(t, err)
				require.NoError(t, b.authenticateAndAuthorize(req)(policy.TokenRequestOptions{Scopes: []string{scope}, EnableMTLSPoP: enabled}))
				expected := "Bearer bearer-request"
				if enabled {
					expected = "Bearer pop-request"
					if bound {
						expected = "mtls_pop pop-request"
					}
				}
				assert.Equal(t, expected, req.Raw().Header.Get("Authorization"))
			}
			assert.Equal(t, 1, bearerRequests)
			assert.Equal(t, 1, popRequests)
		})
	}
}

func TestMTLSPoPAdaptiveCustomTransport(t *testing.T) {
	for _, bound := range []bool{false, true} {
		name := "bearer response"
		if bound {
			name = "bound response"
		}
		t.Run(name, func(t *testing.T) {
			srv, close := mock.NewTLSServer()
			defer close()
			srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized), mock.WithHeader("WWW-Authenticate", `Bearer resource="https://example.com"`))
			srv.AppendResponse()
			var acquisitions int
			cred := mockCredential{getTokenImpl: func(_ context.Context, opts policy.TokenRequestOptions) (exported.AccessToken, error) {
				acquisitions++
				assert.True(t, opts.EnableMTLSPoP)
				tk := exported.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour)}
				if bound {
					tk.BindingCertificate = &tls.Certificate{Certificate: [][]byte{{1}}, PrivateKey: mtlsTestSigner{}}
				}
				return tk, nil
			}}
			p := NewPipeline("test", "v1.0.0", PipelineOptions{
				PerRetry: []policy.Policy{NewBearerTokenPolicy(cred, nil, &policy.BearerTokenOptions{
					AuthorizationHandler: policy.AuthorizationHandler{
						OnRequest: func(*policy.Request, func(policy.TokenRequestOptions) error) error {
							return nil
						},
						OnChallenge: func(_ *policy.Request, _ *http.Response, authNZ func(policy.TokenRequestOptions) error) error {
							return authNZ(policy.TokenRequestOptions{Scopes: []string{scope}, EnableMTLSPoP: true})
						},
					},
				})},
			}, &policy.ClientOptions{Transport: srv})
			req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
			require.NoError(t, err)
			resp, err := p.Do(req)
			if bound {
				require.ErrorContains(t, err, "*http.Client")
				assert.Nil(t, resp)
				assert.Equal(t, 1, srv.Requests(), "the cold challenge must use the caller's transport before token binding is known")
			} else {
				require.NoError(t, err)
				require.NoError(t, resp.Body.Close())
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, 2, srv.Requests(), "the challenge and bearer request must both use the caller's transport")
			}
			assert.Equal(t, 1, acquisitions, "transport incompatibility must not trigger bearer acquisition")
		})
	}
}

// mtlsTestSigner satisfies crypto.Signer for tests that fail before a handshake would call it
type mtlsTestSigner struct{ crypto.Signer }

func requestMTLSPoP(_ *policy.Request, authNZ func(policy.TokenRequestOptions) error) error {
	return authNZ(policy.TokenRequestOptions{Scopes: []string{scope}, EnableMTLSPoP: true})
}

func TestMTLSPoPCertificateExpiration(t *testing.T) {
	now := time.Now()
	req, err := NewRequest(context.Background(), http.MethodGet, "https://example.com")
	require.NoError(t, err)
	for _, lifetime := range []time.Duration{time.Minute, 2 * time.Hour} {
		t.Run(lifetime.String(), func(t *testing.T) {
			tk := exported.AccessToken{
				Token: tokenValue, ExpiresOn: now.Add(time.Hour), RefreshOn: now.Add(30 * time.Minute),
				BindingCertificate: &tls.Certificate{Leaf: &x509.Certificate{NotAfter: now.Add(lifetime)}},
			}
			state := acquiringResourceState{
				req: req,
				p: &BearerTokenPolicy{cred: mockCredential{getTokenImpl: func(context.Context, policy.TokenRequestOptions) (exported.AccessToken, error) {
					return tk, nil
				}}},
			}
			got, expiration, err := acquire(state)
			require.NoError(t, err)
			expected := tk.ExpiresOn
			if lifetime < time.Hour {
				expected = tk.BindingCertificate.Leaf.NotAfter
			}
			assert.Equal(t, expected, expiration)
			assert.Equal(t, tk.ExpiresOn, got.ExpiresOn)
			assert.Equal(t, lifetime < 5*time.Minute, shouldRefresh(got, state))
		})
	}
}
