//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package azeventgrid

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/stretchr/testify/require"
)

func TestBlockCrossHostRedirect(t *testing.T) {
	newReq := func(rawURL string) *http.Request {
		u, _ := url.Parse(rawURL)
		r := &http.Request{URL: u, Header: http.Header{}}
		r.Header.Set("Authorization", "auth-secret")
		r.Header.Set("aeg-sas-key", "secret-key")
		r.Header.Set("aeg-sas-token", "secret-token")
		r.Header.Set("Content-Type", "application/json")
		return r
	}

	t.Run("cross-host redirect is blocked", func(t *testing.T) {
		req := newReq("https://attacker.example/collect")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		err := blockCrossHostRedirect(nil)(req, via)
		require.ErrorIs(t, err, errCrossHostRedirect)
	})

	t.Run("cross-host to a sibling subdomain is blocked", func(t *testing.T) {
		req := newReq("https://attacker.eventgrid.azure.net/collect")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		err := blockCrossHostRedirect(nil)(req, via)
		require.ErrorIs(t, err, errCrossHostRedirect)
	})

	t.Run("same-host redirect is allowed and keeps credential headers", func(t *testing.T) {
		req := newReq("https://topic.eventgrid.azure.net/api/events?redirected=1")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		require.NoError(t, blockCrossHostRedirect(nil)(req, via))
		// credential headers are untouched for a permitted same-host redirect
		require.Equal(t, "secret-key", req.Header.Get("aeg-sas-key"))
		require.Equal(t, "secret-token", req.Header.Get("aeg-sas-token"))
		require.Equal(t, "auth-secret", req.Header.Get("Authorization"))
	})

	t.Run("same-host match is case-insensitive", func(t *testing.T) {
		req := newReq("https://TOPIC.eventgrid.azure.net/api/events")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		require.NoError(t, blockCrossHostRedirect(nil)(req, via))
	})

	t.Run("prior CheckRedirect is chained for same-host redirects", func(t *testing.T) {
		sentinel := errors.New("prior called")
		req := newReq("https://topic.eventgrid.azure.net/api/events?redirected=1")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		err := blockCrossHostRedirect(func(*http.Request, []*http.Request) error {
			return sentinel
		})(req, via)
		require.ErrorIs(t, err, sentinel)
	})

	t.Run("cross-host block takes precedence over prior CheckRedirect", func(t *testing.T) {
		req := newReq("https://attacker.example/collect")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		priorCalled := false
		err := blockCrossHostRedirect(func(*http.Request, []*http.Request) error {
			priorCalled = true
			return nil
		})(req, via)
		require.ErrorIs(t, err, errCrossHostRedirect)
		require.False(t, priorCalled, "prior CheckRedirect must not run for a blocked cross-host redirect")
	})

	t.Run("default redirect limit enforced for same-host redirects", func(t *testing.T) {
		req := newReq("https://topic.eventgrid.azure.net/a")
		via := make([]*http.Request, maxDefaultRedirects)
		for i := range via {
			via[i] = newReq("https://topic.eventgrid.azure.net/api/events")
		}
		require.Error(t, blockCrossHostRedirect(nil)(req, via))
	})
}

