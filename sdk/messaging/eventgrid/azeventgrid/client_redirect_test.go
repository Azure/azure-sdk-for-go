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

func TestNoFollowRedirect(t *testing.T) {
	newReq := func(rawURL string) *http.Request {
		u, _ := url.Parse(rawURL)
		return &http.Request{URL: u, Header: http.Header{}}
	}

	t.Run("cross-host redirect is not followed", func(t *testing.T) {
		req := newReq("https://attacker.example/collect")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		require.ErrorIs(t, noFollowRedirect(req, via), http.ErrUseLastResponse)
	})

	t.Run("same-host redirect is not followed", func(t *testing.T) {
		req := newReq("https://topic.eventgrid.azure.net/api/events?redirected=1")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		require.ErrorIs(t, noFollowRedirect(req, via), http.ErrUseLastResponse)
	})
}

// TestPublisher_DoesNotFollowCrossHostRedirect exercises the full client: a
// "trusted" topic host answers the publish with a 307 redirect to a different
// "attacker" host, and we assert the attacker host is never contacted at all
// (so neither the credential nor the event payload leaks).
func TestPublisher_DoesNotFollowCrossHostRedirect(t *testing.T) {
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

	opts := &ClientOptions{ClientOptions: azcore.ClientOptions{Transport: loopbackTestClient(pool)}}

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
			// The publish must fail rather than follow the redirect.
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
			require.Empty(t, captured["attacker"], "attacker host must not be contacted on a redirect")
		})
	}
}

// TestPublisher_DoesNotFollowSameHostRedirect verifies that even a same-host
// redirect is not followed (matching the .NET no-auto-redirect posture); the
// publish surfaces the 3xx as an error.
func TestPublisher_DoesNotFollowSameHostRedirect(t *testing.T) {
	const host = "trusted-topic.eventgrid.example"

	cert := selfSignedCert(t, host)
	pool := x509.NewCertPool()
	pool.AddCert(cert.Leaf)

	var redirectTargetReached bool

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("redirected") == "" {
			http.Redirect(w, r, r.URL.Path+"?redirected=1", http.StatusTemporaryRedirect)
			return
		}
		redirectTargetReached = true
		w.WriteHeader(http.StatusOK)
	}))
	srv.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	srv.StartTLS()
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	endpoint := fmt.Sprintf("https://%s:%s", host, u.Port())

	opts := &ClientOptions{ClientOptions: azcore.ClientOptions{Transport: loopbackTestClient(pool)}}
	client, err := NewClientWithSharedKeyCredential(endpoint, azcore.NewKeyCredential("supersecret-key"), opts)
	require.NoError(t, err)

	evt, err := messaging.NewCloudEvent("src", "Test.Event", map[string]string{"hello": "world"}, nil)
	require.NoError(t, err)

	_, err = client.PublishCloudEvents(context.Background(), []messaging.CloudEvent{evt}, nil)
	require.Error(t, err)
	require.False(t, redirectTargetReached, "the redirect target must not be reached")
}

// TestWithRedirectProtection_DoesNotMutateSharedOptions verifies the caller's
// ClientOptions (which azcore permits sharing across constructors) is not
// mutated, and that each call produces an independent transport.
func TestWithRedirectProtection_DoesNotMutateSharedOptions(t *testing.T) {
	shared := &ClientOptions{}

	first := withRedirectProtection(shared)
	require.Nil(t, shared.Transport, "caller's options must not be mutated")
	require.NotNil(t, first.Transport, "returned options must have the protective transport")

	second := withRedirectProtection(shared)
	require.Nil(t, shared.Transport, "caller's options must still not be mutated")
	require.NotSame(t, first, second, "each call returns an independent options copy")
	require.NotSame(t, first.Transport, second.Transport, "each call installs an independent client")

	// nil options is handled without panicking.
	require.NotNil(t, withRedirectProtection(nil).Transport)
}

// TestNewDefaultTransport verifies the installed default transport mirrors
// azcore's tuned settings rather than falling back to http.DefaultTransport.
func TestNewDefaultTransport(t *testing.T) {
	tr := newDefaultTransport()
	require.NotSame(t, http.DefaultTransport, tr, "must not reuse the process-global default transport")
	require.True(t, tr.ForceAttemptHTTP2)
	require.Equal(t, 100, tr.MaxIdleConns)
	require.Equal(t, 10, tr.MaxIdleConnsPerHost)
	require.Equal(t, 90*time.Second, tr.IdleConnTimeout)
	require.Equal(t, 10*time.Second, tr.TLSHandshakeTimeout)
	require.Equal(t, time.Second, tr.ExpectContinueTimeout)
	require.NotNil(t, tr.TLSClientConfig)
	require.Equal(t, uint16(tls.VersionTLS12), tr.TLSClientConfig.MinVersion)
	require.NotNil(t, tr.Proxy)
	require.NotNil(t, tr.DialContext)
}

// loopbackTestClient returns an *http.Client whose transport resolves the test
// hostnames to loopback and trusts the test CA. It intentionally sets no
// CheckRedirect, mimicking a plain caller-supplied client; the fix under test
// must disable redirect following.
func loopbackTestClient(pool *x509.CertPool) *http.Client {
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
