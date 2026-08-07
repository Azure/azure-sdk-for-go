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

func TestStripCredentialsOnCrossHostRedirect(t *testing.T) {
	newReq := func(rawURL string) *http.Request {
		u, _ := url.Parse(rawURL)
		r := &http.Request{URL: u, Header: http.Header{}}
		r.Header.Set("Authorization", "Bearer secret")
		r.Header.Set("aeg-sas-key", "secret-key")
		r.Header.Set("aeg-sas-token", "secret-token")
		r.Header.Set("Content-Type", "application/json")
		return r
	}

	t.Run("cross-host strips credential headers", func(t *testing.T) {
		req := newReq("https://attacker.example/collect")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		require.NoError(t, stripCredentialsOnCrossHostRedirect(nil)(req, via))
		require.Empty(t, req.Header.Get("Authorization"))
		require.Empty(t, req.Header.Get("aeg-sas-key"))
		require.Empty(t, req.Header.Get("aeg-sas-token"))
		// non-credential headers are preserved
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))
	})

	t.Run("same-host preserves credential headers", func(t *testing.T) {
		req := newReq("https://topic.eventgrid.azure.net/api/events?redirected=1")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		require.NoError(t, stripCredentialsOnCrossHostRedirect(nil)(req, via))
		require.Equal(t, "secret-key", req.Header.Get("aeg-sas-key"))
		require.Equal(t, "secret-token", req.Header.Get("aeg-sas-token"))
	})

	t.Run("chains prior CheckRedirect", func(t *testing.T) {
		sentinel := errors.New("prior called")
		req := newReq("https://attacker.example/collect")
		via := []*http.Request{newReq("https://topic.eventgrid.azure.net/api/events")}
		err := stripCredentialsOnCrossHostRedirect(func(*http.Request, []*http.Request) error {
			return sentinel
		})(req, via)
		require.ErrorIs(t, err, sentinel)
		require.Empty(t, req.Header.Get("aeg-sas-key")) // still stripped before delegating
	})

	t.Run("default redirect limit enforced", func(t *testing.T) {
		req := newReq("https://topic.eventgrid.azure.net/a")
		via := make([]*http.Request, 10)
		for i := range via {
			via[i] = newReq("https://topic.eventgrid.azure.net/api/events")
		}
		require.Error(t, stripCredentialsOnCrossHostRedirect(nil)(req, via))
	})
}

// TestPublisher_DoesNotLeakCredentialsOnCrossHostRedirect exercises the full
// client: a "trusted" topic host answers the publish with a 307 redirect to a
// different "attacker" host and we assert the credential never reaches it.
func TestPublisher_DoesNotLeakCredentialsOnCrossHostRedirect(t *testing.T) {
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

	// Transport maps the fake hostnames to loopback and trusts the test cert.
	// It has no CheckRedirect, mimicking a plain client; the fix must add one.
	transport := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				_, port, _ := net.SplitHostPort(addr)
				return (&net.Dialer{}).DialContext(ctx, network, "127.0.0.1:"+port)
			},
			TLSClientConfig: &tls.Config{RootCAs: pool},
		},
	}
	opts := &ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}

	evt, err := messaging.NewCloudEvent("src", "Test.Event", map[string]string{"hello": "world"}, nil)
	require.NoError(t, err)
	events := []messaging.CloudEvent{evt}

	cases := []struct {
		name   string
		header string
		secret string
		client func() (*Client, error)
	}{
		{"key", "aeg-sas-key", "supersecret-key", func() (*Client, error) {
			return NewClientWithSharedKeyCredential(endpoint, azcore.NewKeyCredential("supersecret-key"), opts)
		}},
		{"sas", "aeg-sas-token", "supersecret-sas", func() (*Client, error) {
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
			_, _ = client.PublishCloudEvents(context.Background(), events, nil)

			// credential was sent to the trusted host...
			sentToTrusted := false
			for _, h := range captured["trusted"] {
				if h.Get(tc.header) != "" {
					sentToTrusted = true
				}
			}
			require.True(t, sentToTrusted, "expected credential to be sent to the trusted host")

			// ...the redirect was followed to the attacker host...
			require.NotEmpty(t, captured["attacker"], "expected the redirect to reach the attacker host")

			// ...but the credential must NOT have been forwarded.
			for _, h := range captured["attacker"] {
				require.Empty(t, h.Get(tc.header), "credential header %q leaked to attacker host", tc.header)
			}
		})
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
