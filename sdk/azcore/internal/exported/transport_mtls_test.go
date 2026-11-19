// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransportPolicyMTLSClient(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()
	cert := server.TLS.Certificates[0]
	base := server.Client()
	baseTransport, ok := base.Transport.(*http.Transport)
	require.True(t, ok)
	sessionCache := tls.NewLRUClientSessionCache(1)
	baseTransport.TLSClientConfig.ClientSessionCache = sessionCache
	tp := &transportPolicy{trans: base}
	client, err := tp.clientFor(&cert)
	require.NoError(t, err)
	t.Cleanup(client.CloseIdleConnections)
	cached, err := tp.clientFor(&cert)
	require.NoError(t, err)
	require.Same(t, client, cached)
	tr, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	got, err := tr.TLSClientConfig.GetClientCertificate(&tls.CertificateRequestInfo{AcceptableCAs: [][]byte{[]byte("unrelated CA")}})
	require.NoError(t, err)
	assert.Same(t, &cert, got, "CA hints must not filter out the binding certificate")
	assert.Equal(t, tls.RenegotiateOnceAsClient, tr.TLSClientConfig.Renegotiation)
	assert.Equal(t, uint16(tls.VersionTLS12), tr.TLSClientConfig.MaxVersion)
	assert.Nil(t, tr.TLSClientConfig.ClientSessionCache, "session resumption must not cross certificate boundaries")
	assert.Same(t, sessionCache, baseTransport.TLSClientConfig.ClientSessionCache)
	assert.Error(t, client.CheckRedirect(nil, nil))
	copied := cert
	shared, err := tp.clientFor(&copied)
	require.NoError(t, err)
	assert.Same(t, client, shared, "a new copy of the same certificate must share its connection pool")
	got, err = tr.TLSClientConfig.GetClientCertificate(&tls.CertificateRequestInfo{})
	require.NoError(t, err)
	assert.Same(t, &copied, got, "new handshakes must present the most recently acquired copy")
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	require.NoError(t, err)
	for i := range 10 {
		leaf.SerialNumber = big.NewInt(int64(i + 1))
		der, err := x509.CreateCertificate(rand.Reader, leaf, leaf, leaf.PublicKey, cert.PrivateKey)
		require.NoError(t, err)
		rotated, err := tp.clientFor(&tls.Certificate{Certificate: [][]byte{der}, PrivateKey: cert.PrivateKey})
		require.NoError(t, err)
		t.Cleanup(rotated.CloseIdleConnections)
		assert.NotSame(t, client, rotated, "rotation requires a separate connection pool")
	}
	assert.Len(t, tp.clients, 8)
}

func TestTransportPolicyMTLSUnsupportedConfiguration(t *testing.T) {
	cert := &tls.Certificate{Certificate: [][]byte{{1}}, PrivateKey: testSigner{}}
	srv, close := mock.NewServer()
	defer close()
	for _, tc := range []struct {
		name    string
		base    Transporter
		message string
	}{
		{name: "custom transporter", base: srv, message: "*http.Client"},
		{name: "custom round tripper", base: &http.Client{Transport: unexpectedRoundTripper{}}, message: "*http.Transport"},
		{name: "TLS dialer", base: &http.Client{Transport: &http.Transport{
			DialTLSContext: func(context.Context, string, string) (net.Conn, error) { return nil, errors.New("unexpected call") },
		}}, message: "custom TLS dialing"},
		{name: "TLS 1.3 only", base: &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS13}}}, message: "TLS 1.2"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := (&transportPolicy{trans: tc.base}).clientFor(cert)
			require.ErrorContains(t, err, tc.message)
		})
	}
	req, err := NewRequest(context.Background(), http.MethodGet, "http://example.com")
	require.NoError(t, err)
	req.SetOperationValue(MTLSRequest{Certificate: cert})
	_, err = (&transportPolicy{trans: srv}).Do(req)
	require.ErrorContains(t, err, "HTTPS", "bound requests must require HTTPS independently of the client preference")
	require.Zero(t, srv.Requests())
	_, err = (&transportPolicy{}).clientFor(&tls.Certificate{})
	require.ErrorContains(t, err, "private key")
	_, err = (&transportPolicy{}).clientFor(nil)
	require.ErrorContains(t, err, "private key")
	_, err = (&transportPolicy{}).clientFor(&tls.Certificate{Certificate: [][]byte{{1}}, PrivateKey: struct{}{}})
	require.ErrorContains(t, err, "crypto.Signer")
}

func TestTransportPolicyMTLSRedirect(t *testing.T) {
	var redirected atomic.Bool
	target := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { redirected.Store(true) }))
	defer target.Close()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	}))
	defer server.Close()
	cert := server.TLS.Certificates[0]
	client, err := (&transportPolicy{trans: server.Client()}).clientFor(&cert)
	require.NoError(t, err)
	defer client.CloseIdleConnections()
	resp, err := client.Get(server.URL)
	require.ErrorContains(t, err, "redirects")
	if resp != nil {
		require.NoError(t, resp.Body.Close())
	}
	assert.False(t, redirected.Load())
}

func TestTransportPolicyMTLSRespectsServerTrust(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()
	cert := server.TLS.Certificates[0]
	client, err := (&transportPolicy{trans: &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, RootCAs: x509.NewCertPool()},
	}}}).clientFor(&cert)
	require.NoError(t, err)
	defer client.CloseIdleConnections()
	_, err = client.Get(server.URL)
	var verificationError *tls.CertificateVerificationError
	require.ErrorAs(t, err, &verificationError)
}

// testSigner satisfies crypto.Signer for tests that fail before a handshake would call it
type testSigner struct{ crypto.Signer }

type unexpectedRoundTripper struct{}

func (unexpectedRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("unexpected request")
}
