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
	"io"
	"maps"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
	"weak"

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
	require.Same(t, client.Transport, cached.Transport)
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
	assert.Same(t, client.Transport, shared.Transport, "a new copy of the same certificate must share its connection pool")
	got, err = tr.TLSClientConfig.GetClientCertificate(&tls.CertificateRequestInfo{})
	require.NoError(t, err)
	assert.Same(t, &copied, got, "new handshakes must present the most recently acquired copy")
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	require.NoError(t, err)
	var rotations []string
	for i := range 10 {
		leaf.SerialNumber = big.NewInt(int64(i + 1))
		der, err := x509.CreateCertificate(rand.Reader, leaf, leaf, leaf.PublicKey, cert.PrivateKey)
		require.NoError(t, err)
		rotations = append(rotations, string(der))
		rotated, err := tp.clientFor(&tls.Certificate{Certificate: [][]byte{der}, PrivateKey: cert.PrivateKey})
		require.NoError(t, err)
		t.Cleanup(rotated.CloseIdleConnections)
		assert.NotSame(t, client.Transport, rotated.Transport, "rotation requires a separate connection pool")
		// keep the original certificate in use so eviction must choose among the rotations
		_, err = tp.clientFor(&copied)
		require.NoError(t, err)
	}
	cache := cachedMTLSTransports(baseTransport)
	require.Len(t, cache, mtlsTransportsPerBase)
	assert.Contains(t, cache, string(cert.Certificate[0]), "eviction must spare recently used certificates")
	for i, der := range rotations {
		_, ok := cache[der]
		assert.Equal(t, i >= len(rotations)-mtlsTransportsPerBase+1, ok, "eviction must remove the least recently used certificates first")
	}
}

func TestTransportPolicyMTLSSharedAcrossPipelines(t *testing.T) {
	var conns, authenticated atomic.Int32
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
			authenticated.Add(1)
		}
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			conns.Add(1)
		}
	}
	server.TLS = &tls.Config{ClientAuth: tls.RequireAnyClientCert}
	server.StartTLS()
	defer server.Close()
	cert := server.TLS.Certificates[0]
	base := server.Client()
	baseTransport, ok := base.Transport.(*http.Transport)
	require.True(t, ok)
	send := func(tp *transportPolicy) {
		t.Helper()
		req, err := NewRequest(context.Background(), http.MethodGet, server.URL)
		require.NoError(t, err)
		req.SetOperationValue(MTLSRequest{Certificate: &cert})
		resp, err := tp.Do(req)
		require.NoError(t, err)
		_, err = io.Copy(io.Discard, resp.Body)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
	}
	send(&transportPolicy{trans: base})
	send(&transportPolicy{trans: &http.Client{Transport: baseTransport, Timeout: time.Minute}})
	assert.EqualValues(t, 2, authenticated.Load())
	assert.EqualValues(t, 1, conns.Load(), "pipelines sharing a base transport must share its certificate-scoped pool")

	other := &http.Client{Transport: baseTransport.Clone()}
	send(&transportPolicy{trans: other})
	assert.EqualValues(t, 3, authenticated.Load())
	assert.EqualValues(t, 2, conns.Load(), "different base transports must not share connections")
	for _, b := range []*http.Transport{baseTransport, other.Transport.(*http.Transport)} {
		for _, d := range cachedMTLSTransports(b) {
			d.tr.CloseIdleConnections()
		}
	}
}

func TestTransportPolicyMTLSReleasesCollectedBase(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	cert := server.TLS.Certificates[0]
	server.Close()
	key := func() weak.Pointer[http.Transport] {
		base := &http.Transport{}
		_, err := (&transportPolicy{trans: &http.Client{Transport: base}}).clientFor(&cert)
		require.NoError(t, err)
		require.Len(t, cachedMTLSTransports(base), 1)
		return weak.Make(base)
	}()
	require.Eventually(t, func() bool {
		runtime.GC()
		mtlsTransports.mu.Lock()
		defer mtlsTransports.mu.Unlock()
		_, ok := mtlsTransports.bases[key]
		return !ok
	}, 10*time.Second, 10*time.Millisecond, "transports derived from a collected base transport must be released")
}

func TestTransportPolicyMTLSEvictedDuringRequest(t *testing.T) {
	before := mtlsIdleConnTimeout
	mtlsIdleConnTimeout = 50 * time.Millisecond
	t.Cleanup(func() { mtlsIdleConnTimeout = before })
	var closed atomic.Int32
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateClosed {
			closed.Add(1)
		}
	}
	server.StartTLS()
	defer server.Close()
	cert := server.TLS.Certificates[0]
	// httptest's client transport has no idle timeout, so without one on the derived transport,
	// a connection opened after eviction would stay open
	base := server.Client()
	baseTransport, ok := base.Transport.(*http.Transport)
	require.True(t, ok)
	require.Zero(t, baseTransport.IdleConnTimeout)
	client, err := (&transportPolicy{trans: base}).clientFor(&cert)
	require.NoError(t, err)

	// evict the selected transport before the request is sent
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	require.NoError(t, err)
	for i := range mtlsTransportsPerBase {
		leaf.SerialNumber = big.NewInt(int64(i + 1))
		der, err := x509.CreateCertificate(rand.Reader, leaf, leaf, leaf.PublicKey, cert.PrivateKey)
		require.NoError(t, err)
		_, err = mtlsTransportFor(baseTransport, &tls.Certificate{Certificate: [][]byte{der}, PrivateKey: cert.PrivateKey})
		require.NoError(t, err)
	}
	require.NotContains(t, cachedMTLSTransports(baseTransport), string(cert.Certificate[0]))

	resp, err := client.Get(server.URL)
	require.NoError(t, err)
	_, err = io.Copy(io.Discard, resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	require.Eventually(t, func() bool { return closed.Load() == 1 }, 5*time.Second, 10*time.Millisecond,
		"a connection an evicted transport opens must close once idle")
	for _, d := range cachedMTLSTransports(baseTransport) {
		d.tr.CloseIdleConnections()
	}
}

// cachedMTLSTransports returns a copy of the cached transports derived from base
func cachedMTLSTransports(base *http.Transport) map[string]*mtlsTransport {
	mtlsTransports.mu.Lock()
	defer mtlsTransports.mu.Unlock()
	return maps.Clone(mtlsTransports.bases[weak.Make(base)])
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