// TestPublisher_BlocksCrossHostRedirect exercises the full client: a "trusted"
// topic host answers the publish with a 307 redirect to a different "attacker"
// host, and we assert the attacker host is never contacted at all (so neither
// the credential nor the event payload leaks).
func TestPublisher_BlocksCrossHostRedirect(t *testing.T) {
	const trustedHost = "trusted-topic.eventgrid.example"
	const attackerHost = "attacker-collector.example"

	cert := selfSignedCert(t, trustedHost, attackerHost)
	pool := x509.NewCertPool()
	pool.AddCert(cert.Leaf)

	captured := map[string][]http.Header{}
	var attackerURL string

	mkServer := func(role string) *httptest.Server {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			captured[role] = append(captured[role], r.Header.Clone())
			if role == "trusted" {
				http.Redirect(w, r, attackerURL, http.StatusTemporaryRedirect)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		s := httptest.NewUnstartedServer(h)
		s.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
		s.StartTLS()
		return s
	}
	trusted := mkServer("trusted")
	attacker := mkServer("attacker")
	defer trusted.Close()
	defer attacker.Close()

	portOf := func(s *httptest.Server) string { u, _ := url.Parse(s.URL); return u.Port() }
	attackerURL = fmt.Sprintf("https://%s:%s/collect", attackerHost, portOf(attacker))
	endpoint := fmt.Sprintf("https://%s:%s", trustedHost, portOf(trusted))

	transport := crossHostTestTransport(pool)
	opts := &ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}

	evt, err := messaging.NewCloudEvent("src", "Test.Event", map[string]string{"hello": "world"}, nil)
	require.NoError(t, err)
	events := []messaging.CloudEvent{evt}

	cases := []struct {
		name   string
		header string
		client func() (*Client, error)
	}{
		{"key", "aeg-sas-key", func() (*Client, error) {
			return NewClientWithSharedKeyCredential(endpoint, azcore.NewKeyCredential("supersecret-key"), opts)
		}},
		{"sas", "aeg-sas-token", func() (*Client, error) {
			return NewClientWithSAS(endpoint, azcore.NewSASCredential("supersecret-sas"), opts)
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for k := range captured {
				delete(captured, k)
			}
			client, err := tc.client()
			require.NoError(t, err)

			_, err = client.PublishCloudEvents(context.Background(), events, nil)
			// The publish must fail rather than follow the cross-host redirect.
			require.Error(t, err)

			// The credential was sent to the trusted (configured) host on the first request...
			sentToTrusted := false
			for _, h := range captured["trusted"] {
				if h.Get(tc.header) != "" {
					sentToTrusted = true
				}
			}
			require.True(t, sentToTrusted, "expected credential to be sent to the trusted host")

			// ...but the attacker host must NEVER be contacted at all.
			require.Empty(t, captured["attacker"], "attacker host must not be contacted on a cross-host redirect")
		})
	}
}

// TestPublisher_AllowsSameHostRedirect verifies that a same-host redirect is
// still followed (with the credential) and the publish succeeds.
func TestPublisher_AllowsSameHostRedirect(t *testing.T) {
	const host = "trusted-topic.eventgrid.example"

	cert := selfSignedCert(t, host)
	pool := x509.NewCertPool()
	pool.AddCert(cert.Leaf)

	var sawCredentialAfterRedirect bool
	var redirected bool

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("redirected") == "" {
			// Same-host redirect (e.g. path normalization / query addition).
			http.Redirect(w, r, r.URL.Path+"?redirected=1", http.StatusTemporaryRedirect)
			return
		}
		redirected = true
		if r.Header.Get("aeg-sas-key") != "" {
			sawCredentialAfterRedirect = true
		}
		w.WriteHeader(http.StatusOK)
	}))
	srv.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	srv.StartTLS()
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	endpoint := fmt.Sprintf("https://%s:%s", host, u.Port())

	opts := &ClientOptions{ClientOptions: azcore.ClientOptions{Transport: crossHostTestTransport(pool)}}
	client, err := NewClientWithSharedKeyCredential(endpoint, azcore.NewKeyCredential("supersecret-key"), opts)
	require.NoError(t, err)

	evt, err := messaging.NewCloudEvent("src", "Test.Event", map[string]string{"hello": "world"}, nil)
	require.NoError(t, err)

	_, err = client.PublishCloudEvents(context.Background(), []messaging.CloudEvent{evt}, nil)
	require.NoError(t, err)
	require.True(t, redirected, "expected the same-host redirect to be followed")
	require.True(t, sawCredentialAfterRedirect, "expected the credential to be retained on a same-host redirect")
}

// crossHostTestTransport returns an *http.Client whose transport resolves the
// test hostnames to loopback and trusts the test CA. It intentionally sets no
// CheckRedirect, mimicking a plain caller-supplied client; the fix under test
// must add the cross-host protection.
func crossHostTestTransport(pool *x509.CertPool) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				_, port, _ := net.SplitHostPort(addr)
				return (&net.Dialer{}).DialContext(ctx, network, "127.0.0.1:"+port)
			},
			TLSClientConfig: &tls.Config{RootCAs: pool},
		},
	}
}

func selfSignedCert(t *testing.T, hosts ...string) tls.Certificate {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: hosts[0]},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     hosts,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)
	leaf, err := x509.ParseCertificate(der)
	require.NoError(t, err)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key, Leaf: leaf}
}
